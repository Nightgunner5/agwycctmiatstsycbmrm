package main

import (
	"fmt"
	"github.com/Nightgunner5/agwycctmiatstsycbmrm/state"
	"github.com/nsf/termbox-go"
	"time"
)

var (
	gameDay   = []rune("Day")
	gameScore = []rune("Score")
	gameCash  = []rune("Cash")

	gameWidth = func(l ...int) int {
		w := -1

		for _, i := range l {
			if i > w {
				w = i
			}
		}

		return w + 22
	}(len(gameDay), len(gameScore), len(gameCash))
	gameHeight = 3 + 1
)

type gameUI struct {
	parent context
	state  *state.State

	scoreInterp uint64

	currentWorker   int
	inventoryScroll int
	selectedItem    *state.Item
}

func (g *gameUI) paint(w, h int) {
	g.state.Lock()
	for i, r := range []rune(fmt.Sprintf("Worker %d / %d", g.currentWorker+1, len(g.state.Workers))) {
		termbox.SetCell(i, 1, r, termbox.ColorDefault|termbox.AttrBold, termbox.ColorDefault)
	}

	worker := g.state.Workers[g.currentWorker]

	y := 3

	termbox.SetCell(0, y, 'G', termbox.ColorWhite|termbox.AttrBold, termbox.ColorBlack)
	for x, r := range []rune("give item") {
		termbox.SetCell(x+2, y, r, termbox.ColorDefault, termbox.ColorDefault)
	}
	y++

	termbox.SetCell(0, y, 'T', termbox.ColorWhite|termbox.AttrBold, termbox.ColorBlack)
	for x, r := range []rune(fmt.Sprintf("current task: %s (%d%%)", worker.Task, worker.Progress/((1<<16-1)/100))) {
		termbox.SetCell(x+2, y, r, termbox.ColorDefault, termbox.ColorDefault)
	}
	y++

	for _, i := range worker.Inventory {
		y++
		for x, r := range []rune(i.String()) {
			termbox.SetCell(x, y, r, termbox.ColorDefault, termbox.ColorDefault)
		}
	}

	y = gameHeight + 1 - g.inventoryScroll
	if g.inventoryScroll < len(g.state.Inventory) {
		g.selectedItem = g.state.Inventory[g.inventoryScroll]
	} else {
		g.selectedItem = nil
	}
	for _, i := range g.state.Inventory {
		if y > 0 {
			for x, r := range []rune(i.String()) {
				if y == gameHeight+1 {
					termbox.SetCell(x+w*3/5, y, r, termbox.ColorDefault|termbox.AttrBold, termbox.ColorDefault)
				} else {
					termbox.SetCell(x+w*3/5, y, r, termbox.ColorDefault, termbox.ColorDefault)
				}
			}
		}
		y++
	}

	g.state.Unlock()

	for i := 1; i < gameHeight; i++ {
		termbox.SetCell(w-21, i, ' ', termbox.ColorDefault, termbox.ColorDefault)
	}

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

	score := g.interpScore()

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

func (g *gameUI) char(r rune) {
	switch r {
	case 'g', 'G':
		if g.state.RemoveItem(g.selectedItem) {
			g.state.Lock()
			defer g.state.Unlock()

			w := g.state.Workers[g.currentWorker]
			w.Inventory = append(w.Inventory, g.selectedItem)

			repaint()
		}

	case 't', 'T':
		g.state.Lock()
		defer g.state.Unlock()

		w := g.state.Workers[g.currentWorker]
		// TODO: some kind of sub-panel for choosing a task
		w.Task++
		if w.Task == state.NumTasks {
			w.Task = state.TaskNone
		}
		w.Progress = 0

		repaint()
	}
}

func (g *gameUI) key(k termbox.Key) {
	switch k {
	case termbox.KeyArrowUp:
		if g.inventoryScroll > 0 {
			g.inventoryScroll--
			repaint()
		}

	case termbox.KeyArrowDown:
		g.state.Lock()
		if g.inventoryScroll < len(g.state.Inventory)-1 {
			g.inventoryScroll++
			repaint()
		}
		g.state.Unlock()

	case termbox.KeyArrowLeft:
		if g.currentWorker > 0 {
			g.currentWorker--
			repaint()
		}

	case termbox.KeyArrowRight:
		g.state.Lock()
		if g.currentWorker < len(g.state.Workers)-1 {
			g.currentWorker++
			repaint()
		}
		g.state.Unlock()
	}
}

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
		g.state.AdvanceTime()
		repaint()

		time.Sleep(g.state.DayLength / 100)
	}
}

func (g *gameUI) interpScore() uint64 {
	s := g.state.GetScore()

	switch {
	case g.scoreInterp > s:
		g.scoreInterp = s

		fallthrough
	case g.scoreInterp == s:
		return s

	case g.scoreInterp >= s-10 || g.scoreInterp < 10:
		g.scoreInterp += 1
		return g.scoreInterp

	case g.scoreInterp >= s-100 || g.scoreInterp < 100:
		g.scoreInterp += 10
		return g.scoreInterp

	case g.scoreInterp >= s-1000 || g.scoreInterp < 1000:
		g.scoreInterp += 100
		return g.scoreInterp

	case g.scoreInterp >= s-10000 || g.scoreInterp < 10000:
		g.scoreInterp += 1000
		return g.scoreInterp
	}

	g.scoreInterp += 10000
	return g.scoreInterp
}
