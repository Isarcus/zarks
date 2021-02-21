package gif

import (
	"fmt"
	"image"
	"image/color"

	"github.com/Isarcus/zarks/system"
	"github.com/Isarcus/zarks/zimg"
	"github.com/Isarcus/zarks/zmath"
	"github.com/Isarcus/zarks/zmath/zbits"
)

type version string

// Versions
const (
	GIF89a version = "GIF89a"
	GIF87a version = "GIF87a"
)

// Byte codes & presets
const (
	BlockExtension byte = 0x21
	BlockImage     byte = 0x2C
	BlockGraphics  byte = 0xF9
	Trailer        byte = 0x3B

	//                         ABBBCDDD
	DescriptorDefault byte = 0b11110111 // = 247
	// A: existence of global color table
	// B: N, where 2^(N+1) colors are in the color table
	// C: whether color table is sorted by frequency (outdated)
	// D: also size of color table??????? just roll with it I guess
)

// GIF contains all of the data for a GIF in bytes
type GIF struct {
	header      [6]byte
	descriptor  [7]byte
	colorTable  [768]byte
	loopControl [19]byte
	frames      []Frame

	Bounds  zmath.VecInt
	Pallete Pallete
}

// Frame contains color indices
type Frame struct {
	descriptor [10]byte
	data       []byte
}

// NewGIF returns a new empty GIF to which Frames (of the correct dimensions) can be added.
func NewGIF(bounds zmath.VecInt, v version, pal Pallete) *GIF {
	//-----------//
	// I. HEADER //
	//-----------//
	header := [6]byte{}
	for i := range v {
		header[i] = v[i]
	}

	//-------------------------------//
	// II. LOGICAL SCREEN DESCRIPTOR //
	//-------------------------------//
	var (
		width  = zbits.Uint16ToBytes(uint16(bounds.X), zbits.LE)
		height = zbits.Uint16ToBytes(uint16(bounds.Y), zbits.LE)
		desc   = [7]byte{}
	)
	desc[0] = width[0]
	desc[1] = width[1]
	desc[2] = height[0]
	desc[3] = height[1]

	fmt.Println("Descriptor Default: ", DescriptorDefault)

	desc[4] = DescriptorDefault // = 0b11110111
	desc[5] = 0                 // index of background color in color table (outdated)
	desc[6] = 0                 // pixel aspect ratio (outdated)

	//-------------------------//
	// III. GLOBAL COLOR TABLE //
	//-------------------------//
	table := [768]byte{}
	for i, col := range pal {
		idx := i * 3
		table[idx] = col.R
		table[idx+1] = col.G
		table[idx+2] = col.B
	}

	//-----------------------------//
	// IV. ANIMATION (GIF89a ONLY) //
	//-----------------------------//
	lpCtrl := [19]byte{}
	lpCtrl[0] = BlockExtension             // Must be 0x21 to introduce any extension
	lpCtrl[1] = 0xFF                       // "application" code
	lpCtrl[2] = 0x0B                       // length in bytes of information to follow (11 bytes)
	for i, letter := range "NETSCAPE2.0" { // indices 3 to 13, inclusive
		lpCtrl[i+3] = byte(letter)
	}
	lpCtrl[14] = 0x03 // length of informational bytes to follow
	lpCtrl[15] = 0x01
	// bytes 16, 17 are how many times gif should loop (leave as 0 for infinite)
	lpCtrl[16] = 0
	lpCtrl[17] = 0
	// byte 18 is always empty
	lpCtrl[18] = 0

	return &GIF{
		header:      header,
		descriptor:  desc,
		colorTable:  table,
		loopControl: lpCtrl,
		frames:      make([]Frame, 0),

		Pallete: pal,
		Bounds:  bounds,
	}
}

// AddFrame adds an image to an existing GIF. The returned bool is whether it was able to be successfully added.
func (gif *GIF) AddFrame(img *image.RGBA) bool {
	bounds := zmath.VecInt{
		X: img.Bounds().Max.X,
		Y: img.Bounds().Max.Y,
	}
	if bounds != gif.Bounds {
		return false
	}

	//---------------------//
	// I. IMAGE DESCRIPTOR //
	//---------------------//
	desc := [10]byte{}
	desc[0] = BlockImage
	// bytes 1, 2, 3, 4 are for image start position (ignore for simplicity)
	// bytes 5, 6, 7, 8 are for image dimensions
	desc[5] = gif.descriptor[0]
	desc[6] = gif.descriptor[1]
	desc[7] = gif.descriptor[2]
	desc[8] = gif.descriptor[3]
	// byte 9 is a packed field for local color table info (ignore for simplicity)

	//------------------//
	// II. INDEX STREAM //
	//------------------//
	colorIndices := make([]byte, bounds.X*bounds.Y)
	colorIndex := 0

	for y := 0; y < bounds.Y; y++ {
		for x := 0; x < bounds.X; x++ {
			// Find which pallete color is closest to the image color
			var (
				minIdx  int        = -1
				minDist float64    = 500 // above the max possible distance, ~441.7
				imgCol  color.RGBA = img.RGBAAt(x, y)
			)
			for idx, palCol := range gif.Pallete {
				dist := zimg.DistanceToColor(imgCol, palCol)
				if dist < minDist {
					minDist = dist
					minIdx = idx
				}
			}

			if minIdx == -1 {
				fmt.Println("Minimum color distance not found!")
			}

			colorIndices[colorIndex] = uint8(minIdx)
			colorIndex++
		}
	}

	//------------------//
	// III. CODE STREAM //
	//------------------//

	// Make a new codestream & feed it the color indices
	codeStream := newStream(256)
	for _, data := range colorIndices {
		codeStream.addByte(data)
	}

	// compress the data
	compressed := codeStream.finalize()

	// group the data into 255-byte chunks
	var (
		count   = 0
		chunked = make([]byte, 2)
		nextLen = zmath.MinInt(len(compressed), 255)
	)
	chunked[0] = 0x08          // [OK] 8 is the LZW minimum code size for 256 colors
	chunked[1] = byte(nextLen) // length of the first data chunk
	for i, data := range compressed {
		if count == nextLen {
			nextLen = zmath.MinInt(len(compressed)-i-1, 255) // 255
			chunked = append(chunked, byte(nextLen))
			count = 0
		}
		chunked = append(chunked, data)
		count++
	}

	fmt.Println("Length of compressed data: ", len(chunked))

	newFrame := Frame{
		descriptor: desc,
		data:       chunked,
	}

	gif.frames = append(gif.frames, newFrame)

	return true
}

// Save saves a GIF!
func (gif *GIF) Save(path string) {
	f := system.CreateFile(path)
	defer f.Close()

	system.WriteBytes(f, gif.header[:])
	system.WriteBytes(f, gif.descriptor[:])
	system.WriteBytes(f, gif.colorTable[:])
	if gif.header[4] == "9"[0] { // if version 89a, include loop extension
		fmt.Println("[GIF89a] Loop extension added.")
		system.WriteBytes(f, gif.loopControl[:])
	}

	for _, frame := range gif.frames {
		system.WriteBytes(f, frame.descriptor[:])
		system.WriteBytes(f, frame.data[:])
	}

	system.WriteBytes(f, []byte{Trailer})
}
