package ir_data_generated

import (
	rt "github.com/nooga/let-go/pkg/rt"
	vm "github.com/nooga/let-go/pkg/vm"
)

func type_of(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__8003_5 vm.Value
	var arg__8004_6 vm.Value
	var arg__8011_10 vm.Value
	var arg__8012_11 vm.Value
	var arg__8014_12 vm.Value
	var v13 vm.Value
	var callErr error
	_, _, _, _, _, _ = arg__8003_5, arg__8004_6, arg__8011_10, arg__8012_11, arg__8014_12, v13
	arg__8003_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8004_6, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8003_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__8011_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8012_11, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8011_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__8014_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__8012_11, arg0})
	if callErr != nil {
		return nil, callErr
	}
	v13, callErr = rt.InvokeValue(vm.Keyword("type"), []vm.Value{arg__8014_12})
	if callErr != nil {
		return nil, callErr
	}
	return v13, nil
}
func block_set_preds_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__8021_7 vm.Value
	var arg__8030_13 vm.Value
	var v14 vm.Value
	var callErr error
	_, _, _ = arg__8021_7, arg__8030_13, v14
	arg__8021_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("blocks"), arg1, vm.Keyword("preds")})
	if callErr != nil {
		return nil, callErr
	}
	arg__8030_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("blocks"), arg1, vm.Keyword("preds")})
	if callErr != nil {
		return nil, callErr
	}
	v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "assoc-in").Deref(), arg__8030_13, arg2})
	if callErr != nil {
		return nil, callErr
	}
	return v14, nil
}
func set_block_of_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__8038_7 vm.Value
	var arg__8047_13 vm.Value
	var v14 vm.Value
	var callErr error
	_, _, _ = arg__8038_7, arg__8047_13, v14
	arg__8038_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("insts"), arg1, vm.Keyword("block")})
	if callErr != nil {
		return nil, callErr
	}
	arg__8047_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("insts"), arg1, vm.Keyword("block")})
	if callErr != nil {
		return nil, callErr
	}
	v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "assoc-in").Deref(), arg__8047_13, arg2})
	if callErr != nil {
		return nil, callErr
	}
	return v14, nil
}
func block_set_params_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__8055_7 vm.Value
	var arg__8064_13 vm.Value
	var v14 vm.Value
	var callErr error
	_, _, _ = arg__8055_7, arg__8064_13, v14
	arg__8055_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("blocks"), arg1, vm.Keyword("params")})
	if callErr != nil {
		return nil, callErr
	}
	arg__8064_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("blocks"), arg1, vm.Keyword("params")})
	if callErr != nil {
		return nil, callErr
	}
	v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "assoc-in").Deref(), arg__8064_13, arg2})
	if callErr != nil {
		return nil, callErr
	}
	return v14, nil
}
func set_refs_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__8072_7 vm.Value
	var arg__8081_13 vm.Value
	var v14 vm.Value
	var v26 vm.Value
	var callErr error
	_, _, _, _ = arg__8072_7, arg__8081_13, v14, v26
	arg__8072_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("insts"), arg1, vm.Keyword("refs")})
	if callErr != nil {
		return nil, callErr
	}
	arg__8081_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("insts"), arg1, vm.Keyword("refs")})
	if callErr != nil {
		return nil, callErr
	}
	v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "assoc-in").Deref(), arg__8081_13, arg2})
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
	var arg__8101_5 vm.Value
	var arg__8102_6 vm.Value
	var arg__8109_10 vm.Value
	var arg__8110_11 vm.Value
	var arg__8112_12 vm.Value
	var v13 vm.Value
	var callErr error
	_, _, _, _, _, _ = arg__8101_5, arg__8102_6, arg__8109_10, arg__8110_11, arg__8112_12, v13
	arg__8101_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8102_6, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8101_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__8109_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8110_11, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8109_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__8112_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__8110_11, arg0})
	if callErr != nil {
		return nil, callErr
	}
	v13, callErr = rt.InvokeValue(vm.Keyword("refs"), []vm.Value{arg__8112_12})
	if callErr != nil {
		return nil, callErr
	}
	return v13, nil
}
func fn_uses_cache(arg0 vm.Value) (vm.Value, error) {
	var arg__8117_3 vm.Value
	var v4 vm.Value
	var callErr error
	_, _ = arg__8117_3, v4
	arg__8117_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v4, callErr = rt.InvokeValue(vm.Keyword("uses-cache"), []vm.Value{arg__8117_3})
	if callErr != nil {
		return nil, callErr
	}
	return v4, nil
}
func block_insts(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__8123_5 vm.Value
	var arg__8124_6 vm.Value
	var arg__8131_10 vm.Value
	var arg__8132_11 vm.Value
	var arg__8134_12 vm.Value
	var v13 vm.Value
	var callErr error
	_, _, _, _, _, _ = arg__8123_5, arg__8124_6, arg__8131_10, arg__8132_11, arg__8134_12, v13
	arg__8123_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8124_6, callErr = rt.InvokeValue(vm.Keyword("blocks"), []vm.Value{arg__8123_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__8131_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8132_11, callErr = rt.InvokeValue(vm.Keyword("blocks"), []vm.Value{arg__8131_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__8134_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__8132_11, arg0})
	if callErr != nil {
		return nil, callErr
	}
	v13, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8134_12})
	if callErr != nil {
		return nil, callErr
	}
	return v13, nil
}
func set_aux_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__8141_7 vm.Value
	var arg__8150_13 vm.Value
	var v14 vm.Value
	var callErr error
	_, _, _ = arg__8141_7, arg__8150_13, v14
	arg__8141_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("insts"), arg1, vm.Keyword("aux")})
	if callErr != nil {
		return nil, callErr
	}
	arg__8150_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("insts"), arg1, vm.Keyword("aux")})
	if callErr != nil {
		return nil, callErr
	}
	v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "assoc-in").Deref(), arg__8150_13, arg2})
	if callErr != nil {
		return nil, callErr
	}
	return v14, nil
}
func fn_consts(arg0 vm.Value) (vm.Value, error) {
	var arg__8156_3 vm.Value
	var v4 vm.Value
	var callErr error
	_, _ = arg__8156_3, v4
	arg__8156_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v4, callErr = rt.InvokeValue(vm.Keyword("consts"), []vm.Value{arg__8156_3})
	if callErr != nil {
		return nil, callErr
	}
	return v4, nil
}
func block_id(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__8162_5 vm.Value
	var arg__8163_6 vm.Value
	var arg__8170_10 vm.Value
	var arg__8171_11 vm.Value
	var arg__8173_12 vm.Value
	var v13 vm.Value
	var callErr error
	_, _, _, _, _, _ = arg__8162_5, arg__8163_6, arg__8170_10, arg__8171_11, arg__8173_12, v13
	arg__8162_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8163_6, callErr = rt.InvokeValue(vm.Keyword("blocks"), []vm.Value{arg__8162_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__8170_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8171_11, callErr = rt.InvokeValue(vm.Keyword("blocks"), []vm.Value{arg__8170_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__8173_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__8171_11, arg0})
	if callErr != nil {
		return nil, callErr
	}
	v13, callErr = rt.InvokeValue(vm.Keyword("id"), []vm.Value{arg__8173_12})
	if callErr != nil {
		return nil, callErr
	}
	return v13, nil
}
func op(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__8179_5 vm.Value
	var arg__8180_6 vm.Value
	var arg__8187_10 vm.Value
	var arg__8188_11 vm.Value
	var arg__8190_12 vm.Value
	var v13 vm.Value
	var callErr error
	_, _, _, _, _, _ = arg__8179_5, arg__8180_6, arg__8187_10, arg__8188_11, arg__8190_12, v13
	arg__8179_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8180_6, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8179_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__8187_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8188_11, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8187_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__8190_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__8188_11, arg0})
	if callErr != nil {
		return nil, callErr
	}
	v13, callErr = rt.InvokeValue(vm.Keyword("op"), []vm.Value{arg__8190_12})
	if callErr != nil {
		return nil, callErr
	}
	return v13, nil
}
func fn_name(arg0 vm.Value) (vm.Value, error) {
	var arg__8195_3 vm.Value
	var v4 vm.Value
	var callErr error
	_, _ = arg__8195_3, v4
	arg__8195_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v4, callErr = rt.InvokeValue(vm.Keyword("name"), []vm.Value{arg__8195_3})
	if callErr != nil {
		return nil, callErr
	}
	return v4, nil
}
func set_type_of_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__8202_7 vm.Value
	var arg__8211_13 vm.Value
	var v14 vm.Value
	var callErr error
	_, _, _ = arg__8202_7, arg__8211_13, v14
	arg__8202_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("insts"), arg1, vm.Keyword("type")})
	if callErr != nil {
		return nil, callErr
	}
	arg__8211_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("insts"), arg1, vm.Keyword("type")})
	if callErr != nil {
		return nil, callErr
	}
	v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "assoc-in").Deref(), arg__8211_13, arg2})
	if callErr != nil {
		return nil, callErr
	}
	return v14, nil
}
func fn_entry(arg0 vm.Value) (vm.Value, error) {
	var arg__8217_3 vm.Value
	var v4 vm.Value
	var callErr error
	_, _ = arg__8217_3, v4
	arg__8217_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v4, callErr = rt.InvokeValue(vm.Keyword("entry"), []vm.Value{arg__8217_3})
	if callErr != nil {
		return nil, callErr
	}
	return v4, nil
}
func fn_arity(arg0 vm.Value) (vm.Value, error) {
	var arg__8222_3 vm.Value
	var v4 vm.Value
	var callErr error
	_, _ = arg__8222_3, v4
	arg__8222_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v4, callErr = rt.InvokeValue(vm.Keyword("arity"), []vm.Value{arg__8222_3})
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
	var arg__8232_5 vm.Value
	var arg__8233_6 vm.Value
	var arg__8240_10 vm.Value
	var arg__8241_11 vm.Value
	var arg__8243_12 vm.Value
	var v13 vm.Value
	var callErr error
	_, _, _, _, _, _ = arg__8232_5, arg__8233_6, arg__8240_10, arg__8241_11, arg__8243_12, v13
	arg__8232_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8233_6, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8232_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__8240_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8241_11, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8240_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__8243_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__8241_11, arg0})
	if callErr != nil {
		return nil, callErr
	}
	v13, callErr = rt.InvokeValue(vm.Keyword("block"), []vm.Value{arg__8243_12})
	if callErr != nil {
		return nil, callErr
	}
	return v13, nil
}
func fn_insts(arg0 vm.Value) (vm.Value, error) {
	var arg__8248_3 vm.Value
	var v4 vm.Value
	var callErr error
	_, _ = arg__8248_3, v4
	arg__8248_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v4, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8248_3})
	if callErr != nil {
		return nil, callErr
	}
	return v4, nil
}
func set_source_infos_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__8255_7 vm.Value
	var arg__8264_13 vm.Value
	var v14 vm.Value
	var callErr error
	_, _, _ = arg__8255_7, arg__8264_13, v14
	arg__8255_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("insts"), arg1, vm.Keyword("source-infos")})
	if callErr != nil {
		return nil, callErr
	}
	arg__8264_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("insts"), arg1, vm.Keyword("source-infos")})
	if callErr != nil {
		return nil, callErr
	}
	v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "assoc-in").Deref(), arg__8264_13, arg2})
	if callErr != nil {
		return nil, callErr
	}
	return v14, nil
}
func set_op_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__8272_7 vm.Value
	var arg__8281_13 vm.Value
	var v14 vm.Value
	var v26 vm.Value
	var callErr error
	_, _, _, _ = arg__8272_7, arg__8281_13, v14, v26
	arg__8272_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("insts"), arg1, vm.Keyword("op")})
	if callErr != nil {
		return nil, callErr
	}
	arg__8281_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("insts"), arg1, vm.Keyword("op")})
	if callErr != nil {
		return nil, callErr
	}
	v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "assoc-in").Deref(), arg__8281_13, arg2})
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
	var arg__8301_5 vm.Value
	var arg__8302_6 vm.Value
	var arg__8309_10 vm.Value
	var arg__8310_11 vm.Value
	var arg__8312_12 vm.Value
	var v13 vm.Value
	var callErr error
	_, _, _, _, _, _ = arg__8301_5, arg__8302_6, arg__8309_10, arg__8310_11, arg__8312_12, v13
	arg__8301_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8302_6, callErr = rt.InvokeValue(vm.Keyword("blocks"), []vm.Value{arg__8301_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__8309_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8310_11, callErr = rt.InvokeValue(vm.Keyword("blocks"), []vm.Value{arg__8309_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__8312_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__8310_11, arg0})
	if callErr != nil {
		return nil, callErr
	}
	v13, callErr = rt.InvokeValue(vm.Keyword("preds"), []vm.Value{arg__8312_12})
	if callErr != nil {
		return nil, callErr
	}
	return v13, nil
}
func source_infos(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__8318_5 vm.Value
	var arg__8319_6 vm.Value
	var arg__8326_10 vm.Value
	var arg__8327_11 vm.Value
	var arg__8329_12 vm.Value
	var v13 vm.Value
	var callErr error
	_, _, _, _, _, _ = arg__8318_5, arg__8319_6, arg__8326_10, arg__8327_11, arg__8329_12, v13
	arg__8318_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8319_6, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8318_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__8326_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8327_11, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8326_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__8329_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__8327_11, arg0})
	if callErr != nil {
		return nil, callErr
	}
	v13, callErr = rt.InvokeValue(vm.Keyword("source-infos"), []vm.Value{arg__8329_12})
	if callErr != nil {
		return nil, callErr
	}
	return v13, nil
}
func aux(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__8335_5 vm.Value
	var arg__8336_6 vm.Value
	var arg__8343_10 vm.Value
	var arg__8344_11 vm.Value
	var arg__8346_12 vm.Value
	var v13 vm.Value
	var callErr error
	_, _, _, _, _, _ = arg__8335_5, arg__8336_6, arg__8343_10, arg__8344_11, arg__8346_12, v13
	arg__8335_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8336_6, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8335_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__8343_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8344_11, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__8343_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__8346_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__8344_11, arg0})
	if callErr != nil {
		return nil, callErr
	}
	v13, callErr = rt.InvokeValue(vm.Keyword("aux"), []vm.Value{arg__8346_12})
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
	var arg__8355_7 vm.Value
	var arg__8364_13 vm.Value
	var v14 vm.Value
	var callErr error
	_, _, _ = arg__8355_7, arg__8364_13, v14
	arg__8355_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("blocks"), arg1, vm.Keyword("id")})
	if callErr != nil {
		return nil, callErr
	}
	arg__8364_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("blocks"), arg1, vm.Keyword("id")})
	if callErr != nil {
		return nil, callErr
	}
	v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "assoc-in").Deref(), arg__8364_13, arg2})
	if callErr != nil {
		return nil, callErr
	}
	return v14, nil
}
func block_params(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__8371_5 vm.Value
	var arg__8372_6 vm.Value
	var arg__8379_10 vm.Value
	var arg__8380_11 vm.Value
	var arg__8382_12 vm.Value
	var v13 vm.Value
	var callErr error
	_, _, _, _, _, _ = arg__8371_5, arg__8372_6, arg__8379_10, arg__8380_11, arg__8382_12, v13
	arg__8371_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8372_6, callErr = rt.InvokeValue(vm.Keyword("blocks"), []vm.Value{arg__8371_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__8379_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8380_11, callErr = rt.InvokeValue(vm.Keyword("blocks"), []vm.Value{arg__8379_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__8382_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__8380_11, arg0})
	if callErr != nil {
		return nil, callErr
	}
	v13, callErr = rt.InvokeValue(vm.Keyword("params"), []vm.Value{arg__8382_12})
	if callErr != nil {
		return nil, callErr
	}
	return v13, nil
}
func block_set_term_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__8389_7 vm.Value
	var arg__8398_13 vm.Value
	var v14 vm.Value
	var callErr error
	_, _, _ = arg__8389_7, arg__8398_13, v14
	arg__8389_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("blocks"), arg1, vm.Keyword("term")})
	if callErr != nil {
		return nil, callErr
	}
	arg__8398_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("blocks"), arg1, vm.Keyword("term")})
	if callErr != nil {
		return nil, callErr
	}
	v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "assoc-in").Deref(), arg__8398_13, arg2})
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
	var arg__8406_3 vm.Value
	var v4 vm.Value
	var callErr error
	_, _ = arg__8406_3, v4
	arg__8406_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v4, callErr = rt.InvokeValue(vm.Keyword("variadic?"), []vm.Value{arg__8406_3})
	if callErr != nil {
		return nil, callErr
	}
	return v4, nil
}
func block_set_insts_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__8413_7 vm.Value
	var arg__8422_13 vm.Value
	var v14 vm.Value
	var callErr error
	_, _, _ = arg__8413_7, arg__8422_13, v14
	arg__8413_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("blocks"), arg1, vm.Keyword("insts")})
	if callErr != nil {
		return nil, callErr
	}
	arg__8422_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("blocks"), arg1, vm.Keyword("insts")})
	if callErr != nil {
		return nil, callErr
	}
	v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "assoc-in").Deref(), arg__8422_13, arg2})
	if callErr != nil {
		return nil, callErr
	}
	return v14, nil
}
func fn_uses_dirty_QMARK_(arg0 vm.Value) (vm.Value, error) {
	var arg__8428_3 vm.Value
	var v4 vm.Value
	var callErr error
	_, _ = arg__8428_3, v4
	arg__8428_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v4, callErr = rt.InvokeValue(vm.Keyword("uses-dirty?"), []vm.Value{arg__8428_3})
	if callErr != nil {
		return nil, callErr
	}
	return v4, nil
}
func fn_blocks(arg0 vm.Value) (vm.Value, error) {
	var arg__8433_3 vm.Value
	var v4 vm.Value
	var callErr error
	_, _ = arg__8433_3, v4
	arg__8433_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v4, callErr = rt.InvokeValue(vm.Keyword("blocks"), []vm.Value{arg__8433_3})
	if callErr != nil {
		return nil, callErr
	}
	return v4, nil
}
func block_term(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__8439_5 vm.Value
	var arg__8440_6 vm.Value
	var arg__8447_10 vm.Value
	var arg__8448_11 vm.Value
	var arg__8450_12 vm.Value
	var v13 vm.Value
	var callErr error
	_, _, _, _, _, _ = arg__8439_5, arg__8440_6, arg__8447_10, arg__8448_11, arg__8450_12, v13
	arg__8439_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8440_6, callErr = rt.InvokeValue(vm.Keyword("blocks"), []vm.Value{arg__8439_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__8447_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8448_11, callErr = rt.InvokeValue(vm.Keyword("blocks"), []vm.Value{arg__8447_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__8450_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__8448_11, arg0})
	if callErr != nil {
		return nil, callErr
	}
	v13, callErr = rt.InvokeValue(vm.Keyword("term"), []vm.Value{arg__8450_12})
	if callErr != nil {
		return nil, callErr
	}
	return v13, nil
}
