package walk

import (
	rt "github.com/nooga/let-go/pkg/rt"
	vm "github.com/nooga/let-go/pkg/vm"
)

func walk(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var v10 vm.Value
	var inner_3 vm.Value
	var outer_4 vm.Value
	var form_5 vm.Value
	var arg__31072_14 vm.Value
	var arg__31080_18 vm.Value
	var arg__31081_19 vm.Value
	var arg__31089_22 vm.Value
	var arg__31097_26 vm.Value
	var arg__31098_27 vm.Value
	var v28 vm.Value
	var inner_6 vm.Value
	var outer_7 vm.Value
	var form_8 vm.Value
	var v37 vm.Value
	var v163 vm.Value
	var inner_164 vm.Value
	var outer_165 vm.Value
	var form_166 vm.Value
	var inner_30 vm.Value
	var outer_31 vm.Value
	var form_32 vm.Value
	var arg__31105_40 vm.Value
	var arg__31110_42 vm.Value
	var arg__31111_43 vm.Value
	var arg__31115_45 vm.Value
	var arg__31120_47 vm.Value
	var arg__31121_48 vm.Value
	var arg__31126_51 vm.Value
	var arg__31131_53 vm.Value
	var arg__31132_54 vm.Value
	var arg__31136_56 vm.Value
	var arg__31141_58 vm.Value
	var arg__31142_59 vm.Value
	var arg__31143_60 vm.Value
	var arg__31148_62 vm.Value
	var arg__31153_64 vm.Value
	var arg__31154_65 vm.Value
	var arg__31158_67 vm.Value
	var arg__31163_69 vm.Value
	var arg__31164_70 vm.Value
	var arg__31169_73 vm.Value
	var arg__31174_75 vm.Value
	var arg__31175_76 vm.Value
	var arg__31179_78 vm.Value
	var arg__31184_80 vm.Value
	var arg__31185_81 vm.Value
	var arg__31186_82 vm.Value
	var v83 vm.Value
	var inner_33 vm.Value
	var outer_34 vm.Value
	var form_35 vm.Value
	var v92 vm.Value
	var v158 vm.Value
	var inner_159 vm.Value
	var outer_160 vm.Value
	var form_161 vm.Value
	var inner_85 vm.Value
	var outer_86 vm.Value
	var form_87 vm.Value
	var arg__31195_95 vm.Value
	var arg__31202_97 vm.Value
	var v98 vm.Value
	var inner_88 vm.Value
	var outer_89 vm.Value
	var form_90 vm.Value
	var v107 vm.Value
	var v153 vm.Value
	var inner_154 vm.Value
	var outer_155 vm.Value
	var form_156 vm.Value
	var inner_100 vm.Value
	var outer_101 vm.Value
	var form_102 vm.Value
	var arg__31209_110 vm.Value
	var arg__31215_112 vm.Value
	var arg__31220_115 vm.Value
	var arg__31226_117 vm.Value
	var arg__31227_118 vm.Value
	var arg__31232_120 vm.Value
	var arg__31238_122 vm.Value
	var arg__31243_125 vm.Value
	var arg__31249_127 vm.Value
	var arg__31250_128 vm.Value
	var v129 vm.Value
	var inner_103 vm.Value
	var outer_104 vm.Value
	var form_105 vm.Value
	var v148 vm.Value
	var inner_149 vm.Value
	var outer_150 vm.Value
	var form_151 vm.Value
	var inner_131 vm.Value
	var outer_132 vm.Value
	var form_133 vm.Value
	var v139 vm.Value
	var inner_134 vm.Value
	var outer_135 vm.Value
	var form_136 vm.Value
	var v143 vm.Value
	var inner_144 vm.Value
	var outer_145 vm.Value
	var form_146 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = v10, inner_3, outer_4, form_5, arg__31072_14, arg__31080_18, arg__31081_19, arg__31089_22, arg__31097_26, arg__31098_27, v28, inner_6, outer_7, form_8, v37, v163, inner_164, outer_165, form_166, inner_30, outer_31, form_32, arg__31105_40, arg__31110_42, arg__31111_43, arg__31115_45, arg__31120_47, arg__31121_48, arg__31126_51, arg__31131_53, arg__31132_54, arg__31136_56, arg__31141_58, arg__31142_59, arg__31143_60, arg__31148_62, arg__31153_64, arg__31154_65, arg__31158_67, arg__31163_69, arg__31164_70, arg__31169_73, arg__31174_75, arg__31175_76, arg__31179_78, arg__31184_80, arg__31185_81, arg__31186_82, v83, inner_33, outer_34, form_35, v92, v158, inner_159, outer_160, form_161, inner_85, outer_86, form_87, arg__31195_95, arg__31202_97, v98, inner_88, outer_89, form_90, v107, v153, inner_154, outer_155, form_156, inner_100, outer_101, form_102, arg__31209_110, arg__31215_112, arg__31220_115, arg__31226_117, arg__31227_118, arg__31232_120, arg__31238_122, arg__31243_125, arg__31249_127, arg__31250_128, v129, inner_103, outer_104, form_105, v148, inner_149, outer_150, form_151, inner_131, outer_132, form_133, v139, inner_134, outer_135, form_136, v143, inner_144, outer_145, form_146
	v10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "list?").Deref(), []vm.Value{arg2})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v10) {
		inner_3 = arg0
		outer_4 = arg1
		form_5 = arg2
		goto b1
	} else {
		inner_6 = arg0
		outer_7 = arg1
		form_8 = arg2
		goto b2
	}
b1:
	;
	arg__31072_14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{inner_3, form_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__31080_18, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{inner_3, form_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__31081_19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "apply").Deref(), []vm.Value{rt.LookupVar("clojure.core", "list").Deref(), arg__31080_18})
	if callErr != nil {
		return nil, callErr
	}
	arg__31089_22, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{inner_3, form_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__31097_26, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{inner_3, form_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__31098_27, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "apply").Deref(), []vm.Value{rt.LookupVar("clojure.core", "list").Deref(), arg__31097_26})
	if callErr != nil {
		return nil, callErr
	}
	v28, callErr = rt.InvokeValue(outer_4, []vm.Value{arg__31098_27})
	if callErr != nil {
		return nil, callErr
	}
	v163 = v28
	inner_164 = inner_3
	outer_165 = outer_4
	form_166 = form_5
	goto b3
b2:
	;
	v37, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map-entry?").Deref(), []vm.Value{form_8})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v37) {
		inner_30 = inner_6
		outer_31 = outer_7
		form_32 = form_8
		goto b4
	} else {
		inner_33 = inner_6
		outer_34 = outer_7
		form_35 = form_8
		goto b5
	}
b3:
	;
	return v163, nil
b4:
	;
	arg__31105_40, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "key").Deref(), []vm.Value{form_32})
	if callErr != nil {
		return nil, callErr
	}
	arg__31110_42, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "key").Deref(), []vm.Value{form_32})
	if callErr != nil {
		return nil, callErr
	}
	arg__31111_43, callErr = rt.InvokeValue(inner_30, []vm.Value{arg__31110_42})
	if callErr != nil {
		return nil, callErr
	}
	arg__31115_45, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "val").Deref(), []vm.Value{form_32})
	if callErr != nil {
		return nil, callErr
	}
	arg__31120_47, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "val").Deref(), []vm.Value{form_32})
	if callErr != nil {
		return nil, callErr
	}
	arg__31121_48, callErr = rt.InvokeValue(inner_30, []vm.Value{arg__31120_47})
	if callErr != nil {
		return nil, callErr
	}
	arg__31126_51, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "key").Deref(), []vm.Value{form_32})
	if callErr != nil {
		return nil, callErr
	}
	arg__31131_53, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "key").Deref(), []vm.Value{form_32})
	if callErr != nil {
		return nil, callErr
	}
	arg__31132_54, callErr = rt.InvokeValue(inner_30, []vm.Value{arg__31131_53})
	if callErr != nil {
		return nil, callErr
	}
	arg__31136_56, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "val").Deref(), []vm.Value{form_32})
	if callErr != nil {
		return nil, callErr
	}
	arg__31141_58, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "val").Deref(), []vm.Value{form_32})
	if callErr != nil {
		return nil, callErr
	}
	arg__31142_59, callErr = rt.InvokeValue(inner_30, []vm.Value{arg__31141_58})
	if callErr != nil {
		return nil, callErr
	}
	arg__31143_60, callErr = rt.InvokeValue(rt.LookupVar("clojure.lang.MapEntry", "create").Deref(), []vm.Value{arg__31132_54, arg__31142_59})
	if callErr != nil {
		return nil, callErr
	}
	arg__31148_62, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "key").Deref(), []vm.Value{form_32})
	if callErr != nil {
		return nil, callErr
	}
	arg__31153_64, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "key").Deref(), []vm.Value{form_32})
	if callErr != nil {
		return nil, callErr
	}
	arg__31154_65, callErr = rt.InvokeValue(inner_30, []vm.Value{arg__31153_64})
	if callErr != nil {
		return nil, callErr
	}
	arg__31158_67, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "val").Deref(), []vm.Value{form_32})
	if callErr != nil {
		return nil, callErr
	}
	arg__31163_69, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "val").Deref(), []vm.Value{form_32})
	if callErr != nil {
		return nil, callErr
	}
	arg__31164_70, callErr = rt.InvokeValue(inner_30, []vm.Value{arg__31163_69})
	if callErr != nil {
		return nil, callErr
	}
	arg__31169_73, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "key").Deref(), []vm.Value{form_32})
	if callErr != nil {
		return nil, callErr
	}
	arg__31174_75, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "key").Deref(), []vm.Value{form_32})
	if callErr != nil {
		return nil, callErr
	}
	arg__31175_76, callErr = rt.InvokeValue(inner_30, []vm.Value{arg__31174_75})
	if callErr != nil {
		return nil, callErr
	}
	arg__31179_78, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "val").Deref(), []vm.Value{form_32})
	if callErr != nil {
		return nil, callErr
	}
	arg__31184_80, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "val").Deref(), []vm.Value{form_32})
	if callErr != nil {
		return nil, callErr
	}
	arg__31185_81, callErr = rt.InvokeValue(inner_30, []vm.Value{arg__31184_80})
	if callErr != nil {
		return nil, callErr
	}
	arg__31186_82, callErr = rt.InvokeValue(rt.LookupVar("clojure.lang.MapEntry", "create").Deref(), []vm.Value{arg__31175_76, arg__31185_81})
	if callErr != nil {
		return nil, callErr
	}
	v83, callErr = rt.InvokeValue(outer_31, []vm.Value{arg__31186_82})
	if callErr != nil {
		return nil, callErr
	}
	v158 = v83
	inner_159 = inner_30
	outer_160 = outer_31
	form_161 = form_32
	goto b6
b5:
	;
	v92, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq?").Deref(), []vm.Value{form_35})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v92) {
		inner_85 = inner_33
		outer_86 = outer_34
		form_87 = form_35
		goto b7
	} else {
		inner_88 = inner_33
		outer_89 = outer_34
		form_90 = form_35
		goto b8
	}
b6:
	;
	v163 = v158
	inner_164 = inner_159
	outer_165 = outer_160
	form_166 = form_161
	goto b3
b7:
	;
	arg__31195_95, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{inner_85, form_87})
	if callErr != nil {
		return nil, callErr
	}
	arg__31202_97, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{inner_85, form_87})
	if callErr != nil {
		return nil, callErr
	}
	v98, callErr = rt.InvokeValue(outer_86, []vm.Value{arg__31202_97})
	if callErr != nil {
		return nil, callErr
	}
	v153 = v98
	inner_154 = inner_85
	outer_155 = outer_86
	form_156 = form_87
	goto b9
b8:
	;
	v107, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "coll?").Deref(), []vm.Value{form_90})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v107) {
		inner_100 = inner_88
		outer_101 = outer_89
		form_102 = form_90
		goto b10
	} else {
		inner_103 = inner_88
		outer_104 = outer_89
		form_105 = form_90
		goto b11
	}
b9:
	;
	v158 = v153
	inner_159 = inner_154
	outer_160 = outer_155
	form_161 = form_156
	goto b6
b10:
	;
	arg__31209_110, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "empty").Deref(), []vm.Value{form_102})
	if callErr != nil {
		return nil, callErr
	}
	arg__31215_112, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{inner_100, form_102})
	if callErr != nil {
		return nil, callErr
	}
	arg__31220_115, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "empty").Deref(), []vm.Value{form_102})
	if callErr != nil {
		return nil, callErr
	}
	arg__31226_117, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{inner_100, form_102})
	if callErr != nil {
		return nil, callErr
	}
	arg__31227_118, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__31220_115, arg__31226_117})
	if callErr != nil {
		return nil, callErr
	}
	arg__31232_120, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "empty").Deref(), []vm.Value{form_102})
	if callErr != nil {
		return nil, callErr
	}
	arg__31238_122, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{inner_100, form_102})
	if callErr != nil {
		return nil, callErr
	}
	arg__31243_125, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "empty").Deref(), []vm.Value{form_102})
	if callErr != nil {
		return nil, callErr
	}
	arg__31249_127, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{inner_100, form_102})
	if callErr != nil {
		return nil, callErr
	}
	arg__31250_128, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__31243_125, arg__31249_127})
	if callErr != nil {
		return nil, callErr
	}
	v129, callErr = rt.InvokeValue(outer_101, []vm.Value{arg__31250_128})
	if callErr != nil {
		return nil, callErr
	}
	v148 = v129
	inner_149 = inner_100
	outer_150 = outer_101
	form_151 = form_102
	goto b12
b11:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		inner_131 = inner_103
		outer_132 = outer_104
		form_133 = form_105
		goto b13
	} else {
		inner_134 = inner_103
		outer_135 = outer_104
		form_136 = form_105
		goto b14
	}
b12:
	;
	v153 = v148
	inner_154 = inner_149
	outer_155 = outer_150
	form_156 = form_151
	goto b9
b13:
	;
	v139, callErr = rt.InvokeValue(outer_132, []vm.Value{form_133})
	if callErr != nil {
		return nil, callErr
	}
	v143 = v139
	inner_144 = inner_131
	outer_145 = outer_132
	form_146 = form_133
	goto b15
b14:
	;
	v143 = vm.NIL
	inner_144 = inner_134
	outer_145 = outer_135
	form_146 = form_136
	goto b15
b15:
	;
	v148 = v143
	inner_149 = inner_144
	outer_150 = outer_145
	form_151 = form_146
	goto b12
}
func postwalk(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__31259_5 vm.Value
	var arg__31268_10 vm.Value
	var v11 vm.Value
	var callErr error
	_, _, _ = arg__31259_5, arg__31268_10, v11
	arg__31259_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "partial").Deref(), []vm.Value{rt.LookupVar("clojure.walk", "postwalk").Deref(), arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__31268_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "partial").Deref(), []vm.Value{rt.LookupVar("clojure.walk", "postwalk").Deref(), arg0})
	if callErr != nil {
		return nil, callErr
	}
	v11, callErr = rt.InvokeValue(rt.LookupVar("clojure.walk", "walk").Deref(), []vm.Value{arg__31268_10, arg0, arg1})
	if callErr != nil {
		return nil, callErr
	}
	return v11, nil
}
func keywordize_keys(arg0 vm.Value) (vm.Value, error) {
	var v9 vm.Value
	var callErr error
	_ = v9
	v9, callErr = rt.InvokeValue(rt.LookupVar("clojure.walk", "postwalk").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v7 vm.Value
		var x_2 vm.Value
		var f_3 vm.Value
		var arg__31328_11 vm.Value
		var arg__31336_15 vm.Value
		var v16 vm.Value
		var x_4 vm.Value
		var f_5 vm.Value
		var v19 vm.Value
		var x_20 vm.Value
		var f_21 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _, _, _, _ = v7, x_2, f_3, arg__31328_11, arg__31336_15, v16, x_4, f_5, v19, x_20, f_21
		v7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map?").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		if vm.IsTruthy(v7) {
			x_2 = arg0
			f_3 = rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
				var k_6 vm.Value
				var v_12 vm.Value
				var v20 vm.Value
				var vec__31271_13 vm.Value
				var k_14 vm.Value
				var v_15 vm.Value
				var arg__31293_24 vm.Value
				var v25 vm.Value
				var vec__31271_16 vm.Value
				var k_17 vm.Value
				var v_18 vm.Value
				var v28 vm.Value
				var v30 vm.Value
				var vec__31271_31 vm.Value
				var k_32 vm.Value
				var v_33 vm.Value
				var callErr error
				_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = k_6, v_12, v20, vec__31271_13, k_14, v_15, arg__31293_24, v25, vec__31271_16, k_17, v_18, v28, v30, vec__31271_31, k_32, v_33
				k_6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(0), vm.NIL})
				if callErr != nil {
					return nil, callErr
				}
				v_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(1), vm.NIL})
				if callErr != nil {
					return nil, callErr
				}
				v20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "string?").Deref(), []vm.Value{k_6})
				if callErr != nil {
					return nil, callErr
				}
				if vm.IsTruthy(v20) {
					vec__31271_13 = arg0
					k_14 = k_6
					v_15 = v_12
					goto b1
				} else {
					vec__31271_16 = arg0
					k_17 = k_6
					v_18 = v_12
					goto b2
				}
			b1:
				;
				arg__31293_24, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "keyword").Deref(), []vm.Value{k_14})
				if callErr != nil {
					return nil, callErr
				}
				v25, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__31293_24, v_15})
				if callErr != nil {
					return nil, callErr
				}
				v30 = v25
				vec__31271_31 = vec__31271_13
				k_32 = k_14
				v_33 = v_15
				goto b3
			b2:
				;
				v28, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{k_17, v_18})
				if callErr != nil {
					return nil, callErr
				}
				v30 = v28
				vec__31271_31 = vec__31271_16
				k_32 = k_17
				v_33 = v_18
				goto b3
			b3:
				;
				return v30, nil
			})
			goto b1
		} else {
			x_4 = arg0
			f_5 = rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
				var k_6 vm.Value
				var v_12 vm.Value
				var v20 vm.Value
				var vec__31271_13 vm.Value
				var k_14 vm.Value
				var v_15 vm.Value
				var arg__31293_24 vm.Value
				var v25 vm.Value
				var vec__31271_16 vm.Value
				var k_17 vm.Value
				var v_18 vm.Value
				var v28 vm.Value
				var v30 vm.Value
				var vec__31271_31 vm.Value
				var k_32 vm.Value
				var v_33 vm.Value
				var callErr error
				_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = k_6, v_12, v20, vec__31271_13, k_14, v_15, arg__31293_24, v25, vec__31271_16, k_17, v_18, v28, v30, vec__31271_31, k_32, v_33
				k_6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(0), vm.NIL})
				if callErr != nil {
					return nil, callErr
				}
				v_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(1), vm.NIL})
				if callErr != nil {
					return nil, callErr
				}
				v20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "string?").Deref(), []vm.Value{k_6})
				if callErr != nil {
					return nil, callErr
				}
				if vm.IsTruthy(v20) {
					vec__31271_13 = arg0
					k_14 = k_6
					v_15 = v_12
					goto b1
				} else {
					vec__31271_16 = arg0
					k_17 = k_6
					v_18 = v_12
					goto b2
				}
			b1:
				;
				arg__31293_24, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "keyword").Deref(), []vm.Value{k_14})
				if callErr != nil {
					return nil, callErr
				}
				v25, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__31293_24, v_15})
				if callErr != nil {
					return nil, callErr
				}
				v30 = v25
				vec__31271_31 = vec__31271_13
				k_32 = k_14
				v_33 = v_15
				goto b3
			b2:
				;
				v28, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{k_17, v_18})
				if callErr != nil {
					return nil, callErr
				}
				v30 = v28
				vec__31271_31 = vec__31271_16
				k_32 = k_17
				v_33 = v_18
				goto b3
			b3:
				;
				return v30, nil
			})
			goto b2
		}
	b1:
		;
		arg__31328_11, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{f_3, x_2})
		if callErr != nil {
			return nil, callErr
		}
		arg__31336_15, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{f_3, x_2})
		if callErr != nil {
			return nil, callErr
		}
		v16, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{vm.EmptyPersistentMap, arg__31336_15})
		if callErr != nil {
			return nil, callErr
		}
		v19 = v16
		x_20 = x_2
		f_21 = f_3
		goto b3
	b2:
		;
		v19 = x_4
		x_20 = x_4
		f_21 = f_5
		goto b3
	b3:
		;
		return v19, nil
	}), arg0})
	if callErr != nil {
		return nil, callErr
	}
	return v9, nil
}
func prewalk(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__31344_5 vm.Value
	var arg__31349_7 vm.Value
	var arg__31356_12 vm.Value
	var arg__31361_14 vm.Value
	var v15 vm.Value
	var callErr error
	_, _, _, _, _ = arg__31344_5, arg__31349_7, arg__31356_12, arg__31361_14, v15
	arg__31344_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "partial").Deref(), []vm.Value{rt.LookupVar("clojure.walk", "prewalk").Deref(), arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__31349_7, callErr = rt.InvokeValue(arg0, []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__31356_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "partial").Deref(), []vm.Value{rt.LookupVar("clojure.walk", "prewalk").Deref(), arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__31361_14, callErr = rt.InvokeValue(arg0, []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	v15, callErr = rt.InvokeValue(rt.LookupVar("clojure.walk", "walk").Deref(), []vm.Value{arg__31356_12, rt.LookupVar("clojure.core", "identity").Deref(), arg__31361_14})
	if callErr != nil {
		return nil, callErr
	}
	return v15, nil
}
func prewalk_replace(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var v9 vm.Value
	var callErr error
	_ = v9
	v9, callErr = rt.InvokeValue(rt.LookupVar("clojure.walk", "prewalk").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v7 vm.Value
		var x_2 vm.Value
		var smap_3 vm.Value
		var v9 vm.Value
		var x_4 vm.Value
		var smap_5 vm.Value
		var v12 vm.Value
		var x_13 vm.Value
		var smap_14 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _, _ = v7, x_2, smap_3, v9, x_4, smap_5, v12, x_13, smap_14
		v7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "contains?").Deref(), []vm.Value{arg0, arg0})
		if callErr != nil {
			return nil, callErr
		}
		if vm.IsTruthy(v7) {
			x_2 = arg0
			smap_3 = arg0
			goto b1
		} else {
			x_4 = arg0
			smap_5 = arg0
			goto b2
		}
	b1:
		;
		v9, callErr = rt.InvokeValue(smap_3, []vm.Value{x_2})
		if callErr != nil {
			return nil, callErr
		}
		v12 = v9
		x_13 = x_2
		smap_14 = smap_3
		goto b3
	b2:
		;
		v12 = x_4
		x_13 = x_4
		smap_14 = smap_5
		goto b3
	b3:
		;
		return v12, nil
	}), arg1})
	if callErr != nil {
		return nil, callErr
	}
	return v9, nil
}
func stringify_keys(arg0 vm.Value) (vm.Value, error) {
	var v9 vm.Value
	var callErr error
	_ = v9
	v9, callErr = rt.InvokeValue(rt.LookupVar("clojure.walk", "postwalk").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v7 vm.Value
		var x_2 vm.Value
		var f_3 vm.Value
		var arg__31440_11 vm.Value
		var arg__31448_15 vm.Value
		var v16 vm.Value
		var x_4 vm.Value
		var f_5 vm.Value
		var v19 vm.Value
		var x_20 vm.Value
		var f_21 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _, _, _, _ = v7, x_2, f_3, arg__31440_11, arg__31448_15, v16, x_4, f_5, v19, x_20, f_21
		v7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map?").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		if vm.IsTruthy(v7) {
			x_2 = arg0
			f_3 = rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
				var k_6 vm.Value
				var v_12 vm.Value
				var v20 vm.Value
				var vec__31383_13 vm.Value
				var k_14 vm.Value
				var v_15 vm.Value
				var arg__31405_24 vm.Value
				var v25 vm.Value
				var vec__31383_16 vm.Value
				var k_17 vm.Value
				var v_18 vm.Value
				var v28 vm.Value
				var v30 vm.Value
				var vec__31383_31 vm.Value
				var k_32 vm.Value
				var v_33 vm.Value
				var callErr error
				_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = k_6, v_12, v20, vec__31383_13, k_14, v_15, arg__31405_24, v25, vec__31383_16, k_17, v_18, v28, v30, vec__31383_31, k_32, v_33
				k_6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(0), vm.NIL})
				if callErr != nil {
					return nil, callErr
				}
				v_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(1), vm.NIL})
				if callErr != nil {
					return nil, callErr
				}
				v20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "keyword?").Deref(), []vm.Value{k_6})
				if callErr != nil {
					return nil, callErr
				}
				if vm.IsTruthy(v20) {
					vec__31383_13 = arg0
					k_14 = k_6
					v_15 = v_12
					goto b1
				} else {
					vec__31383_16 = arg0
					k_17 = k_6
					v_18 = v_12
					goto b2
				}
			b1:
				;
				arg__31405_24, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "name").Deref(), []vm.Value{k_14})
				if callErr != nil {
					return nil, callErr
				}
				v25, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__31405_24, v_15})
				if callErr != nil {
					return nil, callErr
				}
				v30 = v25
				vec__31383_31 = vec__31383_13
				k_32 = k_14
				v_33 = v_15
				goto b3
			b2:
				;
				v28, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{k_17, v_18})
				if callErr != nil {
					return nil, callErr
				}
				v30 = v28
				vec__31383_31 = vec__31383_16
				k_32 = k_17
				v_33 = v_18
				goto b3
			b3:
				;
				return v30, nil
			})
			goto b1
		} else {
			x_4 = arg0
			f_5 = rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
				var k_6 vm.Value
				var v_12 vm.Value
				var v20 vm.Value
				var vec__31383_13 vm.Value
				var k_14 vm.Value
				var v_15 vm.Value
				var arg__31405_24 vm.Value
				var v25 vm.Value
				var vec__31383_16 vm.Value
				var k_17 vm.Value
				var v_18 vm.Value
				var v28 vm.Value
				var v30 vm.Value
				var vec__31383_31 vm.Value
				var k_32 vm.Value
				var v_33 vm.Value
				var callErr error
				_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = k_6, v_12, v20, vec__31383_13, k_14, v_15, arg__31405_24, v25, vec__31383_16, k_17, v_18, v28, v30, vec__31383_31, k_32, v_33
				k_6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(0), vm.NIL})
				if callErr != nil {
					return nil, callErr
				}
				v_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(1), vm.NIL})
				if callErr != nil {
					return nil, callErr
				}
				v20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "keyword?").Deref(), []vm.Value{k_6})
				if callErr != nil {
					return nil, callErr
				}
				if vm.IsTruthy(v20) {
					vec__31383_13 = arg0
					k_14 = k_6
					v_15 = v_12
					goto b1
				} else {
					vec__31383_16 = arg0
					k_17 = k_6
					v_18 = v_12
					goto b2
				}
			b1:
				;
				arg__31405_24, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "name").Deref(), []vm.Value{k_14})
				if callErr != nil {
					return nil, callErr
				}
				v25, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__31405_24, v_15})
				if callErr != nil {
					return nil, callErr
				}
				v30 = v25
				vec__31383_31 = vec__31383_13
				k_32 = k_14
				v_33 = v_15
				goto b3
			b2:
				;
				v28, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{k_17, v_18})
				if callErr != nil {
					return nil, callErr
				}
				v30 = v28
				vec__31383_31 = vec__31383_16
				k_32 = k_17
				v_33 = v_18
				goto b3
			b3:
				;
				return v30, nil
			})
			goto b2
		}
	b1:
		;
		arg__31440_11, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{f_3, x_2})
		if callErr != nil {
			return nil, callErr
		}
		arg__31448_15, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{f_3, x_2})
		if callErr != nil {
			return nil, callErr
		}
		v16, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{vm.EmptyPersistentMap, arg__31448_15})
		if callErr != nil {
			return nil, callErr
		}
		v19 = v16
		x_20 = x_2
		f_21 = f_3
		goto b3
	b2:
		;
		v19 = x_4
		x_20 = x_4
		f_21 = f_5
		goto b3
	b3:
		;
		return v19, nil
	}), arg0})
	if callErr != nil {
		return nil, callErr
	}
	return v9, nil
}
func postwalk_replace(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var v9 vm.Value
	var callErr error
	_ = v9
	v9, callErr = rt.InvokeValue(rt.LookupVar("clojure.walk", "postwalk").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v7 vm.Value
		var x_2 vm.Value
		var smap_3 vm.Value
		var v9 vm.Value
		var x_4 vm.Value
		var smap_5 vm.Value
		var v12 vm.Value
		var x_13 vm.Value
		var smap_14 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _, _ = v7, x_2, smap_3, v9, x_4, smap_5, v12, x_13, smap_14
		v7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "contains?").Deref(), []vm.Value{arg0, arg0})
		if callErr != nil {
			return nil, callErr
		}
		if vm.IsTruthy(v7) {
			x_2 = arg0
			smap_3 = arg0
			goto b1
		} else {
			x_4 = arg0
			smap_5 = arg0
			goto b2
		}
	b1:
		;
		v9, callErr = rt.InvokeValue(smap_3, []vm.Value{x_2})
		if callErr != nil {
			return nil, callErr
		}
		v12 = v9
		x_13 = x_2
		smap_14 = smap_3
		goto b3
	b2:
		;
		v12 = x_4
		x_13 = x_4
		smap_14 = smap_5
		goto b3
	b3:
		;
		return v12, nil
	}), arg1})
	if callErr != nil {
		return nil, callErr
	}
	return v9, nil
}
