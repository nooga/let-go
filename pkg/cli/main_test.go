/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package cli

import (
	"flag"
	"io"
	"reflect"
	"testing"
)

func TestParseRunFlags(t *testing.T) {
	for _, tt := range []struct {
		name    string
		args    []string
		target  string
		rest    []string
		wantErr bool
	}{
		{name: "terminal", args: []string{"-m", "app", "-v"}, target: "app", rest: []string{"-v"}},
		{name: "long equals", args: []string{"--main=app", "--"}, target: "app", rest: []string{"--"}},
		{name: "bool before main", args: []string{"-custom-bool", "-m", "app", "x"}, target: "app", rest: []string{"x"}},
		{name: "string flag value", args: []string{"-custom-string", "-m", "script.lg", "-v"}, rest: []string{"script.lg", "-v"}},
		{name: "bool equals", args: []string{"-custom-bool=false", "-m", "app"}, target: "app", rest: []string{}},
		{name: "script ends flags", args: []string{"script.lg", "-m", "app"}, rest: []string{"script.lg", "-m", "app"}},
		{name: "double dash ends flags", args: []string{"--", "-m", "app"}, rest: []string{"-m", "app"}},
		{name: "single dash positional", args: []string{"-", "-m", "app"}, rest: []string{"-", "-m", "app"}},
		{name: "unknown flag", args: []string{"-unknown", "-m", "app"}, wantErr: true},
		{name: "missing target", args: []string{"-m"}, wantErr: true},
		{name: "malformed flag", args: []string{"---m=app"}, wantErr: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			fs := flag.NewFlagSet("test", flag.ContinueOnError)
			fs.SetOutput(io.Discard)
			var target string
			fs.StringVar(&target, "m", "", "")
			fs.StringVar(&target, "main", "", "")
			fs.String("custom-string", "", "")
			fs.Bool("custom-bool", false, "")
			err := parseRunFlags(fs, tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if target != tt.target || !reflect.DeepEqual(fs.Args(), tt.rest) {
				t.Fatalf("target %q args %q, want %q %q", target, fs.Args(), tt.target, tt.rest)
			}
		})
	}
}
