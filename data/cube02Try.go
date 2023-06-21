package main

import (
	"fmt"
	"math"
	"time"
)

func main() {
	// Define the dimensions of the console window
	const (
		width  = 80
		height = 24
	)

	// Set up the initial cube parameters
	cubeSize := 10.0
	rotationSpeed := 0.1

	for {
		// Clear the console screen
		fmt.Print("\033[2J")
		fmt.Print("\033[H")

		// Perform cube rotation
		angle := time.Now().UnixNano() / int64(time.Millisecond) * int64(rotationSpeed)
		sin := math.Sin(float64(angle) * math.Pi / 180)
		cos := math.Cos(float64(angle) * math.Pi / 180)

		// Render the cube
		for y := -cubeSize; y <= cubeSize; y++ {
			for x := -cubeSize; x <= cubeSize; x++ {
				point := rotatePoint(x, y, cos, sin)
				if point.z > -cubeSize {
					fmt.Print("#")
				} else {
					fmt.Print(" ")
				}
			}
			fmt.Println()
		}

		// Pause for a short duration before rendering the next frame
		time.Sleep(50 * time.Millisecond)
	}
}

// rotatePoint rotates a 2D point (x, y) by given cos and sin values
// and returns the resulting 3D point (x, y, z)
func rotatePoint(x, y float64, cos, sin float64) Point {
	const (
		zDistance = 10.0 // Distance of the cube from the viewer
	)

	// Rotate the point
	newX := x*cos - y*sin
	newY := x*sin + y*cos

	// Calculate the z coordinate based on the distance
	z := zDistance

	return Point{newX, newY, z}
}

// Point represents a 3D point in space
type Point struct {
	x, y, z float64
}
