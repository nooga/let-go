/*
 * Copyright (c) 2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/nooga/let-go/pkg/vm"
)

// Transit+JSON codec for let-go values.
// Implements the transit-json spec: https://github.com/cognitect/transit-format

// --- Rolling cache (shared by encoder and decoder) ---

const (
	cacheDigits      = 44
	cacheBase        = 48 // '0'
	cacheSize        = cacheDigits * cacheDigits
	minCacheableLen  = 4
	cacheMarker = "^"
	mapMarker   = "^ "
	kwPrefix    = "~:"
	symPrefix   = "~$"
)

type transitCache struct {
	keyToVal map[string]string
	valToKey map[string]string
	idx      int
}

func newTransitCache() *transitCache {
	return &transitCache{
		keyToVal: make(map[string]string),
		valToKey: make(map[string]string),
	}
}

func (c *transitCache) encodeKey(idx int) string {
	hi := idx / cacheDigits
	lo := idx % cacheDigits
	if hi == 0 {
		return cacheMarker + string(rune(lo+cacheBase))
	}
	return cacheMarker + string(rune(hi+cacheBase)) + string(rune(lo+cacheBase))
}

func isCacheRef(s string) bool {
	return len(s) > 0 && s[0] == '^' && s != mapMarker
}

func (c *transitCache) isCacheable(s string, asKey bool) bool {
	if len(s) < minCacheableLen {
		return false
	}
	if asKey {
		return true
	}
	if len(s) >= 2 && s[0] == '~' {
		return s[1] == '#' || s[1] == ':' || s[1] == '$'
	}
	return false
}

// read resolves a cache ref or caches the string and returns the resolved value.
func (c *transitCache) read(s string, asKey bool) string {
	if isCacheRef(s) {
		if v, ok := c.keyToVal[s]; ok {
			return v
		}
		return s
	}
	if c.isCacheable(s, asKey) {
		if len(c.keyToVal) >= cacheSize {
			c.keyToVal = make(map[string]string)
			c.valToKey = make(map[string]string)
			c.idx = 0
		}
		key := c.encodeKey(c.idx)
		c.keyToVal[key] = s
		c.valToKey[s] = key
		c.idx++
	}
	return s
}

// write returns the cache ref for a string if cached, otherwise caches it.
func (c *transitCache) write(s string, asKey bool) string {
	if ref, ok := c.valToKey[s]; ok {
		return ref
	}
	if c.isCacheable(s, asKey) {
		if len(c.keyToVal) >= cacheSize {
			c.keyToVal = make(map[string]string)
			c.valToKey = make(map[string]string)
			c.idx = 0
		}
		key := c.encodeKey(c.idx)
		c.keyToVal[key] = s
		c.valToKey[s] = key
		c.idx++
	}
	return s
}

// --- Decoder ---

// TransitDecoder decodes transit+json into vm.Values.
type TransitDecoder struct {
	cache *transitCache
}

// NewTransitDecoder creates a new decoder with a fresh cache.
func NewTransitDecoder() *TransitDecoder {
	return &TransitDecoder{cache: newTransitCache()}
}

// Decode parses a transit+json string into a vm.Value.
func (d *TransitDecoder) Decode(s string) (vm.Value, error) {
	var raw interface{}
	if err := json.Unmarshal([]byte(s), &raw); err != nil {
		return vm.NIL, fmt.Errorf("transit: invalid JSON: %w", err)
	}
	return d.decodeValue(raw)
}

func (d *TransitDecoder) decodeValue(v interface{}) (vm.Value, error) {
	switch val := v.(type) {
	case string:
		return d.decodeString(val, false), nil
	case float64:
		if val == float64(int64(val)) {
			return vm.Int(int64(val)), nil
		}
		return vm.Float(val), nil
	case bool:
		return vm.Boolean(val), nil
	case nil:
		return vm.NIL, nil
	case []interface{}:
		return d.decodeArray(val)
	case map[string]interface{}:
		return d.decodeObject(val)
	default:
		return vm.NIL, fmt.Errorf("transit: unsupported JSON type %T", v)
	}
}

func (d *TransitDecoder) decodeString(s string, asKey bool) vm.Value {
	s = d.cache.read(s, asKey)
	return parseTransit(s)
}

// parseTransit interprets a resolved transit string.
func parseTransit(s string) vm.Value {
	if len(s) < 2 || s[0] != '~' {
		return vm.String(s)
	}
	switch s[1] {
	case ':':
		return vm.Keyword(s[2:])
	case '$':
		return vm.Symbol(s[2:])
	case '~':
		return vm.String(s[1:]) // escaped ~
	case 'i':
		// Big integer
		n := new(big.Int)
		if _, ok := n.SetString(s[2:], 10); ok {
			if n.IsInt64() {
				return vm.Int(n.Int64())
			}
			return vm.NewBigInt(n)
		}
		return vm.String(s)
	case 'f':
		// Float as string
		var f float64
		if _, err := fmt.Sscanf(s[2:], "%g", &f); err == nil {
			return vm.Float(f)
		}
		return vm.String(s)
	case 'b':
		// Boolean encoded as string (rare)
		if s[2:] == "1" || s[2:] == "t" {
			return vm.TRUE
		}
		return vm.FALSE
	case 'n':
		return vm.NIL
	case '#':
		// Tag marker - not a final value, but we may encounter it standalone
		return vm.String(s)
	default:
		return vm.String(s)
	}
}

func (d *TransitDecoder) decodeArray(arr []interface{}) (vm.Value, error) {
	if len(arr) == 0 {
		return vm.NewArrayVector(nil), nil
	}

	// Check first element for special markers
	if first, ok := arr[0].(string); ok {
		resolved := d.cache.read(first, false)

		// Map-as-array: ["^ ", k1, v1, k2, v2, ...]
		if resolved == mapMarker {
			return d.decodeMapArray(arr)
		}

		// Tagged value: ["~#tag", payload]
		if len(resolved) >= 2 && resolved[0] == '~' && resolved[1] == '#' {
			if len(arr) == 2 {
				return d.decodeTagged(resolved[2:], arr[1])
			}
		}
	}

	// Regular vector
	vals := make([]vm.Value, len(arr))
	for i, item := range arr {
		v, err := d.decodeValue(item)
		if err != nil {
			return vm.NIL, err
		}
		vals[i] = v
	}
	return vm.NewArrayVector(vals), nil
}

func (d *TransitDecoder) decodeMapArray(arr []interface{}) (vm.Value, error) {
	m := vm.EmptyPersistentMap
	for i := 1; i+1 < len(arr); i += 2 {
		k, err := d.decodeKeyValue(arr[i])
		if err != nil {
			return vm.NIL, err
		}
		v, err := d.decodeValue(arr[i+1])
		if err != nil {
			return vm.NIL, err
		}
		m = m.Assoc(k, v).(*vm.PersistentMap)
	}
	return m, nil
}

func (d *TransitDecoder) decodeTagged(tag string, payload interface{}) (vm.Value, error) {
	switch tag {
	case "set":
		if items, ok := payload.([]interface{}); ok {
			set := vm.EmptyPersistentSet
			for _, item := range items {
				v, err := d.decodeValue(item)
				if err != nil {
					return vm.NIL, err
				}
				set = set.Conj(v).(*vm.PersistentSet)
			}
			return set, nil
		}
	case "list":
		if items, ok := payload.([]interface{}); ok {
			vals := make([]vm.Value, len(items))
			for i, item := range items {
				v, err := d.decodeValue(item)
				if err != nil {
					return vm.NIL, err
				}
				vals[i] = v
			}
			var result vm.Value = vm.EmptyList
			for i := len(vals) - 1; i >= 0; i-- {
				result = vm.NewCons(vals[i], result.(vm.Seq))
			}
			return result, nil
		}
	case "'":
		// Quoted value - unwrap
		return d.decodeValue(payload)
	case "cmap":
		// Verbose map (array of k-v pairs when keys aren't stringable)
		if items, ok := payload.([]interface{}); ok {
			m := vm.EmptyPersistentMap
			for i := 0; i+1 < len(items); i += 2 {
				k, err := d.decodeValue(items[i])
				if err != nil {
					return vm.NIL, err
				}
				v, err := d.decodeValue(items[i+1])
				if err != nil {
					return vm.NIL, err
				}
				m = m.Assoc(k, v).(*vm.PersistentMap)
			}
			return m, nil
		}
	}
	// Unknown tag - decode the payload as-is
	return d.decodeValue(payload)
}

// decodeKeyValue decodes a value in key position (enables caching for map keys).
func (d *TransitDecoder) decodeKeyValue(v interface{}) (vm.Value, error) {
	if s, ok := v.(string); ok {
		return d.decodeString(s, true), nil
	}
	return d.decodeValue(v)
}

func (d *TransitDecoder) decodeObject(m map[string]interface{}) (vm.Value, error) {
	// JSON objects are rare in transit+json (maps use array encoding)
	// but can appear for string-keyed maps
	result := vm.EmptyPersistentMap
	for k, v := range m {
		key := d.decodeString(k, true)
		val, err := d.decodeValue(v)
		if err != nil {
			return vm.NIL, err
		}
		result = result.Assoc(key, val).(*vm.PersistentMap)
	}
	return result, nil
}

// --- Encoder ---

// TransitEncoder encodes vm.Values to transit+json.
type TransitEncoder struct {
	cache *transitCache
}

// NewTransitEncoder creates a new encoder with a fresh cache.
func NewTransitEncoder() *TransitEncoder {
	return &TransitEncoder{cache: newTransitCache()}
}

// Encode encodes a vm.Value to a transit+json string.
func (e *TransitEncoder) Encode(v vm.Value) (string, error) {
	raw, err := e.encodeValue(v)
	if err != nil {
		return "", err
	}
	bs, err := json.Marshal(raw)
	if err != nil {
		return "", fmt.Errorf("transit: JSON marshal error: %w", err)
	}
	return string(bs), nil
}

// EncodeList encodes a slice of values as a transit list (used for pod args).
func (e *TransitEncoder) EncodeList(vals []vm.Value) (string, error) {
	items := make([]interface{}, len(vals))
	for i, v := range vals {
		raw, err := e.encodeValue(v)
		if err != nil {
			return "", err
		}
		items[i] = raw
	}
	list := []interface{}{"~#list", items}
	bs, err := json.Marshal(list)
	return string(bs), err
}

func (e *TransitEncoder) encodeValue(v vm.Value) (interface{}, error) {
	switch v.Type() {
	case vm.NilType:
		return nil, nil
	case vm.BooleanType:
		return bool(v.(vm.Boolean)), nil
	case vm.IntType:
		n := int64(v.(vm.Int))
		// Transit spec: ints outside safe JS range use "~i" string encoding
		if n > 1<<53-1 || n < -(1<<53-1) {
			return fmt.Sprintf("~i%d", n), nil
		}
		return n, nil
	case vm.FloatType:
		return float64(v.(vm.Float)), nil
	case vm.StringType:
		s := string(v.(vm.String))
		if len(s) > 0 && (s[0] == '~' || s[0] == '^') {
			return "~" + s, nil // escape
		}
		return s, nil
	case vm.KeywordType:
		raw := kwPrefix + string(v.(vm.Keyword))
		return e.cache.write(raw, false), nil
	case vm.SymbolType:
		raw := symPrefix + string(v.(vm.Symbol))
		return e.cache.write(raw, false), nil
	case vm.BigIntType:
		bi := v.(*vm.BigInt)
		return "~i" + bi.String(), nil

	case vm.ArrayVectorType, vm.PersistentVectorType:
		return e.encodeSeqAsArray(v)
	case vm.MapType, vm.PersistentMapType:
		return e.encodeMap(v)
	case vm.SetType:
		items, err := e.encodeSeqAsArray(v)
		if err != nil {
			return nil, err
		}
		return []interface{}{"~#set", items}, nil
	default:
		// Lists and other seqs -> transit list tag
		if sq, ok := v.(vm.Sequable); ok {
			seq := sq.Seq()
			if seq == nil {
				return []interface{}{"~#list", []interface{}{}}, nil
			}
			var items []interface{}
			for s := seq; s != nil; s = s.Next() {
				item, err := e.encodeValue(s.First())
				if err != nil {
					return nil, err
				}
				items = append(items, item)
			}
			return []interface{}{"~#list", items}, nil
		}
		// Records and other map-like types
		if _, ok := v.(*vm.Record); ok {
			return e.encodeMap(v)
		}
		// Fallback: string representation
		return v.String(), nil
	}
}

func (e *TransitEncoder) encodeSeqAsArray(v vm.Value) ([]interface{}, error) {
	sq, ok := v.(vm.Sequable)
	if !ok {
		return nil, fmt.Errorf("transit: cannot seq %s", v.Type().Name())
	}
	var result []interface{}
	for s := sq.Seq(); s != nil; s = s.Next() {
		item, err := e.encodeValue(s.First())
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	if result == nil {
		result = []interface{}{}
	}
	return result, nil
}

func (e *TransitEncoder) encodeMap(v vm.Value) (interface{}, error) {
	sq, ok := v.(vm.Sequable)
	if !ok {
		return nil, fmt.Errorf("transit: cannot seq map %s", v.Type().Name())
	}
	result := []interface{}{mapMarker}
	for s := sq.Seq(); s != nil; s = s.Next() {
		entry := s.First()
		eSeq, ok := entry.(vm.Sequable)
		if !ok {
			continue
		}
		es := eSeq.Seq()
		k, err := e.encodeKeyValue(es.First())
		if err != nil {
			return nil, err
		}
		val, err := e.encodeValue(es.Next().First())
		if err != nil {
			return nil, err
		}
		result = append(result, k, val)
	}
	return result, nil
}

func (e *TransitEncoder) encodeKeyValue(v vm.Value) (interface{}, error) {
	switch v.Type() {
	case vm.KeywordType:
		raw := kwPrefix + string(v.(vm.Keyword))
		return e.cache.write(raw, true), nil
	case vm.SymbolType:
		raw := symPrefix + string(v.(vm.Symbol))
		return e.cache.write(raw, true), nil
	case vm.StringType:
		s := string(v.(vm.String))
		if len(s) > 0 && (s[0] == '~' || s[0] == '^') {
			s = "~" + s
		}
		return e.cache.write(s, true), nil
	default:
		return e.encodeValue(v)
	}
}

// --- Namespace ---

func installTransitNS() {
	readTransit, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("transit/read: expected 1 argument")
		}
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("transit/read: expected String, got %s", vs[0].Type().Name())
		}
		d := NewTransitDecoder()
		return d.Decode(string(s))
	})

	writeTransit, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("transit/write: expected 1 argument")
		}
		enc := NewTransitEncoder()
		s, err := enc.Encode(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		return vm.String(s), nil
	})

	// read-many: decode a transit string that contains multiple values (e.g. a list)
	// This is mainly for debugging/testing transit payloads.

	ns := vm.NewNamespace("transit")
	ns.Def("read", readTransit)
	ns.Def("write", writeTransit)

	// Expose verbose-write for debugging (no cache, readable output)
	verboseWrite, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("transit/write-verbose: expected 1 argument")
		}
		// Use a fresh encoder for each call = no cache refs in output
		enc := NewTransitEncoder()
		s, err := enc.Encode(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		return vm.String(s), nil
	})
	ns.Def("write-verbose", verboseWrite)

	RegisterNS(ns)
}

// --- Helpers used by pods ---

// TransitEncodeArgs encodes a list of args for pod invocation.
func TransitEncodeArgs(args []vm.Value) (string, error) {
	enc := NewTransitEncoder()
	return enc.EncodeList(args)
}

// TransitDecodeValue decodes a transit+json payload string.
func TransitDecodeValue(s string) (vm.Value, error) {
	dec := NewTransitDecoder()
	return dec.Decode(s)
}

// prStr formats a value as EDN (used for EDN payload encoding).
func prStr(v vm.Value) string {
	switch v.Type() {
	case vm.StringType:
		bs, _ := json.Marshal(string(v.(vm.String)))
		return string(bs)
	case vm.NilType:
		return "nil"
	case vm.BooleanType:
		if bool(v.(vm.Boolean)) {
			return "true"
		}
		return "false"
	default:
		return v.String()
	}
}

// readEDN is set by the compiler package to provide EDN parsing.
var readEDN func(string) (vm.Value, error)

// SetReadEDN sets the EDN reader function (called by compiler package).
func SetReadEDN(fn func(string) (vm.Value, error)) {
	readEDN = fn
}

// evalInNS is set by the compiler package to evaluate code in a namespace.
var evalInNS func(code string, ns *vm.Namespace) (vm.Value, error)

// SetEvalInNS sets the namespace-aware eval function.
func SetEvalInNS(fn func(string, *vm.Namespace) (vm.Value, error)) {
	evalInNS = fn
}

// JSONEncodeArgs encodes args as a plain JSON array string.
func JSONEncodeArgs(args []vm.Value) (string, error) {
	goArgs := make([]interface{}, len(args))
	for i, a := range args {
		v, err := fromValue(a)
		if err != nil {
			return "", err
		}
		goArgs[i] = v
	}
	bs, err := json.Marshal(goArgs)
	return string(bs), err
}

// JSONDecodeValue decodes a JSON string into a vm.Value.
func JSONDecodeValue(s string) (vm.Value, error) {
	var v interface{}
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		return vm.NIL, err
	}
	return toValue(true, v)
}

// EDNEncodeArgs encodes args as an EDN vector string.
func EDNEncodeArgs(args []vm.Value) (string, error) {
	parts := make([]string, len(args))
	for i, a := range args {
		parts[i] = prStr(a)
	}
	return "[" + strings.Join(parts, " ") + "]", nil
}
