package ir_passes_cse

import (
	rt "github.com/nooga/let-go/pkg/rt"
	vm "github.com/nooga/let-go/pkg/vm"
)

func inst_key(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__17882_5 vm.Value
	var arg__17888_7 vm.Value
	var arg__17895_10 vm.Value
	var arg__17896_11 vm.Value
	var arg__17902_13 vm.Value
	var v14 vm.Value
	var callErr error
	_, _, _, _, _, _ = arg__17882_5, arg__17888_7, arg__17895_10, arg__17896_11, arg__17902_13, v14
	arg__17882_5, callErr = rt.InvokeValue(rt.LookupVar("ir", "op").Deref(), []vm.Value{arg0, arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__17888_7, callErr = rt.InvokeValue(rt.LookupVar("ir", "refs").Deref(), []vm.Value{arg0, arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__17895_10, callErr = rt.InvokeValue(rt.LookupVar("ir", "refs").Deref(), []vm.Value{arg0, arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__17896_11, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg__17895_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__17902_13, callErr = rt.InvokeValue(rt.LookupVar("ir", "aux").Deref(), []vm.Value{arg0, arg1})
	if callErr != nil {
		return nil, callErr
	}
	v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__17882_5, arg__17896_11, arg__17902_13})
	if callErr != nil {
		return nil, callErr
	}
	return v14, nil
}
func cse_eligible_QMARK_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value, arg3 vm.Value) (vm.Value, error) {
	var or__x_7 vm.Value
	var op_8 vm.Value
	var nid_9 vm.Value
	var f_10 vm.Value
	var var_facts_11 vm.Value
	var or__x_12 vm.Value
	var op_13 vm.Value
	var nid_14 vm.Value
	var f_15 vm.Value
	var var_facts_16 vm.Value
	var or__x_17 vm.Value
	var and__x_21 bool
	var v51 vm.Value
	var op_52 vm.Value
	var nid_53 vm.Value
	var f_54 vm.Value
	var var_facts_55 vm.Value
	var or__x_56 vm.Value
	var op_22 vm.Value
	var nid_23 vm.Value
	var f_24 vm.Value
	var var_facts_25 vm.Value
	var or__x_26 vm.Value
	var and__x_27 bool
	var arg__17916_36 vm.Value
	var arg__17924_39 vm.Value
	var v40 vm.Value
	var op_28 vm.Value
	var nid_29 vm.Value
	var f_30 vm.Value
	var var_facts_31 vm.Value
	var or__x_32 vm.Value
	var and__x_33 bool
	var v43 vm.Value
	var op_44 vm.Value
	var nid_45 vm.Value
	var f_46 vm.Value
	var var_facts_47 vm.Value
	var or__x_48 vm.Value
	var and__x_49 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = or__x_7, op_8, nid_9, f_10, var_facts_11, or__x_12, op_13, nid_14, f_15, var_facts_16, or__x_17, and__x_21, v51, op_52, nid_53, f_54, var_facts_55, or__x_56, op_22, nid_23, f_24, var_facts_25, or__x_26, and__x_27, arg__17916_36, arg__17924_39, v40, op_28, nid_29, f_30, var_facts_31, or__x_32, and__x_33, v43, op_44, nid_45, f_46, var_facts_47, or__x_48, and__x_49
	or__x_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "contains?").Deref(), []vm.Value{rt.LookupVar("ir.passes.cse", "pure-cse-ops").Deref(), arg0})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(or__x_7) {
		op_8 = arg0
		nid_9 = arg1
		f_10 = arg2
		var_facts_11 = arg3
		or__x_12 = or__x_7
		goto b1
	} else {
		op_13 = arg0
		nid_14 = arg1
		f_15 = arg2
		var_facts_16 = arg3
		or__x_17 = or__x_7
		goto b2
	}
b1:
	;
	v51 = or__x_12
	op_52 = op_8
	nid_53 = nid_9
	f_54 = f_10
	var_facts_55 = var_facts_11
	or__x_56 = or__x_12
	goto b3
b2:
	;
	and__x_21 = op_13 == vm.Keyword("load-var")
	if and__x_21 {
		op_22 = op_13
		nid_23 = nid_14
		f_24 = f_15
		var_facts_25 = var_facts_16
		or__x_26 = or__x_17
		and__x_27 = and__x_21
		goto b4
	} else {
		op_28 = op_13
		nid_29 = nid_14
		f_30 = f_15
		var_facts_31 = var_facts_16
		or__x_32 = or__x_17
		and__x_33 = and__x_21
		goto b5
	}
b3:
	;
	return v51, nil
b4:
	;
	arg__17916_36, callErr = rt.InvokeValue(rt.LookupVar("ir", "aux").Deref(), []vm.Value{nid_23, f_24})
	if callErr != nil {
		return nil, callErr
	}
	arg__17924_39, callErr = rt.InvokeValue(rt.LookupVar("ir", "aux").Deref(), []vm.Value{nid_23, f_24})
	if callErr != nil {
		return nil, callErr
	}
	v40, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.mutability", "stable-load-var?").Deref(), []vm.Value{var_facts_25, arg__17924_39})
	if callErr != nil {
		return nil, callErr
	}
	v43 = v40
	op_44 = op_22
	nid_45 = nid_23
	f_46 = f_24
	var_facts_47 = var_facts_25
	or__x_48 = or__x_26
	and__x_49 = vm.Boolean(and__x_27)
	goto b6
b5:
	;
	v43 = vm.Boolean(and__x_33)
	op_44 = op_28
	nid_45 = nid_29
	f_46 = f_30
	var_facts_47 = var_facts_31
	or__x_48 = or__x_32
	and__x_49 = vm.Boolean(and__x_33)
	goto b6
b6:
	;
	v51 = v43
	op_52 = op_44
	nid_53 = nid_45
	f_54 = f_46
	var_facts_55 = var_facts_47
	or__x_56 = or__x_48
	goto b3
}
func cse(arg0 vm.Value) (vm.Value, error) {
	var var_facts_3 vm.Value
	var arg__17935_5 vm.Value
	var arg__17940_8 vm.Value
	var doseq_seq__17925_9 vm.Value
	var doseq_loop__17926_10 vm.Value
	var var_facts_11 vm.Value
	var f_12 vm.Value
	var v177 vm.Value
	var doseq_seq__17925_14 vm.Value
	var doseq_loop__17926_15 vm.Value
	var var_facts_16 vm.Value
	var f_17 vm.Value
	var v176 vm.Value
	var bid_24 vm.Value
	var seen_28 vm.Value
	var arg__17952_30 vm.Value
	var arg__17959_33 vm.Value
	var arg__17960_34 vm.Value
	var arg__17967_37 vm.Value
	var arg__17974_40 vm.Value
	var arg__17975_41 vm.Value
	var doseq_seq__17927_42 vm.Value
	var doseq_seq__17925_18 vm.Value
	var doseq_loop__17926_19 vm.Value
	var var_facts_20 vm.Value
	var f_21 vm.Value
	var v185 vm.Value
	var v165 vm.Value
	var doseq_seq__17925_166 vm.Value
	var doseq_loop__17926_167 vm.Value
	var var_facts_168 vm.Value
	var f_169 vm.Value
	var doseq_loop__17928_43 vm.Value
	var var_facts_44 vm.Value
	var f_45 vm.Value
	var seen_46 vm.Value
	var v181 vm.Value
	var doseq_seq__17925_48 vm.Value
	var doseq_loop__17926_49 vm.Value
	var bid_50 vm.Value
	var doseq_seq__17927_51 vm.Value
	var doseq_loop__17928_52 vm.Value
	var var_facts_53 vm.Value
	var f_54 vm.Value
	var seen_55 vm.Value
	var v179 vm.Value
	var nid_66 vm.Value
	var op_68 vm.Value
	var v90 vm.Value
	var doseq_seq__17925_56 vm.Value
	var doseq_loop__17926_57 vm.Value
	var bid_58 vm.Value
	var doseq_seq__17927_59 vm.Value
	var doseq_loop__17928_60 vm.Value
	var var_facts_61 vm.Value
	var f_62 vm.Value
	var seen_63 vm.Value
	var v184 vm.Value
	var v151 vm.Value
	var doseq_seq__17925_152 vm.Value
	var doseq_loop__17926_153 vm.Value
	var bid_154 vm.Value
	var doseq_seq__17927_155 vm.Value
	var doseq_loop__17928_156 vm.Value
	var var_facts_157 vm.Value
	var f_158 vm.Value
	var seen_159 vm.Value
	var v175 vm.Value
	var v161 vm.Value
	var doseq_seq__17925_69 vm.Value
	var doseq_loop__17926_70 vm.Value
	var bid_71 vm.Value
	var doseq_seq__17927_72 vm.Value
	var doseq_loop__17928_73 vm.Value
	var var_facts_74 vm.Value
	var f_75 vm.Value
	var seen_76 vm.Value
	var nid_77 vm.Value
	var op_78 vm.Value
	var v182 vm.Value
	var k_93 vm.Value
	var head__18002_95 vm.Value
	var tem__G__0_96 vm.Value
	var doseq_seq__17925_79 vm.Value
	var doseq_loop__17926_80 vm.Value
	var bid_81 vm.Value
	var doseq_seq__17927_82 vm.Value
	var doseq_loop__17928_83 vm.Value
	var var_facts_84 vm.Value
	var f_85 vm.Value
	var seen_86 vm.Value
	var nid_87 vm.Value
	var op_88 vm.Value
	var v178 vm.Value
	var v147 vm.Value
	var doseq_seq__17925_97 vm.Value
	var doseq_loop__17926_98 vm.Value
	var bid_99 vm.Value
	var doseq_seq__17927_100 vm.Value
	var doseq_loop__17928_101 vm.Value
	var var_facts_102 vm.Value
	var f_103 vm.Value
	var seen_104 vm.Value
	var nid_105 vm.Value
	var op_106 vm.Value
	var k_107 vm.Value
	var tem__G__0_108 vm.Value
	var v180 vm.Value
	var v123 vm.Value
	var doseq_seq__17925_109 vm.Value
	var doseq_loop__17926_110 vm.Value
	var bid_111 vm.Value
	var doseq_seq__17927_112 vm.Value
	var doseq_loop__17928_113 vm.Value
	var var_facts_114 vm.Value
	var f_115 vm.Value
	var seen_116 vm.Value
	var nid_117 vm.Value
	var op_118 vm.Value
	var k_119 vm.Value
	var tem__G__0_120 vm.Value
	var v174 vm.Value
	var v128 vm.Value
	var v130 vm.Value
	var doseq_seq__17925_131 vm.Value
	var doseq_loop__17926_132 vm.Value
	var bid_133 vm.Value
	var doseq_seq__17927_134 vm.Value
	var doseq_loop__17928_135 vm.Value
	var var_facts_136 vm.Value
	var f_137 vm.Value
	var seen_138 vm.Value
	var nid_139 vm.Value
	var op_140 vm.Value
	var k_141 vm.Value
	var tem__G__0_142 vm.Value
	var v183 vm.Value
	var v144 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = var_facts_3, arg__17935_5, arg__17940_8, doseq_seq__17925_9, doseq_loop__17926_10, var_facts_11, f_12, v177, doseq_seq__17925_14, doseq_loop__17926_15, var_facts_16, f_17, v176, bid_24, seen_28, arg__17952_30, arg__17959_33, arg__17960_34, arg__17967_37, arg__17974_40, arg__17975_41, doseq_seq__17927_42, doseq_seq__17925_18, doseq_loop__17926_19, var_facts_20, f_21, v185, v165, doseq_seq__17925_166, doseq_loop__17926_167, var_facts_168, f_169, doseq_loop__17928_43, var_facts_44, f_45, seen_46, v181, doseq_seq__17925_48, doseq_loop__17926_49, bid_50, doseq_seq__17927_51, doseq_loop__17928_52, var_facts_53, f_54, seen_55, v179, nid_66, op_68, v90, doseq_seq__17925_56, doseq_loop__17926_57, bid_58, doseq_seq__17927_59, doseq_loop__17928_60, var_facts_61, f_62, seen_63, v184, v151, doseq_seq__17925_152, doseq_loop__17926_153, bid_154, doseq_seq__17927_155, doseq_loop__17928_156, var_facts_157, f_158, seen_159, v175, v161, doseq_seq__17925_69, doseq_loop__17926_70, bid_71, doseq_seq__17927_72, doseq_loop__17928_73, var_facts_74, f_75, seen_76, nid_77, op_78, v182, k_93, head__18002_95, tem__G__0_96, doseq_seq__17925_79, doseq_loop__17926_80, bid_81, doseq_seq__17927_82, doseq_loop__17928_83, var_facts_84, f_85, seen_86, nid_87, op_88, v178, v147, doseq_seq__17925_97, doseq_loop__17926_98, bid_99, doseq_seq__17927_100, doseq_loop__17928_101, var_facts_102, f_103, seen_104, nid_105, op_106, k_107, tem__G__0_108, v180, v123, doseq_seq__17925_109, doseq_loop__17926_110, bid_111, doseq_seq__17927_112, doseq_loop__17928_113, var_facts_114, f_115, seen_116, nid_117, op_118, k_119, tem__G__0_120, v174, v128, v130, doseq_seq__17925_131, doseq_loop__17926_132, bid_133, doseq_seq__17927_134, doseq_loop__17928_135, var_facts_136, f_137, seen_138, nid_139, op_140, k_141, tem__G__0_142, v183, v144
	var_facts_3, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.mutability", "analyze-var-stability").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__17935_5, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__17940_8, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	doseq_seq__17925_9, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{arg__17940_8})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__17926_10 = doseq_seq__17925_9
	var_facts_11 = var_facts_3
	f_12 = arg0
	v177 = vm.EmptyPersistentMap
	goto b1
b1:
	;
	if vm.IsTruthy(doseq_loop__17926_10) {
		doseq_seq__17925_14 = doseq_seq__17925_9
		doseq_loop__17926_15 = doseq_loop__17926_10
		var_facts_16 = var_facts_11
		f_17 = f_12
		v176 = v177
		goto b2
	} else {
		doseq_seq__17925_18 = doseq_seq__17925_9
		doseq_loop__17926_19 = doseq_loop__17926_10
		var_facts_20 = var_facts_11
		f_21 = f_12
		v185 = v177
		goto b3
	}
b2:
	;
	bid_24, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{doseq_loop__17926_15})
	if callErr != nil {
		return nil, callErr
	}
	seen_28, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "atom").Deref(), []vm.Value{v176})
	if callErr != nil {
		return nil, callErr
	}
	arg__17952_30, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-insts").Deref(), []vm.Value{bid_24, f_17})
	if callErr != nil {
		return nil, callErr
	}
	arg__17959_33, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-insts").Deref(), []vm.Value{bid_24, f_17})
	if callErr != nil {
		return nil, callErr
	}
	arg__17960_34, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg__17959_33})
	if callErr != nil {
		return nil, callErr
	}
	arg__17967_37, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-insts").Deref(), []vm.Value{bid_24, f_17})
	if callErr != nil {
		return nil, callErr
	}
	arg__17974_40, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-insts").Deref(), []vm.Value{bid_24, f_17})
	if callErr != nil {
		return nil, callErr
	}
	arg__17975_41, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg__17974_40})
	if callErr != nil {
		return nil, callErr
	}
	doseq_seq__17927_42, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{arg__17975_41})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__17928_43 = doseq_seq__17927_42
	var_facts_44 = var_facts_16
	f_45 = f_17
	seen_46 = seen_28
	v181 = v176
	goto b5
b3:
	;
	v165 = vm.NIL
	doseq_seq__17925_166 = doseq_seq__17925_18
	doseq_loop__17926_167 = doseq_loop__17926_19
	var_facts_168 = var_facts_20
	f_169 = f_21
	goto b4
b4:
	;
	return f_169, nil
b5:
	;
	if vm.IsTruthy(doseq_loop__17928_43) {
		doseq_seq__17925_48 = doseq_seq__17925_14
		doseq_loop__17926_49 = doseq_loop__17926_15
		bid_50 = bid_24
		doseq_seq__17927_51 = doseq_seq__17927_42
		doseq_loop__17928_52 = doseq_loop__17928_43
		var_facts_53 = var_facts_44
		f_54 = f_45
		seen_55 = seen_46
		v179 = v181
		goto b6
	} else {
		doseq_seq__17925_56 = doseq_seq__17925_14
		doseq_loop__17926_57 = doseq_loop__17926_15
		bid_58 = bid_24
		doseq_seq__17927_59 = doseq_seq__17927_42
		doseq_loop__17928_60 = doseq_loop__17928_43
		var_facts_61 = var_facts_44
		f_62 = f_45
		seen_63 = seen_46
		v184 = v181
		goto b7
	}
b6:
	;
	nid_66, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{doseq_loop__17928_52})
	if callErr != nil {
		return nil, callErr
	}
	op_68, callErr = rt.InvokeValue(rt.LookupVar("ir", "op").Deref(), []vm.Value{nid_66, f_54})
	if callErr != nil {
		return nil, callErr
	}
	v90, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.cse", "cse-eligible?").Deref(), []vm.Value{op_68, nid_66, f_54, var_facts_53})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v90) {
		doseq_seq__17925_69 = doseq_seq__17925_48
		doseq_loop__17926_70 = doseq_loop__17926_49
		bid_71 = bid_50
		doseq_seq__17927_72 = doseq_seq__17927_51
		doseq_loop__17928_73 = doseq_loop__17928_52
		var_facts_74 = var_facts_53
		f_75 = f_54
		seen_76 = seen_55
		nid_77 = nid_66
		op_78 = op_68
		v182 = v179
		goto b9
	} else {
		doseq_seq__17925_79 = doseq_seq__17925_48
		doseq_loop__17926_80 = doseq_loop__17926_49
		bid_81 = bid_50
		doseq_seq__17927_82 = doseq_seq__17927_51
		doseq_loop__17928_83 = doseq_loop__17928_52
		var_facts_84 = var_facts_53
		f_85 = f_54
		seen_86 = seen_55
		nid_87 = nid_66
		op_88 = op_68
		v178 = v179
		goto b10
	}
b7:
	;
	v151 = vm.NIL
	doseq_seq__17925_152 = doseq_seq__17925_56
	doseq_loop__17926_153 = doseq_loop__17926_57
	bid_154 = bid_58
	doseq_seq__17927_155 = doseq_seq__17927_59
	doseq_loop__17928_156 = doseq_loop__17928_60
	var_facts_157 = var_facts_61
	f_158 = f_62
	seen_159 = seen_63
	v175 = v184
	goto b8
b8:
	;
	v161, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "next").Deref(), []vm.Value{doseq_loop__17926_153})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__17926_10 = v161
	var_facts_11 = var_facts_157
	f_12 = f_158
	v177 = v175
	goto b1
b9:
	;
	k_93, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.cse", "inst-key").Deref(), []vm.Value{nid_77, f_75})
	if callErr != nil {
		return nil, callErr
	}
	head__18002_95, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{seen_76})
	if callErr != nil {
		return nil, callErr
	}
	tem__G__0_96, callErr = rt.InvokeValue(head__18002_95, []vm.Value{k_93})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(tem__G__0_96) {
		doseq_seq__17925_97 = doseq_seq__17925_69
		doseq_loop__17926_98 = doseq_loop__17926_70
		bid_99 = bid_71
		doseq_seq__17927_100 = doseq_seq__17927_72
		doseq_loop__17928_101 = doseq_loop__17928_73
		var_facts_102 = var_facts_74
		f_103 = f_75
		seen_104 = seen_76
		nid_105 = nid_77
		op_106 = op_78
		k_107 = k_93
		tem__G__0_108 = tem__G__0_96
		v180 = v182
		goto b12
	} else {
		doseq_seq__17925_109 = doseq_seq__17925_69
		doseq_loop__17926_110 = doseq_loop__17926_70
		bid_111 = bid_71
		doseq_seq__17927_112 = doseq_seq__17927_72
		doseq_loop__17928_113 = doseq_loop__17928_73
		var_facts_114 = var_facts_74
		f_115 = f_75
		seen_116 = seen_76
		nid_117 = nid_77
		op_118 = op_78
		k_119 = k_93
		tem__G__0_120 = tem__G__0_96
		v174 = v182
		goto b13
	}
b10:
	;
	v147, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "next").Deref(), []vm.Value{doseq_loop__17928_83})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__17928_43 = v147
	var_facts_44 = var_facts_84
	f_45 = f_85
	seen_46 = seen_86
	v181 = v178
	goto b5
b12:
	;
	v123, callErr = rt.InvokeValue(rt.LookupVar("ir", "replace-all-uses!").Deref(), []vm.Value{f_103, nid_105, tem__G__0_108})
	if callErr != nil {
		return nil, callErr
	}
	v130 = v123
	doseq_seq__17925_131 = doseq_seq__17925_97
	doseq_loop__17926_132 = doseq_loop__17926_98
	bid_133 = bid_99
	doseq_seq__17927_134 = doseq_seq__17927_100
	doseq_loop__17928_135 = doseq_loop__17928_101
	var_facts_136 = var_facts_102
	f_137 = f_103
	seen_138 = seen_104
	nid_139 = nid_105
	op_140 = op_106
	k_141 = k_107
	tem__G__0_142 = tem__G__0_108
	v183 = v180
	goto b14
b13:
	;
	v128, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{seen_116, rt.LookupVar("clojure.core", "assoc").Deref(), k_119, nid_117})
	if callErr != nil {
		return nil, callErr
	}
	v130 = v128
	doseq_seq__17925_131 = doseq_seq__17925_109
	doseq_loop__17926_132 = doseq_loop__17926_110
	bid_133 = bid_111
	doseq_seq__17927_134 = doseq_seq__17927_112
	doseq_loop__17928_135 = doseq_loop__17928_113
	var_facts_136 = var_facts_114
	f_137 = f_115
	seen_138 = seen_116
	nid_139 = nid_117
	op_140 = op_118
	k_141 = k_119
	tem__G__0_142 = tem__G__0_120
	v183 = v174
	goto b14
b14:
	;
	v144, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "next").Deref(), []vm.Value{doseq_loop__17928_135})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__17928_43 = v144
	var_facts_44 = var_facts_136
	f_45 = f_137
	seen_46 = seen_138
	v181 = v183
	goto b5
}
