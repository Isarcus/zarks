package zpixel

import (
	"fmt"
	"image"
	_ "image/jpeg" // necessary for image decoding
	_ "image/png"  // necessary for image decoding
	"io/ioutil"
	"os"
	"strings"

	"github.com/faiface/pixel/pixelgl"

	"github.com/faiface/pixel"
)

// CreateWindow gives you a *pixelgl.Window, and handles errors by itself
func CreateWindow(cfg pixelgl.WindowConfig) *pixelgl.Window {
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	return win
}

// LoadImage returns an image.Image, NOT a pixel.Picture!!!
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

// LoadPicture returns a pixel.Picture at the target path
func LoadPicture(path string) pixel.Picture {
	img := LoadImage(path)
	return pixel.PictureDataFromImage(img)
}

// LoadBaseSprite loads the image at the target path and makes a sprite out of it
func LoadBaseSprite(path string) *pixel.Sprite {
	pic := LoadPicture(path)
	return pixel.NewSprite(pic, pic.Bounds())
}

// CreateSpriteFromPicture takes an existing pixel.Picture and gives you a *pixel.Sprite
func CreateSpriteFromPicture(pic pixel.Picture) *pixel.Sprite {
	return pixel.NewSprite(pic, pic.Bounds())
}

// LoadTextLines loads a file and returns all of its lines as members of a string array.
func LoadTextLines(path string) []string {
	byteContent, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("[ERROR] could not open file at " + path)
		panic(err)
	}
	textContent := strings.Fields(string(byteContent)) // cuts up into individual lines
	cleanText := make([]string, len(textContent))
	for i, s := range textContent {
		cleanText[i] = strings.TrimSuffix(s, "\n")
	}

	return cleanText
}
