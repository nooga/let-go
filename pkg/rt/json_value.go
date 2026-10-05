/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"github.com/nooga/let-go/pkg/vm"
)

// fromValue converts a let-go value to the closed set of Go types the JSON
// encoders accept (see internal/jsonenc). It lives outside json.go because
// js/emit needs it in lg_no_json builds too.
func fromMapValue(v vm.Value) (any, error) {
	r := map[string]any{}
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
			var nk string
			switch k := k.(type) {
			case vm.String:
				nk = string(k)
			case vm.Keyword:
				nk = string(k)
			default:
				nk = k.String()
			}
			r[nk] = vv
		}
	}
	return r, nil
}

func fromSeqValue(s vm.Seq) (any, error) {
	r := []any{}
	for s != nil && s != vm.EmptyList {
		uv, e := fromValue(s.First())
		if e != nil {
			return vm.NIL, e
		}
		r = append(r, uv)
		s = s.Next()
	}
	return r, nil
}

func fromValue(v vm.Value) (any, error) {
	switch v.Type() {
	case vm.StringType:
		return string(v.(vm.String)), nil
	case vm.IntType:
		return int(v.(vm.Int)), nil
	case vm.FloatType:
		// Float and Float32 both report FloatType; assert via Unbox so a Float32
		// (e.g. from `(float x)`) doesn't panic the float64 type assertion.
		return v.Unbox().(float64), nil
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
