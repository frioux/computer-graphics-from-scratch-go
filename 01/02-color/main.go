package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

func main() {
	img := image.NewRGBA(image.Rect(0, 0, 1920, 1080))

	for i := 0; i < 1920; i++ {
		for j := 0; j < 1080; j++ {
			img.Set(i, j, color.Black)
		}
	}

	// red square
	for i := 10; i < 40; i++ {
		for j := 10; j < 40; j++ {
			c := img.At(i, j).(color.RGBA)
			c.R = 255
			img.Set(i, j, c)
		}
	}

	// blue square
	for i := 30; i < 70; i++ {
		for j := 30; j < 70; j++ {
			c := img.At(i, j).(color.RGBA)
			c.G = 255
			img.Set(i, j, c)
		}
	}

	f, err := os.Create("image.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}
