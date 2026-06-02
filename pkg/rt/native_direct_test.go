package rt

import "testing"

func TestNativeRegistryRegisterAndLookup(t *testing.T) {
	RegisterNativeModule(&NativeModule{
		GoPkg:     "example.com/foo",
		Namespace: "test.ns.unique1",
		Fns: map[string]NativeDirectFn{
			"bar": {GoIdent: "Bar", Arity: 2, ParamSpecs: []string{"string", "int"}, ResultSpec: "bool", NeedsError: false},
		},
	})
	d := LookupNativeDirect("test.ns.unique1", "bar")
	if d == nil || d.GoIdent != "Bar" || d.Arity != 2 || d.ResultSpec != "bool" {
		t.Fatalf("bad descriptor: %#v", d)
	}
	if LookupNativeDirect("test.ns.unique1", "missing") != nil {
		t.Fatal("expected nil for missing fn")
	}
	if AllNativeModules()["test.ns.unique1"] == nil {
		t.Fatal("module not in AllNativeModules snapshot")
	}
}
