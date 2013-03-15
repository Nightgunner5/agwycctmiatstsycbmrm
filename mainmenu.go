package main

import (
	"github.com/Nightgunner5/agwycctmiatstsycbmrm/util"
	"github.com/nsf/termbox-go"
)

var (
	mainmenuTitle = util.SplitForUI("a game where you command crafters to make items and then sell them so you can buy more raw materials")

	mainmenuVersion = []rune("v" + util.Version)

	/*
		13:41 < Nightgunner5> guys I need a word that means the same thing as "normal difficulty" but is moderately funny
		13:42 < Nightgunner5> I'll even take punny if that's all you can do
		13:42 < avagrant> Average
		13:42 < Nightgunner5> so far my difficulties are "boring", , "rude", and "the game mode that sets everything on fire"
		13:42 < Zed_Beeonova> vanilla
		13:43 < avagrant> 6 inches
		13:43 < avagrant> (it's normal)
		13:43 < Nightgunner5> should I call it autistic to mess with everyone?
		13:43 < Broodlines> 6 inches
		13:44 < Nightgunner5> okay, it looks like 6 inches is winning
		13:44 <@AngriestIBM> autist and non-autist modes
		13:44 < Chase_Mobile> Autustic should be hard
		13:44 < Nightgunner5> autistic before or after rude?
		13:44 < Chase_Mobile> Or the highest difficulty
		13:44 <@AngriestIBM> grump difficulty
	*/
	mainmenuGameModes = [][][]rune{
		util.SplitForUI("grump"),
		util.SplitForUI("six inches"),
		util.SplitForUI("rude"),
		util.SplitForUI("the game mode that sets everything on fire"),
		util.SplitForUI("autistic"),
	}
)

type mainMenu struct{}

func (m *mainMenu) paint(w, h int) {
	x, y := w/8, 2

	for _, word := range mainmenuTitle {
		if x+len(word)+w/8 >= w {
			x, y = w/8, y+1
		}

		for i, r := range word {
			termbox.SetCell(x+i, y, r, termbox.ColorDefault|termbox.AttrBold, termbox.ColorDefault)
		}

		x += len(word) + 1
	}

	for i, r := range mainmenuVersion {
		termbox.SetCell(x+i, y, r, termbox.ColorDefault, termbox.ColorDefault)
	}

	y++
	for l, line := range mainmenuGameModes {
		x, y = w/8, y+1
		termbox.SetCell(x-w/16-1, y, rune('1'+l), termbox.ColorDefault|termbox.AttrBold, termbox.ColorDefault)

		for _, word := range line {
			if x+len(word)+w/8 >= w {
				x, y = w/8, y+1
			}

			for i, r := range word {
				termbox.SetCell(x+i, y, r, termbox.ColorDefault, termbox.ColorDefault)
			}

			x += len(word) + 1
		}
	}
}

func (m *mainMenu) char(ch rune) {
	switch ch {
	case '1', '2', '3', '4', '5':
		paintCtx = &difficultyConfirm{
			parent:     m,
			difficulty: int(ch - '1'),
		}

		repaint()
	}
}

func (m *mainMenu) key(k termbox.Key) {
}
