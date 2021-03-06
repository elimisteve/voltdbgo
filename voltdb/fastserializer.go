package voltdb

import (
	"encoding/binary"
	"io"
)

// package private methods that perform voltdb compatible
// de/serialization on the base wire protocol types.
// See: http://community.voltdb.com/docs/WireProtocol/index

const (
	vt_ARRAY     int8 = -99 // array (short)(values*)
	vt_NULL      int8 = 1   // null
	vt_BOOL      int8 = 3   // boolean, byte
	vt_SHORT     int8 = 4   // int16
	vt_INT       int8 = 5   // int32
	vt_LONG      int8 = 6   // int64
	vt_FLOAT     int8 = 8   // float64
	vt_STRING    int8 = 9   // string (int32-length-prefix)(utf-8 bytes)
	vt_TIMESTAMP int8 = 11  // int64 timestamp microseconds
	vt_TABLE     int8 = 21  // VoltTable
	vt_DECIMAL   int8 = 22  // fix-scaled, fix-precision decimal
	vt_VARBIN    int8 = 25  // varbinary (int)(bytes)
)

var order = binary.BigEndian

// protoVersion is the implemented VoltDB wireprotocol version.
const protoVersion = 1

func writeProtoVersion(w io.Writer) error {
	var b [1]byte
	b[0] = protoVersion
	_, err := w.Write(b[:1])
	return err
}

func writeBoolean(w io.Writer, d bool) (err error) {
	if d {
		err = writeByte(w, 0x1)
	} else {
		err = writeByte(w, 0x0)
	}
	return
}

func readBoolean(r io.Reader) (bool, error) {
	val, err := readByte(r)
	if err != nil {
		return false, err
	}
	result := val != 0
	return result, nil
}

func writeByte(w io.Writer, d int8) error {
	var b [1]byte
	b[0] = byte(d)
	_, err := w.Write(b[:1])
	return err
}

func readByte(r io.Reader) (int8, error) {
	var b [1]byte
	bs := b[:1]
	_, err := r.Read(bs)
	if err != nil {
		return 0, err
	}
	return int8(b[0]), nil
}

func readByteArray(r io.Reader) ([]int8, error) {
	// byte arrays have 4 byte length prefixes.
	cnt, err := readInt(r)
	if err != nil {
		return nil, err
	}
	arr := make([]int8, cnt)
	for idx := range arr {
		val, err := readByte(r)
		if err != nil {
			return nil, err
		}
		arr[idx] = val
	}
	return arr, nil
}

func writeShort(w io.Writer, d int16) error {
	var b [2]byte
	bs := b[:2]
	order.PutUint16(bs, uint16(d))
	_, err := w.Write(bs)
	return err
}

func readShort(r io.Reader) (int16, error) {
	var b [2]byte
	bs := b[:2]
	_, err := r.Read(bs)
	if err != nil {
		return 0, err
	}
	result := order.Uint16(bs)
	return int16(result), nil
}

func writeInt(w io.Writer, d int32) error {
	var b [4]byte
	bs := b[:4]
	order.PutUint32(bs, uint32(d))
	_, err := w.Write(bs)
	return err
}

func readInt(r io.Reader) (int32, error) {
	var b [4]byte
	bs := b[:4]
	_, err := r.Read(bs)
	if err != nil {
		return 0, err
	}
	result := order.Uint32(bs)
	return int32(result), nil
}

func writeLong(w io.Writer, d int64) error {
	var b [8]byte
	bs := b[:8]
	order.PutUint64(bs, uint64(d))
	_, err := w.Write(bs)
	return err
}

func readLong(r io.Reader) (int64, error) {
	var b [8]byte
	bs := b[:8]
	_, err := r.Read(bs)
	if err != nil {
		return 0, err
	}
	result := order.Uint64(bs)
	return int64(result), nil
}

func writeFloat(w io.Writer, d float64) error {
	var b [8]byte
	bs := b[:8]
	order.PutUint64(bs, uint64(d))
	_, err := w.Write(bs)
	return err
}

func readFloat(r io.Reader) (float64, error) {
	var b [8]byte
	bs := b[:8]
	_, err := r.Read(bs)
	if err != nil {
		return 0, err
	}
	result := order.Uint64(bs)
	return float64(result), nil
}

func writeString(w io.Writer, d string) error {
	writeInt(w, int32(len(d)))
	_, err := io.WriteString(w, d)
	return err
}

func readString(r io.Reader) (result string, err error) {
	result = ""
	length, err := readInt(r)
	if err != nil {
		return
	}
	bs := make([]byte, length)
	_, err = r.Read(bs)
	if err != nil {
		return
	}
	return string(bs), nil
}

func readStringArray(r io.Reader) ([]string, error) {
	cnt, err := readShort(r)
	if err != nil {
		return nil, err
	}
	arr := make([]string, cnt)
	for idx := range arr {
		val, err := readString(r)
		if err != nil {
			return nil, err
		}
		arr[idx] = val
	}
	return arr, nil
}

func writeByteString(w io.Writer, d []byte) error {
	writeInt(w, int32(len(d)))
	_, err := w.Write(d)
	return err
}
