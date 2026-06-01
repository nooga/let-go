package ir_build

import (
	rt "github.com/nooga/let-go/pkg/rt"
	vm "github.com/nooga/let-go/pkg/vm"
)

func add_inst_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value, arg3 vm.Value, arg4 vm.Value) (vm.Value, error) {
	var f_6 vm.Value
	var nid_8 vm.Value
	var arg__251_10 vm.Value
	var arg__257_14 vm.Value
	var si_16 vm.Value
	var ctx_17 vm.Value
	var block_18 vm.Value
	var op_kw_19 vm.Value
	var refs_20 vm.Value
	var aux_21 vm.Value
	var f_22 vm.Value
	var nid_23 vm.Value
	var si_24 vm.Value
	var v35 vm.Value
	var ctx_25 vm.Value
	var block_26 vm.Value
	var op_kw_27 vm.Value
	var refs_28 vm.Value
	var aux_29 vm.Value
	var f_30 vm.Value
	var nid_31 vm.Value
	var si_32 vm.Value
	var v39 vm.Value
	var ctx_40 vm.Value
	var block_41 vm.Value
	var op_kw_42 vm.Value
	var refs_43 vm.Value
	var aux_44 vm.Value
	var f_45 vm.Value
	var nid_46 vm.Value
	var si_47 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = f_6, nid_8, arg__251_10, arg__257_14, si_16, ctx_17, block_18, op_kw_19, refs_20, aux_21, f_22, nid_23, si_24, v35, ctx_25, block_26, op_kw_27, refs_28, aux_29, f_30, nid_31, si_32, v39, ctx_40, block_41, op_kw_42, refs_43, aux_44, f_45, nid_46, si_47
	f_6, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-fn").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	nid_8, callErr = rt.InvokeValue(rt.LookupVar("ir", "add-inst").Deref(), []vm.Value{f_6, arg1, arg2, arg3, arg4})
	if callErr != nil {
		return nil, callErr
	}
	arg__251_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__257_14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	si_16, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg__257_14, vm.Keyword("source-info")})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(si_16) {
		ctx_17 = arg0
		block_18 = arg1
		op_kw_19 = arg2
		refs_20 = arg3
		aux_21 = arg4
		f_22 = f_6
		nid_23 = nid_8
		si_24 = si_16
		goto b1
	} else {
		ctx_25 = arg0
		block_26 = arg1
		op_kw_27 = arg2
		refs_28 = arg3
		aux_29 = arg4
		f_30 = f_6
		nid_31 = nid_8
		si_32 = si_16
		goto b2
	}
b1:
	;
	v35, callErr = rt.InvokeValue(rt.LookupVar("ir", "add-source-info!").Deref(), []vm.Value{f_22, nid_23, si_24})
	if callErr != nil {
		return nil, callErr
	}
	v39 = v35
	ctx_40 = ctx_17
	block_41 = block_18
	op_kw_42 = op_kw_19
	refs_43 = refs_20
	aux_44 = aux_21
	f_45 = f_22
	nid_46 = nid_23
	si_47 = si_24
	goto b3
b2:
	;
	v39 = vm.NIL
	ctx_40 = ctx_25
	block_41 = block_26
	op_kw_42 = op_kw_27
	refs_43 = refs_28
	aux_44 = aux_29
	f_45 = f_30
	nid_46 = nid_31
	si_47 = si_32
	goto b3
b3:
	;
	return nid_46, nil
}
func add_terminator_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value, arg3 vm.Value, arg4 vm.Value) (vm.Value, error) {
	var f_6 vm.Value
	var nid_8 vm.Value
	var arg__283_10 vm.Value
	var arg__289_14 vm.Value
	var si_16 vm.Value
	var ctx_17 vm.Value
	var bid_18 vm.Value
	var op_kw_19 vm.Value
	var refs_20 vm.Value
	var aux_21 vm.Value
	var f_22 vm.Value
	var nid_23 vm.Value
	var si_24 vm.Value
	var v35 vm.Value
	var ctx_25 vm.Value
	var bid_26 vm.Value
	var op_kw_27 vm.Value
	var refs_28 vm.Value
	var aux_29 vm.Value
	var f_30 vm.Value
	var nid_31 vm.Value
	var si_32 vm.Value
	var v39 vm.Value
	var ctx_40 vm.Value
	var bid_41 vm.Value
	var op_kw_42 vm.Value
	var refs_43 vm.Value
	var aux_44 vm.Value
	var f_45 vm.Value
	var nid_46 vm.Value
	var si_47 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = f_6, nid_8, arg__283_10, arg__289_14, si_16, ctx_17, bid_18, op_kw_19, refs_20, aux_21, f_22, nid_23, si_24, v35, ctx_25, bid_26, op_kw_27, refs_28, aux_29, f_30, nid_31, si_32, v39, ctx_40, bid_41, op_kw_42, refs_43, aux_44, f_45, nid_46, si_47
	f_6, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-fn").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	nid_8, callErr = rt.InvokeValue(rt.LookupVar("ir", "add-terminator!").Deref(), []vm.Value{f_6, arg1, arg2, arg3, arg4})
	if callErr != nil {
		return nil, callErr
	}
	arg__283_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__289_14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	si_16, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg__289_14, vm.Keyword("source-info")})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(si_16) {
		ctx_17 = arg0
		bid_18 = arg1
		op_kw_19 = arg2
		refs_20 = arg3
		aux_21 = arg4
		f_22 = f_6
		nid_23 = nid_8
		si_24 = si_16
		goto b1
	} else {
		ctx_25 = arg0
		bid_26 = arg1
		op_kw_27 = arg2
		refs_28 = arg3
		aux_29 = arg4
		f_30 = f_6
		nid_31 = nid_8
		si_32 = si_16
		goto b2
	}
b1:
	;
	v35, callErr = rt.InvokeValue(rt.LookupVar("ir", "add-source-info!").Deref(), []vm.Value{f_22, nid_23, si_24})
	if callErr != nil {
		return nil, callErr
	}
	v39 = v35
	ctx_40 = ctx_17
	bid_41 = bid_18
	op_kw_42 = op_kw_19
	refs_43 = refs_20
	aux_44 = aux_21
	f_45 = f_22
	nid_46 = nid_23
	si_47 = si_24
	goto b3
b2:
	;
	v39 = vm.NIL
	ctx_40 = ctx_25
	bid_41 = bid_26
	op_kw_42 = op_kw_27
	refs_43 = refs_28
	aux_44 = aux_29
	f_45 = f_30
	nid_46 = nid_31
	si_47 = si_32
	goto b3
b3:
	;
	return nid_46, nil
}
func attach_name_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var and__x_11 vm.Value
	var ctx_4 vm.Value
	var inst_id_5 vm.Value
	var sym_6 vm.Value
	var f_32 vm.Value
	var arg__311_34 vm.Value
	var arg__316_37 vm.Value
	var arg__317_38 vm.Value
	var arg__324_41 vm.Value
	var arg__329_44 vm.Value
	var arg__330_45 vm.Value
	var v46 vm.Value
	var ctx_7 vm.Value
	var inst_id_8 vm.Value
	var sym_9 vm.Value
	var v50 vm.Value
	var ctx_51 vm.Value
	var inst_id_52 vm.Value
	var sym_53 vm.Value
	var ctx_12 vm.Value
	var inst_id_13 vm.Value
	var sym_14 vm.Value
	var and__x_15 vm.Value
	var v22 bool
	var ctx_16 vm.Value
	var inst_id_17 vm.Value
	var sym_18 vm.Value
	var and__x_19 vm.Value
	var v25 vm.Value
	var ctx_26 vm.Value
	var inst_id_27 vm.Value
	var sym_28 vm.Value
	var and__x_29 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = and__x_11, ctx_4, inst_id_5, sym_6, f_32, arg__311_34, arg__316_37, arg__317_38, arg__324_41, arg__329_44, arg__330_45, v46, ctx_7, inst_id_8, sym_9, v50, ctx_51, inst_id_52, sym_53, ctx_12, inst_id_13, sym_14, and__x_15, v22, ctx_16, inst_id_17, sym_18, and__x_19, v25, ctx_26, inst_id_27, sym_28, and__x_29
	and__x_11, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "integer?").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(and__x_11) {
		ctx_12 = arg0
		inst_id_13 = arg1
		sym_14 = arg2
		and__x_15 = and__x_11
		goto b4
	} else {
		ctx_16 = arg0
		inst_id_17 = arg1
		sym_18 = arg2
		and__x_19 = and__x_11
		goto b5
	}
b1:
	;
	f_32, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-fn").Deref(), []vm.Value{ctx_4})
	if callErr != nil {
		return nil, callErr
	}
	arg__311_34, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{sym_6})
	if callErr != nil {
		return nil, callErr
	}
	arg__316_37, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{sym_6})
	if callErr != nil {
		return nil, callErr
	}
	arg__317_38, callErr = rt.InvokeValue(rt.LookupVar("ir", "new-named-source-info").Deref(), []vm.Value{arg__316_37})
	if callErr != nil {
		return nil, callErr
	}
	arg__324_41, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{sym_6})
	if callErr != nil {
		return nil, callErr
	}
	arg__329_44, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{sym_6})
	if callErr != nil {
		return nil, callErr
	}
	arg__330_45, callErr = rt.InvokeValue(rt.LookupVar("ir", "new-named-source-info").Deref(), []vm.Value{arg__329_44})
	if callErr != nil {
		return nil, callErr
	}
	v46, callErr = rt.InvokeValue(rt.LookupVar("ir", "add-source-info!").Deref(), []vm.Value{f_32, inst_id_5, arg__330_45})
	if callErr != nil {
		return nil, callErr
	}
	v50 = v46
	ctx_51 = ctx_4
	inst_id_52 = inst_id_5
	sym_53 = sym_6
	goto b3
b2:
	;
	v50 = vm.NIL
	ctx_51 = ctx_7
	inst_id_52 = inst_id_8
	sym_53 = sym_9
	goto b3
b3:
	;
	return v50, nil
b4:
	;
	v22 = rt.GeValue(inst_id_13, vm.Int(0))
	v25 = vm.Boolean(v22)
	ctx_26 = ctx_12
	inst_id_27 = inst_id_13
	sym_28 = sym_14
	and__x_29 = and__x_15
	goto b6
b5:
	;
	v25 = and__x_19
	ctx_26 = ctx_16
	inst_id_27 = inst_id_17
	sym_28 = sym_18
	and__x_29 = and__x_19
	goto b6
b6:
	;
	if vm.IsTruthy(v25) {
		ctx_4 = ctx_26
		inst_id_5 = inst_id_27
		sym_6 = sym_28
		goto b1
	} else {
		ctx_7 = ctx_26
		inst_id_8 = inst_id_27
		sym_9 = sym_28
		goto b2
	}
}
func bind_local_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var s_4 vm.Value
	var stack_8 vm.Value
	var arg__342_10 vm.Value
	var top_idx_11 vm.Value
	var top_13 vm.Value
	var arg__360_16 vm.Value
	var arg__371_19 vm.Value
	var arg__372_20 vm.Value
	var arg__385_24 vm.Value
	var arg__396_27 vm.Value
	var arg__397_28 vm.Value
	var arg__398_29 vm.Value
	var arg__412_33 vm.Value
	var arg__423_36 vm.Value
	var arg__424_37 vm.Value
	var arg__437_41 vm.Value
	var arg__448_44 vm.Value
	var arg__449_45 vm.Value
	var arg__450_46 vm.Value
	var v47 vm.Value
	var v49 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = s_4, stack_8, arg__342_10, top_idx_11, top_13, arg__360_16, arg__371_19, arg__372_20, arg__385_24, arg__396_27, arg__397_28, arg__398_29, arg__412_33, arg__423_36, arg__424_37, arg__437_41, arg__448_44, arg__449_45, arg__450_46, v47, v49
	s_4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	stack_8, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{s_4, vm.Keyword("locals")})
	if callErr != nil {
		return nil, callErr
	}
	arg__342_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{stack_8})
	if callErr != nil {
		return nil, callErr
	}
	top_idx_11 = rt.SubValue(arg__342_10, vm.Int(1))
	top_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{stack_8, top_idx_11})
	if callErr != nil {
		return nil, callErr
	}
	arg__360_16, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "assoc").Deref(), []vm.Value{top_13, arg1, arg2})
	if callErr != nil {
		return nil, callErr
	}
	arg__371_19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "assoc").Deref(), []vm.Value{top_13, arg1, arg2})
	if callErr != nil {
		return nil, callErr
	}
	arg__372_20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "assoc").Deref(), []vm.Value{stack_8, top_idx_11, arg__371_19})
	if callErr != nil {
		return nil, callErr
	}
	arg__385_24, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "assoc").Deref(), []vm.Value{top_13, arg1, arg2})
	if callErr != nil {
		return nil, callErr
	}
	arg__396_27, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "assoc").Deref(), []vm.Value{top_13, arg1, arg2})
	if callErr != nil {
		return nil, callErr
	}
	arg__397_28, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "assoc").Deref(), []vm.Value{stack_8, top_idx_11, arg__396_27})
	if callErr != nil {
		return nil, callErr
	}
	arg__398_29, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "assoc").Deref(), []vm.Value{s_4, vm.Keyword("locals"), arg__397_28})
	if callErr != nil {
		return nil, callErr
	}
	arg__412_33, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "assoc").Deref(), []vm.Value{top_13, arg1, arg2})
	if callErr != nil {
		return nil, callErr
	}
	arg__423_36, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "assoc").Deref(), []vm.Value{top_13, arg1, arg2})
	if callErr != nil {
		return nil, callErr
	}
	arg__424_37, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "assoc").Deref(), []vm.Value{stack_8, top_idx_11, arg__423_36})
	if callErr != nil {
		return nil, callErr
	}
	arg__437_41, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "assoc").Deref(), []vm.Value{top_13, arg1, arg2})
	if callErr != nil {
		return nil, callErr
	}
	arg__448_44, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "assoc").Deref(), []vm.Value{top_13, arg1, arg2})
	if callErr != nil {
		return nil, callErr
	}
	arg__449_45, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "assoc").Deref(), []vm.Value{stack_8, top_idx_11, arg__448_44})
	if callErr != nil {
		return nil, callErr
	}
	arg__450_46, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "assoc").Deref(), []vm.Value{s_4, vm.Keyword("locals"), arg__449_45})
	if callErr != nil {
		return nil, callErr
	}
	v47, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "reset!").Deref(), []vm.Value{arg0, arg__450_46})
	if callErr != nil {
		return nil, callErr
	}
	v49, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "attach-name!").Deref(), []vm.Value{arg0, arg2, arg1})
	if callErr != nil {
		return nil, callErr
	}
	return v49, nil
}
func binding_syms(arg0 vm.Value) (vm.Value, error) {
	var v4 vm.Value
	var sym_1 vm.Value
	var v7 vm.Value
	var sym_2 vm.Value
	var v12 vm.Value
	var v251 vm.Value
	var sym_252 vm.Value
	var sym_9 vm.Value
	var arg__468_15 vm.Value
	var arg__473_18 vm.Value
	var arg__479_22 vm.Value
	var arg__480_23 vm.Value
	var arg__483_26 vm.Value
	var arg__488_29 vm.Value
	var arg__494_33 vm.Value
	var arg__495_34 vm.Value
	var v35 vm.Value
	var sym_10 vm.Value
	var v40 vm.Value
	var v248 vm.Value
	var sym_249 vm.Value
	var sym_37 vm.Value
	var tem__G__0_43 vm.Value
	var sym_38 vm.Value
	var v245 vm.Value
	var sym_246 vm.Value
	var sym_44 vm.Value
	var tem__G__0_45 vm.Value
	var arg__502_50 vm.Value
	var arg__532_54 vm.Value
	var arg__535_57 vm.Value
	var arg__565_61 vm.Value
	var v62 vm.Value
	var sym_46 vm.Value
	var tem__G__0_47 vm.Value
	var keys_syms_66 vm.Value
	var sym_67 vm.Value
	var tem__G__0_68 vm.Value
	var tem__G__0_70 vm.Value
	var keys_syms_71 vm.Value
	var sym_72 vm.Value
	var tem__G__0_73 vm.Value
	var arg__569_79 vm.Value
	var arg__575_83 vm.Value
	var arg__578_86 vm.Value
	var arg__584_90 vm.Value
	var v91 vm.Value
	var keys_syms_74 vm.Value
	var sym_75 vm.Value
	var tem__G__0_76 vm.Value
	var strs_syms_95 vm.Value
	var keys_syms_96 vm.Value
	var sym_97 vm.Value
	var tem__G__0_98 vm.Value
	var tem__G__0_100 vm.Value
	var strs_syms_101 vm.Value
	var keys_syms_102 vm.Value
	var sym_103 vm.Value
	var tem__G__0_104 vm.Value
	var v111 vm.Value
	var strs_syms_105 vm.Value
	var keys_syms_106 vm.Value
	var sym_107 vm.Value
	var tem__G__0_108 vm.Value
	var as_sym_115 vm.Value
	var strs_syms_116 vm.Value
	var keys_syms_117 vm.Value
	var sym_118 vm.Value
	var tem__G__0_119 vm.Value
	var arg__591_121 vm.Value
	var arg__604_132 vm.Value
	var arg__617_143 vm.Value
	var arg__618_144 vm.Value
	var arg__632_156 vm.Value
	var arg__645_167 vm.Value
	var arg__646_168 vm.Value
	var arg__647_169 vm.Value
	var arg__650_172 vm.Value
	var arg__663_183 vm.Value
	var arg__676_194 vm.Value
	var arg__677_195 vm.Value
	var arg__691_207 vm.Value
	var arg__704_218 vm.Value
	var arg__705_219 vm.Value
	var arg__706_220 vm.Value
	var other_syms_221 vm.Value
	var arg__708_223 vm.Value
	var arg__718_225 vm.Value
	var arg__721_228 vm.Value
	var arg__731_230 vm.Value
	var v231 vm.Value
	var sym_233 vm.Value
	var v238 vm.Value
	var sym_234 vm.Value
	var v242 vm.Value
	var sym_243 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = v4, sym_1, v7, sym_2, v12, v251, sym_252, sym_9, arg__468_15, arg__473_18, arg__479_22, arg__480_23, arg__483_26, arg__488_29, arg__494_33, arg__495_34, v35, sym_10, v40, v248, sym_249, sym_37, tem__G__0_43, sym_38, v245, sym_246, sym_44, tem__G__0_45, arg__502_50, arg__532_54, arg__535_57, arg__565_61, v62, sym_46, tem__G__0_47, keys_syms_66, sym_67, tem__G__0_68, tem__G__0_70, keys_syms_71, sym_72, tem__G__0_73, arg__569_79, arg__575_83, arg__578_86, arg__584_90, v91, keys_syms_74, sym_75, tem__G__0_76, strs_syms_95, keys_syms_96, sym_97, tem__G__0_98, tem__G__0_100, strs_syms_101, keys_syms_102, sym_103, tem__G__0_104, v111, strs_syms_105, keys_syms_106, sym_107, tem__G__0_108, as_sym_115, strs_syms_116, keys_syms_117, sym_118, tem__G__0_119, arg__591_121, arg__604_132, arg__617_143, arg__618_144, arg__632_156, arg__645_167, arg__646_168, arg__647_169, arg__650_172, arg__663_183, arg__676_194, arg__677_195, arg__691_207, arg__704_218, arg__705_219, arg__706_220, other_syms_221, arg__708_223, arg__718_225, arg__721_228, arg__731_230, v231, sym_233, v238, sym_234, v242, sym_243
	v4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "symbol?").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v4) {
		sym_1 = arg0
		goto b1
	} else {
		sym_2 = arg0
		goto b2
	}
b1:
	;
	v7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{sym_1})
	if callErr != nil {
		return nil, callErr
	}
	v251 = v7
	sym_252 = sym_1
	goto b3
b2:
	;
	v12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector?").Deref(), []vm.Value{sym_2})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v12) {
		sym_9 = sym_2
		goto b4
	} else {
		sym_10 = sym_2
		goto b5
	}
b3:
	;
	return v251, nil
b4:
	;
	arg__468_15, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__473_18, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "flatten").Deref(), []vm.Value{sym_9})
	if callErr != nil {
		return nil, callErr
	}
	arg__479_22, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "flatten").Deref(), []vm.Value{sym_9})
	if callErr != nil {
		return nil, callErr
	}
	arg__480_23, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "filter").Deref(), []vm.Value{rt.LookupVar("clojure.core", "symbol?").Deref(), arg__479_22})
	if callErr != nil {
		return nil, callErr
	}
	arg__483_26, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__488_29, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "flatten").Deref(), []vm.Value{sym_9})
	if callErr != nil {
		return nil, callErr
	}
	arg__494_33, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "flatten").Deref(), []vm.Value{sym_9})
	if callErr != nil {
		return nil, callErr
	}
	arg__495_34, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "filter").Deref(), []vm.Value{rt.LookupVar("clojure.core", "symbol?").Deref(), arg__494_33})
	if callErr != nil {
		return nil, callErr
	}
	v35, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__483_26, arg__495_34})
	if callErr != nil {
		return nil, callErr
	}
	v248 = v35
	sym_249 = sym_9
	goto b6
b5:
	;
	v40, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map?").Deref(), []vm.Value{sym_10})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v40) {
		sym_37 = sym_10
		goto b7
	} else {
		sym_38 = sym_10
		goto b8
	}
b6:
	;
	v251 = v248
	sym_252 = sym_249
	goto b3
b7:
	;
	tem__G__0_43, callErr = rt.InvokeValue(vm.Keyword("keys"), []vm.Value{sym_37})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(tem__G__0_43) {
		sym_44 = sym_37
		tem__G__0_45 = tem__G__0_43
		goto b10
	} else {
		sym_46 = sym_37
		tem__G__0_47 = tem__G__0_43
		goto b11
	}
b8:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		sym_233 = sym_38
		goto b19
	} else {
		sym_234 = sym_38
		goto b20
	}
b9:
	;
	v248 = v245
	sym_249 = sym_246
	goto b6
b10:
	;
	arg__502_50, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__532_54, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v4 vm.Value
		var x_1 vm.Value
		var arg__524_7 vm.Value
		var arg__529_10 vm.Value
		var v11 vm.Value
		var x_2 vm.Value
		var v14 vm.Value
		var x_15 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _ = v4, x_1, arg__524_7, arg__529_10, v11, x_2, v14, x_15
		v4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "keyword?").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		if vm.IsTruthy(v4) {
			x_1 = arg0
			goto b1
		} else {
			x_2 = arg0
			goto b2
		}
	b1:
		;
		arg__524_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "name").Deref(), []vm.Value{x_1})
		if callErr != nil {
			return nil, callErr
		}
		arg__529_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "name").Deref(), []vm.Value{x_1})
		if callErr != nil {
			return nil, callErr
		}
		v11, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "symbol").Deref(), []vm.Value{arg__529_10})
		if callErr != nil {
			return nil, callErr
		}
		v14 = v11
		x_15 = x_1
		goto b3
	b2:
		;
		v14 = x_2
		x_15 = x_2
		goto b3
	b3:
		;
		return v14, nil
	}), tem__G__0_45})
	if callErr != nil {
		return nil, callErr
	}
	arg__535_57, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__565_61, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v4 vm.Value
		var x_1 vm.Value
		var arg__557_7 vm.Value
		var arg__562_10 vm.Value
		var v11 vm.Value
		var x_2 vm.Value
		var v14 vm.Value
		var x_15 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _ = v4, x_1, arg__557_7, arg__562_10, v11, x_2, v14, x_15
		v4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "keyword?").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		if vm.IsTruthy(v4) {
			x_1 = arg0
			goto b1
		} else {
			x_2 = arg0
			goto b2
		}
	b1:
		;
		arg__557_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "name").Deref(), []vm.Value{x_1})
		if callErr != nil {
			return nil, callErr
		}
		arg__562_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "name").Deref(), []vm.Value{x_1})
		if callErr != nil {
			return nil, callErr
		}
		v11, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "symbol").Deref(), []vm.Value{arg__562_10})
		if callErr != nil {
			return nil, callErr
		}
		v14 = v11
		x_15 = x_1
		goto b3
	b2:
		;
		v14 = x_2
		x_15 = x_2
		goto b3
	b3:
		;
		return v14, nil
	}), tem__G__0_45})
	if callErr != nil {
		return nil, callErr
	}
	v62, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__535_57, arg__565_61})
	if callErr != nil {
		return nil, callErr
	}
	keys_syms_66 = v62
	sym_67 = sym_44
	tem__G__0_68 = tem__G__0_45
	goto b12
b11:
	;
	keys_syms_66 = vm.NIL
	sym_67 = sym_46
	tem__G__0_68 = tem__G__0_47
	goto b12
b12:
	;
	tem__G__0_70, callErr = rt.InvokeValue(vm.Keyword("strs"), []vm.Value{sym_67})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(tem__G__0_70) {
		keys_syms_71 = keys_syms_66
		sym_72 = sym_67
		tem__G__0_73 = tem__G__0_70
		goto b13
	} else {
		keys_syms_74 = keys_syms_66
		sym_75 = sym_67
		tem__G__0_76 = tem__G__0_70
		goto b14
	}
b13:
	;
	arg__569_79, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__575_83, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.LookupVar("clojure.core", "symbol").Deref(), tem__G__0_73})
	if callErr != nil {
		return nil, callErr
	}
	arg__578_86, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__584_90, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.LookupVar("clojure.core", "symbol").Deref(), tem__G__0_73})
	if callErr != nil {
		return nil, callErr
	}
	v91, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__578_86, arg__584_90})
	if callErr != nil {
		return nil, callErr
	}
	strs_syms_95 = v91
	keys_syms_96 = keys_syms_71
	sym_97 = sym_72
	tem__G__0_98 = tem__G__0_73
	goto b15
b14:
	;
	strs_syms_95 = vm.NIL
	keys_syms_96 = keys_syms_74
	sym_97 = sym_75
	tem__G__0_98 = tem__G__0_76
	goto b15
b15:
	;
	tem__G__0_100, callErr = rt.InvokeValue(vm.Keyword("as"), []vm.Value{sym_97})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(tem__G__0_100) {
		strs_syms_101 = strs_syms_95
		keys_syms_102 = keys_syms_96
		sym_103 = sym_97
		tem__G__0_104 = tem__G__0_100
		goto b16
	} else {
		strs_syms_105 = strs_syms_95
		keys_syms_106 = keys_syms_96
		sym_107 = sym_97
		tem__G__0_108 = tem__G__0_100
		goto b17
	}
b16:
	;
	v111, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{tem__G__0_104})
	if callErr != nil {
		return nil, callErr
	}
	as_sym_115 = v111
	strs_syms_116 = strs_syms_101
	keys_syms_117 = keys_syms_102
	sym_118 = sym_103
	tem__G__0_119 = tem__G__0_104
	goto b18
b17:
	;
	as_sym_115 = vm.NIL
	strs_syms_116 = strs_syms_105
	keys_syms_117 = keys_syms_106
	sym_118 = sym_107
	tem__G__0_119 = tem__G__0_108
	goto b18
b18:
	;
	arg__591_121, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__604_132, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "dissoc").Deref(), []vm.Value{sym_118, vm.Keyword("keys"), vm.Keyword("strs"), vm.Keyword("as"), vm.Keyword("or")})
	if callErr != nil {
		return nil, callErr
	}
	arg__617_143, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "dissoc").Deref(), []vm.Value{sym_118, vm.Keyword("keys"), vm.Keyword("strs"), vm.Keyword("as"), vm.Keyword("or")})
	if callErr != nil {
		return nil, callErr
	}
	arg__618_144, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "keys").Deref(), []vm.Value{arg__617_143})
	if callErr != nil {
		return nil, callErr
	}
	arg__632_156, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "dissoc").Deref(), []vm.Value{sym_118, vm.Keyword("keys"), vm.Keyword("strs"), vm.Keyword("as"), vm.Keyword("or")})
	if callErr != nil {
		return nil, callErr
	}
	arg__645_167, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "dissoc").Deref(), []vm.Value{sym_118, vm.Keyword("keys"), vm.Keyword("strs"), vm.Keyword("as"), vm.Keyword("or")})
	if callErr != nil {
		return nil, callErr
	}
	arg__646_168, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "keys").Deref(), []vm.Value{arg__645_167})
	if callErr != nil {
		return nil, callErr
	}
	arg__647_169, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "filter").Deref(), []vm.Value{rt.LookupVar("clojure.core", "symbol?").Deref(), arg__646_168})
	if callErr != nil {
		return nil, callErr
	}
	arg__650_172, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__663_183, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "dissoc").Deref(), []vm.Value{sym_118, vm.Keyword("keys"), vm.Keyword("strs"), vm.Keyword("as"), vm.Keyword("or")})
	if callErr != nil {
		return nil, callErr
	}
	arg__676_194, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "dissoc").Deref(), []vm.Value{sym_118, vm.Keyword("keys"), vm.Keyword("strs"), vm.Keyword("as"), vm.Keyword("or")})
	if callErr != nil {
		return nil, callErr
	}
	arg__677_195, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "keys").Deref(), []vm.Value{arg__676_194})
	if callErr != nil {
		return nil, callErr
	}
	arg__691_207, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "dissoc").Deref(), []vm.Value{sym_118, vm.Keyword("keys"), vm.Keyword("strs"), vm.Keyword("as"), vm.Keyword("or")})
	if callErr != nil {
		return nil, callErr
	}
	arg__704_218, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "dissoc").Deref(), []vm.Value{sym_118, vm.Keyword("keys"), vm.Keyword("strs"), vm.Keyword("as"), vm.Keyword("or")})
	if callErr != nil {
		return nil, callErr
	}
	arg__705_219, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "keys").Deref(), []vm.Value{arg__704_218})
	if callErr != nil {
		return nil, callErr
	}
	arg__706_220, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "filter").Deref(), []vm.Value{rt.LookupVar("clojure.core", "symbol?").Deref(), arg__705_219})
	if callErr != nil {
		return nil, callErr
	}
	other_syms_221, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__650_172, arg__706_220})
	if callErr != nil {
		return nil, callErr
	}
	arg__708_223, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__718_225, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "concat").Deref(), []vm.Value{keys_syms_117, strs_syms_116, as_sym_115, other_syms_221})
	if callErr != nil {
		return nil, callErr
	}
	arg__721_228, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__731_230, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "concat").Deref(), []vm.Value{keys_syms_117, strs_syms_116, as_sym_115, other_syms_221})
	if callErr != nil {
		return nil, callErr
	}
	v231, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__721_228, arg__731_230})
	if callErr != nil {
		return nil, callErr
	}
	v245 = v231
	sym_246 = sym_118
	goto b9
b19:
	;
	v238, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	v242 = v238
	sym_243 = sym_233
	goto b21
b20:
	;
	v242 = vm.NIL
	sym_243 = sym_234
	goto b21
b21:
	;
	v245 = v242
	sym_246 = sym_243
	goto b9
}
func build_args(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var pre_locals_3 vm.Value
	var v5 vm.Value
	var syms_9 vm.Value
	var results_13 vm.Value
	var doseq_seq__733_15 vm.Value
	var doseq_loop__734_16 vm.Value
	var ctx_17 vm.Value
	var syms_18 vm.Value
	var results_19 vm.Value
	var v311 vm.Value
	var v317 string
	var forms_21 vm.Value
	var pre_locals_22 vm.Value
	var doseq_seq__733_23 vm.Value
	var doseq_loop__734_24 vm.Value
	var ctx_25 vm.Value
	var syms_26 vm.Value
	var results_27 vm.Value
	var v310 vm.Value
	var v316 string
	var a_37 vm.Value
	var r_39 vm.Value
	var v59 vm.Value
	var forms_28 vm.Value
	var pre_locals_29 vm.Value
	var doseq_seq__733_30 vm.Value
	var doseq_loop__734_31 vm.Value
	var ctx_32 vm.Value
	var syms_33 vm.Value
	var results_34 vm.Value
	var v315 vm.Value
	var v321 string
	var v102 vm.Value
	var forms_103 vm.Value
	var pre_locals_104 vm.Value
	var doseq_seq__733_105 vm.Value
	var doseq_loop__734_106 vm.Value
	var ctx_107 vm.Value
	var syms_108 vm.Value
	var results_109 vm.Value
	var arg__814_114 vm.Value
	var arg__818_116 vm.Value
	var arg__829_122 vm.Value
	var arg__833_124 vm.Value
	var arg__834_125 vm.Value
	var arg__845_131 vm.Value
	var arg__849_133 vm.Value
	var arg__860_139 vm.Value
	var arg__864_141 vm.Value
	var arg__865_142 vm.Value
	var threaded_143 vm.Value
	var post_locals_145 vm.Value
	var v147 vm.Value
	var doseq_seq__735_149 vm.Value
	var forms_40 vm.Value
	var pre_locals_41 vm.Value
	var doseq_seq__733_42 vm.Value
	var doseq_loop__734_43 vm.Value
	var ctx_44 vm.Value
	var syms_45 vm.Value
	var results_46 vm.Value
	var a_47 vm.Value
	var r_48 vm.Value
	var v313 vm.Value
	var v319 string
	var v66 vm.Value
	var v70 vm.Value
	var forms_49 vm.Value
	var pre_locals_50 vm.Value
	var doseq_seq__733_51 vm.Value
	var doseq_loop__734_52 vm.Value
	var ctx_53 vm.Value
	var syms_54 vm.Value
	var results_55 vm.Value
	var a_56 vm.Value
	var r_57 vm.Value
	var v312 vm.Value
	var v318 string
	var sym_75 vm.Value
	var v77 vm.Value
	var v81 vm.Value
	var v85 vm.Value
	var v87 vm.Value
	var forms_88 vm.Value
	var pre_locals_89 vm.Value
	var doseq_seq__733_90 vm.Value
	var doseq_loop__734_91 vm.Value
	var ctx_92 vm.Value
	var syms_93 vm.Value
	var results_94 vm.Value
	var a_95 vm.Value
	var r_96 vm.Value
	var v314 vm.Value
	var v320 string
	var v98 vm.Value
	var doseq_loop__736_150 vm.Value
	var pre_locals_151 vm.Value
	var ctx_152 vm.Value
	var v333 int
	var v342 vm.Value
	var v351 int
	var forms_154 vm.Value
	var doseq_loop__734_155 vm.Value
	var syms_156 vm.Value
	var results_157 vm.Value
	var threaded_158 vm.Value
	var post_locals_159 vm.Value
	var doseq_seq__735_160 vm.Value
	var doseq_loop__736_161 vm.Value
	var pre_locals_162 vm.Value
	var ctx_163 vm.Value
	var v336 int
	var v345 vm.Value
	var v354 int
	var vec__737_176 vm.Value
	var sym_182 vm.Value
	var val_188 vm.Value
	var and__x_216 vm.Value
	var forms_164 vm.Value
	var doseq_loop__734_165 vm.Value
	var syms_166 vm.Value
	var results_167 vm.Value
	var threaded_168 vm.Value
	var post_locals_169 vm.Value
	var doseq_seq__735_170 vm.Value
	var doseq_loop__736_171 vm.Value
	var pre_locals_172 vm.Value
	var ctx_173 vm.Value
	var v339 int
	var v348 vm.Value
	var v357 int
	var v294 vm.Value
	var forms_295 vm.Value
	var doseq_loop__734_296 vm.Value
	var syms_297 vm.Value
	var results_298 vm.Value
	var threaded_299 vm.Value
	var post_locals_300 vm.Value
	var doseq_seq__735_301 vm.Value
	var doseq_loop__736_302 vm.Value
	var pre_locals_303 vm.Value
	var ctx_304 vm.Value
	var forms_189 vm.Value
	var doseq_loop__734_190 vm.Value
	var syms_191 vm.Value
	var results_192 vm.Value
	var threaded_193 vm.Value
	var post_locals_194 vm.Value
	var doseq_seq__735_195 vm.Value
	var doseq_loop__736_196 vm.Value
	var pre_locals_197 vm.Value
	var ctx_198 vm.Value
	var vec__737_199 vm.Value
	var sym_200 vm.Value
	var val_201 vm.Value
	var v335 int
	var v344 vm.Value
	var v353 int
	var v271 vm.Value
	var forms_202 vm.Value
	var doseq_loop__734_203 vm.Value
	var syms_204 vm.Value
	var results_205 vm.Value
	var threaded_206 vm.Value
	var post_locals_207 vm.Value
	var doseq_seq__735_208 vm.Value
	var doseq_loop__736_209 vm.Value
	var pre_locals_210 vm.Value
	var ctx_211 vm.Value
	var vec__737_212 vm.Value
	var sym_213 vm.Value
	var val_214 vm.Value
	var v332 int
	var v341 vm.Value
	var v350 int
	var v275 vm.Value
	var forms_276 vm.Value
	var doseq_loop__734_277 vm.Value
	var syms_278 vm.Value
	var results_279 vm.Value
	var threaded_280 vm.Value
	var post_locals_281 vm.Value
	var doseq_seq__735_282 vm.Value
	var doseq_loop__736_283 vm.Value
	var pre_locals_284 vm.Value
	var ctx_285 vm.Value
	var vec__737_286 vm.Value
	var sym_287 vm.Value
	var val_288 vm.Value
	var v338 int
	var v347 vm.Value
	var v356 int
	var v290 vm.Value
	var forms_217 vm.Value
	var doseq_loop__734_218 vm.Value
	var syms_219 vm.Value
	var results_220 vm.Value
	var threaded_221 vm.Value
	var post_locals_222 vm.Value
	var doseq_seq__735_223 vm.Value
	var doseq_loop__736_224 vm.Value
	var pre_locals_225 vm.Value
	var ctx_226 vm.Value
	var vec__737_227 vm.Value
	var sym_228 vm.Value
	var val_229 vm.Value
	var and__x_230 vm.Value
	var v334 int
	var v343 vm.Value
	var v352 int
	var arg__903_247 vm.Value
	var arg__911_250 vm.Value
	var v251 vm.Value
	var forms_231 vm.Value
	var doseq_loop__734_232 vm.Value
	var syms_233 vm.Value
	var results_234 vm.Value
	var threaded_235 vm.Value
	var post_locals_236 vm.Value
	var doseq_seq__735_237 vm.Value
	var doseq_loop__736_238 vm.Value
	var pre_locals_239 vm.Value
	var ctx_240 vm.Value
	var vec__737_241 vm.Value
	var sym_242 vm.Value
	var val_243 vm.Value
	var and__x_244 vm.Value
	var v331 int
	var v340 vm.Value
	var v349 int
	var v254 vm.Value
	var forms_255 vm.Value
	var doseq_loop__734_256 vm.Value
	var syms_257 vm.Value
	var results_258 vm.Value
	var threaded_259 vm.Value
	var post_locals_260 vm.Value
	var doseq_seq__735_261 vm.Value
	var doseq_loop__736_262 vm.Value
	var pre_locals_263 vm.Value
	var ctx_264 vm.Value
	var vec__737_265 vm.Value
	var sym_266 vm.Value
	var val_267 vm.Value
	var and__x_268 vm.Value
	var v337 int
	var v346 vm.Value
	var v355 int
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = pre_locals_3, v5, syms_9, results_13, doseq_seq__733_15, doseq_loop__734_16, ctx_17, syms_18, results_19, v311, v317, forms_21, pre_locals_22, doseq_seq__733_23, doseq_loop__734_24, ctx_25, syms_26, results_27, v310, v316, a_37, r_39, v59, forms_28, pre_locals_29, doseq_seq__733_30, doseq_loop__734_31, ctx_32, syms_33, results_34, v315, v321, v102, forms_103, pre_locals_104, doseq_seq__733_105, doseq_loop__734_106, ctx_107, syms_108, results_109, arg__814_114, arg__818_116, arg__829_122, arg__833_124, arg__834_125, arg__845_131, arg__849_133, arg__860_139, arg__864_141, arg__865_142, threaded_143, post_locals_145, v147, doseq_seq__735_149, forms_40, pre_locals_41, doseq_seq__733_42, doseq_loop__734_43, ctx_44, syms_45, results_46, a_47, r_48, v313, v319, v66, v70, forms_49, pre_locals_50, doseq_seq__733_51, doseq_loop__734_52, ctx_53, syms_54, results_55, a_56, r_57, v312, v318, sym_75, v77, v81, v85, v87, forms_88, pre_locals_89, doseq_seq__733_90, doseq_loop__734_91, ctx_92, syms_93, results_94, a_95, r_96, v314, v320, v98, doseq_loop__736_150, pre_locals_151, ctx_152, v333, v342, v351, forms_154, doseq_loop__734_155, syms_156, results_157, threaded_158, post_locals_159, doseq_seq__735_160, doseq_loop__736_161, pre_locals_162, ctx_163, v336, v345, v354, vec__737_176, sym_182, val_188, and__x_216, forms_164, doseq_loop__734_165, syms_166, results_167, threaded_168, post_locals_169, doseq_seq__735_170, doseq_loop__736_171, pre_locals_172, ctx_173, v339, v348, v357, v294, forms_295, doseq_loop__734_296, syms_297, results_298, threaded_299, post_locals_300, doseq_seq__735_301, doseq_loop__736_302, pre_locals_303, ctx_304, forms_189, doseq_loop__734_190, syms_191, results_192, threaded_193, post_locals_194, doseq_seq__735_195, doseq_loop__736_196, pre_locals_197, ctx_198, vec__737_199, sym_200, val_201, v335, v344, v353, v271, forms_202, doseq_loop__734_203, syms_204, results_205, threaded_206, post_locals_207, doseq_seq__735_208, doseq_loop__736_209, pre_locals_210, ctx_211, vec__737_212, sym_213, val_214, v332, v341, v350, v275, forms_276, doseq_loop__734_277, syms_278, results_279, threaded_280, post_locals_281, doseq_seq__735_282, doseq_loop__736_283, pre_locals_284, ctx_285, vec__737_286, sym_287, val_288, v338, v347, v356, v290, forms_217, doseq_loop__734_218, syms_219, results_220, threaded_221, post_locals_222, doseq_seq__735_223, doseq_loop__736_224, pre_locals_225, ctx_226, vec__737_227, sym_228, val_229, and__x_230, v334, v343, v352, arg__903_247, arg__911_250, v251, forms_231, doseq_loop__734_232, syms_233, results_234, threaded_235, post_locals_236, doseq_seq__735_237, doseq_loop__736_238, pre_locals_239, ctx_240, vec__737_241, sym_242, val_243, and__x_244, v331, v340, v349, v254, forms_255, doseq_loop__734_256, syms_257, results_258, threaded_259, post_locals_260, doseq_seq__735_261, doseq_loop__736_262, pre_locals_263, ctx_264, vec__737_265, sym_266, val_267, and__x_268, v337, v346, v355
	pre_locals_3, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "current-locals-flat").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	v5, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "push-locals!").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	syms_9, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "atom").Deref(), []vm.Value{vm.NewArrayVector([]vm.Value{})})
	if callErr != nil {
		return nil, callErr
	}
	results_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "atom").Deref(), []vm.Value{vm.NewArrayVector([]vm.Value{})})
	if callErr != nil {
		return nil, callErr
	}
	doseq_seq__733_15, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__734_16 = doseq_seq__733_15
	ctx_17 = arg1
	syms_18 = syms_9
	results_19 = results_13
	v311 = vm.NIL
	v317 = "arg__"
	goto b1
b1:
	;
	if vm.IsTruthy(doseq_loop__734_16) {
		forms_21 = arg0
		pre_locals_22 = pre_locals_3
		doseq_seq__733_23 = doseq_seq__733_15
		doseq_loop__734_24 = doseq_loop__734_16
		ctx_25 = ctx_17
		syms_26 = syms_18
		results_27 = results_19
		v310 = v311
		v316 = v317
		goto b2
	} else {
		forms_28 = arg0
		pre_locals_29 = pre_locals_3
		doseq_seq__733_30 = doseq_seq__733_15
		doseq_loop__734_31 = doseq_loop__734_16
		ctx_32 = ctx_17
		syms_33 = syms_18
		results_34 = results_19
		v315 = v311
		v321 = v317
		goto b3
	}
b2:
	;
	a_37, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{doseq_loop__734_24})
	if callErr != nil {
		return nil, callErr
	}
	r_39, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-form").Deref(), []vm.Value{a_37, ctx_25})
	if callErr != nil {
		return nil, callErr
	}
	v59, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "terminated?").Deref(), []vm.Value{r_39})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v59) {
		forms_40 = forms_21
		pre_locals_41 = pre_locals_22
		doseq_seq__733_42 = doseq_seq__733_23
		doseq_loop__734_43 = doseq_loop__734_24
		ctx_44 = ctx_25
		syms_45 = syms_26
		results_46 = results_27
		a_47 = a_37
		r_48 = r_39
		v313 = v310
		v319 = v316
		goto b5
	} else {
		forms_49 = forms_21
		pre_locals_50 = pre_locals_22
		doseq_seq__733_51 = doseq_seq__733_23
		doseq_loop__734_52 = doseq_loop__734_24
		ctx_53 = ctx_25
		syms_54 = syms_26
		results_55 = results_27
		a_56 = a_37
		r_57 = r_39
		v312 = v310
		v318 = v316
		goto b6
	}
b3:
	;
	v102 = vm.NIL
	forms_103 = forms_28
	pre_locals_104 = pre_locals_29
	doseq_seq__733_105 = doseq_seq__733_30
	doseq_loop__734_106 = doseq_loop__734_31
	ctx_107 = ctx_32
	syms_108 = syms_33
	results_109 = results_34
	goto b4
b4:
	;
	arg__814_114, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{syms_108})
	if callErr != nil {
		return nil, callErr
	}
	arg__818_116, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{results_109})
	if callErr != nil {
		return nil, callErr
	}
	arg__829_122, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{syms_108})
	if callErr != nil {
		return nil, callErr
	}
	arg__833_124, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{results_109})
	if callErr != nil {
		return nil, callErr
	}
	arg__834_125, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
		var sym_3 vm.Value
		var orig_4 vm.Value
		var ctx_5 vm.Value
		var or__x_11 vm.Value
		var sym_6 vm.Value
		var orig_7 vm.Value
		var ctx_8 vm.Value
		var v30 vm.Value
		var sym_31 vm.Value
		var orig_32 vm.Value
		var ctx_33 vm.Value
		var sym_12 vm.Value
		var orig_13 vm.Value
		var ctx_14 vm.Value
		var or__x_15 vm.Value
		var sym_16 vm.Value
		var orig_17 vm.Value
		var ctx_18 vm.Value
		var or__x_19 vm.Value
		var v23 vm.Value
		var sym_24 vm.Value
		var orig_25 vm.Value
		var ctx_26 vm.Value
		var or__x_27 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = sym_3, orig_4, ctx_5, or__x_11, sym_6, orig_7, ctx_8, v30, sym_31, orig_32, ctx_33, sym_12, orig_13, ctx_14, or__x_15, sym_16, orig_17, ctx_18, or__x_19, v23, sym_24, orig_25, ctx_26, or__x_27
		if vm.IsTruthy(arg0) {
			sym_3 = arg0
			orig_4 = arg1
			ctx_5 = ctx_107
			goto b1
		} else {
			sym_6 = arg0
			orig_7 = arg1
			ctx_8 = ctx_107
			goto b2
		}
	b1:
		;
		or__x_11, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "lookup-local").Deref(), []vm.Value{ctx_5, sym_3})
		if callErr != nil {
			return nil, callErr
		}
		if vm.IsTruthy(or__x_11) {
			sym_12 = sym_3
			orig_13 = orig_4
			ctx_14 = ctx_5
			or__x_15 = or__x_11
			goto b4
		} else {
			sym_16 = sym_3
			orig_17 = orig_4
			ctx_18 = ctx_5
			or__x_19 = or__x_11
			goto b5
		}
	b2:
		;
		v30 = orig_7
		sym_31 = sym_6
		orig_32 = orig_7
		ctx_33 = ctx_8
		goto b3
	b3:
		;
		return v30, nil
	b4:
		;
		v23 = or__x_15
		sym_24 = sym_12
		orig_25 = orig_13
		ctx_26 = ctx_14
		or__x_27 = or__x_15
		goto b6
	b5:
		;
		v23 = orig_17
		sym_24 = sym_16
		orig_25 = orig_17
		ctx_26 = ctx_18
		or__x_27 = or__x_19
		goto b6
	b6:
		;
		v30 = v23
		sym_31 = sym_24
		orig_32 = orig_25
		ctx_33 = ctx_26
		goto b3
	}), arg__829_122, arg__833_124})
	if callErr != nil {
		return nil, callErr
	}
	arg__845_131, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{syms_108})
	if callErr != nil {
		return nil, callErr
	}
	arg__849_133, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{results_109})
	if callErr != nil {
		return nil, callErr
	}
	arg__860_139, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{syms_108})
	if callErr != nil {
		return nil, callErr
	}
	arg__864_141, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{results_109})
	if callErr != nil {
		return nil, callErr
	}
	arg__865_142, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
		var sym_3 vm.Value
		var orig_4 vm.Value
		var ctx_5 vm.Value
		var or__x_11 vm.Value
		var sym_6 vm.Value
		var orig_7 vm.Value
		var ctx_8 vm.Value
		var v30 vm.Value
		var sym_31 vm.Value
		var orig_32 vm.Value
		var ctx_33 vm.Value
		var sym_12 vm.Value
		var orig_13 vm.Value
		var ctx_14 vm.Value
		var or__x_15 vm.Value
		var sym_16 vm.Value
		var orig_17 vm.Value
		var ctx_18 vm.Value
		var or__x_19 vm.Value
		var v23 vm.Value
		var sym_24 vm.Value
		var orig_25 vm.Value
		var ctx_26 vm.Value
		var or__x_27 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = sym_3, orig_4, ctx_5, or__x_11, sym_6, orig_7, ctx_8, v30, sym_31, orig_32, ctx_33, sym_12, orig_13, ctx_14, or__x_15, sym_16, orig_17, ctx_18, or__x_19, v23, sym_24, orig_25, ctx_26, or__x_27
		if vm.IsTruthy(arg0) {
			sym_3 = arg0
			orig_4 = arg1
			ctx_5 = ctx_107
			goto b1
		} else {
			sym_6 = arg0
			orig_7 = arg1
			ctx_8 = ctx_107
			goto b2
		}
	b1:
		;
		or__x_11, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "lookup-local").Deref(), []vm.Value{ctx_5, sym_3})
		if callErr != nil {
			return nil, callErr
		}
		if vm.IsTruthy(or__x_11) {
			sym_12 = sym_3
			orig_13 = orig_4
			ctx_14 = ctx_5
			or__x_15 = or__x_11
			goto b4
		} else {
			sym_16 = sym_3
			orig_17 = orig_4
			ctx_18 = ctx_5
			or__x_19 = or__x_11
			goto b5
		}
	b2:
		;
		v30 = orig_7
		sym_31 = sym_6
		orig_32 = orig_7
		ctx_33 = ctx_8
		goto b3
	b3:
		;
		return v30, nil
	b4:
		;
		v23 = or__x_15
		sym_24 = sym_12
		orig_25 = orig_13
		ctx_26 = ctx_14
		or__x_27 = or__x_15
		goto b6
	b5:
		;
		v23 = orig_17
		sym_24 = sym_16
		orig_25 = orig_17
		ctx_26 = ctx_18
		or__x_27 = or__x_19
		goto b6
	b6:
		;
		v30 = v23
		sym_31 = sym_24
		orig_32 = orig_25
		ctx_33 = ctx_26
		goto b3
	}), arg__860_139, arg__864_141})
	if callErr != nil {
		return nil, callErr
	}
	threaded_143, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg__865_142})
	if callErr != nil {
		return nil, callErr
	}
	post_locals_145, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "current-locals-flat").Deref(), []vm.Value{ctx_107})
	if callErr != nil {
		return nil, callErr
	}
	v147, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "pop-locals!").Deref(), []vm.Value{ctx_107})
	if callErr != nil {
		return nil, callErr
	}
	doseq_seq__735_149, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{post_locals_145})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__736_150 = doseq_seq__735_149
	pre_locals_151 = pre_locals_104
	ctx_152 = ctx_107
	v333 = 0
	v342 = vm.NIL
	v351 = 1
	goto b8
b5:
	;
	v66, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{syms_45, rt.LookupVar("clojure.core", "conj").Deref(), v313})
	if callErr != nil {
		return nil, callErr
	}
	v70, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{results_46, rt.LookupVar("clojure.core", "conj").Deref(), r_48})
	if callErr != nil {
		return nil, callErr
	}
	v87 = v70
	forms_88 = forms_40
	pre_locals_89 = pre_locals_41
	doseq_seq__733_90 = doseq_seq__733_42
	doseq_loop__734_91 = doseq_loop__734_43
	ctx_92 = ctx_44
	syms_93 = syms_45
	results_94 = results_46
	a_95 = a_47
	r_96 = r_48
	v314 = v313
	v320 = v319
	goto b7
b6:
	;
	sym_75, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "gensym").Deref(), []vm.Value{vm.String(v318)})
	if callErr != nil {
		return nil, callErr
	}
	v77, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "bind-local!").Deref(), []vm.Value{ctx_53, sym_75, r_57})
	if callErr != nil {
		return nil, callErr
	}
	v81, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{syms_54, rt.LookupVar("clojure.core", "conj").Deref(), sym_75})
	if callErr != nil {
		return nil, callErr
	}
	v85, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{results_55, rt.LookupVar("clojure.core", "conj").Deref(), r_57})
	if callErr != nil {
		return nil, callErr
	}
	v87 = v85
	forms_88 = forms_49
	pre_locals_89 = pre_locals_50
	doseq_seq__733_90 = doseq_seq__733_51
	doseq_loop__734_91 = doseq_loop__734_52
	ctx_92 = ctx_53
	syms_93 = syms_54
	results_94 = results_55
	a_95 = a_56
	r_96 = r_57
	v314 = v312
	v320 = v318
	goto b7
b7:
	;
	v98, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "next").Deref(), []vm.Value{doseq_loop__734_91})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__734_16 = v98
	ctx_17 = ctx_92
	syms_18 = syms_93
	results_19 = results_94
	v311 = v314
	v317 = v320
	goto b1
b8:
	;
	if vm.IsTruthy(doseq_loop__736_150) {
		forms_154 = forms_103
		doseq_loop__734_155 = doseq_loop__734_106
		syms_156 = syms_108
		results_157 = results_109
		threaded_158 = threaded_143
		post_locals_159 = post_locals_145
		doseq_seq__735_160 = doseq_seq__735_149
		doseq_loop__736_161 = doseq_loop__736_150
		pre_locals_162 = pre_locals_151
		ctx_163 = ctx_152
		v336 = v333
		v345 = v342
		v354 = v351
		goto b9
	} else {
		forms_164 = forms_103
		doseq_loop__734_165 = doseq_loop__734_106
		syms_166 = syms_108
		results_167 = results_109
		threaded_168 = threaded_143
		post_locals_169 = post_locals_145
		doseq_seq__735_170 = doseq_seq__735_149
		doseq_loop__736_171 = doseq_loop__736_150
		pre_locals_172 = pre_locals_151
		ctx_173 = ctx_152
		v339 = v333
		v348 = v342
		v357 = v351
		goto b10
	}
b9:
	;
	vec__737_176, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{doseq_loop__736_161})
	if callErr != nil {
		return nil, callErr
	}
	sym_182, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{vec__737_176, vm.Int(v336), v345})
	if callErr != nil {
		return nil, callErr
	}
	val_188, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{vec__737_176, vm.Int(v354), v345})
	if callErr != nil {
		return nil, callErr
	}
	and__x_216, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "contains?").Deref(), []vm.Value{pre_locals_162, sym_182})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(and__x_216) {
		forms_217 = forms_154
		doseq_loop__734_218 = doseq_loop__734_155
		syms_219 = syms_156
		results_220 = results_157
		threaded_221 = threaded_158
		post_locals_222 = post_locals_159
		doseq_seq__735_223 = doseq_seq__735_160
		doseq_loop__736_224 = doseq_loop__736_161
		pre_locals_225 = pre_locals_162
		ctx_226 = ctx_163
		vec__737_227 = vec__737_176
		sym_228 = sym_182
		val_229 = val_188
		and__x_230 = and__x_216
		v334 = v336
		v343 = v345
		v352 = v354
		goto b15
	} else {
		forms_231 = forms_154
		doseq_loop__734_232 = doseq_loop__734_155
		syms_233 = syms_156
		results_234 = results_157
		threaded_235 = threaded_158
		post_locals_236 = post_locals_159
		doseq_seq__735_237 = doseq_seq__735_160
		doseq_loop__736_238 = doseq_loop__736_161
		pre_locals_239 = pre_locals_162
		ctx_240 = ctx_163
		vec__737_241 = vec__737_176
		sym_242 = sym_182
		val_243 = val_188
		and__x_244 = and__x_216
		v331 = v336
		v340 = v345
		v349 = v354
		goto b16
	}
b10:
	;
	v294 = vm.NIL
	forms_295 = forms_164
	doseq_loop__734_296 = doseq_loop__734_165
	syms_297 = syms_166
	results_298 = results_167
	threaded_299 = threaded_168
	post_locals_300 = post_locals_169
	doseq_seq__735_301 = doseq_seq__735_170
	doseq_loop__736_302 = doseq_loop__736_171
	pre_locals_303 = pre_locals_172
	ctx_304 = ctx_173
	goto b11
b11:
	;
	return threaded_299, nil
b12:
	;
	v271, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "bind-local!").Deref(), []vm.Value{ctx_198, sym_200, val_201})
	if callErr != nil {
		return nil, callErr
	}
	v275 = v271
	forms_276 = forms_189
	doseq_loop__734_277 = doseq_loop__734_190
	syms_278 = syms_191
	results_279 = results_192
	threaded_280 = threaded_193
	post_locals_281 = post_locals_194
	doseq_seq__735_282 = doseq_seq__735_195
	doseq_loop__736_283 = doseq_loop__736_196
	pre_locals_284 = pre_locals_197
	ctx_285 = ctx_198
	vec__737_286 = vec__737_199
	sym_287 = sym_200
	val_288 = val_201
	v338 = v335
	v347 = v344
	v356 = v353
	goto b14
b13:
	;
	v275 = vm.NIL
	forms_276 = forms_202
	doseq_loop__734_277 = doseq_loop__734_203
	syms_278 = syms_204
	results_279 = results_205
	threaded_280 = threaded_206
	post_locals_281 = post_locals_207
	doseq_seq__735_282 = doseq_seq__735_208
	doseq_loop__736_283 = doseq_loop__736_209
	pre_locals_284 = pre_locals_210
	ctx_285 = ctx_211
	vec__737_286 = vec__737_212
	sym_287 = sym_213
	val_288 = val_214
	v338 = v332
	v347 = v341
	v356 = v350
	goto b14
b14:
	;
	v290, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "next").Deref(), []vm.Value{doseq_loop__736_283})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__736_150 = v290
	pre_locals_151 = pre_locals_284
	ctx_152 = ctx_285
	v333 = v338
	v342 = v347
	v351 = v356
	goto b8
b15:
	;
	arg__903_247, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{pre_locals_225, sym_228})
	if callErr != nil {
		return nil, callErr
	}
	arg__911_250, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{pre_locals_225, sym_228})
	if callErr != nil {
		return nil, callErr
	}
	v251, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not=").Deref(), []vm.Value{val_229, arg__911_250})
	if callErr != nil {
		return nil, callErr
	}
	v254 = v251
	forms_255 = forms_217
	doseq_loop__734_256 = doseq_loop__734_218
	syms_257 = syms_219
	results_258 = results_220
	threaded_259 = threaded_221
	post_locals_260 = post_locals_222
	doseq_seq__735_261 = doseq_seq__735_223
	doseq_loop__736_262 = doseq_loop__736_224
	pre_locals_263 = pre_locals_225
	ctx_264 = ctx_226
	vec__737_265 = vec__737_227
	sym_266 = sym_228
	val_267 = val_229
	and__x_268 = and__x_230
	v337 = v334
	v346 = v343
	v355 = v352
	goto b17
b16:
	;
	v254 = and__x_244
	forms_255 = forms_231
	doseq_loop__734_256 = doseq_loop__734_232
	syms_257 = syms_233
	results_258 = results_234
	threaded_259 = threaded_235
	post_locals_260 = post_locals_236
	doseq_seq__735_261 = doseq_seq__735_237
	doseq_loop__736_262 = doseq_loop__736_238
	pre_locals_263 = pre_locals_239
	ctx_264 = ctx_240
	vec__737_265 = vec__737_241
	sym_266 = sym_242
	val_267 = val_243
	and__x_268 = and__x_244
	v337 = v331
	v346 = v340
	v355 = v349
	goto b17
b17:
	;
	if vm.IsTruthy(v254) {
		forms_189 = forms_255
		doseq_loop__734_190 = doseq_loop__734_256
		syms_191 = syms_257
		results_192 = results_258
		threaded_193 = threaded_259
		post_locals_194 = post_locals_260
		doseq_seq__735_195 = doseq_seq__735_261
		doseq_loop__736_196 = doseq_loop__736_262
		pre_locals_197 = pre_locals_263
		ctx_198 = ctx_264
		vec__737_199 = vec__737_265
		sym_200 = sym_266
		val_201 = val_267
		v335 = v337
		v344 = v346
		v353 = v355
		goto b12
	} else {
		forms_202 = forms_255
		doseq_loop__734_203 = doseq_loop__734_256
		syms_204 = syms_257
		results_205 = results_258
		threaded_206 = threaded_259
		post_locals_207 = post_locals_260
		doseq_seq__735_208 = doseq_seq__735_261
		doseq_loop__736_209 = doseq_loop__736_262
		pre_locals_210 = pre_locals_263
		ctx_211 = ctx_264
		vec__737_212 = vec__737_265
		sym_213 = sym_266
		val_214 = val_267
		v332 = v337
		v341 = v346
		v350 = v355
		goto b13
	}
}
func build_builtin_op(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__925_4 vm.Value
	var arg__931_7 vm.Value
	var args_8 vm.Value
	var v20 vm.Value
	var op_kw_9 vm.Value
	var form_10 vm.Value
	var ctx_11 vm.Value
	var args_12 vm.Value
	var arg__942_32 vm.Value
	var v33 bool
	var op_kw_13 vm.Value
	var form_14 vm.Value
	var ctx_15 vm.Value
	var args_16 vm.Value
	var arg__989_79 vm.Value
	var v80 bool
	var v194 vm.Value
	var op_kw_195 vm.Value
	var form_196 vm.Value
	var ctx_197 vm.Value
	var args_198 vm.Value
	var op_kw_22 vm.Value
	var form_23 vm.Value
	var ctx_24 vm.Value
	var args_25 vm.Value
	var arg__946_36 vm.Value
	var arg__950_38 vm.Value
	var arg__958_43 vm.Value
	var arg__959_44 vm.Value
	var arg__965_48 vm.Value
	var arg__969_50 vm.Value
	var arg__977_55 vm.Value
	var arg__978_56 vm.Value
	var v58 vm.Value
	var op_kw_26 vm.Value
	var form_27 vm.Value
	var ctx_28 vm.Value
	var args_29 vm.Value
	var v61 vm.Value
	var v63 vm.Value
	var op_kw_64 vm.Value
	var form_65 vm.Value
	var ctx_66 vm.Value
	var args_67 vm.Value
	var op_kw_69 vm.Value
	var form_70 vm.Value
	var ctx_71 vm.Value
	var args_72 vm.Value
	var v91 bool
	var op_kw_73 vm.Value
	var form_74 vm.Value
	var ctx_75 vm.Value
	var args_76 vm.Value
	var arg__1051_149 vm.Value
	var v150 bool
	var v188 vm.Value
	var op_kw_189 vm.Value
	var form_190 vm.Value
	var ctx_191 vm.Value
	var args_192 vm.Value
	var op_kw_82 vm.Value
	var form_83 vm.Value
	var ctx_84 vm.Value
	var args_85 vm.Value
	var arg__996_94 vm.Value
	var arg__1005_100 vm.Value
	var zero_id_104 vm.Value
	var arg__1013_106 vm.Value
	var arg__1022_112 vm.Value
	var arg__1023_113 vm.Value
	var arg__1030_117 vm.Value
	var arg__1039_123 vm.Value
	var arg__1040_124 vm.Value
	var v126 vm.Value
	var op_kw_86 vm.Value
	var form_87 vm.Value
	var ctx_88 vm.Value
	var args_89 vm.Value
	var v131 vm.Value
	var v133 vm.Value
	var op_kw_134 vm.Value
	var form_135 vm.Value
	var ctx_136 vm.Value
	var args_137 vm.Value
	var op_kw_139 vm.Value
	var form_140 vm.Value
	var ctx_141 vm.Value
	var args_142 vm.Value
	var arg__1056_153 vm.Value
	var arg__1065_157 vm.Value
	var v159 vm.Value
	var op_kw_143 vm.Value
	var form_144 vm.Value
	var ctx_145 vm.Value
	var args_146 vm.Value
	var v182 vm.Value
	var op_kw_183 vm.Value
	var form_184 vm.Value
	var ctx_185 vm.Value
	var args_186 vm.Value
	var op_kw_161 vm.Value
	var form_162 vm.Value
	var ctx_163 vm.Value
	var args_164 vm.Value
	var v172 vm.Value
	var op_kw_165 vm.Value
	var form_166 vm.Value
	var ctx_167 vm.Value
	var args_168 vm.Value
	var v176 vm.Value
	var op_kw_177 vm.Value
	var form_178 vm.Value
	var ctx_179 vm.Value
	var args_180 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__925_4, arg__931_7, args_8, v20, op_kw_9, form_10, ctx_11, args_12, arg__942_32, v33, op_kw_13, form_14, ctx_15, args_16, arg__989_79, v80, v194, op_kw_195, form_196, ctx_197, args_198, op_kw_22, form_23, ctx_24, args_25, arg__946_36, arg__950_38, arg__958_43, arg__959_44, arg__965_48, arg__969_50, arg__977_55, arg__978_56, v58, op_kw_26, form_27, ctx_28, args_29, v61, v63, op_kw_64, form_65, ctx_66, args_67, op_kw_69, form_70, ctx_71, args_72, v91, op_kw_73, form_74, ctx_75, args_76, arg__1051_149, v150, v188, op_kw_189, form_190, ctx_191, args_192, op_kw_82, form_83, ctx_84, args_85, arg__996_94, arg__1005_100, zero_id_104, arg__1013_106, arg__1022_112, arg__1023_113, arg__1030_117, arg__1039_123, arg__1040_124, v126, op_kw_86, form_87, ctx_88, args_89, v131, v133, op_kw_134, form_135, ctx_136, args_137, op_kw_139, form_140, ctx_141, args_142, arg__1056_153, arg__1065_157, v159, op_kw_143, form_144, ctx_145, args_146, v182, op_kw_183, form_184, ctx_185, args_186, op_kw_161, form_162, ctx_163, args_164, v172, op_kw_165, form_166, ctx_167, args_168, v176, op_kw_177, form_178, ctx_179, args_180
	arg__925_4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__931_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	args_8, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-args").Deref(), []vm.Value{arg__931_7, arg2})
	if callErr != nil {
		return nil, callErr
	}
	v20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "contains?").Deref(), []vm.Value{rt.LookupVar("ir.build", "unary-only-ops").Deref(), arg0})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v20) {
		op_kw_9 = arg0
		form_10 = arg1
		ctx_11 = arg2
		args_12 = args_8
		goto b1
	} else {
		op_kw_13 = arg0
		form_14 = arg1
		ctx_15 = arg2
		args_16 = args_8
		goto b2
	}
b1:
	;
	arg__942_32, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{args_12})
	if callErr != nil {
		return nil, callErr
	}
	v33 = arg__942_32 == vm.Int(1)
	if v33 {
		op_kw_22 = op_kw_9
		form_23 = form_10
		ctx_24 = ctx_11
		args_25 = args_12
		goto b4
	} else {
		op_kw_26 = op_kw_9
		form_27 = form_10
		ctx_28 = ctx_11
		args_29 = args_12
		goto b5
	}
b2:
	;
	arg__989_79, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{args_16})
	if callErr != nil {
		return nil, callErr
	}
	v80 = arg__989_79 == vm.Int(1)
	if v80 {
		op_kw_69 = op_kw_13
		form_70 = form_14
		ctx_71 = ctx_15
		args_72 = args_16
		goto b7
	} else {
		op_kw_73 = op_kw_13
		form_74 = form_14
		ctx_75 = ctx_15
		args_76 = args_16
		goto b8
	}
b3:
	;
	return v194, nil
b4:
	;
	arg__946_36, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-fn").Deref(), []vm.Value{ctx_24})
	if callErr != nil {
		return nil, callErr
	}
	arg__950_38, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_24})
	if callErr != nil {
		return nil, callErr
	}
	arg__958_43, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{args_25, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	arg__959_44, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__958_43})
	if callErr != nil {
		return nil, callErr
	}
	arg__965_48, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-fn").Deref(), []vm.Value{ctx_24})
	if callErr != nil {
		return nil, callErr
	}
	arg__969_50, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_24})
	if callErr != nil {
		return nil, callErr
	}
	arg__977_55, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{args_25, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	arg__978_56, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__977_55})
	if callErr != nil {
		return nil, callErr
	}
	v58, callErr = rt.InvokeValue(rt.LookupVar("ir", "add-inst").Deref(), []vm.Value{arg__965_48, arg__969_50, op_kw_22, arg__978_56, vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	v63 = v58
	op_kw_64 = op_kw_22
	form_65 = form_23
	ctx_66 = ctx_24
	args_67 = args_25
	goto b6
b5:
	;
	v61, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-call").Deref(), []vm.Value{form_27, ctx_28})
	if callErr != nil {
		return nil, callErr
	}
	v63 = v61
	op_kw_64 = op_kw_26
	form_65 = form_27
	ctx_66 = ctx_28
	args_67 = args_29
	goto b6
b6:
	;
	v194 = v63
	op_kw_195 = op_kw_64
	form_196 = form_65
	ctx_197 = ctx_66
	args_198 = args_67
	goto b3
b7:
	;
	v91 = op_kw_69 == vm.Keyword("sub")
	if v91 {
		op_kw_82 = op_kw_69
		form_83 = form_70
		ctx_84 = ctx_71
		args_85 = args_72
		goto b10
	} else {
		op_kw_86 = op_kw_69
		form_87 = form_70
		ctx_88 = ctx_71
		args_89 = args_72
		goto b11
	}
b8:
	;
	arg__1051_149, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{args_76})
	if callErr != nil {
		return nil, callErr
	}
	v150 = arg__1051_149 == vm.Int(2)
	if v150 {
		op_kw_139 = op_kw_73
		form_140 = form_74
		ctx_141 = ctx_75
		args_142 = args_76
		goto b13
	} else {
		op_kw_143 = op_kw_73
		form_144 = form_74
		ctx_145 = ctx_75
		args_146 = args_76
		goto b14
	}
b9:
	;
	v194 = v188
	op_kw_195 = op_kw_189
	form_196 = form_190
	ctx_197 = ctx_191
	args_198 = args_192
	goto b3
b10:
	;
	arg__996_94, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_84})
	if callErr != nil {
		return nil, callErr
	}
	arg__1005_100, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_84})
	if callErr != nil {
		return nil, callErr
	}
	zero_id_104, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "add-inst!").Deref(), []vm.Value{ctx_84, arg__1005_100, vm.Keyword("const"), vm.NewArrayVector([]vm.Value{}), vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	arg__1013_106, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_84})
	if callErr != nil {
		return nil, callErr
	}
	arg__1022_112, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{args_85, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	arg__1023_113, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{zero_id_104, arg__1022_112})
	if callErr != nil {
		return nil, callErr
	}
	arg__1030_117, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_84})
	if callErr != nil {
		return nil, callErr
	}
	arg__1039_123, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{args_85, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	arg__1040_124, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{zero_id_104, arg__1039_123})
	if callErr != nil {
		return nil, callErr
	}
	v126, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "add-inst!").Deref(), []vm.Value{ctx_84, arg__1030_117, vm.Keyword("sub"), arg__1040_124, vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	v133 = v126
	op_kw_134 = op_kw_82
	form_135 = form_83
	ctx_136 = ctx_84
	args_137 = args_85
	goto b12
b11:
	;
	v131, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{args_89, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	v133 = v131
	op_kw_134 = op_kw_86
	form_135 = form_87
	ctx_136 = ctx_88
	args_137 = args_89
	goto b12
b12:
	;
	v188 = v133
	op_kw_189 = op_kw_134
	form_190 = form_135
	ctx_191 = ctx_136
	args_192 = args_137
	goto b9
b13:
	;
	arg__1056_153, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_141})
	if callErr != nil {
		return nil, callErr
	}
	arg__1065_157, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_141})
	if callErr != nil {
		return nil, callErr
	}
	v159, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "add-inst!").Deref(), []vm.Value{ctx_141, arg__1065_157, op_kw_139, args_142, vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	v182 = v159
	op_kw_183 = op_kw_139
	form_184 = form_140
	ctx_185 = ctx_141
	args_186 = args_142
	goto b15
b14:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		op_kw_161 = op_kw_143
		form_162 = form_144
		ctx_163 = ctx_145
		args_164 = args_146
		goto b16
	} else {
		op_kw_165 = op_kw_143
		form_166 = form_144
		ctx_167 = ctx_145
		args_168 = args_146
		goto b17
	}
b15:
	;
	v188 = v182
	op_kw_189 = op_kw_183
	form_190 = form_184
	ctx_191 = ctx_185
	args_192 = args_186
	goto b9
b16:
	;
	v172, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "fold-binary-chain").Deref(), []vm.Value{op_kw_161, args_164, ctx_163})
	if callErr != nil {
		return nil, callErr
	}
	v176 = v172
	op_kw_177 = op_kw_161
	form_178 = form_162
	ctx_179 = ctx_163
	args_180 = args_164
	goto b18
b17:
	;
	v176 = vm.NIL
	op_kw_177 = op_kw_165
	form_178 = form_166
	ctx_179 = ctx_167
	args_180 = args_168
	goto b18
b18:
	;
	v182 = v176
	op_kw_183 = op_kw_177
	form_184 = form_178
	ctx_185 = ctx_179
	args_186 = args_180
	goto b15
}
func build_call(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var head_3 vm.Value
	var arg__1082_5 vm.Value
	var arg__1088_8 vm.Value
	var arg_ids_9 vm.Value
	var v19 vm.Value
	var form_10 vm.Value
	var ctx_11 vm.Value
	var head_12 vm.Value
	var arg_ids_13 vm.Value
	var v22 vm.Value
	var form_14 vm.Value
	var ctx_15 vm.Value
	var head_16 vm.Value
	var arg_ids_17 vm.Value
	var v25 vm.Value
	var fn_id_27 vm.Value
	var form_28 vm.Value
	var ctx_29 vm.Value
	var head_30 vm.Value
	var arg_ids_31 vm.Value
	var arg__1107_33 vm.Value
	var arg__1114_36 vm.Value
	var v37 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = head_3, arg__1082_5, arg__1088_8, arg_ids_9, v19, form_10, ctx_11, head_12, arg_ids_13, v22, form_14, ctx_15, head_16, arg_ids_17, v25, fn_id_27, form_28, ctx_29, head_30, arg_ids_31, arg__1107_33, arg__1114_36, v37
	head_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__1082_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__1088_8, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg_ids_9, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-args").Deref(), []vm.Value{arg__1088_8, arg1})
	if callErr != nil {
		return nil, callErr
	}
	v19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "symbol?").Deref(), []vm.Value{head_3})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v19) {
		form_10 = arg0
		ctx_11 = arg1
		head_12 = head_3
		arg_ids_13 = arg_ids_9
		goto b1
	} else {
		form_14 = arg0
		ctx_15 = arg1
		head_16 = head_3
		arg_ids_17 = arg_ids_9
		goto b2
	}
b1:
	;
	v22, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-symbol").Deref(), []vm.Value{head_12, ctx_11})
	if callErr != nil {
		return nil, callErr
	}
	fn_id_27 = v22
	form_28 = form_10
	ctx_29 = ctx_11
	head_30 = head_12
	arg_ids_31 = arg_ids_13
	goto b3
b2:
	;
	v25, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-form").Deref(), []vm.Value{head_16, ctx_15})
	if callErr != nil {
		return nil, callErr
	}
	fn_id_27 = v25
	form_28 = form_14
	ctx_29 = ctx_15
	head_30 = head_16
	arg_ids_31 = arg_ids_17
	goto b3
b3:
	;
	arg__1107_33, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{form_28})
	if callErr != nil {
		return nil, callErr
	}
	arg__1114_36, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{form_28})
	if callErr != nil {
		return nil, callErr
	}
	v37, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-call-with-head").Deref(), []vm.Value{fn_id_27, arg__1114_36, ctx_29})
	if callErr != nil {
		return nil, callErr
	}
	return v37, nil
}
func build_call_with_head(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var v5 vm.Value
	var head_sym_9 vm.Value
	var __11 vm.Value
	var arg_ids_13 vm.Value
	var threaded_15 vm.Value
	var v17 vm.Value
	var arg__1146_19 vm.Value
	var arg__1153_22 vm.Value
	var arg__1160_25 vm.Value
	var arg__1161_26 vm.Value
	var arg__1165_28 vm.Value
	var arg__1171_31 vm.Value
	var arg__1178_34 vm.Value
	var arg__1185_37 vm.Value
	var arg__1186_38 vm.Value
	var arg__1190_40 vm.Value
	var v41 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = v5, head_sym_9, __11, arg_ids_13, threaded_15, v17, arg__1146_19, arg__1153_22, arg__1160_25, arg__1161_26, arg__1165_28, arg__1171_31, arg__1178_34, arg__1185_37, arg__1186_38, arg__1190_40, v41
	v5, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "push-locals!").Deref(), []vm.Value{arg2})
	if callErr != nil {
		return nil, callErr
	}
	head_sym_9, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "gensym").Deref(), []vm.Value{vm.String("head__")})
	if callErr != nil {
		return nil, callErr
	}
	__11, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "bind-local!").Deref(), []vm.Value{arg2, head_sym_9, arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg_ids_13, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-args").Deref(), []vm.Value{arg1, arg2})
	if callErr != nil {
		return nil, callErr
	}
	threaded_15, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "lookup-local").Deref(), []vm.Value{arg2, head_sym_9})
	if callErr != nil {
		return nil, callErr
	}
	v17, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "pop-locals!").Deref(), []vm.Value{arg2})
	if callErr != nil {
		return nil, callErr
	}
	arg__1146_19, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{arg2})
	if callErr != nil {
		return nil, callErr
	}
	arg__1153_22, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "cons").Deref(), []vm.Value{threaded_15, arg_ids_13})
	if callErr != nil {
		return nil, callErr
	}
	arg__1160_25, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "cons").Deref(), []vm.Value{threaded_15, arg_ids_13})
	if callErr != nil {
		return nil, callErr
	}
	arg__1161_26, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg__1160_25})
	if callErr != nil {
		return nil, callErr
	}
	arg__1165_28, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{arg_ids_13})
	if callErr != nil {
		return nil, callErr
	}
	arg__1171_31, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{arg2})
	if callErr != nil {
		return nil, callErr
	}
	arg__1178_34, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "cons").Deref(), []vm.Value{threaded_15, arg_ids_13})
	if callErr != nil {
		return nil, callErr
	}
	arg__1185_37, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "cons").Deref(), []vm.Value{threaded_15, arg_ids_13})
	if callErr != nil {
		return nil, callErr
	}
	arg__1186_38, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg__1185_37})
	if callErr != nil {
		return nil, callErr
	}
	arg__1190_40, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{arg_ids_13})
	if callErr != nil {
		return nil, callErr
	}
	v41, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "add-inst!").Deref(), []vm.Value{arg2, arg__1171_31, vm.Keyword("call"), arg__1186_38, arg__1190_40})
	if callErr != nil {
		return nil, callErr
	}
	return v41, nil
}
func build_do(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var forms_4 vm.Value
	var fs_5 vm.Value
	var last_id_6 vm.Value
	var ctx_7 vm.Value
	var v21 vm.Value
	var form_10 vm.Value
	var forms_11 vm.Value
	var fs_12 vm.Value
	var last_id_13 vm.Value
	var ctx_14 vm.Value
	var v24 vm.Value
	var arg__1203_26 vm.Value
	var arg__1209_29 vm.Value
	var v30 vm.Value
	var form_15 vm.Value
	var forms_16 vm.Value
	var fs_17 vm.Value
	var last_id_18 vm.Value
	var ctx_19 vm.Value
	var v33 vm.Value
	var form_34 vm.Value
	var forms_35 vm.Value
	var fs_36 vm.Value
	var last_id_37 vm.Value
	var ctx_38 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = forms_4, fs_5, last_id_6, ctx_7, v21, form_10, forms_11, fs_12, last_id_13, ctx_14, v24, arg__1203_26, arg__1209_29, v30, form_15, forms_16, fs_17, last_id_18, ctx_19, v33, form_34, forms_35, fs_36, last_id_37, ctx_38
	forms_4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	fs_5 = forms_4
	last_id_6 = vm.NIL
	ctx_7 = arg1
	goto b1
b1:
	;
	v21, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{fs_5})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v21) {
		form_10 = arg0
		forms_11 = forms_4
		fs_12 = fs_5
		last_id_13 = last_id_6
		ctx_14 = ctx_7
		goto b2
	} else {
		form_15 = arg0
		forms_16 = forms_4
		fs_17 = fs_5
		last_id_18 = last_id_6
		ctx_19 = ctx_7
		goto b3
	}
b2:
	;
	v24, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{fs_12})
	if callErr != nil {
		return nil, callErr
	}
	arg__1203_26, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{fs_12})
	if callErr != nil {
		return nil, callErr
	}
	arg__1209_29, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{fs_12})
	if callErr != nil {
		return nil, callErr
	}
	v30, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-form").Deref(), []vm.Value{arg__1209_29, ctx_14})
	if callErr != nil {
		return nil, callErr
	}
	fs_5 = v24
	last_id_6 = v30
	ctx_7 = ctx_14
	goto b1
b3:
	;
	v33 = last_id_18
	form_34 = form_15
	forms_35 = forms_16
	fs_36 = fs_17
	last_id_37 = last_id_18
	ctx_38 = ctx_19
	goto b4
b4:
	;
	return v33, nil
}
func build_fn(arg0 vm.Value) (vm.Value, error) {
	var name_sym_5 vm.Value
	var maybe_doc_9 vm.Value
	var has_doc_QMARK__11 vm.Value
	var defn_form_12 vm.Value
	var name_sym_13 vm.Value
	var maybe_doc_14 vm.Value
	var has_doc_QMARK__15 vm.Value
	var v24 vm.Value
	var defn_form_16 vm.Value
	var name_sym_17 vm.Value
	var maybe_doc_18 vm.Value
	var has_doc_QMARK__19 vm.Value
	var args_vec_27 vm.Value
	var defn_form_28 vm.Value
	var name_sym_29 vm.Value
	var maybe_doc_30 vm.Value
	var has_doc_QMARK__31 vm.Value
	var args_vec_32 vm.Value
	var defn_form_33 vm.Value
	var name_sym_34 vm.Value
	var maybe_doc_35 vm.Value
	var has_doc_QMARK__36 vm.Value
	var args_vec_37 vm.Value
	var defn_form_38 vm.Value
	var name_sym_39 vm.Value
	var maybe_doc_40 vm.Value
	var has_doc_QMARK__41 vm.Value
	var arg__1231_47 int
	var args_vec_48 vm.Value
	var defn_form_49 vm.Value
	var name_sym_50 vm.Value
	var maybe_doc_51 vm.Value
	var has_doc_QMARK__52 vm.Value
	var args_vec_54 vm.Value
	var defn_form_55 vm.Value
	var name_sym_56 vm.Value
	var maybe_doc_57 vm.Value
	var has_doc_QMARK__58 vm.Value
	var head__1233_59 vm.Value
	var args_vec_60 vm.Value
	var defn_form_61 vm.Value
	var name_sym_62 vm.Value
	var maybe_doc_63 vm.Value
	var has_doc_QMARK__64 vm.Value
	var head__1233_65 vm.Value
	var arg__1234_71 int
	var args_vec_72 vm.Value
	var defn_form_73 vm.Value
	var name_sym_74 vm.Value
	var maybe_doc_75 vm.Value
	var has_doc_QMARK__76 vm.Value
	var head__1233_77 vm.Value
	var body_forms_78 vm.Value
	var multi_QMARK__80 vm.Value
	var args_vec_81 vm.Value
	var defn_form_82 vm.Value
	var name_sym_83 vm.Value
	var maybe_doc_84 vm.Value
	var has_doc_QMARK__85 vm.Value
	var body_forms_86 vm.Value
	var multi_QMARK__87 vm.Value
	var arg__1242_97 vm.Value
	var arg__1249_102 vm.Value
	var f_105 vm.Value
	var entry_blk_107 vm.Value
	var ctx_109 vm.Value
	var arities_111 vm.Value
	var expanded_forms_115 vm.Value
	var arg__1373_120 vm.Value
	var arg__1445_126 vm.Value
	var all_caps_127 vm.Value
	var arg__1452_131 vm.Value
	var arg__1459_136 vm.Value
	var captures_137 vm.Value
	var templates_147 vm.Value
	var template_152 vm.Value
	var closure_id_154 vm.Value
	var final_blk_156 vm.Value
	var arg__1537_159 vm.Value
	var arg__1545_164 vm.Value
	var v166 vm.Value
	var args_vec_88 vm.Value
	var defn_form_89 vm.Value
	var name_sym_90 vm.Value
	var maybe_doc_91 vm.Value
	var has_doc_QMARK__92 vm.Value
	var body_forms_93 vm.Value
	var multi_QMARK__94 vm.Value
	var expanded_169 vm.Value
	var v691 vm.Value
	var args_vec_692 vm.Value
	var defn_form_693 vm.Value
	var name_sym_694 vm.Value
	var maybe_doc_695 vm.Value
	var has_doc_QMARK__696 vm.Value
	var body_forms_697 vm.Value
	var multi_QMARK__698 vm.Value
	var args_vec_170 vm.Value
	var defn_form_171 vm.Value
	var name_sym_172 vm.Value
	var maybe_doc_173 vm.Value
	var has_doc_QMARK__174 vm.Value
	var body_forms_175 vm.Value
	var multi_QMARK__176 vm.Value
	var expanded_177 vm.Value
	var v188 vm.Value
	var args_vec_178 vm.Value
	var defn_form_179 vm.Value
	var name_sym_180 vm.Value
	var maybe_doc_181 vm.Value
	var has_doc_QMARK__182 vm.Value
	var body_forms_183 vm.Value
	var multi_QMARK__184 vm.Value
	var expanded_185 vm.Value
	var flat_args_191 vm.Value
	var args_vec_192 vm.Value
	var defn_form_193 vm.Value
	var name_sym_194 vm.Value
	var maybe_doc_195 vm.Value
	var has_doc_QMARK__196 vm.Value
	var body_forms_197 vm.Value
	var multi_QMARK__198 vm.Value
	var expanded_199 vm.Value
	var flat_args_200 vm.Value
	var args_vec_201 vm.Value
	var defn_form_202 vm.Value
	var name_sym_203 vm.Value
	var maybe_doc_204 vm.Value
	var has_doc_QMARK__205 vm.Value
	var body_forms_206 vm.Value
	var multi_QMARK__207 vm.Value
	var expanded_208 vm.Value
	var v220 vm.Value
	var flat_args_209 vm.Value
	var args_vec_210 vm.Value
	var defn_form_211 vm.Value
	var name_sym_212 vm.Value
	var maybe_doc_213 vm.Value
	var has_doc_QMARK__214 vm.Value
	var body_forms_215 vm.Value
	var multi_QMARK__216 vm.Value
	var expanded_217 vm.Value
	var flat_body_223 vm.Value
	var flat_args_224 vm.Value
	var args_vec_225 vm.Value
	var defn_form_226 vm.Value
	var name_sym_227 vm.Value
	var maybe_doc_228 vm.Value
	var has_doc_QMARK__229 vm.Value
	var body_forms_230 vm.Value
	var multi_QMARK__231 vm.Value
	var expanded_232 vm.Value
	var flat_body_233 vm.Value
	var flat_args_234 vm.Value
	var args_vec_235 vm.Value
	var defn_form_236 vm.Value
	var name_sym_237 vm.Value
	var maybe_doc_238 vm.Value
	var has_doc_QMARK__239 vm.Value
	var body_forms_240 vm.Value
	var multi_QMARK__241 vm.Value
	var expanded_242 vm.Value
	var v255 vm.Value
	var flat_body_243 vm.Value
	var flat_args_244 vm.Value
	var args_vec_245 vm.Value
	var defn_form_246 vm.Value
	var name_sym_247 vm.Value
	var maybe_doc_248 vm.Value
	var has_doc_QMARK__249 vm.Value
	var body_forms_250 vm.Value
	var multi_QMARK__251 vm.Value
	var expanded_252 vm.Value
	var variadic_QMARK__259 vm.Value
	var flat_body_260 vm.Value
	var flat_args_261 vm.Value
	var args_vec_262 vm.Value
	var defn_form_263 vm.Value
	var name_sym_264 vm.Value
	var maybe_doc_265 vm.Value
	var has_doc_QMARK__266 vm.Value
	var body_forms_267 vm.Value
	var multi_QMARK__268 vm.Value
	var expanded_269 vm.Value
	var arity_271 vm.Value
	var arg__1564_273 vm.Value
	var arg__1571_276 vm.Value
	var f_277 vm.Value
	var entry_blk_279 vm.Value
	var ctx_281 vm.Value
	var arg__1584_283 vm.Value
	var arg__1589_286 vm.Value
	var arg__1594_289 vm.Value
	var arg__1599_292 vm.Value
	var arg__1600_293 vm.Value
	var arg__1606_296 vm.Value
	var arg__1611_299 vm.Value
	var arg__1616_302 vm.Value
	var arg__1621_305 vm.Value
	var arg__1622_306 vm.Value
	var v307 vm.Value
	var i_308 int
	var arity_309 vm.Value
	var flat_args_310 vm.Value
	var ctx_311 vm.Value
	var entry_blk_312 vm.Value
	var f_313 vm.Value
	var v704 vm.Value
	var v707 vm.Value
	var v348 bool
	var variadic_QMARK__316 vm.Value
	var flat_body_317 vm.Value
	var args_vec_318 vm.Value
	var defn_form_319 vm.Value
	var name_sym_320 vm.Value
	var maybe_doc_321 vm.Value
	var has_doc_QMARK__322 vm.Value
	var body_forms_323 vm.Value
	var multi_QMARK__324 vm.Value
	var expanded_325 vm.Value
	var i_326 int
	var arity_327 vm.Value
	var flat_args_328 vm.Value
	var ctx_329 vm.Value
	var entry_blk_330 vm.Value
	var f_331 vm.Value
	var v705 vm.Value
	var v708 vm.Value
	var arg_id_355 vm.Value
	var arg__1642_357 vm.Value
	var arg__1651_360 vm.Value
	var v361 vm.Value
	var v362 int
	var variadic_QMARK__332 vm.Value
	var flat_body_333 vm.Value
	var args_vec_334 vm.Value
	var defn_form_335 vm.Value
	var name_sym_336 vm.Value
	var maybe_doc_337 vm.Value
	var has_doc_QMARK__338 vm.Value
	var body_forms_339 vm.Value
	var multi_QMARK__340 vm.Value
	var expanded_341 vm.Value
	var i_342 int
	var arity_343 vm.Value
	var flat_args_344 vm.Value
	var ctx_345 vm.Value
	var entry_blk_346 vm.Value
	var f_347 vm.Value
	var v706 vm.Value
	var v709 vm.Value
	var v366 vm.Value
	var variadic_QMARK__367 vm.Value
	var flat_body_368 vm.Value
	var args_vec_369 vm.Value
	var defn_form_370 vm.Value
	var name_sym_371 vm.Value
	var maybe_doc_372 vm.Value
	var has_doc_QMARK__373 vm.Value
	var body_forms_374 vm.Value
	var multi_QMARK__375 vm.Value
	var expanded_376 vm.Value
	var i_377 int
	var arity_378 vm.Value
	var flat_args_379 vm.Value
	var ctx_380 vm.Value
	var entry_blk_381 vm.Value
	var f_382 vm.Value
	var fs_383 vm.Value
	var last_id_384 vm.Value
	var ctx_385 vm.Value
	var v425 vm.Value
	var variadic_QMARK__388 vm.Value
	var flat_body_389 vm.Value
	var args_vec_390 vm.Value
	var defn_form_391 vm.Value
	var name_sym_392 vm.Value
	var maybe_doc_393 vm.Value
	var has_doc_QMARK__394 vm.Value
	var body_forms_395 vm.Value
	var multi_QMARK__396 vm.Value
	var expanded_397 vm.Value
	var i_398 int
	var arity_399 vm.Value
	var flat_args_400 vm.Value
	var entry_blk_401 vm.Value
	var f_402 vm.Value
	var fs_403 vm.Value
	var last_id_404 vm.Value
	var ctx_405 vm.Value
	var v428 vm.Value
	var arg__1663_430 vm.Value
	var arg__1669_433 vm.Value
	var v434 vm.Value
	var variadic_QMARK__406 vm.Value
	var flat_body_407 vm.Value
	var args_vec_408 vm.Value
	var defn_form_409 vm.Value
	var name_sym_410 vm.Value
	var maybe_doc_411 vm.Value
	var has_doc_QMARK__412 vm.Value
	var body_forms_413 vm.Value
	var multi_QMARK__414 vm.Value
	var expanded_415 vm.Value
	var i_416 int
	var arity_417 vm.Value
	var flat_args_418 vm.Value
	var entry_blk_419 vm.Value
	var f_420 vm.Value
	var fs_421 vm.Value
	var last_id_422 vm.Value
	var ctx_423 vm.Value
	var last_val_437 vm.Value
	var variadic_QMARK__438 vm.Value
	var flat_body_439 vm.Value
	var args_vec_440 vm.Value
	var defn_form_441 vm.Value
	var name_sym_442 vm.Value
	var maybe_doc_443 vm.Value
	var has_doc_QMARK__444 vm.Value
	var body_forms_445 vm.Value
	var multi_QMARK__446 vm.Value
	var expanded_447 vm.Value
	var i_448 int
	var arity_449 vm.Value
	var flat_args_450 vm.Value
	var entry_blk_451 vm.Value
	var f_452 vm.Value
	var fs_453 vm.Value
	var last_id_454 vm.Value
	var ctx_455 vm.Value
	var final_blk_457 vm.Value
	var v499 vm.Value
	var last_val_458 vm.Value
	var variadic_QMARK__459 vm.Value
	var flat_body_460 vm.Value
	var args_vec_461 vm.Value
	var defn_form_462 vm.Value
	var name_sym_463 vm.Value
	var maybe_doc_464 vm.Value
	var has_doc_QMARK__465 vm.Value
	var body_forms_466 vm.Value
	var multi_QMARK__467 vm.Value
	var expanded_468 vm.Value
	var i_469 int
	var arity_470 vm.Value
	var flat_args_471 vm.Value
	var entry_blk_472 vm.Value
	var f_473 vm.Value
	var fs_474 vm.Value
	var last_id_475 vm.Value
	var ctx_476 vm.Value
	var final_blk_477 vm.Value
	var last_val_478 vm.Value
	var variadic_QMARK__479 vm.Value
	var flat_body_480 vm.Value
	var args_vec_481 vm.Value
	var defn_form_482 vm.Value
	var name_sym_483 vm.Value
	var maybe_doc_484 vm.Value
	var has_doc_QMARK__485 vm.Value
	var body_forms_486 vm.Value
	var multi_QMARK__487 vm.Value
	var expanded_488 vm.Value
	var i_489 int
	var arity_490 vm.Value
	var flat_args_491 vm.Value
	var entry_blk_492 vm.Value
	var f_493 vm.Value
	var fs_494 vm.Value
	var last_id_495 vm.Value
	var ctx_496 vm.Value
	var final_blk_497 vm.Value
	var v551 vm.Value
	var v669 vm.Value
	var last_val_670 vm.Value
	var variadic_QMARK__671 vm.Value
	var flat_body_672 vm.Value
	var args_vec_673 vm.Value
	var defn_form_674 vm.Value
	var name_sym_675 vm.Value
	var maybe_doc_676 vm.Value
	var has_doc_QMARK__677 vm.Value
	var body_forms_678 vm.Value
	var multi_QMARK__679 vm.Value
	var expanded_680 vm.Value
	var i_681 int
	var arity_682 vm.Value
	var flat_args_683 vm.Value
	var entry_blk_684 vm.Value
	var f_685 vm.Value
	var fs_686 vm.Value
	var last_id_687 vm.Value
	var ctx_688 vm.Value
	var final_blk_689 vm.Value
	var last_val_504 vm.Value
	var variadic_QMARK__505 vm.Value
	var flat_body_506 vm.Value
	var args_vec_507 vm.Value
	var defn_form_508 vm.Value
	var name_sym_509 vm.Value
	var maybe_doc_510 vm.Value
	var has_doc_QMARK__511 vm.Value
	var body_forms_512 vm.Value
	var multi_QMARK__513 vm.Value
	var expanded_514 vm.Value
	var i_515 int
	var arity_516 vm.Value
	var flat_args_517 vm.Value
	var entry_blk_518 vm.Value
	var arg__1677_519 vm.Value
	var f_520 vm.Value
	var fs_521 vm.Value
	var last_id_522 vm.Value
	var ctx_523 vm.Value
	var final_blk_524 vm.Value
	var arg__1678_525 vm.Value
	var arg__1679_526 vm.Value
	var last_val_527 vm.Value
	var variadic_QMARK__528 vm.Value
	var flat_body_529 vm.Value
	var args_vec_530 vm.Value
	var defn_form_531 vm.Value
	var name_sym_532 vm.Value
	var maybe_doc_533 vm.Value
	var has_doc_QMARK__534 vm.Value
	var body_forms_535 vm.Value
	var multi_QMARK__536 vm.Value
	var expanded_537 vm.Value
	var i_538 int
	var arity_539 vm.Value
	var flat_args_540 vm.Value
	var entry_blk_541 vm.Value
	var arg__1677_542 vm.Value
	var f_543 vm.Value
	var fs_544 vm.Value
	var last_id_545 vm.Value
	var ctx_546 vm.Value
	var final_blk_547 vm.Value
	var arg__1678_548 vm.Value
	var arg__1679_549 vm.Value
	var v556 vm.Value
	var arg__1685_558 vm.Value
	var last_val_559 vm.Value
	var variadic_QMARK__560 vm.Value
	var flat_body_561 vm.Value
	var args_vec_562 vm.Value
	var defn_form_563 vm.Value
	var name_sym_564 vm.Value
	var maybe_doc_565 vm.Value
	var has_doc_QMARK__566 vm.Value
	var body_forms_567 vm.Value
	var multi_QMARK__568 vm.Value
	var expanded_569 vm.Value
	var i_570 int
	var arity_571 vm.Value
	var flat_args_572 vm.Value
	var entry_blk_573 vm.Value
	var arg__1677_574 vm.Value
	var f_575 vm.Value
	var fs_576 vm.Value
	var last_id_577 vm.Value
	var ctx_578 vm.Value
	var final_blk_579 vm.Value
	var arg__1678_580 vm.Value
	var arg__1679_581 vm.Value
	var v634 vm.Value
	var last_val_585 vm.Value
	var variadic_QMARK__586 vm.Value
	var flat_body_587 vm.Value
	var args_vec_588 vm.Value
	var defn_form_589 vm.Value
	var name_sym_590 vm.Value
	var maybe_doc_591 vm.Value
	var has_doc_QMARK__592 vm.Value
	var body_forms_593 vm.Value
	var multi_QMARK__594 vm.Value
	var expanded_595 vm.Value
	var i_596 int
	var arity_597 vm.Value
	var flat_args_598 vm.Value
	var entry_blk_599 vm.Value
	var f_600 vm.Value
	var arg__1688_601 vm.Value
	var fs_602 vm.Value
	var last_id_603 vm.Value
	var ctx_604 vm.Value
	var final_blk_605 vm.Value
	var arg__1689_606 vm.Value
	var head__1687_607 vm.Value
	var arg__1690_608 vm.Value
	var last_val_609 vm.Value
	var variadic_QMARK__610 vm.Value
	var flat_body_611 vm.Value
	var args_vec_612 vm.Value
	var defn_form_613 vm.Value
	var name_sym_614 vm.Value
	var maybe_doc_615 vm.Value
	var has_doc_QMARK__616 vm.Value
	var body_forms_617 vm.Value
	var multi_QMARK__618 vm.Value
	var expanded_619 vm.Value
	var i_620 int
	var arity_621 vm.Value
	var flat_args_622 vm.Value
	var entry_blk_623 vm.Value
	var f_624 vm.Value
	var arg__1688_625 vm.Value
	var fs_626 vm.Value
	var last_id_627 vm.Value
	var ctx_628 vm.Value
	var final_blk_629 vm.Value
	var arg__1689_630 vm.Value
	var head__1687_631 vm.Value
	var arg__1690_632 vm.Value
	var v639 vm.Value
	var arg__1696_641 vm.Value
	var last_val_642 vm.Value
	var variadic_QMARK__643 vm.Value
	var flat_body_644 vm.Value
	var args_vec_645 vm.Value
	var defn_form_646 vm.Value
	var name_sym_647 vm.Value
	var maybe_doc_648 vm.Value
	var has_doc_QMARK__649 vm.Value
	var body_forms_650 vm.Value
	var multi_QMARK__651 vm.Value
	var expanded_652 vm.Value
	var i_653 int
	var arity_654 vm.Value
	var flat_args_655 vm.Value
	var entry_blk_656 vm.Value
	var f_657 vm.Value
	var arg__1688_658 vm.Value
	var fs_659 vm.Value
	var last_id_660 vm.Value
	var ctx_661 vm.Value
	var final_blk_662 vm.Value
	var arg__1689_663 vm.Value
	var head__1687_664 vm.Value
	var arg__1690_665 vm.Value
	var v667 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = name_sym_5, maybe_doc_9, has_doc_QMARK__11, defn_form_12, name_sym_13, maybe_doc_14, has_doc_QMARK__15, v24, defn_form_16, name_sym_17, maybe_doc_18, has_doc_QMARK__19, args_vec_27, defn_form_28, name_sym_29, maybe_doc_30, has_doc_QMARK__31, args_vec_32, defn_form_33, name_sym_34, maybe_doc_35, has_doc_QMARK__36, args_vec_37, defn_form_38, name_sym_39, maybe_doc_40, has_doc_QMARK__41, arg__1231_47, args_vec_48, defn_form_49, name_sym_50, maybe_doc_51, has_doc_QMARK__52, args_vec_54, defn_form_55, name_sym_56, maybe_doc_57, has_doc_QMARK__58, head__1233_59, args_vec_60, defn_form_61, name_sym_62, maybe_doc_63, has_doc_QMARK__64, head__1233_65, arg__1234_71, args_vec_72, defn_form_73, name_sym_74, maybe_doc_75, has_doc_QMARK__76, head__1233_77, body_forms_78, multi_QMARK__80, args_vec_81, defn_form_82, name_sym_83, maybe_doc_84, has_doc_QMARK__85, body_forms_86, multi_QMARK__87, arg__1242_97, arg__1249_102, f_105, entry_blk_107, ctx_109, arities_111, expanded_forms_115, arg__1373_120, arg__1445_126, all_caps_127, arg__1452_131, arg__1459_136, captures_137, templates_147, template_152, closure_id_154, final_blk_156, arg__1537_159, arg__1545_164, v166, args_vec_88, defn_form_89, name_sym_90, maybe_doc_91, has_doc_QMARK__92, body_forms_93, multi_QMARK__94, expanded_169, v691, args_vec_692, defn_form_693, name_sym_694, maybe_doc_695, has_doc_QMARK__696, body_forms_697, multi_QMARK__698, args_vec_170, defn_form_171, name_sym_172, maybe_doc_173, has_doc_QMARK__174, body_forms_175, multi_QMARK__176, expanded_177, v188, args_vec_178, defn_form_179, name_sym_180, maybe_doc_181, has_doc_QMARK__182, body_forms_183, multi_QMARK__184, expanded_185, flat_args_191, args_vec_192, defn_form_193, name_sym_194, maybe_doc_195, has_doc_QMARK__196, body_forms_197, multi_QMARK__198, expanded_199, flat_args_200, args_vec_201, defn_form_202, name_sym_203, maybe_doc_204, has_doc_QMARK__205, body_forms_206, multi_QMARK__207, expanded_208, v220, flat_args_209, args_vec_210, defn_form_211, name_sym_212, maybe_doc_213, has_doc_QMARK__214, body_forms_215, multi_QMARK__216, expanded_217, flat_body_223, flat_args_224, args_vec_225, defn_form_226, name_sym_227, maybe_doc_228, has_doc_QMARK__229, body_forms_230, multi_QMARK__231, expanded_232, flat_body_233, flat_args_234, args_vec_235, defn_form_236, name_sym_237, maybe_doc_238, has_doc_QMARK__239, body_forms_240, multi_QMARK__241, expanded_242, v255, flat_body_243, flat_args_244, args_vec_245, defn_form_246, name_sym_247, maybe_doc_248, has_doc_QMARK__249, body_forms_250, multi_QMARK__251, expanded_252, variadic_QMARK__259, flat_body_260, flat_args_261, args_vec_262, defn_form_263, name_sym_264, maybe_doc_265, has_doc_QMARK__266, body_forms_267, multi_QMARK__268, expanded_269, arity_271, arg__1564_273, arg__1571_276, f_277, entry_blk_279, ctx_281, arg__1584_283, arg__1589_286, arg__1594_289, arg__1599_292, arg__1600_293, arg__1606_296, arg__1611_299, arg__1616_302, arg__1621_305, arg__1622_306, v307, i_308, arity_309, flat_args_310, ctx_311, entry_blk_312, f_313, v704, v707, v348, variadic_QMARK__316, flat_body_317, args_vec_318, defn_form_319, name_sym_320, maybe_doc_321, has_doc_QMARK__322, body_forms_323, multi_QMARK__324, expanded_325, i_326, arity_327, flat_args_328, ctx_329, entry_blk_330, f_331, v705, v708, arg_id_355, arg__1642_357, arg__1651_360, v361, v362, variadic_QMARK__332, flat_body_333, args_vec_334, defn_form_335, name_sym_336, maybe_doc_337, has_doc_QMARK__338, body_forms_339, multi_QMARK__340, expanded_341, i_342, arity_343, flat_args_344, ctx_345, entry_blk_346, f_347, v706, v709, v366, variadic_QMARK__367, flat_body_368, args_vec_369, defn_form_370, name_sym_371, maybe_doc_372, has_doc_QMARK__373, body_forms_374, multi_QMARK__375, expanded_376, i_377, arity_378, flat_args_379, ctx_380, entry_blk_381, f_382, fs_383, last_id_384, ctx_385, v425, variadic_QMARK__388, flat_body_389, args_vec_390, defn_form_391, name_sym_392, maybe_doc_393, has_doc_QMARK__394, body_forms_395, multi_QMARK__396, expanded_397, i_398, arity_399, flat_args_400, entry_blk_401, f_402, fs_403, last_id_404, ctx_405, v428, arg__1663_430, arg__1669_433, v434, variadic_QMARK__406, flat_body_407, args_vec_408, defn_form_409, name_sym_410, maybe_doc_411, has_doc_QMARK__412, body_forms_413, multi_QMARK__414, expanded_415, i_416, arity_417, flat_args_418, entry_blk_419, f_420, fs_421, last_id_422, ctx_423, last_val_437, variadic_QMARK__438, flat_body_439, args_vec_440, defn_form_441, name_sym_442, maybe_doc_443, has_doc_QMARK__444, body_forms_445, multi_QMARK__446, expanded_447, i_448, arity_449, flat_args_450, entry_blk_451, f_452, fs_453, last_id_454, ctx_455, final_blk_457, v499, last_val_458, variadic_QMARK__459, flat_body_460, args_vec_461, defn_form_462, name_sym_463, maybe_doc_464, has_doc_QMARK__465, body_forms_466, multi_QMARK__467, expanded_468, i_469, arity_470, flat_args_471, entry_blk_472, f_473, fs_474, last_id_475, ctx_476, final_blk_477, last_val_478, variadic_QMARK__479, flat_body_480, args_vec_481, defn_form_482, name_sym_483, maybe_doc_484, has_doc_QMARK__485, body_forms_486, multi_QMARK__487, expanded_488, i_489, arity_490, flat_args_491, entry_blk_492, f_493, fs_494, last_id_495, ctx_496, final_blk_497, v551, v669, last_val_670, variadic_QMARK__671, flat_body_672, args_vec_673, defn_form_674, name_sym_675, maybe_doc_676, has_doc_QMARK__677, body_forms_678, multi_QMARK__679, expanded_680, i_681, arity_682, flat_args_683, entry_blk_684, f_685, fs_686, last_id_687, ctx_688, final_blk_689, last_val_504, variadic_QMARK__505, flat_body_506, args_vec_507, defn_form_508, name_sym_509, maybe_doc_510, has_doc_QMARK__511, body_forms_512, multi_QMARK__513, expanded_514, i_515, arity_516, flat_args_517, entry_blk_518, arg__1677_519, f_520, fs_521, last_id_522, ctx_523, final_blk_524, arg__1678_525, arg__1679_526, last_val_527, variadic_QMARK__528, flat_body_529, args_vec_530, defn_form_531, name_sym_532, maybe_doc_533, has_doc_QMARK__534, body_forms_535, multi_QMARK__536, expanded_537, i_538, arity_539, flat_args_540, entry_blk_541, arg__1677_542, f_543, fs_544, last_id_545, ctx_546, final_blk_547, arg__1678_548, arg__1679_549, v556, arg__1685_558, last_val_559, variadic_QMARK__560, flat_body_561, args_vec_562, defn_form_563, name_sym_564, maybe_doc_565, has_doc_QMARK__566, body_forms_567, multi_QMARK__568, expanded_569, i_570, arity_571, flat_args_572, entry_blk_573, arg__1677_574, f_575, fs_576, last_id_577, ctx_578, final_blk_579, arg__1678_580, arg__1679_581, v634, last_val_585, variadic_QMARK__586, flat_body_587, args_vec_588, defn_form_589, name_sym_590, maybe_doc_591, has_doc_QMARK__592, body_forms_593, multi_QMARK__594, expanded_595, i_596, arity_597, flat_args_598, entry_blk_599, f_600, arg__1688_601, fs_602, last_id_603, ctx_604, final_blk_605, arg__1689_606, head__1687_607, arg__1690_608, last_val_609, variadic_QMARK__610, flat_body_611, args_vec_612, defn_form_613, name_sym_614, maybe_doc_615, has_doc_QMARK__616, body_forms_617, multi_QMARK__618, expanded_619, i_620, arity_621, flat_args_622, entry_blk_623, f_624, arg__1688_625, fs_626, last_id_627, ctx_628, final_blk_629, arg__1689_630, head__1687_631, arg__1690_632, v639, arg__1696_641, last_val_642, variadic_QMARK__643, flat_body_644, args_vec_645, defn_form_646, name_sym_647, maybe_doc_648, has_doc_QMARK__649, body_forms_650, multi_QMARK__651, expanded_652, i_653, arity_654, flat_args_655, entry_blk_656, f_657, arg__1688_658, fs_659, last_id_660, ctx_661, final_blk_662, arg__1689_663, head__1687_664, arg__1690_665, v667
	name_sym_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(1)})
	if callErr != nil {
		return nil, callErr
	}
	maybe_doc_9, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(2)})
	if callErr != nil {
		return nil, callErr
	}
	has_doc_QMARK__11, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "string?").Deref(), []vm.Value{maybe_doc_9})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(has_doc_QMARK__11) {
		defn_form_12 = arg0
		name_sym_13 = name_sym_5
		maybe_doc_14 = maybe_doc_9
		has_doc_QMARK__15 = has_doc_QMARK__11
		goto b1
	} else {
		defn_form_16 = arg0
		name_sym_17 = name_sym_5
		maybe_doc_18 = maybe_doc_9
		has_doc_QMARK__19 = has_doc_QMARK__11
		goto b2
	}
b1:
	;
	v24, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{defn_form_12, vm.Int(3)})
	if callErr != nil {
		return nil, callErr
	}
	args_vec_27 = v24
	defn_form_28 = defn_form_12
	name_sym_29 = name_sym_13
	maybe_doc_30 = maybe_doc_14
	has_doc_QMARK__31 = has_doc_QMARK__15
	goto b3
b2:
	;
	args_vec_27 = maybe_doc_18
	defn_form_28 = defn_form_16
	name_sym_29 = name_sym_17
	maybe_doc_30 = maybe_doc_18
	has_doc_QMARK__31 = has_doc_QMARK__19
	goto b3
b3:
	;
	if vm.IsTruthy(has_doc_QMARK__31) {
		args_vec_32 = args_vec_27
		defn_form_33 = defn_form_28
		name_sym_34 = name_sym_29
		maybe_doc_35 = maybe_doc_30
		has_doc_QMARK__36 = has_doc_QMARK__31
		goto b4
	} else {
		args_vec_37 = args_vec_27
		defn_form_38 = defn_form_28
		name_sym_39 = name_sym_29
		maybe_doc_40 = maybe_doc_30
		has_doc_QMARK__41 = has_doc_QMARK__31
		goto b5
	}
b4:
	;
	arg__1231_47 = 4
	args_vec_48 = args_vec_32
	defn_form_49 = defn_form_33
	name_sym_50 = name_sym_34
	maybe_doc_51 = maybe_doc_35
	has_doc_QMARK__52 = has_doc_QMARK__36
	goto b6
b5:
	;
	arg__1231_47 = 3
	args_vec_48 = args_vec_37
	defn_form_49 = defn_form_38
	name_sym_50 = name_sym_39
	maybe_doc_51 = maybe_doc_40
	has_doc_QMARK__52 = has_doc_QMARK__41
	goto b6
b6:
	;
	if vm.IsTruthy(has_doc_QMARK__52) {
		args_vec_54 = args_vec_48
		defn_form_55 = defn_form_49
		name_sym_56 = name_sym_50
		maybe_doc_57 = maybe_doc_51
		has_doc_QMARK__58 = has_doc_QMARK__52
		head__1233_59 = rt.LookupVar("clojure.core", "drop").Deref()
		goto b7
	} else {
		args_vec_60 = args_vec_48
		defn_form_61 = defn_form_49
		name_sym_62 = name_sym_50
		maybe_doc_63 = maybe_doc_51
		has_doc_QMARK__64 = has_doc_QMARK__52
		head__1233_65 = rt.LookupVar("clojure.core", "drop").Deref()
		goto b8
	}
b7:
	;
	arg__1234_71 = 4
	args_vec_72 = args_vec_54
	defn_form_73 = defn_form_55
	name_sym_74 = name_sym_56
	maybe_doc_75 = maybe_doc_57
	has_doc_QMARK__76 = has_doc_QMARK__58
	head__1233_77 = head__1233_59
	goto b9
b8:
	;
	arg__1234_71 = 3
	args_vec_72 = args_vec_60
	defn_form_73 = defn_form_61
	name_sym_74 = name_sym_62
	maybe_doc_75 = maybe_doc_63
	has_doc_QMARK__76 = has_doc_QMARK__64
	head__1233_77 = head__1233_65
	goto b9
b9:
	;
	body_forms_78, callErr = rt.InvokeValue(head__1233_77, []vm.Value{vm.Int(arg__1234_71), defn_form_73})
	if callErr != nil {
		return nil, callErr
	}
	multi_QMARK__80, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "list?").Deref(), []vm.Value{args_vec_72})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(multi_QMARK__80) {
		args_vec_81 = args_vec_72
		defn_form_82 = defn_form_73
		name_sym_83 = name_sym_74
		maybe_doc_84 = maybe_doc_75
		has_doc_QMARK__85 = has_doc_QMARK__76
		body_forms_86 = body_forms_78
		multi_QMARK__87 = multi_QMARK__80
		goto b10
	} else {
		args_vec_88 = args_vec_72
		defn_form_89 = defn_form_73
		name_sym_90 = name_sym_74
		maybe_doc_91 = maybe_doc_75
		has_doc_QMARK__92 = has_doc_QMARK__76
		body_forms_93 = body_forms_78
		multi_QMARK__94 = multi_QMARK__80
		goto b11
	}
b10:
	;
	arg__1242_97, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{name_sym_83})
	if callErr != nil {
		return nil, callErr
	}
	arg__1249_102, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{name_sym_83})
	if callErr != nil {
		return nil, callErr
	}
	f_105, callErr = rt.InvokeValue(rt.LookupVar("ir", "new-fn").Deref(), []vm.Value{arg__1249_102, vm.Int(0), vm.Boolean(false)})
	if callErr != nil {
		return nil, callErr
	}
	entry_blk_107, callErr = rt.InvokeValue(rt.LookupVar("ir", "entry-block").Deref(), []vm.Value{f_105})
	if callErr != nil {
		return nil, callErr
	}
	ctx_109, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "new-context").Deref(), []vm.Value{f_105})
	if callErr != nil {
		return nil, callErr
	}
	arities_111, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "cons").Deref(), []vm.Value{args_vec_81, body_forms_86})
	if callErr != nil {
		return nil, callErr
	}
	expanded_forms_115, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapv").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var arg_vec_2 vm.Value
		var body_forms_4 vm.Value
		var e_6 vm.Value
		var arity_form_7 vm.Value
		var arg_vec_8 vm.Value
		var body_forms_9 vm.Value
		var e_10 vm.Value
		var arity_form_11 vm.Value
		var arg_vec_12 vm.Value
		var body_forms_13 vm.Value
		var e_14 vm.Value
		var v22 vm.Value
		var v24 vm.Value
		var arity_form_25 vm.Value
		var arg_vec_26 vm.Value
		var body_forms_27 vm.Value
		var e_28 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg_vec_2, body_forms_4, e_6, arity_form_7, arg_vec_8, body_forms_9, e_10, arity_form_11, arg_vec_12, body_forms_13, e_14, v22, v24, arity_form_25, arg_vec_26, body_forms_27, e_28
		arg_vec_2, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		body_forms_4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		e_6, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "expand-fn-args").Deref(), []vm.Value{arg_vec_2, body_forms_4})
		if callErr != nil {
			return nil, callErr
		}
		if vm.IsTruthy(e_6) {
			arity_form_7 = arg0
			arg_vec_8 = arg_vec_2
			body_forms_9 = body_forms_4
			e_10 = e_6
			goto b1
		} else {
			arity_form_11 = arg0
			arg_vec_12 = arg_vec_2
			body_forms_13 = body_forms_4
			e_14 = e_6
			goto b2
		}
	b1:
		;
		v24 = e_10
		arity_form_25 = arity_form_7
		arg_vec_26 = arg_vec_8
		body_forms_27 = body_forms_9
		e_28 = e_10
		goto b3
	b2:
		;
		v22, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("variadic?"), vm.FALSE, vm.Keyword("flat-args"), arg_vec_12, vm.Keyword("body"), body_forms_13})
		if callErr != nil {
			return nil, callErr
		}
		v24 = v22
		arity_form_25 = arity_form_11
		arg_vec_26 = arg_vec_12
		body_forms_27 = body_forms_13
		e_28 = e_14
		goto b3
	b3:
		;
		return v24, nil
	}), arities_111})
	if callErr != nil {
		return nil, callErr
	}
	arg__1373_120, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__1445_126, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	all_caps_127, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "reduce").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
		var flat_args_6 vm.Value
		var body_10 vm.Value
		var arg__1387_12 vm.Value
		var arg__1391_15 vm.Value
		var arg_set_16 vm.Value
		var arg__1398_20 vm.Value
		var arg__1406_25 vm.Value
		var frees_26 vm.Value
		var arg__1424_34 vm.Value
		var arg__1442_43 vm.Value
		var v44 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _, _, _, _ = flat_args_6, body_10, arg__1387_12, arg__1391_15, arg_set_16, arg__1398_20, arg__1406_25, frees_26, arg__1424_34, arg__1442_43, v44
		flat_args_6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg1, vm.Keyword("flat-args")})
		if callErr != nil {
			return nil, callErr
		}
		body_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg1, vm.Keyword("body")})
		if callErr != nil {
			return nil, callErr
		}
		arg__1387_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
		if callErr != nil {
			return nil, callErr
		}
		arg__1391_15, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
		if callErr != nil {
			return nil, callErr
		}
		arg_set_16, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__1391_15, flat_args_6})
		if callErr != nil {
			return nil, callErr
		}
		arg__1398_20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "cons").Deref(), []vm.Value{vm.Symbol("do"), body_10})
		if callErr != nil {
			return nil, callErr
		}
		arg__1406_25, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "cons").Deref(), []vm.Value{vm.Symbol("do"), body_10})
		if callErr != nil {
			return nil, callErr
		}
		frees_26, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "free-vars").Deref(), []vm.Value{arg__1406_25, arg_set_16})
		if callErr != nil {
			return nil, callErr
		}
		arg__1424_34, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "filter").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
			var v3 vm.Value
			var callErr error
			_ = v3
			v3, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "lookup-local").Deref(), []vm.Value{ctx_109, arg0})
			if callErr != nil {
				return nil, callErr
			}
			return v3, nil
		}), frees_26})
		if callErr != nil {
			return nil, callErr
		}
		arg__1442_43, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "filter").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
			var v3 vm.Value
			var callErr error
			_ = v3
			v3, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "lookup-local").Deref(), []vm.Value{ctx_109, arg0})
			if callErr != nil {
				return nil, callErr
			}
			return v3, nil
		}), frees_26})
		if callErr != nil {
			return nil, callErr
		}
		v44, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg0, arg__1442_43})
		if callErr != nil {
			return nil, callErr
		}
		return v44, nil
	}), arg__1445_126, expanded_forms_115})
	if callErr != nil {
		return nil, callErr
	}
	arg__1452_131, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "sort-by").Deref(), []vm.Value{rt.LookupVar("clojure.core", "str").Deref(), all_caps_127})
	if callErr != nil {
		return nil, callErr
	}
	arg__1459_136, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "sort-by").Deref(), []vm.Value{rt.LookupVar("clojure.core", "str").Deref(), all_caps_127})
	if callErr != nil {
		return nil, callErr
	}
	captures_137, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg__1459_136})
	if callErr != nil {
		return nil, callErr
	}
	templates_147, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapv").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var flat_args_6 vm.Value
		var body_10 vm.Value
		var variadic_QMARK__14 vm.Value
		var v16 vm.Value
		var callErr error
		_, _, _, _ = flat_args_6, body_10, variadic_QMARK__14, v16
		flat_args_6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg0, vm.Keyword("flat-args")})
		if callErr != nil {
			return nil, callErr
		}
		body_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg0, vm.Keyword("body")})
		if callErr != nil {
			return nil, callErr
		}
		variadic_QMARK__14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg0, vm.Keyword("variadic?")})
		if callErr != nil {
			return nil, callErr
		}
		v16, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-inner-fn-template").Deref(), []vm.Value{name_sym_83, flat_args_6, body_10, captures_137, variadic_QMARK__14})
		if callErr != nil {
			return nil, callErr
		}
		return v16, nil
	}), expanded_forms_115})
	if callErr != nil {
		return nil, callErr
	}
	template_152, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("fns"), templates_147, vm.Keyword("kind"), vm.Keyword("multi-fn-template")})
	if callErr != nil {
		return nil, callErr
	}
	closure_id_154, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "emit-template-closure").Deref(), []vm.Value{template_152, captures_137, ctx_109})
	if callErr != nil {
		return nil, callErr
	}
	final_blk_156, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_109})
	if callErr != nil {
		return nil, callErr
	}
	arg__1537_159, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{closure_id_154})
	if callErr != nil {
		return nil, callErr
	}
	arg__1545_164, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{closure_id_154})
	if callErr != nil {
		return nil, callErr
	}
	v166, callErr = rt.InvokeValue(rt.LookupVar("ir", "add-terminator!").Deref(), []vm.Value{f_105, final_blk_156, vm.Keyword("return"), arg__1545_164, vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	v691 = f_105
	args_vec_692 = args_vec_81
	defn_form_693 = defn_form_82
	name_sym_694 = name_sym_83
	maybe_doc_695 = maybe_doc_84
	has_doc_QMARK__696 = has_doc_QMARK__85
	body_forms_697 = body_forms_86
	multi_QMARK__698 = multi_QMARK__87
	goto b12
b11:
	;
	expanded_169, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "expand-fn-args").Deref(), []vm.Value{args_vec_88, body_forms_93})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(expanded_169) {
		args_vec_170 = args_vec_88
		defn_form_171 = defn_form_89
		name_sym_172 = name_sym_90
		maybe_doc_173 = maybe_doc_91
		has_doc_QMARK__174 = has_doc_QMARK__92
		body_forms_175 = body_forms_93
		multi_QMARK__176 = multi_QMARK__94
		expanded_177 = expanded_169
		goto b13
	} else {
		args_vec_178 = args_vec_88
		defn_form_179 = defn_form_89
		name_sym_180 = name_sym_90
		maybe_doc_181 = maybe_doc_91
		has_doc_QMARK__182 = has_doc_QMARK__92
		body_forms_183 = body_forms_93
		multi_QMARK__184 = multi_QMARK__94
		expanded_185 = expanded_169
		goto b14
	}
b12:
	;
	return v691, nil
b13:
	;
	v188, callErr = rt.InvokeValue(vm.Keyword("flat-args"), []vm.Value{expanded_177})
	if callErr != nil {
		return nil, callErr
	}
	flat_args_191 = v188
	args_vec_192 = args_vec_170
	defn_form_193 = defn_form_171
	name_sym_194 = name_sym_172
	maybe_doc_195 = maybe_doc_173
	has_doc_QMARK__196 = has_doc_QMARK__174
	body_forms_197 = body_forms_175
	multi_QMARK__198 = multi_QMARK__176
	expanded_199 = expanded_177
	goto b15
b14:
	;
	flat_args_191 = args_vec_178
	args_vec_192 = args_vec_178
	defn_form_193 = defn_form_179
	name_sym_194 = name_sym_180
	maybe_doc_195 = maybe_doc_181
	has_doc_QMARK__196 = has_doc_QMARK__182
	body_forms_197 = body_forms_183
	multi_QMARK__198 = multi_QMARK__184
	expanded_199 = expanded_185
	goto b15
b15:
	;
	if vm.IsTruthy(expanded_199) {
		flat_args_200 = flat_args_191
		args_vec_201 = args_vec_192
		defn_form_202 = defn_form_193
		name_sym_203 = name_sym_194
		maybe_doc_204 = maybe_doc_195
		has_doc_QMARK__205 = has_doc_QMARK__196
		body_forms_206 = body_forms_197
		multi_QMARK__207 = multi_QMARK__198
		expanded_208 = expanded_199
		goto b16
	} else {
		flat_args_209 = flat_args_191
		args_vec_210 = args_vec_192
		defn_form_211 = defn_form_193
		name_sym_212 = name_sym_194
		maybe_doc_213 = maybe_doc_195
		has_doc_QMARK__214 = has_doc_QMARK__196
		body_forms_215 = body_forms_197
		multi_QMARK__216 = multi_QMARK__198
		expanded_217 = expanded_199
		goto b17
	}
b16:
	;
	v220, callErr = rt.InvokeValue(vm.Keyword("body"), []vm.Value{expanded_208})
	if callErr != nil {
		return nil, callErr
	}
	flat_body_223 = v220
	flat_args_224 = flat_args_200
	args_vec_225 = args_vec_201
	defn_form_226 = defn_form_202
	name_sym_227 = name_sym_203
	maybe_doc_228 = maybe_doc_204
	has_doc_QMARK__229 = has_doc_QMARK__205
	body_forms_230 = body_forms_206
	multi_QMARK__231 = multi_QMARK__207
	expanded_232 = expanded_208
	goto b18
b17:
	;
	flat_body_223 = body_forms_215
	flat_args_224 = flat_args_209
	args_vec_225 = args_vec_210
	defn_form_226 = defn_form_211
	name_sym_227 = name_sym_212
	maybe_doc_228 = maybe_doc_213
	has_doc_QMARK__229 = has_doc_QMARK__214
	body_forms_230 = body_forms_215
	multi_QMARK__231 = multi_QMARK__216
	expanded_232 = expanded_217
	goto b18
b18:
	;
	if vm.IsTruthy(expanded_232) {
		flat_body_233 = flat_body_223
		flat_args_234 = flat_args_224
		args_vec_235 = args_vec_225
		defn_form_236 = defn_form_226
		name_sym_237 = name_sym_227
		maybe_doc_238 = maybe_doc_228
		has_doc_QMARK__239 = has_doc_QMARK__229
		body_forms_240 = body_forms_230
		multi_QMARK__241 = multi_QMARK__231
		expanded_242 = expanded_232
		goto b19
	} else {
		flat_body_243 = flat_body_223
		flat_args_244 = flat_args_224
		args_vec_245 = args_vec_225
		defn_form_246 = defn_form_226
		name_sym_247 = name_sym_227
		maybe_doc_248 = maybe_doc_228
		has_doc_QMARK__249 = has_doc_QMARK__229
		body_forms_250 = body_forms_230
		multi_QMARK__251 = multi_QMARK__231
		expanded_252 = expanded_232
		goto b20
	}
b19:
	;
	v255, callErr = rt.InvokeValue(vm.Keyword("variadic?"), []vm.Value{expanded_242})
	if callErr != nil {
		return nil, callErr
	}
	variadic_QMARK__259 = v255
	flat_body_260 = flat_body_233
	flat_args_261 = flat_args_234
	args_vec_262 = args_vec_235
	defn_form_263 = defn_form_236
	name_sym_264 = name_sym_237
	maybe_doc_265 = maybe_doc_238
	has_doc_QMARK__266 = has_doc_QMARK__239
	body_forms_267 = body_forms_240
	multi_QMARK__268 = multi_QMARK__241
	expanded_269 = expanded_242
	goto b21
b20:
	;
	variadic_QMARK__259 = vm.Boolean(false)
	flat_body_260 = flat_body_243
	flat_args_261 = flat_args_244
	args_vec_262 = args_vec_245
	defn_form_263 = defn_form_246
	name_sym_264 = name_sym_247
	maybe_doc_265 = maybe_doc_248
	has_doc_QMARK__266 = has_doc_QMARK__249
	body_forms_267 = body_forms_250
	multi_QMARK__268 = multi_QMARK__251
	expanded_269 = expanded_252
	goto b21
b21:
	;
	arity_271, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{flat_args_261})
	if callErr != nil {
		return nil, callErr
	}
	arg__1564_273, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{name_sym_264})
	if callErr != nil {
		return nil, callErr
	}
	arg__1571_276, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{name_sym_264})
	if callErr != nil {
		return nil, callErr
	}
	f_277, callErr = rt.InvokeValue(rt.LookupVar("ir", "new-fn").Deref(), []vm.Value{arg__1571_276, arity_271, variadic_QMARK__259})
	if callErr != nil {
		return nil, callErr
	}
	entry_blk_279, callErr = rt.InvokeValue(rt.LookupVar("ir", "entry-block").Deref(), []vm.Value{f_277})
	if callErr != nil {
		return nil, callErr
	}
	ctx_281, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "new-context").Deref(), []vm.Value{f_277})
	if callErr != nil {
		return nil, callErr
	}
	arg__1584_283, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_281})
	if callErr != nil {
		return nil, callErr
	}
	arg__1589_286, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{flat_args_261})
	if callErr != nil {
		return nil, callErr
	}
	arg__1594_289, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_281})
	if callErr != nil {
		return nil, callErr
	}
	arg__1599_292, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{flat_args_261})
	if callErr != nil {
		return nil, callErr
	}
	arg__1600_293, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "assoc").Deref(), []vm.Value{arg__1594_289, vm.Keyword("fn-arg-syms"), arg__1599_292})
	if callErr != nil {
		return nil, callErr
	}
	arg__1606_296, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_281})
	if callErr != nil {
		return nil, callErr
	}
	arg__1611_299, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{flat_args_261})
	if callErr != nil {
		return nil, callErr
	}
	arg__1616_302, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_281})
	if callErr != nil {
		return nil, callErr
	}
	arg__1621_305, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{flat_args_261})
	if callErr != nil {
		return nil, callErr
	}
	arg__1622_306, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "assoc").Deref(), []vm.Value{arg__1616_302, vm.Keyword("fn-arg-syms"), arg__1621_305})
	if callErr != nil {
		return nil, callErr
	}
	v307, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "reset!").Deref(), []vm.Value{ctx_281, arg__1622_306})
	if callErr != nil {
		return nil, callErr
	}
	i_308 = 0
	arity_309 = arity_271
	flat_args_310 = flat_args_261
	ctx_311 = ctx_281
	entry_blk_312 = entry_blk_279
	f_313 = f_277
	v704 = vm.Keyword("load-arg")
	v707 = vm.NewArrayVector([]vm.Value{})
	goto b22
b22:
	;
	v348 = rt.LtValue(vm.Int(i_308), arity_309)
	if v348 {
		variadic_QMARK__316 = variadic_QMARK__259
		flat_body_317 = flat_body_260
		args_vec_318 = args_vec_262
		defn_form_319 = defn_form_263
		name_sym_320 = name_sym_264
		maybe_doc_321 = maybe_doc_265
		has_doc_QMARK__322 = has_doc_QMARK__266
		body_forms_323 = body_forms_267
		multi_QMARK__324 = multi_QMARK__268
		expanded_325 = expanded_269
		i_326 = i_308
		arity_327 = arity_309
		flat_args_328 = flat_args_310
		ctx_329 = ctx_311
		entry_blk_330 = entry_blk_312
		f_331 = f_313
		v705 = v704
		v708 = v707
		goto b23
	} else {
		variadic_QMARK__332 = variadic_QMARK__259
		flat_body_333 = flat_body_260
		args_vec_334 = args_vec_262
		defn_form_335 = defn_form_263
		name_sym_336 = name_sym_264
		maybe_doc_337 = maybe_doc_265
		has_doc_QMARK__338 = has_doc_QMARK__266
		body_forms_339 = body_forms_267
		multi_QMARK__340 = multi_QMARK__268
		expanded_341 = expanded_269
		i_342 = i_308
		arity_343 = arity_309
		flat_args_344 = flat_args_310
		ctx_345 = ctx_311
		entry_blk_346 = entry_blk_312
		f_347 = f_313
		v706 = v704
		v709 = v707
		goto b24
	}
b23:
	;
	arg_id_355, callErr = rt.InvokeValue(rt.LookupVar("ir", "add-inst").Deref(), []vm.Value{f_331, entry_blk_330, v705, v708, vm.Int(i_326)})
	if callErr != nil {
		return nil, callErr
	}
	arg__1642_357, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{flat_args_328, vm.Int(i_326)})
	if callErr != nil {
		return nil, callErr
	}
	arg__1651_360, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{flat_args_328, vm.Int(i_326)})
	if callErr != nil {
		return nil, callErr
	}
	v361, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "bind-local!").Deref(), []vm.Value{ctx_329, arg__1651_360, arg_id_355})
	if callErr != nil {
		return nil, callErr
	}
	v362 = i_326 + 1
	i_308 = v362
	arity_309 = arity_327
	flat_args_310 = flat_args_328
	ctx_311 = ctx_329
	entry_blk_312 = entry_blk_330
	f_313 = f_331
	v704 = v705
	v707 = v708
	goto b22
b24:
	;
	v366 = vm.NIL
	variadic_QMARK__367 = variadic_QMARK__332
	flat_body_368 = flat_body_333
	args_vec_369 = args_vec_334
	defn_form_370 = defn_form_335
	name_sym_371 = name_sym_336
	maybe_doc_372 = maybe_doc_337
	has_doc_QMARK__373 = has_doc_QMARK__338
	body_forms_374 = body_forms_339
	multi_QMARK__375 = multi_QMARK__340
	expanded_376 = expanded_341
	i_377 = i_342
	arity_378 = arity_343
	flat_args_379 = flat_args_344
	ctx_380 = ctx_345
	entry_blk_381 = entry_blk_346
	f_382 = f_347
	goto b25
b25:
	;
	fs_383 = flat_body_368
	last_id_384 = vm.NIL
	ctx_385 = ctx_380
	goto b26
b26:
	;
	v425, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{fs_383})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v425) {
		variadic_QMARK__388 = variadic_QMARK__367
		flat_body_389 = flat_body_368
		args_vec_390 = args_vec_369
		defn_form_391 = defn_form_370
		name_sym_392 = name_sym_371
		maybe_doc_393 = maybe_doc_372
		has_doc_QMARK__394 = has_doc_QMARK__373
		body_forms_395 = body_forms_374
		multi_QMARK__396 = multi_QMARK__375
		expanded_397 = expanded_376
		i_398 = i_377
		arity_399 = arity_378
		flat_args_400 = flat_args_379
		entry_blk_401 = entry_blk_381
		f_402 = f_382
		fs_403 = fs_383
		last_id_404 = last_id_384
		ctx_405 = ctx_385
		goto b27
	} else {
		variadic_QMARK__406 = variadic_QMARK__367
		flat_body_407 = flat_body_368
		args_vec_408 = args_vec_369
		defn_form_409 = defn_form_370
		name_sym_410 = name_sym_371
		maybe_doc_411 = maybe_doc_372
		has_doc_QMARK__412 = has_doc_QMARK__373
		body_forms_413 = body_forms_374
		multi_QMARK__414 = multi_QMARK__375
		expanded_415 = expanded_376
		i_416 = i_377
		arity_417 = arity_378
		flat_args_418 = flat_args_379
		entry_blk_419 = entry_blk_381
		f_420 = f_382
		fs_421 = fs_383
		last_id_422 = last_id_384
		ctx_423 = ctx_385
		goto b28
	}
b27:
	;
	v428, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{fs_403})
	if callErr != nil {
		return nil, callErr
	}
	arg__1663_430, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{fs_403})
	if callErr != nil {
		return nil, callErr
	}
	arg__1669_433, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{fs_403})
	if callErr != nil {
		return nil, callErr
	}
	v434, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-form").Deref(), []vm.Value{arg__1669_433, ctx_405})
	if callErr != nil {
		return nil, callErr
	}
	fs_383 = v428
	last_id_384 = v434
	ctx_385 = ctx_405
	goto b26
b28:
	;
	last_val_437 = last_id_422
	variadic_QMARK__438 = variadic_QMARK__406
	flat_body_439 = flat_body_407
	args_vec_440 = args_vec_408
	defn_form_441 = defn_form_409
	name_sym_442 = name_sym_410
	maybe_doc_443 = maybe_doc_411
	has_doc_QMARK__444 = has_doc_QMARK__412
	body_forms_445 = body_forms_413
	multi_QMARK__446 = multi_QMARK__414
	expanded_447 = expanded_415
	i_448 = i_416
	arity_449 = arity_417
	flat_args_450 = flat_args_418
	entry_blk_451 = entry_blk_419
	f_452 = f_420
	fs_453 = fs_421
	last_id_454 = last_id_422
	ctx_455 = ctx_423
	goto b29
b29:
	;
	final_blk_457, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_455})
	if callErr != nil {
		return nil, callErr
	}
	v499, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "terminated?").Deref(), []vm.Value{last_val_437})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v499) {
		last_val_458 = last_val_437
		variadic_QMARK__459 = variadic_QMARK__438
		flat_body_460 = flat_body_439
		args_vec_461 = args_vec_440
		defn_form_462 = defn_form_441
		name_sym_463 = name_sym_442
		maybe_doc_464 = maybe_doc_443
		has_doc_QMARK__465 = has_doc_QMARK__444
		body_forms_466 = body_forms_445
		multi_QMARK__467 = multi_QMARK__446
		expanded_468 = expanded_447
		i_469 = i_448
		arity_470 = arity_449
		flat_args_471 = flat_args_450
		entry_blk_472 = entry_blk_451
		f_473 = f_452
		fs_474 = fs_453
		last_id_475 = last_id_454
		ctx_476 = ctx_455
		final_blk_477 = final_blk_457
		goto b30
	} else {
		last_val_478 = last_val_437
		variadic_QMARK__479 = variadic_QMARK__438
		flat_body_480 = flat_body_439
		args_vec_481 = args_vec_440
		defn_form_482 = defn_form_441
		name_sym_483 = name_sym_442
		maybe_doc_484 = maybe_doc_443
		has_doc_QMARK__485 = has_doc_QMARK__444
		body_forms_486 = body_forms_445
		multi_QMARK__487 = multi_QMARK__446
		expanded_488 = expanded_447
		i_489 = i_448
		arity_490 = arity_449
		flat_args_491 = flat_args_450
		entry_blk_492 = entry_blk_451
		f_493 = f_452
		fs_494 = fs_453
		last_id_495 = last_id_454
		ctx_496 = ctx_455
		final_blk_497 = final_blk_457
		goto b31
	}
b30:
	;
	v669 = vm.NIL
	last_val_670 = last_val_458
	variadic_QMARK__671 = variadic_QMARK__459
	flat_body_672 = flat_body_460
	args_vec_673 = args_vec_461
	defn_form_674 = defn_form_462
	name_sym_675 = name_sym_463
	maybe_doc_676 = maybe_doc_464
	has_doc_QMARK__677 = has_doc_QMARK__465
	body_forms_678 = body_forms_466
	multi_QMARK__679 = multi_QMARK__467
	expanded_680 = expanded_468
	i_681 = i_469
	arity_682 = arity_470
	flat_args_683 = flat_args_471
	entry_blk_684 = entry_blk_472
	f_685 = f_473
	fs_686 = fs_474
	last_id_687 = last_id_475
	ctx_688 = ctx_476
	final_blk_689 = final_blk_477
	goto b32
b31:
	;
	v551, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{last_val_478})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v551) {
		last_val_504 = last_val_478
		variadic_QMARK__505 = variadic_QMARK__479
		flat_body_506 = flat_body_480
		args_vec_507 = args_vec_481
		defn_form_508 = defn_form_482
		name_sym_509 = name_sym_483
		maybe_doc_510 = maybe_doc_484
		has_doc_QMARK__511 = has_doc_QMARK__485
		body_forms_512 = body_forms_486
		multi_QMARK__513 = multi_QMARK__487
		expanded_514 = expanded_488
		i_515 = i_489
		arity_516 = arity_490
		flat_args_517 = flat_args_491
		entry_blk_518 = entry_blk_492
		arg__1677_519 = f_493
		f_520 = f_493
		fs_521 = fs_494
		last_id_522 = last_id_495
		ctx_523 = ctx_496
		final_blk_524 = final_blk_497
		arg__1678_525 = final_blk_497
		arg__1679_526 = vm.Keyword("return")
		goto b33
	} else {
		last_val_527 = last_val_478
		variadic_QMARK__528 = variadic_QMARK__479
		flat_body_529 = flat_body_480
		args_vec_530 = args_vec_481
		defn_form_531 = defn_form_482
		name_sym_532 = name_sym_483
		maybe_doc_533 = maybe_doc_484
		has_doc_QMARK__534 = has_doc_QMARK__485
		body_forms_535 = body_forms_486
		multi_QMARK__536 = multi_QMARK__487
		expanded_537 = expanded_488
		i_538 = i_489
		arity_539 = arity_490
		flat_args_540 = flat_args_491
		entry_blk_541 = entry_blk_492
		arg__1677_542 = f_493
		f_543 = f_493
		fs_544 = fs_494
		last_id_545 = last_id_495
		ctx_546 = ctx_496
		final_blk_547 = final_blk_497
		arg__1678_548 = final_blk_497
		arg__1679_549 = vm.Keyword("return")
		goto b34
	}
b32:
	;
	v691 = f_685
	args_vec_692 = args_vec_673
	defn_form_693 = defn_form_674
	name_sym_694 = name_sym_675
	maybe_doc_695 = maybe_doc_676
	has_doc_QMARK__696 = has_doc_QMARK__677
	body_forms_697 = body_forms_678
	multi_QMARK__698 = multi_QMARK__679
	goto b12
b33:
	;
	arg__1685_558 = vm.NewArrayVector([]vm.Value{})
	last_val_559 = last_val_504
	variadic_QMARK__560 = variadic_QMARK__505
	flat_body_561 = flat_body_506
	args_vec_562 = args_vec_507
	defn_form_563 = defn_form_508
	name_sym_564 = name_sym_509
	maybe_doc_565 = maybe_doc_510
	has_doc_QMARK__566 = has_doc_QMARK__511
	body_forms_567 = body_forms_512
	multi_QMARK__568 = multi_QMARK__513
	expanded_569 = expanded_514
	i_570 = i_515
	arity_571 = arity_516
	flat_args_572 = flat_args_517
	entry_blk_573 = entry_blk_518
	arg__1677_574 = arg__1677_519
	f_575 = f_520
	fs_576 = fs_521
	last_id_577 = last_id_522
	ctx_578 = ctx_523
	final_blk_579 = final_blk_524
	arg__1678_580 = arg__1678_525
	arg__1679_581 = arg__1679_526
	goto b35
b34:
	;
	v556, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{last_val_527})
	if callErr != nil {
		return nil, callErr
	}
	arg__1685_558 = v556
	last_val_559 = last_val_527
	variadic_QMARK__560 = variadic_QMARK__528
	flat_body_561 = flat_body_529
	args_vec_562 = args_vec_530
	defn_form_563 = defn_form_531
	name_sym_564 = name_sym_532
	maybe_doc_565 = maybe_doc_533
	has_doc_QMARK__566 = has_doc_QMARK__534
	body_forms_567 = body_forms_535
	multi_QMARK__568 = multi_QMARK__536
	expanded_569 = expanded_537
	i_570 = i_538
	arity_571 = arity_539
	flat_args_572 = flat_args_540
	entry_blk_573 = entry_blk_541
	arg__1677_574 = arg__1677_542
	f_575 = f_543
	fs_576 = fs_544
	last_id_577 = last_id_545
	ctx_578 = ctx_546
	final_blk_579 = final_blk_547
	arg__1678_580 = arg__1678_548
	arg__1679_581 = arg__1679_549
	goto b35
b35:
	;
	v634, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{last_val_559})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v634) {
		last_val_585 = last_val_559
		variadic_QMARK__586 = variadic_QMARK__560
		flat_body_587 = flat_body_561
		args_vec_588 = args_vec_562
		defn_form_589 = defn_form_563
		name_sym_590 = name_sym_564
		maybe_doc_591 = maybe_doc_565
		has_doc_QMARK__592 = has_doc_QMARK__566
		body_forms_593 = body_forms_567
		multi_QMARK__594 = multi_QMARK__568
		expanded_595 = expanded_569
		i_596 = i_570
		arity_597 = arity_571
		flat_args_598 = flat_args_572
		entry_blk_599 = entry_blk_573
		f_600 = f_575
		arg__1688_601 = f_575
		fs_602 = fs_576
		last_id_603 = last_id_577
		ctx_604 = ctx_578
		final_blk_605 = final_blk_579
		arg__1689_606 = final_blk_579
		head__1687_607 = rt.LookupVar("ir", "add-terminator!").Deref()
		arg__1690_608 = vm.Keyword("return")
		goto b36
	} else {
		last_val_609 = last_val_559
		variadic_QMARK__610 = variadic_QMARK__560
		flat_body_611 = flat_body_561
		args_vec_612 = args_vec_562
		defn_form_613 = defn_form_563
		name_sym_614 = name_sym_564
		maybe_doc_615 = maybe_doc_565
		has_doc_QMARK__616 = has_doc_QMARK__566
		body_forms_617 = body_forms_567
		multi_QMARK__618 = multi_QMARK__568
		expanded_619 = expanded_569
		i_620 = i_570
		arity_621 = arity_571
		flat_args_622 = flat_args_572
		entry_blk_623 = entry_blk_573
		f_624 = f_575
		arg__1688_625 = f_575
		fs_626 = fs_576
		last_id_627 = last_id_577
		ctx_628 = ctx_578
		final_blk_629 = final_blk_579
		arg__1689_630 = final_blk_579
		head__1687_631 = rt.LookupVar("ir", "add-terminator!").Deref()
		arg__1690_632 = vm.Keyword("return")
		goto b37
	}
b36:
	;
	arg__1696_641 = vm.NewArrayVector([]vm.Value{})
	last_val_642 = last_val_585
	variadic_QMARK__643 = variadic_QMARK__586
	flat_body_644 = flat_body_587
	args_vec_645 = args_vec_588
	defn_form_646 = defn_form_589
	name_sym_647 = name_sym_590
	maybe_doc_648 = maybe_doc_591
	has_doc_QMARK__649 = has_doc_QMARK__592
	body_forms_650 = body_forms_593
	multi_QMARK__651 = multi_QMARK__594
	expanded_652 = expanded_595
	i_653 = i_596
	arity_654 = arity_597
	flat_args_655 = flat_args_598
	entry_blk_656 = entry_blk_599
	f_657 = f_600
	arg__1688_658 = arg__1688_601
	fs_659 = fs_602
	last_id_660 = last_id_603
	ctx_661 = ctx_604
	final_blk_662 = final_blk_605
	arg__1689_663 = arg__1689_606
	head__1687_664 = head__1687_607
	arg__1690_665 = arg__1690_608
	goto b38
b37:
	;
	v639, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{last_val_609})
	if callErr != nil {
		return nil, callErr
	}
	arg__1696_641 = v639
	last_val_642 = last_val_609
	variadic_QMARK__643 = variadic_QMARK__610
	flat_body_644 = flat_body_611
	args_vec_645 = args_vec_612
	defn_form_646 = defn_form_613
	name_sym_647 = name_sym_614
	maybe_doc_648 = maybe_doc_615
	has_doc_QMARK__649 = has_doc_QMARK__616
	body_forms_650 = body_forms_617
	multi_QMARK__651 = multi_QMARK__618
	expanded_652 = expanded_619
	i_653 = i_620
	arity_654 = arity_621
	flat_args_655 = flat_args_622
	entry_blk_656 = entry_blk_623
	f_657 = f_624
	arg__1688_658 = arg__1688_625
	fs_659 = fs_626
	last_id_660 = last_id_627
	ctx_661 = ctx_628
	final_blk_662 = final_blk_629
	arg__1689_663 = arg__1689_630
	head__1687_664 = head__1687_631
	arg__1690_665 = arg__1690_632
	goto b38
b38:
	;
	v667, callErr = rt.InvokeValue(head__1687_664, []vm.Value{arg__1688_658, arg__1689_663, arg__1690_665, arg__1696_641, vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	v669 = v667
	last_val_670 = last_val_642
	variadic_QMARK__671 = variadic_QMARK__643
	flat_body_672 = flat_body_644
	args_vec_673 = args_vec_645
	defn_form_674 = defn_form_646
	name_sym_675 = name_sym_647
	maybe_doc_676 = maybe_doc_648
	has_doc_QMARK__677 = has_doc_QMARK__649
	body_forms_678 = body_forms_650
	multi_QMARK__679 = multi_QMARK__651
	expanded_680 = expanded_652
	i_681 = i_653
	arity_682 = arity_654
	flat_args_683 = flat_args_655
	entry_blk_684 = entry_blk_656
	f_685 = f_657
	fs_686 = fs_659
	last_id_687 = last_id_660
	ctx_688 = ctx_661
	final_blk_689 = final_blk_662
	goto b32
}
func build_fn_STAR_(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var __8 vm.Value
	var maybe_name_14 vm.Value
	var raw_rest_18 vm.Value
	var has_name_QMARK__20 vm.Value
	var form_21 vm.Value
	var vec__1698_22 vm.Value
	var ctx_23 vm.Value
	var __24 vm.Value
	var maybe_name_25 vm.Value
	var raw_rest_26 vm.Value
	var has_name_QMARK__27 vm.Value
	var form_28 vm.Value
	var vec__1698_29 vm.Value
	var ctx_30 vm.Value
	var __31 vm.Value
	var maybe_name_32 vm.Value
	var raw_rest_33 vm.Value
	var has_name_QMARK__34 vm.Value
	var name_sym_39 vm.Value
	var form_40 vm.Value
	var vec__1698_41 vm.Value
	var ctx_42 vm.Value
	var __43 vm.Value
	var maybe_name_44 vm.Value
	var raw_rest_45 vm.Value
	var has_name_QMARK__46 vm.Value
	var name_sym_47 vm.Value
	var form_48 vm.Value
	var vec__1698_49 vm.Value
	var ctx_50 vm.Value
	var __51 vm.Value
	var maybe_name_52 vm.Value
	var raw_rest_53 vm.Value
	var has_name_QMARK__54 vm.Value
	var name_sym_55 vm.Value
	var form_56 vm.Value
	var vec__1698_57 vm.Value
	var ctx_58 vm.Value
	var __59 vm.Value
	var maybe_name_60 vm.Value
	var raw_rest_61 vm.Value
	var has_name_QMARK__62 vm.Value
	var v66 vm.Value
	var rest_forms_68 vm.Value
	var name_sym_69 vm.Value
	var form_70 vm.Value
	var vec__1698_71 vm.Value
	var ctx_72 vm.Value
	var __73 vm.Value
	var maybe_name_74 vm.Value
	var raw_rest_75 vm.Value
	var has_name_QMARK__76 vm.Value
	var and__x_78 vm.Value
	var rest_forms_79 vm.Value
	var name_sym_80 vm.Value
	var form_81 vm.Value
	var vec__1698_82 vm.Value
	var ctx_83 vm.Value
	var __84 vm.Value
	var maybe_name_85 vm.Value
	var raw_rest_86 vm.Value
	var has_name_QMARK__87 vm.Value
	var and__x_88 vm.Value
	var arg__1734_101 vm.Value
	var arg__1739_104 vm.Value
	var v105 vm.Value
	var rest_forms_89 vm.Value
	var name_sym_90 vm.Value
	var form_91 vm.Value
	var vec__1698_92 vm.Value
	var ctx_93 vm.Value
	var __94 vm.Value
	var maybe_name_95 vm.Value
	var raw_rest_96 vm.Value
	var has_name_QMARK__97 vm.Value
	var and__x_98 vm.Value
	var multi_QMARK__108 vm.Value
	var rest_forms_109 vm.Value
	var name_sym_110 vm.Value
	var form_111 vm.Value
	var vec__1698_112 vm.Value
	var ctx_113 vm.Value
	var __114 vm.Value
	var maybe_name_115 vm.Value
	var raw_rest_116 vm.Value
	var has_name_QMARK__117 vm.Value
	var and__x_118 vm.Value
	var multi_QMARK__119 vm.Value
	var rest_forms_120 vm.Value
	var name_sym_121 vm.Value
	var form_122 vm.Value
	var vec__1698_123 vm.Value
	var ctx_124 vm.Value
	var __125 vm.Value
	var maybe_name_126 vm.Value
	var raw_rest_127 vm.Value
	var has_name_QMARK__128 vm.Value
	var expanded_forms_143 vm.Value
	var arg__1850_148 vm.Value
	var arg__1922_154 vm.Value
	var all_caps_155 vm.Value
	var arg__1929_159 vm.Value
	var arg__1936_164 vm.Value
	var captures_165 vm.Value
	var templates_175 vm.Value
	var template_180 vm.Value
	var v182 vm.Value
	var multi_QMARK__129 vm.Value
	var rest_forms_130 vm.Value
	var name_sym_131 vm.Value
	var form_132 vm.Value
	var vec__1698_133 vm.Value
	var ctx_134 vm.Value
	var __135 vm.Value
	var maybe_name_136 vm.Value
	var raw_rest_137 vm.Value
	var has_name_QMARK__138 vm.Value
	var args_vec_185 vm.Value
	var body_forms_187 vm.Value
	var v189 vm.Value
	var v191 vm.Value
	var multi_QMARK__192 vm.Value
	var rest_forms_193 vm.Value
	var name_sym_194 vm.Value
	var form_195 vm.Value
	var vec__1698_196 vm.Value
	var ctx_197 vm.Value
	var __198 vm.Value
	var maybe_name_199 vm.Value
	var raw_rest_200 vm.Value
	var has_name_QMARK__201 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = __8, maybe_name_14, raw_rest_18, has_name_QMARK__20, form_21, vec__1698_22, ctx_23, __24, maybe_name_25, raw_rest_26, has_name_QMARK__27, form_28, vec__1698_29, ctx_30, __31, maybe_name_32, raw_rest_33, has_name_QMARK__34, name_sym_39, form_40, vec__1698_41, ctx_42, __43, maybe_name_44, raw_rest_45, has_name_QMARK__46, name_sym_47, form_48, vec__1698_49, ctx_50, __51, maybe_name_52, raw_rest_53, has_name_QMARK__54, name_sym_55, form_56, vec__1698_57, ctx_58, __59, maybe_name_60, raw_rest_61, has_name_QMARK__62, v66, rest_forms_68, name_sym_69, form_70, vec__1698_71, ctx_72, __73, maybe_name_74, raw_rest_75, has_name_QMARK__76, and__x_78, rest_forms_79, name_sym_80, form_81, vec__1698_82, ctx_83, __84, maybe_name_85, raw_rest_86, has_name_QMARK__87, and__x_88, arg__1734_101, arg__1739_104, v105, rest_forms_89, name_sym_90, form_91, vec__1698_92, ctx_93, __94, maybe_name_95, raw_rest_96, has_name_QMARK__97, and__x_98, multi_QMARK__108, rest_forms_109, name_sym_110, form_111, vec__1698_112, ctx_113, __114, maybe_name_115, raw_rest_116, has_name_QMARK__117, and__x_118, multi_QMARK__119, rest_forms_120, name_sym_121, form_122, vec__1698_123, ctx_124, __125, maybe_name_126, raw_rest_127, has_name_QMARK__128, expanded_forms_143, arg__1850_148, arg__1922_154, all_caps_155, arg__1929_159, arg__1936_164, captures_165, templates_175, template_180, v182, multi_QMARK__129, rest_forms_130, name_sym_131, form_132, vec__1698_133, ctx_134, __135, maybe_name_136, raw_rest_137, has_name_QMARK__138, args_vec_185, body_forms_187, v189, v191, multi_QMARK__192, rest_forms_193, name_sym_194, form_195, vec__1698_196, ctx_197, __198, maybe_name_199, raw_rest_200, has_name_QMARK__201
	__8, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(0), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	maybe_name_14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(1), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	raw_rest_18, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "drop").Deref(), []vm.Value{vm.Int(2), arg0})
	if callErr != nil {
		return nil, callErr
	}
	has_name_QMARK__20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "symbol?").Deref(), []vm.Value{maybe_name_14})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(has_name_QMARK__20) {
		form_21 = arg0
		vec__1698_22 = arg0
		ctx_23 = arg1
		__24 = __8
		maybe_name_25 = maybe_name_14
		raw_rest_26 = raw_rest_18
		has_name_QMARK__27 = has_name_QMARK__20
		goto b1
	} else {
		form_28 = arg0
		vec__1698_29 = arg0
		ctx_30 = arg1
		__31 = __8
		maybe_name_32 = maybe_name_14
		raw_rest_33 = raw_rest_18
		has_name_QMARK__34 = has_name_QMARK__20
		goto b2
	}
b1:
	;
	name_sym_39 = maybe_name_25
	form_40 = form_21
	vec__1698_41 = vec__1698_22
	ctx_42 = ctx_23
	__43 = __24
	maybe_name_44 = maybe_name_25
	raw_rest_45 = raw_rest_26
	has_name_QMARK__46 = has_name_QMARK__27
	goto b3
b2:
	;
	name_sym_39 = vm.String("fn*")
	form_40 = form_28
	vec__1698_41 = vec__1698_29
	ctx_42 = ctx_30
	__43 = __31
	maybe_name_44 = maybe_name_32
	raw_rest_45 = raw_rest_33
	has_name_QMARK__46 = has_name_QMARK__34
	goto b3
b3:
	;
	if vm.IsTruthy(has_name_QMARK__46) {
		name_sym_47 = name_sym_39
		form_48 = form_40
		vec__1698_49 = vec__1698_41
		ctx_50 = ctx_42
		__51 = __43
		maybe_name_52 = maybe_name_44
		raw_rest_53 = raw_rest_45
		has_name_QMARK__54 = has_name_QMARK__46
		goto b4
	} else {
		name_sym_55 = name_sym_39
		form_56 = form_40
		vec__1698_57 = vec__1698_41
		ctx_58 = ctx_42
		__59 = __43
		maybe_name_60 = maybe_name_44
		raw_rest_61 = raw_rest_45
		has_name_QMARK__62 = has_name_QMARK__46
		goto b5
	}
b4:
	;
	rest_forms_68 = raw_rest_53
	name_sym_69 = name_sym_47
	form_70 = form_48
	vec__1698_71 = vec__1698_49
	ctx_72 = ctx_50
	__73 = __51
	maybe_name_74 = maybe_name_52
	raw_rest_75 = raw_rest_53
	has_name_QMARK__76 = has_name_QMARK__54
	goto b6
b5:
	;
	v66, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "cons").Deref(), []vm.Value{maybe_name_60, raw_rest_61})
	if callErr != nil {
		return nil, callErr
	}
	rest_forms_68 = v66
	name_sym_69 = name_sym_55
	form_70 = form_56
	vec__1698_71 = vec__1698_57
	ctx_72 = ctx_58
	__73 = __59
	maybe_name_74 = maybe_name_60
	raw_rest_75 = raw_rest_61
	has_name_QMARK__76 = has_name_QMARK__62
	goto b6
b6:
	;
	and__x_78, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{rest_forms_68})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(and__x_78) {
		rest_forms_79 = rest_forms_68
		name_sym_80 = name_sym_69
		form_81 = form_70
		vec__1698_82 = vec__1698_71
		ctx_83 = ctx_72
		__84 = __73
		maybe_name_85 = maybe_name_74
		raw_rest_86 = raw_rest_75
		has_name_QMARK__87 = has_name_QMARK__76
		and__x_88 = and__x_78
		goto b7
	} else {
		rest_forms_89 = rest_forms_68
		name_sym_90 = name_sym_69
		form_91 = form_70
		vec__1698_92 = vec__1698_71
		ctx_93 = ctx_72
		__94 = __73
		maybe_name_95 = maybe_name_74
		raw_rest_96 = raw_rest_75
		has_name_QMARK__97 = has_name_QMARK__76
		and__x_98 = and__x_78
		goto b8
	}
b7:
	;
	arg__1734_101, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{rest_forms_79})
	if callErr != nil {
		return nil, callErr
	}
	arg__1739_104, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{rest_forms_79})
	if callErr != nil {
		return nil, callErr
	}
	v105, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "list?").Deref(), []vm.Value{arg__1739_104})
	if callErr != nil {
		return nil, callErr
	}
	multi_QMARK__108 = v105
	rest_forms_109 = rest_forms_79
	name_sym_110 = name_sym_80
	form_111 = form_81
	vec__1698_112 = vec__1698_82
	ctx_113 = ctx_83
	__114 = __84
	maybe_name_115 = maybe_name_85
	raw_rest_116 = raw_rest_86
	has_name_QMARK__117 = has_name_QMARK__87
	and__x_118 = and__x_88
	goto b9
b8:
	;
	multi_QMARK__108 = and__x_98
	rest_forms_109 = rest_forms_89
	name_sym_110 = name_sym_90
	form_111 = form_91
	vec__1698_112 = vec__1698_92
	ctx_113 = ctx_93
	__114 = __94
	maybe_name_115 = maybe_name_95
	raw_rest_116 = raw_rest_96
	has_name_QMARK__117 = has_name_QMARK__97
	and__x_118 = and__x_98
	goto b9
b9:
	;
	if vm.IsTruthy(multi_QMARK__108) {
		multi_QMARK__119 = multi_QMARK__108
		rest_forms_120 = rest_forms_109
		name_sym_121 = name_sym_110
		form_122 = form_111
		vec__1698_123 = vec__1698_112
		ctx_124 = ctx_113
		__125 = __114
		maybe_name_126 = maybe_name_115
		raw_rest_127 = raw_rest_116
		has_name_QMARK__128 = has_name_QMARK__117
		goto b10
	} else {
		multi_QMARK__129 = multi_QMARK__108
		rest_forms_130 = rest_forms_109
		name_sym_131 = name_sym_110
		form_132 = form_111
		vec__1698_133 = vec__1698_112
		ctx_134 = ctx_113
		__135 = __114
		maybe_name_136 = maybe_name_115
		raw_rest_137 = raw_rest_116
		has_name_QMARK__138 = has_name_QMARK__117
		goto b11
	}
b10:
	;
	expanded_forms_143, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapv").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var args_vec_2 vm.Value
		var body_forms_4 vm.Value
		var e_6 vm.Value
		var arity_form_7 vm.Value
		var args_vec_8 vm.Value
		var body_forms_9 vm.Value
		var e_10 vm.Value
		var arity_form_11 vm.Value
		var args_vec_12 vm.Value
		var body_forms_13 vm.Value
		var e_14 vm.Value
		var v22 vm.Value
		var v24 vm.Value
		var arity_form_25 vm.Value
		var args_vec_26 vm.Value
		var body_forms_27 vm.Value
		var e_28 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = args_vec_2, body_forms_4, e_6, arity_form_7, args_vec_8, body_forms_9, e_10, arity_form_11, args_vec_12, body_forms_13, e_14, v22, v24, arity_form_25, args_vec_26, body_forms_27, e_28
		args_vec_2, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		body_forms_4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		e_6, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "expand-fn-args").Deref(), []vm.Value{args_vec_2, body_forms_4})
		if callErr != nil {
			return nil, callErr
		}
		if vm.IsTruthy(e_6) {
			arity_form_7 = arg0
			args_vec_8 = args_vec_2
			body_forms_9 = body_forms_4
			e_10 = e_6
			goto b1
		} else {
			arity_form_11 = arg0
			args_vec_12 = args_vec_2
			body_forms_13 = body_forms_4
			e_14 = e_6
			goto b2
		}
	b1:
		;
		v24 = e_10
		arity_form_25 = arity_form_7
		args_vec_26 = args_vec_8
		body_forms_27 = body_forms_9
		e_28 = e_10
		goto b3
	b2:
		;
		v22, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("variadic?"), vm.FALSE, vm.Keyword("flat-args"), args_vec_12, vm.Keyword("body"), body_forms_13})
		if callErr != nil {
			return nil, callErr
		}
		v24 = v22
		arity_form_25 = arity_form_11
		args_vec_26 = args_vec_12
		body_forms_27 = body_forms_13
		e_28 = e_14
		goto b3
	b3:
		;
		return v24, nil
	}), rest_forms_120})
	if callErr != nil {
		return nil, callErr
	}
	arg__1850_148, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__1922_154, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	all_caps_155, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "reduce").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
		var flat_args_6 vm.Value
		var body_10 vm.Value
		var arg__1864_12 vm.Value
		var arg__1868_15 vm.Value
		var arg_set_16 vm.Value
		var arg__1875_20 vm.Value
		var arg__1883_25 vm.Value
		var frees_26 vm.Value
		var arg__1901_34 vm.Value
		var arg__1919_43 vm.Value
		var v44 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _, _, _, _ = flat_args_6, body_10, arg__1864_12, arg__1868_15, arg_set_16, arg__1875_20, arg__1883_25, frees_26, arg__1901_34, arg__1919_43, v44
		flat_args_6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg1, vm.Keyword("flat-args")})
		if callErr != nil {
			return nil, callErr
		}
		body_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg1, vm.Keyword("body")})
		if callErr != nil {
			return nil, callErr
		}
		arg__1864_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
		if callErr != nil {
			return nil, callErr
		}
		arg__1868_15, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
		if callErr != nil {
			return nil, callErr
		}
		arg_set_16, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__1868_15, flat_args_6})
		if callErr != nil {
			return nil, callErr
		}
		arg__1875_20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "cons").Deref(), []vm.Value{vm.Symbol("do"), body_10})
		if callErr != nil {
			return nil, callErr
		}
		arg__1883_25, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "cons").Deref(), []vm.Value{vm.Symbol("do"), body_10})
		if callErr != nil {
			return nil, callErr
		}
		frees_26, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "free-vars").Deref(), []vm.Value{arg__1883_25, arg_set_16})
		if callErr != nil {
			return nil, callErr
		}
		arg__1901_34, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "filter").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
			var v3 vm.Value
			var callErr error
			_ = v3
			v3, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "lookup-local").Deref(), []vm.Value{ctx_124, arg0})
			if callErr != nil {
				return nil, callErr
			}
			return v3, nil
		}), frees_26})
		if callErr != nil {
			return nil, callErr
		}
		arg__1919_43, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "filter").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
			var v3 vm.Value
			var callErr error
			_ = v3
			v3, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "lookup-local").Deref(), []vm.Value{ctx_124, arg0})
			if callErr != nil {
				return nil, callErr
			}
			return v3, nil
		}), frees_26})
		if callErr != nil {
			return nil, callErr
		}
		v44, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg0, arg__1919_43})
		if callErr != nil {
			return nil, callErr
		}
		return v44, nil
	}), arg__1922_154, expanded_forms_143})
	if callErr != nil {
		return nil, callErr
	}
	arg__1929_159, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "sort-by").Deref(), []vm.Value{rt.LookupVar("clojure.core", "str").Deref(), all_caps_155})
	if callErr != nil {
		return nil, callErr
	}
	arg__1936_164, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "sort-by").Deref(), []vm.Value{rt.LookupVar("clojure.core", "str").Deref(), all_caps_155})
	if callErr != nil {
		return nil, callErr
	}
	captures_165, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg__1936_164})
	if callErr != nil {
		return nil, callErr
	}
	templates_175, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapv").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var flat_args_6 vm.Value
		var body_10 vm.Value
		var variadic_QMARK__14 vm.Value
		var v16 vm.Value
		var callErr error
		_, _, _, _ = flat_args_6, body_10, variadic_QMARK__14, v16
		flat_args_6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg0, vm.Keyword("flat-args")})
		if callErr != nil {
			return nil, callErr
		}
		body_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg0, vm.Keyword("body")})
		if callErr != nil {
			return nil, callErr
		}
		variadic_QMARK__14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg0, vm.Keyword("variadic?")})
		if callErr != nil {
			return nil, callErr
		}
		v16, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-inner-fn-template").Deref(), []vm.Value{name_sym_121, flat_args_6, body_10, captures_165, variadic_QMARK__14})
		if callErr != nil {
			return nil, callErr
		}
		return v16, nil
	}), expanded_forms_143})
	if callErr != nil {
		return nil, callErr
	}
	template_180, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("fns"), templates_175, vm.Keyword("kind"), vm.Keyword("multi-fn-template")})
	if callErr != nil {
		return nil, callErr
	}
	v182, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "emit-template-closure").Deref(), []vm.Value{template_180, captures_165, ctx_124})
	if callErr != nil {
		return nil, callErr
	}
	v191 = v182
	multi_QMARK__192 = multi_QMARK__119
	rest_forms_193 = rest_forms_120
	name_sym_194 = name_sym_121
	form_195 = form_122
	vec__1698_196 = vec__1698_123
	ctx_197 = ctx_124
	__198 = __125
	maybe_name_199 = maybe_name_126
	raw_rest_200 = raw_rest_127
	has_name_QMARK__201 = has_name_QMARK__128
	goto b12
b11:
	;
	args_vec_185, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{rest_forms_130})
	if callErr != nil {
		return nil, callErr
	}
	body_forms_187, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{rest_forms_130})
	if callErr != nil {
		return nil, callErr
	}
	v189, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-single-fn*").Deref(), []vm.Value{name_sym_131, args_vec_185, body_forms_187, ctx_134})
	if callErr != nil {
		return nil, callErr
	}
	v191 = v189
	multi_QMARK__192 = multi_QMARK__129
	rest_forms_193 = rest_forms_130
	name_sym_194 = name_sym_131
	form_195 = form_132
	vec__1698_196 = vec__1698_133
	ctx_197 = ctx_134
	__198 = __135
	maybe_name_199 = maybe_name_136
	raw_rest_200 = raw_rest_137
	has_name_QMARK__201 = has_name_QMARK__138
	goto b12
b12:
	;
	return v191, nil
}
func build_form(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var si_4 vm.Value
	var arg__2027_6 vm.Value
	var arg__2033_10 vm.Value
	var old_si_12 vm.Value
	var form_13 vm.Value
	var ctx_14 vm.Value
	var si_15 vm.Value
	var old_si_16 vm.Value
	var v27 vm.Value
	var form_17 vm.Value
	var ctx_18 vm.Value
	var si_19 vm.Value
	var old_si_20 vm.Value
	var __31 vm.Value
	var form_32 vm.Value
	var ctx_33 vm.Value
	var si_34 vm.Value
	var old_si_35 vm.Value
	var v47 vm.Value
	var __36 vm.Value
	var form_37 vm.Value
	var ctx_38 vm.Value
	var si_39 vm.Value
	var old_si_40 vm.Value
	var arg__2051_50 vm.Value
	var arg__2060_56 vm.Value
	var v60 vm.Value
	var __41 vm.Value
	var form_42 vm.Value
	var ctx_43 vm.Value
	var si_44 vm.Value
	var old_si_45 vm.Value
	var v73 vm.Value
	var nid_341 vm.Value
	var __342 vm.Value
	var form_343 vm.Value
	var ctx_344 vm.Value
	var si_345 vm.Value
	var old_si_346 vm.Value
	var __62 vm.Value
	var form_63 vm.Value
	var ctx_64 vm.Value
	var si_65 vm.Value
	var old_si_66 vm.Value
	var arg__2071_76 vm.Value
	var arg__2080_81 vm.Value
	var v84 vm.Value
	var __67 vm.Value
	var form_68 vm.Value
	var ctx_69 vm.Value
	var si_70 vm.Value
	var old_si_71 vm.Value
	var v97 vm.Value
	var v334 vm.Value
	var __335 vm.Value
	var form_336 vm.Value
	var ctx_337 vm.Value
	var si_338 vm.Value
	var old_si_339 vm.Value
	var __86 vm.Value
	var form_87 vm.Value
	var ctx_88 vm.Value
	var si_89 vm.Value
	var old_si_90 vm.Value
	var arg__2091_100 vm.Value
	var arg__2100_105 vm.Value
	var v108 vm.Value
	var __91 vm.Value
	var form_92 vm.Value
	var ctx_93 vm.Value
	var si_94 vm.Value
	var old_si_95 vm.Value
	var v121 vm.Value
	var v327 vm.Value
	var __328 vm.Value
	var form_329 vm.Value
	var ctx_330 vm.Value
	var si_331 vm.Value
	var old_si_332 vm.Value
	var __110 vm.Value
	var form_111 vm.Value
	var ctx_112 vm.Value
	var si_113 vm.Value
	var old_si_114 vm.Value
	var arg__2111_124 vm.Value
	var arg__2120_129 vm.Value
	var v132 vm.Value
	var __115 vm.Value
	var form_116 vm.Value
	var ctx_117 vm.Value
	var si_118 vm.Value
	var old_si_119 vm.Value
	var v145 vm.Value
	var v320 vm.Value
	var __321 vm.Value
	var form_322 vm.Value
	var ctx_323 vm.Value
	var si_324 vm.Value
	var old_si_325 vm.Value
	var __134 vm.Value
	var form_135 vm.Value
	var ctx_136 vm.Value
	var si_137 vm.Value
	var old_si_138 vm.Value
	var arg__2131_148 vm.Value
	var arg__2140_153 vm.Value
	var v156 vm.Value
	var __139 vm.Value
	var form_140 vm.Value
	var ctx_141 vm.Value
	var si_142 vm.Value
	var old_si_143 vm.Value
	var v169 vm.Value
	var v313 vm.Value
	var __314 vm.Value
	var form_315 vm.Value
	var ctx_316 vm.Value
	var si_317 vm.Value
	var old_si_318 vm.Value
	var __158 vm.Value
	var form_159 vm.Value
	var ctx_160 vm.Value
	var si_161 vm.Value
	var old_si_162 vm.Value
	var arg__2151_172 vm.Value
	var arg__2160_177 vm.Value
	var v180 vm.Value
	var __163 vm.Value
	var form_164 vm.Value
	var ctx_165 vm.Value
	var si_166 vm.Value
	var old_si_167 vm.Value
	var v193 vm.Value
	var v306 vm.Value
	var __307 vm.Value
	var form_308 vm.Value
	var ctx_309 vm.Value
	var si_310 vm.Value
	var old_si_311 vm.Value
	var __182 vm.Value
	var form_183 vm.Value
	var ctx_184 vm.Value
	var si_185 vm.Value
	var old_si_186 vm.Value
	var v196 vm.Value
	var __187 vm.Value
	var form_188 vm.Value
	var ctx_189 vm.Value
	var si_190 vm.Value
	var old_si_191 vm.Value
	var v209 vm.Value
	var v299 vm.Value
	var __300 vm.Value
	var form_301 vm.Value
	var ctx_302 vm.Value
	var si_303 vm.Value
	var old_si_304 vm.Value
	var __198 vm.Value
	var form_199 vm.Value
	var ctx_200 vm.Value
	var si_201 vm.Value
	var old_si_202 vm.Value
	var v212 vm.Value
	var __203 vm.Value
	var form_204 vm.Value
	var ctx_205 vm.Value
	var si_206 vm.Value
	var old_si_207 vm.Value
	var v225 vm.Value
	var v292 vm.Value
	var __293 vm.Value
	var form_294 vm.Value
	var ctx_295 vm.Value
	var si_296 vm.Value
	var old_si_297 vm.Value
	var __214 vm.Value
	var form_215 vm.Value
	var ctx_216 vm.Value
	var si_217 vm.Value
	var old_si_218 vm.Value
	var v228 vm.Value
	var __219 vm.Value
	var form_220 vm.Value
	var ctx_221 vm.Value
	var si_222 vm.Value
	var old_si_223 vm.Value
	var v241 vm.Value
	var v285 vm.Value
	var __286 vm.Value
	var form_287 vm.Value
	var ctx_288 vm.Value
	var si_289 vm.Value
	var old_si_290 vm.Value
	var __230 vm.Value
	var form_231 vm.Value
	var ctx_232 vm.Value
	var si_233 vm.Value
	var old_si_234 vm.Value
	var v244 vm.Value
	var __235 vm.Value
	var form_236 vm.Value
	var ctx_237 vm.Value
	var si_238 vm.Value
	var old_si_239 vm.Value
	var v278 vm.Value
	var __279 vm.Value
	var form_280 vm.Value
	var ctx_281 vm.Value
	var si_282 vm.Value
	var old_si_283 vm.Value
	var __246 vm.Value
	var form_247 vm.Value
	var ctx_248 vm.Value
	var si_249 vm.Value
	var old_si_250 vm.Value
	var arg__2201_261 vm.Value
	var arg__2208_266 vm.Value
	var v267 vm.Value
	var __251 vm.Value
	var form_252 vm.Value
	var ctx_253 vm.Value
	var si_254 vm.Value
	var old_si_255 vm.Value
	var v271 vm.Value
	var __272 vm.Value
	var form_273 vm.Value
	var ctx_274 vm.Value
	var si_275 vm.Value
	var old_si_276 vm.Value
	var nid_347 vm.Value
	var __348 vm.Value
	var form_349 vm.Value
	var ctx_350 vm.Value
	var si_351 vm.Value
	var old_si_352 vm.Value
	var v365 vm.Value
	var nid_353 vm.Value
	var __354 vm.Value
	var form_355 vm.Value
	var ctx_356 vm.Value
	var si_357 vm.Value
	var old_si_358 vm.Value
	var v369 vm.Value
	var nid_370 vm.Value
	var __371 vm.Value
	var form_372 vm.Value
	var ctx_373 vm.Value
	var si_374 vm.Value
	var old_si_375 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = si_4, arg__2027_6, arg__2033_10, old_si_12, form_13, ctx_14, si_15, old_si_16, v27, form_17, ctx_18, si_19, old_si_20, __31, form_32, ctx_33, si_34, old_si_35, v47, __36, form_37, ctx_38, si_39, old_si_40, arg__2051_50, arg__2060_56, v60, __41, form_42, ctx_43, si_44, old_si_45, v73, nid_341, __342, form_343, ctx_344, si_345, old_si_346, __62, form_63, ctx_64, si_65, old_si_66, arg__2071_76, arg__2080_81, v84, __67, form_68, ctx_69, si_70, old_si_71, v97, v334, __335, form_336, ctx_337, si_338, old_si_339, __86, form_87, ctx_88, si_89, old_si_90, arg__2091_100, arg__2100_105, v108, __91, form_92, ctx_93, si_94, old_si_95, v121, v327, __328, form_329, ctx_330, si_331, old_si_332, __110, form_111, ctx_112, si_113, old_si_114, arg__2111_124, arg__2120_129, v132, __115, form_116, ctx_117, si_118, old_si_119, v145, v320, __321, form_322, ctx_323, si_324, old_si_325, __134, form_135, ctx_136, si_137, old_si_138, arg__2131_148, arg__2140_153, v156, __139, form_140, ctx_141, si_142, old_si_143, v169, v313, __314, form_315, ctx_316, si_317, old_si_318, __158, form_159, ctx_160, si_161, old_si_162, arg__2151_172, arg__2160_177, v180, __163, form_164, ctx_165, si_166, old_si_167, v193, v306, __307, form_308, ctx_309, si_310, old_si_311, __182, form_183, ctx_184, si_185, old_si_186, v196, __187, form_188, ctx_189, si_190, old_si_191, v209, v299, __300, form_301, ctx_302, si_303, old_si_304, __198, form_199, ctx_200, si_201, old_si_202, v212, __203, form_204, ctx_205, si_206, old_si_207, v225, v292, __293, form_294, ctx_295, si_296, old_si_297, __214, form_215, ctx_216, si_217, old_si_218, v228, __219, form_220, ctx_221, si_222, old_si_223, v241, v285, __286, form_287, ctx_288, si_289, old_si_290, __230, form_231, ctx_232, si_233, old_si_234, v244, __235, form_236, ctx_237, si_238, old_si_239, v278, __279, form_280, ctx_281, si_282, old_si_283, __246, form_247, ctx_248, si_249, old_si_250, arg__2201_261, arg__2208_266, v267, __251, form_252, ctx_253, si_254, old_si_255, v271, __272, form_273, ctx_274, si_275, old_si_276, nid_347, __348, form_349, ctx_350, si_351, old_si_352, v365, nid_353, __354, form_355, ctx_356, si_357, old_si_358, v369, nid_370, __371, form_372, ctx_373, si_374, old_si_375
	si_4, callErr = rt.InvokeValue(rt.LookupVar("ir", "form-source-info").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__2027_6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__2033_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	old_si_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg__2033_10, vm.Keyword("source-info")})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(si_4) {
		form_13 = arg0
		ctx_14 = arg1
		si_15 = si_4
		old_si_16 = old_si_12
		goto b1
	} else {
		form_17 = arg0
		ctx_18 = arg1
		si_19 = si_4
		old_si_20 = old_si_12
		goto b2
	}
b1:
	;
	v27, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{ctx_14, rt.LookupVar("clojure.core", "assoc").Deref(), vm.Keyword("source-info"), si_15})
	if callErr != nil {
		return nil, callErr
	}
	__31 = v27
	form_32 = form_13
	ctx_33 = ctx_14
	si_34 = si_15
	old_si_35 = old_si_16
	goto b3
b2:
	;
	__31 = vm.NIL
	form_32 = form_17
	ctx_33 = ctx_18
	si_34 = si_19
	old_si_35 = old_si_20
	goto b3
b3:
	;
	v47, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{form_32})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v47) {
		__36 = __31
		form_37 = form_32
		ctx_38 = ctx_33
		si_39 = si_34
		old_si_40 = old_si_35
		goto b4
	} else {
		__41 = __31
		form_42 = form_32
		ctx_43 = ctx_33
		si_44 = si_34
		old_si_45 = old_si_35
		goto b5
	}
b4:
	;
	arg__2051_50, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_38})
	if callErr != nil {
		return nil, callErr
	}
	arg__2060_56, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_38})
	if callErr != nil {
		return nil, callErr
	}
	v60, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "add-inst!").Deref(), []vm.Value{ctx_38, arg__2060_56, vm.Keyword("const"), vm.NewArrayVector([]vm.Value{}), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	nid_341 = v60
	__342 = __36
	form_343 = form_37
	ctx_344 = ctx_38
	si_345 = si_39
	old_si_346 = old_si_40
	goto b6
b5:
	;
	v73, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "integer?").Deref(), []vm.Value{form_42})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v73) {
		__62 = __41
		form_63 = form_42
		ctx_64 = ctx_43
		si_65 = si_44
		old_si_66 = old_si_45
		goto b7
	} else {
		__67 = __41
		form_68 = form_42
		ctx_69 = ctx_43
		si_70 = si_44
		old_si_71 = old_si_45
		goto b8
	}
b6:
	;
	if vm.IsTruthy(si_345) {
		nid_347 = nid_341
		__348 = __342
		form_349 = form_343
		ctx_350 = ctx_344
		si_351 = si_345
		old_si_352 = old_si_346
		goto b37
	} else {
		nid_353 = nid_341
		__354 = __342
		form_355 = form_343
		ctx_356 = ctx_344
		si_357 = si_345
		old_si_358 = old_si_346
		goto b38
	}
b7:
	;
	arg__2071_76, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_64})
	if callErr != nil {
		return nil, callErr
	}
	arg__2080_81, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_64})
	if callErr != nil {
		return nil, callErr
	}
	v84, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "add-inst!").Deref(), []vm.Value{ctx_64, arg__2080_81, vm.Keyword("const"), vm.NewArrayVector([]vm.Value{}), form_63})
	if callErr != nil {
		return nil, callErr
	}
	v334 = v84
	__335 = __62
	form_336 = form_63
	ctx_337 = ctx_64
	si_338 = si_65
	old_si_339 = old_si_66
	goto b9
b8:
	;
	v97, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "float?").Deref(), []vm.Value{form_68})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v97) {
		__86 = __67
		form_87 = form_68
		ctx_88 = ctx_69
		si_89 = si_70
		old_si_90 = old_si_71
		goto b10
	} else {
		__91 = __67
		form_92 = form_68
		ctx_93 = ctx_69
		si_94 = si_70
		old_si_95 = old_si_71
		goto b11
	}
b9:
	;
	nid_341 = v334
	__342 = __335
	form_343 = form_336
	ctx_344 = ctx_337
	si_345 = si_338
	old_si_346 = old_si_339
	goto b6
b10:
	;
	arg__2091_100, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_88})
	if callErr != nil {
		return nil, callErr
	}
	arg__2100_105, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_88})
	if callErr != nil {
		return nil, callErr
	}
	v108, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "add-inst!").Deref(), []vm.Value{ctx_88, arg__2100_105, vm.Keyword("const"), vm.NewArrayVector([]vm.Value{}), form_87})
	if callErr != nil {
		return nil, callErr
	}
	v327 = v108
	__328 = __86
	form_329 = form_87
	ctx_330 = ctx_88
	si_331 = si_89
	old_si_332 = old_si_90
	goto b12
b11:
	;
	v121, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "string?").Deref(), []vm.Value{form_92})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v121) {
		__110 = __91
		form_111 = form_92
		ctx_112 = ctx_93
		si_113 = si_94
		old_si_114 = old_si_95
		goto b13
	} else {
		__115 = __91
		form_116 = form_92
		ctx_117 = ctx_93
		si_118 = si_94
		old_si_119 = old_si_95
		goto b14
	}
b12:
	;
	v334 = v327
	__335 = __328
	form_336 = form_329
	ctx_337 = ctx_330
	si_338 = si_331
	old_si_339 = old_si_332
	goto b9
b13:
	;
	arg__2111_124, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_112})
	if callErr != nil {
		return nil, callErr
	}
	arg__2120_129, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_112})
	if callErr != nil {
		return nil, callErr
	}
	v132, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "add-inst!").Deref(), []vm.Value{ctx_112, arg__2120_129, vm.Keyword("const"), vm.NewArrayVector([]vm.Value{}), form_111})
	if callErr != nil {
		return nil, callErr
	}
	v320 = v132
	__321 = __110
	form_322 = form_111
	ctx_323 = ctx_112
	si_324 = si_113
	old_si_325 = old_si_114
	goto b15
b14:
	;
	v145, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "keyword?").Deref(), []vm.Value{form_116})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v145) {
		__134 = __115
		form_135 = form_116
		ctx_136 = ctx_117
		si_137 = si_118
		old_si_138 = old_si_119
		goto b16
	} else {
		__139 = __115
		form_140 = form_116
		ctx_141 = ctx_117
		si_142 = si_118
		old_si_143 = old_si_119
		goto b17
	}
b15:
	;
	v327 = v320
	__328 = __321
	form_329 = form_322
	ctx_330 = ctx_323
	si_331 = si_324
	old_si_332 = old_si_325
	goto b12
b16:
	;
	arg__2131_148, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_136})
	if callErr != nil {
		return nil, callErr
	}
	arg__2140_153, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_136})
	if callErr != nil {
		return nil, callErr
	}
	v156, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "add-inst!").Deref(), []vm.Value{ctx_136, arg__2140_153, vm.Keyword("const"), vm.NewArrayVector([]vm.Value{}), form_135})
	if callErr != nil {
		return nil, callErr
	}
	v313 = v156
	__314 = __134
	form_315 = form_135
	ctx_316 = ctx_136
	si_317 = si_137
	old_si_318 = old_si_138
	goto b18
b17:
	;
	v169, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "boolean?").Deref(), []vm.Value{form_140})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v169) {
		__158 = __139
		form_159 = form_140
		ctx_160 = ctx_141
		si_161 = si_142
		old_si_162 = old_si_143
		goto b19
	} else {
		__163 = __139
		form_164 = form_140
		ctx_165 = ctx_141
		si_166 = si_142
		old_si_167 = old_si_143
		goto b20
	}
b18:
	;
	v320 = v313
	__321 = __314
	form_322 = form_315
	ctx_323 = ctx_316
	si_324 = si_317
	old_si_325 = old_si_318
	goto b15
b19:
	;
	arg__2151_172, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_160})
	if callErr != nil {
		return nil, callErr
	}
	arg__2160_177, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_160})
	if callErr != nil {
		return nil, callErr
	}
	v180, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "add-inst!").Deref(), []vm.Value{ctx_160, arg__2160_177, vm.Keyword("const"), vm.NewArrayVector([]vm.Value{}), form_159})
	if callErr != nil {
		return nil, callErr
	}
	v306 = v180
	__307 = __158
	form_308 = form_159
	ctx_309 = ctx_160
	si_310 = si_161
	old_si_311 = old_si_162
	goto b21
b20:
	;
	v193, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "symbol?").Deref(), []vm.Value{form_164})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v193) {
		__182 = __163
		form_183 = form_164
		ctx_184 = ctx_165
		si_185 = si_166
		old_si_186 = old_si_167
		goto b22
	} else {
		__187 = __163
		form_188 = form_164
		ctx_189 = ctx_165
		si_190 = si_166
		old_si_191 = old_si_167
		goto b23
	}
b21:
	;
	v313 = v306
	__314 = __307
	form_315 = form_308
	ctx_316 = ctx_309
	si_317 = si_310
	old_si_318 = old_si_311
	goto b18
b22:
	;
	v196, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-symbol").Deref(), []vm.Value{form_183, ctx_184})
	if callErr != nil {
		return nil, callErr
	}
	v299 = v196
	__300 = __182
	form_301 = form_183
	ctx_302 = ctx_184
	si_303 = si_185
	old_si_304 = old_si_186
	goto b24
b23:
	;
	v209, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector?").Deref(), []vm.Value{form_188})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v209) {
		__198 = __187
		form_199 = form_188
		ctx_200 = ctx_189
		si_201 = si_190
		old_si_202 = old_si_191
		goto b25
	} else {
		__203 = __187
		form_204 = form_188
		ctx_205 = ctx_189
		si_206 = si_190
		old_si_207 = old_si_191
		goto b26
	}
b24:
	;
	v306 = v299
	__307 = __300
	form_308 = form_301
	ctx_309 = ctx_302
	si_310 = si_303
	old_si_311 = old_si_304
	goto b21
b25:
	;
	v212, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-vector").Deref(), []vm.Value{form_199, ctx_200})
	if callErr != nil {
		return nil, callErr
	}
	v292 = v212
	__293 = __198
	form_294 = form_199
	ctx_295 = ctx_200
	si_296 = si_201
	old_si_297 = old_si_202
	goto b27
b26:
	;
	v225, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map?").Deref(), []vm.Value{form_204})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v225) {
		__214 = __203
		form_215 = form_204
		ctx_216 = ctx_205
		si_217 = si_206
		old_si_218 = old_si_207
		goto b28
	} else {
		__219 = __203
		form_220 = form_204
		ctx_221 = ctx_205
		si_222 = si_206
		old_si_223 = old_si_207
		goto b29
	}
b27:
	;
	v299 = v292
	__300 = __293
	form_301 = form_294
	ctx_302 = ctx_295
	si_303 = si_296
	old_si_304 = old_si_297
	goto b24
b28:
	;
	v228, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-map").Deref(), []vm.Value{form_215, ctx_216})
	if callErr != nil {
		return nil, callErr
	}
	v285 = v228
	__286 = __214
	form_287 = form_215
	ctx_288 = ctx_216
	si_289 = si_217
	old_si_290 = old_si_218
	goto b30
b29:
	;
	v241, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "list?").Deref(), []vm.Value{form_220})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v241) {
		__230 = __219
		form_231 = form_220
		ctx_232 = ctx_221
		si_233 = si_222
		old_si_234 = old_si_223
		goto b31
	} else {
		__235 = __219
		form_236 = form_220
		ctx_237 = ctx_221
		si_238 = si_222
		old_si_239 = old_si_223
		goto b32
	}
b30:
	;
	v292 = v285
	__293 = __286
	form_294 = form_287
	ctx_295 = ctx_288
	si_296 = si_289
	old_si_297 = old_si_290
	goto b27
b31:
	;
	v244, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-list").Deref(), []vm.Value{form_231, ctx_232})
	if callErr != nil {
		return nil, callErr
	}
	v278 = v244
	__279 = __230
	form_280 = form_231
	ctx_281 = ctx_232
	si_282 = si_233
	old_si_283 = old_si_234
	goto b33
b32:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		__246 = __235
		form_247 = form_236
		ctx_248 = ctx_237
		si_249 = si_238
		old_si_250 = old_si_239
		goto b34
	} else {
		__251 = __235
		form_252 = form_236
		ctx_253 = ctx_237
		si_254 = si_238
		old_si_255 = old_si_239
		goto b35
	}
b33:
	;
	v285 = v278
	__286 = __279
	form_287 = form_280
	ctx_288 = ctx_281
	si_289 = si_282
	old_si_290 = old_si_283
	goto b30
b34:
	;
	arg__2201_261, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("build-form: unrecognized form "), form_247})
	if callErr != nil {
		return nil, callErr
	}
	arg__2208_266, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("build-form: unrecognized form "), form_247})
	if callErr != nil {
		return nil, callErr
	}
	v267, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "throw").Deref(), []vm.Value{arg__2208_266})
	if callErr != nil {
		return nil, callErr
	}
	v271 = v267
	__272 = __246
	form_273 = form_247
	ctx_274 = ctx_248
	si_275 = si_249
	old_si_276 = old_si_250
	goto b36
b35:
	;
	v271 = vm.NIL
	__272 = __251
	form_273 = form_252
	ctx_274 = ctx_253
	si_275 = si_254
	old_si_276 = old_si_255
	goto b36
b36:
	;
	v278 = v271
	__279 = __272
	form_280 = form_273
	ctx_281 = ctx_274
	si_282 = si_275
	old_si_283 = old_si_276
	goto b33
b37:
	;
	v365, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{ctx_350, rt.LookupVar("clojure.core", "assoc").Deref(), vm.Keyword("source-info"), old_si_352})
	if callErr != nil {
		return nil, callErr
	}
	v369 = v365
	nid_370 = nid_347
	__371 = __348
	form_372 = form_349
	ctx_373 = ctx_350
	si_374 = si_351
	old_si_375 = old_si_352
	goto b39
b38:
	;
	v369 = vm.NIL
	nid_370 = nid_353
	__371 = __354
	form_372 = form_355
	ctx_373 = ctx_356
	si_374 = si_357
	old_si_375 = old_si_358
	goto b39
b39:
	;
	return nid_370, nil
}
func build_inner_fn_template(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value, arg3 vm.Value, arg4 vm.Value) (vm.Value, error) {
	var arg__3028_6 vm.Value
	var arg__3032_8 vm.Value
	var arg__3038_11 vm.Value
	var arg__3042_13 vm.Value
	var inner_14 vm.Value
	var entry_16 vm.Value
	var inner_ctx_18 vm.Value
	var arg__3054_20 vm.Value
	var arg__3059_23 vm.Value
	var arg__3064_26 vm.Value
	var arg__3069_29 vm.Value
	var arg__3070_30 vm.Value
	var arg__3076_33 vm.Value
	var arg__3081_36 vm.Value
	var arg__3086_39 vm.Value
	var arg__3091_42 vm.Value
	var arg__3092_43 vm.Value
	var v44 vm.Value
	var i_45 int
	var args_vec_46 vm.Value
	var inner_ctx_47 vm.Value
	var entry_48 vm.Value
	var v409 vm.Value
	var v412 vm.Value
	var arg__3097_70 vm.Value
	var v71 bool
	var name_sym_51 vm.Value
	var body_forms_52 vm.Value
	var capture_syms_53 vm.Value
	var variadic_QMARK__54 vm.Value
	var inner_55 vm.Value
	var i_56 int
	var args_vec_57 vm.Value
	var inner_ctx_58 vm.Value
	var entry_59 vm.Value
	var v408 vm.Value
	var v411 vm.Value
	var arg__3104_74 vm.Value
	var arg__3116_80 vm.Value
	var arg__3124_83 vm.Value
	var arg__3136_89 vm.Value
	var v90 vm.Value
	var v91 int
	var name_sym_60 vm.Value
	var body_forms_61 vm.Value
	var capture_syms_62 vm.Value
	var variadic_QMARK__63 vm.Value
	var inner_64 vm.Value
	var i_65 int
	var args_vec_66 vm.Value
	var inner_ctx_67 vm.Value
	var entry_68 vm.Value
	var v410 vm.Value
	var v413 vm.Value
	var v95 vm.Value
	var name_sym_96 vm.Value
	var body_forms_97 vm.Value
	var capture_syms_98 vm.Value
	var variadic_QMARK__99 vm.Value
	var inner_100 vm.Value
	var i_101 int
	var args_vec_102 vm.Value
	var inner_ctx_103 vm.Value
	var entry_104 vm.Value
	var i_105 int
	var capture_syms_106 vm.Value
	var i_107 int
	var inner_108 vm.Value
	var inner_ctx_109 vm.Value
	var entry_110 vm.Value
	var v423 vm.Value
	var v426 vm.Value
	var arg__3142_132 vm.Value
	var v133 bool
	var name_sym_113 vm.Value
	var body_forms_114 vm.Value
	var variadic_QMARK__115 vm.Value
	var args_vec_116 vm.Value
	var capture_syms_117 vm.Value
	var i_118 int
	var inner_119 vm.Value
	var inner_ctx_120 vm.Value
	var entry_121 vm.Value
	var v422 vm.Value
	var v425 vm.Value
	var arg__3149_136 vm.Value
	var arg__3161_142 vm.Value
	var arg__3169_145 vm.Value
	var arg__3181_151 vm.Value
	var v152 vm.Value
	var v153 int
	var name_sym_122 vm.Value
	var body_forms_123 vm.Value
	var variadic_QMARK__124 vm.Value
	var args_vec_125 vm.Value
	var capture_syms_126 vm.Value
	var i_127 int
	var inner_128 vm.Value
	var inner_ctx_129 vm.Value
	var entry_130 vm.Value
	var v424 vm.Value
	var v427 vm.Value
	var v157 vm.Value
	var name_sym_158 vm.Value
	var body_forms_159 vm.Value
	var variadic_QMARK__160 vm.Value
	var args_vec_161 vm.Value
	var capture_syms_162 vm.Value
	var i_163 int
	var inner_164 vm.Value
	var inner_ctx_165 vm.Value
	var entry_166 vm.Value
	var fs_167 vm.Value
	var last_id_168 vm.Value
	var inner_ctx_169 vm.Value
	var v195 vm.Value
	var name_sym_172 vm.Value
	var body_forms_173 vm.Value
	var variadic_QMARK__174 vm.Value
	var args_vec_175 vm.Value
	var capture_syms_176 vm.Value
	var i_177 int
	var inner_178 vm.Value
	var entry_179 vm.Value
	var fs_180 vm.Value
	var last_id_181 vm.Value
	var inner_ctx_182 vm.Value
	var v198 vm.Value
	var arg__3192_200 vm.Value
	var arg__3198_203 vm.Value
	var v204 vm.Value
	var name_sym_183 vm.Value
	var body_forms_184 vm.Value
	var variadic_QMARK__185 vm.Value
	var args_vec_186 vm.Value
	var capture_syms_187 vm.Value
	var i_188 int
	var inner_189 vm.Value
	var entry_190 vm.Value
	var fs_191 vm.Value
	var last_id_192 vm.Value
	var inner_ctx_193 vm.Value
	var last_val_207 vm.Value
	var name_sym_208 vm.Value
	var body_forms_209 vm.Value
	var variadic_QMARK__210 vm.Value
	var args_vec_211 vm.Value
	var capture_syms_212 vm.Value
	var i_213 int
	var inner_214 vm.Value
	var entry_215 vm.Value
	var fs_216 vm.Value
	var last_id_217 vm.Value
	var inner_ctx_218 vm.Value
	var final_blk_220 vm.Value
	var v248 vm.Value
	var last_val_221 vm.Value
	var name_sym_222 vm.Value
	var body_forms_223 vm.Value
	var variadic_QMARK__224 vm.Value
	var args_vec_225 vm.Value
	var capture_syms_226 vm.Value
	var i_227 int
	var inner_228 vm.Value
	var entry_229 vm.Value
	var fs_230 vm.Value
	var last_id_231 vm.Value
	var inner_ctx_232 vm.Value
	var final_blk_233 vm.Value
	var last_val_234 vm.Value
	var name_sym_235 vm.Value
	var body_forms_236 vm.Value
	var variadic_QMARK__237 vm.Value
	var args_vec_238 vm.Value
	var capture_syms_239 vm.Value
	var i_240 int
	var inner_241 vm.Value
	var entry_242 vm.Value
	var fs_243 vm.Value
	var last_id_244 vm.Value
	var inner_ctx_245 vm.Value
	var final_blk_246 vm.Value
	var v286 vm.Value
	var v376 vm.Value
	var last_val_377 vm.Value
	var name_sym_378 vm.Value
	var body_forms_379 vm.Value
	var variadic_QMARK__380 vm.Value
	var args_vec_381 vm.Value
	var capture_syms_382 vm.Value
	var i_383 int
	var inner_384 vm.Value
	var entry_385 vm.Value
	var fs_386 vm.Value
	var last_id_387 vm.Value
	var inner_ctx_388 vm.Value
	var final_blk_389 vm.Value
	var arg__3236_395 vm.Value
	var v398 vm.Value
	var last_val_253 vm.Value
	var name_sym_254 vm.Value
	var body_forms_255 vm.Value
	var variadic_QMARK__256 vm.Value
	var args_vec_257 vm.Value
	var capture_syms_258 vm.Value
	var i_259 int
	var inner_260 vm.Value
	var entry_261 vm.Value
	var fs_262 vm.Value
	var last_id_263 vm.Value
	var arg__3206_264 vm.Value
	var inner_ctx_265 vm.Value
	var final_blk_266 vm.Value
	var arg__3207_267 vm.Value
	var arg__3208_268 vm.Value
	var last_val_269 vm.Value
	var name_sym_270 vm.Value
	var body_forms_271 vm.Value
	var variadic_QMARK__272 vm.Value
	var args_vec_273 vm.Value
	var capture_syms_274 vm.Value
	var i_275 int
	var inner_276 vm.Value
	var entry_277 vm.Value
	var fs_278 vm.Value
	var last_id_279 vm.Value
	var arg__3206_280 vm.Value
	var inner_ctx_281 vm.Value
	var final_blk_282 vm.Value
	var arg__3207_283 vm.Value
	var arg__3208_284 vm.Value
	var v291 vm.Value
	var arg__3214_293 vm.Value
	var last_val_294 vm.Value
	var name_sym_295 vm.Value
	var body_forms_296 vm.Value
	var variadic_QMARK__297 vm.Value
	var args_vec_298 vm.Value
	var capture_syms_299 vm.Value
	var i_300 int
	var inner_301 vm.Value
	var entry_302 vm.Value
	var fs_303 vm.Value
	var last_id_304 vm.Value
	var arg__3206_305 vm.Value
	var inner_ctx_306 vm.Value
	var final_blk_307 vm.Value
	var arg__3207_308 vm.Value
	var arg__3208_309 vm.Value
	var v348 vm.Value
	var last_val_313 vm.Value
	var name_sym_314 vm.Value
	var body_forms_315 vm.Value
	var variadic_QMARK__316 vm.Value
	var args_vec_317 vm.Value
	var capture_syms_318 vm.Value
	var i_319 int
	var inner_320 vm.Value
	var entry_321 vm.Value
	var fs_322 vm.Value
	var last_id_323 vm.Value
	var arg__3217_324 vm.Value
	var inner_ctx_325 vm.Value
	var final_blk_326 vm.Value
	var arg__3218_327 vm.Value
	var head__3216_328 vm.Value
	var arg__3219_329 vm.Value
	var last_val_330 vm.Value
	var name_sym_331 vm.Value
	var body_forms_332 vm.Value
	var variadic_QMARK__333 vm.Value
	var args_vec_334 vm.Value
	var capture_syms_335 vm.Value
	var i_336 int
	var inner_337 vm.Value
	var entry_338 vm.Value
	var fs_339 vm.Value
	var last_id_340 vm.Value
	var arg__3217_341 vm.Value
	var inner_ctx_342 vm.Value
	var final_blk_343 vm.Value
	var arg__3218_344 vm.Value
	var head__3216_345 vm.Value
	var arg__3219_346 vm.Value
	var v353 vm.Value
	var arg__3225_355 vm.Value
	var last_val_356 vm.Value
	var name_sym_357 vm.Value
	var body_forms_358 vm.Value
	var variadic_QMARK__359 vm.Value
	var args_vec_360 vm.Value
	var capture_syms_361 vm.Value
	var i_362 int
	var inner_363 vm.Value
	var entry_364 vm.Value
	var fs_365 vm.Value
	var last_id_366 vm.Value
	var arg__3217_367 vm.Value
	var inner_ctx_368 vm.Value
	var final_blk_369 vm.Value
	var arg__3218_370 vm.Value
	var head__3216_371 vm.Value
	var arg__3219_372 vm.Value
	var v374 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__3028_6, arg__3032_8, arg__3038_11, arg__3042_13, inner_14, entry_16, inner_ctx_18, arg__3054_20, arg__3059_23, arg__3064_26, arg__3069_29, arg__3070_30, arg__3076_33, arg__3081_36, arg__3086_39, arg__3091_42, arg__3092_43, v44, i_45, args_vec_46, inner_ctx_47, entry_48, v409, v412, arg__3097_70, v71, name_sym_51, body_forms_52, capture_syms_53, variadic_QMARK__54, inner_55, i_56, args_vec_57, inner_ctx_58, entry_59, v408, v411, arg__3104_74, arg__3116_80, arg__3124_83, arg__3136_89, v90, v91, name_sym_60, body_forms_61, capture_syms_62, variadic_QMARK__63, inner_64, i_65, args_vec_66, inner_ctx_67, entry_68, v410, v413, v95, name_sym_96, body_forms_97, capture_syms_98, variadic_QMARK__99, inner_100, i_101, args_vec_102, inner_ctx_103, entry_104, i_105, capture_syms_106, i_107, inner_108, inner_ctx_109, entry_110, v423, v426, arg__3142_132, v133, name_sym_113, body_forms_114, variadic_QMARK__115, args_vec_116, capture_syms_117, i_118, inner_119, inner_ctx_120, entry_121, v422, v425, arg__3149_136, arg__3161_142, arg__3169_145, arg__3181_151, v152, v153, name_sym_122, body_forms_123, variadic_QMARK__124, args_vec_125, capture_syms_126, i_127, inner_128, inner_ctx_129, entry_130, v424, v427, v157, name_sym_158, body_forms_159, variadic_QMARK__160, args_vec_161, capture_syms_162, i_163, inner_164, inner_ctx_165, entry_166, fs_167, last_id_168, inner_ctx_169, v195, name_sym_172, body_forms_173, variadic_QMARK__174, args_vec_175, capture_syms_176, i_177, inner_178, entry_179, fs_180, last_id_181, inner_ctx_182, v198, arg__3192_200, arg__3198_203, v204, name_sym_183, body_forms_184, variadic_QMARK__185, args_vec_186, capture_syms_187, i_188, inner_189, entry_190, fs_191, last_id_192, inner_ctx_193, last_val_207, name_sym_208, body_forms_209, variadic_QMARK__210, args_vec_211, capture_syms_212, i_213, inner_214, entry_215, fs_216, last_id_217, inner_ctx_218, final_blk_220, v248, last_val_221, name_sym_222, body_forms_223, variadic_QMARK__224, args_vec_225, capture_syms_226, i_227, inner_228, entry_229, fs_230, last_id_231, inner_ctx_232, final_blk_233, last_val_234, name_sym_235, body_forms_236, variadic_QMARK__237, args_vec_238, capture_syms_239, i_240, inner_241, entry_242, fs_243, last_id_244, inner_ctx_245, final_blk_246, v286, v376, last_val_377, name_sym_378, body_forms_379, variadic_QMARK__380, args_vec_381, capture_syms_382, i_383, inner_384, entry_385, fs_386, last_id_387, inner_ctx_388, final_blk_389, arg__3236_395, v398, last_val_253, name_sym_254, body_forms_255, variadic_QMARK__256, args_vec_257, capture_syms_258, i_259, inner_260, entry_261, fs_262, last_id_263, arg__3206_264, inner_ctx_265, final_blk_266, arg__3207_267, arg__3208_268, last_val_269, name_sym_270, body_forms_271, variadic_QMARK__272, args_vec_273, capture_syms_274, i_275, inner_276, entry_277, fs_278, last_id_279, arg__3206_280, inner_ctx_281, final_blk_282, arg__3207_283, arg__3208_284, v291, arg__3214_293, last_val_294, name_sym_295, body_forms_296, variadic_QMARK__297, args_vec_298, capture_syms_299, i_300, inner_301, entry_302, fs_303, last_id_304, arg__3206_305, inner_ctx_306, final_blk_307, arg__3207_308, arg__3208_309, v348, last_val_313, name_sym_314, body_forms_315, variadic_QMARK__316, args_vec_317, capture_syms_318, i_319, inner_320, entry_321, fs_322, last_id_323, arg__3217_324, inner_ctx_325, final_blk_326, arg__3218_327, head__3216_328, arg__3219_329, last_val_330, name_sym_331, body_forms_332, variadic_QMARK__333, args_vec_334, capture_syms_335, i_336, inner_337, entry_338, fs_339, last_id_340, arg__3217_341, inner_ctx_342, final_blk_343, arg__3218_344, head__3216_345, arg__3219_346, v353, arg__3225_355, last_val_356, name_sym_357, body_forms_358, variadic_QMARK__359, args_vec_360, capture_syms_361, i_362, inner_363, entry_364, fs_365, last_id_366, arg__3217_367, inner_ctx_368, final_blk_369, arg__3218_370, head__3216_371, arg__3219_372, v374
	arg__3028_6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__3032_8, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__3038_11, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__3042_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	inner_14, callErr = rt.InvokeValue(rt.LookupVar("ir", "new-fn").Deref(), []vm.Value{arg__3038_11, arg__3042_13, arg4})
	if callErr != nil {
		return nil, callErr
	}
	entry_16, callErr = rt.InvokeValue(rt.LookupVar("ir", "entry-block").Deref(), []vm.Value{inner_14})
	if callErr != nil {
		return nil, callErr
	}
	inner_ctx_18, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "new-context").Deref(), []vm.Value{inner_14})
	if callErr != nil {
		return nil, callErr
	}
	arg__3054_20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{inner_ctx_18})
	if callErr != nil {
		return nil, callErr
	}
	arg__3059_23, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__3064_26, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{inner_ctx_18})
	if callErr != nil {
		return nil, callErr
	}
	arg__3069_29, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__3070_30, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "assoc").Deref(), []vm.Value{arg__3064_26, vm.Keyword("fn-arg-syms"), arg__3069_29})
	if callErr != nil {
		return nil, callErr
	}
	arg__3076_33, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{inner_ctx_18})
	if callErr != nil {
		return nil, callErr
	}
	arg__3081_36, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__3086_39, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{inner_ctx_18})
	if callErr != nil {
		return nil, callErr
	}
	arg__3091_42, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__3092_43, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "assoc").Deref(), []vm.Value{arg__3086_39, vm.Keyword("fn-arg-syms"), arg__3091_42})
	if callErr != nil {
		return nil, callErr
	}
	v44, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "reset!").Deref(), []vm.Value{inner_ctx_18, arg__3092_43})
	if callErr != nil {
		return nil, callErr
	}
	i_45 = 0
	args_vec_46 = arg1
	inner_ctx_47 = inner_ctx_18
	entry_48 = entry_16
	v409 = vm.Keyword("load-arg")
	v412 = vm.NewArrayVector([]vm.Value{})
	goto b1
b1:
	;
	arg__3097_70, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{args_vec_46})
	if callErr != nil {
		return nil, callErr
	}
	v71 = rt.LtValue(vm.Int(i_45), arg__3097_70)
	if v71 {
		name_sym_51 = arg0
		body_forms_52 = arg2
		capture_syms_53 = arg3
		variadic_QMARK__54 = arg4
		inner_55 = inner_14
		i_56 = i_45
		args_vec_57 = args_vec_46
		inner_ctx_58 = inner_ctx_47
		entry_59 = entry_48
		v408 = v409
		v411 = v412
		goto b2
	} else {
		name_sym_60 = arg0
		body_forms_61 = arg2
		capture_syms_62 = arg3
		variadic_QMARK__63 = arg4
		inner_64 = inner_14
		i_65 = i_45
		args_vec_66 = args_vec_46
		inner_ctx_67 = inner_ctx_47
		entry_68 = entry_48
		v410 = v409
		v413 = v412
		goto b3
	}
b2:
	;
	arg__3104_74, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{args_vec_57, vm.Int(i_56)})
	if callErr != nil {
		return nil, callErr
	}
	arg__3116_80, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "add-inst!").Deref(), []vm.Value{inner_ctx_58, entry_59, v408, v411, vm.Int(i_56)})
	if callErr != nil {
		return nil, callErr
	}
	arg__3124_83, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{args_vec_57, vm.Int(i_56)})
	if callErr != nil {
		return nil, callErr
	}
	arg__3136_89, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "add-inst!").Deref(), []vm.Value{inner_ctx_58, entry_59, v408, v411, vm.Int(i_56)})
	if callErr != nil {
		return nil, callErr
	}
	v90, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "bind-local!").Deref(), []vm.Value{inner_ctx_58, arg__3124_83, arg__3136_89})
	if callErr != nil {
		return nil, callErr
	}
	v91 = i_56 + 1
	i_45 = v91
	args_vec_46 = args_vec_57
	inner_ctx_47 = inner_ctx_58
	entry_48 = entry_59
	v409 = v408
	v412 = v411
	goto b1
b3:
	;
	v95 = vm.NIL
	name_sym_96 = name_sym_60
	body_forms_97 = body_forms_61
	capture_syms_98 = capture_syms_62
	variadic_QMARK__99 = variadic_QMARK__63
	inner_100 = inner_64
	i_101 = i_65
	args_vec_102 = args_vec_66
	inner_ctx_103 = inner_ctx_67
	entry_104 = entry_68
	goto b4
b4:
	;
	i_105 = 0
	capture_syms_106 = capture_syms_98
	i_107 = i_101
	inner_108 = inner_100
	inner_ctx_109 = inner_ctx_103
	entry_110 = entry_104
	v423 = vm.Keyword("load-closed")
	v426 = vm.NewArrayVector([]vm.Value{})
	goto b5
b5:
	;
	arg__3142_132, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{capture_syms_106})
	if callErr != nil {
		return nil, callErr
	}
	v133 = rt.LtValue(vm.Int(i_107), arg__3142_132)
	if v133 {
		name_sym_113 = name_sym_96
		body_forms_114 = body_forms_97
		variadic_QMARK__115 = variadic_QMARK__99
		args_vec_116 = args_vec_102
		capture_syms_117 = capture_syms_106
		i_118 = i_107
		inner_119 = inner_108
		inner_ctx_120 = inner_ctx_109
		entry_121 = entry_110
		v422 = v423
		v425 = v426
		goto b6
	} else {
		name_sym_122 = name_sym_96
		body_forms_123 = body_forms_97
		variadic_QMARK__124 = variadic_QMARK__99
		args_vec_125 = args_vec_102
		capture_syms_126 = capture_syms_106
		i_127 = i_107
		inner_128 = inner_108
		inner_ctx_129 = inner_ctx_109
		entry_130 = entry_110
		v424 = v423
		v427 = v426
		goto b7
	}
b6:
	;
	arg__3149_136, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{capture_syms_117, vm.Int(i_118)})
	if callErr != nil {
		return nil, callErr
	}
	arg__3161_142, callErr = rt.InvokeValue(rt.LookupVar("ir", "add-inst").Deref(), []vm.Value{inner_119, entry_121, v422, v425, vm.Int(i_118)})
	if callErr != nil {
		return nil, callErr
	}
	arg__3169_145, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{capture_syms_117, vm.Int(i_118)})
	if callErr != nil {
		return nil, callErr
	}
	arg__3181_151, callErr = rt.InvokeValue(rt.LookupVar("ir", "add-inst").Deref(), []vm.Value{inner_119, entry_121, v422, v425, vm.Int(i_118)})
	if callErr != nil {
		return nil, callErr
	}
	v152, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "bind-local!").Deref(), []vm.Value{inner_ctx_120, arg__3169_145, arg__3181_151})
	if callErr != nil {
		return nil, callErr
	}
	v153 = i_118 + 1
	i_105 = v153
	capture_syms_106 = capture_syms_117
	i_107 = i_118
	inner_108 = inner_119
	inner_ctx_109 = inner_ctx_120
	entry_110 = entry_121
	v423 = v422
	v426 = v425
	goto b5
b7:
	;
	v157 = vm.NIL
	name_sym_158 = name_sym_122
	body_forms_159 = body_forms_123
	variadic_QMARK__160 = variadic_QMARK__124
	args_vec_161 = args_vec_125
	capture_syms_162 = capture_syms_126
	i_163 = i_127
	inner_164 = inner_128
	inner_ctx_165 = inner_ctx_129
	entry_166 = entry_130
	goto b8
b8:
	;
	fs_167 = body_forms_159
	last_id_168 = vm.NIL
	inner_ctx_169 = inner_ctx_165
	goto b9
b9:
	;
	v195, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{fs_167})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v195) {
		name_sym_172 = name_sym_158
		body_forms_173 = body_forms_159
		variadic_QMARK__174 = variadic_QMARK__160
		args_vec_175 = args_vec_161
		capture_syms_176 = capture_syms_162
		i_177 = i_163
		inner_178 = inner_164
		entry_179 = entry_166
		fs_180 = fs_167
		last_id_181 = last_id_168
		inner_ctx_182 = inner_ctx_169
		goto b10
	} else {
		name_sym_183 = name_sym_158
		body_forms_184 = body_forms_159
		variadic_QMARK__185 = variadic_QMARK__160
		args_vec_186 = args_vec_161
		capture_syms_187 = capture_syms_162
		i_188 = i_163
		inner_189 = inner_164
		entry_190 = entry_166
		fs_191 = fs_167
		last_id_192 = last_id_168
		inner_ctx_193 = inner_ctx_169
		goto b11
	}
b10:
	;
	v198, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{fs_180})
	if callErr != nil {
		return nil, callErr
	}
	arg__3192_200, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{fs_180})
	if callErr != nil {
		return nil, callErr
	}
	arg__3198_203, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{fs_180})
	if callErr != nil {
		return nil, callErr
	}
	v204, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-form").Deref(), []vm.Value{arg__3198_203, inner_ctx_182})
	if callErr != nil {
		return nil, callErr
	}
	fs_167 = v198
	last_id_168 = v204
	inner_ctx_169 = inner_ctx_182
	goto b9
b11:
	;
	last_val_207 = last_id_192
	name_sym_208 = name_sym_183
	body_forms_209 = body_forms_184
	variadic_QMARK__210 = variadic_QMARK__185
	args_vec_211 = args_vec_186
	capture_syms_212 = capture_syms_187
	i_213 = i_188
	inner_214 = inner_189
	entry_215 = entry_190
	fs_216 = fs_191
	last_id_217 = last_id_192
	inner_ctx_218 = inner_ctx_193
	goto b12
b12:
	;
	final_blk_220, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{inner_ctx_218})
	if callErr != nil {
		return nil, callErr
	}
	v248, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "terminated?").Deref(), []vm.Value{last_val_207})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v248) {
		last_val_221 = last_val_207
		name_sym_222 = name_sym_208
		body_forms_223 = body_forms_209
		variadic_QMARK__224 = variadic_QMARK__210
		args_vec_225 = args_vec_211
		capture_syms_226 = capture_syms_212
		i_227 = i_213
		inner_228 = inner_214
		entry_229 = entry_215
		fs_230 = fs_216
		last_id_231 = last_id_217
		inner_ctx_232 = inner_ctx_218
		final_blk_233 = final_blk_220
		goto b13
	} else {
		last_val_234 = last_val_207
		name_sym_235 = name_sym_208
		body_forms_236 = body_forms_209
		variadic_QMARK__237 = variadic_QMARK__210
		args_vec_238 = args_vec_211
		capture_syms_239 = capture_syms_212
		i_240 = i_213
		inner_241 = inner_214
		entry_242 = entry_215
		fs_243 = fs_216
		last_id_244 = last_id_217
		inner_ctx_245 = inner_ctx_218
		final_blk_246 = final_blk_220
		goto b14
	}
b13:
	;
	v376 = vm.NIL
	last_val_377 = last_val_221
	name_sym_378 = name_sym_222
	body_forms_379 = body_forms_223
	variadic_QMARK__380 = variadic_QMARK__224
	args_vec_381 = args_vec_225
	capture_syms_382 = capture_syms_226
	i_383 = i_227
	inner_384 = inner_228
	entry_385 = entry_229
	fs_386 = fs_230
	last_id_387 = last_id_231
	inner_ctx_388 = inner_ctx_232
	final_blk_389 = final_blk_233
	goto b15
b14:
	;
	v286, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{last_val_234})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v286) {
		last_val_253 = last_val_234
		name_sym_254 = name_sym_235
		body_forms_255 = body_forms_236
		variadic_QMARK__256 = variadic_QMARK__237
		args_vec_257 = args_vec_238
		capture_syms_258 = capture_syms_239
		i_259 = i_240
		inner_260 = inner_241
		entry_261 = entry_242
		fs_262 = fs_243
		last_id_263 = last_id_244
		arg__3206_264 = inner_ctx_245
		inner_ctx_265 = inner_ctx_245
		final_blk_266 = final_blk_246
		arg__3207_267 = final_blk_246
		arg__3208_268 = vm.Keyword("return")
		goto b16
	} else {
		last_val_269 = last_val_234
		name_sym_270 = name_sym_235
		body_forms_271 = body_forms_236
		variadic_QMARK__272 = variadic_QMARK__237
		args_vec_273 = args_vec_238
		capture_syms_274 = capture_syms_239
		i_275 = i_240
		inner_276 = inner_241
		entry_277 = entry_242
		fs_278 = fs_243
		last_id_279 = last_id_244
		arg__3206_280 = inner_ctx_245
		inner_ctx_281 = inner_ctx_245
		final_blk_282 = final_blk_246
		arg__3207_283 = final_blk_246
		arg__3208_284 = vm.Keyword("return")
		goto b17
	}
b15:
	;
	arg__3236_395, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{args_vec_381})
	if callErr != nil {
		return nil, callErr
	}
	v398, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("fn"), inner_384, vm.Keyword("variadic?"), variadic_QMARK__380, vm.Keyword("arity"), arg__3236_395, vm.Keyword("kind"), vm.Keyword("fn-template")})
	if callErr != nil {
		return nil, callErr
	}
	return v398, nil
b16:
	;
	arg__3214_293 = vm.NewArrayVector([]vm.Value{})
	last_val_294 = last_val_253
	name_sym_295 = name_sym_254
	body_forms_296 = body_forms_255
	variadic_QMARK__297 = variadic_QMARK__256
	args_vec_298 = args_vec_257
	capture_syms_299 = capture_syms_258
	i_300 = i_259
	inner_301 = inner_260
	entry_302 = entry_261
	fs_303 = fs_262
	last_id_304 = last_id_263
	arg__3206_305 = arg__3206_264
	inner_ctx_306 = inner_ctx_265
	final_blk_307 = final_blk_266
	arg__3207_308 = arg__3207_267
	arg__3208_309 = arg__3208_268
	goto b18
b17:
	;
	v291, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{last_val_269})
	if callErr != nil {
		return nil, callErr
	}
	arg__3214_293 = v291
	last_val_294 = last_val_269
	name_sym_295 = name_sym_270
	body_forms_296 = body_forms_271
	variadic_QMARK__297 = variadic_QMARK__272
	args_vec_298 = args_vec_273
	capture_syms_299 = capture_syms_274
	i_300 = i_275
	inner_301 = inner_276
	entry_302 = entry_277
	fs_303 = fs_278
	last_id_304 = last_id_279
	arg__3206_305 = arg__3206_280
	inner_ctx_306 = inner_ctx_281
	final_blk_307 = final_blk_282
	arg__3207_308 = arg__3207_283
	arg__3208_309 = arg__3208_284
	goto b18
b18:
	;
	v348, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{last_val_294})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v348) {
		last_val_313 = last_val_294
		name_sym_314 = name_sym_295
		body_forms_315 = body_forms_296
		variadic_QMARK__316 = variadic_QMARK__297
		args_vec_317 = args_vec_298
		capture_syms_318 = capture_syms_299
		i_319 = i_300
		inner_320 = inner_301
		entry_321 = entry_302
		fs_322 = fs_303
		last_id_323 = last_id_304
		arg__3217_324 = inner_ctx_306
		inner_ctx_325 = inner_ctx_306
		final_blk_326 = final_blk_307
		arg__3218_327 = final_blk_307
		head__3216_328 = rt.LookupVar("ir.build", "add-terminator!").Deref()
		arg__3219_329 = vm.Keyword("return")
		goto b19
	} else {
		last_val_330 = last_val_294
		name_sym_331 = name_sym_295
		body_forms_332 = body_forms_296
		variadic_QMARK__333 = variadic_QMARK__297
		args_vec_334 = args_vec_298
		capture_syms_335 = capture_syms_299
		i_336 = i_300
		inner_337 = inner_301
		entry_338 = entry_302
		fs_339 = fs_303
		last_id_340 = last_id_304
		arg__3217_341 = inner_ctx_306
		inner_ctx_342 = inner_ctx_306
		final_blk_343 = final_blk_307
		arg__3218_344 = final_blk_307
		head__3216_345 = rt.LookupVar("ir.build", "add-terminator!").Deref()
		arg__3219_346 = vm.Keyword("return")
		goto b20
	}
b19:
	;
	arg__3225_355 = vm.NewArrayVector([]vm.Value{})
	last_val_356 = last_val_313
	name_sym_357 = name_sym_314
	body_forms_358 = body_forms_315
	variadic_QMARK__359 = variadic_QMARK__316
	args_vec_360 = args_vec_317
	capture_syms_361 = capture_syms_318
	i_362 = i_319
	inner_363 = inner_320
	entry_364 = entry_321
	fs_365 = fs_322
	last_id_366 = last_id_323
	arg__3217_367 = arg__3217_324
	inner_ctx_368 = inner_ctx_325
	final_blk_369 = final_blk_326
	arg__3218_370 = arg__3218_327
	head__3216_371 = head__3216_328
	arg__3219_372 = arg__3219_329
	goto b21
b20:
	;
	v353, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{last_val_330})
	if callErr != nil {
		return nil, callErr
	}
	arg__3225_355 = v353
	last_val_356 = last_val_330
	name_sym_357 = name_sym_331
	body_forms_358 = body_forms_332
	variadic_QMARK__359 = variadic_QMARK__333
	args_vec_360 = args_vec_334
	capture_syms_361 = capture_syms_335
	i_362 = i_336
	inner_363 = inner_337
	entry_364 = entry_338
	fs_365 = fs_339
	last_id_366 = last_id_340
	arg__3217_367 = arg__3217_341
	inner_ctx_368 = inner_ctx_342
	final_blk_369 = final_blk_343
	arg__3218_370 = arg__3218_344
	head__3216_371 = head__3216_345
	arg__3219_372 = arg__3219_346
	goto b21
b21:
	;
	v374, callErr = rt.InvokeValue(head__3216_371, []vm.Value{arg__3217_367, arg__3218_370, arg__3219_372, arg__3225_355, vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	v376 = v374
	last_val_377 = last_val_356
	name_sym_378 = name_sym_357
	body_forms_379 = body_forms_358
	variadic_QMARK__380 = variadic_QMARK__359
	args_vec_381 = args_vec_360
	capture_syms_382 = capture_syms_361
	i_383 = i_362
	inner_384 = inner_363
	entry_385 = entry_364
	fs_386 = fs_365
	last_id_387 = last_id_366
	inner_ctx_388 = inner_ctx_368
	final_blk_389 = final_blk_369
	goto b15
}
func build_let(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var __8 vm.Value
	var bindings_14 vm.Value
	var body_18 vm.Value
	var v20 vm.Value
	var i_21 int
	var bindings_22 vm.Value
	var ctx_23 vm.Value
	var v377 int
	var arg__3269_41 vm.Value
	var v42 bool
	var form_26 vm.Value
	var vec__3239_27 vm.Value
	var __28 vm.Value
	var body_29 vm.Value
	var i_30 int
	var bindings_31 vm.Value
	var ctx_32 vm.Value
	var v376 int
	var sym_45 vm.Value
	var arg__3277_46 int
	var expr_49 vm.Value
	var expr_id_51 vm.Value
	var v53 vm.Value
	var v55 int
	var form_33 vm.Value
	var vec__3239_34 vm.Value
	var __35 vm.Value
	var body_36 vm.Value
	var i_37 int
	var bindings_38 vm.Value
	var ctx_39 vm.Value
	var v378 int
	var v59 vm.Value
	var form_60 vm.Value
	var vec__3239_61 vm.Value
	var __62 vm.Value
	var body_63 vm.Value
	var i_64 int
	var bindings_65 vm.Value
	var ctx_66 vm.Value
	var fs_67 vm.Value
	var last_id_68 vm.Value
	var ctx_69 vm.Value
	var v91 vm.Value
	var form_72 vm.Value
	var vec__3239_73 vm.Value
	var __74 vm.Value
	var body_75 vm.Value
	var i_76 int
	var bindings_77 vm.Value
	var fs_78 vm.Value
	var last_id_79 vm.Value
	var ctx_80 vm.Value
	var v94 vm.Value
	var arg__3305_96 vm.Value
	var arg__3311_99 vm.Value
	var v100 vm.Value
	var form_81 vm.Value
	var vec__3239_82 vm.Value
	var __83 vm.Value
	var body_84 vm.Value
	var i_85 int
	var bindings_86 vm.Value
	var fs_87 vm.Value
	var last_id_88 vm.Value
	var ctx_89 vm.Value
	var result_103 vm.Value
	var form_104 vm.Value
	var vec__3239_105 vm.Value
	var __106 vm.Value
	var body_107 vm.Value
	var i_108 int
	var bindings_109 vm.Value
	var fs_110 vm.Value
	var last_id_111 vm.Value
	var ctx_112 vm.Value
	var arg__3314_114 vm.Value
	var arg__3325_120 vm.Value
	var arg__3332_125 vm.Value
	var arg__3334_127 vm.Value
	var arg__3346_134 vm.Value
	var arg__3353_139 vm.Value
	var arg__3355_141 vm.Value
	var arg__3356_142 vm.Value
	var arg__3359_145 vm.Value
	var arg__3370_151 vm.Value
	var arg__3377_156 vm.Value
	var arg__3379_158 vm.Value
	var arg__3391_165 vm.Value
	var arg__3398_170 vm.Value
	var arg__3400_172 vm.Value
	var arg__3401_173 vm.Value
	var let_syms_174 vm.Value
	var post_locals_176 vm.Value
	var v178 vm.Value
	var doseq_seq__3240_180 vm.Value
	var doseq_loop__3241_181 vm.Value
	var ctx_182 vm.Value
	var let_syms_183 vm.Value
	var v393 int
	var v402 vm.Value
	var v411 int
	var result_185 vm.Value
	var form_186 vm.Value
	var vec__3239_187 vm.Value
	var __188 vm.Value
	var body_189 vm.Value
	var i_190 int
	var bindings_191 vm.Value
	var fs_192 vm.Value
	var last_id_193 vm.Value
	var post_locals_194 vm.Value
	var doseq_seq__3240_195 vm.Value
	var doseq_loop__3241_196 vm.Value
	var ctx_197 vm.Value
	var let_syms_198 vm.Value
	var v391 int
	var v400 vm.Value
	var v409 int
	var vec__3242_215 vm.Value
	var sym_221 vm.Value
	var val_227 vm.Value
	var arg__3431_262 vm.Value
	var arg__3436_264 vm.Value
	var and__x_265 vm.Value
	var result_199 vm.Value
	var form_200 vm.Value
	var vec__3239_201 vm.Value
	var __202 vm.Value
	var body_203 vm.Value
	var i_204 int
	var bindings_205 vm.Value
	var fs_206 vm.Value
	var last_id_207 vm.Value
	var post_locals_208 vm.Value
	var doseq_seq__3240_209 vm.Value
	var doseq_loop__3241_210 vm.Value
	var ctx_211 vm.Value
	var let_syms_212 vm.Value
	var v396 int
	var v405 vm.Value
	var v414 int
	var v359 vm.Value
	var result_360 vm.Value
	var form_361 vm.Value
	var vec__3239_362 vm.Value
	var __363 vm.Value
	var body_364 vm.Value
	var i_365 int
	var bindings_366 vm.Value
	var fs_367 vm.Value
	var last_id_368 vm.Value
	var post_locals_369 vm.Value
	var doseq_seq__3240_370 vm.Value
	var doseq_loop__3241_371 vm.Value
	var ctx_372 vm.Value
	var let_syms_373 vm.Value
	var result_228 vm.Value
	var form_229 vm.Value
	var vec__3239_230 vm.Value
	var __231 vm.Value
	var body_232 vm.Value
	var i_233 int
	var bindings_234 vm.Value
	var fs_235 vm.Value
	var last_id_236 vm.Value
	var post_locals_237 vm.Value
	var doseq_seq__3240_238 vm.Value
	var doseq_loop__3241_239 vm.Value
	var ctx_240 vm.Value
	var let_syms_241 vm.Value
	var vec__3242_242 vm.Value
	var sym_243 vm.Value
	var val_244 vm.Value
	var v389 int
	var v398 vm.Value
	var v407 int
	var v332 vm.Value
	var result_245 vm.Value
	var form_246 vm.Value
	var vec__3239_247 vm.Value
	var __248 vm.Value
	var body_249 vm.Value
	var i_250 int
	var bindings_251 vm.Value
	var fs_252 vm.Value
	var last_id_253 vm.Value
	var post_locals_254 vm.Value
	var doseq_seq__3240_255 vm.Value
	var doseq_loop__3241_256 vm.Value
	var ctx_257 vm.Value
	var let_syms_258 vm.Value
	var vec__3242_259 vm.Value
	var sym_260 vm.Value
	var val_261 vm.Value
	var v395 int
	var v404 vm.Value
	var v413 int
	var v336 vm.Value
	var result_337 vm.Value
	var form_338 vm.Value
	var vec__3239_339 vm.Value
	var __340 vm.Value
	var body_341 vm.Value
	var i_342 int
	var bindings_343 vm.Value
	var fs_344 vm.Value
	var last_id_345 vm.Value
	var post_locals_346 vm.Value
	var doseq_seq__3240_347 vm.Value
	var doseq_loop__3241_348 vm.Value
	var ctx_349 vm.Value
	var let_syms_350 vm.Value
	var vec__3242_351 vm.Value
	var sym_352 vm.Value
	var val_353 vm.Value
	var v390 int
	var v399 vm.Value
	var v408 int
	var v355 vm.Value
	var result_266 vm.Value
	var form_267 vm.Value
	var vec__3239_268 vm.Value
	var __269 vm.Value
	var body_270 vm.Value
	var i_271 int
	var bindings_272 vm.Value
	var fs_273 vm.Value
	var last_id_274 vm.Value
	var post_locals_275 vm.Value
	var doseq_seq__3240_276 vm.Value
	var doseq_loop__3241_277 vm.Value
	var ctx_278 vm.Value
	var let_syms_279 vm.Value
	var vec__3242_280 vm.Value
	var sym_281 vm.Value
	var val_282 vm.Value
	var and__x_283 vm.Value
	var v388 int
	var v397 vm.Value
	var v406 int
	var arg__3443_304 vm.Value
	var arg__3451_307 vm.Value
	var v308 vm.Value
	var result_284 vm.Value
	var form_285 vm.Value
	var vec__3239_286 vm.Value
	var __287 vm.Value
	var body_288 vm.Value
	var i_289 int
	var bindings_290 vm.Value
	var fs_291 vm.Value
	var last_id_292 vm.Value
	var post_locals_293 vm.Value
	var doseq_seq__3240_294 vm.Value
	var doseq_loop__3241_295 vm.Value
	var ctx_296 vm.Value
	var let_syms_297 vm.Value
	var vec__3242_298 vm.Value
	var sym_299 vm.Value
	var val_300 vm.Value
	var and__x_301 vm.Value
	var v394 int
	var v403 vm.Value
	var v412 int
	var v311 vm.Value
	var result_312 vm.Value
	var form_313 vm.Value
	var vec__3239_314 vm.Value
	var __315 vm.Value
	var body_316 vm.Value
	var i_317 int
	var bindings_318 vm.Value
	var fs_319 vm.Value
	var last_id_320 vm.Value
	var post_locals_321 vm.Value
	var doseq_seq__3240_322 vm.Value
	var doseq_loop__3241_323 vm.Value
	var ctx_324 vm.Value
	var let_syms_325 vm.Value
	var vec__3242_326 vm.Value
	var sym_327 vm.Value
	var val_328 vm.Value
	var and__x_329 vm.Value
	var v392 int
	var v401 vm.Value
	var v410 int
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = __8, bindings_14, body_18, v20, i_21, bindings_22, ctx_23, v377, arg__3269_41, v42, form_26, vec__3239_27, __28, body_29, i_30, bindings_31, ctx_32, v376, sym_45, arg__3277_46, expr_49, expr_id_51, v53, v55, form_33, vec__3239_34, __35, body_36, i_37, bindings_38, ctx_39, v378, v59, form_60, vec__3239_61, __62, body_63, i_64, bindings_65, ctx_66, fs_67, last_id_68, ctx_69, v91, form_72, vec__3239_73, __74, body_75, i_76, bindings_77, fs_78, last_id_79, ctx_80, v94, arg__3305_96, arg__3311_99, v100, form_81, vec__3239_82, __83, body_84, i_85, bindings_86, fs_87, last_id_88, ctx_89, result_103, form_104, vec__3239_105, __106, body_107, i_108, bindings_109, fs_110, last_id_111, ctx_112, arg__3314_114, arg__3325_120, arg__3332_125, arg__3334_127, arg__3346_134, arg__3353_139, arg__3355_141, arg__3356_142, arg__3359_145, arg__3370_151, arg__3377_156, arg__3379_158, arg__3391_165, arg__3398_170, arg__3400_172, arg__3401_173, let_syms_174, post_locals_176, v178, doseq_seq__3240_180, doseq_loop__3241_181, ctx_182, let_syms_183, v393, v402, v411, result_185, form_186, vec__3239_187, __188, body_189, i_190, bindings_191, fs_192, last_id_193, post_locals_194, doseq_seq__3240_195, doseq_loop__3241_196, ctx_197, let_syms_198, v391, v400, v409, vec__3242_215, sym_221, val_227, arg__3431_262, arg__3436_264, and__x_265, result_199, form_200, vec__3239_201, __202, body_203, i_204, bindings_205, fs_206, last_id_207, post_locals_208, doseq_seq__3240_209, doseq_loop__3241_210, ctx_211, let_syms_212, v396, v405, v414, v359, result_360, form_361, vec__3239_362, __363, body_364, i_365, bindings_366, fs_367, last_id_368, post_locals_369, doseq_seq__3240_370, doseq_loop__3241_371, ctx_372, let_syms_373, result_228, form_229, vec__3239_230, __231, body_232, i_233, bindings_234, fs_235, last_id_236, post_locals_237, doseq_seq__3240_238, doseq_loop__3241_239, ctx_240, let_syms_241, vec__3242_242, sym_243, val_244, v389, v398, v407, v332, result_245, form_246, vec__3239_247, __248, body_249, i_250, bindings_251, fs_252, last_id_253, post_locals_254, doseq_seq__3240_255, doseq_loop__3241_256, ctx_257, let_syms_258, vec__3242_259, sym_260, val_261, v395, v404, v413, v336, result_337, form_338, vec__3239_339, __340, body_341, i_342, bindings_343, fs_344, last_id_345, post_locals_346, doseq_seq__3240_347, doseq_loop__3241_348, ctx_349, let_syms_350, vec__3242_351, sym_352, val_353, v390, v399, v408, v355, result_266, form_267, vec__3239_268, __269, body_270, i_271, bindings_272, fs_273, last_id_274, post_locals_275, doseq_seq__3240_276, doseq_loop__3241_277, ctx_278, let_syms_279, vec__3242_280, sym_281, val_282, and__x_283, v388, v397, v406, arg__3443_304, arg__3451_307, v308, result_284, form_285, vec__3239_286, __287, body_288, i_289, bindings_290, fs_291, last_id_292, post_locals_293, doseq_seq__3240_294, doseq_loop__3241_295, ctx_296, let_syms_297, vec__3242_298, sym_299, val_300, and__x_301, v394, v403, v412, v311, result_312, form_313, vec__3239_314, __315, body_316, i_317, bindings_318, fs_319, last_id_320, post_locals_321, doseq_seq__3240_322, doseq_loop__3241_323, ctx_324, let_syms_325, vec__3242_326, sym_327, val_328, and__x_329, v392, v401, v410
	__8, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(0), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	bindings_14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(1), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	body_18, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "drop").Deref(), []vm.Value{vm.Int(2), arg0})
	if callErr != nil {
		return nil, callErr
	}
	v20, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "push-locals!").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	i_21 = 0
	bindings_22 = bindings_14
	ctx_23 = arg1
	v377 = 2
	goto b1
b1:
	;
	arg__3269_41, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{bindings_22})
	if callErr != nil {
		return nil, callErr
	}
	v42 = rt.LtValue(vm.Int(i_21), arg__3269_41)
	if v42 {
		form_26 = arg0
		vec__3239_27 = arg0
		__28 = __8
		body_29 = body_18
		i_30 = i_21
		bindings_31 = bindings_22
		ctx_32 = ctx_23
		v376 = v377
		goto b2
	} else {
		form_33 = arg0
		vec__3239_34 = arg0
		__35 = __8
		body_36 = body_18
		i_37 = i_21
		bindings_38 = bindings_22
		ctx_39 = ctx_23
		v378 = v377
		goto b3
	}
b2:
	;
	sym_45, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{bindings_31, vm.Int(i_30)})
	if callErr != nil {
		return nil, callErr
	}
	arg__3277_46 = i_30 + 1
	expr_49, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{bindings_31, vm.Int(arg__3277_46)})
	if callErr != nil {
		return nil, callErr
	}
	expr_id_51, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-form").Deref(), []vm.Value{expr_49, ctx_32})
	if callErr != nil {
		return nil, callErr
	}
	v53, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "bind-local!").Deref(), []vm.Value{ctx_32, sym_45, expr_id_51})
	if callErr != nil {
		return nil, callErr
	}
	v55 = i_30 + v376
	i_21 = v55
	bindings_22 = bindings_31
	ctx_23 = ctx_32
	v377 = v376
	goto b1
b3:
	;
	v59 = vm.NIL
	form_60 = form_33
	vec__3239_61 = vec__3239_34
	__62 = __35
	body_63 = body_36
	i_64 = i_37
	bindings_65 = bindings_38
	ctx_66 = ctx_39
	goto b4
b4:
	;
	fs_67 = body_63
	last_id_68 = vm.NIL
	ctx_69 = ctx_66
	goto b5
b5:
	;
	v91, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{fs_67})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v91) {
		form_72 = form_60
		vec__3239_73 = vec__3239_61
		__74 = __62
		body_75 = body_63
		i_76 = i_64
		bindings_77 = bindings_65
		fs_78 = fs_67
		last_id_79 = last_id_68
		ctx_80 = ctx_69
		goto b6
	} else {
		form_81 = form_60
		vec__3239_82 = vec__3239_61
		__83 = __62
		body_84 = body_63
		i_85 = i_64
		bindings_86 = bindings_65
		fs_87 = fs_67
		last_id_88 = last_id_68
		ctx_89 = ctx_69
		goto b7
	}
b6:
	;
	v94, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{fs_78})
	if callErr != nil {
		return nil, callErr
	}
	arg__3305_96, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{fs_78})
	if callErr != nil {
		return nil, callErr
	}
	arg__3311_99, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{fs_78})
	if callErr != nil {
		return nil, callErr
	}
	v100, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-form").Deref(), []vm.Value{arg__3311_99, ctx_80})
	if callErr != nil {
		return nil, callErr
	}
	fs_67 = v94
	last_id_68 = v100
	ctx_69 = ctx_80
	goto b5
b7:
	;
	result_103 = last_id_88
	form_104 = form_81
	vec__3239_105 = vec__3239_82
	__106 = __83
	body_107 = body_84
	i_108 = i_85
	bindings_109 = bindings_86
	fs_110 = fs_87
	last_id_111 = last_id_88
	ctx_112 = ctx_89
	goto b8
b8:
	;
	arg__3314_114, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__3325_120, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{bindings_109})
	if callErr != nil {
		return nil, callErr
	}
	arg__3332_125, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{bindings_109})
	if callErr != nil {
		return nil, callErr
	}
	arg__3334_127, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "range").Deref(), []vm.Value{vm.Int(0), arg__3332_125, vm.Int(2)})
	if callErr != nil {
		return nil, callErr
	}
	arg__3346_134, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{bindings_109})
	if callErr != nil {
		return nil, callErr
	}
	arg__3353_139, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{bindings_109})
	if callErr != nil {
		return nil, callErr
	}
	arg__3355_141, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "range").Deref(), []vm.Value{vm.Int(0), arg__3353_139, vm.Int(2)})
	if callErr != nil {
		return nil, callErr
	}
	arg__3356_142, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{bindings_109, arg0})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), arg__3355_141})
	if callErr != nil {
		return nil, callErr
	}
	arg__3359_145, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__3370_151, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{bindings_109})
	if callErr != nil {
		return nil, callErr
	}
	arg__3377_156, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{bindings_109})
	if callErr != nil {
		return nil, callErr
	}
	arg__3379_158, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "range").Deref(), []vm.Value{vm.Int(0), arg__3377_156, vm.Int(2)})
	if callErr != nil {
		return nil, callErr
	}
	arg__3391_165, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{bindings_109})
	if callErr != nil {
		return nil, callErr
	}
	arg__3398_170, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{bindings_109})
	if callErr != nil {
		return nil, callErr
	}
	arg__3400_172, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "range").Deref(), []vm.Value{vm.Int(0), arg__3398_170, vm.Int(2)})
	if callErr != nil {
		return nil, callErr
	}
	arg__3401_173, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{bindings_109, arg0})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), arg__3400_172})
	if callErr != nil {
		return nil, callErr
	}
	let_syms_174, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__3359_145, arg__3401_173})
	if callErr != nil {
		return nil, callErr
	}
	post_locals_176, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "current-locals-flat").Deref(), []vm.Value{ctx_112})
	if callErr != nil {
		return nil, callErr
	}
	v178, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "pop-locals!").Deref(), []vm.Value{ctx_112})
	if callErr != nil {
		return nil, callErr
	}
	doseq_seq__3240_180, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{post_locals_176})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__3241_181 = doseq_seq__3240_180
	ctx_182 = ctx_112
	let_syms_183 = let_syms_174
	v393 = 0
	v402 = vm.NIL
	v411 = 1
	goto b9
b9:
	;
	if vm.IsTruthy(doseq_loop__3241_181) {
		result_185 = result_103
		form_186 = form_104
		vec__3239_187 = vec__3239_105
		__188 = __106
		body_189 = body_107
		i_190 = i_108
		bindings_191 = bindings_109
		fs_192 = fs_110
		last_id_193 = last_id_111
		post_locals_194 = post_locals_176
		doseq_seq__3240_195 = doseq_seq__3240_180
		doseq_loop__3241_196 = doseq_loop__3241_181
		ctx_197 = ctx_182
		let_syms_198 = let_syms_183
		v391 = v393
		v400 = v402
		v409 = v411
		goto b10
	} else {
		result_199 = result_103
		form_200 = form_104
		vec__3239_201 = vec__3239_105
		__202 = __106
		body_203 = body_107
		i_204 = i_108
		bindings_205 = bindings_109
		fs_206 = fs_110
		last_id_207 = last_id_111
		post_locals_208 = post_locals_176
		doseq_seq__3240_209 = doseq_seq__3240_180
		doseq_loop__3241_210 = doseq_loop__3241_181
		ctx_211 = ctx_182
		let_syms_212 = let_syms_183
		v396 = v393
		v405 = v402
		v414 = v411
		goto b11
	}
b10:
	;
	vec__3242_215, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{doseq_loop__3241_196})
	if callErr != nil {
		return nil, callErr
	}
	sym_221, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{vec__3242_215, vm.Int(v391), v400})
	if callErr != nil {
		return nil, callErr
	}
	val_227, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{vec__3242_215, vm.Int(v409), v400})
	if callErr != nil {
		return nil, callErr
	}
	arg__3431_262, callErr = rt.InvokeValue(let_syms_198, []vm.Value{sym_221})
	if callErr != nil {
		return nil, callErr
	}
	arg__3436_264, callErr = rt.InvokeValue(let_syms_198, []vm.Value{sym_221})
	if callErr != nil {
		return nil, callErr
	}
	and__x_265, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not").Deref(), []vm.Value{arg__3436_264})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(and__x_265) {
		result_266 = result_185
		form_267 = form_186
		vec__3239_268 = vec__3239_187
		__269 = __188
		body_270 = body_189
		i_271 = i_190
		bindings_272 = bindings_191
		fs_273 = fs_192
		last_id_274 = last_id_193
		post_locals_275 = post_locals_194
		doseq_seq__3240_276 = doseq_seq__3240_195
		doseq_loop__3241_277 = doseq_loop__3241_196
		ctx_278 = ctx_197
		let_syms_279 = let_syms_198
		vec__3242_280 = vec__3242_215
		sym_281 = sym_221
		val_282 = val_227
		and__x_283 = and__x_265
		v388 = v391
		v397 = v400
		v406 = v409
		goto b16
	} else {
		result_284 = result_185
		form_285 = form_186
		vec__3239_286 = vec__3239_187
		__287 = __188
		body_288 = body_189
		i_289 = i_190
		bindings_290 = bindings_191
		fs_291 = fs_192
		last_id_292 = last_id_193
		post_locals_293 = post_locals_194
		doseq_seq__3240_294 = doseq_seq__3240_195
		doseq_loop__3241_295 = doseq_loop__3241_196
		ctx_296 = ctx_197
		let_syms_297 = let_syms_198
		vec__3242_298 = vec__3242_215
		sym_299 = sym_221
		val_300 = val_227
		and__x_301 = and__x_265
		v394 = v391
		v403 = v400
		v412 = v409
		goto b17
	}
b11:
	;
	v359 = vm.NIL
	result_360 = result_199
	form_361 = form_200
	vec__3239_362 = vec__3239_201
	__363 = __202
	body_364 = body_203
	i_365 = i_204
	bindings_366 = bindings_205
	fs_367 = fs_206
	last_id_368 = last_id_207
	post_locals_369 = post_locals_208
	doseq_seq__3240_370 = doseq_seq__3240_209
	doseq_loop__3241_371 = doseq_loop__3241_210
	ctx_372 = ctx_211
	let_syms_373 = let_syms_212
	goto b12
b12:
	;
	return result_360, nil
b13:
	;
	v332, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "bind-local!").Deref(), []vm.Value{ctx_240, sym_243, val_244})
	if callErr != nil {
		return nil, callErr
	}
	v336 = v332
	result_337 = result_228
	form_338 = form_229
	vec__3239_339 = vec__3239_230
	__340 = __231
	body_341 = body_232
	i_342 = i_233
	bindings_343 = bindings_234
	fs_344 = fs_235
	last_id_345 = last_id_236
	post_locals_346 = post_locals_237
	doseq_seq__3240_347 = doseq_seq__3240_238
	doseq_loop__3241_348 = doseq_loop__3241_239
	ctx_349 = ctx_240
	let_syms_350 = let_syms_241
	vec__3242_351 = vec__3242_242
	sym_352 = sym_243
	val_353 = val_244
	v390 = v389
	v399 = v398
	v408 = v407
	goto b15
b14:
	;
	v336 = vm.NIL
	result_337 = result_245
	form_338 = form_246
	vec__3239_339 = vec__3239_247
	__340 = __248
	body_341 = body_249
	i_342 = i_250
	bindings_343 = bindings_251
	fs_344 = fs_252
	last_id_345 = last_id_253
	post_locals_346 = post_locals_254
	doseq_seq__3240_347 = doseq_seq__3240_255
	doseq_loop__3241_348 = doseq_loop__3241_256
	ctx_349 = ctx_257
	let_syms_350 = let_syms_258
	vec__3242_351 = vec__3242_259
	sym_352 = sym_260
	val_353 = val_261
	v390 = v395
	v399 = v404
	v408 = v413
	goto b15
b15:
	;
	v355, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "next").Deref(), []vm.Value{doseq_loop__3241_348})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__3241_181 = v355
	ctx_182 = ctx_349
	let_syms_183 = let_syms_350
	v393 = v390
	v402 = v399
	v411 = v408
	goto b9
b16:
	;
	arg__3443_304, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "lookup-local").Deref(), []vm.Value{ctx_278, sym_281})
	if callErr != nil {
		return nil, callErr
	}
	arg__3451_307, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "lookup-local").Deref(), []vm.Value{ctx_278, sym_281})
	if callErr != nil {
		return nil, callErr
	}
	v308, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not=").Deref(), []vm.Value{val_282, arg__3451_307})
	if callErr != nil {
		return nil, callErr
	}
	v311 = v308
	result_312 = result_266
	form_313 = form_267
	vec__3239_314 = vec__3239_268
	__315 = __269
	body_316 = body_270
	i_317 = i_271
	bindings_318 = bindings_272
	fs_319 = fs_273
	last_id_320 = last_id_274
	post_locals_321 = post_locals_275
	doseq_seq__3240_322 = doseq_seq__3240_276
	doseq_loop__3241_323 = doseq_loop__3241_277
	ctx_324 = ctx_278
	let_syms_325 = let_syms_279
	vec__3242_326 = vec__3242_280
	sym_327 = sym_281
	val_328 = val_282
	and__x_329 = and__x_283
	v392 = v388
	v401 = v397
	v410 = v406
	goto b18
b17:
	;
	v311 = and__x_301
	result_312 = result_284
	form_313 = form_285
	vec__3239_314 = vec__3239_286
	__315 = __287
	body_316 = body_288
	i_317 = i_289
	bindings_318 = bindings_290
	fs_319 = fs_291
	last_id_320 = last_id_292
	post_locals_321 = post_locals_293
	doseq_seq__3240_322 = doseq_seq__3240_294
	doseq_loop__3241_323 = doseq_loop__3241_295
	ctx_324 = ctx_296
	let_syms_325 = let_syms_297
	vec__3242_326 = vec__3242_298
	sym_327 = sym_299
	val_328 = val_300
	and__x_329 = and__x_301
	v392 = v394
	v401 = v403
	v410 = v412
	goto b18
b18:
	;
	if vm.IsTruthy(v311) {
		result_228 = result_312
		form_229 = form_313
		vec__3239_230 = vec__3239_314
		__231 = __315
		body_232 = body_316
		i_233 = i_317
		bindings_234 = bindings_318
		fs_235 = fs_319
		last_id_236 = last_id_320
		post_locals_237 = post_locals_321
		doseq_seq__3240_238 = doseq_seq__3240_322
		doseq_loop__3241_239 = doseq_loop__3241_323
		ctx_240 = ctx_324
		let_syms_241 = let_syms_325
		vec__3242_242 = vec__3242_326
		sym_243 = sym_327
		val_244 = val_328
		v389 = v392
		v398 = v401
		v407 = v410
		goto b13
	} else {
		result_245 = result_312
		form_246 = form_313
		vec__3239_247 = vec__3239_314
		__248 = __315
		body_249 = body_316
		i_250 = i_317
		bindings_251 = bindings_318
		fs_252 = fs_319
		last_id_253 = last_id_320
		post_locals_254 = post_locals_321
		doseq_seq__3240_255 = doseq_seq__3240_322
		doseq_loop__3241_256 = doseq_loop__3241_323
		ctx_257 = ctx_324
		let_syms_258 = let_syms_325
		vec__3242_259 = vec__3242_326
		sym_260 = sym_327
		val_261 = val_328
		v395 = v392
		v404 = v401
		v413 = v410
		goto b14
	}
}
func build_list(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var head_3 vm.Value
	var builtin_7 vm.Value
	var v17 bool
	var form_8 vm.Value
	var ctx_9 vm.Value
	var head_10 vm.Value
	var builtin_11 vm.Value
	var v20 vm.Value
	var form_12 vm.Value
	var ctx_13 vm.Value
	var head_14 vm.Value
	var builtin_15 vm.Value
	var v31 bool
	var v329 vm.Value
	var form_330 vm.Value
	var ctx_331 vm.Value
	var head_332 vm.Value
	var builtin_333 vm.Value
	var form_22 vm.Value
	var ctx_23 vm.Value
	var head_24 vm.Value
	var builtin_25 vm.Value
	var v34 vm.Value
	var form_26 vm.Value
	var ctx_27 vm.Value
	var head_28 vm.Value
	var builtin_29 vm.Value
	var v45 bool
	var v323 vm.Value
	var form_324 vm.Value
	var ctx_325 vm.Value
	var head_326 vm.Value
	var builtin_327 vm.Value
	var form_36 vm.Value
	var ctx_37 vm.Value
	var head_38 vm.Value
	var builtin_39 vm.Value
	var v48 vm.Value
	var form_40 vm.Value
	var ctx_41 vm.Value
	var head_42 vm.Value
	var builtin_43 vm.Value
	var v59 bool
	var v317 vm.Value
	var form_318 vm.Value
	var ctx_319 vm.Value
	var head_320 vm.Value
	var builtin_321 vm.Value
	var form_50 vm.Value
	var ctx_51 vm.Value
	var head_52 vm.Value
	var builtin_53 vm.Value
	var v62 vm.Value
	var form_54 vm.Value
	var ctx_55 vm.Value
	var head_56 vm.Value
	var builtin_57 vm.Value
	var v73 bool
	var v311 vm.Value
	var form_312 vm.Value
	var ctx_313 vm.Value
	var head_314 vm.Value
	var builtin_315 vm.Value
	var form_64 vm.Value
	var ctx_65 vm.Value
	var head_66 vm.Value
	var builtin_67 vm.Value
	var v76 vm.Value
	var form_68 vm.Value
	var ctx_69 vm.Value
	var head_70 vm.Value
	var builtin_71 vm.Value
	var v87 bool
	var v305 vm.Value
	var form_306 vm.Value
	var ctx_307 vm.Value
	var head_308 vm.Value
	var builtin_309 vm.Value
	var form_78 vm.Value
	var ctx_79 vm.Value
	var head_80 vm.Value
	var builtin_81 vm.Value
	var v90 vm.Value
	var form_82 vm.Value
	var ctx_83 vm.Value
	var head_84 vm.Value
	var builtin_85 vm.Value
	var v101 bool
	var v299 vm.Value
	var form_300 vm.Value
	var ctx_301 vm.Value
	var head_302 vm.Value
	var builtin_303 vm.Value
	var form_92 vm.Value
	var ctx_93 vm.Value
	var head_94 vm.Value
	var builtin_95 vm.Value
	var v104 vm.Value
	var form_96 vm.Value
	var ctx_97 vm.Value
	var head_98 vm.Value
	var builtin_99 vm.Value
	var v115 bool
	var v293 vm.Value
	var form_294 vm.Value
	var ctx_295 vm.Value
	var head_296 vm.Value
	var builtin_297 vm.Value
	var form_106 vm.Value
	var ctx_107 vm.Value
	var head_108 vm.Value
	var builtin_109 vm.Value
	var v118 vm.Value
	var form_110 vm.Value
	var ctx_111 vm.Value
	var head_112 vm.Value
	var builtin_113 vm.Value
	var v129 bool
	var v287 vm.Value
	var form_288 vm.Value
	var ctx_289 vm.Value
	var head_290 vm.Value
	var builtin_291 vm.Value
	var form_120 vm.Value
	var ctx_121 vm.Value
	var head_122 vm.Value
	var builtin_123 vm.Value
	var v132 vm.Value
	var form_124 vm.Value
	var ctx_125 vm.Value
	var head_126 vm.Value
	var builtin_127 vm.Value
	var v143 bool
	var v281 vm.Value
	var form_282 vm.Value
	var ctx_283 vm.Value
	var head_284 vm.Value
	var builtin_285 vm.Value
	var form_134 vm.Value
	var ctx_135 vm.Value
	var head_136 vm.Value
	var builtin_137 vm.Value
	var v146 vm.Value
	var form_138 vm.Value
	var ctx_139 vm.Value
	var head_140 vm.Value
	var builtin_141 vm.Value
	var v157 bool
	var v275 vm.Value
	var form_276 vm.Value
	var ctx_277 vm.Value
	var head_278 vm.Value
	var builtin_279 vm.Value
	var form_148 vm.Value
	var ctx_149 vm.Value
	var head_150 vm.Value
	var builtin_151 vm.Value
	var v160 vm.Value
	var form_152 vm.Value
	var ctx_153 vm.Value
	var head_154 vm.Value
	var builtin_155 vm.Value
	var v269 vm.Value
	var form_270 vm.Value
	var ctx_271 vm.Value
	var head_272 vm.Value
	var builtin_273 vm.Value
	var form_162 vm.Value
	var ctx_163 vm.Value
	var head_164 vm.Value
	var builtin_165 vm.Value
	var v172 vm.Value
	var form_166 vm.Value
	var ctx_167 vm.Value
	var head_168 vm.Value
	var builtin_169 vm.Value
	var or__x_183 vm.Value
	var v263 vm.Value
	var form_264 vm.Value
	var ctx_265 vm.Value
	var head_266 vm.Value
	var builtin_267 vm.Value
	var form_174 vm.Value
	var ctx_175 vm.Value
	var head_176 vm.Value
	var builtin_177 vm.Value
	var fn_id_228 vm.Value
	var arg__3572_230 vm.Value
	var arg__3579_233 vm.Value
	var v234 vm.Value
	var form_178 vm.Value
	var ctx_179 vm.Value
	var head_180 vm.Value
	var builtin_181 vm.Value
	var v257 vm.Value
	var form_258 vm.Value
	var ctx_259 vm.Value
	var head_260 vm.Value
	var builtin_261 vm.Value
	var form_184 vm.Value
	var ctx_185 vm.Value
	var head_186 vm.Value
	var builtin_187 vm.Value
	var or__x_188 vm.Value
	var form_189 vm.Value
	var ctx_190 vm.Value
	var head_191 vm.Value
	var builtin_192 vm.Value
	var or__x_193 vm.Value
	var or__x_197 vm.Value
	var v220 vm.Value
	var form_221 vm.Value
	var ctx_222 vm.Value
	var head_223 vm.Value
	var builtin_224 vm.Value
	var or__x_225 vm.Value
	var form_198 vm.Value
	var ctx_199 vm.Value
	var head_200 vm.Value
	var builtin_201 vm.Value
	var or__x_202 vm.Value
	var form_203 vm.Value
	var ctx_204 vm.Value
	var head_205 vm.Value
	var builtin_206 vm.Value
	var or__x_207 vm.Value
	var v211 vm.Value
	var v213 vm.Value
	var form_214 vm.Value
	var ctx_215 vm.Value
	var head_216 vm.Value
	var builtin_217 vm.Value
	var or__x_218 vm.Value
	var form_236 vm.Value
	var ctx_237 vm.Value
	var head_238 vm.Value
	var builtin_239 vm.Value
	var v247 vm.Value
	var form_240 vm.Value
	var ctx_241 vm.Value
	var head_242 vm.Value
	var builtin_243 vm.Value
	var v251 vm.Value
	var form_252 vm.Value
	var ctx_253 vm.Value
	var head_254 vm.Value
	var builtin_255 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = head_3, builtin_7, v17, form_8, ctx_9, head_10, builtin_11, v20, form_12, ctx_13, head_14, builtin_15, v31, v329, form_330, ctx_331, head_332, builtin_333, form_22, ctx_23, head_24, builtin_25, v34, form_26, ctx_27, head_28, builtin_29, v45, v323, form_324, ctx_325, head_326, builtin_327, form_36, ctx_37, head_38, builtin_39, v48, form_40, ctx_41, head_42, builtin_43, v59, v317, form_318, ctx_319, head_320, builtin_321, form_50, ctx_51, head_52, builtin_53, v62, form_54, ctx_55, head_56, builtin_57, v73, v311, form_312, ctx_313, head_314, builtin_315, form_64, ctx_65, head_66, builtin_67, v76, form_68, ctx_69, head_70, builtin_71, v87, v305, form_306, ctx_307, head_308, builtin_309, form_78, ctx_79, head_80, builtin_81, v90, form_82, ctx_83, head_84, builtin_85, v101, v299, form_300, ctx_301, head_302, builtin_303, form_92, ctx_93, head_94, builtin_95, v104, form_96, ctx_97, head_98, builtin_99, v115, v293, form_294, ctx_295, head_296, builtin_297, form_106, ctx_107, head_108, builtin_109, v118, form_110, ctx_111, head_112, builtin_113, v129, v287, form_288, ctx_289, head_290, builtin_291, form_120, ctx_121, head_122, builtin_123, v132, form_124, ctx_125, head_126, builtin_127, v143, v281, form_282, ctx_283, head_284, builtin_285, form_134, ctx_135, head_136, builtin_137, v146, form_138, ctx_139, head_140, builtin_141, v157, v275, form_276, ctx_277, head_278, builtin_279, form_148, ctx_149, head_150, builtin_151, v160, form_152, ctx_153, head_154, builtin_155, v269, form_270, ctx_271, head_272, builtin_273, form_162, ctx_163, head_164, builtin_165, v172, form_166, ctx_167, head_168, builtin_169, or__x_183, v263, form_264, ctx_265, head_266, builtin_267, form_174, ctx_175, head_176, builtin_177, fn_id_228, arg__3572_230, arg__3579_233, v234, form_178, ctx_179, head_180, builtin_181, v257, form_258, ctx_259, head_260, builtin_261, form_184, ctx_185, head_186, builtin_187, or__x_188, form_189, ctx_190, head_191, builtin_192, or__x_193, or__x_197, v220, form_221, ctx_222, head_223, builtin_224, or__x_225, form_198, ctx_199, head_200, builtin_201, or__x_202, form_203, ctx_204, head_205, builtin_206, or__x_207, v211, v213, form_214, ctx_215, head_216, builtin_217, or__x_218, form_236, ctx_237, head_238, builtin_239, v247, form_240, ctx_241, head_242, builtin_243, v251, form_252, ctx_253, head_254, builtin_255
	head_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	builtin_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{rt.LookupVar("ir.build", "builtin-ops").Deref(), head_3})
	if callErr != nil {
		return nil, callErr
	}
	v17 = head_3 == vm.Symbol("if")
	if v17 {
		form_8 = arg0
		ctx_9 = arg1
		head_10 = head_3
		builtin_11 = builtin_7
		goto b1
	} else {
		form_12 = arg0
		ctx_13 = arg1
		head_14 = head_3
		builtin_15 = builtin_7
		goto b2
	}
b1:
	;
	v20, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-if").Deref(), []vm.Value{form_8, ctx_9})
	if callErr != nil {
		return nil, callErr
	}
	v329 = v20
	form_330 = form_8
	ctx_331 = ctx_9
	head_332 = head_10
	builtin_333 = builtin_11
	goto b3
b2:
	;
	v31 = head_14 == vm.Symbol("let")
	if v31 {
		form_22 = form_12
		ctx_23 = ctx_13
		head_24 = head_14
		builtin_25 = builtin_15
		goto b4
	} else {
		form_26 = form_12
		ctx_27 = ctx_13
		head_28 = head_14
		builtin_29 = builtin_15
		goto b5
	}
b3:
	;
	return v329, nil
b4:
	;
	v34, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-let").Deref(), []vm.Value{form_22, ctx_23})
	if callErr != nil {
		return nil, callErr
	}
	v323 = v34
	form_324 = form_22
	ctx_325 = ctx_23
	head_326 = head_24
	builtin_327 = builtin_25
	goto b6
b5:
	;
	v45 = head_28 == vm.Symbol("let*")
	if v45 {
		form_36 = form_26
		ctx_37 = ctx_27
		head_38 = head_28
		builtin_39 = builtin_29
		goto b7
	} else {
		form_40 = form_26
		ctx_41 = ctx_27
		head_42 = head_28
		builtin_43 = builtin_29
		goto b8
	}
b6:
	;
	v329 = v323
	form_330 = form_324
	ctx_331 = ctx_325
	head_332 = head_326
	builtin_333 = builtin_327
	goto b3
b7:
	;
	v48, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-let").Deref(), []vm.Value{form_36, ctx_37})
	if callErr != nil {
		return nil, callErr
	}
	v317 = v48
	form_318 = form_36
	ctx_319 = ctx_37
	head_320 = head_38
	builtin_321 = builtin_39
	goto b9
b8:
	;
	v59 = head_42 == vm.Symbol("do")
	if v59 {
		form_50 = form_40
		ctx_51 = ctx_41
		head_52 = head_42
		builtin_53 = builtin_43
		goto b10
	} else {
		form_54 = form_40
		ctx_55 = ctx_41
		head_56 = head_42
		builtin_57 = builtin_43
		goto b11
	}
b9:
	;
	v323 = v317
	form_324 = form_318
	ctx_325 = ctx_319
	head_326 = head_320
	builtin_327 = builtin_321
	goto b6
b10:
	;
	v62, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-do").Deref(), []vm.Value{form_50, ctx_51})
	if callErr != nil {
		return nil, callErr
	}
	v311 = v62
	form_312 = form_50
	ctx_313 = ctx_51
	head_314 = head_52
	builtin_315 = builtin_53
	goto b12
b11:
	;
	v73 = head_56 == vm.Symbol("quote")
	if v73 {
		form_64 = form_54
		ctx_65 = ctx_55
		head_66 = head_56
		builtin_67 = builtin_57
		goto b13
	} else {
		form_68 = form_54
		ctx_69 = ctx_55
		head_70 = head_56
		builtin_71 = builtin_57
		goto b14
	}
b12:
	;
	v317 = v311
	form_318 = form_312
	ctx_319 = ctx_313
	head_320 = head_314
	builtin_321 = builtin_315
	goto b9
b13:
	;
	v76, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-quote").Deref(), []vm.Value{form_64, ctx_65})
	if callErr != nil {
		return nil, callErr
	}
	v305 = v76
	form_306 = form_64
	ctx_307 = ctx_65
	head_308 = head_66
	builtin_309 = builtin_67
	goto b15
b14:
	;
	v87 = head_70 == vm.Symbol("var")
	if v87 {
		form_78 = form_68
		ctx_79 = ctx_69
		head_80 = head_70
		builtin_81 = builtin_71
		goto b16
	} else {
		form_82 = form_68
		ctx_83 = ctx_69
		head_84 = head_70
		builtin_85 = builtin_71
		goto b17
	}
b15:
	;
	v311 = v305
	form_312 = form_306
	ctx_313 = ctx_307
	head_314 = head_308
	builtin_315 = builtin_309
	goto b12
b16:
	;
	v90, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-var").Deref(), []vm.Value{form_78, ctx_79})
	if callErr != nil {
		return nil, callErr
	}
	v299 = v90
	form_300 = form_78
	ctx_301 = ctx_79
	head_302 = head_80
	builtin_303 = builtin_81
	goto b18
b17:
	;
	v101 = head_84 == vm.Symbol("set!")
	if v101 {
		form_92 = form_82
		ctx_93 = ctx_83
		head_94 = head_84
		builtin_95 = builtin_85
		goto b19
	} else {
		form_96 = form_82
		ctx_97 = ctx_83
		head_98 = head_84
		builtin_99 = builtin_85
		goto b20
	}
b18:
	;
	v305 = v299
	form_306 = form_300
	ctx_307 = ctx_301
	head_308 = head_302
	builtin_309 = builtin_303
	goto b15
b19:
	;
	v104, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-set!").Deref(), []vm.Value{form_92, ctx_93})
	if callErr != nil {
		return nil, callErr
	}
	v293 = v104
	form_294 = form_92
	ctx_295 = ctx_93
	head_296 = head_94
	builtin_297 = builtin_95
	goto b21
b20:
	;
	v115 = head_98 == vm.Symbol("loop")
	if v115 {
		form_106 = form_96
		ctx_107 = ctx_97
		head_108 = head_98
		builtin_109 = builtin_99
		goto b22
	} else {
		form_110 = form_96
		ctx_111 = ctx_97
		head_112 = head_98
		builtin_113 = builtin_99
		goto b23
	}
b21:
	;
	v299 = v293
	form_300 = form_294
	ctx_301 = ctx_295
	head_302 = head_296
	builtin_303 = builtin_297
	goto b18
b22:
	;
	v118, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-loop").Deref(), []vm.Value{form_106, ctx_107})
	if callErr != nil {
		return nil, callErr
	}
	v287 = v118
	form_288 = form_106
	ctx_289 = ctx_107
	head_290 = head_108
	builtin_291 = builtin_109
	goto b24
b23:
	;
	v129 = head_112 == vm.Symbol("loop*")
	if v129 {
		form_120 = form_110
		ctx_121 = ctx_111
		head_122 = head_112
		builtin_123 = builtin_113
		goto b25
	} else {
		form_124 = form_110
		ctx_125 = ctx_111
		head_126 = head_112
		builtin_127 = builtin_113
		goto b26
	}
b24:
	;
	v293 = v287
	form_294 = form_288
	ctx_295 = ctx_289
	head_296 = head_290
	builtin_297 = builtin_291
	goto b21
b25:
	;
	v132, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-loop").Deref(), []vm.Value{form_120, ctx_121})
	if callErr != nil {
		return nil, callErr
	}
	v281 = v132
	form_282 = form_120
	ctx_283 = ctx_121
	head_284 = head_122
	builtin_285 = builtin_123
	goto b27
b26:
	;
	v143 = head_126 == vm.Symbol("recur")
	if v143 {
		form_134 = form_124
		ctx_135 = ctx_125
		head_136 = head_126
		builtin_137 = builtin_127
		goto b28
	} else {
		form_138 = form_124
		ctx_139 = ctx_125
		head_140 = head_126
		builtin_141 = builtin_127
		goto b29
	}
b27:
	;
	v287 = v281
	form_288 = form_282
	ctx_289 = ctx_283
	head_290 = head_284
	builtin_291 = builtin_285
	goto b24
b28:
	;
	v146, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-recur").Deref(), []vm.Value{form_134, ctx_135})
	if callErr != nil {
		return nil, callErr
	}
	v275 = v146
	form_276 = form_134
	ctx_277 = ctx_135
	head_278 = head_136
	builtin_279 = builtin_137
	goto b30
b29:
	;
	v157 = head_140 == vm.Symbol("fn*")
	if v157 {
		form_148 = form_138
		ctx_149 = ctx_139
		head_150 = head_140
		builtin_151 = builtin_141
		goto b31
	} else {
		form_152 = form_138
		ctx_153 = ctx_139
		head_154 = head_140
		builtin_155 = builtin_141
		goto b32
	}
b30:
	;
	v281 = v275
	form_282 = form_276
	ctx_283 = ctx_277
	head_284 = head_278
	builtin_285 = builtin_279
	goto b27
b31:
	;
	v160, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-fn*").Deref(), []vm.Value{form_148, ctx_149})
	if callErr != nil {
		return nil, callErr
	}
	v269 = v160
	form_270 = form_148
	ctx_271 = ctx_149
	head_272 = head_150
	builtin_273 = builtin_151
	goto b33
b32:
	;
	if vm.IsTruthy(builtin_155) {
		form_162 = form_152
		ctx_163 = ctx_153
		head_164 = head_154
		builtin_165 = builtin_155
		goto b34
	} else {
		form_166 = form_152
		ctx_167 = ctx_153
		head_168 = head_154
		builtin_169 = builtin_155
		goto b35
	}
b33:
	;
	v275 = v269
	form_276 = form_270
	ctx_277 = ctx_271
	head_278 = head_272
	builtin_279 = builtin_273
	goto b30
b34:
	;
	v172, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-builtin-op").Deref(), []vm.Value{builtin_165, form_162, ctx_163})
	if callErr != nil {
		return nil, callErr
	}
	v263 = v172
	form_264 = form_162
	ctx_265 = ctx_163
	head_266 = head_164
	builtin_267 = builtin_165
	goto b36
b35:
	;
	or__x_183, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "keyword?").Deref(), []vm.Value{head_168})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(or__x_183) {
		form_184 = form_166
		ctx_185 = ctx_167
		head_186 = head_168
		builtin_187 = builtin_169
		or__x_188 = or__x_183
		goto b40
	} else {
		form_189 = form_166
		ctx_190 = ctx_167
		head_191 = head_168
		builtin_192 = builtin_169
		or__x_193 = or__x_183
		goto b41
	}
b36:
	;
	v269 = v263
	form_270 = form_264
	ctx_271 = ctx_265
	head_272 = head_266
	builtin_273 = builtin_267
	goto b33
b37:
	;
	fn_id_228, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-form").Deref(), []vm.Value{head_176, ctx_175})
	if callErr != nil {
		return nil, callErr
	}
	arg__3572_230, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{form_174})
	if callErr != nil {
		return nil, callErr
	}
	arg__3579_233, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{form_174})
	if callErr != nil {
		return nil, callErr
	}
	v234, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-call-with-head").Deref(), []vm.Value{fn_id_228, arg__3579_233, ctx_175})
	if callErr != nil {
		return nil, callErr
	}
	v257 = v234
	form_258 = form_174
	ctx_259 = ctx_175
	head_260 = head_176
	builtin_261 = builtin_177
	goto b39
b38:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		form_236 = form_178
		ctx_237 = ctx_179
		head_238 = head_180
		builtin_239 = builtin_181
		goto b46
	} else {
		form_240 = form_178
		ctx_241 = ctx_179
		head_242 = head_180
		builtin_243 = builtin_181
		goto b47
	}
b39:
	;
	v263 = v257
	form_264 = form_258
	ctx_265 = ctx_259
	head_266 = head_260
	builtin_267 = builtin_261
	goto b36
b40:
	;
	v220 = or__x_188
	form_221 = form_184
	ctx_222 = ctx_185
	head_223 = head_186
	builtin_224 = builtin_187
	or__x_225 = or__x_188
	goto b42
b41:
	;
	or__x_197, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "number?").Deref(), []vm.Value{head_191})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(or__x_197) {
		form_198 = form_189
		ctx_199 = ctx_190
		head_200 = head_191
		builtin_201 = builtin_192
		or__x_202 = or__x_197
		goto b43
	} else {
		form_203 = form_189
		ctx_204 = ctx_190
		head_205 = head_191
		builtin_206 = builtin_192
		or__x_207 = or__x_197
		goto b44
	}
b42:
	;
	if vm.IsTruthy(v220) {
		form_174 = form_221
		ctx_175 = ctx_222
		head_176 = head_223
		builtin_177 = builtin_224
		goto b37
	} else {
		form_178 = form_221
		ctx_179 = ctx_222
		head_180 = head_223
		builtin_181 = builtin_224
		goto b38
	}
b43:
	;
	v213 = or__x_202
	form_214 = form_198
	ctx_215 = ctx_199
	head_216 = head_200
	builtin_217 = builtin_201
	or__x_218 = or__x_202
	goto b45
b44:
	;
	v211, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "string?").Deref(), []vm.Value{head_205})
	if callErr != nil {
		return nil, callErr
	}
	v213 = v211
	form_214 = form_203
	ctx_215 = ctx_204
	head_216 = head_205
	builtin_217 = builtin_206
	or__x_218 = or__x_207
	goto b45
b45:
	;
	v220 = v213
	form_221 = form_214
	ctx_222 = ctx_215
	head_223 = head_216
	builtin_224 = builtin_217
	or__x_225 = or__x_193
	goto b42
b46:
	;
	v247, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-call").Deref(), []vm.Value{form_236, ctx_237})
	if callErr != nil {
		return nil, callErr
	}
	v251 = v247
	form_252 = form_236
	ctx_253 = ctx_237
	head_254 = head_238
	builtin_255 = builtin_239
	goto b48
b47:
	;
	v251 = vm.NIL
	form_252 = form_240
	ctx_253 = ctx_241
	head_254 = head_242
	builtin_255 = builtin_243
	goto b48
b48:
	;
	v257 = v251
	form_258 = form_252
	ctx_259 = ctx_253
	head_260 = head_254
	builtin_261 = builtin_255
	goto b39
}
func build_loop(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var __8 vm.Value
	var bindings_14 vm.Value
	var body_18 vm.Value
	var f_20 vm.Value
	var arg__3620_22 vm.Value
	var arg__3626_26 vm.Value
	var n_slots_28 vm.Value
	var arg__3631_30 vm.Value
	var arg__3637_34 vm.Value
	var old_header_36 vm.Value
	var arg__3642_38 vm.Value
	var arg__3648_42 vm.Value
	var old_caps_44 vm.Value
	var arg__3653_46 vm.Value
	var arg__3658_49 vm.Value
	var arg__3659_50 vm.Value
	var arg__3664_53 vm.Value
	var arg__3669_56 vm.Value
	var arg__3670_57 vm.Value
	var known_58 vm.Value
	var arg__3688_63 vm.Value
	var arg__3708_69 vm.Value
	var arg__3710_70 vm.Value
	var arg__3729_76 vm.Value
	var arg__3749_82 vm.Value
	var arg__3751_83 vm.Value
	var body_caps_84 vm.Value
	var n_caps_86 vm.Value
	var header_88 vm.Value
	var arg__3780_95 vm.Value
	var arg__3804_103 vm.Value
	var loop_param_ids_104 vm.Value
	var arg__3831_112 vm.Value
	var arg__3859_121 vm.Value
	var cap_param_ids_122 vm.Value
	var arg__3891_128 vm.Value
	var arg__3924_135 vm.Value
	var init_vals_136 vm.Value
	var cap_vals_144 vm.Value
	var arg__3946_146 vm.Value
	var arg__3953_149 vm.Value
	var arg__3954_150 vm.Value
	var arg__3962_153 vm.Value
	var arg__3969_156 vm.Value
	var arg__3970_157 vm.Value
	var bt_158 vm.Value
	var entry_end_160 vm.Value
	var __166 vm.Value
	var __168 vm.Value
	var v170 vm.Value
	var v172 vm.Value
	var arg__4003_174 vm.Value
	var arg__4008_177 vm.Value
	var doseq_seq__3587_178 vm.Value
	var doseq_loop__3588_179 vm.Value
	var bindings_180 vm.Value
	var ctx_181 vm.Value
	var loop_param_ids_182 vm.Value
	var v948 int
	var form_184 vm.Value
	var vec__3586_185 vm.Value
	var body_186 vm.Value
	var f_187 vm.Value
	var n_slots_188 vm.Value
	var old_header_189 vm.Value
	var old_caps_190 vm.Value
	var known_191 vm.Value
	var body_caps_192 vm.Value
	var n_caps_193 vm.Value
	var header_194 vm.Value
	var cap_param_ids_195 vm.Value
	var init_vals_196 vm.Value
	var cap_vals_197 vm.Value
	var bt_198 vm.Value
	var entry_end_199 vm.Value
	var __200 vm.Value
	var doseq_seq__3587_201 vm.Value
	var doseq_loop__3588_202 vm.Value
	var bindings_203 vm.Value
	var ctx_204 vm.Value
	var loop_param_ids_205 vm.Value
	var v947 int
	var i_230 vm.Value
	var arg__4016_232 vm.Value
	var arg__4022_236 vm.Value
	var arg__4028_238 vm.Value
	var arg__4040_245 vm.Value
	var arg__4046_247 vm.Value
	var v248 vm.Value
	var v250 vm.Value
	var form_206 vm.Value
	var vec__3586_207 vm.Value
	var body_208 vm.Value
	var f_209 vm.Value
	var n_slots_210 vm.Value
	var old_header_211 vm.Value
	var old_caps_212 vm.Value
	var known_213 vm.Value
	var body_caps_214 vm.Value
	var n_caps_215 vm.Value
	var header_216 vm.Value
	var cap_param_ids_217 vm.Value
	var init_vals_218 vm.Value
	var cap_vals_219 vm.Value
	var bt_220 vm.Value
	var entry_end_221 vm.Value
	var __222 vm.Value
	var doseq_seq__3587_223 vm.Value
	var doseq_loop__3588_224 vm.Value
	var bindings_225 vm.Value
	var ctx_226 vm.Value
	var loop_param_ids_227 vm.Value
	var v949 int
	var v254 vm.Value
	var form_255 vm.Value
	var vec__3586_256 vm.Value
	var body_257 vm.Value
	var f_258 vm.Value
	var n_slots_259 vm.Value
	var old_header_260 vm.Value
	var old_caps_261 vm.Value
	var known_262 vm.Value
	var body_caps_263 vm.Value
	var n_caps_264 vm.Value
	var header_265 vm.Value
	var cap_param_ids_266 vm.Value
	var init_vals_267 vm.Value
	var cap_vals_268 vm.Value
	var bt_269 vm.Value
	var entry_end_270 vm.Value
	var __271 vm.Value
	var doseq_seq__3587_272 vm.Value
	var doseq_loop__3588_273 vm.Value
	var bindings_274 vm.Value
	var ctx_275 vm.Value
	var loop_param_ids_276 vm.Value
	var arg__4057_280 vm.Value
	var arg__4066_285 vm.Value
	var doseq_seq__3589_286 vm.Value
	var doseq_loop__3590_287 vm.Value
	var ctx_288 vm.Value
	var v959 int
	var v962 vm.Value
	var v965 int
	var form_290 vm.Value
	var vec__3586_291 vm.Value
	var body_292 vm.Value
	var f_293 vm.Value
	var n_slots_294 vm.Value
	var old_header_295 vm.Value
	var old_caps_296 vm.Value
	var known_297 vm.Value
	var body_caps_298 vm.Value
	var n_caps_299 vm.Value
	var header_300 vm.Value
	var cap_param_ids_301 vm.Value
	var init_vals_302 vm.Value
	var cap_vals_303 vm.Value
	var bt_304 vm.Value
	var entry_end_305 vm.Value
	var __306 vm.Value
	var doseq_loop__3588_307 vm.Value
	var bindings_308 vm.Value
	var loop_param_ids_309 vm.Value
	var doseq_seq__3589_310 vm.Value
	var doseq_loop__3590_311 vm.Value
	var ctx_312 vm.Value
	var v958 int
	var v961 vm.Value
	var v964 int
	var vec__3591_338 vm.Value
	var sym_344 vm.Value
	var pid_350 vm.Value
	var v352 vm.Value
	var v354 vm.Value
	var form_313 vm.Value
	var vec__3586_314 vm.Value
	var body_315 vm.Value
	var f_316 vm.Value
	var n_slots_317 vm.Value
	var old_header_318 vm.Value
	var old_caps_319 vm.Value
	var known_320 vm.Value
	var body_caps_321 vm.Value
	var n_caps_322 vm.Value
	var header_323 vm.Value
	var cap_param_ids_324 vm.Value
	var init_vals_325 vm.Value
	var cap_vals_326 vm.Value
	var bt_327 vm.Value
	var entry_end_328 vm.Value
	var __329 vm.Value
	var doseq_loop__3588_330 vm.Value
	var bindings_331 vm.Value
	var loop_param_ids_332 vm.Value
	var doseq_seq__3589_333 vm.Value
	var doseq_loop__3590_334 vm.Value
	var ctx_335 vm.Value
	var v960 int
	var v963 vm.Value
	var v966 int
	var v358 vm.Value
	var form_359 vm.Value
	var vec__3586_360 vm.Value
	var body_361 vm.Value
	var f_362 vm.Value
	var n_slots_363 vm.Value
	var old_header_364 vm.Value
	var old_caps_365 vm.Value
	var known_366 vm.Value
	var body_caps_367 vm.Value
	var n_caps_368 vm.Value
	var header_369 vm.Value
	var cap_param_ids_370 vm.Value
	var init_vals_371 vm.Value
	var cap_vals_372 vm.Value
	var bt_373 vm.Value
	var entry_end_374 vm.Value
	var __375 vm.Value
	var doseq_loop__3588_376 vm.Value
	var bindings_377 vm.Value
	var loop_param_ids_378 vm.Value
	var doseq_seq__3589_379 vm.Value
	var doseq_loop__3590_380 vm.Value
	var ctx_381 vm.Value
	var arg__4098_383 vm.Value
	var arg__4110_390 vm.Value
	var arg__4118_395 vm.Value
	var arg__4130_402 vm.Value
	var arg__4142_409 vm.Value
	var arg__4150_414 vm.Value
	var arg__4158_419 vm.Value
	var arg__4164_422 vm.Value
	var arg__4176_429 vm.Value
	var arg__4184_434 vm.Value
	var arg__4196_441 vm.Value
	var arg__4208_448 vm.Value
	var arg__4216_453 vm.Value
	var arg__4224_458 vm.Value
	var v459 vm.Value
	var pre_locals_461 vm.Value
	var fs_462 vm.Value
	var last_val_463 vm.Value
	var ctx_464 vm.Value
	var v518 vm.Value
	var form_467 vm.Value
	var vec__3586_468 vm.Value
	var body_469 vm.Value
	var f_470 vm.Value
	var n_slots_471 vm.Value
	var old_header_472 vm.Value
	var old_caps_473 vm.Value
	var known_474 vm.Value
	var body_caps_475 vm.Value
	var n_caps_476 vm.Value
	var header_477 vm.Value
	var cap_param_ids_478 vm.Value
	var init_vals_479 vm.Value
	var cap_vals_480 vm.Value
	var bt_481 vm.Value
	var entry_end_482 vm.Value
	var __483 vm.Value
	var doseq_loop__3588_484 vm.Value
	var bindings_485 vm.Value
	var loop_param_ids_486 vm.Value
	var doseq_loop__3590_487 vm.Value
	var pre_locals_488 vm.Value
	var fs_489 vm.Value
	var last_val_490 vm.Value
	var ctx_491 vm.Value
	var v521 vm.Value
	var arg__4237_523 vm.Value
	var arg__4243_526 vm.Value
	var v527 vm.Value
	var form_492 vm.Value
	var vec__3586_493 vm.Value
	var body_494 vm.Value
	var f_495 vm.Value
	var n_slots_496 vm.Value
	var old_header_497 vm.Value
	var old_caps_498 vm.Value
	var known_499 vm.Value
	var body_caps_500 vm.Value
	var n_caps_501 vm.Value
	var header_502 vm.Value
	var cap_param_ids_503 vm.Value
	var init_vals_504 vm.Value
	var cap_vals_505 vm.Value
	var bt_506 vm.Value
	var entry_end_507 vm.Value
	var __508 vm.Value
	var doseq_loop__3588_509 vm.Value
	var bindings_510 vm.Value
	var loop_param_ids_511 vm.Value
	var doseq_loop__3590_512 vm.Value
	var pre_locals_513 vm.Value
	var fs_514 vm.Value
	var last_val_515 vm.Value
	var ctx_516 vm.Value
	var result_530 vm.Value
	var form_531 vm.Value
	var vec__3586_532 vm.Value
	var body_533 vm.Value
	var f_534 vm.Value
	var n_slots_535 vm.Value
	var old_header_536 vm.Value
	var old_caps_537 vm.Value
	var known_538 vm.Value
	var body_caps_539 vm.Value
	var n_caps_540 vm.Value
	var header_541 vm.Value
	var cap_param_ids_542 vm.Value
	var init_vals_543 vm.Value
	var cap_vals_544 vm.Value
	var bt_545 vm.Value
	var entry_end_546 vm.Value
	var __547 vm.Value
	var doseq_loop__3588_548 vm.Value
	var bindings_549 vm.Value
	var loop_param_ids_550 vm.Value
	var doseq_loop__3590_551 vm.Value
	var pre_locals_552 vm.Value
	var fs_553 vm.Value
	var last_val_554 vm.Value
	var ctx_555 vm.Value
	var post_locals_557 vm.Value
	var v559 vm.Value
	var doseq_seq__3592_561 vm.Value
	var doseq_loop__3593_562 vm.Value
	var pre_locals_563 vm.Value
	var ctx_564 vm.Value
	var v976 int
	var v985 vm.Value
	var v994 int
	var result_566 vm.Value
	var form_567 vm.Value
	var vec__3586_568 vm.Value
	var body_569 vm.Value
	var f_570 vm.Value
	var n_slots_571 vm.Value
	var old_header_572 vm.Value
	var old_caps_573 vm.Value
	var known_574 vm.Value
	var body_caps_575 vm.Value
	var n_caps_576 vm.Value
	var header_577 vm.Value
	var cap_param_ids_578 vm.Value
	var init_vals_579 vm.Value
	var cap_vals_580 vm.Value
	var bt_581 vm.Value
	var entry_end_582 vm.Value
	var __583 vm.Value
	var doseq_loop__3588_584 vm.Value
	var bindings_585 vm.Value
	var loop_param_ids_586 vm.Value
	var doseq_loop__3590_587 vm.Value
	var fs_588 vm.Value
	var last_val_589 vm.Value
	var post_locals_590 vm.Value
	var doseq_seq__3592_591 vm.Value
	var doseq_loop__3593_592 vm.Value
	var pre_locals_593 vm.Value
	var ctx_594 vm.Value
	var v983 int
	var v992 vm.Value
	var v1001 int
	var vec__3594_626 vm.Value
	var sym_632 vm.Value
	var val_638 vm.Value
	var and__x_704 vm.Value
	var result_595 vm.Value
	var form_596 vm.Value
	var vec__3586_597 vm.Value
	var body_598 vm.Value
	var f_599 vm.Value
	var n_slots_600 vm.Value
	var old_header_601 vm.Value
	var old_caps_602 vm.Value
	var known_603 vm.Value
	var body_caps_604 vm.Value
	var n_caps_605 vm.Value
	var header_606 vm.Value
	var cap_param_ids_607 vm.Value
	var init_vals_608 vm.Value
	var cap_vals_609 vm.Value
	var bt_610 vm.Value
	var entry_end_611 vm.Value
	var __612 vm.Value
	var doseq_loop__3588_613 vm.Value
	var bindings_614 vm.Value
	var loop_param_ids_615 vm.Value
	var doseq_loop__3590_616 vm.Value
	var fs_617 vm.Value
	var last_val_618 vm.Value
	var post_locals_619 vm.Value
	var doseq_seq__3592_620 vm.Value
	var doseq_loop__3593_621 vm.Value
	var pre_locals_622 vm.Value
	var ctx_623 vm.Value
	var v984 int
	var v993 vm.Value
	var v1002 int
	var v858 vm.Value
	var result_859 vm.Value
	var form_860 vm.Value
	var vec__3586_861 vm.Value
	var body_862 vm.Value
	var f_863 vm.Value
	var n_slots_864 vm.Value
	var old_header_865 vm.Value
	var old_caps_866 vm.Value
	var known_867 vm.Value
	var body_caps_868 vm.Value
	var n_caps_869 vm.Value
	var header_870 vm.Value
	var cap_param_ids_871 vm.Value
	var init_vals_872 vm.Value
	var cap_vals_873 vm.Value
	var bt_874 vm.Value
	var entry_end_875 vm.Value
	var __876 vm.Value
	var doseq_loop__3588_877 vm.Value
	var bindings_878 vm.Value
	var loop_param_ids_879 vm.Value
	var doseq_loop__3590_880 vm.Value
	var fs_881 vm.Value
	var last_val_882 vm.Value
	var post_locals_883 vm.Value
	var doseq_seq__3592_884 vm.Value
	var doseq_loop__3593_885 vm.Value
	var pre_locals_886 vm.Value
	var ctx_887 vm.Value
	var arg__4305_889 vm.Value
	var arg__4315_894 vm.Value
	var arg__4321_897 vm.Value
	var arg__4331_902 vm.Value
	var arg__4341_907 vm.Value
	var arg__4347_910 vm.Value
	var arg__4353_913 vm.Value
	var arg__4359_916 vm.Value
	var arg__4369_921 vm.Value
	var arg__4375_924 vm.Value
	var arg__4385_929 vm.Value
	var arg__4395_934 vm.Value
	var arg__4401_937 vm.Value
	var arg__4407_940 vm.Value
	var v941 vm.Value
	var result_639 vm.Value
	var form_640 vm.Value
	var vec__3586_641 vm.Value
	var body_642 vm.Value
	var f_643 vm.Value
	var n_slots_644 vm.Value
	var old_header_645 vm.Value
	var old_caps_646 vm.Value
	var known_647 vm.Value
	var body_caps_648 vm.Value
	var n_caps_649 vm.Value
	var header_650 vm.Value
	var cap_param_ids_651 vm.Value
	var init_vals_652 vm.Value
	var cap_vals_653 vm.Value
	var bt_654 vm.Value
	var entry_end_655 vm.Value
	var __656 vm.Value
	var doseq_loop__3588_657 vm.Value
	var bindings_658 vm.Value
	var loop_param_ids_659 vm.Value
	var doseq_loop__3590_660 vm.Value
	var fs_661 vm.Value
	var last_val_662 vm.Value
	var post_locals_663 vm.Value
	var doseq_seq__3592_664 vm.Value
	var doseq_loop__3593_665 vm.Value
	var pre_locals_666 vm.Value
	var ctx_667 vm.Value
	var vec__3594_668 vm.Value
	var sym_669 vm.Value
	var val_670 vm.Value
	var v980 int
	var v989 vm.Value
	var v998 int
	var v816 vm.Value
	var result_671 vm.Value
	var form_672 vm.Value
	var vec__3586_673 vm.Value
	var body_674 vm.Value
	var f_675 vm.Value
	var n_slots_676 vm.Value
	var old_header_677 vm.Value
	var old_caps_678 vm.Value
	var known_679 vm.Value
	var body_caps_680 vm.Value
	var n_caps_681 vm.Value
	var header_682 vm.Value
	var cap_param_ids_683 vm.Value
	var init_vals_684 vm.Value
	var cap_vals_685 vm.Value
	var bt_686 vm.Value
	var entry_end_687 vm.Value
	var __688 vm.Value
	var doseq_loop__3588_689 vm.Value
	var bindings_690 vm.Value
	var loop_param_ids_691 vm.Value
	var doseq_loop__3590_692 vm.Value
	var fs_693 vm.Value
	var last_val_694 vm.Value
	var post_locals_695 vm.Value
	var doseq_seq__3592_696 vm.Value
	var doseq_loop__3593_697 vm.Value
	var pre_locals_698 vm.Value
	var ctx_699 vm.Value
	var vec__3594_700 vm.Value
	var sym_701 vm.Value
	var val_702 vm.Value
	var v979 int
	var v988 vm.Value
	var v997 int
	var v820 vm.Value
	var result_821 vm.Value
	var form_822 vm.Value
	var vec__3586_823 vm.Value
	var body_824 vm.Value
	var f_825 vm.Value
	var n_slots_826 vm.Value
	var old_header_827 vm.Value
	var old_caps_828 vm.Value
	var known_829 vm.Value
	var body_caps_830 vm.Value
	var n_caps_831 vm.Value
	var header_832 vm.Value
	var cap_param_ids_833 vm.Value
	var init_vals_834 vm.Value
	var cap_vals_835 vm.Value
	var bt_836 vm.Value
	var entry_end_837 vm.Value
	var __838 vm.Value
	var doseq_loop__3588_839 vm.Value
	var bindings_840 vm.Value
	var loop_param_ids_841 vm.Value
	var doseq_loop__3590_842 vm.Value
	var fs_843 vm.Value
	var last_val_844 vm.Value
	var post_locals_845 vm.Value
	var doseq_seq__3592_846 vm.Value
	var doseq_loop__3593_847 vm.Value
	var pre_locals_848 vm.Value
	var ctx_849 vm.Value
	var vec__3594_850 vm.Value
	var sym_851 vm.Value
	var val_852 vm.Value
	var v982 int
	var v991 vm.Value
	var v1000 int
	var v854 vm.Value
	var result_705 vm.Value
	var form_706 vm.Value
	var vec__3586_707 vm.Value
	var body_708 vm.Value
	var f_709 vm.Value
	var n_slots_710 vm.Value
	var old_header_711 vm.Value
	var old_caps_712 vm.Value
	var known_713 vm.Value
	var body_caps_714 vm.Value
	var n_caps_715 vm.Value
	var header_716 vm.Value
	var cap_param_ids_717 vm.Value
	var init_vals_718 vm.Value
	var cap_vals_719 vm.Value
	var bt_720 vm.Value
	var entry_end_721 vm.Value
	var __722 vm.Value
	var doseq_loop__3588_723 vm.Value
	var bindings_724 vm.Value
	var loop_param_ids_725 vm.Value
	var doseq_loop__3590_726 vm.Value
	var fs_727 vm.Value
	var last_val_728 vm.Value
	var post_locals_729 vm.Value
	var doseq_seq__3592_730 vm.Value
	var doseq_loop__3593_731 vm.Value
	var pre_locals_732 vm.Value
	var ctx_733 vm.Value
	var vec__3594_734 vm.Value
	var sym_735 vm.Value
	var val_736 vm.Value
	var and__x_737 vm.Value
	var v977 int
	var v986 vm.Value
	var v995 int
	var arg__4282_773 vm.Value
	var arg__4290_776 vm.Value
	var v777 vm.Value
	var result_738 vm.Value
	var form_739 vm.Value
	var vec__3586_740 vm.Value
	var body_741 vm.Value
	var f_742 vm.Value
	var n_slots_743 vm.Value
	var old_header_744 vm.Value
	var old_caps_745 vm.Value
	var known_746 vm.Value
	var body_caps_747 vm.Value
	var n_caps_748 vm.Value
	var header_749 vm.Value
	var cap_param_ids_750 vm.Value
	var init_vals_751 vm.Value
	var cap_vals_752 vm.Value
	var bt_753 vm.Value
	var entry_end_754 vm.Value
	var __755 vm.Value
	var doseq_loop__3588_756 vm.Value
	var bindings_757 vm.Value
	var loop_param_ids_758 vm.Value
	var doseq_loop__3590_759 vm.Value
	var fs_760 vm.Value
	var last_val_761 vm.Value
	var post_locals_762 vm.Value
	var doseq_seq__3592_763 vm.Value
	var doseq_loop__3593_764 vm.Value
	var pre_locals_765 vm.Value
	var ctx_766 vm.Value
	var vec__3594_767 vm.Value
	var sym_768 vm.Value
	var val_769 vm.Value
	var and__x_770 vm.Value
	var v978 int
	var v987 vm.Value
	var v996 int
	var v780 vm.Value
	var result_781 vm.Value
	var form_782 vm.Value
	var vec__3586_783 vm.Value
	var body_784 vm.Value
	var f_785 vm.Value
	var n_slots_786 vm.Value
	var old_header_787 vm.Value
	var old_caps_788 vm.Value
	var known_789 vm.Value
	var body_caps_790 vm.Value
	var n_caps_791 vm.Value
	var header_792 vm.Value
	var cap_param_ids_793 vm.Value
	var init_vals_794 vm.Value
	var cap_vals_795 vm.Value
	var bt_796 vm.Value
	var entry_end_797 vm.Value
	var __798 vm.Value
	var doseq_loop__3588_799 vm.Value
	var bindings_800 vm.Value
	var loop_param_ids_801 vm.Value
	var doseq_loop__3590_802 vm.Value
	var fs_803 vm.Value
	var last_val_804 vm.Value
	var post_locals_805 vm.Value
	var doseq_seq__3592_806 vm.Value
	var doseq_loop__3593_807 vm.Value
	var pre_locals_808 vm.Value
	var ctx_809 vm.Value
	var vec__3594_810 vm.Value
	var sym_811 vm.Value
	var val_812 vm.Value
	var and__x_813 vm.Value
	var v981 int
	var v990 vm.Value
	var v999 int
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = __8, bindings_14, body_18, f_20, arg__3620_22, arg__3626_26, n_slots_28, arg__3631_30, arg__3637_34, old_header_36, arg__3642_38, arg__3648_42, old_caps_44, arg__3653_46, arg__3658_49, arg__3659_50, arg__3664_53, arg__3669_56, arg__3670_57, known_58, arg__3688_63, arg__3708_69, arg__3710_70, arg__3729_76, arg__3749_82, arg__3751_83, body_caps_84, n_caps_86, header_88, arg__3780_95, arg__3804_103, loop_param_ids_104, arg__3831_112, arg__3859_121, cap_param_ids_122, arg__3891_128, arg__3924_135, init_vals_136, cap_vals_144, arg__3946_146, arg__3953_149, arg__3954_150, arg__3962_153, arg__3969_156, arg__3970_157, bt_158, entry_end_160, __166, __168, v170, v172, arg__4003_174, arg__4008_177, doseq_seq__3587_178, doseq_loop__3588_179, bindings_180, ctx_181, loop_param_ids_182, v948, form_184, vec__3586_185, body_186, f_187, n_slots_188, old_header_189, old_caps_190, known_191, body_caps_192, n_caps_193, header_194, cap_param_ids_195, init_vals_196, cap_vals_197, bt_198, entry_end_199, __200, doseq_seq__3587_201, doseq_loop__3588_202, bindings_203, ctx_204, loop_param_ids_205, v947, i_230, arg__4016_232, arg__4022_236, arg__4028_238, arg__4040_245, arg__4046_247, v248, v250, form_206, vec__3586_207, body_208, f_209, n_slots_210, old_header_211, old_caps_212, known_213, body_caps_214, n_caps_215, header_216, cap_param_ids_217, init_vals_218, cap_vals_219, bt_220, entry_end_221, __222, doseq_seq__3587_223, doseq_loop__3588_224, bindings_225, ctx_226, loop_param_ids_227, v949, v254, form_255, vec__3586_256, body_257, f_258, n_slots_259, old_header_260, old_caps_261, known_262, body_caps_263, n_caps_264, header_265, cap_param_ids_266, init_vals_267, cap_vals_268, bt_269, entry_end_270, __271, doseq_seq__3587_272, doseq_loop__3588_273, bindings_274, ctx_275, loop_param_ids_276, arg__4057_280, arg__4066_285, doseq_seq__3589_286, doseq_loop__3590_287, ctx_288, v959, v962, v965, form_290, vec__3586_291, body_292, f_293, n_slots_294, old_header_295, old_caps_296, known_297, body_caps_298, n_caps_299, header_300, cap_param_ids_301, init_vals_302, cap_vals_303, bt_304, entry_end_305, __306, doseq_loop__3588_307, bindings_308, loop_param_ids_309, doseq_seq__3589_310, doseq_loop__3590_311, ctx_312, v958, v961, v964, vec__3591_338, sym_344, pid_350, v352, v354, form_313, vec__3586_314, body_315, f_316, n_slots_317, old_header_318, old_caps_319, known_320, body_caps_321, n_caps_322, header_323, cap_param_ids_324, init_vals_325, cap_vals_326, bt_327, entry_end_328, __329, doseq_loop__3588_330, bindings_331, loop_param_ids_332, doseq_seq__3589_333, doseq_loop__3590_334, ctx_335, v960, v963, v966, v358, form_359, vec__3586_360, body_361, f_362, n_slots_363, old_header_364, old_caps_365, known_366, body_caps_367, n_caps_368, header_369, cap_param_ids_370, init_vals_371, cap_vals_372, bt_373, entry_end_374, __375, doseq_loop__3588_376, bindings_377, loop_param_ids_378, doseq_seq__3589_379, doseq_loop__3590_380, ctx_381, arg__4098_383, arg__4110_390, arg__4118_395, arg__4130_402, arg__4142_409, arg__4150_414, arg__4158_419, arg__4164_422, arg__4176_429, arg__4184_434, arg__4196_441, arg__4208_448, arg__4216_453, arg__4224_458, v459, pre_locals_461, fs_462, last_val_463, ctx_464, v518, form_467, vec__3586_468, body_469, f_470, n_slots_471, old_header_472, old_caps_473, known_474, body_caps_475, n_caps_476, header_477, cap_param_ids_478, init_vals_479, cap_vals_480, bt_481, entry_end_482, __483, doseq_loop__3588_484, bindings_485, loop_param_ids_486, doseq_loop__3590_487, pre_locals_488, fs_489, last_val_490, ctx_491, v521, arg__4237_523, arg__4243_526, v527, form_492, vec__3586_493, body_494, f_495, n_slots_496, old_header_497, old_caps_498, known_499, body_caps_500, n_caps_501, header_502, cap_param_ids_503, init_vals_504, cap_vals_505, bt_506, entry_end_507, __508, doseq_loop__3588_509, bindings_510, loop_param_ids_511, doseq_loop__3590_512, pre_locals_513, fs_514, last_val_515, ctx_516, result_530, form_531, vec__3586_532, body_533, f_534, n_slots_535, old_header_536, old_caps_537, known_538, body_caps_539, n_caps_540, header_541, cap_param_ids_542, init_vals_543, cap_vals_544, bt_545, entry_end_546, __547, doseq_loop__3588_548, bindings_549, loop_param_ids_550, doseq_loop__3590_551, pre_locals_552, fs_553, last_val_554, ctx_555, post_locals_557, v559, doseq_seq__3592_561, doseq_loop__3593_562, pre_locals_563, ctx_564, v976, v985, v994, result_566, form_567, vec__3586_568, body_569, f_570, n_slots_571, old_header_572, old_caps_573, known_574, body_caps_575, n_caps_576, header_577, cap_param_ids_578, init_vals_579, cap_vals_580, bt_581, entry_end_582, __583, doseq_loop__3588_584, bindings_585, loop_param_ids_586, doseq_loop__3590_587, fs_588, last_val_589, post_locals_590, doseq_seq__3592_591, doseq_loop__3593_592, pre_locals_593, ctx_594, v983, v992, v1001, vec__3594_626, sym_632, val_638, and__x_704, result_595, form_596, vec__3586_597, body_598, f_599, n_slots_600, old_header_601, old_caps_602, known_603, body_caps_604, n_caps_605, header_606, cap_param_ids_607, init_vals_608, cap_vals_609, bt_610, entry_end_611, __612, doseq_loop__3588_613, bindings_614, loop_param_ids_615, doseq_loop__3590_616, fs_617, last_val_618, post_locals_619, doseq_seq__3592_620, doseq_loop__3593_621, pre_locals_622, ctx_623, v984, v993, v1002, v858, result_859, form_860, vec__3586_861, body_862, f_863, n_slots_864, old_header_865, old_caps_866, known_867, body_caps_868, n_caps_869, header_870, cap_param_ids_871, init_vals_872, cap_vals_873, bt_874, entry_end_875, __876, doseq_loop__3588_877, bindings_878, loop_param_ids_879, doseq_loop__3590_880, fs_881, last_val_882, post_locals_883, doseq_seq__3592_884, doseq_loop__3593_885, pre_locals_886, ctx_887, arg__4305_889, arg__4315_894, arg__4321_897, arg__4331_902, arg__4341_907, arg__4347_910, arg__4353_913, arg__4359_916, arg__4369_921, arg__4375_924, arg__4385_929, arg__4395_934, arg__4401_937, arg__4407_940, v941, result_639, form_640, vec__3586_641, body_642, f_643, n_slots_644, old_header_645, old_caps_646, known_647, body_caps_648, n_caps_649, header_650, cap_param_ids_651, init_vals_652, cap_vals_653, bt_654, entry_end_655, __656, doseq_loop__3588_657, bindings_658, loop_param_ids_659, doseq_loop__3590_660, fs_661, last_val_662, post_locals_663, doseq_seq__3592_664, doseq_loop__3593_665, pre_locals_666, ctx_667, vec__3594_668, sym_669, val_670, v980, v989, v998, v816, result_671, form_672, vec__3586_673, body_674, f_675, n_slots_676, old_header_677, old_caps_678, known_679, body_caps_680, n_caps_681, header_682, cap_param_ids_683, init_vals_684, cap_vals_685, bt_686, entry_end_687, __688, doseq_loop__3588_689, bindings_690, loop_param_ids_691, doseq_loop__3590_692, fs_693, last_val_694, post_locals_695, doseq_seq__3592_696, doseq_loop__3593_697, pre_locals_698, ctx_699, vec__3594_700, sym_701, val_702, v979, v988, v997, v820, result_821, form_822, vec__3586_823, body_824, f_825, n_slots_826, old_header_827, old_caps_828, known_829, body_caps_830, n_caps_831, header_832, cap_param_ids_833, init_vals_834, cap_vals_835, bt_836, entry_end_837, __838, doseq_loop__3588_839, bindings_840, loop_param_ids_841, doseq_loop__3590_842, fs_843, last_val_844, post_locals_845, doseq_seq__3592_846, doseq_loop__3593_847, pre_locals_848, ctx_849, vec__3594_850, sym_851, val_852, v982, v991, v1000, v854, result_705, form_706, vec__3586_707, body_708, f_709, n_slots_710, old_header_711, old_caps_712, known_713, body_caps_714, n_caps_715, header_716, cap_param_ids_717, init_vals_718, cap_vals_719, bt_720, entry_end_721, __722, doseq_loop__3588_723, bindings_724, loop_param_ids_725, doseq_loop__3590_726, fs_727, last_val_728, post_locals_729, doseq_seq__3592_730, doseq_loop__3593_731, pre_locals_732, ctx_733, vec__3594_734, sym_735, val_736, and__x_737, v977, v986, v995, arg__4282_773, arg__4290_776, v777, result_738, form_739, vec__3586_740, body_741, f_742, n_slots_743, old_header_744, old_caps_745, known_746, body_caps_747, n_caps_748, header_749, cap_param_ids_750, init_vals_751, cap_vals_752, bt_753, entry_end_754, __755, doseq_loop__3588_756, bindings_757, loop_param_ids_758, doseq_loop__3590_759, fs_760, last_val_761, post_locals_762, doseq_seq__3592_763, doseq_loop__3593_764, pre_locals_765, ctx_766, vec__3594_767, sym_768, val_769, and__x_770, v978, v987, v996, v780, result_781, form_782, vec__3586_783, body_784, f_785, n_slots_786, old_header_787, old_caps_788, known_789, body_caps_790, n_caps_791, header_792, cap_param_ids_793, init_vals_794, cap_vals_795, bt_796, entry_end_797, __798, doseq_loop__3588_799, bindings_800, loop_param_ids_801, doseq_loop__3590_802, fs_803, last_val_804, post_locals_805, doseq_seq__3592_806, doseq_loop__3593_807, pre_locals_808, ctx_809, vec__3594_810, sym_811, val_812, and__x_813, v981, v990, v999
	__8, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(0), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	bindings_14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(1), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	body_18, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "drop").Deref(), []vm.Value{vm.Int(2), arg0})
	if callErr != nil {
		return nil, callErr
	}
	f_20, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-fn").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__3620_22, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{bindings_14})
	if callErr != nil {
		return nil, callErr
	}
	arg__3626_26, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{bindings_14})
	if callErr != nil {
		return nil, callErr
	}
	n_slots_28, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "/").Deref(), []vm.Value{arg__3626_26, vm.Int(2)})
	if callErr != nil {
		return nil, callErr
	}
	arg__3631_30, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__3637_34, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	old_header_36, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg__3637_34, vm.Keyword("loop-header")})
	if callErr != nil {
		return nil, callErr
	}
	arg__3642_38, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__3648_42, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	old_caps_44, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg__3648_42, vm.Keyword("loop-capture-syms")})
	if callErr != nil {
		return nil, callErr
	}
	arg__3653_46, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "current-locals-flat").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__3658_49, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "current-locals-flat").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__3659_50, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "keys").Deref(), []vm.Value{arg__3658_49})
	if callErr != nil {
		return nil, callErr
	}
	arg__3664_53, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "current-locals-flat").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__3669_56, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "current-locals-flat").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__3670_57, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "keys").Deref(), []vm.Value{arg__3669_56})
	if callErr != nil {
		return nil, callErr
	}
	known_58, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "set").Deref(), []vm.Value{arg__3670_57})
	if callErr != nil {
		return nil, callErr
	}
	arg__3688_63, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__3708_69, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__3710_70, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "reduce").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
		var arg__3697_4 vm.Value
		var arg__3705_7 vm.Value
		var v8 vm.Value
		var callErr error
		_, _, _ = arg__3697_4, arg__3705_7, v8
		arg__3697_4, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "captures-of").Deref(), []vm.Value{arg1, known_58})
		if callErr != nil {
			return nil, callErr
		}
		arg__3705_7, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "captures-of").Deref(), []vm.Value{arg1, known_58})
		if callErr != nil {
			return nil, callErr
		}
		v8, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg0, arg__3705_7})
		if callErr != nil {
			return nil, callErr
		}
		return v8, nil
	}), arg__3708_69, body_18})
	if callErr != nil {
		return nil, callErr
	}
	arg__3729_76, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__3749_82, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__3751_83, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "reduce").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
		var arg__3738_4 vm.Value
		var arg__3746_7 vm.Value
		var v8 vm.Value
		var callErr error
		_, _, _ = arg__3738_4, arg__3746_7, v8
		arg__3738_4, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "captures-of").Deref(), []vm.Value{arg1, known_58})
		if callErr != nil {
			return nil, callErr
		}
		arg__3746_7, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "captures-of").Deref(), []vm.Value{arg1, known_58})
		if callErr != nil {
			return nil, callErr
		}
		v8, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg0, arg__3746_7})
		if callErr != nil {
			return nil, callErr
		}
		return v8, nil
	}), arg__3749_82, body_18})
	if callErr != nil {
		return nil, callErr
	}
	body_caps_84, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg__3751_83})
	if callErr != nil {
		return nil, callErr
	}
	n_caps_86, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{body_caps_84})
	if callErr != nil {
		return nil, callErr
	}
	header_88, callErr = rt.InvokeValue(rt.LookupVar("ir", "add-block").Deref(), []vm.Value{f_20})
	if callErr != nil {
		return nil, callErr
	}
	arg__3780_95, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "range").Deref(), []vm.Value{n_slots_28})
	if callErr != nil {
		return nil, callErr
	}
	arg__3804_103, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "range").Deref(), []vm.Value{n_slots_28})
	if callErr != nil {
		return nil, callErr
	}
	loop_param_ids_104, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapv").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var pid_9 vm.Value
		var v11 vm.Value
		var callErr error
		_, _ = pid_9, v11
		pid_9, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "add-inst!").Deref(), []vm.Value{arg1, header_88, vm.Keyword("block-arg"), vm.NewArrayVector([]vm.Value{}), arg0})
		if callErr != nil {
			return nil, callErr
		}
		v11, callErr = rt.InvokeValue(rt.LookupVar("ir", "add-block-param!").Deref(), []vm.Value{f_20, header_88, pid_9})
		if callErr != nil {
			return nil, callErr
		}
		return pid_9, nil
	}), arg__3804_103})
	if callErr != nil {
		return nil, callErr
	}
	arg__3831_112, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "range").Deref(), []vm.Value{n_caps_86})
	if callErr != nil {
		return nil, callErr
	}
	arg__3859_121, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "range").Deref(), []vm.Value{n_caps_86})
	if callErr != nil {
		return nil, callErr
	}
	cap_param_ids_122, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapv").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var arg__3839_7 vm.Value
		var arg__3847_11 vm.Value
		var pid_12 vm.Value
		var v14 vm.Value
		var callErr error
		_, _, _, _ = arg__3839_7, arg__3847_11, pid_12, v14
		arg__3839_7 = rt.AddValue(n_slots_28, arg0)
		arg__3847_11 = rt.AddValue(n_slots_28, arg0)
		pid_12, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "add-inst!").Deref(), []vm.Value{arg1, header_88, vm.Keyword("block-arg"), vm.NewArrayVector([]vm.Value{}), arg__3847_11})
		if callErr != nil {
			return nil, callErr
		}
		v14, callErr = rt.InvokeValue(rt.LookupVar("ir", "add-block-param!").Deref(), []vm.Value{f_20, header_88, pid_12})
		if callErr != nil {
			return nil, callErr
		}
		return pid_12, nil
	}), arg__3859_121})
	if callErr != nil {
		return nil, callErr
	}
	arg__3891_128, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "range").Deref(), []vm.Value{n_slots_28})
	if callErr != nil {
		return nil, callErr
	}
	arg__3924_135, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "range").Deref(), []vm.Value{n_slots_28})
	if callErr != nil {
		return nil, callErr
	}
	init_vals_136, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapv").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var arg__3896_4 vm.Value
		var arg__3897_5 vm.Value
		var arg__3902_8 vm.Value
		var arg__3903_9 vm.Value
		var arg__3904_10 vm.Value
		var arg__3910_13 vm.Value
		var arg__3911_14 vm.Value
		var arg__3916_17 vm.Value
		var arg__3917_18 vm.Value
		var arg__3918_19 vm.Value
		var v20 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _, _, _, _ = arg__3896_4, arg__3897_5, arg__3902_8, arg__3903_9, arg__3904_10, arg__3910_13, arg__3911_14, arg__3916_17, arg__3917_18, arg__3918_19, v20
		arg__3896_4 = rt.MulValue(arg0, vm.Int(2))
		arg__3897_5 = rt.AddValue(arg__3896_4, vm.Int(1))
		arg__3902_8 = rt.MulValue(arg0, vm.Int(2))
		arg__3903_9 = rt.AddValue(arg__3902_8, vm.Int(1))
		arg__3904_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{bindings_14, arg__3903_9})
		if callErr != nil {
			return nil, callErr
		}
		arg__3910_13 = rt.MulValue(arg0, vm.Int(2))
		arg__3911_14 = rt.AddValue(arg__3910_13, vm.Int(1))
		arg__3916_17 = rt.MulValue(arg0, vm.Int(2))
		arg__3917_18 = rt.AddValue(arg__3916_17, vm.Int(1))
		arg__3918_19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{bindings_14, arg__3917_18})
		if callErr != nil {
			return nil, callErr
		}
		v20, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-form").Deref(), []vm.Value{arg__3918_19, arg1})
		if callErr != nil {
			return nil, callErr
		}
		return v20, nil
	}), arg__3924_135})
	if callErr != nil {
		return nil, callErr
	}
	cap_vals_144, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapv").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "lookup-local").Deref(), []vm.Value{arg1, arg0})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), body_caps_84})
	if callErr != nil {
		return nil, callErr
	}
	arg__3946_146, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "concat").Deref(), []vm.Value{init_vals_136, cap_vals_144})
	if callErr != nil {
		return nil, callErr
	}
	arg__3953_149, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "concat").Deref(), []vm.Value{init_vals_136, cap_vals_144})
	if callErr != nil {
		return nil, callErr
	}
	arg__3954_150, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg__3953_149})
	if callErr != nil {
		return nil, callErr
	}
	arg__3962_153, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "concat").Deref(), []vm.Value{init_vals_136, cap_vals_144})
	if callErr != nil {
		return nil, callErr
	}
	arg__3969_156, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "concat").Deref(), []vm.Value{init_vals_136, cap_vals_144})
	if callErr != nil {
		return nil, callErr
	}
	arg__3970_157, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg__3969_156})
	if callErr != nil {
		return nil, callErr
	}
	bt_158, callErr = rt.InvokeValue(rt.LookupVar("ir", "new-branch-target").Deref(), []vm.Value{header_88, arg__3970_157})
	if callErr != nil {
		return nil, callErr
	}
	entry_end_160, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	__166, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "add-terminator!").Deref(), []vm.Value{arg1, entry_end_160, vm.Keyword("branch"), vm.NewArrayVector([]vm.Value{}), bt_158})
	if callErr != nil {
		return nil, callErr
	}
	__168, callErr = rt.InvokeValue(rt.LookupVar("ir", "add-pred!").Deref(), []vm.Value{f_20, header_88, entry_end_160})
	if callErr != nil {
		return nil, callErr
	}
	v170, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-set-block!").Deref(), []vm.Value{arg1, header_88})
	if callErr != nil {
		return nil, callErr
	}
	v172, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "push-locals!").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__4003_174, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "range").Deref(), []vm.Value{n_slots_28})
	if callErr != nil {
		return nil, callErr
	}
	arg__4008_177, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "range").Deref(), []vm.Value{n_slots_28})
	if callErr != nil {
		return nil, callErr
	}
	doseq_seq__3587_178, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{arg__4008_177})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__3588_179 = doseq_seq__3587_178
	bindings_180 = bindings_14
	ctx_181 = arg1
	loop_param_ids_182 = loop_param_ids_104
	v948 = 2
	goto b1
b1:
	;
	if vm.IsTruthy(doseq_loop__3588_179) {
		form_184 = arg0
		vec__3586_185 = arg0
		body_186 = body_18
		f_187 = f_20
		n_slots_188 = n_slots_28
		old_header_189 = old_header_36
		old_caps_190 = old_caps_44
		known_191 = known_58
		body_caps_192 = body_caps_84
		n_caps_193 = n_caps_86
		header_194 = header_88
		cap_param_ids_195 = cap_param_ids_122
		init_vals_196 = init_vals_136
		cap_vals_197 = cap_vals_144
		bt_198 = bt_158
		entry_end_199 = entry_end_160
		__200 = __168
		doseq_seq__3587_201 = doseq_seq__3587_178
		doseq_loop__3588_202 = doseq_loop__3588_179
		bindings_203 = bindings_180
		ctx_204 = ctx_181
		loop_param_ids_205 = loop_param_ids_182
		v947 = v948
		goto b2
	} else {
		form_206 = arg0
		vec__3586_207 = arg0
		body_208 = body_18
		f_209 = f_20
		n_slots_210 = n_slots_28
		old_header_211 = old_header_36
		old_caps_212 = old_caps_44
		known_213 = known_58
		body_caps_214 = body_caps_84
		n_caps_215 = n_caps_86
		header_216 = header_88
		cap_param_ids_217 = cap_param_ids_122
		init_vals_218 = init_vals_136
		cap_vals_219 = cap_vals_144
		bt_220 = bt_158
		entry_end_221 = entry_end_160
		__222 = __168
		doseq_seq__3587_223 = doseq_seq__3587_178
		doseq_loop__3588_224 = doseq_loop__3588_179
		bindings_225 = bindings_180
		ctx_226 = ctx_181
		loop_param_ids_227 = loop_param_ids_182
		v949 = v948
		goto b3
	}
b2:
	;
	i_230, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{doseq_loop__3588_202})
	if callErr != nil {
		return nil, callErr
	}
	arg__4016_232 = rt.MulValue(i_230, vm.Int(v947))
	arg__4022_236, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{bindings_203, arg__4016_232})
	if callErr != nil {
		return nil, callErr
	}
	arg__4028_238, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{loop_param_ids_205, i_230})
	if callErr != nil {
		return nil, callErr
	}
	arg__4040_245, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{bindings_203, arg__4016_232})
	if callErr != nil {
		return nil, callErr
	}
	arg__4046_247, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{loop_param_ids_205, i_230})
	if callErr != nil {
		return nil, callErr
	}
	v248, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "bind-local!").Deref(), []vm.Value{ctx_204, arg__4040_245, arg__4046_247})
	if callErr != nil {
		return nil, callErr
	}
	v250, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "next").Deref(), []vm.Value{doseq_loop__3588_202})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__3588_179 = v250
	bindings_180 = bindings_203
	ctx_181 = ctx_204
	loop_param_ids_182 = loop_param_ids_205
	v948 = v947
	goto b1
b3:
	;
	v254 = vm.NIL
	form_255 = form_206
	vec__3586_256 = vec__3586_207
	body_257 = body_208
	f_258 = f_209
	n_slots_259 = n_slots_210
	old_header_260 = old_header_211
	old_caps_261 = old_caps_212
	known_262 = known_213
	body_caps_263 = body_caps_214
	n_caps_264 = n_caps_215
	header_265 = header_216
	cap_param_ids_266 = cap_param_ids_217
	init_vals_267 = init_vals_218
	cap_vals_268 = cap_vals_219
	bt_269 = bt_220
	entry_end_270 = entry_end_221
	__271 = __222
	doseq_seq__3587_272 = doseq_seq__3587_223
	doseq_loop__3588_273 = doseq_loop__3588_224
	bindings_274 = bindings_225
	ctx_275 = ctx_226
	loop_param_ids_276 = loop_param_ids_227
	goto b4
b4:
	;
	arg__4057_280, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.LookupVar("clojure.core", "vector").Deref(), body_caps_263, cap_param_ids_266})
	if callErr != nil {
		return nil, callErr
	}
	arg__4066_285, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.LookupVar("clojure.core", "vector").Deref(), body_caps_263, cap_param_ids_266})
	if callErr != nil {
		return nil, callErr
	}
	doseq_seq__3589_286, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{arg__4066_285})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__3590_287 = doseq_seq__3589_286
	ctx_288 = ctx_275
	v959 = 0
	v962 = vm.NIL
	v965 = 1
	goto b5
b5:
	;
	if vm.IsTruthy(doseq_loop__3590_287) {
		form_290 = form_255
		vec__3586_291 = vec__3586_256
		body_292 = body_257
		f_293 = f_258
		n_slots_294 = n_slots_259
		old_header_295 = old_header_260
		old_caps_296 = old_caps_261
		known_297 = known_262
		body_caps_298 = body_caps_263
		n_caps_299 = n_caps_264
		header_300 = header_265
		cap_param_ids_301 = cap_param_ids_266
		init_vals_302 = init_vals_267
		cap_vals_303 = cap_vals_268
		bt_304 = bt_269
		entry_end_305 = entry_end_270
		__306 = __271
		doseq_loop__3588_307 = doseq_loop__3588_273
		bindings_308 = bindings_274
		loop_param_ids_309 = loop_param_ids_276
		doseq_seq__3589_310 = doseq_seq__3589_286
		doseq_loop__3590_311 = doseq_loop__3590_287
		ctx_312 = ctx_288
		v958 = v959
		v961 = v962
		v964 = v965
		goto b6
	} else {
		form_313 = form_255
		vec__3586_314 = vec__3586_256
		body_315 = body_257
		f_316 = f_258
		n_slots_317 = n_slots_259
		old_header_318 = old_header_260
		old_caps_319 = old_caps_261
		known_320 = known_262
		body_caps_321 = body_caps_263
		n_caps_322 = n_caps_264
		header_323 = header_265
		cap_param_ids_324 = cap_param_ids_266
		init_vals_325 = init_vals_267
		cap_vals_326 = cap_vals_268
		bt_327 = bt_269
		entry_end_328 = entry_end_270
		__329 = __271
		doseq_loop__3588_330 = doseq_loop__3588_273
		bindings_331 = bindings_274
		loop_param_ids_332 = loop_param_ids_276
		doseq_seq__3589_333 = doseq_seq__3589_286
		doseq_loop__3590_334 = doseq_loop__3590_287
		ctx_335 = ctx_288
		v960 = v959
		v963 = v962
		v966 = v965
		goto b7
	}
b6:
	;
	vec__3591_338, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{doseq_loop__3590_311})
	if callErr != nil {
		return nil, callErr
	}
	sym_344, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{vec__3591_338, vm.Int(v958), v961})
	if callErr != nil {
		return nil, callErr
	}
	pid_350, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{vec__3591_338, vm.Int(v964), v961})
	if callErr != nil {
		return nil, callErr
	}
	v352, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "rebind-local!").Deref(), []vm.Value{ctx_312, sym_344, pid_350})
	if callErr != nil {
		return nil, callErr
	}
	v354, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "next").Deref(), []vm.Value{doseq_loop__3590_311})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__3590_287 = v354
	ctx_288 = ctx_312
	v959 = v958
	v962 = v961
	v965 = v964
	goto b5
b7:
	;
	v358 = vm.NIL
	form_359 = form_313
	vec__3586_360 = vec__3586_314
	body_361 = body_315
	f_362 = f_316
	n_slots_363 = n_slots_317
	old_header_364 = old_header_318
	old_caps_365 = old_caps_319
	known_366 = known_320
	body_caps_367 = body_caps_321
	n_caps_368 = n_caps_322
	header_369 = header_323
	cap_param_ids_370 = cap_param_ids_324
	init_vals_371 = init_vals_325
	cap_vals_372 = cap_vals_326
	bt_373 = bt_327
	entry_end_374 = entry_end_328
	__375 = __329
	doseq_loop__3588_376 = doseq_loop__3588_330
	bindings_377 = bindings_331
	loop_param_ids_378 = loop_param_ids_332
	doseq_seq__3589_379 = doseq_seq__3589_333
	doseq_loop__3590_380 = doseq_loop__3590_334
	ctx_381 = ctx_335
	goto b8
b8:
	;
	arg__4098_383, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_381})
	if callErr != nil {
		return nil, callErr
	}
	arg__4110_390, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_381})
	if callErr != nil {
		return nil, callErr
	}
	arg__4118_395, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "update").Deref(), []vm.Value{arg__4110_390, vm.Keyword("loop-headers"), rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var or__x_2 vm.Value
		var s_3 vm.Value
		var header_4 vm.Value
		var or__x_5 vm.Value
		var s_6 vm.Value
		var header_7 vm.Value
		var arg__4112_12 vm.Value
		var or__x_13 vm.Value
		var s_14 vm.Value
		var header_15 vm.Value
		var or__x_17 vm.Value
		var s_18 vm.Value
		var header_19 vm.Value
		var head__4114_20 vm.Value
		var or__x_21 vm.Value
		var s_22 vm.Value
		var header_23 vm.Value
		var head__4114_24 vm.Value
		var arg__4115_29 vm.Value
		var or__x_30 vm.Value
		var s_31 vm.Value
		var header_32 vm.Value
		var head__4114_33 vm.Value
		var v34 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = or__x_2, s_3, header_4, or__x_5, s_6, header_7, arg__4112_12, or__x_13, s_14, header_15, or__x_17, s_18, header_19, head__4114_20, or__x_21, s_22, header_23, head__4114_24, arg__4115_29, or__x_30, s_31, header_32, head__4114_33, v34
		if vm.IsTruthy(arg0) {
			or__x_2 = arg0
			s_3 = arg0
			header_4 = header_369
			goto b1
		} else {
			or__x_5 = arg0
			s_6 = arg0
			header_7 = header_369
			goto b2
		}
	b1:
		;
		arg__4112_12 = or__x_2
		or__x_13 = or__x_2
		s_14 = s_3
		header_15 = header_4
		goto b3
	b2:
		;
		arg__4112_12 = vm.NewArrayVector([]vm.Value{})
		or__x_13 = or__x_5
		s_14 = s_6
		header_15 = header_7
		goto b3
	b3:
		;
		if vm.IsTruthy(s_14) {
			or__x_17 = s_14
			s_18 = s_14
			header_19 = header_15
			head__4114_20 = rt.LookupVar("clojure.core", "conj").Deref()
			goto b4
		} else {
			or__x_21 = s_14
			s_22 = s_14
			header_23 = header_15
			head__4114_24 = rt.LookupVar("clojure.core", "conj").Deref()
			goto b5
		}
	b4:
		;
		arg__4115_29 = or__x_17
		or__x_30 = or__x_17
		s_31 = s_18
		header_32 = header_19
		head__4114_33 = head__4114_20
		goto b6
	b5:
		;
		arg__4115_29 = vm.NewArrayVector([]vm.Value{})
		or__x_30 = or__x_21
		s_31 = s_22
		header_32 = header_23
		head__4114_33 = head__4114_24
		goto b6
	b6:
		;
		v34, callErr = rt.InvokeValue(head__4114_33, []vm.Value{arg__4115_29, header_32})
		if callErr != nil {
			return nil, callErr
		}
		return v34, nil
	})})
	if callErr != nil {
		return nil, callErr
	}
	arg__4130_402, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_381})
	if callErr != nil {
		return nil, callErr
	}
	arg__4142_409, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_381})
	if callErr != nil {
		return nil, callErr
	}
	arg__4150_414, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "update").Deref(), []vm.Value{arg__4142_409, vm.Keyword("loop-headers"), rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var or__x_2 vm.Value
		var s_3 vm.Value
		var header_4 vm.Value
		var or__x_5 vm.Value
		var s_6 vm.Value
		var header_7 vm.Value
		var arg__4144_12 vm.Value
		var or__x_13 vm.Value
		var s_14 vm.Value
		var header_15 vm.Value
		var or__x_17 vm.Value
		var s_18 vm.Value
		var header_19 vm.Value
		var head__4146_20 vm.Value
		var or__x_21 vm.Value
		var s_22 vm.Value
		var header_23 vm.Value
		var head__4146_24 vm.Value
		var arg__4147_29 vm.Value
		var or__x_30 vm.Value
		var s_31 vm.Value
		var header_32 vm.Value
		var head__4146_33 vm.Value
		var v34 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = or__x_2, s_3, header_4, or__x_5, s_6, header_7, arg__4144_12, or__x_13, s_14, header_15, or__x_17, s_18, header_19, head__4146_20, or__x_21, s_22, header_23, head__4146_24, arg__4147_29, or__x_30, s_31, header_32, head__4146_33, v34
		if vm.IsTruthy(arg0) {
			or__x_2 = arg0
			s_3 = arg0
			header_4 = header_369
			goto b1
		} else {
			or__x_5 = arg0
			s_6 = arg0
			header_7 = header_369
			goto b2
		}
	b1:
		;
		arg__4144_12 = or__x_2
		or__x_13 = or__x_2
		s_14 = s_3
		header_15 = header_4
		goto b3
	b2:
		;
		arg__4144_12 = vm.NewArrayVector([]vm.Value{})
		or__x_13 = or__x_5
		s_14 = s_6
		header_15 = header_7
		goto b3
	b3:
		;
		if vm.IsTruthy(s_14) {
			or__x_17 = s_14
			s_18 = s_14
			header_19 = header_15
			head__4146_20 = rt.LookupVar("clojure.core", "conj").Deref()
			goto b4
		} else {
			or__x_21 = s_14
			s_22 = s_14
			header_23 = header_15
			head__4146_24 = rt.LookupVar("clojure.core", "conj").Deref()
			goto b5
		}
	b4:
		;
		arg__4147_29 = or__x_17
		or__x_30 = or__x_17
		s_31 = s_18
		header_32 = header_19
		head__4146_33 = head__4146_20
		goto b6
	b5:
		;
		arg__4147_29 = vm.NewArrayVector([]vm.Value{})
		or__x_30 = or__x_21
		s_31 = s_22
		header_32 = header_23
		head__4146_33 = head__4146_24
		goto b6
	b6:
		;
		v34, callErr = rt.InvokeValue(head__4146_33, []vm.Value{arg__4147_29, header_32})
		if callErr != nil {
			return nil, callErr
		}
		return v34, nil
	})})
	if callErr != nil {
		return nil, callErr
	}
	arg__4158_419, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "update").Deref(), []vm.Value{arg__4150_414, vm.Keyword("loop-capture-syms-stack"), rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var or__x_2 vm.Value
		var s_3 vm.Value
		var body_caps_4 vm.Value
		var or__x_5 vm.Value
		var s_6 vm.Value
		var body_caps_7 vm.Value
		var arg__4152_12 vm.Value
		var or__x_13 vm.Value
		var s_14 vm.Value
		var body_caps_15 vm.Value
		var or__x_17 vm.Value
		var s_18 vm.Value
		var body_caps_19 vm.Value
		var head__4154_20 vm.Value
		var or__x_21 vm.Value
		var s_22 vm.Value
		var body_caps_23 vm.Value
		var head__4154_24 vm.Value
		var arg__4155_29 vm.Value
		var or__x_30 vm.Value
		var s_31 vm.Value
		var body_caps_32 vm.Value
		var head__4154_33 vm.Value
		var v34 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = or__x_2, s_3, body_caps_4, or__x_5, s_6, body_caps_7, arg__4152_12, or__x_13, s_14, body_caps_15, or__x_17, s_18, body_caps_19, head__4154_20, or__x_21, s_22, body_caps_23, head__4154_24, arg__4155_29, or__x_30, s_31, body_caps_32, head__4154_33, v34
		if vm.IsTruthy(arg0) {
			or__x_2 = arg0
			s_3 = arg0
			body_caps_4 = body_caps_367
			goto b1
		} else {
			or__x_5 = arg0
			s_6 = arg0
			body_caps_7 = body_caps_367
			goto b2
		}
	b1:
		;
		arg__4152_12 = or__x_2
		or__x_13 = or__x_2
		s_14 = s_3
		body_caps_15 = body_caps_4
		goto b3
	b2:
		;
		arg__4152_12 = vm.NewArrayVector([]vm.Value{})
		or__x_13 = or__x_5
		s_14 = s_6
		body_caps_15 = body_caps_7
		goto b3
	b3:
		;
		if vm.IsTruthy(s_14) {
			or__x_17 = s_14
			s_18 = s_14
			body_caps_19 = body_caps_15
			head__4154_20 = rt.LookupVar("clojure.core", "conj").Deref()
			goto b4
		} else {
			or__x_21 = s_14
			s_22 = s_14
			body_caps_23 = body_caps_15
			head__4154_24 = rt.LookupVar("clojure.core", "conj").Deref()
			goto b5
		}
	b4:
		;
		arg__4155_29 = or__x_17
		or__x_30 = or__x_17
		s_31 = s_18
		body_caps_32 = body_caps_19
		head__4154_33 = head__4154_20
		goto b6
	b5:
		;
		arg__4155_29 = vm.NewArrayVector([]vm.Value{})
		or__x_30 = or__x_21
		s_31 = s_22
		body_caps_32 = body_caps_23
		head__4154_33 = head__4154_24
		goto b6
	b6:
		;
		v34, callErr = rt.InvokeValue(head__4154_33, []vm.Value{arg__4155_29, body_caps_32})
		if callErr != nil {
			return nil, callErr
		}
		return v34, nil
	})})
	if callErr != nil {
		return nil, callErr
	}
	arg__4164_422, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_381})
	if callErr != nil {
		return nil, callErr
	}
	arg__4176_429, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_381})
	if callErr != nil {
		return nil, callErr
	}
	arg__4184_434, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "update").Deref(), []vm.Value{arg__4176_429, vm.Keyword("loop-headers"), rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var or__x_2 vm.Value
		var s_3 vm.Value
		var header_4 vm.Value
		var or__x_5 vm.Value
		var s_6 vm.Value
		var header_7 vm.Value
		var arg__4178_12 vm.Value
		var or__x_13 vm.Value
		var s_14 vm.Value
		var header_15 vm.Value
		var or__x_17 vm.Value
		var s_18 vm.Value
		var header_19 vm.Value
		var head__4180_20 vm.Value
		var or__x_21 vm.Value
		var s_22 vm.Value
		var header_23 vm.Value
		var head__4180_24 vm.Value
		var arg__4181_29 vm.Value
		var or__x_30 vm.Value
		var s_31 vm.Value
		var header_32 vm.Value
		var head__4180_33 vm.Value
		var v34 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = or__x_2, s_3, header_4, or__x_5, s_6, header_7, arg__4178_12, or__x_13, s_14, header_15, or__x_17, s_18, header_19, head__4180_20, or__x_21, s_22, header_23, head__4180_24, arg__4181_29, or__x_30, s_31, header_32, head__4180_33, v34
		if vm.IsTruthy(arg0) {
			or__x_2 = arg0
			s_3 = arg0
			header_4 = header_369
			goto b1
		} else {
			or__x_5 = arg0
			s_6 = arg0
			header_7 = header_369
			goto b2
		}
	b1:
		;
		arg__4178_12 = or__x_2
		or__x_13 = or__x_2
		s_14 = s_3
		header_15 = header_4
		goto b3
	b2:
		;
		arg__4178_12 = vm.NewArrayVector([]vm.Value{})
		or__x_13 = or__x_5
		s_14 = s_6
		header_15 = header_7
		goto b3
	b3:
		;
		if vm.IsTruthy(s_14) {
			or__x_17 = s_14
			s_18 = s_14
			header_19 = header_15
			head__4180_20 = rt.LookupVar("clojure.core", "conj").Deref()
			goto b4
		} else {
			or__x_21 = s_14
			s_22 = s_14
			header_23 = header_15
			head__4180_24 = rt.LookupVar("clojure.core", "conj").Deref()
			goto b5
		}
	b4:
		;
		arg__4181_29 = or__x_17
		or__x_30 = or__x_17
		s_31 = s_18
		header_32 = header_19
		head__4180_33 = head__4180_20
		goto b6
	b5:
		;
		arg__4181_29 = vm.NewArrayVector([]vm.Value{})
		or__x_30 = or__x_21
		s_31 = s_22
		header_32 = header_23
		head__4180_33 = head__4180_24
		goto b6
	b6:
		;
		v34, callErr = rt.InvokeValue(head__4180_33, []vm.Value{arg__4181_29, header_32})
		if callErr != nil {
			return nil, callErr
		}
		return v34, nil
	})})
	if callErr != nil {
		return nil, callErr
	}
	arg__4196_441, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_381})
	if callErr != nil {
		return nil, callErr
	}
	arg__4208_448, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_381})
	if callErr != nil {
		return nil, callErr
	}
	arg__4216_453, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "update").Deref(), []vm.Value{arg__4208_448, vm.Keyword("loop-headers"), rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var or__x_2 vm.Value
		var s_3 vm.Value
		var header_4 vm.Value
		var or__x_5 vm.Value
		var s_6 vm.Value
		var header_7 vm.Value
		var arg__4210_12 vm.Value
		var or__x_13 vm.Value
		var s_14 vm.Value
		var header_15 vm.Value
		var or__x_17 vm.Value
		var s_18 vm.Value
		var header_19 vm.Value
		var head__4212_20 vm.Value
		var or__x_21 vm.Value
		var s_22 vm.Value
		var header_23 vm.Value
		var head__4212_24 vm.Value
		var arg__4213_29 vm.Value
		var or__x_30 vm.Value
		var s_31 vm.Value
		var header_32 vm.Value
		var head__4212_33 vm.Value
		var v34 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = or__x_2, s_3, header_4, or__x_5, s_6, header_7, arg__4210_12, or__x_13, s_14, header_15, or__x_17, s_18, header_19, head__4212_20, or__x_21, s_22, header_23, head__4212_24, arg__4213_29, or__x_30, s_31, header_32, head__4212_33, v34
		if vm.IsTruthy(arg0) {
			or__x_2 = arg0
			s_3 = arg0
			header_4 = header_369
			goto b1
		} else {
			or__x_5 = arg0
			s_6 = arg0
			header_7 = header_369
			goto b2
		}
	b1:
		;
		arg__4210_12 = or__x_2
		or__x_13 = or__x_2
		s_14 = s_3
		header_15 = header_4
		goto b3
	b2:
		;
		arg__4210_12 = vm.NewArrayVector([]vm.Value{})
		or__x_13 = or__x_5
		s_14 = s_6
		header_15 = header_7
		goto b3
	b3:
		;
		if vm.IsTruthy(s_14) {
			or__x_17 = s_14
			s_18 = s_14
			header_19 = header_15
			head__4212_20 = rt.LookupVar("clojure.core", "conj").Deref()
			goto b4
		} else {
			or__x_21 = s_14
			s_22 = s_14
			header_23 = header_15
			head__4212_24 = rt.LookupVar("clojure.core", "conj").Deref()
			goto b5
		}
	b4:
		;
		arg__4213_29 = or__x_17
		or__x_30 = or__x_17
		s_31 = s_18
		header_32 = header_19
		head__4212_33 = head__4212_20
		goto b6
	b5:
		;
		arg__4213_29 = vm.NewArrayVector([]vm.Value{})
		or__x_30 = or__x_21
		s_31 = s_22
		header_32 = header_23
		head__4212_33 = head__4212_24
		goto b6
	b6:
		;
		v34, callErr = rt.InvokeValue(head__4212_33, []vm.Value{arg__4213_29, header_32})
		if callErr != nil {
			return nil, callErr
		}
		return v34, nil
	})})
	if callErr != nil {
		return nil, callErr
	}
	arg__4224_458, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "update").Deref(), []vm.Value{arg__4216_453, vm.Keyword("loop-capture-syms-stack"), rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var or__x_2 vm.Value
		var s_3 vm.Value
		var body_caps_4 vm.Value
		var or__x_5 vm.Value
		var s_6 vm.Value
		var body_caps_7 vm.Value
		var arg__4218_12 vm.Value
		var or__x_13 vm.Value
		var s_14 vm.Value
		var body_caps_15 vm.Value
		var or__x_17 vm.Value
		var s_18 vm.Value
		var body_caps_19 vm.Value
		var head__4220_20 vm.Value
		var or__x_21 vm.Value
		var s_22 vm.Value
		var body_caps_23 vm.Value
		var head__4220_24 vm.Value
		var arg__4221_29 vm.Value
		var or__x_30 vm.Value
		var s_31 vm.Value
		var body_caps_32 vm.Value
		var head__4220_33 vm.Value
		var v34 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = or__x_2, s_3, body_caps_4, or__x_5, s_6, body_caps_7, arg__4218_12, or__x_13, s_14, body_caps_15, or__x_17, s_18, body_caps_19, head__4220_20, or__x_21, s_22, body_caps_23, head__4220_24, arg__4221_29, or__x_30, s_31, body_caps_32, head__4220_33, v34
		if vm.IsTruthy(arg0) {
			or__x_2 = arg0
			s_3 = arg0
			body_caps_4 = body_caps_367
			goto b1
		} else {
			or__x_5 = arg0
			s_6 = arg0
			body_caps_7 = body_caps_367
			goto b2
		}
	b1:
		;
		arg__4218_12 = or__x_2
		or__x_13 = or__x_2
		s_14 = s_3
		body_caps_15 = body_caps_4
		goto b3
	b2:
		;
		arg__4218_12 = vm.NewArrayVector([]vm.Value{})
		or__x_13 = or__x_5
		s_14 = s_6
		body_caps_15 = body_caps_7
		goto b3
	b3:
		;
		if vm.IsTruthy(s_14) {
			or__x_17 = s_14
			s_18 = s_14
			body_caps_19 = body_caps_15
			head__4220_20 = rt.LookupVar("clojure.core", "conj").Deref()
			goto b4
		} else {
			or__x_21 = s_14
			s_22 = s_14
			body_caps_23 = body_caps_15
			head__4220_24 = rt.LookupVar("clojure.core", "conj").Deref()
			goto b5
		}
	b4:
		;
		arg__4221_29 = or__x_17
		or__x_30 = or__x_17
		s_31 = s_18
		body_caps_32 = body_caps_19
		head__4220_33 = head__4220_20
		goto b6
	b5:
		;
		arg__4221_29 = vm.NewArrayVector([]vm.Value{})
		or__x_30 = or__x_21
		s_31 = s_22
		body_caps_32 = body_caps_23
		head__4220_33 = head__4220_24
		goto b6
	b6:
		;
		v34, callErr = rt.InvokeValue(head__4220_33, []vm.Value{arg__4221_29, body_caps_32})
		if callErr != nil {
			return nil, callErr
		}
		return v34, nil
	})})
	if callErr != nil {
		return nil, callErr
	}
	v459, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "reset!").Deref(), []vm.Value{ctx_381, arg__4224_458})
	if callErr != nil {
		return nil, callErr
	}
	pre_locals_461, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "current-locals-flat").Deref(), []vm.Value{ctx_381})
	if callErr != nil {
		return nil, callErr
	}
	fs_462 = body_361
	last_val_463 = vm.NIL
	ctx_464 = ctx_381
	goto b9
b9:
	;
	v518, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{fs_462})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v518) {
		form_467 = form_359
		vec__3586_468 = vec__3586_360
		body_469 = body_361
		f_470 = f_362
		n_slots_471 = n_slots_363
		old_header_472 = old_header_364
		old_caps_473 = old_caps_365
		known_474 = known_366
		body_caps_475 = body_caps_367
		n_caps_476 = n_caps_368
		header_477 = header_369
		cap_param_ids_478 = cap_param_ids_370
		init_vals_479 = init_vals_371
		cap_vals_480 = cap_vals_372
		bt_481 = bt_373
		entry_end_482 = entry_end_374
		__483 = __375
		doseq_loop__3588_484 = doseq_loop__3588_376
		bindings_485 = bindings_377
		loop_param_ids_486 = loop_param_ids_378
		doseq_loop__3590_487 = doseq_loop__3590_380
		pre_locals_488 = pre_locals_461
		fs_489 = fs_462
		last_val_490 = last_val_463
		ctx_491 = ctx_464
		goto b10
	} else {
		form_492 = form_359
		vec__3586_493 = vec__3586_360
		body_494 = body_361
		f_495 = f_362
		n_slots_496 = n_slots_363
		old_header_497 = old_header_364
		old_caps_498 = old_caps_365
		known_499 = known_366
		body_caps_500 = body_caps_367
		n_caps_501 = n_caps_368
		header_502 = header_369
		cap_param_ids_503 = cap_param_ids_370
		init_vals_504 = init_vals_371
		cap_vals_505 = cap_vals_372
		bt_506 = bt_373
		entry_end_507 = entry_end_374
		__508 = __375
		doseq_loop__3588_509 = doseq_loop__3588_376
		bindings_510 = bindings_377
		loop_param_ids_511 = loop_param_ids_378
		doseq_loop__3590_512 = doseq_loop__3590_380
		pre_locals_513 = pre_locals_461
		fs_514 = fs_462
		last_val_515 = last_val_463
		ctx_516 = ctx_464
		goto b11
	}
b10:
	;
	v521, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{fs_489})
	if callErr != nil {
		return nil, callErr
	}
	arg__4237_523, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{fs_489})
	if callErr != nil {
		return nil, callErr
	}
	arg__4243_526, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{fs_489})
	if callErr != nil {
		return nil, callErr
	}
	v527, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-form").Deref(), []vm.Value{arg__4243_526, ctx_491})
	if callErr != nil {
		return nil, callErr
	}
	fs_462 = v521
	last_val_463 = v527
	ctx_464 = ctx_491
	goto b9
b11:
	;
	result_530 = last_val_515
	form_531 = form_492
	vec__3586_532 = vec__3586_493
	body_533 = body_494
	f_534 = f_495
	n_slots_535 = n_slots_496
	old_header_536 = old_header_497
	old_caps_537 = old_caps_498
	known_538 = known_499
	body_caps_539 = body_caps_500
	n_caps_540 = n_caps_501
	header_541 = header_502
	cap_param_ids_542 = cap_param_ids_503
	init_vals_543 = init_vals_504
	cap_vals_544 = cap_vals_505
	bt_545 = bt_506
	entry_end_546 = entry_end_507
	__547 = __508
	doseq_loop__3588_548 = doseq_loop__3588_509
	bindings_549 = bindings_510
	loop_param_ids_550 = loop_param_ids_511
	doseq_loop__3590_551 = doseq_loop__3590_512
	pre_locals_552 = pre_locals_513
	fs_553 = fs_514
	last_val_554 = last_val_515
	ctx_555 = ctx_516
	goto b12
b12:
	;
	post_locals_557, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "current-locals-flat").Deref(), []vm.Value{ctx_555})
	if callErr != nil {
		return nil, callErr
	}
	v559, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "pop-locals!").Deref(), []vm.Value{ctx_555})
	if callErr != nil {
		return nil, callErr
	}
	doseq_seq__3592_561, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{post_locals_557})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__3593_562 = doseq_seq__3592_561
	pre_locals_563 = pre_locals_552
	ctx_564 = ctx_555
	v976 = 0
	v985 = vm.NIL
	v994 = 1
	goto b13
b13:
	;
	if vm.IsTruthy(doseq_loop__3593_562) {
		result_566 = result_530
		form_567 = form_531
		vec__3586_568 = vec__3586_532
		body_569 = body_533
		f_570 = f_534
		n_slots_571 = n_slots_535
		old_header_572 = old_header_536
		old_caps_573 = old_caps_537
		known_574 = known_538
		body_caps_575 = body_caps_539
		n_caps_576 = n_caps_540
		header_577 = header_541
		cap_param_ids_578 = cap_param_ids_542
		init_vals_579 = init_vals_543
		cap_vals_580 = cap_vals_544
		bt_581 = bt_545
		entry_end_582 = entry_end_546
		__583 = __547
		doseq_loop__3588_584 = doseq_loop__3588_548
		bindings_585 = bindings_549
		loop_param_ids_586 = loop_param_ids_550
		doseq_loop__3590_587 = doseq_loop__3590_551
		fs_588 = fs_553
		last_val_589 = last_val_554
		post_locals_590 = post_locals_557
		doseq_seq__3592_591 = doseq_seq__3592_561
		doseq_loop__3593_592 = doseq_loop__3593_562
		pre_locals_593 = pre_locals_563
		ctx_594 = ctx_564
		v983 = v976
		v992 = v985
		v1001 = v994
		goto b14
	} else {
		result_595 = result_530
		form_596 = form_531
		vec__3586_597 = vec__3586_532
		body_598 = body_533
		f_599 = f_534
		n_slots_600 = n_slots_535
		old_header_601 = old_header_536
		old_caps_602 = old_caps_537
		known_603 = known_538
		body_caps_604 = body_caps_539
		n_caps_605 = n_caps_540
		header_606 = header_541
		cap_param_ids_607 = cap_param_ids_542
		init_vals_608 = init_vals_543
		cap_vals_609 = cap_vals_544
		bt_610 = bt_545
		entry_end_611 = entry_end_546
		__612 = __547
		doseq_loop__3588_613 = doseq_loop__3588_548
		bindings_614 = bindings_549
		loop_param_ids_615 = loop_param_ids_550
		doseq_loop__3590_616 = doseq_loop__3590_551
		fs_617 = fs_553
		last_val_618 = last_val_554
		post_locals_619 = post_locals_557
		doseq_seq__3592_620 = doseq_seq__3592_561
		doseq_loop__3593_621 = doseq_loop__3593_562
		pre_locals_622 = pre_locals_563
		ctx_623 = ctx_564
		v984 = v976
		v993 = v985
		v1002 = v994
		goto b15
	}
b14:
	;
	vec__3594_626, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{doseq_loop__3593_592})
	if callErr != nil {
		return nil, callErr
	}
	sym_632, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{vec__3594_626, vm.Int(v983), v992})
	if callErr != nil {
		return nil, callErr
	}
	val_638, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{vec__3594_626, vm.Int(v1001), v992})
	if callErr != nil {
		return nil, callErr
	}
	and__x_704, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "contains?").Deref(), []vm.Value{pre_locals_593, sym_632})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(and__x_704) {
		result_705 = result_566
		form_706 = form_567
		vec__3586_707 = vec__3586_568
		body_708 = body_569
		f_709 = f_570
		n_slots_710 = n_slots_571
		old_header_711 = old_header_572
		old_caps_712 = old_caps_573
		known_713 = known_574
		body_caps_714 = body_caps_575
		n_caps_715 = n_caps_576
		header_716 = header_577
		cap_param_ids_717 = cap_param_ids_578
		init_vals_718 = init_vals_579
		cap_vals_719 = cap_vals_580
		bt_720 = bt_581
		entry_end_721 = entry_end_582
		__722 = __583
		doseq_loop__3588_723 = doseq_loop__3588_584
		bindings_724 = bindings_585
		loop_param_ids_725 = loop_param_ids_586
		doseq_loop__3590_726 = doseq_loop__3590_587
		fs_727 = fs_588
		last_val_728 = last_val_589
		post_locals_729 = post_locals_590
		doseq_seq__3592_730 = doseq_seq__3592_591
		doseq_loop__3593_731 = doseq_loop__3593_592
		pre_locals_732 = pre_locals_593
		ctx_733 = ctx_594
		vec__3594_734 = vec__3594_626
		sym_735 = sym_632
		val_736 = val_638
		and__x_737 = and__x_704
		v977 = v983
		v986 = v992
		v995 = v1001
		goto b20
	} else {
		result_738 = result_566
		form_739 = form_567
		vec__3586_740 = vec__3586_568
		body_741 = body_569
		f_742 = f_570
		n_slots_743 = n_slots_571
		old_header_744 = old_header_572
		old_caps_745 = old_caps_573
		known_746 = known_574
		body_caps_747 = body_caps_575
		n_caps_748 = n_caps_576
		header_749 = header_577
		cap_param_ids_750 = cap_param_ids_578
		init_vals_751 = init_vals_579
		cap_vals_752 = cap_vals_580
		bt_753 = bt_581
		entry_end_754 = entry_end_582
		__755 = __583
		doseq_loop__3588_756 = doseq_loop__3588_584
		bindings_757 = bindings_585
		loop_param_ids_758 = loop_param_ids_586
		doseq_loop__3590_759 = doseq_loop__3590_587
		fs_760 = fs_588
		last_val_761 = last_val_589
		post_locals_762 = post_locals_590
		doseq_seq__3592_763 = doseq_seq__3592_591
		doseq_loop__3593_764 = doseq_loop__3593_592
		pre_locals_765 = pre_locals_593
		ctx_766 = ctx_594
		vec__3594_767 = vec__3594_626
		sym_768 = sym_632
		val_769 = val_638
		and__x_770 = and__x_704
		v978 = v983
		v987 = v992
		v996 = v1001
		goto b21
	}
b15:
	;
	v858 = vm.NIL
	result_859 = result_595
	form_860 = form_596
	vec__3586_861 = vec__3586_597
	body_862 = body_598
	f_863 = f_599
	n_slots_864 = n_slots_600
	old_header_865 = old_header_601
	old_caps_866 = old_caps_602
	known_867 = known_603
	body_caps_868 = body_caps_604
	n_caps_869 = n_caps_605
	header_870 = header_606
	cap_param_ids_871 = cap_param_ids_607
	init_vals_872 = init_vals_608
	cap_vals_873 = cap_vals_609
	bt_874 = bt_610
	entry_end_875 = entry_end_611
	__876 = __612
	doseq_loop__3588_877 = doseq_loop__3588_613
	bindings_878 = bindings_614
	loop_param_ids_879 = loop_param_ids_615
	doseq_loop__3590_880 = doseq_loop__3590_616
	fs_881 = fs_617
	last_val_882 = last_val_618
	post_locals_883 = post_locals_619
	doseq_seq__3592_884 = doseq_seq__3592_620
	doseq_loop__3593_885 = doseq_loop__3593_621
	pre_locals_886 = pre_locals_622
	ctx_887 = ctx_623
	goto b16
b16:
	;
	arg__4305_889, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_887})
	if callErr != nil {
		return nil, callErr
	}
	arg__4315_894, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_887})
	if callErr != nil {
		return nil, callErr
	}
	arg__4321_897, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "update").Deref(), []vm.Value{arg__4315_894, vm.Keyword("loop-headers"), rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var or__x_1 vm.Value
		var s_2 vm.Value
		var or__x_3 vm.Value
		var s_4 vm.Value
		var arg__4317_9 vm.Value
		var or__x_10 vm.Value
		var s_11 vm.Value
		var or__x_13 vm.Value
		var s_14 vm.Value
		var head__4318_15 vm.Value
		var or__x_16 vm.Value
		var s_17 vm.Value
		var head__4318_18 vm.Value
		var arg__4319_23 vm.Value
		var or__x_24 vm.Value
		var s_25 vm.Value
		var head__4318_26 vm.Value
		var v27 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = or__x_1, s_2, or__x_3, s_4, arg__4317_9, or__x_10, s_11, or__x_13, s_14, head__4318_15, or__x_16, s_17, head__4318_18, arg__4319_23, or__x_24, s_25, head__4318_26, v27
		if vm.IsTruthy(arg0) {
			or__x_1 = arg0
			s_2 = arg0
			goto b1
		} else {
			or__x_3 = arg0
			s_4 = arg0
			goto b2
		}
	b1:
		;
		arg__4317_9 = or__x_1
		or__x_10 = or__x_1
		s_11 = s_2
		goto b3
	b2:
		;
		arg__4317_9 = vm.NewArrayVector([]vm.Value{})
		or__x_10 = or__x_3
		s_11 = s_4
		goto b3
	b3:
		;
		if vm.IsTruthy(s_11) {
			or__x_13 = s_11
			s_14 = s_11
			head__4318_15 = rt.LookupVar("clojure.core", "pop").Deref()
			goto b4
		} else {
			or__x_16 = s_11
			s_17 = s_11
			head__4318_18 = rt.LookupVar("clojure.core", "pop").Deref()
			goto b5
		}
	b4:
		;
		arg__4319_23 = or__x_13
		or__x_24 = or__x_13
		s_25 = s_14
		head__4318_26 = head__4318_15
		goto b6
	b5:
		;
		arg__4319_23 = vm.NewArrayVector([]vm.Value{})
		or__x_24 = or__x_16
		s_25 = s_17
		head__4318_26 = head__4318_18
		goto b6
	b6:
		;
		v27, callErr = rt.InvokeValue(head__4318_26, []vm.Value{arg__4319_23})
		if callErr != nil {
			return nil, callErr
		}
		return v27, nil
	})})
	if callErr != nil {
		return nil, callErr
	}
	arg__4331_902, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_887})
	if callErr != nil {
		return nil, callErr
	}
	arg__4341_907, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_887})
	if callErr != nil {
		return nil, callErr
	}
	arg__4347_910, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "update").Deref(), []vm.Value{arg__4341_907, vm.Keyword("loop-headers"), rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var or__x_1 vm.Value
		var s_2 vm.Value
		var or__x_3 vm.Value
		var s_4 vm.Value
		var arg__4343_9 vm.Value
		var or__x_10 vm.Value
		var s_11 vm.Value
		var or__x_13 vm.Value
		var s_14 vm.Value
		var head__4344_15 vm.Value
		var or__x_16 vm.Value
		var s_17 vm.Value
		var head__4344_18 vm.Value
		var arg__4345_23 vm.Value
		var or__x_24 vm.Value
		var s_25 vm.Value
		var head__4344_26 vm.Value
		var v27 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = or__x_1, s_2, or__x_3, s_4, arg__4343_9, or__x_10, s_11, or__x_13, s_14, head__4344_15, or__x_16, s_17, head__4344_18, arg__4345_23, or__x_24, s_25, head__4344_26, v27
		if vm.IsTruthy(arg0) {
			or__x_1 = arg0
			s_2 = arg0
			goto b1
		} else {
			or__x_3 = arg0
			s_4 = arg0
			goto b2
		}
	b1:
		;
		arg__4343_9 = or__x_1
		or__x_10 = or__x_1
		s_11 = s_2
		goto b3
	b2:
		;
		arg__4343_9 = vm.NewArrayVector([]vm.Value{})
		or__x_10 = or__x_3
		s_11 = s_4
		goto b3
	b3:
		;
		if vm.IsTruthy(s_11) {
			or__x_13 = s_11
			s_14 = s_11
			head__4344_15 = rt.LookupVar("clojure.core", "pop").Deref()
			goto b4
		} else {
			or__x_16 = s_11
			s_17 = s_11
			head__4344_18 = rt.LookupVar("clojure.core", "pop").Deref()
			goto b5
		}
	b4:
		;
		arg__4345_23 = or__x_13
		or__x_24 = or__x_13
		s_25 = s_14
		head__4344_26 = head__4344_15
		goto b6
	b5:
		;
		arg__4345_23 = vm.NewArrayVector([]vm.Value{})
		or__x_24 = or__x_16
		s_25 = s_17
		head__4344_26 = head__4344_18
		goto b6
	b6:
		;
		v27, callErr = rt.InvokeValue(head__4344_26, []vm.Value{arg__4345_23})
		if callErr != nil {
			return nil, callErr
		}
		return v27, nil
	})})
	if callErr != nil {
		return nil, callErr
	}
	arg__4353_913, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "update").Deref(), []vm.Value{arg__4347_910, vm.Keyword("loop-capture-syms-stack"), rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var or__x_1 vm.Value
		var s_2 vm.Value
		var or__x_3 vm.Value
		var s_4 vm.Value
		var arg__4349_9 vm.Value
		var or__x_10 vm.Value
		var s_11 vm.Value
		var or__x_13 vm.Value
		var s_14 vm.Value
		var head__4350_15 vm.Value
		var or__x_16 vm.Value
		var s_17 vm.Value
		var head__4350_18 vm.Value
		var arg__4351_23 vm.Value
		var or__x_24 vm.Value
		var s_25 vm.Value
		var head__4350_26 vm.Value
		var v27 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = or__x_1, s_2, or__x_3, s_4, arg__4349_9, or__x_10, s_11, or__x_13, s_14, head__4350_15, or__x_16, s_17, head__4350_18, arg__4351_23, or__x_24, s_25, head__4350_26, v27
		if vm.IsTruthy(arg0) {
			or__x_1 = arg0
			s_2 = arg0
			goto b1
		} else {
			or__x_3 = arg0
			s_4 = arg0
			goto b2
		}
	b1:
		;
		arg__4349_9 = or__x_1
		or__x_10 = or__x_1
		s_11 = s_2
		goto b3
	b2:
		;
		arg__4349_9 = vm.NewArrayVector([]vm.Value{})
		or__x_10 = or__x_3
		s_11 = s_4
		goto b3
	b3:
		;
		if vm.IsTruthy(s_11) {
			or__x_13 = s_11
			s_14 = s_11
			head__4350_15 = rt.LookupVar("clojure.core", "pop").Deref()
			goto b4
		} else {
			or__x_16 = s_11
			s_17 = s_11
			head__4350_18 = rt.LookupVar("clojure.core", "pop").Deref()
			goto b5
		}
	b4:
		;
		arg__4351_23 = or__x_13
		or__x_24 = or__x_13
		s_25 = s_14
		head__4350_26 = head__4350_15
		goto b6
	b5:
		;
		arg__4351_23 = vm.NewArrayVector([]vm.Value{})
		or__x_24 = or__x_16
		s_25 = s_17
		head__4350_26 = head__4350_18
		goto b6
	b6:
		;
		v27, callErr = rt.InvokeValue(head__4350_26, []vm.Value{arg__4351_23})
		if callErr != nil {
			return nil, callErr
		}
		return v27, nil
	})})
	if callErr != nil {
		return nil, callErr
	}
	arg__4359_916, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_887})
	if callErr != nil {
		return nil, callErr
	}
	arg__4369_921, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_887})
	if callErr != nil {
		return nil, callErr
	}
	arg__4375_924, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "update").Deref(), []vm.Value{arg__4369_921, vm.Keyword("loop-headers"), rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var or__x_1 vm.Value
		var s_2 vm.Value
		var or__x_3 vm.Value
		var s_4 vm.Value
		var arg__4371_9 vm.Value
		var or__x_10 vm.Value
		var s_11 vm.Value
		var or__x_13 vm.Value
		var s_14 vm.Value
		var head__4372_15 vm.Value
		var or__x_16 vm.Value
		var s_17 vm.Value
		var head__4372_18 vm.Value
		var arg__4373_23 vm.Value
		var or__x_24 vm.Value
		var s_25 vm.Value
		var head__4372_26 vm.Value
		var v27 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = or__x_1, s_2, or__x_3, s_4, arg__4371_9, or__x_10, s_11, or__x_13, s_14, head__4372_15, or__x_16, s_17, head__4372_18, arg__4373_23, or__x_24, s_25, head__4372_26, v27
		if vm.IsTruthy(arg0) {
			or__x_1 = arg0
			s_2 = arg0
			goto b1
		} else {
			or__x_3 = arg0
			s_4 = arg0
			goto b2
		}
	b1:
		;
		arg__4371_9 = or__x_1
		or__x_10 = or__x_1
		s_11 = s_2
		goto b3
	b2:
		;
		arg__4371_9 = vm.NewArrayVector([]vm.Value{})
		or__x_10 = or__x_3
		s_11 = s_4
		goto b3
	b3:
		;
		if vm.IsTruthy(s_11) {
			or__x_13 = s_11
			s_14 = s_11
			head__4372_15 = rt.LookupVar("clojure.core", "pop").Deref()
			goto b4
		} else {
			or__x_16 = s_11
			s_17 = s_11
			head__4372_18 = rt.LookupVar("clojure.core", "pop").Deref()
			goto b5
		}
	b4:
		;
		arg__4373_23 = or__x_13
		or__x_24 = or__x_13
		s_25 = s_14
		head__4372_26 = head__4372_15
		goto b6
	b5:
		;
		arg__4373_23 = vm.NewArrayVector([]vm.Value{})
		or__x_24 = or__x_16
		s_25 = s_17
		head__4372_26 = head__4372_18
		goto b6
	b6:
		;
		v27, callErr = rt.InvokeValue(head__4372_26, []vm.Value{arg__4373_23})
		if callErr != nil {
			return nil, callErr
		}
		return v27, nil
	})})
	if callErr != nil {
		return nil, callErr
	}
	arg__4385_929, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_887})
	if callErr != nil {
		return nil, callErr
	}
	arg__4395_934, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_887})
	if callErr != nil {
		return nil, callErr
	}
	arg__4401_937, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "update").Deref(), []vm.Value{arg__4395_934, vm.Keyword("loop-headers"), rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var or__x_1 vm.Value
		var s_2 vm.Value
		var or__x_3 vm.Value
		var s_4 vm.Value
		var arg__4397_9 vm.Value
		var or__x_10 vm.Value
		var s_11 vm.Value
		var or__x_13 vm.Value
		var s_14 vm.Value
		var head__4398_15 vm.Value
		var or__x_16 vm.Value
		var s_17 vm.Value
		var head__4398_18 vm.Value
		var arg__4399_23 vm.Value
		var or__x_24 vm.Value
		var s_25 vm.Value
		var head__4398_26 vm.Value
		var v27 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = or__x_1, s_2, or__x_3, s_4, arg__4397_9, or__x_10, s_11, or__x_13, s_14, head__4398_15, or__x_16, s_17, head__4398_18, arg__4399_23, or__x_24, s_25, head__4398_26, v27
		if vm.IsTruthy(arg0) {
			or__x_1 = arg0
			s_2 = arg0
			goto b1
		} else {
			or__x_3 = arg0
			s_4 = arg0
			goto b2
		}
	b1:
		;
		arg__4397_9 = or__x_1
		or__x_10 = or__x_1
		s_11 = s_2
		goto b3
	b2:
		;
		arg__4397_9 = vm.NewArrayVector([]vm.Value{})
		or__x_10 = or__x_3
		s_11 = s_4
		goto b3
	b3:
		;
		if vm.IsTruthy(s_11) {
			or__x_13 = s_11
			s_14 = s_11
			head__4398_15 = rt.LookupVar("clojure.core", "pop").Deref()
			goto b4
		} else {
			or__x_16 = s_11
			s_17 = s_11
			head__4398_18 = rt.LookupVar("clojure.core", "pop").Deref()
			goto b5
		}
	b4:
		;
		arg__4399_23 = or__x_13
		or__x_24 = or__x_13
		s_25 = s_14
		head__4398_26 = head__4398_15
		goto b6
	b5:
		;
		arg__4399_23 = vm.NewArrayVector([]vm.Value{})
		or__x_24 = or__x_16
		s_25 = s_17
		head__4398_26 = head__4398_18
		goto b6
	b6:
		;
		v27, callErr = rt.InvokeValue(head__4398_26, []vm.Value{arg__4399_23})
		if callErr != nil {
			return nil, callErr
		}
		return v27, nil
	})})
	if callErr != nil {
		return nil, callErr
	}
	arg__4407_940, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "update").Deref(), []vm.Value{arg__4401_937, vm.Keyword("loop-capture-syms-stack"), rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var or__x_1 vm.Value
		var s_2 vm.Value
		var or__x_3 vm.Value
		var s_4 vm.Value
		var arg__4403_9 vm.Value
		var or__x_10 vm.Value
		var s_11 vm.Value
		var or__x_13 vm.Value
		var s_14 vm.Value
		var head__4404_15 vm.Value
		var or__x_16 vm.Value
		var s_17 vm.Value
		var head__4404_18 vm.Value
		var arg__4405_23 vm.Value
		var or__x_24 vm.Value
		var s_25 vm.Value
		var head__4404_26 vm.Value
		var v27 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = or__x_1, s_2, or__x_3, s_4, arg__4403_9, or__x_10, s_11, or__x_13, s_14, head__4404_15, or__x_16, s_17, head__4404_18, arg__4405_23, or__x_24, s_25, head__4404_26, v27
		if vm.IsTruthy(arg0) {
			or__x_1 = arg0
			s_2 = arg0
			goto b1
		} else {
			or__x_3 = arg0
			s_4 = arg0
			goto b2
		}
	b1:
		;
		arg__4403_9 = or__x_1
		or__x_10 = or__x_1
		s_11 = s_2
		goto b3
	b2:
		;
		arg__4403_9 = vm.NewArrayVector([]vm.Value{})
		or__x_10 = or__x_3
		s_11 = s_4
		goto b3
	b3:
		;
		if vm.IsTruthy(s_11) {
			or__x_13 = s_11
			s_14 = s_11
			head__4404_15 = rt.LookupVar("clojure.core", "pop").Deref()
			goto b4
		} else {
			or__x_16 = s_11
			s_17 = s_11
			head__4404_18 = rt.LookupVar("clojure.core", "pop").Deref()
			goto b5
		}
	b4:
		;
		arg__4405_23 = or__x_13
		or__x_24 = or__x_13
		s_25 = s_14
		head__4404_26 = head__4404_15
		goto b6
	b5:
		;
		arg__4405_23 = vm.NewArrayVector([]vm.Value{})
		or__x_24 = or__x_16
		s_25 = s_17
		head__4404_26 = head__4404_18
		goto b6
	b6:
		;
		v27, callErr = rt.InvokeValue(head__4404_26, []vm.Value{arg__4405_23})
		if callErr != nil {
			return nil, callErr
		}
		return v27, nil
	})})
	if callErr != nil {
		return nil, callErr
	}
	v941, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "reset!").Deref(), []vm.Value{ctx_887, arg__4407_940})
	if callErr != nil {
		return nil, callErr
	}
	return result_859, nil
b17:
	;
	v816, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "bind-local!").Deref(), []vm.Value{ctx_667, sym_669, val_670})
	if callErr != nil {
		return nil, callErr
	}
	v820 = v816
	result_821 = result_639
	form_822 = form_640
	vec__3586_823 = vec__3586_641
	body_824 = body_642
	f_825 = f_643
	n_slots_826 = n_slots_644
	old_header_827 = old_header_645
	old_caps_828 = old_caps_646
	known_829 = known_647
	body_caps_830 = body_caps_648
	n_caps_831 = n_caps_649
	header_832 = header_650
	cap_param_ids_833 = cap_param_ids_651
	init_vals_834 = init_vals_652
	cap_vals_835 = cap_vals_653
	bt_836 = bt_654
	entry_end_837 = entry_end_655
	__838 = __656
	doseq_loop__3588_839 = doseq_loop__3588_657
	bindings_840 = bindings_658
	loop_param_ids_841 = loop_param_ids_659
	doseq_loop__3590_842 = doseq_loop__3590_660
	fs_843 = fs_661
	last_val_844 = last_val_662
	post_locals_845 = post_locals_663
	doseq_seq__3592_846 = doseq_seq__3592_664
	doseq_loop__3593_847 = doseq_loop__3593_665
	pre_locals_848 = pre_locals_666
	ctx_849 = ctx_667
	vec__3594_850 = vec__3594_668
	sym_851 = sym_669
	val_852 = val_670
	v982 = v980
	v991 = v989
	v1000 = v998
	goto b19
b18:
	;
	v820 = vm.NIL
	result_821 = result_671
	form_822 = form_672
	vec__3586_823 = vec__3586_673
	body_824 = body_674
	f_825 = f_675
	n_slots_826 = n_slots_676
	old_header_827 = old_header_677
	old_caps_828 = old_caps_678
	known_829 = known_679
	body_caps_830 = body_caps_680
	n_caps_831 = n_caps_681
	header_832 = header_682
	cap_param_ids_833 = cap_param_ids_683
	init_vals_834 = init_vals_684
	cap_vals_835 = cap_vals_685
	bt_836 = bt_686
	entry_end_837 = entry_end_687
	__838 = __688
	doseq_loop__3588_839 = doseq_loop__3588_689
	bindings_840 = bindings_690
	loop_param_ids_841 = loop_param_ids_691
	doseq_loop__3590_842 = doseq_loop__3590_692
	fs_843 = fs_693
	last_val_844 = last_val_694
	post_locals_845 = post_locals_695
	doseq_seq__3592_846 = doseq_seq__3592_696
	doseq_loop__3593_847 = doseq_loop__3593_697
	pre_locals_848 = pre_locals_698
	ctx_849 = ctx_699
	vec__3594_850 = vec__3594_700
	sym_851 = sym_701
	val_852 = val_702
	v982 = v979
	v991 = v988
	v1000 = v997
	goto b19
b19:
	;
	v854, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "next").Deref(), []vm.Value{doseq_loop__3593_847})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__3593_562 = v854
	pre_locals_563 = pre_locals_848
	ctx_564 = ctx_849
	v976 = v982
	v985 = v991
	v994 = v1000
	goto b13
b20:
	;
	arg__4282_773, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{pre_locals_732, sym_735})
	if callErr != nil {
		return nil, callErr
	}
	arg__4290_776, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{pre_locals_732, sym_735})
	if callErr != nil {
		return nil, callErr
	}
	v777, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not=").Deref(), []vm.Value{val_736, arg__4290_776})
	if callErr != nil {
		return nil, callErr
	}
	v780 = v777
	result_781 = result_705
	form_782 = form_706
	vec__3586_783 = vec__3586_707
	body_784 = body_708
	f_785 = f_709
	n_slots_786 = n_slots_710
	old_header_787 = old_header_711
	old_caps_788 = old_caps_712
	known_789 = known_713
	body_caps_790 = body_caps_714
	n_caps_791 = n_caps_715
	header_792 = header_716
	cap_param_ids_793 = cap_param_ids_717
	init_vals_794 = init_vals_718
	cap_vals_795 = cap_vals_719
	bt_796 = bt_720
	entry_end_797 = entry_end_721
	__798 = __722
	doseq_loop__3588_799 = doseq_loop__3588_723
	bindings_800 = bindings_724
	loop_param_ids_801 = loop_param_ids_725
	doseq_loop__3590_802 = doseq_loop__3590_726
	fs_803 = fs_727
	last_val_804 = last_val_728
	post_locals_805 = post_locals_729
	doseq_seq__3592_806 = doseq_seq__3592_730
	doseq_loop__3593_807 = doseq_loop__3593_731
	pre_locals_808 = pre_locals_732
	ctx_809 = ctx_733
	vec__3594_810 = vec__3594_734
	sym_811 = sym_735
	val_812 = val_736
	and__x_813 = and__x_737
	v981 = v977
	v990 = v986
	v999 = v995
	goto b22
b21:
	;
	v780 = and__x_770
	result_781 = result_738
	form_782 = form_739
	vec__3586_783 = vec__3586_740
	body_784 = body_741
	f_785 = f_742
	n_slots_786 = n_slots_743
	old_header_787 = old_header_744
	old_caps_788 = old_caps_745
	known_789 = known_746
	body_caps_790 = body_caps_747
	n_caps_791 = n_caps_748
	header_792 = header_749
	cap_param_ids_793 = cap_param_ids_750
	init_vals_794 = init_vals_751
	cap_vals_795 = cap_vals_752
	bt_796 = bt_753
	entry_end_797 = entry_end_754
	__798 = __755
	doseq_loop__3588_799 = doseq_loop__3588_756
	bindings_800 = bindings_757
	loop_param_ids_801 = loop_param_ids_758
	doseq_loop__3590_802 = doseq_loop__3590_759
	fs_803 = fs_760
	last_val_804 = last_val_761
	post_locals_805 = post_locals_762
	doseq_seq__3592_806 = doseq_seq__3592_763
	doseq_loop__3593_807 = doseq_loop__3593_764
	pre_locals_808 = pre_locals_765
	ctx_809 = ctx_766
	vec__3594_810 = vec__3594_767
	sym_811 = sym_768
	val_812 = val_769
	and__x_813 = and__x_770
	v981 = v978
	v990 = v987
	v999 = v996
	goto b22
b22:
	;
	if vm.IsTruthy(v780) {
		result_639 = result_781
		form_640 = form_782
		vec__3586_641 = vec__3586_783
		body_642 = body_784
		f_643 = f_785
		n_slots_644 = n_slots_786
		old_header_645 = old_header_787
		old_caps_646 = old_caps_788
		known_647 = known_789
		body_caps_648 = body_caps_790
		n_caps_649 = n_caps_791
		header_650 = header_792
		cap_param_ids_651 = cap_param_ids_793
		init_vals_652 = init_vals_794
		cap_vals_653 = cap_vals_795
		bt_654 = bt_796
		entry_end_655 = entry_end_797
		__656 = __798
		doseq_loop__3588_657 = doseq_loop__3588_799
		bindings_658 = bindings_800
		loop_param_ids_659 = loop_param_ids_801
		doseq_loop__3590_660 = doseq_loop__3590_802
		fs_661 = fs_803
		last_val_662 = last_val_804
		post_locals_663 = post_locals_805
		doseq_seq__3592_664 = doseq_seq__3592_806
		doseq_loop__3593_665 = doseq_loop__3593_807
		pre_locals_666 = pre_locals_808
		ctx_667 = ctx_809
		vec__3594_668 = vec__3594_810
		sym_669 = sym_811
		val_670 = val_812
		v980 = v981
		v989 = v990
		v998 = v999
		goto b17
	} else {
		result_671 = result_781
		form_672 = form_782
		vec__3586_673 = vec__3586_783
		body_674 = body_784
		f_675 = f_785
		n_slots_676 = n_slots_786
		old_header_677 = old_header_787
		old_caps_678 = old_caps_788
		known_679 = known_789
		body_caps_680 = body_caps_790
		n_caps_681 = n_caps_791
		header_682 = header_792
		cap_param_ids_683 = cap_param_ids_793
		init_vals_684 = init_vals_794
		cap_vals_685 = cap_vals_795
		bt_686 = bt_796
		entry_end_687 = entry_end_797
		__688 = __798
		doseq_loop__3588_689 = doseq_loop__3588_799
		bindings_690 = bindings_800
		loop_param_ids_691 = loop_param_ids_801
		doseq_loop__3590_692 = doseq_loop__3590_802
		fs_693 = fs_803
		last_val_694 = last_val_804
		post_locals_695 = post_locals_805
		doseq_seq__3592_696 = doseq_seq__3592_806
		doseq_loop__3593_697 = doseq_loop__3593_807
		pre_locals_698 = pre_locals_808
		ctx_699 = ctx_809
		vec__3594_700 = vec__3594_810
		sym_701 = sym_811
		val_702 = val_812
		v979 = v981
		v988 = v990
		v997 = v999
		goto b18
	}
}
func build_map(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__4411_8 vm.Value
	var arg__4416_11 vm.Value
	var v12 vm.Value
	var form_3 vm.Value
	var ctx_4 vm.Value
	var arg__4421_15 vm.Value
	var arg__4430_20 vm.Value
	var v23 vm.Value
	var form_5 vm.Value
	var ctx_6 vm.Value
	var arg__4456_27 vm.Value
	var arg__4480_31 vm.Value
	var all_const_QMARK__32 vm.Value
	var arg__4494_35 vm.Value
	var arg__4509_39 vm.Value
	var arg__4510_40 vm.Value
	var arg__4525_44 vm.Value
	var arg__4540_48 vm.Value
	var arg__4541_49 vm.Value
	var pairs_50 vm.Value
	var v98 vm.Value
	var form_99 vm.Value
	var ctx_100 vm.Value
	var form_51 vm.Value
	var ctx_52 vm.Value
	var all_const_QMARK__53 vm.Value
	var pairs_54 vm.Value
	var arg__4546_61 vm.Value
	var arg__4555_66 vm.Value
	var v69 vm.Value
	var form_55 vm.Value
	var ctx_56 vm.Value
	var all_const_QMARK__57 vm.Value
	var pairs_58 vm.Value
	var arg__4563_72 vm.Value
	var arg__4569_78 vm.Value
	var arg__4575_81 vm.Value
	var arg__4581_87 vm.Value
	var fn_id_88 vm.Value
	var v90 vm.Value
	var v92 vm.Value
	var form_93 vm.Value
	var ctx_94 vm.Value
	var all_const_QMARK__95 vm.Value
	var pairs_96 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__4411_8, arg__4416_11, v12, form_3, ctx_4, arg__4421_15, arg__4430_20, v23, form_5, ctx_6, arg__4456_27, arg__4480_31, all_const_QMARK__32, arg__4494_35, arg__4509_39, arg__4510_40, arg__4525_44, arg__4540_48, arg__4541_49, pairs_50, v98, form_99, ctx_100, form_51, ctx_52, all_const_QMARK__53, pairs_54, arg__4546_61, arg__4555_66, v69, form_55, ctx_56, all_const_QMARK__57, pairs_58, arg__4563_72, arg__4569_78, arg__4575_81, arg__4581_87, fn_id_88, v90, v92, form_93, ctx_94, all_const_QMARK__95, pairs_96
	arg__4411_8, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__4416_11, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "zero?").Deref(), []vm.Value{arg__4416_11})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v12) {
		form_3 = arg0
		ctx_4 = arg1
		goto b1
	} else {
		form_5 = arg0
		ctx_6 = arg1
		goto b2
	}
b1:
	;
	arg__4421_15, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_4})
	if callErr != nil {
		return nil, callErr
	}
	arg__4430_20, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_4})
	if callErr != nil {
		return nil, callErr
	}
	v23, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "add-inst!").Deref(), []vm.Value{ctx_4, arg__4430_20, vm.Keyword("const"), vm.NewArrayVector([]vm.Value{}), form_3})
	if callErr != nil {
		return nil, callErr
	}
	v98 = v23
	form_99 = form_3
	ctx_100 = ctx_4
	goto b3
b2:
	;
	arg__4456_27, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{form_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__4480_31, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{form_5})
	if callErr != nil {
		return nil, callErr
	}
	all_const_QMARK__32, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "every?").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var arg__4461_2 vm.Value
		var arg__4466_5 vm.Value
		var and__x_6 vm.Value
		var e_7 vm.Value
		var and__x_8 vm.Value
		var arg__4470_13 vm.Value
		var arg__4475_16 vm.Value
		var v17 vm.Value
		var e_9 vm.Value
		var and__x_10 vm.Value
		var v20 vm.Value
		var e_21 vm.Value
		var and__x_22 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _, _, _, _, _, _ = arg__4461_2, arg__4466_5, and__x_6, e_7, and__x_8, arg__4470_13, arg__4475_16, v17, e_9, and__x_10, v20, e_21, and__x_22
		arg__4461_2, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		arg__4466_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		and__x_6, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "is-literal?").Deref(), []vm.Value{arg__4466_5})
		if callErr != nil {
			return nil, callErr
		}
		if vm.IsTruthy(and__x_6) {
			e_7 = arg0
			and__x_8 = and__x_6
			goto b1
		} else {
			e_9 = arg0
			and__x_10 = and__x_6
			goto b2
		}
	b1:
		;
		arg__4470_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "second").Deref(), []vm.Value{e_7})
		if callErr != nil {
			return nil, callErr
		}
		arg__4475_16, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "second").Deref(), []vm.Value{e_7})
		if callErr != nil {
			return nil, callErr
		}
		v17, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "is-literal?").Deref(), []vm.Value{arg__4475_16})
		if callErr != nil {
			return nil, callErr
		}
		v20 = v17
		e_21 = e_7
		and__x_22 = and__x_8
		goto b3
	b2:
		;
		v20 = and__x_10
		e_21 = e_9
		and__x_22 = and__x_10
		goto b3
	b3:
		;
		return v20, nil
	}), arg__4480_31})
	if callErr != nil {
		return nil, callErr
	}
	arg__4494_35, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{form_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__4509_39, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{form_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__4510_40, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapcat").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var arg__4500_3 vm.Value
		var arg__4504_5 vm.Value
		var v6 vm.Value
		var callErr error
		_, _, _ = arg__4500_3, arg__4504_5, v6
		arg__4500_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		arg__4504_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "second").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		v6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__4500_3, arg__4504_5})
		if callErr != nil {
			return nil, callErr
		}
		return v6, nil
	}), arg__4509_39})
	if callErr != nil {
		return nil, callErr
	}
	arg__4525_44, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{form_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__4540_48, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{form_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__4541_49, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapcat").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var arg__4531_3 vm.Value
		var arg__4535_5 vm.Value
		var v6 vm.Value
		var callErr error
		_, _, _ = arg__4531_3, arg__4535_5, v6
		arg__4531_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		arg__4535_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "second").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		v6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__4531_3, arg__4535_5})
		if callErr != nil {
			return nil, callErr
		}
		return v6, nil
	}), arg__4540_48})
	if callErr != nil {
		return nil, callErr
	}
	pairs_50, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg__4541_49})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(all_const_QMARK__32) {
		form_51 = form_5
		ctx_52 = ctx_6
		all_const_QMARK__53 = all_const_QMARK__32
		pairs_54 = pairs_50
		goto b4
	} else {
		form_55 = form_5
		ctx_56 = ctx_6
		all_const_QMARK__57 = all_const_QMARK__32
		pairs_58 = pairs_50
		goto b5
	}
b3:
	;
	return v98, nil
b4:
	;
	arg__4546_61, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_52})
	if callErr != nil {
		return nil, callErr
	}
	arg__4555_66, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_52})
	if callErr != nil {
		return nil, callErr
	}
	v69, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "add-inst!").Deref(), []vm.Value{ctx_52, arg__4555_66, vm.Keyword("const"), vm.NewArrayVector([]vm.Value{}), form_51})
	if callErr != nil {
		return nil, callErr
	}
	v92 = v69
	form_93 = form_51
	ctx_94 = ctx_52
	all_const_QMARK__95 = all_const_QMARK__53
	pairs_96 = pairs_54
	goto b6
b5:
	;
	arg__4563_72, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_56})
	if callErr != nil {
		return nil, callErr
	}
	arg__4569_78, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "resolve").Deref(), []vm.Value{vm.Symbol("array-map")})
	if callErr != nil {
		return nil, callErr
	}
	arg__4575_81, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_56})
	if callErr != nil {
		return nil, callErr
	}
	arg__4581_87, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "resolve").Deref(), []vm.Value{vm.Symbol("array-map")})
	if callErr != nil {
		return nil, callErr
	}
	fn_id_88, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "add-inst!").Deref(), []vm.Value{ctx_56, arg__4575_81, vm.Keyword("load-var"), vm.NewArrayVector([]vm.Value{}), arg__4581_87})
	if callErr != nil {
		return nil, callErr
	}
	v90, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-call-with-head").Deref(), []vm.Value{fn_id_88, pairs_58, ctx_56})
	if callErr != nil {
		return nil, callErr
	}
	v92 = v90
	form_93 = form_55
	ctx_94 = ctx_56
	all_const_QMARK__95 = all_const_QMARK__57
	pairs_96 = pairs_58
	goto b6
b6:
	;
	v98 = v92
	form_99 = form_93
	ctx_100 = ctx_94
	goto b3
}
func build_quote(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var __8 vm.Value
	var x_14 vm.Value
	var arg__4608_16 vm.Value
	var arg__4617_21 vm.Value
	var v24 vm.Value
	var callErr error
	_, _, _, _, _ = __8, x_14, arg__4608_16, arg__4617_21, v24
	__8, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(0), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	x_14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(1), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	arg__4608_16, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__4617_21, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	v24, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "add-inst!").Deref(), []vm.Value{arg1, arg__4617_21, vm.Keyword("const"), vm.NewArrayVector([]vm.Value{}), x_14})
	if callErr != nil {
		return nil, callErr
	}
	return v24, nil
}
func build_recur(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__4624_4 vm.Value
	var arg__4630_8 vm.Value
	var headers_10 vm.Value
	var f_12 vm.Value
	var form_13 vm.Value
	var ctx_14 vm.Value
	var headers_15 vm.Value
	var f_16 vm.Value
	var header_44 vm.Value
	var arg__4644_46 vm.Value
	var arg__4650_50 vm.Value
	var cap_syms_stack_52 vm.Value
	var or__x_54 vm.Value
	var form_17 vm.Value
	var ctx_18 vm.Value
	var headers_19 vm.Value
	var f_20 vm.Value
	var arg__4733_122 vm.Value
	var arg__4739_126 vm.Value
	var arg_syms_128 vm.Value
	var v176 vm.Value
	var form_177 vm.Value
	var ctx_178 vm.Value
	var headers_179 vm.Value
	var f_180 vm.Value
	var form_21 vm.Value
	var ctx_22 vm.Value
	var and__x_23 vm.Value
	var headers_24 vm.Value
	var f_25 vm.Value
	var v33 vm.Value
	var form_26 vm.Value
	var ctx_27 vm.Value
	var and__x_28 vm.Value
	var headers_29 vm.Value
	var f_30 vm.Value
	var v36 vm.Value
	var form_37 vm.Value
	var ctx_38 vm.Value
	var and__x_39 vm.Value
	var headers_40 vm.Value
	var f_41 vm.Value
	var form_55 vm.Value
	var ctx_56 vm.Value
	var headers_57 vm.Value
	var f_58 vm.Value
	var header_59 vm.Value
	var cap_syms_stack_60 vm.Value
	var or__x_61 vm.Value
	var form_62 vm.Value
	var ctx_63 vm.Value
	var headers_64 vm.Value
	var f_65 vm.Value
	var header_66 vm.Value
	var cap_syms_stack_67 vm.Value
	var or__x_68 vm.Value
	var cap_syms_73 vm.Value
	var form_74 vm.Value
	var ctx_75 vm.Value
	var headers_76 vm.Value
	var f_77 vm.Value
	var header_78 vm.Value
	var cap_syms_stack_79 vm.Value
	var or__x_80 vm.Value
	var cap_vals_88 vm.Value
	var arg__4679_93 vm.Value
	var arg__4690_99 vm.Value
	var loop_arg_ids_100 vm.Value
	var arg__4696_102 vm.Value
	var arg__4703_105 vm.Value
	var all_args_106 vm.Value
	var bt_108 vm.Value
	var cur_110 vm.Value
	var v116 vm.Value
	var v118 vm.Value
	var form_129 vm.Value
	var ctx_130 vm.Value
	var headers_131 vm.Value
	var f_132 vm.Value
	var arg_syms_133 vm.Value
	var arg__4750_144 vm.Value
	var arg__4761_150 vm.Value
	var arg_ids_151 vm.Value
	var cur_153 vm.Value
	var arg__4772_156 vm.Value
	var arg__4781_160 vm.Value
	var v161 vm.Value
	var form_134 vm.Value
	var ctx_135 vm.Value
	var headers_136 vm.Value
	var f_137 vm.Value
	var arg_syms_138 vm.Value
	var v167 vm.Value
	var v169 vm.Value
	var form_170 vm.Value
	var ctx_171 vm.Value
	var headers_172 vm.Value
	var f_173 vm.Value
	var arg_syms_174 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__4624_4, arg__4630_8, headers_10, f_12, form_13, ctx_14, headers_15, f_16, header_44, arg__4644_46, arg__4650_50, cap_syms_stack_52, or__x_54, form_17, ctx_18, headers_19, f_20, arg__4733_122, arg__4739_126, arg_syms_128, v176, form_177, ctx_178, headers_179, f_180, form_21, ctx_22, and__x_23, headers_24, f_25, v33, form_26, ctx_27, and__x_28, headers_29, f_30, v36, form_37, ctx_38, and__x_39, headers_40, f_41, form_55, ctx_56, headers_57, f_58, header_59, cap_syms_stack_60, or__x_61, form_62, ctx_63, headers_64, f_65, header_66, cap_syms_stack_67, or__x_68, cap_syms_73, form_74, ctx_75, headers_76, f_77, header_78, cap_syms_stack_79, or__x_80, cap_vals_88, arg__4679_93, arg__4690_99, loop_arg_ids_100, arg__4696_102, arg__4703_105, all_args_106, bt_108, cur_110, v116, v118, form_129, ctx_130, headers_131, f_132, arg_syms_133, arg__4750_144, arg__4761_150, arg_ids_151, cur_153, arg__4772_156, arg__4781_160, v161, form_134, ctx_135, headers_136, f_137, arg_syms_138, v167, v169, form_170, ctx_171, headers_172, f_173, arg_syms_174
	arg__4624_4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__4630_8, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	headers_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg__4630_8, vm.Keyword("loop-headers")})
	if callErr != nil {
		return nil, callErr
	}
	f_12, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-fn").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(headers_10) {
		form_21 = arg0
		ctx_22 = arg1
		and__x_23 = headers_10
		headers_24 = headers_10
		f_25 = f_12
		goto b4
	} else {
		form_26 = arg0
		ctx_27 = arg1
		and__x_28 = headers_10
		headers_29 = headers_10
		f_30 = f_12
		goto b5
	}
b1:
	;
	header_44, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "peek").Deref(), []vm.Value{headers_15})
	if callErr != nil {
		return nil, callErr
	}
	arg__4644_46, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_14})
	if callErr != nil {
		return nil, callErr
	}
	arg__4650_50, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_14})
	if callErr != nil {
		return nil, callErr
	}
	cap_syms_stack_52, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg__4650_50, vm.Keyword("loop-capture-syms-stack")})
	if callErr != nil {
		return nil, callErr
	}
	or__x_54, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "peek").Deref(), []vm.Value{cap_syms_stack_52})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(or__x_54) {
		form_55 = form_13
		ctx_56 = ctx_14
		headers_57 = headers_15
		f_58 = f_16
		header_59 = header_44
		cap_syms_stack_60 = cap_syms_stack_52
		or__x_61 = or__x_54
		goto b7
	} else {
		form_62 = form_13
		ctx_63 = ctx_14
		headers_64 = headers_15
		f_65 = f_16
		header_66 = header_44
		cap_syms_stack_67 = cap_syms_stack_52
		or__x_68 = or__x_54
		goto b8
	}
b2:
	;
	arg__4733_122, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_18})
	if callErr != nil {
		return nil, callErr
	}
	arg__4739_126, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_18})
	if callErr != nil {
		return nil, callErr
	}
	arg_syms_128, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg__4739_126, vm.Keyword("fn-arg-syms")})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(arg_syms_128) {
		form_129 = form_17
		ctx_130 = ctx_18
		headers_131 = headers_19
		f_132 = f_20
		arg_syms_133 = arg_syms_128
		goto b10
	} else {
		form_134 = form_17
		ctx_135 = ctx_18
		headers_136 = headers_19
		f_137 = f_20
		arg_syms_138 = arg_syms_128
		goto b11
	}
b3:
	;
	return v176, nil
b4:
	;
	v33, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{headers_24})
	if callErr != nil {
		return nil, callErr
	}
	v36 = v33
	form_37 = form_21
	ctx_38 = ctx_22
	and__x_39 = and__x_23
	headers_40 = headers_24
	f_41 = f_25
	goto b6
b5:
	;
	v36 = and__x_28
	form_37 = form_26
	ctx_38 = ctx_27
	and__x_39 = and__x_28
	headers_40 = headers_29
	f_41 = f_30
	goto b6
b6:
	;
	if vm.IsTruthy(v36) {
		form_13 = form_37
		ctx_14 = ctx_38
		headers_15 = headers_40
		f_16 = f_41
		goto b1
	} else {
		form_17 = form_37
		ctx_18 = ctx_38
		headers_19 = headers_40
		f_20 = f_41
		goto b2
	}
b7:
	;
	cap_syms_73 = or__x_61
	form_74 = form_55
	ctx_75 = ctx_56
	headers_76 = headers_57
	f_77 = f_58
	header_78 = header_59
	cap_syms_stack_79 = cap_syms_stack_60
	or__x_80 = or__x_61
	goto b9
b8:
	;
	cap_syms_73 = vm.NewArrayVector([]vm.Value{})
	form_74 = form_62
	ctx_75 = ctx_63
	headers_76 = headers_64
	f_77 = f_65
	header_78 = header_66
	cap_syms_stack_79 = cap_syms_stack_67
	or__x_80 = or__x_68
	goto b9
b9:
	;
	cap_vals_88, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapv").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "lookup-local").Deref(), []vm.Value{ctx_75, arg0})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), cap_syms_73})
	if callErr != nil {
		return nil, callErr
	}
	arg__4679_93, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{form_74})
	if callErr != nil {
		return nil, callErr
	}
	arg__4690_99, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{form_74})
	if callErr != nil {
		return nil, callErr
	}
	loop_arg_ids_100, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapv").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-form").Deref(), []vm.Value{arg0, ctx_75})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), arg__4690_99})
	if callErr != nil {
		return nil, callErr
	}
	arg__4696_102, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "concat").Deref(), []vm.Value{loop_arg_ids_100, cap_vals_88})
	if callErr != nil {
		return nil, callErr
	}
	arg__4703_105, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "concat").Deref(), []vm.Value{loop_arg_ids_100, cap_vals_88})
	if callErr != nil {
		return nil, callErr
	}
	all_args_106, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg__4703_105})
	if callErr != nil {
		return nil, callErr
	}
	bt_108, callErr = rt.InvokeValue(rt.LookupVar("ir", "new-branch-target").Deref(), []vm.Value{header_78, all_args_106})
	if callErr != nil {
		return nil, callErr
	}
	cur_110, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_75})
	if callErr != nil {
		return nil, callErr
	}
	v116, callErr = rt.InvokeValue(rt.LookupVar("ir", "add-terminator!").Deref(), []vm.Value{f_77, cur_110, vm.Keyword("branch"), vm.NewArrayVector([]vm.Value{}), bt_108})
	if callErr != nil {
		return nil, callErr
	}
	v118, callErr = rt.InvokeValue(rt.LookupVar("ir", "add-pred!").Deref(), []vm.Value{f_77, header_78, cur_110})
	if callErr != nil {
		return nil, callErr
	}
	v176 = rt.LookupVar("ir.build", "TERMINATED").Deref()
	form_177 = form_74
	ctx_178 = ctx_75
	headers_179 = headers_76
	f_180 = f_77
	goto b3
b10:
	;
	arg__4750_144, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{form_129})
	if callErr != nil {
		return nil, callErr
	}
	arg__4761_150, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{form_129})
	if callErr != nil {
		return nil, callErr
	}
	arg_ids_151, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapv").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-form").Deref(), []vm.Value{arg0, ctx_130})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), arg__4761_150})
	if callErr != nil {
		return nil, callErr
	}
	cur_153, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_130})
	if callErr != nil {
		return nil, callErr
	}
	arg__4772_156, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{arg_ids_151})
	if callErr != nil {
		return nil, callErr
	}
	arg__4781_160, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{arg_ids_151})
	if callErr != nil {
		return nil, callErr
	}
	v161, callErr = rt.InvokeValue(rt.LookupVar("ir", "add-terminator!").Deref(), []vm.Value{f_132, cur_153, vm.Keyword("tail-call"), arg_ids_151, arg__4781_160})
	if callErr != nil {
		return nil, callErr
	}
	v169 = rt.LookupVar("ir.build", "TERMINATED").Deref()
	form_170 = form_129
	ctx_171 = ctx_130
	headers_172 = headers_131
	f_173 = f_132
	arg_syms_174 = arg_syms_133
	goto b12
b11:
	;
	v167, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "throw").Deref(), []vm.Value{vm.String("recur outside of loop or function")})
	if callErr != nil {
		return nil, callErr
	}
	v169 = v167
	form_170 = form_134
	ctx_171 = ctx_135
	headers_172 = headers_136
	f_173 = f_137
	arg_syms_174 = arg_syms_138
	goto b12
b12:
	;
	v176 = v169
	form_177 = form_170
	ctx_178 = ctx_171
	headers_179 = headers_172
	f_180 = f_173
	goto b3
}
func build_set_BANG_(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var __8 vm.Value
	var sym_14 vm.Value
	var val_20 vm.Value
	var v_22 vm.Value
	var v38 vm.Value
	var vec__4785_23 vm.Value
	var form_24 vm.Value
	var ctx_25 vm.Value
	var __26 vm.Value
	var sym_27 vm.Value
	var val_28 vm.Value
	var v_29 vm.Value
	var arg__4818_43 vm.Value
	var arg__4825_48 vm.Value
	var v49 vm.Value
	var vec__4785_30 vm.Value
	var form_31 vm.Value
	var ctx_32 vm.Value
	var __33 vm.Value
	var sym_34 vm.Value
	var val_35 vm.Value
	var v_36 vm.Value
	var v53 vm.Value
	var vec__4785_54 vm.Value
	var form_55 vm.Value
	var ctx_56 vm.Value
	var __57 vm.Value
	var sym_58 vm.Value
	var val_59 vm.Value
	var v_60 vm.Value
	var arg__4830_62 vm.Value
	var arg__4839_67 vm.Value
	var var_nid_70 vm.Value
	var val_nid_72 vm.Value
	var arg__4852_74 vm.Value
	var arg__4857_77 vm.Value
	var arg__4864_81 vm.Value
	var arg__4869_84 vm.Value
	var v86 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = __8, sym_14, val_20, v_22, v38, vec__4785_23, form_24, ctx_25, __26, sym_27, val_28, v_29, arg__4818_43, arg__4825_48, v49, vec__4785_30, form_31, ctx_32, __33, sym_34, val_35, v_36, v53, vec__4785_54, form_55, ctx_56, __57, sym_58, val_59, v_60, arg__4830_62, arg__4839_67, var_nid_70, val_nid_72, arg__4852_74, arg__4857_77, arg__4864_81, arg__4869_84, v86
	__8, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(0), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	sym_14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(1), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	val_20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(2), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	v_22, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "resolve").Deref(), []vm.Value{sym_14})
	if callErr != nil {
		return nil, callErr
	}
	v38, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{v_22})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v38) {
		vec__4785_23 = arg0
		form_24 = arg0
		ctx_25 = arg1
		__26 = __8
		sym_27 = sym_14
		val_28 = val_20
		v_29 = v_22
		goto b1
	} else {
		vec__4785_30 = arg0
		form_31 = arg0
		ctx_32 = arg1
		__33 = __8
		sym_34 = sym_14
		val_35 = val_20
		v_36 = v_22
		goto b2
	}
b1:
	;
	arg__4818_43, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("ir/build: set! can't resolve "), sym_27})
	if callErr != nil {
		return nil, callErr
	}
	arg__4825_48, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("ir/build: set! can't resolve "), sym_27})
	if callErr != nil {
		return nil, callErr
	}
	v49, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "throw").Deref(), []vm.Value{arg__4825_48})
	if callErr != nil {
		return nil, callErr
	}
	v53 = v49
	vec__4785_54 = vec__4785_23
	form_55 = form_24
	ctx_56 = ctx_25
	__57 = __26
	sym_58 = sym_27
	val_59 = val_28
	v_60 = v_29
	goto b3
b2:
	;
	v53 = vm.NIL
	vec__4785_54 = vec__4785_30
	form_55 = form_31
	ctx_56 = ctx_32
	__57 = __33
	sym_58 = sym_34
	val_59 = val_35
	v_60 = v_36
	goto b3
b3:
	;
	arg__4830_62, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_56})
	if callErr != nil {
		return nil, callErr
	}
	arg__4839_67, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_56})
	if callErr != nil {
		return nil, callErr
	}
	var_nid_70, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "add-inst!").Deref(), []vm.Value{ctx_56, arg__4839_67, vm.Keyword("const"), vm.NewArrayVector([]vm.Value{}), v_60})
	if callErr != nil {
		return nil, callErr
	}
	val_nid_72, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-form").Deref(), []vm.Value{val_59, ctx_56})
	if callErr != nil {
		return nil, callErr
	}
	arg__4852_74, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_56})
	if callErr != nil {
		return nil, callErr
	}
	arg__4857_77, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{var_nid_70, val_nid_72})
	if callErr != nil {
		return nil, callErr
	}
	arg__4864_81, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_56})
	if callErr != nil {
		return nil, callErr
	}
	arg__4869_84, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{var_nid_70, val_nid_72})
	if callErr != nil {
		return nil, callErr
	}
	v86, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "add-inst!").Deref(), []vm.Value{ctx_56, arg__4864_81, vm.Keyword("set-var"), arg__4869_84, vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	return v86, nil
}
func build_single_fn_STAR_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value, arg3 vm.Value) (vm.Value, error) {
	var expanded_6 vm.Value
	var name_sym_7 vm.Value
	var args_vec_8 vm.Value
	var body_forms_9 vm.Value
	var ctx_10 vm.Value
	var expanded_11 vm.Value
	var v19 vm.Value
	var name_sym_12 vm.Value
	var args_vec_13 vm.Value
	var body_forms_14 vm.Value
	var ctx_15 vm.Value
	var expanded_16 vm.Value
	var flat_args_22 vm.Value
	var name_sym_23 vm.Value
	var args_vec_24 vm.Value
	var body_forms_25 vm.Value
	var ctx_26 vm.Value
	var expanded_27 vm.Value
	var flat_args_28 vm.Value
	var name_sym_29 vm.Value
	var args_vec_30 vm.Value
	var body_forms_31 vm.Value
	var ctx_32 vm.Value
	var expanded_33 vm.Value
	var v42 vm.Value
	var flat_args_34 vm.Value
	var name_sym_35 vm.Value
	var args_vec_36 vm.Value
	var body_forms_37 vm.Value
	var ctx_38 vm.Value
	var expanded_39 vm.Value
	var flat_body_45 vm.Value
	var flat_args_46 vm.Value
	var name_sym_47 vm.Value
	var args_vec_48 vm.Value
	var body_forms_49 vm.Value
	var ctx_50 vm.Value
	var expanded_51 vm.Value
	var flat_body_52 vm.Value
	var flat_args_53 vm.Value
	var name_sym_54 vm.Value
	var args_vec_55 vm.Value
	var body_forms_56 vm.Value
	var ctx_57 vm.Value
	var expanded_58 vm.Value
	var v68 vm.Value
	var flat_body_59 vm.Value
	var flat_args_60 vm.Value
	var name_sym_61 vm.Value
	var args_vec_62 vm.Value
	var body_forms_63 vm.Value
	var ctx_64 vm.Value
	var expanded_65 vm.Value
	var variadic_QMARK__72 vm.Value
	var flat_body_73 vm.Value
	var flat_args_74 vm.Value
	var name_sym_75 vm.Value
	var args_vec_76 vm.Value
	var body_forms_77 vm.Value
	var ctx_78 vm.Value
	var expanded_79 vm.Value
	var arg__4883_81 vm.Value
	var arg__4887_84 vm.Value
	var arg_set_85 vm.Value
	var arg__4894_89 vm.Value
	var arg__4902_94 vm.Value
	var frees_95 vm.Value
	var arg__4920_104 vm.Value
	var arg__4938_114 vm.Value
	var arg__4939_115 vm.Value
	var arg__4957_125 vm.Value
	var arg__4975_135 vm.Value
	var arg__4976_136 vm.Value
	var captures_137 vm.Value
	var template_139 vm.Value
	var v141 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = expanded_6, name_sym_7, args_vec_8, body_forms_9, ctx_10, expanded_11, v19, name_sym_12, args_vec_13, body_forms_14, ctx_15, expanded_16, flat_args_22, name_sym_23, args_vec_24, body_forms_25, ctx_26, expanded_27, flat_args_28, name_sym_29, args_vec_30, body_forms_31, ctx_32, expanded_33, v42, flat_args_34, name_sym_35, args_vec_36, body_forms_37, ctx_38, expanded_39, flat_body_45, flat_args_46, name_sym_47, args_vec_48, body_forms_49, ctx_50, expanded_51, flat_body_52, flat_args_53, name_sym_54, args_vec_55, body_forms_56, ctx_57, expanded_58, v68, flat_body_59, flat_args_60, name_sym_61, args_vec_62, body_forms_63, ctx_64, expanded_65, variadic_QMARK__72, flat_body_73, flat_args_74, name_sym_75, args_vec_76, body_forms_77, ctx_78, expanded_79, arg__4883_81, arg__4887_84, arg_set_85, arg__4894_89, arg__4902_94, frees_95, arg__4920_104, arg__4938_114, arg__4939_115, arg__4957_125, arg__4975_135, arg__4976_136, captures_137, template_139, v141
	expanded_6, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "expand-fn-args").Deref(), []vm.Value{arg1, arg2})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(expanded_6) {
		name_sym_7 = arg0
		args_vec_8 = arg1
		body_forms_9 = arg2
		ctx_10 = arg3
		expanded_11 = expanded_6
		goto b1
	} else {
		name_sym_12 = arg0
		args_vec_13 = arg1
		body_forms_14 = arg2
		ctx_15 = arg3
		expanded_16 = expanded_6
		goto b2
	}
b1:
	;
	v19, callErr = rt.InvokeValue(vm.Keyword("flat-args"), []vm.Value{expanded_11})
	if callErr != nil {
		return nil, callErr
	}
	flat_args_22 = v19
	name_sym_23 = name_sym_7
	args_vec_24 = args_vec_8
	body_forms_25 = body_forms_9
	ctx_26 = ctx_10
	expanded_27 = expanded_11
	goto b3
b2:
	;
	flat_args_22 = args_vec_13
	name_sym_23 = name_sym_12
	args_vec_24 = args_vec_13
	body_forms_25 = body_forms_14
	ctx_26 = ctx_15
	expanded_27 = expanded_16
	goto b3
b3:
	;
	if vm.IsTruthy(expanded_27) {
		flat_args_28 = flat_args_22
		name_sym_29 = name_sym_23
		args_vec_30 = args_vec_24
		body_forms_31 = body_forms_25
		ctx_32 = ctx_26
		expanded_33 = expanded_27
		goto b4
	} else {
		flat_args_34 = flat_args_22
		name_sym_35 = name_sym_23
		args_vec_36 = args_vec_24
		body_forms_37 = body_forms_25
		ctx_38 = ctx_26
		expanded_39 = expanded_27
		goto b5
	}
b4:
	;
	v42, callErr = rt.InvokeValue(vm.Keyword("body"), []vm.Value{expanded_33})
	if callErr != nil {
		return nil, callErr
	}
	flat_body_45 = v42
	flat_args_46 = flat_args_28
	name_sym_47 = name_sym_29
	args_vec_48 = args_vec_30
	body_forms_49 = body_forms_31
	ctx_50 = ctx_32
	expanded_51 = expanded_33
	goto b6
b5:
	;
	flat_body_45 = body_forms_37
	flat_args_46 = flat_args_34
	name_sym_47 = name_sym_35
	args_vec_48 = args_vec_36
	body_forms_49 = body_forms_37
	ctx_50 = ctx_38
	expanded_51 = expanded_39
	goto b6
b6:
	;
	if vm.IsTruthy(expanded_51) {
		flat_body_52 = flat_body_45
		flat_args_53 = flat_args_46
		name_sym_54 = name_sym_47
		args_vec_55 = args_vec_48
		body_forms_56 = body_forms_49
		ctx_57 = ctx_50
		expanded_58 = expanded_51
		goto b7
	} else {
		flat_body_59 = flat_body_45
		flat_args_60 = flat_args_46
		name_sym_61 = name_sym_47
		args_vec_62 = args_vec_48
		body_forms_63 = body_forms_49
		ctx_64 = ctx_50
		expanded_65 = expanded_51
		goto b8
	}
b7:
	;
	v68, callErr = rt.InvokeValue(vm.Keyword("variadic?"), []vm.Value{expanded_58})
	if callErr != nil {
		return nil, callErr
	}
	variadic_QMARK__72 = v68
	flat_body_73 = flat_body_52
	flat_args_74 = flat_args_53
	name_sym_75 = name_sym_54
	args_vec_76 = args_vec_55
	body_forms_77 = body_forms_56
	ctx_78 = ctx_57
	expanded_79 = expanded_58
	goto b9
b8:
	;
	variadic_QMARK__72 = vm.Boolean(false)
	flat_body_73 = flat_body_59
	flat_args_74 = flat_args_60
	name_sym_75 = name_sym_61
	args_vec_76 = args_vec_62
	body_forms_77 = body_forms_63
	ctx_78 = ctx_64
	expanded_79 = expanded_65
	goto b9
b9:
	;
	arg__4883_81, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__4887_84, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg_set_85, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__4887_84, flat_args_74})
	if callErr != nil {
		return nil, callErr
	}
	arg__4894_89, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "cons").Deref(), []vm.Value{vm.Symbol("do"), flat_body_73})
	if callErr != nil {
		return nil, callErr
	}
	arg__4902_94, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "cons").Deref(), []vm.Value{vm.Symbol("do"), flat_body_73})
	if callErr != nil {
		return nil, callErr
	}
	frees_95, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "free-vars").Deref(), []vm.Value{arg__4902_94, arg_set_85})
	if callErr != nil {
		return nil, callErr
	}
	arg__4920_104, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "filter").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "lookup-local").Deref(), []vm.Value{ctx_78, arg0})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), frees_95})
	if callErr != nil {
		return nil, callErr
	}
	arg__4938_114, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "filter").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "lookup-local").Deref(), []vm.Value{ctx_78, arg0})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), frees_95})
	if callErr != nil {
		return nil, callErr
	}
	arg__4939_115, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "sort-by").Deref(), []vm.Value{rt.LookupVar("clojure.core", "str").Deref(), arg__4938_114})
	if callErr != nil {
		return nil, callErr
	}
	arg__4957_125, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "filter").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "lookup-local").Deref(), []vm.Value{ctx_78, arg0})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), frees_95})
	if callErr != nil {
		return nil, callErr
	}
	arg__4975_135, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "filter").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "lookup-local").Deref(), []vm.Value{ctx_78, arg0})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), frees_95})
	if callErr != nil {
		return nil, callErr
	}
	arg__4976_136, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "sort-by").Deref(), []vm.Value{rt.LookupVar("clojure.core", "str").Deref(), arg__4975_135})
	if callErr != nil {
		return nil, callErr
	}
	captures_137, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg__4976_136})
	if callErr != nil {
		return nil, callErr
	}
	template_139, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-inner-fn-template").Deref(), []vm.Value{name_sym_75, flat_args_74, flat_body_73, captures_137, variadic_QMARK__72})
	if callErr != nil {
		return nil, callErr
	}
	v141, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "emit-template-closure").Deref(), []vm.Value{template_139, captures_137, ctx_78})
	if callErr != nil {
		return nil, callErr
	}
	return v141, nil
}
func build_symbol(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var local_3 vm.Value
	var sym_4 vm.Value
	var ctx_5 vm.Value
	var local_6 vm.Value
	var sym_7 vm.Value
	var ctx_8 vm.Value
	var local_9 vm.Value
	var v_13 vm.Value
	var v23 vm.Value
	var v53 vm.Value
	var sym_54 vm.Value
	var ctx_55 vm.Value
	var local_56 vm.Value
	var sym_14 vm.Value
	var ctx_15 vm.Value
	var local_16 vm.Value
	var v_17 vm.Value
	var arg__5011_28 vm.Value
	var arg__5018_33 vm.Value
	var v34 vm.Value
	var sym_18 vm.Value
	var ctx_19 vm.Value
	var local_20 vm.Value
	var v_21 vm.Value
	var arg__5023_37 vm.Value
	var arg__5032_42 vm.Value
	var v45 vm.Value
	var v47 vm.Value
	var sym_48 vm.Value
	var ctx_49 vm.Value
	var local_50 vm.Value
	var v_51 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = local_3, sym_4, ctx_5, local_6, sym_7, ctx_8, local_9, v_13, v23, v53, sym_54, ctx_55, local_56, sym_14, ctx_15, local_16, v_17, arg__5011_28, arg__5018_33, v34, sym_18, ctx_19, local_20, v_21, arg__5023_37, arg__5032_42, v45, v47, sym_48, ctx_49, local_50, v_51
	local_3, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "lookup-local").Deref(), []vm.Value{arg1, arg0})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(local_3) {
		sym_4 = arg0
		ctx_5 = arg1
		local_6 = local_3
		goto b1
	} else {
		sym_7 = arg0
		ctx_8 = arg1
		local_9 = local_3
		goto b2
	}
b1:
	;
	v53 = local_6
	sym_54 = sym_4
	ctx_55 = ctx_5
	local_56 = local_6
	goto b3
b2:
	;
	v_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "resolve").Deref(), []vm.Value{sym_7})
	if callErr != nil {
		return nil, callErr
	}
	v23, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{v_13})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v23) {
		sym_14 = sym_7
		ctx_15 = ctx_8
		local_16 = local_9
		v_17 = v_13
		goto b4
	} else {
		sym_18 = sym_7
		ctx_19 = ctx_8
		local_20 = local_9
		v_21 = v_13
		goto b5
	}
b3:
	;
	return v53, nil
b4:
	;
	arg__5011_28, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("ir/build: unresolved symbol "), sym_14})
	if callErr != nil {
		return nil, callErr
	}
	arg__5018_33, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("ir/build: unresolved symbol "), sym_14})
	if callErr != nil {
		return nil, callErr
	}
	v34, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "throw").Deref(), []vm.Value{arg__5018_33})
	if callErr != nil {
		return nil, callErr
	}
	v47 = v34
	sym_48 = sym_14
	ctx_49 = ctx_15
	local_50 = local_16
	v_51 = v_17
	goto b6
b5:
	;
	arg__5023_37, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_19})
	if callErr != nil {
		return nil, callErr
	}
	arg__5032_42, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_19})
	if callErr != nil {
		return nil, callErr
	}
	v45, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "add-inst!").Deref(), []vm.Value{ctx_19, arg__5032_42, vm.Keyword("load-var"), vm.NewArrayVector([]vm.Value{}), v_21})
	if callErr != nil {
		return nil, callErr
	}
	v47 = v45
	sym_48 = sym_18
	ctx_49 = ctx_19
	local_50 = local_20
	v_51 = v_21
	goto b6
b6:
	;
	v53 = v47
	sym_54 = sym_48
	ctx_55 = ctx_49
	local_56 = local_50
	goto b3
}
func build_var(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var __8 vm.Value
	var sym_14 vm.Value
	var v_16 vm.Value
	var v30 vm.Value
	var form_17 vm.Value
	var vec__5036_18 vm.Value
	var ctx_19 vm.Value
	var __20 vm.Value
	var sym_21 vm.Value
	var v_22 vm.Value
	var arg__5062_35 vm.Value
	var arg__5069_40 vm.Value
	var v41 vm.Value
	var form_23 vm.Value
	var vec__5036_24 vm.Value
	var ctx_25 vm.Value
	var __26 vm.Value
	var sym_27 vm.Value
	var v_28 vm.Value
	var arg__5074_44 vm.Value
	var arg__5083_49 vm.Value
	var v52 vm.Value
	var v54 vm.Value
	var form_55 vm.Value
	var vec__5036_56 vm.Value
	var ctx_57 vm.Value
	var __58 vm.Value
	var sym_59 vm.Value
	var v_60 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = __8, sym_14, v_16, v30, form_17, vec__5036_18, ctx_19, __20, sym_21, v_22, arg__5062_35, arg__5069_40, v41, form_23, vec__5036_24, ctx_25, __26, sym_27, v_28, arg__5074_44, arg__5083_49, v52, v54, form_55, vec__5036_56, ctx_57, __58, sym_59, v_60
	__8, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(0), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	sym_14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(1), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	v_16, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "resolve").Deref(), []vm.Value{sym_14})
	if callErr != nil {
		return nil, callErr
	}
	v30, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{v_16})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v30) {
		form_17 = arg0
		vec__5036_18 = arg0
		ctx_19 = arg1
		__20 = __8
		sym_21 = sym_14
		v_22 = v_16
		goto b1
	} else {
		form_23 = arg0
		vec__5036_24 = arg0
		ctx_25 = arg1
		__26 = __8
		sym_27 = sym_14
		v_28 = v_16
		goto b2
	}
b1:
	;
	arg__5062_35, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("ir/build: can't resolve var "), sym_21})
	if callErr != nil {
		return nil, callErr
	}
	arg__5069_40, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("ir/build: can't resolve var "), sym_21})
	if callErr != nil {
		return nil, callErr
	}
	v41, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "throw").Deref(), []vm.Value{arg__5069_40})
	if callErr != nil {
		return nil, callErr
	}
	v54 = v41
	form_55 = form_17
	vec__5036_56 = vec__5036_18
	ctx_57 = ctx_19
	__58 = __20
	sym_59 = sym_21
	v_60 = v_22
	goto b3
b2:
	;
	arg__5074_44, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_25})
	if callErr != nil {
		return nil, callErr
	}
	arg__5083_49, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_25})
	if callErr != nil {
		return nil, callErr
	}
	v52, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "add-inst!").Deref(), []vm.Value{ctx_25, arg__5083_49, vm.Keyword("const"), vm.NewArrayVector([]vm.Value{}), v_28})
	if callErr != nil {
		return nil, callErr
	}
	v54 = v52
	form_55 = form_23
	vec__5036_56 = vec__5036_24
	ctx_57 = ctx_25
	__58 = __26
	sym_59 = sym_27
	v_60 = v_28
	goto b3
b3:
	;
	return v54, nil
}
func build_vector(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__5090_8 vm.Value
	var arg__5095_11 vm.Value
	var v12 vm.Value
	var form_3 vm.Value
	var ctx_4 vm.Value
	var arg__5100_15 vm.Value
	var arg__5109_20 vm.Value
	var v23 vm.Value
	var form_5 vm.Value
	var ctx_6 vm.Value
	var arg__5117_26 vm.Value
	var arg__5123_32 vm.Value
	var arg__5129_35 vm.Value
	var arg__5135_41 vm.Value
	var fn_id_42 vm.Value
	var arg__5140_44 vm.Value
	var arg__5147_47 vm.Value
	var v48 vm.Value
	var v50 vm.Value
	var form_51 vm.Value
	var ctx_52 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__5090_8, arg__5095_11, v12, form_3, ctx_4, arg__5100_15, arg__5109_20, v23, form_5, ctx_6, arg__5117_26, arg__5123_32, arg__5129_35, arg__5135_41, fn_id_42, arg__5140_44, arg__5147_47, v48, v50, form_51, ctx_52
	arg__5090_8, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__5095_11, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "zero?").Deref(), []vm.Value{arg__5095_11})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v12) {
		form_3 = arg0
		ctx_4 = arg1
		goto b1
	} else {
		form_5 = arg0
		ctx_6 = arg1
		goto b2
	}
b1:
	;
	arg__5100_15, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_4})
	if callErr != nil {
		return nil, callErr
	}
	arg__5109_20, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_4})
	if callErr != nil {
		return nil, callErr
	}
	v23, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "add-inst!").Deref(), []vm.Value{ctx_4, arg__5109_20, vm.Keyword("const"), vm.NewArrayVector([]vm.Value{}), form_3})
	if callErr != nil {
		return nil, callErr
	}
	v50 = v23
	form_51 = form_3
	ctx_52 = ctx_4
	goto b3
b2:
	;
	arg__5117_26, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_6})
	if callErr != nil {
		return nil, callErr
	}
	arg__5123_32, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "resolve").Deref(), []vm.Value{vm.Symbol("vector")})
	if callErr != nil {
		return nil, callErr
	}
	arg__5129_35, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_6})
	if callErr != nil {
		return nil, callErr
	}
	arg__5135_41, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "resolve").Deref(), []vm.Value{vm.Symbol("vector")})
	if callErr != nil {
		return nil, callErr
	}
	fn_id_42, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "add-inst!").Deref(), []vm.Value{ctx_6, arg__5129_35, vm.Keyword("load-var"), vm.NewArrayVector([]vm.Value{}), arg__5135_41})
	if callErr != nil {
		return nil, callErr
	}
	arg__5140_44, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{form_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__5147_47, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{form_5})
	if callErr != nil {
		return nil, callErr
	}
	v48, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-call-with-head").Deref(), []vm.Value{fn_id_42, arg__5147_47, ctx_6})
	if callErr != nil {
		return nil, callErr
	}
	v50 = v48
	form_51 = form_5
	ctx_52 = ctx_6
	goto b3
b3:
	;
	return v50, nil
}
func captures_of(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var v8 vm.Value
	var form_3 vm.Value
	var known_locals_4 vm.Value
	var v14 vm.Value
	var form_5 vm.Value
	var known_locals_6 vm.Value
	var or__x_31 vm.Value
	var v600 vm.Value
	var form_601 vm.Value
	var known_locals_602 vm.Value
	var form_10 vm.Value
	var known_locals_11 vm.Value
	var v17 vm.Value
	var form_12 vm.Value
	var known_locals_13 vm.Value
	var v20 vm.Value
	var v22 vm.Value
	var form_23 vm.Value
	var known_locals_24 vm.Value
	var form_26 vm.Value
	var known_locals_27 vm.Value
	var head_49 vm.Value
	var v57 bool
	var form_28 vm.Value
	var known_locals_29 vm.Value
	var v552 vm.Value
	var v596 vm.Value
	var form_597 vm.Value
	var known_locals_598 vm.Value
	var form_32 vm.Value
	var known_locals_33 vm.Value
	var or__x_34 vm.Value
	var form_35 vm.Value
	var known_locals_36 vm.Value
	var or__x_37 vm.Value
	var v41 vm.Value
	var v43 vm.Value
	var form_44 vm.Value
	var known_locals_45 vm.Value
	var or__x_46 vm.Value
	var form_50 vm.Value
	var known_locals_51 vm.Value
	var head_52 vm.Value
	var v60 vm.Value
	var form_53 vm.Value
	var known_locals_54 vm.Value
	var head_55 vm.Value
	var v69 bool
	var v542 vm.Value
	var form_543 vm.Value
	var known_locals_544 vm.Value
	var head_545 vm.Value
	var form_62 vm.Value
	var known_locals_63 vm.Value
	var head_64 vm.Value
	var v72 vm.Value
	var form_65 vm.Value
	var known_locals_66 vm.Value
	var head_67 vm.Value
	var v81 bool
	var v537 vm.Value
	var form_538 vm.Value
	var known_locals_539 vm.Value
	var head_540 vm.Value
	var form_74 vm.Value
	var known_locals_75 vm.Value
	var head_76 vm.Value
	var __88 vm.Value
	var _sym_94 vm.Value
	var val_100 vm.Value
	var v102 vm.Value
	var form_77 vm.Value
	var known_locals_78 vm.Value
	var head_79 vm.Value
	var v111 bool
	var v532 vm.Value
	var form_533 vm.Value
	var known_locals_534 vm.Value
	var head_535 vm.Value
	var form_104 vm.Value
	var known_locals_105 vm.Value
	var head_106 vm.Value
	var __118 vm.Value
	var maybe_name_124 vm.Value
	var raw_rest_128 vm.Value
	var has_name_QMARK__130 vm.Value
	var form_107 vm.Value
	var known_locals_108 vm.Value
	var head_109 vm.Value
	var or__x_345 bool
	var v527 vm.Value
	var form_528 vm.Value
	var known_locals_529 vm.Value
	var head_530 vm.Value
	var vec__5150_131 vm.Value
	var form_132 vm.Value
	var known_locals_133 vm.Value
	var head_134 vm.Value
	var __135 vm.Value
	var maybe_name_136 vm.Value
	var raw_rest_137 vm.Value
	var has_name_QMARK__138 vm.Value
	var vec__5150_139 vm.Value
	var form_140 vm.Value
	var known_locals_141 vm.Value
	var head_142 vm.Value
	var __143 vm.Value
	var maybe_name_144 vm.Value
	var raw_rest_145 vm.Value
	var has_name_QMARK__146 vm.Value
	var name_sym_151 vm.Value
	var vec__5150_152 vm.Value
	var form_153 vm.Value
	var known_locals_154 vm.Value
	var head_155 vm.Value
	var __156 vm.Value
	var maybe_name_157 vm.Value
	var raw_rest_158 vm.Value
	var has_name_QMARK__159 vm.Value
	var name_sym_160 vm.Value
	var vec__5150_161 vm.Value
	var form_162 vm.Value
	var known_locals_163 vm.Value
	var head_164 vm.Value
	var __165 vm.Value
	var maybe_name_166 vm.Value
	var raw_rest_167 vm.Value
	var has_name_QMARK__168 vm.Value
	var name_sym_169 vm.Value
	var vec__5150_170 vm.Value
	var form_171 vm.Value
	var known_locals_172 vm.Value
	var head_173 vm.Value
	var __174 vm.Value
	var maybe_name_175 vm.Value
	var raw_rest_176 vm.Value
	var has_name_QMARK__177 vm.Value
	var v181 vm.Value
	var rest_forms_183 vm.Value
	var name_sym_184 vm.Value
	var vec__5150_185 vm.Value
	var form_186 vm.Value
	var known_locals_187 vm.Value
	var head_188 vm.Value
	var __189 vm.Value
	var maybe_name_190 vm.Value
	var raw_rest_191 vm.Value
	var has_name_QMARK__192 vm.Value
	var and__x_194 vm.Value
	var rest_forms_195 vm.Value
	var name_sym_196 vm.Value
	var vec__5150_197 vm.Value
	var form_198 vm.Value
	var known_locals_199 vm.Value
	var head_200 vm.Value
	var __201 vm.Value
	var maybe_name_202 vm.Value
	var raw_rest_203 vm.Value
	var has_name_QMARK__204 vm.Value
	var and__x_205 vm.Value
	var arg__5243_219 vm.Value
	var arg__5248_222 vm.Value
	var v223 vm.Value
	var rest_forms_206 vm.Value
	var name_sym_207 vm.Value
	var vec__5150_208 vm.Value
	var form_209 vm.Value
	var known_locals_210 vm.Value
	var head_211 vm.Value
	var __212 vm.Value
	var maybe_name_213 vm.Value
	var raw_rest_214 vm.Value
	var has_name_QMARK__215 vm.Value
	var and__x_216 vm.Value
	var multi_QMARK__226 vm.Value
	var rest_forms_227 vm.Value
	var name_sym_228 vm.Value
	var vec__5150_229 vm.Value
	var form_230 vm.Value
	var known_locals_231 vm.Value
	var head_232 vm.Value
	var __233 vm.Value
	var maybe_name_234 vm.Value
	var raw_rest_235 vm.Value
	var has_name_QMARK__236 vm.Value
	var and__x_237 vm.Value
	var multi_QMARK__241 vm.Value
	var rest_forms_242 vm.Value
	var name_sym_243 vm.Value
	var vec__5150_244 vm.Value
	var form_245 vm.Value
	var known_locals_246 vm.Value
	var head_247 vm.Value
	var __248 vm.Value
	var maybe_name_249 vm.Value
	var raw_rest_250 vm.Value
	var has_name_QMARK__251 vm.Value
	var shadow_for_252 vm.Value
	var arg__5276_267 vm.Value
	var arg__5346_277 vm.Value
	var arg__5349_280 vm.Value
	var arg__5419_290 vm.Value
	var v291 vm.Value
	var multi_QMARK__253 vm.Value
	var rest_forms_254 vm.Value
	var name_sym_255 vm.Value
	var vec__5150_256 vm.Value
	var form_257 vm.Value
	var known_locals_258 vm.Value
	var head_259 vm.Value
	var __260 vm.Value
	var maybe_name_261 vm.Value
	var raw_rest_262 vm.Value
	var has_name_QMARK__263 vm.Value
	var shadow_for_264 vm.Value
	var args_vec_294 vm.Value
	var body_296 vm.Value
	var arg__5430_297 vm.Value
	var arg__5436_299 vm.Value
	var inner_known_300 vm.Value
	var arg__5438_302 vm.Value
	var arg__5454_310 vm.Value
	var arg__5457_313 vm.Value
	var arg__5473_321 vm.Value
	var v322 vm.Value
	var v324 vm.Value
	var multi_QMARK__325 vm.Value
	var rest_forms_326 vm.Value
	var name_sym_327 vm.Value
	var vec__5150_328 vm.Value
	var form_329 vm.Value
	var known_locals_330 vm.Value
	var head_331 vm.Value
	var __332 vm.Value
	var maybe_name_333 vm.Value
	var raw_rest_334 vm.Value
	var has_name_QMARK__335 vm.Value
	var shadow_for_336 vm.Value
	var form_338 vm.Value
	var known_locals_339 vm.Value
	var head_340 vm.Value
	var __406 vm.Value
	var bindings_412 vm.Value
	var body_416 vm.Value
	var pairs_420 vm.Value
	var arg__5571_426 vm.Value
	var arg__5573_428 vm.Value
	var arg__5574_429 vm.Value
	var arg__5642_436 vm.Value
	var arg__5644_438 vm.Value
	var arg__5645_439 vm.Value
	var vec__5152_440 vm.Value
	var used_446 vm.Value
	var let_bound_set_452 vm.Value
	var new_locals_454 vm.Value
	var arg__5667_456 vm.Value
	var arg__5683_464 vm.Value
	var arg__5686_467 vm.Value
	var arg__5702_475 vm.Value
	var body_captures_476 vm.Value
	var arg__5709_478 vm.Value
	var arg__5717_481 vm.Value
	var v482 vm.Value
	var form_341 vm.Value
	var known_locals_342 vm.Value
	var head_343 vm.Value
	var v522 vm.Value
	var form_523 vm.Value
	var known_locals_524 vm.Value
	var head_525 vm.Value
	var form_346 vm.Value
	var known_locals_347 vm.Value
	var head_348 vm.Value
	var or__x_349 bool
	var form_350 vm.Value
	var known_locals_351 vm.Value
	var head_352 vm.Value
	var or__x_353 bool
	var or__x_357 bool
	var v395 bool
	var form_396 vm.Value
	var known_locals_397 vm.Value
	var head_398 vm.Value
	var or__x_399 vm.Value
	var form_358 vm.Value
	var known_locals_359 vm.Value
	var head_360 vm.Value
	var or__x_361 bool
	var form_362 vm.Value
	var known_locals_363 vm.Value
	var head_364 vm.Value
	var or__x_365 bool
	var or__x_369 bool
	var v389 bool
	var form_390 vm.Value
	var known_locals_391 vm.Value
	var head_392 vm.Value
	var or__x_393 vm.Value
	var form_370 vm.Value
	var known_locals_371 vm.Value
	var head_372 vm.Value
	var or__x_373 bool
	var form_374 vm.Value
	var known_locals_375 vm.Value
	var head_376 vm.Value
	var or__x_377 bool
	var v381 bool
	var v383 bool
	var form_384 vm.Value
	var known_locals_385 vm.Value
	var head_386 vm.Value
	var or__x_387 vm.Value
	var form_484 vm.Value
	var known_locals_485 vm.Value
	var head_486 vm.Value
	var arg__5719_493 vm.Value
	var arg__5735_501 vm.Value
	var arg__5738_504 vm.Value
	var arg__5754_512 vm.Value
	var v513 vm.Value
	var form_487 vm.Value
	var known_locals_488 vm.Value
	var head_489 vm.Value
	var v517 vm.Value
	var form_518 vm.Value
	var known_locals_519 vm.Value
	var head_520 vm.Value
	var form_547 vm.Value
	var known_locals_548 vm.Value
	var arg__5759_555 vm.Value
	var arg__5775_563 vm.Value
	var arg__5778_566 vm.Value
	var arg__5794_574 vm.Value
	var v575 vm.Value
	var form_549 vm.Value
	var known_locals_550 vm.Value
	var v592 vm.Value
	var form_593 vm.Value
	var known_locals_594 vm.Value
	var form_577 vm.Value
	var known_locals_578 vm.Value
	var v584 vm.Value
	var form_579 vm.Value
	var known_locals_580 vm.Value
	var v588 vm.Value
	var form_589 vm.Value
	var known_locals_590 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = v8, form_3, known_locals_4, v14, form_5, known_locals_6, or__x_31, v600, form_601, known_locals_602, form_10, known_locals_11, v17, form_12, known_locals_13, v20, v22, form_23, known_locals_24, form_26, known_locals_27, head_49, v57, form_28, known_locals_29, v552, v596, form_597, known_locals_598, form_32, known_locals_33, or__x_34, form_35, known_locals_36, or__x_37, v41, v43, form_44, known_locals_45, or__x_46, form_50, known_locals_51, head_52, v60, form_53, known_locals_54, head_55, v69, v542, form_543, known_locals_544, head_545, form_62, known_locals_63, head_64, v72, form_65, known_locals_66, head_67, v81, v537, form_538, known_locals_539, head_540, form_74, known_locals_75, head_76, __88, _sym_94, val_100, v102, form_77, known_locals_78, head_79, v111, v532, form_533, known_locals_534, head_535, form_104, known_locals_105, head_106, __118, maybe_name_124, raw_rest_128, has_name_QMARK__130, form_107, known_locals_108, head_109, or__x_345, v527, form_528, known_locals_529, head_530, vec__5150_131, form_132, known_locals_133, head_134, __135, maybe_name_136, raw_rest_137, has_name_QMARK__138, vec__5150_139, form_140, known_locals_141, head_142, __143, maybe_name_144, raw_rest_145, has_name_QMARK__146, name_sym_151, vec__5150_152, form_153, known_locals_154, head_155, __156, maybe_name_157, raw_rest_158, has_name_QMARK__159, name_sym_160, vec__5150_161, form_162, known_locals_163, head_164, __165, maybe_name_166, raw_rest_167, has_name_QMARK__168, name_sym_169, vec__5150_170, form_171, known_locals_172, head_173, __174, maybe_name_175, raw_rest_176, has_name_QMARK__177, v181, rest_forms_183, name_sym_184, vec__5150_185, form_186, known_locals_187, head_188, __189, maybe_name_190, raw_rest_191, has_name_QMARK__192, and__x_194, rest_forms_195, name_sym_196, vec__5150_197, form_198, known_locals_199, head_200, __201, maybe_name_202, raw_rest_203, has_name_QMARK__204, and__x_205, arg__5243_219, arg__5248_222, v223, rest_forms_206, name_sym_207, vec__5150_208, form_209, known_locals_210, head_211, __212, maybe_name_213, raw_rest_214, has_name_QMARK__215, and__x_216, multi_QMARK__226, rest_forms_227, name_sym_228, vec__5150_229, form_230, known_locals_231, head_232, __233, maybe_name_234, raw_rest_235, has_name_QMARK__236, and__x_237, multi_QMARK__241, rest_forms_242, name_sym_243, vec__5150_244, form_245, known_locals_246, head_247, __248, maybe_name_249, raw_rest_250, has_name_QMARK__251, shadow_for_252, arg__5276_267, arg__5346_277, arg__5349_280, arg__5419_290, v291, multi_QMARK__253, rest_forms_254, name_sym_255, vec__5150_256, form_257, known_locals_258, head_259, __260, maybe_name_261, raw_rest_262, has_name_QMARK__263, shadow_for_264, args_vec_294, body_296, arg__5430_297, arg__5436_299, inner_known_300, arg__5438_302, arg__5454_310, arg__5457_313, arg__5473_321, v322, v324, multi_QMARK__325, rest_forms_326, name_sym_327, vec__5150_328, form_329, known_locals_330, head_331, __332, maybe_name_333, raw_rest_334, has_name_QMARK__335, shadow_for_336, form_338, known_locals_339, head_340, __406, bindings_412, body_416, pairs_420, arg__5571_426, arg__5573_428, arg__5574_429, arg__5642_436, arg__5644_438, arg__5645_439, vec__5152_440, used_446, let_bound_set_452, new_locals_454, arg__5667_456, arg__5683_464, arg__5686_467, arg__5702_475, body_captures_476, arg__5709_478, arg__5717_481, v482, form_341, known_locals_342, head_343, v522, form_523, known_locals_524, head_525, form_346, known_locals_347, head_348, or__x_349, form_350, known_locals_351, head_352, or__x_353, or__x_357, v395, form_396, known_locals_397, head_398, or__x_399, form_358, known_locals_359, head_360, or__x_361, form_362, known_locals_363, head_364, or__x_365, or__x_369, v389, form_390, known_locals_391, head_392, or__x_393, form_370, known_locals_371, head_372, or__x_373, form_374, known_locals_375, head_376, or__x_377, v381, v383, form_384, known_locals_385, head_386, or__x_387, form_484, known_locals_485, head_486, arg__5719_493, arg__5735_501, arg__5738_504, arg__5754_512, v513, form_487, known_locals_488, head_489, v517, form_518, known_locals_519, head_520, form_547, known_locals_548, arg__5759_555, arg__5775_563, arg__5778_566, arg__5794_574, v575, form_549, known_locals_550, v592, form_593, known_locals_594, form_577, known_locals_578, v584, form_579, known_locals_580, v588, form_589, known_locals_590
	v8, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "symbol?").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v8) {
		form_3 = arg0
		known_locals_4 = arg1
		goto b1
	} else {
		form_5 = arg0
		known_locals_6 = arg1
		goto b2
	}
b1:
	;
	v14, callErr = rt.InvokeValue(known_locals_4, []vm.Value{form_3})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v14) {
		form_10 = form_3
		known_locals_11 = known_locals_4
		goto b4
	} else {
		form_12 = form_3
		known_locals_13 = known_locals_4
		goto b5
	}
b2:
	;
	or__x_31, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "list?").Deref(), []vm.Value{form_5})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(or__x_31) {
		form_32 = form_5
		known_locals_33 = known_locals_6
		or__x_34 = or__x_31
		goto b10
	} else {
		form_35 = form_5
		known_locals_36 = known_locals_6
		or__x_37 = or__x_31
		goto b11
	}
b3:
	;
	return v600, nil
b4:
	;
	v17, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{form_10})
	if callErr != nil {
		return nil, callErr
	}
	v22 = v17
	form_23 = form_10
	known_locals_24 = known_locals_11
	goto b6
b5:
	;
	v20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	v22 = v20
	form_23 = form_12
	known_locals_24 = known_locals_13
	goto b6
b6:
	;
	v600 = v22
	form_601 = form_23
	known_locals_602 = known_locals_24
	goto b3
b7:
	;
	head_49, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{form_26})
	if callErr != nil {
		return nil, callErr
	}
	v57 = head_49 == vm.Symbol("quote")
	if v57 {
		form_50 = form_26
		known_locals_51 = known_locals_27
		head_52 = head_49
		goto b13
	} else {
		form_53 = form_26
		known_locals_54 = known_locals_27
		head_55 = head_49
		goto b14
	}
b8:
	;
	v552, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector?").Deref(), []vm.Value{form_28})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v552) {
		form_547 = form_28
		known_locals_548 = known_locals_29
		goto b52
	} else {
		form_549 = form_28
		known_locals_550 = known_locals_29
		goto b53
	}
b9:
	;
	v600 = v596
	form_601 = form_597
	known_locals_602 = known_locals_598
	goto b3
b10:
	;
	v43 = or__x_34
	form_44 = form_32
	known_locals_45 = known_locals_33
	or__x_46 = or__x_34
	goto b12
b11:
	;
	v41, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq?").Deref(), []vm.Value{form_35})
	if callErr != nil {
		return nil, callErr
	}
	v43 = v41
	form_44 = form_35
	known_locals_45 = known_locals_36
	or__x_46 = or__x_37
	goto b12
b12:
	;
	if vm.IsTruthy(v43) {
		form_26 = form_44
		known_locals_27 = known_locals_45
		goto b7
	} else {
		form_28 = form_44
		known_locals_29 = known_locals_45
		goto b8
	}
b13:
	;
	v60, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	v542 = v60
	form_543 = form_50
	known_locals_544 = known_locals_51
	head_545 = head_52
	goto b15
b14:
	;
	v69 = head_55 == vm.Symbol("var")
	if v69 {
		form_62 = form_53
		known_locals_63 = known_locals_54
		head_64 = head_55
		goto b16
	} else {
		form_65 = form_53
		known_locals_66 = known_locals_54
		head_67 = head_55
		goto b17
	}
b15:
	;
	v596 = v542
	form_597 = form_543
	known_locals_598 = known_locals_544
	goto b9
b16:
	;
	v72, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	v537 = v72
	form_538 = form_62
	known_locals_539 = known_locals_63
	head_540 = head_64
	goto b18
b17:
	;
	v81 = head_67 == vm.Symbol("set!")
	if v81 {
		form_74 = form_65
		known_locals_75 = known_locals_66
		head_76 = head_67
		goto b19
	} else {
		form_77 = form_65
		known_locals_78 = known_locals_66
		head_79 = head_67
		goto b20
	}
b18:
	;
	v542 = v537
	form_543 = form_538
	known_locals_544 = known_locals_539
	head_545 = head_540
	goto b15
b19:
	;
	__88, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{form_74, vm.Int(0), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	_sym_94, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{form_74, vm.Int(1), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	val_100, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{form_74, vm.Int(2), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	v102, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "captures-of").Deref(), []vm.Value{val_100, known_locals_75})
	if callErr != nil {
		return nil, callErr
	}
	v532 = v102
	form_533 = form_74
	known_locals_534 = known_locals_75
	head_535 = head_76
	goto b21
b20:
	;
	v111 = head_79 == vm.Symbol("fn*")
	if v111 {
		form_104 = form_77
		known_locals_105 = known_locals_78
		head_106 = head_79
		goto b22
	} else {
		form_107 = form_77
		known_locals_108 = known_locals_78
		head_109 = head_79
		goto b23
	}
b21:
	;
	v537 = v532
	form_538 = form_533
	known_locals_539 = known_locals_534
	head_540 = head_535
	goto b18
b22:
	;
	__118, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{form_104, vm.Int(0), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	maybe_name_124, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{form_104, vm.Int(1), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	raw_rest_128, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "drop").Deref(), []vm.Value{vm.Int(2), form_104})
	if callErr != nil {
		return nil, callErr
	}
	has_name_QMARK__130, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "symbol?").Deref(), []vm.Value{maybe_name_124})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(has_name_QMARK__130) {
		vec__5150_131 = form_104
		form_132 = form_104
		known_locals_133 = known_locals_105
		head_134 = head_106
		__135 = __118
		maybe_name_136 = maybe_name_124
		raw_rest_137 = raw_rest_128
		has_name_QMARK__138 = has_name_QMARK__130
		goto b25
	} else {
		vec__5150_139 = form_104
		form_140 = form_104
		known_locals_141 = known_locals_105
		head_142 = head_106
		__143 = __118
		maybe_name_144 = maybe_name_124
		raw_rest_145 = raw_rest_128
		has_name_QMARK__146 = has_name_QMARK__130
		goto b26
	}
b23:
	;
	or__x_345 = head_109 == vm.Symbol("let")
	if or__x_345 {
		form_346 = form_107
		known_locals_347 = known_locals_108
		head_348 = head_109
		or__x_349 = or__x_345
		goto b40
	} else {
		form_350 = form_107
		known_locals_351 = known_locals_108
		head_352 = head_109
		or__x_353 = or__x_345
		goto b41
	}
b24:
	;
	v532 = v527
	form_533 = form_528
	known_locals_534 = known_locals_529
	head_535 = head_530
	goto b21
b25:
	;
	name_sym_151 = maybe_name_136
	vec__5150_152 = vec__5150_131
	form_153 = form_132
	known_locals_154 = known_locals_133
	head_155 = head_134
	__156 = __135
	maybe_name_157 = maybe_name_136
	raw_rest_158 = raw_rest_137
	has_name_QMARK__159 = has_name_QMARK__138
	goto b27
b26:
	;
	name_sym_151 = vm.NIL
	vec__5150_152 = vec__5150_139
	form_153 = form_140
	known_locals_154 = known_locals_141
	head_155 = head_142
	__156 = __143
	maybe_name_157 = maybe_name_144
	raw_rest_158 = raw_rest_145
	has_name_QMARK__159 = has_name_QMARK__146
	goto b27
b27:
	;
	if vm.IsTruthy(has_name_QMARK__159) {
		name_sym_160 = name_sym_151
		vec__5150_161 = vec__5150_152
		form_162 = form_153
		known_locals_163 = known_locals_154
		head_164 = head_155
		__165 = __156
		maybe_name_166 = maybe_name_157
		raw_rest_167 = raw_rest_158
		has_name_QMARK__168 = has_name_QMARK__159
		goto b28
	} else {
		name_sym_169 = name_sym_151
		vec__5150_170 = vec__5150_152
		form_171 = form_153
		known_locals_172 = known_locals_154
		head_173 = head_155
		__174 = __156
		maybe_name_175 = maybe_name_157
		raw_rest_176 = raw_rest_158
		has_name_QMARK__177 = has_name_QMARK__159
		goto b29
	}
b28:
	;
	rest_forms_183 = raw_rest_167
	name_sym_184 = name_sym_160
	vec__5150_185 = vec__5150_161
	form_186 = form_162
	known_locals_187 = known_locals_163
	head_188 = head_164
	__189 = __165
	maybe_name_190 = maybe_name_166
	raw_rest_191 = raw_rest_167
	has_name_QMARK__192 = has_name_QMARK__168
	goto b30
b29:
	;
	v181, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "cons").Deref(), []vm.Value{maybe_name_175, raw_rest_176})
	if callErr != nil {
		return nil, callErr
	}
	rest_forms_183 = v181
	name_sym_184 = name_sym_169
	vec__5150_185 = vec__5150_170
	form_186 = form_171
	known_locals_187 = known_locals_172
	head_188 = head_173
	__189 = __174
	maybe_name_190 = maybe_name_175
	raw_rest_191 = raw_rest_176
	has_name_QMARK__192 = has_name_QMARK__177
	goto b30
b30:
	;
	and__x_194, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{rest_forms_183})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(and__x_194) {
		rest_forms_195 = rest_forms_183
		name_sym_196 = name_sym_184
		vec__5150_197 = vec__5150_185
		form_198 = form_186
		known_locals_199 = known_locals_187
		head_200 = head_188
		__201 = __189
		maybe_name_202 = maybe_name_190
		raw_rest_203 = raw_rest_191
		has_name_QMARK__204 = has_name_QMARK__192
		and__x_205 = and__x_194
		goto b31
	} else {
		rest_forms_206 = rest_forms_183
		name_sym_207 = name_sym_184
		vec__5150_208 = vec__5150_185
		form_209 = form_186
		known_locals_210 = known_locals_187
		head_211 = head_188
		__212 = __189
		maybe_name_213 = maybe_name_190
		raw_rest_214 = raw_rest_191
		has_name_QMARK__215 = has_name_QMARK__192
		and__x_216 = and__x_194
		goto b32
	}
b31:
	;
	arg__5243_219, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{rest_forms_195})
	if callErr != nil {
		return nil, callErr
	}
	arg__5248_222, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{rest_forms_195})
	if callErr != nil {
		return nil, callErr
	}
	v223, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "list?").Deref(), []vm.Value{arg__5248_222})
	if callErr != nil {
		return nil, callErr
	}
	multi_QMARK__226 = v223
	rest_forms_227 = rest_forms_195
	name_sym_228 = name_sym_196
	vec__5150_229 = vec__5150_197
	form_230 = form_198
	known_locals_231 = known_locals_199
	head_232 = head_200
	__233 = __201
	maybe_name_234 = maybe_name_202
	raw_rest_235 = raw_rest_203
	has_name_QMARK__236 = has_name_QMARK__204
	and__x_237 = and__x_205
	goto b33
b32:
	;
	multi_QMARK__226 = and__x_216
	rest_forms_227 = rest_forms_206
	name_sym_228 = name_sym_207
	vec__5150_229 = vec__5150_208
	form_230 = form_209
	known_locals_231 = known_locals_210
	head_232 = head_211
	__233 = __212
	maybe_name_234 = maybe_name_213
	raw_rest_235 = raw_rest_214
	has_name_QMARK__236 = has_name_QMARK__215
	and__x_237 = and__x_216
	goto b33
b33:
	;
	if vm.IsTruthy(multi_QMARK__226) {
		multi_QMARK__241 = multi_QMARK__226
		rest_forms_242 = rest_forms_227
		name_sym_243 = name_sym_228
		vec__5150_244 = vec__5150_229
		form_245 = form_230
		known_locals_246 = known_locals_231
		head_247 = head_232
		__248 = __233
		maybe_name_249 = maybe_name_234
		raw_rest_250 = raw_rest_235
		has_name_QMARK__251 = has_name_QMARK__236
		shadow_for_252 = rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
			var args_vec_2 vm.Value
			var name_sym_3 vm.Value
			var arg__5250_8 vm.Value
			var arg__5254_11 vm.Value
			var arg__5256_12 vm.Value
			var arg__5260_15 vm.Value
			var arg__5264_18 vm.Value
			var arg__5266_19 vm.Value
			var v20 vm.Value
			var args_vec_4 vm.Value
			var name_sym_5 vm.Value
			var arg__5269_23 vm.Value
			var arg__5273_26 vm.Value
			var v27 vm.Value
			var v29 vm.Value
			var args_vec_30 vm.Value
			var name_sym_31 vm.Value
			var callErr error
			_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = args_vec_2, name_sym_3, arg__5250_8, arg__5254_11, arg__5256_12, arg__5260_15, arg__5264_18, arg__5266_19, v20, args_vec_4, name_sym_5, arg__5269_23, arg__5273_26, v27, v29, args_vec_30, name_sym_31
			if vm.IsTruthy(name_sym_228) {
				args_vec_2 = arg0
				name_sym_3 = name_sym_228
				goto b1
			} else {
				args_vec_4 = arg0
				name_sym_5 = name_sym_228
				goto b2
			}
		b1:
			;
			arg__5250_8, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
			if callErr != nil {
				return nil, callErr
			}
			arg__5254_11, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
			if callErr != nil {
				return nil, callErr
			}
			arg__5256_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__5254_11, args_vec_2})
			if callErr != nil {
				return nil, callErr
			}
			arg__5260_15, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
			if callErr != nil {
				return nil, callErr
			}
			arg__5264_18, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
			if callErr != nil {
				return nil, callErr
			}
			arg__5266_19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__5264_18, args_vec_2})
			if callErr != nil {
				return nil, callErr
			}
			v20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{arg__5266_19, name_sym_3})
			if callErr != nil {
				return nil, callErr
			}
			v29 = v20
			args_vec_30 = args_vec_2
			name_sym_31 = name_sym_3
			goto b3
		b2:
			;
			arg__5269_23, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
			if callErr != nil {
				return nil, callErr
			}
			arg__5273_26, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
			if callErr != nil {
				return nil, callErr
			}
			v27, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__5273_26, args_vec_4})
			if callErr != nil {
				return nil, callErr
			}
			v29 = v27
			args_vec_30 = args_vec_4
			name_sym_31 = name_sym_5
			goto b3
		b3:
			;
			return v29, nil
		})
		goto b34
	} else {
		multi_QMARK__253 = multi_QMARK__226
		rest_forms_254 = rest_forms_227
		name_sym_255 = name_sym_228
		vec__5150_256 = vec__5150_229
		form_257 = form_230
		known_locals_258 = known_locals_231
		head_259 = head_232
		__260 = __233
		maybe_name_261 = maybe_name_234
		raw_rest_262 = raw_rest_235
		has_name_QMARK__263 = has_name_QMARK__236
		shadow_for_264 = rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
			var args_vec_2 vm.Value
			var name_sym_3 vm.Value
			var arg__5250_8 vm.Value
			var arg__5254_11 vm.Value
			var arg__5256_12 vm.Value
			var arg__5260_15 vm.Value
			var arg__5264_18 vm.Value
			var arg__5266_19 vm.Value
			var v20 vm.Value
			var args_vec_4 vm.Value
			var name_sym_5 vm.Value
			var arg__5269_23 vm.Value
			var arg__5273_26 vm.Value
			var v27 vm.Value
			var v29 vm.Value
			var args_vec_30 vm.Value
			var name_sym_31 vm.Value
			var callErr error
			_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = args_vec_2, name_sym_3, arg__5250_8, arg__5254_11, arg__5256_12, arg__5260_15, arg__5264_18, arg__5266_19, v20, args_vec_4, name_sym_5, arg__5269_23, arg__5273_26, v27, v29, args_vec_30, name_sym_31
			if vm.IsTruthy(name_sym_228) {
				args_vec_2 = arg0
				name_sym_3 = name_sym_228
				goto b1
			} else {
				args_vec_4 = arg0
				name_sym_5 = name_sym_228
				goto b2
			}
		b1:
			;
			arg__5250_8, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
			if callErr != nil {
				return nil, callErr
			}
			arg__5254_11, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
			if callErr != nil {
				return nil, callErr
			}
			arg__5256_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__5254_11, args_vec_2})
			if callErr != nil {
				return nil, callErr
			}
			arg__5260_15, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
			if callErr != nil {
				return nil, callErr
			}
			arg__5264_18, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
			if callErr != nil {
				return nil, callErr
			}
			arg__5266_19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__5264_18, args_vec_2})
			if callErr != nil {
				return nil, callErr
			}
			v20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{arg__5266_19, name_sym_3})
			if callErr != nil {
				return nil, callErr
			}
			v29 = v20
			args_vec_30 = args_vec_2
			name_sym_31 = name_sym_3
			goto b3
		b2:
			;
			arg__5269_23, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
			if callErr != nil {
				return nil, callErr
			}
			arg__5273_26, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
			if callErr != nil {
				return nil, callErr
			}
			v27, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__5273_26, args_vec_4})
			if callErr != nil {
				return nil, callErr
			}
			v29 = v27
			args_vec_30 = args_vec_4
			name_sym_31 = name_sym_5
			goto b3
		b3:
			;
			return v29, nil
		})
		goto b35
	}
b34:
	;
	arg__5276_267, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__5346_277, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapcat").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var args_vec_4 vm.Value
		var body_6 vm.Value
		var arg__5322_7 vm.Value
		var arg__5328_9 vm.Value
		var inner_known_10 vm.Value
		var v18 vm.Value
		var callErr error
		_, _, _, _, _, _ = args_vec_4, body_6, arg__5322_7, arg__5328_9, inner_known_10, v18
		args_vec_4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		body_6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		arg__5322_7, callErr = rt.InvokeValue(shadow_for_252, []vm.Value{args_vec_4})
		if callErr != nil {
			return nil, callErr
		}
		arg__5328_9, callErr = rt.InvokeValue(shadow_for_252, []vm.Value{args_vec_4})
		if callErr != nil {
			return nil, callErr
		}
		inner_known_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "difference").Deref(), []vm.Value{known_locals_246, arg__5328_9})
		if callErr != nil {
			return nil, callErr
		}
		v18, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapcat").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
			var v3 vm.Value
			var callErr error
			_ = v3
			v3, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "captures-of").Deref(), []vm.Value{arg0, inner_known_10})
			if callErr != nil {
				return nil, callErr
			}
			return v3, nil
		}), body_6})
		if callErr != nil {
			return nil, callErr
		}
		return v18, nil
	}), rest_forms_242})
	if callErr != nil {
		return nil, callErr
	}
	arg__5349_280, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__5419_290, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapcat").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var args_vec_4 vm.Value
		var body_6 vm.Value
		var arg__5395_7 vm.Value
		var arg__5401_9 vm.Value
		var inner_known_10 vm.Value
		var v18 vm.Value
		var callErr error
		_, _, _, _, _, _ = args_vec_4, body_6, arg__5395_7, arg__5401_9, inner_known_10, v18
		args_vec_4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		body_6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		arg__5395_7, callErr = rt.InvokeValue(shadow_for_252, []vm.Value{args_vec_4})
		if callErr != nil {
			return nil, callErr
		}
		arg__5401_9, callErr = rt.InvokeValue(shadow_for_252, []vm.Value{args_vec_4})
		if callErr != nil {
			return nil, callErr
		}
		inner_known_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "difference").Deref(), []vm.Value{known_locals_246, arg__5401_9})
		if callErr != nil {
			return nil, callErr
		}
		v18, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapcat").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
			var v3 vm.Value
			var callErr error
			_ = v3
			v3, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "captures-of").Deref(), []vm.Value{arg0, inner_known_10})
			if callErr != nil {
				return nil, callErr
			}
			return v3, nil
		}), body_6})
		if callErr != nil {
			return nil, callErr
		}
		return v18, nil
	}), rest_forms_242})
	if callErr != nil {
		return nil, callErr
	}
	v291, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__5349_280, arg__5419_290})
	if callErr != nil {
		return nil, callErr
	}
	v324 = v291
	multi_QMARK__325 = multi_QMARK__241
	rest_forms_326 = rest_forms_242
	name_sym_327 = name_sym_243
	vec__5150_328 = vec__5150_244
	form_329 = form_245
	known_locals_330 = known_locals_246
	head_331 = head_247
	__332 = __248
	maybe_name_333 = maybe_name_249
	raw_rest_334 = raw_rest_250
	has_name_QMARK__335 = has_name_QMARK__251
	shadow_for_336 = shadow_for_252
	goto b36
b35:
	;
	args_vec_294, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{rest_forms_254})
	if callErr != nil {
		return nil, callErr
	}
	body_296, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{rest_forms_254})
	if callErr != nil {
		return nil, callErr
	}
	arg__5430_297, callErr = rt.InvokeValue(shadow_for_264, []vm.Value{args_vec_294})
	if callErr != nil {
		return nil, callErr
	}
	arg__5436_299, callErr = rt.InvokeValue(shadow_for_264, []vm.Value{args_vec_294})
	if callErr != nil {
		return nil, callErr
	}
	inner_known_300, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "difference").Deref(), []vm.Value{known_locals_258, arg__5436_299})
	if callErr != nil {
		return nil, callErr
	}
	arg__5438_302, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__5454_310, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapcat").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "captures-of").Deref(), []vm.Value{arg0, inner_known_300})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), body_296})
	if callErr != nil {
		return nil, callErr
	}
	arg__5457_313, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__5473_321, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapcat").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "captures-of").Deref(), []vm.Value{arg0, inner_known_300})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), body_296})
	if callErr != nil {
		return nil, callErr
	}
	v322, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__5457_313, arg__5473_321})
	if callErr != nil {
		return nil, callErr
	}
	v324 = v322
	multi_QMARK__325 = multi_QMARK__253
	rest_forms_326 = rest_forms_254
	name_sym_327 = name_sym_255
	vec__5150_328 = vec__5150_256
	form_329 = form_257
	known_locals_330 = known_locals_258
	head_331 = head_259
	__332 = __260
	maybe_name_333 = maybe_name_261
	raw_rest_334 = raw_rest_262
	has_name_QMARK__335 = has_name_QMARK__263
	shadow_for_336 = shadow_for_264
	goto b36
b36:
	;
	v527 = v324
	form_528 = form_329
	known_locals_529 = known_locals_330
	head_530 = head_331
	goto b24
b37:
	;
	__406, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{form_338, vm.Int(0), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	bindings_412, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{form_338, vm.Int(1), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	body_416, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "drop").Deref(), []vm.Value{vm.Int(2), form_338})
	if callErr != nil {
		return nil, callErr
	}
	pairs_420, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "partition").Deref(), []vm.Value{vm.Int(2), bindings_412})
	if callErr != nil {
		return nil, callErr
	}
	arg__5571_426, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__5573_428, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__5574_429, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__5571_426, arg__5573_428})
	if callErr != nil {
		return nil, callErr
	}
	arg__5642_436, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__5644_438, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__5645_439, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__5642_436, arg__5644_438})
	if callErr != nil {
		return nil, callErr
	}
	vec__5152_440, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "reduce").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
		var caps_8 vm.Value
		var shadowed_14 vm.Value
		var sym_20 vm.Value
		var init_26 vm.Value
		var locs_28 vm.Value
		var arg__5617_31 vm.Value
		var arg__5625_34 vm.Value
		var arg__5626_35 vm.Value
		var arg__5631_37 vm.Value
		var arg__5637_40 vm.Value
		var arg__5638_41 vm.Value
		var v42 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _, _, _, _, _ = caps_8, shadowed_14, sym_20, init_26, locs_28, arg__5617_31, arg__5625_34, arg__5626_35, arg__5631_37, arg__5637_40, arg__5638_41, v42
		caps_8, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(0), vm.NIL})
		if callErr != nil {
			return nil, callErr
		}
		shadowed_14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(1), vm.NIL})
		if callErr != nil {
			return nil, callErr
		}
		sym_20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg1, vm.Int(0), vm.NIL})
		if callErr != nil {
			return nil, callErr
		}
		init_26, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg1, vm.Int(1), vm.NIL})
		if callErr != nil {
			return nil, callErr
		}
		locs_28, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "difference").Deref(), []vm.Value{known_locals_339, shadowed_14})
		if callErr != nil {
			return nil, callErr
		}
		arg__5617_31, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "captures-of").Deref(), []vm.Value{init_26, locs_28})
		if callErr != nil {
			return nil, callErr
		}
		arg__5625_34, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "captures-of").Deref(), []vm.Value{init_26, locs_28})
		if callErr != nil {
			return nil, callErr
		}
		arg__5626_35, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{caps_8, arg__5625_34})
		if callErr != nil {
			return nil, callErr
		}
		arg__5631_37, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "binding-syms").Deref(), []vm.Value{sym_20})
		if callErr != nil {
			return nil, callErr
		}
		arg__5637_40, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "binding-syms").Deref(), []vm.Value{sym_20})
		if callErr != nil {
			return nil, callErr
		}
		arg__5638_41, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{shadowed_14, arg__5637_40})
		if callErr != nil {
			return nil, callErr
		}
		v42, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__5626_35, arg__5638_41})
		if callErr != nil {
			return nil, callErr
		}
		return v42, nil
	}), arg__5645_439, pairs_420})
	if callErr != nil {
		return nil, callErr
	}
	used_446, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{vec__5152_440, vm.Int(0), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	let_bound_set_452, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{vec__5152_440, vm.Int(1), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	new_locals_454, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{known_locals_339, let_bound_set_452})
	if callErr != nil {
		return nil, callErr
	}
	arg__5667_456, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__5683_464, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapcat").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "captures-of").Deref(), []vm.Value{arg0, new_locals_454})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), body_416})
	if callErr != nil {
		return nil, callErr
	}
	arg__5686_467, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__5702_475, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapcat").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "captures-of").Deref(), []vm.Value{arg0, new_locals_454})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), body_416})
	if callErr != nil {
		return nil, callErr
	}
	body_captures_476, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__5686_467, arg__5702_475})
	if callErr != nil {
		return nil, callErr
	}
	arg__5709_478, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "difference").Deref(), []vm.Value{body_captures_476, let_bound_set_452})
	if callErr != nil {
		return nil, callErr
	}
	arg__5717_481, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "difference").Deref(), []vm.Value{body_captures_476, let_bound_set_452})
	if callErr != nil {
		return nil, callErr
	}
	v482, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{used_446, arg__5717_481})
	if callErr != nil {
		return nil, callErr
	}
	v522 = v482
	form_523 = form_338
	known_locals_524 = known_locals_339
	head_525 = head_340
	goto b39
b38:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		form_484 = form_341
		known_locals_485 = known_locals_342
		head_486 = head_343
		goto b49
	} else {
		form_487 = form_341
		known_locals_488 = known_locals_342
		head_489 = head_343
		goto b50
	}
b39:
	;
	v527 = v522
	form_528 = form_523
	known_locals_529 = known_locals_524
	head_530 = head_525
	goto b24
b40:
	;
	v395 = or__x_349
	form_396 = form_346
	known_locals_397 = known_locals_347
	head_398 = head_348
	or__x_399 = vm.Boolean(or__x_349)
	goto b42
b41:
	;
	or__x_357 = head_352 == vm.Symbol("let*")
	if or__x_357 {
		form_358 = form_350
		known_locals_359 = known_locals_351
		head_360 = head_352
		or__x_361 = or__x_357
		goto b43
	} else {
		form_362 = form_350
		known_locals_363 = known_locals_351
		head_364 = head_352
		or__x_365 = or__x_357
		goto b44
	}
b42:
	;
	if v395 {
		form_338 = form_396
		known_locals_339 = known_locals_397
		head_340 = head_398
		goto b37
	} else {
		form_341 = form_396
		known_locals_342 = known_locals_397
		head_343 = head_398
		goto b38
	}
b43:
	;
	v389 = or__x_361
	form_390 = form_358
	known_locals_391 = known_locals_359
	head_392 = head_360
	or__x_393 = vm.Boolean(or__x_361)
	goto b45
b44:
	;
	or__x_369 = head_364 == vm.Symbol("loop")
	if or__x_369 {
		form_370 = form_362
		known_locals_371 = known_locals_363
		head_372 = head_364
		or__x_373 = or__x_369
		goto b46
	} else {
		form_374 = form_362
		known_locals_375 = known_locals_363
		head_376 = head_364
		or__x_377 = or__x_369
		goto b47
	}
b45:
	;
	v395 = v389
	form_396 = form_390
	known_locals_397 = known_locals_391
	head_398 = head_392
	or__x_399 = vm.Boolean(or__x_353)
	goto b42
b46:
	;
	v383 = or__x_373
	form_384 = form_370
	known_locals_385 = known_locals_371
	head_386 = head_372
	or__x_387 = vm.Boolean(or__x_373)
	goto b48
b47:
	;
	v381 = head_376 == vm.Symbol("loop*")
	v383 = v381
	form_384 = form_374
	known_locals_385 = known_locals_375
	head_386 = head_376
	or__x_387 = vm.Boolean(or__x_377)
	goto b48
b48:
	;
	v389 = v383
	form_390 = form_384
	known_locals_391 = known_locals_385
	head_392 = head_386
	or__x_393 = vm.Boolean(or__x_365)
	goto b45
b49:
	;
	arg__5719_493, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__5735_501, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapcat").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "captures-of").Deref(), []vm.Value{arg0, known_locals_485})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), form_484})
	if callErr != nil {
		return nil, callErr
	}
	arg__5738_504, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__5754_512, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapcat").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "captures-of").Deref(), []vm.Value{arg0, known_locals_485})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), form_484})
	if callErr != nil {
		return nil, callErr
	}
	v513, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__5738_504, arg__5754_512})
	if callErr != nil {
		return nil, callErr
	}
	v517 = v513
	form_518 = form_484
	known_locals_519 = known_locals_485
	head_520 = head_486
	goto b51
b50:
	;
	v517 = vm.NIL
	form_518 = form_487
	known_locals_519 = known_locals_488
	head_520 = head_489
	goto b51
b51:
	;
	v522 = v517
	form_523 = form_518
	known_locals_524 = known_locals_519
	head_525 = head_520
	goto b39
b52:
	;
	arg__5759_555, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__5775_563, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapcat").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "captures-of").Deref(), []vm.Value{arg0, known_locals_548})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), form_547})
	if callErr != nil {
		return nil, callErr
	}
	arg__5778_566, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__5794_574, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapcat").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "captures-of").Deref(), []vm.Value{arg0, known_locals_548})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), form_547})
	if callErr != nil {
		return nil, callErr
	}
	v575, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__5778_566, arg__5794_574})
	if callErr != nil {
		return nil, callErr
	}
	v592 = v575
	form_593 = form_547
	known_locals_594 = known_locals_548
	goto b54
b53:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		form_577 = form_549
		known_locals_578 = known_locals_550
		goto b55
	} else {
		form_579 = form_549
		known_locals_580 = known_locals_550
		goto b56
	}
b54:
	;
	v596 = v592
	form_597 = form_593
	known_locals_598 = known_locals_594
	goto b9
b55:
	;
	v584, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	v588 = v584
	form_589 = form_577
	known_locals_590 = known_locals_578
	goto b57
b56:
	;
	v588 = vm.NIL
	form_589 = form_579
	known_locals_590 = known_locals_580
	goto b57
b57:
	;
	v592 = v588
	form_593 = form_589
	known_locals_594 = known_locals_590
	goto b54
}
func ctx_block(arg0 vm.Value) (vm.Value, error) {
	var arg__5799_2 vm.Value
	var arg__5805_6 vm.Value
	var v8 vm.Value
	var callErr error
	_, _, _ = arg__5799_2, arg__5805_6, v8
	arg__5799_2, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__5805_6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v8, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg__5805_6, vm.Keyword("current-block")})
	if callErr != nil {
		return nil, callErr
	}
	return v8, nil
}
func ctx_fn(arg0 vm.Value) (vm.Value, error) {
	var arg__5810_2 vm.Value
	var arg__5816_6 vm.Value
	var v8 vm.Value
	var callErr error
	_, _, _ = arg__5810_2, arg__5816_6, v8
	arg__5810_2, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__5816_6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v8, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg__5816_6, vm.Keyword("fn")})
	if callErr != nil {
		return nil, callErr
	}
	return v8, nil
}
func ctx_set_block_BANG_(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__5822_3 vm.Value
	var arg__5829_7 vm.Value
	var arg__5832_9 vm.Value
	var arg__5838_12 vm.Value
	var arg__5845_16 vm.Value
	var arg__5848_18 vm.Value
	var v19 vm.Value
	var callErr error
	_, _, _, _, _, _, _ = arg__5822_3, arg__5829_7, arg__5832_9, arg__5838_12, arg__5845_16, arg__5848_18, v19
	arg__5822_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__5829_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__5832_9, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "assoc").Deref(), []vm.Value{arg__5829_7, vm.Keyword("current-block"), arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__5838_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__5845_16, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__5848_18, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "assoc").Deref(), []vm.Value{arg__5845_16, vm.Keyword("current-block"), arg1})
	if callErr != nil {
		return nil, callErr
	}
	v19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "reset!").Deref(), []vm.Value{arg0, arg__5848_18})
	if callErr != nil {
		return nil, callErr
	}
	return v19, nil
}
func current_locals_flat(arg0 vm.Value) (vm.Value, error) {
	var arg__5854_5 vm.Value
	var arg__5860_9 vm.Value
	var arg__5862_11 vm.Value
	var arg__5869_16 vm.Value
	var arg__5875_20 vm.Value
	var arg__5877_22 vm.Value
	var v23 vm.Value
	var callErr error
	_, _, _, _, _, _, _ = arg__5854_5, arg__5860_9, arg__5862_11, arg__5869_16, arg__5875_20, arg__5877_22, v23
	arg__5854_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__5860_9, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__5862_11, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg__5860_9, vm.Keyword("locals")})
	if callErr != nil {
		return nil, callErr
	}
	arg__5869_16, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__5875_20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__5877_22, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg__5875_20, vm.Keyword("locals")})
	if callErr != nil {
		return nil, callErr
	}
	v23, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "reduce").Deref(), []vm.Value{rt.LookupVar("clojure.core", "merge").Deref(), vm.EmptyPersistentMap, arg__5877_22})
	if callErr != nil {
		return nil, callErr
	}
	return v23, nil
}
func emit_template_closure(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__5881_4 vm.Value
	var arg__5885_6 vm.Value
	var arg__5893_11 vm.Value
	var arg__5897_13 vm.Value
	var const_id_16 vm.Value
	var v26 vm.Value
	var template_17 vm.Value
	var capture_syms_18 vm.Value
	var ctx_19 vm.Value
	var const_id_20 vm.Value
	var arg__5907_29 vm.Value
	var arg__5911_31 vm.Value
	var arg__5915_34 vm.Value
	var arg__5921_38 vm.Value
	var arg__5925_40 vm.Value
	var arg__5929_43 vm.Value
	var v45 vm.Value
	var template_21 vm.Value
	var capture_syms_22 vm.Value
	var ctx_23 vm.Value
	var const_id_24 vm.Value
	var closure_id_48 vm.Value
	var template_49 vm.Value
	var capture_syms_50 vm.Value
	var ctx_51 vm.Value
	var const_id_52 vm.Value
	var cls_53 vm.Value
	var caps_54 vm.Value
	var ctx_55 vm.Value
	var v113 vm.Value
	var v116 vm.Value
	var v72 vm.Value
	var closure_id_57 vm.Value
	var template_58 vm.Value
	var capture_syms_59 vm.Value
	var const_id_60 vm.Value
	var cls_61 vm.Value
	var caps_62 vm.Value
	var ctx_63 vm.Value
	var v114 vm.Value
	var v117 vm.Value
	var cap_sym_75 vm.Value
	var cap_val_77 vm.Value
	var arg__5945_79 vm.Value
	var arg__5949_81 vm.Value
	var arg__5954_84 vm.Value
	var arg__5960_88 vm.Value
	var arg__5964_90 vm.Value
	var arg__5969_93 vm.Value
	var push_id_95 vm.Value
	var v97 vm.Value
	var closure_id_64 vm.Value
	var template_65 vm.Value
	var capture_syms_66 vm.Value
	var const_id_67 vm.Value
	var cls_68 vm.Value
	var caps_69 vm.Value
	var ctx_70 vm.Value
	var v115 vm.Value
	var v118 vm.Value
	var v100 vm.Value
	var closure_id_101 vm.Value
	var template_102 vm.Value
	var capture_syms_103 vm.Value
	var const_id_104 vm.Value
	var cls_105 vm.Value
	var caps_106 vm.Value
	var ctx_107 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__5881_4, arg__5885_6, arg__5893_11, arg__5897_13, const_id_16, v26, template_17, capture_syms_18, ctx_19, const_id_20, arg__5907_29, arg__5911_31, arg__5915_34, arg__5921_38, arg__5925_40, arg__5929_43, v45, template_21, capture_syms_22, ctx_23, const_id_24, closure_id_48, template_49, capture_syms_50, ctx_51, const_id_52, cls_53, caps_54, ctx_55, v113, v116, v72, closure_id_57, template_58, capture_syms_59, const_id_60, cls_61, caps_62, ctx_63, v114, v117, cap_sym_75, cap_val_77, arg__5945_79, arg__5949_81, arg__5954_84, arg__5960_88, arg__5964_90, arg__5969_93, push_id_95, v97, closure_id_64, template_65, capture_syms_66, const_id_67, cls_68, caps_69, ctx_70, v115, v118, v100, closure_id_101, template_102, capture_syms_103, const_id_104, cls_105, caps_106, ctx_107
	arg__5881_4, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-fn").Deref(), []vm.Value{arg2})
	if callErr != nil {
		return nil, callErr
	}
	arg__5885_6, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{arg2})
	if callErr != nil {
		return nil, callErr
	}
	arg__5893_11, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-fn").Deref(), []vm.Value{arg2})
	if callErr != nil {
		return nil, callErr
	}
	arg__5897_13, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{arg2})
	if callErr != nil {
		return nil, callErr
	}
	const_id_16, callErr = rt.InvokeValue(rt.LookupVar("ir", "add-inst").Deref(), []vm.Value{arg__5893_11, arg__5897_13, vm.Keyword("const"), vm.NewArrayVector([]vm.Value{}), arg0})
	if callErr != nil {
		return nil, callErr
	}
	v26, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v26) {
		template_17 = arg0
		capture_syms_18 = arg1
		ctx_19 = arg2
		const_id_20 = const_id_16
		goto b1
	} else {
		template_21 = arg0
		capture_syms_22 = arg1
		ctx_23 = arg2
		const_id_24 = const_id_16
		goto b2
	}
b1:
	;
	arg__5907_29, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-fn").Deref(), []vm.Value{ctx_19})
	if callErr != nil {
		return nil, callErr
	}
	arg__5911_31, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_19})
	if callErr != nil {
		return nil, callErr
	}
	arg__5915_34, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{const_id_20})
	if callErr != nil {
		return nil, callErr
	}
	arg__5921_38, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-fn").Deref(), []vm.Value{ctx_19})
	if callErr != nil {
		return nil, callErr
	}
	arg__5925_40, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_19})
	if callErr != nil {
		return nil, callErr
	}
	arg__5929_43, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{const_id_20})
	if callErr != nil {
		return nil, callErr
	}
	v45, callErr = rt.InvokeValue(rt.LookupVar("ir", "add-inst").Deref(), []vm.Value{arg__5921_38, arg__5925_40, vm.Keyword("make-closure"), arg__5929_43, vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	closure_id_48 = v45
	template_49 = template_17
	capture_syms_50 = capture_syms_18
	ctx_51 = ctx_19
	const_id_52 = const_id_20
	goto b3
b2:
	;
	closure_id_48 = const_id_24
	template_49 = template_21
	capture_syms_50 = capture_syms_22
	ctx_51 = ctx_23
	const_id_52 = const_id_24
	goto b3
b3:
	;
	cls_53 = closure_id_48
	caps_54 = capture_syms_50
	ctx_55 = ctx_51
	v113 = vm.Keyword("push-closed")
	v116 = vm.NIL
	goto b4
b4:
	;
	v72, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{caps_54})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v72) {
		closure_id_57 = closure_id_48
		template_58 = template_49
		capture_syms_59 = capture_syms_50
		const_id_60 = const_id_52
		cls_61 = cls_53
		caps_62 = caps_54
		ctx_63 = ctx_55
		v114 = v113
		v117 = v116
		goto b5
	} else {
		closure_id_64 = closure_id_48
		template_65 = template_49
		capture_syms_66 = capture_syms_50
		const_id_67 = const_id_52
		cls_68 = cls_53
		caps_69 = caps_54
		ctx_70 = ctx_55
		v115 = v113
		v118 = v116
		goto b6
	}
b5:
	;
	cap_sym_75, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{caps_62})
	if callErr != nil {
		return nil, callErr
	}
	cap_val_77, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-symbol").Deref(), []vm.Value{cap_sym_75, ctx_63})
	if callErr != nil {
		return nil, callErr
	}
	arg__5945_79, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-fn").Deref(), []vm.Value{ctx_63})
	if callErr != nil {
		return nil, callErr
	}
	arg__5949_81, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_63})
	if callErr != nil {
		return nil, callErr
	}
	arg__5954_84, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{cls_61, cap_val_77})
	if callErr != nil {
		return nil, callErr
	}
	arg__5960_88, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-fn").Deref(), []vm.Value{ctx_63})
	if callErr != nil {
		return nil, callErr
	}
	arg__5964_90, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_63})
	if callErr != nil {
		return nil, callErr
	}
	arg__5969_93, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{cls_61, cap_val_77})
	if callErr != nil {
		return nil, callErr
	}
	push_id_95, callErr = rt.InvokeValue(rt.LookupVar("ir", "add-inst").Deref(), []vm.Value{arg__5960_88, arg__5964_90, v114, arg__5969_93, v117})
	if callErr != nil {
		return nil, callErr
	}
	v97, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{caps_62})
	if callErr != nil {
		return nil, callErr
	}
	cls_53 = push_id_95
	caps_54 = v97
	ctx_55 = ctx_63
	v113 = v114
	v116 = v117
	goto b4
b6:
	;
	v100 = cls_68
	closure_id_101 = closure_id_64
	template_102 = template_65
	capture_syms_103 = capture_syms_66
	const_id_104 = const_id_67
	cls_105 = cls_68
	caps_106 = caps_69
	ctx_107 = ctx_70
	goto b7
b7:
	;
	return v100, nil
}
func expand_binding(arg0 vm.Value) (vm.Value, error) {
	var b_1 vm.Value
	var out_2 vm.Value
	var v12 vm.Value
	var bindings_5 vm.Value
	var b_6 vm.Value
	var out_7 vm.Value
	var bindings_8 vm.Value
	var b_9 vm.Value
	var out_10 vm.Value
	var n_16 vm.Value
	var v_18 vm.Value
	var v30 vm.Value
	var v150 vm.Value
	var bindings_151 vm.Value
	var b_152 vm.Value
	var out_153 vm.Value
	var bindings_19 vm.Value
	var b_20 vm.Value
	var out_21 vm.Value
	var n_22 vm.Value
	var v_23 vm.Value
	var gs_35 vm.Value
	var v39 vm.Value
	var arg__6001_41 vm.Value
	var arg__6007_43 vm.Value
	var arg__6016_46 vm.Value
	var arg__6022_48 vm.Value
	var v49 vm.Value
	var bindings_24 vm.Value
	var b_25 vm.Value
	var out_26 vm.Value
	var n_27 vm.Value
	var v_28 vm.Value
	var v62 vm.Value
	var v143 vm.Value
	var bindings_144 vm.Value
	var b_145 vm.Value
	var out_146 vm.Value
	var n_147 vm.Value
	var v_148 vm.Value
	var bindings_51 vm.Value
	var b_52 vm.Value
	var out_53 vm.Value
	var n_54 vm.Value
	var v_55 vm.Value
	var gs_67 vm.Value
	var v71 vm.Value
	var arg__6041_73 vm.Value
	var arg__6047_75 vm.Value
	var arg__6056_78 vm.Value
	var arg__6062_80 vm.Value
	var v81 vm.Value
	var bindings_56 vm.Value
	var b_57 vm.Value
	var out_58 vm.Value
	var n_59 vm.Value
	var v_60 vm.Value
	var v94 vm.Value
	var v136 vm.Value
	var bindings_137 vm.Value
	var b_138 vm.Value
	var out_139 vm.Value
	var n_140 vm.Value
	var v_141 vm.Value
	var bindings_83 vm.Value
	var b_84 vm.Value
	var out_85 vm.Value
	var n_86 vm.Value
	var v_87 vm.Value
	var v99 vm.Value
	var bindings_88 vm.Value
	var b_89 vm.Value
	var out_90 vm.Value
	var n_91 vm.Value
	var v_92 vm.Value
	var v129 vm.Value
	var bindings_130 vm.Value
	var b_131 vm.Value
	var out_132 vm.Value
	var n_133 vm.Value
	var v_134 vm.Value
	var bindings_101 vm.Value
	var b_102 vm.Value
	var out_103 vm.Value
	var n_104 vm.Value
	var v_105 vm.Value
	var v116 vm.Value
	var v118 vm.Value
	var bindings_106 vm.Value
	var b_107 vm.Value
	var out_108 vm.Value
	var n_109 vm.Value
	var v_110 vm.Value
	var v122 vm.Value
	var bindings_123 vm.Value
	var b_124 vm.Value
	var out_125 vm.Value
	var n_126 vm.Value
	var v_127 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = b_1, out_2, v12, bindings_5, b_6, out_7, bindings_8, b_9, out_10, n_16, v_18, v30, v150, bindings_151, b_152, out_153, bindings_19, b_20, out_21, n_22, v_23, gs_35, v39, arg__6001_41, arg__6007_43, arg__6016_46, arg__6022_48, v49, bindings_24, b_25, out_26, n_27, v_28, v62, v143, bindings_144, b_145, out_146, n_147, v_148, bindings_51, b_52, out_53, n_54, v_55, gs_67, v71, arg__6041_73, arg__6047_75, arg__6056_78, arg__6062_80, v81, bindings_56, b_57, out_58, n_59, v_60, v94, v136, bindings_137, b_138, out_139, n_140, v_141, bindings_83, b_84, out_85, n_86, v_87, v99, bindings_88, b_89, out_90, n_91, v_92, v129, bindings_130, b_131, out_132, n_133, v_134, bindings_101, b_102, out_103, n_104, v_105, v116, v118, bindings_106, b_107, out_108, n_109, v_110, v122, bindings_123, b_124, out_125, n_126, v_127
	b_1 = arg0
	out_2 = vm.NewArrayVector([]vm.Value{})
	goto b1
b1:
	;
	v12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "empty?").Deref(), []vm.Value{b_1})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v12) {
		bindings_5 = arg0
		b_6 = b_1
		out_7 = out_2
		goto b2
	} else {
		bindings_8 = arg0
		b_9 = b_1
		out_10 = out_2
		goto b3
	}
b2:
	;
	v150 = out_7
	bindings_151 = bindings_5
	b_152 = b_6
	out_153 = out_7
	goto b4
b3:
	;
	n_16, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{b_9})
	if callErr != nil {
		return nil, callErr
	}
	v_18, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "second").Deref(), []vm.Value{b_9})
	if callErr != nil {
		return nil, callErr
	}
	v30, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector?").Deref(), []vm.Value{n_16})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v30) {
		bindings_19 = bindings_8
		b_20 = b_9
		out_21 = out_10
		n_22 = n_16
		v_23 = v_18
		goto b5
	} else {
		bindings_24 = bindings_8
		b_25 = b_9
		out_26 = out_10
		n_27 = n_16
		v_28 = v_18
		goto b6
	}
b4:
	;
	return v150, nil
b5:
	;
	gs_35, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "gensym").Deref(), []vm.Value{vm.String("vec__")})
	if callErr != nil {
		return nil, callErr
	}
	v39, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "drop").Deref(), []vm.Value{vm.Int(2), b_20})
	if callErr != nil {
		return nil, callErr
	}
	arg__6001_41, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{out_21, gs_35, v_23})
	if callErr != nil {
		return nil, callErr
	}
	arg__6007_43, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "expand-vector-pattern").Deref(), []vm.Value{gs_35, n_22})
	if callErr != nil {
		return nil, callErr
	}
	arg__6016_46, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{out_21, gs_35, v_23})
	if callErr != nil {
		return nil, callErr
	}
	arg__6022_48, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "expand-vector-pattern").Deref(), []vm.Value{gs_35, n_22})
	if callErr != nil {
		return nil, callErr
	}
	v49, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__6016_46, arg__6022_48})
	if callErr != nil {
		return nil, callErr
	}
	b_1 = v39
	out_2 = v49
	goto b1
b6:
	;
	v62, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map?").Deref(), []vm.Value{n_27})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v62) {
		bindings_51 = bindings_24
		b_52 = b_25
		out_53 = out_26
		n_54 = n_27
		v_55 = v_28
		goto b8
	} else {
		bindings_56 = bindings_24
		b_57 = b_25
		out_58 = out_26
		n_59 = n_27
		v_60 = v_28
		goto b9
	}
b7:
	;
	v150 = v143
	bindings_151 = bindings_144
	b_152 = b_145
	out_153 = out_146
	goto b4
b8:
	;
	gs_67, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "gensym").Deref(), []vm.Value{vm.String("map__")})
	if callErr != nil {
		return nil, callErr
	}
	v71, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "drop").Deref(), []vm.Value{vm.Int(2), b_52})
	if callErr != nil {
		return nil, callErr
	}
	arg__6041_73, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{out_53, gs_67, v_55})
	if callErr != nil {
		return nil, callErr
	}
	arg__6047_75, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "expand-map-pattern").Deref(), []vm.Value{gs_67, n_54})
	if callErr != nil {
		return nil, callErr
	}
	arg__6056_78, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{out_53, gs_67, v_55})
	if callErr != nil {
		return nil, callErr
	}
	arg__6062_80, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "expand-map-pattern").Deref(), []vm.Value{gs_67, n_54})
	if callErr != nil {
		return nil, callErr
	}
	v81, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__6056_78, arg__6062_80})
	if callErr != nil {
		return nil, callErr
	}
	b_1 = v71
	out_2 = v81
	goto b1
b9:
	;
	v94, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{n_59})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v94) {
		bindings_83 = bindings_56
		b_84 = b_57
		out_85 = out_58
		n_86 = n_59
		v_87 = v_60
		goto b11
	} else {
		bindings_88 = bindings_56
		b_89 = b_57
		out_90 = out_58
		n_91 = n_59
		v_92 = v_60
		goto b12
	}
b10:
	;
	v143 = v136
	bindings_144 = bindings_137
	b_145 = b_138
	out_146 = out_139
	n_147 = n_140
	v_148 = v_141
	goto b7
b11:
	;
	v99, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "drop").Deref(), []vm.Value{vm.Int(2), b_84})
	if callErr != nil {
		return nil, callErr
	}
	b_1 = v99
	out_2 = out_85
	goto b1
b12:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		bindings_101 = bindings_88
		b_102 = b_89
		out_103 = out_90
		n_104 = n_91
		v_105 = v_92
		goto b14
	} else {
		bindings_106 = bindings_88
		b_107 = b_89
		out_108 = out_90
		n_109 = n_91
		v_110 = v_92
		goto b15
	}
b13:
	;
	v136 = v129
	bindings_137 = bindings_130
	b_138 = b_131
	out_139 = out_132
	n_140 = n_133
	v_141 = v_134
	goto b10
b14:
	;
	v116, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "drop").Deref(), []vm.Value{vm.Int(2), b_102})
	if callErr != nil {
		return nil, callErr
	}
	v118, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{out_103, n_104, v_105})
	if callErr != nil {
		return nil, callErr
	}
	b_1 = v116
	out_2 = v118
	goto b1
b15:
	;
	v122 = vm.NIL
	bindings_123 = bindings_106
	b_124 = b_107
	out_125 = out_108
	n_126 = n_109
	v_127 = v_110
	goto b16
b16:
	;
	v129 = v122
	bindings_130 = bindings_123
	b_131 = b_124
	out_132 = out_125
	n_133 = n_126
	v_134 = v_127
	goto b13
}
func expand_fn_args(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var i_3 int
	var remaining_4 vm.Value
	var v572 vm.Value
	var v16 vm.Value
	var args_vec_7 vm.Value
	var body_forms_8 vm.Value
	var i_9 int
	var remaining_10 vm.Value
	var v575 vm.Value
	var args_vec_11 vm.Value
	var body_forms_12 vm.Value
	var i_13 int
	var remaining_14 vm.Value
	var v573 vm.Value
	var arg__6089_29 vm.Value
	var v31 bool
	var amp_pos_44 int
	var args_vec_45 vm.Value
	var body_forms_46 vm.Value
	var i_47 int
	var remaining_48 vm.Value
	var variadic_QMARK__52 vm.Value
	var args_vec_20 vm.Value
	var body_forms_21 vm.Value
	var i_22 int
	var remaining_23 vm.Value
	var v576 vm.Value
	var args_vec_24 vm.Value
	var body_forms_25 vm.Value
	var i_26 int
	var remaining_27 vm.Value
	var v574 vm.Value
	var v34 int
	var v36 vm.Value
	var v38 int
	var args_vec_39 vm.Value
	var body_forms_40 vm.Value
	var i_41 int
	var remaining_42 vm.Value
	var amp_pos_53 int
	var args_vec_54 vm.Value
	var body_forms_55 vm.Value
	var i_56 int
	var remaining_57 vm.Value
	var variadic_QMARK__58 vm.Value
	var v67 vm.Value
	var amp_pos_59 int
	var args_vec_60 vm.Value
	var body_forms_61 vm.Value
	var i_62 int
	var remaining_63 vm.Value
	var variadic_QMARK__64 vm.Value
	var fixed_args_70 vm.Value
	var amp_pos_71 int
	var args_vec_72 vm.Value
	var body_forms_73 vm.Value
	var i_74 int
	var remaining_75 vm.Value
	var variadic_QMARK__76 vm.Value
	var fixed_args_77 vm.Value
	var amp_pos_78 int
	var args_vec_79 vm.Value
	var body_forms_80 vm.Value
	var i_81 int
	var remaining_82 vm.Value
	var variadic_QMARK__83 vm.Value
	var arg__6107_92 int
	var v97 vm.Value
	var fixed_args_84 vm.Value
	var amp_pos_85 int
	var args_vec_86 vm.Value
	var body_forms_87 vm.Value
	var i_88 int
	var remaining_89 vm.Value
	var variadic_QMARK__90 vm.Value
	var rest_sym_101 vm.Value
	var fixed_args_102 vm.Value
	var amp_pos_103 int
	var args_vec_104 vm.Value
	var body_forms_105 vm.Value
	var i_106 int
	var remaining_107 vm.Value
	var variadic_QMARK__108 vm.Value
	var has_destructure_QMARK__112 vm.Value
	var rest_sym_113 vm.Value
	var fixed_args_114 vm.Value
	var amp_pos_115 int
	var args_vec_116 vm.Value
	var body_forms_117 vm.Value
	var i_118 int
	var remaining_119 vm.Value
	var variadic_QMARK__120 vm.Value
	var has_destructure_QMARK__121 vm.Value
	var rest_sym_122 vm.Value
	var fixed_args_123 vm.Value
	var amp_pos_124 int
	var args_vec_125 vm.Value
	var body_forms_126 vm.Value
	var i_127 int
	var remaining_128 vm.Value
	var variadic_QMARK__129 vm.Value
	var has_destructure_QMARK__130 vm.Value
	var v560 vm.Value
	var rest_sym_561 vm.Value
	var fixed_args_562 vm.Value
	var amp_pos_563 int
	var args_vec_564 vm.Value
	var body_forms_565 vm.Value
	var i_566 int
	var remaining_567 vm.Value
	var variadic_QMARK__568 vm.Value
	var has_destructure_QMARK__569 vm.Value
	var rest_sym_131 vm.Value
	var fixed_args_132 vm.Value
	var amp_pos_133 int
	var args_vec_134 vm.Value
	var body_forms_135 vm.Value
	var i_136 int
	var remaining_137 vm.Value
	var or__x_138 vm.Value
	var variadic_QMARK__139 vm.Value
	var has_destructure_QMARK__140 vm.Value
	var rest_sym_141 vm.Value
	var fixed_args_142 vm.Value
	var amp_pos_143 int
	var args_vec_144 vm.Value
	var body_forms_145 vm.Value
	var i_146 int
	var remaining_147 vm.Value
	var or__x_148 vm.Value
	var variadic_QMARK__149 vm.Value
	var has_destructure_QMARK__150 vm.Value
	var v154 vm.Value
	var rest_sym_155 vm.Value
	var fixed_args_156 vm.Value
	var amp_pos_157 int
	var args_vec_158 vm.Value
	var body_forms_159 vm.Value
	var i_160 int
	var remaining_161 vm.Value
	var or__x_162 vm.Value
	var variadic_QMARK__163 vm.Value
	var has_destructure_QMARK__164 vm.Value
	var remaining_166 vm.Value
	var flat_args_167 vm.Value
	var let_binds_168 vm.Value
	var remaining_169 vm.Value
	var v196 vm.Value
	var rest_sym_173 vm.Value
	var fixed_args_174 vm.Value
	var amp_pos_175 int
	var args_vec_176 vm.Value
	var body_forms_177 vm.Value
	var i_178 int
	var variadic_QMARK__179 vm.Value
	var has_destructure_QMARK__180 vm.Value
	var flat_args_181 vm.Value
	var let_binds_182 vm.Value
	var remaining_183 vm.Value
	var v201 vm.Value
	var rest_sym_184 vm.Value
	var fixed_args_185 vm.Value
	var amp_pos_186 int
	var args_vec_187 vm.Value
	var body_forms_188 vm.Value
	var i_189 int
	var variadic_QMARK__190 vm.Value
	var has_destructure_QMARK__191 vm.Value
	var flat_args_192 vm.Value
	var let_binds_193 vm.Value
	var remaining_194 vm.Value
	var x_204 vm.Value
	var v230 vm.Value
	var result_405 vm.Value
	var rest_sym_406 vm.Value
	var fixed_args_407 vm.Value
	var amp_pos_408 int
	var args_vec_409 vm.Value
	var body_forms_410 vm.Value
	var i_411 int
	var variadic_QMARK__412 vm.Value
	var has_destructure_QMARK__413 vm.Value
	var flat_args_414 vm.Value
	var let_binds_415 vm.Value
	var remaining_416 vm.Value
	var rest_sym_205 vm.Value
	var fixed_args_206 vm.Value
	var amp_pos_207 int
	var args_vec_208 vm.Value
	var body_forms_209 vm.Value
	var i_210 int
	var variadic_QMARK__211 vm.Value
	var has_destructure_QMARK__212 vm.Value
	var flat_args_213 vm.Value
	var let_binds_214 vm.Value
	var remaining_215 vm.Value
	var x_216 vm.Value
	var v233 vm.Value
	var v235 vm.Value
	var rest_sym_217 vm.Value
	var fixed_args_218 vm.Value
	var amp_pos_219 int
	var args_vec_220 vm.Value
	var body_forms_221 vm.Value
	var i_222 int
	var variadic_QMARK__223 vm.Value
	var has_destructure_QMARK__224 vm.Value
	var flat_args_225 vm.Value
	var let_binds_226 vm.Value
	var remaining_227 vm.Value
	var x_228 vm.Value
	var gs_240 vm.Value
	var v268 vm.Value
	var rest_sym_241 vm.Value
	var fixed_args_242 vm.Value
	var amp_pos_243 int
	var args_vec_244 vm.Value
	var body_forms_245 vm.Value
	var i_246 int
	var variadic_QMARK__247 vm.Value
	var has_destructure_QMARK__248 vm.Value
	var flat_args_249 vm.Value
	var let_binds_250 vm.Value
	var remaining_251 vm.Value
	var x_252 vm.Value
	var gs_253 vm.Value
	var v271 vm.Value
	var rest_sym_254 vm.Value
	var fixed_args_255 vm.Value
	var amp_pos_256 int
	var args_vec_257 vm.Value
	var body_forms_258 vm.Value
	var i_259 int
	var variadic_QMARK__260 vm.Value
	var has_destructure_QMARK__261 vm.Value
	var flat_args_262 vm.Value
	var let_binds_263 vm.Value
	var remaining_264 vm.Value
	var x_265 vm.Value
	var gs_266 vm.Value
	var v300 vm.Value
	var binds_384 vm.Value
	var rest_sym_385 vm.Value
	var fixed_args_386 vm.Value
	var amp_pos_387 int
	var args_vec_388 vm.Value
	var body_forms_389 vm.Value
	var i_390 int
	var variadic_QMARK__391 vm.Value
	var has_destructure_QMARK__392 vm.Value
	var flat_args_393 vm.Value
	var let_binds_394 vm.Value
	var remaining_395 vm.Value
	var x_396 vm.Value
	var gs_397 vm.Value
	var v399 vm.Value
	var v401 vm.Value
	var v403 vm.Value
	var rest_sym_273 vm.Value
	var fixed_args_274 vm.Value
	var amp_pos_275 int
	var args_vec_276 vm.Value
	var body_forms_277 vm.Value
	var i_278 int
	var variadic_QMARK__279 vm.Value
	var has_destructure_QMARK__280 vm.Value
	var flat_args_281 vm.Value
	var let_binds_282 vm.Value
	var remaining_283 vm.Value
	var x_284 vm.Value
	var gs_285 vm.Value
	var v303 vm.Value
	var rest_sym_286 vm.Value
	var fixed_args_287 vm.Value
	var amp_pos_288 int
	var args_vec_289 vm.Value
	var body_forms_290 vm.Value
	var i_291 int
	var variadic_QMARK__292 vm.Value
	var has_destructure_QMARK__293 vm.Value
	var flat_args_294 vm.Value
	var let_binds_295 vm.Value
	var remaining_296 vm.Value
	var x_297 vm.Value
	var gs_298 vm.Value
	var v369 vm.Value
	var rest_sym_370 vm.Value
	var fixed_args_371 vm.Value
	var amp_pos_372 int
	var args_vec_373 vm.Value
	var body_forms_374 vm.Value
	var i_375 int
	var variadic_QMARK__376 vm.Value
	var has_destructure_QMARK__377 vm.Value
	var flat_args_378 vm.Value
	var let_binds_379 vm.Value
	var remaining_380 vm.Value
	var x_381 vm.Value
	var gs_382 vm.Value
	var rest_sym_305 vm.Value
	var fixed_args_306 vm.Value
	var amp_pos_307 int
	var args_vec_308 vm.Value
	var body_forms_309 vm.Value
	var i_310 int
	var variadic_QMARK__311 vm.Value
	var has_destructure_QMARK__312 vm.Value
	var flat_args_313 vm.Value
	var let_binds_314 vm.Value
	var remaining_315 vm.Value
	var x_316 vm.Value
	var gs_317 vm.Value
	var arg__6182_335 vm.Value
	var arg__6188_339 vm.Value
	var arg__6189_340 vm.Value
	var arg__6195_344 vm.Value
	var arg__6201_348 vm.Value
	var arg__6202_349 vm.Value
	var v350 vm.Value
	var rest_sym_318 vm.Value
	var fixed_args_319 vm.Value
	var amp_pos_320 int
	var args_vec_321 vm.Value
	var body_forms_322 vm.Value
	var i_323 int
	var variadic_QMARK__324 vm.Value
	var has_destructure_QMARK__325 vm.Value
	var flat_args_326 vm.Value
	var let_binds_327 vm.Value
	var remaining_328 vm.Value
	var x_329 vm.Value
	var gs_330 vm.Value
	var v354 vm.Value
	var rest_sym_355 vm.Value
	var fixed_args_356 vm.Value
	var amp_pos_357 int
	var args_vec_358 vm.Value
	var body_forms_359 vm.Value
	var i_360 int
	var variadic_QMARK__361 vm.Value
	var has_destructure_QMARK__362 vm.Value
	var flat_args_363 vm.Value
	var let_binds_364 vm.Value
	var remaining_365 vm.Value
	var x_366 vm.Value
	var gs_367 vm.Value
	var result_417 vm.Value
	var rest_sym_418 vm.Value
	var fixed_args_419 vm.Value
	var amp_pos_420 int
	var args_vec_421 vm.Value
	var body_forms_422 vm.Value
	var i_423 int
	var variadic_QMARK__424 vm.Value
	var has_destructure_QMARK__425 vm.Value
	var flat_args_426 vm.Value
	var let_binds_427 vm.Value
	var remaining_428 vm.Value
	var arg__6218_443 vm.Value
	var arg__6223_446 vm.Value
	var v447 vm.Value
	var result_429 vm.Value
	var rest_sym_430 vm.Value
	var fixed_args_431 vm.Value
	var amp_pos_432 int
	var args_vec_433 vm.Value
	var body_forms_434 vm.Value
	var i_435 int
	var variadic_QMARK__436 vm.Value
	var has_destructure_QMARK__437 vm.Value
	var flat_args_438 vm.Value
	var let_binds_439 vm.Value
	var remaining_440 vm.Value
	var v450 vm.Value
	var final_flat_args_452 vm.Value
	var result_453 vm.Value
	var rest_sym_454 vm.Value
	var fixed_args_455 vm.Value
	var amp_pos_456 int
	var args_vec_457 vm.Value
	var body_forms_458 vm.Value
	var i_459 int
	var variadic_QMARK__460 vm.Value
	var has_destructure_QMARK__461 vm.Value
	var flat_args_462 vm.Value
	var let_binds_463 vm.Value
	var remaining_464 vm.Value
	var arg__6229_492 vm.Value
	var arg__6233_495 vm.Value
	var v496 vm.Value
	var final_flat_args_465 vm.Value
	var result_466 vm.Value
	var rest_sym_467 vm.Value
	var fixed_args_468 vm.Value
	var amp_pos_469 int
	var args_vec_470 vm.Value
	var body_forms_471 vm.Value
	var i_472 int
	var variadic_QMARK__473 vm.Value
	var has_destructure_QMARK__474 vm.Value
	var flat_args_475 vm.Value
	var let_binds_476 vm.Value
	var remaining_477 vm.Value
	var arg__6238_501 vm.Value
	var arg__6242_504 vm.Value
	var arg__6243_505 vm.Value
	var arg__6250_510 vm.Value
	var arg__6254_513 vm.Value
	var arg__6255_514 vm.Value
	var arg__6257_515 vm.Value
	var arg__6263_520 vm.Value
	var arg__6267_523 vm.Value
	var arg__6268_524 vm.Value
	var arg__6275_529 vm.Value
	var arg__6279_532 vm.Value
	var arg__6280_533 vm.Value
	var arg__6282_534 vm.Value
	var v535 vm.Value
	var final_flat_args_478 vm.Value
	var result_479 vm.Value
	var rest_sym_480 vm.Value
	var fixed_args_481 vm.Value
	var amp_pos_482 int
	var args_vec_483 vm.Value
	var body_forms_484 vm.Value
	var i_485 int
	var variadic_QMARK__486 vm.Value
	var has_destructure_QMARK__487 vm.Value
	var flat_args_488 vm.Value
	var let_binds_489 vm.Value
	var remaining_490 vm.Value
	var body_538 vm.Value
	var final_flat_args_539 vm.Value
	var result_540 vm.Value
	var rest_sym_541 vm.Value
	var fixed_args_542 vm.Value
	var amp_pos_543 int
	var args_vec_544 vm.Value
	var body_forms_545 vm.Value
	var i_546 int
	var variadic_QMARK__547 vm.Value
	var has_destructure_QMARK__548 vm.Value
	var flat_args_549 vm.Value
	var let_binds_550 vm.Value
	var remaining_551 vm.Value
	var v556 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = i_3, remaining_4, v572, v16, args_vec_7, body_forms_8, i_9, remaining_10, v575, args_vec_11, body_forms_12, i_13, remaining_14, v573, arg__6089_29, v31, amp_pos_44, args_vec_45, body_forms_46, i_47, remaining_48, variadic_QMARK__52, args_vec_20, body_forms_21, i_22, remaining_23, v576, args_vec_24, body_forms_25, i_26, remaining_27, v574, v34, v36, v38, args_vec_39, body_forms_40, i_41, remaining_42, amp_pos_53, args_vec_54, body_forms_55, i_56, remaining_57, variadic_QMARK__58, v67, amp_pos_59, args_vec_60, body_forms_61, i_62, remaining_63, variadic_QMARK__64, fixed_args_70, amp_pos_71, args_vec_72, body_forms_73, i_74, remaining_75, variadic_QMARK__76, fixed_args_77, amp_pos_78, args_vec_79, body_forms_80, i_81, remaining_82, variadic_QMARK__83, arg__6107_92, v97, fixed_args_84, amp_pos_85, args_vec_86, body_forms_87, i_88, remaining_89, variadic_QMARK__90, rest_sym_101, fixed_args_102, amp_pos_103, args_vec_104, body_forms_105, i_106, remaining_107, variadic_QMARK__108, has_destructure_QMARK__112, rest_sym_113, fixed_args_114, amp_pos_115, args_vec_116, body_forms_117, i_118, remaining_119, variadic_QMARK__120, has_destructure_QMARK__121, rest_sym_122, fixed_args_123, amp_pos_124, args_vec_125, body_forms_126, i_127, remaining_128, variadic_QMARK__129, has_destructure_QMARK__130, v560, rest_sym_561, fixed_args_562, amp_pos_563, args_vec_564, body_forms_565, i_566, remaining_567, variadic_QMARK__568, has_destructure_QMARK__569, rest_sym_131, fixed_args_132, amp_pos_133, args_vec_134, body_forms_135, i_136, remaining_137, or__x_138, variadic_QMARK__139, has_destructure_QMARK__140, rest_sym_141, fixed_args_142, amp_pos_143, args_vec_144, body_forms_145, i_146, remaining_147, or__x_148, variadic_QMARK__149, has_destructure_QMARK__150, v154, rest_sym_155, fixed_args_156, amp_pos_157, args_vec_158, body_forms_159, i_160, remaining_161, or__x_162, variadic_QMARK__163, has_destructure_QMARK__164, remaining_166, flat_args_167, let_binds_168, remaining_169, v196, rest_sym_173, fixed_args_174, amp_pos_175, args_vec_176, body_forms_177, i_178, variadic_QMARK__179, has_destructure_QMARK__180, flat_args_181, let_binds_182, remaining_183, v201, rest_sym_184, fixed_args_185, amp_pos_186, args_vec_187, body_forms_188, i_189, variadic_QMARK__190, has_destructure_QMARK__191, flat_args_192, let_binds_193, remaining_194, x_204, v230, result_405, rest_sym_406, fixed_args_407, amp_pos_408, args_vec_409, body_forms_410, i_411, variadic_QMARK__412, has_destructure_QMARK__413, flat_args_414, let_binds_415, remaining_416, rest_sym_205, fixed_args_206, amp_pos_207, args_vec_208, body_forms_209, i_210, variadic_QMARK__211, has_destructure_QMARK__212, flat_args_213, let_binds_214, remaining_215, x_216, v233, v235, rest_sym_217, fixed_args_218, amp_pos_219, args_vec_220, body_forms_221, i_222, variadic_QMARK__223, has_destructure_QMARK__224, flat_args_225, let_binds_226, remaining_227, x_228, gs_240, v268, rest_sym_241, fixed_args_242, amp_pos_243, args_vec_244, body_forms_245, i_246, variadic_QMARK__247, has_destructure_QMARK__248, flat_args_249, let_binds_250, remaining_251, x_252, gs_253, v271, rest_sym_254, fixed_args_255, amp_pos_256, args_vec_257, body_forms_258, i_259, variadic_QMARK__260, has_destructure_QMARK__261, flat_args_262, let_binds_263, remaining_264, x_265, gs_266, v300, binds_384, rest_sym_385, fixed_args_386, amp_pos_387, args_vec_388, body_forms_389, i_390, variadic_QMARK__391, has_destructure_QMARK__392, flat_args_393, let_binds_394, remaining_395, x_396, gs_397, v399, v401, v403, rest_sym_273, fixed_args_274, amp_pos_275, args_vec_276, body_forms_277, i_278, variadic_QMARK__279, has_destructure_QMARK__280, flat_args_281, let_binds_282, remaining_283, x_284, gs_285, v303, rest_sym_286, fixed_args_287, amp_pos_288, args_vec_289, body_forms_290, i_291, variadic_QMARK__292, has_destructure_QMARK__293, flat_args_294, let_binds_295, remaining_296, x_297, gs_298, v369, rest_sym_370, fixed_args_371, amp_pos_372, args_vec_373, body_forms_374, i_375, variadic_QMARK__376, has_destructure_QMARK__377, flat_args_378, let_binds_379, remaining_380, x_381, gs_382, rest_sym_305, fixed_args_306, amp_pos_307, args_vec_308, body_forms_309, i_310, variadic_QMARK__311, has_destructure_QMARK__312, flat_args_313, let_binds_314, remaining_315, x_316, gs_317, arg__6182_335, arg__6188_339, arg__6189_340, arg__6195_344, arg__6201_348, arg__6202_349, v350, rest_sym_318, fixed_args_319, amp_pos_320, args_vec_321, body_forms_322, i_323, variadic_QMARK__324, has_destructure_QMARK__325, flat_args_326, let_binds_327, remaining_328, x_329, gs_330, v354, rest_sym_355, fixed_args_356, amp_pos_357, args_vec_358, body_forms_359, i_360, variadic_QMARK__361, has_destructure_QMARK__362, flat_args_363, let_binds_364, remaining_365, x_366, gs_367, result_417, rest_sym_418, fixed_args_419, amp_pos_420, args_vec_421, body_forms_422, i_423, variadic_QMARK__424, has_destructure_QMARK__425, flat_args_426, let_binds_427, remaining_428, arg__6218_443, arg__6223_446, v447, result_429, rest_sym_430, fixed_args_431, amp_pos_432, args_vec_433, body_forms_434, i_435, variadic_QMARK__436, has_destructure_QMARK__437, flat_args_438, let_binds_439, remaining_440, v450, final_flat_args_452, result_453, rest_sym_454, fixed_args_455, amp_pos_456, args_vec_457, body_forms_458, i_459, variadic_QMARK__460, has_destructure_QMARK__461, flat_args_462, let_binds_463, remaining_464, arg__6229_492, arg__6233_495, v496, final_flat_args_465, result_466, rest_sym_467, fixed_args_468, amp_pos_469, args_vec_470, body_forms_471, i_472, variadic_QMARK__473, has_destructure_QMARK__474, flat_args_475, let_binds_476, remaining_477, arg__6238_501, arg__6242_504, arg__6243_505, arg__6250_510, arg__6254_513, arg__6255_514, arg__6257_515, arg__6263_520, arg__6267_523, arg__6268_524, arg__6275_529, arg__6279_532, arg__6280_533, arg__6282_534, v535, final_flat_args_478, result_479, rest_sym_480, fixed_args_481, amp_pos_482, args_vec_483, body_forms_484, i_485, variadic_QMARK__486, has_destructure_QMARK__487, flat_args_488, let_binds_489, remaining_490, body_538, final_flat_args_539, result_540, rest_sym_541, fixed_args_542, amp_pos_543, args_vec_544, body_forms_545, i_546, variadic_QMARK__547, has_destructure_QMARK__548, flat_args_549, let_binds_550, remaining_551, v556
	i_3 = 0
	remaining_4 = arg0
	v572 = vm.Symbol("&")
	goto b1
b1:
	;
	v16, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "empty?").Deref(), []vm.Value{remaining_4})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v16) {
		args_vec_7 = arg0
		body_forms_8 = arg1
		i_9 = i_3
		remaining_10 = remaining_4
		v575 = v572
		goto b2
	} else {
		args_vec_11 = arg0
		body_forms_12 = arg1
		i_13 = i_3
		remaining_14 = remaining_4
		v573 = v572
		goto b3
	}
b2:
	;
	amp_pos_44 = -1
	args_vec_45 = args_vec_7
	body_forms_46 = body_forms_8
	i_47 = i_9
	remaining_48 = remaining_10
	goto b4
b3:
	;
	arg__6089_29, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{remaining_14})
	if callErr != nil {
		return nil, callErr
	}
	v31 = arg__6089_29 == v573
	if v31 {
		args_vec_20 = args_vec_11
		body_forms_21 = body_forms_12
		i_22 = i_13
		remaining_23 = remaining_14
		v576 = v573
		goto b5
	} else {
		args_vec_24 = args_vec_11
		body_forms_25 = body_forms_12
		i_26 = i_13
		remaining_27 = remaining_14
		v574 = v573
		goto b6
	}
b4:
	;
	variadic_QMARK__52, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not=").Deref(), []vm.Value{vm.Int(amp_pos_44), vm.Int(-1)})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(variadic_QMARK__52) {
		amp_pos_53 = amp_pos_44
		args_vec_54 = args_vec_45
		body_forms_55 = body_forms_46
		i_56 = i_47
		remaining_57 = remaining_48
		variadic_QMARK__58 = variadic_QMARK__52
		goto b8
	} else {
		amp_pos_59 = amp_pos_44
		args_vec_60 = args_vec_45
		body_forms_61 = body_forms_46
		i_62 = i_47
		remaining_63 = remaining_48
		variadic_QMARK__64 = variadic_QMARK__52
		goto b9
	}
b5:
	;
	v38 = i_22
	args_vec_39 = args_vec_20
	body_forms_40 = body_forms_21
	i_41 = i_22
	remaining_42 = remaining_23
	goto b7
b6:
	;
	v34 = i_26 + 1
	v36, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{remaining_27})
	if callErr != nil {
		return nil, callErr
	}
	i_3 = v34
	remaining_4 = v36
	v572 = v574
	goto b1
b7:
	;
	amp_pos_44 = v38
	args_vec_45 = args_vec_39
	body_forms_46 = body_forms_40
	i_47 = i_41
	remaining_48 = remaining_42
	goto b4
b8:
	;
	v67, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "take").Deref(), []vm.Value{vm.Int(amp_pos_53), args_vec_54})
	if callErr != nil {
		return nil, callErr
	}
	fixed_args_70 = v67
	amp_pos_71 = amp_pos_53
	args_vec_72 = args_vec_54
	body_forms_73 = body_forms_55
	i_74 = i_56
	remaining_75 = remaining_57
	variadic_QMARK__76 = variadic_QMARK__58
	goto b10
b9:
	;
	fixed_args_70 = args_vec_60
	amp_pos_71 = amp_pos_59
	args_vec_72 = args_vec_60
	body_forms_73 = body_forms_61
	i_74 = i_62
	remaining_75 = remaining_63
	variadic_QMARK__76 = variadic_QMARK__64
	goto b10
b10:
	;
	if vm.IsTruthy(variadic_QMARK__76) {
		fixed_args_77 = fixed_args_70
		amp_pos_78 = amp_pos_71
		args_vec_79 = args_vec_72
		body_forms_80 = body_forms_73
		i_81 = i_74
		remaining_82 = remaining_75
		variadic_QMARK__83 = variadic_QMARK__76
		goto b11
	} else {
		fixed_args_84 = fixed_args_70
		amp_pos_85 = amp_pos_71
		args_vec_86 = args_vec_72
		body_forms_87 = body_forms_73
		i_88 = i_74
		remaining_89 = remaining_75
		variadic_QMARK__90 = variadic_QMARK__76
		goto b12
	}
b11:
	;
	arg__6107_92 = amp_pos_78 + 1
	v97, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{args_vec_79, vm.Int(arg__6107_92), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	rest_sym_101 = v97
	fixed_args_102 = fixed_args_77
	amp_pos_103 = amp_pos_78
	args_vec_104 = args_vec_79
	body_forms_105 = body_forms_80
	i_106 = i_81
	remaining_107 = remaining_82
	variadic_QMARK__108 = variadic_QMARK__83
	goto b13
b12:
	;
	rest_sym_101 = vm.NIL
	fixed_args_102 = fixed_args_84
	amp_pos_103 = amp_pos_85
	args_vec_104 = args_vec_86
	body_forms_105 = body_forms_87
	i_106 = i_88
	remaining_107 = remaining_89
	variadic_QMARK__108 = variadic_QMARK__90
	goto b13
b13:
	;
	has_destructure_QMARK__112, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "some").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var arg__6129_2 vm.Value
		var arg__6134_5 vm.Value
		var v6 vm.Value
		var callErr error
		_, _, _ = arg__6129_2, arg__6134_5, v6
		arg__6129_2, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "symbol?").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		arg__6134_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "symbol?").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		v6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not").Deref(), []vm.Value{arg__6134_5})
		if callErr != nil {
			return nil, callErr
		}
		return v6, nil
	}), fixed_args_102})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(variadic_QMARK__108) {
		rest_sym_131 = rest_sym_101
		fixed_args_132 = fixed_args_102
		amp_pos_133 = amp_pos_103
		args_vec_134 = args_vec_104
		body_forms_135 = body_forms_105
		i_136 = i_106
		remaining_137 = remaining_107
		or__x_138 = variadic_QMARK__108
		variadic_QMARK__139 = variadic_QMARK__108
		has_destructure_QMARK__140 = has_destructure_QMARK__112
		goto b17
	} else {
		rest_sym_141 = rest_sym_101
		fixed_args_142 = fixed_args_102
		amp_pos_143 = amp_pos_103
		args_vec_144 = args_vec_104
		body_forms_145 = body_forms_105
		i_146 = i_106
		remaining_147 = remaining_107
		or__x_148 = variadic_QMARK__108
		variadic_QMARK__149 = variadic_QMARK__108
		has_destructure_QMARK__150 = has_destructure_QMARK__112
		goto b18
	}
b14:
	;
	remaining_166 = fixed_args_114
	flat_args_167 = vm.NewArrayVector([]vm.Value{})
	let_binds_168 = vm.NewArrayVector([]vm.Value{})
	remaining_169 = remaining_119
	goto b20
b15:
	;
	v560 = vm.NIL
	rest_sym_561 = rest_sym_122
	fixed_args_562 = fixed_args_123
	amp_pos_563 = amp_pos_124
	args_vec_564 = args_vec_125
	body_forms_565 = body_forms_126
	i_566 = i_127
	remaining_567 = remaining_128
	variadic_QMARK__568 = variadic_QMARK__129
	has_destructure_QMARK__569 = has_destructure_QMARK__130
	goto b16
b16:
	;
	return v560, nil
b17:
	;
	v154 = or__x_138
	rest_sym_155 = rest_sym_131
	fixed_args_156 = fixed_args_132
	amp_pos_157 = amp_pos_133
	args_vec_158 = args_vec_134
	body_forms_159 = body_forms_135
	i_160 = i_136
	remaining_161 = remaining_137
	or__x_162 = or__x_138
	variadic_QMARK__163 = variadic_QMARK__139
	has_destructure_QMARK__164 = has_destructure_QMARK__140
	goto b19
b18:
	;
	v154 = has_destructure_QMARK__150
	rest_sym_155 = rest_sym_141
	fixed_args_156 = fixed_args_142
	amp_pos_157 = amp_pos_143
	args_vec_158 = args_vec_144
	body_forms_159 = body_forms_145
	i_160 = i_146
	remaining_161 = remaining_147
	or__x_162 = or__x_148
	variadic_QMARK__163 = variadic_QMARK__149
	has_destructure_QMARK__164 = has_destructure_QMARK__150
	goto b19
b19:
	;
	if vm.IsTruthy(v154) {
		rest_sym_113 = rest_sym_155
		fixed_args_114 = fixed_args_156
		amp_pos_115 = amp_pos_157
		args_vec_116 = args_vec_158
		body_forms_117 = body_forms_159
		i_118 = i_160
		remaining_119 = remaining_161
		variadic_QMARK__120 = variadic_QMARK__163
		has_destructure_QMARK__121 = has_destructure_QMARK__164
		goto b14
	} else {
		rest_sym_122 = rest_sym_155
		fixed_args_123 = fixed_args_156
		amp_pos_124 = amp_pos_157
		args_vec_125 = args_vec_158
		body_forms_126 = body_forms_159
		i_127 = i_160
		remaining_128 = remaining_161
		variadic_QMARK__129 = variadic_QMARK__163
		has_destructure_QMARK__130 = has_destructure_QMARK__164
		goto b15
	}
b20:
	;
	v196, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "empty?").Deref(), []vm.Value{remaining_169})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v196) {
		rest_sym_173 = rest_sym_113
		fixed_args_174 = fixed_args_114
		amp_pos_175 = amp_pos_115
		args_vec_176 = args_vec_116
		body_forms_177 = body_forms_117
		i_178 = i_118
		variadic_QMARK__179 = variadic_QMARK__120
		has_destructure_QMARK__180 = has_destructure_QMARK__121
		flat_args_181 = flat_args_167
		let_binds_182 = let_binds_168
		remaining_183 = remaining_169
		goto b21
	} else {
		rest_sym_184 = rest_sym_113
		fixed_args_185 = fixed_args_114
		amp_pos_186 = amp_pos_115
		args_vec_187 = args_vec_116
		body_forms_188 = body_forms_117
		i_189 = i_118
		variadic_QMARK__190 = variadic_QMARK__120
		has_destructure_QMARK__191 = has_destructure_QMARK__121
		flat_args_192 = flat_args_167
		let_binds_193 = let_binds_168
		remaining_194 = remaining_169
		goto b22
	}
b21:
	;
	v201, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("flat-args"), flat_args_181, vm.Keyword("let-binds"), let_binds_182})
	if callErr != nil {
		return nil, callErr
	}
	result_405 = v201
	rest_sym_406 = rest_sym_173
	fixed_args_407 = fixed_args_174
	amp_pos_408 = amp_pos_175
	args_vec_409 = args_vec_176
	body_forms_410 = body_forms_177
	i_411 = i_178
	variadic_QMARK__412 = variadic_QMARK__179
	has_destructure_QMARK__413 = has_destructure_QMARK__180
	flat_args_414 = flat_args_181
	let_binds_415 = let_binds_182
	remaining_416 = remaining_183
	goto b23
b22:
	;
	x_204, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{remaining_194})
	if callErr != nil {
		return nil, callErr
	}
	v230, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "symbol?").Deref(), []vm.Value{x_204})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v230) {
		rest_sym_205 = rest_sym_184
		fixed_args_206 = fixed_args_185
		amp_pos_207 = amp_pos_186
		args_vec_208 = args_vec_187
		body_forms_209 = body_forms_188
		i_210 = i_189
		variadic_QMARK__211 = variadic_QMARK__190
		has_destructure_QMARK__212 = has_destructure_QMARK__191
		flat_args_213 = flat_args_192
		let_binds_214 = let_binds_193
		remaining_215 = remaining_194
		x_216 = x_204
		goto b24
	} else {
		rest_sym_217 = rest_sym_184
		fixed_args_218 = fixed_args_185
		amp_pos_219 = amp_pos_186
		args_vec_220 = args_vec_187
		body_forms_221 = body_forms_188
		i_222 = i_189
		variadic_QMARK__223 = variadic_QMARK__190
		has_destructure_QMARK__224 = has_destructure_QMARK__191
		flat_args_225 = flat_args_192
		let_binds_226 = let_binds_193
		remaining_227 = remaining_194
		x_228 = x_204
		goto b25
	}
b23:
	;
	if vm.IsTruthy(variadic_QMARK__412) {
		result_417 = result_405
		rest_sym_418 = rest_sym_406
		fixed_args_419 = fixed_args_407
		amp_pos_420 = amp_pos_408
		args_vec_421 = args_vec_409
		body_forms_422 = body_forms_410
		i_423 = i_411
		variadic_QMARK__424 = variadic_QMARK__412
		has_destructure_QMARK__425 = has_destructure_QMARK__413
		flat_args_426 = flat_args_414
		let_binds_427 = let_binds_415
		remaining_428 = remaining_416
		goto b36
	} else {
		result_429 = result_405
		rest_sym_430 = rest_sym_406
		fixed_args_431 = fixed_args_407
		amp_pos_432 = amp_pos_408
		args_vec_433 = args_vec_409
		body_forms_434 = body_forms_410
		i_435 = i_411
		variadic_QMARK__436 = variadic_QMARK__412
		has_destructure_QMARK__437 = has_destructure_QMARK__413
		flat_args_438 = flat_args_414
		let_binds_439 = let_binds_415
		remaining_440 = remaining_416
		goto b37
	}
b24:
	;
	v233, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{remaining_215})
	if callErr != nil {
		return nil, callErr
	}
	v235, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{flat_args_213, x_216})
	if callErr != nil {
		return nil, callErr
	}
	remaining_166 = v233
	flat_args_167 = v235
	let_binds_168 = let_binds_214
	remaining_169 = remaining_215
	goto b20
b25:
	;
	gs_240, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "gensym").Deref(), []vm.Value{vm.String("p__")})
	if callErr != nil {
		return nil, callErr
	}
	v268, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector?").Deref(), []vm.Value{x_228})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v268) {
		rest_sym_241 = rest_sym_217
		fixed_args_242 = fixed_args_218
		amp_pos_243 = amp_pos_219
		args_vec_244 = args_vec_220
		body_forms_245 = body_forms_221
		i_246 = i_222
		variadic_QMARK__247 = variadic_QMARK__223
		has_destructure_QMARK__248 = has_destructure_QMARK__224
		flat_args_249 = flat_args_225
		let_binds_250 = let_binds_226
		remaining_251 = remaining_227
		x_252 = x_228
		gs_253 = gs_240
		goto b27
	} else {
		rest_sym_254 = rest_sym_217
		fixed_args_255 = fixed_args_218
		amp_pos_256 = amp_pos_219
		args_vec_257 = args_vec_220
		body_forms_258 = body_forms_221
		i_259 = i_222
		variadic_QMARK__260 = variadic_QMARK__223
		has_destructure_QMARK__261 = has_destructure_QMARK__224
		flat_args_262 = flat_args_225
		let_binds_263 = let_binds_226
		remaining_264 = remaining_227
		x_265 = x_228
		gs_266 = gs_240
		goto b28
	}
b27:
	;
	v271, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "expand-vector-pattern").Deref(), []vm.Value{gs_253, x_252})
	if callErr != nil {
		return nil, callErr
	}
	binds_384 = v271
	rest_sym_385 = rest_sym_241
	fixed_args_386 = fixed_args_242
	amp_pos_387 = amp_pos_243
	args_vec_388 = args_vec_244
	body_forms_389 = body_forms_245
	i_390 = i_246
	variadic_QMARK__391 = variadic_QMARK__247
	has_destructure_QMARK__392 = has_destructure_QMARK__248
	flat_args_393 = flat_args_249
	let_binds_394 = let_binds_250
	remaining_395 = remaining_251
	x_396 = x_252
	gs_397 = gs_253
	goto b29
b28:
	;
	v300, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map?").Deref(), []vm.Value{x_265})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v300) {
		rest_sym_273 = rest_sym_254
		fixed_args_274 = fixed_args_255
		amp_pos_275 = amp_pos_256
		args_vec_276 = args_vec_257
		body_forms_277 = body_forms_258
		i_278 = i_259
		variadic_QMARK__279 = variadic_QMARK__260
		has_destructure_QMARK__280 = has_destructure_QMARK__261
		flat_args_281 = flat_args_262
		let_binds_282 = let_binds_263
		remaining_283 = remaining_264
		x_284 = x_265
		gs_285 = gs_266
		goto b30
	} else {
		rest_sym_286 = rest_sym_254
		fixed_args_287 = fixed_args_255
		amp_pos_288 = amp_pos_256
		args_vec_289 = args_vec_257
		body_forms_290 = body_forms_258
		i_291 = i_259
		variadic_QMARK__292 = variadic_QMARK__260
		has_destructure_QMARK__293 = has_destructure_QMARK__261
		flat_args_294 = flat_args_262
		let_binds_295 = let_binds_263
		remaining_296 = remaining_264
		x_297 = x_265
		gs_298 = gs_266
		goto b31
	}
b29:
	;
	v399, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{remaining_395})
	if callErr != nil {
		return nil, callErr
	}
	v401, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{flat_args_393, gs_397})
	if callErr != nil {
		return nil, callErr
	}
	v403, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{let_binds_394, binds_384})
	if callErr != nil {
		return nil, callErr
	}
	remaining_166 = v399
	flat_args_167 = v401
	let_binds_168 = v403
	remaining_169 = remaining_395
	goto b20
b30:
	;
	v303, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "expand-map-pattern").Deref(), []vm.Value{gs_285, x_284})
	if callErr != nil {
		return nil, callErr
	}
	v369 = v303
	rest_sym_370 = rest_sym_273
	fixed_args_371 = fixed_args_274
	amp_pos_372 = amp_pos_275
	args_vec_373 = args_vec_276
	body_forms_374 = body_forms_277
	i_375 = i_278
	variadic_QMARK__376 = variadic_QMARK__279
	has_destructure_QMARK__377 = has_destructure_QMARK__280
	flat_args_378 = flat_args_281
	let_binds_379 = let_binds_282
	remaining_380 = remaining_283
	x_381 = x_284
	gs_382 = gs_285
	goto b32
b31:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		rest_sym_305 = rest_sym_286
		fixed_args_306 = fixed_args_287
		amp_pos_307 = amp_pos_288
		args_vec_308 = args_vec_289
		body_forms_309 = body_forms_290
		i_310 = i_291
		variadic_QMARK__311 = variadic_QMARK__292
		has_destructure_QMARK__312 = has_destructure_QMARK__293
		flat_args_313 = flat_args_294
		let_binds_314 = let_binds_295
		remaining_315 = remaining_296
		x_316 = x_297
		gs_317 = gs_298
		goto b33
	} else {
		rest_sym_318 = rest_sym_286
		fixed_args_319 = fixed_args_287
		amp_pos_320 = amp_pos_288
		args_vec_321 = args_vec_289
		body_forms_322 = body_forms_290
		i_323 = i_291
		variadic_QMARK__324 = variadic_QMARK__292
		has_destructure_QMARK__325 = has_destructure_QMARK__293
		flat_args_326 = flat_args_294
		let_binds_327 = let_binds_295
		remaining_328 = remaining_296
		x_329 = x_297
		gs_330 = gs_298
		goto b34
	}
b32:
	;
	binds_384 = v369
	rest_sym_385 = rest_sym_370
	fixed_args_386 = fixed_args_371
	amp_pos_387 = amp_pos_372
	args_vec_388 = args_vec_373
	body_forms_389 = body_forms_374
	i_390 = i_375
	variadic_QMARK__391 = variadic_QMARK__376
	has_destructure_QMARK__392 = has_destructure_QMARK__377
	flat_args_393 = flat_args_378
	let_binds_394 = let_binds_379
	remaining_395 = remaining_380
	x_396 = x_381
	gs_397 = gs_382
	goto b29
b33:
	;
	arg__6182_335, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "pr-str").Deref(), []vm.Value{x_316})
	if callErr != nil {
		return nil, callErr
	}
	arg__6188_339, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "pr-str").Deref(), []vm.Value{x_316})
	if callErr != nil {
		return nil, callErr
	}
	arg__6189_340, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("unsupported arg pattern: "), arg__6188_339})
	if callErr != nil {
		return nil, callErr
	}
	arg__6195_344, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "pr-str").Deref(), []vm.Value{x_316})
	if callErr != nil {
		return nil, callErr
	}
	arg__6201_348, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "pr-str").Deref(), []vm.Value{x_316})
	if callErr != nil {
		return nil, callErr
	}
	arg__6202_349, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("unsupported arg pattern: "), arg__6201_348})
	if callErr != nil {
		return nil, callErr
	}
	v350, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "throw").Deref(), []vm.Value{arg__6202_349})
	if callErr != nil {
		return nil, callErr
	}
	v354 = v350
	rest_sym_355 = rest_sym_305
	fixed_args_356 = fixed_args_306
	amp_pos_357 = amp_pos_307
	args_vec_358 = args_vec_308
	body_forms_359 = body_forms_309
	i_360 = i_310
	variadic_QMARK__361 = variadic_QMARK__311
	has_destructure_QMARK__362 = has_destructure_QMARK__312
	flat_args_363 = flat_args_313
	let_binds_364 = let_binds_314
	remaining_365 = remaining_315
	x_366 = x_316
	gs_367 = gs_317
	goto b35
b34:
	;
	v354 = vm.NIL
	rest_sym_355 = rest_sym_318
	fixed_args_356 = fixed_args_319
	amp_pos_357 = amp_pos_320
	args_vec_358 = args_vec_321
	body_forms_359 = body_forms_322
	i_360 = i_323
	variadic_QMARK__361 = variadic_QMARK__324
	has_destructure_QMARK__362 = has_destructure_QMARK__325
	flat_args_363 = flat_args_326
	let_binds_364 = let_binds_327
	remaining_365 = remaining_328
	x_366 = x_329
	gs_367 = gs_330
	goto b35
b35:
	;
	v369 = v354
	rest_sym_370 = rest_sym_355
	fixed_args_371 = fixed_args_356
	amp_pos_372 = amp_pos_357
	args_vec_373 = args_vec_358
	body_forms_374 = body_forms_359
	i_375 = i_360
	variadic_QMARK__376 = variadic_QMARK__361
	has_destructure_QMARK__377 = has_destructure_QMARK__362
	flat_args_378 = flat_args_363
	let_binds_379 = let_binds_364
	remaining_380 = remaining_365
	x_381 = x_366
	gs_382 = gs_367
	goto b32
b36:
	;
	arg__6218_443, callErr = rt.InvokeValue(vm.Keyword("flat-args"), []vm.Value{result_417})
	if callErr != nil {
		return nil, callErr
	}
	arg__6223_446, callErr = rt.InvokeValue(vm.Keyword("flat-args"), []vm.Value{result_417})
	if callErr != nil {
		return nil, callErr
	}
	v447, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{arg__6223_446, rest_sym_418})
	if callErr != nil {
		return nil, callErr
	}
	final_flat_args_452 = v447
	result_453 = result_417
	rest_sym_454 = rest_sym_418
	fixed_args_455 = fixed_args_419
	amp_pos_456 = amp_pos_420
	args_vec_457 = args_vec_421
	body_forms_458 = body_forms_422
	i_459 = i_423
	variadic_QMARK__460 = variadic_QMARK__424
	has_destructure_QMARK__461 = has_destructure_QMARK__425
	flat_args_462 = flat_args_426
	let_binds_463 = let_binds_427
	remaining_464 = remaining_428
	goto b38
b37:
	;
	v450, callErr = rt.InvokeValue(vm.Keyword("flat-args"), []vm.Value{result_429})
	if callErr != nil {
		return nil, callErr
	}
	final_flat_args_452 = v450
	result_453 = result_429
	rest_sym_454 = rest_sym_430
	fixed_args_455 = fixed_args_431
	amp_pos_456 = amp_pos_432
	args_vec_457 = args_vec_433
	body_forms_458 = body_forms_434
	i_459 = i_435
	variadic_QMARK__460 = variadic_QMARK__436
	has_destructure_QMARK__461 = has_destructure_QMARK__437
	flat_args_462 = flat_args_438
	let_binds_463 = let_binds_439
	remaining_464 = remaining_440
	goto b38
b38:
	;
	arg__6229_492, callErr = rt.InvokeValue(vm.Keyword("let-binds"), []vm.Value{result_453})
	if callErr != nil {
		return nil, callErr
	}
	arg__6233_495, callErr = rt.InvokeValue(vm.Keyword("let-binds"), []vm.Value{result_453})
	if callErr != nil {
		return nil, callErr
	}
	v496, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{arg__6233_495})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v496) {
		final_flat_args_465 = final_flat_args_452
		result_466 = result_453
		rest_sym_467 = rest_sym_454
		fixed_args_468 = fixed_args_455
		amp_pos_469 = amp_pos_456
		args_vec_470 = args_vec_457
		body_forms_471 = body_forms_458
		i_472 = i_459
		variadic_QMARK__473 = variadic_QMARK__460
		has_destructure_QMARK__474 = has_destructure_QMARK__461
		flat_args_475 = flat_args_462
		let_binds_476 = let_binds_463
		remaining_477 = remaining_464
		goto b39
	} else {
		final_flat_args_478 = final_flat_args_452
		result_479 = result_453
		rest_sym_480 = rest_sym_454
		fixed_args_481 = fixed_args_455
		amp_pos_482 = amp_pos_456
		args_vec_483 = args_vec_457
		body_forms_484 = body_forms_458
		i_485 = i_459
		variadic_QMARK__486 = variadic_QMARK__460
		has_destructure_QMARK__487 = has_destructure_QMARK__461
		flat_args_488 = flat_args_462
		let_binds_489 = let_binds_463
		remaining_490 = remaining_464
		goto b40
	}
b39:
	;
	arg__6238_501, callErr = rt.InvokeValue(vm.Keyword("let-binds"), []vm.Value{result_466})
	if callErr != nil {
		return nil, callErr
	}
	arg__6242_504, callErr = rt.InvokeValue(vm.Keyword("let-binds"), []vm.Value{result_466})
	if callErr != nil {
		return nil, callErr
	}
	arg__6243_505, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg__6242_504})
	if callErr != nil {
		return nil, callErr
	}
	arg__6250_510, callErr = rt.InvokeValue(vm.Keyword("let-binds"), []vm.Value{result_466})
	if callErr != nil {
		return nil, callErr
	}
	arg__6254_513, callErr = rt.InvokeValue(vm.Keyword("let-binds"), []vm.Value{result_466})
	if callErr != nil {
		return nil, callErr
	}
	arg__6255_514, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg__6254_513})
	if callErr != nil {
		return nil, callErr
	}
	arg__6257_515, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "apply").Deref(), []vm.Value{rt.LookupVar("clojure.core", "list").Deref(), vm.Symbol("let*"), arg__6255_514, body_forms_471})
	if callErr != nil {
		return nil, callErr
	}
	arg__6263_520, callErr = rt.InvokeValue(vm.Keyword("let-binds"), []vm.Value{result_466})
	if callErr != nil {
		return nil, callErr
	}
	arg__6267_523, callErr = rt.InvokeValue(vm.Keyword("let-binds"), []vm.Value{result_466})
	if callErr != nil {
		return nil, callErr
	}
	arg__6268_524, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg__6267_523})
	if callErr != nil {
		return nil, callErr
	}
	arg__6275_529, callErr = rt.InvokeValue(vm.Keyword("let-binds"), []vm.Value{result_466})
	if callErr != nil {
		return nil, callErr
	}
	arg__6279_532, callErr = rt.InvokeValue(vm.Keyword("let-binds"), []vm.Value{result_466})
	if callErr != nil {
		return nil, callErr
	}
	arg__6280_533, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg__6279_532})
	if callErr != nil {
		return nil, callErr
	}
	arg__6282_534, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "apply").Deref(), []vm.Value{rt.LookupVar("clojure.core", "list").Deref(), vm.Symbol("let*"), arg__6280_533, body_forms_471})
	if callErr != nil {
		return nil, callErr
	}
	v535, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "list").Deref(), []vm.Value{arg__6282_534})
	if callErr != nil {
		return nil, callErr
	}
	body_538 = v535
	final_flat_args_539 = final_flat_args_465
	result_540 = result_466
	rest_sym_541 = rest_sym_467
	fixed_args_542 = fixed_args_468
	amp_pos_543 = amp_pos_469
	args_vec_544 = args_vec_470
	body_forms_545 = body_forms_471
	i_546 = i_472
	variadic_QMARK__547 = variadic_QMARK__473
	has_destructure_QMARK__548 = has_destructure_QMARK__474
	flat_args_549 = flat_args_475
	let_binds_550 = let_binds_476
	remaining_551 = remaining_477
	goto b41
b40:
	;
	body_538 = body_forms_484
	final_flat_args_539 = final_flat_args_478
	result_540 = result_479
	rest_sym_541 = rest_sym_480
	fixed_args_542 = fixed_args_481
	amp_pos_543 = amp_pos_482
	args_vec_544 = args_vec_483
	body_forms_545 = body_forms_484
	i_546 = i_485
	variadic_QMARK__547 = variadic_QMARK__486
	has_destructure_QMARK__548 = has_destructure_QMARK__487
	flat_args_549 = flat_args_488
	let_binds_550 = let_binds_489
	remaining_551 = remaining_490
	goto b41
b41:
	;
	v556, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("variadic?"), variadic_QMARK__547, vm.Keyword("flat-args"), final_flat_args_539, vm.Keyword("body"), body_538})
	if callErr != nil {
		return nil, callErr
	}
	v560 = v556
	rest_sym_561 = rest_sym_541
	fixed_args_562 = fixed_args_542
	amp_pos_563 = amp_pos_543
	args_vec_564 = args_vec_544
	body_forms_565 = body_forms_545
	i_566 = i_546
	remaining_567 = remaining_551
	variadic_QMARK__568 = variadic_QMARK__547
	has_destructure_QMARK__569 = has_destructure_QMARK__548
	goto b16
}
func expand_map_pattern(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var defaults_3 vm.Value
	var as_sym_5 vm.Value
	var keys_STAR__7 vm.Value
	var strs_STAR__9 vm.Value
	var ks_13 vm.Value
	var out_14 vm.Value
	var get_expr_15 vm.Value
	var defaults_16 vm.Value
	var v67 vm.Value
	var sym_17 vm.Value
	var pat_18 vm.Value
	var defaults_19 vm.Value
	var as_sym_20 vm.Value
	var or__x_21 vm.Value
	var keys_STAR__22 vm.Value
	var strs_STAR__23 vm.Value
	var get_expr_24 vm.Value
	var sym_25 vm.Value
	var pat_26 vm.Value
	var defaults_27 vm.Value
	var as_sym_28 vm.Value
	var or__x_29 vm.Value
	var keys_STAR__30 vm.Value
	var strs_STAR__31 vm.Value
	var get_expr_32 vm.Value
	var v37 vm.Value
	var sym_38 vm.Value
	var pat_39 vm.Value
	var defaults_40 vm.Value
	var as_sym_41 vm.Value
	var or__x_42 vm.Value
	var keys_STAR__43 vm.Value
	var strs_STAR__44 vm.Value
	var get_expr_45 vm.Value
	var ks_48 vm.Value
	var out_49 vm.Value
	var get_expr_50 vm.Value
	var defaults_51 vm.Value
	var sym_52 vm.Value
	var pat_53 vm.Value
	var as_sym_54 vm.Value
	var keys_STAR__55 vm.Value
	var strs_STAR__56 vm.Value
	var ks_57 vm.Value
	var out_58 vm.Value
	var get_expr_59 vm.Value
	var defaults_60 vm.Value
	var sym_61 vm.Value
	var pat_62 vm.Value
	var as_sym_63 vm.Value
	var keys_STAR__64 vm.Value
	var strs_STAR__65 vm.Value
	var k_71 vm.Value
	var v93 vm.Value
	var binds_129 vm.Value
	var ks_130 vm.Value
	var out_131 vm.Value
	var get_expr_132 vm.Value
	var defaults_133 vm.Value
	var sym_134 vm.Value
	var pat_135 vm.Value
	var as_sym_136 vm.Value
	var keys_STAR__137 vm.Value
	var strs_STAR__138 vm.Value
	var ks_72 vm.Value
	var out_73 vm.Value
	var get_expr_74 vm.Value
	var defaults_75 vm.Value
	var sym_76 vm.Value
	var pat_77 vm.Value
	var as_sym_78 vm.Value
	var keys_STAR__79 vm.Value
	var strs_STAR__80 vm.Value
	var k_81 vm.Value
	var arg__6338_96 vm.Value
	var arg__6343_99 vm.Value
	var v100 vm.Value
	var ks_82 vm.Value
	var out_83 vm.Value
	var get_expr_84 vm.Value
	var defaults_85 vm.Value
	var sym_86 vm.Value
	var pat_87 vm.Value
	var as_sym_88 vm.Value
	var keys_STAR__89 vm.Value
	var strs_STAR__90 vm.Value
	var k_91 vm.Value
	var kn_103 vm.Value
	var ks_104 vm.Value
	var out_105 vm.Value
	var get_expr_106 vm.Value
	var defaults_107 vm.Value
	var sym_108 vm.Value
	var pat_109 vm.Value
	var as_sym_110 vm.Value
	var keys_STAR__111 vm.Value
	var strs_STAR__112 vm.Value
	var k_113 vm.Value
	var v115 vm.Value
	var arg__6355_117 vm.Value
	var arg__6363_119 vm.Value
	var arg__6364_120 vm.Value
	var arg__6374_123 vm.Value
	var arg__6382_125 vm.Value
	var arg__6383_126 vm.Value
	var v127 vm.Value
	var binds_139 vm.Value
	var ks_140 vm.Value
	var out_141 vm.Value
	var get_expr_142 vm.Value
	var defaults_143 vm.Value
	var sym_144 vm.Value
	var pat_145 vm.Value
	var as_sym_146 vm.Value
	var keys_STAR__147 vm.Value
	var strs_STAR__148 vm.Value
	var v169 vm.Value
	var binds_149 vm.Value
	var ks_150 vm.Value
	var out_151 vm.Value
	var get_expr_152 vm.Value
	var defaults_153 vm.Value
	var sym_154 vm.Value
	var pat_155 vm.Value
	var as_sym_156 vm.Value
	var keys_STAR__157 vm.Value
	var strs_STAR__158 vm.Value
	var binds_172 vm.Value
	var binds_173 vm.Value
	var ks_174 vm.Value
	var out_175 vm.Value
	var get_expr_176 vm.Value
	var defaults_177 vm.Value
	var sym_178 vm.Value
	var pat_179 vm.Value
	var as_sym_180 vm.Value
	var keys_STAR__181 vm.Value
	var strs_STAR__182 vm.Value
	var binds_183 vm.Value
	var ks_184 vm.Value
	var out_185 vm.Value
	var get_expr_186 vm.Value
	var defaults_187 vm.Value
	var sym_188 vm.Value
	var pat_189 vm.Value
	var as_sym_190 vm.Value
	var keys_STAR__191 vm.Value
	var strs_STAR__192 vm.Value
	var v205 vm.Value
	var binds_193 vm.Value
	var ks_194 vm.Value
	var out_195 vm.Value
	var get_expr_196 vm.Value
	var defaults_197 vm.Value
	var sym_198 vm.Value
	var pat_199 vm.Value
	var as_sym_200 vm.Value
	var keys_STAR__201 vm.Value
	var strs_STAR__202 vm.Value
	var binds_208 vm.Value
	var binds_209 vm.Value
	var ks_210 vm.Value
	var out_211 vm.Value
	var get_expr_212 vm.Value
	var defaults_213 vm.Value
	var sym_214 vm.Value
	var pat_215 vm.Value
	var as_sym_216 vm.Value
	var keys_STAR__217 vm.Value
	var strs_STAR__218 vm.Value
	var v220 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = defaults_3, as_sym_5, keys_STAR__7, strs_STAR__9, ks_13, out_14, get_expr_15, defaults_16, v67, sym_17, pat_18, defaults_19, as_sym_20, or__x_21, keys_STAR__22, strs_STAR__23, get_expr_24, sym_25, pat_26, defaults_27, as_sym_28, or__x_29, keys_STAR__30, strs_STAR__31, get_expr_32, v37, sym_38, pat_39, defaults_40, as_sym_41, or__x_42, keys_STAR__43, strs_STAR__44, get_expr_45, ks_48, out_49, get_expr_50, defaults_51, sym_52, pat_53, as_sym_54, keys_STAR__55, strs_STAR__56, ks_57, out_58, get_expr_59, defaults_60, sym_61, pat_62, as_sym_63, keys_STAR__64, strs_STAR__65, k_71, v93, binds_129, ks_130, out_131, get_expr_132, defaults_133, sym_134, pat_135, as_sym_136, keys_STAR__137, strs_STAR__138, ks_72, out_73, get_expr_74, defaults_75, sym_76, pat_77, as_sym_78, keys_STAR__79, strs_STAR__80, k_81, arg__6338_96, arg__6343_99, v100, ks_82, out_83, get_expr_84, defaults_85, sym_86, pat_87, as_sym_88, keys_STAR__89, strs_STAR__90, k_91, kn_103, ks_104, out_105, get_expr_106, defaults_107, sym_108, pat_109, as_sym_110, keys_STAR__111, strs_STAR__112, k_113, v115, arg__6355_117, arg__6363_119, arg__6364_120, arg__6374_123, arg__6382_125, arg__6383_126, v127, binds_139, ks_140, out_141, get_expr_142, defaults_143, sym_144, pat_145, as_sym_146, keys_STAR__147, strs_STAR__148, v169, binds_149, ks_150, out_151, get_expr_152, defaults_153, sym_154, pat_155, as_sym_156, keys_STAR__157, strs_STAR__158, binds_172, binds_173, ks_174, out_175, get_expr_176, defaults_177, sym_178, pat_179, as_sym_180, keys_STAR__181, strs_STAR__182, binds_183, ks_184, out_185, get_expr_186, defaults_187, sym_188, pat_189, as_sym_190, keys_STAR__191, strs_STAR__192, v205, binds_193, ks_194, out_195, get_expr_196, defaults_197, sym_198, pat_199, as_sym_200, keys_STAR__201, strs_STAR__202, binds_208, binds_209, ks_210, out_211, get_expr_212, defaults_213, sym_214, pat_215, as_sym_216, keys_STAR__217, strs_STAR__218, v220
	defaults_3, callErr = rt.InvokeValue(vm.Keyword("or"), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	as_sym_5, callErr = rt.InvokeValue(vm.Keyword("as"), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	keys_STAR__7, callErr = rt.InvokeValue(vm.Keyword("keys"), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	strs_STAR__9, callErr = rt.InvokeValue(vm.Keyword("strs"), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(keys_STAR__7) {
		sym_17 = arg0
		pat_18 = arg1
		defaults_19 = defaults_3
		as_sym_20 = as_sym_5
		or__x_21 = keys_STAR__7
		keys_STAR__22 = keys_STAR__7
		strs_STAR__23 = strs_STAR__9
		get_expr_24 = rt.BoxNativeFn(func(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
			var kn_3 vm.Value
			var default__4 vm.Value
			var sym_5 vm.Value
			var arg__6303_12 vm.Value
			var arg__6311_16 vm.Value
			var v17 vm.Value
			var kn_6 vm.Value
			var default__7 vm.Value
			var sym_8 vm.Value
			var arg__6318_21 vm.Value
			var arg__6325_25 vm.Value
			var v26 vm.Value
			var v28 vm.Value
			var kn_29 vm.Value
			var default__30 vm.Value
			var sym_31 vm.Value
			var callErr error
			_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = kn_3, default__4, sym_5, arg__6303_12, arg__6311_16, v17, kn_6, default__7, sym_8, arg__6318_21, arg__6325_25, v26, v28, kn_29, default__30, sym_31
			if vm.IsTruthy(arg1) {
				kn_3 = arg0
				default__4 = arg1
				sym_5 = arg0
				goto b1
			} else {
				kn_6 = arg0
				default__7 = arg1
				sym_8 = arg0
				goto b2
			}
		b1:
			;
			arg__6303_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "keyword").Deref(), []vm.Value{kn_3})
			if callErr != nil {
				return nil, callErr
			}
			arg__6311_16, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "keyword").Deref(), []vm.Value{kn_3})
			if callErr != nil {
				return nil, callErr
			}
			v17, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "list").Deref(), []vm.Value{vm.Symbol("get"), sym_5, arg__6311_16, default__4})
			if callErr != nil {
				return nil, callErr
			}
			v28 = v17
			kn_29 = kn_3
			default__30 = default__4
			sym_31 = sym_5
			goto b3
		b2:
			;
			arg__6318_21, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "keyword").Deref(), []vm.Value{kn_6})
			if callErr != nil {
				return nil, callErr
			}
			arg__6325_25, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "keyword").Deref(), []vm.Value{kn_6})
			if callErr != nil {
				return nil, callErr
			}
			v26, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "list").Deref(), []vm.Value{vm.Symbol("get"), sym_8, arg__6325_25})
			if callErr != nil {
				return nil, callErr
			}
			v28 = v26
			kn_29 = kn_6
			default__30 = default__7
			sym_31 = sym_8
			goto b3
		b3:
			;
			return v28, nil
		})
		goto b2
	} else {
		sym_25 = arg0
		pat_26 = arg1
		defaults_27 = defaults_3
		as_sym_28 = as_sym_5
		or__x_29 = keys_STAR__7
		keys_STAR__30 = keys_STAR__7
		strs_STAR__31 = strs_STAR__9
		get_expr_32 = rt.BoxNativeFn(func(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
			var kn_3 vm.Value
			var default__4 vm.Value
			var sym_5 vm.Value
			var arg__6303_12 vm.Value
			var arg__6311_16 vm.Value
			var v17 vm.Value
			var kn_6 vm.Value
			var default__7 vm.Value
			var sym_8 vm.Value
			var arg__6318_21 vm.Value
			var arg__6325_25 vm.Value
			var v26 vm.Value
			var v28 vm.Value
			var kn_29 vm.Value
			var default__30 vm.Value
			var sym_31 vm.Value
			var callErr error
			_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = kn_3, default__4, sym_5, arg__6303_12, arg__6311_16, v17, kn_6, default__7, sym_8, arg__6318_21, arg__6325_25, v26, v28, kn_29, default__30, sym_31
			if vm.IsTruthy(arg1) {
				kn_3 = arg0
				default__4 = arg1
				sym_5 = arg0
				goto b1
			} else {
				kn_6 = arg0
				default__7 = arg1
				sym_8 = arg0
				goto b2
			}
		b1:
			;
			arg__6303_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "keyword").Deref(), []vm.Value{kn_3})
			if callErr != nil {
				return nil, callErr
			}
			arg__6311_16, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "keyword").Deref(), []vm.Value{kn_3})
			if callErr != nil {
				return nil, callErr
			}
			v17, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "list").Deref(), []vm.Value{vm.Symbol("get"), sym_5, arg__6311_16, default__4})
			if callErr != nil {
				return nil, callErr
			}
			v28 = v17
			kn_29 = kn_3
			default__30 = default__4
			sym_31 = sym_5
			goto b3
		b2:
			;
			arg__6318_21, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "keyword").Deref(), []vm.Value{kn_6})
			if callErr != nil {
				return nil, callErr
			}
			arg__6325_25, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "keyword").Deref(), []vm.Value{kn_6})
			if callErr != nil {
				return nil, callErr
			}
			v26, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "list").Deref(), []vm.Value{vm.Symbol("get"), sym_8, arg__6325_25})
			if callErr != nil {
				return nil, callErr
			}
			v28 = v26
			kn_29 = kn_6
			default__30 = default__7
			sym_31 = sym_8
			goto b3
		b3:
			;
			return v28, nil
		})
		goto b3
	}
b1:
	;
	v67, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "empty?").Deref(), []vm.Value{ks_13})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v67) {
		ks_48 = ks_13
		out_49 = out_14
		get_expr_50 = get_expr_15
		defaults_51 = defaults_16
		sym_52 = sym_38
		pat_53 = pat_39
		as_sym_54 = as_sym_41
		keys_STAR__55 = keys_STAR__43
		strs_STAR__56 = strs_STAR__44
		goto b5
	} else {
		ks_57 = ks_13
		out_58 = out_14
		get_expr_59 = get_expr_15
		defaults_60 = defaults_16
		sym_61 = sym_38
		pat_62 = pat_39
		as_sym_63 = as_sym_41
		keys_STAR__64 = keys_STAR__43
		strs_STAR__65 = strs_STAR__44
		goto b6
	}
b2:
	;
	v37 = or__x_21
	sym_38 = sym_17
	pat_39 = pat_18
	defaults_40 = defaults_19
	as_sym_41 = as_sym_20
	or__x_42 = or__x_21
	keys_STAR__43 = keys_STAR__22
	strs_STAR__44 = strs_STAR__23
	get_expr_45 = get_expr_24
	goto b4
b3:
	;
	v37 = vm.NewArrayVector([]vm.Value{})
	sym_38 = sym_25
	pat_39 = pat_26
	defaults_40 = defaults_27
	as_sym_41 = as_sym_28
	or__x_42 = or__x_29
	keys_STAR__43 = keys_STAR__30
	strs_STAR__44 = strs_STAR__31
	get_expr_45 = get_expr_32
	goto b4
b4:
	;
	ks_13 = v37
	out_14 = vm.NewArrayVector([]vm.Value{})
	get_expr_15 = get_expr_45
	defaults_16 = defaults_40
	goto b1
b5:
	;
	binds_129 = out_49
	ks_130 = ks_48
	out_131 = out_49
	get_expr_132 = get_expr_50
	defaults_133 = defaults_51
	sym_134 = sym_52
	pat_135 = pat_53
	as_sym_136 = as_sym_54
	keys_STAR__137 = keys_STAR__55
	strs_STAR__138 = strs_STAR__56
	goto b7
b6:
	;
	k_71, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{ks_57})
	if callErr != nil {
		return nil, callErr
	}
	v93, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "keyword?").Deref(), []vm.Value{k_71})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v93) {
		ks_72 = ks_57
		out_73 = out_58
		get_expr_74 = get_expr_59
		defaults_75 = defaults_60
		sym_76 = sym_61
		pat_77 = pat_62
		as_sym_78 = as_sym_63
		keys_STAR__79 = keys_STAR__64
		strs_STAR__80 = strs_STAR__65
		k_81 = k_71
		goto b8
	} else {
		ks_82 = ks_57
		out_83 = out_58
		get_expr_84 = get_expr_59
		defaults_85 = defaults_60
		sym_86 = sym_61
		pat_87 = pat_62
		as_sym_88 = as_sym_63
		keys_STAR__89 = keys_STAR__64
		strs_STAR__90 = strs_STAR__65
		k_91 = k_71
		goto b9
	}
b7:
	;
	if vm.IsTruthy(strs_STAR__138) {
		binds_139 = binds_129
		ks_140 = ks_130
		out_141 = out_131
		get_expr_142 = get_expr_132
		defaults_143 = defaults_133
		sym_144 = sym_134
		pat_145 = pat_135
		as_sym_146 = as_sym_136
		keys_STAR__147 = keys_STAR__137
		strs_STAR__148 = strs_STAR__138
		goto b11
	} else {
		binds_149 = binds_129
		ks_150 = ks_130
		out_151 = out_131
		get_expr_152 = get_expr_132
		defaults_153 = defaults_133
		sym_154 = sym_134
		pat_155 = pat_135
		as_sym_156 = as_sym_136
		keys_STAR__157 = keys_STAR__137
		strs_STAR__158 = strs_STAR__138
		goto b12
	}
b8:
	;
	arg__6338_96, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "name").Deref(), []vm.Value{k_81})
	if callErr != nil {
		return nil, callErr
	}
	arg__6343_99, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "name").Deref(), []vm.Value{k_81})
	if callErr != nil {
		return nil, callErr
	}
	v100, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "symbol").Deref(), []vm.Value{arg__6343_99})
	if callErr != nil {
		return nil, callErr
	}
	kn_103 = v100
	ks_104 = ks_72
	out_105 = out_73
	get_expr_106 = get_expr_74
	defaults_107 = defaults_75
	sym_108 = sym_76
	pat_109 = pat_77
	as_sym_110 = as_sym_78
	keys_STAR__111 = keys_STAR__79
	strs_STAR__112 = strs_STAR__80
	k_113 = k_81
	goto b10
b9:
	;
	kn_103 = k_91
	ks_104 = ks_82
	out_105 = out_83
	get_expr_106 = get_expr_84
	defaults_107 = defaults_85
	sym_108 = sym_86
	pat_109 = pat_87
	as_sym_110 = as_sym_88
	keys_STAR__111 = keys_STAR__89
	strs_STAR__112 = strs_STAR__90
	k_113 = k_91
	goto b10
b10:
	;
	v115, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{ks_104})
	if callErr != nil {
		return nil, callErr
	}
	arg__6355_117, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{defaults_107, kn_103})
	if callErr != nil {
		return nil, callErr
	}
	arg__6363_119, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{defaults_107, kn_103})
	if callErr != nil {
		return nil, callErr
	}
	arg__6364_120, callErr = rt.InvokeValue(get_expr_106, []vm.Value{kn_103, arg__6363_119})
	if callErr != nil {
		return nil, callErr
	}
	arg__6374_123, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{defaults_107, kn_103})
	if callErr != nil {
		return nil, callErr
	}
	arg__6382_125, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{defaults_107, kn_103})
	if callErr != nil {
		return nil, callErr
	}
	arg__6383_126, callErr = rt.InvokeValue(get_expr_106, []vm.Value{kn_103, arg__6382_125})
	if callErr != nil {
		return nil, callErr
	}
	v127, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{out_105, kn_103, arg__6383_126})
	if callErr != nil {
		return nil, callErr
	}
	ks_13 = v115
	out_14 = v127
	get_expr_15 = get_expr_106
	defaults_16 = defaults_107
	goto b1
b11:
	;
	v169, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "reduce").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
		var v13 vm.Value
		var out_4 vm.Value
		var s_5 vm.Value
		var defaults_6 vm.Value
		var get_expr_7 vm.Value
		var arg__6443_16 vm.Value
		var arg__6448_19 vm.Value
		var v20 vm.Value
		var out_8 vm.Value
		var s_9 vm.Value
		var defaults_10 vm.Value
		var get_expr_11 vm.Value
		var kn_23 vm.Value
		var out_24 vm.Value
		var s_25 vm.Value
		var defaults_26 vm.Value
		var get_expr_27 vm.Value
		var arg__6457_29 vm.Value
		var arg__6465_31 vm.Value
		var arg__6466_32 vm.Value
		var arg__6476_35 vm.Value
		var arg__6484_37 vm.Value
		var arg__6485_38 vm.Value
		var v39 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = v13, out_4, s_5, defaults_6, get_expr_7, arg__6443_16, arg__6448_19, v20, out_8, s_9, defaults_10, get_expr_11, kn_23, out_24, s_25, defaults_26, get_expr_27, arg__6457_29, arg__6465_31, arg__6466_32, arg__6476_35, arg__6484_37, arg__6485_38, v39
		v13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "keyword?").Deref(), []vm.Value{arg1})
		if callErr != nil {
			return nil, callErr
		}
		if vm.IsTruthy(v13) {
			out_4 = arg0
			s_5 = arg1
			defaults_6 = defaults_143
			get_expr_7 = get_expr_142
			goto b1
		} else {
			out_8 = arg0
			s_9 = arg1
			defaults_10 = defaults_143
			get_expr_11 = get_expr_142
			goto b2
		}
	b1:
		;
		arg__6443_16, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "name").Deref(), []vm.Value{s_5})
		if callErr != nil {
			return nil, callErr
		}
		arg__6448_19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "name").Deref(), []vm.Value{s_5})
		if callErr != nil {
			return nil, callErr
		}
		v20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "symbol").Deref(), []vm.Value{arg__6448_19})
		if callErr != nil {
			return nil, callErr
		}
		kn_23 = v20
		out_24 = out_4
		s_25 = s_5
		defaults_26 = defaults_6
		get_expr_27 = get_expr_7
		goto b3
	b2:
		;
		kn_23 = s_9
		out_24 = out_8
		s_25 = s_9
		defaults_26 = defaults_10
		get_expr_27 = get_expr_11
		goto b3
	b3:
		;
		arg__6457_29, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{defaults_26, kn_23})
		if callErr != nil {
			return nil, callErr
		}
		arg__6465_31, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{defaults_26, kn_23})
		if callErr != nil {
			return nil, callErr
		}
		arg__6466_32, callErr = rt.InvokeValue(get_expr_27, []vm.Value{kn_23, arg__6465_31})
		if callErr != nil {
			return nil, callErr
		}
		arg__6476_35, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{defaults_26, kn_23})
		if callErr != nil {
			return nil, callErr
		}
		arg__6484_37, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{defaults_26, kn_23})
		if callErr != nil {
			return nil, callErr
		}
		arg__6485_38, callErr = rt.InvokeValue(get_expr_27, []vm.Value{kn_23, arg__6484_37})
		if callErr != nil {
			return nil, callErr
		}
		v39, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{out_24, kn_23, arg__6485_38})
		if callErr != nil {
			return nil, callErr
		}
		return v39, nil
	}), binds_139, strs_STAR__148})
	if callErr != nil {
		return nil, callErr
	}
	binds_172 = v169
	binds_173 = binds_139
	ks_174 = ks_140
	out_175 = out_141
	get_expr_176 = get_expr_142
	defaults_177 = defaults_143
	sym_178 = sym_144
	pat_179 = pat_145
	as_sym_180 = as_sym_146
	keys_STAR__181 = keys_STAR__147
	strs_STAR__182 = strs_STAR__148
	goto b13
b12:
	;
	binds_172 = binds_149
	binds_173 = binds_149
	ks_174 = ks_150
	out_175 = out_151
	get_expr_176 = get_expr_152
	defaults_177 = defaults_153
	sym_178 = sym_154
	pat_179 = pat_155
	as_sym_180 = as_sym_156
	keys_STAR__181 = keys_STAR__157
	strs_STAR__182 = strs_STAR__158
	goto b13
b13:
	;
	if vm.IsTruthy(as_sym_180) {
		binds_183 = binds_172
		ks_184 = ks_174
		out_185 = out_175
		get_expr_186 = get_expr_176
		defaults_187 = defaults_177
		sym_188 = sym_178
		pat_189 = pat_179
		as_sym_190 = as_sym_180
		keys_STAR__191 = keys_STAR__181
		strs_STAR__192 = strs_STAR__182
		goto b14
	} else {
		binds_193 = binds_172
		ks_194 = ks_174
		out_195 = out_175
		get_expr_196 = get_expr_176
		defaults_197 = defaults_177
		sym_198 = sym_178
		pat_199 = pat_179
		as_sym_200 = as_sym_180
		keys_STAR__201 = keys_STAR__181
		strs_STAR__202 = strs_STAR__182
		goto b15
	}
b14:
	;
	v205, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{binds_183, as_sym_190, sym_188})
	if callErr != nil {
		return nil, callErr
	}
	binds_208 = v205
	binds_209 = binds_183
	ks_210 = ks_184
	out_211 = out_185
	get_expr_212 = get_expr_186
	defaults_213 = defaults_187
	sym_214 = sym_188
	pat_215 = pat_189
	as_sym_216 = as_sym_190
	keys_STAR__217 = keys_STAR__191
	strs_STAR__218 = strs_STAR__192
	goto b16
b15:
	;
	binds_208 = binds_193
	binds_209 = binds_193
	ks_210 = ks_194
	out_211 = out_195
	get_expr_212 = get_expr_196
	defaults_213 = defaults_197
	sym_214 = sym_198
	pat_215 = pat_199
	as_sym_216 = as_sym_200
	keys_STAR__217 = keys_STAR__201
	strs_STAR__218 = strs_STAR__202
	goto b16
b16:
	;
	v220, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "expand-binding").Deref(), []vm.Value{binds_208})
	if callErr != nil {
		return nil, callErr
	}
	return v220, nil
}
func expand_vector_pattern(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var i_2 int
	var remaining_3 vm.Value
	var out_4 vm.Value
	var sym_5 vm.Value
	var v20 vm.Value
	var pat_9 vm.Value
	var i_10 int
	var remaining_11 vm.Value
	var out_12 vm.Value
	var sym_13 vm.Value
	var v23 vm.Value
	var pat_14 vm.Value
	var i_15 int
	var remaining_16 vm.Value
	var out_17 vm.Value
	var sym_18 vm.Value
	var x_26 vm.Value
	var v40 bool
	var v444 vm.Value
	var pat_445 vm.Value
	var i_446 int
	var remaining_447 vm.Value
	var out_448 vm.Value
	var sym_449 vm.Value
	var pat_27 vm.Value
	var i_28 int
	var remaining_29 vm.Value
	var out_30 vm.Value
	var sym_31 vm.Value
	var x_32 vm.Value
	var rest_pat_43 vm.Value
	var rem_STAR__47 vm.Value
	var rest_expr_51 vm.Value
	var v71 vm.Value
	var pat_33 vm.Value
	var i_34 int
	var remaining_35 vm.Value
	var out_36 vm.Value
	var sym_37 vm.Value
	var x_38 vm.Value
	var v198 bool
	var v436 vm.Value
	var pat_437 vm.Value
	var i_438 int
	var remaining_439 vm.Value
	var out_440 vm.Value
	var sym_441 vm.Value
	var x_442 vm.Value
	var pat_52 vm.Value
	var i_53 int
	var remaining_54 vm.Value
	var out_55 vm.Value
	var sym_56 vm.Value
	var x_57 vm.Value
	var rest_pat_58 vm.Value
	var rem_STAR__59 vm.Value
	var rest_expr_60 vm.Value
	var v74 vm.Value
	var pat_61 vm.Value
	var i_62 int
	var remaining_63 vm.Value
	var out_64 vm.Value
	var sym_65 vm.Value
	var x_66 vm.Value
	var rest_pat_67 vm.Value
	var rem_STAR__68 vm.Value
	var rest_expr_69 vm.Value
	var gs_79 vm.Value
	var arg__6545_81 vm.Value
	var arg__6549_83 vm.Value
	var arg__6554_86 vm.Value
	var arg__6555_87 vm.Value
	var arg__6564_90 vm.Value
	var arg__6568_92 vm.Value
	var arg__6573_95 vm.Value
	var arg__6574_96 vm.Value
	var v97 vm.Value
	var out_STAR__99 vm.Value
	var pat_100 vm.Value
	var i_101 int
	var remaining_102 vm.Value
	var out_103 vm.Value
	var sym_104 vm.Value
	var x_105 vm.Value
	var rest_pat_106 vm.Value
	var rem_STAR__107 vm.Value
	var rest_expr_108 vm.Value
	var and__x_130 vm.Value
	var out_STAR__109 vm.Value
	var pat_110 vm.Value
	var i_111 int
	var remaining_112 vm.Value
	var out_113 vm.Value
	var sym_114 vm.Value
	var x_115 vm.Value
	var rest_pat_116 vm.Value
	var rem_STAR__117 vm.Value
	var rest_expr_118 vm.Value
	var v176 vm.Value
	var arg__6592_178 vm.Value
	var arg__6599_181 vm.Value
	var v182 vm.Value
	var out_STAR__119 vm.Value
	var pat_120 vm.Value
	var i_121 int
	var remaining_122 vm.Value
	var out_123 vm.Value
	var sym_124 vm.Value
	var x_125 vm.Value
	var rest_pat_126 vm.Value
	var rem_STAR__127 vm.Value
	var rest_expr_128 vm.Value
	var out_STAR__131 vm.Value
	var pat_132 vm.Value
	var i_133 int
	var remaining_134 vm.Value
	var out_135 vm.Value
	var sym_136 vm.Value
	var x_137 vm.Value
	var rest_pat_138 vm.Value
	var rem_STAR__139 vm.Value
	var rest_expr_140 vm.Value
	var and__x_141 vm.Value
	var arg__6581_155 vm.Value
	var v157 bool
	var out_STAR__142 vm.Value
	var pat_143 vm.Value
	var i_144 int
	var remaining_145 vm.Value
	var out_146 vm.Value
	var sym_147 vm.Value
	var x_148 vm.Value
	var rest_pat_149 vm.Value
	var rem_STAR__150 vm.Value
	var rest_expr_151 vm.Value
	var and__x_152 vm.Value
	var v160 vm.Value
	var out_STAR__161 vm.Value
	var pat_162 vm.Value
	var i_163 int
	var remaining_164 vm.Value
	var out_165 vm.Value
	var sym_166 vm.Value
	var x_167 vm.Value
	var rest_pat_168 vm.Value
	var rem_STAR__169 vm.Value
	var rest_expr_170 vm.Value
	var and__x_171 vm.Value
	var pat_185 vm.Value
	var i_186 int
	var remaining_187 vm.Value
	var out_188 vm.Value
	var sym_189 vm.Value
	var x_190 vm.Value
	var v203 vm.Value
	var arg__6612_205 vm.Value
	var arg__6619_208 vm.Value
	var v209 vm.Value
	var pat_191 vm.Value
	var i_192 int
	var remaining_193 vm.Value
	var out_194 vm.Value
	var sym_195 vm.Value
	var x_196 vm.Value
	var v224 vm.Value
	var v428 vm.Value
	var pat_429 vm.Value
	var i_430 int
	var remaining_431 vm.Value
	var out_432 vm.Value
	var sym_433 vm.Value
	var x_434 vm.Value
	var pat_211 vm.Value
	var i_212 int
	var remaining_213 vm.Value
	var out_214 vm.Value
	var sym_215 vm.Value
	var x_216 vm.Value
	var v226 int
	var v228 vm.Value
	var arg__6639_234 vm.Value
	var arg__6652_241 vm.Value
	var v242 vm.Value
	var pat_217 vm.Value
	var i_218 int
	var remaining_219 vm.Value
	var out_220 vm.Value
	var sym_221 vm.Value
	var x_222 vm.Value
	var or__x_257 vm.Value
	var v420 vm.Value
	var pat_421 vm.Value
	var i_422 int
	var remaining_423 vm.Value
	var out_424 vm.Value
	var sym_425 vm.Value
	var x_426 vm.Value
	var pat_244 vm.Value
	var i_245 int
	var remaining_246 vm.Value
	var out_247 vm.Value
	var sym_248 vm.Value
	var x_249 vm.Value
	var gs_289 vm.Value
	var v305 vm.Value
	var pat_250 vm.Value
	var i_251 int
	var remaining_252 vm.Value
	var out_253 vm.Value
	var sym_254 vm.Value
	var x_255 vm.Value
	var v412 vm.Value
	var pat_413 vm.Value
	var i_414 int
	var remaining_415 vm.Value
	var out_416 vm.Value
	var sym_417 vm.Value
	var x_418 vm.Value
	var pat_258 vm.Value
	var i_259 int
	var remaining_260 vm.Value
	var out_261 vm.Value
	var sym_262 vm.Value
	var x_263 vm.Value
	var or__x_264 vm.Value
	var pat_265 vm.Value
	var i_266 int
	var remaining_267 vm.Value
	var out_268 vm.Value
	var sym_269 vm.Value
	var x_270 vm.Value
	var or__x_271 vm.Value
	var v275 vm.Value
	var v277 vm.Value
	var pat_278 vm.Value
	var i_279 int
	var remaining_280 vm.Value
	var out_281 vm.Value
	var sym_282 vm.Value
	var x_283 vm.Value
	var or__x_284 vm.Value
	var pat_290 vm.Value
	var i_291 int
	var remaining_292 vm.Value
	var out_293 vm.Value
	var sym_294 vm.Value
	var x_295 vm.Value
	var gs_296 vm.Value
	var v308 vm.Value
	var pat_297 vm.Value
	var i_298 int
	var remaining_299 vm.Value
	var out_300 vm.Value
	var sym_301 vm.Value
	var x_302 vm.Value
	var gs_303 vm.Value
	var v325 vm.Value
	var nested_341 vm.Value
	var pat_342 vm.Value
	var i_343 int
	var remaining_344 vm.Value
	var out_345 vm.Value
	var sym_346 vm.Value
	var x_347 vm.Value
	var gs_348 vm.Value
	var arg__6689_354 vm.Value
	var arg__6702_361 vm.Value
	var out_STAR__362 vm.Value
	var v363 int
	var v365 vm.Value
	var v367 vm.Value
	var pat_310 vm.Value
	var i_311 int
	var remaining_312 vm.Value
	var out_313 vm.Value
	var sym_314 vm.Value
	var x_315 vm.Value
	var gs_316 vm.Value
	var v328 vm.Value
	var pat_317 vm.Value
	var i_318 int
	var remaining_319 vm.Value
	var out_320 vm.Value
	var sym_321 vm.Value
	var x_322 vm.Value
	var gs_323 vm.Value
	var v332 vm.Value
	var pat_333 vm.Value
	var i_334 int
	var remaining_335 vm.Value
	var out_336 vm.Value
	var sym_337 vm.Value
	var x_338 vm.Value
	var gs_339 vm.Value
	var pat_369 vm.Value
	var i_370 int
	var remaining_371 vm.Value
	var out_372 vm.Value
	var sym_373 vm.Value
	var x_374 vm.Value
	var arg__6716_385 vm.Value
	var arg__6722_389 vm.Value
	var arg__6723_390 vm.Value
	var arg__6729_394 vm.Value
	var arg__6735_398 vm.Value
	var arg__6736_399 vm.Value
	var v400 vm.Value
	var pat_375 vm.Value
	var i_376 int
	var remaining_377 vm.Value
	var out_378 vm.Value
	var sym_379 vm.Value
	var x_380 vm.Value
	var v404 vm.Value
	var pat_405 vm.Value
	var i_406 int
	var remaining_407 vm.Value
	var out_408 vm.Value
	var sym_409 vm.Value
	var x_410 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = i_2, remaining_3, out_4, sym_5, v20, pat_9, i_10, remaining_11, out_12, sym_13, v23, pat_14, i_15, remaining_16, out_17, sym_18, x_26, v40, v444, pat_445, i_446, remaining_447, out_448, sym_449, pat_27, i_28, remaining_29, out_30, sym_31, x_32, rest_pat_43, rem_STAR__47, rest_expr_51, v71, pat_33, i_34, remaining_35, out_36, sym_37, x_38, v198, v436, pat_437, i_438, remaining_439, out_440, sym_441, x_442, pat_52, i_53, remaining_54, out_55, sym_56, x_57, rest_pat_58, rem_STAR__59, rest_expr_60, v74, pat_61, i_62, remaining_63, out_64, sym_65, x_66, rest_pat_67, rem_STAR__68, rest_expr_69, gs_79, arg__6545_81, arg__6549_83, arg__6554_86, arg__6555_87, arg__6564_90, arg__6568_92, arg__6573_95, arg__6574_96, v97, out_STAR__99, pat_100, i_101, remaining_102, out_103, sym_104, x_105, rest_pat_106, rem_STAR__107, rest_expr_108, and__x_130, out_STAR__109, pat_110, i_111, remaining_112, out_113, sym_114, x_115, rest_pat_116, rem_STAR__117, rest_expr_118, v176, arg__6592_178, arg__6599_181, v182, out_STAR__119, pat_120, i_121, remaining_122, out_123, sym_124, x_125, rest_pat_126, rem_STAR__127, rest_expr_128, out_STAR__131, pat_132, i_133, remaining_134, out_135, sym_136, x_137, rest_pat_138, rem_STAR__139, rest_expr_140, and__x_141, arg__6581_155, v157, out_STAR__142, pat_143, i_144, remaining_145, out_146, sym_147, x_148, rest_pat_149, rem_STAR__150, rest_expr_151, and__x_152, v160, out_STAR__161, pat_162, i_163, remaining_164, out_165, sym_166, x_167, rest_pat_168, rem_STAR__169, rest_expr_170, and__x_171, pat_185, i_186, remaining_187, out_188, sym_189, x_190, v203, arg__6612_205, arg__6619_208, v209, pat_191, i_192, remaining_193, out_194, sym_195, x_196, v224, v428, pat_429, i_430, remaining_431, out_432, sym_433, x_434, pat_211, i_212, remaining_213, out_214, sym_215, x_216, v226, v228, arg__6639_234, arg__6652_241, v242, pat_217, i_218, remaining_219, out_220, sym_221, x_222, or__x_257, v420, pat_421, i_422, remaining_423, out_424, sym_425, x_426, pat_244, i_245, remaining_246, out_247, sym_248, x_249, gs_289, v305, pat_250, i_251, remaining_252, out_253, sym_254, x_255, v412, pat_413, i_414, remaining_415, out_416, sym_417, x_418, pat_258, i_259, remaining_260, out_261, sym_262, x_263, or__x_264, pat_265, i_266, remaining_267, out_268, sym_269, x_270, or__x_271, v275, v277, pat_278, i_279, remaining_280, out_281, sym_282, x_283, or__x_284, pat_290, i_291, remaining_292, out_293, sym_294, x_295, gs_296, v308, pat_297, i_298, remaining_299, out_300, sym_301, x_302, gs_303, v325, nested_341, pat_342, i_343, remaining_344, out_345, sym_346, x_347, gs_348, arg__6689_354, arg__6702_361, out_STAR__362, v363, v365, v367, pat_310, i_311, remaining_312, out_313, sym_314, x_315, gs_316, v328, pat_317, i_318, remaining_319, out_320, sym_321, x_322, gs_323, v332, pat_333, i_334, remaining_335, out_336, sym_337, x_338, gs_339, pat_369, i_370, remaining_371, out_372, sym_373, x_374, arg__6716_385, arg__6722_389, arg__6723_390, arg__6729_394, arg__6735_398, arg__6736_399, v400, pat_375, i_376, remaining_377, out_378, sym_379, x_380, v404, pat_405, i_406, remaining_407, out_408, sym_409, x_410
	i_2 = 0
	remaining_3 = arg1
	out_4 = vm.NewArrayVector([]vm.Value{})
	sym_5 = arg0
	goto b1
b1:
	;
	v20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "empty?").Deref(), []vm.Value{remaining_3})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v20) {
		pat_9 = arg1
		i_10 = i_2
		remaining_11 = remaining_3
		out_12 = out_4
		sym_13 = sym_5
		goto b2
	} else {
		pat_14 = arg1
		i_15 = i_2
		remaining_16 = remaining_3
		out_17 = out_4
		sym_18 = sym_5
		goto b3
	}
b2:
	;
	v23, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "expand-binding").Deref(), []vm.Value{out_12})
	if callErr != nil {
		return nil, callErr
	}
	v444 = v23
	pat_445 = pat_9
	i_446 = i_10
	remaining_447 = remaining_11
	out_448 = out_12
	sym_449 = sym_13
	goto b4
b3:
	;
	x_26, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{remaining_16})
	if callErr != nil {
		return nil, callErr
	}
	v40 = x_26 == vm.Symbol("&")
	if v40 {
		pat_27 = pat_14
		i_28 = i_15
		remaining_29 = remaining_16
		out_30 = out_17
		sym_31 = sym_18
		x_32 = x_26
		goto b5
	} else {
		pat_33 = pat_14
		i_34 = i_15
		remaining_35 = remaining_16
		out_36 = out_17
		sym_37 = sym_18
		x_38 = x_26
		goto b6
	}
b4:
	;
	return v444, nil
b5:
	;
	rest_pat_43, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "second").Deref(), []vm.Value{remaining_29})
	if callErr != nil {
		return nil, callErr
	}
	rem_STAR__47, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "drop").Deref(), []vm.Value{vm.Int(2), remaining_29})
	if callErr != nil {
		return nil, callErr
	}
	rest_expr_51, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "list").Deref(), []vm.Value{vm.Symbol("drop"), vm.Int(i_28), sym_31})
	if callErr != nil {
		return nil, callErr
	}
	v71, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "symbol?").Deref(), []vm.Value{rest_pat_43})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v71) {
		pat_52 = pat_27
		i_53 = i_28
		remaining_54 = remaining_29
		out_55 = out_30
		sym_56 = sym_31
		x_57 = x_32
		rest_pat_58 = rest_pat_43
		rem_STAR__59 = rem_STAR__47
		rest_expr_60 = rest_expr_51
		goto b8
	} else {
		pat_61 = pat_27
		i_62 = i_28
		remaining_63 = remaining_29
		out_64 = out_30
		sym_65 = sym_31
		x_66 = x_32
		rest_pat_67 = rest_pat_43
		rem_STAR__68 = rem_STAR__47
		rest_expr_69 = rest_expr_51
		goto b9
	}
b6:
	;
	v198 = x_38 == vm.Keyword("as")
	if v198 {
		pat_185 = pat_33
		i_186 = i_34
		remaining_187 = remaining_35
		out_188 = out_36
		sym_189 = sym_37
		x_190 = x_38
		goto b17
	} else {
		pat_191 = pat_33
		i_192 = i_34
		remaining_193 = remaining_35
		out_194 = out_36
		sym_195 = sym_37
		x_196 = x_38
		goto b18
	}
b7:
	;
	v444 = v436
	pat_445 = pat_437
	i_446 = i_438
	remaining_447 = remaining_439
	out_448 = out_440
	sym_449 = sym_441
	goto b4
b8:
	;
	v74, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{out_55, rest_pat_58, rest_expr_60})
	if callErr != nil {
		return nil, callErr
	}
	out_STAR__99 = v74
	pat_100 = pat_52
	i_101 = i_53
	remaining_102 = remaining_54
	out_103 = out_55
	sym_104 = sym_56
	x_105 = x_57
	rest_pat_106 = rest_pat_58
	rem_STAR__107 = rem_STAR__59
	rest_expr_108 = rest_expr_60
	goto b10
b9:
	;
	gs_79, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "gensym").Deref(), []vm.Value{vm.String("rest__")})
	if callErr != nil {
		return nil, callErr
	}
	arg__6545_81, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{out_64, gs_79, rest_expr_69})
	if callErr != nil {
		return nil, callErr
	}
	arg__6549_83, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{gs_79, rest_pat_67})
	if callErr != nil {
		return nil, callErr
	}
	arg__6554_86, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{gs_79, rest_pat_67})
	if callErr != nil {
		return nil, callErr
	}
	arg__6555_87, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "expand-binding").Deref(), []vm.Value{arg__6554_86})
	if callErr != nil {
		return nil, callErr
	}
	arg__6564_90, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{out_64, gs_79, rest_expr_69})
	if callErr != nil {
		return nil, callErr
	}
	arg__6568_92, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{gs_79, rest_pat_67})
	if callErr != nil {
		return nil, callErr
	}
	arg__6573_95, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{gs_79, rest_pat_67})
	if callErr != nil {
		return nil, callErr
	}
	arg__6574_96, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "expand-binding").Deref(), []vm.Value{arg__6573_95})
	if callErr != nil {
		return nil, callErr
	}
	v97, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__6564_90, arg__6574_96})
	if callErr != nil {
		return nil, callErr
	}
	out_STAR__99 = v97
	pat_100 = pat_61
	i_101 = i_62
	remaining_102 = remaining_63
	out_103 = out_64
	sym_104 = sym_65
	x_105 = x_66
	rest_pat_106 = rest_pat_67
	rem_STAR__107 = rem_STAR__68
	rest_expr_108 = rest_expr_69
	goto b10
b10:
	;
	and__x_130, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{rem_STAR__107})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(and__x_130) {
		out_STAR__131 = out_STAR__99
		pat_132 = pat_100
		i_133 = i_101
		remaining_134 = remaining_102
		out_135 = out_103
		sym_136 = sym_104
		x_137 = x_105
		rest_pat_138 = rest_pat_106
		rem_STAR__139 = rem_STAR__107
		rest_expr_140 = rest_expr_108
		and__x_141 = and__x_130
		goto b14
	} else {
		out_STAR__142 = out_STAR__99
		pat_143 = pat_100
		i_144 = i_101
		remaining_145 = remaining_102
		out_146 = out_103
		sym_147 = sym_104
		x_148 = x_105
		rest_pat_149 = rest_pat_106
		rem_STAR__150 = rem_STAR__107
		rest_expr_151 = rest_expr_108
		and__x_152 = and__x_130
		goto b15
	}
b11:
	;
	v176, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "drop").Deref(), []vm.Value{vm.Int(2), rem_STAR__117})
	if callErr != nil {
		return nil, callErr
	}
	arg__6592_178, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "second").Deref(), []vm.Value{rem_STAR__117})
	if callErr != nil {
		return nil, callErr
	}
	arg__6599_181, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "second").Deref(), []vm.Value{rem_STAR__117})
	if callErr != nil {
		return nil, callErr
	}
	v182, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{out_STAR__109, arg__6599_181, sym_114})
	if callErr != nil {
		return nil, callErr
	}
	i_2 = i_111
	remaining_3 = v176
	out_4 = v182
	sym_5 = sym_114
	goto b1
b12:
	;
	i_2 = i_121
	remaining_3 = rem_STAR__127
	out_4 = out_STAR__119
	sym_5 = sym_124
	goto b1
b14:
	;
	arg__6581_155, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{rem_STAR__139})
	if callErr != nil {
		return nil, callErr
	}
	v157 = arg__6581_155 == vm.Keyword("as")
	v160 = vm.Boolean(v157)
	out_STAR__161 = out_STAR__131
	pat_162 = pat_132
	i_163 = i_133
	remaining_164 = remaining_134
	out_165 = out_135
	sym_166 = sym_136
	x_167 = x_137
	rest_pat_168 = rest_pat_138
	rem_STAR__169 = rem_STAR__139
	rest_expr_170 = rest_expr_140
	and__x_171 = and__x_141
	goto b16
b15:
	;
	v160 = and__x_152
	out_STAR__161 = out_STAR__142
	pat_162 = pat_143
	i_163 = i_144
	remaining_164 = remaining_145
	out_165 = out_146
	sym_166 = sym_147
	x_167 = x_148
	rest_pat_168 = rest_pat_149
	rem_STAR__169 = rem_STAR__150
	rest_expr_170 = rest_expr_151
	and__x_171 = and__x_152
	goto b16
b16:
	;
	if vm.IsTruthy(v160) {
		out_STAR__109 = out_STAR__161
		pat_110 = pat_162
		i_111 = i_163
		remaining_112 = remaining_164
		out_113 = out_165
		sym_114 = sym_166
		x_115 = x_167
		rest_pat_116 = rest_pat_168
		rem_STAR__117 = rem_STAR__169
		rest_expr_118 = rest_expr_170
		goto b11
	} else {
		out_STAR__119 = out_STAR__161
		pat_120 = pat_162
		i_121 = i_163
		remaining_122 = remaining_164
		out_123 = out_165
		sym_124 = sym_166
		x_125 = x_167
		rest_pat_126 = rest_pat_168
		rem_STAR__127 = rem_STAR__169
		rest_expr_128 = rest_expr_170
		goto b12
	}
b17:
	;
	v203, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "drop").Deref(), []vm.Value{vm.Int(2), remaining_187})
	if callErr != nil {
		return nil, callErr
	}
	arg__6612_205, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "second").Deref(), []vm.Value{remaining_187})
	if callErr != nil {
		return nil, callErr
	}
	arg__6619_208, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "second").Deref(), []vm.Value{remaining_187})
	if callErr != nil {
		return nil, callErr
	}
	v209, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{out_188, arg__6619_208, sym_189})
	if callErr != nil {
		return nil, callErr
	}
	i_2 = i_186
	remaining_3 = v203
	out_4 = v209
	sym_5 = sym_189
	goto b1
b18:
	;
	v224, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "symbol?").Deref(), []vm.Value{x_196})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v224) {
		pat_211 = pat_191
		i_212 = i_192
		remaining_213 = remaining_193
		out_214 = out_194
		sym_215 = sym_195
		x_216 = x_196
		goto b20
	} else {
		pat_217 = pat_191
		i_218 = i_192
		remaining_219 = remaining_193
		out_220 = out_194
		sym_221 = sym_195
		x_222 = x_196
		goto b21
	}
b19:
	;
	v436 = v428
	pat_437 = pat_429
	i_438 = i_430
	remaining_439 = remaining_431
	out_440 = out_432
	sym_441 = sym_433
	x_442 = x_434
	goto b7
b20:
	;
	v226 = i_212 + 1
	v228, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{remaining_213})
	if callErr != nil {
		return nil, callErr
	}
	arg__6639_234, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "list").Deref(), []vm.Value{vm.Symbol("nth"), sym_215, vm.Int(i_212), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	arg__6652_241, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "list").Deref(), []vm.Value{vm.Symbol("nth"), sym_215, vm.Int(i_212), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	v242, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{out_214, x_216, arg__6652_241})
	if callErr != nil {
		return nil, callErr
	}
	i_2 = v226
	remaining_3 = v228
	out_4 = v242
	sym_5 = sym_215
	goto b1
b21:
	;
	or__x_257, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector?").Deref(), []vm.Value{x_222})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(or__x_257) {
		pat_258 = pat_217
		i_259 = i_218
		remaining_260 = remaining_219
		out_261 = out_220
		sym_262 = sym_221
		x_263 = x_222
		or__x_264 = or__x_257
		goto b26
	} else {
		pat_265 = pat_217
		i_266 = i_218
		remaining_267 = remaining_219
		out_268 = out_220
		sym_269 = sym_221
		x_270 = x_222
		or__x_271 = or__x_257
		goto b27
	}
b22:
	;
	v428 = v420
	pat_429 = pat_421
	i_430 = i_422
	remaining_431 = remaining_423
	out_432 = out_424
	sym_433 = sym_425
	x_434 = x_426
	goto b19
b23:
	;
	gs_289, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "gensym").Deref(), []vm.Value{vm.String("v__")})
	if callErr != nil {
		return nil, callErr
	}
	v305, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector?").Deref(), []vm.Value{x_249})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v305) {
		pat_290 = pat_244
		i_291 = i_245
		remaining_292 = remaining_246
		out_293 = out_247
		sym_294 = sym_248
		x_295 = x_249
		gs_296 = gs_289
		goto b29
	} else {
		pat_297 = pat_244
		i_298 = i_245
		remaining_299 = remaining_246
		out_300 = out_247
		sym_301 = sym_248
		x_302 = x_249
		gs_303 = gs_289
		goto b30
	}
b24:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		pat_369 = pat_250
		i_370 = i_251
		remaining_371 = remaining_252
		out_372 = out_253
		sym_373 = sym_254
		x_374 = x_255
		goto b35
	} else {
		pat_375 = pat_250
		i_376 = i_251
		remaining_377 = remaining_252
		out_378 = out_253
		sym_379 = sym_254
		x_380 = x_255
		goto b36
	}
b25:
	;
	v420 = v412
	pat_421 = pat_413
	i_422 = i_414
	remaining_423 = remaining_415
	out_424 = out_416
	sym_425 = sym_417
	x_426 = x_418
	goto b22
b26:
	;
	v277 = or__x_264
	pat_278 = pat_258
	i_279 = i_259
	remaining_280 = remaining_260
	out_281 = out_261
	sym_282 = sym_262
	x_283 = x_263
	or__x_284 = or__x_264
	goto b28
b27:
	;
	v275, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map?").Deref(), []vm.Value{x_270})
	if callErr != nil {
		return nil, callErr
	}
	v277 = v275
	pat_278 = pat_265
	i_279 = i_266
	remaining_280 = remaining_267
	out_281 = out_268
	sym_282 = sym_269
	x_283 = x_270
	or__x_284 = or__x_271
	goto b28
b28:
	;
	if vm.IsTruthy(v277) {
		pat_244 = pat_278
		i_245 = i_279
		remaining_246 = remaining_280
		out_247 = out_281
		sym_248 = sym_282
		x_249 = x_283
		goto b23
	} else {
		pat_250 = pat_278
		i_251 = i_279
		remaining_252 = remaining_280
		out_253 = out_281
		sym_254 = sym_282
		x_255 = x_283
		goto b24
	}
b29:
	;
	v308, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "expand-vector-pattern").Deref(), []vm.Value{gs_296, x_295})
	if callErr != nil {
		return nil, callErr
	}
	nested_341 = v308
	pat_342 = pat_290
	i_343 = i_291
	remaining_344 = remaining_292
	out_345 = out_293
	sym_346 = sym_294
	x_347 = x_295
	gs_348 = gs_296
	goto b31
b30:
	;
	v325, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map?").Deref(), []vm.Value{x_302})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v325) {
		pat_310 = pat_297
		i_311 = i_298
		remaining_312 = remaining_299
		out_313 = out_300
		sym_314 = sym_301
		x_315 = x_302
		gs_316 = gs_303
		goto b32
	} else {
		pat_317 = pat_297
		i_318 = i_298
		remaining_319 = remaining_299
		out_320 = out_300
		sym_321 = sym_301
		x_322 = x_302
		gs_323 = gs_303
		goto b33
	}
b31:
	;
	arg__6689_354, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "list").Deref(), []vm.Value{vm.Symbol("nth"), sym_346, vm.Int(i_343), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	arg__6702_361, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "list").Deref(), []vm.Value{vm.Symbol("nth"), sym_346, vm.Int(i_343), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	out_STAR__362, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{out_345, gs_348, arg__6702_361})
	if callErr != nil {
		return nil, callErr
	}
	v363 = i_343 + 1
	v365, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{remaining_344})
	if callErr != nil {
		return nil, callErr
	}
	v367, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{out_STAR__362, nested_341})
	if callErr != nil {
		return nil, callErr
	}
	i_2 = v363
	remaining_3 = v365
	out_4 = v367
	sym_5 = sym_346
	goto b1
b32:
	;
	v328, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "expand-map-pattern").Deref(), []vm.Value{gs_316, x_315})
	if callErr != nil {
		return nil, callErr
	}
	v332 = v328
	pat_333 = pat_310
	i_334 = i_311
	remaining_335 = remaining_312
	out_336 = out_313
	sym_337 = sym_314
	x_338 = x_315
	gs_339 = gs_316
	goto b34
b33:
	;
	v332 = vm.NIL
	pat_333 = pat_317
	i_334 = i_318
	remaining_335 = remaining_319
	out_336 = out_320
	sym_337 = sym_321
	x_338 = x_322
	gs_339 = gs_323
	goto b34
b34:
	;
	nested_341 = v332
	pat_342 = pat_333
	i_343 = i_334
	remaining_344 = remaining_335
	out_345 = out_336
	sym_346 = sym_337
	x_347 = x_338
	gs_348 = gs_339
	goto b31
b35:
	;
	arg__6716_385, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "pr-str").Deref(), []vm.Value{x_374})
	if callErr != nil {
		return nil, callErr
	}
	arg__6722_389, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "pr-str").Deref(), []vm.Value{x_374})
	if callErr != nil {
		return nil, callErr
	}
	arg__6723_390, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("unsupported destructuring pattern: "), arg__6722_389})
	if callErr != nil {
		return nil, callErr
	}
	arg__6729_394, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "pr-str").Deref(), []vm.Value{x_374})
	if callErr != nil {
		return nil, callErr
	}
	arg__6735_398, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "pr-str").Deref(), []vm.Value{x_374})
	if callErr != nil {
		return nil, callErr
	}
	arg__6736_399, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("unsupported destructuring pattern: "), arg__6735_398})
	if callErr != nil {
		return nil, callErr
	}
	v400, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "throw").Deref(), []vm.Value{arg__6736_399})
	if callErr != nil {
		return nil, callErr
	}
	v404 = v400
	pat_405 = pat_369
	i_406 = i_370
	remaining_407 = remaining_371
	out_408 = out_372
	sym_409 = sym_373
	x_410 = x_374
	goto b37
b36:
	;
	v404 = vm.NIL
	pat_405 = pat_375
	i_406 = i_376
	remaining_407 = remaining_377
	out_408 = out_378
	sym_409 = sym_379
	x_410 = x_380
	goto b37
b37:
	;
	v412 = v404
	pat_413 = pat_405
	i_414 = i_406
	remaining_415 = remaining_407
	out_416 = out_408
	sym_417 = sym_409
	x_418 = x_410
	goto b25
}
func fold_binary_chain(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__6741_10 vm.Value
	var arg__6749_15 vm.Value
	var arg__6755_19 vm.Value
	var arg__6756_20 vm.Value
	var arg__6763_24 vm.Value
	var arg__6771_29 vm.Value
	var arg__6777_33 vm.Value
	var arg__6778_34 vm.Value
	var v36 vm.Value
	var acc_4 vm.Value
	var i_5 int
	var op_kw_6 vm.Value
	var ctx_7 vm.Value
	var args_8 vm.Value
	var v81 vm.Value
	var arg__6784_50 vm.Value
	var v51 bool
	var acc_39 vm.Value
	var i_40 int
	var op_kw_41 vm.Value
	var ctx_42 vm.Value
	var args_43 vm.Value
	var v83 vm.Value
	var acc_44 vm.Value
	var i_45 int
	var op_kw_46 vm.Value
	var ctx_47 vm.Value
	var args_48 vm.Value
	var v82 vm.Value
	var arg__6789_55 vm.Value
	var arg__6798_58 vm.Value
	var arg__6799_59 vm.Value
	var arg__6806_63 vm.Value
	var arg__6815_66 vm.Value
	var arg__6816_67 vm.Value
	var v69 vm.Value
	var v70 int
	var v72 vm.Value
	var acc_73 vm.Value
	var i_74 int
	var op_kw_75 vm.Value
	var ctx_76 vm.Value
	var args_77 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__6741_10, arg__6749_15, arg__6755_19, arg__6756_20, arg__6763_24, arg__6771_29, arg__6777_33, arg__6778_34, v36, acc_4, i_5, op_kw_6, ctx_7, args_8, v81, arg__6784_50, v51, acc_39, i_40, op_kw_41, ctx_42, args_43, v83, acc_44, i_45, op_kw_46, ctx_47, args_48, v82, arg__6789_55, arg__6798_58, arg__6799_59, arg__6806_63, arg__6815_66, arg__6816_67, v69, v70, v72, acc_73, i_74, op_kw_75, ctx_76, args_77
	arg__6741_10, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{arg2})
	if callErr != nil {
		return nil, callErr
	}
	arg__6749_15, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg1, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	arg__6755_19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg1, vm.Int(1)})
	if callErr != nil {
		return nil, callErr
	}
	arg__6756_20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__6749_15, arg__6755_19})
	if callErr != nil {
		return nil, callErr
	}
	arg__6763_24, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{arg2})
	if callErr != nil {
		return nil, callErr
	}
	arg__6771_29, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg1, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	arg__6777_33, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg1, vm.Int(1)})
	if callErr != nil {
		return nil, callErr
	}
	arg__6778_34, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__6771_29, arg__6777_33})
	if callErr != nil {
		return nil, callErr
	}
	v36, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "add-inst!").Deref(), []vm.Value{arg2, arg__6763_24, arg0, arg__6778_34, vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	acc_4 = v36
	i_5 = 2
	op_kw_6 = arg0
	ctx_7 = arg2
	args_8 = arg1
	v81 = vm.NIL
	goto b1
b1:
	;
	arg__6784_50, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{args_8})
	if callErr != nil {
		return nil, callErr
	}
	v51 = rt.GeValue(vm.Int(i_5), arg__6784_50)
	if v51 {
		acc_39 = acc_4
		i_40 = i_5
		op_kw_41 = op_kw_6
		ctx_42 = ctx_7
		args_43 = args_8
		v83 = v81
		goto b2
	} else {
		acc_44 = acc_4
		i_45 = i_5
		op_kw_46 = op_kw_6
		ctx_47 = ctx_7
		args_48 = args_8
		v82 = v81
		goto b3
	}
b2:
	;
	v72 = acc_39
	acc_73 = acc_39
	i_74 = i_40
	op_kw_75 = op_kw_41
	ctx_76 = ctx_42
	args_77 = args_43
	goto b4
b3:
	;
	arg__6789_55, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_47})
	if callErr != nil {
		return nil, callErr
	}
	arg__6798_58, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{args_48, vm.Int(i_45)})
	if callErr != nil {
		return nil, callErr
	}
	arg__6799_59, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{acc_44, arg__6798_58})
	if callErr != nil {
		return nil, callErr
	}
	arg__6806_63, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "ctx-block").Deref(), []vm.Value{ctx_47})
	if callErr != nil {
		return nil, callErr
	}
	arg__6815_66, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{args_48, vm.Int(i_45)})
	if callErr != nil {
		return nil, callErr
	}
	arg__6816_67, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{acc_44, arg__6815_66})
	if callErr != nil {
		return nil, callErr
	}
	v69, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "add-inst!").Deref(), []vm.Value{ctx_47, arg__6806_63, op_kw_46, arg__6816_67, v82})
	if callErr != nil {
		return nil, callErr
	}
	v70 = i_45 + 1
	acc_4 = v69
	i_5 = v70
	op_kw_6 = op_kw_46
	ctx_7 = ctx_47
	args_8 = args_48
	v81 = v82
	goto b1
b4:
	;
	return v72, nil
}
func free_vars(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var v8 vm.Value
	var form_3 vm.Value
	var bound_4 vm.Value
	var v14 vm.Value
	var form_5 vm.Value
	var bound_6 vm.Value
	var or__x_31 vm.Value
	var v629 vm.Value
	var form_630 vm.Value
	var bound_631 vm.Value
	var form_10 vm.Value
	var bound_11 vm.Value
	var v17 vm.Value
	var form_12 vm.Value
	var bound_13 vm.Value
	var v20 vm.Value
	var v22 vm.Value
	var form_23 vm.Value
	var bound_24 vm.Value
	var form_26 vm.Value
	var bound_27 vm.Value
	var head_49 vm.Value
	var v57 bool
	var form_28 vm.Value
	var bound_29 vm.Value
	var v581 vm.Value
	var v625 vm.Value
	var form_626 vm.Value
	var bound_627 vm.Value
	var form_32 vm.Value
	var bound_33 vm.Value
	var or__x_34 vm.Value
	var form_35 vm.Value
	var bound_36 vm.Value
	var or__x_37 vm.Value
	var v41 vm.Value
	var v43 vm.Value
	var form_44 vm.Value
	var bound_45 vm.Value
	var or__x_46 vm.Value
	var form_50 vm.Value
	var bound_51 vm.Value
	var head_52 vm.Value
	var v60 vm.Value
	var form_53 vm.Value
	var bound_54 vm.Value
	var head_55 vm.Value
	var v69 bool
	var v571 vm.Value
	var form_572 vm.Value
	var bound_573 vm.Value
	var head_574 vm.Value
	var form_62 vm.Value
	var bound_63 vm.Value
	var head_64 vm.Value
	var v72 vm.Value
	var form_65 vm.Value
	var bound_66 vm.Value
	var head_67 vm.Value
	var v81 bool
	var v566 vm.Value
	var form_567 vm.Value
	var bound_568 vm.Value
	var head_569 vm.Value
	var form_74 vm.Value
	var bound_75 vm.Value
	var head_76 vm.Value
	var __88 vm.Value
	var _sym_94 vm.Value
	var val_100 vm.Value
	var v102 vm.Value
	var form_77 vm.Value
	var bound_78 vm.Value
	var head_79 vm.Value
	var v111 bool
	var v561 vm.Value
	var form_562 vm.Value
	var bound_563 vm.Value
	var head_564 vm.Value
	var form_104 vm.Value
	var bound_105 vm.Value
	var head_106 vm.Value
	var __118 vm.Value
	var maybe_name_124 vm.Value
	var raw_rest_128 vm.Value
	var has_name_QMARK__130 vm.Value
	var form_107 vm.Value
	var bound_108 vm.Value
	var head_109 vm.Value
	var or__x_394 bool
	var v556 vm.Value
	var form_557 vm.Value
	var bound_558 vm.Value
	var head_559 vm.Value
	var form_131 vm.Value
	var vec__6820_132 vm.Value
	var bound_133 vm.Value
	var head_134 vm.Value
	var __135 vm.Value
	var maybe_name_136 vm.Value
	var raw_rest_137 vm.Value
	var has_name_QMARK__138 vm.Value
	var form_139 vm.Value
	var vec__6820_140 vm.Value
	var bound_141 vm.Value
	var head_142 vm.Value
	var __143 vm.Value
	var maybe_name_144 vm.Value
	var raw_rest_145 vm.Value
	var has_name_QMARK__146 vm.Value
	var name_sym_151 vm.Value
	var form_152 vm.Value
	var vec__6820_153 vm.Value
	var bound_154 vm.Value
	var head_155 vm.Value
	var __156 vm.Value
	var maybe_name_157 vm.Value
	var raw_rest_158 vm.Value
	var has_name_QMARK__159 vm.Value
	var name_sym_160 vm.Value
	var form_161 vm.Value
	var vec__6820_162 vm.Value
	var bound_163 vm.Value
	var head_164 vm.Value
	var __165 vm.Value
	var maybe_name_166 vm.Value
	var raw_rest_167 vm.Value
	var has_name_QMARK__168 vm.Value
	var name_sym_169 vm.Value
	var form_170 vm.Value
	var vec__6820_171 vm.Value
	var bound_172 vm.Value
	var head_173 vm.Value
	var __174 vm.Value
	var maybe_name_175 vm.Value
	var raw_rest_176 vm.Value
	var has_name_QMARK__177 vm.Value
	var v181 vm.Value
	var rest_forms_183 vm.Value
	var name_sym_184 vm.Value
	var form_185 vm.Value
	var vec__6820_186 vm.Value
	var bound_187 vm.Value
	var head_188 vm.Value
	var __189 vm.Value
	var maybe_name_190 vm.Value
	var raw_rest_191 vm.Value
	var has_name_QMARK__192 vm.Value
	var and__x_194 vm.Value
	var rest_forms_195 vm.Value
	var name_sym_196 vm.Value
	var form_197 vm.Value
	var vec__6820_198 vm.Value
	var bound_199 vm.Value
	var head_200 vm.Value
	var __201 vm.Value
	var maybe_name_202 vm.Value
	var raw_rest_203 vm.Value
	var has_name_QMARK__204 vm.Value
	var and__x_205 vm.Value
	var arg__6913_219 vm.Value
	var arg__6918_222 vm.Value
	var v223 vm.Value
	var rest_forms_206 vm.Value
	var name_sym_207 vm.Value
	var form_208 vm.Value
	var vec__6820_209 vm.Value
	var bound_210 vm.Value
	var head_211 vm.Value
	var __212 vm.Value
	var maybe_name_213 vm.Value
	var raw_rest_214 vm.Value
	var has_name_QMARK__215 vm.Value
	var and__x_216 vm.Value
	var multi_QMARK__226 vm.Value
	var rest_forms_227 vm.Value
	var name_sym_228 vm.Value
	var form_229 vm.Value
	var vec__6820_230 vm.Value
	var bound_231 vm.Value
	var head_232 vm.Value
	var __233 vm.Value
	var maybe_name_234 vm.Value
	var raw_rest_235 vm.Value
	var has_name_QMARK__236 vm.Value
	var and__x_237 vm.Value
	var multi_QMARK__238 vm.Value
	var rest_forms_239 vm.Value
	var name_sym_240 vm.Value
	var form_241 vm.Value
	var vec__6820_242 vm.Value
	var bound_243 vm.Value
	var head_244 vm.Value
	var __245 vm.Value
	var maybe_name_246 vm.Value
	var raw_rest_247 vm.Value
	var has_name_QMARK__248 vm.Value
	var arg__6920_262 vm.Value
	var arg__7020_270 vm.Value
	var arg__7023_273 vm.Value
	var arg__7123_281 vm.Value
	var v282 vm.Value
	var multi_QMARK__249 vm.Value
	var rest_forms_250 vm.Value
	var name_sym_251 vm.Value
	var form_252 vm.Value
	var vec__6820_253 vm.Value
	var bound_254 vm.Value
	var head_255 vm.Value
	var __256 vm.Value
	var maybe_name_257 vm.Value
	var raw_rest_258 vm.Value
	var has_name_QMARK__259 vm.Value
	var args_vec_285 vm.Value
	var body_287 vm.Value
	var v374 vm.Value
	var multi_QMARK__375 vm.Value
	var rest_forms_376 vm.Value
	var name_sym_377 vm.Value
	var form_378 vm.Value
	var vec__6820_379 vm.Value
	var bound_380 vm.Value
	var head_381 vm.Value
	var __382 vm.Value
	var maybe_name_383 vm.Value
	var raw_rest_384 vm.Value
	var has_name_QMARK__385 vm.Value
	var multi_QMARK__288 vm.Value
	var rest_forms_289 vm.Value
	var name_sym_290 vm.Value
	var form_291 vm.Value
	var vec__6820_292 vm.Value
	var bound_293 vm.Value
	var head_294 vm.Value
	var __295 vm.Value
	var maybe_name_296 vm.Value
	var raw_rest_297 vm.Value
	var has_name_QMARK__298 vm.Value
	var args_vec_299 vm.Value
	var body_300 vm.Value
	var arg__7131_316 vm.Value
	var arg__7135_319 vm.Value
	var arg__7137_320 vm.Value
	var arg__7141_323 vm.Value
	var arg__7145_326 vm.Value
	var arg__7147_327 vm.Value
	var v328 vm.Value
	var multi_QMARK__301 vm.Value
	var rest_forms_302 vm.Value
	var name_sym_303 vm.Value
	var form_304 vm.Value
	var vec__6820_305 vm.Value
	var bound_306 vm.Value
	var head_307 vm.Value
	var __308 vm.Value
	var maybe_name_309 vm.Value
	var raw_rest_310 vm.Value
	var has_name_QMARK__311 vm.Value
	var args_vec_312 vm.Value
	var body_313 vm.Value
	var arg__7150_331 vm.Value
	var arg__7154_334 vm.Value
	var v335 vm.Value
	var arg_set_337 vm.Value
	var multi_QMARK__338 vm.Value
	var rest_forms_339 vm.Value
	var name_sym_340 vm.Value
	var form_341 vm.Value
	var vec__6820_342 vm.Value
	var bound_343 vm.Value
	var head_344 vm.Value
	var __345 vm.Value
	var maybe_name_346 vm.Value
	var raw_rest_347 vm.Value
	var has_name_QMARK__348 vm.Value
	var args_vec_349 vm.Value
	var body_350 vm.Value
	var arg__7157_352 vm.Value
	var arg__7173_360 vm.Value
	var arg__7176_363 vm.Value
	var arg__7192_371 vm.Value
	var v372 vm.Value
	var form_387 vm.Value
	var bound_388 vm.Value
	var head_389 vm.Value
	var __455 vm.Value
	var bindings_461 vm.Value
	var body_465 vm.Value
	var pairs_469 vm.Value
	var arg__7285_473 vm.Value
	var arg__7287_474 vm.Value
	var arg__7350_479 vm.Value
	var arg__7352_480 vm.Value
	var vec__6822_481 vm.Value
	var used_487 vm.Value
	var new_bound_493 vm.Value
	var arg__7384_501 vm.Value
	var arg__7402_510 vm.Value
	var v511 vm.Value
	var form_390 vm.Value
	var bound_391 vm.Value
	var head_392 vm.Value
	var v551 vm.Value
	var form_552 vm.Value
	var bound_553 vm.Value
	var head_554 vm.Value
	var form_395 vm.Value
	var bound_396 vm.Value
	var head_397 vm.Value
	var or__x_398 bool
	var form_399 vm.Value
	var bound_400 vm.Value
	var head_401 vm.Value
	var or__x_402 bool
	var or__x_406 bool
	var v444 bool
	var form_445 vm.Value
	var bound_446 vm.Value
	var head_447 vm.Value
	var or__x_448 vm.Value
	var form_407 vm.Value
	var bound_408 vm.Value
	var head_409 vm.Value
	var or__x_410 bool
	var form_411 vm.Value
	var bound_412 vm.Value
	var head_413 vm.Value
	var or__x_414 bool
	var or__x_418 bool
	var v438 bool
	var form_439 vm.Value
	var bound_440 vm.Value
	var head_441 vm.Value
	var or__x_442 vm.Value
	var form_419 vm.Value
	var bound_420 vm.Value
	var head_421 vm.Value
	var or__x_422 bool
	var form_423 vm.Value
	var bound_424 vm.Value
	var head_425 vm.Value
	var or__x_426 bool
	var v430 bool
	var v432 bool
	var form_433 vm.Value
	var bound_434 vm.Value
	var head_435 vm.Value
	var or__x_436 vm.Value
	var form_513 vm.Value
	var bound_514 vm.Value
	var head_515 vm.Value
	var arg__7404_522 vm.Value
	var arg__7420_530 vm.Value
	var arg__7423_533 vm.Value
	var arg__7439_541 vm.Value
	var v542 vm.Value
	var form_516 vm.Value
	var bound_517 vm.Value
	var head_518 vm.Value
	var v546 vm.Value
	var form_547 vm.Value
	var bound_548 vm.Value
	var head_549 vm.Value
	var form_576 vm.Value
	var bound_577 vm.Value
	var arg__7444_584 vm.Value
	var arg__7460_592 vm.Value
	var arg__7463_595 vm.Value
	var arg__7479_603 vm.Value
	var v604 vm.Value
	var form_578 vm.Value
	var bound_579 vm.Value
	var v621 vm.Value
	var form_622 vm.Value
	var bound_623 vm.Value
	var form_606 vm.Value
	var bound_607 vm.Value
	var v613 vm.Value
	var form_608 vm.Value
	var bound_609 vm.Value
	var v617 vm.Value
	var form_618 vm.Value
	var bound_619 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = v8, form_3, bound_4, v14, form_5, bound_6, or__x_31, v629, form_630, bound_631, form_10, bound_11, v17, form_12, bound_13, v20, v22, form_23, bound_24, form_26, bound_27, head_49, v57, form_28, bound_29, v581, v625, form_626, bound_627, form_32, bound_33, or__x_34, form_35, bound_36, or__x_37, v41, v43, form_44, bound_45, or__x_46, form_50, bound_51, head_52, v60, form_53, bound_54, head_55, v69, v571, form_572, bound_573, head_574, form_62, bound_63, head_64, v72, form_65, bound_66, head_67, v81, v566, form_567, bound_568, head_569, form_74, bound_75, head_76, __88, _sym_94, val_100, v102, form_77, bound_78, head_79, v111, v561, form_562, bound_563, head_564, form_104, bound_105, head_106, __118, maybe_name_124, raw_rest_128, has_name_QMARK__130, form_107, bound_108, head_109, or__x_394, v556, form_557, bound_558, head_559, form_131, vec__6820_132, bound_133, head_134, __135, maybe_name_136, raw_rest_137, has_name_QMARK__138, form_139, vec__6820_140, bound_141, head_142, __143, maybe_name_144, raw_rest_145, has_name_QMARK__146, name_sym_151, form_152, vec__6820_153, bound_154, head_155, __156, maybe_name_157, raw_rest_158, has_name_QMARK__159, name_sym_160, form_161, vec__6820_162, bound_163, head_164, __165, maybe_name_166, raw_rest_167, has_name_QMARK__168, name_sym_169, form_170, vec__6820_171, bound_172, head_173, __174, maybe_name_175, raw_rest_176, has_name_QMARK__177, v181, rest_forms_183, name_sym_184, form_185, vec__6820_186, bound_187, head_188, __189, maybe_name_190, raw_rest_191, has_name_QMARK__192, and__x_194, rest_forms_195, name_sym_196, form_197, vec__6820_198, bound_199, head_200, __201, maybe_name_202, raw_rest_203, has_name_QMARK__204, and__x_205, arg__6913_219, arg__6918_222, v223, rest_forms_206, name_sym_207, form_208, vec__6820_209, bound_210, head_211, __212, maybe_name_213, raw_rest_214, has_name_QMARK__215, and__x_216, multi_QMARK__226, rest_forms_227, name_sym_228, form_229, vec__6820_230, bound_231, head_232, __233, maybe_name_234, raw_rest_235, has_name_QMARK__236, and__x_237, multi_QMARK__238, rest_forms_239, name_sym_240, form_241, vec__6820_242, bound_243, head_244, __245, maybe_name_246, raw_rest_247, has_name_QMARK__248, arg__6920_262, arg__7020_270, arg__7023_273, arg__7123_281, v282, multi_QMARK__249, rest_forms_250, name_sym_251, form_252, vec__6820_253, bound_254, head_255, __256, maybe_name_257, raw_rest_258, has_name_QMARK__259, args_vec_285, body_287, v374, multi_QMARK__375, rest_forms_376, name_sym_377, form_378, vec__6820_379, bound_380, head_381, __382, maybe_name_383, raw_rest_384, has_name_QMARK__385, multi_QMARK__288, rest_forms_289, name_sym_290, form_291, vec__6820_292, bound_293, head_294, __295, maybe_name_296, raw_rest_297, has_name_QMARK__298, args_vec_299, body_300, arg__7131_316, arg__7135_319, arg__7137_320, arg__7141_323, arg__7145_326, arg__7147_327, v328, multi_QMARK__301, rest_forms_302, name_sym_303, form_304, vec__6820_305, bound_306, head_307, __308, maybe_name_309, raw_rest_310, has_name_QMARK__311, args_vec_312, body_313, arg__7150_331, arg__7154_334, v335, arg_set_337, multi_QMARK__338, rest_forms_339, name_sym_340, form_341, vec__6820_342, bound_343, head_344, __345, maybe_name_346, raw_rest_347, has_name_QMARK__348, args_vec_349, body_350, arg__7157_352, arg__7173_360, arg__7176_363, arg__7192_371, v372, form_387, bound_388, head_389, __455, bindings_461, body_465, pairs_469, arg__7285_473, arg__7287_474, arg__7350_479, arg__7352_480, vec__6822_481, used_487, new_bound_493, arg__7384_501, arg__7402_510, v511, form_390, bound_391, head_392, v551, form_552, bound_553, head_554, form_395, bound_396, head_397, or__x_398, form_399, bound_400, head_401, or__x_402, or__x_406, v444, form_445, bound_446, head_447, or__x_448, form_407, bound_408, head_409, or__x_410, form_411, bound_412, head_413, or__x_414, or__x_418, v438, form_439, bound_440, head_441, or__x_442, form_419, bound_420, head_421, or__x_422, form_423, bound_424, head_425, or__x_426, v430, v432, form_433, bound_434, head_435, or__x_436, form_513, bound_514, head_515, arg__7404_522, arg__7420_530, arg__7423_533, arg__7439_541, v542, form_516, bound_517, head_518, v546, form_547, bound_548, head_549, form_576, bound_577, arg__7444_584, arg__7460_592, arg__7463_595, arg__7479_603, v604, form_578, bound_579, v621, form_622, bound_623, form_606, bound_607, v613, form_608, bound_609, v617, form_618, bound_619
	v8, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "symbol?").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v8) {
		form_3 = arg0
		bound_4 = arg1
		goto b1
	} else {
		form_5 = arg0
		bound_6 = arg1
		goto b2
	}
b1:
	;
	v14, callErr = rt.InvokeValue(bound_4, []vm.Value{form_3})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v14) {
		form_10 = form_3
		bound_11 = bound_4
		goto b4
	} else {
		form_12 = form_3
		bound_13 = bound_4
		goto b5
	}
b2:
	;
	or__x_31, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "list?").Deref(), []vm.Value{form_5})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(or__x_31) {
		form_32 = form_5
		bound_33 = bound_6
		or__x_34 = or__x_31
		goto b10
	} else {
		form_35 = form_5
		bound_36 = bound_6
		or__x_37 = or__x_31
		goto b11
	}
b3:
	;
	return v629, nil
b4:
	;
	v17, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	v22 = v17
	form_23 = form_10
	bound_24 = bound_11
	goto b6
b5:
	;
	v20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{form_12})
	if callErr != nil {
		return nil, callErr
	}
	v22 = v20
	form_23 = form_12
	bound_24 = bound_13
	goto b6
b6:
	;
	v629 = v22
	form_630 = form_23
	bound_631 = bound_24
	goto b3
b7:
	;
	head_49, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{form_26})
	if callErr != nil {
		return nil, callErr
	}
	v57 = head_49 == vm.Symbol("quote")
	if v57 {
		form_50 = form_26
		bound_51 = bound_27
		head_52 = head_49
		goto b13
	} else {
		form_53 = form_26
		bound_54 = bound_27
		head_55 = head_49
		goto b14
	}
b8:
	;
	v581, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector?").Deref(), []vm.Value{form_28})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v581) {
		form_576 = form_28
		bound_577 = bound_29
		goto b55
	} else {
		form_578 = form_28
		bound_579 = bound_29
		goto b56
	}
b9:
	;
	v629 = v625
	form_630 = form_626
	bound_631 = bound_627
	goto b3
b10:
	;
	v43 = or__x_34
	form_44 = form_32
	bound_45 = bound_33
	or__x_46 = or__x_34
	goto b12
b11:
	;
	v41, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq?").Deref(), []vm.Value{form_35})
	if callErr != nil {
		return nil, callErr
	}
	v43 = v41
	form_44 = form_35
	bound_45 = bound_36
	or__x_46 = or__x_37
	goto b12
b12:
	;
	if vm.IsTruthy(v43) {
		form_26 = form_44
		bound_27 = bound_45
		goto b7
	} else {
		form_28 = form_44
		bound_29 = bound_45
		goto b8
	}
b13:
	;
	v60, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	v571 = v60
	form_572 = form_50
	bound_573 = bound_51
	head_574 = head_52
	goto b15
b14:
	;
	v69 = head_55 == vm.Symbol("var")
	if v69 {
		form_62 = form_53
		bound_63 = bound_54
		head_64 = head_55
		goto b16
	} else {
		form_65 = form_53
		bound_66 = bound_54
		head_67 = head_55
		goto b17
	}
b15:
	;
	v625 = v571
	form_626 = form_572
	bound_627 = bound_573
	goto b9
b16:
	;
	v72, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	v566 = v72
	form_567 = form_62
	bound_568 = bound_63
	head_569 = head_64
	goto b18
b17:
	;
	v81 = head_67 == vm.Symbol("set!")
	if v81 {
		form_74 = form_65
		bound_75 = bound_66
		head_76 = head_67
		goto b19
	} else {
		form_77 = form_65
		bound_78 = bound_66
		head_79 = head_67
		goto b20
	}
b18:
	;
	v571 = v566
	form_572 = form_567
	bound_573 = bound_568
	head_574 = head_569
	goto b15
b19:
	;
	__88, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{form_74, vm.Int(0), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	_sym_94, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{form_74, vm.Int(1), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	val_100, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{form_74, vm.Int(2), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	v102, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "free-vars").Deref(), []vm.Value{val_100, bound_75})
	if callErr != nil {
		return nil, callErr
	}
	v561 = v102
	form_562 = form_74
	bound_563 = bound_75
	head_564 = head_76
	goto b21
b20:
	;
	v111 = head_79 == vm.Symbol("fn*")
	if v111 {
		form_104 = form_77
		bound_105 = bound_78
		head_106 = head_79
		goto b22
	} else {
		form_107 = form_77
		bound_108 = bound_78
		head_109 = head_79
		goto b23
	}
b21:
	;
	v566 = v561
	form_567 = form_562
	bound_568 = bound_563
	head_569 = head_564
	goto b18
b22:
	;
	__118, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{form_104, vm.Int(0), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	maybe_name_124, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{form_104, vm.Int(1), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	raw_rest_128, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "drop").Deref(), []vm.Value{vm.Int(2), form_104})
	if callErr != nil {
		return nil, callErr
	}
	has_name_QMARK__130, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "symbol?").Deref(), []vm.Value{maybe_name_124})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(has_name_QMARK__130) {
		form_131 = form_104
		vec__6820_132 = form_104
		bound_133 = bound_105
		head_134 = head_106
		__135 = __118
		maybe_name_136 = maybe_name_124
		raw_rest_137 = raw_rest_128
		has_name_QMARK__138 = has_name_QMARK__130
		goto b25
	} else {
		form_139 = form_104
		vec__6820_140 = form_104
		bound_141 = bound_105
		head_142 = head_106
		__143 = __118
		maybe_name_144 = maybe_name_124
		raw_rest_145 = raw_rest_128
		has_name_QMARK__146 = has_name_QMARK__130
		goto b26
	}
b23:
	;
	or__x_394 = head_109 == vm.Symbol("let")
	if or__x_394 {
		form_395 = form_107
		bound_396 = bound_108
		head_397 = head_109
		or__x_398 = or__x_394
		goto b43
	} else {
		form_399 = form_107
		bound_400 = bound_108
		head_401 = head_109
		or__x_402 = or__x_394
		goto b44
	}
b24:
	;
	v561 = v556
	form_562 = form_557
	bound_563 = bound_558
	head_564 = head_559
	goto b21
b25:
	;
	name_sym_151 = maybe_name_136
	form_152 = form_131
	vec__6820_153 = vec__6820_132
	bound_154 = bound_133
	head_155 = head_134
	__156 = __135
	maybe_name_157 = maybe_name_136
	raw_rest_158 = raw_rest_137
	has_name_QMARK__159 = has_name_QMARK__138
	goto b27
b26:
	;
	name_sym_151 = vm.NIL
	form_152 = form_139
	vec__6820_153 = vec__6820_140
	bound_154 = bound_141
	head_155 = head_142
	__156 = __143
	maybe_name_157 = maybe_name_144
	raw_rest_158 = raw_rest_145
	has_name_QMARK__159 = has_name_QMARK__146
	goto b27
b27:
	;
	if vm.IsTruthy(has_name_QMARK__159) {
		name_sym_160 = name_sym_151
		form_161 = form_152
		vec__6820_162 = vec__6820_153
		bound_163 = bound_154
		head_164 = head_155
		__165 = __156
		maybe_name_166 = maybe_name_157
		raw_rest_167 = raw_rest_158
		has_name_QMARK__168 = has_name_QMARK__159
		goto b28
	} else {
		name_sym_169 = name_sym_151
		form_170 = form_152
		vec__6820_171 = vec__6820_153
		bound_172 = bound_154
		head_173 = head_155
		__174 = __156
		maybe_name_175 = maybe_name_157
		raw_rest_176 = raw_rest_158
		has_name_QMARK__177 = has_name_QMARK__159
		goto b29
	}
b28:
	;
	rest_forms_183 = raw_rest_167
	name_sym_184 = name_sym_160
	form_185 = form_161
	vec__6820_186 = vec__6820_162
	bound_187 = bound_163
	head_188 = head_164
	__189 = __165
	maybe_name_190 = maybe_name_166
	raw_rest_191 = raw_rest_167
	has_name_QMARK__192 = has_name_QMARK__168
	goto b30
b29:
	;
	v181, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "cons").Deref(), []vm.Value{maybe_name_175, raw_rest_176})
	if callErr != nil {
		return nil, callErr
	}
	rest_forms_183 = v181
	name_sym_184 = name_sym_169
	form_185 = form_170
	vec__6820_186 = vec__6820_171
	bound_187 = bound_172
	head_188 = head_173
	__189 = __174
	maybe_name_190 = maybe_name_175
	raw_rest_191 = raw_rest_176
	has_name_QMARK__192 = has_name_QMARK__177
	goto b30
b30:
	;
	and__x_194, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{rest_forms_183})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(and__x_194) {
		rest_forms_195 = rest_forms_183
		name_sym_196 = name_sym_184
		form_197 = form_185
		vec__6820_198 = vec__6820_186
		bound_199 = bound_187
		head_200 = head_188
		__201 = __189
		maybe_name_202 = maybe_name_190
		raw_rest_203 = raw_rest_191
		has_name_QMARK__204 = has_name_QMARK__192
		and__x_205 = and__x_194
		goto b31
	} else {
		rest_forms_206 = rest_forms_183
		name_sym_207 = name_sym_184
		form_208 = form_185
		vec__6820_209 = vec__6820_186
		bound_210 = bound_187
		head_211 = head_188
		__212 = __189
		maybe_name_213 = maybe_name_190
		raw_rest_214 = raw_rest_191
		has_name_QMARK__215 = has_name_QMARK__192
		and__x_216 = and__x_194
		goto b32
	}
b31:
	;
	arg__6913_219, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{rest_forms_195})
	if callErr != nil {
		return nil, callErr
	}
	arg__6918_222, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{rest_forms_195})
	if callErr != nil {
		return nil, callErr
	}
	v223, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "list?").Deref(), []vm.Value{arg__6918_222})
	if callErr != nil {
		return nil, callErr
	}
	multi_QMARK__226 = v223
	rest_forms_227 = rest_forms_195
	name_sym_228 = name_sym_196
	form_229 = form_197
	vec__6820_230 = vec__6820_198
	bound_231 = bound_199
	head_232 = head_200
	__233 = __201
	maybe_name_234 = maybe_name_202
	raw_rest_235 = raw_rest_203
	has_name_QMARK__236 = has_name_QMARK__204
	and__x_237 = and__x_205
	goto b33
b32:
	;
	multi_QMARK__226 = and__x_216
	rest_forms_227 = rest_forms_206
	name_sym_228 = name_sym_207
	form_229 = form_208
	vec__6820_230 = vec__6820_209
	bound_231 = bound_210
	head_232 = head_211
	__233 = __212
	maybe_name_234 = maybe_name_213
	raw_rest_235 = raw_rest_214
	has_name_QMARK__236 = has_name_QMARK__215
	and__x_237 = and__x_216
	goto b33
b33:
	;
	if vm.IsTruthy(multi_QMARK__226) {
		multi_QMARK__238 = multi_QMARK__226
		rest_forms_239 = rest_forms_227
		name_sym_240 = name_sym_228
		form_241 = form_229
		vec__6820_242 = vec__6820_230
		bound_243 = bound_231
		head_244 = head_232
		__245 = __233
		maybe_name_246 = maybe_name_234
		raw_rest_247 = raw_rest_235
		has_name_QMARK__248 = has_name_QMARK__236
		goto b34
	} else {
		multi_QMARK__249 = multi_QMARK__226
		rest_forms_250 = rest_forms_227
		name_sym_251 = name_sym_228
		form_252 = form_229
		vec__6820_253 = vec__6820_230
		bound_254 = bound_231
		head_255 = head_232
		__256 = __233
		maybe_name_257 = maybe_name_234
		raw_rest_258 = raw_rest_235
		has_name_QMARK__259 = has_name_QMARK__236
		goto b35
	}
b34:
	;
	arg__6920_262, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__7020_270, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapcat").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var args_vec_3 vm.Value
		var body_5 vm.Value
		var arity_form_6 vm.Value
		var name_sym_7 vm.Value
		var args_vec_8 vm.Value
		var body_9 vm.Value
		var arg__6978_16 vm.Value
		var arg__6982_19 vm.Value
		var arg__6984_20 vm.Value
		var arg__6988_23 vm.Value
		var arg__6992_26 vm.Value
		var arg__6994_27 vm.Value
		var v28 vm.Value
		var arity_form_10 vm.Value
		var name_sym_11 vm.Value
		var args_vec_12 vm.Value
		var body_13 vm.Value
		var arg__6997_31 vm.Value
		var arg__7001_34 vm.Value
		var v35 vm.Value
		var arg_set_37 vm.Value
		var arity_form_38 vm.Value
		var name_sym_39 vm.Value
		var args_vec_40 vm.Value
		var body_41 vm.Value
		var v49 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = args_vec_3, body_5, arity_form_6, name_sym_7, args_vec_8, body_9, arg__6978_16, arg__6982_19, arg__6984_20, arg__6988_23, arg__6992_26, arg__6994_27, v28, arity_form_10, name_sym_11, args_vec_12, body_13, arg__6997_31, arg__7001_34, v35, arg_set_37, arity_form_38, name_sym_39, args_vec_40, body_41, v49
		args_vec_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		body_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		if vm.IsTruthy(name_sym_240) {
			arity_form_6 = arg0
			name_sym_7 = name_sym_240
			args_vec_8 = args_vec_3
			body_9 = body_5
			goto b1
		} else {
			arity_form_10 = arg0
			name_sym_11 = name_sym_240
			args_vec_12 = args_vec_3
			body_13 = body_5
			goto b2
		}
	b1:
		;
		arg__6978_16, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
		if callErr != nil {
			return nil, callErr
		}
		arg__6982_19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
		if callErr != nil {
			return nil, callErr
		}
		arg__6984_20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__6982_19, args_vec_8})
		if callErr != nil {
			return nil, callErr
		}
		arg__6988_23, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
		if callErr != nil {
			return nil, callErr
		}
		arg__6992_26, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
		if callErr != nil {
			return nil, callErr
		}
		arg__6994_27, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__6992_26, args_vec_8})
		if callErr != nil {
			return nil, callErr
		}
		v28, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{arg__6994_27, name_sym_7})
		if callErr != nil {
			return nil, callErr
		}
		arg_set_37 = v28
		arity_form_38 = arity_form_6
		name_sym_39 = name_sym_7
		args_vec_40 = args_vec_8
		body_41 = body_9
		goto b3
	b2:
		;
		arg__6997_31, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
		if callErr != nil {
			return nil, callErr
		}
		arg__7001_34, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
		if callErr != nil {
			return nil, callErr
		}
		v35, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__7001_34, args_vec_12})
		if callErr != nil {
			return nil, callErr
		}
		arg_set_37 = v35
		arity_form_38 = arity_form_10
		name_sym_39 = name_sym_11
		args_vec_40 = args_vec_12
		body_41 = body_13
		goto b3
	b3:
		;
		v49, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapcat").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
			var v3 vm.Value
			var callErr error
			_ = v3
			v3, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "free-vars").Deref(), []vm.Value{arg0, arg_set_37})
			if callErr != nil {
				return nil, callErr
			}
			return v3, nil
		}), body_41})
		if callErr != nil {
			return nil, callErr
		}
		return v49, nil
	}), rest_forms_239})
	if callErr != nil {
		return nil, callErr
	}
	arg__7023_273, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__7123_281, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapcat").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var args_vec_3 vm.Value
		var body_5 vm.Value
		var arity_form_6 vm.Value
		var name_sym_7 vm.Value
		var args_vec_8 vm.Value
		var body_9 vm.Value
		var arg__7081_16 vm.Value
		var arg__7085_19 vm.Value
		var arg__7087_20 vm.Value
		var arg__7091_23 vm.Value
		var arg__7095_26 vm.Value
		var arg__7097_27 vm.Value
		var v28 vm.Value
		var arity_form_10 vm.Value
		var name_sym_11 vm.Value
		var args_vec_12 vm.Value
		var body_13 vm.Value
		var arg__7100_31 vm.Value
		var arg__7104_34 vm.Value
		var v35 vm.Value
		var arg_set_37 vm.Value
		var arity_form_38 vm.Value
		var name_sym_39 vm.Value
		var args_vec_40 vm.Value
		var body_41 vm.Value
		var v49 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = args_vec_3, body_5, arity_form_6, name_sym_7, args_vec_8, body_9, arg__7081_16, arg__7085_19, arg__7087_20, arg__7091_23, arg__7095_26, arg__7097_27, v28, arity_form_10, name_sym_11, args_vec_12, body_13, arg__7100_31, arg__7104_34, v35, arg_set_37, arity_form_38, name_sym_39, args_vec_40, body_41, v49
		args_vec_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		body_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		if vm.IsTruthy(name_sym_240) {
			arity_form_6 = arg0
			name_sym_7 = name_sym_240
			args_vec_8 = args_vec_3
			body_9 = body_5
			goto b1
		} else {
			arity_form_10 = arg0
			name_sym_11 = name_sym_240
			args_vec_12 = args_vec_3
			body_13 = body_5
			goto b2
		}
	b1:
		;
		arg__7081_16, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
		if callErr != nil {
			return nil, callErr
		}
		arg__7085_19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
		if callErr != nil {
			return nil, callErr
		}
		arg__7087_20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__7085_19, args_vec_8})
		if callErr != nil {
			return nil, callErr
		}
		arg__7091_23, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
		if callErr != nil {
			return nil, callErr
		}
		arg__7095_26, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
		if callErr != nil {
			return nil, callErr
		}
		arg__7097_27, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__7095_26, args_vec_8})
		if callErr != nil {
			return nil, callErr
		}
		v28, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{arg__7097_27, name_sym_7})
		if callErr != nil {
			return nil, callErr
		}
		arg_set_37 = v28
		arity_form_38 = arity_form_6
		name_sym_39 = name_sym_7
		args_vec_40 = args_vec_8
		body_41 = body_9
		goto b3
	b2:
		;
		arg__7100_31, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
		if callErr != nil {
			return nil, callErr
		}
		arg__7104_34, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
		if callErr != nil {
			return nil, callErr
		}
		v35, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__7104_34, args_vec_12})
		if callErr != nil {
			return nil, callErr
		}
		arg_set_37 = v35
		arity_form_38 = arity_form_10
		name_sym_39 = name_sym_11
		args_vec_40 = args_vec_12
		body_41 = body_13
		goto b3
	b3:
		;
		v49, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapcat").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
			var v3 vm.Value
			var callErr error
			_ = v3
			v3, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "free-vars").Deref(), []vm.Value{arg0, arg_set_37})
			if callErr != nil {
				return nil, callErr
			}
			return v3, nil
		}), body_41})
		if callErr != nil {
			return nil, callErr
		}
		return v49, nil
	}), rest_forms_239})
	if callErr != nil {
		return nil, callErr
	}
	v282, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__7023_273, arg__7123_281})
	if callErr != nil {
		return nil, callErr
	}
	v374 = v282
	multi_QMARK__375 = multi_QMARK__238
	rest_forms_376 = rest_forms_239
	name_sym_377 = name_sym_240
	form_378 = form_241
	vec__6820_379 = vec__6820_242
	bound_380 = bound_243
	head_381 = head_244
	__382 = __245
	maybe_name_383 = maybe_name_246
	raw_rest_384 = raw_rest_247
	has_name_QMARK__385 = has_name_QMARK__248
	goto b36
b35:
	;
	args_vec_285, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{rest_forms_250})
	if callErr != nil {
		return nil, callErr
	}
	body_287, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{rest_forms_250})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(name_sym_251) {
		multi_QMARK__288 = multi_QMARK__249
		rest_forms_289 = rest_forms_250
		name_sym_290 = name_sym_251
		form_291 = form_252
		vec__6820_292 = vec__6820_253
		bound_293 = bound_254
		head_294 = head_255
		__295 = __256
		maybe_name_296 = maybe_name_257
		raw_rest_297 = raw_rest_258
		has_name_QMARK__298 = has_name_QMARK__259
		args_vec_299 = args_vec_285
		body_300 = body_287
		goto b37
	} else {
		multi_QMARK__301 = multi_QMARK__249
		rest_forms_302 = rest_forms_250
		name_sym_303 = name_sym_251
		form_304 = form_252
		vec__6820_305 = vec__6820_253
		bound_306 = bound_254
		head_307 = head_255
		__308 = __256
		maybe_name_309 = maybe_name_257
		raw_rest_310 = raw_rest_258
		has_name_QMARK__311 = has_name_QMARK__259
		args_vec_312 = args_vec_285
		body_313 = body_287
		goto b38
	}
b36:
	;
	v556 = v374
	form_557 = form_378
	bound_558 = bound_380
	head_559 = head_381
	goto b24
b37:
	;
	arg__7131_316, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__7135_319, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__7137_320, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__7135_319, args_vec_299})
	if callErr != nil {
		return nil, callErr
	}
	arg__7141_323, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__7145_326, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__7147_327, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__7145_326, args_vec_299})
	if callErr != nil {
		return nil, callErr
	}
	v328, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{arg__7147_327, name_sym_290})
	if callErr != nil {
		return nil, callErr
	}
	arg_set_337 = v328
	multi_QMARK__338 = multi_QMARK__288
	rest_forms_339 = rest_forms_289
	name_sym_340 = name_sym_290
	form_341 = form_291
	vec__6820_342 = vec__6820_292
	bound_343 = bound_293
	head_344 = head_294
	__345 = __295
	maybe_name_346 = maybe_name_296
	raw_rest_347 = raw_rest_297
	has_name_QMARK__348 = has_name_QMARK__298
	args_vec_349 = args_vec_299
	body_350 = body_300
	goto b39
b38:
	;
	arg__7150_331, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__7154_334, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	v335, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__7154_334, args_vec_312})
	if callErr != nil {
		return nil, callErr
	}
	arg_set_337 = v335
	multi_QMARK__338 = multi_QMARK__301
	rest_forms_339 = rest_forms_302
	name_sym_340 = name_sym_303
	form_341 = form_304
	vec__6820_342 = vec__6820_305
	bound_343 = bound_306
	head_344 = head_307
	__345 = __308
	maybe_name_346 = maybe_name_309
	raw_rest_347 = raw_rest_310
	has_name_QMARK__348 = has_name_QMARK__311
	args_vec_349 = args_vec_312
	body_350 = body_313
	goto b39
b39:
	;
	arg__7157_352, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__7173_360, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapcat").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "free-vars").Deref(), []vm.Value{arg0, arg_set_337})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), body_350})
	if callErr != nil {
		return nil, callErr
	}
	arg__7176_363, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__7192_371, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapcat").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "free-vars").Deref(), []vm.Value{arg0, arg_set_337})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), body_350})
	if callErr != nil {
		return nil, callErr
	}
	v372, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__7176_363, arg__7192_371})
	if callErr != nil {
		return nil, callErr
	}
	v374 = v372
	multi_QMARK__375 = multi_QMARK__338
	rest_forms_376 = rest_forms_339
	name_sym_377 = name_sym_340
	form_378 = form_341
	vec__6820_379 = vec__6820_342
	bound_380 = bound_343
	head_381 = head_344
	__382 = __345
	maybe_name_383 = maybe_name_346
	raw_rest_384 = raw_rest_347
	has_name_QMARK__385 = has_name_QMARK__348
	goto b36
b40:
	;
	__455, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{form_387, vm.Int(0), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	bindings_461, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{form_387, vm.Int(1), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	body_465, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "drop").Deref(), []vm.Value{vm.Int(2), form_387})
	if callErr != nil {
		return nil, callErr
	}
	pairs_469, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "partition").Deref(), []vm.Value{vm.Int(2), bindings_461})
	if callErr != nil {
		return nil, callErr
	}
	arg__7285_473, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__7287_474, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__7285_473, bound_388})
	if callErr != nil {
		return nil, callErr
	}
	arg__7350_479, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__7352_480, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__7350_479, bound_388})
	if callErr != nil {
		return nil, callErr
	}
	vec__6822_481, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "reduce").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
		var caps_7 vm.Value
		var b_13 vm.Value
		var sym_19 vm.Value
		var init_25 vm.Value
		var arg__7325_28 vm.Value
		var arg__7333_31 vm.Value
		var arg__7334_32 vm.Value
		var arg__7339_34 vm.Value
		var arg__7345_37 vm.Value
		var arg__7346_38 vm.Value
		var v39 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _, _, _, _ = caps_7, b_13, sym_19, init_25, arg__7325_28, arg__7333_31, arg__7334_32, arg__7339_34, arg__7345_37, arg__7346_38, v39
		caps_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(0), vm.NIL})
		if callErr != nil {
			return nil, callErr
		}
		b_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(1), vm.NIL})
		if callErr != nil {
			return nil, callErr
		}
		sym_19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg1, vm.Int(0), vm.NIL})
		if callErr != nil {
			return nil, callErr
		}
		init_25, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg1, vm.Int(1), vm.NIL})
		if callErr != nil {
			return nil, callErr
		}
		arg__7325_28, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "free-vars").Deref(), []vm.Value{init_25, b_13})
		if callErr != nil {
			return nil, callErr
		}
		arg__7333_31, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "free-vars").Deref(), []vm.Value{init_25, b_13})
		if callErr != nil {
			return nil, callErr
		}
		arg__7334_32, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{caps_7, arg__7333_31})
		if callErr != nil {
			return nil, callErr
		}
		arg__7339_34, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "binding-syms").Deref(), []vm.Value{sym_19})
		if callErr != nil {
			return nil, callErr
		}
		arg__7345_37, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "binding-syms").Deref(), []vm.Value{sym_19})
		if callErr != nil {
			return nil, callErr
		}
		arg__7346_38, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{b_13, arg__7345_37})
		if callErr != nil {
			return nil, callErr
		}
		v39, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__7334_32, arg__7346_38})
		if callErr != nil {
			return nil, callErr
		}
		return v39, nil
	}), arg__7352_480, pairs_469})
	if callErr != nil {
		return nil, callErr
	}
	used_487, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{vec__6822_481, vm.Int(0), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	new_bound_493, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{vec__6822_481, vm.Int(1), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	arg__7384_501, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapcat").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "free-vars").Deref(), []vm.Value{arg0, new_bound_493})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), body_465})
	if callErr != nil {
		return nil, callErr
	}
	arg__7402_510, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapcat").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "free-vars").Deref(), []vm.Value{arg0, new_bound_493})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), body_465})
	if callErr != nil {
		return nil, callErr
	}
	v511, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{used_487, arg__7402_510})
	if callErr != nil {
		return nil, callErr
	}
	v551 = v511
	form_552 = form_387
	bound_553 = bound_388
	head_554 = head_389
	goto b42
b41:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		form_513 = form_390
		bound_514 = bound_391
		head_515 = head_392
		goto b52
	} else {
		form_516 = form_390
		bound_517 = bound_391
		head_518 = head_392
		goto b53
	}
b42:
	;
	v556 = v551
	form_557 = form_552
	bound_558 = bound_553
	head_559 = head_554
	goto b24
b43:
	;
	v444 = or__x_398
	form_445 = form_395
	bound_446 = bound_396
	head_447 = head_397
	or__x_448 = vm.Boolean(or__x_398)
	goto b45
b44:
	;
	or__x_406 = head_401 == vm.Symbol("let*")
	if or__x_406 {
		form_407 = form_399
		bound_408 = bound_400
		head_409 = head_401
		or__x_410 = or__x_406
		goto b46
	} else {
		form_411 = form_399
		bound_412 = bound_400
		head_413 = head_401
		or__x_414 = or__x_406
		goto b47
	}
b45:
	;
	if v444 {
		form_387 = form_445
		bound_388 = bound_446
		head_389 = head_447
		goto b40
	} else {
		form_390 = form_445
		bound_391 = bound_446
		head_392 = head_447
		goto b41
	}
b46:
	;
	v438 = or__x_410
	form_439 = form_407
	bound_440 = bound_408
	head_441 = head_409
	or__x_442 = vm.Boolean(or__x_410)
	goto b48
b47:
	;
	or__x_418 = head_413 == vm.Symbol("loop")
	if or__x_418 {
		form_419 = form_411
		bound_420 = bound_412
		head_421 = head_413
		or__x_422 = or__x_418
		goto b49
	} else {
		form_423 = form_411
		bound_424 = bound_412
		head_425 = head_413
		or__x_426 = or__x_418
		goto b50
	}
b48:
	;
	v444 = v438
	form_445 = form_439
	bound_446 = bound_440
	head_447 = head_441
	or__x_448 = vm.Boolean(or__x_402)
	goto b45
b49:
	;
	v432 = or__x_422
	form_433 = form_419
	bound_434 = bound_420
	head_435 = head_421
	or__x_436 = vm.Boolean(or__x_422)
	goto b51
b50:
	;
	v430 = head_425 == vm.Symbol("loop*")
	v432 = v430
	form_433 = form_423
	bound_434 = bound_424
	head_435 = head_425
	or__x_436 = vm.Boolean(or__x_426)
	goto b51
b51:
	;
	v438 = v432
	form_439 = form_433
	bound_440 = bound_434
	head_441 = head_435
	or__x_442 = vm.Boolean(or__x_414)
	goto b48
b52:
	;
	arg__7404_522, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__7420_530, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapcat").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "free-vars").Deref(), []vm.Value{arg0, bound_514})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), form_513})
	if callErr != nil {
		return nil, callErr
	}
	arg__7423_533, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__7439_541, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapcat").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "free-vars").Deref(), []vm.Value{arg0, bound_514})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), form_513})
	if callErr != nil {
		return nil, callErr
	}
	v542, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__7423_533, arg__7439_541})
	if callErr != nil {
		return nil, callErr
	}
	v546 = v542
	form_547 = form_513
	bound_548 = bound_514
	head_549 = head_515
	goto b54
b53:
	;
	v546 = vm.NIL
	form_547 = form_516
	bound_548 = bound_517
	head_549 = head_518
	goto b54
b54:
	;
	v551 = v546
	form_552 = form_547
	bound_553 = bound_548
	head_554 = head_549
	goto b42
b55:
	;
	arg__7444_584, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__7460_592, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapcat").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "free-vars").Deref(), []vm.Value{arg0, bound_577})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), form_576})
	if callErr != nil {
		return nil, callErr
	}
	arg__7463_595, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__7479_603, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapcat").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "free-vars").Deref(), []vm.Value{arg0, bound_577})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), form_576})
	if callErr != nil {
		return nil, callErr
	}
	v604, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__7463_595, arg__7479_603})
	if callErr != nil {
		return nil, callErr
	}
	v621 = v604
	form_622 = form_576
	bound_623 = bound_577
	goto b57
b56:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		form_606 = form_578
		bound_607 = bound_579
		goto b58
	} else {
		form_608 = form_578
		bound_609 = bound_579
		goto b59
	}
b57:
	;
	v625 = v621
	form_626 = form_622
	bound_627 = bound_623
	goto b9
b58:
	;
	v613, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	v617 = v613
	form_618 = form_606
	bound_619 = bound_607
	goto b60
b59:
	;
	v617 = vm.NIL
	form_618 = form_608
	bound_619 = bound_609
	goto b60
b60:
	;
	v621 = v617
	form_622 = form_618
	bound_623 = bound_619
	goto b57
}
func is_literal_QMARK_(arg0 vm.Value) (vm.Value, error) {
	var or__x_3 vm.Value
	var v_4 vm.Value
	var or__x_5 vm.Value
	var v_6 vm.Value
	var or__x_7 vm.Value
	var or__x_11 vm.Value
	var v61 vm.Value
	var v_62 vm.Value
	var or__x_63 vm.Value
	var v_12 vm.Value
	var or__x_13 vm.Value
	var v_14 vm.Value
	var or__x_15 vm.Value
	var or__x_19 vm.Value
	var v57 vm.Value
	var v_58 vm.Value
	var or__x_59 vm.Value
	var v_20 vm.Value
	var or__x_21 vm.Value
	var v_22 vm.Value
	var or__x_23 vm.Value
	var or__x_27 vm.Value
	var v53 vm.Value
	var v_54 vm.Value
	var or__x_55 vm.Value
	var v_28 vm.Value
	var or__x_29 vm.Value
	var v_30 vm.Value
	var or__x_31 vm.Value
	var or__x_35 vm.Value
	var v49 vm.Value
	var v_50 vm.Value
	var or__x_51 vm.Value
	var v_36 vm.Value
	var or__x_37 vm.Value
	var v_38 vm.Value
	var or__x_39 vm.Value
	var v43 vm.Value
	var v45 vm.Value
	var v_46 vm.Value
	var or__x_47 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = or__x_3, v_4, or__x_5, v_6, or__x_7, or__x_11, v61, v_62, or__x_63, v_12, or__x_13, v_14, or__x_15, or__x_19, v57, v_58, or__x_59, v_20, or__x_21, v_22, or__x_23, or__x_27, v53, v_54, or__x_55, v_28, or__x_29, v_30, or__x_31, or__x_35, v49, v_50, or__x_51, v_36, or__x_37, v_38, or__x_39, v43, v45, v_46, or__x_47
	or__x_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(or__x_3) {
		v_4 = arg0
		or__x_5 = or__x_3
		goto b1
	} else {
		v_6 = arg0
		or__x_7 = or__x_3
		goto b2
	}
b1:
	;
	v61 = or__x_5
	v_62 = v_4
	or__x_63 = or__x_5
	goto b3
b2:
	;
	or__x_11, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "integer?").Deref(), []vm.Value{v_6})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(or__x_11) {
		v_12 = v_6
		or__x_13 = or__x_11
		goto b4
	} else {
		v_14 = v_6
		or__x_15 = or__x_11
		goto b5
	}
b3:
	;
	return v61, nil
b4:
	;
	v57 = or__x_13
	v_58 = v_12
	or__x_59 = or__x_13
	goto b6
b5:
	;
	or__x_19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "float?").Deref(), []vm.Value{v_14})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(or__x_19) {
		v_20 = v_14
		or__x_21 = or__x_19
		goto b7
	} else {
		v_22 = v_14
		or__x_23 = or__x_19
		goto b8
	}
b6:
	;
	v61 = v57
	v_62 = v_58
	or__x_63 = or__x_7
	goto b3
b7:
	;
	v53 = or__x_21
	v_54 = v_20
	or__x_55 = or__x_21
	goto b9
b8:
	;
	or__x_27, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "string?").Deref(), []vm.Value{v_22})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(or__x_27) {
		v_28 = v_22
		or__x_29 = or__x_27
		goto b10
	} else {
		v_30 = v_22
		or__x_31 = or__x_27
		goto b11
	}
b9:
	;
	v57 = v53
	v_58 = v_54
	or__x_59 = or__x_15
	goto b6
b10:
	;
	v49 = or__x_29
	v_50 = v_28
	or__x_51 = or__x_29
	goto b12
b11:
	;
	or__x_35, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "keyword?").Deref(), []vm.Value{v_30})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(or__x_35) {
		v_36 = v_30
		or__x_37 = or__x_35
		goto b13
	} else {
		v_38 = v_30
		or__x_39 = or__x_35
		goto b14
	}
b12:
	;
	v53 = v49
	v_54 = v_50
	or__x_55 = or__x_23
	goto b9
b13:
	;
	v45 = or__x_37
	v_46 = v_36
	or__x_47 = or__x_37
	goto b15
b14:
	;
	v43, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "boolean?").Deref(), []vm.Value{v_38})
	if callErr != nil {
		return nil, callErr
	}
	v45 = v43
	v_46 = v_38
	or__x_47 = or__x_39
	goto b15
b15:
	;
	v49 = v45
	v_50 = v_46
	or__x_51 = or__x_31
	goto b12
}
func lookup_local(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__7502_6 vm.Value
	var arg__7508_10 vm.Value
	var arg__7510_12 vm.Value
	var arg__7515_15 vm.Value
	var arg__7521_19 vm.Value
	var arg__7523_21 vm.Value
	var arg__7524_22 vm.Value
	var v23 vm.Value
	var i_2 vm.Value
	var sym_3 vm.Value
	var ctx_4 vm.Value
	var v158 int
	var v165 vm.Value
	var v172 vm.Value
	var v32 bool
	var i_25 vm.Value
	var sym_26 vm.Value
	var ctx_27 vm.Value
	var v161 int
	var v168 vm.Value
	var v175 vm.Value
	var i_28 vm.Value
	var sym_29 vm.Value
	var ctx_30 vm.Value
	var v159 int
	var v166 vm.Value
	var v173 vm.Value
	var arg__7530_43 vm.Value
	var arg__7536_47 vm.Value
	var arg__7538_49 vm.Value
	var arg__7544_52 vm.Value
	var arg__7550_56 vm.Value
	var arg__7552_58 vm.Value
	var arg__7554_59 vm.Value
	var arg__7560_62 vm.Value
	var arg__7566_66 vm.Value
	var arg__7568_68 vm.Value
	var arg__7574_71 vm.Value
	var arg__7580_75 vm.Value
	var arg__7582_77 vm.Value
	var arg__7584_78 vm.Value
	var v79 vm.Value
	var v142 vm.Value
	var i_143 vm.Value
	var sym_144 vm.Value
	var ctx_145 vm.Value
	var i_36 vm.Value
	var sym_37 vm.Value
	var ctx_38 vm.Value
	var v162 int
	var v169 vm.Value
	var v176 vm.Value
	var arg__7589_82 vm.Value
	var arg__7595_86 vm.Value
	var arg__7597_88 vm.Value
	var arg__7603_91 vm.Value
	var arg__7609_95 vm.Value
	var arg__7611_97 vm.Value
	var arg__7613_98 vm.Value
	var arg__7619_101 vm.Value
	var arg__7625_105 vm.Value
	var arg__7627_107 vm.Value
	var arg__7633_110 vm.Value
	var arg__7639_114 vm.Value
	var arg__7641_116 vm.Value
	var arg__7643_117 vm.Value
	var v118 vm.Value
	var i_39 vm.Value
	var sym_40 vm.Value
	var ctx_41 vm.Value
	var v160 int
	var v167 vm.Value
	var v174 vm.Value
	var v137 vm.Value
	var i_138 vm.Value
	var sym_139 vm.Value
	var ctx_140 vm.Value
	var i_120 vm.Value
	var sym_121 vm.Value
	var ctx_122 vm.Value
	var v157 int
	var v164 vm.Value
	var v171 vm.Value
	var v128 vm.Value
	var i_123 vm.Value
	var sym_124 vm.Value
	var ctx_125 vm.Value
	var v163 int
	var v170 vm.Value
	var v177 vm.Value
	var v132 vm.Value
	var i_133 vm.Value
	var sym_134 vm.Value
	var ctx_135 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__7502_6, arg__7508_10, arg__7510_12, arg__7515_15, arg__7521_19, arg__7523_21, arg__7524_22, v23, i_2, sym_3, ctx_4, v158, v165, v172, v32, i_25, sym_26, ctx_27, v161, v168, v175, i_28, sym_29, ctx_30, v159, v166, v173, arg__7530_43, arg__7536_47, arg__7538_49, arg__7544_52, arg__7550_56, arg__7552_58, arg__7554_59, arg__7560_62, arg__7566_66, arg__7568_68, arg__7574_71, arg__7580_75, arg__7582_77, arg__7584_78, v79, v142, i_143, sym_144, ctx_145, i_36, sym_37, ctx_38, v162, v169, v176, arg__7589_82, arg__7595_86, arg__7597_88, arg__7603_91, arg__7609_95, arg__7611_97, arg__7613_98, arg__7619_101, arg__7625_105, arg__7627_107, arg__7633_110, arg__7639_114, arg__7641_116, arg__7643_117, v118, i_39, sym_40, ctx_41, v160, v167, v174, v137, i_138, sym_139, ctx_140, i_120, sym_121, ctx_122, v157, v164, v171, v128, i_123, sym_124, ctx_125, v163, v170, v177, v132, i_133, sym_134, ctx_135
	arg__7502_6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__7508_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__7510_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg__7508_10, vm.Keyword("locals")})
	if callErr != nil {
		return nil, callErr
	}
	arg__7515_15, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__7521_19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__7523_21, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg__7521_19, vm.Keyword("locals")})
	if callErr != nil {
		return nil, callErr
	}
	arg__7524_22, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{arg__7523_21})
	if callErr != nil {
		return nil, callErr
	}
	v23 = rt.SubValue(arg__7524_22, vm.Int(1))
	i_2 = v23
	sym_3 = arg1
	ctx_4 = arg0
	v158 = 0
	v165 = vm.Keyword("locals")
	v172 = vm.Keyword("else")
	goto b1
b1:
	;
	v32 = rt.LtValue(i_2, vm.Int(v158))
	if v32 {
		i_25 = i_2
		sym_26 = sym_3
		ctx_27 = ctx_4
		v161 = v158
		v168 = v165
		v175 = v172
		goto b2
	} else {
		i_28 = i_2
		sym_29 = sym_3
		ctx_30 = ctx_4
		v159 = v158
		v166 = v165
		v173 = v172
		goto b3
	}
b2:
	;
	v142 = vm.NIL
	i_143 = i_25
	sym_144 = sym_26
	ctx_145 = ctx_27
	goto b4
b3:
	;
	arg__7530_43, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_30})
	if callErr != nil {
		return nil, callErr
	}
	arg__7536_47, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_30})
	if callErr != nil {
		return nil, callErr
	}
	arg__7538_49, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg__7536_47, v166})
	if callErr != nil {
		return nil, callErr
	}
	arg__7544_52, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_30})
	if callErr != nil {
		return nil, callErr
	}
	arg__7550_56, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_30})
	if callErr != nil {
		return nil, callErr
	}
	arg__7552_58, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg__7550_56, v166})
	if callErr != nil {
		return nil, callErr
	}
	arg__7554_59, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__7552_58, i_28})
	if callErr != nil {
		return nil, callErr
	}
	arg__7560_62, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_30})
	if callErr != nil {
		return nil, callErr
	}
	arg__7566_66, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_30})
	if callErr != nil {
		return nil, callErr
	}
	arg__7568_68, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg__7566_66, v166})
	if callErr != nil {
		return nil, callErr
	}
	arg__7574_71, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_30})
	if callErr != nil {
		return nil, callErr
	}
	arg__7580_75, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_30})
	if callErr != nil {
		return nil, callErr
	}
	arg__7582_77, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg__7580_75, v166})
	if callErr != nil {
		return nil, callErr
	}
	arg__7584_78, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__7582_77, i_28})
	if callErr != nil {
		return nil, callErr
	}
	v79, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "contains?").Deref(), []vm.Value{arg__7584_78, sym_29})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v79) {
		i_36 = i_28
		sym_37 = sym_29
		ctx_38 = ctx_30
		v162 = v159
		v169 = v166
		v176 = v173
		goto b5
	} else {
		i_39 = i_28
		sym_40 = sym_29
		ctx_41 = ctx_30
		v160 = v159
		v167 = v166
		v174 = v173
		goto b6
	}
b4:
	;
	return v142, nil
b5:
	;
	arg__7589_82, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_38})
	if callErr != nil {
		return nil, callErr
	}
	arg__7595_86, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_38})
	if callErr != nil {
		return nil, callErr
	}
	arg__7597_88, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg__7595_86, vm.Keyword("locals")})
	if callErr != nil {
		return nil, callErr
	}
	arg__7603_91, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_38})
	if callErr != nil {
		return nil, callErr
	}
	arg__7609_95, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_38})
	if callErr != nil {
		return nil, callErr
	}
	arg__7611_97, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg__7609_95, vm.Keyword("locals")})
	if callErr != nil {
		return nil, callErr
	}
	arg__7613_98, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__7611_97, i_36})
	if callErr != nil {
		return nil, callErr
	}
	arg__7619_101, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_38})
	if callErr != nil {
		return nil, callErr
	}
	arg__7625_105, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_38})
	if callErr != nil {
		return nil, callErr
	}
	arg__7627_107, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg__7625_105, vm.Keyword("locals")})
	if callErr != nil {
		return nil, callErr
	}
	arg__7633_110, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_38})
	if callErr != nil {
		return nil, callErr
	}
	arg__7639_114, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ctx_38})
	if callErr != nil {
		return nil, callErr
	}
	arg__7641_116, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg__7639_114, vm.Keyword("locals")})
	if callErr != nil {
		return nil, callErr
	}
	arg__7643_117, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__7641_116, i_36})
	if callErr != nil {
		return nil, callErr
	}
	v118, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg__7643_117, sym_37})
	if callErr != nil {
		return nil, callErr
	}
	v137 = v118
	i_138 = i_36
	sym_139 = sym_37
	ctx_140 = ctx_38
	goto b7
b6:
	;
	if vm.IsTruthy(v174) {
		i_120 = i_39
		sym_121 = sym_40
		ctx_122 = ctx_41
		v157 = v160
		v164 = v167
		v171 = v174
		goto b8
	} else {
		i_123 = i_39
		sym_124 = sym_40
		ctx_125 = ctx_41
		v163 = v160
		v170 = v167
		v177 = v174
		goto b9
	}
b7:
	;
	v142 = v137
	i_143 = i_138
	sym_144 = sym_139
	ctx_145 = ctx_140
	goto b4
b8:
	;
	v128 = rt.SubValue(i_120, vm.Int(1))
	i_2 = v128
	sym_3 = sym_121
	ctx_4 = ctx_122
	v158 = v157
	v165 = v164
	v172 = v171
	goto b1
b9:
	;
	v132 = vm.NIL
	i_133 = i_123
	sym_134 = sym_124
	ctx_135 = ctx_125
	goto b10
b10:
	;
	v137 = v132
	i_138 = i_133
	sym_139 = sym_134
	ctx_140 = ctx_135
	goto b7
}
func new_context(arg0 vm.Value) (vm.Value, error) {
	var arg__7653_5 vm.Value
	var arg__7657_9 vm.Value
	var arg__7658_10 vm.Value
	var arg__7667_16 vm.Value
	var arg__7671_20 vm.Value
	var arg__7672_21 vm.Value
	var v22 vm.Value
	var callErr error
	_, _, _, _, _, _, _ = arg__7653_5, arg__7657_9, arg__7658_10, arg__7667_16, arg__7671_20, arg__7672_21, v22
	arg__7653_5, callErr = rt.InvokeValue(rt.LookupVar("ir", "entry-block").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__7657_9, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.EmptyPersistentMap})
	if callErr != nil {
		return nil, callErr
	}
	arg__7658_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("fn"), arg0, vm.Keyword("current-block"), arg__7653_5, vm.Keyword("locals"), arg__7657_9})
	if callErr != nil {
		return nil, callErr
	}
	arg__7667_16, callErr = rt.InvokeValue(rt.LookupVar("ir", "entry-block").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__7671_20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.EmptyPersistentMap})
	if callErr != nil {
		return nil, callErr
	}
	arg__7672_21, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("fn"), arg0, vm.Keyword("current-block"), arg__7667_16, vm.Keyword("locals"), arg__7671_20})
	if callErr != nil {
		return nil, callErr
	}
	v22, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "atom").Deref(), []vm.Value{arg__7672_21})
	if callErr != nil {
		return nil, callErr
	}
	return v22, nil
}
func parse_defn_args_body(arg0 vm.Value) (vm.Value, error) {
	var x2_6 vm.Value
	var x3_12 vm.Value
	var x4_18 vm.Value
	var v28 vm.Value
	var form_19 vm.Value
	var x2_20 vm.Value
	var x3_21 vm.Value
	var x4_22 vm.Value
	var v39 vm.Value
	var form_23 vm.Value
	var x2_24 vm.Value
	var x3_25 vm.Value
	var x4_26 vm.Value
	var v70 vm.Value
	var v110 vm.Value
	var form_111 vm.Value
	var x2_112 vm.Value
	var x3_113 vm.Value
	var x4_114 vm.Value
	var form_30 vm.Value
	var x2_31 vm.Value
	var x3_32 vm.Value
	var x4_33 vm.Value
	var arg__7707_45 vm.Value
	var v46 vm.Value
	var form_34 vm.Value
	var x2_35 vm.Value
	var x3_36 vm.Value
	var x4_37 vm.Value
	var arg__7715_52 vm.Value
	var v53 vm.Value
	var v55 vm.Value
	var form_56 vm.Value
	var x2_57 vm.Value
	var x3_58 vm.Value
	var x4_59 vm.Value
	var form_61 vm.Value
	var x2_62 vm.Value
	var x3_63 vm.Value
	var x4_64 vm.Value
	var arg__7726_76 vm.Value
	var v77 vm.Value
	var form_65 vm.Value
	var x2_66 vm.Value
	var x3_67 vm.Value
	var x4_68 vm.Value
	var v104 vm.Value
	var form_105 vm.Value
	var x2_106 vm.Value
	var x3_107 vm.Value
	var x4_108 vm.Value
	var form_79 vm.Value
	var x2_80 vm.Value
	var x3_81 vm.Value
	var x4_82 vm.Value
	var arg__7734_93 vm.Value
	var v94 vm.Value
	var form_83 vm.Value
	var x2_84 vm.Value
	var x3_85 vm.Value
	var x4_86 vm.Value
	var v98 vm.Value
	var form_99 vm.Value
	var x2_100 vm.Value
	var x3_101 vm.Value
	var x4_102 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = x2_6, x3_12, x4_18, v28, form_19, x2_20, x3_21, x4_22, v39, form_23, x2_24, x3_25, x4_26, v70, v110, form_111, x2_112, x3_113, x4_114, form_30, x2_31, x3_32, x4_33, arg__7707_45, v46, form_34, x2_35, x3_36, x4_37, arg__7715_52, v53, v55, form_56, x2_57, x3_58, x4_59, form_61, x2_62, x3_63, x4_64, arg__7726_76, v77, form_65, x2_66, x3_67, x4_68, v104, form_105, x2_106, x3_107, x4_108, form_79, x2_80, x3_81, x4_82, arg__7734_93, v94, form_83, x2_84, x3_85, x4_86, v98, form_99, x2_100, x3_101, x4_102
	x2_6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(2), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	x3_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(3), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	x4_18, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(4), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	v28, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "string?").Deref(), []vm.Value{x2_6})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v28) {
		form_19 = arg0
		x2_20 = x2_6
		x3_21 = x3_12
		x4_22 = x4_18
		goto b1
	} else {
		form_23 = arg0
		x2_24 = x2_6
		x3_25 = x3_12
		x4_26 = x4_18
		goto b2
	}
b1:
	;
	v39, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map?").Deref(), []vm.Value{x3_21})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v39) {
		form_30 = form_19
		x2_31 = x2_20
		x3_32 = x3_21
		x4_33 = x4_22
		goto b4
	} else {
		form_34 = form_19
		x2_35 = x2_20
		x3_36 = x3_21
		x4_37 = x4_22
		goto b5
	}
b2:
	;
	v70, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map?").Deref(), []vm.Value{x2_24})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v70) {
		form_61 = form_23
		x2_62 = x2_24
		x3_63 = x3_25
		x4_64 = x4_26
		goto b7
	} else {
		form_65 = form_23
		x2_66 = x2_24
		x3_67 = x3_25
		x4_68 = x4_26
		goto b8
	}
b3:
	;
	return v110, nil
b4:
	;
	arg__7707_45, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "drop").Deref(), []vm.Value{vm.Int(5), form_30})
	if callErr != nil {
		return nil, callErr
	}
	v46, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{x4_33, arg__7707_45})
	if callErr != nil {
		return nil, callErr
	}
	v55 = v46
	form_56 = form_30
	x2_57 = x2_31
	x3_58 = x3_32
	x4_59 = x4_33
	goto b6
b5:
	;
	arg__7715_52, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "drop").Deref(), []vm.Value{vm.Int(4), form_34})
	if callErr != nil {
		return nil, callErr
	}
	v53, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{x3_36, arg__7715_52})
	if callErr != nil {
		return nil, callErr
	}
	v55 = v53
	form_56 = form_34
	x2_57 = x2_35
	x3_58 = x3_36
	x4_59 = x4_37
	goto b6
b6:
	;
	v110 = v55
	form_111 = form_56
	x2_112 = x2_57
	x3_113 = x3_58
	x4_114 = x4_59
	goto b3
b7:
	;
	arg__7726_76, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "drop").Deref(), []vm.Value{vm.Int(4), form_61})
	if callErr != nil {
		return nil, callErr
	}
	v77, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{x3_63, arg__7726_76})
	if callErr != nil {
		return nil, callErr
	}
	v104 = v77
	form_105 = form_61
	x2_106 = x2_62
	x3_107 = x3_63
	x4_108 = x4_64
	goto b9
b8:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		form_79 = form_65
		x2_80 = x2_66
		x3_81 = x3_67
		x4_82 = x4_68
		goto b10
	} else {
		form_83 = form_65
		x2_84 = x2_66
		x3_85 = x3_67
		x4_86 = x4_68
		goto b11
	}
b9:
	;
	v110 = v104
	form_111 = form_105
	x2_112 = x2_106
	x3_113 = x3_107
	x4_114 = x4_108
	goto b3
b10:
	;
	arg__7734_93, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "drop").Deref(), []vm.Value{vm.Int(3), form_79})
	if callErr != nil {
		return nil, callErr
	}
	v94, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{x2_80, arg__7734_93})
	if callErr != nil {
		return nil, callErr
	}
	v98 = v94
	form_99 = form_79
	x2_100 = x2_80
	x3_101 = x3_81
	x4_102 = x4_82
	goto b12
b11:
	;
	v98 = vm.NIL
	form_99 = form_83
	x2_100 = x2_84
	x3_101 = x3_85
	x4_102 = x4_86
	goto b12
b12:
	;
	v104 = v98
	form_105 = form_99
	x2_106 = x2_100
	x3_107 = x3_101
	x4_108 = x4_102
	goto b9
}
func pop_locals_BANG_(arg0 vm.Value) (vm.Value, error) {
	var s_2 vm.Value
	var stack_6 vm.Value
	var arg__7751_10 vm.Value
	var arg__7759_15 vm.Value
	var arg__7760_16 vm.Value
	var arg__7761_17 vm.Value
	var arg__7770_22 vm.Value
	var arg__7778_27 vm.Value
	var arg__7779_28 vm.Value
	var arg__7780_29 vm.Value
	var arg__7781_30 vm.Value
	var arg__7791_35 vm.Value
	var arg__7799_40 vm.Value
	var arg__7800_41 vm.Value
	var arg__7801_42 vm.Value
	var arg__7810_47 vm.Value
	var arg__7818_52 vm.Value
	var arg__7819_53 vm.Value
	var arg__7820_54 vm.Value
	var arg__7821_55 vm.Value
	var v56 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = s_2, stack_6, arg__7751_10, arg__7759_15, arg__7760_16, arg__7761_17, arg__7770_22, arg__7778_27, arg__7779_28, arg__7780_29, arg__7781_30, arg__7791_35, arg__7799_40, arg__7800_41, arg__7801_42, arg__7810_47, arg__7818_52, arg__7819_53, arg__7820_54, arg__7821_55, v56
	s_2, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	stack_6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{s_2, vm.Keyword("locals")})
	if callErr != nil {
		return nil, callErr
	}
	arg__7751_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{stack_6})
	if callErr != nil {
		return nil, callErr
	}
	arg__7759_15, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{stack_6})
	if callErr != nil {
		return nil, callErr
	}
	arg__7760_16 = rt.SubValue(arg__7759_15, vm.Int(1))
	arg__7761_17, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "subvec").Deref(), []vm.Value{stack_6, vm.Int(0), arg__7760_16})
	if callErr != nil {
		return nil, callErr
	}
	arg__7770_22, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{stack_6})
	if callErr != nil {
		return nil, callErr
	}
	arg__7778_27, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{stack_6})
	if callErr != nil {
		return nil, callErr
	}
	arg__7779_28 = rt.SubValue(arg__7778_27, vm.Int(1))
	arg__7780_29, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "subvec").Deref(), []vm.Value{stack_6, vm.Int(0), arg__7779_28})
	if callErr != nil {
		return nil, callErr
	}
	arg__7781_30, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "assoc").Deref(), []vm.Value{s_2, vm.Keyword("locals"), arg__7780_29})
	if callErr != nil {
		return nil, callErr
	}
	arg__7791_35, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{stack_6})
	if callErr != nil {
		return nil, callErr
	}
	arg__7799_40, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{stack_6})
	if callErr != nil {
		return nil, callErr
	}
	arg__7800_41 = rt.SubValue(arg__7799_40, vm.Int(1))
	arg__7801_42, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "subvec").Deref(), []vm.Value{stack_6, vm.Int(0), arg__7800_41})
	if callErr != nil {
		return nil, callErr
	}
	arg__7810_47, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{stack_6})
	if callErr != nil {
		return nil, callErr
	}
	arg__7818_52, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{stack_6})
	if callErr != nil {
		return nil, callErr
	}
	arg__7819_53 = rt.SubValue(arg__7818_52, vm.Int(1))
	arg__7820_54, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "subvec").Deref(), []vm.Value{stack_6, vm.Int(0), arg__7819_53})
	if callErr != nil {
		return nil, callErr
	}
	arg__7821_55, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "assoc").Deref(), []vm.Value{s_2, vm.Keyword("locals"), arg__7820_54})
	if callErr != nil {
		return nil, callErr
	}
	v56, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "reset!").Deref(), []vm.Value{arg0, arg__7821_55})
	if callErr != nil {
		return nil, callErr
	}
	return v56, nil
}
func push_locals_BANG_(arg0 vm.Value) (vm.Value, error) {
	var s_2 vm.Value
	var arg__7833_7 vm.Value
	var arg__7841_13 vm.Value
	var arg__7843_15 vm.Value
	var arg__7852_21 vm.Value
	var arg__7860_27 vm.Value
	var arg__7862_29 vm.Value
	var arg__7863_30 vm.Value
	var arg__7873_36 vm.Value
	var arg__7881_42 vm.Value
	var arg__7883_44 vm.Value
	var arg__7892_50 vm.Value
	var arg__7900_56 vm.Value
	var arg__7902_58 vm.Value
	var arg__7903_59 vm.Value
	var v60 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = s_2, arg__7833_7, arg__7841_13, arg__7843_15, arg__7852_21, arg__7860_27, arg__7862_29, arg__7863_30, arg__7873_36, arg__7881_42, arg__7883_44, arg__7892_50, arg__7900_56, arg__7902_58, arg__7903_59, v60
	s_2, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__7833_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{s_2, vm.Keyword("locals")})
	if callErr != nil {
		return nil, callErr
	}
	arg__7841_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{s_2, vm.Keyword("locals")})
	if callErr != nil {
		return nil, callErr
	}
	arg__7843_15, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{arg__7841_13, vm.EmptyPersistentMap})
	if callErr != nil {
		return nil, callErr
	}
	arg__7852_21, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{s_2, vm.Keyword("locals")})
	if callErr != nil {
		return nil, callErr
	}
	arg__7860_27, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{s_2, vm.Keyword("locals")})
	if callErr != nil {
		return nil, callErr
	}
	arg__7862_29, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{arg__7860_27, vm.EmptyPersistentMap})
	if callErr != nil {
		return nil, callErr
	}
	arg__7863_30, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "assoc").Deref(), []vm.Value{s_2, vm.Keyword("locals"), arg__7862_29})
	if callErr != nil {
		return nil, callErr
	}
	arg__7873_36, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{s_2, vm.Keyword("locals")})
	if callErr != nil {
		return nil, callErr
	}
	arg__7881_42, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{s_2, vm.Keyword("locals")})
	if callErr != nil {
		return nil, callErr
	}
	arg__7883_44, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{arg__7881_42, vm.EmptyPersistentMap})
	if callErr != nil {
		return nil, callErr
	}
	arg__7892_50, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{s_2, vm.Keyword("locals")})
	if callErr != nil {
		return nil, callErr
	}
	arg__7900_56, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{s_2, vm.Keyword("locals")})
	if callErr != nil {
		return nil, callErr
	}
	arg__7902_58, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{arg__7900_56, vm.EmptyPersistentMap})
	if callErr != nil {
		return nil, callErr
	}
	arg__7903_59, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "assoc").Deref(), []vm.Value{s_2, vm.Keyword("locals"), arg__7902_58})
	if callErr != nil {
		return nil, callErr
	}
	v60, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "reset!").Deref(), []vm.Value{arg0, arg__7903_59})
	if callErr != nil {
		return nil, callErr
	}
	return v60, nil
}
func rebind_local_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var s_5 vm.Value
	var stack_9 vm.Value
	var arg__7915_17 vm.Value
	var v18 vm.Value
	var i_10 vm.Value
	var sym_11 vm.Value
	var ctx_12 vm.Value
	var inst_id_13 vm.Value
	var s_14 vm.Value
	var stack_15 vm.Value
	var v123 int
	var v130 vm.Value
	var v33 bool
	var i_20 vm.Value
	var sym_21 vm.Value
	var ctx_22 vm.Value
	var inst_id_23 vm.Value
	var s_24 vm.Value
	var stack_25 vm.Value
	var v126 int
	var v133 vm.Value
	var v36 vm.Value
	var i_26 vm.Value
	var sym_27 vm.Value
	var ctx_28 vm.Value
	var inst_id_29 vm.Value
	var s_30 vm.Value
	var stack_31 vm.Value
	var v124 int
	var v131 vm.Value
	var arg__7930_51 vm.Value
	var arg__7938_54 vm.Value
	var v55 vm.Value
	var v112 vm.Value
	var i_113 vm.Value
	var sym_114 vm.Value
	var ctx_115 vm.Value
	var inst_id_116 vm.Value
	var s_117 vm.Value
	var stack_118 vm.Value
	var i_38 vm.Value
	var sym_39 vm.Value
	var ctx_40 vm.Value
	var inst_id_41 vm.Value
	var s_42 vm.Value
	var stack_43 vm.Value
	var v127 int
	var v134 vm.Value
	var arg__7945_58 vm.Value
	var arg__7954_61 vm.Value
	var updated_frame_62 vm.Value
	var new_stack_64 vm.Value
	var arg__7972_68 vm.Value
	var arg__7982_73 vm.Value
	var v74 vm.Value
	var v76 vm.Value
	var i_44 vm.Value
	var sym_45 vm.Value
	var ctx_46 vm.Value
	var inst_id_47 vm.Value
	var s_48 vm.Value
	var stack_49 vm.Value
	var v125 int
	var v132 vm.Value
	var v104 vm.Value
	var i_105 vm.Value
	var sym_106 vm.Value
	var ctx_107 vm.Value
	var inst_id_108 vm.Value
	var s_109 vm.Value
	var stack_110 vm.Value
	var i_78 vm.Value
	var sym_79 vm.Value
	var ctx_80 vm.Value
	var inst_id_81 vm.Value
	var s_82 vm.Value
	var stack_83 vm.Value
	var v122 int
	var v129 vm.Value
	var v92 vm.Value
	var i_84 vm.Value
	var sym_85 vm.Value
	var ctx_86 vm.Value
	var inst_id_87 vm.Value
	var s_88 vm.Value
	var stack_89 vm.Value
	var v128 int
	var v135 vm.Value
	var v96 vm.Value
	var i_97 vm.Value
	var sym_98 vm.Value
	var ctx_99 vm.Value
	var inst_id_100 vm.Value
	var s_101 vm.Value
	var stack_102 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = s_5, stack_9, arg__7915_17, v18, i_10, sym_11, ctx_12, inst_id_13, s_14, stack_15, v123, v130, v33, i_20, sym_21, ctx_22, inst_id_23, s_24, stack_25, v126, v133, v36, i_26, sym_27, ctx_28, inst_id_29, s_30, stack_31, v124, v131, arg__7930_51, arg__7938_54, v55, v112, i_113, sym_114, ctx_115, inst_id_116, s_117, stack_118, i_38, sym_39, ctx_40, inst_id_41, s_42, stack_43, v127, v134, arg__7945_58, arg__7954_61, updated_frame_62, new_stack_64, arg__7972_68, arg__7982_73, v74, v76, i_44, sym_45, ctx_46, inst_id_47, s_48, stack_49, v125, v132, v104, i_105, sym_106, ctx_107, inst_id_108, s_109, stack_110, i_78, sym_79, ctx_80, inst_id_81, s_82, stack_83, v122, v129, v92, i_84, sym_85, ctx_86, inst_id_87, s_88, stack_89, v128, v135, v96, i_97, sym_98, ctx_99, inst_id_100, s_101, stack_102
	s_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	stack_9, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{s_5, vm.Keyword("locals")})
	if callErr != nil {
		return nil, callErr
	}
	arg__7915_17, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{stack_9})
	if callErr != nil {
		return nil, callErr
	}
	v18 = rt.SubValue(arg__7915_17, vm.Int(1))
	i_10 = v18
	sym_11 = arg1
	ctx_12 = arg0
	inst_id_13 = arg2
	s_14 = s_5
	stack_15 = stack_9
	v123 = 0
	v130 = vm.Keyword("else")
	goto b1
b1:
	;
	v33 = rt.LtValue(i_10, vm.Int(v123))
	if v33 {
		i_20 = i_10
		sym_21 = sym_11
		ctx_22 = ctx_12
		inst_id_23 = inst_id_13
		s_24 = s_14
		stack_25 = stack_15
		v126 = v123
		v133 = v130
		goto b2
	} else {
		i_26 = i_10
		sym_27 = sym_11
		ctx_28 = ctx_12
		inst_id_29 = inst_id_13
		s_30 = s_14
		stack_31 = stack_15
		v124 = v123
		v131 = v130
		goto b3
	}
b2:
	;
	v36, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "bind-local!").Deref(), []vm.Value{ctx_22, sym_21, inst_id_23})
	if callErr != nil {
		return nil, callErr
	}
	v112 = v36
	i_113 = i_20
	sym_114 = sym_21
	ctx_115 = ctx_22
	inst_id_116 = inst_id_23
	s_117 = s_24
	stack_118 = stack_25
	goto b4
b3:
	;
	arg__7930_51, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{stack_31, i_26})
	if callErr != nil {
		return nil, callErr
	}
	arg__7938_54, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{stack_31, i_26})
	if callErr != nil {
		return nil, callErr
	}
	v55, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "contains?").Deref(), []vm.Value{arg__7938_54, sym_27})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v55) {
		i_38 = i_26
		sym_39 = sym_27
		ctx_40 = ctx_28
		inst_id_41 = inst_id_29
		s_42 = s_30
		stack_43 = stack_31
		v127 = v124
		v134 = v131
		goto b5
	} else {
		i_44 = i_26
		sym_45 = sym_27
		ctx_46 = ctx_28
		inst_id_47 = inst_id_29
		s_48 = s_30
		stack_49 = stack_31
		v125 = v124
		v132 = v131
		goto b6
	}
b4:
	;
	return v112, nil
b5:
	;
	arg__7945_58, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{stack_43, i_38})
	if callErr != nil {
		return nil, callErr
	}
	arg__7954_61, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{stack_43, i_38})
	if callErr != nil {
		return nil, callErr
	}
	updated_frame_62, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "assoc").Deref(), []vm.Value{arg__7954_61, sym_39, inst_id_41})
	if callErr != nil {
		return nil, callErr
	}
	new_stack_64, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "assoc").Deref(), []vm.Value{stack_43, i_38, updated_frame_62})
	if callErr != nil {
		return nil, callErr
	}
	arg__7972_68, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "assoc").Deref(), []vm.Value{s_42, vm.Keyword("locals"), new_stack_64})
	if callErr != nil {
		return nil, callErr
	}
	arg__7982_73, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "assoc").Deref(), []vm.Value{s_42, vm.Keyword("locals"), new_stack_64})
	if callErr != nil {
		return nil, callErr
	}
	v74, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "reset!").Deref(), []vm.Value{ctx_40, arg__7982_73})
	if callErr != nil {
		return nil, callErr
	}
	v76, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "attach-name!").Deref(), []vm.Value{ctx_40, inst_id_41, sym_39})
	if callErr != nil {
		return nil, callErr
	}
	v104 = v76
	i_105 = i_38
	sym_106 = sym_39
	ctx_107 = ctx_40
	inst_id_108 = inst_id_41
	s_109 = s_42
	stack_110 = stack_43
	goto b7
b6:
	;
	if vm.IsTruthy(v132) {
		i_78 = i_44
		sym_79 = sym_45
		ctx_80 = ctx_46
		inst_id_81 = inst_id_47
		s_82 = s_48
		stack_83 = stack_49
		v122 = v125
		v129 = v132
		goto b8
	} else {
		i_84 = i_44
		sym_85 = sym_45
		ctx_86 = ctx_46
		inst_id_87 = inst_id_47
		s_88 = s_48
		stack_89 = stack_49
		v128 = v125
		v135 = v132
		goto b9
	}
b7:
	;
	v112 = v104
	i_113 = i_105
	sym_114 = sym_106
	ctx_115 = ctx_107
	inst_id_116 = inst_id_108
	s_117 = s_109
	stack_118 = stack_110
	goto b4
b8:
	;
	v92 = rt.SubValue(i_78, vm.Int(1))
	i_10 = v92
	sym_11 = sym_79
	ctx_12 = ctx_80
	inst_id_13 = inst_id_81
	s_14 = s_82
	stack_15 = stack_83
	v123 = v122
	v130 = v129
	goto b1
b9:
	;
	v96 = vm.NIL
	i_97 = i_84
	sym_98 = sym_85
	ctx_99 = ctx_86
	inst_id_100 = inst_id_87
	s_101 = s_88
	stack_102 = stack_89
	goto b10
b10:
	;
	v104 = v96
	i_105 = i_97
	sym_106 = sym_98
	ctx_107 = ctx_99
	inst_id_108 = inst_id_100
	s_109 = s_101
	stack_110 = stack_102
	goto b7
}
func terminated_QMARK_(arg0 vm.Value) bool {
	var v2 bool
	_ = v2
	v2 = arg0 == rt.LookupVar("ir.build", "TERMINATED").Deref()
	return v2
}
