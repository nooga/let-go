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
// origin-wide localStorage keyspace.
type HostStorage struct {
	prefix string
}

func NewHostStorage(storeID string) *HostStorage {
	if storeID == "" {
		storeID = "default"
	}
	return &HostStorage{prefix: "let-go:" + storeID + ":"}
}

func (s *HostStorage) physicalKey(key string) string {
	return s.prefix + key
}

func (s *HostStorage) Get(key string) (value string, ok bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			value, ok, err = "", false, fmt.Errorf("storage/get: %v", r)
		}
	}()
	v := js.Global().Get("localStorage").Call("getItem", s.physicalKey(key))
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
	js.Global().Get("localStorage").Call("setItem", s.physicalKey(key), value)
	return nil
}

func (s *HostStorage) Remove(key string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("storage/remove: %v", r)
		}
	}()
	js.Global().Get("localStorage").Call("removeItem", s.physicalKey(key))
	return nil
}

func (s *HostStorage) Keys(prefix string) (keys []string, err error) {
	defer func() {
		if r := recover(); r != nil {
			keys, err = nil, fmt.Errorf("storage/keys: %v", r)
		}
	}()
	ls := js.Global().Get("localStorage")
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
