//go:build !bootstrap

package rt

import _ "embed"

//go:generate go run -tags bootstrap ../../cmd/lgbgen core_compiled.lgb

//go:embed core_compiled.lgb
var CoreCompiledLGB []byte
