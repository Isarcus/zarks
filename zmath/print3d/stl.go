package print3d

import (
	"fmt"
	"math"

	"github.com/Isarcus/zarks/system"
	"github.com/Isarcus/zarks/zmath"
	"github.com/Isarcus/zarks/zmath/zbits"
)

type vec3 struct {
	X, Y, Z float32
}

type triangle struct {
	normal, v1, v2, v3 vec3
	attributes         uint16 // should always be zero
}

// STLData contains all of the information for an STL-encoded 2D shape
type STLData struct {
	Header    [80]byte
	Length    uint32
	Triangles []triangle
}

var normalDown = vec3{0, 0, -1}
var normalZero = vec3{0, 0, 0}

// MapToSTLData takes in a zmath.Map and converts it into bytes
func MapToSTLData(data zmath.Map, title string, makeSolid bool) *STLData {
	// First create the header data
	header := [80]byte{}
	copy(header[:], []byte(title))

	// Get correct bounds
	bounds := data.Bounds()
	bounds.X--
	bounds.Y--

	length := uint32(bounds.X * bounds.Y * 2)
	if makeSolid {
		length += uint32(4 * (bounds.X + bounds.Y)) // for the sides
		length += 2                                 // for the bottom
	}
	triangles := make([]triangle, length)

	// Set the heightmap triangle data
	idx := 0
	for x := 0; x < bounds.X; x++ {
		for y := 0; y < bounds.Y; y++ {
			var v1, v2, v3 vec3
			// first triangle
			v1 = vec3{
				float32(x),
				float32(y),
				float32(data[x][y]),
			}
			v2 = vec3{
				float32(x + 1),
				float32(y),
				float32(data[x+1][y]),
			}
			v3 = vec3{
				float32(x + 1),
				float32(y + 1),
				float32(data[x+1][y+1]),
			}
			triangles[idx] = triangle{
				v1: v1,
				v2: v2,
				v3: v3,
			}
			idx++

			// second triangle
			v2 = vec3{
				float32(x),
				float32(y + 1),
				float32(data[x][y+1]),
			}
			triangles[idx] = triangle{
				v1: v3,
				v2: v2,
				v3: v1,
			}
			idx++
		}
	}

	// If a solid shape, create the sides too
	if makeSolid {
		var v1, v2, v3 vec3
		for y := 0; y <= bounds.Y; y += bounds.Y { // Yes, this line is correct!
			for x := 0; x < bounds.X; x++ {
				v1 = vec3{
					float32(x),
					float32(y),
					0,
				}
				v2 = vec3{
					float32(x),
					float32(y),
					float32(data[x][y]),
				}
				v3 = vec3{
					float32(x + 1),
					float32(y),
					float32(data[x+1][y]),
				}
				triangles[idx] = triangle{
					v1: v1,
					v2: v2,
					v3: v3,
				}
				idx++

				v2 = vec3{
					float32(x + 1),
					float32(y),
					0,
				}
				triangles[idx] = triangle{
					v1: v3,
					v2: v2,
					v3: v1,
				}
				idx++
			}
		}
		for x := 0; x <= bounds.X; x += bounds.X {
			for y := 0; y < bounds.Y; y++ {
				v1 = vec3{
					float32(x),
					float32(y),
					0,
				}
				v2 = vec3{
					float32(x),
					float32(y),
					float32(data[x][y]),
				}
				v3 = vec3{
					float32(x),
					float32(y + 1),
					float32(data[x][y+1]),
				}
				triangles[idx] = triangle{
					v1: v1,
					v2: v2,
					v3: v3,
				}
				idx++

				v2 = vec3{
					float32(x),
					float32(y + 1),
					0,
				}
				triangles[idx] = triangle{
					v1: v3,
					v2: v2,
					v3: v1,
				}
				idx++
			}
		}

		// The bottom 2 triangles
		v1 = vec3{0, 0, 0}
		v2 = vec3{float32(bounds.X), 0, 0}
		v3 = vec3{float32(bounds.X), float32(bounds.Y), 0}
		triangles[idx] = triangle{
			v1: v1,
			v2: v2,
			v3: v3,
		}
		idx++
		v2 = vec3{0, float32(bounds.Y), 0}
		triangles[idx] = triangle{
			v1: v3,
			v2: v2,
			v3: v1,
		}
	}

	return &STLData{
		Header:    header,
		Length:    length,
		Triangles: triangles,
	}
}

// Save creates an STL file from an STLData struct.
func (data *STLData) Save(path string) {
	f := system.CreateFile(path)
	defer f.Close()

	// Write the header (80)
	system.WriteBytes(f, data.Header[:])

	// Write the number of triangles (4)
	system.WriteBytes(f, zbits.Uint32ToBytes(data.Length, zbits.LE))

	// Write the triangle data (50n)
	for _, tri := range data.Triangles {
		system.WriteBytes(f, tri.toBytes())
	}

	fmt.Println("File saved at ", path)
}

func (t triangle) toBytes() []byte {
	bytesNormal := t.normal.toBytes()
	bytesV1 := t.v1.toBytes()
	bytesV2 := t.v2.toBytes()
	bytesV3 := t.v3.toBytes()
	bytesAttributes := zbits.Uint16ToBytes(t.attributes, zbits.LE)

	final := append(bytesNormal, bytesV1...)
	final = append(final, bytesV2...)
	final = append(final, bytesV3...)
	final = append(final, bytesAttributes...)

	return final
}

func (v3 vec3) toBytes() []byte {
	var bits uint32

	bits = math.Float32bits(v3.X)
	bytesX := zbits.Uint32ToBytes(bits, zbits.LE)
	bits = math.Float32bits(v3.Y)
	bytesY := zbits.Uint32ToBytes(bits, zbits.LE)
	bits = math.Float32bits(v3.Z)
	bytesZ := zbits.Uint32ToBytes(bits, zbits.LE)

	return append(append(bytesX, bytesY...), bytesZ...)
}
