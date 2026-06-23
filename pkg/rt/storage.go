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
type FileStorage struct {
	root string
}

func NewFileStorage(root string) (*FileStorage, error) {
	if root == "" {
		return nil, fmt.Errorf("storage: empty file storage root")
	}
	if err := os.MkdirAll(root, 0700); err != nil {
		return nil, err
	}
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
	return os.WriteFile(s.pathFor(key), []byte(value), 0600)
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

type nopStorage struct{}

func (nopStorage) Get(string) (string, bool, error) { return "", false, nil }
func (nopStorage) Set(string, string) error         { return nil }
func (nopStorage) Remove(string) error              { return nil }
func (nopStorage) Keys(string) ([]string, error)    { return nil, nil }

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
	return nopStorage{}
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
