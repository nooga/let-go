package ir_passes_mutability

import (
	rt "github.com/nooga/let-go/pkg/rt"
	vm "github.com/nooga/let-go/pkg/vm"
)

func mutating_var_call_QMARK_(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__17721_9 vm.Value
	var v10 bool
	var nid_3 vm.Value
	var f_4 vm.Value
	var refs_13 vm.Value
	var v21 vm.Value
	var nid_5 vm.Value
	var f_6 vm.Value
	var v93 vm.Value
	var nid_94 vm.Value
	var f_95 vm.Value
	var nid_14 vm.Value
	var f_15 vm.Value
	var refs_16 vm.Value
	var v24 vm.Value
	var nid_17 vm.Value
	var f_18 vm.Value
	var refs_19 vm.Value
	var callee_28 vm.Value
	var nid_29 vm.Value
	var f_30 vm.Value
	var refs_31 vm.Value
	var and__x_32 vm.Value
	var callee_33 vm.Value
	var nid_34 vm.Value
	var f_35 vm.Value
	var refs_36 vm.Value
	var arg__17739_45 vm.Value
	var and__x_46 bool
	var and__x_37 vm.Value
	var callee_38 vm.Value
	var nid_39 vm.Value
	var f_40 vm.Value
	var refs_41 vm.Value
	var v84 vm.Value
	var and__x_85 vm.Value
	var callee_86 vm.Value
	var nid_87 vm.Value
	var f_88 vm.Value
	var refs_89 vm.Value
	var callee_47 vm.Value
	var nid_48 vm.Value
	var f_49 vm.Value
	var refs_50 vm.Value
	var and__x_51 bool
	var arg__17746_60 vm.Value
	var arg__17753_63 vm.Value
	var arg__17754_64 vm.Value
	var arg__17762_68 vm.Value
	var arg__17769_71 vm.Value
	var arg__17770_72 vm.Value
	var v73 vm.Value
	var callee_52 vm.Value
	var nid_53 vm.Value
	var f_54 vm.Value
	var refs_55 vm.Value
	var and__x_56 bool
	var v76 vm.Value
	var callee_77 vm.Value
	var nid_78 vm.Value
	var f_79 vm.Value
	var refs_80 vm.Value
	var and__x_81 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__17721_9, v10, nid_3, f_4, refs_13, v21, nid_5, f_6, v93, nid_94, f_95, nid_14, f_15, refs_16, v24, nid_17, f_18, refs_19, callee_28, nid_29, f_30, refs_31, and__x_32, callee_33, nid_34, f_35, refs_36, arg__17739_45, and__x_46, and__x_37, callee_38, nid_39, f_40, refs_41, v84, and__x_85, callee_86, nid_87, f_88, refs_89, callee_47, nid_48, f_49, refs_50, and__x_51, arg__17746_60, arg__17753_63, arg__17754_64, arg__17762_68, arg__17769_71, arg__17770_72, v73, callee_52, nid_53, f_54, refs_55, and__x_56, v76, callee_77, nid_78, f_79, refs_80, and__x_81
	arg__17721_9, callErr = rt.InvokeValue(rt.LookupVar("ir", "op").Deref(), []vm.Value{arg0, arg1})
	if callErr != nil {
		return nil, callErr
	}
	v10 = arg__17721_9 == vm.Keyword("call")
	if v10 {
		nid_3 = arg0
		f_4 = arg1
		goto b1
	} else {
		nid_5 = arg0
		f_6 = arg1
		goto b2
	}
b1:
	;
	refs_13, callErr = rt.InvokeValue(rt.LookupVar("ir", "refs").Deref(), []vm.Value{nid_3, f_4})
	if callErr != nil {
		return nil, callErr
	}
	v21, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{refs_13})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v21) {
		nid_14 = nid_3
		f_15 = f_4
		refs_16 = refs_13
		goto b4
	} else {
		nid_17 = nid_3
		f_18 = f_4
		refs_19 = refs_13
		goto b5
	}
b2:
	;
	v93 = vm.NIL
	nid_94 = nid_5
	f_95 = f_6
	goto b3
b3:
	;
	return v93, nil
b4:
	;
	v24, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{refs_16})
	if callErr != nil {
		return nil, callErr
	}
	callee_28 = v24
	nid_29 = nid_14
	f_30 = f_15
	refs_31 = refs_16
	goto b6
b5:
	;
	callee_28 = vm.NIL
	nid_29 = nid_17
	f_30 = f_18
	refs_31 = refs_19
	goto b6
b6:
	;
	if vm.IsTruthy(callee_28) {
		and__x_32 = callee_28
		callee_33 = callee_28
		nid_34 = nid_29
		f_35 = f_30
		refs_36 = refs_31
		goto b7
	} else {
		and__x_37 = callee_28
		callee_38 = callee_28
		nid_39 = nid_29
		f_40 = f_30
		refs_41 = refs_31
		goto b8
	}
b7:
	;
	arg__17739_45, callErr = rt.InvokeValue(rt.LookupVar("ir", "op").Deref(), []vm.Value{callee_33, f_35})
	if callErr != nil {
		return nil, callErr
	}
	and__x_46 = arg__17739_45 == vm.Keyword("load-var")
	if and__x_46 {
		callee_47 = callee_33
		nid_48 = nid_34
		f_49 = f_35
		refs_50 = refs_36
		and__x_51 = and__x_46
		goto b10
	} else {
		callee_52 = callee_33
		nid_53 = nid_34
		f_54 = f_35
		refs_55 = refs_36
		and__x_56 = and__x_46
		goto b11
	}
b8:
	;
	v84 = and__x_37
	and__x_85 = and__x_37
	callee_86 = callee_38
	nid_87 = nid_39
	f_88 = f_40
	refs_89 = refs_41
	goto b9
b9:
	;
	v93 = v84
	nid_94 = nid_87
	f_95 = f_88
	goto b3
b10:
	;
	arg__17746_60, callErr = rt.InvokeValue(rt.LookupVar("ir", "aux").Deref(), []vm.Value{callee_47, f_49})
	if callErr != nil {
		return nil, callErr
	}
	arg__17753_63, callErr = rt.InvokeValue(rt.LookupVar("ir", "aux").Deref(), []vm.Value{callee_47, f_49})
	if callErr != nil {
		return nil, callErr
	}
	arg__17754_64, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{arg__17753_63})
	if callErr != nil {
		return nil, callErr
	}
	arg__17762_68, callErr = rt.InvokeValue(rt.LookupVar("ir", "aux").Deref(), []vm.Value{callee_47, f_49})
	if callErr != nil {
		return nil, callErr
	}
	arg__17769_71, callErr = rt.InvokeValue(rt.LookupVar("ir", "aux").Deref(), []vm.Value{callee_47, f_49})
	if callErr != nil {
		return nil, callErr
	}
	arg__17770_72, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{arg__17769_71})
	if callErr != nil {
		return nil, callErr
	}
	v73, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "contains?").Deref(), []vm.Value{rt.LookupVar("ir.passes.mutability", "known-var-mutating-builtins").Deref(), arg__17770_72})
	if callErr != nil {
		return nil, callErr
	}
	v76 = v73
	callee_77 = callee_47
	nid_78 = nid_48
	f_79 = f_49
	refs_80 = refs_50
	and__x_81 = vm.Boolean(and__x_51)
	goto b12
b11:
	;
	v76 = vm.Boolean(and__x_56)
	callee_77 = callee_52
	nid_78 = nid_53
	f_79 = f_54
	refs_80 = refs_55
	and__x_81 = vm.Boolean(and__x_56)
	goto b12
b12:
	;
	v84 = v76
	and__x_85 = and__x_32
	callee_86 = callee_77
	nid_87 = nid_78
	f_88 = f_79
	refs_89 = refs_80
	goto b9
}
func stable_load_var_QMARK_(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__17773_4 vm.Value
	var arg__17777_7 vm.Value
	var and__x_8 vm.Value
	var facts_9 vm.Value
	var var_value_10 vm.Value
	var and__x_11 vm.Value
	var arg__17780_17 vm.Value
	var arg__17785_20 vm.Value
	var arg__17787_21 vm.Value
	var arg__17791_24 vm.Value
	var arg__17796_27 vm.Value
	var arg__17798_28 vm.Value
	var v29 vm.Value
	var facts_12 vm.Value
	var var_value_13 vm.Value
	var and__x_14 vm.Value
	var v32 vm.Value
	var facts_33 vm.Value
	var var_value_34 vm.Value
	var and__x_35 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__17773_4, arg__17777_7, and__x_8, facts_9, var_value_10, and__x_11, arg__17780_17, arg__17785_20, arg__17787_21, arg__17791_24, arg__17796_27, arg__17798_28, v29, facts_12, var_value_13, and__x_14, v32, facts_33, var_value_34, and__x_35
	arg__17773_4, callErr = rt.InvokeValue(vm.Keyword("unknown-all?"), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__17777_7, callErr = rt.InvokeValue(vm.Keyword("unknown-all?"), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	and__x_8, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not").Deref(), []vm.Value{arg__17777_7})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(and__x_8) {
		facts_9 = arg0
		var_value_10 = arg1
		and__x_11 = and__x_8
		goto b1
	} else {
		facts_12 = arg0
		var_value_13 = arg1
		and__x_14 = and__x_8
		goto b2
	}
b1:
	;
	arg__17780_17, callErr = rt.InvokeValue(vm.Keyword("mutated-vars"), []vm.Value{facts_9})
	if callErr != nil {
		return nil, callErr
	}
	arg__17785_20, callErr = rt.InvokeValue(vm.Keyword("mutated-vars"), []vm.Value{facts_9})
	if callErr != nil {
		return nil, callErr
	}
	arg__17787_21, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "contains?").Deref(), []vm.Value{arg__17785_20, var_value_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__17791_24, callErr = rt.InvokeValue(vm.Keyword("mutated-vars"), []vm.Value{facts_9})
	if callErr != nil {
		return nil, callErr
	}
	arg__17796_27, callErr = rt.InvokeValue(vm.Keyword("mutated-vars"), []vm.Value{facts_9})
	if callErr != nil {
		return nil, callErr
	}
	arg__17798_28, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "contains?").Deref(), []vm.Value{arg__17796_27, var_value_10})
	if callErr != nil {
		return nil, callErr
	}
	v29, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not").Deref(), []vm.Value{arg__17798_28})
	if callErr != nil {
		return nil, callErr
	}
	v32 = v29
	facts_33 = facts_9
	var_value_34 = var_value_10
	and__x_35 = and__x_11
	goto b3
b2:
	;
	v32 = and__x_14
	facts_33 = facts_12
	var_value_34 = var_value_13
	and__x_35 = and__x_14
	goto b3
b3:
	;
	return v32, nil
}
func analyze_var_stability(arg0 vm.Value) (vm.Value, error) {
	var arg__17803_8 vm.Value
	var arg__17804_9 vm.Value
	var arg__17810_13 vm.Value
	var arg__17811_14 vm.Value
	var arg__17812_15 vm.Value
	var arg__17818_19 vm.Value
	var arg__17819_20 vm.Value
	var arg__17825_24 vm.Value
	var arg__17826_25 vm.Value
	var arg__17827_26 vm.Value
	var v27 vm.Value
	var v29 vm.Value
	var ids_2 vm.Value
	var mutated_vars_3 vm.Value
	var unknown_all_QMARK__4 vm.Value
	var f_5 vm.Value
	var ids_32 vm.Value
	var mutated_vars_33 vm.Value
	var unknown_all_QMARK__34 vm.Value
	var f_35 vm.Value
	var v65 vm.Value
	var ids_36 vm.Value
	var mutated_vars_37 vm.Value
	var unknown_all_QMARK__38 vm.Value
	var f_39 vm.Value
	var nid_68 vm.Value
	var op_70 vm.Value
	var v84 bool
	var v157 vm.Value
	var ids_158 vm.Value
	var mutated_vars_159 vm.Value
	var unknown_all_QMARK__160 vm.Value
	var f_161 vm.Value
	var ids_40 vm.Value
	var mutated_vars_41 vm.Value
	var or__x_42 bool
	var unknown_all_QMARK__43 bool
	var f_44 vm.Value
	var ids_45 vm.Value
	var mutated_vars_46 vm.Value
	var or__x_47 bool
	var unknown_all_QMARK__48 bool
	var f_49 vm.Value
	var v53 vm.Value
	var v55 vm.Value
	var ids_56 vm.Value
	var mutated_vars_57 vm.Value
	var or__x_58 vm.Value
	var unknown_all_QMARK__59 vm.Value
	var f_60 vm.Value
	var ids_71 vm.Value
	var mutated_vars_72 vm.Value
	var unknown_all_QMARK__73 vm.Value
	var f_74 vm.Value
	var nid_75 vm.Value
	var op_76 vm.Value
	var v87 vm.Value
	var arg__17856_89 vm.Value
	var arg__17864_92 vm.Value
	var v93 vm.Value
	var ids_77 vm.Value
	var mutated_vars_78 vm.Value
	var unknown_all_QMARK__79 vm.Value
	var f_80 vm.Value
	var nid_81 vm.Value
	var op_82 vm.Value
	var v108 vm.Value
	var v149 vm.Value
	var ids_150 vm.Value
	var mutated_vars_151 vm.Value
	var unknown_all_QMARK__152 vm.Value
	var f_153 vm.Value
	var nid_154 vm.Value
	var op_155 vm.Value
	var ids_95 vm.Value
	var mutated_vars_96 vm.Value
	var unknown_all_QMARK__97 vm.Value
	var f_98 vm.Value
	var nid_99 vm.Value
	var op_100 vm.Value
	var v111 vm.Value
	var ids_101 vm.Value
	var mutated_vars_102 vm.Value
	var unknown_all_QMARK__103 vm.Value
	var f_104 vm.Value
	var nid_105 vm.Value
	var op_106 vm.Value
	var v141 vm.Value
	var ids_142 vm.Value
	var mutated_vars_143 vm.Value
	var unknown_all_QMARK__144 vm.Value
	var f_145 vm.Value
	var nid_146 vm.Value
	var op_147 vm.Value
	var ids_114 vm.Value
	var mutated_vars_115 vm.Value
	var unknown_all_QMARK__116 vm.Value
	var f_117 vm.Value
	var nid_118 vm.Value
	var op_119 vm.Value
	var v129 vm.Value
	var ids_120 vm.Value
	var mutated_vars_121 vm.Value
	var unknown_all_QMARK__122 vm.Value
	var f_123 vm.Value
	var nid_124 vm.Value
	var op_125 vm.Value
	var v133 vm.Value
	var ids_134 vm.Value
	var mutated_vars_135 vm.Value
	var unknown_all_QMARK__136 vm.Value
	var f_137 vm.Value
	var nid_138 vm.Value
	var op_139 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__17803_8, arg__17804_9, arg__17810_13, arg__17811_14, arg__17812_15, arg__17818_19, arg__17819_20, arg__17825_24, arg__17826_25, arg__17827_26, v27, v29, ids_2, mutated_vars_3, unknown_all_QMARK__4, f_5, ids_32, mutated_vars_33, unknown_all_QMARK__34, f_35, v65, ids_36, mutated_vars_37, unknown_all_QMARK__38, f_39, nid_68, op_70, v84, v157, ids_158, mutated_vars_159, unknown_all_QMARK__160, f_161, ids_40, mutated_vars_41, or__x_42, unknown_all_QMARK__43, f_44, ids_45, mutated_vars_46, or__x_47, unknown_all_QMARK__48, f_49, v53, v55, ids_56, mutated_vars_57, or__x_58, unknown_all_QMARK__59, f_60, ids_71, mutated_vars_72, unknown_all_QMARK__73, f_74, nid_75, op_76, v87, arg__17856_89, arg__17864_92, v93, ids_77, mutated_vars_78, unknown_all_QMARK__79, f_80, nid_81, op_82, v108, v149, ids_150, mutated_vars_151, unknown_all_QMARK__152, f_153, nid_154, op_155, ids_95, mutated_vars_96, unknown_all_QMARK__97, f_98, nid_99, op_100, v111, ids_101, mutated_vars_102, unknown_all_QMARK__103, f_104, nid_105, op_106, v141, ids_142, mutated_vars_143, unknown_all_QMARK__144, f_145, nid_146, op_147, ids_114, mutated_vars_115, unknown_all_QMARK__116, f_117, nid_118, op_119, v129, ids_120, mutated_vars_121, unknown_all_QMARK__122, f_123, nid_124, op_125, v133, ids_134, mutated_vars_135, unknown_all_QMARK__136, f_137, nid_138, op_139
	arg__17803_8, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__17804_9, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__17803_8})
	if callErr != nil {
		return nil, callErr
	}
	arg__17810_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__17811_14, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__17810_13})
	if callErr != nil {
		return nil, callErr
	}
	arg__17812_15, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{arg__17811_14})
	if callErr != nil {
		return nil, callErr
	}
	arg__17818_19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__17819_20, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__17818_19})
	if callErr != nil {
		return nil, callErr
	}
	arg__17825_24, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__17826_25, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__17825_24})
	if callErr != nil {
		return nil, callErr
	}
	arg__17827_26, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{arg__17826_25})
	if callErr != nil {
		return nil, callErr
	}
	v27, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "range").Deref(), []vm.Value{arg__17827_26})
	if callErr != nil {
		return nil, callErr
	}
	v29, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	ids_2 = v27
	mutated_vars_3 = v29
	unknown_all_QMARK__4 = vm.Boolean(false)
	f_5 = arg0
	goto b1
b1:
	;
	if vm.IsTruthy(unknown_all_QMARK__4) {
		ids_40 = ids_2
		mutated_vars_41 = mutated_vars_3
		or__x_42 = vm.IsTruthy(unknown_all_QMARK__4)
		unknown_all_QMARK__43 = vm.IsTruthy(unknown_all_QMARK__4)
		f_44 = f_5
		goto b5
	} else {
		ids_45 = ids_2
		mutated_vars_46 = mutated_vars_3
		or__x_47 = vm.IsTruthy(unknown_all_QMARK__4)
		unknown_all_QMARK__48 = vm.IsTruthy(unknown_all_QMARK__4)
		f_49 = f_5
		goto b6
	}
b2:
	;
	v65, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("mutated-vars"), mutated_vars_33, vm.Keyword("unknown-all?"), unknown_all_QMARK__34})
	if callErr != nil {
		return nil, callErr
	}
	v157 = v65
	ids_158 = ids_32
	mutated_vars_159 = mutated_vars_33
	unknown_all_QMARK__160 = unknown_all_QMARK__34
	f_161 = f_35
	goto b4
b3:
	;
	nid_68, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{ids_36})
	if callErr != nil {
		return nil, callErr
	}
	op_70, callErr = rt.InvokeValue(rt.LookupVar("ir", "op").Deref(), []vm.Value{nid_68, f_39})
	if callErr != nil {
		return nil, callErr
	}
	v84 = op_70 == vm.Keyword("set-var")
	if v84 {
		ids_71 = ids_36
		mutated_vars_72 = mutated_vars_37
		unknown_all_QMARK__73 = unknown_all_QMARK__38
		f_74 = f_39
		nid_75 = nid_68
		op_76 = op_70
		goto b8
	} else {
		ids_77 = ids_36
		mutated_vars_78 = mutated_vars_37
		unknown_all_QMARK__79 = unknown_all_QMARK__38
		f_80 = f_39
		nid_81 = nid_68
		op_82 = op_70
		goto b9
	}
b4:
	;
	return v157, nil
b5:
	;
	v55 = vm.Boolean(or__x_42)
	ids_56 = ids_40
	mutated_vars_57 = mutated_vars_41
	or__x_58 = vm.Boolean(or__x_42)
	unknown_all_QMARK__59 = vm.Boolean(unknown_all_QMARK__43)
	f_60 = f_44
	goto b7
b6:
	;
	v53, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "empty?").Deref(), []vm.Value{ids_45})
	if callErr != nil {
		return nil, callErr
	}
	v55 = v53
	ids_56 = ids_45
	mutated_vars_57 = mutated_vars_46
	or__x_58 = vm.Boolean(or__x_47)
	unknown_all_QMARK__59 = vm.Boolean(unknown_all_QMARK__48)
	f_60 = f_49
	goto b7
b7:
	;
	if vm.IsTruthy(v55) {
		ids_32 = ids_56
		mutated_vars_33 = mutated_vars_57
		unknown_all_QMARK__34 = unknown_all_QMARK__59
		f_35 = f_60
		goto b2
	} else {
		ids_36 = ids_56
		mutated_vars_37 = mutated_vars_57
		unknown_all_QMARK__38 = unknown_all_QMARK__59
		f_39 = f_60
		goto b3
	}
b8:
	;
	v87, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{ids_71})
	if callErr != nil {
		return nil, callErr
	}
	arg__17856_89, callErr = rt.InvokeValue(rt.LookupVar("ir", "aux").Deref(), []vm.Value{nid_75, f_74})
	if callErr != nil {
		return nil, callErr
	}
	arg__17864_92, callErr = rt.InvokeValue(rt.LookupVar("ir", "aux").Deref(), []vm.Value{nid_75, f_74})
	if callErr != nil {
		return nil, callErr
	}
	v93, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{mutated_vars_72, arg__17864_92})
	if callErr != nil {
		return nil, callErr
	}
	ids_2 = v87
	mutated_vars_3 = v93
	unknown_all_QMARK__4 = unknown_all_QMARK__73
	f_5 = f_74
	goto b1
b9:
	;
	v108, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.mutability", "mutating-var-call?").Deref(), []vm.Value{nid_81, f_80})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v108) {
		ids_95 = ids_77
		mutated_vars_96 = mutated_vars_78
		unknown_all_QMARK__97 = unknown_all_QMARK__79
		f_98 = f_80
		nid_99 = nid_81
		op_100 = op_82
		goto b11
	} else {
		ids_101 = ids_77
		mutated_vars_102 = mutated_vars_78
		unknown_all_QMARK__103 = unknown_all_QMARK__79
		f_104 = f_80
		nid_105 = nid_81
		op_106 = op_82
		goto b12
	}
b10:
	;
	v157 = v149
	ids_158 = ids_150
	mutated_vars_159 = mutated_vars_151
	unknown_all_QMARK__160 = unknown_all_QMARK__152
	f_161 = f_153
	goto b4
b11:
	;
	v111, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{ids_95})
	if callErr != nil {
		return nil, callErr
	}
	ids_2 = v111
	mutated_vars_3 = mutated_vars_96
	unknown_all_QMARK__4 = vm.Boolean(true)
	f_5 = f_98
	goto b1
b12:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		ids_114 = ids_101
		mutated_vars_115 = mutated_vars_102
		unknown_all_QMARK__116 = unknown_all_QMARK__103
		f_117 = f_104
		nid_118 = nid_105
		op_119 = op_106
		goto b14
	} else {
		ids_120 = ids_101
		mutated_vars_121 = mutated_vars_102
		unknown_all_QMARK__122 = unknown_all_QMARK__103
		f_123 = f_104
		nid_124 = nid_105
		op_125 = op_106
		goto b15
	}
b13:
	;
	v149 = v141
	ids_150 = ids_142
	mutated_vars_151 = mutated_vars_143
	unknown_all_QMARK__152 = unknown_all_QMARK__144
	f_153 = f_145
	nid_154 = nid_146
	op_155 = op_147
	goto b10
b14:
	;
	v129, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{ids_114})
	if callErr != nil {
		return nil, callErr
	}
	ids_2 = v129
	mutated_vars_3 = mutated_vars_115
	unknown_all_QMARK__4 = unknown_all_QMARK__116
	f_5 = f_117
	goto b1
b15:
	;
	v133 = vm.NIL
	ids_134 = ids_120
	mutated_vars_135 = mutated_vars_121
	unknown_all_QMARK__136 = unknown_all_QMARK__122
	f_137 = f_123
	nid_138 = nid_124
	op_139 = op_125
	goto b16
b16:
	;
	v141 = v133
	ids_142 = ids_134
	mutated_vars_143 = mutated_vars_135
	unknown_all_QMARK__144 = unknown_all_QMARK__136
	f_145 = f_137
	nid_146 = nid_138
	op_147 = op_139
	goto b13
}
