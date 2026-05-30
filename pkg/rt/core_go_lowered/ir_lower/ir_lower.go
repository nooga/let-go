package ir_lower

import (
	rt "github.com/nooga/let-go/pkg/rt"
	vm "github.com/nooga/let-go/pkg/vm"
)

func add_patch_BANG_(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var v9 vm.Value
	var callErr error
	_ = v9
	v9, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "update").Deref(), vm.Keyword("patches"), rt.LookupVar("clojure.core", "conj").Deref(), arg1})
	if callErr != nil {
		return nil, callErr
	}
	return v9, nil
}
func args_at_top_QMARK_(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var n_4 vm.Value
	var v12 vm.Value
	var l_5 vm.Value
	var args_6 vm.Value
	var n_7 vm.Value
	var l_8 vm.Value
	var args_9 vm.Value
	var n_10 vm.Value
	var v161 vm.Value
	var l_162 vm.Value
	var args_163 vm.Value
	var n_164 vm.Value
	var l_16 vm.Value
	var args_17 vm.Value
	var n_18 vm.Value
	var arg__10308_25 vm.Value
	var base_26 vm.Value
	var v36 vm.Value
	var l_19 vm.Value
	var args_20 vm.Value
	var n_21 vm.Value
	var v156 vm.Value
	var l_157 vm.Value
	var args_158 vm.Value
	var n_159 vm.Value
	var l_27 vm.Value
	var args_28 vm.Value
	var n_29 vm.Value
	var base_30 vm.Value
	var l_31 vm.Value
	var args_32 vm.Value
	var n_33 vm.Value
	var base_34 vm.Value
	var v148 vm.Value
	var l_149 vm.Value
	var args_150 vm.Value
	var n_151 vm.Value
	var base_152 vm.Value
	var i_40 int
	var n_41 vm.Value
	var base_42 vm.Value
	var args_43 vm.Value
	var l_44 vm.Value
	var v169 vm.Value
	var v57 bool
	var i_47 int
	var n_48 vm.Value
	var base_49 vm.Value
	var args_50 vm.Value
	var l_51 vm.Value
	var v175 vm.Value
	var i_52 int
	var n_53 vm.Value
	var base_54 vm.Value
	var args_55 vm.Value
	var l_56 vm.Value
	var v171 vm.Value
	var arg__10321_72 vm.Value
	var arg__10329_75 vm.Value
	var pos_76 vm.Value
	var or__x_78 vm.Value
	var v141 vm.Value
	var i_142 int
	var n_143 vm.Value
	var base_144 vm.Value
	var args_145 vm.Value
	var l_146 vm.Value
	var i_61 int
	var n_62 vm.Value
	var base_63 vm.Value
	var args_64 vm.Value
	var l_65 vm.Value
	var v176 vm.Value
	var i_66 int
	var n_67 vm.Value
	var base_68 vm.Value
	var args_69 vm.Value
	var l_70 vm.Value
	var v168 vm.Value
	var v134 vm.Value
	var i_135 int
	var n_136 vm.Value
	var base_137 vm.Value
	var args_138 vm.Value
	var l_139 vm.Value
	var i_79 int
	var n_80 vm.Value
	var base_81 vm.Value
	var args_82 vm.Value
	var l_83 vm.Value
	var pos_84 vm.Value
	var or__x_85 vm.Value
	var v172 vm.Value
	var i_86 int
	var n_87 vm.Value
	var base_88 vm.Value
	var args_89 vm.Value
	var l_90 vm.Value
	var pos_91 vm.Value
	var or__x_92 vm.Value
	var v170 vm.Value
	var arg__10336_95 vm.Value
	var v98 vm.Value
	var v100 vm.Value
	var i_101 int
	var n_102 vm.Value
	var base_103 vm.Value
	var args_104 vm.Value
	var l_105 vm.Value
	var pos_106 vm.Value
	var or__x_107 vm.Value
	var v173 vm.Value
	var i_111 int
	var n_112 vm.Value
	var base_113 vm.Value
	var args_114 vm.Value
	var l_115 vm.Value
	var v167 vm.Value
	var v123 int
	var i_116 int
	var n_117 vm.Value
	var base_118 vm.Value
	var args_119 vm.Value
	var l_120 vm.Value
	var v174 vm.Value
	var v127 vm.Value
	var i_128 int
	var n_129 vm.Value
	var base_130 vm.Value
	var args_131 vm.Value
	var l_132 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = n_4, v12, l_5, args_6, n_7, l_8, args_9, n_10, v161, l_162, args_163, n_164, l_16, args_17, n_18, arg__10308_25, base_26, v36, l_19, args_20, n_21, v156, l_157, args_158, n_159, l_27, args_28, n_29, base_30, l_31, args_32, n_33, base_34, v148, l_149, args_150, n_151, base_152, i_40, n_41, base_42, args_43, l_44, v169, v57, i_47, n_48, base_49, args_50, l_51, v175, i_52, n_53, base_54, args_55, l_56, v171, arg__10321_72, arg__10329_75, pos_76, or__x_78, v141, i_142, n_143, base_144, args_145, l_146, i_61, n_62, base_63, args_64, l_65, v176, i_66, n_67, base_68, args_69, l_70, v168, v134, i_135, n_136, base_137, args_138, l_139, i_79, n_80, base_81, args_82, l_83, pos_84, or__x_85, v172, i_86, n_87, base_88, args_89, l_90, pos_91, or__x_92, v170, arg__10336_95, v98, v100, i_101, n_102, base_103, args_104, l_105, pos_106, or__x_107, v173, i_111, n_112, base_113, args_114, l_115, v167, v123, i_116, n_117, base_118, args_119, l_120, v174, v127, i_128, n_129, base_130, args_131, l_132
	n_4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	v12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "zero?").Deref(), []vm.Value{n_4})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v12) {
		l_5 = arg0
		args_6 = arg1
		n_7 = n_4
		goto b1
	} else {
		l_8 = arg0
		args_9 = arg1
		n_10 = n_4
		goto b2
	}
b1:
	;
	v161 = vm.Boolean(true)
	l_162 = l_5
	args_163 = args_6
	n_164 = n_7
	goto b3
b2:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		l_16 = l_8
		args_17 = args_9
		n_18 = n_10
		goto b4
	} else {
		l_19 = l_8
		args_20 = args_9
		n_21 = n_10
		goto b5
	}
b3:
	;
	return v161, nil
b4:
	;
	arg__10308_25, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "stack-sp").Deref(), []vm.Value{l_16})
	if callErr != nil {
		return nil, callErr
	}
	base_26 = rt.SubValue(arg__10308_25, n_18)
	v36, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "neg?").Deref(), []vm.Value{base_26})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v36) {
		l_27 = l_16
		args_28 = args_17
		n_29 = n_18
		base_30 = base_26
		goto b7
	} else {
		l_31 = l_16
		args_32 = args_17
		n_33 = n_18
		base_34 = base_26
		goto b8
	}
b5:
	;
	v156 = vm.NIL
	l_157 = l_19
	args_158 = args_20
	n_159 = n_21
	goto b6
b6:
	;
	v161 = v156
	l_162 = l_157
	args_163 = args_158
	n_164 = n_159
	goto b3
b7:
	;
	v148 = vm.Boolean(false)
	l_149 = l_27
	args_150 = args_28
	n_151 = n_29
	base_152 = base_30
	goto b9
b8:
	;
	i_40 = 0
	n_41 = n_33
	base_42 = base_34
	args_43 = args_32
	l_44 = l_31
	v169 = vm.Keyword("else")
	goto b10
b9:
	;
	v156 = v148
	l_157 = l_149
	args_158 = args_150
	n_159 = n_151
	goto b6
b10:
	;
	v57 = rt.GeValue(vm.Int(i_40), n_41)
	if v57 {
		i_47 = i_40
		n_48 = n_41
		base_49 = base_42
		args_50 = args_43
		l_51 = l_44
		v175 = v169
		goto b11
	} else {
		i_52 = i_40
		n_53 = n_41
		base_54 = base_42
		args_55 = args_43
		l_56 = l_44
		v171 = v169
		goto b12
	}
b11:
	;
	v141 = vm.Boolean(true)
	i_142 = i_47
	n_143 = n_48
	base_144 = base_49
	args_145 = args_50
	l_146 = l_51
	goto b13
b12:
	;
	arg__10321_72, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{args_55, vm.Int(i_52)})
	if callErr != nil {
		return nil, callErr
	}
	arg__10329_75, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{args_55, vm.Int(i_52)})
	if callErr != nil {
		return nil, callErr
	}
	pos_76, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "value-pos-of").Deref(), []vm.Value{l_56, arg__10329_75})
	if callErr != nil {
		return nil, callErr
	}
	or__x_78, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{pos_76})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(or__x_78) {
		i_79 = i_52
		n_80 = n_53
		base_81 = base_54
		args_82 = args_55
		l_83 = l_56
		pos_84 = pos_76
		or__x_85 = or__x_78
		v172 = v171
		goto b17
	} else {
		i_86 = i_52
		n_87 = n_53
		base_88 = base_54
		args_89 = args_55
		l_90 = l_56
		pos_91 = pos_76
		or__x_92 = or__x_78
		v170 = v171
		goto b18
	}
b13:
	;
	v148 = v141
	l_149 = l_146
	args_150 = args_145
	n_151 = n_143
	base_152 = base_144
	goto b9
b14:
	;
	v134 = vm.Boolean(false)
	i_135 = i_61
	n_136 = n_62
	base_137 = base_63
	args_138 = args_64
	l_139 = l_65
	goto b16
b15:
	;
	if vm.IsTruthy(v168) {
		i_111 = i_66
		n_112 = n_67
		base_113 = base_68
		args_114 = args_69
		l_115 = l_70
		v167 = v168
		goto b20
	} else {
		i_116 = i_66
		n_117 = n_67
		base_118 = base_68
		args_119 = args_69
		l_120 = l_70
		v174 = v168
		goto b21
	}
b16:
	;
	v141 = v134
	i_142 = i_135
	n_143 = n_136
	base_144 = base_137
	args_145 = args_138
	l_146 = l_139
	goto b13
b17:
	;
	v100 = or__x_85
	i_101 = i_79
	n_102 = n_80
	base_103 = base_81
	args_104 = args_82
	l_105 = l_83
	pos_106 = pos_84
	or__x_107 = or__x_85
	v173 = v172
	goto b19
b18:
	;
	arg__10336_95 = rt.AddValue(base_88, vm.Int(i_86))
	v98, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not=").Deref(), []vm.Value{pos_91, arg__10336_95})
	if callErr != nil {
		return nil, callErr
	}
	v100 = v98
	i_101 = i_86
	n_102 = n_87
	base_103 = base_88
	args_104 = args_89
	l_105 = l_90
	pos_106 = pos_91
	or__x_107 = or__x_92
	v173 = v170
	goto b19
b19:
	;
	if vm.IsTruthy(v100) {
		i_61 = i_101
		n_62 = n_102
		base_63 = base_103
		args_64 = args_104
		l_65 = l_105
		v176 = v173
		goto b14
	} else {
		i_66 = i_101
		n_67 = n_102
		base_68 = base_103
		args_69 = args_104
		l_70 = l_105
		v168 = v173
		goto b15
	}
b20:
	;
	v123 = i_111 + 1
	i_40 = v123
	n_41 = n_112
	base_42 = base_113
	args_43 = args_114
	l_44 = l_115
	v169 = v167
	goto b10
b21:
	;
	v127 = vm.NIL
	i_128 = i_116
	n_129 = n_117
	base_130 = base_118
	args_131 = args_119
	l_132 = l_120
	goto b22
b22:
	;
	v134 = v127
	i_135 = i_128
	n_136 = n_129
	base_137 = base_130
	args_138 = args_131
	l_139 = l_132
	goto b16
}
func block_ip(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__10347_4 vm.Value
	var arg__10348_5 vm.Value
	var arg__10355_9 vm.Value
	var arg__10356_10 vm.Value
	var v11 vm.Value
	var callErr error
	_, _, _, _, _ = arg__10347_4, arg__10348_5, arg__10355_9, arg__10356_10, v11
	arg__10347_4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__10348_5, callErr = rt.InvokeValue(vm.Keyword("block-ips"), []vm.Value{arg__10347_4})
	if callErr != nil {
		return nil, callErr
	}
	arg__10355_9, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__10356_10, callErr = rt.InvokeValue(vm.Keyword("block-ips"), []vm.Value{arg__10355_9})
	if callErr != nil {
		return nil, callErr
	}
	v11, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__10356_10, arg1})
	if callErr != nil {
		return nil, callErr
	}
	return v11, nil
}
func block_junk_of(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__10362_4 vm.Value
	var arg__10363_5 vm.Value
	var arg__10371_10 vm.Value
	var arg__10372_11 vm.Value
	var v13 vm.Value
	var callErr error
	_, _, _, _, _ = arg__10362_4, arg__10363_5, arg__10371_10, arg__10372_11, v13
	arg__10362_4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__10363_5, callErr = rt.InvokeValue(vm.Keyword("block-junk"), []vm.Value{arg__10362_4})
	if callErr != nil {
		return nil, callErr
	}
	arg__10371_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__10372_11, callErr = rt.InvokeValue(vm.Keyword("block-junk"), []vm.Value{arg__10371_10})
	if callErr != nil {
		return nil, callErr
	}
	v13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg__10372_11, arg1, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	return v13, nil
}
func bump_max_stack_BANG_(arg0 vm.Value) (vm.Value, error) {
	var c_3 vm.Value
	var arg__10381_5 vm.Value
	var arg__10387_8 vm.Value
	var arg__10388_9 vm.Value
	var arg__10395_13 vm.Value
	var arg__10396_14 vm.Value
	var arg__10397_15 vm.Value
	var rt_sp_16 vm.Value
	var arg__10402_24 vm.Value
	var v25 bool
	var l_17 vm.Value
	var c_18 vm.Value
	var rt_sp_19 vm.Value
	var v28 vm.Value
	var l_20 vm.Value
	var c_21 vm.Value
	var rt_sp_22 vm.Value
	var v32 vm.Value
	var l_33 vm.Value
	var c_34 vm.Value
	var rt_sp_35 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = c_3, arg__10381_5, arg__10387_8, arg__10388_9, arg__10395_13, arg__10396_14, arg__10397_15, rt_sp_16, arg__10402_24, v25, l_17, c_18, rt_sp_19, v28, l_20, c_21, rt_sp_22, v32, l_33, c_34, rt_sp_35
	c_3, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "chunk-of").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__10381_5, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "stack-sp").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__10387_8, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__10388_9, callErr = rt.InvokeValue(vm.Keyword("current-block"), []vm.Value{arg__10387_8})
	if callErr != nil {
		return nil, callErr
	}
	arg__10395_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__10396_14, callErr = rt.InvokeValue(vm.Keyword("current-block"), []vm.Value{arg__10395_13})
	if callErr != nil {
		return nil, callErr
	}
	arg__10397_15, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "block-junk-of").Deref(), []vm.Value{arg0, arg__10396_14})
	if callErr != nil {
		return nil, callErr
	}
	rt_sp_16 = rt.AddValue(arg__10381_5, arg__10397_15)
	arg__10402_24, callErr = rt.InvokeValue(rt.LookupVar("ir", "chunk-max-stack").Deref(), []vm.Value{c_3})
	if callErr != nil {
		return nil, callErr
	}
	v25 = rt.GtValue(rt_sp_16, arg__10402_24)
	if v25 {
		l_17 = arg0
		c_18 = c_3
		rt_sp_19 = rt_sp_16
		goto b1
	} else {
		l_20 = arg0
		c_21 = c_3
		rt_sp_22 = rt_sp_16
		goto b2
	}
b1:
	;
	v28, callErr = rt.InvokeValue(rt.LookupVar("ir", "chunk-set-max-stack!").Deref(), []vm.Value{c_18, rt_sp_19})
	if callErr != nil {
		return nil, callErr
	}
	v32 = v28
	l_33 = l_17
	c_34 = c_18
	rt_sp_35 = rt_sp_19
	goto b3
b2:
	;
	v32 = vm.NIL
	l_33 = l_20
	c_34 = c_21
	rt_sp_35 = rt_sp_22
	goto b3
b3:
	;
	return v32, nil
}
func bump_stack_sp_BANG_(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var v9 vm.Value
	var callErr error
	_ = v9
	v9, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "update").Deref(), vm.Keyword("stack-sp"), rt.LookupVar("clojure.core", "+").Deref(), arg1})
	if callErr != nil {
		return nil, callErr
	}
	return v9, nil
}
func chunk_of(arg0 vm.Value) (vm.Value, error) {
	var arg__10575_3 vm.Value
	var v4 vm.Value
	var callErr error
	_, _ = arg__10575_3, v4
	arg__10575_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v4, callErr = rt.InvokeValue(vm.Keyword("chunk"), []vm.Value{arg__10575_3})
	if callErr != nil {
		return nil, callErr
	}
	return v4, nil
}
func consume_refs_in_place_BANG_(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var doseq_seq__10576_4 vm.Value
	var doseq_loop__10577_5 vm.Value
	var l_6 vm.Value
	var refs_8 vm.Value
	var doseq_seq__10576_9 vm.Value
	var doseq_loop__10577_10 vm.Value
	var l_11 vm.Value
	var r_18 vm.Value
	var v20 vm.Value
	var v22 vm.Value
	var refs_12 vm.Value
	var doseq_seq__10576_13 vm.Value
	var doseq_loop__10577_14 vm.Value
	var l_15 vm.Value
	var v26 vm.Value
	var refs_27 vm.Value
	var doseq_seq__10576_28 vm.Value
	var doseq_loop__10577_29 vm.Value
	var l_30 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = doseq_seq__10576_4, doseq_loop__10577_5, l_6, refs_8, doseq_seq__10576_9, doseq_loop__10577_10, l_11, r_18, v20, v22, refs_12, doseq_seq__10576_13, doseq_loop__10577_14, l_15, v26, refs_27, doseq_seq__10576_28, doseq_loop__10577_29, l_30
	doseq_seq__10576_4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__10577_5 = doseq_seq__10576_4
	l_6 = arg0
	goto b1
b1:
	;
	if vm.IsTruthy(doseq_loop__10577_5) {
		refs_8 = arg1
		doseq_seq__10576_9 = doseq_seq__10576_4
		doseq_loop__10577_10 = doseq_loop__10577_5
		l_11 = l_6
		goto b2
	} else {
		refs_12 = arg1
		doseq_seq__10576_13 = doseq_seq__10576_4
		doseq_loop__10577_14 = doseq_loop__10577_5
		l_15 = l_6
		goto b3
	}
b2:
	;
	r_18, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{doseq_loop__10577_10})
	if callErr != nil {
		return nil, callErr
	}
	v20, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "decrement-use!").Deref(), []vm.Value{l_11, r_18})
	if callErr != nil {
		return nil, callErr
	}
	v22, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "next").Deref(), []vm.Value{doseq_loop__10577_10})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__10577_5 = v22
	l_6 = l_11
	goto b1
b3:
	;
	v26 = vm.NIL
	refs_27 = refs_12
	doseq_seq__10576_28 = doseq_seq__10576_13
	doseq_loop__10577_29 = doseq_loop__10577_14
	l_30 = l_15
	goto b4
b4:
	;
	return v26, nil
}
func contains_after_QMARK_(arg0 vm.Value, arg1 int, arg2 vm.Value) (vm.Value, error) {
	var v7 int
	var j_4 int
	var refs_5 vm.Value
	var target_6 vm.Value
	var v71 vm.Value
	var arg__10597_18 vm.Value
	var v19 bool
	var i_9 int
	var j_10 int
	var refs_11 vm.Value
	var target_12 vm.Value
	var v74 vm.Value
	var i_13 int
	var j_14 int
	var refs_15 vm.Value
	var target_16 vm.Value
	var v72 vm.Value
	var arg__10603_32 vm.Value
	var v33 bool
	var v63 vm.Value
	var i_64 int
	var j_65 int
	var refs_66 vm.Value
	var target_67 vm.Value
	var i_23 int
	var j_24 int
	var refs_25 vm.Value
	var target_26 vm.Value
	var v75 vm.Value
	var i_27 int
	var j_28 int
	var refs_29 vm.Value
	var target_30 vm.Value
	var v73 vm.Value
	var v57 vm.Value
	var i_58 int
	var j_59 int
	var refs_60 vm.Value
	var target_61 vm.Value
	var i_37 int
	var j_38 int
	var refs_39 vm.Value
	var target_40 vm.Value
	var v70 vm.Value
	var v47 int
	var i_41 int
	var j_42 int
	var refs_43 vm.Value
	var target_44 vm.Value
	var v76 vm.Value
	var v51 vm.Value
	var i_52 int
	var j_53 int
	var refs_54 vm.Value
	var target_55 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = v7, j_4, refs_5, target_6, v71, arg__10597_18, v19, i_9, j_10, refs_11, target_12, v74, i_13, j_14, refs_15, target_16, v72, arg__10603_32, v33, v63, i_64, j_65, refs_66, target_67, i_23, j_24, refs_25, target_26, v75, i_27, j_28, refs_29, target_30, v73, v57, i_58, j_59, refs_60, target_61, i_37, j_38, refs_39, target_40, v70, v47, i_41, j_42, refs_43, target_44, v76, v51, i_52, j_53, refs_54, target_55
	v7 = arg1 + 1
	j_4 = v7
	refs_5 = arg0
	target_6 = arg2
	v71 = vm.Keyword("else")
	goto b1
b1:
	;
	arg__10597_18, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{refs_5})
	if callErr != nil {
		return nil, callErr
	}
	v19 = rt.GeValue(vm.Int(j_4), arg__10597_18)
	if v19 {
		i_9 = arg1
		j_10 = j_4
		refs_11 = refs_5
		target_12 = target_6
		v74 = v71
		goto b2
	} else {
		i_13 = arg1
		j_14 = j_4
		refs_15 = refs_5
		target_16 = target_6
		v72 = v71
		goto b3
	}
b2:
	;
	v63 = vm.Boolean(false)
	i_64 = i_9
	j_65 = j_10
	refs_66 = refs_11
	target_67 = target_12
	goto b4
b3:
	;
	arg__10603_32, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{refs_15, vm.Int(j_14)})
	if callErr != nil {
		return nil, callErr
	}
	v33 = arg__10603_32 == target_16
	if v33 {
		i_23 = i_13
		j_24 = j_14
		refs_25 = refs_15
		target_26 = target_16
		v75 = v72
		goto b5
	} else {
		i_27 = i_13
		j_28 = j_14
		refs_29 = refs_15
		target_30 = target_16
		v73 = v72
		goto b6
	}
b4:
	;
	return v63, nil
b5:
	;
	v57 = vm.Boolean(true)
	i_58 = i_23
	j_59 = j_24
	refs_60 = refs_25
	target_61 = target_26
	goto b7
b6:
	;
	if vm.IsTruthy(v73) {
		i_37 = i_27
		j_38 = j_28
		refs_39 = refs_29
		target_40 = target_30
		v70 = v73
		goto b8
	} else {
		i_41 = i_27
		j_42 = j_28
		refs_43 = refs_29
		target_44 = target_30
		v76 = v73
		goto b9
	}
b7:
	;
	v63 = v57
	i_64 = i_58
	j_65 = j_59
	refs_66 = refs_60
	target_67 = target_61
	goto b4
b8:
	;
	v47 = j_38 + 1
	j_4 = v47
	refs_5 = refs_39
	target_6 = target_40
	v71 = v70
	goto b1
b9:
	;
	v51 = vm.NIL
	i_52 = i_41
	j_53 = j_42
	refs_54 = refs_43
	target_55 = target_44
	goto b10
b10:
	;
	v57 = v51
	i_58 = i_52
	j_59 = j_53
	refs_60 = refs_54
	target_61 = target_55
	goto b7
}
func decrement_use_BANG_(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var cur_4 vm.Value
	var v12 vm.Value
	var l_5 vm.Value
	var nid_6 vm.Value
	var cur_7 vm.Value
	var arg__10620_17 vm.Value
	var v23 vm.Value
	var l_8 vm.Value
	var nid_9 vm.Value
	var cur_10 vm.Value
	var v27 vm.Value
	var l_28 vm.Value
	var nid_29 vm.Value
	var cur_30 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _ = cur_4, v12, l_5, nid_6, cur_7, arg__10620_17, v23, l_8, nid_9, cur_10, v27, l_28, nid_29, cur_30
	cur_4, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "use-count-of").Deref(), []vm.Value{arg0, arg1})
	if callErr != nil {
		return nil, callErr
	}
	v12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "pos?").Deref(), []vm.Value{cur_4})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v12) {
		l_5 = arg0
		nid_6 = arg1
		cur_7 = cur_4
		goto b1
	} else {
		l_8 = arg0
		nid_9 = arg1
		cur_10 = cur_4
		goto b2
	}
b1:
	;
	arg__10620_17 = rt.SubValue(cur_7, vm.Int(1))
	v23, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{l_5, rt.LookupVar("clojure.core", "update").Deref(), vm.Keyword("use-count"), rt.LookupVar("clojure.core", "assoc").Deref(), nid_6, arg__10620_17})
	if callErr != nil {
		return nil, callErr
	}
	v27 = v23
	l_28 = l_5
	nid_29 = nid_6
	cur_30 = cur_7
	goto b3
b2:
	;
	v27 = vm.NIL
	l_28 = l_8
	nid_29 = nid_9
	cur_30 = cur_10
	goto b3
b3:
	;
	return v27, nil
}
func deferrable_branch_if_cond_QMARK_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var l_4 vm.Value
	var term_5 vm.Value
	var cond_ref_6 vm.Value
	var arg__10635_13 vm.Value
	var and__x_14 bool
	var l_7 vm.Value
	var term_8 vm.Value
	var cond_ref_9 vm.Value
	var v116 vm.Value
	var l_117 vm.Value
	var term_118 vm.Value
	var cond_ref_119 vm.Value
	var l_15 vm.Value
	var term_16 vm.Value
	var cond_ref_17 vm.Value
	var and__x_18 bool
	var arg__10640_26 vm.Value
	var uses_27 vm.Value
	var arg__10645_39 vm.Value
	var v40 bool
	var l_19 vm.Value
	var term_20 vm.Value
	var cond_ref_21 vm.Value
	var and__x_22 bool
	var v108 vm.Value
	var l_109 vm.Value
	var term_110 vm.Value
	var cond_ref_111 vm.Value
	var and__x_112 vm.Value
	var l_28 vm.Value
	var term_29 vm.Value
	var cond_ref_30 vm.Value
	var and__x_31 bool
	var uses_32 vm.Value
	var v43 vm.Value
	var l_33 vm.Value
	var term_34 vm.Value
	var cond_ref_35 vm.Value
	var and__x_36 bool
	var uses_37 vm.Value
	var us_47 vm.Value
	var l_48 vm.Value
	var term_49 vm.Value
	var cond_ref_50 vm.Value
	var and__x_51 bool
	var uses_52 vm.Value
	var and__x_53 vm.Value
	var us_54 vm.Value
	var l_55 vm.Value
	var term_56 vm.Value
	var cond_ref_57 vm.Value
	var uses_58 vm.Value
	var arg__10654_67 vm.Value
	var arg__10659_70 vm.Value
	var and__x_71 vm.Value
	var and__x_59 vm.Value
	var us_60 vm.Value
	var l_61 vm.Value
	var term_62 vm.Value
	var cond_ref_63 vm.Value
	var uses_64 vm.Value
	var v99 vm.Value
	var and__x_100 vm.Value
	var us_101 vm.Value
	var l_102 vm.Value
	var term_103 vm.Value
	var cond_ref_104 vm.Value
	var uses_105 vm.Value
	var us_72 vm.Value
	var l_73 vm.Value
	var term_74 vm.Value
	var cond_ref_75 vm.Value
	var uses_76 vm.Value
	var and__x_77 vm.Value
	var arg__10664_86 vm.Value
	var v87 bool
	var us_78 vm.Value
	var l_79 vm.Value
	var term_80 vm.Value
	var cond_ref_81 vm.Value
	var uses_82 vm.Value
	var and__x_83 vm.Value
	var v90 vm.Value
	var us_91 vm.Value
	var l_92 vm.Value
	var term_93 vm.Value
	var cond_ref_94 vm.Value
	var uses_95 vm.Value
	var and__x_96 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = l_4, term_5, cond_ref_6, arg__10635_13, and__x_14, l_7, term_8, cond_ref_9, v116, l_117, term_118, cond_ref_119, l_15, term_16, cond_ref_17, and__x_18, arg__10640_26, uses_27, arg__10645_39, v40, l_19, term_20, cond_ref_21, and__x_22, v108, l_109, term_110, cond_ref_111, and__x_112, l_28, term_29, cond_ref_30, and__x_31, uses_32, v43, l_33, term_34, cond_ref_35, and__x_36, uses_37, us_47, l_48, term_49, cond_ref_50, and__x_51, uses_52, and__x_53, us_54, l_55, term_56, cond_ref_57, uses_58, arg__10654_67, arg__10659_70, and__x_71, and__x_59, us_60, l_61, term_62, cond_ref_63, uses_64, v99, and__x_100, us_101, l_102, term_103, cond_ref_104, uses_105, us_72, l_73, term_74, cond_ref_75, uses_76, and__x_77, arg__10664_86, v87, us_78, l_79, term_80, cond_ref_81, uses_82, and__x_83, v90, us_91, l_92, term_93, cond_ref_94, uses_95, and__x_96
	if vm.IsTruthy(arg2) {
		l_4 = arg0
		term_5 = arg1
		cond_ref_6 = arg2
		goto b1
	} else {
		l_7 = arg0
		term_8 = arg1
		cond_ref_9 = arg2
		goto b2
	}
b1:
	;
	arg__10635_13, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "use-count-of").Deref(), []vm.Value{l_4, cond_ref_6})
	if callErr != nil {
		return nil, callErr
	}
	and__x_14 = arg__10635_13 == vm.Int(1)
	if and__x_14 {
		l_15 = l_4
		term_16 = term_5
		cond_ref_17 = cond_ref_6
		and__x_18 = and__x_14
		goto b4
	} else {
		l_19 = l_4
		term_20 = term_5
		cond_ref_21 = cond_ref_6
		and__x_22 = and__x_14
		goto b5
	}
b2:
	;
	v116 = vm.NIL
	l_117 = l_7
	term_118 = term_8
	cond_ref_119 = cond_ref_9
	goto b3
b3:
	;
	return v116, nil
b4:
	;
	arg__10640_26, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{l_15})
	if callErr != nil {
		return nil, callErr
	}
	uses_27, callErr = rt.InvokeValue(vm.Keyword("uses"), []vm.Value{arg__10640_26})
	if callErr != nil {
		return nil, callErr
	}
	arg__10645_39, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{uses_27})
	if callErr != nil {
		return nil, callErr
	}
	v40 = rt.LtValue(cond_ref_17, arg__10645_39)
	if v40 {
		l_28 = l_15
		term_29 = term_16
		cond_ref_30 = cond_ref_17
		and__x_31 = and__x_18
		uses_32 = uses_27
		goto b7
	} else {
		l_33 = l_15
		term_34 = term_16
		cond_ref_35 = cond_ref_17
		and__x_36 = and__x_18
		uses_37 = uses_27
		goto b8
	}
b5:
	;
	v108 = vm.Boolean(and__x_22)
	l_109 = l_19
	term_110 = term_20
	cond_ref_111 = cond_ref_21
	and__x_112 = vm.Boolean(and__x_22)
	goto b6
b6:
	;
	v116 = v108
	l_117 = l_109
	term_118 = term_110
	cond_ref_119 = cond_ref_111
	goto b3
b7:
	;
	v43, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{uses_32, cond_ref_30})
	if callErr != nil {
		return nil, callErr
	}
	us_47 = v43
	l_48 = l_28
	term_49 = term_29
	cond_ref_50 = cond_ref_30
	and__x_51 = and__x_31
	uses_52 = uses_32
	goto b9
b8:
	;
	us_47 = vm.NIL
	l_48 = l_33
	term_49 = term_34
	cond_ref_50 = cond_ref_35
	and__x_51 = and__x_36
	uses_52 = uses_37
	goto b9
b9:
	;
	if vm.IsTruthy(us_47) {
		and__x_53 = us_47
		us_54 = us_47
		l_55 = l_48
		term_56 = term_49
		cond_ref_57 = cond_ref_50
		uses_58 = uses_52
		goto b10
	} else {
		and__x_59 = us_47
		us_60 = us_47
		l_61 = l_48
		term_62 = term_49
		cond_ref_63 = cond_ref_50
		uses_64 = uses_52
		goto b11
	}
b10:
	;
	arg__10654_67, callErr = rt.InvokeValue(rt.LookupVar("ir", "uses-empty?").Deref(), []vm.Value{us_54})
	if callErr != nil {
		return nil, callErr
	}
	arg__10659_70, callErr = rt.InvokeValue(rt.LookupVar("ir", "uses-empty?").Deref(), []vm.Value{us_54})
	if callErr != nil {
		return nil, callErr
	}
	and__x_71, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not").Deref(), []vm.Value{arg__10659_70})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(and__x_71) {
		us_72 = us_54
		l_73 = l_55
		term_74 = term_56
		cond_ref_75 = cond_ref_57
		uses_76 = uses_58
		and__x_77 = and__x_71
		goto b13
	} else {
		us_78 = us_54
		l_79 = l_55
		term_80 = term_56
		cond_ref_81 = cond_ref_57
		uses_82 = uses_58
		and__x_83 = and__x_71
		goto b14
	}
b11:
	;
	v99 = and__x_59
	and__x_100 = and__x_59
	us_101 = us_60
	l_102 = l_61
	term_103 = term_62
	cond_ref_104 = cond_ref_63
	uses_105 = uses_64
	goto b12
b12:
	;
	v108 = v99
	l_109 = l_102
	term_110 = term_103
	cond_ref_111 = cond_ref_104
	and__x_112 = vm.Boolean(and__x_51)
	goto b6
b13:
	;
	arg__10664_86, callErr = rt.InvokeValue(rt.LookupVar("ir", "uses-first").Deref(), []vm.Value{us_72})
	if callErr != nil {
		return nil, callErr
	}
	v87 = term_74 == arg__10664_86
	v90 = vm.Boolean(v87)
	us_91 = us_72
	l_92 = l_73
	term_93 = term_74
	cond_ref_94 = cond_ref_75
	uses_95 = uses_76
	and__x_96 = and__x_77
	goto b15
b14:
	;
	v90 = and__x_83
	us_91 = us_78
	l_92 = l_79
	term_93 = term_80
	cond_ref_94 = cond_ref_81
	uses_95 = uses_82
	and__x_96 = and__x_83
	goto b15
b15:
	;
	v99 = v90
	and__x_100 = and__x_53
	us_101 = us_91
	l_102 = l_92
	term_103 = term_93
	cond_ref_104 = cond_ref_94
	uses_105 = uses_95
	goto b12
}
func emit_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var l_3 vm.Value
	var inst_id_4 vm.Value
	var op_kw_5 vm.Value
	var v11 vm.Value
	var l_6 vm.Value
	var inst_id_7 vm.Value
	var op_kw_8 vm.Value
	var v15 vm.Value
	var l_16 vm.Value
	var inst_id_17 vm.Value
	var op_kw_18 vm.Value
	var arg__10673_20 vm.Value
	var arg__10678_22 vm.Value
	var arg__10683_25 vm.Value
	var arg__10688_27 vm.Value
	var v28 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = l_3, inst_id_4, op_kw_5, v11, l_6, inst_id_7, op_kw_8, v15, l_16, inst_id_17, op_kw_18, arg__10673_20, arg__10678_22, arg__10683_25, arg__10688_27, v28
	if vm.IsTruthy(arg1) {
		l_3 = arg0
		inst_id_4 = arg1
		op_kw_5 = arg2
		goto b1
	} else {
		l_6 = arg0
		inst_id_7 = arg1
		op_kw_8 = arg2
		goto b2
	}
b1:
	;
	v11, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "record-source-info!").Deref(), []vm.Value{l_3, inst_id_4})
	if callErr != nil {
		return nil, callErr
	}
	v15 = v11
	l_16 = l_3
	inst_id_17 = inst_id_4
	op_kw_18 = op_kw_5
	goto b3
b2:
	;
	v15 = vm.NIL
	l_16 = l_6
	inst_id_17 = inst_id_7
	op_kw_18 = op_kw_8
	goto b3
b3:
	;
	arg__10673_20, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "chunk-of").Deref(), []vm.Value{l_16})
	if callErr != nil {
		return nil, callErr
	}
	arg__10678_22, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "stack-sp").Deref(), []vm.Value{l_16})
	if callErr != nil {
		return nil, callErr
	}
	arg__10683_25, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "chunk-of").Deref(), []vm.Value{l_16})
	if callErr != nil {
		return nil, callErr
	}
	arg__10688_27, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "stack-sp").Deref(), []vm.Value{l_16})
	if callErr != nil {
		return nil, callErr
	}
	v28, callErr = rt.InvokeValue(rt.LookupVar("ir", "chunk-emit").Deref(), []vm.Value{arg__10683_25, op_kw_18, arg__10688_27})
	if callErr != nil {
		return nil, callErr
	}
	return v28, nil
}
func emit_inst_BANG_(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var f_4 vm.Value
	var op_6 vm.Value
	var refs_8 vm.Value
	var v20 vm.Value
	var l_9 vm.Value
	var nid_10 vm.Value
	var f_11 vm.Value
	var op_12 vm.Value
	var refs_13 vm.Value
	var v23 vm.Value
	var l_14 vm.Value
	var nid_15 vm.Value
	var f_16 vm.Value
	var op_17 vm.Value
	var refs_18 vm.Value
	var v26 vm.Value
	var v28 vm.Value
	var l_29 vm.Value
	var nid_30 vm.Value
	var f_31 vm.Value
	var op_32 vm.Value
	var refs_33 vm.Value
	var v35 vm.Value
	var arg__10726_48 vm.Value
	var v49 bool
	var l_36 vm.Value
	var nid_37 vm.Value
	var f_38 vm.Value
	var op_39 vm.Value
	var refs_40 vm.Value
	var arg__10732_52 vm.Value
	var arg__10740_56 vm.Value
	var arg__10741_57 vm.Value
	var v58 vm.Value
	var l_41 vm.Value
	var nid_42 vm.Value
	var f_43 vm.Value
	var op_44 vm.Value
	var refs_45 vm.Value
	var v62 vm.Value
	var l_63 vm.Value
	var nid_64 vm.Value
	var f_65 vm.Value
	var op_66 vm.Value
	var refs_67 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = f_4, op_6, refs_8, v20, l_9, nid_10, f_11, op_12, refs_13, v23, l_14, nid_15, f_16, op_17, refs_18, v26, v28, l_29, nid_30, f_31, op_32, refs_33, v35, arg__10726_48, v49, l_36, nid_37, f_38, op_39, refs_40, arg__10732_52, arg__10740_56, arg__10741_57, v58, l_41, nid_42, f_43, op_44, refs_45, v62, l_63, nid_64, f_65, op_66, refs_67
	f_4, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "f-of").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	op_6, callErr = rt.InvokeValue(rt.LookupVar("ir", "op").Deref(), []vm.Value{arg1, f_4})
	if callErr != nil {
		return nil, callErr
	}
	refs_8, callErr = rt.InvokeValue(rt.LookupVar("ir", "refs").Deref(), []vm.Value{arg1, f_4})
	if callErr != nil {
		return nil, callErr
	}
	v20, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "refs-at-top-last-use?").Deref(), []vm.Value{arg0, refs_8})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v20) {
		l_9 = arg0
		nid_10 = arg1
		f_11 = f_4
		op_12 = op_6
		refs_13 = refs_8
		goto b1
	} else {
		l_14 = arg0
		nid_15 = arg1
		f_16 = f_4
		op_17 = op_6
		refs_18 = refs_8
		goto b2
	}
b1:
	;
	v23, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "consume-refs-in-place!").Deref(), []vm.Value{l_9, refs_13})
	if callErr != nil {
		return nil, callErr
	}
	v28 = v23
	l_29 = l_9
	nid_30 = nid_10
	f_31 = f_11
	op_32 = op_12
	refs_33 = refs_13
	goto b3
b2:
	;
	v26, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "materialize-refs!").Deref(), []vm.Value{l_14, refs_18})
	if callErr != nil {
		return nil, callErr
	}
	v28 = v26
	l_29 = l_14
	nid_30 = nid_15
	f_31 = f_16
	op_32 = op_17
	refs_33 = refs_18
	goto b3
b3:
	;
	v35, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "lower-node!").Deref(), []vm.Value{l_29, nid_30})
	if callErr != nil {
		return nil, callErr
	}
	arg__10726_48, callErr = rt.InvokeValue(rt.LookupVar("ir", "op-stack-out").Deref(), []vm.Value{op_32})
	if callErr != nil {
		return nil, callErr
	}
	v49 = arg__10726_48 == vm.Int(1)
	if v49 {
		l_36 = l_29
		nid_37 = nid_30
		f_38 = f_31
		op_39 = op_32
		refs_40 = refs_33
		goto b4
	} else {
		l_41 = l_29
		nid_42 = nid_30
		f_43 = f_31
		op_44 = op_32
		refs_45 = refs_33
		goto b5
	}
b4:
	;
	arg__10732_52, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "stack-sp").Deref(), []vm.Value{l_36})
	if callErr != nil {
		return nil, callErr
	}
	arg__10740_56, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "stack-sp").Deref(), []vm.Value{l_36})
	if callErr != nil {
		return nil, callErr
	}
	arg__10741_57 = rt.SubValue(arg__10740_56, vm.Int(1))
	v58, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "set-value-pos!").Deref(), []vm.Value{l_36, nid_37, arg__10741_57})
	if callErr != nil {
		return nil, callErr
	}
	v62 = v58
	l_63 = l_36
	nid_64 = nid_37
	f_65 = f_38
	op_66 = op_39
	refs_67 = refs_40
	goto b6
b5:
	;
	v62 = vm.NIL
	l_63 = l_41
	nid_64 = nid_42
	f_65 = f_43
	op_66 = op_44
	refs_67 = refs_45
	goto b6
b6:
	;
	return v62, nil
}
func emit_placeholder_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var l_4 vm.Value
	var inst_id_5 vm.Value
	var op_kw_6 vm.Value
	var v12 vm.Value
	var l_7 vm.Value
	var inst_id_8 vm.Value
	var op_kw_9 vm.Value
	var v16 vm.Value
	var l_17 vm.Value
	var inst_id_18 vm.Value
	var op_kw_19 vm.Value
	var arg__10750_21 vm.Value
	var arg__10755_23 vm.Value
	var arg__10760_26 vm.Value
	var arg__10765_28 vm.Value
	var v29 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = l_4, inst_id_5, op_kw_6, v12, l_7, inst_id_8, op_kw_9, v16, l_17, inst_id_18, op_kw_19, arg__10750_21, arg__10755_23, arg__10760_26, arg__10765_28, v29
	if vm.IsTruthy(arg1) {
		l_4 = arg0
		inst_id_5 = arg1
		op_kw_6 = arg2
		goto b1
	} else {
		l_7 = arg0
		inst_id_8 = arg1
		op_kw_9 = arg2
		goto b2
	}
b1:
	;
	v12, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "record-source-info!").Deref(), []vm.Value{l_4, inst_id_5})
	if callErr != nil {
		return nil, callErr
	}
	v16 = v12
	l_17 = l_4
	inst_id_18 = inst_id_5
	op_kw_19 = op_kw_6
	goto b3
b2:
	;
	v16 = vm.NIL
	l_17 = l_7
	inst_id_18 = inst_id_8
	op_kw_19 = op_kw_9
	goto b3
b3:
	;
	arg__10750_21, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "chunk-of").Deref(), []vm.Value{l_17})
	if callErr != nil {
		return nil, callErr
	}
	arg__10755_23, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "stack-sp").Deref(), []vm.Value{l_17})
	if callErr != nil {
		return nil, callErr
	}
	arg__10760_26, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "chunk-of").Deref(), []vm.Value{l_17})
	if callErr != nil {
		return nil, callErr
	}
	arg__10765_28, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "stack-sp").Deref(), []vm.Value{l_17})
	if callErr != nil {
		return nil, callErr
	}
	v29, callErr = rt.InvokeValue(rt.LookupVar("ir", "chunk-emit-placeholder").Deref(), []vm.Value{arg__10760_26, op_kw_19, arg__10765_28})
	if callErr != nil {
		return nil, callErr
	}
	return v29, nil
}
func emit_with_arg_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value, arg3 vm.Value) (vm.Value, error) {
	var l_4 vm.Value
	var inst_id_5 vm.Value
	var op_kw_6 vm.Value
	var arg_7 vm.Value
	var v14 vm.Value
	var l_8 vm.Value
	var inst_id_9 vm.Value
	var op_kw_10 vm.Value
	var arg_11 vm.Value
	var v18 vm.Value
	var l_19 vm.Value
	var inst_id_20 vm.Value
	var op_kw_21 vm.Value
	var arg_22 vm.Value
	var arg__10774_24 vm.Value
	var arg__10779_26 vm.Value
	var arg__10785_29 vm.Value
	var arg__10790_31 vm.Value
	var v32 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = l_4, inst_id_5, op_kw_6, arg_7, v14, l_8, inst_id_9, op_kw_10, arg_11, v18, l_19, inst_id_20, op_kw_21, arg_22, arg__10774_24, arg__10779_26, arg__10785_29, arg__10790_31, v32
	if vm.IsTruthy(arg1) {
		l_4 = arg0
		inst_id_5 = arg1
		op_kw_6 = arg2
		arg_7 = arg3
		goto b1
	} else {
		l_8 = arg0
		inst_id_9 = arg1
		op_kw_10 = arg2
		arg_11 = arg3
		goto b2
	}
b1:
	;
	v14, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "record-source-info!").Deref(), []vm.Value{l_4, inst_id_5})
	if callErr != nil {
		return nil, callErr
	}
	v18 = v14
	l_19 = l_4
	inst_id_20 = inst_id_5
	op_kw_21 = op_kw_6
	arg_22 = arg_7
	goto b3
b2:
	;
	v18 = vm.NIL
	l_19 = l_8
	inst_id_20 = inst_id_9
	op_kw_21 = op_kw_10
	arg_22 = arg_11
	goto b3
b3:
	;
	arg__10774_24, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "chunk-of").Deref(), []vm.Value{l_19})
	if callErr != nil {
		return nil, callErr
	}
	arg__10779_26, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "stack-sp").Deref(), []vm.Value{l_19})
	if callErr != nil {
		return nil, callErr
	}
	arg__10785_29, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "chunk-of").Deref(), []vm.Value{l_19})
	if callErr != nil {
		return nil, callErr
	}
	arg__10790_31, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "stack-sp").Deref(), []vm.Value{l_19})
	if callErr != nil {
		return nil, callErr
	}
	v32, callErr = rt.InvokeValue(rt.LookupVar("ir", "chunk-emit-with-arg").Deref(), []vm.Value{arg__10785_29, op_kw_21, arg__10790_31, arg_22})
	if callErr != nil {
		return nil, callErr
	}
	return v32, nil
}
func f_of(arg0 vm.Value) (vm.Value, error) {
	var arg__10796_3 vm.Value
	var v4 vm.Value
	var callErr error
	_, _ = arg__10796_3, v4
	arg__10796_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v4, callErr = rt.InvokeValue(vm.Keyword("f"), []vm.Value{arg__10796_3})
	if callErr != nil {
		return nil, callErr
	}
	return v4, nil
}
func is_terminator_branch_arg_use_QMARK_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var op_5 vm.Value
	var aux_7 vm.Value
	var v19 bool
	var term_id_8 vm.Value
	var target_id_9 vm.Value
	var f_10 vm.Value
	var op_11 vm.Value
	var aux_12 vm.Value
	var arg__10815_25 vm.Value
	var arg__10823_31 vm.Value
	var v32 vm.Value
	var term_id_13 vm.Value
	var target_id_14 vm.Value
	var f_15 vm.Value
	var op_16 vm.Value
	var aux_17 vm.Value
	var v45 bool
	var v134 vm.Value
	var term_id_135 vm.Value
	var target_id_136 vm.Value
	var f_137 vm.Value
	var op_138 vm.Value
	var aux_139 vm.Value
	var term_id_34 vm.Value
	var target_id_35 vm.Value
	var f_36 vm.Value
	var op_37 vm.Value
	var aux_38 vm.Value
	var t_48 vm.Value
	var e_50 vm.Value
	var arg__10838_55 vm.Value
	var arg__10846_61 vm.Value
	var or__x_62 vm.Value
	var term_id_39 vm.Value
	var target_id_40 vm.Value
	var f_41 vm.Value
	var op_42 vm.Value
	var aux_43 vm.Value
	var v127 vm.Value
	var term_id_128 vm.Value
	var target_id_129 vm.Value
	var f_130 vm.Value
	var op_131 vm.Value
	var aux_132 vm.Value
	var term_id_63 vm.Value
	var target_id_64 vm.Value
	var f_65 vm.Value
	var op_66 vm.Value
	var aux_67 vm.Value
	var t_68 vm.Value
	var e_69 vm.Value
	var or__x_70 vm.Value
	var term_id_71 vm.Value
	var target_id_72 vm.Value
	var f_73 vm.Value
	var op_74 vm.Value
	var aux_75 vm.Value
	var t_76 vm.Value
	var e_77 vm.Value
	var or__x_78 vm.Value
	var arg__10853_85 vm.Value
	var arg__10861_91 vm.Value
	var v92 vm.Value
	var v94 vm.Value
	var term_id_95 vm.Value
	var target_id_96 vm.Value
	var f_97 vm.Value
	var op_98 vm.Value
	var aux_99 vm.Value
	var t_100 vm.Value
	var e_101 vm.Value
	var or__x_102 vm.Value
	var term_id_104 vm.Value
	var target_id_105 vm.Value
	var f_106 vm.Value
	var op_107 vm.Value
	var aux_108 vm.Value
	var term_id_109 vm.Value
	var target_id_110 vm.Value
	var f_111 vm.Value
	var op_112 vm.Value
	var aux_113 vm.Value
	var v120 vm.Value
	var term_id_121 vm.Value
	var target_id_122 vm.Value
	var f_123 vm.Value
	var op_124 vm.Value
	var aux_125 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = op_5, aux_7, v19, term_id_8, target_id_9, f_10, op_11, aux_12, arg__10815_25, arg__10823_31, v32, term_id_13, target_id_14, f_15, op_16, aux_17, v45, v134, term_id_135, target_id_136, f_137, op_138, aux_139, term_id_34, target_id_35, f_36, op_37, aux_38, t_48, e_50, arg__10838_55, arg__10846_61, or__x_62, term_id_39, target_id_40, f_41, op_42, aux_43, v127, term_id_128, target_id_129, f_130, op_131, aux_132, term_id_63, target_id_64, f_65, op_66, aux_67, t_68, e_69, or__x_70, term_id_71, target_id_72, f_73, op_74, aux_75, t_76, e_77, or__x_78, arg__10853_85, arg__10861_91, v92, v94, term_id_95, target_id_96, f_97, op_98, aux_99, t_100, e_101, or__x_102, term_id_104, target_id_105, f_106, op_107, aux_108, term_id_109, target_id_110, f_111, op_112, aux_113, v120, term_id_121, target_id_122, f_123, op_124, aux_125
	op_5, callErr = rt.InvokeValue(rt.LookupVar("ir", "op").Deref(), []vm.Value{arg0, arg2})
	if callErr != nil {
		return nil, callErr
	}
	aux_7, callErr = rt.InvokeValue(rt.LookupVar("ir", "aux").Deref(), []vm.Value{arg0, arg2})
	if callErr != nil {
		return nil, callErr
	}
	v19 = op_5 == vm.Keyword("branch")
	if v19 {
		term_id_8 = arg0
		target_id_9 = arg1
		f_10 = arg2
		op_11 = op_5
		aux_12 = aux_7
		goto b1
	} else {
		term_id_13 = arg0
		target_id_14 = arg1
		f_15 = arg2
		op_16 = op_5
		aux_17 = aux_7
		goto b2
	}
b1:
	;
	arg__10815_25, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-args").Deref(), []vm.Value{aux_12})
	if callErr != nil {
		return nil, callErr
	}
	arg__10823_31, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-args").Deref(), []vm.Value{aux_12})
	if callErr != nil {
		return nil, callErr
	}
	v32, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "some").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) vm.Value {
		var v2 vm.Value
		_ = v2
		v2 = vm.Boolean(arg0 == target_id_9)
		return v2
	}), arg__10823_31})
	if callErr != nil {
		return nil, callErr
	}
	v134 = v32
	term_id_135 = term_id_8
	target_id_136 = target_id_9
	f_137 = f_10
	op_138 = op_11
	aux_139 = aux_12
	goto b3
b2:
	;
	v45 = op_16 == vm.Keyword("branch-if")
	if v45 {
		term_id_34 = term_id_13
		target_id_35 = target_id_14
		f_36 = f_15
		op_37 = op_16
		aux_38 = aux_17
		goto b4
	} else {
		term_id_39 = term_id_13
		target_id_40 = target_id_14
		f_41 = f_15
		op_42 = op_16
		aux_43 = aux_17
		goto b5
	}
b3:
	;
	return v134, nil
b4:
	;
	t_48, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-true").Deref(), []vm.Value{aux_38})
	if callErr != nil {
		return nil, callErr
	}
	e_50, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-false").Deref(), []vm.Value{aux_38})
	if callErr != nil {
		return nil, callErr
	}
	arg__10838_55, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-args").Deref(), []vm.Value{t_48})
	if callErr != nil {
		return nil, callErr
	}
	arg__10846_61, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-args").Deref(), []vm.Value{t_48})
	if callErr != nil {
		return nil, callErr
	}
	or__x_62, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "some").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) vm.Value {
		var v2 vm.Value
		_ = v2
		v2 = vm.Boolean(arg0 == target_id_35)
		return v2
	}), arg__10846_61})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(or__x_62) {
		term_id_63 = term_id_34
		target_id_64 = target_id_35
		f_65 = f_36
		op_66 = op_37
		aux_67 = aux_38
		t_68 = t_48
		e_69 = e_50
		or__x_70 = or__x_62
		goto b7
	} else {
		term_id_71 = term_id_34
		target_id_72 = target_id_35
		f_73 = f_36
		op_74 = op_37
		aux_75 = aux_38
		t_76 = t_48
		e_77 = e_50
		or__x_78 = or__x_62
		goto b8
	}
b5:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		term_id_104 = term_id_39
		target_id_105 = target_id_40
		f_106 = f_41
		op_107 = op_42
		aux_108 = aux_43
		goto b10
	} else {
		term_id_109 = term_id_39
		target_id_110 = target_id_40
		f_111 = f_41
		op_112 = op_42
		aux_113 = aux_43
		goto b11
	}
b6:
	;
	v134 = v127
	term_id_135 = term_id_128
	target_id_136 = target_id_129
	f_137 = f_130
	op_138 = op_131
	aux_139 = aux_132
	goto b3
b7:
	;
	v94 = or__x_70
	term_id_95 = term_id_63
	target_id_96 = target_id_64
	f_97 = f_65
	op_98 = op_66
	aux_99 = aux_67
	t_100 = t_68
	e_101 = e_69
	or__x_102 = or__x_70
	goto b9
b8:
	;
	arg__10853_85, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-args").Deref(), []vm.Value{e_77})
	if callErr != nil {
		return nil, callErr
	}
	arg__10861_91, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-args").Deref(), []vm.Value{e_77})
	if callErr != nil {
		return nil, callErr
	}
	v92, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "some").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) vm.Value {
		var v2 vm.Value
		_ = v2
		v2 = vm.Boolean(arg0 == target_id_72)
		return v2
	}), arg__10861_91})
	if callErr != nil {
		return nil, callErr
	}
	v94 = v92
	term_id_95 = term_id_71
	target_id_96 = target_id_72
	f_97 = f_73
	op_98 = op_74
	aux_99 = aux_75
	t_100 = t_76
	e_101 = e_77
	or__x_102 = or__x_78
	goto b9
b9:
	;
	v127 = v94
	term_id_128 = term_id_95
	target_id_129 = target_id_96
	f_130 = f_97
	op_131 = op_98
	aux_132 = aux_99
	goto b6
b10:
	;
	v120 = vm.Boolean(false)
	term_id_121 = term_id_104
	target_id_122 = target_id_105
	f_123 = f_106
	op_124 = op_107
	aux_125 = aux_108
	goto b12
b11:
	;
	v120 = vm.NIL
	term_id_121 = term_id_109
	target_id_122 = target_id_110
	f_123 = f_111
	op_124 = op_112
	aux_125 = aux_113
	goto b12
b12:
	;
	v127 = v120
	term_id_128 = term_id_121
	target_id_129 = target_id_122
	f_130 = f_123
	op_131 = op_124
	aux_132 = aux_125
	goto b6
}
func lower(arg0 vm.Value) (vm.Value, error) {
	var v5 vm.Value
	var f_2 vm.Value
	var v10 vm.Value
	var f_3 vm.Value
	var v14 vm.Value
	var f_15 vm.Value
	var l_17 vm.Value
	var arg__10876_20 vm.Value
	var arg__10877_21 vm.Value
	var arg__10884_25 vm.Value
	var arg__10885_26 vm.Value
	var v27 vm.Value
	var arg__10889_29 vm.Value
	var arg__10894_32 vm.Value
	var n_blocks_33 vm.Value
	var bid_34 int
	var l_35 vm.Value
	var n_blocks_36 vm.Value
	var v47 bool
	var f_39 vm.Value
	var bid_40 int
	var l_41 vm.Value
	var n_blocks_42 vm.Value
	var arg__10902_50 vm.Value
	var arg__10907_53 vm.Value
	var arg__10908_54 vm.Value
	var arg__10915_57 vm.Value
	var arg__10920_60 vm.Value
	var arg__10921_61 vm.Value
	var v62 vm.Value
	var v64 vm.Value
	var v65 int
	var f_43 vm.Value
	var bid_44 int
	var l_45 vm.Value
	var n_blocks_46 vm.Value
	var v69 vm.Value
	var f_70 vm.Value
	var bid_71 int
	var l_72 vm.Value
	var n_blocks_73 vm.Value
	var v75 vm.Value
	var c_77 vm.Value
	var arg__10937_79 vm.Value
	var arg__10942_82 vm.Value
	var n_blocks_83 vm.Value
	var bid_84 int
	var acc_85 vm.Value
	var f_86 vm.Value
	var bid_87 int
	var n_blocks_88 vm.Value
	var v104 bool
	var l_92 vm.Value
	var c_93 vm.Value
	var acc_94 vm.Value
	var f_95 vm.Value
	var bid_96 int
	var n_blocks_97 vm.Value
	var l_98 vm.Value
	var c_99 vm.Value
	var acc_100 vm.Value
	var f_101 vm.Value
	var bid_102 int
	var n_blocks_103 vm.Value
	var arg__10950_108 vm.Value
	var arg__10957_111 vm.Value
	var p_112 vm.Value
	var v113 int
	var v128 bool
	var max_params_141 vm.Value
	var l_142 vm.Value
	var c_143 vm.Value
	var acc_144 vm.Value
	var f_145 vm.Value
	var bid_146 int
	var n_blocks_147 vm.Value
	var arg__10965_149 vm.Value
	var arg__10973_153 vm.Value
	var arg__10975_154 vm.Value
	var v155 vm.Value
	var v157 vm.Value
	var l_114 vm.Value
	var c_115 vm.Value
	var acc_116 vm.Value
	var f_117 vm.Value
	var bid_118 int
	var n_blocks_119 vm.Value
	var p_120 vm.Value
	var l_121 vm.Value
	var c_122 vm.Value
	var acc_123 vm.Value
	var f_124 vm.Value
	var bid_125 int
	var n_blocks_126 vm.Value
	var p_127 vm.Value
	var v132 vm.Value
	var l_133 vm.Value
	var c_134 vm.Value
	var acc_135 vm.Value
	var f_136 vm.Value
	var bid_137 int
	var n_blocks_138 vm.Value
	var p_139 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = v5, f_2, v10, f_3, v14, f_15, l_17, arg__10876_20, arg__10877_21, arg__10884_25, arg__10885_26, v27, arg__10889_29, arg__10894_32, n_blocks_33, bid_34, l_35, n_blocks_36, v47, f_39, bid_40, l_41, n_blocks_42, arg__10902_50, arg__10907_53, arg__10908_54, arg__10915_57, arg__10920_60, arg__10921_61, v62, v64, v65, f_43, bid_44, l_45, n_blocks_46, v69, f_70, bid_71, l_72, n_blocks_73, v75, c_77, arg__10937_79, arg__10942_82, n_blocks_83, bid_84, acc_85, f_86, bid_87, n_blocks_88, v104, l_92, c_93, acc_94, f_95, bid_96, n_blocks_97, l_98, c_99, acc_100, f_101, bid_102, n_blocks_103, arg__10950_108, arg__10957_111, p_112, v113, v128, max_params_141, l_142, c_143, acc_144, f_145, bid_146, n_blocks_147, arg__10965_149, arg__10973_153, arg__10975_154, v155, v157, l_114, c_115, acc_116, f_117, bid_118, n_blocks_119, p_120, l_121, c_122, acc_123, f_124, bid_125, n_blocks_126, p_127, v132, l_133, c_134, acc_135, f_136, bid_137, n_blocks_138, p_139
	v5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v5) {
		f_2 = arg0
		goto b1
	} else {
		f_3 = arg0
		goto b2
	}
b1:
	;
	v10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "throw").Deref(), []vm.Value{vm.String("ir/lower: nil function")})
	if callErr != nil {
		return nil, callErr
	}
	v14 = v10
	f_15 = f_2
	goto b3
b2:
	;
	v14 = vm.NIL
	f_15 = f_3
	goto b3
b3:
	;
	l_17, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "new-lowerer").Deref(), []vm.Value{f_15})
	if callErr != nil {
		return nil, callErr
	}
	arg__10876_20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{l_17})
	if callErr != nil {
		return nil, callErr
	}
	arg__10877_21, callErr = rt.InvokeValue(vm.Keyword("uses"), []vm.Value{arg__10876_20})
	if callErr != nil {
		return nil, callErr
	}
	arg__10884_25, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{l_17})
	if callErr != nil {
		return nil, callErr
	}
	arg__10885_26, callErr = rt.InvokeValue(vm.Keyword("uses"), []vm.Value{arg__10884_25})
	if callErr != nil {
		return nil, callErr
	}
	v27, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "check-cross-block!").Deref(), []vm.Value{f_15, arg__10885_26})
	if callErr != nil {
		return nil, callErr
	}
	arg__10889_29, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{f_15})
	if callErr != nil {
		return nil, callErr
	}
	arg__10894_32, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{f_15})
	if callErr != nil {
		return nil, callErr
	}
	n_blocks_33, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{arg__10894_32})
	if callErr != nil {
		return nil, callErr
	}
	bid_34 = 0
	l_35 = l_17
	n_blocks_36 = n_blocks_33
	goto b4
b4:
	;
	v47 = rt.LtValue(vm.Int(bid_34), n_blocks_36)
	if v47 {
		f_39 = f_15
		bid_40 = bid_34
		l_41 = l_35
		n_blocks_42 = n_blocks_36
		goto b5
	} else {
		f_43 = f_15
		bid_44 = bid_34
		l_45 = l_35
		n_blocks_46 = n_blocks_36
		goto b6
	}
b5:
	;
	arg__10902_50, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "chunk-of").Deref(), []vm.Value{l_41})
	if callErr != nil {
		return nil, callErr
	}
	arg__10907_53, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "chunk-of").Deref(), []vm.Value{l_41})
	if callErr != nil {
		return nil, callErr
	}
	arg__10908_54, callErr = rt.InvokeValue(rt.LookupVar("ir", "chunk-length").Deref(), []vm.Value{arg__10907_53})
	if callErr != nil {
		return nil, callErr
	}
	arg__10915_57, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "chunk-of").Deref(), []vm.Value{l_41})
	if callErr != nil {
		return nil, callErr
	}
	arg__10920_60, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "chunk-of").Deref(), []vm.Value{l_41})
	if callErr != nil {
		return nil, callErr
	}
	arg__10921_61, callErr = rt.InvokeValue(rt.LookupVar("ir", "chunk-length").Deref(), []vm.Value{arg__10920_60})
	if callErr != nil {
		return nil, callErr
	}
	v62, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "set-block-ip!").Deref(), []vm.Value{l_41, vm.Int(bid_40), arg__10921_61})
	if callErr != nil {
		return nil, callErr
	}
	v64, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "lower-block!").Deref(), []vm.Value{l_41, vm.Int(bid_40)})
	if callErr != nil {
		return nil, callErr
	}
	v65 = bid_40 + 1
	bid_34 = v65
	l_35 = l_41
	n_blocks_36 = n_blocks_42
	goto b4
b6:
	;
	v69 = vm.NIL
	f_70 = f_43
	bid_71 = bid_44
	l_72 = l_45
	n_blocks_73 = n_blocks_46
	goto b7
b7:
	;
	v75, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "patch-branches!").Deref(), []vm.Value{l_72})
	if callErr != nil {
		return nil, callErr
	}
	c_77, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "chunk-of").Deref(), []vm.Value{l_72})
	if callErr != nil {
		return nil, callErr
	}
	arg__10937_79, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{f_70})
	if callErr != nil {
		return nil, callErr
	}
	arg__10942_82, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{f_70})
	if callErr != nil {
		return nil, callErr
	}
	n_blocks_83, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{arg__10942_82})
	if callErr != nil {
		return nil, callErr
	}
	bid_84 = 0
	acc_85 = vm.Int(0)
	f_86 = f_70
	bid_87 = bid_71
	n_blocks_88 = n_blocks_83
	goto b8
b8:
	;
	v104 = rt.GeValue(vm.Int(bid_87), n_blocks_88)
	if v104 {
		l_92 = l_72
		c_93 = c_77
		acc_94 = acc_85
		f_95 = f_86
		bid_96 = bid_87
		n_blocks_97 = n_blocks_88
		goto b9
	} else {
		l_98 = l_72
		c_99 = c_77
		acc_100 = acc_85
		f_101 = f_86
		bid_102 = bid_87
		n_blocks_103 = n_blocks_88
		goto b10
	}
b9:
	;
	max_params_141 = acc_94
	l_142 = l_92
	c_143 = c_93
	acc_144 = acc_94
	f_145 = f_95
	bid_146 = bid_96
	n_blocks_147 = n_blocks_97
	goto b11
b10:
	;
	arg__10950_108, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-params").Deref(), []vm.Value{vm.Int(bid_102), f_101})
	if callErr != nil {
		return nil, callErr
	}
	arg__10957_111, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-params").Deref(), []vm.Value{vm.Int(bid_102), f_101})
	if callErr != nil {
		return nil, callErr
	}
	p_112, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{arg__10957_111})
	if callErr != nil {
		return nil, callErr
	}
	v113 = bid_102 + 1
	v128 = rt.GtValue(p_112, acc_100)
	if v128 {
		l_114 = l_98
		c_115 = c_99
		acc_116 = acc_100
		f_117 = f_101
		bid_118 = bid_102
		n_blocks_119 = n_blocks_103
		p_120 = p_112
		goto b12
	} else {
		l_121 = l_98
		c_122 = c_99
		acc_123 = acc_100
		f_124 = f_101
		bid_125 = bid_102
		n_blocks_126 = n_blocks_103
		p_127 = p_112
		goto b13
	}
b11:
	;
	arg__10965_149, callErr = rt.InvokeValue(rt.LookupVar("ir", "chunk-max-stack").Deref(), []vm.Value{c_143})
	if callErr != nil {
		return nil, callErr
	}
	arg__10973_153, callErr = rt.InvokeValue(rt.LookupVar("ir", "chunk-max-stack").Deref(), []vm.Value{c_143})
	if callErr != nil {
		return nil, callErr
	}
	arg__10975_154 = rt.AddValue(arg__10973_153, max_params_141)
	v155, callErr = rt.InvokeValue(rt.LookupVar("ir", "chunk-set-max-stack!").Deref(), []vm.Value{c_143, arg__10975_154})
	if callErr != nil {
		return nil, callErr
	}
	v157, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "chunk-of").Deref(), []vm.Value{l_142})
	if callErr != nil {
		return nil, callErr
	}
	return v157, nil
b12:
	;
	v132 = p_120
	l_133 = l_114
	c_134 = c_115
	acc_135 = acc_116
	f_136 = f_117
	bid_137 = bid_118
	n_blocks_138 = n_blocks_119
	p_139 = p_120
	goto b14
b13:
	;
	v132 = acc_123
	l_133 = l_121
	c_134 = c_122
	acc_135 = acc_123
	f_136 = f_124
	bid_137 = bid_125
	n_blocks_138 = n_blocks_126
	p_139 = p_127
	goto b14
b14:
	;
	bid_84 = v113
	acc_85 = v132
	f_86 = f_101
	bid_87 = bid_102
	n_blocks_88 = n_blocks_103
	goto b8
}
func lower_block_BANG_(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var v7 vm.Value
	var f_9 vm.Value
	var params_11 vm.Value
	var insts_13 vm.Value
	var term_15 vm.Value
	var l_16 vm.Value
	var bid_17 vm.Value
	var f_18 vm.Value
	var params_19 vm.Value
	var insts_20 vm.Value
	var term_21 vm.Value
	var v61 vm.Value
	var l_22 vm.Value
	var bid_23 vm.Value
	var f_24 vm.Value
	var params_25 vm.Value
	var insts_26 vm.Value
	var term_27 vm.Value
	var term_op_65 vm.Value
	var l_66 vm.Value
	var bid_67 vm.Value
	var f_68 vm.Value
	var params_69 vm.Value
	var insts_70 vm.Value
	var term_71 vm.Value
	var v87 bool
	var l_28 vm.Value
	var bid_29 vm.Value
	var f_30 vm.Value
	var params_31 vm.Value
	var insts_32 vm.Value
	var and__x_33 vm.Value
	var term_34 vm.Value
	var arg__11013_44 vm.Value
	var arg__11018_47 vm.Value
	var v48 vm.Value
	var l_35 vm.Value
	var bid_36 vm.Value
	var f_37 vm.Value
	var params_38 vm.Value
	var insts_39 vm.Value
	var and__x_40 vm.Value
	var term_41 vm.Value
	var v51 vm.Value
	var l_52 vm.Value
	var bid_53 vm.Value
	var f_54 vm.Value
	var params_55 vm.Value
	var insts_56 vm.Value
	var and__x_57 vm.Value
	var term_58 vm.Value
	var term_op_72 vm.Value
	var l_73 vm.Value
	var bid_74 vm.Value
	var f_75 vm.Value
	var params_76 vm.Value
	var insts_77 vm.Value
	var term_78 vm.Value
	var arg__11031_90 vm.Value
	var arg__11038_93 vm.Value
	var cond_ref_94 vm.Value
	var v112 vm.Value
	var term_op_79 vm.Value
	var l_80 vm.Value
	var bid_81 vm.Value
	var f_82 vm.Value
	var params_83 vm.Value
	var insts_84 vm.Value
	var term_85 vm.Value
	var deferred_cond_129 vm.Value
	var term_op_130 vm.Value
	var l_131 vm.Value
	var bid_132 vm.Value
	var f_133 vm.Value
	var params_134 vm.Value
	var insts_135 vm.Value
	var term_136 vm.Value
	var arg__11050_138 vm.Value
	var arg__11056_141 vm.Value
	var v142 vm.Value
	var term_op_95 vm.Value
	var l_96 vm.Value
	var bid_97 vm.Value
	var f_98 vm.Value
	var params_99 vm.Value
	var insts_100 vm.Value
	var term_101 vm.Value
	var cond_ref_102 vm.Value
	var term_op_103 vm.Value
	var l_104 vm.Value
	var bid_105 vm.Value
	var f_106 vm.Value
	var params_107 vm.Value
	var insts_108 vm.Value
	var term_109 vm.Value
	var cond_ref_110 vm.Value
	var v117 vm.Value
	var term_op_118 vm.Value
	var l_119 vm.Value
	var bid_120 vm.Value
	var f_121 vm.Value
	var params_122 vm.Value
	var insts_123 vm.Value
	var term_124 vm.Value
	var cond_ref_125 vm.Value
	var i_143 int
	var params_144 vm.Value
	var l_145 vm.Value
	var arg__11061_167 vm.Value
	var v168 bool
	var deferred_cond_148 vm.Value
	var term_op_149 vm.Value
	var bid_150 vm.Value
	var f_151 vm.Value
	var insts_152 vm.Value
	var term_153 vm.Value
	var i_154 int
	var params_155 vm.Value
	var l_156 vm.Value
	var arg__11068_171 vm.Value
	var arg__11077_174 vm.Value
	var v175 vm.Value
	var v176 int
	var deferred_cond_157 vm.Value
	var term_op_158 vm.Value
	var bid_159 vm.Value
	var f_160 vm.Value
	var insts_161 vm.Value
	var term_162 vm.Value
	var i_163 int
	var params_164 vm.Value
	var l_165 vm.Value
	var v180 vm.Value
	var deferred_cond_181 vm.Value
	var term_op_182 vm.Value
	var bid_183 vm.Value
	var f_184 vm.Value
	var insts_185 vm.Value
	var term_186 vm.Value
	var i_187 int
	var params_188 vm.Value
	var l_189 vm.Value
	var doseq_seq__10979_191 vm.Value
	var doseq_loop__10980_192 vm.Value
	var deferred_cond_193 vm.Value
	var f_194 vm.Value
	var l_195 vm.Value
	var v1109 vm.Value
	var v1140 vm.Value
	var v1171 vm.Value
	var v1202 vm.Value
	var term_op_197 vm.Value
	var bid_198 vm.Value
	var insts_199 vm.Value
	var term_200 vm.Value
	var i_201 int
	var params_202 vm.Value
	var doseq_seq__10979_203 vm.Value
	var doseq_loop__10980_204 vm.Value
	var deferred_cond_205 vm.Value
	var f_206 vm.Value
	var l_207 vm.Value
	var v1107 vm.Value
	var v1138 vm.Value
	var v1169 vm.Value
	var v1200 vm.Value
	var nid_221 vm.Value
	var op_223 vm.Value
	var v251 bool
	var term_op_208 vm.Value
	var bid_209 vm.Value
	var insts_210 vm.Value
	var term_211 vm.Value
	var i_212 int
	var params_213 vm.Value
	var doseq_seq__10979_214 vm.Value
	var doseq_loop__10980_215 vm.Value
	var deferred_cond_216 vm.Value
	var f_217 vm.Value
	var l_218 vm.Value
	var v1119 vm.Value
	var v1150 vm.Value
	var v1181 vm.Value
	var v1212 vm.Value
	var v667 vm.Value
	var term_op_668 vm.Value
	var bid_669 vm.Value
	var insts_670 vm.Value
	var term_671 vm.Value
	var i_672 int
	var params_673 vm.Value
	var doseq_seq__10979_674 vm.Value
	var doseq_loop__10980_675 vm.Value
	var deferred_cond_676 vm.Value
	var f_677 vm.Value
	var l_678 vm.Value
	var term_op_224 vm.Value
	var bid_225 vm.Value
	var insts_226 vm.Value
	var term_227 vm.Value
	var i_228 int
	var params_229 vm.Value
	var doseq_seq__10979_230 vm.Value
	var doseq_loop__10980_231 vm.Value
	var deferred_cond_232 vm.Value
	var f_233 vm.Value
	var l_234 vm.Value
	var nid_235 vm.Value
	var op_236 vm.Value
	var v1105 vm.Value
	var v1136 vm.Value
	var v1167 vm.Value
	var v1198 vm.Value
	var term_op_237 vm.Value
	var bid_238 vm.Value
	var insts_239 vm.Value
	var term_240 vm.Value
	var i_241 int
	var params_242 vm.Value
	var doseq_seq__10979_243 vm.Value
	var doseq_loop__10980_244 vm.Value
	var deferred_cond_245 vm.Value
	var f_246 vm.Value
	var l_247 vm.Value
	var nid_248 vm.Value
	var op_249 vm.Value
	var v1110 vm.Value
	var v1141 vm.Value
	var v1172 vm.Value
	var v1203 vm.Value
	var v282 bool
	var v648 vm.Value
	var term_op_649 vm.Value
	var bid_650 vm.Value
	var insts_651 vm.Value
	var term_652 vm.Value
	var i_653 int
	var params_654 vm.Value
	var doseq_seq__10979_655 vm.Value
	var doseq_loop__10980_656 vm.Value
	var deferred_cond_657 vm.Value
	var f_658 vm.Value
	var l_659 vm.Value
	var nid_660 vm.Value
	var op_661 vm.Value
	var v1112 vm.Value
	var v1143 vm.Value
	var v1174 vm.Value
	var v1205 vm.Value
	var v663 vm.Value
	var term_op_255 vm.Value
	var bid_256 vm.Value
	var insts_257 vm.Value
	var term_258 vm.Value
	var i_259 int
	var params_260 vm.Value
	var doseq_seq__10979_261 vm.Value
	var doseq_loop__10980_262 vm.Value
	var deferred_cond_263 vm.Value
	var f_264 vm.Value
	var l_265 vm.Value
	var nid_266 vm.Value
	var op_267 vm.Value
	var v1091 vm.Value
	var v1122 vm.Value
	var v1153 vm.Value
	var v1184 vm.Value
	var term_op_268 vm.Value
	var bid_269 vm.Value
	var insts_270 vm.Value
	var term_271 vm.Value
	var i_272 int
	var params_273 vm.Value
	var doseq_seq__10979_274 vm.Value
	var doseq_loop__10980_275 vm.Value
	var deferred_cond_276 vm.Value
	var f_277 vm.Value
	var l_278 vm.Value
	var nid_279 vm.Value
	var op_280 vm.Value
	var v1116 vm.Value
	var v1147 vm.Value
	var v1178 vm.Value
	var v1209 vm.Value
	var v313 bool
	var v633 vm.Value
	var term_op_634 vm.Value
	var bid_635 vm.Value
	var insts_636 vm.Value
	var term_637 vm.Value
	var i_638 int
	var params_639 vm.Value
	var doseq_seq__10979_640 vm.Value
	var doseq_loop__10980_641 vm.Value
	var deferred_cond_642 vm.Value
	var f_643 vm.Value
	var l_644 vm.Value
	var nid_645 vm.Value
	var op_646 vm.Value
	var v1098 vm.Value
	var v1129 vm.Value
	var v1160 vm.Value
	var v1191 vm.Value
	var term_op_286 vm.Value
	var bid_287 vm.Value
	var insts_288 vm.Value
	var term_289 vm.Value
	var i_290 int
	var params_291 vm.Value
	var doseq_seq__10979_292 vm.Value
	var doseq_loop__10980_293 vm.Value
	var deferred_cond_294 vm.Value
	var f_295 vm.Value
	var l_296 vm.Value
	var nid_297 vm.Value
	var op_298 vm.Value
	var v1097 vm.Value
	var v1128 vm.Value
	var v1159 vm.Value
	var v1190 vm.Value
	var arg__11102_316 vm.Value
	var arg__11109_319 vm.Value
	var doseq_seq__10981_320 vm.Value
	var term_op_299 vm.Value
	var bid_300 vm.Value
	var insts_301 vm.Value
	var term_302 vm.Value
	var i_303 int
	var params_304 vm.Value
	var doseq_seq__10979_305 vm.Value
	var doseq_loop__10980_306 vm.Value
	var deferred_cond_307 vm.Value
	var f_308 vm.Value
	var l_309 vm.Value
	var nid_310 vm.Value
	var op_311 vm.Value
	var v1102 vm.Value
	var v1133 vm.Value
	var v1164 vm.Value
	var v1195 vm.Value
	var v618 vm.Value
	var term_op_619 vm.Value
	var bid_620 vm.Value
	var insts_621 vm.Value
	var term_622 vm.Value
	var i_623 int
	var params_624 vm.Value
	var doseq_seq__10979_625 vm.Value
	var doseq_loop__10980_626 vm.Value
	var deferred_cond_627 vm.Value
	var f_628 vm.Value
	var l_629 vm.Value
	var nid_630 vm.Value
	var op_631 vm.Value
	var v1118 vm.Value
	var v1149 vm.Value
	var v1180 vm.Value
	var v1211 vm.Value
	var doseq_loop__10982_321 vm.Value
	var l_322 vm.Value
	var v1101 vm.Value
	var v1132 vm.Value
	var v1163 vm.Value
	var v1194 vm.Value
	var term_op_324 vm.Value
	var bid_325 vm.Value
	var insts_326 vm.Value
	var term_327 vm.Value
	var i_328 int
	var params_329 vm.Value
	var doseq_seq__10979_330 vm.Value
	var doseq_loop__10980_331 vm.Value
	var deferred_cond_332 vm.Value
	var f_333 vm.Value
	var nid_334 vm.Value
	var op_335 vm.Value
	var doseq_seq__10981_336 vm.Value
	var doseq_loop__10982_337 vm.Value
	var l_338 vm.Value
	var v1092 vm.Value
	var v1123 vm.Value
	var v1154 vm.Value
	var v1185 vm.Value
	var r_356 vm.Value
	var v358 vm.Value
	var v360 vm.Value
	var term_op_339 vm.Value
	var bid_340 vm.Value
	var insts_341 vm.Value
	var term_342 vm.Value
	var i_343 int
	var params_344 vm.Value
	var doseq_seq__10979_345 vm.Value
	var doseq_loop__10980_346 vm.Value
	var deferred_cond_347 vm.Value
	var f_348 vm.Value
	var nid_349 vm.Value
	var op_350 vm.Value
	var doseq_seq__10981_351 vm.Value
	var doseq_loop__10982_352 vm.Value
	var l_353 vm.Value
	var v1117 vm.Value
	var v1148 vm.Value
	var v1179 vm.Value
	var v1210 vm.Value
	var v364 vm.Value
	var term_op_365 vm.Value
	var bid_366 vm.Value
	var insts_367 vm.Value
	var term_368 vm.Value
	var i_369 int
	var params_370 vm.Value
	var doseq_seq__10979_371 vm.Value
	var doseq_loop__10980_372 vm.Value
	var deferred_cond_373 vm.Value
	var f_374 vm.Value
	var nid_375 vm.Value
	var op_376 vm.Value
	var doseq_seq__10981_377 vm.Value
	var doseq_loop__10982_378 vm.Value
	var l_379 vm.Value
	var v1100 vm.Value
	var v1131 vm.Value
	var v1162 vm.Value
	var v1193 vm.Value
	var term_op_381 vm.Value
	var bid_382 vm.Value
	var insts_383 vm.Value
	var term_384 vm.Value
	var i_385 int
	var params_386 vm.Value
	var doseq_seq__10979_387 vm.Value
	var doseq_loop__10980_388 vm.Value
	var deferred_cond_389 vm.Value
	var f_390 vm.Value
	var l_391 vm.Value
	var nid_392 vm.Value
	var op_393 vm.Value
	var v1113 vm.Value
	var v1144 vm.Value
	var v1175 vm.Value
	var v1206 vm.Value
	var term_op_394 vm.Value
	var bid_395 vm.Value
	var insts_396 vm.Value
	var term_397 vm.Value
	var i_398 int
	var params_399 vm.Value
	var doseq_seq__10979_400 vm.Value
	var doseq_loop__10980_401 vm.Value
	var deferred_cond_402 vm.Value
	var f_403 vm.Value
	var l_404 vm.Value
	var nid_405 vm.Value
	var op_406 vm.Value
	var v1093 vm.Value
	var v1124 vm.Value
	var v1155 vm.Value
	var v1186 vm.Value
	var and__x_484 vm.Value
	var v603 vm.Value
	var term_op_604 vm.Value
	var bid_605 vm.Value
	var insts_606 vm.Value
	var term_607 vm.Value
	var i_608 int
	var params_609 vm.Value
	var doseq_seq__10979_610 vm.Value
	var doseq_loop__10980_611 vm.Value
	var deferred_cond_612 vm.Value
	var f_613 vm.Value
	var l_614 vm.Value
	var nid_615 vm.Value
	var op_616 vm.Value
	var v1115 vm.Value
	var v1146 vm.Value
	var v1177 vm.Value
	var v1208 vm.Value
	var term_op_407 vm.Value
	var bid_408 vm.Value
	var insts_409 vm.Value
	var term_410 vm.Value
	var i_411 int
	var params_412 vm.Value
	var doseq_seq__10979_413 vm.Value
	var doseq_loop__10980_414 vm.Value
	var deferred_cond_415 vm.Value
	var and__x_416 vm.Value
	var f_417 vm.Value
	var l_418 vm.Value
	var nid_419 vm.Value
	var op_420 vm.Value
	var v1094 vm.Value
	var v1125 vm.Value
	var v1156 vm.Value
	var v1187 vm.Value
	var v436 bool
	var term_op_421 vm.Value
	var bid_422 vm.Value
	var insts_423 vm.Value
	var term_424 vm.Value
	var i_425 int
	var params_426 vm.Value
	var doseq_seq__10979_427 vm.Value
	var doseq_loop__10980_428 vm.Value
	var deferred_cond_429 vm.Value
	var and__x_430 vm.Value
	var f_431 vm.Value
	var l_432 vm.Value
	var nid_433 vm.Value
	var op_434 vm.Value
	var v1099 vm.Value
	var v1130 vm.Value
	var v1161 vm.Value
	var v1192 vm.Value
	var v439 vm.Value
	var term_op_440 vm.Value
	var bid_441 vm.Value
	var insts_442 vm.Value
	var term_443 vm.Value
	var i_444 int
	var params_445 vm.Value
	var doseq_seq__10979_446 vm.Value
	var doseq_loop__10980_447 vm.Value
	var deferred_cond_448 vm.Value
	var and__x_449 vm.Value
	var f_450 vm.Value
	var l_451 vm.Value
	var nid_452 vm.Value
	var op_453 vm.Value
	var v1095 vm.Value
	var v1126 vm.Value
	var v1157 vm.Value
	var v1188 vm.Value
	var term_op_457 vm.Value
	var bid_458 vm.Value
	var insts_459 vm.Value
	var term_460 vm.Value
	var i_461 int
	var params_462 vm.Value
	var doseq_seq__10979_463 vm.Value
	var doseq_loop__10980_464 vm.Value
	var deferred_cond_465 vm.Value
	var f_466 vm.Value
	var l_467 vm.Value
	var nid_468 vm.Value
	var op_469 vm.Value
	var v1111 vm.Value
	var v1142 vm.Value
	var v1173 vm.Value
	var v1204 vm.Value
	var term_op_470 vm.Value
	var bid_471 vm.Value
	var insts_472 vm.Value
	var term_473 vm.Value
	var i_474 int
	var params_475 vm.Value
	var doseq_seq__10979_476 vm.Value
	var doseq_loop__10980_477 vm.Value
	var deferred_cond_478 vm.Value
	var f_479 vm.Value
	var l_480 vm.Value
	var nid_481 vm.Value
	var op_482 vm.Value
	var v1106 vm.Value
	var v1137 vm.Value
	var v1168 vm.Value
	var v1199 vm.Value
	var v588 vm.Value
	var term_op_589 vm.Value
	var bid_590 vm.Value
	var insts_591 vm.Value
	var term_592 vm.Value
	var i_593 int
	var params_594 vm.Value
	var doseq_seq__10979_595 vm.Value
	var doseq_loop__10980_596 vm.Value
	var deferred_cond_597 vm.Value
	var f_598 vm.Value
	var l_599 vm.Value
	var nid_600 vm.Value
	var op_601 vm.Value
	var v1096 vm.Value
	var v1127 vm.Value
	var v1158 vm.Value
	var v1189 vm.Value
	var term_op_485 vm.Value
	var bid_486 vm.Value
	var insts_487 vm.Value
	var term_488 vm.Value
	var i_489 int
	var params_490 vm.Value
	var doseq_seq__10979_491 vm.Value
	var doseq_loop__10980_492 vm.Value
	var deferred_cond_493 vm.Value
	var f_494 vm.Value
	var l_495 vm.Value
	var nid_496 vm.Value
	var op_497 vm.Value
	var and__x_498 vm.Value
	var v1108 vm.Value
	var v1139 vm.Value
	var v1170 vm.Value
	var v1201 vm.Value
	var arg__11131_515 vm.Value
	var arg__11138_518 vm.Value
	var v519 vm.Value
	var term_op_499 vm.Value
	var bid_500 vm.Value
	var insts_501 vm.Value
	var term_502 vm.Value
	var i_503 int
	var params_504 vm.Value
	var doseq_seq__10979_505 vm.Value
	var doseq_loop__10980_506 vm.Value
	var deferred_cond_507 vm.Value
	var f_508 vm.Value
	var l_509 vm.Value
	var nid_510 vm.Value
	var op_511 vm.Value
	var and__x_512 vm.Value
	var v1114 vm.Value
	var v1145 vm.Value
	var v1176 vm.Value
	var v1207 vm.Value
	var v522 vm.Value
	var term_op_523 vm.Value
	var bid_524 vm.Value
	var insts_525 vm.Value
	var term_526 vm.Value
	var i_527 int
	var params_528 vm.Value
	var doseq_seq__10979_529 vm.Value
	var doseq_loop__10980_530 vm.Value
	var deferred_cond_531 vm.Value
	var f_532 vm.Value
	var l_533 vm.Value
	var nid_534 vm.Value
	var op_535 vm.Value
	var and__x_536 vm.Value
	var v1089 vm.Value
	var v1120 vm.Value
	var v1151 vm.Value
	var v1182 vm.Value
	var term_op_540 vm.Value
	var bid_541 vm.Value
	var insts_542 vm.Value
	var term_543 vm.Value
	var i_544 int
	var params_545 vm.Value
	var doseq_seq__10979_546 vm.Value
	var doseq_loop__10980_547 vm.Value
	var deferred_cond_548 vm.Value
	var f_549 vm.Value
	var l_550 vm.Value
	var nid_551 vm.Value
	var op_552 vm.Value
	var v1090 vm.Value
	var v1121 vm.Value
	var v1152 vm.Value
	var v1183 vm.Value
	var v569 vm.Value
	var term_op_553 vm.Value
	var bid_554 vm.Value
	var insts_555 vm.Value
	var term_556 vm.Value
	var i_557 int
	var params_558 vm.Value
	var doseq_seq__10979_559 vm.Value
	var doseq_loop__10980_560 vm.Value
	var deferred_cond_561 vm.Value
	var f_562 vm.Value
	var l_563 vm.Value
	var nid_564 vm.Value
	var op_565 vm.Value
	var v1103 vm.Value
	var v1134 vm.Value
	var v1165 vm.Value
	var v1196 vm.Value
	var v573 vm.Value
	var term_op_574 vm.Value
	var bid_575 vm.Value
	var insts_576 vm.Value
	var term_577 vm.Value
	var i_578 int
	var params_579 vm.Value
	var doseq_seq__10979_580 vm.Value
	var doseq_loop__10980_581 vm.Value
	var deferred_cond_582 vm.Value
	var f_583 vm.Value
	var l_584 vm.Value
	var nid_585 vm.Value
	var op_586 vm.Value
	var v1104 vm.Value
	var v1135 vm.Value
	var v1166 vm.Value
	var v1197 vm.Value
	var term_op_679 vm.Value
	var bid_680 vm.Value
	var insts_681 vm.Value
	var term_682 vm.Value
	var i_683 int
	var params_684 vm.Value
	var doseq_loop__10980_685 vm.Value
	var deferred_cond_686 vm.Value
	var f_687 vm.Value
	var l_688 vm.Value
	var v764 bool
	var term_op_689 vm.Value
	var bid_690 vm.Value
	var insts_691 vm.Value
	var term_692 vm.Value
	var i_693 int
	var params_694 vm.Value
	var doseq_loop__10980_695 vm.Value
	var deferred_cond_696 vm.Value
	var f_697 vm.Value
	var l_698 vm.Value
	var v1067 vm.Value
	var term_op_1068 vm.Value
	var bid_1069 vm.Value
	var insts_1070 vm.Value
	var term_1071 vm.Value
	var i_1072 int
	var params_1073 vm.Value
	var doseq_loop__10980_1074 vm.Value
	var deferred_cond_1075 vm.Value
	var f_1076 vm.Value
	var l_1077 vm.Value
	var term_op_699 vm.Value
	var bid_700 vm.Value
	var insts_701 vm.Value
	var and__x_702 vm.Value
	var term_703 vm.Value
	var i_704 int
	var params_705 vm.Value
	var doseq_loop__10980_706 vm.Value
	var deferred_cond_707 vm.Value
	var f_708 vm.Value
	var l_709 vm.Value
	var arg__11150_723 vm.Value
	var arg__11155_726 vm.Value
	var v727 vm.Value
	var term_op_710 vm.Value
	var bid_711 vm.Value
	var insts_712 vm.Value
	var and__x_713 vm.Value
	var term_714 vm.Value
	var i_715 int
	var params_716 vm.Value
	var doseq_loop__10980_717 vm.Value
	var deferred_cond_718 vm.Value
	var f_719 vm.Value
	var l_720 vm.Value
	var v730 vm.Value
	var term_op_731 vm.Value
	var bid_732 vm.Value
	var insts_733 vm.Value
	var and__x_734 vm.Value
	var term_735 vm.Value
	var i_736 int
	var params_737 vm.Value
	var doseq_loop__10980_738 vm.Value
	var deferred_cond_739 vm.Value
	var f_740 vm.Value
	var l_741 vm.Value
	var term_op_743 vm.Value
	var bid_744 vm.Value
	var insts_745 vm.Value
	var term_746 vm.Value
	var i_747 int
	var params_748 vm.Value
	var doseq_loop__10980_749 vm.Value
	var deferred_cond_750 vm.Value
	var f_751 vm.Value
	var l_752 vm.Value
	var cond_aux_767 vm.Value
	var t_bt_769 vm.Value
	var args_771 vm.Value
	var v773 vm.Value
	var term_op_753 vm.Value
	var bid_754 vm.Value
	var insts_755 vm.Value
	var term_756 vm.Value
	var i_757 int
	var params_758 vm.Value
	var doseq_loop__10980_759 vm.Value
	var deferred_cond_760 vm.Value
	var f_761 vm.Value
	var l_762 vm.Value
	var v1053 vm.Value
	var term_op_1054 vm.Value
	var bid_1055 vm.Value
	var insts_1056 vm.Value
	var term_1057 vm.Value
	var i_1058 int
	var params_1059 vm.Value
	var doseq_loop__10980_1060 vm.Value
	var deferred_cond_1061 vm.Value
	var f_1062 vm.Value
	var l_1063 vm.Value
	var term_op_774 vm.Value
	var bid_775 vm.Value
	var insts_776 vm.Value
	var term_777 vm.Value
	var i_778 int
	var params_779 vm.Value
	var doseq_loop__10980_780 vm.Value
	var deferred_cond_781 vm.Value
	var f_782 vm.Value
	var l_783 vm.Value
	var cond_aux_784 vm.Value
	var t_bt_785 vm.Value
	var args_786 vm.Value
	var v802 vm.Value
	var term_op_787 vm.Value
	var bid_788 vm.Value
	var insts_789 vm.Value
	var term_790 vm.Value
	var i_791 int
	var params_792 vm.Value
	var doseq_loop__10980_793 vm.Value
	var deferred_cond_794 vm.Value
	var f_795 vm.Value
	var l_796 vm.Value
	var cond_aux_797 vm.Value
	var t_bt_798 vm.Value
	var args_799 vm.Value
	var v904 vm.Value
	var term_op_905 vm.Value
	var bid_906 vm.Value
	var insts_907 vm.Value
	var term_908 vm.Value
	var i_909 int
	var params_910 vm.Value
	var doseq_loop__10980_911 vm.Value
	var deferred_cond_912 vm.Value
	var f_913 vm.Value
	var l_914 vm.Value
	var cond_aux_915 vm.Value
	var t_bt_916 vm.Value
	var args_917 vm.Value
	var v919 vm.Value
	var term_op_804 vm.Value
	var bid_805 vm.Value
	var insts_806 vm.Value
	var term_807 vm.Value
	var i_808 int
	var params_809 vm.Value
	var doseq_loop__10980_810 vm.Value
	var deferred_cond_811 vm.Value
	var f_812 vm.Value
	var l_813 vm.Value
	var cond_aux_814 vm.Value
	var t_bt_815 vm.Value
	var args_816 vm.Value
	var term_refs_833 vm.Value
	var v863 vm.Value
	var term_op_817 vm.Value
	var bid_818 vm.Value
	var insts_819 vm.Value
	var term_820 vm.Value
	var i_821 int
	var params_822 vm.Value
	var doseq_loop__10980_823 vm.Value
	var deferred_cond_824 vm.Value
	var f_825 vm.Value
	var l_826 vm.Value
	var cond_aux_827 vm.Value
	var t_bt_828 vm.Value
	var args_829 vm.Value
	var v889 vm.Value
	var term_op_890 vm.Value
	var bid_891 vm.Value
	var insts_892 vm.Value
	var term_893 vm.Value
	var i_894 int
	var params_895 vm.Value
	var doseq_loop__10980_896 vm.Value
	var deferred_cond_897 vm.Value
	var f_898 vm.Value
	var l_899 vm.Value
	var cond_aux_900 vm.Value
	var t_bt_901 vm.Value
	var args_902 vm.Value
	var term_op_834 vm.Value
	var bid_835 vm.Value
	var insts_836 vm.Value
	var term_837 vm.Value
	var i_838 int
	var params_839 vm.Value
	var doseq_loop__10980_840 vm.Value
	var deferred_cond_841 vm.Value
	var f_842 vm.Value
	var l_843 vm.Value
	var cond_aux_844 vm.Value
	var t_bt_845 vm.Value
	var args_846 vm.Value
	var term_refs_847 vm.Value
	var v866 vm.Value
	var term_op_848 vm.Value
	var bid_849 vm.Value
	var insts_850 vm.Value
	var term_851 vm.Value
	var i_852 int
	var params_853 vm.Value
	var doseq_loop__10980_854 vm.Value
	var deferred_cond_855 vm.Value
	var f_856 vm.Value
	var l_857 vm.Value
	var cond_aux_858 vm.Value
	var t_bt_859 vm.Value
	var args_860 vm.Value
	var term_refs_861 vm.Value
	var v869 vm.Value
	var v871 vm.Value
	var term_op_872 vm.Value
	var bid_873 vm.Value
	var insts_874 vm.Value
	var term_875 vm.Value
	var i_876 int
	var params_877 vm.Value
	var doseq_loop__10980_878 vm.Value
	var deferred_cond_879 vm.Value
	var f_880 vm.Value
	var l_881 vm.Value
	var cond_aux_882 vm.Value
	var t_bt_883 vm.Value
	var args_884 vm.Value
	var term_refs_885 vm.Value
	var term_op_921 vm.Value
	var bid_922 vm.Value
	var insts_923 vm.Value
	var term_924 vm.Value
	var i_925 int
	var params_926 vm.Value
	var doseq_loop__10980_927 vm.Value
	var deferred_cond_928 vm.Value
	var f_929 vm.Value
	var l_930 vm.Value
	var term_refs_944 vm.Value
	var v968 vm.Value
	var term_op_931 vm.Value
	var bid_932 vm.Value
	var insts_933 vm.Value
	var term_934 vm.Value
	var i_935 int
	var params_936 vm.Value
	var doseq_loop__10980_937 vm.Value
	var deferred_cond_938 vm.Value
	var f_939 vm.Value
	var l_940 vm.Value
	var v1041 vm.Value
	var term_op_1042 vm.Value
	var bid_1043 vm.Value
	var insts_1044 vm.Value
	var term_1045 vm.Value
	var i_1046 int
	var params_1047 vm.Value
	var doseq_loop__10980_1048 vm.Value
	var deferred_cond_1049 vm.Value
	var f_1050 vm.Value
	var l_1051 vm.Value
	var term_op_945 vm.Value
	var bid_946 vm.Value
	var insts_947 vm.Value
	var term_948 vm.Value
	var i_949 int
	var params_950 vm.Value
	var doseq_loop__10980_951 vm.Value
	var deferred_cond_952 vm.Value
	var f_953 vm.Value
	var l_954 vm.Value
	var term_refs_955 vm.Value
	var v971 vm.Value
	var term_op_956 vm.Value
	var bid_957 vm.Value
	var insts_958 vm.Value
	var term_959 vm.Value
	var i_960 int
	var params_961 vm.Value
	var doseq_loop__10980_962 vm.Value
	var deferred_cond_963 vm.Value
	var f_964 vm.Value
	var l_965 vm.Value
	var term_refs_966 vm.Value
	var v974 vm.Value
	var v976 vm.Value
	var term_op_977 vm.Value
	var bid_978 vm.Value
	var insts_979 vm.Value
	var term_980 vm.Value
	var i_981 int
	var params_982 vm.Value
	var doseq_loop__10980_983 vm.Value
	var deferred_cond_984 vm.Value
	var f_985 vm.Value
	var l_986 vm.Value
	var term_refs_987 vm.Value
	var v1011 bool
	var term_op_988 vm.Value
	var bid_989 vm.Value
	var insts_990 vm.Value
	var term_991 vm.Value
	var i_992 int
	var params_993 vm.Value
	var doseq_loop__10980_994 vm.Value
	var deferred_cond_995 vm.Value
	var f_996 vm.Value
	var l_997 vm.Value
	var term_refs_998 vm.Value
	var bt_1014 vm.Value
	var arg__11235_1016 vm.Value
	var arg__11241_1019 vm.Value
	var v1020 vm.Value
	var term_op_999 vm.Value
	var bid_1000 vm.Value
	var insts_1001 vm.Value
	var term_1002 vm.Value
	var i_1003 int
	var params_1004 vm.Value
	var doseq_loop__10980_1005 vm.Value
	var deferred_cond_1006 vm.Value
	var f_1007 vm.Value
	var l_1008 vm.Value
	var term_refs_1009 vm.Value
	var v1024 vm.Value
	var term_op_1025 vm.Value
	var bid_1026 vm.Value
	var insts_1027 vm.Value
	var term_1028 vm.Value
	var i_1029 int
	var params_1030 vm.Value
	var doseq_loop__10980_1031 vm.Value
	var deferred_cond_1032 vm.Value
	var f_1033 vm.Value
	var l_1034 vm.Value
	var term_refs_1035 vm.Value
	var v1037 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = v7, f_9, params_11, insts_13, term_15, l_16, bid_17, f_18, params_19, insts_20, term_21, v61, l_22, bid_23, f_24, params_25, insts_26, term_27, term_op_65, l_66, bid_67, f_68, params_69, insts_70, term_71, v87, l_28, bid_29, f_30, params_31, insts_32, and__x_33, term_34, arg__11013_44, arg__11018_47, v48, l_35, bid_36, f_37, params_38, insts_39, and__x_40, term_41, v51, l_52, bid_53, f_54, params_55, insts_56, and__x_57, term_58, term_op_72, l_73, bid_74, f_75, params_76, insts_77, term_78, arg__11031_90, arg__11038_93, cond_ref_94, v112, term_op_79, l_80, bid_81, f_82, params_83, insts_84, term_85, deferred_cond_129, term_op_130, l_131, bid_132, f_133, params_134, insts_135, term_136, arg__11050_138, arg__11056_141, v142, term_op_95, l_96, bid_97, f_98, params_99, insts_100, term_101, cond_ref_102, term_op_103, l_104, bid_105, f_106, params_107, insts_108, term_109, cond_ref_110, v117, term_op_118, l_119, bid_120, f_121, params_122, insts_123, term_124, cond_ref_125, i_143, params_144, l_145, arg__11061_167, v168, deferred_cond_148, term_op_149, bid_150, f_151, insts_152, term_153, i_154, params_155, l_156, arg__11068_171, arg__11077_174, v175, v176, deferred_cond_157, term_op_158, bid_159, f_160, insts_161, term_162, i_163, params_164, l_165, v180, deferred_cond_181, term_op_182, bid_183, f_184, insts_185, term_186, i_187, params_188, l_189, doseq_seq__10979_191, doseq_loop__10980_192, deferred_cond_193, f_194, l_195, v1109, v1140, v1171, v1202, term_op_197, bid_198, insts_199, term_200, i_201, params_202, doseq_seq__10979_203, doseq_loop__10980_204, deferred_cond_205, f_206, l_207, v1107, v1138, v1169, v1200, nid_221, op_223, v251, term_op_208, bid_209, insts_210, term_211, i_212, params_213, doseq_seq__10979_214, doseq_loop__10980_215, deferred_cond_216, f_217, l_218, v1119, v1150, v1181, v1212, v667, term_op_668, bid_669, insts_670, term_671, i_672, params_673, doseq_seq__10979_674, doseq_loop__10980_675, deferred_cond_676, f_677, l_678, term_op_224, bid_225, insts_226, term_227, i_228, params_229, doseq_seq__10979_230, doseq_loop__10980_231, deferred_cond_232, f_233, l_234, nid_235, op_236, v1105, v1136, v1167, v1198, term_op_237, bid_238, insts_239, term_240, i_241, params_242, doseq_seq__10979_243, doseq_loop__10980_244, deferred_cond_245, f_246, l_247, nid_248, op_249, v1110, v1141, v1172, v1203, v282, v648, term_op_649, bid_650, insts_651, term_652, i_653, params_654, doseq_seq__10979_655, doseq_loop__10980_656, deferred_cond_657, f_658, l_659, nid_660, op_661, v1112, v1143, v1174, v1205, v663, term_op_255, bid_256, insts_257, term_258, i_259, params_260, doseq_seq__10979_261, doseq_loop__10980_262, deferred_cond_263, f_264, l_265, nid_266, op_267, v1091, v1122, v1153, v1184, term_op_268, bid_269, insts_270, term_271, i_272, params_273, doseq_seq__10979_274, doseq_loop__10980_275, deferred_cond_276, f_277, l_278, nid_279, op_280, v1116, v1147, v1178, v1209, v313, v633, term_op_634, bid_635, insts_636, term_637, i_638, params_639, doseq_seq__10979_640, doseq_loop__10980_641, deferred_cond_642, f_643, l_644, nid_645, op_646, v1098, v1129, v1160, v1191, term_op_286, bid_287, insts_288, term_289, i_290, params_291, doseq_seq__10979_292, doseq_loop__10980_293, deferred_cond_294, f_295, l_296, nid_297, op_298, v1097, v1128, v1159, v1190, arg__11102_316, arg__11109_319, doseq_seq__10981_320, term_op_299, bid_300, insts_301, term_302, i_303, params_304, doseq_seq__10979_305, doseq_loop__10980_306, deferred_cond_307, f_308, l_309, nid_310, op_311, v1102, v1133, v1164, v1195, v618, term_op_619, bid_620, insts_621, term_622, i_623, params_624, doseq_seq__10979_625, doseq_loop__10980_626, deferred_cond_627, f_628, l_629, nid_630, op_631, v1118, v1149, v1180, v1211, doseq_loop__10982_321, l_322, v1101, v1132, v1163, v1194, term_op_324, bid_325, insts_326, term_327, i_328, params_329, doseq_seq__10979_330, doseq_loop__10980_331, deferred_cond_332, f_333, nid_334, op_335, doseq_seq__10981_336, doseq_loop__10982_337, l_338, v1092, v1123, v1154, v1185, r_356, v358, v360, term_op_339, bid_340, insts_341, term_342, i_343, params_344, doseq_seq__10979_345, doseq_loop__10980_346, deferred_cond_347, f_348, nid_349, op_350, doseq_seq__10981_351, doseq_loop__10982_352, l_353, v1117, v1148, v1179, v1210, v364, term_op_365, bid_366, insts_367, term_368, i_369, params_370, doseq_seq__10979_371, doseq_loop__10980_372, deferred_cond_373, f_374, nid_375, op_376, doseq_seq__10981_377, doseq_loop__10982_378, l_379, v1100, v1131, v1162, v1193, term_op_381, bid_382, insts_383, term_384, i_385, params_386, doseq_seq__10979_387, doseq_loop__10980_388, deferred_cond_389, f_390, l_391, nid_392, op_393, v1113, v1144, v1175, v1206, term_op_394, bid_395, insts_396, term_397, i_398, params_399, doseq_seq__10979_400, doseq_loop__10980_401, deferred_cond_402, f_403, l_404, nid_405, op_406, v1093, v1124, v1155, v1186, and__x_484, v603, term_op_604, bid_605, insts_606, term_607, i_608, params_609, doseq_seq__10979_610, doseq_loop__10980_611, deferred_cond_612, f_613, l_614, nid_615, op_616, v1115, v1146, v1177, v1208, term_op_407, bid_408, insts_409, term_410, i_411, params_412, doseq_seq__10979_413, doseq_loop__10980_414, deferred_cond_415, and__x_416, f_417, l_418, nid_419, op_420, v1094, v1125, v1156, v1187, v436, term_op_421, bid_422, insts_423, term_424, i_425, params_426, doseq_seq__10979_427, doseq_loop__10980_428, deferred_cond_429, and__x_430, f_431, l_432, nid_433, op_434, v1099, v1130, v1161, v1192, v439, term_op_440, bid_441, insts_442, term_443, i_444, params_445, doseq_seq__10979_446, doseq_loop__10980_447, deferred_cond_448, and__x_449, f_450, l_451, nid_452, op_453, v1095, v1126, v1157, v1188, term_op_457, bid_458, insts_459, term_460, i_461, params_462, doseq_seq__10979_463, doseq_loop__10980_464, deferred_cond_465, f_466, l_467, nid_468, op_469, v1111, v1142, v1173, v1204, term_op_470, bid_471, insts_472, term_473, i_474, params_475, doseq_seq__10979_476, doseq_loop__10980_477, deferred_cond_478, f_479, l_480, nid_481, op_482, v1106, v1137, v1168, v1199, v588, term_op_589, bid_590, insts_591, term_592, i_593, params_594, doseq_seq__10979_595, doseq_loop__10980_596, deferred_cond_597, f_598, l_599, nid_600, op_601, v1096, v1127, v1158, v1189, term_op_485, bid_486, insts_487, term_488, i_489, params_490, doseq_seq__10979_491, doseq_loop__10980_492, deferred_cond_493, f_494, l_495, nid_496, op_497, and__x_498, v1108, v1139, v1170, v1201, arg__11131_515, arg__11138_518, v519, term_op_499, bid_500, insts_501, term_502, i_503, params_504, doseq_seq__10979_505, doseq_loop__10980_506, deferred_cond_507, f_508, l_509, nid_510, op_511, and__x_512, v1114, v1145, v1176, v1207, v522, term_op_523, bid_524, insts_525, term_526, i_527, params_528, doseq_seq__10979_529, doseq_loop__10980_530, deferred_cond_531, f_532, l_533, nid_534, op_535, and__x_536, v1089, v1120, v1151, v1182, term_op_540, bid_541, insts_542, term_543, i_544, params_545, doseq_seq__10979_546, doseq_loop__10980_547, deferred_cond_548, f_549, l_550, nid_551, op_552, v1090, v1121, v1152, v1183, v569, term_op_553, bid_554, insts_555, term_556, i_557, params_558, doseq_seq__10979_559, doseq_loop__10980_560, deferred_cond_561, f_562, l_563, nid_564, op_565, v1103, v1134, v1165, v1196, v573, term_op_574, bid_575, insts_576, term_577, i_578, params_579, doseq_seq__10979_580, doseq_loop__10980_581, deferred_cond_582, f_583, l_584, nid_585, op_586, v1104, v1135, v1166, v1197, term_op_679, bid_680, insts_681, term_682, i_683, params_684, doseq_loop__10980_685, deferred_cond_686, f_687, l_688, v764, term_op_689, bid_690, insts_691, term_692, i_693, params_694, doseq_loop__10980_695, deferred_cond_696, f_697, l_698, v1067, term_op_1068, bid_1069, insts_1070, term_1071, i_1072, params_1073, doseq_loop__10980_1074, deferred_cond_1075, f_1076, l_1077, term_op_699, bid_700, insts_701, and__x_702, term_703, i_704, params_705, doseq_loop__10980_706, deferred_cond_707, f_708, l_709, arg__11150_723, arg__11155_726, v727, term_op_710, bid_711, insts_712, and__x_713, term_714, i_715, params_716, doseq_loop__10980_717, deferred_cond_718, f_719, l_720, v730, term_op_731, bid_732, insts_733, and__x_734, term_735, i_736, params_737, doseq_loop__10980_738, deferred_cond_739, f_740, l_741, term_op_743, bid_744, insts_745, term_746, i_747, params_748, doseq_loop__10980_749, deferred_cond_750, f_751, l_752, cond_aux_767, t_bt_769, args_771, v773, term_op_753, bid_754, insts_755, term_756, i_757, params_758, doseq_loop__10980_759, deferred_cond_760, f_761, l_762, v1053, term_op_1054, bid_1055, insts_1056, term_1057, i_1058, params_1059, doseq_loop__10980_1060, deferred_cond_1061, f_1062, l_1063, term_op_774, bid_775, insts_776, term_777, i_778, params_779, doseq_loop__10980_780, deferred_cond_781, f_782, l_783, cond_aux_784, t_bt_785, args_786, v802, term_op_787, bid_788, insts_789, term_790, i_791, params_792, doseq_loop__10980_793, deferred_cond_794, f_795, l_796, cond_aux_797, t_bt_798, args_799, v904, term_op_905, bid_906, insts_907, term_908, i_909, params_910, doseq_loop__10980_911, deferred_cond_912, f_913, l_914, cond_aux_915, t_bt_916, args_917, v919, term_op_804, bid_805, insts_806, term_807, i_808, params_809, doseq_loop__10980_810, deferred_cond_811, f_812, l_813, cond_aux_814, t_bt_815, args_816, term_refs_833, v863, term_op_817, bid_818, insts_819, term_820, i_821, params_822, doseq_loop__10980_823, deferred_cond_824, f_825, l_826, cond_aux_827, t_bt_828, args_829, v889, term_op_890, bid_891, insts_892, term_893, i_894, params_895, doseq_loop__10980_896, deferred_cond_897, f_898, l_899, cond_aux_900, t_bt_901, args_902, term_op_834, bid_835, insts_836, term_837, i_838, params_839, doseq_loop__10980_840, deferred_cond_841, f_842, l_843, cond_aux_844, t_bt_845, args_846, term_refs_847, v866, term_op_848, bid_849, insts_850, term_851, i_852, params_853, doseq_loop__10980_854, deferred_cond_855, f_856, l_857, cond_aux_858, t_bt_859, args_860, term_refs_861, v869, v871, term_op_872, bid_873, insts_874, term_875, i_876, params_877, doseq_loop__10980_878, deferred_cond_879, f_880, l_881, cond_aux_882, t_bt_883, args_884, term_refs_885, term_op_921, bid_922, insts_923, term_924, i_925, params_926, doseq_loop__10980_927, deferred_cond_928, f_929, l_930, term_refs_944, v968, term_op_931, bid_932, insts_933, term_934, i_935, params_936, doseq_loop__10980_937, deferred_cond_938, f_939, l_940, v1041, term_op_1042, bid_1043, insts_1044, term_1045, i_1046, params_1047, doseq_loop__10980_1048, deferred_cond_1049, f_1050, l_1051, term_op_945, bid_946, insts_947, term_948, i_949, params_950, doseq_loop__10980_951, deferred_cond_952, f_953, l_954, term_refs_955, v971, term_op_956, bid_957, insts_958, term_959, i_960, params_961, doseq_loop__10980_962, deferred_cond_963, f_964, l_965, term_refs_966, v974, v976, term_op_977, bid_978, insts_979, term_980, i_981, params_982, doseq_loop__10980_983, deferred_cond_984, f_985, l_986, term_refs_987, v1011, term_op_988, bid_989, insts_990, term_991, i_992, params_993, doseq_loop__10980_994, deferred_cond_995, f_996, l_997, term_refs_998, bt_1014, arg__11235_1016, arg__11241_1019, v1020, term_op_999, bid_1000, insts_1001, term_1002, i_1003, params_1004, doseq_loop__10980_1005, deferred_cond_1006, f_1007, l_1008, term_refs_1009, v1024, term_op_1025, bid_1026, insts_1027, term_1028, i_1029, params_1030, doseq_loop__10980_1031, deferred_cond_1032, f_1033, l_1034, term_refs_1035, v1037
	v7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "assoc").Deref(), vm.Keyword("current-block"), arg1})
	if callErr != nil {
		return nil, callErr
	}
	f_9, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "f-of").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	params_11, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-params").Deref(), []vm.Value{arg1, f_9})
	if callErr != nil {
		return nil, callErr
	}
	insts_13, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-insts").Deref(), []vm.Value{arg1, f_9})
	if callErr != nil {
		return nil, callErr
	}
	term_15, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-term").Deref(), []vm.Value{arg1, f_9})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(term_15) {
		l_28 = arg0
		bid_29 = arg1
		f_30 = f_9
		params_31 = params_11
		insts_32 = insts_13
		and__x_33 = term_15
		term_34 = term_15
		goto b4
	} else {
		l_35 = arg0
		bid_36 = arg1
		f_37 = f_9
		params_38 = params_11
		insts_39 = insts_13
		and__x_40 = term_15
		term_41 = term_15
		goto b5
	}
b1:
	;
	v61, callErr = rt.InvokeValue(rt.LookupVar("ir", "op").Deref(), []vm.Value{term_21, f_18})
	if callErr != nil {
		return nil, callErr
	}
	term_op_65 = v61
	l_66 = l_16
	bid_67 = bid_17
	f_68 = f_18
	params_69 = params_19
	insts_70 = insts_20
	term_71 = term_21
	goto b3
b2:
	;
	term_op_65 = vm.NIL
	l_66 = l_22
	bid_67 = bid_23
	f_68 = f_24
	params_69 = params_25
	insts_70 = insts_26
	term_71 = term_27
	goto b3
b3:
	;
	v87 = term_op_65 == vm.Keyword("branch-if")
	if v87 {
		term_op_72 = term_op_65
		l_73 = l_66
		bid_74 = bid_67
		f_75 = f_68
		params_76 = params_69
		insts_77 = insts_70
		term_78 = term_71
		goto b7
	} else {
		term_op_79 = term_op_65
		l_80 = l_66
		bid_81 = bid_67
		f_82 = f_68
		params_83 = params_69
		insts_84 = insts_70
		term_85 = term_71
		goto b8
	}
b4:
	;
	arg__11013_44, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "zero?").Deref(), []vm.Value{term_34})
	if callErr != nil {
		return nil, callErr
	}
	arg__11018_47, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "zero?").Deref(), []vm.Value{term_34})
	if callErr != nil {
		return nil, callErr
	}
	v48, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not").Deref(), []vm.Value{arg__11018_47})
	if callErr != nil {
		return nil, callErr
	}
	v51 = v48
	l_52 = l_28
	bid_53 = bid_29
	f_54 = f_30
	params_55 = params_31
	insts_56 = insts_32
	and__x_57 = and__x_33
	term_58 = term_34
	goto b6
b5:
	;
	v51 = and__x_40
	l_52 = l_35
	bid_53 = bid_36
	f_54 = f_37
	params_55 = params_38
	insts_56 = insts_39
	and__x_57 = and__x_40
	term_58 = term_41
	goto b6
b6:
	;
	if vm.IsTruthy(v51) {
		l_16 = l_52
		bid_17 = bid_53
		f_18 = f_54
		params_19 = params_55
		insts_20 = insts_56
		term_21 = term_58
		goto b1
	} else {
		l_22 = l_52
		bid_23 = bid_53
		f_24 = f_54
		params_25 = params_55
		insts_26 = insts_56
		term_27 = term_58
		goto b2
	}
b7:
	;
	arg__11031_90, callErr = rt.InvokeValue(rt.LookupVar("ir", "refs").Deref(), []vm.Value{term_78, f_75})
	if callErr != nil {
		return nil, callErr
	}
	arg__11038_93, callErr = rt.InvokeValue(rt.LookupVar("ir", "refs").Deref(), []vm.Value{term_78, f_75})
	if callErr != nil {
		return nil, callErr
	}
	cond_ref_94, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{arg__11038_93})
	if callErr != nil {
		return nil, callErr
	}
	v112, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "deferrable-branch-if-cond?").Deref(), []vm.Value{l_73, term_78, cond_ref_94})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v112) {
		term_op_95 = term_op_72
		l_96 = l_73
		bid_97 = bid_74
		f_98 = f_75
		params_99 = params_76
		insts_100 = insts_77
		term_101 = term_78
		cond_ref_102 = cond_ref_94
		goto b10
	} else {
		term_op_103 = term_op_72
		l_104 = l_73
		bid_105 = bid_74
		f_106 = f_75
		params_107 = params_76
		insts_108 = insts_77
		term_109 = term_78
		cond_ref_110 = cond_ref_94
		goto b11
	}
b8:
	;
	deferred_cond_129 = vm.NIL
	term_op_130 = term_op_79
	l_131 = l_80
	bid_132 = bid_81
	f_133 = f_82
	params_134 = params_83
	insts_135 = insts_84
	term_136 = term_85
	goto b9
b9:
	;
	arg__11050_138, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{params_134})
	if callErr != nil {
		return nil, callErr
	}
	arg__11056_141, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{params_134})
	if callErr != nil {
		return nil, callErr
	}
	v142, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "set-stack-sp!").Deref(), []vm.Value{l_131, arg__11056_141})
	if callErr != nil {
		return nil, callErr
	}
	i_143 = 0
	params_144 = params_134
	l_145 = l_131
	goto b13
b10:
	;
	v117 = cond_ref_102
	term_op_118 = term_op_95
	l_119 = l_96
	bid_120 = bid_97
	f_121 = f_98
	params_122 = params_99
	insts_123 = insts_100
	term_124 = term_101
	cond_ref_125 = cond_ref_102
	goto b12
b11:
	;
	v117 = vm.NIL
	term_op_118 = term_op_103
	l_119 = l_104
	bid_120 = bid_105
	f_121 = f_106
	params_122 = params_107
	insts_123 = insts_108
	term_124 = term_109
	cond_ref_125 = cond_ref_110
	goto b12
b12:
	;
	deferred_cond_129 = v117
	term_op_130 = term_op_118
	l_131 = l_119
	bid_132 = bid_120
	f_133 = f_121
	params_134 = params_122
	insts_135 = insts_123
	term_136 = term_124
	goto b9
b13:
	;
	arg__11061_167, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{params_144})
	if callErr != nil {
		return nil, callErr
	}
	v168 = rt.LtValue(vm.Int(i_143), arg__11061_167)
	if v168 {
		deferred_cond_148 = deferred_cond_129
		term_op_149 = term_op_130
		bid_150 = bid_132
		f_151 = f_133
		insts_152 = insts_135
		term_153 = term_136
		i_154 = i_143
		params_155 = params_144
		l_156 = l_145
		goto b14
	} else {
		deferred_cond_157 = deferred_cond_129
		term_op_158 = term_op_130
		bid_159 = bid_132
		f_160 = f_133
		insts_161 = insts_135
		term_162 = term_136
		i_163 = i_143
		params_164 = params_144
		l_165 = l_145
		goto b15
	}
b14:
	;
	arg__11068_171, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{params_155, vm.Int(i_154)})
	if callErr != nil {
		return nil, callErr
	}
	arg__11077_174, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{params_155, vm.Int(i_154)})
	if callErr != nil {
		return nil, callErr
	}
	v175, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "set-value-pos!").Deref(), []vm.Value{l_156, arg__11077_174, vm.Int(i_154)})
	if callErr != nil {
		return nil, callErr
	}
	v176 = i_154 + 1
	i_143 = v176
	params_144 = params_155
	l_145 = l_156
	goto b13
b15:
	;
	v180 = vm.NIL
	deferred_cond_181 = deferred_cond_157
	term_op_182 = term_op_158
	bid_183 = bid_159
	f_184 = f_160
	insts_185 = insts_161
	term_186 = term_162
	i_187 = i_163
	params_188 = params_164
	l_189 = l_165
	goto b16
b16:
	;
	doseq_seq__10979_191, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{insts_185})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__10980_192 = doseq_seq__10979_191
	deferred_cond_193 = deferred_cond_181
	f_194 = f_184
	l_195 = l_189
	v1109 = vm.Keyword("invalid")
	v1140 = vm.Keyword("block-arg")
	v1171 = vm.Keyword("pop")
	v1202 = vm.Keyword("else")
	goto b17
b17:
	;
	if vm.IsTruthy(doseq_loop__10980_192) {
		term_op_197 = term_op_182
		bid_198 = bid_183
		insts_199 = insts_185
		term_200 = term_186
		i_201 = i_187
		params_202 = params_188
		doseq_seq__10979_203 = doseq_seq__10979_191
		doseq_loop__10980_204 = doseq_loop__10980_192
		deferred_cond_205 = deferred_cond_193
		f_206 = f_194
		l_207 = l_195
		v1107 = v1109
		v1138 = v1140
		v1169 = v1171
		v1200 = v1202
		goto b18
	} else {
		term_op_208 = term_op_182
		bid_209 = bid_183
		insts_210 = insts_185
		term_211 = term_186
		i_212 = i_187
		params_213 = params_188
		doseq_seq__10979_214 = doseq_seq__10979_191
		doseq_loop__10980_215 = doseq_loop__10980_192
		deferred_cond_216 = deferred_cond_193
		f_217 = f_194
		l_218 = l_195
		v1119 = v1109
		v1150 = v1140
		v1181 = v1171
		v1212 = v1202
		goto b19
	}
b18:
	;
	nid_221, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{doseq_loop__10980_204})
	if callErr != nil {
		return nil, callErr
	}
	op_223, callErr = rt.InvokeValue(rt.LookupVar("ir", "op").Deref(), []vm.Value{nid_221, f_206})
	if callErr != nil {
		return nil, callErr
	}
	v251 = op_223 == v1107
	if v251 {
		term_op_224 = term_op_197
		bid_225 = bid_198
		insts_226 = insts_199
		term_227 = term_200
		i_228 = i_201
		params_229 = params_202
		doseq_seq__10979_230 = doseq_seq__10979_203
		doseq_loop__10980_231 = doseq_loop__10980_204
		deferred_cond_232 = deferred_cond_205
		f_233 = f_206
		l_234 = l_207
		nid_235 = nid_221
		op_236 = op_223
		v1105 = v1107
		v1136 = v1138
		v1167 = v1169
		v1198 = v1200
		goto b21
	} else {
		term_op_237 = term_op_197
		bid_238 = bid_198
		insts_239 = insts_199
		term_240 = term_200
		i_241 = i_201
		params_242 = params_202
		doseq_seq__10979_243 = doseq_seq__10979_203
		doseq_loop__10980_244 = doseq_loop__10980_204
		deferred_cond_245 = deferred_cond_205
		f_246 = f_206
		l_247 = l_207
		nid_248 = nid_221
		op_249 = op_223
		v1110 = v1107
		v1141 = v1138
		v1172 = v1169
		v1203 = v1200
		goto b22
	}
b19:
	;
	v667 = vm.NIL
	term_op_668 = term_op_208
	bid_669 = bid_209
	insts_670 = insts_210
	term_671 = term_211
	i_672 = i_212
	params_673 = params_213
	doseq_seq__10979_674 = doseq_seq__10979_214
	doseq_loop__10980_675 = doseq_loop__10980_215
	deferred_cond_676 = deferred_cond_216
	f_677 = f_217
	l_678 = l_218
	goto b20
b20:
	;
	if vm.IsTruthy(term_671) {
		term_op_699 = term_op_668
		bid_700 = bid_669
		insts_701 = insts_670
		and__x_702 = term_671
		term_703 = term_671
		i_704 = i_672
		params_705 = params_673
		doseq_loop__10980_706 = doseq_loop__10980_675
		deferred_cond_707 = deferred_cond_676
		f_708 = f_677
		l_709 = l_678
		goto b52
	} else {
		term_op_710 = term_op_668
		bid_711 = bid_669
		insts_712 = insts_670
		and__x_713 = term_671
		term_714 = term_671
		i_715 = i_672
		params_716 = params_673
		doseq_loop__10980_717 = doseq_loop__10980_675
		deferred_cond_718 = deferred_cond_676
		f_719 = f_677
		l_720 = l_678
		goto b53
	}
b21:
	;
	v648 = vm.NIL
	term_op_649 = term_op_224
	bid_650 = bid_225
	insts_651 = insts_226
	term_652 = term_227
	i_653 = i_228
	params_654 = params_229
	doseq_seq__10979_655 = doseq_seq__10979_230
	doseq_loop__10980_656 = doseq_loop__10980_231
	deferred_cond_657 = deferred_cond_232
	f_658 = f_233
	l_659 = l_234
	nid_660 = nid_235
	op_661 = op_236
	v1112 = v1105
	v1143 = v1136
	v1174 = v1167
	v1205 = v1198
	goto b23
b22:
	;
	v282 = op_249 == v1141
	if v282 {
		term_op_255 = term_op_237
		bid_256 = bid_238
		insts_257 = insts_239
		term_258 = term_240
		i_259 = i_241
		params_260 = params_242
		doseq_seq__10979_261 = doseq_seq__10979_243
		doseq_loop__10980_262 = doseq_loop__10980_244
		deferred_cond_263 = deferred_cond_245
		f_264 = f_246
		l_265 = l_247
		nid_266 = nid_248
		op_267 = op_249
		v1091 = v1110
		v1122 = v1141
		v1153 = v1172
		v1184 = v1203
		goto b24
	} else {
		term_op_268 = term_op_237
		bid_269 = bid_238
		insts_270 = insts_239
		term_271 = term_240
		i_272 = i_241
		params_273 = params_242
		doseq_seq__10979_274 = doseq_seq__10979_243
		doseq_loop__10980_275 = doseq_loop__10980_244
		deferred_cond_276 = deferred_cond_245
		f_277 = f_246
		l_278 = l_247
		nid_279 = nid_248
		op_280 = op_249
		v1116 = v1110
		v1147 = v1141
		v1178 = v1172
		v1209 = v1203
		goto b25
	}
b23:
	;
	v663, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "next").Deref(), []vm.Value{doseq_loop__10980_656})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__10980_192 = v663
	deferred_cond_193 = deferred_cond_657
	f_194 = f_658
	l_195 = l_659
	v1109 = v1112
	v1140 = v1143
	v1171 = v1174
	v1202 = v1205
	goto b17
b24:
	;
	v633 = vm.NIL
	term_op_634 = term_op_255
	bid_635 = bid_256
	insts_636 = insts_257
	term_637 = term_258
	i_638 = i_259
	params_639 = params_260
	doseq_seq__10979_640 = doseq_seq__10979_261
	doseq_loop__10980_641 = doseq_loop__10980_262
	deferred_cond_642 = deferred_cond_263
	f_643 = f_264
	l_644 = l_265
	nid_645 = nid_266
	op_646 = op_267
	v1098 = v1091
	v1129 = v1122
	v1160 = v1153
	v1191 = v1184
	goto b26
b25:
	;
	v313 = op_280 == v1178
	if v313 {
		term_op_286 = term_op_268
		bid_287 = bid_269
		insts_288 = insts_270
		term_289 = term_271
		i_290 = i_272
		params_291 = params_273
		doseq_seq__10979_292 = doseq_seq__10979_274
		doseq_loop__10980_293 = doseq_loop__10980_275
		deferred_cond_294 = deferred_cond_276
		f_295 = f_277
		l_296 = l_278
		nid_297 = nid_279
		op_298 = op_280
		v1097 = v1116
		v1128 = v1147
		v1159 = v1178
		v1190 = v1209
		goto b27
	} else {
		term_op_299 = term_op_268
		bid_300 = bid_269
		insts_301 = insts_270
		term_302 = term_271
		i_303 = i_272
		params_304 = params_273
		doseq_seq__10979_305 = doseq_seq__10979_274
		doseq_loop__10980_306 = doseq_loop__10980_275
		deferred_cond_307 = deferred_cond_276
		f_308 = f_277
		l_309 = l_278
		nid_310 = nid_279
		op_311 = op_280
		v1102 = v1116
		v1133 = v1147
		v1164 = v1178
		v1195 = v1209
		goto b28
	}
b26:
	;
	v648 = v633
	term_op_649 = term_op_634
	bid_650 = bid_635
	insts_651 = insts_636
	term_652 = term_637
	i_653 = i_638
	params_654 = params_639
	doseq_seq__10979_655 = doseq_seq__10979_640
	doseq_loop__10980_656 = doseq_loop__10980_641
	deferred_cond_657 = deferred_cond_642
	f_658 = f_643
	l_659 = l_644
	nid_660 = nid_645
	op_661 = op_646
	v1112 = v1098
	v1143 = v1129
	v1174 = v1160
	v1205 = v1191
	goto b23
b27:
	;
	arg__11102_316, callErr = rt.InvokeValue(rt.LookupVar("ir", "refs").Deref(), []vm.Value{nid_297, f_295})
	if callErr != nil {
		return nil, callErr
	}
	arg__11109_319, callErr = rt.InvokeValue(rt.LookupVar("ir", "refs").Deref(), []vm.Value{nid_297, f_295})
	if callErr != nil {
		return nil, callErr
	}
	doseq_seq__10981_320, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{arg__11109_319})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__10982_321 = doseq_seq__10981_320
	l_322 = l_296
	v1101 = v1097
	v1132 = v1128
	v1163 = v1159
	v1194 = v1190
	goto b30
b28:
	;
	if vm.IsTruthy(deferred_cond_307) {
		term_op_407 = term_op_299
		bid_408 = bid_300
		insts_409 = insts_301
		term_410 = term_302
		i_411 = i_303
		params_412 = params_304
		doseq_seq__10979_413 = doseq_seq__10979_305
		doseq_loop__10980_414 = doseq_loop__10980_306
		deferred_cond_415 = deferred_cond_307
		and__x_416 = deferred_cond_307
		f_417 = f_308
		l_418 = l_309
		nid_419 = nid_310
		op_420 = op_311
		v1094 = v1102
		v1125 = v1133
		v1156 = v1164
		v1187 = v1195
		goto b37
	} else {
		term_op_421 = term_op_299
		bid_422 = bid_300
		insts_423 = insts_301
		term_424 = term_302
		i_425 = i_303
		params_426 = params_304
		doseq_seq__10979_427 = doseq_seq__10979_305
		doseq_loop__10980_428 = doseq_loop__10980_306
		deferred_cond_429 = deferred_cond_307
		and__x_430 = deferred_cond_307
		f_431 = f_308
		l_432 = l_309
		nid_433 = nid_310
		op_434 = op_311
		v1099 = v1102
		v1130 = v1133
		v1161 = v1164
		v1192 = v1195
		goto b38
	}
b29:
	;
	v633 = v618
	term_op_634 = term_op_619
	bid_635 = bid_620
	insts_636 = insts_621
	term_637 = term_622
	i_638 = i_623
	params_639 = params_624
	doseq_seq__10979_640 = doseq_seq__10979_625
	doseq_loop__10980_641 = doseq_loop__10980_626
	deferred_cond_642 = deferred_cond_627
	f_643 = f_628
	l_644 = l_629
	nid_645 = nid_630
	op_646 = op_631
	v1098 = v1118
	v1129 = v1149
	v1160 = v1180
	v1191 = v1211
	goto b26
b30:
	;
	if vm.IsTruthy(doseq_loop__10982_321) {
		term_op_324 = term_op_286
		bid_325 = bid_287
		insts_326 = insts_288
		term_327 = term_289
		i_328 = i_290
		params_329 = params_291
		doseq_seq__10979_330 = doseq_seq__10979_292
		doseq_loop__10980_331 = doseq_loop__10980_293
		deferred_cond_332 = deferred_cond_294
		f_333 = f_295
		nid_334 = nid_297
		op_335 = op_298
		doseq_seq__10981_336 = doseq_seq__10981_320
		doseq_loop__10982_337 = doseq_loop__10982_321
		l_338 = l_322
		v1092 = v1101
		v1123 = v1132
		v1154 = v1163
		v1185 = v1194
		goto b31
	} else {
		term_op_339 = term_op_286
		bid_340 = bid_287
		insts_341 = insts_288
		term_342 = term_289
		i_343 = i_290
		params_344 = params_291
		doseq_seq__10979_345 = doseq_seq__10979_292
		doseq_loop__10980_346 = doseq_loop__10980_293
		deferred_cond_347 = deferred_cond_294
		f_348 = f_295
		nid_349 = nid_297
		op_350 = op_298
		doseq_seq__10981_351 = doseq_seq__10981_320
		doseq_loop__10982_352 = doseq_loop__10982_321
		l_353 = l_322
		v1117 = v1101
		v1148 = v1132
		v1179 = v1163
		v1210 = v1194
		goto b32
	}
b31:
	;
	r_356, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{doseq_loop__10982_337})
	if callErr != nil {
		return nil, callErr
	}
	v358, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "decrement-use!").Deref(), []vm.Value{l_338, r_356})
	if callErr != nil {
		return nil, callErr
	}
	v360, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "next").Deref(), []vm.Value{doseq_loop__10982_337})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__10982_321 = v360
	l_322 = l_338
	v1101 = v1092
	v1132 = v1123
	v1163 = v1154
	v1194 = v1185
	goto b30
b32:
	;
	v364 = vm.NIL
	term_op_365 = term_op_339
	bid_366 = bid_340
	insts_367 = insts_341
	term_368 = term_342
	i_369 = i_343
	params_370 = params_344
	doseq_seq__10979_371 = doseq_seq__10979_345
	doseq_loop__10980_372 = doseq_loop__10980_346
	deferred_cond_373 = deferred_cond_347
	f_374 = f_348
	nid_375 = nid_349
	op_376 = op_350
	doseq_seq__10981_377 = doseq_seq__10981_351
	doseq_loop__10982_378 = doseq_loop__10982_352
	l_379 = l_353
	v1100 = v1117
	v1131 = v1148
	v1162 = v1179
	v1193 = v1210
	goto b33
b33:
	;
	v618 = v364
	term_op_619 = term_op_365
	bid_620 = bid_366
	insts_621 = insts_367
	term_622 = term_368
	i_623 = i_369
	params_624 = params_370
	doseq_seq__10979_625 = doseq_seq__10979_371
	doseq_loop__10980_626 = doseq_loop__10980_372
	deferred_cond_627 = deferred_cond_373
	f_628 = f_374
	l_629 = l_379
	nid_630 = nid_375
	op_631 = op_376
	v1118 = v1100
	v1149 = v1131
	v1180 = v1162
	v1211 = v1193
	goto b29
b34:
	;
	v603 = vm.NIL
	term_op_604 = term_op_381
	bid_605 = bid_382
	insts_606 = insts_383
	term_607 = term_384
	i_608 = i_385
	params_609 = params_386
	doseq_seq__10979_610 = doseq_seq__10979_387
	doseq_loop__10980_611 = doseq_loop__10980_388
	deferred_cond_612 = deferred_cond_389
	f_613 = f_390
	l_614 = l_391
	nid_615 = nid_392
	op_616 = op_393
	v1115 = v1113
	v1146 = v1144
	v1177 = v1175
	v1208 = v1206
	goto b36
b35:
	;
	and__x_484, callErr = rt.InvokeValue(rt.LookupVar("ir", "op-cheap-load?").Deref(), []vm.Value{op_406})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(and__x_484) {
		term_op_485 = term_op_394
		bid_486 = bid_395
		insts_487 = insts_396
		term_488 = term_397
		i_489 = i_398
		params_490 = params_399
		doseq_seq__10979_491 = doseq_seq__10979_400
		doseq_loop__10980_492 = doseq_loop__10980_401
		deferred_cond_493 = deferred_cond_402
		f_494 = f_403
		l_495 = l_404
		nid_496 = nid_405
		op_497 = op_406
		and__x_498 = and__x_484
		v1108 = v1093
		v1139 = v1124
		v1170 = v1155
		v1201 = v1186
		goto b43
	} else {
		term_op_499 = term_op_394
		bid_500 = bid_395
		insts_501 = insts_396
		term_502 = term_397
		i_503 = i_398
		params_504 = params_399
		doseq_seq__10979_505 = doseq_seq__10979_400
		doseq_loop__10980_506 = doseq_loop__10980_401
		deferred_cond_507 = deferred_cond_402
		f_508 = f_403
		l_509 = l_404
		nid_510 = nid_405
		op_511 = op_406
		and__x_512 = and__x_484
		v1114 = v1093
		v1145 = v1124
		v1176 = v1155
		v1207 = v1186
		goto b44
	}
b36:
	;
	v618 = v603
	term_op_619 = term_op_604
	bid_620 = bid_605
	insts_621 = insts_606
	term_622 = term_607
	i_623 = i_608
	params_624 = params_609
	doseq_seq__10979_625 = doseq_seq__10979_610
	doseq_loop__10980_626 = doseq_loop__10980_611
	deferred_cond_627 = deferred_cond_612
	f_628 = f_613
	l_629 = l_614
	nid_630 = nid_615
	op_631 = op_616
	v1118 = v1115
	v1149 = v1146
	v1180 = v1177
	v1211 = v1208
	goto b29
b37:
	;
	v436 = nid_419 == deferred_cond_415
	v439 = vm.Boolean(v436)
	term_op_440 = term_op_407
	bid_441 = bid_408
	insts_442 = insts_409
	term_443 = term_410
	i_444 = i_411
	params_445 = params_412
	doseq_seq__10979_446 = doseq_seq__10979_413
	doseq_loop__10980_447 = doseq_loop__10980_414
	deferred_cond_448 = deferred_cond_415
	and__x_449 = and__x_416
	f_450 = f_417
	l_451 = l_418
	nid_452 = nid_419
	op_453 = op_420
	v1095 = v1094
	v1126 = v1125
	v1157 = v1156
	v1188 = v1187
	goto b39
b38:
	;
	v439 = and__x_430
	term_op_440 = term_op_421
	bid_441 = bid_422
	insts_442 = insts_423
	term_443 = term_424
	i_444 = i_425
	params_445 = params_426
	doseq_seq__10979_446 = doseq_seq__10979_427
	doseq_loop__10980_447 = doseq_loop__10980_428
	deferred_cond_448 = deferred_cond_429
	and__x_449 = and__x_430
	f_450 = f_431
	l_451 = l_432
	nid_452 = nid_433
	op_453 = op_434
	v1095 = v1099
	v1126 = v1130
	v1157 = v1161
	v1188 = v1192
	goto b39
b39:
	;
	if vm.IsTruthy(v439) {
		term_op_381 = term_op_440
		bid_382 = bid_441
		insts_383 = insts_442
		term_384 = term_443
		i_385 = i_444
		params_386 = params_445
		doseq_seq__10979_387 = doseq_seq__10979_446
		doseq_loop__10980_388 = doseq_loop__10980_447
		deferred_cond_389 = deferred_cond_448
		f_390 = f_450
		l_391 = l_451
		nid_392 = nid_452
		op_393 = op_453
		v1113 = v1095
		v1144 = v1126
		v1175 = v1157
		v1206 = v1188
		goto b34
	} else {
		term_op_394 = term_op_440
		bid_395 = bid_441
		insts_396 = insts_442
		term_397 = term_443
		i_398 = i_444
		params_399 = params_445
		doseq_seq__10979_400 = doseq_seq__10979_446
		doseq_loop__10980_401 = doseq_loop__10980_447
		deferred_cond_402 = deferred_cond_448
		f_403 = f_450
		l_404 = l_451
		nid_405 = nid_452
		op_406 = op_453
		v1093 = v1095
		v1124 = v1126
		v1155 = v1157
		v1186 = v1188
		goto b35
	}
b40:
	;
	v588 = vm.NIL
	term_op_589 = term_op_457
	bid_590 = bid_458
	insts_591 = insts_459
	term_592 = term_460
	i_593 = i_461
	params_594 = params_462
	doseq_seq__10979_595 = doseq_seq__10979_463
	doseq_loop__10980_596 = doseq_loop__10980_464
	deferred_cond_597 = deferred_cond_465
	f_598 = f_466
	l_599 = l_467
	nid_600 = nid_468
	op_601 = op_469
	v1096 = v1111
	v1127 = v1142
	v1158 = v1173
	v1189 = v1204
	goto b42
b41:
	;
	if vm.IsTruthy(v1199) {
		term_op_540 = term_op_470
		bid_541 = bid_471
		insts_542 = insts_472
		term_543 = term_473
		i_544 = i_474
		params_545 = params_475
		doseq_seq__10979_546 = doseq_seq__10979_476
		doseq_loop__10980_547 = doseq_loop__10980_477
		deferred_cond_548 = deferred_cond_478
		f_549 = f_479
		l_550 = l_480
		nid_551 = nid_481
		op_552 = op_482
		v1090 = v1106
		v1121 = v1137
		v1152 = v1168
		v1183 = v1199
		goto b46
	} else {
		term_op_553 = term_op_470
		bid_554 = bid_471
		insts_555 = insts_472
		term_556 = term_473
		i_557 = i_474
		params_558 = params_475
		doseq_seq__10979_559 = doseq_seq__10979_476
		doseq_loop__10980_560 = doseq_loop__10980_477
		deferred_cond_561 = deferred_cond_478
		f_562 = f_479
		l_563 = l_480
		nid_564 = nid_481
		op_565 = op_482
		v1103 = v1106
		v1134 = v1137
		v1165 = v1168
		v1196 = v1199
		goto b47
	}
b42:
	;
	v603 = v588
	term_op_604 = term_op_589
	bid_605 = bid_590
	insts_606 = insts_591
	term_607 = term_592
	i_608 = i_593
	params_609 = params_594
	doseq_seq__10979_610 = doseq_seq__10979_595
	doseq_loop__10980_611 = doseq_loop__10980_596
	deferred_cond_612 = deferred_cond_597
	f_613 = f_598
	l_614 = l_599
	nid_615 = nid_600
	op_616 = op_601
	v1115 = v1096
	v1146 = v1127
	v1177 = v1158
	v1208 = v1189
	goto b36
b43:
	;
	arg__11131_515, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "should-body-emit-cheap?").Deref(), []vm.Value{l_495, nid_496})
	if callErr != nil {
		return nil, callErr
	}
	arg__11138_518, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "should-body-emit-cheap?").Deref(), []vm.Value{l_495, nid_496})
	if callErr != nil {
		return nil, callErr
	}
	v519, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not").Deref(), []vm.Value{arg__11138_518})
	if callErr != nil {
		return nil, callErr
	}
	v522 = v519
	term_op_523 = term_op_485
	bid_524 = bid_486
	insts_525 = insts_487
	term_526 = term_488
	i_527 = i_489
	params_528 = params_490
	doseq_seq__10979_529 = doseq_seq__10979_491
	doseq_loop__10980_530 = doseq_loop__10980_492
	deferred_cond_531 = deferred_cond_493
	f_532 = f_494
	l_533 = l_495
	nid_534 = nid_496
	op_535 = op_497
	and__x_536 = and__x_498
	v1089 = v1108
	v1120 = v1139
	v1151 = v1170
	v1182 = v1201
	goto b45
b44:
	;
	v522 = and__x_512
	term_op_523 = term_op_499
	bid_524 = bid_500
	insts_525 = insts_501
	term_526 = term_502
	i_527 = i_503
	params_528 = params_504
	doseq_seq__10979_529 = doseq_seq__10979_505
	doseq_loop__10980_530 = doseq_loop__10980_506
	deferred_cond_531 = deferred_cond_507
	f_532 = f_508
	l_533 = l_509
	nid_534 = nid_510
	op_535 = op_511
	and__x_536 = and__x_512
	v1089 = v1114
	v1120 = v1145
	v1151 = v1176
	v1182 = v1207
	goto b45
b45:
	;
	if vm.IsTruthy(v522) {
		term_op_457 = term_op_523
		bid_458 = bid_524
		insts_459 = insts_525
		term_460 = term_526
		i_461 = i_527
		params_462 = params_528
		doseq_seq__10979_463 = doseq_seq__10979_529
		doseq_loop__10980_464 = doseq_loop__10980_530
		deferred_cond_465 = deferred_cond_531
		f_466 = f_532
		l_467 = l_533
		nid_468 = nid_534
		op_469 = op_535
		v1111 = v1089
		v1142 = v1120
		v1173 = v1151
		v1204 = v1182
		goto b40
	} else {
		term_op_470 = term_op_523
		bid_471 = bid_524
		insts_472 = insts_525
		term_473 = term_526
		i_474 = i_527
		params_475 = params_528
		doseq_seq__10979_476 = doseq_seq__10979_529
		doseq_loop__10980_477 = doseq_loop__10980_530
		deferred_cond_478 = deferred_cond_531
		f_479 = f_532
		l_480 = l_533
		nid_481 = nid_534
		op_482 = op_535
		v1106 = v1089
		v1137 = v1120
		v1168 = v1151
		v1199 = v1182
		goto b41
	}
b46:
	;
	v569, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "emit-inst!").Deref(), []vm.Value{l_550, nid_551})
	if callErr != nil {
		return nil, callErr
	}
	v573 = v569
	term_op_574 = term_op_540
	bid_575 = bid_541
	insts_576 = insts_542
	term_577 = term_543
	i_578 = i_544
	params_579 = params_545
	doseq_seq__10979_580 = doseq_seq__10979_546
	doseq_loop__10980_581 = doseq_loop__10980_547
	deferred_cond_582 = deferred_cond_548
	f_583 = f_549
	l_584 = l_550
	nid_585 = nid_551
	op_586 = op_552
	v1104 = v1090
	v1135 = v1121
	v1166 = v1152
	v1197 = v1183
	goto b48
b47:
	;
	v573 = vm.NIL
	term_op_574 = term_op_553
	bid_575 = bid_554
	insts_576 = insts_555
	term_577 = term_556
	i_578 = i_557
	params_579 = params_558
	doseq_seq__10979_580 = doseq_seq__10979_559
	doseq_loop__10980_581 = doseq_loop__10980_560
	deferred_cond_582 = deferred_cond_561
	f_583 = f_562
	l_584 = l_563
	nid_585 = nid_564
	op_586 = op_565
	v1104 = v1103
	v1135 = v1134
	v1166 = v1165
	v1197 = v1196
	goto b48
b48:
	;
	v588 = v573
	term_op_589 = term_op_574
	bid_590 = bid_575
	insts_591 = insts_576
	term_592 = term_577
	i_593 = i_578
	params_594 = params_579
	doseq_seq__10979_595 = doseq_seq__10979_580
	doseq_loop__10980_596 = doseq_loop__10980_581
	deferred_cond_597 = deferred_cond_582
	f_598 = f_583
	l_599 = l_584
	nid_600 = nid_585
	op_601 = op_586
	v1096 = v1104
	v1127 = v1135
	v1158 = v1166
	v1189 = v1197
	goto b42
b49:
	;
	v764 = term_op_679 == vm.Keyword("branch-if")
	if v764 {
		term_op_743 = term_op_679
		bid_744 = bid_680
		insts_745 = insts_681
		term_746 = term_682
		i_747 = i_683
		params_748 = params_684
		doseq_loop__10980_749 = doseq_loop__10980_685
		deferred_cond_750 = deferred_cond_686
		f_751 = f_687
		l_752 = l_688
		goto b55
	} else {
		term_op_753 = term_op_679
		bid_754 = bid_680
		insts_755 = insts_681
		term_756 = term_682
		i_757 = i_683
		params_758 = params_684
		doseq_loop__10980_759 = doseq_loop__10980_685
		deferred_cond_760 = deferred_cond_686
		f_761 = f_687
		l_762 = l_688
		goto b56
	}
b50:
	;
	v1067 = vm.NIL
	term_op_1068 = term_op_689
	bid_1069 = bid_690
	insts_1070 = insts_691
	term_1071 = term_692
	i_1072 = i_693
	params_1073 = params_694
	doseq_loop__10980_1074 = doseq_loop__10980_695
	deferred_cond_1075 = deferred_cond_696
	f_1076 = f_697
	l_1077 = l_698
	goto b51
b51:
	;
	return v1067, nil
b52:
	;
	arg__11150_723, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "zero?").Deref(), []vm.Value{term_703})
	if callErr != nil {
		return nil, callErr
	}
	arg__11155_726, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "zero?").Deref(), []vm.Value{term_703})
	if callErr != nil {
		return nil, callErr
	}
	v727, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not").Deref(), []vm.Value{arg__11155_726})
	if callErr != nil {
		return nil, callErr
	}
	v730 = v727
	term_op_731 = term_op_699
	bid_732 = bid_700
	insts_733 = insts_701
	and__x_734 = and__x_702
	term_735 = term_703
	i_736 = i_704
	params_737 = params_705
	doseq_loop__10980_738 = doseq_loop__10980_706
	deferred_cond_739 = deferred_cond_707
	f_740 = f_708
	l_741 = l_709
	goto b54
b53:
	;
	v730 = and__x_713
	term_op_731 = term_op_710
	bid_732 = bid_711
	insts_733 = insts_712
	and__x_734 = and__x_713
	term_735 = term_714
	i_736 = i_715
	params_737 = params_716
	doseq_loop__10980_738 = doseq_loop__10980_717
	deferred_cond_739 = deferred_cond_718
	f_740 = f_719
	l_741 = l_720
	goto b54
b54:
	;
	if vm.IsTruthy(v730) {
		term_op_679 = term_op_731
		bid_680 = bid_732
		insts_681 = insts_733
		term_682 = term_735
		i_683 = i_736
		params_684 = params_737
		doseq_loop__10980_685 = doseq_loop__10980_738
		deferred_cond_686 = deferred_cond_739
		f_687 = f_740
		l_688 = l_741
		goto b49
	} else {
		term_op_689 = term_op_731
		bid_690 = bid_732
		insts_691 = insts_733
		term_692 = term_735
		i_693 = i_736
		params_694 = params_737
		doseq_loop__10980_695 = doseq_loop__10980_738
		deferred_cond_696 = deferred_cond_739
		f_697 = f_740
		l_698 = l_741
		goto b50
	}
b55:
	;
	cond_aux_767, callErr = rt.InvokeValue(rt.LookupVar("ir", "aux").Deref(), []vm.Value{term_746, f_751})
	if callErr != nil {
		return nil, callErr
	}
	t_bt_769, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-true").Deref(), []vm.Value{cond_aux_767})
	if callErr != nil {
		return nil, callErr
	}
	args_771, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-args").Deref(), []vm.Value{t_bt_769})
	if callErr != nil {
		return nil, callErr
	}
	v773, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "materialize-branch-args!").Deref(), []vm.Value{l_752, args_771})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(deferred_cond_750) {
		term_op_774 = term_op_743
		bid_775 = bid_744
		insts_776 = insts_745
		term_777 = term_746
		i_778 = i_747
		params_779 = params_748
		doseq_loop__10980_780 = doseq_loop__10980_749
		deferred_cond_781 = deferred_cond_750
		f_782 = f_751
		l_783 = l_752
		cond_aux_784 = cond_aux_767
		t_bt_785 = t_bt_769
		args_786 = args_771
		goto b58
	} else {
		term_op_787 = term_op_743
		bid_788 = bid_744
		insts_789 = insts_745
		term_790 = term_746
		i_791 = i_747
		params_792 = params_748
		doseq_loop__10980_793 = doseq_loop__10980_749
		deferred_cond_794 = deferred_cond_750
		f_795 = f_751
		l_796 = l_752
		cond_aux_797 = cond_aux_767
		t_bt_798 = t_bt_769
		args_799 = args_771
		goto b59
	}
b56:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		term_op_921 = term_op_753
		bid_922 = bid_754
		insts_923 = insts_755
		term_924 = term_756
		i_925 = i_757
		params_926 = params_758
		doseq_loop__10980_927 = doseq_loop__10980_759
		deferred_cond_928 = deferred_cond_760
		f_929 = f_761
		l_930 = l_762
		goto b67
	} else {
		term_op_931 = term_op_753
		bid_932 = bid_754
		insts_933 = insts_755
		term_934 = term_756
		i_935 = i_757
		params_936 = params_758
		doseq_loop__10980_937 = doseq_loop__10980_759
		deferred_cond_938 = deferred_cond_760
		f_939 = f_761
		l_940 = l_762
		goto b68
	}
b57:
	;
	v1067 = v1053
	term_op_1068 = term_op_1054
	bid_1069 = bid_1055
	insts_1070 = insts_1056
	term_1071 = term_1057
	i_1072 = i_1058
	params_1073 = params_1059
	doseq_loop__10980_1074 = doseq_loop__10980_1060
	deferred_cond_1075 = deferred_cond_1061
	f_1076 = f_1062
	l_1077 = l_1063
	goto b51
b58:
	;
	v802, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "emit-inst!").Deref(), []vm.Value{l_783, deferred_cond_781})
	if callErr != nil {
		return nil, callErr
	}
	v904 = v802
	term_op_905 = term_op_774
	bid_906 = bid_775
	insts_907 = insts_776
	term_908 = term_777
	i_909 = i_778
	params_910 = params_779
	doseq_loop__10980_911 = doseq_loop__10980_780
	deferred_cond_912 = deferred_cond_781
	f_913 = f_782
	l_914 = l_783
	cond_aux_915 = cond_aux_784
	t_bt_916 = t_bt_785
	args_917 = args_786
	goto b60
b59:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		term_op_804 = term_op_787
		bid_805 = bid_788
		insts_806 = insts_789
		term_807 = term_790
		i_808 = i_791
		params_809 = params_792
		doseq_loop__10980_810 = doseq_loop__10980_793
		deferred_cond_811 = deferred_cond_794
		f_812 = f_795
		l_813 = l_796
		cond_aux_814 = cond_aux_797
		t_bt_815 = t_bt_798
		args_816 = args_799
		goto b61
	} else {
		term_op_817 = term_op_787
		bid_818 = bid_788
		insts_819 = insts_789
		term_820 = term_790
		i_821 = i_791
		params_822 = params_792
		doseq_loop__10980_823 = doseq_loop__10980_793
		deferred_cond_824 = deferred_cond_794
		f_825 = f_795
		l_826 = l_796
		cond_aux_827 = cond_aux_797
		t_bt_828 = t_bt_798
		args_829 = args_799
		goto b62
	}
b60:
	;
	v919, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "lower-node!").Deref(), []vm.Value{l_914, term_908})
	if callErr != nil {
		return nil, callErr
	}
	v1053 = v919
	term_op_1054 = term_op_905
	bid_1055 = bid_906
	insts_1056 = insts_907
	term_1057 = term_908
	i_1058 = i_909
	params_1059 = params_910
	doseq_loop__10980_1060 = doseq_loop__10980_911
	deferred_cond_1061 = deferred_cond_912
	f_1062 = f_913
	l_1063 = l_914
	goto b57
b61:
	;
	term_refs_833, callErr = rt.InvokeValue(rt.LookupVar("ir", "refs").Deref(), []vm.Value{term_807, f_812})
	if callErr != nil {
		return nil, callErr
	}
	v863, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "refs-at-top-last-use?").Deref(), []vm.Value{l_813, term_refs_833})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v863) {
		term_op_834 = term_op_804
		bid_835 = bid_805
		insts_836 = insts_806
		term_837 = term_807
		i_838 = i_808
		params_839 = params_809
		doseq_loop__10980_840 = doseq_loop__10980_810
		deferred_cond_841 = deferred_cond_811
		f_842 = f_812
		l_843 = l_813
		cond_aux_844 = cond_aux_814
		t_bt_845 = t_bt_815
		args_846 = args_816
		term_refs_847 = term_refs_833
		goto b64
	} else {
		term_op_848 = term_op_804
		bid_849 = bid_805
		insts_850 = insts_806
		term_851 = term_807
		i_852 = i_808
		params_853 = params_809
		doseq_loop__10980_854 = doseq_loop__10980_810
		deferred_cond_855 = deferred_cond_811
		f_856 = f_812
		l_857 = l_813
		cond_aux_858 = cond_aux_814
		t_bt_859 = t_bt_815
		args_860 = args_816
		term_refs_861 = term_refs_833
		goto b65
	}
b62:
	;
	v889 = vm.NIL
	term_op_890 = term_op_817
	bid_891 = bid_818
	insts_892 = insts_819
	term_893 = term_820
	i_894 = i_821
	params_895 = params_822
	doseq_loop__10980_896 = doseq_loop__10980_823
	deferred_cond_897 = deferred_cond_824
	f_898 = f_825
	l_899 = l_826
	cond_aux_900 = cond_aux_827
	t_bt_901 = t_bt_828
	args_902 = args_829
	goto b63
b63:
	;
	v904 = v889
	term_op_905 = term_op_890
	bid_906 = bid_891
	insts_907 = insts_892
	term_908 = term_893
	i_909 = i_894
	params_910 = params_895
	doseq_loop__10980_911 = doseq_loop__10980_896
	deferred_cond_912 = deferred_cond_897
	f_913 = f_898
	l_914 = l_899
	cond_aux_915 = cond_aux_900
	t_bt_916 = t_bt_901
	args_917 = args_902
	goto b60
b64:
	;
	v866, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "consume-refs-in-place!").Deref(), []vm.Value{l_843, term_refs_847})
	if callErr != nil {
		return nil, callErr
	}
	v871 = v866
	term_op_872 = term_op_834
	bid_873 = bid_835
	insts_874 = insts_836
	term_875 = term_837
	i_876 = i_838
	params_877 = params_839
	doseq_loop__10980_878 = doseq_loop__10980_840
	deferred_cond_879 = deferred_cond_841
	f_880 = f_842
	l_881 = l_843
	cond_aux_882 = cond_aux_844
	t_bt_883 = t_bt_845
	args_884 = args_846
	term_refs_885 = term_refs_847
	goto b66
b65:
	;
	v869, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "materialize-refs!").Deref(), []vm.Value{l_857, term_refs_861})
	if callErr != nil {
		return nil, callErr
	}
	v871 = v869
	term_op_872 = term_op_848
	bid_873 = bid_849
	insts_874 = insts_850
	term_875 = term_851
	i_876 = i_852
	params_877 = params_853
	doseq_loop__10980_878 = doseq_loop__10980_854
	deferred_cond_879 = deferred_cond_855
	f_880 = f_856
	l_881 = l_857
	cond_aux_882 = cond_aux_858
	t_bt_883 = t_bt_859
	args_884 = args_860
	term_refs_885 = term_refs_861
	goto b66
b66:
	;
	v889 = v871
	term_op_890 = term_op_872
	bid_891 = bid_873
	insts_892 = insts_874
	term_893 = term_875
	i_894 = i_876
	params_895 = params_877
	doseq_loop__10980_896 = doseq_loop__10980_878
	deferred_cond_897 = deferred_cond_879
	f_898 = f_880
	l_899 = l_881
	cond_aux_900 = cond_aux_882
	t_bt_901 = t_bt_883
	args_902 = args_884
	goto b63
b67:
	;
	term_refs_944, callErr = rt.InvokeValue(rt.LookupVar("ir", "refs").Deref(), []vm.Value{term_924, f_929})
	if callErr != nil {
		return nil, callErr
	}
	v968, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "refs-at-top-last-use?").Deref(), []vm.Value{l_930, term_refs_944})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v968) {
		term_op_945 = term_op_921
		bid_946 = bid_922
		insts_947 = insts_923
		term_948 = term_924
		i_949 = i_925
		params_950 = params_926
		doseq_loop__10980_951 = doseq_loop__10980_927
		deferred_cond_952 = deferred_cond_928
		f_953 = f_929
		l_954 = l_930
		term_refs_955 = term_refs_944
		goto b70
	} else {
		term_op_956 = term_op_921
		bid_957 = bid_922
		insts_958 = insts_923
		term_959 = term_924
		i_960 = i_925
		params_961 = params_926
		doseq_loop__10980_962 = doseq_loop__10980_927
		deferred_cond_963 = deferred_cond_928
		f_964 = f_929
		l_965 = l_930
		term_refs_966 = term_refs_944
		goto b71
	}
b68:
	;
	v1041 = vm.NIL
	term_op_1042 = term_op_931
	bid_1043 = bid_932
	insts_1044 = insts_933
	term_1045 = term_934
	i_1046 = i_935
	params_1047 = params_936
	doseq_loop__10980_1048 = doseq_loop__10980_937
	deferred_cond_1049 = deferred_cond_938
	f_1050 = f_939
	l_1051 = l_940
	goto b69
b69:
	;
	v1053 = v1041
	term_op_1054 = term_op_1042
	bid_1055 = bid_1043
	insts_1056 = insts_1044
	term_1057 = term_1045
	i_1058 = i_1046
	params_1059 = params_1047
	doseq_loop__10980_1060 = doseq_loop__10980_1048
	deferred_cond_1061 = deferred_cond_1049
	f_1062 = f_1050
	l_1063 = l_1051
	goto b57
b70:
	;
	v971, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "consume-refs-in-place!").Deref(), []vm.Value{l_954, term_refs_955})
	if callErr != nil {
		return nil, callErr
	}
	v976 = v971
	term_op_977 = term_op_945
	bid_978 = bid_946
	insts_979 = insts_947
	term_980 = term_948
	i_981 = i_949
	params_982 = params_950
	doseq_loop__10980_983 = doseq_loop__10980_951
	deferred_cond_984 = deferred_cond_952
	f_985 = f_953
	l_986 = l_954
	term_refs_987 = term_refs_955
	goto b72
b71:
	;
	v974, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "materialize-refs!").Deref(), []vm.Value{l_965, term_refs_966})
	if callErr != nil {
		return nil, callErr
	}
	v976 = v974
	term_op_977 = term_op_956
	bid_978 = bid_957
	insts_979 = insts_958
	term_980 = term_959
	i_981 = i_960
	params_982 = params_961
	doseq_loop__10980_983 = doseq_loop__10980_962
	deferred_cond_984 = deferred_cond_963
	f_985 = f_964
	l_986 = l_965
	term_refs_987 = term_refs_966
	goto b72
b72:
	;
	v1011 = term_op_977 == vm.Keyword("branch")
	if v1011 {
		term_op_988 = term_op_977
		bid_989 = bid_978
		insts_990 = insts_979
		term_991 = term_980
		i_992 = i_981
		params_993 = params_982
		doseq_loop__10980_994 = doseq_loop__10980_983
		deferred_cond_995 = deferred_cond_984
		f_996 = f_985
		l_997 = l_986
		term_refs_998 = term_refs_987
		goto b73
	} else {
		term_op_999 = term_op_977
		bid_1000 = bid_978
		insts_1001 = insts_979
		term_1002 = term_980
		i_1003 = i_981
		params_1004 = params_982
		doseq_loop__10980_1005 = doseq_loop__10980_983
		deferred_cond_1006 = deferred_cond_984
		f_1007 = f_985
		l_1008 = l_986
		term_refs_1009 = term_refs_987
		goto b74
	}
b73:
	;
	bt_1014, callErr = rt.InvokeValue(rt.LookupVar("ir", "aux").Deref(), []vm.Value{term_991, f_996})
	if callErr != nil {
		return nil, callErr
	}
	arg__11235_1016, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-args").Deref(), []vm.Value{bt_1014})
	if callErr != nil {
		return nil, callErr
	}
	arg__11241_1019, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-args").Deref(), []vm.Value{bt_1014})
	if callErr != nil {
		return nil, callErr
	}
	v1020, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "materialize-branch-args!").Deref(), []vm.Value{l_997, arg__11241_1019})
	if callErr != nil {
		return nil, callErr
	}
	v1024 = v1020
	term_op_1025 = term_op_988
	bid_1026 = bid_989
	insts_1027 = insts_990
	term_1028 = term_991
	i_1029 = i_992
	params_1030 = params_993
	doseq_loop__10980_1031 = doseq_loop__10980_994
	deferred_cond_1032 = deferred_cond_995
	f_1033 = f_996
	l_1034 = l_997
	term_refs_1035 = term_refs_998
	goto b75
b74:
	;
	v1024 = vm.NIL
	term_op_1025 = term_op_999
	bid_1026 = bid_1000
	insts_1027 = insts_1001
	term_1028 = term_1002
	i_1029 = i_1003
	params_1030 = params_1004
	doseq_loop__10980_1031 = doseq_loop__10980_1005
	deferred_cond_1032 = deferred_cond_1006
	f_1033 = f_1007
	l_1034 = l_1008
	term_refs_1035 = term_refs_1009
	goto b75
b75:
	;
	v1037, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "lower-node!").Deref(), []vm.Value{l_1034, term_1028})
	if callErr != nil {
		return nil, callErr
	}
	v1041 = v1037
	term_op_1042 = term_op_1025
	bid_1043 = bid_1026
	insts_1044 = insts_1027
	term_1045 = term_1028
	i_1046 = i_1029
	params_1047 = params_1030
	doseq_loop__10980_1048 = doseq_loop__10980_1031
	deferred_cond_1049 = deferred_cond_1032
	f_1050 = f_1033
	l_1051 = l_1034
	goto b69
}
func lower_node_BANG_(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var f_4 vm.Value
	var op_6 vm.Value
	var aux_8 vm.Value
	var v22 bool
	var l_9 vm.Value
	var nid_10 vm.Value
	var f_11 vm.Value
	var case__11247_12 vm.Value
	var op_13 vm.Value
	var aux_14 vm.Value
	var l_15 vm.Value
	var nid_16 vm.Value
	var f_17 vm.Value
	var case__11247_18 vm.Value
	var op_19 vm.Value
	var aux_20 vm.Value
	var v39 bool
	var v1057 vm.Value
	var l_1058 vm.Value
	var nid_1059 vm.Value
	var f_1060 vm.Value
	var case__11247_1061 vm.Value
	var op_1062 vm.Value
	var aux_1063 vm.Value
	var l_26 vm.Value
	var nid_27 vm.Value
	var f_28 vm.Value
	var case__11247_29 vm.Value
	var op_30 vm.Value
	var aux_31 vm.Value
	var arg__11268_42 vm.Value
	var arg__11272_44 vm.Value
	var arg__11277_47 vm.Value
	var arg__11281_49 vm.Value
	var idx_50 vm.Value
	var v54 vm.Value
	var v58 vm.Value
	var v60 vm.Value
	var l_32 vm.Value
	var nid_33 vm.Value
	var f_34 vm.Value
	var case__11247_35 vm.Value
	var op_36 vm.Value
	var aux_37 vm.Value
	var v75 bool
	var v1049 vm.Value
	var l_1050 vm.Value
	var nid_1051 vm.Value
	var f_1052 vm.Value
	var case__11247_1053 vm.Value
	var op_1054 vm.Value
	var aux_1055 vm.Value
	var l_62 vm.Value
	var nid_63 vm.Value
	var f_64 vm.Value
	var case__11247_65 vm.Value
	var op_66 vm.Value
	var aux_67 vm.Value
	var arg__11307_79 vm.Value
	var arg__11315_83 vm.Value
	var v84 vm.Value
	var v88 vm.Value
	var v90 vm.Value
	var l_68 vm.Value
	var nid_69 vm.Value
	var f_70 vm.Value
	var case__11247_71 vm.Value
	var op_72 vm.Value
	var aux_73 vm.Value
	var v105 bool
	var v1041 vm.Value
	var l_1042 vm.Value
	var nid_1043 vm.Value
	var f_1044 vm.Value
	var case__11247_1045 vm.Value
	var op_1046 vm.Value
	var aux_1047 vm.Value
	var l_92 vm.Value
	var nid_93 vm.Value
	var f_94 vm.Value
	var case__11247_95 vm.Value
	var op_96 vm.Value
	var aux_97 vm.Value
	var arg__11329_108 vm.Value
	var arg__11335_111 vm.Value
	var idx_112 vm.Value
	var v116 vm.Value
	var v120 vm.Value
	var v122 vm.Value
	var l_98 vm.Value
	var nid_99 vm.Value
	var f_100 vm.Value
	var case__11247_101 vm.Value
	var op_102 vm.Value
	var aux_103 vm.Value
	var or__x_137 bool
	var v1033 vm.Value
	var l_1034 vm.Value
	var nid_1035 vm.Value
	var f_1036 vm.Value
	var case__11247_1037 vm.Value
	var op_1038 vm.Value
	var aux_1039 vm.Value
	var l_124 vm.Value
	var nid_125 vm.Value
	var f_126 vm.Value
	var case__11247_127 vm.Value
	var op_128 vm.Value
	var aux_129 vm.Value
	var v329 vm.Value
	var v333 vm.Value
	var l_130 vm.Value
	var nid_131 vm.Value
	var f_132 vm.Value
	var case__11247_133 vm.Value
	var op_134 vm.Value
	var aux_135 vm.Value
	var or__x_348 bool
	var v1025 vm.Value
	var l_1026 vm.Value
	var nid_1027 vm.Value
	var f_1028 vm.Value
	var case__11247_1029 vm.Value
	var op_1030 vm.Value
	var aux_1031 vm.Value
	var l_138 vm.Value
	var nid_139 vm.Value
	var f_140 vm.Value
	var case__11247_141 vm.Value
	var op_142 vm.Value
	var aux_143 vm.Value
	var or__x_144 bool
	var l_145 vm.Value
	var nid_146 vm.Value
	var f_147 vm.Value
	var case__11247_148 vm.Value
	var op_149 vm.Value
	var aux_150 vm.Value
	var or__x_151 bool
	var or__x_155 bool
	var v319 bool
	var l_320 vm.Value
	var nid_321 vm.Value
	var f_322 vm.Value
	var case__11247_323 vm.Value
	var op_324 vm.Value
	var aux_325 vm.Value
	var or__x_326 vm.Value
	var l_156 vm.Value
	var nid_157 vm.Value
	var f_158 vm.Value
	var case__11247_159 vm.Value
	var op_160 vm.Value
	var aux_161 vm.Value
	var or__x_162 bool
	var l_163 vm.Value
	var nid_164 vm.Value
	var f_165 vm.Value
	var case__11247_166 vm.Value
	var op_167 vm.Value
	var aux_168 vm.Value
	var or__x_169 bool
	var or__x_173 bool
	var v310 bool
	var l_311 vm.Value
	var nid_312 vm.Value
	var f_313 vm.Value
	var case__11247_314 vm.Value
	var op_315 vm.Value
	var aux_316 vm.Value
	var or__x_317 vm.Value
	var l_174 vm.Value
	var nid_175 vm.Value
	var f_176 vm.Value
	var case__11247_177 vm.Value
	var op_178 vm.Value
	var aux_179 vm.Value
	var or__x_180 bool
	var l_181 vm.Value
	var nid_182 vm.Value
	var f_183 vm.Value
	var case__11247_184 vm.Value
	var op_185 vm.Value
	var aux_186 vm.Value
	var or__x_187 bool
	var or__x_191 bool
	var v301 bool
	var l_302 vm.Value
	var nid_303 vm.Value
	var f_304 vm.Value
	var case__11247_305 vm.Value
	var op_306 vm.Value
	var aux_307 vm.Value
	var or__x_308 vm.Value
	var l_192 vm.Value
	var nid_193 vm.Value
	var f_194 vm.Value
	var case__11247_195 vm.Value
	var op_196 vm.Value
	var aux_197 vm.Value
	var or__x_198 bool
	var l_199 vm.Value
	var nid_200 vm.Value
	var f_201 vm.Value
	var case__11247_202 vm.Value
	var op_203 vm.Value
	var aux_204 vm.Value
	var or__x_205 bool
	var or__x_209 bool
	var v292 bool
	var l_293 vm.Value
	var nid_294 vm.Value
	var f_295 vm.Value
	var case__11247_296 vm.Value
	var op_297 vm.Value
	var aux_298 vm.Value
	var or__x_299 vm.Value
	var l_210 vm.Value
	var nid_211 vm.Value
	var f_212 vm.Value
	var case__11247_213 vm.Value
	var op_214 vm.Value
	var aux_215 vm.Value
	var or__x_216 bool
	var l_217 vm.Value
	var nid_218 vm.Value
	var f_219 vm.Value
	var case__11247_220 vm.Value
	var op_221 vm.Value
	var aux_222 vm.Value
	var or__x_223 bool
	var or__x_227 bool
	var v283 bool
	var l_284 vm.Value
	var nid_285 vm.Value
	var f_286 vm.Value
	var case__11247_287 vm.Value
	var op_288 vm.Value
	var aux_289 vm.Value
	var or__x_290 vm.Value
	var l_228 vm.Value
	var nid_229 vm.Value
	var f_230 vm.Value
	var case__11247_231 vm.Value
	var op_232 vm.Value
	var aux_233 vm.Value
	var or__x_234 bool
	var l_235 vm.Value
	var nid_236 vm.Value
	var f_237 vm.Value
	var case__11247_238 vm.Value
	var op_239 vm.Value
	var aux_240 vm.Value
	var or__x_241 bool
	var or__x_245 bool
	var v274 bool
	var l_275 vm.Value
	var nid_276 vm.Value
	var f_277 vm.Value
	var case__11247_278 vm.Value
	var op_279 vm.Value
	var aux_280 vm.Value
	var or__x_281 vm.Value
	var l_246 vm.Value
	var nid_247 vm.Value
	var f_248 vm.Value
	var case__11247_249 vm.Value
	var op_250 vm.Value
	var aux_251 vm.Value
	var or__x_252 bool
	var l_253 vm.Value
	var nid_254 vm.Value
	var f_255 vm.Value
	var case__11247_256 vm.Value
	var op_257 vm.Value
	var aux_258 vm.Value
	var or__x_259 bool
	var v263 bool
	var v265 bool
	var l_266 vm.Value
	var nid_267 vm.Value
	var f_268 vm.Value
	var case__11247_269 vm.Value
	var op_270 vm.Value
	var aux_271 vm.Value
	var or__x_272 vm.Value
	var l_335 vm.Value
	var nid_336 vm.Value
	var f_337 vm.Value
	var case__11247_338 vm.Value
	var op_339 vm.Value
	var aux_340 vm.Value
	var v378 vm.Value
	var l_341 vm.Value
	var nid_342 vm.Value
	var f_343 vm.Value
	var case__11247_344 vm.Value
	var op_345 vm.Value
	var aux_346 vm.Value
	var v393 bool
	var v1017 vm.Value
	var l_1018 vm.Value
	var nid_1019 vm.Value
	var f_1020 vm.Value
	var case__11247_1021 vm.Value
	var op_1022 vm.Value
	var aux_1023 vm.Value
	var l_349 vm.Value
	var nid_350 vm.Value
	var f_351 vm.Value
	var case__11247_352 vm.Value
	var op_353 vm.Value
	var aux_354 vm.Value
	var or__x_355 bool
	var l_356 vm.Value
	var nid_357 vm.Value
	var f_358 vm.Value
	var case__11247_359 vm.Value
	var op_360 vm.Value
	var aux_361 vm.Value
	var or__x_362 bool
	var v366 bool
	var v368 bool
	var l_369 vm.Value
	var nid_370 vm.Value
	var f_371 vm.Value
	var case__11247_372 vm.Value
	var op_373 vm.Value
	var aux_374 vm.Value
	var or__x_375 vm.Value
	var l_380 vm.Value
	var nid_381 vm.Value
	var f_382 vm.Value
	var case__11247_383 vm.Value
	var op_384 vm.Value
	var aux_385 vm.Value
	var arg__11401_397 vm.Value
	var arg__11409_401 vm.Value
	var v402 vm.Value
	var arg__11414_404 vm.Value
	var arg__11421_409 vm.Value
	var arg__11422_411 vm.Value
	var v412 vm.Value
	var l_386 vm.Value
	var nid_387 vm.Value
	var f_388 vm.Value
	var case__11247_389 vm.Value
	var op_390 vm.Value
	var aux_391 vm.Value
	var v427 bool
	var v1009 vm.Value
	var l_1010 vm.Value
	var nid_1011 vm.Value
	var f_1012 vm.Value
	var case__11247_1013 vm.Value
	var op_1014 vm.Value
	var aux_1015 vm.Value
	var l_414 vm.Value
	var nid_415 vm.Value
	var f_416 vm.Value
	var case__11247_417 vm.Value
	var op_418 vm.Value
	var aux_419 vm.Value
	var v432 vm.Value
	var v436 vm.Value
	var l_420 vm.Value
	var nid_421 vm.Value
	var f_422 vm.Value
	var case__11247_423 vm.Value
	var op_424 vm.Value
	var aux_425 vm.Value
	var v451 bool
	var v1001 vm.Value
	var l_1002 vm.Value
	var nid_1003 vm.Value
	var f_1004 vm.Value
	var case__11247_1005 vm.Value
	var op_1006 vm.Value
	var aux_1007 vm.Value
	var l_438 vm.Value
	var nid_439 vm.Value
	var f_440 vm.Value
	var case__11247_441 vm.Value
	var op_442 vm.Value
	var aux_443 vm.Value
	var arg__11445_455 vm.Value
	var arg__11453_459 vm.Value
	var v460 vm.Value
	var l_444 vm.Value
	var nid_445 vm.Value
	var f_446 vm.Value
	var case__11247_447 vm.Value
	var op_448 vm.Value
	var aux_449 vm.Value
	var v475 bool
	var v993 vm.Value
	var l_994 vm.Value
	var nid_995 vm.Value
	var f_996 vm.Value
	var case__11247_997 vm.Value
	var op_998 vm.Value
	var aux_999 vm.Value
	var l_462 vm.Value
	var nid_463 vm.Value
	var f_464 vm.Value
	var case__11247_465 vm.Value
	var op_466 vm.Value
	var aux_467 vm.Value
	var v480 vm.Value
	var l_468 vm.Value
	var nid_469 vm.Value
	var f_470 vm.Value
	var case__11247_471 vm.Value
	var op_472 vm.Value
	var aux_473 vm.Value
	var v495 bool
	var v985 vm.Value
	var l_986 vm.Value
	var nid_987 vm.Value
	var f_988 vm.Value
	var case__11247_989 vm.Value
	var op_990 vm.Value
	var aux_991 vm.Value
	var l_482 vm.Value
	var nid_483 vm.Value
	var f_484 vm.Value
	var case__11247_485 vm.Value
	var op_486 vm.Value
	var aux_487 vm.Value
	var args_498 vm.Value
	var target_500 vm.Value
	var argc_502 vm.Value
	var arg__11479_504 vm.Value
	var arg__11486_507 vm.Value
	var target_params_508 vm.Value
	var cur_sp_510 vm.Value
	var drop_count_511 vm.Value
	var arg__11497_514 vm.Value
	var arg__11498_515 vm.Value
	var arg__11505_519 vm.Value
	var arg__11506_520 vm.Value
	var cur_junk_521 vm.Value
	var v551 vm.Value
	var l_488 vm.Value
	var nid_489 vm.Value
	var f_490 vm.Value
	var case__11247_491 vm.Value
	var op_492 vm.Value
	var aux_493 vm.Value
	var v689 bool
	var v977 vm.Value
	var l_978 vm.Value
	var nid_979 vm.Value
	var f_980 vm.Value
	var case__11247_981 vm.Value
	var op_982 vm.Value
	var aux_983 vm.Value
	var l_522 vm.Value
	var nid_523 vm.Value
	var f_524 vm.Value
	var case__11247_525 vm.Value
	var op_526 vm.Value
	var bt_527 vm.Value
	var aux_528 vm.Value
	var args_529 vm.Value
	var target_530 vm.Value
	var argc_531 vm.Value
	var target_params_532 vm.Value
	var cur_sp_533 vm.Value
	var drop_count_534 vm.Value
	var cur_junk_535 vm.Value
	var ignore_553 vm.Value
	var arg__11515_555 vm.Value
	var arg__11523_558 vm.Value
	var off_ip_559 vm.Value
	var v561 vm.Value
	var arg__11536_564 vm.Value
	var arg__11543_570 vm.Value
	var arg__11556_580 vm.Value
	var v581 vm.Value
	var v583 vm.Value
	var l_536 vm.Value
	var nid_537 vm.Value
	var f_538 vm.Value
	var case__11247_539 vm.Value
	var op_540 vm.Value
	var bt_541 vm.Value
	var aux_542 vm.Value
	var args_543 vm.Value
	var target_544 vm.Value
	var argc_545 vm.Value
	var target_params_546 vm.Value
	var cur_sp_547 vm.Value
	var drop_count_548 vm.Value
	var cur_junk_549 vm.Value
	var v660 vm.Value
	var l_661 vm.Value
	var nid_662 vm.Value
	var f_663 vm.Value
	var case__11247_664 vm.Value
	var op_665 vm.Value
	var bt_666 vm.Value
	var aux_667 vm.Value
	var args_668 vm.Value
	var target_669 vm.Value
	var argc_670 vm.Value
	var target_params_671 vm.Value
	var cur_sp_672 vm.Value
	var drop_count_673 vm.Value
	var cur_junk_674 vm.Value
	var l_585 vm.Value
	var nid_586 vm.Value
	var f_587 vm.Value
	var case__11247_588 vm.Value
	var op_589 vm.Value
	var bt_590 vm.Value
	var aux_591 vm.Value
	var args_592 vm.Value
	var target_593 vm.Value
	var argc_594 vm.Value
	var target_params_595 vm.Value
	var cur_sp_596 vm.Value
	var drop_count_597 vm.Value
	var cur_junk_598 vm.Value
	var arg_ip_618 vm.Value
	var v620 vm.Value
	var arg__11580_623 vm.Value
	var arg__11587_629 vm.Value
	var arg__11600_639 vm.Value
	var v640 vm.Value
	var l_599 vm.Value
	var nid_600 vm.Value
	var f_601 vm.Value
	var case__11247_602 vm.Value
	var op_603 vm.Value
	var bt_604 vm.Value
	var aux_605 vm.Value
	var args_606 vm.Value
	var target_607 vm.Value
	var argc_608 vm.Value
	var target_params_609 vm.Value
	var cur_sp_610 vm.Value
	var drop_count_611 vm.Value
	var cur_junk_612 vm.Value
	var v644 vm.Value
	var l_645 vm.Value
	var nid_646 vm.Value
	var f_647 vm.Value
	var case__11247_648 vm.Value
	var op_649 vm.Value
	var bt_650 vm.Value
	var aux_651 vm.Value
	var args_652 vm.Value
	var target_653 vm.Value
	var argc_654 vm.Value
	var target_params_655 vm.Value
	var cur_sp_656 vm.Value
	var drop_count_657 vm.Value
	var cur_junk_658 vm.Value
	var l_676 vm.Value
	var nid_677 vm.Value
	var f_678 vm.Value
	var case__11247_679 vm.Value
	var op_680 vm.Value
	var aux_681 vm.Value
	var ft_692 vm.Value
	var tt_694 vm.Value
	var ft_target_696 vm.Value
	var tt_target_698 vm.Value
	var arg_ip_702 vm.Value
	var arg__11627_705 vm.Value
	var arg__11628_706 vm.Value
	var arg__11635_710 vm.Value
	var arg__11636_711 vm.Value
	var cur_junk_712 vm.Value
	var arg__11640_714 vm.Value
	var arg__11642_715 vm.Value
	var arg__11649_718 vm.Value
	var arg__11656_721 vm.Value
	var arg__11657_722 vm.Value
	var v723 vm.Value
	var target_junk_724 vm.Value
	var v726 vm.Value
	var v728 vm.Value
	var arg__11676_731 vm.Value
	var arg__11683_737 vm.Value
	var arg__11696_747 vm.Value
	var v748 vm.Value
	var v752 vm.Value
	var my_block_754 vm.Value
	var next_block_id_755 vm.Value
	var v789 vm.Value
	var l_682 vm.Value
	var nid_683 vm.Value
	var f_684 vm.Value
	var case__11247_685 vm.Value
	var op_686 vm.Value
	var aux_687 vm.Value
	var v849 bool
	var v969 vm.Value
	var l_970 vm.Value
	var nid_971 vm.Value
	var f_972 vm.Value
	var case__11247_973 vm.Value
	var op_974 vm.Value
	var aux_975 vm.Value
	var l_756 vm.Value
	var nid_757 vm.Value
	var f_758 vm.Value
	var case__11247_759 vm.Value
	var op_760 vm.Value
	var ct_761 vm.Value
	var aux_762 vm.Value
	var ft_763 vm.Value
	var tt_764 vm.Value
	var ft_target_765 vm.Value
	var tt_target_766 vm.Value
	var arg_ip_767 vm.Value
	var cur_junk_768 vm.Value
	var target_junk_769 vm.Value
	var my_block_770 vm.Value
	var next_block_id_771 vm.Value
	var arg_ip2_794 vm.Value
	var arg__11724_797 vm.Value
	var arg__11731_803 vm.Value
	var arg__11744_813 vm.Value
	var v814 vm.Value
	var l_772 vm.Value
	var nid_773 vm.Value
	var f_774 vm.Value
	var case__11247_775 vm.Value
	var op_776 vm.Value
	var ct_777 vm.Value
	var aux_778 vm.Value
	var ft_779 vm.Value
	var tt_780 vm.Value
	var ft_target_781 vm.Value
	var tt_target_782 vm.Value
	var arg_ip_783 vm.Value
	var cur_junk_784 vm.Value
	var target_junk_785 vm.Value
	var my_block_786 vm.Value
	var next_block_id_787 vm.Value
	var v818 vm.Value
	var l_819 vm.Value
	var nid_820 vm.Value
	var f_821 vm.Value
	var case__11247_822 vm.Value
	var op_823 vm.Value
	var ct_824 vm.Value
	var aux_825 vm.Value
	var ft_826 vm.Value
	var tt_827 vm.Value
	var ft_target_828 vm.Value
	var tt_target_829 vm.Value
	var arg_ip_830 vm.Value
	var cur_junk_831 vm.Value
	var target_junk_832 vm.Value
	var my_block_833 vm.Value
	var next_block_id_834 vm.Value
	var l_836 vm.Value
	var nid_837 vm.Value
	var f_838 vm.Value
	var case__11247_839 vm.Value
	var op_840 vm.Value
	var aux_841 vm.Value
	var arg__11753_853 vm.Value
	var arg__11761_857 vm.Value
	var v858 vm.Value
	var v862 vm.Value
	var v864 vm.Value
	var l_842 vm.Value
	var nid_843 vm.Value
	var f_844 vm.Value
	var case__11247_845 vm.Value
	var op_846 vm.Value
	var aux_847 vm.Value
	var v879 bool
	var v961 vm.Value
	var l_962 vm.Value
	var nid_963 vm.Value
	var f_964 vm.Value
	var case__11247_965 vm.Value
	var op_966 vm.Value
	var aux_967 vm.Value
	var l_866 vm.Value
	var nid_867 vm.Value
	var f_868 vm.Value
	var case__11247_869 vm.Value
	var op_870 vm.Value
	var aux_871 vm.Value
	var v884 vm.Value
	var l_872 vm.Value
	var nid_873 vm.Value
	var f_874 vm.Value
	var case__11247_875 vm.Value
	var op_876 vm.Value
	var aux_877 vm.Value
	var v899 bool
	var v953 vm.Value
	var l_954 vm.Value
	var nid_955 vm.Value
	var f_956 vm.Value
	var case__11247_957 vm.Value
	var op_958 vm.Value
	var aux_959 vm.Value
	var l_886 vm.Value
	var nid_887 vm.Value
	var f_888 vm.Value
	var case__11247_889 vm.Value
	var op_890 vm.Value
	var aux_891 vm.Value
	var v904 vm.Value
	var v908 vm.Value
	var l_892 vm.Value
	var nid_893 vm.Value
	var f_894 vm.Value
	var case__11247_895 vm.Value
	var op_896 vm.Value
	var aux_897 vm.Value
	var v945 vm.Value
	var l_946 vm.Value
	var nid_947 vm.Value
	var f_948 vm.Value
	var case__11247_949 vm.Value
	var op_950 vm.Value
	var aux_951 vm.Value
	var l_910 vm.Value
	var nid_911 vm.Value
	var f_912 vm.Value
	var case__11247_913 vm.Value
	var op_914 vm.Value
	var aux_915 vm.Value
	var arg__11798_927 vm.Value
	var arg__11805_932 vm.Value
	var v933 vm.Value
	var l_916 vm.Value
	var nid_917 vm.Value
	var f_918 vm.Value
	var case__11247_919 vm.Value
	var op_920 vm.Value
	var aux_921 vm.Value
	var v937 vm.Value
	var l_938 vm.Value
	var nid_939 vm.Value
	var f_940 vm.Value
	var case__11247_941 vm.Value
	var op_942 vm.Value
	var aux_943 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = f_4, op_6, aux_8, v22, l_9, nid_10, f_11, case__11247_12, op_13, aux_14, l_15, nid_16, f_17, case__11247_18, op_19, aux_20, v39, v1057, l_1058, nid_1059, f_1060, case__11247_1061, op_1062, aux_1063, l_26, nid_27, f_28, case__11247_29, op_30, aux_31, arg__11268_42, arg__11272_44, arg__11277_47, arg__11281_49, idx_50, v54, v58, v60, l_32, nid_33, f_34, case__11247_35, op_36, aux_37, v75, v1049, l_1050, nid_1051, f_1052, case__11247_1053, op_1054, aux_1055, l_62, nid_63, f_64, case__11247_65, op_66, aux_67, arg__11307_79, arg__11315_83, v84, v88, v90, l_68, nid_69, f_70, case__11247_71, op_72, aux_73, v105, v1041, l_1042, nid_1043, f_1044, case__11247_1045, op_1046, aux_1047, l_92, nid_93, f_94, case__11247_95, op_96, aux_97, arg__11329_108, arg__11335_111, idx_112, v116, v120, v122, l_98, nid_99, f_100, case__11247_101, op_102, aux_103, or__x_137, v1033, l_1034, nid_1035, f_1036, case__11247_1037, op_1038, aux_1039, l_124, nid_125, f_126, case__11247_127, op_128, aux_129, v329, v333, l_130, nid_131, f_132, case__11247_133, op_134, aux_135, or__x_348, v1025, l_1026, nid_1027, f_1028, case__11247_1029, op_1030, aux_1031, l_138, nid_139, f_140, case__11247_141, op_142, aux_143, or__x_144, l_145, nid_146, f_147, case__11247_148, op_149, aux_150, or__x_151, or__x_155, v319, l_320, nid_321, f_322, case__11247_323, op_324, aux_325, or__x_326, l_156, nid_157, f_158, case__11247_159, op_160, aux_161, or__x_162, l_163, nid_164, f_165, case__11247_166, op_167, aux_168, or__x_169, or__x_173, v310, l_311, nid_312, f_313, case__11247_314, op_315, aux_316, or__x_317, l_174, nid_175, f_176, case__11247_177, op_178, aux_179, or__x_180, l_181, nid_182, f_183, case__11247_184, op_185, aux_186, or__x_187, or__x_191, v301, l_302, nid_303, f_304, case__11247_305, op_306, aux_307, or__x_308, l_192, nid_193, f_194, case__11247_195, op_196, aux_197, or__x_198, l_199, nid_200, f_201, case__11247_202, op_203, aux_204, or__x_205, or__x_209, v292, l_293, nid_294, f_295, case__11247_296, op_297, aux_298, or__x_299, l_210, nid_211, f_212, case__11247_213, op_214, aux_215, or__x_216, l_217, nid_218, f_219, case__11247_220, op_221, aux_222, or__x_223, or__x_227, v283, l_284, nid_285, f_286, case__11247_287, op_288, aux_289, or__x_290, l_228, nid_229, f_230, case__11247_231, op_232, aux_233, or__x_234, l_235, nid_236, f_237, case__11247_238, op_239, aux_240, or__x_241, or__x_245, v274, l_275, nid_276, f_277, case__11247_278, op_279, aux_280, or__x_281, l_246, nid_247, f_248, case__11247_249, op_250, aux_251, or__x_252, l_253, nid_254, f_255, case__11247_256, op_257, aux_258, or__x_259, v263, v265, l_266, nid_267, f_268, case__11247_269, op_270, aux_271, or__x_272, l_335, nid_336, f_337, case__11247_338, op_339, aux_340, v378, l_341, nid_342, f_343, case__11247_344, op_345, aux_346, v393, v1017, l_1018, nid_1019, f_1020, case__11247_1021, op_1022, aux_1023, l_349, nid_350, f_351, case__11247_352, op_353, aux_354, or__x_355, l_356, nid_357, f_358, case__11247_359, op_360, aux_361, or__x_362, v366, v368, l_369, nid_370, f_371, case__11247_372, op_373, aux_374, or__x_375, l_380, nid_381, f_382, case__11247_383, op_384, aux_385, arg__11401_397, arg__11409_401, v402, arg__11414_404, arg__11421_409, arg__11422_411, v412, l_386, nid_387, f_388, case__11247_389, op_390, aux_391, v427, v1009, l_1010, nid_1011, f_1012, case__11247_1013, op_1014, aux_1015, l_414, nid_415, f_416, case__11247_417, op_418, aux_419, v432, v436, l_420, nid_421, f_422, case__11247_423, op_424, aux_425, v451, v1001, l_1002, nid_1003, f_1004, case__11247_1005, op_1006, aux_1007, l_438, nid_439, f_440, case__11247_441, op_442, aux_443, arg__11445_455, arg__11453_459, v460, l_444, nid_445, f_446, case__11247_447, op_448, aux_449, v475, v993, l_994, nid_995, f_996, case__11247_997, op_998, aux_999, l_462, nid_463, f_464, case__11247_465, op_466, aux_467, v480, l_468, nid_469, f_470, case__11247_471, op_472, aux_473, v495, v985, l_986, nid_987, f_988, case__11247_989, op_990, aux_991, l_482, nid_483, f_484, case__11247_485, op_486, aux_487, args_498, target_500, argc_502, arg__11479_504, arg__11486_507, target_params_508, cur_sp_510, drop_count_511, arg__11497_514, arg__11498_515, arg__11505_519, arg__11506_520, cur_junk_521, v551, l_488, nid_489, f_490, case__11247_491, op_492, aux_493, v689, v977, l_978, nid_979, f_980, case__11247_981, op_982, aux_983, l_522, nid_523, f_524, case__11247_525, op_526, bt_527, aux_528, args_529, target_530, argc_531, target_params_532, cur_sp_533, drop_count_534, cur_junk_535, ignore_553, arg__11515_555, arg__11523_558, off_ip_559, v561, arg__11536_564, arg__11543_570, arg__11556_580, v581, v583, l_536, nid_537, f_538, case__11247_539, op_540, bt_541, aux_542, args_543, target_544, argc_545, target_params_546, cur_sp_547, drop_count_548, cur_junk_549, v660, l_661, nid_662, f_663, case__11247_664, op_665, bt_666, aux_667, args_668, target_669, argc_670, target_params_671, cur_sp_672, drop_count_673, cur_junk_674, l_585, nid_586, f_587, case__11247_588, op_589, bt_590, aux_591, args_592, target_593, argc_594, target_params_595, cur_sp_596, drop_count_597, cur_junk_598, arg_ip_618, v620, arg__11580_623, arg__11587_629, arg__11600_639, v640, l_599, nid_600, f_601, case__11247_602, op_603, bt_604, aux_605, args_606, target_607, argc_608, target_params_609, cur_sp_610, drop_count_611, cur_junk_612, v644, l_645, nid_646, f_647, case__11247_648, op_649, bt_650, aux_651, args_652, target_653, argc_654, target_params_655, cur_sp_656, drop_count_657, cur_junk_658, l_676, nid_677, f_678, case__11247_679, op_680, aux_681, ft_692, tt_694, ft_target_696, tt_target_698, arg_ip_702, arg__11627_705, arg__11628_706, arg__11635_710, arg__11636_711, cur_junk_712, arg__11640_714, arg__11642_715, arg__11649_718, arg__11656_721, arg__11657_722, v723, target_junk_724, v726, v728, arg__11676_731, arg__11683_737, arg__11696_747, v748, v752, my_block_754, next_block_id_755, v789, l_682, nid_683, f_684, case__11247_685, op_686, aux_687, v849, v969, l_970, nid_971, f_972, case__11247_973, op_974, aux_975, l_756, nid_757, f_758, case__11247_759, op_760, ct_761, aux_762, ft_763, tt_764, ft_target_765, tt_target_766, arg_ip_767, cur_junk_768, target_junk_769, my_block_770, next_block_id_771, arg_ip2_794, arg__11724_797, arg__11731_803, arg__11744_813, v814, l_772, nid_773, f_774, case__11247_775, op_776, ct_777, aux_778, ft_779, tt_780, ft_target_781, tt_target_782, arg_ip_783, cur_junk_784, target_junk_785, my_block_786, next_block_id_787, v818, l_819, nid_820, f_821, case__11247_822, op_823, ct_824, aux_825, ft_826, tt_827, ft_target_828, tt_target_829, arg_ip_830, cur_junk_831, target_junk_832, my_block_833, next_block_id_834, l_836, nid_837, f_838, case__11247_839, op_840, aux_841, arg__11753_853, arg__11761_857, v858, v862, v864, l_842, nid_843, f_844, case__11247_845, op_846, aux_847, v879, v961, l_962, nid_963, f_964, case__11247_965, op_966, aux_967, l_866, nid_867, f_868, case__11247_869, op_870, aux_871, v884, l_872, nid_873, f_874, case__11247_875, op_876, aux_877, v899, v953, l_954, nid_955, f_956, case__11247_957, op_958, aux_959, l_886, nid_887, f_888, case__11247_889, op_890, aux_891, v904, v908, l_892, nid_893, f_894, case__11247_895, op_896, aux_897, v945, l_946, nid_947, f_948, case__11247_949, op_950, aux_951, l_910, nid_911, f_912, case__11247_913, op_914, aux_915, arg__11798_927, arg__11805_932, v933, l_916, nid_917, f_918, case__11247_919, op_920, aux_921, v937, l_938, nid_939, f_940, case__11247_941, op_942, aux_943
	f_4, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "f-of").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	op_6, callErr = rt.InvokeValue(rt.LookupVar("ir", "op").Deref(), []vm.Value{arg1, f_4})
	if callErr != nil {
		return nil, callErr
	}
	aux_8, callErr = rt.InvokeValue(rt.LookupVar("ir", "aux").Deref(), []vm.Value{arg1, f_4})
	if callErr != nil {
		return nil, callErr
	}
	v22 = op_6 == vm.Keyword("block-arg")
	if v22 {
		l_9 = arg0
		nid_10 = arg1
		f_11 = f_4
		case__11247_12 = op_6
		op_13 = op_6
		aux_14 = aux_8
		goto b1
	} else {
		l_15 = arg0
		nid_16 = arg1
		f_17 = f_4
		case__11247_18 = op_6
		op_19 = op_6
		aux_20 = aux_8
		goto b2
	}
b1:
	;
	v1057 = vm.NIL
	l_1058 = l_9
	nid_1059 = nid_10
	f_1060 = f_11
	case__11247_1061 = case__11247_12
	op_1062 = op_13
	aux_1063 = aux_14
	goto b3
b2:
	;
	v39 = case__11247_18 == vm.Keyword("const")
	if v39 {
		l_26 = l_15
		nid_27 = nid_16
		f_28 = f_17
		case__11247_29 = case__11247_18
		op_30 = op_19
		aux_31 = aux_20
		goto b4
	} else {
		l_32 = l_15
		nid_33 = nid_16
		f_34 = f_17
		case__11247_35 = case__11247_18
		op_36 = op_19
		aux_37 = aux_20
		goto b5
	}
b3:
	;
	return v1057, nil
b4:
	;
	arg__11268_42, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "chunk-of").Deref(), []vm.Value{l_26})
	if callErr != nil {
		return nil, callErr
	}
	arg__11272_44, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "template-value").Deref(), []vm.Value{aux_31})
	if callErr != nil {
		return nil, callErr
	}
	arg__11277_47, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "chunk-of").Deref(), []vm.Value{l_26})
	if callErr != nil {
		return nil, callErr
	}
	arg__11281_49, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "template-value").Deref(), []vm.Value{aux_31})
	if callErr != nil {
		return nil, callErr
	}
	idx_50, callErr = rt.InvokeValue(rt.LookupVar("ir", "chunk-intern-const").Deref(), []vm.Value{arg__11277_47, arg__11281_49})
	if callErr != nil {
		return nil, callErr
	}
	v54, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "emit-with-arg!").Deref(), []vm.Value{l_26, nid_27, vm.Keyword("const"), idx_50})
	if callErr != nil {
		return nil, callErr
	}
	v58, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "bump-stack-sp!").Deref(), []vm.Value{l_26, vm.Int(1)})
	if callErr != nil {
		return nil, callErr
	}
	v60, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "bump-max-stack!").Deref(), []vm.Value{l_26})
	if callErr != nil {
		return nil, callErr
	}
	v1049 = v60
	l_1050 = l_26
	nid_1051 = nid_27
	f_1052 = f_28
	case__11247_1053 = case__11247_29
	op_1054 = op_30
	aux_1055 = aux_31
	goto b6
b5:
	;
	v75 = case__11247_35 == vm.Keyword("load-arg")
	if v75 {
		l_62 = l_32
		nid_63 = nid_33
		f_64 = f_34
		case__11247_65 = case__11247_35
		op_66 = op_36
		aux_67 = aux_37
		goto b7
	} else {
		l_68 = l_32
		nid_69 = nid_33
		f_70 = f_34
		case__11247_71 = case__11247_35
		op_72 = op_36
		aux_73 = aux_37
		goto b8
	}
b6:
	;
	v1057 = v1049
	l_1058 = l_1050
	nid_1059 = nid_1051
	f_1060 = f_1052
	case__11247_1061 = case__11247_1053
	op_1062 = op_1054
	aux_1063 = aux_1055
	goto b3
b7:
	;
	arg__11307_79, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "int").Deref(), []vm.Value{aux_67})
	if callErr != nil {
		return nil, callErr
	}
	arg__11315_83, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "int").Deref(), []vm.Value{aux_67})
	if callErr != nil {
		return nil, callErr
	}
	v84, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "emit-with-arg!").Deref(), []vm.Value{l_62, nid_63, vm.Keyword("load-arg"), arg__11315_83})
	if callErr != nil {
		return nil, callErr
	}
	v88, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "bump-stack-sp!").Deref(), []vm.Value{l_62, vm.Int(1)})
	if callErr != nil {
		return nil, callErr
	}
	v90, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "bump-max-stack!").Deref(), []vm.Value{l_62})
	if callErr != nil {
		return nil, callErr
	}
	v1041 = v90
	l_1042 = l_62
	nid_1043 = nid_63
	f_1044 = f_64
	case__11247_1045 = case__11247_65
	op_1046 = op_66
	aux_1047 = aux_67
	goto b9
b8:
	;
	v105 = case__11247_71 == vm.Keyword("load-var")
	if v105 {
		l_92 = l_68
		nid_93 = nid_69
		f_94 = f_70
		case__11247_95 = case__11247_71
		op_96 = op_72
		aux_97 = aux_73
		goto b10
	} else {
		l_98 = l_68
		nid_99 = nid_69
		f_100 = f_70
		case__11247_101 = case__11247_71
		op_102 = op_72
		aux_103 = aux_73
		goto b11
	}
b9:
	;
	v1049 = v1041
	l_1050 = l_1042
	nid_1051 = nid_1043
	f_1052 = f_1044
	case__11247_1053 = case__11247_1045
	op_1054 = op_1046
	aux_1055 = aux_1047
	goto b6
b10:
	;
	arg__11329_108, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "chunk-of").Deref(), []vm.Value{l_92})
	if callErr != nil {
		return nil, callErr
	}
	arg__11335_111, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "chunk-of").Deref(), []vm.Value{l_92})
	if callErr != nil {
		return nil, callErr
	}
	idx_112, callErr = rt.InvokeValue(rt.LookupVar("ir", "chunk-intern-const").Deref(), []vm.Value{arg__11335_111, aux_97})
	if callErr != nil {
		return nil, callErr
	}
	v116, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "emit-with-arg!").Deref(), []vm.Value{l_92, nid_93, vm.Keyword("load-var"), idx_112})
	if callErr != nil {
		return nil, callErr
	}
	v120, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "bump-stack-sp!").Deref(), []vm.Value{l_92, vm.Int(1)})
	if callErr != nil {
		return nil, callErr
	}
	v122, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "bump-max-stack!").Deref(), []vm.Value{l_92})
	if callErr != nil {
		return nil, callErr
	}
	v1033 = v122
	l_1034 = l_92
	nid_1035 = nid_93
	f_1036 = f_94
	case__11247_1037 = case__11247_95
	op_1038 = op_96
	aux_1039 = aux_97
	goto b12
b11:
	;
	or__x_137 = case__11247_101 == vm.Keyword("add")
	if or__x_137 {
		l_138 = l_98
		nid_139 = nid_99
		f_140 = f_100
		case__11247_141 = case__11247_101
		op_142 = op_102
		aux_143 = aux_103
		or__x_144 = or__x_137
		goto b16
	} else {
		l_145 = l_98
		nid_146 = nid_99
		f_147 = f_100
		case__11247_148 = case__11247_101
		op_149 = op_102
		aux_150 = aux_103
		or__x_151 = or__x_137
		goto b17
	}
b12:
	;
	v1041 = v1033
	l_1042 = l_1034
	nid_1043 = nid_1035
	f_1044 = f_1036
	case__11247_1045 = case__11247_1037
	op_1046 = op_1038
	aux_1047 = aux_1039
	goto b9
b13:
	;
	v329, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "emit!").Deref(), []vm.Value{l_124, nid_125, op_128})
	if callErr != nil {
		return nil, callErr
	}
	v333, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "bump-stack-sp!").Deref(), []vm.Value{l_124, vm.Int(-1)})
	if callErr != nil {
		return nil, callErr
	}
	v1025 = v333
	l_1026 = l_124
	nid_1027 = nid_125
	f_1028 = f_126
	case__11247_1029 = case__11247_127
	op_1030 = op_128
	aux_1031 = aux_129
	goto b15
b14:
	;
	or__x_348 = case__11247_133 == vm.Keyword("inc")
	if or__x_348 {
		l_349 = l_130
		nid_350 = nid_131
		f_351 = f_132
		case__11247_352 = case__11247_133
		op_353 = op_134
		aux_354 = aux_135
		or__x_355 = or__x_348
		goto b40
	} else {
		l_356 = l_130
		nid_357 = nid_131
		f_358 = f_132
		case__11247_359 = case__11247_133
		op_360 = op_134
		aux_361 = aux_135
		or__x_362 = or__x_348
		goto b41
	}
b15:
	;
	v1033 = v1025
	l_1034 = l_1026
	nid_1035 = nid_1027
	f_1036 = f_1028
	case__11247_1037 = case__11247_1029
	op_1038 = op_1030
	aux_1039 = aux_1031
	goto b12
b16:
	;
	v319 = or__x_144
	l_320 = l_138
	nid_321 = nid_139
	f_322 = f_140
	case__11247_323 = case__11247_141
	op_324 = op_142
	aux_325 = aux_143
	or__x_326 = vm.Boolean(or__x_144)
	goto b18
b17:
	;
	or__x_155 = case__11247_148 == vm.Keyword("sub")
	if or__x_155 {
		l_156 = l_145
		nid_157 = nid_146
		f_158 = f_147
		case__11247_159 = case__11247_148
		op_160 = op_149
		aux_161 = aux_150
		or__x_162 = or__x_155
		goto b19
	} else {
		l_163 = l_145
		nid_164 = nid_146
		f_165 = f_147
		case__11247_166 = case__11247_148
		op_167 = op_149
		aux_168 = aux_150
		or__x_169 = or__x_155
		goto b20
	}
b18:
	;
	if v319 {
		l_124 = l_320
		nid_125 = nid_321
		f_126 = f_322
		case__11247_127 = case__11247_323
		op_128 = op_324
		aux_129 = aux_325
		goto b13
	} else {
		l_130 = l_320
		nid_131 = nid_321
		f_132 = f_322
		case__11247_133 = case__11247_323
		op_134 = op_324
		aux_135 = aux_325
		goto b14
	}
b19:
	;
	v310 = or__x_162
	l_311 = l_156
	nid_312 = nid_157
	f_313 = f_158
	case__11247_314 = case__11247_159
	op_315 = op_160
	aux_316 = aux_161
	or__x_317 = vm.Boolean(or__x_162)
	goto b21
b20:
	;
	or__x_173 = case__11247_166 == vm.Keyword("mul")
	if or__x_173 {
		l_174 = l_163
		nid_175 = nid_164
		f_176 = f_165
		case__11247_177 = case__11247_166
		op_178 = op_167
		aux_179 = aux_168
		or__x_180 = or__x_173
		goto b22
	} else {
		l_181 = l_163
		nid_182 = nid_164
		f_183 = f_165
		case__11247_184 = case__11247_166
		op_185 = op_167
		aux_186 = aux_168
		or__x_187 = or__x_173
		goto b23
	}
b21:
	;
	v319 = v310
	l_320 = l_311
	nid_321 = nid_312
	f_322 = f_313
	case__11247_323 = case__11247_314
	op_324 = op_315
	aux_325 = aux_316
	or__x_326 = vm.Boolean(or__x_151)
	goto b18
b22:
	;
	v301 = or__x_180
	l_302 = l_174
	nid_303 = nid_175
	f_304 = f_176
	case__11247_305 = case__11247_177
	op_306 = op_178
	aux_307 = aux_179
	or__x_308 = vm.Boolean(or__x_180)
	goto b24
b23:
	;
	or__x_191 = case__11247_184 == vm.Keyword("lt")
	if or__x_191 {
		l_192 = l_181
		nid_193 = nid_182
		f_194 = f_183
		case__11247_195 = case__11247_184
		op_196 = op_185
		aux_197 = aux_186
		or__x_198 = or__x_191
		goto b25
	} else {
		l_199 = l_181
		nid_200 = nid_182
		f_201 = f_183
		case__11247_202 = case__11247_184
		op_203 = op_185
		aux_204 = aux_186
		or__x_205 = or__x_191
		goto b26
	}
b24:
	;
	v310 = v301
	l_311 = l_302
	nid_312 = nid_303
	f_313 = f_304
	case__11247_314 = case__11247_305
	op_315 = op_306
	aux_316 = aux_307
	or__x_317 = vm.Boolean(or__x_169)
	goto b21
b25:
	;
	v292 = or__x_198
	l_293 = l_192
	nid_294 = nid_193
	f_295 = f_194
	case__11247_296 = case__11247_195
	op_297 = op_196
	aux_298 = aux_197
	or__x_299 = vm.Boolean(or__x_198)
	goto b27
b26:
	;
	or__x_209 = case__11247_202 == vm.Keyword("lte")
	if or__x_209 {
		l_210 = l_199
		nid_211 = nid_200
		f_212 = f_201
		case__11247_213 = case__11247_202
		op_214 = op_203
		aux_215 = aux_204
		or__x_216 = or__x_209
		goto b28
	} else {
		l_217 = l_199
		nid_218 = nid_200
		f_219 = f_201
		case__11247_220 = case__11247_202
		op_221 = op_203
		aux_222 = aux_204
		or__x_223 = or__x_209
		goto b29
	}
b27:
	;
	v301 = v292
	l_302 = l_293
	nid_303 = nid_294
	f_304 = f_295
	case__11247_305 = case__11247_296
	op_306 = op_297
	aux_307 = aux_298
	or__x_308 = vm.Boolean(or__x_187)
	goto b24
b28:
	;
	v283 = or__x_216
	l_284 = l_210
	nid_285 = nid_211
	f_286 = f_212
	case__11247_287 = case__11247_213
	op_288 = op_214
	aux_289 = aux_215
	or__x_290 = vm.Boolean(or__x_216)
	goto b30
b29:
	;
	or__x_227 = case__11247_220 == vm.Keyword("gt")
	if or__x_227 {
		l_228 = l_217
		nid_229 = nid_218
		f_230 = f_219
		case__11247_231 = case__11247_220
		op_232 = op_221
		aux_233 = aux_222
		or__x_234 = or__x_227
		goto b31
	} else {
		l_235 = l_217
		nid_236 = nid_218
		f_237 = f_219
		case__11247_238 = case__11247_220
		op_239 = op_221
		aux_240 = aux_222
		or__x_241 = or__x_227
		goto b32
	}
b30:
	;
	v292 = v283
	l_293 = l_284
	nid_294 = nid_285
	f_295 = f_286
	case__11247_296 = case__11247_287
	op_297 = op_288
	aux_298 = aux_289
	or__x_299 = vm.Boolean(or__x_205)
	goto b27
b31:
	;
	v274 = or__x_234
	l_275 = l_228
	nid_276 = nid_229
	f_277 = f_230
	case__11247_278 = case__11247_231
	op_279 = op_232
	aux_280 = aux_233
	or__x_281 = vm.Boolean(or__x_234)
	goto b33
b32:
	;
	or__x_245 = case__11247_238 == vm.Keyword("gte")
	if or__x_245 {
		l_246 = l_235
		nid_247 = nid_236
		f_248 = f_237
		case__11247_249 = case__11247_238
		op_250 = op_239
		aux_251 = aux_240
		or__x_252 = or__x_245
		goto b34
	} else {
		l_253 = l_235
		nid_254 = nid_236
		f_255 = f_237
		case__11247_256 = case__11247_238
		op_257 = op_239
		aux_258 = aux_240
		or__x_259 = or__x_245
		goto b35
	}
b33:
	;
	v283 = v274
	l_284 = l_275
	nid_285 = nid_276
	f_286 = f_277
	case__11247_287 = case__11247_278
	op_288 = op_279
	aux_289 = aux_280
	or__x_290 = vm.Boolean(or__x_223)
	goto b30
b34:
	;
	v265 = or__x_252
	l_266 = l_246
	nid_267 = nid_247
	f_268 = f_248
	case__11247_269 = case__11247_249
	op_270 = op_250
	aux_271 = aux_251
	or__x_272 = vm.Boolean(or__x_252)
	goto b36
b35:
	;
	v263 = case__11247_256 == vm.Keyword("eq")
	v265 = v263
	l_266 = l_253
	nid_267 = nid_254
	f_268 = f_255
	case__11247_269 = case__11247_256
	op_270 = op_257
	aux_271 = aux_258
	or__x_272 = vm.Boolean(or__x_259)
	goto b36
b36:
	;
	v274 = v265
	l_275 = l_266
	nid_276 = nid_267
	f_277 = f_268
	case__11247_278 = case__11247_269
	op_279 = op_270
	aux_280 = aux_271
	or__x_281 = vm.Boolean(or__x_241)
	goto b33
b37:
	;
	v378, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "emit!").Deref(), []vm.Value{l_335, nid_336, op_339})
	if callErr != nil {
		return nil, callErr
	}
	v1017 = v378
	l_1018 = l_335
	nid_1019 = nid_336
	f_1020 = f_337
	case__11247_1021 = case__11247_338
	op_1022 = op_339
	aux_1023 = aux_340
	goto b39
b38:
	;
	v393 = case__11247_344 == vm.Keyword("call")
	if v393 {
		l_380 = l_341
		nid_381 = nid_342
		f_382 = f_343
		case__11247_383 = case__11247_344
		op_384 = op_345
		aux_385 = aux_346
		goto b43
	} else {
		l_386 = l_341
		nid_387 = nid_342
		f_388 = f_343
		case__11247_389 = case__11247_344
		op_390 = op_345
		aux_391 = aux_346
		goto b44
	}
b39:
	;
	v1025 = v1017
	l_1026 = l_1018
	nid_1027 = nid_1019
	f_1028 = f_1020
	case__11247_1029 = case__11247_1021
	op_1030 = op_1022
	aux_1031 = aux_1023
	goto b15
b40:
	;
	v368 = or__x_355
	l_369 = l_349
	nid_370 = nid_350
	f_371 = f_351
	case__11247_372 = case__11247_352
	op_373 = op_353
	aux_374 = aux_354
	or__x_375 = vm.Boolean(or__x_355)
	goto b42
b41:
	;
	v366 = case__11247_359 == vm.Keyword("dec")
	v368 = v366
	l_369 = l_356
	nid_370 = nid_357
	f_371 = f_358
	case__11247_372 = case__11247_359
	op_373 = op_360
	aux_374 = aux_361
	or__x_375 = vm.Boolean(or__x_362)
	goto b42
b42:
	;
	if v368 {
		l_335 = l_369
		nid_336 = nid_370
		f_337 = f_371
		case__11247_338 = case__11247_372
		op_339 = op_373
		aux_340 = aux_374
		goto b37
	} else {
		l_341 = l_369
		nid_342 = nid_370
		f_343 = f_371
		case__11247_344 = case__11247_372
		op_345 = op_373
		aux_346 = aux_374
		goto b38
	}
b43:
	;
	arg__11401_397, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "int").Deref(), []vm.Value{aux_385})
	if callErr != nil {
		return nil, callErr
	}
	arg__11409_401, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "int").Deref(), []vm.Value{aux_385})
	if callErr != nil {
		return nil, callErr
	}
	v402, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "emit-with-arg!").Deref(), []vm.Value{l_380, nid_381, vm.Keyword("call"), arg__11409_401})
	if callErr != nil {
		return nil, callErr
	}
	arg__11414_404, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "int").Deref(), []vm.Value{aux_385})
	if callErr != nil {
		return nil, callErr
	}
	arg__11421_409, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "int").Deref(), []vm.Value{aux_385})
	if callErr != nil {
		return nil, callErr
	}
	arg__11422_411 = rt.SubValue(vm.Int(0), arg__11421_409)
	v412, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "bump-stack-sp!").Deref(), []vm.Value{l_380, arg__11422_411})
	if callErr != nil {
		return nil, callErr
	}
	v1009 = v412
	l_1010 = l_380
	nid_1011 = nid_381
	f_1012 = f_382
	case__11247_1013 = case__11247_383
	op_1014 = op_384
	aux_1015 = aux_385
	goto b45
b44:
	;
	v427 = case__11247_389 == vm.Keyword("set-var")
	if v427 {
		l_414 = l_386
		nid_415 = nid_387
		f_416 = f_388
		case__11247_417 = case__11247_389
		op_418 = op_390
		aux_419 = aux_391
		goto b46
	} else {
		l_420 = l_386
		nid_421 = nid_387
		f_422 = f_388
		case__11247_423 = case__11247_389
		op_424 = op_390
		aux_425 = aux_391
		goto b47
	}
b45:
	;
	v1017 = v1009
	l_1018 = l_1010
	nid_1019 = nid_1011
	f_1020 = f_1012
	case__11247_1021 = case__11247_1013
	op_1022 = op_1014
	aux_1023 = aux_1015
	goto b39
b46:
	;
	v432, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "emit!").Deref(), []vm.Value{l_414, nid_415, vm.Keyword("set-var")})
	if callErr != nil {
		return nil, callErr
	}
	v436, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "bump-stack-sp!").Deref(), []vm.Value{l_414, vm.Int(-1)})
	if callErr != nil {
		return nil, callErr
	}
	v1001 = v436
	l_1002 = l_414
	nid_1003 = nid_415
	f_1004 = f_416
	case__11247_1005 = case__11247_417
	op_1006 = op_418
	aux_1007 = aux_419
	goto b48
b47:
	;
	v451 = case__11247_423 == vm.Keyword("tail-call")
	if v451 {
		l_438 = l_420
		nid_439 = nid_421
		f_440 = f_422
		case__11247_441 = case__11247_423
		op_442 = op_424
		aux_443 = aux_425
		goto b49
	} else {
		l_444 = l_420
		nid_445 = nid_421
		f_446 = f_422
		case__11247_447 = case__11247_423
		op_448 = op_424
		aux_449 = aux_425
		goto b50
	}
b48:
	;
	v1009 = v1001
	l_1010 = l_1002
	nid_1011 = nid_1003
	f_1012 = f_1004
	case__11247_1013 = case__11247_1005
	op_1014 = op_1006
	aux_1015 = aux_1007
	goto b45
b49:
	;
	arg__11445_455, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "int").Deref(), []vm.Value{aux_443})
	if callErr != nil {
		return nil, callErr
	}
	arg__11453_459, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "int").Deref(), []vm.Value{aux_443})
	if callErr != nil {
		return nil, callErr
	}
	v460, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "emit-with-arg!").Deref(), []vm.Value{l_438, nid_439, vm.Keyword("tail-call"), arg__11453_459})
	if callErr != nil {
		return nil, callErr
	}
	v993 = v460
	l_994 = l_438
	nid_995 = nid_439
	f_996 = f_440
	case__11247_997 = case__11247_441
	op_998 = op_442
	aux_999 = aux_443
	goto b51
b50:
	;
	v475 = case__11247_447 == vm.Keyword("return")
	if v475 {
		l_462 = l_444
		nid_463 = nid_445
		f_464 = f_446
		case__11247_465 = case__11247_447
		op_466 = op_448
		aux_467 = aux_449
		goto b52
	} else {
		l_468 = l_444
		nid_469 = nid_445
		f_470 = f_446
		case__11247_471 = case__11247_447
		op_472 = op_448
		aux_473 = aux_449
		goto b53
	}
b51:
	;
	v1001 = v993
	l_1002 = l_994
	nid_1003 = nid_995
	f_1004 = f_996
	case__11247_1005 = case__11247_997
	op_1006 = op_998
	aux_1007 = aux_999
	goto b48
b52:
	;
	v480, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "emit!").Deref(), []vm.Value{l_462, nid_463, vm.Keyword("return")})
	if callErr != nil {
		return nil, callErr
	}
	v985 = v480
	l_986 = l_462
	nid_987 = nid_463
	f_988 = f_464
	case__11247_989 = case__11247_465
	op_990 = op_466
	aux_991 = aux_467
	goto b54
b53:
	;
	v495 = case__11247_471 == vm.Keyword("branch")
	if v495 {
		l_482 = l_468
		nid_483 = nid_469
		f_484 = f_470
		case__11247_485 = case__11247_471
		op_486 = op_472
		aux_487 = aux_473
		goto b55
	} else {
		l_488 = l_468
		nid_489 = nid_469
		f_490 = f_470
		case__11247_491 = case__11247_471
		op_492 = op_472
		aux_493 = aux_473
		goto b56
	}
b54:
	;
	v993 = v985
	l_994 = l_986
	nid_995 = nid_987
	f_996 = f_988
	case__11247_997 = case__11247_989
	op_998 = op_990
	aux_999 = aux_991
	goto b51
b55:
	;
	args_498, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-args").Deref(), []vm.Value{aux_487})
	if callErr != nil {
		return nil, callErr
	}
	target_500, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-target").Deref(), []vm.Value{aux_487})
	if callErr != nil {
		return nil, callErr
	}
	argc_502, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{args_498})
	if callErr != nil {
		return nil, callErr
	}
	arg__11479_504, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-params").Deref(), []vm.Value{target_500, f_484})
	if callErr != nil {
		return nil, callErr
	}
	arg__11486_507, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-params").Deref(), []vm.Value{target_500, f_484})
	if callErr != nil {
		return nil, callErr
	}
	target_params_508, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{arg__11486_507})
	if callErr != nil {
		return nil, callErr
	}
	cur_sp_510, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "stack-sp").Deref(), []vm.Value{l_482})
	if callErr != nil {
		return nil, callErr
	}
	drop_count_511 = rt.SubValue(cur_sp_510, target_params_508)
	arg__11497_514, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{l_482})
	if callErr != nil {
		return nil, callErr
	}
	arg__11498_515, callErr = rt.InvokeValue(vm.Keyword("current-block"), []vm.Value{arg__11497_514})
	if callErr != nil {
		return nil, callErr
	}
	arg__11505_519, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{l_482})
	if callErr != nil {
		return nil, callErr
	}
	arg__11506_520, callErr = rt.InvokeValue(vm.Keyword("current-block"), []vm.Value{arg__11505_519})
	if callErr != nil {
		return nil, callErr
	}
	cur_junk_521, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "block-junk-of").Deref(), []vm.Value{l_482, arg__11506_520})
	if callErr != nil {
		return nil, callErr
	}
	v551, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "pos?").Deref(), []vm.Value{drop_count_511})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v551) {
		l_522 = l_482
		nid_523 = nid_483
		f_524 = f_484
		case__11247_525 = case__11247_485
		op_526 = op_486
		bt_527 = aux_487
		aux_528 = aux_487
		args_529 = args_498
		target_530 = target_500
		argc_531 = argc_502
		target_params_532 = target_params_508
		cur_sp_533 = cur_sp_510
		drop_count_534 = drop_count_511
		cur_junk_535 = cur_junk_521
		goto b58
	} else {
		l_536 = l_482
		nid_537 = nid_483
		f_538 = f_484
		case__11247_539 = case__11247_485
		op_540 = op_486
		bt_541 = aux_487
		aux_542 = aux_487
		args_543 = args_498
		target_544 = target_500
		argc_545 = argc_502
		target_params_546 = target_params_508
		cur_sp_547 = cur_sp_510
		drop_count_548 = drop_count_511
		cur_junk_549 = cur_junk_521
		goto b59
	}
b56:
	;
	v689 = case__11247_491 == vm.Keyword("branch-if")
	if v689 {
		l_676 = l_488
		nid_677 = nid_489
		f_678 = f_490
		case__11247_679 = case__11247_491
		op_680 = op_492
		aux_681 = aux_493
		goto b64
	} else {
		l_682 = l_488
		nid_683 = nid_489
		f_684 = f_490
		case__11247_685 = case__11247_491
		op_686 = op_492
		aux_687 = aux_493
		goto b65
	}
b57:
	;
	v985 = v977
	l_986 = l_978
	nid_987 = nid_979
	f_988 = f_980
	case__11247_989 = case__11247_981
	op_990 = op_982
	aux_991 = aux_983
	goto b54
b58:
	;
	ignore_553 = rt.SubValue(drop_count_534, argc_531)
	arg__11515_555, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "chunk-of").Deref(), []vm.Value{l_522})
	if callErr != nil {
		return nil, callErr
	}
	arg__11523_558, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "chunk-of").Deref(), []vm.Value{l_522})
	if callErr != nil {
		return nil, callErr
	}
	off_ip_559, callErr = rt.InvokeValue(rt.LookupVar("ir", "chunk-emit-recur").Deref(), []vm.Value{arg__11523_558, cur_sp_533, argc_531, ignore_553})
	if callErr != nil {
		return nil, callErr
	}
	v561, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "record-source-info!").Deref(), []vm.Value{l_522, nid_523})
	if callErr != nil {
		return nil, callErr
	}
	arg__11536_564 = rt.SubValue(off_ip_559, vm.Int(1))
	arg__11543_570, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("src-ip"), arg__11536_564, vm.Keyword("offset-slot"), vm.Int(1), vm.Keyword("negate?"), vm.Boolean(true), vm.Keyword("target-block"), target_530})
	if callErr != nil {
		return nil, callErr
	}
	arg__11556_580, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("src-ip"), arg__11536_564, vm.Keyword("offset-slot"), vm.Int(1), vm.Keyword("negate?"), vm.Boolean(true), vm.Keyword("target-block"), target_530})
	if callErr != nil {
		return nil, callErr
	}
	v581, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "add-patch!").Deref(), []vm.Value{l_522, arg__11556_580})
	if callErr != nil {
		return nil, callErr
	}
	v583, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "set-stack-sp!").Deref(), []vm.Value{l_522, target_params_532})
	if callErr != nil {
		return nil, callErr
	}
	v660 = v583
	l_661 = l_522
	nid_662 = nid_523
	f_663 = f_524
	case__11247_664 = case__11247_525
	op_665 = op_526
	bt_666 = bt_527
	aux_667 = aux_528
	args_668 = args_529
	target_669 = target_530
	argc_670 = argc_531
	target_params_671 = target_params_532
	cur_sp_672 = cur_sp_533
	drop_count_673 = drop_count_534
	cur_junk_674 = cur_junk_535
	goto b60
b59:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		l_585 = l_536
		nid_586 = nid_537
		f_587 = f_538
		case__11247_588 = case__11247_539
		op_589 = op_540
		bt_590 = bt_541
		aux_591 = aux_542
		args_592 = args_543
		target_593 = target_544
		argc_594 = argc_545
		target_params_595 = target_params_546
		cur_sp_596 = cur_sp_547
		drop_count_597 = drop_count_548
		cur_junk_598 = cur_junk_549
		goto b61
	} else {
		l_599 = l_536
		nid_600 = nid_537
		f_601 = f_538
		case__11247_602 = case__11247_539
		op_603 = op_540
		bt_604 = bt_541
		aux_605 = aux_542
		args_606 = args_543
		target_607 = target_544
		argc_608 = argc_545
		target_params_609 = target_params_546
		cur_sp_610 = cur_sp_547
		drop_count_611 = drop_count_548
		cur_junk_612 = cur_junk_549
		goto b62
	}
b60:
	;
	v977 = v660
	l_978 = l_661
	nid_979 = nid_662
	f_980 = f_663
	case__11247_981 = case__11247_664
	op_982 = op_665
	aux_983 = aux_667
	goto b57
b61:
	;
	arg_ip_618, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "emit-placeholder!").Deref(), []vm.Value{l_585, nid_586, vm.Keyword("branch")})
	if callErr != nil {
		return nil, callErr
	}
	v620, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "record-block-junk!").Deref(), []vm.Value{l_585, target_593, cur_junk_598})
	if callErr != nil {
		return nil, callErr
	}
	arg__11580_623 = rt.SubValue(arg_ip_618, vm.Int(1))
	arg__11587_629, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("src-ip"), arg__11580_623, vm.Keyword("offset-slot"), vm.Int(1), vm.Keyword("negate?"), vm.Boolean(false), vm.Keyword("target-block"), target_593})
	if callErr != nil {
		return nil, callErr
	}
	arg__11600_639, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("src-ip"), arg__11580_623, vm.Keyword("offset-slot"), vm.Int(1), vm.Keyword("negate?"), vm.Boolean(false), vm.Keyword("target-block"), target_593})
	if callErr != nil {
		return nil, callErr
	}
	v640, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "add-patch!").Deref(), []vm.Value{l_585, arg__11600_639})
	if callErr != nil {
		return nil, callErr
	}
	v644 = v640
	l_645 = l_585
	nid_646 = nid_586
	f_647 = f_587
	case__11247_648 = case__11247_588
	op_649 = op_589
	bt_650 = bt_590
	aux_651 = aux_591
	args_652 = args_592
	target_653 = target_593
	argc_654 = argc_594
	target_params_655 = target_params_595
	cur_sp_656 = cur_sp_596
	drop_count_657 = drop_count_597
	cur_junk_658 = cur_junk_598
	goto b63
b62:
	;
	v644 = vm.NIL
	l_645 = l_599
	nid_646 = nid_600
	f_647 = f_601
	case__11247_648 = case__11247_602
	op_649 = op_603
	bt_650 = bt_604
	aux_651 = aux_605
	args_652 = args_606
	target_653 = target_607
	argc_654 = argc_608
	target_params_655 = target_params_609
	cur_sp_656 = cur_sp_610
	drop_count_657 = drop_count_611
	cur_junk_658 = cur_junk_612
	goto b63
b63:
	;
	v660 = v644
	l_661 = l_645
	nid_662 = nid_646
	f_663 = f_647
	case__11247_664 = case__11247_648
	op_665 = op_649
	bt_666 = bt_650
	aux_667 = aux_651
	args_668 = args_652
	target_669 = target_653
	argc_670 = argc_654
	target_params_671 = target_params_655
	cur_sp_672 = cur_sp_656
	drop_count_673 = drop_count_657
	cur_junk_674 = cur_junk_658
	goto b60
b64:
	;
	ft_692, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-false").Deref(), []vm.Value{aux_681})
	if callErr != nil {
		return nil, callErr
	}
	tt_694, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-true").Deref(), []vm.Value{aux_681})
	if callErr != nil {
		return nil, callErr
	}
	ft_target_696, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-target").Deref(), []vm.Value{ft_692})
	if callErr != nil {
		return nil, callErr
	}
	tt_target_698, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-target").Deref(), []vm.Value{tt_694})
	if callErr != nil {
		return nil, callErr
	}
	arg_ip_702, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "emit-placeholder!").Deref(), []vm.Value{l_676, nid_677, vm.Keyword("branch-if")})
	if callErr != nil {
		return nil, callErr
	}
	arg__11627_705, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{l_676})
	if callErr != nil {
		return nil, callErr
	}
	arg__11628_706, callErr = rt.InvokeValue(vm.Keyword("current-block"), []vm.Value{arg__11627_705})
	if callErr != nil {
		return nil, callErr
	}
	arg__11635_710, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{l_676})
	if callErr != nil {
		return nil, callErr
	}
	arg__11636_711, callErr = rt.InvokeValue(vm.Keyword("current-block"), []vm.Value{arg__11635_710})
	if callErr != nil {
		return nil, callErr
	}
	cur_junk_712, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "block-junk-of").Deref(), []vm.Value{l_676, arg__11636_711})
	if callErr != nil {
		return nil, callErr
	}
	arg__11640_714, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "stack-sp").Deref(), []vm.Value{l_676})
	if callErr != nil {
		return nil, callErr
	}
	arg__11642_715 = rt.AddValue(arg__11640_714, cur_junk_712)
	arg__11649_718, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-params").Deref(), []vm.Value{ft_target_696, f_678})
	if callErr != nil {
		return nil, callErr
	}
	arg__11656_721, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-params").Deref(), []vm.Value{ft_target_696, f_678})
	if callErr != nil {
		return nil, callErr
	}
	arg__11657_722, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{arg__11656_721})
	if callErr != nil {
		return nil, callErr
	}
	v723 = rt.SubValue(arg__11642_715, vm.Int(1))
	target_junk_724 = rt.SubValue(v723, arg__11657_722)
	v726, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "record-block-junk!").Deref(), []vm.Value{l_676, ft_target_696, target_junk_724})
	if callErr != nil {
		return nil, callErr
	}
	v728, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "record-block-junk!").Deref(), []vm.Value{l_676, tt_target_698, target_junk_724})
	if callErr != nil {
		return nil, callErr
	}
	arg__11676_731 = rt.SubValue(arg_ip_702, vm.Int(1))
	arg__11683_737, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("src-ip"), arg__11676_731, vm.Keyword("offset-slot"), vm.Int(1), vm.Keyword("negate?"), vm.Boolean(false), vm.Keyword("target-block"), ft_target_696})
	if callErr != nil {
		return nil, callErr
	}
	arg__11696_747, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("src-ip"), arg__11676_731, vm.Keyword("offset-slot"), vm.Int(1), vm.Keyword("negate?"), vm.Boolean(false), vm.Keyword("target-block"), ft_target_696})
	if callErr != nil {
		return nil, callErr
	}
	v748, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "add-patch!").Deref(), []vm.Value{l_676, arg__11696_747})
	if callErr != nil {
		return nil, callErr
	}
	v752, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "bump-stack-sp!").Deref(), []vm.Value{l_676, vm.Int(-1)})
	if callErr != nil {
		return nil, callErr
	}
	my_block_754, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-of").Deref(), []vm.Value{nid_677, f_678})
	if callErr != nil {
		return nil, callErr
	}
	next_block_id_755 = rt.AddValue(my_block_754, vm.Int(1))
	v789, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not=").Deref(), []vm.Value{tt_target_698, next_block_id_755})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v789) {
		l_756 = l_676
		nid_757 = nid_677
		f_758 = f_678
		case__11247_759 = case__11247_679
		op_760 = op_680
		ct_761 = aux_681
		aux_762 = aux_681
		ft_763 = ft_692
		tt_764 = tt_694
		ft_target_765 = ft_target_696
		tt_target_766 = tt_target_698
		arg_ip_767 = arg_ip_702
		cur_junk_768 = cur_junk_712
		target_junk_769 = target_junk_724
		my_block_770 = my_block_754
		next_block_id_771 = next_block_id_755
		goto b67
	} else {
		l_772 = l_676
		nid_773 = nid_677
		f_774 = f_678
		case__11247_775 = case__11247_679
		op_776 = op_680
		ct_777 = aux_681
		aux_778 = aux_681
		ft_779 = ft_692
		tt_780 = tt_694
		ft_target_781 = ft_target_696
		tt_target_782 = tt_target_698
		arg_ip_783 = arg_ip_702
		cur_junk_784 = cur_junk_712
		target_junk_785 = target_junk_724
		my_block_786 = my_block_754
		next_block_id_787 = next_block_id_755
		goto b68
	}
b65:
	;
	v849 = case__11247_685 == vm.Keyword("load-closed")
	if v849 {
		l_836 = l_682
		nid_837 = nid_683
		f_838 = f_684
		case__11247_839 = case__11247_685
		op_840 = op_686
		aux_841 = aux_687
		goto b70
	} else {
		l_842 = l_682
		nid_843 = nid_683
		f_844 = f_684
		case__11247_845 = case__11247_685
		op_846 = op_686
		aux_847 = aux_687
		goto b71
	}
b66:
	;
	v977 = v969
	l_978 = l_970
	nid_979 = nid_971
	f_980 = f_972
	case__11247_981 = case__11247_973
	op_982 = op_974
	aux_983 = aux_975
	goto b57
b67:
	;
	arg_ip2_794, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "emit-placeholder!").Deref(), []vm.Value{l_756, nid_757, vm.Keyword("branch")})
	if callErr != nil {
		return nil, callErr
	}
	arg__11724_797 = rt.SubValue(arg_ip2_794, vm.Int(1))
	arg__11731_803, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("src-ip"), arg__11724_797, vm.Keyword("offset-slot"), vm.Int(1), vm.Keyword("negate?"), vm.Boolean(false), vm.Keyword("target-block"), tt_target_766})
	if callErr != nil {
		return nil, callErr
	}
	arg__11744_813, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("src-ip"), arg__11724_797, vm.Keyword("offset-slot"), vm.Int(1), vm.Keyword("negate?"), vm.Boolean(false), vm.Keyword("target-block"), tt_target_766})
	if callErr != nil {
		return nil, callErr
	}
	v814, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "add-patch!").Deref(), []vm.Value{l_756, arg__11744_813})
	if callErr != nil {
		return nil, callErr
	}
	v818 = v814
	l_819 = l_756
	nid_820 = nid_757
	f_821 = f_758
	case__11247_822 = case__11247_759
	op_823 = op_760
	ct_824 = ct_761
	aux_825 = aux_762
	ft_826 = ft_763
	tt_827 = tt_764
	ft_target_828 = ft_target_765
	tt_target_829 = tt_target_766
	arg_ip_830 = arg_ip_767
	cur_junk_831 = cur_junk_768
	target_junk_832 = target_junk_769
	my_block_833 = my_block_770
	next_block_id_834 = next_block_id_771
	goto b69
b68:
	;
	v818 = vm.NIL
	l_819 = l_772
	nid_820 = nid_773
	f_821 = f_774
	case__11247_822 = case__11247_775
	op_823 = op_776
	ct_824 = ct_777
	aux_825 = aux_778
	ft_826 = ft_779
	tt_827 = tt_780
	ft_target_828 = ft_target_781
	tt_target_829 = tt_target_782
	arg_ip_830 = arg_ip_783
	cur_junk_831 = cur_junk_784
	target_junk_832 = target_junk_785
	my_block_833 = my_block_786
	next_block_id_834 = next_block_id_787
	goto b69
b69:
	;
	v969 = v818
	l_970 = l_819
	nid_971 = nid_820
	f_972 = f_821
	case__11247_973 = case__11247_822
	op_974 = op_823
	aux_975 = aux_825
	goto b66
b70:
	;
	arg__11753_853, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "int").Deref(), []vm.Value{aux_841})
	if callErr != nil {
		return nil, callErr
	}
	arg__11761_857, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "int").Deref(), []vm.Value{aux_841})
	if callErr != nil {
		return nil, callErr
	}
	v858, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "emit-with-arg!").Deref(), []vm.Value{l_836, nid_837, vm.Keyword("load-closed"), arg__11761_857})
	if callErr != nil {
		return nil, callErr
	}
	v862, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "bump-stack-sp!").Deref(), []vm.Value{l_836, vm.Int(1)})
	if callErr != nil {
		return nil, callErr
	}
	v864, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "bump-max-stack!").Deref(), []vm.Value{l_836})
	if callErr != nil {
		return nil, callErr
	}
	v961 = v864
	l_962 = l_836
	nid_963 = nid_837
	f_964 = f_838
	case__11247_965 = case__11247_839
	op_966 = op_840
	aux_967 = aux_841
	goto b72
b71:
	;
	v879 = case__11247_845 == vm.Keyword("make-closure")
	if v879 {
		l_866 = l_842
		nid_867 = nid_843
		f_868 = f_844
		case__11247_869 = case__11247_845
		op_870 = op_846
		aux_871 = aux_847
		goto b73
	} else {
		l_872 = l_842
		nid_873 = nid_843
		f_874 = f_844
		case__11247_875 = case__11247_845
		op_876 = op_846
		aux_877 = aux_847
		goto b74
	}
b72:
	;
	v969 = v961
	l_970 = l_962
	nid_971 = nid_963
	f_972 = f_964
	case__11247_973 = case__11247_965
	op_974 = op_966
	aux_975 = aux_967
	goto b66
b73:
	;
	v884, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "emit!").Deref(), []vm.Value{l_866, nid_867, vm.Keyword("make-closure")})
	if callErr != nil {
		return nil, callErr
	}
	v953 = v884
	l_954 = l_866
	nid_955 = nid_867
	f_956 = f_868
	case__11247_957 = case__11247_869
	op_958 = op_870
	aux_959 = aux_871
	goto b75
b74:
	;
	v899 = case__11247_875 == vm.Keyword("push-closed")
	if v899 {
		l_886 = l_872
		nid_887 = nid_873
		f_888 = f_874
		case__11247_889 = case__11247_875
		op_890 = op_876
		aux_891 = aux_877
		goto b76
	} else {
		l_892 = l_872
		nid_893 = nid_873
		f_894 = f_874
		case__11247_895 = case__11247_875
		op_896 = op_876
		aux_897 = aux_877
		goto b77
	}
b75:
	;
	v961 = v953
	l_962 = l_954
	nid_963 = nid_955
	f_964 = f_956
	case__11247_965 = case__11247_957
	op_966 = op_958
	aux_967 = aux_959
	goto b72
b76:
	;
	v904, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "emit!").Deref(), []vm.Value{l_886, nid_887, vm.Keyword("push-closed")})
	if callErr != nil {
		return nil, callErr
	}
	v908, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "bump-stack-sp!").Deref(), []vm.Value{l_886, vm.Int(-1)})
	if callErr != nil {
		return nil, callErr
	}
	v945 = v908
	l_946 = l_886
	nid_947 = nid_887
	f_948 = f_888
	case__11247_949 = case__11247_889
	op_950 = op_890
	aux_951 = aux_891
	goto b78
b77:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		l_910 = l_892
		nid_911 = nid_893
		f_912 = f_894
		case__11247_913 = case__11247_895
		op_914 = op_896
		aux_915 = aux_897
		goto b79
	} else {
		l_916 = l_892
		nid_917 = nid_893
		f_918 = f_894
		case__11247_919 = case__11247_895
		op_920 = op_896
		aux_921 = aux_897
		goto b80
	}
b78:
	;
	v953 = v945
	l_954 = l_946
	nid_955 = nid_947
	f_956 = f_948
	case__11247_957 = case__11247_949
	op_958 = op_950
	aux_959 = aux_951
	goto b75
b79:
	;
	arg__11798_927, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("ir/lower: unsupported op for lowering: "), op_914})
	if callErr != nil {
		return nil, callErr
	}
	arg__11805_932, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("ir/lower: unsupported op for lowering: "), op_914})
	if callErr != nil {
		return nil, callErr
	}
	v933, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "throw").Deref(), []vm.Value{arg__11805_932})
	if callErr != nil {
		return nil, callErr
	}
	v937 = v933
	l_938 = l_910
	nid_939 = nid_911
	f_940 = f_912
	case__11247_941 = case__11247_913
	op_942 = op_914
	aux_943 = aux_915
	goto b81
b80:
	;
	v937 = vm.NIL
	l_938 = l_916
	nid_939 = nid_917
	f_940 = f_918
	case__11247_941 = case__11247_919
	op_942 = op_920
	aux_943 = aux_921
	goto b81
b81:
	;
	v945 = v937
	l_946 = l_938
	nid_947 = nid_939
	f_948 = f_940
	case__11247_949 = case__11247_941
	op_950 = op_942
	aux_951 = aux_943
	goto b78
}
func materialize_BANG_(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var f_4 vm.Value
	var pos_6 vm.Value
	var l_7 vm.Value
	var nid_8 vm.Value
	var f_9 vm.Value
	var pos_10 vm.Value
	var arg__11819_25 vm.Value
	var arg__11820_26 vm.Value
	var and__x_27 bool
	var l_11 vm.Value
	var nid_12 vm.Value
	var f_13 vm.Value
	var pos_14 vm.Value
	var op_82 vm.Value
	var aux_84 vm.Value
	var v100 bool
	var v303 vm.Value
	var l_304 vm.Value
	var nid_305 vm.Value
	var f_306 vm.Value
	var pos_307 vm.Value
	var l_16 vm.Value
	var nid_17 vm.Value
	var f_18 vm.Value
	var pos_19 vm.Value
	var l_20 vm.Value
	var nid_21 vm.Value
	var f_22 vm.Value
	var pos_23 vm.Value
	var arg__11831_55 vm.Value
	var arg__11832_56 vm.Value
	var nth_arg_57 vm.Value
	var arg__11837_59 vm.Value
	var arg__11841_61 vm.Value
	var arg__11847_64 vm.Value
	var arg__11851_66 vm.Value
	var v67 vm.Value
	var v71 vm.Value
	var v73 vm.Value
	var v75 vm.Value
	var l_76 vm.Value
	var nid_77 vm.Value
	var f_78 vm.Value
	var pos_79 vm.Value
	var l_28 vm.Value
	var nid_29 vm.Value
	var f_30 vm.Value
	var pos_31 vm.Value
	var and__x_32 bool
	var arg__11827_41 vm.Value
	var v42 bool
	var l_33 vm.Value
	var nid_34 vm.Value
	var f_35 vm.Value
	var pos_36 vm.Value
	var and__x_37 bool
	var v45 bool
	var l_46 vm.Value
	var nid_47 vm.Value
	var f_48 vm.Value
	var pos_49 vm.Value
	var and__x_50 vm.Value
	var l_85 vm.Value
	var nid_86 vm.Value
	var f_87 vm.Value
	var pos_88 vm.Value
	var case__11806_89 vm.Value
	var op_90 vm.Value
	var aux_91 vm.Value
	var arg__11876_103 vm.Value
	var arg__11880_105 vm.Value
	var arg__11885_108 vm.Value
	var arg__11889_110 vm.Value
	var idx_111 vm.Value
	var v115 vm.Value
	var v119 vm.Value
	var v121 vm.Value
	var l_92 vm.Value
	var nid_93 vm.Value
	var f_94 vm.Value
	var pos_95 vm.Value
	var case__11806_96 vm.Value
	var op_97 vm.Value
	var aux_98 vm.Value
	var v138 bool
	var v294 vm.Value
	var l_295 vm.Value
	var nid_296 vm.Value
	var f_297 vm.Value
	var pos_298 vm.Value
	var case__11806_299 vm.Value
	var op_300 vm.Value
	var aux_301 vm.Value
	var l_123 vm.Value
	var nid_124 vm.Value
	var f_125 vm.Value
	var pos_126 vm.Value
	var case__11806_127 vm.Value
	var op_128 vm.Value
	var aux_129 vm.Value
	var arg__11915_142 vm.Value
	var arg__11923_146 vm.Value
	var v147 vm.Value
	var v151 vm.Value
	var v153 vm.Value
	var l_130 vm.Value
	var nid_131 vm.Value
	var f_132 vm.Value
	var pos_133 vm.Value
	var case__11806_134 vm.Value
	var op_135 vm.Value
	var aux_136 vm.Value
	var v170 bool
	var v285 vm.Value
	var l_286 vm.Value
	var nid_287 vm.Value
	var f_288 vm.Value
	var pos_289 vm.Value
	var case__11806_290 vm.Value
	var op_291 vm.Value
	var aux_292 vm.Value
	var l_155 vm.Value
	var nid_156 vm.Value
	var f_157 vm.Value
	var pos_158 vm.Value
	var case__11806_159 vm.Value
	var op_160 vm.Value
	var aux_161 vm.Value
	var arg__11937_173 vm.Value
	var arg__11943_176 vm.Value
	var idx_177 vm.Value
	var v181 vm.Value
	var v185 vm.Value
	var v187 vm.Value
	var l_162 vm.Value
	var nid_163 vm.Value
	var f_164 vm.Value
	var pos_165 vm.Value
	var case__11806_166 vm.Value
	var op_167 vm.Value
	var aux_168 vm.Value
	var v204 bool
	var v276 vm.Value
	var l_277 vm.Value
	var nid_278 vm.Value
	var f_279 vm.Value
	var pos_280 vm.Value
	var case__11806_281 vm.Value
	var op_282 vm.Value
	var aux_283 vm.Value
	var l_189 vm.Value
	var nid_190 vm.Value
	var f_191 vm.Value
	var pos_192 vm.Value
	var case__11806_193 vm.Value
	var op_194 vm.Value
	var aux_195 vm.Value
	var arg__11970_208 vm.Value
	var arg__11978_212 vm.Value
	var v213 vm.Value
	var v217 vm.Value
	var v219 vm.Value
	var l_196 vm.Value
	var nid_197 vm.Value
	var f_198 vm.Value
	var pos_199 vm.Value
	var case__11806_200 vm.Value
	var op_201 vm.Value
	var aux_202 vm.Value
	var v267 vm.Value
	var l_268 vm.Value
	var nid_269 vm.Value
	var f_270 vm.Value
	var pos_271 vm.Value
	var case__11806_272 vm.Value
	var op_273 vm.Value
	var aux_274 vm.Value
	var l_221 vm.Value
	var nid_222 vm.Value
	var f_223 vm.Value
	var pos_224 vm.Value
	var case__11806_225 vm.Value
	var op_226 vm.Value
	var aux_227 vm.Value
	var arg__11998_244 vm.Value
	var arg__12011_253 vm.Value
	var v254 vm.Value
	var l_228 vm.Value
	var nid_229 vm.Value
	var f_230 vm.Value
	var pos_231 vm.Value
	var case__11806_232 vm.Value
	var op_233 vm.Value
	var aux_234 vm.Value
	var v258 vm.Value
	var l_259 vm.Value
	var nid_260 vm.Value
	var f_261 vm.Value
	var pos_262 vm.Value
	var case__11806_263 vm.Value
	var op_264 vm.Value
	var aux_265 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = f_4, pos_6, l_7, nid_8, f_9, pos_10, arg__11819_25, arg__11820_26, and__x_27, l_11, nid_12, f_13, pos_14, op_82, aux_84, v100, v303, l_304, nid_305, f_306, pos_307, l_16, nid_17, f_18, pos_19, l_20, nid_21, f_22, pos_23, arg__11831_55, arg__11832_56, nth_arg_57, arg__11837_59, arg__11841_61, arg__11847_64, arg__11851_66, v67, v71, v73, v75, l_76, nid_77, f_78, pos_79, l_28, nid_29, f_30, pos_31, and__x_32, arg__11827_41, v42, l_33, nid_34, f_35, pos_36, and__x_37, v45, l_46, nid_47, f_48, pos_49, and__x_50, l_85, nid_86, f_87, pos_88, case__11806_89, op_90, aux_91, arg__11876_103, arg__11880_105, arg__11885_108, arg__11889_110, idx_111, v115, v119, v121, l_92, nid_93, f_94, pos_95, case__11806_96, op_97, aux_98, v138, v294, l_295, nid_296, f_297, pos_298, case__11806_299, op_300, aux_301, l_123, nid_124, f_125, pos_126, case__11806_127, op_128, aux_129, arg__11915_142, arg__11923_146, v147, v151, v153, l_130, nid_131, f_132, pos_133, case__11806_134, op_135, aux_136, v170, v285, l_286, nid_287, f_288, pos_289, case__11806_290, op_291, aux_292, l_155, nid_156, f_157, pos_158, case__11806_159, op_160, aux_161, arg__11937_173, arg__11943_176, idx_177, v181, v185, v187, l_162, nid_163, f_164, pos_165, case__11806_166, op_167, aux_168, v204, v276, l_277, nid_278, f_279, pos_280, case__11806_281, op_282, aux_283, l_189, nid_190, f_191, pos_192, case__11806_193, op_194, aux_195, arg__11970_208, arg__11978_212, v213, v217, v219, l_196, nid_197, f_198, pos_199, case__11806_200, op_201, aux_202, v267, l_268, nid_269, f_270, pos_271, case__11806_272, op_273, aux_274, l_221, nid_222, f_223, pos_224, case__11806_225, op_226, aux_227, arg__11998_244, arg__12011_253, v254, l_228, nid_229, f_230, pos_231, case__11806_232, op_233, aux_234, v258, l_259, nid_260, f_261, pos_262, case__11806_263, op_264, aux_265
	f_4, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "f-of").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	pos_6, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "value-pos-of").Deref(), []vm.Value{arg0, arg1})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(pos_6) {
		l_7 = arg0
		nid_8 = arg1
		f_9 = f_4
		pos_10 = pos_6
		goto b1
	} else {
		l_11 = arg0
		nid_12 = arg1
		f_13 = f_4
		pos_14 = pos_6
		goto b2
	}
b1:
	;
	arg__11819_25, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "stack-sp").Deref(), []vm.Value{l_7})
	if callErr != nil {
		return nil, callErr
	}
	arg__11820_26 = rt.SubValue(arg__11819_25, vm.Int(1))
	and__x_27 = pos_10 == arg__11820_26
	if and__x_27 {
		l_28 = l_7
		nid_29 = nid_8
		f_30 = f_9
		pos_31 = pos_10
		and__x_32 = and__x_27
		goto b7
	} else {
		l_33 = l_7
		nid_34 = nid_8
		f_35 = f_9
		pos_36 = pos_10
		and__x_37 = and__x_27
		goto b8
	}
b2:
	;
	op_82, callErr = rt.InvokeValue(rt.LookupVar("ir", "op").Deref(), []vm.Value{nid_12, f_13})
	if callErr != nil {
		return nil, callErr
	}
	aux_84, callErr = rt.InvokeValue(rt.LookupVar("ir", "aux").Deref(), []vm.Value{nid_12, f_13})
	if callErr != nil {
		return nil, callErr
	}
	v100 = op_82 == vm.Keyword("const")
	if v100 {
		l_85 = l_11
		nid_86 = nid_12
		f_87 = f_13
		pos_88 = pos_14
		case__11806_89 = op_82
		op_90 = op_82
		aux_91 = aux_84
		goto b10
	} else {
		l_92 = l_11
		nid_93 = nid_12
		f_94 = f_13
		pos_95 = pos_14
		case__11806_96 = op_82
		op_97 = op_82
		aux_98 = aux_84
		goto b11
	}
b3:
	;
	return v303, nil
b4:
	;
	v75 = vm.NIL
	l_76 = l_16
	nid_77 = nid_17
	f_78 = f_18
	pos_79 = pos_19
	goto b6
b5:
	;
	arg__11831_55, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "stack-sp").Deref(), []vm.Value{l_20})
	if callErr != nil {
		return nil, callErr
	}
	arg__11832_56 = rt.SubValue(arg__11831_55, vm.Int(1))
	nth_arg_57 = rt.SubValue(arg__11832_56, pos_23)
	arg__11837_59, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "chunk-of").Deref(), []vm.Value{l_20})
	if callErr != nil {
		return nil, callErr
	}
	arg__11841_61, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "stack-sp").Deref(), []vm.Value{l_20})
	if callErr != nil {
		return nil, callErr
	}
	arg__11847_64, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "chunk-of").Deref(), []vm.Value{l_20})
	if callErr != nil {
		return nil, callErr
	}
	arg__11851_66, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "stack-sp").Deref(), []vm.Value{l_20})
	if callErr != nil {
		return nil, callErr
	}
	v67, callErr = rt.InvokeValue(rt.LookupVar("ir", "chunk-emit-dup-nth").Deref(), []vm.Value{arg__11847_64, arg__11851_66, nth_arg_57})
	if callErr != nil {
		return nil, callErr
	}
	v71, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "bump-stack-sp!").Deref(), []vm.Value{l_20, vm.Int(1)})
	if callErr != nil {
		return nil, callErr
	}
	v73, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "bump-max-stack!").Deref(), []vm.Value{l_20})
	if callErr != nil {
		return nil, callErr
	}
	v75 = v73
	l_76 = l_20
	nid_77 = nid_21
	f_78 = f_22
	pos_79 = pos_23
	goto b6
b6:
	;
	v303 = v75
	l_304 = l_76
	nid_305 = nid_77
	f_306 = f_78
	pos_307 = pos_79
	goto b3
b7:
	;
	arg__11827_41, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "use-count-of").Deref(), []vm.Value{l_28, nid_29})
	if callErr != nil {
		return nil, callErr
	}
	v42 = arg__11827_41 == vm.Int(1)
	v45 = v42
	l_46 = l_28
	nid_47 = nid_29
	f_48 = f_30
	pos_49 = pos_31
	and__x_50 = vm.Boolean(and__x_32)
	goto b9
b8:
	;
	v45 = and__x_37
	l_46 = l_33
	nid_47 = nid_34
	f_48 = f_35
	pos_49 = pos_36
	and__x_50 = vm.Boolean(and__x_37)
	goto b9
b9:
	;
	if v45 {
		l_16 = l_46
		nid_17 = nid_47
		f_18 = f_48
		pos_19 = pos_49
		goto b4
	} else {
		l_20 = l_46
		nid_21 = nid_47
		f_22 = f_48
		pos_23 = pos_49
		goto b5
	}
b10:
	;
	arg__11876_103, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "chunk-of").Deref(), []vm.Value{l_85})
	if callErr != nil {
		return nil, callErr
	}
	arg__11880_105, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "template-value").Deref(), []vm.Value{aux_91})
	if callErr != nil {
		return nil, callErr
	}
	arg__11885_108, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "chunk-of").Deref(), []vm.Value{l_85})
	if callErr != nil {
		return nil, callErr
	}
	arg__11889_110, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "template-value").Deref(), []vm.Value{aux_91})
	if callErr != nil {
		return nil, callErr
	}
	idx_111, callErr = rt.InvokeValue(rt.LookupVar("ir", "chunk-intern-const").Deref(), []vm.Value{arg__11885_108, arg__11889_110})
	if callErr != nil {
		return nil, callErr
	}
	v115, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "emit-with-arg!").Deref(), []vm.Value{l_85, nid_86, vm.Keyword("const"), idx_111})
	if callErr != nil {
		return nil, callErr
	}
	v119, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "bump-stack-sp!").Deref(), []vm.Value{l_85, vm.Int(1)})
	if callErr != nil {
		return nil, callErr
	}
	v121, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "bump-max-stack!").Deref(), []vm.Value{l_85})
	if callErr != nil {
		return nil, callErr
	}
	v294 = v121
	l_295 = l_85
	nid_296 = nid_86
	f_297 = f_87
	pos_298 = pos_88
	case__11806_299 = case__11806_89
	op_300 = op_90
	aux_301 = aux_91
	goto b12
b11:
	;
	v138 = case__11806_96 == vm.Keyword("load-arg")
	if v138 {
		l_123 = l_92
		nid_124 = nid_93
		f_125 = f_94
		pos_126 = pos_95
		case__11806_127 = case__11806_96
		op_128 = op_97
		aux_129 = aux_98
		goto b13
	} else {
		l_130 = l_92
		nid_131 = nid_93
		f_132 = f_94
		pos_133 = pos_95
		case__11806_134 = case__11806_96
		op_135 = op_97
		aux_136 = aux_98
		goto b14
	}
b12:
	;
	v303 = v294
	l_304 = l_295
	nid_305 = nid_296
	f_306 = f_297
	pos_307 = pos_298
	goto b3
b13:
	;
	arg__11915_142, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "int").Deref(), []vm.Value{aux_129})
	if callErr != nil {
		return nil, callErr
	}
	arg__11923_146, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "int").Deref(), []vm.Value{aux_129})
	if callErr != nil {
		return nil, callErr
	}
	v147, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "emit-with-arg!").Deref(), []vm.Value{l_123, nid_124, vm.Keyword("load-arg"), arg__11923_146})
	if callErr != nil {
		return nil, callErr
	}
	v151, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "bump-stack-sp!").Deref(), []vm.Value{l_123, vm.Int(1)})
	if callErr != nil {
		return nil, callErr
	}
	v153, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "bump-max-stack!").Deref(), []vm.Value{l_123})
	if callErr != nil {
		return nil, callErr
	}
	v285 = v153
	l_286 = l_123
	nid_287 = nid_124
	f_288 = f_125
	pos_289 = pos_126
	case__11806_290 = case__11806_127
	op_291 = op_128
	aux_292 = aux_129
	goto b15
b14:
	;
	v170 = case__11806_134 == vm.Keyword("load-var")
	if v170 {
		l_155 = l_130
		nid_156 = nid_131
		f_157 = f_132
		pos_158 = pos_133
		case__11806_159 = case__11806_134
		op_160 = op_135
		aux_161 = aux_136
		goto b16
	} else {
		l_162 = l_130
		nid_163 = nid_131
		f_164 = f_132
		pos_165 = pos_133
		case__11806_166 = case__11806_134
		op_167 = op_135
		aux_168 = aux_136
		goto b17
	}
b15:
	;
	v294 = v285
	l_295 = l_286
	nid_296 = nid_287
	f_297 = f_288
	pos_298 = pos_289
	case__11806_299 = case__11806_290
	op_300 = op_291
	aux_301 = aux_292
	goto b12
b16:
	;
	arg__11937_173, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "chunk-of").Deref(), []vm.Value{l_155})
	if callErr != nil {
		return nil, callErr
	}
	arg__11943_176, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "chunk-of").Deref(), []vm.Value{l_155})
	if callErr != nil {
		return nil, callErr
	}
	idx_177, callErr = rt.InvokeValue(rt.LookupVar("ir", "chunk-intern-const").Deref(), []vm.Value{arg__11943_176, aux_161})
	if callErr != nil {
		return nil, callErr
	}
	v181, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "emit-with-arg!").Deref(), []vm.Value{l_155, nid_156, vm.Keyword("load-var"), idx_177})
	if callErr != nil {
		return nil, callErr
	}
	v185, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "bump-stack-sp!").Deref(), []vm.Value{l_155, vm.Int(1)})
	if callErr != nil {
		return nil, callErr
	}
	v187, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "bump-max-stack!").Deref(), []vm.Value{l_155})
	if callErr != nil {
		return nil, callErr
	}
	v276 = v187
	l_277 = l_155
	nid_278 = nid_156
	f_279 = f_157
	pos_280 = pos_158
	case__11806_281 = case__11806_159
	op_282 = op_160
	aux_283 = aux_161
	goto b18
b17:
	;
	v204 = case__11806_166 == vm.Keyword("load-closed")
	if v204 {
		l_189 = l_162
		nid_190 = nid_163
		f_191 = f_164
		pos_192 = pos_165
		case__11806_193 = case__11806_166
		op_194 = op_167
		aux_195 = aux_168
		goto b19
	} else {
		l_196 = l_162
		nid_197 = nid_163
		f_198 = f_164
		pos_199 = pos_165
		case__11806_200 = case__11806_166
		op_201 = op_167
		aux_202 = aux_168
		goto b20
	}
b18:
	;
	v285 = v276
	l_286 = l_277
	nid_287 = nid_278
	f_288 = f_279
	pos_289 = pos_280
	case__11806_290 = case__11806_281
	op_291 = op_282
	aux_292 = aux_283
	goto b15
b19:
	;
	arg__11970_208, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "int").Deref(), []vm.Value{aux_195})
	if callErr != nil {
		return nil, callErr
	}
	arg__11978_212, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "int").Deref(), []vm.Value{aux_195})
	if callErr != nil {
		return nil, callErr
	}
	v213, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "emit-with-arg!").Deref(), []vm.Value{l_189, nid_190, vm.Keyword("load-closed"), arg__11978_212})
	if callErr != nil {
		return nil, callErr
	}
	v217, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "bump-stack-sp!").Deref(), []vm.Value{l_189, vm.Int(1)})
	if callErr != nil {
		return nil, callErr
	}
	v219, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "bump-max-stack!").Deref(), []vm.Value{l_189})
	if callErr != nil {
		return nil, callErr
	}
	v267 = v219
	l_268 = l_189
	nid_269 = nid_190
	f_270 = f_191
	pos_271 = pos_192
	case__11806_272 = case__11806_193
	op_273 = op_194
	aux_274 = aux_195
	goto b21
b20:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		l_221 = l_196
		nid_222 = nid_197
		f_223 = f_198
		pos_224 = pos_199
		case__11806_225 = case__11806_200
		op_226 = op_201
		aux_227 = aux_202
		goto b22
	} else {
		l_228 = l_196
		nid_229 = nid_197
		f_230 = f_198
		pos_231 = pos_199
		case__11806_232 = case__11806_200
		op_233 = op_201
		aux_234 = aux_202
		goto b23
	}
b21:
	;
	v276 = v267
	l_277 = l_268
	nid_278 = nid_269
	f_279 = f_270
	pos_280 = pos_271
	case__11806_281 = case__11806_272
	op_282 = op_273
	aux_283 = aux_274
	goto b18
b22:
	;
	arg__11998_244, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("ir/lower: value %"), nid_222, vm.String(" not on stack for materialize (op="), op_226, vm.String(")")})
	if callErr != nil {
		return nil, callErr
	}
	arg__12011_253, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("ir/lower: value %"), nid_222, vm.String(" not on stack for materialize (op="), op_226, vm.String(")")})
	if callErr != nil {
		return nil, callErr
	}
	v254, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "throw").Deref(), []vm.Value{arg__12011_253})
	if callErr != nil {
		return nil, callErr
	}
	v258 = v254
	l_259 = l_221
	nid_260 = nid_222
	f_261 = f_223
	pos_262 = pos_224
	case__11806_263 = case__11806_225
	op_264 = op_226
	aux_265 = aux_227
	goto b24
b23:
	;
	v258 = vm.NIL
	l_259 = l_228
	nid_260 = nid_229
	f_261 = f_230
	pos_262 = pos_231
	case__11806_263 = case__11806_232
	op_264 = op_233
	aux_265 = aux_234
	goto b24
b24:
	;
	v267 = v258
	l_268 = l_259
	nid_269 = nid_260
	f_270 = f_261
	pos_271 = pos_262
	case__11806_272 = case__11806_263
	op_273 = op_264
	aux_274 = aux_265
	goto b21
}
func materialize_branch_args_BANG_(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var argc_4 vm.Value
	var k_5 int
	var argc_6 vm.Value
	var args_7 vm.Value
	var l_8 vm.Value
	var v19 bool
	var k_11 int
	var argc_12 vm.Value
	var args_13 vm.Value
	var l_14 vm.Value
	var k_15 int
	var argc_16 vm.Value
	var args_17 vm.Value
	var l_18 vm.Value
	var arg__12025_23 vm.Value
	var arg__12033_26 vm.Value
	var pos_27 vm.Value
	var skip_72 int
	var k_73 int
	var argc_74 vm.Value
	var args_75 vm.Value
	var l_76 vm.Value
	var v87 bool
	var k_28 int
	var argc_29 vm.Value
	var args_30 vm.Value
	var l_31 vm.Value
	var pos_32 vm.Value
	var v62 int
	var k_33 int
	var argc_34 vm.Value
	var args_35 vm.Value
	var l_36 vm.Value
	var pos_37 vm.Value
	var v65 int
	var k_66 int
	var argc_67 vm.Value
	var args_68 vm.Value
	var l_69 vm.Value
	var pos_70 vm.Value
	var k_38 int
	var argc_39 vm.Value
	var args_40 vm.Value
	var l_41 vm.Value
	var and__x_42 vm.Value
	var pos_43 vm.Value
	var v51 bool
	var k_44 int
	var argc_45 vm.Value
	var args_46 vm.Value
	var l_47 vm.Value
	var and__x_48 vm.Value
	var pos_49 vm.Value
	var v54 vm.Value
	var k_55 int
	var argc_56 vm.Value
	var args_57 vm.Value
	var l_58 vm.Value
	var and__x_59 vm.Value
	var pos_60 vm.Value
	var skip_77 int
	var k_78 int
	var argc_79 vm.Value
	var args_80 vm.Value
	var l_81 vm.Value
	var v90 vm.Value
	var skip_82 int
	var k_83 int
	var argc_84 vm.Value
	var args_85 vm.Value
	var l_86 vm.Value
	var v103 vm.Value
	var v226 vm.Value
	var skip_227 int
	var k_228 int
	var argc_229 vm.Value
	var args_230 vm.Value
	var l_231 vm.Value
	var skip_92 int
	var k_93 int
	var argc_94 vm.Value
	var args_95 vm.Value
	var l_96 vm.Value
	var v106 vm.Value
	var skip_97 int
	var k_98 int
	var argc_99 vm.Value
	var args_100 vm.Value
	var l_101 vm.Value
	var v219 vm.Value
	var skip_220 int
	var k_221 int
	var argc_222 vm.Value
	var args_223 vm.Value
	var l_224 vm.Value
	var skip_108 int
	var k_109 int
	var argc_110 vm.Value
	var args_111 vm.Value
	var l_112 vm.Value
	var arg__12057_121 vm.Value
	var arg__12063_124 vm.Value
	var arg__12065_125 vm.Value
	var arg__12070_128 vm.Value
	var arg__12076_131 vm.Value
	var arg__12078_132 vm.Value
	var doseq_seq__12012_133 vm.Value
	var skip_113 int
	var k_114 int
	var argc_115 vm.Value
	var args_116 vm.Value
	var l_117 vm.Value
	var v212 vm.Value
	var skip_213 int
	var k_214 int
	var argc_215 vm.Value
	var args_216 vm.Value
	var l_217 vm.Value
	var doseq_loop__12013_134 vm.Value
	var l_135 vm.Value
	var skip_137 int
	var k_138 int
	var argc_139 vm.Value
	var args_140 vm.Value
	var doseq_seq__12012_141 vm.Value
	var doseq_loop__12013_142 vm.Value
	var l_143 vm.Value
	var a_153 vm.Value
	var v155 vm.Value
	var v157 vm.Value
	var v159 vm.Value
	var skip_144 int
	var k_145 int
	var argc_146 vm.Value
	var args_147 vm.Value
	var doseq_seq__12012_148 vm.Value
	var doseq_loop__12013_149 vm.Value
	var l_150 vm.Value
	var v163 vm.Value
	var skip_164 int
	var k_165 int
	var argc_166 vm.Value
	var args_167 vm.Value
	var doseq_seq__12012_168 vm.Value
	var doseq_loop__12013_169 vm.Value
	var l_170 vm.Value
	var k_171 int
	var k_172 int
	var args_173 vm.Value
	var l_174 vm.Value
	var skip_175 int
	var v190 bool
	var argc_178 vm.Value
	var doseq_loop__12013_179 vm.Value
	var k_180 int
	var args_181 vm.Value
	var l_182 vm.Value
	var skip_183 int
	var arg__12103_193 vm.Value
	var arg__12111_196 vm.Value
	var v197 vm.Value
	var v198 int
	var argc_184 vm.Value
	var doseq_loop__12013_185 vm.Value
	var k_186 int
	var args_187 vm.Value
	var l_188 vm.Value
	var skip_189 int
	var v202 vm.Value
	var argc_203 vm.Value
	var doseq_loop__12013_204 vm.Value
	var k_205 int
	var args_206 vm.Value
	var l_207 vm.Value
	var skip_208 int
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = argc_4, k_5, argc_6, args_7, l_8, v19, k_11, argc_12, args_13, l_14, k_15, argc_16, args_17, l_18, arg__12025_23, arg__12033_26, pos_27, skip_72, k_73, argc_74, args_75, l_76, v87, k_28, argc_29, args_30, l_31, pos_32, v62, k_33, argc_34, args_35, l_36, pos_37, v65, k_66, argc_67, args_68, l_69, pos_70, k_38, argc_39, args_40, l_41, and__x_42, pos_43, v51, k_44, argc_45, args_46, l_47, and__x_48, pos_49, v54, k_55, argc_56, args_57, l_58, and__x_59, pos_60, skip_77, k_78, argc_79, args_80, l_81, v90, skip_82, k_83, argc_84, args_85, l_86, v103, v226, skip_227, k_228, argc_229, args_230, l_231, skip_92, k_93, argc_94, args_95, l_96, v106, skip_97, k_98, argc_99, args_100, l_101, v219, skip_220, k_221, argc_222, args_223, l_224, skip_108, k_109, argc_110, args_111, l_112, arg__12057_121, arg__12063_124, arg__12065_125, arg__12070_128, arg__12076_131, arg__12078_132, doseq_seq__12012_133, skip_113, k_114, argc_115, args_116, l_117, v212, skip_213, k_214, argc_215, args_216, l_217, doseq_loop__12013_134, l_135, skip_137, k_138, argc_139, args_140, doseq_seq__12012_141, doseq_loop__12013_142, l_143, a_153, v155, v157, v159, skip_144, k_145, argc_146, args_147, doseq_seq__12012_148, doseq_loop__12013_149, l_150, v163, skip_164, k_165, argc_166, args_167, doseq_seq__12012_168, doseq_loop__12013_169, l_170, k_171, k_172, args_173, l_174, skip_175, v190, argc_178, doseq_loop__12013_179, k_180, args_181, l_182, skip_183, arg__12103_193, arg__12111_196, v197, v198, argc_184, doseq_loop__12013_185, k_186, args_187, l_188, skip_189, v202, argc_203, doseq_loop__12013_204, k_205, args_206, l_207, skip_208
	argc_4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	k_5 = 0
	argc_6 = argc_4
	args_7 = arg1
	l_8 = arg0
	goto b1
b1:
	;
	v19 = rt.GeValue(vm.Int(k_5), argc_6)
	if v19 {
		k_11 = k_5
		argc_12 = argc_6
		args_13 = args_7
		l_14 = l_8
		goto b2
	} else {
		k_15 = k_5
		argc_16 = argc_6
		args_17 = args_7
		l_18 = l_8
		goto b3
	}
b2:
	;
	skip_72 = k_11
	k_73 = k_11
	argc_74 = argc_12
	args_75 = args_13
	l_76 = l_14
	goto b4
b3:
	;
	arg__12025_23, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{args_17, vm.Int(k_15)})
	if callErr != nil {
		return nil, callErr
	}
	arg__12033_26, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{args_17, vm.Int(k_15)})
	if callErr != nil {
		return nil, callErr
	}
	pos_27, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "value-pos-of").Deref(), []vm.Value{l_18, arg__12033_26})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(pos_27) {
		k_38 = k_15
		argc_39 = argc_16
		args_40 = args_17
		l_41 = l_18
		and__x_42 = pos_27
		pos_43 = pos_27
		goto b8
	} else {
		k_44 = k_15
		argc_45 = argc_16
		args_46 = args_17
		l_47 = l_18
		and__x_48 = pos_27
		pos_49 = pos_27
		goto b9
	}
b4:
	;
	v87 = vm.Int(skip_72) == argc_74
	if v87 {
		skip_77 = skip_72
		k_78 = k_73
		argc_79 = argc_74
		args_80 = args_75
		l_81 = l_76
		goto b11
	} else {
		skip_82 = skip_72
		k_83 = k_73
		argc_84 = argc_74
		args_85 = args_75
		l_86 = l_76
		goto b12
	}
b5:
	;
	v62 = k_28 + 1
	k_5 = v62
	argc_6 = argc_29
	args_7 = args_30
	l_8 = l_31
	goto b1
b6:
	;
	v65 = k_33
	k_66 = k_33
	argc_67 = argc_34
	args_68 = args_35
	l_69 = l_36
	pos_70 = pos_37
	goto b7
b7:
	;
	skip_72 = v65
	k_73 = k_66
	argc_74 = argc_67
	args_75 = args_68
	l_76 = l_69
	goto b4
b8:
	;
	v51 = pos_43 == vm.Int(k_38)
	v54 = vm.Boolean(v51)
	k_55 = k_38
	argc_56 = argc_39
	args_57 = args_40
	l_58 = l_41
	and__x_59 = and__x_42
	pos_60 = pos_43
	goto b10
b9:
	;
	v54 = and__x_48
	k_55 = k_44
	argc_56 = argc_45
	args_57 = args_46
	l_58 = l_47
	and__x_59 = and__x_48
	pos_60 = pos_49
	goto b10
b10:
	;
	if vm.IsTruthy(v54) {
		k_28 = k_55
		argc_29 = argc_56
		args_30 = args_57
		l_31 = l_58
		pos_32 = pos_60
		goto b5
	} else {
		k_33 = k_55
		argc_34 = argc_56
		args_35 = args_57
		l_36 = l_58
		pos_37 = pos_60
		goto b6
	}
b11:
	;
	v90, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "consume-refs-in-place!").Deref(), []vm.Value{l_81, args_80})
	if callErr != nil {
		return nil, callErr
	}
	v226 = v90
	skip_227 = skip_77
	k_228 = k_78
	argc_229 = argc_79
	args_230 = args_80
	l_231 = l_81
	goto b13
b12:
	;
	v103, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "args-at-top?").Deref(), []vm.Value{l_86, args_85})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v103) {
		skip_92 = skip_82
		k_93 = k_83
		argc_94 = argc_84
		args_95 = args_85
		l_96 = l_86
		goto b14
	} else {
		skip_97 = skip_82
		k_98 = k_83
		argc_99 = argc_84
		args_100 = args_85
		l_101 = l_86
		goto b15
	}
b13:
	;
	return v226, nil
b14:
	;
	v106, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "consume-refs-in-place!").Deref(), []vm.Value{l_96, args_95})
	if callErr != nil {
		return nil, callErr
	}
	v219 = v106
	skip_220 = skip_92
	k_221 = k_93
	argc_222 = argc_94
	args_223 = args_95
	l_224 = l_96
	goto b16
b15:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		skip_108 = skip_97
		k_109 = k_98
		argc_110 = argc_99
		args_111 = args_100
		l_112 = l_101
		goto b17
	} else {
		skip_113 = skip_97
		k_114 = k_98
		argc_115 = argc_99
		args_116 = args_100
		l_117 = l_101
		goto b18
	}
b16:
	;
	v226 = v219
	skip_227 = skip_220
	k_228 = k_221
	argc_229 = argc_222
	args_230 = args_223
	l_231 = l_224
	goto b13
b17:
	;
	arg__12057_121, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{args_111})
	if callErr != nil {
		return nil, callErr
	}
	arg__12063_124, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{args_111})
	if callErr != nil {
		return nil, callErr
	}
	arg__12065_125, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "subvec").Deref(), []vm.Value{arg__12063_124, vm.Int(skip_108)})
	if callErr != nil {
		return nil, callErr
	}
	arg__12070_128, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{args_111})
	if callErr != nil {
		return nil, callErr
	}
	arg__12076_131, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{args_111})
	if callErr != nil {
		return nil, callErr
	}
	arg__12078_132, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "subvec").Deref(), []vm.Value{arg__12076_131, vm.Int(skip_108)})
	if callErr != nil {
		return nil, callErr
	}
	doseq_seq__12012_133, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{arg__12078_132})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__12013_134 = doseq_seq__12012_133
	l_135 = l_112
	goto b20
b18:
	;
	v212 = vm.NIL
	skip_213 = skip_113
	k_214 = k_114
	argc_215 = argc_115
	args_216 = args_116
	l_217 = l_117
	goto b19
b19:
	;
	v219 = v212
	skip_220 = skip_213
	k_221 = k_214
	argc_222 = argc_215
	args_223 = args_216
	l_224 = l_217
	goto b16
b20:
	;
	if vm.IsTruthy(doseq_loop__12013_134) {
		skip_137 = skip_108
		k_138 = k_109
		argc_139 = argc_110
		args_140 = args_111
		doseq_seq__12012_141 = doseq_seq__12012_133
		doseq_loop__12013_142 = doseq_loop__12013_134
		l_143 = l_135
		goto b21
	} else {
		skip_144 = skip_108
		k_145 = k_109
		argc_146 = argc_110
		args_147 = args_111
		doseq_seq__12012_148 = doseq_seq__12012_133
		doseq_loop__12013_149 = doseq_loop__12013_134
		l_150 = l_135
		goto b22
	}
b21:
	;
	a_153, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{doseq_loop__12013_142})
	if callErr != nil {
		return nil, callErr
	}
	v155, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "materialize!").Deref(), []vm.Value{l_143, a_153})
	if callErr != nil {
		return nil, callErr
	}
	v157, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "decrement-use!").Deref(), []vm.Value{l_143, a_153})
	if callErr != nil {
		return nil, callErr
	}
	v159, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "next").Deref(), []vm.Value{doseq_loop__12013_142})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__12013_134 = v159
	l_135 = l_143
	goto b20
b22:
	;
	v163 = vm.NIL
	skip_164 = skip_144
	k_165 = k_145
	argc_166 = argc_146
	args_167 = args_147
	doseq_seq__12012_168 = doseq_seq__12012_148
	doseq_loop__12013_169 = doseq_loop__12013_149
	l_170 = l_150
	goto b23
b23:
	;
	k_171 = 0
	k_172 = k_165
	args_173 = args_167
	l_174 = l_170
	skip_175 = skip_164
	goto b24
b24:
	;
	v190 = k_172 < skip_175
	if v190 {
		argc_178 = argc_166
		doseq_loop__12013_179 = doseq_loop__12013_169
		k_180 = k_172
		args_181 = args_173
		l_182 = l_174
		skip_183 = skip_175
		goto b25
	} else {
		argc_184 = argc_166
		doseq_loop__12013_185 = doseq_loop__12013_169
		k_186 = k_172
		args_187 = args_173
		l_188 = l_174
		skip_189 = skip_175
		goto b26
	}
b25:
	;
	arg__12103_193, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{args_181, vm.Int(k_180)})
	if callErr != nil {
		return nil, callErr
	}
	arg__12111_196, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{args_181, vm.Int(k_180)})
	if callErr != nil {
		return nil, callErr
	}
	v197, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "decrement-use!").Deref(), []vm.Value{l_182, arg__12111_196})
	if callErr != nil {
		return nil, callErr
	}
	v198 = k_180 + 1
	k_171 = v198
	k_172 = k_180
	args_173 = args_181
	l_174 = l_182
	skip_175 = skip_183
	goto b24
b26:
	;
	v202 = vm.NIL
	argc_203 = argc_184
	doseq_loop__12013_204 = doseq_loop__12013_185
	k_205 = k_186
	args_206 = args_187
	l_207 = l_188
	skip_208 = skip_189
	goto b27
b27:
	;
	v212 = v202
	skip_213 = skip_208
	k_214 = k_205
	argc_215 = argc_203
	args_216 = args_206
	l_217 = l_207
	goto b19
}
func materialize_refs_BANG_(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var doseq_seq__12113_4 vm.Value
	var doseq_loop__12114_5 vm.Value
	var l_6 vm.Value
	var refs_8 vm.Value
	var doseq_seq__12113_9 vm.Value
	var doseq_loop__12114_10 vm.Value
	var l_11 vm.Value
	var r_18 vm.Value
	var v20 vm.Value
	var v22 vm.Value
	var v24 vm.Value
	var refs_12 vm.Value
	var doseq_seq__12113_13 vm.Value
	var doseq_loop__12114_14 vm.Value
	var l_15 vm.Value
	var v28 vm.Value
	var refs_29 vm.Value
	var doseq_seq__12113_30 vm.Value
	var doseq_loop__12114_31 vm.Value
	var l_32 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = doseq_seq__12113_4, doseq_loop__12114_5, l_6, refs_8, doseq_seq__12113_9, doseq_loop__12114_10, l_11, r_18, v20, v22, v24, refs_12, doseq_seq__12113_13, doseq_loop__12114_14, l_15, v28, refs_29, doseq_seq__12113_30, doseq_loop__12114_31, l_32
	doseq_seq__12113_4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__12114_5 = doseq_seq__12113_4
	l_6 = arg0
	goto b1
b1:
	;
	if vm.IsTruthy(doseq_loop__12114_5) {
		refs_8 = arg1
		doseq_seq__12113_9 = doseq_seq__12113_4
		doseq_loop__12114_10 = doseq_loop__12114_5
		l_11 = l_6
		goto b2
	} else {
		refs_12 = arg1
		doseq_seq__12113_13 = doseq_seq__12113_4
		doseq_loop__12114_14 = doseq_loop__12114_5
		l_15 = l_6
		goto b3
	}
b2:
	;
	r_18, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{doseq_loop__12114_10})
	if callErr != nil {
		return nil, callErr
	}
	v20, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "materialize!").Deref(), []vm.Value{l_11, r_18})
	if callErr != nil {
		return nil, callErr
	}
	v22, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "decrement-use!").Deref(), []vm.Value{l_11, r_18})
	if callErr != nil {
		return nil, callErr
	}
	v24, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "next").Deref(), []vm.Value{doseq_loop__12114_10})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__12114_5 = v24
	l_6 = l_11
	goto b1
b3:
	;
	v28 = vm.NIL
	refs_29 = refs_12
	doseq_seq__12113_30 = doseq_seq__12113_13
	doseq_loop__12114_31 = doseq_loop__12114_14
	l_32 = l_15
	goto b4
b4:
	;
	return v28, nil
}
func new_lowerer(arg0 vm.Value) (vm.Value, error) {
	var uses_3 vm.Value
	var arg__12140_5 vm.Value
	var arg__12145_8 vm.Value
	var n_blocks_9 vm.Value
	var i_10 int
	var acc_11 vm.Value
	var uses_12 vm.Value
	var arg__12150_27 vm.Value
	var v28 bool
	var f_16 vm.Value
	var n_blocks_17 vm.Value
	var i_18 int
	var acc_19 vm.Value
	var uses_20 vm.Value
	var f_21 vm.Value
	var n_blocks_22 vm.Value
	var i_23 int
	var acc_24 vm.Value
	var uses_25 vm.Value
	var us_32 vm.Value
	var v33 int
	var v47 vm.Value
	var use_count_65 vm.Value
	var f_66 vm.Value
	var n_blocks_67 vm.Value
	var i_68 int
	var acc_69 vm.Value
	var uses_70 vm.Value
	var arg__12180_76 vm.Value
	var arg__12185_79 vm.Value
	var arg__12186_80 vm.Value
	var arg__12197_87 vm.Value
	var arg__12204_92 vm.Value
	var arg__12205_93 vm.Value
	var arg__12212_98 vm.Value
	var arg__12219_103 vm.Value
	var arg__12220_104 vm.Value
	var arg__12229_112 vm.Value
	var arg__12238_119 vm.Value
	var arg__12243_122 vm.Value
	var arg__12244_123 vm.Value
	var arg__12255_130 vm.Value
	var arg__12262_135 vm.Value
	var arg__12263_136 vm.Value
	var arg__12270_141 vm.Value
	var arg__12277_146 vm.Value
	var arg__12278_147 vm.Value
	var arg__12287_155 vm.Value
	var v156 vm.Value
	var f_34 vm.Value
	var n_blocks_35 vm.Value
	var i_36 int
	var acc_37 vm.Value
	var uses_38 vm.Value
	var us_39 vm.Value
	var f_40 vm.Value
	var n_blocks_41 vm.Value
	var i_42 int
	var acc_43 vm.Value
	var uses_44 vm.Value
	var us_45 vm.Value
	var arg__12165_51 vm.Value
	var arg__12172_54 vm.Value
	var v55 vm.Value
	var v57 vm.Value
	var f_58 vm.Value
	var n_blocks_59 vm.Value
	var i_60 int
	var acc_61 vm.Value
	var uses_62 vm.Value
	var us_63 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = uses_3, arg__12140_5, arg__12145_8, n_blocks_9, i_10, acc_11, uses_12, arg__12150_27, v28, f_16, n_blocks_17, i_18, acc_19, uses_20, f_21, n_blocks_22, i_23, acc_24, uses_25, us_32, v33, v47, use_count_65, f_66, n_blocks_67, i_68, acc_69, uses_70, arg__12180_76, arg__12185_79, arg__12186_80, arg__12197_87, arg__12204_92, arg__12205_93, arg__12212_98, arg__12219_103, arg__12220_104, arg__12229_112, arg__12238_119, arg__12243_122, arg__12244_123, arg__12255_130, arg__12262_135, arg__12263_136, arg__12270_141, arg__12277_146, arg__12278_147, arg__12287_155, v156, f_34, n_blocks_35, i_36, acc_37, uses_38, us_39, f_40, n_blocks_41, i_42, acc_43, uses_44, us_45, arg__12165_51, arg__12172_54, v55, v57, f_58, n_blocks_59, i_60, acc_61, uses_62, us_63
	uses_3, callErr = rt.InvokeValue(rt.LookupVar("ir", "uses").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__12140_5, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__12145_8, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	n_blocks_9, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{arg__12145_8})
	if callErr != nil {
		return nil, callErr
	}
	i_10 = 0
	acc_11 = vm.EmptyPersistentMap
	uses_12 = uses_3
	goto b1
b1:
	;
	arg__12150_27, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{uses_12})
	if callErr != nil {
		return nil, callErr
	}
	v28 = rt.GeValue(vm.Int(i_10), arg__12150_27)
	if v28 {
		f_16 = arg0
		n_blocks_17 = n_blocks_9
		i_18 = i_10
		acc_19 = acc_11
		uses_20 = uses_12
		goto b2
	} else {
		f_21 = arg0
		n_blocks_22 = n_blocks_9
		i_23 = i_10
		acc_24 = acc_11
		uses_25 = uses_12
		goto b3
	}
b2:
	;
	use_count_65 = acc_19
	f_66 = f_16
	n_blocks_67 = n_blocks_17
	i_68 = i_18
	acc_69 = acc_19
	uses_70 = uses_20
	goto b4
b3:
	;
	us_32, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{uses_25, vm.Int(i_23)})
	if callErr != nil {
		return nil, callErr
	}
	v33 = i_23 + 1
	v47, callErr = rt.InvokeValue(rt.LookupVar("ir", "uses-empty?").Deref(), []vm.Value{us_32})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v47) {
		f_34 = f_21
		n_blocks_35 = n_blocks_22
		i_36 = i_23
		acc_37 = acc_24
		uses_38 = uses_25
		us_39 = us_32
		goto b5
	} else {
		f_40 = f_21
		n_blocks_41 = n_blocks_22
		i_42 = i_23
		acc_43 = acc_24
		uses_44 = uses_25
		us_45 = us_32
		goto b6
	}
b4:
	;
	arg__12180_76, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-consts").Deref(), []vm.Value{f_66})
	if callErr != nil {
		return nil, callErr
	}
	arg__12185_79, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-consts").Deref(), []vm.Value{f_66})
	if callErr != nil {
		return nil, callErr
	}
	arg__12186_80, callErr = rt.InvokeValue(rt.LookupVar("ir", "new-chunk").Deref(), []vm.Value{arg__12185_79})
	if callErr != nil {
		return nil, callErr
	}
	arg__12197_87, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "repeat").Deref(), []vm.Value{n_blocks_67, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	arg__12204_92, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "repeat").Deref(), []vm.Value{n_blocks_67, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	arg__12205_93, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg__12204_92})
	if callErr != nil {
		return nil, callErr
	}
	arg__12212_98, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "repeat").Deref(), []vm.Value{n_blocks_67, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	arg__12219_103, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "repeat").Deref(), []vm.Value{n_blocks_67, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	arg__12220_104, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg__12219_103})
	if callErr != nil {
		return nil, callErr
	}
	arg__12229_112, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("value-stack-pos"), vm.EmptyPersistentMap, vm.Keyword("chunk"), arg__12186_80, vm.Keyword("uses"), uses_70, vm.Keyword("f"), f_66, vm.Keyword("block-junk"), arg__12205_93, vm.Keyword("block-ips"), arg__12220_104, vm.Keyword("patches"), vm.NewArrayVector([]vm.Value{}), vm.Keyword("use-count"), use_count_65, vm.Keyword("current-block"), vm.Int(0), vm.Keyword("stack-sp"), vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	arg__12238_119, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-consts").Deref(), []vm.Value{f_66})
	if callErr != nil {
		return nil, callErr
	}
	arg__12243_122, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-consts").Deref(), []vm.Value{f_66})
	if callErr != nil {
		return nil, callErr
	}
	arg__12244_123, callErr = rt.InvokeValue(rt.LookupVar("ir", "new-chunk").Deref(), []vm.Value{arg__12243_122})
	if callErr != nil {
		return nil, callErr
	}
	arg__12255_130, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "repeat").Deref(), []vm.Value{n_blocks_67, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	arg__12262_135, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "repeat").Deref(), []vm.Value{n_blocks_67, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	arg__12263_136, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg__12262_135})
	if callErr != nil {
		return nil, callErr
	}
	arg__12270_141, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "repeat").Deref(), []vm.Value{n_blocks_67, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	arg__12277_146, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "repeat").Deref(), []vm.Value{n_blocks_67, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	arg__12278_147, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg__12277_146})
	if callErr != nil {
		return nil, callErr
	}
	arg__12287_155, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("value-stack-pos"), vm.EmptyPersistentMap, vm.Keyword("chunk"), arg__12244_123, vm.Keyword("uses"), uses_70, vm.Keyword("f"), f_66, vm.Keyword("block-junk"), arg__12263_136, vm.Keyword("block-ips"), arg__12278_147, vm.Keyword("patches"), vm.NewArrayVector([]vm.Value{}), vm.Keyword("use-count"), use_count_65, vm.Keyword("current-block"), vm.Int(0), vm.Keyword("stack-sp"), vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	v156, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "atom").Deref(), []vm.Value{arg__12287_155})
	if callErr != nil {
		return nil, callErr
	}
	return v156, nil
b5:
	;
	v57 = acc_37
	f_58 = f_34
	n_blocks_59 = n_blocks_35
	i_60 = i_36
	acc_61 = acc_37
	uses_62 = uses_38
	us_63 = us_39
	goto b7
b6:
	;
	arg__12165_51, callErr = rt.InvokeValue(rt.LookupVar("ir", "uses-count").Deref(), []vm.Value{us_45})
	if callErr != nil {
		return nil, callErr
	}
	arg__12172_54, callErr = rt.InvokeValue(rt.LookupVar("ir", "uses-count").Deref(), []vm.Value{us_45})
	if callErr != nil {
		return nil, callErr
	}
	v55, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "assoc").Deref(), []vm.Value{acc_43, vm.Int(i_42), arg__12172_54})
	if callErr != nil {
		return nil, callErr
	}
	v57 = v55
	f_58 = f_40
	n_blocks_59 = n_blocks_41
	i_60 = i_42
	acc_61 = acc_43
	uses_62 = uses_44
	us_63 = us_45
	goto b7
b7:
	;
	i_10 = v33
	acc_11 = v57
	uses_12 = uses_25
	goto b1
}
func patch_branches_BANG_(arg0 vm.Value) (vm.Value, error) {
	var arg__12294_3 vm.Value
	var arg__12295_4 vm.Value
	var arg__12301_8 vm.Value
	var arg__12302_9 vm.Value
	var doseq_seq__12288_10 vm.Value
	var doseq_loop__12289_11 vm.Value
	var l_12 vm.Value
	var v86 vm.Value
	var v92 vm.Value
	var v98 vm.Value
	var v104 vm.Value
	var doseq_seq__12288_14 vm.Value
	var doseq_loop__12289_15 vm.Value
	var l_16 vm.Value
	var v85 vm.Value
	var v91 vm.Value
	var v97 vm.Value
	var v103 vm.Value
	var p_22 vm.Value
	var arg__12309_24 vm.Value
	var arg__12314_27 vm.Value
	var target_ip_28 vm.Value
	var src_ip_30 vm.Value
	var v44 vm.Value
	var doseq_seq__12288_17 vm.Value
	var doseq_loop__12289_18 vm.Value
	var l_19 vm.Value
	var v90 vm.Value
	var v96 vm.Value
	var v102 vm.Value
	var v108 vm.Value
	var v74 vm.Value
	var doseq_seq__12288_75 vm.Value
	var doseq_loop__12289_76 vm.Value
	var l_77 vm.Value
	var doseq_seq__12288_31 vm.Value
	var doseq_loop__12289_32 vm.Value
	var l_33 vm.Value
	var p_34 vm.Value
	var target_ip_35 vm.Value
	var src_ip_36 vm.Value
	var v88 vm.Value
	var v94 vm.Value
	var v100 vm.Value
	var v106 vm.Value
	var v46 vm.Value
	var doseq_seq__12288_37 vm.Value
	var doseq_loop__12289_38 vm.Value
	var l_39 vm.Value
	var p_40 vm.Value
	var target_ip_41 vm.Value
	var src_ip_42 vm.Value
	var v87 vm.Value
	var v93 vm.Value
	var v99 vm.Value
	var v105 vm.Value
	var v48 vm.Value
	var offset_50 vm.Value
	var doseq_seq__12288_51 vm.Value
	var doseq_loop__12289_52 vm.Value
	var l_53 vm.Value
	var p_54 vm.Value
	var target_ip_55 vm.Value
	var src_ip_56 vm.Value
	var v89 vm.Value
	var v95 vm.Value
	var v101 vm.Value
	var v107 vm.Value
	var arg__12326_58 vm.Value
	var arg__12330_60 vm.Value
	var arg__12337_64 vm.Value
	var arg__12341_66 vm.Value
	var arg__12342_67 vm.Value
	var v68 vm.Value
	var v70 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__12294_3, arg__12295_4, arg__12301_8, arg__12302_9, doseq_seq__12288_10, doseq_loop__12289_11, l_12, v86, v92, v98, v104, doseq_seq__12288_14, doseq_loop__12289_15, l_16, v85, v91, v97, v103, p_22, arg__12309_24, arg__12314_27, target_ip_28, src_ip_30, v44, doseq_seq__12288_17, doseq_loop__12289_18, l_19, v90, v96, v102, v108, v74, doseq_seq__12288_75, doseq_loop__12289_76, l_77, doseq_seq__12288_31, doseq_loop__12289_32, l_33, p_34, target_ip_35, src_ip_36, v88, v94, v100, v106, v46, doseq_seq__12288_37, doseq_loop__12289_38, l_39, p_40, target_ip_41, src_ip_42, v87, v93, v99, v105, v48, offset_50, doseq_seq__12288_51, doseq_loop__12289_52, l_53, p_54, target_ip_55, src_ip_56, v89, v95, v101, v107, arg__12326_58, arg__12330_60, arg__12337_64, arg__12341_66, arg__12342_67, v68, v70
	arg__12294_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__12295_4, callErr = rt.InvokeValue(vm.Keyword("patches"), []vm.Value{arg__12294_3})
	if callErr != nil {
		return nil, callErr
	}
	arg__12301_8, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__12302_9, callErr = rt.InvokeValue(vm.Keyword("patches"), []vm.Value{arg__12301_8})
	if callErr != nil {
		return nil, callErr
	}
	doseq_seq__12288_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{arg__12302_9})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__12289_11 = doseq_seq__12288_10
	l_12 = arg0
	v86 = vm.Keyword("target-block")
	v92 = vm.Keyword("src-ip")
	v98 = vm.Keyword("negate?")
	v104 = vm.Keyword("offset-slot")
	goto b1
b1:
	;
	if vm.IsTruthy(doseq_loop__12289_11) {
		doseq_seq__12288_14 = doseq_seq__12288_10
		doseq_loop__12289_15 = doseq_loop__12289_11
		l_16 = l_12
		v85 = v86
		v91 = v92
		v97 = v98
		v103 = v104
		goto b2
	} else {
		doseq_seq__12288_17 = doseq_seq__12288_10
		doseq_loop__12289_18 = doseq_loop__12289_11
		l_19 = l_12
		v90 = v86
		v96 = v92
		v102 = v98
		v108 = v104
		goto b3
	}
b2:
	;
	p_22, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{doseq_loop__12289_15})
	if callErr != nil {
		return nil, callErr
	}
	arg__12309_24, callErr = rt.InvokeValue(v85, []vm.Value{p_22})
	if callErr != nil {
		return nil, callErr
	}
	arg__12314_27, callErr = rt.InvokeValue(v85, []vm.Value{p_22})
	if callErr != nil {
		return nil, callErr
	}
	target_ip_28, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "block-ip").Deref(), []vm.Value{l_16, arg__12314_27})
	if callErr != nil {
		return nil, callErr
	}
	src_ip_30, callErr = rt.InvokeValue(v91, []vm.Value{p_22})
	if callErr != nil {
		return nil, callErr
	}
	v44, callErr = rt.InvokeValue(v97, []vm.Value{p_22})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v44) {
		doseq_seq__12288_31 = doseq_seq__12288_14
		doseq_loop__12289_32 = doseq_loop__12289_15
		l_33 = l_16
		p_34 = p_22
		target_ip_35 = target_ip_28
		src_ip_36 = src_ip_30
		v88 = v85
		v94 = v91
		v100 = v97
		v106 = v103
		goto b5
	} else {
		doseq_seq__12288_37 = doseq_seq__12288_14
		doseq_loop__12289_38 = doseq_loop__12289_15
		l_39 = l_16
		p_40 = p_22
		target_ip_41 = target_ip_28
		src_ip_42 = src_ip_30
		v87 = v85
		v93 = v91
		v99 = v97
		v105 = v103
		goto b6
	}
b3:
	;
	v74 = vm.NIL
	doseq_seq__12288_75 = doseq_seq__12288_17
	doseq_loop__12289_76 = doseq_loop__12289_18
	l_77 = l_19
	goto b4
b4:
	;
	return v74, nil
b5:
	;
	v46 = rt.SubValue(src_ip_36, target_ip_35)
	offset_50 = v46
	doseq_seq__12288_51 = doseq_seq__12288_31
	doseq_loop__12289_52 = doseq_loop__12289_32
	l_53 = l_33
	p_54 = p_34
	target_ip_55 = target_ip_35
	src_ip_56 = src_ip_36
	v89 = v88
	v95 = v94
	v101 = v100
	v107 = v106
	goto b7
b6:
	;
	v48 = rt.SubValue(target_ip_41, src_ip_42)
	offset_50 = v48
	doseq_seq__12288_51 = doseq_seq__12288_37
	doseq_loop__12289_52 = doseq_loop__12289_38
	l_53 = l_39
	p_54 = p_40
	target_ip_55 = target_ip_41
	src_ip_56 = src_ip_42
	v89 = v87
	v95 = v93
	v101 = v99
	v107 = v105
	goto b7
b7:
	;
	arg__12326_58, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "chunk-of").Deref(), []vm.Value{l_53})
	if callErr != nil {
		return nil, callErr
	}
	arg__12330_60, callErr = rt.InvokeValue(v107, []vm.Value{p_54})
	if callErr != nil {
		return nil, callErr
	}
	arg__12337_64, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "chunk-of").Deref(), []vm.Value{l_53})
	if callErr != nil {
		return nil, callErr
	}
	arg__12341_66, callErr = rt.InvokeValue(v107, []vm.Value{p_54})
	if callErr != nil {
		return nil, callErr
	}
	arg__12342_67 = rt.AddValue(src_ip_56, arg__12341_66)
	v68, callErr = rt.InvokeValue(rt.LookupVar("ir", "chunk-update!").Deref(), []vm.Value{arg__12337_64, arg__12342_67, offset_50})
	if callErr != nil {
		return nil, callErr
	}
	v70, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "next").Deref(), []vm.Value{doseq_loop__12289_52})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__12289_11 = v70
	l_12 = l_53
	v86 = v89
	v92 = v95
	v98 = v101
	v104 = v107
	goto b1
}
func record_block_junk_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 int) (vm.Value, error) {
	var cur_5 vm.Value
	var v14 bool
	var l_6 vm.Value
	var bid_7 vm.Value
	var junk_8 int
	var cur_9 vm.Value
	var v23 vm.Value
	var l_10 vm.Value
	var bid_11 vm.Value
	var junk_12 int
	var cur_13 vm.Value
	var v27 vm.Value
	var l_28 vm.Value
	var bid_29 vm.Value
	var junk_30 int
	var cur_31 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = cur_5, v14, l_6, bid_7, junk_8, cur_9, v23, l_10, bid_11, junk_12, cur_13, v27, l_28, bid_29, junk_30, cur_31
	cur_5, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "block-junk-of").Deref(), []vm.Value{arg0, arg1})
	if callErr != nil {
		return nil, callErr
	}
	v14 = rt.GtValue(vm.Int(arg2), cur_5)
	if v14 {
		l_6 = arg0
		bid_7 = arg1
		junk_8 = arg2
		cur_9 = cur_5
		goto b1
	} else {
		l_10 = arg0
		bid_11 = arg1
		junk_12 = arg2
		cur_13 = cur_5
		goto b2
	}
b1:
	;
	v23, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{l_6, rt.LookupVar("clojure.core", "update").Deref(), vm.Keyword("block-junk"), rt.LookupVar("clojure.core", "assoc").Deref(), bid_7, vm.Int(junk_8)})
	if callErr != nil {
		return nil, callErr
	}
	v27 = v23
	l_28 = l_6
	bid_29 = bid_7
	junk_30 = junk_8
	cur_31 = cur_9
	goto b3
b2:
	;
	v27 = vm.NIL
	l_28 = l_10
	bid_29 = bid_11
	junk_30 = junk_12
	cur_31 = cur_13
	goto b3
b3:
	;
	return v27, nil
}
func record_source_info_BANG_(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var l_2 vm.Value
	var inst_id_3 vm.Value
	var c_8 vm.Value
	var arg__12376_10 vm.Value
	var arg__12382_13 vm.Value
	var arg__12383_14 vm.Value
	var arg__12389_17 vm.Value
	var arg__12395_20 vm.Value
	var arg__12396_21 vm.Value
	var doseq_seq__12367_22 vm.Value
	var l_4 vm.Value
	var inst_id_5 vm.Value
	var v55 vm.Value
	var l_56 vm.Value
	var inst_id_57 vm.Value
	var doseq_loop__12368_23 vm.Value
	var c_24 vm.Value
	var l_26 vm.Value
	var inst_id_27 vm.Value
	var doseq_seq__12367_28 vm.Value
	var doseq_loop__12368_29 vm.Value
	var c_30 vm.Value
	var si_38 vm.Value
	var v40 vm.Value
	var v42 vm.Value
	var l_31 vm.Value
	var inst_id_32 vm.Value
	var doseq_seq__12367_33 vm.Value
	var doseq_loop__12368_34 vm.Value
	var c_35 vm.Value
	var v46 vm.Value
	var l_47 vm.Value
	var inst_id_48 vm.Value
	var doseq_seq__12367_49 vm.Value
	var doseq_loop__12368_50 vm.Value
	var c_51 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = l_2, inst_id_3, c_8, arg__12376_10, arg__12382_13, arg__12383_14, arg__12389_17, arg__12395_20, arg__12396_21, doseq_seq__12367_22, l_4, inst_id_5, v55, l_56, inst_id_57, doseq_loop__12368_23, c_24, l_26, inst_id_27, doseq_seq__12367_28, doseq_loop__12368_29, c_30, si_38, v40, v42, l_31, inst_id_32, doseq_seq__12367_33, doseq_loop__12368_34, c_35, v46, l_47, inst_id_48, doseq_seq__12367_49, doseq_loop__12368_50, c_51
	if vm.IsTruthy(arg1) {
		l_2 = arg0
		inst_id_3 = arg1
		goto b1
	} else {
		l_4 = arg0
		inst_id_5 = arg1
		goto b2
	}
b1:
	;
	c_8, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "chunk-of").Deref(), []vm.Value{l_2})
	if callErr != nil {
		return nil, callErr
	}
	arg__12376_10, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "f-of").Deref(), []vm.Value{l_2})
	if callErr != nil {
		return nil, callErr
	}
	arg__12382_13, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "f-of").Deref(), []vm.Value{l_2})
	if callErr != nil {
		return nil, callErr
	}
	arg__12383_14, callErr = rt.InvokeValue(rt.LookupVar("ir", "source-infos").Deref(), []vm.Value{inst_id_3, arg__12382_13})
	if callErr != nil {
		return nil, callErr
	}
	arg__12389_17, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "f-of").Deref(), []vm.Value{l_2})
	if callErr != nil {
		return nil, callErr
	}
	arg__12395_20, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "f-of").Deref(), []vm.Value{l_2})
	if callErr != nil {
		return nil, callErr
	}
	arg__12396_21, callErr = rt.InvokeValue(rt.LookupVar("ir", "source-infos").Deref(), []vm.Value{inst_id_3, arg__12395_20})
	if callErr != nil {
		return nil, callErr
	}
	doseq_seq__12367_22, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{arg__12396_21})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__12368_23 = doseq_seq__12367_22
	c_24 = c_8
	goto b4
b2:
	;
	v55 = vm.NIL
	l_56 = l_4
	inst_id_57 = inst_id_5
	goto b3
b3:
	;
	return v55, nil
b4:
	;
	if vm.IsTruthy(doseq_loop__12368_23) {
		l_26 = l_2
		inst_id_27 = inst_id_3
		doseq_seq__12367_28 = doseq_seq__12367_22
		doseq_loop__12368_29 = doseq_loop__12368_23
		c_30 = c_24
		goto b5
	} else {
		l_31 = l_2
		inst_id_32 = inst_id_3
		doseq_seq__12367_33 = doseq_seq__12367_22
		doseq_loop__12368_34 = doseq_loop__12368_23
		c_35 = c_24
		goto b6
	}
b5:
	;
	si_38, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{doseq_loop__12368_29})
	if callErr != nil {
		return nil, callErr
	}
	v40, callErr = rt.InvokeValue(rt.LookupVar("ir", "chunk-add-source-info!").Deref(), []vm.Value{c_30, si_38})
	if callErr != nil {
		return nil, callErr
	}
	v42, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "next").Deref(), []vm.Value{doseq_loop__12368_29})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__12368_23 = v42
	c_24 = c_30
	goto b4
b6:
	;
	v46 = vm.NIL
	l_47 = l_31
	inst_id_48 = inst_id_32
	doseq_seq__12367_49 = doseq_seq__12367_33
	doseq_loop__12368_50 = doseq_loop__12368_34
	c_51 = c_35
	goto b7
b7:
	;
	v55 = v46
	l_56 = l_47
	inst_id_57 = inst_id_48
	goto b3
}
func refs_at_top_last_use_QMARK_(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__12411_8 vm.Value
	var arg__12416_11 vm.Value
	var v12 vm.Value
	var l_3 vm.Value
	var refs_4 vm.Value
	var l_5 vm.Value
	var refs_6 vm.Value
	var arg__12422_21 vm.Value
	var arg__12429_24 vm.Value
	var v25 vm.Value
	var v134 vm.Value
	var l_135 vm.Value
	var refs_136 vm.Value
	var l_16 vm.Value
	var refs_17 vm.Value
	var l_18 vm.Value
	var refs_19 vm.Value
	var v130 vm.Value
	var l_131 vm.Value
	var refs_132 vm.Value
	var l_29 vm.Value
	var refs_30 vm.Value
	var l_31 vm.Value
	var refs_32 vm.Value
	var v126 vm.Value
	var l_127 vm.Value
	var refs_128 vm.Value
	var i_35 int
	var refs_36 vm.Value
	var l_37 vm.Value
	var v143 int
	var v152 vm.Value
	var arg__12434_47 vm.Value
	var v48 bool
	var i_40 int
	var refs_41 vm.Value
	var l_42 vm.Value
	var v147 int
	var v156 vm.Value
	var i_43 int
	var refs_44 vm.Value
	var l_45 vm.Value
	var v145 int
	var v154 vm.Value
	var arg__12442_60 vm.Value
	var arg__12450_63 vm.Value
	var arg__12451_64 vm.Value
	var arg__12460_68 vm.Value
	var arg__12468_71 vm.Value
	var arg__12469_72 vm.Value
	var v73 vm.Value
	var v119 vm.Value
	var i_120 int
	var refs_121 vm.Value
	var l_122 vm.Value
	var i_52 int
	var refs_53 vm.Value
	var l_54 vm.Value
	var v149 int
	var v158 vm.Value
	var i_55 int
	var refs_56 vm.Value
	var l_57 vm.Value
	var v142 int
	var v151 vm.Value
	var arg__12477_84 vm.Value
	var arg__12486_87 vm.Value
	var v88 vm.Value
	var v114 vm.Value
	var i_115 int
	var refs_116 vm.Value
	var l_117 vm.Value
	var i_77 int
	var refs_78 vm.Value
	var l_79 vm.Value
	var v146 int
	var v155 vm.Value
	var i_80 int
	var refs_81 vm.Value
	var l_82 vm.Value
	var v144 int
	var v153 vm.Value
	var v109 vm.Value
	var i_110 int
	var refs_111 vm.Value
	var l_112 vm.Value
	var i_92 int
	var refs_93 vm.Value
	var l_94 vm.Value
	var v141 int
	var v150 vm.Value
	var v100 int
	var i_95 int
	var refs_96 vm.Value
	var l_97 vm.Value
	var v148 int
	var v157 vm.Value
	var v104 vm.Value
	var i_105 int
	var refs_106 vm.Value
	var l_107 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__12411_8, arg__12416_11, v12, l_3, refs_4, l_5, refs_6, arg__12422_21, arg__12429_24, v25, v134, l_135, refs_136, l_16, refs_17, l_18, refs_19, v130, l_131, refs_132, l_29, refs_30, l_31, refs_32, v126, l_127, refs_128, i_35, refs_36, l_37, v143, v152, arg__12434_47, v48, i_40, refs_41, l_42, v147, v156, i_43, refs_44, l_45, v145, v154, arg__12442_60, arg__12450_63, arg__12451_64, arg__12460_68, arg__12468_71, arg__12469_72, v73, v119, i_120, refs_121, l_122, i_52, refs_53, l_54, v149, v158, i_55, refs_56, l_57, v142, v151, arg__12477_84, arg__12486_87, v88, v114, i_115, refs_116, l_117, i_77, refs_78, l_79, v146, v155, i_80, refs_81, l_82, v144, v153, v109, i_110, refs_111, l_112, i_92, refs_93, l_94, v141, v150, v100, i_95, refs_96, l_97, v148, v157, v104, i_105, refs_106, l_107
	arg__12411_8, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__12416_11, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	v12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "zero?").Deref(), []vm.Value{arg__12416_11})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v12) {
		l_3 = arg0
		refs_4 = arg1
		goto b1
	} else {
		l_5 = arg0
		refs_6 = arg1
		goto b2
	}
b1:
	;
	v134 = vm.Boolean(true)
	l_135 = l_3
	refs_136 = refs_4
	goto b3
b2:
	;
	arg__12422_21, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "args-at-top?").Deref(), []vm.Value{l_5, refs_6})
	if callErr != nil {
		return nil, callErr
	}
	arg__12429_24, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "args-at-top?").Deref(), []vm.Value{l_5, refs_6})
	if callErr != nil {
		return nil, callErr
	}
	v25, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not").Deref(), []vm.Value{arg__12429_24})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v25) {
		l_16 = l_5
		refs_17 = refs_6
		goto b4
	} else {
		l_18 = l_5
		refs_19 = refs_6
		goto b5
	}
b3:
	;
	return v134, nil
b4:
	;
	v130 = vm.Boolean(false)
	l_131 = l_16
	refs_132 = refs_17
	goto b6
b5:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		l_29 = l_18
		refs_30 = refs_19
		goto b7
	} else {
		l_31 = l_18
		refs_32 = refs_19
		goto b8
	}
b6:
	;
	v134 = v130
	l_135 = l_131
	refs_136 = refs_132
	goto b3
b7:
	;
	i_35 = 0
	refs_36 = refs_30
	l_37 = l_29
	v143 = 1
	v152 = vm.Keyword("else")
	goto b10
b8:
	;
	v126 = vm.NIL
	l_127 = l_31
	refs_128 = refs_32
	goto b9
b9:
	;
	v130 = v126
	l_131 = l_127
	refs_132 = refs_128
	goto b6
b10:
	;
	arg__12434_47, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{refs_36})
	if callErr != nil {
		return nil, callErr
	}
	v48 = rt.GeValue(vm.Int(i_35), arg__12434_47)
	if v48 {
		i_40 = i_35
		refs_41 = refs_36
		l_42 = l_37
		v147 = v143
		v156 = v152
		goto b11
	} else {
		i_43 = i_35
		refs_44 = refs_36
		l_45 = l_37
		v145 = v143
		v154 = v152
		goto b12
	}
b11:
	;
	v119 = vm.Boolean(true)
	i_120 = i_40
	refs_121 = refs_41
	l_122 = l_42
	goto b13
b12:
	;
	arg__12442_60, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{refs_44, vm.Int(i_43)})
	if callErr != nil {
		return nil, callErr
	}
	arg__12450_63, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{refs_44, vm.Int(i_43)})
	if callErr != nil {
		return nil, callErr
	}
	arg__12451_64, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "use-count-of").Deref(), []vm.Value{l_45, arg__12450_63})
	if callErr != nil {
		return nil, callErr
	}
	arg__12460_68, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{refs_44, vm.Int(i_43)})
	if callErr != nil {
		return nil, callErr
	}
	arg__12468_71, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{refs_44, vm.Int(i_43)})
	if callErr != nil {
		return nil, callErr
	}
	arg__12469_72, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "use-count-of").Deref(), []vm.Value{l_45, arg__12468_71})
	if callErr != nil {
		return nil, callErr
	}
	v73, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not=").Deref(), []vm.Value{vm.Int(v145), arg__12469_72})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v73) {
		i_52 = i_43
		refs_53 = refs_44
		l_54 = l_45
		v149 = v145
		v158 = v154
		goto b14
	} else {
		i_55 = i_43
		refs_56 = refs_44
		l_57 = l_45
		v142 = v145
		v151 = v154
		goto b15
	}
b13:
	;
	v126 = v119
	l_127 = l_122
	refs_128 = refs_121
	goto b9
b14:
	;
	v114 = vm.Boolean(false)
	i_115 = i_52
	refs_116 = refs_53
	l_117 = l_54
	goto b16
b15:
	;
	arg__12477_84, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{refs_56, vm.Int(i_55)})
	if callErr != nil {
		return nil, callErr
	}
	arg__12486_87, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{refs_56, vm.Int(i_55)})
	if callErr != nil {
		return nil, callErr
	}
	v88, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "contains-after?").Deref(), []vm.Value{refs_56, vm.Int(i_55), arg__12486_87})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v88) {
		i_77 = i_55
		refs_78 = refs_56
		l_79 = l_57
		v146 = v142
		v155 = v151
		goto b17
	} else {
		i_80 = i_55
		refs_81 = refs_56
		l_82 = l_57
		v144 = v142
		v153 = v151
		goto b18
	}
b16:
	;
	v119 = v114
	i_120 = i_115
	refs_121 = refs_116
	l_122 = l_117
	goto b13
b17:
	;
	v109 = vm.Boolean(false)
	i_110 = i_77
	refs_111 = refs_78
	l_112 = l_79
	goto b19
b18:
	;
	if vm.IsTruthy(v153) {
		i_92 = i_80
		refs_93 = refs_81
		l_94 = l_82
		v141 = v144
		v150 = v153
		goto b20
	} else {
		i_95 = i_80
		refs_96 = refs_81
		l_97 = l_82
		v148 = v144
		v157 = v153
		goto b21
	}
b19:
	;
	v114 = v109
	i_115 = i_110
	refs_116 = refs_111
	l_117 = l_112
	goto b16
b20:
	;
	v100 = i_92 + 1
	i_35 = v100
	refs_36 = refs_93
	l_37 = l_94
	v143 = v141
	v152 = v150
	goto b10
b21:
	;
	v104 = vm.NIL
	i_105 = i_95
	refs_106 = refs_96
	l_107 = l_97
	goto b22
b22:
	;
	v109 = v104
	i_110 = i_105
	refs_111 = refs_106
	l_112 = l_107
	goto b19
}
func set_block_ip_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var v16 vm.Value
	var callErr error
	_ = v16
	v16, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "update").Deref(), vm.Keyword("block-ips"), rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v4 vm.Value
		var callErr error
		_ = v4
		v4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "assoc").Deref(), []vm.Value{arg0, arg1, arg2})
		if callErr != nil {
			return nil, callErr
		}
		return v4, nil
	})})
	if callErr != nil {
		return nil, callErr
	}
	return v16, nil
}
func set_stack_sp_BANG_(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var v7 vm.Value
	var callErr error
	_ = v7
	v7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "assoc").Deref(), vm.Keyword("stack-sp"), arg1})
	if callErr != nil {
		return nil, callErr
	}
	return v7, nil
}
func set_value_pos_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var v10 vm.Value
	var callErr error
	_ = v10
	v10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "update").Deref(), vm.Keyword("value-stack-pos"), rt.LookupVar("clojure.core", "assoc").Deref(), arg1, arg2})
	if callErr != nil {
		return nil, callErr
	}
	return v10, nil
}
func should_body_emit_cheap_QMARK_(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var uc_4 vm.Value
	var v12 vm.Value
	var l_5 vm.Value
	var nid_6 vm.Value
	var uc_7 vm.Value
	var l_8 vm.Value
	var nid_9 vm.Value
	var uc_10 vm.Value
	var v23 bool
	var v211 vm.Value
	var l_212 vm.Value
	var nid_213 vm.Value
	var uc_214 vm.Value
	var l_16 vm.Value
	var nid_17 vm.Value
	var uc_18 vm.Value
	var l_19 vm.Value
	var nid_20 vm.Value
	var uc_21 vm.Value
	var v206 vm.Value
	var l_207 vm.Value
	var nid_208 vm.Value
	var uc_209 vm.Value
	var l_27 vm.Value
	var nid_28 vm.Value
	var uc_29 vm.Value
	var arg__12547_37 vm.Value
	var uses_38 vm.Value
	var arg__12552_48 vm.Value
	var v49 bool
	var l_30 vm.Value
	var nid_31 vm.Value
	var uc_32 vm.Value
	var v201 vm.Value
	var l_202 vm.Value
	var nid_203 vm.Value
	var uc_204 vm.Value
	var l_39 vm.Value
	var nid_40 vm.Value
	var uc_41 vm.Value
	var uses_42 vm.Value
	var v52 vm.Value
	var l_43 vm.Value
	var nid_44 vm.Value
	var uc_45 vm.Value
	var uses_46 vm.Value
	var us_56 vm.Value
	var l_57 vm.Value
	var nid_58 vm.Value
	var uc_59 vm.Value
	var uses_60 vm.Value
	var or__x_72 vm.Value
	var us_61 vm.Value
	var l_62 vm.Value
	var nid_63 vm.Value
	var uc_64 vm.Value
	var uses_65 vm.Value
	var us_66 vm.Value
	var l_67 vm.Value
	var nid_68 vm.Value
	var uc_69 vm.Value
	var uses_70 vm.Value
	var v192 vm.Value
	var us_193 vm.Value
	var l_194 vm.Value
	var nid_195 vm.Value
	var uc_196 vm.Value
	var uses_197 vm.Value
	var us_73 vm.Value
	var l_74 vm.Value
	var nid_75 vm.Value
	var uc_76 vm.Value
	var uses_77 vm.Value
	var or__x_78 vm.Value
	var us_79 vm.Value
	var l_80 vm.Value
	var nid_81 vm.Value
	var uc_82 vm.Value
	var uses_83 vm.Value
	var or__x_84 vm.Value
	var v88 vm.Value
	var v90 vm.Value
	var us_91 vm.Value
	var l_92 vm.Value
	var nid_93 vm.Value
	var uc_94 vm.Value
	var uses_95 vm.Value
	var or__x_96 vm.Value
	var us_100 vm.Value
	var l_101 vm.Value
	var nid_102 vm.Value
	var uc_103 vm.Value
	var uses_104 vm.Value
	var user_id_113 vm.Value
	var arg__12571_115 vm.Value
	var arg__12577_118 vm.Value
	var user_refs_119 vm.Value
	var arg__12581_135 vm.Value
	var arg__12586_138 vm.Value
	var v139 vm.Value
	var us_105 vm.Value
	var l_106 vm.Value
	var nid_107 vm.Value
	var uc_108 vm.Value
	var uses_109 vm.Value
	var v185 vm.Value
	var us_186 vm.Value
	var l_187 vm.Value
	var nid_188 vm.Value
	var uc_189 vm.Value
	var uses_190 vm.Value
	var us_120 vm.Value
	var l_121 vm.Value
	var nid_122 vm.Value
	var uc_123 vm.Value
	var uses_124 vm.Value
	var user_id_125 vm.Value
	var user_refs_126 vm.Value
	var us_127 vm.Value
	var l_128 vm.Value
	var nid_129 vm.Value
	var uc_130 vm.Value
	var uses_131 vm.Value
	var user_id_132 vm.Value
	var user_refs_133 vm.Value
	var v174 vm.Value
	var us_175 vm.Value
	var l_176 vm.Value
	var nid_177 vm.Value
	var uc_178 vm.Value
	var uses_179 vm.Value
	var user_id_180 vm.Value
	var user_refs_181 vm.Value
	var us_143 vm.Value
	var l_144 vm.Value
	var nid_145 vm.Value
	var uc_146 vm.Value
	var uses_147 vm.Value
	var user_id_148 vm.Value
	var user_refs_149 vm.Value
	var arg__12590_160 vm.Value
	var v161 bool
	var us_150 vm.Value
	var l_151 vm.Value
	var nid_152 vm.Value
	var uc_153 vm.Value
	var uses_154 vm.Value
	var user_id_155 vm.Value
	var user_refs_156 vm.Value
	var v165 vm.Value
	var us_166 vm.Value
	var l_167 vm.Value
	var nid_168 vm.Value
	var uc_169 vm.Value
	var uses_170 vm.Value
	var user_id_171 vm.Value
	var user_refs_172 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = uc_4, v12, l_5, nid_6, uc_7, l_8, nid_9, uc_10, v23, v211, l_212, nid_213, uc_214, l_16, nid_17, uc_18, l_19, nid_20, uc_21, v206, l_207, nid_208, uc_209, l_27, nid_28, uc_29, arg__12547_37, uses_38, arg__12552_48, v49, l_30, nid_31, uc_32, v201, l_202, nid_203, uc_204, l_39, nid_40, uc_41, uses_42, v52, l_43, nid_44, uc_45, uses_46, us_56, l_57, nid_58, uc_59, uses_60, or__x_72, us_61, l_62, nid_63, uc_64, uses_65, us_66, l_67, nid_68, uc_69, uses_70, v192, us_193, l_194, nid_195, uc_196, uses_197, us_73, l_74, nid_75, uc_76, uses_77, or__x_78, us_79, l_80, nid_81, uc_82, uses_83, or__x_84, v88, v90, us_91, l_92, nid_93, uc_94, uses_95, or__x_96, us_100, l_101, nid_102, uc_103, uses_104, user_id_113, arg__12571_115, arg__12577_118, user_refs_119, arg__12581_135, arg__12586_138, v139, us_105, l_106, nid_107, uc_108, uses_109, v185, us_186, l_187, nid_188, uc_189, uses_190, us_120, l_121, nid_122, uc_123, uses_124, user_id_125, user_refs_126, us_127, l_128, nid_129, uc_130, uses_131, user_id_132, user_refs_133, v174, us_175, l_176, nid_177, uc_178, uses_179, user_id_180, user_refs_181, us_143, l_144, nid_145, uc_146, uses_147, user_id_148, user_refs_149, arg__12590_160, v161, us_150, l_151, nid_152, uc_153, uses_154, user_id_155, user_refs_156, v165, us_166, l_167, nid_168, uc_169, uses_170, user_id_171, user_refs_172
	uc_4, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "use-count-of").Deref(), []vm.Value{arg0, arg1})
	if callErr != nil {
		return nil, callErr
	}
	v12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "zero?").Deref(), []vm.Value{uc_4})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v12) {
		l_5 = arg0
		nid_6 = arg1
		uc_7 = uc_4
		goto b1
	} else {
		l_8 = arg0
		nid_9 = arg1
		uc_10 = uc_4
		goto b2
	}
b1:
	;
	v211 = vm.Boolean(false)
	l_212 = l_5
	nid_213 = nid_6
	uc_214 = uc_7
	goto b3
b2:
	;
	v23 = rt.GtValue(uc_10, vm.Int(1))
	if v23 {
		l_16 = l_8
		nid_17 = nid_9
		uc_18 = uc_10
		goto b4
	} else {
		l_19 = l_8
		nid_20 = nid_9
		uc_21 = uc_10
		goto b5
	}
b3:
	;
	return v211, nil
b4:
	;
	v206 = vm.Boolean(false)
	l_207 = l_16
	nid_208 = nid_17
	uc_209 = uc_18
	goto b6
b5:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		l_27 = l_19
		nid_28 = nid_20
		uc_29 = uc_21
		goto b7
	} else {
		l_30 = l_19
		nid_31 = nid_20
		uc_32 = uc_21
		goto b8
	}
b6:
	;
	v211 = v206
	l_212 = l_207
	nid_213 = nid_208
	uc_214 = uc_209
	goto b3
b7:
	;
	arg__12547_37, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{l_27})
	if callErr != nil {
		return nil, callErr
	}
	uses_38, callErr = rt.InvokeValue(vm.Keyword("uses"), []vm.Value{arg__12547_37})
	if callErr != nil {
		return nil, callErr
	}
	arg__12552_48, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{uses_38})
	if callErr != nil {
		return nil, callErr
	}
	v49 = rt.LtValue(nid_28, arg__12552_48)
	if v49 {
		l_39 = l_27
		nid_40 = nid_28
		uc_41 = uc_29
		uses_42 = uses_38
		goto b10
	} else {
		l_43 = l_27
		nid_44 = nid_28
		uc_45 = uc_29
		uses_46 = uses_38
		goto b11
	}
b8:
	;
	v201 = vm.NIL
	l_202 = l_30
	nid_203 = nid_31
	uc_204 = uc_32
	goto b9
b9:
	;
	v206 = v201
	l_207 = l_202
	nid_208 = nid_203
	uc_209 = uc_204
	goto b6
b10:
	;
	v52, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{uses_42, nid_40})
	if callErr != nil {
		return nil, callErr
	}
	us_56 = v52
	l_57 = l_39
	nid_58 = nid_40
	uc_59 = uc_41
	uses_60 = uses_42
	goto b12
b11:
	;
	us_56 = vm.NIL
	l_57 = l_43
	nid_58 = nid_44
	uc_59 = uc_45
	uses_60 = uses_46
	goto b12
b12:
	;
	or__x_72, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{us_56})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(or__x_72) {
		us_73 = us_56
		l_74 = l_57
		nid_75 = nid_58
		uc_76 = uc_59
		uses_77 = uses_60
		or__x_78 = or__x_72
		goto b16
	} else {
		us_79 = us_56
		l_80 = l_57
		nid_81 = nid_58
		uc_82 = uc_59
		uses_83 = uses_60
		or__x_84 = or__x_72
		goto b17
	}
b13:
	;
	v192 = vm.Boolean(false)
	us_193 = us_61
	l_194 = l_62
	nid_195 = nid_63
	uc_196 = uc_64
	uses_197 = uses_65
	goto b15
b14:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		us_100 = us_66
		l_101 = l_67
		nid_102 = nid_68
		uc_103 = uc_69
		uses_104 = uses_70
		goto b19
	} else {
		us_105 = us_66
		l_106 = l_67
		nid_107 = nid_68
		uc_108 = uc_69
		uses_109 = uses_70
		goto b20
	}
b15:
	;
	v201 = v192
	l_202 = l_194
	nid_203 = nid_195
	uc_204 = uc_196
	goto b9
b16:
	;
	v90 = or__x_78
	us_91 = us_73
	l_92 = l_74
	nid_93 = nid_75
	uc_94 = uc_76
	uses_95 = uses_77
	or__x_96 = or__x_78
	goto b18
b17:
	;
	v88, callErr = rt.InvokeValue(rt.LookupVar("ir", "uses-empty?").Deref(), []vm.Value{us_79})
	if callErr != nil {
		return nil, callErr
	}
	v90 = v88
	us_91 = us_79
	l_92 = l_80
	nid_93 = nid_81
	uc_94 = uc_82
	uses_95 = uses_83
	or__x_96 = or__x_84
	goto b18
b18:
	;
	if vm.IsTruthy(v90) {
		us_61 = us_91
		l_62 = l_92
		nid_63 = nid_93
		uc_64 = uc_94
		uses_65 = uses_95
		goto b13
	} else {
		us_66 = us_91
		l_67 = l_92
		nid_68 = nid_93
		uc_69 = uc_94
		uses_70 = uses_95
		goto b14
	}
b19:
	;
	user_id_113, callErr = rt.InvokeValue(rt.LookupVar("ir", "uses-first").Deref(), []vm.Value{us_100})
	if callErr != nil {
		return nil, callErr
	}
	arg__12571_115, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "f-of").Deref(), []vm.Value{l_101})
	if callErr != nil {
		return nil, callErr
	}
	arg__12577_118, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "f-of").Deref(), []vm.Value{l_101})
	if callErr != nil {
		return nil, callErr
	}
	user_refs_119, callErr = rt.InvokeValue(rt.LookupVar("ir", "refs").Deref(), []vm.Value{user_id_113, arg__12577_118})
	if callErr != nil {
		return nil, callErr
	}
	arg__12581_135, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{user_refs_119})
	if callErr != nil {
		return nil, callErr
	}
	arg__12586_138, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{user_refs_119})
	if callErr != nil {
		return nil, callErr
	}
	v139, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "zero?").Deref(), []vm.Value{arg__12586_138})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v139) {
		us_120 = us_100
		l_121 = l_101
		nid_122 = nid_102
		uc_123 = uc_103
		uses_124 = uses_104
		user_id_125 = user_id_113
		user_refs_126 = user_refs_119
		goto b22
	} else {
		us_127 = us_100
		l_128 = l_101
		nid_129 = nid_102
		uc_130 = uc_103
		uses_131 = uses_104
		user_id_132 = user_id_113
		user_refs_133 = user_refs_119
		goto b23
	}
b20:
	;
	v185 = vm.NIL
	us_186 = us_105
	l_187 = l_106
	nid_188 = nid_107
	uc_189 = uc_108
	uses_190 = uses_109
	goto b21
b21:
	;
	v192 = v185
	us_193 = us_186
	l_194 = l_187
	nid_195 = nid_188
	uc_196 = uc_189
	uses_197 = uses_190
	goto b15
b22:
	;
	v174 = vm.Boolean(true)
	us_175 = us_120
	l_176 = l_121
	nid_177 = nid_122
	uc_178 = uc_123
	uses_179 = uses_124
	user_id_180 = user_id_125
	user_refs_181 = user_refs_126
	goto b24
b23:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		us_143 = us_127
		l_144 = l_128
		nid_145 = nid_129
		uc_146 = uc_130
		uses_147 = uses_131
		user_id_148 = user_id_132
		user_refs_149 = user_refs_133
		goto b25
	} else {
		us_150 = us_127
		l_151 = l_128
		nid_152 = nid_129
		uc_153 = uc_130
		uses_154 = uses_131
		user_id_155 = user_id_132
		user_refs_156 = user_refs_133
		goto b26
	}
b24:
	;
	v185 = v174
	us_186 = us_175
	l_187 = l_176
	nid_188 = nid_177
	uc_189 = uc_178
	uses_190 = uses_179
	goto b21
b25:
	;
	arg__12590_160, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{user_refs_149})
	if callErr != nil {
		return nil, callErr
	}
	v161 = arg__12590_160 == nid_145
	v165 = vm.Boolean(v161)
	us_166 = us_143
	l_167 = l_144
	nid_168 = nid_145
	uc_169 = uc_146
	uses_170 = uses_147
	user_id_171 = user_id_148
	user_refs_172 = user_refs_149
	goto b27
b26:
	;
	v165 = vm.NIL
	us_166 = us_150
	l_167 = l_151
	nid_168 = nid_152
	uc_169 = uc_153
	uses_170 = uses_154
	user_id_171 = user_id_155
	user_refs_172 = user_refs_156
	goto b27
b27:
	;
	v174 = v165
	us_175 = us_166
	l_176 = l_167
	nid_177 = nid_168
	uc_178 = uc_169
	uses_179 = uses_170
	user_id_180 = user_id_171
	user_refs_181 = user_refs_172
	goto b24
}
func stack_sp(arg0 vm.Value) (vm.Value, error) {
	var arg__12596_3 vm.Value
	var v4 vm.Value
	var callErr error
	_, _ = arg__12596_3, v4
	arg__12596_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v4, callErr = rt.InvokeValue(vm.Keyword("stack-sp"), []vm.Value{arg__12596_3})
	if callErr != nil {
		return nil, callErr
	}
	return v4, nil
}
func template_value(arg0 vm.Value) (vm.Value, error) {
	var and__x_4 vm.Value
	var aux_1 vm.Value
	var inner_21 vm.Value
	var arity_23 vm.Value
	var variadic_QMARK__25 vm.Value
	var chunk_27 vm.Value
	var v29 vm.Value
	var aux_2 vm.Value
	var and__x_34 vm.Value
	var v82 vm.Value
	var aux_83 vm.Value
	var aux_5 vm.Value
	var and__x_6 vm.Value
	var arg__12602_11 vm.Value
	var v13 bool
	var aux_7 vm.Value
	var and__x_8 vm.Value
	var v16 vm.Value
	var aux_17 vm.Value
	var and__x_18 vm.Value
	var aux_31 vm.Value
	var arg__12630_52 vm.Value
	var arg__12635_56 vm.Value
	var fn_vals_57 vm.Value
	var arg__12641_61 vm.Value
	var arg__12648_66 vm.Value
	var v67 vm.Value
	var aux_32 vm.Value
	var v79 vm.Value
	var aux_80 vm.Value
	var aux_35 vm.Value
	var and__x_36 vm.Value
	var arg__12625_41 vm.Value
	var v43 bool
	var aux_37 vm.Value
	var and__x_38 vm.Value
	var v46 vm.Value
	var aux_47 vm.Value
	var and__x_48 vm.Value
	var aux_69 vm.Value
	var aux_70 vm.Value
	var v76 vm.Value
	var aux_77 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = and__x_4, aux_1, inner_21, arity_23, variadic_QMARK__25, chunk_27, v29, aux_2, and__x_34, v82, aux_83, aux_5, and__x_6, arg__12602_11, v13, aux_7, and__x_8, v16, aux_17, and__x_18, aux_31, arg__12630_52, arg__12635_56, fn_vals_57, arg__12641_61, arg__12648_66, v67, aux_32, v79, aux_80, aux_35, and__x_36, arg__12625_41, v43, aux_37, and__x_38, v46, aux_47, and__x_48, aux_69, aux_70, v76, aux_77
	and__x_4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map?").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(and__x_4) {
		aux_5 = arg0
		and__x_6 = and__x_4
		goto b4
	} else {
		aux_7 = arg0
		and__x_8 = and__x_4
		goto b5
	}
b1:
	;
	inner_21, callErr = rt.InvokeValue(vm.Keyword("fn"), []vm.Value{aux_1})
	if callErr != nil {
		return nil, callErr
	}
	arity_23, callErr = rt.InvokeValue(vm.Keyword("arity"), []vm.Value{aux_1})
	if callErr != nil {
		return nil, callErr
	}
	variadic_QMARK__25, callErr = rt.InvokeValue(vm.Keyword("variadic?"), []vm.Value{aux_1})
	if callErr != nil {
		return nil, callErr
	}
	chunk_27, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "lower").Deref(), []vm.Value{inner_21})
	if callErr != nil {
		return nil, callErr
	}
	v29, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "chunk->fn").Deref(), []vm.Value{arity_23, variadic_QMARK__25, chunk_27})
	if callErr != nil {
		return nil, callErr
	}
	v82 = v29
	aux_83 = aux_1
	goto b3
b2:
	;
	and__x_34, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map?").Deref(), []vm.Value{aux_2})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(and__x_34) {
		aux_35 = aux_2
		and__x_36 = and__x_34
		goto b10
	} else {
		aux_37 = aux_2
		and__x_38 = and__x_34
		goto b11
	}
b3:
	;
	return v82, nil
b4:
	;
	arg__12602_11, callErr = rt.InvokeValue(vm.Keyword("kind"), []vm.Value{aux_5})
	if callErr != nil {
		return nil, callErr
	}
	v13 = arg__12602_11 == vm.Keyword("fn-template")
	v16 = vm.Boolean(v13)
	aux_17 = aux_5
	and__x_18 = and__x_6
	goto b6
b5:
	;
	v16 = and__x_8
	aux_17 = aux_7
	and__x_18 = and__x_8
	goto b6
b6:
	;
	if vm.IsTruthy(v16) {
		aux_1 = aux_17
		goto b1
	} else {
		aux_2 = aux_17
		goto b2
	}
b7:
	;
	arg__12630_52, callErr = rt.InvokeValue(vm.Keyword("fns"), []vm.Value{aux_31})
	if callErr != nil {
		return nil, callErr
	}
	arg__12635_56, callErr = rt.InvokeValue(vm.Keyword("fns"), []vm.Value{aux_31})
	if callErr != nil {
		return nil, callErr
	}
	fn_vals_57, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapv").Deref(), []vm.Value{rt.LookupVar("ir.lower", "template-value").Deref(), arg__12635_56})
	if callErr != nil {
		return nil, callErr
	}
	arg__12641_61, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "apply").Deref(), []vm.Value{rt.LookupVar("clojure.core", "list").Deref(), fn_vals_57})
	if callErr != nil {
		return nil, callErr
	}
	arg__12648_66, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "apply").Deref(), []vm.Value{rt.LookupVar("clojure.core", "list").Deref(), fn_vals_57})
	if callErr != nil {
		return nil, callErr
	}
	v67, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "make-multi-arity").Deref(), []vm.Value{arg__12648_66})
	if callErr != nil {
		return nil, callErr
	}
	v79 = v67
	aux_80 = aux_31
	goto b9
b8:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		aux_69 = aux_32
		goto b13
	} else {
		aux_70 = aux_32
		goto b14
	}
b9:
	;
	v82 = v79
	aux_83 = aux_80
	goto b3
b10:
	;
	arg__12625_41, callErr = rt.InvokeValue(vm.Keyword("kind"), []vm.Value{aux_35})
	if callErr != nil {
		return nil, callErr
	}
	v43 = arg__12625_41 == vm.Keyword("multi-fn-template")
	v46 = vm.Boolean(v43)
	aux_47 = aux_35
	and__x_48 = and__x_36
	goto b12
b11:
	;
	v46 = and__x_38
	aux_47 = aux_37
	and__x_48 = and__x_38
	goto b12
b12:
	;
	if vm.IsTruthy(v46) {
		aux_31 = aux_47
		goto b7
	} else {
		aux_32 = aux_47
		goto b8
	}
b13:
	;
	v76 = aux_69
	aux_77 = aux_69
	goto b15
b14:
	;
	v76 = vm.NIL
	aux_77 = aux_70
	goto b15
b15:
	;
	v79 = v76
	aux_80 = aux_77
	goto b9
}
func use_count_of(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__12653_4 vm.Value
	var arg__12654_5 vm.Value
	var arg__12661_9 vm.Value
	var arg__12662_10 vm.Value
	var or__x_11 vm.Value
	var l_12 vm.Value
	var nid_13 vm.Value
	var or__x_14 vm.Value
	var l_15 vm.Value
	var nid_16 vm.Value
	var or__x_17 vm.Value
	var v22 vm.Value
	var l_23 vm.Value
	var nid_24 vm.Value
	var or__x_25 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__12653_4, arg__12654_5, arg__12661_9, arg__12662_10, or__x_11, l_12, nid_13, or__x_14, l_15, nid_16, or__x_17, v22, l_23, nid_24, or__x_25
	arg__12653_4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__12654_5, callErr = rt.InvokeValue(vm.Keyword("use-count"), []vm.Value{arg__12653_4})
	if callErr != nil {
		return nil, callErr
	}
	arg__12661_9, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__12662_10, callErr = rt.InvokeValue(vm.Keyword("use-count"), []vm.Value{arg__12661_9})
	if callErr != nil {
		return nil, callErr
	}
	or__x_11, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg__12662_10, arg1})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(or__x_11) {
		l_12 = arg0
		nid_13 = arg1
		or__x_14 = or__x_11
		goto b1
	} else {
		l_15 = arg0
		nid_16 = arg1
		or__x_17 = or__x_11
		goto b2
	}
b1:
	;
	v22 = or__x_14
	l_23 = l_12
	nid_24 = nid_13
	or__x_25 = or__x_14
	goto b3
b2:
	;
	v22 = vm.Int(0)
	l_23 = l_15
	nid_24 = nid_16
	or__x_25 = or__x_17
	goto b3
b3:
	;
	return v22, nil
}
func value_pos_of(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__12668_4 vm.Value
	var arg__12669_5 vm.Value
	var arg__12676_9 vm.Value
	var arg__12677_10 vm.Value
	var v11 vm.Value
	var callErr error
	_, _, _, _, _ = arg__12668_4, arg__12669_5, arg__12676_9, arg__12677_10, v11
	arg__12668_4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__12669_5, callErr = rt.InvokeValue(vm.Keyword("value-stack-pos"), []vm.Value{arg__12668_4})
	if callErr != nil {
		return nil, callErr
	}
	arg__12676_9, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__12677_10, callErr = rt.InvokeValue(vm.Keyword("value-stack-pos"), []vm.Value{arg__12676_9})
	if callErr != nil {
		return nil, callErr
	}
	v11, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg__12677_10, arg1})
	if callErr != nil {
		return nil, callErr
	}
	return v11, nil
}
