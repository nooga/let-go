package ir_data_generated

import (
	rt "github.com/nooga/let-go/pkg/rt"
	vm "github.com/nooga/let-go/pkg/vm"
)

func type_of(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__7998_5 vm.Value
	var arg__7999_6 vm.Value
	var arg__8006_10 vm.Value
	var arg__8007_11 vm.Value
	var arg__8009_12 vm.Value
	var v13 vm.Value
	var callErr error
	_, _, _, _, _, _ = arg__7998_5, arg__7999_6, arg__8006_10, arg__8007_11, arg__8009_12, v13
	arg__7998_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__7999_6, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__7998_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__8006_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8007_11, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8006_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__8009_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__8007_11, arg0})
	if callErr != nil {
		return nil, callErr
	}
	v13, callErr = rt.InvokeValue(vm.Keyword("type"), []vm.Value{arg__8009_12})
	if callErr != nil {
		return nil, callErr
	}
	return v13, nil
}
func block_set_preds_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__8016_7 vm.Value
	var arg__8025_13 vm.Value
	var v14 vm.Value
	var callErr error
	_, _, _ = arg__8016_7, arg__8025_13, v14
	arg__8016_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("blocks"), arg1, vm.Keyword("preds")})
	if callErr != nil {
		return nil, callErr
	}
	arg__8025_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("blocks"), arg1, vm.Keyword("preds")})
	if callErr != nil {
		return nil, callErr
	}
	v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "assoc-in").Deref(), arg__8025_13, arg2})
	if callErr != nil {
		return nil, callErr
	}
	return v14, nil
}
func set_block_of_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__8033_7 vm.Value
	var arg__8042_13 vm.Value
	var v14 vm.Value
	var callErr error
	_, _, _ = arg__8033_7, arg__8042_13, v14
	arg__8033_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("insts"), arg1, vm.Keyword("block")})
	if callErr != nil {
		return nil, callErr
	}
	arg__8042_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("insts"), arg1, vm.Keyword("block")})
	if callErr != nil {
		return nil, callErr
	}
	v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "assoc-in").Deref(), arg__8042_13, arg2})
	if callErr != nil {
		return nil, callErr
	}
	return v14, nil
}
func block_set_params_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__8050_7 vm.Value
	var arg__8059_13 vm.Value
	var v14 vm.Value
	var callErr error
	_, _, _ = arg__8050_7, arg__8059_13, v14
	arg__8050_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("blocks"), arg1, vm.Keyword("params")})
	if callErr != nil {
		return nil, callErr
	}
	arg__8059_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("blocks"), arg1, vm.Keyword("params")})
	if callErr != nil {
		return nil, callErr
	}
	v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "assoc-in").Deref(), arg__8059_13, arg2})
	if callErr != nil {
		return nil, callErr
	}
	return v14, nil
}
func set_refs_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__8067_7 vm.Value
	var arg__8076_13 vm.Value
	var v14 vm.Value
	var v26 vm.Value
	var callErr error
	_, _, _, _ = arg__8067_7, arg__8076_13, v14, v26
	arg__8067_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("insts"), arg1, vm.Keyword("refs")})
	if callErr != nil {
		return nil, callErr
	}
	arg__8076_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("insts"), arg1, vm.Keyword("refs")})
	if callErr != nil {
		return nil, callErr
	}
	v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "assoc-in").Deref(), arg__8076_13, arg2})
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
	var arg__8096_5 vm.Value
	var arg__8097_6 vm.Value
	var arg__8104_10 vm.Value
	var arg__8105_11 vm.Value
	var arg__8107_12 vm.Value
	var v13 vm.Value
	var callErr error
	_, _, _, _, _, _ = arg__8096_5, arg__8097_6, arg__8104_10, arg__8105_11, arg__8107_12, v13
	arg__8096_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8097_6, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8096_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__8104_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8105_11, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8104_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__8107_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__8105_11, arg0})
	if callErr != nil {
		return nil, callErr
	}
	v13, callErr = rt.InvokeValue(vm.Keyword("refs"), []vm.Value{arg__8107_12})
	if callErr != nil {
		return nil, callErr
	}
	return v13, nil
}
func fn_uses_cache(arg0 vm.Value) (vm.Value, error) {
	var arg__8112_3 vm.Value
	var v4 vm.Value
	var callErr error
	_, _ = arg__8112_3, v4
	arg__8112_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v4, callErr = rt.InvokeValue(vm.Keyword("uses-cache"), []vm.Value{arg__8112_3})
	if callErr != nil {
		return nil, callErr
	}
	return v4, nil
}
func block_insts(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__8118_5 vm.Value
	var arg__8119_6 vm.Value
	var arg__8126_10 vm.Value
	var arg__8127_11 vm.Value
	var arg__8129_12 vm.Value
	var v13 vm.Value
	var callErr error
	_, _, _, _, _, _ = arg__8118_5, arg__8119_6, arg__8126_10, arg__8127_11, arg__8129_12, v13
	arg__8118_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8119_6, callErr = rt.InvokeValue(vm.Keyword("blocks"), []vm.Value{arg__8118_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__8126_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8127_11, callErr = rt.InvokeValue(vm.Keyword("blocks"), []vm.Value{arg__8126_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__8129_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__8127_11, arg0})
	if callErr != nil {
		return nil, callErr
	}
	v13, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8129_12})
	if callErr != nil {
		return nil, callErr
	}
	return v13, nil
}
func set_aux_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__8136_7 vm.Value
	var arg__8145_13 vm.Value
	var v14 vm.Value
	var callErr error
	_, _, _ = arg__8136_7, arg__8145_13, v14
	arg__8136_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("insts"), arg1, vm.Keyword("aux")})
	if callErr != nil {
		return nil, callErr
	}
	arg__8145_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("insts"), arg1, vm.Keyword("aux")})
	if callErr != nil {
		return nil, callErr
	}
	v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "assoc-in").Deref(), arg__8145_13, arg2})
	if callErr != nil {
		return nil, callErr
	}
	return v14, nil
}
func fn_consts(arg0 vm.Value) (vm.Value, error) {
	var arg__8151_3 vm.Value
	var v4 vm.Value
	var callErr error
	_, _ = arg__8151_3, v4
	arg__8151_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v4, callErr = rt.InvokeValue(vm.Keyword("consts"), []vm.Value{arg__8151_3})
	if callErr != nil {
		return nil, callErr
	}
	return v4, nil
}
func block_id(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__8157_5 vm.Value
	var arg__8158_6 vm.Value
	var arg__8165_10 vm.Value
	var arg__8166_11 vm.Value
	var arg__8168_12 vm.Value
	var v13 vm.Value
	var callErr error
	_, _, _, _, _, _ = arg__8157_5, arg__8158_6, arg__8165_10, arg__8166_11, arg__8168_12, v13
	arg__8157_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8158_6, callErr = rt.InvokeValue(vm.Keyword("blocks"), []vm.Value{arg__8157_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__8165_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8166_11, callErr = rt.InvokeValue(vm.Keyword("blocks"), []vm.Value{arg__8165_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__8168_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__8166_11, arg0})
	if callErr != nil {
		return nil, callErr
	}
	v13, callErr = rt.InvokeValue(vm.Keyword("id"), []vm.Value{arg__8168_12})
	if callErr != nil {
		return nil, callErr
	}
	return v13, nil
}
func op(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__8174_5 vm.Value
	var arg__8175_6 vm.Value
	var arg__8182_10 vm.Value
	var arg__8183_11 vm.Value
	var arg__8185_12 vm.Value
	var v13 vm.Value
	var callErr error
	_, _, _, _, _, _ = arg__8174_5, arg__8175_6, arg__8182_10, arg__8183_11, arg__8185_12, v13
	arg__8174_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8175_6, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8174_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__8182_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8183_11, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8182_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__8185_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__8183_11, arg0})
	if callErr != nil {
		return nil, callErr
	}
	v13, callErr = rt.InvokeValue(vm.Keyword("op"), []vm.Value{arg__8185_12})
	if callErr != nil {
		return nil, callErr
	}
	return v13, nil
}
func fn_name(arg0 vm.Value) (vm.Value, error) {
	var arg__8190_3 vm.Value
	var v4 vm.Value
	var callErr error
	_, _ = arg__8190_3, v4
	arg__8190_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v4, callErr = rt.InvokeValue(vm.Keyword("name"), []vm.Value{arg__8190_3})
	if callErr != nil {
		return nil, callErr
	}
	return v4, nil
}
func set_type_of_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__8197_7 vm.Value
	var arg__8206_13 vm.Value
	var v14 vm.Value
	var callErr error
	_, _, _ = arg__8197_7, arg__8206_13, v14
	arg__8197_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("insts"), arg1, vm.Keyword("type")})
	if callErr != nil {
		return nil, callErr
	}
	arg__8206_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("insts"), arg1, vm.Keyword("type")})
	if callErr != nil {
		return nil, callErr
	}
	v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "assoc-in").Deref(), arg__8206_13, arg2})
	if callErr != nil {
		return nil, callErr
	}
	return v14, nil
}
func fn_entry(arg0 vm.Value) (vm.Value, error) {
	var arg__8212_3 vm.Value
	var v4 vm.Value
	var callErr error
	_, _ = arg__8212_3, v4
	arg__8212_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v4, callErr = rt.InvokeValue(vm.Keyword("entry"), []vm.Value{arg__8212_3})
	if callErr != nil {
		return nil, callErr
	}
	return v4, nil
}
func fn_arity(arg0 vm.Value) (vm.Value, error) {
	var arg__8217_3 vm.Value
	var v4 vm.Value
	var callErr error
	_, _ = arg__8217_3, v4
	arg__8217_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v4, callErr = rt.InvokeValue(vm.Keyword("arity"), []vm.Value{arg__8217_3})
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
	var arg__8227_5 vm.Value
	var arg__8228_6 vm.Value
	var arg__8235_10 vm.Value
	var arg__8236_11 vm.Value
	var arg__8238_12 vm.Value
	var v13 vm.Value
	var callErr error
	_, _, _, _, _, _ = arg__8227_5, arg__8228_6, arg__8235_10, arg__8236_11, arg__8238_12, v13
	arg__8227_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8228_6, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8227_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__8235_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8236_11, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8235_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__8238_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__8236_11, arg0})
	if callErr != nil {
		return nil, callErr
	}
	v13, callErr = rt.InvokeValue(vm.Keyword("block"), []vm.Value{arg__8238_12})
	if callErr != nil {
		return nil, callErr
	}
	return v13, nil
}
func fn_insts(arg0 vm.Value) (vm.Value, error) {
	var arg__8243_3 vm.Value
	var v4 vm.Value
	var callErr error
	_, _ = arg__8243_3, v4
	arg__8243_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v4, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8243_3})
	if callErr != nil {
		return nil, callErr
	}
	return v4, nil
}
func set_source_infos_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__8250_7 vm.Value
	var arg__8259_13 vm.Value
	var v14 vm.Value
	var callErr error
	_, _, _ = arg__8250_7, arg__8259_13, v14
	arg__8250_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("insts"), arg1, vm.Keyword("source-infos")})
	if callErr != nil {
		return nil, callErr
	}
	arg__8259_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("insts"), arg1, vm.Keyword("source-infos")})
	if callErr != nil {
		return nil, callErr
	}
	v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "assoc-in").Deref(), arg__8259_13, arg2})
	if callErr != nil {
		return nil, callErr
	}
	return v14, nil
}
func set_op_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__8267_7 vm.Value
	var arg__8276_13 vm.Value
	var v14 vm.Value
	var v26 vm.Value
	var callErr error
	_, _, _, _ = arg__8267_7, arg__8276_13, v14, v26
	arg__8267_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("insts"), arg1, vm.Keyword("op")})
	if callErr != nil {
		return nil, callErr
	}
	arg__8276_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("insts"), arg1, vm.Keyword("op")})
	if callErr != nil {
		return nil, callErr
	}
	v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "assoc-in").Deref(), arg__8276_13, arg2})
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
	var arg__8296_5 vm.Value
	var arg__8297_6 vm.Value
	var arg__8304_10 vm.Value
	var arg__8305_11 vm.Value
	var arg__8307_12 vm.Value
	var v13 vm.Value
	var callErr error
	_, _, _, _, _, _ = arg__8296_5, arg__8297_6, arg__8304_10, arg__8305_11, arg__8307_12, v13
	arg__8296_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8297_6, callErr = rt.InvokeValue(vm.Keyword("blocks"), []vm.Value{arg__8296_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__8304_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8305_11, callErr = rt.InvokeValue(vm.Keyword("blocks"), []vm.Value{arg__8304_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__8307_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__8305_11, arg0})
	if callErr != nil {
		return nil, callErr
	}
	v13, callErr = rt.InvokeValue(vm.Keyword("preds"), []vm.Value{arg__8307_12})
	if callErr != nil {
		return nil, callErr
	}
	return v13, nil
}
func source_infos(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__8313_5 vm.Value
	var arg__8314_6 vm.Value
	var arg__8321_10 vm.Value
	var arg__8322_11 vm.Value
	var arg__8324_12 vm.Value
	var v13 vm.Value
	var callErr error
	_, _, _, _, _, _ = arg__8313_5, arg__8314_6, arg__8321_10, arg__8322_11, arg__8324_12, v13
	arg__8313_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8314_6, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8313_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__8321_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8322_11, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8321_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__8324_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__8322_11, arg0})
	if callErr != nil {
		return nil, callErr
	}
	v13, callErr = rt.InvokeValue(vm.Keyword("source-infos"), []vm.Value{arg__8324_12})
	if callErr != nil {
		return nil, callErr
	}
	return v13, nil
}
func aux(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__8330_5 vm.Value
	var arg__8331_6 vm.Value
	var arg__8338_10 vm.Value
	var arg__8339_11 vm.Value
	var arg__8341_12 vm.Value
	var v13 vm.Value
	var callErr error
	_, _, _, _, _, _ = arg__8330_5, arg__8331_6, arg__8338_10, arg__8339_11, arg__8341_12, v13
	arg__8330_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8331_6, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8330_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__8338_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8339_11, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8338_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__8341_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__8339_11, arg0})
	if callErr != nil {
		return nil, callErr
	}
	v13, callErr = rt.InvokeValue(vm.Keyword("aux"), []vm.Value{arg__8341_12})
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
	var arg__8350_7 vm.Value
	var arg__8359_13 vm.Value
	var v14 vm.Value
	var callErr error
	_, _, _ = arg__8350_7, arg__8359_13, v14
	arg__8350_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("blocks"), arg1, vm.Keyword("id")})
	if callErr != nil {
		return nil, callErr
	}
	arg__8359_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("blocks"), arg1, vm.Keyword("id")})
	if callErr != nil {
		return nil, callErr
	}
	v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "assoc-in").Deref(), arg__8359_13, arg2})
	if callErr != nil {
		return nil, callErr
	}
	return v14, nil
}
func block_params(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__8366_5 vm.Value
	var arg__8367_6 vm.Value
	var arg__8374_10 vm.Value
	var arg__8375_11 vm.Value
	var arg__8377_12 vm.Value
	var v13 vm.Value
	var callErr error
	_, _, _, _, _, _ = arg__8366_5, arg__8367_6, arg__8374_10, arg__8375_11, arg__8377_12, v13
	arg__8366_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8367_6, callErr = rt.InvokeValue(vm.Keyword("blocks"), []vm.Value{arg__8366_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__8374_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8375_11, callErr = rt.InvokeValue(vm.Keyword("blocks"), []vm.Value{arg__8374_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__8377_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__8375_11, arg0})
	if callErr != nil {
		return nil, callErr
	}
	v13, callErr = rt.InvokeValue(vm.Keyword("params"), []vm.Value{arg__8377_12})
	if callErr != nil {
		return nil, callErr
	}
	return v13, nil
}
func block_set_term_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__8384_7 vm.Value
	var arg__8393_13 vm.Value
	var v14 vm.Value
	var callErr error
	_, _, _ = arg__8384_7, arg__8393_13, v14
	arg__8384_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("blocks"), arg1, vm.Keyword("term")})
	if callErr != nil {
		return nil, callErr
	}
	arg__8393_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("blocks"), arg1, vm.Keyword("term")})
	if callErr != nil {
		return nil, callErr
	}
	v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "assoc-in").Deref(), arg__8393_13, arg2})
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
	var arg__8401_3 vm.Value
	var v4 vm.Value
	var callErr error
	_, _ = arg__8401_3, v4
	arg__8401_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v4, callErr = rt.InvokeValue(vm.Keyword("variadic?"), []vm.Value{arg__8401_3})
	if callErr != nil {
		return nil, callErr
	}
	return v4, nil
}
func block_set_insts_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__8408_7 vm.Value
	var arg__8417_13 vm.Value
	var v14 vm.Value
	var callErr error
	_, _, _ = arg__8408_7, arg__8417_13, v14
	arg__8408_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("blocks"), arg1, vm.Keyword("insts")})
	if callErr != nil {
		return nil, callErr
	}
	arg__8417_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("blocks"), arg1, vm.Keyword("insts")})
	if callErr != nil {
		return nil, callErr
	}
	v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "assoc-in").Deref(), arg__8417_13, arg2})
	if callErr != nil {
		return nil, callErr
	}
	return v14, nil
}
func fn_uses_dirty_QMARK_(arg0 vm.Value) (vm.Value, error) {
	var arg__8423_3 vm.Value
	var v4 vm.Value
	var callErr error
	_, _ = arg__8423_3, v4
	arg__8423_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v4, callErr = rt.InvokeValue(vm.Keyword("uses-dirty?"), []vm.Value{arg__8423_3})
	if callErr != nil {
		return nil, callErr
	}
	return v4, nil
}
func fn_blocks(arg0 vm.Value) (vm.Value, error) {
	var arg__8428_3 vm.Value
	var v4 vm.Value
	var callErr error
	_, _ = arg__8428_3, v4
	arg__8428_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v4, callErr = rt.InvokeValue(vm.Keyword("blocks"), []vm.Value{arg__8428_3})
	if callErr != nil {
		return nil, callErr
	}
	return v4, nil
}
func block_term(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__8434_5 vm.Value
	var arg__8435_6 vm.Value
	var arg__8442_10 vm.Value
	var arg__8443_11 vm.Value
	var arg__8445_12 vm.Value
	var v13 vm.Value
	var callErr error
	_, _, _, _, _, _ = arg__8434_5, arg__8435_6, arg__8442_10, arg__8443_11, arg__8445_12, v13
	arg__8434_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8435_6, callErr = rt.InvokeValue(vm.Keyword("blocks"), []vm.Value{arg__8434_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__8442_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8443_11, callErr = rt.InvokeValue(vm.Keyword("blocks"), []vm.Value{arg__8442_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__8445_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__8443_11, arg0})
	if callErr != nil {
		return nil, callErr
	}
	v13, callErr = rt.InvokeValue(vm.Keyword("term"), []vm.Value{arg__8445_12})
	if callErr != nil {
		return nil, callErr
	}
	return v13, nil
}
