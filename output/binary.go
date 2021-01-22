package output

import (
	"encoding/binary"
)

type endian int

// Endian bit orders
const (
	BigEndian endian = iota
	LittleEndian
)

// Uint64ToBytes converts a float64 to an array of eight bytes
func Uint64ToBytes(bits uint64, bitOrder endian) []byte {
	buf := [8]byte{}
	switch bitOrder {
	case BigEndian:
		binary.BigEndian.PutUint64(buf[:], bits)
		break
	case LittleEndian:
		binary.LittleEndian.PutUint64(buf[:], bits)
		break
	}

	return buf[:]
}

// Uint32ToBytes converts a float32 to an array of four bytes
func Uint32ToBytes(bits uint32, bitOrder endian) []byte {
	buf := [4]byte{}

	switch bitOrder {
	case BigEndian:
		binary.BigEndian.PutUint32(buf[:], bits)
		break
	case LittleEndian:
		binary.LittleEndian.PutUint32(buf[:], bits)
		break
	}

	return buf[:]
}

// Uint16ToBytes converts a float32 to an array of two bytes
func Uint16ToBytes(bits uint16, bitOrder endian) []byte {
	buf := [2]byte{}

	switch bitOrder {
	case BigEndian:
		binary.BigEndian.PutUint16(buf[:], bits)
		break
	case LittleEndian:
		binary.LittleEndian.PutUint16(buf[:], bits)
		break
	}

	return buf[:]
}
