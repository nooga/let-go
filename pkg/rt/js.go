/*
 * Copyright (c) 2026 Matt Parrett
 * SPDX-License-Identifier: MIT
 *
 * Platform-independent half of the js namespace. Argument validation and
 * marshaling live here so that bugs (wrong arity, bad event-name type) are
 * caught the same way whether .lg code is running native or in WASM. The
 * platform-specific files (js_wasm.go, js_other.go) only provide the
 * dispatch primitive.
 */

package rt

import (
	"encoding/json"
	"fmt"

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
	buf, err := json.Marshal(data)
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
