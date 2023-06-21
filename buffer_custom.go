package main

import (
	termbox "github.com/nsf/termbox-go"
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	xMax, yMax := termbox.Size()
	// Create your custom buffer
	customBuffer := make([][]termbox.Cell, yMax)
	for i := range customBuffer {
		customBuffer[i] = make([]termbox.Cell, xMax)
	}

	// Modify your custom buffer
	for y := range customBuffer {
		for x := range customBuffer[y] {
			customBuffer[y][x] = termbox.Cell{
				Ch: 'A',
				Fg: termbox.ColorRed,
				Bg: termbox.ColorBlue,
			}
		}
	}

	// Copy the custom buffer to the Termbox screen buffer
	for y := range customBuffer {
		for x, cell := range customBuffer[y] {
			termbox.SetCell(x, y, cell.Ch, cell.Fg, cell.Bg)
		}
	}

	termbox.Flush()

	// Wait for a key press to exit
	for {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey {
			break
		}
	}
}
