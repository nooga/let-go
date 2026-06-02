/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"github.com/nooga/let-go/pkg/vm"
)

// installNativeDirectNS registers `(rt/native-modules)` on the "rt"
// namespace so Lisp code — chiefly ir.passes.pipeline/lower-ns-to-go —
// can inspect the registry and seed *lowered-registry* with native
// direct-callable entries.
//
// Shape of the returned value (a let-go map):
//
//	{<namespace-string> [<module>...]}
//
// where each <module> is a map:
//
//	{:go-pkg   <import-path-string>
//	 :ns       <namespace-string>
//	 :fns      {<lisp-name-string>
//	              {:go-ident    <string>
//	               :arity       <int>
//	               :variadic?   <bool>
//	               :param-specs [<string>...]
//	               :result-spec <string>
//	               :needs-error? <bool>}}}
//
// Called from lang.go's installLangNS so it runs after the "rt"
// namespace exists.
func installNativeDirectNS() {
	ns := LookupOrRegisterNS("rt")
	fn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		_ = vs
		return nativeModulesToLisp(), nil
	})
	ns.Def("native-modules", fn)

	lookup, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, nil
		}
		nsName, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, nil
		}
		name, ok := vs[1].(vm.String)
		if !ok {
			return vm.NIL, nil
		}
		d := LookupNativeDirect(string(nsName), string(name))
		if d == nil {
			return vm.NIL, nil
		}
		return descriptorToLisp(d), nil
	})
	ns.Def("native-direct", lookup)
}

func nativeModulesToLisp() vm.Value {
	all := AllNativeModules()
	out := vm.EmptyPersistentMap
	for nsName, mods := range all {
		modVec := make([]vm.Value, 0, len(mods))
		for _, m := range mods {
			modVec = append(modVec, moduleToLisp(m))
		}
		out = out.Assoc(vm.String(nsName), vm.NewPersistentVector(modVec)).(*vm.PersistentMap)
	}
	return out
}

func moduleToLisp(m *NativeModule) vm.Value {
	fnsMap := vm.EmptyPersistentMap
	for name, d := range m.Fns {
		d := d
		fnsMap = fnsMap.Assoc(vm.String(name), descriptorToLisp(&d)).(*vm.PersistentMap)
	}
	out := vm.EmptyPersistentMap
	out = out.Assoc(vm.Keyword("go-pkg"), vm.String(m.GoPkg)).(*vm.PersistentMap)
	out = out.Assoc(vm.Keyword("ns"), vm.String(m.Namespace)).(*vm.PersistentMap)
	out = out.Assoc(vm.Keyword("fns"), fnsMap).(*vm.PersistentMap)
	return out
}

func descriptorToLisp(d *NativeDirectFn) vm.Value {
	specs := make([]vm.Value, len(d.ParamSpecs))
	for i, s := range d.ParamSpecs {
		specs[i] = vm.String(s)
	}
	out := vm.EmptyPersistentMap
	out = out.Assoc(vm.Keyword("go-ident"), vm.String(d.GoIdent)).(*vm.PersistentMap)
	out = out.Assoc(vm.Keyword("arity"), vm.Int(int64(d.Arity))).(*vm.PersistentMap)
	out = out.Assoc(vm.Keyword("variadic?"), vm.Boolean(d.Variadic)).(*vm.PersistentMap)
	out = out.Assoc(vm.Keyword("param-specs"), vm.NewPersistentVector(specs)).(*vm.PersistentMap)
	out = out.Assoc(vm.Keyword("result-spec"), vm.String(d.ResultSpec)).(*vm.PersistentMap)
	out = out.Assoc(vm.Keyword("needs-error?"), vm.Boolean(d.NeedsError)).(*vm.PersistentMap)
	return out
}
