//go:build !bootstrap

/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package rt

import _ "embed"

// aotDriverSrc is the AOT compile driver (nooga/let-go#596) as embedded
// namespace source, so a shipped binary can drive lowering with no let-go
// checkout on disk. scripts/lg-compile is now a thin CLI shim over it.
//
// Embedded from pkg/rt/aotdriver/ rather than pkg/rt/core/ for the same reason
// gogen is (see gogen_src.go): placing it under core/ would enroll it in the
// embedded-core lowering universe (EmbeddedNSNames) that the self-hosting
// bootstrap enumerates, and would pull it into genmanifest's sourceSpecs,
// making every edit here churn pkg/rt/generated.sums. Registered as an
// auxiliary embedded source instead, so it resolves like embedded core without
// joining that universe.
//
// Gated !bootstrap to match gogen_src.go: the bootstrap build resolves these
// namespaces via classpath, and the driver is not needed to build core.
//
//go:embed aotdriver/aot_driver.lg
var aotDriverSrc string

func init() { registerEmbeddedSource("lg.aotdriver", aotDriverSrc) }
