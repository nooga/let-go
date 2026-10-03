//go:build !lg_no_json

/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"encoding/json"
	"fmt"

	"github.com/nooga/let-go/pkg/vm"
)

func toValue(keywordize bool, i any) (vm.Value, error) {
	switch i := i.(type) {
	case string:
		return vm.String(i), nil
	case bool:
		return vm.Boolean(i), nil
	case float64:
		if i == float64(int64(i)) {
			return vm.Int(int64(i)), nil
		}
		return vm.Float(i), nil
	case nil:
		return vm.NIL, nil
	case []any:
		r := make([]vm.Value, len(i))
		for j := range i {
			v, e := toValue(keywordize, i[j])
			if e != nil {
				return vm.NIL, e
			}
			r[j] = v
		}
		return vm.NewArrayVector(r), nil
	case map[string]any:
		newmap := vm.EmptyPersistentMap
		for k, v := range i {
			ve, e := toValue(keywordize, v)
			if e != nil {
				return vm.NIL, e
			}
			if keywordize {
				newmap = newmap.Assoc(vm.Keyword(k), ve).(*vm.PersistentMap)
			} else {
				newmap = newmap.Assoc(vm.String(k), ve).(*vm.PersistentMap)
			}

		}
		return newmap, nil
	default:
		return vm.NIL, vm.NewExecutionError("invalid JSON value")
	}
}

func optionsHaveKeywordize(opts vm.Value) (bool, error) {
	o, ok := opts.(vm.Lookup)
	if !ok {
		return false, fmt.Errorf("read-json options are not Map")
	}
	return vm.IsTruthy(o.ValueAt(vm.Keyword("keywords?"))), nil
}

func init() { RegisterInstaller(installJSONNS) }

// nolint
func installJSONNS() {
	readJson, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 || len(vs) > 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}

		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("read-json expected String")
		}

		keywordize := false
		var err error
		if len(vs) == 2 {
			keywordize, err = optionsHaveKeywordize(vs[1])
			if err != nil {
				return vm.NIL, err
			}
		}

		var v any
		err = json.Unmarshal([]byte(s), &v)
		if err != nil {
			return vm.NIL, err
		}

		return toValue(keywordize, v)
	})

	writeJson, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		v, err := fromValue(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		s, err := json.Marshal(v)
		if err != nil {
			return vm.NIL, err
		}
		return vm.String(s), nil
	})

	if err != nil {
		panic("json NS init failed")
	}

	ns := vm.NewNamespace("json")

	ns.Def("read-json", readJson)
	ns.Def("write-json", writeJson)
	RegisterNS(ns)
}
