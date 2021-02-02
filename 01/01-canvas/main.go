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

	for i := 10; i < 40; i++ {
		for j := 10; j < 40; j++ {
			img.Set(i, j, color.RGBA{R: 255, A: 255})
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
