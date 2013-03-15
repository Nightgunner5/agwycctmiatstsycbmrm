package main

import (
	"github.com/nsf/termbox-go"
)

func main() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	var shouldExit bool

	go painter()

	for {
		e := termbox.PollEvent()
		switch e.Type {
		case termbox.EventError:
			panic(e.Err) // TODO: better error handling?

		case termbox.EventResize:
			repaint()

		case termbox.EventKey:
			if e.Key != termbox.KeyEsc {
				shouldExit = false
			}

			if e.Key == 0 {
				paintLock.Lock()
				paintCtx.char(e.Ch)
				paintLock.Unlock()
			} else if e.Key == termbox.KeyEsc {
				if shouldExit {
					return
				} else {
					shouldExit = true
				}
			} else {
				paintLock.Lock()
				paintCtx.key(e.Key)
				paintLock.Unlock()
			}
		}
	}
}
