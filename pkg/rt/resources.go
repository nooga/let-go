/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 *
 * Resource subsystem: locating non-source resource bytes (templates, static
 * web assets, data files) for io/resource. In dev these come from the
 * filesystem (the -resource-paths roots); in a -b standalone binary they come
 * from a gzip archive appended to the executable.
 */

package rt

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

// ResourceProvider locates resource bytes by a normalized, slash-separated
// name. io/resource consults the currently-installed provider. There are two
// implementations: a filesystem provider (dev, backed by -resource-paths
// roots) and an embedded provider (a -b standalone binary, backed by the
// archive appended to the executable).
type ResourceProvider interface {
	// Open returns a reader over the named resource, or ok=false if absent.
	Open(name string) (io.ReadCloser, bool)
}

// resourceProvider is the process-wide provider, installed by the runtime
// entry point. Mirrors the nsLoader seam in lang.go.
var resourceProvider ResourceProvider

// SetResourceProvider installs the provider consulted by io/resource.
func SetResourceProvider(p ResourceProvider) { resourceProvider = p }

// GetResourceProvider returns the currently-installed provider (may be nil).
func GetResourceProvider() ResourceProvider { return resourceProvider }

// FSResourceProvider resolves resources against an ordered list of filesystem
// roots; the first root containing the name wins.
type FSResourceProvider struct {
	roots []string
}

// NewFSResourceProvider builds a provider over the given resource roots.
func NewFSResourceProvider(roots []string) *FSResourceProvider {
	return &FSResourceProvider{roots: roots}
}

// Open implements ResourceProvider. Errors and non-regular files (e.g.
// directories) are treated as "not found", so a lookup never panics on a
// permission race or a name that happens to be a directory.
func (p *FSResourceProvider) Open(name string) (io.ReadCloser, bool) {
	clean, ok := NormalizeResourceName(name)
	if !ok {
		return nil, false
	}
	rel := filepath.FromSlash(clean)
	for _, root := range p.roots {
		realRoot, ok := resolveReal(root)
		if !ok {
			continue // missing or dangling root
		}
		real, ok := resolveReal(filepath.Join(realRoot, rel))
		if !ok {
			continue // missing or dangling candidate
		}
		// Containment: a symlink component may resolve outside the root
		// (e.g. `up -> ..`), which would bypass the `..` rejection in
		// NormalizeResourceName. Only serve candidates that stay inside the
		// resolved root.
		if !WithinRoot(realRoot, real) {
			continue
		}
		// Stat before Open: only expose regular files. Directories, FIFOs,
		// devices, and sockets are treated as not found. This must happen
		// *before* Open, because opening a FIFO blocks until a writer appears —
		// so an Open-then-check approach would hang the lookup.
		info, err := os.Stat(real)
		if err != nil || !info.Mode().IsRegular() {
			continue
		}
		f, err := os.Open(real)
		if err != nil {
			continue
		}
		return f, true
	}
	return nil, false
}

// resolveReal returns the absolute, symlink-resolved form of p, or ok=false if
// p does not exist (or cannot be resolved). Used to enforce root containment.
func resolveReal(p string) (string, bool) {
	abs, err := filepath.Abs(p)
	if err != nil {
		return "", false
	}
	real, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return "", false
	}
	return real, true
}

// WithinRoot reports whether the resolved path p is the resolved directory
// root or lies inside it. Both arguments must already be absolute and
// symlink-resolved (see resolveReal).
func WithinRoot(root, p string) bool {
	if p == root {
		return true
	}
	if !strings.HasSuffix(root, string(os.PathSeparator)) {
		root += string(os.PathSeparator)
	}
	return strings.HasPrefix(p, root)
}

// EmbeddedResourceProvider serves resources from an in-memory map decoded from
// the archive appended to a standalone binary.
type EmbeddedResourceProvider struct {
	files map[string][]byte
}

// NewEmbeddedResourceProvider builds a provider over a decoded archive map.
func NewEmbeddedResourceProvider(files map[string][]byte) *EmbeddedResourceProvider {
	return &EmbeddedResourceProvider{files: files}
}

// Open implements ResourceProvider against the embedded map.
func (p *EmbeddedResourceProvider) Open(name string) (io.ReadCloser, bool) {
	clean, ok := NormalizeResourceName(name)
	if !ok {
		return nil, false
	}
	data, ok := p.files[clean]
	if !ok {
		return nil, false
	}
	return io.NopCloser(bytes.NewReader(data)), true
}

// resourceArchiveMagic prefixes the (pre-gzip) framed resource archive.
var resourceArchiveMagic = [4]byte{'L', 'G', 'R', 'A'}

// Resource is the handle returned by io/resource. It is reader-coercible (the
// io namespace extends IReadable for it), so it flows through io/slurp,
// io/reader, and io/line-seq. It re-opens through its bound provider on each
// read, so a handle can be read more than once.
type Resource struct {
	Name     string
	provider ResourceProvider
}

func (r *Resource) String() string { return fmt.Sprintf("#<resource %s>", r.Name) }

// Open returns a fresh reader over the resource, or ok=false if it has since
// disappeared (a dev-mode filesystem race).
func (r *Resource) Open() (io.ReadCloser, bool) { return r.provider.Open(r.Name) }

// NewResource builds a resource handle bound to a provider.
func NewResource(name string, p ResourceProvider) *Resource {
	return &Resource{Name: name, provider: p}
}

// NormalizeResourceName canonicalizes a resource lookup name into a
// slash-separated path relative to a resource root. It strips a leading "/"
// or "./", collapses "." and in-root ".." segments, and reports ok=false for
// empty names or names that try to escape the root via "..".
func NormalizeResourceName(name string) (string, bool) {
	if name == "" {
		return "", false
	}
	n := strings.ReplaceAll(name, "\\", "/")
	// Strip every leading slash (not just one) so "//x" and "/x" both reduce to
	// the relative key "x", then clean *relatively* — that way an escaping ".."
	// survives as a "../" prefix we can reject, instead of being silently
	// neutralized to root by an absolute clean.
	n = strings.TrimLeft(n, "/")
	cleaned := path.Clean(n)
	if cleaned == "." || cleaned == "" {
		return "", false
	}
	if cleaned == ".." || strings.HasPrefix(cleaned, "../") {
		return "", false
	}
	return cleaned, true
}

// EncodeResourceArchive serializes a name->bytes map into a gzip-compressed,
// self-describing archive. Entries are written in sorted name order so the
// output is deterministic for a given input.
func EncodeResourceArchive(files map[string][]byte) ([]byte, error) {
	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)

	var framed bytes.Buffer
	framed.Write(resourceArchiveMagic[:])
	var u32 [4]byte
	var u64 [8]byte
	binary.LittleEndian.PutUint32(u32[:], uint32(len(names)))
	framed.Write(u32[:])
	for _, name := range names {
		data := files[name]
		binary.LittleEndian.PutUint32(u32[:], uint32(len(name)))
		framed.Write(u32[:])
		framed.WriteString(name)
		binary.LittleEndian.PutUint64(u64[:], uint64(len(data)))
		framed.Write(u64[:])
		framed.Write(data)
	}

	var out bytes.Buffer
	gz := gzip.NewWriter(&out)
	if _, err := gz.Write(framed.Bytes()); err != nil {
		return nil, fmt.Errorf("gzip resource archive: %w", err)
	}
	if err := gz.Close(); err != nil {
		return nil, fmt.Errorf("gzip resource archive: %w", err)
	}
	return out.Bytes(), nil
}

// DecodeResourceArchive reverses EncodeResourceArchive.
func DecodeResourceArchive(blob []byte) (map[string][]byte, error) {
	gz, err := gzip.NewReader(bytes.NewReader(blob))
	if err != nil {
		return nil, fmt.Errorf("read resource archive: %w", err)
	}
	defer gz.Close()
	framed, err := io.ReadAll(gz)
	if err != nil {
		return nil, fmt.Errorf("read resource archive: %w", err)
	}

	r := bytes.NewReader(framed)
	var magic [4]byte
	if _, err := io.ReadFull(r, magic[:]); err != nil || magic != resourceArchiveMagic {
		return nil, fmt.Errorf("invalid resource archive header")
	}
	var u32 [4]byte
	var u64 [8]byte
	if _, err := io.ReadFull(r, u32[:]); err != nil {
		return nil, fmt.Errorf("read resource archive count: %w", err)
	}
	count := binary.LittleEndian.Uint32(u32[:])

	files := make(map[string][]byte, count)
	for i := uint32(0); i < count; i++ {
		if _, err := io.ReadFull(r, u32[:]); err != nil {
			return nil, fmt.Errorf("read resource name length: %w", err)
		}
		nameLen := binary.LittleEndian.Uint32(u32[:])
		name := make([]byte, nameLen)
		if _, err := io.ReadFull(r, name); err != nil {
			return nil, fmt.Errorf("read resource name: %w", err)
		}
		if _, err := io.ReadFull(r, u64[:]); err != nil {
			return nil, fmt.Errorf("read resource data length: %w", err)
		}
		dataLen := binary.LittleEndian.Uint64(u64[:])
		data := make([]byte, dataLen)
		if _, err := io.ReadFull(r, data); err != nil {
			return nil, fmt.Errorf("read resource data: %w", err)
		}
		files[string(name)] = data
	}
	return files, nil
}
