package ir_dominance

import (
	rt "github.com/nooga/let-go/pkg/rt"
	vm "github.com/nooga/let-go/pkg/vm"
)

func intersect(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value, arg3 vm.Value) (vm.Value, error) {
	var b1_5 vm.Value
	var b2_6 vm.Value
	var b2_7 vm.Value
	var b1_8 vm.Value
	var rpo_idx_9 vm.Value
	var idom_10 vm.Value
	var v20 bool
	var b2_12 vm.Value
	var b1_13 vm.Value
	var rpo_idx_14 vm.Value
	var idom_15 vm.Value
	var b2_16 vm.Value
	var b1_17 vm.Value
	var rpo_idx_18 vm.Value
	var idom_19 vm.Value
	var arg__8453_32 vm.Value
	var arg__8459_34 vm.Value
	var v35 bool
	var v88 vm.Value
	var b2_89 vm.Value
	var b1_90 vm.Value
	var rpo_idx_91 vm.Value
	var idom_92 vm.Value
	var b2_23 vm.Value
	var b1_24 vm.Value
	var rpo_idx_25 vm.Value
	var idom_26 vm.Value
	var v38 vm.Value
	var b2_27 vm.Value
	var b1_28 vm.Value
	var rpo_idx_29 vm.Value
	var idom_30 vm.Value
	var arg__8470_49 vm.Value
	var arg__8476_51 vm.Value
	var v52 bool
	var v82 vm.Value
	var b2_83 vm.Value
	var b1_84 vm.Value
	var rpo_idx_85 vm.Value
	var idom_86 vm.Value
	var b2_40 vm.Value
	var b1_41 vm.Value
	var rpo_idx_42 vm.Value
	var idom_43 vm.Value
	var v55 vm.Value
	var b2_44 vm.Value
	var b1_45 vm.Value
	var rpo_idx_46 vm.Value
	var idom_47 vm.Value
	var v76 vm.Value
	var b2_77 vm.Value
	var b1_78 vm.Value
	var rpo_idx_79 vm.Value
	var idom_80 vm.Value
	var b2_57 vm.Value
	var b1_58 vm.Value
	var rpo_idx_59 vm.Value
	var idom_60 vm.Value
	var b2_61 vm.Value
	var b1_62 vm.Value
	var rpo_idx_63 vm.Value
	var idom_64 vm.Value
	var v70 vm.Value
	var b2_71 vm.Value
	var b1_72 vm.Value
	var rpo_idx_73 vm.Value
	var idom_74 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = b1_5, b2_6, b2_7, b1_8, rpo_idx_9, idom_10, v20, b2_12, b1_13, rpo_idx_14, idom_15, b2_16, b1_17, rpo_idx_18, idom_19, arg__8453_32, arg__8459_34, v35, v88, b2_89, b1_90, rpo_idx_91, idom_92, b2_23, b1_24, rpo_idx_25, idom_26, v38, b2_27, b1_28, rpo_idx_29, idom_30, arg__8470_49, arg__8476_51, v52, v82, b2_83, b1_84, rpo_idx_85, idom_86, b2_40, b1_41, rpo_idx_42, idom_43, v55, b2_44, b1_45, rpo_idx_46, idom_47, v76, b2_77, b1_78, rpo_idx_79, idom_80, b2_57, b1_58, rpo_idx_59, idom_60, b2_61, b1_62, rpo_idx_63, idom_64, v70, b2_71, b1_72, rpo_idx_73, idom_74
	b1_5 = arg2
	b2_6 = arg3
	b2_7 = arg3
	b1_8 = arg2
	rpo_idx_9 = arg1
	idom_10 = arg0
	goto b1
b1:
	;
	v20 = b1_8 == b2_7
	if v20 {
		b2_12 = b2_7
		b1_13 = b1_8
		rpo_idx_14 = rpo_idx_9
		idom_15 = idom_10
		goto b2
	} else {
		b2_16 = b2_7
		b1_17 = b1_8
		rpo_idx_18 = rpo_idx_9
		idom_19 = idom_10
		goto b3
	}
b2:
	;
	v88 = b1_13
	b2_89 = b2_12
	b1_90 = b1_13
	rpo_idx_91 = rpo_idx_14
	idom_92 = idom_15
	goto b4
b3:
	;
	arg__8453_32, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{rpo_idx_18, b1_17})
	if callErr != nil {
		return nil, callErr
	}
	arg__8459_34, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{rpo_idx_18, b2_16})
	if callErr != nil {
		return nil, callErr
	}
	v35 = rt.GtValue(arg__8453_32, arg__8459_34)
	if v35 {
		b2_23 = b2_16
		b1_24 = b1_17
		rpo_idx_25 = rpo_idx_18
		idom_26 = idom_19
		goto b5
	} else {
		b2_27 = b2_16
		b1_28 = b1_17
		rpo_idx_29 = rpo_idx_18
		idom_30 = idom_19
		goto b6
	}
b4:
	;
	return v88, nil
b5:
	;
	v38, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{idom_26, b1_24})
	if callErr != nil {
		return nil, callErr
	}
	b1_5 = v38
	b2_6 = b2_23
	b2_7 = b2_23
	b1_8 = b1_24
	rpo_idx_9 = rpo_idx_25
	idom_10 = idom_26
	goto b1
b6:
	;
	arg__8470_49, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{rpo_idx_29, b2_27})
	if callErr != nil {
		return nil, callErr
	}
	arg__8476_51, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{rpo_idx_29, b1_28})
	if callErr != nil {
		return nil, callErr
	}
	v52 = rt.GtValue(arg__8470_49, arg__8476_51)
	if v52 {
		b2_40 = b2_27
		b1_41 = b1_28
		rpo_idx_42 = rpo_idx_29
		idom_43 = idom_30
		goto b8
	} else {
		b2_44 = b2_27
		b1_45 = b1_28
		rpo_idx_46 = rpo_idx_29
		idom_47 = idom_30
		goto b9
	}
b7:
	;
	v88 = v82
	b2_89 = b2_83
	b1_90 = b1_84
	rpo_idx_91 = rpo_idx_85
	idom_92 = idom_86
	goto b4
b8:
	;
	v55, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{idom_43, b2_40})
	if callErr != nil {
		return nil, callErr
	}
	b1_5 = b1_41
	b2_6 = v55
	b2_7 = b2_40
	b1_8 = b1_41
	rpo_idx_9 = rpo_idx_42
	idom_10 = idom_43
	goto b1
b9:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		b2_57 = b2_44
		b1_58 = b1_45
		rpo_idx_59 = rpo_idx_46
		idom_60 = idom_47
		goto b11
	} else {
		b2_61 = b2_44
		b1_62 = b1_45
		rpo_idx_63 = rpo_idx_46
		idom_64 = idom_47
		goto b12
	}
b10:
	;
	v82 = v76
	b2_83 = b2_77
	b1_84 = b1_78
	rpo_idx_85 = rpo_idx_79
	idom_86 = idom_80
	goto b7
b11:
	;
	v70 = b1_58
	b2_71 = b2_57
	b1_72 = b1_58
	rpo_idx_73 = rpo_idx_59
	idom_74 = idom_60
	goto b13
b12:
	;
	v70 = vm.NIL
	b2_71 = b2_61
	b1_72 = b1_62
	rpo_idx_73 = rpo_idx_63
	idom_74 = idom_64
	goto b13
b13:
	;
	v76 = v70
	b2_77 = b2_71
	b1_78 = b1_72
	rpo_idx_79 = rpo_idx_73
	idom_80 = idom_74
	goto b10
}
func successors(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var term_4 vm.Value
	var op_6 vm.Value
	var aux_8 vm.Value
	var v20 bool
	var f_9 vm.Value
	var bid_10 vm.Value
	var term_11 vm.Value
	var op_12 vm.Value
	var aux_13 vm.Value
	var arg__8503_24 vm.Value
	var v25 vm.Value
	var f_14 vm.Value
	var bid_15 vm.Value
	var term_16 vm.Value
	var op_17 vm.Value
	var aux_18 vm.Value
	var v38 bool
	var v81 vm.Value
	var f_82 vm.Value
	var bid_83 vm.Value
	var term_84 vm.Value
	var op_85 vm.Value
	var aux_86 vm.Value
	var f_27 vm.Value
	var bid_28 vm.Value
	var term_29 vm.Value
	var op_30 vm.Value
	var aux_31 vm.Value
	var t_41 vm.Value
	var e_43 vm.Value
	var arg__8516_46 vm.Value
	var arg__8520_48 vm.Value
	var v49 vm.Value
	var f_32 vm.Value
	var bid_33 vm.Value
	var term_34 vm.Value
	var op_35 vm.Value
	var aux_36 vm.Value
	var v74 vm.Value
	var f_75 vm.Value
	var bid_76 vm.Value
	var term_77 vm.Value
	var op_78 vm.Value
	var aux_79 vm.Value
	var f_51 vm.Value
	var bid_52 vm.Value
	var term_53 vm.Value
	var op_54 vm.Value
	var aux_55 vm.Value
	var f_56 vm.Value
	var bid_57 vm.Value
	var term_58 vm.Value
	var op_59 vm.Value
	var aux_60 vm.Value
	var v67 vm.Value
	var f_68 vm.Value
	var bid_69 vm.Value
	var term_70 vm.Value
	var op_71 vm.Value
	var aux_72 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = term_4, op_6, aux_8, v20, f_9, bid_10, term_11, op_12, aux_13, arg__8503_24, v25, f_14, bid_15, term_16, op_17, aux_18, v38, v81, f_82, bid_83, term_84, op_85, aux_86, f_27, bid_28, term_29, op_30, aux_31, t_41, e_43, arg__8516_46, arg__8520_48, v49, f_32, bid_33, term_34, op_35, aux_36, v74, f_75, bid_76, term_77, op_78, aux_79, f_51, bid_52, term_53, op_54, aux_55, f_56, bid_57, term_58, op_59, aux_60, v67, f_68, bid_69, term_70, op_71, aux_72
	term_4, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-term").Deref(), []vm.Value{arg1, arg0})
	if callErr != nil {
		return nil, callErr
	}
	op_6, callErr = rt.InvokeValue(rt.LookupVar("ir", "op").Deref(), []vm.Value{term_4, arg0})
	if callErr != nil {
		return nil, callErr
	}
	aux_8, callErr = rt.InvokeValue(rt.LookupVar("ir", "aux").Deref(), []vm.Value{term_4, arg0})
	if callErr != nil {
		return nil, callErr
	}
	v20 = op_6 == vm.Keyword("branch")
	if v20 {
		f_9 = arg0
		bid_10 = arg1
		term_11 = term_4
		op_12 = op_6
		aux_13 = aux_8
		goto b1
	} else {
		f_14 = arg0
		bid_15 = arg1
		term_16 = term_4
		op_17 = op_6
		aux_18 = aux_8
		goto b2
	}
b1:
	;
	arg__8503_24, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-target").Deref(), []vm.Value{aux_13})
	if callErr != nil {
		return nil, callErr
	}
	v25, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__8503_24})
	if callErr != nil {
		return nil, callErr
	}
	v81 = v25
	f_82 = f_9
	bid_83 = bid_10
	term_84 = term_11
	op_85 = op_12
	aux_86 = aux_13
	goto b3
b2:
	;
	v38 = op_17 == vm.Keyword("branch-if")
	if v38 {
		f_27 = f_14
		bid_28 = bid_15
		term_29 = term_16
		op_30 = op_17
		aux_31 = aux_18
		goto b4
	} else {
		f_32 = f_14
		bid_33 = bid_15
		term_34 = term_16
		op_35 = op_17
		aux_36 = aux_18
		goto b5
	}
b3:
	;
	return v81, nil
b4:
	;
	t_41, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-true").Deref(), []vm.Value{aux_31})
	if callErr != nil {
		return nil, callErr
	}
	e_43, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-false").Deref(), []vm.Value{aux_31})
	if callErr != nil {
		return nil, callErr
	}
	arg__8516_46, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-target").Deref(), []vm.Value{t_41})
	if callErr != nil {
		return nil, callErr
	}
	arg__8520_48, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-target").Deref(), []vm.Value{e_43})
	if callErr != nil {
		return nil, callErr
	}
	v49, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__8516_46, arg__8520_48})
	if callErr != nil {
		return nil, callErr
	}
	v74 = v49
	f_75 = f_27
	bid_76 = bid_28
	term_77 = term_29
	op_78 = op_30
	aux_79 = aux_31
	goto b6
b5:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		f_51 = f_32
		bid_52 = bid_33
		term_53 = term_34
		op_54 = op_35
		aux_55 = aux_36
		goto b7
	} else {
		f_56 = f_32
		bid_57 = bid_33
		term_58 = term_34
		op_59 = op_35
		aux_60 = aux_36
		goto b8
	}
b6:
	;
	v81 = v74
	f_82 = f_75
	bid_83 = bid_76
	term_84 = term_77
	op_85 = op_78
	aux_86 = aux_79
	goto b3
b7:
	;
	v67 = vm.NewArrayVector([]vm.Value{})
	f_68 = f_51
	bid_69 = bid_52
	term_70 = term_53
	op_71 = op_54
	aux_72 = aux_55
	goto b9
b8:
	;
	v67 = vm.NIL
	f_68 = f_56
	bid_69 = bid_57
	term_70 = term_58
	op_71 = op_59
	aux_72 = aux_60
	goto b9
b9:
	;
	v74 = v67
	f_75 = f_68
	bid_76 = bid_69
	term_77 = term_70
	op_78 = op_71
	aux_79 = aux_72
	goto b6
}
func dfs_postorder(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__8529_10 vm.Value
	var arg__8531_12 vm.Value
	var v13 vm.Value
	var v15 vm.Value
	var stack_3 vm.Value
	var visited_4 vm.Value
	var post_5 vm.Value
	var f_6 vm.Value
	var v29 vm.Value
	var entry_18 vm.Value
	var stack_19 vm.Value
	var visited_20 vm.Value
	var post_21 vm.Value
	var f_22 vm.Value
	var entry_23 vm.Value
	var stack_24 vm.Value
	var visited_25 vm.Value
	var post_26 vm.Value
	var f_27 vm.Value
	var top_33 vm.Value
	var bid_37 vm.Value
	var succs_41 vm.Value
	var si_45 vm.Value
	var arg__8560_65 vm.Value
	var v66 bool
	var v130 vm.Value
	var entry_131 vm.Value
	var stack_132 vm.Value
	var visited_133 vm.Value
	var post_134 vm.Value
	var f_135 vm.Value
	var entry_46 vm.Value
	var stack_47 vm.Value
	var visited_48 vm.Value
	var post_49 vm.Value
	var f_50 vm.Value
	var top_51 vm.Value
	var bid_52 vm.Value
	var succs_53 vm.Value
	var si_54 vm.Value
	var s_69 vm.Value
	var arg__8570_71 vm.Value
	var arg__8576_74 vm.Value
	var arg__8577_75 vm.Value
	var arg__8583_78 vm.Value
	var arg__8584_79 vm.Value
	var arg__8590_82 vm.Value
	var stack_PRIME__83 vm.Value
	var v107 vm.Value
	var entry_55 vm.Value
	var stack_56 vm.Value
	var visited_57 vm.Value
	var post_58 vm.Value
	var f_59 vm.Value
	var top_60 vm.Value
	var bid_61 vm.Value
	var succs_62 vm.Value
	var si_63 vm.Value
	var v126 vm.Value
	var v128 vm.Value
	var entry_84 vm.Value
	var stack_85 vm.Value
	var visited_86 vm.Value
	var post_87 vm.Value
	var f_88 vm.Value
	var top_89 vm.Value
	var bid_90 vm.Value
	var succs_91 vm.Value
	var si_92 vm.Value
	var s_93 vm.Value
	var stack_PRIME__94 vm.Value
	var entry_95 vm.Value
	var stack_96 vm.Value
	var visited_97 vm.Value
	var post_98 vm.Value
	var f_99 vm.Value
	var top_100 vm.Value
	var bid_101 vm.Value
	var succs_102 vm.Value
	var si_103 vm.Value
	var s_104 vm.Value
	var stack_PRIME__105 vm.Value
	var arg__8604_112 vm.Value
	var arg__8606_114 vm.Value
	var arg__8616_118 vm.Value
	var arg__8618_120 vm.Value
	var v121 vm.Value
	var v123 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__8529_10, arg__8531_12, v13, v15, stack_3, visited_4, post_5, f_6, v29, entry_18, stack_19, visited_20, post_21, f_22, entry_23, stack_24, visited_25, post_26, f_27, top_33, bid_37, succs_41, si_45, arg__8560_65, v66, v130, entry_131, stack_132, visited_133, post_134, f_135, entry_46, stack_47, visited_48, post_49, f_50, top_51, bid_52, succs_53, si_54, s_69, arg__8570_71, arg__8576_74, arg__8577_75, arg__8583_78, arg__8584_79, arg__8590_82, stack_PRIME__83, v107, entry_55, stack_56, visited_57, post_58, f_59, top_60, bid_61, succs_62, si_63, v126, v128, entry_84, stack_85, visited_86, post_87, f_88, top_89, bid_90, succs_91, si_92, s_93, stack_PRIME__94, entry_95, stack_96, visited_97, post_98, f_99, top_100, bid_101, succs_102, si_103, s_104, stack_PRIME__105, arg__8604_112, arg__8606_114, arg__8616_118, arg__8618_120, v121, v123
	arg__8529_10, callErr = rt.InvokeValue(rt.LookupVar("ir.dominance", "successors").Deref(), []vm.Value{arg0, arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__8531_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg1, arg__8529_10, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	v13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__8531_12})
	if callErr != nil {
		return nil, callErr
	}
	v15, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	stack_3 = v13
	visited_4 = v15
	post_5 = vm.NewArrayVector([]vm.Value{})
	f_6 = arg0
	goto b1
b1:
	;
	v29, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "empty?").Deref(), []vm.Value{stack_3})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v29) {
		entry_18 = arg1
		stack_19 = stack_3
		visited_20 = visited_4
		post_21 = post_5
		f_22 = f_6
		goto b2
	} else {
		entry_23 = arg1
		stack_24 = stack_3
		visited_25 = visited_4
		post_26 = post_5
		f_27 = f_6
		goto b3
	}
b2:
	;
	v130 = post_21
	entry_131 = entry_18
	stack_132 = stack_19
	visited_133 = visited_20
	post_134 = post_21
	f_135 = f_22
	goto b4
b3:
	;
	top_33, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "peek").Deref(), []vm.Value{stack_24})
	if callErr != nil {
		return nil, callErr
	}
	bid_37, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{top_33, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	succs_41, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{top_33, vm.Int(1)})
	if callErr != nil {
		return nil, callErr
	}
	si_45, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{top_33, vm.Int(2)})
	if callErr != nil {
		return nil, callErr
	}
	arg__8560_65, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{succs_41})
	if callErr != nil {
		return nil, callErr
	}
	v66 = rt.LtValue(si_45, arg__8560_65)
	if v66 {
		entry_46 = entry_23
		stack_47 = stack_24
		visited_48 = visited_25
		post_49 = post_26
		f_50 = f_27
		top_51 = top_33
		bid_52 = bid_37
		succs_53 = succs_41
		si_54 = si_45
		goto b5
	} else {
		entry_55 = entry_23
		stack_56 = stack_24
		visited_57 = visited_25
		post_58 = post_26
		f_59 = f_27
		top_60 = top_33
		bid_61 = bid_37
		succs_62 = succs_41
		si_63 = si_45
		goto b6
	}
b4:
	;
	return v130, nil
b5:
	;
	s_69, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{succs_53, si_54})
	if callErr != nil {
		return nil, callErr
	}
	arg__8570_71, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{stack_47})
	if callErr != nil {
		return nil, callErr
	}
	arg__8576_74 = rt.AddValue(si_54, vm.Int(1))
	arg__8577_75, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{bid_52, succs_53, arg__8576_74})
	if callErr != nil {
		return nil, callErr
	}
	arg__8583_78, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{stack_47})
	if callErr != nil {
		return nil, callErr
	}
	arg__8584_79 = rt.SubValue(arg__8583_78, vm.Int(1))
	arg__8590_82, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{bid_52, succs_53, arg__8576_74})
	if callErr != nil {
		return nil, callErr
	}
	stack_PRIME__83, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "assoc").Deref(), []vm.Value{stack_47, arg__8584_79, arg__8590_82})
	if callErr != nil {
		return nil, callErr
	}
	v107, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "contains?").Deref(), []vm.Value{visited_48, s_69})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v107) {
		entry_84 = entry_46
		stack_85 = stack_47
		visited_86 = visited_48
		post_87 = post_49
		f_88 = f_50
		top_89 = top_51
		bid_90 = bid_52
		succs_91 = succs_53
		si_92 = si_54
		s_93 = s_69
		stack_PRIME__94 = stack_PRIME__83
		goto b8
	} else {
		entry_95 = entry_46
		stack_96 = stack_47
		visited_97 = visited_48
		post_98 = post_49
		f_99 = f_50
		top_100 = top_51
		bid_101 = bid_52
		succs_102 = succs_53
		si_103 = si_54
		s_104 = s_69
		stack_PRIME__105 = stack_PRIME__83
		goto b9
	}
b6:
	;
	v126, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "pop").Deref(), []vm.Value{stack_56})
	if callErr != nil {
		return nil, callErr
	}
	v128, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{post_58, bid_61})
	if callErr != nil {
		return nil, callErr
	}
	stack_3 = v126
	visited_4 = visited_57
	post_5 = v128
	f_6 = f_59
	goto b1
b8:
	;
	stack_3 = stack_PRIME__94
	visited_4 = visited_86
	post_5 = post_87
	f_6 = f_88
	goto b1
b9:
	;
	arg__8604_112, callErr = rt.InvokeValue(rt.LookupVar("ir.dominance", "successors").Deref(), []vm.Value{f_99, s_104})
	if callErr != nil {
		return nil, callErr
	}
	arg__8606_114, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{s_104, arg__8604_112, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	arg__8616_118, callErr = rt.InvokeValue(rt.LookupVar("ir.dominance", "successors").Deref(), []vm.Value{f_99, s_104})
	if callErr != nil {
		return nil, callErr
	}
	arg__8618_120, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{s_104, arg__8616_118, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	v121, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{stack_PRIME__105, arg__8618_120})
	if callErr != nil {
		return nil, callErr
	}
	v123, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{visited_97, s_104})
	if callErr != nil {
		return nil, callErr
	}
	stack_3 = v121
	visited_4 = v123
	post_5 = post_98
	f_6 = f_99
	goto b1
}
func reverse_postorder(arg0 vm.Value) (vm.Value, error) {
	var arg__8636_3 vm.Value
	var arg__8642_6 vm.Value
	var arg__8643_7 vm.Value
	var arg__8649_10 vm.Value
	var arg__8655_13 vm.Value
	var arg__8656_14 vm.Value
	var arg__8657_15 vm.Value
	var arg__8663_18 vm.Value
	var arg__8669_21 vm.Value
	var arg__8670_22 vm.Value
	var arg__8676_25 vm.Value
	var arg__8682_28 vm.Value
	var arg__8683_29 vm.Value
	var arg__8684_30 vm.Value
	var v31 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__8636_3, arg__8642_6, arg__8643_7, arg__8649_10, arg__8655_13, arg__8656_14, arg__8657_15, arg__8663_18, arg__8669_21, arg__8670_22, arg__8676_25, arg__8682_28, arg__8683_29, arg__8684_30, v31
	arg__8636_3, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-entry").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__8642_6, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-entry").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__8643_7, callErr = rt.InvokeValue(rt.LookupVar("ir.dominance", "dfs-postorder").Deref(), []vm.Value{arg0, arg__8642_6})
	if callErr != nil {
		return nil, callErr
	}
	arg__8649_10, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-entry").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__8655_13, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-entry").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__8656_14, callErr = rt.InvokeValue(rt.LookupVar("ir.dominance", "dfs-postorder").Deref(), []vm.Value{arg0, arg__8655_13})
	if callErr != nil {
		return nil, callErr
	}
	arg__8657_15, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "reverse").Deref(), []vm.Value{arg__8656_14})
	if callErr != nil {
		return nil, callErr
	}
	arg__8663_18, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-entry").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__8669_21, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-entry").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__8670_22, callErr = rt.InvokeValue(rt.LookupVar("ir.dominance", "dfs-postorder").Deref(), []vm.Value{arg0, arg__8669_21})
	if callErr != nil {
		return nil, callErr
	}
	arg__8676_25, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-entry").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__8682_28, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-entry").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__8683_29, callErr = rt.InvokeValue(rt.LookupVar("ir.dominance", "dfs-postorder").Deref(), []vm.Value{arg0, arg__8682_28})
	if callErr != nil {
		return nil, callErr
	}
	arg__8684_30, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "reverse").Deref(), []vm.Value{arg__8683_29})
	if callErr != nil {
		return nil, callErr
	}
	v31, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg__8684_30})
	if callErr != nil {
		return nil, callErr
	}
	return v31, nil
}
func refine_idom(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value, arg3 vm.Value) (vm.Value, error) {
	var preds_6 vm.Value
	var ps_7 vm.Value
	var new_id_8 vm.Value
	var rpo_idx_9 vm.Value
	var idom_10 vm.Value
	var v28 vm.Value
	var f_13 vm.Value
	var bid_14 vm.Value
	var preds_15 vm.Value
	var ps_16 vm.Value
	var new_id_17 vm.Value
	var rpo_idx_18 vm.Value
	var idom_19 vm.Value
	var f_20 vm.Value
	var bid_21 vm.Value
	var preds_22 vm.Value
	var ps_23 vm.Value
	var new_id_24 vm.Value
	var rpo_idx_25 vm.Value
	var idom_26 vm.Value
	var p_32 vm.Value
	var arg__8702_51 vm.Value
	var v52 bool
	var v134 vm.Value
	var f_135 vm.Value
	var bid_136 vm.Value
	var preds_137 vm.Value
	var ps_138 vm.Value
	var new_id_139 vm.Value
	var rpo_idx_140 vm.Value
	var idom_141 vm.Value
	var f_33 vm.Value
	var bid_34 vm.Value
	var preds_35 vm.Value
	var ps_36 vm.Value
	var new_id_37 vm.Value
	var rpo_idx_38 vm.Value
	var idom_39 vm.Value
	var p_40 vm.Value
	var v55 vm.Value
	var f_41 vm.Value
	var bid_42 vm.Value
	var preds_43 vm.Value
	var ps_44 vm.Value
	var new_id_45 vm.Value
	var rpo_idx_46 vm.Value
	var idom_47 vm.Value
	var p_48 vm.Value
	var v74 bool
	var v124 vm.Value
	var f_125 vm.Value
	var bid_126 vm.Value
	var preds_127 vm.Value
	var ps_128 vm.Value
	var new_id_129 vm.Value
	var rpo_idx_130 vm.Value
	var idom_131 vm.Value
	var p_132 vm.Value
	var f_57 vm.Value
	var bid_58 vm.Value
	var preds_59 vm.Value
	var ps_60 vm.Value
	var new_id_61 vm.Value
	var rpo_idx_62 vm.Value
	var idom_63 vm.Value
	var p_64 vm.Value
	var v77 vm.Value
	var f_65 vm.Value
	var bid_66 vm.Value
	var preds_67 vm.Value
	var ps_68 vm.Value
	var new_id_69 vm.Value
	var rpo_idx_70 vm.Value
	var idom_71 vm.Value
	var p_72 vm.Value
	var v114 vm.Value
	var f_115 vm.Value
	var bid_116 vm.Value
	var preds_117 vm.Value
	var ps_118 vm.Value
	var new_id_119 vm.Value
	var rpo_idx_120 vm.Value
	var idom_121 vm.Value
	var p_122 vm.Value
	var f_79 vm.Value
	var bid_80 vm.Value
	var preds_81 vm.Value
	var ps_82 vm.Value
	var new_id_83 vm.Value
	var rpo_idx_84 vm.Value
	var idom_85 vm.Value
	var p_86 vm.Value
	var v98 vm.Value
	var v100 vm.Value
	var f_87 vm.Value
	var bid_88 vm.Value
	var preds_89 vm.Value
	var ps_90 vm.Value
	var new_id_91 vm.Value
	var rpo_idx_92 vm.Value
	var idom_93 vm.Value
	var p_94 vm.Value
	var v104 vm.Value
	var f_105 vm.Value
	var bid_106 vm.Value
	var preds_107 vm.Value
	var ps_108 vm.Value
	var new_id_109 vm.Value
	var rpo_idx_110 vm.Value
	var idom_111 vm.Value
	var p_112 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = preds_6, ps_7, new_id_8, rpo_idx_9, idom_10, v28, f_13, bid_14, preds_15, ps_16, new_id_17, rpo_idx_18, idom_19, f_20, bid_21, preds_22, ps_23, new_id_24, rpo_idx_25, idom_26, p_32, arg__8702_51, v52, v134, f_135, bid_136, preds_137, ps_138, new_id_139, rpo_idx_140, idom_141, f_33, bid_34, preds_35, ps_36, new_id_37, rpo_idx_38, idom_39, p_40, v55, f_41, bid_42, preds_43, ps_44, new_id_45, rpo_idx_46, idom_47, p_48, v74, v124, f_125, bid_126, preds_127, ps_128, new_id_129, rpo_idx_130, idom_131, p_132, f_57, bid_58, preds_59, ps_60, new_id_61, rpo_idx_62, idom_63, p_64, v77, f_65, bid_66, preds_67, ps_68, new_id_69, rpo_idx_70, idom_71, p_72, v114, f_115, bid_116, preds_117, ps_118, new_id_119, rpo_idx_120, idom_121, p_122, f_79, bid_80, preds_81, ps_82, new_id_83, rpo_idx_84, idom_85, p_86, v98, v100, f_87, bid_88, preds_89, ps_90, new_id_91, rpo_idx_92, idom_93, p_94, v104, f_105, bid_106, preds_107, ps_108, new_id_109, rpo_idx_110, idom_111, p_112
	preds_6, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-preds").Deref(), []vm.Value{arg1, arg0})
	if callErr != nil {
		return nil, callErr
	}
	ps_7 = preds_6
	new_id_8 = vm.Int(-1)
	rpo_idx_9 = arg3
	idom_10 = arg2
	goto b1
b1:
	;
	v28, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "empty?").Deref(), []vm.Value{ps_7})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v28) {
		f_13 = arg0
		bid_14 = arg1
		preds_15 = preds_6
		ps_16 = ps_7
		new_id_17 = new_id_8
		rpo_idx_18 = rpo_idx_9
		idom_19 = idom_10
		goto b2
	} else {
		f_20 = arg0
		bid_21 = arg1
		preds_22 = preds_6
		ps_23 = ps_7
		new_id_24 = new_id_8
		rpo_idx_25 = rpo_idx_9
		idom_26 = idom_10
		goto b3
	}
b2:
	;
	v134 = new_id_17
	f_135 = f_13
	bid_136 = bid_14
	preds_137 = preds_15
	ps_138 = ps_16
	new_id_139 = new_id_17
	rpo_idx_140 = rpo_idx_18
	idom_141 = idom_19
	goto b4
b3:
	;
	p_32, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{ps_23})
	if callErr != nil {
		return nil, callErr
	}
	arg__8702_51, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{idom_26, p_32})
	if callErr != nil {
		return nil, callErr
	}
	v52 = arg__8702_51 == vm.Int(-1)
	if v52 {
		f_33 = f_20
		bid_34 = bid_21
		preds_35 = preds_22
		ps_36 = ps_23
		new_id_37 = new_id_24
		rpo_idx_38 = rpo_idx_25
		idom_39 = idom_26
		p_40 = p_32
		goto b5
	} else {
		f_41 = f_20
		bid_42 = bid_21
		preds_43 = preds_22
		ps_44 = ps_23
		new_id_45 = new_id_24
		rpo_idx_46 = rpo_idx_25
		idom_47 = idom_26
		p_48 = p_32
		goto b6
	}
b4:
	;
	return v134, nil
b5:
	;
	v55, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{ps_36})
	if callErr != nil {
		return nil, callErr
	}
	ps_7 = v55
	new_id_8 = new_id_37
	rpo_idx_9 = rpo_idx_38
	idom_10 = idom_39
	goto b1
b6:
	;
	v74 = new_id_45 == vm.Int(-1)
	if v74 {
		f_57 = f_41
		bid_58 = bid_42
		preds_59 = preds_43
		ps_60 = ps_44
		new_id_61 = new_id_45
		rpo_idx_62 = rpo_idx_46
		idom_63 = idom_47
		p_64 = p_48
		goto b8
	} else {
		f_65 = f_41
		bid_66 = bid_42
		preds_67 = preds_43
		ps_68 = ps_44
		new_id_69 = new_id_45
		rpo_idx_70 = rpo_idx_46
		idom_71 = idom_47
		p_72 = p_48
		goto b9
	}
b7:
	;
	v134 = v124
	f_135 = f_125
	bid_136 = bid_126
	preds_137 = preds_127
	ps_138 = ps_128
	new_id_139 = new_id_129
	rpo_idx_140 = rpo_idx_130
	idom_141 = idom_131
	goto b4
b8:
	;
	v77, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{ps_60})
	if callErr != nil {
		return nil, callErr
	}
	ps_7 = v77
	new_id_8 = p_64
	rpo_idx_9 = rpo_idx_62
	idom_10 = idom_63
	goto b1
b9:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		f_79 = f_65
		bid_80 = bid_66
		preds_81 = preds_67
		ps_82 = ps_68
		new_id_83 = new_id_69
		rpo_idx_84 = rpo_idx_70
		idom_85 = idom_71
		p_86 = p_72
		goto b11
	} else {
		f_87 = f_65
		bid_88 = bid_66
		preds_89 = preds_67
		ps_90 = ps_68
		new_id_91 = new_id_69
		rpo_idx_92 = rpo_idx_70
		idom_93 = idom_71
		p_94 = p_72
		goto b12
	}
b10:
	;
	v124 = v114
	f_125 = f_115
	bid_126 = bid_116
	preds_127 = preds_117
	ps_128 = ps_118
	new_id_129 = new_id_119
	rpo_idx_130 = rpo_idx_120
	idom_131 = idom_121
	p_132 = p_122
	goto b7
b11:
	;
	v98, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{ps_82})
	if callErr != nil {
		return nil, callErr
	}
	v100, callErr = rt.InvokeValue(rt.LookupVar("ir.dominance", "intersect").Deref(), []vm.Value{idom_85, rpo_idx_84, p_86, new_id_83})
	if callErr != nil {
		return nil, callErr
	}
	ps_7 = v98
	new_id_8 = v100
	rpo_idx_9 = rpo_idx_84
	idom_10 = idom_85
	goto b1
b12:
	;
	v104 = vm.NIL
	f_105 = f_87
	bid_106 = bid_88
	preds_107 = preds_89
	ps_108 = ps_90
	new_id_109 = new_id_91
	rpo_idx_110 = rpo_idx_92
	idom_111 = idom_93
	p_112 = p_94
	goto b13
b13:
	;
	v114 = v104
	f_115 = f_105
	bid_116 = bid_106
	preds_117 = preds_107
	ps_118 = ps_108
	new_id_119 = new_id_109
	rpo_idx_120 = rpo_idx_110
	idom_121 = idom_111
	p_122 = p_112
	goto b10
}
func dominators(arg0 vm.Value) (vm.Value, error) {
	var arg__8726_3 vm.Value
	var arg__8731_6 vm.Value
	var n_7 vm.Value
	var entry_9 vm.Value
	var rpo_11 vm.Value
	var arg__8743_19 vm.Value
	var arg__8750_24 vm.Value
	var v25 vm.Value
	var i_12 int
	var idx_13 vm.Value
	var rpo_14 vm.Value
	var arg__8755_40 vm.Value
	var v41 bool
	var f_27 vm.Value
	var n_28 vm.Value
	var entry_29 vm.Value
	var i_30 int
	var idx_31 vm.Value
	var rpo_32 vm.Value
	var f_33 vm.Value
	var n_34 vm.Value
	var entry_35 vm.Value
	var i_36 int
	var idx_37 vm.Value
	var rpo_38 vm.Value
	var v44 int
	var arg__8763_46 vm.Value
	var arg__8772_49 vm.Value
	var v50 vm.Value
	var rpo_idx_52 vm.Value
	var f_53 vm.Value
	var n_54 vm.Value
	var entry_55 vm.Value
	var i_56 int
	var idx_57 vm.Value
	var rpo_58 vm.Value
	var arg__8779_62 vm.Value
	var arg__8786_67 vm.Value
	var arg__8787_68 vm.Value
	var arg__8796_73 vm.Value
	var arg__8803_78 vm.Value
	var arg__8804_79 vm.Value
	var idom0_80 vm.Value
	var idom_81 vm.Value
	var rpo_idx_82 vm.Value
	var rpo_83 vm.Value
	var f_84 vm.Value
	var entry_85 vm.Value
	var v327 int
	var v341 int
	var v355 int
	var bs_87 vm.Value
	var idom_88 vm.Value
	var changed_QMARK__89 vm.Value
	var rpo_idx_90 vm.Value
	var idom_91 vm.Value
	var f_92 vm.Value
	var entry_93 vm.Value
	var v325 int
	var v339 int
	var v353 int
	var v119 vm.Value
	var n_96 vm.Value
	var i_97 int
	var idx_98 vm.Value
	var idom0_99 vm.Value
	var rpo_100 vm.Value
	var bs_101 vm.Value
	var changed_QMARK__102 vm.Value
	var rpo_idx_103 vm.Value
	var idom_104 vm.Value
	var f_105 vm.Value
	var entry_106 vm.Value
	var v332 int
	var v346 int
	var v360 int
	var v122 vm.Value
	var n_107 vm.Value
	var i_108 int
	var idx_109 vm.Value
	var idom0_110 vm.Value
	var rpo_111 vm.Value
	var bs_112 vm.Value
	var changed_QMARK__113 vm.Value
	var rpo_idx_114 vm.Value
	var idom_115 vm.Value
	var f_116 vm.Value
	var entry_117 vm.Value
	var v323 int
	var v337 int
	var v351 int
	var b_125 vm.Value
	var v150 bool
	var step_244 vm.Value
	var n_245 vm.Value
	var i_246 int
	var idx_247 vm.Value
	var idom0_248 vm.Value
	var rpo_249 vm.Value
	var bs_250 vm.Value
	var changed_QMARK__251 vm.Value
	var rpo_idx_252 vm.Value
	var idom_253 vm.Value
	var f_254 vm.Value
	var entry_255 vm.Value
	var v328 int
	var v342 int
	var v356 int
	var v283 vm.Value
	var n_126 vm.Value
	var i_127 int
	var idx_128 vm.Value
	var idom0_129 vm.Value
	var rpo_130 vm.Value
	var bs_131 vm.Value
	var changed_QMARK__132 vm.Value
	var rpo_idx_133 vm.Value
	var idom_134 vm.Value
	var f_135 vm.Value
	var entry_136 vm.Value
	var b_137 vm.Value
	var v324 int
	var v338 int
	var v352 int
	var v153 vm.Value
	var n_138 vm.Value
	var i_139 int
	var idx_140 vm.Value
	var idom0_141 vm.Value
	var rpo_142 vm.Value
	var bs_143 vm.Value
	var changed_QMARK__144 vm.Value
	var rpo_idx_145 vm.Value
	var idom_146 vm.Value
	var f_147 vm.Value
	var entry_148 vm.Value
	var b_149 vm.Value
	var v321 int
	var v335 int
	var v349 int
	var ni_156 vm.Value
	var or__x_184 bool
	var n_157 vm.Value
	var i_158 int
	var idx_159 vm.Value
	var idom0_160 vm.Value
	var rpo_161 vm.Value
	var bs_162 vm.Value
	var changed_QMARK__163 vm.Value
	var rpo_idx_164 vm.Value
	var idom_165 vm.Value
	var f_166 vm.Value
	var entry_167 vm.Value
	var b_168 vm.Value
	var ni_169 vm.Value
	var v322 int
	var v336 int
	var v350 int
	var v236 vm.Value
	var n_170 vm.Value
	var i_171 int
	var idx_172 vm.Value
	var idom0_173 vm.Value
	var rpo_174 vm.Value
	var bs_175 vm.Value
	var changed_QMARK__176 vm.Value
	var rpo_idx_177 vm.Value
	var idom_178 vm.Value
	var f_179 vm.Value
	var entry_180 vm.Value
	var b_181 vm.Value
	var ni_182 vm.Value
	var v331 int
	var v345 int
	var v359 int
	var v239 vm.Value
	var v241 vm.Value
	var n_185 vm.Value
	var i_186 int
	var idx_187 vm.Value
	var idom0_188 vm.Value
	var rpo_189 vm.Value
	var bs_190 vm.Value
	var changed_QMARK__191 vm.Value
	var rpo_idx_192 vm.Value
	var idom_193 vm.Value
	var f_194 vm.Value
	var entry_195 vm.Value
	var b_196 vm.Value
	var ni_197 vm.Value
	var or__x_198 bool
	var v320 int
	var v334 int
	var v348 int
	var n_199 vm.Value
	var i_200 int
	var idx_201 vm.Value
	var idom0_202 vm.Value
	var rpo_203 vm.Value
	var bs_204 vm.Value
	var changed_QMARK__205 vm.Value
	var rpo_idx_206 vm.Value
	var idom_207 vm.Value
	var f_208 vm.Value
	var entry_209 vm.Value
	var b_210 vm.Value
	var ni_211 vm.Value
	var or__x_212 bool
	var v329 int
	var v343 int
	var v357 int
	var arg__8837_216 vm.Value
	var v217 bool
	var v219 bool
	var n_220 vm.Value
	var i_221 int
	var idx_222 vm.Value
	var idom0_223 vm.Value
	var rpo_224 vm.Value
	var bs_225 vm.Value
	var changed_QMARK__226 vm.Value
	var rpo_idx_227 vm.Value
	var idom_228 vm.Value
	var f_229 vm.Value
	var entry_230 vm.Value
	var b_231 vm.Value
	var ni_232 vm.Value
	var or__x_233 vm.Value
	var v326 int
	var v340 int
	var v354 int
	var step_256 vm.Value
	var n_257 vm.Value
	var i_258 int
	var idx_259 vm.Value
	var idom0_260 vm.Value
	var rpo_261 vm.Value
	var bs_262 vm.Value
	var changed_QMARK__263 vm.Value
	var rpo_idx_264 vm.Value
	var idom_265 vm.Value
	var f_266 vm.Value
	var entry_267 vm.Value
	var v330 int
	var v344 int
	var v358 int
	var v288 vm.Value
	var step_268 vm.Value
	var n_269 vm.Value
	var i_270 int
	var idx_271 vm.Value
	var idom0_272 vm.Value
	var rpo_273 vm.Value
	var bs_274 vm.Value
	var changed_QMARK__275 vm.Value
	var rpo_idx_276 vm.Value
	var idom_277 vm.Value
	var f_278 vm.Value
	var entry_279 vm.Value
	var v333 int
	var v347 int
	var v361 int
	var v293 vm.Value
	var final_295 vm.Value
	var step_296 vm.Value
	var n_297 vm.Value
	var i_298 int
	var idx_299 vm.Value
	var idom0_300 vm.Value
	var rpo_301 vm.Value
	var bs_302 vm.Value
	var changed_QMARK__303 vm.Value
	var rpo_idx_304 vm.Value
	var idom_305 vm.Value
	var f_306 vm.Value
	var entry_307 vm.Value
	var v311 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__8726_3, arg__8731_6, n_7, entry_9, rpo_11, arg__8743_19, arg__8750_24, v25, i_12, idx_13, rpo_14, arg__8755_40, v41, f_27, n_28, entry_29, i_30, idx_31, rpo_32, f_33, n_34, entry_35, i_36, idx_37, rpo_38, v44, arg__8763_46, arg__8772_49, v50, rpo_idx_52, f_53, n_54, entry_55, i_56, idx_57, rpo_58, arg__8779_62, arg__8786_67, arg__8787_68, arg__8796_73, arg__8803_78, arg__8804_79, idom0_80, idom_81, rpo_idx_82, rpo_83, f_84, entry_85, v327, v341, v355, bs_87, idom_88, changed_QMARK__89, rpo_idx_90, idom_91, f_92, entry_93, v325, v339, v353, v119, n_96, i_97, idx_98, idom0_99, rpo_100, bs_101, changed_QMARK__102, rpo_idx_103, idom_104, f_105, entry_106, v332, v346, v360, v122, n_107, i_108, idx_109, idom0_110, rpo_111, bs_112, changed_QMARK__113, rpo_idx_114, idom_115, f_116, entry_117, v323, v337, v351, b_125, v150, step_244, n_245, i_246, idx_247, idom0_248, rpo_249, bs_250, changed_QMARK__251, rpo_idx_252, idom_253, f_254, entry_255, v328, v342, v356, v283, n_126, i_127, idx_128, idom0_129, rpo_130, bs_131, changed_QMARK__132, rpo_idx_133, idom_134, f_135, entry_136, b_137, v324, v338, v352, v153, n_138, i_139, idx_140, idom0_141, rpo_142, bs_143, changed_QMARK__144, rpo_idx_145, idom_146, f_147, entry_148, b_149, v321, v335, v349, ni_156, or__x_184, n_157, i_158, idx_159, idom0_160, rpo_161, bs_162, changed_QMARK__163, rpo_idx_164, idom_165, f_166, entry_167, b_168, ni_169, v322, v336, v350, v236, n_170, i_171, idx_172, idom0_173, rpo_174, bs_175, changed_QMARK__176, rpo_idx_177, idom_178, f_179, entry_180, b_181, ni_182, v331, v345, v359, v239, v241, n_185, i_186, idx_187, idom0_188, rpo_189, bs_190, changed_QMARK__191, rpo_idx_192, idom_193, f_194, entry_195, b_196, ni_197, or__x_198, v320, v334, v348, n_199, i_200, idx_201, idom0_202, rpo_203, bs_204, changed_QMARK__205, rpo_idx_206, idom_207, f_208, entry_209, b_210, ni_211, or__x_212, v329, v343, v357, arg__8837_216, v217, v219, n_220, i_221, idx_222, idom0_223, rpo_224, bs_225, changed_QMARK__226, rpo_idx_227, idom_228, f_229, entry_230, b_231, ni_232, or__x_233, v326, v340, v354, step_256, n_257, i_258, idx_259, idom0_260, rpo_261, bs_262, changed_QMARK__263, rpo_idx_264, idom_265, f_266, entry_267, v330, v344, v358, v288, step_268, n_269, i_270, idx_271, idom0_272, rpo_273, bs_274, changed_QMARK__275, rpo_idx_276, idom_277, f_278, entry_279, v333, v347, v361, v293, final_295, step_296, n_297, i_298, idx_299, idom0_300, rpo_301, bs_302, changed_QMARK__303, rpo_idx_304, idom_305, f_306, entry_307, v311
	arg__8726_3, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__8731_6, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	n_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{arg__8731_6})
	if callErr != nil {
		return nil, callErr
	}
	entry_9, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-entry").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	rpo_11, callErr = rt.InvokeValue(rt.LookupVar("ir.dominance", "reverse-postorder").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__8743_19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "repeat").Deref(), []vm.Value{n_7, vm.Int(-1)})
	if callErr != nil {
		return nil, callErr
	}
	arg__8750_24, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "repeat").Deref(), []vm.Value{n_7, vm.Int(-1)})
	if callErr != nil {
		return nil, callErr
	}
	v25, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg__8750_24})
	if callErr != nil {
		return nil, callErr
	}
	i_12 = 0
	idx_13 = v25
	rpo_14 = rpo_11
	goto b1
b1:
	;
	arg__8755_40, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{rpo_14})
	if callErr != nil {
		return nil, callErr
	}
	v41 = rt.GeValue(vm.Int(i_12), arg__8755_40)
	if v41 {
		f_27 = arg0
		n_28 = n_7
		entry_29 = entry_9
		i_30 = i_12
		idx_31 = idx_13
		rpo_32 = rpo_14
		goto b2
	} else {
		f_33 = arg0
		n_34 = n_7
		entry_35 = entry_9
		i_36 = i_12
		idx_37 = idx_13
		rpo_38 = rpo_14
		goto b3
	}
b2:
	;
	rpo_idx_52 = idx_31
	f_53 = f_27
	n_54 = n_28
	entry_55 = entry_29
	i_56 = i_30
	idx_57 = idx_31
	rpo_58 = rpo_32
	goto b4
b3:
	;
	v44 = i_36 + 1
	arg__8763_46, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{rpo_38, vm.Int(i_36)})
	if callErr != nil {
		return nil, callErr
	}
	arg__8772_49, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{rpo_38, vm.Int(i_36)})
	if callErr != nil {
		return nil, callErr
	}
	v50, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "assoc").Deref(), []vm.Value{idx_37, arg__8772_49, vm.Int(i_36)})
	if callErr != nil {
		return nil, callErr
	}
	i_12 = v44
	idx_13 = v50
	rpo_14 = rpo_38
	goto b1
b4:
	;
	arg__8779_62, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "repeat").Deref(), []vm.Value{n_54, vm.Int(-1)})
	if callErr != nil {
		return nil, callErr
	}
	arg__8786_67, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "repeat").Deref(), []vm.Value{n_54, vm.Int(-1)})
	if callErr != nil {
		return nil, callErr
	}
	arg__8787_68, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg__8786_67})
	if callErr != nil {
		return nil, callErr
	}
	arg__8796_73, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "repeat").Deref(), []vm.Value{n_54, vm.Int(-1)})
	if callErr != nil {
		return nil, callErr
	}
	arg__8803_78, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "repeat").Deref(), []vm.Value{n_54, vm.Int(-1)})
	if callErr != nil {
		return nil, callErr
	}
	arg__8804_79, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg__8803_78})
	if callErr != nil {
		return nil, callErr
	}
	idom0_80, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "assoc").Deref(), []vm.Value{arg__8804_79, entry_55, entry_55})
	if callErr != nil {
		return nil, callErr
	}
	idom_81 = idom0_80
	rpo_idx_82 = rpo_idx_52
	rpo_83 = rpo_58
	f_84 = f_53
	entry_85 = entry_55
	v327 = -1
	v341 = 1
	v355 = 0
	goto b5
b5:
	;
	bs_87 = rpo_83
	idom_88 = idom_81
	changed_QMARK__89 = vm.Boolean(false)
	rpo_idx_90 = rpo_idx_82
	idom_91 = idom_81
	f_92 = f_84
	entry_93 = entry_85
	v325 = v327
	v339 = v341
	v353 = v355
	goto b6
b6:
	;
	v119, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "empty?").Deref(), []vm.Value{bs_87})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v119) {
		n_96 = n_54
		i_97 = i_56
		idx_98 = idx_57
		idom0_99 = idom0_80
		rpo_100 = rpo_83
		bs_101 = bs_87
		changed_QMARK__102 = changed_QMARK__89
		rpo_idx_103 = rpo_idx_90
		idom_104 = idom_91
		f_105 = f_92
		entry_106 = entry_93
		v332 = v325
		v346 = v339
		v360 = v353
		goto b7
	} else {
		n_107 = n_54
		i_108 = i_56
		idx_109 = idx_57
		idom0_110 = idom0_80
		rpo_111 = rpo_83
		bs_112 = bs_87
		changed_QMARK__113 = changed_QMARK__89
		rpo_idx_114 = rpo_idx_90
		idom_115 = idom_91
		f_116 = f_92
		entry_117 = entry_93
		v323 = v325
		v337 = v339
		v351 = v353
		goto b8
	}
b7:
	;
	v122, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{idom_104, changed_QMARK__102})
	if callErr != nil {
		return nil, callErr
	}
	step_244 = v122
	n_245 = n_96
	i_246 = i_97
	idx_247 = idx_98
	idom0_248 = idom0_99
	rpo_249 = rpo_100
	bs_250 = bs_101
	changed_QMARK__251 = changed_QMARK__102
	rpo_idx_252 = rpo_idx_103
	idom_253 = idom_104
	f_254 = f_105
	entry_255 = entry_106
	v328 = v332
	v342 = v346
	v356 = v360
	goto b9
b8:
	;
	b_125, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{bs_112})
	if callErr != nil {
		return nil, callErr
	}
	v150 = b_125 == entry_117
	if v150 {
		n_126 = n_107
		i_127 = i_108
		idx_128 = idx_109
		idom0_129 = idom0_110
		rpo_130 = rpo_111
		bs_131 = bs_112
		changed_QMARK__132 = changed_QMARK__113
		rpo_idx_133 = rpo_idx_114
		idom_134 = idom_115
		f_135 = f_116
		entry_136 = entry_117
		b_137 = b_125
		v324 = v323
		v338 = v337
		v352 = v351
		goto b10
	} else {
		n_138 = n_107
		i_139 = i_108
		idx_140 = idx_109
		idom0_141 = idom0_110
		rpo_142 = rpo_111
		bs_143 = bs_112
		changed_QMARK__144 = changed_QMARK__113
		rpo_idx_145 = rpo_idx_114
		idom_146 = idom_115
		f_147 = f_116
		entry_148 = entry_117
		b_149 = b_125
		v321 = v323
		v335 = v337
		v349 = v351
		goto b11
	}
b9:
	;
	v283, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{step_244, vm.Int(v342)})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v283) {
		step_256 = step_244
		n_257 = n_245
		i_258 = i_246
		idx_259 = idx_247
		idom0_260 = idom0_248
		rpo_261 = rpo_249
		bs_262 = bs_250
		changed_QMARK__263 = changed_QMARK__251
		rpo_idx_264 = rpo_idx_252
		idom_265 = idom_253
		f_266 = f_254
		entry_267 = entry_255
		v330 = v328
		v344 = v342
		v358 = v356
		goto b19
	} else {
		step_268 = step_244
		n_269 = n_245
		i_270 = i_246
		idx_271 = idx_247
		idom0_272 = idom0_248
		rpo_273 = rpo_249
		bs_274 = bs_250
		changed_QMARK__275 = changed_QMARK__251
		rpo_idx_276 = rpo_idx_252
		idom_277 = idom_253
		f_278 = f_254
		entry_279 = entry_255
		v333 = v328
		v347 = v342
		v361 = v356
		goto b20
	}
b10:
	;
	v153, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{bs_131})
	if callErr != nil {
		return nil, callErr
	}
	bs_87 = v153
	idom_88 = idom_134
	changed_QMARK__89 = changed_QMARK__132
	rpo_idx_90 = rpo_idx_133
	idom_91 = idom_134
	f_92 = f_135
	entry_93 = entry_136
	v325 = v324
	v339 = v338
	v353 = v352
	goto b6
b11:
	;
	ni_156, callErr = rt.InvokeValue(rt.LookupVar("ir.dominance", "refine-idom").Deref(), []vm.Value{f_147, b_149, idom_146, rpo_idx_145})
	if callErr != nil {
		return nil, callErr
	}
	or__x_184 = ni_156 == vm.Int(v321)
	if or__x_184 {
		n_185 = n_138
		i_186 = i_139
		idx_187 = idx_140
		idom0_188 = idom0_141
		rpo_189 = rpo_142
		bs_190 = bs_143
		changed_QMARK__191 = changed_QMARK__144
		rpo_idx_192 = rpo_idx_145
		idom_193 = idom_146
		f_194 = f_147
		entry_195 = entry_148
		b_196 = b_149
		ni_197 = ni_156
		or__x_198 = or__x_184
		v320 = v321
		v334 = v335
		v348 = v349
		goto b16
	} else {
		n_199 = n_138
		i_200 = i_139
		idx_201 = idx_140
		idom0_202 = idom0_141
		rpo_203 = rpo_142
		bs_204 = bs_143
		changed_QMARK__205 = changed_QMARK__144
		rpo_idx_206 = rpo_idx_145
		idom_207 = idom_146
		f_208 = f_147
		entry_209 = entry_148
		b_210 = b_149
		ni_211 = ni_156
		or__x_212 = or__x_184
		v329 = v321
		v343 = v335
		v357 = v349
		goto b17
	}
b13:
	;
	v236, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{bs_162})
	if callErr != nil {
		return nil, callErr
	}
	bs_87 = v236
	idom_88 = idom_165
	changed_QMARK__89 = changed_QMARK__163
	rpo_idx_90 = rpo_idx_164
	idom_91 = idom_165
	f_92 = f_166
	entry_93 = entry_167
	v325 = v322
	v339 = v336
	v353 = v350
	goto b6
b14:
	;
	v239, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{bs_175})
	if callErr != nil {
		return nil, callErr
	}
	v241, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "assoc").Deref(), []vm.Value{idom_178, b_181, ni_182})
	if callErr != nil {
		return nil, callErr
	}
	bs_87 = v239
	idom_88 = v241
	changed_QMARK__89 = vm.Boolean(true)
	rpo_idx_90 = rpo_idx_177
	idom_91 = idom_178
	f_92 = f_179
	entry_93 = entry_180
	v325 = v331
	v339 = v345
	v353 = v359
	goto b6
b16:
	;
	v219 = or__x_198
	n_220 = n_185
	i_221 = i_186
	idx_222 = idx_187
	idom0_223 = idom0_188
	rpo_224 = rpo_189
	bs_225 = bs_190
	changed_QMARK__226 = changed_QMARK__191
	rpo_idx_227 = rpo_idx_192
	idom_228 = idom_193
	f_229 = f_194
	entry_230 = entry_195
	b_231 = b_196
	ni_232 = ni_197
	or__x_233 = vm.Boolean(or__x_198)
	v326 = v320
	v340 = v334
	v354 = v348
	goto b18
b17:
	;
	arg__8837_216, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{idom_207, b_210})
	if callErr != nil {
		return nil, callErr
	}
	v217 = arg__8837_216 == ni_211
	v219 = v217
	n_220 = n_199
	i_221 = i_200
	idx_222 = idx_201
	idom0_223 = idom0_202
	rpo_224 = rpo_203
	bs_225 = bs_204
	changed_QMARK__226 = changed_QMARK__205
	rpo_idx_227 = rpo_idx_206
	idom_228 = idom_207
	f_229 = f_208
	entry_230 = entry_209
	b_231 = b_210
	ni_232 = ni_211
	or__x_233 = vm.Boolean(or__x_212)
	v326 = v329
	v340 = v343
	v354 = v357
	goto b18
b18:
	;
	if v219 {
		n_157 = n_220
		i_158 = i_221
		idx_159 = idx_222
		idom0_160 = idom0_223
		rpo_161 = rpo_224
		bs_162 = bs_225
		changed_QMARK__163 = changed_QMARK__226
		rpo_idx_164 = rpo_idx_227
		idom_165 = idom_228
		f_166 = f_229
		entry_167 = entry_230
		b_168 = b_231
		ni_169 = ni_232
		v322 = v326
		v336 = v340
		v350 = v354
		goto b13
	} else {
		n_170 = n_220
		i_171 = i_221
		idx_172 = idx_222
		idom0_173 = idom0_223
		rpo_174 = rpo_224
		bs_175 = bs_225
		changed_QMARK__176 = changed_QMARK__226
		rpo_idx_177 = rpo_idx_227
		idom_178 = idom_228
		f_179 = f_229
		entry_180 = entry_230
		b_181 = b_231
		ni_182 = ni_232
		v331 = v326
		v345 = v340
		v359 = v354
		goto b14
	}
b19:
	;
	v288, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{step_256, vm.Int(v358)})
	if callErr != nil {
		return nil, callErr
	}
	idom_81 = v288
	rpo_idx_82 = rpo_idx_264
	rpo_83 = rpo_261
	f_84 = f_266
	entry_85 = entry_267
	v327 = v330
	v341 = v344
	v355 = v358
	goto b5
b20:
	;
	v293, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{step_268, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	final_295 = v293
	step_296 = step_268
	n_297 = n_269
	i_298 = i_270
	idx_299 = idx_271
	idom0_300 = idom0_272
	rpo_301 = rpo_273
	bs_302 = bs_274
	changed_QMARK__303 = changed_QMARK__275
	rpo_idx_304 = rpo_idx_276
	idom_305 = idom_277
	f_306 = f_278
	entry_307 = entry_279
	goto b21
b21:
	;
	v311, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "assoc").Deref(), []vm.Value{final_295, entry_307, vm.Int(-1)})
	if callErr != nil {
		return nil, callErr
	}
	return v311, nil
}
func dominates_QMARK_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var idom_5 vm.Value
	var b_6 vm.Value
	var idom_7 vm.Value
	var b_8 vm.Value
	var a_9 vm.Value
	var v72 int
	var v79 vm.Value
	var v20 bool
	var f_11 vm.Value
	var idom_12 vm.Value
	var b_13 vm.Value
	var a_14 vm.Value
	var v75 int
	var v82 vm.Value
	var f_15 vm.Value
	var idom_16 vm.Value
	var b_17 vm.Value
	var a_18 vm.Value
	var v73 int
	var v80 vm.Value
	var v32 bool
	var v63 vm.Value
	var f_64 vm.Value
	var idom_65 vm.Value
	var b_66 vm.Value
	var a_67 vm.Value
	var f_24 vm.Value
	var idom_25 vm.Value
	var b_26 vm.Value
	var a_27 vm.Value
	var v76 int
	var v83 vm.Value
	var f_28 vm.Value
	var idom_29 vm.Value
	var b_30 vm.Value
	var a_31 vm.Value
	var v74 int
	var v81 vm.Value
	var v57 vm.Value
	var f_58 vm.Value
	var idom_59 vm.Value
	var b_60 vm.Value
	var a_61 vm.Value
	var f_36 vm.Value
	var idom_37 vm.Value
	var b_38 vm.Value
	var a_39 vm.Value
	var v71 int
	var v78 vm.Value
	var v47 vm.Value
	var f_40 vm.Value
	var idom_41 vm.Value
	var b_42 vm.Value
	var a_43 vm.Value
	var v77 int
	var v84 vm.Value
	var v51 vm.Value
	var f_52 vm.Value
	var idom_53 vm.Value
	var b_54 vm.Value
	var a_55 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = idom_5, b_6, idom_7, b_8, a_9, v72, v79, v20, f_11, idom_12, b_13, a_14, v75, v82, f_15, idom_16, b_17, a_18, v73, v80, v32, v63, f_64, idom_65, b_66, a_67, f_24, idom_25, b_26, a_27, v76, v83, f_28, idom_29, b_30, a_31, v74, v81, v57, f_58, idom_59, b_60, a_61, f_36, idom_37, b_38, a_39, v71, v78, v47, f_40, idom_41, b_42, a_43, v77, v84, v51, f_52, idom_53, b_54, a_55
	idom_5, callErr = rt.InvokeValue(rt.LookupVar("ir.dominance", "dominators").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	b_6 = arg2
	idom_7 = idom_5
	b_8 = arg2
	a_9 = arg1
	v72 = -1
	v79 = vm.Keyword("else")
	goto b1
b1:
	;
	v20 = b_8 == vm.Int(v72)
	if v20 {
		f_11 = arg0
		idom_12 = idom_7
		b_13 = b_8
		a_14 = a_9
		v75 = v72
		v82 = v79
		goto b2
	} else {
		f_15 = arg0
		idom_16 = idom_7
		b_17 = b_8
		a_18 = a_9
		v73 = v72
		v80 = v79
		goto b3
	}
b2:
	;
	v63 = vm.Boolean(false)
	f_64 = f_11
	idom_65 = idom_12
	b_66 = b_13
	a_67 = a_14
	goto b4
b3:
	;
	v32 = a_18 == b_17
	if v32 {
		f_24 = f_15
		idom_25 = idom_16
		b_26 = b_17
		a_27 = a_18
		v76 = v73
		v83 = v80
		goto b5
	} else {
		f_28 = f_15
		idom_29 = idom_16
		b_30 = b_17
		a_31 = a_18
		v74 = v73
		v81 = v80
		goto b6
	}
b4:
	;
	return v63, nil
b5:
	;
	v57 = vm.Boolean(true)
	f_58 = f_24
	idom_59 = idom_25
	b_60 = b_26
	a_61 = a_27
	goto b7
b6:
	;
	if vm.IsTruthy(v81) {
		f_36 = f_28
		idom_37 = idom_29
		b_38 = b_30
		a_39 = a_31
		v71 = v74
		v78 = v81
		goto b8
	} else {
		f_40 = f_28
		idom_41 = idom_29
		b_42 = b_30
		a_43 = a_31
		v77 = v74
		v84 = v81
		goto b9
	}
b7:
	;
	v63 = v57
	f_64 = f_58
	idom_65 = idom_59
	b_66 = b_60
	a_67 = a_61
	goto b4
b8:
	;
	v47, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{idom_37, b_38})
	if callErr != nil {
		return nil, callErr
	}
	b_6 = v47
	idom_7 = idom_37
	b_8 = b_38
	a_9 = a_39
	v72 = v71
	v79 = v78
	goto b1
b9:
	;
	v51 = vm.NIL
	f_52 = f_40
	idom_53 = idom_41
	b_54 = b_42
	a_55 = a_43
	goto b10
b10:
	;
	v57 = v51
	f_58 = f_52
	idom_59 = idom_53
	b_60 = b_54
	a_61 = a_55
	goto b7
}
