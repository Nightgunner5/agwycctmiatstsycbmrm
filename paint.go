package main

import (
	"github.com/nsf/termbox-go"
	"sync"
)

type context interface {
	paint(w, h int)
	char(ch rune)
	key(termbox.Key)
}

var (
	paintCtx  context = new(mainMenu)
	paintLock sync.Mutex
	paintCond = make(chan struct{}, 1)
)

func repaint() {
	select {
	case paintCond <- struct{}{}:
	default:
	}
}

var (
	backHint = []rune("To go back, push Q.")
)

func painter() {
	var (
		gameTitle = []rune("AGWYCCTMIATSTSYCBMRM")
		exitHint  = []rune("To exit, push ESC twice.")
	)

	for {
		paintLock.Lock()
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		w, h := termbox.Size()
		_ = h

		// color top bar
		for x := 0; x < w; x++ {
			termbox.SetCell(x, 0, ' ', termbox.ColorBlack, termbox.ColorWhite)
		}

		// game title, left top
		if w > len(gameTitle)*2+len(exitHint) {
			for i, r := range gameTitle {
				termbox.SetCell(i*2, 0, r, termbox.ColorBlack|termbox.AttrBold, termbox.ColorWhite)
				termbox.SetCell(i*2+1, 0, '.', termbox.ColorBlack|termbox.AttrBold, termbox.ColorWhite)
			}
		} else {
			for i, r := range gameTitle {
				termbox.SetCell(i, 0, r, termbox.ColorBlack|termbox.AttrBold, termbox.ColorWhite)
			}
		}

		// "To exit, push ESC twice.", right top
		for i, r := range exitHint {
			termbox.SetCell(w-len(exitHint)+i, 0, r, termbox.ColorBlack, termbox.ColorWhite)
		}

		paintCtx.paint(w, h)

		termbox.Flush()

		paintLock.Unlock()

		<-paintCond
	}
}
