package main

import "fmt"

func main() {
	char := 'A'

	// Using fmt.Sprint and escape sequence
	inverted := fmt.Sprintf("\u001b[7m%c\u001b[0m", char)
	fmt.Println(inverted)

	// Using fmt.Sprint and Unicode character
	inverted = fmt.Sprintf("%c\u0336", char)
	fmt.Println(inverted)
}
