package main

import (
	_ "embed"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/nsf/termbox-go"
)

//go:embed "data/small8.64c"
var data []byte

// this func input Ascii capital letter byte code and returns 8x8 font consist of 0 and 1 - 8 strings by 8x Zeros or Ones
func renderChar(b byte) *[]string {

	var items = []rune{
		'@', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S',
		'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '[', '~', ']', '|', '\\', ' ', '!', '"', '#', '$', '%', '&',
		'\'', '(', ')', '*', '+', ',', '-', '.', '/', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		':', ';', '<', '=', '>', '?',
	}

	//this func maps code of letters with indices of pixels data regarding shape of certain semigraphics image of the letter
	translator := make(map[byte]int)
	for i, x := range items {
		translator[byte(x)] = i * 8

	}
	//charset data starts from the 3rd byte
	charset := data[2:]
	var rendered = make([]string, 0, 8) //create space for semigraphics image consist of Zeros and Ones!

	for y := 0; y < 8; y++ {

		t := translator[b]
		z := t + y
		x := charset[z]
		struna := fmt.Sprintf("%08b", x)
		rendered = append(rendered, struna)
	}
	return &rendered //return address of semigraphics "Big" image 8x8 cursors size.
}

// this func shifts left 0,1 make them 0,1,R or L font metacode for further exchanging by semigraphics
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

// this func exchange 0,1,L,R by " "█"▐"▌" - what delivers posibility to scroll by half a cursor
func makeSemigraphic(subString string) string {
	subString = strings.ReplaceAll(subString, "0", ` `)
	subString = strings.ReplaceAll(subString, "1", `█`)
	subString = strings.ReplaceAll(subString, "R", `▐`)
	subString = strings.ReplaceAll(subString, "L", `▌`)
	return subString
}

type scroller struct {
	messageString       string              //text to scroll
	colorMessage        []termbox.Attribute //text colorization but of termbox.Attribute type
	messageStringMatrix [][]string          //two buffers for cursor fonts consist of 0,1 or in second buf 0,1,L,R
	messageBuferCells   [][][]termbox.Cell  //two buffers consist of termbox cells
	xMax                int
	yMax                int
	progress            int
	index               int
	interval            int
	lcx                 int
	lcy                 int
	rcx                 int
	rcy                 int
	frame               int
	speed               int
}

// This func is for loading text to be scrolled and in the meantime changes all letters to be uppercase - embeded data requires it
func (s *scroller) loadMessage(ss string) {

	s.messageString = strings.ToUpper(ss)

}
func (s *scroller) scrollerInit() { //Initialize scroll tables;colors; ascii table; shift half cursor semigraphics; create buffers4termbox
	//colorize ascii text

	var randomNumber int
	s.colorMessage = make([]termbox.Attribute, len(s.messageString))

	//allocate [1][7] matrix dimension (2x8 strings)
	s.messageStringMatrix = append(s.messageStringMatrix, []string{"", "", "", "", "", "", "", ""})
	s.messageStringMatrix = append(s.messageStringMatrix, []string{"", "", "", "", "", "", "", ""})
	oldRandom := 0
	for i, v := range s.messageString {
		if v == ' ' {
			// Generate a random number between 8 and 16 (termbox light colors numbers)
		onceAgain:
			randomNumber = rand.Intn(9) + 8
			if (randomNumber == oldRandom) || (randomNumber == 9) {
				goto onceAgain
			}
			oldRandom = randomNumber

		}
		s.colorMessage[i] = termbox.Attribute(randomNumber)
		rendered8x8 := renderChar(byte(v))
		for y, str := range *rendered8x8 {
			s.messageStringMatrix[0][y] = s.messageStringMatrix[0][y] + str
		}

	}

	for i, str := range s.messageStringMatrix[0] {
		s.messageStringMatrix[1][i] = createAlternativeString(str)
	}
	//fmt.Printf("długość Matrix:%d\n", len(s.messageStringMatrix[0][0]))
	s.messageBuferCells = make([][][]termbox.Cell, 2)

	for i := 0; i < 2; i++ {
		s.messageBuferCells[i] = make([][]termbox.Cell, 8)
		//var ixi *int
		for j := 0; j < 8; j++ {
			for _, xx := range s.messageStringMatrix[i][j] {

				var yy termbox.Cell
				vi := int(xx)

				yy.Ch = changeCharacter(vi)
				s.messageBuferCells[i][j] = append(s.messageBuferCells[i][j], yy)
				//		ixi = &ix
			}

		}
		//fmt.Printf("xx:% 3d\n", *ixi)

	}
	//colorization
	kolorIndex := 0
	for k := 0; k < len(s.messageBuferCells[0][0]); k = k + 8 {
		kolor := s.colorMessage[kolorIndex]
		for repeat := 0; repeat < 8; repeat++ {
			s.messageBuferCells[0][0][k+repeat].Fg = kolor
			s.messageBuferCells[0][1][k+repeat].Fg = kolor - 8
			s.messageBuferCells[0][2][k+repeat].Fg = kolor
			s.messageBuferCells[0][3][k+repeat].Fg = kolor - 8
			s.messageBuferCells[0][4][k+repeat].Fg = kolor
			s.messageBuferCells[0][5][k+repeat].Fg = kolor - 8
			s.messageBuferCells[0][6][k+repeat].Fg = kolor
			s.messageBuferCells[0][7][k+repeat].Fg = kolor - 8

		}
		kolorIndex++
		if kolorIndex >= len(s.colorMessage) {
			kolorIndex = kolorIndex - len(s.colorMessage)
		}

	}
	/*
		source := s.messageBuferCells[0][0][1]
		destination := s.messageBuferCells[1][0][0]
		n := copy(destination, source)
		printAt(50, s.yMax, termbox.ColorDefault, "liczba skopiowanych: %d", n)
		//copy colors to frame 1 Fg
	*/
	for r := 0; r <= 7; r++ {
		for i, object := range s.messageBuferCells[0][r][1:] {
			s.messageBuferCells[1][r][i].Fg = object.Fg

		}
	}
	s.xMax, s.yMax = termbox.Size()
	s.index = 0
}

func onExit() {
	exec.Command("/usr/bin/bash")
}

func changeCharacter(i int) rune {

	switch i {
	case 0:
		return ' '
	case 49:
		return '█'
	case 82:
		return '▐'
	case 76:
		return '▌'
	default:
		return ' '
	}

}

func printAt(x, y, c int, format string, args ...interface{}) error {

	termbox.SetCursor(x, y)
	termbox.SetCell(x, y, ' ', termbox.Attribute(c), termbox.ColorDefault)

	// Format and print the string
	str := fmt.Sprintf(format, args...)
	for i, char := range str {
		termbox.SetCell(x+i+1, y, char, termbox.Attribute(c), termbox.ColorDefault)
	}

	termbox.Flush()
	return nil
}
func (s *scroller) scrolling() {
	//var frame int

	indexing := s.index
	frame := s.frame
	lenXbuffer := len(s.messageBuferCells[1][0])
	/*	termbox.Close()
		fmt.Println("len", lenXbuffer, "indexing", indexing, "s.frame", s.frame)
		os.Exit(128)
	*/
	//lenXbuffer = lenXbuffer

	for x := s.lcx; x <= s.rcx; x++ {
		znak := s.messageBuferCells[frame][0][indexing]
		if znak.Ch != ' ' {
			termbox.SetCell(x, 0+s.lcy, znak.Ch, znak.Fg, termbox.ColorDefault)
		}
		znak = s.messageBuferCells[frame][1][indexing]
		if znak.Ch != ' ' {
			termbox.SetCell(x, 1+s.lcy, znak.Ch, znak.Fg, termbox.ColorDefault)
		}
		znak = s.messageBuferCells[frame][2][indexing]
		if znak.Ch != ' ' {
			termbox.SetCell(x, 2+s.lcy, znak.Ch, znak.Fg, termbox.ColorDefault)
		}
		znak = s.messageBuferCells[frame][3][indexing]
		if znak.Ch != ' ' {
			termbox.SetCell(x, 3+s.lcy, znak.Ch, znak.Fg, termbox.ColorDefault)
		}
		znak = s.messageBuferCells[frame][4][indexing]
		if znak.Ch != ' ' {
			termbox.SetCell(x, 4+s.lcy, znak.Ch, znak.Fg, termbox.ColorDefault)
		}
		znak = s.messageBuferCells[frame][5][indexing]
		if znak.Ch != ' ' {
			termbox.SetCell(x, 5+s.lcy, znak.Ch, znak.Fg, termbox.ColorDefault)
		}
		znak = s.messageBuferCells[frame][6][indexing]
		if znak.Ch != ' ' {
			termbox.SetCell(x, 6+s.lcy, znak.Ch, znak.Fg, termbox.ColorDefault)
		}
		znak = s.messageBuferCells[frame][7][indexing]
		if znak.Ch != ' ' {
			termbox.SetCell(x, 7+s.lcy, znak.Ch, znak.Fg, termbox.ColorDefault)
		}

		indexing++
		if indexing >= lenXbuffer {
			indexing = indexing - lenXbuffer
		}
	}

	/*switch frame {

	case 0:
		frame = 1
	case 1:
		frame = 0
	}*/
	//s.index++
	s.frame = s.frame + s.speed

looper:
	if s.frame >= 2 {
		s.frame = s.frame - 2
		s.index++
		if s.index >= len(s.messageBuferCells[0][0]) {
			s.index = s.index - len(s.messageBuferCells[0][0])
		}
	}
	if s.frame >= 2 {
		goto looper
	}
	/*
		for yyy := 0; yyy <= s.yMax; yyy++ {
			for xxx := 0; xxx <= s.xMax; xxx++ {
				termbox.SetBg(xxx, yyy, 0)
			}
		}*/

}

func findAllOccurrences(input string, pattern string) []int {
	var occurrences []int
	startIndex := 0

	for {
		index := strings.Index(input[startIndex:], pattern)
		if index == -1 {
			break
		}

		// Adjust the index based on the startIndex
		index += startIndex

		occurrences = append(occurrences, index)

		// Update the startIndex for the next search
		startIndex = index + 1
	}

	return occurrences
}

type chess struct {
	yy   int //init value for every frame
	xx   int //init value for every frame
	tx   int //weight chassField
	ty   int //height chassField
	c1   int //color 1
	c2   int //color 2
	xMax int
	yMax int
}

/*
func (c *chess) drawChessBoard(fx, fy int) {

		yy, xx := c.yy, c.xx
		lama1, lama2 := c.c1, c.c2
		drawingColor := lama1
		for y := 0; y <= fy; y++ {
			yy++
			if yy == c.ty {
				if drawingColor == lama1 {
					drawingColor = lama2
				} else {
					drawingColor = lama1
				}
				yy = 0
			}
			for x := 0; x <= fx; x++ {
				xx++
				if xx == c.tx {
					if drawingColor == lama2 {
						drawingColor = lama1

					} else {
						drawingColor = lama2
					}
					xx = 0
				}
				termbox.SetBg(x, y, termbox.Attribute(drawingColor))
			}
		}
	}
*/
func (c *chess) drawChessBoard(fxmin, fymin, fxmax, fymax int) {
	bgcolour := c.c1
	bgcolour2 := c.c2
	counter := 0
	turtle := true
	turtleoff := false
	for y := fymin; y <= fymax; y++ {
		counter++
		if counter == 8 {
			counter = 0
			bgcolour, bgcolour2 = bgcolour2, bgcolour
		}
		for x := fxmin; x <= fxmax; x++ {
			termbox.SetBg(x, y, termbox.Attribute(bgcolour))
		}
	}
	counter = 0
	for hx := fxmin; hx < fxmax; hx++ {
		counter++
		if counter == 16 {
			counter = 0
			turtle, turtleoff = turtleoff, turtle

		}
		for hy := fymin; hy < fymax; hy++ {
			if turtle {
				bgx := termbox.GetCell(hx, hy).Bg

				if bgx == termbox.Attribute(c.c1) {
					bgcolour = c.c2
					termbox.SetBg(hx, hy, termbox.Attribute(bgcolour))
				}
				if bgx == termbox.Attribute(c.c2) {
					bgcolour = c.c1
					termbox.SetBg(hx, hy, termbox.Attribute(bgcolour))
				}
			}
		}
	}

}

/*
	func (c *chess) drawChessBoard(fx, fy int) {
		for dy := 0; dy <= c.yMax; dy = dy + c.ty {
			for bly := dy; bly < dy+c.ty; bly++ {

				for dx := 0; dx <= c.xMax; dx++ {
					if dx > c.xMax {
						break
					}
					if bly > c.yMax {
						break
					}
					termbox.SetBg(dx, bly, termbox.Attribute(c.c1))

				}
			}
			dy = dy + c.ty
		}

		for ddx := 0; ddx <= c.xMax; ddx = ddx + c.tx {
			for blx := ddx; blx < ddx+c.tx; blx++ {
				for ddy := 0; ddy < c.yMax; ddy++ {
					if blx > c.xMax {

						break
					}
					if ddy > c.yMax {
						break
					}

					bgbg := termbox.GetCell(blx, ddy).Bg

					if bgbg == termbox.Attribute(c.c1) {
						termbox.SetBg(blx, ddy, termbox.Attribute(c.c2))
					} else {
						termbox.SetBg(blx, ddy, termbox.Attribute(c.c1))
					}
				}
			}

			ddx = ddx + c.tx
		}

}
*/
func main() {

	//defer onExit() ; this is a trick to get bash after information about property and responsibility but it doesn't work because goroutine dies after exit and new process should be created to parent process 1 tobe working fine.
	rand.Seed(time.Now().UnixNano())

	err := termbox.Init()
	if err != nil {
		log.Fatal("init() error", err)
	}

	var chessy chess
	chessy.xx = 0
	chessy.yy = 0
	chessy.ty = 8
	chessy.tx = 14
	chessy.c1 = 0
	chessy.c2 = 1
	chessy.xMax, chessy.yMax = termbox.Size()
	//create scroller struct
	var sc scroller
	var sc2, sc3 scroller
	info := "                                          Hi !            " +
		"this is the termbox library demo with a bit of tricks included.   " +
		"the library is responsible for a terminal display and its colorfull semi+graphics.   " +
		"sure it isn't for windows because it has no a colorfull console.   " +
		"I hope you understand that Linux rules ;-)   " +
		"you can get bored by this scroller.   " +
		"It may happen. KINDLY PLEASE DON'T COMMIT a SUICIDE for any reasons.   " +
		"It's not my intention to bring you to the death. I'm glad you're here.   " +
		"Don't be afraid if something is not as you wanted to be.   " +
		"My life is probably not so complicated as yours but I'm very empatic and I'm etusiast of humanbeings.   " +
		"Please remember! Don't give up. Don't loose your faith and your hope ! ! !   " +
		"Nothing's gona be wrong if your success is a bit postponed in time till you get enough of knowledge at the topics you're interested in.   " +
		"Be patient and systematic   ! ! !"
	info2 := "       10 9 8 7 6 5 4 3 2 1 * :) ;) XD       This is demo program coded by Andy! In case of any questions don't hasitate to ping me on linkedin. Go with Andy :) Cheers !    "
	//info2 := "               Andy Andy andy Andy !!! ! ! ! !! !! !!"
	sc.loadMessage(info2)
	sc2.loadMessage(info)
	sc3.loadMessage(info)

	sc.scrollerInit()
	sc2.scrollerInit()
	sc3.scrollerInit()

	termbox.Flush()

	defer termbox.Close()
	//load info to the message

	//	printAt(10, 10, 4, "x:%d y:%d", sc.xMax, sc.yMax)

	sc.lcx = 0 //left corner x coordinate
	sc2.lcx = 0
	sc3.lcx = 0
	sc2.lcy = (sc.yMax-8)/2 - 5 //left corner y coordinate
	sc.lcy = (sc.yMax - 8) / 2
	sc3.lcy = (sc.yMax-8)/2 + 5

	sc.rcx = sc.xMax //right corner x max coordinate
	sc2.rcx = sc.xMax
	sc3.rcx = sc.xMax
	//printAt(10, 20, 2, "to jest sc.rcx:%d\n", sc.rcx)
	sc.frame = 0
	sc2.frame = 0
	sc3.frame = 0
	sc.interval = 16 + 17
	//	sc2.interval = 30
	//	sc3.interval = 30
	sc.speed = 1
	sc2.speed = 2
	sc3.speed = 2

	indicesof := findAllOccurrences(info2, "Andy")
	for _, index := range indicesof {

		positionStart := index * 8
		positionEnd := positionStart + 8*len("Andy")
		gradient := []int{5, 13, 7, 15, 7, 13, 5, 5}
		for background := positionStart; background <= positionEnd; background++ {
			for y := 0; y <= 7; y++ {
				sc.messageBuferCells[0][y][background].Fg = termbox.Attribute(gradient[y])
				sc.messageBuferCells[1][y][background-1].Fg = termbox.Attribute(gradient[y])
			}
		}

	}
	//	...
	eventCh := make(chan termbox.Event)
	errCh := make(chan error)
	go func() {
		for {
			event := termbox.PollEvent()
			eventCh <- event
		}
	}()
	// Render all scrollers in one goroutine because flush cannot be executed in separate routines and parallization doesn't improve the performance
	//direction := 1
	go func() {
		for {
			sc.scrolling()
			sc2.scrolling()
			sc3.scrolling()

			chessy.drawChessBoard(0, 0, sc.xMax, sc.yMax)
			//termbox.Sync()

			termbox.Flush()
			termbox.Clear(0, 0)

			/*		chessy.tx = chessy.tx + 2*direction
					chessy.ty = chessy.ty + 1*direction
					if chessy.tx == 2 {
						direction = 1
					}
					if chessy.tx >= 16 {
						direction = -1

					}
			*/
			/*		chessy.yy--
					if chessy.yy == -chessy.ty {
						chessy.yy = chessy.yy + chessy.ty
					}*/
			// termbox.Clear(termbox.ColorDefault, 0) ; just testing. Now it's unnecessary

			time.Sleep(time.Millisecond * time.Duration(sc.interval))
		}
	}()

	//the below code is for exit from this program by key press [ESC] only
	for {
		select {
		case event := <-eventCh:
			// Handle the event
			switch event.Type {
			case termbox.EventKey:
				// Handle key press event
				if event.Key == termbox.KeyEsc {
					// Exit the loop if the Escape key is pressed
					return
				}
				// Handle other key events as needed
			case termbox.EventMouse:
				// Handle mouse event
				// ...
			case termbox.EventResize:
				termbox.Close()
				os.Exit(0)
			}

			// Perform other operations based on the event if any ...

		case err := <-errCh:
			// Handle errors, if any
			panic(err)

		default:
			// Perform other non-blocking operations
			// ...

		}
	}
}
