//go:build js && wasm

/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"
	"sort"
	"strings"
	"syscall/js"
)

// HostStorage binds the storage namespace to browser localStorage. Logical
// keys are scoped under a host-selected store id before they reach the
// origin-wide localStorage keyspace. The id is run through encodeStorageKey
// first — same as the native FileStorage root — so it can't contain the ':'
// separator and the (store-id, logical-key) -> physical-key mapping stays
// unambiguous when two bundles share an origin.
//
// Each method recovers a thrown localStorage access (Safari private mode,
// storage disabled, quota exceeded) into a returned error rather than a nil,
// honoring the Storage contract: a genuinely absent key reads back as nil,
// but an unavailable or failing backend surfaces as a guest exception
// instead of a silent miss or dropped write. A future refinement could probe
// availability once at construction so private mode fails once, clearly,
// rather than on every call.
type HostStorage struct {
	prefix string
}

func NewHostStorage(storeID string) *HostStorage {
	if storeID == "" {
		storeID = "default"
	}
	return &HostStorage{prefix: "let-go:" + encodeStorageKey(storeID) + ":"}
}

func (s *HostStorage) physicalKey(key string) string {
	return s.prefix + key
}

// localStorage returns the page's localStorage, or an error where there is
// none: a Web Worker, which is where lg -w runs the program (#1078). It is
// checked rather than left to the recover in each method because TinyGo cannot
// recover a syscall/js panic on wasm, so the call would trap the page.
func localStorage(op string) (js.Value, error) {
	ls := js.Global().Get("localStorage")
	if ls.IsUndefined() || ls.IsNull() {
		return js.Value{}, fmt.Errorf("%s: no localStorage here (a Web Worker has none)", op)
	}
	return ls, nil
}

func (s *HostStorage) Get(key string) (value string, ok bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			value, ok, err = "", false, fmt.Errorf("storage/get: %v", r)
		}
	}()
	ls, err := localStorage("storage/get")
	if err != nil {
		return "", false, err
	}
	v := ls.Call("getItem", s.physicalKey(key))
	if v.IsNull() || v.IsUndefined() {
		return "", false, nil
	}
	return v.String(), true, nil
}

func (s *HostStorage) Set(key, value string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("storage/set: %v", r)
		}
	}()
	ls, err := localStorage("storage/set")
	if err != nil {
		return err
	}
	ls.Call("setItem", s.physicalKey(key), value)
	return nil
}

func (s *HostStorage) Remove(key string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("storage/remove: %v", r)
		}
	}()
	ls, err := localStorage("storage/remove")
	if err != nil {
		return err
	}
	ls.Call("removeItem", s.physicalKey(key))
	return nil
}

func (s *HostStorage) Keys(prefix string) (keys []string, err error) {
	defer func() {
		if r := recover(); r != nil {
			keys, err = nil, fmt.Errorf("storage/keys: %v", r)
		}
	}()
	ls, err := localStorage("storage/keys")
	if err != nil {
		return nil, err
	}
	n := ls.Get("length").Int()
	out := make([]string, 0, n)
	for i := 0; i < n; i++ {
		k := ls.Call("key", i)
		if k.IsNull() || k.IsUndefined() {
			continue
		}
		physical := k.String()
		if !strings.HasPrefix(physical, s.prefix) {
			continue
		}
		logical := strings.TrimPrefix(physical, s.prefix)
		if strings.HasPrefix(logical, prefix) {
			out = append(out, logical)
		}
	}
	sort.Strings(out)
	return out, nil
}
