package main

import (
	"github.com/Nightgunner5/agwycctmiatstsycbmrm/state"
	"github.com/nsf/termbox-go"
	"time"
)

var (
	gameDay   = []rune("Day")
	gameScore = []rune("Score")
	gameCash  = []rune("Cash")
)

type gameUI struct {
	parent context
	state  *state.State
}

func (g *gameUI) paint(w, h int) {
	day := g.state.GetDay()
	for i := 0; i < 20; i++ {
		if day == 0 && i != 0 {
			termbox.SetCell(w-1-i, 1, '0', termbox.ColorDefault, termbox.ColorDefault)
		} else {
			termbox.SetCell(w-1-i, 1, rune('0'+(day%10)), termbox.ColorDefault|termbox.AttrBold, termbox.ColorDefault)
		}
		day /= 10
	}

	for i, r := range gameDay {
		termbox.SetCell(w-21-len(gameDay)+i, 1, r, termbox.ColorDefault|termbox.AttrBold, termbox.ColorDefault)
	}

	score, shouldRepaint := g.state.GetInterpolatedScore()

	if shouldRepaint {
		go func() {
			time.Sleep(time.Second / 10)
			repaint()
		}()
	}

	for i := 0; i < 20; i++ {
		if score == 0 && i != 0 {
			termbox.SetCell(w-1-i, 2, '0', termbox.ColorDefault, termbox.ColorDefault)
		} else {
			termbox.SetCell(w-1-i, 2, rune('0'+(score%10)), termbox.ColorDefault|termbox.AttrBold, termbox.ColorDefault)
		}
		score /= 10
	}

	for i, r := range gameScore {
		termbox.SetCell(w-21-len(gameScore)+i, 2, r, termbox.ColorDefault|termbox.AttrBold, termbox.ColorDefault)
	}

	cash := g.state.GetCash()

	for i := 0; i < 20; i++ {
		if cash == 0 && i != 0 {
			termbox.SetCell(w-1-i, 3, '0', termbox.ColorDefault, termbox.ColorDefault)
		} else {
			termbox.SetCell(w-1-i, 3, rune('0'+(cash%10)), termbox.ColorDefault|termbox.AttrBold, termbox.ColorDefault)
		}
		cash /= 10
	}

	for i, r := range gameCash {
		termbox.SetCell(w-21-len(gameCash)+i, 3, r, termbox.ColorDefault|termbox.AttrBold, termbox.ColorDefault)
	}
}

func (g *gameUI) char(r rune) {}

func (g *gameUI) key(k termbox.Key) {}

func (g *gameUI) process() {
	go g.advanceDate()

	//g.state.IncrementScore(100000)
	//repaint()

	for {
		time.Sleep(time.Second / 10)
	}
}

func (g *gameUI) advanceDate() {
	for {
		time.Sleep(time.Second * 5)
		g.state.AdvanceDay()
		repaint()
	}
}
