package ir_passes_infer_arg_types

import (
	rt "github.com/nooga/let-go/pkg/rt"
	vm "github.com/nooga/let-go/pkg/vm"
)

func join_types(arg0 vm.Value, arg1 vm.Value) vm.Value {
	var v8 bool
	var a_3 vm.Value
	var b_4 vm.Value
	var a_5 vm.Value
	var b_6 vm.Value
	var v16 bool
	var v48 vm.Value
	var a_49 vm.Value
	var b_50 vm.Value
	var a_11 vm.Value
	var b_12 vm.Value
	var a_13 vm.Value
	var b_14 vm.Value
	var v23 bool
	var v44 vm.Value
	var a_45 vm.Value
	var b_46 vm.Value
	var a_19 vm.Value
	var b_20 vm.Value
	var a_21 vm.Value
	var b_22 vm.Value
	var v40 vm.Value
	var a_41 vm.Value
	var b_42 vm.Value
	var a_26 vm.Value
	var b_27 vm.Value
	var a_28 vm.Value
	var b_29 vm.Value
	var v36 vm.Value
	var a_37 vm.Value
	var b_38 vm.Value
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = v8, a_3, b_4, a_5, b_6, v16, v48, a_49, b_50, a_11, b_12, a_13, b_14, v23, v44, a_45, b_46, a_19, b_20, a_21, b_22, v40, a_41, b_42, a_26, b_27, a_28, b_29, v36, a_37, b_38
	v8 = arg0 == vm.Keyword("unknown")
	if v8 {
		a_3 = arg0
		b_4 = arg1
		goto b1
	} else {
		a_5 = arg0
		b_6 = arg1
		goto b2
	}
b1:
	;
	v48 = b_4
	a_49 = a_3
	b_50 = b_4
	goto b3
b2:
	;
	v16 = b_6 == vm.Keyword("unknown")
	if v16 {
		a_11 = a_5
		b_12 = b_6
		goto b4
	} else {
		a_13 = a_5
		b_14 = b_6
		goto b5
	}
b3:
	;
	return v48
b4:
	;
	v44 = a_11
	a_45 = a_11
	b_46 = b_12
	goto b6
b5:
	;
	v23 = a_13 == b_14
	if v23 {
		a_19 = a_13
		b_20 = b_14
		goto b7
	} else {
		a_21 = a_13
		b_22 = b_14
		goto b8
	}
b6:
	;
	v48 = v44
	a_49 = a_45
	b_50 = b_46
	goto b3
b7:
	;
	v40 = a_19
	a_41 = a_19
	b_42 = b_20
	goto b9
b8:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		a_26 = a_21
		b_27 = b_22
		goto b10
	} else {
		a_28 = a_21
		b_29 = b_22
		goto b11
	}
b9:
	;
	v44 = v40
	a_45 = a_41
	b_46 = b_42
	goto b6
b10:
	;
	v36 = vm.Keyword("unknown")
	a_37 = a_26
	b_38 = b_27
	goto b12
b11:
	;
	v36 = vm.NIL
	a_37 = a_28
	b_38 = b_29
	goto b12
b12:
	;
	v40 = v36
	a_41 = a_37
	b_42 = b_38
	goto b9
}
func constraint_from_user(arg0 vm.Value) (vm.Value, error) {
	var v7 vm.Value
	var user_op_2 vm.Value
	var user_op_3 vm.Value
	var v22 vm.Value
	var user_op_23 vm.Value
	var user_op_11 vm.Value
	var user_op_12 vm.Value
	var v19 vm.Value
	var user_op_20 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _ = v7, user_op_2, user_op_3, v22, user_op_23, user_op_11, user_op_12, v19, user_op_20
	v7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "contains?").Deref(), []vm.Value{rt.LookupVar("ir.passes.infer-arg-types", "int-constraint-ops").Deref(), arg0})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v7) {
		user_op_2 = arg0
		goto b1
	} else {
		user_op_3 = arg0
		goto b2
	}
b1:
	;
	v22 = vm.Keyword("int")
	user_op_23 = user_op_2
	goto b3
b2:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		user_op_11 = user_op_3
		goto b4
	} else {
		user_op_12 = user_op_3
		goto b5
	}
b3:
	;
	return v22, nil
b4:
	;
	v19 = vm.Keyword("unknown")
	user_op_20 = user_op_11
	goto b6
b5:
	;
	v19 = vm.NIL
	user_op_20 = user_op_12
	goto b6
b6:
	;
	v22 = v19
	user_op_23 = user_op_20
	goto b3
}
func infer_one_load_arg_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var bs_5 vm.Value
	var joined_9 vm.Value
	var v21 vm.Value
	var f_10 vm.Value
	var use_index_11 vm.Value
	var nid_12 vm.Value
	var bs_13 vm.Value
	var joined_14 vm.Value
	var f_15 vm.Value
	var use_index_16 vm.Value
	var nid_17 vm.Value
	var bs_18 vm.Value
	var joined_19 vm.Value
	var v34 vm.Value
	var arg__18279_47 vm.Value
	var arg__18285_51 vm.Value
	var v52 vm.Value
	var v70 vm.Value
	var f_71 vm.Value
	var use_index_72 vm.Value
	var nid_73 vm.Value
	var bs_74 vm.Value
	var joined_75 vm.Value
	var f_35 vm.Value
	var use_index_36 vm.Value
	var nid_37 vm.Value
	var bs_38 vm.Value
	var joined_39 vm.Value
	var arg__18291_55 vm.Value
	var arg__18298_58 vm.Value
	var v59 vm.Value
	var f_40 vm.Value
	var use_index_41 vm.Value
	var nid_42 vm.Value
	var bs_43 vm.Value
	var joined_44 vm.Value
	var v63 vm.Value
	var f_64 vm.Value
	var use_index_65 vm.Value
	var nid_66 vm.Value
	var bs_67 vm.Value
	var joined_68 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = bs_5, joined_9, v21, f_10, use_index_11, nid_12, bs_13, joined_14, f_15, use_index_16, nid_17, bs_18, joined_19, v34, arg__18279_47, arg__18285_51, v52, v70, f_71, use_index_72, nid_73, bs_74, joined_75, f_35, use_index_36, nid_37, bs_38, joined_39, arg__18291_55, arg__18298_58, v59, f_40, use_index_41, nid_42, bs_43, joined_44, v63, f_64, use_index_65, nid_66, bs_67, joined_68
	bs_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg1, arg2})
	if callErr != nil {
		return nil, callErr
	}
	joined_9, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "atom").Deref(), []vm.Value{vm.Keyword("unknown")})
	if callErr != nil {
		return nil, callErr
	}
	v21, callErr = rt.InvokeValue(rt.LookupVar("ir", "uses-empty?").Deref(), []vm.Value{bs_5})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v21) {
		f_10 = arg0
		use_index_11 = arg1
		nid_12 = arg2
		bs_13 = bs_5
		joined_14 = joined_9
		goto b1
	} else {
		f_15 = arg0
		use_index_16 = arg1
		nid_17 = arg2
		bs_18 = bs_5
		joined_19 = joined_9
		goto b2
	}
b1:
	;
	v70 = vm.NIL
	f_71 = f_10
	use_index_72 = use_index_11
	nid_73 = nid_12
	bs_74 = bs_13
	joined_75 = joined_14
	goto b3
b2:
	;
	v34, callErr = rt.InvokeValue(rt.LookupVar("ir", "uses-for-each").Deref(), []vm.Value{bs_18, rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var user_op_4 vm.Value
		var c_6 vm.Value
		var v10 vm.Value
		var callErr error
		_, _, _ = user_op_4, c_6, v10
		user_op_4, callErr = rt.InvokeValue(rt.LookupVar("ir", "op").Deref(), []vm.Value{arg0, f_15})
		if callErr != nil {
			return nil, callErr
		}
		c_6, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.infer-arg-types", "constraint-from-user").Deref(), []vm.Value{user_op_4})
		if callErr != nil {
			return nil, callErr
		}
		v10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{joined_19, rt.LookupVar("ir.passes.infer-arg-types", "join-types").Deref(), c_6})
		if callErr != nil {
			return nil, callErr
		}
		return v10, nil
	})})
	if callErr != nil {
		return nil, callErr
	}
	arg__18279_47, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{joined_19})
	if callErr != nil {
		return nil, callErr
	}
	arg__18285_51, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{joined_19})
	if callErr != nil {
		return nil, callErr
	}
	v52, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not=").Deref(), []vm.Value{vm.Keyword("unknown"), arg__18285_51})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v52) {
		f_35 = f_15
		use_index_36 = use_index_16
		nid_37 = nid_17
		bs_38 = bs_18
		joined_39 = joined_19
		goto b4
	} else {
		f_40 = f_15
		use_index_41 = use_index_16
		nid_42 = nid_17
		bs_43 = bs_18
		joined_44 = joined_19
		goto b5
	}
b3:
	;
	return v70, nil
b4:
	;
	arg__18291_55, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{joined_39})
	if callErr != nil {
		return nil, callErr
	}
	arg__18298_58, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{joined_39})
	if callErr != nil {
		return nil, callErr
	}
	v59, callErr = rt.InvokeValue(rt.LookupVar("ir", "set-type!").Deref(), []vm.Value{f_35, nid_37, arg__18298_58})
	if callErr != nil {
		return nil, callErr
	}
	v63 = v59
	f_64 = f_35
	use_index_65 = use_index_36
	nid_66 = nid_37
	bs_67 = bs_38
	joined_68 = joined_39
	goto b6
b5:
	;
	v63 = vm.NIL
	f_64 = f_40
	use_index_65 = use_index_41
	nid_66 = nid_42
	bs_67 = bs_43
	joined_68 = joined_44
	goto b6
b6:
	;
	v70 = v63
	f_71 = f_64
	use_index_72 = use_index_65
	nid_73 = nid_66
	bs_74 = bs_67
	joined_75 = joined_68
	goto b3
}
func infer_arg_types(arg0 vm.Value) (vm.Value, error) {
	var entry_bid_3 vm.Value
	var insts_5 vm.Value
	var use_index_7 vm.Value
	var doseq_seq__18299_9 vm.Value
	var doseq_loop__18300_10 vm.Value
	var f_11 vm.Value
	var use_index_12 vm.Value
	var v77 vm.Value
	var entry_bid_14 vm.Value
	var insts_15 vm.Value
	var doseq_seq__18299_16 vm.Value
	var doseq_loop__18300_17 vm.Value
	var f_18 vm.Value
	var use_index_19 vm.Value
	var v76 vm.Value
	var nid_28 vm.Value
	var arg__18324_45 vm.Value
	var v46 bool
	var entry_bid_20 vm.Value
	var insts_21 vm.Value
	var doseq_seq__18299_22 vm.Value
	var doseq_loop__18300_23 vm.Value
	var f_24 vm.Value
	var use_index_25 vm.Value
	var v81 vm.Value
	var v66 vm.Value
	var entry_bid_67 vm.Value
	var insts_68 vm.Value
	var doseq_seq__18299_69 vm.Value
	var doseq_loop__18300_70 vm.Value
	var f_71 vm.Value
	var use_index_72 vm.Value
	var entry_bid_29 vm.Value
	var insts_30 vm.Value
	var doseq_seq__18299_31 vm.Value
	var doseq_loop__18300_32 vm.Value
	var f_33 vm.Value
	var use_index_34 vm.Value
	var nid_35 vm.Value
	var v79 vm.Value
	var v49 vm.Value
	var entry_bid_36 vm.Value
	var insts_37 vm.Value
	var doseq_seq__18299_38 vm.Value
	var doseq_loop__18300_39 vm.Value
	var f_40 vm.Value
	var use_index_41 vm.Value
	var nid_42 vm.Value
	var v78 vm.Value
	var v53 vm.Value
	var entry_bid_54 vm.Value
	var insts_55 vm.Value
	var doseq_seq__18299_56 vm.Value
	var doseq_loop__18300_57 vm.Value
	var f_58 vm.Value
	var use_index_59 vm.Value
	var nid_60 vm.Value
	var v80 vm.Value
	var v62 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = entry_bid_3, insts_5, use_index_7, doseq_seq__18299_9, doseq_loop__18300_10, f_11, use_index_12, v77, entry_bid_14, insts_15, doseq_seq__18299_16, doseq_loop__18300_17, f_18, use_index_19, v76, nid_28, arg__18324_45, v46, entry_bid_20, insts_21, doseq_seq__18299_22, doseq_loop__18300_23, f_24, use_index_25, v81, v66, entry_bid_67, insts_68, doseq_seq__18299_69, doseq_loop__18300_70, f_71, use_index_72, entry_bid_29, insts_30, doseq_seq__18299_31, doseq_loop__18300_32, f_33, use_index_34, nid_35, v79, v49, entry_bid_36, insts_37, doseq_seq__18299_38, doseq_loop__18300_39, f_40, use_index_41, nid_42, v78, v53, entry_bid_54, insts_55, doseq_seq__18299_56, doseq_loop__18300_57, f_58, use_index_59, nid_60, v80, v62
	entry_bid_3, callErr = rt.InvokeValue(rt.LookupVar("ir", "entry-block").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	insts_5, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-insts").Deref(), []vm.Value{entry_bid_3, arg0})
	if callErr != nil {
		return nil, callErr
	}
	use_index_7, callErr = rt.InvokeValue(rt.LookupVar("ir", "uses").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	doseq_seq__18299_9, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{insts_5})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__18300_10 = doseq_seq__18299_9
	f_11 = arg0
	use_index_12 = use_index_7
	v77 = vm.Keyword("load-arg")
	goto b1
b1:
	;
	if vm.IsTruthy(doseq_loop__18300_10) {
		entry_bid_14 = entry_bid_3
		insts_15 = insts_5
		doseq_seq__18299_16 = doseq_seq__18299_9
		doseq_loop__18300_17 = doseq_loop__18300_10
		f_18 = f_11
		use_index_19 = use_index_12
		v76 = v77
		goto b2
	} else {
		entry_bid_20 = entry_bid_3
		insts_21 = insts_5
		doseq_seq__18299_22 = doseq_seq__18299_9
		doseq_loop__18300_23 = doseq_loop__18300_10
		f_24 = f_11
		use_index_25 = use_index_12
		v81 = v77
		goto b3
	}
b2:
	;
	nid_28, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{doseq_loop__18300_17})
	if callErr != nil {
		return nil, callErr
	}
	arg__18324_45, callErr = rt.InvokeValue(rt.LookupVar("ir", "op").Deref(), []vm.Value{nid_28, f_18})
	if callErr != nil {
		return nil, callErr
	}
	v46 = arg__18324_45 == v76
	if v46 {
		entry_bid_29 = entry_bid_14
		insts_30 = insts_15
		doseq_seq__18299_31 = doseq_seq__18299_16
		doseq_loop__18300_32 = doseq_loop__18300_17
		f_33 = f_18
		use_index_34 = use_index_19
		nid_35 = nid_28
		v79 = v76
		goto b5
	} else {
		entry_bid_36 = entry_bid_14
		insts_37 = insts_15
		doseq_seq__18299_38 = doseq_seq__18299_16
		doseq_loop__18300_39 = doseq_loop__18300_17
		f_40 = f_18
		use_index_41 = use_index_19
		nid_42 = nid_28
		v78 = v76
		goto b6
	}
b3:
	;
	v66 = vm.NIL
	entry_bid_67 = entry_bid_20
	insts_68 = insts_21
	doseq_seq__18299_69 = doseq_seq__18299_22
	doseq_loop__18300_70 = doseq_loop__18300_23
	f_71 = f_24
	use_index_72 = use_index_25
	goto b4
b4:
	;
	return f_71, nil
b5:
	;
	v49, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.infer-arg-types", "infer-one-load-arg!").Deref(), []vm.Value{f_33, use_index_34, nid_35})
	if callErr != nil {
		return nil, callErr
	}
	v53 = v49
	entry_bid_54 = entry_bid_29
	insts_55 = insts_30
	doseq_seq__18299_56 = doseq_seq__18299_31
	doseq_loop__18300_57 = doseq_loop__18300_32
	f_58 = f_33
	use_index_59 = use_index_34
	nid_60 = nid_35
	v80 = v79
	goto b7
b6:
	;
	v53 = vm.NIL
	entry_bid_54 = entry_bid_36
	insts_55 = insts_37
	doseq_seq__18299_56 = doseq_seq__18299_38
	doseq_loop__18300_57 = doseq_loop__18300_39
	f_58 = f_40
	use_index_59 = use_index_41
	nid_60 = nid_42
	v80 = v78
	goto b7
b7:
	;
	v62, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "next").Deref(), []vm.Value{doseq_loop__18300_57})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__18300_10 = v62
	f_11 = f_58
	use_index_12 = use_index_59
	v77 = v80
	goto b1
}
