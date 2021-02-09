package zbits

import (
	"encoding/binary"
)

type endian int

const (
	// BigEndian means big-endian byte order
	BigEndian endian = iota
	// LittleEndian means little-endian byte order
	LittleEndian

	// BE means big-endian byte order
	BE = BigEndian
	// LE means little-endian byte order
	LE = LittleEndian
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
	case LittleEndian:
		binary.LittleEndian.PutUint16(buf[:], bits)
	}

	return buf[:]
}

// BytesToUint16 converts a slice of 2 bytes into a uint16
func BytesToUint16(data []byte, bitOrder endian) uint16 {
	var num uint16

	switch bitOrder {
	case BigEndian:
		num = binary.BigEndian.Uint16(data)
	case LittleEndian:
		num = binary.LittleEndian.Uint16(data)
	}

	return num
}

// BytesToUint32 converts a slice of 4 bytes into a uint32
func BytesToUint32(data []byte, bitOrder endian) uint32 {
	var num uint32

	switch bitOrder {
	case BigEndian:
		num = binary.BigEndian.Uint32(data)
	case LittleEndian:
		num = binary.LittleEndian.Uint32(data)
	}

	return num
}

// BytesToUint64 converts a slice of 8 bytes into a uint64
func BytesToUint64(data []byte, bitOrder endian) uint64 {
	var num uint64

	switch bitOrder {
	case BigEndian:
		num = binary.BigEndian.Uint64(data)
	case LittleEndian:
		num = binary.LittleEndian.Uint64(data)
	}

	return num
}
