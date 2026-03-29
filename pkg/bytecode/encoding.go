package bytecode

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math"
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

// WriteInt32 writes a little-endian int32.
func (w *Writer) WriteInt32(v int32) error {
	binary.LittleEndian.PutUint32(w.buf[:4], uint32(v))
	_, err := w.w.Write(w.buf[:4])
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
}

// NewReader creates a Reader wrapping r.
func NewReader(r io.Reader) *Reader {
	if br, ok := r.(*bufio.Reader); ok {
		return &Reader{r: br}
	}
	return &Reader{r: bufio.NewReader(r)}
}

// ReadVarint reads an unsigned LEB128-encoded integer.
func (r *Reader) ReadVarint() (uint64, error) {
	var result uint64
	var shift uint
	for i := 0; i < 10; i++ {
		b, err := r.r.ReadByte()
		if err != nil {
			return 0, err
		}
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
	return r.r.ReadByte()
}

// ReadBytes reads exactly n bytes.
func (r *Reader) ReadBytes(n int) ([]byte, error) {
	buf := make([]byte, n)
	_, err := io.ReadFull(r.r, buf)
	return buf, err
}

// ReadUint16 reads a little-endian uint16.
func (r *Reader) ReadUint16() (uint16, error) {
	if _, err := io.ReadFull(r.r, r.buf[:2]); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(r.buf[:2]), nil
}

// ReadInt32 reads a little-endian int32.
func (r *Reader) ReadInt32() (int32, error) {
	if _, err := io.ReadFull(r.r, r.buf[:4]); err != nil {
		return 0, err
	}
	return int32(binary.LittleEndian.Uint32(r.buf[:4])), nil
}

// ReadFloat64 reads an IEEE 754 little-endian float64.
func (r *Reader) ReadFloat64() (float64, error) {
	if _, err := io.ReadFull(r.r, r.buf[:8]); err != nil {
		return 0, err
	}
	return math.Float64frombits(binary.LittleEndian.Uint64(r.buf[:8])), nil
}
