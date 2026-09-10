/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

// Package empty deliberately exports nothing. It exists so tests can pin
// cmd/lginterop's skip accounting: a scanned package with no eligible exports
// must be reported as skipped, not counted as generated output that was never
// written.
package empty

var unexported = 0

var _ = unexported
