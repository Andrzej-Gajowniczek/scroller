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

	eventCh := make(chan termbox.Event)
	errCh := make(chan error)

	go func() {
		for {
			event, err := termbox.PollEvent()
			if err != nil {
				errCh <- err
				return
			}
			eventCh <- event
		}
	}()

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
			}

			// Perform other operations based on the event
			// ...

		case err := <-errCh:
			// Handle errors, if any
			panic(err)

		default:
			// Perform other non-blocking operations
			// ...
		}
	}
}
