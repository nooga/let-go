package ir_dump

import (
	rt "github.com/nooga/let-go/pkg/rt"
	vm "github.com/nooga/let-go/pkg/vm"
)

func format_args(arg0 vm.Value) (vm.Value, error) {
	var arg__8907_5 vm.Value
	var arg__8925_11 vm.Value
	var v12 vm.Value
	var callErr error
	_, _, _ = arg__8907_5, arg__8925_11, v12
	arg__8907_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v4 vm.Value
		var callErr error
		_ = v4
		v4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("v"), arg0})
		if callErr != nil {
			return nil, callErr
		}
		return v4, nil
	}), arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__8925_11, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v4 vm.Value
		var callErr error
		_ = v4
		v4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("v"), arg0})
		if callErr != nil {
			return nil, callErr
		}
		return v4, nil
	}), arg0})
	if callErr != nil {
		return nil, callErr
	}
	v12, callErr = rt.InvokeValue(rt.LookupVar("clojure.string", "join").Deref(), []vm.Value{vm.String(", "), arg__8925_11})
	if callErr != nil {
		return nil, callErr
	}
	return v12, nil
}
func scalar_type_display(arg0 vm.Value) vm.Value {
	var v6 bool
	var case__8926_1 vm.Value
	var t_2 vm.Value
	var case__8926_3 vm.Value
	var t_4 vm.Value
	var v15 bool
	var v154 vm.Value
	var case__8926_155 vm.Value
	var t_156 vm.Value
	var case__8926_10 vm.Value
	var t_11 vm.Value
	var case__8926_12 vm.Value
	var t_13 vm.Value
	var v24 bool
	var v150 vm.Value
	var case__8926_151 vm.Value
	var t_152 vm.Value
	var case__8926_19 vm.Value
	var t_20 vm.Value
	var case__8926_21 vm.Value
	var t_22 vm.Value
	var v33 bool
	var v146 vm.Value
	var case__8926_147 vm.Value
	var t_148 vm.Value
	var case__8926_28 vm.Value
	var t_29 vm.Value
	var case__8926_30 vm.Value
	var t_31 vm.Value
	var v42 bool
	var v142 vm.Value
	var case__8926_143 vm.Value
	var t_144 vm.Value
	var case__8926_37 vm.Value
	var t_38 vm.Value
	var case__8926_39 vm.Value
	var t_40 vm.Value
	var v51 bool
	var v138 vm.Value
	var case__8926_139 vm.Value
	var t_140 vm.Value
	var case__8926_46 vm.Value
	var t_47 vm.Value
	var case__8926_48 vm.Value
	var t_49 vm.Value
	var v60 bool
	var v134 vm.Value
	var case__8926_135 vm.Value
	var t_136 vm.Value
	var case__8926_55 vm.Value
	var t_56 vm.Value
	var case__8926_57 vm.Value
	var t_58 vm.Value
	var v69 bool
	var v130 vm.Value
	var case__8926_131 vm.Value
	var t_132 vm.Value
	var case__8926_64 vm.Value
	var t_65 vm.Value
	var case__8926_66 vm.Value
	var t_67 vm.Value
	var v78 bool
	var v126 vm.Value
	var case__8926_127 vm.Value
	var t_128 vm.Value
	var case__8926_73 vm.Value
	var t_74 vm.Value
	var case__8926_75 vm.Value
	var t_76 vm.Value
	var v87 bool
	var v122 vm.Value
	var case__8926_123 vm.Value
	var t_124 vm.Value
	var case__8926_82 vm.Value
	var t_83 vm.Value
	var case__8926_84 vm.Value
	var t_85 vm.Value
	var v96 bool
	var v118 vm.Value
	var case__8926_119 vm.Value
	var t_120 vm.Value
	var case__8926_91 vm.Value
	var t_92 vm.Value
	var case__8926_93 vm.Value
	var t_94 vm.Value
	var v114 vm.Value
	var case__8926_115 vm.Value
	var t_116 vm.Value
	var case__8926_100 vm.Value
	var t_101 vm.Value
	var case__8926_102 vm.Value
	var t_103 vm.Value
	var v110 vm.Value
	var case__8926_111 vm.Value
	var t_112 vm.Value
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = v6, case__8926_1, t_2, case__8926_3, t_4, v15, v154, case__8926_155, t_156, case__8926_10, t_11, case__8926_12, t_13, v24, v150, case__8926_151, t_152, case__8926_19, t_20, case__8926_21, t_22, v33, v146, case__8926_147, t_148, case__8926_28, t_29, case__8926_30, t_31, v42, v142, case__8926_143, t_144, case__8926_37, t_38, case__8926_39, t_40, v51, v138, case__8926_139, t_140, case__8926_46, t_47, case__8926_48, t_49, v60, v134, case__8926_135, t_136, case__8926_55, t_56, case__8926_57, t_58, v69, v130, case__8926_131, t_132, case__8926_64, t_65, case__8926_66, t_67, v78, v126, case__8926_127, t_128, case__8926_73, t_74, case__8926_75, t_76, v87, v122, case__8926_123, t_124, case__8926_82, t_83, case__8926_84, t_85, v96, v118, case__8926_119, t_120, case__8926_91, t_92, case__8926_93, t_94, v114, case__8926_115, t_116, case__8926_100, t_101, case__8926_102, t_103, v110, case__8926_111, t_112
	v6 = arg0 == vm.Keyword("unknown")
	if v6 {
		case__8926_1 = arg0
		t_2 = arg0
		goto b1
	} else {
		case__8926_3 = arg0
		t_4 = arg0
		goto b2
	}
b1:
	;
	v154 = vm.String("unknown")
	case__8926_155 = case__8926_1
	t_156 = t_2
	goto b3
b2:
	;
	v15 = case__8926_3 == vm.Keyword("bottom")
	if v15 {
		case__8926_10 = case__8926_3
		t_11 = t_4
		goto b4
	} else {
		case__8926_12 = case__8926_3
		t_13 = t_4
		goto b5
	}
b3:
	;
	return v154
b4:
	;
	v150 = vm.String("bottom")
	case__8926_151 = case__8926_10
	t_152 = t_11
	goto b6
b5:
	;
	v24 = case__8926_12 == vm.Keyword("true")
	if v24 {
		case__8926_19 = case__8926_12
		t_20 = t_13
		goto b7
	} else {
		case__8926_21 = case__8926_12
		t_22 = t_13
		goto b8
	}
b6:
	;
	v154 = v150
	case__8926_155 = case__8926_151
	t_156 = t_152
	goto b3
b7:
	;
	v146 = vm.String("true")
	case__8926_147 = case__8926_19
	t_148 = t_20
	goto b9
b8:
	;
	v33 = case__8926_21 == vm.Keyword("false")
	if v33 {
		case__8926_28 = case__8926_21
		t_29 = t_22
		goto b10
	} else {
		case__8926_30 = case__8926_21
		t_31 = t_22
		goto b11
	}
b9:
	;
	v150 = v146
	case__8926_151 = case__8926_147
	t_152 = t_148
	goto b6
b10:
	;
	v142 = vm.String("false")
	case__8926_143 = case__8926_28
	t_144 = t_29
	goto b12
b11:
	;
	v42 = case__8926_30 == vm.Keyword("int")
	if v42 {
		case__8926_37 = case__8926_30
		t_38 = t_31
		goto b13
	} else {
		case__8926_39 = case__8926_30
		t_40 = t_31
		goto b14
	}
b12:
	;
	v146 = v142
	case__8926_147 = case__8926_143
	t_148 = t_144
	goto b9
b13:
	;
	v138 = vm.String("int")
	case__8926_139 = case__8926_37
	t_140 = t_38
	goto b15
b14:
	;
	v51 = case__8926_39 == vm.Keyword("float")
	if v51 {
		case__8926_46 = case__8926_39
		t_47 = t_40
		goto b16
	} else {
		case__8926_48 = case__8926_39
		t_49 = t_40
		goto b17
	}
b15:
	;
	v142 = v138
	case__8926_143 = case__8926_139
	t_144 = t_140
	goto b12
b16:
	;
	v134 = vm.String("float")
	case__8926_135 = case__8926_46
	t_136 = t_47
	goto b18
b17:
	;
	v60 = case__8926_48 == vm.Keyword("number")
	if v60 {
		case__8926_55 = case__8926_48
		t_56 = t_49
		goto b19
	} else {
		case__8926_57 = case__8926_48
		t_58 = t_49
		goto b20
	}
b18:
	;
	v138 = v134
	case__8926_139 = case__8926_135
	t_140 = t_136
	goto b15
b19:
	;
	v130 = vm.String("number")
	case__8926_131 = case__8926_55
	t_132 = t_56
	goto b21
b20:
	;
	v69 = case__8926_57 == vm.Keyword("bool")
	if v69 {
		case__8926_64 = case__8926_57
		t_65 = t_58
		goto b22
	} else {
		case__8926_66 = case__8926_57
		t_67 = t_58
		goto b23
	}
b21:
	;
	v134 = v130
	case__8926_135 = case__8926_131
	t_136 = t_132
	goto b18
b22:
	;
	v126 = vm.String("bool")
	case__8926_127 = case__8926_64
	t_128 = t_65
	goto b24
b23:
	;
	v78 = case__8926_66 == vm.Keyword("nil")
	if v78 {
		case__8926_73 = case__8926_66
		t_74 = t_67
		goto b25
	} else {
		case__8926_75 = case__8926_66
		t_76 = t_67
		goto b26
	}
b24:
	;
	v130 = v126
	case__8926_131 = case__8926_127
	t_132 = t_128
	goto b21
b25:
	;
	v122 = vm.String("nil")
	case__8926_123 = case__8926_73
	t_124 = t_74
	goto b27
b26:
	;
	v87 = case__8926_75 == vm.Keyword("string")
	if v87 {
		case__8926_82 = case__8926_75
		t_83 = t_76
		goto b28
	} else {
		case__8926_84 = case__8926_75
		t_85 = t_76
		goto b29
	}
b27:
	;
	v126 = v122
	case__8926_127 = case__8926_123
	t_128 = t_124
	goto b24
b28:
	;
	v118 = vm.String("string")
	case__8926_119 = case__8926_82
	t_120 = t_83
	goto b30
b29:
	;
	v96 = case__8926_84 == vm.Keyword("any")
	if v96 {
		case__8926_91 = case__8926_84
		t_92 = t_85
		goto b31
	} else {
		case__8926_93 = case__8926_84
		t_94 = t_85
		goto b32
	}
b30:
	;
	v122 = v118
	case__8926_123 = case__8926_119
	t_124 = t_120
	goto b27
b31:
	;
	v114 = vm.String("any")
	case__8926_115 = case__8926_91
	t_116 = t_92
	goto b33
b32:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		case__8926_100 = case__8926_93
		t_101 = t_94
		goto b34
	} else {
		case__8926_102 = case__8926_93
		t_103 = t_94
		goto b35
	}
b33:
	;
	v118 = v114
	case__8926_119 = case__8926_115
	t_120 = t_116
	goto b30
b34:
	;
	v110 = vm.String("??")
	case__8926_111 = case__8926_100
	t_112 = t_101
	goto b36
b35:
	;
	v110 = vm.NIL
	case__8926_111 = case__8926_102
	t_112 = t_103
	goto b36
b36:
	;
	v114 = v110
	case__8926_115 = case__8926_111
	t_116 = t_112
	goto b33
}
func type_display(arg0 vm.Value) (vm.Value, error) {
	var v5 vm.Value
	var t_2 vm.Value
	var v8 vm.Value
	var t_3 vm.Value
	var and__x_13 vm.Value
	var v236 vm.Value
	var t_237 vm.Value
	var t_10 vm.Value
	var arg__8978_34 vm.Value
	var arg__8991_38 vm.Value
	var arg__8992_39 vm.Value
	var arg__9006_44 vm.Value
	var arg__9019_48 vm.Value
	var arg__9020_49 vm.Value
	var arg__9021_50 vm.Value
	var arg__9036_56 vm.Value
	var arg__9049_60 vm.Value
	var arg__9050_61 vm.Value
	var arg__9064_66 vm.Value
	var arg__9077_70 vm.Value
	var arg__9078_71 vm.Value
	var arg__9079_72 vm.Value
	var arg__9080_73 vm.Value
	var arg__9097_81 vm.Value
	var arg__9110_85 vm.Value
	var arg__9111_86 vm.Value
	var arg__9125_91 vm.Value
	var arg__9138_95 vm.Value
	var arg__9139_96 vm.Value
	var arg__9140_97 vm.Value
	var arg__9155_103 vm.Value
	var arg__9168_107 vm.Value
	var arg__9169_108 vm.Value
	var arg__9183_113 vm.Value
	var arg__9196_117 vm.Value
	var arg__9197_118 vm.Value
	var arg__9198_119 vm.Value
	var arg__9199_120 vm.Value
	var v122 vm.Value
	var t_11 vm.Value
	var and__x_127 vm.Value
	var v233 vm.Value
	var t_234 vm.Value
	var t_14 vm.Value
	var and__x_15 vm.Value
	var arg__8963_21 vm.Value
	var v22 bool
	var t_16 vm.Value
	var and__x_17 vm.Value
	var v25 vm.Value
	var t_26 vm.Value
	var and__x_27 vm.Value
	var t_124 vm.Value
	var tag_146 vm.Value
	var v_150 vm.Value
	var v160 bool
	var t_125 vm.Value
	var v230 vm.Value
	var t_231 vm.Value
	var t_128 vm.Value
	var and__x_129 vm.Value
	var arg__9208_135 vm.Value
	var v136 bool
	var t_130 vm.Value
	var and__x_131 vm.Value
	var v139 vm.Value
	var t_140 vm.Value
	var and__x_141 vm.Value
	var t_151 vm.Value
	var tag_152 vm.Value
	var case__8949_153 vm.Value
	var v_154 vm.Value
	var v167 vm.Value
	var t_155 vm.Value
	var tag_156 vm.Value
	var case__8949_157 vm.Value
	var v_158 vm.Value
	var v178 bool
	var v213 vm.Value
	var t_214 vm.Value
	var tag_215 vm.Value
	var case__8949_216 vm.Value
	var v_217 vm.Value
	var t_169 vm.Value
	var tag_170 vm.Value
	var case__8949_171 vm.Value
	var v_172 vm.Value
	var v185 vm.Value
	var t_173 vm.Value
	var tag_174 vm.Value
	var case__8949_175 vm.Value
	var v_176 vm.Value
	var v207 vm.Value
	var t_208 vm.Value
	var tag_209 vm.Value
	var case__8949_210 vm.Value
	var v_211 vm.Value
	var t_187 vm.Value
	var tag_188 vm.Value
	var case__8949_189 vm.Value
	var v_190 vm.Value
	var t_191 vm.Value
	var tag_192 vm.Value
	var case__8949_193 vm.Value
	var v_194 vm.Value
	var v201 vm.Value
	var t_202 vm.Value
	var tag_203 vm.Value
	var case__8949_204 vm.Value
	var v_205 vm.Value
	var t_219 vm.Value
	var t_220 vm.Value
	var v227 vm.Value
	var t_228 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = v5, t_2, v8, t_3, and__x_13, v236, t_237, t_10, arg__8978_34, arg__8991_38, arg__8992_39, arg__9006_44, arg__9019_48, arg__9020_49, arg__9021_50, arg__9036_56, arg__9049_60, arg__9050_61, arg__9064_66, arg__9077_70, arg__9078_71, arg__9079_72, arg__9080_73, arg__9097_81, arg__9110_85, arg__9111_86, arg__9125_91, arg__9138_95, arg__9139_96, arg__9140_97, arg__9155_103, arg__9168_107, arg__9169_108, arg__9183_113, arg__9196_117, arg__9197_118, arg__9198_119, arg__9199_120, v122, t_11, and__x_127, v233, t_234, t_14, and__x_15, arg__8963_21, v22, t_16, and__x_17, v25, t_26, and__x_27, t_124, tag_146, v_150, v160, t_125, v230, t_231, t_128, and__x_129, arg__9208_135, v136, t_130, and__x_131, v139, t_140, and__x_141, t_151, tag_152, case__8949_153, v_154, v167, t_155, tag_156, case__8949_157, v_158, v178, v213, t_214, tag_215, case__8949_216, v_217, t_169, tag_170, case__8949_171, v_172, v185, t_173, tag_174, case__8949_175, v_176, v207, t_208, tag_209, case__8949_210, v_211, t_187, tag_188, case__8949_189, v_190, t_191, tag_192, case__8949_193, v_194, v201, t_202, tag_203, case__8949_204, v_205, t_219, t_220, v227, t_228
	v5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "keyword?").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v5) {
		t_2 = arg0
		goto b1
	} else {
		t_3 = arg0
		goto b2
	}
b1:
	;
	v8, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "scalar-type-display").Deref(), []vm.Value{t_2})
	if callErr != nil {
		return nil, callErr
	}
	v236 = v8
	t_237 = t_2
	goto b3
b2:
	;
	and__x_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector?").Deref(), []vm.Value{t_3})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(and__x_13) {
		t_14 = t_3
		and__x_15 = and__x_13
		goto b7
	} else {
		t_16 = t_3
		and__x_17 = and__x_13
		goto b8
	}
b3:
	;
	return v236, nil
b4:
	;
	arg__8978_34, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__8991_38, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__8992_39, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "sort-by").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v6 vm.Value
		var callErr error
		_ = v6
		v6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{rt.LookupVar("ir.dump", "type-order").Deref(), arg0, vm.Int(100)})
		if callErr != nil {
			return nil, callErr
		}
		return v6, nil
	}), arg__8991_38})
	if callErr != nil {
		return nil, callErr
	}
	arg__9006_44, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__9019_48, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__9020_49, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "sort-by").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v6 vm.Value
		var callErr error
		_ = v6
		v6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{rt.LookupVar("ir.dump", "type-order").Deref(), arg0, vm.Int(100)})
		if callErr != nil {
			return nil, callErr
		}
		return v6, nil
	}), arg__9019_48})
	if callErr != nil {
		return nil, callErr
	}
	arg__9021_50, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.LookupVar("ir.dump", "type-display").Deref(), arg__9020_49})
	if callErr != nil {
		return nil, callErr
	}
	arg__9036_56, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__9049_60, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__9050_61, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "sort-by").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v6 vm.Value
		var callErr error
		_ = v6
		v6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{rt.LookupVar("ir.dump", "type-order").Deref(), arg0, vm.Int(100)})
		if callErr != nil {
			return nil, callErr
		}
		return v6, nil
	}), arg__9049_60})
	if callErr != nil {
		return nil, callErr
	}
	arg__9064_66, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__9077_70, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__9078_71, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "sort-by").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v6 vm.Value
		var callErr error
		_ = v6
		v6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{rt.LookupVar("ir.dump", "type-order").Deref(), arg0, vm.Int(100)})
		if callErr != nil {
			return nil, callErr
		}
		return v6, nil
	}), arg__9077_70})
	if callErr != nil {
		return nil, callErr
	}
	arg__9079_72, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.LookupVar("ir.dump", "type-display").Deref(), arg__9078_71})
	if callErr != nil {
		return nil, callErr
	}
	arg__9080_73, callErr = rt.InvokeValue(rt.LookupVar("clojure.string", "join").Deref(), []vm.Value{vm.String(","), arg__9079_72})
	if callErr != nil {
		return nil, callErr
	}
	arg__9097_81, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__9110_85, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__9111_86, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "sort-by").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v6 vm.Value
		var callErr error
		_ = v6
		v6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{rt.LookupVar("ir.dump", "type-order").Deref(), arg0, vm.Int(100)})
		if callErr != nil {
			return nil, callErr
		}
		return v6, nil
	}), arg__9110_85})
	if callErr != nil {
		return nil, callErr
	}
	arg__9125_91, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__9138_95, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__9139_96, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "sort-by").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v6 vm.Value
		var callErr error
		_ = v6
		v6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{rt.LookupVar("ir.dump", "type-order").Deref(), arg0, vm.Int(100)})
		if callErr != nil {
			return nil, callErr
		}
		return v6, nil
	}), arg__9138_95})
	if callErr != nil {
		return nil, callErr
	}
	arg__9140_97, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.LookupVar("ir.dump", "type-display").Deref(), arg__9139_96})
	if callErr != nil {
		return nil, callErr
	}
	arg__9155_103, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__9168_107, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__9169_108, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "sort-by").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v6 vm.Value
		var callErr error
		_ = v6
		v6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{rt.LookupVar("ir.dump", "type-order").Deref(), arg0, vm.Int(100)})
		if callErr != nil {
			return nil, callErr
		}
		return v6, nil
	}), arg__9168_107})
	if callErr != nil {
		return nil, callErr
	}
	arg__9183_113, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__9196_117, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__9197_118, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "sort-by").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v6 vm.Value
		var callErr error
		_ = v6
		v6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{rt.LookupVar("ir.dump", "type-order").Deref(), arg0, vm.Int(100)})
		if callErr != nil {
			return nil, callErr
		}
		return v6, nil
	}), arg__9196_117})
	if callErr != nil {
		return nil, callErr
	}
	arg__9198_119, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.LookupVar("ir.dump", "type-display").Deref(), arg__9197_118})
	if callErr != nil {
		return nil, callErr
	}
	arg__9199_120, callErr = rt.InvokeValue(rt.LookupVar("clojure.string", "join").Deref(), []vm.Value{vm.String(","), arg__9198_119})
	if callErr != nil {
		return nil, callErr
	}
	v122, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("union{"), arg__9199_120, vm.String("}")})
	if callErr != nil {
		return nil, callErr
	}
	v233 = v122
	t_234 = t_10
	goto b6
b5:
	;
	and__x_127, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector?").Deref(), []vm.Value{t_11})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(and__x_127) {
		t_128 = t_11
		and__x_129 = and__x_127
		goto b13
	} else {
		t_130 = t_11
		and__x_131 = and__x_127
		goto b14
	}
b6:
	;
	v236 = v233
	t_237 = t_234
	goto b3
b7:
	;
	arg__8963_21, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{t_14})
	if callErr != nil {
		return nil, callErr
	}
	v22 = arg__8963_21 == vm.Keyword("union")
	v25 = vm.Boolean(v22)
	t_26 = t_14
	and__x_27 = and__x_15
	goto b9
b8:
	;
	v25 = and__x_17
	t_26 = t_16
	and__x_27 = and__x_17
	goto b9
b9:
	;
	if vm.IsTruthy(v25) {
		t_10 = t_26
		goto b4
	} else {
		t_11 = t_26
		goto b5
	}
b10:
	;
	tag_146, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{t_124, vm.Int(1)})
	if callErr != nil {
		return nil, callErr
	}
	v_150, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{t_124, vm.Int(2)})
	if callErr != nil {
		return nil, callErr
	}
	v160 = tag_146 == vm.Keyword("int")
	if v160 {
		t_151 = t_124
		tag_152 = tag_146
		case__8949_153 = tag_146
		v_154 = v_150
		goto b16
	} else {
		t_155 = t_124
		tag_156 = tag_146
		case__8949_157 = tag_146
		v_158 = v_150
		goto b17
	}
b11:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		t_219 = t_125
		goto b25
	} else {
		t_220 = t_125
		goto b26
	}
b12:
	;
	v233 = v230
	t_234 = t_231
	goto b6
b13:
	;
	arg__9208_135, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{t_128})
	if callErr != nil {
		return nil, callErr
	}
	v136 = arg__9208_135 == vm.Keyword("const")
	v139 = vm.Boolean(v136)
	t_140 = t_128
	and__x_141 = and__x_129
	goto b15
b14:
	;
	v139 = and__x_131
	t_140 = t_130
	and__x_141 = and__x_131
	goto b15
b15:
	;
	if vm.IsTruthy(v139) {
		t_124 = t_140
		goto b10
	} else {
		t_125 = t_140
		goto b11
	}
b16:
	;
	v167, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("int("), v_154, vm.String(")")})
	if callErr != nil {
		return nil, callErr
	}
	v213 = v167
	t_214 = t_151
	tag_215 = tag_152
	case__8949_216 = case__8949_153
	v_217 = v_154
	goto b18
b17:
	;
	v178 = case__8949_157 == vm.Keyword("float")
	if v178 {
		t_169 = t_155
		tag_170 = tag_156
		case__8949_171 = case__8949_157
		v_172 = v_158
		goto b19
	} else {
		t_173 = t_155
		tag_174 = tag_156
		case__8949_175 = case__8949_157
		v_176 = v_158
		goto b20
	}
b18:
	;
	v230 = v213
	t_231 = t_214
	goto b12
b19:
	;
	v185, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("float("), v_172, vm.String(")")})
	if callErr != nil {
		return nil, callErr
	}
	v207 = v185
	t_208 = t_169
	tag_209 = tag_170
	case__8949_210 = case__8949_171
	v_211 = v_172
	goto b21
b20:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		t_187 = t_173
		tag_188 = tag_174
		case__8949_189 = case__8949_175
		v_190 = v_176
		goto b22
	} else {
		t_191 = t_173
		tag_192 = tag_174
		case__8949_193 = case__8949_175
		v_194 = v_176
		goto b23
	}
b21:
	;
	v213 = v207
	t_214 = t_208
	tag_215 = tag_209
	case__8949_216 = case__8949_210
	v_217 = v_211
	goto b18
b22:
	;
	v201 = vm.String("??")
	t_202 = t_187
	tag_203 = tag_188
	case__8949_204 = case__8949_189
	v_205 = v_190
	goto b24
b23:
	;
	v201 = vm.NIL
	t_202 = t_191
	tag_203 = tag_192
	case__8949_204 = case__8949_193
	v_205 = v_194
	goto b24
b24:
	;
	v207 = v201
	t_208 = t_202
	tag_209 = tag_203
	case__8949_210 = case__8949_204
	v_211 = v_205
	goto b21
b25:
	;
	v227 = vm.String("??")
	t_228 = t_219
	goto b27
b26:
	;
	v227 = vm.NIL
	t_228 = t_220
	goto b27
b27:
	;
	v230 = v227
	t_231 = t_228
	goto b12
}
func op_display_name(arg0 vm.Value) (vm.Value, error) {
	var arg__9242_5 vm.Value
	var arg__9246_9 vm.Value
	var arg__9251_12 vm.Value
	var arg__9255_16 vm.Value
	var arg__9256_17 vm.Value
	var arg__9262_21 vm.Value
	var arg__9266_25 vm.Value
	var arg__9271_28 vm.Value
	var arg__9275_32 vm.Value
	var arg__9276_33 vm.Value
	var arg__9277_34 vm.Value
	var arg__9284_39 vm.Value
	var arg__9288_43 vm.Value
	var arg__9293_46 vm.Value
	var arg__9297_50 vm.Value
	var arg__9298_51 vm.Value
	var arg__9304_55 vm.Value
	var arg__9308_59 vm.Value
	var arg__9313_62 vm.Value
	var arg__9317_66 vm.Value
	var arg__9318_67 vm.Value
	var arg__9319_68 vm.Value
	var v69 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__9242_5, arg__9246_9, arg__9251_12, arg__9255_16, arg__9256_17, arg__9262_21, arg__9266_25, arg__9271_28, arg__9275_32, arg__9276_33, arg__9277_34, arg__9284_39, arg__9288_43, arg__9293_46, arg__9297_50, arg__9298_51, arg__9304_55, arg__9308_59, arg__9313_62, arg__9317_66, arg__9318_67, arg__9319_68, v69
	arg__9242_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "name").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__9246_9, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "re-pattern").Deref(), []vm.Value{vm.String("-")})
	if callErr != nil {
		return nil, callErr
	}
	arg__9251_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "name").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__9255_16, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "re-pattern").Deref(), []vm.Value{vm.String("-")})
	if callErr != nil {
		return nil, callErr
	}
	arg__9256_17, callErr = rt.InvokeValue(rt.LookupVar("clojure.string", "split").Deref(), []vm.Value{arg__9251_12, arg__9255_16})
	if callErr != nil {
		return nil, callErr
	}
	arg__9262_21, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "name").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__9266_25, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "re-pattern").Deref(), []vm.Value{vm.String("-")})
	if callErr != nil {
		return nil, callErr
	}
	arg__9271_28, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "name").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__9275_32, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "re-pattern").Deref(), []vm.Value{vm.String("-")})
	if callErr != nil {
		return nil, callErr
	}
	arg__9276_33, callErr = rt.InvokeValue(rt.LookupVar("clojure.string", "split").Deref(), []vm.Value{arg__9271_28, arg__9275_32})
	if callErr != nil {
		return nil, callErr
	}
	arg__9277_34, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.LookupVar("clojure.string", "capitalize").Deref(), arg__9276_33})
	if callErr != nil {
		return nil, callErr
	}
	arg__9284_39, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "name").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__9288_43, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "re-pattern").Deref(), []vm.Value{vm.String("-")})
	if callErr != nil {
		return nil, callErr
	}
	arg__9293_46, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "name").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__9297_50, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "re-pattern").Deref(), []vm.Value{vm.String("-")})
	if callErr != nil {
		return nil, callErr
	}
	arg__9298_51, callErr = rt.InvokeValue(rt.LookupVar("clojure.string", "split").Deref(), []vm.Value{arg__9293_46, arg__9297_50})
	if callErr != nil {
		return nil, callErr
	}
	arg__9304_55, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "name").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__9308_59, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "re-pattern").Deref(), []vm.Value{vm.String("-")})
	if callErr != nil {
		return nil, callErr
	}
	arg__9313_62, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "name").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__9317_66, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "re-pattern").Deref(), []vm.Value{vm.String("-")})
	if callErr != nil {
		return nil, callErr
	}
	arg__9318_67, callErr = rt.InvokeValue(rt.LookupVar("clojure.string", "split").Deref(), []vm.Value{arg__9313_62, arg__9317_66})
	if callErr != nil {
		return nil, callErr
	}
	arg__9319_68, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.LookupVar("clojure.string", "capitalize").Deref(), arg__9318_67})
	if callErr != nil {
		return nil, callErr
	}
	v69, callErr = rt.InvokeValue(rt.LookupVar("clojure.string", "join").Deref(), []vm.Value{vm.String(""), arg__9319_68})
	if callErr != nil {
		return nil, callErr
	}
	return v69, nil
}
func format_refs(arg0 vm.Value) (vm.Value, error) {
	var arg__9336_5 vm.Value
	var arg__9354_11 vm.Value
	var v12 vm.Value
	var callErr error
	_, _, _ = arg__9336_5, arg__9354_11, v12
	arg__9336_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v4 vm.Value
		var callErr error
		_ = v4
		v4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(" v"), arg0})
		if callErr != nil {
			return nil, callErr
		}
		return v4, nil
	}), arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__9354_11, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v4 vm.Value
		var callErr error
		_ = v4
		v4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(" v"), arg0})
		if callErr != nil {
			return nil, callErr
		}
		return v4, nil
	}), arg0})
	if callErr != nil {
		return nil, callErr
	}
	v12, callErr = rt.InvokeValue(rt.LookupVar("clojure.string", "join").Deref(), []vm.Value{vm.String(""), arg__9354_11})
	if callErr != nil {
		return nil, callErr
	}
	return v12, nil
}
func format_target(arg0 vm.Value) (vm.Value, error) {
	var arg__9359_3 vm.Value
	var arg__9364_6 vm.Value
	var arg__9369_9 vm.Value
	var arg__9370_10 vm.Value
	var arg__9377_15 vm.Value
	var arg__9382_18 vm.Value
	var arg__9387_21 vm.Value
	var arg__9388_22 vm.Value
	var v24 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _ = arg__9359_3, arg__9364_6, arg__9369_9, arg__9370_10, arg__9377_15, arg__9382_18, arg__9387_21, arg__9388_22, v24
	arg__9359_3, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-target").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__9364_6, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-args").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__9369_9, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-args").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__9370_10, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "format-args").Deref(), []vm.Value{arg__9369_9})
	if callErr != nil {
		return nil, callErr
	}
	arg__9377_15, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-target").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__9382_18, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-args").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__9387_21, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-args").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__9388_22, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "format-args").Deref(), []vm.Value{arg__9387_21})
	if callErr != nil {
		return nil, callErr
	}
	v24, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("b"), arg__9377_15, vm.String("("), arg__9388_22, vm.String(")")})
	if callErr != nil {
		return nil, callErr
	}
	return v24, nil
}
func terminator_targets_str(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var v7 bool
	var op_2 vm.Value
	var aux_3 vm.Value
	var arg__9396_11 vm.Value
	var arg__9402_15 vm.Value
	var v16 vm.Value
	var op_4 vm.Value
	var aux_5 vm.Value
	var v23 bool
	var v74 vm.Value
	var op_75 vm.Value
	var aux_76 vm.Value
	var op_18 vm.Value
	var aux_19 vm.Value
	var arg__9409_27 vm.Value
	var arg__9414_30 vm.Value
	var arg__9415_31 vm.Value
	var arg__9420_34 vm.Value
	var arg__9425_37 vm.Value
	var arg__9426_38 vm.Value
	var arg__9432_42 vm.Value
	var arg__9437_45 vm.Value
	var arg__9438_46 vm.Value
	var arg__9443_49 vm.Value
	var arg__9448_52 vm.Value
	var arg__9449_53 vm.Value
	var v54 vm.Value
	var op_20 vm.Value
	var aux_21 vm.Value
	var v70 vm.Value
	var op_71 vm.Value
	var aux_72 vm.Value
	var op_56 vm.Value
	var aux_57 vm.Value
	var op_58 vm.Value
	var aux_59 vm.Value
	var v66 vm.Value
	var op_67 vm.Value
	var aux_68 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = v7, op_2, aux_3, arg__9396_11, arg__9402_15, v16, op_4, aux_5, v23, v74, op_75, aux_76, op_18, aux_19, arg__9409_27, arg__9414_30, arg__9415_31, arg__9420_34, arg__9425_37, arg__9426_38, arg__9432_42, arg__9437_45, arg__9438_46, arg__9443_49, arg__9448_52, arg__9449_53, v54, op_20, aux_21, v70, op_71, aux_72, op_56, aux_57, op_58, aux_59, v66, op_67, aux_68
	v7 = arg0 == vm.Keyword("branch")
	if v7 {
		op_2 = arg0
		aux_3 = arg1
		goto b1
	} else {
		op_4 = arg0
		aux_5 = arg1
		goto b2
	}
b1:
	;
	arg__9396_11, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "format-target").Deref(), []vm.Value{aux_3})
	if callErr != nil {
		return nil, callErr
	}
	arg__9402_15, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "format-target").Deref(), []vm.Value{aux_3})
	if callErr != nil {
		return nil, callErr
	}
	v16, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(" -> "), arg__9402_15})
	if callErr != nil {
		return nil, callErr
	}
	v74 = v16
	op_75 = op_2
	aux_76 = aux_3
	goto b3
b2:
	;
	v23 = op_4 == vm.Keyword("branch-if")
	if v23 {
		op_18 = op_4
		aux_19 = aux_5
		goto b4
	} else {
		op_20 = op_4
		aux_21 = aux_5
		goto b5
	}
b3:
	;
	return v74, nil
b4:
	;
	arg__9409_27, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-true").Deref(), []vm.Value{aux_19})
	if callErr != nil {
		return nil, callErr
	}
	arg__9414_30, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-true").Deref(), []vm.Value{aux_19})
	if callErr != nil {
		return nil, callErr
	}
	arg__9415_31, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "format-target").Deref(), []vm.Value{arg__9414_30})
	if callErr != nil {
		return nil, callErr
	}
	arg__9420_34, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-false").Deref(), []vm.Value{aux_19})
	if callErr != nil {
		return nil, callErr
	}
	arg__9425_37, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-false").Deref(), []vm.Value{aux_19})
	if callErr != nil {
		return nil, callErr
	}
	arg__9426_38, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "format-target").Deref(), []vm.Value{arg__9425_37})
	if callErr != nil {
		return nil, callErr
	}
	arg__9432_42, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-true").Deref(), []vm.Value{aux_19})
	if callErr != nil {
		return nil, callErr
	}
	arg__9437_45, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-true").Deref(), []vm.Value{aux_19})
	if callErr != nil {
		return nil, callErr
	}
	arg__9438_46, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "format-target").Deref(), []vm.Value{arg__9437_45})
	if callErr != nil {
		return nil, callErr
	}
	arg__9443_49, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-false").Deref(), []vm.Value{aux_19})
	if callErr != nil {
		return nil, callErr
	}
	arg__9448_52, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-false").Deref(), []vm.Value{aux_19})
	if callErr != nil {
		return nil, callErr
	}
	arg__9449_53, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "format-target").Deref(), []vm.Value{arg__9448_52})
	if callErr != nil {
		return nil, callErr
	}
	v54, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(" -> "), arg__9438_46, vm.String(" : "), arg__9449_53})
	if callErr != nil {
		return nil, callErr
	}
	v70 = v54
	op_71 = op_18
	aux_72 = aux_19
	goto b6
b5:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		op_56 = op_20
		aux_57 = aux_21
		goto b7
	} else {
		op_58 = op_20
		aux_59 = aux_21
		goto b8
	}
b6:
	;
	v74 = v70
	op_75 = op_71
	aux_76 = aux_72
	goto b3
b7:
	;
	v66 = vm.String("")
	op_67 = op_56
	aux_68 = aux_57
	goto b9
b8:
	;
	v66 = vm.NIL
	op_67 = op_58
	aux_68 = aux_59
	goto b9
b9:
	;
	v70 = v66
	op_71 = op_67
	aux_72 = aux_68
	goto b6
}
func write_node(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var op_3 vm.Value
	var refs_5 vm.Value
	var aux_7 vm.Value
	var v21 vm.Value
	var f_8 vm.Value
	var id_9 vm.Value
	var op_10 vm.Value
	var refs_11 vm.Value
	var aux_12 vm.Value
	var arg__9474_25 vm.Value
	var arg__9478_27 vm.Value
	var arg__9484_29 vm.Value
	var arg__9491_34 vm.Value
	var arg__9495_36 vm.Value
	var arg__9501_38 vm.Value
	var v40 vm.Value
	var f_13 vm.Value
	var id_14 vm.Value
	var op_15 vm.Value
	var refs_16 vm.Value
	var aux_17 vm.Value
	var t_43 vm.Value
	var arg__9514_47 vm.Value
	var arg__9518_49 vm.Value
	var arg__9522_73 vm.Value
	var arg__9527_76 vm.Value
	var v77 vm.Value
	var v269 vm.Value
	var f_270 vm.Value
	var id_271 vm.Value
	var op_272 vm.Value
	var refs_273 vm.Value
	var aux_274 vm.Value
	var f_50 vm.Value
	var id_51 vm.Value
	var arg__9509_52 vm.Value
	var op_53 vm.Value
	var refs_54 vm.Value
	var aux_55 vm.Value
	var t_56 vm.Value
	var arg__9508_57 string
	var arg__9510_58 string
	var arg__9514_59 vm.Value
	var arg__9518_60 vm.Value
	var v82 vm.Value
	var f_61 vm.Value
	var id_62 vm.Value
	var arg__9509_63 vm.Value
	var op_64 vm.Value
	var refs_65 vm.Value
	var aux_66 vm.Value
	var t_67 vm.Value
	var arg__9508_68 string
	var arg__9510_69 string
	var arg__9514_70 vm.Value
	var arg__9518_71 vm.Value
	var arg__9533_86 vm.Value
	var f_87 vm.Value
	var id_88 vm.Value
	var arg__9509_89 vm.Value
	var op_90 vm.Value
	var refs_91 vm.Value
	var aux_92 vm.Value
	var t_93 vm.Value
	var arg__9508_94 string
	var arg__9510_95 string
	var arg__9514_96 vm.Value
	var arg__9518_97 vm.Value
	var v125 vm.Value
	var arg__9533_98 vm.Value
	var f_99 vm.Value
	var id_100 vm.Value
	var arg__9509_101 vm.Value
	var op_102 vm.Value
	var refs_103 vm.Value
	var aux_104 vm.Value
	var t_105 vm.Value
	var arg__9508_106 string
	var arg__9510_107 string
	var arg__9514_108 vm.Value
	var arg__9518_109 vm.Value
	var arg__9543_129 vm.Value
	var arg__9549_133 vm.Value
	var v134 vm.Value
	var arg__9533_110 vm.Value
	var f_111 vm.Value
	var id_112 vm.Value
	var arg__9509_113 vm.Value
	var op_114 vm.Value
	var refs_115 vm.Value
	var aux_116 vm.Value
	var t_117 vm.Value
	var arg__9508_118 string
	var arg__9510_119 string
	var arg__9514_120 vm.Value
	var arg__9518_121 vm.Value
	var arg__9550_138 vm.Value
	var arg__9533_139 vm.Value
	var f_140 vm.Value
	var id_141 vm.Value
	var arg__9509_142 vm.Value
	var op_143 vm.Value
	var refs_144 vm.Value
	var aux_145 vm.Value
	var t_146 vm.Value
	var arg__9508_147 string
	var arg__9510_148 string
	var arg__9514_149 vm.Value
	var arg__9518_150 vm.Value
	var arg__9559_156 vm.Value
	var arg__9563_158 vm.Value
	var arg__9567_184 vm.Value
	var arg__9572_187 vm.Value
	var v188 vm.Value
	var f_159 vm.Value
	var arg__9554_160 vm.Value
	var id_161 vm.Value
	var op_162 vm.Value
	var refs_163 vm.Value
	var aux_164 vm.Value
	var t_165 vm.Value
	var head__9552_166 vm.Value
	var arg__9553_167 string
	var arg__9555_168 string
	var arg__9559_169 vm.Value
	var arg__9563_170 vm.Value
	var v193 vm.Value
	var f_171 vm.Value
	var arg__9554_172 vm.Value
	var id_173 vm.Value
	var op_174 vm.Value
	var refs_175 vm.Value
	var aux_176 vm.Value
	var t_177 vm.Value
	var head__9552_178 vm.Value
	var arg__9553_179 string
	var arg__9555_180 string
	var arg__9559_181 vm.Value
	var arg__9563_182 vm.Value
	var arg__9578_197 vm.Value
	var f_198 vm.Value
	var arg__9554_199 vm.Value
	var id_200 vm.Value
	var op_201 vm.Value
	var refs_202 vm.Value
	var aux_203 vm.Value
	var t_204 vm.Value
	var head__9552_205 vm.Value
	var arg__9553_206 string
	var arg__9555_207 string
	var arg__9559_208 vm.Value
	var arg__9563_209 vm.Value
	var v239 vm.Value
	var arg__9578_210 vm.Value
	var f_211 vm.Value
	var arg__9554_212 vm.Value
	var id_213 vm.Value
	var op_214 vm.Value
	var refs_215 vm.Value
	var aux_216 vm.Value
	var t_217 vm.Value
	var head__9552_218 vm.Value
	var arg__9553_219 string
	var arg__9555_220 string
	var arg__9559_221 vm.Value
	var arg__9563_222 vm.Value
	var arg__9588_243 vm.Value
	var arg__9594_247 vm.Value
	var v248 vm.Value
	var arg__9578_223 vm.Value
	var f_224 vm.Value
	var arg__9554_225 vm.Value
	var id_226 vm.Value
	var op_227 vm.Value
	var refs_228 vm.Value
	var aux_229 vm.Value
	var t_230 vm.Value
	var head__9552_231 vm.Value
	var arg__9553_232 string
	var arg__9555_233 string
	var arg__9559_234 vm.Value
	var arg__9563_235 vm.Value
	var arg__9595_252 vm.Value
	var arg__9578_253 vm.Value
	var f_254 vm.Value
	var arg__9554_255 vm.Value
	var id_256 vm.Value
	var op_257 vm.Value
	var refs_258 vm.Value
	var aux_259 vm.Value
	var t_260 vm.Value
	var head__9552_261 vm.Value
	var arg__9553_262 string
	var arg__9555_263 string
	var arg__9559_264 vm.Value
	var arg__9563_265 vm.Value
	var v267 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = op_3, refs_5, aux_7, v21, f_8, id_9, op_10, refs_11, aux_12, arg__9474_25, arg__9478_27, arg__9484_29, arg__9491_34, arg__9495_36, arg__9501_38, v40, f_13, id_14, op_15, refs_16, aux_17, t_43, arg__9514_47, arg__9518_49, arg__9522_73, arg__9527_76, v77, v269, f_270, id_271, op_272, refs_273, aux_274, f_50, id_51, arg__9509_52, op_53, refs_54, aux_55, t_56, arg__9508_57, arg__9510_58, arg__9514_59, arg__9518_60, v82, f_61, id_62, arg__9509_63, op_64, refs_65, aux_66, t_67, arg__9508_68, arg__9510_69, arg__9514_70, arg__9518_71, arg__9533_86, f_87, id_88, arg__9509_89, op_90, refs_91, aux_92, t_93, arg__9508_94, arg__9510_95, arg__9514_96, arg__9518_97, v125, arg__9533_98, f_99, id_100, arg__9509_101, op_102, refs_103, aux_104, t_105, arg__9508_106, arg__9510_107, arg__9514_108, arg__9518_109, arg__9543_129, arg__9549_133, v134, arg__9533_110, f_111, id_112, arg__9509_113, op_114, refs_115, aux_116, t_117, arg__9508_118, arg__9510_119, arg__9514_120, arg__9518_121, arg__9550_138, arg__9533_139, f_140, id_141, arg__9509_142, op_143, refs_144, aux_145, t_146, arg__9508_147, arg__9510_148, arg__9514_149, arg__9518_150, arg__9559_156, arg__9563_158, arg__9567_184, arg__9572_187, v188, f_159, arg__9554_160, id_161, op_162, refs_163, aux_164, t_165, head__9552_166, arg__9553_167, arg__9555_168, arg__9559_169, arg__9563_170, v193, f_171, arg__9554_172, id_173, op_174, refs_175, aux_176, t_177, head__9552_178, arg__9553_179, arg__9555_180, arg__9559_181, arg__9563_182, arg__9578_197, f_198, arg__9554_199, id_200, op_201, refs_202, aux_203, t_204, head__9552_205, arg__9553_206, arg__9555_207, arg__9559_208, arg__9563_209, v239, arg__9578_210, f_211, arg__9554_212, id_213, op_214, refs_215, aux_216, t_217, head__9552_218, arg__9553_219, arg__9555_220, arg__9559_221, arg__9563_222, arg__9588_243, arg__9594_247, v248, arg__9578_223, f_224, arg__9554_225, id_226, op_227, refs_228, aux_229, t_230, head__9552_231, arg__9553_232, arg__9555_233, arg__9559_234, arg__9563_235, arg__9595_252, arg__9578_253, f_254, arg__9554_255, id_256, op_257, refs_258, aux_259, t_260, head__9552_261, arg__9553_262, arg__9555_263, arg__9559_264, arg__9563_265, v267
	op_3, callErr = rt.InvokeValue(rt.LookupVar("ir", "op").Deref(), []vm.Value{arg1, arg0})
	if callErr != nil {
		return nil, callErr
	}
	refs_5, callErr = rt.InvokeValue(rt.LookupVar("ir", "refs").Deref(), []vm.Value{arg1, arg0})
	if callErr != nil {
		return nil, callErr
	}
	aux_7, callErr = rt.InvokeValue(rt.LookupVar("ir", "aux").Deref(), []vm.Value{arg1, arg0})
	if callErr != nil {
		return nil, callErr
	}
	v21, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "contains?").Deref(), []vm.Value{rt.LookupVar("ir.dump", "terminator-ops").Deref(), op_3})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v21) {
		f_8 = arg0
		id_9 = arg1
		op_10 = op_3
		refs_11 = refs_5
		aux_12 = aux_7
		goto b1
	} else {
		f_13 = arg0
		id_14 = arg1
		op_15 = op_3
		refs_16 = refs_5
		aux_17 = aux_7
		goto b2
	}
b1:
	;
	arg__9474_25, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "op-display-name").Deref(), []vm.Value{op_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__9478_27, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "format-refs").Deref(), []vm.Value{refs_11})
	if callErr != nil {
		return nil, callErr
	}
	arg__9484_29, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "terminator-targets-str").Deref(), []vm.Value{op_10, aux_12})
	if callErr != nil {
		return nil, callErr
	}
	arg__9491_34, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "op-display-name").Deref(), []vm.Value{op_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__9495_36, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "format-refs").Deref(), []vm.Value{refs_11})
	if callErr != nil {
		return nil, callErr
	}
	arg__9501_38, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "terminator-targets-str").Deref(), []vm.Value{op_10, aux_12})
	if callErr != nil {
		return nil, callErr
	}
	v40, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("    "), arg__9491_34, arg__9495_36, arg__9501_38, vm.String("\n")})
	if callErr != nil {
		return nil, callErr
	}
	v269 = v40
	f_270 = f_8
	id_271 = id_9
	op_272 = op_10
	refs_273 = refs_11
	aux_274 = aux_12
	goto b3
b2:
	;
	t_43, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{id_14, f_13})
	if callErr != nil {
		return nil, callErr
	}
	arg__9514_47, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "op-display-name").Deref(), []vm.Value{op_15})
	if callErr != nil {
		return nil, callErr
	}
	arg__9518_49, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "format-refs").Deref(), []vm.Value{refs_16})
	if callErr != nil {
		return nil, callErr
	}
	arg__9522_73, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{aux_17})
	if callErr != nil {
		return nil, callErr
	}
	arg__9527_76, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{aux_17})
	if callErr != nil {
		return nil, callErr
	}
	v77, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not").Deref(), []vm.Value{arg__9527_76})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v77) {
		f_50 = f_13
		id_51 = id_14
		arg__9509_52 = id_14
		op_53 = op_15
		refs_54 = refs_16
		aux_55 = aux_17
		t_56 = t_43
		arg__9508_57 = "    v"
		arg__9510_58 = " = "
		arg__9514_59 = arg__9514_47
		arg__9518_60 = arg__9518_49
		goto b4
	} else {
		f_61 = f_13
		id_62 = id_14
		arg__9509_63 = id_14
		op_64 = op_15
		refs_65 = refs_16
		aux_66 = aux_17
		t_67 = t_43
		arg__9508_68 = "    v"
		arg__9510_69 = " = "
		arg__9514_70 = arg__9514_47
		arg__9518_71 = arg__9518_49
		goto b5
	}
b3:
	;
	return v269, nil
b4:
	;
	v82, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(" ; "), aux_55})
	if callErr != nil {
		return nil, callErr
	}
	arg__9533_86 = v82
	f_87 = f_50
	id_88 = id_51
	arg__9509_89 = arg__9509_52
	op_90 = op_53
	refs_91 = refs_54
	aux_92 = aux_55
	t_93 = t_56
	arg__9508_94 = arg__9508_57
	arg__9510_95 = arg__9510_58
	arg__9514_96 = arg__9514_59
	arg__9518_97 = arg__9518_60
	goto b6
b5:
	;
	arg__9533_86 = vm.String("")
	f_87 = f_61
	id_88 = id_62
	arg__9509_89 = arg__9509_63
	op_90 = op_64
	refs_91 = refs_65
	aux_92 = aux_66
	t_93 = t_67
	arg__9508_94 = arg__9508_68
	arg__9510_95 = arg__9510_69
	arg__9514_96 = arg__9514_70
	arg__9518_97 = arg__9518_71
	goto b6
b6:
	;
	v125, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not=").Deref(), []vm.Value{t_93, vm.Keyword("unknown")})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v125) {
		arg__9533_98 = arg__9533_86
		f_99 = f_87
		id_100 = id_88
		arg__9509_101 = arg__9509_89
		op_102 = op_90
		refs_103 = refs_91
		aux_104 = aux_92
		t_105 = t_93
		arg__9508_106 = arg__9508_94
		arg__9510_107 = arg__9510_95
		arg__9514_108 = arg__9514_96
		arg__9518_109 = arg__9518_97
		goto b7
	} else {
		arg__9533_110 = arg__9533_86
		f_111 = f_87
		id_112 = id_88
		arg__9509_113 = arg__9509_89
		op_114 = op_90
		refs_115 = refs_91
		aux_116 = aux_92
		t_117 = t_93
		arg__9508_118 = arg__9508_94
		arg__9510_119 = arg__9510_95
		arg__9514_120 = arg__9514_96
		arg__9518_121 = arg__9518_97
		goto b8
	}
b7:
	;
	arg__9543_129, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "type-display").Deref(), []vm.Value{t_105})
	if callErr != nil {
		return nil, callErr
	}
	arg__9549_133, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "type-display").Deref(), []vm.Value{t_105})
	if callErr != nil {
		return nil, callErr
	}
	v134, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(" : "), arg__9549_133})
	if callErr != nil {
		return nil, callErr
	}
	arg__9550_138 = v134
	arg__9533_139 = arg__9533_98
	f_140 = f_99
	id_141 = id_100
	arg__9509_142 = arg__9509_101
	op_143 = op_102
	refs_144 = refs_103
	aux_145 = aux_104
	t_146 = t_105
	arg__9508_147 = arg__9508_106
	arg__9510_148 = arg__9510_107
	arg__9514_149 = arg__9514_108
	arg__9518_150 = arg__9518_109
	goto b9
b8:
	;
	arg__9550_138 = vm.String("")
	arg__9533_139 = arg__9533_110
	f_140 = f_111
	id_141 = id_112
	arg__9509_142 = arg__9509_113
	op_143 = op_114
	refs_144 = refs_115
	aux_145 = aux_116
	t_146 = t_117
	arg__9508_147 = arg__9508_118
	arg__9510_148 = arg__9510_119
	arg__9514_149 = arg__9514_120
	arg__9518_150 = arg__9518_121
	goto b9
b9:
	;
	arg__9559_156, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "op-display-name").Deref(), []vm.Value{op_143})
	if callErr != nil {
		return nil, callErr
	}
	arg__9563_158, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "format-refs").Deref(), []vm.Value{refs_144})
	if callErr != nil {
		return nil, callErr
	}
	arg__9567_184, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{aux_145})
	if callErr != nil {
		return nil, callErr
	}
	arg__9572_187, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{aux_145})
	if callErr != nil {
		return nil, callErr
	}
	v188, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not").Deref(), []vm.Value{arg__9572_187})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v188) {
		f_159 = f_140
		arg__9554_160 = id_141
		id_161 = id_141
		op_162 = op_143
		refs_163 = refs_144
		aux_164 = aux_145
		t_165 = t_146
		head__9552_166 = rt.LookupVar("clojure.core", "str").Deref()
		arg__9553_167 = "    v"
		arg__9555_168 = " = "
		arg__9559_169 = arg__9559_156
		arg__9563_170 = arg__9563_158
		goto b10
	} else {
		f_171 = f_140
		arg__9554_172 = id_141
		id_173 = id_141
		op_174 = op_143
		refs_175 = refs_144
		aux_176 = aux_145
		t_177 = t_146
		head__9552_178 = rt.LookupVar("clojure.core", "str").Deref()
		arg__9553_179 = "    v"
		arg__9555_180 = " = "
		arg__9559_181 = arg__9559_156
		arg__9563_182 = arg__9563_158
		goto b11
	}
b10:
	;
	v193, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(" ; "), aux_164})
	if callErr != nil {
		return nil, callErr
	}
	arg__9578_197 = v193
	f_198 = f_159
	arg__9554_199 = arg__9554_160
	id_200 = id_161
	op_201 = op_162
	refs_202 = refs_163
	aux_203 = aux_164
	t_204 = t_165
	head__9552_205 = head__9552_166
	arg__9553_206 = arg__9553_167
	arg__9555_207 = arg__9555_168
	arg__9559_208 = arg__9559_169
	arg__9563_209 = arg__9563_170
	goto b12
b11:
	;
	arg__9578_197 = vm.String("")
	f_198 = f_171
	arg__9554_199 = arg__9554_172
	id_200 = id_173
	op_201 = op_174
	refs_202 = refs_175
	aux_203 = aux_176
	t_204 = t_177
	head__9552_205 = head__9552_178
	arg__9553_206 = arg__9553_179
	arg__9555_207 = arg__9555_180
	arg__9559_208 = arg__9559_181
	arg__9563_209 = arg__9563_182
	goto b12
b12:
	;
	v239, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not=").Deref(), []vm.Value{t_204, vm.Keyword("unknown")})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v239) {
		arg__9578_210 = arg__9578_197
		f_211 = f_198
		arg__9554_212 = arg__9554_199
		id_213 = id_200
		op_214 = op_201
		refs_215 = refs_202
		aux_216 = aux_203
		t_217 = t_204
		head__9552_218 = head__9552_205
		arg__9553_219 = arg__9553_206
		arg__9555_220 = arg__9555_207
		arg__9559_221 = arg__9559_208
		arg__9563_222 = arg__9563_209
		goto b13
	} else {
		arg__9578_223 = arg__9578_197
		f_224 = f_198
		arg__9554_225 = arg__9554_199
		id_226 = id_200
		op_227 = op_201
		refs_228 = refs_202
		aux_229 = aux_203
		t_230 = t_204
		head__9552_231 = head__9552_205
		arg__9553_232 = arg__9553_206
		arg__9555_233 = arg__9555_207
		arg__9559_234 = arg__9559_208
		arg__9563_235 = arg__9563_209
		goto b14
	}
b13:
	;
	arg__9588_243, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "type-display").Deref(), []vm.Value{t_217})
	if callErr != nil {
		return nil, callErr
	}
	arg__9594_247, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "type-display").Deref(), []vm.Value{t_217})
	if callErr != nil {
		return nil, callErr
	}
	v248, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(" : "), arg__9594_247})
	if callErr != nil {
		return nil, callErr
	}
	arg__9595_252 = v248
	arg__9578_253 = arg__9578_210
	f_254 = f_211
	arg__9554_255 = arg__9554_212
	id_256 = id_213
	op_257 = op_214
	refs_258 = refs_215
	aux_259 = aux_216
	t_260 = t_217
	head__9552_261 = head__9552_218
	arg__9553_262 = arg__9553_219
	arg__9555_263 = arg__9555_220
	arg__9559_264 = arg__9559_221
	arg__9563_265 = arg__9563_222
	goto b15
b14:
	;
	arg__9595_252 = vm.String("")
	arg__9578_253 = arg__9578_223
	f_254 = f_224
	arg__9554_255 = arg__9554_225
	id_256 = id_226
	op_257 = op_227
	refs_258 = refs_228
	aux_259 = aux_229
	t_260 = t_230
	head__9552_261 = head__9552_231
	arg__9553_262 = arg__9553_232
	arg__9555_263 = arg__9555_233
	arg__9559_264 = arg__9559_234
	arg__9563_265 = arg__9563_235
	goto b15
b15:
	;
	v267, callErr = rt.InvokeValue(head__9552_261, []vm.Value{vm.String(arg__9553_262), arg__9554_255, vm.String(arg__9555_263), arg__9559_264, arg__9563_265, arg__9578_253, arg__9595_252, vm.String("\n")})
	if callErr != nil {
		return nil, callErr
	}
	v269 = v267
	f_270 = f_254
	id_271 = id_256
	op_272 = op_257
	refs_273 = refs_258
	aux_274 = aux_259
	goto b3
}
func write_block(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var params_3 vm.Value
	var preds_5 vm.Value
	var insts_7 vm.Value
	var term_9 vm.Value
	var arg__9621_11 vm.Value
	var entry_QMARK__12 bool
	var f_14 vm.Value
	var bid_15 vm.Value
	var params_16 vm.Value
	var preds_17 vm.Value
	var insts_18 vm.Value
	var term_19 vm.Value
	var entry_QMARK__20 bool
	var arg__9622_21 string
	var f_22 vm.Value
	var bid_23 vm.Value
	var params_24 vm.Value
	var preds_25 vm.Value
	var insts_26 vm.Value
	var term_27 vm.Value
	var entry_QMARK__28 bool
	var arg__9622_29 string
	var arg__9623_35 string
	var f_36 vm.Value
	var bid_37 vm.Value
	var params_38 vm.Value
	var preds_39 vm.Value
	var insts_40 vm.Value
	var term_41 vm.Value
	var entry_QMARK__42 vm.Value
	var arg__9622_43 string
	var arg__9703_54 vm.Value
	var arg__9781_64 vm.Value
	var arg__9782_65 vm.Value
	var v96 vm.Value
	var arg__9623_67 string
	var f_68 vm.Value
	var arg__9625_69 vm.Value
	var bid_70 vm.Value
	var params_71 vm.Value
	var preds_72 vm.Value
	var insts_73 vm.Value
	var term_74 vm.Value
	var entry_QMARK__75 vm.Value
	var arg__9622_76 string
	var arg__9624_77 string
	var arg__9626_78 string
	var arg__9782_79 vm.Value
	var arg__9783_80 string
	var arg__9804_103 vm.Value
	var arg__9822_109 vm.Value
	var arg__9823_110 vm.Value
	var arg__9842_117 vm.Value
	var arg__9860_123 vm.Value
	var arg__9861_124 vm.Value
	var v125 vm.Value
	var arg__9623_81 string
	var f_82 vm.Value
	var arg__9625_83 vm.Value
	var bid_84 vm.Value
	var params_85 vm.Value
	var preds_86 vm.Value
	var insts_87 vm.Value
	var term_88 vm.Value
	var entry_QMARK__89 vm.Value
	var arg__9622_90 string
	var arg__9624_91 string
	var arg__9626_92 string
	var arg__9782_93 vm.Value
	var arg__9783_94 string
	var arg__9862_129 vm.Value
	var arg__9623_130 string
	var f_131 vm.Value
	var arg__9625_132 vm.Value
	var bid_133 vm.Value
	var params_134 vm.Value
	var preds_135 vm.Value
	var insts_136 vm.Value
	var term_137 vm.Value
	var entry_QMARK__138 vm.Value
	var arg__9622_139 string
	var arg__9624_140 string
	var arg__9626_141 string
	var arg__9782_142 vm.Value
	var arg__9783_143 string
	var f_147 vm.Value
	var bid_148 vm.Value
	var params_149 vm.Value
	var preds_150 vm.Value
	var insts_151 vm.Value
	var term_152 vm.Value
	var entry_QMARK__153 bool
	var head__9864_154 vm.Value
	var arg__9865_155 string
	var f_156 vm.Value
	var bid_157 vm.Value
	var params_158 vm.Value
	var preds_159 vm.Value
	var insts_160 vm.Value
	var term_161 vm.Value
	var entry_QMARK__162 bool
	var head__9864_163 vm.Value
	var arg__9865_164 string
	var arg__9866_170 string
	var f_171 vm.Value
	var bid_172 vm.Value
	var params_173 vm.Value
	var preds_174 vm.Value
	var insts_175 vm.Value
	var term_176 vm.Value
	var entry_QMARK__177 vm.Value
	var head__9864_178 vm.Value
	var arg__9865_179 string
	var arg__9946_190 vm.Value
	var arg__10024_200 vm.Value
	var arg__10025_201 vm.Value
	var v234 vm.Value
	var arg__9866_203 string
	var f_204 vm.Value
	var arg__9868_205 vm.Value
	var bid_206 vm.Value
	var params_207 vm.Value
	var preds_208 vm.Value
	var insts_209 vm.Value
	var term_210 vm.Value
	var entry_QMARK__211 vm.Value
	var head__9864_212 vm.Value
	var arg__9865_213 string
	var arg__9867_214 string
	var arg__9869_215 string
	var arg__10025_216 vm.Value
	var arg__10026_217 string
	var arg__10047_241 vm.Value
	var arg__10065_247 vm.Value
	var arg__10066_248 vm.Value
	var arg__10085_255 vm.Value
	var arg__10103_261 vm.Value
	var arg__10104_262 vm.Value
	var v263 vm.Value
	var arg__9866_218 string
	var f_219 vm.Value
	var arg__9868_220 vm.Value
	var bid_221 vm.Value
	var params_222 vm.Value
	var preds_223 vm.Value
	var insts_224 vm.Value
	var term_225 vm.Value
	var entry_QMARK__226 vm.Value
	var head__9864_227 vm.Value
	var arg__9865_228 string
	var arg__9867_229 string
	var arg__9869_230 string
	var arg__10025_231 vm.Value
	var arg__10026_232 string
	var arg__10105_267 vm.Value
	var arg__9866_268 string
	var f_269 vm.Value
	var arg__9868_270 vm.Value
	var bid_271 vm.Value
	var params_272 vm.Value
	var preds_273 vm.Value
	var insts_274 vm.Value
	var term_275 vm.Value
	var entry_QMARK__276 vm.Value
	var head__9864_277 vm.Value
	var arg__9865_278 string
	var arg__9867_279 string
	var arg__9869_280 string
	var arg__10025_281 vm.Value
	var arg__10026_282 string
	var header_284 vm.Value
	var arg__10123_293 vm.Value
	var arg__10141_303 vm.Value
	var body_304 vm.Value
	var or__x_326 vm.Value
	var f_305 vm.Value
	var bid_306 vm.Value
	var params_307 vm.Value
	var preds_308 vm.Value
	var insts_309 vm.Value
	var term_310 vm.Value
	var entry_QMARK__311 vm.Value
	var header_312 vm.Value
	var body_313 vm.Value
	var v365 vm.Value
	var f_314 vm.Value
	var bid_315 vm.Value
	var params_316 vm.Value
	var preds_317 vm.Value
	var insts_318 vm.Value
	var term_319 vm.Value
	var entry_QMARK__320 vm.Value
	var header_321 vm.Value
	var body_322 vm.Value
	var term_line_369 vm.Value
	var f_370 vm.Value
	var bid_371 vm.Value
	var params_372 vm.Value
	var preds_373 vm.Value
	var insts_374 vm.Value
	var term_375 vm.Value
	var entry_QMARK__376 vm.Value
	var header_377 vm.Value
	var body_378 vm.Value
	var v382 vm.Value
	var f_327 vm.Value
	var bid_328 vm.Value
	var params_329 vm.Value
	var preds_330 vm.Value
	var insts_331 vm.Value
	var term_332 vm.Value
	var entry_QMARK__333 vm.Value
	var header_334 vm.Value
	var body_335 vm.Value
	var or__x_336 vm.Value
	var f_337 vm.Value
	var bid_338 vm.Value
	var params_339 vm.Value
	var preds_340 vm.Value
	var insts_341 vm.Value
	var term_342 vm.Value
	var entry_QMARK__343 vm.Value
	var header_344 vm.Value
	var body_345 vm.Value
	var or__x_346 vm.Value
	var v350 vm.Value
	var v352 vm.Value
	var f_353 vm.Value
	var bid_354 vm.Value
	var params_355 vm.Value
	var preds_356 vm.Value
	var insts_357 vm.Value
	var term_358 vm.Value
	var entry_QMARK__359 vm.Value
	var header_360 vm.Value
	var body_361 vm.Value
	var or__x_362 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = params_3, preds_5, insts_7, term_9, arg__9621_11, entry_QMARK__12, f_14, bid_15, params_16, preds_17, insts_18, term_19, entry_QMARK__20, arg__9622_21, f_22, bid_23, params_24, preds_25, insts_26, term_27, entry_QMARK__28, arg__9622_29, arg__9623_35, f_36, bid_37, params_38, preds_39, insts_40, term_41, entry_QMARK__42, arg__9622_43, arg__9703_54, arg__9781_64, arg__9782_65, v96, arg__9623_67, f_68, arg__9625_69, bid_70, params_71, preds_72, insts_73, term_74, entry_QMARK__75, arg__9622_76, arg__9624_77, arg__9626_78, arg__9782_79, arg__9783_80, arg__9804_103, arg__9822_109, arg__9823_110, arg__9842_117, arg__9860_123, arg__9861_124, v125, arg__9623_81, f_82, arg__9625_83, bid_84, params_85, preds_86, insts_87, term_88, entry_QMARK__89, arg__9622_90, arg__9624_91, arg__9626_92, arg__9782_93, arg__9783_94, arg__9862_129, arg__9623_130, f_131, arg__9625_132, bid_133, params_134, preds_135, insts_136, term_137, entry_QMARK__138, arg__9622_139, arg__9624_140, arg__9626_141, arg__9782_142, arg__9783_143, f_147, bid_148, params_149, preds_150, insts_151, term_152, entry_QMARK__153, head__9864_154, arg__9865_155, f_156, bid_157, params_158, preds_159, insts_160, term_161, entry_QMARK__162, head__9864_163, arg__9865_164, arg__9866_170, f_171, bid_172, params_173, preds_174, insts_175, term_176, entry_QMARK__177, head__9864_178, arg__9865_179, arg__9946_190, arg__10024_200, arg__10025_201, v234, arg__9866_203, f_204, arg__9868_205, bid_206, params_207, preds_208, insts_209, term_210, entry_QMARK__211, head__9864_212, arg__9865_213, arg__9867_214, arg__9869_215, arg__10025_216, arg__10026_217, arg__10047_241, arg__10065_247, arg__10066_248, arg__10085_255, arg__10103_261, arg__10104_262, v263, arg__9866_218, f_219, arg__9868_220, bid_221, params_222, preds_223, insts_224, term_225, entry_QMARK__226, head__9864_227, arg__9865_228, arg__9867_229, arg__9869_230, arg__10025_231, arg__10026_232, arg__10105_267, arg__9866_268, f_269, arg__9868_270, bid_271, params_272, preds_273, insts_274, term_275, entry_QMARK__276, head__9864_277, arg__9865_278, arg__9867_279, arg__9869_280, arg__10025_281, arg__10026_282, header_284, arg__10123_293, arg__10141_303, body_304, or__x_326, f_305, bid_306, params_307, preds_308, insts_309, term_310, entry_QMARK__311, header_312, body_313, v365, f_314, bid_315, params_316, preds_317, insts_318, term_319, entry_QMARK__320, header_321, body_322, term_line_369, f_370, bid_371, params_372, preds_373, insts_374, term_375, entry_QMARK__376, header_377, body_378, v382, f_327, bid_328, params_329, preds_330, insts_331, term_332, entry_QMARK__333, header_334, body_335, or__x_336, f_337, bid_338, params_339, preds_340, insts_341, term_342, entry_QMARK__343, header_344, body_345, or__x_346, v350, v352, f_353, bid_354, params_355, preds_356, insts_357, term_358, entry_QMARK__359, header_360, body_361, or__x_362
	params_3, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-params").Deref(), []vm.Value{arg1, arg0})
	if callErr != nil {
		return nil, callErr
	}
	preds_5, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-preds").Deref(), []vm.Value{arg1, arg0})
	if callErr != nil {
		return nil, callErr
	}
	insts_7, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-insts").Deref(), []vm.Value{arg1, arg0})
	if callErr != nil {
		return nil, callErr
	}
	term_9, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-term").Deref(), []vm.Value{arg1, arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__9621_11, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-entry").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	entry_QMARK__12 = arg1 == arg__9621_11
	if entry_QMARK__12 {
		f_14 = arg0
		bid_15 = arg1
		params_16 = params_3
		preds_17 = preds_5
		insts_18 = insts_7
		term_19 = term_9
		entry_QMARK__20 = entry_QMARK__12
		arg__9622_21 = "  "
		goto b1
	} else {
		f_22 = arg0
		bid_23 = arg1
		params_24 = params_3
		preds_25 = preds_5
		insts_26 = insts_7
		term_27 = term_9
		entry_QMARK__28 = entry_QMARK__12
		arg__9622_29 = "  "
		goto b2
	}
b1:
	;
	arg__9623_35 = "entry "
	f_36 = f_14
	bid_37 = bid_15
	params_38 = params_16
	preds_39 = preds_17
	insts_40 = insts_18
	term_41 = term_19
	entry_QMARK__42 = vm.Boolean(entry_QMARK__20)
	arg__9622_43 = arg__9622_21
	goto b3
b2:
	;
	arg__9623_35 = ""
	f_36 = f_22
	bid_37 = bid_23
	params_38 = params_24
	preds_39 = preds_25
	insts_40 = insts_26
	term_41 = term_27
	entry_QMARK__42 = vm.Boolean(entry_QMARK__28)
	arg__9622_43 = arg__9622_29
	goto b3
b3:
	;
	arg__9703_54, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var arg__9674_5 vm.Value
		var arg__9681_8 vm.Value
		var arg__9682_9 vm.Value
		var arg__9692_14 vm.Value
		var arg__9699_17 vm.Value
		var arg__9700_18 vm.Value
		var v19 vm.Value
		var callErr error
		_, _, _, _, _, _, _ = arg__9674_5, arg__9681_8, arg__9682_9, arg__9692_14, arg__9699_17, arg__9700_18, v19
		arg__9674_5, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_36})
		if callErr != nil {
			return nil, callErr
		}
		arg__9681_8, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_36})
		if callErr != nil {
			return nil, callErr
		}
		arg__9682_9, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "type-display").Deref(), []vm.Value{arg__9681_8})
		if callErr != nil {
			return nil, callErr
		}
		arg__9692_14, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_36})
		if callErr != nil {
			return nil, callErr
		}
		arg__9699_17, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_36})
		if callErr != nil {
			return nil, callErr
		}
		arg__9700_18, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "type-display").Deref(), []vm.Value{arg__9699_17})
		if callErr != nil {
			return nil, callErr
		}
		v19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("v"), arg0, vm.String(": "), arg__9700_18})
		if callErr != nil {
			return nil, callErr
		}
		return v19, nil
	}), params_38})
	if callErr != nil {
		return nil, callErr
	}
	arg__9781_64, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var arg__9752_5 vm.Value
		var arg__9759_8 vm.Value
		var arg__9760_9 vm.Value
		var arg__9770_14 vm.Value
		var arg__9777_17 vm.Value
		var arg__9778_18 vm.Value
		var v19 vm.Value
		var callErr error
		_, _, _, _, _, _, _ = arg__9752_5, arg__9759_8, arg__9760_9, arg__9770_14, arg__9777_17, arg__9778_18, v19
		arg__9752_5, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_36})
		if callErr != nil {
			return nil, callErr
		}
		arg__9759_8, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_36})
		if callErr != nil {
			return nil, callErr
		}
		arg__9760_9, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "type-display").Deref(), []vm.Value{arg__9759_8})
		if callErr != nil {
			return nil, callErr
		}
		arg__9770_14, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_36})
		if callErr != nil {
			return nil, callErr
		}
		arg__9777_17, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_36})
		if callErr != nil {
			return nil, callErr
		}
		arg__9778_18, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "type-display").Deref(), []vm.Value{arg__9777_17})
		if callErr != nil {
			return nil, callErr
		}
		v19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("v"), arg0, vm.String(": "), arg__9778_18})
		if callErr != nil {
			return nil, callErr
		}
		return v19, nil
	}), params_38})
	if callErr != nil {
		return nil, callErr
	}
	arg__9782_65, callErr = rt.InvokeValue(rt.LookupVar("clojure.string", "join").Deref(), []vm.Value{vm.String(", "), arg__9781_64})
	if callErr != nil {
		return nil, callErr
	}
	v96, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{preds_39})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v96) {
		arg__9623_67 = arg__9623_35
		f_68 = f_36
		arg__9625_69 = bid_37
		bid_70 = bid_37
		params_71 = params_38
		preds_72 = preds_39
		insts_73 = insts_40
		term_74 = term_41
		entry_QMARK__75 = entry_QMARK__42
		arg__9622_76 = arg__9622_43
		arg__9624_77 = "b"
		arg__9626_78 = "("
		arg__9782_79 = arg__9782_65
		arg__9783_80 = "):"
		goto b4
	} else {
		arg__9623_81 = arg__9623_35
		f_82 = f_36
		arg__9625_83 = bid_37
		bid_84 = bid_37
		params_85 = params_38
		preds_86 = preds_39
		insts_87 = insts_40
		term_88 = term_41
		entry_QMARK__89 = entry_QMARK__42
		arg__9622_90 = arg__9622_43
		arg__9624_91 = "b"
		arg__9626_92 = "("
		arg__9782_93 = arg__9782_65
		arg__9783_94 = "):"
		goto b5
	}
b4:
	;
	arg__9804_103, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v4 vm.Value
		var callErr error
		_ = v4
		v4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("b"), arg0})
		if callErr != nil {
			return nil, callErr
		}
		return v4, nil
	}), preds_72})
	if callErr != nil {
		return nil, callErr
	}
	arg__9822_109, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v4 vm.Value
		var callErr error
		_ = v4
		v4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("b"), arg0})
		if callErr != nil {
			return nil, callErr
		}
		return v4, nil
	}), preds_72})
	if callErr != nil {
		return nil, callErr
	}
	arg__9823_110, callErr = rt.InvokeValue(rt.LookupVar("clojure.string", "join").Deref(), []vm.Value{vm.String(", "), arg__9822_109})
	if callErr != nil {
		return nil, callErr
	}
	arg__9842_117, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v4 vm.Value
		var callErr error
		_ = v4
		v4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("b"), arg0})
		if callErr != nil {
			return nil, callErr
		}
		return v4, nil
	}), preds_72})
	if callErr != nil {
		return nil, callErr
	}
	arg__9860_123, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v4 vm.Value
		var callErr error
		_ = v4
		v4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("b"), arg0})
		if callErr != nil {
			return nil, callErr
		}
		return v4, nil
	}), preds_72})
	if callErr != nil {
		return nil, callErr
	}
	arg__9861_124, callErr = rt.InvokeValue(rt.LookupVar("clojure.string", "join").Deref(), []vm.Value{vm.String(", "), arg__9860_123})
	if callErr != nil {
		return nil, callErr
	}
	v125, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("    ; preds: "), arg__9861_124})
	if callErr != nil {
		return nil, callErr
	}
	arg__9862_129 = v125
	arg__9623_130 = arg__9623_67
	f_131 = f_68
	arg__9625_132 = arg__9625_69
	bid_133 = bid_70
	params_134 = params_71
	preds_135 = preds_72
	insts_136 = insts_73
	term_137 = term_74
	entry_QMARK__138 = entry_QMARK__75
	arg__9622_139 = arg__9622_76
	arg__9624_140 = arg__9624_77
	arg__9626_141 = arg__9626_78
	arg__9782_142 = arg__9782_79
	arg__9783_143 = arg__9783_80
	goto b6
b5:
	;
	arg__9862_129 = vm.String("")
	arg__9623_130 = arg__9623_81
	f_131 = f_82
	arg__9625_132 = arg__9625_83
	bid_133 = bid_84
	params_134 = params_85
	preds_135 = preds_86
	insts_136 = insts_87
	term_137 = term_88
	entry_QMARK__138 = entry_QMARK__89
	arg__9622_139 = arg__9622_90
	arg__9624_140 = arg__9624_91
	arg__9626_141 = arg__9626_92
	arg__9782_142 = arg__9782_93
	arg__9783_143 = arg__9783_94
	goto b6
b6:
	;
	if vm.IsTruthy(entry_QMARK__138) {
		f_147 = f_131
		bid_148 = bid_133
		params_149 = params_134
		preds_150 = preds_135
		insts_151 = insts_136
		term_152 = term_137
		entry_QMARK__153 = vm.IsTruthy(entry_QMARK__138)
		head__9864_154 = rt.LookupVar("clojure.core", "str").Deref()
		arg__9865_155 = "  "
		goto b7
	} else {
		f_156 = f_131
		bid_157 = bid_133
		params_158 = params_134
		preds_159 = preds_135
		insts_160 = insts_136
		term_161 = term_137
		entry_QMARK__162 = vm.IsTruthy(entry_QMARK__138)
		head__9864_163 = rt.LookupVar("clojure.core", "str").Deref()
		arg__9865_164 = "  "
		goto b8
	}
b7:
	;
	arg__9866_170 = "entry "
	f_171 = f_147
	bid_172 = bid_148
	params_173 = params_149
	preds_174 = preds_150
	insts_175 = insts_151
	term_176 = term_152
	entry_QMARK__177 = vm.Boolean(entry_QMARK__153)
	head__9864_178 = head__9864_154
	arg__9865_179 = arg__9865_155
	goto b9
b8:
	;
	arg__9866_170 = ""
	f_171 = f_156
	bid_172 = bid_157
	params_173 = params_158
	preds_174 = preds_159
	insts_175 = insts_160
	term_176 = term_161
	entry_QMARK__177 = vm.Boolean(entry_QMARK__162)
	head__9864_178 = head__9864_163
	arg__9865_179 = arg__9865_164
	goto b9
b9:
	;
	arg__9946_190, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var arg__9917_5 vm.Value
		var arg__9924_8 vm.Value
		var arg__9925_9 vm.Value
		var arg__9935_14 vm.Value
		var arg__9942_17 vm.Value
		var arg__9943_18 vm.Value
		var v19 vm.Value
		var callErr error
		_, _, _, _, _, _, _ = arg__9917_5, arg__9924_8, arg__9925_9, arg__9935_14, arg__9942_17, arg__9943_18, v19
		arg__9917_5, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_171})
		if callErr != nil {
			return nil, callErr
		}
		arg__9924_8, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_171})
		if callErr != nil {
			return nil, callErr
		}
		arg__9925_9, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "type-display").Deref(), []vm.Value{arg__9924_8})
		if callErr != nil {
			return nil, callErr
		}
		arg__9935_14, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_171})
		if callErr != nil {
			return nil, callErr
		}
		arg__9942_17, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_171})
		if callErr != nil {
			return nil, callErr
		}
		arg__9943_18, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "type-display").Deref(), []vm.Value{arg__9942_17})
		if callErr != nil {
			return nil, callErr
		}
		v19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("v"), arg0, vm.String(": "), arg__9943_18})
		if callErr != nil {
			return nil, callErr
		}
		return v19, nil
	}), params_173})
	if callErr != nil {
		return nil, callErr
	}
	arg__10024_200, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var arg__9995_5 vm.Value
		var arg__10002_8 vm.Value
		var arg__10003_9 vm.Value
		var arg__10013_14 vm.Value
		var arg__10020_17 vm.Value
		var arg__10021_18 vm.Value
		var v19 vm.Value
		var callErr error
		_, _, _, _, _, _, _ = arg__9995_5, arg__10002_8, arg__10003_9, arg__10013_14, arg__10020_17, arg__10021_18, v19
		arg__9995_5, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_171})
		if callErr != nil {
			return nil, callErr
		}
		arg__10002_8, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_171})
		if callErr != nil {
			return nil, callErr
		}
		arg__10003_9, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "type-display").Deref(), []vm.Value{arg__10002_8})
		if callErr != nil {
			return nil, callErr
		}
		arg__10013_14, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_171})
		if callErr != nil {
			return nil, callErr
		}
		arg__10020_17, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_171})
		if callErr != nil {
			return nil, callErr
		}
		arg__10021_18, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "type-display").Deref(), []vm.Value{arg__10020_17})
		if callErr != nil {
			return nil, callErr
		}
		v19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("v"), arg0, vm.String(": "), arg__10021_18})
		if callErr != nil {
			return nil, callErr
		}
		return v19, nil
	}), params_173})
	if callErr != nil {
		return nil, callErr
	}
	arg__10025_201, callErr = rt.InvokeValue(rt.LookupVar("clojure.string", "join").Deref(), []vm.Value{vm.String(", "), arg__10024_200})
	if callErr != nil {
		return nil, callErr
	}
	v234, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{preds_174})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v234) {
		arg__9866_203 = arg__9866_170
		f_204 = f_171
		arg__9868_205 = bid_172
		bid_206 = bid_172
		params_207 = params_173
		preds_208 = preds_174
		insts_209 = insts_175
		term_210 = term_176
		entry_QMARK__211 = entry_QMARK__177
		head__9864_212 = head__9864_178
		arg__9865_213 = arg__9865_179
		arg__9867_214 = "b"
		arg__9869_215 = "("
		arg__10025_216 = arg__10025_201
		arg__10026_217 = "):"
		goto b10
	} else {
		arg__9866_218 = arg__9866_170
		f_219 = f_171
		arg__9868_220 = bid_172
		bid_221 = bid_172
		params_222 = params_173
		preds_223 = preds_174
		insts_224 = insts_175
		term_225 = term_176
		entry_QMARK__226 = entry_QMARK__177
		head__9864_227 = head__9864_178
		arg__9865_228 = arg__9865_179
		arg__9867_229 = "b"
		arg__9869_230 = "("
		arg__10025_231 = arg__10025_201
		arg__10026_232 = "):"
		goto b11
	}
b10:
	;
	arg__10047_241, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v4 vm.Value
		var callErr error
		_ = v4
		v4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("b"), arg0})
		if callErr != nil {
			return nil, callErr
		}
		return v4, nil
	}), preds_208})
	if callErr != nil {
		return nil, callErr
	}
	arg__10065_247, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v4 vm.Value
		var callErr error
		_ = v4
		v4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("b"), arg0})
		if callErr != nil {
			return nil, callErr
		}
		return v4, nil
	}), preds_208})
	if callErr != nil {
		return nil, callErr
	}
	arg__10066_248, callErr = rt.InvokeValue(rt.LookupVar("clojure.string", "join").Deref(), []vm.Value{vm.String(", "), arg__10065_247})
	if callErr != nil {
		return nil, callErr
	}
	arg__10085_255, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v4 vm.Value
		var callErr error
		_ = v4
		v4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("b"), arg0})
		if callErr != nil {
			return nil, callErr
		}
		return v4, nil
	}), preds_208})
	if callErr != nil {
		return nil, callErr
	}
	arg__10103_261, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v4 vm.Value
		var callErr error
		_ = v4
		v4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("b"), arg0})
		if callErr != nil {
			return nil, callErr
		}
		return v4, nil
	}), preds_208})
	if callErr != nil {
		return nil, callErr
	}
	arg__10104_262, callErr = rt.InvokeValue(rt.LookupVar("clojure.string", "join").Deref(), []vm.Value{vm.String(", "), arg__10103_261})
	if callErr != nil {
		return nil, callErr
	}
	v263, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("    ; preds: "), arg__10104_262})
	if callErr != nil {
		return nil, callErr
	}
	arg__10105_267 = v263
	arg__9866_268 = arg__9866_203
	f_269 = f_204
	arg__9868_270 = arg__9868_205
	bid_271 = bid_206
	params_272 = params_207
	preds_273 = preds_208
	insts_274 = insts_209
	term_275 = term_210
	entry_QMARK__276 = entry_QMARK__211
	head__9864_277 = head__9864_212
	arg__9865_278 = arg__9865_213
	arg__9867_279 = arg__9867_214
	arg__9869_280 = arg__9869_215
	arg__10025_281 = arg__10025_216
	arg__10026_282 = arg__10026_217
	goto b12
b11:
	;
	arg__10105_267 = vm.String("")
	arg__9866_268 = arg__9866_218
	f_269 = f_219
	arg__9868_270 = arg__9868_220
	bid_271 = bid_221
	params_272 = params_222
	preds_273 = preds_223
	insts_274 = insts_224
	term_275 = term_225
	entry_QMARK__276 = entry_QMARK__226
	head__9864_277 = head__9864_227
	arg__9865_278 = arg__9865_228
	arg__9867_279 = arg__9867_229
	arg__9869_280 = arg__9869_230
	arg__10025_281 = arg__10025_231
	arg__10026_282 = arg__10026_232
	goto b12
b12:
	;
	header_284, callErr = rt.InvokeValue(head__9864_277, []vm.Value{vm.String(arg__9865_278), vm.String(arg__9866_268), vm.String(arg__9867_279), arg__9868_270, vm.String(arg__9869_280), arg__10025_281, vm.String(arg__10026_282), arg__10105_267, vm.String("\n")})
	if callErr != nil {
		return nil, callErr
	}
	arg__10123_293, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "write-node").Deref(), []vm.Value{f_269, arg0})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), insts_274})
	if callErr != nil {
		return nil, callErr
	}
	arg__10141_303, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "write-node").Deref(), []vm.Value{f_269, arg0})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), insts_274})
	if callErr != nil {
		return nil, callErr
	}
	body_304, callErr = rt.InvokeValue(rt.LookupVar("clojure.string", "join").Deref(), []vm.Value{vm.String(""), arg__10141_303})
	if callErr != nil {
		return nil, callErr
	}
	or__x_326, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not=").Deref(), []vm.Value{term_275, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(or__x_326) {
		f_327 = f_269
		bid_328 = bid_271
		params_329 = params_272
		preds_330 = preds_273
		insts_331 = insts_274
		term_332 = term_275
		entry_QMARK__333 = entry_QMARK__276
		header_334 = header_284
		body_335 = body_304
		or__x_336 = or__x_326
		goto b16
	} else {
		f_337 = f_269
		bid_338 = bid_271
		params_339 = params_272
		preds_340 = preds_273
		insts_341 = insts_274
		term_342 = term_275
		entry_QMARK__343 = entry_QMARK__276
		header_344 = header_284
		body_345 = body_304
		or__x_346 = or__x_326
		goto b17
	}
b13:
	;
	v365, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "write-node").Deref(), []vm.Value{f_305, term_310})
	if callErr != nil {
		return nil, callErr
	}
	term_line_369 = v365
	f_370 = f_305
	bid_371 = bid_306
	params_372 = params_307
	preds_373 = preds_308
	insts_374 = insts_309
	term_375 = term_310
	entry_QMARK__376 = entry_QMARK__311
	header_377 = header_312
	body_378 = body_313
	goto b15
b14:
	;
	term_line_369 = vm.String("")
	f_370 = f_314
	bid_371 = bid_315
	params_372 = params_316
	preds_373 = preds_317
	insts_374 = insts_318
	term_375 = term_319
	entry_QMARK__376 = entry_QMARK__320
	header_377 = header_321
	body_378 = body_322
	goto b15
b15:
	;
	v382, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{header_377, body_378, term_line_369, vm.String("\n")})
	if callErr != nil {
		return nil, callErr
	}
	return v382, nil
b16:
	;
	v352 = or__x_336
	f_353 = f_327
	bid_354 = bid_328
	params_355 = params_329
	preds_356 = preds_330
	insts_357 = insts_331
	term_358 = term_332
	entry_QMARK__359 = entry_QMARK__333
	header_360 = header_334
	body_361 = body_335
	or__x_362 = or__x_336
	goto b18
b17:
	;
	v350, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{insts_341})
	if callErr != nil {
		return nil, callErr
	}
	v352 = v350
	f_353 = f_337
	bid_354 = bid_338
	params_355 = params_339
	preds_356 = preds_340
	insts_357 = insts_341
	term_358 = term_342
	entry_QMARK__359 = entry_QMARK__343
	header_360 = header_344
	body_361 = body_345
	or__x_362 = or__x_346
	goto b18
b18:
	;
	if vm.IsTruthy(v352) {
		f_305 = f_353
		bid_306 = bid_354
		params_307 = params_355
		preds_308 = preds_356
		insts_309 = insts_357
		term_310 = term_358
		entry_QMARK__311 = entry_QMARK__359
		header_312 = header_360
		body_313 = body_361
		goto b13
	} else {
		f_314 = f_353
		bid_315 = bid_354
		params_316 = params_355
		preds_317 = preds_356
		insts_318 = insts_357
		term_319 = term_358
		entry_QMARK__320 = entry_QMARK__359
		header_321 = header_360
		body_322 = body_361
		goto b14
	}
}
func dump(arg0 vm.Value) (vm.Value, error) {
	var arg__10168_4 vm.Value
	var arg__10173_7 vm.Value
	var arg__10178_10 vm.Value
	var arg__10190_17 vm.Value
	var arg__10201_23 vm.Value
	var arg__10202_24 vm.Value
	var arg__10214_31 vm.Value
	var arg__10225_37 vm.Value
	var arg__10226_38 vm.Value
	var arg__10227_39 vm.Value
	var arg__10233_43 vm.Value
	var arg__10238_46 vm.Value
	var arg__10243_49 vm.Value
	var arg__10255_56 vm.Value
	var arg__10266_62 vm.Value
	var arg__10267_63 vm.Value
	var arg__10279_70 vm.Value
	var arg__10290_76 vm.Value
	var arg__10291_77 vm.Value
	var arg__10292_78 vm.Value
	var v79 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__10168_4, arg__10173_7, arg__10178_10, arg__10190_17, arg__10201_23, arg__10202_24, arg__10214_31, arg__10225_37, arg__10226_38, arg__10227_39, arg__10233_43, arg__10238_46, arg__10243_49, arg__10255_56, arg__10266_62, arg__10267_63, arg__10279_70, arg__10290_76, arg__10291_77, arg__10292_78, v79
	arg__10168_4, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-name").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__10173_7, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-arity").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__10178_10, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-variadic?").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__10190_17, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__10201_23, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__10202_24, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "write-block").Deref(), []vm.Value{arg0, arg0})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), arg__10201_23})
	if callErr != nil {
		return nil, callErr
	}
	arg__10214_31, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__10225_37, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__10226_38, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "write-block").Deref(), []vm.Value{arg0, arg0})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), arg__10225_37})
	if callErr != nil {
		return nil, callErr
	}
	arg__10227_39, callErr = rt.InvokeValue(rt.LookupVar("clojure.string", "join").Deref(), []vm.Value{vm.String(""), arg__10226_38})
	if callErr != nil {
		return nil, callErr
	}
	arg__10233_43, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-name").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__10238_46, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-arity").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__10243_49, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-variadic?").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__10255_56, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__10266_62, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__10267_63, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "write-block").Deref(), []vm.Value{arg0, arg0})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), arg__10266_62})
	if callErr != nil {
		return nil, callErr
	}
	arg__10279_70, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__10290_76, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__10291_77, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "write-block").Deref(), []vm.Value{arg0, arg0})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), arg__10290_76})
	if callErr != nil {
		return nil, callErr
	}
	arg__10292_78, callErr = rt.InvokeValue(rt.LookupVar("clojure.string", "join").Deref(), []vm.Value{vm.String(""), arg__10291_77})
	if callErr != nil {
		return nil, callErr
	}
	v79, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("fn "), arg__10233_43, vm.String("(arity="), arg__10238_46, vm.String(", variadic="), arg__10243_49, vm.String("):\n"), arg__10292_78})
	if callErr != nil {
		return nil, callErr
	}
	return v79, nil
}
