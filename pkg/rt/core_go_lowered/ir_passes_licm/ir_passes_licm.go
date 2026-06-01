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
	var for__a18465_5 vm.Value
	var __15 vm.Value
	var for__iter18464_17 vm.Value
	var arg__21010_19 vm.Value
	var arg__21015_21 vm.Value
	var arg__21016_22 vm.Value
	var for__a18465_27 vm.Value
	var __37 vm.Value
	var for__iter18464_39 vm.Value
	var arg__23562_41 vm.Value
	var arg__23567_43 vm.Value
	var arg__23568_44 vm.Value
	var v45 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _ = for__a18465_5, __15, for__iter18464_17, arg__21010_19, arg__21015_21, arg__21016_22, for__a18465_27, __37, for__iter18464_39, arg__23562_41, arg__23567_43, arg__23568_44, v45
	for__a18465_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "atom").Deref(), []vm.Value{vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	__15, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "reset!").Deref(), []vm.Value{for__a18465_5, rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v14 vm.Value
		var callErr error
		_ = v14
		v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "lazy-seq*").Deref(), []vm.Value{rt.BoxNativeFn(func() (vm.Value, error) {
			var tem__G__0_4 vm.Value
			var f_5 vm.Value
			var for__a18465_6 vm.Value
			var for__s_7 vm.Value
			var tem__G__0_8 vm.Value
			var header_15 vm.Value
			var for__a18463_19 vm.Value
			var __33 vm.Value
			var for__iter18462_35 vm.Value
			var arg__20667_37 vm.Value
			var arg__20674_39 vm.Value
			var arg__20675_40 vm.Value
			var arg__20679_42 vm.Value
			var head__20683_44 vm.Value
			var arg__20687_46 vm.Value
			var arg__20688_47 vm.Value
			var for__a18463_52 vm.Value
			var __66 vm.Value
			var for__iter18462_68 vm.Value
			var arg__20980_70 vm.Value
			var arg__20987_72 vm.Value
			var arg__20988_73 vm.Value
			var arg__20992_75 vm.Value
			var head__20996_77 vm.Value
			var arg__21000_79 vm.Value
			var arg__21001_80 vm.Value
			var v81 vm.Value
			var f_9 vm.Value
			var for__a18465_10 vm.Value
			var for__s_11 vm.Value
			var tem__G__0_12 vm.Value
			var v85 vm.Value
			var f_86 vm.Value
			var for__a18465_87 vm.Value
			var for__s_88 vm.Value
			var tem__G__0_89 vm.Value
			var callErr error
			_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = tem__G__0_4, f_5, for__a18465_6, for__s_7, tem__G__0_8, header_15, for__a18463_19, __33, for__iter18462_35, arg__20667_37, arg__20674_39, arg__20675_40, arg__20679_42, head__20683_44, arg__20687_46, arg__20688_47, for__a18463_52, __66, for__iter18462_68, arg__20980_70, arg__20987_72, arg__20988_73, arg__20992_75, head__20996_77, arg__21000_79, arg__21001_80, v81, f_9, for__a18465_10, for__s_11, tem__G__0_12, v85, f_86, for__a18465_87, for__s_88, tem__G__0_89
			tem__G__0_4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{arg0})
			if callErr != nil {
				return nil, callErr
			}
			if vm.IsTruthy(tem__G__0_4) {
				f_5 = arg0
				for__a18465_6 = for__a18465_5
				for__s_7 = arg0
				tem__G__0_8 = tem__G__0_4
				goto b1
			} else {
				f_9 = arg0
				for__a18465_10 = for__a18465_5
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
			for__a18463_19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "atom").Deref(), []vm.Value{vm.NIL})
			if callErr != nil {
				return nil, callErr
			}
			__33, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "reset!").Deref(), []vm.Value{for__a18463_19, rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
				var v18 vm.Value
				var callErr error
				_ = v18
				v18, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "lazy-seq*").Deref(), []vm.Value{rt.BoxNativeFn(func() (vm.Value, error) {
					var tem__G__0_5 vm.Value
					var f_6 vm.Value
					var for__a18463_7 vm.Value
					var for__s_8 vm.Value
					var header_9 vm.Value
					var tem__G__0_10 vm.Value
					var pred_18 vm.Value
					var v34 vm.Value
					var f_11 vm.Value
					var for__a18463_12 vm.Value
					var for__s_13 vm.Value
					var header_14 vm.Value
					var tem__G__0_15 vm.Value
					var v109 vm.Value
					var f_110 vm.Value
					var for__a18463_111 vm.Value
					var for__s_112 vm.Value
					var header_113 vm.Value
					var tem__G__0_114 vm.Value
					var f_19 vm.Value
					var for__a18463_20 vm.Value
					var for__s_21 vm.Value
					var header_22 vm.Value
					var tem__G__0_23 vm.Value
					var for__xs_24 vm.Value
					var pred_25 vm.Value
					var arg__20606_37 vm.Value
					var arg__20611_40 vm.Value
					var v41 vm.Value
					var f_26 vm.Value
					var for__a18463_27 vm.Value
					var for__s_28 vm.Value
					var header_29 vm.Value
					var tem__G__0_30 vm.Value
					var for__xs_31 vm.Value
					var pred_32 vm.Value
					var arg__20612_45 vm.Value
					var f_46 vm.Value
					var for__a18463_47 vm.Value
					var for__s_48 vm.Value
					var header_49 vm.Value
					var tem__G__0_50 vm.Value
					var for__xs_51 vm.Value
					var pred_52 vm.Value
					var arg__20616_54 vm.Value
					var head__20620_56 vm.Value
					var arg__20624_58 vm.Value
					var arg__20625_59 vm.Value
					var v78 vm.Value
					var f_61 vm.Value
					var for__a18463_62 vm.Value
					var for__s_63 vm.Value
					var header_64 vm.Value
					var tem__G__0_65 vm.Value
					var for__xs_66 vm.Value
					var pred_67 vm.Value
					var head__20626_68 vm.Value
					var arg__20637_81 vm.Value
					var arg__20642_84 vm.Value
					var v85 vm.Value
					var f_69 vm.Value
					var for__a18463_70 vm.Value
					var for__s_71 vm.Value
					var header_72 vm.Value
					var tem__G__0_73 vm.Value
					var for__xs_74 vm.Value
					var pred_75 vm.Value
					var head__20626_76 vm.Value
					var arg__20643_89 vm.Value
					var f_90 vm.Value
					var for__a18463_91 vm.Value
					var for__s_92 vm.Value
					var header_93 vm.Value
					var tem__G__0_94 vm.Value
					var for__xs_95 vm.Value
					var pred_96 vm.Value
					var head__20626_97 vm.Value
					var arg__20647_99 vm.Value
					var head__20651_101 vm.Value
					var arg__20655_103 vm.Value
					var arg__20656_104 vm.Value
					var v105 vm.Value
					var callErr error
					_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = tem__G__0_5, f_6, for__a18463_7, for__s_8, header_9, tem__G__0_10, pred_18, v34, f_11, for__a18463_12, for__s_13, header_14, tem__G__0_15, v109, f_110, for__a18463_111, for__s_112, header_113, tem__G__0_114, f_19, for__a18463_20, for__s_21, header_22, tem__G__0_23, for__xs_24, pred_25, arg__20606_37, arg__20611_40, v41, f_26, for__a18463_27, for__s_28, header_29, tem__G__0_30, for__xs_31, pred_32, arg__20612_45, f_46, for__a18463_47, for__s_48, header_49, tem__G__0_50, for__xs_51, pred_52, arg__20616_54, head__20620_56, arg__20624_58, arg__20625_59, v78, f_61, for__a18463_62, for__s_63, header_64, tem__G__0_65, for__xs_66, pred_67, head__20626_68, arg__20637_81, arg__20642_84, v85, f_69, for__a18463_70, for__s_71, header_72, tem__G__0_73, for__xs_74, pred_75, head__20626_76, arg__20643_89, f_90, for__a18463_91, for__s_92, header_93, tem__G__0_94, for__xs_95, pred_96, head__20626_97, arg__20647_99, head__20651_101, arg__20655_103, arg__20656_104, v105
					tem__G__0_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{for__s_7})
					if callErr != nil {
						return nil, callErr
					}
					if vm.IsTruthy(tem__G__0_5) {
						f_6 = f_5
						for__a18463_7 = for__a18463_19
						for__s_8 = for__s_7
						header_9 = header_15
						tem__G__0_10 = tem__G__0_5
						goto b1
					} else {
						f_11 = f_5
						for__a18463_12 = for__a18463_19
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
						for__a18463_20 = for__a18463_7
						for__s_21 = for__s_8
						header_22 = header_9
						tem__G__0_23 = tem__G__0_10
						for__xs_24 = tem__G__0_10
						pred_25 = pred_18
						goto b4
					} else {
						f_26 = f_6
						for__a18463_27 = for__a18463_7
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
					for__a18463_111 = for__a18463_12
					for__s_112 = for__s_13
					header_113 = header_14
					tem__G__0_114 = tem__G__0_15
					goto b3
				b3:
					;
					return v109, nil
				b4:
					;
					arg__20606_37, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{pred_25, header_22})
					if callErr != nil {
						return nil, callErr
					}
					arg__20611_40, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{pred_25, header_22})
					if callErr != nil {
						return nil, callErr
					}
					v41, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "list").Deref(), []vm.Value{arg__20611_40})
					if callErr != nil {
						return nil, callErr
					}
					arg__20612_45 = v41
					f_46 = f_19
					for__a18463_47 = for__a18463_20
					for__s_48 = for__s_21
					header_49 = header_22
					tem__G__0_50 = tem__G__0_23
					for__xs_51 = for__xs_24
					pred_52 = pred_25
					goto b6
				b5:
					;
					arg__20612_45 = vm.NIL
					f_46 = f_26
					for__a18463_47 = for__a18463_27
					for__s_48 = for__s_28
					header_49 = header_29
					tem__G__0_50 = tem__G__0_30
					for__xs_51 = for__xs_31
					pred_52 = pred_32
					goto b6
				b6:
					;
					arg__20616_54, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{for__xs_51})
					if callErr != nil {
						return nil, callErr
					}
					head__20620_56, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{for__a18463_47})
					if callErr != nil {
						return nil, callErr
					}
					arg__20624_58, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{for__xs_51})
					if callErr != nil {
						return nil, callErr
					}
					arg__20625_59, callErr = rt.InvokeValue(head__20620_56, []vm.Value{arg__20624_58})
					if callErr != nil {
						return nil, callErr
					}
					v78, callErr = rt.InvokeValue(rt.LookupVar("ir.dominance", "dominates?").Deref(), []vm.Value{f_46, header_49, pred_52})
					if callErr != nil {
						return nil, callErr
					}
					if vm.IsTruthy(v78) {
						f_61 = f_46
						for__a18463_62 = for__a18463_47
						for__s_63 = for__s_48
						header_64 = header_49
						tem__G__0_65 = tem__G__0_50
						for__xs_66 = for__xs_51
						pred_67 = pred_52
						head__20626_68 = rt.LookupVar("clojure.core", "concat-list").Deref()
						goto b7
					} else {
						f_69 = f_46
						for__a18463_70 = for__a18463_47
						for__s_71 = for__s_48
						header_72 = header_49
						tem__G__0_73 = tem__G__0_50
						for__xs_74 = for__xs_51
						pred_75 = pred_52
						head__20626_76 = rt.LookupVar("clojure.core", "concat-list").Deref()
						goto b8
					}
				b7:
					;
					arg__20637_81, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{pred_67, header_64})
					if callErr != nil {
						return nil, callErr
					}
					arg__20642_84, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{pred_67, header_64})
					if callErr != nil {
						return nil, callErr
					}
					v85, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "list").Deref(), []vm.Value{arg__20642_84})
					if callErr != nil {
						return nil, callErr
					}
					arg__20643_89 = v85
					f_90 = f_61
					for__a18463_91 = for__a18463_62
					for__s_92 = for__s_63
					header_93 = header_64
					tem__G__0_94 = tem__G__0_65
					for__xs_95 = for__xs_66
					pred_96 = pred_67
					head__20626_97 = head__20626_68
					goto b9
				b8:
					;
					arg__20643_89 = vm.NIL
					f_90 = f_69
					for__a18463_91 = for__a18463_70
					for__s_92 = for__s_71
					header_93 = header_72
					tem__G__0_94 = tem__G__0_73
					for__xs_95 = for__xs_74
					pred_96 = pred_75
					head__20626_97 = head__20626_76
					goto b9
				b9:
					;
					arg__20647_99, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{for__xs_95})
					if callErr != nil {
						return nil, callErr
					}
					head__20651_101, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{for__a18463_91})
					if callErr != nil {
						return nil, callErr
					}
					arg__20655_103, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{for__xs_95})
					if callErr != nil {
						return nil, callErr
					}
					arg__20656_104, callErr = rt.InvokeValue(head__20651_101, []vm.Value{arg__20655_103})
					if callErr != nil {
						return nil, callErr
					}
					v105, callErr = rt.InvokeValue(head__20626_97, []vm.Value{arg__20643_89, arg__20656_104})
					if callErr != nil {
						return nil, callErr
					}
					v109 = v105
					f_110 = f_90
					for__a18463_111 = for__a18463_91
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
			for__iter18462_35, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{for__a18463_19})
			if callErr != nil {
				return nil, callErr
			}
			arg__20667_37, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-preds").Deref(), []vm.Value{header_15, f_5})
			if callErr != nil {
				return nil, callErr
			}
			arg__20674_39, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-preds").Deref(), []vm.Value{header_15, f_5})
			if callErr != nil {
				return nil, callErr
			}
			arg__20675_40, callErr = rt.InvokeValue(for__iter18462_35, []vm.Value{arg__20674_39})
			if callErr != nil {
				return nil, callErr
			}
			arg__20679_42, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{tem__G__0_8})
			if callErr != nil {
				return nil, callErr
			}
			head__20683_44, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{for__a18465_6})
			if callErr != nil {
				return nil, callErr
			}
			arg__20687_46, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{tem__G__0_8})
			if callErr != nil {
				return nil, callErr
			}
			arg__20688_47, callErr = rt.InvokeValue(head__20683_44, []vm.Value{arg__20687_46})
			if callErr != nil {
				return nil, callErr
			}
			for__a18463_52, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "atom").Deref(), []vm.Value{vm.NIL})
			if callErr != nil {
				return nil, callErr
			}
			__66, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "reset!").Deref(), []vm.Value{for__a18463_52, rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
				var v18 vm.Value
				var callErr error
				_ = v18
				v18, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "lazy-seq*").Deref(), []vm.Value{rt.BoxNativeFn(func() (vm.Value, error) {
					var tem__G__0_5 vm.Value
					var f_6 vm.Value
					var for__a18463_7 vm.Value
					var for__s_8 vm.Value
					var header_9 vm.Value
					var tem__G__0_10 vm.Value
					var pred_18 vm.Value
					var v34 vm.Value
					var f_11 vm.Value
					var for__a18463_12 vm.Value
					var for__s_13 vm.Value
					var header_14 vm.Value
					var tem__G__0_15 vm.Value
					var v109 vm.Value
					var f_110 vm.Value
					var for__a18463_111 vm.Value
					var for__s_112 vm.Value
					var header_113 vm.Value
					var tem__G__0_114 vm.Value
					var f_19 vm.Value
					var for__a18463_20 vm.Value
					var for__s_21 vm.Value
					var header_22 vm.Value
					var tem__G__0_23 vm.Value
					var for__xs_24 vm.Value
					var pred_25 vm.Value
					var arg__20919_37 vm.Value
					var arg__20924_40 vm.Value
					var v41 vm.Value
					var f_26 vm.Value
					var for__a18463_27 vm.Value
					var for__s_28 vm.Value
					var header_29 vm.Value
					var tem__G__0_30 vm.Value
					var for__xs_31 vm.Value
					var pred_32 vm.Value
					var arg__20925_45 vm.Value
					var f_46 vm.Value
					var for__a18463_47 vm.Value
					var for__s_48 vm.Value
					var header_49 vm.Value
					var tem__G__0_50 vm.Value
					var for__xs_51 vm.Value
					var pred_52 vm.Value
					var arg__20929_54 vm.Value
					var head__20933_56 vm.Value
					var arg__20937_58 vm.Value
					var arg__20938_59 vm.Value
					var v78 vm.Value
					var f_61 vm.Value
					var for__a18463_62 vm.Value
					var for__s_63 vm.Value
					var header_64 vm.Value
					var tem__G__0_65 vm.Value
					var for__xs_66 vm.Value
					var pred_67 vm.Value
					var head__20939_68 vm.Value
					var arg__20950_81 vm.Value
					var arg__20955_84 vm.Value
					var v85 vm.Value
					var f_69 vm.Value
					var for__a18463_70 vm.Value
					var for__s_71 vm.Value
					var header_72 vm.Value
					var tem__G__0_73 vm.Value
					var for__xs_74 vm.Value
					var pred_75 vm.Value
					var head__20939_76 vm.Value
					var arg__20956_89 vm.Value
					var f_90 vm.Value
					var for__a18463_91 vm.Value
					var for__s_92 vm.Value
					var header_93 vm.Value
					var tem__G__0_94 vm.Value
					var for__xs_95 vm.Value
					var pred_96 vm.Value
					var head__20939_97 vm.Value
					var arg__20960_99 vm.Value
					var head__20964_101 vm.Value
					var arg__20968_103 vm.Value
					var arg__20969_104 vm.Value
					var v105 vm.Value
					var callErr error
					_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = tem__G__0_5, f_6, for__a18463_7, for__s_8, header_9, tem__G__0_10, pred_18, v34, f_11, for__a18463_12, for__s_13, header_14, tem__G__0_15, v109, f_110, for__a18463_111, for__s_112, header_113, tem__G__0_114, f_19, for__a18463_20, for__s_21, header_22, tem__G__0_23, for__xs_24, pred_25, arg__20919_37, arg__20924_40, v41, f_26, for__a18463_27, for__s_28, header_29, tem__G__0_30, for__xs_31, pred_32, arg__20925_45, f_46, for__a18463_47, for__s_48, header_49, tem__G__0_50, for__xs_51, pred_52, arg__20929_54, head__20933_56, arg__20937_58, arg__20938_59, v78, f_61, for__a18463_62, for__s_63, header_64, tem__G__0_65, for__xs_66, pred_67, head__20939_68, arg__20950_81, arg__20955_84, v85, f_69, for__a18463_70, for__s_71, header_72, tem__G__0_73, for__xs_74, pred_75, head__20939_76, arg__20956_89, f_90, for__a18463_91, for__s_92, header_93, tem__G__0_94, for__xs_95, pred_96, head__20939_97, arg__20960_99, head__20964_101, arg__20968_103, arg__20969_104, v105
					tem__G__0_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{for__s_7})
					if callErr != nil {
						return nil, callErr
					}
					if vm.IsTruthy(tem__G__0_5) {
						f_6 = f_5
						for__a18463_7 = for__a18463_52
						for__s_8 = for__s_7
						header_9 = header_15
						tem__G__0_10 = tem__G__0_5
						goto b1
					} else {
						f_11 = f_5
						for__a18463_12 = for__a18463_52
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
						for__a18463_20 = for__a18463_7
						for__s_21 = for__s_8
						header_22 = header_9
						tem__G__0_23 = tem__G__0_10
						for__xs_24 = tem__G__0_10
						pred_25 = pred_18
						goto b4
					} else {
						f_26 = f_6
						for__a18463_27 = for__a18463_7
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
					for__a18463_111 = for__a18463_12
					for__s_112 = for__s_13
					header_113 = header_14
					tem__G__0_114 = tem__G__0_15
					goto b3
				b3:
					;
					return v109, nil
				b4:
					;
					arg__20919_37, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{pred_25, header_22})
					if callErr != nil {
						return nil, callErr
					}
					arg__20924_40, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{pred_25, header_22})
					if callErr != nil {
						return nil, callErr
					}
					v41, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "list").Deref(), []vm.Value{arg__20924_40})
					if callErr != nil {
						return nil, callErr
					}
					arg__20925_45 = v41
					f_46 = f_19
					for__a18463_47 = for__a18463_20
					for__s_48 = for__s_21
					header_49 = header_22
					tem__G__0_50 = tem__G__0_23
					for__xs_51 = for__xs_24
					pred_52 = pred_25
					goto b6
				b5:
					;
					arg__20925_45 = vm.NIL
					f_46 = f_26
					for__a18463_47 = for__a18463_27
					for__s_48 = for__s_28
					header_49 = header_29
					tem__G__0_50 = tem__G__0_30
					for__xs_51 = for__xs_31
					pred_52 = pred_32
					goto b6
				b6:
					;
					arg__20929_54, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{for__xs_51})
					if callErr != nil {
						return nil, callErr
					}
					head__20933_56, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{for__a18463_47})
					if callErr != nil {
						return nil, callErr
					}
					arg__20937_58, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{for__xs_51})
					if callErr != nil {
						return nil, callErr
					}
					arg__20938_59, callErr = rt.InvokeValue(head__20933_56, []vm.Value{arg__20937_58})
					if callErr != nil {
						return nil, callErr
					}
					v78, callErr = rt.InvokeValue(rt.LookupVar("ir.dominance", "dominates?").Deref(), []vm.Value{f_46, header_49, pred_52})
					if callErr != nil {
						return nil, callErr
					}
					if vm.IsTruthy(v78) {
						f_61 = f_46
						for__a18463_62 = for__a18463_47
						for__s_63 = for__s_48
						header_64 = header_49
						tem__G__0_65 = tem__G__0_50
						for__xs_66 = for__xs_51
						pred_67 = pred_52
						head__20939_68 = rt.LookupVar("clojure.core", "concat-list").Deref()
						goto b7
					} else {
						f_69 = f_46
						for__a18463_70 = for__a18463_47
						for__s_71 = for__s_48
						header_72 = header_49
						tem__G__0_73 = tem__G__0_50
						for__xs_74 = for__xs_51
						pred_75 = pred_52
						head__20939_76 = rt.LookupVar("clojure.core", "concat-list").Deref()
						goto b8
					}
				b7:
					;
					arg__20950_81, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{pred_67, header_64})
					if callErr != nil {
						return nil, callErr
					}
					arg__20955_84, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{pred_67, header_64})
					if callErr != nil {
						return nil, callErr
					}
					v85, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "list").Deref(), []vm.Value{arg__20955_84})
					if callErr != nil {
						return nil, callErr
					}
					arg__20956_89 = v85
					f_90 = f_61
					for__a18463_91 = for__a18463_62
					for__s_92 = for__s_63
					header_93 = header_64
					tem__G__0_94 = tem__G__0_65
					for__xs_95 = for__xs_66
					pred_96 = pred_67
					head__20939_97 = head__20939_68
					goto b9
				b8:
					;
					arg__20956_89 = vm.NIL
					f_90 = f_69
					for__a18463_91 = for__a18463_70
					for__s_92 = for__s_71
					header_93 = header_72
					tem__G__0_94 = tem__G__0_73
					for__xs_95 = for__xs_74
					pred_96 = pred_75
					head__20939_97 = head__20939_76
					goto b9
				b9:
					;
					arg__20960_99, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{for__xs_95})
					if callErr != nil {
						return nil, callErr
					}
					head__20964_101, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{for__a18463_91})
					if callErr != nil {
						return nil, callErr
					}
					arg__20968_103, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{for__xs_95})
					if callErr != nil {
						return nil, callErr
					}
					arg__20969_104, callErr = rt.InvokeValue(head__20964_101, []vm.Value{arg__20968_103})
					if callErr != nil {
						return nil, callErr
					}
					v105, callErr = rt.InvokeValue(head__20939_97, []vm.Value{arg__20956_89, arg__20969_104})
					if callErr != nil {
						return nil, callErr
					}
					v109 = v105
					f_110 = f_90
					for__a18463_111 = for__a18463_91
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
			for__iter18462_68, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{for__a18463_52})
			if callErr != nil {
				return nil, callErr
			}
			arg__20980_70, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-preds").Deref(), []vm.Value{header_15, f_5})
			if callErr != nil {
				return nil, callErr
			}
			arg__20987_72, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-preds").Deref(), []vm.Value{header_15, f_5})
			if callErr != nil {
				return nil, callErr
			}
			arg__20988_73, callErr = rt.InvokeValue(for__iter18462_68, []vm.Value{arg__20987_72})
			if callErr != nil {
				return nil, callErr
			}
			arg__20992_75, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{tem__G__0_8})
			if callErr != nil {
				return nil, callErr
			}
			head__20996_77, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{for__a18465_6})
			if callErr != nil {
				return nil, callErr
			}
			arg__21000_79, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{tem__G__0_8})
			if callErr != nil {
				return nil, callErr
			}
			arg__21001_80, callErr = rt.InvokeValue(head__20996_77, []vm.Value{arg__21000_79})
			if callErr != nil {
				return nil, callErr
			}
			v81, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "concat-list").Deref(), []vm.Value{arg__20988_73, arg__21001_80})
			if callErr != nil {
				return nil, callErr
			}
			v85 = v81
			f_86 = f_5
			for__a18465_87 = for__a18465_6
			for__s_88 = for__s_7
			tem__G__0_89 = tem__G__0_8
			goto b3
		b2:
			;
			v85 = vm.NIL
			f_86 = f_9
			for__a18465_87 = for__a18465_10
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
	for__iter18464_17, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{for__a18465_5})
	if callErr != nil {
		return nil, callErr
	}
	arg__21010_19, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__21015_21, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__21016_22, callErr = rt.InvokeValue(for__iter18464_17, []vm.Value{arg__21015_21})
	if callErr != nil {
		return nil, callErr
	}
	for__a18465_27, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "atom").Deref(), []vm.Value{vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	__37, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "reset!").Deref(), []vm.Value{for__a18465_27, rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var v14 vm.Value
		var callErr error
		_ = v14
		v14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "lazy-seq*").Deref(), []vm.Value{rt.BoxNativeFn(func() (vm.Value, error) {
			var tem__G__0_4 vm.Value
			var f_5 vm.Value
			var for__a18465_6 vm.Value
			var for__s_7 vm.Value
			var tem__G__0_8 vm.Value
			var header_15 vm.Value
			var for__a18463_19 vm.Value
			var __33 vm.Value
			var for__iter18462_35 vm.Value
			var arg__23219_37 vm.Value
			var arg__23226_39 vm.Value
			var arg__23227_40 vm.Value
			var arg__23231_42 vm.Value
			var head__23235_44 vm.Value
			var arg__23239_46 vm.Value
			var arg__23240_47 vm.Value
			var for__a18463_52 vm.Value
			var __66 vm.Value
			var for__iter18462_68 vm.Value
			var arg__23532_70 vm.Value
			var arg__23539_72 vm.Value
			var arg__23540_73 vm.Value
			var arg__23544_75 vm.Value
			var head__23548_77 vm.Value
			var arg__23552_79 vm.Value
			var arg__23553_80 vm.Value
			var v81 vm.Value
			var f_9 vm.Value
			var for__a18465_10 vm.Value
			var for__s_11 vm.Value
			var tem__G__0_12 vm.Value
			var v85 vm.Value
			var f_86 vm.Value
			var for__a18465_87 vm.Value
			var for__s_88 vm.Value
			var tem__G__0_89 vm.Value
			var callErr error
			_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = tem__G__0_4, f_5, for__a18465_6, for__s_7, tem__G__0_8, header_15, for__a18463_19, __33, for__iter18462_35, arg__23219_37, arg__23226_39, arg__23227_40, arg__23231_42, head__23235_44, arg__23239_46, arg__23240_47, for__a18463_52, __66, for__iter18462_68, arg__23532_70, arg__23539_72, arg__23540_73, arg__23544_75, head__23548_77, arg__23552_79, arg__23553_80, v81, f_9, for__a18465_10, for__s_11, tem__G__0_12, v85, f_86, for__a18465_87, for__s_88, tem__G__0_89
			tem__G__0_4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{arg0})
			if callErr != nil {
				return nil, callErr
			}
			if vm.IsTruthy(tem__G__0_4) {
				f_5 = arg0
				for__a18465_6 = for__a18465_27
				for__s_7 = arg0
				tem__G__0_8 = tem__G__0_4
				goto b1
			} else {
				f_9 = arg0
				for__a18465_10 = for__a18465_27
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
			for__a18463_19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "atom").Deref(), []vm.Value{vm.NIL})
			if callErr != nil {
				return nil, callErr
			}
			__33, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "reset!").Deref(), []vm.Value{for__a18463_19, rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
				var v18 vm.Value
				var callErr error
				_ = v18
				v18, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "lazy-seq*").Deref(), []vm.Value{rt.BoxNativeFn(func() (vm.Value, error) {
					var tem__G__0_5 vm.Value
					var f_6 vm.Value
					var for__a18463_7 vm.Value
					var for__s_8 vm.Value
					var header_9 vm.Value
					var tem__G__0_10 vm.Value
					var pred_18 vm.Value
					var v34 vm.Value
					var f_11 vm.Value
					var for__a18463_12 vm.Value
					var for__s_13 vm.Value
					var header_14 vm.Value
					var tem__G__0_15 vm.Value
					var v109 vm.Value
					var f_110 vm.Value
					var for__a18463_111 vm.Value
					var for__s_112 vm.Value
					var header_113 vm.Value
					var tem__G__0_114 vm.Value
					var f_19 vm.Value
					var for__a18463_20 vm.Value
					var for__s_21 vm.Value
					var header_22 vm.Value
					var tem__G__0_23 vm.Value
					var for__xs_24 vm.Value
					var pred_25 vm.Value
					var arg__23158_37 vm.Value
					var arg__23163_40 vm.Value
					var v41 vm.Value
					var f_26 vm.Value
					var for__a18463_27 vm.Value
					var for__s_28 vm.Value
					var header_29 vm.Value
					var tem__G__0_30 vm.Value
					var for__xs_31 vm.Value
					var pred_32 vm.Value
					var arg__23164_45 vm.Value
					var f_46 vm.Value
					var for__a18463_47 vm.Value
					var for__s_48 vm.Value
					var header_49 vm.Value
					var tem__G__0_50 vm.Value
					var for__xs_51 vm.Value
					var pred_52 vm.Value
					var arg__23168_54 vm.Value
					var head__23172_56 vm.Value
					var arg__23176_58 vm.Value
					var arg__23177_59 vm.Value
					var v78 vm.Value
					var f_61 vm.Value
					var for__a18463_62 vm.Value
					var for__s_63 vm.Value
					var header_64 vm.Value
					var tem__G__0_65 vm.Value
					var for__xs_66 vm.Value
					var pred_67 vm.Value
					var head__23178_68 vm.Value
					var arg__23189_81 vm.Value
					var arg__23194_84 vm.Value
					var v85 vm.Value
					var f_69 vm.Value
					var for__a18463_70 vm.Value
					var for__s_71 vm.Value
					var header_72 vm.Value
					var tem__G__0_73 vm.Value
					var for__xs_74 vm.Value
					var pred_75 vm.Value
					var head__23178_76 vm.Value
					var arg__23195_89 vm.Value
					var f_90 vm.Value
					var for__a18463_91 vm.Value
					var for__s_92 vm.Value
					var header_93 vm.Value
					var tem__G__0_94 vm.Value
					var for__xs_95 vm.Value
					var pred_96 vm.Value
					var head__23178_97 vm.Value
					var arg__23199_99 vm.Value
					var head__23203_101 vm.Value
					var arg__23207_103 vm.Value
					var arg__23208_104 vm.Value
					var v105 vm.Value
					var callErr error
					_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = tem__G__0_5, f_6, for__a18463_7, for__s_8, header_9, tem__G__0_10, pred_18, v34, f_11, for__a18463_12, for__s_13, header_14, tem__G__0_15, v109, f_110, for__a18463_111, for__s_112, header_113, tem__G__0_114, f_19, for__a18463_20, for__s_21, header_22, tem__G__0_23, for__xs_24, pred_25, arg__23158_37, arg__23163_40, v41, f_26, for__a18463_27, for__s_28, header_29, tem__G__0_30, for__xs_31, pred_32, arg__23164_45, f_46, for__a18463_47, for__s_48, header_49, tem__G__0_50, for__xs_51, pred_52, arg__23168_54, head__23172_56, arg__23176_58, arg__23177_59, v78, f_61, for__a18463_62, for__s_63, header_64, tem__G__0_65, for__xs_66, pred_67, head__23178_68, arg__23189_81, arg__23194_84, v85, f_69, for__a18463_70, for__s_71, header_72, tem__G__0_73, for__xs_74, pred_75, head__23178_76, arg__23195_89, f_90, for__a18463_91, for__s_92, header_93, tem__G__0_94, for__xs_95, pred_96, head__23178_97, arg__23199_99, head__23203_101, arg__23207_103, arg__23208_104, v105
					tem__G__0_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{for__s_7})
					if callErr != nil {
						return nil, callErr
					}
					if vm.IsTruthy(tem__G__0_5) {
						f_6 = f_5
						for__a18463_7 = for__a18463_19
						for__s_8 = for__s_7
						header_9 = header_15
						tem__G__0_10 = tem__G__0_5
						goto b1
					} else {
						f_11 = f_5
						for__a18463_12 = for__a18463_19
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
						for__a18463_20 = for__a18463_7
						for__s_21 = for__s_8
						header_22 = header_9
						tem__G__0_23 = tem__G__0_10
						for__xs_24 = tem__G__0_10
						pred_25 = pred_18
						goto b4
					} else {
						f_26 = f_6
						for__a18463_27 = for__a18463_7
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
					for__a18463_111 = for__a18463_12
					for__s_112 = for__s_13
					header_113 = header_14
					tem__G__0_114 = tem__G__0_15
					goto b3
				b3:
					;
					return v109, nil
				b4:
					;
					arg__23158_37, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{pred_25, header_22})
					if callErr != nil {
						return nil, callErr
					}
					arg__23163_40, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{pred_25, header_22})
					if callErr != nil {
						return nil, callErr
					}
					v41, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "list").Deref(), []vm.Value{arg__23163_40})
					if callErr != nil {
						return nil, callErr
					}
					arg__23164_45 = v41
					f_46 = f_19
					for__a18463_47 = for__a18463_20
					for__s_48 = for__s_21
					header_49 = header_22
					tem__G__0_50 = tem__G__0_23
					for__xs_51 = for__xs_24
					pred_52 = pred_25
					goto b6
				b5:
					;
					arg__23164_45 = vm.NIL
					f_46 = f_26
					for__a18463_47 = for__a18463_27
					for__s_48 = for__s_28
					header_49 = header_29
					tem__G__0_50 = tem__G__0_30
					for__xs_51 = for__xs_31
					pred_52 = pred_32
					goto b6
				b6:
					;
					arg__23168_54, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{for__xs_51})
					if callErr != nil {
						return nil, callErr
					}
					head__23172_56, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{for__a18463_47})
					if callErr != nil {
						return nil, callErr
					}
					arg__23176_58, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{for__xs_51})
					if callErr != nil {
						return nil, callErr
					}
					arg__23177_59, callErr = rt.InvokeValue(head__23172_56, []vm.Value{arg__23176_58})
					if callErr != nil {
						return nil, callErr
					}
					v78, callErr = rt.InvokeValue(rt.LookupVar("ir.dominance", "dominates?").Deref(), []vm.Value{f_46, header_49, pred_52})
					if callErr != nil {
						return nil, callErr
					}
					if vm.IsTruthy(v78) {
						f_61 = f_46
						for__a18463_62 = for__a18463_47
						for__s_63 = for__s_48
						header_64 = header_49
						tem__G__0_65 = tem__G__0_50
						for__xs_66 = for__xs_51
						pred_67 = pred_52
						head__23178_68 = rt.LookupVar("clojure.core", "concat-list").Deref()
						goto b7
					} else {
						f_69 = f_46
						for__a18463_70 = for__a18463_47
						for__s_71 = for__s_48
						header_72 = header_49
						tem__G__0_73 = tem__G__0_50
						for__xs_74 = for__xs_51
						pred_75 = pred_52
						head__23178_76 = rt.LookupVar("clojure.core", "concat-list").Deref()
						goto b8
					}
				b7:
					;
					arg__23189_81, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{pred_67, header_64})
					if callErr != nil {
						return nil, callErr
					}
					arg__23194_84, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{pred_67, header_64})
					if callErr != nil {
						return nil, callErr
					}
					v85, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "list").Deref(), []vm.Value{arg__23194_84})
					if callErr != nil {
						return nil, callErr
					}
					arg__23195_89 = v85
					f_90 = f_61
					for__a18463_91 = for__a18463_62
					for__s_92 = for__s_63
					header_93 = header_64
					tem__G__0_94 = tem__G__0_65
					for__xs_95 = for__xs_66
					pred_96 = pred_67
					head__23178_97 = head__23178_68
					goto b9
				b8:
					;
					arg__23195_89 = vm.NIL
					f_90 = f_69
					for__a18463_91 = for__a18463_70
					for__s_92 = for__s_71
					header_93 = header_72
					tem__G__0_94 = tem__G__0_73
					for__xs_95 = for__xs_74
					pred_96 = pred_75
					head__23178_97 = head__23178_76
					goto b9
				b9:
					;
					arg__23199_99, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{for__xs_95})
					if callErr != nil {
						return nil, callErr
					}
					head__23203_101, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{for__a18463_91})
					if callErr != nil {
						return nil, callErr
					}
					arg__23207_103, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{for__xs_95})
					if callErr != nil {
						return nil, callErr
					}
					arg__23208_104, callErr = rt.InvokeValue(head__23203_101, []vm.Value{arg__23207_103})
					if callErr != nil {
						return nil, callErr
					}
					v105, callErr = rt.InvokeValue(head__23178_97, []vm.Value{arg__23195_89, arg__23208_104})
					if callErr != nil {
						return nil, callErr
					}
					v109 = v105
					f_110 = f_90
					for__a18463_111 = for__a18463_91
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
			for__iter18462_35, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{for__a18463_19})
			if callErr != nil {
				return nil, callErr
			}
			arg__23219_37, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-preds").Deref(), []vm.Value{header_15, f_5})
			if callErr != nil {
				return nil, callErr
			}
			arg__23226_39, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-preds").Deref(), []vm.Value{header_15, f_5})
			if callErr != nil {
				return nil, callErr
			}
			arg__23227_40, callErr = rt.InvokeValue(for__iter18462_35, []vm.Value{arg__23226_39})
			if callErr != nil {
				return nil, callErr
			}
			arg__23231_42, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{tem__G__0_8})
			if callErr != nil {
				return nil, callErr
			}
			head__23235_44, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{for__a18465_6})
			if callErr != nil {
				return nil, callErr
			}
			arg__23239_46, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{tem__G__0_8})
			if callErr != nil {
				return nil, callErr
			}
			arg__23240_47, callErr = rt.InvokeValue(head__23235_44, []vm.Value{arg__23239_46})
			if callErr != nil {
				return nil, callErr
			}
			for__a18463_52, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "atom").Deref(), []vm.Value{vm.NIL})
			if callErr != nil {
				return nil, callErr
			}
			__66, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "reset!").Deref(), []vm.Value{for__a18463_52, rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
				var v18 vm.Value
				var callErr error
				_ = v18
				v18, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "lazy-seq*").Deref(), []vm.Value{rt.BoxNativeFn(func() (vm.Value, error) {
					var tem__G__0_5 vm.Value
					var f_6 vm.Value
					var for__a18463_7 vm.Value
					var for__s_8 vm.Value
					var header_9 vm.Value
					var tem__G__0_10 vm.Value
					var pred_18 vm.Value
					var v34 vm.Value
					var f_11 vm.Value
					var for__a18463_12 vm.Value
					var for__s_13 vm.Value
					var header_14 vm.Value
					var tem__G__0_15 vm.Value
					var v109 vm.Value
					var f_110 vm.Value
					var for__a18463_111 vm.Value
					var for__s_112 vm.Value
					var header_113 vm.Value
					var tem__G__0_114 vm.Value
					var f_19 vm.Value
					var for__a18463_20 vm.Value
					var for__s_21 vm.Value
					var header_22 vm.Value
					var tem__G__0_23 vm.Value
					var for__xs_24 vm.Value
					var pred_25 vm.Value
					var arg__23471_37 vm.Value
					var arg__23476_40 vm.Value
					var v41 vm.Value
					var f_26 vm.Value
					var for__a18463_27 vm.Value
					var for__s_28 vm.Value
					var header_29 vm.Value
					var tem__G__0_30 vm.Value
					var for__xs_31 vm.Value
					var pred_32 vm.Value
					var arg__23477_45 vm.Value
					var f_46 vm.Value
					var for__a18463_47 vm.Value
					var for__s_48 vm.Value
					var header_49 vm.Value
					var tem__G__0_50 vm.Value
					var for__xs_51 vm.Value
					var pred_52 vm.Value
					var arg__23481_54 vm.Value
					var head__23485_56 vm.Value
					var arg__23489_58 vm.Value
					var arg__23490_59 vm.Value
					var v78 vm.Value
					var f_61 vm.Value
					var for__a18463_62 vm.Value
					var for__s_63 vm.Value
					var header_64 vm.Value
					var tem__G__0_65 vm.Value
					var for__xs_66 vm.Value
					var pred_67 vm.Value
					var head__23491_68 vm.Value
					var arg__23502_81 vm.Value
					var arg__23507_84 vm.Value
					var v85 vm.Value
					var f_69 vm.Value
					var for__a18463_70 vm.Value
					var for__s_71 vm.Value
					var header_72 vm.Value
					var tem__G__0_73 vm.Value
					var for__xs_74 vm.Value
					var pred_75 vm.Value
					var head__23491_76 vm.Value
					var arg__23508_89 vm.Value
					var f_90 vm.Value
					var for__a18463_91 vm.Value
					var for__s_92 vm.Value
					var header_93 vm.Value
					var tem__G__0_94 vm.Value
					var for__xs_95 vm.Value
					var pred_96 vm.Value
					var head__23491_97 vm.Value
					var arg__23512_99 vm.Value
					var head__23516_101 vm.Value
					var arg__23520_103 vm.Value
					var arg__23521_104 vm.Value
					var v105 vm.Value
					var callErr error
					_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = tem__G__0_5, f_6, for__a18463_7, for__s_8, header_9, tem__G__0_10, pred_18, v34, f_11, for__a18463_12, for__s_13, header_14, tem__G__0_15, v109, f_110, for__a18463_111, for__s_112, header_113, tem__G__0_114, f_19, for__a18463_20, for__s_21, header_22, tem__G__0_23, for__xs_24, pred_25, arg__23471_37, arg__23476_40, v41, f_26, for__a18463_27, for__s_28, header_29, tem__G__0_30, for__xs_31, pred_32, arg__23477_45, f_46, for__a18463_47, for__s_48, header_49, tem__G__0_50, for__xs_51, pred_52, arg__23481_54, head__23485_56, arg__23489_58, arg__23490_59, v78, f_61, for__a18463_62, for__s_63, header_64, tem__G__0_65, for__xs_66, pred_67, head__23491_68, arg__23502_81, arg__23507_84, v85, f_69, for__a18463_70, for__s_71, header_72, tem__G__0_73, for__xs_74, pred_75, head__23491_76, arg__23508_89, f_90, for__a18463_91, for__s_92, header_93, tem__G__0_94, for__xs_95, pred_96, head__23491_97, arg__23512_99, head__23516_101, arg__23520_103, arg__23521_104, v105
					tem__G__0_5, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{for__s_7})
					if callErr != nil {
						return nil, callErr
					}
					if vm.IsTruthy(tem__G__0_5) {
						f_6 = f_5
						for__a18463_7 = for__a18463_52
						for__s_8 = for__s_7
						header_9 = header_15
						tem__G__0_10 = tem__G__0_5
						goto b1
					} else {
						f_11 = f_5
						for__a18463_12 = for__a18463_52
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
						for__a18463_20 = for__a18463_7
						for__s_21 = for__s_8
						header_22 = header_9
						tem__G__0_23 = tem__G__0_10
						for__xs_24 = tem__G__0_10
						pred_25 = pred_18
						goto b4
					} else {
						f_26 = f_6
						for__a18463_27 = for__a18463_7
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
					for__a18463_111 = for__a18463_12
					for__s_112 = for__s_13
					header_113 = header_14
					tem__G__0_114 = tem__G__0_15
					goto b3
				b3:
					;
					return v109, nil
				b4:
					;
					arg__23471_37, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{pred_25, header_22})
					if callErr != nil {
						return nil, callErr
					}
					arg__23476_40, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{pred_25, header_22})
					if callErr != nil {
						return nil, callErr
					}
					v41, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "list").Deref(), []vm.Value{arg__23476_40})
					if callErr != nil {
						return nil, callErr
					}
					arg__23477_45 = v41
					f_46 = f_19
					for__a18463_47 = for__a18463_20
					for__s_48 = for__s_21
					header_49 = header_22
					tem__G__0_50 = tem__G__0_23
					for__xs_51 = for__xs_24
					pred_52 = pred_25
					goto b6
				b5:
					;
					arg__23477_45 = vm.NIL
					f_46 = f_26
					for__a18463_47 = for__a18463_27
					for__s_48 = for__s_28
					header_49 = header_29
					tem__G__0_50 = tem__G__0_30
					for__xs_51 = for__xs_31
					pred_52 = pred_32
					goto b6
				b6:
					;
					arg__23481_54, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{for__xs_51})
					if callErr != nil {
						return nil, callErr
					}
					head__23485_56, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{for__a18463_47})
					if callErr != nil {
						return nil, callErr
					}
					arg__23489_58, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{for__xs_51})
					if callErr != nil {
						return nil, callErr
					}
					arg__23490_59, callErr = rt.InvokeValue(head__23485_56, []vm.Value{arg__23489_58})
					if callErr != nil {
						return nil, callErr
					}
					v78, callErr = rt.InvokeValue(rt.LookupVar("ir.dominance", "dominates?").Deref(), []vm.Value{f_46, header_49, pred_52})
					if callErr != nil {
						return nil, callErr
					}
					if vm.IsTruthy(v78) {
						f_61 = f_46
						for__a18463_62 = for__a18463_47
						for__s_63 = for__s_48
						header_64 = header_49
						tem__G__0_65 = tem__G__0_50
						for__xs_66 = for__xs_51
						pred_67 = pred_52
						head__23491_68 = rt.LookupVar("clojure.core", "concat-list").Deref()
						goto b7
					} else {
						f_69 = f_46
						for__a18463_70 = for__a18463_47
						for__s_71 = for__s_48
						header_72 = header_49
						tem__G__0_73 = tem__G__0_50
						for__xs_74 = for__xs_51
						pred_75 = pred_52
						head__23491_76 = rt.LookupVar("clojure.core", "concat-list").Deref()
						goto b8
					}
				b7:
					;
					arg__23502_81, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{pred_67, header_64})
					if callErr != nil {
						return nil, callErr
					}
					arg__23507_84, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{pred_67, header_64})
					if callErr != nil {
						return nil, callErr
					}
					v85, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "list").Deref(), []vm.Value{arg__23507_84})
					if callErr != nil {
						return nil, callErr
					}
					arg__23508_89 = v85
					f_90 = f_61
					for__a18463_91 = for__a18463_62
					for__s_92 = for__s_63
					header_93 = header_64
					tem__G__0_94 = tem__G__0_65
					for__xs_95 = for__xs_66
					pred_96 = pred_67
					head__23491_97 = head__23491_68
					goto b9
				b8:
					;
					arg__23508_89 = vm.NIL
					f_90 = f_69
					for__a18463_91 = for__a18463_70
					for__s_92 = for__s_71
					header_93 = header_72
					tem__G__0_94 = tem__G__0_73
					for__xs_95 = for__xs_74
					pred_96 = pred_75
					head__23491_97 = head__23491_76
					goto b9
				b9:
					;
					arg__23512_99, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{for__xs_95})
					if callErr != nil {
						return nil, callErr
					}
					head__23516_101, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{for__a18463_91})
					if callErr != nil {
						return nil, callErr
					}
					arg__23520_103, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{for__xs_95})
					if callErr != nil {
						return nil, callErr
					}
					arg__23521_104, callErr = rt.InvokeValue(head__23516_101, []vm.Value{arg__23520_103})
					if callErr != nil {
						return nil, callErr
					}
					v105, callErr = rt.InvokeValue(head__23491_97, []vm.Value{arg__23508_89, arg__23521_104})
					if callErr != nil {
						return nil, callErr
					}
					v109 = v105
					f_110 = f_90
					for__a18463_111 = for__a18463_91
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
			for__iter18462_68, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{for__a18463_52})
			if callErr != nil {
				return nil, callErr
			}
			arg__23532_70, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-preds").Deref(), []vm.Value{header_15, f_5})
			if callErr != nil {
				return nil, callErr
			}
			arg__23539_72, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-preds").Deref(), []vm.Value{header_15, f_5})
			if callErr != nil {
				return nil, callErr
			}
			arg__23540_73, callErr = rt.InvokeValue(for__iter18462_68, []vm.Value{arg__23539_72})
			if callErr != nil {
				return nil, callErr
			}
			arg__23544_75, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{tem__G__0_8})
			if callErr != nil {
				return nil, callErr
			}
			head__23548_77, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{for__a18465_6})
			if callErr != nil {
				return nil, callErr
			}
			arg__23552_79, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "rest").Deref(), []vm.Value{tem__G__0_8})
			if callErr != nil {
				return nil, callErr
			}
			arg__23553_80, callErr = rt.InvokeValue(head__23548_77, []vm.Value{arg__23552_79})
			if callErr != nil {
				return nil, callErr
			}
			v81, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "concat-list").Deref(), []vm.Value{arg__23540_73, arg__23553_80})
			if callErr != nil {
				return nil, callErr
			}
			v85 = v81
			f_86 = f_5
			for__a18465_87 = for__a18465_6
			for__s_88 = for__s_7
			tem__G__0_89 = tem__G__0_8
			goto b3
		b2:
			;
			v85 = vm.NIL
			f_86 = f_9
			for__a18465_87 = for__a18465_10
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
	for__iter18464_39, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "deref").Deref(), []vm.Value{for__a18465_27})
	if callErr != nil {
		return nil, callErr
	}
	arg__23562_41, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__23567_43, callErr = rt.InvokeValue(rt.LookupVar("ir", "blocks").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__23568_44, callErr = rt.InvokeValue(for__iter18464_39, []vm.Value{arg__23567_43})
	if callErr != nil {
		return nil, callErr
	}
	v45, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg__23568_44})
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
	var arg__23580_5 vm.Value
	var arg__23588_8 vm.Value
	var outs_9 vm.Value
	var arg__23593_20 vm.Value
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
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__23580_5, arg__23588_8, outs_9, arg__23593_20, v21, f_10, header_11, loop_set_12, outs_13, v24, f_14, header_15, loop_set_16, outs_17, v28, f_29, header_30, loop_set_31, outs_32
	arg__23580_5, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-preds").Deref(), []vm.Value{arg1, arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__23588_8, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-preds").Deref(), []vm.Value{arg1, arg0})
	if callErr != nil {
		return nil, callErr
	}
	outs_9, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "remove").Deref(), []vm.Value{arg2, arg__23588_8})
	if callErr != nil {
		return nil, callErr
	}
	arg__23593_20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "count").Deref(), []vm.Value{outs_9})
	if callErr != nil {
		return nil, callErr
	}
	v21 = arg__23593_20 == vm.Int(1)
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
	var arg__24329_16 vm.Value
	var arg__24339_19 vm.Value
	var v20 vm.Value
	var callErr error
	_, _, _, _, _ = pred_7, header_13, arg__24329_16, arg__24339_19, v20
	pred_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg1, vm.Int(0), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	header_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg1, vm.Int(1), vm.NIL})
	if callErr != nil {
		return nil, callErr
	}
	arg__24329_16, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.licm", "reachable-without").Deref(), []vm.Value{arg0, pred_7, header_13})
	if callErr != nil {
		return nil, callErr
	}
	arg__24339_19, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.licm", "reachable-without").Deref(), []vm.Value{arg0, pred_7, header_13})
	if callErr != nil {
		return nil, callErr
	}
	v20, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "conj").Deref(), []vm.Value{arg__24339_19, header_13})
	if callErr != nil {
		return nil, callErr
	}
	return v20, nil
}
func find_in_loop_users(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__24344_5 vm.Value
	var arg__24350_8 vm.Value
	var clone_users_9 vm.Value
	var acc_13 vm.Value
	var v27 vm.Value
	var v29 vm.Value
	var callErr error
	_, _, _, _, _, _ = arg__24344_5, arg__24350_8, clone_users_9, acc_13, v27, v29
	arg__24344_5, callErr = rt.InvokeValue(rt.LookupVar("ir", "uses").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__24350_8, callErr = rt.InvokeValue(rt.LookupVar("ir", "uses").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	clone_users_9, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "nth").Deref(), []vm.Value{arg__24350_8, arg1})
	if callErr != nil {
		return nil, callErr
	}
	acc_13, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "atom").Deref(), []vm.Value{vm.NewArrayVector([]vm.Value{})})
	if callErr != nil {
		return nil, callErr
	}
	v27, callErr = rt.InvokeValue(rt.LookupVar("ir", "uses-for-each").Deref(), []vm.Value{clone_users_9, rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var arg__24419_17 vm.Value
		var arg__24427_21 vm.Value
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
		var arg__24433_37 vm.Value
		var arg__24440_39 vm.Value
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
		var arg__24449_58 vm.Value
		var arg__24459_64 vm.Value
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
		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__24419_17, arg__24427_21, and__x_22, nid_5, acc_6, clone_7, f_8, loop_set_9, v88, nid_10, acc_11, clone_12, f_13, loop_set_14, v92, nid_93, acc_94, clone_95, f_96, loop_set_97, nid_23, acc_24, clone_25, f_26, loop_set_27, and__x_28, arg__24433_37, arg__24440_39, and__x_40, nid_29, acc_30, clone_31, f_32, loop_set_33, and__x_34, v77, nid_78, acc_79, clone_80, f_81, loop_set_82, and__x_83, nid_41, acc_42, clone_43, f_44, loop_set_45, and__x_46, arg__24449_58, arg__24459_64, v65, nid_47, acc_48, clone_49, f_50, loop_set_51, and__x_52, v68, nid_69, acc_70, clone_71, f_72, loop_set_73, and__x_74
		arg__24419_17, callErr = rt.InvokeValue(rt.LookupVar("ir", "op").Deref(), []vm.Value{arg0, arg0})
		if callErr != nil {
			return nil, callErr
		}
		arg__24427_21, callErr = rt.InvokeValue(rt.LookupVar("ir", "op").Deref(), []vm.Value{arg0, arg0})
		if callErr != nil {
			return nil, callErr
		}
		and__x_22, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not=").Deref(), []vm.Value{vm.Keyword("invalid"), arg__24427_21})
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
		arg__24433_37, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-of").Deref(), []vm.Value{nid_23, f_26})
		if callErr != nil {
			return nil, callErr
		}
		arg__24440_39, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-of").Deref(), []vm.Value{nid_23, f_26})
		if callErr != nil {
			return nil, callErr
		}
		and__x_40, callErr = rt.InvokeValue(loop_set_27, []vm.Value{arg__24440_39})
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
		arg__24449_58, callErr = rt.InvokeValue(rt.LookupVar("ir", "refs").Deref(), []vm.Value{nid_41, f_44})
		if callErr != nil {
			return nil, callErr
		}
		arg__24459_64, callErr = rt.InvokeValue(rt.LookupVar("ir", "refs").Deref(), []vm.Value{nid_41, f_44})
		if callErr != nil {
			return nil, callErr
		}
		v65, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "some").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) vm.Value {
			var v2 vm.Value
			_ = v2
			v2 = vm.Boolean(arg0 == clone_43)
			return v2
		}), arg__24459_64})
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
	var arg__25120_29 vm.Value
	var arg__25136_36 vm.Value
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
	var arg__25148_59 vm.Value
	var arg__25159_65 vm.Value
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
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = header_6, loop_set_8, pre_header_10, f_11, back_edge_12, header_13, loop_set_14, pre_header_15, hoistable_23, arg__25120_29, arg__25136_36, hoisted_pairs_37, v53, f_16, back_edge_17, header_18, loop_set_19, pre_header_20, v81, f_82, back_edge_83, header_84, loop_set_85, pre_header_86, f_38, back_edge_39, header_40, loop_set_41, pre_header_42, hoistable_43, hoisted_pairs_44, arg__25148_59, arg__25159_65, v66, f_45, back_edge_46, header_47, loop_set_48, pre_header_49, hoistable_50, hoisted_pairs_51, v70, f_71, back_edge_72, header_73, loop_set_74, pre_header_75, hoistable_76, hoisted_pairs_77
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
	arg__25120_29, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "sort").Deref(), []vm.Value{hoistable_23})
	if callErr != nil {
		return nil, callErr
	}
	arg__25136_36, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "sort").Deref(), []vm.Value{hoistable_23})
	if callErr != nil {
		return nil, callErr
	}
	hoisted_pairs_37, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "mapv").Deref(), []vm.Value{rt.BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) {
		var arg__25131_5 vm.Value
		var v6 vm.Value
		var callErr error
		_, _ = arg__25131_5, v6
		arg__25131_5, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.licm", "hoist-one!").Deref(), []vm.Value{f_11, arg0, pre_header_15})
		if callErr != nil {
			return nil, callErr
		}
		v6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vector").Deref(), []vm.Value{arg0, arg__25131_5})
		if callErr != nil {
			return nil, callErr
		}
		return v6, nil
	}), arg__25136_36})
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
	arg__25148_59, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("body"), loop_set_41, vm.Keyword("preheader"), pre_header_42, vm.Keyword("header"), header_40})
	if callErr != nil {
		return nil, callErr
	}
	arg__25159_65, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "array-map").Deref(), []vm.Value{vm.Keyword("body"), loop_set_41, vm.Keyword("preheader"), pre_header_42, vm.Keyword("header"), header_40})
	if callErr != nil {
		return nil, callErr
	}
	v66, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.licm", "thread-hoisted-through-header!").Deref(), []vm.Value{f_38, arg__25159_65, hoisted_pairs_44})
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
		var arg__25208_4 vm.Value
		var arg__25216_7 vm.Value
		var arg__25217_8 vm.Value
		var arg__25225_11 vm.Value
		var arg__25233_14 vm.Value
		var arg__25234_15 vm.Value
		var v16 vm.Value
		var callErr error
		_, _, _, _, _, _, _ = arg__25208_4, arg__25216_7, arg__25217_8, arg__25225_11, arg__25233_14, arg__25234_15, v16
		arg__25208_4, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-of").Deref(), []vm.Value{arg0, arg0})
		if callErr != nil {
			return nil, callErr
		}
		arg__25216_7, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-of").Deref(), []vm.Value{arg0, arg0})
		if callErr != nil {
			return nil, callErr
		}
		arg__25217_8, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "contains?").Deref(), []vm.Value{arg2, arg__25216_7})
		if callErr != nil {
			return nil, callErr
		}
		arg__25225_11, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-of").Deref(), []vm.Value{arg0, arg0})
		if callErr != nil {
			return nil, callErr
		}
		arg__25233_14, callErr = rt.InvokeValue(rt.LookupVar("ir", "block-of").Deref(), []vm.Value{arg0, arg0})
		if callErr != nil {
			return nil, callErr
		}
		arg__25234_15, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "contains?").Deref(), []vm.Value{arg2, arg__25233_14})
		if callErr != nil {
			return nil, callErr
		}
		v16, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "not").Deref(), []vm.Value{arg__25234_15})
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
	var arg__25242_3 vm.Value
	var arg__25247_6 vm.Value
	var doseq_seq__25237_7 vm.Value
	var doseq_loop__25238_8 vm.Value
	var f_9 vm.Value
	var doseq_seq__25237_11 vm.Value
	var doseq_loop__25238_12 vm.Value
	var f_13 vm.Value
	var be_19 vm.Value
	var v21 vm.Value
	var v23 vm.Value
	var doseq_seq__25237_14 vm.Value
	var doseq_loop__25238_15 vm.Value
	var f_16 vm.Value
	var v27 vm.Value
	var doseq_seq__25237_28 vm.Value
	var doseq_loop__25238_29 vm.Value
	var f_30 vm.Value
	var callErr error
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = arg__25242_3, arg__25247_6, doseq_seq__25237_7, doseq_loop__25238_8, f_9, doseq_seq__25237_11, doseq_loop__25238_12, f_13, be_19, v21, v23, doseq_seq__25237_14, doseq_loop__25238_15, f_16, v27, doseq_seq__25237_28, doseq_loop__25238_29, f_30
	arg__25242_3, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.licm", "back-edges").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__25247_6, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.licm", "back-edges").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	doseq_seq__25237_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "seq").Deref(), []vm.Value{arg__25247_6})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__25238_8 = doseq_seq__25237_7
	f_9 = arg0
	goto b1
b1:
	;
	if vm.IsTruthy(doseq_loop__25238_8) {
		doseq_seq__25237_11 = doseq_seq__25237_7
		doseq_loop__25238_12 = doseq_loop__25238_8
		f_13 = f_9
		goto b2
	} else {
		doseq_seq__25237_14 = doseq_seq__25237_7
		doseq_loop__25238_15 = doseq_loop__25238_8
		f_16 = f_9
		goto b3
	}
b2:
	;
	be_19, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "first").Deref(), []vm.Value{doseq_loop__25238_12})
	if callErr != nil {
		return nil, callErr
	}
	v21, callErr = rt.InvokeValue(rt.LookupVar("ir.passes.licm", "licm-one-loop").Deref(), []vm.Value{f_13, be_19})
	if callErr != nil {
		return nil, callErr
	}
	v23, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "next").Deref(), []vm.Value{doseq_loop__25238_12})
	if callErr != nil {
		return nil, callErr
	}
	doseq_loop__25238_8 = v23
	f_9 = f_13
	goto b1
b3:
	;
	v27 = vm.NIL
	doseq_seq__25237_28 = doseq_seq__25237_14
	doseq_loop__25238_29 = doseq_loop__25238_15
	f_30 = f_16
	goto b4
b4:
	;
	return f_30, nil
}
