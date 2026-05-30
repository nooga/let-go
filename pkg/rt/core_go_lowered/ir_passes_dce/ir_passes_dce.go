package ir_passes_dce

import (
	rt "github.com/nooga/let-go/pkg/rt"
	vm "github.com/nooga/let-go/pkg/vm"
)

func removable_QMARK_(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var arg__18034_3 vm.Value
	var arg__18041_6 vm.Value
	var and__x_7 vm.Value
	var nid_8 vm.Value
	var f_9 vm.Value
	var and__x_10 vm.Value
	var arg__18045_16 vm.Value
	var arg__18051_19 vm.Value
	var arg__18053_20 vm.Value
	var arg__18058_23 vm.Value
	var arg__18064_26 vm.Value
	var arg__18066_27 vm.Value
	var v28 vm.Value
	var nid_11 vm.Value
	var f_12 vm.Value
	var and__x_13 vm.Value
	var v31 vm.Value
	var nid_32 vm.Value
	var f_33 vm.Value
	var and__x_34 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__18034_3, arg__18041_6, and__x_7, nid_8, f_9, and__x_10, arg__18045_16, arg__18051_19, arg__18053_20, arg__18058_23, arg__18064_26, arg__18066_27, v28, nid_11, f_12, and__x_13, v31, nid_32, f_33, and__x_34
	arg__18034_3, callErr = rt.InvokeValue(rt.LookupVar("ir", "op").Deref(), []vm.Value{arg0, arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__18041_6, callErr = rt.InvokeValue(rt.LookupVar("ir", "op").Deref(), []vm.Value{arg0, arg1})
	if callErr != nil {
		return nil, callErr
	}
	and__x_7, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.dce", "pure-ops").Deref(), []vm.Value{arg__18041_6})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(and__x_7) {
		nid_8 = arg0
		f_9 = arg1
		and__x_10 = and__x_7
		goto b1
	} else {
		nid_11 = arg0
		f_12 = arg1
		and__x_13 = and__x_7
		goto b2
	}
b1:
	;
	arg__18045_16, callErr = rt.InvokeValue(rt.LookupVar("ir", "uses").Deref(), []vm.Value{f_9})
	if callErr != nil {
		return nil, callErr
	}
	arg__18051_19, callErr = rt.InvokeValue(rt.LookupVar("ir", "uses").Deref(), []vm.Value{f_9})
	if callErr != nil {
		return nil, callErr
	}
	arg__18053_20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__18051_19, nid_8})
	if callErr != nil {
		return nil, callErr
	}
	arg__18058_23, callErr = rt.InvokeValue(rt.LookupVar("ir", "uses").Deref(), []vm.Value{f_9})
	if callErr != nil {
		return nil, callErr
	}
	arg__18064_26, callErr = rt.InvokeValue(rt.LookupVar("ir", "uses").Deref(), []vm.Value{f_9})
	if callErr != nil {
		return nil, callErr
	}
	arg__18066_27, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__18064_26, nid_8})
	if callErr != nil {
		return nil, callErr
	}
	v28, callErr = rt.InvokeValue(rt.LookupVar("ir", "uses-empty?").Deref(), []vm.Value{arg__18066_27})
	if callErr != nil {
		return nil, callErr
	}
	v31 = v28
	nid_32 = nid_8
	f_33 = f_9
	and__x_34 = and__x_10
	goto b3
b2:
	;
	v31 = and__x_13
	nid_32 = nid_11
	f_33 = f_12
	and__x_34 = and__x_13
	goto b3
b3:
	;
	return v31, nil
}
func one_pass(arg0 vm.Value) (vm.Value, error) {
	var removed_5 vm.Value
	var arg__18079_7 vm.Value
	var arg__18084_10 vm.Value
	var doseq_seq__18069_11 vm.Value
	var doseq_loop__18070_12 vm.Value
	var f_13 vm.Value
	var removed_14 vm.Value
	var v170 string
	var v179 string
	var v188 bool
	var doseq_seq__18069_16 vm.Value
	var doseq_loop__18070_17 vm.Value
	var f_18 vm.Value
	var removed_19 vm.Value
	var v169 string
	var v178 string
	var v187 bool
	var bid_26 vm.Value
	var arg__18093_28 vm.Value
	var arg__18100_31 vm.Value
	var arg__18101_32 vm.Value
	var arg__18108_35 vm.Value
	var arg__18115_38 vm.Value
	var arg__18116_39 vm.Value
	var doseq_seq__18071_40 vm.Value
	var doseq_seq__18069_20 vm.Value
	var doseq_loop__18070_21 vm.Value
	var f_22 vm.Value
	var removed_23 vm.Value
	var v176 string
	var v185 string
	var v194 bool
	var v131 vm.Value
	var doseq_seq__18069_132 vm.Value
	var doseq_loop__18070_133 vm.Value
	var f_134 vm.Value
	var removed_135 vm.Value
	var arg__18191_138 vm.Value
	var arg__18197_142 vm.Value
	var arg__18198_143 vm.Value
	var arg__18204_147 vm.Value
	var arg__18210_151 vm.Value
	var arg__18211_152 vm.Value
	var v153 vm.Value
	var v155 vm.Value
	var doseq_loop__18072_41 vm.Value
	var f_42 vm.Value
	var removed_43 vm.Value
	var bid_44 vm.Value
	var v173 string
	var v182 string
	var v191 bool
	var doseq_seq__18069_46 vm.Value
	var doseq_loop__18070_47 vm.Value
	var doseq_seq__18071_48 vm.Value
	var doseq_loop__18072_49 vm.Value
	var f_50 vm.Value
	var removed_51 vm.Value
	var bid_52 vm.Value
	var v172 string
	var v181 string
	var v190 bool
	var nid_62 vm.Value
	var v80 vm.Value
	var doseq_seq__18069_53 vm.Value
	var doseq_loop__18070_54 vm.Value
	var doseq_seq__18071_55 vm.Value
	var doseq_loop__18072_56 vm.Value
	var f_57 vm.Value
	var removed_58 vm.Value
	var bid_59 vm.Value
	var v175 string
	var v184 string
	var v193 bool
	var v118 vm.Value
	var doseq_seq__18069_119 vm.Value
	var doseq_loop__18070_120 vm.Value
	var doseq_seq__18071_121 vm.Value
	var doseq_loop__18072_122 vm.Value
	var f_123 vm.Value
	var removed_124 vm.Value
	var bid_125 vm.Value
	var v168 string
	var v177 string
	var v186 bool
	var v127 vm.Value
	var doseq_seq__18069_63 vm.Value
	var doseq_loop__18070_64 vm.Value
	var doseq_seq__18071_65 vm.Value
	var doseq_loop__18072_66 vm.Value
	var f_67 vm.Value
	var removed_68 vm.Value
	var bid_69 vm.Value
	var nid_70 vm.Value
	var v174 string
	var v183 string
	var v192 bool
	var arg__18133_85 vm.Value
	var arg__18143_90 vm.Value
	var arg__18144_91 vm.Value
	var arg__18154_96 vm.Value
	var arg__18164_101 vm.Value
	var arg__18165_102 vm.Value
	var v103 vm.Value
	var v105 vm.Value
	var v109 vm.Value
	var v111 vm.Value
	var doseq_seq__18069_71 vm.Value
	var doseq_loop__18070_72 vm.Value
	var doseq_seq__18071_73 vm.Value
	var doseq_loop__18072_74 vm.Value
	var f_75 vm.Value
	var removed_76 vm.Value
	var bid_77 vm.Value
	var nid_78 vm.Value
	var v171 string
	var v180 string
	var v189 bool
	var v114 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = removed_5, arg__18079_7, arg__18084_10, doseq_seq__18069_11, doseq_loop__18070_12, f_13, removed_14, v170, v179, v188, doseq_seq__18069_16, doseq_loop__18070_17, f_18, removed_19, v169, v178, v187, bid_26, arg__18093_28, arg__18100_31, arg__18101_32, arg__18108_35, arg__18115_38, arg__18116_39, doseq_seq__18071_40, doseq_seq__18069_20, doseq_loop__18070_21, f_22, removed_23, v176, v185, v194, v131, doseq_seq__18069_132, doseq_loop__18070_133, f_134, removed_135, arg__18191_138, arg__18197_142, arg__18198_143, arg__18204_147, arg__18210_151, arg__18211_152, v153, v155, doseq_loop__18072_41, f_42, removed_43, bid_44, v173, v182, v191, doseq_seq__18069_46, doseq_loop__18070_47, doseq_seq__18071_48, doseq_loop__18072_49, f_50, removed_51, bid_52, v172, v181, v190, nid_62, v80, doseq_seq__18069_53, doseq_loop__18070_54, doseq_seq__18071_55, doseq_loop__18072_56, f_57, removed_58, bid_59, v175, v184, v193, v118, doseq_seq__18069_119, doseq_loop__18070_120, doseq_seq__18071_121, doseq_loop__18072_122, f_123, removed_124, bid_125, v168, v177, v186, v127, doseq_seq__18069_63, doseq_loop__18070_64, doseq_seq__18071_65, doseq_loop__18072_66, f_67, removed_68, bid_69, nid_70, v174, v183, v192, arg__18133_85, arg__18143_90, arg__18144_91, arg__18154_96, arg__18164_101, arg__18165_102, v103, v105, v109, v111, doseq_seq__18069_71, doseq_loop__18070_72, doseq_seq__18071_73, doseq_loop__18072_74, f_75, removed_76, bid_77, nid_78, v171, v180, v189, v114
	removed_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "atom").Deref(), []vm.Value{vm.Boolean(false)})
	if callErr != nil {
		return nil, callErr
	}
	arg__18079_7, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__18084_10, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	doseq_seq__18069_11, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{arg__18084_10})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__18070_12 = doseq_seq__18069_11
	f_13 = arg0
	removed_14 = removed_5
	v170 = "DEBUG: removing inst "
	v179 = " op="
	v188 = true
	goto b1
b1:
	;
	if vm.IsTruthy(doseq_loop__18070_12) {
		doseq_seq__18069_16 = doseq_seq__18069_11
		doseq_loop__18070_17 = doseq_loop__18070_12
		f_18 = f_13
		removed_19 = removed_14
		v169 = v170
		v178 = v179
		v187 = v188
		goto b2
	} else {
		doseq_seq__18069_20 = doseq_seq__18069_11
		doseq_loop__18070_21 = doseq_loop__18070_12
		f_22 = f_13
		removed_23 = removed_14
		v176 = v170
		v185 = v179
		v194 = v188
		goto b3
	}
b2:
	;
	bid_26, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{doseq_loop__18070_17})
	if callErr != nil {
		return nil, callErr
	}
	arg__18093_28, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-insts").Deref(), []vm.Value{bid_26, f_18})
	if callErr != nil {
		return nil, callErr
	}
	arg__18100_31, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-insts").Deref(), []vm.Value{bid_26, f_18})
	if callErr != nil {
		return nil, callErr
	}
	arg__18101_32, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg__18100_31})
	if callErr != nil {
		return nil, callErr
	}
	arg__18108_35, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-insts").Deref(), []vm.Value{bid_26, f_18})
	if callErr != nil {
		return nil, callErr
	}
	arg__18115_38, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-insts").Deref(), []vm.Value{bid_26, f_18})
	if callErr != nil {
		return nil, callErr
	}
	arg__18116_39, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg__18115_38})
	if callErr != nil {
		return nil, callErr
	}
	doseq_seq__18071_40, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{arg__18116_39})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__18072_41 = doseq_seq__18071_40
	f_42 = f_18
	removed_43 = removed_19
	bid_44 = bid_26
	v173 = v169
	v182 = v178
	v191 = v187
	goto b5
b3:
	;
	v131 = vm.NIL
	doseq_seq__18069_132 = doseq_seq__18069_20
	doseq_loop__18070_133 = doseq_loop__18070_21
	f_134 = f_22
	removed_135 = removed_23
	goto b4
b4:
	;
	arg__18191_138, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{removed_135})
	if callErr != nil {
		return nil, callErr
	}
	arg__18197_142, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{removed_135})
	if callErr != nil {
		return nil, callErr
	}
	arg__18198_143, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("DEBUG: one-pass removed="), arg__18197_142})
	if callErr != nil {
		return nil, callErr
	}
	arg__18204_147, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{removed_135})
	if callErr != nil {
		return nil, callErr
	}
	arg__18210_151, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{removed_135})
	if callErr != nil {
		return nil, callErr
	}
	arg__18211_152, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String("DEBUG: one-pass removed="), arg__18210_151})
	if callErr != nil {
		return nil, callErr
	}
	v153, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "println").Deref(), []vm.Value{arg__18211_152})
	if callErr != nil {
		return nil, callErr
	}
	v155, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{removed_135})
	if callErr != nil {
		return nil, callErr
	}
	return v155, nil
b5:
	;
	if vm.IsTruthy(doseq_loop__18072_41) {
		doseq_seq__18069_46 = doseq_seq__18069_16
		doseq_loop__18070_47 = doseq_loop__18070_17
		doseq_seq__18071_48 = doseq_seq__18071_40
		doseq_loop__18072_49 = doseq_loop__18072_41
		f_50 = f_42
		removed_51 = removed_43
		bid_52 = bid_44
		v172 = v173
		v181 = v182
		v190 = v191
		goto b6
	} else {
		doseq_seq__18069_53 = doseq_seq__18069_16
		doseq_loop__18070_54 = doseq_loop__18070_17
		doseq_seq__18071_55 = doseq_seq__18071_40
		doseq_loop__18072_56 = doseq_loop__18072_41
		f_57 = f_42
		removed_58 = removed_43
		bid_59 = bid_44
		v175 = v173
		v184 = v182
		v193 = v191
		goto b7
	}
b6:
	;
	nid_62, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{doseq_loop__18072_49})
	if callErr != nil {
		return nil, callErr
	}
	v80, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.dce", "removable?").Deref(), []vm.Value{nid_62, f_50})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v80) {
		doseq_seq__18069_63 = doseq_seq__18069_46
		doseq_loop__18070_64 = doseq_loop__18070_47
		doseq_seq__18071_65 = doseq_seq__18071_48
		doseq_loop__18072_66 = doseq_loop__18072_49
		f_67 = f_50
		removed_68 = removed_51
		bid_69 = bid_52
		nid_70 = nid_62
		v174 = v172
		v183 = v181
		v192 = v190
		goto b9
	} else {
		doseq_seq__18069_71 = doseq_seq__18069_46
		doseq_loop__18070_72 = doseq_loop__18070_47
		doseq_seq__18071_73 = doseq_seq__18071_48
		doseq_loop__18072_74 = doseq_loop__18072_49
		f_75 = f_50
		removed_76 = removed_51
		bid_77 = bid_52
		nid_78 = nid_62
		v171 = v172
		v180 = v181
		v189 = v190
		goto b10
	}
b7:
	;
	v118 = vm.NIL
	doseq_seq__18069_119 = doseq_seq__18069_53
	doseq_loop__18070_120 = doseq_loop__18070_54
	doseq_seq__18071_121 = doseq_seq__18071_55
	doseq_loop__18072_122 = doseq_loop__18072_56
	f_123 = f_57
	removed_124 = removed_58
	bid_125 = bid_59
	v168 = v175
	v177 = v184
	v186 = v193
	goto b8
b8:
	;
	v127, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "next").Deref(), []vm.Value{doseq_loop__18070_120})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__18070_12 = v127
	f_13 = f_123
	removed_14 = removed_124
	v170 = v168
	v179 = v177
	v188 = v186
	goto b1
b9:
	;
	arg__18133_85, callErr = rt.InvokeValue(rt.LookupVar("ir", "op").Deref(), []vm.Value{nid_70, f_67})
	if callErr != nil {
		return nil, callErr
	}
	arg__18143_90, callErr = rt.InvokeValue(rt.LookupVar("ir", "op").Deref(), []vm.Value{nid_70, f_67})
	if callErr != nil {
		return nil, callErr
	}
	arg__18144_91, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(v174), nid_70, vm.String(v183), arg__18143_90})
	if callErr != nil {
		return nil, callErr
	}
	arg__18154_96, callErr = rt.InvokeValue(rt.LookupVar("ir", "op").Deref(), []vm.Value{nid_70, f_67})
	if callErr != nil {
		return nil, callErr
	}
	arg__18164_101, callErr = rt.InvokeValue(rt.LookupVar("ir", "op").Deref(), []vm.Value{nid_70, f_67})
	if callErr != nil {
		return nil, callErr
	}
	arg__18165_102, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{vm.String(v174), nid_70, vm.String(v183), arg__18164_101})
	if callErr != nil {
		return nil, callErr
	}
	v103, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "println").Deref(), []vm.Value{arg__18165_102})
	if callErr != nil {
		return nil, callErr
	}
	v105, callErr = rt.InvokeValue(rt.LookupVar("ir", "remove-inst!").Deref(), []vm.Value{f_67, bid_69, nid_70})
	if callErr != nil {
		return nil, callErr
	}
	v109, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "reset!").Deref(), []vm.Value{removed_68, vm.Boolean(v192)})
	if callErr != nil {
		return nil, callErr
	}
	v111, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "next").Deref(), []vm.Value{doseq_loop__18072_66})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__18072_41 = v111
	f_42 = f_67
	removed_43 = removed_68
	bid_44 = bid_69
	v173 = v174
	v182 = v183
	v191 = v192
	goto b5
b10:
	;
	v114, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "next").Deref(), []vm.Value{doseq_loop__18072_74})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__18072_41 = v114
	f_42 = f_75
	removed_43 = removed_76
	bid_44 = bid_77
	v173 = v171
	v182 = v180
	v191 = v189
	goto b5
}
func dce(arg0 vm.Value) (vm.Value, error) {
	var f_2 vm.Value
	var v7 vm.Value
	var f_4 vm.Value
	var f_5 vm.Value
	var v12 vm.Value
	var f_13 vm.Value
	var callErr error
	_, _, _, _, _, _ = f_2, v7, f_4, f_5, v12, f_13
	f_2 = arg0
	goto b1
b1:
	;
	v7, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.dce", "one-pass").Deref(), []vm.Value{f_2})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v7) {
		f_4 = f_2
		goto b2
	} else {
		f_5 = f_2
		goto b3
	}
b2:
	;
	f_2 = f_4
	goto b1
b3:
	;
	v12 = vm.NIL
	f_13 = f_5
	goto b4
b4:
	;
	return f_13, nil
}
