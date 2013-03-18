package main

import (
	"github.com/Nightgunner5/agwycctmiatstsycbmrm/state"
	"github.com/nsf/termbox-go"
	"sort"
	"time"
	"unsafe"
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

type sortItems []*state.Item

func (s sortItems) Len() int {
	return len(s)
}

func (s sortItems) Less(i, j int) bool {
	si, sj := s[i], s[j]

	if ci, cj := si.Category.Category(), sj.Category.Category(); ci != cj {
		return ci < cj
	}

	if si.Category != sj.Category {
		return si.Category < sj.Category
	}

	if si.Name != sj.Name {
		return si.Name < sj.Name
	}

	return uintptr(unsafe.Pointer(si.Components)) < uintptr(unsafe.Pointer(sj.Components))
}

func (s sortItems) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type gameUI struct {
	parent context
	state  *state.State

	inventoryScroll int
	inventoryCache  sortItems
}

func (g *gameUI) paint(w, h int) {
	{
		var items sortItems
		g.state.Lock()
		for i, count := range g.state.Inventory {
			item := i
			for i := uint(0); i < count; i++ {
				items = append(items, &item)
			}
		}
		g.state.Unlock()

		sort.Sort(items)

		g.inventoryCache = items

		y := gameHeight + 1 - g.inventoryScroll
		for _, i := range items {
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
	}

	for x := w - gameWidth; x < w; x++ {
		for y := 1; y <= gameHeight; y++ {
			termbox.SetCell(x, y, ' ', termbox.ColorDefault, termbox.ColorDefault)
		}
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

func (g *gameUI) key(k termbox.Key) {
	switch k {
	case termbox.KeyArrowDown:
		if g.inventoryScroll < len(g.inventoryCache)-1 {
			g.inventoryScroll++
			repaint()
		}

	case termbox.KeyArrowUp:
		if g.inventoryScroll > 0 {
			g.inventoryScroll--
			repaint()
		}
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
		g.state.AdvanceDay()
		repaint()

		time.Sleep(g.state.DayLength)
	}
}
