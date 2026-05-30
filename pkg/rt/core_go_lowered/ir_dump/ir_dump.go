package ir_dump

import (
	rt "github.com/nooga/let-go/pkg/rt"
	vm "github.com/nooga/let-go/pkg/vm"
)

func format_args(arg0 vm.Value) (vm.Value, error) {
	var arg__8902_5 vm.Value
	var arg__8920_11 vm.Value
	var v12 vm.Value
	var callErr error
	_, _, _ = arg__8902_5, arg__8920_11, v12
	arg__8902_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
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
	arg__8920_11, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
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
	v12, callErr = rt.InvokeValue(rt.LookupVar("clojure.string", "join").Deref(), []vm.Value{vm.String(", "), arg__8920_11})
	if callErr != nil {
		return nil, callErr
	}
	return v12, nil
}
func scalar_type_display(arg0 vm.Value) vm.Value {
	var v6 bool
	var case__8921_1 vm.Value
	var t_2 vm.Value
	var case__8921_3 vm.Value
	var t_4 vm.Value
	var v15 bool
	var v154 vm.Value
	var case__8921_155 vm.Value
	var t_156 vm.Value
	var case__8921_10 vm.Value
	var t_11 vm.Value
	var case__8921_12 vm.Value
	var t_13 vm.Value
	var v24 bool
	var v150 vm.Value
	var case__8921_151 vm.Value
	var t_152 vm.Value
	var case__8921_19 vm.Value
	var t_20 vm.Value
	var case__8921_21 vm.Value
	var t_22 vm.Value
	var v33 bool
	var v146 vm.Value
	var case__8921_147 vm.Value
	var t_148 vm.Value
	var case__8921_28 vm.Value
	var t_29 vm.Value
	var case__8921_30 vm.Value
	var t_31 vm.Value
	var v42 bool
	var v142 vm.Value
	var case__8921_143 vm.Value
	var t_144 vm.Value
	var case__8921_37 vm.Value
	var t_38 vm.Value
	var case__8921_39 vm.Value
	var t_40 vm.Value
	var v51 bool
	var v138 vm.Value
	var case__8921_139 vm.Value
	var t_140 vm.Value
	var case__8921_46 vm.Value
	var t_47 vm.Value
	var case__8921_48 vm.Value
	var t_49 vm.Value
	var v60 bool
	var v134 vm.Value
	var case__8921_135 vm.Value
	var t_136 vm.Value
	var case__8921_55 vm.Value
	var t_56 vm.Value
	var case__8921_57 vm.Value
	var t_58 vm.Value
	var v69 bool
	var v130 vm.Value
	var case__8921_131 vm.Value
	var t_132 vm.Value
	var case__8921_64 vm.Value
	var t_65 vm.Value
	var case__8921_66 vm.Value
	var t_67 vm.Value
	var v78 bool
	var v126 vm.Value
	var case__8921_127 vm.Value
	var t_128 vm.Value
	var case__8921_73 vm.Value
	var t_74 vm.Value
	var case__8921_75 vm.Value
	var t_76 vm.Value
	var v87 bool
	var v122 vm.Value
	var case__8921_123 vm.Value
	var t_124 vm.Value
	var case__8921_82 vm.Value
	var t_83 vm.Value
	var case__8921_84 vm.Value
	var t_85 vm.Value
	var v96 bool
	var v118 vm.Value
	var case__8921_119 vm.Value
	var t_120 vm.Value
	var case__8921_91 vm.Value
	var t_92 vm.Value
	var case__8921_93 vm.Value
	var t_94 vm.Value
	var v114 vm.Value
	var case__8921_115 vm.Value
	var t_116 vm.Value
	var case__8921_100 vm.Value
	var t_101 vm.Value
	var case__8921_102 vm.Value
	var t_103 vm.Value
	var v110 vm.Value
	var case__8921_111 vm.Value
	var t_112 vm.Value
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = v6, case__8921_1, t_2, case__8921_3, t_4, v15, v154, case__8921_155, t_156, case__8921_10, t_11, case__8921_12, t_13, v24, v150, case__8921_151, t_152, case__8921_19, t_20, case__8921_21, t_22, v33, v146, case__8921_147, t_148, case__8921_28, t_29, case__8921_30, t_31, v42, v142, case__8921_143, t_144, case__8921_37, t_38, case__8921_39, t_40, v51, v138, case__8921_139, t_140, case__8921_46, t_47, case__8921_48, t_49, v60, v134, case__8921_135, t_136, case__8921_55, t_56, case__8921_57, t_58, v69, v130, case__8921_131, t_132, case__8921_64, t_65, case__8921_66, t_67, v78, v126, case__8921_127, t_128, case__8921_73, t_74, case__8921_75, t_76, v87, v122, case__8921_123, t_124, case__8921_82, t_83, case__8921_84, t_85, v96, v118, case__8921_119, t_120, case__8921_91, t_92, case__8921_93, t_94, v114, case__8921_115, t_116, case__8921_100, t_101, case__8921_102, t_103, v110, case__8921_111, t_112
	v6 = arg0 == vm.Keyword("unknown")
	if v6 {
		case__8921_1 = arg0
		t_2 = arg0
		goto b1
	} else {
		case__8921_3 = arg0
		t_4 = arg0
		goto b2
	}
b1:
	;
	v154 = vm.String("unknown")
	case__8921_155 = case__8921_1
	t_156 = t_2
	goto b3
b2:
	;
	v15 = case__8921_3 == vm.Keyword("bottom")
	if v15 {
		case__8921_10 = case__8921_3
		t_11 = t_4
		goto b4
	} else {
		case__8921_12 = case__8921_3
		t_13 = t_4
		goto b5
	}
b3:
	;
	return v154
b4:
	;
	v150 = vm.String("bottom")
	case__8921_151 = case__8921_10
	t_152 = t_11
	goto b6
b5:
	;
	v24 = case__8921_12 == vm.Keyword("true")
	if v24 {
		case__8921_19 = case__8921_12
		t_20 = t_13
		goto b7
	} else {
		case__8921_21 = case__8921_12
		t_22 = t_13
		goto b8
	}
b6:
	;
	v154 = v150
	case__8921_155 = case__8921_151
	t_156 = t_152
	goto b3
b7:
	;
	v146 = vm.String("true")
	case__8921_147 = case__8921_19
	t_148 = t_20
	goto b9
b8:
	;
	v33 = case__8921_21 == vm.Keyword("false")
	if v33 {
		case__8921_28 = case__8921_21
		t_29 = t_22
		goto b10
	} else {
		case__8921_30 = case__8921_21
		t_31 = t_22
		goto b11
	}
b9:
	;
	v150 = v146
	case__8921_151 = case__8921_147
	t_152 = t_148
	goto b6
b10:
	;
	v142 = vm.String("false")
	case__8921_143 = case__8921_28
	t_144 = t_29
	goto b12
b11:
	;
	v42 = case__8921_30 == vm.Keyword("int")
	if v42 {
		case__8921_37 = case__8921_30
		t_38 = t_31
		goto b13
	} else {
		case__8921_39 = case__8921_30
		t_40 = t_31
		goto b14
	}
b12:
	;
	v146 = v142
	case__8921_147 = case__8921_143
	t_148 = t_144
	goto b9
b13:
	;
	v138 = vm.String("int")
	case__8921_139 = case__8921_37
	t_140 = t_38
	goto b15
b14:
	;
	v51 = case__8921_39 == vm.Keyword("float")
	if v51 {
		case__8921_46 = case__8921_39
		t_47 = t_40
		goto b16
	} else {
		case__8921_48 = case__8921_39
		t_49 = t_40
		goto b17
	}
b15:
	;
	v142 = v138
	case__8921_143 = case__8921_139
	t_144 = t_140
	goto b12
b16:
	;
	v134 = vm.String("float")
	case__8921_135 = case__8921_46
	t_136 = t_47
	goto b18
b17:
	;
	v60 = case__8921_48 == vm.Keyword("number")
	if v60 {
		case__8921_55 = case__8921_48
		t_56 = t_49
		goto b19
	} else {
		case__8921_57 = case__8921_48
		t_58 = t_49
		goto b20
	}
b18:
	;
	v138 = v134
	case__8921_139 = case__8921_135
	t_140 = t_136
	goto b15
b19:
	;
	v130 = vm.String("number")
	case__8921_131 = case__8921_55
	t_132 = t_56
	goto b21
b20:
	;
	v69 = case__8921_57 == vm.Keyword("bool")
	if v69 {
		case__8921_64 = case__8921_57
		t_65 = t_58
		goto b22
	} else {
		case__8921_66 = case__8921_57
		t_67 = t_58
		goto b23
	}
b21:
	;
	v134 = v130
	case__8921_135 = case__8921_131
	t_136 = t_132
	goto b18
b22:
	;
	v126 = vm.String("bool")
	case__8921_127 = case__8921_64
	t_128 = t_65
	goto b24
b23:
	;
	v78 = case__8921_66 == vm.Keyword("nil")
	if v78 {
		case__8921_73 = case__8921_66
		t_74 = t_67
		goto b25
	} else {
		case__8921_75 = case__8921_66
		t_76 = t_67
		goto b26
	}
b24:
	;
	v130 = v126
	case__8921_131 = case__8921_127
	t_132 = t_128
	goto b21
b25:
	;
	v122 = vm.String("nil")
	case__8921_123 = case__8921_73
	t_124 = t_74
	goto b27
b26:
	;
	v87 = case__8921_75 == vm.Keyword("string")
	if v87 {
		case__8921_82 = case__8921_75
		t_83 = t_76
		goto b28
	} else {
		case__8921_84 = case__8921_75
		t_85 = t_76
		goto b29
	}
b27:
	;
	v126 = v122
	case__8921_127 = case__8921_123
	t_128 = t_124
	goto b24
b28:
	;
	v118 = vm.String("string")
	case__8921_119 = case__8921_82
	t_120 = t_83
	goto b30
b29:
	;
	v96 = case__8921_84 == vm.Keyword("any")
	if v96 {
		case__8921_91 = case__8921_84
		t_92 = t_85
		goto b31
	} else {
		case__8921_93 = case__8921_84
		t_94 = t_85
		goto b32
	}
b30:
	;
	v122 = v118
	case__8921_123 = case__8921_119
	t_124 = t_120
	goto b27
b31:
	;
	v114 = vm.String("any")
	case__8921_115 = case__8921_91
	t_116 = t_92
	goto b33
b32:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		case__8921_100 = case__8921_93
		t_101 = t_94
		goto b34
	} else {
		case__8921_102 = case__8921_93
		t_103 = t_94
		goto b35
	}
b33:
	;
	v118 = v114
	case__8921_119 = case__8921_115
	t_120 = t_116
	goto b30
b34:
	;
	v110 = vm.String("??")
	case__8921_111 = case__8921_100
	t_112 = t_101
	goto b36
b35:
	;
	v110 = vm.NIL
	case__8921_111 = case__8921_102
	t_112 = t_103
	goto b36
b36:
	;
	v114 = v110
	case__8921_115 = case__8921_111
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
	var arg__8973_34 vm.Value
	var arg__8986_38 vm.Value
	var arg__8987_39 vm.Value
	var arg__9001_44 vm.Value
	var arg__9014_48 vm.Value
	var arg__9015_49 vm.Value
	var arg__9016_50 vm.Value
	var arg__9031_56 vm.Value
	var arg__9044_60 vm.Value
	var arg__9045_61 vm.Value
	var arg__9059_66 vm.Value
	var arg__9072_70 vm.Value
	var arg__9073_71 vm.Value
	var arg__9074_72 vm.Value
	var arg__9075_73 vm.Value
	var arg__9092_81 vm.Value
	var arg__9105_85 vm.Value
	var arg__9106_86 vm.Value
	var arg__9120_91 vm.Value
	var arg__9133_95 vm.Value
	var arg__9134_96 vm.Value
	var arg__9135_97 vm.Value
	var arg__9150_103 vm.Value
	var arg__9163_107 vm.Value
	var arg__9164_108 vm.Value
	var arg__9178_113 vm.Value
	var arg__9191_117 vm.Value
	var arg__9192_118 vm.Value
	var arg__9193_119 vm.Value
	var arg__9194_120 vm.Value
	var v122 vm.Value
	var t_11 vm.Value
	var and__x_127 vm.Value
	var v233 vm.Value
	var t_234 vm.Value
	var t_14 vm.Value
	var and__x_15 vm.Value
	var arg__8958_21 vm.Value
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
	var arg__9203_135 vm.Value
	var v136 bool
	var t_130 vm.Value
	var and__x_131 vm.Value
	var v139 vm.Value
	var t_140 vm.Value
	var and__x_141 vm.Value
	var t_151 vm.Value
	var case__8944_152 vm.Value
	var tag_153 vm.Value
	var v_154 vm.Value
	var v167 vm.Value
	var t_155 vm.Value
	var case__8944_156 vm.Value
	var tag_157 vm.Value
	var v_158 vm.Value
	var v178 bool
	var v213 vm.Value
	var t_214 vm.Value
	var case__8944_215 vm.Value
	var tag_216 vm.Value
	var v_217 vm.Value
	var t_169 vm.Value
	var case__8944_170 vm.Value
	var tag_171 vm.Value
	var v_172 vm.Value
	var v185 vm.Value
	var t_173 vm.Value
	var case__8944_174 vm.Value
	var tag_175 vm.Value
	var v_176 vm.Value
	var v207 vm.Value
	var t_208 vm.Value
	var case__8944_209 vm.Value
	var tag_210 vm.Value
	var v_211 vm.Value
	var t_187 vm.Value
	var case__8944_188 vm.Value
	var tag_189 vm.Value
	var v_190 vm.Value
	var t_191 vm.Value
	var case__8944_192 vm.Value
	var tag_193 vm.Value
	var v_194 vm.Value
	var v201 vm.Value
	var t_202 vm.Value
	var case__8944_203 vm.Value
	var tag_204 vm.Value
	var v_205 vm.Value
	var t_219 vm.Value
	var t_220 vm.Value
	var v227 vm.Value
	var t_228 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = v5, t_2, v8, t_3, and__x_13, v236, t_237, t_10, arg__8973_34, arg__8986_38, arg__8987_39, arg__9001_44, arg__9014_48, arg__9015_49, arg__9016_50, arg__9031_56, arg__9044_60, arg__9045_61, arg__9059_66, arg__9072_70, arg__9073_71, arg__9074_72, arg__9075_73, arg__9092_81, arg__9105_85, arg__9106_86, arg__9120_91, arg__9133_95, arg__9134_96, arg__9135_97, arg__9150_103, arg__9163_107, arg__9164_108, arg__9178_113, arg__9191_117, arg__9192_118, arg__9193_119, arg__9194_120, v122, t_11, and__x_127, v233, t_234, t_14, and__x_15, arg__8958_21, v22, t_16, and__x_17, v25, t_26, and__x_27, t_124, tag_146, v_150, v160, t_125, v230, t_231, t_128, and__x_129, arg__9203_135, v136, t_130, and__x_131, v139, t_140, and__x_141, t_151, case__8944_152, tag_153, v_154, v167, t_155, case__8944_156, tag_157, v_158, v178, v213, t_214, case__8944_215, tag_216, v_217, t_169, case__8944_170, tag_171, v_172, v185, t_173, case__8944_174, tag_175, v_176, v207, t_208, case__8944_209, tag_210, v_211, t_187, case__8944_188, tag_189, v_190, t_191, case__8944_192, tag_193, v_194, v201, t_202, case__8944_203, tag_204, v_205, t_219, t_220, v227, t_228
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
	arg__8973_34, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__8986_38, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__8987_39, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "sort-by").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v6 vm.Value
		var callErr error
		_ = v6
		v6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{rt.LookupVar("ir.dump", "type-order").Deref(), arg0, vm.Int(100)})
		if callErr != nil {
			return nil, callErr
		}
		return v6, nil
	}), arg__8986_38})
	if callErr != nil {
		return nil, callErr
	}
	arg__9001_44, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__9014_48, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__9015_49, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "sort-by").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v6 vm.Value
		var callErr error
		_ = v6
		v6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{rt.LookupVar("ir.dump", "type-order").Deref(), arg0, vm.Int(100)})
		if callErr != nil {
			return nil, callErr
		}
		return v6, nil
	}), arg__9014_48})
	if callErr != nil {
		return nil, callErr
	}
	arg__9016_50, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.LookupVar("ir.dump", "type-display").Deref(), arg__9015_49})
	if callErr != nil {
		return nil, callErr
	}
	arg__9031_56, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__9044_60, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__9045_61, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "sort-by").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v6 vm.Value
		var callErr error
		_ = v6
		v6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{rt.LookupVar("ir.dump", "type-order").Deref(), arg0, vm.Int(100)})
		if callErr != nil {
			return nil, callErr
		}
		return v6, nil
	}), arg__9044_60})
	if callErr != nil {
		return nil, callErr
	}
	arg__9059_66, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__9072_70, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__9073_71, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "sort-by").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v6 vm.Value
		var callErr error
		_ = v6
		v6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{rt.LookupVar("ir.dump", "type-order").Deref(), arg0, vm.Int(100)})
		if callErr != nil {
			return nil, callErr
		}
		return v6, nil
	}), arg__9072_70})
	if callErr != nil {
		return nil, callErr
	}
	arg__9074_72, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.LookupVar("ir.dump", "type-display").Deref(), arg__9073_71})
	if callErr != nil {
		return nil, callErr
	}
	arg__9075_73, callErr = rt.InvokeValue(rt.LookupVar("clojure.string", "join").Deref(), []vm.Value{vm.String(","), arg__9074_72})
	if callErr != nil {
		return nil, callErr
	}
	arg__9092_81, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__9105_85, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__9106_86, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "sort-by").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v6 vm.Value
		var callErr error
		_ = v6
		v6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{rt.LookupVar("ir.dump", "type-order").Deref(), arg0, vm.Int(100)})
		if callErr != nil {
			return nil, callErr
		}
		return v6, nil
	}), arg__9105_85})
	if callErr != nil {
		return nil, callErr
	}
	arg__9120_91, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__9133_95, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__9134_96, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "sort-by").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v6 vm.Value
		var callErr error
		_ = v6
		v6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{rt.LookupVar("ir.dump", "type-order").Deref(), arg0, vm.Int(100)})
		if callErr != nil {
			return nil, callErr
		}
		return v6, nil
	}), arg__9133_95})
	if callErr != nil {
		return nil, callErr
	}
	arg__9135_97, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.LookupVar("ir.dump", "type-display").Deref(), arg__9134_96})
	if callErr != nil {
		return nil, callErr
	}
	arg__9150_103, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__9163_107, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__9164_108, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "sort-by").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v6 vm.Value
		var callErr error
		_ = v6
		v6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{rt.LookupVar("ir.dump", "type-order").Deref(), arg0, vm.Int(100)})
		if callErr != nil {
			return nil, callErr
		}
		return v6, nil
	}), arg__9163_107})
	if callErr != nil {
		return nil, callErr
	}
	arg__9178_113, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__9191_117, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__9192_118, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "sort-by").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v6 vm.Value
		var callErr error
		_ = v6
		v6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{rt.LookupVar("ir.dump", "type-order").Deref(), arg0, vm.Int(100)})
		if callErr != nil {
			return nil, callErr
		}
		return v6, nil
	}), arg__9191_117})
	if callErr != nil {
		return nil, callErr
	}
	arg__9193_119, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.LookupVar("ir.dump", "type-display").Deref(), arg__9192_118})
	if callErr != nil {
		return nil, callErr
	}
	arg__9194_120, callErr = rt.InvokeValue(rt.LookupVar("clojure.string", "join").Deref(), []vm.Value{vm.String(","), arg__9193_119})
	if callErr != nil {
		return nil, callErr
	}
	v122, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("union{"), arg__9194_120, vm.String("}")})
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
	arg__8958_21, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{t_14})
	if callErr != nil {
		return nil, callErr
	}
	v22 = arg__8958_21 == vm.Keyword("union")
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
		case__8944_152 = tag_146
		tag_153 = tag_146
		v_154 = v_150
		goto b16
	} else {
		t_155 = t_124
		case__8944_156 = tag_146
		tag_157 = tag_146
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
	arg__9203_135, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{t_128})
	if callErr != nil {
		return nil, callErr
	}
	v136 = arg__9203_135 == vm.Keyword("const")
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
	case__8944_215 = case__8944_152
	tag_216 = tag_153
	v_217 = v_154
	goto b18
b17:
	;
	v178 = case__8944_156 == vm.Keyword("float")
	if v178 {
		t_169 = t_155
		case__8944_170 = case__8944_156
		tag_171 = tag_157
		v_172 = v_158
		goto b19
	} else {
		t_173 = t_155
		case__8944_174 = case__8944_156
		tag_175 = tag_157
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
	case__8944_209 = case__8944_170
	tag_210 = tag_171
	v_211 = v_172
	goto b21
b20:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		t_187 = t_173
		case__8944_188 = case__8944_174
		tag_189 = tag_175
		v_190 = v_176
		goto b22
	} else {
		t_191 = t_173
		case__8944_192 = case__8944_174
		tag_193 = tag_175
		v_194 = v_176
		goto b23
	}
b21:
	;
	v213 = v207
	t_214 = t_208
	case__8944_215 = case__8944_209
	tag_216 = tag_210
	v_217 = v_211
	goto b18
b22:
	;
	v201 = vm.String("??")
	t_202 = t_187
	case__8944_203 = case__8944_188
	tag_204 = tag_189
	v_205 = v_190
	goto b24
b23:
	;
	v201 = vm.NIL
	t_202 = t_191
	case__8944_203 = case__8944_192
	tag_204 = tag_193
	v_205 = v_194
	goto b24
b24:
	;
	v207 = v201
	t_208 = t_202
	case__8944_209 = case__8944_203
	tag_210 = tag_204
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
	var arg__9237_5 vm.Value
	var arg__9241_9 vm.Value
	var arg__9246_12 vm.Value
	var arg__9250_16 vm.Value
	var arg__9251_17 vm.Value
	var arg__9257_21 vm.Value
	var arg__9261_25 vm.Value
	var arg__9266_28 vm.Value
	var arg__9270_32 vm.Value
	var arg__9271_33 vm.Value
	var arg__9272_34 vm.Value
	var arg__9279_39 vm.Value
	var arg__9283_43 vm.Value
	var arg__9288_46 vm.Value
	var arg__9292_50 vm.Value
	var arg__9293_51 vm.Value
	var arg__9299_55 vm.Value
	var arg__9303_59 vm.Value
	var arg__9308_62 vm.Value
	var arg__9312_66 vm.Value
	var arg__9313_67 vm.Value
	var arg__9314_68 vm.Value
	var v69 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__9237_5, arg__9241_9, arg__9246_12, arg__9250_16, arg__9251_17, arg__9257_21, arg__9261_25, arg__9266_28, arg__9270_32, arg__9271_33, arg__9272_34, arg__9279_39, arg__9283_43, arg__9288_46, arg__9292_50, arg__9293_51, arg__9299_55, arg__9303_59, arg__9308_62, arg__9312_66, arg__9313_67, arg__9314_68, v69
	arg__9237_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "name").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__9241_9, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "re-pattern").Deref(), []vm.Value{vm.String("-")})
	if callErr != nil {
		return nil, callErr
	}
	arg__9246_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "name").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__9250_16, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "re-pattern").Deref(), []vm.Value{vm.String("-")})
	if callErr != nil {
		return nil, callErr
	}
	arg__9251_17, callErr = rt.InvokeValue(rt.LookupVar("clojure.string", "split").Deref(), []vm.Value{arg__9246_12, arg__9250_16})
	if callErr != nil {
		return nil, callErr
	}
	arg__9257_21, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "name").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__9261_25, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "re-pattern").Deref(), []vm.Value{vm.String("-")})
	if callErr != nil {
		return nil, callErr
	}
	arg__9266_28, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "name").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__9270_32, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "re-pattern").Deref(), []vm.Value{vm.String("-")})
	if callErr != nil {
		return nil, callErr
	}
	arg__9271_33, callErr = rt.InvokeValue(rt.LookupVar("clojure.string", "split").Deref(), []vm.Value{arg__9266_28, arg__9270_32})
	if callErr != nil {
		return nil, callErr
	}
	arg__9272_34, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.LookupVar("clojure.string", "capitalize").Deref(), arg__9271_33})
	if callErr != nil {
		return nil, callErr
	}
	arg__9279_39, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "name").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__9283_43, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "re-pattern").Deref(), []vm.Value{vm.String("-")})
	if callErr != nil {
		return nil, callErr
	}
	arg__9288_46, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "name").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__9292_50, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "re-pattern").Deref(), []vm.Value{vm.String("-")})
	if callErr != nil {
		return nil, callErr
	}
	arg__9293_51, callErr = rt.InvokeValue(rt.LookupVar("clojure.string", "split").Deref(), []vm.Value{arg__9288_46, arg__9292_50})
	if callErr != nil {
		return nil, callErr
	}
	arg__9299_55, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "name").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__9303_59, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "re-pattern").Deref(), []vm.Value{vm.String("-")})
	if callErr != nil {
		return nil, callErr
	}
	arg__9308_62, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "name").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__9312_66, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "re-pattern").Deref(), []vm.Value{vm.String("-")})
	if callErr != nil {
		return nil, callErr
	}
	arg__9313_67, callErr = rt.InvokeValue(rt.LookupVar("clojure.string", "split").Deref(), []vm.Value{arg__9308_62, arg__9312_66})
	if callErr != nil {
		return nil, callErr
	}
	arg__9314_68, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.LookupVar("clojure.string", "capitalize").Deref(), arg__9313_67})
	if callErr != nil {
		return nil, callErr
	}
	v69, callErr = rt.InvokeValue(rt.LookupVar("clojure.string", "join").Deref(), []vm.Value{vm.String(""), arg__9314_68})
	if callErr != nil {
		return nil, callErr
	}
	return v69, nil
}
func format_refs(arg0 vm.Value) (vm.Value, error) {
	var arg__9331_5 vm.Value
	var arg__9349_11 vm.Value
	var v12 vm.Value
	var callErr error
	_, _, _ = arg__9331_5, arg__9349_11, v12
	arg__9331_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
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
	arg__9349_11, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
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
	v12, callErr = rt.InvokeValue(rt.LookupVar("clojure.string", "join").Deref(), []vm.Value{vm.String(""), arg__9349_11})
	if callErr != nil {
		return nil, callErr
	}
	return v12, nil
}
func format_target(arg0 vm.Value) (vm.Value, error) {
	var arg__9354_3 vm.Value
	var arg__9359_6 vm.Value
	var arg__9364_9 vm.Value
	var arg__9365_10 vm.Value
	var arg__9372_15 vm.Value
	var arg__9377_18 vm.Value
	var arg__9382_21 vm.Value
	var arg__9383_22 vm.Value
	var v24 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _ = arg__9354_3, arg__9359_6, arg__9364_9, arg__9365_10, arg__9372_15, arg__9377_18, arg__9382_21, arg__9383_22, v24
	arg__9354_3, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-target").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__9359_6, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-args").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__9364_9, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-args").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__9365_10, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "format-args").Deref(), []vm.Value{arg__9364_9})
	if callErr != nil {
		return nil, callErr
	}
	arg__9372_15, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-target").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__9377_18, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-args").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__9382_21, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-args").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__9383_22, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "format-args").Deref(), []vm.Value{arg__9382_21})
	if callErr != nil {
		return nil, callErr
	}
	v24, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("b"), arg__9372_15, vm.String("("), arg__9383_22, vm.String(")")})
	if callErr != nil {
		return nil, callErr
	}
	return v24, nil
}
func terminator_targets_str(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var v7 bool
	var op_2 vm.Value
	var aux_3 vm.Value
	var arg__9391_11 vm.Value
	var arg__9397_15 vm.Value
	var v16 vm.Value
	var op_4 vm.Value
	var aux_5 vm.Value
	var v23 bool
	var v74 vm.Value
	var op_75 vm.Value
	var aux_76 vm.Value
	var op_18 vm.Value
	var aux_19 vm.Value
	var arg__9404_27 vm.Value
	var arg__9409_30 vm.Value
	var arg__9410_31 vm.Value
	var arg__9415_34 vm.Value
	var arg__9420_37 vm.Value
	var arg__9421_38 vm.Value
	var arg__9427_42 vm.Value
	var arg__9432_45 vm.Value
	var arg__9433_46 vm.Value
	var arg__9438_49 vm.Value
	var arg__9443_52 vm.Value
	var arg__9444_53 vm.Value
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
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = v7, op_2, aux_3, arg__9391_11, arg__9397_15, v16, op_4, aux_5, v23, v74, op_75, aux_76, op_18, aux_19, arg__9404_27, arg__9409_30, arg__9410_31, arg__9415_34, arg__9420_37, arg__9421_38, arg__9427_42, arg__9432_45, arg__9433_46, arg__9438_49, arg__9443_52, arg__9444_53, v54, op_20, aux_21, v70, op_71, aux_72, op_56, aux_57, op_58, aux_59, v66, op_67, aux_68
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
	arg__9391_11, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "format-target").Deref(), []vm.Value{aux_3})
	if callErr != nil {
		return nil, callErr
	}
	arg__9397_15, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "format-target").Deref(), []vm.Value{aux_3})
	if callErr != nil {
		return nil, callErr
	}
	v16, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(" -> "), arg__9397_15})
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
	arg__9404_27, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-true").Deref(), []vm.Value{aux_19})
	if callErr != nil {
		return nil, callErr
	}
	arg__9409_30, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-true").Deref(), []vm.Value{aux_19})
	if callErr != nil {
		return nil, callErr
	}
	arg__9410_31, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "format-target").Deref(), []vm.Value{arg__9409_30})
	if callErr != nil {
		return nil, callErr
	}
	arg__9415_34, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-false").Deref(), []vm.Value{aux_19})
	if callErr != nil {
		return nil, callErr
	}
	arg__9420_37, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-false").Deref(), []vm.Value{aux_19})
	if callErr != nil {
		return nil, callErr
	}
	arg__9421_38, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "format-target").Deref(), []vm.Value{arg__9420_37})
	if callErr != nil {
		return nil, callErr
	}
	arg__9427_42, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-true").Deref(), []vm.Value{aux_19})
	if callErr != nil {
		return nil, callErr
	}
	arg__9432_45, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-true").Deref(), []vm.Value{aux_19})
	if callErr != nil {
		return nil, callErr
	}
	arg__9433_46, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "format-target").Deref(), []vm.Value{arg__9432_45})
	if callErr != nil {
		return nil, callErr
	}
	arg__9438_49, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-false").Deref(), []vm.Value{aux_19})
	if callErr != nil {
		return nil, callErr
	}
	arg__9443_52, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-false").Deref(), []vm.Value{aux_19})
	if callErr != nil {
		return nil, callErr
	}
	arg__9444_53, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "format-target").Deref(), []vm.Value{arg__9443_52})
	if callErr != nil {
		return nil, callErr
	}
	v54, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(" -> "), arg__9433_46, vm.String(" : "), arg__9444_53})
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
	var arg__9469_25 vm.Value
	var arg__9473_27 vm.Value
	var arg__9479_29 vm.Value
	var arg__9486_34 vm.Value
	var arg__9490_36 vm.Value
	var arg__9496_38 vm.Value
	var v40 vm.Value
	var f_13 vm.Value
	var id_14 vm.Value
	var op_15 vm.Value
	var refs_16 vm.Value
	var aux_17 vm.Value
	var t_43 vm.Value
	var arg__9509_47 vm.Value
	var arg__9513_49 vm.Value
	var arg__9517_73 vm.Value
	var arg__9522_76 vm.Value
	var v77 vm.Value
	var v269 vm.Value
	var f_270 vm.Value
	var id_271 vm.Value
	var op_272 vm.Value
	var refs_273 vm.Value
	var aux_274 vm.Value
	var f_50 vm.Value
	var id_51 vm.Value
	var arg__9504_52 vm.Value
	var op_53 vm.Value
	var refs_54 vm.Value
	var aux_55 vm.Value
	var t_56 vm.Value
	var arg__9503_57 string
	var arg__9505_58 string
	var arg__9509_59 vm.Value
	var arg__9513_60 vm.Value
	var v82 vm.Value
	var f_61 vm.Value
	var id_62 vm.Value
	var arg__9504_63 vm.Value
	var op_64 vm.Value
	var refs_65 vm.Value
	var aux_66 vm.Value
	var t_67 vm.Value
	var arg__9503_68 string
	var arg__9505_69 string
	var arg__9509_70 vm.Value
	var arg__9513_71 vm.Value
	var arg__9528_86 vm.Value
	var f_87 vm.Value
	var id_88 vm.Value
	var arg__9504_89 vm.Value
	var op_90 vm.Value
	var refs_91 vm.Value
	var aux_92 vm.Value
	var t_93 vm.Value
	var arg__9503_94 string
	var arg__9505_95 string
	var arg__9509_96 vm.Value
	var arg__9513_97 vm.Value
	var v125 vm.Value
	var arg__9528_98 vm.Value
	var f_99 vm.Value
	var id_100 vm.Value
	var arg__9504_101 vm.Value
	var op_102 vm.Value
	var refs_103 vm.Value
	var aux_104 vm.Value
	var t_105 vm.Value
	var arg__9503_106 string
	var arg__9505_107 string
	var arg__9509_108 vm.Value
	var arg__9513_109 vm.Value
	var arg__9538_129 vm.Value
	var arg__9544_133 vm.Value
	var v134 vm.Value
	var arg__9528_110 vm.Value
	var f_111 vm.Value
	var id_112 vm.Value
	var arg__9504_113 vm.Value
	var op_114 vm.Value
	var refs_115 vm.Value
	var aux_116 vm.Value
	var t_117 vm.Value
	var arg__9503_118 string
	var arg__9505_119 string
	var arg__9509_120 vm.Value
	var arg__9513_121 vm.Value
	var arg__9545_138 vm.Value
	var arg__9528_139 vm.Value
	var f_140 vm.Value
	var id_141 vm.Value
	var arg__9504_142 vm.Value
	var op_143 vm.Value
	var refs_144 vm.Value
	var aux_145 vm.Value
	var t_146 vm.Value
	var arg__9503_147 string
	var arg__9505_148 string
	var arg__9509_149 vm.Value
	var arg__9513_150 vm.Value
	var arg__9554_156 vm.Value
	var arg__9558_158 vm.Value
	var arg__9562_184 vm.Value
	var arg__9567_187 vm.Value
	var v188 vm.Value
	var f_159 vm.Value
	var id_160 vm.Value
	var arg__9549_161 vm.Value
	var op_162 vm.Value
	var refs_163 vm.Value
	var aux_164 vm.Value
	var t_165 vm.Value
	var head__9547_166 vm.Value
	var arg__9548_167 string
	var arg__9550_168 string
	var arg__9554_169 vm.Value
	var arg__9558_170 vm.Value
	var v193 vm.Value
	var f_171 vm.Value
	var id_172 vm.Value
	var arg__9549_173 vm.Value
	var op_174 vm.Value
	var refs_175 vm.Value
	var aux_176 vm.Value
	var t_177 vm.Value
	var head__9547_178 vm.Value
	var arg__9548_179 string
	var arg__9550_180 string
	var arg__9554_181 vm.Value
	var arg__9558_182 vm.Value
	var arg__9573_197 vm.Value
	var f_198 vm.Value
	var id_199 vm.Value
	var arg__9549_200 vm.Value
	var op_201 vm.Value
	var refs_202 vm.Value
	var aux_203 vm.Value
	var t_204 vm.Value
	var head__9547_205 vm.Value
	var arg__9548_206 string
	var arg__9550_207 string
	var arg__9554_208 vm.Value
	var arg__9558_209 vm.Value
	var v239 vm.Value
	var arg__9573_210 vm.Value
	var f_211 vm.Value
	var id_212 vm.Value
	var arg__9549_213 vm.Value
	var op_214 vm.Value
	var refs_215 vm.Value
	var aux_216 vm.Value
	var t_217 vm.Value
	var head__9547_218 vm.Value
	var arg__9548_219 string
	var arg__9550_220 string
	var arg__9554_221 vm.Value
	var arg__9558_222 vm.Value
	var arg__9583_243 vm.Value
	var arg__9589_247 vm.Value
	var v248 vm.Value
	var arg__9573_223 vm.Value
	var f_224 vm.Value
	var id_225 vm.Value
	var arg__9549_226 vm.Value
	var op_227 vm.Value
	var refs_228 vm.Value
	var aux_229 vm.Value
	var t_230 vm.Value
	var head__9547_231 vm.Value
	var arg__9548_232 string
	var arg__9550_233 string
	var arg__9554_234 vm.Value
	var arg__9558_235 vm.Value
	var arg__9590_252 vm.Value
	var arg__9573_253 vm.Value
	var f_254 vm.Value
	var id_255 vm.Value
	var arg__9549_256 vm.Value
	var op_257 vm.Value
	var refs_258 vm.Value
	var aux_259 vm.Value
	var t_260 vm.Value
	var head__9547_261 vm.Value
	var arg__9548_262 string
	var arg__9550_263 string
	var arg__9554_264 vm.Value
	var arg__9558_265 vm.Value
	var v267 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = op_3, refs_5, aux_7, v21, f_8, id_9, op_10, refs_11, aux_12, arg__9469_25, arg__9473_27, arg__9479_29, arg__9486_34, arg__9490_36, arg__9496_38, v40, f_13, id_14, op_15, refs_16, aux_17, t_43, arg__9509_47, arg__9513_49, arg__9517_73, arg__9522_76, v77, v269, f_270, id_271, op_272, refs_273, aux_274, f_50, id_51, arg__9504_52, op_53, refs_54, aux_55, t_56, arg__9503_57, arg__9505_58, arg__9509_59, arg__9513_60, v82, f_61, id_62, arg__9504_63, op_64, refs_65, aux_66, t_67, arg__9503_68, arg__9505_69, arg__9509_70, arg__9513_71, arg__9528_86, f_87, id_88, arg__9504_89, op_90, refs_91, aux_92, t_93, arg__9503_94, arg__9505_95, arg__9509_96, arg__9513_97, v125, arg__9528_98, f_99, id_100, arg__9504_101, op_102, refs_103, aux_104, t_105, arg__9503_106, arg__9505_107, arg__9509_108, arg__9513_109, arg__9538_129, arg__9544_133, v134, arg__9528_110, f_111, id_112, arg__9504_113, op_114, refs_115, aux_116, t_117, arg__9503_118, arg__9505_119, arg__9509_120, arg__9513_121, arg__9545_138, arg__9528_139, f_140, id_141, arg__9504_142, op_143, refs_144, aux_145, t_146, arg__9503_147, arg__9505_148, arg__9509_149, arg__9513_150, arg__9554_156, arg__9558_158, arg__9562_184, arg__9567_187, v188, f_159, id_160, arg__9549_161, op_162, refs_163, aux_164, t_165, head__9547_166, arg__9548_167, arg__9550_168, arg__9554_169, arg__9558_170, v193, f_171, id_172, arg__9549_173, op_174, refs_175, aux_176, t_177, head__9547_178, arg__9548_179, arg__9550_180, arg__9554_181, arg__9558_182, arg__9573_197, f_198, id_199, arg__9549_200, op_201, refs_202, aux_203, t_204, head__9547_205, arg__9548_206, arg__9550_207, arg__9554_208, arg__9558_209, v239, arg__9573_210, f_211, id_212, arg__9549_213, op_214, refs_215, aux_216, t_217, head__9547_218, arg__9548_219, arg__9550_220, arg__9554_221, arg__9558_222, arg__9583_243, arg__9589_247, v248, arg__9573_223, f_224, id_225, arg__9549_226, op_227, refs_228, aux_229, t_230, head__9547_231, arg__9548_232, arg__9550_233, arg__9554_234, arg__9558_235, arg__9590_252, arg__9573_253, f_254, id_255, arg__9549_256, op_257, refs_258, aux_259, t_260, head__9547_261, arg__9548_262, arg__9550_263, arg__9554_264, arg__9558_265, v267
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
	arg__9469_25, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "op-display-name").Deref(), []vm.Value{op_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__9473_27, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "format-refs").Deref(), []vm.Value{refs_11})
	if callErr != nil {
		return nil, callErr
	}
	arg__9479_29, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "terminator-targets-str").Deref(), []vm.Value{op_10, aux_12})
	if callErr != nil {
		return nil, callErr
	}
	arg__9486_34, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "op-display-name").Deref(), []vm.Value{op_10})
	if callErr != nil {
		return nil, callErr
	}
	arg__9490_36, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "format-refs").Deref(), []vm.Value{refs_11})
	if callErr != nil {
		return nil, callErr
	}
	arg__9496_38, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "terminator-targets-str").Deref(), []vm.Value{op_10, aux_12})
	if callErr != nil {
		return nil, callErr
	}
	v40, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("    "), arg__9486_34, arg__9490_36, arg__9496_38, vm.String("\n")})
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
	arg__9509_47, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "op-display-name").Deref(), []vm.Value{op_15})
	if callErr != nil {
		return nil, callErr
	}
	arg__9513_49, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "format-refs").Deref(), []vm.Value{refs_16})
	if callErr != nil {
		return nil, callErr
	}
	arg__9517_73, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{aux_17})
	if callErr != nil {
		return nil, callErr
	}
	arg__9522_76, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{aux_17})
	if callErr != nil {
		return nil, callErr
	}
	v77, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not").Deref(), []vm.Value{arg__9522_76})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v77) {
		f_50 = f_13
		id_51 = id_14
		arg__9504_52 = id_14
		op_53 = op_15
		refs_54 = refs_16
		aux_55 = aux_17
		t_56 = t_43
		arg__9503_57 = "    v"
		arg__9505_58 = " = "
		arg__9509_59 = arg__9509_47
		arg__9513_60 = arg__9513_49
		goto b4
	} else {
		f_61 = f_13
		id_62 = id_14
		arg__9504_63 = id_14
		op_64 = op_15
		refs_65 = refs_16
		aux_66 = aux_17
		t_67 = t_43
		arg__9503_68 = "    v"
		arg__9505_69 = " = "
		arg__9509_70 = arg__9509_47
		arg__9513_71 = arg__9513_49
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
	arg__9528_86 = v82
	f_87 = f_50
	id_88 = id_51
	arg__9504_89 = arg__9504_52
	op_90 = op_53
	refs_91 = refs_54
	aux_92 = aux_55
	t_93 = t_56
	arg__9503_94 = arg__9503_57
	arg__9505_95 = arg__9505_58
	arg__9509_96 = arg__9509_59
	arg__9513_97 = arg__9513_60
	goto b6
b5:
	;
	arg__9528_86 = vm.String("")
	f_87 = f_61
	id_88 = id_62
	arg__9504_89 = arg__9504_63
	op_90 = op_64
	refs_91 = refs_65
	aux_92 = aux_66
	t_93 = t_67
	arg__9503_94 = arg__9503_68
	arg__9505_95 = arg__9505_69
	arg__9509_96 = arg__9509_70
	arg__9513_97 = arg__9513_71
	goto b6
b6:
	;
	v125, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not=").Deref(), []vm.Value{t_93, vm.Keyword("unknown")})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v125) {
		arg__9528_98 = arg__9528_86
		f_99 = f_87
		id_100 = id_88
		arg__9504_101 = arg__9504_89
		op_102 = op_90
		refs_103 = refs_91
		aux_104 = aux_92
		t_105 = t_93
		arg__9503_106 = arg__9503_94
		arg__9505_107 = arg__9505_95
		arg__9509_108 = arg__9509_96
		arg__9513_109 = arg__9513_97
		goto b7
	} else {
		arg__9528_110 = arg__9528_86
		f_111 = f_87
		id_112 = id_88
		arg__9504_113 = arg__9504_89
		op_114 = op_90
		refs_115 = refs_91
		aux_116 = aux_92
		t_117 = t_93
		arg__9503_118 = arg__9503_94
		arg__9505_119 = arg__9505_95
		arg__9509_120 = arg__9509_96
		arg__9513_121 = arg__9513_97
		goto b8
	}
b7:
	;
	arg__9538_129, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "type-display").Deref(), []vm.Value{t_105})
	if callErr != nil {
		return nil, callErr
	}
	arg__9544_133, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "type-display").Deref(), []vm.Value{t_105})
	if callErr != nil {
		return nil, callErr
	}
	v134, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(" : "), arg__9544_133})
	if callErr != nil {
		return nil, callErr
	}
	arg__9545_138 = v134
	arg__9528_139 = arg__9528_98
	f_140 = f_99
	id_141 = id_100
	arg__9504_142 = arg__9504_101
	op_143 = op_102
	refs_144 = refs_103
	aux_145 = aux_104
	t_146 = t_105
	arg__9503_147 = arg__9503_106
	arg__9505_148 = arg__9505_107
	arg__9509_149 = arg__9509_108
	arg__9513_150 = arg__9513_109
	goto b9
b8:
	;
	arg__9545_138 = vm.String("")
	arg__9528_139 = arg__9528_110
	f_140 = f_111
	id_141 = id_112
	arg__9504_142 = arg__9504_113
	op_143 = op_114
	refs_144 = refs_115
	aux_145 = aux_116
	t_146 = t_117
	arg__9503_147 = arg__9503_118
	arg__9505_148 = arg__9505_119
	arg__9509_149 = arg__9509_120
	arg__9513_150 = arg__9513_121
	goto b9
b9:
	;
	arg__9554_156, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "op-display-name").Deref(), []vm.Value{op_143})
	if callErr != nil {
		return nil, callErr
	}
	arg__9558_158, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "format-refs").Deref(), []vm.Value{refs_144})
	if callErr != nil {
		return nil, callErr
	}
	arg__9562_184, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{aux_145})
	if callErr != nil {
		return nil, callErr
	}
	arg__9567_187, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{aux_145})
	if callErr != nil {
		return nil, callErr
	}
	v188, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not").Deref(), []vm.Value{arg__9567_187})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v188) {
		f_159 = f_140
		id_160 = id_141
		arg__9549_161 = id_141
		op_162 = op_143
		refs_163 = refs_144
		aux_164 = aux_145
		t_165 = t_146
		head__9547_166 = rt.LookupVar("clojure.core", "str").Deref()
		arg__9548_167 = "    v"
		arg__9550_168 = " = "
		arg__9554_169 = arg__9554_156
		arg__9558_170 = arg__9558_158
		goto b10
	} else {
		f_171 = f_140
		id_172 = id_141
		arg__9549_173 = id_141
		op_174 = op_143
		refs_175 = refs_144
		aux_176 = aux_145
		t_177 = t_146
		head__9547_178 = rt.LookupVar("clojure.core", "str").Deref()
		arg__9548_179 = "    v"
		arg__9550_180 = " = "
		arg__9554_181 = arg__9554_156
		arg__9558_182 = arg__9558_158
		goto b11
	}
b10:
	;
	v193, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(" ; "), aux_164})
	if callErr != nil {
		return nil, callErr
	}
	arg__9573_197 = v193
	f_198 = f_159
	id_199 = id_160
	arg__9549_200 = arg__9549_161
	op_201 = op_162
	refs_202 = refs_163
	aux_203 = aux_164
	t_204 = t_165
	head__9547_205 = head__9547_166
	arg__9548_206 = arg__9548_167
	arg__9550_207 = arg__9550_168
	arg__9554_208 = arg__9554_169
	arg__9558_209 = arg__9558_170
	goto b12
b11:
	;
	arg__9573_197 = vm.String("")
	f_198 = f_171
	id_199 = id_172
	arg__9549_200 = arg__9549_173
	op_201 = op_174
	refs_202 = refs_175
	aux_203 = aux_176
	t_204 = t_177
	head__9547_205 = head__9547_178
	arg__9548_206 = arg__9548_179
	arg__9550_207 = arg__9550_180
	arg__9554_208 = arg__9554_181
	arg__9558_209 = arg__9558_182
	goto b12
b12:
	;
	v239, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not=").Deref(), []vm.Value{t_204, vm.Keyword("unknown")})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v239) {
		arg__9573_210 = arg__9573_197
		f_211 = f_198
		id_212 = id_199
		arg__9549_213 = arg__9549_200
		op_214 = op_201
		refs_215 = refs_202
		aux_216 = aux_203
		t_217 = t_204
		head__9547_218 = head__9547_205
		arg__9548_219 = arg__9548_206
		arg__9550_220 = arg__9550_207
		arg__9554_221 = arg__9554_208
		arg__9558_222 = arg__9558_209
		goto b13
	} else {
		arg__9573_223 = arg__9573_197
		f_224 = f_198
		id_225 = id_199
		arg__9549_226 = arg__9549_200
		op_227 = op_201
		refs_228 = refs_202
		aux_229 = aux_203
		t_230 = t_204
		head__9547_231 = head__9547_205
		arg__9548_232 = arg__9548_206
		arg__9550_233 = arg__9550_207
		arg__9554_234 = arg__9554_208
		arg__9558_235 = arg__9558_209
		goto b14
	}
b13:
	;
	arg__9583_243, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "type-display").Deref(), []vm.Value{t_217})
	if callErr != nil {
		return nil, callErr
	}
	arg__9589_247, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "type-display").Deref(), []vm.Value{t_217})
	if callErr != nil {
		return nil, callErr
	}
	v248, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(" : "), arg__9589_247})
	if callErr != nil {
		return nil, callErr
	}
	arg__9590_252 = v248
	arg__9573_253 = arg__9573_210
	f_254 = f_211
	id_255 = id_212
	arg__9549_256 = arg__9549_213
	op_257 = op_214
	refs_258 = refs_215
	aux_259 = aux_216
	t_260 = t_217
	head__9547_261 = head__9547_218
	arg__9548_262 = arg__9548_219
	arg__9550_263 = arg__9550_220
	arg__9554_264 = arg__9554_221
	arg__9558_265 = arg__9558_222
	goto b15
b14:
	;
	arg__9590_252 = vm.String("")
	arg__9573_253 = arg__9573_223
	f_254 = f_224
	id_255 = id_225
	arg__9549_256 = arg__9549_226
	op_257 = op_227
	refs_258 = refs_228
	aux_259 = aux_229
	t_260 = t_230
	head__9547_261 = head__9547_231
	arg__9548_262 = arg__9548_232
	arg__9550_263 = arg__9550_233
	arg__9554_264 = arg__9554_234
	arg__9558_265 = arg__9558_235
	goto b15
b15:
	;
	v267, callErr = rt.InvokeValue(head__9547_261, []vm.Value{vm.String(arg__9548_262), arg__9549_256, vm.String(arg__9550_263), arg__9554_264, arg__9558_265, arg__9573_253, arg__9590_252, vm.String("\n")})
	if callErr != nil {
		return nil, callErr
	}
	v269 = v267
	f_270 = f_254
	id_271 = id_255
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
	var arg__9616_11 vm.Value
	var entry_QMARK__12 bool
	var f_14 vm.Value
	var bid_15 vm.Value
	var params_16 vm.Value
	var preds_17 vm.Value
	var insts_18 vm.Value
	var term_19 vm.Value
	var entry_QMARK__20 bool
	var arg__9617_21 string
	var f_22 vm.Value
	var bid_23 vm.Value
	var params_24 vm.Value
	var preds_25 vm.Value
	var insts_26 vm.Value
	var term_27 vm.Value
	var entry_QMARK__28 bool
	var arg__9617_29 string
	var arg__9618_35 string
	var f_36 vm.Value
	var bid_37 vm.Value
	var params_38 vm.Value
	var preds_39 vm.Value
	var insts_40 vm.Value
	var term_41 vm.Value
	var entry_QMARK__42 vm.Value
	var arg__9617_43 string
	var arg__9698_54 vm.Value
	var arg__9776_64 vm.Value
	var arg__9777_65 vm.Value
	var v96 vm.Value
	var arg__9618_67 string
	var f_68 vm.Value
	var arg__9620_69 vm.Value
	var bid_70 vm.Value
	var params_71 vm.Value
	var preds_72 vm.Value
	var insts_73 vm.Value
	var term_74 vm.Value
	var entry_QMARK__75 vm.Value
	var arg__9617_76 string
	var arg__9619_77 string
	var arg__9621_78 string
	var arg__9777_79 vm.Value
	var arg__9778_80 string
	var arg__9799_103 vm.Value
	var arg__9817_109 vm.Value
	var arg__9818_110 vm.Value
	var arg__9837_117 vm.Value
	var arg__9855_123 vm.Value
	var arg__9856_124 vm.Value
	var v125 vm.Value
	var arg__9618_81 string
	var f_82 vm.Value
	var arg__9620_83 vm.Value
	var bid_84 vm.Value
	var params_85 vm.Value
	var preds_86 vm.Value
	var insts_87 vm.Value
	var term_88 vm.Value
	var entry_QMARK__89 vm.Value
	var arg__9617_90 string
	var arg__9619_91 string
	var arg__9621_92 string
	var arg__9777_93 vm.Value
	var arg__9778_94 string
	var arg__9857_129 vm.Value
	var arg__9618_130 string
	var f_131 vm.Value
	var arg__9620_132 vm.Value
	var bid_133 vm.Value
	var params_134 vm.Value
	var preds_135 vm.Value
	var insts_136 vm.Value
	var term_137 vm.Value
	var entry_QMARK__138 vm.Value
	var arg__9617_139 string
	var arg__9619_140 string
	var arg__9621_141 string
	var arg__9777_142 vm.Value
	var arg__9778_143 string
	var f_147 vm.Value
	var bid_148 vm.Value
	var params_149 vm.Value
	var preds_150 vm.Value
	var insts_151 vm.Value
	var term_152 vm.Value
	var entry_QMARK__153 bool
	var head__9859_154 vm.Value
	var arg__9860_155 string
	var f_156 vm.Value
	var bid_157 vm.Value
	var params_158 vm.Value
	var preds_159 vm.Value
	var insts_160 vm.Value
	var term_161 vm.Value
	var entry_QMARK__162 bool
	var head__9859_163 vm.Value
	var arg__9860_164 string
	var arg__9861_170 string
	var f_171 vm.Value
	var bid_172 vm.Value
	var params_173 vm.Value
	var preds_174 vm.Value
	var insts_175 vm.Value
	var term_176 vm.Value
	var entry_QMARK__177 vm.Value
	var head__9859_178 vm.Value
	var arg__9860_179 string
	var arg__9941_190 vm.Value
	var arg__10019_200 vm.Value
	var arg__10020_201 vm.Value
	var v234 vm.Value
	var arg__9861_203 string
	var f_204 vm.Value
	var arg__9863_205 vm.Value
	var bid_206 vm.Value
	var params_207 vm.Value
	var preds_208 vm.Value
	var insts_209 vm.Value
	var term_210 vm.Value
	var entry_QMARK__211 vm.Value
	var head__9859_212 vm.Value
	var arg__9860_213 string
	var arg__9862_214 string
	var arg__9864_215 string
	var arg__10020_216 vm.Value
	var arg__10021_217 string
	var arg__10042_241 vm.Value
	var arg__10060_247 vm.Value
	var arg__10061_248 vm.Value
	var arg__10080_255 vm.Value
	var arg__10098_261 vm.Value
	var arg__10099_262 vm.Value
	var v263 vm.Value
	var arg__9861_218 string
	var f_219 vm.Value
	var arg__9863_220 vm.Value
	var bid_221 vm.Value
	var params_222 vm.Value
	var preds_223 vm.Value
	var insts_224 vm.Value
	var term_225 vm.Value
	var entry_QMARK__226 vm.Value
	var head__9859_227 vm.Value
	var arg__9860_228 string
	var arg__9862_229 string
	var arg__9864_230 string
	var arg__10020_231 vm.Value
	var arg__10021_232 string
	var arg__10100_267 vm.Value
	var arg__9861_268 string
	var f_269 vm.Value
	var arg__9863_270 vm.Value
	var bid_271 vm.Value
	var params_272 vm.Value
	var preds_273 vm.Value
	var insts_274 vm.Value
	var term_275 vm.Value
	var entry_QMARK__276 vm.Value
	var head__9859_277 vm.Value
	var arg__9860_278 string
	var arg__9862_279 string
	var arg__9864_280 string
	var arg__10020_281 vm.Value
	var arg__10021_282 string
	var header_284 vm.Value
	var arg__10118_293 vm.Value
	var arg__10136_303 vm.Value
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
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = params_3, preds_5, insts_7, term_9, arg__9616_11, entry_QMARK__12, f_14, bid_15, params_16, preds_17, insts_18, term_19, entry_QMARK__20, arg__9617_21, f_22, bid_23, params_24, preds_25, insts_26, term_27, entry_QMARK__28, arg__9617_29, arg__9618_35, f_36, bid_37, params_38, preds_39, insts_40, term_41, entry_QMARK__42, arg__9617_43, arg__9698_54, arg__9776_64, arg__9777_65, v96, arg__9618_67, f_68, arg__9620_69, bid_70, params_71, preds_72, insts_73, term_74, entry_QMARK__75, arg__9617_76, arg__9619_77, arg__9621_78, arg__9777_79, arg__9778_80, arg__9799_103, arg__9817_109, arg__9818_110, arg__9837_117, arg__9855_123, arg__9856_124, v125, arg__9618_81, f_82, arg__9620_83, bid_84, params_85, preds_86, insts_87, term_88, entry_QMARK__89, arg__9617_90, arg__9619_91, arg__9621_92, arg__9777_93, arg__9778_94, arg__9857_129, arg__9618_130, f_131, arg__9620_132, bid_133, params_134, preds_135, insts_136, term_137, entry_QMARK__138, arg__9617_139, arg__9619_140, arg__9621_141, arg__9777_142, arg__9778_143, f_147, bid_148, params_149, preds_150, insts_151, term_152, entry_QMARK__153, head__9859_154, arg__9860_155, f_156, bid_157, params_158, preds_159, insts_160, term_161, entry_QMARK__162, head__9859_163, arg__9860_164, arg__9861_170, f_171, bid_172, params_173, preds_174, insts_175, term_176, entry_QMARK__177, head__9859_178, arg__9860_179, arg__9941_190, arg__10019_200, arg__10020_201, v234, arg__9861_203, f_204, arg__9863_205, bid_206, params_207, preds_208, insts_209, term_210, entry_QMARK__211, head__9859_212, arg__9860_213, arg__9862_214, arg__9864_215, arg__10020_216, arg__10021_217, arg__10042_241, arg__10060_247, arg__10061_248, arg__10080_255, arg__10098_261, arg__10099_262, v263, arg__9861_218, f_219, arg__9863_220, bid_221, params_222, preds_223, insts_224, term_225, entry_QMARK__226, head__9859_227, arg__9860_228, arg__9862_229, arg__9864_230, arg__10020_231, arg__10021_232, arg__10100_267, arg__9861_268, f_269, arg__9863_270, bid_271, params_272, preds_273, insts_274, term_275, entry_QMARK__276, head__9859_277, arg__9860_278, arg__9862_279, arg__9864_280, arg__10020_281, arg__10021_282, header_284, arg__10118_293, arg__10136_303, body_304, or__x_326, f_305, bid_306, params_307, preds_308, insts_309, term_310, entry_QMARK__311, header_312, body_313, v365, f_314, bid_315, params_316, preds_317, insts_318, term_319, entry_QMARK__320, header_321, body_322, term_line_369, f_370, bid_371, params_372, preds_373, insts_374, term_375, entry_QMARK__376, header_377, body_378, v382, f_327, bid_328, params_329, preds_330, insts_331, term_332, entry_QMARK__333, header_334, body_335, or__x_336, f_337, bid_338, params_339, preds_340, insts_341, term_342, entry_QMARK__343, header_344, body_345, or__x_346, v350, v352, f_353, bid_354, params_355, preds_356, insts_357, term_358, entry_QMARK__359, header_360, body_361, or__x_362
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
	arg__9616_11, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-entry").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	entry_QMARK__12 = arg1 == arg__9616_11
	if entry_QMARK__12 {
		f_14 = arg0
		bid_15 = arg1
		params_16 = params_3
		preds_17 = preds_5
		insts_18 = insts_7
		term_19 = term_9
		entry_QMARK__20 = entry_QMARK__12
		arg__9617_21 = "  "
		goto b1
	} else {
		f_22 = arg0
		bid_23 = arg1
		params_24 = params_3
		preds_25 = preds_5
		insts_26 = insts_7
		term_27 = term_9
		entry_QMARK__28 = entry_QMARK__12
		arg__9617_29 = "  "
		goto b2
	}
b1:
	;
	arg__9618_35 = "entry "
	f_36 = f_14
	bid_37 = bid_15
	params_38 = params_16
	preds_39 = preds_17
	insts_40 = insts_18
	term_41 = term_19
	entry_QMARK__42 = vm.Boolean(entry_QMARK__20)
	arg__9617_43 = arg__9617_21
	goto b3
b2:
	;
	arg__9618_35 = ""
	f_36 = f_22
	bid_37 = bid_23
	params_38 = params_24
	preds_39 = preds_25
	insts_40 = insts_26
	term_41 = term_27
	entry_QMARK__42 = vm.Boolean(entry_QMARK__28)
	arg__9617_43 = arg__9617_29
	goto b3
b3:
	;
	arg__9698_54, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var arg__9669_5 vm.Value
		var arg__9676_8 vm.Value
		var arg__9677_9 vm.Value
		var arg__9687_14 vm.Value
		var arg__9694_17 vm.Value
		var arg__9695_18 vm.Value
		var v19 vm.Value
		var callErr error
		_, _, _, _, _, _, _ = arg__9669_5, arg__9676_8, arg__9677_9, arg__9687_14, arg__9694_17, arg__9695_18, v19
		arg__9669_5, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_36})
		if callErr != nil {
			return nil, callErr
		}
		arg__9676_8, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_36})
		if callErr != nil {
			return nil, callErr
		}
		arg__9677_9, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "type-display").Deref(), []vm.Value{arg__9676_8})
		if callErr != nil {
			return nil, callErr
		}
		arg__9687_14, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_36})
		if callErr != nil {
			return nil, callErr
		}
		arg__9694_17, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_36})
		if callErr != nil {
			return nil, callErr
		}
		arg__9695_18, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "type-display").Deref(), []vm.Value{arg__9694_17})
		if callErr != nil {
			return nil, callErr
		}
		v19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("v"), arg0, vm.String(": "), arg__9695_18})
		if callErr != nil {
			return nil, callErr
		}
		return v19, nil
	}), params_38})
	if callErr != nil {
		return nil, callErr
	}
	arg__9776_64, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var arg__9747_5 vm.Value
		var arg__9754_8 vm.Value
		var arg__9755_9 vm.Value
		var arg__9765_14 vm.Value
		var arg__9772_17 vm.Value
		var arg__9773_18 vm.Value
		var v19 vm.Value
		var callErr error
		_, _, _, _, _, _, _ = arg__9747_5, arg__9754_8, arg__9755_9, arg__9765_14, arg__9772_17, arg__9773_18, v19
		arg__9747_5, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_36})
		if callErr != nil {
			return nil, callErr
		}
		arg__9754_8, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_36})
		if callErr != nil {
			return nil, callErr
		}
		arg__9755_9, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "type-display").Deref(), []vm.Value{arg__9754_8})
		if callErr != nil {
			return nil, callErr
		}
		arg__9765_14, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_36})
		if callErr != nil {
			return nil, callErr
		}
		arg__9772_17, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_36})
		if callErr != nil {
			return nil, callErr
		}
		arg__9773_18, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "type-display").Deref(), []vm.Value{arg__9772_17})
		if callErr != nil {
			return nil, callErr
		}
		v19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("v"), arg0, vm.String(": "), arg__9773_18})
		if callErr != nil {
			return nil, callErr
		}
		return v19, nil
	}), params_38})
	if callErr != nil {
		return nil, callErr
	}
	arg__9777_65, callErr = rt.InvokeValue(rt.LookupVar("clojure.string", "join").Deref(), []vm.Value{vm.String(", "), arg__9776_64})
	if callErr != nil {
		return nil, callErr
	}
	v96, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{preds_39})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v96) {
		arg__9618_67 = arg__9618_35
		f_68 = f_36
		arg__9620_69 = bid_37
		bid_70 = bid_37
		params_71 = params_38
		preds_72 = preds_39
		insts_73 = insts_40
		term_74 = term_41
		entry_QMARK__75 = entry_QMARK__42
		arg__9617_76 = arg__9617_43
		arg__9619_77 = "b"
		arg__9621_78 = "("
		arg__9777_79 = arg__9777_65
		arg__9778_80 = "):"
		goto b4
	} else {
		arg__9618_81 = arg__9618_35
		f_82 = f_36
		arg__9620_83 = bid_37
		bid_84 = bid_37
		params_85 = params_38
		preds_86 = preds_39
		insts_87 = insts_40
		term_88 = term_41
		entry_QMARK__89 = entry_QMARK__42
		arg__9617_90 = arg__9617_43
		arg__9619_91 = "b"
		arg__9621_92 = "("
		arg__9777_93 = arg__9777_65
		arg__9778_94 = "):"
		goto b5
	}
b4:
	;
	arg__9799_103, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
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
	arg__9817_109, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
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
	arg__9818_110, callErr = rt.InvokeValue(rt.LookupVar("clojure.string", "join").Deref(), []vm.Value{vm.String(", "), arg__9817_109})
	if callErr != nil {
		return nil, callErr
	}
	arg__9837_117, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
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
	arg__9855_123, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
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
	arg__9856_124, callErr = rt.InvokeValue(rt.LookupVar("clojure.string", "join").Deref(), []vm.Value{vm.String(", "), arg__9855_123})
	if callErr != nil {
		return nil, callErr
	}
	v125, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("    ; preds: "), arg__9856_124})
	if callErr != nil {
		return nil, callErr
	}
	arg__9857_129 = v125
	arg__9618_130 = arg__9618_67
	f_131 = f_68
	arg__9620_132 = arg__9620_69
	bid_133 = bid_70
	params_134 = params_71
	preds_135 = preds_72
	insts_136 = insts_73
	term_137 = term_74
	entry_QMARK__138 = entry_QMARK__75
	arg__9617_139 = arg__9617_76
	arg__9619_140 = arg__9619_77
	arg__9621_141 = arg__9621_78
	arg__9777_142 = arg__9777_79
	arg__9778_143 = arg__9778_80
	goto b6
b5:
	;
	arg__9857_129 = vm.String("")
	arg__9618_130 = arg__9618_81
	f_131 = f_82
	arg__9620_132 = arg__9620_83
	bid_133 = bid_84
	params_134 = params_85
	preds_135 = preds_86
	insts_136 = insts_87
	term_137 = term_88
	entry_QMARK__138 = entry_QMARK__89
	arg__9617_139 = arg__9617_90
	arg__9619_140 = arg__9619_91
	arg__9621_141 = arg__9621_92
	arg__9777_142 = arg__9777_93
	arg__9778_143 = arg__9778_94
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
		head__9859_154 = rt.LookupVar("clojure.core", "str").Deref()
		arg__9860_155 = "  "
		goto b7
	} else {
		f_156 = f_131
		bid_157 = bid_133
		params_158 = params_134
		preds_159 = preds_135
		insts_160 = insts_136
		term_161 = term_137
		entry_QMARK__162 = vm.IsTruthy(entry_QMARK__138)
		head__9859_163 = rt.LookupVar("clojure.core", "str").Deref()
		arg__9860_164 = "  "
		goto b8
	}
b7:
	;
	arg__9861_170 = "entry "
	f_171 = f_147
	bid_172 = bid_148
	params_173 = params_149
	preds_174 = preds_150
	insts_175 = insts_151
	term_176 = term_152
	entry_QMARK__177 = vm.Boolean(entry_QMARK__153)
	head__9859_178 = head__9859_154
	arg__9860_179 = arg__9860_155
	goto b9
b8:
	;
	arg__9861_170 = ""
	f_171 = f_156
	bid_172 = bid_157
	params_173 = params_158
	preds_174 = preds_159
	insts_175 = insts_160
	term_176 = term_161
	entry_QMARK__177 = vm.Boolean(entry_QMARK__162)
	head__9859_178 = head__9859_163
	arg__9860_179 = arg__9860_164
	goto b9
b9:
	;
	arg__9941_190, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var arg__9912_5 vm.Value
		var arg__9919_8 vm.Value
		var arg__9920_9 vm.Value
		var arg__9930_14 vm.Value
		var arg__9937_17 vm.Value
		var arg__9938_18 vm.Value
		var v19 vm.Value
		var callErr error
		_, _, _, _, _, _, _ = arg__9912_5, arg__9919_8, arg__9920_9, arg__9930_14, arg__9937_17, arg__9938_18, v19
		arg__9912_5, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_171})
		if callErr != nil {
			return nil, callErr
		}
		arg__9919_8, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_171})
		if callErr != nil {
			return nil, callErr
		}
		arg__9920_9, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "type-display").Deref(), []vm.Value{arg__9919_8})
		if callErr != nil {
			return nil, callErr
		}
		arg__9930_14, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_171})
		if callErr != nil {
			return nil, callErr
		}
		arg__9937_17, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_171})
		if callErr != nil {
			return nil, callErr
		}
		arg__9938_18, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "type-display").Deref(), []vm.Value{arg__9937_17})
		if callErr != nil {
			return nil, callErr
		}
		v19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("v"), arg0, vm.String(": "), arg__9938_18})
		if callErr != nil {
			return nil, callErr
		}
		return v19, nil
	}), params_173})
	if callErr != nil {
		return nil, callErr
	}
	arg__10019_200, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var arg__9990_5 vm.Value
		var arg__9997_8 vm.Value
		var arg__9998_9 vm.Value
		var arg__10008_14 vm.Value
		var arg__10015_17 vm.Value
		var arg__10016_18 vm.Value
		var v19 vm.Value
		var callErr error
		_, _, _, _, _, _, _ = arg__9990_5, arg__9997_8, arg__9998_9, arg__10008_14, arg__10015_17, arg__10016_18, v19
		arg__9990_5, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_171})
		if callErr != nil {
			return nil, callErr
		}
		arg__9997_8, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_171})
		if callErr != nil {
			return nil, callErr
		}
		arg__9998_9, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "type-display").Deref(), []vm.Value{arg__9997_8})
		if callErr != nil {
			return nil, callErr
		}
		arg__10008_14, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_171})
		if callErr != nil {
			return nil, callErr
		}
		arg__10015_17, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_171})
		if callErr != nil {
			return nil, callErr
		}
		arg__10016_18, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "type-display").Deref(), []vm.Value{arg__10015_17})
		if callErr != nil {
			return nil, callErr
		}
		v19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("v"), arg0, vm.String(": "), arg__10016_18})
		if callErr != nil {
			return nil, callErr
		}
		return v19, nil
	}), params_173})
	if callErr != nil {
		return nil, callErr
	}
	arg__10020_201, callErr = rt.InvokeValue(rt.LookupVar("clojure.string", "join").Deref(), []vm.Value{vm.String(", "), arg__10019_200})
	if callErr != nil {
		return nil, callErr
	}
	v234, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{preds_174})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v234) {
		arg__9861_203 = arg__9861_170
		f_204 = f_171
		arg__9863_205 = bid_172
		bid_206 = bid_172
		params_207 = params_173
		preds_208 = preds_174
		insts_209 = insts_175
		term_210 = term_176
		entry_QMARK__211 = entry_QMARK__177
		head__9859_212 = head__9859_178
		arg__9860_213 = arg__9860_179
		arg__9862_214 = "b"
		arg__9864_215 = "("
		arg__10020_216 = arg__10020_201
		arg__10021_217 = "):"
		goto b10
	} else {
		arg__9861_218 = arg__9861_170
		f_219 = f_171
		arg__9863_220 = bid_172
		bid_221 = bid_172
		params_222 = params_173
		preds_223 = preds_174
		insts_224 = insts_175
		term_225 = term_176
		entry_QMARK__226 = entry_QMARK__177
		head__9859_227 = head__9859_178
		arg__9860_228 = arg__9860_179
		arg__9862_229 = "b"
		arg__9864_230 = "("
		arg__10020_231 = arg__10020_201
		arg__10021_232 = "):"
		goto b11
	}
b10:
	;
	arg__10042_241, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
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
	arg__10060_247, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
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
	arg__10061_248, callErr = rt.InvokeValue(rt.LookupVar("clojure.string", "join").Deref(), []vm.Value{vm.String(", "), arg__10060_247})
	if callErr != nil {
		return nil, callErr
	}
	arg__10080_255, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
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
	arg__10098_261, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
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
	arg__10099_262, callErr = rt.InvokeValue(rt.LookupVar("clojure.string", "join").Deref(), []vm.Value{vm.String(", "), arg__10098_261})
	if callErr != nil {
		return nil, callErr
	}
	v263, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("    ; preds: "), arg__10099_262})
	if callErr != nil {
		return nil, callErr
	}
	arg__10100_267 = v263
	arg__9861_268 = arg__9861_203
	f_269 = f_204
	arg__9863_270 = arg__9863_205
	bid_271 = bid_206
	params_272 = params_207
	preds_273 = preds_208
	insts_274 = insts_209
	term_275 = term_210
	entry_QMARK__276 = entry_QMARK__211
	head__9859_277 = head__9859_212
	arg__9860_278 = arg__9860_213
	arg__9862_279 = arg__9862_214
	arg__9864_280 = arg__9864_215
	arg__10020_281 = arg__10020_216
	arg__10021_282 = arg__10021_217
	goto b12
b11:
	;
	arg__10100_267 = vm.String("")
	arg__9861_268 = arg__9861_218
	f_269 = f_219
	arg__9863_270 = arg__9863_220
	bid_271 = bid_221
	params_272 = params_222
	preds_273 = preds_223
	insts_274 = insts_224
	term_275 = term_225
	entry_QMARK__276 = entry_QMARK__226
	head__9859_277 = head__9859_227
	arg__9860_278 = arg__9860_228
	arg__9862_279 = arg__9862_229
	arg__9864_280 = arg__9864_230
	arg__10020_281 = arg__10020_231
	arg__10021_282 = arg__10021_232
	goto b12
b12:
	;
	header_284, callErr = rt.InvokeValue(head__9859_277, []vm.Value{vm.String(arg__9860_278), vm.String(arg__9861_268), vm.String(arg__9862_279), arg__9863_270, vm.String(arg__9864_280), arg__10020_281, vm.String(arg__10021_282), arg__10100_267, vm.String("\n")})
	if callErr != nil {
		return nil, callErr
	}
	arg__10118_293, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
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
	arg__10136_303, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
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
	body_304, callErr = rt.InvokeValue(rt.LookupVar("clojure.string", "join").Deref(), []vm.Value{vm.String(""), arg__10136_303})
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
	var arg__10163_4 vm.Value
	var arg__10168_7 vm.Value
	var arg__10173_10 vm.Value
	var arg__10185_17 vm.Value
	var arg__10196_23 vm.Value
	var arg__10197_24 vm.Value
	var arg__10209_31 vm.Value
	var arg__10220_37 vm.Value
	var arg__10221_38 vm.Value
	var arg__10222_39 vm.Value
	var arg__10228_43 vm.Value
	var arg__10233_46 vm.Value
	var arg__10238_49 vm.Value
	var arg__10250_56 vm.Value
	var arg__10261_62 vm.Value
	var arg__10262_63 vm.Value
	var arg__10274_70 vm.Value
	var arg__10285_76 vm.Value
	var arg__10286_77 vm.Value
	var arg__10287_78 vm.Value
	var v79 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__10163_4, arg__10168_7, arg__10173_10, arg__10185_17, arg__10196_23, arg__10197_24, arg__10209_31, arg__10220_37, arg__10221_38, arg__10222_39, arg__10228_43, arg__10233_46, arg__10238_49, arg__10250_56, arg__10261_62, arg__10262_63, arg__10274_70, arg__10285_76, arg__10286_77, arg__10287_78, v79
	arg__10163_4, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-name").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__10168_7, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-arity").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__10173_10, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-variadic?").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__10185_17, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__10196_23, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__10197_24, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "write-block").Deref(), []vm.Value{arg0, arg0})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), arg__10196_23})
	if callErr != nil {
		return nil, callErr
	}
	arg__10209_31, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__10220_37, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__10221_38, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "write-block").Deref(), []vm.Value{arg0, arg0})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), arg__10220_37})
	if callErr != nil {
		return nil, callErr
	}
	arg__10222_39, callErr = rt.InvokeValue(rt.LookupVar("clojure.string", "join").Deref(), []vm.Value{vm.String(""), arg__10221_38})
	if callErr != nil {
		return nil, callErr
	}
	arg__10228_43, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-name").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__10233_46, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-arity").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__10238_49, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-variadic?").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__10250_56, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__10261_62, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__10262_63, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "write-block").Deref(), []vm.Value{arg0, arg0})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), arg__10261_62})
	if callErr != nil {
		return nil, callErr
	}
	arg__10274_70, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__10285_76, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__10286_77, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("ir.dump", "write-block").Deref(), []vm.Value{arg0, arg0})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), arg__10285_76})
	if callErr != nil {
		return nil, callErr
	}
	arg__10287_78, callErr = rt.InvokeValue(rt.LookupVar("clojure.string", "join").Deref(), []vm.Value{vm.String(""), arg__10286_77})
	if callErr != nil {
		return nil, callErr
	}
	v79, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("fn "), arg__10228_43, vm.String("(arity="), arg__10233_46, vm.String(", variadic="), arg__10238_49, vm.String("):\n"), arg__10287_78})
	if callErr != nil {
		return nil, callErr
	}
	return v79, nil
}
