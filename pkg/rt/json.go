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

func toValue(keywordize bool, i interface{}) (vm.Value, error) {
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

func fromMapValue(v vm.Value) (interface{}, error) {
	r := map[string]interface{}{}
	if sq, ok := v.(vm.Sequable); ok {
		for s := sq.Seq(); s != nil && s != vm.EmptyList; s = s.Next() {
			entry := s.First()
			// Get key and value from the entry using Seq interface
			eSeq, ok := entry.(vm.Sequable)
			if !ok {
				return vm.NIL, vm.NewExecutionError("invalid map entry")
			}
			es := eSeq.Seq()
			k := es.First()
			ov := es.Next().First()
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
	}
	return r, nil
}

func fromSeqValue(s vm.Seq) (interface{}, error) {
	r := []interface{}{}
	for s != nil {
		uv, e := fromValue(s.First())
		if e != nil {
			return vm.NIL, e
		}
		r = append(r, uv)
		s = s.Next()
	}
	return r, nil
}

func fromValue(v vm.Value) (interface{}, error) {
	switch v.Type() {
	case vm.StringType:
		return string(v.(vm.String)), nil
	case vm.IntType:
		return int(v.(vm.Int)), nil
	case vm.FloatType:
		return float64(v.(vm.Float)), nil
	case vm.BooleanType:
		return bool(v.(vm.Boolean)), nil
	case vm.MapType, vm.PersistentMapType:
		return fromMapValue(v)
	case vm.KeywordType:
		kw := string(v.(vm.Keyword))
		return kw, nil
	case vm.NilType:
		return nil, nil
	case vm.ArrayVectorType, vm.PersistentVectorType:
		if sq, ok := v.(vm.Sequable); ok {
			return fromSeqValue(sq.Seq())
		}
		return v.String(), nil
	default:
		// Records and other map-like types
		if _, ok := v.(*vm.Record); ok {
			return fromMapValue(v)
		}
		s, ok := v.(vm.Seq)
		if !ok {
			if sq, ok := v.(vm.Sequable); ok {
				return fromSeqValue(sq.Seq())
			}
			return v.String(), nil
		}
		return fromSeqValue(s)
	}
}

func optionsHaveKeywordize(opts vm.Value) (bool, error) {
	o, ok := opts.(vm.Lookup)
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

	ns.Def("read-json", readJson)
	ns.Def("write-json", writeJson)
	RegisterNS(ns)
}
