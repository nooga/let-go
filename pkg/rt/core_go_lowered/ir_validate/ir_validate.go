package ir_validate

import (
	rt "github.com/nooga/let-go/pkg/rt"
	vm "github.com/nooga/let-go/pkg/vm"
)

func check_no_cross_block_refs_BANG_(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__28334_5 vm.Value
	var insts_6 vm.Value
	var arg__28340_10 vm.Value
	var arg__28347_15 vm.Value
	var doseq_seq__28325_16 vm.Value
	var doseq_loop__28326_17 vm.Value
	var label_18 vm.Value
	var insts_19 vm.Value
	var f_21 vm.Value
	var doseq_seq__28325_22 vm.Value
	var doseq_loop__28326_23 vm.Value
	var label_24 vm.Value
	var insts_25 vm.Value
	var vec__28329_33 vm.Value
	var i_39 vm.Value
	var ins_45 vm.Value
	var op_49 vm.Value
	var v71 vm.Value
	var f_26 vm.Value
	var doseq_seq__28325_27 vm.Value
	var doseq_loop__28326_28 vm.Value
	var label_29 vm.Value
	var insts_30 vm.Value
	var v289 vm.Value
	var f_290 vm.Value
	var doseq_seq__28325_291 vm.Value
	var doseq_loop__28326_292 vm.Value
	var label_293 vm.Value
	var insts_294 vm.Value
	var f_50 vm.Value
	var doseq_seq__28325_51 vm.Value
	var doseq_loop__28326_52 vm.Value
	var label_53 vm.Value
	var insts_54 vm.Value
	var vec__28329_55 vm.Value
	var i_56 vm.Value
	var ins_57 vm.Value
	var op_58 vm.Value
	var arg__28380_76 vm.Value
	var arg__28387_81 vm.Value
	var doseq_seq__28327_82 vm.Value
	var f_59 vm.Value
	var doseq_seq__28325_60 vm.Value
	var doseq_loop__28326_61 vm.Value
	var label_62 vm.Value
	var insts_63 vm.Value
	var vec__28329_64 vm.Value
	var i_65 vm.Value
	var ins_66 vm.Value
	var op_67 vm.Value
	var v285 vm.Value
	var doseq_loop__28328_83 vm.Value
	var label_84 vm.Value
	var i_85 vm.Value
	var ins_86 vm.Value
	var insts_87 vm.Value
	var op_88 vm.Value
	var f_90 vm.Value
	var doseq_seq__28325_91 vm.Value
	var doseq_loop__28326_92 vm.Value
	var vec__28329_93 vm.Value
	var doseq_seq__28327_94 vm.Value
	var doseq_loop__28328_95 vm.Value
	var label_96 vm.Value
	var i_97 vm.Value
	var ins_98 vm.Value
	var insts_99 vm.Value
	var op_100 vm.Value
	var r_114 vm.Value
	var referent_116 vm.Value
	var ref_op_120 vm.Value
	var def_block_124 vm.Value
	var use_block_128 vm.Value
	var and__x_162 vm.Value
	var f_101 vm.Value
	var doseq_seq__28325_102 vm.Value
	var doseq_loop__28326_103 vm.Value
	var vec__28329_104 vm.Value
	var doseq_seq__28327_105 vm.Value
	var doseq_loop__28328_106 vm.Value
	var label_107 vm.Value
	var i_108 vm.Value
	var ins_109 vm.Value
	var insts_110 vm.Value
	var op_111 vm.Value
	var v269 vm.Value
	var f_270 vm.Value
	var doseq_seq__28325_271 vm.Value
	var doseq_loop__28326_272 vm.Value
	var vec__28329_273 vm.Value
	var doseq_seq__28327_274 vm.Value
	var doseq_loop__28328_275 vm.Value
	var label_276 vm.Value
	var i_277 vm.Value
	var ins_278 vm.Value
	var insts_279 vm.Value
	var op_280 vm.Value
	var v282 vm.Value
	var f_129 vm.Value
	var doseq_seq__28325_130 vm.Value
	var doseq_loop__28326_131 vm.Value
	var vec__28329_132 vm.Value
	var doseq_seq__28327_133 vm.Value
	var doseq_loop__28328_134 vm.Value
	var label_135 vm.Value
	var i_136 vm.Value
	var ins_137 vm.Value
	var insts_138 vm.Value
	var op_139 vm.Value
	var r_140 vm.Value
	var referent_141 vm.Value
	var ref_op_142 vm.Value
	var def_block_143 vm.Value
	var use_block_144 vm.Value
	var arg__28450_240 vm.Value
	var arg__28481_259 vm.Value
	var v260 vm.Value
	var v262 vm.Value
	var f_145 vm.Value
	var doseq_seq__28325_146 vm.Value
	var doseq_loop__28326_147 vm.Value
	var vec__28329_148 vm.Value
	var doseq_seq__28327_149 vm.Value
	var doseq_loop__28328_150 vm.Value
	var label_151 vm.Value
	var i_152 vm.Value
	var ins_153 vm.Value
	var insts_154 vm.Value
	var op_155 vm.Value
	var r_156 vm.Value
	var referent_157 vm.Value
	var ref_op_158 vm.Value
	var def_block_159 vm.Value
	var use_block_160 vm.Value
	var v265 vm.Value
	var f_163 vm.Value
	var doseq_seq__28325_164 vm.Value
	var doseq_loop__28326_165 vm.Value
	var vec__28329_166 vm.Value
	var doseq_seq__28327_167 vm.Value
	var doseq_loop__28328_168 vm.Value
	var label_169 vm.Value
	var i_170 vm.Value
	var ins_171 vm.Value
	var insts_172 vm.Value
	var op_173 vm.Value
	var r_174 vm.Value
	var referent_175 vm.Value
	var ref_op_176 vm.Value
	var def_block_177 vm.Value
	var use_block_178 vm.Value
	var and__x_179 vm.Value
	var v201 vm.Value
	var f_180 vm.Value
	var doseq_seq__28325_181 vm.Value
	var doseq_loop__28326_182 vm.Value
	var vec__28329_183 vm.Value
	var doseq_seq__28327_184 vm.Value
	var doseq_loop__28328_185 vm.Value
	var label_186 vm.Value
	var i_187 vm.Value
	var ins_188 vm.Value
	var insts_189 vm.Value
	var op_190 vm.Value
	var r_191 vm.Value
	var referent_192 vm.Value
	var ref_op_193 vm.Value
	var def_block_194 vm.Value
	var use_block_195 vm.Value
	var and__x_196 vm.Value
	var v204 vm.Value
	var f_205 vm.Value
	var doseq_seq__28325_206 vm.Value
	var doseq_loop__28326_207 vm.Value
	var vec__28329_208 vm.Value
	var doseq_seq__28327_209 vm.Value
	var doseq_loop__28328_210 vm.Value
	var label_211 vm.Value
	var i_212 vm.Value
	var ins_213 vm.Value
	var insts_214 vm.Value
	var op_215 vm.Value
	var r_216 vm.Value
	var referent_217 vm.Value
	var ref_op_218 vm.Value
	var def_block_219 vm.Value
	var use_block_220 vm.Value
	var and__x_221 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__28334_5, insts_6, arg__28340_10, arg__28347_15, doseq_seq__28325_16, doseq_loop__28326_17, label_18, insts_19, f_21, doseq_seq__28325_22, doseq_loop__28326_23, label_24, insts_25, vec__28329_33, i_39, ins_45, op_49, v71, f_26, doseq_seq__28325_27, doseq_loop__28326_28, label_29, insts_30, v289, f_290, doseq_seq__28325_291, doseq_loop__28326_292, label_293, insts_294, f_50, doseq_seq__28325_51, doseq_loop__28326_52, label_53, insts_54, vec__28329_55, i_56, ins_57, op_58, arg__28380_76, arg__28387_81, doseq_seq__28327_82, f_59, doseq_seq__28325_60, doseq_loop__28326_61, label_62, insts_63, vec__28329_64, i_65, ins_66, op_67, v285, doseq_loop__28328_83, label_84, i_85, ins_86, insts_87, op_88, f_90, doseq_seq__28325_91, doseq_loop__28326_92, vec__28329_93, doseq_seq__28327_94, doseq_loop__28328_95, label_96, i_97, ins_98, insts_99, op_100, r_114, referent_116, ref_op_120, def_block_124, use_block_128, and__x_162, f_101, doseq_seq__28325_102, doseq_loop__28326_103, vec__28329_104, doseq_seq__28327_105, doseq_loop__28328_106, label_107, i_108, ins_109, insts_110, op_111, v269, f_270, doseq_seq__28325_271, doseq_loop__28326_272, vec__28329_273, doseq_seq__28327_274, doseq_loop__28328_275, label_276, i_277, ins_278, insts_279, op_280, v282, f_129, doseq_seq__28325_130, doseq_loop__28326_131, vec__28329_132, doseq_seq__28327_133, doseq_loop__28328_134, label_135, i_136, ins_137, insts_138, op_139, r_140, referent_141, ref_op_142, def_block_143, use_block_144, arg__28450_240, arg__28481_259, v260, v262, f_145, doseq_seq__28325_146, doseq_loop__28326_147, vec__28329_148, doseq_seq__28327_149, doseq_loop__28328_150, label_151, i_152, ins_153, insts_154, op_155, r_156, referent_157, ref_op_158, def_block_159, use_block_160, v265, f_163, doseq_seq__28325_164, doseq_loop__28326_165, vec__28329_166, doseq_seq__28327_167, doseq_loop__28328_168, label_169, i_170, ins_171, insts_172, op_173, r_174, referent_175, ref_op_176, def_block_177, use_block_178, and__x_179, v201, f_180, doseq_seq__28325_181, doseq_loop__28326_182, vec__28329_183, doseq_seq__28327_184, doseq_loop__28328_185, label_186, i_187, ins_188, insts_189, op_190, r_191, referent_192, ref_op_193, def_block_194, use_block_195, and__x_196, v204, f_205, doseq_seq__28325_206, doseq_loop__28326_207, vec__28329_208, doseq_seq__28327_209, doseq_loop__28328_210, label_211, i_212, ins_213, insts_214, op_215, r_216, referent_217, ref_op_218, def_block_219, use_block_220, and__x_221
	arg__28334_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	insts_6, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__28334_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__28340_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map-indexed").Deref(), []vm.Value{rt.LookupVar("clojure.core", "vector").Deref(), insts_6})
	if callErr != nil {
		return nil, callErr
	}
	arg__28347_15, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map-indexed").Deref(), []vm.Value{rt.LookupVar("clojure.core", "vector").Deref(), insts_6})
	if callErr != nil {
		return nil, callErr
	}
	doseq_seq__28325_16, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{arg__28347_15})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__28326_17 = doseq_seq__28325_16
	label_18 = arg1
	insts_19 = insts_6
	goto b1
b1:
	;
	if vm.IsTruthy(doseq_loop__28326_17) {
		f_21 = arg0
		doseq_seq__28325_22 = doseq_seq__28325_16
		doseq_loop__28326_23 = doseq_loop__28326_17
		label_24 = label_18
		insts_25 = insts_19
		goto b2
	} else {
		f_26 = arg0
		doseq_seq__28325_27 = doseq_seq__28325_16
		doseq_loop__28326_28 = doseq_loop__28326_17
		label_29 = label_18
		insts_30 = insts_19
		goto b3
	}
b2:
	;
	vec__28329_33, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{doseq_loop__28326_23})
	if callErr != nil {
		return nil, callErr
	}
	i_39, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{vec__28329_33, vm.Int(0), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	ins_45, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{vec__28329_33, vm.Int(1), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	op_49, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{ins_45, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	v71, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not=").Deref(), []vm.Value{op_49, vm.Keyword("invalid")})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v71) {
		f_50 = f_21
		doseq_seq__28325_51 = doseq_seq__28325_22
		doseq_loop__28326_52 = doseq_loop__28326_23
		label_53 = label_24
		insts_54 = insts_25
		vec__28329_55 = vec__28329_33
		i_56 = i_39
		ins_57 = ins_45
		op_58 = op_49
		goto b5
	} else {
		f_59 = f_21
		doseq_seq__28325_60 = doseq_seq__28325_22
		doseq_loop__28326_61 = doseq_loop__28326_23
		label_62 = label_24
		insts_63 = insts_25
		vec__28329_64 = vec__28329_33
		i_65 = i_39
		ins_66 = ins_45
		op_67 = op_49
		goto b6
	}
b3:
	;
	v289 = vm.NIL
	f_290 = f_26
	doseq_seq__28325_291 = doseq_seq__28325_27
	doseq_loop__28326_292 = doseq_loop__28326_28
	label_293 = label_29
	insts_294 = insts_30
	goto b4
b4:
	;
	return v289, nil
b5:
	;
	arg__28380_76, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{ins_57, vm.Int(1)})
	if callErr != nil {
		return nil, callErr
	}
	arg__28387_81, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{ins_57, vm.Int(1)})
	if callErr != nil {
		return nil, callErr
	}
	doseq_seq__28327_82, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{arg__28387_81})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__28328_83 = doseq_seq__28327_82
	label_84 = label_53
	i_85 = i_56
	ins_86 = ins_57
	insts_87 = insts_54
	op_88 = op_58
	goto b8
b6:
	;
	v285, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "next").Deref(), []vm.Value{doseq_loop__28326_61})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__28326_17 = v285
	label_18 = label_62
	insts_19 = insts_63
	goto b1
b8:
	;
	if vm.IsTruthy(doseq_loop__28328_83) {
		f_90 = f_50
		doseq_seq__28325_91 = doseq_seq__28325_51
		doseq_loop__28326_92 = doseq_loop__28326_52
		vec__28329_93 = vec__28329_55
		doseq_seq__28327_94 = doseq_seq__28327_82
		doseq_loop__28328_95 = doseq_loop__28328_83
		label_96 = label_84
		i_97 = i_85
		ins_98 = ins_86
		insts_99 = insts_87
		op_100 = op_88
		goto b9
	} else {
		f_101 = f_50
		doseq_seq__28325_102 = doseq_seq__28325_51
		doseq_loop__28326_103 = doseq_loop__28326_52
		vec__28329_104 = vec__28329_55
		doseq_seq__28327_105 = doseq_seq__28327_82
		doseq_loop__28328_106 = doseq_loop__28328_83
		label_107 = label_84
		i_108 = i_85
		ins_109 = ins_86
		insts_110 = insts_87
		op_111 = op_88
		goto b10
	}
b9:
	;
	r_114, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{doseq_loop__28328_95})
	if callErr != nil {
		return nil, callErr
	}
	referent_116, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{insts_99, r_114})
	if callErr != nil {
		return nil, callErr
	}
	ref_op_120, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{referent_116, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	def_block_124, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{referent_116, vm.Int(3)})
	if callErr != nil {
		return nil, callErr
	}
	use_block_128, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{ins_98, vm.Int(3)})
	if callErr != nil {
		return nil, callErr
	}
	and__x_162, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not=").Deref(), []vm.Value{def_block_124, use_block_128})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(and__x_162) {
		f_163 = f_90
		doseq_seq__28325_164 = doseq_seq__28325_91
		doseq_loop__28326_165 = doseq_loop__28326_92
		vec__28329_166 = vec__28329_93
		doseq_seq__28327_167 = doseq_seq__28327_94
		doseq_loop__28328_168 = doseq_loop__28328_95
		label_169 = label_96
		i_170 = i_97
		ins_171 = ins_98
		insts_172 = insts_99
		op_173 = op_100
		r_174 = r_114
		referent_175 = referent_116
		ref_op_176 = ref_op_120
		def_block_177 = def_block_124
		use_block_178 = use_block_128
		and__x_179 = and__x_162
		goto b15
	} else {
		f_180 = f_90
		doseq_seq__28325_181 = doseq_seq__28325_91
		doseq_loop__28326_182 = doseq_loop__28326_92
		vec__28329_183 = vec__28329_93
		doseq_seq__28327_184 = doseq_seq__28327_94
		doseq_loop__28328_185 = doseq_loop__28328_95
		label_186 = label_96
		i_187 = i_97
		ins_188 = ins_98
		insts_189 = insts_99
		op_190 = op_100
		r_191 = r_114
		referent_192 = referent_116
		ref_op_193 = ref_op_120
		def_block_194 = def_block_124
		use_block_195 = use_block_128
		and__x_196 = and__x_162
		goto b16
	}
b10:
	;
	v269 = vm.NIL
	f_270 = f_101
	doseq_seq__28325_271 = doseq_seq__28325_102
	doseq_loop__28326_272 = doseq_loop__28326_103
	vec__28329_273 = vec__28329_104
	doseq_seq__28327_274 = doseq_seq__28327_105
	doseq_loop__28328_275 = doseq_loop__28328_106
	label_276 = label_107
	i_277 = i_108
	ins_278 = ins_109
	insts_279 = insts_110
	op_280 = op_111
	goto b11
b11:
	;
	v282, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "next").Deref(), []vm.Value{doseq_loop__28326_272})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__28326_17 = v282
	label_18 = label_276
	insts_19 = insts_279
	goto b1
b12:
	;
	arg__28450_240, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("validate after "), label_135, vm.String(": Inst #"), i_136, vm.String(" (op="), op_139, vm.String(", block="), use_block_144, vm.String(")"), vm.String(" has cross-block ref "), r_140, vm.String(" (defined in block "), def_block_143, vm.String(")")})
	if callErr != nil {
		return nil, callErr
	}
	arg__28481_259, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("validate after "), label_135, vm.String(": Inst #"), i_136, vm.String(" (op="), op_139, vm.String(", block="), use_block_144, vm.String(")"), vm.String(" has cross-block ref "), r_140, vm.String(" (defined in block "), def_block_143, vm.String(")")})
	if callErr != nil {
		return nil, callErr
	}
	v260, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "throw").Deref(), []vm.Value{arg__28481_259})
	if callErr != nil {
		return nil, callErr
	}
	v262, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "next").Deref(), []vm.Value{doseq_loop__28328_134})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__28328_83 = v262
	label_84 = label_135
	i_85 = i_136
	ins_86 = ins_137
	insts_87 = insts_138
	op_88 = op_139
	goto b8
b13:
	;
	v265, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "next").Deref(), []vm.Value{doseq_loop__28328_150})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__28328_83 = v265
	label_84 = label_151
	i_85 = i_152
	ins_86 = ins_153
	insts_87 = insts_154
	op_88 = op_155
	goto b8
b15:
	;
	v201, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not=").Deref(), []vm.Value{vm.Keyword("block-arg"), ref_op_176})
	if callErr != nil {
		return nil, callErr
	}
	v204 = v201
	f_205 = f_163
	doseq_seq__28325_206 = doseq_seq__28325_164
	doseq_loop__28326_207 = doseq_loop__28326_165
	vec__28329_208 = vec__28329_166
	doseq_seq__28327_209 = doseq_seq__28327_167
	doseq_loop__28328_210 = doseq_loop__28328_168
	label_211 = label_169
	i_212 = i_170
	ins_213 = ins_171
	insts_214 = insts_172
	op_215 = op_173
	r_216 = r_174
	referent_217 = referent_175
	ref_op_218 = ref_op_176
	def_block_219 = def_block_177
	use_block_220 = use_block_178
	and__x_221 = and__x_179
	goto b17
b16:
	;
	v204 = and__x_196
	f_205 = f_180
	doseq_seq__28325_206 = doseq_seq__28325_181
	doseq_loop__28326_207 = doseq_loop__28326_182
	vec__28329_208 = vec__28329_183
	doseq_seq__28327_209 = doseq_seq__28327_184
	doseq_loop__28328_210 = doseq_loop__28328_185
	label_211 = label_186
	i_212 = i_187
	ins_213 = ins_188
	insts_214 = insts_189
	op_215 = op_190
	r_216 = r_191
	referent_217 = referent_192
	ref_op_218 = ref_op_193
	def_block_219 = def_block_194
	use_block_220 = use_block_195
	and__x_221 = and__x_196
	goto b17
b17:
	;
	if vm.IsTruthy(v204) {
		f_129 = f_205
		doseq_seq__28325_130 = doseq_seq__28325_206
		doseq_loop__28326_131 = doseq_loop__28326_207
		vec__28329_132 = vec__28329_208
		doseq_seq__28327_133 = doseq_seq__28327_209
		doseq_loop__28328_134 = doseq_loop__28328_210
		label_135 = label_211
		i_136 = i_212
		ins_137 = ins_213
		insts_138 = insts_214
		op_139 = op_215
		r_140 = r_216
		referent_141 = referent_217
		ref_op_142 = ref_op_218
		def_block_143 = def_block_219
		use_block_144 = use_block_220
		goto b12
	} else {
		f_145 = f_205
		doseq_seq__28325_146 = doseq_seq__28325_206
		doseq_loop__28326_147 = doseq_loop__28326_207
		vec__28329_148 = vec__28329_208
		doseq_seq__28327_149 = doseq_seq__28327_209
		doseq_loop__28328_150 = doseq_loop__28328_210
		label_151 = label_211
		i_152 = i_212
		ins_153 = ins_213
		insts_154 = insts_214
		op_155 = op_215
		r_156 = r_216
		referent_157 = referent_217
		ref_op_158 = ref_op_218
		def_block_159 = def_block_219
		use_block_160 = use_block_220
		goto b13
	}
}
func check_branch_if_symmetric_args_BANG_(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__28498_5 vm.Value
	var blocks_6 vm.Value
	var arg__28503_9 vm.Value
	var insts_10 vm.Value
	var i_11 int
	var label_12 vm.Value
	var insts_13 vm.Value
	var blocks_14 vm.Value
	var v265 vm.Value
	var v280 int
	var v295 int
	var v310 vm.Value
	var v325 int
	var v340 string
	var v355 string
	var v370 string
	var v385 string
	var v400 string
	var v415 string
	var v430 string
	var arg__28508_28 vm.Value
	var v29 bool
	var f_17 vm.Value
	var i_18 int
	var label_19 vm.Value
	var insts_20 vm.Value
	var blocks_21 vm.Value
	var v264 vm.Value
	var v279 int
	var v294 int
	var v309 vm.Value
	var v324 int
	var v339 string
	var v354 string
	var v369 string
	var v384 string
	var v399 string
	var v414 string
	var v429 string
	var arg__28515_33 vm.Value
	var term_id_34 vm.Value
	var f_22 vm.Value
	var i_23 int
	var label_24 vm.Value
	var insts_25 vm.Value
	var blocks_26 vm.Value
	var v274 vm.Value
	var v289 int
	var v304 int
	var v319 vm.Value
	var v334 int
	var v349 string
	var v364 string
	var v379 string
	var v394 string
	var v409 string
	var v424 string
	var v439 string
	var v215 vm.Value
	var f_216 vm.Value
	var i_217 int
	var label_218 vm.Value
	var insts_219 vm.Value
	var blocks_220 vm.Value
	var f_35 vm.Value
	var i_36 int
	var label_37 vm.Value
	var insts_38 vm.Value
	var blocks_39 vm.Value
	var term_id_40 vm.Value
	var v270 vm.Value
	var v285 int
	var v300 int
	var v315 vm.Value
	var v330 int
	var v345 string
	var v360 string
	var v375 string
	var v390 string
	var v405 string
	var v420 string
	var v435 string
	var term_76 vm.Value
	var op_80 vm.Value
	var v98 bool
	var f_41 vm.Value
	var i_42 int
	var label_43 vm.Value
	var insts_44 vm.Value
	var blocks_45 vm.Value
	var term_id_46 vm.Value
	var v268 vm.Value
	var v283 int
	var v298 int
	var v313 vm.Value
	var v328 int
	var v343 string
	var v358 string
	var v373 string
	var v388 string
	var v403 string
	var v418 string
	var v433 string
	var v204 vm.Value
	var f_205 vm.Value
	var i_206 int
	var label_207 vm.Value
	var insts_208 vm.Value
	var blocks_209 vm.Value
	var term_id_210 vm.Value
	var v273 vm.Value
	var v288 int
	var v303 int
	var v318 vm.Value
	var v333 int
	var v348 string
	var v363 string
	var v378 string
	var v393 string
	var v408 string
	var v423 string
	var v438 string
	var v211 int
	var f_47 vm.Value
	var i_48 int
	var label_49 vm.Value
	var insts_50 vm.Value
	var blocks_51 vm.Value
	var and__x_52 vm.Value
	var term_id_53 vm.Value
	var v263 vm.Value
	var v278 int
	var v293 int
	var v308 vm.Value
	var v323 int
	var v338 string
	var v353 string
	var v368 string
	var v383 string
	var v398 string
	var v413 string
	var v428 string
	var v63 bool
	var f_54 vm.Value
	var i_55 int
	var label_56 vm.Value
	var insts_57 vm.Value
	var blocks_58 vm.Value
	var and__x_59 vm.Value
	var term_id_60 vm.Value
	var v271 vm.Value
	var v286 int
	var v301 int
	var v316 vm.Value
	var v331 int
	var v346 string
	var v361 string
	var v376 string
	var v391 string
	var v406 string
	var v421 string
	var v436 string
	var v66 vm.Value
	var f_67 vm.Value
	var i_68 int
	var label_69 vm.Value
	var insts_70 vm.Value
	var blocks_71 vm.Value
	var and__x_72 vm.Value
	var term_id_73 vm.Value
	var v267 vm.Value
	var v282 int
	var v297 int
	var v312 vm.Value
	var v327 int
	var v342 string
	var v357 string
	var v372 string
	var v387 string
	var v402 string
	var v417 string
	var v432 string
	var f_81 vm.Value
	var i_82 int
	var label_83 vm.Value
	var insts_84 vm.Value
	var blocks_85 vm.Value
	var term_id_86 vm.Value
	var term_87 vm.Value
	var op_88 vm.Value
	var v261 vm.Value
	var v276 int
	var v291 int
	var v306 vm.Value
	var v321 int
	var v336 string
	var v351 string
	var v366 string
	var v381 string
	var v396 string
	var v411 string
	var v426 string
	var aux_103 vm.Value
	var arg__28538_105 vm.Value
	var arg__28543_108 vm.Value
	var t_args_109 vm.Value
	var arg__28547_111 vm.Value
	var arg__28552_114 vm.Value
	var f_args_115 vm.Value
	var v138 bool
	var f_89 vm.Value
	var i_90 int
	var label_91 vm.Value
	var insts_92 vm.Value
	var blocks_93 vm.Value
	var term_id_94 vm.Value
	var term_95 vm.Value
	var op_96 vm.Value
	var v269 vm.Value
	var v284 int
	var v299 int
	var v314 vm.Value
	var v329 int
	var v344 string
	var v359 string
	var v374 string
	var v389 string
	var v404 string
	var v419 string
	var v434 string
	var v192 vm.Value
	var f_193 vm.Value
	var i_194 int
	var label_195 vm.Value
	var insts_196 vm.Value
	var blocks_197 vm.Value
	var term_id_198 vm.Value
	var term_199 vm.Value
	var op_200 vm.Value
	var v262 vm.Value
	var v277 int
	var v292 int
	var v307 vm.Value
	var v322 int
	var v337 string
	var v352 string
	var v367 string
	var v382 string
	var v397 string
	var v412 string
	var v427 string
	var f_116 vm.Value
	var i_117 int
	var label_118 vm.Value
	var insts_119 vm.Value
	var blocks_120 vm.Value
	var term_id_121 vm.Value
	var term_122 vm.Value
	var op_123 vm.Value
	var aux_124 vm.Value
	var t_args_125 vm.Value
	var f_args_126 vm.Value
	var v272 vm.Value
	var v287 int
	var v302 int
	var v317 vm.Value
	var v332 int
	var v347 string
	var v362 string
	var v377 string
	var v392 string
	var v407 string
	var v422 string
	var v437 string
	var f_127 vm.Value
	var i_128 int
	var label_129 vm.Value
	var insts_130 vm.Value
	var blocks_131 vm.Value
	var term_id_132 vm.Value
	var term_133 vm.Value
	var op_134 vm.Value
	var aux_135 vm.Value
	var t_args_136 vm.Value
	var f_args_137 vm.Value
	var v266 vm.Value
	var v281 int
	var v296 int
	var v311 vm.Value
	var v326 int
	var v341 string
	var v356 string
	var v371 string
	var v386 string
	var v401 string
	var v416 string
	var v431 string
	var arg__28578_157 vm.Value
	var arg__28603_174 vm.Value
	var v175 vm.Value
	var v177 vm.Value
	var f_178 vm.Value
	var i_179 int
	var label_180 vm.Value
	var insts_181 vm.Value
	var blocks_182 vm.Value
	var term_id_183 vm.Value
	var term_184 vm.Value
	var op_185 vm.Value
	var aux_186 vm.Value
	var t_args_187 vm.Value
	var f_args_188 vm.Value
	var v260 vm.Value
	var v275 int
	var v290 int
	var v305 vm.Value
	var v320 int
	var v335 string
	var v350 string
	var v365 string
	var v380 string
	var v395 string
	var v410 string
	var v425 string
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__28498_5, blocks_6, arg__28503_9, insts_10, i_11, label_12, insts_13, blocks_14, v265, v280, v295, v310, v325, v340, v355, v370, v385, v400, v415, v430, arg__28508_28, v29, f_17, i_18, label_19, insts_20, blocks_21, v264, v279, v294, v309, v324, v339, v354, v369, v384, v399, v414, v429, arg__28515_33, term_id_34, f_22, i_23, label_24, insts_25, blocks_26, v274, v289, v304, v319, v334, v349, v364, v379, v394, v409, v424, v439, v215, f_216, i_217, label_218, insts_219, blocks_220, f_35, i_36, label_37, insts_38, blocks_39, term_id_40, v270, v285, v300, v315, v330, v345, v360, v375, v390, v405, v420, v435, term_76, op_80, v98, f_41, i_42, label_43, insts_44, blocks_45, term_id_46, v268, v283, v298, v313, v328, v343, v358, v373, v388, v403, v418, v433, v204, f_205, i_206, label_207, insts_208, blocks_209, term_id_210, v273, v288, v303, v318, v333, v348, v363, v378, v393, v408, v423, v438, v211, f_47, i_48, label_49, insts_50, blocks_51, and__x_52, term_id_53, v263, v278, v293, v308, v323, v338, v353, v368, v383, v398, v413, v428, v63, f_54, i_55, label_56, insts_57, blocks_58, and__x_59, term_id_60, v271, v286, v301, v316, v331, v346, v361, v376, v391, v406, v421, v436, v66, f_67, i_68, label_69, insts_70, blocks_71, and__x_72, term_id_73, v267, v282, v297, v312, v327, v342, v357, v372, v387, v402, v417, v432, f_81, i_82, label_83, insts_84, blocks_85, term_id_86, term_87, op_88, v261, v276, v291, v306, v321, v336, v351, v366, v381, v396, v411, v426, aux_103, arg__28538_105, arg__28543_108, t_args_109, arg__28547_111, arg__28552_114, f_args_115, v138, f_89, i_90, label_91, insts_92, blocks_93, term_id_94, term_95, op_96, v269, v284, v299, v314, v329, v344, v359, v374, v389, v404, v419, v434, v192, f_193, i_194, label_195, insts_196, blocks_197, term_id_198, term_199, op_200, v262, v277, v292, v307, v322, v337, v352, v367, v382, v397, v412, v427, f_116, i_117, label_118, insts_119, blocks_120, term_id_121, term_122, op_123, aux_124, t_args_125, f_args_126, v272, v287, v302, v317, v332, v347, v362, v377, v392, v407, v422, v437, f_127, i_128, label_129, insts_130, blocks_131, term_id_132, term_133, op_134, aux_135, t_args_136, f_args_137, v266, v281, v296, v311, v326, v341, v356, v371, v386, v401, v416, v431, arg__28578_157, arg__28603_174, v175, v177, f_178, i_179, label_180, insts_181, blocks_182, term_id_183, term_184, op_185, aux_186, t_args_187, f_args_188, v260, v275, v290, v305, v320, v335, v350, v365, v380, v395, v410, v425
	arg__28498_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	blocks_6, callErr = rt.InvokeValue(vm.Keyword("blocks"), []vm.Value{arg__28498_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__28503_9, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	insts_10, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__28503_9})
	if callErr != nil {
		return nil, callErr
	}
	i_11 = 0
	label_12 = arg1
	insts_13 = insts_10
	blocks_14 = blocks_6
	v265 = vm.Keyword("term")
	v280 = 0
	v295 = 0
	v310 = vm.Keyword("branch-if")
	v325 = 2
	v340 = "validate after "
	v355 = ": block #"
	v370 = " :branch-if has asymmetric"
	v385 = " branch-target args (true="
	v400 = ", false="
	v415 = ") — lower-block! requires"
	v430 = " true.args == false.args"
	goto b1
b1:
	;
	arg__28508_28, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{blocks_14})
	if callErr != nil {
		return nil, callErr
	}
	v29 = rt.LtValue(vm.Int(i_11), arg__28508_28)
	if v29 {
		f_17 = arg0
		i_18 = i_11
		label_19 = label_12
		insts_20 = insts_13
		blocks_21 = blocks_14
		v264 = v265
		v279 = v280
		v294 = v295
		v309 = v310
		v324 = v325
		v339 = v340
		v354 = v355
		v369 = v370
		v384 = v385
		v399 = v400
		v414 = v415
		v429 = v430
		goto b2
	} else {
		f_22 = arg0
		i_23 = i_11
		label_24 = label_12
		insts_25 = insts_13
		blocks_26 = blocks_14
		v274 = v265
		v289 = v280
		v304 = v295
		v319 = v310
		v334 = v325
		v349 = v340
		v364 = v355
		v379 = v370
		v394 = v385
		v409 = v400
		v424 = v415
		v439 = v430
		goto b3
	}
b2:
	;
	arg__28515_33, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{blocks_21, vm.Int(i_18)})
	if callErr != nil {
		return nil, callErr
	}
	term_id_34, callErr = rt.InvokeValue(v264, []vm.Value{arg__28515_33})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(term_id_34) {
		f_47 = f_17
		i_48 = i_18
		label_49 = label_19
		insts_50 = insts_20
		blocks_51 = blocks_21
		and__x_52 = term_id_34
		term_id_53 = term_id_34
		v263 = v264
		v278 = v279
		v293 = v294
		v308 = v309
		v323 = v324
		v338 = v339
		v353 = v354
		v368 = v369
		v383 = v384
		v398 = v399
		v413 = v414
		v428 = v429
		goto b8
	} else {
		f_54 = f_17
		i_55 = i_18
		label_56 = label_19
		insts_57 = insts_20
		blocks_58 = blocks_21
		and__x_59 = term_id_34
		term_id_60 = term_id_34
		v271 = v264
		v286 = v279
		v301 = v294
		v316 = v309
		v331 = v324
		v346 = v339
		v361 = v354
		v376 = v369
		v391 = v384
		v406 = v399
		v421 = v414
		v436 = v429
		goto b9
	}
b3:
	;
	v215 = vm.NIL
	f_216 = f_22
	i_217 = i_23
	label_218 = label_24
	insts_219 = insts_25
	blocks_220 = blocks_26
	goto b4
b4:
	;
	return v215, nil
b5:
	;
	term_76, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{insts_38, term_id_40})
	if callErr != nil {
		return nil, callErr
	}
	op_80, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{term_76, vm.Int(v300)})
	if callErr != nil {
		return nil, callErr
	}
	v98 = op_80 == v315
	if v98 {
		f_81 = f_35
		i_82 = i_36
		label_83 = label_37
		insts_84 = insts_38
		blocks_85 = blocks_39
		term_id_86 = term_id_40
		term_87 = term_76
		op_88 = op_80
		v261 = v270
		v276 = v285
		v291 = v300
		v306 = v315
		v321 = v330
		v336 = v345
		v351 = v360
		v366 = v375
		v381 = v390
		v396 = v405
		v411 = v420
		v426 = v435
		goto b11
	} else {
		f_89 = f_35
		i_90 = i_36
		label_91 = label_37
		insts_92 = insts_38
		blocks_93 = blocks_39
		term_id_94 = term_id_40
		term_95 = term_76
		op_96 = op_80
		v269 = v270
		v284 = v285
		v299 = v300
		v314 = v315
		v329 = v330
		v344 = v345
		v359 = v360
		v374 = v375
		v389 = v390
		v404 = v405
		v419 = v420
		v434 = v435
		goto b12
	}
b6:
	;
	v204 = vm.NIL
	f_205 = f_41
	i_206 = i_42
	label_207 = label_43
	insts_208 = insts_44
	blocks_209 = blocks_45
	term_id_210 = term_id_46
	v273 = v268
	v288 = v283
	v303 = v298
	v318 = v313
	v333 = v328
	v348 = v343
	v363 = v358
	v378 = v373
	v393 = v388
	v408 = v403
	v423 = v418
	v438 = v433
	goto b7
b7:
	;
	v211 = i_206 + 1
	i_11 = v211
	label_12 = label_207
	insts_13 = insts_208
	blocks_14 = blocks_209
	v265 = v273
	v280 = v288
	v295 = v303
	v310 = v318
	v325 = v333
	v340 = v348
	v355 = v363
	v370 = v378
	v385 = v393
	v400 = v408
	v415 = v423
	v430 = v438
	goto b1
b8:
	;
	v63 = rt.GtValue(term_id_53, vm.Int(v278))
	v66 = vm.Boolean(v63)
	f_67 = f_47
	i_68 = i_48
	label_69 = label_49
	insts_70 = insts_50
	blocks_71 = blocks_51
	and__x_72 = and__x_52
	term_id_73 = term_id_53
	v267 = v263
	v282 = v278
	v297 = v293
	v312 = v308
	v327 = v323
	v342 = v338
	v357 = v353
	v372 = v368
	v387 = v383
	v402 = v398
	v417 = v413
	v432 = v428
	goto b10
b9:
	;
	v66 = and__x_59
	f_67 = f_54
	i_68 = i_55
	label_69 = label_56
	insts_70 = insts_57
	blocks_71 = blocks_58
	and__x_72 = and__x_59
	term_id_73 = term_id_60
	v267 = v271
	v282 = v286
	v297 = v301
	v312 = v316
	v327 = v331
	v342 = v346
	v357 = v361
	v372 = v376
	v387 = v391
	v402 = v406
	v417 = v421
	v432 = v436
	goto b10
b10:
	;
	if vm.IsTruthy(v66) {
		f_35 = f_67
		i_36 = i_68
		label_37 = label_69
		insts_38 = insts_70
		blocks_39 = blocks_71
		term_id_40 = term_id_73
		v270 = v267
		v285 = v282
		v300 = v297
		v315 = v312
		v330 = v327
		v345 = v342
		v360 = v357
		v375 = v372
		v390 = v387
		v405 = v402
		v420 = v417
		v435 = v432
		goto b5
	} else {
		f_41 = f_67
		i_42 = i_68
		label_43 = label_69
		insts_44 = insts_70
		blocks_45 = blocks_71
		term_id_46 = term_id_73
		v268 = v267
		v283 = v282
		v298 = v297
		v313 = v312
		v328 = v327
		v343 = v342
		v358 = v357
		v373 = v372
		v388 = v387
		v403 = v402
		v418 = v417
		v433 = v432
		goto b6
	}
b11:
	;
	aux_103, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{term_87, vm.Int(v321)})
	if callErr != nil {
		return nil, callErr
	}
	arg__28538_105, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-true").Deref(), []vm.Value{aux_103})
	if callErr != nil {
		return nil, callErr
	}
	arg__28543_108, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-true").Deref(), []vm.Value{aux_103})
	if callErr != nil {
		return nil, callErr
	}
	t_args_109, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-args").Deref(), []vm.Value{arg__28543_108})
	if callErr != nil {
		return nil, callErr
	}
	arg__28547_111, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-false").Deref(), []vm.Value{aux_103})
	if callErr != nil {
		return nil, callErr
	}
	arg__28552_114, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-false").Deref(), []vm.Value{aux_103})
	if callErr != nil {
		return nil, callErr
	}
	f_args_115, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-args").Deref(), []vm.Value{arg__28552_114})
	if callErr != nil {
		return nil, callErr
	}
	v138 = t_args_109 == f_args_115
	if v138 {
		f_116 = f_81
		i_117 = i_82
		label_118 = label_83
		insts_119 = insts_84
		blocks_120 = blocks_85
		term_id_121 = term_id_86
		term_122 = term_87
		op_123 = op_88
		aux_124 = aux_103
		t_args_125 = t_args_109
		f_args_126 = f_args_115
		v272 = v261
		v287 = v276
		v302 = v291
		v317 = v306
		v332 = v321
		v347 = v336
		v362 = v351
		v377 = v366
		v392 = v381
		v407 = v396
		v422 = v411
		v437 = v426
		goto b14
	} else {
		f_127 = f_81
		i_128 = i_82
		label_129 = label_83
		insts_130 = insts_84
		blocks_131 = blocks_85
		term_id_132 = term_id_86
		term_133 = term_87
		op_134 = op_88
		aux_135 = aux_103
		t_args_136 = t_args_109
		f_args_137 = f_args_115
		v266 = v261
		v281 = v276
		v296 = v291
		v311 = v306
		v326 = v321
		v341 = v336
		v356 = v351
		v371 = v366
		v386 = v381
		v401 = v396
		v416 = v411
		v431 = v426
		goto b15
	}
b12:
	;
	v192 = vm.NIL
	f_193 = f_89
	i_194 = i_90
	label_195 = label_91
	insts_196 = insts_92
	blocks_197 = blocks_93
	term_id_198 = term_id_94
	term_199 = term_95
	op_200 = op_96
	v262 = v269
	v277 = v284
	v292 = v299
	v307 = v314
	v322 = v329
	v337 = v344
	v352 = v359
	v367 = v374
	v382 = v389
	v397 = v404
	v412 = v419
	v427 = v434
	goto b13
b13:
	;
	v204 = v192
	f_205 = f_193
	i_206 = i_194
	label_207 = label_195
	insts_208 = insts_196
	blocks_209 = blocks_197
	term_id_210 = term_id_198
	v273 = v262
	v288 = v277
	v303 = v292
	v318 = v307
	v333 = v322
	v348 = v337
	v363 = v352
	v378 = v367
	v393 = v382
	v408 = v397
	v423 = v412
	v438 = v427
	goto b7
b14:
	;
	v177 = vm.NIL
	f_178 = f_116
	i_179 = i_117
	label_180 = label_118
	insts_181 = insts_119
	blocks_182 = blocks_120
	term_id_183 = term_id_121
	term_184 = term_122
	op_185 = op_123
	aux_186 = aux_124
	t_args_187 = t_args_125
	f_args_188 = f_args_126
	v260 = v272
	v275 = v287
	v290 = v302
	v305 = v317
	v320 = v332
	v335 = v347
	v350 = v362
	v365 = v377
	v380 = v392
	v395 = v407
	v410 = v422
	v425 = v437
	goto b16
b15:
	;
	arg__28578_157, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(v341), label_129, vm.String(v356), vm.Int(i_128), vm.String(v371), vm.String(v386), t_args_136, vm.String(v401), f_args_137, vm.String(v416), vm.String(v431)})
	if callErr != nil {
		return nil, callErr
	}
	arg__28603_174, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(v341), label_129, vm.String(v356), vm.Int(i_128), vm.String(v371), vm.String(v386), t_args_136, vm.String(v401), f_args_137, vm.String(v416), vm.String(v431)})
	if callErr != nil {
		return nil, callErr
	}
	v175, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "throw").Deref(), []vm.Value{arg__28603_174})
	if callErr != nil {
		return nil, callErr
	}
	v177 = v175
	f_178 = f_127
	i_179 = i_128
	label_180 = label_129
	insts_181 = insts_130
	blocks_182 = blocks_131
	term_id_183 = term_id_132
	term_184 = term_133
	op_185 = op_134
	aux_186 = aux_135
	t_args_187 = t_args_136
	f_args_188 = f_args_137
	v260 = v266
	v275 = v281
	v290 = v296
	v305 = v311
	v320 = v326
	v335 = v341
	v350 = v356
	v365 = v371
	v380 = v386
	v395 = v401
	v410 = v416
	v425 = v431
	goto b16
b16:
	;
	v192 = v177
	f_193 = f_178
	i_194 = i_179
	label_195 = label_180
	insts_196 = insts_181
	blocks_197 = blocks_182
	term_id_198 = term_id_183
	term_199 = term_184
	op_200 = op_185
	v262 = v260
	v277 = v275
	v292 = v290
	v307 = v305
	v322 = v320
	v337 = v335
	v352 = v350
	v367 = v365
	v382 = v380
	v397 = v395
	v412 = v410
	v427 = v425
	goto b13
}
func check_inst_shapes_BANG_(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var i_2 int
	var label_3 vm.Value
	var insts_4 vm.Value
	var v104 int
	var v113 string
	var v122 string
	var v131 string
	var arg__28609_14 vm.Value
	var v15 bool
	var i_7 int
	var label_8 vm.Value
	var insts_9 vm.Value
	var v103 int
	var v112 string
	var v121 string
	var v130 string
	var ins_18 vm.Value
	var and__x_28 vm.Value
	var i_10 int
	var label_11 vm.Value
	var insts_12 vm.Value
	var v110 int
	var v119 string
	var v128 string
	var v137 string
	var v83 vm.Value
	var i_84 int
	var label_85 vm.Value
	var insts_86 vm.Value
	var i_19 int
	var label_20 vm.Value
	var insts_21 vm.Value
	var ins_22 vm.Value
	var v107 int
	var v116 string
	var v125 string
	var v134 string
	var i_23 int
	var label_24 vm.Value
	var insts_25 vm.Value
	var ins_26 vm.Value
	var v106 int
	var v115 string
	var v124 string
	var v133 string
	var arg__28636_62 vm.Value
	var arg__28651_71 vm.Value
	var v72 vm.Value
	var v74 vm.Value
	var i_75 int
	var label_76 vm.Value
	var insts_77 vm.Value
	var ins_78 vm.Value
	var v109 int
	var v118 string
	var v127 string
	var v136 string
	var v79 int
	var i_29 int
	var label_30 vm.Value
	var insts_31 vm.Value
	var ins_32 vm.Value
	var and__x_33 vm.Value
	var v102 int
	var v111 string
	var v120 string
	var v129 string
	var arg__28622_42 vm.Value
	var v43 bool
	var i_34 int
	var label_35 vm.Value
	var insts_36 vm.Value
	var ins_37 vm.Value
	var and__x_38 vm.Value
	var v108 int
	var v117 string
	var v126 string
	var v135 string
	var v46 vm.Value
	var i_47 int
	var label_48 vm.Value
	var insts_49 vm.Value
	var ins_50 vm.Value
	var and__x_51 vm.Value
	var v105 int
	var v114 string
	var v123 string
	var v132 string
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = i_2, label_3, insts_4, v104, v113, v122, v131, arg__28609_14, v15, i_7, label_8, insts_9, v103, v112, v121, v130, ins_18, and__x_28, i_10, label_11, insts_12, v110, v119, v128, v137, v83, i_84, label_85, insts_86, i_19, label_20, insts_21, ins_22, v107, v116, v125, v134, i_23, label_24, insts_25, ins_26, v106, v115, v124, v133, arg__28636_62, arg__28651_71, v72, v74, i_75, label_76, insts_77, ins_78, v109, v118, v127, v136, v79, i_29, label_30, insts_31, ins_32, and__x_33, v102, v111, v120, v129, arg__28622_42, v43, i_34, label_35, insts_36, ins_37, and__x_38, v108, v117, v126, v135, v46, i_47, label_48, insts_49, ins_50, and__x_51, v105, v114, v123, v132
	i_2 = 0
	label_3 = arg1
	insts_4 = arg0
	v104 = 6
	v113 = "validate after "
	v122 = ": Inst #"
	v131 = " has bad shape: "
	goto b1
b1:
	;
	arg__28609_14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{insts_4})
	if callErr != nil {
		return nil, callErr
	}
	v15 = rt.LtValue(vm.Int(i_2), arg__28609_14)
	if v15 {
		i_7 = i_2
		label_8 = label_3
		insts_9 = insts_4
		v103 = v104
		v112 = v113
		v121 = v122
		v130 = v131
		goto b2
	} else {
		i_10 = i_2
		label_11 = label_3
		insts_12 = insts_4
		v110 = v104
		v119 = v113
		v128 = v122
		v137 = v131
		goto b3
	}
b2:
	;
	ins_18, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{insts_9, vm.Int(i_7)})
	if callErr != nil {
		return nil, callErr
	}
	and__x_28, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector?").Deref(), []vm.Value{ins_18})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(and__x_28) {
		i_29 = i_7
		label_30 = label_8
		insts_31 = insts_9
		ins_32 = ins_18
		and__x_33 = and__x_28
		v102 = v103
		v111 = v112
		v120 = v121
		v129 = v130
		goto b8
	} else {
		i_34 = i_7
		label_35 = label_8
		insts_36 = insts_9
		ins_37 = ins_18
		and__x_38 = and__x_28
		v108 = v103
		v117 = v112
		v126 = v121
		v135 = v130
		goto b9
	}
b3:
	;
	v83 = vm.NIL
	i_84 = i_10
	label_85 = label_11
	insts_86 = insts_12
	goto b4
b4:
	;
	return v83, nil
b5:
	;
	v74 = vm.NIL
	i_75 = i_19
	label_76 = label_20
	insts_77 = insts_21
	ins_78 = ins_22
	v109 = v107
	v118 = v116
	v127 = v125
	v136 = v134
	goto b7
b6:
	;
	arg__28636_62, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(v115), label_24, vm.String(v124), vm.Int(i_23), vm.String(v133), ins_26})
	if callErr != nil {
		return nil, callErr
	}
	arg__28651_71, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(v115), label_24, vm.String(v124), vm.Int(i_23), vm.String(v133), ins_26})
	if callErr != nil {
		return nil, callErr
	}
	v72, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "throw").Deref(), []vm.Value{arg__28651_71})
	if callErr != nil {
		return nil, callErr
	}
	v74 = v72
	i_75 = i_23
	label_76 = label_24
	insts_77 = insts_25
	ins_78 = ins_26
	v109 = v106
	v118 = v115
	v127 = v124
	v136 = v133
	goto b7
b7:
	;
	v79 = i_75 + 1
	i_2 = v79
	label_3 = label_76
	insts_4 = insts_77
	v104 = v109
	v113 = v118
	v122 = v127
	v131 = v136
	goto b1
b8:
	;
	arg__28622_42, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{ins_32})
	if callErr != nil {
		return nil, callErr
	}
	v43 = arg__28622_42 == vm.Int(v102)
	v46 = vm.Boolean(v43)
	i_47 = i_29
	label_48 = label_30
	insts_49 = insts_31
	ins_50 = ins_32
	and__x_51 = and__x_33
	v105 = v102
	v114 = v111
	v123 = v120
	v132 = v129
	goto b10
b9:
	;
	v46 = and__x_38
	i_47 = i_34
	label_48 = label_35
	insts_49 = insts_36
	ins_50 = ins_37
	and__x_51 = and__x_38
	v105 = v108
	v114 = v117
	v123 = v126
	v132 = v135
	goto b10
b10:
	;
	if vm.IsTruthy(v46) {
		i_19 = i_47
		label_20 = label_48
		insts_21 = insts_49
		ins_22 = ins_50
		v107 = v105
		v116 = v114
		v125 = v123
		v134 = v132
		goto b5
	} else {
		i_23 = i_47
		label_24 = label_48
		insts_25 = insts_49
		ins_26 = ins_50
		v106 = v105
		v115 = v114
		v124 = v123
		v133 = v132
		goto b6
	}
}
func check_refs_in_range_BANG_(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var n_3 vm.Value
	var i_4 int
	var label_5 vm.Value
	var insts_6 vm.Value
	var n_7 vm.Value
	var v306 int
	var v325 int
	var v344 vm.Value
	var v363 int
	var v382 string
	var v401 string
	var v420 string
	var v439 string
	var v458 string
	var v477 string
	var v18 bool
	var i_10 int
	var label_11 vm.Value
	var insts_12 vm.Value
	var n_13 vm.Value
	var v305 int
	var v324 int
	var v343 vm.Value
	var v362 int
	var v381 string
	var v400 string
	var v419 string
	var v438 string
	var v457 string
	var v476 string
	var ins_21 vm.Value
	var op_25 vm.Value
	var refs_29 vm.Value
	var v47 vm.Value
	var i_14 int
	var label_15 vm.Value
	var insts_16 vm.Value
	var n_17 vm.Value
	var v318 int
	var v337 int
	var v356 vm.Value
	var v375 int
	var v394 string
	var v413 string
	var v432 string
	var v451 string
	var v470 string
	var v489 string
	var v256 vm.Value
	var i_257 int
	var label_258 vm.Value
	var insts_259 vm.Value
	var n_260 vm.Value
	var i_30 int
	var label_31 vm.Value
	var insts_32 vm.Value
	var n_33 vm.Value
	var ins_34 vm.Value
	var op_35 vm.Value
	var refs_36 vm.Value
	var v312 int
	var v331 int
	var v350 vm.Value
	var v369 int
	var v388 string
	var v407 string
	var v426 string
	var v445 string
	var v464 string
	var v483 string
	var doseq_seq__28653_50 vm.Value
	var i_37 int
	var label_38 vm.Value
	var insts_39 vm.Value
	var n_40 vm.Value
	var ins_41 vm.Value
	var op_42 vm.Value
	var refs_43 vm.Value
	var v309 int
	var v328 int
	var v347 vm.Value
	var v366 int
	var v385 string
	var v404 string
	var v423 string
	var v442 string
	var v461 string
	var v480 string
	var v244 vm.Value
	var i_245 int
	var label_246 vm.Value
	var insts_247 vm.Value
	var n_248 vm.Value
	var ins_249 vm.Value
	var op_250 vm.Value
	var refs_251 vm.Value
	var v317 int
	var v336 int
	var v355 vm.Value
	var v374 int
	var v393 string
	var v412 string
	var v431 string
	var v450 string
	var v469 string
	var v488 string
	var v252 int
	var doseq_loop__28654_51 vm.Value
	var label_52 vm.Value
	var i_53 int
	var n_54 vm.Value
	var op_55 vm.Value
	var v304 int
	var v323 int
	var v342 vm.Value
	var v361 int
	var v380 string
	var v399 string
	var v418 string
	var v437 string
	var v456 string
	var v475 string
	var insts_57 vm.Value
	var ins_58 vm.Value
	var refs_59 vm.Value
	var doseq_seq__28653_60 vm.Value
	var doseq_loop__28654_61 vm.Value
	var label_62 vm.Value
	var i_63 int
	var n_64 vm.Value
	var op_65 vm.Value
	var v313 int
	var v332 int
	var v351 vm.Value
	var v370 int
	var v389 string
	var v408 string
	var v427 string
	var v446 string
	var v465 string
	var v484 string
	var r_77 vm.Value
	var or__x_99 vm.Value
	var insts_66 vm.Value
	var ins_67 vm.Value
	var refs_68 vm.Value
	var doseq_seq__28653_69 vm.Value
	var doseq_loop__28654_70 vm.Value
	var label_71 vm.Value
	var i_72 int
	var n_73 vm.Value
	var op_74 vm.Value
	var v308 int
	var v327 int
	var v346 vm.Value
	var v365 int
	var v384 string
	var v403 string
	var v422 string
	var v441 string
	var v460 string
	var v479 string
	var v231 vm.Value
	var insts_232 vm.Value
	var ins_233 vm.Value
	var refs_234 vm.Value
	var doseq_seq__28653_235 vm.Value
	var doseq_loop__28654_236 vm.Value
	var label_237 vm.Value
	var i_238 int
	var n_239 vm.Value
	var op_240 vm.Value
	var v301 int
	var v320 int
	var v339 vm.Value
	var v358 int
	var v377 string
	var v396 string
	var v415 string
	var v434 string
	var v453 string
	var v472 string
	var insts_78 vm.Value
	var ins_79 vm.Value
	var refs_80 vm.Value
	var doseq_seq__28653_81 vm.Value
	var doseq_loop__28654_82 vm.Value
	var label_83 vm.Value
	var i_84 int
	var n_85 vm.Value
	var op_86 vm.Value
	var r_87 vm.Value
	var v311 int
	var v330 int
	var v349 vm.Value
	var v368 int
	var v387 string
	var v406 string
	var v425 string
	var v444 string
	var v463 string
	var v482 string
	var arg__28718_193 vm.Value
	var arg__28745_210 vm.Value
	var v211 vm.Value
	var insts_88 vm.Value
	var ins_89 vm.Value
	var refs_90 vm.Value
	var doseq_seq__28653_91 vm.Value
	var doseq_loop__28654_92 vm.Value
	var label_93 vm.Value
	var i_94 int
	var n_95 vm.Value
	var op_96 vm.Value
	var r_97 vm.Value
	var v302 int
	var v321 int
	var v340 vm.Value
	var v359 int
	var v378 string
	var v397 string
	var v416 string
	var v435 string
	var v454 string
	var v473 string
	var v215 vm.Value
	var insts_216 vm.Value
	var ins_217 vm.Value
	var refs_218 vm.Value
	var doseq_seq__28653_219 vm.Value
	var doseq_loop__28654_220 vm.Value
	var label_221 vm.Value
	var i_222 int
	var n_223 vm.Value
	var op_224 vm.Value
	var r_225 vm.Value
	var v316 int
	var v335 int
	var v354 vm.Value
	var v373 int
	var v392 string
	var v411 string
	var v430 string
	var v449 string
	var v468 string
	var v487 string
	var v227 vm.Value
	var insts_100 vm.Value
	var ins_101 vm.Value
	var refs_102 vm.Value
	var doseq_seq__28653_103 vm.Value
	var doseq_loop__28654_104 vm.Value
	var label_105 vm.Value
	var i_106 int
	var n_107 vm.Value
	var op_108 vm.Value
	var r_109 vm.Value
	var or__x_110 vm.Value
	var v307 int
	var v326 int
	var v345 vm.Value
	var v364 int
	var v383 string
	var v402 string
	var v421 string
	var v440 string
	var v459 string
	var v478 string
	var insts_111 vm.Value
	var ins_112 vm.Value
	var refs_113 vm.Value
	var doseq_seq__28653_114 vm.Value
	var doseq_loop__28654_115 vm.Value
	var label_116 vm.Value
	var i_117 int
	var n_118 vm.Value
	var op_119 vm.Value
	var r_120 vm.Value
	var or__x_121 vm.Value
	var v300 int
	var v319 int
	var v338 vm.Value
	var v357 int
	var v376 string
	var v395 string
	var v414 string
	var v433 string
	var v452 string
	var v471 string
	var or__x_125 bool
	var v165 vm.Value
	var insts_166 vm.Value
	var ins_167 vm.Value
	var refs_168 vm.Value
	var doseq_seq__28653_169 vm.Value
	var doseq_loop__28654_170 vm.Value
	var label_171 vm.Value
	var i_172 int
	var n_173 vm.Value
	var op_174 vm.Value
	var r_175 vm.Value
	var or__x_176 vm.Value
	var v314 int
	var v333 int
	var v352 vm.Value
	var v371 int
	var v390 string
	var v409 string
	var v428 string
	var v447 string
	var v466 string
	var v485 string
	var insts_126 vm.Value
	var ins_127 vm.Value
	var refs_128 vm.Value
	var doseq_seq__28653_129 vm.Value
	var doseq_loop__28654_130 vm.Value
	var label_131 vm.Value
	var i_132 int
	var n_133 vm.Value
	var op_134 vm.Value
	var r_135 vm.Value
	var or__x_136 bool
	var v310 int
	var v329 int
	var v348 vm.Value
	var v367 int
	var v386 string
	var v405 string
	var v424 string
	var v443 string
	var v462 string
	var v481 string
	var insts_137 vm.Value
	var ins_138 vm.Value
	var refs_139 vm.Value
	var doseq_seq__28653_140 vm.Value
	var doseq_loop__28654_141 vm.Value
	var label_142 vm.Value
	var i_143 int
	var n_144 vm.Value
	var op_145 vm.Value
	var r_146 vm.Value
	var or__x_147 bool
	var v315 int
	var v334 int
	var v353 vm.Value
	var v372 int
	var v391 string
	var v410 string
	var v429 string
	var v448 string
	var v467 string
	var v486 string
	var v150 bool
	var v152 bool
	var insts_153 vm.Value
	var ins_154 vm.Value
	var refs_155 vm.Value
	var doseq_seq__28653_156 vm.Value
	var doseq_loop__28654_157 vm.Value
	var label_158 vm.Value
	var i_159 int
	var n_160 vm.Value
	var op_161 vm.Value
	var r_162 vm.Value
	var or__x_163 vm.Value
	var v303 int
	var v322 int
	var v341 vm.Value
	var v360 int
	var v379 string
	var v398 string
	var v417 string
	var v436 string
	var v455 string
	var v474 string
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = n_3, i_4, label_5, insts_6, n_7, v306, v325, v344, v363, v382, v401, v420, v439, v458, v477, v18, i_10, label_11, insts_12, n_13, v305, v324, v343, v362, v381, v400, v419, v438, v457, v476, ins_21, op_25, refs_29, v47, i_14, label_15, insts_16, n_17, v318, v337, v356, v375, v394, v413, v432, v451, v470, v489, v256, i_257, label_258, insts_259, n_260, i_30, label_31, insts_32, n_33, ins_34, op_35, refs_36, v312, v331, v350, v369, v388, v407, v426, v445, v464, v483, doseq_seq__28653_50, i_37, label_38, insts_39, n_40, ins_41, op_42, refs_43, v309, v328, v347, v366, v385, v404, v423, v442, v461, v480, v244, i_245, label_246, insts_247, n_248, ins_249, op_250, refs_251, v317, v336, v355, v374, v393, v412, v431, v450, v469, v488, v252, doseq_loop__28654_51, label_52, i_53, n_54, op_55, v304, v323, v342, v361, v380, v399, v418, v437, v456, v475, insts_57, ins_58, refs_59, doseq_seq__28653_60, doseq_loop__28654_61, label_62, i_63, n_64, op_65, v313, v332, v351, v370, v389, v408, v427, v446, v465, v484, r_77, or__x_99, insts_66, ins_67, refs_68, doseq_seq__28653_69, doseq_loop__28654_70, label_71, i_72, n_73, op_74, v308, v327, v346, v365, v384, v403, v422, v441, v460, v479, v231, insts_232, ins_233, refs_234, doseq_seq__28653_235, doseq_loop__28654_236, label_237, i_238, n_239, op_240, v301, v320, v339, v358, v377, v396, v415, v434, v453, v472, insts_78, ins_79, refs_80, doseq_seq__28653_81, doseq_loop__28654_82, label_83, i_84, n_85, op_86, r_87, v311, v330, v349, v368, v387, v406, v425, v444, v463, v482, arg__28718_193, arg__28745_210, v211, insts_88, ins_89, refs_90, doseq_seq__28653_91, doseq_loop__28654_92, label_93, i_94, n_95, op_96, r_97, v302, v321, v340, v359, v378, v397, v416, v435, v454, v473, v215, insts_216, ins_217, refs_218, doseq_seq__28653_219, doseq_loop__28654_220, label_221, i_222, n_223, op_224, r_225, v316, v335, v354, v373, v392, v411, v430, v449, v468, v487, v227, insts_100, ins_101, refs_102, doseq_seq__28653_103, doseq_loop__28654_104, label_105, i_106, n_107, op_108, r_109, or__x_110, v307, v326, v345, v364, v383, v402, v421, v440, v459, v478, insts_111, ins_112, refs_113, doseq_seq__28653_114, doseq_loop__28654_115, label_116, i_117, n_118, op_119, r_120, or__x_121, v300, v319, v338, v357, v376, v395, v414, v433, v452, v471, or__x_125, v165, insts_166, ins_167, refs_168, doseq_seq__28653_169, doseq_loop__28654_170, label_171, i_172, n_173, op_174, r_175, or__x_176, v314, v333, v352, v371, v390, v409, v428, v447, v466, v485, insts_126, ins_127, refs_128, doseq_seq__28653_129, doseq_loop__28654_130, label_131, i_132, n_133, op_134, r_135, or__x_136, v310, v329, v348, v367, v386, v405, v424, v443, v462, v481, insts_137, ins_138, refs_139, doseq_seq__28653_140, doseq_loop__28654_141, label_142, i_143, n_144, op_145, r_146, or__x_147, v315, v334, v353, v372, v391, v410, v429, v448, v467, v486, v150, v152, insts_153, ins_154, refs_155, doseq_seq__28653_156, doseq_loop__28654_157, label_158, i_159, n_160, op_161, r_162, or__x_163, v303, v322, v341, v360, v379, v398, v417, v436, v455, v474
	n_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	i_4 = 0
	label_5 = arg1
	insts_6 = arg0
	n_7 = n_3
	v306 = 0
	v325 = 1
	v344 = vm.Keyword("invalid")
	v363 = 0
	v382 = "validate after "
	v401 = ": Inst #"
	v420 = " (op="
	v439 = ")"
	v458 = " has out-of-range ref "
	v477 = " (insts count="
	goto b1
b1:
	;
	v18 = rt.LtValue(vm.Int(i_4), n_7)
	if v18 {
		i_10 = i_4
		label_11 = label_5
		insts_12 = insts_6
		n_13 = n_7
		v305 = v306
		v324 = v325
		v343 = v344
		v362 = v363
		v381 = v382
		v400 = v401
		v419 = v420
		v438 = v439
		v457 = v458
		v476 = v477
		goto b2
	} else {
		i_14 = i_4
		label_15 = label_5
		insts_16 = insts_6
		n_17 = n_7
		v318 = v306
		v337 = v325
		v356 = v344
		v375 = v363
		v394 = v382
		v413 = v401
		v432 = v420
		v451 = v439
		v470 = v458
		v489 = v477
		goto b3
	}
b2:
	;
	ins_21, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{insts_12, vm.Int(i_10)})
	if callErr != nil {
		return nil, callErr
	}
	op_25, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{ins_21, vm.Int(v305)})
	if callErr != nil {
		return nil, callErr
	}
	refs_29, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{ins_21, vm.Int(v324)})
	if callErr != nil {
		return nil, callErr
	}
	v47, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not=").Deref(), []vm.Value{op_25, v343})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v47) {
		i_30 = i_10
		label_31 = label_11
		insts_32 = insts_12
		n_33 = n_13
		ins_34 = ins_21
		op_35 = op_25
		refs_36 = refs_29
		v312 = v305
		v331 = v324
		v350 = v343
		v369 = v362
		v388 = v381
		v407 = v400
		v426 = v419
		v445 = v438
		v464 = v457
		v483 = v476
		goto b5
	} else {
		i_37 = i_10
		label_38 = label_11
		insts_39 = insts_12
		n_40 = n_13
		ins_41 = ins_21
		op_42 = op_25
		refs_43 = refs_29
		v309 = v305
		v328 = v324
		v347 = v343
		v366 = v362
		v385 = v381
		v404 = v400
		v423 = v419
		v442 = v438
		v461 = v457
		v480 = v476
		goto b6
	}
b3:
	;
	v256 = vm.NIL
	i_257 = i_14
	label_258 = label_15
	insts_259 = insts_16
	n_260 = n_17
	goto b4
b4:
	;
	return v256, nil
b5:
	;
	doseq_seq__28653_50, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{refs_36})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__28654_51 = doseq_seq__28653_50
	label_52 = label_31
	i_53 = i_30
	n_54 = n_33
	op_55 = op_35
	v304 = v312
	v323 = v331
	v342 = v350
	v361 = v369
	v380 = v388
	v399 = v407
	v418 = v426
	v437 = v445
	v456 = v464
	v475 = v483
	goto b8
b6:
	;
	v244 = vm.NIL
	i_245 = i_37
	label_246 = label_38
	insts_247 = insts_39
	n_248 = n_40
	ins_249 = ins_41
	op_250 = op_42
	refs_251 = refs_43
	v317 = v309
	v336 = v328
	v355 = v347
	v374 = v366
	v393 = v385
	v412 = v404
	v431 = v423
	v450 = v442
	v469 = v461
	v488 = v480
	goto b7
b7:
	;
	v252 = i_245 + 1
	i_4 = v252
	label_5 = label_246
	insts_6 = insts_247
	n_7 = n_248
	v306 = v317
	v325 = v336
	v344 = v355
	v363 = v374
	v382 = v393
	v401 = v412
	v420 = v431
	v439 = v450
	v458 = v469
	v477 = v488
	goto b1
b8:
	;
	if vm.IsTruthy(doseq_loop__28654_51) {
		insts_57 = insts_32
		ins_58 = ins_34
		refs_59 = refs_36
		doseq_seq__28653_60 = doseq_seq__28653_50
		doseq_loop__28654_61 = doseq_loop__28654_51
		label_62 = label_52
		i_63 = i_53
		n_64 = n_54
		op_65 = op_55
		v313 = v304
		v332 = v323
		v351 = v342
		v370 = v361
		v389 = v380
		v408 = v399
		v427 = v418
		v446 = v437
		v465 = v456
		v484 = v475
		goto b9
	} else {
		insts_66 = insts_32
		ins_67 = ins_34
		refs_68 = refs_36
		doseq_seq__28653_69 = doseq_seq__28653_50
		doseq_loop__28654_70 = doseq_loop__28654_51
		label_71 = label_52
		i_72 = i_53
		n_73 = n_54
		op_74 = op_55
		v308 = v304
		v327 = v323
		v346 = v342
		v365 = v361
		v384 = v380
		v403 = v399
		v422 = v418
		v441 = v437
		v460 = v456
		v479 = v475
		goto b10
	}
b9:
	;
	r_77, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{doseq_loop__28654_61})
	if callErr != nil {
		return nil, callErr
	}
	or__x_99, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{r_77})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(or__x_99) {
		insts_100 = insts_57
		ins_101 = ins_58
		refs_102 = refs_59
		doseq_seq__28653_103 = doseq_seq__28653_60
		doseq_loop__28654_104 = doseq_loop__28654_61
		label_105 = label_62
		i_106 = i_63
		n_107 = n_64
		op_108 = op_65
		r_109 = r_77
		or__x_110 = or__x_99
		v307 = v313
		v326 = v332
		v345 = v351
		v364 = v370
		v383 = v389
		v402 = v408
		v421 = v427
		v440 = v446
		v459 = v465
		v478 = v484
		goto b15
	} else {
		insts_111 = insts_57
		ins_112 = ins_58
		refs_113 = refs_59
		doseq_seq__28653_114 = doseq_seq__28653_60
		doseq_loop__28654_115 = doseq_loop__28654_61
		label_116 = label_62
		i_117 = i_63
		n_118 = n_64
		op_119 = op_65
		r_120 = r_77
		or__x_121 = or__x_99
		v300 = v313
		v319 = v332
		v338 = v351
		v357 = v370
		v376 = v389
		v395 = v408
		v414 = v427
		v433 = v446
		v452 = v465
		v471 = v484
		goto b16
	}
b10:
	;
	v231 = vm.NIL
	insts_232 = insts_66
	ins_233 = ins_67
	refs_234 = refs_68
	doseq_seq__28653_235 = doseq_seq__28653_69
	doseq_loop__28654_236 = doseq_loop__28654_70
	label_237 = label_71
	i_238 = i_72
	n_239 = n_73
	op_240 = op_74
	v301 = v308
	v320 = v327
	v339 = v346
	v358 = v365
	v377 = v384
	v396 = v403
	v415 = v422
	v434 = v441
	v453 = v460
	v472 = v479
	goto b11
b11:
	;
	v244 = v231
	i_245 = i_238
	label_246 = label_237
	insts_247 = insts_232
	n_248 = n_239
	ins_249 = ins_233
	op_250 = op_240
	refs_251 = refs_234
	v317 = v301
	v336 = v320
	v355 = v339
	v374 = v358
	v393 = v377
	v412 = v396
	v431 = v415
	v450 = v434
	v469 = v453
	v488 = v472
	goto b7
b12:
	;
	arg__28718_193, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(v387), label_83, vm.String(v406), vm.Int(i_84), vm.String(v425), op_86, vm.String(v444), vm.String(v463), r_87, vm.String(v482), n_85, vm.String(v444)})
	if callErr != nil {
		return nil, callErr
	}
	arg__28745_210, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(v387), label_83, vm.String(v406), vm.Int(i_84), vm.String(v425), op_86, vm.String(v444), vm.String(v463), r_87, vm.String(v482), n_85, vm.String(v444)})
	if callErr != nil {
		return nil, callErr
	}
	v211, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "throw").Deref(), []vm.Value{arg__28745_210})
	if callErr != nil {
		return nil, callErr
	}
	v215 = v211
	insts_216 = insts_78
	ins_217 = ins_79
	refs_218 = refs_80
	doseq_seq__28653_219 = doseq_seq__28653_81
	doseq_loop__28654_220 = doseq_loop__28654_82
	label_221 = label_83
	i_222 = i_84
	n_223 = n_85
	op_224 = op_86
	r_225 = r_87
	v316 = v311
	v335 = v330
	v354 = v349
	v373 = v368
	v392 = v387
	v411 = v406
	v430 = v425
	v449 = v444
	v468 = v463
	v487 = v482
	goto b14
b13:
	;
	v215 = vm.NIL
	insts_216 = insts_88
	ins_217 = ins_89
	refs_218 = refs_90
	doseq_seq__28653_219 = doseq_seq__28653_91
	doseq_loop__28654_220 = doseq_loop__28654_92
	label_221 = label_93
	i_222 = i_94
	n_223 = n_95
	op_224 = op_96
	r_225 = r_97
	v316 = v302
	v335 = v321
	v354 = v340
	v373 = v359
	v392 = v378
	v411 = v397
	v430 = v416
	v449 = v435
	v468 = v454
	v487 = v473
	goto b14
b14:
	;
	v227, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "next").Deref(), []vm.Value{doseq_loop__28654_220})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__28654_51 = v227
	label_52 = label_221
	i_53 = i_222
	n_54 = n_223
	op_55 = op_224
	v304 = v316
	v323 = v335
	v342 = v354
	v361 = v373
	v380 = v392
	v399 = v411
	v418 = v430
	v437 = v449
	v456 = v468
	v475 = v487
	goto b8
b15:
	;
	v165 = or__x_110
	insts_166 = insts_100
	ins_167 = ins_101
	refs_168 = refs_102
	doseq_seq__28653_169 = doseq_seq__28653_103
	doseq_loop__28654_170 = doseq_loop__28654_104
	label_171 = label_105
	i_172 = i_106
	n_173 = n_107
	op_174 = op_108
	r_175 = r_109
	or__x_176 = or__x_110
	v314 = v307
	v333 = v326
	v352 = v345
	v371 = v364
	v390 = v383
	v409 = v402
	v428 = v421
	v447 = v440
	v466 = v459
	v485 = v478
	goto b17
b16:
	;
	or__x_125 = rt.LtValue(r_120, vm.Int(v357))
	if or__x_125 {
		insts_126 = insts_111
		ins_127 = ins_112
		refs_128 = refs_113
		doseq_seq__28653_129 = doseq_seq__28653_114
		doseq_loop__28654_130 = doseq_loop__28654_115
		label_131 = label_116
		i_132 = i_117
		n_133 = n_118
		op_134 = op_119
		r_135 = r_120
		or__x_136 = or__x_125
		v310 = v300
		v329 = v319
		v348 = v338
		v367 = v357
		v386 = v376
		v405 = v395
		v424 = v414
		v443 = v433
		v462 = v452
		v481 = v471
		goto b18
	} else {
		insts_137 = insts_111
		ins_138 = ins_112
		refs_139 = refs_113
		doseq_seq__28653_140 = doseq_seq__28653_114
		doseq_loop__28654_141 = doseq_loop__28654_115
		label_142 = label_116
		i_143 = i_117
		n_144 = n_118
		op_145 = op_119
		r_146 = r_120
		or__x_147 = or__x_125
		v315 = v300
		v334 = v319
		v353 = v338
		v372 = v357
		v391 = v376
		v410 = v395
		v429 = v414
		v448 = v433
		v467 = v452
		v486 = v471
		goto b19
	}
b17:
	;
	if vm.IsTruthy(v165) {
		insts_78 = insts_166
		ins_79 = ins_167
		refs_80 = refs_168
		doseq_seq__28653_81 = doseq_seq__28653_169
		doseq_loop__28654_82 = doseq_loop__28654_170
		label_83 = label_171
		i_84 = i_172
		n_85 = n_173
		op_86 = op_174
		r_87 = r_175
		v311 = v314
		v330 = v333
		v349 = v352
		v368 = v371
		v387 = v390
		v406 = v409
		v425 = v428
		v444 = v447
		v463 = v466
		v482 = v485
		goto b12
	} else {
		insts_88 = insts_166
		ins_89 = ins_167
		refs_90 = refs_168
		doseq_seq__28653_91 = doseq_seq__28653_169
		doseq_loop__28654_92 = doseq_loop__28654_170
		label_93 = label_171
		i_94 = i_172
		n_95 = n_173
		op_96 = op_174
		r_97 = r_175
		v302 = v314
		v321 = v333
		v340 = v352
		v359 = v371
		v378 = v390
		v397 = v409
		v416 = v428
		v435 = v447
		v454 = v466
		v473 = v485
		goto b13
	}
b18:
	;
	v152 = or__x_136
	insts_153 = insts_126
	ins_154 = ins_127
	refs_155 = refs_128
	doseq_seq__28653_156 = doseq_seq__28653_129
	doseq_loop__28654_157 = doseq_loop__28654_130
	label_158 = label_131
	i_159 = i_132
	n_160 = n_133
	op_161 = op_134
	r_162 = r_135
	or__x_163 = vm.Boolean(or__x_136)
	v303 = v310
	v322 = v329
	v341 = v348
	v360 = v367
	v379 = v386
	v398 = v405
	v417 = v424
	v436 = v443
	v455 = v462
	v474 = v481
	goto b20
b19:
	;
	v150 = rt.GeValue(r_146, n_144)
	v152 = v150
	insts_153 = insts_137
	ins_154 = ins_138
	refs_155 = refs_139
	doseq_seq__28653_156 = doseq_seq__28653_140
	doseq_loop__28654_157 = doseq_loop__28654_141
	label_158 = label_142
	i_159 = i_143
	n_160 = n_144
	op_161 = op_145
	r_162 = r_146
	or__x_163 = vm.Boolean(or__x_147)
	v303 = v315
	v322 = v334
	v341 = v353
	v360 = v372
	v379 = v391
	v398 = v410
	v417 = v429
	v436 = v448
	v455 = v467
	v474 = v486
	goto b20
b20:
	;
	v165 = vm.Boolean(v152)
	insts_166 = insts_153
	ins_167 = ins_154
	refs_168 = refs_155
	doseq_seq__28653_169 = doseq_seq__28653_156
	doseq_loop__28654_170 = doseq_loop__28654_157
	label_171 = label_158
	i_172 = i_159
	n_173 = n_160
	op_174 = op_161
	r_175 = r_162
	or__x_176 = or__x_121
	v314 = v303
	v333 = v322
	v352 = v341
	v371 = v360
	v390 = v379
	v409 = v398
	v428 = v417
	v447 = v436
	v466 = v455
	v485 = v474
	goto b17
}
func check_branch_arg_arity_for_target_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value, arg3 vm.Value) (vm.Value, error) {
	var target_6 vm.Value
	var args_8 vm.Value
	var tgt_blk_10 vm.Value
	var params_12 vm.Value
	var arg__28766_30 vm.Value
	var arg__28770_32 vm.Value
	var v33 bool
	var label_13 vm.Value
	var src_block_14 vm.Value
	var target_bt_15 vm.Value
	var blocks_16 vm.Value
	var target_17 vm.Value
	var args_18 vm.Value
	var tgt_blk_19 vm.Value
	var params_20 vm.Value
	var label_21 vm.Value
	var src_block_22 vm.Value
	var target_bt_23 vm.Value
	var blocks_24 vm.Value
	var target_25 vm.Value
	var args_26 vm.Value
	var tgt_blk_27 vm.Value
	var params_28 vm.Value
	var arg__28781_42 vm.Value
	var arg__28788_46 vm.Value
	var arg__28801_54 vm.Value
	var arg__28808_58 vm.Value
	var arg__28810_60 vm.Value
	var arg__28822_67 vm.Value
	var arg__28829_71 vm.Value
	var arg__28842_79 vm.Value
	var arg__28849_83 vm.Value
	var arg__28851_85 vm.Value
	var v86 vm.Value
	var v88 vm.Value
	var label_89 vm.Value
	var src_block_90 vm.Value
	var target_bt_91 vm.Value
	var blocks_92 vm.Value
	var target_93 vm.Value
	var args_94 vm.Value
	var tgt_blk_95 vm.Value
	var params_96 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = target_6, args_8, tgt_blk_10, params_12, arg__28766_30, arg__28770_32, v33, label_13, src_block_14, target_bt_15, blocks_16, target_17, args_18, tgt_blk_19, params_20, label_21, src_block_22, target_bt_23, blocks_24, target_25, args_26, tgt_blk_27, params_28, arg__28781_42, arg__28788_46, arg__28801_54, arg__28808_58, arg__28810_60, arg__28822_67, arg__28829_71, arg__28842_79, arg__28849_83, arg__28851_85, v86, v88, label_89, src_block_90, target_bt_91, blocks_92, target_93, args_94, tgt_blk_95, params_96
	target_6, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-target").Deref(), []vm.Value{arg2})
	if callErr != nil {
		return nil, callErr
	}
	args_8, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-args").Deref(), []vm.Value{arg2})
	if callErr != nil {
		return nil, callErr
	}
	tgt_blk_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg3, target_6})
	if callErr != nil {
		return nil, callErr
	}
	params_12, callErr = rt.InvokeValue(vm.Keyword("params"), []vm.Value{tgt_blk_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__28766_30, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{args_8})
	if callErr != nil {
		return nil, callErr
	}
	arg__28770_32, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{params_12})
	if callErr != nil {
		return nil, callErr
	}
	v33 = arg__28766_30 == arg__28770_32
	if v33 {
		label_13 = arg0
		src_block_14 = arg1
		target_bt_15 = arg2
		blocks_16 = arg3
		target_17 = target_6
		args_18 = args_8
		tgt_blk_19 = tgt_blk_10
		params_20 = params_12
		goto b1
	} else {
		label_21 = arg0
		src_block_22 = arg1
		target_bt_23 = arg2
		blocks_24 = arg3
		target_25 = target_6
		args_26 = args_8
		tgt_blk_27 = tgt_blk_10
		params_28 = params_12
		goto b2
	}
b1:
	;
	v88 = vm.NIL
	label_89 = label_13
	src_block_90 = src_block_14
	target_bt_91 = target_bt_15
	blocks_92 = blocks_16
	target_93 = target_17
	args_94 = args_18
	tgt_blk_95 = tgt_blk_19
	params_96 = params_20
	goto b3
b2:
	;
	arg__28781_42, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{args_26})
	if callErr != nil {
		return nil, callErr
	}
	arg__28788_46, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{params_28})
	if callErr != nil {
		return nil, callErr
	}
	arg__28801_54, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{args_26})
	if callErr != nil {
		return nil, callErr
	}
	arg__28808_58, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{params_28})
	if callErr != nil {
		return nil, callErr
	}
	arg__28810_60, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("validate after "), label_21, vm.String(": block #"), src_block_22, vm.String(" branch to b"), target_25, vm.String(" passes "), arg__28801_54, vm.String(" args but b"), target_25, vm.String(" has "), arg__28808_58, vm.String(" params")})
	if callErr != nil {
		return nil, callErr
	}
	arg__28822_67, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{args_26})
	if callErr != nil {
		return nil, callErr
	}
	arg__28829_71, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{params_28})
	if callErr != nil {
		return nil, callErr
	}
	arg__28842_79, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{args_26})
	if callErr != nil {
		return nil, callErr
	}
	arg__28849_83, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{params_28})
	if callErr != nil {
		return nil, callErr
	}
	arg__28851_85, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("validate after "), label_21, vm.String(": block #"), src_block_22, vm.String(" branch to b"), target_25, vm.String(" passes "), arg__28842_79, vm.String(" args but b"), target_25, vm.String(" has "), arg__28849_83, vm.String(" params")})
	if callErr != nil {
		return nil, callErr
	}
	v86, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "throw").Deref(), []vm.Value{arg__28851_85})
	if callErr != nil {
		return nil, callErr
	}
	v88 = v86
	label_89 = label_21
	src_block_90 = src_block_22
	target_bt_91 = target_bt_23
	blocks_92 = blocks_24
	target_93 = target_25
	args_94 = args_26
	tgt_blk_95 = tgt_blk_27
	params_96 = params_28
	goto b3
b3:
	;
	return v88, nil
}
func check_blocks_terminated_BANG_(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__28856_4 vm.Value
	var blocks_5 vm.Value
	var arg__28861_8 vm.Value
	var insts_9 vm.Value
	var i_10 int
	var label_11 vm.Value
	var insts_12 vm.Value
	var blocks_13 vm.Value
	var v388 vm.Value
	var v409 vm.Value
	var v430 int
	var v451 string
	var v472 string
	var v493 string
	var v514 int
	var v535 string
	var v556 string
	var v577 string
	var v598 string
	var v619 string
	var v640 int
	var v661 string
	var v682 string
	var v703 string
	var v724 string
	var arg__28866_27 vm.Value
	var v28 bool
	var f_16 vm.Value
	var i_17 int
	var label_18 vm.Value
	var insts_19 vm.Value
	var blocks_20 vm.Value
	var v387 vm.Value
	var v408 vm.Value
	var v429 int
	var v450 string
	var v471 string
	var v492 string
	var v513 int
	var v534 string
	var v555 string
	var v576 string
	var v597 string
	var v618 string
	var v639 int
	var v660 string
	var v681 string
	var v702 string
	var v723 string
	var blk_31 vm.Value
	var preds_33 vm.Value
	var term_id_35 vm.Value
	var or__x_53 bool
	var f_21 vm.Value
	var i_22 int
	var label_23 vm.Value
	var insts_24 vm.Value
	var blocks_25 vm.Value
	var v402 vm.Value
	var v423 vm.Value
	var v444 int
	var v465 string
	var v486 string
	var v507 string
	var v528 int
	var v549 string
	var v570 string
	var v591 string
	var v612 string
	var v633 string
	var v654 int
	var v675 string
	var v696 string
	var v717 string
	var v738 string
	var v317 vm.Value
	var f_318 vm.Value
	var i_319 int
	var label_320 vm.Value
	var insts_321 vm.Value
	var blocks_322 vm.Value
	var f_36 vm.Value
	var i_37 int
	var label_38 vm.Value
	var insts_39 vm.Value
	var blocks_40 vm.Value
	var blk_41 vm.Value
	var preds_42 vm.Value
	var term_id_43 vm.Value
	var v395 vm.Value
	var v416 vm.Value
	var v437 int
	var v458 string
	var v479 string
	var v500 string
	var v521 int
	var v542 string
	var v563 string
	var v584 string
	var v605 string
	var v626 string
	var v647 int
	var v668 string
	var v689 string
	var v710 string
	var v731 string
	var v105 vm.Value
	var f_44 vm.Value
	var i_45 int
	var label_46 vm.Value
	var insts_47 vm.Value
	var blocks_48 vm.Value
	var blk_49 vm.Value
	var preds_50 vm.Value
	var term_id_51 vm.Value
	var v391 vm.Value
	var v412 vm.Value
	var v433 int
	var v454 string
	var v475 string
	var v496 string
	var v517 int
	var v538 string
	var v559 string
	var v580 string
	var v601 string
	var v622 string
	var v643 int
	var v664 string
	var v685 string
	var v706 string
	var v727 string
	var v304 vm.Value
	var f_305 vm.Value
	var i_306 int
	var label_307 vm.Value
	var insts_308 vm.Value
	var blocks_309 vm.Value
	var blk_310 vm.Value
	var preds_311 vm.Value
	var term_id_312 vm.Value
	var v401 vm.Value
	var v422 vm.Value
	var v443 int
	var v464 string
	var v485 string
	var v506 string
	var v527 int
	var v548 string
	var v569 string
	var v590 string
	var v611 string
	var v632 string
	var v653 int
	var v674 string
	var v695 string
	var v716 string
	var v737 string
	var v313 int
	var f_54 vm.Value
	var i_55 int
	var label_56 vm.Value
	var insts_57 vm.Value
	var blocks_58 vm.Value
	var blk_59 vm.Value
	var preds_60 vm.Value
	var term_id_61 vm.Value
	var or__x_62 bool
	var v386 vm.Value
	var v407 vm.Value
	var v428 int
	var v449 string
	var v470 string
	var v491 string
	var v512 int
	var v533 string
	var v554 string
	var v575 string
	var v596 string
	var v617 string
	var v638 int
	var v659 string
	var v680 string
	var v701 string
	var v722 string
	var f_63 vm.Value
	var i_64 int
	var label_65 vm.Value
	var insts_66 vm.Value
	var blocks_67 vm.Value
	var blk_68 vm.Value
	var preds_69 vm.Value
	var term_id_70 vm.Value
	var or__x_71 bool
	var v396 vm.Value
	var v417 vm.Value
	var v438 int
	var v459 string
	var v480 string
	var v501 string
	var v522 int
	var v543 string
	var v564 string
	var v585 string
	var v606 string
	var v627 string
	var v648 int
	var v669 string
	var v690 string
	var v711 string
	var v732 string
	var v75 vm.Value
	var v77 vm.Value
	var f_78 vm.Value
	var i_79 int
	var label_80 vm.Value
	var insts_81 vm.Value
	var blocks_82 vm.Value
	var blk_83 vm.Value
	var preds_84 vm.Value
	var term_id_85 vm.Value
	var or__x_86 vm.Value
	var v390 vm.Value
	var v411 vm.Value
	var v432 int
	var v453 string
	var v474 string
	var v495 string
	var v516 int
	var v537 string
	var v558 string
	var v579 string
	var v600 string
	var v621 string
	var v642 int
	var v663 string
	var v684 string
	var v705 string
	var v726 string
	var f_88 vm.Value
	var i_89 int
	var label_90 vm.Value
	var insts_91 vm.Value
	var blocks_92 vm.Value
	var blk_93 vm.Value
	var preds_94 vm.Value
	var term_id_95 vm.Value
	var v383 vm.Value
	var v404 vm.Value
	var v425 int
	var v446 string
	var v467 string
	var v488 string
	var v509 int
	var v530 string
	var v551 string
	var v572 string
	var v593 string
	var v614 string
	var v635 int
	var v656 string
	var v677 string
	var v698 string
	var v719 string
	var arg__28895_114 vm.Value
	var arg__28908_123 vm.Value
	var v124 vm.Value
	var f_96 vm.Value
	var i_97 int
	var label_98 vm.Value
	var insts_99 vm.Value
	var blocks_100 vm.Value
	var blk_101 vm.Value
	var preds_102 vm.Value
	var term_id_103 vm.Value
	var v394 vm.Value
	var v415 vm.Value
	var v436 int
	var v457 string
	var v478 string
	var v499 string
	var v520 int
	var v541 string
	var v562 string
	var v583 string
	var v604 string
	var v625 string
	var v646 int
	var v667 string
	var v688 string
	var v709 string
	var v730 string
	var v128 vm.Value
	var f_129 vm.Value
	var i_130 int
	var label_131 vm.Value
	var insts_132 vm.Value
	var blocks_133 vm.Value
	var blk_134 vm.Value
	var preds_135 vm.Value
	var term_id_136 vm.Value
	var v384 vm.Value
	var v405 vm.Value
	var v426 int
	var v447 string
	var v468 string
	var v489 string
	var v510 int
	var v531 string
	var v552 string
	var v573 string
	var v594 string
	var v615 string
	var v636 int
	var v657 string
	var v678 string
	var v699 string
	var v720 string
	var or__x_154 bool
	var f_137 vm.Value
	var i_138 int
	var label_139 vm.Value
	var insts_140 vm.Value
	var blocks_141 vm.Value
	var blk_142 vm.Value
	var preds_143 vm.Value
	var term_id_144 vm.Value
	var v400 vm.Value
	var v421 vm.Value
	var v442 int
	var v463 string
	var v484 string
	var v505 string
	var v526 int
	var v547 string
	var v568 string
	var v589 string
	var v610 string
	var v631 string
	var v652 int
	var v673 string
	var v694 string
	var v715 string
	var v736 string
	var arg__28926_195 vm.Value
	var arg__28939_203 vm.Value
	var arg__28941_205 vm.Value
	var arg__28953_212 vm.Value
	var arg__28966_220 vm.Value
	var arg__28968_222 vm.Value
	var v223 vm.Value
	var f_145 vm.Value
	var i_146 int
	var label_147 vm.Value
	var insts_148 vm.Value
	var blocks_149 vm.Value
	var blk_150 vm.Value
	var preds_151 vm.Value
	var term_id_152 vm.Value
	var v389 vm.Value
	var v410 vm.Value
	var v431 int
	var v452 string
	var v473 string
	var v494 string
	var v515 int
	var v536 string
	var v557 string
	var v578 string
	var v599 string
	var v620 string
	var v641 int
	var v662 string
	var v683 string
	var v704 string
	var v725 string
	var v227 vm.Value
	var f_228 vm.Value
	var i_229 int
	var label_230 vm.Value
	var insts_231 vm.Value
	var blocks_232 vm.Value
	var blk_233 vm.Value
	var preds_234 vm.Value
	var term_id_235 vm.Value
	var v382 vm.Value
	var v403 vm.Value
	var v424 int
	var v445 string
	var v466 string
	var v487 string
	var v508 int
	var v529 string
	var v550 string
	var v571 string
	var v592 string
	var v613 string
	var v634 int
	var v655 string
	var v676 string
	var v697 string
	var v718 string
	var term_inst_237 vm.Value
	var term_op_241 vm.Value
	var v263 vm.Value
	var f_155 vm.Value
	var i_156 int
	var label_157 vm.Value
	var insts_158 vm.Value
	var blocks_159 vm.Value
	var blk_160 vm.Value
	var preds_161 vm.Value
	var term_id_162 vm.Value
	var or__x_163 bool
	var v397 vm.Value
	var v418 vm.Value
	var v439 int
	var v460 string
	var v481 string
	var v502 string
	var v523 int
	var v544 string
	var v565 string
	var v586 string
	var v607 string
	var v628 string
	var v649 int
	var v670 string
	var v691 string
	var v712 string
	var v733 string
	var f_164 vm.Value
	var i_165 int
	var label_166 vm.Value
	var insts_167 vm.Value
	var blocks_168 vm.Value
	var blk_169 vm.Value
	var preds_170 vm.Value
	var term_id_171 vm.Value
	var or__x_172 bool
	var v393 vm.Value
	var v414 vm.Value
	var v435 int
	var v456 string
	var v477 string
	var v498 string
	var v519 int
	var v540 string
	var v561 string
	var v582 string
	var v603 string
	var v624 string
	var v645 int
	var v666 string
	var v687 string
	var v708 string
	var v729 string
	var arg__28915_176 vm.Value
	var v177 bool
	var v179 bool
	var f_180 vm.Value
	var i_181 int
	var label_182 vm.Value
	var insts_183 vm.Value
	var blocks_184 vm.Value
	var blk_185 vm.Value
	var preds_186 vm.Value
	var term_id_187 vm.Value
	var or__x_188 vm.Value
	var v399 vm.Value
	var v420 vm.Value
	var v441 int
	var v462 string
	var v483 string
	var v504 string
	var v525 int
	var v546 string
	var v567 string
	var v588 string
	var v609 string
	var v630 string
	var v651 int
	var v672 string
	var v693 string
	var v714 string
	var v735 string
	var f_242 vm.Value
	var i_243 int
	var label_244 vm.Value
	var insts_245 vm.Value
	var blocks_246 vm.Value
	var blk_247 vm.Value
	var preds_248 vm.Value
	var term_id_249 vm.Value
	var term_inst_250 vm.Value
	var term_op_251 vm.Value
	var v385 vm.Value
	var v406 vm.Value
	var v427 int
	var v448 string
	var v469 string
	var v490 string
	var v511 int
	var v532 string
	var v553 string
	var v574 string
	var v595 string
	var v616 string
	var v637 int
	var v658 string
	var v679 string
	var v700 string
	var v721 string
	var f_252 vm.Value
	var i_253 int
	var label_254 vm.Value
	var insts_255 vm.Value
	var blocks_256 vm.Value
	var blk_257 vm.Value
	var preds_258 vm.Value
	var term_id_259 vm.Value
	var term_inst_260 vm.Value
	var term_op_261 vm.Value
	var v392 vm.Value
	var v413 vm.Value
	var v434 int
	var v455 string
	var v476 string
	var v497 string
	var v518 int
	var v539 string
	var v560 string
	var v581 string
	var v602 string
	var v623 string
	var v644 int
	var v665 string
	var v686 string
	var v707 string
	var v728 string
	var arg__28999_276 vm.Value
	var arg__29018_287 vm.Value
	var v288 vm.Value
	var v290 vm.Value
	var f_291 vm.Value
	var i_292 int
	var label_293 vm.Value
	var insts_294 vm.Value
	var blocks_295 vm.Value
	var blk_296 vm.Value
	var preds_297 vm.Value
	var term_id_298 vm.Value
	var term_inst_299 vm.Value
	var term_op_300 vm.Value
	var v398 vm.Value
	var v419 vm.Value
	var v440 int
	var v461 string
	var v482 string
	var v503 string
	var v524 int
	var v545 string
	var v566 string
	var v587 string
	var v608 string
	var v629 string
	var v650 int
	var v671 string
	var v692 string
	var v713 string
	var v734 string
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__28856_4, blocks_5, arg__28861_8, insts_9, i_10, label_11, insts_12, blocks_13, v388, v409, v430, v451, v472, v493, v514, v535, v556, v577, v598, v619, v640, v661, v682, v703, v724, arg__28866_27, v28, f_16, i_17, label_18, insts_19, blocks_20, v387, v408, v429, v450, v471, v492, v513, v534, v555, v576, v597, v618, v639, v660, v681, v702, v723, blk_31, preds_33, term_id_35, or__x_53, f_21, i_22, label_23, insts_24, blocks_25, v402, v423, v444, v465, v486, v507, v528, v549, v570, v591, v612, v633, v654, v675, v696, v717, v738, v317, f_318, i_319, label_320, insts_321, blocks_322, f_36, i_37, label_38, insts_39, blocks_40, blk_41, preds_42, term_id_43, v395, v416, v437, v458, v479, v500, v521, v542, v563, v584, v605, v626, v647, v668, v689, v710, v731, v105, f_44, i_45, label_46, insts_47, blocks_48, blk_49, preds_50, term_id_51, v391, v412, v433, v454, v475, v496, v517, v538, v559, v580, v601, v622, v643, v664, v685, v706, v727, v304, f_305, i_306, label_307, insts_308, blocks_309, blk_310, preds_311, term_id_312, v401, v422, v443, v464, v485, v506, v527, v548, v569, v590, v611, v632, v653, v674, v695, v716, v737, v313, f_54, i_55, label_56, insts_57, blocks_58, blk_59, preds_60, term_id_61, or__x_62, v386, v407, v428, v449, v470, v491, v512, v533, v554, v575, v596, v617, v638, v659, v680, v701, v722, f_63, i_64, label_65, insts_66, blocks_67, blk_68, preds_69, term_id_70, or__x_71, v396, v417, v438, v459, v480, v501, v522, v543, v564, v585, v606, v627, v648, v669, v690, v711, v732, v75, v77, f_78, i_79, label_80, insts_81, blocks_82, blk_83, preds_84, term_id_85, or__x_86, v390, v411, v432, v453, v474, v495, v516, v537, v558, v579, v600, v621, v642, v663, v684, v705, v726, f_88, i_89, label_90, insts_91, blocks_92, blk_93, preds_94, term_id_95, v383, v404, v425, v446, v467, v488, v509, v530, v551, v572, v593, v614, v635, v656, v677, v698, v719, arg__28895_114, arg__28908_123, v124, f_96, i_97, label_98, insts_99, blocks_100, blk_101, preds_102, term_id_103, v394, v415, v436, v457, v478, v499, v520, v541, v562, v583, v604, v625, v646, v667, v688, v709, v730, v128, f_129, i_130, label_131, insts_132, blocks_133, blk_134, preds_135, term_id_136, v384, v405, v426, v447, v468, v489, v510, v531, v552, v573, v594, v615, v636, v657, v678, v699, v720, or__x_154, f_137, i_138, label_139, insts_140, blocks_141, blk_142, preds_143, term_id_144, v400, v421, v442, v463, v484, v505, v526, v547, v568, v589, v610, v631, v652, v673, v694, v715, v736, arg__28926_195, arg__28939_203, arg__28941_205, arg__28953_212, arg__28966_220, arg__28968_222, v223, f_145, i_146, label_147, insts_148, blocks_149, blk_150, preds_151, term_id_152, v389, v410, v431, v452, v473, v494, v515, v536, v557, v578, v599, v620, v641, v662, v683, v704, v725, v227, f_228, i_229, label_230, insts_231, blocks_232, blk_233, preds_234, term_id_235, v382, v403, v424, v445, v466, v487, v508, v529, v550, v571, v592, v613, v634, v655, v676, v697, v718, term_inst_237, term_op_241, v263, f_155, i_156, label_157, insts_158, blocks_159, blk_160, preds_161, term_id_162, or__x_163, v397, v418, v439, v460, v481, v502, v523, v544, v565, v586, v607, v628, v649, v670, v691, v712, v733, f_164, i_165, label_166, insts_167, blocks_168, blk_169, preds_170, term_id_171, or__x_172, v393, v414, v435, v456, v477, v498, v519, v540, v561, v582, v603, v624, v645, v666, v687, v708, v729, arg__28915_176, v177, v179, f_180, i_181, label_182, insts_183, blocks_184, blk_185, preds_186, term_id_187, or__x_188, v399, v420, v441, v462, v483, v504, v525, v546, v567, v588, v609, v630, v651, v672, v693, v714, v735, f_242, i_243, label_244, insts_245, blocks_246, blk_247, preds_248, term_id_249, term_inst_250, term_op_251, v385, v406, v427, v448, v469, v490, v511, v532, v553, v574, v595, v616, v637, v658, v679, v700, v721, f_252, i_253, label_254, insts_255, blocks_256, blk_257, preds_258, term_id_259, term_inst_260, term_op_261, v392, v413, v434, v455, v476, v497, v518, v539, v560, v581, v602, v623, v644, v665, v686, v707, v728, arg__28999_276, arg__29018_287, v288, v290, f_291, i_292, label_293, insts_294, blocks_295, blk_296, preds_297, term_id_298, term_inst_299, term_op_300, v398, v419, v440, v461, v482, v503, v524, v545, v566, v587, v608, v629, v650, v671, v692, v713, v734
	arg__28856_4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	blocks_5, callErr = rt.InvokeValue(vm.Keyword("blocks"), []vm.Value{arg__28856_4})
	if callErr != nil {
		return nil, callErr
	}
	arg__28861_8, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	insts_9, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__28861_8})
	if callErr != nil {
		return nil, callErr
	}
	i_10 = 0
	label_11 = arg1
	insts_12 = insts_9
	blocks_13 = blocks_5
	v388 = vm.Keyword("preds")
	v409 = vm.Keyword("term")
	v430 = 0
	v451 = "validate after "
	v472 = ": block #"
	v493 = " has no :term field"
	v514 = 0
	v535 = "validate after "
	v556 = ": block #"
	v577 = " :term="
	v598 = " out of insts range [0,"
	v619 = ")"
	v640 = 0
	v661 = "validate after "
	v682 = ": block #"
	v703 = " :term Inst #"
	v724 = " has non-terminator op "
	goto b1
b1:
	;
	arg__28866_27, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{blocks_13})
	if callErr != nil {
		return nil, callErr
	}
	v28 = rt.LtValue(vm.Int(i_10), arg__28866_27)
	if v28 {
		f_16 = arg0
		i_17 = i_10
		label_18 = label_11
		insts_19 = insts_12
		blocks_20 = blocks_13
		v387 = v388
		v408 = v409
		v429 = v430
		v450 = v451
		v471 = v472
		v492 = v493
		v513 = v514
		v534 = v535
		v555 = v556
		v576 = v577
		v597 = v598
		v618 = v619
		v639 = v640
		v660 = v661
		v681 = v682
		v702 = v703
		v723 = v724
		goto b2
	} else {
		f_21 = arg0
		i_22 = i_10
		label_23 = label_11
		insts_24 = insts_12
		blocks_25 = blocks_13
		v402 = v388
		v423 = v409
		v444 = v430
		v465 = v451
		v486 = v472
		v507 = v493
		v528 = v514
		v549 = v535
		v570 = v556
		v591 = v577
		v612 = v598
		v633 = v619
		v654 = v640
		v675 = v661
		v696 = v682
		v717 = v703
		v738 = v724
		goto b3
	}
b2:
	;
	blk_31, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{blocks_20, vm.Int(i_17)})
	if callErr != nil {
		return nil, callErr
	}
	preds_33, callErr = rt.InvokeValue(v387, []vm.Value{blk_31})
	if callErr != nil {
		return nil, callErr
	}
	term_id_35, callErr = rt.InvokeValue(v408, []vm.Value{blk_31})
	if callErr != nil {
		return nil, callErr
	}
	or__x_53 = i_17 == v429
	if or__x_53 {
		f_54 = f_16
		i_55 = i_17
		label_56 = label_18
		insts_57 = insts_19
		blocks_58 = blocks_20
		blk_59 = blk_31
		preds_60 = preds_33
		term_id_61 = term_id_35
		or__x_62 = or__x_53
		v386 = v387
		v407 = v408
		v428 = v429
		v449 = v450
		v470 = v471
		v491 = v492
		v512 = v513
		v533 = v534
		v554 = v555
		v575 = v576
		v596 = v597
		v617 = v618
		v638 = v639
		v659 = v660
		v680 = v681
		v701 = v702
		v722 = v723
		goto b8
	} else {
		f_63 = f_16
		i_64 = i_17
		label_65 = label_18
		insts_66 = insts_19
		blocks_67 = blocks_20
		blk_68 = blk_31
		preds_69 = preds_33
		term_id_70 = term_id_35
		or__x_71 = or__x_53
		v396 = v387
		v417 = v408
		v438 = v429
		v459 = v450
		v480 = v471
		v501 = v492
		v522 = v513
		v543 = v534
		v564 = v555
		v585 = v576
		v606 = v597
		v627 = v618
		v648 = v639
		v669 = v660
		v690 = v681
		v711 = v702
		v732 = v723
		goto b9
	}
b3:
	;
	v317 = vm.NIL
	f_318 = f_21
	i_319 = i_22
	label_320 = label_23
	insts_321 = insts_24
	blocks_322 = blocks_25
	goto b4
b4:
	;
	return v317, nil
b5:
	;
	v105, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{term_id_43})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v105) {
		f_88 = f_36
		i_89 = i_37
		label_90 = label_38
		insts_91 = insts_39
		blocks_92 = blocks_40
		blk_93 = blk_41
		preds_94 = preds_42
		term_id_95 = term_id_43
		v383 = v395
		v404 = v416
		v425 = v437
		v446 = v458
		v467 = v479
		v488 = v500
		v509 = v521
		v530 = v542
		v551 = v563
		v572 = v584
		v593 = v605
		v614 = v626
		v635 = v647
		v656 = v668
		v677 = v689
		v698 = v710
		v719 = v731
		goto b11
	} else {
		f_96 = f_36
		i_97 = i_37
		label_98 = label_38
		insts_99 = insts_39
		blocks_100 = blocks_40
		blk_101 = blk_41
		preds_102 = preds_42
		term_id_103 = term_id_43
		v394 = v395
		v415 = v416
		v436 = v437
		v457 = v458
		v478 = v479
		v499 = v500
		v520 = v521
		v541 = v542
		v562 = v563
		v583 = v584
		v604 = v605
		v625 = v626
		v646 = v647
		v667 = v668
		v688 = v689
		v709 = v710
		v730 = v731
		goto b12
	}
b6:
	;
	v304 = vm.NIL
	f_305 = f_44
	i_306 = i_45
	label_307 = label_46
	insts_308 = insts_47
	blocks_309 = blocks_48
	blk_310 = blk_49
	preds_311 = preds_50
	term_id_312 = term_id_51
	v401 = v391
	v422 = v412
	v443 = v433
	v464 = v454
	v485 = v475
	v506 = v496
	v527 = v517
	v548 = v538
	v569 = v559
	v590 = v580
	v611 = v601
	v632 = v622
	v653 = v643
	v674 = v664
	v695 = v685
	v716 = v706
	v737 = v727
	goto b7
b7:
	;
	v313 = i_306 + 1
	i_10 = v313
	label_11 = label_307
	insts_12 = insts_308
	blocks_13 = blocks_309
	v388 = v401
	v409 = v422
	v430 = v443
	v451 = v464
	v472 = v485
	v493 = v506
	v514 = v527
	v535 = v548
	v556 = v569
	v577 = v590
	v598 = v611
	v619 = v632
	v640 = v653
	v661 = v674
	v682 = v695
	v703 = v716
	v724 = v737
	goto b1
b8:
	;
	v77 = vm.Boolean(or__x_62)
	f_78 = f_54
	i_79 = i_55
	label_80 = label_56
	insts_81 = insts_57
	blocks_82 = blocks_58
	blk_83 = blk_59
	preds_84 = preds_60
	term_id_85 = term_id_61
	or__x_86 = vm.Boolean(or__x_62)
	v390 = v386
	v411 = v407
	v432 = v428
	v453 = v449
	v474 = v470
	v495 = v491
	v516 = v512
	v537 = v533
	v558 = v554
	v579 = v575
	v600 = v596
	v621 = v617
	v642 = v638
	v663 = v659
	v684 = v680
	v705 = v701
	v726 = v722
	goto b10
b9:
	;
	v75, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{preds_69})
	if callErr != nil {
		return nil, callErr
	}
	v77 = v75
	f_78 = f_63
	i_79 = i_64
	label_80 = label_65
	insts_81 = insts_66
	blocks_82 = blocks_67
	blk_83 = blk_68
	preds_84 = preds_69
	term_id_85 = term_id_70
	or__x_86 = vm.Boolean(or__x_71)
	v390 = v396
	v411 = v417
	v432 = v438
	v453 = v459
	v474 = v480
	v495 = v501
	v516 = v522
	v537 = v543
	v558 = v564
	v579 = v585
	v600 = v606
	v621 = v627
	v642 = v648
	v663 = v669
	v684 = v690
	v705 = v711
	v726 = v732
	goto b10
b10:
	;
	if vm.IsTruthy(v77) {
		f_36 = f_78
		i_37 = i_79
		label_38 = label_80
		insts_39 = insts_81
		blocks_40 = blocks_82
		blk_41 = blk_83
		preds_42 = preds_84
		term_id_43 = term_id_85
		v395 = v390
		v416 = v411
		v437 = v432
		v458 = v453
		v479 = v474
		v500 = v495
		v521 = v516
		v542 = v537
		v563 = v558
		v584 = v579
		v605 = v600
		v626 = v621
		v647 = v642
		v668 = v663
		v689 = v684
		v710 = v705
		v731 = v726
		goto b5
	} else {
		f_44 = f_78
		i_45 = i_79
		label_46 = label_80
		insts_47 = insts_81
		blocks_48 = blocks_82
		blk_49 = blk_83
		preds_50 = preds_84
		term_id_51 = term_id_85
		v391 = v390
		v412 = v411
		v433 = v432
		v454 = v453
		v475 = v474
		v496 = v495
		v517 = v516
		v538 = v537
		v559 = v558
		v580 = v579
		v601 = v600
		v622 = v621
		v643 = v642
		v664 = v663
		v685 = v684
		v706 = v705
		v727 = v726
		goto b6
	}
b11:
	;
	arg__28895_114, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(v446), label_90, vm.String(v467), vm.Int(i_89), vm.String(v488)})
	if callErr != nil {
		return nil, callErr
	}
	arg__28908_123, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(v446), label_90, vm.String(v467), vm.Int(i_89), vm.String(v488)})
	if callErr != nil {
		return nil, callErr
	}
	v124, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "throw").Deref(), []vm.Value{arg__28908_123})
	if callErr != nil {
		return nil, callErr
	}
	v128 = v124
	f_129 = f_88
	i_130 = i_89
	label_131 = label_90
	insts_132 = insts_91
	blocks_133 = blocks_92
	blk_134 = blk_93
	preds_135 = preds_94
	term_id_136 = term_id_95
	v384 = v383
	v405 = v404
	v426 = v425
	v447 = v446
	v468 = v467
	v489 = v488
	v510 = v509
	v531 = v530
	v552 = v551
	v573 = v572
	v594 = v593
	v615 = v614
	v636 = v635
	v657 = v656
	v678 = v677
	v699 = v698
	v720 = v719
	goto b13
b12:
	;
	v128 = vm.NIL
	f_129 = f_96
	i_130 = i_97
	label_131 = label_98
	insts_132 = insts_99
	blocks_133 = blocks_100
	blk_134 = blk_101
	preds_135 = preds_102
	term_id_136 = term_id_103
	v384 = v394
	v405 = v415
	v426 = v436
	v447 = v457
	v468 = v478
	v489 = v499
	v510 = v520
	v531 = v541
	v552 = v562
	v573 = v583
	v594 = v604
	v615 = v625
	v636 = v646
	v657 = v667
	v678 = v688
	v699 = v709
	v720 = v730
	goto b13
b13:
	;
	or__x_154 = rt.LtValue(term_id_136, vm.Int(v510))
	if or__x_154 {
		f_155 = f_129
		i_156 = i_130
		label_157 = label_131
		insts_158 = insts_132
		blocks_159 = blocks_133
		blk_160 = blk_134
		preds_161 = preds_135
		term_id_162 = term_id_136
		or__x_163 = or__x_154
		v397 = v384
		v418 = v405
		v439 = v426
		v460 = v447
		v481 = v468
		v502 = v489
		v523 = v510
		v544 = v531
		v565 = v552
		v586 = v573
		v607 = v594
		v628 = v615
		v649 = v636
		v670 = v657
		v691 = v678
		v712 = v699
		v733 = v720
		goto b17
	} else {
		f_164 = f_129
		i_165 = i_130
		label_166 = label_131
		insts_167 = insts_132
		blocks_168 = blocks_133
		blk_169 = blk_134
		preds_170 = preds_135
		term_id_171 = term_id_136
		or__x_172 = or__x_154
		v393 = v384
		v414 = v405
		v435 = v426
		v456 = v447
		v477 = v468
		v498 = v489
		v519 = v510
		v540 = v531
		v561 = v552
		v582 = v573
		v603 = v594
		v624 = v615
		v645 = v636
		v666 = v657
		v687 = v678
		v708 = v699
		v729 = v720
		goto b18
	}
b14:
	;
	arg__28926_195, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{insts_140})
	if callErr != nil {
		return nil, callErr
	}
	arg__28939_203, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{insts_140})
	if callErr != nil {
		return nil, callErr
	}
	arg__28941_205, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(v547), label_139, vm.String(v568), vm.Int(i_138), vm.String(v589), term_id_144, vm.String(v610), arg__28939_203, vm.String(v631)})
	if callErr != nil {
		return nil, callErr
	}
	arg__28953_212, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{insts_140})
	if callErr != nil {
		return nil, callErr
	}
	arg__28966_220, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{insts_140})
	if callErr != nil {
		return nil, callErr
	}
	arg__28968_222, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(v547), label_139, vm.String(v568), vm.Int(i_138), vm.String(v589), term_id_144, vm.String(v610), arg__28966_220, vm.String(v631)})
	if callErr != nil {
		return nil, callErr
	}
	v223, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "throw").Deref(), []vm.Value{arg__28968_222})
	if callErr != nil {
		return nil, callErr
	}
	v227 = v223
	f_228 = f_137
	i_229 = i_138
	label_230 = label_139
	insts_231 = insts_140
	blocks_232 = blocks_141
	blk_233 = blk_142
	preds_234 = preds_143
	term_id_235 = term_id_144
	v382 = v400
	v403 = v421
	v424 = v442
	v445 = v463
	v466 = v484
	v487 = v505
	v508 = v526
	v529 = v547
	v550 = v568
	v571 = v589
	v592 = v610
	v613 = v631
	v634 = v652
	v655 = v673
	v676 = v694
	v697 = v715
	v718 = v736
	goto b16
b15:
	;
	v227 = vm.NIL
	f_228 = f_145
	i_229 = i_146
	label_230 = label_147
	insts_231 = insts_148
	blocks_232 = blocks_149
	blk_233 = blk_150
	preds_234 = preds_151
	term_id_235 = term_id_152
	v382 = v389
	v403 = v410
	v424 = v431
	v445 = v452
	v466 = v473
	v487 = v494
	v508 = v515
	v529 = v536
	v550 = v557
	v571 = v578
	v592 = v599
	v613 = v620
	v634 = v641
	v655 = v662
	v676 = v683
	v697 = v704
	v718 = v725
	goto b16
b16:
	;
	term_inst_237, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{insts_231, term_id_235})
	if callErr != nil {
		return nil, callErr
	}
	term_op_241, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{term_inst_237, vm.Int(v634)})
	if callErr != nil {
		return nil, callErr
	}
	v263, callErr = rt.InvokeValue(rt.LookupVar("ir", "op-terminator?").Deref(), []vm.Value{term_op_241})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v263) {
		f_242 = f_228
		i_243 = i_229
		label_244 = label_230
		insts_245 = insts_231
		blocks_246 = blocks_232
		blk_247 = blk_233
		preds_248 = preds_234
		term_id_249 = term_id_235
		term_inst_250 = term_inst_237
		term_op_251 = term_op_241
		v385 = v382
		v406 = v403
		v427 = v424
		v448 = v445
		v469 = v466
		v490 = v487
		v511 = v508
		v532 = v529
		v553 = v550
		v574 = v571
		v595 = v592
		v616 = v613
		v637 = v634
		v658 = v655
		v679 = v676
		v700 = v697
		v721 = v718
		goto b20
	} else {
		f_252 = f_228
		i_253 = i_229
		label_254 = label_230
		insts_255 = insts_231
		blocks_256 = blocks_232
		blk_257 = blk_233
		preds_258 = preds_234
		term_id_259 = term_id_235
		term_inst_260 = term_inst_237
		term_op_261 = term_op_241
		v392 = v382
		v413 = v403
		v434 = v424
		v455 = v445
		v476 = v466
		v497 = v487
		v518 = v508
		v539 = v529
		v560 = v550
		v581 = v571
		v602 = v592
		v623 = v613
		v644 = v634
		v665 = v655
		v686 = v676
		v707 = v697
		v728 = v718
		goto b21
	}
b17:
	;
	v179 = or__x_163
	f_180 = f_155
	i_181 = i_156
	label_182 = label_157
	insts_183 = insts_158
	blocks_184 = blocks_159
	blk_185 = blk_160
	preds_186 = preds_161
	term_id_187 = term_id_162
	or__x_188 = vm.Boolean(or__x_163)
	v399 = v397
	v420 = v418
	v441 = v439
	v462 = v460
	v483 = v481
	v504 = v502
	v525 = v523
	v546 = v544
	v567 = v565
	v588 = v586
	v609 = v607
	v630 = v628
	v651 = v649
	v672 = v670
	v693 = v691
	v714 = v712
	v735 = v733
	goto b19
b18:
	;
	arg__28915_176, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{insts_167})
	if callErr != nil {
		return nil, callErr
	}
	v177 = rt.GeValue(term_id_171, arg__28915_176)
	v179 = v177
	f_180 = f_164
	i_181 = i_165
	label_182 = label_166
	insts_183 = insts_167
	blocks_184 = blocks_168
	blk_185 = blk_169
	preds_186 = preds_170
	term_id_187 = term_id_171
	or__x_188 = vm.Boolean(or__x_172)
	v399 = v393
	v420 = v414
	v441 = v435
	v462 = v456
	v483 = v477
	v504 = v498
	v525 = v519
	v546 = v540
	v567 = v561
	v588 = v582
	v609 = v603
	v630 = v624
	v651 = v645
	v672 = v666
	v693 = v687
	v714 = v708
	v735 = v729
	goto b19
b19:
	;
	if v179 {
		f_137 = f_180
		i_138 = i_181
		label_139 = label_182
		insts_140 = insts_183
		blocks_141 = blocks_184
		blk_142 = blk_185
		preds_143 = preds_186
		term_id_144 = term_id_187
		v400 = v399
		v421 = v420
		v442 = v441
		v463 = v462
		v484 = v483
		v505 = v504
		v526 = v525
		v547 = v546
		v568 = v567
		v589 = v588
		v610 = v609
		v631 = v630
		v652 = v651
		v673 = v672
		v694 = v693
		v715 = v714
		v736 = v735
		goto b14
	} else {
		f_145 = f_180
		i_146 = i_181
		label_147 = label_182
		insts_148 = insts_183
		blocks_149 = blocks_184
		blk_150 = blk_185
		preds_151 = preds_186
		term_id_152 = term_id_187
		v389 = v399
		v410 = v420
		v431 = v441
		v452 = v462
		v473 = v483
		v494 = v504
		v515 = v525
		v536 = v546
		v557 = v567
		v578 = v588
		v599 = v609
		v620 = v630
		v641 = v651
		v662 = v672
		v683 = v693
		v704 = v714
		v725 = v735
		goto b15
	}
b20:
	;
	v290 = vm.NIL
	f_291 = f_242
	i_292 = i_243
	label_293 = label_244
	insts_294 = insts_245
	blocks_295 = blocks_246
	blk_296 = blk_247
	preds_297 = preds_248
	term_id_298 = term_id_249
	term_inst_299 = term_inst_250
	term_op_300 = term_op_251
	v398 = v385
	v419 = v406
	v440 = v427
	v461 = v448
	v482 = v469
	v503 = v490
	v524 = v511
	v545 = v532
	v566 = v553
	v587 = v574
	v608 = v595
	v629 = v616
	v650 = v637
	v671 = v658
	v692 = v679
	v713 = v700
	v734 = v721
	goto b22
b21:
	;
	arg__28999_276, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(v665), label_254, vm.String(v686), vm.Int(i_253), vm.String(v707), term_id_259, vm.String(v728), term_op_261})
	if callErr != nil {
		return nil, callErr
	}
	arg__29018_287, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(v665), label_254, vm.String(v686), vm.Int(i_253), vm.String(v707), term_id_259, vm.String(v728), term_op_261})
	if callErr != nil {
		return nil, callErr
	}
	v288, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "throw").Deref(), []vm.Value{arg__29018_287})
	if callErr != nil {
		return nil, callErr
	}
	v290 = v288
	f_291 = f_252
	i_292 = i_253
	label_293 = label_254
	insts_294 = insts_255
	blocks_295 = blocks_256
	blk_296 = blk_257
	preds_297 = preds_258
	term_id_298 = term_id_259
	term_inst_299 = term_inst_260
	term_op_300 = term_op_261
	v398 = v392
	v419 = v413
	v440 = v434
	v461 = v455
	v482 = v476
	v503 = v497
	v524 = v518
	v545 = v539
	v566 = v560
	v587 = v581
	v608 = v602
	v629 = v623
	v650 = v644
	v671 = v665
	v692 = v686
	v713 = v707
	v734 = v728
	goto b22
b22:
	;
	v304 = v290
	f_305 = f_291
	i_306 = i_292
	label_307 = label_293
	insts_308 = insts_294
	blocks_309 = blocks_295
	blk_310 = blk_296
	preds_311 = preds_297
	term_id_312 = term_id_298
	v401 = v398
	v422 = v419
	v443 = v440
	v464 = v461
	v485 = v482
	v506 = v503
	v527 = v524
	v548 = v545
	v569 = v566
	v590 = v587
	v611 = v608
	v632 = v629
	v653 = v650
	v674 = v671
	v695 = v692
	v716 = v713
	v737 = v734
	goto b7
}
func check_branch_arg_arities_BANG_(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__29024_5 vm.Value
	var blocks_6 vm.Value
	var arg__29029_9 vm.Value
	var insts_10 vm.Value
	var i_11 int
	var label_12 vm.Value
	var insts_13 vm.Value
	var blocks_14 vm.Value
	var v182 vm.Value
	var v194 int
	var v206 int
	var v218 vm.Value
	var v230 vm.Value
	var arg__29034_28 vm.Value
	var v29 bool
	var f_17 vm.Value
	var i_18 int
	var label_19 vm.Value
	var insts_20 vm.Value
	var blocks_21 vm.Value
	var v181 vm.Value
	var v193 int
	var v205 int
	var v217 vm.Value
	var v229 vm.Value
	var blk_32 vm.Value
	var term_id_34 vm.Value
	var f_22 vm.Value
	var i_23 int
	var label_24 vm.Value
	var insts_25 vm.Value
	var blocks_26 vm.Value
	var v189 vm.Value
	var v201 int
	var v213 int
	var v225 vm.Value
	var v237 vm.Value
	var v162 vm.Value
	var f_163 vm.Value
	var i_164 int
	var label_165 vm.Value
	var insts_166 vm.Value
	var blocks_167 vm.Value
	var f_35 vm.Value
	var i_36 int
	var label_37 vm.Value
	var insts_38 vm.Value
	var blocks_39 vm.Value
	var blk_40 vm.Value
	var term_id_41 vm.Value
	var v186 vm.Value
	var v198 int
	var v210 int
	var v222 vm.Value
	var v234 vm.Value
	var term_51 vm.Value
	var op_55 vm.Value
	var aux_59 vm.Value
	var v81 bool
	var f_42 vm.Value
	var i_43 int
	var label_44 vm.Value
	var insts_45 vm.Value
	var blocks_46 vm.Value
	var blk_47 vm.Value
	var term_id_48 vm.Value
	var v184 vm.Value
	var v196 int
	var v208 int
	var v220 vm.Value
	var v232 vm.Value
	var v150 vm.Value
	var f_151 vm.Value
	var i_152 int
	var label_153 vm.Value
	var insts_154 vm.Value
	var blocks_155 vm.Value
	var blk_156 vm.Value
	var term_id_157 vm.Value
	var v188 vm.Value
	var v200 int
	var v212 int
	var v224 vm.Value
	var v236 vm.Value
	var v158 int
	var f_60 vm.Value
	var i_61 int
	var label_62 vm.Value
	var insts_63 vm.Value
	var blocks_64 vm.Value
	var blk_65 vm.Value
	var term_id_66 vm.Value
	var term_67 vm.Value
	var op_68 vm.Value
	var aux_69 vm.Value
	var v180 vm.Value
	var v192 int
	var v204 int
	var v216 vm.Value
	var v228 vm.Value
	var v84 vm.Value
	var f_70 vm.Value
	var i_71 int
	var label_72 vm.Value
	var insts_73 vm.Value
	var blocks_74 vm.Value
	var blk_75 vm.Value
	var term_id_76 vm.Value
	var term_77 vm.Value
	var op_78 vm.Value
	var aux_79 vm.Value
	var v187 vm.Value
	var v199 int
	var v211 int
	var v223 vm.Value
	var v235 vm.Value
	var v107 bool
	var v136 vm.Value
	var f_137 vm.Value
	var i_138 int
	var label_139 vm.Value
	var insts_140 vm.Value
	var blocks_141 vm.Value
	var blk_142 vm.Value
	var term_id_143 vm.Value
	var term_144 vm.Value
	var op_145 vm.Value
	var aux_146 vm.Value
	var v183 vm.Value
	var v195 int
	var v207 int
	var v219 vm.Value
	var v231 vm.Value
	var f_86 vm.Value
	var i_87 int
	var label_88 vm.Value
	var insts_89 vm.Value
	var blocks_90 vm.Value
	var blk_91 vm.Value
	var term_id_92 vm.Value
	var term_93 vm.Value
	var op_94 vm.Value
	var aux_95 vm.Value
	var v178 vm.Value
	var v190 int
	var v202 int
	var v214 vm.Value
	var v226 vm.Value
	var arg__29075_110 vm.Value
	var arg__29083_113 vm.Value
	var v114 vm.Value
	var arg__29090_116 vm.Value
	var arg__29098_119 vm.Value
	var v120 vm.Value
	var f_96 vm.Value
	var i_97 int
	var label_98 vm.Value
	var insts_99 vm.Value
	var blocks_100 vm.Value
	var blk_101 vm.Value
	var term_id_102 vm.Value
	var term_103 vm.Value
	var op_104 vm.Value
	var aux_105 vm.Value
	var v185 vm.Value
	var v197 int
	var v209 int
	var v221 vm.Value
	var v233 vm.Value
	var v124 vm.Value
	var f_125 vm.Value
	var i_126 int
	var label_127 vm.Value
	var insts_128 vm.Value
	var blocks_129 vm.Value
	var blk_130 vm.Value
	var term_id_131 vm.Value
	var term_132 vm.Value
	var op_133 vm.Value
	var aux_134 vm.Value
	var v179 vm.Value
	var v191 int
	var v203 int
	var v215 vm.Value
	var v227 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__29024_5, blocks_6, arg__29029_9, insts_10, i_11, label_12, insts_13, blocks_14, v182, v194, v206, v218, v230, arg__29034_28, v29, f_17, i_18, label_19, insts_20, blocks_21, v181, v193, v205, v217, v229, blk_32, term_id_34, f_22, i_23, label_24, insts_25, blocks_26, v189, v201, v213, v225, v237, v162, f_163, i_164, label_165, insts_166, blocks_167, f_35, i_36, label_37, insts_38, blocks_39, blk_40, term_id_41, v186, v198, v210, v222, v234, term_51, op_55, aux_59, v81, f_42, i_43, label_44, insts_45, blocks_46, blk_47, term_id_48, v184, v196, v208, v220, v232, v150, f_151, i_152, label_153, insts_154, blocks_155, blk_156, term_id_157, v188, v200, v212, v224, v236, v158, f_60, i_61, label_62, insts_63, blocks_64, blk_65, term_id_66, term_67, op_68, aux_69, v180, v192, v204, v216, v228, v84, f_70, i_71, label_72, insts_73, blocks_74, blk_75, term_id_76, term_77, op_78, aux_79, v187, v199, v211, v223, v235, v107, v136, f_137, i_138, label_139, insts_140, blocks_141, blk_142, term_id_143, term_144, op_145, aux_146, v183, v195, v207, v219, v231, f_86, i_87, label_88, insts_89, blocks_90, blk_91, term_id_92, term_93, op_94, aux_95, v178, v190, v202, v214, v226, arg__29075_110, arg__29083_113, v114, arg__29090_116, arg__29098_119, v120, f_96, i_97, label_98, insts_99, blocks_100, blk_101, term_id_102, term_103, op_104, aux_105, v185, v197, v209, v221, v233, v124, f_125, i_126, label_127, insts_128, blocks_129, blk_130, term_id_131, term_132, op_133, aux_134, v179, v191, v203, v215, v227
	arg__29024_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	blocks_6, callErr = rt.InvokeValue(vm.Keyword("blocks"), []vm.Value{arg__29024_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__29029_9, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	insts_10, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__29029_9})
	if callErr != nil {
		return nil, callErr
	}
	i_11 = 0
	label_12 = arg1
	insts_13 = insts_10
	blocks_14 = blocks_6
	v182 = vm.Keyword("term")
	v194 = 0
	v206 = 2
	v218 = vm.Keyword("branch")
	v230 = vm.Keyword("branch-if")
	goto b1
b1:
	;
	arg__29034_28, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{blocks_14})
	if callErr != nil {
		return nil, callErr
	}
	v29 = rt.LtValue(vm.Int(i_11), arg__29034_28)
	if v29 {
		f_17 = arg0
		i_18 = i_11
		label_19 = label_12
		insts_20 = insts_13
		blocks_21 = blocks_14
		v181 = v182
		v193 = v194
		v205 = v206
		v217 = v218
		v229 = v230
		goto b2
	} else {
		f_22 = arg0
		i_23 = i_11
		label_24 = label_12
		insts_25 = insts_13
		blocks_26 = blocks_14
		v189 = v182
		v201 = v194
		v213 = v206
		v225 = v218
		v237 = v230
		goto b3
	}
b2:
	;
	blk_32, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{blocks_21, vm.Int(i_18)})
	if callErr != nil {
		return nil, callErr
	}
	term_id_34, callErr = rt.InvokeValue(v181, []vm.Value{blk_32})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(term_id_34) {
		f_35 = f_17
		i_36 = i_18
		label_37 = label_19
		insts_38 = insts_20
		blocks_39 = blocks_21
		blk_40 = blk_32
		term_id_41 = term_id_34
		v186 = v181
		v198 = v193
		v210 = v205
		v222 = v217
		v234 = v229
		goto b5
	} else {
		f_42 = f_17
		i_43 = i_18
		label_44 = label_19
		insts_45 = insts_20
		blocks_46 = blocks_21
		blk_47 = blk_32
		term_id_48 = term_id_34
		v184 = v181
		v196 = v193
		v208 = v205
		v220 = v217
		v232 = v229
		goto b6
	}
b3:
	;
	v162 = vm.NIL
	f_163 = f_22
	i_164 = i_23
	label_165 = label_24
	insts_166 = insts_25
	blocks_167 = blocks_26
	goto b4
b4:
	;
	return v162, nil
b5:
	;
	term_51, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{insts_38, term_id_41})
	if callErr != nil {
		return nil, callErr
	}
	op_55, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{term_51, vm.Int(v198)})
	if callErr != nil {
		return nil, callErr
	}
	aux_59, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{term_51, vm.Int(v210)})
	if callErr != nil {
		return nil, callErr
	}
	v81 = op_55 == v222
	if v81 {
		f_60 = f_35
		i_61 = i_36
		label_62 = label_37
		insts_63 = insts_38
		blocks_64 = blocks_39
		blk_65 = blk_40
		term_id_66 = term_id_41
		term_67 = term_51
		op_68 = op_55
		aux_69 = aux_59
		v180 = v186
		v192 = v198
		v204 = v210
		v216 = v222
		v228 = v234
		goto b8
	} else {
		f_70 = f_35
		i_71 = i_36
		label_72 = label_37
		insts_73 = insts_38
		blocks_74 = blocks_39
		blk_75 = blk_40
		term_id_76 = term_id_41
		term_77 = term_51
		op_78 = op_55
		aux_79 = aux_59
		v187 = v186
		v199 = v198
		v211 = v210
		v223 = v222
		v235 = v234
		goto b9
	}
b6:
	;
	v150 = vm.NIL
	f_151 = f_42
	i_152 = i_43
	label_153 = label_44
	insts_154 = insts_45
	blocks_155 = blocks_46
	blk_156 = blk_47
	term_id_157 = term_id_48
	v188 = v184
	v200 = v196
	v212 = v208
	v224 = v220
	v236 = v232
	goto b7
b7:
	;
	v158 = i_152 + 1
	i_11 = v158
	label_12 = label_153
	insts_13 = insts_154
	blocks_14 = blocks_155
	v182 = v188
	v194 = v200
	v206 = v212
	v218 = v224
	v230 = v236
	goto b1
b8:
	;
	v84, callErr = rt.InvokeValue(rt.LookupVar("ir.validate", "check-branch-arg-arity-for-target!").Deref(), []vm.Value{label_62, vm.Int(i_61), aux_69, blocks_64})
	if callErr != nil {
		return nil, callErr
	}
	v136 = v84
	f_137 = f_60
	i_138 = i_61
	label_139 = label_62
	insts_140 = insts_63
	blocks_141 = blocks_64
	blk_142 = blk_65
	term_id_143 = term_id_66
	term_144 = term_67
	op_145 = op_68
	aux_146 = aux_69
	v183 = v180
	v195 = v192
	v207 = v204
	v219 = v216
	v231 = v228
	goto b10
b9:
	;
	v107 = op_78 == v235
	if v107 {
		f_86 = f_70
		i_87 = i_71
		label_88 = label_72
		insts_89 = insts_73
		blocks_90 = blocks_74
		blk_91 = blk_75
		term_id_92 = term_id_76
		term_93 = term_77
		op_94 = op_78
		aux_95 = aux_79
		v178 = v187
		v190 = v199
		v202 = v211
		v214 = v223
		v226 = v235
		goto b11
	} else {
		f_96 = f_70
		i_97 = i_71
		label_98 = label_72
		insts_99 = insts_73
		blocks_100 = blocks_74
		blk_101 = blk_75
		term_id_102 = term_id_76
		term_103 = term_77
		op_104 = op_78
		aux_105 = aux_79
		v185 = v187
		v197 = v199
		v209 = v211
		v221 = v223
		v233 = v235
		goto b12
	}
b10:
	;
	v150 = v136
	f_151 = f_137
	i_152 = i_138
	label_153 = label_139
	insts_154 = insts_140
	blocks_155 = blocks_141
	blk_156 = blk_142
	term_id_157 = term_id_143
	v188 = v183
	v200 = v195
	v212 = v207
	v224 = v219
	v236 = v231
	goto b7
b11:
	;
	arg__29075_110, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-true").Deref(), []vm.Value{aux_95})
	if callErr != nil {
		return nil, callErr
	}
	arg__29083_113, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-true").Deref(), []vm.Value{aux_95})
	if callErr != nil {
		return nil, callErr
	}
	v114, callErr = rt.InvokeValue(rt.LookupVar("ir.validate", "check-branch-arg-arity-for-target!").Deref(), []vm.Value{label_88, vm.Int(i_87), arg__29083_113, blocks_90})
	if callErr != nil {
		return nil, callErr
	}
	arg__29090_116, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-false").Deref(), []vm.Value{aux_95})
	if callErr != nil {
		return nil, callErr
	}
	arg__29098_119, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-false").Deref(), []vm.Value{aux_95})
	if callErr != nil {
		return nil, callErr
	}
	v120, callErr = rt.InvokeValue(rt.LookupVar("ir.validate", "check-branch-arg-arity-for-target!").Deref(), []vm.Value{label_88, vm.Int(i_87), arg__29098_119, blocks_90})
	if callErr != nil {
		return nil, callErr
	}
	v124 = v120
	f_125 = f_86
	i_126 = i_87
	label_127 = label_88
	insts_128 = insts_89
	blocks_129 = blocks_90
	blk_130 = blk_91
	term_id_131 = term_id_92
	term_132 = term_93
	op_133 = op_94
	aux_134 = aux_95
	v179 = v178
	v191 = v190
	v203 = v202
	v215 = v214
	v227 = v226
	goto b13
b12:
	;
	v124 = vm.NIL
	f_125 = f_96
	i_126 = i_97
	label_127 = label_98
	insts_128 = insts_99
	blocks_129 = blocks_100
	blk_130 = blk_101
	term_id_131 = term_id_102
	term_132 = term_103
	op_133 = op_104
	aux_134 = aux_105
	v179 = v185
	v191 = v197
	v203 = v209
	v215 = v221
	v227 = v233
	goto b13
b13:
	;
	v136 = v124
	f_137 = f_125
	i_138 = i_126
	label_139 = label_127
	insts_140 = insts_128
	blocks_141 = blocks_129
	blk_142 = blk_130
	term_id_143 = term_id_131
	term_144 = term_132
	op_145 = op_133
	aux_146 = aux_134
	v183 = v179
	v195 = v191
	v207 = v203
	v219 = v215
	v231 = v227
	goto b10
}
func validate_fn_BANG_(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__29105_5 vm.Value
	var insts_6 vm.Value
	var v8 vm.Value
	var v10 vm.Value
	var v12 vm.Value
	var v14 vm.Value
	var v16 vm.Value
	var v18 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _ = arg__29105_5, insts_6, v8, v10, v12, v14, v16, v18
	arg__29105_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	insts_6, callErr = rt.InvokeValue(vm.Keyword("insts"), []vm.Value{arg__29105_5})
	if callErr != nil {
		return nil, callErr
	}
	v8, callErr = rt.InvokeValue(rt.LookupVar("ir.validate", "check-inst-shapes!").Deref(), []vm.Value{insts_6, arg1})
	if callErr != nil {
		return nil, callErr
	}
	v10, callErr = rt.InvokeValue(rt.LookupVar("ir.validate", "check-refs-in-range!").Deref(), []vm.Value{insts_6, arg1})
	if callErr != nil {
		return nil, callErr
	}
	v12, callErr = rt.InvokeValue(rt.LookupVar("ir.validate", "check-blocks-terminated!").Deref(), []vm.Value{arg0, arg1})
	if callErr != nil {
		return nil, callErr
	}
	v14, callErr = rt.InvokeValue(rt.LookupVar("ir.validate", "check-no-cross-block-refs!").Deref(), []vm.Value{arg0, arg1})
	if callErr != nil {
		return nil, callErr
	}
	v16, callErr = rt.InvokeValue(rt.LookupVar("ir.validate", "check-branch-arg-arities!").Deref(), []vm.Value{arg0, arg1})
	if callErr != nil {
		return nil, callErr
	}
	v18, callErr = rt.InvokeValue(rt.LookupVar("ir.validate", "check-branch-if-symmetric-args!").Deref(), []vm.Value{arg0, arg1})
	if callErr != nil {
		return nil, callErr
	}
	return arg0, nil
}
