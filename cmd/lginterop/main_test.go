/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package main

import (
	"strings"
	"testing"
)

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
		{"main is not importable", "main", true},
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

// validateEntries must refuse inputs that would emit uncompilable output
// (aliases landing on the file's own imports) or silently overwrite one
// generated file with another (distinct aliases normalizing to one filename).
func TestValidateEntries(t *testing.T) {
	tests := []struct {
		name    string
		entries []interopEntry
		wantErr string // substring; empty means valid
	}{
		{
			name:    "distinct aliases pass",
			entries: []interopEntry{{pkg: "hash/crc32"}, {pkg: "image/color"}},
		},
		{
			name:    "alias vm collides with the emitted vm import",
			entries: []interopEntry{{pkg: "hash/crc32", alias: "vm"}},
			wantErr: "vm import",
		},
		{
			name:    "package path ending in /vm collides via its default alias",
			entries: []interopEntry{{pkg: "example.com/some/vm"}},
			wantErr: "vm import",
		},
		{
			name:    "alias fmt collides in smart mode",
			entries: []interopEntry{{pkg: "example.com/some/fmt", smart: true}},
			wantErr: "fmt import",
		},
		{
			name:    "alias fmt without smart mode is fine",
			entries: []interopEntry{{pkg: "example.com/some/fmt"}},
		},
		{
			name:    "same alias twice collides",
			entries: []interopEntry{{pkg: "a/x", alias: "dup"}, {pkg: "b/y", alias: "dup"}},
			wantErr: "collide",
		},
		{
			name: "distinct aliases normalizing to one filename collide",
			entries: []interopEntry{
				{pkg: "a/x", alias: "foo-bar"},
				{pkg: "b/y", alias: "foo_bar"},
			},
			wantErr: "interop_foo_bar.go",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEntries(tt.entries)
			if tt.wantErr == "" {
				if err != nil {
					t.Errorf("validateEntries() = %v, want nil", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("validateEntries() = %v, want error containing %q", err, tt.wantErr)
			}
		})
	}
}
