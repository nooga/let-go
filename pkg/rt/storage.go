/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/nooga/let-go/pkg/vm"
)

// Storage is the host seam for persistent string key/value storage.
// Guest code names logical keys; the host owns the physical backend and
// namespace. Values are strings so callers keep serialization policy.
//
// Error contract, uniform across backends and surfaced to guest code:
//   - Get of an absent key returns ok=false with a nil error, so storage/get
//     yields nil. Get returns a non-nil error only when the backend itself
//     fails (unreadable file, unavailable localStorage); that throws in the
//     guest. Absence and failure are deliberately distinguishable.
//   - Set and Remove return an error on backend failure (full disk, quota,
//     private-mode localStorage), which throws. They never silently drop a
//     write — losing a save with no signal is the failure mode to avoid.
type Storage interface {
	Get(key string) (value string, ok bool, err error)
	Set(key, value string) error
	Remove(key string) error
	Keys(prefix string) ([]string, error)
}

// MemoryStorage is the default inert-but-honest backend and the test/embedder
// convenience store. It is per-instance and does not persist.
type MemoryStorage struct {
	mu sync.RWMutex
	m  map[string]string
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{m: make(map[string]string)}
}

func (s *MemoryStorage) Get(key string) (string, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.m[key]
	return v, ok, nil
}

func (s *MemoryStorage) Set(key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[key] = value
	return nil
}

func (s *MemoryStorage) Remove(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, key)
	return nil
}

func (s *MemoryStorage) Keys(prefix string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]string, 0, len(s.m))
	for k := range s.m {
		if strings.HasPrefix(k, prefix) {
			out = append(out, k)
		}
	}
	sort.Strings(out)
	return out, nil
}

// FileStorage persists keys as files below one host-selected root. Logical
// keys are encoded, never used as raw filenames, so callers can use portable
// string keys without learning filesystem path rules.
//
// No in-process lock: native lg drives it from a single guest thread, and
// Set's atomic temp-then-rename keeps a concurrent reader from observing a
// torn value. Cross-process coordination is the host's responsibility.
type FileStorage struct {
	root string
}

func NewFileStorage(root string) (*FileStorage, error) {
	if root == "" {
		return nil, fmt.Errorf("storage: empty file storage root")
	}
	// The root is created lazily by Set, not here: installing the store is a
	// side effect of every lg run, so eager MkdirAll would litter the config
	// dir with an empty store dir for scripts that never touch storage. Read
	// paths (Get/Keys/Remove) already treat a missing root as empty.
	return &FileStorage{root: root}, nil
}

func NewDefaultFileStorage(storeID string) (*FileStorage, error) {
	if storeID == "" {
		storeID = "default"
	}
	base, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	return NewFileStorage(filepath.Join(base, "let-go", "storage", encodeStorageKey(storeID)))
}

func (s *FileStorage) pathFor(key string) string {
	return filepath.Join(s.root, encodeStorageKey(key))
}

func (s *FileStorage) Get(key string) (string, bool, error) {
	data, err := os.ReadFile(s.pathFor(key))
	if os.IsNotExist(err) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return string(data), true, nil
}

func (s *FileStorage) Set(key, value string) error {
	if err := os.MkdirAll(s.root, 0700); err != nil {
		return err
	}
	// Atomic replace: write to a temp file in the same directory, then rename
	// over the target. A crash or full disk mid-write leaves an orphan temp
	// file behind, never a truncated existing value — the failure mode that
	// matters for save data. The "tmp-" prefix keeps these out of Keys/Get,
	// which only see "k-"-encoded names.
	tmp, err := os.CreateTemp(s.root, "tmp-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	if _, err := tmp.Write([]byte(value)); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpName)
		return err
	}
	if err := os.Rename(tmpName, s.pathFor(key)); err != nil {
		os.Remove(tmpName)
		return err
	}
	return nil
}

func (s *FileStorage) Remove(key string) error {
	err := os.Remove(s.pathFor(key))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func (s *FileStorage) Keys(prefix string) ([]string, error) {
	entries, err := os.ReadDir(s.root)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		// decodeStorageKey rejects any name without the "k-" prefix, which is
		// also what filters out Set's "tmp-*" temp files (orphaned by a crash
		// between CreateTemp and Rename). Keep that prefix in sync with Set.
		k, ok := decodeStorageKey(e.Name())
		if !ok || !strings.HasPrefix(k, prefix) {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys, nil
}

func encodeStorageKey(key string) string {
	return "k-" + base64.RawURLEncoding.EncodeToString([]byte(key))
}

func decodeStorageKey(name string) (string, bool) {
	if !strings.HasPrefix(name, "k-") {
		return "", false
	}
	data, err := base64.RawURLEncoding.DecodeString(strings.TrimPrefix(name, "k-"))
	if err != nil {
		return "", false
	}
	return string(data), true
}

// fallbackStorage backs *storage* when no host binding resolves (the var is
// missing or not yet installed). A shared in-memory store keeps "no host
// bound" semantics identical to the default root binding installed in
// iort.go: writes are visible within the process, never silently dropped.
//
// This deliberately differs from the nop fallbacks on *emit*/*keys*, whose
// root defaults are themselves nops — matching their roots. *storage*'s root
// is a MemoryStorage, so its fallback matches that, not a nop.
var fallbackStorage = NewMemoryStorage()

func resolveStorageVar(ec *vm.ExecContext, varName string) Storage {
	ns := lookupNSCached(NameCoreNS)
	if ns == nil {
		return nil
	}
	v := ns.LookupLocal(vm.Symbol(varName))
	if v == nil {
		return nil
	}
	b, ok := ec.Deref(v).(*vm.Boxed)
	if !ok {
		return nil
	}
	if s, ok := b.Unbox().(Storage); ok {
		return s
	}
	return nil
}

func boundStorage(ec *vm.ExecContext) Storage {
	if s := resolveStorageVar(ec, "*storage*"); s != nil {
		return s
	}
	return fallbackStorage
}

func storageStringArg(name string, v vm.Value) (string, error) {
	s, ok := v.(vm.String)
	if !ok {
		return "", fmt.Errorf("storage/%s expected string key", name)
	}
	return string(s), nil
}

func init() { RegisterInstaller(installStorageNS) }

func installStorageNS() {
	ns := vm.NewNamespace("storage")

	getFn := vm.NewCtxNativeFn("get", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("storage/get expects 1 arg")
		}
		key, err := storageStringArg("get", vs[0])
		if err != nil {
			return vm.NIL, err
		}
		value, ok, err := boundStorage(ec).Get(key)
		if err != nil {
			return vm.NIL, err
		}
		if !ok {
			return vm.NIL, nil
		}
		return vm.String(value), nil
	})

	setFn := vm.NewCtxNativeFn("set", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("storage/set expects 2 args")
		}
		key, err := storageStringArg("set", vs[0])
		if err != nil {
			return vm.NIL, err
		}
		value, ok := vs[1].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("storage/set expected string value")
		}
		return vm.NIL, boundStorage(ec).Set(key, string(value))
	})

	removeFn := vm.NewCtxNativeFn("remove", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("storage/remove expects 1 arg")
		}
		key, err := storageStringArg("remove", vs[0])
		if err != nil {
			return vm.NIL, err
		}
		return vm.NIL, boundStorage(ec).Remove(key)
	})

	keysFn := vm.NewCtxNativeFn("keys", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		if len(vs) > 1 {
			return vm.NIL, fmt.Errorf("storage/keys expects 0 or 1 args")
		}
		prefix := ""
		if len(vs) == 1 {
			var err error
			prefix, err = storageStringArg("keys", vs[0])
			if err != nil {
				return vm.NIL, err
			}
		}
		keys, err := boundStorage(ec).Keys(prefix)
		if err != nil {
			return vm.NIL, err
		}
		vals := make([]vm.Value, len(keys))
		for i, key := range keys {
			vals[i] = vm.String(key)
		}
		return vm.NewArrayVector(vals), nil
	})

	ns.Def("get", getFn)
	ns.Def("set", setFn)
	ns.Def("remove", removeFn)
	ns.Def("keys", keysFn)
	RegisterNS(ns)
}
