package ir_passes_trace

import (
	rt "github.com/nooga/let-go/pkg/rt"
	vm "github.com/nooga/let-go/pkg/vm"
)

func live_inst_count(arg0 vm.Value) (vm.Value, error) {
	var arg__30525_7 vm.Value
	var arg__30544_13 vm.Value
	var arg__30545_14 vm.Value
	var arg__30565_21 vm.Value
	var arg__30584_27 vm.Value
	var arg__30585_28 vm.Value
	var v29 vm.Value
	var callErr error
	_, _, _, _, _, _, _ = arg__30525_7, arg__30544_13, arg__30545_14, arg__30565_21, arg__30584_27, arg__30585_28, v29
	arg__30525_7, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__30544_13, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__30545_14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var arg__30532_3 vm.Value
		var arg__30539_6 vm.Value
		var v7 vm.Value
		var callErr error
		_, _, _ = arg__30532_3, arg__30539_6, v7
		arg__30532_3, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-insts").Deref(), []vm.Value{arg0, arg0})
		if callErr != nil {
			return nil, callErr
		}
		arg__30539_6, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-insts").Deref(), []vm.Value{arg0, arg0})
		if callErr != nil {
			return nil, callErr
		}
		v7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{arg__30539_6})
		if callErr != nil {
			return nil, callErr
		}
		return v7, nil
	}), arg__30544_13})
	if callErr != nil {
		return nil, callErr
	}
	arg__30565_21, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__30584_27, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__30585_28, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var arg__30572_3 vm.Value
		var arg__30579_6 vm.Value
		var v7 vm.Value
		var callErr error
		_, _, _ = arg__30572_3, arg__30579_6, v7
		arg__30572_3, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-insts").Deref(), []vm.Value{arg0, arg0})
		if callErr != nil {
			return nil, callErr
		}
		arg__30579_6, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-insts").Deref(), []vm.Value{arg0, arg0})
		if callErr != nil {
			return nil, callErr
		}
		v7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{arg__30579_6})
		if callErr != nil {
			return nil, callErr
		}
		return v7, nil
	}), arg__30584_27})
	if callErr != nil {
		return nil, callErr
	}
	v29, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "reduce").Deref(), []vm.Value{rt.LookupVar("clojure.core", "+").Deref(), arg__30585_28})
	if callErr != nil {
		return nil, callErr
	}
	return v29, nil
}
func ns_now() (vm.Value, error) {
	var v1 vm.Value
	var callErr error
	_ = v1
	v1, callErr = rt.InvokeValue(rt.LookupVar("System", "nanoTime").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	return v1, nil
}
func dump_pair(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var before_dump_5 vm.Value
	var before_cnt_7 vm.Value
	var t0_9 vm.Value
	var __10 vm.Value
	var t1_12 vm.Value
	var after_cnt_14 vm.Value
	var after_dump_16 vm.Value
	var arg__30608_19 vm.Value
	var arg__30615_24 vm.Value
	var arg__30623_28 vm.Value
	var v30 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _ = before_dump_5, before_cnt_7, t0_9, __10, t1_12, after_cnt_14, after_dump_16, arg__30608_19, arg__30615_24, arg__30623_28, v30
	before_dump_5, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "dump").Deref(), []vm.Value{arg2})
	if callErr != nil {
		return nil, callErr
	}
	before_cnt_7, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.trace", "live-inst-count").Deref(), []vm.Value{arg2})
	if callErr != nil {
		return nil, callErr
	}
	t0_9, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.trace", "ns-now").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	__10, callErr = rt.InvokeValue(arg1, []vm.Value{arg2})
	if callErr != nil {
		return nil, callErr
	}
	t1_12, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.trace", "ns-now").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	after_cnt_14, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.trace", "live-inst-count").Deref(), []vm.Value{arg2})
	if callErr != nil {
		return nil, callErr
	}
	after_dump_16, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "dump").Deref(), []vm.Value{arg2})
	if callErr != nil {
		return nil, callErr
	}
	arg__30608_19 = rt.SubValue(t1_12, t0_9)
	arg__30615_24, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "/").Deref(), []vm.Value{arg__30608_19, vm.Float(1e+06)})
	if callErr != nil {
		return nil, callErr
	}
	arg__30623_28 = rt.SubValue(before_cnt_7, after_cnt_14)
	v30, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("ms"), arg__30615_24, vm.Keyword("pass"), arg0, vm.Keyword("after"), after_dump_16, vm.Keyword("delta"), arg__30623_28, vm.Keyword("before"), before_dump_5})
	if callErr != nil {
		return nil, callErr
	}
	return v30, nil
}
func optimize_fn_traced(arg0 vm.Value) (vm.Value, error) {
	var v5 vm.Value
	var arg__30636_9 vm.Value
	var arg__30644_14 vm.Value
	var v15 vm.Value
	var iter_16 int
	var f_17 vm.Value
	var v202 string
	var v209 string
	var v216 string
	var v223 string
	var v230 string
	var v237 string
	var v244 string
	var v251 string
	var v258 int
	var v265 vm.Value
	var before_21 vm.Value
	var arg__30654_24 vm.Value
	var arg__30662_28 vm.Value
	var v29 vm.Value
	var arg__30670_33 vm.Value
	var arg__30678_38 vm.Value
	var v39 vm.Value
	var arg__30684_42 vm.Value
	var arg__30692_46 vm.Value
	var v47 vm.Value
	var arg__30700_51 vm.Value
	var arg__30708_56 vm.Value
	var v57 vm.Value
	var arg__30714_60 vm.Value
	var arg__30722_64 vm.Value
	var v65 vm.Value
	var arg__30730_69 vm.Value
	var arg__30738_74 vm.Value
	var v75 vm.Value
	var arg__30744_78 vm.Value
	var arg__30752_82 vm.Value
	var v83 vm.Value
	var arg__30760_87 vm.Value
	var arg__30768_92 vm.Value
	var v93 vm.Value
	var after_95 vm.Value
	var v104 bool
	var iter_96 int
	var f_97 vm.Value
	var before_98 vm.Value
	var after_99 vm.Value
	var v205 string
	var v212 string
	var v219 string
	var v226 string
	var v233 string
	var v240 string
	var v247 string
	var v254 string
	var v261 int
	var v268 vm.Value
	var iter_100 int
	var f_101 vm.Value
	var before_102 vm.Value
	var after_103 vm.Value
	var v203 string
	var v210 string
	var v217 string
	var v224 string
	var v231 string
	var v238 string
	var v245 string
	var v252 string
	var v259 int
	var v266 vm.Value
	var v116 bool
	var v159 vm.Value
	var iter_160 int
	var f_161 vm.Value
	var before_162 vm.Value
	var after_163 vm.Value
	var arg__30799_167 vm.Value
	var arg__30807_172 vm.Value
	var v173 vm.Value
	var iter_107 int
	var f_108 vm.Value
	var before_109 vm.Value
	var after_110 vm.Value
	var v206 string
	var v213 string
	var v220 string
	var v227 string
	var v234 string
	var v241 string
	var v248 string
	var v255 string
	var v262 int
	var v269 vm.Value
	var arg__30783_123 vm.Value
	var arg__30792_130 vm.Value
	var v131 vm.Value
	var iter_111 int
	var f_112 vm.Value
	var before_113 vm.Value
	var after_114 vm.Value
	var v204 string
	var v211 string
	var v218 string
	var v225 string
	var v232 string
	var v239 string
	var v246 string
	var v253 string
	var v260 int
	var v267 vm.Value
	var v153 vm.Value
	var iter_154 int
	var f_155 vm.Value
	var before_156 vm.Value
	var after_157 vm.Value
	var iter_133 int
	var f_134 vm.Value
	var before_135 vm.Value
	var after_136 vm.Value
	var v201 string
	var v208 string
	var v215 string
	var v222 string
	var v229 string
	var v236 string
	var v243 string
	var v250 string
	var v257 int
	var v264 vm.Value
	var v143 int
	var iter_137 int
	var f_138 vm.Value
	var before_139 vm.Value
	var after_140 vm.Value
	var v207 string
	var v214 string
	var v221 string
	var v228 string
	var v235 string
	var v242 string
	var v249 string
	var v256 string
	var v263 int
	var v270 vm.Value
	var v147 vm.Value
	var iter_148 int
	var f_149 vm.Value
	var before_150 vm.Value
	var after_151 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = v5, arg__30636_9, arg__30644_14, v15, iter_16, f_17, v202, v209, v216, v223, v230, v237, v244, v251, v258, v265, before_21, arg__30654_24, arg__30662_28, v29, arg__30670_33, arg__30678_38, v39, arg__30684_42, arg__30692_46, v47, arg__30700_51, arg__30708_56, v57, arg__30714_60, arg__30722_64, v65, arg__30730_69, arg__30738_74, v75, arg__30744_78, arg__30752_82, v83, arg__30760_87, arg__30768_92, v93, after_95, v104, iter_96, f_97, before_98, after_99, v205, v212, v219, v226, v233, v240, v247, v254, v261, v268, iter_100, f_101, before_102, after_103, v203, v210, v217, v224, v231, v238, v245, v252, v259, v266, v116, v159, iter_160, f_161, before_162, after_163, arg__30799_167, arg__30807_172, v173, iter_107, f_108, before_109, after_110, v206, v213, v220, v227, v234, v241, v248, v255, v262, v269, arg__30783_123, arg__30792_130, v131, iter_111, f_112, before_113, after_114, v204, v211, v218, v225, v232, v239, v246, v253, v260, v267, v153, iter_154, f_155, before_156, after_157, iter_133, f_134, before_135, after_136, v201, v208, v215, v222, v229, v236, v243, v250, v257, v264, v143, iter_137, f_138, before_139, after_140, v207, v214, v221, v228, v235, v242, v249, v256, v263, v270, v147, iter_148, f_149, before_150, after_151
	v5, callErr = rt.InvokeValue(rt.LookupVar("ir.validate", "validate-fn!").Deref(), []vm.Value{arg0, vm.String("build")})
	if callErr != nil {
		return nil, callErr
	}
	arg__30636_9, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "typeinfer").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__30644_14, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "typeinfer").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v15, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.trace", "trace-pass").Deref(), []vm.Value{vm.String("typeinfer-pre"), vm.Int(-1), arg__30644_14, arg0})
	if callErr != nil {
		return nil, callErr
	}
	iter_16 = 0
	f_17 = arg0
	v202 = "constfold"
	v209 = "constfold/"
	v216 = "cse"
	v223 = "cse/"
	v230 = "licm"
	v237 = "licm/"
	v244 = "dce"
	v251 = "dce/"
	v258 = 15
	v265 = vm.Keyword("else")
	goto b1
b1:
	;
	before_21, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.trace", "live-inst-count").Deref(), []vm.Value{f_17})
	if callErr != nil {
		return nil, callErr
	}
	arg__30654_24, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.constfold", "constfold").Deref(), []vm.Value{f_17})
	if callErr != nil {
		return nil, callErr
	}
	arg__30662_28, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.constfold", "constfold").Deref(), []vm.Value{f_17})
	if callErr != nil {
		return nil, callErr
	}
	v29, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.trace", "trace-pass").Deref(), []vm.Value{vm.String(v202), vm.Int(iter_16), arg__30662_28, f_17})
	if callErr != nil {
		return nil, callErr
	}
	arg__30670_33, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(v209), vm.Int(iter_16)})
	if callErr != nil {
		return nil, callErr
	}
	arg__30678_38, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(v209), vm.Int(iter_16)})
	if callErr != nil {
		return nil, callErr
	}
	v39, callErr = rt.InvokeValue(rt.LookupVar("ir.validate", "validate-fn!").Deref(), []vm.Value{f_17, arg__30678_38})
	if callErr != nil {
		return nil, callErr
	}
	arg__30684_42, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.cse", "cse").Deref(), []vm.Value{f_17})
	if callErr != nil {
		return nil, callErr
	}
	arg__30692_46, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.cse", "cse").Deref(), []vm.Value{f_17})
	if callErr != nil {
		return nil, callErr
	}
	v47, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.trace", "trace-pass").Deref(), []vm.Value{vm.String(v216), vm.Int(iter_16), arg__30692_46, f_17})
	if callErr != nil {
		return nil, callErr
	}
	arg__30700_51, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(v223), vm.Int(iter_16)})
	if callErr != nil {
		return nil, callErr
	}
	arg__30708_56, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(v223), vm.Int(iter_16)})
	if callErr != nil {
		return nil, callErr
	}
	v57, callErr = rt.InvokeValue(rt.LookupVar("ir.validate", "validate-fn!").Deref(), []vm.Value{f_17, arg__30708_56})
	if callErr != nil {
		return nil, callErr
	}
	arg__30714_60, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.licm", "licm").Deref(), []vm.Value{f_17})
	if callErr != nil {
		return nil, callErr
	}
	arg__30722_64, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.licm", "licm").Deref(), []vm.Value{f_17})
	if callErr != nil {
		return nil, callErr
	}
	v65, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.trace", "trace-pass").Deref(), []vm.Value{vm.String(v230), vm.Int(iter_16), arg__30722_64, f_17})
	if callErr != nil {
		return nil, callErr
	}
	arg__30730_69, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(v237), vm.Int(iter_16)})
	if callErr != nil {
		return nil, callErr
	}
	arg__30738_74, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(v237), vm.Int(iter_16)})
	if callErr != nil {
		return nil, callErr
	}
	v75, callErr = rt.InvokeValue(rt.LookupVar("ir.validate", "validate-fn!").Deref(), []vm.Value{f_17, arg__30738_74})
	if callErr != nil {
		return nil, callErr
	}
	arg__30744_78, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.dce", "dce").Deref(), []vm.Value{f_17})
	if callErr != nil {
		return nil, callErr
	}
	arg__30752_82, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.dce", "dce").Deref(), []vm.Value{f_17})
	if callErr != nil {
		return nil, callErr
	}
	v83, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.trace", "trace-pass").Deref(), []vm.Value{vm.String(v244), vm.Int(iter_16), arg__30752_82, f_17})
	if callErr != nil {
		return nil, callErr
	}
	arg__30760_87, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(v251), vm.Int(iter_16)})
	if callErr != nil {
		return nil, callErr
	}
	arg__30768_92, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(v251), vm.Int(iter_16)})
	if callErr != nil {
		return nil, callErr
	}
	v93, callErr = rt.InvokeValue(rt.LookupVar("ir.validate", "validate-fn!").Deref(), []vm.Value{f_17, arg__30768_92})
	if callErr != nil {
		return nil, callErr
	}
	after_95, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.trace", "live-inst-count").Deref(), []vm.Value{f_17})
	if callErr != nil {
		return nil, callErr
	}
	v104 = before_21 == after_95
	if v104 {
		iter_96 = iter_16
		f_97 = f_17
		before_98 = before_21
		after_99 = after_95
		v205 = v202
		v212 = v209
		v219 = v216
		v226 = v223
		v233 = v230
		v240 = v237
		v247 = v244
		v254 = v251
		v261 = v258
		v268 = v265
		goto b2
	} else {
		iter_100 = iter_16
		f_101 = f_17
		before_102 = before_21
		after_103 = after_95
		v203 = v202
		v210 = v209
		v217 = v216
		v224 = v223
		v231 = v230
		v238 = v237
		v245 = v244
		v252 = v251
		v259 = v258
		v266 = v265
		goto b3
	}
b2:
	;
	v159 = f_97
	iter_160 = iter_96
	f_161 = f_97
	before_162 = before_98
	after_163 = after_99
	goto b4
b3:
	;
	v116 = iter_100 >= v259
	if v116 {
		iter_107 = iter_100
		f_108 = f_101
		before_109 = before_102
		after_110 = after_103
		v206 = v203
		v213 = v210
		v220 = v217
		v227 = v224
		v234 = v231
		v241 = v238
		v248 = v245
		v255 = v252
		v262 = v259
		v269 = v266
		goto b5
	} else {
		iter_111 = iter_100
		f_112 = f_101
		before_113 = before_102
		after_114 = after_103
		v204 = v203
		v211 = v210
		v218 = v217
		v225 = v224
		v232 = v231
		v239 = v238
		v246 = v245
		v253 = v252
		v260 = v259
		v267 = v266
		goto b6
	}
b4:
	;
	arg__30799_167, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "typeinfer").Deref(), []vm.Value{f_161})
	if callErr != nil {
		return nil, callErr
	}
	arg__30807_172, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "typeinfer").Deref(), []vm.Value{f_161})
	if callErr != nil {
		return nil, callErr
	}
	v173, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.trace", "trace-pass").Deref(), []vm.Value{vm.String("typeinfer-post"), vm.Int(-1), arg__30807_172, f_161})
	if callErr != nil {
		return nil, callErr
	}
	return f_161, nil
b5:
	;
	arg__30783_123, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("warn: optimize-fn-traced max iters (16) reached, "), before_109, vm.String(" insts after 16 cycles")})
	if callErr != nil {
		return nil, callErr
	}
	arg__30792_130, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("warn: optimize-fn-traced max iters (16) reached, "), before_109, vm.String(" insts after 16 cycles")})
	if callErr != nil {
		return nil, callErr
	}
	v131, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "println").Deref(), []vm.Value{arg__30792_130})
	if callErr != nil {
		return nil, callErr
	}
	v153 = f_108
	iter_154 = iter_107
	f_155 = f_108
	before_156 = before_109
	after_157 = after_110
	goto b7
b6:
	;
	if vm.IsTruthy(v267) {
		iter_133 = iter_111
		f_134 = f_112
		before_135 = before_113
		after_136 = after_114
		v201 = v204
		v208 = v211
		v215 = v218
		v222 = v225
		v229 = v232
		v236 = v239
		v243 = v246
		v250 = v253
		v257 = v260
		v264 = v267
		goto b8
	} else {
		iter_137 = iter_111
		f_138 = f_112
		before_139 = before_113
		after_140 = after_114
		v207 = v204
		v214 = v211
		v221 = v218
		v228 = v225
		v235 = v232
		v242 = v239
		v249 = v246
		v256 = v253
		v263 = v260
		v270 = v267
		goto b9
	}
b7:
	;
	v159 = v153
	iter_160 = iter_154
	f_161 = f_155
	before_162 = before_156
	after_163 = after_157
	goto b4
b8:
	;
	v143 = iter_133 + 1
	iter_16 = v143
	f_17 = f_134
	v202 = v201
	v209 = v208
	v216 = v215
	v223 = v222
	v230 = v229
	v237 = v236
	v244 = v243
	v251 = v250
	v258 = v257
	v265 = v264
	goto b1
b9:
	;
	v147 = vm.NIL
	iter_148 = iter_137
	f_149 = f_138
	before_150 = before_139
	after_151 = after_140
	goto b10
b10:
	;
	v153 = v147
	iter_154 = iter_148
	f_155 = f_149
	before_156 = before_150
	after_157 = after_151
	goto b7
}
func print_trace(arg0 vm.Value) (vm.Value, error) {
	var arg__30826_17 vm.Value
	var arg__30843_34 vm.Value
	var v35 vm.Value
	var arg__30850_42 vm.Value
	var arg__30858_50 vm.Value
	var arg__30859_51 vm.Value
	var arg__30867_59 vm.Value
	var arg__30875_67 vm.Value
	var arg__30876_68 vm.Value
	var v69 vm.Value
	var doseq_seq__30809_71 vm.Value
	var doseq_loop__30810_72 vm.Value
	var v256 string
	var v259 vm.Value
	var v262 vm.Value
	var v265 vm.Value
	var v268 vm.Value
	var v271 vm.Value
	var v274 vm.Value
	var trace_74 vm.Value
	var doseq_seq__30809_75 vm.Value
	var doseq_loop__30810_76 vm.Value
	var v255 string
	var v258 vm.Value
	var v261 vm.Value
	var v264 vm.Value
	var v267 vm.Value
	var v270 vm.Value
	var v273 vm.Value
	var e_82 vm.Value
	var arg__30886_85 vm.Value
	var arg__30889_87 vm.Value
	var arg__30892_89 vm.Value
	var arg__30895_91 vm.Value
	var arg__30898_93 vm.Value
	var arg__30901_95 vm.Value
	var arg__30906_99 vm.Value
	var arg__30909_101 vm.Value
	var arg__30912_103 vm.Value
	var arg__30915_105 vm.Value
	var arg__30918_107 vm.Value
	var arg__30921_109 vm.Value
	var arg__30922_110 vm.Value
	var arg__30927_114 vm.Value
	var arg__30930_116 vm.Value
	var arg__30933_118 vm.Value
	var arg__30936_120 vm.Value
	var arg__30939_122 vm.Value
	var arg__30942_124 vm.Value
	var arg__30947_128 vm.Value
	var arg__30950_130 vm.Value
	var arg__30953_132 vm.Value
	var arg__30956_134 vm.Value
	var arg__30959_136 vm.Value
	var arg__30962_138 vm.Value
	var arg__30963_139 vm.Value
	var v140 vm.Value
	var v142 vm.Value
	var trace_77 vm.Value
	var doseq_seq__30809_78 vm.Value
	var doseq_loop__30810_79 vm.Value
	var v257 string
	var v260 vm.Value
	var v263 vm.Value
	var v266 vm.Value
	var v269 vm.Value
	var v272 vm.Value
	var v275 vm.Value
	var v146 vm.Value
	var trace_147 vm.Value
	var doseq_seq__30809_148 vm.Value
	var doseq_loop__30810_149 vm.Value
	var arg__30973_154 vm.Value
	var arg__30981_160 vm.Value
	var total_ms_161 vm.Value
	var arg__30988_166 vm.Value
	var arg__30996_172 vm.Value
	var total_removed_173 vm.Value
	var arg__31003_180 vm.Value
	var arg__31011_188 vm.Value
	var arg__31012_189 vm.Value
	var arg__31020_197 vm.Value
	var arg__31028_205 vm.Value
	var arg__31029_206 vm.Value
	var v207 vm.Value
	var arg__31034_210 vm.Value
	var arg__31042_214 vm.Value
	var arg__31045_215 vm.Value
	var arg__31051_219 vm.Value
	var arg__31059_223 vm.Value
	var arg__31062_224 vm.Value
	var v225 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__30826_17, arg__30843_34, v35, arg__30850_42, arg__30858_50, arg__30859_51, arg__30867_59, arg__30875_67, arg__30876_68, v69, doseq_seq__30809_71, doseq_loop__30810_72, v256, v259, v262, v265, v268, v271, v274, trace_74, doseq_seq__30809_75, doseq_loop__30810_76, v255, v258, v261, v264, v267, v270, v273, e_82, arg__30886_85, arg__30889_87, arg__30892_89, arg__30895_91, arg__30898_93, arg__30901_95, arg__30906_99, arg__30909_101, arg__30912_103, arg__30915_105, arg__30918_107, arg__30921_109, arg__30922_110, arg__30927_114, arg__30930_116, arg__30933_118, arg__30936_120, arg__30939_122, arg__30942_124, arg__30947_128, arg__30950_130, arg__30953_132, arg__30956_134, arg__30959_136, arg__30962_138, arg__30963_139, v140, v142, trace_77, doseq_seq__30809_78, doseq_loop__30810_79, v257, v260, v263, v266, v269, v272, v275, v146, trace_147, doseq_seq__30809_148, doseq_loop__30810_149, arg__30973_154, arg__30981_160, total_ms_161, arg__30988_166, arg__30996_172, total_removed_173, arg__31003_180, arg__31011_188, arg__31012_189, arg__31020_197, arg__31028_205, arg__31029_206, v207, arg__31034_210, arg__31042_214, arg__31045_215, arg__31051_219, arg__31059_223, arg__31062_224, v225
	arg__30826_17, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "format").Deref(), []vm.Value{vm.String("%4s  %-18s  %6s  %6s  %6s  %7s"), vm.String("iter"), vm.String("pass"), vm.String("before"), vm.String("after"), vm.String("delta"), vm.String("ms")})
	if callErr != nil {
		return nil, callErr
	}
	arg__30843_34, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "format").Deref(), []vm.Value{vm.String("%4s  %-18s  %6s  %6s  %6s  %7s"), vm.String("iter"), vm.String("pass"), vm.String("before"), vm.String("after"), vm.String("delta"), vm.String("ms")})
	if callErr != nil {
		return nil, callErr
	}
	v35, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "println").Deref(), []vm.Value{arg__30843_34})
	if callErr != nil {
		return nil, callErr
	}
	arg__30850_42, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "repeat").Deref(), []vm.Value{vm.Int(55), vm.String("-")})
	if callErr != nil {
		return nil, callErr
	}
	arg__30858_50, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "repeat").Deref(), []vm.Value{vm.Int(55), vm.String("-")})
	if callErr != nil {
		return nil, callErr
	}
	arg__30859_51, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "apply").Deref(), []vm.Value{rt.LookupVar("clojure.core", "str").Deref(), arg__30858_50})
	if callErr != nil {
		return nil, callErr
	}
	arg__30867_59, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "repeat").Deref(), []vm.Value{vm.Int(55), vm.String("-")})
	if callErr != nil {
		return nil, callErr
	}
	arg__30875_67, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "repeat").Deref(), []vm.Value{vm.Int(55), vm.String("-")})
	if callErr != nil {
		return nil, callErr
	}
	arg__30876_68, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "apply").Deref(), []vm.Value{rt.LookupVar("clojure.core", "str").Deref(), arg__30875_67})
	if callErr != nil {
		return nil, callErr
	}
	v69, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "println").Deref(), []vm.Value{arg__30876_68})
	if callErr != nil {
		return nil, callErr
	}
	doseq_seq__30809_71, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__30810_72 = doseq_seq__30809_71
	v256 = "%4d  %-18s  %6d  %6d  %+6d  %7.2f"
	v259 = vm.Keyword("iter")
	v262 = vm.Keyword("pass")
	v265 = vm.Keyword("before")
	v268 = vm.Keyword("after")
	v271 = vm.Keyword("delta")
	v274 = vm.Keyword("ms")
	goto b1
b1:
	;
	if vm.IsTruthy(doseq_loop__30810_72) {
		trace_74 = arg0
		doseq_seq__30809_75 = doseq_seq__30809_71
		doseq_loop__30810_76 = doseq_loop__30810_72
		v255 = v256
		v258 = v259
		v261 = v262
		v264 = v265
		v267 = v268
		v270 = v271
		v273 = v274
		goto b2
	} else {
		trace_77 = arg0
		doseq_seq__30809_78 = doseq_seq__30809_71
		doseq_loop__30810_79 = doseq_loop__30810_72
		v257 = v256
		v260 = v259
		v263 = v262
		v266 = v265
		v269 = v268
		v272 = v271
		v275 = v274
		goto b3
	}
b2:
	;
	e_82, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{doseq_loop__30810_76})
	if callErr != nil {
		return nil, callErr
	}
	arg__30886_85, callErr = rt.InvokeValue(v258, []vm.Value{e_82})
	if callErr != nil {
		return nil, callErr
	}
	arg__30889_87, callErr = rt.InvokeValue(v261, []vm.Value{e_82})
	if callErr != nil {
		return nil, callErr
	}
	arg__30892_89, callErr = rt.InvokeValue(v264, []vm.Value{e_82})
	if callErr != nil {
		return nil, callErr
	}
	arg__30895_91, callErr = rt.InvokeValue(v267, []vm.Value{e_82})
	if callErr != nil {
		return nil, callErr
	}
	arg__30898_93, callErr = rt.InvokeValue(v270, []vm.Value{e_82})
	if callErr != nil {
		return nil, callErr
	}
	arg__30901_95, callErr = rt.InvokeValue(v273, []vm.Value{e_82})
	if callErr != nil {
		return nil, callErr
	}
	arg__30906_99, callErr = rt.InvokeValue(v258, []vm.Value{e_82})
	if callErr != nil {
		return nil, callErr
	}
	arg__30909_101, callErr = rt.InvokeValue(v261, []vm.Value{e_82})
	if callErr != nil {
		return nil, callErr
	}
	arg__30912_103, callErr = rt.InvokeValue(v264, []vm.Value{e_82})
	if callErr != nil {
		return nil, callErr
	}
	arg__30915_105, callErr = rt.InvokeValue(v267, []vm.Value{e_82})
	if callErr != nil {
		return nil, callErr
	}
	arg__30918_107, callErr = rt.InvokeValue(v270, []vm.Value{e_82})
	if callErr != nil {
		return nil, callErr
	}
	arg__30921_109, callErr = rt.InvokeValue(v273, []vm.Value{e_82})
	if callErr != nil {
		return nil, callErr
	}
	arg__30922_110, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "format").Deref(), []vm.Value{vm.String(v255), arg__30906_99, arg__30909_101, arg__30912_103, arg__30915_105, arg__30918_107, arg__30921_109})
	if callErr != nil {
		return nil, callErr
	}
	arg__30927_114, callErr = rt.InvokeValue(v258, []vm.Value{e_82})
	if callErr != nil {
		return nil, callErr
	}
	arg__30930_116, callErr = rt.InvokeValue(v261, []vm.Value{e_82})
	if callErr != nil {
		return nil, callErr
	}
	arg__30933_118, callErr = rt.InvokeValue(v264, []vm.Value{e_82})
	if callErr != nil {
		return nil, callErr
	}
	arg__30936_120, callErr = rt.InvokeValue(v267, []vm.Value{e_82})
	if callErr != nil {
		return nil, callErr
	}
	arg__30939_122, callErr = rt.InvokeValue(v270, []vm.Value{e_82})
	if callErr != nil {
		return nil, callErr
	}
	arg__30942_124, callErr = rt.InvokeValue(v273, []vm.Value{e_82})
	if callErr != nil {
		return nil, callErr
	}
	arg__30947_128, callErr = rt.InvokeValue(v258, []vm.Value{e_82})
	if callErr != nil {
		return nil, callErr
	}
	arg__30950_130, callErr = rt.InvokeValue(v261, []vm.Value{e_82})
	if callErr != nil {
		return nil, callErr
	}
	arg__30953_132, callErr = rt.InvokeValue(v264, []vm.Value{e_82})
	if callErr != nil {
		return nil, callErr
	}
	arg__30956_134, callErr = rt.InvokeValue(v267, []vm.Value{e_82})
	if callErr != nil {
		return nil, callErr
	}
	arg__30959_136, callErr = rt.InvokeValue(v270, []vm.Value{e_82})
	if callErr != nil {
		return nil, callErr
	}
	arg__30962_138, callErr = rt.InvokeValue(v273, []vm.Value{e_82})
	if callErr != nil {
		return nil, callErr
	}
	arg__30963_139, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "format").Deref(), []vm.Value{vm.String(v255), arg__30947_128, arg__30950_130, arg__30953_132, arg__30956_134, arg__30959_136, arg__30962_138})
	if callErr != nil {
		return nil, callErr
	}
	v140, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "println").Deref(), []vm.Value{arg__30963_139})
	if callErr != nil {
		return nil, callErr
	}
	v142, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "next").Deref(), []vm.Value{doseq_loop__30810_76})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__30810_72 = v142
	v256 = v255
	v259 = v258
	v262 = v261
	v265 = v264
	v268 = v267
	v271 = v270
	v274 = v273
	goto b1
b3:
	;
	v146 = vm.NIL
	trace_147 = trace_77
	doseq_seq__30809_148 = doseq_seq__30809_78
	doseq_loop__30810_149 = doseq_loop__30810_79
	goto b4
b4:
	;
	arg__30973_154, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{vm.Keyword("ms"), trace_147})
	if callErr != nil {
		return nil, callErr
	}
	arg__30981_160, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{vm.Keyword("ms"), trace_147})
	if callErr != nil {
		return nil, callErr
	}
	total_ms_161, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "reduce").Deref(), []vm.Value{rt.LookupVar("clojure.core", "+").Deref(), arg__30981_160})
	if callErr != nil {
		return nil, callErr
	}
	arg__30988_166, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{vm.Keyword("delta"), trace_147})
	if callErr != nil {
		return nil, callErr
	}
	arg__30996_172, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{vm.Keyword("delta"), trace_147})
	if callErr != nil {
		return nil, callErr
	}
	total_removed_173, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "reduce").Deref(), []vm.Value{rt.LookupVar("clojure.core", "+").Deref(), arg__30996_172})
	if callErr != nil {
		return nil, callErr
	}
	arg__31003_180, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "repeat").Deref(), []vm.Value{vm.Int(55), vm.String("-")})
	if callErr != nil {
		return nil, callErr
	}
	arg__31011_188, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "repeat").Deref(), []vm.Value{vm.Int(55), vm.String("-")})
	if callErr != nil {
		return nil, callErr
	}
	arg__31012_189, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "apply").Deref(), []vm.Value{rt.LookupVar("clojure.core", "str").Deref(), arg__31011_188})
	if callErr != nil {
		return nil, callErr
	}
	arg__31020_197, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "repeat").Deref(), []vm.Value{vm.Int(55), vm.String("-")})
	if callErr != nil {
		return nil, callErr
	}
	arg__31028_205, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "repeat").Deref(), []vm.Value{vm.Int(55), vm.String("-")})
	if callErr != nil {
		return nil, callErr
	}
	arg__31029_206, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "apply").Deref(), []vm.Value{rt.LookupVar("clojure.core", "str").Deref(), arg__31028_205})
	if callErr != nil {
		return nil, callErr
	}
	v207, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "println").Deref(), []vm.Value{arg__31029_206})
	if callErr != nil {
		return nil, callErr
	}
	arg__31034_210, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{trace_147})
	if callErr != nil {
		return nil, callErr
	}
	arg__31042_214, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{trace_147})
	if callErr != nil {
		return nil, callErr
	}
	arg__31045_215, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "format").Deref(), []vm.Value{vm.String("  %d passes, %d insts removed, %.2f ms total"), arg__31042_214, total_removed_173, total_ms_161})
	if callErr != nil {
		return nil, callErr
	}
	arg__31051_219, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{trace_147})
	if callErr != nil {
		return nil, callErr
	}
	arg__31059_223, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{trace_147})
	if callErr != nil {
		return nil, callErr
	}
	arg__31062_224, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "format").Deref(), []vm.Value{vm.String("  %d passes, %d insts removed, %.2f ms total"), arg__31059_223, total_removed_173, total_ms_161})
	if callErr != nil {
		return nil, callErr
	}
	v225, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "println").Deref(), []vm.Value{arg__31062_224})
	if callErr != nil {
		return nil, callErr
	}
	return v225, nil
}
