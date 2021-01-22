package gif

import (
	"fmt"
	"math"
	"zarks/output"
	"zarks/zmath"
)

// MaxTableSize is the maximum code table size allowed by GIFs
const MaxTableSize = 4095 // 2^12 - 1

// Code is a code for the codestream to be encoded a few bits at a time
type Code uint16
type bit uint8
type key string

// SizedCode is a Code with a size marker
type SizedCode struct {
	code Code
	size int // bit size; max 12
}

type stream struct {
	codeTable map[key]Code
	codeLen   int // up to 12
	CC, EOI   key // clear code & end of information code keys
	colors    int

	buffer []byte
	data   []SizedCode
	bitLen int // size of the data, in bits
}

func newStream(colors int) *stream {
	// First, initialize the code table. The keys here actually match the color codes they correspond to,
	// but this will not be the case when additional codes are added.
	// The second returned value is the starting length, in bits, of a code.
	codeTable, codeLen := newCodeTable(colors)

	// Determine the indices of the 'Clear' and 'End of Information' codes
	var (
		cc  = colors     // 256
		eoi = colors + 1 // 257
	)

	data := make([]SizedCode, 1)
	data[0] = newSizedCode(codeTable[getKey(cc)], codeLen)

	// return the new stream
	return &stream{
		codeTable: codeTable,
		codeLen:   codeLen,
		CC:        getKey(cc),
		EOI:       getKey(eoi),
		colors:    colors,

		buffer: make([]byte, 0),

		data:   data,
		bitLen: codeLen, // bitLen starts out equal to codeLen because the code array always starts with one CC code
	}
}

func newCodeTable(colors int) (map[key]Code, int) {
	codeTable := make(map[key]Code)
	for i := 0; i < colors+2; i++ {
		bytes := Code(i)
		curKey := getKey(i)
		codeTable[curKey] = bytes
		//fmt.Println(strconv.FormatInt(int64(bytes[1]), 2))
	}
	codeLen := int(math.Log2(float64(colors))) + 1

	// Check the codetable for double-matches (shouldn't happen; just for debugging)
	for _, cod := range codeTable {
		matches := 0
		for _, match := range codeTable {
			if cod == match {
				matches++
			}
		}

		if matches != 1 {
			fmt.Printf("Found %v matches for the value %v\n", matches, cod)
		}
	}

	return codeTable, codeLen
}

func newSizedCode(c Code, size int) SizedCode {
	return SizedCode{
		code: c,
		size: size,
	}
}

func getCode(num int) Code {
	return Code(num)
}

func getBytes(num int) []byte {
	bits := output.Uint16ToBytes(uint16(num), output.BigEndian)
	bytes := []byte{
		bits[0],
		bits[1],
	}
	return bytes
}

func getKey(num int) key {
	// this key prefix is a difficult (impossible?) string of bits to achieve with the normal key-
	// generation formula. TODO: Make it better & faster
	var keyPrefix = []byte{
		255,
		255,
		0,
		255,
	}
	return key(append(keyPrefix, getBytes(num)...))
}

func (s *stream) addByte(K byte) {
	// if this is the first call for a data block, just initialize the buffer and wait for the next call
	if len(s.buffer) == 0 {
		s.buffer = append(s.buffer, K)
		//fmt.Println("Stream has received the first byte.")
		return
	}

	// Create the key for whose existence we must query the codetable.         <- i amaze myself sometimes
	// This can't be created with a simple append(), because we don't want to modify the buffer itself yet.
	testbuf := make([]byte, len(s.buffer)+1) // "buffer + K"
	copy(testbuf, s.buffer)
	testbuf[len(s.buffer)] = K
	testkey := key(testbuf)

	// Now we get to the fun stuff!
	_, ok := s.codeTable[testkey]
	tableSize := len(s.codeTable)
	if ok { // if the key already exists, append to the buffer and just continue on
		s.buffer = testbuf
		//fmt.Println("Code recognized at index ", len(s.data))

	} else if len(s.codeTable) < MaxTableSize { // if the codetable has room, write the buffer's (pre-existing) code
		// append to codetable
		s.codeTable[testkey] = Code(tableSize)
		s.codeLen = int(math.Log2(float64(tableSize-1))) + 1
		// Isaac you fucking genius I can't believe you tried tableSize-1 as the term jesus christ it fucking works

		if s.codeLen > 12 {
			panic("AAAAAHHHHH THE BITS ARE TOO LONG")
		}

		// append to data
		s.appendBuffer()

		// reset the buffer to the current byte, K
		s.buffer = make([]byte, 1)
		s.buffer[0] = K

	} else { // if the codetable needs to be reset, reset it!
		// Append to data before erasing codetable.
		s.appendBuffer()

		// Also append a CC code so the decoder knows the codetable is being reset.
		s.appendCC()

		// reset the buffer to the current byte, K
		s.buffer = make([]byte, 1)
		s.buffer[0] = K

		// erase all that hard work and start anew!
		s.codeTable, s.codeLen = newCodeTable(s.colors)

		//fmt.Println("Code table reset!")
	}
}

func (s *stream) appendBuffer() {
	var bufKey key
	if len(s.buffer) == 1 {
		bufKey = getKey(int(s.buffer[0]))
	} else {
		bufKey = key(s.buffer)
	}

	appendMe, ok := s.codeTable[bufKey]
	if ok {
		s.data = append(s.data, newSizedCode(appendMe, s.codeLen))
		s.bitLen += s.codeLen
	} else {
		fmt.Print("Current buffer: ")
		for _, b := range s.buffer {
			fmt.Print(int(b), " ")
		}
		fmt.Println(" (Length of ", len(s.buffer), ")")

		fmt.Println("Code not found for: ", string(s.buffer))
		fmt.Println("Current table size: ", len(s.codeTable))
		panic("^ what they said")
	}
}

func (s *stream) appendCC() {
	cc, ok := s.codeTable[s.CC]
	if ok {
		s.data = append(s.data, newSizedCode(cc, s.codeLen))
		s.bitLen += s.codeLen
		//fmt.Println("CC Code added at index ", len(s.data)-1)
	} else {
		panic("No CC Code???????")
	}
}

// Finalize puts the last byte(s) of the buffer into the codestream and converts the codestream to bytes.
// Note on 1/8/2021: this is extremely mentally taxing
func (s *stream) finalize() []byte {
	var (
		data       = make([]byte, 0)
		bitsFilled int
		buf        byte
	)

	// first, add the last of the buffer to the codestream
	s.appendBuffer()
	s.data = append(s.data, newSizedCode(s.codeTable[s.EOI], s.codeLen)) // EOI IS IMPORTANT!

	for i, c := range s.data {
		if c.code == s.codeTable[s.CC] {
			//fmt.Println("CC Code encountered at index ", i)
		}
		var (
			bits      = c.code
			leftovers = c.size
			used      = 0
		)

		for leftovers > 0 {
			shiftBy := bitsFilled - used // shiftBy to the LEFT

			var temp byte
			if shiftBy > 0 {
				temp = byte(bits << shiftBy)
			} else {
				temp = byte(bits >> -shiftBy)
			}
			buf = buf | temp

			dUsed := zmath.MinInt(leftovers, 8-bitsFilled)
			used += dUsed
			bitsFilled += dUsed
			leftovers -= dUsed

			if bitsFilled == 8 {
				data = append(data, buf)
				//fmt.Println(strconv.FormatInt(int64(buf), 2))
				buf = 0
				bitsFilled = 0
			} else if bitsFilled > 8 {
				panic("AAAAAAHHHH TOO MANY BITS!!!!")
			}
		}

		// if there's some buffer left over on the last one, just add it
		if i == (len(s.data)-1) && bitsFilled != 0 {
			data = append(data, buf)
		}
	}

	return data
}

/*
// SHIFT THEM BITS
for leftovers > 0 {
	shiftBy := leftovers + bitsFilled - 8
	if shiftBy > 0 {
		buf = buf | byte(bits>>shiftBy)
	} else {
		buf = buf | byte(bits<<(-shiftBy))
	}

	dBits := zmath.MinInt(8, bitsFilled+leftovers) - bitsFilled
	bitsFilled += dBits
	leftovers -= dBits

	if bitsFilled == 8 {
		data = append(data, buf)
		buf = 0
		bitsFilled = 0
	} else if bitsFilled > 8 {
		panic("AAAAAAHHHH TOO MANY BITS!!!!")
	}
}
*/
