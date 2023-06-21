package main

import (
	_ "embed"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/nsf/termbox-go"
)

var at = termbox.SetCell

func printAt(x, y int, s string, colors ...termbox.Attribute) {
	ink, bg := termbox.ColorDefault, termbox.ColorDefault
	switch len(colors) {
	case 1:
		ink = colors[0]
		bg = termbox.ColorDefault
	case 2:
		ink = colors[0]
		bg = colors[1]
	default:

	}

	for i, c := range s {
		at(x+i, y, c, ink, bg)
	}
}

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

	var duplicatShifter = []string{"", "", "", "", "", "", "", ""}
	for i, msg := range onScreenBuffer {

		alternative := createAlternativeString(msg)
		duplicatShifter[i] += alternative
	}

	err := termbox.Init()
	if err != nil {
		log.Fatal("kicha z termboksem", err)
	}
	defer termbox.Close()
	//

	xMax, yMax := termbox.Size()
	infoXY := fmt.Sprintf("x:%3d; y:%3d", xMax, yMax)
	printAt(0, 0, infoXY, termbox.ColorDefault, termbox.ColorDefault)
	/*
		for c := 0; c < 32; c++ {
			printAt(0, 1+c, "Eliza", termbox.Attribute(c), 0)
		}
		printAt(10, 10, "Migotka", termbox.AttrBlink|termbox.ColorLightCyan)
		at(12, 10, '*', 0, termbox.ColorDefault)
		termbox.Flush()
		termbox.SetFg(12, 10, termbox.ColorLightMagenta)
	*/
	termbox.Sync()
	customBuffer := make([][]termbox.Cell, len(onScreenBuffer))
	for i := range customBuffer {
		customBuffer[i] = make([]termbox.Cell, len(onScreenBuffer[0]))
	}
	customBuffer2 := customBuffer
	copy(customBuffer, customBuffer2)

	rozmiarBufora := fmt.Sprintf("len(customBuffer2):%v", len(customBuffer2)*len(customBuffer2[0]))
	printAt(20, 20, rozmiarBufora)
	termbox.Sync()

	//copy strings from onScreenBuffer -> customeBuffer of type termbox.Cell
	for yy, s := range onScreenBuffer {
		r := []rune(makeSemigraphic(s))
		for xx, rr := range r {
			customBuffer[yy][xx].Ch = rr
		}
	}
	//copy customBuffer to internalBuffer of termbox

	for {
		for yy := 0; yy < 8; yy++ {
			for xx := 0; xx < xMax; xx++ {
				termbox.SetCell(xx, yy, customBuffer[yy][xx].Ch, termbox.ColorDefault, termbox.ColorDefault)
			}
		}
		termbox.Sync()
		time.Sleep(time.Millisecond * (1000 / 60))

		eV := termbox.PollRawEvent()

		if eV.Key == termbox.KeyEsc {
			return

		}

	}
}
