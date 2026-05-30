package ir_passes_pipeline

import (
	rt "github.com/nooga/let-go/pkg/rt"
	vm "github.com/nooga/let-go/pkg/vm"
)

func unwrap_name(arg0 vm.Value) (vm.Value, error) {
	var and__x_5 vm.Value
	var raw_name_2 vm.Value
	var v22 vm.Value
	var raw_name_3 vm.Value
	var v25 vm.Value
	var raw_name_26 vm.Value
	var raw_name_6 vm.Value
	var and__x_7 vm.Value
	var arg__29143_13 vm.Value
	var v14 bool
	var raw_name_8 vm.Value
	var and__x_9 vm.Value
	var v17 vm.Value
	var raw_name_18 vm.Value
	var and__x_19 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = and__x_5, raw_name_2, v22, raw_name_3, v25, raw_name_26, raw_name_6, and__x_7, arg__29143_13, v14, raw_name_8, and__x_9, v17, raw_name_18, and__x_19
	and__x_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq?").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(and__x_5) {
		raw_name_6 = arg0
		and__x_7 = and__x_5
		goto b4
	} else {
		raw_name_8 = arg0
		and__x_9 = and__x_5
		goto b5
	}
b1:
	;
	v22, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "second").Deref(), []vm.Value{raw_name_2})
	if callErr != nil {
		return nil, callErr
	}
	v25 = v22
	raw_name_26 = raw_name_2
	goto b3
b2:
	;
	v25 = raw_name_3
	raw_name_26 = raw_name_3
	goto b3
b3:
	;
	return v25, nil
b4:
	;
	arg__29143_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{raw_name_6})
	if callErr != nil {
		return nil, callErr
	}
	v14 = arg__29143_13 == vm.Symbol("with-meta")
	v17 = vm.Boolean(v14)
	raw_name_18 = raw_name_6
	and__x_19 = and__x_7
	goto b6
b5:
	;
	v17 = and__x_9
	raw_name_18 = raw_name_8
	and__x_19 = and__x_9
	goto b6
b6:
	;
	if vm.IsTruthy(v17) {
		raw_name_2 = raw_name_18
		goto b1
	} else {
		raw_name_3 = raw_name_18
		goto b2
	}
}
func collect_call_targets(arg0 vm.Value) (vm.Value, error) {
	var arg__29155_2 vm.Value
	var arg__29158_5 vm.Value
	var acc_6 vm.Value
	var letfn__29147_10 vm.Value
	var v23 vm.Value
	var v24 vm.Value
	var v26 vm.Value
	var callErr error
	_, _, _, _, _, _, _ = arg__29155_2, arg__29158_5, acc_6, letfn__29147_10, v23, v24, v26
	arg__29155_2, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__29158_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	acc_6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "atom").Deref(), []vm.Value{arg__29158_5})
	if callErr != nil {
		return nil, callErr
	}
	letfn__29147_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "atom").Deref(), []vm.Value{vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	v23, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "reset!").Deref(), []vm.Value{letfn__29147_10, rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v10 vm.Value
		var x_3 vm.Value
		var acc_4 vm.Value
		var visit_5 vm.Value
		var head_13 vm.Value
		var v23 vm.Value
		var x_6 vm.Value
		var acc_7 vm.Value
		var visit_8 vm.Value
		var v86 vm.Value
		var v195 vm.Value
		var x_196 vm.Value
		var acc_197 vm.Value
		var visit_198 vm.Value
		var x_14 vm.Value
		var acc_15 vm.Value
		var visit_16 vm.Value
		var head_17 vm.Value
		var arg__29276_27 vm.Value
		var arg__29283_31 vm.Value
		var v32 vm.Value
		var x_18 vm.Value
		var acc_19 vm.Value
		var visit_20 vm.Value
		var head_21 vm.Value
		var v36 vm.Value
		var x_37 vm.Value
		var acc_38 vm.Value
		var visit_39 vm.Value
		var head_40 vm.Value
		var arg__29287_42 vm.Value
		var arg__29292_45 vm.Value
		var doseq_seq__29148_46 vm.Value
		var doseq_loop__29149_47 vm.Value
		var visit_48 vm.Value
		var x_50 vm.Value
		var acc_51 vm.Value
		var head_52 vm.Value
		var doseq_seq__29148_53 vm.Value
		var doseq_loop__29149_54 vm.Value
		var visit_55 vm.Value
		var el_64 vm.Value
		var v65 vm.Value
		var v67 vm.Value
		var x_56 vm.Value
		var acc_57 vm.Value
		var head_58 vm.Value
		var doseq_seq__29148_59 vm.Value
		var doseq_loop__29149_60 vm.Value
		var visit_61 vm.Value
		var v71 vm.Value
		var x_72 vm.Value
		var acc_73 vm.Value
		var head_74 vm.Value
		var doseq_seq__29148_75 vm.Value
		var doseq_loop__29149_76 vm.Value
		var visit_77 vm.Value
		var x_79 vm.Value
		var acc_80 vm.Value
		var visit_81 vm.Value
		var doseq_seq__29150_89 vm.Value
		var x_82 vm.Value
		var acc_83 vm.Value
		var visit_84 vm.Value
		var v126 vm.Value
		var v190 vm.Value
		var x_191 vm.Value
		var acc_192 vm.Value
		var visit_193 vm.Value
		var doseq_loop__29151_90 vm.Value
		var visit_91 vm.Value
		var x_93 vm.Value
		var acc_94 vm.Value
		var doseq_seq__29150_95 vm.Value
		var doseq_loop__29151_96 vm.Value
		var visit_97 vm.Value
		var el_105 vm.Value
		var v106 vm.Value
		var v108 vm.Value
		var x_98 vm.Value
		var acc_99 vm.Value
		var doseq_seq__29150_100 vm.Value
		var doseq_loop__29151_101 vm.Value
		var visit_102 vm.Value
		var v112 vm.Value
		var x_113 vm.Value
		var acc_114 vm.Value
		var doseq_seq__29150_115 vm.Value
		var doseq_loop__29151_116 vm.Value
		var visit_117 vm.Value
		var x_119 vm.Value
		var acc_120 vm.Value
		var visit_121 vm.Value
		var doseq_seq__29152_129 vm.Value
		var x_122 vm.Value
		var acc_123 vm.Value
		var visit_124 vm.Value
		var v185 vm.Value
		var x_186 vm.Value
		var acc_187 vm.Value
		var visit_188 vm.Value
		var doseq_loop__29153_130 vm.Value
		var visit_131 vm.Value
		var x_133 vm.Value
		var acc_134 vm.Value
		var doseq_seq__29152_135 vm.Value
		var doseq_loop__29153_136 vm.Value
		var visit_137 vm.Value
		var entry_145 vm.Value
		var arg__29329_147 vm.Value
		var arg__29334_149 vm.Value
		var v150 vm.Value
		var arg__29338_152 vm.Value
		var arg__29343_154 vm.Value
		var v155 vm.Value
		var v157 vm.Value
		var x_138 vm.Value
		var acc_139 vm.Value
		var doseq_seq__29152_140 vm.Value
		var doseq_loop__29153_141 vm.Value
		var visit_142 vm.Value
		var v161 vm.Value
		var x_162 vm.Value
		var acc_163 vm.Value
		var doseq_seq__29152_164 vm.Value
		var doseq_loop__29153_165 vm.Value
		var visit_166 vm.Value
		var x_168 vm.Value
		var acc_169 vm.Value
		var visit_170 vm.Value
		var x_171 vm.Value
		var acc_172 vm.Value
		var visit_173 vm.Value
		var v180 vm.Value
		var x_181 vm.Value
		var acc_182 vm.Value
		var visit_183 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = v10, x_3, acc_4, visit_5, head_13, v23, x_6, acc_7, visit_8, v86, v195, x_196, acc_197, visit_198, x_14, acc_15, visit_16, head_17, arg__29276_27, arg__29283_31, v32, x_18, acc_19, visit_20, head_21, v36, x_37, acc_38, visit_39, head_40, arg__29287_42, arg__29292_45, doseq_seq__29148_46, doseq_loop__29149_47, visit_48, x_50, acc_51, head_52, doseq_seq__29148_53, doseq_loop__29149_54, visit_55, el_64, v65, v67, x_56, acc_57, head_58, doseq_seq__29148_59, doseq_loop__29149_60, visit_61, v71, x_72, acc_73, head_74, doseq_seq__29148_75, doseq_loop__29149_76, visit_77, x_79, acc_80, visit_81, doseq_seq__29150_89, x_82, acc_83, visit_84, v126, v190, x_191, acc_192, visit_193, doseq_loop__29151_90, visit_91, x_93, acc_94, doseq_seq__29150_95, doseq_loop__29151_96, visit_97, el_105, v106, v108, x_98, acc_99, doseq_seq__29150_100, doseq_loop__29151_101, visit_102, v112, x_113, acc_114, doseq_seq__29150_115, doseq_loop__29151_116, visit_117, x_119, acc_120, visit_121, doseq_seq__29152_129, x_122, acc_123, visit_124, v185, x_186, acc_187, visit_188, doseq_loop__29153_130, visit_131, x_133, acc_134, doseq_seq__29152_135, doseq_loop__29153_136, visit_137, entry_145, arg__29329_147, arg__29334_149, v150, arg__29338_152, arg__29343_154, v155, v157, x_138, acc_139, doseq_seq__29152_140, doseq_loop__29153_141, visit_142, v161, x_162, acc_163, doseq_seq__29152_164, doseq_loop__29153_165, visit_166, x_168, acc_169, visit_170, x_171, acc_172, visit_173, v180, x_181, acc_182, visit_183
		v10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq?").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		if vm.IsTruthy(v10) {
			x_3 = arg0
			acc_4 = acc_6
			visit_5 = rt.BoxNativeFn(func(args ...vm.Value) (vm.Value, error) {
				var arg__29165_3 vm.Value
				var arg__29171_6 vm.Value
				var v7 vm.Value
				var callErr error
				_, _, _ = arg__29165_3, arg__29171_6, v7
				arg__29165_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{letfn__29147_10})
				if callErr != nil {
					return nil, callErr
				}
				arg__29171_6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{letfn__29147_10})
				if callErr != nil {
					return nil, callErr
				}
				v7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "apply").Deref(), []vm.Value{arg__29171_6, rt.BoxRestArgs(args)})
				if callErr != nil {
					return nil, callErr
				}
				return v7, nil
			})
			goto b1
		} else {
			x_6 = arg0
			acc_7 = acc_6
			visit_8 = rt.BoxNativeFn(func(args ...vm.Value) (vm.Value, error) {
				var arg__29165_3 vm.Value
				var arg__29171_6 vm.Value
				var v7 vm.Value
				var callErr error
				_, _, _ = arg__29165_3, arg__29171_6, v7
				arg__29165_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{letfn__29147_10})
				if callErr != nil {
					return nil, callErr
				}
				arg__29171_6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{letfn__29147_10})
				if callErr != nil {
					return nil, callErr
				}
				v7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "apply").Deref(), []vm.Value{arg__29171_6, rt.BoxRestArgs(args)})
				if callErr != nil {
					return nil, callErr
				}
				return v7, nil
			})
			goto b2
		}
	b1:
		;
		head_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{x_3})
		if callErr != nil {
			return nil, callErr
		}
		v23, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "symbol?").Deref(), []vm.Value{head_13})
		if callErr != nil {
			return nil, callErr
		}
		if vm.IsTruthy(v23) {
			x_14 = x_3
			acc_15 = acc_4
			visit_16 = visit_5
			head_17 = head_13
			goto b4
		} else {
			x_18 = x_3
			acc_19 = acc_4
			visit_20 = visit_5
			head_21 = head_13
			goto b5
		}
	b2:
		;
		v86, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector?").Deref(), []vm.Value{x_6})
		if callErr != nil {
			return nil, callErr
		}
		if vm.IsTruthy(v86) {
			x_79 = x_6
			acc_80 = acc_7
			visit_81 = visit_8
			goto b11
		} else {
			x_82 = x_6
			acc_83 = acc_7
			visit_84 = visit_8
			goto b12
		}
	b3:
		;
		return v195, nil
	b4:
		;
		arg__29276_27, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{head_17})
		if callErr != nil {
			return nil, callErr
		}
		arg__29283_31, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{head_17})
		if callErr != nil {
			return nil, callErr
		}
		v32, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{acc_15, rt.LookupVar("clojure.core", "conj").Deref(), arg__29283_31})
		if callErr != nil {
			return nil, callErr
		}
		v36 = v32
		x_37 = x_14
		acc_38 = acc_15
		visit_39 = visit_16
		head_40 = head_17
		goto b6
	b5:
		;
		v36 = vm.NIL
		x_37 = x_18
		acc_38 = acc_19
		visit_39 = visit_20
		head_40 = head_21
		goto b6
	b6:
		;
		arg__29287_42, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{x_37})
		if callErr != nil {
			return nil, callErr
		}
		arg__29292_45, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{x_37})
		if callErr != nil {
			return nil, callErr
		}
		doseq_seq__29148_46, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{arg__29292_45})
		if callErr != nil {
			return nil, callErr
		}
		doseq_loop__29149_47 = doseq_seq__29148_46
		visit_48 = visit_39
		goto b7
	b7:
		;
		if vm.IsTruthy(doseq_loop__29149_47) {
			x_50 = x_37
			acc_51 = acc_38
			head_52 = head_40
			doseq_seq__29148_53 = doseq_seq__29148_46
			doseq_loop__29149_54 = doseq_loop__29149_47
			visit_55 = visit_48
			goto b8
		} else {
			x_56 = x_37
			acc_57 = acc_38
			head_58 = head_40
			doseq_seq__29148_59 = doseq_seq__29148_46
			doseq_loop__29149_60 = doseq_loop__29149_47
			visit_61 = visit_48
			goto b9
		}
	b8:
		;
		el_64, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{doseq_loop__29149_54})
		if callErr != nil {
			return nil, callErr
		}
		v65, callErr = rt.InvokeValue(visit_55, []vm.Value{el_64})
		if callErr != nil {
			return nil, callErr
		}
		v67, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "next").Deref(), []vm.Value{doseq_loop__29149_54})
		if callErr != nil {
			return nil, callErr
		}
		doseq_loop__29149_47 = v67
		visit_48 = visit_55
		goto b7
	b9:
		;
		v71 = vm.NIL
		x_72 = x_56
		acc_73 = acc_57
		head_74 = head_58
		doseq_seq__29148_75 = doseq_seq__29148_59
		doseq_loop__29149_76 = doseq_loop__29149_60
		visit_77 = visit_61
		goto b10
	b10:
		;
		v195 = v71
		x_196 = x_72
		acc_197 = acc_73
		visit_198 = visit_77
		goto b3
	b11:
		;
		doseq_seq__29150_89, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{x_79})
		if callErr != nil {
			return nil, callErr
		}
		doseq_loop__29151_90 = doseq_seq__29150_89
		visit_91 = visit_81
		goto b14
	b12:
		;
		v126, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map?").Deref(), []vm.Value{x_82})
		if callErr != nil {
			return nil, callErr
		}
		if vm.IsTruthy(v126) {
			x_119 = x_82
			acc_120 = acc_83
			visit_121 = visit_84
			goto b18
		} else {
			x_122 = x_82
			acc_123 = acc_83
			visit_124 = visit_84
			goto b19
		}
	b13:
		;
		v195 = v190
		x_196 = x_191
		acc_197 = acc_192
		visit_198 = visit_193
		goto b3
	b14:
		;
		if vm.IsTruthy(doseq_loop__29151_90) {
			x_93 = x_79
			acc_94 = acc_80
			doseq_seq__29150_95 = doseq_seq__29150_89
			doseq_loop__29151_96 = doseq_loop__29151_90
			visit_97 = visit_91
			goto b15
		} else {
			x_98 = x_79
			acc_99 = acc_80
			doseq_seq__29150_100 = doseq_seq__29150_89
			doseq_loop__29151_101 = doseq_loop__29151_90
			visit_102 = visit_91
			goto b16
		}
	b15:
		;
		el_105, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{doseq_loop__29151_96})
		if callErr != nil {
			return nil, callErr
		}
		v106, callErr = rt.InvokeValue(visit_97, []vm.Value{el_105})
		if callErr != nil {
			return nil, callErr
		}
		v108, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "next").Deref(), []vm.Value{doseq_loop__29151_96})
		if callErr != nil {
			return nil, callErr
		}
		doseq_loop__29151_90 = v108
		visit_91 = visit_97
		goto b14
	b16:
		;
		v112 = vm.NIL
		x_113 = x_98
		acc_114 = acc_99
		doseq_seq__29150_115 = doseq_seq__29150_100
		doseq_loop__29151_116 = doseq_loop__29151_101
		visit_117 = visit_102
		goto b17
	b17:
		;
		v190 = v112
		x_191 = x_113
		acc_192 = acc_114
		visit_193 = visit_117
		goto b13
	b18:
		;
		doseq_seq__29152_129, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{x_119})
		if callErr != nil {
			return nil, callErr
		}
		doseq_loop__29153_130 = doseq_seq__29152_129
		visit_131 = visit_121
		goto b21
	b19:
		;
		if vm.IsTruthy(vm.Keyword("else")) {
			x_168 = x_122
			acc_169 = acc_123
			visit_170 = visit_124
			goto b25
		} else {
			x_171 = x_122
			acc_172 = acc_123
			visit_173 = visit_124
			goto b26
		}
	b20:
		;
		v190 = v185
		x_191 = x_186
		acc_192 = acc_187
		visit_193 = visit_188
		goto b13
	b21:
		;
		if vm.IsTruthy(doseq_loop__29153_130) {
			x_133 = x_119
			acc_134 = acc_120
			doseq_seq__29152_135 = doseq_seq__29152_129
			doseq_loop__29153_136 = doseq_loop__29153_130
			visit_137 = visit_131
			goto b22
		} else {
			x_138 = x_119
			acc_139 = acc_120
			doseq_seq__29152_140 = doseq_seq__29152_129
			doseq_loop__29153_141 = doseq_loop__29153_130
			visit_142 = visit_131
			goto b23
		}
	b22:
		;
		entry_145, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{doseq_loop__29153_136})
		if callErr != nil {
			return nil, callErr
		}
		arg__29329_147, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "key").Deref(), []vm.Value{entry_145})
		if callErr != nil {
			return nil, callErr
		}
		arg__29334_149, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "key").Deref(), []vm.Value{entry_145})
		if callErr != nil {
			return nil, callErr
		}
		v150, callErr = rt.InvokeValue(visit_137, []vm.Value{arg__29334_149})
		if callErr != nil {
			return nil, callErr
		}
		arg__29338_152, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "val").Deref(), []vm.Value{entry_145})
		if callErr != nil {
			return nil, callErr
		}
		arg__29343_154, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "val").Deref(), []vm.Value{entry_145})
		if callErr != nil {
			return nil, callErr
		}
		v155, callErr = rt.InvokeValue(visit_137, []vm.Value{arg__29343_154})
		if callErr != nil {
			return nil, callErr
		}
		v157, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "next").Deref(), []vm.Value{doseq_loop__29153_136})
		if callErr != nil {
			return nil, callErr
		}
		doseq_loop__29153_130 = v157
		visit_131 = visit_137
		goto b21
	b23:
		;
		v161 = vm.NIL
		x_162 = x_138
		acc_163 = acc_139
		doseq_seq__29152_164 = doseq_seq__29152_140
		doseq_loop__29153_165 = doseq_loop__29153_141
		visit_166 = visit_142
		goto b24
	b24:
		;
		v185 = v161
		x_186 = x_162
		acc_187 = acc_163
		visit_188 = visit_166
		goto b20
	b25:
		;
		v180 = vm.NIL
		x_181 = x_168
		acc_182 = acc_169
		visit_183 = visit_170
		goto b27
	b26:
		;
		v180 = vm.NIL
		x_181 = x_171
		acc_182 = acc_172
		visit_183 = visit_173
		goto b27
	b27:
		;
		v185 = v180
		x_186 = x_181
		acc_187 = acc_182
		visit_188 = visit_183
		goto b20
	})})
	if callErr != nil {
		return nil, callErr
	}
	v24, callErr = rt.InvokeValue(rt.BoxNativeFn(func(args ...vm.Value) (vm.Value, error) {
		var arg__29165_3 vm.Value
		var arg__29171_6 vm.Value
		var v7 vm.Value
		var callErr error
		_, _, _ = arg__29165_3, arg__29171_6, v7
		arg__29165_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{letfn__29147_10})
		if callErr != nil {
			return nil, callErr
		}
		arg__29171_6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{letfn__29147_10})
		if callErr != nil {
			return nil, callErr
		}
		v7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "apply").Deref(), []vm.Value{arg__29171_6, rt.BoxRestArgs(args)})
		if callErr != nil {
			return nil, callErr
		}
		return v7, nil
	}), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v26, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{acc_6})
	if callErr != nil {
		return nil, callErr
	}
	return v26, nil
}
func defn_key(arg0 vm.Value) (vm.Value, error) {
	var and__x_4 vm.Value
	var form_1 vm.Value
	var raw_name_23 vm.Value
	var name_sym_25 vm.Value
	var name_str_27 vm.Value
	var maybe_doc_31 vm.Value
	var has_doc_QMARK__33 vm.Value
	var form_2 vm.Value
	var v179 vm.Value
	var form_180 vm.Value
	var form_5 vm.Value
	var and__x_6 vm.Value
	var arg__29360_11 vm.Value
	var v13 bool
	var form_7 vm.Value
	var and__x_8 vm.Value
	var v16 vm.Value
	var form_17 vm.Value
	var and__x_18 vm.Value
	var form_34 vm.Value
	var raw_name_35 vm.Value
	var name_sym_36 vm.Value
	var name_str_37 vm.Value
	var maybe_doc_38 vm.Value
	var has_doc_QMARK__39 vm.Value
	var arg__29384_60 vm.Value
	var v62 bool
	var form_40 vm.Value
	var raw_name_41 vm.Value
	var name_sym_42 vm.Value
	var name_str_43 vm.Value
	var maybe_doc_44 vm.Value
	var has_doc_QMARK__45 vm.Value
	var args_or_arity_80 vm.Value
	var form_81 vm.Value
	var raw_name_82 vm.Value
	var name_sym_83 vm.Value
	var name_str_84 vm.Value
	var maybe_doc_85 vm.Value
	var has_doc_QMARK__86 vm.Value
	var v102 vm.Value
	var form_47 vm.Value
	var raw_name_48 vm.Value
	var name_sym_49 vm.Value
	var name_str_50 vm.Value
	var maybe_doc_51 vm.Value
	var has_doc_QMARK__52 vm.Value
	var v67 vm.Value
	var form_53 vm.Value
	var raw_name_54 vm.Value
	var name_sym_55 vm.Value
	var name_str_56 vm.Value
	var maybe_doc_57 vm.Value
	var has_doc_QMARK__58 vm.Value
	var v71 vm.Value
	var form_72 vm.Value
	var raw_name_73 vm.Value
	var name_sym_74 vm.Value
	var name_str_75 vm.Value
	var maybe_doc_76 vm.Value
	var has_doc_QMARK__77 vm.Value
	var args_or_arity_87 vm.Value
	var form_88 vm.Value
	var raw_name_89 vm.Value
	var name_sym_90 vm.Value
	var name_str_91 vm.Value
	var maybe_doc_92 vm.Value
	var has_doc_QMARK__93 vm.Value
	var arg__29399_106 vm.Value
	var v107 vm.Value
	var args_or_arity_94 vm.Value
	var form_95 vm.Value
	var raw_name_96 vm.Value
	var name_sym_97 vm.Value
	var name_str_98 vm.Value
	var maybe_doc_99 vm.Value
	var has_doc_QMARK__100 vm.Value
	var v124 vm.Value
	var v168 vm.Value
	var args_or_arity_169 vm.Value
	var form_170 vm.Value
	var raw_name_171 vm.Value
	var name_sym_172 vm.Value
	var name_str_173 vm.Value
	var maybe_doc_174 vm.Value
	var has_doc_QMARK__175 vm.Value
	var args_or_arity_109 vm.Value
	var form_110 vm.Value
	var raw_name_111 vm.Value
	var name_sym_112 vm.Value
	var name_str_113 vm.Value
	var maybe_doc_114 vm.Value
	var has_doc_QMARK__115 vm.Value
	var v128 vm.Value
	var args_or_arity_116 vm.Value
	var form_117 vm.Value
	var raw_name_118 vm.Value
	var name_sym_119 vm.Value
	var name_str_120 vm.Value
	var maybe_doc_121 vm.Value
	var has_doc_QMARK__122 vm.Value
	var v159 vm.Value
	var args_or_arity_160 vm.Value
	var form_161 vm.Value
	var raw_name_162 vm.Value
	var name_sym_163 vm.Value
	var name_str_164 vm.Value
	var maybe_doc_165 vm.Value
	var has_doc_QMARK__166 vm.Value
	var args_or_arity_130 vm.Value
	var form_131 vm.Value
	var raw_name_132 vm.Value
	var name_sym_133 vm.Value
	var name_str_134 vm.Value
	var maybe_doc_135 vm.Value
	var has_doc_QMARK__136 vm.Value
	var args_or_arity_137 vm.Value
	var form_138 vm.Value
	var raw_name_139 vm.Value
	var name_sym_140 vm.Value
	var name_str_141 vm.Value
	var maybe_doc_142 vm.Value
	var has_doc_QMARK__143 vm.Value
	var v150 vm.Value
	var args_or_arity_151 vm.Value
	var form_152 vm.Value
	var raw_name_153 vm.Value
	var name_sym_154 vm.Value
	var name_str_155 vm.Value
	var maybe_doc_156 vm.Value
	var has_doc_QMARK__157 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = and__x_4, form_1, raw_name_23, name_sym_25, name_str_27, maybe_doc_31, has_doc_QMARK__33, form_2, v179, form_180, form_5, and__x_6, arg__29360_11, v13, form_7, and__x_8, v16, form_17, and__x_18, form_34, raw_name_35, name_sym_36, name_str_37, maybe_doc_38, has_doc_QMARK__39, arg__29384_60, v62, form_40, raw_name_41, name_sym_42, name_str_43, maybe_doc_44, has_doc_QMARK__45, args_or_arity_80, form_81, raw_name_82, name_sym_83, name_str_84, maybe_doc_85, has_doc_QMARK__86, v102, form_47, raw_name_48, name_sym_49, name_str_50, maybe_doc_51, has_doc_QMARK__52, v67, form_53, raw_name_54, name_sym_55, name_str_56, maybe_doc_57, has_doc_QMARK__58, v71, form_72, raw_name_73, name_sym_74, name_str_75, maybe_doc_76, has_doc_QMARK__77, args_or_arity_87, form_88, raw_name_89, name_sym_90, name_str_91, maybe_doc_92, has_doc_QMARK__93, arg__29399_106, v107, args_or_arity_94, form_95, raw_name_96, name_sym_97, name_str_98, maybe_doc_99, has_doc_QMARK__100, v124, v168, args_or_arity_169, form_170, raw_name_171, name_sym_172, name_str_173, maybe_doc_174, has_doc_QMARK__175, args_or_arity_109, form_110, raw_name_111, name_sym_112, name_str_113, maybe_doc_114, has_doc_QMARK__115, v128, args_or_arity_116, form_117, raw_name_118, name_sym_119, name_str_120, maybe_doc_121, has_doc_QMARK__122, v159, args_or_arity_160, form_161, raw_name_162, name_sym_163, name_str_164, maybe_doc_165, has_doc_QMARK__166, args_or_arity_130, form_131, raw_name_132, name_sym_133, name_str_134, maybe_doc_135, has_doc_QMARK__136, args_or_arity_137, form_138, raw_name_139, name_sym_140, name_str_141, maybe_doc_142, has_doc_QMARK__143, v150, args_or_arity_151, form_152, raw_name_153, name_sym_154, name_str_155, maybe_doc_156, has_doc_QMARK__157
	and__x_4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq?").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(and__x_4) {
		form_5 = arg0
		and__x_6 = and__x_4
		goto b4
	} else {
		form_7 = arg0
		and__x_8 = and__x_4
		goto b5
	}
b1:
	;
	raw_name_23, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{form_1, vm.Int(1)})
	if callErr != nil {
		return nil, callErr
	}
	name_sym_25, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.pipeline", "unwrap-name").Deref(), []vm.Value{raw_name_23})
	if callErr != nil {
		return nil, callErr
	}
	name_str_27, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{name_sym_25})
	if callErr != nil {
		return nil, callErr
	}
	maybe_doc_31, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{form_1, vm.Int(2)})
	if callErr != nil {
		return nil, callErr
	}
	has_doc_QMARK__33, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "string?").Deref(), []vm.Value{maybe_doc_31})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(has_doc_QMARK__33) {
		form_34 = form_1
		raw_name_35 = raw_name_23
		name_sym_36 = name_sym_25
		name_str_37 = name_str_27
		maybe_doc_38 = maybe_doc_31
		has_doc_QMARK__39 = has_doc_QMARK__33
		goto b7
	} else {
		form_40 = form_1
		raw_name_41 = raw_name_23
		name_sym_42 = name_sym_25
		name_str_43 = name_str_27
		maybe_doc_44 = maybe_doc_31
		has_doc_QMARK__45 = has_doc_QMARK__33
		goto b8
	}
b2:
	;
	v179 = vm.NIL
	form_180 = form_2
	goto b3
b3:
	;
	return v179, nil
b4:
	;
	arg__29360_11, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{form_5})
	if callErr != nil {
		return nil, callErr
	}
	v13 = rt.GeValue(arg__29360_11, vm.Int(3))
	v16 = vm.Boolean(v13)
	form_17 = form_5
	and__x_18 = and__x_6
	goto b6
b5:
	;
	v16 = and__x_8
	form_17 = form_7
	and__x_18 = and__x_8
	goto b6
b6:
	;
	if vm.IsTruthy(v16) {
		form_1 = form_17
		goto b1
	} else {
		form_2 = form_17
		goto b2
	}
b7:
	;
	arg__29384_60, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{form_34})
	if callErr != nil {
		return nil, callErr
	}
	v62 = rt.GeValue(arg__29384_60, vm.Int(4))
	if v62 {
		form_47 = form_34
		raw_name_48 = raw_name_35
		name_sym_49 = name_sym_36
		name_str_50 = name_str_37
		maybe_doc_51 = maybe_doc_38
		has_doc_QMARK__52 = has_doc_QMARK__39
		goto b10
	} else {
		form_53 = form_34
		raw_name_54 = raw_name_35
		name_sym_55 = name_sym_36
		name_str_56 = name_str_37
		maybe_doc_57 = maybe_doc_38
		has_doc_QMARK__58 = has_doc_QMARK__39
		goto b11
	}
b8:
	;
	args_or_arity_80 = maybe_doc_44
	form_81 = form_40
	raw_name_82 = raw_name_41
	name_sym_83 = name_sym_42
	name_str_84 = name_str_43
	maybe_doc_85 = maybe_doc_44
	has_doc_QMARK__86 = has_doc_QMARK__45
	goto b9
b9:
	;
	v102, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector?").Deref(), []vm.Value{args_or_arity_80})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v102) {
		args_or_arity_87 = args_or_arity_80
		form_88 = form_81
		raw_name_89 = raw_name_82
		name_sym_90 = name_sym_83
		name_str_91 = name_str_84
		maybe_doc_92 = maybe_doc_85
		has_doc_QMARK__93 = has_doc_QMARK__86
		goto b13
	} else {
		args_or_arity_94 = args_or_arity_80
		form_95 = form_81
		raw_name_96 = raw_name_82
		name_sym_97 = name_sym_83
		name_str_98 = name_str_84
		maybe_doc_99 = maybe_doc_85
		has_doc_QMARK__100 = has_doc_QMARK__86
		goto b14
	}
b10:
	;
	v67, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{form_47, vm.Int(3)})
	if callErr != nil {
		return nil, callErr
	}
	v71 = v67
	form_72 = form_47
	raw_name_73 = raw_name_48
	name_sym_74 = name_sym_49
	name_str_75 = name_str_50
	maybe_doc_76 = maybe_doc_51
	has_doc_QMARK__77 = has_doc_QMARK__52
	goto b12
b11:
	;
	v71 = vm.NIL
	form_72 = form_53
	raw_name_73 = raw_name_54
	name_sym_74 = name_sym_55
	name_str_75 = name_str_56
	maybe_doc_76 = maybe_doc_57
	has_doc_QMARK__77 = has_doc_QMARK__58
	goto b12
b12:
	;
	args_or_arity_80 = v71
	form_81 = form_72
	raw_name_82 = raw_name_73
	name_sym_83 = name_sym_74
	name_str_84 = name_str_75
	maybe_doc_85 = maybe_doc_76
	has_doc_QMARK__86 = has_doc_QMARK__77
	goto b9
b13:
	;
	arg__29399_106, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{args_or_arity_87})
	if callErr != nil {
		return nil, callErr
	}
	v107, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{name_str_91, arg__29399_106})
	if callErr != nil {
		return nil, callErr
	}
	v168 = v107
	args_or_arity_169 = args_or_arity_87
	form_170 = form_88
	raw_name_171 = raw_name_89
	name_sym_172 = name_sym_90
	name_str_173 = name_str_91
	maybe_doc_174 = maybe_doc_92
	has_doc_QMARK__175 = has_doc_QMARK__93
	goto b15
b14:
	;
	v124, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "list?").Deref(), []vm.Value{args_or_arity_94})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v124) {
		args_or_arity_109 = args_or_arity_94
		form_110 = form_95
		raw_name_111 = raw_name_96
		name_sym_112 = name_sym_97
		name_str_113 = name_str_98
		maybe_doc_114 = maybe_doc_99
		has_doc_QMARK__115 = has_doc_QMARK__100
		goto b16
	} else {
		args_or_arity_116 = args_or_arity_94
		form_117 = form_95
		raw_name_118 = raw_name_96
		name_sym_119 = name_sym_97
		name_str_120 = name_str_98
		maybe_doc_121 = maybe_doc_99
		has_doc_QMARK__122 = has_doc_QMARK__100
		goto b17
	}
b15:
	;
	v179 = v168
	form_180 = form_170
	goto b3
b16:
	;
	v128, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{name_str_113, vm.Keyword("multi")})
	if callErr != nil {
		return nil, callErr
	}
	v159 = v128
	args_or_arity_160 = args_or_arity_109
	form_161 = form_110
	raw_name_162 = raw_name_111
	name_sym_163 = name_sym_112
	name_str_164 = name_str_113
	maybe_doc_165 = maybe_doc_114
	has_doc_QMARK__166 = has_doc_QMARK__115
	goto b18
b17:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		args_or_arity_130 = args_or_arity_116
		form_131 = form_117
		raw_name_132 = raw_name_118
		name_sym_133 = name_sym_119
		name_str_134 = name_str_120
		maybe_doc_135 = maybe_doc_121
		has_doc_QMARK__136 = has_doc_QMARK__122
		goto b19
	} else {
		args_or_arity_137 = args_or_arity_116
		form_138 = form_117
		raw_name_139 = raw_name_118
		name_sym_140 = name_sym_119
		name_str_141 = name_str_120
		maybe_doc_142 = maybe_doc_121
		has_doc_QMARK__143 = has_doc_QMARK__122
		goto b20
	}
b18:
	;
	v168 = v159
	args_or_arity_169 = args_or_arity_160
	form_170 = form_161
	raw_name_171 = raw_name_162
	name_sym_172 = name_sym_163
	name_str_173 = name_str_164
	maybe_doc_174 = maybe_doc_165
	has_doc_QMARK__175 = has_doc_QMARK__166
	goto b15
b19:
	;
	v150 = vm.NIL
	args_or_arity_151 = args_or_arity_130
	form_152 = form_131
	raw_name_153 = raw_name_132
	name_sym_154 = name_sym_133
	name_str_155 = name_str_134
	maybe_doc_156 = maybe_doc_135
	has_doc_QMARK__157 = has_doc_QMARK__136
	goto b21
b20:
	;
	v150 = vm.NIL
	args_or_arity_151 = args_or_arity_137
	form_152 = form_138
	raw_name_153 = raw_name_139
	name_sym_154 = name_sym_140
	name_str_155 = name_str_141
	maybe_doc_156 = maybe_doc_142
	has_doc_QMARK__157 = has_doc_QMARK__143
	goto b21
b21:
	;
	v159 = v150
	args_or_arity_160 = args_or_arity_151
	form_161 = form_152
	raw_name_162 = raw_name_153
	name_sym_163 = name_sym_154
	name_str_164 = name_str_155
	maybe_doc_165 = maybe_doc_156
	has_doc_QMARK__166 = has_doc_QMARK__157
	goto b18
}
func forms_by_arity(arg0 vm.Value) (vm.Value, error) {
	var remaining_1 vm.Value
	var acc_2 vm.Value
	var v12 vm.Value
	var forms_5 vm.Value
	var remaining_6 vm.Value
	var acc_7 vm.Value
	var forms_8 vm.Value
	var remaining_9 vm.Value
	var acc_10 vm.Value
	var form_16 vm.Value
	var k_18 vm.Value
	var v20 vm.Value
	var v43 vm.Value
	var forms_44 vm.Value
	var remaining_45 vm.Value
	var acc_46 vm.Value
	var forms_21 vm.Value
	var remaining_22 vm.Value
	var acc_23 vm.Value
	var form_24 vm.Value
	var k_25 vm.Value
	var v33 vm.Value
	var forms_26 vm.Value
	var remaining_27 vm.Value
	var acc_28 vm.Value
	var form_29 vm.Value
	var k_30 vm.Value
	var v36 vm.Value
	var forms_37 vm.Value
	var remaining_38 vm.Value
	var acc_39 vm.Value
	var form_40 vm.Value
	var k_41 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = remaining_1, acc_2, v12, forms_5, remaining_6, acc_7, forms_8, remaining_9, acc_10, form_16, k_18, v20, v43, forms_44, remaining_45, acc_46, forms_21, remaining_22, acc_23, form_24, k_25, v33, forms_26, remaining_27, acc_28, form_29, k_30, v36, forms_37, remaining_38, acc_39, form_40, k_41
	remaining_1 = arg0
	acc_2 = vm.EmptyPersistentMap
	goto b1
b1:
	;
	v12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "empty?").Deref(), []vm.Value{remaining_1})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v12) {
		forms_5 = arg0
		remaining_6 = remaining_1
		acc_7 = acc_2
		goto b2
	} else {
		forms_8 = arg0
		remaining_9 = remaining_1
		acc_10 = acc_2
		goto b3
	}
b2:
	;
	v43 = acc_7
	forms_44 = forms_5
	remaining_45 = remaining_6
	acc_46 = acc_7
	goto b4
b3:
	;
	form_16, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{remaining_9})
	if callErr != nil {
		return nil, callErr
	}
	k_18, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.pipeline", "defn-key").Deref(), []vm.Value{form_16})
	if callErr != nil {
		return nil, callErr
	}
	v20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{remaining_9})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(k_18) {
		forms_21 = forms_8
		remaining_22 = remaining_9
		acc_23 = acc_10
		form_24 = form_16
		k_25 = k_18
		goto b5
	} else {
		forms_26 = forms_8
		remaining_27 = remaining_9
		acc_28 = acc_10
		form_29 = form_16
		k_30 = k_18
		goto b6
	}
b4:
	;
	return v43, nil
b5:
	;
	v33, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "assoc").Deref(), []vm.Value{acc_23, k_25, form_24})
	if callErr != nil {
		return nil, callErr
	}
	v36 = v33
	forms_37 = forms_21
	remaining_38 = remaining_22
	acc_39 = acc_23
	form_40 = form_24
	k_41 = k_25
	goto b7
b6:
	;
	v36 = acc_28
	forms_37 = forms_26
	remaining_38 = remaining_27
	acc_39 = acc_28
	form_40 = form_29
	k_41 = k_30
	goto b7
b7:
	;
	remaining_1 = v20
	acc_2 = v36
	goto b1
}
func order_defn_forms(arg0 vm.Value) (vm.Value, error) {
	var keyed_2 vm.Value
	var arg__29434_5 vm.Value
	var arg__29440_9 vm.Value
	var arg__29441_10 vm.Value
	var arg__29447_14 vm.Value
	var arg__29453_18 vm.Value
	var arg__29454_19 vm.Value
	var own_names_20 vm.Value
	var arg__29512_23 vm.Value
	var arg__29518_26 vm.Value
	var arg__29519_27 vm.Value
	var arg__29525_30 vm.Value
	var arg__29531_33 vm.Value
	var arg__29532_34 vm.Value
	var sorted_ks_35 vm.Value
	var sorted_forms_37 vm.Value
	var topo_44 vm.Value
	var forms_45 vm.Value
	var keyed_46 vm.Value
	var own_names_47 vm.Value
	var cmp_48 vm.Value
	var sorted_ks_49 vm.Value
	var sorted_forms_50 vm.Value
	var deps_fn_51 vm.Value
	var or__x_52 vm.Value
	var topo_53 vm.Value
	var forms_54 vm.Value
	var keyed_55 vm.Value
	var own_names_56 vm.Value
	var cmp_57 vm.Value
	var sorted_ks_58 vm.Value
	var sorted_forms_59 vm.Value
	var deps_fn_60 vm.Value
	var or__x_61 vm.Value
	var topo_62 vm.Value
	var ordered_66 vm.Value
	var forms_67 vm.Value
	var keyed_68 vm.Value
	var own_names_69 vm.Value
	var cmp_70 vm.Value
	var sorted_ks_71 vm.Value
	var sorted_forms_72 vm.Value
	var deps_fn_73 vm.Value
	var or__x_74 vm.Value
	var topo_75 vm.Value
	var arg__29622_79 vm.Value
	var arg__29647_84 vm.Value
	var non_defn_85 vm.Value
	var v87 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = keyed_2, arg__29434_5, arg__29440_9, arg__29441_10, arg__29447_14, arg__29453_18, arg__29454_19, own_names_20, arg__29512_23, arg__29518_26, arg__29519_27, arg__29525_30, arg__29531_33, arg__29532_34, sorted_ks_35, sorted_forms_37, topo_44, forms_45, keyed_46, own_names_47, cmp_48, sorted_ks_49, sorted_forms_50, deps_fn_51, or__x_52, topo_53, forms_54, keyed_55, own_names_56, cmp_57, sorted_ks_58, sorted_forms_59, deps_fn_60, or__x_61, topo_62, ordered_66, forms_67, keyed_68, own_names_69, cmp_70, sorted_ks_71, sorted_forms_72, deps_fn_73, or__x_74, topo_75, arg__29622_79, arg__29647_84, non_defn_85, v87
	keyed_2, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.pipeline", "forms-by-arity").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__29434_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "keys").Deref(), []vm.Value{keyed_2})
	if callErr != nil {
		return nil, callErr
	}
	arg__29440_9, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "keys").Deref(), []vm.Value{keyed_2})
	if callErr != nil {
		return nil, callErr
	}
	arg__29441_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.LookupVar("clojure.core", "first").Deref(), arg__29440_9})
	if callErr != nil {
		return nil, callErr
	}
	arg__29447_14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "keys").Deref(), []vm.Value{keyed_2})
	if callErr != nil {
		return nil, callErr
	}
	arg__29453_18, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "keys").Deref(), []vm.Value{keyed_2})
	if callErr != nil {
		return nil, callErr
	}
	arg__29454_19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.LookupVar("clojure.core", "first").Deref(), arg__29453_18})
	if callErr != nil {
		return nil, callErr
	}
	own_names_20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "set").Deref(), []vm.Value{arg__29454_19})
	if callErr != nil {
		return nil, callErr
	}
	arg__29512_23, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "keys").Deref(), []vm.Value{keyed_2})
	if callErr != nil {
		return nil, callErr
	}
	arg__29518_26, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "keys").Deref(), []vm.Value{keyed_2})
	if callErr != nil {
		return nil, callErr
	}
	arg__29519_27, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "sort").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
		var na_7 vm.Value
		var aa_13 vm.Value
		var nb_19 vm.Value
		var ab_25 vm.Value
		var c_27 vm.Value
		var v43 vm.Value
		var vec__29425_28 vm.Value
		var vec__29426_29 vm.Value
		var na_30 vm.Value
		var aa_31 vm.Value
		var nb_32 vm.Value
		var ab_33 vm.Value
		var c_34 vm.Value
		var arg__29494_46 vm.Value
		var arg__29498_48 vm.Value
		var arg__29503_51 vm.Value
		var arg__29507_53 vm.Value
		var v54 vm.Value
		var vec__29425_35 vm.Value
		var vec__29426_36 vm.Value
		var na_37 vm.Value
		var aa_38 vm.Value
		var nb_39 vm.Value
		var ab_40 vm.Value
		var c_41 vm.Value
		var v57 vm.Value
		var vec__29425_58 vm.Value
		var vec__29426_59 vm.Value
		var na_60 vm.Value
		var aa_61 vm.Value
		var nb_62 vm.Value
		var ab_63 vm.Value
		var c_64 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = na_7, aa_13, nb_19, ab_25, c_27, v43, vec__29425_28, vec__29426_29, na_30, aa_31, nb_32, ab_33, c_34, arg__29494_46, arg__29498_48, arg__29503_51, arg__29507_53, v54, vec__29425_35, vec__29426_36, na_37, aa_38, nb_39, ab_40, c_41, v57, vec__29425_58, vec__29426_59, na_60, aa_61, nb_62, ab_63, c_64
		na_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(0), vm.NIL})
		if callErr != nil {
			return nil, callErr
		}
		aa_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(1), vm.NIL})
		if callErr != nil {
			return nil, callErr
		}
		nb_19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg1, vm.Int(0), vm.NIL})
		if callErr != nil {
			return nil, callErr
		}
		ab_25, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg1, vm.Int(1), vm.NIL})
		if callErr != nil {
			return nil, callErr
		}
		c_27, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "compare").Deref(), []vm.Value{na_7, nb_19})
		if callErr != nil {
			return nil, callErr
		}
		v43, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "zero?").Deref(), []vm.Value{c_27})
		if callErr != nil {
			return nil, callErr
		}
		if vm.IsTruthy(v43) {
			vec__29425_28 = arg0
			vec__29426_29 = arg1
			na_30 = na_7
			aa_31 = aa_13
			nb_32 = nb_19
			ab_33 = ab_25
			c_34 = c_27
			goto b1
		} else {
			vec__29425_35 = arg0
			vec__29426_36 = arg1
			na_37 = na_7
			aa_38 = aa_13
			nb_39 = nb_19
			ab_40 = ab_25
			c_41 = c_27
			goto b2
		}
	b1:
		;
		arg__29494_46, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{aa_31})
		if callErr != nil {
			return nil, callErr
		}
		arg__29498_48, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{ab_33})
		if callErr != nil {
			return nil, callErr
		}
		arg__29503_51, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{aa_31})
		if callErr != nil {
			return nil, callErr
		}
		arg__29507_53, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{ab_33})
		if callErr != nil {
			return nil, callErr
		}
		v54, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "compare").Deref(), []vm.Value{arg__29503_51, arg__29507_53})
		if callErr != nil {
			return nil, callErr
		}
		v57 = v54
		vec__29425_58 = vec__29425_28
		vec__29426_59 = vec__29426_29
		na_60 = na_30
		aa_61 = aa_31
		nb_62 = nb_32
		ab_63 = ab_33
		c_64 = c_34
		goto b3
	b2:
		;
		v57 = c_41
		vec__29425_58 = vec__29425_35
		vec__29426_59 = vec__29426_36
		na_60 = na_37
		aa_61 = aa_38
		nb_62 = nb_39
		ab_63 = ab_40
		c_64 = c_41
		goto b3
	b3:
		;
		return v57, nil
	}), arg__29518_26})
	if callErr != nil {
		return nil, callErr
	}
	arg__29525_30, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "keys").Deref(), []vm.Value{keyed_2})
	if callErr != nil {
		return nil, callErr
	}
	arg__29531_33, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "keys").Deref(), []vm.Value{keyed_2})
	if callErr != nil {
		return nil, callErr
	}
	arg__29532_34, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "sort").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
		var na_7 vm.Value
		var aa_13 vm.Value
		var nb_19 vm.Value
		var ab_25 vm.Value
		var c_27 vm.Value
		var v43 vm.Value
		var vec__29425_28 vm.Value
		var vec__29426_29 vm.Value
		var na_30 vm.Value
		var aa_31 vm.Value
		var nb_32 vm.Value
		var ab_33 vm.Value
		var c_34 vm.Value
		var arg__29494_46 vm.Value
		var arg__29498_48 vm.Value
		var arg__29503_51 vm.Value
		var arg__29507_53 vm.Value
		var v54 vm.Value
		var vec__29425_35 vm.Value
		var vec__29426_36 vm.Value
		var na_37 vm.Value
		var aa_38 vm.Value
		var nb_39 vm.Value
		var ab_40 vm.Value
		var c_41 vm.Value
		var v57 vm.Value
		var vec__29425_58 vm.Value
		var vec__29426_59 vm.Value
		var na_60 vm.Value
		var aa_61 vm.Value
		var nb_62 vm.Value
		var ab_63 vm.Value
		var c_64 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = na_7, aa_13, nb_19, ab_25, c_27, v43, vec__29425_28, vec__29426_29, na_30, aa_31, nb_32, ab_33, c_34, arg__29494_46, arg__29498_48, arg__29503_51, arg__29507_53, v54, vec__29425_35, vec__29426_36, na_37, aa_38, nb_39, ab_40, c_41, v57, vec__29425_58, vec__29426_59, na_60, aa_61, nb_62, ab_63, c_64
		na_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(0), vm.NIL})
		if callErr != nil {
			return nil, callErr
		}
		aa_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(1), vm.NIL})
		if callErr != nil {
			return nil, callErr
		}
		nb_19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg1, vm.Int(0), vm.NIL})
		if callErr != nil {
			return nil, callErr
		}
		ab_25, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg1, vm.Int(1), vm.NIL})
		if callErr != nil {
			return nil, callErr
		}
		c_27, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "compare").Deref(), []vm.Value{na_7, nb_19})
		if callErr != nil {
			return nil, callErr
		}
		v43, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "zero?").Deref(), []vm.Value{c_27})
		if callErr != nil {
			return nil, callErr
		}
		if vm.IsTruthy(v43) {
			vec__29425_28 = arg0
			vec__29426_29 = arg1
			na_30 = na_7
			aa_31 = aa_13
			nb_32 = nb_19
			ab_33 = ab_25
			c_34 = c_27
			goto b1
		} else {
			vec__29425_35 = arg0
			vec__29426_36 = arg1
			na_37 = na_7
			aa_38 = aa_13
			nb_39 = nb_19
			ab_40 = ab_25
			c_41 = c_27
			goto b2
		}
	b1:
		;
		arg__29494_46, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{aa_31})
		if callErr != nil {
			return nil, callErr
		}
		arg__29498_48, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{ab_33})
		if callErr != nil {
			return nil, callErr
		}
		arg__29503_51, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{aa_31})
		if callErr != nil {
			return nil, callErr
		}
		arg__29507_53, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{ab_33})
		if callErr != nil {
			return nil, callErr
		}
		v54, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "compare").Deref(), []vm.Value{arg__29503_51, arg__29507_53})
		if callErr != nil {
			return nil, callErr
		}
		v57 = v54
		vec__29425_58 = vec__29425_28
		vec__29426_59 = vec__29426_29
		na_60 = na_30
		aa_61 = aa_31
		nb_62 = nb_32
		ab_63 = ab_33
		c_64 = c_34
		goto b3
	b2:
		;
		v57 = c_41
		vec__29425_58 = vec__29425_35
		vec__29426_59 = vec__29426_36
		na_60 = na_37
		aa_61 = aa_38
		nb_62 = nb_39
		ab_63 = ab_40
		c_64 = c_41
		goto b3
	b3:
		;
		return v57, nil
	}), arg__29531_33})
	if callErr != nil {
		return nil, callErr
	}
	sorted_ks_35, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg__29532_34})
	if callErr != nil {
		return nil, callErr
	}
	sorted_forms_37, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapv").Deref(), []vm.Value{keyed_2, sorted_ks_35})
	if callErr != nil {
		return nil, callErr
	}
	topo_44, callErr = rt.InvokeValue(rt.LookupVar("graph", "toposort-by").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var arg__29590_2 vm.Value
		var arg__29595_5 vm.Value
		var v6 vm.Value
		var callErr error
		_, _, _ = arg__29590_2, arg__29595_5, v6
		arg__29590_2, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.pipeline", "defn-key").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		arg__29595_5, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.pipeline", "defn-key").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		v6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{arg__29595_5})
		if callErr != nil {
			return nil, callErr
		}
		return v6, nil
	}), sorted_forms_37, rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var arg__29541_3 vm.Value
		var arg__29546_6 vm.Value
		var self_name_7 vm.Value
		var arg__29559_13 vm.Value
		var arg__29573_20 vm.Value
		var v21 vm.Value
		var callErr error
		_, _, _, _, _, _ = arg__29541_3, arg__29546_6, self_name_7, arg__29559_13, arg__29573_20, v21
		arg__29541_3, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.pipeline", "defn-key").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		arg__29546_6, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.pipeline", "defn-key").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		self_name_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{arg__29546_6})
		if callErr != nil {
			return nil, callErr
		}
		arg__29559_13, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.pipeline", "collect-call-targets").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		arg__29573_20, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.pipeline", "collect-call-targets").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		v21, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "filterv").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
			var and__x_4 vm.Value
			var n_5 vm.Value
			var own_names_6 vm.Value
			var self_name_7 vm.Value
			var and__x_8 vm.Value
			var v14 vm.Value
			var n_9 vm.Value
			var own_names_10 vm.Value
			var self_name_11 vm.Value
			var and__x_12 vm.Value
			var v17 vm.Value
			var n_18 vm.Value
			var own_names_19 vm.Value
			var self_name_20 vm.Value
			var and__x_21 vm.Value
			var callErr error
			_, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = and__x_4, n_5, own_names_6, self_name_7, and__x_8, v14, n_9, own_names_10, self_name_11, and__x_12, v17, n_18, own_names_19, self_name_20, and__x_21
			and__x_4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not=").Deref(), []vm.Value{arg0, self_name_7})
			if callErr != nil {
				return nil, callErr
			}
			if vm.IsTruthy(and__x_4) {
				n_5 = arg0
				own_names_6 = own_names_20
				self_name_7 = self_name_7
				and__x_8 = and__x_4
				goto b1
			} else {
				n_9 = arg0
				own_names_10 = own_names_20
				self_name_11 = self_name_7
				and__x_12 = and__x_4
				goto b2
			}
		b1:
			;
			v14, callErr = rt.InvokeValue(own_names_6, []vm.Value{n_5})
			if callErr != nil {
				return nil, callErr
			}
			v17 = v14
			n_18 = n_5
			own_names_19 = own_names_6
			self_name_20 = self_name_7
			and__x_21 = and__x_8
			goto b3
		b2:
			;
			v17 = and__x_12
			n_18 = n_9
			own_names_19 = own_names_10
			self_name_20 = self_name_11
			and__x_21 = and__x_12
			goto b3
		b3:
			;
			return v17, nil
		}), arg__29573_20})
		if callErr != nil {
			return nil, callErr
		}
		return v21, nil
	})})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(topo_44) {
		forms_45 = arg0
		keyed_46 = keyed_2
		own_names_47 = own_names_20
		cmp_48 = rt.BoxNativeFn(func(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
			var na_7 vm.Value
			var aa_13 vm.Value
			var nb_19 vm.Value
			var ab_25 vm.Value
			var c_27 vm.Value
			var v43 vm.Value
			var vec__29425_28 vm.Value
			var vec__29426_29 vm.Value
			var na_30 vm.Value
			var aa_31 vm.Value
			var nb_32 vm.Value
			var ab_33 vm.Value
			var c_34 vm.Value
			var arg__29494_46 vm.Value
			var arg__29498_48 vm.Value
			var arg__29503_51 vm.Value
			var arg__29507_53 vm.Value
			var v54 vm.Value
			var vec__29425_35 vm.Value
			var vec__29426_36 vm.Value
			var na_37 vm.Value
			var aa_38 vm.Value
			var nb_39 vm.Value
			var ab_40 vm.Value
			var c_41 vm.Value
			var v57 vm.Value
			var vec__29425_58 vm.Value
			var vec__29426_59 vm.Value
			var na_60 vm.Value
			var aa_61 vm.Value
			var nb_62 vm.Value
			var ab_63 vm.Value
			var c_64 vm.Value
			var callErr error
			_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = na_7, aa_13, nb_19, ab_25, c_27, v43, vec__29425_28, vec__29426_29, na_30, aa_31, nb_32, ab_33, c_34, arg__29494_46, arg__29498_48, arg__29503_51, arg__29507_53, v54, vec__29425_35, vec__29426_36, na_37, aa_38, nb_39, ab_40, c_41, v57, vec__29425_58, vec__29426_59, na_60, aa_61, nb_62, ab_63, c_64
			na_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(0), vm.NIL})
			if callErr != nil {
				return nil, callErr
			}
			aa_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(1), vm.NIL})
			if callErr != nil {
				return nil, callErr
			}
			nb_19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg1, vm.Int(0), vm.NIL})
			if callErr != nil {
				return nil, callErr
			}
			ab_25, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg1, vm.Int(1), vm.NIL})
			if callErr != nil {
				return nil, callErr
			}
			c_27, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "compare").Deref(), []vm.Value{na_7, nb_19})
			if callErr != nil {
				return nil, callErr
			}
			v43, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "zero?").Deref(), []vm.Value{c_27})
			if callErr != nil {
				return nil, callErr
			}
			if vm.IsTruthy(v43) {
				vec__29425_28 = arg0
				vec__29426_29 = arg1
				na_30 = na_7
				aa_31 = aa_13
				nb_32 = nb_19
				ab_33 = ab_25
				c_34 = c_27
				goto b1
			} else {
				vec__29425_35 = arg0
				vec__29426_36 = arg1
				na_37 = na_7
				aa_38 = aa_13
				nb_39 = nb_19
				ab_40 = ab_25
				c_41 = c_27
				goto b2
			}
		b1:
			;
			arg__29494_46, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{aa_31})
			if callErr != nil {
				return nil, callErr
			}
			arg__29498_48, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{ab_33})
			if callErr != nil {
				return nil, callErr
			}
			arg__29503_51, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{aa_31})
			if callErr != nil {
				return nil, callErr
			}
			arg__29507_53, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{ab_33})
			if callErr != nil {
				return nil, callErr
			}
			v54, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "compare").Deref(), []vm.Value{arg__29503_51, arg__29507_53})
			if callErr != nil {
				return nil, callErr
			}
			v57 = v54
			vec__29425_58 = vec__29425_28
			vec__29426_59 = vec__29426_29
			na_60 = na_30
			aa_61 = aa_31
			nb_62 = nb_32
			ab_63 = ab_33
			c_64 = c_34
			goto b3
		b2:
			;
			v57 = c_41
			vec__29425_58 = vec__29425_35
			vec__29426_59 = vec__29426_36
			na_60 = na_37
			aa_61 = aa_38
			nb_62 = nb_39
			ab_63 = ab_40
			c_64 = c_41
			goto b3
		b3:
			;
			return v57, nil
		})
		sorted_ks_49 = sorted_ks_35
		sorted_forms_50 = sorted_forms_37
		deps_fn_51 = rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
			var arg__29541_3 vm.Value
			var arg__29546_6 vm.Value
			var self_name_7 vm.Value
			var arg__29559_13 vm.Value
			var arg__29573_20 vm.Value
			var v21 vm.Value
			var callErr error
			_, _, _, _, _, _ = arg__29541_3, arg__29546_6, self_name_7, arg__29559_13, arg__29573_20, v21
			arg__29541_3, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.pipeline", "defn-key").Deref(), []vm.Value{arg0})
			if callErr != nil {
				return nil, callErr
			}
			arg__29546_6, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.pipeline", "defn-key").Deref(), []vm.Value{arg0})
			if callErr != nil {
				return nil, callErr
			}
			self_name_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{arg__29546_6})
			if callErr != nil {
				return nil, callErr
			}
			arg__29559_13, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.pipeline", "collect-call-targets").Deref(), []vm.Value{arg0})
			if callErr != nil {
				return nil, callErr
			}
			arg__29573_20, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.pipeline", "collect-call-targets").Deref(), []vm.Value{arg0})
			if callErr != nil {
				return nil, callErr
			}
			v21, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "filterv").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
				var and__x_4 vm.Value
				var n_5 vm.Value
				var own_names_6 vm.Value
				var self_name_7 vm.Value
				var and__x_8 vm.Value
				var v14 vm.Value
				var n_9 vm.Value
				var own_names_10 vm.Value
				var self_name_11 vm.Value
				var and__x_12 vm.Value
				var v17 vm.Value
				var n_18 vm.Value
				var own_names_19 vm.Value
				var self_name_20 vm.Value
				var and__x_21 vm.Value
				var callErr error
				_, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = and__x_4, n_5, own_names_6, self_name_7, and__x_8, v14, n_9, own_names_10, self_name_11, and__x_12, v17, n_18, own_names_19, self_name_20, and__x_21
				and__x_4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not=").Deref(), []vm.Value{arg0, self_name_7})
				if callErr != nil {
					return nil, callErr
				}
				if vm.IsTruthy(and__x_4) {
					n_5 = arg0
					own_names_6 = own_names_20
					self_name_7 = self_name_7
					and__x_8 = and__x_4
					goto b1
				} else {
					n_9 = arg0
					own_names_10 = own_names_20
					self_name_11 = self_name_7
					and__x_12 = and__x_4
					goto b2
				}
			b1:
				;
				v14, callErr = rt.InvokeValue(own_names_6, []vm.Value{n_5})
				if callErr != nil {
					return nil, callErr
				}
				v17 = v14
				n_18 = n_5
				own_names_19 = own_names_6
				self_name_20 = self_name_7
				and__x_21 = and__x_8
				goto b3
			b2:
				;
				v17 = and__x_12
				n_18 = n_9
				own_names_19 = own_names_10
				self_name_20 = self_name_11
				and__x_21 = and__x_12
				goto b3
			b3:
				;
				return v17, nil
			}), arg__29573_20})
			if callErr != nil {
				return nil, callErr
			}
			return v21, nil
		})
		or__x_52 = topo_44
		topo_53 = topo_44
		goto b1
	} else {
		forms_54 = arg0
		keyed_55 = keyed_2
		own_names_56 = own_names_20
		cmp_57 = rt.BoxNativeFn(func(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
			var na_7 vm.Value
			var aa_13 vm.Value
			var nb_19 vm.Value
			var ab_25 vm.Value
			var c_27 vm.Value
			var v43 vm.Value
			var vec__29425_28 vm.Value
			var vec__29426_29 vm.Value
			var na_30 vm.Value
			var aa_31 vm.Value
			var nb_32 vm.Value
			var ab_33 vm.Value
			var c_34 vm.Value
			var arg__29494_46 vm.Value
			var arg__29498_48 vm.Value
			var arg__29503_51 vm.Value
			var arg__29507_53 vm.Value
			var v54 vm.Value
			var vec__29425_35 vm.Value
			var vec__29426_36 vm.Value
			var na_37 vm.Value
			var aa_38 vm.Value
			var nb_39 vm.Value
			var ab_40 vm.Value
			var c_41 vm.Value
			var v57 vm.Value
			var vec__29425_58 vm.Value
			var vec__29426_59 vm.Value
			var na_60 vm.Value
			var aa_61 vm.Value
			var nb_62 vm.Value
			var ab_63 vm.Value
			var c_64 vm.Value
			var callErr error
			_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = na_7, aa_13, nb_19, ab_25, c_27, v43, vec__29425_28, vec__29426_29, na_30, aa_31, nb_32, ab_33, c_34, arg__29494_46, arg__29498_48, arg__29503_51, arg__29507_53, v54, vec__29425_35, vec__29426_36, na_37, aa_38, nb_39, ab_40, c_41, v57, vec__29425_58, vec__29426_59, na_60, aa_61, nb_62, ab_63, c_64
			na_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(0), vm.NIL})
			if callErr != nil {
				return nil, callErr
			}
			aa_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(1), vm.NIL})
			if callErr != nil {
				return nil, callErr
			}
			nb_19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg1, vm.Int(0), vm.NIL})
			if callErr != nil {
				return nil, callErr
			}
			ab_25, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg1, vm.Int(1), vm.NIL})
			if callErr != nil {
				return nil, callErr
			}
			c_27, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "compare").Deref(), []vm.Value{na_7, nb_19})
			if callErr != nil {
				return nil, callErr
			}
			v43, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "zero?").Deref(), []vm.Value{c_27})
			if callErr != nil {
				return nil, callErr
			}
			if vm.IsTruthy(v43) {
				vec__29425_28 = arg0
				vec__29426_29 = arg1
				na_30 = na_7
				aa_31 = aa_13
				nb_32 = nb_19
				ab_33 = ab_25
				c_34 = c_27
				goto b1
			} else {
				vec__29425_35 = arg0
				vec__29426_36 = arg1
				na_37 = na_7
				aa_38 = aa_13
				nb_39 = nb_19
				ab_40 = ab_25
				c_41 = c_27
				goto b2
			}
		b1:
			;
			arg__29494_46, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{aa_31})
			if callErr != nil {
				return nil, callErr
			}
			arg__29498_48, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{ab_33})
			if callErr != nil {
				return nil, callErr
			}
			arg__29503_51, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{aa_31})
			if callErr != nil {
				return nil, callErr
			}
			arg__29507_53, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{ab_33})
			if callErr != nil {
				return nil, callErr
			}
			v54, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "compare").Deref(), []vm.Value{arg__29503_51, arg__29507_53})
			if callErr != nil {
				return nil, callErr
			}
			v57 = v54
			vec__29425_58 = vec__29425_28
			vec__29426_59 = vec__29426_29
			na_60 = na_30
			aa_61 = aa_31
			nb_62 = nb_32
			ab_63 = ab_33
			c_64 = c_34
			goto b3
		b2:
			;
			v57 = c_41
			vec__29425_58 = vec__29425_35
			vec__29426_59 = vec__29426_36
			na_60 = na_37
			aa_61 = aa_38
			nb_62 = nb_39
			ab_63 = ab_40
			c_64 = c_41
			goto b3
		b3:
			;
			return v57, nil
		})
		sorted_ks_58 = sorted_ks_35
		sorted_forms_59 = sorted_forms_37
		deps_fn_60 = rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
			var arg__29541_3 vm.Value
			var arg__29546_6 vm.Value
			var self_name_7 vm.Value
			var arg__29559_13 vm.Value
			var arg__29573_20 vm.Value
			var v21 vm.Value
			var callErr error
			_, _, _, _, _, _ = arg__29541_3, arg__29546_6, self_name_7, arg__29559_13, arg__29573_20, v21
			arg__29541_3, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.pipeline", "defn-key").Deref(), []vm.Value{arg0})
			if callErr != nil {
				return nil, callErr
			}
			arg__29546_6, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.pipeline", "defn-key").Deref(), []vm.Value{arg0})
			if callErr != nil {
				return nil, callErr
			}
			self_name_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{arg__29546_6})
			if callErr != nil {
				return nil, callErr
			}
			arg__29559_13, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.pipeline", "collect-call-targets").Deref(), []vm.Value{arg0})
			if callErr != nil {
				return nil, callErr
			}
			arg__29573_20, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.pipeline", "collect-call-targets").Deref(), []vm.Value{arg0})
			if callErr != nil {
				return nil, callErr
			}
			v21, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "filterv").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
				var and__x_4 vm.Value
				var n_5 vm.Value
				var own_names_6 vm.Value
				var self_name_7 vm.Value
				var and__x_8 vm.Value
				var v14 vm.Value
				var n_9 vm.Value
				var own_names_10 vm.Value
				var self_name_11 vm.Value
				var and__x_12 vm.Value
				var v17 vm.Value
				var n_18 vm.Value
				var own_names_19 vm.Value
				var self_name_20 vm.Value
				var and__x_21 vm.Value
				var callErr error
				_, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = and__x_4, n_5, own_names_6, self_name_7, and__x_8, v14, n_9, own_names_10, self_name_11, and__x_12, v17, n_18, own_names_19, self_name_20, and__x_21
				and__x_4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not=").Deref(), []vm.Value{arg0, self_name_7})
				if callErr != nil {
					return nil, callErr
				}
				if vm.IsTruthy(and__x_4) {
					n_5 = arg0
					own_names_6 = own_names_20
					self_name_7 = self_name_7
					and__x_8 = and__x_4
					goto b1
				} else {
					n_9 = arg0
					own_names_10 = own_names_20
					self_name_11 = self_name_7
					and__x_12 = and__x_4
					goto b2
				}
			b1:
				;
				v14, callErr = rt.InvokeValue(own_names_6, []vm.Value{n_5})
				if callErr != nil {
					return nil, callErr
				}
				v17 = v14
				n_18 = n_5
				own_names_19 = own_names_6
				self_name_20 = self_name_7
				and__x_21 = and__x_8
				goto b3
			b2:
				;
				v17 = and__x_12
				n_18 = n_9
				own_names_19 = own_names_10
				self_name_20 = self_name_11
				and__x_21 = and__x_12
				goto b3
			b3:
				;
				return v17, nil
			}), arg__29573_20})
			if callErr != nil {
				return nil, callErr
			}
			return v21, nil
		})
		or__x_61 = topo_44
		topo_62 = topo_44
		goto b2
	}
b1:
	;
	ordered_66 = or__x_52
	forms_67 = forms_45
	keyed_68 = keyed_46
	own_names_69 = own_names_47
	cmp_70 = cmp_48
	sorted_ks_71 = sorted_ks_49
	sorted_forms_72 = sorted_forms_50
	deps_fn_73 = deps_fn_51
	or__x_74 = or__x_52
	topo_75 = topo_53
	goto b3
b2:
	;
	ordered_66 = sorted_forms_59
	forms_67 = forms_54
	keyed_68 = keyed_55
	own_names_69 = own_names_56
	cmp_70 = cmp_57
	sorted_ks_71 = sorted_ks_58
	sorted_forms_72 = sorted_forms_59
	deps_fn_73 = deps_fn_60
	or__x_74 = or__x_61
	topo_75 = topo_62
	goto b3
b3:
	;
	arg__29622_79, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "remove").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var arg__29614_2 vm.Value
		var arg__29619_5 vm.Value
		var v6 vm.Value
		var callErr error
		_, _, _ = arg__29614_2, arg__29619_5, v6
		arg__29614_2, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.pipeline", "defn-key").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		arg__29619_5, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.pipeline", "defn-key").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		v6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "some?").Deref(), []vm.Value{arg__29619_5})
		if callErr != nil {
			return nil, callErr
		}
		return v6, nil
	}), forms_67})
	if callErr != nil {
		return nil, callErr
	}
	arg__29647_84, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "remove").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var arg__29639_2 vm.Value
		var arg__29644_5 vm.Value
		var v6 vm.Value
		var callErr error
		_, _, _ = arg__29639_2, arg__29644_5, v6
		arg__29639_2, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.pipeline", "defn-key").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		arg__29644_5, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.pipeline", "defn-key").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		v6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "some?").Deref(), []vm.Value{arg__29644_5})
		if callErr != nil {
			return nil, callErr
		}
		return v6, nil
	}), forms_67})
	if callErr != nil {
		return nil, callErr
	}
	non_defn_85, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg__29647_84})
	if callErr != nil {
		return nil, callErr
	}
	v87, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{ordered_66, non_defn_85})
	if callErr != nil {
		return nil, callErr
	}
	return v87, nil
}
func compile_form_STAR_(arg0 vm.Value) (vm.Value, error) {
	var arg__29658_5 vm.Value
	var arg__29665_10 vm.Value
	var name_sym_11 vm.Value
	var maybe_doc_15 vm.Value
	var has_doc_QMARK__17 vm.Value
	var form_18 vm.Value
	var name_sym_19 vm.Value
	var maybe_doc_20 vm.Value
	var has_doc_QMARK__21 vm.Value
	var v30 vm.Value
	var form_22 vm.Value
	var name_sym_23 vm.Value
	var maybe_doc_24 vm.Value
	var has_doc_QMARK__25 vm.Value
	var first_form_33 vm.Value
	var form_34 vm.Value
	var name_sym_35 vm.Value
	var maybe_doc_36 vm.Value
	var has_doc_QMARK__37 vm.Value
	var go_target_QMARK__40 bool
	var v54 vm.Value
	var first_form_41 vm.Value
	var form_42 vm.Value
	var name_sym_43 vm.Value
	var maybe_doc_44 vm.Value
	var has_doc_QMARK__45 vm.Value
	var go_target_QMARK__46 bool
	var first_form_47 vm.Value
	var form_48 vm.Value
	var name_sym_49 vm.Value
	var maybe_doc_50 vm.Value
	var has_doc_QMARK__51 vm.Value
	var go_target_QMARK__52 bool
	var v359 vm.Value
	var first_form_360 vm.Value
	var form_361 vm.Value
	var name_sym_362 vm.Value
	var maybe_doc_363 vm.Value
	var has_doc_QMARK__364 vm.Value
	var go_target_QMARK__365 vm.Value
	var args_vec_57 vm.Value
	var first_form_58 vm.Value
	var form_59 vm.Value
	var name_sym_60 vm.Value
	var maybe_doc_61 vm.Value
	var has_doc_QMARK__62 vm.Value
	var go_target_QMARK__63 bool
	var arg__29684_64 vm.Value
	var args_vec_65 vm.Value
	var first_form_66 vm.Value
	var form_67 vm.Value
	var name_sym_68 vm.Value
	var maybe_doc_69 vm.Value
	var has_doc_QMARK__70 vm.Value
	var go_target_QMARK__71 bool
	var arg__29684_72 vm.Value
	var arg__29685_78 int
	var args_vec_79 vm.Value
	var first_form_80 vm.Value
	var form_81 vm.Value
	var name_sym_82 vm.Value
	var maybe_doc_83 vm.Value
	var has_doc_QMARK__84 vm.Value
	var go_target_QMARK__85 bool
	var arg__29684_86 vm.Value
	var args_vec_88 vm.Value
	var first_form_89 vm.Value
	var form_90 vm.Value
	var name_sym_91 vm.Value
	var maybe_doc_92 vm.Value
	var has_doc_QMARK__93 vm.Value
	var go_target_QMARK__94 bool
	var arg__29684_95 vm.Value
	var head__29687_96 vm.Value
	var args_vec_97 vm.Value
	var first_form_98 vm.Value
	var form_99 vm.Value
	var name_sym_100 vm.Value
	var maybe_doc_101 vm.Value
	var has_doc_QMARK__102 vm.Value
	var go_target_QMARK__103 bool
	var arg__29684_104 vm.Value
	var head__29687_105 vm.Value
	var arg__29688_111 int
	var args_vec_112 vm.Value
	var first_form_113 vm.Value
	var form_114 vm.Value
	var name_sym_115 vm.Value
	var maybe_doc_116 vm.Value
	var has_doc_QMARK__117 vm.Value
	var go_target_QMARK__118 bool
	var arg__29684_119 vm.Value
	var head__29687_120 vm.Value
	var arg__29690_121 vm.Value
	var args_vec_124 vm.Value
	var first_form_125 vm.Value
	var form_126 vm.Value
	var name_sym_127 vm.Value
	var maybe_doc_128 vm.Value
	var has_doc_QMARK__129 vm.Value
	var go_target_QMARK__130 bool
	var head__29691_131 vm.Value
	var arg__29692_132 vm.Value
	var args_vec_133 vm.Value
	var first_form_134 vm.Value
	var form_135 vm.Value
	var name_sym_136 vm.Value
	var maybe_doc_137 vm.Value
	var has_doc_QMARK__138 vm.Value
	var go_target_QMARK__139 bool
	var head__29691_140 vm.Value
	var arg__29692_141 vm.Value
	var arg__29693_147 int
	var args_vec_148 vm.Value
	var first_form_149 vm.Value
	var form_150 vm.Value
	var name_sym_151 vm.Value
	var maybe_doc_152 vm.Value
	var has_doc_QMARK__153 vm.Value
	var go_target_QMARK__154 bool
	var head__29691_155 vm.Value
	var arg__29692_156 vm.Value
	var args_vec_158 vm.Value
	var first_form_159 vm.Value
	var form_160 vm.Value
	var name_sym_161 vm.Value
	var maybe_doc_162 vm.Value
	var has_doc_QMARK__163 vm.Value
	var go_target_QMARK__164 bool
	var head__29691_165 vm.Value
	var arg__29692_166 vm.Value
	var head__29695_167 vm.Value
	var args_vec_168 vm.Value
	var first_form_169 vm.Value
	var form_170 vm.Value
	var name_sym_171 vm.Value
	var maybe_doc_172 vm.Value
	var has_doc_QMARK__173 vm.Value
	var go_target_QMARK__174 bool
	var head__29691_175 vm.Value
	var arg__29692_176 vm.Value
	var head__29695_177 vm.Value
	var arg__29696_183 int
	var args_vec_184 vm.Value
	var first_form_185 vm.Value
	var form_186 vm.Value
	var name_sym_187 vm.Value
	var maybe_doc_188 vm.Value
	var has_doc_QMARK__189 vm.Value
	var go_target_QMARK__190 bool
	var head__29691_191 vm.Value
	var arg__29692_192 vm.Value
	var head__29695_193 vm.Value
	var arg__29698_194 vm.Value
	var body_forms_195 vm.Value
	var expanded_201 vm.Value
	var arg__29713_203 vm.Value
	var arg__29718_206 vm.Value
	var ir_fn_207 vm.Value
	var args_vec_208 vm.Value
	var first_form_209 vm.Value
	var form_210 vm.Value
	var name_sym_211 vm.Value
	var maybe_doc_212 vm.Value
	var has_doc_QMARK__213 vm.Value
	var go_target_QMARK__214 bool
	var body_forms_215 vm.Value
	var expanded_216 vm.Value
	var ir_fn_217 vm.Value
	var v232 vm.Value
	var args_vec_218 vm.Value
	var first_form_219 vm.Value
	var form_220 vm.Value
	var name_sym_221 vm.Value
	var maybe_doc_222 vm.Value
	var has_doc_QMARK__223 vm.Value
	var go_target_QMARK__224 bool
	var body_forms_225 vm.Value
	var expanded_226 vm.Value
	var ir_fn_227 vm.Value
	var v235 vm.Value
	var v237 vm.Value
	var args_vec_238 vm.Value
	var first_form_239 vm.Value
	var form_240 vm.Value
	var name_sym_241 vm.Value
	var maybe_doc_242 vm.Value
	var has_doc_QMARK__243 vm.Value
	var go_target_QMARK__244 vm.Value
	var body_forms_245 vm.Value
	var expanded_246 vm.Value
	var ir_fn_247 vm.Value
	var first_form_249 vm.Value
	var form_250 vm.Value
	var name_sym_251 vm.Value
	var maybe_doc_252 vm.Value
	var has_doc_QMARK__253 vm.Value
	var go_target_QMARK__254 bool
	var first_form_255 vm.Value
	var form_256 vm.Value
	var name_sym_257 vm.Value
	var maybe_doc_258 vm.Value
	var has_doc_QMARK__259 vm.Value
	var go_target_QMARK__260 bool
	var arg__29727_266 int
	var first_form_267 vm.Value
	var form_268 vm.Value
	var name_sym_269 vm.Value
	var maybe_doc_270 vm.Value
	var has_doc_QMARK__271 vm.Value
	var go_target_QMARK__272 bool
	var first_form_274 vm.Value
	var form_275 vm.Value
	var name_sym_276 vm.Value
	var maybe_doc_277 vm.Value
	var has_doc_QMARK__278 vm.Value
	var go_target_QMARK__279 bool
	var head__29729_280 vm.Value
	var first_form_281 vm.Value
	var form_282 vm.Value
	var name_sym_283 vm.Value
	var maybe_doc_284 vm.Value
	var has_doc_QMARK__285 vm.Value
	var go_target_QMARK__286 bool
	var head__29729_287 vm.Value
	var arg__29730_293 int
	var first_form_294 vm.Value
	var form_295 vm.Value
	var name_sym_296 vm.Value
	var maybe_doc_297 vm.Value
	var has_doc_QMARK__298 vm.Value
	var go_target_QMARK__299 bool
	var head__29729_300 vm.Value
	var arities_301 vm.Value
	var first_form_302 vm.Value
	var form_303 vm.Value
	var name_sym_304 vm.Value
	var maybe_doc_305 vm.Value
	var has_doc_QMARK__306 vm.Value
	var go_target_QMARK__307 bool
	var arities_308 vm.Value
	var fn_templates_324 vm.Value
	var v329 vm.Value
	var first_form_309 vm.Value
	var form_310 vm.Value
	var name_sym_311 vm.Value
	var maybe_doc_312 vm.Value
	var has_doc_QMARK__313 vm.Value
	var go_target_QMARK__314 bool
	var arities_315 vm.Value
	var fn_vals_338 vm.Value
	var arg__29966_342 vm.Value
	var arg__29973_347 vm.Value
	var v348 vm.Value
	var v350 vm.Value
	var first_form_351 vm.Value
	var form_352 vm.Value
	var name_sym_353 vm.Value
	var maybe_doc_354 vm.Value
	var has_doc_QMARK__355 vm.Value
	var go_target_QMARK__356 vm.Value
	var arities_357 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__29658_5, arg__29665_10, name_sym_11, maybe_doc_15, has_doc_QMARK__17, form_18, name_sym_19, maybe_doc_20, has_doc_QMARK__21, v30, form_22, name_sym_23, maybe_doc_24, has_doc_QMARK__25, first_form_33, form_34, name_sym_35, maybe_doc_36, has_doc_QMARK__37, go_target_QMARK__40, v54, first_form_41, form_42, name_sym_43, maybe_doc_44, has_doc_QMARK__45, go_target_QMARK__46, first_form_47, form_48, name_sym_49, maybe_doc_50, has_doc_QMARK__51, go_target_QMARK__52, v359, first_form_360, form_361, name_sym_362, maybe_doc_363, has_doc_QMARK__364, go_target_QMARK__365, args_vec_57, first_form_58, form_59, name_sym_60, maybe_doc_61, has_doc_QMARK__62, go_target_QMARK__63, arg__29684_64, args_vec_65, first_form_66, form_67, name_sym_68, maybe_doc_69, has_doc_QMARK__70, go_target_QMARK__71, arg__29684_72, arg__29685_78, args_vec_79, first_form_80, form_81, name_sym_82, maybe_doc_83, has_doc_QMARK__84, go_target_QMARK__85, arg__29684_86, args_vec_88, first_form_89, form_90, name_sym_91, maybe_doc_92, has_doc_QMARK__93, go_target_QMARK__94, arg__29684_95, head__29687_96, args_vec_97, first_form_98, form_99, name_sym_100, maybe_doc_101, has_doc_QMARK__102, go_target_QMARK__103, arg__29684_104, head__29687_105, arg__29688_111, args_vec_112, first_form_113, form_114, name_sym_115, maybe_doc_116, has_doc_QMARK__117, go_target_QMARK__118, arg__29684_119, head__29687_120, arg__29690_121, args_vec_124, first_form_125, form_126, name_sym_127, maybe_doc_128, has_doc_QMARK__129, go_target_QMARK__130, head__29691_131, arg__29692_132, args_vec_133, first_form_134, form_135, name_sym_136, maybe_doc_137, has_doc_QMARK__138, go_target_QMARK__139, head__29691_140, arg__29692_141, arg__29693_147, args_vec_148, first_form_149, form_150, name_sym_151, maybe_doc_152, has_doc_QMARK__153, go_target_QMARK__154, head__29691_155, arg__29692_156, args_vec_158, first_form_159, form_160, name_sym_161, maybe_doc_162, has_doc_QMARK__163, go_target_QMARK__164, head__29691_165, arg__29692_166, head__29695_167, args_vec_168, first_form_169, form_170, name_sym_171, maybe_doc_172, has_doc_QMARK__173, go_target_QMARK__174, head__29691_175, arg__29692_176, head__29695_177, arg__29696_183, args_vec_184, first_form_185, form_186, name_sym_187, maybe_doc_188, has_doc_QMARK__189, go_target_QMARK__190, head__29691_191, arg__29692_192, head__29695_193, arg__29698_194, body_forms_195, expanded_201, arg__29713_203, arg__29718_206, ir_fn_207, args_vec_208, first_form_209, form_210, name_sym_211, maybe_doc_212, has_doc_QMARK__213, go_target_QMARK__214, body_forms_215, expanded_216, ir_fn_217, v232, args_vec_218, first_form_219, form_220, name_sym_221, maybe_doc_222, has_doc_QMARK__223, go_target_QMARK__224, body_forms_225, expanded_226, ir_fn_227, v235, v237, args_vec_238, first_form_239, form_240, name_sym_241, maybe_doc_242, has_doc_QMARK__243, go_target_QMARK__244, body_forms_245, expanded_246, ir_fn_247, first_form_249, form_250, name_sym_251, maybe_doc_252, has_doc_QMARK__253, go_target_QMARK__254, first_form_255, form_256, name_sym_257, maybe_doc_258, has_doc_QMARK__259, go_target_QMARK__260, arg__29727_266, first_form_267, form_268, name_sym_269, maybe_doc_270, has_doc_QMARK__271, go_target_QMARK__272, first_form_274, form_275, name_sym_276, maybe_doc_277, has_doc_QMARK__278, go_target_QMARK__279, head__29729_280, first_form_281, form_282, name_sym_283, maybe_doc_284, has_doc_QMARK__285, go_target_QMARK__286, head__29729_287, arg__29730_293, first_form_294, form_295, name_sym_296, maybe_doc_297, has_doc_QMARK__298, go_target_QMARK__299, head__29729_300, arities_301, first_form_302, form_303, name_sym_304, maybe_doc_305, has_doc_QMARK__306, go_target_QMARK__307, arities_308, fn_templates_324, v329, first_form_309, form_310, name_sym_311, maybe_doc_312, has_doc_QMARK__313, go_target_QMARK__314, arities_315, fn_vals_338, arg__29966_342, arg__29973_347, v348, v350, first_form_351, form_352, name_sym_353, maybe_doc_354, has_doc_QMARK__355, go_target_QMARK__356, arities_357
	arg__29658_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(1)})
	if callErr != nil {
		return nil, callErr
	}
	arg__29665_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(1)})
	if callErr != nil {
		return nil, callErr
	}
	name_sym_11, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.pipeline", "unwrap-name").Deref(), []vm.Value{arg__29665_10})
	if callErr != nil {
		return nil, callErr
	}
	maybe_doc_15, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(2)})
	if callErr != nil {
		return nil, callErr
	}
	has_doc_QMARK__17, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "string?").Deref(), []vm.Value{maybe_doc_15})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(has_doc_QMARK__17) {
		form_18 = arg0
		name_sym_19 = name_sym_11
		maybe_doc_20 = maybe_doc_15
		has_doc_QMARK__21 = has_doc_QMARK__17
		goto b1
	} else {
		form_22 = arg0
		name_sym_23 = name_sym_11
		maybe_doc_24 = maybe_doc_15
		has_doc_QMARK__25 = has_doc_QMARK__17
		goto b2
	}
b1:
	;
	v30, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{form_18, vm.Int(3)})
	if callErr != nil {
		return nil, callErr
	}
	first_form_33 = v30
	form_34 = form_18
	name_sym_35 = name_sym_19
	maybe_doc_36 = maybe_doc_20
	has_doc_QMARK__37 = has_doc_QMARK__21
	goto b3
b2:
	;
	first_form_33 = maybe_doc_24
	form_34 = form_22
	name_sym_35 = name_sym_23
	maybe_doc_36 = maybe_doc_24
	has_doc_QMARK__37 = has_doc_QMARK__25
	goto b3
b3:
	;
	go_target_QMARK__40 = rt.LookupVar("ir.passes.pipeline", "*target*").Deref() == vm.Keyword("go")
	v54, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector?").Deref(), []vm.Value{first_form_33})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v54) {
		first_form_41 = first_form_33
		form_42 = form_34
		name_sym_43 = name_sym_35
		maybe_doc_44 = maybe_doc_36
		has_doc_QMARK__45 = has_doc_QMARK__37
		go_target_QMARK__46 = go_target_QMARK__40
		goto b4
	} else {
		first_form_47 = first_form_33
		form_48 = form_34
		name_sym_49 = name_sym_35
		maybe_doc_50 = maybe_doc_36
		has_doc_QMARK__51 = has_doc_QMARK__37
		go_target_QMARK__52 = go_target_QMARK__40
		goto b5
	}
b4:
	;
	if vm.IsTruthy(has_doc_QMARK__45) {
		args_vec_57 = first_form_41
		first_form_58 = first_form_41
		form_59 = form_42
		name_sym_60 = name_sym_43
		maybe_doc_61 = maybe_doc_44
		has_doc_QMARK__62 = has_doc_QMARK__45
		go_target_QMARK__63 = go_target_QMARK__46
		arg__29684_64 = rt.LookupVar("ir.passes.pipeline", "expand-all").Deref()
		goto b7
	} else {
		args_vec_65 = first_form_41
		first_form_66 = first_form_41
		form_67 = form_42
		name_sym_68 = name_sym_43
		maybe_doc_69 = maybe_doc_44
		has_doc_QMARK__70 = has_doc_QMARK__45
		go_target_QMARK__71 = go_target_QMARK__46
		arg__29684_72 = rt.LookupVar("ir.passes.pipeline", "expand-all").Deref()
		goto b8
	}
b5:
	;
	if vm.IsTruthy(has_doc_QMARK__51) {
		first_form_249 = first_form_47
		form_250 = form_48
		name_sym_251 = name_sym_49
		maybe_doc_252 = maybe_doc_50
		has_doc_QMARK__253 = has_doc_QMARK__51
		go_target_QMARK__254 = go_target_QMARK__52
		goto b22
	} else {
		first_form_255 = first_form_47
		form_256 = form_48
		name_sym_257 = name_sym_49
		maybe_doc_258 = maybe_doc_50
		has_doc_QMARK__259 = has_doc_QMARK__51
		go_target_QMARK__260 = go_target_QMARK__52
		goto b23
	}
b6:
	;
	return v359, nil
b7:
	;
	arg__29685_78 = 4
	args_vec_79 = args_vec_57
	first_form_80 = first_form_58
	form_81 = form_59
	name_sym_82 = name_sym_60
	maybe_doc_83 = maybe_doc_61
	has_doc_QMARK__84 = has_doc_QMARK__62
	go_target_QMARK__85 = go_target_QMARK__63
	arg__29684_86 = arg__29684_64
	goto b9
b8:
	;
	arg__29685_78 = 3
	args_vec_79 = args_vec_65
	first_form_80 = first_form_66
	form_81 = form_67
	name_sym_82 = name_sym_68
	maybe_doc_83 = maybe_doc_69
	has_doc_QMARK__84 = has_doc_QMARK__70
	go_target_QMARK__85 = go_target_QMARK__71
	arg__29684_86 = arg__29684_72
	goto b9
b9:
	;
	if vm.IsTruthy(has_doc_QMARK__84) {
		args_vec_88 = args_vec_79
		first_form_89 = first_form_80
		form_90 = form_81
		name_sym_91 = name_sym_82
		maybe_doc_92 = maybe_doc_83
		has_doc_QMARK__93 = has_doc_QMARK__84
		go_target_QMARK__94 = go_target_QMARK__85
		arg__29684_95 = arg__29684_86
		head__29687_96 = rt.LookupVar("clojure.core", "drop").Deref()
		goto b10
	} else {
		args_vec_97 = args_vec_79
		first_form_98 = first_form_80
		form_99 = form_81
		name_sym_100 = name_sym_82
		maybe_doc_101 = maybe_doc_83
		has_doc_QMARK__102 = has_doc_QMARK__84
		go_target_QMARK__103 = go_target_QMARK__85
		arg__29684_104 = arg__29684_86
		head__29687_105 = rt.LookupVar("clojure.core", "drop").Deref()
		goto b11
	}
b10:
	;
	arg__29688_111 = 4
	args_vec_112 = args_vec_88
	first_form_113 = first_form_89
	form_114 = form_90
	name_sym_115 = name_sym_91
	maybe_doc_116 = maybe_doc_92
	has_doc_QMARK__117 = has_doc_QMARK__93
	go_target_QMARK__118 = go_target_QMARK__94
	arg__29684_119 = arg__29684_95
	head__29687_120 = head__29687_96
	goto b12
b11:
	;
	arg__29688_111 = 3
	args_vec_112 = args_vec_97
	first_form_113 = first_form_98
	form_114 = form_99
	name_sym_115 = name_sym_100
	maybe_doc_116 = maybe_doc_101
	has_doc_QMARK__117 = has_doc_QMARK__102
	go_target_QMARK__118 = go_target_QMARK__103
	arg__29684_119 = arg__29684_104
	head__29687_120 = head__29687_105
	goto b12
b12:
	;
	arg__29690_121, callErr = rt.InvokeValue(head__29687_120, []vm.Value{vm.Int(arg__29688_111), form_114})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(has_doc_QMARK__117) {
		args_vec_124 = args_vec_112
		first_form_125 = first_form_113
		form_126 = form_114
		name_sym_127 = name_sym_115
		maybe_doc_128 = maybe_doc_116
		has_doc_QMARK__129 = has_doc_QMARK__117
		go_target_QMARK__130 = go_target_QMARK__118
		head__29691_131 = rt.LookupVar("clojure.core", "map").Deref()
		arg__29692_132 = rt.LookupVar("ir.passes.pipeline", "expand-all").Deref()
		goto b13
	} else {
		args_vec_133 = args_vec_112
		first_form_134 = first_form_113
		form_135 = form_114
		name_sym_136 = name_sym_115
		maybe_doc_137 = maybe_doc_116
		has_doc_QMARK__138 = has_doc_QMARK__117
		go_target_QMARK__139 = go_target_QMARK__118
		head__29691_140 = rt.LookupVar("clojure.core", "map").Deref()
		arg__29692_141 = rt.LookupVar("ir.passes.pipeline", "expand-all").Deref()
		goto b14
	}
b13:
	;
	arg__29693_147 = 4
	args_vec_148 = args_vec_124
	first_form_149 = first_form_125
	form_150 = form_126
	name_sym_151 = name_sym_127
	maybe_doc_152 = maybe_doc_128
	has_doc_QMARK__153 = has_doc_QMARK__129
	go_target_QMARK__154 = go_target_QMARK__130
	head__29691_155 = head__29691_131
	arg__29692_156 = arg__29692_132
	goto b15
b14:
	;
	arg__29693_147 = 3
	args_vec_148 = args_vec_133
	first_form_149 = first_form_134
	form_150 = form_135
	name_sym_151 = name_sym_136
	maybe_doc_152 = maybe_doc_137
	has_doc_QMARK__153 = has_doc_QMARK__138
	go_target_QMARK__154 = go_target_QMARK__139
	head__29691_155 = head__29691_140
	arg__29692_156 = arg__29692_141
	goto b15
b15:
	;
	if vm.IsTruthy(has_doc_QMARK__153) {
		args_vec_158 = args_vec_148
		first_form_159 = first_form_149
		form_160 = form_150
		name_sym_161 = name_sym_151
		maybe_doc_162 = maybe_doc_152
		has_doc_QMARK__163 = has_doc_QMARK__153
		go_target_QMARK__164 = go_target_QMARK__154
		head__29691_165 = head__29691_155
		arg__29692_166 = arg__29692_156
		head__29695_167 = rt.LookupVar("clojure.core", "drop").Deref()
		goto b16
	} else {
		args_vec_168 = args_vec_148
		first_form_169 = first_form_149
		form_170 = form_150
		name_sym_171 = name_sym_151
		maybe_doc_172 = maybe_doc_152
		has_doc_QMARK__173 = has_doc_QMARK__153
		go_target_QMARK__174 = go_target_QMARK__154
		head__29691_175 = head__29691_155
		arg__29692_176 = arg__29692_156
		head__29695_177 = rt.LookupVar("clojure.core", "drop").Deref()
		goto b17
	}
b16:
	;
	arg__29696_183 = 4
	args_vec_184 = args_vec_158
	first_form_185 = first_form_159
	form_186 = form_160
	name_sym_187 = name_sym_161
	maybe_doc_188 = maybe_doc_162
	has_doc_QMARK__189 = has_doc_QMARK__163
	go_target_QMARK__190 = go_target_QMARK__164
	head__29691_191 = head__29691_165
	arg__29692_192 = arg__29692_166
	head__29695_193 = head__29695_167
	goto b18
b17:
	;
	arg__29696_183 = 3
	args_vec_184 = args_vec_168
	first_form_185 = first_form_169
	form_186 = form_170
	name_sym_187 = name_sym_171
	maybe_doc_188 = maybe_doc_172
	has_doc_QMARK__189 = has_doc_QMARK__173
	go_target_QMARK__190 = go_target_QMARK__174
	head__29691_191 = head__29691_175
	arg__29692_192 = arg__29692_176
	head__29695_193 = head__29695_177
	goto b18
b18:
	;
	arg__29698_194, callErr = rt.InvokeValue(head__29695_193, []vm.Value{vm.Int(arg__29696_183), form_186})
	if callErr != nil {
		return nil, callErr
	}
	body_forms_195, callErr = rt.InvokeValue(head__29691_191, []vm.Value{arg__29692_192, arg__29698_194})
	if callErr != nil {
		return nil, callErr
	}
	expanded_201, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "apply").Deref(), []vm.Value{rt.LookupVar("clojure.core", "list").Deref(), vm.Symbol("defn"), name_sym_151, args_vec_148, body_forms_195})
	if callErr != nil {
		return nil, callErr
	}
	arg__29713_203, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-fn").Deref(), []vm.Value{expanded_201})
	if callErr != nil {
		return nil, callErr
	}
	arg__29718_206, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-fn").Deref(), []vm.Value{expanded_201})
	if callErr != nil {
		return nil, callErr
	}
	ir_fn_207, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.pipeline", "optimize-fn").Deref(), []vm.Value{arg__29718_206})
	if callErr != nil {
		return nil, callErr
	}
	if go_target_QMARK__154 {
		args_vec_208 = args_vec_148
		first_form_209 = first_form_149
		form_210 = form_150
		name_sym_211 = name_sym_151
		maybe_doc_212 = maybe_doc_152
		has_doc_QMARK__213 = has_doc_QMARK__153
		go_target_QMARK__214 = go_target_QMARK__154
		body_forms_215 = body_forms_195
		expanded_216 = expanded_201
		ir_fn_217 = ir_fn_207
		goto b19
	} else {
		args_vec_218 = args_vec_148
		first_form_219 = first_form_149
		form_220 = form_150
		name_sym_221 = name_sym_151
		maybe_doc_222 = maybe_doc_152
		has_doc_QMARK__223 = has_doc_QMARK__153
		go_target_QMARK__224 = go_target_QMARK__154
		body_forms_225 = body_forms_195
		expanded_226 = expanded_201
		ir_fn_227 = ir_fn_207
		goto b20
	}
b19:
	;
	v232, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "lower").Deref(), []vm.Value{ir_fn_217, vm.Keyword("bridge")})
	if callErr != nil {
		return nil, callErr
	}
	v237 = v232
	args_vec_238 = args_vec_208
	first_form_239 = first_form_209
	form_240 = form_210
	name_sym_241 = name_sym_211
	maybe_doc_242 = maybe_doc_212
	has_doc_QMARK__243 = has_doc_QMARK__213
	go_target_QMARK__244 = vm.Boolean(go_target_QMARK__214)
	body_forms_245 = body_forms_215
	expanded_246 = expanded_216
	ir_fn_247 = ir_fn_217
	goto b21
b20:
	;
	v235, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "lower").Deref(), []vm.Value{ir_fn_227})
	if callErr != nil {
		return nil, callErr
	}
	v237 = v235
	args_vec_238 = args_vec_218
	first_form_239 = first_form_219
	form_240 = form_220
	name_sym_241 = name_sym_221
	maybe_doc_242 = maybe_doc_222
	has_doc_QMARK__243 = has_doc_QMARK__223
	go_target_QMARK__244 = vm.Boolean(go_target_QMARK__224)
	body_forms_245 = body_forms_225
	expanded_246 = expanded_226
	ir_fn_247 = ir_fn_227
	goto b21
b21:
	;
	v359 = v237
	first_form_360 = first_form_239
	form_361 = form_240
	name_sym_362 = name_sym_241
	maybe_doc_363 = maybe_doc_242
	has_doc_QMARK__364 = has_doc_QMARK__243
	go_target_QMARK__365 = go_target_QMARK__244
	goto b6
b22:
	;
	arg__29727_266 = 3
	first_form_267 = first_form_249
	form_268 = form_250
	name_sym_269 = name_sym_251
	maybe_doc_270 = maybe_doc_252
	has_doc_QMARK__271 = has_doc_QMARK__253
	go_target_QMARK__272 = go_target_QMARK__254
	goto b24
b23:
	;
	arg__29727_266 = 2
	first_form_267 = first_form_255
	form_268 = form_256
	name_sym_269 = name_sym_257
	maybe_doc_270 = maybe_doc_258
	has_doc_QMARK__271 = has_doc_QMARK__259
	go_target_QMARK__272 = go_target_QMARK__260
	goto b24
b24:
	;
	if vm.IsTruthy(has_doc_QMARK__271) {
		first_form_274 = first_form_267
		form_275 = form_268
		name_sym_276 = name_sym_269
		maybe_doc_277 = maybe_doc_270
		has_doc_QMARK__278 = has_doc_QMARK__271
		go_target_QMARK__279 = go_target_QMARK__272
		head__29729_280 = rt.LookupVar("clojure.core", "drop").Deref()
		goto b25
	} else {
		first_form_281 = first_form_267
		form_282 = form_268
		name_sym_283 = name_sym_269
		maybe_doc_284 = maybe_doc_270
		has_doc_QMARK__285 = has_doc_QMARK__271
		go_target_QMARK__286 = go_target_QMARK__272
		head__29729_287 = rt.LookupVar("clojure.core", "drop").Deref()
		goto b26
	}
b25:
	;
	arg__29730_293 = 3
	first_form_294 = first_form_274
	form_295 = form_275
	name_sym_296 = name_sym_276
	maybe_doc_297 = maybe_doc_277
	has_doc_QMARK__298 = has_doc_QMARK__278
	go_target_QMARK__299 = go_target_QMARK__279
	head__29729_300 = head__29729_280
	goto b27
b26:
	;
	arg__29730_293 = 2
	first_form_294 = first_form_281
	form_295 = form_282
	name_sym_296 = name_sym_283
	maybe_doc_297 = maybe_doc_284
	has_doc_QMARK__298 = has_doc_QMARK__285
	go_target_QMARK__299 = go_target_QMARK__286
	head__29729_300 = head__29729_287
	goto b27
b27:
	;
	arities_301, callErr = rt.InvokeValue(head__29729_300, []vm.Value{vm.Int(arg__29730_293), form_295})
	if callErr != nil {
		return nil, callErr
	}
	if go_target_QMARK__299 {
		first_form_302 = first_form_294
		form_303 = form_295
		name_sym_304 = name_sym_296
		maybe_doc_305 = maybe_doc_297
		has_doc_QMARK__306 = has_doc_QMARK__298
		go_target_QMARK__307 = go_target_QMARK__299
		arities_308 = arities_301
		goto b28
	} else {
		first_form_309 = first_form_294
		form_310 = form_295
		name_sym_311 = name_sym_296
		maybe_doc_312 = maybe_doc_297
		has_doc_QMARK__313 = has_doc_QMARK__298
		go_target_QMARK__314 = go_target_QMARK__299
		arities_315 = arities_301
		goto b29
	}
b28:
	;
	fn_templates_324, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapv").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var args_vec_3 vm.Value
		var arg__29790_6 vm.Value
		var arg__29796_10 vm.Value
		var body_forms_11 vm.Value
		var expanded_17 vm.Value
		var arg__29811_19 vm.Value
		var arg__29816_22 vm.Value
		var ir_fn_23 vm.Value
		var result_27 vm.Value
		var arg__29825_44 vm.Value
		var v45 vm.Value
		var arity_form_28 vm.Value
		var name_sym_29 vm.Value
		var args_vec_30 vm.Value
		var body_forms_31 vm.Value
		var expanded_32 vm.Value
		var ir_fn_33 vm.Value
		var result_34 vm.Value
		var arg__29830_50 vm.Value
		var v51 vm.Value
		var arity_form_35 vm.Value
		var name_sym_36 vm.Value
		var args_vec_37 vm.Value
		var body_forms_38 vm.Value
		var expanded_39 vm.Value
		var ir_fn_40 vm.Value
		var result_41 vm.Value
		var v54 vm.Value
		var arity_form_55 vm.Value
		var name_sym_56 vm.Value
		var args_vec_57 vm.Value
		var body_forms_58 vm.Value
		var expanded_59 vm.Value
		var ir_fn_60 vm.Value
		var result_61 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = args_vec_3, arg__29790_6, arg__29796_10, body_forms_11, expanded_17, arg__29811_19, arg__29816_22, ir_fn_23, result_27, arg__29825_44, v45, arity_form_28, name_sym_29, args_vec_30, body_forms_31, expanded_32, ir_fn_33, result_34, arg__29830_50, v51, arity_form_35, name_sym_36, args_vec_37, body_forms_38, expanded_39, ir_fn_40, result_41, v54, arity_form_55, name_sym_56, args_vec_57, body_forms_58, expanded_59, ir_fn_60, result_61
		args_vec_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		arg__29790_6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		arg__29796_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		body_forms_11, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.LookupVar("ir.passes.pipeline", "expand-all").Deref(), arg__29796_10})
		if callErr != nil {
			return nil, callErr
		}
		expanded_17, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "apply").Deref(), []vm.Value{rt.LookupVar("clojure.core", "list").Deref(), vm.Symbol("defn"), name_sym_304, args_vec_3, body_forms_11})
		if callErr != nil {
			return nil, callErr
		}
		arg__29811_19, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-fn").Deref(), []vm.Value{expanded_17})
		if callErr != nil {
			return nil, callErr
		}
		arg__29816_22, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-fn").Deref(), []vm.Value{expanded_17})
		if callErr != nil {
			return nil, callErr
		}
		ir_fn_23, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.pipeline", "optimize-fn").Deref(), []vm.Value{arg__29816_22})
		if callErr != nil {
			return nil, callErr
		}
		result_27, callErr = rt.InvokeValue(rt.LookupVar("ir.lower-go", "lower").Deref(), []vm.Value{ir_fn_23, vm.Keyword("bridge")})
		if callErr != nil {
			return nil, callErr
		}
		arg__29825_44, callErr = rt.InvokeValue(vm.Keyword("status"), []vm.Value{result_27})
		if callErr != nil {
			return nil, callErr
		}
		v45 = vm.Boolean(vm.Keyword("lowered") == arg__29825_44)
		if vm.IsTruthy(v45) {
			arity_form_28 = arg0
			name_sym_29 = name_sym_304
			args_vec_30 = args_vec_3
			body_forms_31 = body_forms_11
			expanded_32 = expanded_17
			ir_fn_33 = ir_fn_23
			result_34 = result_27
			goto b1
		} else {
			arity_form_35 = arg0
			name_sym_36 = name_sym_304
			args_vec_37 = args_vec_3
			body_forms_38 = body_forms_11
			expanded_39 = expanded_17
			ir_fn_40 = ir_fn_23
			result_41 = result_27
			goto b2
		}
	b1:
		;
		arg__29830_50, callErr = rt.InvokeValue(vm.Keyword("decl"), []vm.Value{result_34})
		if callErr != nil {
			return nil, callErr
		}
		v51, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("fn"), arg__29830_50})
		if callErr != nil {
			return nil, callErr
		}
		v54 = v51
		arity_form_55 = arity_form_28
		name_sym_56 = name_sym_29
		args_vec_57 = args_vec_30
		body_forms_58 = body_forms_31
		expanded_59 = expanded_32
		ir_fn_60 = ir_fn_33
		result_61 = result_34
		goto b3
	b2:
		;
		v54 = result_41
		arity_form_55 = arity_form_35
		name_sym_56 = name_sym_36
		args_vec_57 = args_vec_37
		body_forms_58 = body_forms_38
		expanded_59 = expanded_39
		ir_fn_60 = ir_fn_40
		result_61 = result_41
		goto b3
	b3:
		;
		return v54, nil
	}), arities_308})
	if callErr != nil {
		return nil, callErr
	}
	v329, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("fns"), fn_templates_324, vm.Keyword("kind"), vm.Keyword("multi-fn-template")})
	if callErr != nil {
		return nil, callErr
	}
	v350 = v329
	first_form_351 = first_form_302
	form_352 = form_303
	name_sym_353 = name_sym_304
	maybe_doc_354 = maybe_doc_305
	has_doc_QMARK__355 = has_doc_QMARK__306
	go_target_QMARK__356 = vm.Boolean(go_target_QMARK__307)
	arities_357 = arities_308
	goto b30
b29:
	;
	fn_vals_338, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapv").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var args_vec_3 vm.Value
		var arg__29907_6 vm.Value
		var arg__29913_10 vm.Value
		var body_forms_11 vm.Value
		var expanded_17 vm.Value
		var arg__29928_19 vm.Value
		var arg__29933_22 vm.Value
		var arg__29934_23 vm.Value
		var arg__29939_26 vm.Value
		var arg__29944_29 vm.Value
		var arg__29945_30 vm.Value
		var chunk_31 vm.Value
		var arg__29949_33 vm.Value
		var arg__29956_37 vm.Value
		var v39 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = args_vec_3, arg__29907_6, arg__29913_10, body_forms_11, expanded_17, arg__29928_19, arg__29933_22, arg__29934_23, arg__29939_26, arg__29944_29, arg__29945_30, chunk_31, arg__29949_33, arg__29956_37, v39
		args_vec_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		arg__29907_6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		arg__29913_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		body_forms_11, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.LookupVar("ir.passes.pipeline", "expand-all").Deref(), arg__29913_10})
		if callErr != nil {
			return nil, callErr
		}
		expanded_17, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "apply").Deref(), []vm.Value{rt.LookupVar("clojure.core", "list").Deref(), vm.Symbol("defn"), name_sym_311, args_vec_3, body_forms_11})
		if callErr != nil {
			return nil, callErr
		}
		arg__29928_19, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-fn").Deref(), []vm.Value{expanded_17})
		if callErr != nil {
			return nil, callErr
		}
		arg__29933_22, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-fn").Deref(), []vm.Value{expanded_17})
		if callErr != nil {
			return nil, callErr
		}
		arg__29934_23, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.pipeline", "optimize-fn").Deref(), []vm.Value{arg__29933_22})
		if callErr != nil {
			return nil, callErr
		}
		arg__29939_26, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-fn").Deref(), []vm.Value{expanded_17})
		if callErr != nil {
			return nil, callErr
		}
		arg__29944_29, callErr = rt.InvokeValue(rt.LookupVar("ir.build", "build-fn").Deref(), []vm.Value{expanded_17})
		if callErr != nil {
			return nil, callErr
		}
		arg__29945_30, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.pipeline", "optimize-fn").Deref(), []vm.Value{arg__29944_29})
		if callErr != nil {
			return nil, callErr
		}
		chunk_31, callErr = rt.InvokeValue(rt.LookupVar("ir.lower", "lower").Deref(), []vm.Value{arg__29945_30})
		if callErr != nil {
			return nil, callErr
		}
		arg__29949_33, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{args_vec_3})
		if callErr != nil {
			return nil, callErr
		}
		arg__29956_37, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{args_vec_3})
		if callErr != nil {
			return nil, callErr
		}
		v39, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "chunk->fn").Deref(), []vm.Value{arg__29956_37, vm.FALSE, chunk_31})
		if callErr != nil {
			return nil, callErr
		}
		return v39, nil
	}), arities_315})
	if callErr != nil {
		return nil, callErr
	}
	arg__29966_342, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "apply").Deref(), []vm.Value{rt.LookupVar("clojure.core", "list").Deref(), fn_vals_338})
	if callErr != nil {
		return nil, callErr
	}
	arg__29973_347, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "apply").Deref(), []vm.Value{rt.LookupVar("clojure.core", "list").Deref(), fn_vals_338})
	if callErr != nil {
		return nil, callErr
	}
	v348, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "make-multi-arity").Deref(), []vm.Value{arg__29973_347})
	if callErr != nil {
		return nil, callErr
	}
	v350 = v348
	first_form_351 = first_form_309
	form_352 = form_310
	name_sym_353 = name_sym_311
	maybe_doc_354 = maybe_doc_312
	has_doc_QMARK__355 = has_doc_QMARK__313
	go_target_QMARK__356 = vm.Boolean(go_target_QMARK__314)
	arities_357 = arities_315
	goto b30
b30:
	;
	v359 = v350
	first_form_360 = first_form_351
	form_361 = form_352
	name_sym_362 = name_sym_353
	maybe_doc_363 = maybe_doc_354
	has_doc_QMARK__364 = has_doc_QMARK__355
	go_target_QMARK__365 = go_target_QMARK__356
	goto b6
}
func expand_fully(arg0 vm.Value) (vm.Value, error) {
	var f_1 vm.Value
	var e_4 vm.Value
	var v11 bool
	var form_5 vm.Value
	var f_6 vm.Value
	var e_7 vm.Value
	var form_8 vm.Value
	var f_9 vm.Value
	var e_10 vm.Value
	var v15 vm.Value
	var form_16 vm.Value
	var f_17 vm.Value
	var e_18 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _ = f_1, e_4, v11, form_5, f_6, e_7, form_8, f_9, e_10, v15, form_16, f_17, e_18
	f_1 = arg0
	goto b1
b1:
	;
	e_4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "macroexpand").Deref(), []vm.Value{f_1})
	if callErr != nil {
		return nil, callErr
	}
	v11 = e_4 == f_1
	if v11 {
		form_5 = arg0
		f_6 = f_1
		e_7 = e_4
		goto b2
	} else {
		form_8 = arg0
		f_9 = f_1
		e_10 = e_4
		goto b3
	}
b2:
	;
	v15 = e_7
	form_16 = form_5
	f_17 = f_6
	e_18 = e_7
	goto b4
b3:
	;
	f_1 = e_10
	goto b1
b4:
	;
	return v15, nil
}
func optimize_fn(arg0 vm.Value) (vm.Value, error) {
	var v4 vm.Value
	var v6 vm.Value
	var v10 vm.Value
	var v12 vm.Value
	var v14 vm.Value
	var v18 vm.Value
	var v20 vm.Value
	var v24 vm.Value
	var v26 vm.Value
	var v30 vm.Value
	var v32 vm.Value
	var v36 vm.Value
	var v38 vm.Value
	var v42 vm.Value
	var v44 vm.Value
	var v48 vm.Value
	var v50 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = v4, v6, v10, v12, v14, v18, v20, v24, v26, v30, v32, v36, v38, v42, v44, v48, v50
	v4, callErr = rt.InvokeValue(rt.LookupVar("ir.validate", "validate-fn!").Deref(), []vm.Value{arg0, vm.String("build")})
	if callErr != nil {
		return nil, callErr
	}
	v6, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.infer-arg-types", "infer-arg-types").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v10, callErr = rt.InvokeValue(rt.LookupVar("ir.validate", "validate-fn!").Deref(), []vm.Value{arg0, vm.String("infer-arg-types")})
	if callErr != nil {
		return nil, callErr
	}
	v12, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "typeinfer").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v14, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.constfold", "constfold").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v18, callErr = rt.InvokeValue(rt.LookupVar("ir.validate", "validate-fn!").Deref(), []vm.Value{arg0, vm.String("constfold-1")})
	if callErr != nil {
		return nil, callErr
	}
	v20, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.cse", "cse").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v24, callErr = rt.InvokeValue(rt.LookupVar("ir.validate", "validate-fn!").Deref(), []vm.Value{arg0, vm.String("cse")})
	if callErr != nil {
		return nil, callErr
	}
	v26, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.constfold", "constfold").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v30, callErr = rt.InvokeValue(rt.LookupVar("ir.validate", "validate-fn!").Deref(), []vm.Value{arg0, vm.String("constfold-2")})
	if callErr != nil {
		return nil, callErr
	}
	v32, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.licm", "licm").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v36, callErr = rt.InvokeValue(rt.LookupVar("ir.validate", "validate-fn!").Deref(), []vm.Value{arg0, vm.String("licm")})
	if callErr != nil {
		return nil, callErr
	}
	v38, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.constfold", "constfold").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v42, callErr = rt.InvokeValue(rt.LookupVar("ir.validate", "validate-fn!").Deref(), []vm.Value{arg0, vm.String("constfold-3")})
	if callErr != nil {
		return nil, callErr
	}
	v44, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.dce", "dce").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v48, callErr = rt.InvokeValue(rt.LookupVar("ir.validate", "validate-fn!").Deref(), []vm.Value{arg0, vm.String("dce")})
	if callErr != nil {
		return nil, callErr
	}
	v50, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "typeinfer").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	return arg0, nil
}
func expand_all(arg0 vm.Value) (vm.Value, error) {
	var v4 vm.Value
	var form_1 vm.Value
	var head_7 vm.Value
	var v15 vm.Value
	var form_2 vm.Value
	var v146 vm.Value
	var v217 vm.Value
	var form_218 vm.Value
	var form_8 vm.Value
	var head_9 vm.Value
	var arg__30124_20 vm.Value
	var arg__30130_24 vm.Value
	var arg__30131_25 vm.Value
	var arg__30137_29 vm.Value
	var arg__30143_33 vm.Value
	var arg__30144_34 vm.Value
	var arg__30145_35 vm.Value
	var arg__30153_40 vm.Value
	var arg__30159_44 vm.Value
	var arg__30160_45 vm.Value
	var arg__30166_49 vm.Value
	var arg__30172_53 vm.Value
	var arg__30173_54 vm.Value
	var arg__30174_55 vm.Value
	var v56 vm.Value
	var form_10 vm.Value
	var head_11 vm.Value
	var expanded_59 vm.Value
	var v67 vm.Value
	var v139 vm.Value
	var form_140 vm.Value
	var head_141 vm.Value
	var form_60 vm.Value
	var head_61 vm.Value
	var expanded_62 vm.Value
	var arg__30185_77 vm.Value
	var v78 bool
	var form_63 vm.Value
	var head_64 vm.Value
	var expanded_65 vm.Value
	var v132 vm.Value
	var v134 vm.Value
	var form_135 vm.Value
	var head_136 vm.Value
	var expanded_137 vm.Value
	var form_69 vm.Value
	var head_70 vm.Value
	var expanded_71 vm.Value
	var form_72 vm.Value
	var head_73 vm.Value
	var expanded_74 vm.Value
	var arg__30190_83 vm.Value
	var arg__30195_86 vm.Value
	var arg__30201_90 vm.Value
	var arg__30202_91 vm.Value
	var arg__30208_95 vm.Value
	var arg__30214_99 vm.Value
	var arg__30215_100 vm.Value
	var arg__30216_101 vm.Value
	var arg__30222_105 vm.Value
	var arg__30227_108 vm.Value
	var arg__30233_112 vm.Value
	var arg__30234_113 vm.Value
	var arg__30240_117 vm.Value
	var arg__30246_121 vm.Value
	var arg__30247_122 vm.Value
	var arg__30248_123 vm.Value
	var v124 vm.Value
	var v126 vm.Value
	var form_127 vm.Value
	var head_128 vm.Value
	var expanded_129 vm.Value
	var form_143 vm.Value
	var arg__30260_151 vm.Value
	var arg__30267_156 vm.Value
	var arg__30268_157 vm.Value
	var arg__30275_162 vm.Value
	var arg__30282_167 vm.Value
	var arg__30283_168 vm.Value
	var v169 vm.Value
	var form_144 vm.Value
	var v174 vm.Value
	var v214 vm.Value
	var form_215 vm.Value
	var form_171 vm.Value
	var arg__30339_180 vm.Value
	var arg__30392_185 vm.Value
	var arg__30393_186 vm.Value
	var arg__30447_192 vm.Value
	var arg__30500_197 vm.Value
	var arg__30501_198 vm.Value
	var v199 vm.Value
	var form_172 vm.Value
	var v211 vm.Value
	var form_212 vm.Value
	var form_201 vm.Value
	var form_202 vm.Value
	var v208 vm.Value
	var form_209 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = v4, form_1, head_7, v15, form_2, v146, v217, form_218, form_8, head_9, arg__30124_20, arg__30130_24, arg__30131_25, arg__30137_29, arg__30143_33, arg__30144_34, arg__30145_35, arg__30153_40, arg__30159_44, arg__30160_45, arg__30166_49, arg__30172_53, arg__30173_54, arg__30174_55, v56, form_10, head_11, expanded_59, v67, v139, form_140, head_141, form_60, head_61, expanded_62, arg__30185_77, v78, form_63, head_64, expanded_65, v132, v134, form_135, head_136, expanded_137, form_69, head_70, expanded_71, form_72, head_73, expanded_74, arg__30190_83, arg__30195_86, arg__30201_90, arg__30202_91, arg__30208_95, arg__30214_99, arg__30215_100, arg__30216_101, arg__30222_105, arg__30227_108, arg__30233_112, arg__30234_113, arg__30240_117, arg__30246_121, arg__30247_122, arg__30248_123, v124, v126, form_127, head_128, expanded_129, form_143, arg__30260_151, arg__30267_156, arg__30268_157, arg__30275_162, arg__30282_167, arg__30283_168, v169, form_144, v174, v214, form_215, form_171, arg__30339_180, arg__30392_185, arg__30393_186, arg__30447_192, arg__30500_197, arg__30501_198, v199, form_172, v211, form_212, form_201, form_202, v208, form_209
	v4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq?").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v4) {
		form_1 = arg0
		goto b1
	} else {
		form_2 = arg0
		goto b2
	}
b1:
	;
	head_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{form_1})
	if callErr != nil {
		return nil, callErr
	}
	v15, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "contains?").Deref(), []vm.Value{rt.LookupVar("ir.passes.pipeline", "build-known-heads").Deref(), head_7})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v15) {
		form_8 = form_1
		head_9 = head_7
		goto b4
	} else {
		form_10 = form_1
		head_11 = head_7
		goto b5
	}
b2:
	;
	v146, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector?").Deref(), []vm.Value{form_2})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v146) {
		form_143 = form_2
		goto b13
	} else {
		form_144 = form_2
		goto b14
	}
b3:
	;
	return v217, nil
b4:
	;
	arg__30124_20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{form_8})
	if callErr != nil {
		return nil, callErr
	}
	arg__30130_24, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{form_8})
	if callErr != nil {
		return nil, callErr
	}
	arg__30131_25, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.LookupVar("ir.passes.pipeline", "expand-all").Deref(), arg__30130_24})
	if callErr != nil {
		return nil, callErr
	}
	arg__30137_29, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{form_8})
	if callErr != nil {
		return nil, callErr
	}
	arg__30143_33, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{form_8})
	if callErr != nil {
		return nil, callErr
	}
	arg__30144_34, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.LookupVar("ir.passes.pipeline", "expand-all").Deref(), arg__30143_33})
	if callErr != nil {
		return nil, callErr
	}
	arg__30145_35, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "doall").Deref(), []vm.Value{arg__30144_34})
	if callErr != nil {
		return nil, callErr
	}
	arg__30153_40, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{form_8})
	if callErr != nil {
		return nil, callErr
	}
	arg__30159_44, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{form_8})
	if callErr != nil {
		return nil, callErr
	}
	arg__30160_45, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.LookupVar("ir.passes.pipeline", "expand-all").Deref(), arg__30159_44})
	if callErr != nil {
		return nil, callErr
	}
	arg__30166_49, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{form_8})
	if callErr != nil {
		return nil, callErr
	}
	arg__30172_53, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{form_8})
	if callErr != nil {
		return nil, callErr
	}
	arg__30173_54, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.LookupVar("ir.passes.pipeline", "expand-all").Deref(), arg__30172_53})
	if callErr != nil {
		return nil, callErr
	}
	arg__30174_55, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "doall").Deref(), []vm.Value{arg__30173_54})
	if callErr != nil {
		return nil, callErr
	}
	v56, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "apply").Deref(), []vm.Value{rt.LookupVar("clojure.core", "list").Deref(), head_9, arg__30174_55})
	if callErr != nil {
		return nil, callErr
	}
	v139 = v56
	form_140 = form_8
	head_141 = head_9
	goto b6
b5:
	;
	expanded_59, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.pipeline", "expand-fully").Deref(), []vm.Value{form_10})
	if callErr != nil {
		return nil, callErr
	}
	v67, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq?").Deref(), []vm.Value{expanded_59})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v67) {
		form_60 = form_10
		head_61 = head_11
		expanded_62 = expanded_59
		goto b7
	} else {
		form_63 = form_10
		head_64 = head_11
		expanded_65 = expanded_59
		goto b8
	}
b6:
	;
	v217 = v139
	form_218 = form_140
	goto b3
b7:
	;
	arg__30185_77, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{expanded_62})
	if callErr != nil {
		return nil, callErr
	}
	v78 = arg__30185_77 == vm.Symbol("quote")
	if v78 {
		form_69 = form_60
		head_70 = head_61
		expanded_71 = expanded_62
		goto b10
	} else {
		form_72 = form_60
		head_73 = head_61
		expanded_74 = expanded_62
		goto b11
	}
b8:
	;
	v132, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.pipeline", "expand-all").Deref(), []vm.Value{expanded_65})
	if callErr != nil {
		return nil, callErr
	}
	v134 = v132
	form_135 = form_63
	head_136 = head_64
	expanded_137 = expanded_65
	goto b9
b9:
	;
	v139 = v134
	form_140 = form_135
	head_141 = head_136
	goto b6
b10:
	;
	v126 = expanded_71
	form_127 = form_69
	head_128 = head_70
	expanded_129 = expanded_71
	goto b12
b11:
	;
	arg__30190_83, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{expanded_74})
	if callErr != nil {
		return nil, callErr
	}
	arg__30195_86, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{expanded_74})
	if callErr != nil {
		return nil, callErr
	}
	arg__30201_90, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{expanded_74})
	if callErr != nil {
		return nil, callErr
	}
	arg__30202_91, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.LookupVar("ir.passes.pipeline", "expand-all").Deref(), arg__30201_90})
	if callErr != nil {
		return nil, callErr
	}
	arg__30208_95, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{expanded_74})
	if callErr != nil {
		return nil, callErr
	}
	arg__30214_99, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{expanded_74})
	if callErr != nil {
		return nil, callErr
	}
	arg__30215_100, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.LookupVar("ir.passes.pipeline", "expand-all").Deref(), arg__30214_99})
	if callErr != nil {
		return nil, callErr
	}
	arg__30216_101, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "doall").Deref(), []vm.Value{arg__30215_100})
	if callErr != nil {
		return nil, callErr
	}
	arg__30222_105, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{expanded_74})
	if callErr != nil {
		return nil, callErr
	}
	arg__30227_108, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{expanded_74})
	if callErr != nil {
		return nil, callErr
	}
	arg__30233_112, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{expanded_74})
	if callErr != nil {
		return nil, callErr
	}
	arg__30234_113, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.LookupVar("ir.passes.pipeline", "expand-all").Deref(), arg__30233_112})
	if callErr != nil {
		return nil, callErr
	}
	arg__30240_117, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{expanded_74})
	if callErr != nil {
		return nil, callErr
	}
	arg__30246_121, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{expanded_74})
	if callErr != nil {
		return nil, callErr
	}
	arg__30247_122, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.LookupVar("ir.passes.pipeline", "expand-all").Deref(), arg__30246_121})
	if callErr != nil {
		return nil, callErr
	}
	arg__30248_123, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "doall").Deref(), []vm.Value{arg__30247_122})
	if callErr != nil {
		return nil, callErr
	}
	v124, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "apply").Deref(), []vm.Value{rt.LookupVar("clojure.core", "list").Deref(), arg__30222_105, arg__30248_123})
	if callErr != nil {
		return nil, callErr
	}
	v126 = v124
	form_127 = form_72
	head_128 = head_73
	expanded_129 = expanded_74
	goto b12
b12:
	;
	v134 = v126
	form_135 = form_127
	head_136 = head_128
	expanded_137 = expanded_129
	goto b9
b13:
	;
	arg__30260_151, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.LookupVar("ir.passes.pipeline", "expand-all").Deref(), form_143})
	if callErr != nil {
		return nil, callErr
	}
	arg__30267_156, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.LookupVar("ir.passes.pipeline", "expand-all").Deref(), form_143})
	if callErr != nil {
		return nil, callErr
	}
	arg__30268_157, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "doall").Deref(), []vm.Value{arg__30267_156})
	if callErr != nil {
		return nil, callErr
	}
	arg__30275_162, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.LookupVar("ir.passes.pipeline", "expand-all").Deref(), form_143})
	if callErr != nil {
		return nil, callErr
	}
	arg__30282_167, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.LookupVar("ir.passes.pipeline", "expand-all").Deref(), form_143})
	if callErr != nil {
		return nil, callErr
	}
	arg__30283_168, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "doall").Deref(), []vm.Value{arg__30282_167})
	if callErr != nil {
		return nil, callErr
	}
	v169, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg__30283_168})
	if callErr != nil {
		return nil, callErr
	}
	v214 = v169
	form_215 = form_143
	goto b15
b14:
	;
	v174, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map?").Deref(), []vm.Value{form_144})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v174) {
		form_171 = form_144
		goto b16
	} else {
		form_172 = form_144
		goto b17
	}
b15:
	;
	v217 = v214
	form_218 = form_215
	goto b3
b16:
	;
	arg__30339_180, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var k_6 vm.Value
		var v_12 vm.Value
		var arg__30332_15 vm.Value
		var arg__30336_17 vm.Value
		var v18 vm.Value
		var callErr error
		_, _, _, _, _ = k_6, v_12, arg__30332_15, arg__30336_17, v18
		k_6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(0), vm.NIL})
		if callErr != nil {
			return nil, callErr
		}
		v_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(1), vm.NIL})
		if callErr != nil {
			return nil, callErr
		}
		arg__30332_15, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.pipeline", "expand-all").Deref(), []vm.Value{k_6})
		if callErr != nil {
			return nil, callErr
		}
		arg__30336_17, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.pipeline", "expand-all").Deref(), []vm.Value{v_12})
		if callErr != nil {
			return nil, callErr
		}
		v18, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__30332_15, arg__30336_17})
		if callErr != nil {
			return nil, callErr
		}
		return v18, nil
	}), form_171})
	if callErr != nil {
		return nil, callErr
	}
	arg__30392_185, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var k_6 vm.Value
		var v_12 vm.Value
		var arg__30385_15 vm.Value
		var arg__30389_17 vm.Value
		var v18 vm.Value
		var callErr error
		_, _, _, _, _ = k_6, v_12, arg__30385_15, arg__30389_17, v18
		k_6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(0), vm.NIL})
		if callErr != nil {
			return nil, callErr
		}
		v_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(1), vm.NIL})
		if callErr != nil {
			return nil, callErr
		}
		arg__30385_15, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.pipeline", "expand-all").Deref(), []vm.Value{k_6})
		if callErr != nil {
			return nil, callErr
		}
		arg__30389_17, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.pipeline", "expand-all").Deref(), []vm.Value{v_12})
		if callErr != nil {
			return nil, callErr
		}
		v18, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__30385_15, arg__30389_17})
		if callErr != nil {
			return nil, callErr
		}
		return v18, nil
	}), form_171})
	if callErr != nil {
		return nil, callErr
	}
	arg__30393_186, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "doall").Deref(), []vm.Value{arg__30392_185})
	if callErr != nil {
		return nil, callErr
	}
	arg__30447_192, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var k_6 vm.Value
		var v_12 vm.Value
		var arg__30440_15 vm.Value
		var arg__30444_17 vm.Value
		var v18 vm.Value
		var callErr error
		_, _, _, _, _ = k_6, v_12, arg__30440_15, arg__30444_17, v18
		k_6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(0), vm.NIL})
		if callErr != nil {
			return nil, callErr
		}
		v_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(1), vm.NIL})
		if callErr != nil {
			return nil, callErr
		}
		arg__30440_15, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.pipeline", "expand-all").Deref(), []vm.Value{k_6})
		if callErr != nil {
			return nil, callErr
		}
		arg__30444_17, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.pipeline", "expand-all").Deref(), []vm.Value{v_12})
		if callErr != nil {
			return nil, callErr
		}
		v18, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__30440_15, arg__30444_17})
		if callErr != nil {
			return nil, callErr
		}
		return v18, nil
	}), form_171})
	if callErr != nil {
		return nil, callErr
	}
	arg__30500_197, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var k_6 vm.Value
		var v_12 vm.Value
		var arg__30493_15 vm.Value
		var arg__30497_17 vm.Value
		var v18 vm.Value
		var callErr error
		_, _, _, _, _ = k_6, v_12, arg__30493_15, arg__30497_17, v18
		k_6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(0), vm.NIL})
		if callErr != nil {
			return nil, callErr
		}
		v_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(1), vm.NIL})
		if callErr != nil {
			return nil, callErr
		}
		arg__30493_15, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.pipeline", "expand-all").Deref(), []vm.Value{k_6})
		if callErr != nil {
			return nil, callErr
		}
		arg__30497_17, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.pipeline", "expand-all").Deref(), []vm.Value{v_12})
		if callErr != nil {
			return nil, callErr
		}
		v18, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__30493_15, arg__30497_17})
		if callErr != nil {
			return nil, callErr
		}
		return v18, nil
	}), form_171})
	if callErr != nil {
		return nil, callErr
	}
	arg__30501_198, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "doall").Deref(), []vm.Value{arg__30500_197})
	if callErr != nil {
		return nil, callErr
	}
	v199, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{vm.EmptyPersistentMap, arg__30501_198})
	if callErr != nil {
		return nil, callErr
	}
	v211 = v199
	form_212 = form_171
	goto b18
b17:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		form_201 = form_172
		goto b19
	} else {
		form_202 = form_172
		goto b20
	}
b18:
	;
	v214 = v211
	form_215 = form_212
	goto b15
b19:
	;
	v208 = form_201
	form_209 = form_201
	goto b21
b20:
	;
	v208 = vm.NIL
	form_209 = form_202
	goto b21
b21:
	;
	v211 = v208
	form_212 = form_209
	goto b18
}
