package zbits

// Stream is an array of bytes with some helpful member functions
type Stream struct {
	data []byte
}

// NewStream returns a new stream of the desired size
func NewStream(size int) *Stream {
	return &Stream{
		data: make([]byte, size),
	}
}

// Size returns the length of the called Stream's data
func (s *Stream) Size() int {
	return len(s.data)
}

// Append appends to the stream
func (s *Stream) Append(b byte) {
	s.data = append(s.data, b)
}

// AppendX appends several bytes to the desired
func (s *Stream) AppendX(b []byte) {
	s.data = append(s.data, b...)
}

// Paste puts the desired data into a stream, starting at the desired index.
// This function WILL expand the called Stream if the pasted data wouldn't otherwise fit.
func (s *Stream) Paste(startingAt int, pasteMe []byte) {
	dLen := startingAt + len(pasteMe) - len(s.data)
	if dLen > 0 {
		s.data = append(s.data, make([]byte, dLen)...)
	}
	for i, d := range pasteMe {
		s.data[startingAt+i] = d
	}
}

// GetData returns a slice of the called Stream's data
func (s *Stream) GetData() []byte {
	return s.data
}
