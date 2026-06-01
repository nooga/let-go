package ir_data_generated

import (
	rt "github.com/nooga/let-go/pkg/rt"
	vm "github.com/nooga/let-go/pkg/vm"
)

func type_of(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__7999_5 vm.Value
	var arg__8000_6 vm.Value
	var arg__8007_10 vm.Value
	var arg__8008_11 vm.Value
	var arg__8010_12 vm.Value
	var v13 vm.Value
	var callErr error
	_, _, _, _, _, _ = arg__7999_5, arg__8000_6, arg__8007_10, arg__8008_11, arg__8010_12, v13
	arg__7999_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8000_6, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__7999_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__8007_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8008_11, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8007_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__8010_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__8008_11, arg0})
	if callErr != nil {
		return nil, callErr
	}
	v13, callErr = rt.InvokeValue(vm.Keyword("type"), []vm.Value{arg__8010_12})
	if callErr != nil {
		return nil, callErr
	}
	return v13, nil
}
func block_set_preds_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__8017_7 vm.Value
	var arg__8026_13 vm.Value
	var v14 vm.Value
	var callErr error
	_, _, _ = arg__8017_7, arg__8026_13, v14
	arg__8017_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("blocks"), arg1, vm.Keyword("preds")})
	if callErr != nil {
		return nil, callErr
	}
	arg__8026_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("blocks"), arg1, vm.Keyword("preds")})
	if callErr != nil {
		return nil, callErr
	}
	v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "assoc-in").Deref(), arg__8026_13, arg2})
	if callErr != nil {
		return nil, callErr
	}
	return v14, nil
}
func set_block_of_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__8034_7 vm.Value
	var arg__8043_13 vm.Value
	var v14 vm.Value
	var callErr error
	_, _, _ = arg__8034_7, arg__8043_13, v14
	arg__8034_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("insts"), arg1, vm.Keyword("block")})
	if callErr != nil {
		return nil, callErr
	}
	arg__8043_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("insts"), arg1, vm.Keyword("block")})
	if callErr != nil {
		return nil, callErr
	}
	v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "assoc-in").Deref(), arg__8043_13, arg2})
	if callErr != nil {
		return nil, callErr
	}
	return v14, nil
}
func block_set_params_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__8051_7 vm.Value
	var arg__8060_13 vm.Value
	var v14 vm.Value
	var callErr error
	_, _, _ = arg__8051_7, arg__8060_13, v14
	arg__8051_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("blocks"), arg1, vm.Keyword("params")})
	if callErr != nil {
		return nil, callErr
	}
	arg__8060_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("blocks"), arg1, vm.Keyword("params")})
	if callErr != nil {
		return nil, callErr
	}
	v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "assoc-in").Deref(), arg__8060_13, arg2})
	if callErr != nil {
		return nil, callErr
	}
	return v14, nil
}
func set_refs_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__8068_7 vm.Value
	var arg__8077_13 vm.Value
	var v14 vm.Value
	var v26 vm.Value
	var callErr error
	_, _, _, _ = arg__8068_7, arg__8077_13, v14, v26
	arg__8068_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("insts"), arg1, vm.Keyword("refs")})
	if callErr != nil {
		return nil, callErr
	}
	arg__8077_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("insts"), arg1, vm.Keyword("refs")})
	if callErr != nil {
		return nil, callErr
	}
	v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "assoc-in").Deref(), arg__8077_13, arg2})
	if callErr != nil {
		return nil, callErr
	}
	v26, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "assoc").Deref(), vm.Keyword("uses-dirty?"), vm.Boolean(true), vm.Keyword("uses-cache"), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	return v26, nil
}
func refs(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__8097_5 vm.Value
	var arg__8098_6 vm.Value
	var arg__8105_10 vm.Value
	var arg__8106_11 vm.Value
	var arg__8108_12 vm.Value
	var v13 vm.Value
	var callErr error
	_, _, _, _, _, _ = arg__8097_5, arg__8098_6, arg__8105_10, arg__8106_11, arg__8108_12, v13
	arg__8097_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8098_6, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8097_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__8105_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8106_11, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8105_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__8108_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__8106_11, arg0})
	if callErr != nil {
		return nil, callErr
	}
	v13, callErr = rt.InvokeValue(vm.Keyword("refs"), []vm.Value{arg__8108_12})
	if callErr != nil {
		return nil, callErr
	}
	return v13, nil
}
func fn_uses_cache(arg0 vm.Value) (vm.Value, error) {
	var arg__8113_3 vm.Value
	var v4 vm.Value
	var callErr error
	_, _ = arg__8113_3, v4
	arg__8113_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v4, callErr = rt.InvokeValue(vm.Keyword("uses-cache"), []vm.Value{arg__8113_3})
	if callErr != nil {
		return nil, callErr
	}
	return v4, nil
}
func block_insts(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__8119_5 vm.Value
	var arg__8120_6 vm.Value
	var arg__8127_10 vm.Value
	var arg__8128_11 vm.Value
	var arg__8130_12 vm.Value
	var v13 vm.Value
	var callErr error
	_, _, _, _, _, _ = arg__8119_5, arg__8120_6, arg__8127_10, arg__8128_11, arg__8130_12, v13
	arg__8119_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8120_6, callErr = rt.InvokeValue(vm.Keyword("blocks"), []vm.Value{arg__8119_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__8127_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8128_11, callErr = rt.InvokeValue(vm.Keyword("blocks"), []vm.Value{arg__8127_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__8130_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__8128_11, arg0})
	if callErr != nil {
		return nil, callErr
	}
	v13, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8130_12})
	if callErr != nil {
		return nil, callErr
	}
	return v13, nil
}
func set_aux_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__8137_7 vm.Value
	var arg__8146_13 vm.Value
	var v14 vm.Value
	var callErr error
	_, _, _ = arg__8137_7, arg__8146_13, v14
	arg__8137_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("insts"), arg1, vm.Keyword("aux")})
	if callErr != nil {
		return nil, callErr
	}
	arg__8146_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("insts"), arg1, vm.Keyword("aux")})
	if callErr != nil {
		return nil, callErr
	}
	v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "assoc-in").Deref(), arg__8146_13, arg2})
	if callErr != nil {
		return nil, callErr
	}
	return v14, nil
}
func fn_consts(arg0 vm.Value) (vm.Value, error) {
	var arg__8152_3 vm.Value
	var v4 vm.Value
	var callErr error
	_, _ = arg__8152_3, v4
	arg__8152_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v4, callErr = rt.InvokeValue(vm.Keyword("consts"), []vm.Value{arg__8152_3})
	if callErr != nil {
		return nil, callErr
	}
	return v4, nil
}
func block_id(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__8158_5 vm.Value
	var arg__8159_6 vm.Value
	var arg__8166_10 vm.Value
	var arg__8167_11 vm.Value
	var arg__8169_12 vm.Value
	var v13 vm.Value
	var callErr error
	_, _, _, _, _, _ = arg__8158_5, arg__8159_6, arg__8166_10, arg__8167_11, arg__8169_12, v13
	arg__8158_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8159_6, callErr = rt.InvokeValue(vm.Keyword("blocks"), []vm.Value{arg__8158_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__8166_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8167_11, callErr = rt.InvokeValue(vm.Keyword("blocks"), []vm.Value{arg__8166_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__8169_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__8167_11, arg0})
	if callErr != nil {
		return nil, callErr
	}
	v13, callErr = rt.InvokeValue(vm.Keyword("id"), []vm.Value{arg__8169_12})
	if callErr != nil {
		return nil, callErr
	}
	return v13, nil
}
func op(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__8175_5 vm.Value
	var arg__8176_6 vm.Value
	var arg__8183_10 vm.Value
	var arg__8184_11 vm.Value
	var arg__8186_12 vm.Value
	var v13 vm.Value
	var callErr error
	_, _, _, _, _, _ = arg__8175_5, arg__8176_6, arg__8183_10, arg__8184_11, arg__8186_12, v13
	arg__8175_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8176_6, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8175_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__8183_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8184_11, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8183_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__8186_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__8184_11, arg0})
	if callErr != nil {
		return nil, callErr
	}
	v13, callErr = rt.InvokeValue(vm.Keyword("op"), []vm.Value{arg__8186_12})
	if callErr != nil {
		return nil, callErr
	}
	return v13, nil
}
func fn_name(arg0 vm.Value) (vm.Value, error) {
	var arg__8191_3 vm.Value
	var v4 vm.Value
	var callErr error
	_, _ = arg__8191_3, v4
	arg__8191_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v4, callErr = rt.InvokeValue(vm.Keyword("name"), []vm.Value{arg__8191_3})
	if callErr != nil {
		return nil, callErr
	}
	return v4, nil
}
func set_type_of_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__8198_7 vm.Value
	var arg__8207_13 vm.Value
	var v14 vm.Value
	var callErr error
	_, _, _ = arg__8198_7, arg__8207_13, v14
	arg__8198_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("insts"), arg1, vm.Keyword("type")})
	if callErr != nil {
		return nil, callErr
	}
	arg__8207_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("insts"), arg1, vm.Keyword("type")})
	if callErr != nil {
		return nil, callErr
	}
	v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "assoc-in").Deref(), arg__8207_13, arg2})
	if callErr != nil {
		return nil, callErr
	}
	return v14, nil
}
func fn_entry(arg0 vm.Value) (vm.Value, error) {
	var arg__8213_3 vm.Value
	var v4 vm.Value
	var callErr error
	_, _ = arg__8213_3, v4
	arg__8213_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v4, callErr = rt.InvokeValue(vm.Keyword("entry"), []vm.Value{arg__8213_3})
	if callErr != nil {
		return nil, callErr
	}
	return v4, nil
}
func fn_arity(arg0 vm.Value) (vm.Value, error) {
	var arg__8218_3 vm.Value
	var v4 vm.Value
	var callErr error
	_, _ = arg__8218_3, v4
	arg__8218_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v4, callErr = rt.InvokeValue(vm.Keyword("arity"), []vm.Value{arg__8218_3})
	if callErr != nil {
		return nil, callErr
	}
	return v4, nil
}
func branch_target_target(arg0 vm.Value) (vm.Value, error) {
	var v2 vm.Value
	var callErr error
	_ = v2
	v2, callErr = rt.InvokeValue(vm.Keyword("target"), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	return v2, nil
}
func cond_target_false(arg0 vm.Value) (vm.Value, error) {
	var v2 vm.Value
	var callErr error
	_ = v2
	v2, callErr = rt.InvokeValue(vm.Keyword("false-target"), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	return v2, nil
}
func block_of(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__8228_5 vm.Value
	var arg__8229_6 vm.Value
	var arg__8236_10 vm.Value
	var arg__8237_11 vm.Value
	var arg__8239_12 vm.Value
	var v13 vm.Value
	var callErr error
	_, _, _, _, _, _ = arg__8228_5, arg__8229_6, arg__8236_10, arg__8237_11, arg__8239_12, v13
	arg__8228_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8229_6, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8228_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__8236_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8237_11, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8236_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__8239_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__8237_11, arg0})
	if callErr != nil {
		return nil, callErr
	}
	v13, callErr = rt.InvokeValue(vm.Keyword("block"), []vm.Value{arg__8239_12})
	if callErr != nil {
		return nil, callErr
	}
	return v13, nil
}
func fn_insts(arg0 vm.Value) (vm.Value, error) {
	var arg__8244_3 vm.Value
	var v4 vm.Value
	var callErr error
	_, _ = arg__8244_3, v4
	arg__8244_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v4, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8244_3})
	if callErr != nil {
		return nil, callErr
	}
	return v4, nil
}
func set_source_infos_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__8251_7 vm.Value
	var arg__8260_13 vm.Value
	var v14 vm.Value
	var callErr error
	_, _, _ = arg__8251_7, arg__8260_13, v14
	arg__8251_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("insts"), arg1, vm.Keyword("source-infos")})
	if callErr != nil {
		return nil, callErr
	}
	arg__8260_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("insts"), arg1, vm.Keyword("source-infos")})
	if callErr != nil {
		return nil, callErr
	}
	v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "assoc-in").Deref(), arg__8260_13, arg2})
	if callErr != nil {
		return nil, callErr
	}
	return v14, nil
}
func set_op_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__8268_7 vm.Value
	var arg__8277_13 vm.Value
	var v14 vm.Value
	var v26 vm.Value
	var callErr error
	_, _, _, _ = arg__8268_7, arg__8277_13, v14, v26
	arg__8268_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("insts"), arg1, vm.Keyword("op")})
	if callErr != nil {
		return nil, callErr
	}
	arg__8277_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("insts"), arg1, vm.Keyword("op")})
	if callErr != nil {
		return nil, callErr
	}
	v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "assoc-in").Deref(), arg__8277_13, arg2})
	if callErr != nil {
		return nil, callErr
	}
	v26, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "assoc").Deref(), vm.Keyword("uses-dirty?"), vm.Boolean(true), vm.Keyword("uses-cache"), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	return v26, nil
}
func block_preds(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__8297_5 vm.Value
	var arg__8298_6 vm.Value
	var arg__8305_10 vm.Value
	var arg__8306_11 vm.Value
	var arg__8308_12 vm.Value
	var v13 vm.Value
	var callErr error
	_, _, _, _, _, _ = arg__8297_5, arg__8298_6, arg__8305_10, arg__8306_11, arg__8308_12, v13
	arg__8297_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8298_6, callErr = rt.InvokeValue(vm.Keyword("blocks"), []vm.Value{arg__8297_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__8305_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8306_11, callErr = rt.InvokeValue(vm.Keyword("blocks"), []vm.Value{arg__8305_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__8308_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__8306_11, arg0})
	if callErr != nil {
		return nil, callErr
	}
	v13, callErr = rt.InvokeValue(vm.Keyword("preds"), []vm.Value{arg__8308_12})
	if callErr != nil {
		return nil, callErr
	}
	return v13, nil
}
func source_infos(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__8314_5 vm.Value
	var arg__8315_6 vm.Value
	var arg__8322_10 vm.Value
	var arg__8323_11 vm.Value
	var arg__8325_12 vm.Value
	var v13 vm.Value
	var callErr error
	_, _, _, _, _, _ = arg__8314_5, arg__8315_6, arg__8322_10, arg__8323_11, arg__8325_12, v13
	arg__8314_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8315_6, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8314_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__8322_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8323_11, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8322_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__8325_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__8323_11, arg0})
	if callErr != nil {
		return nil, callErr
	}
	v13, callErr = rt.InvokeValue(vm.Keyword("source-infos"), []vm.Value{arg__8325_12})
	if callErr != nil {
		return nil, callErr
	}
	return v13, nil
}
func aux(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__8331_5 vm.Value
	var arg__8332_6 vm.Value
	var arg__8339_10 vm.Value
	var arg__8340_11 vm.Value
	var arg__8342_12 vm.Value
	var v13 vm.Value
	var callErr error
	_, _, _, _, _, _ = arg__8331_5, arg__8332_6, arg__8339_10, arg__8340_11, arg__8342_12, v13
	arg__8331_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8332_6, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8331_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__8339_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8340_11, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8339_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__8342_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__8340_11, arg0})
	if callErr != nil {
		return nil, callErr
	}
	v13, callErr = rt.InvokeValue(vm.Keyword("aux"), []vm.Value{arg__8342_12})
	if callErr != nil {
		return nil, callErr
	}
	return v13, nil
}
func branch_target_args(arg0 vm.Value) (vm.Value, error) {
	var v2 vm.Value
	var callErr error
	_ = v2
	v2, callErr = rt.InvokeValue(vm.Keyword("args"), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	return v2, nil
}
func block_set_id_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__8351_7 vm.Value
	var arg__8360_13 vm.Value
	var v14 vm.Value
	var callErr error
	_, _, _ = arg__8351_7, arg__8360_13, v14
	arg__8351_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("blocks"), arg1, vm.Keyword("id")})
	if callErr != nil {
		return nil, callErr
	}
	arg__8360_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("blocks"), arg1, vm.Keyword("id")})
	if callErr != nil {
		return nil, callErr
	}
	v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "assoc-in").Deref(), arg__8360_13, arg2})
	if callErr != nil {
		return nil, callErr
	}
	return v14, nil
}
func block_params(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__8367_5 vm.Value
	var arg__8368_6 vm.Value
	var arg__8375_10 vm.Value
	var arg__8376_11 vm.Value
	var arg__8378_12 vm.Value
	var v13 vm.Value
	var callErr error
	_, _, _, _, _, _ = arg__8367_5, arg__8368_6, arg__8375_10, arg__8376_11, arg__8378_12, v13
	arg__8367_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8368_6, callErr = rt.InvokeValue(vm.Keyword("blocks"), []vm.Value{arg__8367_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__8375_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8376_11, callErr = rt.InvokeValue(vm.Keyword("blocks"), []vm.Value{arg__8375_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__8378_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__8376_11, arg0})
	if callErr != nil {
		return nil, callErr
	}
	v13, callErr = rt.InvokeValue(vm.Keyword("params"), []vm.Value{arg__8378_12})
	if callErr != nil {
		return nil, callErr
	}
	return v13, nil
}
func block_set_term_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__8385_7 vm.Value
	var arg__8394_13 vm.Value
	var v14 vm.Value
	var callErr error
	_, _, _ = arg__8385_7, arg__8394_13, v14
	arg__8385_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("blocks"), arg1, vm.Keyword("term")})
	if callErr != nil {
		return nil, callErr
	}
	arg__8394_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("blocks"), arg1, vm.Keyword("term")})
	if callErr != nil {
		return nil, callErr
	}
	v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "assoc-in").Deref(), arg__8394_13, arg2})
	if callErr != nil {
		return nil, callErr
	}
	return v14, nil
}
func cond_target_true(arg0 vm.Value) (vm.Value, error) {
	var v2 vm.Value
	var callErr error
	_ = v2
	v2, callErr = rt.InvokeValue(vm.Keyword("true-target"), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	return v2, nil
}
func fn_variadic_QMARK_(arg0 vm.Value) (vm.Value, error) {
	var arg__8402_3 vm.Value
	var v4 vm.Value
	var callErr error
	_, _ = arg__8402_3, v4
	arg__8402_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v4, callErr = rt.InvokeValue(vm.Keyword("variadic?"), []vm.Value{arg__8402_3})
	if callErr != nil {
		return nil, callErr
	}
	return v4, nil
}
func block_set_insts_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__8409_7 vm.Value
	var arg__8418_13 vm.Value
	var v14 vm.Value
	var callErr error
	_, _, _ = arg__8409_7, arg__8418_13, v14
	arg__8409_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("blocks"), arg1, vm.Keyword("insts")})
	if callErr != nil {
		return nil, callErr
	}
	arg__8418_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("blocks"), arg1, vm.Keyword("insts")})
	if callErr != nil {
		return nil, callErr
	}
	v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "assoc-in").Deref(), arg__8418_13, arg2})
	if callErr != nil {
		return nil, callErr
	}
	return v14, nil
}
func fn_uses_dirty_QMARK_(arg0 vm.Value) (vm.Value, error) {
	var arg__8424_3 vm.Value
	var v4 vm.Value
	var callErr error
	_, _ = arg__8424_3, v4
	arg__8424_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v4, callErr = rt.InvokeValue(vm.Keyword("uses-dirty?"), []vm.Value{arg__8424_3})
	if callErr != nil {
		return nil, callErr
	}
	return v4, nil
}
func fn_blocks(arg0 vm.Value) (vm.Value, error) {
	var arg__8429_3 vm.Value
	var v4 vm.Value
	var callErr error
	_, _ = arg__8429_3, v4
	arg__8429_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v4, callErr = rt.InvokeValue(vm.Keyword("blocks"), []vm.Value{arg__8429_3})
	if callErr != nil {
		return nil, callErr
	}
	return v4, nil
}
func block_term(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__8435_5 vm.Value
	var arg__8436_6 vm.Value
	var arg__8443_10 vm.Value
	var arg__8444_11 vm.Value
	var arg__8446_12 vm.Value
	var v13 vm.Value
	var callErr error
	_, _, _, _, _, _ = arg__8435_5, arg__8436_6, arg__8443_10, arg__8444_11, arg__8446_12, v13
	arg__8435_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8436_6, callErr = rt.InvokeValue(vm.Keyword("blocks"), []vm.Value{arg__8435_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__8443_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8444_11, callErr = rt.InvokeValue(vm.Keyword("blocks"), []vm.Value{arg__8443_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__8446_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__8444_11, arg0})
	if callErr != nil {
		return nil, callErr
	}
	v13, callErr = rt.InvokeValue(vm.Keyword("term"), []vm.Value{arg__8446_12})
	if callErr != nil {
		return nil, callErr
	}
	return v13, nil
}
