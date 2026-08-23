/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package main

import "testing"

func TestValidateOutPkg(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		wantErr bool
	}{
		{"empty is the in-tree default", "", false},
		{"ordinary name", "interop", false},
		{"underscores and digits", "go_interop2", false},
		{"leading underscore", "_interop", false},
		{"mixed case", "myInterop", false},
		{"rt is ambiguous with the in-tree output", "rt", true},
		{"blank identifier", "_", true},
		{"go keyword", "package", true},
		{"another go keyword", "func", true},
		{"hyphen is not an identifier char", "my-pkg", true},
		{"dot is not an identifier char", "my.pkg", true},
		{"leading digit", "1interop", true},
		{"space", "my pkg", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateOutPkg(tt.in)
			if tt.wantErr && err == nil {
				t.Errorf("validateOutPkg(%q) = nil, want error", tt.in)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("validateOutPkg(%q) = %v, want nil", tt.in, err)
			}
		})
	}
}
