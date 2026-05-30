package ir_lower_go

import (
	rt "github.com/nooga/let-go/pkg/rt"
	vm "github.com/nooga/let-go/pkg/vm"
)

func any_fn_template_QMARK_(arg0 vm.Value) (vm.Value, error) {
	var or__x_2 vm.Value
	var aux_3 vm.Value
	var or__x_4 vm.Value
	var aux_5 vm.Value
	var or__x_6 vm.Value
	var v10 vm.Value
	var v12 vm.Value
	var aux_13 vm.Value
	var or__x_14 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _ = or__x_2, aux_3, or__x_4, aux_5, or__x_6, v10, v12, aux_13, or__x_14
	or__x_2, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "fn-template?").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(or__x_2) {
		aux_3 = arg0
		or__x_4 = or__x_2
		goto b1
	} else {
		aux_5 = arg0
		or__x_6 = or__x_2
		goto b2
	}
b1:
	;
	v12 = or__x_4
	aux_13 = aux_3
	or__x_14 = or__x_4
	goto b3
b2:
	;
	v10, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "multi-fn-template?").Deref(), []vm.Value{aux_5})
	if callErr != nil {
		return nil, callErr
	}
	v12 = v10
	aux_13 = aux_5
	or__x_14 = or__x_6
	goto b3
b3:
	;
	return v12, nil
}
func arg_local_decls(arg0 vm.Value) (vm.Value, error) {
	var arity_3 vm.Value
	var variadic_QMARK__5 vm.Value
	var f_6 vm.Value
	var arity_7 vm.Value
	var variadic_QMARK__8 vm.Value
	var v13 vm.Value
	var f_9 vm.Value
	var arity_10 vm.Value
	var variadic_QMARK__11 vm.Value
	var fixed_arity_16 vm.Value
	var f_17 vm.Value
	var arity_18 vm.Value
	var variadic_QMARK__19 vm.Value
	var i_20 int
	var decls_21 vm.Value
	var assigns_22 vm.Value
	var variadic_QMARK__23 vm.Value
	var fixed_arity_24 vm.Value
	var v245 string
	var v248 string
	var v251 string
	var v254 vm.Value
	var v257 string
	var v43 bool
	var f_29 vm.Value
	var arity_30 vm.Value
	var i_31 int
	var decls_32 vm.Value
	var assigns_33 vm.Value
	var variadic_QMARK__34 vm.Value
	var fixed_arity_35 vm.Value
	var v246 string
	var v249 string
	var v252 string
	var v255 vm.Value
	var v258 string
	var f_36 vm.Value
	var arity_37 vm.Value
	var i_38 int
	var decls_39 vm.Value
	var assigns_40 vm.Value
	var variadic_QMARK__41 vm.Value
	var fixed_arity_42 vm.Value
	var v244 string
	var v247 string
	var v250 string
	var v253 vm.Value
	var v256 string
	var sym_156 vm.Value
	var arg_sym_160 vm.Value
	var v161 int
	var arg__12790_165 vm.Value
	var arg__12797_171 vm.Value
	var arg__12799_173 vm.Value
	var arg__12806_178 vm.Value
	var arg__12813_184 vm.Value
	var arg__12815_186 vm.Value
	var v187 vm.Value
	var arg__12821_190 vm.Value
	var arg__12825_192 vm.Value
	var arg__12831_196 vm.Value
	var arg__12835_198 vm.Value
	var arg__12836_199 vm.Value
	var arg__12843_203 vm.Value
	var arg__12847_205 vm.Value
	var arg__12853_209 vm.Value
	var arg__12857_211 vm.Value
	var arg__12858_212 vm.Value
	var v213 vm.Value
	var v215 vm.Value
	var f_216 vm.Value
	var arity_217 vm.Value
	var i_218 int
	var decls_219 vm.Value
	var assigns_220 vm.Value
	var variadic_QMARK__221 vm.Value
	var fixed_arity_222 vm.Value
	var f_45 vm.Value
	var arity_46 vm.Value
	var i_47 int
	var decls_48 vm.Value
	var final_decls_49 vm.Value
	var final_assigns_50 vm.Value
	var assigns_51 vm.Value
	var variadic_QMARK__52 vm.Value
	var fixed_arity_53 vm.Value
	var arg__12700_69 vm.Value
	var arg__12707_76 vm.Value
	var arg__12709_78 vm.Value
	var arg__12716_84 vm.Value
	var arg__12723_91 vm.Value
	var arg__12725_93 vm.Value
	var arg__12726_94 vm.Value
	var arg__12732_99 vm.Value
	var arg__12736_103 vm.Value
	var arg__12742_109 vm.Value
	var arg__12746_113 vm.Value
	var arg__12747_114 vm.Value
	var arg__12754_120 vm.Value
	var arg__12758_124 vm.Value
	var arg__12764_130 vm.Value
	var arg__12768_134 vm.Value
	var arg__12769_135 vm.Value
	var arg__12770_136 vm.Value
	var v137 vm.Value
	var f_54 vm.Value
	var arity_55 vm.Value
	var i_56 int
	var decls_57 vm.Value
	var final_decls_58 vm.Value
	var final_assigns_59 vm.Value
	var assigns_60 vm.Value
	var variadic_QMARK__61 vm.Value
	var fixed_arity_62 vm.Value
	var v140 vm.Value
	var v142 vm.Value
	var f_143 vm.Value
	var arity_144 vm.Value
	var i_145 int
	var decls_146 vm.Value
	var final_decls_147 vm.Value
	var final_assigns_148 vm.Value
	var assigns_149 vm.Value
	var variadic_QMARK__150 vm.Value
	var fixed_arity_151 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arity_3, variadic_QMARK__5, f_6, arity_7, variadic_QMARK__8, v13, f_9, arity_10, variadic_QMARK__11, fixed_arity_16, f_17, arity_18, variadic_QMARK__19, i_20, decls_21, assigns_22, variadic_QMARK__23, fixed_arity_24, v245, v248, v251, v254, v257, v43, f_29, arity_30, i_31, decls_32, assigns_33, variadic_QMARK__34, fixed_arity_35, v246, v249, v252, v255, v258, f_36, arity_37, i_38, decls_39, assigns_40, variadic_QMARK__41, fixed_arity_42, v244, v247, v250, v253, v256, sym_156, arg_sym_160, v161, arg__12790_165, arg__12797_171, arg__12799_173, arg__12806_178, arg__12813_184, arg__12815_186, v187, arg__12821_190, arg__12825_192, arg__12831_196, arg__12835_198, arg__12836_199, arg__12843_203, arg__12847_205, arg__12853_209, arg__12857_211, arg__12858_212, v213, v215, f_216, arity_217, i_218, decls_219, assigns_220, variadic_QMARK__221, fixed_arity_222, f_45, arity_46, i_47, decls_48, final_decls_49, final_assigns_50, assigns_51, variadic_QMARK__52, fixed_arity_53, arg__12700_69, arg__12707_76, arg__12709_78, arg__12716_84, arg__12723_91, arg__12725_93, arg__12726_94, arg__12732_99, arg__12736_103, arg__12742_109, arg__12746_113, arg__12747_114, arg__12754_120, arg__12758_124, arg__12764_130, arg__12768_134, arg__12769_135, arg__12770_136, v137, f_54, arity_55, i_56, decls_57, final_decls_58, final_assigns_59, assigns_60, variadic_QMARK__61, fixed_arity_62, v140, v142, f_143, arity_144, i_145, decls_146, final_decls_147, final_assigns_148, assigns_149, variadic_QMARK__150, fixed_arity_151
	arity_3, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-arity").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	variadic_QMARK__5, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-variadic?").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(variadic_QMARK__5) {
		f_6 = arg0
		arity_7 = arity_3
		variadic_QMARK__8 = variadic_QMARK__5
		goto b1
	} else {
		f_9 = arg0
		arity_10 = arity_3
		variadic_QMARK__11 = variadic_QMARK__5
		goto b2
	}
b1:
	;
	v13 = rt.SubValue(arity_7, vm.Int(1))
	fixed_arity_16 = v13
	f_17 = f_6
	arity_18 = arity_7
	variadic_QMARK__19 = variadic_QMARK__8
	goto b3
b2:
	;
	fixed_arity_16 = arity_10
	f_17 = f_9
	arity_18 = arity_10
	variadic_QMARK__19 = variadic_QMARK__11
	goto b3
b3:
	;
	i_20 = 0
	decls_21 = vm.NewArrayVector([]vm.Value{})
	assigns_22 = vm.NewArrayVector([]vm.Value{})
	variadic_QMARK__23 = variadic_QMARK__19
	fixed_arity_24 = fixed_arity_16
	v245 = "a"
	v248 = "arg"
	v251 = "vm.Value"
	v254 = vm.NIL
	v257 = "="
	goto b4
b4:
	;
	v43 = rt.GeValue(vm.Int(i_20), fixed_arity_24)
	if v43 {
		f_29 = f_17
		arity_30 = arity_18
		i_31 = i_20
		decls_32 = decls_21
		assigns_33 = assigns_22
		variadic_QMARK__34 = variadic_QMARK__23
		fixed_arity_35 = fixed_arity_24
		v246 = v245
		v249 = v248
		v252 = v251
		v255 = v254
		v258 = v257
		goto b5
	} else {
		f_36 = f_17
		arity_37 = arity_18
		i_38 = i_20
		decls_39 = decls_21
		assigns_40 = assigns_22
		variadic_QMARK__41 = variadic_QMARK__23
		fixed_arity_42 = fixed_arity_24
		v244 = v245
		v247 = v248
		v250 = v251
		v253 = v254
		v256 = v257
		goto b6
	}
b5:
	;
	if vm.IsTruthy(variadic_QMARK__34) {
		f_45 = f_29
		arity_46 = arity_30
		i_47 = i_31
		decls_48 = decls_32
		final_decls_49 = decls_32
		final_assigns_50 = assigns_33
		assigns_51 = assigns_33
		variadic_QMARK__52 = variadic_QMARK__34
		fixed_arity_53 = fixed_arity_35
		goto b8
	} else {
		f_54 = f_29
		arity_55 = arity_30
		i_56 = i_31
		decls_57 = decls_32
		final_decls_58 = decls_32
		final_assigns_59 = assigns_33
		assigns_60 = assigns_33
		variadic_QMARK__61 = variadic_QMARK__34
		fixed_arity_62 = fixed_arity_35
		goto b9
	}
b6:
	;
	sym_156, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(v244), vm.Int(i_38)})
	if callErr != nil {
		return nil, callErr
	}
	arg_sym_160, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(v247), vm.Int(i_38)})
	if callErr != nil {
		return nil, callErr
	}
	v161 = i_38 + 1
	arg__12790_165, callErr = rt.InvokeValue(rt.LookupVar("gogen", "type").Deref(), []vm.Value{vm.String(v250)})
	if callErr != nil {
		return nil, callErr
	}
	arg__12797_171, callErr = rt.InvokeValue(rt.LookupVar("gogen", "type").Deref(), []vm.Value{vm.String(v250)})
	if callErr != nil {
		return nil, callErr
	}
	arg__12799_173, callErr = rt.InvokeValue(rt.LookupVar("gogen", "var-decl").Deref(), []vm.Value{sym_156, arg__12797_171, v253})
	if callErr != nil {
		return nil, callErr
	}
	arg__12806_178, callErr = rt.InvokeValue(rt.LookupVar("gogen", "type").Deref(), []vm.Value{vm.String(v250)})
	if callErr != nil {
		return nil, callErr
	}
	arg__12813_184, callErr = rt.InvokeValue(rt.LookupVar("gogen", "type").Deref(), []vm.Value{vm.String(v250)})
	if callErr != nil {
		return nil, callErr
	}
	arg__12815_186, callErr = rt.InvokeValue(rt.LookupVar("gogen", "var-decl").Deref(), []vm.Value{sym_156, arg__12813_184, v253})
	if callErr != nil {
		return nil, callErr
	}
	v187, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{decls_39, arg__12815_186})
	if callErr != nil {
		return nil, callErr
	}
	arg__12821_190, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{sym_156})
	if callErr != nil {
		return nil, callErr
	}
	arg__12825_192, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{arg_sym_160})
	if callErr != nil {
		return nil, callErr
	}
	arg__12831_196, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{sym_156})
	if callErr != nil {
		return nil, callErr
	}
	arg__12835_198, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{arg_sym_160})
	if callErr != nil {
		return nil, callErr
	}
	arg__12836_199, callErr = rt.InvokeValue(rt.LookupVar("gogen", "assign").Deref(), []vm.Value{vm.String(v256), arg__12831_196, arg__12835_198})
	if callErr != nil {
		return nil, callErr
	}
	arg__12843_203, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{sym_156})
	if callErr != nil {
		return nil, callErr
	}
	arg__12847_205, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{arg_sym_160})
	if callErr != nil {
		return nil, callErr
	}
	arg__12853_209, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{sym_156})
	if callErr != nil {
		return nil, callErr
	}
	arg__12857_211, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{arg_sym_160})
	if callErr != nil {
		return nil, callErr
	}
	arg__12858_212, callErr = rt.InvokeValue(rt.LookupVar("gogen", "assign").Deref(), []vm.Value{vm.String(v256), arg__12853_209, arg__12857_211})
	if callErr != nil {
		return nil, callErr
	}
	v213, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{assigns_40, arg__12858_212})
	if callErr != nil {
		return nil, callErr
	}
	i_20 = v161
	decls_21 = v187
	assigns_22 = v213
	variadic_QMARK__23 = variadic_QMARK__41
	fixed_arity_24 = fixed_arity_42
	v245 = v244
	v248 = v247
	v251 = v250
	v254 = v253
	v257 = v256
	goto b4
b7:
	;
	return v215, nil
b8:
	;
	arg__12700_69, callErr = rt.InvokeValue(rt.LookupVar("gogen", "type").Deref(), []vm.Value{vm.String("[]vm.Value")})
	if callErr != nil {
		return nil, callErr
	}
	arg__12707_76, callErr = rt.InvokeValue(rt.LookupVar("gogen", "type").Deref(), []vm.Value{vm.String("[]vm.Value")})
	if callErr != nil {
		return nil, callErr
	}
	arg__12709_78, callErr = rt.InvokeValue(rt.LookupVar("gogen", "var-decl").Deref(), []vm.Value{vm.String("args0"), arg__12707_76, vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	arg__12716_84, callErr = rt.InvokeValue(rt.LookupVar("gogen", "type").Deref(), []vm.Value{vm.String("[]vm.Value")})
	if callErr != nil {
		return nil, callErr
	}
	arg__12723_91, callErr = rt.InvokeValue(rt.LookupVar("gogen", "type").Deref(), []vm.Value{vm.String("[]vm.Value")})
	if callErr != nil {
		return nil, callErr
	}
	arg__12725_93, callErr = rt.InvokeValue(rt.LookupVar("gogen", "var-decl").Deref(), []vm.Value{vm.String("args0"), arg__12723_91, vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	arg__12726_94, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{final_decls_49, arg__12725_93})
	if callErr != nil {
		return nil, callErr
	}
	arg__12732_99, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("args0")})
	if callErr != nil {
		return nil, callErr
	}
	arg__12736_103, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("args")})
	if callErr != nil {
		return nil, callErr
	}
	arg__12742_109, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("args0")})
	if callErr != nil {
		return nil, callErr
	}
	arg__12746_113, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("args")})
	if callErr != nil {
		return nil, callErr
	}
	arg__12747_114, callErr = rt.InvokeValue(rt.LookupVar("gogen", "assign").Deref(), []vm.Value{vm.String("="), arg__12742_109, arg__12746_113})
	if callErr != nil {
		return nil, callErr
	}
	arg__12754_120, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("args0")})
	if callErr != nil {
		return nil, callErr
	}
	arg__12758_124, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("args")})
	if callErr != nil {
		return nil, callErr
	}
	arg__12764_130, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("args0")})
	if callErr != nil {
		return nil, callErr
	}
	arg__12768_134, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("args")})
	if callErr != nil {
		return nil, callErr
	}
	arg__12769_135, callErr = rt.InvokeValue(rt.LookupVar("gogen", "assign").Deref(), []vm.Value{vm.String("="), arg__12764_130, arg__12768_134})
	if callErr != nil {
		return nil, callErr
	}
	arg__12770_136, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{final_assigns_50, arg__12769_135})
	if callErr != nil {
		return nil, callErr
	}
	v137, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__12726_94, arg__12770_136})
	if callErr != nil {
		return nil, callErr
	}
	v142 = v137
	f_143 = f_45
	arity_144 = arity_46
	i_145 = i_47
	decls_146 = decls_48
	final_decls_147 = final_decls_49
	final_assigns_148 = final_assigns_50
	assigns_149 = assigns_51
	variadic_QMARK__150 = variadic_QMARK__52
	fixed_arity_151 = fixed_arity_53
	goto b10
b9:
	;
	v140, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{final_decls_58, final_assigns_59})
	if callErr != nil {
		return nil, callErr
	}
	v142 = v140
	f_143 = f_54
	arity_144 = arity_55
	i_145 = i_56
	decls_146 = decls_57
	final_decls_147 = final_decls_58
	final_assigns_148 = final_assigns_59
	assigns_149 = assigns_60
	variadic_QMARK__150 = variadic_QMARK__61
	fixed_arity_151 = fixed_arity_62
	goto b10
b10:
	;
	v215 = v142
	f_216 = f_143
	arity_217 = arity_144
	i_218 = i_145
	decls_219 = decls_146
	assigns_220 = assigns_149
	variadic_QMARK__221 = variadic_QMARK__150
	fixed_arity_222 = fixed_arity_151
	goto b7
}
func box_as_value(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var t_4 vm.Value
	var spec_6 vm.Value
	var expr_8 vm.Value
	var v22 vm.Value
	var f_9 vm.Value
	var closed_exprs_10 vm.Value
	var nid_11 vm.Value
	var t_12 vm.Value
	var spec_13 vm.Value
	var expr_14 vm.Value
	var f_15 vm.Value
	var closed_exprs_16 vm.Value
	var nid_17 vm.Value
	var t_18 vm.Value
	var spec_19 vm.Value
	var expr_20 vm.Value
	var v39 bool
	var v225 vm.Value
	var f_226 vm.Value
	var closed_exprs_227 vm.Value
	var nid_228 vm.Value
	var t_229 vm.Value
	var spec_230 vm.Value
	var expr_231 vm.Value
	var f_26 vm.Value
	var closed_exprs_27 vm.Value
	var nid_28 vm.Value
	var t_29 vm.Value
	var spec_30 vm.Value
	var expr_31 vm.Value
	var f_32 vm.Value
	var closed_exprs_33 vm.Value
	var nid_34 vm.Value
	var t_35 vm.Value
	var spec_36 vm.Value
	var expr_37 vm.Value
	var v55 bool
	var v217 vm.Value
	var f_218 vm.Value
	var closed_exprs_219 vm.Value
	var nid_220 vm.Value
	var t_221 vm.Value
	var spec_222 vm.Value
	var expr_223 vm.Value
	var f_42 vm.Value
	var closed_exprs_43 vm.Value
	var nid_44 vm.Value
	var t_45 vm.Value
	var spec_46 vm.Value
	var expr_47 vm.Value
	var arg__12884_60 vm.Value
	var arg__12890_66 vm.Value
	var arg__12892_68 vm.Value
	var arg__12895_70 vm.Value
	var arg__12900_75 vm.Value
	var arg__12906_81 vm.Value
	var arg__12908_83 vm.Value
	var arg__12911_85 vm.Value
	var v86 vm.Value
	var f_48 vm.Value
	var closed_exprs_49 vm.Value
	var nid_50 vm.Value
	var t_51 vm.Value
	var spec_52 vm.Value
	var expr_53 vm.Value
	var v101 bool
	var v209 vm.Value
	var f_210 vm.Value
	var closed_exprs_211 vm.Value
	var nid_212 vm.Value
	var t_213 vm.Value
	var spec_214 vm.Value
	var expr_215 vm.Value
	var f_88 vm.Value
	var closed_exprs_89 vm.Value
	var nid_90 vm.Value
	var t_91 vm.Value
	var spec_92 vm.Value
	var expr_93 vm.Value
	var arg__12917_105 vm.Value
	var arg__12922_109 vm.Value
	var v110 vm.Value
	var f_94 vm.Value
	var closed_exprs_95 vm.Value
	var nid_96 vm.Value
	var t_97 vm.Value
	var spec_98 vm.Value
	var expr_99 vm.Value
	var v125 bool
	var v201 vm.Value
	var f_202 vm.Value
	var closed_exprs_203 vm.Value
	var nid_204 vm.Value
	var t_205 vm.Value
	var spec_206 vm.Value
	var expr_207 vm.Value
	var f_112 vm.Value
	var closed_exprs_113 vm.Value
	var nid_114 vm.Value
	var t_115 vm.Value
	var spec_116 vm.Value
	var expr_117 vm.Value
	var arg__12928_129 vm.Value
	var arg__12933_133 vm.Value
	var v134 vm.Value
	var f_118 vm.Value
	var closed_exprs_119 vm.Value
	var nid_120 vm.Value
	var t_121 vm.Value
	var spec_122 vm.Value
	var expr_123 vm.Value
	var v149 bool
	var v193 vm.Value
	var f_194 vm.Value
	var closed_exprs_195 vm.Value
	var nid_196 vm.Value
	var t_197 vm.Value
	var spec_198 vm.Value
	var expr_199 vm.Value
	var f_136 vm.Value
	var closed_exprs_137 vm.Value
	var nid_138 vm.Value
	var t_139 vm.Value
	var spec_140 vm.Value
	var expr_141 vm.Value
	var arg__12939_153 vm.Value
	var arg__12944_157 vm.Value
	var v158 vm.Value
	var f_142 vm.Value
	var closed_exprs_143 vm.Value
	var nid_144 vm.Value
	var t_145 vm.Value
	var spec_146 vm.Value
	var expr_147 vm.Value
	var v185 vm.Value
	var f_186 vm.Value
	var closed_exprs_187 vm.Value
	var nid_188 vm.Value
	var t_189 vm.Value
	var spec_190 vm.Value
	var expr_191 vm.Value
	var f_160 vm.Value
	var closed_exprs_161 vm.Value
	var nid_162 vm.Value
	var t_163 vm.Value
	var spec_164 vm.Value
	var expr_165 vm.Value
	var f_166 vm.Value
	var closed_exprs_167 vm.Value
	var nid_168 vm.Value
	var t_169 vm.Value
	var spec_170 vm.Value
	var expr_171 vm.Value
	var v177 vm.Value
	var f_178 vm.Value
	var closed_exprs_179 vm.Value
	var nid_180 vm.Value
	var t_181 vm.Value
	var spec_182 vm.Value
	var expr_183 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = t_4, spec_6, expr_8, v22, f_9, closed_exprs_10, nid_11, t_12, spec_13, expr_14, f_15, closed_exprs_16, nid_17, t_18, spec_19, expr_20, v39, v225, f_226, closed_exprs_227, nid_228, t_229, spec_230, expr_231, f_26, closed_exprs_27, nid_28, t_29, spec_30, expr_31, f_32, closed_exprs_33, nid_34, t_35, spec_36, expr_37, v55, v217, f_218, closed_exprs_219, nid_220, t_221, spec_222, expr_223, f_42, closed_exprs_43, nid_44, t_45, spec_46, expr_47, arg__12884_60, arg__12890_66, arg__12892_68, arg__12895_70, arg__12900_75, arg__12906_81, arg__12908_83, arg__12911_85, v86, f_48, closed_exprs_49, nid_50, t_51, spec_52, expr_53, v101, v209, f_210, closed_exprs_211, nid_212, t_213, spec_214, expr_215, f_88, closed_exprs_89, nid_90, t_91, spec_92, expr_93, arg__12917_105, arg__12922_109, v110, f_94, closed_exprs_95, nid_96, t_97, spec_98, expr_99, v125, v201, f_202, closed_exprs_203, nid_204, t_205, spec_206, expr_207, f_112, closed_exprs_113, nid_114, t_115, spec_116, expr_117, arg__12928_129, arg__12933_133, v134, f_118, closed_exprs_119, nid_120, t_121, spec_122, expr_123, v149, v193, f_194, closed_exprs_195, nid_196, t_197, spec_198, expr_199, f_136, closed_exprs_137, nid_138, t_139, spec_140, expr_141, arg__12939_153, arg__12944_157, v158, f_142, closed_exprs_143, nid_144, t_145, spec_146, expr_147, v185, f_186, closed_exprs_187, nid_188, t_189, spec_190, expr_191, f_160, closed_exprs_161, nid_162, t_163, spec_164, expr_165, f_166, closed_exprs_167, nid_168, t_169, spec_170, expr_171, v177, f_178, closed_exprs_179, nid_180, t_181, spec_182, expr_183
	t_4, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg2, arg0})
	if callErr != nil {
		return nil, callErr
	}
	spec_6, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "go-type-spec").Deref(), []vm.Value{t_4})
	if callErr != nil {
		return nil, callErr
	}
	expr_8, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "value-expr").Deref(), []vm.Value{arg0, arg1, arg2})
	if callErr != nil {
		return nil, callErr
	}
	v22, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{expr_8})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v22) {
		f_9 = arg0
		closed_exprs_10 = arg1
		nid_11 = arg2
		t_12 = t_4
		spec_13 = spec_6
		expr_14 = expr_8
		goto b1
	} else {
		f_15 = arg0
		closed_exprs_16 = arg1
		nid_17 = arg2
		t_18 = t_4
		spec_19 = spec_6
		expr_20 = expr_8
		goto b2
	}
b1:
	;
	v225 = vm.NIL
	f_226 = f_9
	closed_exprs_227 = closed_exprs_10
	nid_228 = nid_11
	t_229 = t_12
	spec_230 = spec_13
	expr_231 = expr_14
	goto b3
b2:
	;
	v39 = spec_19 == vm.String("vm.Value")
	if v39 {
		f_26 = f_15
		closed_exprs_27 = closed_exprs_16
		nid_28 = nid_17
		t_29 = t_18
		spec_30 = spec_19
		expr_31 = expr_20
		goto b4
	} else {
		f_32 = f_15
		closed_exprs_33 = closed_exprs_16
		nid_34 = nid_17
		t_35 = t_18
		spec_36 = spec_19
		expr_37 = expr_20
		goto b5
	}
b3:
	;
	return v225, nil
b4:
	;
	v217 = expr_31
	f_218 = f_26
	closed_exprs_219 = closed_exprs_27
	nid_220 = nid_28
	t_221 = t_29
	spec_222 = spec_30
	expr_223 = expr_31
	goto b6
b5:
	;
	v55 = spec_36 == vm.String("bool")
	if v55 {
		f_42 = f_32
		closed_exprs_43 = closed_exprs_33
		nid_44 = nid_34
		t_45 = t_35
		spec_46 = spec_36
		expr_47 = expr_37
		goto b7
	} else {
		f_48 = f_32
		closed_exprs_49 = closed_exprs_33
		nid_50 = nid_34
		t_51 = t_35
		spec_52 = spec_36
		expr_53 = expr_37
		goto b8
	}
b6:
	;
	v225 = v217
	f_226 = f_218
	closed_exprs_227 = closed_exprs_219
	nid_228 = nid_220
	t_229 = t_221
	spec_230 = spec_222
	expr_231 = expr_223
	goto b3
b7:
	;
	arg__12884_60, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("vm")})
	if callErr != nil {
		return nil, callErr
	}
	arg__12890_66, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("vm")})
	if callErr != nil {
		return nil, callErr
	}
	arg__12892_68, callErr = rt.InvokeValue(rt.LookupVar("gogen", "field-sel").Deref(), []vm.Value{arg__12890_66, vm.String("Boolean")})
	if callErr != nil {
		return nil, callErr
	}
	arg__12895_70, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{expr_47})
	if callErr != nil {
		return nil, callErr
	}
	arg__12900_75, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("vm")})
	if callErr != nil {
		return nil, callErr
	}
	arg__12906_81, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("vm")})
	if callErr != nil {
		return nil, callErr
	}
	arg__12908_83, callErr = rt.InvokeValue(rt.LookupVar("gogen", "field-sel").Deref(), []vm.Value{arg__12906_81, vm.String("Boolean")})
	if callErr != nil {
		return nil, callErr
	}
	arg__12911_85, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{expr_47})
	if callErr != nil {
		return nil, callErr
	}
	v86, callErr = rt.InvokeValue(rt.LookupVar("gogen", "call").Deref(), []vm.Value{arg__12908_83, arg__12911_85})
	if callErr != nil {
		return nil, callErr
	}
	v209 = v86
	f_210 = f_42
	closed_exprs_211 = closed_exprs_43
	nid_212 = nid_44
	t_213 = t_45
	spec_214 = spec_46
	expr_215 = expr_47
	goto b9
b8:
	;
	v101 = spec_52 == vm.String("int")
	if v101 {
		f_88 = f_48
		closed_exprs_89 = closed_exprs_49
		nid_90 = nid_50
		t_91 = t_51
		spec_92 = spec_52
		expr_93 = expr_53
		goto b10
	} else {
		f_94 = f_48
		closed_exprs_95 = closed_exprs_49
		nid_96 = nid_50
		t_97 = t_51
		spec_98 = spec_52
		expr_99 = expr_53
		goto b11
	}
b9:
	;
	v217 = v209
	f_218 = f_210
	closed_exprs_219 = closed_exprs_211
	nid_220 = nid_212
	t_221 = t_213
	spec_222 = spec_214
	expr_223 = expr_215
	goto b6
b10:
	;
	arg__12917_105, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{expr_93})
	if callErr != nil {
		return nil, callErr
	}
	arg__12922_109, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{expr_93})
	if callErr != nil {
		return nil, callErr
	}
	v110, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "vm-call").Deref(), []vm.Value{vm.String("Int"), arg__12922_109})
	if callErr != nil {
		return nil, callErr
	}
	v201 = v110
	f_202 = f_88
	closed_exprs_203 = closed_exprs_89
	nid_204 = nid_90
	t_205 = t_91
	spec_206 = spec_92
	expr_207 = expr_93
	goto b12
b11:
	;
	v125 = spec_98 == vm.String("float64")
	if v125 {
		f_112 = f_94
		closed_exprs_113 = closed_exprs_95
		nid_114 = nid_96
		t_115 = t_97
		spec_116 = spec_98
		expr_117 = expr_99
		goto b13
	} else {
		f_118 = f_94
		closed_exprs_119 = closed_exprs_95
		nid_120 = nid_96
		t_121 = t_97
		spec_122 = spec_98
		expr_123 = expr_99
		goto b14
	}
b12:
	;
	v209 = v201
	f_210 = f_202
	closed_exprs_211 = closed_exprs_203
	nid_212 = nid_204
	t_213 = t_205
	spec_214 = spec_206
	expr_215 = expr_207
	goto b9
b13:
	;
	arg__12928_129, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{expr_117})
	if callErr != nil {
		return nil, callErr
	}
	arg__12933_133, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{expr_117})
	if callErr != nil {
		return nil, callErr
	}
	v134, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "vm-call").Deref(), []vm.Value{vm.String("Float"), arg__12933_133})
	if callErr != nil {
		return nil, callErr
	}
	v193 = v134
	f_194 = f_112
	closed_exprs_195 = closed_exprs_113
	nid_196 = nid_114
	t_197 = t_115
	spec_198 = spec_116
	expr_199 = expr_117
	goto b15
b14:
	;
	v149 = spec_122 == vm.String("string")
	if v149 {
		f_136 = f_118
		closed_exprs_137 = closed_exprs_119
		nid_138 = nid_120
		t_139 = t_121
		spec_140 = spec_122
		expr_141 = expr_123
		goto b16
	} else {
		f_142 = f_118
		closed_exprs_143 = closed_exprs_119
		nid_144 = nid_120
		t_145 = t_121
		spec_146 = spec_122
		expr_147 = expr_123
		goto b17
	}
b15:
	;
	v201 = v193
	f_202 = f_194
	closed_exprs_203 = closed_exprs_195
	nid_204 = nid_196
	t_205 = t_197
	spec_206 = spec_198
	expr_207 = expr_199
	goto b12
b16:
	;
	arg__12939_153, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{expr_141})
	if callErr != nil {
		return nil, callErr
	}
	arg__12944_157, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{expr_141})
	if callErr != nil {
		return nil, callErr
	}
	v158, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "vm-call").Deref(), []vm.Value{vm.String("String"), arg__12944_157})
	if callErr != nil {
		return nil, callErr
	}
	v185 = v158
	f_186 = f_136
	closed_exprs_187 = closed_exprs_137
	nid_188 = nid_138
	t_189 = t_139
	spec_190 = spec_140
	expr_191 = expr_141
	goto b18
b17:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		f_160 = f_142
		closed_exprs_161 = closed_exprs_143
		nid_162 = nid_144
		t_163 = t_145
		spec_164 = spec_146
		expr_165 = expr_147
		goto b19
	} else {
		f_166 = f_142
		closed_exprs_167 = closed_exprs_143
		nid_168 = nid_144
		t_169 = t_145
		spec_170 = spec_146
		expr_171 = expr_147
		goto b20
	}
b18:
	;
	v193 = v185
	f_194 = f_186
	closed_exprs_195 = closed_exprs_187
	nid_196 = nid_188
	t_197 = t_189
	spec_198 = spec_190
	expr_199 = expr_191
	goto b15
b19:
	;
	v177 = expr_165
	f_178 = f_160
	closed_exprs_179 = closed_exprs_161
	nid_180 = nid_162
	t_181 = t_163
	spec_182 = spec_164
	expr_183 = expr_165
	goto b21
b20:
	;
	v177 = vm.NIL
	f_178 = f_166
	closed_exprs_179 = closed_exprs_167
	nid_180 = nid_168
	t_181 = t_169
	spec_182 = spec_170
	expr_183 = expr_171
	goto b21
b21:
	;
	v185 = v177
	f_186 = f_178
	closed_exprs_187 = closed_exprs_179
	nid_188 = nid_180
	t_189 = t_181
	spec_190 = spec_182
	expr_191 = expr_183
	goto b18
}
func boxed_list_expr(arg0 vm.Value) (vm.Value, error) {
	var arg__12948_4 vm.Value
	var arg__12953_7 vm.Value
	var v8 vm.Value
	var xs_1 vm.Value
	var v13 vm.Value
	var xs_2 vm.Value
	var arg__12961_17 vm.Value
	var arg__12967_21 vm.Value
	var elems_22 vm.Value
	var v30 vm.Value
	var v67 vm.Value
	var xs_68 vm.Value
	var xs_23 vm.Value
	var elems_24 vm.Value
	var arg__12978_37 vm.Value
	var arg__12984_42 vm.Value
	var arg__12986_43 vm.Value
	var arg__12987_44 vm.Value
	var arg__12994_51 vm.Value
	var arg__13000_56 vm.Value
	var arg__13002_57 vm.Value
	var arg__13003_58 vm.Value
	var v59 vm.Value
	var xs_25 vm.Value
	var elems_26 vm.Value
	var v63 vm.Value
	var xs_64 vm.Value
	var elems_65 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__12948_4, arg__12953_7, v8, xs_1, v13, xs_2, arg__12961_17, arg__12967_21, elems_22, v30, v67, xs_68, xs_23, elems_24, arg__12978_37, arg__12984_42, arg__12986_43, arg__12987_44, arg__12994_51, arg__13000_56, arg__13002_57, arg__13003_58, v59, xs_25, elems_26, v63, xs_64, elems_65
	arg__12948_4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__12953_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v8, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "zero?").Deref(), []vm.Value{arg__12953_7})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v8) {
		xs_1 = arg0
		goto b1
	} else {
		xs_2 = arg0
		goto b2
	}
b1:
	;
	v13, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "vm-sel").Deref(), []vm.Value{vm.String("EmptyList")})
	if callErr != nil {
		return nil, callErr
	}
	v67 = v13
	xs_68 = xs_1
	goto b3
b2:
	;
	arg__12961_17, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{xs_2})
	if callErr != nil {
		return nil, callErr
	}
	arg__12967_21, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{xs_2})
	if callErr != nil {
		return nil, callErr
	}
	elems_22, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapv").Deref(), []vm.Value{rt.LookupVar("ir.lower-go", "boxed-value-expr").Deref(), arg__12967_21})
	if callErr != nil {
		return nil, callErr
	}
	v30, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "every?").Deref(), []vm.Value{rt.LookupVar("clojure.core", "identity").Deref(), elems_22})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v30) {
		xs_23 = xs_2
		elems_24 = elems_22
		goto b4
	} else {
		xs_25 = xs_2
		elems_26 = elems_22
		goto b5
	}
b3:
	;
	return v67, nil
b4:
	;
	arg__12978_37, callErr = rt.InvokeValue(rt.LookupVar("gogen", "type").Deref(), []vm.Value{vm.String("[]vm.Value")})
	if callErr != nil {
		return nil, callErr
	}
	arg__12984_42, callErr = rt.InvokeValue(rt.LookupVar("gogen", "type").Deref(), []vm.Value{vm.String("[]vm.Value")})
	if callErr != nil {
		return nil, callErr
	}
	arg__12986_43, callErr = rt.InvokeValue(rt.LookupVar("gogen", "composite-lit").Deref(), []vm.Value{arg__12984_42, elems_24})
	if callErr != nil {
		return nil, callErr
	}
	arg__12987_44, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__12986_43})
	if callErr != nil {
		return nil, callErr
	}
	arg__12994_51, callErr = rt.InvokeValue(rt.LookupVar("gogen", "type").Deref(), []vm.Value{vm.String("[]vm.Value")})
	if callErr != nil {
		return nil, callErr
	}
	arg__13000_56, callErr = rt.InvokeValue(rt.LookupVar("gogen", "type").Deref(), []vm.Value{vm.String("[]vm.Value")})
	if callErr != nil {
		return nil, callErr
	}
	arg__13002_57, callErr = rt.InvokeValue(rt.LookupVar("gogen", "composite-lit").Deref(), []vm.Value{arg__13000_56, elems_24})
	if callErr != nil {
		return nil, callErr
	}
	arg__13003_58, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__13002_57})
	if callErr != nil {
		return nil, callErr
	}
	v59, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "vm-call").Deref(), []vm.Value{vm.String("NewList"), arg__13003_58})
	if callErr != nil {
		return nil, callErr
	}
	v63 = v59
	xs_64 = xs_23
	elems_65 = elems_24
	goto b6
b5:
	;
	v63 = vm.NIL
	xs_64 = xs_25
	elems_65 = elems_26
	goto b6
b6:
	;
	v67 = v63
	xs_68 = xs_64
	goto b3
}
func boxed_map_expr(arg0 vm.Value) (vm.Value, error) {
	var arg__13007_4 vm.Value
	var arg__13012_7 vm.Value
	var v8 vm.Value
	var m_1 vm.Value
	var v13 vm.Value
	var m_2 vm.Value
	var entries_16 vm.Value
	var v164 vm.Value
	var m_165 vm.Value
	var remaining_17 vm.Value
	var out_18 vm.Value
	var v30 vm.Value
	var m_21 vm.Value
	var entries_22 vm.Value
	var remaining_23 vm.Value
	var out_24 vm.Value
	var m_25 vm.Value
	var entries_26 vm.Value
	var remaining_27 vm.Value
	var out_28 vm.Value
	var entry_34 vm.Value
	var arg__13028_36 vm.Value
	var arg__13033_39 vm.Value
	var kexpr_40 vm.Value
	var arg__13037_42 vm.Value
	var arg__13042_45 vm.Value
	var vexpr_46 vm.Value
	var or__x_62 vm.Value
	var elems_110 vm.Value
	var m_111 vm.Value
	var entries_112 vm.Value
	var remaining_113 vm.Value
	var out_114 vm.Value
	var m_47 vm.Value
	var entries_48 vm.Value
	var remaining_49 vm.Value
	var out_50 vm.Value
	var entry_51 vm.Value
	var kexpr_52 vm.Value
	var vexpr_53 vm.Value
	var m_54 vm.Value
	var entries_55 vm.Value
	var remaining_56 vm.Value
	var out_57 vm.Value
	var entry_58 vm.Value
	var kexpr_59 vm.Value
	var vexpr_60 vm.Value
	var v97 vm.Value
	var v99 vm.Value
	var v101 vm.Value
	var m_102 vm.Value
	var entries_103 vm.Value
	var remaining_104 vm.Value
	var out_105 vm.Value
	var entry_106 vm.Value
	var kexpr_107 vm.Value
	var vexpr_108 vm.Value
	var m_63 vm.Value
	var entries_64 vm.Value
	var remaining_65 vm.Value
	var out_66 vm.Value
	var entry_67 vm.Value
	var kexpr_68 vm.Value
	var vexpr_69 vm.Value
	var or__x_70 vm.Value
	var m_71 vm.Value
	var entries_72 vm.Value
	var remaining_73 vm.Value
	var out_74 vm.Value
	var entry_75 vm.Value
	var kexpr_76 vm.Value
	var vexpr_77 vm.Value
	var or__x_78 vm.Value
	var v82 vm.Value
	var v84 vm.Value
	var m_85 vm.Value
	var entries_86 vm.Value
	var remaining_87 vm.Value
	var out_88 vm.Value
	var entry_89 vm.Value
	var kexpr_90 vm.Value
	var vexpr_91 vm.Value
	var or__x_92 vm.Value
	var elems_115 vm.Value
	var m_116 vm.Value
	var entries_117 vm.Value
	var remaining_118 vm.Value
	var out_119 vm.Value
	var arg__13064_131 vm.Value
	var arg__13070_136 vm.Value
	var arg__13072_137 vm.Value
	var arg__13073_138 vm.Value
	var arg__13080_145 vm.Value
	var arg__13086_150 vm.Value
	var arg__13088_151 vm.Value
	var arg__13089_152 vm.Value
	var v153 vm.Value
	var elems_120 vm.Value
	var m_121 vm.Value
	var entries_122 vm.Value
	var remaining_123 vm.Value
	var out_124 vm.Value
	var v157 vm.Value
	var elems_158 vm.Value
	var m_159 vm.Value
	var entries_160 vm.Value
	var remaining_161 vm.Value
	var out_162 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__13007_4, arg__13012_7, v8, m_1, v13, m_2, entries_16, v164, m_165, remaining_17, out_18, v30, m_21, entries_22, remaining_23, out_24, m_25, entries_26, remaining_27, out_28, entry_34, arg__13028_36, arg__13033_39, kexpr_40, arg__13037_42, arg__13042_45, vexpr_46, or__x_62, elems_110, m_111, entries_112, remaining_113, out_114, m_47, entries_48, remaining_49, out_50, entry_51, kexpr_52, vexpr_53, m_54, entries_55, remaining_56, out_57, entry_58, kexpr_59, vexpr_60, v97, v99, v101, m_102, entries_103, remaining_104, out_105, entry_106, kexpr_107, vexpr_108, m_63, entries_64, remaining_65, out_66, entry_67, kexpr_68, vexpr_69, or__x_70, m_71, entries_72, remaining_73, out_74, entry_75, kexpr_76, vexpr_77, or__x_78, v82, v84, m_85, entries_86, remaining_87, out_88, entry_89, kexpr_90, vexpr_91, or__x_92, elems_115, m_116, entries_117, remaining_118, out_119, arg__13064_131, arg__13070_136, arg__13072_137, arg__13073_138, arg__13080_145, arg__13086_150, arg__13088_151, arg__13089_152, v153, elems_120, m_121, entries_122, remaining_123, out_124, v157, elems_158, m_159, entries_160, remaining_161, out_162
	arg__13007_4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__13012_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v8, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "zero?").Deref(), []vm.Value{arg__13012_7})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v8) {
		m_1 = arg0
		goto b1
	} else {
		m_2 = arg0
		goto b2
	}
b1:
	;
	v13, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "vm-sel").Deref(), []vm.Value{vm.String("EmptyPersistentMap")})
	if callErr != nil {
		return nil, callErr
	}
	v164 = v13
	m_165 = m_1
	goto b3
b2:
	;
	entries_16, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{m_2})
	if callErr != nil {
		return nil, callErr
	}
	remaining_17 = entries_16
	out_18 = vm.NewArrayVector([]vm.Value{})
	goto b4
b3:
	;
	return v164, nil
b4:
	;
	v30, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{remaining_17})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v30) {
		m_21 = m_2
		entries_22 = entries_16
		remaining_23 = remaining_17
		out_24 = out_18
		goto b5
	} else {
		m_25 = m_2
		entries_26 = entries_16
		remaining_27 = remaining_17
		out_28 = out_18
		goto b6
	}
b5:
	;
	elems_110 = out_24
	m_111 = m_21
	entries_112 = entries_22
	remaining_113 = remaining_23
	out_114 = out_24
	goto b7
b6:
	;
	entry_34, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{remaining_27})
	if callErr != nil {
		return nil, callErr
	}
	arg__13028_36, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{entry_34})
	if callErr != nil {
		return nil, callErr
	}
	arg__13033_39, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{entry_34})
	if callErr != nil {
		return nil, callErr
	}
	kexpr_40, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "boxed-value-expr").Deref(), []vm.Value{arg__13033_39})
	if callErr != nil {
		return nil, callErr
	}
	arg__13037_42, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "second").Deref(), []vm.Value{entry_34})
	if callErr != nil {
		return nil, callErr
	}
	arg__13042_45, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "second").Deref(), []vm.Value{entry_34})
	if callErr != nil {
		return nil, callErr
	}
	vexpr_46, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "boxed-value-expr").Deref(), []vm.Value{arg__13042_45})
	if callErr != nil {
		return nil, callErr
	}
	or__x_62, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{kexpr_40})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(or__x_62) {
		m_63 = m_25
		entries_64 = entries_26
		remaining_65 = remaining_27
		out_66 = out_28
		entry_67 = entry_34
		kexpr_68 = kexpr_40
		vexpr_69 = vexpr_46
		or__x_70 = or__x_62
		goto b11
	} else {
		m_71 = m_25
		entries_72 = entries_26
		remaining_73 = remaining_27
		out_74 = out_28
		entry_75 = entry_34
		kexpr_76 = kexpr_40
		vexpr_77 = vexpr_46
		or__x_78 = or__x_62
		goto b12
	}
b7:
	;
	if vm.IsTruthy(elems_110) {
		elems_115 = elems_110
		m_116 = m_111
		entries_117 = entries_112
		remaining_118 = remaining_113
		out_119 = out_114
		goto b14
	} else {
		elems_120 = elems_110
		m_121 = m_111
		entries_122 = entries_112
		remaining_123 = remaining_113
		out_124 = out_114
		goto b15
	}
b8:
	;
	v101 = vm.NIL
	m_102 = m_47
	entries_103 = entries_48
	remaining_104 = remaining_49
	out_105 = out_50
	entry_106 = entry_51
	kexpr_107 = kexpr_52
	vexpr_108 = vexpr_53
	goto b10
b9:
	;
	v97, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "next").Deref(), []vm.Value{remaining_56})
	if callErr != nil {
		return nil, callErr
	}
	v99, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{out_57, kexpr_59, vexpr_60})
	if callErr != nil {
		return nil, callErr
	}
	remaining_17 = v97
	out_18 = v99
	goto b4
b10:
	;
	elems_110 = v101
	m_111 = m_102
	entries_112 = entries_103
	remaining_113 = remaining_104
	out_114 = out_105
	goto b7
b11:
	;
	v84 = or__x_70
	m_85 = m_63
	entries_86 = entries_64
	remaining_87 = remaining_65
	out_88 = out_66
	entry_89 = entry_67
	kexpr_90 = kexpr_68
	vexpr_91 = vexpr_69
	or__x_92 = or__x_70
	goto b13
b12:
	;
	v82, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{vexpr_77})
	if callErr != nil {
		return nil, callErr
	}
	v84 = v82
	m_85 = m_71
	entries_86 = entries_72
	remaining_87 = remaining_73
	out_88 = out_74
	entry_89 = entry_75
	kexpr_90 = kexpr_76
	vexpr_91 = vexpr_77
	or__x_92 = or__x_78
	goto b13
b13:
	;
	if vm.IsTruthy(v84) {
		m_47 = m_85
		entries_48 = entries_86
		remaining_49 = remaining_87
		out_50 = out_88
		entry_51 = entry_89
		kexpr_52 = kexpr_90
		vexpr_53 = vexpr_91
		goto b8
	} else {
		m_54 = m_85
		entries_55 = entries_86
		remaining_56 = remaining_87
		out_57 = out_88
		entry_58 = entry_89
		kexpr_59 = kexpr_90
		vexpr_60 = vexpr_91
		goto b9
	}
b14:
	;
	arg__13064_131, callErr = rt.InvokeValue(rt.LookupVar("gogen", "type").Deref(), []vm.Value{vm.String("[]vm.Value")})
	if callErr != nil {
		return nil, callErr
	}
	arg__13070_136, callErr = rt.InvokeValue(rt.LookupVar("gogen", "type").Deref(), []vm.Value{vm.String("[]vm.Value")})
	if callErr != nil {
		return nil, callErr
	}
	arg__13072_137, callErr = rt.InvokeValue(rt.LookupVar("gogen", "composite-lit").Deref(), []vm.Value{arg__13070_136, elems_115})
	if callErr != nil {
		return nil, callErr
	}
	arg__13073_138, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__13072_137})
	if callErr != nil {
		return nil, callErr
	}
	arg__13080_145, callErr = rt.InvokeValue(rt.LookupVar("gogen", "type").Deref(), []vm.Value{vm.String("[]vm.Value")})
	if callErr != nil {
		return nil, callErr
	}
	arg__13086_150, callErr = rt.InvokeValue(rt.LookupVar("gogen", "type").Deref(), []vm.Value{vm.String("[]vm.Value")})
	if callErr != nil {
		return nil, callErr
	}
	arg__13088_151, callErr = rt.InvokeValue(rt.LookupVar("gogen", "composite-lit").Deref(), []vm.Value{arg__13086_150, elems_115})
	if callErr != nil {
		return nil, callErr
	}
	arg__13089_152, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__13088_151})
	if callErr != nil {
		return nil, callErr
	}
	v153, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "vm-call").Deref(), []vm.Value{vm.String("NewPersistentMap"), arg__13089_152})
	if callErr != nil {
		return nil, callErr
	}
	v157 = v153
	elems_158 = elems_115
	m_159 = m_116
	entries_160 = entries_117
	remaining_161 = remaining_118
	out_162 = out_119
	goto b16
b15:
	;
	v157 = vm.NIL
	elems_158 = elems_120
	m_159 = m_121
	entries_160 = entries_122
	remaining_161 = remaining_123
	out_162 = out_124
	goto b16
b16:
	;
	v164 = v157
	m_165 = m_159
	goto b3
}
func boxed_value_expr(arg0 vm.Value) (vm.Value, error) {
	var v4 vm.Value
	var v_1 vm.Value
	var v9 vm.Value
	var v_2 vm.Value
	var v14 bool
	var v247 vm.Value
	var v_248 vm.Value
	var v_11 vm.Value
	var v19 vm.Value
	var v_12 vm.Value
	var v24 bool
	var v244 vm.Value
	var v_245 vm.Value
	var v_21 vm.Value
	var v29 vm.Value
	var v_22 vm.Value
	var v34 vm.Value
	var v241 vm.Value
	var v_242 vm.Value
	var v_31 vm.Value
	var arg__13114_39 vm.Value
	var arg__13115_40 vm.Value
	var arg__13122_45 vm.Value
	var arg__13123_46 vm.Value
	var v47 vm.Value
	var v_32 vm.Value
	var v52 vm.Value
	var v238 vm.Value
	var v_239 vm.Value
	var v_49 vm.Value
	var arg__13132_57 vm.Value
	var arg__13133_58 vm.Value
	var arg__13140_63 vm.Value
	var arg__13141_64 vm.Value
	var v65 vm.Value
	var v_50 vm.Value
	var v70 vm.Value
	var v235 vm.Value
	var v_236 vm.Value
	var v_67 vm.Value
	var arg__13150_75 vm.Value
	var arg__13151_76 vm.Value
	var arg__13158_81 vm.Value
	var arg__13159_82 vm.Value
	var v83 vm.Value
	var v_68 vm.Value
	var v88 vm.Value
	var v232 vm.Value
	var v_233 vm.Value
	var v_85 vm.Value
	var arg__13168_93 vm.Value
	var arg__13169_94 vm.Value
	var arg__13176_99 vm.Value
	var arg__13177_100 vm.Value
	var v101 vm.Value
	var v_86 vm.Value
	var v106 vm.Value
	var v229 vm.Value
	var v_230 vm.Value
	var v_103 vm.Value
	var arg__13186_111 vm.Value
	var arg__13192_115 vm.Value
	var arg__13194_117 vm.Value
	var arg__13199_120 vm.Value
	var arg__13205_124 vm.Value
	var arg__13207_126 vm.Value
	var arg__13208_127 vm.Value
	var arg__13209_128 vm.Value
	var arg__13216_133 vm.Value
	var arg__13222_137 vm.Value
	var arg__13224_139 vm.Value
	var arg__13229_142 vm.Value
	var arg__13235_146 vm.Value
	var arg__13237_148 vm.Value
	var arg__13238_149 vm.Value
	var arg__13239_150 vm.Value
	var v151 vm.Value
	var v_104 vm.Value
	var v156 vm.Value
	var v226 vm.Value
	var v_227 vm.Value
	var v_153 vm.Value
	var arg__13248_161 vm.Value
	var arg__13253_164 vm.Value
	var arg__13254_165 vm.Value
	var arg__13255_166 vm.Value
	var arg__13262_171 vm.Value
	var arg__13267_174 vm.Value
	var arg__13268_175 vm.Value
	var arg__13269_176 vm.Value
	var v177 vm.Value
	var v_154 vm.Value
	var v182 vm.Value
	var v223 vm.Value
	var v_224 vm.Value
	var v_179 vm.Value
	var v185 vm.Value
	var v_180 vm.Value
	var v190 vm.Value
	var v220 vm.Value
	var v_221 vm.Value
	var v_187 vm.Value
	var v193 vm.Value
	var v_188 vm.Value
	var v198 vm.Value
	var v217 vm.Value
	var v_218 vm.Value
	var v_195 vm.Value
	var v201 vm.Value
	var v_196 vm.Value
	var v214 vm.Value
	var v_215 vm.Value
	var v_203 vm.Value
	var v_204 vm.Value
	var v211 vm.Value
	var v_212 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = v4, v_1, v9, v_2, v14, v247, v_248, v_11, v19, v_12, v24, v244, v_245, v_21, v29, v_22, v34, v241, v_242, v_31, arg__13114_39, arg__13115_40, arg__13122_45, arg__13123_46, v47, v_32, v52, v238, v_239, v_49, arg__13132_57, arg__13133_58, arg__13140_63, arg__13141_64, v65, v_50, v70, v235, v_236, v_67, arg__13150_75, arg__13151_76, arg__13158_81, arg__13159_82, v83, v_68, v88, v232, v_233, v_85, arg__13168_93, arg__13169_94, arg__13176_99, arg__13177_100, v101, v_86, v106, v229, v_230, v_103, arg__13186_111, arg__13192_115, arg__13194_117, arg__13199_120, arg__13205_124, arg__13207_126, arg__13208_127, arg__13209_128, arg__13216_133, arg__13222_137, arg__13224_139, arg__13229_142, arg__13235_146, arg__13237_148, arg__13238_149, arg__13239_150, v151, v_104, v156, v226, v_227, v_153, arg__13248_161, arg__13253_164, arg__13254_165, arg__13255_166, arg__13262_171, arg__13267_174, arg__13268_175, arg__13269_176, v177, v_154, v182, v223, v_224, v_179, v185, v_180, v190, v220, v_221, v_187, v193, v_188, v198, v217, v_218, v_195, v201, v_196, v214, v_215, v_203, v_204, v211, v_212
	v4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v4) {
		v_1 = arg0
		goto b1
	} else {
		v_2 = arg0
		goto b2
	}
b1:
	;
	v9, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "vm-sel").Deref(), []vm.Value{vm.String("NIL")})
	if callErr != nil {
		return nil, callErr
	}
	v247 = v9
	v_248 = v_1
	goto b3
b2:
	;
	v14 = v_2 == vm.Boolean(true)
	if v14 {
		v_11 = v_2
		goto b4
	} else {
		v_12 = v_2
		goto b5
	}
b3:
	;
	return v247, nil
b4:
	;
	v19, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "vm-sel").Deref(), []vm.Value{vm.String("TRUE")})
	if callErr != nil {
		return nil, callErr
	}
	v244 = v19
	v_245 = v_11
	goto b6
b5:
	;
	v24 = v_12 == vm.Boolean(false)
	if v24 {
		v_21 = v_12
		goto b7
	} else {
		v_22 = v_12
		goto b8
	}
b6:
	;
	v247 = v244
	v_248 = v_245
	goto b3
b7:
	;
	v29, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "vm-sel").Deref(), []vm.Value{vm.String("FALSE")})
	if callErr != nil {
		return nil, callErr
	}
	v241 = v29
	v_242 = v_21
	goto b9
b8:
	;
	v34, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "integer?").Deref(), []vm.Value{v_22})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v34) {
		v_31 = v_22
		goto b10
	} else {
		v_32 = v_22
		goto b11
	}
b9:
	;
	v244 = v241
	v_245 = v_242
	goto b6
b10:
	;
	arg__13114_39, callErr = rt.InvokeValue(rt.LookupVar("gogen", "int-lit").Deref(), []vm.Value{v_31})
	if callErr != nil {
		return nil, callErr
	}
	arg__13115_40, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__13114_39})
	if callErr != nil {
		return nil, callErr
	}
	arg__13122_45, callErr = rt.InvokeValue(rt.LookupVar("gogen", "int-lit").Deref(), []vm.Value{v_31})
	if callErr != nil {
		return nil, callErr
	}
	arg__13123_46, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__13122_45})
	if callErr != nil {
		return nil, callErr
	}
	v47, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "vm-call").Deref(), []vm.Value{vm.String("Int"), arg__13123_46})
	if callErr != nil {
		return nil, callErr
	}
	v238 = v47
	v_239 = v_31
	goto b12
b11:
	;
	v52, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "float?").Deref(), []vm.Value{v_32})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v52) {
		v_49 = v_32
		goto b13
	} else {
		v_50 = v_32
		goto b14
	}
b12:
	;
	v241 = v238
	v_242 = v_239
	goto b9
b13:
	;
	arg__13132_57, callErr = rt.InvokeValue(rt.LookupVar("gogen", "float-lit").Deref(), []vm.Value{v_49})
	if callErr != nil {
		return nil, callErr
	}
	arg__13133_58, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__13132_57})
	if callErr != nil {
		return nil, callErr
	}
	arg__13140_63, callErr = rt.InvokeValue(rt.LookupVar("gogen", "float-lit").Deref(), []vm.Value{v_49})
	if callErr != nil {
		return nil, callErr
	}
	arg__13141_64, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__13140_63})
	if callErr != nil {
		return nil, callErr
	}
	v65, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "vm-call").Deref(), []vm.Value{vm.String("Float"), arg__13141_64})
	if callErr != nil {
		return nil, callErr
	}
	v235 = v65
	v_236 = v_49
	goto b15
b14:
	;
	v70, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "string?").Deref(), []vm.Value{v_50})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v70) {
		v_67 = v_50
		goto b16
	} else {
		v_68 = v_50
		goto b17
	}
b15:
	;
	v238 = v235
	v_239 = v_236
	goto b12
b16:
	;
	arg__13150_75, callErr = rt.InvokeValue(rt.LookupVar("gogen", "string-lit").Deref(), []vm.Value{v_67})
	if callErr != nil {
		return nil, callErr
	}
	arg__13151_76, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__13150_75})
	if callErr != nil {
		return nil, callErr
	}
	arg__13158_81, callErr = rt.InvokeValue(rt.LookupVar("gogen", "string-lit").Deref(), []vm.Value{v_67})
	if callErr != nil {
		return nil, callErr
	}
	arg__13159_82, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__13158_81})
	if callErr != nil {
		return nil, callErr
	}
	v83, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "vm-call").Deref(), []vm.Value{vm.String("String"), arg__13159_82})
	if callErr != nil {
		return nil, callErr
	}
	v232 = v83
	v_233 = v_67
	goto b18
b17:
	;
	v88, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "char?").Deref(), []vm.Value{v_68})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v88) {
		v_85 = v_68
		goto b19
	} else {
		v_86 = v_68
		goto b20
	}
b18:
	;
	v235 = v232
	v_236 = v_233
	goto b15
b19:
	;
	arg__13168_93, callErr = rt.InvokeValue(rt.LookupVar("gogen", "char-lit").Deref(), []vm.Value{v_85})
	if callErr != nil {
		return nil, callErr
	}
	arg__13169_94, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__13168_93})
	if callErr != nil {
		return nil, callErr
	}
	arg__13176_99, callErr = rt.InvokeValue(rt.LookupVar("gogen", "char-lit").Deref(), []vm.Value{v_85})
	if callErr != nil {
		return nil, callErr
	}
	arg__13177_100, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__13176_99})
	if callErr != nil {
		return nil, callErr
	}
	v101, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "vm-call").Deref(), []vm.Value{vm.String("Char"), arg__13177_100})
	if callErr != nil {
		return nil, callErr
	}
	v229 = v101
	v_230 = v_85
	goto b21
b20:
	;
	v106, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "keyword?").Deref(), []vm.Value{v_86})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v106) {
		v_103 = v_86
		goto b22
	} else {
		v_104 = v_86
		goto b23
	}
b21:
	;
	v232 = v229
	v_233 = v_230
	goto b18
b22:
	;
	arg__13186_111, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{v_103})
	if callErr != nil {
		return nil, callErr
	}
	arg__13192_115, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{v_103})
	if callErr != nil {
		return nil, callErr
	}
	arg__13194_117, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "subs").Deref(), []vm.Value{arg__13192_115, vm.Int(1)})
	if callErr != nil {
		return nil, callErr
	}
	arg__13199_120, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{v_103})
	if callErr != nil {
		return nil, callErr
	}
	arg__13205_124, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{v_103})
	if callErr != nil {
		return nil, callErr
	}
	arg__13207_126, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "subs").Deref(), []vm.Value{arg__13205_124, vm.Int(1)})
	if callErr != nil {
		return nil, callErr
	}
	arg__13208_127, callErr = rt.InvokeValue(rt.LookupVar("gogen", "string-lit").Deref(), []vm.Value{arg__13207_126})
	if callErr != nil {
		return nil, callErr
	}
	arg__13209_128, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__13208_127})
	if callErr != nil {
		return nil, callErr
	}
	arg__13216_133, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{v_103})
	if callErr != nil {
		return nil, callErr
	}
	arg__13222_137, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{v_103})
	if callErr != nil {
		return nil, callErr
	}
	arg__13224_139, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "subs").Deref(), []vm.Value{arg__13222_137, vm.Int(1)})
	if callErr != nil {
		return nil, callErr
	}
	arg__13229_142, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{v_103})
	if callErr != nil {
		return nil, callErr
	}
	arg__13235_146, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{v_103})
	if callErr != nil {
		return nil, callErr
	}
	arg__13237_148, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "subs").Deref(), []vm.Value{arg__13235_146, vm.Int(1)})
	if callErr != nil {
		return nil, callErr
	}
	arg__13238_149, callErr = rt.InvokeValue(rt.LookupVar("gogen", "string-lit").Deref(), []vm.Value{arg__13237_148})
	if callErr != nil {
		return nil, callErr
	}
	arg__13239_150, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__13238_149})
	if callErr != nil {
		return nil, callErr
	}
	v151, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "vm-call").Deref(), []vm.Value{vm.String("Keyword"), arg__13239_150})
	if callErr != nil {
		return nil, callErr
	}
	v226 = v151
	v_227 = v_103
	goto b24
b23:
	;
	v156, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "symbol?").Deref(), []vm.Value{v_104})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v156) {
		v_153 = v_104
		goto b25
	} else {
		v_154 = v_104
		goto b26
	}
b24:
	;
	v229 = v226
	v_230 = v_227
	goto b21
b25:
	;
	arg__13248_161, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{v_153})
	if callErr != nil {
		return nil, callErr
	}
	arg__13253_164, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{v_153})
	if callErr != nil {
		return nil, callErr
	}
	arg__13254_165, callErr = rt.InvokeValue(rt.LookupVar("gogen", "string-lit").Deref(), []vm.Value{arg__13253_164})
	if callErr != nil {
		return nil, callErr
	}
	arg__13255_166, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__13254_165})
	if callErr != nil {
		return nil, callErr
	}
	arg__13262_171, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{v_153})
	if callErr != nil {
		return nil, callErr
	}
	arg__13267_174, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{v_153})
	if callErr != nil {
		return nil, callErr
	}
	arg__13268_175, callErr = rt.InvokeValue(rt.LookupVar("gogen", "string-lit").Deref(), []vm.Value{arg__13267_174})
	if callErr != nil {
		return nil, callErr
	}
	arg__13269_176, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__13268_175})
	if callErr != nil {
		return nil, callErr
	}
	v177, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "vm-call").Deref(), []vm.Value{vm.String("Symbol"), arg__13269_176})
	if callErr != nil {
		return nil, callErr
	}
	v223 = v177
	v_224 = v_153
	goto b27
b26:
	;
	v182, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "list?").Deref(), []vm.Value{v_154})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v182) {
		v_179 = v_154
		goto b28
	} else {
		v_180 = v_154
		goto b29
	}
b27:
	;
	v226 = v223
	v_227 = v_224
	goto b24
b28:
	;
	v185, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "boxed-list-expr").Deref(), []vm.Value{v_179})
	if callErr != nil {
		return nil, callErr
	}
	v220 = v185
	v_221 = v_179
	goto b30
b29:
	;
	v190, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector?").Deref(), []vm.Value{v_180})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v190) {
		v_187 = v_180
		goto b31
	} else {
		v_188 = v_180
		goto b32
	}
b30:
	;
	v223 = v220
	v_224 = v_221
	goto b27
b31:
	;
	v193, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "boxed-vector-expr").Deref(), []vm.Value{v_187})
	if callErr != nil {
		return nil, callErr
	}
	v217 = v193
	v_218 = v_187
	goto b33
b32:
	;
	v198, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map?").Deref(), []vm.Value{v_188})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v198) {
		v_195 = v_188
		goto b34
	} else {
		v_196 = v_188
		goto b35
	}
b33:
	;
	v220 = v217
	v_221 = v_218
	goto b30
b34:
	;
	v201, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "boxed-map-expr").Deref(), []vm.Value{v_195})
	if callErr != nil {
		return nil, callErr
	}
	v214 = v201
	v_215 = v_195
	goto b36
b35:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		v_203 = v_196
		goto b37
	} else {
		v_204 = v_196
		goto b38
	}
b36:
	;
	v217 = v214
	v_218 = v_215
	goto b33
b37:
	;
	v211 = vm.NIL
	v_212 = v_203
	goto b39
b38:
	;
	v211 = vm.NIL
	v_212 = v_204
	goto b39
b39:
	;
	v214 = v211
	v_215 = v_212
	goto b36
}
func boxed_vector_expr(arg0 vm.Value) (vm.Value, error) {
	var elems_4 vm.Value
	var v12 vm.Value
	var xs_5 vm.Value
	var elems_6 vm.Value
	var arg__13303_19 vm.Value
	var arg__13309_24 vm.Value
	var arg__13311_25 vm.Value
	var arg__13312_26 vm.Value
	var arg__13319_33 vm.Value
	var arg__13325_38 vm.Value
	var arg__13327_39 vm.Value
	var arg__13328_40 vm.Value
	var v41 vm.Value
	var xs_7 vm.Value
	var elems_8 vm.Value
	var v45 vm.Value
	var xs_46 vm.Value
	var elems_47 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = elems_4, v12, xs_5, elems_6, arg__13303_19, arg__13309_24, arg__13311_25, arg__13312_26, arg__13319_33, arg__13325_38, arg__13327_39, arg__13328_40, v41, xs_7, elems_8, v45, xs_46, elems_47
	elems_4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapv").Deref(), []vm.Value{rt.LookupVar("ir.lower-go", "boxed-value-expr").Deref(), arg0})
	if callErr != nil {
		return nil, callErr
	}
	v12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "every?").Deref(), []vm.Value{rt.LookupVar("clojure.core", "identity").Deref(), elems_4})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v12) {
		xs_5 = arg0
		elems_6 = elems_4
		goto b1
	} else {
		xs_7 = arg0
		elems_8 = elems_4
		goto b2
	}
b1:
	;
	arg__13303_19, callErr = rt.InvokeValue(rt.LookupVar("gogen", "type").Deref(), []vm.Value{vm.String("[]vm.Value")})
	if callErr != nil {
		return nil, callErr
	}
	arg__13309_24, callErr = rt.InvokeValue(rt.LookupVar("gogen", "type").Deref(), []vm.Value{vm.String("[]vm.Value")})
	if callErr != nil {
		return nil, callErr
	}
	arg__13311_25, callErr = rt.InvokeValue(rt.LookupVar("gogen", "composite-lit").Deref(), []vm.Value{arg__13309_24, elems_6})
	if callErr != nil {
		return nil, callErr
	}
	arg__13312_26, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__13311_25})
	if callErr != nil {
		return nil, callErr
	}
	arg__13319_33, callErr = rt.InvokeValue(rt.LookupVar("gogen", "type").Deref(), []vm.Value{vm.String("[]vm.Value")})
	if callErr != nil {
		return nil, callErr
	}
	arg__13325_38, callErr = rt.InvokeValue(rt.LookupVar("gogen", "type").Deref(), []vm.Value{vm.String("[]vm.Value")})
	if callErr != nil {
		return nil, callErr
	}
	arg__13327_39, callErr = rt.InvokeValue(rt.LookupVar("gogen", "composite-lit").Deref(), []vm.Value{arg__13325_38, elems_6})
	if callErr != nil {
		return nil, callErr
	}
	arg__13328_40, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__13327_39})
	if callErr != nil {
		return nil, callErr
	}
	v41, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "vm-call").Deref(), []vm.Value{vm.String("NewArrayVector"), arg__13328_40})
	if callErr != nil {
		return nil, callErr
	}
	v45 = v41
	xs_46 = xs_5
	elems_47 = elems_6
	goto b3
b2:
	;
	v45 = vm.NIL
	xs_46 = xs_7
	elems_47 = elems_8
	goto b3
b3:
	;
	return v45, nil
}
func call_assign_stmts(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var refs_4 vm.Value
	var arg__13340_6 vm.Value
	var arg__13347_9 vm.Value
	var callee_10 vm.Value
	var arg__13359_16 vm.Value
	var arg__13372_23 vm.Value
	var arg_exprs_24 vm.Value
	var arg__13376_28 vm.Value
	var arg__13382_33 vm.Value
	var args_slice_34 vm.Value
	var arg__13387_38 vm.Value
	var arg__13393_44 vm.Value
	var arg__13395_46 vm.Value
	var arg__13399_48 vm.Value
	var arg__13404_53 vm.Value
	var arg__13410_59 vm.Value
	var arg__13412_61 vm.Value
	var arg__13416_63 vm.Value
	var invoke_expr_64 vm.Value
	var err_id_68 vm.Value
	var f_69 vm.Value
	var closed_exprs_70 vm.Value
	var nid_71 vm.Value
	var refs_72 vm.Value
	var callee_73 vm.Value
	var arg_exprs_74 vm.Value
	var args_slice_75 vm.Value
	var invoke_expr_76 vm.Value
	var err_id_77 vm.Value
	var v145 vm.Value
	var f_78 vm.Value
	var closed_exprs_79 vm.Value
	var nid_80 vm.Value
	var refs_81 vm.Value
	var callee_82 vm.Value
	var arg_exprs_83 vm.Value
	var args_slice_84 vm.Value
	var invoke_expr_85 vm.Value
	var err_id_86 vm.Value
	var v463 vm.Value
	var f_464 vm.Value
	var closed_exprs_465 vm.Value
	var nid_466 vm.Value
	var refs_467 vm.Value
	var callee_468 vm.Value
	var arg_exprs_469 vm.Value
	var args_slice_470 vm.Value
	var invoke_expr_471 vm.Value
	var err_id_472 vm.Value
	var f_87 vm.Value
	var closed_exprs_88 vm.Value
	var nid_89 vm.Value
	var refs_90 vm.Value
	var and__x_91 vm.Value
	var callee_92 vm.Value
	var arg_exprs_93 vm.Value
	var args_slice_94 vm.Value
	var invoke_expr_95 vm.Value
	var err_id_96 vm.Value
	var v111 vm.Value
	var f_97 vm.Value
	var closed_exprs_98 vm.Value
	var nid_99 vm.Value
	var refs_100 vm.Value
	var and__x_101 vm.Value
	var callee_102 vm.Value
	var arg_exprs_103 vm.Value
	var args_slice_104 vm.Value
	var invoke_expr_105 vm.Value
	var err_id_106 vm.Value
	var v114 vm.Value
	var f_115 vm.Value
	var closed_exprs_116 vm.Value
	var nid_117 vm.Value
	var refs_118 vm.Value
	var and__x_119 vm.Value
	var callee_120 vm.Value
	var arg_exprs_121 vm.Value
	var args_slice_122 vm.Value
	var invoke_expr_123 vm.Value
	var err_id_124 vm.Value
	var f_126 vm.Value
	var closed_exprs_127 vm.Value
	var nid_128 vm.Value
	var refs_129 vm.Value
	var callee_130 vm.Value
	var arg_exprs_131 vm.Value
	var args_slice_132 vm.Value
	var invoke_expr_133 vm.Value
	var err_id_134 vm.Value
	var arg__13437_150 vm.Value
	var arg__13439_151 vm.Value
	var arg__13442_153 vm.Value
	var arg__13451_158 vm.Value
	var arg__13453_159 vm.Value
	var arg__13456_161 vm.Value
	var v162 vm.Value
	var f_135 vm.Value
	var closed_exprs_136 vm.Value
	var nid_137 vm.Value
	var refs_138 vm.Value
	var callee_139 vm.Value
	var arg_exprs_140 vm.Value
	var args_slice_141 vm.Value
	var invoke_expr_142 vm.Value
	var err_id_143 vm.Value
	var arg__13462_169 vm.Value
	var arg__13464_170 vm.Value
	var arg__13467_172 vm.Value
	var arg__13474_179 vm.Value
	var arg__13476_180 vm.Value
	var arg__13479_182 vm.Value
	var v183 vm.Value
	var assign_185 vm.Value
	var f_186 vm.Value
	var closed_exprs_187 vm.Value
	var nid_188 vm.Value
	var refs_189 vm.Value
	var callee_190 vm.Value
	var arg_exprs_191 vm.Value
	var args_slice_192 vm.Value
	var invoke_expr_193 vm.Value
	var err_id_194 vm.Value
	var case__13329_196 vm.Value
	var v220 bool
	var assign_197 vm.Value
	var f_198 vm.Value
	var closed_exprs_199 vm.Value
	var nid_200 vm.Value
	var refs_201 vm.Value
	var callee_202 vm.Value
	var arg_exprs_203 vm.Value
	var args_slice_204 vm.Value
	var invoke_expr_205 vm.Value
	var err_id_206 vm.Value
	var case__13329_207 vm.Value
	var v225 vm.Value
	var assign_208 vm.Value
	var f_209 vm.Value
	var closed_exprs_210 vm.Value
	var nid_211 vm.Value
	var refs_212 vm.Value
	var callee_213 vm.Value
	var arg_exprs_214 vm.Value
	var args_slice_215 vm.Value
	var invoke_expr_216 vm.Value
	var err_id_217 vm.Value
	var case__13329_218 vm.Value
	var v250 bool
	var zero_expr_400 vm.Value
	var assign_401 vm.Value
	var f_402 vm.Value
	var closed_exprs_403 vm.Value
	var nid_404 vm.Value
	var refs_405 vm.Value
	var callee_406 vm.Value
	var arg_exprs_407 vm.Value
	var args_slice_408 vm.Value
	var invoke_expr_409 vm.Value
	var err_id_410 vm.Value
	var case__13329_411 vm.Value
	var arg__13512_417 vm.Value
	var arg__13519_423 vm.Value
	var arg__13520_424 vm.Value
	var arg__13525_427 vm.Value
	var arg__13530_430 vm.Value
	var arg__13531_431 vm.Value
	var arg__13532_432 vm.Value
	var arg__13541_440 vm.Value
	var arg__13548_446 vm.Value
	var arg__13549_447 vm.Value
	var arg__13554_450 vm.Value
	var arg__13559_453 vm.Value
	var arg__13560_454 vm.Value
	var arg__13561_455 vm.Value
	var err_check_457 vm.Value
	var v459 vm.Value
	var assign_227 vm.Value
	var f_228 vm.Value
	var closed_exprs_229 vm.Value
	var nid_230 vm.Value
	var refs_231 vm.Value
	var callee_232 vm.Value
	var arg_exprs_233 vm.Value
	var args_slice_234 vm.Value
	var invoke_expr_235 vm.Value
	var err_id_236 vm.Value
	var case__13329_237 vm.Value
	var v255 vm.Value
	var assign_238 vm.Value
	var f_239 vm.Value
	var closed_exprs_240 vm.Value
	var nid_241 vm.Value
	var refs_242 vm.Value
	var callee_243 vm.Value
	var arg_exprs_244 vm.Value
	var args_slice_245 vm.Value
	var invoke_expr_246 vm.Value
	var err_id_247 vm.Value
	var case__13329_248 vm.Value
	var v280 bool
	var v387 vm.Value
	var assign_388 vm.Value
	var f_389 vm.Value
	var closed_exprs_390 vm.Value
	var nid_391 vm.Value
	var refs_392 vm.Value
	var callee_393 vm.Value
	var arg_exprs_394 vm.Value
	var args_slice_395 vm.Value
	var invoke_expr_396 vm.Value
	var err_id_397 vm.Value
	var case__13329_398 vm.Value
	var assign_257 vm.Value
	var f_258 vm.Value
	var closed_exprs_259 vm.Value
	var nid_260 vm.Value
	var refs_261 vm.Value
	var callee_262 vm.Value
	var arg_exprs_263 vm.Value
	var args_slice_264 vm.Value
	var invoke_expr_265 vm.Value
	var err_id_266 vm.Value
	var case__13329_267 vm.Value
	var v285 vm.Value
	var assign_268 vm.Value
	var f_269 vm.Value
	var closed_exprs_270 vm.Value
	var nid_271 vm.Value
	var refs_272 vm.Value
	var callee_273 vm.Value
	var arg_exprs_274 vm.Value
	var args_slice_275 vm.Value
	var invoke_expr_276 vm.Value
	var err_id_277 vm.Value
	var case__13329_278 vm.Value
	var v310 bool
	var v374 vm.Value
	var assign_375 vm.Value
	var f_376 vm.Value
	var closed_exprs_377 vm.Value
	var nid_378 vm.Value
	var refs_379 vm.Value
	var callee_380 vm.Value
	var arg_exprs_381 vm.Value
	var args_slice_382 vm.Value
	var invoke_expr_383 vm.Value
	var err_id_384 vm.Value
	var case__13329_385 vm.Value
	var assign_287 vm.Value
	var f_288 vm.Value
	var closed_exprs_289 vm.Value
	var nid_290 vm.Value
	var refs_291 vm.Value
	var callee_292 vm.Value
	var arg_exprs_293 vm.Value
	var args_slice_294 vm.Value
	var invoke_expr_295 vm.Value
	var err_id_296 vm.Value
	var case__13329_297 vm.Value
	var v315 vm.Value
	var assign_298 vm.Value
	var f_299 vm.Value
	var closed_exprs_300 vm.Value
	var nid_301 vm.Value
	var refs_302 vm.Value
	var callee_303 vm.Value
	var arg_exprs_304 vm.Value
	var args_slice_305 vm.Value
	var invoke_expr_306 vm.Value
	var err_id_307 vm.Value
	var case__13329_308 vm.Value
	var v361 vm.Value
	var assign_362 vm.Value
	var f_363 vm.Value
	var closed_exprs_364 vm.Value
	var nid_365 vm.Value
	var refs_366 vm.Value
	var callee_367 vm.Value
	var arg_exprs_368 vm.Value
	var args_slice_369 vm.Value
	var invoke_expr_370 vm.Value
	var err_id_371 vm.Value
	var case__13329_372 vm.Value
	var assign_317 vm.Value
	var f_318 vm.Value
	var closed_exprs_319 vm.Value
	var nid_320 vm.Value
	var refs_321 vm.Value
	var callee_322 vm.Value
	var arg_exprs_323 vm.Value
	var args_slice_324 vm.Value
	var invoke_expr_325 vm.Value
	var err_id_326 vm.Value
	var case__13329_327 vm.Value
	var v344 vm.Value
	var assign_328 vm.Value
	var f_329 vm.Value
	var closed_exprs_330 vm.Value
	var nid_331 vm.Value
	var refs_332 vm.Value
	var callee_333 vm.Value
	var arg_exprs_334 vm.Value
	var args_slice_335 vm.Value
	var invoke_expr_336 vm.Value
	var err_id_337 vm.Value
	var case__13329_338 vm.Value
	var v348 vm.Value
	var assign_349 vm.Value
	var f_350 vm.Value
	var closed_exprs_351 vm.Value
	var nid_352 vm.Value
	var refs_353 vm.Value
	var callee_354 vm.Value
	var arg_exprs_355 vm.Value
	var args_slice_356 vm.Value
	var invoke_expr_357 vm.Value
	var err_id_358 vm.Value
	var case__13329_359 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = refs_4, arg__13340_6, arg__13347_9, callee_10, arg__13359_16, arg__13372_23, arg_exprs_24, arg__13376_28, arg__13382_33, args_slice_34, arg__13387_38, arg__13393_44, arg__13395_46, arg__13399_48, arg__13404_53, arg__13410_59, arg__13412_61, arg__13416_63, invoke_expr_64, err_id_68, f_69, closed_exprs_70, nid_71, refs_72, callee_73, arg_exprs_74, args_slice_75, invoke_expr_76, err_id_77, v145, f_78, closed_exprs_79, nid_80, refs_81, callee_82, arg_exprs_83, args_slice_84, invoke_expr_85, err_id_86, v463, f_464, closed_exprs_465, nid_466, refs_467, callee_468, arg_exprs_469, args_slice_470, invoke_expr_471, err_id_472, f_87, closed_exprs_88, nid_89, refs_90, and__x_91, callee_92, arg_exprs_93, args_slice_94, invoke_expr_95, err_id_96, v111, f_97, closed_exprs_98, nid_99, refs_100, and__x_101, callee_102, arg_exprs_103, args_slice_104, invoke_expr_105, err_id_106, v114, f_115, closed_exprs_116, nid_117, refs_118, and__x_119, callee_120, arg_exprs_121, args_slice_122, invoke_expr_123, err_id_124, f_126, closed_exprs_127, nid_128, refs_129, callee_130, arg_exprs_131, args_slice_132, invoke_expr_133, err_id_134, arg__13437_150, arg__13439_151, arg__13442_153, arg__13451_158, arg__13453_159, arg__13456_161, v162, f_135, closed_exprs_136, nid_137, refs_138, callee_139, arg_exprs_140, args_slice_141, invoke_expr_142, err_id_143, arg__13462_169, arg__13464_170, arg__13467_172, arg__13474_179, arg__13476_180, arg__13479_182, v183, assign_185, f_186, closed_exprs_187, nid_188, refs_189, callee_190, arg_exprs_191, args_slice_192, invoke_expr_193, err_id_194, case__13329_196, v220, assign_197, f_198, closed_exprs_199, nid_200, refs_201, callee_202, arg_exprs_203, args_slice_204, invoke_expr_205, err_id_206, case__13329_207, v225, assign_208, f_209, closed_exprs_210, nid_211, refs_212, callee_213, arg_exprs_214, args_slice_215, invoke_expr_216, err_id_217, case__13329_218, v250, zero_expr_400, assign_401, f_402, closed_exprs_403, nid_404, refs_405, callee_406, arg_exprs_407, args_slice_408, invoke_expr_409, err_id_410, case__13329_411, arg__13512_417, arg__13519_423, arg__13520_424, arg__13525_427, arg__13530_430, arg__13531_431, arg__13532_432, arg__13541_440, arg__13548_446, arg__13549_447, arg__13554_450, arg__13559_453, arg__13560_454, arg__13561_455, err_check_457, v459, assign_227, f_228, closed_exprs_229, nid_230, refs_231, callee_232, arg_exprs_233, args_slice_234, invoke_expr_235, err_id_236, case__13329_237, v255, assign_238, f_239, closed_exprs_240, nid_241, refs_242, callee_243, arg_exprs_244, args_slice_245, invoke_expr_246, err_id_247, case__13329_248, v280, v387, assign_388, f_389, closed_exprs_390, nid_391, refs_392, callee_393, arg_exprs_394, args_slice_395, invoke_expr_396, err_id_397, case__13329_398, assign_257, f_258, closed_exprs_259, nid_260, refs_261, callee_262, arg_exprs_263, args_slice_264, invoke_expr_265, err_id_266, case__13329_267, v285, assign_268, f_269, closed_exprs_270, nid_271, refs_272, callee_273, arg_exprs_274, args_slice_275, invoke_expr_276, err_id_277, case__13329_278, v310, v374, assign_375, f_376, closed_exprs_377, nid_378, refs_379, callee_380, arg_exprs_381, args_slice_382, invoke_expr_383, err_id_384, case__13329_385, assign_287, f_288, closed_exprs_289, nid_290, refs_291, callee_292, arg_exprs_293, args_slice_294, invoke_expr_295, err_id_296, case__13329_297, v315, assign_298, f_299, closed_exprs_300, nid_301, refs_302, callee_303, arg_exprs_304, args_slice_305, invoke_expr_306, err_id_307, case__13329_308, v361, assign_362, f_363, closed_exprs_364, nid_365, refs_366, callee_367, arg_exprs_368, args_slice_369, invoke_expr_370, err_id_371, case__13329_372, assign_317, f_318, closed_exprs_319, nid_320, refs_321, callee_322, arg_exprs_323, args_slice_324, invoke_expr_325, err_id_326, case__13329_327, v344, assign_328, f_329, closed_exprs_330, nid_331, refs_332, callee_333, arg_exprs_334, args_slice_335, invoke_expr_336, err_id_337, case__13329_338, v348, assign_349, f_350, closed_exprs_351, nid_352, refs_353, callee_354, arg_exprs_355, args_slice_356, invoke_expr_357, err_id_358, case__13329_359
	refs_4, callErr = rt.InvokeValue(rt.LookupVar("ir", "refs").Deref(), []vm.Value{arg2, arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__13340_6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{refs_4})
	if callErr != nil {
		return nil, callErr
	}
	arg__13347_9, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{refs_4})
	if callErr != nil {
		return nil, callErr
	}
	callee_10, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "box-as-value").Deref(), []vm.Value{arg0, arg1, arg__13347_9})
	if callErr != nil {
		return nil, callErr
	}
	arg__13359_16, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{refs_4})
	if callErr != nil {
		return nil, callErr
	}
	arg__13372_23, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{refs_4})
	if callErr != nil {
		return nil, callErr
	}
	arg_exprs_24, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapv").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v4 vm.Value
		var callErr error
		_ = v4
		v4, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "box-as-value").Deref(), []vm.Value{arg0, arg1, arg0})
		if callErr != nil {
			return nil, callErr
		}
		return v4, nil
	}), arg__13372_23})
	if callErr != nil {
		return nil, callErr
	}
	arg__13376_28, callErr = rt.InvokeValue(rt.LookupVar("gogen", "type").Deref(), []vm.Value{vm.String("[]vm.Value")})
	if callErr != nil {
		return nil, callErr
	}
	arg__13382_33, callErr = rt.InvokeValue(rt.LookupVar("gogen", "type").Deref(), []vm.Value{vm.String("[]vm.Value")})
	if callErr != nil {
		return nil, callErr
	}
	args_slice_34, callErr = rt.InvokeValue(rt.LookupVar("gogen", "composite-lit").Deref(), []vm.Value{arg__13382_33, arg_exprs_24})
	if callErr != nil {
		return nil, callErr
	}
	arg__13387_38, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("rt")})
	if callErr != nil {
		return nil, callErr
	}
	arg__13393_44, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("rt")})
	if callErr != nil {
		return nil, callErr
	}
	arg__13395_46, callErr = rt.InvokeValue(rt.LookupVar("gogen", "field-sel").Deref(), []vm.Value{arg__13393_44, vm.String("InvokeValue")})
	if callErr != nil {
		return nil, callErr
	}
	arg__13399_48, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{callee_10, args_slice_34})
	if callErr != nil {
		return nil, callErr
	}
	arg__13404_53, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("rt")})
	if callErr != nil {
		return nil, callErr
	}
	arg__13410_59, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("rt")})
	if callErr != nil {
		return nil, callErr
	}
	arg__13412_61, callErr = rt.InvokeValue(rt.LookupVar("gogen", "field-sel").Deref(), []vm.Value{arg__13410_59, vm.String("InvokeValue")})
	if callErr != nil {
		return nil, callErr
	}
	arg__13416_63, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{callee_10, args_slice_34})
	if callErr != nil {
		return nil, callErr
	}
	invoke_expr_64, callErr = rt.InvokeValue(rt.LookupVar("gogen", "call").Deref(), []vm.Value{arg__13412_61, arg__13416_63})
	if callErr != nil {
		return nil, callErr
	}
	err_id_68, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("callErr")})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(callee_10) {
		f_87 = arg0
		closed_exprs_88 = arg1
		nid_89 = arg2
		refs_90 = refs_4
		and__x_91 = callee_10
		callee_92 = callee_10
		arg_exprs_93 = arg_exprs_24
		args_slice_94 = args_slice_34
		invoke_expr_95 = invoke_expr_64
		err_id_96 = err_id_68
		goto b4
	} else {
		f_97 = arg0
		closed_exprs_98 = arg1
		nid_99 = arg2
		refs_100 = refs_4
		and__x_101 = callee_10
		callee_102 = callee_10
		arg_exprs_103 = arg_exprs_24
		args_slice_104 = args_slice_34
		invoke_expr_105 = invoke_expr_64
		err_id_106 = err_id_68
		goto b5
	}
b1:
	;
	v145, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "used?").Deref(), []vm.Value{f_69, nid_71})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v145) {
		f_126 = f_69
		closed_exprs_127 = closed_exprs_70
		nid_128 = nid_71
		refs_129 = refs_72
		callee_130 = callee_73
		arg_exprs_131 = arg_exprs_74
		args_slice_132 = args_slice_75
		invoke_expr_133 = invoke_expr_76
		err_id_134 = err_id_77
		goto b7
	} else {
		f_135 = f_69
		closed_exprs_136 = closed_exprs_70
		nid_137 = nid_71
		refs_138 = refs_72
		callee_139 = callee_73
		arg_exprs_140 = arg_exprs_74
		args_slice_141 = args_slice_75
		invoke_expr_142 = invoke_expr_76
		err_id_143 = err_id_77
		goto b8
	}
b2:
	;
	v463 = vm.NIL
	f_464 = f_78
	closed_exprs_465 = closed_exprs_79
	nid_466 = nid_80
	refs_467 = refs_81
	callee_468 = callee_82
	arg_exprs_469 = arg_exprs_83
	args_slice_470 = args_slice_84
	invoke_expr_471 = invoke_expr_85
	err_id_472 = err_id_86
	goto b3
b3:
	;
	return v463, nil
b4:
	;
	v111, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "every?").Deref(), []vm.Value{rt.LookupVar("clojure.core", "identity").Deref(), arg_exprs_93})
	if callErr != nil {
		return nil, callErr
	}
	v114 = v111
	f_115 = f_87
	closed_exprs_116 = closed_exprs_88
	nid_117 = nid_89
	refs_118 = refs_90
	and__x_119 = and__x_91
	callee_120 = callee_92
	arg_exprs_121 = arg_exprs_93
	args_slice_122 = args_slice_94
	invoke_expr_123 = invoke_expr_95
	err_id_124 = err_id_96
	goto b6
b5:
	;
	v114 = and__x_101
	f_115 = f_97
	closed_exprs_116 = closed_exprs_98
	nid_117 = nid_99
	refs_118 = refs_100
	and__x_119 = and__x_101
	callee_120 = callee_102
	arg_exprs_121 = arg_exprs_103
	args_slice_122 = args_slice_104
	invoke_expr_123 = invoke_expr_105
	err_id_124 = err_id_106
	goto b6
b6:
	;
	if vm.IsTruthy(v114) {
		f_69 = f_115
		closed_exprs_70 = closed_exprs_116
		nid_71 = nid_117
		refs_72 = refs_118
		callee_73 = callee_120
		arg_exprs_74 = arg_exprs_121
		args_slice_75 = args_slice_122
		invoke_expr_76 = invoke_expr_123
		err_id_77 = err_id_124
		goto b1
	} else {
		f_78 = f_115
		closed_exprs_79 = closed_exprs_116
		nid_80 = nid_117
		refs_81 = refs_118
		callee_82 = callee_120
		arg_exprs_83 = arg_exprs_121
		args_slice_84 = args_slice_122
		invoke_expr_85 = invoke_expr_123
		err_id_86 = err_id_124
		goto b2
	}
b7:
	;
	arg__13437_150, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "value-ident").Deref(), []vm.Value{f_126, nid_128})
	if callErr != nil {
		return nil, callErr
	}
	arg__13439_151, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__13437_150, err_id_134})
	if callErr != nil {
		return nil, callErr
	}
	arg__13442_153, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{invoke_expr_133})
	if callErr != nil {
		return nil, callErr
	}
	arg__13451_158, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "value-ident").Deref(), []vm.Value{f_126, nid_128})
	if callErr != nil {
		return nil, callErr
	}
	arg__13453_159, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__13451_158, err_id_134})
	if callErr != nil {
		return nil, callErr
	}
	arg__13456_161, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{invoke_expr_133})
	if callErr != nil {
		return nil, callErr
	}
	v162, callErr = rt.InvokeValue(rt.LookupVar("gogen", "multi-assign").Deref(), []vm.Value{vm.String("="), arg__13453_159, arg__13456_161})
	if callErr != nil {
		return nil, callErr
	}
	assign_185 = v162
	f_186 = f_126
	closed_exprs_187 = closed_exprs_127
	nid_188 = nid_128
	refs_189 = refs_129
	callee_190 = callee_130
	arg_exprs_191 = arg_exprs_131
	args_slice_192 = args_slice_132
	invoke_expr_193 = invoke_expr_133
	err_id_194 = err_id_134
	goto b9
b8:
	;
	arg__13462_169, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("_")})
	if callErr != nil {
		return nil, callErr
	}
	arg__13464_170, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__13462_169, err_id_143})
	if callErr != nil {
		return nil, callErr
	}
	arg__13467_172, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{invoke_expr_142})
	if callErr != nil {
		return nil, callErr
	}
	arg__13474_179, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("_")})
	if callErr != nil {
		return nil, callErr
	}
	arg__13476_180, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__13474_179, err_id_143})
	if callErr != nil {
		return nil, callErr
	}
	arg__13479_182, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{invoke_expr_142})
	if callErr != nil {
		return nil, callErr
	}
	v183, callErr = rt.InvokeValue(rt.LookupVar("gogen", "multi-assign").Deref(), []vm.Value{vm.String("="), arg__13476_180, arg__13479_182})
	if callErr != nil {
		return nil, callErr
	}
	assign_185 = v183
	f_186 = f_135
	closed_exprs_187 = closed_exprs_136
	nid_188 = nid_137
	refs_189 = refs_138
	callee_190 = callee_139
	arg_exprs_191 = arg_exprs_140
	args_slice_192 = args_slice_141
	invoke_expr_193 = invoke_expr_142
	err_id_194 = err_id_143
	goto b9
b9:
	;
	case__13329_196, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "function-return-spec").Deref(), []vm.Value{f_186})
	if callErr != nil {
		return nil, callErr
	}
	v220 = case__13329_196 == vm.String("bool")
	if v220 {
		assign_197 = assign_185
		f_198 = f_186
		closed_exprs_199 = closed_exprs_187
		nid_200 = nid_188
		refs_201 = refs_189
		callee_202 = callee_190
		arg_exprs_203 = arg_exprs_191
		args_slice_204 = args_slice_192
		invoke_expr_205 = invoke_expr_193
		err_id_206 = err_id_194
		case__13329_207 = case__13329_196
		goto b10
	} else {
		assign_208 = assign_185
		f_209 = f_186
		closed_exprs_210 = closed_exprs_187
		nid_211 = nid_188
		refs_212 = refs_189
		callee_213 = callee_190
		arg_exprs_214 = arg_exprs_191
		args_slice_215 = args_slice_192
		invoke_expr_216 = invoke_expr_193
		err_id_217 = err_id_194
		case__13329_218 = case__13329_196
		goto b11
	}
b10:
	;
	v225, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("false")})
	if callErr != nil {
		return nil, callErr
	}
	zero_expr_400 = v225
	assign_401 = assign_197
	f_402 = f_198
	closed_exprs_403 = closed_exprs_199
	nid_404 = nid_200
	refs_405 = refs_201
	callee_406 = callee_202
	arg_exprs_407 = arg_exprs_203
	args_slice_408 = args_slice_204
	invoke_expr_409 = invoke_expr_205
	err_id_410 = err_id_206
	case__13329_411 = case__13329_207
	goto b12
b11:
	;
	v250 = case__13329_218 == vm.String("int")
	if v250 {
		assign_227 = assign_208
		f_228 = f_209
		closed_exprs_229 = closed_exprs_210
		nid_230 = nid_211
		refs_231 = refs_212
		callee_232 = callee_213
		arg_exprs_233 = arg_exprs_214
		args_slice_234 = args_slice_215
		invoke_expr_235 = invoke_expr_216
		err_id_236 = err_id_217
		case__13329_237 = case__13329_218
		goto b13
	} else {
		assign_238 = assign_208
		f_239 = f_209
		closed_exprs_240 = closed_exprs_210
		nid_241 = nid_211
		refs_242 = refs_212
		callee_243 = callee_213
		arg_exprs_244 = arg_exprs_214
		args_slice_245 = args_slice_215
		invoke_expr_246 = invoke_expr_216
		err_id_247 = err_id_217
		case__13329_248 = case__13329_218
		goto b14
	}
b12:
	;
	arg__13512_417, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("nil")})
	if callErr != nil {
		return nil, callErr
	}
	arg__13519_423, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("nil")})
	if callErr != nil {
		return nil, callErr
	}
	arg__13520_424, callErr = rt.InvokeValue(rt.LookupVar("gogen", "binary").Deref(), []vm.Value{vm.String("!="), err_id_410, arg__13519_423})
	if callErr != nil {
		return nil, callErr
	}
	arg__13525_427, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{zero_expr_400, err_id_410})
	if callErr != nil {
		return nil, callErr
	}
	arg__13530_430, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{zero_expr_400, err_id_410})
	if callErr != nil {
		return nil, callErr
	}
	arg__13531_431, callErr = rt.InvokeValue(rt.LookupVar("gogen", "return-stmt").Deref(), []vm.Value{arg__13530_430})
	if callErr != nil {
		return nil, callErr
	}
	arg__13532_432, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__13531_431})
	if callErr != nil {
		return nil, callErr
	}
	arg__13541_440, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("nil")})
	if callErr != nil {
		return nil, callErr
	}
	arg__13548_446, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("nil")})
	if callErr != nil {
		return nil, callErr
	}
	arg__13549_447, callErr = rt.InvokeValue(rt.LookupVar("gogen", "binary").Deref(), []vm.Value{vm.String("!="), err_id_410, arg__13548_446})
	if callErr != nil {
		return nil, callErr
	}
	arg__13554_450, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{zero_expr_400, err_id_410})
	if callErr != nil {
		return nil, callErr
	}
	arg__13559_453, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{zero_expr_400, err_id_410})
	if callErr != nil {
		return nil, callErr
	}
	arg__13560_454, callErr = rt.InvokeValue(rt.LookupVar("gogen", "return-stmt").Deref(), []vm.Value{arg__13559_453})
	if callErr != nil {
		return nil, callErr
	}
	arg__13561_455, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__13560_454})
	if callErr != nil {
		return nil, callErr
	}
	err_check_457, callErr = rt.InvokeValue(rt.LookupVar("gogen", "if-stmt").Deref(), []vm.Value{vm.NIL, arg__13549_447, arg__13561_455, vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	v459, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{assign_401, err_check_457})
	if callErr != nil {
		return nil, callErr
	}
	v463 = v459
	f_464 = f_402
	closed_exprs_465 = closed_exprs_403
	nid_466 = nid_404
	refs_467 = refs_405
	callee_468 = callee_406
	arg_exprs_469 = arg_exprs_407
	args_slice_470 = args_slice_408
	invoke_expr_471 = invoke_expr_409
	err_id_472 = err_id_410
	goto b3
b13:
	;
	v255, callErr = rt.InvokeValue(rt.LookupVar("gogen", "int-lit").Deref(), []vm.Value{vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	v387 = v255
	assign_388 = assign_227
	f_389 = f_228
	closed_exprs_390 = closed_exprs_229
	nid_391 = nid_230
	refs_392 = refs_231
	callee_393 = callee_232
	arg_exprs_394 = arg_exprs_233
	args_slice_395 = args_slice_234
	invoke_expr_396 = invoke_expr_235
	err_id_397 = err_id_236
	case__13329_398 = case__13329_237
	goto b15
b14:
	;
	v280 = case__13329_248 == vm.String("float64")
	if v280 {
		assign_257 = assign_238
		f_258 = f_239
		closed_exprs_259 = closed_exprs_240
		nid_260 = nid_241
		refs_261 = refs_242
		callee_262 = callee_243
		arg_exprs_263 = arg_exprs_244
		args_slice_264 = args_slice_245
		invoke_expr_265 = invoke_expr_246
		err_id_266 = err_id_247
		case__13329_267 = case__13329_248
		goto b16
	} else {
		assign_268 = assign_238
		f_269 = f_239
		closed_exprs_270 = closed_exprs_240
		nid_271 = nid_241
		refs_272 = refs_242
		callee_273 = callee_243
		arg_exprs_274 = arg_exprs_244
		args_slice_275 = args_slice_245
		invoke_expr_276 = invoke_expr_246
		err_id_277 = err_id_247
		case__13329_278 = case__13329_248
		goto b17
	}
b15:
	;
	zero_expr_400 = v387
	assign_401 = assign_388
	f_402 = f_389
	closed_exprs_403 = closed_exprs_390
	nid_404 = nid_391
	refs_405 = refs_392
	callee_406 = callee_393
	arg_exprs_407 = arg_exprs_394
	args_slice_408 = args_slice_395
	invoke_expr_409 = invoke_expr_396
	err_id_410 = err_id_397
	case__13329_411 = case__13329_398
	goto b12
b16:
	;
	v285, callErr = rt.InvokeValue(rt.LookupVar("gogen", "float-lit").Deref(), []vm.Value{vm.Float(0)})
	if callErr != nil {
		return nil, callErr
	}
	v374 = v285
	assign_375 = assign_257
	f_376 = f_258
	closed_exprs_377 = closed_exprs_259
	nid_378 = nid_260
	refs_379 = refs_261
	callee_380 = callee_262
	arg_exprs_381 = arg_exprs_263
	args_slice_382 = args_slice_264
	invoke_expr_383 = invoke_expr_265
	err_id_384 = err_id_266
	case__13329_385 = case__13329_267
	goto b18
b17:
	;
	v310 = case__13329_278 == vm.String("string")
	if v310 {
		assign_287 = assign_268
		f_288 = f_269
		closed_exprs_289 = closed_exprs_270
		nid_290 = nid_271
		refs_291 = refs_272
		callee_292 = callee_273
		arg_exprs_293 = arg_exprs_274
		args_slice_294 = args_slice_275
		invoke_expr_295 = invoke_expr_276
		err_id_296 = err_id_277
		case__13329_297 = case__13329_278
		goto b19
	} else {
		assign_298 = assign_268
		f_299 = f_269
		closed_exprs_300 = closed_exprs_270
		nid_301 = nid_271
		refs_302 = refs_272
		callee_303 = callee_273
		arg_exprs_304 = arg_exprs_274
		args_slice_305 = args_slice_275
		invoke_expr_306 = invoke_expr_276
		err_id_307 = err_id_277
		case__13329_308 = case__13329_278
		goto b20
	}
b18:
	;
	v387 = v374
	assign_388 = assign_375
	f_389 = f_376
	closed_exprs_390 = closed_exprs_377
	nid_391 = nid_378
	refs_392 = refs_379
	callee_393 = callee_380
	arg_exprs_394 = arg_exprs_381
	args_slice_395 = args_slice_382
	invoke_expr_396 = invoke_expr_383
	err_id_397 = err_id_384
	case__13329_398 = case__13329_385
	goto b15
b19:
	;
	v315, callErr = rt.InvokeValue(rt.LookupVar("gogen", "string-lit").Deref(), []vm.Value{vm.String("")})
	if callErr != nil {
		return nil, callErr
	}
	v361 = v315
	assign_362 = assign_287
	f_363 = f_288
	closed_exprs_364 = closed_exprs_289
	nid_365 = nid_290
	refs_366 = refs_291
	callee_367 = callee_292
	arg_exprs_368 = arg_exprs_293
	args_slice_369 = args_slice_294
	invoke_expr_370 = invoke_expr_295
	err_id_371 = err_id_296
	case__13329_372 = case__13329_297
	goto b21
b20:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		assign_317 = assign_298
		f_318 = f_299
		closed_exprs_319 = closed_exprs_300
		nid_320 = nid_301
		refs_321 = refs_302
		callee_322 = callee_303
		arg_exprs_323 = arg_exprs_304
		args_slice_324 = args_slice_305
		invoke_expr_325 = invoke_expr_306
		err_id_326 = err_id_307
		case__13329_327 = case__13329_308
		goto b22
	} else {
		assign_328 = assign_298
		f_329 = f_299
		closed_exprs_330 = closed_exprs_300
		nid_331 = nid_301
		refs_332 = refs_302
		callee_333 = callee_303
		arg_exprs_334 = arg_exprs_304
		args_slice_335 = args_slice_305
		invoke_expr_336 = invoke_expr_306
		err_id_337 = err_id_307
		case__13329_338 = case__13329_308
		goto b23
	}
b21:
	;
	v374 = v361
	assign_375 = assign_362
	f_376 = f_363
	closed_exprs_377 = closed_exprs_364
	nid_378 = nid_365
	refs_379 = refs_366
	callee_380 = callee_367
	arg_exprs_381 = arg_exprs_368
	args_slice_382 = args_slice_369
	invoke_expr_383 = invoke_expr_370
	err_id_384 = err_id_371
	case__13329_385 = case__13329_372
	goto b18
b22:
	;
	v344, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("nil")})
	if callErr != nil {
		return nil, callErr
	}
	v348 = v344
	assign_349 = assign_317
	f_350 = f_318
	closed_exprs_351 = closed_exprs_319
	nid_352 = nid_320
	refs_353 = refs_321
	callee_354 = callee_322
	arg_exprs_355 = arg_exprs_323
	args_slice_356 = args_slice_324
	invoke_expr_357 = invoke_expr_325
	err_id_358 = err_id_326
	case__13329_359 = case__13329_327
	goto b24
b23:
	;
	v348 = vm.NIL
	assign_349 = assign_328
	f_350 = f_329
	closed_exprs_351 = closed_exprs_330
	nid_352 = nid_331
	refs_353 = refs_332
	callee_354 = callee_333
	arg_exprs_355 = arg_exprs_334
	args_slice_356 = args_slice_335
	invoke_expr_357 = invoke_expr_336
	err_id_358 = err_id_337
	case__13329_359 = case__13329_338
	goto b24
b24:
	;
	v361 = v348
	assign_362 = assign_349
	f_363 = f_350
	closed_exprs_364 = closed_exprs_351
	nid_365 = nid_352
	refs_366 = refs_353
	callee_367 = callee_354
	arg_exprs_368 = arg_exprs_355
	args_slice_369 = args_slice_356
	invoke_expr_370 = invoke_expr_357
	err_id_371 = err_id_358
	case__13329_372 = case__13329_359
	goto b21
}
func closure_expr(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var info_4 vm.Value
	var f_5 vm.Value
	var closed_exprs_6 vm.Value
	var nid_7 vm.Value
	var info_8 vm.Value
	var arg__13581_19 vm.Value
	var arg__13593_26 vm.Value
	var capture_exprs_27 vm.Value
	var v41 vm.Value
	var f_9 vm.Value
	var closed_exprs_10 vm.Value
	var nid_11 vm.Value
	var info_12 vm.Value
	var v61 vm.Value
	var f_62 vm.Value
	var closed_exprs_63 vm.Value
	var nid_64 vm.Value
	var info_65 vm.Value
	var f_28 vm.Value
	var closed_exprs_29 vm.Value
	var nid_30 vm.Value
	var info_31 vm.Value
	var capture_exprs_32 vm.Value
	var arg__13601_44 vm.Value
	var arg__13606_47 vm.Value
	var v48 vm.Value
	var f_33 vm.Value
	var closed_exprs_34 vm.Value
	var nid_35 vm.Value
	var info_36 vm.Value
	var capture_exprs_37 vm.Value
	var v52 vm.Value
	var f_53 vm.Value
	var closed_exprs_54 vm.Value
	var nid_55 vm.Value
	var info_56 vm.Value
	var capture_exprs_57 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = info_4, f_5, closed_exprs_6, nid_7, info_8, arg__13581_19, arg__13593_26, capture_exprs_27, v41, f_9, closed_exprs_10, nid_11, info_12, v61, f_62, closed_exprs_63, nid_64, info_65, f_28, closed_exprs_29, nid_30, info_31, capture_exprs_32, arg__13601_44, arg__13606_47, v48, f_33, closed_exprs_34, nid_35, info_36, capture_exprs_37, v52, f_53, closed_exprs_54, nid_55, info_56, capture_exprs_57
	info_4, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "closure-info").Deref(), []vm.Value{arg0, arg2})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(info_4) {
		f_5 = arg0
		closed_exprs_6 = arg1
		nid_7 = arg2
		info_8 = info_4
		goto b1
	} else {
		f_9 = arg0
		closed_exprs_10 = arg1
		nid_11 = arg2
		info_12 = info_4
		goto b2
	}
b1:
	;
	arg__13581_19, callErr = rt.InvokeValue(vm.Keyword("captures"), []vm.Value{info_8})
	if callErr != nil {
		return nil, callErr
	}
	arg__13593_26, callErr = rt.InvokeValue(vm.Keyword("captures"), []vm.Value{info_8})
	if callErr != nil {
		return nil, callErr
	}
	capture_exprs_27, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapv").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v4 vm.Value
		var callErr error
		_ = v4
		v4, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "value-expr").Deref(), []vm.Value{f_5, closed_exprs_6, arg0})
		if callErr != nil {
			return nil, callErr
		}
		return v4, nil
	}), arg__13593_26})
	if callErr != nil {
		return nil, callErr
	}
	v41, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "every?").Deref(), []vm.Value{rt.LookupVar("clojure.core", "identity").Deref(), capture_exprs_27})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v41) {
		f_28 = f_5
		closed_exprs_29 = closed_exprs_6
		nid_30 = nid_7
		info_31 = info_8
		capture_exprs_32 = capture_exprs_27
		goto b4
	} else {
		f_33 = f_5
		closed_exprs_34 = closed_exprs_6
		nid_35 = nid_7
		info_36 = info_8
		capture_exprs_37 = capture_exprs_27
		goto b5
	}
b2:
	;
	v61 = vm.NIL
	f_62 = f_9
	closed_exprs_63 = closed_exprs_10
	nid_64 = nid_11
	info_65 = info_12
	goto b3
b3:
	;
	return v61, nil
b4:
	;
	arg__13601_44, callErr = rt.InvokeValue(vm.Keyword("template"), []vm.Value{info_31})
	if callErr != nil {
		return nil, callErr
	}
	arg__13606_47, callErr = rt.InvokeValue(vm.Keyword("template"), []vm.Value{info_31})
	if callErr != nil {
		return nil, callErr
	}
	v48, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "lower-template-closure-expr").Deref(), []vm.Value{arg__13606_47, capture_exprs_32})
	if callErr != nil {
		return nil, callErr
	}
	v52 = v48
	f_53 = f_28
	closed_exprs_54 = closed_exprs_29
	nid_55 = nid_30
	info_56 = info_31
	capture_exprs_57 = capture_exprs_32
	goto b6
b5:
	;
	v52 = vm.NIL
	f_53 = f_33
	closed_exprs_54 = closed_exprs_34
	nid_55 = nid_35
	info_56 = info_36
	capture_exprs_57 = capture_exprs_37
	goto b6
b6:
	;
	v61 = v52
	f_62 = f_53
	closed_exprs_63 = closed_exprs_54
	nid_64 = nid_55
	info_65 = info_56
	goto b3
}
func closure_info(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var op_3 vm.Value
	var aux_5 vm.Value
	var refs_7 vm.Value
	var and__x_19 bool
	var f_8 vm.Value
	var nid_9 vm.Value
	var op_10 vm.Value
	var aux_11 vm.Value
	var refs_12 vm.Value
	var v49 vm.Value
	var f_13 vm.Value
	var nid_14 vm.Value
	var op_15 vm.Value
	var aux_16 vm.Value
	var refs_17 vm.Value
	var v62 bool
	var v239 vm.Value
	var f_240 vm.Value
	var nid_241 vm.Value
	var op_242 vm.Value
	var aux_243 vm.Value
	var refs_244 vm.Value
	var f_20 vm.Value
	var nid_21 vm.Value
	var op_22 vm.Value
	var aux_23 vm.Value
	var refs_24 vm.Value
	var and__x_25 bool
	var v34 vm.Value
	var f_26 vm.Value
	var nid_27 vm.Value
	var op_28 vm.Value
	var aux_29 vm.Value
	var refs_30 vm.Value
	var and__x_31 bool
	var v37 vm.Value
	var f_38 vm.Value
	var nid_39 vm.Value
	var op_40 vm.Value
	var aux_41 vm.Value
	var refs_42 vm.Value
	var and__x_43 vm.Value
	var f_51 vm.Value
	var nid_52 vm.Value
	var op_53 vm.Value
	var aux_54 vm.Value
	var refs_55 vm.Value
	var arg__13639_76 vm.Value
	var v77 bool
	var f_56 vm.Value
	var nid_57 vm.Value
	var op_58 vm.Value
	var aux_59 vm.Value
	var refs_60 vm.Value
	var v110 bool
	var v232 vm.Value
	var f_233 vm.Value
	var nid_234 vm.Value
	var op_235 vm.Value
	var aux_236 vm.Value
	var refs_237 vm.Value
	var f_64 vm.Value
	var nid_65 vm.Value
	var op_66 vm.Value
	var aux_67 vm.Value
	var refs_68 vm.Value
	var arg__13646_82 vm.Value
	var arg__13654_87 vm.Value
	var v88 vm.Value
	var f_69 vm.Value
	var nid_70 vm.Value
	var op_71 vm.Value
	var aux_72 vm.Value
	var refs_73 vm.Value
	var v92 vm.Value
	var f_93 vm.Value
	var nid_94 vm.Value
	var op_95 vm.Value
	var aux_96 vm.Value
	var refs_97 vm.Value
	var f_99 vm.Value
	var nid_100 vm.Value
	var op_101 vm.Value
	var aux_102 vm.Value
	var refs_103 vm.Value
	var arg__13661_124 vm.Value
	var v125 bool
	var f_104 vm.Value
	var nid_105 vm.Value
	var op_106 vm.Value
	var aux_107 vm.Value
	var refs_108 vm.Value
	var v225 vm.Value
	var f_226 vm.Value
	var nid_227 vm.Value
	var op_228 vm.Value
	var aux_229 vm.Value
	var refs_230 vm.Value
	var f_112 vm.Value
	var nid_113 vm.Value
	var op_114 vm.Value
	var aux_115 vm.Value
	var refs_116 vm.Value
	var arg__13668_130 vm.Value
	var arg__13676_135 vm.Value
	var base_136 vm.Value
	var f_117 vm.Value
	var nid_118 vm.Value
	var op_119 vm.Value
	var aux_120 vm.Value
	var refs_121 vm.Value
	var v195 vm.Value
	var f_196 vm.Value
	var nid_197 vm.Value
	var op_198 vm.Value
	var aux_199 vm.Value
	var refs_200 vm.Value
	var f_137 vm.Value
	var nid_138 vm.Value
	var op_139 vm.Value
	var aux_140 vm.Value
	var refs_141 vm.Value
	var base_142 vm.Value
	var arg__13681_152 vm.Value
	var arg__13687_156 vm.Value
	var arg__13691_159 vm.Value
	var arg__13697_163 vm.Value
	var arg__13698_164 vm.Value
	var arg__13704_168 vm.Value
	var arg__13710_172 vm.Value
	var arg__13714_175 vm.Value
	var arg__13720_179 vm.Value
	var arg__13721_180 vm.Value
	var v181 vm.Value
	var f_143 vm.Value
	var nid_144 vm.Value
	var op_145 vm.Value
	var aux_146 vm.Value
	var refs_147 vm.Value
	var base_148 vm.Value
	var v185 vm.Value
	var f_186 vm.Value
	var nid_187 vm.Value
	var op_188 vm.Value
	var aux_189 vm.Value
	var refs_190 vm.Value
	var base_191 vm.Value
	var f_202 vm.Value
	var nid_203 vm.Value
	var op_204 vm.Value
	var aux_205 vm.Value
	var refs_206 vm.Value
	var f_207 vm.Value
	var nid_208 vm.Value
	var op_209 vm.Value
	var aux_210 vm.Value
	var refs_211 vm.Value
	var v218 vm.Value
	var f_219 vm.Value
	var nid_220 vm.Value
	var op_221 vm.Value
	var aux_222 vm.Value
	var refs_223 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = op_3, aux_5, refs_7, and__x_19, f_8, nid_9, op_10, aux_11, refs_12, v49, f_13, nid_14, op_15, aux_16, refs_17, v62, v239, f_240, nid_241, op_242, aux_243, refs_244, f_20, nid_21, op_22, aux_23, refs_24, and__x_25, v34, f_26, nid_27, op_28, aux_29, refs_30, and__x_31, v37, f_38, nid_39, op_40, aux_41, refs_42, and__x_43, f_51, nid_52, op_53, aux_54, refs_55, arg__13639_76, v77, f_56, nid_57, op_58, aux_59, refs_60, v110, v232, f_233, nid_234, op_235, aux_236, refs_237, f_64, nid_65, op_66, aux_67, refs_68, arg__13646_82, arg__13654_87, v88, f_69, nid_70, op_71, aux_72, refs_73, v92, f_93, nid_94, op_95, aux_96, refs_97, f_99, nid_100, op_101, aux_102, refs_103, arg__13661_124, v125, f_104, nid_105, op_106, aux_107, refs_108, v225, f_226, nid_227, op_228, aux_229, refs_230, f_112, nid_113, op_114, aux_115, refs_116, arg__13668_130, arg__13676_135, base_136, f_117, nid_118, op_119, aux_120, refs_121, v195, f_196, nid_197, op_198, aux_199, refs_200, f_137, nid_138, op_139, aux_140, refs_141, base_142, arg__13681_152, arg__13687_156, arg__13691_159, arg__13697_163, arg__13698_164, arg__13704_168, arg__13710_172, arg__13714_175, arg__13720_179, arg__13721_180, v181, f_143, nid_144, op_145, aux_146, refs_147, base_148, v185, f_186, nid_187, op_188, aux_189, refs_190, base_191, f_202, nid_203, op_204, aux_205, refs_206, f_207, nid_208, op_209, aux_210, refs_211, v218, f_219, nid_220, op_221, aux_222, refs_223
	op_3, callErr = rt.InvokeValue(rt.LookupVar("ir", "op").Deref(), []vm.Value{arg1, arg0})
	if callErr != nil {
		return nil, callErr
	}
	aux_5, callErr = rt.InvokeValue(rt.LookupVar("ir", "aux").Deref(), []vm.Value{arg1, arg0})
	if callErr != nil {
		return nil, callErr
	}
	refs_7, callErr = rt.InvokeValue(rt.LookupVar("ir", "refs").Deref(), []vm.Value{arg1, arg0})
	if callErr != nil {
		return nil, callErr
	}
	and__x_19 = op_3 == vm.Keyword("const")
	if and__x_19 {
		f_20 = arg0
		nid_21 = arg1
		op_22 = op_3
		aux_23 = aux_5
		refs_24 = refs_7
		and__x_25 = and__x_19
		goto b4
	} else {
		f_26 = arg0
		nid_27 = arg1
		op_28 = op_3
		aux_29 = aux_5
		refs_30 = refs_7
		and__x_31 = and__x_19
		goto b5
	}
b1:
	;
	v49, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("template"), aux_11, vm.Keyword("captures"), vm.NewArrayVector([]vm.Value{})})
	if callErr != nil {
		return nil, callErr
	}
	v239 = v49
	f_240 = f_8
	nid_241 = nid_9
	op_242 = op_10
	aux_243 = aux_11
	refs_244 = refs_12
	goto b3
b2:
	;
	v62 = op_15 == vm.Keyword("make-closure")
	if v62 {
		f_51 = f_13
		nid_52 = nid_14
		op_53 = op_15
		aux_54 = aux_16
		refs_55 = refs_17
		goto b7
	} else {
		f_56 = f_13
		nid_57 = nid_14
		op_58 = op_15
		aux_59 = aux_16
		refs_60 = refs_17
		goto b8
	}
b3:
	;
	return v239, nil
b4:
	;
	v34, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "any-fn-template?").Deref(), []vm.Value{aux_23})
	if callErr != nil {
		return nil, callErr
	}
	v37 = v34
	f_38 = f_20
	nid_39 = nid_21
	op_40 = op_22
	aux_41 = aux_23
	refs_42 = refs_24
	and__x_43 = vm.Boolean(and__x_25)
	goto b6
b5:
	;
	v37 = vm.Boolean(and__x_31)
	f_38 = f_26
	nid_39 = nid_27
	op_40 = op_28
	aux_41 = aux_29
	refs_42 = refs_30
	and__x_43 = vm.Boolean(and__x_31)
	goto b6
b6:
	;
	if vm.IsTruthy(v37) {
		f_8 = f_38
		nid_9 = nid_39
		op_10 = op_40
		aux_11 = aux_41
		refs_12 = refs_42
		goto b1
	} else {
		f_13 = f_38
		nid_14 = nid_39
		op_15 = op_40
		aux_16 = aux_41
		refs_17 = refs_42
		goto b2
	}
b7:
	;
	arg__13639_76, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{refs_55})
	if callErr != nil {
		return nil, callErr
	}
	v77 = arg__13639_76 == vm.Int(1)
	if v77 {
		f_64 = f_51
		nid_65 = nid_52
		op_66 = op_53
		aux_67 = aux_54
		refs_68 = refs_55
		goto b10
	} else {
		f_69 = f_51
		nid_70 = nid_52
		op_71 = op_53
		aux_72 = aux_54
		refs_73 = refs_55
		goto b11
	}
b8:
	;
	v110 = op_58 == vm.Keyword("push-closed")
	if v110 {
		f_99 = f_56
		nid_100 = nid_57
		op_101 = op_58
		aux_102 = aux_59
		refs_103 = refs_60
		goto b13
	} else {
		f_104 = f_56
		nid_105 = nid_57
		op_106 = op_58
		aux_107 = aux_59
		refs_108 = refs_60
		goto b14
	}
b9:
	;
	v239 = v232
	f_240 = f_233
	nid_241 = nid_234
	op_242 = op_235
	aux_243 = aux_236
	refs_244 = refs_237
	goto b3
b10:
	;
	arg__13646_82, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{refs_68, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	arg__13654_87, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{refs_68, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	v88, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "closure-info").Deref(), []vm.Value{f_64, arg__13654_87})
	if callErr != nil {
		return nil, callErr
	}
	v92 = v88
	f_93 = f_64
	nid_94 = nid_65
	op_95 = op_66
	aux_96 = aux_67
	refs_97 = refs_68
	goto b12
b11:
	;
	v92 = vm.NIL
	f_93 = f_69
	nid_94 = nid_70
	op_95 = op_71
	aux_96 = aux_72
	refs_97 = refs_73
	goto b12
b12:
	;
	v232 = v92
	f_233 = f_93
	nid_234 = nid_94
	op_235 = op_95
	aux_236 = aux_96
	refs_237 = refs_97
	goto b9
b13:
	;
	arg__13661_124, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{refs_103})
	if callErr != nil {
		return nil, callErr
	}
	v125 = arg__13661_124 == vm.Int(2)
	if v125 {
		f_112 = f_99
		nid_113 = nid_100
		op_114 = op_101
		aux_115 = aux_102
		refs_116 = refs_103
		goto b16
	} else {
		f_117 = f_99
		nid_118 = nid_100
		op_119 = op_101
		aux_120 = aux_102
		refs_121 = refs_103
		goto b17
	}
b14:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		f_202 = f_104
		nid_203 = nid_105
		op_204 = op_106
		aux_205 = aux_107
		refs_206 = refs_108
		goto b22
	} else {
		f_207 = f_104
		nid_208 = nid_105
		op_209 = op_106
		aux_210 = aux_107
		refs_211 = refs_108
		goto b23
	}
b15:
	;
	v232 = v225
	f_233 = f_226
	nid_234 = nid_227
	op_235 = op_228
	aux_236 = aux_229
	refs_237 = refs_230
	goto b9
b16:
	;
	arg__13668_130, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{refs_116, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	arg__13676_135, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{refs_116, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	base_136, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "closure-info").Deref(), []vm.Value{f_112, arg__13676_135})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(base_136) {
		f_137 = f_112
		nid_138 = nid_113
		op_139 = op_114
		aux_140 = aux_115
		refs_141 = refs_116
		base_142 = base_136
		goto b19
	} else {
		f_143 = f_112
		nid_144 = nid_113
		op_145 = op_114
		aux_146 = aux_115
		refs_147 = refs_116
		base_148 = base_136
		goto b20
	}
b17:
	;
	v195 = vm.NIL
	f_196 = f_117
	nid_197 = nid_118
	op_198 = op_119
	aux_199 = aux_120
	refs_200 = refs_121
	goto b18
b18:
	;
	v225 = v195
	f_226 = f_196
	nid_227 = nid_197
	op_228 = op_198
	aux_229 = aux_199
	refs_230 = refs_200
	goto b15
b19:
	;
	arg__13681_152, callErr = rt.InvokeValue(vm.Keyword("captures"), []vm.Value{base_142})
	if callErr != nil {
		return nil, callErr
	}
	arg__13687_156, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{refs_141, vm.Int(1)})
	if callErr != nil {
		return nil, callErr
	}
	arg__13691_159, callErr = rt.InvokeValue(vm.Keyword("captures"), []vm.Value{base_142})
	if callErr != nil {
		return nil, callErr
	}
	arg__13697_163, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{refs_141, vm.Int(1)})
	if callErr != nil {
		return nil, callErr
	}
	arg__13698_164, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{arg__13691_159, arg__13697_163})
	if callErr != nil {
		return nil, callErr
	}
	arg__13704_168, callErr = rt.InvokeValue(vm.Keyword("captures"), []vm.Value{base_142})
	if callErr != nil {
		return nil, callErr
	}
	arg__13710_172, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{refs_141, vm.Int(1)})
	if callErr != nil {
		return nil, callErr
	}
	arg__13714_175, callErr = rt.InvokeValue(vm.Keyword("captures"), []vm.Value{base_142})
	if callErr != nil {
		return nil, callErr
	}
	arg__13720_179, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{refs_141, vm.Int(1)})
	if callErr != nil {
		return nil, callErr
	}
	arg__13721_180, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{arg__13714_175, arg__13720_179})
	if callErr != nil {
		return nil, callErr
	}
	v181, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "assoc").Deref(), []vm.Value{base_142, vm.Keyword("captures"), arg__13721_180})
	if callErr != nil {
		return nil, callErr
	}
	v185 = v181
	f_186 = f_137
	nid_187 = nid_138
	op_188 = op_139
	aux_189 = aux_140
	refs_190 = refs_141
	base_191 = base_142
	goto b21
b20:
	;
	v185 = vm.NIL
	f_186 = f_143
	nid_187 = nid_144
	op_188 = op_145
	aux_189 = aux_146
	refs_190 = refs_147
	base_191 = base_148
	goto b21
b21:
	;
	v195 = v185
	f_196 = f_186
	nid_197 = nid_187
	op_198 = op_188
	aux_199 = aux_189
	refs_200 = refs_190
	goto b18
b22:
	;
	v218 = vm.NIL
	f_219 = f_202
	nid_220 = nid_203
	op_221 = op_204
	aux_222 = aux_205
	refs_223 = refs_206
	goto b24
b23:
	;
	v218 = vm.NIL
	f_219 = f_207
	nid_220 = nid_208
	op_221 = op_209
	aux_222 = aux_210
	refs_223 = refs_211
	goto b24
b24:
	;
	v225 = v218
	f_226 = f_219
	nid_227 = nid_220
	op_228 = op_221
	aux_229 = aux_222
	refs_230 = refs_223
	goto b15
}
func collect_local_names(arg0 vm.Value) (vm.Value, error) {
	var arg__13880_2 vm.Value
	var arg__13885_5 vm.Value
	var ids_6 vm.Value
	var remaining_7 vm.Value
	var out_8 vm.Value
	var f_9 vm.Value
	var v21 vm.Value
	var ids_12 vm.Value
	var remaining_13 vm.Value
	var out_14 vm.Value
	var f_15 vm.Value
	var ids_16 vm.Value
	var remaining_17 vm.Value
	var out_18 vm.Value
	var f_19 vm.Value
	var nid_25 vm.Value
	var arg__13897_27 vm.Value
	var arg__13904_30 vm.Value
	var go_type_31 vm.Value
	var v45 vm.Value
	var v66 vm.Value
	var ids_67 vm.Value
	var remaining_68 vm.Value
	var out_69 vm.Value
	var f_70 vm.Value
	var ids_32 vm.Value
	var remaining_33 vm.Value
	var out_34 vm.Value
	var f_35 vm.Value
	var nid_36 vm.Value
	var go_type_37 vm.Value
	var ids_38 vm.Value
	var remaining_39 vm.Value
	var out_40 vm.Value
	var f_41 vm.Value
	var nid_42 vm.Value
	var go_type_43 vm.Value
	var v50 vm.Value
	var arg__13917_52 vm.Value
	var arg__13925_55 vm.Value
	var v56 vm.Value
	var v58 vm.Value
	var ids_59 vm.Value
	var remaining_60 vm.Value
	var out_61 vm.Value
	var f_62 vm.Value
	var nid_63 vm.Value
	var go_type_64 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__13880_2, arg__13885_5, ids_6, remaining_7, out_8, f_9, v21, ids_12, remaining_13, out_14, f_15, ids_16, remaining_17, out_18, f_19, nid_25, arg__13897_27, arg__13904_30, go_type_31, v45, v66, ids_67, remaining_68, out_69, f_70, ids_32, remaining_33, out_34, f_35, nid_36, go_type_37, ids_38, remaining_39, out_40, f_41, nid_42, go_type_43, v50, arg__13917_52, arg__13925_55, v56, v58, ids_59, remaining_60, out_61, f_62, nid_63, go_type_64
	arg__13880_2, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "collect-local-ids").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__13885_5, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "collect-local-ids").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	ids_6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "distinct").Deref(), []vm.Value{arg__13885_5})
	if callErr != nil {
		return nil, callErr
	}
	remaining_7 = ids_6
	out_8 = vm.NewArrayVector([]vm.Value{})
	f_9 = arg0
	goto b1
b1:
	;
	v21, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "empty?").Deref(), []vm.Value{remaining_7})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v21) {
		ids_12 = ids_6
		remaining_13 = remaining_7
		out_14 = out_8
		f_15 = f_9
		goto b2
	} else {
		ids_16 = ids_6
		remaining_17 = remaining_7
		out_18 = out_8
		f_19 = f_9
		goto b3
	}
b2:
	;
	v66 = out_14
	ids_67 = ids_12
	remaining_68 = remaining_13
	out_69 = out_14
	f_70 = f_15
	goto b4
b3:
	;
	nid_25, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{remaining_17})
	if callErr != nil {
		return nil, callErr
	}
	arg__13897_27, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{nid_25, f_19})
	if callErr != nil {
		return nil, callErr
	}
	arg__13904_30, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{nid_25, f_19})
	if callErr != nil {
		return nil, callErr
	}
	go_type_31, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "infer-go-type").Deref(), []vm.Value{arg__13904_30})
	if callErr != nil {
		return nil, callErr
	}
	v45, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{go_type_31})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v45) {
		ids_32 = ids_16
		remaining_33 = remaining_17
		out_34 = out_18
		f_35 = f_19
		nid_36 = nid_25
		go_type_37 = go_type_31
		goto b5
	} else {
		ids_38 = ids_16
		remaining_39 = remaining_17
		out_40 = out_18
		f_41 = f_19
		nid_42 = nid_25
		go_type_43 = go_type_31
		goto b6
	}
b4:
	;
	return v66, nil
b5:
	;
	v58 = vm.NIL
	ids_59 = ids_32
	remaining_60 = remaining_33
	out_61 = out_34
	f_62 = f_35
	nid_63 = nid_36
	go_type_64 = go_type_37
	goto b7
b6:
	;
	v50, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{remaining_39})
	if callErr != nil {
		return nil, callErr
	}
	arg__13917_52, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "value-name").Deref(), []vm.Value{f_41, nid_42})
	if callErr != nil {
		return nil, callErr
	}
	arg__13925_55, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "value-name").Deref(), []vm.Value{f_41, nid_42})
	if callErr != nil {
		return nil, callErr
	}
	v56, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{out_40, arg__13925_55})
	if callErr != nil {
		return nil, callErr
	}
	remaining_7 = v50
	out_8 = v56
	f_9 = f_41
	goto b1
b7:
	;
	v66 = v58
	ids_67 = ids_59
	remaining_68 = remaining_60
	out_69 = out_61
	f_70 = f_62
	goto b4
}
func const_expr(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var v7 vm.Value
	var aux_2 vm.Value
	var t_3 vm.Value
	var v12 vm.Value
	var aux_4 vm.Value
	var t_5 vm.Value
	var v19 bool
	var v148 vm.Value
	var aux_149 vm.Value
	var t_150 vm.Value
	var aux_14 vm.Value
	var t_15 vm.Value
	var v24 vm.Value
	var aux_16 vm.Value
	var t_17 vm.Value
	var v31 bool
	var v144 vm.Value
	var aux_145 vm.Value
	var t_146 vm.Value
	var aux_26 vm.Value
	var t_27 vm.Value
	var v36 vm.Value
	var aux_28 vm.Value
	var t_29 vm.Value
	var v43 vm.Value
	var v140 vm.Value
	var aux_141 vm.Value
	var t_142 vm.Value
	var aux_38 vm.Value
	var t_39 vm.Value
	var v46 vm.Value
	var aux_40 vm.Value
	var t_41 vm.Value
	var v53 vm.Value
	var v136 vm.Value
	var aux_137 vm.Value
	var t_138 vm.Value
	var aux_48 vm.Value
	var t_49 vm.Value
	var aux_50 vm.Value
	var t_51 vm.Value
	var or__x_62 bool
	var v132 vm.Value
	var aux_133 vm.Value
	var t_134 vm.Value
	var aux_57 vm.Value
	var t_58 vm.Value
	var v84 vm.Value
	var aux_59 vm.Value
	var t_60 vm.Value
	var v91 vm.Value
	var v128 vm.Value
	var aux_129 vm.Value
	var t_130 vm.Value
	var aux_63 vm.Value
	var t_64 vm.Value
	var or__x_65 bool
	var aux_66 vm.Value
	var t_67 vm.Value
	var or__x_68 bool
	var arg__13958_75 vm.Value
	var v76 bool
	var v78 bool
	var aux_79 vm.Value
	var t_80 vm.Value
	var or__x_81 vm.Value
	var aux_86 vm.Value
	var t_87 vm.Value
	var v94 vm.Value
	var aux_88 vm.Value
	var t_89 vm.Value
	var v101 vm.Value
	var v124 vm.Value
	var aux_125 vm.Value
	var t_126 vm.Value
	var aux_96 vm.Value
	var t_97 vm.Value
	var v104 vm.Value
	var aux_98 vm.Value
	var t_99 vm.Value
	var v120 vm.Value
	var aux_121 vm.Value
	var t_122 vm.Value
	var aux_106 vm.Value
	var t_107 vm.Value
	var aux_108 vm.Value
	var t_109 vm.Value
	var v116 vm.Value
	var aux_117 vm.Value
	var t_118 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = v7, aux_2, t_3, v12, aux_4, t_5, v19, v148, aux_149, t_150, aux_14, t_15, v24, aux_16, t_17, v31, v144, aux_145, t_146, aux_26, t_27, v36, aux_28, t_29, v43, v140, aux_141, t_142, aux_38, t_39, v46, aux_40, t_41, v53, v136, aux_137, t_138, aux_48, t_49, aux_50, t_51, or__x_62, v132, aux_133, t_134, aux_57, t_58, v84, aux_59, t_60, v91, v128, aux_129, t_130, aux_63, t_64, or__x_65, aux_66, t_67, or__x_68, arg__13958_75, v76, v78, aux_79, t_80, or__x_81, aux_86, t_87, v94, aux_88, t_89, v101, v124, aux_125, t_126, aux_96, t_97, v104, aux_98, t_99, v120, aux_121, t_122, aux_106, t_107, aux_108, t_109, v116, aux_117, t_118
	v7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v7) {
		aux_2 = arg0
		t_3 = arg1
		goto b1
	} else {
		aux_4 = arg0
		t_5 = arg1
		goto b2
	}
b1:
	;
	v12, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("nil")})
	if callErr != nil {
		return nil, callErr
	}
	v148 = v12
	aux_149 = aux_2
	t_150 = t_3
	goto b3
b2:
	;
	v19 = aux_4 == vm.Boolean(true)
	if v19 {
		aux_14 = aux_4
		t_15 = t_5
		goto b4
	} else {
		aux_16 = aux_4
		t_17 = t_5
		goto b5
	}
b3:
	;
	return v148, nil
b4:
	;
	v24, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("true")})
	if callErr != nil {
		return nil, callErr
	}
	v144 = v24
	aux_145 = aux_14
	t_146 = t_15
	goto b6
b5:
	;
	v31 = aux_16 == vm.Boolean(false)
	if v31 {
		aux_26 = aux_16
		t_27 = t_17
		goto b7
	} else {
		aux_28 = aux_16
		t_29 = t_17
		goto b8
	}
b6:
	;
	v148 = v144
	aux_149 = aux_145
	t_150 = t_146
	goto b3
b7:
	;
	v36, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("false")})
	if callErr != nil {
		return nil, callErr
	}
	v140 = v36
	aux_141 = aux_26
	t_142 = t_27
	goto b9
b8:
	;
	v43, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "string?").Deref(), []vm.Value{aux_28})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v43) {
		aux_38 = aux_28
		t_39 = t_29
		goto b10
	} else {
		aux_40 = aux_28
		t_41 = t_29
		goto b11
	}
b9:
	;
	v144 = v140
	aux_145 = aux_141
	t_146 = t_142
	goto b6
b10:
	;
	v46, callErr = rt.InvokeValue(rt.LookupVar("gogen", "string-lit").Deref(), []vm.Value{aux_38})
	if callErr != nil {
		return nil, callErr
	}
	v136 = v46
	aux_137 = aux_38
	t_138 = t_39
	goto b12
b11:
	;
	v53, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "keyword?").Deref(), []vm.Value{aux_40})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v53) {
		aux_48 = aux_40
		t_49 = t_41
		goto b13
	} else {
		aux_50 = aux_40
		t_51 = t_41
		goto b14
	}
b12:
	;
	v140 = v136
	aux_141 = aux_137
	t_142 = t_138
	goto b9
b13:
	;
	v132 = vm.NIL
	aux_133 = aux_48
	t_134 = t_49
	goto b15
b14:
	;
	or__x_62 = t_51 == vm.Keyword("float")
	if or__x_62 {
		aux_63 = aux_50
		t_64 = t_51
		or__x_65 = or__x_62
		goto b19
	} else {
		aux_66 = aux_50
		t_67 = t_51
		or__x_68 = or__x_62
		goto b20
	}
b15:
	;
	v136 = v132
	aux_137 = aux_133
	t_138 = t_134
	goto b12
b16:
	;
	v84, callErr = rt.InvokeValue(rt.LookupVar("gogen", "float-lit").Deref(), []vm.Value{aux_57})
	if callErr != nil {
		return nil, callErr
	}
	v128 = v84
	aux_129 = aux_57
	t_130 = t_58
	goto b18
b17:
	;
	v91, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "integer?").Deref(), []vm.Value{aux_59})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v91) {
		aux_86 = aux_59
		t_87 = t_60
		goto b22
	} else {
		aux_88 = aux_59
		t_89 = t_60
		goto b23
	}
b18:
	;
	v132 = v128
	aux_133 = aux_129
	t_134 = t_130
	goto b15
b19:
	;
	v78 = or__x_65
	aux_79 = aux_63
	t_80 = t_64
	or__x_81 = vm.Boolean(or__x_65)
	goto b21
b20:
	;
	arg__13958_75, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("const"), vm.Keyword("float"), vm.Float(0)})
	if callErr != nil {
		return nil, callErr
	}
	v76 = t_67 == arg__13958_75
	v78 = v76
	aux_79 = aux_66
	t_80 = t_67
	or__x_81 = vm.Boolean(or__x_68)
	goto b21
b21:
	;
	if v78 {
		aux_57 = aux_79
		t_58 = t_80
		goto b16
	} else {
		aux_59 = aux_79
		t_60 = t_80
		goto b17
	}
b22:
	;
	v94, callErr = rt.InvokeValue(rt.LookupVar("gogen", "int-lit").Deref(), []vm.Value{aux_86})
	if callErr != nil {
		return nil, callErr
	}
	v124 = v94
	aux_125 = aux_86
	t_126 = t_87
	goto b24
b23:
	;
	v101, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "float?").Deref(), []vm.Value{aux_88})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v101) {
		aux_96 = aux_88
		t_97 = t_89
		goto b25
	} else {
		aux_98 = aux_88
		t_99 = t_89
		goto b26
	}
b24:
	;
	v128 = v124
	aux_129 = aux_125
	t_130 = t_126
	goto b18
b25:
	;
	v104, callErr = rt.InvokeValue(rt.LookupVar("gogen", "float-lit").Deref(), []vm.Value{aux_96})
	if callErr != nil {
		return nil, callErr
	}
	v120 = v104
	aux_121 = aux_96
	t_122 = t_97
	goto b27
b26:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		aux_106 = aux_98
		t_107 = t_99
		goto b28
	} else {
		aux_108 = aux_98
		t_109 = t_99
		goto b29
	}
b27:
	;
	v124 = v120
	aux_125 = aux_121
	t_126 = t_122
	goto b24
b28:
	;
	v116 = vm.NIL
	aux_117 = aux_106
	t_118 = t_107
	goto b30
b29:
	;
	v116 = vm.NIL
	aux_117 = aux_108
	t_118 = t_109
	goto b30
b30:
	;
	v120 = v116
	aux_121 = aux_117
	t_122 = t_118
	goto b27
}
func discard_locals_stmt(arg0 vm.Value) (vm.Value, error) {
	var v4 vm.Value
	var names_1 vm.Value
	var lhs_9 vm.Value
	var rhs_13 vm.Value
	var v17 vm.Value
	var names_2 vm.Value
	var v21 vm.Value
	var names_22 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _ = v4, names_1, lhs_9, rhs_13, v17, names_2, v21, names_22
	v4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v4) {
		names_1 = arg0
		goto b1
	} else {
		names_2 = arg0
		goto b2
	}
b1:
	;
	lhs_9, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapv").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v4 vm.Value
		var callErr error
		_ = v4
		v4, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("_")})
		if callErr != nil {
			return nil, callErr
		}
		return v4, nil
	}), names_1})
	if callErr != nil {
		return nil, callErr
	}
	rhs_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapv").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v2 vm.Value
		var callErr error
		_ = v2
		v2, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		return v2, nil
	}), names_1})
	if callErr != nil {
		return nil, callErr
	}
	v17, callErr = rt.InvokeValue(rt.LookupVar("gogen", "multi-assign").Deref(), []vm.Value{vm.String("="), lhs_9, rhs_13})
	if callErr != nil {
		return nil, callErr
	}
	v21 = v17
	names_22 = names_1
	goto b3
b2:
	;
	v21 = vm.NIL
	names_22 = names_2
	goto b3
b3:
	;
	return v21, nil
}
func distinct_imports(arg0 vm.Value) (vm.Value, error) {
	var v5 vm.Value
	var remaining_1 vm.Value
	var seen_2 vm.Value
	var out_3 vm.Value
	var v17 vm.Value
	var entries_8 vm.Value
	var remaining_9 vm.Value
	var seen_10 vm.Value
	var out_11 vm.Value
	var entries_12 vm.Value
	var remaining_13 vm.Value
	var seen_14 vm.Value
	var out_15 vm.Value
	var entry_21 vm.Value
	var arg__14016_24 vm.Value
	var arg__14019_26 vm.Value
	var key_27 vm.Value
	var v41 vm.Value
	var v53 vm.Value
	var entries_54 vm.Value
	var remaining_55 vm.Value
	var seen_56 vm.Value
	var out_57 vm.Value
	var entries_28 vm.Value
	var remaining_29 vm.Value
	var seen_30 vm.Value
	var out_31 vm.Value
	var entry_32 vm.Value
	var key_33 vm.Value
	var v44 vm.Value
	var entries_34 vm.Value
	var remaining_35 vm.Value
	var seen_36 vm.Value
	var out_37 vm.Value
	var entry_38 vm.Value
	var key_39 vm.Value
	var v47 vm.Value
	var v49 vm.Value
	var v51 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = v5, remaining_1, seen_2, out_3, v17, entries_8, remaining_9, seen_10, out_11, entries_12, remaining_13, seen_14, out_15, entry_21, arg__14016_24, arg__14019_26, key_27, v41, v53, entries_54, remaining_55, seen_56, out_57, entries_28, remaining_29, seen_30, out_31, entry_32, key_33, v44, entries_34, remaining_35, seen_36, out_37, entry_38, key_39, v47, v49, v51
	v5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	remaining_1 = arg0
	seen_2 = v5
	out_3 = vm.NewArrayVector([]vm.Value{})
	goto b1
b1:
	;
	v17, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "empty?").Deref(), []vm.Value{remaining_1})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v17) {
		entries_8 = arg0
		remaining_9 = remaining_1
		seen_10 = seen_2
		out_11 = out_3
		goto b2
	} else {
		entries_12 = arg0
		remaining_13 = remaining_1
		seen_14 = seen_2
		out_15 = out_3
		goto b3
	}
b2:
	;
	v53 = out_11
	entries_54 = entries_8
	remaining_55 = remaining_9
	seen_56 = seen_10
	out_57 = out_11
	goto b4
b3:
	;
	entry_21, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{remaining_13})
	if callErr != nil {
		return nil, callErr
	}
	arg__14016_24, callErr = rt.InvokeValue(vm.Keyword("path"), []vm.Value{entry_21})
	if callErr != nil {
		return nil, callErr
	}
	arg__14019_26, callErr = rt.InvokeValue(vm.Keyword("alias"), []vm.Value{entry_21})
	if callErr != nil {
		return nil, callErr
	}
	key_27, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__14016_24, arg__14019_26})
	if callErr != nil {
		return nil, callErr
	}
	v41, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "contains?").Deref(), []vm.Value{seen_14, key_27})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v41) {
		entries_28 = entries_12
		remaining_29 = remaining_13
		seen_30 = seen_14
		out_31 = out_15
		entry_32 = entry_21
		key_33 = key_27
		goto b5
	} else {
		entries_34 = entries_12
		remaining_35 = remaining_13
		seen_36 = seen_14
		out_37 = out_15
		entry_38 = entry_21
		key_39 = key_27
		goto b6
	}
b4:
	;
	return v53, nil
b5:
	;
	v44, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{remaining_29})
	if callErr != nil {
		return nil, callErr
	}
	remaining_1 = v44
	seen_2 = seen_30
	out_3 = out_31
	goto b1
b6:
	;
	v47, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{remaining_35})
	if callErr != nil {
		return nil, callErr
	}
	v49, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{seen_36, key_39})
	if callErr != nil {
		return nil, callErr
	}
	v51, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{out_37, entry_38})
	if callErr != nil {
		return nil, callErr
	}
	remaining_1 = v47
	seen_2 = v49
	out_3 = v51
	goto b1
}
func emit_assignments_for_target(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__14044_4 vm.Value
	var arg__14050_7 vm.Value
	var params_8 vm.Value
	var args_10 vm.Value
	var i_11 int
	var out_12 vm.Value
	var params_13 vm.Value
	var f_14 vm.Value
	var closed_exprs_15 vm.Value
	var args_16 vm.Value
	var arg__14059_35 vm.Value
	var v36 bool
	var bt_20 vm.Value
	var i_21 int
	var out_22 vm.Value
	var params_23 vm.Value
	var f_24 vm.Value
	var closed_exprs_25 vm.Value
	var args_26 vm.Value
	var bt_27 vm.Value
	var i_28 int
	var out_29 vm.Value
	var params_30 vm.Value
	var f_31 vm.Value
	var closed_exprs_32 vm.Value
	var args_33 vm.Value
	var pid_40 vm.Value
	var arg__14070_58 vm.Value
	var arg__14077_61 vm.Value
	var v62 vm.Value
	var v401 vm.Value
	var bt_402 vm.Value
	var i_403 int
	var out_404 vm.Value
	var params_405 vm.Value
	var f_406 vm.Value
	var closed_exprs_407 vm.Value
	var args_408 vm.Value
	var bt_41 vm.Value
	var i_42 int
	var out_43 vm.Value
	var params_44 vm.Value
	var f_45 vm.Value
	var closed_exprs_46 vm.Value
	var args_47 vm.Value
	var pid_48 vm.Value
	var v64 int
	var bt_49 vm.Value
	var i_50 int
	var out_51 vm.Value
	var params_52 vm.Value
	var f_53 vm.Value
	var closed_exprs_54 vm.Value
	var args_55 vm.Value
	var pid_56 vm.Value
	var lhs_67 vm.Value
	var arg__14089_69 vm.Value
	var arg__14096_72 vm.Value
	var lhs_spec_73 vm.Value
	var arg_nid_75 vm.Value
	var arg__14107_77 vm.Value
	var arg__14114_80 vm.Value
	var arg_spec_81 vm.Value
	var v107 bool
	var v391 vm.Value
	var bt_392 vm.Value
	var i_393 int
	var out_394 vm.Value
	var params_395 vm.Value
	var f_396 vm.Value
	var closed_exprs_397 vm.Value
	var args_398 vm.Value
	var pid_399 vm.Value
	var bt_82 vm.Value
	var i_83 int
	var out_84 vm.Value
	var params_85 vm.Value
	var f_86 vm.Value
	var closed_exprs_87 vm.Value
	var args_88 vm.Value
	var pid_89 vm.Value
	var lhs_90 vm.Value
	var lhs_spec_91 vm.Value
	var arg_nid_92 vm.Value
	var arg_spec_93 vm.Value
	var v110 vm.Value
	var bt_94 vm.Value
	var i_95 int
	var out_96 vm.Value
	var params_97 vm.Value
	var f_98 vm.Value
	var closed_exprs_99 vm.Value
	var args_100 vm.Value
	var pid_101 vm.Value
	var lhs_102 vm.Value
	var lhs_spec_103 vm.Value
	var arg_nid_104 vm.Value
	var arg_spec_105 vm.Value
	var and__x_137 bool
	var rhs_320 vm.Value
	var bt_321 vm.Value
	var i_322 int
	var out_323 vm.Value
	var params_324 vm.Value
	var f_325 vm.Value
	var closed_exprs_326 vm.Value
	var args_327 vm.Value
	var pid_328 vm.Value
	var lhs_329 vm.Value
	var lhs_spec_330 vm.Value
	var arg_nid_331 vm.Value
	var arg_spec_332 vm.Value
	var v360 vm.Value
	var bt_112 vm.Value
	var i_113 int
	var out_114 vm.Value
	var params_115 vm.Value
	var f_116 vm.Value
	var closed_exprs_117 vm.Value
	var args_118 vm.Value
	var pid_119 vm.Value
	var lhs_120 vm.Value
	var lhs_spec_121 vm.Value
	var arg_nid_122 vm.Value
	var arg_spec_123 vm.Value
	var raw_185 vm.Value
	var bt_124 vm.Value
	var i_125 int
	var out_126 vm.Value
	var params_127 vm.Value
	var f_128 vm.Value
	var closed_exprs_129 vm.Value
	var args_130 vm.Value
	var pid_131 vm.Value
	var lhs_132 vm.Value
	var lhs_spec_133 vm.Value
	var arg_nid_134 vm.Value
	var arg_spec_135 vm.Value
	var v306 vm.Value
	var bt_307 vm.Value
	var i_308 int
	var out_309 vm.Value
	var params_310 vm.Value
	var f_311 vm.Value
	var closed_exprs_312 vm.Value
	var args_313 vm.Value
	var pid_314 vm.Value
	var lhs_315 vm.Value
	var lhs_spec_316 vm.Value
	var arg_nid_317 vm.Value
	var arg_spec_318 vm.Value
	var bt_138 vm.Value
	var i_139 int
	var out_140 vm.Value
	var params_141 vm.Value
	var f_142 vm.Value
	var closed_exprs_143 vm.Value
	var args_144 vm.Value
	var pid_145 vm.Value
	var lhs_146 vm.Value
	var lhs_spec_147 vm.Value
	var arg_nid_148 vm.Value
	var arg_spec_149 vm.Value
	var and__x_150 bool
	var v166 bool
	var bt_151 vm.Value
	var i_152 int
	var out_153 vm.Value
	var params_154 vm.Value
	var f_155 vm.Value
	var closed_exprs_156 vm.Value
	var args_157 vm.Value
	var pid_158 vm.Value
	var lhs_159 vm.Value
	var lhs_spec_160 vm.Value
	var arg_nid_161 vm.Value
	var arg_spec_162 vm.Value
	var and__x_163 bool
	var v169 bool
	var bt_170 vm.Value
	var i_171 int
	var out_172 vm.Value
	var params_173 vm.Value
	var f_174 vm.Value
	var closed_exprs_175 vm.Value
	var args_176 vm.Value
	var pid_177 vm.Value
	var lhs_178 vm.Value
	var lhs_spec_179 vm.Value
	var arg_nid_180 vm.Value
	var arg_spec_181 vm.Value
	var and__x_182 vm.Value
	var bt_186 vm.Value
	var i_187 int
	var out_188 vm.Value
	var params_189 vm.Value
	var f_190 vm.Value
	var closed_exprs_191 vm.Value
	var args_192 vm.Value
	var pid_193 vm.Value
	var lhs_194 vm.Value
	var lhs_spec_195 vm.Value
	var arg_nid_196 vm.Value
	var arg_spec_197 vm.Value
	var raw_198 vm.Value
	var arg__14138_216 vm.Value
	var arg__14144_222 vm.Value
	var arg__14146_224 vm.Value
	var arg__14149_226 vm.Value
	var arg__14154_231 vm.Value
	var arg__14160_237 vm.Value
	var arg__14162_239 vm.Value
	var arg__14165_241 vm.Value
	var v242 vm.Value
	var bt_199 vm.Value
	var i_200 int
	var out_201 vm.Value
	var params_202 vm.Value
	var f_203 vm.Value
	var closed_exprs_204 vm.Value
	var args_205 vm.Value
	var pid_206 vm.Value
	var lhs_207 vm.Value
	var lhs_spec_208 vm.Value
	var arg_nid_209 vm.Value
	var arg_spec_210 vm.Value
	var raw_211 vm.Value
	var v246 vm.Value
	var bt_247 vm.Value
	var i_248 int
	var out_249 vm.Value
	var params_250 vm.Value
	var f_251 vm.Value
	var closed_exprs_252 vm.Value
	var args_253 vm.Value
	var pid_254 vm.Value
	var lhs_255 vm.Value
	var lhs_spec_256 vm.Value
	var arg_nid_257 vm.Value
	var arg_spec_258 vm.Value
	var raw_259 vm.Value
	var bt_261 vm.Value
	var i_262 int
	var out_263 vm.Value
	var params_264 vm.Value
	var f_265 vm.Value
	var closed_exprs_266 vm.Value
	var args_267 vm.Value
	var pid_268 vm.Value
	var lhs_269 vm.Value
	var lhs_spec_270 vm.Value
	var arg_nid_271 vm.Value
	var arg_spec_272 vm.Value
	var v288 vm.Value
	var bt_273 vm.Value
	var i_274 int
	var out_275 vm.Value
	var params_276 vm.Value
	var f_277 vm.Value
	var closed_exprs_278 vm.Value
	var args_279 vm.Value
	var pid_280 vm.Value
	var lhs_281 vm.Value
	var lhs_spec_282 vm.Value
	var arg_nid_283 vm.Value
	var arg_spec_284 vm.Value
	var v292 vm.Value
	var bt_293 vm.Value
	var i_294 int
	var out_295 vm.Value
	var params_296 vm.Value
	var f_297 vm.Value
	var closed_exprs_298 vm.Value
	var args_299 vm.Value
	var pid_300 vm.Value
	var lhs_301 vm.Value
	var lhs_spec_302 vm.Value
	var arg_nid_303 vm.Value
	var arg_spec_304 vm.Value
	var rhs_333 vm.Value
	var bt_334 vm.Value
	var i_335 int
	var out_336 vm.Value
	var params_337 vm.Value
	var f_338 vm.Value
	var closed_exprs_339 vm.Value
	var args_340 vm.Value
	var pid_341 vm.Value
	var lhs_342 vm.Value
	var lhs_spec_343 vm.Value
	var arg_nid_344 vm.Value
	var arg_spec_345 vm.Value
	var rhs_346 vm.Value
	var bt_347 vm.Value
	var i_348 int
	var out_349 vm.Value
	var params_350 vm.Value
	var f_351 vm.Value
	var closed_exprs_352 vm.Value
	var args_353 vm.Value
	var pid_354 vm.Value
	var lhs_355 vm.Value
	var lhs_spec_356 vm.Value
	var arg_nid_357 vm.Value
	var arg_spec_358 vm.Value
	var v364 int
	var arg__14185_368 vm.Value
	var arg__14195_373 vm.Value
	var v374 vm.Value
	var v376 vm.Value
	var rhs_377 vm.Value
	var bt_378 vm.Value
	var i_379 int
	var out_380 vm.Value
	var params_381 vm.Value
	var f_382 vm.Value
	var closed_exprs_383 vm.Value
	var args_384 vm.Value
	var pid_385 vm.Value
	var lhs_386 vm.Value
	var lhs_spec_387 vm.Value
	var arg_nid_388 vm.Value
	var arg_spec_389 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__14044_4, arg__14050_7, params_8, args_10, i_11, out_12, params_13, f_14, closed_exprs_15, args_16, arg__14059_35, v36, bt_20, i_21, out_22, params_23, f_24, closed_exprs_25, args_26, bt_27, i_28, out_29, params_30, f_31, closed_exprs_32, args_33, pid_40, arg__14070_58, arg__14077_61, v62, v401, bt_402, i_403, out_404, params_405, f_406, closed_exprs_407, args_408, bt_41, i_42, out_43, params_44, f_45, closed_exprs_46, args_47, pid_48, v64, bt_49, i_50, out_51, params_52, f_53, closed_exprs_54, args_55, pid_56, lhs_67, arg__14089_69, arg__14096_72, lhs_spec_73, arg_nid_75, arg__14107_77, arg__14114_80, arg_spec_81, v107, v391, bt_392, i_393, out_394, params_395, f_396, closed_exprs_397, args_398, pid_399, bt_82, i_83, out_84, params_85, f_86, closed_exprs_87, args_88, pid_89, lhs_90, lhs_spec_91, arg_nid_92, arg_spec_93, v110, bt_94, i_95, out_96, params_97, f_98, closed_exprs_99, args_100, pid_101, lhs_102, lhs_spec_103, arg_nid_104, arg_spec_105, and__x_137, rhs_320, bt_321, i_322, out_323, params_324, f_325, closed_exprs_326, args_327, pid_328, lhs_329, lhs_spec_330, arg_nid_331, arg_spec_332, v360, bt_112, i_113, out_114, params_115, f_116, closed_exprs_117, args_118, pid_119, lhs_120, lhs_spec_121, arg_nid_122, arg_spec_123, raw_185, bt_124, i_125, out_126, params_127, f_128, closed_exprs_129, args_130, pid_131, lhs_132, lhs_spec_133, arg_nid_134, arg_spec_135, v306, bt_307, i_308, out_309, params_310, f_311, closed_exprs_312, args_313, pid_314, lhs_315, lhs_spec_316, arg_nid_317, arg_spec_318, bt_138, i_139, out_140, params_141, f_142, closed_exprs_143, args_144, pid_145, lhs_146, lhs_spec_147, arg_nid_148, arg_spec_149, and__x_150, v166, bt_151, i_152, out_153, params_154, f_155, closed_exprs_156, args_157, pid_158, lhs_159, lhs_spec_160, arg_nid_161, arg_spec_162, and__x_163, v169, bt_170, i_171, out_172, params_173, f_174, closed_exprs_175, args_176, pid_177, lhs_178, lhs_spec_179, arg_nid_180, arg_spec_181, and__x_182, bt_186, i_187, out_188, params_189, f_190, closed_exprs_191, args_192, pid_193, lhs_194, lhs_spec_195, arg_nid_196, arg_spec_197, raw_198, arg__14138_216, arg__14144_222, arg__14146_224, arg__14149_226, arg__14154_231, arg__14160_237, arg__14162_239, arg__14165_241, v242, bt_199, i_200, out_201, params_202, f_203, closed_exprs_204, args_205, pid_206, lhs_207, lhs_spec_208, arg_nid_209, arg_spec_210, raw_211, v246, bt_247, i_248, out_249, params_250, f_251, closed_exprs_252, args_253, pid_254, lhs_255, lhs_spec_256, arg_nid_257, arg_spec_258, raw_259, bt_261, i_262, out_263, params_264, f_265, closed_exprs_266, args_267, pid_268, lhs_269, lhs_spec_270, arg_nid_271, arg_spec_272, v288, bt_273, i_274, out_275, params_276, f_277, closed_exprs_278, args_279, pid_280, lhs_281, lhs_spec_282, arg_nid_283, arg_spec_284, v292, bt_293, i_294, out_295, params_296, f_297, closed_exprs_298, args_299, pid_300, lhs_301, lhs_spec_302, arg_nid_303, arg_spec_304, rhs_333, bt_334, i_335, out_336, params_337, f_338, closed_exprs_339, args_340, pid_341, lhs_342, lhs_spec_343, arg_nid_344, arg_spec_345, rhs_346, bt_347, i_348, out_349, params_350, f_351, closed_exprs_352, args_353, pid_354, lhs_355, lhs_spec_356, arg_nid_357, arg_spec_358, v364, arg__14185_368, arg__14195_373, v374, v376, rhs_377, bt_378, i_379, out_380, params_381, f_382, closed_exprs_383, args_384, pid_385, lhs_386, lhs_spec_387, arg_nid_388, arg_spec_389
	arg__14044_4, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-target").Deref(), []vm.Value{arg2})
	if callErr != nil {
		return nil, callErr
	}
	arg__14050_7, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-target").Deref(), []vm.Value{arg2})
	if callErr != nil {
		return nil, callErr
	}
	params_8, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-params").Deref(), []vm.Value{arg__14050_7, arg0})
	if callErr != nil {
		return nil, callErr
	}
	args_10, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-args").Deref(), []vm.Value{arg2})
	if callErr != nil {
		return nil, callErr
	}
	i_11 = 0
	out_12 = vm.NewArrayVector([]vm.Value{})
	params_13 = params_8
	f_14 = arg0
	closed_exprs_15 = arg1
	args_16 = args_10
	goto b1
b1:
	;
	arg__14059_35, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{params_13})
	if callErr != nil {
		return nil, callErr
	}
	v36 = rt.GeValue(vm.Int(i_11), arg__14059_35)
	if v36 {
		bt_20 = arg2
		i_21 = i_11
		out_22 = out_12
		params_23 = params_13
		f_24 = f_14
		closed_exprs_25 = closed_exprs_15
		args_26 = args_16
		goto b2
	} else {
		bt_27 = arg2
		i_28 = i_11
		out_29 = out_12
		params_30 = params_13
		f_31 = f_14
		closed_exprs_32 = closed_exprs_15
		args_33 = args_16
		goto b3
	}
b2:
	;
	v401 = out_22
	bt_402 = bt_20
	i_403 = i_21
	out_404 = out_22
	params_405 = params_23
	f_406 = f_24
	closed_exprs_407 = closed_exprs_25
	args_408 = args_26
	goto b4
b3:
	;
	pid_40, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{params_30, vm.Int(i_28)})
	if callErr != nil {
		return nil, callErr
	}
	arg__14070_58, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "used?").Deref(), []vm.Value{f_31, pid_40})
	if callErr != nil {
		return nil, callErr
	}
	arg__14077_61, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "used?").Deref(), []vm.Value{f_31, pid_40})
	if callErr != nil {
		return nil, callErr
	}
	v62, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not").Deref(), []vm.Value{arg__14077_61})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v62) {
		bt_41 = bt_27
		i_42 = i_28
		out_43 = out_29
		params_44 = params_30
		f_45 = f_31
		closed_exprs_46 = closed_exprs_32
		args_47 = args_33
		pid_48 = pid_40
		goto b5
	} else {
		bt_49 = bt_27
		i_50 = i_28
		out_51 = out_29
		params_52 = params_30
		f_53 = f_31
		closed_exprs_54 = closed_exprs_32
		args_55 = args_33
		pid_56 = pid_40
		goto b6
	}
b4:
	;
	return v401, nil
b5:
	;
	v64 = i_42 + 1
	i_11 = v64
	out_12 = out_43
	params_13 = params_44
	f_14 = f_45
	closed_exprs_15 = closed_exprs_46
	args_16 = args_47
	goto b1
b6:
	;
	lhs_67, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "value-ident").Deref(), []vm.Value{f_53, pid_56})
	if callErr != nil {
		return nil, callErr
	}
	arg__14089_69, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{pid_56, f_53})
	if callErr != nil {
		return nil, callErr
	}
	arg__14096_72, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{pid_56, f_53})
	if callErr != nil {
		return nil, callErr
	}
	lhs_spec_73, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "go-type-spec").Deref(), []vm.Value{arg__14096_72})
	if callErr != nil {
		return nil, callErr
	}
	arg_nid_75, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{args_55, vm.Int(i_50)})
	if callErr != nil {
		return nil, callErr
	}
	arg__14107_77, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg_nid_75, f_53})
	if callErr != nil {
		return nil, callErr
	}
	arg__14114_80, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg_nid_75, f_53})
	if callErr != nil {
		return nil, callErr
	}
	arg_spec_81, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "go-type-spec").Deref(), []vm.Value{arg__14114_80})
	if callErr != nil {
		return nil, callErr
	}
	v107 = lhs_spec_73 == vm.String("vm.Value")
	if v107 {
		bt_82 = bt_49
		i_83 = i_50
		out_84 = out_51
		params_85 = params_52
		f_86 = f_53
		closed_exprs_87 = closed_exprs_54
		args_88 = args_55
		pid_89 = pid_56
		lhs_90 = lhs_67
		lhs_spec_91 = lhs_spec_73
		arg_nid_92 = arg_nid_75
		arg_spec_93 = arg_spec_81
		goto b8
	} else {
		bt_94 = bt_49
		i_95 = i_50
		out_96 = out_51
		params_97 = params_52
		f_98 = f_53
		closed_exprs_99 = closed_exprs_54
		args_100 = args_55
		pid_101 = pid_56
		lhs_102 = lhs_67
		lhs_spec_103 = lhs_spec_73
		arg_nid_104 = arg_nid_75
		arg_spec_105 = arg_spec_81
		goto b9
	}
b7:
	;
	v401 = v391
	bt_402 = bt_392
	i_403 = i_393
	out_404 = out_394
	params_405 = params_395
	f_406 = f_396
	closed_exprs_407 = closed_exprs_397
	args_408 = args_398
	goto b4
b8:
	;
	v110, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "box-as-value").Deref(), []vm.Value{f_86, closed_exprs_87, arg_nid_92})
	if callErr != nil {
		return nil, callErr
	}
	rhs_320 = v110
	bt_321 = bt_82
	i_322 = i_83
	out_323 = out_84
	params_324 = params_85
	f_325 = f_86
	closed_exprs_326 = closed_exprs_87
	args_327 = args_88
	pid_328 = pid_89
	lhs_329 = lhs_90
	lhs_spec_330 = lhs_spec_91
	arg_nid_331 = arg_nid_92
	arg_spec_332 = arg_spec_93
	goto b10
b9:
	;
	and__x_137 = lhs_spec_103 == vm.String("bool")
	if and__x_137 {
		bt_138 = bt_94
		i_139 = i_95
		out_140 = out_96
		params_141 = params_97
		f_142 = f_98
		closed_exprs_143 = closed_exprs_99
		args_144 = args_100
		pid_145 = pid_101
		lhs_146 = lhs_102
		lhs_spec_147 = lhs_spec_103
		arg_nid_148 = arg_nid_104
		arg_spec_149 = arg_spec_105
		and__x_150 = and__x_137
		goto b14
	} else {
		bt_151 = bt_94
		i_152 = i_95
		out_153 = out_96
		params_154 = params_97
		f_155 = f_98
		closed_exprs_156 = closed_exprs_99
		args_157 = args_100
		pid_158 = pid_101
		lhs_159 = lhs_102
		lhs_spec_160 = lhs_spec_103
		arg_nid_161 = arg_nid_104
		arg_spec_162 = arg_spec_105
		and__x_163 = and__x_137
		goto b15
	}
b10:
	;
	v360, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{rhs_320})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v360) {
		rhs_333 = rhs_320
		bt_334 = bt_321
		i_335 = i_322
		out_336 = out_323
		params_337 = params_324
		f_338 = f_325
		closed_exprs_339 = closed_exprs_326
		args_340 = args_327
		pid_341 = pid_328
		lhs_342 = lhs_329
		lhs_spec_343 = lhs_spec_330
		arg_nid_344 = arg_nid_331
		arg_spec_345 = arg_spec_332
		goto b23
	} else {
		rhs_346 = rhs_320
		bt_347 = bt_321
		i_348 = i_322
		out_349 = out_323
		params_350 = params_324
		f_351 = f_325
		closed_exprs_352 = closed_exprs_326
		args_353 = args_327
		pid_354 = pid_328
		lhs_355 = lhs_329
		lhs_spec_356 = lhs_spec_330
		arg_nid_357 = arg_nid_331
		arg_spec_358 = arg_spec_332
		goto b24
	}
b11:
	;
	raw_185, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "value-expr").Deref(), []vm.Value{f_116, closed_exprs_117, arg_nid_122})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(raw_185) {
		bt_186 = bt_112
		i_187 = i_113
		out_188 = out_114
		params_189 = params_115
		f_190 = f_116
		closed_exprs_191 = closed_exprs_117
		args_192 = args_118
		pid_193 = pid_119
		lhs_194 = lhs_120
		lhs_spec_195 = lhs_spec_121
		arg_nid_196 = arg_nid_122
		arg_spec_197 = arg_spec_123
		raw_198 = raw_185
		goto b17
	} else {
		bt_199 = bt_112
		i_200 = i_113
		out_201 = out_114
		params_202 = params_115
		f_203 = f_116
		closed_exprs_204 = closed_exprs_117
		args_205 = args_118
		pid_206 = pid_119
		lhs_207 = lhs_120
		lhs_spec_208 = lhs_spec_121
		arg_nid_209 = arg_nid_122
		arg_spec_210 = arg_spec_123
		raw_211 = raw_185
		goto b18
	}
b12:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		bt_261 = bt_124
		i_262 = i_125
		out_263 = out_126
		params_264 = params_127
		f_265 = f_128
		closed_exprs_266 = closed_exprs_129
		args_267 = args_130
		pid_268 = pid_131
		lhs_269 = lhs_132
		lhs_spec_270 = lhs_spec_133
		arg_nid_271 = arg_nid_134
		arg_spec_272 = arg_spec_135
		goto b20
	} else {
		bt_273 = bt_124
		i_274 = i_125
		out_275 = out_126
		params_276 = params_127
		f_277 = f_128
		closed_exprs_278 = closed_exprs_129
		args_279 = args_130
		pid_280 = pid_131
		lhs_281 = lhs_132
		lhs_spec_282 = lhs_spec_133
		arg_nid_283 = arg_nid_134
		arg_spec_284 = arg_spec_135
		goto b21
	}
b13:
	;
	rhs_320 = v306
	bt_321 = bt_307
	i_322 = i_308
	out_323 = out_309
	params_324 = params_310
	f_325 = f_311
	closed_exprs_326 = closed_exprs_312
	args_327 = args_313
	pid_328 = pid_314
	lhs_329 = lhs_315
	lhs_spec_330 = lhs_spec_316
	arg_nid_331 = arg_nid_317
	arg_spec_332 = arg_spec_318
	goto b10
b14:
	;
	v166 = arg_spec_149 == vm.String("vm.Value")
	v169 = v166
	bt_170 = bt_138
	i_171 = i_139
	out_172 = out_140
	params_173 = params_141
	f_174 = f_142
	closed_exprs_175 = closed_exprs_143
	args_176 = args_144
	pid_177 = pid_145
	lhs_178 = lhs_146
	lhs_spec_179 = lhs_spec_147
	arg_nid_180 = arg_nid_148
	arg_spec_181 = arg_spec_149
	and__x_182 = vm.Boolean(and__x_150)
	goto b16
b15:
	;
	v169 = and__x_163
	bt_170 = bt_151
	i_171 = i_152
	out_172 = out_153
	params_173 = params_154
	f_174 = f_155
	closed_exprs_175 = closed_exprs_156
	args_176 = args_157
	pid_177 = pid_158
	lhs_178 = lhs_159
	lhs_spec_179 = lhs_spec_160
	arg_nid_180 = arg_nid_161
	arg_spec_181 = arg_spec_162
	and__x_182 = vm.Boolean(and__x_163)
	goto b16
b16:
	;
	if v169 {
		bt_112 = bt_170
		i_113 = i_171
		out_114 = out_172
		params_115 = params_173
		f_116 = f_174
		closed_exprs_117 = closed_exprs_175
		args_118 = args_176
		pid_119 = pid_177
		lhs_120 = lhs_178
		lhs_spec_121 = lhs_spec_179
		arg_nid_122 = arg_nid_180
		arg_spec_123 = arg_spec_181
		goto b11
	} else {
		bt_124 = bt_170
		i_125 = i_171
		out_126 = out_172
		params_127 = params_173
		f_128 = f_174
		closed_exprs_129 = closed_exprs_175
		args_130 = args_176
		pid_131 = pid_177
		lhs_132 = lhs_178
		lhs_spec_133 = lhs_spec_179
		arg_nid_134 = arg_nid_180
		arg_spec_135 = arg_spec_181
		goto b12
	}
b17:
	;
	arg__14138_216, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("vm")})
	if callErr != nil {
		return nil, callErr
	}
	arg__14144_222, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("vm")})
	if callErr != nil {
		return nil, callErr
	}
	arg__14146_224, callErr = rt.InvokeValue(rt.LookupVar("gogen", "field-sel").Deref(), []vm.Value{arg__14144_222, vm.String("IsTruthy")})
	if callErr != nil {
		return nil, callErr
	}
	arg__14149_226, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{raw_198})
	if callErr != nil {
		return nil, callErr
	}
	arg__14154_231, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("vm")})
	if callErr != nil {
		return nil, callErr
	}
	arg__14160_237, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("vm")})
	if callErr != nil {
		return nil, callErr
	}
	arg__14162_239, callErr = rt.InvokeValue(rt.LookupVar("gogen", "field-sel").Deref(), []vm.Value{arg__14160_237, vm.String("IsTruthy")})
	if callErr != nil {
		return nil, callErr
	}
	arg__14165_241, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{raw_198})
	if callErr != nil {
		return nil, callErr
	}
	v242, callErr = rt.InvokeValue(rt.LookupVar("gogen", "call").Deref(), []vm.Value{arg__14162_239, arg__14165_241})
	if callErr != nil {
		return nil, callErr
	}
	v246 = v242
	bt_247 = bt_186
	i_248 = i_187
	out_249 = out_188
	params_250 = params_189
	f_251 = f_190
	closed_exprs_252 = closed_exprs_191
	args_253 = args_192
	pid_254 = pid_193
	lhs_255 = lhs_194
	lhs_spec_256 = lhs_spec_195
	arg_nid_257 = arg_nid_196
	arg_spec_258 = arg_spec_197
	raw_259 = raw_198
	goto b19
b18:
	;
	v246 = vm.NIL
	bt_247 = bt_199
	i_248 = i_200
	out_249 = out_201
	params_250 = params_202
	f_251 = f_203
	closed_exprs_252 = closed_exprs_204
	args_253 = args_205
	pid_254 = pid_206
	lhs_255 = lhs_207
	lhs_spec_256 = lhs_spec_208
	arg_nid_257 = arg_nid_209
	arg_spec_258 = arg_spec_210
	raw_259 = raw_211
	goto b19
b19:
	;
	v306 = v246
	bt_307 = bt_247
	i_308 = i_248
	out_309 = out_249
	params_310 = params_250
	f_311 = f_251
	closed_exprs_312 = closed_exprs_252
	args_313 = args_253
	pid_314 = pid_254
	lhs_315 = lhs_255
	lhs_spec_316 = lhs_spec_256
	arg_nid_317 = arg_nid_257
	arg_spec_318 = arg_spec_258
	goto b13
b20:
	;
	v288, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "value-expr").Deref(), []vm.Value{f_265, closed_exprs_266, arg_nid_271})
	if callErr != nil {
		return nil, callErr
	}
	v292 = v288
	bt_293 = bt_261
	i_294 = i_262
	out_295 = out_263
	params_296 = params_264
	f_297 = f_265
	closed_exprs_298 = closed_exprs_266
	args_299 = args_267
	pid_300 = pid_268
	lhs_301 = lhs_269
	lhs_spec_302 = lhs_spec_270
	arg_nid_303 = arg_nid_271
	arg_spec_304 = arg_spec_272
	goto b22
b21:
	;
	v292 = vm.NIL
	bt_293 = bt_273
	i_294 = i_274
	out_295 = out_275
	params_296 = params_276
	f_297 = f_277
	closed_exprs_298 = closed_exprs_278
	args_299 = args_279
	pid_300 = pid_280
	lhs_301 = lhs_281
	lhs_spec_302 = lhs_spec_282
	arg_nid_303 = arg_nid_283
	arg_spec_304 = arg_spec_284
	goto b22
b22:
	;
	v306 = v292
	bt_307 = bt_293
	i_308 = i_294
	out_309 = out_295
	params_310 = params_296
	f_311 = f_297
	closed_exprs_312 = closed_exprs_298
	args_313 = args_299
	pid_314 = pid_300
	lhs_315 = lhs_301
	lhs_spec_316 = lhs_spec_302
	arg_nid_317 = arg_nid_303
	arg_spec_318 = arg_spec_304
	goto b13
b23:
	;
	v376 = vm.NIL
	rhs_377 = rhs_333
	bt_378 = bt_334
	i_379 = i_335
	out_380 = out_336
	params_381 = params_337
	f_382 = f_338
	closed_exprs_383 = closed_exprs_339
	args_384 = args_340
	pid_385 = pid_341
	lhs_386 = lhs_342
	lhs_spec_387 = lhs_spec_343
	arg_nid_388 = arg_nid_344
	arg_spec_389 = arg_spec_345
	goto b25
b24:
	;
	v364 = i_348 + 1
	arg__14185_368, callErr = rt.InvokeValue(rt.LookupVar("gogen", "assign").Deref(), []vm.Value{vm.String("="), lhs_355, rhs_346})
	if callErr != nil {
		return nil, callErr
	}
	arg__14195_373, callErr = rt.InvokeValue(rt.LookupVar("gogen", "assign").Deref(), []vm.Value{vm.String("="), lhs_355, rhs_346})
	if callErr != nil {
		return nil, callErr
	}
	v374, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{out_349, arg__14195_373})
	if callErr != nil {
		return nil, callErr
	}
	i_11 = v364
	out_12 = v374
	params_13 = params_350
	f_14 = f_351
	closed_exprs_15 = closed_exprs_352
	args_16 = args_353
	goto b1
b25:
	;
	v391 = v376
	bt_392 = bt_378
	i_393 = i_379
	out_394 = out_380
	params_395 = params_381
	f_396 = f_382
	closed_exprs_397 = closed_exprs_383
	args_398 = args_384
	pid_399 = pid_385
	goto b7
}
func file(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var results_3 vm.Value
	var arg__14208_7 vm.Value
	var arg__14219_12 vm.Value
	var imports_13 vm.Value
	var decls_17 vm.Value
	var import_nodes_21 vm.Value
	var v23 vm.Value
	var callErr error
	_, _, _, _, _, _, _ = results_3, arg__14208_7, arg__14219_12, imports_13, decls_17, import_nodes_21, v23
	results_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__14208_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapcat").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var entries_2 vm.Value
		var result_3 vm.Value
		var entries_4 vm.Value
		var result_5 vm.Value
		var entries_6 vm.Value
		var v11 vm.Value
		var result_12 vm.Value
		var entries_13 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _ = entries_2, result_3, entries_4, result_5, entries_6, v11, result_12, entries_13
		entries_2, callErr = rt.InvokeValue(vm.Keyword("imports"), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		if vm.IsTruthy(entries_2) {
			result_3 = arg0
			entries_4 = entries_2
			goto b1
		} else {
			result_5 = arg0
			entries_6 = entries_2
			goto b2
		}
	b1:
		;
		v11 = entries_4
		result_12 = result_3
		entries_13 = entries_4
		goto b3
	b2:
		;
		v11 = vm.NewArrayVector([]vm.Value{})
		result_12 = result_5
		entries_13 = entries_6
		goto b3
	b3:
		;
		return v11, nil
	}), results_3})
	if callErr != nil {
		return nil, callErr
	}
	arg__14219_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapcat").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var entries_2 vm.Value
		var result_3 vm.Value
		var entries_4 vm.Value
		var result_5 vm.Value
		var entries_6 vm.Value
		var v11 vm.Value
		var result_12 vm.Value
		var entries_13 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _ = entries_2, result_3, entries_4, result_5, entries_6, v11, result_12, entries_13
		entries_2, callErr = rt.InvokeValue(vm.Keyword("imports"), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		if vm.IsTruthy(entries_2) {
			result_3 = arg0
			entries_4 = entries_2
			goto b1
		} else {
			result_5 = arg0
			entries_6 = entries_2
			goto b2
		}
	b1:
		;
		v11 = entries_4
		result_12 = result_3
		entries_13 = entries_4
		goto b3
	b2:
		;
		v11 = vm.NewArrayVector([]vm.Value{})
		result_12 = result_5
		entries_13 = entries_6
		goto b3
	b3:
		;
		return v11, nil
	}), results_3})
	if callErr != nil {
		return nil, callErr
	}
	imports_13, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "distinct-imports").Deref(), []vm.Value{arg__14219_12})
	if callErr != nil {
		return nil, callErr
	}
	decls_17, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapv").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v2 vm.Value
		var callErr error
		_ = v2
		v2, callErr = rt.InvokeValue(vm.Keyword("decl"), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		return v2, nil
	}), results_3})
	if callErr != nil {
		return nil, callErr
	}
	import_nodes_21, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapv").Deref(), []vm.Value{rt.LookupVar("ir.lower-go", "import-spec-node").Deref(), imports_13})
	if callErr != nil {
		return nil, callErr
	}
	v23, callErr = rt.InvokeValue(rt.LookupVar("gogen", "file").Deref(), []vm.Value{arg0, import_nodes_21, decls_17})
	if callErr != nil {
		return nil, callErr
	}
	return v23, nil
}
func fn_template_QMARK_(arg0 vm.Value) (vm.Value, error) {
	var and__x_2 vm.Value
	var aux_3 vm.Value
	var and__x_4 vm.Value
	var arg__14246_9 vm.Value
	var v11 bool
	var aux_5 vm.Value
	var and__x_6 vm.Value
	var v14 vm.Value
	var aux_15 vm.Value
	var and__x_16 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _ = and__x_2, aux_3, and__x_4, arg__14246_9, v11, aux_5, and__x_6, v14, aux_15, and__x_16
	and__x_2, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map?").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(and__x_2) {
		aux_3 = arg0
		and__x_4 = and__x_2
		goto b1
	} else {
		aux_5 = arg0
		and__x_6 = and__x_2
		goto b2
	}
b1:
	;
	arg__14246_9, callErr = rt.InvokeValue(vm.Keyword("kind"), []vm.Value{aux_3})
	if callErr != nil {
		return nil, callErr
	}
	v11 = arg__14246_9 == vm.Keyword("fn-template")
	v14 = vm.Boolean(v11)
	aux_15 = aux_3
	and__x_16 = and__x_4
	goto b3
b2:
	;
	v14 = and__x_6
	aux_15 = aux_5
	and__x_16 = and__x_6
	goto b3
b3:
	;
	return v14, nil
}
func function_needs_tail_call_QMARK_(arg0 vm.Value) (vm.Value, error) {
	var v4 vm.Value
	var blocks_1 vm.Value
	var f_2 vm.Value
	var v11 vm.Value
	var blocks_6 vm.Value
	var f_7 vm.Value
	var blocks_8 vm.Value
	var f_9 vm.Value
	var v86 vm.Value
	var blocks_87 vm.Value
	var f_88 vm.Value
	var blocks_15 vm.Value
	var f_16 vm.Value
	var arg__14639_22 vm.Value
	var arg__14645_25 vm.Value
	var term_26 vm.Value
	var v34 vm.Value
	var blocks_17 vm.Value
	var f_18 vm.Value
	var v82 vm.Value
	var blocks_83 vm.Value
	var f_84 vm.Value
	var blocks_27 vm.Value
	var f_28 vm.Value
	var term_29 vm.Value
	var v37 vm.Value
	var blocks_30 vm.Value
	var f_31 vm.Value
	var term_32 vm.Value
	var arg__14659_47 vm.Value
	var v48 bool
	var v75 vm.Value
	var blocks_76 vm.Value
	var f_77 vm.Value
	var term_78 vm.Value
	var blocks_39 vm.Value
	var f_40 vm.Value
	var term_41 vm.Value
	var blocks_42 vm.Value
	var f_43 vm.Value
	var term_44 vm.Value
	var v70 vm.Value
	var blocks_71 vm.Value
	var f_72 vm.Value
	var term_73 vm.Value
	var blocks_52 vm.Value
	var f_53 vm.Value
	var term_54 vm.Value
	var v61 vm.Value
	var blocks_55 vm.Value
	var f_56 vm.Value
	var term_57 vm.Value
	var v65 vm.Value
	var blocks_66 vm.Value
	var f_67 vm.Value
	var term_68 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = v4, blocks_1, f_2, v11, blocks_6, f_7, blocks_8, f_9, v86, blocks_87, f_88, blocks_15, f_16, arg__14639_22, arg__14645_25, term_26, v34, blocks_17, f_18, v82, blocks_83, f_84, blocks_27, f_28, term_29, v37, blocks_30, f_31, term_32, arg__14659_47, v48, v75, blocks_76, f_77, term_78, blocks_39, f_40, term_41, blocks_42, f_43, term_44, v70, blocks_71, f_72, term_73, blocks_52, f_53, term_54, v61, blocks_55, f_56, term_57, v65, blocks_66, f_67, term_68
	v4, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	blocks_1 = v4
	f_2 = arg0
	goto b1
b1:
	;
	v11, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "empty?").Deref(), []vm.Value{blocks_1})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v11) {
		blocks_6 = blocks_1
		f_7 = f_2
		goto b2
	} else {
		blocks_8 = blocks_1
		f_9 = f_2
		goto b3
	}
b2:
	;
	v86 = vm.Boolean(false)
	blocks_87 = blocks_6
	f_88 = f_7
	goto b4
b3:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		blocks_15 = blocks_8
		f_16 = f_9
		goto b5
	} else {
		blocks_17 = blocks_8
		f_18 = f_9
		goto b6
	}
b4:
	;
	return v86, nil
b5:
	;
	arg__14639_22, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{blocks_15})
	if callErr != nil {
		return nil, callErr
	}
	arg__14645_25, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{blocks_15})
	if callErr != nil {
		return nil, callErr
	}
	term_26, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-term").Deref(), []vm.Value{arg__14645_25, f_16})
	if callErr != nil {
		return nil, callErr
	}
	v34, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{term_26})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v34) {
		blocks_27 = blocks_15
		f_28 = f_16
		term_29 = term_26
		goto b8
	} else {
		blocks_30 = blocks_15
		f_31 = f_16
		term_32 = term_26
		goto b9
	}
b6:
	;
	v82 = vm.NIL
	blocks_83 = blocks_17
	f_84 = f_18
	goto b7
b7:
	;
	v86 = v82
	blocks_87 = blocks_83
	f_88 = f_84
	goto b4
b8:
	;
	v37, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{blocks_27})
	if callErr != nil {
		return nil, callErr
	}
	blocks_1 = v37
	f_2 = f_28
	goto b1
b9:
	;
	arg__14659_47, callErr = rt.InvokeValue(rt.LookupVar("ir", "op").Deref(), []vm.Value{term_32, f_31})
	if callErr != nil {
		return nil, callErr
	}
	v48 = arg__14659_47 == vm.Keyword("tail-call")
	if v48 {
		blocks_39 = blocks_30
		f_40 = f_31
		term_41 = term_32
		goto b11
	} else {
		blocks_42 = blocks_30
		f_43 = f_31
		term_44 = term_32
		goto b12
	}
b10:
	;
	v82 = v75
	blocks_83 = blocks_76
	f_84 = f_77
	goto b7
b11:
	;
	v70 = vm.Boolean(true)
	blocks_71 = blocks_39
	f_72 = f_40
	term_73 = term_41
	goto b13
b12:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		blocks_52 = blocks_42
		f_53 = f_43
		term_54 = term_44
		goto b14
	} else {
		blocks_55 = blocks_42
		f_56 = f_43
		term_57 = term_44
		goto b15
	}
b13:
	;
	v75 = v70
	blocks_76 = blocks_71
	f_77 = f_72
	term_78 = term_73
	goto b10
b14:
	;
	v61, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{blocks_52})
	if callErr != nil {
		return nil, callErr
	}
	blocks_1 = v61
	f_2 = f_53
	goto b1
b15:
	;
	v65 = vm.NIL
	blocks_66 = blocks_55
	f_67 = f_56
	term_68 = term_57
	goto b16
b16:
	;
	v70 = v65
	blocks_71 = blocks_66
	f_72 = f_67
	term_73 = term_68
	goto b13
}
func function_needs_vm_QMARK_(arg0 vm.Value) (vm.Value, error) {
	var or__x_2 vm.Value
	var f_3 vm.Value
	var or__x_4 vm.Value
	var f_5 vm.Value
	var or__x_6 vm.Value
	var arg_types_11 vm.Value
	var f_12 vm.Value
	var or__x_13 vm.Value
	var or__x_15 vm.Value
	var arg_types_16 vm.Value
	var f_17 vm.Value
	var or__x_18 vm.Value
	var arg_types_19 vm.Value
	var f_20 vm.Value
	var or__x_21 vm.Value
	var ret_ids_26 vm.Value
	var arg_types_27 vm.Value
	var f_28 vm.Value
	var or__x_29 vm.Value
	var or__x_31 vm.Value
	var ret_ids_32 vm.Value
	var arg_types_33 vm.Value
	var f_34 vm.Value
	var or__x_35 vm.Value
	var ret_ids_36 vm.Value
	var arg_types_37 vm.Value
	var f_38 vm.Value
	var or__x_39 vm.Value
	var or__x_43 vm.Value
	var v113 vm.Value
	var ret_ids_114 vm.Value
	var arg_types_115 vm.Value
	var f_116 vm.Value
	var or__x_117 vm.Value
	var ret_ids_44 vm.Value
	var arg_types_45 vm.Value
	var f_46 vm.Value
	var or__x_47 vm.Value
	var ret_ids_48 vm.Value
	var arg_types_49 vm.Value
	var f_50 vm.Value
	var or__x_51 vm.Value
	var or__x_57 vm.Value
	var v107 vm.Value
	var ret_ids_108 vm.Value
	var arg_types_109 vm.Value
	var f_110 vm.Value
	var or__x_111 vm.Value
	var ret_ids_58 vm.Value
	var arg_types_59 vm.Value
	var f_60 vm.Value
	var or__x_61 vm.Value
	var ret_ids_62 vm.Value
	var arg_types_63 vm.Value
	var f_64 vm.Value
	var or__x_65 vm.Value
	var or__x_75 vm.Value
	var v101 vm.Value
	var ret_ids_102 vm.Value
	var arg_types_103 vm.Value
	var f_104 vm.Value
	var or__x_105 vm.Value
	var ret_ids_76 vm.Value
	var arg_types_77 vm.Value
	var f_78 vm.Value
	var or__x_79 vm.Value
	var ret_ids_80 vm.Value
	var arg_types_81 vm.Value
	var f_82 vm.Value
	var or__x_83 vm.Value
	var arg__14729_88 vm.Value
	var arg__14735_92 vm.Value
	var v93 vm.Value
	var v95 vm.Value
	var ret_ids_96 vm.Value
	var arg_types_97 vm.Value
	var f_98 vm.Value
	var or__x_99 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = or__x_2, f_3, or__x_4, f_5, or__x_6, arg_types_11, f_12, or__x_13, or__x_15, arg_types_16, f_17, or__x_18, arg_types_19, f_20, or__x_21, ret_ids_26, arg_types_27, f_28, or__x_29, or__x_31, ret_ids_32, arg_types_33, f_34, or__x_35, ret_ids_36, arg_types_37, f_38, or__x_39, or__x_43, v113, ret_ids_114, arg_types_115, f_116, or__x_117, ret_ids_44, arg_types_45, f_46, or__x_47, ret_ids_48, arg_types_49, f_50, or__x_51, or__x_57, v107, ret_ids_108, arg_types_109, f_110, or__x_111, ret_ids_58, arg_types_59, f_60, or__x_61, ret_ids_62, arg_types_63, f_64, or__x_65, or__x_75, v101, ret_ids_102, arg_types_103, f_104, or__x_105, ret_ids_76, arg_types_77, f_78, or__x_79, ret_ids_80, arg_types_81, f_82, or__x_83, arg__14729_88, arg__14735_92, v93, v95, ret_ids_96, arg_types_97, f_98, or__x_99
	or__x_2, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-arg-types").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(or__x_2) {
		f_3 = arg0
		or__x_4 = or__x_2
		goto b1
	} else {
		f_5 = arg0
		or__x_6 = or__x_2
		goto b2
	}
b1:
	;
	arg_types_11 = or__x_4
	f_12 = f_3
	or__x_13 = or__x_4
	goto b3
b2:
	;
	arg_types_11 = vm.NewArrayVector([]vm.Value{})
	f_12 = f_5
	or__x_13 = or__x_6
	goto b3
b3:
	;
	or__x_15, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "return-ref-ids").Deref(), []vm.Value{f_12})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(or__x_15) {
		arg_types_16 = arg_types_11
		f_17 = f_12
		or__x_18 = or__x_15
		goto b4
	} else {
		arg_types_19 = arg_types_11
		f_20 = f_12
		or__x_21 = or__x_15
		goto b5
	}
b4:
	;
	ret_ids_26 = or__x_18
	arg_types_27 = arg_types_16
	f_28 = f_17
	or__x_29 = or__x_18
	goto b6
b5:
	;
	ret_ids_26 = vm.NewArrayVector([]vm.Value{})
	arg_types_27 = arg_types_19
	f_28 = f_20
	or__x_29 = or__x_21
	goto b6
b6:
	;
	or__x_31, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-variadic?").Deref(), []vm.Value{f_28})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(or__x_31) {
		ret_ids_32 = ret_ids_26
		arg_types_33 = arg_types_27
		f_34 = f_28
		or__x_35 = or__x_31
		goto b7
	} else {
		ret_ids_36 = ret_ids_26
		arg_types_37 = arg_types_27
		f_38 = f_28
		or__x_39 = or__x_31
		goto b8
	}
b7:
	;
	v113 = or__x_35
	ret_ids_114 = ret_ids_32
	arg_types_115 = arg_types_33
	f_116 = f_34
	or__x_117 = or__x_35
	goto b9
b8:
	;
	or__x_43, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "function-needs-rt?").Deref(), []vm.Value{f_38})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(or__x_43) {
		ret_ids_44 = ret_ids_36
		arg_types_45 = arg_types_37
		f_46 = f_38
		or__x_47 = or__x_43
		goto b10
	} else {
		ret_ids_48 = ret_ids_36
		arg_types_49 = arg_types_37
		f_50 = f_38
		or__x_51 = or__x_43
		goto b11
	}
b9:
	;
	return v113, nil
b10:
	;
	v107 = or__x_47
	ret_ids_108 = ret_ids_44
	arg_types_109 = arg_types_45
	f_110 = f_46
	or__x_111 = or__x_47
	goto b12
b11:
	;
	or__x_57, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "some").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var arg__14687_3 vm.Value
		var v4 vm.Value
		var callErr error
		_, _ = arg__14687_3, v4
		arg__14687_3, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "go-type-spec").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		v4 = vm.Boolean(vm.String("vm.Value") == arg__14687_3)
		return v4, nil
	}), arg_types_49})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(or__x_57) {
		ret_ids_58 = ret_ids_48
		arg_types_59 = arg_types_49
		f_60 = f_50
		or__x_61 = or__x_57
		goto b13
	} else {
		ret_ids_62 = ret_ids_48
		arg_types_63 = arg_types_49
		f_64 = f_50
		or__x_65 = or__x_57
		goto b14
	}
b12:
	;
	v113 = v107
	ret_ids_114 = ret_ids_108
	arg_types_115 = arg_types_109
	f_116 = f_110
	or__x_117 = or__x_39
	goto b9
b13:
	;
	v101 = or__x_61
	ret_ids_102 = ret_ids_58
	arg_types_103 = arg_types_59
	f_104 = f_60
	or__x_105 = or__x_61
	goto b15
b14:
	;
	or__x_75, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "some").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var arg__14714_4 vm.Value
		var arg__14721_7 vm.Value
		var arg__14722_8 vm.Value
		var v9 vm.Value
		var callErr error
		_, _, _, _ = arg__14714_4, arg__14721_7, arg__14722_8, v9
		arg__14714_4, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_64})
		if callErr != nil {
			return nil, callErr
		}
		arg__14721_7, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_64})
		if callErr != nil {
			return nil, callErr
		}
		arg__14722_8, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "go-type-spec").Deref(), []vm.Value{arg__14721_7})
		if callErr != nil {
			return nil, callErr
		}
		v9 = vm.Boolean(vm.String("vm.Value") == arg__14722_8)
		return v9, nil
	}), ret_ids_62})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(or__x_75) {
		ret_ids_76 = ret_ids_62
		arg_types_77 = arg_types_63
		f_78 = f_64
		or__x_79 = or__x_75
		goto b16
	} else {
		ret_ids_80 = ret_ids_62
		arg_types_81 = arg_types_63
		f_82 = f_64
		or__x_83 = or__x_75
		goto b17
	}
b15:
	;
	v107 = v101
	ret_ids_108 = ret_ids_102
	arg_types_109 = arg_types_103
	f_110 = f_104
	or__x_111 = or__x_51
	goto b12
b16:
	;
	v95 = or__x_79
	ret_ids_96 = ret_ids_76
	arg_types_97 = arg_types_77
	f_98 = f_78
	or__x_99 = or__x_79
	goto b18
b17:
	;
	arg__14729_88, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "nested-template-fns").Deref(), []vm.Value{f_82})
	if callErr != nil {
		return nil, callErr
	}
	arg__14735_92, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "nested-template-fns").Deref(), []vm.Value{f_82})
	if callErr != nil {
		return nil, callErr
	}
	v93, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "some").Deref(), []vm.Value{rt.LookupVar("ir.lower-go", "function-needs-vm?").Deref(), arg__14735_92})
	if callErr != nil {
		return nil, callErr
	}
	v95 = v93
	ret_ids_96 = ret_ids_80
	arg_types_97 = arg_types_81
	f_98 = f_82
	or__x_99 = or__x_83
	goto b18
b18:
	;
	v101 = v95
	ret_ids_102 = ret_ids_96
	arg_types_103 = arg_types_97
	f_104 = f_98
	or__x_105 = or__x_65
	goto b15
}
func function_return_spec(arg0 vm.Value) (vm.Value, error) {
	var ret_ids_2 vm.Value
	var f_3 vm.Value
	var ret_ids_4 vm.Value
	var specs_34 vm.Value
	var and__x_44 vm.Value
	var f_5 vm.Value
	var ret_ids_6 vm.Value
	var v82 vm.Value
	var f_83 vm.Value
	var ret_ids_84 vm.Value
	var f_7 vm.Value
	var and__x_8 vm.Value
	var ret_ids_9 vm.Value
	var arg__14742_15 vm.Value
	var arg__14747_18 vm.Value
	var v19 vm.Value
	var f_10 vm.Value
	var and__x_11 vm.Value
	var ret_ids_12 vm.Value
	var v22 vm.Value
	var f_23 vm.Value
	var and__x_24 vm.Value
	var ret_ids_25 vm.Value
	var f_35 vm.Value
	var ret_ids_36 vm.Value
	var specs_37 vm.Value
	var v71 vm.Value
	var f_38 vm.Value
	var ret_ids_39 vm.Value
	var specs_40 vm.Value
	var v75 vm.Value
	var f_76 vm.Value
	var ret_ids_77 vm.Value
	var specs_78 vm.Value
	var f_45 vm.Value
	var ret_ids_46 vm.Value
	var specs_47 vm.Value
	var and__x_48 vm.Value
	var arg__14788_56 vm.Value
	var arg__14793_59 vm.Value
	var arg__14794_60 vm.Value
	var v61 bool
	var f_49 vm.Value
	var ret_ids_50 vm.Value
	var specs_51 vm.Value
	var and__x_52 vm.Value
	var v64 vm.Value
	var f_65 vm.Value
	var ret_ids_66 vm.Value
	var specs_67 vm.Value
	var and__x_68 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = ret_ids_2, f_3, ret_ids_4, specs_34, and__x_44, f_5, ret_ids_6, v82, f_83, ret_ids_84, f_7, and__x_8, ret_ids_9, arg__14742_15, arg__14747_18, v19, f_10, and__x_11, ret_ids_12, v22, f_23, and__x_24, ret_ids_25, f_35, ret_ids_36, specs_37, v71, f_38, ret_ids_39, specs_40, v75, f_76, ret_ids_77, specs_78, f_45, ret_ids_46, specs_47, and__x_48, arg__14788_56, arg__14793_59, arg__14794_60, v61, f_49, ret_ids_50, specs_51, and__x_52, v64, f_65, ret_ids_66, specs_67, and__x_68
	ret_ids_2, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "return-ref-ids").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(ret_ids_2) {
		f_7 = arg0
		and__x_8 = ret_ids_2
		ret_ids_9 = ret_ids_2
		goto b4
	} else {
		f_10 = arg0
		and__x_11 = ret_ids_2
		ret_ids_12 = ret_ids_2
		goto b5
	}
b1:
	;
	specs_34, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var arg__14769_3 vm.Value
		var arg__14776_6 vm.Value
		var v7 vm.Value
		var callErr error
		_, _, _ = arg__14769_3, arg__14776_6, v7
		arg__14769_3, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_3})
		if callErr != nil {
			return nil, callErr
		}
		arg__14776_6, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_3})
		if callErr != nil {
			return nil, callErr
		}
		v7, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "go-type-spec").Deref(), []vm.Value{arg__14776_6})
		if callErr != nil {
			return nil, callErr
		}
		return v7, nil
	}), ret_ids_4})
	if callErr != nil {
		return nil, callErr
	}
	and__x_44, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "every?").Deref(), []vm.Value{rt.LookupVar("clojure.core", "identity").Deref(), specs_34})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(and__x_44) {
		f_45 = f_3
		ret_ids_46 = ret_ids_4
		specs_47 = specs_34
		and__x_48 = and__x_44
		goto b10
	} else {
		f_49 = f_3
		ret_ids_50 = ret_ids_4
		specs_51 = specs_34
		and__x_52 = and__x_44
		goto b11
	}
b2:
	;
	v82 = vm.NIL
	f_83 = f_5
	ret_ids_84 = ret_ids_6
	goto b3
b3:
	;
	return v82, nil
b4:
	;
	arg__14742_15, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{ret_ids_9})
	if callErr != nil {
		return nil, callErr
	}
	arg__14747_18, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{ret_ids_9})
	if callErr != nil {
		return nil, callErr
	}
	v19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "pos?").Deref(), []vm.Value{arg__14747_18})
	if callErr != nil {
		return nil, callErr
	}
	v22 = v19
	f_23 = f_7
	and__x_24 = and__x_8
	ret_ids_25 = ret_ids_9
	goto b6
b5:
	;
	v22 = and__x_11
	f_23 = f_10
	and__x_24 = and__x_11
	ret_ids_25 = ret_ids_12
	goto b6
b6:
	;
	if vm.IsTruthy(v22) {
		f_3 = f_23
		ret_ids_4 = ret_ids_25
		goto b1
	} else {
		f_5 = f_23
		ret_ids_6 = ret_ids_25
		goto b2
	}
b7:
	;
	v71, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{specs_37})
	if callErr != nil {
		return nil, callErr
	}
	v75 = v71
	f_76 = f_35
	ret_ids_77 = ret_ids_36
	specs_78 = specs_37
	goto b9
b8:
	;
	v75 = vm.NIL
	f_76 = f_38
	ret_ids_77 = ret_ids_39
	specs_78 = specs_40
	goto b9
b9:
	;
	v82 = v75
	f_83 = f_76
	ret_ids_84 = ret_ids_77
	goto b3
b10:
	;
	arg__14788_56, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "distinct").Deref(), []vm.Value{specs_47})
	if callErr != nil {
		return nil, callErr
	}
	arg__14793_59, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "distinct").Deref(), []vm.Value{specs_47})
	if callErr != nil {
		return nil, callErr
	}
	arg__14794_60, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{arg__14793_59})
	if callErr != nil {
		return nil, callErr
	}
	v61 = arg__14794_60 == vm.Int(1)
	v64 = vm.Boolean(v61)
	f_65 = f_45
	ret_ids_66 = ret_ids_46
	specs_67 = specs_47
	and__x_68 = and__x_48
	goto b12
b11:
	;
	v64 = and__x_52
	f_65 = f_49
	ret_ids_66 = ret_ids_50
	specs_67 = specs_51
	and__x_68 = and__x_52
	goto b12
b12:
	;
	if vm.IsTruthy(v64) {
		f_35 = f_65
		ret_ids_36 = ret_ids_66
		specs_37 = specs_67
		goto b7
	} else {
		f_38 = f_65
		ret_ids_39 = ret_ids_66
		specs_40 = specs_67
		goto b8
	}
}
func go_name(arg0 vm.Value) (vm.Value, error) {
	var arg__14814_5 vm.Value
	var arg__14832_11 vm.Value
	var munged_12 vm.Value
	var v20 vm.Value
	var s_13 vm.Value
	var munged_14 vm.Value
	var v25 vm.Value
	var s_15 vm.Value
	var munged_16 vm.Value
	var v28 vm.Value
	var s_29 vm.Value
	var munged_30 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _ = arg__14814_5, arg__14832_11, munged_12, v20, s_13, munged_14, v25, s_15, munged_16, v28, s_29, munged_30
	arg__14814_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var or__x_4 vm.Value
		var ch_5 vm.Value
		var or__x_6 vm.Value
		var ch_7 vm.Value
		var or__x_8 vm.Value
		var v12 vm.Value
		var ch_13 vm.Value
		var or__x_14 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _ = or__x_4, ch_5, or__x_6, ch_7, or__x_8, v12, ch_13, or__x_14
		or__x_4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{rt.LookupVar("ir.lower-go", "go-name-munge-map").Deref(), arg0})
		if callErr != nil {
			return nil, callErr
		}
		if vm.IsTruthy(or__x_4) {
			ch_5 = arg0
			or__x_6 = or__x_4
			goto b1
		} else {
			ch_7 = arg0
			or__x_8 = or__x_4
			goto b2
		}
	b1:
		;
		v12 = or__x_6
		ch_13 = ch_5
		or__x_14 = or__x_6
		goto b3
	b2:
		;
		v12 = ch_7
		ch_13 = ch_7
		or__x_14 = or__x_8
		goto b3
	b3:
		;
		return v12, nil
	}), arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__14832_11, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var or__x_4 vm.Value
		var ch_5 vm.Value
		var or__x_6 vm.Value
		var ch_7 vm.Value
		var or__x_8 vm.Value
		var v12 vm.Value
		var ch_13 vm.Value
		var or__x_14 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _ = or__x_4, ch_5, or__x_6, ch_7, or__x_8, v12, ch_13, or__x_14
		or__x_4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{rt.LookupVar("ir.lower-go", "go-name-munge-map").Deref(), arg0})
		if callErr != nil {
			return nil, callErr
		}
		if vm.IsTruthy(or__x_4) {
			ch_5 = arg0
			or__x_6 = or__x_4
			goto b1
		} else {
			ch_7 = arg0
			or__x_8 = or__x_4
			goto b2
		}
	b1:
		;
		v12 = or__x_6
		ch_13 = ch_5
		or__x_14 = or__x_6
		goto b3
	b2:
		;
		v12 = ch_7
		ch_13 = ch_7
		or__x_14 = or__x_8
		goto b3
	b3:
		;
		return v12, nil
	}), arg0})
	if callErr != nil {
		return nil, callErr
	}
	munged_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "apply").Deref(), []vm.Value{rt.LookupVar("clojure.core", "str").Deref(), arg__14832_11})
	if callErr != nil {
		return nil, callErr
	}
	v20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "contains?").Deref(), []vm.Value{rt.LookupVar("ir.lower-go", "go-reserved-words").Deref(), munged_12})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v20) {
		s_13 = arg0
		munged_14 = munged_12
		goto b1
	} else {
		s_15 = arg0
		munged_16 = munged_12
		goto b2
	}
b1:
	;
	v25, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{munged_14, vm.String("_")})
	if callErr != nil {
		return nil, callErr
	}
	v28 = v25
	s_29 = s_13
	munged_30 = munged_14
	goto b3
b2:
	;
	v28 = munged_16
	s_29 = s_15
	munged_30 = munged_16
	goto b3
b3:
	;
	return v28, nil
}
func go_type_spec(arg0 vm.Value) (vm.Value, error) {
	var or__x_4 bool
	var t_1 vm.Value
	var t_2 vm.Value
	var or__x_27 bool
	var v163 vm.Value
	var t_164 vm.Value
	var t_5 vm.Value
	var or__x_6 bool
	var t_7 vm.Value
	var or__x_8 bool
	var arg__14850_15 vm.Value
	var v16 bool
	var v18 bool
	var t_19 vm.Value
	var or__x_20 vm.Value
	var t_24 vm.Value
	var t_25 vm.Value
	var or__x_50 bool
	var v160 vm.Value
	var t_161 vm.Value
	var t_28 vm.Value
	var or__x_29 bool
	var t_30 vm.Value
	var or__x_31 bool
	var arg__14858_38 vm.Value
	var v39 bool
	var v41 bool
	var t_42 vm.Value
	var or__x_43 vm.Value
	var t_47 vm.Value
	var t_48 vm.Value
	var v81 bool
	var v157 vm.Value
	var t_158 vm.Value
	var t_51 vm.Value
	var or__x_52 bool
	var t_53 vm.Value
	var or__x_54 bool
	var or__x_58 bool
	var v72 bool
	var t_73 vm.Value
	var or__x_74 vm.Value
	var t_59 vm.Value
	var or__x_60 bool
	var t_61 vm.Value
	var or__x_62 bool
	var v66 bool
	var v68 bool
	var t_69 vm.Value
	var or__x_70 vm.Value
	var t_78 vm.Value
	var t_79 vm.Value
	var or__x_88 bool
	var v154 vm.Value
	var t_155 vm.Value
	var t_85 vm.Value
	var t_86 vm.Value
	var and__x_119 vm.Value
	var v151 vm.Value
	var t_152 vm.Value
	var t_89 vm.Value
	var or__x_90 bool
	var t_91 vm.Value
	var or__x_92 bool
	var or__x_96 bool
	var v110 bool
	var t_111 vm.Value
	var or__x_112 vm.Value
	var t_97 vm.Value
	var or__x_98 bool
	var t_99 vm.Value
	var or__x_100 bool
	var v104 bool
	var v106 bool
	var t_107 vm.Value
	var or__x_108 vm.Value
	var t_116 vm.Value
	var t_117 vm.Value
	var v148 vm.Value
	var t_149 vm.Value
	var t_120 vm.Value
	var and__x_121 vm.Value
	var arg__14879_126 vm.Value
	var v128 bool
	var t_122 vm.Value
	var and__x_123 vm.Value
	var v131 vm.Value
	var t_132 vm.Value
	var and__x_133 vm.Value
	var t_137 vm.Value
	var t_138 vm.Value
	var v145 vm.Value
	var t_146 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = or__x_4, t_1, t_2, or__x_27, v163, t_164, t_5, or__x_6, t_7, or__x_8, arg__14850_15, v16, v18, t_19, or__x_20, t_24, t_25, or__x_50, v160, t_161, t_28, or__x_29, t_30, or__x_31, arg__14858_38, v39, v41, t_42, or__x_43, t_47, t_48, v81, v157, t_158, t_51, or__x_52, t_53, or__x_54, or__x_58, v72, t_73, or__x_74, t_59, or__x_60, t_61, or__x_62, v66, v68, t_69, or__x_70, t_78, t_79, or__x_88, v154, t_155, t_85, t_86, and__x_119, v151, t_152, t_89, or__x_90, t_91, or__x_92, or__x_96, v110, t_111, or__x_112, t_97, or__x_98, t_99, or__x_100, v104, v106, t_107, or__x_108, t_116, t_117, v148, t_149, t_120, and__x_121, arg__14879_126, v128, t_122, and__x_123, v131, t_132, and__x_133, t_137, t_138, v145, t_146
	or__x_4 = arg0 == vm.Keyword("int")
	if or__x_4 {
		t_5 = arg0
		or__x_6 = or__x_4
		goto b4
	} else {
		t_7 = arg0
		or__x_8 = or__x_4
		goto b5
	}
b1:
	;
	v163 = vm.String("int")
	t_164 = t_1
	goto b3
b2:
	;
	or__x_27 = t_2 == vm.Keyword("float")
	if or__x_27 {
		t_28 = t_2
		or__x_29 = or__x_27
		goto b10
	} else {
		t_30 = t_2
		or__x_31 = or__x_27
		goto b11
	}
b3:
	;
	return v163, nil
b4:
	;
	v18 = or__x_6
	t_19 = t_5
	or__x_20 = vm.Boolean(or__x_6)
	goto b6
b5:
	;
	arg__14850_15, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("const"), vm.Keyword("int"), vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	v16 = t_7 == arg__14850_15
	v18 = v16
	t_19 = t_7
	or__x_20 = vm.Boolean(or__x_8)
	goto b6
b6:
	;
	if v18 {
		t_1 = t_19
		goto b1
	} else {
		t_2 = t_19
		goto b2
	}
b7:
	;
	v160 = vm.String("float64")
	t_161 = t_24
	goto b9
b8:
	;
	or__x_50 = t_25 == vm.Keyword("bool")
	if or__x_50 {
		t_51 = t_25
		or__x_52 = or__x_50
		goto b16
	} else {
		t_53 = t_25
		or__x_54 = or__x_50
		goto b17
	}
b9:
	;
	v163 = v160
	t_164 = t_161
	goto b3
b10:
	;
	v41 = or__x_29
	t_42 = t_28
	or__x_43 = vm.Boolean(or__x_29)
	goto b12
b11:
	;
	arg__14858_38, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("const"), vm.Keyword("float"), vm.Float(0)})
	if callErr != nil {
		return nil, callErr
	}
	v39 = t_30 == arg__14858_38
	v41 = v39
	t_42 = t_30
	or__x_43 = vm.Boolean(or__x_31)
	goto b12
b12:
	;
	if v41 {
		t_24 = t_42
		goto b7
	} else {
		t_25 = t_42
		goto b8
	}
b13:
	;
	v157 = vm.String("bool")
	t_158 = t_47
	goto b15
b14:
	;
	v81 = t_48 == vm.Keyword("string")
	if v81 {
		t_78 = t_48
		goto b22
	} else {
		t_79 = t_48
		goto b23
	}
b15:
	;
	v160 = v157
	t_161 = t_158
	goto b9
b16:
	;
	v72 = or__x_52
	t_73 = t_51
	or__x_74 = vm.Boolean(or__x_52)
	goto b18
b17:
	;
	or__x_58 = t_53 == vm.Keyword("true")
	if or__x_58 {
		t_59 = t_53
		or__x_60 = or__x_58
		goto b19
	} else {
		t_61 = t_53
		or__x_62 = or__x_58
		goto b20
	}
b18:
	;
	if v72 {
		t_47 = t_73
		goto b13
	} else {
		t_48 = t_73
		goto b14
	}
b19:
	;
	v68 = or__x_60
	t_69 = t_59
	or__x_70 = vm.Boolean(or__x_60)
	goto b21
b20:
	;
	v66 = t_61 == vm.Keyword("false")
	v68 = v66
	t_69 = t_61
	or__x_70 = vm.Boolean(or__x_62)
	goto b21
b21:
	;
	v72 = v68
	t_73 = t_69
	or__x_74 = vm.Boolean(or__x_54)
	goto b18
b22:
	;
	v154 = vm.String("string")
	t_155 = t_78
	goto b24
b23:
	;
	or__x_88 = t_79 == vm.Keyword("unknown")
	if or__x_88 {
		t_89 = t_79
		or__x_90 = or__x_88
		goto b28
	} else {
		t_91 = t_79
		or__x_92 = or__x_88
		goto b29
	}
b24:
	;
	v157 = v154
	t_158 = t_155
	goto b15
b25:
	;
	v151 = vm.String("vm.Value")
	t_152 = t_85
	goto b27
b26:
	;
	and__x_119, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector?").Deref(), []vm.Value{t_86})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(and__x_119) {
		t_120 = t_86
		and__x_121 = and__x_119
		goto b37
	} else {
		t_122 = t_86
		and__x_123 = and__x_119
		goto b38
	}
b27:
	;
	v154 = v151
	t_155 = t_152
	goto b24
b28:
	;
	v110 = or__x_90
	t_111 = t_89
	or__x_112 = vm.Boolean(or__x_90)
	goto b30
b29:
	;
	or__x_96 = t_91 == vm.Keyword("any")
	if or__x_96 {
		t_97 = t_91
		or__x_98 = or__x_96
		goto b31
	} else {
		t_99 = t_91
		or__x_100 = or__x_96
		goto b32
	}
b30:
	;
	if v110 {
		t_85 = t_111
		goto b25
	} else {
		t_86 = t_111
		goto b26
	}
b31:
	;
	v106 = or__x_98
	t_107 = t_97
	or__x_108 = vm.Boolean(or__x_98)
	goto b33
b32:
	;
	v104 = t_99 == vm.Keyword("nil")
	v106 = v104
	t_107 = t_99
	or__x_108 = vm.Boolean(or__x_100)
	goto b33
b33:
	;
	v110 = v106
	t_111 = t_107
	or__x_112 = vm.Boolean(or__x_92)
	goto b30
b34:
	;
	v148 = vm.String("vm.Value")
	t_149 = t_116
	goto b36
b35:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		t_137 = t_117
		goto b40
	} else {
		t_138 = t_117
		goto b41
	}
b36:
	;
	v151 = v148
	t_152 = t_149
	goto b27
b37:
	;
	arg__14879_126, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{t_120})
	if callErr != nil {
		return nil, callErr
	}
	v128 = arg__14879_126 == vm.Keyword("union")
	v131 = vm.Boolean(v128)
	t_132 = t_120
	and__x_133 = and__x_121
	goto b39
b38:
	;
	v131 = and__x_123
	t_132 = t_122
	and__x_133 = and__x_123
	goto b39
b39:
	;
	if vm.IsTruthy(v131) {
		t_116 = t_132
		goto b34
	} else {
		t_117 = t_132
		goto b35
	}
b40:
	;
	v145 = vm.NIL
	t_146 = t_137
	goto b42
b41:
	;
	v145 = vm.NIL
	t_146 = t_138
	goto b42
b42:
	;
	v148 = v145
	t_149 = t_146
	goto b36
}
func import_entry(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var v5 vm.Value
	var callErr error
	_ = v5
	v5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("path"), arg0, vm.Keyword("alias"), arg1})
	if callErr != nil {
		return nil, callErr
	}
	return v5, nil
}
func import_spec_node(arg0 vm.Value) (vm.Value, error) {
	var path_2 vm.Value
	var alias_4 vm.Value
	var entry_5 vm.Value
	var path_6 vm.Value
	var alias_7 vm.Value
	var v13 vm.Value
	var entry_8 vm.Value
	var path_9 vm.Value
	var alias_10 vm.Value
	var v16 vm.Value
	var v18 vm.Value
	var entry_19 vm.Value
	var path_20 vm.Value
	var alias_21 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _ = path_2, alias_4, entry_5, path_6, alias_7, v13, entry_8, path_9, alias_10, v16, v18, entry_19, path_20, alias_21
	path_2, callErr = rt.InvokeValue(vm.Keyword("path"), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	alias_4, callErr = rt.InvokeValue(vm.Keyword("alias"), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(alias_4) {
		entry_5 = arg0
		path_6 = path_2
		alias_7 = alias_4
		goto b1
	} else {
		entry_8 = arg0
		path_9 = path_2
		alias_10 = alias_4
		goto b2
	}
b1:
	;
	v13, callErr = rt.InvokeValue(rt.LookupVar("gogen", "import-spec").Deref(), []vm.Value{path_6, alias_7})
	if callErr != nil {
		return nil, callErr
	}
	v18 = v13
	entry_19 = entry_5
	path_20 = path_6
	alias_21 = alias_7
	goto b3
b2:
	;
	v16, callErr = rt.InvokeValue(rt.LookupVar("gogen", "import-spec").Deref(), []vm.Value{path_9})
	if callErr != nil {
		return nil, callErr
	}
	v18 = v16
	entry_19 = entry_8
	path_20 = path_9
	alias_21 = alias_10
	goto b3
b3:
	;
	return v18, nil
}
func infer_go_type(arg0 vm.Value) (vm.Value, error) {
	var spec_2 vm.Value
	var t_3 vm.Value
	var spec_4 vm.Value
	var v9 vm.Value
	var t_5 vm.Value
	var spec_6 vm.Value
	var v13 vm.Value
	var t_14 vm.Value
	var spec_15 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _ = spec_2, t_3, spec_4, v9, t_5, spec_6, v13, t_14, spec_15
	spec_2, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "go-type-spec").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(spec_2) {
		t_3 = arg0
		spec_4 = spec_2
		goto b1
	} else {
		t_5 = arg0
		spec_6 = spec_2
		goto b2
	}
b1:
	;
	v9, callErr = rt.InvokeValue(rt.LookupVar("gogen", "type").Deref(), []vm.Value{spec_4})
	if callErr != nil {
		return nil, callErr
	}
	v13 = v9
	t_14 = t_3
	spec_15 = spec_4
	goto b3
b2:
	;
	v13 = vm.NIL
	t_14 = t_5
	spec_15 = spec_6
	goto b3
b3:
	;
	return v13, nil
}
func label_name(arg0 vm.Value) (vm.Value, error) {
	var v4 vm.Value
	var callErr error
	_ = v4
	v4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("b"), arg0})
	if callErr != nil {
		return nil, callErr
	}
	return v4, nil
}
func local_carrying_op_QMARK_(arg0 vm.Value) (vm.Value, error) {
	var or__x_2 bool
	var op_3 vm.Value
	var or__x_4 bool
	var op_5 vm.Value
	var or__x_6 bool
	var or__x_10 bool
	var v50 vm.Value
	var op_51 vm.Value
	var or__x_52 vm.Value
	var op_11 vm.Value
	var or__x_12 bool
	var op_13 vm.Value
	var or__x_14 bool
	var or__x_18 bool
	var v46 vm.Value
	var op_47 vm.Value
	var or__x_48 vm.Value
	var op_19 vm.Value
	var or__x_20 bool
	var op_21 vm.Value
	var or__x_22 bool
	var or__x_26 bool
	var v42 vm.Value
	var op_43 vm.Value
	var or__x_44 vm.Value
	var op_27 vm.Value
	var or__x_28 bool
	var op_29 vm.Value
	var or__x_30 bool
	var v36 vm.Value
	var v38 vm.Value
	var op_39 vm.Value
	var or__x_40 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = or__x_2, op_3, or__x_4, op_5, or__x_6, or__x_10, v50, op_51, or__x_52, op_11, or__x_12, op_13, or__x_14, or__x_18, v46, op_47, or__x_48, op_19, or__x_20, op_21, or__x_22, or__x_26, v42, op_43, or__x_44, op_27, or__x_28, op_29, or__x_30, v36, v38, op_39, or__x_40
	or__x_2 = arg0 == vm.Keyword("block-arg")
	if or__x_2 {
		op_3 = arg0
		or__x_4 = or__x_2
		goto b1
	} else {
		op_5 = arg0
		or__x_6 = or__x_2
		goto b2
	}
b1:
	;
	v50 = vm.Boolean(or__x_4)
	op_51 = op_3
	or__x_52 = vm.Boolean(or__x_4)
	goto b3
b2:
	;
	or__x_10 = op_5 == vm.Keyword("call")
	if or__x_10 {
		op_11 = op_5
		or__x_12 = or__x_10
		goto b4
	} else {
		op_13 = op_5
		or__x_14 = or__x_10
		goto b5
	}
b3:
	;
	return v50, nil
b4:
	;
	v46 = vm.Boolean(or__x_12)
	op_47 = op_11
	or__x_48 = vm.Boolean(or__x_12)
	goto b6
b5:
	;
	or__x_18 = op_13 == vm.Keyword("inc")
	if or__x_18 {
		op_19 = op_13
		or__x_20 = or__x_18
		goto b7
	} else {
		op_21 = op_13
		or__x_22 = or__x_18
		goto b8
	}
b6:
	;
	v50 = v46
	op_51 = op_47
	or__x_52 = vm.Boolean(or__x_6)
	goto b3
b7:
	;
	v42 = vm.Boolean(or__x_20)
	op_43 = op_19
	or__x_44 = vm.Boolean(or__x_20)
	goto b9
b8:
	;
	or__x_26 = op_21 == vm.Keyword("dec")
	if or__x_26 {
		op_27 = op_21
		or__x_28 = or__x_26
		goto b10
	} else {
		op_29 = op_21
		or__x_30 = or__x_26
		goto b11
	}
b9:
	;
	v46 = v42
	op_47 = op_43
	or__x_48 = vm.Boolean(or__x_14)
	goto b6
b10:
	;
	v38 = vm.Boolean(or__x_28)
	op_39 = op_27
	or__x_40 = vm.Boolean(or__x_28)
	goto b12
b11:
	;
	v36, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "contains?").Deref(), []vm.Value{rt.LookupVar("ir.lower-go", "binary-op").Deref(), op_29})
	if callErr != nil {
		return nil, callErr
	}
	v38 = v36
	op_39 = op_29
	or__x_40 = vm.Boolean(or__x_30)
	goto b12
b12:
	;
	v42 = v38
	op_43 = op_39
	or__x_44 = vm.Boolean(or__x_22)
	goto b9
}
func local_decls(arg0 vm.Value) (vm.Value, error) {
	var arg__15417_2 vm.Value
	var arg__15422_5 vm.Value
	var ids_6 vm.Value
	var remaining_7 vm.Value
	var out_8 vm.Value
	var f_9 vm.Value
	var v88 vm.Value
	var v21 vm.Value
	var ids_12 vm.Value
	var remaining_13 vm.Value
	var out_14 vm.Value
	var f_15 vm.Value
	var v91 vm.Value
	var ids_16 vm.Value
	var remaining_17 vm.Value
	var out_18 vm.Value
	var f_19 vm.Value
	var v89 vm.Value
	var nid_25 vm.Value
	var arg__15434_27 vm.Value
	var arg__15441_30 vm.Value
	var go_type_31 vm.Value
	var v45 vm.Value
	var v78 vm.Value
	var ids_79 vm.Value
	var remaining_80 vm.Value
	var out_81 vm.Value
	var f_82 vm.Value
	var ids_32 vm.Value
	var remaining_33 vm.Value
	var out_34 vm.Value
	var f_35 vm.Value
	var nid_36 vm.Value
	var go_type_37 vm.Value
	var v92 vm.Value
	var ids_38 vm.Value
	var remaining_39 vm.Value
	var out_40 vm.Value
	var f_41 vm.Value
	var nid_42 vm.Value
	var go_type_43 vm.Value
	var v90 vm.Value
	var v50 vm.Value
	var arg__15454_52 vm.Value
	var arg__15463_56 vm.Value
	var arg__15466_58 vm.Value
	var arg__15474_61 vm.Value
	var arg__15483_65 vm.Value
	var arg__15486_67 vm.Value
	var v68 vm.Value
	var v70 vm.Value
	var ids_71 vm.Value
	var remaining_72 vm.Value
	var out_73 vm.Value
	var f_74 vm.Value
	var nid_75 vm.Value
	var go_type_76 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__15417_2, arg__15422_5, ids_6, remaining_7, out_8, f_9, v88, v21, ids_12, remaining_13, out_14, f_15, v91, ids_16, remaining_17, out_18, f_19, v89, nid_25, arg__15434_27, arg__15441_30, go_type_31, v45, v78, ids_79, remaining_80, out_81, f_82, ids_32, remaining_33, out_34, f_35, nid_36, go_type_37, v92, ids_38, remaining_39, out_40, f_41, nid_42, go_type_43, v90, v50, arg__15454_52, arg__15463_56, arg__15466_58, arg__15474_61, arg__15483_65, arg__15486_67, v68, v70, ids_71, remaining_72, out_73, f_74, nid_75, go_type_76
	arg__15417_2, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "collect-local-ids").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__15422_5, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "collect-local-ids").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	ids_6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "distinct").Deref(), []vm.Value{arg__15422_5})
	if callErr != nil {
		return nil, callErr
	}
	remaining_7 = ids_6
	out_8 = vm.NewArrayVector([]vm.Value{})
	f_9 = arg0
	v88 = vm.NIL
	goto b1
b1:
	;
	v21, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "empty?").Deref(), []vm.Value{remaining_7})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v21) {
		ids_12 = ids_6
		remaining_13 = remaining_7
		out_14 = out_8
		f_15 = f_9
		v91 = v88
		goto b2
	} else {
		ids_16 = ids_6
		remaining_17 = remaining_7
		out_18 = out_8
		f_19 = f_9
		v89 = v88
		goto b3
	}
b2:
	;
	v78 = out_14
	ids_79 = ids_12
	remaining_80 = remaining_13
	out_81 = out_14
	f_82 = f_15
	goto b4
b3:
	;
	nid_25, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{remaining_17})
	if callErr != nil {
		return nil, callErr
	}
	arg__15434_27, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{nid_25, f_19})
	if callErr != nil {
		return nil, callErr
	}
	arg__15441_30, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{nid_25, f_19})
	if callErr != nil {
		return nil, callErr
	}
	go_type_31, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "infer-go-type").Deref(), []vm.Value{arg__15441_30})
	if callErr != nil {
		return nil, callErr
	}
	v45, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{go_type_31})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v45) {
		ids_32 = ids_16
		remaining_33 = remaining_17
		out_34 = out_18
		f_35 = f_19
		nid_36 = nid_25
		go_type_37 = go_type_31
		v92 = v89
		goto b5
	} else {
		ids_38 = ids_16
		remaining_39 = remaining_17
		out_40 = out_18
		f_41 = f_19
		nid_42 = nid_25
		go_type_43 = go_type_31
		v90 = v89
		goto b6
	}
b4:
	;
	return v78, nil
b5:
	;
	v70 = vm.NIL
	ids_71 = ids_32
	remaining_72 = remaining_33
	out_73 = out_34
	f_74 = f_35
	nid_75 = nid_36
	go_type_76 = go_type_37
	goto b7
b6:
	;
	v50, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{remaining_39})
	if callErr != nil {
		return nil, callErr
	}
	arg__15454_52, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "value-name").Deref(), []vm.Value{f_41, nid_42})
	if callErr != nil {
		return nil, callErr
	}
	arg__15463_56, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "value-name").Deref(), []vm.Value{f_41, nid_42})
	if callErr != nil {
		return nil, callErr
	}
	arg__15466_58, callErr = rt.InvokeValue(rt.LookupVar("gogen", "var-decl").Deref(), []vm.Value{arg__15463_56, go_type_43, v90})
	if callErr != nil {
		return nil, callErr
	}
	arg__15474_61, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "value-name").Deref(), []vm.Value{f_41, nid_42})
	if callErr != nil {
		return nil, callErr
	}
	arg__15483_65, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "value-name").Deref(), []vm.Value{f_41, nid_42})
	if callErr != nil {
		return nil, callErr
	}
	arg__15486_67, callErr = rt.InvokeValue(rt.LookupVar("gogen", "var-decl").Deref(), []vm.Value{arg__15483_65, go_type_43, v90})
	if callErr != nil {
		return nil, callErr
	}
	v68, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{out_40, arg__15486_67})
	if callErr != nil {
		return nil, callErr
	}
	remaining_7 = v50
	out_8 = v68
	f_9 = f_41
	v88 = v90
	goto b1
b7:
	;
	v78 = v70
	ids_79 = ids_71
	remaining_80 = remaining_72
	out_81 = out_73
	f_82 = f_74
	goto b4
}
func lower_block_stmts(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__15503_10 vm.Value
	var v11 bool
	var f_3 vm.Value
	var closed_exprs_4 vm.Value
	var bid_5 vm.Value
	var f_6 vm.Value
	var closed_exprs_7 vm.Value
	var bid_8 vm.Value
	var arg__15508_17 vm.Value
	var arg__15514_21 vm.Value
	var arg__15516_23 vm.Value
	var v24 vm.Value
	var label_stmts_26 vm.Value
	var f_27 vm.Value
	var closed_exprs_28 vm.Value
	var bid_29 vm.Value
	var insts_31 vm.Value
	var remaining_32 vm.Value
	var out_33 vm.Value
	var f_34 vm.Value
	var closed_exprs_35 vm.Value
	var v53 vm.Value
	var label_stmts_38 vm.Value
	var bid_39 vm.Value
	var insts_40 vm.Value
	var remaining_41 vm.Value
	var out_42 vm.Value
	var f_43 vm.Value
	var closed_exprs_44 vm.Value
	var label_stmts_45 vm.Value
	var bid_46 vm.Value
	var insts_47 vm.Value
	var remaining_48 vm.Value
	var out_49 vm.Value
	var f_50 vm.Value
	var closed_exprs_51 vm.Value
	var arg__15530_57 vm.Value
	var arg__15537_60 vm.Value
	var stmts_61 vm.Value
	var v79 vm.Value
	var inst_stmts_331 vm.Value
	var label_stmts_332 vm.Value
	var bid_333 vm.Value
	var insts_334 vm.Value
	var remaining_335 vm.Value
	var out_336 vm.Value
	var f_337 vm.Value
	var closed_exprs_338 vm.Value
	var term_stmts_340 vm.Value
	var label_stmts_62 vm.Value
	var bid_63 vm.Value
	var insts_64 vm.Value
	var remaining_65 vm.Value
	var out_66 vm.Value
	var f_67 vm.Value
	var closed_exprs_68 vm.Value
	var stmts_69 vm.Value
	var arg__15544_82 vm.Value
	var arg__15550_85 vm.Value
	var op_86 vm.Value
	var or__x_106 bool
	var label_stmts_70 vm.Value
	var bid_71 vm.Value
	var insts_72 vm.Value
	var remaining_73 vm.Value
	var out_74 vm.Value
	var f_75 vm.Value
	var closed_exprs_76 vm.Value
	var stmts_77 vm.Value
	var v317 vm.Value
	var v319 vm.Value
	var v321 vm.Value
	var label_stmts_322 vm.Value
	var bid_323 vm.Value
	var insts_324 vm.Value
	var remaining_325 vm.Value
	var out_326 vm.Value
	var f_327 vm.Value
	var closed_exprs_328 vm.Value
	var stmts_329 vm.Value
	var label_stmts_87 vm.Value
	var bid_88 vm.Value
	var insts_89 vm.Value
	var remaining_90 vm.Value
	var out_91 vm.Value
	var f_92 vm.Value
	var closed_exprs_93 vm.Value
	var stmts_94 vm.Value
	var op_95 vm.Value
	var v301 vm.Value
	var label_stmts_96 vm.Value
	var bid_97 vm.Value
	var insts_98 vm.Value
	var remaining_99 vm.Value
	var out_100 vm.Value
	var f_101 vm.Value
	var closed_exprs_102 vm.Value
	var stmts_103 vm.Value
	var op_104 vm.Value
	var v305 vm.Value
	var label_stmts_306 vm.Value
	var bid_307 vm.Value
	var insts_308 vm.Value
	var remaining_309 vm.Value
	var out_310 vm.Value
	var f_311 vm.Value
	var closed_exprs_312 vm.Value
	var stmts_313 vm.Value
	var op_314 vm.Value
	var label_stmts_107 vm.Value
	var bid_108 vm.Value
	var insts_109 vm.Value
	var remaining_110 vm.Value
	var out_111 vm.Value
	var f_112 vm.Value
	var closed_exprs_113 vm.Value
	var stmts_114 vm.Value
	var op_115 vm.Value
	var or__x_116 bool
	var label_stmts_117 vm.Value
	var bid_118 vm.Value
	var insts_119 vm.Value
	var remaining_120 vm.Value
	var out_121 vm.Value
	var f_122 vm.Value
	var closed_exprs_123 vm.Value
	var stmts_124 vm.Value
	var op_125 vm.Value
	var or__x_126 bool
	var or__x_130 bool
	var v288 vm.Value
	var label_stmts_289 vm.Value
	var bid_290 vm.Value
	var insts_291 vm.Value
	var remaining_292 vm.Value
	var out_293 vm.Value
	var f_294 vm.Value
	var closed_exprs_295 vm.Value
	var stmts_296 vm.Value
	var op_297 vm.Value
	var or__x_298 vm.Value
	var label_stmts_131 vm.Value
	var bid_132 vm.Value
	var insts_133 vm.Value
	var remaining_134 vm.Value
	var out_135 vm.Value
	var f_136 vm.Value
	var closed_exprs_137 vm.Value
	var stmts_138 vm.Value
	var op_139 vm.Value
	var or__x_140 bool
	var label_stmts_141 vm.Value
	var bid_142 vm.Value
	var insts_143 vm.Value
	var remaining_144 vm.Value
	var out_145 vm.Value
	var f_146 vm.Value
	var closed_exprs_147 vm.Value
	var stmts_148 vm.Value
	var op_149 vm.Value
	var or__x_150 bool
	var or__x_154 bool
	var v276 vm.Value
	var label_stmts_277 vm.Value
	var bid_278 vm.Value
	var insts_279 vm.Value
	var remaining_280 vm.Value
	var out_281 vm.Value
	var f_282 vm.Value
	var closed_exprs_283 vm.Value
	var stmts_284 vm.Value
	var op_285 vm.Value
	var or__x_286 vm.Value
	var label_stmts_155 vm.Value
	var bid_156 vm.Value
	var insts_157 vm.Value
	var remaining_158 vm.Value
	var out_159 vm.Value
	var f_160 vm.Value
	var closed_exprs_161 vm.Value
	var stmts_162 vm.Value
	var op_163 vm.Value
	var or__x_164 bool
	var label_stmts_165 vm.Value
	var bid_166 vm.Value
	var insts_167 vm.Value
	var remaining_168 vm.Value
	var out_169 vm.Value
	var f_170 vm.Value
	var closed_exprs_171 vm.Value
	var stmts_172 vm.Value
	var op_173 vm.Value
	var or__x_174 bool
	var or__x_178 bool
	var v264 vm.Value
	var label_stmts_265 vm.Value
	var bid_266 vm.Value
	var insts_267 vm.Value
	var remaining_268 vm.Value
	var out_269 vm.Value
	var f_270 vm.Value
	var closed_exprs_271 vm.Value
	var stmts_272 vm.Value
	var op_273 vm.Value
	var or__x_274 vm.Value
	var label_stmts_179 vm.Value
	var bid_180 vm.Value
	var insts_181 vm.Value
	var remaining_182 vm.Value
	var out_183 vm.Value
	var f_184 vm.Value
	var closed_exprs_185 vm.Value
	var stmts_186 vm.Value
	var op_187 vm.Value
	var or__x_188 bool
	var label_stmts_189 vm.Value
	var bid_190 vm.Value
	var insts_191 vm.Value
	var remaining_192 vm.Value
	var out_193 vm.Value
	var f_194 vm.Value
	var closed_exprs_195 vm.Value
	var stmts_196 vm.Value
	var op_197 vm.Value
	var or__x_198 bool
	var or__x_202 bool
	var v252 vm.Value
	var label_stmts_253 vm.Value
	var bid_254 vm.Value
	var insts_255 vm.Value
	var remaining_256 vm.Value
	var out_257 vm.Value
	var f_258 vm.Value
	var closed_exprs_259 vm.Value
	var stmts_260 vm.Value
	var op_261 vm.Value
	var or__x_262 vm.Value
	var label_stmts_203 vm.Value
	var bid_204 vm.Value
	var insts_205 vm.Value
	var remaining_206 vm.Value
	var out_207 vm.Value
	var f_208 vm.Value
	var closed_exprs_209 vm.Value
	var stmts_210 vm.Value
	var op_211 vm.Value
	var or__x_212 bool
	var label_stmts_213 vm.Value
	var bid_214 vm.Value
	var insts_215 vm.Value
	var remaining_216 vm.Value
	var out_217 vm.Value
	var f_218 vm.Value
	var closed_exprs_219 vm.Value
	var stmts_220 vm.Value
	var op_221 vm.Value
	var or__x_222 bool
	var arg__15566_226 vm.Value
	var arg__15572_229 vm.Value
	var arg__15573_230 vm.Value
	var arg__15579_233 vm.Value
	var arg__15585_236 vm.Value
	var arg__15586_237 vm.Value
	var v238 vm.Value
	var v240 vm.Value
	var label_stmts_241 vm.Value
	var bid_242 vm.Value
	var insts_243 vm.Value
	var remaining_244 vm.Value
	var out_245 vm.Value
	var f_246 vm.Value
	var closed_exprs_247 vm.Value
	var stmts_248 vm.Value
	var op_249 vm.Value
	var or__x_250 vm.Value
	var inst_stmts_341 vm.Value
	var label_stmts_342 vm.Value
	var bid_343 vm.Value
	var insts_344 vm.Value
	var remaining_345 vm.Value
	var out_346 vm.Value
	var f_347 vm.Value
	var closed_exprs_348 vm.Value
	var term_stmts_349 vm.Value
	var arg__15612_395 vm.Value
	var arg__15621_398 vm.Value
	var v399 vm.Value
	var inst_stmts_350 vm.Value
	var label_stmts_351 vm.Value
	var bid_352 vm.Value
	var insts_353 vm.Value
	var remaining_354 vm.Value
	var out_355 vm.Value
	var f_356 vm.Value
	var closed_exprs_357 vm.Value
	var term_stmts_358 vm.Value
	var v403 vm.Value
	var inst_stmts_404 vm.Value
	var label_stmts_405 vm.Value
	var bid_406 vm.Value
	var insts_407 vm.Value
	var remaining_408 vm.Value
	var out_409 vm.Value
	var f_410 vm.Value
	var closed_exprs_411 vm.Value
	var term_stmts_412 vm.Value
	var and__x_359 vm.Value
	var inst_stmts_360 vm.Value
	var label_stmts_361 vm.Value
	var bid_362 vm.Value
	var insts_363 vm.Value
	var remaining_364 vm.Value
	var out_365 vm.Value
	var f_366 vm.Value
	var closed_exprs_367 vm.Value
	var term_stmts_368 vm.Value
	var and__x_369 vm.Value
	var inst_stmts_370 vm.Value
	var label_stmts_371 vm.Value
	var bid_372 vm.Value
	var insts_373 vm.Value
	var remaining_374 vm.Value
	var out_375 vm.Value
	var f_376 vm.Value
	var closed_exprs_377 vm.Value
	var term_stmts_378 vm.Value
	var v382 vm.Value
	var and__x_383 vm.Value
	var inst_stmts_384 vm.Value
	var label_stmts_385 vm.Value
	var bid_386 vm.Value
	var insts_387 vm.Value
	var remaining_388 vm.Value
	var out_389 vm.Value
	var f_390 vm.Value
	var closed_exprs_391 vm.Value
	var term_stmts_392 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__15503_10, v11, f_3, closed_exprs_4, bid_5, f_6, closed_exprs_7, bid_8, arg__15508_17, arg__15514_21, arg__15516_23, v24, label_stmts_26, f_27, closed_exprs_28, bid_29, insts_31, remaining_32, out_33, f_34, closed_exprs_35, v53, label_stmts_38, bid_39, insts_40, remaining_41, out_42, f_43, closed_exprs_44, label_stmts_45, bid_46, insts_47, remaining_48, out_49, f_50, closed_exprs_51, arg__15530_57, arg__15537_60, stmts_61, v79, inst_stmts_331, label_stmts_332, bid_333, insts_334, remaining_335, out_336, f_337, closed_exprs_338, term_stmts_340, label_stmts_62, bid_63, insts_64, remaining_65, out_66, f_67, closed_exprs_68, stmts_69, arg__15544_82, arg__15550_85, op_86, or__x_106, label_stmts_70, bid_71, insts_72, remaining_73, out_74, f_75, closed_exprs_76, stmts_77, v317, v319, v321, label_stmts_322, bid_323, insts_324, remaining_325, out_326, f_327, closed_exprs_328, stmts_329, label_stmts_87, bid_88, insts_89, remaining_90, out_91, f_92, closed_exprs_93, stmts_94, op_95, v301, label_stmts_96, bid_97, insts_98, remaining_99, out_100, f_101, closed_exprs_102, stmts_103, op_104, v305, label_stmts_306, bid_307, insts_308, remaining_309, out_310, f_311, closed_exprs_312, stmts_313, op_314, label_stmts_107, bid_108, insts_109, remaining_110, out_111, f_112, closed_exprs_113, stmts_114, op_115, or__x_116, label_stmts_117, bid_118, insts_119, remaining_120, out_121, f_122, closed_exprs_123, stmts_124, op_125, or__x_126, or__x_130, v288, label_stmts_289, bid_290, insts_291, remaining_292, out_293, f_294, closed_exprs_295, stmts_296, op_297, or__x_298, label_stmts_131, bid_132, insts_133, remaining_134, out_135, f_136, closed_exprs_137, stmts_138, op_139, or__x_140, label_stmts_141, bid_142, insts_143, remaining_144, out_145, f_146, closed_exprs_147, stmts_148, op_149, or__x_150, or__x_154, v276, label_stmts_277, bid_278, insts_279, remaining_280, out_281, f_282, closed_exprs_283, stmts_284, op_285, or__x_286, label_stmts_155, bid_156, insts_157, remaining_158, out_159, f_160, closed_exprs_161, stmts_162, op_163, or__x_164, label_stmts_165, bid_166, insts_167, remaining_168, out_169, f_170, closed_exprs_171, stmts_172, op_173, or__x_174, or__x_178, v264, label_stmts_265, bid_266, insts_267, remaining_268, out_269, f_270, closed_exprs_271, stmts_272, op_273, or__x_274, label_stmts_179, bid_180, insts_181, remaining_182, out_183, f_184, closed_exprs_185, stmts_186, op_187, or__x_188, label_stmts_189, bid_190, insts_191, remaining_192, out_193, f_194, closed_exprs_195, stmts_196, op_197, or__x_198, or__x_202, v252, label_stmts_253, bid_254, insts_255, remaining_256, out_257, f_258, closed_exprs_259, stmts_260, op_261, or__x_262, label_stmts_203, bid_204, insts_205, remaining_206, out_207, f_208, closed_exprs_209, stmts_210, op_211, or__x_212, label_stmts_213, bid_214, insts_215, remaining_216, out_217, f_218, closed_exprs_219, stmts_220, op_221, or__x_222, arg__15566_226, arg__15572_229, arg__15573_230, arg__15579_233, arg__15585_236, arg__15586_237, v238, v240, label_stmts_241, bid_242, insts_243, remaining_244, out_245, f_246, closed_exprs_247, stmts_248, op_249, or__x_250, inst_stmts_341, label_stmts_342, bid_343, insts_344, remaining_345, out_346, f_347, closed_exprs_348, term_stmts_349, arg__15612_395, arg__15621_398, v399, inst_stmts_350, label_stmts_351, bid_352, insts_353, remaining_354, out_355, f_356, closed_exprs_357, term_stmts_358, v403, inst_stmts_404, label_stmts_405, bid_406, insts_407, remaining_408, out_409, f_410, closed_exprs_411, term_stmts_412, and__x_359, inst_stmts_360, label_stmts_361, bid_362, insts_363, remaining_364, out_365, f_366, closed_exprs_367, term_stmts_368, and__x_369, inst_stmts_370, label_stmts_371, bid_372, insts_373, remaining_374, out_375, f_376, closed_exprs_377, term_stmts_378, v382, and__x_383, inst_stmts_384, label_stmts_385, bid_386, insts_387, remaining_388, out_389, f_390, closed_exprs_391, term_stmts_392
	arg__15503_10, callErr = rt.InvokeValue(rt.LookupVar("ir", "entry-block").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v11 = arg2 == arg__15503_10
	if v11 {
		f_3 = arg0
		closed_exprs_4 = arg1
		bid_5 = arg2
		goto b1
	} else {
		f_6 = arg0
		closed_exprs_7 = arg1
		bid_8 = arg2
		goto b2
	}
b1:
	;
	label_stmts_26 = vm.NewArrayVector([]vm.Value{})
	f_27 = f_3
	closed_exprs_28 = closed_exprs_4
	bid_29 = bid_5
	goto b3
b2:
	;
	arg__15508_17, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "label-name").Deref(), []vm.Value{bid_8})
	if callErr != nil {
		return nil, callErr
	}
	arg__15514_21, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "label-name").Deref(), []vm.Value{bid_8})
	if callErr != nil {
		return nil, callErr
	}
	arg__15516_23, callErr = rt.InvokeValue(rt.LookupVar("gogen", "label-stmt").Deref(), []vm.Value{arg__15514_21, vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	v24, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__15516_23})
	if callErr != nil {
		return nil, callErr
	}
	label_stmts_26 = v24
	f_27 = f_6
	closed_exprs_28 = closed_exprs_7
	bid_29 = bid_8
	goto b3
b3:
	;
	insts_31, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-insts").Deref(), []vm.Value{bid_29, f_27})
	if callErr != nil {
		return nil, callErr
	}
	remaining_32 = insts_31
	out_33 = vm.NewArrayVector([]vm.Value{})
	f_34 = f_27
	closed_exprs_35 = closed_exprs_28
	goto b4
b4:
	;
	v53, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "empty?").Deref(), []vm.Value{remaining_32})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v53) {
		label_stmts_38 = label_stmts_26
		bid_39 = bid_29
		insts_40 = insts_31
		remaining_41 = remaining_32
		out_42 = out_33
		f_43 = f_34
		closed_exprs_44 = closed_exprs_35
		goto b5
	} else {
		label_stmts_45 = label_stmts_26
		bid_46 = bid_29
		insts_47 = insts_31
		remaining_48 = remaining_32
		out_49 = out_33
		f_50 = f_34
		closed_exprs_51 = closed_exprs_35
		goto b6
	}
b5:
	;
	inst_stmts_331 = out_42
	label_stmts_332 = label_stmts_38
	bid_333 = bid_39
	insts_334 = insts_40
	remaining_335 = remaining_41
	out_336 = out_42
	f_337 = f_43
	closed_exprs_338 = closed_exprs_44
	goto b7
b6:
	;
	arg__15530_57, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{remaining_48})
	if callErr != nil {
		return nil, callErr
	}
	arg__15537_60, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{remaining_48})
	if callErr != nil {
		return nil, callErr
	}
	stmts_61, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "lower-inst-stmts").Deref(), []vm.Value{f_50, closed_exprs_51, arg__15537_60})
	if callErr != nil {
		return nil, callErr
	}
	v79, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{stmts_61})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v79) {
		label_stmts_62 = label_stmts_45
		bid_63 = bid_46
		insts_64 = insts_47
		remaining_65 = remaining_48
		out_66 = out_49
		f_67 = f_50
		closed_exprs_68 = closed_exprs_51
		stmts_69 = stmts_61
		goto b8
	} else {
		label_stmts_70 = label_stmts_45
		bid_71 = bid_46
		insts_72 = insts_47
		remaining_73 = remaining_48
		out_74 = out_49
		f_75 = f_50
		closed_exprs_76 = closed_exprs_51
		stmts_77 = stmts_61
		goto b9
	}
b7:
	;
	term_stmts_340, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "lower-terminator").Deref(), []vm.Value{f_337, closed_exprs_338, bid_333})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(inst_stmts_331) {
		and__x_359 = inst_stmts_331
		inst_stmts_360 = inst_stmts_331
		label_stmts_361 = label_stmts_332
		bid_362 = bid_333
		insts_363 = insts_334
		remaining_364 = remaining_335
		out_365 = out_336
		f_366 = f_337
		closed_exprs_367 = closed_exprs_338
		term_stmts_368 = term_stmts_340
		goto b32
	} else {
		and__x_369 = inst_stmts_331
		inst_stmts_370 = inst_stmts_331
		label_stmts_371 = label_stmts_332
		bid_372 = bid_333
		insts_373 = insts_334
		remaining_374 = remaining_335
		out_375 = out_336
		f_376 = f_337
		closed_exprs_377 = closed_exprs_338
		term_stmts_378 = term_stmts_340
		goto b33
	}
b8:
	;
	arg__15544_82, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{remaining_65})
	if callErr != nil {
		return nil, callErr
	}
	arg__15550_85, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{remaining_65})
	if callErr != nil {
		return nil, callErr
	}
	op_86, callErr = rt.InvokeValue(rt.LookupVar("ir", "op").Deref(), []vm.Value{arg__15550_85, f_67})
	if callErr != nil {
		return nil, callErr
	}
	or__x_106 = op_86 == vm.Keyword("load-arg")
	if or__x_106 {
		label_stmts_107 = label_stmts_62
		bid_108 = bid_63
		insts_109 = insts_64
		remaining_110 = remaining_65
		out_111 = out_66
		f_112 = f_67
		closed_exprs_113 = closed_exprs_68
		stmts_114 = stmts_69
		op_115 = op_86
		or__x_116 = or__x_106
		goto b14
	} else {
		label_stmts_117 = label_stmts_62
		bid_118 = bid_63
		insts_119 = insts_64
		remaining_120 = remaining_65
		out_121 = out_66
		f_122 = f_67
		closed_exprs_123 = closed_exprs_68
		stmts_124 = stmts_69
		op_125 = op_86
		or__x_126 = or__x_106
		goto b15
	}
b9:
	;
	v317, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{remaining_73})
	if callErr != nil {
		return nil, callErr
	}
	v319, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{out_74, stmts_77})
	if callErr != nil {
		return nil, callErr
	}
	remaining_32 = v317
	out_33 = v319
	f_34 = f_75
	closed_exprs_35 = closed_exprs_76
	goto b4
b10:
	;
	inst_stmts_331 = v321
	label_stmts_332 = label_stmts_322
	bid_333 = bid_323
	insts_334 = insts_324
	remaining_335 = remaining_325
	out_336 = out_326
	f_337 = f_327
	closed_exprs_338 = closed_exprs_328
	goto b7
b11:
	;
	v301, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{remaining_90})
	if callErr != nil {
		return nil, callErr
	}
	remaining_32 = v301
	out_33 = out_91
	f_34 = f_92
	closed_exprs_35 = closed_exprs_93
	goto b4
b12:
	;
	v305 = vm.NIL
	label_stmts_306 = label_stmts_96
	bid_307 = bid_97
	insts_308 = insts_98
	remaining_309 = remaining_99
	out_310 = out_100
	f_311 = f_101
	closed_exprs_312 = closed_exprs_102
	stmts_313 = stmts_103
	op_314 = op_104
	goto b13
b13:
	;
	v321 = v305
	label_stmts_322 = label_stmts_306
	bid_323 = bid_307
	insts_324 = insts_308
	remaining_325 = remaining_309
	out_326 = out_310
	f_327 = f_311
	closed_exprs_328 = closed_exprs_312
	stmts_329 = stmts_313
	goto b10
b14:
	;
	v288 = vm.Boolean(or__x_116)
	label_stmts_289 = label_stmts_107
	bid_290 = bid_108
	insts_291 = insts_109
	remaining_292 = remaining_110
	out_293 = out_111
	f_294 = f_112
	closed_exprs_295 = closed_exprs_113
	stmts_296 = stmts_114
	op_297 = op_115
	or__x_298 = vm.Boolean(or__x_116)
	goto b16
b15:
	;
	or__x_130 = op_125 == vm.Keyword("load-var")
	if or__x_130 {
		label_stmts_131 = label_stmts_117
		bid_132 = bid_118
		insts_133 = insts_119
		remaining_134 = remaining_120
		out_135 = out_121
		f_136 = f_122
		closed_exprs_137 = closed_exprs_123
		stmts_138 = stmts_124
		op_139 = op_125
		or__x_140 = or__x_130
		goto b17
	} else {
		label_stmts_141 = label_stmts_117
		bid_142 = bid_118
		insts_143 = insts_119
		remaining_144 = remaining_120
		out_145 = out_121
		f_146 = f_122
		closed_exprs_147 = closed_exprs_123
		stmts_148 = stmts_124
		op_149 = op_125
		or__x_150 = or__x_130
		goto b18
	}
b16:
	;
	if vm.IsTruthy(v288) {
		label_stmts_87 = label_stmts_289
		bid_88 = bid_290
		insts_89 = insts_291
		remaining_90 = remaining_292
		out_91 = out_293
		f_92 = f_294
		closed_exprs_93 = closed_exprs_295
		stmts_94 = stmts_296
		op_95 = op_297
		goto b11
	} else {
		label_stmts_96 = label_stmts_289
		bid_97 = bid_290
		insts_98 = insts_291
		remaining_99 = remaining_292
		out_100 = out_293
		f_101 = f_294
		closed_exprs_102 = closed_exprs_295
		stmts_103 = stmts_296
		op_104 = op_297
		goto b12
	}
b17:
	;
	v276 = vm.Boolean(or__x_140)
	label_stmts_277 = label_stmts_131
	bid_278 = bid_132
	insts_279 = insts_133
	remaining_280 = remaining_134
	out_281 = out_135
	f_282 = f_136
	closed_exprs_283 = closed_exprs_137
	stmts_284 = stmts_138
	op_285 = op_139
	or__x_286 = vm.Boolean(or__x_140)
	goto b19
b18:
	;
	or__x_154 = op_149 == vm.Keyword("load-closed")
	if or__x_154 {
		label_stmts_155 = label_stmts_141
		bid_156 = bid_142
		insts_157 = insts_143
		remaining_158 = remaining_144
		out_159 = out_145
		f_160 = f_146
		closed_exprs_161 = closed_exprs_147
		stmts_162 = stmts_148
		op_163 = op_149
		or__x_164 = or__x_154
		goto b20
	} else {
		label_stmts_165 = label_stmts_141
		bid_166 = bid_142
		insts_167 = insts_143
		remaining_168 = remaining_144
		out_169 = out_145
		f_170 = f_146
		closed_exprs_171 = closed_exprs_147
		stmts_172 = stmts_148
		op_173 = op_149
		or__x_174 = or__x_154
		goto b21
	}
b19:
	;
	v288 = v276
	label_stmts_289 = label_stmts_277
	bid_290 = bid_278
	insts_291 = insts_279
	remaining_292 = remaining_280
	out_293 = out_281
	f_294 = f_282
	closed_exprs_295 = closed_exprs_283
	stmts_296 = stmts_284
	op_297 = op_285
	or__x_298 = vm.Boolean(or__x_126)
	goto b16
b20:
	;
	v264 = vm.Boolean(or__x_164)
	label_stmts_265 = label_stmts_155
	bid_266 = bid_156
	insts_267 = insts_157
	remaining_268 = remaining_158
	out_269 = out_159
	f_270 = f_160
	closed_exprs_271 = closed_exprs_161
	stmts_272 = stmts_162
	op_273 = op_163
	or__x_274 = vm.Boolean(or__x_164)
	goto b22
b21:
	;
	or__x_178 = op_173 == vm.Keyword("const")
	if or__x_178 {
		label_stmts_179 = label_stmts_165
		bid_180 = bid_166
		insts_181 = insts_167
		remaining_182 = remaining_168
		out_183 = out_169
		f_184 = f_170
		closed_exprs_185 = closed_exprs_171
		stmts_186 = stmts_172
		op_187 = op_173
		or__x_188 = or__x_178
		goto b23
	} else {
		label_stmts_189 = label_stmts_165
		bid_190 = bid_166
		insts_191 = insts_167
		remaining_192 = remaining_168
		out_193 = out_169
		f_194 = f_170
		closed_exprs_195 = closed_exprs_171
		stmts_196 = stmts_172
		op_197 = op_173
		or__x_198 = or__x_178
		goto b24
	}
b22:
	;
	v276 = v264
	label_stmts_277 = label_stmts_265
	bid_278 = bid_266
	insts_279 = insts_267
	remaining_280 = remaining_268
	out_281 = out_269
	f_282 = f_270
	closed_exprs_283 = closed_exprs_271
	stmts_284 = stmts_272
	op_285 = op_273
	or__x_286 = vm.Boolean(or__x_150)
	goto b19
b23:
	;
	v252 = vm.Boolean(or__x_188)
	label_stmts_253 = label_stmts_179
	bid_254 = bid_180
	insts_255 = insts_181
	remaining_256 = remaining_182
	out_257 = out_183
	f_258 = f_184
	closed_exprs_259 = closed_exprs_185
	stmts_260 = stmts_186
	op_261 = op_187
	or__x_262 = vm.Boolean(or__x_188)
	goto b25
b24:
	;
	or__x_202 = op_197 == vm.Keyword("block-arg")
	if or__x_202 {
		label_stmts_203 = label_stmts_189
		bid_204 = bid_190
		insts_205 = insts_191
		remaining_206 = remaining_192
		out_207 = out_193
		f_208 = f_194
		closed_exprs_209 = closed_exprs_195
		stmts_210 = stmts_196
		op_211 = op_197
		or__x_212 = or__x_202
		goto b26
	} else {
		label_stmts_213 = label_stmts_189
		bid_214 = bid_190
		insts_215 = insts_191
		remaining_216 = remaining_192
		out_217 = out_193
		f_218 = f_194
		closed_exprs_219 = closed_exprs_195
		stmts_220 = stmts_196
		op_221 = op_197
		or__x_222 = or__x_202
		goto b27
	}
b25:
	;
	v264 = v252
	label_stmts_265 = label_stmts_253
	bid_266 = bid_254
	insts_267 = insts_255
	remaining_268 = remaining_256
	out_269 = out_257
	f_270 = f_258
	closed_exprs_271 = closed_exprs_259
	stmts_272 = stmts_260
	op_273 = op_261
	or__x_274 = vm.Boolean(or__x_174)
	goto b22
b26:
	;
	v240 = vm.Boolean(or__x_212)
	label_stmts_241 = label_stmts_203
	bid_242 = bid_204
	insts_243 = insts_205
	remaining_244 = remaining_206
	out_245 = out_207
	f_246 = f_208
	closed_exprs_247 = closed_exprs_209
	stmts_248 = stmts_210
	op_249 = op_211
	or__x_250 = vm.Boolean(or__x_212)
	goto b28
b27:
	;
	arg__15566_226, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{remaining_216})
	if callErr != nil {
		return nil, callErr
	}
	arg__15572_229, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{remaining_216})
	if callErr != nil {
		return nil, callErr
	}
	arg__15573_230, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "used?").Deref(), []vm.Value{f_218, arg__15572_229})
	if callErr != nil {
		return nil, callErr
	}
	arg__15579_233, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{remaining_216})
	if callErr != nil {
		return nil, callErr
	}
	arg__15585_236, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{remaining_216})
	if callErr != nil {
		return nil, callErr
	}
	arg__15586_237, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "used?").Deref(), []vm.Value{f_218, arg__15585_236})
	if callErr != nil {
		return nil, callErr
	}
	v238, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not").Deref(), []vm.Value{arg__15586_237})
	if callErr != nil {
		return nil, callErr
	}
	v240 = v238
	label_stmts_241 = label_stmts_213
	bid_242 = bid_214
	insts_243 = insts_215
	remaining_244 = remaining_216
	out_245 = out_217
	f_246 = f_218
	closed_exprs_247 = closed_exprs_219
	stmts_248 = stmts_220
	op_249 = op_221
	or__x_250 = vm.Boolean(or__x_222)
	goto b28
b28:
	;
	v252 = v240
	label_stmts_253 = label_stmts_241
	bid_254 = bid_242
	insts_255 = insts_243
	remaining_256 = remaining_244
	out_257 = out_245
	f_258 = f_246
	closed_exprs_259 = closed_exprs_247
	stmts_260 = stmts_248
	op_261 = op_249
	or__x_262 = vm.Boolean(or__x_198)
	goto b25
b29:
	;
	arg__15612_395, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "concat").Deref(), []vm.Value{label_stmts_342, inst_stmts_341, term_stmts_349})
	if callErr != nil {
		return nil, callErr
	}
	arg__15621_398, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "concat").Deref(), []vm.Value{label_stmts_342, inst_stmts_341, term_stmts_349})
	if callErr != nil {
		return nil, callErr
	}
	v399, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg__15621_398})
	if callErr != nil {
		return nil, callErr
	}
	v403 = v399
	inst_stmts_404 = inst_stmts_341
	label_stmts_405 = label_stmts_342
	bid_406 = bid_343
	insts_407 = insts_344
	remaining_408 = remaining_345
	out_409 = out_346
	f_410 = f_347
	closed_exprs_411 = closed_exprs_348
	term_stmts_412 = term_stmts_349
	goto b31
b30:
	;
	v403 = vm.NIL
	inst_stmts_404 = inst_stmts_350
	label_stmts_405 = label_stmts_351
	bid_406 = bid_352
	insts_407 = insts_353
	remaining_408 = remaining_354
	out_409 = out_355
	f_410 = f_356
	closed_exprs_411 = closed_exprs_357
	term_stmts_412 = term_stmts_358
	goto b31
b31:
	;
	return v403, nil
b32:
	;
	v382 = term_stmts_368
	and__x_383 = and__x_359
	inst_stmts_384 = inst_stmts_360
	label_stmts_385 = label_stmts_361
	bid_386 = bid_362
	insts_387 = insts_363
	remaining_388 = remaining_364
	out_389 = out_365
	f_390 = f_366
	closed_exprs_391 = closed_exprs_367
	term_stmts_392 = term_stmts_368
	goto b34
b33:
	;
	v382 = and__x_369
	and__x_383 = and__x_369
	inst_stmts_384 = inst_stmts_370
	label_stmts_385 = label_stmts_371
	bid_386 = bid_372
	insts_387 = insts_373
	remaining_388 = remaining_374
	out_389 = out_375
	f_390 = f_376
	closed_exprs_391 = closed_exprs_377
	term_stmts_392 = term_stmts_378
	goto b34
b34:
	;
	if vm.IsTruthy(v382) {
		inst_stmts_341 = inst_stmts_384
		label_stmts_342 = label_stmts_385
		bid_343 = bid_386
		insts_344 = insts_387
		remaining_345 = remaining_388
		out_346 = out_389
		f_347 = f_390
		closed_exprs_348 = closed_exprs_391
		term_stmts_349 = term_stmts_392
		goto b29
	} else {
		inst_stmts_350 = inst_stmts_384
		label_stmts_351 = label_stmts_385
		bid_352 = bid_386
		insts_353 = insts_387
		remaining_354 = remaining_388
		out_355 = out_389
		f_356 = f_390
		closed_exprs_357 = closed_exprs_391
		term_stmts_358 = term_stmts_392
		goto b30
	}
}
func lower_fn_STAR_(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var params_3 vm.Value
	var results_5 vm.Value
	var body_7 vm.Value
	var needs_rt_QMARK__9 vm.Value
	var needs_vm_QMARK__11 vm.Value
	var f_12 vm.Value
	var mode_13 vm.Value
	var params_14 vm.Value
	var results_15 vm.Value
	var body_16 vm.Value
	var needs_rt_QMARK__17 vm.Value
	var needs_vm_QMARK__18 vm.Value
	var arg__15643_33 vm.Value
	var v34 vm.Value
	var f_19 vm.Value
	var mode_20 vm.Value
	var params_21 vm.Value
	var results_22 vm.Value
	var body_23 vm.Value
	var needs_rt_QMARK__24 vm.Value
	var needs_vm_QMARK__25 vm.Value
	var arg__15644_38 vm.Value
	var f_39 vm.Value
	var mode_40 vm.Value
	var params_41 vm.Value
	var results_42 vm.Value
	var body_43 vm.Value
	var needs_rt_QMARK__44 vm.Value
	var needs_vm_QMARK__45 vm.Value
	var arg__15644_46 vm.Value
	var f_47 vm.Value
	var mode_48 vm.Value
	var params_49 vm.Value
	var results_50 vm.Value
	var body_51 vm.Value
	var needs_rt_QMARK__52 vm.Value
	var needs_vm_QMARK__53 vm.Value
	var arg__15651_69 vm.Value
	var v70 vm.Value
	var arg__15644_54 vm.Value
	var f_55 vm.Value
	var mode_56 vm.Value
	var params_57 vm.Value
	var results_58 vm.Value
	var body_59 vm.Value
	var needs_rt_QMARK__60 vm.Value
	var needs_vm_QMARK__61 vm.Value
	var arg__15652_74 vm.Value
	var arg__15644_75 vm.Value
	var f_76 vm.Value
	var mode_77 vm.Value
	var params_78 vm.Value
	var results_79 vm.Value
	var body_80 vm.Value
	var needs_rt_QMARK__81 vm.Value
	var needs_vm_QMARK__82 vm.Value
	var f_84 vm.Value
	var mode_85 vm.Value
	var params_86 vm.Value
	var results_87 vm.Value
	var body_88 vm.Value
	var needs_rt_QMARK__89 vm.Value
	var needs_vm_QMARK__90 vm.Value
	var head__15653_91 vm.Value
	var arg__15660_107 vm.Value
	var v108 vm.Value
	var f_92 vm.Value
	var mode_93 vm.Value
	var params_94 vm.Value
	var results_95 vm.Value
	var body_96 vm.Value
	var needs_rt_QMARK__97 vm.Value
	var needs_vm_QMARK__98 vm.Value
	var head__15653_99 vm.Value
	var arg__15661_112 vm.Value
	var f_113 vm.Value
	var mode_114 vm.Value
	var params_115 vm.Value
	var results_116 vm.Value
	var body_117 vm.Value
	var needs_rt_QMARK__118 vm.Value
	var needs_vm_QMARK__119 vm.Value
	var head__15653_120 vm.Value
	var arg__15661_121 vm.Value
	var f_122 vm.Value
	var mode_123 vm.Value
	var params_124 vm.Value
	var results_125 vm.Value
	var body_126 vm.Value
	var needs_rt_QMARK__127 vm.Value
	var needs_vm_QMARK__128 vm.Value
	var head__15653_129 vm.Value
	var arg__15668_146 vm.Value
	var v147 vm.Value
	var arg__15661_130 vm.Value
	var f_131 vm.Value
	var mode_132 vm.Value
	var params_133 vm.Value
	var results_134 vm.Value
	var body_135 vm.Value
	var needs_rt_QMARK__136 vm.Value
	var needs_vm_QMARK__137 vm.Value
	var head__15653_138 vm.Value
	var arg__15669_151 vm.Value
	var arg__15661_152 vm.Value
	var f_153 vm.Value
	var mode_154 vm.Value
	var params_155 vm.Value
	var results_156 vm.Value
	var body_157 vm.Value
	var needs_rt_QMARK__158 vm.Value
	var needs_vm_QMARK__159 vm.Value
	var head__15653_160 vm.Value
	var arg__15670_161 vm.Value
	var f_163 vm.Value
	var mode_164 vm.Value
	var params_165 vm.Value
	var results_166 vm.Value
	var body_167 vm.Value
	var needs_rt_QMARK__168 vm.Value
	var needs_vm_QMARK__169 vm.Value
	var head__15671_170 vm.Value
	var arg__15678_186 vm.Value
	var v187 vm.Value
	var f_171 vm.Value
	var mode_172 vm.Value
	var params_173 vm.Value
	var results_174 vm.Value
	var body_175 vm.Value
	var needs_rt_QMARK__176 vm.Value
	var needs_vm_QMARK__177 vm.Value
	var head__15671_178 vm.Value
	var arg__15679_191 vm.Value
	var f_192 vm.Value
	var mode_193 vm.Value
	var params_194 vm.Value
	var results_195 vm.Value
	var body_196 vm.Value
	var needs_rt_QMARK__197 vm.Value
	var needs_vm_QMARK__198 vm.Value
	var head__15671_199 vm.Value
	var arg__15679_200 vm.Value
	var f_201 vm.Value
	var mode_202 vm.Value
	var params_203 vm.Value
	var results_204 vm.Value
	var body_205 vm.Value
	var needs_rt_QMARK__206 vm.Value
	var needs_vm_QMARK__207 vm.Value
	var head__15671_208 vm.Value
	var arg__15686_225 vm.Value
	var v226 vm.Value
	var arg__15679_209 vm.Value
	var f_210 vm.Value
	var mode_211 vm.Value
	var params_212 vm.Value
	var results_213 vm.Value
	var body_214 vm.Value
	var needs_rt_QMARK__215 vm.Value
	var needs_vm_QMARK__216 vm.Value
	var head__15671_217 vm.Value
	var arg__15687_230 vm.Value
	var arg__15679_231 vm.Value
	var f_232 vm.Value
	var mode_233 vm.Value
	var params_234 vm.Value
	var results_235 vm.Value
	var body_236 vm.Value
	var needs_rt_QMARK__237 vm.Value
	var needs_vm_QMARK__238 vm.Value
	var head__15671_239 vm.Value
	var f_241 vm.Value
	var mode_242 vm.Value
	var params_243 vm.Value
	var results_244 vm.Value
	var body_245 vm.Value
	var needs_rt_QMARK__246 vm.Value
	var needs_vm_QMARK__247 vm.Value
	var head__15671_248 vm.Value
	var head__15688_249 vm.Value
	var arg__15695_266 vm.Value
	var v267 vm.Value
	var f_250 vm.Value
	var mode_251 vm.Value
	var params_252 vm.Value
	var results_253 vm.Value
	var body_254 vm.Value
	var needs_rt_QMARK__255 vm.Value
	var needs_vm_QMARK__256 vm.Value
	var head__15671_257 vm.Value
	var head__15688_258 vm.Value
	var arg__15696_271 vm.Value
	var f_272 vm.Value
	var mode_273 vm.Value
	var params_274 vm.Value
	var results_275 vm.Value
	var body_276 vm.Value
	var needs_rt_QMARK__277 vm.Value
	var needs_vm_QMARK__278 vm.Value
	var head__15671_279 vm.Value
	var head__15688_280 vm.Value
	var arg__15696_281 vm.Value
	var f_282 vm.Value
	var mode_283 vm.Value
	var params_284 vm.Value
	var results_285 vm.Value
	var body_286 vm.Value
	var needs_rt_QMARK__287 vm.Value
	var needs_vm_QMARK__288 vm.Value
	var head__15671_289 vm.Value
	var head__15688_290 vm.Value
	var arg__15703_308 vm.Value
	var v309 vm.Value
	var arg__15696_291 vm.Value
	var f_292 vm.Value
	var mode_293 vm.Value
	var params_294 vm.Value
	var results_295 vm.Value
	var body_296 vm.Value
	var needs_rt_QMARK__297 vm.Value
	var needs_vm_QMARK__298 vm.Value
	var head__15671_299 vm.Value
	var head__15688_300 vm.Value
	var arg__15704_313 vm.Value
	var arg__15696_314 vm.Value
	var f_315 vm.Value
	var mode_316 vm.Value
	var params_317 vm.Value
	var results_318 vm.Value
	var body_319 vm.Value
	var needs_rt_QMARK__320 vm.Value
	var needs_vm_QMARK__321 vm.Value
	var head__15671_322 vm.Value
	var head__15688_323 vm.Value
	var arg__15705_324 vm.Value
	var imports_325 vm.Value
	var v343 vm.Value
	var f_326 vm.Value
	var mode_327 vm.Value
	var params_328 vm.Value
	var results_329 vm.Value
	var body_330 vm.Value
	var needs_rt_QMARK__331 vm.Value
	var needs_vm_QMARK__332 vm.Value
	var imports_333 vm.Value
	var v348 vm.Value
	var f_334 vm.Value
	var mode_335 vm.Value
	var params_336 vm.Value
	var results_337 vm.Value
	var body_338 vm.Value
	var needs_rt_QMARK__339 vm.Value
	var needs_vm_QMARK__340 vm.Value
	var imports_341 vm.Value
	var v367 vm.Value
	var v469 vm.Value
	var f_470 vm.Value
	var mode_471 vm.Value
	var params_472 vm.Value
	var results_473 vm.Value
	var body_474 vm.Value
	var needs_rt_QMARK__475 vm.Value
	var needs_vm_QMARK__476 vm.Value
	var imports_477 vm.Value
	var f_350 vm.Value
	var mode_351 vm.Value
	var params_352 vm.Value
	var results_353 vm.Value
	var body_354 vm.Value
	var needs_rt_QMARK__355 vm.Value
	var needs_vm_QMARK__356 vm.Value
	var imports_357 vm.Value
	var v372 vm.Value
	var f_358 vm.Value
	var mode_359 vm.Value
	var params_360 vm.Value
	var results_361 vm.Value
	var body_362 vm.Value
	var needs_rt_QMARK__363 vm.Value
	var needs_vm_QMARK__364 vm.Value
	var imports_365 vm.Value
	var v391 vm.Value
	var v459 vm.Value
	var f_460 vm.Value
	var mode_461 vm.Value
	var params_462 vm.Value
	var results_463 vm.Value
	var body_464 vm.Value
	var needs_rt_QMARK__465 vm.Value
	var needs_vm_QMARK__466 vm.Value
	var imports_467 vm.Value
	var f_374 vm.Value
	var mode_375 vm.Value
	var params_376 vm.Value
	var results_377 vm.Value
	var body_378 vm.Value
	var needs_rt_QMARK__379 vm.Value
	var needs_vm_QMARK__380 vm.Value
	var imports_381 vm.Value
	var v396 vm.Value
	var f_382 vm.Value
	var mode_383 vm.Value
	var params_384 vm.Value
	var results_385 vm.Value
	var body_386 vm.Value
	var needs_rt_QMARK__387 vm.Value
	var needs_vm_QMARK__388 vm.Value
	var imports_389 vm.Value
	var v449 vm.Value
	var f_450 vm.Value
	var mode_451 vm.Value
	var params_452 vm.Value
	var results_453 vm.Value
	var body_454 vm.Value
	var needs_rt_QMARK__455 vm.Value
	var needs_vm_QMARK__456 vm.Value
	var imports_457 vm.Value
	var f_398 vm.Value
	var mode_399 vm.Value
	var params_400 vm.Value
	var results_401 vm.Value
	var body_402 vm.Value
	var needs_rt_QMARK__403 vm.Value
	var needs_vm_QMARK__404 vm.Value
	var imports_405 vm.Value
	var arg__15739_422 vm.Value
	var arg__15744_425 vm.Value
	var arg__15745_426 vm.Value
	var arg__15753_429 vm.Value
	var arg__15758_432 vm.Value
	var arg__15759_433 vm.Value
	var arg__15763_434 vm.Value
	var v435 vm.Value
	var f_406 vm.Value
	var mode_407 vm.Value
	var params_408 vm.Value
	var results_409 vm.Value
	var body_410 vm.Value
	var needs_rt_QMARK__411 vm.Value
	var needs_vm_QMARK__412 vm.Value
	var imports_413 vm.Value
	var v439 vm.Value
	var f_440 vm.Value
	var mode_441 vm.Value
	var params_442 vm.Value
	var results_443 vm.Value
	var body_444 vm.Value
	var needs_rt_QMARK__445 vm.Value
	var needs_vm_QMARK__446 vm.Value
	var imports_447 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = params_3, results_5, body_7, needs_rt_QMARK__9, needs_vm_QMARK__11, f_12, mode_13, params_14, results_15, body_16, needs_rt_QMARK__17, needs_vm_QMARK__18, arg__15643_33, v34, f_19, mode_20, params_21, results_22, body_23, needs_rt_QMARK__24, needs_vm_QMARK__25, arg__15644_38, f_39, mode_40, params_41, results_42, body_43, needs_rt_QMARK__44, needs_vm_QMARK__45, arg__15644_46, f_47, mode_48, params_49, results_50, body_51, needs_rt_QMARK__52, needs_vm_QMARK__53, arg__15651_69, v70, arg__15644_54, f_55, mode_56, params_57, results_58, body_59, needs_rt_QMARK__60, needs_vm_QMARK__61, arg__15652_74, arg__15644_75, f_76, mode_77, params_78, results_79, body_80, needs_rt_QMARK__81, needs_vm_QMARK__82, f_84, mode_85, params_86, results_87, body_88, needs_rt_QMARK__89, needs_vm_QMARK__90, head__15653_91, arg__15660_107, v108, f_92, mode_93, params_94, results_95, body_96, needs_rt_QMARK__97, needs_vm_QMARK__98, head__15653_99, arg__15661_112, f_113, mode_114, params_115, results_116, body_117, needs_rt_QMARK__118, needs_vm_QMARK__119, head__15653_120, arg__15661_121, f_122, mode_123, params_124, results_125, body_126, needs_rt_QMARK__127, needs_vm_QMARK__128, head__15653_129, arg__15668_146, v147, arg__15661_130, f_131, mode_132, params_133, results_134, body_135, needs_rt_QMARK__136, needs_vm_QMARK__137, head__15653_138, arg__15669_151, arg__15661_152, f_153, mode_154, params_155, results_156, body_157, needs_rt_QMARK__158, needs_vm_QMARK__159, head__15653_160, arg__15670_161, f_163, mode_164, params_165, results_166, body_167, needs_rt_QMARK__168, needs_vm_QMARK__169, head__15671_170, arg__15678_186, v187, f_171, mode_172, params_173, results_174, body_175, needs_rt_QMARK__176, needs_vm_QMARK__177, head__15671_178, arg__15679_191, f_192, mode_193, params_194, results_195, body_196, needs_rt_QMARK__197, needs_vm_QMARK__198, head__15671_199, arg__15679_200, f_201, mode_202, params_203, results_204, body_205, needs_rt_QMARK__206, needs_vm_QMARK__207, head__15671_208, arg__15686_225, v226, arg__15679_209, f_210, mode_211, params_212, results_213, body_214, needs_rt_QMARK__215, needs_vm_QMARK__216, head__15671_217, arg__15687_230, arg__15679_231, f_232, mode_233, params_234, results_235, body_236, needs_rt_QMARK__237, needs_vm_QMARK__238, head__15671_239, f_241, mode_242, params_243, results_244, body_245, needs_rt_QMARK__246, needs_vm_QMARK__247, head__15671_248, head__15688_249, arg__15695_266, v267, f_250, mode_251, params_252, results_253, body_254, needs_rt_QMARK__255, needs_vm_QMARK__256, head__15671_257, head__15688_258, arg__15696_271, f_272, mode_273, params_274, results_275, body_276, needs_rt_QMARK__277, needs_vm_QMARK__278, head__15671_279, head__15688_280, arg__15696_281, f_282, mode_283, params_284, results_285, body_286, needs_rt_QMARK__287, needs_vm_QMARK__288, head__15671_289, head__15688_290, arg__15703_308, v309, arg__15696_291, f_292, mode_293, params_294, results_295, body_296, needs_rt_QMARK__297, needs_vm_QMARK__298, head__15671_299, head__15688_300, arg__15704_313, arg__15696_314, f_315, mode_316, params_317, results_318, body_319, needs_rt_QMARK__320, needs_vm_QMARK__321, head__15671_322, head__15688_323, arg__15705_324, imports_325, v343, f_326, mode_327, params_328, results_329, body_330, needs_rt_QMARK__331, needs_vm_QMARK__332, imports_333, v348, f_334, mode_335, params_336, results_337, body_338, needs_rt_QMARK__339, needs_vm_QMARK__340, imports_341, v367, v469, f_470, mode_471, params_472, results_473, body_474, needs_rt_QMARK__475, needs_vm_QMARK__476, imports_477, f_350, mode_351, params_352, results_353, body_354, needs_rt_QMARK__355, needs_vm_QMARK__356, imports_357, v372, f_358, mode_359, params_360, results_361, body_362, needs_rt_QMARK__363, needs_vm_QMARK__364, imports_365, v391, v459, f_460, mode_461, params_462, results_463, body_464, needs_rt_QMARK__465, needs_vm_QMARK__466, imports_467, f_374, mode_375, params_376, results_377, body_378, needs_rt_QMARK__379, needs_vm_QMARK__380, imports_381, v396, f_382, mode_383, params_384, results_385, body_386, needs_rt_QMARK__387, needs_vm_QMARK__388, imports_389, v449, f_450, mode_451, params_452, results_453, body_454, needs_rt_QMARK__455, needs_vm_QMARK__456, imports_457, f_398, mode_399, params_400, results_401, body_402, needs_rt_QMARK__403, needs_vm_QMARK__404, imports_405, arg__15739_422, arg__15744_425, arg__15745_426, arg__15753_429, arg__15758_432, arg__15759_433, arg__15763_434, v435, f_406, mode_407, params_408, results_409, body_410, needs_rt_QMARK__411, needs_vm_QMARK__412, imports_413, v439, f_440, mode_441, params_442, results_443, body_444, needs_rt_QMARK__445, needs_vm_QMARK__446, imports_447
	params_3, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "params-for").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	results_5, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "result-node").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	body_7, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "lower-body").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	needs_rt_QMARK__9, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "function-needs-rt?").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	needs_vm_QMARK__11, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "function-needs-vm?").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(needs_rt_QMARK__9) {
		f_12 = arg0
		mode_13 = arg1
		params_14 = params_3
		results_15 = results_5
		body_16 = body_7
		needs_rt_QMARK__17 = needs_rt_QMARK__9
		needs_vm_QMARK__18 = needs_vm_QMARK__11
		goto b1
	} else {
		f_19 = arg0
		mode_20 = arg1
		params_21 = params_3
		results_22 = results_5
		body_23 = body_7
		needs_rt_QMARK__24 = needs_rt_QMARK__9
		needs_vm_QMARK__25 = needs_vm_QMARK__11
		goto b2
	}
b1:
	;
	arg__15643_33, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "import-entry").Deref(), []vm.Value{vm.String("github.com/nooga/let-go/pkg/rt"), vm.String("rt")})
	if callErr != nil {
		return nil, callErr
	}
	v34, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__15643_33})
	if callErr != nil {
		return nil, callErr
	}
	arg__15644_38 = v34
	f_39 = f_12
	mode_40 = mode_13
	params_41 = params_14
	results_42 = results_15
	body_43 = body_16
	needs_rt_QMARK__44 = needs_rt_QMARK__17
	needs_vm_QMARK__45 = needs_vm_QMARK__18
	goto b3
b2:
	;
	arg__15644_38 = vm.NewArrayVector([]vm.Value{})
	f_39 = f_19
	mode_40 = mode_20
	params_41 = params_21
	results_42 = results_22
	body_43 = body_23
	needs_rt_QMARK__44 = needs_rt_QMARK__24
	needs_vm_QMARK__45 = needs_vm_QMARK__25
	goto b3
b3:
	;
	if vm.IsTruthy(needs_vm_QMARK__45) {
		arg__15644_46 = arg__15644_38
		f_47 = f_39
		mode_48 = mode_40
		params_49 = params_41
		results_50 = results_42
		body_51 = body_43
		needs_rt_QMARK__52 = needs_rt_QMARK__44
		needs_vm_QMARK__53 = needs_vm_QMARK__45
		goto b4
	} else {
		arg__15644_54 = arg__15644_38
		f_55 = f_39
		mode_56 = mode_40
		params_57 = params_41
		results_58 = results_42
		body_59 = body_43
		needs_rt_QMARK__60 = needs_rt_QMARK__44
		needs_vm_QMARK__61 = needs_vm_QMARK__45
		goto b5
	}
b4:
	;
	arg__15651_69, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "import-entry").Deref(), []vm.Value{vm.String("github.com/nooga/let-go/pkg/vm"), vm.String("vm")})
	if callErr != nil {
		return nil, callErr
	}
	v70, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__15651_69})
	if callErr != nil {
		return nil, callErr
	}
	arg__15652_74 = v70
	arg__15644_75 = arg__15644_46
	f_76 = f_47
	mode_77 = mode_48
	params_78 = params_49
	results_79 = results_50
	body_80 = body_51
	needs_rt_QMARK__81 = needs_rt_QMARK__52
	needs_vm_QMARK__82 = needs_vm_QMARK__53
	goto b6
b5:
	;
	arg__15652_74 = vm.NewArrayVector([]vm.Value{})
	arg__15644_75 = arg__15644_54
	f_76 = f_55
	mode_77 = mode_56
	params_78 = params_57
	results_79 = results_58
	body_80 = body_59
	needs_rt_QMARK__81 = needs_rt_QMARK__60
	needs_vm_QMARK__82 = needs_vm_QMARK__61
	goto b6
b6:
	;
	if vm.IsTruthy(needs_rt_QMARK__81) {
		f_84 = f_76
		mode_85 = mode_77
		params_86 = params_78
		results_87 = results_79
		body_88 = body_80
		needs_rt_QMARK__89 = needs_rt_QMARK__81
		needs_vm_QMARK__90 = needs_vm_QMARK__82
		head__15653_91 = rt.LookupVar("clojure.core", "concat").Deref()
		goto b7
	} else {
		f_92 = f_76
		mode_93 = mode_77
		params_94 = params_78
		results_95 = results_79
		body_96 = body_80
		needs_rt_QMARK__97 = needs_rt_QMARK__81
		needs_vm_QMARK__98 = needs_vm_QMARK__82
		head__15653_99 = rt.LookupVar("clojure.core", "concat").Deref()
		goto b8
	}
b7:
	;
	arg__15660_107, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "import-entry").Deref(), []vm.Value{vm.String("github.com/nooga/let-go/pkg/rt"), vm.String("rt")})
	if callErr != nil {
		return nil, callErr
	}
	v108, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__15660_107})
	if callErr != nil {
		return nil, callErr
	}
	arg__15661_112 = v108
	f_113 = f_84
	mode_114 = mode_85
	params_115 = params_86
	results_116 = results_87
	body_117 = body_88
	needs_rt_QMARK__118 = needs_rt_QMARK__89
	needs_vm_QMARK__119 = needs_vm_QMARK__90
	head__15653_120 = head__15653_91
	goto b9
b8:
	;
	arg__15661_112 = vm.NewArrayVector([]vm.Value{})
	f_113 = f_92
	mode_114 = mode_93
	params_115 = params_94
	results_116 = results_95
	body_117 = body_96
	needs_rt_QMARK__118 = needs_rt_QMARK__97
	needs_vm_QMARK__119 = needs_vm_QMARK__98
	head__15653_120 = head__15653_99
	goto b9
b9:
	;
	if vm.IsTruthy(needs_vm_QMARK__119) {
		arg__15661_121 = arg__15661_112
		f_122 = f_113
		mode_123 = mode_114
		params_124 = params_115
		results_125 = results_116
		body_126 = body_117
		needs_rt_QMARK__127 = needs_rt_QMARK__118
		needs_vm_QMARK__128 = needs_vm_QMARK__119
		head__15653_129 = head__15653_120
		goto b10
	} else {
		arg__15661_130 = arg__15661_112
		f_131 = f_113
		mode_132 = mode_114
		params_133 = params_115
		results_134 = results_116
		body_135 = body_117
		needs_rt_QMARK__136 = needs_rt_QMARK__118
		needs_vm_QMARK__137 = needs_vm_QMARK__119
		head__15653_138 = head__15653_120
		goto b11
	}
b10:
	;
	arg__15668_146, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "import-entry").Deref(), []vm.Value{vm.String("github.com/nooga/let-go/pkg/vm"), vm.String("vm")})
	if callErr != nil {
		return nil, callErr
	}
	v147, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__15668_146})
	if callErr != nil {
		return nil, callErr
	}
	arg__15669_151 = v147
	arg__15661_152 = arg__15661_121
	f_153 = f_122
	mode_154 = mode_123
	params_155 = params_124
	results_156 = results_125
	body_157 = body_126
	needs_rt_QMARK__158 = needs_rt_QMARK__127
	needs_vm_QMARK__159 = needs_vm_QMARK__128
	head__15653_160 = head__15653_129
	goto b12
b11:
	;
	arg__15669_151 = vm.NewArrayVector([]vm.Value{})
	arg__15661_152 = arg__15661_130
	f_153 = f_131
	mode_154 = mode_132
	params_155 = params_133
	results_156 = results_134
	body_157 = body_135
	needs_rt_QMARK__158 = needs_rt_QMARK__136
	needs_vm_QMARK__159 = needs_vm_QMARK__137
	head__15653_160 = head__15653_138
	goto b12
b12:
	;
	arg__15670_161, callErr = rt.InvokeValue(head__15653_160, []vm.Value{arg__15661_152, arg__15669_151})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(needs_rt_QMARK__158) {
		f_163 = f_153
		mode_164 = mode_154
		params_165 = params_155
		results_166 = results_156
		body_167 = body_157
		needs_rt_QMARK__168 = needs_rt_QMARK__158
		needs_vm_QMARK__169 = needs_vm_QMARK__159
		head__15671_170 = rt.LookupVar("clojure.core", "vec").Deref()
		goto b13
	} else {
		f_171 = f_153
		mode_172 = mode_154
		params_173 = params_155
		results_174 = results_156
		body_175 = body_157
		needs_rt_QMARK__176 = needs_rt_QMARK__158
		needs_vm_QMARK__177 = needs_vm_QMARK__159
		head__15671_178 = rt.LookupVar("clojure.core", "vec").Deref()
		goto b14
	}
b13:
	;
	arg__15678_186, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "import-entry").Deref(), []vm.Value{vm.String("github.com/nooga/let-go/pkg/rt"), vm.String("rt")})
	if callErr != nil {
		return nil, callErr
	}
	v187, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__15678_186})
	if callErr != nil {
		return nil, callErr
	}
	arg__15679_191 = v187
	f_192 = f_163
	mode_193 = mode_164
	params_194 = params_165
	results_195 = results_166
	body_196 = body_167
	needs_rt_QMARK__197 = needs_rt_QMARK__168
	needs_vm_QMARK__198 = needs_vm_QMARK__169
	head__15671_199 = head__15671_170
	goto b15
b14:
	;
	arg__15679_191 = vm.NewArrayVector([]vm.Value{})
	f_192 = f_171
	mode_193 = mode_172
	params_194 = params_173
	results_195 = results_174
	body_196 = body_175
	needs_rt_QMARK__197 = needs_rt_QMARK__176
	needs_vm_QMARK__198 = needs_vm_QMARK__177
	head__15671_199 = head__15671_178
	goto b15
b15:
	;
	if vm.IsTruthy(needs_vm_QMARK__198) {
		arg__15679_200 = arg__15679_191
		f_201 = f_192
		mode_202 = mode_193
		params_203 = params_194
		results_204 = results_195
		body_205 = body_196
		needs_rt_QMARK__206 = needs_rt_QMARK__197
		needs_vm_QMARK__207 = needs_vm_QMARK__198
		head__15671_208 = head__15671_199
		goto b16
	} else {
		arg__15679_209 = arg__15679_191
		f_210 = f_192
		mode_211 = mode_193
		params_212 = params_194
		results_213 = results_195
		body_214 = body_196
		needs_rt_QMARK__215 = needs_rt_QMARK__197
		needs_vm_QMARK__216 = needs_vm_QMARK__198
		head__15671_217 = head__15671_199
		goto b17
	}
b16:
	;
	arg__15686_225, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "import-entry").Deref(), []vm.Value{vm.String("github.com/nooga/let-go/pkg/vm"), vm.String("vm")})
	if callErr != nil {
		return nil, callErr
	}
	v226, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__15686_225})
	if callErr != nil {
		return nil, callErr
	}
	arg__15687_230 = v226
	arg__15679_231 = arg__15679_200
	f_232 = f_201
	mode_233 = mode_202
	params_234 = params_203
	results_235 = results_204
	body_236 = body_205
	needs_rt_QMARK__237 = needs_rt_QMARK__206
	needs_vm_QMARK__238 = needs_vm_QMARK__207
	head__15671_239 = head__15671_208
	goto b18
b17:
	;
	arg__15687_230 = vm.NewArrayVector([]vm.Value{})
	arg__15679_231 = arg__15679_209
	f_232 = f_210
	mode_233 = mode_211
	params_234 = params_212
	results_235 = results_213
	body_236 = body_214
	needs_rt_QMARK__237 = needs_rt_QMARK__215
	needs_vm_QMARK__238 = needs_vm_QMARK__216
	head__15671_239 = head__15671_217
	goto b18
b18:
	;
	if vm.IsTruthy(needs_rt_QMARK__237) {
		f_241 = f_232
		mode_242 = mode_233
		params_243 = params_234
		results_244 = results_235
		body_245 = body_236
		needs_rt_QMARK__246 = needs_rt_QMARK__237
		needs_vm_QMARK__247 = needs_vm_QMARK__238
		head__15671_248 = head__15671_239
		head__15688_249 = rt.LookupVar("clojure.core", "concat").Deref()
		goto b19
	} else {
		f_250 = f_232
		mode_251 = mode_233
		params_252 = params_234
		results_253 = results_235
		body_254 = body_236
		needs_rt_QMARK__255 = needs_rt_QMARK__237
		needs_vm_QMARK__256 = needs_vm_QMARK__238
		head__15671_257 = head__15671_239
		head__15688_258 = rt.LookupVar("clojure.core", "concat").Deref()
		goto b20
	}
b19:
	;
	arg__15695_266, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "import-entry").Deref(), []vm.Value{vm.String("github.com/nooga/let-go/pkg/rt"), vm.String("rt")})
	if callErr != nil {
		return nil, callErr
	}
	v267, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__15695_266})
	if callErr != nil {
		return nil, callErr
	}
	arg__15696_271 = v267
	f_272 = f_241
	mode_273 = mode_242
	params_274 = params_243
	results_275 = results_244
	body_276 = body_245
	needs_rt_QMARK__277 = needs_rt_QMARK__246
	needs_vm_QMARK__278 = needs_vm_QMARK__247
	head__15671_279 = head__15671_248
	head__15688_280 = head__15688_249
	goto b21
b20:
	;
	arg__15696_271 = vm.NewArrayVector([]vm.Value{})
	f_272 = f_250
	mode_273 = mode_251
	params_274 = params_252
	results_275 = results_253
	body_276 = body_254
	needs_rt_QMARK__277 = needs_rt_QMARK__255
	needs_vm_QMARK__278 = needs_vm_QMARK__256
	head__15671_279 = head__15671_257
	head__15688_280 = head__15688_258
	goto b21
b21:
	;
	if vm.IsTruthy(needs_vm_QMARK__278) {
		arg__15696_281 = arg__15696_271
		f_282 = f_272
		mode_283 = mode_273
		params_284 = params_274
		results_285 = results_275
		body_286 = body_276
		needs_rt_QMARK__287 = needs_rt_QMARK__277
		needs_vm_QMARK__288 = needs_vm_QMARK__278
		head__15671_289 = head__15671_279
		head__15688_290 = head__15688_280
		goto b22
	} else {
		arg__15696_291 = arg__15696_271
		f_292 = f_272
		mode_293 = mode_273
		params_294 = params_274
		results_295 = results_275
		body_296 = body_276
		needs_rt_QMARK__297 = needs_rt_QMARK__277
		needs_vm_QMARK__298 = needs_vm_QMARK__278
		head__15671_299 = head__15671_279
		head__15688_300 = head__15688_280
		goto b23
	}
b22:
	;
	arg__15703_308, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "import-entry").Deref(), []vm.Value{vm.String("github.com/nooga/let-go/pkg/vm"), vm.String("vm")})
	if callErr != nil {
		return nil, callErr
	}
	v309, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__15703_308})
	if callErr != nil {
		return nil, callErr
	}
	arg__15704_313 = v309
	arg__15696_314 = arg__15696_281
	f_315 = f_282
	mode_316 = mode_283
	params_317 = params_284
	results_318 = results_285
	body_319 = body_286
	needs_rt_QMARK__320 = needs_rt_QMARK__287
	needs_vm_QMARK__321 = needs_vm_QMARK__288
	head__15671_322 = head__15671_289
	head__15688_323 = head__15688_290
	goto b24
b23:
	;
	arg__15704_313 = vm.NewArrayVector([]vm.Value{})
	arg__15696_314 = arg__15696_291
	f_315 = f_292
	mode_316 = mode_293
	params_317 = params_294
	results_318 = results_295
	body_319 = body_296
	needs_rt_QMARK__320 = needs_rt_QMARK__297
	needs_vm_QMARK__321 = needs_vm_QMARK__298
	head__15671_322 = head__15671_299
	head__15688_323 = head__15688_300
	goto b24
b24:
	;
	arg__15705_324, callErr = rt.InvokeValue(head__15688_323, []vm.Value{arg__15696_314, arg__15704_313})
	if callErr != nil {
		return nil, callErr
	}
	imports_325, callErr = rt.InvokeValue(head__15671_322, []vm.Value{arg__15705_324})
	if callErr != nil {
		return nil, callErr
	}
	v343, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{params_234})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v343) {
		f_326 = f_232
		mode_327 = mode_233
		params_328 = params_234
		results_329 = results_235
		body_330 = body_236
		needs_rt_QMARK__331 = needs_rt_QMARK__237
		needs_vm_QMARK__332 = needs_vm_QMARK__238
		imports_333 = imports_325
		goto b25
	} else {
		f_334 = f_232
		mode_335 = mode_233
		params_336 = params_234
		results_337 = results_235
		body_338 = body_236
		needs_rt_QMARK__339 = needs_rt_QMARK__237
		needs_vm_QMARK__340 = needs_vm_QMARK__238
		imports_341 = imports_325
		goto b26
	}
b25:
	;
	v348, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "unsupported").Deref(), []vm.Value{mode_327, vm.String("unsupported parameter types")})
	if callErr != nil {
		return nil, callErr
	}
	v469 = v348
	f_470 = f_326
	mode_471 = mode_327
	params_472 = params_328
	results_473 = results_329
	body_474 = body_330
	needs_rt_QMARK__475 = needs_rt_QMARK__331
	needs_vm_QMARK__476 = needs_vm_QMARK__332
	imports_477 = imports_333
	goto b27
b26:
	;
	v367, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{results_337})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v367) {
		f_350 = f_334
		mode_351 = mode_335
		params_352 = params_336
		results_353 = results_337
		body_354 = body_338
		needs_rt_QMARK__355 = needs_rt_QMARK__339
		needs_vm_QMARK__356 = needs_vm_QMARK__340
		imports_357 = imports_341
		goto b28
	} else {
		f_358 = f_334
		mode_359 = mode_335
		params_360 = params_336
		results_361 = results_337
		body_362 = body_338
		needs_rt_QMARK__363 = needs_rt_QMARK__339
		needs_vm_QMARK__364 = needs_vm_QMARK__340
		imports_365 = imports_341
		goto b29
	}
b27:
	;
	return v469, nil
b28:
	;
	v372, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "unsupported").Deref(), []vm.Value{mode_351, vm.String("unsupported result type")})
	if callErr != nil {
		return nil, callErr
	}
	v459 = v372
	f_460 = f_350
	mode_461 = mode_351
	params_462 = params_352
	results_463 = results_353
	body_464 = body_354
	needs_rt_QMARK__465 = needs_rt_QMARK__355
	needs_vm_QMARK__466 = needs_vm_QMARK__356
	imports_467 = imports_357
	goto b30
b29:
	;
	v391, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{body_362})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v391) {
		f_374 = f_358
		mode_375 = mode_359
		params_376 = params_360
		results_377 = results_361
		body_378 = body_362
		needs_rt_QMARK__379 = needs_rt_QMARK__363
		needs_vm_QMARK__380 = needs_vm_QMARK__364
		imports_381 = imports_365
		goto b31
	} else {
		f_382 = f_358
		mode_383 = mode_359
		params_384 = params_360
		results_385 = results_361
		body_386 = body_362
		needs_rt_QMARK__387 = needs_rt_QMARK__363
		needs_vm_QMARK__388 = needs_vm_QMARK__364
		imports_389 = imports_365
		goto b32
	}
b30:
	;
	v469 = v459
	f_470 = f_460
	mode_471 = mode_461
	params_472 = params_462
	results_473 = results_463
	body_474 = body_464
	needs_rt_QMARK__475 = needs_rt_QMARK__465
	needs_vm_QMARK__476 = needs_vm_QMARK__466
	imports_477 = imports_467
	goto b27
b31:
	;
	v396, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "unsupported").Deref(), []vm.Value{mode_375, vm.String("unsupported function body shape")})
	if callErr != nil {
		return nil, callErr
	}
	v449 = v396
	f_450 = f_374
	mode_451 = mode_375
	params_452 = params_376
	results_453 = results_377
	body_454 = body_378
	needs_rt_QMARK__455 = needs_rt_QMARK__379
	needs_vm_QMARK__456 = needs_vm_QMARK__380
	imports_457 = imports_381
	goto b33
b32:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		f_398 = f_382
		mode_399 = mode_383
		params_400 = params_384
		results_401 = results_385
		body_402 = body_386
		needs_rt_QMARK__403 = needs_rt_QMARK__387
		needs_vm_QMARK__404 = needs_vm_QMARK__388
		imports_405 = imports_389
		goto b34
	} else {
		f_406 = f_382
		mode_407 = mode_383
		params_408 = params_384
		results_409 = results_385
		body_410 = body_386
		needs_rt_QMARK__411 = needs_rt_QMARK__387
		needs_vm_QMARK__412 = needs_vm_QMARK__388
		imports_413 = imports_389
		goto b35
	}
b33:
	;
	v459 = v449
	f_460 = f_450
	mode_461 = mode_451
	params_462 = params_452
	results_463 = results_453
	body_464 = body_454
	needs_rt_QMARK__465 = needs_rt_QMARK__455
	needs_vm_QMARK__466 = needs_vm_QMARK__456
	imports_467 = imports_457
	goto b30
b34:
	;
	arg__15739_422, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-name").Deref(), []vm.Value{f_398})
	if callErr != nil {
		return nil, callErr
	}
	arg__15744_425, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-name").Deref(), []vm.Value{f_398})
	if callErr != nil {
		return nil, callErr
	}
	arg__15745_426, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "go-name").Deref(), []vm.Value{arg__15744_425})
	if callErr != nil {
		return nil, callErr
	}
	arg__15753_429, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-name").Deref(), []vm.Value{f_398})
	if callErr != nil {
		return nil, callErr
	}
	arg__15758_432, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-name").Deref(), []vm.Value{f_398})
	if callErr != nil {
		return nil, callErr
	}
	arg__15759_433, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "go-name").Deref(), []vm.Value{arg__15758_432})
	if callErr != nil {
		return nil, callErr
	}
	arg__15763_434, callErr = rt.InvokeValue(rt.LookupVar("gogen", "func-decl").Deref(), []vm.Value{arg__15759_433, params_400, results_401, body_402})
	if callErr != nil {
		return nil, callErr
	}
	v435, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("imports"), imports_405, vm.Keyword("status"), vm.Keyword("lowered"), vm.Keyword("decl"), arg__15763_434})
	if callErr != nil {
		return nil, callErr
	}
	v439 = v435
	f_440 = f_398
	mode_441 = mode_399
	params_442 = params_400
	results_443 = results_401
	body_444 = body_402
	needs_rt_QMARK__445 = needs_rt_QMARK__403
	needs_vm_QMARK__446 = needs_vm_QMARK__404
	imports_447 = imports_405
	goto b36
b35:
	;
	v439 = vm.NIL
	f_440 = f_406
	mode_441 = mode_407
	params_442 = params_408
	results_443 = results_409
	body_444 = body_410
	needs_rt_QMARK__445 = needs_rt_QMARK__411
	needs_vm_QMARK__446 = needs_vm_QMARK__412
	imports_447 = imports_413
	goto b36
b36:
	;
	v449 = v439
	f_450 = f_440
	mode_451 = mode_441
	params_452 = params_442
	results_453 = results_443
	body_454 = body_444
	needs_rt_QMARK__455 = needs_rt_QMARK__445
	needs_vm_QMARK__456 = needs_vm_QMARK__446
	imports_457 = imports_447
	goto b33
}
func lower_fn_lit(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var params_3 vm.Value
	var results_5 vm.Value
	var body_7 vm.Value
	var f_8 vm.Value
	var closed_exprs_9 vm.Value
	var params_10 vm.Value
	var results_11 vm.Value
	var body_12 vm.Value
	var v64 vm.Value
	var f_13 vm.Value
	var closed_exprs_14 vm.Value
	var params_15 vm.Value
	var results_16 vm.Value
	var body_17 vm.Value
	var v68 vm.Value
	var f_69 vm.Value
	var closed_exprs_70 vm.Value
	var params_71 vm.Value
	var results_72 vm.Value
	var body_73 vm.Value
	var f_18 vm.Value
	var closed_exprs_19 vm.Value
	var and__x_20 vm.Value
	var params_21 vm.Value
	var results_22 vm.Value
	var body_23 vm.Value
	var f_24 vm.Value
	var closed_exprs_25 vm.Value
	var and__x_26 vm.Value
	var params_27 vm.Value
	var results_28 vm.Value
	var body_29 vm.Value
	var v55 vm.Value
	var f_56 vm.Value
	var closed_exprs_57 vm.Value
	var and__x_58 vm.Value
	var params_59 vm.Value
	var results_60 vm.Value
	var body_61 vm.Value
	var f_31 vm.Value
	var closed_exprs_32 vm.Value
	var params_33 vm.Value
	var and__x_34 vm.Value
	var results_35 vm.Value
	var body_36 vm.Value
	var f_37 vm.Value
	var closed_exprs_38 vm.Value
	var params_39 vm.Value
	var and__x_40 vm.Value
	var results_41 vm.Value
	var body_42 vm.Value
	var v46 vm.Value
	var f_47 vm.Value
	var closed_exprs_48 vm.Value
	var params_49 vm.Value
	var and__x_50 vm.Value
	var results_51 vm.Value
	var body_52 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = params_3, results_5, body_7, f_8, closed_exprs_9, params_10, results_11, body_12, v64, f_13, closed_exprs_14, params_15, results_16, body_17, v68, f_69, closed_exprs_70, params_71, results_72, body_73, f_18, closed_exprs_19, and__x_20, params_21, results_22, body_23, f_24, closed_exprs_25, and__x_26, params_27, results_28, body_29, v55, f_56, closed_exprs_57, and__x_58, params_59, results_60, body_61, f_31, closed_exprs_32, params_33, and__x_34, results_35, body_36, f_37, closed_exprs_38, params_39, and__x_40, results_41, body_42, v46, f_47, closed_exprs_48, params_49, and__x_50, results_51, body_52
	params_3, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "params-for").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	results_5, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "result-node").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	body_7, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "lower-body").Deref(), []vm.Value{arg0, arg1})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(params_3) {
		f_18 = arg0
		closed_exprs_19 = arg1
		and__x_20 = params_3
		params_21 = params_3
		results_22 = results_5
		body_23 = body_7
		goto b4
	} else {
		f_24 = arg0
		closed_exprs_25 = arg1
		and__x_26 = params_3
		params_27 = params_3
		results_28 = results_5
		body_29 = body_7
		goto b5
	}
b1:
	;
	v64, callErr = rt.InvokeValue(rt.LookupVar("gogen", "func-lit").Deref(), []vm.Value{params_10, results_11, body_12})
	if callErr != nil {
		return nil, callErr
	}
	v68 = v64
	f_69 = f_8
	closed_exprs_70 = closed_exprs_9
	params_71 = params_10
	results_72 = results_11
	body_73 = body_12
	goto b3
b2:
	;
	v68 = vm.NIL
	f_69 = f_13
	closed_exprs_70 = closed_exprs_14
	params_71 = params_15
	results_72 = results_16
	body_73 = body_17
	goto b3
b3:
	;
	return v68, nil
b4:
	;
	if vm.IsTruthy(results_22) {
		f_31 = f_18
		closed_exprs_32 = closed_exprs_19
		params_33 = params_21
		and__x_34 = results_22
		results_35 = results_22
		body_36 = body_23
		goto b7
	} else {
		f_37 = f_18
		closed_exprs_38 = closed_exprs_19
		params_39 = params_21
		and__x_40 = results_22
		results_41 = results_22
		body_42 = body_23
		goto b8
	}
b5:
	;
	v55 = and__x_26
	f_56 = f_24
	closed_exprs_57 = closed_exprs_25
	and__x_58 = and__x_26
	params_59 = params_27
	results_60 = results_28
	body_61 = body_29
	goto b6
b6:
	;
	if vm.IsTruthy(v55) {
		f_8 = f_56
		closed_exprs_9 = closed_exprs_57
		params_10 = params_59
		results_11 = results_60
		body_12 = body_61
		goto b1
	} else {
		f_13 = f_56
		closed_exprs_14 = closed_exprs_57
		params_15 = params_59
		results_16 = results_60
		body_17 = body_61
		goto b2
	}
b7:
	;
	v46 = body_36
	f_47 = f_31
	closed_exprs_48 = closed_exprs_32
	params_49 = params_33
	and__x_50 = and__x_34
	results_51 = results_35
	body_52 = body_36
	goto b9
b8:
	;
	v46 = and__x_40
	f_47 = f_37
	closed_exprs_48 = closed_exprs_38
	params_49 = params_39
	and__x_50 = and__x_40
	results_51 = results_41
	body_52 = body_42
	goto b9
b9:
	;
	v55 = v46
	f_56 = f_47
	closed_exprs_57 = closed_exprs_48
	and__x_58 = and__x_20
	params_59 = params_49
	results_60 = results_51
	body_61 = body_52
	goto b6
}
func lower_inst_stmts(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var op_4 vm.Value
	var aux_6 vm.Value
	var v18 bool
	var f_7 vm.Value
	var closed_exprs_8 vm.Value
	var nid_9 vm.Value
	var op_10 vm.Value
	var aux_11 vm.Value
	var f_12 vm.Value
	var closed_exprs_13 vm.Value
	var nid_14 vm.Value
	var op_15 vm.Value
	var aux_16 vm.Value
	var v33 bool
	var v428 vm.Value
	var f_429 vm.Value
	var closed_exprs_430 vm.Value
	var nid_431 vm.Value
	var op_432 vm.Value
	var aux_433 vm.Value
	var f_22 vm.Value
	var closed_exprs_23 vm.Value
	var nid_24 vm.Value
	var op_25 vm.Value
	var aux_26 vm.Value
	var f_27 vm.Value
	var closed_exprs_28 vm.Value
	var nid_29 vm.Value
	var op_30 vm.Value
	var aux_31 vm.Value
	var v48 bool
	var v421 vm.Value
	var f_422 vm.Value
	var closed_exprs_423 vm.Value
	var nid_424 vm.Value
	var op_425 vm.Value
	var aux_426 vm.Value
	var f_37 vm.Value
	var closed_exprs_38 vm.Value
	var nid_39 vm.Value
	var op_40 vm.Value
	var aux_41 vm.Value
	var f_42 vm.Value
	var closed_exprs_43 vm.Value
	var nid_44 vm.Value
	var op_45 vm.Value
	var aux_46 vm.Value
	var and__x_63 bool
	var v414 vm.Value
	var f_415 vm.Value
	var closed_exprs_416 vm.Value
	var nid_417 vm.Value
	var op_418 vm.Value
	var aux_419 vm.Value
	var f_52 vm.Value
	var closed_exprs_53 vm.Value
	var nid_54 vm.Value
	var op_55 vm.Value
	var aux_56 vm.Value
	var f_57 vm.Value
	var closed_exprs_58 vm.Value
	var nid_59 vm.Value
	var op_60 vm.Value
	var aux_61 vm.Value
	var v102 bool
	var v407 vm.Value
	var f_408 vm.Value
	var closed_exprs_409 vm.Value
	var nid_410 vm.Value
	var op_411 vm.Value
	var aux_412 vm.Value
	var f_64 vm.Value
	var closed_exprs_65 vm.Value
	var nid_66 vm.Value
	var op_67 vm.Value
	var aux_68 vm.Value
	var and__x_69 bool
	var v78 vm.Value
	var f_70 vm.Value
	var closed_exprs_71 vm.Value
	var nid_72 vm.Value
	var op_73 vm.Value
	var aux_74 vm.Value
	var and__x_75 bool
	var v81 vm.Value
	var f_82 vm.Value
	var closed_exprs_83 vm.Value
	var nid_84 vm.Value
	var op_85 vm.Value
	var aux_86 vm.Value
	var and__x_87 vm.Value
	var f_91 vm.Value
	var closed_exprs_92 vm.Value
	var nid_93 vm.Value
	var op_94 vm.Value
	var aux_95 vm.Value
	var f_96 vm.Value
	var closed_exprs_97 vm.Value
	var nid_98 vm.Value
	var op_99 vm.Value
	var aux_100 vm.Value
	var v117 bool
	var v400 vm.Value
	var f_401 vm.Value
	var closed_exprs_402 vm.Value
	var nid_403 vm.Value
	var op_404 vm.Value
	var aux_405 vm.Value
	var f_106 vm.Value
	var closed_exprs_107 vm.Value
	var nid_108 vm.Value
	var op_109 vm.Value
	var aux_110 vm.Value
	var f_111 vm.Value
	var closed_exprs_112 vm.Value
	var nid_113 vm.Value
	var op_114 vm.Value
	var aux_115 vm.Value
	var v132 bool
	var v393 vm.Value
	var f_394 vm.Value
	var closed_exprs_395 vm.Value
	var nid_396 vm.Value
	var op_397 vm.Value
	var aux_398 vm.Value
	var f_121 vm.Value
	var closed_exprs_122 vm.Value
	var nid_123 vm.Value
	var op_124 vm.Value
	var aux_125 vm.Value
	var f_126 vm.Value
	var closed_exprs_127 vm.Value
	var nid_128 vm.Value
	var op_129 vm.Value
	var aux_130 vm.Value
	var v147 bool
	var v386 vm.Value
	var f_387 vm.Value
	var closed_exprs_388 vm.Value
	var nid_389 vm.Value
	var op_390 vm.Value
	var aux_391 vm.Value
	var f_136 vm.Value
	var closed_exprs_137 vm.Value
	var nid_138 vm.Value
	var op_139 vm.Value
	var aux_140 vm.Value
	var f_141 vm.Value
	var closed_exprs_142 vm.Value
	var nid_143 vm.Value
	var op_144 vm.Value
	var aux_145 vm.Value
	var v164 vm.Value
	var v379 vm.Value
	var f_380 vm.Value
	var closed_exprs_381 vm.Value
	var nid_382 vm.Value
	var op_383 vm.Value
	var aux_384 vm.Value
	var f_151 vm.Value
	var closed_exprs_152 vm.Value
	var nid_153 vm.Value
	var op_154 vm.Value
	var aux_155 vm.Value
	var v177 vm.Value
	var f_156 vm.Value
	var closed_exprs_157 vm.Value
	var nid_158 vm.Value
	var op_159 vm.Value
	var aux_160 vm.Value
	var or__x_235 bool
	var v372 vm.Value
	var f_373 vm.Value
	var closed_exprs_374 vm.Value
	var nid_375 vm.Value
	var op_376 vm.Value
	var aux_377 vm.Value
	var f_166 vm.Value
	var closed_exprs_167 vm.Value
	var nid_168 vm.Value
	var op_169 vm.Value
	var aux_170 vm.Value
	var rhs_180 vm.Value
	var f_171 vm.Value
	var closed_exprs_172 vm.Value
	var nid_173 vm.Value
	var op_174 vm.Value
	var aux_175 vm.Value
	var v217 vm.Value
	var f_218 vm.Value
	var closed_exprs_219 vm.Value
	var nid_220 vm.Value
	var op_221 vm.Value
	var aux_222 vm.Value
	var f_181 vm.Value
	var closed_exprs_182 vm.Value
	var nid_183 vm.Value
	var op_184 vm.Value
	var aux_185 vm.Value
	var rhs_186 vm.Value
	var arg__15835_197 vm.Value
	var arg__15844_201 vm.Value
	var arg__15846_202 vm.Value
	var v203 vm.Value
	var f_187 vm.Value
	var closed_exprs_188 vm.Value
	var nid_189 vm.Value
	var op_190 vm.Value
	var aux_191 vm.Value
	var rhs_192 vm.Value
	var v207 vm.Value
	var f_208 vm.Value
	var closed_exprs_209 vm.Value
	var nid_210 vm.Value
	var op_211 vm.Value
	var aux_212 vm.Value
	var rhs_213 vm.Value
	var f_224 vm.Value
	var closed_exprs_225 vm.Value
	var nid_226 vm.Value
	var op_227 vm.Value
	var aux_228 vm.Value
	var v272 vm.Value
	var f_229 vm.Value
	var closed_exprs_230 vm.Value
	var nid_231 vm.Value
	var op_232 vm.Value
	var aux_233 vm.Value
	var v330 bool
	var v365 vm.Value
	var f_366 vm.Value
	var closed_exprs_367 vm.Value
	var nid_368 vm.Value
	var op_369 vm.Value
	var aux_370 vm.Value
	var f_236 vm.Value
	var closed_exprs_237 vm.Value
	var nid_238 vm.Value
	var op_239 vm.Value
	var aux_240 vm.Value
	var or__x_241 bool
	var f_242 vm.Value
	var closed_exprs_243 vm.Value
	var nid_244 vm.Value
	var op_245 vm.Value
	var aux_246 vm.Value
	var or__x_247 bool
	var v251 bool
	var v253 bool
	var f_254 vm.Value
	var closed_exprs_255 vm.Value
	var nid_256 vm.Value
	var op_257 vm.Value
	var aux_258 vm.Value
	var or__x_259 vm.Value
	var f_261 vm.Value
	var closed_exprs_262 vm.Value
	var nid_263 vm.Value
	var op_264 vm.Value
	var aux_265 vm.Value
	var rhs_275 vm.Value
	var f_266 vm.Value
	var closed_exprs_267 vm.Value
	var nid_268 vm.Value
	var op_269 vm.Value
	var aux_270 vm.Value
	var v312 vm.Value
	var f_313 vm.Value
	var closed_exprs_314 vm.Value
	var nid_315 vm.Value
	var op_316 vm.Value
	var aux_317 vm.Value
	var f_276 vm.Value
	var closed_exprs_277 vm.Value
	var nid_278 vm.Value
	var op_279 vm.Value
	var aux_280 vm.Value
	var rhs_281 vm.Value
	var arg__15870_292 vm.Value
	var arg__15879_296 vm.Value
	var arg__15881_297 vm.Value
	var v298 vm.Value
	var f_282 vm.Value
	var closed_exprs_283 vm.Value
	var nid_284 vm.Value
	var op_285 vm.Value
	var aux_286 vm.Value
	var rhs_287 vm.Value
	var v302 vm.Value
	var f_303 vm.Value
	var closed_exprs_304 vm.Value
	var nid_305 vm.Value
	var op_306 vm.Value
	var aux_307 vm.Value
	var rhs_308 vm.Value
	var f_319 vm.Value
	var closed_exprs_320 vm.Value
	var nid_321 vm.Value
	var op_322 vm.Value
	var aux_323 vm.Value
	var v333 vm.Value
	var f_324 vm.Value
	var closed_exprs_325 vm.Value
	var nid_326 vm.Value
	var op_327 vm.Value
	var aux_328 vm.Value
	var v358 vm.Value
	var f_359 vm.Value
	var closed_exprs_360 vm.Value
	var nid_361 vm.Value
	var op_362 vm.Value
	var aux_363 vm.Value
	var f_335 vm.Value
	var closed_exprs_336 vm.Value
	var nid_337 vm.Value
	var op_338 vm.Value
	var aux_339 vm.Value
	var f_340 vm.Value
	var closed_exprs_341 vm.Value
	var nid_342 vm.Value
	var op_343 vm.Value
	var aux_344 vm.Value
	var v351 vm.Value
	var f_352 vm.Value
	var closed_exprs_353 vm.Value
	var nid_354 vm.Value
	var op_355 vm.Value
	var aux_356 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = op_4, aux_6, v18, f_7, closed_exprs_8, nid_9, op_10, aux_11, f_12, closed_exprs_13, nid_14, op_15, aux_16, v33, v428, f_429, closed_exprs_430, nid_431, op_432, aux_433, f_22, closed_exprs_23, nid_24, op_25, aux_26, f_27, closed_exprs_28, nid_29, op_30, aux_31, v48, v421, f_422, closed_exprs_423, nid_424, op_425, aux_426, f_37, closed_exprs_38, nid_39, op_40, aux_41, f_42, closed_exprs_43, nid_44, op_45, aux_46, and__x_63, v414, f_415, closed_exprs_416, nid_417, op_418, aux_419, f_52, closed_exprs_53, nid_54, op_55, aux_56, f_57, closed_exprs_58, nid_59, op_60, aux_61, v102, v407, f_408, closed_exprs_409, nid_410, op_411, aux_412, f_64, closed_exprs_65, nid_66, op_67, aux_68, and__x_69, v78, f_70, closed_exprs_71, nid_72, op_73, aux_74, and__x_75, v81, f_82, closed_exprs_83, nid_84, op_85, aux_86, and__x_87, f_91, closed_exprs_92, nid_93, op_94, aux_95, f_96, closed_exprs_97, nid_98, op_99, aux_100, v117, v400, f_401, closed_exprs_402, nid_403, op_404, aux_405, f_106, closed_exprs_107, nid_108, op_109, aux_110, f_111, closed_exprs_112, nid_113, op_114, aux_115, v132, v393, f_394, closed_exprs_395, nid_396, op_397, aux_398, f_121, closed_exprs_122, nid_123, op_124, aux_125, f_126, closed_exprs_127, nid_128, op_129, aux_130, v147, v386, f_387, closed_exprs_388, nid_389, op_390, aux_391, f_136, closed_exprs_137, nid_138, op_139, aux_140, f_141, closed_exprs_142, nid_143, op_144, aux_145, v164, v379, f_380, closed_exprs_381, nid_382, op_383, aux_384, f_151, closed_exprs_152, nid_153, op_154, aux_155, v177, f_156, closed_exprs_157, nid_158, op_159, aux_160, or__x_235, v372, f_373, closed_exprs_374, nid_375, op_376, aux_377, f_166, closed_exprs_167, nid_168, op_169, aux_170, rhs_180, f_171, closed_exprs_172, nid_173, op_174, aux_175, v217, f_218, closed_exprs_219, nid_220, op_221, aux_222, f_181, closed_exprs_182, nid_183, op_184, aux_185, rhs_186, arg__15835_197, arg__15844_201, arg__15846_202, v203, f_187, closed_exprs_188, nid_189, op_190, aux_191, rhs_192, v207, f_208, closed_exprs_209, nid_210, op_211, aux_212, rhs_213, f_224, closed_exprs_225, nid_226, op_227, aux_228, v272, f_229, closed_exprs_230, nid_231, op_232, aux_233, v330, v365, f_366, closed_exprs_367, nid_368, op_369, aux_370, f_236, closed_exprs_237, nid_238, op_239, aux_240, or__x_241, f_242, closed_exprs_243, nid_244, op_245, aux_246, or__x_247, v251, v253, f_254, closed_exprs_255, nid_256, op_257, aux_258, or__x_259, f_261, closed_exprs_262, nid_263, op_264, aux_265, rhs_275, f_266, closed_exprs_267, nid_268, op_269, aux_270, v312, f_313, closed_exprs_314, nid_315, op_316, aux_317, f_276, closed_exprs_277, nid_278, op_279, aux_280, rhs_281, arg__15870_292, arg__15879_296, arg__15881_297, v298, f_282, closed_exprs_283, nid_284, op_285, aux_286, rhs_287, v302, f_303, closed_exprs_304, nid_305, op_306, aux_307, rhs_308, f_319, closed_exprs_320, nid_321, op_322, aux_323, v333, f_324, closed_exprs_325, nid_326, op_327, aux_328, v358, f_359, closed_exprs_360, nid_361, op_362, aux_363, f_335, closed_exprs_336, nid_337, op_338, aux_339, f_340, closed_exprs_341, nid_342, op_343, aux_344, v351, f_352, closed_exprs_353, nid_354, op_355, aux_356
	op_4, callErr = rt.InvokeValue(rt.LookupVar("ir", "op").Deref(), []vm.Value{arg2, arg0})
	if callErr != nil {
		return nil, callErr
	}
	aux_6, callErr = rt.InvokeValue(rt.LookupVar("ir", "aux").Deref(), []vm.Value{arg2, arg0})
	if callErr != nil {
		return nil, callErr
	}
	v18 = op_4 == vm.Keyword("load-arg")
	if v18 {
		f_7 = arg0
		closed_exprs_8 = arg1
		nid_9 = arg2
		op_10 = op_4
		aux_11 = aux_6
		goto b1
	} else {
		f_12 = arg0
		closed_exprs_13 = arg1
		nid_14 = arg2
		op_15 = op_4
		aux_16 = aux_6
		goto b2
	}
b1:
	;
	v428 = vm.NewArrayVector([]vm.Value{})
	f_429 = f_7
	closed_exprs_430 = closed_exprs_8
	nid_431 = nid_9
	op_432 = op_10
	aux_433 = aux_11
	goto b3
b2:
	;
	v33 = op_15 == vm.Keyword("load-var")
	if v33 {
		f_22 = f_12
		closed_exprs_23 = closed_exprs_13
		nid_24 = nid_14
		op_25 = op_15
		aux_26 = aux_16
		goto b4
	} else {
		f_27 = f_12
		closed_exprs_28 = closed_exprs_13
		nid_29 = nid_14
		op_30 = op_15
		aux_31 = aux_16
		goto b5
	}
b3:
	;
	return v428, nil
b4:
	;
	v421 = vm.NewArrayVector([]vm.Value{})
	f_422 = f_22
	closed_exprs_423 = closed_exprs_23
	nid_424 = nid_24
	op_425 = op_25
	aux_426 = aux_26
	goto b6
b5:
	;
	v48 = op_30 == vm.Keyword("load-closed")
	if v48 {
		f_37 = f_27
		closed_exprs_38 = closed_exprs_28
		nid_39 = nid_29
		op_40 = op_30
		aux_41 = aux_31
		goto b7
	} else {
		f_42 = f_27
		closed_exprs_43 = closed_exprs_28
		nid_44 = nid_29
		op_45 = op_30
		aux_46 = aux_31
		goto b8
	}
b6:
	;
	v428 = v421
	f_429 = f_422
	closed_exprs_430 = closed_exprs_423
	nid_431 = nid_424
	op_432 = op_425
	aux_433 = aux_426
	goto b3
b7:
	;
	v414 = vm.NewArrayVector([]vm.Value{})
	f_415 = f_37
	closed_exprs_416 = closed_exprs_38
	nid_417 = nid_39
	op_418 = op_40
	aux_419 = aux_41
	goto b9
b8:
	;
	and__x_63 = op_45 == vm.Keyword("const")
	if and__x_63 {
		f_64 = f_42
		closed_exprs_65 = closed_exprs_43
		nid_66 = nid_44
		op_67 = op_45
		aux_68 = aux_46
		and__x_69 = and__x_63
		goto b13
	} else {
		f_70 = f_42
		closed_exprs_71 = closed_exprs_43
		nid_72 = nid_44
		op_73 = op_45
		aux_74 = aux_46
		and__x_75 = and__x_63
		goto b14
	}
b9:
	;
	v421 = v414
	f_422 = f_415
	closed_exprs_423 = closed_exprs_416
	nid_424 = nid_417
	op_425 = op_418
	aux_426 = aux_419
	goto b6
b10:
	;
	v407 = vm.NewArrayVector([]vm.Value{})
	f_408 = f_52
	closed_exprs_409 = closed_exprs_53
	nid_410 = nid_54
	op_411 = op_55
	aux_412 = aux_56
	goto b12
b11:
	;
	v102 = op_60 == vm.Keyword("const")
	if v102 {
		f_91 = f_57
		closed_exprs_92 = closed_exprs_58
		nid_93 = nid_59
		op_94 = op_60
		aux_95 = aux_61
		goto b16
	} else {
		f_96 = f_57
		closed_exprs_97 = closed_exprs_58
		nid_98 = nid_59
		op_99 = op_60
		aux_100 = aux_61
		goto b17
	}
b12:
	;
	v414 = v407
	f_415 = f_408
	closed_exprs_416 = closed_exprs_409
	nid_417 = nid_410
	op_418 = op_411
	aux_419 = aux_412
	goto b9
b13:
	;
	v78, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "any-fn-template?").Deref(), []vm.Value{aux_68})
	if callErr != nil {
		return nil, callErr
	}
	v81 = v78
	f_82 = f_64
	closed_exprs_83 = closed_exprs_65
	nid_84 = nid_66
	op_85 = op_67
	aux_86 = aux_68
	and__x_87 = vm.Boolean(and__x_69)
	goto b15
b14:
	;
	v81 = vm.Boolean(and__x_75)
	f_82 = f_70
	closed_exprs_83 = closed_exprs_71
	nid_84 = nid_72
	op_85 = op_73
	aux_86 = aux_74
	and__x_87 = vm.Boolean(and__x_75)
	goto b15
b15:
	;
	if vm.IsTruthy(v81) {
		f_52 = f_82
		closed_exprs_53 = closed_exprs_83
		nid_54 = nid_84
		op_55 = op_85
		aux_56 = aux_86
		goto b10
	} else {
		f_57 = f_82
		closed_exprs_58 = closed_exprs_83
		nid_59 = nid_84
		op_60 = op_85
		aux_61 = aux_86
		goto b11
	}
b16:
	;
	v400 = vm.NewArrayVector([]vm.Value{})
	f_401 = f_91
	closed_exprs_402 = closed_exprs_92
	nid_403 = nid_93
	op_404 = op_94
	aux_405 = aux_95
	goto b18
b17:
	;
	v117 = op_99 == vm.Keyword("make-closure")
	if v117 {
		f_106 = f_96
		closed_exprs_107 = closed_exprs_97
		nid_108 = nid_98
		op_109 = op_99
		aux_110 = aux_100
		goto b19
	} else {
		f_111 = f_96
		closed_exprs_112 = closed_exprs_97
		nid_113 = nid_98
		op_114 = op_99
		aux_115 = aux_100
		goto b20
	}
b18:
	;
	v407 = v400
	f_408 = f_401
	closed_exprs_409 = closed_exprs_402
	nid_410 = nid_403
	op_411 = op_404
	aux_412 = aux_405
	goto b12
b19:
	;
	v393 = vm.NewArrayVector([]vm.Value{})
	f_394 = f_106
	closed_exprs_395 = closed_exprs_107
	nid_396 = nid_108
	op_397 = op_109
	aux_398 = aux_110
	goto b21
b20:
	;
	v132 = op_114 == vm.Keyword("push-closed")
	if v132 {
		f_121 = f_111
		closed_exprs_122 = closed_exprs_112
		nid_123 = nid_113
		op_124 = op_114
		aux_125 = aux_115
		goto b22
	} else {
		f_126 = f_111
		closed_exprs_127 = closed_exprs_112
		nid_128 = nid_113
		op_129 = op_114
		aux_130 = aux_115
		goto b23
	}
b21:
	;
	v400 = v393
	f_401 = f_394
	closed_exprs_402 = closed_exprs_395
	nid_403 = nid_396
	op_404 = op_397
	aux_405 = aux_398
	goto b18
b22:
	;
	v386 = vm.NewArrayVector([]vm.Value{})
	f_387 = f_121
	closed_exprs_388 = closed_exprs_122
	nid_389 = nid_123
	op_390 = op_124
	aux_391 = aux_125
	goto b24
b23:
	;
	v147 = op_129 == vm.Keyword("block-arg")
	if v147 {
		f_136 = f_126
		closed_exprs_137 = closed_exprs_127
		nid_138 = nid_128
		op_139 = op_129
		aux_140 = aux_130
		goto b25
	} else {
		f_141 = f_126
		closed_exprs_142 = closed_exprs_127
		nid_143 = nid_128
		op_144 = op_129
		aux_145 = aux_130
		goto b26
	}
b24:
	;
	v393 = v386
	f_394 = f_387
	closed_exprs_395 = closed_exprs_388
	nid_396 = nid_389
	op_397 = op_390
	aux_398 = aux_391
	goto b21
b25:
	;
	v379 = vm.NewArrayVector([]vm.Value{})
	f_380 = f_136
	closed_exprs_381 = closed_exprs_137
	nid_382 = nid_138
	op_383 = op_139
	aux_384 = aux_140
	goto b27
b26:
	;
	v164, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "contains?").Deref(), []vm.Value{rt.LookupVar("ir.lower-go", "binary-op").Deref(), op_144})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v164) {
		f_151 = f_141
		closed_exprs_152 = closed_exprs_142
		nid_153 = nid_143
		op_154 = op_144
		aux_155 = aux_145
		goto b28
	} else {
		f_156 = f_141
		closed_exprs_157 = closed_exprs_142
		nid_158 = nid_143
		op_159 = op_144
		aux_160 = aux_145
		goto b29
	}
b27:
	;
	v386 = v379
	f_387 = f_380
	closed_exprs_388 = closed_exprs_381
	nid_389 = nid_382
	op_390 = op_383
	aux_391 = aux_384
	goto b24
b28:
	;
	v177, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "used?").Deref(), []vm.Value{f_151, nid_153})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v177) {
		f_166 = f_151
		closed_exprs_167 = closed_exprs_152
		nid_168 = nid_153
		op_169 = op_154
		aux_170 = aux_155
		goto b31
	} else {
		f_171 = f_151
		closed_exprs_172 = closed_exprs_152
		nid_173 = nid_153
		op_174 = op_154
		aux_175 = aux_155
		goto b32
	}
b29:
	;
	or__x_235 = op_159 == vm.Keyword("inc")
	if or__x_235 {
		f_236 = f_156
		closed_exprs_237 = closed_exprs_157
		nid_238 = nid_158
		op_239 = op_159
		aux_240 = aux_160
		or__x_241 = or__x_235
		goto b40
	} else {
		f_242 = f_156
		closed_exprs_243 = closed_exprs_157
		nid_244 = nid_158
		op_245 = op_159
		aux_246 = aux_160
		or__x_247 = or__x_235
		goto b41
	}
b30:
	;
	v379 = v372
	f_380 = f_373
	closed_exprs_381 = closed_exprs_374
	nid_382 = nid_375
	op_383 = op_376
	aux_384 = aux_377
	goto b27
b31:
	;
	rhs_180, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "inst-rhs").Deref(), []vm.Value{f_166, closed_exprs_167, nid_168})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(rhs_180) {
		f_181 = f_166
		closed_exprs_182 = closed_exprs_167
		nid_183 = nid_168
		op_184 = op_169
		aux_185 = aux_170
		rhs_186 = rhs_180
		goto b34
	} else {
		f_187 = f_166
		closed_exprs_188 = closed_exprs_167
		nid_189 = nid_168
		op_190 = op_169
		aux_191 = aux_170
		rhs_192 = rhs_180
		goto b35
	}
b32:
	;
	v217 = vm.NewArrayVector([]vm.Value{})
	f_218 = f_171
	closed_exprs_219 = closed_exprs_172
	nid_220 = nid_173
	op_221 = op_174
	aux_222 = aux_175
	goto b33
b33:
	;
	v372 = v217
	f_373 = f_218
	closed_exprs_374 = closed_exprs_219
	nid_375 = nid_220
	op_376 = op_221
	aux_377 = aux_222
	goto b30
b34:
	;
	arg__15835_197, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "value-ident").Deref(), []vm.Value{f_181, nid_183})
	if callErr != nil {
		return nil, callErr
	}
	arg__15844_201, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "value-ident").Deref(), []vm.Value{f_181, nid_183})
	if callErr != nil {
		return nil, callErr
	}
	arg__15846_202, callErr = rt.InvokeValue(rt.LookupVar("gogen", "assign").Deref(), []vm.Value{vm.String("="), arg__15844_201, rhs_186})
	if callErr != nil {
		return nil, callErr
	}
	v203, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__15846_202})
	if callErr != nil {
		return nil, callErr
	}
	v207 = v203
	f_208 = f_181
	closed_exprs_209 = closed_exprs_182
	nid_210 = nid_183
	op_211 = op_184
	aux_212 = aux_185
	rhs_213 = rhs_186
	goto b36
b35:
	;
	v207 = vm.NIL
	f_208 = f_187
	closed_exprs_209 = closed_exprs_188
	nid_210 = nid_189
	op_211 = op_190
	aux_212 = aux_191
	rhs_213 = rhs_192
	goto b36
b36:
	;
	v217 = v207
	f_218 = f_208
	closed_exprs_219 = closed_exprs_209
	nid_220 = nid_210
	op_221 = op_211
	aux_222 = aux_212
	goto b33
b37:
	;
	v272, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "used?").Deref(), []vm.Value{f_224, nid_226})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v272) {
		f_261 = f_224
		closed_exprs_262 = closed_exprs_225
		nid_263 = nid_226
		op_264 = op_227
		aux_265 = aux_228
		goto b43
	} else {
		f_266 = f_224
		closed_exprs_267 = closed_exprs_225
		nid_268 = nid_226
		op_269 = op_227
		aux_270 = aux_228
		goto b44
	}
b38:
	;
	v330 = op_232 == vm.Keyword("call")
	if v330 {
		f_319 = f_229
		closed_exprs_320 = closed_exprs_230
		nid_321 = nid_231
		op_322 = op_232
		aux_323 = aux_233
		goto b49
	} else {
		f_324 = f_229
		closed_exprs_325 = closed_exprs_230
		nid_326 = nid_231
		op_327 = op_232
		aux_328 = aux_233
		goto b50
	}
b39:
	;
	v372 = v365
	f_373 = f_366
	closed_exprs_374 = closed_exprs_367
	nid_375 = nid_368
	op_376 = op_369
	aux_377 = aux_370
	goto b30
b40:
	;
	v253 = or__x_241
	f_254 = f_236
	closed_exprs_255 = closed_exprs_237
	nid_256 = nid_238
	op_257 = op_239
	aux_258 = aux_240
	or__x_259 = vm.Boolean(or__x_241)
	goto b42
b41:
	;
	v251 = op_245 == vm.Keyword("dec")
	v253 = v251
	f_254 = f_242
	closed_exprs_255 = closed_exprs_243
	nid_256 = nid_244
	op_257 = op_245
	aux_258 = aux_246
	or__x_259 = vm.Boolean(or__x_247)
	goto b42
b42:
	;
	if v253 {
		f_224 = f_254
		closed_exprs_225 = closed_exprs_255
		nid_226 = nid_256
		op_227 = op_257
		aux_228 = aux_258
		goto b37
	} else {
		f_229 = f_254
		closed_exprs_230 = closed_exprs_255
		nid_231 = nid_256
		op_232 = op_257
		aux_233 = aux_258
		goto b38
	}
b43:
	;
	rhs_275, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "inst-rhs").Deref(), []vm.Value{f_261, closed_exprs_262, nid_263})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(rhs_275) {
		f_276 = f_261
		closed_exprs_277 = closed_exprs_262
		nid_278 = nid_263
		op_279 = op_264
		aux_280 = aux_265
		rhs_281 = rhs_275
		goto b46
	} else {
		f_282 = f_261
		closed_exprs_283 = closed_exprs_262
		nid_284 = nid_263
		op_285 = op_264
		aux_286 = aux_265
		rhs_287 = rhs_275
		goto b47
	}
b44:
	;
	v312 = vm.NewArrayVector([]vm.Value{})
	f_313 = f_266
	closed_exprs_314 = closed_exprs_267
	nid_315 = nid_268
	op_316 = op_269
	aux_317 = aux_270
	goto b45
b45:
	;
	v365 = v312
	f_366 = f_313
	closed_exprs_367 = closed_exprs_314
	nid_368 = nid_315
	op_369 = op_316
	aux_370 = aux_317
	goto b39
b46:
	;
	arg__15870_292, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "value-ident").Deref(), []vm.Value{f_276, nid_278})
	if callErr != nil {
		return nil, callErr
	}
	arg__15879_296, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "value-ident").Deref(), []vm.Value{f_276, nid_278})
	if callErr != nil {
		return nil, callErr
	}
	arg__15881_297, callErr = rt.InvokeValue(rt.LookupVar("gogen", "assign").Deref(), []vm.Value{vm.String("="), arg__15879_296, rhs_281})
	if callErr != nil {
		return nil, callErr
	}
	v298, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__15881_297})
	if callErr != nil {
		return nil, callErr
	}
	v302 = v298
	f_303 = f_276
	closed_exprs_304 = closed_exprs_277
	nid_305 = nid_278
	op_306 = op_279
	aux_307 = aux_280
	rhs_308 = rhs_281
	goto b48
b47:
	;
	v302 = vm.NIL
	f_303 = f_282
	closed_exprs_304 = closed_exprs_283
	nid_305 = nid_284
	op_306 = op_285
	aux_307 = aux_286
	rhs_308 = rhs_287
	goto b48
b48:
	;
	v312 = v302
	f_313 = f_303
	closed_exprs_314 = closed_exprs_304
	nid_315 = nid_305
	op_316 = op_306
	aux_317 = aux_307
	goto b45
b49:
	;
	v333, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "call-assign-stmts").Deref(), []vm.Value{f_319, closed_exprs_320, nid_321})
	if callErr != nil {
		return nil, callErr
	}
	v358 = v333
	f_359 = f_319
	closed_exprs_360 = closed_exprs_320
	nid_361 = nid_321
	op_362 = op_322
	aux_363 = aux_323
	goto b51
b50:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		f_335 = f_324
		closed_exprs_336 = closed_exprs_325
		nid_337 = nid_326
		op_338 = op_327
		aux_339 = aux_328
		goto b52
	} else {
		f_340 = f_324
		closed_exprs_341 = closed_exprs_325
		nid_342 = nid_326
		op_343 = op_327
		aux_344 = aux_328
		goto b53
	}
b51:
	;
	v365 = v358
	f_366 = f_359
	closed_exprs_367 = closed_exprs_360
	nid_368 = nid_361
	op_369 = op_362
	aux_370 = aux_363
	goto b39
b52:
	;
	v351 = vm.NIL
	f_352 = f_335
	closed_exprs_353 = closed_exprs_336
	nid_354 = nid_337
	op_355 = op_338
	aux_356 = aux_339
	goto b54
b53:
	;
	v351 = vm.NIL
	f_352 = f_340
	closed_exprs_353 = closed_exprs_341
	nid_354 = nid_342
	op_355 = op_343
	aux_356 = aux_344
	goto b54
b54:
	;
	v358 = v351
	f_359 = f_352
	closed_exprs_360 = closed_exprs_353
	nid_361 = nid_354
	op_362 = op_355
	aux_363 = aux_356
	goto b51
}
func lower_template_closure_expr(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var v7 vm.Value
	var template_2 vm.Value
	var capture_exprs_3 vm.Value
	var arg__15896_10 vm.Value
	var arg__15901_13 vm.Value
	var inner_STAR__14 vm.Value
	var template_4 vm.Value
	var capture_exprs_5 vm.Value
	var v65 vm.Value
	var v170 vm.Value
	var template_171 vm.Value
	var capture_exprs_172 vm.Value
	var template_15 vm.Value
	var capture_exprs_16 vm.Value
	var inner_STAR__17 vm.Value
	var arg__15906_25 vm.Value
	var arg__15912_31 vm.Value
	var arg__15914_33 vm.Value
	var arg__15917_35 vm.Value
	var arg__15922_40 vm.Value
	var arg__15928_46 vm.Value
	var arg__15930_48 vm.Value
	var arg__15933_50 vm.Value
	var v51 vm.Value
	var template_18 vm.Value
	var capture_exprs_19 vm.Value
	var inner_STAR__20 vm.Value
	var v55 vm.Value
	var template_56 vm.Value
	var capture_exprs_57 vm.Value
	var inner_STAR__58 vm.Value
	var template_60 vm.Value
	var capture_exprs_61 vm.Value
	var arg__15949_71 vm.Value
	var arg__15963_77 vm.Value
	var inners_78 vm.Value
	var v88 vm.Value
	var template_62 vm.Value
	var capture_exprs_63 vm.Value
	var v166 vm.Value
	var template_167 vm.Value
	var capture_exprs_168 vm.Value
	var template_79 vm.Value
	var capture_exprs_80 vm.Value
	var inners_81 vm.Value
	var boxed_93 vm.Value
	var arg__16039_97 vm.Value
	var arg__16045_103 vm.Value
	var arg__16047_105 vm.Value
	var arg__16052_110 vm.Value
	var arg__16058_115 vm.Value
	var arg__16060_116 vm.Value
	var arg__16061_117 vm.Value
	var arg__16066_122 vm.Value
	var arg__16072_128 vm.Value
	var arg__16074_130 vm.Value
	var arg__16079_135 vm.Value
	var arg__16085_140 vm.Value
	var arg__16087_141 vm.Value
	var arg__16088_142 vm.Value
	var v143 vm.Value
	var template_82 vm.Value
	var capture_exprs_83 vm.Value
	var inners_84 vm.Value
	var v147 vm.Value
	var template_148 vm.Value
	var capture_exprs_149 vm.Value
	var inners_150 vm.Value
	var template_152 vm.Value
	var capture_exprs_153 vm.Value
	var template_154 vm.Value
	var capture_exprs_155 vm.Value
	var v162 vm.Value
	var template_163 vm.Value
	var capture_exprs_164 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = v7, template_2, capture_exprs_3, arg__15896_10, arg__15901_13, inner_STAR__14, template_4, capture_exprs_5, v65, v170, template_171, capture_exprs_172, template_15, capture_exprs_16, inner_STAR__17, arg__15906_25, arg__15912_31, arg__15914_33, arg__15917_35, arg__15922_40, arg__15928_46, arg__15930_48, arg__15933_50, v51, template_18, capture_exprs_19, inner_STAR__20, v55, template_56, capture_exprs_57, inner_STAR__58, template_60, capture_exprs_61, arg__15949_71, arg__15963_77, inners_78, v88, template_62, capture_exprs_63, v166, template_167, capture_exprs_168, template_79, capture_exprs_80, inners_81, boxed_93, arg__16039_97, arg__16045_103, arg__16047_105, arg__16052_110, arg__16058_115, arg__16060_116, arg__16061_117, arg__16066_122, arg__16072_128, arg__16074_130, arg__16079_135, arg__16085_140, arg__16087_141, arg__16088_142, v143, template_82, capture_exprs_83, inners_84, v147, template_148, capture_exprs_149, inners_150, template_152, capture_exprs_153, template_154, capture_exprs_155, v162, template_163, capture_exprs_164
	v7, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "fn-template?").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v7) {
		template_2 = arg0
		capture_exprs_3 = arg1
		goto b1
	} else {
		template_4 = arg0
		capture_exprs_5 = arg1
		goto b2
	}
b1:
	;
	arg__15896_10, callErr = rt.InvokeValue(vm.Keyword("fn"), []vm.Value{template_2})
	if callErr != nil {
		return nil, callErr
	}
	arg__15901_13, callErr = rt.InvokeValue(vm.Keyword("fn"), []vm.Value{template_2})
	if callErr != nil {
		return nil, callErr
	}
	inner_STAR__14, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "lower-fn-lit").Deref(), []vm.Value{arg__15901_13, capture_exprs_3})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(inner_STAR__14) {
		template_15 = template_2
		capture_exprs_16 = capture_exprs_3
		inner_STAR__17 = inner_STAR__14
		goto b4
	} else {
		template_18 = template_2
		capture_exprs_19 = capture_exprs_3
		inner_STAR__20 = inner_STAR__14
		goto b5
	}
b2:
	;
	v65, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "multi-fn-template?").Deref(), []vm.Value{template_4})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v65) {
		template_60 = template_4
		capture_exprs_61 = capture_exprs_5
		goto b7
	} else {
		template_62 = template_4
		capture_exprs_63 = capture_exprs_5
		goto b8
	}
b3:
	;
	return v170, nil
b4:
	;
	arg__15906_25, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("rt")})
	if callErr != nil {
		return nil, callErr
	}
	arg__15912_31, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("rt")})
	if callErr != nil {
		return nil, callErr
	}
	arg__15914_33, callErr = rt.InvokeValue(rt.LookupVar("gogen", "field-sel").Deref(), []vm.Value{arg__15912_31, vm.String("BoxNativeFn")})
	if callErr != nil {
		return nil, callErr
	}
	arg__15917_35, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{inner_STAR__17})
	if callErr != nil {
		return nil, callErr
	}
	arg__15922_40, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("rt")})
	if callErr != nil {
		return nil, callErr
	}
	arg__15928_46, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("rt")})
	if callErr != nil {
		return nil, callErr
	}
	arg__15930_48, callErr = rt.InvokeValue(rt.LookupVar("gogen", "field-sel").Deref(), []vm.Value{arg__15928_46, vm.String("BoxNativeFn")})
	if callErr != nil {
		return nil, callErr
	}
	arg__15933_50, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{inner_STAR__17})
	if callErr != nil {
		return nil, callErr
	}
	v51, callErr = rt.InvokeValue(rt.LookupVar("gogen", "call").Deref(), []vm.Value{arg__15930_48, arg__15933_50})
	if callErr != nil {
		return nil, callErr
	}
	v55 = v51
	template_56 = template_15
	capture_exprs_57 = capture_exprs_16
	inner_STAR__58 = inner_STAR__17
	goto b6
b5:
	;
	v55 = vm.NIL
	template_56 = template_18
	capture_exprs_57 = capture_exprs_19
	inner_STAR__58 = inner_STAR__20
	goto b6
b6:
	;
	v170 = v55
	template_171 = template_56
	capture_exprs_172 = capture_exprs_57
	goto b3
b7:
	;
	arg__15949_71, callErr = rt.InvokeValue(vm.Keyword("fns"), []vm.Value{template_60})
	if callErr != nil {
		return nil, callErr
	}
	arg__15963_77, callErr = rt.InvokeValue(vm.Keyword("fns"), []vm.Value{template_60})
	if callErr != nil {
		return nil, callErr
	}
	inners_78, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapv").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var arg__15953_3 vm.Value
		var arg__15958_6 vm.Value
		var v7 vm.Value
		var callErr error
		_, _, _ = arg__15953_3, arg__15958_6, v7
		arg__15953_3, callErr = rt.InvokeValue(vm.Keyword("fn"), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		arg__15958_6, callErr = rt.InvokeValue(vm.Keyword("fn"), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		v7, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "lower-fn-lit").Deref(), []vm.Value{arg__15958_6, capture_exprs_61})
		if callErr != nil {
			return nil, callErr
		}
		return v7, nil
	}), arg__15963_77})
	if callErr != nil {
		return nil, callErr
	}
	v88, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "every?").Deref(), []vm.Value{rt.LookupVar("clojure.core", "identity").Deref(), inners_78})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v88) {
		template_79 = template_60
		capture_exprs_80 = capture_exprs_61
		inners_81 = inners_78
		goto b10
	} else {
		template_82 = template_60
		capture_exprs_83 = capture_exprs_61
		inners_84 = inners_78
		goto b11
	}
b8:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		template_152 = template_62
		capture_exprs_153 = capture_exprs_63
		goto b13
	} else {
		template_154 = template_62
		capture_exprs_155 = capture_exprs_63
		goto b14
	}
b9:
	;
	v170 = v166
	template_171 = template_167
	capture_exprs_172 = capture_exprs_168
	goto b3
b10:
	;
	boxed_93, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapv").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var arg__16006_4 vm.Value
		var arg__16012_10 vm.Value
		var arg__16014_12 vm.Value
		var arg__16017_14 vm.Value
		var arg__16022_19 vm.Value
		var arg__16028_25 vm.Value
		var arg__16030_27 vm.Value
		var arg__16033_29 vm.Value
		var v30 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _, _ = arg__16006_4, arg__16012_10, arg__16014_12, arg__16017_14, arg__16022_19, arg__16028_25, arg__16030_27, arg__16033_29, v30
		arg__16006_4, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("rt")})
		if callErr != nil {
			return nil, callErr
		}
		arg__16012_10, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("rt")})
		if callErr != nil {
			return nil, callErr
		}
		arg__16014_12, callErr = rt.InvokeValue(rt.LookupVar("gogen", "field-sel").Deref(), []vm.Value{arg__16012_10, vm.String("BoxNativeFn")})
		if callErr != nil {
			return nil, callErr
		}
		arg__16017_14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		arg__16022_19, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("rt")})
		if callErr != nil {
			return nil, callErr
		}
		arg__16028_25, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("rt")})
		if callErr != nil {
			return nil, callErr
		}
		arg__16030_27, callErr = rt.InvokeValue(rt.LookupVar("gogen", "field-sel").Deref(), []vm.Value{arg__16028_25, vm.String("BoxNativeFn")})
		if callErr != nil {
			return nil, callErr
		}
		arg__16033_29, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		v30, callErr = rt.InvokeValue(rt.LookupVar("gogen", "call").Deref(), []vm.Value{arg__16030_27, arg__16033_29})
		if callErr != nil {
			return nil, callErr
		}
		return v30, nil
	}), inners_81})
	if callErr != nil {
		return nil, callErr
	}
	arg__16039_97, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("rt")})
	if callErr != nil {
		return nil, callErr
	}
	arg__16045_103, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("rt")})
	if callErr != nil {
		return nil, callErr
	}
	arg__16047_105, callErr = rt.InvokeValue(rt.LookupVar("gogen", "field-sel").Deref(), []vm.Value{arg__16045_103, vm.String("MakeNativeMultiArity")})
	if callErr != nil {
		return nil, callErr
	}
	arg__16052_110, callErr = rt.InvokeValue(rt.LookupVar("gogen", "type").Deref(), []vm.Value{vm.String("[]vm.Value")})
	if callErr != nil {
		return nil, callErr
	}
	arg__16058_115, callErr = rt.InvokeValue(rt.LookupVar("gogen", "type").Deref(), []vm.Value{vm.String("[]vm.Value")})
	if callErr != nil {
		return nil, callErr
	}
	arg__16060_116, callErr = rt.InvokeValue(rt.LookupVar("gogen", "composite-lit").Deref(), []vm.Value{arg__16058_115, boxed_93})
	if callErr != nil {
		return nil, callErr
	}
	arg__16061_117, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__16060_116})
	if callErr != nil {
		return nil, callErr
	}
	arg__16066_122, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("rt")})
	if callErr != nil {
		return nil, callErr
	}
	arg__16072_128, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("rt")})
	if callErr != nil {
		return nil, callErr
	}
	arg__16074_130, callErr = rt.InvokeValue(rt.LookupVar("gogen", "field-sel").Deref(), []vm.Value{arg__16072_128, vm.String("MakeNativeMultiArity")})
	if callErr != nil {
		return nil, callErr
	}
	arg__16079_135, callErr = rt.InvokeValue(rt.LookupVar("gogen", "type").Deref(), []vm.Value{vm.String("[]vm.Value")})
	if callErr != nil {
		return nil, callErr
	}
	arg__16085_140, callErr = rt.InvokeValue(rt.LookupVar("gogen", "type").Deref(), []vm.Value{vm.String("[]vm.Value")})
	if callErr != nil {
		return nil, callErr
	}
	arg__16087_141, callErr = rt.InvokeValue(rt.LookupVar("gogen", "composite-lit").Deref(), []vm.Value{arg__16085_140, boxed_93})
	if callErr != nil {
		return nil, callErr
	}
	arg__16088_142, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__16087_141})
	if callErr != nil {
		return nil, callErr
	}
	v143, callErr = rt.InvokeValue(rt.LookupVar("gogen", "call").Deref(), []vm.Value{arg__16074_130, arg__16088_142})
	if callErr != nil {
		return nil, callErr
	}
	v147 = v143
	template_148 = template_79
	capture_exprs_149 = capture_exprs_80
	inners_150 = inners_81
	goto b12
b11:
	;
	v147 = vm.NIL
	template_148 = template_82
	capture_exprs_149 = capture_exprs_83
	inners_150 = inners_84
	goto b12
b12:
	;
	v166 = v147
	template_167 = template_148
	capture_exprs_168 = capture_exprs_149
	goto b9
b13:
	;
	v162 = vm.NIL
	template_163 = template_152
	capture_exprs_164 = capture_exprs_153
	goto b15
b14:
	;
	v162 = vm.NIL
	template_163 = template_154
	capture_exprs_164 = capture_exprs_155
	goto b15
b15:
	;
	v166 = v162
	template_167 = template_163
	capture_exprs_168 = capture_exprs_164
	goto b9
}
func lower_terminator(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var term_4 vm.Value
	var op_6 vm.Value
	var refs_8 vm.Value
	var aux_10 vm.Value
	var needs_error_QMARK__12 vm.Value
	var v30 bool
	var f_13 vm.Value
	var closed_exprs_14 vm.Value
	var bid_15 vm.Value
	var term_16 vm.Value
	var op_17 vm.Value
	var refs_18 vm.Value
	var aux_19 vm.Value
	var needs_error_QMARK__20 vm.Value
	var arg__16118_50 vm.Value
	var v51 bool
	var f_21 vm.Value
	var closed_exprs_22 vm.Value
	var bid_23 vm.Value
	var term_24 vm.Value
	var op_25 vm.Value
	var refs_26 vm.Value
	var aux_27 vm.Value
	var needs_error_QMARK__28 vm.Value
	var v270 bool
	var v1024 vm.Value
	var f_1025 vm.Value
	var closed_exprs_1026 vm.Value
	var bid_1027 vm.Value
	var term_1028 vm.Value
	var op_1029 vm.Value
	var refs_1030 vm.Value
	var aux_1031 vm.Value
	var needs_error_QMARK__1032 vm.Value
	var f_32 vm.Value
	var closed_exprs_33 vm.Value
	var bid_34 vm.Value
	var term_35 vm.Value
	var op_36 vm.Value
	var refs_37 vm.Value
	var aux_38 vm.Value
	var needs_error_QMARK__39 vm.Value
	var ret_spec_54 vm.Value
	var v74 bool
	var f_40 vm.Value
	var closed_exprs_41 vm.Value
	var bid_42 vm.Value
	var term_43 vm.Value
	var op_44 vm.Value
	var refs_45 vm.Value
	var aux_46 vm.Value
	var needs_error_QMARK__47 vm.Value
	var v243 vm.Value
	var f_244 vm.Value
	var closed_exprs_245 vm.Value
	var bid_246 vm.Value
	var term_247 vm.Value
	var op_248 vm.Value
	var refs_249 vm.Value
	var aux_250 vm.Value
	var needs_error_QMARK__251 vm.Value
	var f_55 vm.Value
	var closed_exprs_56 vm.Value
	var bid_57 vm.Value
	var term_58 vm.Value
	var op_59 vm.Value
	var refs_60 vm.Value
	var aux_61 vm.Value
	var needs_error_QMARK__62 vm.Value
	var ret_spec_63 vm.Value
	var arg__16131_79 vm.Value
	var arg__16140_84 vm.Value
	var v85 vm.Value
	var f_64 vm.Value
	var closed_exprs_65 vm.Value
	var bid_66 vm.Value
	var term_67 vm.Value
	var op_68 vm.Value
	var refs_69 vm.Value
	var aux_70 vm.Value
	var needs_error_QMARK__71 vm.Value
	var ret_spec_72 vm.Value
	var arg__16148_90 vm.Value
	var arg__16157_95 vm.Value
	var v96 vm.Value
	var expr_98 vm.Value
	var f_99 vm.Value
	var closed_exprs_100 vm.Value
	var bid_101 vm.Value
	var term_102 vm.Value
	var op_103 vm.Value
	var refs_104 vm.Value
	var aux_105 vm.Value
	var needs_error_QMARK__106 vm.Value
	var ret_spec_107 vm.Value
	var expr_108 vm.Value
	var f_109 vm.Value
	var closed_exprs_110 vm.Value
	var bid_111 vm.Value
	var term_112 vm.Value
	var op_113 vm.Value
	var refs_114 vm.Value
	var aux_115 vm.Value
	var needs_error_QMARK__116 vm.Value
	var ret_spec_117 vm.Value
	var expr_118 vm.Value
	var f_119 vm.Value
	var closed_exprs_120 vm.Value
	var bid_121 vm.Value
	var term_122 vm.Value
	var op_123 vm.Value
	var refs_124 vm.Value
	var aux_125 vm.Value
	var needs_error_QMARK__126 vm.Value
	var ret_spec_127 vm.Value
	var v229 vm.Value
	var expr_230 vm.Value
	var f_231 vm.Value
	var closed_exprs_232 vm.Value
	var bid_233 vm.Value
	var term_234 vm.Value
	var op_235 vm.Value
	var refs_236 vm.Value
	var aux_237 vm.Value
	var needs_error_QMARK__238 vm.Value
	var ret_spec_239 vm.Value
	var expr_130 vm.Value
	var f_131 vm.Value
	var closed_exprs_132 vm.Value
	var bid_133 vm.Value
	var term_134 vm.Value
	var op_135 vm.Value
	var refs_136 vm.Value
	var aux_137 vm.Value
	var needs_error_QMARK__138 vm.Value
	var ret_spec_139 vm.Value
	var head__16158_140 vm.Value
	var arg__16164_157 vm.Value
	var v158 vm.Value
	var expr_141 vm.Value
	var f_142 vm.Value
	var closed_exprs_143 vm.Value
	var bid_144 vm.Value
	var term_145 vm.Value
	var op_146 vm.Value
	var refs_147 vm.Value
	var aux_148 vm.Value
	var needs_error_QMARK__149 vm.Value
	var ret_spec_150 vm.Value
	var head__16158_151 vm.Value
	var v161 vm.Value
	var arg__16167_163 vm.Value
	var expr_164 vm.Value
	var f_165 vm.Value
	var closed_exprs_166 vm.Value
	var bid_167 vm.Value
	var term_168 vm.Value
	var op_169 vm.Value
	var refs_170 vm.Value
	var aux_171 vm.Value
	var needs_error_QMARK__172 vm.Value
	var ret_spec_173 vm.Value
	var head__16158_174 vm.Value
	var expr_176 vm.Value
	var f_177 vm.Value
	var closed_exprs_178 vm.Value
	var bid_179 vm.Value
	var term_180 vm.Value
	var op_181 vm.Value
	var refs_182 vm.Value
	var aux_183 vm.Value
	var needs_error_QMARK__184 vm.Value
	var ret_spec_185 vm.Value
	var head__16158_186 vm.Value
	var head__16168_187 vm.Value
	var arg__16174_205 vm.Value
	var v206 vm.Value
	var expr_188 vm.Value
	var f_189 vm.Value
	var closed_exprs_190 vm.Value
	var bid_191 vm.Value
	var term_192 vm.Value
	var op_193 vm.Value
	var refs_194 vm.Value
	var aux_195 vm.Value
	var needs_error_QMARK__196 vm.Value
	var ret_spec_197 vm.Value
	var head__16158_198 vm.Value
	var head__16168_199 vm.Value
	var v209 vm.Value
	var arg__16177_211 vm.Value
	var expr_212 vm.Value
	var f_213 vm.Value
	var closed_exprs_214 vm.Value
	var bid_215 vm.Value
	var term_216 vm.Value
	var op_217 vm.Value
	var refs_218 vm.Value
	var aux_219 vm.Value
	var needs_error_QMARK__220 vm.Value
	var ret_spec_221 vm.Value
	var head__16158_222 vm.Value
	var head__16168_223 vm.Value
	var arg__16178_224 vm.Value
	var v225 vm.Value
	var f_253 vm.Value
	var closed_exprs_254 vm.Value
	var bid_255 vm.Value
	var term_256 vm.Value
	var op_257 vm.Value
	var refs_258 vm.Value
	var aux_259 vm.Value
	var needs_error_QMARK__260 vm.Value
	var v273 vm.Value
	var f_261 vm.Value
	var closed_exprs_262 vm.Value
	var bid_263 vm.Value
	var term_264 vm.Value
	var op_265 vm.Value
	var refs_266 vm.Value
	var aux_267 vm.Value
	var needs_error_QMARK__268 vm.Value
	var v292 bool
	var v1014 vm.Value
	var f_1015 vm.Value
	var closed_exprs_1016 vm.Value
	var bid_1017 vm.Value
	var term_1018 vm.Value
	var op_1019 vm.Value
	var refs_1020 vm.Value
	var aux_1021 vm.Value
	var needs_error_QMARK__1022 vm.Value
	var f_275 vm.Value
	var closed_exprs_276 vm.Value
	var bid_277 vm.Value
	var term_278 vm.Value
	var op_279 vm.Value
	var refs_280 vm.Value
	var aux_281 vm.Value
	var needs_error_QMARK__282 vm.Value
	var arg__16194_312 vm.Value
	var v313 bool
	var f_283 vm.Value
	var closed_exprs_284 vm.Value
	var bid_285 vm.Value
	var term_286 vm.Value
	var op_287 vm.Value
	var refs_288 vm.Value
	var aux_289 vm.Value
	var needs_error_QMARK__290 vm.Value
	var v627 bool
	var v1004 vm.Value
	var f_1005 vm.Value
	var closed_exprs_1006 vm.Value
	var bid_1007 vm.Value
	var term_1008 vm.Value
	var op_1009 vm.Value
	var refs_1010 vm.Value
	var aux_1011 vm.Value
	var needs_error_QMARK__1012 vm.Value
	var f_294 vm.Value
	var closed_exprs_295 vm.Value
	var bid_296 vm.Value
	var term_297 vm.Value
	var op_298 vm.Value
	var refs_299 vm.Value
	var aux_300 vm.Value
	var needs_error_QMARK__301 vm.Value
	var cond_nid_318 vm.Value
	var arg__16205_320 vm.Value
	var arg__16212_323 vm.Value
	var cond_spec_324 vm.Value
	var raw_cond_326 vm.Value
	var v350 bool
	var f_302 vm.Value
	var closed_exprs_303 vm.Value
	var bid_304 vm.Value
	var term_305 vm.Value
	var op_306 vm.Value
	var refs_307 vm.Value
	var aux_308 vm.Value
	var needs_error_QMARK__309 vm.Value
	var v600 vm.Value
	var f_601 vm.Value
	var closed_exprs_602 vm.Value
	var bid_603 vm.Value
	var term_604 vm.Value
	var op_605 vm.Value
	var refs_606 vm.Value
	var aux_607 vm.Value
	var needs_error_QMARK__608 vm.Value
	var f_327 vm.Value
	var closed_exprs_328 vm.Value
	var bid_329 vm.Value
	var term_330 vm.Value
	var op_331 vm.Value
	var refs_332 vm.Value
	var aux_333 vm.Value
	var needs_error_QMARK__334 vm.Value
	var cond_nid_335 vm.Value
	var cond_spec_336 vm.Value
	var raw_cond_337 vm.Value
	var f_338 vm.Value
	var closed_exprs_339 vm.Value
	var bid_340 vm.Value
	var term_341 vm.Value
	var op_342 vm.Value
	var refs_343 vm.Value
	var aux_344 vm.Value
	var needs_error_QMARK__345 vm.Value
	var cond_nid_346 vm.Value
	var cond_spec_347 vm.Value
	var raw_cond_348 vm.Value
	var cond_expr_422 vm.Value
	var f_423 vm.Value
	var closed_exprs_424 vm.Value
	var bid_425 vm.Value
	var term_426 vm.Value
	var op_427 vm.Value
	var refs_428 vm.Value
	var aux_429 vm.Value
	var needs_error_QMARK__430 vm.Value
	var cond_nid_431 vm.Value
	var cond_spec_432 vm.Value
	var raw_cond_433 vm.Value
	var arg__16258_435 vm.Value
	var arg__16265_438 vm.Value
	var true_stmts_439 vm.Value
	var arg__16271_441 vm.Value
	var arg__16278_444 vm.Value
	var false_stmts_445 vm.Value
	var f_353 vm.Value
	var closed_exprs_354 vm.Value
	var bid_355 vm.Value
	var term_356 vm.Value
	var op_357 vm.Value
	var refs_358 vm.Value
	var aux_359 vm.Value
	var needs_error_QMARK__360 vm.Value
	var cond_nid_361 vm.Value
	var cond_spec_362 vm.Value
	var raw_cond_363 vm.Value
	var arg__16225_379 vm.Value
	var arg__16231_385 vm.Value
	var arg__16233_387 vm.Value
	var arg__16236_389 vm.Value
	var arg__16241_394 vm.Value
	var arg__16247_400 vm.Value
	var arg__16249_402 vm.Value
	var arg__16252_404 vm.Value
	var v405 vm.Value
	var f_364 vm.Value
	var closed_exprs_365 vm.Value
	var bid_366 vm.Value
	var term_367 vm.Value
	var op_368 vm.Value
	var refs_369 vm.Value
	var aux_370 vm.Value
	var needs_error_QMARK__371 vm.Value
	var cond_nid_372 vm.Value
	var cond_spec_373 vm.Value
	var raw_cond_374 vm.Value
	var v409 vm.Value
	var f_410 vm.Value
	var closed_exprs_411 vm.Value
	var bid_412 vm.Value
	var term_413 vm.Value
	var op_414 vm.Value
	var refs_415 vm.Value
	var aux_416 vm.Value
	var needs_error_QMARK__417 vm.Value
	var cond_nid_418 vm.Value
	var cond_spec_419 vm.Value
	var raw_cond_420 vm.Value
	var cond_expr_446 vm.Value
	var f_447 vm.Value
	var closed_exprs_448 vm.Value
	var bid_449 vm.Value
	var term_450 vm.Value
	var op_451 vm.Value
	var refs_452 vm.Value
	var aux_453 vm.Value
	var needs_error_QMARK__454 vm.Value
	var cond_nid_455 vm.Value
	var cond_spec_456 vm.Value
	var raw_cond_457 vm.Value
	var true_stmts_458 vm.Value
	var false_stmts_459 vm.Value
	var arg__16289_577 vm.Value
	var v578 vm.Value
	var cond_expr_460 vm.Value
	var f_461 vm.Value
	var closed_exprs_462 vm.Value
	var bid_463 vm.Value
	var term_464 vm.Value
	var op_465 vm.Value
	var refs_466 vm.Value
	var aux_467 vm.Value
	var needs_error_QMARK__468 vm.Value
	var cond_nid_469 vm.Value
	var cond_spec_470 vm.Value
	var raw_cond_471 vm.Value
	var true_stmts_472 vm.Value
	var false_stmts_473 vm.Value
	var v582 vm.Value
	var cond_expr_583 vm.Value
	var f_584 vm.Value
	var closed_exprs_585 vm.Value
	var bid_586 vm.Value
	var term_587 vm.Value
	var op_588 vm.Value
	var refs_589 vm.Value
	var aux_590 vm.Value
	var needs_error_QMARK__591 vm.Value
	var cond_nid_592 vm.Value
	var cond_spec_593 vm.Value
	var raw_cond_594 vm.Value
	var true_stmts_595 vm.Value
	var false_stmts_596 vm.Value
	var cond_expr_474 vm.Value
	var and__x_475 vm.Value
	var f_476 vm.Value
	var closed_exprs_477 vm.Value
	var bid_478 vm.Value
	var term_479 vm.Value
	var op_480 vm.Value
	var refs_481 vm.Value
	var aux_482 vm.Value
	var needs_error_QMARK__483 vm.Value
	var cond_nid_484 vm.Value
	var cond_spec_485 vm.Value
	var raw_cond_486 vm.Value
	var true_stmts_487 vm.Value
	var false_stmts_488 vm.Value
	var cond_expr_489 vm.Value
	var and__x_490 vm.Value
	var f_491 vm.Value
	var closed_exprs_492 vm.Value
	var bid_493 vm.Value
	var term_494 vm.Value
	var op_495 vm.Value
	var refs_496 vm.Value
	var aux_497 vm.Value
	var needs_error_QMARK__498 vm.Value
	var cond_nid_499 vm.Value
	var cond_spec_500 vm.Value
	var raw_cond_501 vm.Value
	var true_stmts_502 vm.Value
	var false_stmts_503 vm.Value
	var v556 vm.Value
	var cond_expr_557 vm.Value
	var and__x_558 vm.Value
	var f_559 vm.Value
	var closed_exprs_560 vm.Value
	var bid_561 vm.Value
	var term_562 vm.Value
	var op_563 vm.Value
	var refs_564 vm.Value
	var aux_565 vm.Value
	var needs_error_QMARK__566 vm.Value
	var cond_nid_567 vm.Value
	var cond_spec_568 vm.Value
	var raw_cond_569 vm.Value
	var true_stmts_570 vm.Value
	var false_stmts_571 vm.Value
	var cond_expr_505 vm.Value
	var f_506 vm.Value
	var closed_exprs_507 vm.Value
	var bid_508 vm.Value
	var term_509 vm.Value
	var op_510 vm.Value
	var refs_511 vm.Value
	var aux_512 vm.Value
	var needs_error_QMARK__513 vm.Value
	var cond_nid_514 vm.Value
	var cond_spec_515 vm.Value
	var raw_cond_516 vm.Value
	var and__x_517 vm.Value
	var true_stmts_518 vm.Value
	var false_stmts_519 vm.Value
	var cond_expr_520 vm.Value
	var f_521 vm.Value
	var closed_exprs_522 vm.Value
	var bid_523 vm.Value
	var term_524 vm.Value
	var op_525 vm.Value
	var refs_526 vm.Value
	var aux_527 vm.Value
	var needs_error_QMARK__528 vm.Value
	var cond_nid_529 vm.Value
	var cond_spec_530 vm.Value
	var raw_cond_531 vm.Value
	var and__x_532 vm.Value
	var true_stmts_533 vm.Value
	var false_stmts_534 vm.Value
	var v538 vm.Value
	var cond_expr_539 vm.Value
	var f_540 vm.Value
	var closed_exprs_541 vm.Value
	var bid_542 vm.Value
	var term_543 vm.Value
	var op_544 vm.Value
	var refs_545 vm.Value
	var aux_546 vm.Value
	var needs_error_QMARK__547 vm.Value
	var cond_nid_548 vm.Value
	var cond_spec_549 vm.Value
	var raw_cond_550 vm.Value
	var and__x_551 vm.Value
	var true_stmts_552 vm.Value
	var false_stmts_553 vm.Value
	var f_610 vm.Value
	var closed_exprs_611 vm.Value
	var bid_612 vm.Value
	var term_613 vm.Value
	var op_614 vm.Value
	var refs_615 vm.Value
	var aux_616 vm.Value
	var needs_error_QMARK__617 vm.Value
	var v648 vm.Value
	var f_618 vm.Value
	var closed_exprs_619 vm.Value
	var bid_620 vm.Value
	var term_621 vm.Value
	var op_622 vm.Value
	var refs_623 vm.Value
	var aux_624 vm.Value
	var needs_error_QMARK__625 vm.Value
	var v994 vm.Value
	var f_995 vm.Value
	var closed_exprs_996 vm.Value
	var bid_997 vm.Value
	var term_998 vm.Value
	var op_999 vm.Value
	var refs_1000 vm.Value
	var aux_1001 vm.Value
	var needs_error_QMARK__1002 vm.Value
	var f_629 vm.Value
	var closed_exprs_630 vm.Value
	var bid_631 vm.Value
	var term_632 vm.Value
	var op_633 vm.Value
	var refs_634 vm.Value
	var tail_args_635 vm.Value
	var aux_636 vm.Value
	var needs_error_QMARK__637 vm.Value
	var arg__16298_651 vm.Value
	var v652 vm.Value
	var f_638 vm.Value
	var closed_exprs_639 vm.Value
	var bid_640 vm.Value
	var term_641 vm.Value
	var op_642 vm.Value
	var refs_643 vm.Value
	var tail_args_644 vm.Value
	var aux_645 vm.Value
	var needs_error_QMARK__646 vm.Value
	var v655 vm.Value
	var fixed_arity_657 vm.Value
	var f_658 vm.Value
	var closed_exprs_659 vm.Value
	var bid_660 vm.Value
	var term_661 vm.Value
	var op_662 vm.Value
	var refs_663 vm.Value
	var tail_args_664 vm.Value
	var aux_665 vm.Value
	var needs_error_QMARK__666 vm.Value
	var i_667 int
	var remaining_668 vm.Value
	var out_669 vm.Value
	var fixed_arity_670 vm.Value
	var f_671 vm.Value
	var closed_exprs_672 vm.Value
	var v1050 string
	var v1061 string
	var v1072 string
	var v703 vm.Value
	var bid_676 vm.Value
	var term_677 vm.Value
	var op_678 vm.Value
	var refs_679 vm.Value
	var tail_args_680 vm.Value
	var aux_681 vm.Value
	var needs_error_QMARK__682 vm.Value
	var i_683 int
	var remaining_684 vm.Value
	var out_685 vm.Value
	var fixed_arity_686 vm.Value
	var f_687 vm.Value
	var closed_exprs_688 vm.Value
	var v1053 string
	var v1064 string
	var v1075 string
	var bid_689 vm.Value
	var term_690 vm.Value
	var op_691 vm.Value
	var refs_692 vm.Value
	var tail_args_693 vm.Value
	var aux_694 vm.Value
	var needs_error_QMARK__695 vm.Value
	var i_696 int
	var remaining_697 vm.Value
	var out_698 vm.Value
	var fixed_arity_699 vm.Value
	var f_700 vm.Value
	var closed_exprs_701 vm.Value
	var v1047 string
	var v1058 string
	var v1069 string
	var arg__16310_707 vm.Value
	var arg__16317_710 vm.Value
	var val_expr_711 vm.Value
	var and__x_741 vm.Value
	var assigns_890 vm.Value
	var bid_891 vm.Value
	var term_892 vm.Value
	var op_893 vm.Value
	var refs_894 vm.Value
	var tail_args_895 vm.Value
	var aux_896 vm.Value
	var needs_error_QMARK__897 vm.Value
	var i_898 int
	var remaining_899 vm.Value
	var out_900 vm.Value
	var fixed_arity_901 vm.Value
	var f_902 vm.Value
	var closed_exprs_903 vm.Value
	var bid_712 vm.Value
	var term_713 vm.Value
	var op_714 vm.Value
	var refs_715 vm.Value
	var tail_args_716 vm.Value
	var aux_717 vm.Value
	var needs_error_QMARK__718 vm.Value
	var i_719 int
	var remaining_720 vm.Value
	var out_721 vm.Value
	var fixed_arity_722 vm.Value
	var f_723 vm.Value
	var closed_exprs_724 vm.Value
	var val_expr_725 vm.Value
	var v1049 string
	var v1060 string
	var v1071 string
	var v796 vm.Value
	var bid_726 vm.Value
	var term_727 vm.Value
	var op_728 vm.Value
	var refs_729 vm.Value
	var tail_args_730 vm.Value
	var aux_731 vm.Value
	var needs_error_QMARK__732 vm.Value
	var i_733 int
	var remaining_734 vm.Value
	var out_735 vm.Value
	var fixed_arity_736 vm.Value
	var f_737 vm.Value
	var closed_exprs_738 vm.Value
	var val_expr_739 vm.Value
	var v1051 string
	var v1062 string
	var v1073 string
	var arg__16331_801 vm.Value
	var arg__16338_806 vm.Value
	var v807 vm.Value
	var lhs_809 vm.Value
	var bid_810 vm.Value
	var term_811 vm.Value
	var op_812 vm.Value
	var refs_813 vm.Value
	var tail_args_814 vm.Value
	var aux_815 vm.Value
	var needs_error_QMARK__816 vm.Value
	var i_817 int
	var remaining_818 vm.Value
	var out_819 vm.Value
	var fixed_arity_820 vm.Value
	var f_821 vm.Value
	var closed_exprs_822 vm.Value
	var val_expr_823 vm.Value
	var v1052 string
	var v1063 string
	var v1074 string
	var v855 vm.Value
	var bid_742 vm.Value
	var term_743 vm.Value
	var op_744 vm.Value
	var refs_745 vm.Value
	var tail_args_746 vm.Value
	var aux_747 vm.Value
	var needs_error_QMARK__748 vm.Value
	var i_749 int
	var remaining_750 vm.Value
	var out_751 vm.Value
	var fixed_arity_752 vm.Value
	var f_753 vm.Value
	var closed_exprs_754 vm.Value
	var val_expr_755 vm.Value
	var and__x_756 vm.Value
	var v1045 string
	var v1056 string
	var v1067 string
	var v773 bool
	var bid_757 vm.Value
	var term_758 vm.Value
	var op_759 vm.Value
	var refs_760 vm.Value
	var tail_args_761 vm.Value
	var aux_762 vm.Value
	var needs_error_QMARK__763 vm.Value
	var i_764 int
	var remaining_765 vm.Value
	var out_766 vm.Value
	var fixed_arity_767 vm.Value
	var f_768 vm.Value
	var closed_exprs_769 vm.Value
	var val_expr_770 vm.Value
	var and__x_771 vm.Value
	var v1048 string
	var v1059 string
	var v1070 string
	var v776 vm.Value
	var bid_777 vm.Value
	var term_778 vm.Value
	var op_779 vm.Value
	var refs_780 vm.Value
	var tail_args_781 vm.Value
	var aux_782 vm.Value
	var needs_error_QMARK__783 vm.Value
	var i_784 int
	var remaining_785 vm.Value
	var out_786 vm.Value
	var fixed_arity_787 vm.Value
	var f_788 vm.Value
	var closed_exprs_789 vm.Value
	var val_expr_790 vm.Value
	var and__x_791 vm.Value
	var v1044 string
	var v1055 string
	var v1066 string
	var lhs_824 vm.Value
	var bid_825 vm.Value
	var term_826 vm.Value
	var op_827 vm.Value
	var refs_828 vm.Value
	var tail_args_829 vm.Value
	var aux_830 vm.Value
	var needs_error_QMARK__831 vm.Value
	var i_832 int
	var remaining_833 vm.Value
	var out_834 vm.Value
	var fixed_arity_835 vm.Value
	var f_836 vm.Value
	var closed_exprs_837 vm.Value
	var val_expr_838 vm.Value
	var v1054 string
	var v1065 string
	var v1076 string
	var lhs_839 vm.Value
	var bid_840 vm.Value
	var term_841 vm.Value
	var op_842 vm.Value
	var refs_843 vm.Value
	var tail_args_844 vm.Value
	var aux_845 vm.Value
	var needs_error_QMARK__846 vm.Value
	var i_847 int
	var remaining_848 vm.Value
	var out_849 vm.Value
	var fixed_arity_850 vm.Value
	var f_851 vm.Value
	var closed_exprs_852 vm.Value
	var val_expr_853 vm.Value
	var v1046 string
	var v1057 string
	var v1068 string
	var v859 int
	var v861 vm.Value
	var arg__16354_865 vm.Value
	var arg__16364_870 vm.Value
	var v871 vm.Value
	var v873 vm.Value
	var lhs_874 vm.Value
	var bid_875 vm.Value
	var term_876 vm.Value
	var op_877 vm.Value
	var refs_878 vm.Value
	var tail_args_879 vm.Value
	var aux_880 vm.Value
	var needs_error_QMARK__881 vm.Value
	var i_882 int
	var remaining_883 vm.Value
	var out_884 vm.Value
	var fixed_arity_885 vm.Value
	var f_886 vm.Value
	var closed_exprs_887 vm.Value
	var val_expr_888 vm.Value
	var assigns_904 vm.Value
	var bid_905 vm.Value
	var term_906 vm.Value
	var op_907 vm.Value
	var refs_908 vm.Value
	var tail_args_909 vm.Value
	var aux_910 vm.Value
	var needs_error_QMARK__911 vm.Value
	var i_912 int
	var remaining_913 vm.Value
	var out_914 vm.Value
	var fixed_arity_915 vm.Value
	var f_916 vm.Value
	var closed_exprs_917 vm.Value
	var arg__16369_936 vm.Value
	var arg__16375_941 vm.Value
	var v942 vm.Value
	var assigns_918 vm.Value
	var bid_919 vm.Value
	var term_920 vm.Value
	var op_921 vm.Value
	var refs_922 vm.Value
	var tail_args_923 vm.Value
	var aux_924 vm.Value
	var needs_error_QMARK__925 vm.Value
	var i_926 int
	var remaining_927 vm.Value
	var out_928 vm.Value
	var fixed_arity_929 vm.Value
	var f_930 vm.Value
	var closed_exprs_931 vm.Value
	var v946 vm.Value
	var assigns_947 vm.Value
	var bid_948 vm.Value
	var term_949 vm.Value
	var op_950 vm.Value
	var refs_951 vm.Value
	var tail_args_952 vm.Value
	var aux_953 vm.Value
	var needs_error_QMARK__954 vm.Value
	var i_955 int
	var remaining_956 vm.Value
	var out_957 vm.Value
	var fixed_arity_958 vm.Value
	var f_959 vm.Value
	var closed_exprs_960 vm.Value
	var f_962 vm.Value
	var closed_exprs_963 vm.Value
	var bid_964 vm.Value
	var term_965 vm.Value
	var op_966 vm.Value
	var refs_967 vm.Value
	var aux_968 vm.Value
	var needs_error_QMARK__969 vm.Value
	var f_970 vm.Value
	var closed_exprs_971 vm.Value
	var bid_972 vm.Value
	var term_973 vm.Value
	var op_974 vm.Value
	var refs_975 vm.Value
	var aux_976 vm.Value
	var needs_error_QMARK__977 vm.Value
	var v984 vm.Value
	var f_985 vm.Value
	var closed_exprs_986 vm.Value
	var bid_987 vm.Value
	var term_988 vm.Value
	var op_989 vm.Value
	var refs_990 vm.Value
	var aux_991 vm.Value
	var needs_error_QMARK__992 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = term_4, op_6, refs_8, aux_10, needs_error_QMARK__12, v30, f_13, closed_exprs_14, bid_15, term_16, op_17, refs_18, aux_19, needs_error_QMARK__20, arg__16118_50, v51, f_21, closed_exprs_22, bid_23, term_24, op_25, refs_26, aux_27, needs_error_QMARK__28, v270, v1024, f_1025, closed_exprs_1026, bid_1027, term_1028, op_1029, refs_1030, aux_1031, needs_error_QMARK__1032, f_32, closed_exprs_33, bid_34, term_35, op_36, refs_37, aux_38, needs_error_QMARK__39, ret_spec_54, v74, f_40, closed_exprs_41, bid_42, term_43, op_44, refs_45, aux_46, needs_error_QMARK__47, v243, f_244, closed_exprs_245, bid_246, term_247, op_248, refs_249, aux_250, needs_error_QMARK__251, f_55, closed_exprs_56, bid_57, term_58, op_59, refs_60, aux_61, needs_error_QMARK__62, ret_spec_63, arg__16131_79, arg__16140_84, v85, f_64, closed_exprs_65, bid_66, term_67, op_68, refs_69, aux_70, needs_error_QMARK__71, ret_spec_72, arg__16148_90, arg__16157_95, v96, expr_98, f_99, closed_exprs_100, bid_101, term_102, op_103, refs_104, aux_105, needs_error_QMARK__106, ret_spec_107, expr_108, f_109, closed_exprs_110, bid_111, term_112, op_113, refs_114, aux_115, needs_error_QMARK__116, ret_spec_117, expr_118, f_119, closed_exprs_120, bid_121, term_122, op_123, refs_124, aux_125, needs_error_QMARK__126, ret_spec_127, v229, expr_230, f_231, closed_exprs_232, bid_233, term_234, op_235, refs_236, aux_237, needs_error_QMARK__238, ret_spec_239, expr_130, f_131, closed_exprs_132, bid_133, term_134, op_135, refs_136, aux_137, needs_error_QMARK__138, ret_spec_139, head__16158_140, arg__16164_157, v158, expr_141, f_142, closed_exprs_143, bid_144, term_145, op_146, refs_147, aux_148, needs_error_QMARK__149, ret_spec_150, head__16158_151, v161, arg__16167_163, expr_164, f_165, closed_exprs_166, bid_167, term_168, op_169, refs_170, aux_171, needs_error_QMARK__172, ret_spec_173, head__16158_174, expr_176, f_177, closed_exprs_178, bid_179, term_180, op_181, refs_182, aux_183, needs_error_QMARK__184, ret_spec_185, head__16158_186, head__16168_187, arg__16174_205, v206, expr_188, f_189, closed_exprs_190, bid_191, term_192, op_193, refs_194, aux_195, needs_error_QMARK__196, ret_spec_197, head__16158_198, head__16168_199, v209, arg__16177_211, expr_212, f_213, closed_exprs_214, bid_215, term_216, op_217, refs_218, aux_219, needs_error_QMARK__220, ret_spec_221, head__16158_222, head__16168_223, arg__16178_224, v225, f_253, closed_exprs_254, bid_255, term_256, op_257, refs_258, aux_259, needs_error_QMARK__260, v273, f_261, closed_exprs_262, bid_263, term_264, op_265, refs_266, aux_267, needs_error_QMARK__268, v292, v1014, f_1015, closed_exprs_1016, bid_1017, term_1018, op_1019, refs_1020, aux_1021, needs_error_QMARK__1022, f_275, closed_exprs_276, bid_277, term_278, op_279, refs_280, aux_281, needs_error_QMARK__282, arg__16194_312, v313, f_283, closed_exprs_284, bid_285, term_286, op_287, refs_288, aux_289, needs_error_QMARK__290, v627, v1004, f_1005, closed_exprs_1006, bid_1007, term_1008, op_1009, refs_1010, aux_1011, needs_error_QMARK__1012, f_294, closed_exprs_295, bid_296, term_297, op_298, refs_299, aux_300, needs_error_QMARK__301, cond_nid_318, arg__16205_320, arg__16212_323, cond_spec_324, raw_cond_326, v350, f_302, closed_exprs_303, bid_304, term_305, op_306, refs_307, aux_308, needs_error_QMARK__309, v600, f_601, closed_exprs_602, bid_603, term_604, op_605, refs_606, aux_607, needs_error_QMARK__608, f_327, closed_exprs_328, bid_329, term_330, op_331, refs_332, aux_333, needs_error_QMARK__334, cond_nid_335, cond_spec_336, raw_cond_337, f_338, closed_exprs_339, bid_340, term_341, op_342, refs_343, aux_344, needs_error_QMARK__345, cond_nid_346, cond_spec_347, raw_cond_348, cond_expr_422, f_423, closed_exprs_424, bid_425, term_426, op_427, refs_428, aux_429, needs_error_QMARK__430, cond_nid_431, cond_spec_432, raw_cond_433, arg__16258_435, arg__16265_438, true_stmts_439, arg__16271_441, arg__16278_444, false_stmts_445, f_353, closed_exprs_354, bid_355, term_356, op_357, refs_358, aux_359, needs_error_QMARK__360, cond_nid_361, cond_spec_362, raw_cond_363, arg__16225_379, arg__16231_385, arg__16233_387, arg__16236_389, arg__16241_394, arg__16247_400, arg__16249_402, arg__16252_404, v405, f_364, closed_exprs_365, bid_366, term_367, op_368, refs_369, aux_370, needs_error_QMARK__371, cond_nid_372, cond_spec_373, raw_cond_374, v409, f_410, closed_exprs_411, bid_412, term_413, op_414, refs_415, aux_416, needs_error_QMARK__417, cond_nid_418, cond_spec_419, raw_cond_420, cond_expr_446, f_447, closed_exprs_448, bid_449, term_450, op_451, refs_452, aux_453, needs_error_QMARK__454, cond_nid_455, cond_spec_456, raw_cond_457, true_stmts_458, false_stmts_459, arg__16289_577, v578, cond_expr_460, f_461, closed_exprs_462, bid_463, term_464, op_465, refs_466, aux_467, needs_error_QMARK__468, cond_nid_469, cond_spec_470, raw_cond_471, true_stmts_472, false_stmts_473, v582, cond_expr_583, f_584, closed_exprs_585, bid_586, term_587, op_588, refs_589, aux_590, needs_error_QMARK__591, cond_nid_592, cond_spec_593, raw_cond_594, true_stmts_595, false_stmts_596, cond_expr_474, and__x_475, f_476, closed_exprs_477, bid_478, term_479, op_480, refs_481, aux_482, needs_error_QMARK__483, cond_nid_484, cond_spec_485, raw_cond_486, true_stmts_487, false_stmts_488, cond_expr_489, and__x_490, f_491, closed_exprs_492, bid_493, term_494, op_495, refs_496, aux_497, needs_error_QMARK__498, cond_nid_499, cond_spec_500, raw_cond_501, true_stmts_502, false_stmts_503, v556, cond_expr_557, and__x_558, f_559, closed_exprs_560, bid_561, term_562, op_563, refs_564, aux_565, needs_error_QMARK__566, cond_nid_567, cond_spec_568, raw_cond_569, true_stmts_570, false_stmts_571, cond_expr_505, f_506, closed_exprs_507, bid_508, term_509, op_510, refs_511, aux_512, needs_error_QMARK__513, cond_nid_514, cond_spec_515, raw_cond_516, and__x_517, true_stmts_518, false_stmts_519, cond_expr_520, f_521, closed_exprs_522, bid_523, term_524, op_525, refs_526, aux_527, needs_error_QMARK__528, cond_nid_529, cond_spec_530, raw_cond_531, and__x_532, true_stmts_533, false_stmts_534, v538, cond_expr_539, f_540, closed_exprs_541, bid_542, term_543, op_544, refs_545, aux_546, needs_error_QMARK__547, cond_nid_548, cond_spec_549, raw_cond_550, and__x_551, true_stmts_552, false_stmts_553, f_610, closed_exprs_611, bid_612, term_613, op_614, refs_615, aux_616, needs_error_QMARK__617, v648, f_618, closed_exprs_619, bid_620, term_621, op_622, refs_623, aux_624, needs_error_QMARK__625, v994, f_995, closed_exprs_996, bid_997, term_998, op_999, refs_1000, aux_1001, needs_error_QMARK__1002, f_629, closed_exprs_630, bid_631, term_632, op_633, refs_634, tail_args_635, aux_636, needs_error_QMARK__637, arg__16298_651, v652, f_638, closed_exprs_639, bid_640, term_641, op_642, refs_643, tail_args_644, aux_645, needs_error_QMARK__646, v655, fixed_arity_657, f_658, closed_exprs_659, bid_660, term_661, op_662, refs_663, tail_args_664, aux_665, needs_error_QMARK__666, i_667, remaining_668, out_669, fixed_arity_670, f_671, closed_exprs_672, v1050, v1061, v1072, v703, bid_676, term_677, op_678, refs_679, tail_args_680, aux_681, needs_error_QMARK__682, i_683, remaining_684, out_685, fixed_arity_686, f_687, closed_exprs_688, v1053, v1064, v1075, bid_689, term_690, op_691, refs_692, tail_args_693, aux_694, needs_error_QMARK__695, i_696, remaining_697, out_698, fixed_arity_699, f_700, closed_exprs_701, v1047, v1058, v1069, arg__16310_707, arg__16317_710, val_expr_711, and__x_741, assigns_890, bid_891, term_892, op_893, refs_894, tail_args_895, aux_896, needs_error_QMARK__897, i_898, remaining_899, out_900, fixed_arity_901, f_902, closed_exprs_903, bid_712, term_713, op_714, refs_715, tail_args_716, aux_717, needs_error_QMARK__718, i_719, remaining_720, out_721, fixed_arity_722, f_723, closed_exprs_724, val_expr_725, v1049, v1060, v1071, v796, bid_726, term_727, op_728, refs_729, tail_args_730, aux_731, needs_error_QMARK__732, i_733, remaining_734, out_735, fixed_arity_736, f_737, closed_exprs_738, val_expr_739, v1051, v1062, v1073, arg__16331_801, arg__16338_806, v807, lhs_809, bid_810, term_811, op_812, refs_813, tail_args_814, aux_815, needs_error_QMARK__816, i_817, remaining_818, out_819, fixed_arity_820, f_821, closed_exprs_822, val_expr_823, v1052, v1063, v1074, v855, bid_742, term_743, op_744, refs_745, tail_args_746, aux_747, needs_error_QMARK__748, i_749, remaining_750, out_751, fixed_arity_752, f_753, closed_exprs_754, val_expr_755, and__x_756, v1045, v1056, v1067, v773, bid_757, term_758, op_759, refs_760, tail_args_761, aux_762, needs_error_QMARK__763, i_764, remaining_765, out_766, fixed_arity_767, f_768, closed_exprs_769, val_expr_770, and__x_771, v1048, v1059, v1070, v776, bid_777, term_778, op_779, refs_780, tail_args_781, aux_782, needs_error_QMARK__783, i_784, remaining_785, out_786, fixed_arity_787, f_788, closed_exprs_789, val_expr_790, and__x_791, v1044, v1055, v1066, lhs_824, bid_825, term_826, op_827, refs_828, tail_args_829, aux_830, needs_error_QMARK__831, i_832, remaining_833, out_834, fixed_arity_835, f_836, closed_exprs_837, val_expr_838, v1054, v1065, v1076, lhs_839, bid_840, term_841, op_842, refs_843, tail_args_844, aux_845, needs_error_QMARK__846, i_847, remaining_848, out_849, fixed_arity_850, f_851, closed_exprs_852, val_expr_853, v1046, v1057, v1068, v859, v861, arg__16354_865, arg__16364_870, v871, v873, lhs_874, bid_875, term_876, op_877, refs_878, tail_args_879, aux_880, needs_error_QMARK__881, i_882, remaining_883, out_884, fixed_arity_885, f_886, closed_exprs_887, val_expr_888, assigns_904, bid_905, term_906, op_907, refs_908, tail_args_909, aux_910, needs_error_QMARK__911, i_912, remaining_913, out_914, fixed_arity_915, f_916, closed_exprs_917, arg__16369_936, arg__16375_941, v942, assigns_918, bid_919, term_920, op_921, refs_922, tail_args_923, aux_924, needs_error_QMARK__925, i_926, remaining_927, out_928, fixed_arity_929, f_930, closed_exprs_931, v946, assigns_947, bid_948, term_949, op_950, refs_951, tail_args_952, aux_953, needs_error_QMARK__954, i_955, remaining_956, out_957, fixed_arity_958, f_959, closed_exprs_960, f_962, closed_exprs_963, bid_964, term_965, op_966, refs_967, aux_968, needs_error_QMARK__969, f_970, closed_exprs_971, bid_972, term_973, op_974, refs_975, aux_976, needs_error_QMARK__977, v984, f_985, closed_exprs_986, bid_987, term_988, op_989, refs_990, aux_991, needs_error_QMARK__992
	term_4, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-term").Deref(), []vm.Value{arg2, arg0})
	if callErr != nil {
		return nil, callErr
	}
	op_6, callErr = rt.InvokeValue(rt.LookupVar("ir", "op").Deref(), []vm.Value{term_4, arg0})
	if callErr != nil {
		return nil, callErr
	}
	refs_8, callErr = rt.InvokeValue(rt.LookupVar("ir", "refs").Deref(), []vm.Value{term_4, arg0})
	if callErr != nil {
		return nil, callErr
	}
	aux_10, callErr = rt.InvokeValue(rt.LookupVar("ir", "aux").Deref(), []vm.Value{term_4, arg0})
	if callErr != nil {
		return nil, callErr
	}
	needs_error_QMARK__12, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "function-needs-error?").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v30 = op_6 == vm.Keyword("return")
	if v30 {
		f_13 = arg0
		closed_exprs_14 = arg1
		bid_15 = arg2
		term_16 = term_4
		op_17 = op_6
		refs_18 = refs_8
		aux_19 = aux_10
		needs_error_QMARK__20 = needs_error_QMARK__12
		goto b1
	} else {
		f_21 = arg0
		closed_exprs_22 = arg1
		bid_23 = arg2
		term_24 = term_4
		op_25 = op_6
		refs_26 = refs_8
		aux_27 = aux_10
		needs_error_QMARK__28 = needs_error_QMARK__12
		goto b2
	}
b1:
	;
	arg__16118_50, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{refs_18})
	if callErr != nil {
		return nil, callErr
	}
	v51 = arg__16118_50 == vm.Int(1)
	if v51 {
		f_32 = f_13
		closed_exprs_33 = closed_exprs_14
		bid_34 = bid_15
		term_35 = term_16
		op_36 = op_17
		refs_37 = refs_18
		aux_38 = aux_19
		needs_error_QMARK__39 = needs_error_QMARK__20
		goto b4
	} else {
		f_40 = f_13
		closed_exprs_41 = closed_exprs_14
		bid_42 = bid_15
		term_43 = term_16
		op_44 = op_17
		refs_45 = refs_18
		aux_46 = aux_19
		needs_error_QMARK__47 = needs_error_QMARK__20
		goto b5
	}
b2:
	;
	v270 = op_25 == vm.Keyword("branch")
	if v270 {
		f_253 = f_21
		closed_exprs_254 = closed_exprs_22
		bid_255 = bid_23
		term_256 = term_24
		op_257 = op_25
		refs_258 = refs_26
		aux_259 = aux_27
		needs_error_QMARK__260 = needs_error_QMARK__28
		goto b19
	} else {
		f_261 = f_21
		closed_exprs_262 = closed_exprs_22
		bid_263 = bid_23
		term_264 = term_24
		op_265 = op_25
		refs_266 = refs_26
		aux_267 = aux_27
		needs_error_QMARK__268 = needs_error_QMARK__28
		goto b20
	}
b3:
	;
	return v1024, nil
b4:
	;
	ret_spec_54, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "function-return-spec").Deref(), []vm.Value{f_32})
	if callErr != nil {
		return nil, callErr
	}
	v74 = ret_spec_54 == vm.String("vm.Value")
	if v74 {
		f_55 = f_32
		closed_exprs_56 = closed_exprs_33
		bid_57 = bid_34
		term_58 = term_35
		op_59 = op_36
		refs_60 = refs_37
		aux_61 = aux_38
		needs_error_QMARK__62 = needs_error_QMARK__39
		ret_spec_63 = ret_spec_54
		goto b7
	} else {
		f_64 = f_32
		closed_exprs_65 = closed_exprs_33
		bid_66 = bid_34
		term_67 = term_35
		op_68 = op_36
		refs_69 = refs_37
		aux_70 = aux_38
		needs_error_QMARK__71 = needs_error_QMARK__39
		ret_spec_72 = ret_spec_54
		goto b8
	}
b5:
	;
	v243 = vm.NIL
	f_244 = f_40
	closed_exprs_245 = closed_exprs_41
	bid_246 = bid_42
	term_247 = term_43
	op_248 = op_44
	refs_249 = refs_45
	aux_250 = aux_46
	needs_error_QMARK__251 = needs_error_QMARK__47
	goto b6
b6:
	;
	v1024 = v243
	f_1025 = f_244
	closed_exprs_1026 = closed_exprs_245
	bid_1027 = bid_246
	term_1028 = term_247
	op_1029 = op_248
	refs_1030 = refs_249
	aux_1031 = aux_250
	needs_error_QMARK__1032 = needs_error_QMARK__251
	goto b3
b7:
	;
	arg__16131_79, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{refs_60, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	arg__16140_84, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{refs_60, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	v85, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "box-as-value").Deref(), []vm.Value{f_55, closed_exprs_56, arg__16140_84})
	if callErr != nil {
		return nil, callErr
	}
	expr_98 = v85
	f_99 = f_55
	closed_exprs_100 = closed_exprs_56
	bid_101 = bid_57
	term_102 = term_58
	op_103 = op_59
	refs_104 = refs_60
	aux_105 = aux_61
	needs_error_QMARK__106 = needs_error_QMARK__62
	ret_spec_107 = ret_spec_63
	goto b9
b8:
	;
	arg__16148_90, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{refs_69, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	arg__16157_95, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{refs_69, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	v96, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "value-expr").Deref(), []vm.Value{f_64, closed_exprs_65, arg__16157_95})
	if callErr != nil {
		return nil, callErr
	}
	expr_98 = v96
	f_99 = f_64
	closed_exprs_100 = closed_exprs_65
	bid_101 = bid_66
	term_102 = term_67
	op_103 = op_68
	refs_104 = refs_69
	aux_105 = aux_70
	needs_error_QMARK__106 = needs_error_QMARK__71
	ret_spec_107 = ret_spec_72
	goto b9
b9:
	;
	if vm.IsTruthy(expr_98) {
		expr_108 = expr_98
		f_109 = f_99
		closed_exprs_110 = closed_exprs_100
		bid_111 = bid_101
		term_112 = term_102
		op_113 = op_103
		refs_114 = refs_104
		aux_115 = aux_105
		needs_error_QMARK__116 = needs_error_QMARK__106
		ret_spec_117 = ret_spec_107
		goto b10
	} else {
		expr_118 = expr_98
		f_119 = f_99
		closed_exprs_120 = closed_exprs_100
		bid_121 = bid_101
		term_122 = term_102
		op_123 = op_103
		refs_124 = refs_104
		aux_125 = aux_105
		needs_error_QMARK__126 = needs_error_QMARK__106
		ret_spec_127 = ret_spec_107
		goto b11
	}
b10:
	;
	if vm.IsTruthy(needs_error_QMARK__116) {
		expr_130 = expr_108
		f_131 = f_109
		closed_exprs_132 = closed_exprs_110
		bid_133 = bid_111
		term_134 = term_112
		op_135 = op_113
		refs_136 = refs_114
		aux_137 = aux_115
		needs_error_QMARK__138 = needs_error_QMARK__116
		ret_spec_139 = ret_spec_117
		head__16158_140 = rt.LookupVar("clojure.core", "vector").Deref()
		goto b13
	} else {
		expr_141 = expr_108
		f_142 = f_109
		closed_exprs_143 = closed_exprs_110
		bid_144 = bid_111
		term_145 = term_112
		op_146 = op_113
		refs_147 = refs_114
		aux_148 = aux_115
		needs_error_QMARK__149 = needs_error_QMARK__116
		ret_spec_150 = ret_spec_117
		head__16158_151 = rt.LookupVar("clojure.core", "vector").Deref()
		goto b14
	}
b11:
	;
	v229 = vm.NIL
	expr_230 = expr_118
	f_231 = f_119
	closed_exprs_232 = closed_exprs_120
	bid_233 = bid_121
	term_234 = term_122
	op_235 = op_123
	refs_236 = refs_124
	aux_237 = aux_125
	needs_error_QMARK__238 = needs_error_QMARK__126
	ret_spec_239 = ret_spec_127
	goto b12
b12:
	;
	v243 = v229
	f_244 = f_231
	closed_exprs_245 = closed_exprs_232
	bid_246 = bid_233
	term_247 = term_234
	op_248 = op_235
	refs_249 = refs_236
	aux_250 = aux_237
	needs_error_QMARK__251 = needs_error_QMARK__238
	goto b6
b13:
	;
	arg__16164_157, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("nil")})
	if callErr != nil {
		return nil, callErr
	}
	v158, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{expr_130, arg__16164_157})
	if callErr != nil {
		return nil, callErr
	}
	arg__16167_163 = v158
	expr_164 = expr_130
	f_165 = f_131
	closed_exprs_166 = closed_exprs_132
	bid_167 = bid_133
	term_168 = term_134
	op_169 = op_135
	refs_170 = refs_136
	aux_171 = aux_137
	needs_error_QMARK__172 = needs_error_QMARK__138
	ret_spec_173 = ret_spec_139
	head__16158_174 = head__16158_140
	goto b15
b14:
	;
	v161, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{expr_141})
	if callErr != nil {
		return nil, callErr
	}
	arg__16167_163 = v161
	expr_164 = expr_141
	f_165 = f_142
	closed_exprs_166 = closed_exprs_143
	bid_167 = bid_144
	term_168 = term_145
	op_169 = op_146
	refs_170 = refs_147
	aux_171 = aux_148
	needs_error_QMARK__172 = needs_error_QMARK__149
	ret_spec_173 = ret_spec_150
	head__16158_174 = head__16158_151
	goto b15
b15:
	;
	if vm.IsTruthy(needs_error_QMARK__172) {
		expr_176 = expr_164
		f_177 = f_165
		closed_exprs_178 = closed_exprs_166
		bid_179 = bid_167
		term_180 = term_168
		op_181 = op_169
		refs_182 = refs_170
		aux_183 = aux_171
		needs_error_QMARK__184 = needs_error_QMARK__172
		ret_spec_185 = ret_spec_173
		head__16158_186 = head__16158_174
		head__16168_187 = rt.LookupVar("gogen", "return-stmt").Deref()
		goto b16
	} else {
		expr_188 = expr_164
		f_189 = f_165
		closed_exprs_190 = closed_exprs_166
		bid_191 = bid_167
		term_192 = term_168
		op_193 = op_169
		refs_194 = refs_170
		aux_195 = aux_171
		needs_error_QMARK__196 = needs_error_QMARK__172
		ret_spec_197 = ret_spec_173
		head__16158_198 = head__16158_174
		head__16168_199 = rt.LookupVar("gogen", "return-stmt").Deref()
		goto b17
	}
b16:
	;
	arg__16174_205, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("nil")})
	if callErr != nil {
		return nil, callErr
	}
	v206, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{expr_176, arg__16174_205})
	if callErr != nil {
		return nil, callErr
	}
	arg__16177_211 = v206
	expr_212 = expr_176
	f_213 = f_177
	closed_exprs_214 = closed_exprs_178
	bid_215 = bid_179
	term_216 = term_180
	op_217 = op_181
	refs_218 = refs_182
	aux_219 = aux_183
	needs_error_QMARK__220 = needs_error_QMARK__184
	ret_spec_221 = ret_spec_185
	head__16158_222 = head__16158_186
	head__16168_223 = head__16168_187
	goto b18
b17:
	;
	v209, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{expr_188})
	if callErr != nil {
		return nil, callErr
	}
	arg__16177_211 = v209
	expr_212 = expr_188
	f_213 = f_189
	closed_exprs_214 = closed_exprs_190
	bid_215 = bid_191
	term_216 = term_192
	op_217 = op_193
	refs_218 = refs_194
	aux_219 = aux_195
	needs_error_QMARK__220 = needs_error_QMARK__196
	ret_spec_221 = ret_spec_197
	head__16158_222 = head__16158_198
	head__16168_223 = head__16168_199
	goto b18
b18:
	;
	arg__16178_224, callErr = rt.InvokeValue(head__16168_223, []vm.Value{arg__16177_211})
	if callErr != nil {
		return nil, callErr
	}
	v225, callErr = rt.InvokeValue(head__16158_222, []vm.Value{arg__16178_224})
	if callErr != nil {
		return nil, callErr
	}
	v229 = v225
	expr_230 = expr_164
	f_231 = f_165
	closed_exprs_232 = closed_exprs_166
	bid_233 = bid_167
	term_234 = term_168
	op_235 = op_169
	refs_236 = refs_170
	aux_237 = aux_171
	needs_error_QMARK__238 = needs_error_QMARK__172
	ret_spec_239 = ret_spec_173
	goto b12
b19:
	;
	v273, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "transfer-stmts").Deref(), []vm.Value{f_253, closed_exprs_254, aux_259})
	if callErr != nil {
		return nil, callErr
	}
	v1014 = v273
	f_1015 = f_253
	closed_exprs_1016 = closed_exprs_254
	bid_1017 = bid_255
	term_1018 = term_256
	op_1019 = op_257
	refs_1020 = refs_258
	aux_1021 = aux_259
	needs_error_QMARK__1022 = needs_error_QMARK__260
	goto b21
b20:
	;
	v292 = op_265 == vm.Keyword("branch-if")
	if v292 {
		f_275 = f_261
		closed_exprs_276 = closed_exprs_262
		bid_277 = bid_263
		term_278 = term_264
		op_279 = op_265
		refs_280 = refs_266
		aux_281 = aux_267
		needs_error_QMARK__282 = needs_error_QMARK__268
		goto b22
	} else {
		f_283 = f_261
		closed_exprs_284 = closed_exprs_262
		bid_285 = bid_263
		term_286 = term_264
		op_287 = op_265
		refs_288 = refs_266
		aux_289 = aux_267
		needs_error_QMARK__290 = needs_error_QMARK__268
		goto b23
	}
b21:
	;
	v1024 = v1014
	f_1025 = f_1015
	closed_exprs_1026 = closed_exprs_1016
	bid_1027 = bid_1017
	term_1028 = term_1018
	op_1029 = op_1019
	refs_1030 = refs_1020
	aux_1031 = aux_1021
	needs_error_QMARK__1032 = needs_error_QMARK__1022
	goto b3
b22:
	;
	arg__16194_312, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{refs_280})
	if callErr != nil {
		return nil, callErr
	}
	v313 = arg__16194_312 == vm.Int(1)
	if v313 {
		f_294 = f_275
		closed_exprs_295 = closed_exprs_276
		bid_296 = bid_277
		term_297 = term_278
		op_298 = op_279
		refs_299 = refs_280
		aux_300 = aux_281
		needs_error_QMARK__301 = needs_error_QMARK__282
		goto b25
	} else {
		f_302 = f_275
		closed_exprs_303 = closed_exprs_276
		bid_304 = bid_277
		term_305 = term_278
		op_306 = op_279
		refs_307 = refs_280
		aux_308 = aux_281
		needs_error_QMARK__309 = needs_error_QMARK__282
		goto b26
	}
b23:
	;
	v627 = op_287 == vm.Keyword("tail-call")
	if v627 {
		f_610 = f_283
		closed_exprs_611 = closed_exprs_284
		bid_612 = bid_285
		term_613 = term_286
		op_614 = op_287
		refs_615 = refs_288
		aux_616 = aux_289
		needs_error_QMARK__617 = needs_error_QMARK__290
		goto b43
	} else {
		f_618 = f_283
		closed_exprs_619 = closed_exprs_284
		bid_620 = bid_285
		term_621 = term_286
		op_622 = op_287
		refs_623 = refs_288
		aux_624 = aux_289
		needs_error_QMARK__625 = needs_error_QMARK__290
		goto b44
	}
b24:
	;
	v1014 = v1004
	f_1015 = f_1005
	closed_exprs_1016 = closed_exprs_1006
	bid_1017 = bid_1007
	term_1018 = term_1008
	op_1019 = op_1009
	refs_1020 = refs_1010
	aux_1021 = aux_1011
	needs_error_QMARK__1022 = needs_error_QMARK__1012
	goto b21
b25:
	;
	cond_nid_318, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{refs_299, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	arg__16205_320, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{cond_nid_318, f_294})
	if callErr != nil {
		return nil, callErr
	}
	arg__16212_323, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{cond_nid_318, f_294})
	if callErr != nil {
		return nil, callErr
	}
	cond_spec_324, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "go-type-spec").Deref(), []vm.Value{arg__16212_323})
	if callErr != nil {
		return nil, callErr
	}
	raw_cond_326, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "value-expr").Deref(), []vm.Value{f_294, closed_exprs_295, cond_nid_318})
	if callErr != nil {
		return nil, callErr
	}
	v350 = cond_spec_324 == vm.String("bool")
	if v350 {
		f_327 = f_294
		closed_exprs_328 = closed_exprs_295
		bid_329 = bid_296
		term_330 = term_297
		op_331 = op_298
		refs_332 = refs_299
		aux_333 = aux_300
		needs_error_QMARK__334 = needs_error_QMARK__301
		cond_nid_335 = cond_nid_318
		cond_spec_336 = cond_spec_324
		raw_cond_337 = raw_cond_326
		goto b28
	} else {
		f_338 = f_294
		closed_exprs_339 = closed_exprs_295
		bid_340 = bid_296
		term_341 = term_297
		op_342 = op_298
		refs_343 = refs_299
		aux_344 = aux_300
		needs_error_QMARK__345 = needs_error_QMARK__301
		cond_nid_346 = cond_nid_318
		cond_spec_347 = cond_spec_324
		raw_cond_348 = raw_cond_326
		goto b29
	}
b26:
	;
	v600 = vm.NIL
	f_601 = f_302
	closed_exprs_602 = closed_exprs_303
	bid_603 = bid_304
	term_604 = term_305
	op_605 = op_306
	refs_606 = refs_307
	aux_607 = aux_308
	needs_error_QMARK__608 = needs_error_QMARK__309
	goto b27
b27:
	;
	v1004 = v600
	f_1005 = f_601
	closed_exprs_1006 = closed_exprs_602
	bid_1007 = bid_603
	term_1008 = term_604
	op_1009 = op_605
	refs_1010 = refs_606
	aux_1011 = aux_607
	needs_error_QMARK__1012 = needs_error_QMARK__608
	goto b24
b28:
	;
	cond_expr_422 = raw_cond_337
	f_423 = f_327
	closed_exprs_424 = closed_exprs_328
	bid_425 = bid_329
	term_426 = term_330
	op_427 = op_331
	refs_428 = refs_332
	aux_429 = aux_333
	needs_error_QMARK__430 = needs_error_QMARK__334
	cond_nid_431 = cond_nid_335
	cond_spec_432 = cond_spec_336
	raw_cond_433 = raw_cond_337
	goto b30
b29:
	;
	if vm.IsTruthy(raw_cond_348) {
		f_353 = f_338
		closed_exprs_354 = closed_exprs_339
		bid_355 = bid_340
		term_356 = term_341
		op_357 = op_342
		refs_358 = refs_343
		aux_359 = aux_344
		needs_error_QMARK__360 = needs_error_QMARK__345
		cond_nid_361 = cond_nid_346
		cond_spec_362 = cond_spec_347
		raw_cond_363 = raw_cond_348
		goto b31
	} else {
		f_364 = f_338
		closed_exprs_365 = closed_exprs_339
		bid_366 = bid_340
		term_367 = term_341
		op_368 = op_342
		refs_369 = refs_343
		aux_370 = aux_344
		needs_error_QMARK__371 = needs_error_QMARK__345
		cond_nid_372 = cond_nid_346
		cond_spec_373 = cond_spec_347
		raw_cond_374 = raw_cond_348
		goto b32
	}
b30:
	;
	arg__16258_435, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-true").Deref(), []vm.Value{aux_429})
	if callErr != nil {
		return nil, callErr
	}
	arg__16265_438, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-true").Deref(), []vm.Value{aux_429})
	if callErr != nil {
		return nil, callErr
	}
	true_stmts_439, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "transfer-stmts").Deref(), []vm.Value{f_423, closed_exprs_424, arg__16265_438})
	if callErr != nil {
		return nil, callErr
	}
	arg__16271_441, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-false").Deref(), []vm.Value{aux_429})
	if callErr != nil {
		return nil, callErr
	}
	arg__16278_444, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-false").Deref(), []vm.Value{aux_429})
	if callErr != nil {
		return nil, callErr
	}
	false_stmts_445, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "transfer-stmts").Deref(), []vm.Value{f_423, closed_exprs_424, arg__16278_444})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(cond_expr_422) {
		cond_expr_474 = cond_expr_422
		and__x_475 = cond_expr_422
		f_476 = f_423
		closed_exprs_477 = closed_exprs_424
		bid_478 = bid_425
		term_479 = term_426
		op_480 = op_427
		refs_481 = refs_428
		aux_482 = aux_429
		needs_error_QMARK__483 = needs_error_QMARK__430
		cond_nid_484 = cond_nid_431
		cond_spec_485 = cond_spec_432
		raw_cond_486 = raw_cond_433
		true_stmts_487 = true_stmts_439
		false_stmts_488 = false_stmts_445
		goto b37
	} else {
		cond_expr_489 = cond_expr_422
		and__x_490 = cond_expr_422
		f_491 = f_423
		closed_exprs_492 = closed_exprs_424
		bid_493 = bid_425
		term_494 = term_426
		op_495 = op_427
		refs_496 = refs_428
		aux_497 = aux_429
		needs_error_QMARK__498 = needs_error_QMARK__430
		cond_nid_499 = cond_nid_431
		cond_spec_500 = cond_spec_432
		raw_cond_501 = raw_cond_433
		true_stmts_502 = true_stmts_439
		false_stmts_503 = false_stmts_445
		goto b38
	}
b31:
	;
	arg__16225_379, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("vm")})
	if callErr != nil {
		return nil, callErr
	}
	arg__16231_385, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("vm")})
	if callErr != nil {
		return nil, callErr
	}
	arg__16233_387, callErr = rt.InvokeValue(rt.LookupVar("gogen", "field-sel").Deref(), []vm.Value{arg__16231_385, vm.String("IsTruthy")})
	if callErr != nil {
		return nil, callErr
	}
	arg__16236_389, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{raw_cond_363})
	if callErr != nil {
		return nil, callErr
	}
	arg__16241_394, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("vm")})
	if callErr != nil {
		return nil, callErr
	}
	arg__16247_400, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("vm")})
	if callErr != nil {
		return nil, callErr
	}
	arg__16249_402, callErr = rt.InvokeValue(rt.LookupVar("gogen", "field-sel").Deref(), []vm.Value{arg__16247_400, vm.String("IsTruthy")})
	if callErr != nil {
		return nil, callErr
	}
	arg__16252_404, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{raw_cond_363})
	if callErr != nil {
		return nil, callErr
	}
	v405, callErr = rt.InvokeValue(rt.LookupVar("gogen", "call").Deref(), []vm.Value{arg__16249_402, arg__16252_404})
	if callErr != nil {
		return nil, callErr
	}
	v409 = v405
	f_410 = f_353
	closed_exprs_411 = closed_exprs_354
	bid_412 = bid_355
	term_413 = term_356
	op_414 = op_357
	refs_415 = refs_358
	aux_416 = aux_359
	needs_error_QMARK__417 = needs_error_QMARK__360
	cond_nid_418 = cond_nid_361
	cond_spec_419 = cond_spec_362
	raw_cond_420 = raw_cond_363
	goto b33
b32:
	;
	v409 = vm.NIL
	f_410 = f_364
	closed_exprs_411 = closed_exprs_365
	bid_412 = bid_366
	term_413 = term_367
	op_414 = op_368
	refs_415 = refs_369
	aux_416 = aux_370
	needs_error_QMARK__417 = needs_error_QMARK__371
	cond_nid_418 = cond_nid_372
	cond_spec_419 = cond_spec_373
	raw_cond_420 = raw_cond_374
	goto b33
b33:
	;
	cond_expr_422 = v409
	f_423 = f_410
	closed_exprs_424 = closed_exprs_411
	bid_425 = bid_412
	term_426 = term_413
	op_427 = op_414
	refs_428 = refs_415
	aux_429 = aux_416
	needs_error_QMARK__430 = needs_error_QMARK__417
	cond_nid_431 = cond_nid_418
	cond_spec_432 = cond_spec_419
	raw_cond_433 = raw_cond_420
	goto b30
b34:
	;
	arg__16289_577, callErr = rt.InvokeValue(rt.LookupVar("gogen", "if-stmt").Deref(), []vm.Value{vm.NIL, cond_expr_446, true_stmts_458, false_stmts_459})
	if callErr != nil {
		return nil, callErr
	}
	v578, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__16289_577})
	if callErr != nil {
		return nil, callErr
	}
	v582 = v578
	cond_expr_583 = cond_expr_446
	f_584 = f_447
	closed_exprs_585 = closed_exprs_448
	bid_586 = bid_449
	term_587 = term_450
	op_588 = op_451
	refs_589 = refs_452
	aux_590 = aux_453
	needs_error_QMARK__591 = needs_error_QMARK__454
	cond_nid_592 = cond_nid_455
	cond_spec_593 = cond_spec_456
	raw_cond_594 = raw_cond_457
	true_stmts_595 = true_stmts_458
	false_stmts_596 = false_stmts_459
	goto b36
b35:
	;
	v582 = vm.NIL
	cond_expr_583 = cond_expr_460
	f_584 = f_461
	closed_exprs_585 = closed_exprs_462
	bid_586 = bid_463
	term_587 = term_464
	op_588 = op_465
	refs_589 = refs_466
	aux_590 = aux_467
	needs_error_QMARK__591 = needs_error_QMARK__468
	cond_nid_592 = cond_nid_469
	cond_spec_593 = cond_spec_470
	raw_cond_594 = raw_cond_471
	true_stmts_595 = true_stmts_472
	false_stmts_596 = false_stmts_473
	goto b36
b36:
	;
	v600 = v582
	f_601 = f_584
	closed_exprs_602 = closed_exprs_585
	bid_603 = bid_586
	term_604 = term_587
	op_605 = op_588
	refs_606 = refs_589
	aux_607 = aux_590
	needs_error_QMARK__608 = needs_error_QMARK__591
	goto b27
b37:
	;
	if vm.IsTruthy(true_stmts_487) {
		cond_expr_505 = cond_expr_474
		f_506 = f_476
		closed_exprs_507 = closed_exprs_477
		bid_508 = bid_478
		term_509 = term_479
		op_510 = op_480
		refs_511 = refs_481
		aux_512 = aux_482
		needs_error_QMARK__513 = needs_error_QMARK__483
		cond_nid_514 = cond_nid_484
		cond_spec_515 = cond_spec_485
		raw_cond_516 = raw_cond_486
		and__x_517 = true_stmts_487
		true_stmts_518 = true_stmts_487
		false_stmts_519 = false_stmts_488
		goto b40
	} else {
		cond_expr_520 = cond_expr_474
		f_521 = f_476
		closed_exprs_522 = closed_exprs_477
		bid_523 = bid_478
		term_524 = term_479
		op_525 = op_480
		refs_526 = refs_481
		aux_527 = aux_482
		needs_error_QMARK__528 = needs_error_QMARK__483
		cond_nid_529 = cond_nid_484
		cond_spec_530 = cond_spec_485
		raw_cond_531 = raw_cond_486
		and__x_532 = true_stmts_487
		true_stmts_533 = true_stmts_487
		false_stmts_534 = false_stmts_488
		goto b41
	}
b38:
	;
	v556 = and__x_490
	cond_expr_557 = cond_expr_489
	and__x_558 = and__x_490
	f_559 = f_491
	closed_exprs_560 = closed_exprs_492
	bid_561 = bid_493
	term_562 = term_494
	op_563 = op_495
	refs_564 = refs_496
	aux_565 = aux_497
	needs_error_QMARK__566 = needs_error_QMARK__498
	cond_nid_567 = cond_nid_499
	cond_spec_568 = cond_spec_500
	raw_cond_569 = raw_cond_501
	true_stmts_570 = true_stmts_502
	false_stmts_571 = false_stmts_503
	goto b39
b39:
	;
	if vm.IsTruthy(v556) {
		cond_expr_446 = cond_expr_557
		f_447 = f_559
		closed_exprs_448 = closed_exprs_560
		bid_449 = bid_561
		term_450 = term_562
		op_451 = op_563
		refs_452 = refs_564
		aux_453 = aux_565
		needs_error_QMARK__454 = needs_error_QMARK__566
		cond_nid_455 = cond_nid_567
		cond_spec_456 = cond_spec_568
		raw_cond_457 = raw_cond_569
		true_stmts_458 = true_stmts_570
		false_stmts_459 = false_stmts_571
		goto b34
	} else {
		cond_expr_460 = cond_expr_557
		f_461 = f_559
		closed_exprs_462 = closed_exprs_560
		bid_463 = bid_561
		term_464 = term_562
		op_465 = op_563
		refs_466 = refs_564
		aux_467 = aux_565
		needs_error_QMARK__468 = needs_error_QMARK__566
		cond_nid_469 = cond_nid_567
		cond_spec_470 = cond_spec_568
		raw_cond_471 = raw_cond_569
		true_stmts_472 = true_stmts_570
		false_stmts_473 = false_stmts_571
		goto b35
	}
b40:
	;
	v538 = false_stmts_519
	cond_expr_539 = cond_expr_505
	f_540 = f_506
	closed_exprs_541 = closed_exprs_507
	bid_542 = bid_508
	term_543 = term_509
	op_544 = op_510
	refs_545 = refs_511
	aux_546 = aux_512
	needs_error_QMARK__547 = needs_error_QMARK__513
	cond_nid_548 = cond_nid_514
	cond_spec_549 = cond_spec_515
	raw_cond_550 = raw_cond_516
	and__x_551 = and__x_517
	true_stmts_552 = true_stmts_518
	false_stmts_553 = false_stmts_519
	goto b42
b41:
	;
	v538 = and__x_532
	cond_expr_539 = cond_expr_520
	f_540 = f_521
	closed_exprs_541 = closed_exprs_522
	bid_542 = bid_523
	term_543 = term_524
	op_544 = op_525
	refs_545 = refs_526
	aux_546 = aux_527
	needs_error_QMARK__547 = needs_error_QMARK__528
	cond_nid_548 = cond_nid_529
	cond_spec_549 = cond_spec_530
	raw_cond_550 = raw_cond_531
	and__x_551 = and__x_532
	true_stmts_552 = true_stmts_533
	false_stmts_553 = false_stmts_534
	goto b42
b42:
	;
	v556 = v538
	cond_expr_557 = cond_expr_539
	and__x_558 = and__x_475
	f_559 = f_540
	closed_exprs_560 = closed_exprs_541
	bid_561 = bid_542
	term_562 = term_543
	op_563 = op_544
	refs_564 = refs_545
	aux_565 = aux_546
	needs_error_QMARK__566 = needs_error_QMARK__547
	cond_nid_567 = cond_nid_548
	cond_spec_568 = cond_spec_549
	raw_cond_569 = raw_cond_550
	true_stmts_570 = true_stmts_552
	false_stmts_571 = false_stmts_553
	goto b39
b43:
	;
	v648, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-variadic?").Deref(), []vm.Value{f_610})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v648) {
		f_629 = f_610
		closed_exprs_630 = closed_exprs_611
		bid_631 = bid_612
		term_632 = term_613
		op_633 = op_614
		refs_634 = refs_615
		tail_args_635 = refs_615
		aux_636 = aux_616
		needs_error_QMARK__637 = needs_error_QMARK__617
		goto b46
	} else {
		f_638 = f_610
		closed_exprs_639 = closed_exprs_611
		bid_640 = bid_612
		term_641 = term_613
		op_642 = op_614
		refs_643 = refs_615
		tail_args_644 = refs_615
		aux_645 = aux_616
		needs_error_QMARK__646 = needs_error_QMARK__617
		goto b47
	}
b44:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		f_962 = f_618
		closed_exprs_963 = closed_exprs_619
		bid_964 = bid_620
		term_965 = term_621
		op_966 = op_622
		refs_967 = refs_623
		aux_968 = aux_624
		needs_error_QMARK__969 = needs_error_QMARK__625
		goto b65
	} else {
		f_970 = f_618
		closed_exprs_971 = closed_exprs_619
		bid_972 = bid_620
		term_973 = term_621
		op_974 = op_622
		refs_975 = refs_623
		aux_976 = aux_624
		needs_error_QMARK__977 = needs_error_QMARK__625
		goto b66
	}
b45:
	;
	v1004 = v994
	f_1005 = f_995
	closed_exprs_1006 = closed_exprs_996
	bid_1007 = bid_997
	term_1008 = term_998
	op_1009 = op_999
	refs_1010 = refs_1000
	aux_1011 = aux_1001
	needs_error_QMARK__1012 = needs_error_QMARK__1002
	goto b24
b46:
	;
	arg__16298_651, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-arity").Deref(), []vm.Value{f_629})
	if callErr != nil {
		return nil, callErr
	}
	v652 = rt.SubValue(arg__16298_651, vm.Int(1))
	fixed_arity_657 = v652
	f_658 = f_629
	closed_exprs_659 = closed_exprs_630
	bid_660 = bid_631
	term_661 = term_632
	op_662 = op_633
	refs_663 = refs_634
	tail_args_664 = tail_args_635
	aux_665 = aux_636
	needs_error_QMARK__666 = needs_error_QMARK__637
	goto b48
b47:
	;
	v655, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-arity").Deref(), []vm.Value{f_638})
	if callErr != nil {
		return nil, callErr
	}
	fixed_arity_657 = v655
	f_658 = f_638
	closed_exprs_659 = closed_exprs_639
	bid_660 = bid_640
	term_661 = term_641
	op_662 = op_642
	refs_663 = refs_643
	tail_args_664 = tail_args_644
	aux_665 = aux_645
	needs_error_QMARK__666 = needs_error_QMARK__646
	goto b48
b48:
	;
	i_667 = 0
	remaining_668 = tail_args_664
	out_669 = vm.NewArrayVector([]vm.Value{})
	fixed_arity_670 = fixed_arity_657
	f_671 = f_658
	closed_exprs_672 = closed_exprs_659
	v1050 = "args0"
	v1061 = "a"
	v1072 = "="
	goto b49
b49:
	;
	v703, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "empty?").Deref(), []vm.Value{remaining_668})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v703) {
		bid_676 = bid_660
		term_677 = term_661
		op_678 = op_662
		refs_679 = refs_663
		tail_args_680 = tail_args_664
		aux_681 = aux_665
		needs_error_QMARK__682 = needs_error_QMARK__666
		i_683 = i_667
		remaining_684 = remaining_668
		out_685 = out_669
		fixed_arity_686 = fixed_arity_670
		f_687 = f_671
		closed_exprs_688 = closed_exprs_672
		v1053 = v1050
		v1064 = v1061
		v1075 = v1072
		goto b50
	} else {
		bid_689 = bid_660
		term_690 = term_661
		op_691 = op_662
		refs_692 = refs_663
		tail_args_693 = tail_args_664
		aux_694 = aux_665
		needs_error_QMARK__695 = needs_error_QMARK__666
		i_696 = i_667
		remaining_697 = remaining_668
		out_698 = out_669
		fixed_arity_699 = fixed_arity_670
		f_700 = f_671
		closed_exprs_701 = closed_exprs_672
		v1047 = v1050
		v1058 = v1061
		v1069 = v1072
		goto b51
	}
b50:
	;
	assigns_890 = out_685
	bid_891 = bid_676
	term_892 = term_677
	op_893 = op_678
	refs_894 = refs_679
	tail_args_895 = tail_args_680
	aux_896 = aux_681
	needs_error_QMARK__897 = needs_error_QMARK__682
	i_898 = i_683
	remaining_899 = remaining_684
	out_900 = out_685
	fixed_arity_901 = fixed_arity_686
	f_902 = f_687
	closed_exprs_903 = closed_exprs_688
	goto b52
b51:
	;
	arg__16310_707, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{remaining_697})
	if callErr != nil {
		return nil, callErr
	}
	arg__16317_710, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{remaining_697})
	if callErr != nil {
		return nil, callErr
	}
	val_expr_711, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "box-as-value").Deref(), []vm.Value{f_700, closed_exprs_701, arg__16317_710})
	if callErr != nil {
		return nil, callErr
	}
	and__x_741, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-variadic?").Deref(), []vm.Value{f_700})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(and__x_741) {
		bid_742 = bid_689
		term_743 = term_690
		op_744 = op_691
		refs_745 = refs_692
		tail_args_746 = tail_args_693
		aux_747 = aux_694
		needs_error_QMARK__748 = needs_error_QMARK__695
		i_749 = i_696
		remaining_750 = remaining_697
		out_751 = out_698
		fixed_arity_752 = fixed_arity_699
		f_753 = f_700
		closed_exprs_754 = closed_exprs_701
		val_expr_755 = val_expr_711
		and__x_756 = and__x_741
		v1045 = v1047
		v1056 = v1058
		v1067 = v1069
		goto b56
	} else {
		bid_757 = bid_689
		term_758 = term_690
		op_759 = op_691
		refs_760 = refs_692
		tail_args_761 = tail_args_693
		aux_762 = aux_694
		needs_error_QMARK__763 = needs_error_QMARK__695
		i_764 = i_696
		remaining_765 = remaining_697
		out_766 = out_698
		fixed_arity_767 = fixed_arity_699
		f_768 = f_700
		closed_exprs_769 = closed_exprs_701
		val_expr_770 = val_expr_711
		and__x_771 = and__x_741
		v1048 = v1047
		v1059 = v1058
		v1070 = v1069
		goto b57
	}
b52:
	;
	if vm.IsTruthy(assigns_890) {
		assigns_904 = assigns_890
		bid_905 = bid_891
		term_906 = term_892
		op_907 = op_893
		refs_908 = refs_894
		tail_args_909 = tail_args_895
		aux_910 = aux_896
		needs_error_QMARK__911 = needs_error_QMARK__897
		i_912 = i_898
		remaining_913 = remaining_899
		out_914 = out_900
		fixed_arity_915 = fixed_arity_901
		f_916 = f_902
		closed_exprs_917 = closed_exprs_903
		goto b62
	} else {
		assigns_918 = assigns_890
		bid_919 = bid_891
		term_920 = term_892
		op_921 = op_893
		refs_922 = refs_894
		tail_args_923 = tail_args_895
		aux_924 = aux_896
		needs_error_QMARK__925 = needs_error_QMARK__897
		i_926 = i_898
		remaining_927 = remaining_899
		out_928 = out_900
		fixed_arity_929 = fixed_arity_901
		f_930 = f_902
		closed_exprs_931 = closed_exprs_903
		goto b63
	}
b53:
	;
	v796, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String(v1049)})
	if callErr != nil {
		return nil, callErr
	}
	lhs_809 = v796
	bid_810 = bid_712
	term_811 = term_713
	op_812 = op_714
	refs_813 = refs_715
	tail_args_814 = tail_args_716
	aux_815 = aux_717
	needs_error_QMARK__816 = needs_error_QMARK__718
	i_817 = i_719
	remaining_818 = remaining_720
	out_819 = out_721
	fixed_arity_820 = fixed_arity_722
	f_821 = f_723
	closed_exprs_822 = closed_exprs_724
	val_expr_823 = val_expr_725
	v1052 = v1049
	v1063 = v1060
	v1074 = v1071
	goto b55
b54:
	;
	arg__16331_801, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(v1062), vm.Int(i_733)})
	if callErr != nil {
		return nil, callErr
	}
	arg__16338_806, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(v1062), vm.Int(i_733)})
	if callErr != nil {
		return nil, callErr
	}
	v807, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{arg__16338_806})
	if callErr != nil {
		return nil, callErr
	}
	lhs_809 = v807
	bid_810 = bid_726
	term_811 = term_727
	op_812 = op_728
	refs_813 = refs_729
	tail_args_814 = tail_args_730
	aux_815 = aux_731
	needs_error_QMARK__816 = needs_error_QMARK__732
	i_817 = i_733
	remaining_818 = remaining_734
	out_819 = out_735
	fixed_arity_820 = fixed_arity_736
	f_821 = f_737
	closed_exprs_822 = closed_exprs_738
	val_expr_823 = val_expr_739
	v1052 = v1051
	v1063 = v1062
	v1074 = v1073
	goto b55
b55:
	;
	v855, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{val_expr_823})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v855) {
		lhs_824 = lhs_809
		bid_825 = bid_810
		term_826 = term_811
		op_827 = op_812
		refs_828 = refs_813
		tail_args_829 = tail_args_814
		aux_830 = aux_815
		needs_error_QMARK__831 = needs_error_QMARK__816
		i_832 = i_817
		remaining_833 = remaining_818
		out_834 = out_819
		fixed_arity_835 = fixed_arity_820
		f_836 = f_821
		closed_exprs_837 = closed_exprs_822
		val_expr_838 = val_expr_823
		v1054 = v1052
		v1065 = v1063
		v1076 = v1074
		goto b59
	} else {
		lhs_839 = lhs_809
		bid_840 = bid_810
		term_841 = term_811
		op_842 = op_812
		refs_843 = refs_813
		tail_args_844 = tail_args_814
		aux_845 = aux_815
		needs_error_QMARK__846 = needs_error_QMARK__816
		i_847 = i_817
		remaining_848 = remaining_818
		out_849 = out_819
		fixed_arity_850 = fixed_arity_820
		f_851 = f_821
		closed_exprs_852 = closed_exprs_822
		val_expr_853 = val_expr_823
		v1046 = v1052
		v1057 = v1063
		v1068 = v1074
		goto b60
	}
b56:
	;
	v773 = vm.Int(i_749) == fixed_arity_752
	v776 = vm.Boolean(v773)
	bid_777 = bid_742
	term_778 = term_743
	op_779 = op_744
	refs_780 = refs_745
	tail_args_781 = tail_args_746
	aux_782 = aux_747
	needs_error_QMARK__783 = needs_error_QMARK__748
	i_784 = i_749
	remaining_785 = remaining_750
	out_786 = out_751
	fixed_arity_787 = fixed_arity_752
	f_788 = f_753
	closed_exprs_789 = closed_exprs_754
	val_expr_790 = val_expr_755
	and__x_791 = and__x_756
	v1044 = v1045
	v1055 = v1056
	v1066 = v1067
	goto b58
b57:
	;
	v776 = and__x_771
	bid_777 = bid_757
	term_778 = term_758
	op_779 = op_759
	refs_780 = refs_760
	tail_args_781 = tail_args_761
	aux_782 = aux_762
	needs_error_QMARK__783 = needs_error_QMARK__763
	i_784 = i_764
	remaining_785 = remaining_765
	out_786 = out_766
	fixed_arity_787 = fixed_arity_767
	f_788 = f_768
	closed_exprs_789 = closed_exprs_769
	val_expr_790 = val_expr_770
	and__x_791 = and__x_771
	v1044 = v1048
	v1055 = v1059
	v1066 = v1070
	goto b58
b58:
	;
	if vm.IsTruthy(v776) {
		bid_712 = bid_777
		term_713 = term_778
		op_714 = op_779
		refs_715 = refs_780
		tail_args_716 = tail_args_781
		aux_717 = aux_782
		needs_error_QMARK__718 = needs_error_QMARK__783
		i_719 = i_784
		remaining_720 = remaining_785
		out_721 = out_786
		fixed_arity_722 = fixed_arity_787
		f_723 = f_788
		closed_exprs_724 = closed_exprs_789
		val_expr_725 = val_expr_790
		v1049 = v1044
		v1060 = v1055
		v1071 = v1066
		goto b53
	} else {
		bid_726 = bid_777
		term_727 = term_778
		op_728 = op_779
		refs_729 = refs_780
		tail_args_730 = tail_args_781
		aux_731 = aux_782
		needs_error_QMARK__732 = needs_error_QMARK__783
		i_733 = i_784
		remaining_734 = remaining_785
		out_735 = out_786
		fixed_arity_736 = fixed_arity_787
		f_737 = f_788
		closed_exprs_738 = closed_exprs_789
		val_expr_739 = val_expr_790
		v1051 = v1044
		v1062 = v1055
		v1073 = v1066
		goto b54
	}
b59:
	;
	v873 = vm.NIL
	lhs_874 = lhs_824
	bid_875 = bid_825
	term_876 = term_826
	op_877 = op_827
	refs_878 = refs_828
	tail_args_879 = tail_args_829
	aux_880 = aux_830
	needs_error_QMARK__881 = needs_error_QMARK__831
	i_882 = i_832
	remaining_883 = remaining_833
	out_884 = out_834
	fixed_arity_885 = fixed_arity_835
	f_886 = f_836
	closed_exprs_887 = closed_exprs_837
	val_expr_888 = val_expr_838
	goto b61
b60:
	;
	v859 = i_847 + 1
	v861, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{remaining_848})
	if callErr != nil {
		return nil, callErr
	}
	arg__16354_865, callErr = rt.InvokeValue(rt.LookupVar("gogen", "assign").Deref(), []vm.Value{vm.String(v1068), lhs_839, val_expr_853})
	if callErr != nil {
		return nil, callErr
	}
	arg__16364_870, callErr = rt.InvokeValue(rt.LookupVar("gogen", "assign").Deref(), []vm.Value{vm.String(v1068), lhs_839, val_expr_853})
	if callErr != nil {
		return nil, callErr
	}
	v871, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{out_849, arg__16364_870})
	if callErr != nil {
		return nil, callErr
	}
	i_667 = v859
	remaining_668 = v861
	out_669 = v871
	fixed_arity_670 = fixed_arity_850
	f_671 = f_851
	closed_exprs_672 = closed_exprs_852
	v1050 = v1046
	v1061 = v1057
	v1072 = v1068
	goto b49
b61:
	;
	assigns_890 = v873
	bid_891 = bid_875
	term_892 = term_876
	op_893 = op_877
	refs_894 = refs_878
	tail_args_895 = tail_args_879
	aux_896 = aux_880
	needs_error_QMARK__897 = needs_error_QMARK__881
	i_898 = i_882
	remaining_899 = remaining_883
	out_900 = out_884
	fixed_arity_901 = fixed_arity_885
	f_902 = f_886
	closed_exprs_903 = closed_exprs_887
	goto b52
b62:
	;
	arg__16369_936, callErr = rt.InvokeValue(rt.LookupVar("gogen", "goto-stmt").Deref(), []vm.Value{vm.String("func_entry")})
	if callErr != nil {
		return nil, callErr
	}
	arg__16375_941, callErr = rt.InvokeValue(rt.LookupVar("gogen", "goto-stmt").Deref(), []vm.Value{vm.String("func_entry")})
	if callErr != nil {
		return nil, callErr
	}
	v942, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{assigns_904, arg__16375_941})
	if callErr != nil {
		return nil, callErr
	}
	v946 = v942
	assigns_947 = assigns_904
	bid_948 = bid_905
	term_949 = term_906
	op_950 = op_907
	refs_951 = refs_908
	tail_args_952 = tail_args_909
	aux_953 = aux_910
	needs_error_QMARK__954 = needs_error_QMARK__911
	i_955 = i_912
	remaining_956 = remaining_913
	out_957 = out_914
	fixed_arity_958 = fixed_arity_915
	f_959 = f_916
	closed_exprs_960 = closed_exprs_917
	goto b64
b63:
	;
	v946 = vm.NIL
	assigns_947 = assigns_918
	bid_948 = bid_919
	term_949 = term_920
	op_950 = op_921
	refs_951 = refs_922
	tail_args_952 = tail_args_923
	aux_953 = aux_924
	needs_error_QMARK__954 = needs_error_QMARK__925
	i_955 = i_926
	remaining_956 = remaining_927
	out_957 = out_928
	fixed_arity_958 = fixed_arity_929
	f_959 = f_930
	closed_exprs_960 = closed_exprs_931
	goto b64
b64:
	;
	v994 = v946
	f_995 = f_959
	closed_exprs_996 = closed_exprs_960
	bid_997 = bid_948
	term_998 = term_949
	op_999 = op_950
	refs_1000 = refs_951
	aux_1001 = aux_953
	needs_error_QMARK__1002 = needs_error_QMARK__954
	goto b45
b65:
	;
	v984 = vm.NIL
	f_985 = f_962
	closed_exprs_986 = closed_exprs_963
	bid_987 = bid_964
	term_988 = term_965
	op_989 = op_966
	refs_990 = refs_967
	aux_991 = aux_968
	needs_error_QMARK__992 = needs_error_QMARK__969
	goto b67
b66:
	;
	v984 = vm.NIL
	f_985 = f_970
	closed_exprs_986 = closed_exprs_971
	bid_987 = bid_972
	term_988 = term_973
	op_989 = op_974
	refs_990 = refs_975
	aux_991 = aux_976
	needs_error_QMARK__992 = needs_error_QMARK__977
	goto b67
b67:
	;
	v994 = v984
	f_995 = f_985
	closed_exprs_996 = closed_exprs_986
	bid_997 = bid_987
	term_998 = term_988
	op_999 = op_989
	refs_1000 = refs_990
	aux_1001 = aux_991
	needs_error_QMARK__1002 = needs_error_QMARK__992
	goto b45
}
func multi_fn_template_QMARK_(arg0 vm.Value) (vm.Value, error) {
	var and__x_2 vm.Value
	var aux_3 vm.Value
	var and__x_4 vm.Value
	var arg__16381_9 vm.Value
	var v11 bool
	var aux_5 vm.Value
	var and__x_6 vm.Value
	var v14 vm.Value
	var aux_15 vm.Value
	var and__x_16 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _ = and__x_2, aux_3, and__x_4, arg__16381_9, v11, aux_5, and__x_6, v14, aux_15, and__x_16
	and__x_2, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map?").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(and__x_2) {
		aux_3 = arg0
		and__x_4 = and__x_2
		goto b1
	} else {
		aux_5 = arg0
		and__x_6 = and__x_2
		goto b2
	}
b1:
	;
	arg__16381_9, callErr = rt.InvokeValue(vm.Keyword("kind"), []vm.Value{aux_3})
	if callErr != nil {
		return nil, callErr
	}
	v11 = arg__16381_9 == vm.Keyword("multi-fn-template")
	v14 = vm.Boolean(v11)
	aux_15 = aux_3
	and__x_16 = and__x_4
	goto b3
b2:
	;
	v14 = and__x_6
	aux_15 = aux_5
	and__x_16 = and__x_6
	goto b3
b3:
	;
	return v14, nil
}
func params_for(arg0 vm.Value) (vm.Value, error) {
	var arity_2 vm.Value
	var variadic_QMARK__4 vm.Value
	var f_5 vm.Value
	var arity_6 vm.Value
	var variadic_QMARK__7 vm.Value
	var v12 vm.Value
	var f_8 vm.Value
	var arity_9 vm.Value
	var variadic_QMARK__10 vm.Value
	var fixed_arity_15 vm.Value
	var f_16 vm.Value
	var arity_17 vm.Value
	var variadic_QMARK__18 vm.Value
	var or__x_20 vm.Value
	var fixed_arity_21 vm.Value
	var f_22 vm.Value
	var arity_23 vm.Value
	var variadic_QMARK__24 vm.Value
	var or__x_25 vm.Value
	var fixed_arity_26 vm.Value
	var f_27 vm.Value
	var arity_28 vm.Value
	var variadic_QMARK__29 vm.Value
	var or__x_30 vm.Value
	var arg_types_35 vm.Value
	var fixed_arity_36 vm.Value
	var f_37 vm.Value
	var arity_38 vm.Value
	var variadic_QMARK__39 vm.Value
	var or__x_40 vm.Value
	var i_41 int
	var out_42 vm.Value
	var arg_types_43 vm.Value
	var variadic_QMARK__44 vm.Value
	var fixed_arity_45 vm.Value
	var f_46 vm.Value
	var v321 vm.Value
	var v335 vm.Value
	var v349 string
	var v64 bool
	var arity_50 vm.Value
	var i_51 int
	var out_52 vm.Value
	var arg_types_53 vm.Value
	var variadic_QMARK__54 vm.Value
	var fixed_arity_55 vm.Value
	var f_56 vm.Value
	var v323 vm.Value
	var v337 vm.Value
	var v351 string
	var arity_57 vm.Value
	var i_58 int
	var out_59 vm.Value
	var arg_types_60 vm.Value
	var variadic_QMARK__61 vm.Value
	var fixed_arity_62 vm.Value
	var f_63 vm.Value
	var v316 vm.Value
	var v330 vm.Value
	var v344 string
	var load_arg_t_119 vm.Value
	var arg__16501_137 vm.Value
	var v138 bool
	var v288 vm.Value
	var arity_289 vm.Value
	var i_290 int
	var out_291 vm.Value
	var arg_types_292 vm.Value
	var variadic_QMARK__293 vm.Value
	var fixed_arity_294 vm.Value
	var f_295 vm.Value
	var arity_66 vm.Value
	var i_67 int
	var out_68 vm.Value
	var arg_types_69 vm.Value
	var variadic_QMARK__70 vm.Value
	var fixed_arity_71 vm.Value
	var f_72 vm.Value
	var arg__16470_85 vm.Value
	var arg__16476_91 vm.Value
	var arg__16477_92 vm.Value
	var arg__16484_98 vm.Value
	var arg__16490_104 vm.Value
	var arg__16491_105 vm.Value
	var v106 vm.Value
	var arity_73 vm.Value
	var i_74 int
	var out_75 vm.Value
	var arg_types_76 vm.Value
	var variadic_QMARK__77 vm.Value
	var fixed_arity_78 vm.Value
	var f_79 vm.Value
	var v109 vm.Value
	var arity_110 vm.Value
	var i_111 int
	var out_112 vm.Value
	var arg_types_113 vm.Value
	var variadic_QMARK__114 vm.Value
	var fixed_arity_115 vm.Value
	var f_116 vm.Value
	var arity_120 vm.Value
	var i_121 int
	var out_122 vm.Value
	var arg_types_123 vm.Value
	var variadic_QMARK__124 vm.Value
	var fixed_arity_125 vm.Value
	var f_126 vm.Value
	var load_arg_t_127 vm.Value
	var v320 vm.Value
	var v334 vm.Value
	var v348 string
	var v141 vm.Value
	var arity_128 vm.Value
	var i_129 int
	var out_130 vm.Value
	var arg_types_131 vm.Value
	var variadic_QMARK__132 vm.Value
	var fixed_arity_133 vm.Value
	var f_134 vm.Value
	var load_arg_t_135 vm.Value
	var v313 vm.Value
	var v327 vm.Value
	var v341 string
	var meta_t_145 vm.Value
	var arity_146 vm.Value
	var i_147 int
	var out_148 vm.Value
	var arg_types_149 vm.Value
	var variadic_QMARK__150 vm.Value
	var fixed_arity_151 vm.Value
	var f_152 vm.Value
	var load_arg_t_153 vm.Value
	var v310 vm.Value
	var v324 vm.Value
	var v338 string
	var v175 vm.Value
	var meta_t_154 vm.Value
	var arity_155 vm.Value
	var i_156 int
	var out_157 vm.Value
	var arg_types_158 vm.Value
	var variadic_QMARK__159 vm.Value
	var fixed_arity_160 vm.Value
	var f_161 vm.Value
	var load_arg_t_162 vm.Value
	var v317 vm.Value
	var v331 vm.Value
	var v345 string
	var meta_t_163 vm.Value
	var arity_164 vm.Value
	var i_165 int
	var out_166 vm.Value
	var arg_types_167 vm.Value
	var variadic_QMARK__168 vm.Value
	var fixed_arity_169 vm.Value
	var f_170 vm.Value
	var load_arg_t_171 vm.Value
	var v315 vm.Value
	var v329 vm.Value
	var v343 string
	var t_212 vm.Value
	var meta_t_213 vm.Value
	var arity_214 vm.Value
	var i_215 int
	var out_216 vm.Value
	var arg_types_217 vm.Value
	var variadic_QMARK__218 vm.Value
	var fixed_arity_219 vm.Value
	var f_220 vm.Value
	var load_arg_t_221 vm.Value
	var v319 vm.Value
	var v333 vm.Value
	var v347 string
	var go_type_223 vm.Value
	var v247 vm.Value
	var meta_t_178 vm.Value
	var arity_179 vm.Value
	var i_180 int
	var out_181 vm.Value
	var arg_types_182 vm.Value
	var variadic_QMARK__183 vm.Value
	var fixed_arity_184 vm.Value
	var f_185 vm.Value
	var load_arg_t_186 vm.Value
	var v311 vm.Value
	var v325 vm.Value
	var v339 string
	var meta_t_187 vm.Value
	var arity_188 vm.Value
	var i_189 int
	var out_190 vm.Value
	var arg_types_191 vm.Value
	var variadic_QMARK__192 vm.Value
	var fixed_arity_193 vm.Value
	var f_194 vm.Value
	var load_arg_t_195 vm.Value
	var v314 vm.Value
	var v328 vm.Value
	var v342 string
	var v201 vm.Value
	var meta_t_202 vm.Value
	var arity_203 vm.Value
	var i_204 int
	var out_205 vm.Value
	var arg_types_206 vm.Value
	var variadic_QMARK__207 vm.Value
	var fixed_arity_208 vm.Value
	var f_209 vm.Value
	var load_arg_t_210 vm.Value
	var v318 vm.Value
	var v332 vm.Value
	var v346 string
	var t_224 vm.Value
	var meta_t_225 vm.Value
	var arity_226 vm.Value
	var i_227 int
	var out_228 vm.Value
	var arg_types_229 vm.Value
	var variadic_QMARK__230 vm.Value
	var fixed_arity_231 vm.Value
	var f_232 vm.Value
	var load_arg_t_233 vm.Value
	var go_type_234 vm.Value
	var v322 vm.Value
	var v336 vm.Value
	var v350 string
	var t_235 vm.Value
	var meta_t_236 vm.Value
	var arity_237 vm.Value
	var i_238 int
	var out_239 vm.Value
	var arg_types_240 vm.Value
	var variadic_QMARK__241 vm.Value
	var fixed_arity_242 vm.Value
	var f_243 vm.Value
	var load_arg_t_244 vm.Value
	var go_type_245 vm.Value
	var v312 vm.Value
	var v326 vm.Value
	var v340 string
	var v251 int
	var arg__16525_255 vm.Value
	var arg__16533_260 vm.Value
	var arg__16535_261 vm.Value
	var arg__16543_266 vm.Value
	var arg__16551_271 vm.Value
	var arg__16553_272 vm.Value
	var v273 vm.Value
	var v275 vm.Value
	var t_276 vm.Value
	var meta_t_277 vm.Value
	var arity_278 vm.Value
	var i_279 int
	var out_280 vm.Value
	var arg_types_281 vm.Value
	var variadic_QMARK__282 vm.Value
	var fixed_arity_283 vm.Value
	var f_284 vm.Value
	var load_arg_t_285 vm.Value
	var go_type_286 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arity_2, variadic_QMARK__4, f_5, arity_6, variadic_QMARK__7, v12, f_8, arity_9, variadic_QMARK__10, fixed_arity_15, f_16, arity_17, variadic_QMARK__18, or__x_20, fixed_arity_21, f_22, arity_23, variadic_QMARK__24, or__x_25, fixed_arity_26, f_27, arity_28, variadic_QMARK__29, or__x_30, arg_types_35, fixed_arity_36, f_37, arity_38, variadic_QMARK__39, or__x_40, i_41, out_42, arg_types_43, variadic_QMARK__44, fixed_arity_45, f_46, v321, v335, v349, v64, arity_50, i_51, out_52, arg_types_53, variadic_QMARK__54, fixed_arity_55, f_56, v323, v337, v351, arity_57, i_58, out_59, arg_types_60, variadic_QMARK__61, fixed_arity_62, f_63, v316, v330, v344, load_arg_t_119, arg__16501_137, v138, v288, arity_289, i_290, out_291, arg_types_292, variadic_QMARK__293, fixed_arity_294, f_295, arity_66, i_67, out_68, arg_types_69, variadic_QMARK__70, fixed_arity_71, f_72, arg__16470_85, arg__16476_91, arg__16477_92, arg__16484_98, arg__16490_104, arg__16491_105, v106, arity_73, i_74, out_75, arg_types_76, variadic_QMARK__77, fixed_arity_78, f_79, v109, arity_110, i_111, out_112, arg_types_113, variadic_QMARK__114, fixed_arity_115, f_116, arity_120, i_121, out_122, arg_types_123, variadic_QMARK__124, fixed_arity_125, f_126, load_arg_t_127, v320, v334, v348, v141, arity_128, i_129, out_130, arg_types_131, variadic_QMARK__132, fixed_arity_133, f_134, load_arg_t_135, v313, v327, v341, meta_t_145, arity_146, i_147, out_148, arg_types_149, variadic_QMARK__150, fixed_arity_151, f_152, load_arg_t_153, v310, v324, v338, v175, meta_t_154, arity_155, i_156, out_157, arg_types_158, variadic_QMARK__159, fixed_arity_160, f_161, load_arg_t_162, v317, v331, v345, meta_t_163, arity_164, i_165, out_166, arg_types_167, variadic_QMARK__168, fixed_arity_169, f_170, load_arg_t_171, v315, v329, v343, t_212, meta_t_213, arity_214, i_215, out_216, arg_types_217, variadic_QMARK__218, fixed_arity_219, f_220, load_arg_t_221, v319, v333, v347, go_type_223, v247, meta_t_178, arity_179, i_180, out_181, arg_types_182, variadic_QMARK__183, fixed_arity_184, f_185, load_arg_t_186, v311, v325, v339, meta_t_187, arity_188, i_189, out_190, arg_types_191, variadic_QMARK__192, fixed_arity_193, f_194, load_arg_t_195, v314, v328, v342, v201, meta_t_202, arity_203, i_204, out_205, arg_types_206, variadic_QMARK__207, fixed_arity_208, f_209, load_arg_t_210, v318, v332, v346, t_224, meta_t_225, arity_226, i_227, out_228, arg_types_229, variadic_QMARK__230, fixed_arity_231, f_232, load_arg_t_233, go_type_234, v322, v336, v350, t_235, meta_t_236, arity_237, i_238, out_239, arg_types_240, variadic_QMARK__241, fixed_arity_242, f_243, load_arg_t_244, go_type_245, v312, v326, v340, v251, arg__16525_255, arg__16533_260, arg__16535_261, arg__16543_266, arg__16551_271, arg__16553_272, v273, v275, t_276, meta_t_277, arity_278, i_279, out_280, arg_types_281, variadic_QMARK__282, fixed_arity_283, f_284, load_arg_t_285, go_type_286
	arity_2, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-arity").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	variadic_QMARK__4, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-variadic?").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(variadic_QMARK__4) {
		f_5 = arg0
		arity_6 = arity_2
		variadic_QMARK__7 = variadic_QMARK__4
		goto b1
	} else {
		f_8 = arg0
		arity_9 = arity_2
		variadic_QMARK__10 = variadic_QMARK__4
		goto b2
	}
b1:
	;
	v12 = rt.SubValue(arity_6, vm.Int(1))
	fixed_arity_15 = v12
	f_16 = f_5
	arity_17 = arity_6
	variadic_QMARK__18 = variadic_QMARK__7
	goto b3
b2:
	;
	fixed_arity_15 = arity_9
	f_16 = f_8
	arity_17 = arity_9
	variadic_QMARK__18 = variadic_QMARK__10
	goto b3
b3:
	;
	or__x_20, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-arg-types").Deref(), []vm.Value{f_16})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(or__x_20) {
		fixed_arity_21 = fixed_arity_15
		f_22 = f_16
		arity_23 = arity_17
		variadic_QMARK__24 = variadic_QMARK__18
		or__x_25 = or__x_20
		goto b4
	} else {
		fixed_arity_26 = fixed_arity_15
		f_27 = f_16
		arity_28 = arity_17
		variadic_QMARK__29 = variadic_QMARK__18
		or__x_30 = or__x_20
		goto b5
	}
b4:
	;
	arg_types_35 = or__x_25
	fixed_arity_36 = fixed_arity_21
	f_37 = f_22
	arity_38 = arity_23
	variadic_QMARK__39 = variadic_QMARK__24
	or__x_40 = or__x_25
	goto b6
b5:
	;
	arg_types_35 = vm.NewArrayVector([]vm.Value{})
	fixed_arity_36 = fixed_arity_26
	f_37 = f_27
	arity_38 = arity_28
	variadic_QMARK__39 = variadic_QMARK__29
	or__x_40 = or__x_30
	goto b6
b6:
	;
	i_41 = 0
	out_42 = vm.NewArrayVector([]vm.Value{})
	arg_types_43 = arg_types_35
	variadic_QMARK__44 = variadic_QMARK__39
	fixed_arity_45 = fixed_arity_36
	f_46 = f_37
	v321 = vm.Keyword("unknown")
	v335 = vm.Keyword("else")
	v349 = "arg"
	goto b7
b7:
	;
	v64 = rt.GeValue(vm.Int(i_41), fixed_arity_45)
	if v64 {
		arity_50 = arity_38
		i_51 = i_41
		out_52 = out_42
		arg_types_53 = arg_types_43
		variadic_QMARK__54 = variadic_QMARK__44
		fixed_arity_55 = fixed_arity_45
		f_56 = f_46
		v323 = v321
		v337 = v335
		v351 = v349
		goto b8
	} else {
		arity_57 = arity_38
		i_58 = i_41
		out_59 = out_42
		arg_types_60 = arg_types_43
		variadic_QMARK__61 = variadic_QMARK__44
		fixed_arity_62 = fixed_arity_45
		f_63 = f_46
		v316 = v321
		v330 = v335
		v344 = v349
		goto b9
	}
b8:
	;
	if vm.IsTruthy(variadic_QMARK__54) {
		arity_66 = arity_50
		i_67 = i_51
		out_68 = out_52
		arg_types_69 = arg_types_53
		variadic_QMARK__70 = variadic_QMARK__54
		fixed_arity_71 = fixed_arity_55
		f_72 = f_56
		goto b11
	} else {
		arity_73 = arity_50
		i_74 = i_51
		out_75 = out_52
		arg_types_76 = arg_types_53
		variadic_QMARK__77 = variadic_QMARK__54
		fixed_arity_78 = fixed_arity_55
		f_79 = f_56
		goto b12
	}
b9:
	;
	load_arg_t_119, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-load-arg-type").Deref(), []vm.Value{f_63, vm.Int(i_58)})
	if callErr != nil {
		return nil, callErr
	}
	arg__16501_137, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{arg_types_60})
	if callErr != nil {
		return nil, callErr
	}
	v138 = rt.LtValue(vm.Int(i_58), arg__16501_137)
	if v138 {
		arity_120 = arity_57
		i_121 = i_58
		out_122 = out_59
		arg_types_123 = arg_types_60
		variadic_QMARK__124 = variadic_QMARK__61
		fixed_arity_125 = fixed_arity_62
		f_126 = f_63
		load_arg_t_127 = load_arg_t_119
		v320 = v316
		v334 = v330
		v348 = v344
		goto b14
	} else {
		arity_128 = arity_57
		i_129 = i_58
		out_130 = out_59
		arg_types_131 = arg_types_60
		variadic_QMARK__132 = variadic_QMARK__61
		fixed_arity_133 = fixed_arity_62
		f_134 = f_63
		load_arg_t_135 = load_arg_t_119
		v313 = v316
		v327 = v330
		v341 = v344
		goto b15
	}
b10:
	;
	return v288, nil
b11:
	;
	arg__16470_85, callErr = rt.InvokeValue(rt.LookupVar("gogen", "type").Deref(), []vm.Value{vm.String("vm.Value")})
	if callErr != nil {
		return nil, callErr
	}
	arg__16476_91, callErr = rt.InvokeValue(rt.LookupVar("gogen", "type").Deref(), []vm.Value{vm.String("vm.Value")})
	if callErr != nil {
		return nil, callErr
	}
	arg__16477_92, callErr = rt.InvokeValue(rt.LookupVar("gogen", "variadic-param").Deref(), []vm.Value{vm.String("args"), arg__16476_91})
	if callErr != nil {
		return nil, callErr
	}
	arg__16484_98, callErr = rt.InvokeValue(rt.LookupVar("gogen", "type").Deref(), []vm.Value{vm.String("vm.Value")})
	if callErr != nil {
		return nil, callErr
	}
	arg__16490_104, callErr = rt.InvokeValue(rt.LookupVar("gogen", "type").Deref(), []vm.Value{vm.String("vm.Value")})
	if callErr != nil {
		return nil, callErr
	}
	arg__16491_105, callErr = rt.InvokeValue(rt.LookupVar("gogen", "variadic-param").Deref(), []vm.Value{vm.String("args"), arg__16490_104})
	if callErr != nil {
		return nil, callErr
	}
	v106, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{out_68, arg__16491_105})
	if callErr != nil {
		return nil, callErr
	}
	v109 = v106
	arity_110 = arity_66
	i_111 = i_67
	out_112 = out_68
	arg_types_113 = arg_types_69
	variadic_QMARK__114 = variadic_QMARK__70
	fixed_arity_115 = fixed_arity_71
	f_116 = f_72
	goto b13
b12:
	;
	v109 = out_75
	arity_110 = arity_73
	i_111 = i_74
	out_112 = out_75
	arg_types_113 = arg_types_76
	variadic_QMARK__114 = variadic_QMARK__77
	fixed_arity_115 = fixed_arity_78
	f_116 = f_79
	goto b13
b13:
	;
	v288 = v109
	arity_289 = arity_110
	i_290 = i_111
	out_291 = out_112
	arg_types_292 = arg_types_113
	variadic_QMARK__293 = variadic_QMARK__114
	fixed_arity_294 = fixed_arity_115
	f_295 = f_116
	goto b10
b14:
	;
	v141, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg_types_123, vm.Int(i_121)})
	if callErr != nil {
		return nil, callErr
	}
	meta_t_145 = v141
	arity_146 = arity_120
	i_147 = i_121
	out_148 = out_122
	arg_types_149 = arg_types_123
	variadic_QMARK__150 = variadic_QMARK__124
	fixed_arity_151 = fixed_arity_125
	f_152 = f_126
	load_arg_t_153 = load_arg_t_127
	v310 = v320
	v324 = v334
	v338 = v348
	goto b16
b15:
	;
	meta_t_145 = vm.Keyword("unknown")
	arity_146 = arity_128
	i_147 = i_129
	out_148 = out_130
	arg_types_149 = arg_types_131
	variadic_QMARK__150 = variadic_QMARK__132
	fixed_arity_151 = fixed_arity_133
	f_152 = f_134
	load_arg_t_153 = load_arg_t_135
	v310 = v313
	v324 = v327
	v338 = v341
	goto b16
b16:
	;
	v175, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not=").Deref(), []vm.Value{v310, load_arg_t_153})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v175) {
		meta_t_154 = meta_t_145
		arity_155 = arity_146
		i_156 = i_147
		out_157 = out_148
		arg_types_158 = arg_types_149
		variadic_QMARK__159 = variadic_QMARK__150
		fixed_arity_160 = fixed_arity_151
		f_161 = f_152
		load_arg_t_162 = load_arg_t_153
		v317 = v310
		v331 = v324
		v345 = v338
		goto b17
	} else {
		meta_t_163 = meta_t_145
		arity_164 = arity_146
		i_165 = i_147
		out_166 = out_148
		arg_types_167 = arg_types_149
		variadic_QMARK__168 = variadic_QMARK__150
		fixed_arity_169 = fixed_arity_151
		f_170 = f_152
		load_arg_t_171 = load_arg_t_153
		v315 = v310
		v329 = v324
		v343 = v338
		goto b18
	}
b17:
	;
	t_212 = load_arg_t_162
	meta_t_213 = meta_t_154
	arity_214 = arity_155
	i_215 = i_156
	out_216 = out_157
	arg_types_217 = arg_types_158
	variadic_QMARK__218 = variadic_QMARK__159
	fixed_arity_219 = fixed_arity_160
	f_220 = f_161
	load_arg_t_221 = load_arg_t_162
	v319 = v317
	v333 = v331
	v347 = v345
	goto b19
b18:
	;
	if vm.IsTruthy(v329) {
		meta_t_178 = meta_t_163
		arity_179 = arity_164
		i_180 = i_165
		out_181 = out_166
		arg_types_182 = arg_types_167
		variadic_QMARK__183 = variadic_QMARK__168
		fixed_arity_184 = fixed_arity_169
		f_185 = f_170
		load_arg_t_186 = load_arg_t_171
		v311 = v315
		v325 = v329
		v339 = v343
		goto b20
	} else {
		meta_t_187 = meta_t_163
		arity_188 = arity_164
		i_189 = i_165
		out_190 = out_166
		arg_types_191 = arg_types_167
		variadic_QMARK__192 = variadic_QMARK__168
		fixed_arity_193 = fixed_arity_169
		f_194 = f_170
		load_arg_t_195 = load_arg_t_171
		v314 = v315
		v328 = v329
		v342 = v343
		goto b21
	}
b19:
	;
	go_type_223, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "infer-go-type").Deref(), []vm.Value{t_212})
	if callErr != nil {
		return nil, callErr
	}
	v247, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{go_type_223})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v247) {
		t_224 = t_212
		meta_t_225 = meta_t_213
		arity_226 = arity_214
		i_227 = i_215
		out_228 = out_216
		arg_types_229 = arg_types_217
		variadic_QMARK__230 = variadic_QMARK__218
		fixed_arity_231 = fixed_arity_219
		f_232 = f_220
		load_arg_t_233 = load_arg_t_221
		go_type_234 = go_type_223
		v322 = v319
		v336 = v333
		v350 = v347
		goto b23
	} else {
		t_235 = t_212
		meta_t_236 = meta_t_213
		arity_237 = arity_214
		i_238 = i_215
		out_239 = out_216
		arg_types_240 = arg_types_217
		variadic_QMARK__241 = variadic_QMARK__218
		fixed_arity_242 = fixed_arity_219
		f_243 = f_220
		load_arg_t_244 = load_arg_t_221
		go_type_245 = go_type_223
		v312 = v319
		v326 = v333
		v340 = v347
		goto b24
	}
b20:
	;
	v201 = meta_t_178
	meta_t_202 = meta_t_178
	arity_203 = arity_179
	i_204 = i_180
	out_205 = out_181
	arg_types_206 = arg_types_182
	variadic_QMARK__207 = variadic_QMARK__183
	fixed_arity_208 = fixed_arity_184
	f_209 = f_185
	load_arg_t_210 = load_arg_t_186
	v318 = v311
	v332 = v325
	v346 = v339
	goto b22
b21:
	;
	v201 = vm.NIL
	meta_t_202 = meta_t_187
	arity_203 = arity_188
	i_204 = i_189
	out_205 = out_190
	arg_types_206 = arg_types_191
	variadic_QMARK__207 = variadic_QMARK__192
	fixed_arity_208 = fixed_arity_193
	f_209 = f_194
	load_arg_t_210 = load_arg_t_195
	v318 = v314
	v332 = v328
	v346 = v342
	goto b22
b22:
	;
	t_212 = v201
	meta_t_213 = meta_t_202
	arity_214 = arity_203
	i_215 = i_204
	out_216 = out_205
	arg_types_217 = arg_types_206
	variadic_QMARK__218 = variadic_QMARK__207
	fixed_arity_219 = fixed_arity_208
	f_220 = f_209
	load_arg_t_221 = load_arg_t_210
	v319 = v318
	v333 = v332
	v347 = v346
	goto b19
b23:
	;
	v275 = vm.NIL
	t_276 = t_224
	meta_t_277 = meta_t_225
	arity_278 = arity_226
	i_279 = i_227
	out_280 = out_228
	arg_types_281 = arg_types_229
	variadic_QMARK__282 = variadic_QMARK__230
	fixed_arity_283 = fixed_arity_231
	f_284 = f_232
	load_arg_t_285 = load_arg_t_233
	go_type_286 = go_type_234
	goto b25
b24:
	;
	v251 = i_238 + 1
	arg__16525_255, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(v340), vm.Int(i_238)})
	if callErr != nil {
		return nil, callErr
	}
	arg__16533_260, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(v340), vm.Int(i_238)})
	if callErr != nil {
		return nil, callErr
	}
	arg__16535_261, callErr = rt.InvokeValue(rt.LookupVar("gogen", "param").Deref(), []vm.Value{arg__16533_260, go_type_245})
	if callErr != nil {
		return nil, callErr
	}
	arg__16543_266, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(v340), vm.Int(i_238)})
	if callErr != nil {
		return nil, callErr
	}
	arg__16551_271, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(v340), vm.Int(i_238)})
	if callErr != nil {
		return nil, callErr
	}
	arg__16553_272, callErr = rt.InvokeValue(rt.LookupVar("gogen", "param").Deref(), []vm.Value{arg__16551_271, go_type_245})
	if callErr != nil {
		return nil, callErr
	}
	v273, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{out_239, arg__16553_272})
	if callErr != nil {
		return nil, callErr
	}
	i_41 = v251
	out_42 = v273
	arg_types_43 = arg_types_240
	variadic_QMARK__44 = variadic_QMARK__241
	fixed_arity_45 = fixed_arity_242
	f_46 = f_243
	v321 = v312
	v335 = v326
	v349 = v340
	goto b7
b25:
	;
	v288 = v275
	arity_289 = arity_278
	i_290 = i_279
	out_291 = out_280
	arg_types_292 = arg_types_281
	variadic_QMARK__293 = variadic_QMARK__282
	fixed_arity_294 = fixed_arity_283
	f_295 = f_284
	goto b10
}
func result_node(arg0 vm.Value) (vm.Value, error) {
	var ret_ids_2 vm.Value
	var f_3 vm.Value
	var ret_ids_4 vm.Value
	var specs_34 vm.Value
	var and__x_44 vm.Value
	var f_5 vm.Value
	var ret_ids_6 vm.Value
	var v136 vm.Value
	var f_137 vm.Value
	var ret_ids_138 vm.Value
	var f_7 vm.Value
	var and__x_8 vm.Value
	var ret_ids_9 vm.Value
	var arg__16560_15 vm.Value
	var arg__16565_18 vm.Value
	var v19 vm.Value
	var f_10 vm.Value
	var and__x_11 vm.Value
	var ret_ids_12 vm.Value
	var v22 vm.Value
	var f_23 vm.Value
	var and__x_24 vm.Value
	var ret_ids_25 vm.Value
	var f_35 vm.Value
	var ret_ids_36 vm.Value
	var specs_37 vm.Value
	var arg__16617_72 vm.Value
	var arg__16622_75 vm.Value
	var arg__16623_76 vm.Value
	var arg__16628_79 vm.Value
	var arg__16633_82 vm.Value
	var arg__16634_83 vm.Value
	var arg__16635_84 vm.Value
	var results_85 vm.Value
	var v95 vm.Value
	var f_38 vm.Value
	var ret_ids_39 vm.Value
	var specs_40 vm.Value
	var v129 vm.Value
	var f_130 vm.Value
	var ret_ids_131 vm.Value
	var specs_132 vm.Value
	var f_45 vm.Value
	var ret_ids_46 vm.Value
	var specs_47 vm.Value
	var and__x_48 vm.Value
	var arg__16606_56 vm.Value
	var arg__16611_59 vm.Value
	var arg__16612_60 vm.Value
	var v61 bool
	var f_49 vm.Value
	var ret_ids_50 vm.Value
	var specs_51 vm.Value
	var and__x_52 vm.Value
	var v64 vm.Value
	var f_65 vm.Value
	var ret_ids_66 vm.Value
	var specs_67 vm.Value
	var and__x_68 vm.Value
	var f_86 vm.Value
	var ret_ids_87 vm.Value
	var specs_88 vm.Value
	var results_89 vm.Value
	var arg__16643_100 vm.Value
	var arg__16648_105 vm.Value
	var arg__16649_106 vm.Value
	var arg__16655_111 vm.Value
	var arg__16660_116 vm.Value
	var arg__16661_117 vm.Value
	var v118 vm.Value
	var f_90 vm.Value
	var ret_ids_91 vm.Value
	var specs_92 vm.Value
	var results_93 vm.Value
	var v121 vm.Value
	var f_122 vm.Value
	var ret_ids_123 vm.Value
	var specs_124 vm.Value
	var results_125 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = ret_ids_2, f_3, ret_ids_4, specs_34, and__x_44, f_5, ret_ids_6, v136, f_137, ret_ids_138, f_7, and__x_8, ret_ids_9, arg__16560_15, arg__16565_18, v19, f_10, and__x_11, ret_ids_12, v22, f_23, and__x_24, ret_ids_25, f_35, ret_ids_36, specs_37, arg__16617_72, arg__16622_75, arg__16623_76, arg__16628_79, arg__16633_82, arg__16634_83, arg__16635_84, results_85, v95, f_38, ret_ids_39, specs_40, v129, f_130, ret_ids_131, specs_132, f_45, ret_ids_46, specs_47, and__x_48, arg__16606_56, arg__16611_59, arg__16612_60, v61, f_49, ret_ids_50, specs_51, and__x_52, v64, f_65, ret_ids_66, specs_67, and__x_68, f_86, ret_ids_87, specs_88, results_89, arg__16643_100, arg__16648_105, arg__16649_106, arg__16655_111, arg__16660_116, arg__16661_117, v118, f_90, ret_ids_91, specs_92, results_93, v121, f_122, ret_ids_123, specs_124, results_125
	ret_ids_2, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "return-ref-ids").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(ret_ids_2) {
		f_7 = arg0
		and__x_8 = ret_ids_2
		ret_ids_9 = ret_ids_2
		goto b4
	} else {
		f_10 = arg0
		and__x_11 = ret_ids_2
		ret_ids_12 = ret_ids_2
		goto b5
	}
b1:
	;
	specs_34, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var arg__16587_3 vm.Value
		var arg__16594_6 vm.Value
		var v7 vm.Value
		var callErr error
		_, _, _ = arg__16587_3, arg__16594_6, v7
		arg__16587_3, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_3})
		if callErr != nil {
			return nil, callErr
		}
		arg__16594_6, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_3})
		if callErr != nil {
			return nil, callErr
		}
		v7, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "go-type-spec").Deref(), []vm.Value{arg__16594_6})
		if callErr != nil {
			return nil, callErr
		}
		return v7, nil
	}), ret_ids_4})
	if callErr != nil {
		return nil, callErr
	}
	and__x_44, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "every?").Deref(), []vm.Value{rt.LookupVar("clojure.core", "identity").Deref(), specs_34})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(and__x_44) {
		f_45 = f_3
		ret_ids_46 = ret_ids_4
		specs_47 = specs_34
		and__x_48 = and__x_44
		goto b10
	} else {
		f_49 = f_3
		ret_ids_50 = ret_ids_4
		specs_51 = specs_34
		and__x_52 = and__x_44
		goto b11
	}
b2:
	;
	v136 = vm.NIL
	f_137 = f_5
	ret_ids_138 = ret_ids_6
	goto b3
b3:
	;
	return v136, nil
b4:
	;
	arg__16560_15, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{ret_ids_9})
	if callErr != nil {
		return nil, callErr
	}
	arg__16565_18, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{ret_ids_9})
	if callErr != nil {
		return nil, callErr
	}
	v19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "pos?").Deref(), []vm.Value{arg__16565_18})
	if callErr != nil {
		return nil, callErr
	}
	v22 = v19
	f_23 = f_7
	and__x_24 = and__x_8
	ret_ids_25 = ret_ids_9
	goto b6
b5:
	;
	v22 = and__x_11
	f_23 = f_10
	and__x_24 = and__x_11
	ret_ids_25 = ret_ids_12
	goto b6
b6:
	;
	if vm.IsTruthy(v22) {
		f_3 = f_23
		ret_ids_4 = ret_ids_25
		goto b1
	} else {
		f_5 = f_23
		ret_ids_6 = ret_ids_25
		goto b2
	}
b7:
	;
	arg__16617_72, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{specs_37})
	if callErr != nil {
		return nil, callErr
	}
	arg__16622_75, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{specs_37})
	if callErr != nil {
		return nil, callErr
	}
	arg__16623_76, callErr = rt.InvokeValue(rt.LookupVar("gogen", "type").Deref(), []vm.Value{arg__16622_75})
	if callErr != nil {
		return nil, callErr
	}
	arg__16628_79, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{specs_37})
	if callErr != nil {
		return nil, callErr
	}
	arg__16633_82, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{specs_37})
	if callErr != nil {
		return nil, callErr
	}
	arg__16634_83, callErr = rt.InvokeValue(rt.LookupVar("gogen", "type").Deref(), []vm.Value{arg__16633_82})
	if callErr != nil {
		return nil, callErr
	}
	arg__16635_84, callErr = rt.InvokeValue(rt.LookupVar("gogen", "result").Deref(), []vm.Value{arg__16634_83})
	if callErr != nil {
		return nil, callErr
	}
	results_85, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__16635_84})
	if callErr != nil {
		return nil, callErr
	}
	v95, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "function-needs-error?").Deref(), []vm.Value{f_35})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v95) {
		f_86 = f_35
		ret_ids_87 = ret_ids_36
		specs_88 = specs_37
		results_89 = results_85
		goto b13
	} else {
		f_90 = f_35
		ret_ids_91 = ret_ids_36
		specs_92 = specs_37
		results_93 = results_85
		goto b14
	}
b8:
	;
	v129 = vm.NIL
	f_130 = f_38
	ret_ids_131 = ret_ids_39
	specs_132 = specs_40
	goto b9
b9:
	;
	v136 = v129
	f_137 = f_130
	ret_ids_138 = ret_ids_131
	goto b3
b10:
	;
	arg__16606_56, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "distinct").Deref(), []vm.Value{specs_47})
	if callErr != nil {
		return nil, callErr
	}
	arg__16611_59, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "distinct").Deref(), []vm.Value{specs_47})
	if callErr != nil {
		return nil, callErr
	}
	arg__16612_60, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{arg__16611_59})
	if callErr != nil {
		return nil, callErr
	}
	v61 = arg__16612_60 == vm.Int(1)
	v64 = vm.Boolean(v61)
	f_65 = f_45
	ret_ids_66 = ret_ids_46
	specs_67 = specs_47
	and__x_68 = and__x_48
	goto b12
b11:
	;
	v64 = and__x_52
	f_65 = f_49
	ret_ids_66 = ret_ids_50
	specs_67 = specs_51
	and__x_68 = and__x_52
	goto b12
b12:
	;
	if vm.IsTruthy(v64) {
		f_35 = f_65
		ret_ids_36 = ret_ids_66
		specs_37 = specs_67
		goto b7
	} else {
		f_38 = f_65
		ret_ids_39 = ret_ids_66
		specs_40 = specs_67
		goto b8
	}
b13:
	;
	arg__16643_100, callErr = rt.InvokeValue(rt.LookupVar("gogen", "type").Deref(), []vm.Value{vm.String("error")})
	if callErr != nil {
		return nil, callErr
	}
	arg__16648_105, callErr = rt.InvokeValue(rt.LookupVar("gogen", "type").Deref(), []vm.Value{vm.String("error")})
	if callErr != nil {
		return nil, callErr
	}
	arg__16649_106, callErr = rt.InvokeValue(rt.LookupVar("gogen", "result").Deref(), []vm.Value{arg__16648_105})
	if callErr != nil {
		return nil, callErr
	}
	arg__16655_111, callErr = rt.InvokeValue(rt.LookupVar("gogen", "type").Deref(), []vm.Value{vm.String("error")})
	if callErr != nil {
		return nil, callErr
	}
	arg__16660_116, callErr = rt.InvokeValue(rt.LookupVar("gogen", "type").Deref(), []vm.Value{vm.String("error")})
	if callErr != nil {
		return nil, callErr
	}
	arg__16661_117, callErr = rt.InvokeValue(rt.LookupVar("gogen", "result").Deref(), []vm.Value{arg__16660_116})
	if callErr != nil {
		return nil, callErr
	}
	v118, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{results_89, arg__16661_117})
	if callErr != nil {
		return nil, callErr
	}
	v121 = v118
	f_122 = f_86
	ret_ids_123 = ret_ids_87
	specs_124 = specs_88
	results_125 = results_89
	goto b15
b14:
	;
	v121 = results_93
	f_122 = f_90
	ret_ids_123 = ret_ids_91
	specs_124 = specs_92
	results_125 = results_93
	goto b15
b15:
	;
	v129 = v121
	f_130 = f_122
	ret_ids_131 = ret_ids_123
	specs_132 = specs_124
	goto b9
}
func return_ref_ids(arg0 vm.Value) (vm.Value, error) {
	var v5 vm.Value
	var blocks_1 vm.Value
	var out_2 vm.Value
	var f_3 vm.Value
	var v15 vm.Value
	var blocks_8 vm.Value
	var out_9 vm.Value
	var f_10 vm.Value
	var blocks_11 vm.Value
	var out_12 vm.Value
	var f_13 vm.Value
	var arg__16671_19 vm.Value
	var arg__16677_22 vm.Value
	var term_23 vm.Value
	var v33 vm.Value
	var v123 vm.Value
	var blocks_124 vm.Value
	var out_125 vm.Value
	var f_126 vm.Value
	var blocks_24 vm.Value
	var out_25 vm.Value
	var f_26 vm.Value
	var term_27 vm.Value
	var v36 vm.Value
	var blocks_28 vm.Value
	var out_29 vm.Value
	var f_30 vm.Value
	var term_31 vm.Value
	var arg__16690_47 vm.Value
	var v49 bool
	var v117 vm.Value
	var blocks_118 vm.Value
	var out_119 vm.Value
	var f_120 vm.Value
	var term_121 vm.Value
	var blocks_38 vm.Value
	var out_39 vm.Value
	var f_40 vm.Value
	var term_41 vm.Value
	var refs_52 vm.Value
	var arg__16701_65 vm.Value
	var v66 bool
	var blocks_42 vm.Value
	var out_43 vm.Value
	var f_44 vm.Value
	var term_45 vm.Value
	var v111 vm.Value
	var blocks_112 vm.Value
	var out_113 vm.Value
	var f_114 vm.Value
	var term_115 vm.Value
	var blocks_53 vm.Value
	var out_54 vm.Value
	var f_55 vm.Value
	var term_56 vm.Value
	var refs_57 vm.Value
	var v69 vm.Value
	var arg__16711_73 vm.Value
	var arg__16719_78 vm.Value
	var v79 vm.Value
	var blocks_58 vm.Value
	var out_59 vm.Value
	var f_60 vm.Value
	var term_61 vm.Value
	var refs_62 vm.Value
	var v83 vm.Value
	var blocks_84 vm.Value
	var out_85 vm.Value
	var f_86 vm.Value
	var term_87 vm.Value
	var refs_88 vm.Value
	var blocks_90 vm.Value
	var out_91 vm.Value
	var f_92 vm.Value
	var term_93 vm.Value
	var v101 vm.Value
	var blocks_94 vm.Value
	var out_95 vm.Value
	var f_96 vm.Value
	var term_97 vm.Value
	var v105 vm.Value
	var blocks_106 vm.Value
	var out_107 vm.Value
	var f_108 vm.Value
	var term_109 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = v5, blocks_1, out_2, f_3, v15, blocks_8, out_9, f_10, blocks_11, out_12, f_13, arg__16671_19, arg__16677_22, term_23, v33, v123, blocks_124, out_125, f_126, blocks_24, out_25, f_26, term_27, v36, blocks_28, out_29, f_30, term_31, arg__16690_47, v49, v117, blocks_118, out_119, f_120, term_121, blocks_38, out_39, f_40, term_41, refs_52, arg__16701_65, v66, blocks_42, out_43, f_44, term_45, v111, blocks_112, out_113, f_114, term_115, blocks_53, out_54, f_55, term_56, refs_57, v69, arg__16711_73, arg__16719_78, v79, blocks_58, out_59, f_60, term_61, refs_62, v83, blocks_84, out_85, f_86, term_87, refs_88, blocks_90, out_91, f_92, term_93, v101, blocks_94, out_95, f_96, term_97, v105, blocks_106, out_107, f_108, term_109
	v5, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	blocks_1 = v5
	out_2 = vm.NewArrayVector([]vm.Value{})
	f_3 = arg0
	goto b1
b1:
	;
	v15, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "empty?").Deref(), []vm.Value{blocks_1})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v15) {
		blocks_8 = blocks_1
		out_9 = out_2
		f_10 = f_3
		goto b2
	} else {
		blocks_11 = blocks_1
		out_12 = out_2
		f_13 = f_3
		goto b3
	}
b2:
	;
	v123 = out_9
	blocks_124 = blocks_8
	out_125 = out_9
	f_126 = f_10
	goto b4
b3:
	;
	arg__16671_19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{blocks_11})
	if callErr != nil {
		return nil, callErr
	}
	arg__16677_22, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{blocks_11})
	if callErr != nil {
		return nil, callErr
	}
	term_23, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-term").Deref(), []vm.Value{arg__16677_22, f_13})
	if callErr != nil {
		return nil, callErr
	}
	v33, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{term_23})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v33) {
		blocks_24 = blocks_11
		out_25 = out_12
		f_26 = f_13
		term_27 = term_23
		goto b5
	} else {
		blocks_28 = blocks_11
		out_29 = out_12
		f_30 = f_13
		term_31 = term_23
		goto b6
	}
b4:
	;
	return v123, nil
b5:
	;
	v36, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{blocks_24})
	if callErr != nil {
		return nil, callErr
	}
	blocks_1 = v36
	out_2 = out_25
	f_3 = f_26
	goto b1
b6:
	;
	arg__16690_47, callErr = rt.InvokeValue(rt.LookupVar("ir", "op").Deref(), []vm.Value{term_31, f_30})
	if callErr != nil {
		return nil, callErr
	}
	v49 = arg__16690_47 == vm.Keyword("return")
	if v49 {
		blocks_38 = blocks_28
		out_39 = out_29
		f_40 = f_30
		term_41 = term_31
		goto b8
	} else {
		blocks_42 = blocks_28
		out_43 = out_29
		f_44 = f_30
		term_45 = term_31
		goto b9
	}
b7:
	;
	v123 = v117
	blocks_124 = blocks_118
	out_125 = out_119
	f_126 = f_120
	goto b4
b8:
	;
	refs_52, callErr = rt.InvokeValue(rt.LookupVar("ir", "refs").Deref(), []vm.Value{term_41, f_40})
	if callErr != nil {
		return nil, callErr
	}
	arg__16701_65, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{refs_52})
	if callErr != nil {
		return nil, callErr
	}
	v66 = arg__16701_65 == vm.Int(1)
	if v66 {
		blocks_53 = blocks_38
		out_54 = out_39
		f_55 = f_40
		term_56 = term_41
		refs_57 = refs_52
		goto b11
	} else {
		blocks_58 = blocks_38
		out_59 = out_39
		f_60 = f_40
		term_61 = term_41
		refs_62 = refs_52
		goto b12
	}
b9:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		blocks_90 = blocks_42
		out_91 = out_43
		f_92 = f_44
		term_93 = term_45
		goto b14
	} else {
		blocks_94 = blocks_42
		out_95 = out_43
		f_96 = f_44
		term_97 = term_45
		goto b15
	}
b10:
	;
	v117 = v111
	blocks_118 = blocks_112
	out_119 = out_113
	f_120 = f_114
	term_121 = term_115
	goto b7
b11:
	;
	v69, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{blocks_53})
	if callErr != nil {
		return nil, callErr
	}
	arg__16711_73, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{refs_57, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	arg__16719_78, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{refs_57, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	v79, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{out_54, arg__16719_78})
	if callErr != nil {
		return nil, callErr
	}
	blocks_1 = v69
	out_2 = v79
	f_3 = f_55
	goto b1
b12:
	;
	v83 = vm.NIL
	blocks_84 = blocks_58
	out_85 = out_59
	f_86 = f_60
	term_87 = term_61
	refs_88 = refs_62
	goto b13
b13:
	;
	v111 = v83
	blocks_112 = blocks_84
	out_113 = out_85
	f_114 = f_86
	term_115 = term_87
	goto b10
b14:
	;
	v101, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{blocks_90})
	if callErr != nil {
		return nil, callErr
	}
	blocks_1 = v101
	out_2 = out_91
	f_3 = f_92
	goto b1
b15:
	;
	v105 = vm.NIL
	blocks_106 = blocks_94
	out_107 = out_95
	f_108 = f_96
	term_109 = term_97
	goto b16
b16:
	;
	v111 = v105
	blocks_112 = blocks_106
	out_113 = out_107
	f_114 = f_108
	term_115 = term_109
	goto b10
}
func source_name_of(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var sis_4 vm.Value
	var remaining_5 vm.Value
	var v16 vm.Value
	var f_7 vm.Value
	var nid_8 vm.Value
	var sis_9 vm.Value
	var remaining_10 vm.Value
	var arg__16734_19 vm.Value
	var arg__16739_22 vm.Value
	var sym_23 vm.Value
	var and__x_35 vm.Value
	var f_11 vm.Value
	var nid_12 vm.Value
	var sis_13 vm.Value
	var remaining_14 vm.Value
	var v78 vm.Value
	var f_79 vm.Value
	var nid_80 vm.Value
	var sis_81 vm.Value
	var remaining_82 vm.Value
	var f_24 vm.Value
	var nid_25 vm.Value
	var sis_26 vm.Value
	var remaining_27 vm.Value
	var sym_28 vm.Value
	var f_29 vm.Value
	var nid_30 vm.Value
	var sis_31 vm.Value
	var remaining_32 vm.Value
	var sym_33 vm.Value
	var v67 vm.Value
	var v69 vm.Value
	var f_70 vm.Value
	var nid_71 vm.Value
	var sis_72 vm.Value
	var remaining_73 vm.Value
	var sym_74 vm.Value
	var f_36 vm.Value
	var nid_37 vm.Value
	var sis_38 vm.Value
	var remaining_39 vm.Value
	var sym_40 vm.Value
	var and__x_41 vm.Value
	var arg__16746_50 vm.Value
	var arg__16751_53 vm.Value
	var v54 vm.Value
	var f_42 vm.Value
	var nid_43 vm.Value
	var sis_44 vm.Value
	var remaining_45 vm.Value
	var sym_46 vm.Value
	var and__x_47 vm.Value
	var v57 vm.Value
	var f_58 vm.Value
	var nid_59 vm.Value
	var sis_60 vm.Value
	var remaining_61 vm.Value
	var sym_62 vm.Value
	var and__x_63 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = sis_4, remaining_5, v16, f_7, nid_8, sis_9, remaining_10, arg__16734_19, arg__16739_22, sym_23, and__x_35, f_11, nid_12, sis_13, remaining_14, v78, f_79, nid_80, sis_81, remaining_82, f_24, nid_25, sis_26, remaining_27, sym_28, f_29, nid_30, sis_31, remaining_32, sym_33, v67, v69, f_70, nid_71, sis_72, remaining_73, sym_74, f_36, nid_37, sis_38, remaining_39, sym_40, and__x_41, arg__16746_50, arg__16751_53, v54, f_42, nid_43, sis_44, remaining_45, sym_46, and__x_47, v57, f_58, nid_59, sis_60, remaining_61, sym_62, and__x_63
	sis_4, callErr = rt.InvokeValue(rt.LookupVar("ir", "source-infos").Deref(), []vm.Value{arg1, arg0})
	if callErr != nil {
		return nil, callErr
	}
	remaining_5 = sis_4
	goto b1
b1:
	;
	v16, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{remaining_5})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v16) {
		f_7 = arg0
		nid_8 = arg1
		sis_9 = sis_4
		remaining_10 = remaining_5
		goto b2
	} else {
		f_11 = arg0
		nid_12 = arg1
		sis_13 = sis_4
		remaining_14 = remaining_5
		goto b3
	}
b2:
	;
	arg__16734_19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{remaining_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__16739_22, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{remaining_10})
	if callErr != nil {
		return nil, callErr
	}
	sym_23, callErr = rt.InvokeValue(rt.LookupVar("ir", "source-info-symbol").Deref(), []vm.Value{arg__16739_22})
	if callErr != nil {
		return nil, callErr
	}
	and__x_35, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "string?").Deref(), []vm.Value{sym_23})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(and__x_35) {
		f_36 = f_7
		nid_37 = nid_8
		sis_38 = sis_9
		remaining_39 = remaining_10
		sym_40 = sym_23
		and__x_41 = and__x_35
		goto b8
	} else {
		f_42 = f_7
		nid_43 = nid_8
		sis_44 = sis_9
		remaining_45 = remaining_10
		sym_46 = sym_23
		and__x_47 = and__x_35
		goto b9
	}
b3:
	;
	v78 = vm.NIL
	f_79 = f_11
	nid_80 = nid_12
	sis_81 = sis_13
	remaining_82 = remaining_14
	goto b4
b4:
	;
	return v78, nil
b5:
	;
	v69 = sym_28
	f_70 = f_24
	nid_71 = nid_25
	sis_72 = sis_26
	remaining_73 = remaining_27
	sym_74 = sym_28
	goto b7
b6:
	;
	v67, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "next").Deref(), []vm.Value{remaining_32})
	if callErr != nil {
		return nil, callErr
	}
	remaining_5 = v67
	goto b1
b7:
	;
	v78 = v69
	f_79 = f_70
	nid_80 = nid_71
	sis_81 = sis_72
	remaining_82 = remaining_73
	goto b4
b8:
	;
	arg__16746_50, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{sym_40})
	if callErr != nil {
		return nil, callErr
	}
	arg__16751_53, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{sym_40})
	if callErr != nil {
		return nil, callErr
	}
	v54, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "pos?").Deref(), []vm.Value{arg__16751_53})
	if callErr != nil {
		return nil, callErr
	}
	v57 = v54
	f_58 = f_36
	nid_59 = nid_37
	sis_60 = sis_38
	remaining_61 = remaining_39
	sym_62 = sym_40
	and__x_63 = and__x_41
	goto b10
b9:
	;
	v57 = and__x_47
	f_58 = f_42
	nid_59 = nid_43
	sis_60 = sis_44
	remaining_61 = remaining_45
	sym_62 = sym_46
	and__x_63 = and__x_47
	goto b10
b10:
	;
	if vm.IsTruthy(v57) {
		f_24 = f_58
		nid_25 = nid_59
		sis_26 = sis_60
		remaining_27 = remaining_61
		sym_28 = sym_62
		goto b5
	} else {
		f_29 = f_58
		nid_30 = nid_59
		sis_31 = sis_60
		remaining_32 = remaining_61
		sym_33 = sym_62
		goto b6
	}
}
func template_fns(arg0 vm.Value) (vm.Value, error) {
	var v4 vm.Value
	var aux_1 vm.Value
	var arg__16761_8 vm.Value
	var v9 vm.Value
	var aux_2 vm.Value
	var v14 vm.Value
	var v39 vm.Value
	var aux_40 vm.Value
	var aux_11 vm.Value
	var arg__16768_18 vm.Value
	var arg__16773_22 vm.Value
	var v23 vm.Value
	var aux_12 vm.Value
	var v36 vm.Value
	var aux_37 vm.Value
	var aux_25 vm.Value
	var aux_26 vm.Value
	var v33 vm.Value
	var aux_34 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = v4, aux_1, arg__16761_8, v9, aux_2, v14, v39, aux_40, aux_11, arg__16768_18, arg__16773_22, v23, aux_12, v36, aux_37, aux_25, aux_26, v33, aux_34
	v4, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "fn-template?").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v4) {
		aux_1 = arg0
		goto b1
	} else {
		aux_2 = arg0
		goto b2
	}
b1:
	;
	arg__16761_8, callErr = rt.InvokeValue(vm.Keyword("fn"), []vm.Value{aux_1})
	if callErr != nil {
		return nil, callErr
	}
	v9, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__16761_8})
	if callErr != nil {
		return nil, callErr
	}
	v39 = v9
	aux_40 = aux_1
	goto b3
b2:
	;
	v14, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "multi-fn-template?").Deref(), []vm.Value{aux_2})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v14) {
		aux_11 = aux_2
		goto b4
	} else {
		aux_12 = aux_2
		goto b5
	}
b3:
	;
	return v39, nil
b4:
	;
	arg__16768_18, callErr = rt.InvokeValue(vm.Keyword("fns"), []vm.Value{aux_11})
	if callErr != nil {
		return nil, callErr
	}
	arg__16773_22, callErr = rt.InvokeValue(vm.Keyword("fns"), []vm.Value{aux_11})
	if callErr != nil {
		return nil, callErr
	}
	v23, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapv").Deref(), []vm.Value{vm.Keyword("fn"), arg__16773_22})
	if callErr != nil {
		return nil, callErr
	}
	v36 = v23
	aux_37 = aux_11
	goto b6
b5:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		aux_25 = aux_12
		goto b7
	} else {
		aux_26 = aux_12
		goto b8
	}
b6:
	;
	v39 = v36
	aux_40 = aux_37
	goto b3
b7:
	;
	v33 = vm.NewArrayVector([]vm.Value{})
	aux_34 = aux_25
	goto b9
b8:
	;
	v33 = vm.NIL
	aux_34 = aux_26
	goto b9
b9:
	;
	v36 = v33
	aux_37 = aux_34
	goto b6
}
func transfer_stmts(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var assigns_4 vm.Value
	var f_5 vm.Value
	var closed_exprs_6 vm.Value
	var bt_7 vm.Value
	var assigns_8 vm.Value
	var arg__16785_15 vm.Value
	var arg__16790_18 vm.Value
	var arg__16791_19 vm.Value
	var arg__16796_22 vm.Value
	var arg__16801_25 vm.Value
	var arg__16802_26 vm.Value
	var arg__16803_27 vm.Value
	var arg__16809_30 vm.Value
	var arg__16814_33 vm.Value
	var arg__16815_34 vm.Value
	var arg__16820_37 vm.Value
	var arg__16825_40 vm.Value
	var arg__16826_41 vm.Value
	var arg__16827_42 vm.Value
	var v43 vm.Value
	var f_9 vm.Value
	var closed_exprs_10 vm.Value
	var bt_11 vm.Value
	var assigns_12 vm.Value
	var v47 vm.Value
	var f_48 vm.Value
	var closed_exprs_49 vm.Value
	var bt_50 vm.Value
	var assigns_51 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = assigns_4, f_5, closed_exprs_6, bt_7, assigns_8, arg__16785_15, arg__16790_18, arg__16791_19, arg__16796_22, arg__16801_25, arg__16802_26, arg__16803_27, arg__16809_30, arg__16814_33, arg__16815_34, arg__16820_37, arg__16825_40, arg__16826_41, arg__16827_42, v43, f_9, closed_exprs_10, bt_11, assigns_12, v47, f_48, closed_exprs_49, bt_50, assigns_51
	assigns_4, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "emit-assignments-for-target").Deref(), []vm.Value{arg0, arg1, arg2})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(assigns_4) {
		f_5 = arg0
		closed_exprs_6 = arg1
		bt_7 = arg2
		assigns_8 = assigns_4
		goto b1
	} else {
		f_9 = arg0
		closed_exprs_10 = arg1
		bt_11 = arg2
		assigns_12 = assigns_4
		goto b2
	}
b1:
	;
	arg__16785_15, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-target").Deref(), []vm.Value{bt_7})
	if callErr != nil {
		return nil, callErr
	}
	arg__16790_18, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-target").Deref(), []vm.Value{bt_7})
	if callErr != nil {
		return nil, callErr
	}
	arg__16791_19, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "label-name").Deref(), []vm.Value{arg__16790_18})
	if callErr != nil {
		return nil, callErr
	}
	arg__16796_22, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-target").Deref(), []vm.Value{bt_7})
	if callErr != nil {
		return nil, callErr
	}
	arg__16801_25, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-target").Deref(), []vm.Value{bt_7})
	if callErr != nil {
		return nil, callErr
	}
	arg__16802_26, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "label-name").Deref(), []vm.Value{arg__16801_25})
	if callErr != nil {
		return nil, callErr
	}
	arg__16803_27, callErr = rt.InvokeValue(rt.LookupVar("gogen", "goto-stmt").Deref(), []vm.Value{arg__16802_26})
	if callErr != nil {
		return nil, callErr
	}
	arg__16809_30, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-target").Deref(), []vm.Value{bt_7})
	if callErr != nil {
		return nil, callErr
	}
	arg__16814_33, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-target").Deref(), []vm.Value{bt_7})
	if callErr != nil {
		return nil, callErr
	}
	arg__16815_34, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "label-name").Deref(), []vm.Value{arg__16814_33})
	if callErr != nil {
		return nil, callErr
	}
	arg__16820_37, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-target").Deref(), []vm.Value{bt_7})
	if callErr != nil {
		return nil, callErr
	}
	arg__16825_40, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-target").Deref(), []vm.Value{bt_7})
	if callErr != nil {
		return nil, callErr
	}
	arg__16826_41, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "label-name").Deref(), []vm.Value{arg__16825_40})
	if callErr != nil {
		return nil, callErr
	}
	arg__16827_42, callErr = rt.InvokeValue(rt.LookupVar("gogen", "goto-stmt").Deref(), []vm.Value{arg__16826_41})
	if callErr != nil {
		return nil, callErr
	}
	v43, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{assigns_8, arg__16827_42})
	if callErr != nil {
		return nil, callErr
	}
	v47 = v43
	f_48 = f_5
	closed_exprs_49 = closed_exprs_6
	bt_50 = bt_7
	assigns_51 = assigns_8
	goto b3
b2:
	;
	v47 = vm.NIL
	f_48 = f_9
	closed_exprs_49 = closed_exprs_10
	bt_50 = bt_11
	assigns_51 = assigns_12
	goto b3
b3:
	;
	return v47, nil
}
func unsupported(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var v7 bool
	var mode_2 vm.Value
	var msg_3 vm.Value
	var v15 vm.Value
	var mode_4 vm.Value
	var msg_5 vm.Value
	var arg__16842_20 vm.Value
	var arg__16849_25 vm.Value
	var v26 vm.Value
	var v28 vm.Value
	var mode_29 vm.Value
	var msg_30 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _ = v7, mode_2, msg_3, v15, mode_4, msg_5, arg__16842_20, arg__16849_25, v26, v28, mode_29, msg_30
	v7 = arg0 == vm.Keyword("bridge")
	if v7 {
		mode_2 = arg0
		msg_3 = arg1
		goto b1
	} else {
		mode_4 = arg0
		msg_5 = arg1
		goto b2
	}
b1:
	;
	v15, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("status"), vm.Keyword("fallback"), vm.Keyword("decl"), vm.NIL, vm.Keyword("reason"), msg_3})
	if callErr != nil {
		return nil, callErr
	}
	v28 = v15
	mode_29 = mode_2
	msg_30 = msg_3
	goto b3
b2:
	;
	arg__16842_20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("ir.lower-go: "), msg_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__16849_25, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("ir.lower-go: "), msg_5})
	if callErr != nil {
		return nil, callErr
	}
	v26, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "throw").Deref(), []vm.Value{arg__16849_25})
	if callErr != nil {
		return nil, callErr
	}
	v28 = v26
	mode_29 = mode_4
	msg_30 = msg_5
	goto b3
b3:
	;
	return v28, nil
}
func used_QMARK_(arg0 vm.Value, arg1 int) (vm.Value, error) {
	var uses_3 vm.Value
	var arg__16857_5 vm.Value
	var and__x_6 bool
	var f_7 vm.Value
	var nid_8 int
	var uses_9 vm.Value
	var and__x_10 bool
	var arg__16863_17 vm.Value
	var arg__16870_20 vm.Value
	var arg__16871_21 vm.Value
	var arg__16878_24 vm.Value
	var arg__16885_27 vm.Value
	var arg__16886_28 vm.Value
	var v29 vm.Value
	var f_11 vm.Value
	var nid_12 int
	var uses_13 vm.Value
	var and__x_14 bool
	var v32 vm.Value
	var f_33 vm.Value
	var nid_34 int
	var uses_35 vm.Value
	var and__x_36 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = uses_3, arg__16857_5, and__x_6, f_7, nid_8, uses_9, and__x_10, arg__16863_17, arg__16870_20, arg__16871_21, arg__16878_24, arg__16885_27, arg__16886_28, v29, f_11, nid_12, uses_13, and__x_14, v32, f_33, nid_34, uses_35, and__x_36
	uses_3, callErr = rt.InvokeValue(rt.LookupVar("ir", "uses").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__16857_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{uses_3})
	if callErr != nil {
		return nil, callErr
	}
	and__x_6 = rt.LtValue(vm.Int(arg1), arg__16857_5)
	if and__x_6 {
		f_7 = arg0
		nid_8 = arg1
		uses_9 = uses_3
		and__x_10 = and__x_6
		goto b1
	} else {
		f_11 = arg0
		nid_12 = arg1
		uses_13 = uses_3
		and__x_14 = and__x_6
		goto b2
	}
b1:
	;
	arg__16863_17, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{uses_9, vm.Int(nid_8)})
	if callErr != nil {
		return nil, callErr
	}
	arg__16870_20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{uses_9, vm.Int(nid_8)})
	if callErr != nil {
		return nil, callErr
	}
	arg__16871_21, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{arg__16870_20})
	if callErr != nil {
		return nil, callErr
	}
	arg__16878_24, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{uses_9, vm.Int(nid_8)})
	if callErr != nil {
		return nil, callErr
	}
	arg__16885_27, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{uses_9, vm.Int(nid_8)})
	if callErr != nil {
		return nil, callErr
	}
	arg__16886_28, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{arg__16885_27})
	if callErr != nil {
		return nil, callErr
	}
	v29, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "pos?").Deref(), []vm.Value{arg__16886_28})
	if callErr != nil {
		return nil, callErr
	}
	v32 = v29
	f_33 = f_7
	nid_34 = nid_8
	uses_35 = uses_9
	and__x_36 = vm.Boolean(and__x_10)
	goto b3
b2:
	;
	v32 = vm.Boolean(and__x_14)
	f_33 = f_11
	nid_34 = nid_12
	uses_35 = uses_13
	and__x_36 = vm.Boolean(and__x_14)
	goto b3
b3:
	;
	return v32, nil
}
func vm_call(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__17278_3 vm.Value
	var arg__17284_6 vm.Value
	var v7 vm.Value
	var callErr error
	_, _, _ = arg__17278_3, arg__17284_6, v7
	arg__17278_3, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "vm-sel").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__17284_6, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "vm-sel").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v7, callErr = rt.InvokeValue(rt.LookupVar("gogen", "call").Deref(), []vm.Value{arg__17284_6, arg1})
	if callErr != nil {
		return nil, callErr
	}
	return v7, nil
}
func vm_sel(arg0 vm.Value) (vm.Value, error) {
	var arg__17289_4 vm.Value
	var arg__17295_9 vm.Value
	var v10 vm.Value
	var callErr error
	_, _, _ = arg__17289_4, arg__17295_9, v10
	arg__17289_4, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("vm")})
	if callErr != nil {
		return nil, callErr
	}
	arg__17295_9, callErr = rt.InvokeValue(rt.LookupVar("gogen", "ident").Deref(), []vm.Value{vm.String("vm")})
	if callErr != nil {
		return nil, callErr
	}
	v10, callErr = rt.InvokeValue(rt.LookupVar("gogen", "field-sel").Deref(), []vm.Value{arg__17295_9, arg0})
	if callErr != nil {
		return nil, callErr
	}
	return v10, nil
}
