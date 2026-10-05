//go:build !lg_no_json

/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/nooga/let-go/pkg/vm"
)

// decodeJSON parses one JSON value with UseNumber, so numbers reach toValue as
// their literal text. json.Unmarshal into any yields float64 for every number,
// which has already rounded integers above 2^53 before anything can convert
// them. Trailing non-whitespace is rejected, as json.Unmarshal does.
func decodeJSON(s string) (any, error) {
	dec := json.NewDecoder(strings.NewReader(s))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("unexpected end of JSON input")
		}
		return nil, err
	}
	if rest := strings.TrimLeft(s[dec.InputOffset():], " \t\r\n"); rest != "" {
		return nil, fmt.Errorf("invalid character %q after top-level value", rest[0])
	}
	return v, nil
}

// numberToValue keeps integer literals exact: int64 when they fit, BigInt past
// that, matching the reader. Decimal and exponent forms go through float64,
// and a whole float still becomes an Int when it is in int64 range.
func numberToValue(n json.Number) (vm.Value, error) {
	s := string(n)
	if !strings.ContainsAny(s, ".eE") {
		if i, err := n.Int64(); err == nil {
			return vm.Int(i), nil
		}
		if b, ok := vm.NewBigIntFromString(s); ok {
			return b, nil
		}
	}
	f, err := n.Float64()
	if err != nil {
		return vm.NIL, fmt.Errorf("json: number %s out of range", s)
	}
	return floatToValue(f), nil
}

// floatToValue converts a whole float to Int only inside int64 range:
// float64(int64(f)) is platform-defined past it (arm64 saturates, amd64 wraps),
// so the unguarded round-trip check misreads 2^63 as an Int on arm64.
func floatToValue(f float64) vm.Value {
	if f == math.Trunc(f) && f >= -(1<<63) && f < 1<<63 {
		return vm.Int(int64(f))
	}
	return vm.Float(f)
}

func toValue(keywordize bool, i any) (vm.Value, error) {
	switch i := i.(type) {
	case string:
		return vm.String(i), nil
	case bool:
		return vm.Boolean(i), nil
	case json.Number:
		return numberToValue(i)
	case float64:
		return floatToValue(i), nil
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

		v, err := decodeJSON(string(s))
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
