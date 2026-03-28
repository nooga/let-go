/*
 * io namespace — let-go's equivalent of clojure.java.io
 *
 * Built on Go's io.Reader/io.Writer interfaces.
 * Provides: reader/writer coercion, line-seq, byte buffers,
 * copy, encoding (base64/hex/url), and slurp/spit polymorphism.
 */

package rt

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/nooga/let-go/pkg/vm"
)

// --- Reader/Writer wrapper types ---

// LGReader wraps any io.Reader with optional buffering and close.
type LGReader struct {
	raw    io.Reader
	br     *bufio.Reader
	closer io.Closer // non-nil if the underlying source is closeable
}

func newLGReader(r io.Reader, closer io.Closer) *LGReader {
	return &LGReader{raw: r, closer: closer}
}

func (r *LGReader) Buffered() *bufio.Reader {
	if r.br == nil {
		r.br = bufio.NewReader(r.raw)
	}
	return r.br
}

func (r *LGReader) Read(p []byte) (int, error) {
	return r.Buffered().Read(p)
}

func (r *LGReader) Close() error {
	if r.closer != nil {
		return r.closer.Close()
	}
	return nil
}

func (r *LGReader) String() string {
	return "#<io/reader>"
}

// LGWriter wraps any io.Writer with optional buffering and close.
type LGWriter struct {
	raw    io.Writer
	bw     *bufio.Writer
	closer io.Closer
}

func newLGWriter(w io.Writer, closer io.Closer) *LGWriter {
	return &LGWriter{raw: w, closer: closer}
}

func (w *LGWriter) Buffered() *bufio.Writer {
	if w.bw == nil {
		w.bw = bufio.NewWriter(w.raw)
	}
	return w.bw
}

func (w *LGWriter) Write(p []byte) (int, error) {
	return w.raw.Write(p)
}

func (w *LGWriter) Flush() error {
	if w.bw != nil {
		return w.bw.Flush()
	}
	return nil
}

func (w *LGWriter) Close() error {
	if w.bw != nil {
		w.bw.Flush()
	}
	if w.closer != nil {
		return w.closer.Close()
	}
	return nil
}

func (w *LGWriter) String() string {
	return "#<io/writer>"
}

// LGBuffer wraps bytes.Buffer as a readable/writable value.
type LGBuffer struct {
	Buf *bytes.Buffer
}

func newLGBuffer(initial []byte) *LGBuffer {
	if initial != nil {
		return &LGBuffer{Buf: bytes.NewBuffer(initial)}
	}
	return &LGBuffer{Buf: &bytes.Buffer{}}
}

func (b *LGBuffer) Read(p []byte) (int, error)  { return b.Buf.Read(p) }
func (b *LGBuffer) Write(p []byte) (int, error) { return b.Buf.Write(p) }
func (b *LGBuffer) String() string              { return fmt.Sprintf("#<io/buffer %d bytes>", b.Buf.Len()) }

// rawStr extracts a plain Go string from a vm.Value (without quoting).
func rawStr(v vm.Value) string {
	if v == nil || v == vm.NIL {
		return ""
	}
	if s, ok := v.(vm.String); ok {
		return string(s)
	}
	return v.String()
}

// --- Protocols ---

// IReadable and IWritable are the IO coercion protocols.
// Extending these allows user-defined types to participate in the IO system.
var (
	ReadableProto *vm.Protocol // method: -make-reader [x] → reader
	WritableProto *vm.Protocol // method: -make-writer [x] → writer
)

// --- Coercion helpers ---

// coerceReader extracts an io.Reader from various let-go value types.
// Checks protocol first, then falls back to built-in Go type dispatch.
func coerceReader(v vm.Value) (*LGReader, error) {
	// Protocol dispatch: if the value satisfies IReadable, use it
	if ReadableProto != nil && ReadableProto.Satisfies(v) {
		fn, ok := ReadableProto.Lookup("make-reader", v)
		if ok {
			result, err := fn.Invoke([]vm.Value{v})
			if err != nil {
				return nil, err
			}
			// The protocol method should return a boxed LGReader
			if b, ok := result.(*vm.Boxed); ok {
				if r, ok := b.Unbox().(*LGReader); ok {
					return r, nil
				}
			}
			// If it returned something else, try to coerce that (one level only)
			return coerceReaderBuiltin(result)
		}
	}
	return coerceReaderBuiltin(v)
}

// coerceReaderBuiltin handles built-in Go types without protocol dispatch.
func coerceReaderBuiltin(v vm.Value) (*LGReader, error) {
	// Already an LGReader?
	if b, ok := v.(*vm.Boxed); ok {
		if r, ok := b.Unbox().(*LGReader); ok {
			return r, nil
		}
		if buf, ok := b.Unbox().(*LGBuffer); ok {
			return newLGReader(buf, nil), nil
		}
		if h, ok := b.Unbox().(*IOHandle); ok {
			return newLGReader(h.Reader(), nil), nil
		}
		// Any boxed io.Reader
		if r, ok := b.Unbox().(io.Reader); ok {
			var closer io.Closer
			if c, ok := b.Unbox().(io.Closer); ok {
				closer = c
			}
			return newLGReader(r, closer), nil
		}
	}
	// String → file path
	if s, ok := v.(vm.String); ok {
		f, err := os.Open(string(s))
		if err != nil {
			return nil, err
		}
		return newLGReader(f, f), nil
	}
	return nil, fmt.Errorf("cannot coerce %s to reader", v.Type().Name())
}

// coerceWriter extracts an io.Writer from various let-go value types.
func coerceWriter(v vm.Value) (*LGWriter, error) {
	// Protocol dispatch
	if WritableProto != nil && WritableProto.Satisfies(v) {
		fn, ok := WritableProto.Lookup("make-writer", v)
		if ok {
			result, err := fn.Invoke([]vm.Value{v})
			if err != nil {
				return nil, err
			}
			if b, ok := result.(*vm.Boxed); ok {
				if w, ok := b.Unbox().(*LGWriter); ok {
					return w, nil
				}
			}
			return coerceWriterBuiltin(result)
		}
	}
	return coerceWriterBuiltin(v)
}

// coerceWriterBuiltin handles built-in Go types without protocol dispatch.
func coerceWriterBuiltin(v vm.Value) (*LGWriter, error) {
	if b, ok := v.(*vm.Boxed); ok {
		if w, ok := b.Unbox().(*LGWriter); ok {
			return w, nil
		}
		if buf, ok := b.Unbox().(*LGBuffer); ok {
			return newLGWriter(buf, nil), nil
		}
		if h, ok := b.Unbox().(*IOHandle); ok {
			return newLGWriter(h.File, nil), nil
		}
		if w, ok := b.Unbox().(io.Writer); ok {
			var closer io.Closer
			if c, ok := b.Unbox().(io.Closer); ok {
				closer = c
			}
			return newLGWriter(w, closer), nil
		}
	}
	// String → file path (write mode)
	if s, ok := v.(vm.String); ok {
		f, err := os.OpenFile(string(s), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return nil, err
		}
		return newLGWriter(f, f), nil
	}
	return nil, fmt.Errorf("cannot coerce %s to writer", v.Type().Name())
}

// coerceResponseBody reads an HTTP response body from various value types.
// Supports: String, LGReader, LGBuffer, IOHandle, Seq of strings.
func coerceResponseBody(body vm.Value) ([]byte, error) {
	if body == vm.NIL {
		return nil, nil
	}
	// Fast path: string
	if s, ok := body.(vm.String); ok {
		return []byte(s), nil
	}
	// Reader-coercible: read all
	if r, err := coerceReaderBuiltin(body); err == nil {
		defer r.Close()
		return io.ReadAll(r)
	}
	// Seq of strings
	if sq, ok := body.(vm.Sequable); ok {
		var buf bytes.Buffer
		for s := sq.Seq(); s != nil; s = s.Next() {
			v := s.First()
			if str, ok := v.(vm.String); ok {
				buf.WriteString(string(str))
			} else {
				buf.WriteString(v.String())
			}
		}
		return buf.Bytes(), nil
	}
	// Fallback
	return []byte(body.String()), nil
}

// --- Lazy line-seq from Go ---

func makeLineSeq(br *bufio.Reader) *vm.LazySeq {
	var buildThunk func() vm.Fn
	buildThunk = func() vm.Fn {
		fn, _ := vm.NativeFnType.Wrap(func(_ []vm.Value) (vm.Value, error) {
			line, err := br.ReadString('\n')
			if err != nil {
				if len(line) > 0 {
					// Last line without trailing newline
					line = strings.TrimRight(line, "\n\r")
					return vm.NewCons(vm.String(line), vm.EmptyList), nil
				}
				return nil, nil // EOF → empty seq
			}
			line = strings.TrimRight(line, "\n\r")
			return vm.NewCons(vm.String(line), vm.NewLazySeq(buildThunk())), nil
		})
		return fn.(vm.Fn)
	}
	return vm.NewLazySeq(buildThunk())
}

// helper to build a protocol impl map from Go
func protoImplMap(method vm.Keyword, fn vm.Value) *vm.PersistentMap {
	return vm.EmptyPersistentMap.Assoc(method, fn).(*vm.PersistentMap)
}

// nolint
func installIoNS() {
	// --- Protocols ---
	ReadableProto = vm.NewProtocol("io/IReadable", []vm.Symbol{"make-reader"})
	WritableProto = vm.NewProtocol("io/IWritable", []vm.Symbol{"make-writer"})

	// Extend IReadable for String (path → file reader)
	stringReaderImpl, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		s := vs[0].(vm.String)
		f, err := os.Open(string(s))
		if err != nil {
			return vm.NIL, err
		}
		return vm.NewBoxed(newLGReader(f, f)), nil
	})
	ReadableProto.Extend(vm.StringType, protoImplMap("make-reader", stringReaderImpl))

	// Extend IWritable for String (path → file writer)
	stringWriterImpl, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		s := vs[0].(vm.String)
		f, err := os.OpenFile(string(s), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return vm.NIL, err
		}
		return vm.NewBoxed(newLGWriter(f, f)), nil
	})
	WritableProto.Extend(vm.StringType, protoImplMap("make-writer", stringWriterImpl))

	// Extend IReadable for io/URL records — HTTP GET the URL
	urlReaderImpl, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		rec, ok := vs[0].(*vm.Record)
		if !ok {
			return vm.NIL, fmt.Errorf("make-reader expected URL record")
		}
		rawURL := rec.ValueAt(vm.Keyword("raw"))
		if rawURL == vm.NIL {
			return vm.NIL, fmt.Errorf("URL has no :raw field")
		}
		var urlStr string
		if s, ok := rawURL.(vm.String); ok {
			urlStr = string(s)
		} else {
			urlStr = rawURL.String()
		}
		resp, err := http.Get(urlStr)
		if err != nil {
			return vm.NIL, err
		}
		return vm.NewBoxed(newLGReader(resp.Body, resp.Body)), nil
	})
	ReadableProto.Extend(urlMapping.RecType, protoImplMap("make-reader", urlReaderImpl))

	// io/reader — coerce to reader
	reader, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("io/reader expects 1 arg")
		}
		r, err := coerceReader(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		return vm.NewBoxed(r), nil
	})

	// io/writer — coerce to writer
	writer, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("io/writer expects 1 arg")
		}
		w, err := coerceWriter(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		return vm.NewBoxed(w), nil
	})

	// io/line-seq — lazy seq of lines from a reader
	lineSeq, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("io/line-seq expects 1 arg")
		}
		r, err := coerceReader(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		return makeLineSeq(r.Buffered()), nil
	})

	// io/close — close a reader, writer, or handle
	closef, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("io/close expects 1 arg")
		}
		b, ok := vs[0].(*vm.Boxed)
		if !ok {
			return vm.NIL, fmt.Errorf("io/close expected closeable value")
		}
		if c, ok := b.Unbox().(io.Closer); ok {
			return vm.NIL, c.Close()
		}
		return vm.NIL, nil
	})

	// io/copy — stream from reader to writer, returns bytes copied
	copyf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("io/copy expects 2 args (reader, writer)")
		}
		r, err := coerceReader(vs[0])
		if err != nil {
			return vm.NIL, fmt.Errorf("io/copy reader: %w", err)
		}
		w, err := coerceWriter(vs[1])
		if err != nil {
			return vm.NIL, fmt.Errorf("io/copy writer: %w", err)
		}
		n, err := io.Copy(w, r)
		if err != nil {
			return vm.NIL, err
		}
		return vm.Int(n), nil
	})

	// io/buffer — create a byte buffer
	buffer, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) == 0 {
			return vm.NewBoxed(newLGBuffer(nil)), nil
		}
		if s, ok := vs[0].(vm.String); ok {
			return vm.NewBoxed(newLGBuffer([]byte(s))), nil
		}
		return vm.NewBoxed(newLGBuffer(nil)), nil
	})

	// io/buffer-str — get buffer contents as string
	bufferStr, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("io/buffer-str expects 1 arg")
		}
		b, ok := vs[0].(*vm.Boxed)
		if !ok {
			return vm.NIL, fmt.Errorf("io/buffer-str expected buffer")
		}
		buf, ok := b.Unbox().(*LGBuffer)
		if !ok {
			return vm.NIL, fmt.Errorf("io/buffer-str expected buffer")
		}
		return vm.String(buf.Buf.String()), nil
	})

	// io/buffer-bytes — get buffer length
	bufferLen, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("io/buffer-len expects 1 arg")
		}
		b, ok := vs[0].(*vm.Boxed)
		if !ok {
			return vm.NIL, fmt.Errorf("io/buffer-len expected buffer")
		}
		buf, ok := b.Unbox().(*LGBuffer)
		if !ok {
			return vm.NIL, fmt.Errorf("io/buffer-len expected buffer")
		}
		return vm.Int(buf.Buf.Len()), nil
	})

	// io/slurp — read everything from a reader-coercible source
	slurpf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("io/slurp expects 1 arg")
		}
		r, err := coerceReader(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		defer r.Close()
		data, err := io.ReadAll(r)
		if err != nil {
			return vm.NIL, err
		}
		return vm.String(data), nil
	})

	// io/spit — write string to a writer-coercible target
	spitf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("io/spit expects 2 args")
		}
		w, err := coerceWriter(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		defer w.Close()
		var data string
		if s, ok := vs[1].(vm.String); ok {
			data = string(s)
		} else {
			data = vs[1].String()
		}
		_, err = w.Write([]byte(data))
		return vm.NIL, err
	})

	// io/write-str — write a string to a writer (no close)
	writeStr, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("io/write-str expects 2 args (writer, str)")
		}
		w, err := coerceWriter(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		var data string
		if s, ok := vs[1].(vm.String); ok {
			data = string(s)
		} else {
			data = vs[1].String()
		}
		_, err = w.Write([]byte(data))
		return vm.NIL, err
	})

	// io/flush — flush a writer
	flushf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("io/flush expects 1 arg")
		}
		w, err := coerceWriter(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		return vm.NIL, w.Flush()
	})

	// io/string-reader — create a reader from a string (reads string contents, not a file path)
	stringReader, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("io/string-reader expects 1 arg")
		}
		var s string
		if str, ok := vs[0].(vm.String); ok {
			s = string(str)
		} else {
			s = vs[0].String()
		}
		return vm.NewBoxed(newLGReader(strings.NewReader(s), nil)), nil
	})

	// io/string-writer — create a writer backed by a buffer, retrieve with buffer-str
	stringWriter, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		buf := newLGBuffer(nil)
		return vm.NewBoxed(newLGWriter(buf, nil)), nil
	})

	// --- Encoding functions ---

	// io/encode — (io/encode :base64 data) or (io/encode :hex data) or (io/encode :url data)
	encode, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("io/encode expects 2 args (encoding, data)")
		}
		enc, ok := vs[0].(vm.Keyword)
		if !ok {
			return vm.NIL, fmt.Errorf("io/encode expected keyword encoding")
		}
		var data string
		if s, ok := vs[1].(vm.String); ok {
			data = string(s)
		} else {
			data = vs[1].String()
		}
		switch enc {
		case "base64":
			return vm.String(base64.StdEncoding.EncodeToString([]byte(data))), nil
		case "base64url":
			return vm.String(base64.URLEncoding.EncodeToString([]byte(data))), nil
		case "hex":
			return vm.String(hex.EncodeToString([]byte(data))), nil
		case "url":
			return vm.String(url.QueryEscape(data)), nil
		default:
			return vm.NIL, fmt.Errorf("io/encode: unknown encoding %s", enc)
		}
	})

	// io/decode — (io/decode :base64 data) etc.
	decode, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("io/decode expects 2 args (encoding, data)")
		}
		enc, ok := vs[0].(vm.Keyword)
		if !ok {
			return vm.NIL, fmt.Errorf("io/decode expected keyword encoding")
		}
		var data string
		if s, ok := vs[1].(vm.String); ok {
			data = string(s)
		} else {
			data = vs[1].String()
		}
		switch enc {
		case "base64":
			b, err := base64.StdEncoding.DecodeString(data)
			if err != nil {
				return vm.NIL, err
			}
			return vm.String(b), nil
		case "base64url":
			b, err := base64.URLEncoding.DecodeString(data)
			if err != nil {
				return vm.NIL, err
			}
			return vm.String(b), nil
		case "hex":
			b, err := hex.DecodeString(data)
			if err != nil {
				return vm.NIL, err
			}
			return vm.String(b), nil
		case "url":
			s, err := url.QueryUnescape(data)
			if err != nil {
				return vm.NIL, err
			}
			return vm.String(s), nil
		default:
			return vm.NIL, fmt.Errorf("io/decode: unknown encoding %s", enc)
		}
	})

	// --- read-lines / write-lines (port from io.lg) ---

	readLines, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("io/read-lines expects 1 arg")
		}
		r, err := coerceReader(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		defer r.Close()
		var lines []vm.Value
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			lines = append(lines, vm.String(scanner.Text()))
		}
		if err := scanner.Err(); err != nil {
			return vm.NIL, err
		}
		return vm.NewArrayVector(lines), nil
	})

	writeLines, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("io/write-lines expects 2 args (target, lines)")
		}
		w, err := coerceWriter(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		defer w.Close()
		seq, ok := vs[1].(vm.Sequable)
		if !ok {
			return vm.NIL, fmt.Errorf("io/write-lines expected sequable lines")
		}
		first := true
		for s := seq.Seq(); s != nil; s = s.Next() {
			if !first {
				w.Write([]byte("\n"))
			}
			val := s.First()
			if str, ok := val.(vm.String); ok {
				w.Write([]byte(str))
			} else {
				w.Write([]byte(val.String()))
			}
			first = false
		}
		return vm.NIL, nil
	})

	ns := vm.NewNamespace("io")

	// Protocols
	ns.Def("IReadable", ReadableProto)
	ns.Def("IWritable", WritableProto)
	ns.Def("make-reader", vm.NewProtocolFn(ReadableProto, "make-reader"))
	ns.Def("make-writer", vm.NewProtocolFn(WritableProto, "make-writer"))

	// Core reader/writer
	ns.Def("reader", reader)
	ns.Def("writer", writer)
	ns.Def("string-reader", stringReader)
	ns.Def("string-writer", stringWriter)
	ns.Def("close", closef)
	ns.Def("flush", flushf)

	// Streaming
	ns.Def("line-seq", lineSeq)
	ns.Def("copy", copyf)
	ns.Def("slurp", slurpf)
	ns.Def("spit", spitf)
	ns.Def("write-str", writeStr)

	// Byte buffers
	ns.Def("buffer", buffer)
	ns.Def("buffer-str", bufferStr)
	ns.Def("buffer-len", bufferLen)

	// Encoding
	ns.Def("encode", encode)
	ns.Def("decode", decode)

	// File helpers
	ns.Def("read-lines", readLines)
	ns.Def("write-lines", writeLines)

	// --- URL ---

	// io/url — parse a URL string into a URL record
	urlf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("io/url expects 1 arg")
		}
		var raw string
		if s, ok := vs[0].(vm.String); ok {
			raw = string(s)
		} else {
			raw = vs[0].String()
		}
		u, err := url.Parse(raw)
		if err != nil {
			return vm.NIL, err
		}
		var userInfo string
		if u.User != nil {
			userInfo = u.User.String()
		}
		return urlMapping.StructToRecord(LGURL{
			Scheme:   u.Scheme,
			Host:     u.Hostname(),
			Port:     u.Port(),
			Path:     u.Path,
			Query:    u.RawQuery,
			Fragment: u.Fragment,
			UserInfo: userInfo,
			Raw:      raw,
		}), nil
	})

	// io/url? — check if value is a URL record
	urlPred, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.FALSE, nil
		}
		if r, ok := vs[0].(*vm.Record); ok {
			return vm.Boolean(r.Type() == urlMapping.RecType), nil
		}
		return vm.FALSE, nil
	})

	// io/url-str — format a URL record back to a string
	urlStr, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("io/url-str expects 1 arg")
		}
		rec, ok := vs[0].(*vm.Record)
		if !ok {
			return vm.NIL, fmt.Errorf("io/url-str expected URL record")
		}
		u := &url.URL{
			Scheme:   rawStr(rec.ValueAt(vm.Keyword("scheme"))),
			Host:     rawStr(rec.ValueAt(vm.Keyword("host"))),
			Path:     rawStr(rec.ValueAt(vm.Keyword("path"))),
			RawQuery: rawStr(rec.ValueAt(vm.Keyword("query"))),
			Fragment: rawStr(rec.ValueAt(vm.Keyword("fragment"))),
		}
		if port := rawStr(rec.ValueAt(vm.Keyword("port"))); port != "" {
			u.Host = u.Host + ":" + port
		}
		if ui := rawStr(rec.ValueAt(vm.Keyword("user-info"))); ui != "" {
			u.User = url.User(ui)
		}
		return vm.String(u.String()), nil
	})

	ns.Def("url", urlf)
	ns.Def("url?", urlPred)
	ns.Def("url-str", urlStr)

	RegisterNS(ns)
}
