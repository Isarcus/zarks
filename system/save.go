package system

import (
	"bufio"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
)

// CreateFile creates a file at the target location
func CreateFile(path string) *os.File {
	if FileExists(path) {
		reader := bufio.NewReader(os.Stdin)
		if !QueryYN(reader, "File at "+path+" already exists. Overwrite?") {
			fmt.Println("Save aborted.")
			return nil
		}
	}
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	return file
}

// WriteBytes writes the passed bytes to the file, trying up to 2 more times to write all the bytes
func WriteBytes(f *os.File, bytes []byte) {
	written := 0
	tries := 0
	for written < len(bytes) && tries < 2 {
		l, err := f.Write(bytes[written:])
		written += l
		if err != nil {
			fmt.Println("Error writing in try #", tries, ": ", err)
		}
		tries++
	}

}

// SaveImage saves the image as a PNG at the provided path
func SaveImage(path string, img *image.RGBA) {
	if FileExists(path) {
		reader := bufio.NewReader(os.Stdin)
		if !QueryYN(reader, "File at "+path+" already exists. Overwrite?") {
			fmt.Println("Save aborted.")
			return
		}
	}
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		fmt.Println("[ERROR] Save aborted.")
		fmt.Println(err)
		return
	}

	// III. Save image
	err = png.Encode(file, img)
	if err != nil {
		fmt.Println("[ERROR] Save aborted.")
		fmt.Println(err)
		return
	}

	fmt.Println("Image saved at ", path)
}

// SaveImageJPEG saves an image as a JPEG at the target location. Quality should range from 1 to 100 inclusive,
// but note that saving as a jpeg is never lossless.
func SaveImageJPEG(path string, img *image.RGBA, quality int, checkIfExists bool) {
	if checkIfExists && FileExists(path) {
		reader := bufio.NewReader(os.Stdin)
		if !QueryYN(reader, "File at "+path+" already exists. Overwrite?") {
			fmt.Println("Save aborted.")
			return
		}
	}
	file, err := os.Create(path)
	if err != nil {
		fmt.Println("[ERROR] Save aborted.")
		fmt.Println(err)
		return
	}

	// III. Save image
	options := jpeg.Options{
		Quality: quality,
	}
	err = jpeg.Encode(file, img, &options)
	if err != nil {
		fmt.Println("[ERROR] Save aborted.")
		fmt.Println(err)
		return
	}

	fmt.Println("Image saved at ", path)
}
