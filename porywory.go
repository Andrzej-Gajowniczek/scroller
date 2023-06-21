package main

import (
	"fmt"
	"log"

	"github.com/nsf/termbox-go"
)

func printAt(x, y int, s string, par ...termbox.Attribute) {

	ink := termbox.ColorDefault
	bg := termbox.ColorDefault
	switch len(par) {
	case 2:
		ink = par[0]
		bg = par[1]
	case 1:
		ink = par[0]
	default:

	}
	r := []rune(s)
	for i, v := range r {

		termbox.SetCell(x+i, y, v, ink, bg)
	}

}

func main() {
	err := termbox.Init()
	if err != nil {
		log.Fatal("termbox error", err)
	}
	defer termbox.Close()
	defer termbox.Flush()

	xMax, yMax := termbox.Size()
	napis := "Lubię Go za jego prostotę!"
	lenno := len(napis)
	poz := (xMax - lenno) / 2
	printAt(poz, yMax/2, napis, 2, 0)
	info := fmt.Sprintf("xMax:%d yMax:%d len:%d napie: %s ", xMax, yMax, lenno, napis)
	printAt(0, 0, info, 0, 2)

	termbox.Flush()
	for {
		ev := termbox.PollEvent()
		if ev.Key == termbox.KeyEsc {
			return
		}
	}
}
