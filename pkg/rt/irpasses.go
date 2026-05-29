/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package rt

import _ "embed"

//go:embed core/ir/zipper.lg
var IRZipperSrc string

//go:embed core/ir/passes.lg
var IRPassesSrc string

//go:embed core/ir/passes/dce.lg
var IRPassDCESrc string

//go:embed core/ir/passes/constfold.lg
var IRPassConstFoldSrc string

//go:embed core/ir/passes/cse.lg
var IRPassCSESrc string

//go:embed core/ir/passes/mutability.lg
var IRPassMutabilitySrc string

//go:embed core/ir/passes/typeinfer.lg
var IRPassTypeInferSrc string

//go:embed core/ir/passes/licm.lg
var IRPassLICMSrc string

//go:embed core/ir/passes/pipeline.lg
var IRPassPipelineSrc string

//go:embed core/ir/passes/trace.lg
var IRPassTraceSrc string

//go:embed core/ir/dump.lg
var IRDumpSrc string

//go:embed core/ir/dominance.lg
var IRDominanceSrc string

//go:embed core/ir/lower.lg
var IRLowerSrc string

//go:embed core/ir/lower_go.lg
var IRLowerGoSrc string

//go:embed core/ir/data.lg
var IRDataSrc string

// force-rebuild-marker-2
//
//go:embed core/ir/validate.lg
var IRValidateSrc string
