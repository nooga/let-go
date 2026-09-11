/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package test

import "testing"

func TestTypedArraysAcceptEmptySeqableInputs(t *testing.T) {
	got, err := evalCoreCompat(`[(alength (byte-array []))
	                             (alength (int-array []))
	                             (alength (long-array []))
	                             (alength (double-array []))
	                             (alength (float-array []))
	                             (alength (object-array []))
	                             (alength (int-array '()))
	                             (alength (int-array nil))]`)
	if err != nil {
		t.Fatal(err)
	}
	if got.String() != "[0 0 0 0 0 0 0 0]" {
		t.Fatalf("empty seqable inputs must produce empty arrays: got %s", got)
	}
}

func TestByteArrayMutationSurvivesGoInterop(t *testing.T) {
	got, err := evalCoreCompat(`(do
		(require '[io :as io])
		(let [r (io/string-reader "ABC")
		      b (byte-array 3)
		      n (.Read r b)]
		  [n (vec b)]))`)
	if err != nil {
		t.Fatal(err)
	}
	if got.String() != "[3 [65 66 67]]" {
		t.Fatalf("Go mutation must remain visible through the typed array: got %s", got)
	}
}

func TestContainsRejectsUnsupportedCollectionsAndArrayKeys(t *testing.T) {
	cases := map[string]string{
		"array nil key":     `(contains? (int-array [0 1 2]) nil)`,
		"array keyword key": `(contains? (int-array [0 1 2]) :a)`,
		"array vector key":  `(contains? (int-array [0 1 2]) [0 1 2])`,
		"list":              `(contains? '(1 2 3) 0)`,
		"number":            `(contains? 42 0)`,
		"keyword":           `(contains? :a :a)`,
	}
	for name, expr := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := evalCoreCompat(expr); err == nil {
				t.Fatalf("expected contains? to reject %s", expr)
			}
		})
	}
}

func TestContainsSupportsMapEntryIndexes(t *testing.T) {
	got, err := evalCoreCompat(`[(contains? (first {:a 1}) 0)
	                             (contains? (first {:a 1}) 1)
	                             (contains? (first {:a 1}) 2)]`)
	if err != nil {
		t.Fatal(err)
	}
	if got.String() != "[true true false]" {
		t.Fatalf("map entries must remain indexed: got %s", got)
	}
}
