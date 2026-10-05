/*
 * Copyright (c) 2026 Matt Parrett
 * SPDX-License-Identifier: MIT
 *
 * The js namespace. (js/emit ...) validates and marshals its args here, then
 * dispatches through the *emit* host seam (see emitter.go) — so the native
 * code path is identical on every platform. The platform difference is only
 * which Emitter sits at the *emit* root: HostEmitter in the WASM bundle
 * (hostemitter_js_wasm.go), a FuncEmitter via api.WithEmit for Go embedders,
 * or the default no-op otherwise.
 */

package rt

import (
	"fmt"

	"github.com/nooga/let-go/pkg/rt/internal/jsonenc"
	"github.com/nooga/let-go/pkg/vm"
)

// prepareEmit validates args for (js/emit event-name data) and returns the
// event name and the JSON-marshaled data ready to hand to the platform
// dispatcher. Same contract on every platform.
func prepareEmit(vs []vm.Value) (string, string, error) {
	if len(vs) != 2 {
		return "", "", fmt.Errorf("js/emit expects 2 args (event-name data), got %d", len(vs))
	}
	name, err := eventName(vs[0])
	if err != nil {
		return "", "", err
	}
	data, err := fromValue(vs[1])
	if err != nil {
		return "", "", err
	}
	buf, err := jsonenc.Marshal(data)
	if err != nil {
		return "", "", fmt.Errorf("js/emit: %w", err)
	}
	return name, string(buf), nil
}

// eventName coerces a let-go value into the string event name passed to
// CustomEvent. Accepts keyword (:stats), symbol (stats), or string ("stats").
func eventName(v vm.Value) (string, error) {
	switch v.Type() {
	case vm.KeywordType:
		return string(v.(vm.Keyword)), nil
	case vm.SymbolType:
		return v.(vm.Symbol).String(), nil
	case vm.StringType:
		return string(v.(vm.String)), nil
	default:
		return "", fmt.Errorf("js/emit event-name must be keyword, symbol, or string; got %s", v.Type().Name())
	}
}

func init() { RegisterInstaller(installJSNS) }

func installJSNS() {
	ns := vm.NewNamespace("js")

	// (js/emit event-name data) -> nil. Fire-and-forget through the *emit*
	// host seam. Arg validation runs on every platform so type bugs surface
	// in native dev, not just at runtime in the browser. Ctx-aware so the
	// dispatch respects the current (binding [*emit* ...]) / api.WithEmit.
	emitFn := vm.NewCtxNativeFn("emit", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		name, dataJSON, err := prepareEmit(vs)
		if err != nil {
			return vm.NIL, err
		}
		EmitVia(ec, name, dataJSON)
		return vm.NIL, nil
	})
	ns.Def("emit", emitFn)

	// (js/url-param name) -> string-or-nil. Reads a query parameter from the
	// page URL that loaded the bundle; nil if absent or off-browser. Lets .lg
	// boot code parameterize from the URL (?seed=42, ?bench=fast) without JS
	// plumbing. The page URL reaches the worker via the host seam below
	// (see hostURLParam).
	urlParamFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, nil
		}
		name, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, nil
		}
		if v, ok := hostURLParam(string(name)); ok {
			return vm.String(v), nil
		}
		return vm.NIL, nil
	})
	ns.Def("url-param", urlParamFn)

	RegisterNS(ns)
}
