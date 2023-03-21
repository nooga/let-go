/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"encoding/json"
	"fmt"

	"github.com/nooga/let-go/pkg/vm"
)

func toValue(keywordize bool, i interface{}) (vm.Value, error) {
	switch i := i.(type) {
	case string:
		return vm.String(i), nil
	case bool:
		return vm.Boolean(i), nil
	case float64:
		return vm.Int(i), nil
	case nil:
		return vm.NIL, nil
	case []interface{}:
		r := make([]vm.Value, len(i))
		for j := 0; j < len(i); j++ {
			v, e := toValue(keywordize, i[j])
			if e != nil {
				return vm.NIL, e
			}
			r[j] = v
		}
		return vm.NewArrayVector(r), nil
	case map[string]interface{}:
		newmap := make(vm.Map)
		for k, v := range i {
			ve, e := toValue(keywordize, v)
			if e != nil {
				return vm.NIL, e
			}
			if keywordize {
				newmap[vm.Keyword(k)] = ve
			} else {
				newmap[vm.String(k)] = ve
			}

		}
		return newmap, nil
	default:
		return vm.NIL, vm.NewExecutionError("invalid JSON value")
	}
}

func fromValue(v vm.Value) (interface{}, error) {
	switch v.Type() {
	case vm.StringType:
		return string(v.(vm.String)), nil
	case vm.IntType:
		return int(v.(vm.Int)), nil
	case vm.BooleanType:
		return bool(v.(vm.Boolean)), nil
	case vm.MapType:
		r := map[string]interface{}{}
		for k, ov := range v.(vm.Map) {
			vv, e := fromValue(ov)
			if e != nil {
				return vm.NIL, vm.NewExecutionError("invalid VM value")
			}
			nk := k.String()
			if k.Type() == vm.KeywordType {
				nk = nk[1:]
			}
			r[nk] = vv
		}
		return r, nil
	case vm.KeywordType:
		return string(v.(vm.Keyword)), nil
	case vm.NilType:
		return nil, nil
	default:
		s, ok := v.(vm.Seq)
		if !ok {
			return v.String(), nil
		}
		r := []interface{}{}
		for s != vm.EmptyList {
			uv, e := fromValue(s.First())
			if e != nil {
				return vm.NIL, e
			}
			r = append(r, uv)
			s = s.Next()
		}
		return r, nil
	}
}

func optionsHaveKeywordize(opts vm.Value) (bool, error) {
	o, ok := opts.(vm.Map)
	if !ok {
		return false, fmt.Errorf("read-json options are not Map")
	}
	return vm.IsTruthy(o.ValueAt(vm.Keyword("keywords?"))), nil
}

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

		var v interface{}
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

	// vars
	CurrentNS = ns.Def("*ns*", ns)

	ns.Def("read-json", readJson)
	ns.Def("write-json", writeJson)
	RegisterNS(ns)
}
