package system

import (
	"fmt"
	"image"
	_ "image/jpeg" // necessary for image decoding
	_ "image/png"  // necessary for image decoding
	"os"
)

// LoadImage returns an *image.Image
func LoadImage(path string) image.Image {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("[ERROR] Could not open file at " + path)
		panic(err)
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("[ERROR] Opened, but could not decode file at " + path)
		panic(err)
	}
	return img
}
