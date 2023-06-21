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

	message := "Hello, Termbox!"
	x, y := 2, 2
	fgColor := termbox.ColorWhite
	bgColor := termbox.ColorBlack

	drawText(x, y, message, fgColor, bgColor)
	termbox.Flush()

	// Change the color of the rune after a delay
	// Here, we change the foreground color to yellow
	// and the background color to blue
	termbox.Sync()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	fgColor = termbox.ColorYellow
	bgColor = termbox.ColorBlue
	drawText(x, y, message, fgColor, bgColor)
	termbox.Flush()

	// Wait for a key press to exit
	for {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey {
			break
		}
	}
}

func drawText(x, y int, text string, fg, bg termbox.Attribute) {
	for i, ch := range text {
		termbox.SetCell(x+i, y, ch, fg, bg)
	}
}
