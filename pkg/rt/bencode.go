/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

/*
 * bencode namespace — message framing over a net conn.
 *
 *   (bencode/write! conn value)       → nil. Encodes one bencode value.
 *   (bencode/read! conn)              → decoded value, nil on clean EOF. Blocks.
 *   (bencode/read! conn timeout-ms)   → same, but errors with a message
 *                                       containing "timeout" if nothing arrives
 *                                       in time (SetReadDeadline, cleared after).
 *
 * Value conversion (pinned both directions):
 *
 *   lg → bencode                         bencode → lg
 *   String        → byte string          byte string → String
 *   keyword       → its name             int         → Int
 *   Int           → int                  list        → vector
 *   vector/list   → list                 dict        → map with String keys
 *   map           → dict (String/keyword keys)
 *   nil/bool/float → error (no bencode representation)
 */

package rt

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/nooga/let-go/pkg/vm"
	"github.com/zeebo/bencode"
)

func init() { RegisterInstaller(installBencodeNS) }

// bencodeFromLG converts an lg value into a Go value the bencode encoder
// accepts. nil/bool/float are rejected: bencode has no representation for them
// and silently coercing would corrupt nREPL messages.
func bencodeFromLG(v vm.Value) (interface{}, error) {
	switch v.Type() {
	case vm.StringType:
		return string(v.(vm.String)), nil
	case vm.KeywordType:
		return string(v.(vm.Keyword)), nil // :op → "op"
	case vm.IntType:
		return int64(v.(vm.Int)), nil
	case vm.MapType, vm.PersistentMapType:
		return bencodeMapToDict(v)
	case vm.ArrayVectorType, vm.PersistentVectorType, vm.ListType:
		if sq, ok := v.(vm.Sequable); ok {
			return bencodeSeqToList(sq.Seq())
		}
	}
	return nil, fmt.Errorf("bencode/write! cannot encode %s value", v.Type().Name())
}

func bencodeSeqToList(s vm.Seq) ([]interface{}, error) {
	out := []interface{}{}
	for s != nil && s != vm.EmptyList {
		e, err := bencodeFromLG(s.First())
		if err != nil {
			return nil, err
		}
		out = append(out, e)
		s = s.Next()
	}
	return out, nil
}

func bencodeMapToDict(v vm.Value) (map[string]interface{}, error) {
	sq, ok := v.(vm.Sequable)
	if !ok {
		return nil, fmt.Errorf("bencode/write! cannot iterate %s", v.Type().Name())
	}
	out := map[string]interface{}{}
	for s := sq.Seq(); s != nil && s != vm.EmptyList; s = s.Next() {
		entry, ok := s.First().(vm.Sequable)
		if !ok {
			return nil, fmt.Errorf("bencode/write! invalid map entry")
		}
		es := entry.Seq()
		key, err := bencodeMapKey(es.First())
		if err != nil {
			return nil, err
		}
		val, err := bencodeFromLG(es.Next().First())
		if err != nil {
			return nil, err
		}
		out[key] = val
	}
	return out, nil
}

func bencodeMapKey(k vm.Value) (string, error) {
	switch k.Type() {
	case vm.StringType:
		return string(k.(vm.String)), nil
	case vm.KeywordType:
		return string(k.(vm.Keyword)), nil
	}
	return "", fmt.Errorf("bencode/write! map keys must be String or keyword, got %s", k.Type().Name())
}

// bencodeToLG converts a decoded bencode value (string/int64/[]interface{}/
// map[string]interface{}) into lg values. Dicts become maps with String keys.
func bencodeToLG(i interface{}) (vm.Value, error) {
	switch x := i.(type) {
	case string:
		return vm.String(x), nil
	case []byte:
		return vm.String(x), nil
	case int64:
		return vm.Int(x), nil
	case []interface{}:
		out := make([]vm.Value, len(x))
		for j := range x {
			e, err := bencodeToLG(x[j])
			if err != nil {
				return vm.NIL, err
			}
			out[j] = e
		}
		return vm.NewArrayVector(out), nil
	case map[string]interface{}:
		m := vm.EmptyPersistentMap
		for k, v := range x {
			cv, err := bencodeToLG(v)
			if err != nil {
				return vm.NIL, err
			}
			m = m.Assoc(vm.String(k), cv).(*vm.PersistentMap)
		}
		return m, nil
	case nil:
		return vm.NIL, nil
	}
	return vm.NIL, fmt.Errorf("bencode/read! decoded unsupported value %T", i)
}

func installBencodeNS() {
	writeFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("bencode/write! expects 2 args (conn value)")
		}
		c, err := unboxNetConn(vs[0])
		if err != nil {
			return vm.NIL, fmt.Errorf("bencode/write!: %v", err)
		}
		goVal, err := bencodeFromLG(vs[1])
		if err != nil {
			return vm.NIL, err
		}
		enc, err := bencode.EncodeBytes(goVal)
		if err != nil {
			return vm.NIL, fmt.Errorf("bencode/write!: %v", err)
		}
		if _, err := c.conn.Write(enc); err != nil {
			return vm.NIL, fmt.Errorf("bencode/write!: %v", err)
		}
		return vm.NIL, nil
	})

	readFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 || len(vs) > 2 {
			return vm.NIL, fmt.Errorf("bencode/read! expects 1-2 args (conn [timeout-ms])")
		}
		c, err := unboxNetConn(vs[0])
		if err != nil {
			return vm.NIL, fmt.Errorf("bencode/read!: %v", err)
		}
		timeoutMs := int64(-1)
		if len(vs) == 2 {
			t, ok := vs[1].(vm.Int)
			if !ok {
				return vm.NIL, fmt.Errorf("bencode/read! expected Int timeout-ms")
			}
			timeoutMs = int64(t)
		}

		// Serialize the stateful read path per conn: decoder init, deadline,
		// and Decode must not run concurrently on the same connection.
		c.readMu.Lock()
		defer c.readMu.Unlock()

		// Reuse the conn's decoder so bytes its bufio.Reader buffers past one
		// message survive to the next read.
		if c.dec == nil {
			c.dec = bencode.NewDecoder(c.conn)
		}
		if timeoutMs >= 0 {
			if err := c.conn.SetReadDeadline(time.Now().Add(time.Duration(timeoutMs) * time.Millisecond)); err != nil {
				return vm.NIL, fmt.Errorf("bencode/read!: %v", err)
			}
			defer func() { _ = c.conn.SetReadDeadline(time.Time{}) }()
		}

		var raw interface{}
		if err := c.dec.Decode(&raw); err != nil {
			if errors.Is(err, io.EOF) {
				return vm.NIL, nil
			}
			var ne net.Error
			if errors.As(err, &ne) && ne.Timeout() {
				return vm.NIL, fmt.Errorf("bencode/read! timeout after %dms", timeoutMs)
			}
			return vm.NIL, fmt.Errorf("bencode/read!: %v", err)
		}
		return bencodeToLG(raw)
	})

	ns := vm.NewNamespace("bencode")
	ns.Def("write!", writeFn)
	ns.Def("read!", readFn)
	RegisterNS(ns)
}
