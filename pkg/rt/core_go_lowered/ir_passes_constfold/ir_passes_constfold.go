package ir_passes_constfold

import (
	rt "github.com/nooga/let-go/pkg/rt"
	vm "github.com/nooga/let-go/pkg/vm"
)

func const_QMARK_(arg0 vm.Value, arg1 vm.Value) (bool, error) {
	var arg__17330_4 vm.Value
	var v5 bool
	var callErr error
	_, _ = arg__17330_4, v5
	arg__17330_4, callErr = rt.InvokeValue(rt.LookupVar("ir", "op").Deref(), []vm.Value{arg0, arg1})
	if callErr != nil {
		return false, callErr
	}
	v5 = arg__17330_4 == vm.Keyword("const")
	return v5, nil
}
func const_val(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var v7 vm.Value
	var nid_2 vm.Value
	var f_3 vm.Value
	var v10 vm.Value
	var nid_4 vm.Value
	var f_5 vm.Value
	var v14 vm.Value
	var nid_15 vm.Value
	var f_16 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _ = v7, nid_2, f_3, v10, nid_4, f_5, v14, nid_15, f_16
	v7, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.constfold", "const?").Deref(), []vm.Value{arg0, arg1})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v7) {
		nid_2 = arg0
		f_3 = arg1
		goto b1
	} else {
		nid_4 = arg0
		f_5 = arg1
		goto b2
	}
b1:
	;
	v10, callErr = rt.InvokeValue(rt.LookupVar("ir", "aux").Deref(), []vm.Value{nid_2, f_3})
	if callErr != nil {
		return nil, callErr
	}
	v14 = v10
	nid_15 = nid_2
	f_16 = f_3
	goto b3
b2:
	;
	v14 = vm.NIL
	nid_15 = nid_4
	f_16 = f_5
	goto b3
b3:
	;
	return v14, nil
}
func apply_action_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var v11 vm.Value
	var f_4 vm.Value
	var nid_5 vm.Value
	var action_6 vm.Value
	var map__17468_14 vm.Value
	var op_18 vm.Value
	var refs_22 vm.Value
	var aux_26 vm.Value
	var v28 vm.Value
	var v30 vm.Value
	var v32 vm.Value
	var f_7 vm.Value
	var nid_8 vm.Value
	var action_9 vm.Value
	var v41 vm.Value
	var v78 vm.Value
	var f_79 vm.Value
	var nid_80 vm.Value
	var action_81 vm.Value
	var f_34 vm.Value
	var nid_35 vm.Value
	var action_36 vm.Value
	var arg__17515_44 vm.Value
	var arg__17521_47 vm.Value
	var v48 vm.Value
	var f_37 vm.Value
	var nid_38 vm.Value
	var action_39 vm.Value
	var v57 vm.Value
	var v73 vm.Value
	var f_74 vm.Value
	var nid_75 vm.Value
	var action_76 vm.Value
	var f_50 vm.Value
	var nid_51 vm.Value
	var action_52 vm.Value
	var arg__17528_60 vm.Value
	var arg__17534_63 vm.Value
	var v64 vm.Value
	var f_53 vm.Value
	var nid_54 vm.Value
	var action_55 vm.Value
	var v68 vm.Value
	var f_69 vm.Value
	var nid_70 vm.Value
	var action_71 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = v11, f_4, nid_5, action_6, map__17468_14, op_18, refs_22, aux_26, v28, v30, v32, f_7, nid_8, action_9, v41, v78, f_79, nid_80, action_81, f_34, nid_35, action_36, arg__17515_44, arg__17521_47, v48, f_37, nid_38, action_39, v57, v73, f_74, nid_75, action_76, f_50, nid_51, action_52, arg__17528_60, arg__17534_63, v64, f_53, nid_54, action_55, v68, f_69, nid_70, action_71
	v11, callErr = rt.InvokeValue(vm.Keyword("replace-with"), []vm.Value{arg2})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v11) {
		f_4 = arg0
		nid_5 = arg1
		action_6 = arg2
		goto b1
	} else {
		f_7 = arg0
		nid_8 = arg1
		action_9 = arg2
		goto b2
	}
b1:
	;
	map__17468_14, callErr = rt.InvokeValue(vm.Keyword("replace-with"), []vm.Value{action_6})
	if callErr != nil {
		return nil, callErr
	}
	op_18, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{map__17468_14, vm.Keyword("op")})
	if callErr != nil {
		return nil, callErr
	}
	refs_22, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{map__17468_14, vm.Keyword("refs")})
	if callErr != nil {
		return nil, callErr
	}
	aux_26, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{map__17468_14, vm.Keyword("aux")})
	if callErr != nil {
		return nil, callErr
	}
	v28, callErr = rt.InvokeValue(rt.LookupVar("ir", "replace-op!").Deref(), []vm.Value{f_4, nid_5, op_18})
	if callErr != nil {
		return nil, callErr
	}
	v30, callErr = rt.InvokeValue(rt.LookupVar("ir", "set-refs!").Deref(), []vm.Value{f_4, nid_5, refs_22})
	if callErr != nil {
		return nil, callErr
	}
	v32, callErr = rt.InvokeValue(rt.LookupVar("ir", "set-aux!").Deref(), []vm.Value{f_4, nid_5, aux_26})
	if callErr != nil {
		return nil, callErr
	}
	v78 = v32
	f_79 = f_4
	nid_80 = nid_5
	action_81 = action_6
	goto b3
b2:
	;
	v41, callErr = rt.InvokeValue(vm.Keyword("replace-uses"), []vm.Value{action_9})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v41) {
		f_34 = f_7
		nid_35 = nid_8
		action_36 = action_9
		goto b4
	} else {
		f_37 = f_7
		nid_38 = nid_8
		action_39 = action_9
		goto b5
	}
b3:
	;
	return v78, nil
b4:
	;
	arg__17515_44, callErr = rt.InvokeValue(vm.Keyword("replace-uses"), []vm.Value{action_36})
	if callErr != nil {
		return nil, callErr
	}
	arg__17521_47, callErr = rt.InvokeValue(vm.Keyword("replace-uses"), []vm.Value{action_36})
	if callErr != nil {
		return nil, callErr
	}
	v48, callErr = rt.InvokeValue(rt.LookupVar("ir", "replace-all-uses!").Deref(), []vm.Value{f_34, nid_35, arg__17521_47})
	if callErr != nil {
		return nil, callErr
	}
	v73 = v48
	f_74 = f_34
	nid_75 = nid_35
	action_76 = action_36
	goto b6
b5:
	;
	v57, callErr = rt.InvokeValue(vm.Keyword("swap-refs"), []vm.Value{action_39})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v57) {
		f_50 = f_37
		nid_51 = nid_38
		action_52 = action_39
		goto b7
	} else {
		f_53 = f_37
		nid_54 = nid_38
		action_55 = action_39
		goto b8
	}
b6:
	;
	v78 = v73
	f_79 = f_74
	nid_80 = nid_75
	action_81 = action_76
	goto b3
b7:
	;
	arg__17528_60, callErr = rt.InvokeValue(vm.Keyword("swap-refs"), []vm.Value{action_52})
	if callErr != nil {
		return nil, callErr
	}
	arg__17534_63, callErr = rt.InvokeValue(vm.Keyword("swap-refs"), []vm.Value{action_52})
	if callErr != nil {
		return nil, callErr
	}
	v64, callErr = rt.InvokeValue(rt.LookupVar("ir", "set-refs!").Deref(), []vm.Value{f_50, nid_51, arg__17534_63})
	if callErr != nil {
		return nil, callErr
	}
	v68 = v64
	f_69 = f_50
	nid_70 = nid_51
	action_71 = action_52
	goto b9
b8:
	;
	v68 = vm.NIL
	f_69 = f_53
	nid_70 = nid_54
	action_71 = action_55
	goto b9
b9:
	;
	v73 = v68
	f_74 = f_69
	nid_75 = nid_70
	action_76 = action_71
	goto b6
}
func try_identity(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var op_4 vm.Value
	var refs_6 vm.Value
	var arg__17550_17 vm.Value
	var v18 bool
	var nid_7 vm.Value
	var f_8 vm.Value
	var op_9 vm.Value
	var refs_10 vm.Value
	var r0_23 vm.Value
	var r1_27 vm.Value
	var v0_29 vm.Value
	var v1_31 vm.Value
	var v57 bool
	var nid_11 vm.Value
	var f_12 vm.Value
	var op_13 vm.Value
	var refs_14 vm.Value
	var v570 vm.Value
	var nid_571 vm.Value
	var f_572 vm.Value
	var op_573 vm.Value
	var refs_574 vm.Value
	var nid_34 vm.Value
	var f_35 vm.Value
	var case__17535_36 vm.Value
	var op_37 vm.Value
	var refs_38 vm.Value
	var r0_39 vm.Value
	var r1_40 vm.Value
	var v0_41 vm.Value
	var v1_42 vm.Value
	var zero_QMARK__43 vm.Value
	var one_QMARK__44 vm.Value
	var v81 vm.Value
	var nid_45 vm.Value
	var f_46 vm.Value
	var case__17535_47 vm.Value
	var op_48 vm.Value
	var refs_49 vm.Value
	var r0_50 vm.Value
	var r1_51 vm.Value
	var v0_52 vm.Value
	var v1_53 vm.Value
	var zero_QMARK__54 vm.Value
	var one_QMARK__55 vm.Value
	var v205 bool
	var v555 vm.Value
	var nid_556 vm.Value
	var f_557 vm.Value
	var case__17535_558 vm.Value
	var op_559 vm.Value
	var refs_560 vm.Value
	var r0_561 vm.Value
	var r1_562 vm.Value
	var v0_563 vm.Value
	var v1_564 vm.Value
	var zero_QMARK__565 vm.Value
	var one_QMARK__566 vm.Value
	var nid_59 vm.Value
	var f_60 vm.Value
	var case__17535_61 vm.Value
	var op_62 vm.Value
	var refs_63 vm.Value
	var r0_64 vm.Value
	var r1_65 vm.Value
	var v0_66 vm.Value
	var v1_67 vm.Value
	var zero_QMARK__68 vm.Value
	var one_QMARK__69 vm.Value
	var v85 vm.Value
	var nid_70 vm.Value
	var f_71 vm.Value
	var case__17535_72 vm.Value
	var op_73 vm.Value
	var refs_74 vm.Value
	var r0_75 vm.Value
	var r1_76 vm.Value
	var v0_77 vm.Value
	var v1_78 vm.Value
	var zero_QMARK__79 vm.Value
	var one_QMARK__80 vm.Value
	var v109 vm.Value
	var v169 vm.Value
	var nid_170 vm.Value
	var f_171 vm.Value
	var case__17535_172 vm.Value
	var op_173 vm.Value
	var refs_174 vm.Value
	var r0_175 vm.Value
	var r1_176 vm.Value
	var v0_177 vm.Value
	var v1_178 vm.Value
	var zero_QMARK__179 vm.Value
	var one_QMARK__180 vm.Value
	var nid_87 vm.Value
	var f_88 vm.Value
	var case__17535_89 vm.Value
	var op_90 vm.Value
	var refs_91 vm.Value
	var r0_92 vm.Value
	var r1_93 vm.Value
	var v0_94 vm.Value
	var v1_95 vm.Value
	var zero_QMARK__96 vm.Value
	var one_QMARK__97 vm.Value
	var v113 vm.Value
	var nid_98 vm.Value
	var f_99 vm.Value
	var case__17535_100 vm.Value
	var op_101 vm.Value
	var refs_102 vm.Value
	var r0_103 vm.Value
	var r1_104 vm.Value
	var v0_105 vm.Value
	var v1_106 vm.Value
	var zero_QMARK__107 vm.Value
	var one_QMARK__108 vm.Value
	var v156 vm.Value
	var nid_157 vm.Value
	var f_158 vm.Value
	var case__17535_159 vm.Value
	var op_160 vm.Value
	var refs_161 vm.Value
	var r0_162 vm.Value
	var r1_163 vm.Value
	var v0_164 vm.Value
	var v1_165 vm.Value
	var zero_QMARK__166 vm.Value
	var one_QMARK__167 vm.Value
	var nid_115 vm.Value
	var f_116 vm.Value
	var case__17535_117 vm.Value
	var op_118 vm.Value
	var refs_119 vm.Value
	var r0_120 vm.Value
	var r1_121 vm.Value
	var v0_122 vm.Value
	var v1_123 vm.Value
	var zero_QMARK__124 vm.Value
	var one_QMARK__125 vm.Value
	var nid_126 vm.Value
	var f_127 vm.Value
	var case__17535_128 vm.Value
	var op_129 vm.Value
	var refs_130 vm.Value
	var r0_131 vm.Value
	var r1_132 vm.Value
	var v0_133 vm.Value
	var v1_134 vm.Value
	var zero_QMARK__135 vm.Value
	var one_QMARK__136 vm.Value
	var v143 vm.Value
	var nid_144 vm.Value
	var f_145 vm.Value
	var case__17535_146 vm.Value
	var op_147 vm.Value
	var refs_148 vm.Value
	var r0_149 vm.Value
	var r1_150 vm.Value
	var v0_151 vm.Value
	var v1_152 vm.Value
	var zero_QMARK__153 vm.Value
	var one_QMARK__154 vm.Value
	var nid_182 vm.Value
	var f_183 vm.Value
	var case__17535_184 vm.Value
	var op_185 vm.Value
	var refs_186 vm.Value
	var r0_187 vm.Value
	var r1_188 vm.Value
	var v0_189 vm.Value
	var v1_190 vm.Value
	var zero_QMARK__191 vm.Value
	var one_QMARK__192 vm.Value
	var v229 vm.Value
	var nid_193 vm.Value
	var f_194 vm.Value
	var case__17535_195 vm.Value
	var op_196 vm.Value
	var refs_197 vm.Value
	var r0_198 vm.Value
	var r1_199 vm.Value
	var v0_200 vm.Value
	var v1_201 vm.Value
	var zero_QMARK__202 vm.Value
	var one_QMARK__203 vm.Value
	var v273 bool
	var v542 vm.Value
	var nid_543 vm.Value
	var f_544 vm.Value
	var case__17535_545 vm.Value
	var op_546 vm.Value
	var refs_547 vm.Value
	var r0_548 vm.Value
	var r1_549 vm.Value
	var v0_550 vm.Value
	var v1_551 vm.Value
	var zero_QMARK__552 vm.Value
	var one_QMARK__553 vm.Value
	var nid_207 vm.Value
	var f_208 vm.Value
	var case__17535_209 vm.Value
	var op_210 vm.Value
	var refs_211 vm.Value
	var r0_212 vm.Value
	var r1_213 vm.Value
	var v0_214 vm.Value
	var v1_215 vm.Value
	var zero_QMARK__216 vm.Value
	var one_QMARK__217 vm.Value
	var v233 vm.Value
	var nid_218 vm.Value
	var f_219 vm.Value
	var case__17535_220 vm.Value
	var op_221 vm.Value
	var refs_222 vm.Value
	var r0_223 vm.Value
	var r1_224 vm.Value
	var v0_225 vm.Value
	var v1_226 vm.Value
	var zero_QMARK__227 vm.Value
	var one_QMARK__228 vm.Value
	var v237 vm.Value
	var nid_238 vm.Value
	var f_239 vm.Value
	var case__17535_240 vm.Value
	var op_241 vm.Value
	var refs_242 vm.Value
	var r0_243 vm.Value
	var r1_244 vm.Value
	var v0_245 vm.Value
	var v1_246 vm.Value
	var zero_QMARK__247 vm.Value
	var one_QMARK__248 vm.Value
	var nid_250 vm.Value
	var f_251 vm.Value
	var case__17535_252 vm.Value
	var op_253 vm.Value
	var refs_254 vm.Value
	var r0_255 vm.Value
	var r1_256 vm.Value
	var v0_257 vm.Value
	var v1_258 vm.Value
	var zero_QMARK__259 vm.Value
	var one_QMARK__260 vm.Value
	var v297 vm.Value
	var nid_261 vm.Value
	var f_262 vm.Value
	var case__17535_263 vm.Value
	var op_264 vm.Value
	var refs_265 vm.Value
	var r0_266 vm.Value
	var r1_267 vm.Value
	var v0_268 vm.Value
	var v1_269 vm.Value
	var zero_QMARK__270 vm.Value
	var one_QMARK__271 vm.Value
	var v529 vm.Value
	var nid_530 vm.Value
	var f_531 vm.Value
	var case__17535_532 vm.Value
	var op_533 vm.Value
	var refs_534 vm.Value
	var r0_535 vm.Value
	var r1_536 vm.Value
	var v0_537 vm.Value
	var v1_538 vm.Value
	var zero_QMARK__539 vm.Value
	var one_QMARK__540 vm.Value
	var nid_275 vm.Value
	var f_276 vm.Value
	var case__17535_277 vm.Value
	var op_278 vm.Value
	var refs_279 vm.Value
	var r0_280 vm.Value
	var r1_281 vm.Value
	var v0_282 vm.Value
	var v1_283 vm.Value
	var zero_QMARK__284 vm.Value
	var one_QMARK__285 vm.Value
	var v301 vm.Value
	var nid_286 vm.Value
	var f_287 vm.Value
	var case__17535_288 vm.Value
	var op_289 vm.Value
	var refs_290 vm.Value
	var r0_291 vm.Value
	var r1_292 vm.Value
	var v0_293 vm.Value
	var v1_294 vm.Value
	var zero_QMARK__295 vm.Value
	var one_QMARK__296 vm.Value
	var v325 vm.Value
	var v475 vm.Value
	var nid_476 vm.Value
	var f_477 vm.Value
	var case__17535_478 vm.Value
	var op_479 vm.Value
	var refs_480 vm.Value
	var r0_481 vm.Value
	var r1_482 vm.Value
	var v0_483 vm.Value
	var v1_484 vm.Value
	var zero_QMARK__485 vm.Value
	var one_QMARK__486 vm.Value
	var nid_303 vm.Value
	var f_304 vm.Value
	var case__17535_305 vm.Value
	var op_306 vm.Value
	var refs_307 vm.Value
	var r0_308 vm.Value
	var r1_309 vm.Value
	var v0_310 vm.Value
	var v1_311 vm.Value
	var zero_QMARK__312 vm.Value
	var one_QMARK__313 vm.Value
	var v329 vm.Value
	var nid_314 vm.Value
	var f_315 vm.Value
	var case__17535_316 vm.Value
	var op_317 vm.Value
	var refs_318 vm.Value
	var r0_319 vm.Value
	var r1_320 vm.Value
	var v0_321 vm.Value
	var v1_322 vm.Value
	var zero_QMARK__323 vm.Value
	var one_QMARK__324 vm.Value
	var or__x_353 vm.Value
	var v462 vm.Value
	var nid_463 vm.Value
	var f_464 vm.Value
	var case__17535_465 vm.Value
	var op_466 vm.Value
	var refs_467 vm.Value
	var r0_468 vm.Value
	var r1_469 vm.Value
	var v0_470 vm.Value
	var v1_471 vm.Value
	var zero_QMARK__472 vm.Value
	var one_QMARK__473 vm.Value
	var nid_331 vm.Value
	var f_332 vm.Value
	var case__17535_333 vm.Value
	var op_334 vm.Value
	var refs_335 vm.Value
	var r0_336 vm.Value
	var r1_337 vm.Value
	var v0_338 vm.Value
	var v1_339 vm.Value
	var zero_QMARK__340 vm.Value
	var one_QMARK__341 vm.Value
	var arg__17632_405 vm.Value
	var v406 vm.Value
	var nid_342 vm.Value
	var f_343 vm.Value
	var case__17535_344 vm.Value
	var op_345 vm.Value
	var refs_346 vm.Value
	var r0_347 vm.Value
	var r1_348 vm.Value
	var v0_349 vm.Value
	var v1_350 vm.Value
	var zero_QMARK__351 vm.Value
	var one_QMARK__352 vm.Value
	var v449 vm.Value
	var nid_450 vm.Value
	var f_451 vm.Value
	var case__17535_452 vm.Value
	var op_453 vm.Value
	var refs_454 vm.Value
	var r0_455 vm.Value
	var r1_456 vm.Value
	var v0_457 vm.Value
	var v1_458 vm.Value
	var zero_QMARK__459 vm.Value
	var one_QMARK__460 vm.Value
	var nid_354 vm.Value
	var f_355 vm.Value
	var case__17535_356 vm.Value
	var op_357 vm.Value
	var refs_358 vm.Value
	var r0_359 vm.Value
	var r1_360 vm.Value
	var v0_361 vm.Value
	var v1_362 vm.Value
	var zero_QMARK__363 vm.Value
	var one_QMARK__364 vm.Value
	var or__x_365 vm.Value
	var nid_366 vm.Value
	var f_367 vm.Value
	var case__17535_368 vm.Value
	var op_369 vm.Value
	var refs_370 vm.Value
	var r0_371 vm.Value
	var r1_372 vm.Value
	var v0_373 vm.Value
	var v1_374 vm.Value
	var zero_QMARK__375 vm.Value
	var one_QMARK__376 vm.Value
	var or__x_377 vm.Value
	var v380 vm.Value
	var v382 vm.Value
	var nid_383 vm.Value
	var f_384 vm.Value
	var case__17535_385 vm.Value
	var op_386 vm.Value
	var refs_387 vm.Value
	var r0_388 vm.Value
	var r1_389 vm.Value
	var v0_390 vm.Value
	var v1_391 vm.Value
	var zero_QMARK__392 vm.Value
	var one_QMARK__393 vm.Value
	var or__x_394 vm.Value
	var nid_408 vm.Value
	var f_409 vm.Value
	var case__17535_410 vm.Value
	var op_411 vm.Value
	var refs_412 vm.Value
	var r0_413 vm.Value
	var r1_414 vm.Value
	var v0_415 vm.Value
	var v1_416 vm.Value
	var zero_QMARK__417 vm.Value
	var one_QMARK__418 vm.Value
	var nid_419 vm.Value
	var f_420 vm.Value
	var case__17535_421 vm.Value
	var op_422 vm.Value
	var refs_423 vm.Value
	var r0_424 vm.Value
	var r1_425 vm.Value
	var v0_426 vm.Value
	var v1_427 vm.Value
	var zero_QMARK__428 vm.Value
	var one_QMARK__429 vm.Value
	var v436 vm.Value
	var nid_437 vm.Value
	var f_438 vm.Value
	var case__17535_439 vm.Value
	var op_440 vm.Value
	var refs_441 vm.Value
	var r0_442 vm.Value
	var r1_443 vm.Value
	var v0_444 vm.Value
	var v1_445 vm.Value
	var zero_QMARK__446 vm.Value
	var one_QMARK__447 vm.Value
	var nid_488 vm.Value
	var f_489 vm.Value
	var case__17535_490 vm.Value
	var op_491 vm.Value
	var refs_492 vm.Value
	var r0_493 vm.Value
	var r1_494 vm.Value
	var v0_495 vm.Value
	var v1_496 vm.Value
	var zero_QMARK__497 vm.Value
	var one_QMARK__498 vm.Value
	var nid_499 vm.Value
	var f_500 vm.Value
	var case__17535_501 vm.Value
	var op_502 vm.Value
	var refs_503 vm.Value
	var r0_504 vm.Value
	var r1_505 vm.Value
	var v0_506 vm.Value
	var v1_507 vm.Value
	var zero_QMARK__508 vm.Value
	var one_QMARK__509 vm.Value
	var v516 vm.Value
	var nid_517 vm.Value
	var f_518 vm.Value
	var case__17535_519 vm.Value
	var op_520 vm.Value
	var refs_521 vm.Value
	var r0_522 vm.Value
	var r1_523 vm.Value
	var v0_524 vm.Value
	var v1_525 vm.Value
	var zero_QMARK__526 vm.Value
	var one_QMARK__527 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = op_4, refs_6, arg__17550_17, v18, nid_7, f_8, op_9, refs_10, r0_23, r1_27, v0_29, v1_31, v57, nid_11, f_12, op_13, refs_14, v570, nid_571, f_572, op_573, refs_574, nid_34, f_35, case__17535_36, op_37, refs_38, r0_39, r1_40, v0_41, v1_42, zero_QMARK__43, one_QMARK__44, v81, nid_45, f_46, case__17535_47, op_48, refs_49, r0_50, r1_51, v0_52, v1_53, zero_QMARK__54, one_QMARK__55, v205, v555, nid_556, f_557, case__17535_558, op_559, refs_560, r0_561, r1_562, v0_563, v1_564, zero_QMARK__565, one_QMARK__566, nid_59, f_60, case__17535_61, op_62, refs_63, r0_64, r1_65, v0_66, v1_67, zero_QMARK__68, one_QMARK__69, v85, nid_70, f_71, case__17535_72, op_73, refs_74, r0_75, r1_76, v0_77, v1_78, zero_QMARK__79, one_QMARK__80, v109, v169, nid_170, f_171, case__17535_172, op_173, refs_174, r0_175, r1_176, v0_177, v1_178, zero_QMARK__179, one_QMARK__180, nid_87, f_88, case__17535_89, op_90, refs_91, r0_92, r1_93, v0_94, v1_95, zero_QMARK__96, one_QMARK__97, v113, nid_98, f_99, case__17535_100, op_101, refs_102, r0_103, r1_104, v0_105, v1_106, zero_QMARK__107, one_QMARK__108, v156, nid_157, f_158, case__17535_159, op_160, refs_161, r0_162, r1_163, v0_164, v1_165, zero_QMARK__166, one_QMARK__167, nid_115, f_116, case__17535_117, op_118, refs_119, r0_120, r1_121, v0_122, v1_123, zero_QMARK__124, one_QMARK__125, nid_126, f_127, case__17535_128, op_129, refs_130, r0_131, r1_132, v0_133, v1_134, zero_QMARK__135, one_QMARK__136, v143, nid_144, f_145, case__17535_146, op_147, refs_148, r0_149, r1_150, v0_151, v1_152, zero_QMARK__153, one_QMARK__154, nid_182, f_183, case__17535_184, op_185, refs_186, r0_187, r1_188, v0_189, v1_190, zero_QMARK__191, one_QMARK__192, v229, nid_193, f_194, case__17535_195, op_196, refs_197, r0_198, r1_199, v0_200, v1_201, zero_QMARK__202, one_QMARK__203, v273, v542, nid_543, f_544, case__17535_545, op_546, refs_547, r0_548, r1_549, v0_550, v1_551, zero_QMARK__552, one_QMARK__553, nid_207, f_208, case__17535_209, op_210, refs_211, r0_212, r1_213, v0_214, v1_215, zero_QMARK__216, one_QMARK__217, v233, nid_218, f_219, case__17535_220, op_221, refs_222, r0_223, r1_224, v0_225, v1_226, zero_QMARK__227, one_QMARK__228, v237, nid_238, f_239, case__17535_240, op_241, refs_242, r0_243, r1_244, v0_245, v1_246, zero_QMARK__247, one_QMARK__248, nid_250, f_251, case__17535_252, op_253, refs_254, r0_255, r1_256, v0_257, v1_258, zero_QMARK__259, one_QMARK__260, v297, nid_261, f_262, case__17535_263, op_264, refs_265, r0_266, r1_267, v0_268, v1_269, zero_QMARK__270, one_QMARK__271, v529, nid_530, f_531, case__17535_532, op_533, refs_534, r0_535, r1_536, v0_537, v1_538, zero_QMARK__539, one_QMARK__540, nid_275, f_276, case__17535_277, op_278, refs_279, r0_280, r1_281, v0_282, v1_283, zero_QMARK__284, one_QMARK__285, v301, nid_286, f_287, case__17535_288, op_289, refs_290, r0_291, r1_292, v0_293, v1_294, zero_QMARK__295, one_QMARK__296, v325, v475, nid_476, f_477, case__17535_478, op_479, refs_480, r0_481, r1_482, v0_483, v1_484, zero_QMARK__485, one_QMARK__486, nid_303, f_304, case__17535_305, op_306, refs_307, r0_308, r1_309, v0_310, v1_311, zero_QMARK__312, one_QMARK__313, v329, nid_314, f_315, case__17535_316, op_317, refs_318, r0_319, r1_320, v0_321, v1_322, zero_QMARK__323, one_QMARK__324, or__x_353, v462, nid_463, f_464, case__17535_465, op_466, refs_467, r0_468, r1_469, v0_470, v1_471, zero_QMARK__472, one_QMARK__473, nid_331, f_332, case__17535_333, op_334, refs_335, r0_336, r1_337, v0_338, v1_339, zero_QMARK__340, one_QMARK__341, arg__17632_405, v406, nid_342, f_343, case__17535_344, op_345, refs_346, r0_347, r1_348, v0_349, v1_350, zero_QMARK__351, one_QMARK__352, v449, nid_450, f_451, case__17535_452, op_453, refs_454, r0_455, r1_456, v0_457, v1_458, zero_QMARK__459, one_QMARK__460, nid_354, f_355, case__17535_356, op_357, refs_358, r0_359, r1_360, v0_361, v1_362, zero_QMARK__363, one_QMARK__364, or__x_365, nid_366, f_367, case__17535_368, op_369, refs_370, r0_371, r1_372, v0_373, v1_374, zero_QMARK__375, one_QMARK__376, or__x_377, v380, v382, nid_383, f_384, case__17535_385, op_386, refs_387, r0_388, r1_389, v0_390, v1_391, zero_QMARK__392, one_QMARK__393, or__x_394, nid_408, f_409, case__17535_410, op_411, refs_412, r0_413, r1_414, v0_415, v1_416, zero_QMARK__417, one_QMARK__418, nid_419, f_420, case__17535_421, op_422, refs_423, r0_424, r1_425, v0_426, v1_427, zero_QMARK__428, one_QMARK__429, v436, nid_437, f_438, case__17535_439, op_440, refs_441, r0_442, r1_443, v0_444, v1_445, zero_QMARK__446, one_QMARK__447, nid_488, f_489, case__17535_490, op_491, refs_492, r0_493, r1_494, v0_495, v1_496, zero_QMARK__497, one_QMARK__498, nid_499, f_500, case__17535_501, op_502, refs_503, r0_504, r1_505, v0_506, v1_507, zero_QMARK__508, one_QMARK__509, v516, nid_517, f_518, case__17535_519, op_520, refs_521, r0_522, r1_523, v0_524, v1_525, zero_QMARK__526, one_QMARK__527
	op_4, callErr = rt.InvokeValue(rt.LookupVar("ir", "op").Deref(), []vm.Value{arg0, arg1})
	if callErr != nil {
		return nil, callErr
	}
	refs_6, callErr = rt.InvokeValue(rt.LookupVar("ir", "refs").Deref(), []vm.Value{arg0, arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__17550_17, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{refs_6})
	if callErr != nil {
		return nil, callErr
	}
	v18 = arg__17550_17 == vm.Int(2)
	if v18 {
		nid_7 = arg0
		f_8 = arg1
		op_9 = op_4
		refs_10 = refs_6
		goto b1
	} else {
		nid_11 = arg0
		f_12 = arg1
		op_13 = op_4
		refs_14 = refs_6
		goto b2
	}
b1:
	;
	r0_23, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{refs_10, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	r1_27, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{refs_10, vm.Int(1)})
	if callErr != nil {
		return nil, callErr
	}
	v0_29, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.constfold", "const-val").Deref(), []vm.Value{r0_23, f_8})
	if callErr != nil {
		return nil, callErr
	}
	v1_31, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.constfold", "const-val").Deref(), []vm.Value{r1_27, f_8})
	if callErr != nil {
		return nil, callErr
	}
	v57 = op_9 == vm.Keyword("add")
	if v57 {
		nid_34 = nid_7
		f_35 = f_8
		case__17535_36 = op_9
		op_37 = op_9
		refs_38 = refs_10
		r0_39 = r0_23
		r1_40 = r1_27
		v0_41 = v0_29
		v1_42 = v1_31
		zero_QMARK__43 = rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
			var and__x_2 vm.Value
			var v_3 vm.Value
			var and__x_4 vm.Value
			var v9 vm.Value
			var v_5 vm.Value
			var and__x_6 vm.Value
			var v12 vm.Value
			var v_13 vm.Value
			var and__x_14 vm.Value
			var callErr error
			_, _, _, _, _, _, _, _, _ = and__x_2, v_3, and__x_4, v9, v_5, and__x_6, v12, v_13, and__x_14
			and__x_2, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "number?").Deref(), []vm.Value{arg0})
			if callErr != nil {
				return nil, callErr
			}
			if vm.IsTruthy(and__x_2) {
				v_3 = arg0
				and__x_4 = and__x_2
				goto b1
			} else {
				v_5 = arg0
				and__x_6 = and__x_2
				goto b2
			}
		b1:
			;
			v9 = vm.Boolean(vm.Int(0) == v_3)
			v12 = v9
			v_13 = v_3
			and__x_14 = and__x_4
			goto b3
		b2:
			;
			v12 = and__x_6
			v_13 = v_5
			and__x_14 = and__x_6
			goto b3
		b3:
			;
			return v12, nil
		})
		one_QMARK__44 = rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
			var and__x_2 vm.Value
			var v_3 vm.Value
			var and__x_4 vm.Value
			var v9 vm.Value
			var v_5 vm.Value
			var and__x_6 vm.Value
			var v12 vm.Value
			var v_13 vm.Value
			var and__x_14 vm.Value
			var callErr error
			_, _, _, _, _, _, _, _, _ = and__x_2, v_3, and__x_4, v9, v_5, and__x_6, v12, v_13, and__x_14
			and__x_2, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "number?").Deref(), []vm.Value{arg0})
			if callErr != nil {
				return nil, callErr
			}
			if vm.IsTruthy(and__x_2) {
				v_3 = arg0
				and__x_4 = and__x_2
				goto b1
			} else {
				v_5 = arg0
				and__x_6 = and__x_2
				goto b2
			}
		b1:
			;
			v9 = vm.Boolean(vm.Int(1) == v_3)
			v12 = v9
			v_13 = v_3
			and__x_14 = and__x_4
			goto b3
		b2:
			;
			v12 = and__x_6
			v_13 = v_5
			and__x_14 = and__x_6
			goto b3
		b3:
			;
			return v12, nil
		})
		goto b4
	} else {
		nid_45 = nid_7
		f_46 = f_8
		case__17535_47 = op_9
		op_48 = op_9
		refs_49 = refs_10
		r0_50 = r0_23
		r1_51 = r1_27
		v0_52 = v0_29
		v1_53 = v1_31
		zero_QMARK__54 = rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
			var and__x_2 vm.Value
			var v_3 vm.Value
			var and__x_4 vm.Value
			var v9 vm.Value
			var v_5 vm.Value
			var and__x_6 vm.Value
			var v12 vm.Value
			var v_13 vm.Value
			var and__x_14 vm.Value
			var callErr error
			_, _, _, _, _, _, _, _, _ = and__x_2, v_3, and__x_4, v9, v_5, and__x_6, v12, v_13, and__x_14
			and__x_2, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "number?").Deref(), []vm.Value{arg0})
			if callErr != nil {
				return nil, callErr
			}
			if vm.IsTruthy(and__x_2) {
				v_3 = arg0
				and__x_4 = and__x_2
				goto b1
			} else {
				v_5 = arg0
				and__x_6 = and__x_2
				goto b2
			}
		b1:
			;
			v9 = vm.Boolean(vm.Int(0) == v_3)
			v12 = v9
			v_13 = v_3
			and__x_14 = and__x_4
			goto b3
		b2:
			;
			v12 = and__x_6
			v_13 = v_5
			and__x_14 = and__x_6
			goto b3
		b3:
			;
			return v12, nil
		})
		one_QMARK__55 = rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
			var and__x_2 vm.Value
			var v_3 vm.Value
			var and__x_4 vm.Value
			var v9 vm.Value
			var v_5 vm.Value
			var and__x_6 vm.Value
			var v12 vm.Value
			var v_13 vm.Value
			var and__x_14 vm.Value
			var callErr error
			_, _, _, _, _, _, _, _, _ = and__x_2, v_3, and__x_4, v9, v_5, and__x_6, v12, v_13, and__x_14
			and__x_2, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "number?").Deref(), []vm.Value{arg0})
			if callErr != nil {
				return nil, callErr
			}
			if vm.IsTruthy(and__x_2) {
				v_3 = arg0
				and__x_4 = and__x_2
				goto b1
			} else {
				v_5 = arg0
				and__x_6 = and__x_2
				goto b2
			}
		b1:
			;
			v9 = vm.Boolean(vm.Int(1) == v_3)
			v12 = v9
			v_13 = v_3
			and__x_14 = and__x_4
			goto b3
		b2:
			;
			v12 = and__x_6
			v_13 = v_5
			and__x_14 = and__x_6
			goto b3
		b3:
			;
			return v12, nil
		})
		goto b5
	}
b2:
	;
	v570 = vm.NIL
	nid_571 = nid_11
	f_572 = f_12
	op_573 = op_13
	refs_574 = refs_14
	goto b3
b3:
	;
	return v570, nil
b4:
	;
	v81, callErr = rt.InvokeValue(zero_QMARK__43, []vm.Value{v1_42})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v81) {
		nid_59 = nid_34
		f_60 = f_35
		case__17535_61 = case__17535_36
		op_62 = op_37
		refs_63 = refs_38
		r0_64 = r0_39
		r1_65 = r1_40
		v0_66 = v0_41
		v1_67 = v1_42
		zero_QMARK__68 = zero_QMARK__43
		one_QMARK__69 = one_QMARK__44
		goto b7
	} else {
		nid_70 = nid_34
		f_71 = f_35
		case__17535_72 = case__17535_36
		op_73 = op_37
		refs_74 = refs_38
		r0_75 = r0_39
		r1_76 = r1_40
		v0_77 = v0_41
		v1_78 = v1_42
		zero_QMARK__79 = zero_QMARK__43
		one_QMARK__80 = one_QMARK__44
		goto b8
	}
b5:
	;
	v205 = case__17535_47 == vm.Keyword("sub")
	if v205 {
		nid_182 = nid_45
		f_183 = f_46
		case__17535_184 = case__17535_47
		op_185 = op_48
		refs_186 = refs_49
		r0_187 = r0_50
		r1_188 = r1_51
		v0_189 = v0_52
		v1_190 = v1_53
		zero_QMARK__191 = zero_QMARK__54
		one_QMARK__192 = one_QMARK__55
		goto b16
	} else {
		nid_193 = nid_45
		f_194 = f_46
		case__17535_195 = case__17535_47
		op_196 = op_48
		refs_197 = refs_49
		r0_198 = r0_50
		r1_199 = r1_51
		v0_200 = v0_52
		v1_201 = v1_53
		zero_QMARK__202 = zero_QMARK__54
		one_QMARK__203 = one_QMARK__55
		goto b17
	}
b6:
	;
	v570 = v555
	nid_571 = nid_556
	f_572 = f_557
	op_573 = op_559
	refs_574 = refs_560
	goto b3
b7:
	;
	v85, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("replace-uses"), r0_64})
	if callErr != nil {
		return nil, callErr
	}
	v169 = v85
	nid_170 = nid_59
	f_171 = f_60
	case__17535_172 = case__17535_61
	op_173 = op_62
	refs_174 = refs_63
	r0_175 = r0_64
	r1_176 = r1_65
	v0_177 = v0_66
	v1_178 = v1_67
	zero_QMARK__179 = zero_QMARK__68
	one_QMARK__180 = one_QMARK__69
	goto b9
b8:
	;
	v109, callErr = rt.InvokeValue(zero_QMARK__79, []vm.Value{v0_77})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v109) {
		nid_87 = nid_70
		f_88 = f_71
		case__17535_89 = case__17535_72
		op_90 = op_73
		refs_91 = refs_74
		r0_92 = r0_75
		r1_93 = r1_76
		v0_94 = v0_77
		v1_95 = v1_78
		zero_QMARK__96 = zero_QMARK__79
		one_QMARK__97 = one_QMARK__80
		goto b10
	} else {
		nid_98 = nid_70
		f_99 = f_71
		case__17535_100 = case__17535_72
		op_101 = op_73
		refs_102 = refs_74
		r0_103 = r0_75
		r1_104 = r1_76
		v0_105 = v0_77
		v1_106 = v1_78
		zero_QMARK__107 = zero_QMARK__79
		one_QMARK__108 = one_QMARK__80
		goto b11
	}
b9:
	;
	v555 = v169
	nid_556 = nid_170
	f_557 = f_171
	case__17535_558 = case__17535_172
	op_559 = op_173
	refs_560 = refs_174
	r0_561 = r0_175
	r1_562 = r1_176
	v0_563 = v0_177
	v1_564 = v1_178
	zero_QMARK__565 = zero_QMARK__179
	one_QMARK__566 = one_QMARK__180
	goto b6
b10:
	;
	v113, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("replace-uses"), r1_93})
	if callErr != nil {
		return nil, callErr
	}
	v156 = v113
	nid_157 = nid_87
	f_158 = f_88
	case__17535_159 = case__17535_89
	op_160 = op_90
	refs_161 = refs_91
	r0_162 = r0_92
	r1_163 = r1_93
	v0_164 = v0_94
	v1_165 = v1_95
	zero_QMARK__166 = zero_QMARK__96
	one_QMARK__167 = one_QMARK__97
	goto b12
b11:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		nid_115 = nid_98
		f_116 = f_99
		case__17535_117 = case__17535_100
		op_118 = op_101
		refs_119 = refs_102
		r0_120 = r0_103
		r1_121 = r1_104
		v0_122 = v0_105
		v1_123 = v1_106
		zero_QMARK__124 = zero_QMARK__107
		one_QMARK__125 = one_QMARK__108
		goto b13
	} else {
		nid_126 = nid_98
		f_127 = f_99
		case__17535_128 = case__17535_100
		op_129 = op_101
		refs_130 = refs_102
		r0_131 = r0_103
		r1_132 = r1_104
		v0_133 = v0_105
		v1_134 = v1_106
		zero_QMARK__135 = zero_QMARK__107
		one_QMARK__136 = one_QMARK__108
		goto b14
	}
b12:
	;
	v169 = v156
	nid_170 = nid_157
	f_171 = f_158
	case__17535_172 = case__17535_159
	op_173 = op_160
	refs_174 = refs_161
	r0_175 = r0_162
	r1_176 = r1_163
	v0_177 = v0_164
	v1_178 = v1_165
	zero_QMARK__179 = zero_QMARK__166
	one_QMARK__180 = one_QMARK__167
	goto b9
b13:
	;
	v143 = vm.NIL
	nid_144 = nid_115
	f_145 = f_116
	case__17535_146 = case__17535_117
	op_147 = op_118
	refs_148 = refs_119
	r0_149 = r0_120
	r1_150 = r1_121
	v0_151 = v0_122
	v1_152 = v1_123
	zero_QMARK__153 = zero_QMARK__124
	one_QMARK__154 = one_QMARK__125
	goto b15
b14:
	;
	v143 = vm.NIL
	nid_144 = nid_126
	f_145 = f_127
	case__17535_146 = case__17535_128
	op_147 = op_129
	refs_148 = refs_130
	r0_149 = r0_131
	r1_150 = r1_132
	v0_151 = v0_133
	v1_152 = v1_134
	zero_QMARK__153 = zero_QMARK__135
	one_QMARK__154 = one_QMARK__136
	goto b15
b15:
	;
	v156 = v143
	nid_157 = nid_144
	f_158 = f_145
	case__17535_159 = case__17535_146
	op_160 = op_147
	refs_161 = refs_148
	r0_162 = r0_149
	r1_163 = r1_150
	v0_164 = v0_151
	v1_165 = v1_152
	zero_QMARK__166 = zero_QMARK__153
	one_QMARK__167 = one_QMARK__154
	goto b12
b16:
	;
	v229, callErr = rt.InvokeValue(zero_QMARK__191, []vm.Value{v1_190})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v229) {
		nid_207 = nid_182
		f_208 = f_183
		case__17535_209 = case__17535_184
		op_210 = op_185
		refs_211 = refs_186
		r0_212 = r0_187
		r1_213 = r1_188
		v0_214 = v0_189
		v1_215 = v1_190
		zero_QMARK__216 = zero_QMARK__191
		one_QMARK__217 = one_QMARK__192
		goto b19
	} else {
		nid_218 = nid_182
		f_219 = f_183
		case__17535_220 = case__17535_184
		op_221 = op_185
		refs_222 = refs_186
		r0_223 = r0_187
		r1_224 = r1_188
		v0_225 = v0_189
		v1_226 = v1_190
		zero_QMARK__227 = zero_QMARK__191
		one_QMARK__228 = one_QMARK__192
		goto b20
	}
b17:
	;
	v273 = case__17535_195 == vm.Keyword("mul")
	if v273 {
		nid_250 = nid_193
		f_251 = f_194
		case__17535_252 = case__17535_195
		op_253 = op_196
		refs_254 = refs_197
		r0_255 = r0_198
		r1_256 = r1_199
		v0_257 = v0_200
		v1_258 = v1_201
		zero_QMARK__259 = zero_QMARK__202
		one_QMARK__260 = one_QMARK__203
		goto b22
	} else {
		nid_261 = nid_193
		f_262 = f_194
		case__17535_263 = case__17535_195
		op_264 = op_196
		refs_265 = refs_197
		r0_266 = r0_198
		r1_267 = r1_199
		v0_268 = v0_200
		v1_269 = v1_201
		zero_QMARK__270 = zero_QMARK__202
		one_QMARK__271 = one_QMARK__203
		goto b23
	}
b18:
	;
	v555 = v542
	nid_556 = nid_543
	f_557 = f_544
	case__17535_558 = case__17535_545
	op_559 = op_546
	refs_560 = refs_547
	r0_561 = r0_548
	r1_562 = r1_549
	v0_563 = v0_550
	v1_564 = v1_551
	zero_QMARK__565 = zero_QMARK__552
	one_QMARK__566 = one_QMARK__553
	goto b6
b19:
	;
	v233, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("replace-uses"), r0_212})
	if callErr != nil {
		return nil, callErr
	}
	v237 = v233
	nid_238 = nid_207
	f_239 = f_208
	case__17535_240 = case__17535_209
	op_241 = op_210
	refs_242 = refs_211
	r0_243 = r0_212
	r1_244 = r1_213
	v0_245 = v0_214
	v1_246 = v1_215
	zero_QMARK__247 = zero_QMARK__216
	one_QMARK__248 = one_QMARK__217
	goto b21
b20:
	;
	v237 = vm.NIL
	nid_238 = nid_218
	f_239 = f_219
	case__17535_240 = case__17535_220
	op_241 = op_221
	refs_242 = refs_222
	r0_243 = r0_223
	r1_244 = r1_224
	v0_245 = v0_225
	v1_246 = v1_226
	zero_QMARK__247 = zero_QMARK__227
	one_QMARK__248 = one_QMARK__228
	goto b21
b21:
	;
	v542 = v237
	nid_543 = nid_238
	f_544 = f_239
	case__17535_545 = case__17535_240
	op_546 = op_241
	refs_547 = refs_242
	r0_548 = r0_243
	r1_549 = r1_244
	v0_550 = v0_245
	v1_551 = v1_246
	zero_QMARK__552 = zero_QMARK__247
	one_QMARK__553 = one_QMARK__248
	goto b18
b22:
	;
	v297, callErr = rt.InvokeValue(one_QMARK__260, []vm.Value{v1_258})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v297) {
		nid_275 = nid_250
		f_276 = f_251
		case__17535_277 = case__17535_252
		op_278 = op_253
		refs_279 = refs_254
		r0_280 = r0_255
		r1_281 = r1_256
		v0_282 = v0_257
		v1_283 = v1_258
		zero_QMARK__284 = zero_QMARK__259
		one_QMARK__285 = one_QMARK__260
		goto b25
	} else {
		nid_286 = nid_250
		f_287 = f_251
		case__17535_288 = case__17535_252
		op_289 = op_253
		refs_290 = refs_254
		r0_291 = r0_255
		r1_292 = r1_256
		v0_293 = v0_257
		v1_294 = v1_258
		zero_QMARK__295 = zero_QMARK__259
		one_QMARK__296 = one_QMARK__260
		goto b26
	}
b23:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		nid_488 = nid_261
		f_489 = f_262
		case__17535_490 = case__17535_263
		op_491 = op_264
		refs_492 = refs_265
		r0_493 = r0_266
		r1_494 = r1_267
		v0_495 = v0_268
		v1_496 = v1_269
		zero_QMARK__497 = zero_QMARK__270
		one_QMARK__498 = one_QMARK__271
		goto b40
	} else {
		nid_499 = nid_261
		f_500 = f_262
		case__17535_501 = case__17535_263
		op_502 = op_264
		refs_503 = refs_265
		r0_504 = r0_266
		r1_505 = r1_267
		v0_506 = v0_268
		v1_507 = v1_269
		zero_QMARK__508 = zero_QMARK__270
		one_QMARK__509 = one_QMARK__271
		goto b41
	}
b24:
	;
	v542 = v529
	nid_543 = nid_530
	f_544 = f_531
	case__17535_545 = case__17535_532
	op_546 = op_533
	refs_547 = refs_534
	r0_548 = r0_535
	r1_549 = r1_536
	v0_550 = v0_537
	v1_551 = v1_538
	zero_QMARK__552 = zero_QMARK__539
	one_QMARK__553 = one_QMARK__540
	goto b18
b25:
	;
	v301, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("replace-uses"), r0_280})
	if callErr != nil {
		return nil, callErr
	}
	v475 = v301
	nid_476 = nid_275
	f_477 = f_276
	case__17535_478 = case__17535_277
	op_479 = op_278
	refs_480 = refs_279
	r0_481 = r0_280
	r1_482 = r1_281
	v0_483 = v0_282
	v1_484 = v1_283
	zero_QMARK__485 = zero_QMARK__284
	one_QMARK__486 = one_QMARK__285
	goto b27
b26:
	;
	v325, callErr = rt.InvokeValue(one_QMARK__296, []vm.Value{v0_293})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v325) {
		nid_303 = nid_286
		f_304 = f_287
		case__17535_305 = case__17535_288
		op_306 = op_289
		refs_307 = refs_290
		r0_308 = r0_291
		r1_309 = r1_292
		v0_310 = v0_293
		v1_311 = v1_294
		zero_QMARK__312 = zero_QMARK__295
		one_QMARK__313 = one_QMARK__296
		goto b28
	} else {
		nid_314 = nid_286
		f_315 = f_287
		case__17535_316 = case__17535_288
		op_317 = op_289
		refs_318 = refs_290
		r0_319 = r0_291
		r1_320 = r1_292
		v0_321 = v0_293
		v1_322 = v1_294
		zero_QMARK__323 = zero_QMARK__295
		one_QMARK__324 = one_QMARK__296
		goto b29
	}
b27:
	;
	v529 = v475
	nid_530 = nid_476
	f_531 = f_477
	case__17535_532 = case__17535_478
	op_533 = op_479
	refs_534 = refs_480
	r0_535 = r0_481
	r1_536 = r1_482
	v0_537 = v0_483
	v1_538 = v1_484
	zero_QMARK__539 = zero_QMARK__485
	one_QMARK__540 = one_QMARK__486
	goto b24
b28:
	;
	v329, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("replace-uses"), r1_309})
	if callErr != nil {
		return nil, callErr
	}
	v462 = v329
	nid_463 = nid_303
	f_464 = f_304
	case__17535_465 = case__17535_305
	op_466 = op_306
	refs_467 = refs_307
	r0_468 = r0_308
	r1_469 = r1_309
	v0_470 = v0_310
	v1_471 = v1_311
	zero_QMARK__472 = zero_QMARK__312
	one_QMARK__473 = one_QMARK__313
	goto b30
b29:
	;
	or__x_353, callErr = rt.InvokeValue(zero_QMARK__323, []vm.Value{v0_321})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(or__x_353) {
		nid_354 = nid_314
		f_355 = f_315
		case__17535_356 = case__17535_316
		op_357 = op_317
		refs_358 = refs_318
		r0_359 = r0_319
		r1_360 = r1_320
		v0_361 = v0_321
		v1_362 = v1_322
		zero_QMARK__363 = zero_QMARK__323
		one_QMARK__364 = one_QMARK__324
		or__x_365 = or__x_353
		goto b34
	} else {
		nid_366 = nid_314
		f_367 = f_315
		case__17535_368 = case__17535_316
		op_369 = op_317
		refs_370 = refs_318
		r0_371 = r0_319
		r1_372 = r1_320
		v0_373 = v0_321
		v1_374 = v1_322
		zero_QMARK__375 = zero_QMARK__323
		one_QMARK__376 = one_QMARK__324
		or__x_377 = or__x_353
		goto b35
	}
b30:
	;
	v475 = v462
	nid_476 = nid_463
	f_477 = f_464
	case__17535_478 = case__17535_465
	op_479 = op_466
	refs_480 = refs_467
	r0_481 = r0_468
	r1_482 = r1_469
	v0_483 = v0_470
	v1_484 = v1_471
	zero_QMARK__485 = zero_QMARK__472
	one_QMARK__486 = one_QMARK__473
	goto b27
b31:
	;
	arg__17632_405, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("op"), vm.Keyword("const"), vm.Keyword("aux"), vm.Int(0), vm.Keyword("refs"), vm.NewArrayVector([]vm.Value{})})
	if callErr != nil {
		return nil, callErr
	}
	v406, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("replace-with"), arg__17632_405})
	if callErr != nil {
		return nil, callErr
	}
	v449 = v406
	nid_450 = nid_331
	f_451 = f_332
	case__17535_452 = case__17535_333
	op_453 = op_334
	refs_454 = refs_335
	r0_455 = r0_336
	r1_456 = r1_337
	v0_457 = v0_338
	v1_458 = v1_339
	zero_QMARK__459 = zero_QMARK__340
	one_QMARK__460 = one_QMARK__341
	goto b33
b32:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		nid_408 = nid_342
		f_409 = f_343
		case__17535_410 = case__17535_344
		op_411 = op_345
		refs_412 = refs_346
		r0_413 = r0_347
		r1_414 = r1_348
		v0_415 = v0_349
		v1_416 = v1_350
		zero_QMARK__417 = zero_QMARK__351
		one_QMARK__418 = one_QMARK__352
		goto b37
	} else {
		nid_419 = nid_342
		f_420 = f_343
		case__17535_421 = case__17535_344
		op_422 = op_345
		refs_423 = refs_346
		r0_424 = r0_347
		r1_425 = r1_348
		v0_426 = v0_349
		v1_427 = v1_350
		zero_QMARK__428 = zero_QMARK__351
		one_QMARK__429 = one_QMARK__352
		goto b38
	}
b33:
	;
	v462 = v449
	nid_463 = nid_450
	f_464 = f_451
	case__17535_465 = case__17535_452
	op_466 = op_453
	refs_467 = refs_454
	r0_468 = r0_455
	r1_469 = r1_456
	v0_470 = v0_457
	v1_471 = v1_458
	zero_QMARK__472 = zero_QMARK__459
	one_QMARK__473 = one_QMARK__460
	goto b30
b34:
	;
	v382 = or__x_365
	nid_383 = nid_354
	f_384 = f_355
	case__17535_385 = case__17535_356
	op_386 = op_357
	refs_387 = refs_358
	r0_388 = r0_359
	r1_389 = r1_360
	v0_390 = v0_361
	v1_391 = v1_362
	zero_QMARK__392 = zero_QMARK__363
	one_QMARK__393 = one_QMARK__364
	or__x_394 = or__x_365
	goto b36
b35:
	;
	v380, callErr = rt.InvokeValue(zero_QMARK__375, []vm.Value{v1_374})
	if callErr != nil {
		return nil, callErr
	}
	v382 = v380
	nid_383 = nid_366
	f_384 = f_367
	case__17535_385 = case__17535_368
	op_386 = op_369
	refs_387 = refs_370
	r0_388 = r0_371
	r1_389 = r1_372
	v0_390 = v0_373
	v1_391 = v1_374
	zero_QMARK__392 = zero_QMARK__375
	one_QMARK__393 = one_QMARK__376
	or__x_394 = or__x_377
	goto b36
b36:
	;
	if vm.IsTruthy(v382) {
		nid_331 = nid_383
		f_332 = f_384
		case__17535_333 = case__17535_385
		op_334 = op_386
		refs_335 = refs_387
		r0_336 = r0_388
		r1_337 = r1_389
		v0_338 = v0_390
		v1_339 = v1_391
		zero_QMARK__340 = zero_QMARK__392
		one_QMARK__341 = one_QMARK__393
		goto b31
	} else {
		nid_342 = nid_383
		f_343 = f_384
		case__17535_344 = case__17535_385
		op_345 = op_386
		refs_346 = refs_387
		r0_347 = r0_388
		r1_348 = r1_389
		v0_349 = v0_390
		v1_350 = v1_391
		zero_QMARK__351 = zero_QMARK__392
		one_QMARK__352 = one_QMARK__393
		goto b32
	}
b37:
	;
	v436 = vm.NIL
	nid_437 = nid_408
	f_438 = f_409
	case__17535_439 = case__17535_410
	op_440 = op_411
	refs_441 = refs_412
	r0_442 = r0_413
	r1_443 = r1_414
	v0_444 = v0_415
	v1_445 = v1_416
	zero_QMARK__446 = zero_QMARK__417
	one_QMARK__447 = one_QMARK__418
	goto b39
b38:
	;
	v436 = vm.NIL
	nid_437 = nid_419
	f_438 = f_420
	case__17535_439 = case__17535_421
	op_440 = op_422
	refs_441 = refs_423
	r0_442 = r0_424
	r1_443 = r1_425
	v0_444 = v0_426
	v1_445 = v1_427
	zero_QMARK__446 = zero_QMARK__428
	one_QMARK__447 = one_QMARK__429
	goto b39
b39:
	;
	v449 = v436
	nid_450 = nid_437
	f_451 = f_438
	case__17535_452 = case__17535_439
	op_453 = op_440
	refs_454 = refs_441
	r0_455 = r0_442
	r1_456 = r1_443
	v0_457 = v0_444
	v1_458 = v1_445
	zero_QMARK__459 = zero_QMARK__446
	one_QMARK__460 = one_QMARK__447
	goto b33
b40:
	;
	v516 = vm.NIL
	nid_517 = nid_488
	f_518 = f_489
	case__17535_519 = case__17535_490
	op_520 = op_491
	refs_521 = refs_492
	r0_522 = r0_493
	r1_523 = r1_494
	v0_524 = v0_495
	v1_525 = v1_496
	zero_QMARK__526 = zero_QMARK__497
	one_QMARK__527 = one_QMARK__498
	goto b42
b41:
	;
	v516 = vm.NIL
	nid_517 = nid_499
	f_518 = f_500
	case__17535_519 = case__17535_501
	op_520 = op_502
	refs_521 = refs_503
	r0_522 = r0_504
	r1_523 = r1_505
	v0_524 = v0_506
	v1_525 = v1_507
	zero_QMARK__526 = zero_QMARK__508
	one_QMARK__527 = one_QMARK__509
	goto b42
b42:
	;
	v529 = v516
	nid_530 = nid_517
	f_531 = f_518
	case__17535_532 = case__17535_519
	op_533 = op_520
	refs_534 = refs_521
	r0_535 = r0_522
	r1_536 = r1_523
	v0_537 = v0_524
	v1_538 = v1_525
	zero_QMARK__539 = zero_QMARK__526
	one_QMARK__540 = one_QMARK__527
	goto b24
}
func try_canonicalize(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var op_4 vm.Value
	var refs_6 vm.Value
	var and__x_16 vm.Value
	var nid_7 vm.Value
	var f_8 vm.Value
	var op_9 vm.Value
	var refs_10 vm.Value
	var arg__17707_117 vm.Value
	var arg__17713_121 vm.Value
	var arg__17714_122 vm.Value
	var v123 vm.Value
	var nid_11 vm.Value
	var f_12 vm.Value
	var op_13 vm.Value
	var refs_14 vm.Value
	var v127 vm.Value
	var nid_128 vm.Value
	var f_129 vm.Value
	var op_130 vm.Value
	var refs_131 vm.Value
	var nid_17 vm.Value
	var f_18 vm.Value
	var op_19 vm.Value
	var refs_20 vm.Value
	var and__x_21 vm.Value
	var arg__17650_30 vm.Value
	var and__x_31 bool
	var nid_22 vm.Value
	var f_23 vm.Value
	var op_24 vm.Value
	var refs_25 vm.Value
	var and__x_26 vm.Value
	var v104 vm.Value
	var nid_105 vm.Value
	var f_106 vm.Value
	var op_107 vm.Value
	var refs_108 vm.Value
	var and__x_109 vm.Value
	var nid_32 vm.Value
	var f_33 vm.Value
	var op_34 vm.Value
	var refs_35 vm.Value
	var and__x_36 bool
	var arg__17656_46 vm.Value
	var arg__17664_51 vm.Value
	var and__x_52 vm.Value
	var nid_37 vm.Value
	var f_38 vm.Value
	var op_39 vm.Value
	var refs_40 vm.Value
	var and__x_41 bool
	var v96 vm.Value
	var nid_97 vm.Value
	var f_98 vm.Value
	var op_99 vm.Value
	var refs_100 vm.Value
	var and__x_101 vm.Value
	var nid_53 vm.Value
	var f_54 vm.Value
	var op_55 vm.Value
	var refs_56 vm.Value
	var and__x_57 vm.Value
	var arg__17671_67 vm.Value
	var arg__17679_72 vm.Value
	var arg__17681_73 vm.Value
	var arg__17688_78 vm.Value
	var arg__17696_83 vm.Value
	var arg__17698_84 vm.Value
	var v85 vm.Value
	var nid_58 vm.Value
	var f_59 vm.Value
	var op_60 vm.Value
	var refs_61 vm.Value
	var and__x_62 vm.Value
	var v88 vm.Value
	var nid_89 vm.Value
	var f_90 vm.Value
	var op_91 vm.Value
	var refs_92 vm.Value
	var and__x_93 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = op_4, refs_6, and__x_16, nid_7, f_8, op_9, refs_10, arg__17707_117, arg__17713_121, arg__17714_122, v123, nid_11, f_12, op_13, refs_14, v127, nid_128, f_129, op_130, refs_131, nid_17, f_18, op_19, refs_20, and__x_21, arg__17650_30, and__x_31, nid_22, f_23, op_24, refs_25, and__x_26, v104, nid_105, f_106, op_107, refs_108, and__x_109, nid_32, f_33, op_34, refs_35, and__x_36, arg__17656_46, arg__17664_51, and__x_52, nid_37, f_38, op_39, refs_40, and__x_41, v96, nid_97, f_98, op_99, refs_100, and__x_101, nid_53, f_54, op_55, refs_56, and__x_57, arg__17671_67, arg__17679_72, arg__17681_73, arg__17688_78, arg__17696_83, arg__17698_84, v85, nid_58, f_59, op_60, refs_61, and__x_62, v88, nid_89, f_90, op_91, refs_92, and__x_93
	op_4, callErr = rt.InvokeValue(rt.LookupVar("ir", "op").Deref(), []vm.Value{arg0, arg1})
	if callErr != nil {
		return nil, callErr
	}
	refs_6, callErr = rt.InvokeValue(rt.LookupVar("ir", "refs").Deref(), []vm.Value{arg0, arg1})
	if callErr != nil {
		return nil, callErr
	}
	and__x_16, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.constfold", "commutative").Deref(), []vm.Value{op_4})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(and__x_16) {
		nid_17 = arg0
		f_18 = arg1
		op_19 = op_4
		refs_20 = refs_6
		and__x_21 = and__x_16
		goto b4
	} else {
		nid_22 = arg0
		f_23 = arg1
		op_24 = op_4
		refs_25 = refs_6
		and__x_26 = and__x_16
		goto b5
	}
b1:
	;
	arg__17707_117, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{refs_10, vm.Int(1)})
	if callErr != nil {
		return nil, callErr
	}
	arg__17713_121, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{refs_10, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	arg__17714_122, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__17707_117, arg__17713_121})
	if callErr != nil {
		return nil, callErr
	}
	v123, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("swap-refs"), arg__17714_122})
	if callErr != nil {
		return nil, callErr
	}
	v127 = v123
	nid_128 = nid_7
	f_129 = f_8
	op_130 = op_9
	refs_131 = refs_10
	goto b3
b2:
	;
	v127 = vm.NIL
	nid_128 = nid_11
	f_129 = f_12
	op_130 = op_13
	refs_131 = refs_14
	goto b3
b3:
	;
	return v127, nil
b4:
	;
	arg__17650_30, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{refs_20})
	if callErr != nil {
		return nil, callErr
	}
	and__x_31 = arg__17650_30 == vm.Int(2)
	if and__x_31 {
		nid_32 = nid_17
		f_33 = f_18
		op_34 = op_19
		refs_35 = refs_20
		and__x_36 = and__x_31
		goto b7
	} else {
		nid_37 = nid_17
		f_38 = f_18
		op_39 = op_19
		refs_40 = refs_20
		and__x_41 = and__x_31
		goto b8
	}
b5:
	;
	v104 = and__x_26
	nid_105 = nid_22
	f_106 = f_23
	op_107 = op_24
	refs_108 = refs_25
	and__x_109 = and__x_26
	goto b6
b6:
	;
	if vm.IsTruthy(v104) {
		nid_7 = nid_105
		f_8 = f_106
		op_9 = op_107
		refs_10 = refs_108
		goto b1
	} else {
		nid_11 = nid_105
		f_12 = f_106
		op_13 = op_107
		refs_14 = refs_108
		goto b2
	}
b7:
	;
	arg__17656_46, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{refs_35, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	arg__17664_51, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{refs_35, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	and__x_52, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.constfold", "const?").Deref(), []vm.Value{arg__17664_51, f_33})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(and__x_52) {
		nid_53 = nid_32
		f_54 = f_33
		op_55 = op_34
		refs_56 = refs_35
		and__x_57 = and__x_52
		goto b10
	} else {
		nid_58 = nid_32
		f_59 = f_33
		op_60 = op_34
		refs_61 = refs_35
		and__x_62 = and__x_52
		goto b11
	}
b8:
	;
	v96 = vm.Boolean(and__x_41)
	nid_97 = nid_37
	f_98 = f_38
	op_99 = op_39
	refs_100 = refs_40
	and__x_101 = vm.Boolean(and__x_41)
	goto b9
b9:
	;
	v104 = v96
	nid_105 = nid_97
	f_106 = f_98
	op_107 = op_99
	refs_108 = refs_100
	and__x_109 = and__x_21
	goto b6
b10:
	;
	arg__17671_67, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{refs_56, vm.Int(1)})
	if callErr != nil {
		return nil, callErr
	}
	arg__17679_72, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{refs_56, vm.Int(1)})
	if callErr != nil {
		return nil, callErr
	}
	arg__17681_73, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.constfold", "const?").Deref(), []vm.Value{arg__17679_72, f_54})
	if callErr != nil {
		return nil, callErr
	}
	arg__17688_78, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{refs_56, vm.Int(1)})
	if callErr != nil {
		return nil, callErr
	}
	arg__17696_83, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{refs_56, vm.Int(1)})
	if callErr != nil {
		return nil, callErr
	}
	arg__17698_84, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.constfold", "const?").Deref(), []vm.Value{arg__17696_83, f_54})
	if callErr != nil {
		return nil, callErr
	}
	v85, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not").Deref(), []vm.Value{arg__17698_84})
	if callErr != nil {
		return nil, callErr
	}
	v88 = v85
	nid_89 = nid_53
	f_90 = f_54
	op_91 = op_55
	refs_92 = refs_56
	and__x_93 = and__x_57
	goto b12
b11:
	;
	v88 = and__x_62
	nid_89 = nid_58
	f_90 = f_59
	op_91 = op_60
	refs_92 = refs_61
	and__x_93 = and__x_62
	goto b12
b12:
	;
	v96 = v88
	nid_97 = nid_89
	f_98 = f_90
	op_99 = op_91
	refs_100 = refs_92
	and__x_101 = vm.Boolean(and__x_36)
	goto b9
}
