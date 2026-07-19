//go:build !bootstrap

/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import _ "embed"

// gogenSrc is the gogen macro layer, embedded into the binary so `(require
// 'gogen)` resolves with no external source path (nooga/let-go#425: self-contain
// gogen — it was formerly an external classpath dir, which broke `lg-compile`
// from the wrong cwd / with older binaries).
//
// It is embedded from pkg/rt/gogen/ rather than pkg/rt/core/ on purpose: gogen
// is the AOT Go-emitter, not core content. Placing it under core/ would enroll
// it in the embedded-core lowering universe (EmbeddedNSNames), which the
// self-hosting bootstrap enumerates — enough to perturb the interpreted-vs-
// native lowering fixpoint (TestLoweringDeterminism). Registered as an
// auxiliary embedded source (see EmbeddedSource) instead, so it resolves like
// embedded core without joining that universe.
//
//go:embed gogen/gogen.lg
var gogenSrc string

func init() { registerEmbeddedSource("gogen", gogenSrc) }
