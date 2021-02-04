package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

type point struct {
	x, y, z float64
}

func (p point) Sub(s point) point {
	return point{
		x: p.x - s.x,
		y: p.y - s.y,
		z: p.z - s.z,
	}
}

func (p point) Dot(s point) float64 { return p.x*s.x + p.y*s.y + p.z*s.z }

type sphere struct {
	color.Color
	point
	radius float64
}

var (
	Cw float64 = 500
	Ch float64 = 500

	Vw float64 = 500
	Vh float64 = 500

	d     float64 = 1
	scene         = struct {
		spheres []*sphere
	}{
		spheres: []*sphere{
			// {Color: color.RGBA{R: 255, A: 255}, point: point{20, 20, 20}, radius: 100},
			{Color: color.RGBA{R: 255, A: 255}, point: point{0, -10, 30}, radius: 100},
			// {Color: color.RGBA{B: 255, A: 255}, point: point{20, 0, 40}, radius: 100},
			{Color: color.RGBA{G: 255, A: 255}, point: point{-20, 0, 40}, radius: 100},
		},
	}

	BACKGROUND_COLOR = color.White
)

func main() {
	img := image.NewRGBA(image.Rect(int(-Cw/2), int(-Ch/2), int(Cw/2), int(Ch/2)))

	O := point{}
	for x := -Cw / 2; x < Cw/2; x++ {
		for y := -Ch / 2; y < Ch/2; y++ {
			D := CanvasToViewport(x, y)
			color := TraceRay(O, D, 1, math.Inf(1))
			img.Set(int(x), int(y), color)
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

func width(img *image.RGBA) int {
	return img.Bounds().Max.X - img.Bounds().Min.X
}

func height(img *image.RGBA) int {
	return img.Bounds().Max.Y - img.Bounds().Min.Y
}

func CanvasToViewport(x, y float64) point {
	return point{x * Vw / Cw, y * Vh / Ch, d}
}

func TraceRay(O, D point, tMin, tMax float64) color.Color {
	closestT := math.Inf(1)
	var closestSphere *sphere

	for _, s := range scene.spheres {
		t1, t2 := IntersectRaySphere(O, D, s)
		if tMin < t1 && t1 < tMax && t1 < closestT {
			closestT = t1
			closestSphere = s
		}
		if tMin < t2 && t2 < tMax && t2 < closestT {
			closestT = t2
			closestSphere = s
		}
	}
	if closestSphere == nil {
		return BACKGROUND_COLOR
	}
	return closestSphere.Color
}

func IntersectRaySphere(O, D point, s *sphere) (float64, float64) {
	CO := O.Sub(s.point)

	a := D.Dot(D)
	b := 2 * CO.Dot(D)
	c := CO.Dot(CO) - s.radius*s.radius

	discriminant := b*b - 4*a*c
	if discriminant < 0 {
		return math.Inf(1), math.Inf(1)
	}

	t1 := (-b + math.Sqrt(discriminant)) / (2 * a)
	t2 := (-b - math.Sqrt(discriminant)) / (2 * a)
	return t1, t2
}
