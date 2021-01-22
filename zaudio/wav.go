package zaudio

import "zarks/zmath/zbits"

type wav struct {
	cfg    Config
	header *zbits.Stream
	bytes  *zbits.Stream
}

func newWAV(cfg Config) *wav {
	var (
		byteRate   = cfg.Channels * cfg.SampleRate * cfg.BitsPerSample / 8
		blockAlign = cfg.BitsPerSample * cfg.Channels
	)

	h := zbits.NewStream(44)
	h.Paste(0, []byte("RIFF")) //                                         // mark file as a RIFF
	/**/                                                                  // bytes 4-7 are for the total filesize
	h.Paste(8, []byte("WAVE"))                                            // file format (WAV)
	h.Paste(12, []byte("fmt "))                                           // format chunk marker
	h.Paste(16, zbits.Uint32ToBytes(16, zbits.LE))                        // size of format chunk
	h.Paste(20, zbits.Uint16ToBytes(1, zbits.LE))                         // encoding format (1 = PCM)
	h.Paste(22, zbits.Uint16ToBytes(uint16(cfg.Channels), zbits.LE))      // channels
	h.Paste(24, zbits.Uint32ToBytes(uint32(cfg.SampleRate), zbits.LE))    // sample rate
	h.Paste(28, zbits.Uint32ToBytes(uint32(byteRate), zbits.LE))          // byte rate
	h.Paste(32, zbits.Uint16ToBytes(uint16(blockAlign), zbits.LE))        //
	h.Paste(34, zbits.Uint16ToBytes(uint16(cfg.BitsPerSample), zbits.LE)) //
	h.Paste(36, []byte("data"))

	return &wav{
		cfg:    cfg,
		header: h,
		bytes:  zbits.NewStream(0),
	}
}

func (w *wav) AddData(val uint64) {
	switch w.cfg.BitsPerSample {
	case 8:
		w.bytes.Append(byte(val))
	case 16:
		w.bytes.AppendX(zbits.Uint16ToBytes(uint16(val), zbits.LE))
	case 32:
		w.bytes.AppendX(zbits.Uint32ToBytes(uint32(val), zbits.LE))
	case 64:
		w.bytes.AppendX(zbits.Uint64ToBytes(val, zbits.LE))
	}
}

func (w *wav) AddDatastream(data []uint64) {
	for _, d := range data {
		w.AddData(d)
	}
}

func (w *wav) ConvertTo(f FileFormat) AudioFile { // TODO
	return w
}

func (w *wav) GetData() []byte {
	var (
		dataSize = uint32(w.bytes.Size())
		fileSize = dataSize + 44
	)

	w.header.Paste(4, zbits.Uint32ToBytes(fileSize, zbits.LE))
	w.header.Paste(40, zbits.Uint32ToBytes(dataSize, zbits.LE))

	stream := zbits.NewStream(0)
	stream.AppendX(w.header.GetData())
	stream.AppendX(w.bytes.GetData())

	return stream.GetData()
}

func (w *wav) GetFormat() FileFormat {
	return w.cfg.Format
}
