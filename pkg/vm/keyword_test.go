/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import "testing"

func TestKeywordNamespacedSplitsAtFirstSlash(t *testing.T) {
	keyword := Keyword("a/b/c")

	rawNS, rawName, hasNS := keyword.NamespacedRaw()
	if rawNS != Keyword("a") || rawName != Keyword("b/c") || !hasNS {
		t.Fatalf("NamespacedRaw() = (%v, %v, %t), want (a, b/c, true)", rawNS, rawName, hasNS)
	}
	ns, name := keyword.Namespaced()
	if ns != Symbol("a") || name != Symbol("b/c") {
		t.Fatalf("Namespaced() = (%v, %v), want (a, b/c)", ns, name)
	}
	if got := keyword.Namespace(); got != String("a") {
		t.Fatalf("Namespace() = %v, want a", got)
	}
	if got := keyword.Name(); got != String("b/c") {
		t.Fatalf("Name() = %v, want b/c", got)
	}
}

func TestKeywordNamespacedMatchesSymbolSlashSemantics(t *testing.T) {
	tests := []struct {
		keyword   Keyword
		rawNS     Keyword
		rawName   Keyword
		hasNS     bool
		ns        Value
		name      Value
		namespace Value
		nameValue Value
	}{
		{keyword: "plain", rawName: "plain", ns: NIL, name: Symbol("plain"), namespace: NIL, nameValue: String("plain")},
		{keyword: "/", rawName: "/", ns: NIL, name: Symbol("/"), namespace: NIL, nameValue: String("/")},
		{keyword: "/name", rawName: "name", hasNS: true, ns: Symbol(""), name: Symbol("name"), namespace: String(""), nameValue: String("name")},
		{keyword: "ns/", rawNS: "ns", hasNS: true, ns: Symbol("ns"), name: Symbol(""), namespace: String("ns"), nameValue: String("")},
		{keyword: "", ns: NIL, name: Symbol(""), namespace: NIL, nameValue: String("")},
		{keyword: "a//", rawNS: "a", rawName: "/", hasNS: true, ns: Symbol("a"), name: Symbol("/"), namespace: String("a"), nameValue: String("/")},
	}

	for _, test := range tests {
		t.Run(string(test.keyword), func(t *testing.T) {
			rawNS, rawName, hasNS := test.keyword.NamespacedRaw()
			if rawNS != test.rawNS || rawName != test.rawName || hasNS != test.hasNS {
				t.Fatalf("NamespacedRaw() = (%v, %v, %t), want (%v, %v, %t)", rawNS, rawName, hasNS, test.rawNS, test.rawName, test.hasNS)
			}
			ns, name := test.keyword.Namespaced()
			if ns != test.ns || name != test.name {
				t.Fatalf("Namespaced() = (%v, %v), want (%v, %v)", ns, name, test.ns, test.name)
			}
			if namespace := test.keyword.Namespace(); namespace != test.namespace {
				t.Fatalf("Namespace() = %v, want %v", namespace, test.namespace)
			}
			if nameValue := test.keyword.Name(); nameValue != test.nameValue {
				t.Fatalf("Name() = %v, want %v", nameValue, test.nameValue)
			}
		})
	}
}
