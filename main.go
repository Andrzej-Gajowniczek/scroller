package main

import (
	_ "embed"
	"fmt"
	"strings"
	"time"
)

//go:embed "data/small8.64c"
var data []byte

func renderChar(b byte) *[]string {

	var items = []rune{
		'@', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S',
		'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '[', '~', ']', '|', '\\', ' ', '!', '"', '#', '$', '%', '&',
		'\'', '(', ')', '*', '+', ',', '-', '.', '/', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		':', ';', '<', '=', '>', '?',
	}

	translator := make(map[byte]int)
	for i, x := range items {
		translator[byte(x)] = i * 8

	}

	/*
		for i := 64; i < 128; i++ {
			translator[byte(i)] = (i - 64) * 8
		}
	*/

	charset := data[2:]
	var rendered = make([]string, 0, 8)

	for y := 0; y < 8; y++ {

		t := translator[b]
		z := t + y
		x := charset[z]
		struna := fmt.Sprintf("%08b", x)
		//struna = strings.ReplaceAll(struna, "0", ` `)
		//struna = strings.ReplaceAll(struna, "1", `█`)
		rendered = append(rendered, struna)
	}
	return &rendered
}

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

func makeSemigraphic(subString string) string {
	subString = strings.ReplaceAll(subString, "0", ` `)
	subString = strings.ReplaceAll(subString, "1", `█`)
	subString = strings.ReplaceAll(subString, "R", `▐`)
	subString = strings.ReplaceAll(subString, "L", `▌`)
	return subString
}
func main() {

	message := "KOCHAM PAULINKE MOJA ZONE"

	// this buffer is consist forom Zeros and ones {0,1} representing bits in: var data []byte - which is charset loaded from file before compilation to be part of binary
	var onScreenBuffer = []string{"", "", "", "", "", "", "", ""}
	// this would be buffer consist of 1,0,R,L depends on schematic:
	// i:0,i+:0 -> 0 which represents " "
	// i:0,i+:1 -> R which represents "▐"
	// i:1,i+:0 -> L which represents "▌"
	// i:1,i+:1 -> 1 which represents "█"
	// It is created to shift left full block █ a half field for left ▐▌ - so as you can see the second cursor consists of two fields and it's location is half field shifted to left

	for _, u := range message {
		tablica := renderChar(byte(u))
		for i, c := range *tablica {
			onScreenBuffer[i] += c
		}
	}
	fmt.Println(onScreenBuffer)
	var duplicatShifter = []string{"", "", "", "", "", "", "", ""}
	for i, msg := range onScreenBuffer {

		alternative := createAlternativeString(msg)
		duplicatShifter[i] += alternative

	}
	fmt.Println(duplicatShifter)

	//var camoBuffer [8][]rune
	//fmt.Println(camoBuffer, reflect.TypeOf(camoBuffer))
	for p := 0; p < (8 * len(message)); p++ {

		for _, info := range onScreenBuffer {
			//runeStartIndex := utf8.RuneCountInString(info[:p])
			//subString := []rune(info)[runeStartIndex:]
			subString := info[0:p]
			subString = strings.ReplaceAll(subString, "0", ` `)
			subString = strings.ReplaceAll(subString, "1", `█`)

			fmt.Println(string(subString))
		}
		time.Sleep(1000 * time.Millisecond / 30)

		for _, info := range duplicatShifter {
			//runeStartIndex := utf8.RuneCountInString(info[:p])
			//subString := []rune(info)[runeStartIndex:]
			subString := info[0:p]
			subString = strings.ReplaceAll(subString, "0", ` `)
			subString = strings.ReplaceAll(subString, "1", `█`)
			subString = strings.ReplaceAll(subString, "R", `▐`)

			subString = strings.ReplaceAll(subString, "L", `▌`)
			fmt.Println(string(subString))
		}

		time.Sleep(1000 * time.Millisecond / 30)
	}
	testówka := "00011000010011011010101111"
	fmt.Println(makeSemigraphic(testówka))
	fmt.Println(makeSemigraphic(createAlternativeString(testówka)))
	//	}

}
