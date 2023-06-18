package main

import (
	"fmt"
)

func createAlternativeString(message string) string {
	length := len(message)
	alternative := make([]byte, length)

	for i := 0; i < length-1; i++ {
		switch {
		case message[i] == '0' && message[i+1] == '0':
			alternative[i] = '0'
		case message[i] == '0' && message[i+1] == '1':
			alternative[i] = 'R'
		case message[i] == '1' && message[i+1] == '0':
			alternative[i] = 'L'
		case message[i] == '1' && message[i+1] == '1':
			alternative[i] = '1'
		}
	}

	return string(alternative)
}

func main() {
	message := "011010011"
	alternative := createAlternativeString(message)
	fmt.Println("Message:", message)
	fmt.Println("Alternative:", alternative)
}
