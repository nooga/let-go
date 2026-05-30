package ir_passes_typeinfer

import (
	rt "github.com/nooga/let-go/pkg/rt"
	vm "github.com/nooga/let-go/pkg/vm"
)

func const_type_QMARK_(arg0 vm.Value) (vm.Value, error) {
	var and__x_2 vm.Value
	var t_3 vm.Value
	var and__x_4 vm.Value
	var arg__25261_10 vm.Value
	var v11 bool
	var t_5 vm.Value
	var and__x_6 vm.Value
	var v14 vm.Value
	var t_15 vm.Value
	var and__x_16 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _ = and__x_2, t_3, and__x_4, arg__25261_10, v11, t_5, and__x_6, v14, t_15, and__x_16
	and__x_2, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector?").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(and__x_2) {
		t_3 = arg0
		and__x_4 = and__x_2
		goto b1
	} else {
		t_5 = arg0
		and__x_6 = and__x_2
		goto b2
	}
b1:
	;
	arg__25261_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{t_3})
	if callErr != nil {
		return nil, callErr
	}
	v11 = arg__25261_10 == vm.Keyword("const")
	v14 = vm.Boolean(v11)
	t_15 = t_3
	and__x_16 = and__x_4
	goto b3
b2:
	;
	v14 = and__x_6
	t_15 = t_5
	and__x_16 = and__x_6
	goto b3
b3:
	;
	return v14, nil
}
func sort_members(arg0 vm.Value) (vm.Value, error) {
	var v4 vm.Value
	var callErr error
	_ = v4
	v4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "sort-by").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v4 vm.Value
		var member_1 vm.Value
		var v9 vm.Value
		var member_2 vm.Value
		var key_12 vm.Value
		var member_13 vm.Value
		var v19 vm.Value
		var callErr error
		_, _, _, _, _, _, _ = v4, member_1, v9, member_2, key_12, member_13, v19
		v4, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "const-type?").Deref(), []vm.Value{arg0})
		if callErr != nil {
			return nil, callErr
		}
		if vm.IsTruthy(v4) {
			member_1 = arg0
			goto b1
		} else {
			member_2 = arg0
			goto b2
		}
	b1:
		;
		v9, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{member_1, vm.Int(1)})
		if callErr != nil {
			return nil, callErr
		}
		key_12 = v9
		member_13 = member_1
		goto b3
	b2:
		;
		key_12 = member_2
		member_13 = member_2
		goto b3
	b3:
		;
		v19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{rt.LookupVar("ir.passes.typeinfer", "type-order").Deref(), key_12, vm.Int(100)})
		if callErr != nil {
			return nil, callErr
		}
		return v19, nil
	}), arg0})
	if callErr != nil {
		return nil, callErr
	}
	return v4, nil
}
func union_type_QMARK_(arg0 vm.Value) (vm.Value, error) {
	var and__x_2 vm.Value
	var t_3 vm.Value
	var and__x_4 vm.Value
	var arg__25304_10 vm.Value
	var v11 bool
	var t_5 vm.Value
	var and__x_6 vm.Value
	var v14 vm.Value
	var t_15 vm.Value
	var and__x_16 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _ = and__x_2, t_3, and__x_4, arg__25304_10, v11, t_5, and__x_6, v14, t_15, and__x_16
	and__x_2, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector?").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(and__x_2) {
		t_3 = arg0
		and__x_4 = and__x_2
		goto b1
	} else {
		t_5 = arg0
		and__x_6 = and__x_2
		goto b2
	}
b1:
	;
	arg__25304_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{t_3})
	if callErr != nil {
		return nil, callErr
	}
	v11 = arg__25304_10 == vm.Keyword("union")
	v14 = vm.Boolean(v11)
	t_15 = t_3
	and__x_16 = and__x_4
	goto b3
b2:
	;
	v14 = and__x_6
	t_15 = t_5
	and__x_16 = and__x_6
	goto b3
b3:
	;
	return v14, nil
}
func ti_inc_BANG_(arg0 vm.Value) (vm.Value, error) {
	var k_1 vm.Value
	var arg__25313_12 vm.Value
	var arg__25323_21 vm.Value
	var v22 vm.Value
	var k_2 vm.Value
	var v26 vm.Value
	var k_27 vm.Value
	var callErr error
	_, _, _, _, _, _, _ = k_1, arg__25313_12, arg__25323_21, v22, k_2, v26, k_27
	if vm.IsTruthy(rt.LookupVar("ir.passes.typeinfer", "*ti-counters*").Deref()) {
		k_1 = arg0
		goto b1
	} else {
		k_2 = arg0
		goto b2
	}
b1:
	;
	arg__25313_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "fnil").Deref(), []vm.Value{rt.LookupVar("clojure.core", "inc").Deref(), vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	arg__25323_21, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "fnil").Deref(), []vm.Value{rt.LookupVar("clojure.core", "inc").Deref(), vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	v22, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{rt.LookupVar("ir.passes.typeinfer", "*ti-counters*").Deref(), rt.LookupVar("clojure.core", "update").Deref(), k_1, arg__25323_21})
	if callErr != nil {
		return nil, callErr
	}
	v26 = v22
	k_27 = k_1
	goto b3
b2:
	;
	v26 = vm.NIL
	k_27 = k_2
	goto b3
b3:
	;
	return v26, nil
}
func refine_truthy(arg0 vm.Value) (vm.Value, error) {
	var t_2 vm.Value
	var v6 bool
	var t_3 vm.Value
	var t_4 vm.Value
	var v13 bool
	var v99 vm.Value
	var t_100 vm.Value
	var t_10 vm.Value
	var t_11 vm.Value
	var v20 bool
	var v96 vm.Value
	var t_97 vm.Value
	var t_17 vm.Value
	var t_18 vm.Value
	var v27 vm.Value
	var v93 vm.Value
	var t_94 vm.Value
	var t_24 vm.Value
	var arg__25752_31 vm.Value
	var arg__25761_34 vm.Value
	var arg__25771_38 vm.Value
	var arg__25772_39 vm.Value
	var arg__25776_43 vm.Value
	var arg__25785_46 vm.Value
	var arg__25795_50 vm.Value
	var arg__25796_51 vm.Value
	var arg__25797_52 vm.Value
	var arg__25801_56 vm.Value
	var arg__25810_59 vm.Value
	var arg__25820_63 vm.Value
	var arg__25821_64 vm.Value
	var arg__25825_68 vm.Value
	var arg__25834_71 vm.Value
	var arg__25844_75 vm.Value
	var arg__25845_76 vm.Value
	var arg__25846_77 vm.Value
	var v78 vm.Value
	var t_25 vm.Value
	var v90 vm.Value
	var t_91 vm.Value
	var t_80 vm.Value
	var t_81 vm.Value
	var v87 vm.Value
	var t_88 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = t_2, v6, t_3, t_4, v13, v99, t_100, t_10, t_11, v20, v96, t_97, t_17, t_18, v27, v93, t_94, t_24, arg__25752_31, arg__25761_34, arg__25771_38, arg__25772_39, arg__25776_43, arg__25785_46, arg__25795_50, arg__25796_51, arg__25797_52, arg__25801_56, arg__25810_59, arg__25820_63, arg__25821_64, arg__25825_68, arg__25834_71, arg__25844_75, arg__25845_76, arg__25846_77, v78, t_25, v90, t_91, t_80, t_81, v87, t_88
	t_2, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "normalize-type").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v6 = t_2 == vm.Keyword("bool")
	if v6 {
		t_3 = t_2
		goto b1
	} else {
		t_4 = t_2
		goto b2
	}
b1:
	;
	v99 = vm.Keyword("true")
	t_100 = t_3
	goto b3
b2:
	;
	v13 = t_4 == vm.Keyword("false")
	if v13 {
		t_10 = t_4
		goto b4
	} else {
		t_11 = t_4
		goto b5
	}
b3:
	;
	return v99, nil
b4:
	;
	v96 = vm.Keyword("bottom")
	t_97 = t_10
	goto b6
b5:
	;
	v20 = t_11 == vm.Keyword("nil")
	if v20 {
		t_17 = t_11
		goto b7
	} else {
		t_18 = t_11
		goto b8
	}
b6:
	;
	v99 = v96
	t_100 = t_97
	goto b3
b7:
	;
	v93 = vm.Keyword("bottom")
	t_94 = t_17
	goto b9
b8:
	;
	v27, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "union-type?").Deref(), []vm.Value{t_18})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v27) {
		t_24 = t_18
		goto b10
	} else {
		t_25 = t_18
		goto b11
	}
b9:
	;
	v96 = v93
	t_97 = t_94
	goto b6
b10:
	;
	arg__25752_31, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("union")})
	if callErr != nil {
		return nil, callErr
	}
	arg__25761_34, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_24})
	if callErr != nil {
		return nil, callErr
	}
	arg__25771_38, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_24})
	if callErr != nil {
		return nil, callErr
	}
	arg__25772_39, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "remove").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) vm.Value {
		var or__x_2 vm.Value
		var m_3 vm.Value
		var or__x_4 vm.Value
		var m_5 vm.Value
		var or__x_6 vm.Value
		var v10 vm.Value
		var v12 vm.Value
		var m_13 vm.Value
		var or__x_14 vm.Value
		_, _, _, _, _, _, _, _, _ = or__x_2, m_3, or__x_4, m_5, or__x_6, v10, v12, m_13, or__x_14
		or__x_2 = vm.Boolean(arg0 == vm.Keyword("nil"))
		if vm.IsTruthy(or__x_2) {
			m_3 = arg0
			or__x_4 = or__x_2
			goto b1
		} else {
			m_5 = arg0
			or__x_6 = or__x_2
			goto b2
		}
	b1:
		;
		v12 = or__x_4
		m_13 = m_3
		or__x_14 = or__x_4
		goto b3
	b2:
		;
		v10 = vm.Boolean(m_5 == vm.Keyword("false"))
		v12 = v10
		m_13 = m_5
		or__x_14 = or__x_6
		goto b3
	b3:
		;
		return v12
	}), arg__25771_38})
	if callErr != nil {
		return nil, callErr
	}
	arg__25776_43, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("union")})
	if callErr != nil {
		return nil, callErr
	}
	arg__25785_46, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_24})
	if callErr != nil {
		return nil, callErr
	}
	arg__25795_50, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_24})
	if callErr != nil {
		return nil, callErr
	}
	arg__25796_51, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "remove").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) vm.Value {
		var or__x_2 vm.Value
		var m_3 vm.Value
		var or__x_4 vm.Value
		var m_5 vm.Value
		var or__x_6 vm.Value
		var v10 vm.Value
		var v12 vm.Value
		var m_13 vm.Value
		var or__x_14 vm.Value
		_, _, _, _, _, _, _, _, _ = or__x_2, m_3, or__x_4, m_5, or__x_6, v10, v12, m_13, or__x_14
		or__x_2 = vm.Boolean(arg0 == vm.Keyword("nil"))
		if vm.IsTruthy(or__x_2) {
			m_3 = arg0
			or__x_4 = or__x_2
			goto b1
		} else {
			m_5 = arg0
			or__x_6 = or__x_2
			goto b2
		}
	b1:
		;
		v12 = or__x_4
		m_13 = m_3
		or__x_14 = or__x_4
		goto b3
	b2:
		;
		v10 = vm.Boolean(m_5 == vm.Keyword("false"))
		v12 = v10
		m_13 = m_5
		or__x_14 = or__x_6
		goto b3
	b3:
		;
		return v12
	}), arg__25795_50})
	if callErr != nil {
		return nil, callErr
	}
	arg__25797_52, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__25776_43, arg__25796_51})
	if callErr != nil {
		return nil, callErr
	}
	arg__25801_56, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("union")})
	if callErr != nil {
		return nil, callErr
	}
	arg__25810_59, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_24})
	if callErr != nil {
		return nil, callErr
	}
	arg__25820_63, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_24})
	if callErr != nil {
		return nil, callErr
	}
	arg__25821_64, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "remove").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) vm.Value {
		var or__x_2 vm.Value
		var m_3 vm.Value
		var or__x_4 vm.Value
		var m_5 vm.Value
		var or__x_6 vm.Value
		var v10 vm.Value
		var v12 vm.Value
		var m_13 vm.Value
		var or__x_14 vm.Value
		_, _, _, _, _, _, _, _, _ = or__x_2, m_3, or__x_4, m_5, or__x_6, v10, v12, m_13, or__x_14
		or__x_2 = vm.Boolean(arg0 == vm.Keyword("nil"))
		if vm.IsTruthy(or__x_2) {
			m_3 = arg0
			or__x_4 = or__x_2
			goto b1
		} else {
			m_5 = arg0
			or__x_6 = or__x_2
			goto b2
		}
	b1:
		;
		v12 = or__x_4
		m_13 = m_3
		or__x_14 = or__x_4
		goto b3
	b2:
		;
		v10 = vm.Boolean(m_5 == vm.Keyword("false"))
		v12 = v10
		m_13 = m_5
		or__x_14 = or__x_6
		goto b3
	b3:
		;
		return v12
	}), arg__25820_63})
	if callErr != nil {
		return nil, callErr
	}
	arg__25825_68, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("union")})
	if callErr != nil {
		return nil, callErr
	}
	arg__25834_71, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_24})
	if callErr != nil {
		return nil, callErr
	}
	arg__25844_75, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_24})
	if callErr != nil {
		return nil, callErr
	}
	arg__25845_76, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "remove").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) vm.Value {
		var or__x_2 vm.Value
		var m_3 vm.Value
		var or__x_4 vm.Value
		var m_5 vm.Value
		var or__x_6 vm.Value
		var v10 vm.Value
		var v12 vm.Value
		var m_13 vm.Value
		var or__x_14 vm.Value
		_, _, _, _, _, _, _, _, _ = or__x_2, m_3, or__x_4, m_5, or__x_6, v10, v12, m_13, or__x_14
		or__x_2 = vm.Boolean(arg0 == vm.Keyword("nil"))
		if vm.IsTruthy(or__x_2) {
			m_3 = arg0
			or__x_4 = or__x_2
			goto b1
		} else {
			m_5 = arg0
			or__x_6 = or__x_2
			goto b2
		}
	b1:
		;
		v12 = or__x_4
		m_13 = m_3
		or__x_14 = or__x_4
		goto b3
	b2:
		;
		v10 = vm.Boolean(m_5 == vm.Keyword("false"))
		v12 = v10
		m_13 = m_5
		or__x_14 = or__x_6
		goto b3
	b3:
		;
		return v12
	}), arg__25844_75})
	if callErr != nil {
		return nil, callErr
	}
	arg__25846_77, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__25825_68, arg__25845_76})
	if callErr != nil {
		return nil, callErr
	}
	v78, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "normalize-type").Deref(), []vm.Value{arg__25846_77})
	if callErr != nil {
		return nil, callErr
	}
	v90 = v78
	t_91 = t_24
	goto b12
b11:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		t_80 = t_25
		goto b13
	} else {
		t_81 = t_25
		goto b14
	}
b12:
	;
	v93 = v90
	t_94 = t_91
	goto b9
b13:
	;
	v87 = t_80
	t_88 = t_80
	goto b15
b14:
	;
	v87 = vm.NIL
	t_88 = t_81
	goto b15
b15:
	;
	v90 = v87
	t_91 = t_88
	goto b12
}
func truthy_type_QMARK_(arg0 vm.Value) (vm.Value, error) {
	var t_2 vm.Value
	var v6 bool
	var t_3 vm.Value
	var t_4 vm.Value
	var v13 bool
	var v126 vm.Value
	var t_127 vm.Value
	var t_10 vm.Value
	var t_11 vm.Value
	var v20 bool
	var v123 vm.Value
	var t_124 vm.Value
	var t_17 vm.Value
	var t_18 vm.Value
	var v27 bool
	var v120 vm.Value
	var t_121 vm.Value
	var t_24 vm.Value
	var t_25 vm.Value
	var v34 bool
	var v117 vm.Value
	var t_118 vm.Value
	var t_31 vm.Value
	var t_32 vm.Value
	var v41 bool
	var v114 vm.Value
	var t_115 vm.Value
	var t_38 vm.Value
	var t_39 vm.Value
	var v48 bool
	var v111 vm.Value
	var t_112 vm.Value
	var t_45 vm.Value
	var t_46 vm.Value
	var arg__25869_58 vm.Value
	var v59 bool
	var v108 vm.Value
	var t_109 vm.Value
	var t_52 vm.Value
	var t_53 vm.Value
	var arg__25875_69 vm.Value
	var v70 bool
	var v105 vm.Value
	var t_106 vm.Value
	var t_63 vm.Value
	var t_64 vm.Value
	var v77 vm.Value
	var v102 vm.Value
	var t_103 vm.Value
	var t_74 vm.Value
	var arg__25883_81 vm.Value
	var arg__25889_85 vm.Value
	var v86 vm.Value
	var t_75 vm.Value
	var v99 vm.Value
	var t_100 vm.Value
	var t_88 vm.Value
	var t_89 vm.Value
	var v96 vm.Value
	var t_97 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = t_2, v6, t_3, t_4, v13, v126, t_127, t_10, t_11, v20, v123, t_124, t_17, t_18, v27, v120, t_121, t_24, t_25, v34, v117, t_118, t_31, t_32, v41, v114, t_115, t_38, t_39, v48, v111, t_112, t_45, t_46, arg__25869_58, v59, v108, t_109, t_52, t_53, arg__25875_69, v70, v105, t_106, t_63, t_64, v77, v102, t_103, t_74, arg__25883_81, arg__25889_85, v86, t_75, v99, t_100, t_88, t_89, v96, t_97
	t_2, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "normalize-type").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v6 = t_2 == vm.Keyword("true")
	if v6 {
		t_3 = t_2
		goto b1
	} else {
		t_4 = t_2
		goto b2
	}
b1:
	;
	v126 = vm.Boolean(true)
	t_127 = t_3
	goto b3
b2:
	;
	v13 = t_4 == vm.Keyword("false")
	if v13 {
		t_10 = t_4
		goto b4
	} else {
		t_11 = t_4
		goto b5
	}
b3:
	;
	return v126, nil
b4:
	;
	v123 = vm.Boolean(false)
	t_124 = t_10
	goto b6
b5:
	;
	v20 = t_11 == vm.Keyword("nil")
	if v20 {
		t_17 = t_11
		goto b7
	} else {
		t_18 = t_11
		goto b8
	}
b6:
	;
	v126 = v123
	t_127 = t_124
	goto b3
b7:
	;
	v120 = vm.Boolean(false)
	t_121 = t_17
	goto b9
b8:
	;
	v27 = t_18 == vm.Keyword("int")
	if v27 {
		t_24 = t_18
		goto b10
	} else {
		t_25 = t_18
		goto b11
	}
b9:
	;
	v123 = v120
	t_124 = t_121
	goto b6
b10:
	;
	v117 = vm.Boolean(true)
	t_118 = t_24
	goto b12
b11:
	;
	v34 = t_25 == vm.Keyword("float")
	if v34 {
		t_31 = t_25
		goto b13
	} else {
		t_32 = t_25
		goto b14
	}
b12:
	;
	v120 = v117
	t_121 = t_118
	goto b9
b13:
	;
	v114 = vm.Boolean(true)
	t_115 = t_31
	goto b15
b14:
	;
	v41 = t_32 == vm.Keyword("number")
	if v41 {
		t_38 = t_32
		goto b16
	} else {
		t_39 = t_32
		goto b17
	}
b15:
	;
	v117 = v114
	t_118 = t_115
	goto b12
b16:
	;
	v111 = vm.Boolean(true)
	t_112 = t_38
	goto b18
b17:
	;
	v48 = t_39 == vm.Keyword("string")
	if v48 {
		t_45 = t_39
		goto b19
	} else {
		t_46 = t_39
		goto b20
	}
b18:
	;
	v114 = v111
	t_115 = t_112
	goto b15
b19:
	;
	v108 = vm.Boolean(true)
	t_109 = t_45
	goto b21
b20:
	;
	arg__25869_58, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("const"), vm.Keyword("int"), vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	v59 = t_46 == arg__25869_58
	if v59 {
		t_52 = t_46
		goto b22
	} else {
		t_53 = t_46
		goto b23
	}
b21:
	;
	v111 = v108
	t_112 = t_109
	goto b18
b22:
	;
	v105 = vm.Boolean(true)
	t_106 = t_52
	goto b24
b23:
	;
	arg__25875_69, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("const"), vm.Keyword("float"), vm.Float(0)})
	if callErr != nil {
		return nil, callErr
	}
	v70 = t_53 == arg__25875_69
	if v70 {
		t_63 = t_53
		goto b25
	} else {
		t_64 = t_53
		goto b26
	}
b24:
	;
	v108 = v105
	t_109 = t_106
	goto b21
b25:
	;
	v102 = vm.Boolean(true)
	t_103 = t_63
	goto b27
b26:
	;
	v77, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "union-type?").Deref(), []vm.Value{t_64})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v77) {
		t_74 = t_64
		goto b28
	} else {
		t_75 = t_64
		goto b29
	}
b27:
	;
	v105 = v102
	t_106 = t_103
	goto b24
b28:
	;
	arg__25883_81, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_74})
	if callErr != nil {
		return nil, callErr
	}
	arg__25889_85, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_74})
	if callErr != nil {
		return nil, callErr
	}
	v86, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "every?").Deref(), []vm.Value{rt.LookupVar("ir.passes.typeinfer", "truthy-type?").Deref(), arg__25889_85})
	if callErr != nil {
		return nil, callErr
	}
	v99 = v86
	t_100 = t_74
	goto b30
b29:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		t_88 = t_75
		goto b31
	} else {
		t_89 = t_75
		goto b32
	}
b30:
	;
	v102 = v99
	t_103 = t_100
	goto b27
b31:
	;
	v96 = vm.Boolean(false)
	t_97 = t_88
	goto b33
b32:
	;
	v96 = vm.NIL
	t_97 = t_89
	goto b33
b33:
	;
	v99 = v96
	t_100 = t_97
	goto b30
}
func refine_falsey(arg0 vm.Value) (vm.Value, error) {
	var t_2 vm.Value
	var v6 bool
	var t_3 vm.Value
	var t_4 vm.Value
	var v13 bool
	var v99 vm.Value
	var t_100 vm.Value
	var t_10 vm.Value
	var t_11 vm.Value
	var v20 vm.Value
	var v96 vm.Value
	var t_97 vm.Value
	var t_17 vm.Value
	var t_18 vm.Value
	var v27 vm.Value
	var v93 vm.Value
	var t_94 vm.Value
	var t_24 vm.Value
	var arg__25905_31 vm.Value
	var arg__25914_34 vm.Value
	var arg__25924_38 vm.Value
	var arg__25925_39 vm.Value
	var arg__25929_43 vm.Value
	var arg__25938_46 vm.Value
	var arg__25948_50 vm.Value
	var arg__25949_51 vm.Value
	var arg__25950_52 vm.Value
	var arg__25954_56 vm.Value
	var arg__25963_59 vm.Value
	var arg__25973_63 vm.Value
	var arg__25974_64 vm.Value
	var arg__25978_68 vm.Value
	var arg__25987_71 vm.Value
	var arg__25997_75 vm.Value
	var arg__25998_76 vm.Value
	var arg__25999_77 vm.Value
	var v78 vm.Value
	var t_25 vm.Value
	var v90 vm.Value
	var t_91 vm.Value
	var t_80 vm.Value
	var t_81 vm.Value
	var v87 vm.Value
	var t_88 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = t_2, v6, t_3, t_4, v13, v99, t_100, t_10, t_11, v20, v96, t_97, t_17, t_18, v27, v93, t_94, t_24, arg__25905_31, arg__25914_34, arg__25924_38, arg__25925_39, arg__25929_43, arg__25938_46, arg__25948_50, arg__25949_51, arg__25950_52, arg__25954_56, arg__25963_59, arg__25973_63, arg__25974_64, arg__25978_68, arg__25987_71, arg__25997_75, arg__25998_76, arg__25999_77, v78, t_25, v90, t_91, t_80, t_81, v87, t_88
	t_2, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "normalize-type").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v6 = t_2 == vm.Keyword("bool")
	if v6 {
		t_3 = t_2
		goto b1
	} else {
		t_4 = t_2
		goto b2
	}
b1:
	;
	v99 = vm.Keyword("false")
	t_100 = t_3
	goto b3
b2:
	;
	v13 = t_4 == vm.Keyword("true")
	if v13 {
		t_10 = t_4
		goto b4
	} else {
		t_11 = t_4
		goto b5
	}
b3:
	;
	return v99, nil
b4:
	;
	v96 = vm.Keyword("bottom")
	t_97 = t_10
	goto b6
b5:
	;
	v20, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "truthy-type?").Deref(), []vm.Value{t_11})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v20) {
		t_17 = t_11
		goto b7
	} else {
		t_18 = t_11
		goto b8
	}
b6:
	;
	v99 = v96
	t_100 = t_97
	goto b3
b7:
	;
	v93 = vm.Keyword("bottom")
	t_94 = t_17
	goto b9
b8:
	;
	v27, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "union-type?").Deref(), []vm.Value{t_18})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v27) {
		t_24 = t_18
		goto b10
	} else {
		t_25 = t_18
		goto b11
	}
b9:
	;
	v96 = v93
	t_97 = t_94
	goto b6
b10:
	;
	arg__25905_31, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("union")})
	if callErr != nil {
		return nil, callErr
	}
	arg__25914_34, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_24})
	if callErr != nil {
		return nil, callErr
	}
	arg__25924_38, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_24})
	if callErr != nil {
		return nil, callErr
	}
	arg__25925_39, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "filter").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) vm.Value {
		var or__x_2 vm.Value
		var m_3 vm.Value
		var or__x_4 vm.Value
		var m_5 vm.Value
		var or__x_6 vm.Value
		var v10 vm.Value
		var v12 vm.Value
		var m_13 vm.Value
		var or__x_14 vm.Value
		_, _, _, _, _, _, _, _, _ = or__x_2, m_3, or__x_4, m_5, or__x_6, v10, v12, m_13, or__x_14
		or__x_2 = vm.Boolean(arg0 == vm.Keyword("nil"))
		if vm.IsTruthy(or__x_2) {
			m_3 = arg0
			or__x_4 = or__x_2
			goto b1
		} else {
			m_5 = arg0
			or__x_6 = or__x_2
			goto b2
		}
	b1:
		;
		v12 = or__x_4
		m_13 = m_3
		or__x_14 = or__x_4
		goto b3
	b2:
		;
		v10 = vm.Boolean(m_5 == vm.Keyword("false"))
		v12 = v10
		m_13 = m_5
		or__x_14 = or__x_6
		goto b3
	b3:
		;
		return v12
	}), arg__25924_38})
	if callErr != nil {
		return nil, callErr
	}
	arg__25929_43, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("union")})
	if callErr != nil {
		return nil, callErr
	}
	arg__25938_46, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_24})
	if callErr != nil {
		return nil, callErr
	}
	arg__25948_50, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_24})
	if callErr != nil {
		return nil, callErr
	}
	arg__25949_51, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "filter").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) vm.Value {
		var or__x_2 vm.Value
		var m_3 vm.Value
		var or__x_4 vm.Value
		var m_5 vm.Value
		var or__x_6 vm.Value
		var v10 vm.Value
		var v12 vm.Value
		var m_13 vm.Value
		var or__x_14 vm.Value
		_, _, _, _, _, _, _, _, _ = or__x_2, m_3, or__x_4, m_5, or__x_6, v10, v12, m_13, or__x_14
		or__x_2 = vm.Boolean(arg0 == vm.Keyword("nil"))
		if vm.IsTruthy(or__x_2) {
			m_3 = arg0
			or__x_4 = or__x_2
			goto b1
		} else {
			m_5 = arg0
			or__x_6 = or__x_2
			goto b2
		}
	b1:
		;
		v12 = or__x_4
		m_13 = m_3
		or__x_14 = or__x_4
		goto b3
	b2:
		;
		v10 = vm.Boolean(m_5 == vm.Keyword("false"))
		v12 = v10
		m_13 = m_5
		or__x_14 = or__x_6
		goto b3
	b3:
		;
		return v12
	}), arg__25948_50})
	if callErr != nil {
		return nil, callErr
	}
	arg__25950_52, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__25929_43, arg__25949_51})
	if callErr != nil {
		return nil, callErr
	}
	arg__25954_56, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("union")})
	if callErr != nil {
		return nil, callErr
	}
	arg__25963_59, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_24})
	if callErr != nil {
		return nil, callErr
	}
	arg__25973_63, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_24})
	if callErr != nil {
		return nil, callErr
	}
	arg__25974_64, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "filter").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) vm.Value {
		var or__x_2 vm.Value
		var m_3 vm.Value
		var or__x_4 vm.Value
		var m_5 vm.Value
		var or__x_6 vm.Value
		var v10 vm.Value
		var v12 vm.Value
		var m_13 vm.Value
		var or__x_14 vm.Value
		_, _, _, _, _, _, _, _, _ = or__x_2, m_3, or__x_4, m_5, or__x_6, v10, v12, m_13, or__x_14
		or__x_2 = vm.Boolean(arg0 == vm.Keyword("nil"))
		if vm.IsTruthy(or__x_2) {
			m_3 = arg0
			or__x_4 = or__x_2
			goto b1
		} else {
			m_5 = arg0
			or__x_6 = or__x_2
			goto b2
		}
	b1:
		;
		v12 = or__x_4
		m_13 = m_3
		or__x_14 = or__x_4
		goto b3
	b2:
		;
		v10 = vm.Boolean(m_5 == vm.Keyword("false"))
		v12 = v10
		m_13 = m_5
		or__x_14 = or__x_6
		goto b3
	b3:
		;
		return v12
	}), arg__25973_63})
	if callErr != nil {
		return nil, callErr
	}
	arg__25978_68, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("union")})
	if callErr != nil {
		return nil, callErr
	}
	arg__25987_71, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_24})
	if callErr != nil {
		return nil, callErr
	}
	arg__25997_75, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_24})
	if callErr != nil {
		return nil, callErr
	}
	arg__25998_76, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "filter").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) vm.Value {
		var or__x_2 vm.Value
		var m_3 vm.Value
		var or__x_4 vm.Value
		var m_5 vm.Value
		var or__x_6 vm.Value
		var v10 vm.Value
		var v12 vm.Value
		var m_13 vm.Value
		var or__x_14 vm.Value
		_, _, _, _, _, _, _, _, _ = or__x_2, m_3, or__x_4, m_5, or__x_6, v10, v12, m_13, or__x_14
		or__x_2 = vm.Boolean(arg0 == vm.Keyword("nil"))
		if vm.IsTruthy(or__x_2) {
			m_3 = arg0
			or__x_4 = or__x_2
			goto b1
		} else {
			m_5 = arg0
			or__x_6 = or__x_2
			goto b2
		}
	b1:
		;
		v12 = or__x_4
		m_13 = m_3
		or__x_14 = or__x_4
		goto b3
	b2:
		;
		v10 = vm.Boolean(m_5 == vm.Keyword("false"))
		v12 = v10
		m_13 = m_5
		or__x_14 = or__x_6
		goto b3
	b3:
		;
		return v12
	}), arg__25997_75})
	if callErr != nil {
		return nil, callErr
	}
	arg__25999_77, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "into").Deref(), []vm.Value{arg__25978_68, arg__25998_76})
	if callErr != nil {
		return nil, callErr
	}
	v78, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "normalize-type").Deref(), []vm.Value{arg__25999_77})
	if callErr != nil {
		return nil, callErr
	}
	v90 = v78
	t_91 = t_24
	goto b12
b11:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		t_80 = t_25
		goto b13
	} else {
		t_81 = t_25
		goto b14
	}
b12:
	;
	v93 = v90
	t_94 = t_91
	goto b9
b13:
	;
	v87 = t_80
	t_88 = t_80
	goto b15
b14:
	;
	v87 = vm.NIL
	t_88 = t_81
	goto b15
b15:
	;
	v90 = v87
	t_91 = t_88
	goto b12
}
func refine_edge_arg_type(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value, arg3 vm.Value, arg4 vm.Value) (vm.Value, error) {
	var term_6 vm.Value
	var v20 vm.Value
	var pred_bid_7 vm.Value
	var target_bid_8 vm.Value
	var arg_id_9 vm.Value
	var arg_type_10 vm.Value
	var f_11 vm.Value
	var term_12 vm.Value
	var pred_bid_13 vm.Value
	var target_bid_14 vm.Value
	var arg_id_15 vm.Value
	var arg_type_16 vm.Value
	var f_17 vm.Value
	var term_18 vm.Value
	var op_24 vm.Value
	var refs_26 vm.Value
	var aux_28 vm.Value
	var or__x_50 vm.Value
	var v304 vm.Value
	var pred_bid_305 vm.Value
	var target_bid_306 vm.Value
	var arg_id_307 vm.Value
	var arg_type_308 vm.Value
	var f_309 vm.Value
	var term_310 vm.Value
	var pred_bid_29 vm.Value
	var target_bid_30 vm.Value
	var arg_id_31 vm.Value
	var arg_type_32 vm.Value
	var f_33 vm.Value
	var term_34 vm.Value
	var op_35 vm.Value
	var refs_36 vm.Value
	var aux_37 vm.Value
	var pred_bid_38 vm.Value
	var target_bid_39 vm.Value
	var arg_id_40 vm.Value
	var arg_type_41 vm.Value
	var f_42 vm.Value
	var term_43 vm.Value
	var op_44 vm.Value
	var refs_45 vm.Value
	var aux_46 vm.Value
	var cond_id_90 vm.Value
	var tt_92 vm.Value
	var ft_94 vm.Value
	var arg__26044_96 vm.Value
	var true_edge_QMARK__97 bool
	var arg__26049_99 vm.Value
	var false_edge_QMARK__100 bool
	var v129 bool
	var v293 vm.Value
	var pred_bid_294 vm.Value
	var target_bid_295 vm.Value
	var arg_id_296 vm.Value
	var arg_type_297 vm.Value
	var f_298 vm.Value
	var term_299 vm.Value
	var op_300 vm.Value
	var refs_301 vm.Value
	var aux_302 vm.Value
	var pred_bid_51 vm.Value
	var target_bid_52 vm.Value
	var arg_id_53 vm.Value
	var arg_type_54 vm.Value
	var f_55 vm.Value
	var term_56 vm.Value
	var op_57 vm.Value
	var refs_58 vm.Value
	var aux_59 vm.Value
	var or__x_60 vm.Value
	var pred_bid_61 vm.Value
	var target_bid_62 vm.Value
	var arg_id_63 vm.Value
	var arg_type_64 vm.Value
	var f_65 vm.Value
	var term_66 vm.Value
	var op_67 vm.Value
	var refs_68 vm.Value
	var aux_69 vm.Value
	var or__x_70 vm.Value
	var v74 vm.Value
	var v76 vm.Value
	var pred_bid_77 vm.Value
	var target_bid_78 vm.Value
	var arg_id_79 vm.Value
	var arg_type_80 vm.Value
	var f_81 vm.Value
	var term_82 vm.Value
	var op_83 vm.Value
	var refs_84 vm.Value
	var aux_85 vm.Value
	var or__x_86 vm.Value
	var pred_bid_101 vm.Value
	var target_bid_102 vm.Value
	var arg_id_103 vm.Value
	var arg_type_104 vm.Value
	var f_105 vm.Value
	var term_106 vm.Value
	var op_107 vm.Value
	var refs_108 vm.Value
	var aux_109 vm.Value
	var cond_id_110 vm.Value
	var tt_111 vm.Value
	var ft_112 vm.Value
	var true_edge_QMARK__113 bool
	var false_edge_QMARK__114 bool
	var pred_bid_115 vm.Value
	var target_bid_116 vm.Value
	var arg_id_117 vm.Value
	var arg_type_118 vm.Value
	var f_119 vm.Value
	var term_120 vm.Value
	var op_121 vm.Value
	var refs_122 vm.Value
	var aux_123 vm.Value
	var cond_id_124 vm.Value
	var tt_125 vm.Value
	var ft_126 vm.Value
	var true_edge_QMARK__127 bool
	var false_edge_QMARK__128 bool
	var v277 vm.Value
	var pred_bid_278 vm.Value
	var target_bid_279 vm.Value
	var arg_id_280 vm.Value
	var arg_type_281 vm.Value
	var f_282 vm.Value
	var term_283 vm.Value
	var op_284 vm.Value
	var refs_285 vm.Value
	var aux_286 vm.Value
	var cond_id_287 vm.Value
	var tt_288 vm.Value
	var ft_289 vm.Value
	var true_edge_QMARK__290 bool
	var false_edge_QMARK__291 bool
	var pred_bid_131 vm.Value
	var target_bid_132 vm.Value
	var arg_id_133 vm.Value
	var arg_type_134 vm.Value
	var f_135 vm.Value
	var term_136 vm.Value
	var op_137 vm.Value
	var refs_138 vm.Value
	var aux_139 vm.Value
	var cond_id_140 vm.Value
	var tt_141 vm.Value
	var ft_142 vm.Value
	var true_edge_QMARK__143 bool
	var false_edge_QMARK__144 bool
	var v161 vm.Value
	var pred_bid_145 vm.Value
	var target_bid_146 vm.Value
	var arg_id_147 vm.Value
	var arg_type_148 vm.Value
	var f_149 vm.Value
	var term_150 vm.Value
	var op_151 vm.Value
	var refs_152 vm.Value
	var aux_153 vm.Value
	var cond_id_154 vm.Value
	var tt_155 vm.Value
	var ft_156 vm.Value
	var true_edge_QMARK__157 bool
	var false_edge_QMARK__158 bool
	var v260 vm.Value
	var pred_bid_261 vm.Value
	var target_bid_262 vm.Value
	var arg_id_263 vm.Value
	var arg_type_264 vm.Value
	var f_265 vm.Value
	var term_266 vm.Value
	var op_267 vm.Value
	var refs_268 vm.Value
	var aux_269 vm.Value
	var cond_id_270 vm.Value
	var tt_271 vm.Value
	var ft_272 vm.Value
	var true_edge_QMARK__273 vm.Value
	var false_edge_QMARK__274 bool
	var pred_bid_163 vm.Value
	var target_bid_164 vm.Value
	var arg_id_165 vm.Value
	var arg_type_166 vm.Value
	var f_167 vm.Value
	var term_168 vm.Value
	var op_169 vm.Value
	var refs_170 vm.Value
	var aux_171 vm.Value
	var cond_id_172 vm.Value
	var tt_173 vm.Value
	var ft_174 vm.Value
	var true_edge_QMARK__175 bool
	var false_edge_QMARK__176 bool
	var v193 vm.Value
	var pred_bid_177 vm.Value
	var target_bid_178 vm.Value
	var arg_id_179 vm.Value
	var arg_type_180 vm.Value
	var f_181 vm.Value
	var term_182 vm.Value
	var op_183 vm.Value
	var refs_184 vm.Value
	var aux_185 vm.Value
	var cond_id_186 vm.Value
	var tt_187 vm.Value
	var ft_188 vm.Value
	var true_edge_QMARK__189 bool
	var false_edge_QMARK__190 bool
	var v244 vm.Value
	var pred_bid_245 vm.Value
	var target_bid_246 vm.Value
	var arg_id_247 vm.Value
	var arg_type_248 vm.Value
	var f_249 vm.Value
	var term_250 vm.Value
	var op_251 vm.Value
	var refs_252 vm.Value
	var aux_253 vm.Value
	var cond_id_254 vm.Value
	var tt_255 vm.Value
	var ft_256 vm.Value
	var true_edge_QMARK__257 bool
	var false_edge_QMARK__258 vm.Value
	var pred_bid_195 vm.Value
	var target_bid_196 vm.Value
	var arg_id_197 vm.Value
	var arg_type_198 vm.Value
	var f_199 vm.Value
	var term_200 vm.Value
	var op_201 vm.Value
	var refs_202 vm.Value
	var aux_203 vm.Value
	var cond_id_204 vm.Value
	var tt_205 vm.Value
	var ft_206 vm.Value
	var true_edge_QMARK__207 bool
	var false_edge_QMARK__208 bool
	var pred_bid_209 vm.Value
	var target_bid_210 vm.Value
	var arg_id_211 vm.Value
	var arg_type_212 vm.Value
	var f_213 vm.Value
	var term_214 vm.Value
	var op_215 vm.Value
	var refs_216 vm.Value
	var aux_217 vm.Value
	var cond_id_218 vm.Value
	var tt_219 vm.Value
	var ft_220 vm.Value
	var true_edge_QMARK__221 bool
	var false_edge_QMARK__222 bool
	var v228 vm.Value
	var pred_bid_229 vm.Value
	var target_bid_230 vm.Value
	var arg_id_231 vm.Value
	var arg_type_232 vm.Value
	var f_233 vm.Value
	var term_234 vm.Value
	var op_235 vm.Value
	var refs_236 vm.Value
	var aux_237 vm.Value
	var cond_id_238 vm.Value
	var tt_239 vm.Value
	var ft_240 vm.Value
	var true_edge_QMARK__241 bool
	var false_edge_QMARK__242 bool
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = term_6, v20, pred_bid_7, target_bid_8, arg_id_9, arg_type_10, f_11, term_12, pred_bid_13, target_bid_14, arg_id_15, arg_type_16, f_17, term_18, op_24, refs_26, aux_28, or__x_50, v304, pred_bid_305, target_bid_306, arg_id_307, arg_type_308, f_309, term_310, pred_bid_29, target_bid_30, arg_id_31, arg_type_32, f_33, term_34, op_35, refs_36, aux_37, pred_bid_38, target_bid_39, arg_id_40, arg_type_41, f_42, term_43, op_44, refs_45, aux_46, cond_id_90, tt_92, ft_94, arg__26044_96, true_edge_QMARK__97, arg__26049_99, false_edge_QMARK__100, v129, v293, pred_bid_294, target_bid_295, arg_id_296, arg_type_297, f_298, term_299, op_300, refs_301, aux_302, pred_bid_51, target_bid_52, arg_id_53, arg_type_54, f_55, term_56, op_57, refs_58, aux_59, or__x_60, pred_bid_61, target_bid_62, arg_id_63, arg_type_64, f_65, term_66, op_67, refs_68, aux_69, or__x_70, v74, v76, pred_bid_77, target_bid_78, arg_id_79, arg_type_80, f_81, term_82, op_83, refs_84, aux_85, or__x_86, pred_bid_101, target_bid_102, arg_id_103, arg_type_104, f_105, term_106, op_107, refs_108, aux_109, cond_id_110, tt_111, ft_112, true_edge_QMARK__113, false_edge_QMARK__114, pred_bid_115, target_bid_116, arg_id_117, arg_type_118, f_119, term_120, op_121, refs_122, aux_123, cond_id_124, tt_125, ft_126, true_edge_QMARK__127, false_edge_QMARK__128, v277, pred_bid_278, target_bid_279, arg_id_280, arg_type_281, f_282, term_283, op_284, refs_285, aux_286, cond_id_287, tt_288, ft_289, true_edge_QMARK__290, false_edge_QMARK__291, pred_bid_131, target_bid_132, arg_id_133, arg_type_134, f_135, term_136, op_137, refs_138, aux_139, cond_id_140, tt_141, ft_142, true_edge_QMARK__143, false_edge_QMARK__144, v161, pred_bid_145, target_bid_146, arg_id_147, arg_type_148, f_149, term_150, op_151, refs_152, aux_153, cond_id_154, tt_155, ft_156, true_edge_QMARK__157, false_edge_QMARK__158, v260, pred_bid_261, target_bid_262, arg_id_263, arg_type_264, f_265, term_266, op_267, refs_268, aux_269, cond_id_270, tt_271, ft_272, true_edge_QMARK__273, false_edge_QMARK__274, pred_bid_163, target_bid_164, arg_id_165, arg_type_166, f_167, term_168, op_169, refs_170, aux_171, cond_id_172, tt_173, ft_174, true_edge_QMARK__175, false_edge_QMARK__176, v193, pred_bid_177, target_bid_178, arg_id_179, arg_type_180, f_181, term_182, op_183, refs_184, aux_185, cond_id_186, tt_187, ft_188, true_edge_QMARK__189, false_edge_QMARK__190, v244, pred_bid_245, target_bid_246, arg_id_247, arg_type_248, f_249, term_250, op_251, refs_252, aux_253, cond_id_254, tt_255, ft_256, true_edge_QMARK__257, false_edge_QMARK__258, pred_bid_195, target_bid_196, arg_id_197, arg_type_198, f_199, term_200, op_201, refs_202, aux_203, cond_id_204, tt_205, ft_206, true_edge_QMARK__207, false_edge_QMARK__208, pred_bid_209, target_bid_210, arg_id_211, arg_type_212, f_213, term_214, op_215, refs_216, aux_217, cond_id_218, tt_219, ft_220, true_edge_QMARK__221, false_edge_QMARK__222, v228, pred_bid_229, target_bid_230, arg_id_231, arg_type_232, f_233, term_234, op_235, refs_236, aux_237, cond_id_238, tt_239, ft_240, true_edge_QMARK__241, false_edge_QMARK__242
	term_6, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-term").Deref(), []vm.Value{arg0, arg4})
	if callErr != nil {
		return nil, callErr
	}
	v20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{term_6})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v20) {
		pred_bid_7 = arg0
		target_bid_8 = arg1
		arg_id_9 = arg2
		arg_type_10 = arg3
		f_11 = arg4
		term_12 = term_6
		goto b1
	} else {
		pred_bid_13 = arg0
		target_bid_14 = arg1
		arg_id_15 = arg2
		arg_type_16 = arg3
		f_17 = arg4
		term_18 = term_6
		goto b2
	}
b1:
	;
	v304 = arg_type_10
	pred_bid_305 = pred_bid_7
	target_bid_306 = target_bid_8
	arg_id_307 = arg_id_9
	arg_type_308 = arg_type_10
	f_309 = f_11
	term_310 = term_12
	goto b3
b2:
	;
	op_24, callErr = rt.InvokeValue(rt.LookupVar("ir", "op").Deref(), []vm.Value{term_18, f_17})
	if callErr != nil {
		return nil, callErr
	}
	refs_26, callErr = rt.InvokeValue(rt.LookupVar("ir", "refs").Deref(), []vm.Value{term_18, f_17})
	if callErr != nil {
		return nil, callErr
	}
	aux_28, callErr = rt.InvokeValue(rt.LookupVar("ir", "aux").Deref(), []vm.Value{term_18, f_17})
	if callErr != nil {
		return nil, callErr
	}
	or__x_50, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not=").Deref(), []vm.Value{op_24, vm.Keyword("branch-if")})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(or__x_50) {
		pred_bid_51 = pred_bid_13
		target_bid_52 = target_bid_14
		arg_id_53 = arg_id_15
		arg_type_54 = arg_type_16
		f_55 = f_17
		term_56 = term_18
		op_57 = op_24
		refs_58 = refs_26
		aux_59 = aux_28
		or__x_60 = or__x_50
		goto b7
	} else {
		pred_bid_61 = pred_bid_13
		target_bid_62 = target_bid_14
		arg_id_63 = arg_id_15
		arg_type_64 = arg_type_16
		f_65 = f_17
		term_66 = term_18
		op_67 = op_24
		refs_68 = refs_26
		aux_69 = aux_28
		or__x_70 = or__x_50
		goto b8
	}
b3:
	;
	return v304, nil
b4:
	;
	v293 = arg_type_32
	pred_bid_294 = pred_bid_29
	target_bid_295 = target_bid_30
	arg_id_296 = arg_id_31
	arg_type_297 = arg_type_32
	f_298 = f_33
	term_299 = term_34
	op_300 = op_35
	refs_301 = refs_36
	aux_302 = aux_37
	goto b6
b5:
	;
	cond_id_90, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{refs_45})
	if callErr != nil {
		return nil, callErr
	}
	tt_92, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-true").Deref(), []vm.Value{aux_46})
	if callErr != nil {
		return nil, callErr
	}
	ft_94, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-false").Deref(), []vm.Value{aux_46})
	if callErr != nil {
		return nil, callErr
	}
	arg__26044_96, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-target").Deref(), []vm.Value{tt_92})
	if callErr != nil {
		return nil, callErr
	}
	true_edge_QMARK__97 = target_bid_39 == arg__26044_96
	arg__26049_99, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-target").Deref(), []vm.Value{ft_94})
	if callErr != nil {
		return nil, callErr
	}
	false_edge_QMARK__100 = target_bid_39 == arg__26049_99
	v129 = arg_id_40 == cond_id_90
	if v129 {
		pred_bid_101 = pred_bid_38
		target_bid_102 = target_bid_39
		arg_id_103 = arg_id_40
		arg_type_104 = arg_type_41
		f_105 = f_42
		term_106 = term_43
		op_107 = op_44
		refs_108 = refs_45
		aux_109 = aux_46
		cond_id_110 = cond_id_90
		tt_111 = tt_92
		ft_112 = ft_94
		true_edge_QMARK__113 = true_edge_QMARK__97
		false_edge_QMARK__114 = false_edge_QMARK__100
		goto b10
	} else {
		pred_bid_115 = pred_bid_38
		target_bid_116 = target_bid_39
		arg_id_117 = arg_id_40
		arg_type_118 = arg_type_41
		f_119 = f_42
		term_120 = term_43
		op_121 = op_44
		refs_122 = refs_45
		aux_123 = aux_46
		cond_id_124 = cond_id_90
		tt_125 = tt_92
		ft_126 = ft_94
		true_edge_QMARK__127 = true_edge_QMARK__97
		false_edge_QMARK__128 = false_edge_QMARK__100
		goto b11
	}
b6:
	;
	v304 = v293
	pred_bid_305 = pred_bid_294
	target_bid_306 = target_bid_295
	arg_id_307 = arg_id_296
	arg_type_308 = arg_type_297
	f_309 = f_298
	term_310 = term_299
	goto b3
b7:
	;
	v76 = or__x_60
	pred_bid_77 = pred_bid_51
	target_bid_78 = target_bid_52
	arg_id_79 = arg_id_53
	arg_type_80 = arg_type_54
	f_81 = f_55
	term_82 = term_56
	op_83 = op_57
	refs_84 = refs_58
	aux_85 = aux_59
	or__x_86 = or__x_60
	goto b9
b8:
	;
	v74, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "empty?").Deref(), []vm.Value{refs_68})
	if callErr != nil {
		return nil, callErr
	}
	v76 = v74
	pred_bid_77 = pred_bid_61
	target_bid_78 = target_bid_62
	arg_id_79 = arg_id_63
	arg_type_80 = arg_type_64
	f_81 = f_65
	term_82 = term_66
	op_83 = op_67
	refs_84 = refs_68
	aux_85 = aux_69
	or__x_86 = or__x_70
	goto b9
b9:
	;
	if vm.IsTruthy(v76) {
		pred_bid_29 = pred_bid_77
		target_bid_30 = target_bid_78
		arg_id_31 = arg_id_79
		arg_type_32 = arg_type_80
		f_33 = f_81
		term_34 = term_82
		op_35 = op_83
		refs_36 = refs_84
		aux_37 = aux_85
		goto b4
	} else {
		pred_bid_38 = pred_bid_77
		target_bid_39 = target_bid_78
		arg_id_40 = arg_id_79
		arg_type_41 = arg_type_80
		f_42 = f_81
		term_43 = term_82
		op_44 = op_83
		refs_45 = refs_84
		aux_46 = aux_85
		goto b5
	}
b10:
	;
	if true_edge_QMARK__113 {
		pred_bid_131 = pred_bid_101
		target_bid_132 = target_bid_102
		arg_id_133 = arg_id_103
		arg_type_134 = arg_type_104
		f_135 = f_105
		term_136 = term_106
		op_137 = op_107
		refs_138 = refs_108
		aux_139 = aux_109
		cond_id_140 = cond_id_110
		tt_141 = tt_111
		ft_142 = ft_112
		true_edge_QMARK__143 = true_edge_QMARK__113
		false_edge_QMARK__144 = false_edge_QMARK__114
		goto b13
	} else {
		pred_bid_145 = pred_bid_101
		target_bid_146 = target_bid_102
		arg_id_147 = arg_id_103
		arg_type_148 = arg_type_104
		f_149 = f_105
		term_150 = term_106
		op_151 = op_107
		refs_152 = refs_108
		aux_153 = aux_109
		cond_id_154 = cond_id_110
		tt_155 = tt_111
		ft_156 = ft_112
		true_edge_QMARK__157 = true_edge_QMARK__113
		false_edge_QMARK__158 = false_edge_QMARK__114
		goto b14
	}
b11:
	;
	v277 = arg_type_118
	pred_bid_278 = pred_bid_115
	target_bid_279 = target_bid_116
	arg_id_280 = arg_id_117
	arg_type_281 = arg_type_118
	f_282 = f_119
	term_283 = term_120
	op_284 = op_121
	refs_285 = refs_122
	aux_286 = aux_123
	cond_id_287 = cond_id_124
	tt_288 = tt_125
	ft_289 = ft_126
	true_edge_QMARK__290 = true_edge_QMARK__127
	false_edge_QMARK__291 = false_edge_QMARK__128
	goto b12
b12:
	;
	v293 = v277
	pred_bid_294 = pred_bid_278
	target_bid_295 = target_bid_279
	arg_id_296 = arg_id_280
	arg_type_297 = arg_type_281
	f_298 = f_282
	term_299 = term_283
	op_300 = op_284
	refs_301 = refs_285
	aux_302 = aux_286
	goto b6
b13:
	;
	v161, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "refine-truthy").Deref(), []vm.Value{arg_type_134})
	if callErr != nil {
		return nil, callErr
	}
	v260 = v161
	pred_bid_261 = pred_bid_131
	target_bid_262 = target_bid_132
	arg_id_263 = arg_id_133
	arg_type_264 = arg_type_134
	f_265 = f_135
	term_266 = term_136
	op_267 = op_137
	refs_268 = refs_138
	aux_269 = aux_139
	cond_id_270 = cond_id_140
	tt_271 = tt_141
	ft_272 = ft_142
	true_edge_QMARK__273 = vm.Boolean(true_edge_QMARK__143)
	false_edge_QMARK__274 = false_edge_QMARK__144
	goto b15
b14:
	;
	if false_edge_QMARK__158 {
		pred_bid_163 = pred_bid_145
		target_bid_164 = target_bid_146
		arg_id_165 = arg_id_147
		arg_type_166 = arg_type_148
		f_167 = f_149
		term_168 = term_150
		op_169 = op_151
		refs_170 = refs_152
		aux_171 = aux_153
		cond_id_172 = cond_id_154
		tt_173 = tt_155
		ft_174 = ft_156
		true_edge_QMARK__175 = true_edge_QMARK__157
		false_edge_QMARK__176 = false_edge_QMARK__158
		goto b16
	} else {
		pred_bid_177 = pred_bid_145
		target_bid_178 = target_bid_146
		arg_id_179 = arg_id_147
		arg_type_180 = arg_type_148
		f_181 = f_149
		term_182 = term_150
		op_183 = op_151
		refs_184 = refs_152
		aux_185 = aux_153
		cond_id_186 = cond_id_154
		tt_187 = tt_155
		ft_188 = ft_156
		true_edge_QMARK__189 = true_edge_QMARK__157
		false_edge_QMARK__190 = false_edge_QMARK__158
		goto b17
	}
b15:
	;
	v277 = v260
	pred_bid_278 = pred_bid_261
	target_bid_279 = target_bid_262
	arg_id_280 = arg_id_263
	arg_type_281 = arg_type_264
	f_282 = f_265
	term_283 = term_266
	op_284 = op_267
	refs_285 = refs_268
	aux_286 = aux_269
	cond_id_287 = cond_id_270
	tt_288 = tt_271
	ft_289 = ft_272
	true_edge_QMARK__290 = vm.IsTruthy(true_edge_QMARK__273)
	false_edge_QMARK__291 = false_edge_QMARK__274
	goto b12
b16:
	;
	v193, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "refine-falsey").Deref(), []vm.Value{arg_type_166})
	if callErr != nil {
		return nil, callErr
	}
	v244 = v193
	pred_bid_245 = pred_bid_163
	target_bid_246 = target_bid_164
	arg_id_247 = arg_id_165
	arg_type_248 = arg_type_166
	f_249 = f_167
	term_250 = term_168
	op_251 = op_169
	refs_252 = refs_170
	aux_253 = aux_171
	cond_id_254 = cond_id_172
	tt_255 = tt_173
	ft_256 = ft_174
	true_edge_QMARK__257 = true_edge_QMARK__175
	false_edge_QMARK__258 = vm.Boolean(false_edge_QMARK__176)
	goto b18
b17:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		pred_bid_195 = pred_bid_177
		target_bid_196 = target_bid_178
		arg_id_197 = arg_id_179
		arg_type_198 = arg_type_180
		f_199 = f_181
		term_200 = term_182
		op_201 = op_183
		refs_202 = refs_184
		aux_203 = aux_185
		cond_id_204 = cond_id_186
		tt_205 = tt_187
		ft_206 = ft_188
		true_edge_QMARK__207 = true_edge_QMARK__189
		false_edge_QMARK__208 = false_edge_QMARK__190
		goto b19
	} else {
		pred_bid_209 = pred_bid_177
		target_bid_210 = target_bid_178
		arg_id_211 = arg_id_179
		arg_type_212 = arg_type_180
		f_213 = f_181
		term_214 = term_182
		op_215 = op_183
		refs_216 = refs_184
		aux_217 = aux_185
		cond_id_218 = cond_id_186
		tt_219 = tt_187
		ft_220 = ft_188
		true_edge_QMARK__221 = true_edge_QMARK__189
		false_edge_QMARK__222 = false_edge_QMARK__190
		goto b20
	}
b18:
	;
	v260 = v244
	pred_bid_261 = pred_bid_245
	target_bid_262 = target_bid_246
	arg_id_263 = arg_id_247
	arg_type_264 = arg_type_248
	f_265 = f_249
	term_266 = term_250
	op_267 = op_251
	refs_268 = refs_252
	aux_269 = aux_253
	cond_id_270 = cond_id_254
	tt_271 = tt_255
	ft_272 = ft_256
	true_edge_QMARK__273 = vm.Boolean(true_edge_QMARK__257)
	false_edge_QMARK__274 = vm.IsTruthy(false_edge_QMARK__258)
	goto b15
b19:
	;
	v228 = arg_type_198
	pred_bid_229 = pred_bid_195
	target_bid_230 = target_bid_196
	arg_id_231 = arg_id_197
	arg_type_232 = arg_type_198
	f_233 = f_199
	term_234 = term_200
	op_235 = op_201
	refs_236 = refs_202
	aux_237 = aux_203
	cond_id_238 = cond_id_204
	tt_239 = tt_205
	ft_240 = ft_206
	true_edge_QMARK__241 = true_edge_QMARK__207
	false_edge_QMARK__242 = false_edge_QMARK__208
	goto b21
b20:
	;
	v228 = vm.NIL
	pred_bid_229 = pred_bid_209
	target_bid_230 = target_bid_210
	arg_id_231 = arg_id_211
	arg_type_232 = arg_type_212
	f_233 = f_213
	term_234 = term_214
	op_235 = op_215
	refs_236 = refs_216
	aux_237 = aux_217
	cond_id_238 = cond_id_218
	tt_239 = tt_219
	ft_240 = ft_220
	true_edge_QMARK__241 = true_edge_QMARK__221
	false_edge_QMARK__242 = false_edge_QMARK__222
	goto b21
b21:
	;
	v244 = v228
	pred_bid_245 = pred_bid_229
	target_bid_246 = target_bid_230
	arg_id_247 = arg_id_231
	arg_type_248 = arg_type_232
	f_249 = f_233
	term_250 = term_234
	op_251 = op_235
	refs_252 = refs_236
	aux_253 = aux_237
	cond_id_254 = cond_id_238
	tt_255 = tt_239
	ft_256 = ft_240
	true_edge_QMARK__257 = true_edge_QMARK__241
	false_edge_QMARK__258 = vm.Boolean(false_edge_QMARK__242)
	goto b18
}
func target_arg_types(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value, arg3 vm.Value) (vm.Value, error) {
	var term_5 vm.Value
	var v17 vm.Value
	var pred_bid_6 vm.Value
	var target_bid_7 vm.Value
	var param_pos_8 vm.Value
	var f_9 vm.Value
	var term_10 vm.Value
	var pred_bid_11 vm.Value
	var target_bid_12 vm.Value
	var param_pos_13 vm.Value
	var f_14 vm.Value
	var term_15 vm.Value
	var op_22 vm.Value
	var aux_24 vm.Value
	var v42 bool
	var v282 vm.Value
	var pred_bid_283 vm.Value
	var target_bid_284 vm.Value
	var param_pos_285 vm.Value
	var f_286 vm.Value
	var term_287 vm.Value
	var pred_bid_25 vm.Value
	var target_bid_26 vm.Value
	var param_pos_27 vm.Value
	var f_28 vm.Value
	var term_29 vm.Value
	var case__26058_30 vm.Value
	var op_31 vm.Value
	var aux_32 vm.Value
	var arg__26083_61 vm.Value
	var v62 bool
	var pred_bid_33 vm.Value
	var target_bid_34 vm.Value
	var param_pos_35 vm.Value
	var f_36 vm.Value
	var term_37 vm.Value
	var case__26058_38 vm.Value
	var op_39 vm.Value
	var aux_40 vm.Value
	var v106 bool
	var v272 vm.Value
	var pred_bid_273 vm.Value
	var target_bid_274 vm.Value
	var param_pos_275 vm.Value
	var f_276 vm.Value
	var term_277 vm.Value
	var case__26058_278 vm.Value
	var op_279 vm.Value
	var aux_280 vm.Value
	var pred_bid_44 vm.Value
	var target_bid_45 vm.Value
	var param_pos_46 vm.Value
	var f_47 vm.Value
	var term_48 vm.Value
	var case__26058_49 vm.Value
	var op_50 vm.Value
	var aux_51 vm.Value
	var arg__26087_65 vm.Value
	var arg__26093_68 vm.Value
	var arg_id_69 vm.Value
	var arg_type_71 vm.Value
	var arg__26104_74 vm.Value
	var v75 vm.Value
	var pred_bid_52 vm.Value
	var target_bid_53 vm.Value
	var param_pos_54 vm.Value
	var f_55 vm.Value
	var term_56 vm.Value
	var case__26058_57 vm.Value
	var op_58 vm.Value
	var aux_59 vm.Value
	var v79 vm.Value
	var pred_bid_80 vm.Value
	var target_bid_81 vm.Value
	var param_pos_82 vm.Value
	var f_83 vm.Value
	var term_84 vm.Value
	var case__26058_85 vm.Value
	var op_86 vm.Value
	var aux_87 vm.Value
	var pred_bid_89 vm.Value
	var target_bid_90 vm.Value
	var param_pos_91 vm.Value
	var f_92 vm.Value
	var term_93 vm.Value
	var case__26058_94 vm.Value
	var op_95 vm.Value
	var aux_96 vm.Value
	var tt_109 vm.Value
	var ft_111 vm.Value
	var arg__26117_133 vm.Value
	var v134 bool
	var pred_bid_97 vm.Value
	var target_bid_98 vm.Value
	var param_pos_99 vm.Value
	var f_100 vm.Value
	var term_101 vm.Value
	var case__26058_102 vm.Value
	var op_103 vm.Value
	var aux_104 vm.Value
	var v262 vm.Value
	var pred_bid_263 vm.Value
	var target_bid_264 vm.Value
	var param_pos_265 vm.Value
	var f_266 vm.Value
	var term_267 vm.Value
	var case__26058_268 vm.Value
	var op_269 vm.Value
	var aux_270 vm.Value
	var pred_bid_112 vm.Value
	var target_bid_113 vm.Value
	var param_pos_114 vm.Value
	var f_115 vm.Value
	var term_116 vm.Value
	var case__26058_117 vm.Value
	var op_118 vm.Value
	var aux_119 vm.Value
	var tt_120 vm.Value
	var ft_121 vm.Value
	var arg__26121_137 vm.Value
	var arg__26127_140 vm.Value
	var arg_id_141 vm.Value
	var arg_type_143 vm.Value
	var arg__26146_146 vm.Value
	var arg__26159_149 vm.Value
	var arg__26160_150 vm.Value
	var v151 vm.Value
	var pred_bid_122 vm.Value
	var target_bid_123 vm.Value
	var param_pos_124 vm.Value
	var f_125 vm.Value
	var term_126 vm.Value
	var case__26058_127 vm.Value
	var op_128 vm.Value
	var aux_129 vm.Value
	var tt_130 vm.Value
	var ft_131 vm.Value
	var tt_types_155 vm.Value
	var pred_bid_156 vm.Value
	var target_bid_157 vm.Value
	var param_pos_158 vm.Value
	var f_159 vm.Value
	var term_160 vm.Value
	var case__26058_161 vm.Value
	var op_162 vm.Value
	var aux_163 vm.Value
	var tt_164 vm.Value
	var ft_165 vm.Value
	var arg__26165_189 vm.Value
	var v190 bool
	var tt_types_166 vm.Value
	var pred_bid_167 vm.Value
	var target_bid_168 vm.Value
	var param_pos_169 vm.Value
	var f_170 vm.Value
	var term_171 vm.Value
	var case__26058_172 vm.Value
	var op_173 vm.Value
	var aux_174 vm.Value
	var tt_175 vm.Value
	var ft_176 vm.Value
	var arg__26169_193 vm.Value
	var arg__26175_196 vm.Value
	var arg_id_197 vm.Value
	var arg_type_199 vm.Value
	var arg__26194_202 vm.Value
	var arg__26207_205 vm.Value
	var arg__26208_206 vm.Value
	var v207 vm.Value
	var tt_types_177 vm.Value
	var pred_bid_178 vm.Value
	var target_bid_179 vm.Value
	var param_pos_180 vm.Value
	var f_181 vm.Value
	var term_182 vm.Value
	var case__26058_183 vm.Value
	var op_184 vm.Value
	var aux_185 vm.Value
	var tt_186 vm.Value
	var ft_187 vm.Value
	var ft_types_211 vm.Value
	var tt_types_212 vm.Value
	var pred_bid_213 vm.Value
	var target_bid_214 vm.Value
	var param_pos_215 vm.Value
	var f_216 vm.Value
	var term_217 vm.Value
	var case__26058_218 vm.Value
	var op_219 vm.Value
	var aux_220 vm.Value
	var tt_221 vm.Value
	var ft_222 vm.Value
	var arg__26214_224 vm.Value
	var arg__26221_227 vm.Value
	var v228 vm.Value
	var pred_bid_230 vm.Value
	var target_bid_231 vm.Value
	var param_pos_232 vm.Value
	var f_233 vm.Value
	var term_234 vm.Value
	var case__26058_235 vm.Value
	var op_236 vm.Value
	var aux_237 vm.Value
	var pred_bid_238 vm.Value
	var target_bid_239 vm.Value
	var param_pos_240 vm.Value
	var f_241 vm.Value
	var term_242 vm.Value
	var case__26058_243 vm.Value
	var op_244 vm.Value
	var aux_245 vm.Value
	var v252 vm.Value
	var pred_bid_253 vm.Value
	var target_bid_254 vm.Value
	var param_pos_255 vm.Value
	var f_256 vm.Value
	var term_257 vm.Value
	var case__26058_258 vm.Value
	var op_259 vm.Value
	var aux_260 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = term_5, v17, pred_bid_6, target_bid_7, param_pos_8, f_9, term_10, pred_bid_11, target_bid_12, param_pos_13, f_14, term_15, op_22, aux_24, v42, v282, pred_bid_283, target_bid_284, param_pos_285, f_286, term_287, pred_bid_25, target_bid_26, param_pos_27, f_28, term_29, case__26058_30, op_31, aux_32, arg__26083_61, v62, pred_bid_33, target_bid_34, param_pos_35, f_36, term_37, case__26058_38, op_39, aux_40, v106, v272, pred_bid_273, target_bid_274, param_pos_275, f_276, term_277, case__26058_278, op_279, aux_280, pred_bid_44, target_bid_45, param_pos_46, f_47, term_48, case__26058_49, op_50, aux_51, arg__26087_65, arg__26093_68, arg_id_69, arg_type_71, arg__26104_74, v75, pred_bid_52, target_bid_53, param_pos_54, f_55, term_56, case__26058_57, op_58, aux_59, v79, pred_bid_80, target_bid_81, param_pos_82, f_83, term_84, case__26058_85, op_86, aux_87, pred_bid_89, target_bid_90, param_pos_91, f_92, term_93, case__26058_94, op_95, aux_96, tt_109, ft_111, arg__26117_133, v134, pred_bid_97, target_bid_98, param_pos_99, f_100, term_101, case__26058_102, op_103, aux_104, v262, pred_bid_263, target_bid_264, param_pos_265, f_266, term_267, case__26058_268, op_269, aux_270, pred_bid_112, target_bid_113, param_pos_114, f_115, term_116, case__26058_117, op_118, aux_119, tt_120, ft_121, arg__26121_137, arg__26127_140, arg_id_141, arg_type_143, arg__26146_146, arg__26159_149, arg__26160_150, v151, pred_bid_122, target_bid_123, param_pos_124, f_125, term_126, case__26058_127, op_128, aux_129, tt_130, ft_131, tt_types_155, pred_bid_156, target_bid_157, param_pos_158, f_159, term_160, case__26058_161, op_162, aux_163, tt_164, ft_165, arg__26165_189, v190, tt_types_166, pred_bid_167, target_bid_168, param_pos_169, f_170, term_171, case__26058_172, op_173, aux_174, tt_175, ft_176, arg__26169_193, arg__26175_196, arg_id_197, arg_type_199, arg__26194_202, arg__26207_205, arg__26208_206, v207, tt_types_177, pred_bid_178, target_bid_179, param_pos_180, f_181, term_182, case__26058_183, op_184, aux_185, tt_186, ft_187, ft_types_211, tt_types_212, pred_bid_213, target_bid_214, param_pos_215, f_216, term_217, case__26058_218, op_219, aux_220, tt_221, ft_222, arg__26214_224, arg__26221_227, v228, pred_bid_230, target_bid_231, param_pos_232, f_233, term_234, case__26058_235, op_236, aux_237, pred_bid_238, target_bid_239, param_pos_240, f_241, term_242, case__26058_243, op_244, aux_245, v252, pred_bid_253, target_bid_254, param_pos_255, f_256, term_257, case__26058_258, op_259, aux_260
	term_5, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-term").Deref(), []vm.Value{arg0, arg3})
	if callErr != nil {
		return nil, callErr
	}
	v17, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{term_5})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v17) {
		pred_bid_6 = arg0
		target_bid_7 = arg1
		param_pos_8 = arg2
		f_9 = arg3
		term_10 = term_5
		goto b1
	} else {
		pred_bid_11 = arg0
		target_bid_12 = arg1
		param_pos_13 = arg2
		f_14 = arg3
		term_15 = term_5
		goto b2
	}
b1:
	;
	v282 = vm.NewArrayVector([]vm.Value{})
	pred_bid_283 = pred_bid_6
	target_bid_284 = target_bid_7
	param_pos_285 = param_pos_8
	f_286 = f_9
	term_287 = term_10
	goto b3
b2:
	;
	op_22, callErr = rt.InvokeValue(rt.LookupVar("ir", "op").Deref(), []vm.Value{term_15, f_14})
	if callErr != nil {
		return nil, callErr
	}
	aux_24, callErr = rt.InvokeValue(rt.LookupVar("ir", "aux").Deref(), []vm.Value{term_15, f_14})
	if callErr != nil {
		return nil, callErr
	}
	v42 = op_22 == vm.Keyword("branch")
	if v42 {
		pred_bid_25 = pred_bid_11
		target_bid_26 = target_bid_12
		param_pos_27 = param_pos_13
		f_28 = f_14
		term_29 = term_15
		case__26058_30 = op_22
		op_31 = op_22
		aux_32 = aux_24
		goto b4
	} else {
		pred_bid_33 = pred_bid_11
		target_bid_34 = target_bid_12
		param_pos_35 = param_pos_13
		f_36 = f_14
		term_37 = term_15
		case__26058_38 = op_22
		op_39 = op_22
		aux_40 = aux_24
		goto b5
	}
b3:
	;
	return v282, nil
b4:
	;
	arg__26083_61, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-target").Deref(), []vm.Value{aux_32})
	if callErr != nil {
		return nil, callErr
	}
	v62 = target_bid_26 == arg__26083_61
	if v62 {
		pred_bid_44 = pred_bid_25
		target_bid_45 = target_bid_26
		param_pos_46 = param_pos_27
		f_47 = f_28
		term_48 = term_29
		case__26058_49 = case__26058_30
		op_50 = op_31
		aux_51 = aux_32
		goto b7
	} else {
		pred_bid_52 = pred_bid_25
		target_bid_53 = target_bid_26
		param_pos_54 = param_pos_27
		f_55 = f_28
		term_56 = term_29
		case__26058_57 = case__26058_30
		op_58 = op_31
		aux_59 = aux_32
		goto b8
	}
b5:
	;
	v106 = case__26058_38 == vm.Keyword("branch-if")
	if v106 {
		pred_bid_89 = pred_bid_33
		target_bid_90 = target_bid_34
		param_pos_91 = param_pos_35
		f_92 = f_36
		term_93 = term_37
		case__26058_94 = case__26058_38
		op_95 = op_39
		aux_96 = aux_40
		goto b10
	} else {
		pred_bid_97 = pred_bid_33
		target_bid_98 = target_bid_34
		param_pos_99 = param_pos_35
		f_100 = f_36
		term_101 = term_37
		case__26058_102 = case__26058_38
		op_103 = op_39
		aux_104 = aux_40
		goto b11
	}
b6:
	;
	v282 = v272
	pred_bid_283 = pred_bid_273
	target_bid_284 = target_bid_274
	param_pos_285 = param_pos_275
	f_286 = f_276
	term_287 = term_277
	goto b3
b7:
	;
	arg__26087_65, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-args").Deref(), []vm.Value{aux_51})
	if callErr != nil {
		return nil, callErr
	}
	arg__26093_68, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-args").Deref(), []vm.Value{aux_51})
	if callErr != nil {
		return nil, callErr
	}
	arg_id_69, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__26093_68, param_pos_46})
	if callErr != nil {
		return nil, callErr
	}
	arg_type_71, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg_id_69, f_47})
	if callErr != nil {
		return nil, callErr
	}
	arg__26104_74, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "normalize-type").Deref(), []vm.Value{arg_type_71})
	if callErr != nil {
		return nil, callErr
	}
	v75, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__26104_74})
	if callErr != nil {
		return nil, callErr
	}
	v79 = v75
	pred_bid_80 = pred_bid_44
	target_bid_81 = target_bid_45
	param_pos_82 = param_pos_46
	f_83 = f_47
	term_84 = term_48
	case__26058_85 = case__26058_49
	op_86 = op_50
	aux_87 = aux_51
	goto b9
b8:
	;
	v79 = vm.NewArrayVector([]vm.Value{})
	pred_bid_80 = pred_bid_52
	target_bid_81 = target_bid_53
	param_pos_82 = param_pos_54
	f_83 = f_55
	term_84 = term_56
	case__26058_85 = case__26058_57
	op_86 = op_58
	aux_87 = aux_59
	goto b9
b9:
	;
	v272 = v79
	pred_bid_273 = pred_bid_80
	target_bid_274 = target_bid_81
	param_pos_275 = param_pos_82
	f_276 = f_83
	term_277 = term_84
	case__26058_278 = case__26058_85
	op_279 = op_86
	aux_280 = aux_87
	goto b6
b10:
	;
	tt_109, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-true").Deref(), []vm.Value{aux_96})
	if callErr != nil {
		return nil, callErr
	}
	ft_111, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-false").Deref(), []vm.Value{aux_96})
	if callErr != nil {
		return nil, callErr
	}
	arg__26117_133, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-target").Deref(), []vm.Value{tt_109})
	if callErr != nil {
		return nil, callErr
	}
	v134 = target_bid_90 == arg__26117_133
	if v134 {
		pred_bid_112 = pred_bid_89
		target_bid_113 = target_bid_90
		param_pos_114 = param_pos_91
		f_115 = f_92
		term_116 = term_93
		case__26058_117 = case__26058_94
		op_118 = op_95
		aux_119 = aux_96
		tt_120 = tt_109
		ft_121 = ft_111
		goto b13
	} else {
		pred_bid_122 = pred_bid_89
		target_bid_123 = target_bid_90
		param_pos_124 = param_pos_91
		f_125 = f_92
		term_126 = term_93
		case__26058_127 = case__26058_94
		op_128 = op_95
		aux_129 = aux_96
		tt_130 = tt_109
		ft_131 = ft_111
		goto b14
	}
b11:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		pred_bid_230 = pred_bid_97
		target_bid_231 = target_bid_98
		param_pos_232 = param_pos_99
		f_233 = f_100
		term_234 = term_101
		case__26058_235 = case__26058_102
		op_236 = op_103
		aux_237 = aux_104
		goto b19
	} else {
		pred_bid_238 = pred_bid_97
		target_bid_239 = target_bid_98
		param_pos_240 = param_pos_99
		f_241 = f_100
		term_242 = term_101
		case__26058_243 = case__26058_102
		op_244 = op_103
		aux_245 = aux_104
		goto b20
	}
b12:
	;
	v272 = v262
	pred_bid_273 = pred_bid_263
	target_bid_274 = target_bid_264
	param_pos_275 = param_pos_265
	f_276 = f_266
	term_277 = term_267
	case__26058_278 = case__26058_268
	op_279 = op_269
	aux_280 = aux_270
	goto b6
b13:
	;
	arg__26121_137, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-args").Deref(), []vm.Value{tt_120})
	if callErr != nil {
		return nil, callErr
	}
	arg__26127_140, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-args").Deref(), []vm.Value{tt_120})
	if callErr != nil {
		return nil, callErr
	}
	arg_id_141, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__26127_140, param_pos_114})
	if callErr != nil {
		return nil, callErr
	}
	arg_type_143, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg_id_141, f_115})
	if callErr != nil {
		return nil, callErr
	}
	arg__26146_146, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "refine-edge-arg-type").Deref(), []vm.Value{pred_bid_112, target_bid_113, arg_id_141, arg_type_143, f_115})
	if callErr != nil {
		return nil, callErr
	}
	arg__26159_149, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "refine-edge-arg-type").Deref(), []vm.Value{pred_bid_112, target_bid_113, arg_id_141, arg_type_143, f_115})
	if callErr != nil {
		return nil, callErr
	}
	arg__26160_150, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "normalize-type").Deref(), []vm.Value{arg__26159_149})
	if callErr != nil {
		return nil, callErr
	}
	v151, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__26160_150})
	if callErr != nil {
		return nil, callErr
	}
	tt_types_155 = v151
	pred_bid_156 = pred_bid_112
	target_bid_157 = target_bid_113
	param_pos_158 = param_pos_114
	f_159 = f_115
	term_160 = term_116
	case__26058_161 = case__26058_117
	op_162 = op_118
	aux_163 = aux_119
	tt_164 = tt_120
	ft_165 = ft_121
	goto b15
b14:
	;
	tt_types_155 = vm.NewArrayVector([]vm.Value{})
	pred_bid_156 = pred_bid_122
	target_bid_157 = target_bid_123
	param_pos_158 = param_pos_124
	f_159 = f_125
	term_160 = term_126
	case__26058_161 = case__26058_127
	op_162 = op_128
	aux_163 = aux_129
	tt_164 = tt_130
	ft_165 = ft_131
	goto b15
b15:
	;
	arg__26165_189, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-target").Deref(), []vm.Value{ft_165})
	if callErr != nil {
		return nil, callErr
	}
	v190 = target_bid_157 == arg__26165_189
	if v190 {
		tt_types_166 = tt_types_155
		pred_bid_167 = pred_bid_156
		target_bid_168 = target_bid_157
		param_pos_169 = param_pos_158
		f_170 = f_159
		term_171 = term_160
		case__26058_172 = case__26058_161
		op_173 = op_162
		aux_174 = aux_163
		tt_175 = tt_164
		ft_176 = ft_165
		goto b16
	} else {
		tt_types_177 = tt_types_155
		pred_bid_178 = pred_bid_156
		target_bid_179 = target_bid_157
		param_pos_180 = param_pos_158
		f_181 = f_159
		term_182 = term_160
		case__26058_183 = case__26058_161
		op_184 = op_162
		aux_185 = aux_163
		tt_186 = tt_164
		ft_187 = ft_165
		goto b17
	}
b16:
	;
	arg__26169_193, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-args").Deref(), []vm.Value{ft_176})
	if callErr != nil {
		return nil, callErr
	}
	arg__26175_196, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-args").Deref(), []vm.Value{ft_176})
	if callErr != nil {
		return nil, callErr
	}
	arg_id_197, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__26175_196, param_pos_169})
	if callErr != nil {
		return nil, callErr
	}
	arg_type_199, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg_id_197, f_170})
	if callErr != nil {
		return nil, callErr
	}
	arg__26194_202, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "refine-edge-arg-type").Deref(), []vm.Value{pred_bid_167, target_bid_168, arg_id_197, arg_type_199, f_170})
	if callErr != nil {
		return nil, callErr
	}
	arg__26207_205, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "refine-edge-arg-type").Deref(), []vm.Value{pred_bid_167, target_bid_168, arg_id_197, arg_type_199, f_170})
	if callErr != nil {
		return nil, callErr
	}
	arg__26208_206, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "normalize-type").Deref(), []vm.Value{arg__26207_205})
	if callErr != nil {
		return nil, callErr
	}
	v207, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__26208_206})
	if callErr != nil {
		return nil, callErr
	}
	ft_types_211 = v207
	tt_types_212 = tt_types_166
	pred_bid_213 = pred_bid_167
	target_bid_214 = target_bid_168
	param_pos_215 = param_pos_169
	f_216 = f_170
	term_217 = term_171
	case__26058_218 = case__26058_172
	op_219 = op_173
	aux_220 = aux_174
	tt_221 = tt_175
	ft_222 = ft_176
	goto b18
b17:
	;
	ft_types_211 = vm.NewArrayVector([]vm.Value{})
	tt_types_212 = tt_types_177
	pred_bid_213 = pred_bid_178
	target_bid_214 = target_bid_179
	param_pos_215 = param_pos_180
	f_216 = f_181
	term_217 = term_182
	case__26058_218 = case__26058_183
	op_219 = op_184
	aux_220 = aux_185
	tt_221 = tt_186
	ft_222 = ft_187
	goto b18
b18:
	;
	arg__26214_224, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "concat").Deref(), []vm.Value{tt_types_212, ft_types_211})
	if callErr != nil {
		return nil, callErr
	}
	arg__26221_227, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "concat").Deref(), []vm.Value{tt_types_212, ft_types_211})
	if callErr != nil {
		return nil, callErr
	}
	v228, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg__26221_227})
	if callErr != nil {
		return nil, callErr
	}
	v262 = v228
	pred_bid_263 = pred_bid_213
	target_bid_264 = target_bid_214
	param_pos_265 = param_pos_215
	f_266 = f_216
	term_267 = term_217
	case__26058_268 = case__26058_218
	op_269 = op_219
	aux_270 = aux_220
	goto b12
b19:
	;
	v252 = vm.NewArrayVector([]vm.Value{})
	pred_bid_253 = pred_bid_230
	target_bid_254 = target_bid_231
	param_pos_255 = param_pos_232
	f_256 = f_233
	term_257 = term_234
	case__26058_258 = case__26058_235
	op_259 = op_236
	aux_260 = aux_237
	goto b21
b20:
	;
	v252 = vm.NIL
	pred_bid_253 = pred_bid_238
	target_bid_254 = target_bid_239
	param_pos_255 = param_pos_240
	f_256 = f_241
	term_257 = term_242
	case__26058_258 = case__26058_243
	op_259 = op_244
	aux_260 = aux_245
	goto b21
b21:
	;
	v262 = v252
	pred_bid_263 = pred_bid_253
	target_bid_264 = target_bid_254
	param_pos_265 = param_pos_255
	f_266 = f_256
	term_267 = term_257
	case__26058_268 = case__26058_258
	op_269 = op_259
	aux_270 = aux_260
	goto b12
}
func param_has_ready_source_QMARK_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value, arg3 vm.Value) (vm.Value, error) {
	var arg__26260_8 vm.Value
	var arg__26264_10 vm.Value
	var arg__26268_13 vm.Value
	var arg__26272_15 vm.Value
	var arg__26273_16 vm.Value
	var arg__26313_22 vm.Value
	var arg__26317_24 vm.Value
	var arg__26321_27 vm.Value
	var arg__26325_29 vm.Value
	var arg__26326_30 vm.Value
	var v31 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _ = arg__26260_8, arg__26264_10, arg__26268_13, arg__26272_15, arg__26273_16, arg__26313_22, arg__26317_24, arg__26321_27, arg__26325_29, arg__26326_30, v31
	arg__26260_8, callErr = rt.InvokeValue(vm.Keyword("param-sources"), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__26264_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg1, arg2})
	if callErr != nil {
		return nil, callErr
	}
	arg__26268_13, callErr = rt.InvokeValue(vm.Keyword("param-sources"), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__26272_15, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg1, arg2})
	if callErr != nil {
		return nil, callErr
	}
	arg__26273_16, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg__26268_13, arg__26272_15})
	if callErr != nil {
		return nil, callErr
	}
	arg__26313_22, callErr = rt.InvokeValue(vm.Keyword("param-sources"), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__26317_24, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg1, arg2})
	if callErr != nil {
		return nil, callErr
	}
	arg__26321_27, callErr = rt.InvokeValue(vm.Keyword("param-sources"), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__26325_29, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg1, arg2})
	if callErr != nil {
		return nil, callErr
	}
	arg__26326_30, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg__26321_27, arg__26325_29})
	if callErr != nil {
		return nil, callErr
	}
	v31, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "some").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var arg__26281_6 vm.Value
		var arg__26289_11 vm.Value
		var arg__26291_12 vm.Value
		var arg__26299_18 vm.Value
		var arg__26307_23 vm.Value
		var arg__26309_24 vm.Value
		var v25 vm.Value
		var callErr error
		_, _, _, _, _, _, _ = arg__26281_6, arg__26289_11, arg__26291_12, arg__26299_18, arg__26307_23, arg__26309_24, v25
		arg__26281_6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(2)})
		if callErr != nil {
			return nil, callErr
		}
		arg__26289_11, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(2)})
		if callErr != nil {
			return nil, callErr
		}
		arg__26291_12, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg__26289_11, arg3})
		if callErr != nil {
			return nil, callErr
		}
		arg__26299_18, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(2)})
		if callErr != nil {
			return nil, callErr
		}
		arg__26307_23, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(2)})
		if callErr != nil {
			return nil, callErr
		}
		arg__26309_24, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg__26307_23, arg3})
		if callErr != nil {
			return nil, callErr
		}
		v25, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not=").Deref(), []vm.Value{vm.Keyword("bottom"), arg__26309_24})
		if callErr != nil {
			return nil, callErr
		}
		return v25, nil
	}), arg__26326_30})
	if callErr != nil {
		return nil, callErr
	}
	return v31, nil
}
func enqueue_entry(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var v6 vm.Value
	var v14 vm.Value
	var queue_7 vm.Value
	var queued_8 vm.Value
	var entry_9 vm.Value
	var v17 vm.Value
	var queue_10 vm.Value
	var queued_11 vm.Value
	var entry_12 vm.Value
	var v22 vm.Value
	var arg__26347_25 vm.Value
	var arg__26353_27 vm.Value
	var v28 vm.Value
	var v30 vm.Value
	var queue_31 vm.Value
	var queued_32 vm.Value
	var entry_33 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = v6, v14, queue_7, queued_8, entry_9, v17, queue_10, queued_11, entry_12, v22, arg__26347_25, arg__26353_27, v28, v30, queue_31, queued_32, entry_33
	v6, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "ti-inc!").Deref(), []vm.Value{vm.Keyword("enqueue-attempt")})
	if callErr != nil {
		return nil, callErr
	}
	v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "contains?").Deref(), []vm.Value{arg1, arg2})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v14) {
		queue_7 = arg0
		queued_8 = arg1
		entry_9 = arg2
		goto b1
	} else {
		queue_10 = arg0
		queued_11 = arg1
		entry_12 = arg2
		goto b2
	}
b1:
	;
	v17, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{queue_7, queued_8})
	if callErr != nil {
		return nil, callErr
	}
	v30 = v17
	queue_31 = queue_7
	queued_32 = queued_8
	entry_33 = entry_9
	goto b3
b2:
	;
	v22, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "ti-inc!").Deref(), []vm.Value{vm.Keyword("enqueue")})
	if callErr != nil {
		return nil, callErr
	}
	arg__26347_25, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{queue_10, entry_12})
	if callErr != nil {
		return nil, callErr
	}
	arg__26353_27, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{queued_11, entry_12})
	if callErr != nil {
		return nil, callErr
	}
	v28, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg__26347_25, arg__26353_27})
	if callErr != nil {
		return nil, callErr
	}
	v30 = v28
	queue_31 = queue_10
	queued_32 = queued_11
	entry_33 = entry_12
	goto b3
b3:
	;
	return v30, nil
}
func join_all(arg0 vm.Value) (vm.Value, error) {
	var v6 vm.Value
	var callErr error
	_ = v6
	v6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "reduce").Deref(), []vm.Value{rt.LookupVar("ir.passes.typeinfer", "type-join").Deref(), vm.Keyword("bottom"), arg0})
	if callErr != nil {
		return nil, callErr
	}
	return v6, nil
}
func numeric_type_QMARK_(arg0 vm.Value) (bool, error) {
	var or__x_2 bool
	var t_3 vm.Value
	var or__x_4 bool
	var t_5 vm.Value
	var or__x_6 bool
	var or__x_10 bool
	var v56 bool
	var t_57 vm.Value
	var or__x_58 vm.Value
	var t_11 vm.Value
	var or__x_12 bool
	var t_13 vm.Value
	var or__x_14 bool
	var or__x_18 bool
	var v52 bool
	var t_53 vm.Value
	var or__x_54 vm.Value
	var t_19 vm.Value
	var or__x_20 bool
	var t_21 vm.Value
	var or__x_22 bool
	var arg__26372_29 vm.Value
	var or__x_30 bool
	var v48 bool
	var t_49 vm.Value
	var or__x_50 vm.Value
	var t_31 vm.Value
	var or__x_32 bool
	var t_33 vm.Value
	var or__x_34 bool
	var arg__26378_41 vm.Value
	var v42 bool
	var v44 bool
	var t_45 vm.Value
	var or__x_46 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = or__x_2, t_3, or__x_4, t_5, or__x_6, or__x_10, v56, t_57, or__x_58, t_11, or__x_12, t_13, or__x_14, or__x_18, v52, t_53, or__x_54, t_19, or__x_20, t_21, or__x_22, arg__26372_29, or__x_30, v48, t_49, or__x_50, t_31, or__x_32, t_33, or__x_34, arg__26378_41, v42, v44, t_45, or__x_46
	or__x_2 = arg0 == vm.Keyword("int")
	if or__x_2 {
		t_3 = arg0
		or__x_4 = or__x_2
		goto b1
	} else {
		t_5 = arg0
		or__x_6 = or__x_2
		goto b2
	}
b1:
	;
	v56 = or__x_4
	t_57 = t_3
	or__x_58 = vm.Boolean(or__x_4)
	goto b3
b2:
	;
	or__x_10 = t_5 == vm.Keyword("float")
	if or__x_10 {
		t_11 = t_5
		or__x_12 = or__x_10
		goto b4
	} else {
		t_13 = t_5
		or__x_14 = or__x_10
		goto b5
	}
b3:
	;
	return v56, nil
b4:
	;
	v52 = or__x_12
	t_53 = t_11
	or__x_54 = vm.Boolean(or__x_12)
	goto b6
b5:
	;
	or__x_18 = t_13 == vm.Keyword("number")
	if or__x_18 {
		t_19 = t_13
		or__x_20 = or__x_18
		goto b7
	} else {
		t_21 = t_13
		or__x_22 = or__x_18
		goto b8
	}
b6:
	;
	v56 = v52
	t_57 = t_53
	or__x_58 = vm.Boolean(or__x_6)
	goto b3
b7:
	;
	v48 = or__x_20
	t_49 = t_19
	or__x_50 = vm.Boolean(or__x_20)
	goto b9
b8:
	;
	arg__26372_29, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("const"), vm.Keyword("int"), vm.Int(0)})
	if callErr != nil {
		return false, callErr
	}
	or__x_30 = t_21 == arg__26372_29
	if or__x_30 {
		t_31 = t_21
		or__x_32 = or__x_30
		goto b10
	} else {
		t_33 = t_21
		or__x_34 = or__x_30
		goto b11
	}
b9:
	;
	v52 = v48
	t_53 = t_49
	or__x_54 = vm.Boolean(or__x_14)
	goto b6
b10:
	;
	v44 = or__x_32
	t_45 = t_31
	or__x_46 = vm.Boolean(or__x_32)
	goto b12
b11:
	;
	arg__26378_41, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("const"), vm.Keyword("float"), vm.Float(0)})
	if callErr != nil {
		return false, callErr
	}
	v42 = t_33 == arg__26378_41
	v44 = v42
	t_45 = t_33
	or__x_46 = vm.Boolean(or__x_34)
	goto b12
b12:
	;
	v48 = v44
	t_49 = t_45
	or__x_50 = vm.Boolean(or__x_22)
	goto b9
}
func type_join(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var v5 vm.Value
	var v10 bool
	var a_6 vm.Value
	var b_7 vm.Value
	var a_8 vm.Value
	var b_9 vm.Value
	var v18 bool
	var v675 vm.Value
	var a_676 vm.Value
	var b_677 vm.Value
	var a_13 vm.Value
	var b_14 vm.Value
	var a_15 vm.Value
	var b_16 vm.Value
	var v26 bool
	var v671 vm.Value
	var a_672 vm.Value
	var b_673 vm.Value
	var a_21 vm.Value
	var b_22 vm.Value
	var a_23 vm.Value
	var b_24 vm.Value
	var v34 bool
	var v667 vm.Value
	var a_668 vm.Value
	var b_669 vm.Value
	var a_29 vm.Value
	var b_30 vm.Value
	var a_31 vm.Value
	var b_32 vm.Value
	var v43 bool
	var v663 vm.Value
	var a_664 vm.Value
	var b_665 vm.Value
	var a_38 vm.Value
	var b_39 vm.Value
	var a_40 vm.Value
	var b_41 vm.Value
	var v52 bool
	var v659 vm.Value
	var a_660 vm.Value
	var b_661 vm.Value
	var a_47 vm.Value
	var b_48 vm.Value
	var a_49 vm.Value
	var b_50 vm.Value
	var v61 bool
	var v655 vm.Value
	var a_656 vm.Value
	var b_657 vm.Value
	var a_56 vm.Value
	var b_57 vm.Value
	var a_58 vm.Value
	var b_59 vm.Value
	var and__x_70 bool
	var v651 vm.Value
	var a_652 vm.Value
	var b_653 vm.Value
	var a_65 vm.Value
	var b_66 vm.Value
	var a_67 vm.Value
	var b_68 vm.Value
	var and__x_98 bool
	var v647 vm.Value
	var a_648 vm.Value
	var b_649 vm.Value
	var a_71 vm.Value
	var b_72 vm.Value
	var and__x_73 bool
	var arg__26403_82 vm.Value
	var v83 bool
	var a_74 vm.Value
	var b_75 vm.Value
	var and__x_76 bool
	var v86 bool
	var a_87 vm.Value
	var b_88 vm.Value
	var and__x_89 vm.Value
	var a_93 vm.Value
	var b_94 vm.Value
	var a_95 vm.Value
	var b_96 vm.Value
	var and__x_126 bool
	var v643 vm.Value
	var a_644 vm.Value
	var b_645 vm.Value
	var a_99 vm.Value
	var b_100 vm.Value
	var and__x_101 bool
	var arg__26411_110 vm.Value
	var v111 bool
	var a_102 vm.Value
	var b_103 vm.Value
	var and__x_104 bool
	var v114 bool
	var a_115 vm.Value
	var b_116 vm.Value
	var and__x_117 vm.Value
	var a_121 vm.Value
	var b_122 vm.Value
	var a_123 vm.Value
	var b_124 vm.Value
	var and__x_154 bool
	var v639 vm.Value
	var a_640 vm.Value
	var b_641 vm.Value
	var a_127 vm.Value
	var b_128 vm.Value
	var and__x_129 bool
	var arg__26419_138 vm.Value
	var v139 bool
	var a_130 vm.Value
	var b_131 vm.Value
	var and__x_132 bool
	var v142 bool
	var a_143 vm.Value
	var b_144 vm.Value
	var and__x_145 vm.Value
	var a_149 vm.Value
	var b_150 vm.Value
	var a_151 vm.Value
	var b_152 vm.Value
	var and__x_182 bool
	var v635 vm.Value
	var a_636 vm.Value
	var b_637 vm.Value
	var a_155 vm.Value
	var b_156 vm.Value
	var and__x_157 bool
	var arg__26427_166 vm.Value
	var v167 bool
	var a_158 vm.Value
	var b_159 vm.Value
	var and__x_160 bool
	var v170 bool
	var a_171 vm.Value
	var b_172 vm.Value
	var and__x_173 vm.Value
	var a_177 vm.Value
	var b_178 vm.Value
	var a_179 vm.Value
	var b_180 vm.Value
	var and__x_224 bool
	var v631 vm.Value
	var a_632 vm.Value
	var b_633 vm.Value
	var a_183 vm.Value
	var b_184 vm.Value
	var and__x_185 bool
	var or__x_191 bool
	var a_186 vm.Value
	var b_187 vm.Value
	var and__x_188 bool
	var v212 bool
	var a_213 vm.Value
	var b_214 vm.Value
	var and__x_215 vm.Value
	var a_192 vm.Value
	var b_193 vm.Value
	var and__x_194 bool
	var or__x_195 bool
	var a_196 vm.Value
	var b_197 vm.Value
	var and__x_198 bool
	var or__x_199 bool
	var v203 bool
	var v205 bool
	var a_206 vm.Value
	var b_207 vm.Value
	var and__x_208 bool
	var or__x_209 vm.Value
	var a_219 vm.Value
	var b_220 vm.Value
	var a_221 vm.Value
	var b_222 vm.Value
	var v266 bool
	var v627 vm.Value
	var a_628 vm.Value
	var b_629 vm.Value
	var a_225 vm.Value
	var b_226 vm.Value
	var and__x_227 bool
	var or__x_233 bool
	var a_228 vm.Value
	var b_229 vm.Value
	var and__x_230 bool
	var v254 bool
	var a_255 vm.Value
	var b_256 vm.Value
	var and__x_257 vm.Value
	var a_234 vm.Value
	var b_235 vm.Value
	var and__x_236 bool
	var or__x_237 bool
	var a_238 vm.Value
	var b_239 vm.Value
	var and__x_240 bool
	var or__x_241 bool
	var v245 bool
	var v247 bool
	var a_248 vm.Value
	var b_249 vm.Value
	var and__x_250 bool
	var or__x_251 vm.Value
	var a_261 vm.Value
	var b_262 vm.Value
	var or__x_273 bool
	var a_263 vm.Value
	var b_264 vm.Value
	var v363 bool
	var v623 vm.Value
	var a_624 vm.Value
	var b_625 vm.Value
	var a_268 vm.Value
	var b_269 vm.Value
	var a_270 vm.Value
	var b_271 vm.Value
	var arg__26464_347 vm.Value
	var arg__26470_351 vm.Value
	var v352 vm.Value
	var v354 vm.Value
	var a_355 vm.Value
	var b_356 vm.Value
	var a_274 vm.Value
	var b_275 vm.Value
	var or__x_276 bool
	var a_277 vm.Value
	var b_278 vm.Value
	var or__x_279 bool
	var or__x_283 bool
	var v338 bool
	var a_339 vm.Value
	var b_340 vm.Value
	var or__x_341 vm.Value
	var a_284 vm.Value
	var b_285 vm.Value
	var or__x_286 bool
	var a_287 vm.Value
	var b_288 vm.Value
	var or__x_289 bool
	var or__x_293 bool
	var v333 bool
	var a_334 vm.Value
	var b_335 vm.Value
	var or__x_336 vm.Value
	var a_294 vm.Value
	var b_295 vm.Value
	var or__x_296 bool
	var a_297 vm.Value
	var b_298 vm.Value
	var or__x_299 bool
	var arg__26453_306 vm.Value
	var or__x_307 bool
	var v328 bool
	var a_329 vm.Value
	var b_330 vm.Value
	var or__x_331 vm.Value
	var a_308 vm.Value
	var b_309 vm.Value
	var or__x_310 bool
	var a_311 vm.Value
	var b_312 vm.Value
	var or__x_313 bool
	var arg__26459_320 vm.Value
	var v321 bool
	var v323 bool
	var a_324 vm.Value
	var b_325 vm.Value
	var or__x_326 vm.Value
	var a_358 vm.Value
	var b_359 vm.Value
	var or__x_370 bool
	var a_360 vm.Value
	var b_361 vm.Value
	var or__x_460 bool
	var v619 vm.Value
	var a_620 vm.Value
	var b_621 vm.Value
	var a_365 vm.Value
	var b_366 vm.Value
	var a_367 vm.Value
	var b_368 vm.Value
	var arg__26495_444 vm.Value
	var arg__26501_448 vm.Value
	var v449 vm.Value
	var v451 vm.Value
	var a_452 vm.Value
	var b_453 vm.Value
	var a_371 vm.Value
	var b_372 vm.Value
	var or__x_373 bool
	var a_374 vm.Value
	var b_375 vm.Value
	var or__x_376 bool
	var or__x_380 bool
	var v435 bool
	var a_436 vm.Value
	var b_437 vm.Value
	var or__x_438 vm.Value
	var a_381 vm.Value
	var b_382 vm.Value
	var or__x_383 bool
	var a_384 vm.Value
	var b_385 vm.Value
	var or__x_386 bool
	var or__x_390 bool
	var v430 bool
	var a_431 vm.Value
	var b_432 vm.Value
	var or__x_433 vm.Value
	var a_391 vm.Value
	var b_392 vm.Value
	var or__x_393 bool
	var a_394 vm.Value
	var b_395 vm.Value
	var or__x_396 bool
	var arg__26484_403 vm.Value
	var or__x_404 bool
	var v425 bool
	var a_426 vm.Value
	var b_427 vm.Value
	var or__x_428 vm.Value
	var a_405 vm.Value
	var b_406 vm.Value
	var or__x_407 bool
	var a_408 vm.Value
	var b_409 vm.Value
	var or__x_410 bool
	var arg__26490_417 vm.Value
	var v418 bool
	var v420 bool
	var a_421 vm.Value
	var b_422 vm.Value
	var or__x_423 vm.Value
	var a_455 vm.Value
	var b_456 vm.Value
	var a_457 vm.Value
	var b_458 vm.Value
	var v615 vm.Value
	var a_616 vm.Value
	var b_617 vm.Value
	var a_461 vm.Value
	var b_462 vm.Value
	var or__x_463 bool
	var a_464 vm.Value
	var b_465 vm.Value
	var or__x_466 bool
	var arg__26509_473 vm.Value
	var v474 bool
	var and__x_476 bool
	var a_477 vm.Value
	var b_478 vm.Value
	var or__x_479 vm.Value
	var and__x_480 bool
	var a_481 vm.Value
	var b_482 vm.Value
	var or__x_488 bool
	var and__x_483 bool
	var a_484 vm.Value
	var b_485 vm.Value
	var or__x_513 bool
	var and__x_514 vm.Value
	var a_515 vm.Value
	var b_516 vm.Value
	var and__x_489 bool
	var a_490 vm.Value
	var b_491 vm.Value
	var or__x_492 bool
	var and__x_493 bool
	var a_494 vm.Value
	var b_495 vm.Value
	var or__x_496 bool
	var arg__26517_503 vm.Value
	var v504 bool
	var v506 bool
	var and__x_507 bool
	var a_508 vm.Value
	var b_509 vm.Value
	var or__x_510 vm.Value
	var or__x_517 bool
	var a_518 vm.Value
	var b_519 vm.Value
	var or__x_520 bool
	var a_521 vm.Value
	var b_522 vm.Value
	var or__x_526 bool
	var v587 bool
	var or__x_588 vm.Value
	var a_589 vm.Value
	var b_590 vm.Value
	var a_527 vm.Value
	var b_528 vm.Value
	var or__x_529 bool
	var a_530 vm.Value
	var b_531 vm.Value
	var or__x_532 bool
	var arg__26525_539 vm.Value
	var v540 bool
	var and__x_542 bool
	var a_543 vm.Value
	var b_544 vm.Value
	var or__x_545 vm.Value
	var or__x_546 bool
	var and__x_547 bool
	var a_548 vm.Value
	var b_549 vm.Value
	var or__x_556 bool
	var or__x_550 bool
	var and__x_551 bool
	var a_552 vm.Value
	var b_553 vm.Value
	var v581 bool
	var or__x_582 bool
	var and__x_583 vm.Value
	var a_584 vm.Value
	var b_585 vm.Value
	var and__x_557 bool
	var a_558 vm.Value
	var b_559 vm.Value
	var or__x_560 bool
	var and__x_561 bool
	var a_562 vm.Value
	var b_563 vm.Value
	var or__x_564 bool
	var arg__26533_571 vm.Value
	var v572 bool
	var v574 bool
	var and__x_575 bool
	var a_576 vm.Value
	var b_577 vm.Value
	var or__x_578 vm.Value
	var a_594 vm.Value
	var b_595 vm.Value
	var arg__26538_602 vm.Value
	var arg__26544_606 vm.Value
	var v607 vm.Value
	var a_596 vm.Value
	var b_597 vm.Value
	var v611 vm.Value
	var a_612 vm.Value
	var b_613 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = v5, v10, a_6, b_7, a_8, b_9, v18, v675, a_676, b_677, a_13, b_14, a_15, b_16, v26, v671, a_672, b_673, a_21, b_22, a_23, b_24, v34, v667, a_668, b_669, a_29, b_30, a_31, b_32, v43, v663, a_664, b_665, a_38, b_39, a_40, b_41, v52, v659, a_660, b_661, a_47, b_48, a_49, b_50, v61, v655, a_656, b_657, a_56, b_57, a_58, b_59, and__x_70, v651, a_652, b_653, a_65, b_66, a_67, b_68, and__x_98, v647, a_648, b_649, a_71, b_72, and__x_73, arg__26403_82, v83, a_74, b_75, and__x_76, v86, a_87, b_88, and__x_89, a_93, b_94, a_95, b_96, and__x_126, v643, a_644, b_645, a_99, b_100, and__x_101, arg__26411_110, v111, a_102, b_103, and__x_104, v114, a_115, b_116, and__x_117, a_121, b_122, a_123, b_124, and__x_154, v639, a_640, b_641, a_127, b_128, and__x_129, arg__26419_138, v139, a_130, b_131, and__x_132, v142, a_143, b_144, and__x_145, a_149, b_150, a_151, b_152, and__x_182, v635, a_636, b_637, a_155, b_156, and__x_157, arg__26427_166, v167, a_158, b_159, and__x_160, v170, a_171, b_172, and__x_173, a_177, b_178, a_179, b_180, and__x_224, v631, a_632, b_633, a_183, b_184, and__x_185, or__x_191, a_186, b_187, and__x_188, v212, a_213, b_214, and__x_215, a_192, b_193, and__x_194, or__x_195, a_196, b_197, and__x_198, or__x_199, v203, v205, a_206, b_207, and__x_208, or__x_209, a_219, b_220, a_221, b_222, v266, v627, a_628, b_629, a_225, b_226, and__x_227, or__x_233, a_228, b_229, and__x_230, v254, a_255, b_256, and__x_257, a_234, b_235, and__x_236, or__x_237, a_238, b_239, and__x_240, or__x_241, v245, v247, a_248, b_249, and__x_250, or__x_251, a_261, b_262, or__x_273, a_263, b_264, v363, v623, a_624, b_625, a_268, b_269, a_270, b_271, arg__26464_347, arg__26470_351, v352, v354, a_355, b_356, a_274, b_275, or__x_276, a_277, b_278, or__x_279, or__x_283, v338, a_339, b_340, or__x_341, a_284, b_285, or__x_286, a_287, b_288, or__x_289, or__x_293, v333, a_334, b_335, or__x_336, a_294, b_295, or__x_296, a_297, b_298, or__x_299, arg__26453_306, or__x_307, v328, a_329, b_330, or__x_331, a_308, b_309, or__x_310, a_311, b_312, or__x_313, arg__26459_320, v321, v323, a_324, b_325, or__x_326, a_358, b_359, or__x_370, a_360, b_361, or__x_460, v619, a_620, b_621, a_365, b_366, a_367, b_368, arg__26495_444, arg__26501_448, v449, v451, a_452, b_453, a_371, b_372, or__x_373, a_374, b_375, or__x_376, or__x_380, v435, a_436, b_437, or__x_438, a_381, b_382, or__x_383, a_384, b_385, or__x_386, or__x_390, v430, a_431, b_432, or__x_433, a_391, b_392, or__x_393, a_394, b_395, or__x_396, arg__26484_403, or__x_404, v425, a_426, b_427, or__x_428, a_405, b_406, or__x_407, a_408, b_409, or__x_410, arg__26490_417, v418, v420, a_421, b_422, or__x_423, a_455, b_456, a_457, b_458, v615, a_616, b_617, a_461, b_462, or__x_463, a_464, b_465, or__x_466, arg__26509_473, v474, and__x_476, a_477, b_478, or__x_479, and__x_480, a_481, b_482, or__x_488, and__x_483, a_484, b_485, or__x_513, and__x_514, a_515, b_516, and__x_489, a_490, b_491, or__x_492, and__x_493, a_494, b_495, or__x_496, arg__26517_503, v504, v506, and__x_507, a_508, b_509, or__x_510, or__x_517, a_518, b_519, or__x_520, a_521, b_522, or__x_526, v587, or__x_588, a_589, b_590, a_527, b_528, or__x_529, a_530, b_531, or__x_532, arg__26525_539, v540, and__x_542, a_543, b_544, or__x_545, or__x_546, and__x_547, a_548, b_549, or__x_556, or__x_550, and__x_551, a_552, b_553, v581, or__x_582, and__x_583, a_584, b_585, and__x_557, a_558, b_559, or__x_560, and__x_561, a_562, b_563, or__x_564, arg__26533_571, v572, v574, and__x_575, a_576, b_577, or__x_578, a_594, b_595, arg__26538_602, arg__26544_606, v607, a_596, b_597, v611, a_612, b_613
	v5, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "ti-inc!").Deref(), []vm.Value{vm.Keyword("type-join")})
	if callErr != nil {
		return nil, callErr
	}
	v10 = arg0 == arg1
	if v10 {
		a_6 = arg0
		b_7 = arg1
		goto b1
	} else {
		a_8 = arg0
		b_9 = arg1
		goto b2
	}
b1:
	;
	v675 = a_6
	a_676 = a_6
	b_677 = b_7
	goto b3
b2:
	;
	v18 = a_8 == vm.Keyword("bottom")
	if v18 {
		a_13 = a_8
		b_14 = b_9
		goto b4
	} else {
		a_15 = a_8
		b_16 = b_9
		goto b5
	}
b3:
	;
	return v675, nil
b4:
	;
	v671 = b_14
	a_672 = a_13
	b_673 = b_14
	goto b6
b5:
	;
	v26 = b_16 == vm.Keyword("bottom")
	if v26 {
		a_21 = a_15
		b_22 = b_16
		goto b7
	} else {
		a_23 = a_15
		b_24 = b_16
		goto b8
	}
b6:
	;
	v675 = v671
	a_676 = a_672
	b_677 = b_673
	goto b3
b7:
	;
	v667 = a_21
	a_668 = a_21
	b_669 = b_22
	goto b9
b8:
	;
	v34 = a_23 == vm.Keyword("any")
	if v34 {
		a_29 = a_23
		b_30 = b_24
		goto b10
	} else {
		a_31 = a_23
		b_32 = b_24
		goto b11
	}
b9:
	;
	v671 = v667
	a_672 = a_668
	b_673 = b_669
	goto b6
b10:
	;
	v663 = vm.Keyword("any")
	a_664 = a_29
	b_665 = b_30
	goto b12
b11:
	;
	v43 = b_32 == vm.Keyword("any")
	if v43 {
		a_38 = a_31
		b_39 = b_32
		goto b13
	} else {
		a_40 = a_31
		b_41 = b_32
		goto b14
	}
b12:
	;
	v667 = v663
	a_668 = a_664
	b_669 = b_665
	goto b9
b13:
	;
	v659 = vm.Keyword("any")
	a_660 = a_38
	b_661 = b_39
	goto b15
b14:
	;
	v52 = a_40 == vm.Keyword("unknown")
	if v52 {
		a_47 = a_40
		b_48 = b_41
		goto b16
	} else {
		a_49 = a_40
		b_50 = b_41
		goto b17
	}
b15:
	;
	v663 = v659
	a_664 = a_660
	b_665 = b_661
	goto b12
b16:
	;
	v655 = vm.Keyword("unknown")
	a_656 = a_47
	b_657 = b_48
	goto b18
b17:
	;
	v61 = b_50 == vm.Keyword("unknown")
	if v61 {
		a_56 = a_49
		b_57 = b_50
		goto b19
	} else {
		a_58 = a_49
		b_59 = b_50
		goto b20
	}
b18:
	;
	v659 = v655
	a_660 = a_656
	b_661 = b_657
	goto b15
b19:
	;
	v651 = vm.Keyword("unknown")
	a_652 = a_56
	b_653 = b_57
	goto b21
b20:
	;
	and__x_70 = a_58 == vm.Keyword("int")
	if and__x_70 {
		a_71 = a_58
		b_72 = b_59
		and__x_73 = and__x_70
		goto b25
	} else {
		a_74 = a_58
		b_75 = b_59
		and__x_76 = and__x_70
		goto b26
	}
b21:
	;
	v655 = v651
	a_656 = a_652
	b_657 = b_653
	goto b18
b22:
	;
	v647 = vm.Keyword("int")
	a_648 = a_65
	b_649 = b_66
	goto b24
b23:
	;
	and__x_98 = b_68 == vm.Keyword("int")
	if and__x_98 {
		a_99 = a_67
		b_100 = b_68
		and__x_101 = and__x_98
		goto b31
	} else {
		a_102 = a_67
		b_103 = b_68
		and__x_104 = and__x_98
		goto b32
	}
b24:
	;
	v651 = v647
	a_652 = a_648
	b_653 = b_649
	goto b21
b25:
	;
	arg__26403_82, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("const"), vm.Keyword("int"), vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	v83 = b_72 == arg__26403_82
	v86 = v83
	a_87 = a_71
	b_88 = b_72
	and__x_89 = vm.Boolean(and__x_73)
	goto b27
b26:
	;
	v86 = and__x_76
	a_87 = a_74
	b_88 = b_75
	and__x_89 = vm.Boolean(and__x_76)
	goto b27
b27:
	;
	if v86 {
		a_65 = a_87
		b_66 = b_88
		goto b22
	} else {
		a_67 = a_87
		b_68 = b_88
		goto b23
	}
b28:
	;
	v643 = vm.Keyword("int")
	a_644 = a_93
	b_645 = b_94
	goto b30
b29:
	;
	and__x_126 = a_95 == vm.Keyword("float")
	if and__x_126 {
		a_127 = a_95
		b_128 = b_96
		and__x_129 = and__x_126
		goto b37
	} else {
		a_130 = a_95
		b_131 = b_96
		and__x_132 = and__x_126
		goto b38
	}
b30:
	;
	v647 = v643
	a_648 = a_644
	b_649 = b_645
	goto b24
b31:
	;
	arg__26411_110, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("const"), vm.Keyword("int"), vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	v111 = a_99 == arg__26411_110
	v114 = v111
	a_115 = a_99
	b_116 = b_100
	and__x_117 = vm.Boolean(and__x_101)
	goto b33
b32:
	;
	v114 = and__x_104
	a_115 = a_102
	b_116 = b_103
	and__x_117 = vm.Boolean(and__x_104)
	goto b33
b33:
	;
	if v114 {
		a_93 = a_115
		b_94 = b_116
		goto b28
	} else {
		a_95 = a_115
		b_96 = b_116
		goto b29
	}
b34:
	;
	v639 = vm.Keyword("float")
	a_640 = a_121
	b_641 = b_122
	goto b36
b35:
	;
	and__x_154 = b_124 == vm.Keyword("float")
	if and__x_154 {
		a_155 = a_123
		b_156 = b_124
		and__x_157 = and__x_154
		goto b43
	} else {
		a_158 = a_123
		b_159 = b_124
		and__x_160 = and__x_154
		goto b44
	}
b36:
	;
	v643 = v639
	a_644 = a_640
	b_645 = b_641
	goto b30
b37:
	;
	arg__26419_138, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("const"), vm.Keyword("float"), vm.Float(0)})
	if callErr != nil {
		return nil, callErr
	}
	v139 = b_128 == arg__26419_138
	v142 = v139
	a_143 = a_127
	b_144 = b_128
	and__x_145 = vm.Boolean(and__x_129)
	goto b39
b38:
	;
	v142 = and__x_132
	a_143 = a_130
	b_144 = b_131
	and__x_145 = vm.Boolean(and__x_132)
	goto b39
b39:
	;
	if v142 {
		a_121 = a_143
		b_122 = b_144
		goto b34
	} else {
		a_123 = a_143
		b_124 = b_144
		goto b35
	}
b40:
	;
	v635 = vm.Keyword("float")
	a_636 = a_149
	b_637 = b_150
	goto b42
b41:
	;
	and__x_182 = a_151 == vm.Keyword("bool")
	if and__x_182 {
		a_183 = a_151
		b_184 = b_152
		and__x_185 = and__x_182
		goto b49
	} else {
		a_186 = a_151
		b_187 = b_152
		and__x_188 = and__x_182
		goto b50
	}
b42:
	;
	v639 = v635
	a_640 = a_636
	b_641 = b_637
	goto b36
b43:
	;
	arg__26427_166, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("const"), vm.Keyword("float"), vm.Float(0)})
	if callErr != nil {
		return nil, callErr
	}
	v167 = a_155 == arg__26427_166
	v170 = v167
	a_171 = a_155
	b_172 = b_156
	and__x_173 = vm.Boolean(and__x_157)
	goto b45
b44:
	;
	v170 = and__x_160
	a_171 = a_158
	b_172 = b_159
	and__x_173 = vm.Boolean(and__x_160)
	goto b45
b45:
	;
	if v170 {
		a_149 = a_171
		b_150 = b_172
		goto b40
	} else {
		a_151 = a_171
		b_152 = b_172
		goto b41
	}
b46:
	;
	v631 = vm.Keyword("bool")
	a_632 = a_177
	b_633 = b_178
	goto b48
b47:
	;
	and__x_224 = b_180 == vm.Keyword("bool")
	if and__x_224 {
		a_225 = a_179
		b_226 = b_180
		and__x_227 = and__x_224
		goto b58
	} else {
		a_228 = a_179
		b_229 = b_180
		and__x_230 = and__x_224
		goto b59
	}
b48:
	;
	v635 = v631
	a_636 = a_632
	b_637 = b_633
	goto b42
b49:
	;
	or__x_191 = b_184 == vm.Keyword("true")
	if or__x_191 {
		a_192 = a_183
		b_193 = b_184
		and__x_194 = and__x_185
		or__x_195 = or__x_191
		goto b52
	} else {
		a_196 = a_183
		b_197 = b_184
		and__x_198 = and__x_185
		or__x_199 = or__x_191
		goto b53
	}
b50:
	;
	v212 = and__x_188
	a_213 = a_186
	b_214 = b_187
	and__x_215 = vm.Boolean(and__x_188)
	goto b51
b51:
	;
	if v212 {
		a_177 = a_213
		b_178 = b_214
		goto b46
	} else {
		a_179 = a_213
		b_180 = b_214
		goto b47
	}
b52:
	;
	v205 = or__x_195
	a_206 = a_192
	b_207 = b_193
	and__x_208 = and__x_194
	or__x_209 = vm.Boolean(or__x_195)
	goto b54
b53:
	;
	v203 = b_197 == vm.Keyword("false")
	v205 = v203
	a_206 = a_196
	b_207 = b_197
	and__x_208 = and__x_198
	or__x_209 = vm.Boolean(or__x_199)
	goto b54
b54:
	;
	v212 = v205
	a_213 = a_206
	b_214 = b_207
	and__x_215 = vm.Boolean(and__x_208)
	goto b51
b55:
	;
	v627 = vm.Keyword("bool")
	a_628 = a_219
	b_629 = b_220
	goto b57
b56:
	;
	v266 = a_221 == vm.Keyword("number")
	if v266 {
		a_261 = a_221
		b_262 = b_222
		goto b64
	} else {
		a_263 = a_221
		b_264 = b_222
		goto b65
	}
b57:
	;
	v631 = v627
	a_632 = a_628
	b_633 = b_629
	goto b48
b58:
	;
	or__x_233 = a_225 == vm.Keyword("true")
	if or__x_233 {
		a_234 = a_225
		b_235 = b_226
		and__x_236 = and__x_227
		or__x_237 = or__x_233
		goto b61
	} else {
		a_238 = a_225
		b_239 = b_226
		and__x_240 = and__x_227
		or__x_241 = or__x_233
		goto b62
	}
b59:
	;
	v254 = and__x_230
	a_255 = a_228
	b_256 = b_229
	and__x_257 = vm.Boolean(and__x_230)
	goto b60
b60:
	;
	if v254 {
		a_219 = a_255
		b_220 = b_256
		goto b55
	} else {
		a_221 = a_255
		b_222 = b_256
		goto b56
	}
b61:
	;
	v247 = or__x_237
	a_248 = a_234
	b_249 = b_235
	and__x_250 = and__x_236
	or__x_251 = vm.Boolean(or__x_237)
	goto b63
b62:
	;
	v245 = a_238 == vm.Keyword("false")
	v247 = v245
	a_248 = a_238
	b_249 = b_239
	and__x_250 = and__x_240
	or__x_251 = vm.Boolean(or__x_241)
	goto b63
b63:
	;
	v254 = v247
	a_255 = a_248
	b_256 = b_249
	and__x_257 = vm.Boolean(and__x_250)
	goto b60
b64:
	;
	or__x_273 = b_262 == vm.Keyword("int")
	if or__x_273 {
		a_274 = a_261
		b_275 = b_262
		or__x_276 = or__x_273
		goto b70
	} else {
		a_277 = a_261
		b_278 = b_262
		or__x_279 = or__x_273
		goto b71
	}
b65:
	;
	v363 = b_264 == vm.Keyword("number")
	if v363 {
		a_358 = a_263
		b_359 = b_264
		goto b82
	} else {
		a_360 = a_263
		b_361 = b_264
		goto b83
	}
b66:
	;
	v627 = v623
	a_628 = a_624
	b_629 = b_625
	goto b57
b67:
	;
	v354 = vm.Keyword("number")
	a_355 = a_268
	b_356 = b_269
	goto b69
b68:
	;
	arg__26464_347, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("union"), a_270, b_271})
	if callErr != nil {
		return nil, callErr
	}
	arg__26470_351, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("union"), a_270, b_271})
	if callErr != nil {
		return nil, callErr
	}
	v352, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "normalize-type").Deref(), []vm.Value{arg__26470_351})
	if callErr != nil {
		return nil, callErr
	}
	v354 = v352
	a_355 = a_270
	b_356 = b_271
	goto b69
b69:
	;
	v623 = v354
	a_624 = a_355
	b_625 = b_356
	goto b66
b70:
	;
	v338 = or__x_276
	a_339 = a_274
	b_340 = b_275
	or__x_341 = vm.Boolean(or__x_276)
	goto b72
b71:
	;
	or__x_283 = b_278 == vm.Keyword("float")
	if or__x_283 {
		a_284 = a_277
		b_285 = b_278
		or__x_286 = or__x_283
		goto b73
	} else {
		a_287 = a_277
		b_288 = b_278
		or__x_289 = or__x_283
		goto b74
	}
b72:
	;
	if v338 {
		a_268 = a_339
		b_269 = b_340
		goto b67
	} else {
		a_270 = a_339
		b_271 = b_340
		goto b68
	}
b73:
	;
	v333 = or__x_286
	a_334 = a_284
	b_335 = b_285
	or__x_336 = vm.Boolean(or__x_286)
	goto b75
b74:
	;
	or__x_293 = b_288 == vm.Keyword("number")
	if or__x_293 {
		a_294 = a_287
		b_295 = b_288
		or__x_296 = or__x_293
		goto b76
	} else {
		a_297 = a_287
		b_298 = b_288
		or__x_299 = or__x_293
		goto b77
	}
b75:
	;
	v338 = v333
	a_339 = a_334
	b_340 = b_335
	or__x_341 = vm.Boolean(or__x_279)
	goto b72
b76:
	;
	v328 = or__x_296
	a_329 = a_294
	b_330 = b_295
	or__x_331 = vm.Boolean(or__x_296)
	goto b78
b77:
	;
	arg__26453_306, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("const"), vm.Keyword("int"), vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	or__x_307 = b_298 == arg__26453_306
	if or__x_307 {
		a_308 = a_297
		b_309 = b_298
		or__x_310 = or__x_307
		goto b79
	} else {
		a_311 = a_297
		b_312 = b_298
		or__x_313 = or__x_307
		goto b80
	}
b78:
	;
	v333 = v328
	a_334 = a_329
	b_335 = b_330
	or__x_336 = vm.Boolean(or__x_289)
	goto b75
b79:
	;
	v323 = or__x_310
	a_324 = a_308
	b_325 = b_309
	or__x_326 = vm.Boolean(or__x_310)
	goto b81
b80:
	;
	arg__26459_320, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("const"), vm.Keyword("float"), vm.Float(0)})
	if callErr != nil {
		return nil, callErr
	}
	v321 = b_312 == arg__26459_320
	v323 = v321
	a_324 = a_311
	b_325 = b_312
	or__x_326 = vm.Boolean(or__x_313)
	goto b81
b81:
	;
	v328 = v323
	a_329 = a_324
	b_330 = b_325
	or__x_331 = vm.Boolean(or__x_299)
	goto b78
b82:
	;
	or__x_370 = a_358 == vm.Keyword("int")
	if or__x_370 {
		a_371 = a_358
		b_372 = b_359
		or__x_373 = or__x_370
		goto b88
	} else {
		a_374 = a_358
		b_375 = b_359
		or__x_376 = or__x_370
		goto b89
	}
b83:
	;
	or__x_460 = a_360 == vm.Keyword("int")
	if or__x_460 {
		a_461 = a_360
		b_462 = b_361
		or__x_463 = or__x_460
		goto b103
	} else {
		a_464 = a_360
		b_465 = b_361
		or__x_466 = or__x_460
		goto b104
	}
b84:
	;
	v623 = v619
	a_624 = a_620
	b_625 = b_621
	goto b66
b85:
	;
	v451 = vm.Keyword("number")
	a_452 = a_365
	b_453 = b_366
	goto b87
b86:
	;
	arg__26495_444, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("union"), a_367, b_368})
	if callErr != nil {
		return nil, callErr
	}
	arg__26501_448, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("union"), a_367, b_368})
	if callErr != nil {
		return nil, callErr
	}
	v449, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "normalize-type").Deref(), []vm.Value{arg__26501_448})
	if callErr != nil {
		return nil, callErr
	}
	v451 = v449
	a_452 = a_367
	b_453 = b_368
	goto b87
b87:
	;
	v619 = v451
	a_620 = a_452
	b_621 = b_453
	goto b84
b88:
	;
	v435 = or__x_373
	a_436 = a_371
	b_437 = b_372
	or__x_438 = vm.Boolean(or__x_373)
	goto b90
b89:
	;
	or__x_380 = a_374 == vm.Keyword("float")
	if or__x_380 {
		a_381 = a_374
		b_382 = b_375
		or__x_383 = or__x_380
		goto b91
	} else {
		a_384 = a_374
		b_385 = b_375
		or__x_386 = or__x_380
		goto b92
	}
b90:
	;
	if v435 {
		a_365 = a_436
		b_366 = b_437
		goto b85
	} else {
		a_367 = a_436
		b_368 = b_437
		goto b86
	}
b91:
	;
	v430 = or__x_383
	a_431 = a_381
	b_432 = b_382
	or__x_433 = vm.Boolean(or__x_383)
	goto b93
b92:
	;
	or__x_390 = a_384 == vm.Keyword("number")
	if or__x_390 {
		a_391 = a_384
		b_392 = b_385
		or__x_393 = or__x_390
		goto b94
	} else {
		a_394 = a_384
		b_395 = b_385
		or__x_396 = or__x_390
		goto b95
	}
b93:
	;
	v435 = v430
	a_436 = a_431
	b_437 = b_432
	or__x_438 = vm.Boolean(or__x_376)
	goto b90
b94:
	;
	v425 = or__x_393
	a_426 = a_391
	b_427 = b_392
	or__x_428 = vm.Boolean(or__x_393)
	goto b96
b95:
	;
	arg__26484_403, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("const"), vm.Keyword("int"), vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	or__x_404 = a_394 == arg__26484_403
	if or__x_404 {
		a_405 = a_394
		b_406 = b_395
		or__x_407 = or__x_404
		goto b97
	} else {
		a_408 = a_394
		b_409 = b_395
		or__x_410 = or__x_404
		goto b98
	}
b96:
	;
	v430 = v425
	a_431 = a_426
	b_432 = b_427
	or__x_433 = vm.Boolean(or__x_386)
	goto b93
b97:
	;
	v420 = or__x_407
	a_421 = a_405
	b_422 = b_406
	or__x_423 = vm.Boolean(or__x_407)
	goto b99
b98:
	;
	arg__26490_417, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("const"), vm.Keyword("float"), vm.Float(0)})
	if callErr != nil {
		return nil, callErr
	}
	v418 = a_408 == arg__26490_417
	v420 = v418
	a_421 = a_408
	b_422 = b_409
	or__x_423 = vm.Boolean(or__x_410)
	goto b99
b99:
	;
	v425 = v420
	a_426 = a_421
	b_427 = b_422
	or__x_428 = vm.Boolean(or__x_396)
	goto b96
b100:
	;
	v615 = vm.Keyword("number")
	a_616 = a_455
	b_617 = b_456
	goto b102
b101:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		a_594 = a_457
		b_595 = b_458
		goto b124
	} else {
		a_596 = a_457
		b_597 = b_458
		goto b125
	}
b102:
	;
	v619 = v615
	a_620 = a_616
	b_621 = b_617
	goto b84
b103:
	;
	and__x_476 = or__x_463
	a_477 = a_461
	b_478 = b_462
	or__x_479 = vm.Boolean(or__x_463)
	goto b105
b104:
	;
	arg__26509_473, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("const"), vm.Keyword("int"), vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	v474 = a_464 == arg__26509_473
	and__x_476 = v474
	a_477 = a_464
	b_478 = b_465
	or__x_479 = vm.Boolean(or__x_466)
	goto b105
b105:
	;
	if and__x_476 {
		and__x_480 = and__x_476
		a_481 = a_477
		b_482 = b_478
		goto b106
	} else {
		and__x_483 = and__x_476
		a_484 = a_477
		b_485 = b_478
		goto b107
	}
b106:
	;
	or__x_488 = b_482 == vm.Keyword("float")
	if or__x_488 {
		and__x_489 = and__x_480
		a_490 = a_481
		b_491 = b_482
		or__x_492 = or__x_488
		goto b109
	} else {
		and__x_493 = and__x_480
		a_494 = a_481
		b_495 = b_482
		or__x_496 = or__x_488
		goto b110
	}
b107:
	;
	or__x_513 = and__x_483
	and__x_514 = vm.Boolean(and__x_483)
	a_515 = a_484
	b_516 = b_485
	goto b108
b108:
	;
	if or__x_513 {
		or__x_517 = or__x_513
		a_518 = a_515
		b_519 = b_516
		goto b112
	} else {
		or__x_520 = or__x_513
		a_521 = a_515
		b_522 = b_516
		goto b113
	}
b109:
	;
	v506 = or__x_492
	and__x_507 = and__x_489
	a_508 = a_490
	b_509 = b_491
	or__x_510 = vm.Boolean(or__x_492)
	goto b111
b110:
	;
	arg__26517_503, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("const"), vm.Keyword("float"), vm.Float(0)})
	if callErr != nil {
		return nil, callErr
	}
	v504 = b_495 == arg__26517_503
	v506 = v504
	and__x_507 = and__x_493
	a_508 = a_494
	b_509 = b_495
	or__x_510 = vm.Boolean(or__x_496)
	goto b111
b111:
	;
	or__x_513 = v506
	and__x_514 = vm.Boolean(and__x_507)
	a_515 = a_508
	b_516 = b_509
	goto b108
b112:
	;
	v587 = or__x_517
	or__x_588 = vm.Boolean(or__x_517)
	a_589 = a_518
	b_590 = b_519
	goto b114
b113:
	;
	or__x_526 = a_521 == vm.Keyword("float")
	if or__x_526 {
		a_527 = a_521
		b_528 = b_522
		or__x_529 = or__x_526
		goto b115
	} else {
		a_530 = a_521
		b_531 = b_522
		or__x_532 = or__x_526
		goto b116
	}
b114:
	;
	if v587 {
		a_455 = a_589
		b_456 = b_590
		goto b100
	} else {
		a_457 = a_589
		b_458 = b_590
		goto b101
	}
b115:
	;
	and__x_542 = or__x_529
	a_543 = a_527
	b_544 = b_528
	or__x_545 = vm.Boolean(or__x_529)
	goto b117
b116:
	;
	arg__26525_539, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("const"), vm.Keyword("float"), vm.Float(0)})
	if callErr != nil {
		return nil, callErr
	}
	v540 = a_530 == arg__26525_539
	and__x_542 = v540
	a_543 = a_530
	b_544 = b_531
	or__x_545 = vm.Boolean(or__x_532)
	goto b117
b117:
	;
	if and__x_542 {
		or__x_546 = or__x_520
		and__x_547 = and__x_542
		a_548 = a_543
		b_549 = b_544
		goto b118
	} else {
		or__x_550 = or__x_520
		and__x_551 = and__x_542
		a_552 = a_543
		b_553 = b_544
		goto b119
	}
b118:
	;
	or__x_556 = b_549 == vm.Keyword("int")
	if or__x_556 {
		and__x_557 = and__x_547
		a_558 = a_548
		b_559 = b_549
		or__x_560 = or__x_556
		goto b121
	} else {
		and__x_561 = and__x_547
		a_562 = a_548
		b_563 = b_549
		or__x_564 = or__x_556
		goto b122
	}
b119:
	;
	v581 = and__x_551
	or__x_582 = or__x_550
	and__x_583 = vm.Boolean(and__x_551)
	a_584 = a_552
	b_585 = b_553
	goto b120
b120:
	;
	v587 = v581
	or__x_588 = vm.Boolean(or__x_582)
	a_589 = a_584
	b_590 = b_585
	goto b114
b121:
	;
	v574 = or__x_560
	and__x_575 = and__x_557
	a_576 = a_558
	b_577 = b_559
	or__x_578 = vm.Boolean(or__x_560)
	goto b123
b122:
	;
	arg__26533_571, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("const"), vm.Keyword("int"), vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	v572 = b_563 == arg__26533_571
	v574 = v572
	and__x_575 = and__x_561
	a_576 = a_562
	b_577 = b_563
	or__x_578 = vm.Boolean(or__x_564)
	goto b123
b123:
	;
	v581 = v574
	or__x_582 = or__x_546
	and__x_583 = vm.Boolean(and__x_575)
	a_584 = a_576
	b_585 = b_577
	goto b120
b124:
	;
	arg__26538_602, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("union"), a_594, b_595})
	if callErr != nil {
		return nil, callErr
	}
	arg__26544_606, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("union"), a_594, b_595})
	if callErr != nil {
		return nil, callErr
	}
	v607, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "normalize-type").Deref(), []vm.Value{arg__26544_606})
	if callErr != nil {
		return nil, callErr
	}
	v611 = v607
	a_612 = a_594
	b_613 = b_595
	goto b126
b125:
	;
	v611 = vm.NIL
	a_612 = a_596
	b_613 = b_597
	goto b126
b126:
	;
	v615 = v611
	a_616 = a_612
	b_617 = b_613
	goto b102
}
func set_type_if_changed_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var new_type_4 vm.Value
	var old_type_6 vm.Value
	var v16 bool
	var f_7 vm.Value
	var inst_8 vm.Value
	var new_type_9 vm.Value
	var old_type_10 vm.Value
	var f_11 vm.Value
	var inst_12 vm.Value
	var new_type_13 vm.Value
	var old_type_14 vm.Value
	var v20 vm.Value
	var joined_22 vm.Value
	var f_23 vm.Value
	var inst_24 vm.Value
	var new_type_25 vm.Value
	var old_type_26 vm.Value
	var v37 bool
	var joined_27 vm.Value
	var f_28 vm.Value
	var inst_29 vm.Value
	var new_type_30 vm.Value
	var old_type_31 vm.Value
	var joined_32 vm.Value
	var f_33 vm.Value
	var inst_34 vm.Value
	var new_type_35 vm.Value
	var old_type_36 vm.Value
	var v42 vm.Value
	var v45 vm.Value
	var joined_46 vm.Value
	var f_47 vm.Value
	var inst_48 vm.Value
	var new_type_49 vm.Value
	var old_type_50 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = new_type_4, old_type_6, v16, f_7, inst_8, new_type_9, old_type_10, f_11, inst_12, new_type_13, old_type_14, v20, joined_22, f_23, inst_24, new_type_25, old_type_26, v37, joined_27, f_28, inst_29, new_type_30, old_type_31, joined_32, f_33, inst_34, new_type_35, old_type_36, v42, v45, joined_46, f_47, inst_48, new_type_49, old_type_50
	new_type_4, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "normalize-type").Deref(), []vm.Value{arg2})
	if callErr != nil {
		return nil, callErr
	}
	old_type_6, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg1, arg0})
	if callErr != nil {
		return nil, callErr
	}
	v16 = old_type_6 == vm.Keyword("bottom")
	if v16 {
		f_7 = arg0
		inst_8 = arg1
		new_type_9 = new_type_4
		old_type_10 = old_type_6
		goto b1
	} else {
		f_11 = arg0
		inst_12 = arg1
		new_type_13 = new_type_4
		old_type_14 = old_type_6
		goto b2
	}
b1:
	;
	joined_22 = new_type_9
	f_23 = f_7
	inst_24 = inst_8
	new_type_25 = new_type_9
	old_type_26 = old_type_10
	goto b3
b2:
	;
	v20, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "type-join").Deref(), []vm.Value{old_type_14, new_type_13})
	if callErr != nil {
		return nil, callErr
	}
	joined_22 = v20
	f_23 = f_11
	inst_24 = inst_12
	new_type_25 = new_type_13
	old_type_26 = old_type_14
	goto b3
b3:
	;
	v37 = old_type_26 == joined_22
	if v37 {
		joined_27 = joined_22
		f_28 = f_23
		inst_29 = inst_24
		new_type_30 = new_type_25
		old_type_31 = old_type_26
		goto b4
	} else {
		joined_32 = joined_22
		f_33 = f_23
		inst_34 = inst_24
		new_type_35 = new_type_25
		old_type_36 = old_type_26
		goto b5
	}
b4:
	;
	v45 = vm.NIL
	joined_46 = joined_27
	f_47 = f_28
	inst_48 = inst_29
	new_type_49 = new_type_30
	old_type_50 = old_type_31
	goto b6
b5:
	;
	v42, callErr = rt.InvokeValue(rt.LookupVar("ir", "set-type!").Deref(), []vm.Value{f_33, inst_34, joined_32})
	if callErr != nil {
		return nil, callErr
	}
	v45 = vm.Boolean(true)
	joined_46 = joined_32
	f_47 = f_33
	inst_48 = inst_34
	new_type_49 = new_type_35
	old_type_50 = old_type_36
	goto b6
b6:
	;
	return v45, nil
}
func source_arg_type(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var pred_bid_5 vm.Value
	var target_bid_9 vm.Value
	var arg_id_13 vm.Value
	var arg_type_15 vm.Value
	var arg__26600_17 vm.Value
	var arg__26613_20 vm.Value
	var v21 vm.Value
	var callErr error
	_, _, _, _, _, _, _ = pred_bid_5, target_bid_9, arg_id_13, arg_type_15, arg__26600_17, arg__26613_20, v21
	pred_bid_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	target_bid_9, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(1)})
	if callErr != nil {
		return nil, callErr
	}
	arg_id_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg0, vm.Int(2)})
	if callErr != nil {
		return nil, callErr
	}
	arg_type_15, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg_id_13, arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__26600_17, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "refine-edge-arg-type").Deref(), []vm.Value{pred_bid_5, target_bid_9, arg_id_13, arg_type_15, arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__26613_20, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "refine-edge-arg-type").Deref(), []vm.Value{pred_bid_5, target_bid_9, arg_id_13, arg_type_15, arg1})
	if callErr != nil {
		return nil, callErr
	}
	v21, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "normalize-type").Deref(), []vm.Value{arg__26613_20})
	if callErr != nil {
		return nil, callErr
	}
	return v21, nil
}
func infer_block_param_from_deps(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value, arg3 vm.Value) (vm.Value, error) {
	var arg__26616_5 vm.Value
	var arg__26620_7 vm.Value
	var arg__26624_10 vm.Value
	var arg__26628_12 vm.Value
	var sources_13 vm.Value
	var v25 vm.Value
	var deps_14 vm.Value
	var bid_15 vm.Value
	var param_pos_16 vm.Value
	var f_17 vm.Value
	var sources_18 vm.Value
	var arg__26647_34 vm.Value
	var arg__26664_43 vm.Value
	var v44 vm.Value
	var deps_19 vm.Value
	var bid_20 vm.Value
	var param_pos_21 vm.Value
	var f_22 vm.Value
	var sources_23 vm.Value
	var v48 vm.Value
	var deps_49 vm.Value
	var bid_50 vm.Value
	var param_pos_51 vm.Value
	var f_52 vm.Value
	var sources_53 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__26616_5, arg__26620_7, arg__26624_10, arg__26628_12, sources_13, v25, deps_14, bid_15, param_pos_16, f_17, sources_18, arg__26647_34, arg__26664_43, v44, deps_19, bid_20, param_pos_21, f_22, sources_23, v48, deps_49, bid_50, param_pos_51, f_52, sources_53
	arg__26616_5, callErr = rt.InvokeValue(vm.Keyword("param-sources"), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__26620_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg1, arg2})
	if callErr != nil {
		return nil, callErr
	}
	arg__26624_10, callErr = rt.InvokeValue(vm.Keyword("param-sources"), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__26628_12, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg1, arg2})
	if callErr != nil {
		return nil, callErr
	}
	sources_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg__26624_10, arg__26628_12})
	if callErr != nil {
		return nil, callErr
	}
	v25, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{sources_13})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v25) {
		deps_14 = arg0
		bid_15 = arg1
		param_pos_16 = arg2
		f_17 = arg3
		sources_18 = sources_13
		goto b1
	} else {
		deps_19 = arg0
		bid_20 = arg1
		param_pos_21 = arg2
		f_22 = arg3
		sources_23 = sources_13
		goto b2
	}
b1:
	;
	arg__26647_34, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "source-arg-type").Deref(), []vm.Value{arg0, f_17})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), sources_18})
	if callErr != nil {
		return nil, callErr
	}
	arg__26664_43, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "source-arg-type").Deref(), []vm.Value{arg0, f_17})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), sources_18})
	if callErr != nil {
		return nil, callErr
	}
	v44, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "join-all").Deref(), []vm.Value{arg__26664_43})
	if callErr != nil {
		return nil, callErr
	}
	v48 = v44
	deps_49 = deps_14
	bid_50 = bid_15
	param_pos_51 = param_pos_16
	f_52 = f_17
	sources_53 = sources_18
	goto b3
b2:
	;
	v48 = vm.Keyword("bottom")
	deps_49 = deps_19
	bid_50 = bid_20
	param_pos_51 = param_pos_21
	f_52 = f_22
	sources_53 = sources_23
	goto b3
b3:
	;
	return v48, nil
}
func infer_block_param(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var preds_4 vm.Value
	var v14 vm.Value
	var bid_5 vm.Value
	var param_pos_6 vm.Value
	var f_7 vm.Value
	var preds_8 vm.Value
	var arg__26696_27 vm.Value
	var arg__26721_40 vm.Value
	var v41 vm.Value
	var bid_9 vm.Value
	var param_pos_10 vm.Value
	var f_11 vm.Value
	var preds_12 vm.Value
	var v45 vm.Value
	var bid_46 vm.Value
	var param_pos_47 vm.Value
	var f_48 vm.Value
	var preds_49 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = preds_4, v14, bid_5, param_pos_6, f_7, preds_8, arg__26696_27, arg__26721_40, v41, bid_9, param_pos_10, f_11, preds_12, v45, bid_46, param_pos_47, f_48, preds_49
	preds_4, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-preds").Deref(), []vm.Value{arg0, arg2})
	if callErr != nil {
		return nil, callErr
	}
	v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{preds_4})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v14) {
		bid_5 = arg0
		param_pos_6 = arg1
		f_7 = arg2
		preds_8 = preds_4
		goto b1
	} else {
		bid_9 = arg0
		param_pos_10 = arg1
		f_11 = arg2
		preds_12 = preds_4
		goto b2
	}
b1:
	;
	arg__26696_27, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapcat").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v5 vm.Value
		var callErr error
		_ = v5
		v5, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "target-arg-types").Deref(), []vm.Value{arg0, bid_5, param_pos_6, f_7})
		if callErr != nil {
			return nil, callErr
		}
		return v5, nil
	}), preds_8})
	if callErr != nil {
		return nil, callErr
	}
	arg__26721_40, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapcat").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v5 vm.Value
		var callErr error
		_ = v5
		v5, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "target-arg-types").Deref(), []vm.Value{arg0, bid_5, param_pos_6, f_7})
		if callErr != nil {
			return nil, callErr
		}
		return v5, nil
	}), preds_8})
	if callErr != nil {
		return nil, callErr
	}
	v41, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "join-all").Deref(), []vm.Value{arg__26721_40})
	if callErr != nil {
		return nil, callErr
	}
	v45 = v41
	bid_46 = bid_5
	param_pos_47 = param_pos_6
	f_48 = f_7
	preds_49 = preds_8
	goto b3
b2:
	;
	v45 = vm.Keyword("bottom")
	bid_46 = bid_9
	param_pos_47 = param_pos_10
	f_48 = f_11
	preds_49 = preds_12
	goto b3
b3:
	;
	return v45, nil
}
func build_deps(arg0 vm.Value) (vm.Value, error) {
	var bu_4 vm.Value
	var ps_8 vm.Value
	var arg__26794_16 vm.Value
	var arg__26799_19 vm.Value
	var doseq_seq__26722_20 vm.Value
	var doseq_loop__26723_21 vm.Value
	var f_22 vm.Value
	var add_branch_args_23 vm.Value
	var v276 vm.Value
	var v291 vm.Value
	var v306 vm.Value
	var bu_25 vm.Value
	var ps_26 vm.Value
	var add_27 vm.Value
	var doseq_seq__26722_28 vm.Value
	var doseq_loop__26723_29 vm.Value
	var f_30 vm.Value
	var add_branch_args_31 vm.Value
	var v275 vm.Value
	var v290 vm.Value
	var v305 vm.Value
	var bid_41 vm.Value
	var term_43 vm.Value
	var v63 vm.Value
	var bu_32 vm.Value
	var ps_33 vm.Value
	var add_34 vm.Value
	var doseq_seq__26722_35 vm.Value
	var doseq_loop__26723_36 vm.Value
	var f_37 vm.Value
	var add_branch_args_38 vm.Value
	var v285 vm.Value
	var v300 vm.Value
	var v315 vm.Value
	var v245 vm.Value
	var bu_246 vm.Value
	var ps_247 vm.Value
	var add_248 vm.Value
	var doseq_seq__26722_249 vm.Value
	var doseq_loop__26723_250 vm.Value
	var f_251 vm.Value
	var add_branch_args_252 vm.Value
	var arg__26896_256 vm.Value
	var arg__26901_259 vm.Value
	var arg__26906_262 vm.Value
	var v263 vm.Value
	var bu_44 vm.Value
	var ps_45 vm.Value
	var add_46 vm.Value
	var doseq_seq__26722_47 vm.Value
	var doseq_loop__26723_48 vm.Value
	var f_49 vm.Value
	var add_branch_args_50 vm.Value
	var bid_51 vm.Value
	var term_52 vm.Value
	var v281 vm.Value
	var v296 vm.Value
	var v311 vm.Value
	var bu_53 vm.Value
	var ps_54 vm.Value
	var add_55 vm.Value
	var doseq_seq__26722_56 vm.Value
	var doseq_loop__26723_57 vm.Value
	var f_58 vm.Value
	var add_branch_args_59 vm.Value
	var bid_60 vm.Value
	var term_61 vm.Value
	var v279 vm.Value
	var v294 vm.Value
	var v309 vm.Value
	var op_68 vm.Value
	var aux_70 vm.Value
	var v96 bool
	var v230 vm.Value
	var bu_231 vm.Value
	var ps_232 vm.Value
	var add_233 vm.Value
	var doseq_seq__26722_234 vm.Value
	var doseq_loop__26723_235 vm.Value
	var f_236 vm.Value
	var add_branch_args_237 vm.Value
	var bid_238 vm.Value
	var term_239 vm.Value
	var v284 vm.Value
	var v299 vm.Value
	var v314 vm.Value
	var v241 vm.Value
	var bu_71 vm.Value
	var ps_72 vm.Value
	var add_73 vm.Value
	var doseq_seq__26722_74 vm.Value
	var doseq_loop__26723_75 vm.Value
	var f_76 vm.Value
	var add_branch_args_77 vm.Value
	var bid_78 vm.Value
	var term_79 vm.Value
	var case__26724_80 vm.Value
	var op_81 vm.Value
	var aux_82 vm.Value
	var v274 vm.Value
	var v289 vm.Value
	var v304 vm.Value
	var arg__26827_99 vm.Value
	var arg__26831_101 vm.Value
	var arg__26837_103 vm.Value
	var arg__26841_105 vm.Value
	var v106 vm.Value
	var bu_83 vm.Value
	var ps_84 vm.Value
	var add_85 vm.Value
	var doseq_seq__26722_86 vm.Value
	var doseq_loop__26723_87 vm.Value
	var f_88 vm.Value
	var add_branch_args_89 vm.Value
	var bid_90 vm.Value
	var term_91 vm.Value
	var case__26724_92 vm.Value
	var op_93 vm.Value
	var aux_94 vm.Value
	var v282 vm.Value
	var v297 vm.Value
	var v312 vm.Value
	var v133 bool
	var v216 vm.Value
	var bu_217 vm.Value
	var ps_218 vm.Value
	var add_219 vm.Value
	var doseq_seq__26722_220 vm.Value
	var doseq_loop__26723_221 vm.Value
	var f_222 vm.Value
	var add_branch_args_223 vm.Value
	var bid_224 vm.Value
	var term_225 vm.Value
	var case__26724_226 vm.Value
	var op_227 vm.Value
	var aux_228 vm.Value
	var v278 vm.Value
	var v293 vm.Value
	var v308 vm.Value
	var bu_108 vm.Value
	var ps_109 vm.Value
	var add_110 vm.Value
	var doseq_seq__26722_111 vm.Value
	var doseq_loop__26723_112 vm.Value
	var f_113 vm.Value
	var add_branch_args_114 vm.Value
	var bid_115 vm.Value
	var term_116 vm.Value
	var case__26724_117 vm.Value
	var op_118 vm.Value
	var aux_119 vm.Value
	var v272 vm.Value
	var v287 vm.Value
	var v302 vm.Value
	var tt_136 vm.Value
	var ft_138 vm.Value
	var arg__26854_140 vm.Value
	var arg__26858_142 vm.Value
	var arg__26864_144 vm.Value
	var arg__26868_146 vm.Value
	var v147 vm.Value
	var arg__26873_149 vm.Value
	var arg__26877_151 vm.Value
	var arg__26883_153 vm.Value
	var arg__26887_155 vm.Value
	var v156 vm.Value
	var bu_120 vm.Value
	var ps_121 vm.Value
	var add_122 vm.Value
	var doseq_seq__26722_123 vm.Value
	var doseq_loop__26723_124 vm.Value
	var f_125 vm.Value
	var add_branch_args_126 vm.Value
	var bid_127 vm.Value
	var term_128 vm.Value
	var case__26724_129 vm.Value
	var op_130 vm.Value
	var aux_131 vm.Value
	var v280 vm.Value
	var v295 vm.Value
	var v310 vm.Value
	var v202 vm.Value
	var bu_203 vm.Value
	var ps_204 vm.Value
	var add_205 vm.Value
	var doseq_seq__26722_206 vm.Value
	var doseq_loop__26723_207 vm.Value
	var f_208 vm.Value
	var add_branch_args_209 vm.Value
	var bid_210 vm.Value
	var term_211 vm.Value
	var case__26724_212 vm.Value
	var op_213 vm.Value
	var aux_214 vm.Value
	var v273 vm.Value
	var v288 vm.Value
	var v303 vm.Value
	var bu_158 vm.Value
	var ps_159 vm.Value
	var add_160 vm.Value
	var doseq_seq__26722_161 vm.Value
	var doseq_loop__26723_162 vm.Value
	var f_163 vm.Value
	var add_branch_args_164 vm.Value
	var bid_165 vm.Value
	var term_166 vm.Value
	var case__26724_167 vm.Value
	var op_168 vm.Value
	var aux_169 vm.Value
	var v283 vm.Value
	var v298 vm.Value
	var v313 vm.Value
	var bu_170 vm.Value
	var ps_171 vm.Value
	var add_172 vm.Value
	var doseq_seq__26722_173 vm.Value
	var doseq_loop__26723_174 vm.Value
	var f_175 vm.Value
	var add_branch_args_176 vm.Value
	var bid_177 vm.Value
	var term_178 vm.Value
	var case__26724_179 vm.Value
	var op_180 vm.Value
	var aux_181 vm.Value
	var v277 vm.Value
	var v292 vm.Value
	var v307 vm.Value
	var v188 vm.Value
	var bu_189 vm.Value
	var ps_190 vm.Value
	var add_191 vm.Value
	var doseq_seq__26722_192 vm.Value
	var doseq_loop__26723_193 vm.Value
	var f_194 vm.Value
	var add_branch_args_195 vm.Value
	var bid_196 vm.Value
	var term_197 vm.Value
	var case__26724_198 vm.Value
	var op_199 vm.Value
	var aux_200 vm.Value
	var v271 vm.Value
	var v286 vm.Value
	var v301 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = bu_4, ps_8, arg__26794_16, arg__26799_19, doseq_seq__26722_20, doseq_loop__26723_21, f_22, add_branch_args_23, v276, v291, v306, bu_25, ps_26, add_27, doseq_seq__26722_28, doseq_loop__26723_29, f_30, add_branch_args_31, v275, v290, v305, bid_41, term_43, v63, bu_32, ps_33, add_34, doseq_seq__26722_35, doseq_loop__26723_36, f_37, add_branch_args_38, v285, v300, v315, v245, bu_246, ps_247, add_248, doseq_seq__26722_249, doseq_loop__26723_250, f_251, add_branch_args_252, arg__26896_256, arg__26901_259, arg__26906_262, v263, bu_44, ps_45, add_46, doseq_seq__26722_47, doseq_loop__26723_48, f_49, add_branch_args_50, bid_51, term_52, v281, v296, v311, bu_53, ps_54, add_55, doseq_seq__26722_56, doseq_loop__26723_57, f_58, add_branch_args_59, bid_60, term_61, v279, v294, v309, op_68, aux_70, v96, v230, bu_231, ps_232, add_233, doseq_seq__26722_234, doseq_loop__26723_235, f_236, add_branch_args_237, bid_238, term_239, v284, v299, v314, v241, bu_71, ps_72, add_73, doseq_seq__26722_74, doseq_loop__26723_75, f_76, add_branch_args_77, bid_78, term_79, case__26724_80, op_81, aux_82, v274, v289, v304, arg__26827_99, arg__26831_101, arg__26837_103, arg__26841_105, v106, bu_83, ps_84, add_85, doseq_seq__26722_86, doseq_loop__26723_87, f_88, add_branch_args_89, bid_90, term_91, case__26724_92, op_93, aux_94, v282, v297, v312, v133, v216, bu_217, ps_218, add_219, doseq_seq__26722_220, doseq_loop__26723_221, f_222, add_branch_args_223, bid_224, term_225, case__26724_226, op_227, aux_228, v278, v293, v308, bu_108, ps_109, add_110, doseq_seq__26722_111, doseq_loop__26723_112, f_113, add_branch_args_114, bid_115, term_116, case__26724_117, op_118, aux_119, v272, v287, v302, tt_136, ft_138, arg__26854_140, arg__26858_142, arg__26864_144, arg__26868_146, v147, arg__26873_149, arg__26877_151, arg__26883_153, arg__26887_155, v156, bu_120, ps_121, add_122, doseq_seq__26722_123, doseq_loop__26723_124, f_125, add_branch_args_126, bid_127, term_128, case__26724_129, op_130, aux_131, v280, v295, v310, v202, bu_203, ps_204, add_205, doseq_seq__26722_206, doseq_loop__26723_207, f_208, add_branch_args_209, bid_210, term_211, case__26724_212, op_213, aux_214, v273, v288, v303, bu_158, ps_159, add_160, doseq_seq__26722_161, doseq_loop__26723_162, f_163, add_branch_args_164, bid_165, term_166, case__26724_167, op_168, aux_169, v283, v298, v313, bu_170, ps_171, add_172, doseq_seq__26722_173, doseq_loop__26723_174, f_175, add_branch_args_176, bid_177, term_178, case__26724_179, op_180, aux_181, v277, v292, v307, v188, bu_189, ps_190, add_191, doseq_seq__26722_192, doseq_loop__26723_193, f_194, add_branch_args_195, bid_196, term_197, case__26724_198, op_199, aux_200, v271, v286, v301
	bu_4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "atom").Deref(), []vm.Value{vm.EmptyPersistentMap})
	if callErr != nil {
		return nil, callErr
	}
	ps_8, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "atom").Deref(), []vm.Value{vm.EmptyPersistentMap})
	if callErr != nil {
		return nil, callErr
	}
	arg__26794_16, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__26799_19, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	doseq_seq__26722_20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{arg__26799_19})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__26723_21 = doseq_seq__26722_20
	f_22 = arg0
	add_branch_args_23 = rt.BoxNativeFn(func(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
		var as_6 vm.Value
		var pos_7 vm.Value
		var add_8 vm.Value
		var target_9 vm.Value
		var ps_10 vm.Value
		var bu_11 vm.Value
		var pred_12 vm.Value
		var v32 vm.Value
		var args_15 vm.Value
		var as_16 vm.Value
		var pos_17 vm.Value
		var add_18 vm.Value
		var target_19 vm.Value
		var ps_20 vm.Value
		var bu_21 vm.Value
		var pred_22 vm.Value
		var arg_id_35 vm.Value
		var pp_37 vm.Value
		var v38 vm.Value
		var arg__26778_40 vm.Value
		var arg__26786_42 vm.Value
		var v43 vm.Value
		var v45 vm.Value
		var v46 vm.Value
		var args_23 vm.Value
		var as_24 vm.Value
		var pos_25 vm.Value
		var add_26 vm.Value
		var target_27 vm.Value
		var ps_28 vm.Value
		var bu_29 vm.Value
		var pred_30 vm.Value
		var v50 vm.Value
		var args_51 vm.Value
		var as_52 vm.Value
		var pos_53 vm.Value
		var add_54 vm.Value
		var target_55 vm.Value
		var ps_56 vm.Value
		var bu_57 vm.Value
		var pred_58 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = as_6, pos_7, add_8, target_9, ps_10, bu_11, pred_12, v32, args_15, as_16, pos_17, add_18, target_19, ps_20, bu_21, pred_22, arg_id_35, pp_37, v38, arg__26778_40, arg__26786_42, v43, v45, v46, args_23, as_24, pos_25, add_26, target_27, ps_28, bu_29, pred_30, v50, args_51, as_52, pos_53, add_54, target_55, ps_56, bu_57, pred_58
		as_6 = arg2
		pos_7 = vm.Int(0)
		add_8 = rt.BoxNativeFn(func(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
			var arg__26736_6 vm.Value
			var arg__26740_10 vm.Value
			var arg__26741_11 vm.Value
			var arg__26749_16 vm.Value
			var arg__26753_20 vm.Value
			var arg__26754_21 vm.Value
			var v22 vm.Value
			var callErr error
			_, _, _, _, _, _, _ = arg__26736_6, arg__26740_10, arg__26741_11, arg__26749_16, arg__26753_20, arg__26754_21, v22
			arg__26736_6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
			if callErr != nil {
				return nil, callErr
			}
			arg__26740_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
			if callErr != nil {
				return nil, callErr
			}
			arg__26741_11, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "fnil").Deref(), []vm.Value{rt.LookupVar("clojure.core", "conj").Deref(), arg__26740_10})
			if callErr != nil {
				return nil, callErr
			}
			arg__26749_16, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
			if callErr != nil {
				return nil, callErr
			}
			arg__26753_20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
			if callErr != nil {
				return nil, callErr
			}
			arg__26754_21, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "fnil").Deref(), []vm.Value{rt.LookupVar("clojure.core", "conj").Deref(), arg__26753_20})
			if callErr != nil {
				return nil, callErr
			}
			v22, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "update").Deref(), arg1, arg__26754_21, arg2})
			if callErr != nil {
				return nil, callErr
			}
			return v22, nil
		})
		target_9 = arg1
		ps_10 = ps_8
		bu_11 = bu_4
		pred_12 = arg0
		goto b1
	b1:
		;
		v32, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{as_6})
		if callErr != nil {
			return nil, callErr
		}
		if vm.IsTruthy(v32) {
			args_15 = arg2
			as_16 = as_6
			pos_17 = pos_7
			add_18 = add_8
			target_19 = target_9
			ps_20 = ps_10
			bu_21 = bu_11
			pred_22 = pred_12
			goto b2
		} else {
			args_23 = arg2
			as_24 = as_6
			pos_25 = pos_7
			add_26 = add_8
			target_27 = target_9
			ps_28 = ps_10
			bu_29 = bu_11
			pred_30 = pred_12
			goto b3
		}
	b2:
		;
		arg_id_35, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{as_16})
		if callErr != nil {
			return nil, callErr
		}
		pp_37, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{target_19, pos_17})
		if callErr != nil {
			return nil, callErr
		}
		v38, callErr = rt.InvokeValue(add_18, []vm.Value{bu_21, arg_id_35, pp_37})
		if callErr != nil {
			return nil, callErr
		}
		arg__26778_40, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{pred_22, target_19, arg_id_35})
		if callErr != nil {
			return nil, callErr
		}
		arg__26786_42, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{pred_22, target_19, arg_id_35})
		if callErr != nil {
			return nil, callErr
		}
		v43, callErr = rt.InvokeValue(add_18, []vm.Value{ps_20, pp_37, arg__26786_42})
		if callErr != nil {
			return nil, callErr
		}
		v45, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{as_16})
		if callErr != nil {
			return nil, callErr
		}
		v46 = rt.AddValue(pos_17, vm.Int(1))
		as_6 = v45
		pos_7 = v46
		add_8 = add_18
		target_9 = target_19
		ps_10 = ps_20
		bu_11 = bu_21
		pred_12 = pred_22
		goto b1
	b3:
		;
		v50 = vm.NIL
		args_51 = args_23
		as_52 = as_24
		pos_53 = pos_25
		add_54 = add_26
		target_55 = target_27
		ps_56 = ps_28
		bu_57 = bu_29
		pred_58 = pred_30
		goto b4
	b4:
		;
		return v50, nil
	})
	v276 = vm.Keyword("branch")
	v291 = vm.Keyword("branch-if")
	v306 = vm.Keyword("else")
	goto b1
b1:
	;
	if vm.IsTruthy(doseq_loop__26723_21) {
		bu_25 = bu_4
		ps_26 = ps_8
		add_27 = rt.BoxNativeFn(func(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
			var arg__26736_6 vm.Value
			var arg__26740_10 vm.Value
			var arg__26741_11 vm.Value
			var arg__26749_16 vm.Value
			var arg__26753_20 vm.Value
			var arg__26754_21 vm.Value
			var v22 vm.Value
			var callErr error
			_, _, _, _, _, _, _ = arg__26736_6, arg__26740_10, arg__26741_11, arg__26749_16, arg__26753_20, arg__26754_21, v22
			arg__26736_6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
			if callErr != nil {
				return nil, callErr
			}
			arg__26740_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
			if callErr != nil {
				return nil, callErr
			}
			arg__26741_11, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "fnil").Deref(), []vm.Value{rt.LookupVar("clojure.core", "conj").Deref(), arg__26740_10})
			if callErr != nil {
				return nil, callErr
			}
			arg__26749_16, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
			if callErr != nil {
				return nil, callErr
			}
			arg__26753_20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
			if callErr != nil {
				return nil, callErr
			}
			arg__26754_21, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "fnil").Deref(), []vm.Value{rt.LookupVar("clojure.core", "conj").Deref(), arg__26753_20})
			if callErr != nil {
				return nil, callErr
			}
			v22, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "update").Deref(), arg1, arg__26754_21, arg2})
			if callErr != nil {
				return nil, callErr
			}
			return v22, nil
		})
		doseq_seq__26722_28 = doseq_seq__26722_20
		doseq_loop__26723_29 = doseq_loop__26723_21
		f_30 = f_22
		add_branch_args_31 = add_branch_args_23
		v275 = v276
		v290 = v291
		v305 = v306
		goto b2
	} else {
		bu_32 = bu_4
		ps_33 = ps_8
		add_34 = rt.BoxNativeFn(func(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
			var arg__26736_6 vm.Value
			var arg__26740_10 vm.Value
			var arg__26741_11 vm.Value
			var arg__26749_16 vm.Value
			var arg__26753_20 vm.Value
			var arg__26754_21 vm.Value
			var v22 vm.Value
			var callErr error
			_, _, _, _, _, _, _ = arg__26736_6, arg__26740_10, arg__26741_11, arg__26749_16, arg__26753_20, arg__26754_21, v22
			arg__26736_6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
			if callErr != nil {
				return nil, callErr
			}
			arg__26740_10, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
			if callErr != nil {
				return nil, callErr
			}
			arg__26741_11, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "fnil").Deref(), []vm.Value{rt.LookupVar("clojure.core", "conj").Deref(), arg__26740_10})
			if callErr != nil {
				return nil, callErr
			}
			arg__26749_16, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
			if callErr != nil {
				return nil, callErr
			}
			arg__26753_20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "hash-set").Deref(), []vm.Value{})
			if callErr != nil {
				return nil, callErr
			}
			arg__26754_21, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "fnil").Deref(), []vm.Value{rt.LookupVar("clojure.core", "conj").Deref(), arg__26753_20})
			if callErr != nil {
				return nil, callErr
			}
			v22, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{arg0, rt.LookupVar("clojure.core", "update").Deref(), arg1, arg__26754_21, arg2})
			if callErr != nil {
				return nil, callErr
			}
			return v22, nil
		})
		doseq_seq__26722_35 = doseq_seq__26722_20
		doseq_loop__26723_36 = doseq_loop__26723_21
		f_37 = f_22
		add_branch_args_38 = add_branch_args_23
		v285 = v276
		v300 = v291
		v315 = v306
		goto b3
	}
b2:
	;
	bid_41, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{doseq_loop__26723_29})
	if callErr != nil {
		return nil, callErr
	}
	term_43, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-term").Deref(), []vm.Value{bid_41, f_30})
	if callErr != nil {
		return nil, callErr
	}
	v63, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{term_43})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v63) {
		bu_44 = bu_25
		ps_45 = ps_26
		add_46 = add_27
		doseq_seq__26722_47 = doseq_seq__26722_28
		doseq_loop__26723_48 = doseq_loop__26723_29
		f_49 = f_30
		add_branch_args_50 = add_branch_args_31
		bid_51 = bid_41
		term_52 = term_43
		v281 = v275
		v296 = v290
		v311 = v305
		goto b5
	} else {
		bu_53 = bu_25
		ps_54 = ps_26
		add_55 = add_27
		doseq_seq__26722_56 = doseq_seq__26722_28
		doseq_loop__26723_57 = doseq_loop__26723_29
		f_58 = f_30
		add_branch_args_59 = add_branch_args_31
		bid_60 = bid_41
		term_61 = term_43
		v279 = v275
		v294 = v290
		v309 = v305
		goto b6
	}
b3:
	;
	v245 = vm.NIL
	bu_246 = bu_32
	ps_247 = ps_33
	add_248 = add_34
	doseq_seq__26722_249 = doseq_seq__26722_35
	doseq_loop__26723_250 = doseq_loop__26723_36
	f_251 = f_37
	add_branch_args_252 = add_branch_args_38
	goto b4
b4:
	;
	arg__26896_256, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{ps_247})
	if callErr != nil {
		return nil, callErr
	}
	arg__26901_259, callErr = rt.InvokeValue(rt.LookupVar("ir", "uses").Deref(), []vm.Value{f_251})
	if callErr != nil {
		return nil, callErr
	}
	arg__26906_262, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{bu_246})
	if callErr != nil {
		return nil, callErr
	}
	v263, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("param-sources"), arg__26896_256, vm.Keyword("uses"), arg__26901_259, vm.Keyword("branch-arg-users"), arg__26906_262})
	if callErr != nil {
		return nil, callErr
	}
	return v263, nil
b5:
	;
	v230 = vm.NIL
	bu_231 = bu_44
	ps_232 = ps_45
	add_233 = add_46
	doseq_seq__26722_234 = doseq_seq__26722_47
	doseq_loop__26723_235 = doseq_loop__26723_48
	f_236 = f_49
	add_branch_args_237 = add_branch_args_50
	bid_238 = bid_51
	term_239 = term_52
	v284 = v281
	v299 = v296
	v314 = v311
	goto b7
b6:
	;
	op_68, callErr = rt.InvokeValue(rt.LookupVar("ir", "op").Deref(), []vm.Value{term_61, f_58})
	if callErr != nil {
		return nil, callErr
	}
	aux_70, callErr = rt.InvokeValue(rt.LookupVar("ir", "aux").Deref(), []vm.Value{term_61, f_58})
	if callErr != nil {
		return nil, callErr
	}
	v96 = op_68 == v279
	if v96 {
		bu_71 = bu_53
		ps_72 = ps_54
		add_73 = add_55
		doseq_seq__26722_74 = doseq_seq__26722_56
		doseq_loop__26723_75 = doseq_loop__26723_57
		f_76 = f_58
		add_branch_args_77 = add_branch_args_59
		bid_78 = bid_60
		term_79 = term_61
		case__26724_80 = op_68
		op_81 = op_68
		aux_82 = aux_70
		v274 = v279
		v289 = v294
		v304 = v309
		goto b8
	} else {
		bu_83 = bu_53
		ps_84 = ps_54
		add_85 = add_55
		doseq_seq__26722_86 = doseq_seq__26722_56
		doseq_loop__26723_87 = doseq_loop__26723_57
		f_88 = f_58
		add_branch_args_89 = add_branch_args_59
		bid_90 = bid_60
		term_91 = term_61
		case__26724_92 = op_68
		op_93 = op_68
		aux_94 = aux_70
		v282 = v279
		v297 = v294
		v312 = v309
		goto b9
	}
b7:
	;
	v241, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "next").Deref(), []vm.Value{doseq_loop__26723_235})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__26723_21 = v241
	f_22 = f_236
	add_branch_args_23 = add_branch_args_237
	v276 = v284
	v291 = v299
	v306 = v314
	goto b1
b8:
	;
	arg__26827_99, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-target").Deref(), []vm.Value{aux_82})
	if callErr != nil {
		return nil, callErr
	}
	arg__26831_101, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-args").Deref(), []vm.Value{aux_82})
	if callErr != nil {
		return nil, callErr
	}
	arg__26837_103, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-target").Deref(), []vm.Value{aux_82})
	if callErr != nil {
		return nil, callErr
	}
	arg__26841_105, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-args").Deref(), []vm.Value{aux_82})
	if callErr != nil {
		return nil, callErr
	}
	v106, callErr = rt.InvokeValue(add_branch_args_77, []vm.Value{bid_78, arg__26837_103, arg__26841_105})
	if callErr != nil {
		return nil, callErr
	}
	v216 = v106
	bu_217 = bu_71
	ps_218 = ps_72
	add_219 = add_73
	doseq_seq__26722_220 = doseq_seq__26722_74
	doseq_loop__26723_221 = doseq_loop__26723_75
	f_222 = f_76
	add_branch_args_223 = add_branch_args_77
	bid_224 = bid_78
	term_225 = term_79
	case__26724_226 = case__26724_80
	op_227 = op_81
	aux_228 = aux_82
	v278 = v274
	v293 = v289
	v308 = v304
	goto b10
b9:
	;
	v133 = case__26724_92 == v297
	if v133 {
		bu_108 = bu_83
		ps_109 = ps_84
		add_110 = add_85
		doseq_seq__26722_111 = doseq_seq__26722_86
		doseq_loop__26723_112 = doseq_loop__26723_87
		f_113 = f_88
		add_branch_args_114 = add_branch_args_89
		bid_115 = bid_90
		term_116 = term_91
		case__26724_117 = case__26724_92
		op_118 = op_93
		aux_119 = aux_94
		v272 = v282
		v287 = v297
		v302 = v312
		goto b11
	} else {
		bu_120 = bu_83
		ps_121 = ps_84
		add_122 = add_85
		doseq_seq__26722_123 = doseq_seq__26722_86
		doseq_loop__26723_124 = doseq_loop__26723_87
		f_125 = f_88
		add_branch_args_126 = add_branch_args_89
		bid_127 = bid_90
		term_128 = term_91
		case__26724_129 = case__26724_92
		op_130 = op_93
		aux_131 = aux_94
		v280 = v282
		v295 = v297
		v310 = v312
		goto b12
	}
b10:
	;
	v230 = v216
	bu_231 = bu_217
	ps_232 = ps_218
	add_233 = add_219
	doseq_seq__26722_234 = doseq_seq__26722_220
	doseq_loop__26723_235 = doseq_loop__26723_221
	f_236 = f_222
	add_branch_args_237 = add_branch_args_223
	bid_238 = bid_224
	term_239 = term_225
	v284 = v278
	v299 = v293
	v314 = v308
	goto b7
b11:
	;
	tt_136, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-true").Deref(), []vm.Value{aux_119})
	if callErr != nil {
		return nil, callErr
	}
	ft_138, callErr = rt.InvokeValue(rt.LookupVar("ir", "cond-target-false").Deref(), []vm.Value{aux_119})
	if callErr != nil {
		return nil, callErr
	}
	arg__26854_140, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-target").Deref(), []vm.Value{tt_136})
	if callErr != nil {
		return nil, callErr
	}
	arg__26858_142, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-args").Deref(), []vm.Value{tt_136})
	if callErr != nil {
		return nil, callErr
	}
	arg__26864_144, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-target").Deref(), []vm.Value{tt_136})
	if callErr != nil {
		return nil, callErr
	}
	arg__26868_146, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-args").Deref(), []vm.Value{tt_136})
	if callErr != nil {
		return nil, callErr
	}
	v147, callErr = rt.InvokeValue(add_branch_args_114, []vm.Value{bid_115, arg__26864_144, arg__26868_146})
	if callErr != nil {
		return nil, callErr
	}
	arg__26873_149, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-target").Deref(), []vm.Value{ft_138})
	if callErr != nil {
		return nil, callErr
	}
	arg__26877_151, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-args").Deref(), []vm.Value{ft_138})
	if callErr != nil {
		return nil, callErr
	}
	arg__26883_153, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-target").Deref(), []vm.Value{ft_138})
	if callErr != nil {
		return nil, callErr
	}
	arg__26887_155, callErr = rt.InvokeValue(rt.LookupVar("ir", "branch-target-args").Deref(), []vm.Value{ft_138})
	if callErr != nil {
		return nil, callErr
	}
	v156, callErr = rt.InvokeValue(add_branch_args_114, []vm.Value{bid_115, arg__26883_153, arg__26887_155})
	if callErr != nil {
		return nil, callErr
	}
	v202 = v156
	bu_203 = bu_108
	ps_204 = ps_109
	add_205 = add_110
	doseq_seq__26722_206 = doseq_seq__26722_111
	doseq_loop__26723_207 = doseq_loop__26723_112
	f_208 = f_113
	add_branch_args_209 = add_branch_args_114
	bid_210 = bid_115
	term_211 = term_116
	case__26724_212 = case__26724_117
	op_213 = op_118
	aux_214 = aux_119
	v273 = v272
	v288 = v287
	v303 = v302
	goto b13
b12:
	;
	if vm.IsTruthy(v310) {
		bu_158 = bu_120
		ps_159 = ps_121
		add_160 = add_122
		doseq_seq__26722_161 = doseq_seq__26722_123
		doseq_loop__26723_162 = doseq_loop__26723_124
		f_163 = f_125
		add_branch_args_164 = add_branch_args_126
		bid_165 = bid_127
		term_166 = term_128
		case__26724_167 = case__26724_129
		op_168 = op_130
		aux_169 = aux_131
		v283 = v280
		v298 = v295
		v313 = v310
		goto b14
	} else {
		bu_170 = bu_120
		ps_171 = ps_121
		add_172 = add_122
		doseq_seq__26722_173 = doseq_seq__26722_123
		doseq_loop__26723_174 = doseq_loop__26723_124
		f_175 = f_125
		add_branch_args_176 = add_branch_args_126
		bid_177 = bid_127
		term_178 = term_128
		case__26724_179 = case__26724_129
		op_180 = op_130
		aux_181 = aux_131
		v277 = v280
		v292 = v295
		v307 = v310
		goto b15
	}
b13:
	;
	v216 = v202
	bu_217 = bu_203
	ps_218 = ps_204
	add_219 = add_205
	doseq_seq__26722_220 = doseq_seq__26722_206
	doseq_loop__26723_221 = doseq_loop__26723_207
	f_222 = f_208
	add_branch_args_223 = add_branch_args_209
	bid_224 = bid_210
	term_225 = term_211
	case__26724_226 = case__26724_212
	op_227 = op_213
	aux_228 = aux_214
	v278 = v273
	v293 = v288
	v308 = v303
	goto b10
b14:
	;
	v188 = vm.NIL
	bu_189 = bu_158
	ps_190 = ps_159
	add_191 = add_160
	doseq_seq__26722_192 = doseq_seq__26722_161
	doseq_loop__26723_193 = doseq_loop__26723_162
	f_194 = f_163
	add_branch_args_195 = add_branch_args_164
	bid_196 = bid_165
	term_197 = term_166
	case__26724_198 = case__26724_167
	op_199 = op_168
	aux_200 = aux_169
	v271 = v283
	v286 = v298
	v301 = v313
	goto b16
b15:
	;
	v188 = vm.NIL
	bu_189 = bu_170
	ps_190 = ps_171
	add_191 = add_172
	doseq_seq__26722_192 = doseq_seq__26722_173
	doseq_loop__26723_193 = doseq_loop__26723_174
	f_194 = f_175
	add_branch_args_195 = add_branch_args_176
	bid_196 = bid_177
	term_197 = term_178
	case__26724_198 = case__26724_179
	op_199 = op_180
	aux_200 = aux_181
	v271 = v277
	v286 = v292
	v301 = v307
	goto b16
b16:
	;
	v202 = v188
	bu_203 = bu_189
	ps_204 = ps_190
	add_205 = add_191
	doseq_seq__26722_206 = doseq_seq__26722_192
	doseq_loop__26723_207 = doseq_loop__26723_193
	f_208 = f_194
	add_branch_args_209 = add_branch_args_195
	bid_210 = bid_196
	term_211 = term_197
	case__26724_212 = case__26724_198
	op_213 = op_199
	aux_214 = aux_200
	v273 = v271
	v288 = v286
	v303 = v301
	goto b13
}
func classify_const(arg0 vm.Value) (vm.Value, error) {
	var v5 vm.Value
	var aux_2 vm.Value
	var aux_3 vm.Value
	var v12 bool
	var v122 vm.Value
	var aux_123 vm.Value
	var aux_9 vm.Value
	var aux_10 vm.Value
	var v19 bool
	var v119 vm.Value
	var aux_120 vm.Value
	var aux_16 vm.Value
	var aux_17 vm.Value
	var and__x_26 vm.Value
	var v116 vm.Value
	var aux_117 vm.Value
	var aux_23 vm.Value
	var v44 vm.Value
	var aux_24 vm.Value
	var v49 vm.Value
	var v113 vm.Value
	var aux_114 vm.Value
	var aux_27 vm.Value
	var and__x_28 vm.Value
	var v33 bool
	var aux_29 vm.Value
	var and__x_30 vm.Value
	var v36 vm.Value
	var aux_37 vm.Value
	var and__x_38 vm.Value
	var aux_46 vm.Value
	var aux_47 vm.Value
	var and__x_56 vm.Value
	var v110 vm.Value
	var aux_111 vm.Value
	var aux_53 vm.Value
	var v74 vm.Value
	var aux_54 vm.Value
	var v79 vm.Value
	var v107 vm.Value
	var aux_108 vm.Value
	var aux_57 vm.Value
	var and__x_58 vm.Value
	var v63 bool
	var aux_59 vm.Value
	var and__x_60 vm.Value
	var v66 vm.Value
	var aux_67 vm.Value
	var and__x_68 vm.Value
	var aux_76 vm.Value
	var aux_77 vm.Value
	var v86 vm.Value
	var v104 vm.Value
	var aux_105 vm.Value
	var aux_83 vm.Value
	var aux_84 vm.Value
	var v101 vm.Value
	var aux_102 vm.Value
	var aux_90 vm.Value
	var aux_91 vm.Value
	var v98 vm.Value
	var aux_99 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = v5, aux_2, aux_3, v12, v122, aux_123, aux_9, aux_10, v19, v119, aux_120, aux_16, aux_17, and__x_26, v116, aux_117, aux_23, v44, aux_24, v49, v113, aux_114, aux_27, and__x_28, v33, aux_29, and__x_30, v36, aux_37, and__x_38, aux_46, aux_47, and__x_56, v110, aux_111, aux_53, v74, aux_54, v79, v107, aux_108, aux_57, and__x_58, v63, aux_59, and__x_60, v66, aux_67, and__x_68, aux_76, aux_77, v86, v104, aux_105, aux_83, aux_84, v101, aux_102, aux_90, aux_91, v98, aux_99
	v5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v5) {
		aux_2 = arg0
		goto b1
	} else {
		aux_3 = arg0
		goto b2
	}
b1:
	;
	v122 = vm.Keyword("nil")
	aux_123 = aux_2
	goto b3
b2:
	;
	v12 = aux_3 == vm.Boolean(true)
	if v12 {
		aux_9 = aux_3
		goto b4
	} else {
		aux_10 = aux_3
		goto b5
	}
b3:
	;
	return v122, nil
b4:
	;
	v119 = vm.Keyword("true")
	aux_120 = aux_9
	goto b6
b5:
	;
	v19 = aux_10 == vm.Boolean(false)
	if v19 {
		aux_16 = aux_10
		goto b7
	} else {
		aux_17 = aux_10
		goto b8
	}
b6:
	;
	v122 = v119
	aux_123 = aux_120
	goto b3
b7:
	;
	v116 = vm.Keyword("false")
	aux_117 = aux_16
	goto b9
b8:
	;
	and__x_26, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "integer?").Deref(), []vm.Value{aux_17})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(and__x_26) {
		aux_27 = aux_17
		and__x_28 = and__x_26
		goto b13
	} else {
		aux_29 = aux_17
		and__x_30 = and__x_26
		goto b14
	}
b9:
	;
	v119 = v116
	aux_120 = aux_117
	goto b6
b10:
	;
	v44, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("const"), vm.Keyword("int"), vm.Int(0)})
	if callErr != nil {
		return nil, callErr
	}
	v113 = v44
	aux_114 = aux_23
	goto b12
b11:
	;
	v49, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "integer?").Deref(), []vm.Value{aux_24})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v49) {
		aux_46 = aux_24
		goto b16
	} else {
		aux_47 = aux_24
		goto b17
	}
b12:
	;
	v116 = v113
	aux_117 = aux_114
	goto b9
b13:
	;
	v33 = aux_27 == vm.Int(0)
	v36 = vm.Boolean(v33)
	aux_37 = aux_27
	and__x_38 = and__x_28
	goto b15
b14:
	;
	v36 = and__x_30
	aux_37 = aux_29
	and__x_38 = and__x_30
	goto b15
b15:
	;
	if vm.IsTruthy(v36) {
		aux_23 = aux_37
		goto b10
	} else {
		aux_24 = aux_37
		goto b11
	}
b16:
	;
	v110 = vm.Keyword("int")
	aux_111 = aux_46
	goto b18
b17:
	;
	and__x_56, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "float?").Deref(), []vm.Value{aux_47})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(and__x_56) {
		aux_57 = aux_47
		and__x_58 = and__x_56
		goto b22
	} else {
		aux_59 = aux_47
		and__x_60 = and__x_56
		goto b23
	}
b18:
	;
	v113 = v110
	aux_114 = aux_111
	goto b12
b19:
	;
	v74, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("const"), vm.Keyword("float"), vm.Float(0)})
	if callErr != nil {
		return nil, callErr
	}
	v107 = v74
	aux_108 = aux_53
	goto b21
b20:
	;
	v79, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "float?").Deref(), []vm.Value{aux_54})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v79) {
		aux_76 = aux_54
		goto b25
	} else {
		aux_77 = aux_54
		goto b26
	}
b21:
	;
	v110 = v107
	aux_111 = aux_108
	goto b18
b22:
	;
	v63 = aux_57 == vm.Float(0)
	v66 = vm.Boolean(v63)
	aux_67 = aux_57
	and__x_68 = and__x_58
	goto b24
b23:
	;
	v66 = and__x_60
	aux_67 = aux_59
	and__x_68 = and__x_60
	goto b24
b24:
	;
	if vm.IsTruthy(v66) {
		aux_53 = aux_67
		goto b19
	} else {
		aux_54 = aux_67
		goto b20
	}
b25:
	;
	v104 = vm.Keyword("float")
	aux_105 = aux_76
	goto b27
b26:
	;
	v86, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "string?").Deref(), []vm.Value{aux_77})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v86) {
		aux_83 = aux_77
		goto b28
	} else {
		aux_84 = aux_77
		goto b29
	}
b27:
	;
	v107 = v104
	aux_108 = aux_105
	goto b21
b28:
	;
	v101 = vm.Keyword("string")
	aux_102 = aux_83
	goto b30
b29:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		aux_90 = aux_84
		goto b31
	} else {
		aux_91 = aux_84
		goto b32
	}
b30:
	;
	v104 = v101
	aux_105 = aux_102
	goto b27
b31:
	;
	v98 = vm.Keyword("any")
	aux_99 = aux_90
	goto b33
b32:
	;
	v98 = vm.NIL
	aux_99 = aux_91
	goto b33
b33:
	;
	v101 = v98
	aux_102 = aux_99
	goto b30
}
func arg_seed(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var stored_3 vm.Value
	var arg_types_5 vm.Value
	var idx_7 vm.Value
	var and__x_21 vm.Value
	var inst_8 vm.Value
	var f_9 vm.Value
	var stored_10 vm.Value
	var arg_types_11 vm.Value
	var idx_12 vm.Value
	var v50 vm.Value
	var inst_13 vm.Value
	var f_14 vm.Value
	var stored_15 vm.Value
	var arg_types_16 vm.Value
	var idx_17 vm.Value
	var and__x_63 vm.Value
	var v137 vm.Value
	var inst_138 vm.Value
	var f_139 vm.Value
	var stored_140 vm.Value
	var arg_types_141 vm.Value
	var idx_142 vm.Value
	var inst_22 vm.Value
	var f_23 vm.Value
	var stored_24 vm.Value
	var arg_types_25 vm.Value
	var idx_26 vm.Value
	var and__x_27 vm.Value
	var v38 vm.Value
	var inst_28 vm.Value
	var f_29 vm.Value
	var stored_30 vm.Value
	var arg_types_31 vm.Value
	var idx_32 vm.Value
	var and__x_33 vm.Value
	var v41 vm.Value
	var inst_42 vm.Value
	var f_43 vm.Value
	var stored_44 vm.Value
	var arg_types_45 vm.Value
	var idx_46 vm.Value
	var and__x_47 vm.Value
	var inst_52 vm.Value
	var f_53 vm.Value
	var stored_54 vm.Value
	var arg_types_55 vm.Value
	var idx_56 vm.Value
	var arg__27110_93 vm.Value
	var arg__27116_96 vm.Value
	var arg__27117_97 vm.Value
	var arg__27123_100 vm.Value
	var arg__27129_103 vm.Value
	var arg__27130_104 vm.Value
	var v105 vm.Value
	var inst_57 vm.Value
	var f_58 vm.Value
	var stored_59 vm.Value
	var arg_types_60 vm.Value
	var idx_61 vm.Value
	var v130 vm.Value
	var inst_131 vm.Value
	var f_132 vm.Value
	var stored_133 vm.Value
	var arg_types_134 vm.Value
	var idx_135 vm.Value
	var inst_64 vm.Value
	var f_65 vm.Value
	var stored_66 vm.Value
	var arg_types_67 vm.Value
	var idx_68 vm.Value
	var and__x_69 vm.Value
	var arg__27101_78 vm.Value
	var arg__27105_80 vm.Value
	var v81 bool
	var inst_70 vm.Value
	var f_71 vm.Value
	var stored_72 vm.Value
	var arg_types_73 vm.Value
	var idx_74 vm.Value
	var and__x_75 vm.Value
	var v84 vm.Value
	var inst_85 vm.Value
	var f_86 vm.Value
	var stored_87 vm.Value
	var arg_types_88 vm.Value
	var idx_89 vm.Value
	var and__x_90 vm.Value
	var inst_107 vm.Value
	var f_108 vm.Value
	var stored_109 vm.Value
	var arg_types_110 vm.Value
	var idx_111 vm.Value
	var inst_112 vm.Value
	var f_113 vm.Value
	var stored_114 vm.Value
	var arg_types_115 vm.Value
	var idx_116 vm.Value
	var v123 vm.Value
	var inst_124 vm.Value
	var f_125 vm.Value
	var stored_126 vm.Value
	var arg_types_127 vm.Value
	var idx_128 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = stored_3, arg_types_5, idx_7, and__x_21, inst_8, f_9, stored_10, arg_types_11, idx_12, v50, inst_13, f_14, stored_15, arg_types_16, idx_17, and__x_63, v137, inst_138, f_139, stored_140, arg_types_141, idx_142, inst_22, f_23, stored_24, arg_types_25, idx_26, and__x_27, v38, inst_28, f_29, stored_30, arg_types_31, idx_32, and__x_33, v41, inst_42, f_43, stored_44, arg_types_45, idx_46, and__x_47, inst_52, f_53, stored_54, arg_types_55, idx_56, arg__27110_93, arg__27116_96, arg__27117_97, arg__27123_100, arg__27129_103, arg__27130_104, v105, inst_57, f_58, stored_59, arg_types_60, idx_61, v130, inst_131, f_132, stored_133, arg_types_134, idx_135, inst_64, f_65, stored_66, arg_types_67, idx_68, and__x_69, arg__27101_78, arg__27105_80, v81, inst_70, f_71, stored_72, arg_types_73, idx_74, and__x_75, v84, inst_85, f_86, stored_87, arg_types_88, idx_89, and__x_90, inst_107, f_108, stored_109, arg_types_110, idx_111, inst_112, f_113, stored_114, arg_types_115, idx_116, v123, inst_124, f_125, stored_126, arg_types_127, idx_128
	stored_3, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg_types_5, callErr = rt.InvokeValue(rt.LookupVar("ir", "fn-arg-types").Deref(), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	idx_7, callErr = rt.InvokeValue(rt.LookupVar("ir", "aux").Deref(), []vm.Value{arg0, arg1})
	if callErr != nil {
		return nil, callErr
	}
	and__x_21, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not=").Deref(), []vm.Value{stored_3, vm.Keyword("unknown")})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(and__x_21) {
		inst_22 = arg0
		f_23 = arg1
		stored_24 = stored_3
		arg_types_25 = arg_types_5
		idx_26 = idx_7
		and__x_27 = and__x_21
		goto b4
	} else {
		inst_28 = arg0
		f_29 = arg1
		stored_30 = stored_3
		arg_types_31 = arg_types_5
		idx_32 = idx_7
		and__x_33 = and__x_21
		goto b5
	}
b1:
	;
	v50, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "normalize-type").Deref(), []vm.Value{stored_10})
	if callErr != nil {
		return nil, callErr
	}
	v137 = v50
	inst_138 = inst_8
	f_139 = f_9
	stored_140 = stored_10
	arg_types_141 = arg_types_11
	idx_142 = idx_12
	goto b3
b2:
	;
	and__x_63, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector?").Deref(), []vm.Value{arg_types_16})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(and__x_63) {
		inst_64 = inst_13
		f_65 = f_14
		stored_66 = stored_15
		arg_types_67 = arg_types_16
		idx_68 = idx_17
		and__x_69 = and__x_63
		goto b10
	} else {
		inst_70 = inst_13
		f_71 = f_14
		stored_72 = stored_15
		arg_types_73 = arg_types_16
		idx_74 = idx_17
		and__x_75 = and__x_63
		goto b11
	}
b3:
	;
	return v137, nil
b4:
	;
	v38, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not=").Deref(), []vm.Value{stored_24, vm.Keyword("bottom")})
	if callErr != nil {
		return nil, callErr
	}
	v41 = v38
	inst_42 = inst_22
	f_43 = f_23
	stored_44 = stored_24
	arg_types_45 = arg_types_25
	idx_46 = idx_26
	and__x_47 = and__x_27
	goto b6
b5:
	;
	v41 = and__x_33
	inst_42 = inst_28
	f_43 = f_29
	stored_44 = stored_30
	arg_types_45 = arg_types_31
	idx_46 = idx_32
	and__x_47 = and__x_33
	goto b6
b6:
	;
	if vm.IsTruthy(v41) {
		inst_8 = inst_42
		f_9 = f_43
		stored_10 = stored_44
		arg_types_11 = arg_types_45
		idx_12 = idx_46
		goto b1
	} else {
		inst_13 = inst_42
		f_14 = f_43
		stored_15 = stored_44
		arg_types_16 = arg_types_45
		idx_17 = idx_46
		goto b2
	}
b7:
	;
	arg__27110_93, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "int").Deref(), []vm.Value{idx_56})
	if callErr != nil {
		return nil, callErr
	}
	arg__27116_96, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "int").Deref(), []vm.Value{idx_56})
	if callErr != nil {
		return nil, callErr
	}
	arg__27117_97, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg_types_55, arg__27116_96})
	if callErr != nil {
		return nil, callErr
	}
	arg__27123_100, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "int").Deref(), []vm.Value{idx_56})
	if callErr != nil {
		return nil, callErr
	}
	arg__27129_103, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "int").Deref(), []vm.Value{idx_56})
	if callErr != nil {
		return nil, callErr
	}
	arg__27130_104, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg_types_55, arg__27129_103})
	if callErr != nil {
		return nil, callErr
	}
	v105, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "normalize-type").Deref(), []vm.Value{arg__27130_104})
	if callErr != nil {
		return nil, callErr
	}
	v130 = v105
	inst_131 = inst_52
	f_132 = f_53
	stored_133 = stored_54
	arg_types_134 = arg_types_55
	idx_135 = idx_56
	goto b9
b8:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		inst_107 = inst_57
		f_108 = f_58
		stored_109 = stored_59
		arg_types_110 = arg_types_60
		idx_111 = idx_61
		goto b13
	} else {
		inst_112 = inst_57
		f_113 = f_58
		stored_114 = stored_59
		arg_types_115 = arg_types_60
		idx_116 = idx_61
		goto b14
	}
b9:
	;
	v137 = v130
	inst_138 = inst_131
	f_139 = f_132
	stored_140 = stored_133
	arg_types_141 = arg_types_134
	idx_142 = idx_135
	goto b3
b10:
	;
	arg__27101_78, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "int").Deref(), []vm.Value{idx_68})
	if callErr != nil {
		return nil, callErr
	}
	arg__27105_80, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{arg_types_67})
	if callErr != nil {
		return nil, callErr
	}
	v81 = rt.LtValue(arg__27101_78, arg__27105_80)
	v84 = vm.Boolean(v81)
	inst_85 = inst_64
	f_86 = f_65
	stored_87 = stored_66
	arg_types_88 = arg_types_67
	idx_89 = idx_68
	and__x_90 = and__x_69
	goto b12
b11:
	;
	v84 = and__x_75
	inst_85 = inst_70
	f_86 = f_71
	stored_87 = stored_72
	arg_types_88 = arg_types_73
	idx_89 = idx_74
	and__x_90 = and__x_75
	goto b12
b12:
	;
	if vm.IsTruthy(v84) {
		inst_52 = inst_85
		f_53 = f_86
		stored_54 = stored_87
		arg_types_55 = arg_types_88
		idx_56 = idx_89
		goto b7
	} else {
		inst_57 = inst_85
		f_58 = f_86
		stored_59 = stored_87
		arg_types_60 = arg_types_88
		idx_61 = idx_89
		goto b8
	}
b13:
	;
	v123 = vm.Keyword("unknown")
	inst_124 = inst_107
	f_125 = f_108
	stored_126 = stored_109
	arg_types_127 = arg_types_110
	idx_128 = idx_111
	goto b15
b14:
	;
	v123 = vm.NIL
	inst_124 = inst_112
	f_125 = f_113
	stored_126 = stored_114
	arg_types_127 = arg_types_115
	idx_128 = idx_116
	goto b15
b15:
	;
	v130 = v123
	inst_131 = inst_124
	f_132 = f_125
	stored_133 = stored_126
	arg_types_134 = arg_types_127
	idx_135 = idx_128
	goto b9
}
func propagate_changed_BANG_(arg0 int, arg1 vm.Value, arg2 vm.Value, arg3 vm.Value) (vm.Value, error) {
	var uses_5 vm.Value
	var arg__27140_17 vm.Value
	var v18 bool
	var v_6 int
	var deps_7 vm.Value
	var queue_8 vm.Value
	var queued_9 vm.Value
	var uses_10 vm.Value
	var v21 vm.Value
	var v_11 int
	var deps_12 vm.Value
	var queue_13 vm.Value
	var queued_14 vm.Value
	var uses_15 vm.Value
	var us_25 vm.Value
	var v_26 int
	var deps_27 vm.Value
	var queue_28 vm.Value
	var queued_29 vm.Value
	var uses_30 vm.Value
	var arg__27149_32 vm.Value
	var arg__27154_35 vm.Value
	var state_36 vm.Value
	var us_37 vm.Value
	var v_38 int
	var deps_39 vm.Value
	var queue_40 vm.Value
	var queued_41 vm.Value
	var uses_42 vm.Value
	var state_43 vm.Value
	var v93 vm.Value
	var us_44 vm.Value
	var v_45 int
	var deps_46 vm.Value
	var queue_47 vm.Value
	var queued_48 vm.Value
	var uses_49 vm.Value
	var state_50 vm.Value
	var v97 vm.Value
	var us_98 vm.Value
	var v_99 int
	var deps_100 vm.Value
	var queue_101 vm.Value
	var queued_102 vm.Value
	var uses_103 vm.Value
	var state_104 vm.Value
	var vec__27131_106 vm.Value
	var queue_112 vm.Value
	var queued_118 vm.Value
	var arg__27284_123 vm.Value
	var arg__27289_126 vm.Value
	var v127 vm.Value
	var and__x_51 vm.Value
	var us_52 vm.Value
	var v_53 int
	var deps_54 vm.Value
	var queue_55 vm.Value
	var queued_56 vm.Value
	var uses_57 vm.Value
	var state_58 vm.Value
	var arg__27158_69 vm.Value
	var arg__27163_72 vm.Value
	var v73 vm.Value
	var and__x_59 vm.Value
	var us_60 vm.Value
	var v_61 int
	var deps_62 vm.Value
	var queue_63 vm.Value
	var queued_64 vm.Value
	var uses_65 vm.Value
	var state_66 vm.Value
	var v76 vm.Value
	var and__x_77 vm.Value
	var us_78 vm.Value
	var v_79 int
	var deps_80 vm.Value
	var queue_81 vm.Value
	var queued_82 vm.Value
	var uses_83 vm.Value
	var state_84 vm.Value
	var params_119 vm.Value
	var q_120 vm.Value
	var s_121 vm.Value
	var v199 vm.Value
	var v202 int
	var v205 vm.Value
	var v208 int
	var v146 vm.Value
	var v_129 int
	var deps_130 vm.Value
	var vec__27131_131 vm.Value
	var queue_132 vm.Value
	var queued_133 vm.Value
	var params_134 vm.Value
	var q_135 vm.Value
	var s_136 vm.Value
	var v198 vm.Value
	var v201 int
	var v204 vm.Value
	var v207 int
	var arg__27301_151 vm.Value
	var arg__27302_152 vm.Value
	var arg__27311_157 vm.Value
	var arg__27312_158 vm.Value
	var vec__27133_159 vm.Value
	var q2_165 vm.Value
	var s2_171 vm.Value
	var v173 vm.Value
	var v_137 int
	var deps_138 vm.Value
	var vec__27131_139 vm.Value
	var queue_140 vm.Value
	var queued_141 vm.Value
	var params_142 vm.Value
	var q_143 vm.Value
	var s_144 vm.Value
	var v200 vm.Value
	var v203 int
	var v206 vm.Value
	var v209 int
	var v176 vm.Value
	var v178 vm.Value
	var v_179 int
	var deps_180 vm.Value
	var vec__27131_181 vm.Value
	var queue_182 vm.Value
	var queued_183 vm.Value
	var params_184 vm.Value
	var q_185 vm.Value
	var s_186 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = uses_5, arg__27140_17, v18, v_6, deps_7, queue_8, queued_9, uses_10, v21, v_11, deps_12, queue_13, queued_14, uses_15, us_25, v_26, deps_27, queue_28, queued_29, uses_30, arg__27149_32, arg__27154_35, state_36, us_37, v_38, deps_39, queue_40, queued_41, uses_42, state_43, v93, us_44, v_45, deps_46, queue_47, queued_48, uses_49, state_50, v97, us_98, v_99, deps_100, queue_101, queued_102, uses_103, state_104, vec__27131_106, queue_112, queued_118, arg__27284_123, arg__27289_126, v127, and__x_51, us_52, v_53, deps_54, queue_55, queued_56, uses_57, state_58, arg__27158_69, arg__27163_72, v73, and__x_59, us_60, v_61, deps_62, queue_63, queued_64, uses_65, state_66, v76, and__x_77, us_78, v_79, deps_80, queue_81, queued_82, uses_83, state_84, params_119, q_120, s_121, v199, v202, v205, v208, v146, v_129, deps_130, vec__27131_131, queue_132, queued_133, params_134, q_135, s_136, v198, v201, v204, v207, arg__27301_151, arg__27302_152, arg__27311_157, arg__27312_158, vec__27133_159, q2_165, s2_171, v173, v_137, deps_138, vec__27131_139, queue_140, queued_141, params_142, q_143, s_144, v200, v203, v206, v209, v176, v178, v_179, deps_180, vec__27131_181, queue_182, queued_183, params_184, q_185, s_186
	uses_5, callErr = rt.InvokeValue(vm.Keyword("uses"), []vm.Value{arg1})
	if callErr != nil {
		return nil, callErr
	}
	arg__27140_17, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{uses_5})
	if callErr != nil {
		return nil, callErr
	}
	v18 = rt.LtValue(vm.Int(arg0), arg__27140_17)
	if v18 {
		v_6 = arg0
		deps_7 = arg1
		queue_8 = arg2
		queued_9 = arg3
		uses_10 = uses_5
		goto b1
	} else {
		v_11 = arg0
		deps_12 = arg1
		queue_13 = arg2
		queued_14 = arg3
		uses_15 = uses_5
		goto b2
	}
b1:
	;
	v21, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{uses_10, vm.Int(v_6)})
	if callErr != nil {
		return nil, callErr
	}
	us_25 = v21
	v_26 = v_6
	deps_27 = deps_7
	queue_28 = queue_8
	queued_29 = queued_9
	uses_30 = uses_10
	goto b3
b2:
	;
	us_25 = vm.NIL
	v_26 = v_11
	deps_27 = deps_12
	queue_28 = queue_13
	queued_29 = queued_14
	uses_30 = uses_15
	goto b3
b3:
	;
	arg__27149_32, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{queue_28, queued_29})
	if callErr != nil {
		return nil, callErr
	}
	arg__27154_35, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{queue_28, queued_29})
	if callErr != nil {
		return nil, callErr
	}
	state_36, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "atom").Deref(), []vm.Value{arg__27154_35})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(us_25) {
		and__x_51 = us_25
		us_52 = us_25
		v_53 = v_26
		deps_54 = deps_27
		queue_55 = queue_28
		queued_56 = queued_29
		uses_57 = uses_30
		state_58 = state_36
		goto b7
	} else {
		and__x_59 = us_25
		us_60 = us_25
		v_61 = v_26
		deps_62 = deps_27
		queue_63 = queue_28
		queued_64 = queued_29
		uses_65 = uses_30
		state_66 = state_36
		goto b8
	}
b4:
	;
	v93, callErr = rt.InvokeValue(rt.LookupVar("ir", "uses-for-each").Deref(), []vm.Value{us_37, rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var vec__27132_3 vm.Value
		var q_9 vm.Value
		var s_15 vm.Value
		var arg__27239_18 vm.Value
		var arg__27246_22 vm.Value
		var arg__27247_23 vm.Value
		var arg__27255_27 vm.Value
		var arg__27262_31 vm.Value
		var arg__27263_32 vm.Value
		var v33 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _, _, _ = vec__27132_3, q_9, s_15, arg__27239_18, arg__27246_22, arg__27247_23, arg__27255_27, arg__27262_31, arg__27263_32, v33
		vec__27132_3, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{state_43})
		if callErr != nil {
			return nil, callErr
		}
		q_9, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{vec__27132_3, vm.Int(0), vm.NIL})
		if callErr != nil {
			return nil, callErr
		}
		s_15, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{vec__27132_3, vm.Int(1), vm.NIL})
		if callErr != nil {
			return nil, callErr
		}
		arg__27239_18, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("inst"), arg0})
		if callErr != nil {
			return nil, callErr
		}
		arg__27246_22, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("inst"), arg0})
		if callErr != nil {
			return nil, callErr
		}
		arg__27247_23, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "enqueue-entry").Deref(), []vm.Value{q_9, s_15, arg__27246_22})
		if callErr != nil {
			return nil, callErr
		}
		arg__27255_27, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("inst"), arg0})
		if callErr != nil {
			return nil, callErr
		}
		arg__27262_31, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{vm.Keyword("inst"), arg0})
		if callErr != nil {
			return nil, callErr
		}
		arg__27263_32, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "enqueue-entry").Deref(), []vm.Value{q_9, s_15, arg__27262_31})
		if callErr != nil {
			return nil, callErr
		}
		v33, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "reset!").Deref(), []vm.Value{state_43, arg__27263_32})
		if callErr != nil {
			return nil, callErr
		}
		return v33, nil
	})})
	if callErr != nil {
		return nil, callErr
	}
	v97 = v93
	us_98 = us_37
	v_99 = v_38
	deps_100 = deps_39
	queue_101 = queue_40
	queued_102 = queued_41
	uses_103 = uses_42
	state_104 = state_43
	goto b6
b5:
	;
	v97 = vm.NIL
	us_98 = us_44
	v_99 = v_45
	deps_100 = deps_46
	queue_101 = queue_47
	queued_102 = queued_48
	uses_103 = uses_49
	state_104 = state_50
	goto b6
b6:
	;
	vec__27131_106, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{state_104})
	if callErr != nil {
		return nil, callErr
	}
	queue_112, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{vec__27131_106, vm.Int(0), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	queued_118, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{vec__27131_106, vm.Int(1), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	arg__27284_123, callErr = rt.InvokeValue(vm.Keyword("branch-arg-users"), []vm.Value{deps_100})
	if callErr != nil {
		return nil, callErr
	}
	arg__27289_126, callErr = rt.InvokeValue(vm.Keyword("branch-arg-users"), []vm.Value{deps_100})
	if callErr != nil {
		return nil, callErr
	}
	v127, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "get").Deref(), []vm.Value{arg__27289_126, vm.Int(v_99)})
	if callErr != nil {
		return nil, callErr
	}
	params_119 = v127
	q_120 = queue_112
	s_121 = queued_118
	v199 = vm.Keyword("param")
	v202 = 0
	v205 = vm.NIL
	v208 = 1
	goto b10
b7:
	;
	arg__27158_69, callErr = rt.InvokeValue(rt.LookupVar("ir", "uses-empty?").Deref(), []vm.Value{us_52})
	if callErr != nil {
		return nil, callErr
	}
	arg__27163_72, callErr = rt.InvokeValue(rt.LookupVar("ir", "uses-empty?").Deref(), []vm.Value{us_52})
	if callErr != nil {
		return nil, callErr
	}
	v73, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not").Deref(), []vm.Value{arg__27163_72})
	if callErr != nil {
		return nil, callErr
	}
	v76 = v73
	and__x_77 = and__x_51
	us_78 = us_52
	v_79 = v_53
	deps_80 = deps_54
	queue_81 = queue_55
	queued_82 = queued_56
	uses_83 = uses_57
	state_84 = state_58
	goto b9
b8:
	;
	v76 = and__x_59
	and__x_77 = and__x_59
	us_78 = us_60
	v_79 = v_61
	deps_80 = deps_62
	queue_81 = queue_63
	queued_82 = queued_64
	uses_83 = uses_65
	state_84 = state_66
	goto b9
b9:
	;
	if vm.IsTruthy(v76) {
		us_37 = us_78
		v_38 = v_79
		deps_39 = deps_80
		queue_40 = queue_81
		queued_41 = queued_82
		uses_42 = uses_83
		state_43 = state_84
		goto b4
	} else {
		us_44 = us_78
		v_45 = v_79
		deps_46 = deps_80
		queue_47 = queue_81
		queued_48 = queued_82
		uses_49 = uses_83
		state_50 = state_84
		goto b5
	}
b10:
	;
	v146, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{params_119})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v146) {
		v_129 = v_99
		deps_130 = deps_100
		vec__27131_131 = vec__27131_106
		queue_132 = queue_112
		queued_133 = queued_118
		params_134 = params_119
		q_135 = q_120
		s_136 = s_121
		v198 = v199
		v201 = v202
		v204 = v205
		v207 = v208
		goto b11
	} else {
		v_137 = v_99
		deps_138 = deps_100
		vec__27131_139 = vec__27131_106
		queue_140 = queue_112
		queued_141 = queued_118
		params_142 = params_119
		q_143 = q_120
		s_144 = s_121
		v200 = v199
		v203 = v202
		v206 = v205
		v209 = v208
		goto b12
	}
b11:
	;
	arg__27301_151, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{params_134})
	if callErr != nil {
		return nil, callErr
	}
	arg__27302_152, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{v198, arg__27301_151})
	if callErr != nil {
		return nil, callErr
	}
	arg__27311_157, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{params_134})
	if callErr != nil {
		return nil, callErr
	}
	arg__27312_158, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{v198, arg__27311_157})
	if callErr != nil {
		return nil, callErr
	}
	vec__27133_159, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "enqueue-entry").Deref(), []vm.Value{q_135, s_136, arg__27312_158})
	if callErr != nil {
		return nil, callErr
	}
	q2_165, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{vec__27133_159, vm.Int(v201), v204})
	if callErr != nil {
		return nil, callErr
	}
	s2_171, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{vec__27133_159, vm.Int(v207), v204})
	if callErr != nil {
		return nil, callErr
	}
	v173, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{params_134})
	if callErr != nil {
		return nil, callErr
	}
	params_119 = v173
	q_120 = q2_165
	s_121 = s2_171
	v199 = v198
	v202 = v201
	v205 = v204
	v208 = v207
	goto b10
b12:
	;
	v176, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{q_143, s_144})
	if callErr != nil {
		return nil, callErr
	}
	v178 = v176
	v_179 = v_137
	deps_180 = deps_138
	vec__27131_181 = vec__27131_139
	queue_182 = queue_140
	queued_183 = queued_141
	params_184 = params_142
	q_185 = q_143
	s_186 = s_144
	goto b13
b13:
	;
	return v178, nil
}
func infer_one(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var op_3 vm.Value
	var refs_5 vm.Value
	var v15 bool
	var inst_6 vm.Value
	var f_7 vm.Value
	var op_8 vm.Value
	var refs_9 vm.Value
	var arg__27350_18 vm.Value
	var arg__27357_21 vm.Value
	var v22 vm.Value
	var inst_10 vm.Value
	var f_11 vm.Value
	var op_12 vm.Value
	var refs_13 vm.Value
	var v33 bool
	var v225 vm.Value
	var inst_226 vm.Value
	var f_227 vm.Value
	var op_228 vm.Value
	var refs_229 vm.Value
	var inst_24 vm.Value
	var f_25 vm.Value
	var op_26 vm.Value
	var refs_27 vm.Value
	var v36 vm.Value
	var inst_28 vm.Value
	var f_29 vm.Value
	var op_30 vm.Value
	var refs_31 vm.Value
	var v47 bool
	var v219 vm.Value
	var inst_220 vm.Value
	var f_221 vm.Value
	var op_222 vm.Value
	var refs_223 vm.Value
	var inst_38 vm.Value
	var f_39 vm.Value
	var op_40 vm.Value
	var refs_41 vm.Value
	var inst_42 vm.Value
	var f_43 vm.Value
	var op_44 vm.Value
	var refs_45 vm.Value
	var v60 bool
	var v213 vm.Value
	var inst_214 vm.Value
	var f_215 vm.Value
	var op_216 vm.Value
	var refs_217 vm.Value
	var inst_51 vm.Value
	var f_52 vm.Value
	var op_53 vm.Value
	var refs_54 vm.Value
	var inst_55 vm.Value
	var f_56 vm.Value
	var op_57 vm.Value
	var refs_58 vm.Value
	var v73 bool
	var v207 vm.Value
	var inst_208 vm.Value
	var f_209 vm.Value
	var op_210 vm.Value
	var refs_211 vm.Value
	var inst_64 vm.Value
	var f_65 vm.Value
	var op_66 vm.Value
	var refs_67 vm.Value
	var inst_68 vm.Value
	var f_69 vm.Value
	var op_70 vm.Value
	var refs_71 vm.Value
	var v86 bool
	var v201 vm.Value
	var inst_202 vm.Value
	var f_203 vm.Value
	var op_204 vm.Value
	var refs_205 vm.Value
	var inst_77 vm.Value
	var f_78 vm.Value
	var op_79 vm.Value
	var refs_80 vm.Value
	var arg__27378_89 vm.Value
	var arg__27385_92 vm.Value
	var v93 vm.Value
	var inst_81 vm.Value
	var f_82 vm.Value
	var op_83 vm.Value
	var refs_84 vm.Value
	var v106 vm.Value
	var v195 vm.Value
	var inst_196 vm.Value
	var f_197 vm.Value
	var op_198 vm.Value
	var refs_199 vm.Value
	var inst_95 vm.Value
	var f_96 vm.Value
	var op_97 vm.Value
	var refs_98 vm.Value
	var inst_99 vm.Value
	var f_100 vm.Value
	var op_101 vm.Value
	var refs_102 vm.Value
	var v121 vm.Value
	var v189 vm.Value
	var inst_190 vm.Value
	var f_191 vm.Value
	var op_192 vm.Value
	var refs_193 vm.Value
	var inst_110 vm.Value
	var f_111 vm.Value
	var op_112 vm.Value
	var refs_113 vm.Value
	var inst_114 vm.Value
	var f_115 vm.Value
	var op_116 vm.Value
	var refs_117 vm.Value
	var v136 vm.Value
	var v183 vm.Value
	var inst_184 vm.Value
	var f_185 vm.Value
	var op_186 vm.Value
	var refs_187 vm.Value
	var inst_125 vm.Value
	var f_126 vm.Value
	var op_127 vm.Value
	var refs_128 vm.Value
	var arg__27416_145 vm.Value
	var arg__27433_154 vm.Value
	var v155 vm.Value
	var inst_129 vm.Value
	var f_130 vm.Value
	var op_131 vm.Value
	var refs_132 vm.Value
	var v177 vm.Value
	var inst_178 vm.Value
	var f_179 vm.Value
	var op_180 vm.Value
	var refs_181 vm.Value
	var inst_157 vm.Value
	var f_158 vm.Value
	var op_159 vm.Value
	var refs_160 vm.Value
	var inst_161 vm.Value
	var f_162 vm.Value
	var op_163 vm.Value
	var refs_164 vm.Value
	var v171 vm.Value
	var inst_172 vm.Value
	var f_173 vm.Value
	var op_174 vm.Value
	var refs_175 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = op_3, refs_5, v15, inst_6, f_7, op_8, refs_9, arg__27350_18, arg__27357_21, v22, inst_10, f_11, op_12, refs_13, v33, v225, inst_226, f_227, op_228, refs_229, inst_24, f_25, op_26, refs_27, v36, inst_28, f_29, op_30, refs_31, v47, v219, inst_220, f_221, op_222, refs_223, inst_38, f_39, op_40, refs_41, inst_42, f_43, op_44, refs_45, v60, v213, inst_214, f_215, op_216, refs_217, inst_51, f_52, op_53, refs_54, inst_55, f_56, op_57, refs_58, v73, v207, inst_208, f_209, op_210, refs_211, inst_64, f_65, op_66, refs_67, inst_68, f_69, op_70, refs_71, v86, v201, inst_202, f_203, op_204, refs_205, inst_77, f_78, op_79, refs_80, arg__27378_89, arg__27385_92, v93, inst_81, f_82, op_83, refs_84, v106, v195, inst_196, f_197, op_198, refs_199, inst_95, f_96, op_97, refs_98, inst_99, f_100, op_101, refs_102, v121, v189, inst_190, f_191, op_192, refs_193, inst_110, f_111, op_112, refs_113, inst_114, f_115, op_116, refs_117, v136, v183, inst_184, f_185, op_186, refs_187, inst_125, f_126, op_127, refs_128, arg__27416_145, arg__27433_154, v155, inst_129, f_130, op_131, refs_132, v177, inst_178, f_179, op_180, refs_181, inst_157, f_158, op_159, refs_160, inst_161, f_162, op_163, refs_164, v171, inst_172, f_173, op_174, refs_175
	op_3, callErr = rt.InvokeValue(rt.LookupVar("ir", "op").Deref(), []vm.Value{arg0, arg1})
	if callErr != nil {
		return nil, callErr
	}
	refs_5, callErr = rt.InvokeValue(rt.LookupVar("ir", "refs").Deref(), []vm.Value{arg0, arg1})
	if callErr != nil {
		return nil, callErr
	}
	v15 = op_3 == vm.Keyword("const")
	if v15 {
		inst_6 = arg0
		f_7 = arg1
		op_8 = op_3
		refs_9 = refs_5
		goto b1
	} else {
		inst_10 = arg0
		f_11 = arg1
		op_12 = op_3
		refs_13 = refs_5
		goto b2
	}
b1:
	;
	arg__27350_18, callErr = rt.InvokeValue(rt.LookupVar("ir", "aux").Deref(), []vm.Value{inst_6, f_7})
	if callErr != nil {
		return nil, callErr
	}
	arg__27357_21, callErr = rt.InvokeValue(rt.LookupVar("ir", "aux").Deref(), []vm.Value{inst_6, f_7})
	if callErr != nil {
		return nil, callErr
	}
	v22, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "classify-const").Deref(), []vm.Value{arg__27357_21})
	if callErr != nil {
		return nil, callErr
	}
	v225 = v22
	inst_226 = inst_6
	f_227 = f_7
	op_228 = op_8
	refs_229 = refs_9
	goto b3
b2:
	;
	v33 = op_12 == vm.Keyword("load-arg")
	if v33 {
		inst_24 = inst_10
		f_25 = f_11
		op_26 = op_12
		refs_27 = refs_13
		goto b4
	} else {
		inst_28 = inst_10
		f_29 = f_11
		op_30 = op_12
		refs_31 = refs_13
		goto b5
	}
b3:
	;
	return v225, nil
b4:
	;
	v36, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "arg-seed").Deref(), []vm.Value{inst_24, f_25})
	if callErr != nil {
		return nil, callErr
	}
	v219 = v36
	inst_220 = inst_24
	f_221 = f_25
	op_222 = op_26
	refs_223 = refs_27
	goto b6
b5:
	;
	v47 = op_30 == vm.Keyword("load-var")
	if v47 {
		inst_38 = inst_28
		f_39 = f_29
		op_40 = op_30
		refs_41 = refs_31
		goto b7
	} else {
		inst_42 = inst_28
		f_43 = f_29
		op_44 = op_30
		refs_45 = refs_31
		goto b8
	}
b6:
	;
	v225 = v219
	inst_226 = inst_220
	f_227 = f_221
	op_228 = op_222
	refs_229 = refs_223
	goto b3
b7:
	;
	v213 = vm.Keyword("unknown")
	inst_214 = inst_38
	f_215 = f_39
	op_216 = op_40
	refs_217 = refs_41
	goto b9
b8:
	;
	v60 = op_44 == vm.Keyword("load-closed")
	if v60 {
		inst_51 = inst_42
		f_52 = f_43
		op_53 = op_44
		refs_54 = refs_45
		goto b10
	} else {
		inst_55 = inst_42
		f_56 = f_43
		op_57 = op_44
		refs_58 = refs_45
		goto b11
	}
b9:
	;
	v219 = v213
	inst_220 = inst_214
	f_221 = f_215
	op_222 = op_216
	refs_223 = refs_217
	goto b6
b10:
	;
	v207 = vm.Keyword("unknown")
	inst_208 = inst_51
	f_209 = f_52
	op_210 = op_53
	refs_211 = refs_54
	goto b12
b11:
	;
	v73 = op_57 == vm.Keyword("call")
	if v73 {
		inst_64 = inst_55
		f_65 = f_56
		op_66 = op_57
		refs_67 = refs_58
		goto b13
	} else {
		inst_68 = inst_55
		f_69 = f_56
		op_70 = op_57
		refs_71 = refs_58
		goto b14
	}
b12:
	;
	v213 = v207
	inst_214 = inst_208
	f_215 = f_209
	op_216 = op_210
	refs_217 = refs_211
	goto b9
b13:
	;
	v201 = vm.Keyword("unknown")
	inst_202 = inst_64
	f_203 = f_65
	op_204 = op_66
	refs_205 = refs_67
	goto b15
b14:
	;
	v86 = op_70 == vm.Keyword("block-arg")
	if v86 {
		inst_77 = inst_68
		f_78 = f_69
		op_79 = op_70
		refs_80 = refs_71
		goto b16
	} else {
		inst_81 = inst_68
		f_82 = f_69
		op_83 = op_70
		refs_84 = refs_71
		goto b17
	}
b15:
	;
	v207 = v201
	inst_208 = inst_202
	f_209 = f_203
	op_210 = op_204
	refs_211 = refs_205
	goto b12
b16:
	;
	arg__27378_89, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{inst_77, f_78})
	if callErr != nil {
		return nil, callErr
	}
	arg__27385_92, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{inst_77, f_78})
	if callErr != nil {
		return nil, callErr
	}
	v93, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "normalize-type").Deref(), []vm.Value{arg__27385_92})
	if callErr != nil {
		return nil, callErr
	}
	v195 = v93
	inst_196 = inst_77
	f_197 = f_78
	op_198 = op_79
	refs_199 = refs_80
	goto b18
b17:
	;
	v106, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "contains?").Deref(), []vm.Value{rt.LookupVar("ir.passes.typeinfer", "terminator-ops").Deref(), op_83})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v106) {
		inst_95 = inst_81
		f_96 = f_82
		op_97 = op_83
		refs_98 = refs_84
		goto b19
	} else {
		inst_99 = inst_81
		f_100 = f_82
		op_101 = op_83
		refs_102 = refs_84
		goto b20
	}
b18:
	;
	v201 = v195
	inst_202 = inst_196
	f_203 = f_197
	op_204 = op_198
	refs_205 = refs_199
	goto b15
b19:
	;
	v189 = vm.Keyword("nil")
	inst_190 = inst_95
	f_191 = f_96
	op_192 = op_97
	refs_193 = refs_98
	goto b21
b20:
	;
	v121, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "contains?").Deref(), []vm.Value{rt.LookupVar("ir.passes.typeinfer", "comparison-ops").Deref(), op_101})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v121) {
		inst_110 = inst_99
		f_111 = f_100
		op_112 = op_101
		refs_113 = refs_102
		goto b22
	} else {
		inst_114 = inst_99
		f_115 = f_100
		op_116 = op_101
		refs_117 = refs_102
		goto b23
	}
b21:
	;
	v195 = v189
	inst_196 = inst_190
	f_197 = f_191
	op_198 = op_192
	refs_199 = refs_193
	goto b18
b22:
	;
	v183 = vm.Keyword("bool")
	inst_184 = inst_110
	f_185 = f_111
	op_186 = op_112
	refs_187 = refs_113
	goto b24
b23:
	;
	v136, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "contains?").Deref(), []vm.Value{rt.LookupVar("ir.passes.typeinfer", "numeric-ops").Deref(), op_116})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v136) {
		inst_125 = inst_114
		f_126 = f_115
		op_127 = op_116
		refs_128 = refs_117
		goto b25
	} else {
		inst_129 = inst_114
		f_130 = f_115
		op_131 = op_116
		refs_132 = refs_117
		goto b26
	}
b24:
	;
	v189 = v183
	inst_190 = inst_184
	f_191 = f_185
	op_192 = op_186
	refs_193 = refs_187
	goto b21
b25:
	;
	arg__27416_145, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_126})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), refs_128})
	if callErr != nil {
		return nil, callErr
	}
	arg__27433_154, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "map").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v3 vm.Value
		var callErr error
		_ = v3
		v3, callErr = rt.InvokeValue(rt.LookupVar("ir", "type-of").Deref(), []vm.Value{arg0, f_126})
		if callErr != nil {
			return nil, callErr
		}
		return v3, nil
	}), refs_128})
	if callErr != nil {
		return nil, callErr
	}
	v155, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "numeric-op-type").Deref(), []vm.Value{arg__27433_154})
	if callErr != nil {
		return nil, callErr
	}
	v177 = v155
	inst_178 = inst_125
	f_179 = f_126
	op_180 = op_127
	refs_181 = refs_128
	goto b27
b26:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		inst_157 = inst_129
		f_158 = f_130
		op_159 = op_131
		refs_160 = refs_132
		goto b28
	} else {
		inst_161 = inst_129
		f_162 = f_130
		op_163 = op_131
		refs_164 = refs_132
		goto b29
	}
b27:
	;
	v183 = v177
	inst_184 = inst_178
	f_185 = f_179
	op_186 = op_180
	refs_187 = refs_181
	goto b24
b28:
	;
	v171 = vm.Keyword("unknown")
	inst_172 = inst_157
	f_173 = f_158
	op_174 = op_159
	refs_175 = refs_160
	goto b30
b29:
	;
	v171 = vm.NIL
	inst_172 = inst_161
	f_173 = f_162
	op_174 = op_163
	refs_175 = refs_164
	goto b30
b30:
	;
	v177 = v171
	inst_178 = inst_172
	f_179 = f_173
	op_180 = op_174
	refs_181 = refs_175
	goto b27
}
func analyze_block_BANG_(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var v5 vm.Value
	var params_7 vm.Value
	var ps_8 vm.Value
	var idx_9 int
	var changed_QMARK__10 vm.Value
	var f_11 vm.Value
	var bid_12 vm.Value
	var v29 vm.Value
	var params_16 vm.Value
	var ps_17 vm.Value
	var idx_18 int
	var changed_QMARK__19 vm.Value
	var f_20 vm.Value
	var bid_21 vm.Value
	var param_id_32 vm.Value
	var param_type_34 vm.Value
	var param_changed_QMARK__36 vm.Value
	var v38 vm.Value
	var v39 int
	var params_22 vm.Value
	var ps_23 vm.Value
	var idx_24 int
	var changed_QMARK__25 vm.Value
	var f_26 vm.Value
	var bid_27 vm.Value
	var changed_param_QMARK__76 vm.Value
	var params_77 vm.Value
	var ps_78 vm.Value
	var idx_79 int
	var changed_QMARK__80 vm.Value
	var f_81 vm.Value
	var bid_82 vm.Value
	var v88 vm.Value
	var params_40 vm.Value
	var ps_41 vm.Value
	var idx_42 int
	var or__x_43 vm.Value
	var changed_QMARK__44 vm.Value
	var f_45 vm.Value
	var bid_46 vm.Value
	var param_id_47 vm.Value
	var param_type_48 vm.Value
	var param_changed_QMARK__49 vm.Value
	var params_50 vm.Value
	var ps_51 vm.Value
	var idx_52 int
	var or__x_53 vm.Value
	var changed_QMARK__54 vm.Value
	var f_55 vm.Value
	var bid_56 vm.Value
	var param_id_57 vm.Value
	var param_type_58 vm.Value
	var param_changed_QMARK__59 vm.Value
	var v63 vm.Value
	var params_64 vm.Value
	var ps_65 vm.Value
	var idx_66 int
	var or__x_67 vm.Value
	var changed_QMARK__68 vm.Value
	var f_69 vm.Value
	var bid_70 vm.Value
	var param_id_71 vm.Value
	var param_type_72 vm.Value
	var param_changed_QMARK__73 vm.Value
	var insts_83 vm.Value
	var changed_QMARK__84 vm.Value
	var f_85 vm.Value
	var changed_QMARK__86 vm.Value
	var v108 vm.Value
	var changed_param_QMARK__91 vm.Value
	var params_92 vm.Value
	var ps_93 vm.Value
	var idx_94 int
	var bid_95 vm.Value
	var insts_96 vm.Value
	var f_97 vm.Value
	var changed_QMARK__98 vm.Value
	var inst_111 vm.Value
	var arg__28264_113 vm.Value
	var arg__28273_116 vm.Value
	var inst_changed_QMARK__117 vm.Value
	var v119 vm.Value
	var changed_param_QMARK__99 vm.Value
	var params_100 vm.Value
	var ps_101 vm.Value
	var idx_102 int
	var bid_103 vm.Value
	var insts_104 vm.Value
	var f_105 vm.Value
	var changed_QMARK__106 vm.Value
	var changed_inst_QMARK__159 vm.Value
	var changed_param_QMARK__160 vm.Value
	var params_161 vm.Value
	var ps_162 vm.Value
	var idx_163 int
	var bid_164 vm.Value
	var insts_165 vm.Value
	var f_166 vm.Value
	var changed_QMARK__167 vm.Value
	var term_169 vm.Value
	var v191 vm.Value
	var changed_param_QMARK__120 vm.Value
	var params_121 vm.Value
	var ps_122 vm.Value
	var idx_123 int
	var bid_124 vm.Value
	var insts_125 vm.Value
	var f_126 vm.Value
	var or__x_127 vm.Value
	var changed_QMARK__128 vm.Value
	var inst_129 vm.Value
	var inst_changed_QMARK__130 vm.Value
	var changed_param_QMARK__131 vm.Value
	var params_132 vm.Value
	var ps_133 vm.Value
	var idx_134 int
	var bid_135 vm.Value
	var insts_136 vm.Value
	var f_137 vm.Value
	var or__x_138 vm.Value
	var changed_QMARK__139 vm.Value
	var inst_140 vm.Value
	var inst_changed_QMARK__141 vm.Value
	var v145 vm.Value
	var changed_param_QMARK__146 vm.Value
	var params_147 vm.Value
	var ps_148 vm.Value
	var idx_149 int
	var bid_150 vm.Value
	var insts_151 vm.Value
	var f_152 vm.Value
	var or__x_153 vm.Value
	var changed_QMARK__154 vm.Value
	var inst_155 vm.Value
	var inst_changed_QMARK__156 vm.Value
	var changed_inst_QMARK__170 vm.Value
	var changed_param_QMARK__171 vm.Value
	var params_172 vm.Value
	var ps_173 vm.Value
	var idx_174 int
	var bid_175 vm.Value
	var insts_176 vm.Value
	var f_177 vm.Value
	var changed_QMARK__178 vm.Value
	var term_179 vm.Value
	var changed_inst_QMARK__180 vm.Value
	var changed_param_QMARK__181 vm.Value
	var params_182 vm.Value
	var ps_183 vm.Value
	var idx_184 int
	var bid_185 vm.Value
	var insts_186 vm.Value
	var f_187 vm.Value
	var changed_QMARK__188 vm.Value
	var term_189 vm.Value
	var arg__28292_196 vm.Value
	var arg__28301_199 vm.Value
	var v200 vm.Value
	var changed_term_QMARK__202 vm.Value
	var changed_inst_QMARK__203 vm.Value
	var changed_param_QMARK__204 vm.Value
	var params_205 vm.Value
	var ps_206 vm.Value
	var idx_207 int
	var bid_208 vm.Value
	var insts_209 vm.Value
	var f_210 vm.Value
	var changed_QMARK__211 vm.Value
	var term_212 vm.Value
	var changed_term_QMARK__213 vm.Value
	var changed_inst_QMARK__214 vm.Value
	var or__x_215 vm.Value
	var changed_param_QMARK__216 vm.Value
	var params_217 vm.Value
	var ps_218 vm.Value
	var idx_219 int
	var bid_220 vm.Value
	var insts_221 vm.Value
	var f_222 vm.Value
	var changed_QMARK__223 vm.Value
	var term_224 vm.Value
	var changed_term_QMARK__225 vm.Value
	var changed_inst_QMARK__226 vm.Value
	var or__x_227 vm.Value
	var changed_param_QMARK__228 vm.Value
	var params_229 vm.Value
	var ps_230 vm.Value
	var idx_231 int
	var bid_232 vm.Value
	var insts_233 vm.Value
	var f_234 vm.Value
	var changed_QMARK__235 vm.Value
	var term_236 vm.Value
	var v280 vm.Value
	var changed_term_QMARK__281 vm.Value
	var changed_inst_QMARK__282 vm.Value
	var or__x_283 vm.Value
	var changed_param_QMARK__284 vm.Value
	var params_285 vm.Value
	var ps_286 vm.Value
	var idx_287 int
	var bid_288 vm.Value
	var insts_289 vm.Value
	var f_290 vm.Value
	var changed_QMARK__291 vm.Value
	var term_292 vm.Value
	var changed_term_QMARK__239 vm.Value
	var or__x_240 vm.Value
	var changed_inst_QMARK__241 vm.Value
	var changed_param_QMARK__242 vm.Value
	var params_243 vm.Value
	var ps_244 vm.Value
	var idx_245 int
	var bid_246 vm.Value
	var insts_247 vm.Value
	var f_248 vm.Value
	var changed_QMARK__249 vm.Value
	var term_250 vm.Value
	var changed_term_QMARK__251 vm.Value
	var or__x_252 vm.Value
	var changed_inst_QMARK__253 vm.Value
	var changed_param_QMARK__254 vm.Value
	var params_255 vm.Value
	var ps_256 vm.Value
	var idx_257 int
	var bid_258 vm.Value
	var insts_259 vm.Value
	var f_260 vm.Value
	var changed_QMARK__261 vm.Value
	var term_262 vm.Value
	var v266 vm.Value
	var changed_term_QMARK__267 vm.Value
	var or__x_268 vm.Value
	var changed_inst_QMARK__269 vm.Value
	var changed_param_QMARK__270 vm.Value
	var params_271 vm.Value
	var ps_272 vm.Value
	var idx_273 int
	var bid_274 vm.Value
	var insts_275 vm.Value
	var f_276 vm.Value
	var changed_QMARK__277 vm.Value
	var term_278 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = v5, params_7, ps_8, idx_9, changed_QMARK__10, f_11, bid_12, v29, params_16, ps_17, idx_18, changed_QMARK__19, f_20, bid_21, param_id_32, param_type_34, param_changed_QMARK__36, v38, v39, params_22, ps_23, idx_24, changed_QMARK__25, f_26, bid_27, changed_param_QMARK__76, params_77, ps_78, idx_79, changed_QMARK__80, f_81, bid_82, v88, params_40, ps_41, idx_42, or__x_43, changed_QMARK__44, f_45, bid_46, param_id_47, param_type_48, param_changed_QMARK__49, params_50, ps_51, idx_52, or__x_53, changed_QMARK__54, f_55, bid_56, param_id_57, param_type_58, param_changed_QMARK__59, v63, params_64, ps_65, idx_66, or__x_67, changed_QMARK__68, f_69, bid_70, param_id_71, param_type_72, param_changed_QMARK__73, insts_83, changed_QMARK__84, f_85, changed_QMARK__86, v108, changed_param_QMARK__91, params_92, ps_93, idx_94, bid_95, insts_96, f_97, changed_QMARK__98, inst_111, arg__28264_113, arg__28273_116, inst_changed_QMARK__117, v119, changed_param_QMARK__99, params_100, ps_101, idx_102, bid_103, insts_104, f_105, changed_QMARK__106, changed_inst_QMARK__159, changed_param_QMARK__160, params_161, ps_162, idx_163, bid_164, insts_165, f_166, changed_QMARK__167, term_169, v191, changed_param_QMARK__120, params_121, ps_122, idx_123, bid_124, insts_125, f_126, or__x_127, changed_QMARK__128, inst_129, inst_changed_QMARK__130, changed_param_QMARK__131, params_132, ps_133, idx_134, bid_135, insts_136, f_137, or__x_138, changed_QMARK__139, inst_140, inst_changed_QMARK__141, v145, changed_param_QMARK__146, params_147, ps_148, idx_149, bid_150, insts_151, f_152, or__x_153, changed_QMARK__154, inst_155, inst_changed_QMARK__156, changed_inst_QMARK__170, changed_param_QMARK__171, params_172, ps_173, idx_174, bid_175, insts_176, f_177, changed_QMARK__178, term_179, changed_inst_QMARK__180, changed_param_QMARK__181, params_182, ps_183, idx_184, bid_185, insts_186, f_187, changed_QMARK__188, term_189, arg__28292_196, arg__28301_199, v200, changed_term_QMARK__202, changed_inst_QMARK__203, changed_param_QMARK__204, params_205, ps_206, idx_207, bid_208, insts_209, f_210, changed_QMARK__211, term_212, changed_term_QMARK__213, changed_inst_QMARK__214, or__x_215, changed_param_QMARK__216, params_217, ps_218, idx_219, bid_220, insts_221, f_222, changed_QMARK__223, term_224, changed_term_QMARK__225, changed_inst_QMARK__226, or__x_227, changed_param_QMARK__228, params_229, ps_230, idx_231, bid_232, insts_233, f_234, changed_QMARK__235, term_236, v280, changed_term_QMARK__281, changed_inst_QMARK__282, or__x_283, changed_param_QMARK__284, params_285, ps_286, idx_287, bid_288, insts_289, f_290, changed_QMARK__291, term_292, changed_term_QMARK__239, or__x_240, changed_inst_QMARK__241, changed_param_QMARK__242, params_243, ps_244, idx_245, bid_246, insts_247, f_248, changed_QMARK__249, term_250, changed_term_QMARK__251, or__x_252, changed_inst_QMARK__253, changed_param_QMARK__254, params_255, ps_256, idx_257, bid_258, insts_259, f_260, changed_QMARK__261, term_262, v266, changed_term_QMARK__267, or__x_268, changed_inst_QMARK__269, changed_param_QMARK__270, params_271, ps_272, idx_273, bid_274, insts_275, f_276, changed_QMARK__277, term_278
	v5, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "ti-inc!").Deref(), []vm.Value{vm.Keyword("analyze-block")})
	if callErr != nil {
		return nil, callErr
	}
	params_7, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-params").Deref(), []vm.Value{arg0, arg1})
	if callErr != nil {
		return nil, callErr
	}
	ps_8 = params_7
	idx_9 = 0
	changed_QMARK__10 = vm.Boolean(false)
	f_11 = arg1
	bid_12 = arg0
	goto b1
b1:
	;
	v29, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{ps_8})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v29) {
		params_16 = params_7
		ps_17 = ps_8
		idx_18 = idx_9
		changed_QMARK__19 = changed_QMARK__10
		f_20 = f_11
		bid_21 = bid_12
		goto b2
	} else {
		params_22 = params_7
		ps_23 = ps_8
		idx_24 = idx_9
		changed_QMARK__25 = changed_QMARK__10
		f_26 = f_11
		bid_27 = bid_12
		goto b3
	}
b2:
	;
	param_id_32, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{ps_17})
	if callErr != nil {
		return nil, callErr
	}
	param_type_34, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "infer-block-param").Deref(), []vm.Value{bid_21, vm.Int(idx_18), f_20})
	if callErr != nil {
		return nil, callErr
	}
	param_changed_QMARK__36, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "set-type-if-changed!").Deref(), []vm.Value{f_20, param_id_32, param_type_34})
	if callErr != nil {
		return nil, callErr
	}
	v38, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{ps_17})
	if callErr != nil {
		return nil, callErr
	}
	v39 = idx_18 + 1
	if vm.IsTruthy(changed_QMARK__19) {
		params_40 = params_16
		ps_41 = ps_17
		idx_42 = idx_18
		or__x_43 = changed_QMARK__19
		changed_QMARK__44 = changed_QMARK__19
		f_45 = f_20
		bid_46 = bid_21
		param_id_47 = param_id_32
		param_type_48 = param_type_34
		param_changed_QMARK__49 = param_changed_QMARK__36
		goto b5
	} else {
		params_50 = params_16
		ps_51 = ps_17
		idx_52 = idx_18
		or__x_53 = changed_QMARK__19
		changed_QMARK__54 = changed_QMARK__19
		f_55 = f_20
		bid_56 = bid_21
		param_id_57 = param_id_32
		param_type_58 = param_type_34
		param_changed_QMARK__59 = param_changed_QMARK__36
		goto b6
	}
b3:
	;
	changed_param_QMARK__76 = changed_QMARK__25
	params_77 = params_22
	ps_78 = ps_23
	idx_79 = idx_24
	changed_QMARK__80 = changed_QMARK__25
	f_81 = f_26
	bid_82 = bid_27
	goto b4
b4:
	;
	v88, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-insts").Deref(), []vm.Value{bid_82, f_81})
	if callErr != nil {
		return nil, callErr
	}
	insts_83 = v88
	changed_QMARK__84 = vm.Boolean(false)
	f_85 = f_81
	changed_QMARK__86 = changed_QMARK__80
	goto b8
b5:
	;
	v63 = or__x_43
	params_64 = params_40
	ps_65 = ps_41
	idx_66 = idx_42
	or__x_67 = or__x_43
	changed_QMARK__68 = changed_QMARK__44
	f_69 = f_45
	bid_70 = bid_46
	param_id_71 = param_id_47
	param_type_72 = param_type_48
	param_changed_QMARK__73 = param_changed_QMARK__49
	goto b7
b6:
	;
	v63 = param_changed_QMARK__59
	params_64 = params_50
	ps_65 = ps_51
	idx_66 = idx_52
	or__x_67 = or__x_53
	changed_QMARK__68 = changed_QMARK__54
	f_69 = f_55
	bid_70 = bid_56
	param_id_71 = param_id_57
	param_type_72 = param_type_58
	param_changed_QMARK__73 = param_changed_QMARK__59
	goto b7
b7:
	;
	ps_8 = v38
	idx_9 = v39
	changed_QMARK__10 = v63
	f_11 = f_20
	bid_12 = bid_21
	goto b1
b8:
	;
	v108, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{insts_83})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v108) {
		changed_param_QMARK__91 = changed_param_QMARK__76
		params_92 = params_77
		ps_93 = ps_78
		idx_94 = idx_79
		bid_95 = bid_82
		insts_96 = insts_83
		f_97 = f_85
		changed_QMARK__98 = changed_QMARK__86
		goto b9
	} else {
		changed_param_QMARK__99 = changed_param_QMARK__76
		params_100 = params_77
		ps_101 = ps_78
		idx_102 = idx_79
		bid_103 = bid_82
		insts_104 = insts_83
		f_105 = f_85
		changed_QMARK__106 = changed_QMARK__86
		goto b10
	}
b9:
	;
	inst_111, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{insts_96})
	if callErr != nil {
		return nil, callErr
	}
	arg__28264_113, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "infer-one").Deref(), []vm.Value{inst_111, f_97})
	if callErr != nil {
		return nil, callErr
	}
	arg__28273_116, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "infer-one").Deref(), []vm.Value{inst_111, f_97})
	if callErr != nil {
		return nil, callErr
	}
	inst_changed_QMARK__117, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "set-type-if-changed!").Deref(), []vm.Value{f_97, inst_111, arg__28273_116})
	if callErr != nil {
		return nil, callErr
	}
	v119, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{insts_96})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(changed_QMARK__98) {
		changed_param_QMARK__120 = changed_param_QMARK__91
		params_121 = params_92
		ps_122 = ps_93
		idx_123 = idx_94
		bid_124 = bid_95
		insts_125 = insts_96
		f_126 = f_97
		or__x_127 = changed_QMARK__98
		changed_QMARK__128 = changed_QMARK__98
		inst_129 = inst_111
		inst_changed_QMARK__130 = inst_changed_QMARK__117
		goto b12
	} else {
		changed_param_QMARK__131 = changed_param_QMARK__91
		params_132 = params_92
		ps_133 = ps_93
		idx_134 = idx_94
		bid_135 = bid_95
		insts_136 = insts_96
		f_137 = f_97
		or__x_138 = changed_QMARK__98
		changed_QMARK__139 = changed_QMARK__98
		inst_140 = inst_111
		inst_changed_QMARK__141 = inst_changed_QMARK__117
		goto b13
	}
b10:
	;
	changed_inst_QMARK__159 = changed_QMARK__106
	changed_param_QMARK__160 = changed_param_QMARK__99
	params_161 = params_100
	ps_162 = ps_101
	idx_163 = idx_102
	bid_164 = bid_103
	insts_165 = insts_104
	f_166 = f_105
	changed_QMARK__167 = changed_QMARK__106
	goto b11
b11:
	;
	term_169, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-term").Deref(), []vm.Value{bid_164, f_166})
	if callErr != nil {
		return nil, callErr
	}
	v191, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nil?").Deref(), []vm.Value{term_169})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v191) {
		changed_inst_QMARK__170 = changed_inst_QMARK__159
		changed_param_QMARK__171 = changed_param_QMARK__160
		params_172 = params_161
		ps_173 = ps_162
		idx_174 = idx_163
		bid_175 = bid_164
		insts_176 = insts_165
		f_177 = f_166
		changed_QMARK__178 = changed_QMARK__167
		term_179 = term_169
		goto b15
	} else {
		changed_inst_QMARK__180 = changed_inst_QMARK__159
		changed_param_QMARK__181 = changed_param_QMARK__160
		params_182 = params_161
		ps_183 = ps_162
		idx_184 = idx_163
		bid_185 = bid_164
		insts_186 = insts_165
		f_187 = f_166
		changed_QMARK__188 = changed_QMARK__167
		term_189 = term_169
		goto b16
	}
b12:
	;
	v145 = or__x_127
	changed_param_QMARK__146 = changed_param_QMARK__120
	params_147 = params_121
	ps_148 = ps_122
	idx_149 = idx_123
	bid_150 = bid_124
	insts_151 = insts_125
	f_152 = f_126
	or__x_153 = or__x_127
	changed_QMARK__154 = changed_QMARK__128
	inst_155 = inst_129
	inst_changed_QMARK__156 = inst_changed_QMARK__130
	goto b14
b13:
	;
	v145 = inst_changed_QMARK__141
	changed_param_QMARK__146 = changed_param_QMARK__131
	params_147 = params_132
	ps_148 = ps_133
	idx_149 = idx_134
	bid_150 = bid_135
	insts_151 = insts_136
	f_152 = f_137
	or__x_153 = or__x_138
	changed_QMARK__154 = changed_QMARK__139
	inst_155 = inst_140
	inst_changed_QMARK__156 = inst_changed_QMARK__141
	goto b14
b14:
	;
	insts_83 = v119
	changed_QMARK__84 = v145
	f_85 = f_97
	changed_QMARK__86 = changed_QMARK__98
	goto b8
b15:
	;
	changed_term_QMARK__202 = vm.Boolean(false)
	changed_inst_QMARK__203 = changed_inst_QMARK__170
	changed_param_QMARK__204 = changed_param_QMARK__171
	params_205 = params_172
	ps_206 = ps_173
	idx_207 = idx_174
	bid_208 = bid_175
	insts_209 = insts_176
	f_210 = f_177
	changed_QMARK__211 = changed_QMARK__178
	term_212 = term_179
	goto b17
b16:
	;
	arg__28292_196, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "infer-one").Deref(), []vm.Value{term_189, f_187})
	if callErr != nil {
		return nil, callErr
	}
	arg__28301_199, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "infer-one").Deref(), []vm.Value{term_189, f_187})
	if callErr != nil {
		return nil, callErr
	}
	v200, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "set-type-if-changed!").Deref(), []vm.Value{f_187, term_189, arg__28301_199})
	if callErr != nil {
		return nil, callErr
	}
	changed_term_QMARK__202 = v200
	changed_inst_QMARK__203 = changed_inst_QMARK__180
	changed_param_QMARK__204 = changed_param_QMARK__181
	params_205 = params_182
	ps_206 = ps_183
	idx_207 = idx_184
	bid_208 = bid_185
	insts_209 = insts_186
	f_210 = f_187
	changed_QMARK__211 = changed_QMARK__188
	term_212 = term_189
	goto b17
b17:
	;
	if vm.IsTruthy(changed_param_QMARK__204) {
		changed_term_QMARK__213 = changed_term_QMARK__202
		changed_inst_QMARK__214 = changed_inst_QMARK__203
		or__x_215 = changed_param_QMARK__204
		changed_param_QMARK__216 = changed_param_QMARK__204
		params_217 = params_205
		ps_218 = ps_206
		idx_219 = idx_207
		bid_220 = bid_208
		insts_221 = insts_209
		f_222 = f_210
		changed_QMARK__223 = changed_QMARK__211
		term_224 = term_212
		goto b18
	} else {
		changed_term_QMARK__225 = changed_term_QMARK__202
		changed_inst_QMARK__226 = changed_inst_QMARK__203
		or__x_227 = changed_param_QMARK__204
		changed_param_QMARK__228 = changed_param_QMARK__204
		params_229 = params_205
		ps_230 = ps_206
		idx_231 = idx_207
		bid_232 = bid_208
		insts_233 = insts_209
		f_234 = f_210
		changed_QMARK__235 = changed_QMARK__211
		term_236 = term_212
		goto b19
	}
b18:
	;
	v280 = or__x_215
	changed_term_QMARK__281 = changed_term_QMARK__213
	changed_inst_QMARK__282 = changed_inst_QMARK__214
	or__x_283 = or__x_215
	changed_param_QMARK__284 = changed_param_QMARK__216
	params_285 = params_217
	ps_286 = ps_218
	idx_287 = idx_219
	bid_288 = bid_220
	insts_289 = insts_221
	f_290 = f_222
	changed_QMARK__291 = changed_QMARK__223
	term_292 = term_224
	goto b20
b19:
	;
	if vm.IsTruthy(changed_inst_QMARK__226) {
		changed_term_QMARK__239 = changed_term_QMARK__225
		or__x_240 = changed_inst_QMARK__226
		changed_inst_QMARK__241 = changed_inst_QMARK__226
		changed_param_QMARK__242 = changed_param_QMARK__228
		params_243 = params_229
		ps_244 = ps_230
		idx_245 = idx_231
		bid_246 = bid_232
		insts_247 = insts_233
		f_248 = f_234
		changed_QMARK__249 = changed_QMARK__235
		term_250 = term_236
		goto b21
	} else {
		changed_term_QMARK__251 = changed_term_QMARK__225
		or__x_252 = changed_inst_QMARK__226
		changed_inst_QMARK__253 = changed_inst_QMARK__226
		changed_param_QMARK__254 = changed_param_QMARK__228
		params_255 = params_229
		ps_256 = ps_230
		idx_257 = idx_231
		bid_258 = bid_232
		insts_259 = insts_233
		f_260 = f_234
		changed_QMARK__261 = changed_QMARK__235
		term_262 = term_236
		goto b22
	}
b20:
	;
	return v280, nil
b21:
	;
	v266 = or__x_240
	changed_term_QMARK__267 = changed_term_QMARK__239
	or__x_268 = or__x_240
	changed_inst_QMARK__269 = changed_inst_QMARK__241
	changed_param_QMARK__270 = changed_param_QMARK__242
	params_271 = params_243
	ps_272 = ps_244
	idx_273 = idx_245
	bid_274 = bid_246
	insts_275 = insts_247
	f_276 = f_248
	changed_QMARK__277 = changed_QMARK__249
	term_278 = term_250
	goto b23
b22:
	;
	v266 = changed_term_QMARK__251
	changed_term_QMARK__267 = changed_term_QMARK__251
	or__x_268 = or__x_252
	changed_inst_QMARK__269 = changed_inst_QMARK__253
	changed_param_QMARK__270 = changed_param_QMARK__254
	params_271 = params_255
	ps_272 = ps_256
	idx_273 = idx_257
	bid_274 = bid_258
	insts_275 = insts_259
	f_276 = f_260
	changed_QMARK__277 = changed_QMARK__261
	term_278 = term_262
	goto b23
b23:
	;
	v280 = v266
	changed_term_QMARK__281 = changed_term_QMARK__267
	changed_inst_QMARK__282 = changed_inst_QMARK__269
	or__x_283 = or__x_227
	changed_param_QMARK__284 = changed_param_QMARK__270
	params_285 = params_271
	ps_286 = ps_272
	idx_287 = idx_273
	bid_288 = bid_274
	insts_289 = insts_275
	f_290 = f_276
	changed_QMARK__291 = changed_QMARK__277
	term_292 = term_278
	goto b20
}
func falsey_type_QMARK_(arg0 vm.Value) (vm.Value, error) {
	var t_2 vm.Value
	var v6 bool
	var t_3 vm.Value
	var t_4 vm.Value
	var v13 bool
	var v48 vm.Value
	var t_49 vm.Value
	var t_10 vm.Value
	var t_11 vm.Value
	var v20 vm.Value
	var v45 vm.Value
	var t_46 vm.Value
	var t_17 vm.Value
	var arg__28316_24 vm.Value
	var arg__28322_28 vm.Value
	var v29 vm.Value
	var t_18 vm.Value
	var v42 vm.Value
	var t_43 vm.Value
	var t_31 vm.Value
	var t_32 vm.Value
	var v39 vm.Value
	var t_40 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = t_2, v6, t_3, t_4, v13, v48, t_49, t_10, t_11, v20, v45, t_46, t_17, arg__28316_24, arg__28322_28, v29, t_18, v42, t_43, t_31, t_32, v39, t_40
	t_2, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "normalize-type").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v6 = t_2 == vm.Keyword("nil")
	if v6 {
		t_3 = t_2
		goto b1
	} else {
		t_4 = t_2
		goto b2
	}
b1:
	;
	v48 = vm.Boolean(true)
	t_49 = t_3
	goto b3
b2:
	;
	v13 = t_4 == vm.Keyword("false")
	if v13 {
		t_10 = t_4
		goto b4
	} else {
		t_11 = t_4
		goto b5
	}
b3:
	;
	return v48, nil
b4:
	;
	v45 = vm.Boolean(true)
	t_46 = t_10
	goto b6
b5:
	;
	v20, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.typeinfer", "union-type?").Deref(), []vm.Value{t_11})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v20) {
		t_17 = t_11
		goto b7
	} else {
		t_18 = t_11
		goto b8
	}
b6:
	;
	v48 = v45
	t_49 = t_46
	goto b3
b7:
	;
	arg__28316_24, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_17})
	if callErr != nil {
		return nil, callErr
	}
	arg__28322_28, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{t_17})
	if callErr != nil {
		return nil, callErr
	}
	v29, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "every?").Deref(), []vm.Value{rt.LookupVar("ir.passes.typeinfer", "falsey-type?").Deref(), arg__28322_28})
	if callErr != nil {
		return nil, callErr
	}
	v42 = v29
	t_43 = t_17
	goto b9
b8:
	;
	if vm.IsTruthy(vm.Keyword("else")) {
		t_31 = t_18
		goto b10
	} else {
		t_32 = t_18
		goto b11
	}
b9:
	;
	v45 = v42
	t_46 = t_43
	goto b6
b10:
	;
	v39 = vm.Boolean(false)
	t_40 = t_31
	goto b12
b11:
	;
	v39 = vm.NIL
	t_40 = t_32
	goto b12
b12:
	;
	v42 = v39
	t_43 = t_40
	goto b9
}
