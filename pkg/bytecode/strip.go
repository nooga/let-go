/*
 * Copyright (c) 2026 Matt Parrett <matt.parrett@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package bytecode

import "bytes"

// StripDebug re-encodes an .lgb without its debug sections: per-chunk source
// maps and local-variable tables. Execution is unaffected — the VM consults
// both only for error reporting and debugging — so a stripped bundle trades
// source-located errors for size (about 20% on the core bundle). Existing
// encoding flags such as FlagCompressed are preserved, and FlagDebugSplit is
// added so a companion produced by SplitDebug can safely restore the tables.
func StripDebug(data []byte) ([]byte, error) {
	m, err := Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	for _, c := range m.Chunks {
		c.SourceMap = nil
		c.LocalVars = nil
	}
	// The encoder echoes Flags as given; with every table emptied the
	// localvars section should be dropped, not written as zero counts.
	m.Flags &^= FlagLocalVars
	m.Flags |= FlagDebugSplit
	var buf bytes.Buffer
	if err := Encode(&buf, m); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
