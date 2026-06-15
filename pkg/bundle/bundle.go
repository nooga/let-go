/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

// Package bundle implements the standalone-binary trailer format and resource
// collection used by `lg -b`. A standalone binary is the lg executable with a
// compiled bytecode payload (and optionally a resource archive) appended,
// followed by a fixed trailer the runtime reads back at startup.
//
// Trailers appended to standalone binaries:
//
//	Legacy (no resources): [lgb data][8-byte lgbSize][4-byte 'LGBX']  (12-byte trailer)
//	v2    (resources):     [lgb data][resource archive][8-byte lgbSize][8-byte resSize][4-byte 'LGB2']  (20-byte trailer)
//
// A v2 trailer is written only when the bundle carries resources; resource-less
// bundles keep the byte-identical legacy trailer. Readers recognize both.
package bundle

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/nooga/let-go/pkg/rt"
)

var bundleMagic = [4]byte{'L', 'G', 'B', 'X'}
var bundleMagicV2 = [4]byte{'L', 'G', 'B', '2'}

// bundleKind classifies a standalone binary's appended trailer.
type bundleKind int

const (
	bundleNone   bundleKind = iota // no recognized trailer (a plain, non-bundled binary)
	bundleLegacy                   // 12-byte 'LGBX' trailer (lgb only)
	bundleV2                       // 20-byte 'LGB2' trailer (lgb + resource archive)
)

// trailerLen returns the on-disk size of the trailer for this kind.
func (k bundleKind) trailerLen() int64 {
	switch k {
	case bundleV2:
		return 20
	case bundleLegacy:
		return 12
	default:
		return 0
	}
}

// parseBundleTrailer reads and validates the trailer appended to f, the single
// place that discriminates the LGB2/LGBX formats. It returns bundleNone with a
// nil error when f carries no recognized trailer (a normal, non-bundled
// binary). For a recognized trailer it validates that the claimed payload fits
// within the file and returns a "corrupt bundle" error otherwise — so callers
// never seek to a bogus offset or allocate a garbage-sized slice.
func parseBundleTrailer(f *os.File) (kind bundleKind, lgbSize, resSize uint64, err error) {
	fi, err := f.Stat()
	if err != nil {
		return bundleNone, 0, 0, err
	}
	total := fi.Size()
	if total < bundleLegacy.trailerLen() {
		return bundleNone, 0, 0, nil
	}

	// Discriminate by the trailing 4-byte magic.
	if _, err := f.Seek(-4, io.SeekEnd); err != nil {
		return bundleNone, 0, 0, err
	}
	var magic [4]byte
	if _, err := io.ReadFull(f, magic[:]); err != nil {
		return bundleNone, 0, 0, err
	}

	switch magic {
	case bundleMagicV2:
		if total < bundleV2.trailerLen() {
			return bundleNone, 0, 0, nil
		}
		if _, err := f.Seek(-bundleV2.trailerLen(), io.SeekEnd); err != nil {
			return bundleNone, 0, 0, err
		}
		var tr [20]byte
		if _, err := io.ReadFull(f, tr[:]); err != nil {
			return bundleNone, 0, 0, err
		}
		kind = bundleV2
		lgbSize = binary.LittleEndian.Uint64(tr[0:8])
		resSize = binary.LittleEndian.Uint64(tr[8:16])
	case bundleMagic:
		if _, err := f.Seek(-bundleLegacy.trailerLen(), io.SeekEnd); err != nil {
			return bundleNone, 0, 0, err
		}
		var tr [12]byte
		if _, err := io.ReadFull(f, tr[:]); err != nil {
			return bundleNone, 0, 0, err
		}
		kind = bundleLegacy
		lgbSize = binary.LittleEndian.Uint64(tr[0:8])
	default:
		return bundleNone, 0, 0, nil
	}

	// Size guard: the claimed payload plus trailer must fit within the file.
	// A crafted size that fails this can no longer reach a make([]byte, lgbSize)
	// or a negative seek offset.
	if !payloadFitsFile(lgbSize, resSize, kind.trailerLen(), total) {
		return bundleNone, 0, 0, fmt.Errorf("corrupt bundle: payload size exceeds file size")
	}
	return kind, lgbSize, resSize, nil
}

// payloadFitsFile reports whether a payload of lgbSize + resSize bytes plus a
// trailerLen-byte trailer fits within a total-byte file. It subtracts step by
// step instead of summing, so the test can't overflow uint64 even on a huge
// (e.g. sparse) file where the individual sizes are valid but their sum wraps.
func payloadFitsFile(lgbSize, resSize uint64, trailerLen, total int64) bool {
	if total < 0 || trailerLen < 0 {
		return false
	}
	avail := uint64(total)
	if lgbSize > avail {
		return false
	}
	avail -= lgbSize
	if resSize > avail {
		return false
	}
	avail -= resSize
	return uint64(trailerLen) <= avail
}

// ReadBundled extracts the appended payload from the file at path. It
// recognizes both the v2 (resource-carrying) trailer and the legacy trailer;
// res is nil for legacy bundles. Returns nil, nil for a non-bundle or corrupt
// file.
func ReadBundled(path string) (lgb []byte, res []byte) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil
	}
	defer f.Close()

	kind, lgbSize, resSize, err := parseBundleTrailer(f)
	if err != nil || kind == bundleNone {
		return nil, nil // not a bundle, or a corrupt one — behave as no payload
	}

	// Payload layout: [lgb][resArc][trailer]. Sizes are validated to fit the
	// file, so the seek offset is a valid negative and make() can't overrun.
	if _, err := f.Seek(-kind.trailerLen()-int64(resSize)-int64(lgbSize), io.SeekEnd); err != nil {
		return nil, nil
	}
	lgb = make([]byte, lgbSize)
	if _, err := io.ReadFull(f, lgb); err != nil {
		return nil, nil
	}
	if resSize > 0 {
		res = make([]byte, resSize)
		if _, err := io.ReadFull(f, res); err != nil {
			return nil, nil
		}
	}
	return lgb, res
}

// BaseBinarySize returns the size of the lg binary without any appended bundle,
// so re-bundling can strip an existing payload. A corrupt trailer surfaces as
// an error rather than a silently wrong size.
func BaseBinarySize(f *os.File) (int64, error) {
	kind, lgbSize, resSize, err := parseBundleTrailer(f)
	if err != nil {
		return 0, err
	}
	fi, err := f.Stat()
	if err != nil {
		return 0, err
	}
	total := fi.Size()
	if kind == bundleNone {
		return total, nil
	}
	// Sizes are validated to fit the file, so this can't go negative.
	return total - int64(lgbSize) - int64(resSize) - kind.trailerLen(), nil
}

// AppendTrailer writes the bundle payload to out: the lgb bytes, an optional
// resource archive, and the appropriate trailer. When resArc is non-empty it
// emits the v2 trailer; otherwise the byte-identical legacy trailer. This is the
// only writer of the trailer format (parseBundleTrailer is the only reader).
func AppendTrailer(out io.Writer, lgbData, resArc []byte) error {
	if _, err := out.Write(lgbData); err != nil {
		return err
	}

	// Embed the resource archive (if any) between the lgb and the trailer.
	if len(resArc) > 0 {
		if _, err := out.Write(resArc); err != nil {
			return err
		}
		// v2 trailer: [8-byte lgbSize][8-byte resSize][4-byte 'LGB2']
		var tr [20]byte
		binary.LittleEndian.PutUint64(tr[0:8], uint64(len(lgbData)))
		binary.LittleEndian.PutUint64(tr[8:16], uint64(len(resArc)))
		copy(tr[16:], bundleMagicV2[:])
		_, err := out.Write(tr[:])
		return err
	}

	// No resources: legacy trailer [8-byte lgbSize][4-byte 'LGBX'].
	var footer [12]byte
	binary.LittleEndian.PutUint64(footer[:8], uint64(len(lgbData)))
	copy(footer[8:], bundleMagic[:])
	_, err := out.Write(footer[:])
	return err
}

// CollectResources returns a map of slash-relative path → file bytes for every
// regular file reachable under the resource roots. It follows symlinks
// everywhere the dev FS provider does (which resolves names with os.Stat) — to
// symlinked roots, symlinked sub-directories, and symlinked files — so a -b
// bundle embeds exactly what dev lookup would find. Symlink cycles are guarded
// by real-path. When the same relative path exists under multiple roots the
// first root wins (matching the provider's precedence). The bundle's own output
// file (excludeAbs) is never embedded, so a dst inside a resource root can't
// embed itself.
func CollectResources(roots []string, excludeAbs string) (map[string][]byte, error) {
	files := map[string][]byte{}
	var exclude os.FileInfo
	if excludeAbs != "" {
		// os.Stat (not Lstat) so the comparison is symlink/hardlink robust.
		if fi, err := os.Stat(excludeAbs); err == nil {
			exclude = fi
		}
	}
	for _, root := range roots {
		realRoot, err := filepath.EvalSymlinks(root)
		if err != nil {
			continue // missing or dangling root
		}
		absRoot, err := filepath.Abs(realRoot)
		if err != nil {
			continue
		}
		if info, err := os.Stat(absRoot); err != nil || !info.IsDir() {
			continue
		}
		ancestors := map[string]bool{}
		if err := collectResourceDir(absRoot, absRoot, "", files, exclude, ancestors); err != nil {
			return nil, err
		}
	}
	return files, nil
}

// collectResourceDir recursively collects regular files under dir (an absolute,
// symlink-resolved directory inside rootReal) into files, keyed by relPrefix +
// entry name (slash-separated). Symlinks to files and sub-directories are
// followed, but only when they resolve to a path inside rootReal — a symlink
// escaping the root (e.g. `up -> ..`) is skipped, so a bundle never embeds
// files outside the declared resource tree. ancestors holds the resolved
// directory paths on the current descent path; revisiting one is a cycle and is
// skipped (while distinct names aliasing the same in-root dir are both kept).
func collectResourceDir(dir, rootReal, relPrefix string, files map[string][]byte, exclude os.FileInfo, ancestors map[string]bool) error {
	if ancestors[dir] {
		return nil // symlink cycle
	}
	ancestors[dir] = true
	defer delete(ancestors, dir)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		p := filepath.Join(dir, e.Name())
		key := e.Name()
		if relPrefix != "" {
			key = relPrefix + "/" + e.Name()
		}
		real, err := filepath.EvalSymlinks(p)
		if err != nil {
			continue // dangling symlink, race — skip
		}
		// Containment: never follow a symlink that escapes the resource root.
		if !rt.WithinRoot(rootReal, real) {
			continue
		}
		info, err := os.Stat(real)
		if err != nil {
			continue
		}
		if info.IsDir() {
			if err := collectResourceDir(real, rootReal, key, files, exclude, ancestors); err != nil {
				return err
			}
			continue
		}
		if !info.Mode().IsRegular() {
			continue // FIFO, device, socket
		}
		if exclude != nil && os.SameFile(info, exclude) {
			continue // never embed the bundle's own output file
		}
		if _, exists := files[key]; exists {
			continue // first root wins
		}
		data, err := os.ReadFile(real)
		if err != nil {
			return err
		}
		files[key] = data
	}
	return nil
}
