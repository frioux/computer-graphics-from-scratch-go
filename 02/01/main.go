package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

type point [3]float64

// Dot product of two 3D vectors.
func DotProduct(v1, v2 point) float64 {
	return v1[0]*v2[0] + v1[1]*v2[1] + v1[2]*v2[2]
}

// Computes v1 - v2.
func Subtract(v1, v2 point) point {
	return point{v1[0] - v2[0], v1[1] - v2[1], v1[2] - v2[2]}
}

// A Sphere.
type Sphere struct {
	center point
	radius float64
	color.Color
}

// Scene setup.
var viewport_size float64 = 1
var projection_plane_z float64 = 1
var camera_position point
var background_color = color.White
var spheres = []Sphere{
	{center: point{0, -1, 3}, radius: 1, Color: color.RGBA{R: 255, A: 255}},
	{center: point{2, 0, 4}, radius: 1, Color: color.RGBA{B: 255, A: 255}},
	{center: point{-2, 0, 4}, radius: 1, Color: color.RGBA{G: 255, A: 255}},
}

var canvas *image.RGBA

func width(img *image.RGBA) float64 {
	return float64(img.Bounds().Max.X - img.Bounds().Min.X)
}

func height(img *image.RGBA) float64 {
	return float64(img.Bounds().Max.Y - img.Bounds().Min.Y)
}

// Converts 2D canvas coordinates to 3D viewport coordinates.
func CanvasToViewport(p2d point) point {
	return point{p2d[0] * viewport_size / width(canvas),
		p2d[1] * viewport_size / height(canvas),
		projection_plane_z}
}

// Computes the intersection of a ray and a sphere. Returns the values
// of t for the intersections.
func IntersectRaySphere(origin, direction point, sphere Sphere) (float64, float64) {
	oc := Subtract(origin, sphere.center)

	k1 := DotProduct(direction, direction)
	k2 := 2 * DotProduct(oc, direction)
	k3 := DotProduct(oc, oc) - sphere.radius*sphere.radius

	discriminant := k2*k2 - 4*k1*k3
	if discriminant < 0 {
		return math.Inf(1), math.Inf(1)
	}

	t1 := (-k2 + math.Sqrt(discriminant)) / (2 * k1)
	t2 := (-k2 - math.Sqrt(discriminant)) / (2 * k1)
	return t1, t2
}

// Traces a ray against the set of spheres in the scene.
func TraceRay(origin, direction point, min_t, max_t float64) color.Color {
	closest_t := math.Inf(1)
	var closest_sphere Sphere

	for i := range spheres {
		t0, t1 := IntersectRaySphere(origin, direction, spheres[i])
		if t0 < closest_t && min_t < t0 && t0 < max_t {
			closest_t = t0
			closest_sphere = spheres[i]
		}
		if t1 < closest_t && min_t < t1 && t1 < max_t {
			closest_t = t1
			closest_sphere = spheres[i]
		}
	}

	if closest_sphere.radius == 0 {
		return background_color
	}

	return closest_sphere.Color
}

//
// Main loop.
//
func main() {
	canvas = image.NewRGBA(image.Rect(-300, -300, 300, 300))

	for x := -width(canvas) / 2; x < width(canvas)/2; x++ {
		for y := -height(canvas) / 2; y < height(canvas)/2; y++ {
			direction := CanvasToViewport(point{x, y, 0})
			color := TraceRay(camera_position, direction, 1, math.Inf(1))
			canvas.Set(int(math.Round(x)), int(-math.Round(y))-1, color)
		}
	}

	f, err := os.Create("chapter-02.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := png.Encode(f, canvas); err != nil {
		panic(err)
	}
}
