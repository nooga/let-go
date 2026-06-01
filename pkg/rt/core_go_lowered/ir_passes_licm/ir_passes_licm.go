package ir_passes_licm

import (
	rt "github.com/nooga/let-go/pkg/rt"
	vm "github.com/nooga/let-go/pkg/vm"
)

func hoist_one_BANG_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var op_5 vm.Value
	var refs_7 vm.Value
	var aux_9 vm.Value
	var clone_11 vm.Value
	var from_block_13 vm.Value
	var v15 vm.Value
	var v17 vm.Value
	var callErr error
	_, _, _, _, _, _, _ = op_5, refs_7, aux_9, clone_11, from_block_13, v15, v17
	op_5, callErr = rt.InvokeValue(rt.LookupVar("ir", "op").Deref(), []vm.Value{arg1, arg0})
	if callErr != nil {
		return nil, callErr
	}
	refs_7, callErr = rt.InvokeValue(rt.LookupVar("ir", "refs").Deref(), []vm.Value{arg1, arg0})
	if callErr != nil {
		return nil, callErr
	}
	aux_9, callErr = rt.InvokeValue(rt.LookupVar("ir", "aux").Deref(), []vm.Value{arg1, arg0})
	if callErr != nil {
		return nil, callErr
	}
	clone_11, callErr = rt.InvokeValue(rt.LookupVar("ir", "add-inst").Deref(), []vm.Value{arg0, arg2, op_5, refs_7, aux_9})
	if callErr != nil {
		return nil, callErr
	}
	from_block_13, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-of").Deref(), []vm.Value{arg1, arg0})
	if callErr != nil {
		return nil, callErr
	}
	v15, callErr = rt.InvokeValue(rt.LookupVar("ir", "replace-all-uses!").Deref(), []vm.Value{arg0, arg1, clone_11})
	if callErr != nil {
		return nil, callErr
	}
	v17, callErr = rt.InvokeValue(rt.LookupVar("ir", "remove-inst!").Deref(), []vm.Value{arg0, from_block_13, arg1})
	if callErr != nil {
		return nil, callErr
	}
	return clone_11, nil
}
func back_edges(arg0 vm.Value) (vm.Value, error) {
	var for__a18461_5 vm.Value
	var __15 vm.Value
	var for__iter18460_17 vm.Value
	var arg__21006_19 vm.Value
	var arg__21011_21 vm.Value
	var arg__21012_22 vm.Value
	var for__a18461_27 vm.Value
	var __37 vm.Value
	var for__iter18460_39 vm.Value
	var arg__23558_41 vm.Value
	var arg__23563_43 vm.Value
	var arg__23564_44 vm.Value
	var v45 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _ = for__a18461_5, __15, for__iter18460_17, arg__21006_19, arg__21011_21, arg__21012_22, for__a18461_27, __37, for__iter18460_39, arg__23558_41, arg__23563_43, arg__23564_44, v45
	for__a18461_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "atom").Deref(), []vm.Value{vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	__15, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "reset!").Deref(), []vm.Value{for__a18461_5, rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v14 vm.Value
		var callErr error
		_ = v14
		v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "lazy-seq*").Deref(), []vm.Value{rt.BoxNativeFn(func() (vm.Value, error) {
			var tem__G__0_4 vm.Value
			var f_5 vm.Value
			var for__a18461_6 vm.Value
			var for__s_7 vm.Value
			var tem__G__0_8 vm.Value
			var header_15 vm.Value
			var for__a18459_19 vm.Value
			var __33 vm.Value
			var for__iter18458_35 vm.Value
			var arg__20663_37 vm.Value
			var arg__20670_39 vm.Value
			var arg__20671_40 vm.Value
			var arg__20675_42 vm.Value
			var head__20679_44 vm.Value
			var arg__20683_46 vm.Value
			var arg__20684_47 vm.Value
			var for__a18459_52 vm.Value
			var __66 vm.Value
			var for__iter18458_68 vm.Value
			var arg__20976_70 vm.Value
			var arg__20983_72 vm.Value
			var arg__20984_73 vm.Value
			var arg__20988_75 vm.Value
			var head__20992_77 vm.Value
			var arg__20996_79 vm.Value
			var arg__20997_80 vm.Value
			var v81 vm.Value
			var f_9 vm.Value
			var for__a18461_10 vm.Value
			var for__s_11 vm.Value
			var tem__G__0_12 vm.Value
			var v85 vm.Value
			var f_86 vm.Value
			var for__a18461_87 vm.Value
			var for__s_88 vm.Value
			var tem__G__0_89 vm.Value
			var callErr error
			_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = tem__G__0_4, f_5, for__a18461_6, for__s_7, tem__G__0_8, header_15, for__a18459_19, __33, for__iter18458_35, arg__20663_37, arg__20670_39, arg__20671_40, arg__20675_42, head__20679_44, arg__20683_46, arg__20684_47, for__a18459_52, __66, for__iter18458_68, arg__20976_70, arg__20983_72, arg__20984_73, arg__20988_75, head__20992_77, arg__20996_79, arg__20997_80, v81, f_9, for__a18461_10, for__s_11, tem__G__0_12, v85, f_86, for__a18461_87, for__s_88, tem__G__0_89
			tem__G__0_4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{arg0})
			if callErr != nil {
				return nil, callErr
			}
			if vm.IsTruthy(tem__G__0_4) {
				f_5 = arg0
				for__a18461_6 = for__a18461_5
				for__s_7 = arg0
				tem__G__0_8 = tem__G__0_4
				goto b1
			} else {
				f_9 = arg0
				for__a18461_10 = for__a18461_5
				for__s_11 = arg0
				tem__G__0_12 = tem__G__0_4
				goto b2
			}
		b1:
			;
			header_15, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{tem__G__0_8})
			if callErr != nil {
				return nil, callErr
			}
			for__a18459_19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "atom").Deref(), []vm.Value{vm.NIL})
			if callErr != nil {
				return nil, callErr
			}
			__33, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "reset!").Deref(), []vm.Value{for__a18459_19, rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
				var v18 vm.Value
				var callErr error
				_ = v18
				v18, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "lazy-seq*").Deref(), []vm.Value{rt.BoxNativeFn(func() (vm.Value, error) {
					var tem__G__0_5 vm.Value
					var f_6 vm.Value
					var for__a18459_7 vm.Value
					var for__s_8 vm.Value
					var header_9 vm.Value
					var tem__G__0_10 vm.Value
					var pred_18 vm.Value
					var v34 vm.Value
					var f_11 vm.Value
					var for__a18459_12 vm.Value
					var for__s_13 vm.Value
					var header_14 vm.Value
					var tem__G__0_15 vm.Value
					var v109 vm.Value
					var f_110 vm.Value
					var for__a18459_111 vm.Value
					var for__s_112 vm.Value
					var header_113 vm.Value
					var tem__G__0_114 vm.Value
					var f_19 vm.Value
					var for__a18459_20 vm.Value
					var for__s_21 vm.Value
					var header_22 vm.Value
					var tem__G__0_23 vm.Value
					var for__xs_24 vm.Value
					var pred_25 vm.Value
					var arg__20602_37 vm.Value
					var arg__20607_40 vm.Value
					var v41 vm.Value
					var f_26 vm.Value
					var for__a18459_27 vm.Value
					var for__s_28 vm.Value
					var header_29 vm.Value
					var tem__G__0_30 vm.Value
					var for__xs_31 vm.Value
					var pred_32 vm.Value
					var arg__20608_45 vm.Value
					var f_46 vm.Value
					var for__a18459_47 vm.Value
					var for__s_48 vm.Value
					var header_49 vm.Value
					var tem__G__0_50 vm.Value
					var for__xs_51 vm.Value
					var pred_52 vm.Value
					var arg__20612_54 vm.Value
					var head__20616_56 vm.Value
					var arg__20620_58 vm.Value
					var arg__20621_59 vm.Value
					var v78 vm.Value
					var f_61 vm.Value
					var for__a18459_62 vm.Value
					var for__s_63 vm.Value
					var header_64 vm.Value
					var tem__G__0_65 vm.Value
					var for__xs_66 vm.Value
					var pred_67 vm.Value
					var head__20622_68 vm.Value
					var arg__20633_81 vm.Value
					var arg__20638_84 vm.Value
					var v85 vm.Value
					var f_69 vm.Value
					var for__a18459_70 vm.Value
					var for__s_71 vm.Value
					var header_72 vm.Value
					var tem__G__0_73 vm.Value
					var for__xs_74 vm.Value
					var pred_75 vm.Value
					var head__20622_76 vm.Value
					var arg__20639_89 vm.Value
					var f_90 vm.Value
					var for__a18459_91 vm.Value
					var for__s_92 vm.Value
					var header_93 vm.Value
					var tem__G__0_94 vm.Value
					var for__xs_95 vm.Value
					var pred_96 vm.Value
					var head__20622_97 vm.Value
					var arg__20643_99 vm.Value
					var head__20647_101 vm.Value
					var arg__20651_103 vm.Value
					var arg__20652_104 vm.Value
					var v105 vm.Value
					var callErr error
					_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = tem__G__0_5, f_6, for__a18459_7, for__s_8, header_9, tem__G__0_10, pred_18, v34, f_11, for__a18459_12, for__s_13, header_14, tem__G__0_15, v109, f_110, for__a18459_111, for__s_112, header_113, tem__G__0_114, f_19, for__a18459_20, for__s_21, header_22, tem__G__0_23, for__xs_24, pred_25, arg__20602_37, arg__20607_40, v41, f_26, for__a18459_27, for__s_28, header_29, tem__G__0_30, for__xs_31, pred_32, arg__20608_45, f_46, for__a18459_47, for__s_48, header_49, tem__G__0_50, for__xs_51, pred_52, arg__20612_54, head__20616_56, arg__20620_58, arg__20621_59, v78, f_61, for__a18459_62, for__s_63, header_64, tem__G__0_65, for__xs_66, pred_67, head__20622_68, arg__20633_81, arg__20638_84, v85, f_69, for__a18459_70, for__s_71, header_72, tem__G__0_73, for__xs_74, pred_75, head__20622_76, arg__20639_89, f_90, for__a18459_91, for__s_92, header_93, tem__G__0_94, for__xs_95, pred_96, head__20622_97, arg__20643_99, head__20647_101, arg__20651_103, arg__20652_104, v105
					tem__G__0_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{for__s_7})
					if callErr != nil {
						return nil, callErr
					}
					if vm.IsTruthy(tem__G__0_5) {
						f_6 = f_5
						for__a18459_7 = for__a18459_19
						for__s_8 = for__s_7
						header_9 = header_15
						tem__G__0_10 = tem__G__0_5
						goto b1
					} else {
						f_11 = f_5
						for__a18459_12 = for__a18459_19
						for__s_13 = for__s_7
						header_14 = header_15
						tem__G__0_15 = tem__G__0_5
						goto b2
					}
				b1:
					;
					pred_18, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{tem__G__0_10})
					if callErr != nil {
						return nil, callErr
					}
					v34, callErr = rt.InvokeValue(rt.LookupVar("ir.dominance", "dominates?").Deref(), []vm.Value{f_6, header_9, pred_18})
					if callErr != nil {
						return nil, callErr
					}
					if vm.IsTruthy(v34) {
						f_19 = f_6
						for__a18459_20 = for__a18459_7
						for__s_21 = for__s_8
						header_22 = header_9
						tem__G__0_23 = tem__G__0_10
						for__xs_24 = tem__G__0_10
						pred_25 = pred_18
						goto b4
					} else {
						f_26 = f_6
						for__a18459_27 = for__a18459_7
						for__s_28 = for__s_8
						header_29 = header_9
						tem__G__0_30 = tem__G__0_10
						for__xs_31 = tem__G__0_10
						pred_32 = pred_18
						goto b5
					}
				b2:
					;
					v109 = vm.NIL
					f_110 = f_11
					for__a18459_111 = for__a18459_12
					for__s_112 = for__s_13
					header_113 = header_14
					tem__G__0_114 = tem__G__0_15
					goto b3
				b3:
					;
					return v109, nil
				b4:
					;
					arg__20602_37, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{pred_25, header_22})
					if callErr != nil {
						return nil, callErr
					}
					arg__20607_40, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{pred_25, header_22})
					if callErr != nil {
						return nil, callErr
					}
					v41, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "list").Deref(), []vm.Value{arg__20607_40})
					if callErr != nil {
						return nil, callErr
					}
					arg__20608_45 = v41
					f_46 = f_19
					for__a18459_47 = for__a18459_20
					for__s_48 = for__s_21
					header_49 = header_22
					tem__G__0_50 = tem__G__0_23
					for__xs_51 = for__xs_24
					pred_52 = pred_25
					goto b6
				b5:
					;
					arg__20608_45 = vm.NIL
					f_46 = f_26
					for__a18459_47 = for__a18459_27
					for__s_48 = for__s_28
					header_49 = header_29
					tem__G__0_50 = tem__G__0_30
					for__xs_51 = for__xs_31
					pred_52 = pred_32
					goto b6
				b6:
					;
					arg__20612_54, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{for__xs_51})
					if callErr != nil {
						return nil, callErr
					}
					head__20616_56, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{for__a18459_47})
					if callErr != nil {
						return nil, callErr
					}
					arg__20620_58, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{for__xs_51})
					if callErr != nil {
						return nil, callErr
					}
					arg__20621_59, callErr = rt.InvokeValue(head__20616_56, []vm.Value{arg__20620_58})
					if callErr != nil {
						return nil, callErr
					}
					v78, callErr = rt.InvokeValue(rt.LookupVar("ir.dominance", "dominates?").Deref(), []vm.Value{f_46, header_49, pred_52})
					if callErr != nil {
						return nil, callErr
					}
					if vm.IsTruthy(v78) {
						f_61 = f_46
						for__a18459_62 = for__a18459_47
						for__s_63 = for__s_48
						header_64 = header_49
						tem__G__0_65 = tem__G__0_50
						for__xs_66 = for__xs_51
						pred_67 = pred_52
						head__20622_68 = rt.LookupVar("clojure.core", "concat-list").Deref()
						goto b7
					} else {
						f_69 = f_46
						for__a18459_70 = for__a18459_47
						for__s_71 = for__s_48
						header_72 = header_49
						tem__G__0_73 = tem__G__0_50
						for__xs_74 = for__xs_51
						pred_75 = pred_52
						head__20622_76 = rt.LookupVar("clojure.core", "concat-list").Deref()
						goto b8
					}
				b7:
					;
					arg__20633_81, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{pred_67, header_64})
					if callErr != nil {
						return nil, callErr
					}
					arg__20638_84, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{pred_67, header_64})
					if callErr != nil {
						return nil, callErr
					}
					v85, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "list").Deref(), []vm.Value{arg__20638_84})
					if callErr != nil {
						return nil, callErr
					}
					arg__20639_89 = v85
					f_90 = f_61
					for__a18459_91 = for__a18459_62
					for__s_92 = for__s_63
					header_93 = header_64
					tem__G__0_94 = tem__G__0_65
					for__xs_95 = for__xs_66
					pred_96 = pred_67
					head__20622_97 = head__20622_68
					goto b9
				b8:
					;
					arg__20639_89 = vm.NIL
					f_90 = f_69
					for__a18459_91 = for__a18459_70
					for__s_92 = for__s_71
					header_93 = header_72
					tem__G__0_94 = tem__G__0_73
					for__xs_95 = for__xs_74
					pred_96 = pred_75
					head__20622_97 = head__20622_76
					goto b9
				b9:
					;
					arg__20643_99, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{for__xs_95})
					if callErr != nil {
						return nil, callErr
					}
					head__20647_101, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{for__a18459_91})
					if callErr != nil {
						return nil, callErr
					}
					arg__20651_103, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{for__xs_95})
					if callErr != nil {
						return nil, callErr
					}
					arg__20652_104, callErr = rt.InvokeValue(head__20647_101, []vm.Value{arg__20651_103})
					if callErr != nil {
						return nil, callErr
					}
					v105, callErr = rt.InvokeValue(head__20622_97, []vm.Value{arg__20639_89, arg__20652_104})
					if callErr != nil {
						return nil, callErr
					}
					v109 = v105
					f_110 = f_90
					for__a18459_111 = for__a18459_91
					for__s_112 = for__s_92
					header_113 = header_93
					tem__G__0_114 = tem__G__0_94
					goto b3
				})})
				if callErr != nil {
					return nil, callErr
				}
				return v18, nil
			})})
			if callErr != nil {
				return nil, callErr
			}
			for__iter18458_35, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{for__a18459_19})
			if callErr != nil {
				return nil, callErr
			}
			arg__20663_37, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-preds").Deref(), []vm.Value{header_15, f_5})
			if callErr != nil {
				return nil, callErr
			}
			arg__20670_39, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-preds").Deref(), []vm.Value{header_15, f_5})
			if callErr != nil {
				return nil, callErr
			}
			arg__20671_40, callErr = rt.InvokeValue(for__iter18458_35, []vm.Value{arg__20670_39})
			if callErr != nil {
				return nil, callErr
			}
			arg__20675_42, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{tem__G__0_8})
			if callErr != nil {
				return nil, callErr
			}
			head__20679_44, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{for__a18461_6})
			if callErr != nil {
				return nil, callErr
			}
			arg__20683_46, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{tem__G__0_8})
			if callErr != nil {
				return nil, callErr
			}
			arg__20684_47, callErr = rt.InvokeValue(head__20679_44, []vm.Value{arg__20683_46})
			if callErr != nil {
				return nil, callErr
			}
			for__a18459_52, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "atom").Deref(), []vm.Value{vm.NIL})
			if callErr != nil {
				return nil, callErr
			}
			__66, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "reset!").Deref(), []vm.Value{for__a18459_52, rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
				var v18 vm.Value
				var callErr error
				_ = v18
				v18, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "lazy-seq*").Deref(), []vm.Value{rt.BoxNativeFn(func() (vm.Value, error) {
					var tem__G__0_5 vm.Value
					var f_6 vm.Value
					var for__a18459_7 vm.Value
					var for__s_8 vm.Value
					var header_9 vm.Value
					var tem__G__0_10 vm.Value
					var pred_18 vm.Value
					var v34 vm.Value
					var f_11 vm.Value
					var for__a18459_12 vm.Value
					var for__s_13 vm.Value
					var header_14 vm.Value
					var tem__G__0_15 vm.Value
					var v109 vm.Value
					var f_110 vm.Value
					var for__a18459_111 vm.Value
					var for__s_112 vm.Value
					var header_113 vm.Value
					var tem__G__0_114 vm.Value
					var f_19 vm.Value
					var for__a18459_20 vm.Value
					var for__s_21 vm.Value
					var header_22 vm.Value
					var tem__G__0_23 vm.Value
					var for__xs_24 vm.Value
					var pred_25 vm.Value
					var arg__20915_37 vm.Value
					var arg__20920_40 vm.Value
					var v41 vm.Value
					var f_26 vm.Value
					var for__a18459_27 vm.Value
					var for__s_28 vm.Value
					var header_29 vm.Value
					var tem__G__0_30 vm.Value
					var for__xs_31 vm.Value
					var pred_32 vm.Value
					var arg__20921_45 vm.Value
					var f_46 vm.Value
					var for__a18459_47 vm.Value
					var for__s_48 vm.Value
					var header_49 vm.Value
					var tem__G__0_50 vm.Value
					var for__xs_51 vm.Value
					var pred_52 vm.Value
					var arg__20925_54 vm.Value
					var head__20929_56 vm.Value
					var arg__20933_58 vm.Value
					var arg__20934_59 vm.Value
					var v78 vm.Value
					var f_61 vm.Value
					var for__a18459_62 vm.Value
					var for__s_63 vm.Value
					var header_64 vm.Value
					var tem__G__0_65 vm.Value
					var for__xs_66 vm.Value
					var pred_67 vm.Value
					var head__20935_68 vm.Value
					var arg__20946_81 vm.Value
					var arg__20951_84 vm.Value
					var v85 vm.Value
					var f_69 vm.Value
					var for__a18459_70 vm.Value
					var for__s_71 vm.Value
					var header_72 vm.Value
					var tem__G__0_73 vm.Value
					var for__xs_74 vm.Value
					var pred_75 vm.Value
					var head__20935_76 vm.Value
					var arg__20952_89 vm.Value
					var f_90 vm.Value
					var for__a18459_91 vm.Value
					var for__s_92 vm.Value
					var header_93 vm.Value
					var tem__G__0_94 vm.Value
					var for__xs_95 vm.Value
					var pred_96 vm.Value
					var head__20935_97 vm.Value
					var arg__20956_99 vm.Value
					var head__20960_101 vm.Value
					var arg__20964_103 vm.Value
					var arg__20965_104 vm.Value
					var v105 vm.Value
					var callErr error
					_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = tem__G__0_5, f_6, for__a18459_7, for__s_8, header_9, tem__G__0_10, pred_18, v34, f_11, for__a18459_12, for__s_13, header_14, tem__G__0_15, v109, f_110, for__a18459_111, for__s_112, header_113, tem__G__0_114, f_19, for__a18459_20, for__s_21, header_22, tem__G__0_23, for__xs_24, pred_25, arg__20915_37, arg__20920_40, v41, f_26, for__a18459_27, for__s_28, header_29, tem__G__0_30, for__xs_31, pred_32, arg__20921_45, f_46, for__a18459_47, for__s_48, header_49, tem__G__0_50, for__xs_51, pred_52, arg__20925_54, head__20929_56, arg__20933_58, arg__20934_59, v78, f_61, for__a18459_62, for__s_63, header_64, tem__G__0_65, for__xs_66, pred_67, head__20935_68, arg__20946_81, arg__20951_84, v85, f_69, for__a18459_70, for__s_71, header_72, tem__G__0_73, for__xs_74, pred_75, head__20935_76, arg__20952_89, f_90, for__a18459_91, for__s_92, header_93, tem__G__0_94, for__xs_95, pred_96, head__20935_97, arg__20956_99, head__20960_101, arg__20964_103, arg__20965_104, v105
					tem__G__0_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{for__s_7})
					if callErr != nil {
						return nil, callErr
					}
					if vm.IsTruthy(tem__G__0_5) {
						f_6 = f_5
						for__a18459_7 = for__a18459_52
						for__s_8 = for__s_7
						header_9 = header_15
						tem__G__0_10 = tem__G__0_5
						goto b1
					} else {
						f_11 = f_5
						for__a18459_12 = for__a18459_52
						for__s_13 = for__s_7
						header_14 = header_15
						tem__G__0_15 = tem__G__0_5
						goto b2
					}
				b1:
					;
					pred_18, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{tem__G__0_10})
					if callErr != nil {
						return nil, callErr
					}
					v34, callErr = rt.InvokeValue(rt.LookupVar("ir.dominance", "dominates?").Deref(), []vm.Value{f_6, header_9, pred_18})
					if callErr != nil {
						return nil, callErr
					}
					if vm.IsTruthy(v34) {
						f_19 = f_6
						for__a18459_20 = for__a18459_7
						for__s_21 = for__s_8
						header_22 = header_9
						tem__G__0_23 = tem__G__0_10
						for__xs_24 = tem__G__0_10
						pred_25 = pred_18
						goto b4
					} else {
						f_26 = f_6
						for__a18459_27 = for__a18459_7
						for__s_28 = for__s_8
						header_29 = header_9
						tem__G__0_30 = tem__G__0_10
						for__xs_31 = tem__G__0_10
						pred_32 = pred_18
						goto b5
					}
				b2:
					;
					v109 = vm.NIL
					f_110 = f_11
					for__a18459_111 = for__a18459_12
					for__s_112 = for__s_13
					header_113 = header_14
					tem__G__0_114 = tem__G__0_15
					goto b3
				b3:
					;
					return v109, nil
				b4:
					;
					arg__20915_37, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{pred_25, header_22})
					if callErr != nil {
						return nil, callErr
					}
					arg__20920_40, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{pred_25, header_22})
					if callErr != nil {
						return nil, callErr
					}
					v41, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "list").Deref(), []vm.Value{arg__20920_40})
					if callErr != nil {
						return nil, callErr
					}
					arg__20921_45 = v41
					f_46 = f_19
					for__a18459_47 = for__a18459_20
					for__s_48 = for__s_21
					header_49 = header_22
					tem__G__0_50 = tem__G__0_23
					for__xs_51 = for__xs_24
					pred_52 = pred_25
					goto b6
				b5:
					;
					arg__20921_45 = vm.NIL
					f_46 = f_26
					for__a18459_47 = for__a18459_27
					for__s_48 = for__s_28
					header_49 = header_29
					tem__G__0_50 = tem__G__0_30
					for__xs_51 = for__xs_31
					pred_52 = pred_32
					goto b6
				b6:
					;
					arg__20925_54, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{for__xs_51})
					if callErr != nil {
						return nil, callErr
					}
					head__20929_56, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{for__a18459_47})
					if callErr != nil {
						return nil, callErr
					}
					arg__20933_58, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{for__xs_51})
					if callErr != nil {
						return nil, callErr
					}
					arg__20934_59, callErr = rt.InvokeValue(head__20929_56, []vm.Value{arg__20933_58})
					if callErr != nil {
						return nil, callErr
					}
					v78, callErr = rt.InvokeValue(rt.LookupVar("ir.dominance", "dominates?").Deref(), []vm.Value{f_46, header_49, pred_52})
					if callErr != nil {
						return nil, callErr
					}
					if vm.IsTruthy(v78) {
						f_61 = f_46
						for__a18459_62 = for__a18459_47
						for__s_63 = for__s_48
						header_64 = header_49
						tem__G__0_65 = tem__G__0_50
						for__xs_66 = for__xs_51
						pred_67 = pred_52
						head__20935_68 = rt.LookupVar("clojure.core", "concat-list").Deref()
						goto b7
					} else {
						f_69 = f_46
						for__a18459_70 = for__a18459_47
						for__s_71 = for__s_48
						header_72 = header_49
						tem__G__0_73 = tem__G__0_50
						for__xs_74 = for__xs_51
						pred_75 = pred_52
						head__20935_76 = rt.LookupVar("clojure.core", "concat-list").Deref()
						goto b8
					}
				b7:
					;
					arg__20946_81, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{pred_67, header_64})
					if callErr != nil {
						return nil, callErr
					}
					arg__20951_84, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{pred_67, header_64})
					if callErr != nil {
						return nil, callErr
					}
					v85, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "list").Deref(), []vm.Value{arg__20951_84})
					if callErr != nil {
						return nil, callErr
					}
					arg__20952_89 = v85
					f_90 = f_61
					for__a18459_91 = for__a18459_62
					for__s_92 = for__s_63
					header_93 = header_64
					tem__G__0_94 = tem__G__0_65
					for__xs_95 = for__xs_66
					pred_96 = pred_67
					head__20935_97 = head__20935_68
					goto b9
				b8:
					;
					arg__20952_89 = vm.NIL
					f_90 = f_69
					for__a18459_91 = for__a18459_70
					for__s_92 = for__s_71
					header_93 = header_72
					tem__G__0_94 = tem__G__0_73
					for__xs_95 = for__xs_74
					pred_96 = pred_75
					head__20935_97 = head__20935_76
					goto b9
				b9:
					;
					arg__20956_99, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{for__xs_95})
					if callErr != nil {
						return nil, callErr
					}
					head__20960_101, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{for__a18459_91})
					if callErr != nil {
						return nil, callErr
					}
					arg__20964_103, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{for__xs_95})
					if callErr != nil {
						return nil, callErr
					}
					arg__20965_104, callErr = rt.InvokeValue(head__20960_101, []vm.Value{arg__20964_103})
					if callErr != nil {
						return nil, callErr
					}
					v105, callErr = rt.InvokeValue(head__20935_97, []vm.Value{arg__20952_89, arg__20965_104})
					if callErr != nil {
						return nil, callErr
					}
					v109 = v105
					f_110 = f_90
					for__a18459_111 = for__a18459_91
					for__s_112 = for__s_92
					header_113 = header_93
					tem__G__0_114 = tem__G__0_94
					goto b3
				})})
				if callErr != nil {
					return nil, callErr
				}
				return v18, nil
			})})
			if callErr != nil {
				return nil, callErr
			}
			for__iter18458_68, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{for__a18459_52})
			if callErr != nil {
				return nil, callErr
			}
			arg__20976_70, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-preds").Deref(), []vm.Value{header_15, f_5})
			if callErr != nil {
				return nil, callErr
			}
			arg__20983_72, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-preds").Deref(), []vm.Value{header_15, f_5})
			if callErr != nil {
				return nil, callErr
			}
			arg__20984_73, callErr = rt.InvokeValue(for__iter18458_68, []vm.Value{arg__20983_72})
			if callErr != nil {
				return nil, callErr
			}
			arg__20988_75, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{tem__G__0_8})
			if callErr != nil {
				return nil, callErr
			}
			head__20992_77, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{for__a18461_6})
			if callErr != nil {
				return nil, callErr
			}
			arg__20996_79, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{tem__G__0_8})
			if callErr != nil {
				return nil, callErr
			}
			arg__20997_80, callErr = rt.InvokeValue(head__20992_77, []vm.Value{arg__20996_79})
			if callErr != nil {
				return nil, callErr
			}
			v81, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "concat-list").Deref(), []vm.Value{arg__20984_73, arg__20997_80})
			if callErr != nil {
				return nil, callErr
			}
			v85 = v81
			f_86 = f_5
			for__a18461_87 = for__a18461_6
			for__s_88 = for__s_7
			tem__G__0_89 = tem__G__0_8
			goto b3
		b2:
			;
			v85 = vm.NIL
			f_86 = f_9
			for__a18461_87 = for__a18461_10
			for__s_88 = for__s_11
			tem__G__0_89 = tem__G__0_12
			goto b3
		b3:
			;
			return v85, nil
		})})
		if callErr != nil {
			return nil, callErr
		}
		return v14, nil
	})})
	if callErr != nil {
		return nil, callErr
	}
	for__iter18460_17, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{for__a18461_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__21006_19, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__21011_21, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__21012_22, callErr = rt.InvokeValue(for__iter18460_17, []vm.Value{arg__21011_21})
	if callErr != nil {
		return nil, callErr
	}
	for__a18461_27, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "atom").Deref(), []vm.Value{vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	__37, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "reset!").Deref(), []vm.Value{for__a18461_27, rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v14 vm.Value
		var callErr error
		_ = v14
		v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "lazy-seq*").Deref(), []vm.Value{rt.BoxNativeFn(func() (vm.Value, error) {
			var tem__G__0_4 vm.Value
			var f_5 vm.Value
			var for__a18461_6 vm.Value
			var for__s_7 vm.Value
			var tem__G__0_8 vm.Value
			var header_15 vm.Value
			var for__a18459_19 vm.Value
			var __33 vm.Value
			var for__iter18458_35 vm.Value
			var arg__23215_37 vm.Value
			var arg__23222_39 vm.Value
			var arg__23223_40 vm.Value
			var arg__23227_42 vm.Value
			var head__23231_44 vm.Value
			var arg__23235_46 vm.Value
			var arg__23236_47 vm.Value
			var for__a18459_52 vm.Value
			var __66 vm.Value
			var for__iter18458_68 vm.Value
			var arg__23528_70 vm.Value
			var arg__23535_72 vm.Value
			var arg__23536_73 vm.Value
			var arg__23540_75 vm.Value
			var head__23544_77 vm.Value
			var arg__23548_79 vm.Value
			var arg__23549_80 vm.Value
			var v81 vm.Value
			var f_9 vm.Value
			var for__a18461_10 vm.Value
			var for__s_11 vm.Value
			var tem__G__0_12 vm.Value
			var v85 vm.Value
			var f_86 vm.Value
			var for__a18461_87 vm.Value
			var for__s_88 vm.Value
			var tem__G__0_89 vm.Value
			var callErr error
			_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = tem__G__0_4, f_5, for__a18461_6, for__s_7, tem__G__0_8, header_15, for__a18459_19, __33, for__iter18458_35, arg__23215_37, arg__23222_39, arg__23223_40, arg__23227_42, head__23231_44, arg__23235_46, arg__23236_47, for__a18459_52, __66, for__iter18458_68, arg__23528_70, arg__23535_72, arg__23536_73, arg__23540_75, head__23544_77, arg__23548_79, arg__23549_80, v81, f_9, for__a18461_10, for__s_11, tem__G__0_12, v85, f_86, for__a18461_87, for__s_88, tem__G__0_89
			tem__G__0_4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{arg0})
			if callErr != nil {
				return nil, callErr
			}
			if vm.IsTruthy(tem__G__0_4) {
				f_5 = arg0
				for__a18461_6 = for__a18461_27
				for__s_7 = arg0
				tem__G__0_8 = tem__G__0_4
				goto b1
			} else {
				f_9 = arg0
				for__a18461_10 = for__a18461_27
				for__s_11 = arg0
				tem__G__0_12 = tem__G__0_4
				goto b2
			}
		b1:
			;
			header_15, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{tem__G__0_8})
			if callErr != nil {
				return nil, callErr
			}
			for__a18459_19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "atom").Deref(), []vm.Value{vm.NIL})
			if callErr != nil {
				return nil, callErr
			}
			__33, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "reset!").Deref(), []vm.Value{for__a18459_19, rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
				var v18 vm.Value
				var callErr error
				_ = v18
				v18, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "lazy-seq*").Deref(), []vm.Value{rt.BoxNativeFn(func() (vm.Value, error) {
					var tem__G__0_5 vm.Value
					var f_6 vm.Value
					var for__a18459_7 vm.Value
					var for__s_8 vm.Value
					var header_9 vm.Value
					var tem__G__0_10 vm.Value
					var pred_18 vm.Value
					var v34 vm.Value
					var f_11 vm.Value
					var for__a18459_12 vm.Value
					var for__s_13 vm.Value
					var header_14 vm.Value
					var tem__G__0_15 vm.Value
					var v109 vm.Value
					var f_110 vm.Value
					var for__a18459_111 vm.Value
					var for__s_112 vm.Value
					var header_113 vm.Value
					var tem__G__0_114 vm.Value
					var f_19 vm.Value
					var for__a18459_20 vm.Value
					var for__s_21 vm.Value
					var header_22 vm.Value
					var tem__G__0_23 vm.Value
					var for__xs_24 vm.Value
					var pred_25 vm.Value
					var arg__23154_37 vm.Value
					var arg__23159_40 vm.Value
					var v41 vm.Value
					var f_26 vm.Value
					var for__a18459_27 vm.Value
					var for__s_28 vm.Value
					var header_29 vm.Value
					var tem__G__0_30 vm.Value
					var for__xs_31 vm.Value
					var pred_32 vm.Value
					var arg__23160_45 vm.Value
					var f_46 vm.Value
					var for__a18459_47 vm.Value
					var for__s_48 vm.Value
					var header_49 vm.Value
					var tem__G__0_50 vm.Value
					var for__xs_51 vm.Value
					var pred_52 vm.Value
					var arg__23164_54 vm.Value
					var head__23168_56 vm.Value
					var arg__23172_58 vm.Value
					var arg__23173_59 vm.Value
					var v78 vm.Value
					var f_61 vm.Value
					var for__a18459_62 vm.Value
					var for__s_63 vm.Value
					var header_64 vm.Value
					var tem__G__0_65 vm.Value
					var for__xs_66 vm.Value
					var pred_67 vm.Value
					var head__23174_68 vm.Value
					var arg__23185_81 vm.Value
					var arg__23190_84 vm.Value
					var v85 vm.Value
					var f_69 vm.Value
					var for__a18459_70 vm.Value
					var for__s_71 vm.Value
					var header_72 vm.Value
					var tem__G__0_73 vm.Value
					var for__xs_74 vm.Value
					var pred_75 vm.Value
					var head__23174_76 vm.Value
					var arg__23191_89 vm.Value
					var f_90 vm.Value
					var for__a18459_91 vm.Value
					var for__s_92 vm.Value
					var header_93 vm.Value
					var tem__G__0_94 vm.Value
					var for__xs_95 vm.Value
					var pred_96 vm.Value
					var head__23174_97 vm.Value
					var arg__23195_99 vm.Value
					var head__23199_101 vm.Value
					var arg__23203_103 vm.Value
					var arg__23204_104 vm.Value
					var v105 vm.Value
					var callErr error
					_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = tem__G__0_5, f_6, for__a18459_7, for__s_8, header_9, tem__G__0_10, pred_18, v34, f_11, for__a18459_12, for__s_13, header_14, tem__G__0_15, v109, f_110, for__a18459_111, for__s_112, header_113, tem__G__0_114, f_19, for__a18459_20, for__s_21, header_22, tem__G__0_23, for__xs_24, pred_25, arg__23154_37, arg__23159_40, v41, f_26, for__a18459_27, for__s_28, header_29, tem__G__0_30, for__xs_31, pred_32, arg__23160_45, f_46, for__a18459_47, for__s_48, header_49, tem__G__0_50, for__xs_51, pred_52, arg__23164_54, head__23168_56, arg__23172_58, arg__23173_59, v78, f_61, for__a18459_62, for__s_63, header_64, tem__G__0_65, for__xs_66, pred_67, head__23174_68, arg__23185_81, arg__23190_84, v85, f_69, for__a18459_70, for__s_71, header_72, tem__G__0_73, for__xs_74, pred_75, head__23174_76, arg__23191_89, f_90, for__a18459_91, for__s_92, header_93, tem__G__0_94, for__xs_95, pred_96, head__23174_97, arg__23195_99, head__23199_101, arg__23203_103, arg__23204_104, v105
					tem__G__0_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{for__s_7})
					if callErr != nil {
						return nil, callErr
					}
					if vm.IsTruthy(tem__G__0_5) {
						f_6 = f_5
						for__a18459_7 = for__a18459_19
						for__s_8 = for__s_7
						header_9 = header_15
						tem__G__0_10 = tem__G__0_5
						goto b1
					} else {
						f_11 = f_5
						for__a18459_12 = for__a18459_19
						for__s_13 = for__s_7
						header_14 = header_15
						tem__G__0_15 = tem__G__0_5
						goto b2
					}
				b1:
					;
					pred_18, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{tem__G__0_10})
					if callErr != nil {
						return nil, callErr
					}
					v34, callErr = rt.InvokeValue(rt.LookupVar("ir.dominance", "dominates?").Deref(), []vm.Value{f_6, header_9, pred_18})
					if callErr != nil {
						return nil, callErr
					}
					if vm.IsTruthy(v34) {
						f_19 = f_6
						for__a18459_20 = for__a18459_7
						for__s_21 = for__s_8
						header_22 = header_9
						tem__G__0_23 = tem__G__0_10
						for__xs_24 = tem__G__0_10
						pred_25 = pred_18
						goto b4
					} else {
						f_26 = f_6
						for__a18459_27 = for__a18459_7
						for__s_28 = for__s_8
						header_29 = header_9
						tem__G__0_30 = tem__G__0_10
						for__xs_31 = tem__G__0_10
						pred_32 = pred_18
						goto b5
					}
				b2:
					;
					v109 = vm.NIL
					f_110 = f_11
					for__a18459_111 = for__a18459_12
					for__s_112 = for__s_13
					header_113 = header_14
					tem__G__0_114 = tem__G__0_15
					goto b3
				b3:
					;
					return v109, nil
				b4:
					;
					arg__23154_37, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{pred_25, header_22})
					if callErr != nil {
						return nil, callErr
					}
					arg__23159_40, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{pred_25, header_22})
					if callErr != nil {
						return nil, callErr
					}
					v41, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "list").Deref(), []vm.Value{arg__23159_40})
					if callErr != nil {
						return nil, callErr
					}
					arg__23160_45 = v41
					f_46 = f_19
					for__a18459_47 = for__a18459_20
					for__s_48 = for__s_21
					header_49 = header_22
					tem__G__0_50 = tem__G__0_23
					for__xs_51 = for__xs_24
					pred_52 = pred_25
					goto b6
				b5:
					;
					arg__23160_45 = vm.NIL
					f_46 = f_26
					for__a18459_47 = for__a18459_27
					for__s_48 = for__s_28
					header_49 = header_29
					tem__G__0_50 = tem__G__0_30
					for__xs_51 = for__xs_31
					pred_52 = pred_32
					goto b6
				b6:
					;
					arg__23164_54, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{for__xs_51})
					if callErr != nil {
						return nil, callErr
					}
					head__23168_56, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{for__a18459_47})
					if callErr != nil {
						return nil, callErr
					}
					arg__23172_58, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{for__xs_51})
					if callErr != nil {
						return nil, callErr
					}
					arg__23173_59, callErr = rt.InvokeValue(head__23168_56, []vm.Value{arg__23172_58})
					if callErr != nil {
						return nil, callErr
					}
					v78, callErr = rt.InvokeValue(rt.LookupVar("ir.dominance", "dominates?").Deref(), []vm.Value{f_46, header_49, pred_52})
					if callErr != nil {
						return nil, callErr
					}
					if vm.IsTruthy(v78) {
						f_61 = f_46
						for__a18459_62 = for__a18459_47
						for__s_63 = for__s_48
						header_64 = header_49
						tem__G__0_65 = tem__G__0_50
						for__xs_66 = for__xs_51
						pred_67 = pred_52
						head__23174_68 = rt.LookupVar("clojure.core", "concat-list").Deref()
						goto b7
					} else {
						f_69 = f_46
						for__a18459_70 = for__a18459_47
						for__s_71 = for__s_48
						header_72 = header_49
						tem__G__0_73 = tem__G__0_50
						for__xs_74 = for__xs_51
						pred_75 = pred_52
						head__23174_76 = rt.LookupVar("clojure.core", "concat-list").Deref()
						goto b8
					}
				b7:
					;
					arg__23185_81, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{pred_67, header_64})
					if callErr != nil {
						return nil, callErr
					}
					arg__23190_84, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{pred_67, header_64})
					if callErr != nil {
						return nil, callErr
					}
					v85, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "list").Deref(), []vm.Value{arg__23190_84})
					if callErr != nil {
						return nil, callErr
					}
					arg__23191_89 = v85
					f_90 = f_61
					for__a18459_91 = for__a18459_62
					for__s_92 = for__s_63
					header_93 = header_64
					tem__G__0_94 = tem__G__0_65
					for__xs_95 = for__xs_66
					pred_96 = pred_67
					head__23174_97 = head__23174_68
					goto b9
				b8:
					;
					arg__23191_89 = vm.NIL
					f_90 = f_69
					for__a18459_91 = for__a18459_70
					for__s_92 = for__s_71
					header_93 = header_72
					tem__G__0_94 = tem__G__0_73
					for__xs_95 = for__xs_74
					pred_96 = pred_75
					head__23174_97 = head__23174_76
					goto b9
				b9:
					;
					arg__23195_99, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{for__xs_95})
					if callErr != nil {
						return nil, callErr
					}
					head__23199_101, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{for__a18459_91})
					if callErr != nil {
						return nil, callErr
					}
					arg__23203_103, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{for__xs_95})
					if callErr != nil {
						return nil, callErr
					}
					arg__23204_104, callErr = rt.InvokeValue(head__23199_101, []vm.Value{arg__23203_103})
					if callErr != nil {
						return nil, callErr
					}
					v105, callErr = rt.InvokeValue(head__23174_97, []vm.Value{arg__23191_89, arg__23204_104})
					if callErr != nil {
						return nil, callErr
					}
					v109 = v105
					f_110 = f_90
					for__a18459_111 = for__a18459_91
					for__s_112 = for__s_92
					header_113 = header_93
					tem__G__0_114 = tem__G__0_94
					goto b3
				})})
				if callErr != nil {
					return nil, callErr
				}
				return v18, nil
			})})
			if callErr != nil {
				return nil, callErr
			}
			for__iter18458_35, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{for__a18459_19})
			if callErr != nil {
				return nil, callErr
			}
			arg__23215_37, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-preds").Deref(), []vm.Value{header_15, f_5})
			if callErr != nil {
				return nil, callErr
			}
			arg__23222_39, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-preds").Deref(), []vm.Value{header_15, f_5})
			if callErr != nil {
				return nil, callErr
			}
			arg__23223_40, callErr = rt.InvokeValue(for__iter18458_35, []vm.Value{arg__23222_39})
			if callErr != nil {
				return nil, callErr
			}
			arg__23227_42, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{tem__G__0_8})
			if callErr != nil {
				return nil, callErr
			}
			head__23231_44, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{for__a18461_6})
			if callErr != nil {
				return nil, callErr
			}
			arg__23235_46, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{tem__G__0_8})
			if callErr != nil {
				return nil, callErr
			}
			arg__23236_47, callErr = rt.InvokeValue(head__23231_44, []vm.Value{arg__23235_46})
			if callErr != nil {
				return nil, callErr
			}
			for__a18459_52, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "atom").Deref(), []vm.Value{vm.NIL})
			if callErr != nil {
				return nil, callErr
			}
			__66, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "reset!").Deref(), []vm.Value{for__a18459_52, rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
				var v18 vm.Value
				var callErr error
				_ = v18
				v18, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "lazy-seq*").Deref(), []vm.Value{rt.BoxNativeFn(func() (vm.Value, error) {
					var tem__G__0_5 vm.Value
					var f_6 vm.Value
					var for__a18459_7 vm.Value
					var for__s_8 vm.Value
					var header_9 vm.Value
					var tem__G__0_10 vm.Value
					var pred_18 vm.Value
					var v34 vm.Value
					var f_11 vm.Value
					var for__a18459_12 vm.Value
					var for__s_13 vm.Value
					var header_14 vm.Value
					var tem__G__0_15 vm.Value
					var v109 vm.Value
					var f_110 vm.Value
					var for__a18459_111 vm.Value
					var for__s_112 vm.Value
					var header_113 vm.Value
					var tem__G__0_114 vm.Value
					var f_19 vm.Value
					var for__a18459_20 vm.Value
					var for__s_21 vm.Value
					var header_22 vm.Value
					var tem__G__0_23 vm.Value
					var for__xs_24 vm.Value
					var pred_25 vm.Value
					var arg__23467_37 vm.Value
					var arg__23472_40 vm.Value
					var v41 vm.Value
					var f_26 vm.Value
					var for__a18459_27 vm.Value
					var for__s_28 vm.Value
					var header_29 vm.Value
					var tem__G__0_30 vm.Value
					var for__xs_31 vm.Value
					var pred_32 vm.Value
					var arg__23473_45 vm.Value
					var f_46 vm.Value
					var for__a18459_47 vm.Value
					var for__s_48 vm.Value
					var header_49 vm.Value
					var tem__G__0_50 vm.Value
					var for__xs_51 vm.Value
					var pred_52 vm.Value
					var arg__23477_54 vm.Value
					var head__23481_56 vm.Value
					var arg__23485_58 vm.Value
					var arg__23486_59 vm.Value
					var v78 vm.Value
					var f_61 vm.Value
					var for__a18459_62 vm.Value
					var for__s_63 vm.Value
					var header_64 vm.Value
					var tem__G__0_65 vm.Value
					var for__xs_66 vm.Value
					var pred_67 vm.Value
					var head__23487_68 vm.Value
					var arg__23498_81 vm.Value
					var arg__23503_84 vm.Value
					var v85 vm.Value
					var f_69 vm.Value
					var for__a18459_70 vm.Value
					var for__s_71 vm.Value
					var header_72 vm.Value
					var tem__G__0_73 vm.Value
					var for__xs_74 vm.Value
					var pred_75 vm.Value
					var head__23487_76 vm.Value
					var arg__23504_89 vm.Value
					var f_90 vm.Value
					var for__a18459_91 vm.Value
					var for__s_92 vm.Value
					var header_93 vm.Value
					var tem__G__0_94 vm.Value
					var for__xs_95 vm.Value
					var pred_96 vm.Value
					var head__23487_97 vm.Value
					var arg__23508_99 vm.Value
					var head__23512_101 vm.Value
					var arg__23516_103 vm.Value
					var arg__23517_104 vm.Value
					var v105 vm.Value
					var callErr error
					_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = tem__G__0_5, f_6, for__a18459_7, for__s_8, header_9, tem__G__0_10, pred_18, v34, f_11, for__a18459_12, for__s_13, header_14, tem__G__0_15, v109, f_110, for__a18459_111, for__s_112, header_113, tem__G__0_114, f_19, for__a18459_20, for__s_21, header_22, tem__G__0_23, for__xs_24, pred_25, arg__23467_37, arg__23472_40, v41, f_26, for__a18459_27, for__s_28, header_29, tem__G__0_30, for__xs_31, pred_32, arg__23473_45, f_46, for__a18459_47, for__s_48, header_49, tem__G__0_50, for__xs_51, pred_52, arg__23477_54, head__23481_56, arg__23485_58, arg__23486_59, v78, f_61, for__a18459_62, for__s_63, header_64, tem__G__0_65, for__xs_66, pred_67, head__23487_68, arg__23498_81, arg__23503_84, v85, f_69, for__a18459_70, for__s_71, header_72, tem__G__0_73, for__xs_74, pred_75, head__23487_76, arg__23504_89, f_90, for__a18459_91, for__s_92, header_93, tem__G__0_94, for__xs_95, pred_96, head__23487_97, arg__23508_99, head__23512_101, arg__23516_103, arg__23517_104, v105
					tem__G__0_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{for__s_7})
					if callErr != nil {
						return nil, callErr
					}
					if vm.IsTruthy(tem__G__0_5) {
						f_6 = f_5
						for__a18459_7 = for__a18459_52
						for__s_8 = for__s_7
						header_9 = header_15
						tem__G__0_10 = tem__G__0_5
						goto b1
					} else {
						f_11 = f_5
						for__a18459_12 = for__a18459_52
						for__s_13 = for__s_7
						header_14 = header_15
						tem__G__0_15 = tem__G__0_5
						goto b2
					}
				b1:
					;
					pred_18, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{tem__G__0_10})
					if callErr != nil {
						return nil, callErr
					}
					v34, callErr = rt.InvokeValue(rt.LookupVar("ir.dominance", "dominates?").Deref(), []vm.Value{f_6, header_9, pred_18})
					if callErr != nil {
						return nil, callErr
					}
					if vm.IsTruthy(v34) {
						f_19 = f_6
						for__a18459_20 = for__a18459_7
						for__s_21 = for__s_8
						header_22 = header_9
						tem__G__0_23 = tem__G__0_10
						for__xs_24 = tem__G__0_10
						pred_25 = pred_18
						goto b4
					} else {
						f_26 = f_6
						for__a18459_27 = for__a18459_7
						for__s_28 = for__s_8
						header_29 = header_9
						tem__G__0_30 = tem__G__0_10
						for__xs_31 = tem__G__0_10
						pred_32 = pred_18
						goto b5
					}
				b2:
					;
					v109 = vm.NIL
					f_110 = f_11
					for__a18459_111 = for__a18459_12
					for__s_112 = for__s_13
					header_113 = header_14
					tem__G__0_114 = tem__G__0_15
					goto b3
				b3:
					;
					return v109, nil
				b4:
					;
					arg__23467_37, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{pred_25, header_22})
					if callErr != nil {
						return nil, callErr
					}
					arg__23472_40, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{pred_25, header_22})
					if callErr != nil {
						return nil, callErr
					}
					v41, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "list").Deref(), []vm.Value{arg__23472_40})
					if callErr != nil {
						return nil, callErr
					}
					arg__23473_45 = v41
					f_46 = f_19
					for__a18459_47 = for__a18459_20
					for__s_48 = for__s_21
					header_49 = header_22
					tem__G__0_50 = tem__G__0_23
					for__xs_51 = for__xs_24
					pred_52 = pred_25
					goto b6
				b5:
					;
					arg__23473_45 = vm.NIL
					f_46 = f_26
					for__a18459_47 = for__a18459_27
					for__s_48 = for__s_28
					header_49 = header_29
					tem__G__0_50 = tem__G__0_30
					for__xs_51 = for__xs_31
					pred_52 = pred_32
					goto b6
				b6:
					;
					arg__23477_54, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{for__xs_51})
					if callErr != nil {
						return nil, callErr
					}
					head__23481_56, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{for__a18459_47})
					if callErr != nil {
						return nil, callErr
					}
					arg__23485_58, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{for__xs_51})
					if callErr != nil {
						return nil, callErr
					}
					arg__23486_59, callErr = rt.InvokeValue(head__23481_56, []vm.Value{arg__23485_58})
					if callErr != nil {
						return nil, callErr
					}
					v78, callErr = rt.InvokeValue(rt.LookupVar("ir.dominance", "dominates?").Deref(), []vm.Value{f_46, header_49, pred_52})
					if callErr != nil {
						return nil, callErr
					}
					if vm.IsTruthy(v78) {
						f_61 = f_46
						for__a18459_62 = for__a18459_47
						for__s_63 = for__s_48
						header_64 = header_49
						tem__G__0_65 = tem__G__0_50
						for__xs_66 = for__xs_51
						pred_67 = pred_52
						head__23487_68 = rt.LookupVar("clojure.core", "concat-list").Deref()
						goto b7
					} else {
						f_69 = f_46
						for__a18459_70 = for__a18459_47
						for__s_71 = for__s_48
						header_72 = header_49
						tem__G__0_73 = tem__G__0_50
						for__xs_74 = for__xs_51
						pred_75 = pred_52
						head__23487_76 = rt.LookupVar("clojure.core", "concat-list").Deref()
						goto b8
					}
				b7:
					;
					arg__23498_81, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{pred_67, header_64})
					if callErr != nil {
						return nil, callErr
					}
					arg__23503_84, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{pred_67, header_64})
					if callErr != nil {
						return nil, callErr
					}
					v85, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "list").Deref(), []vm.Value{arg__23503_84})
					if callErr != nil {
						return nil, callErr
					}
					arg__23504_89 = v85
					f_90 = f_61
					for__a18459_91 = for__a18459_62
					for__s_92 = for__s_63
					header_93 = header_64
					tem__G__0_94 = tem__G__0_65
					for__xs_95 = for__xs_66
					pred_96 = pred_67
					head__23487_97 = head__23487_68
					goto b9
				b8:
					;
					arg__23504_89 = vm.NIL
					f_90 = f_69
					for__a18459_91 = for__a18459_70
					for__s_92 = for__s_71
					header_93 = header_72
					tem__G__0_94 = tem__G__0_73
					for__xs_95 = for__xs_74
					pred_96 = pred_75
					head__23487_97 = head__23487_76
					goto b9
				b9:
					;
					arg__23508_99, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{for__xs_95})
					if callErr != nil {
						return nil, callErr
					}
					head__23512_101, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{for__a18459_91})
					if callErr != nil {
						return nil, callErr
					}
					arg__23516_103, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{for__xs_95})
					if callErr != nil {
						return nil, callErr
					}
					arg__23517_104, callErr = rt.InvokeValue(head__23512_101, []vm.Value{arg__23516_103})
					if callErr != nil {
						return nil, callErr
					}
					v105, callErr = rt.InvokeValue(head__23487_97, []vm.Value{arg__23504_89, arg__23517_104})
					if callErr != nil {
						return nil, callErr
					}
					v109 = v105
					f_110 = f_90
					for__a18459_111 = for__a18459_91
					for__s_112 = for__s_92
					header_113 = header_93
					tem__G__0_114 = tem__G__0_94
					goto b3
				})})
				if callErr != nil {
					return nil, callErr
				}
				return v18, nil
			})})
			if callErr != nil {
				return nil, callErr
			}
			for__iter18458_68, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{for__a18459_52})
			if callErr != nil {
				return nil, callErr
			}
			arg__23528_70, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-preds").Deref(), []vm.Value{header_15, f_5})
			if callErr != nil {
				return nil, callErr
			}
			arg__23535_72, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-preds").Deref(), []vm.Value{header_15, f_5})
			if callErr != nil {
				return nil, callErr
			}
			arg__23536_73, callErr = rt.InvokeValue(for__iter18458_68, []vm.Value{arg__23535_72})
			if callErr != nil {
				return nil, callErr
			}
			arg__23540_75, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{tem__G__0_8})
			if callErr != nil {
				return nil, callErr
			}
			head__23544_77, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{for__a18461_6})
			if callErr != nil {
				return nil, callErr
			}
			arg__23548_79, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{tem__G__0_8})
			if callErr != nil {
				return nil, callErr
			}
			arg__23549_80, callErr = rt.InvokeValue(head__23544_77, []vm.Value{arg__23548_79})
			if callErr != nil {
				return nil, callErr
			}
			v81, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "concat-list").Deref(), []vm.Value{arg__23536_73, arg__23549_80})
			if callErr != nil {
				return nil, callErr
			}
			v85 = v81
			f_86 = f_5
			for__a18461_87 = for__a18461_6
			for__s_88 = for__s_7
			tem__G__0_89 = tem__G__0_8
			goto b3
		b2:
			;
			v85 = vm.NIL
			f_86 = f_9
			for__a18461_87 = for__a18461_10
			for__s_88 = for__s_11
			tem__G__0_89 = tem__G__0_12
			goto b3
		b3:
			;
			return v85, nil
		})})
		if callErr != nil {
			return nil, callErr
		}
		return v14, nil
	})})
	if callErr != nil {
		return nil, callErr
	}
	for__iter18460_39, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{for__a18461_27})
	if callErr != nil {
		return nil, callErr
	}
	arg__23558_41, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__23563_43, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__23564_44, callErr = rt.InvokeValue(for__iter18460_39, []vm.Value{arg__23563_43})
	if callErr != nil {
		return nil, callErr
	}
	v45, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg__23564_44})
	if callErr != nil {
		return nil, callErr
	}
	return v45, nil
}
func pure_op_QMARK_(arg0 vm.Value) (vm.Value, error) {
	var v4 vm.Value
	var callErr error
	_ = v4
	v4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "contains?").Deref(), []vm.Value{rt.LookupVar("ir.passes.licm", "pure-ops").Deref(), arg0})
	if callErr != nil {
		return nil, callErr
	}
	return v4, nil
}
func unique_pre_header(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__23576_5 vm.Value
	var arg__23584_8 vm.Value
	var outs_9 vm.Value
	var arg__23589_20 vm.Value
	var v21 bool
	var f_10 vm.Value
	var header_11 vm.Value
	var loop_set_12 vm.Value
	var outs_13 vm.Value
	var v24 vm.Value
	var f_14 vm.Value
	var header_15 vm.Value
	var loop_set_16 vm.Value
	var outs_17 vm.Value
	var v28 vm.Value
	var f_29 vm.Value
	var header_30 vm.Value
	var loop_set_31 vm.Value
	var outs_32 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__23576_5, arg__23584_8, outs_9, arg__23589_20, v21, f_10, header_11, loop_set_12, outs_13, v24, f_14, header_15, loop_set_16, outs_17, v28, f_29, header_30, loop_set_31, outs_32
	arg__23576_5, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-preds").Deref(), []vm.Value{arg1, arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__23584_8, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-preds").Deref(), []vm.Value{arg1, arg0})
	if callErr != nil {
		return nil, callErr
	}
	outs_9, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "remove").Deref(), []vm.Value{arg2, arg__23584_8})
	if callErr != nil {
		return nil, callErr
	}
	arg__23589_20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{outs_9})
	if callErr != nil {
		return nil, callErr
	}
	v21 = arg__23589_20 == vm.Int(1)
	if v21 {
		f_10 = arg0
		header_11 = arg1
		loop_set_12 = arg2
		outs_13 = outs_9
		goto b1
	} else {
		f_14 = arg0
		header_15 = arg1
		loop_set_16 = arg2
		outs_17 = outs_9
		goto b2
	}
b1:
	;
	v24, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{outs_13})
	if callErr != nil {
		return nil, callErr
	}
	v28 = v24
	f_29 = f_10
	header_30 = header_11
	loop_set_31 = loop_set_12
	outs_32 = outs_13
	goto b3
b2:
	;
	v28 = vm.NIL
	f_29 = f_14
	header_30 = header_15
	loop_set_31 = loop_set_16
	outs_32 = outs_17
	goto b3
b3:
	;
	return v28, nil
}
func loop_blocks(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var pred_7 vm.Value
	var header_13 vm.Value
	var arg__24325_16 vm.Value
	var arg__24335_19 vm.Value
	var v20 vm.Value
	var callErr error
	_, _, _, _, _ = pred_7, header_13, arg__24325_16, arg__24335_19, v20
	pred_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg1, vm.Int(0), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	header_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg1, vm.Int(1), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	arg__24325_16, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.licm", "reachable-without").Deref(), []vm.Value{arg0, pred_7, header_13})
	if callErr != nil {
		return nil, callErr
	}
	arg__24335_19, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.licm", "reachable-without").Deref(), []vm.Value{arg0, pred_7, header_13})
	if callErr != nil {
		return nil, callErr
	}
	v20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{arg__24335_19, header_13})
	if callErr != nil {
		return nil, callErr
	}
	return v20, nil
}
func find_in_loop_users(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__24340_5 vm.Value
	var arg__24346_8 vm.Value
	var clone_users_9 vm.Value
	var acc_13 vm.Value
	var v27 vm.Value
	var v29 vm.Value
	var callErr error
	_, _, _, _, _, _ = arg__24340_5, arg__24346_8, clone_users_9, acc_13, v27, v29
	arg__24340_5, callErr = rt.InvokeValue(rt.LookupVar("ir", "uses").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__24346_8, callErr = rt.InvokeValue(rt.LookupVar("ir", "uses").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	clone_users_9, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__24346_8, arg1})
	if callErr != nil {
		return nil, callErr
	}
	acc_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "atom").Deref(), []vm.Value{vm.NewArrayVector([]vm.Value{})})
	if callErr != nil {
		return nil, callErr
	}
	v27, callErr = rt.InvokeValue(rt.LookupVar("ir", "uses-for-each").Deref(), []vm.Value{clone_users_9, rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var arg__24415_17 vm.Value
		var arg__24423_21 vm.Value
		var and__x_22 vm.Value
		var nid_5 vm.Value
		var acc_6 vm.Value
		var clone_7 vm.Value
		var f_8 vm.Value
		var loop_set_9 vm.Value
		var v88 vm.Value
		var nid_10 vm.Value
		var acc_11 vm.Value
		var clone_12 vm.Value
		var f_13 vm.Value
		var loop_set_14 vm.Value
		var v92 vm.Value
		var nid_93 vm.Value
		var acc_94 vm.Value
		var clone_95 vm.Value
		var f_96 vm.Value
		var loop_set_97 vm.Value
		var nid_23 vm.Value
		var acc_24 vm.Value
		var clone_25 vm.Value
		var f_26 vm.Value
		var loop_set_27 vm.Value
		var and__x_28 vm.Value
		var arg__24429_37 vm.Value
		var arg__24436_39 vm.Value
		var and__x_40 vm.Value
		var nid_29 vm.Value
		var acc_30 vm.Value
		var clone_31 vm.Value
		var f_32 vm.Value
		var loop_set_33 vm.Value
		var and__x_34 vm.Value
		var v77 vm.Value
		var nid_78 vm.Value
		var acc_79 vm.Value
		var clone_80 vm.Value
		var f_81 vm.Value
		var loop_set_82 vm.Value
		var and__x_83 vm.Value
		var nid_41 vm.Value
		var acc_42 vm.Value
		var clone_43 vm.Value
		var f_44 vm.Value
		var loop_set_45 vm.Value
		var and__x_46 vm.Value
		var arg__24445_58 vm.Value
		var arg__24455_64 vm.Value
		var v65 vm.Value
		var nid_47 vm.Value
		var acc_48 vm.Value
		var clone_49 vm.Value
		var f_50 vm.Value
		var loop_set_51 vm.Value
		var and__x_52 vm.Value
		var v68 vm.Value
		var nid_69 vm.Value
		var acc_70 vm.Value
		var clone_71 vm.Value
		var f_72 vm.Value
		var loop_set_73 vm.Value
		var and__x_74 vm.Value
		var callErr error
		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__24415_17, arg__24423_21, and__x_22, nid_5, acc_6, clone_7, f_8, loop_set_9, v88, nid_10, acc_11, clone_12, f_13, loop_set_14, v92, nid_93, acc_94, clone_95, f_96, loop_set_97, nid_23, acc_24, clone_25, f_26, loop_set_27, and__x_28, arg__24429_37, arg__24436_39, and__x_40, nid_29, acc_30, clone_31, f_32, loop_set_33, and__x_34, v77, nid_78, acc_79, clone_80, f_81, loop_set_82, and__x_83, nid_41, acc_42, clone_43, f_44, loop_set_45, and__x_46, arg__24445_58, arg__24455_64, v65, nid_47, acc_48, clone_49, f_50, loop_set_51, and__x_52, v68, nid_69, acc_70, clone_71, f_72, loop_set_73, and__x_74
		arg__24415_17, callErr = rt.InvokeValue(rt.LookupVar("ir", "op").Deref(), []vm.Value{arg0, arg0})
		if callErr != nil {
			return nil, callErr
		}
		arg__24423_21, callErr = rt.InvokeValue(rt.LookupVar("ir", "op").Deref(), []vm.Value{arg0, arg0})
		if callErr != nil {
			return nil, callErr
		}
		and__x_22, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not=").Deref(), []vm.Value{vm.Keyword("invalid"), arg__24423_21})
		if callErr != nil {
			return nil, callErr
		}
		if vm.IsTruthy(and__x_22) {
			nid_23 = arg0
			acc_24 = acc_13
			clone_25 = arg1
			f_26 = arg0
			loop_set_27 = arg2
			and__x_28 = and__x_22
			goto b4
		} else {
			nid_29 = arg0
			acc_30 = acc_13
			clone_31 = arg1
			f_32 = arg0
			loop_set_33 = arg2
			and__x_34 = and__x_22
			goto b5
		}
	b1:
		;
		v88, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "swap!").Deref(), []vm.Value{acc_6, rt.LookupVar("clojure.core", "conj").Deref(), nid_5})
		if callErr != nil {
			return nil, callErr
		}
		v92 = v88
		nid_93 = nid_5
		acc_94 = acc_6
		clone_95 = clone_7
		f_96 = f_8
		loop_set_97 = loop_set_9
		goto b3
	b2:
		;
		v92 = vm.NIL
		nid_93 = nid_10
		acc_94 = acc_11
		clone_95 = clone_12
		f_96 = f_13
		loop_set_97 = loop_set_14
		goto b3
	b3:
		;
		return v92, nil
	b4:
		;
		arg__24429_37, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-of").Deref(), []vm.Value{nid_23, f_26})
		if callErr != nil {
			return nil, callErr
		}
		arg__24436_39, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-of").Deref(), []vm.Value{nid_23, f_26})
		if callErr != nil {
			return nil, callErr
		}
		and__x_40, callErr = rt.InvokeValue(loop_set_27, []vm.Value{arg__24436_39})
		if callErr != nil {
			return nil, callErr
		}
		if vm.IsTruthy(and__x_40) {
			nid_41 = nid_23
			acc_42 = acc_24
			clone_43 = clone_25
			f_44 = f_26
			loop_set_45 = loop_set_27
			and__x_46 = and__x_40
			goto b7
		} else {
			nid_47 = nid_23
			acc_48 = acc_24
			clone_49 = clone_25
			f_50 = f_26
			loop_set_51 = loop_set_27
			and__x_52 = and__x_40
			goto b8
		}
	b5:
		;
		v77 = and__x_34
		nid_78 = nid_29
		acc_79 = acc_30
		clone_80 = clone_31
		f_81 = f_32
		loop_set_82 = loop_set_33
		and__x_83 = and__x_34
		goto b6
	b6:
		;
		if vm.IsTruthy(v77) {
			nid_5 = nid_78
			acc_6 = acc_79
			clone_7 = clone_80
			f_8 = f_81
			loop_set_9 = loop_set_82
			goto b1
		} else {
			nid_10 = nid_78
			acc_11 = acc_79
			clone_12 = clone_80
			f_13 = f_81
			loop_set_14 = loop_set_82
			goto b2
		}
	b7:
		;
		arg__24445_58, callErr = rt.InvokeValue(rt.LookupVar("ir", "refs").Deref(), []vm.Value{nid_41, f_44})
		if callErr != nil {
			return nil, callErr
		}
		arg__24455_64, callErr = rt.InvokeValue(rt.LookupVar("ir", "refs").Deref(), []vm.Value{nid_41, f_44})
		if callErr != nil {
			return nil, callErr
		}
		v65, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "some").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) vm.Value {
			var v2 vm.Value
			_ = v2
			v2 = vm.Boolean(arg0 == clone_43)
			return v2
		}), arg__24455_64})
		if callErr != nil {
			return nil, callErr
		}
		v68 = v65
		nid_69 = nid_41
		acc_70 = acc_42
		clone_71 = clone_43
		f_72 = f_44
		loop_set_73 = loop_set_45
		and__x_74 = and__x_46
		goto b9
	b8:
		;
		v68 = and__x_52
		nid_69 = nid_47
		acc_70 = acc_48
		clone_71 = clone_49
		f_72 = f_50
		loop_set_73 = loop_set_51
		and__x_74 = and__x_52
		goto b9
	b9:
		;
		v77 = v68
		nid_78 = nid_69
		acc_79 = acc_70
		clone_80 = clone_71
		f_81 = f_72
		loop_set_82 = loop_set_73
		and__x_83 = and__x_28
		goto b6
	})})
	if callErr != nil {
		return nil, callErr
	}
	v29, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{acc_13})
	if callErr != nil {
		return nil, callErr
	}
	return v29, nil
}
func licm_one_loop(arg0 vm.Value, arg1 vm.Value) (vm.Value, error) {
	var header_6 vm.Value
	var loop_set_8 vm.Value
	var pre_header_10 vm.Value
	var f_11 vm.Value
	var back_edge_12 vm.Value
	var header_13 vm.Value
	var loop_set_14 vm.Value
	var pre_header_15 vm.Value
	var hoistable_23 vm.Value
	var arg__25116_29 vm.Value
	var arg__25132_36 vm.Value
	var hoisted_pairs_37 vm.Value
	var v53 vm.Value
	var f_16 vm.Value
	var back_edge_17 vm.Value
	var header_18 vm.Value
	var loop_set_19 vm.Value
	var pre_header_20 vm.Value
	var v81 vm.Value
	var f_82 vm.Value
	var back_edge_83 vm.Value
	var header_84 vm.Value
	var loop_set_85 vm.Value
	var pre_header_86 vm.Value
	var f_38 vm.Value
	var back_edge_39 vm.Value
	var header_40 vm.Value
	var loop_set_41 vm.Value
	var pre_header_42 vm.Value
	var hoistable_43 vm.Value
	var hoisted_pairs_44 vm.Value
	var arg__25144_59 vm.Value
	var arg__25155_65 vm.Value
	var v66 vm.Value
	var f_45 vm.Value
	var back_edge_46 vm.Value
	var header_47 vm.Value
	var loop_set_48 vm.Value
	var pre_header_49 vm.Value
	var hoistable_50 vm.Value
	var hoisted_pairs_51 vm.Value
	var v70 vm.Value
	var f_71 vm.Value
	var back_edge_72 vm.Value
	var header_73 vm.Value
	var loop_set_74 vm.Value
	var pre_header_75 vm.Value
	var hoistable_76 vm.Value
	var hoisted_pairs_77 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = header_6, loop_set_8, pre_header_10, f_11, back_edge_12, header_13, loop_set_14, pre_header_15, hoistable_23, arg__25116_29, arg__25132_36, hoisted_pairs_37, v53, f_16, back_edge_17, header_18, loop_set_19, pre_header_20, v81, f_82, back_edge_83, header_84, loop_set_85, pre_header_86, f_38, back_edge_39, header_40, loop_set_41, pre_header_42, hoistable_43, hoisted_pairs_44, arg__25144_59, arg__25155_65, v66, f_45, back_edge_46, header_47, loop_set_48, pre_header_49, hoistable_50, hoisted_pairs_51, v70, f_71, back_edge_72, header_73, loop_set_74, pre_header_75, hoistable_76, hoisted_pairs_77
	header_6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg1, vm.Int(1)})
	if callErr != nil {
		return nil, callErr
	}
	loop_set_8, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.licm", "loop-blocks").Deref(), []vm.Value{arg0, arg1})
	if callErr != nil {
		return nil, callErr
	}
	pre_header_10, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.licm", "unique-pre-header").Deref(), []vm.Value{arg0, header_6, loop_set_8})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(pre_header_10) {
		f_11 = arg0
		back_edge_12 = arg1
		header_13 = header_6
		loop_set_14 = loop_set_8
		pre_header_15 = pre_header_10
		goto b1
	} else {
		f_16 = arg0
		back_edge_17 = arg1
		header_18 = header_6
		loop_set_19 = loop_set_8
		pre_header_20 = pre_header_10
		goto b2
	}
b1:
	;
	hoistable_23, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.licm", "collect-hoistable").Deref(), []vm.Value{f_11, loop_set_14})
	if callErr != nil {
		return nil, callErr
	}
	arg__25116_29, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "sort").Deref(), []vm.Value{hoistable_23})
	if callErr != nil {
		return nil, callErr
	}
	arg__25132_36, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "sort").Deref(), []vm.Value{hoistable_23})
	if callErr != nil {
		return nil, callErr
	}
	hoisted_pairs_37, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapv").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var arg__25127_5 vm.Value
		var v6 vm.Value
		var callErr error
		_, _ = arg__25127_5, v6
		arg__25127_5, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.licm", "hoist-one!").Deref(), []vm.Value{f_11, arg0, pre_header_15})
		if callErr != nil {
			return nil, callErr
		}
		v6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg0, arg__25127_5})
		if callErr != nil {
			return nil, callErr
		}
		return v6, nil
	}), arg__25132_36})
	if callErr != nil {
		return nil, callErr
	}
	v53, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{hoisted_pairs_37})
	if callErr != nil {
		return nil, callErr
	}
	if vm.IsTruthy(v53) {
		f_38 = f_11
		back_edge_39 = back_edge_12
		header_40 = header_13
		loop_set_41 = loop_set_14
		pre_header_42 = pre_header_15
		hoistable_43 = hoistable_23
		hoisted_pairs_44 = hoisted_pairs_37
		goto b4
	} else {
		f_45 = f_11
		back_edge_46 = back_edge_12
		header_47 = header_13
		loop_set_48 = loop_set_14
		pre_header_49 = pre_header_15
		hoistable_50 = hoistable_23
		hoisted_pairs_51 = hoisted_pairs_37
		goto b5
	}
b2:
	;
	v81 = vm.NIL
	f_82 = f_16
	back_edge_83 = back_edge_17
	header_84 = header_18
	loop_set_85 = loop_set_19
	pre_header_86 = pre_header_20
	goto b3
b3:
	;
	return v81, nil
b4:
	;
	arg__25144_59, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("body"), loop_set_41, vm.Keyword("preheader"), pre_header_42, vm.Keyword("header"), header_40})
	if callErr != nil {
		return nil, callErr
	}
	arg__25155_65, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("body"), loop_set_41, vm.Keyword("preheader"), pre_header_42, vm.Keyword("header"), header_40})
	if callErr != nil {
		return nil, callErr
	}
	v66, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.licm", "thread-hoisted-through-header!").Deref(), []vm.Value{f_38, arg__25155_65, hoisted_pairs_44})
	if callErr != nil {
		return nil, callErr
	}
	v70 = v66
	f_71 = f_38
	back_edge_72 = back_edge_39
	header_73 = header_40
	loop_set_74 = loop_set_41
	pre_header_75 = pre_header_42
	hoistable_76 = hoistable_43
	hoisted_pairs_77 = hoisted_pairs_44
	goto b6
b5:
	;
	v70 = vm.NIL
	f_71 = f_45
	back_edge_72 = back_edge_46
	header_73 = header_47
	loop_set_74 = loop_set_48
	pre_header_75 = pre_header_49
	hoistable_76 = hoistable_50
	hoisted_pairs_77 = hoisted_pairs_51
	goto b6
b6:
	;
	v81 = v70
	f_82 = f_71
	back_edge_83 = back_edge_72
	header_84 = header_73
	loop_set_85 = loop_set_74
	pre_header_86 = pre_header_75
	goto b3
}
func operand_defined_outside_QMARK_(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var refs_5 vm.Value
	var v15 vm.Value
	var callErr error
	_, _ = refs_5, v15
	refs_5, callErr = rt.InvokeValue(rt.LookupVar("ir", "refs").Deref(), []vm.Value{arg1, arg0})
	if callErr != nil {
		return nil, callErr
	}
	v15, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "every?").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var arg__25204_4 vm.Value
		var arg__25212_7 vm.Value
		var arg__25213_8 vm.Value
		var arg__25221_11 vm.Value
		var arg__25229_14 vm.Value
		var arg__25230_15 vm.Value
		var v16 vm.Value
		var callErr error
		_, _, _, _, _, _, _ = arg__25204_4, arg__25212_7, arg__25213_8, arg__25221_11, arg__25229_14, arg__25230_15, v16
		arg__25204_4, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-of").Deref(), []vm.Value{arg0, arg0})
		if callErr != nil {
			return nil, callErr
		}
		arg__25212_7, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-of").Deref(), []vm.Value{arg0, arg0})
		if callErr != nil {
			return nil, callErr
		}
		arg__25213_8, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "contains?").Deref(), []vm.Value{arg2, arg__25212_7})
		if callErr != nil {
			return nil, callErr
		}
		arg__25221_11, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-of").Deref(), []vm.Value{arg0, arg0})
		if callErr != nil {
			return nil, callErr
		}
		arg__25229_14, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-of").Deref(), []vm.Value{arg0, arg0})
		if callErr != nil {
			return nil, callErr
		}
		arg__25230_15, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "contains?").Deref(), []vm.Value{arg2, arg__25229_14})
		if callErr != nil {
			return nil, callErr
		}
		v16, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not").Deref(), []vm.Value{arg__25230_15})
		if callErr != nil {
			return nil, callErr
		}
		return v16, nil
	}), refs_5})
	if callErr != nil {
		return nil, callErr
	}
	return v15, nil
}
func licm(arg0 vm.Value) (vm.Value, error) {
	var arg__25238_3 vm.Value
	var arg__25243_6 vm.Value
	var doseq_seq__25233_7 vm.Value
	var doseq_loop__25234_8 vm.Value
	var f_9 vm.Value
	var doseq_seq__25233_11 vm.Value
	var doseq_loop__25234_12 vm.Value
	var f_13 vm.Value
	var be_19 vm.Value
	var v21 vm.Value
	var v23 vm.Value
	var doseq_seq__25233_14 vm.Value
	var doseq_loop__25234_15 vm.Value
	var f_16 vm.Value
	var v27 vm.Value
	var doseq_seq__25233_28 vm.Value
	var doseq_loop__25234_29 vm.Value
	var f_30 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__25238_3, arg__25243_6, doseq_seq__25233_7, doseq_loop__25234_8, f_9, doseq_seq__25233_11, doseq_loop__25234_12, f_13, be_19, v21, v23, doseq_seq__25233_14, doseq_loop__25234_15, f_16, v27, doseq_seq__25233_28, doseq_loop__25234_29, f_30
	arg__25238_3, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.licm", "back-edges").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__25243_6, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.licm", "back-edges").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	doseq_seq__25233_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{arg__25243_6})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__25234_8 = doseq_seq__25233_7
	f_9 = arg0
	goto b1
b1:
	;
	if vm.IsTruthy(doseq_loop__25234_8) {
		doseq_seq__25233_11 = doseq_seq__25233_7
		doseq_loop__25234_12 = doseq_loop__25234_8
		f_13 = f_9
		goto b2
	} else {
		doseq_seq__25233_14 = doseq_seq__25233_7
		doseq_loop__25234_15 = doseq_loop__25234_8
		f_16 = f_9
		goto b3
	}
b2:
	;
	be_19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{doseq_loop__25234_12})
	if callErr != nil {
		return nil, callErr
	}
	v21, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.licm", "licm-one-loop").Deref(), []vm.Value{f_13, be_19})
	if callErr != nil {
		return nil, callErr
	}
	v23, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "next").Deref(), []vm.Value{doseq_loop__25234_12})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__25234_8 = v23
	f_9 = f_13
	goto b1
b3:
	;
	v27 = vm.NIL
	doseq_seq__25233_28 = doseq_seq__25233_14
	doseq_loop__25234_29 = doseq_loop__25234_15
	f_30 = f_16
	goto b4
b4:
	;
	return f_30, nil
}
