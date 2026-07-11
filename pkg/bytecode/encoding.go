package bytecode

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"unsafe"
)

// Writer wraps a buffered writer with binary encoding helpers.
type Writer struct {
	w   *bufio.Writer
	buf [8]byte // scratch buffer for fixed-size writes
}

// NewWriter creates a Writer wrapping w.
func NewWriter(w io.Writer) *Writer {
	if bw, ok := w.(*bufio.Writer); ok {
		return &Writer{w: bw}
	}
	return &Writer{w: bufio.NewWriter(w)}
}

// Flush flushes the underlying buffer.
func (w *Writer) Flush() error {
	return w.w.Flush()
}

// WriteVarint writes an unsigned LEB128-encoded integer.
func (w *Writer) WriteVarint(v uint64) error {
	for {
		b := byte(v & 0x7f)
		v >>= 7
		if v != 0 {
			b |= 0x80
		}
		if err := w.w.WriteByte(b); err != nil {
			return err
		}
		if v == 0 {
			return nil
		}
	}
}

// WriteSvarint writes a signed zigzag-encoded varint.
func (w *Writer) WriteSvarint(v int64) error {
	uv := uint64(v<<1) ^ uint64(v>>63) // zigzag encode
	return w.WriteVarint(uv)
}

// WriteByte writes a single byte.
func (w *Writer) WriteByte(b byte) error {
	return w.w.WriteByte(b)
}

// WriteBytes writes a byte slice.
func (w *Writer) WriteBytes(b []byte) error {
	_, err := w.w.Write(b)
	return err
}

// WriteUint16 writes a little-endian uint16.
func (w *Writer) WriteUint16(v uint16) error {
	binary.LittleEndian.PutUint16(w.buf[:2], v)
	_, err := w.w.Write(w.buf[:2])
	return err
}

// WriteUint32 writes a little-endian uint32.
func (w *Writer) WriteUint32(v uint32) error {
	binary.LittleEndian.PutUint32(w.buf[:4], v)
	_, err := w.w.Write(w.buf[:4])
	return err
}

// WriteInt32 writes a little-endian int32.
func (w *Writer) WriteInt32(v int32) error {
	binary.LittleEndian.PutUint32(w.buf[:4], uint32(v))
	_, err := w.w.Write(w.buf[:4])
	return err
}

// WriteUint64 writes a little-endian uint64.
func (w *Writer) WriteUint64(v uint64) error {
	binary.LittleEndian.PutUint64(w.buf[:8], v)
	_, err := w.w.Write(w.buf[:8])
	return err
}

// WriteFloat64 writes an IEEE 754 little-endian float64.
func (w *Writer) WriteFloat64(v float64) error {
	binary.LittleEndian.PutUint64(w.buf[:8], math.Float64bits(v))
	_, err := w.w.Write(w.buf[:8])
	return err
}

// Reader wraps a buffered reader with binary decoding helpers.
type Reader struct {
	r   *bufio.Reader
	buf [8]byte
	// pos is the absolute byte offset consumed so far. Every read method
	// advances it, so Offset() is exact regardless of bufio buffering. It lets
	// the decoder capture a section's raw byte span for deferred parsing.
	pos int
	// data, when non-nil, is the full backing buffer the reader was created over
	// (Reader was built from a []byte). It enables zero-copy Slice(start,end)
	// of already-resident bytes — used to defer source-map materialization.
	data []byte
}

// NewReader creates a Reader wrapping r.
func NewReader(r io.Reader) *Reader {
	if br, ok := r.(*bufio.Reader); ok {
		return &Reader{r: br}
	}
	return &Reader{r: bufio.NewReader(r)}
}

// NewReaderBytes creates a Reader over an in-memory buffer. Reads still stream
// through bufio, but the reader retains `data` so the decoder can Slice() a
// section's raw bytes zero-copy (the buffer stays resident) and defer decoding.
func NewReaderBytes(data []byte) *Reader {
	return &Reader{r: bufio.NewReader(bytes.NewReader(data)), data: data}
}

// Offset returns the absolute number of bytes consumed so far.
func (r *Reader) Offset() int { return r.pos }

// HasBackingData reports whether the reader can Slice() raw bytes (built via
// NewReaderBytes).
func (r *Reader) HasBackingData() bool { return r.data != nil }

// Slice returns data[start:end] of the backing buffer (zero-copy). Only valid
// when HasBackingData() is true and [start,end] is within already-consumed input.
func (r *Reader) Slice(start, end int) []byte { return r.data[start:end] }

// ReadVarint reads an unsigned LEB128-encoded integer.
func (r *Reader) ReadVarint() (uint64, error) {
	var result uint64
	var shift uint
	for range 10 {
		b, err := r.r.ReadByte()
		if err != nil {
			return 0, err
		}
		r.pos++
		result |= uint64(b&0x7f) << shift
		if b&0x80 == 0 {
			return result, nil
		}
		shift += 7
	}
	return 0, fmt.Errorf("varint exceeds 10 bytes")
}

// ReadSvarint reads a signed zigzag-encoded varint.
func (r *Reader) ReadSvarint() (int64, error) {
	uv, err := r.ReadVarint()
	if err != nil {
		return 0, err
	}
	// zigzag decode
	return int64(uv>>1) ^ -int64(uv&1), nil
}

// ReadByte reads a single byte.
func (r *Reader) ReadByte() (byte, error) {
	b, err := r.r.ReadByte()
	if err == nil {
		r.pos++
	}
	return b, err
}

// ReadBytes reads exactly n bytes.
func (r *Reader) ReadBytes(n int) ([]byte, error) {
	buf := make([]byte, n)
	m, err := io.ReadFull(r.r, buf)
	r.pos += m
	return buf, err
}

// ReadString reads exactly n bytes and returns them as a string with a single
// backing allocation. The returned string owns its storage and is safe to
// retain after the reader advances.
func (r *Reader) ReadString(n int) (string, error) {
	buf := make([]byte, n)
	m, err := io.ReadFull(r.r, buf)
	r.pos += m
	if err != nil {
		return "", err
	}
	if len(buf) == 0 {
		return "", nil
	}
	return unsafe.String(unsafe.SliceData(buf), len(buf)), nil
}

// ReadUint16 reads a little-endian uint16.
func (r *Reader) ReadUint16() (uint16, error) {
	if _, err := io.ReadFull(r.r, r.buf[:2]); err != nil {
		return 0, err
	}
	r.pos += 2
	return binary.LittleEndian.Uint16(r.buf[:2]), nil
}

// ReadInt32 reads a little-endian int32.
func (r *Reader) ReadInt32() (int32, error) {
	if _, err := io.ReadFull(r.r, r.buf[:4]); err != nil {
		return 0, err
	}
	r.pos += 4
	return int32(binary.LittleEndian.Uint32(r.buf[:4])), nil
}

// ReadUint32 reads a little-endian uint32.
func (r *Reader) ReadUint32() (uint32, error) {
	if _, err := io.ReadFull(r.r, r.buf[:4]); err != nil {
		return 0, err
	}
	r.pos += 4
	return binary.LittleEndian.Uint32(r.buf[:4]), nil
}

// ReadUint64 reads a little-endian uint64.
func (r *Reader) ReadUint64() (uint64, error) {
	if _, err := io.ReadFull(r.r, r.buf[:8]); err != nil {
		return 0, err
	}
	r.pos += 8
	return binary.LittleEndian.Uint64(r.buf[:8]), nil
}

// ReadFloat64 reads an IEEE 754 little-endian float64.
func (r *Reader) ReadFloat64() (float64, error) {
	if _, err := io.ReadFull(r.r, r.buf[:8]); err != nil {
		return 0, err
	}
	r.pos += 8
	return math.Float64frombits(binary.LittleEndian.Uint64(r.buf[:8])), nil
}
