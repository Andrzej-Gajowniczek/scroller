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
	rotationSpeed := 0.05

	// Calculate the half dimensions of the cube
	halfCubeSize := cubeSize / 2.0

	// Initialize the animation loop
	for {
		// Clear the console screen
		fmt.Print("\033[2J")
		fmt.Print("\033[H")

		// Perform cube rotation
		angle := time.Now().UnixNano() / int64(time.Millisecond) * int64(rotationSpeed)
		sin := math.Sin(float64(angle) * math.Pi / 180)
		cos := math.Cos(float64(angle) * math.Pi / 180)

		// Calculate the rotated coordinates of the cube
		xRotated := halfCubeSize * cos
		yRotated := halfCubeSize * sin

		// Render the cube
		for y := -halfCubeSize; y <= halfCubeSize; y++ {
			for x := -halfCubeSize; x <= halfCubeSize; x++ {
				if (math.Abs(x-xRotated) <= 0.5) && (math.Abs(y-yRotated) <= 0.5) {
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
