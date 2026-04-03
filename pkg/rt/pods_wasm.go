//go:build js && wasm

/*
 * Copyright (c) 2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

// Pods are not supported in WASM — subprocess execution is unavailable.

func installPodsNS() {}

func ShutdownAllPods() {}
