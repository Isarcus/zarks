package zaudio

import (
	"zarks/output"
)

// FileFormat is a type for file formats
type FileFormat int

// File type
const (
	WAV FileFormat = iota
	MP3
)

// Config contains all of the necessary information to generate a new audio file, other than the actual audio data
type Config struct {
	Format        FileFormat
	Channels      int
	SampleRate    int
	BitsPerSample int
}

// AudioFile is an audio file!
type AudioFile interface {
	AddData(uint64)
	AddDatastream([]uint64)

	ConvertTo(FileFormat) AudioFile

	GetFormat() FileFormat
	GetData() []byte
}

// NewAudioFile makes a new audio file!
func NewAudioFile(cfg Config) AudioFile {
	switch cfg.Format {
	default:
		return newWAV(cfg)
	case WAV:
		return newWAV(cfg)
	}
}

// SaveAudioFile saves an audio file!
func SaveAudioFile(af AudioFile, path string) {
	switch af.GetFormat() {
	case WAV:
		path += ".wav"
	}

	f := output.CreateFile(path)
	defer f.Close()

	output.WriteBytes(f, af.GetData())
}
