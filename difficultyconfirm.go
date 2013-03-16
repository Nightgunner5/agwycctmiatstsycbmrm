package main

import (
	"fmt"
	"github.com/Nightgunner5/agwycctmiatstsycbmrm/state"
	"github.com/Nightgunner5/agwycctmiatstsycbmrm/util"
	"github.com/nsf/termbox-go"
)

var (
	difficultyDescription = [][][]rune{
		util.SplitForUI("Grump mode, for when you're down in the grumps. Alternatively, it can be used for your gramps. Grandpa! No! That is not the right button. Grandpa, stop. You're horrible at this. No Grandpa, you can't check your AOL on this. It's a video game. No, that is not a film projector."),
		util.SplitForUI("Six inches is a normal mode for normal people. Alternatively, you can replace 'normal' with 'boring'. Six inches is a boring mode for boring people. See what I did there? Get it? Wh- why are you walking away? Guys?"),
		util.SplitForUI("This game mode is rude. It does rude things and brings its rude friends to be rude near you. Recommended for torture-seekers. And rude people."),
		util.SplitForUI("AAAAAAAAAAAAAAAAAAAAAAAAA PUT IT OUT PUT IT OUT PUT IT OUT AAAAAAAAAAAAAAAAAA MAKE IT STOP IT BURNS AAAAAAAAAAAAAAAAAAAA HELP ME SOMEONE HELP AAAAAAAAAAAAAAAA HELP ME EVERYTHING IS ON FIRE"),
		util.SplitForUI("Hey guys let's make a game where you play a game and then there's space and stuff and then you get items or something but really the whole thing is a metaphor for the triviality of life. Ooh, the pizza is here."),
	}

	difficultyStartingMoney = []rune("Starting Cash")
	difficultyGimmick       = []rune("Gimmick")

	difficultyTuning = []struct {
		startingmoney uint64
		gimmick       [][]rune
	}{
		{
			startingmoney: 1000000000,
			gimmick:       util.SplitForUI("grandpa-friendly"),
		},
		{
			startingmoney: 10000000,
			gimmick:       util.SplitForUI("boring"),
		},
		{
			startingmoney: 1000000,
			gimmick:       util.SplitForUI("the game insults you"),
		},
		{
			startingmoney: 100000,
			gimmick:       util.SplitForUI("everything is on fire"),
		},
		{
			startingmoney: 250,
			gimmick:       util.SplitForUI("slave labor required"),
		},
	}

	difficultyStartHint = util.SplitForUI("Push B to begin.")
)

type difficultyConfirm struct {
	parent     context
	difficulty int
}

func (d *difficultyConfirm) paint(w, h int) {
	for i, r := range backHint {
		termbox.SetCell(w-len(backHint)+i, 1, r, termbox.ColorDefault, termbox.ColorDefault)
	}

	x, y := 1, 3

	for _, word := range difficultyDescription[d.difficulty] {
		if x+len(word)+1 >= w {
			x, y = 1, y+1
		}

		for i, r := range word {
			termbox.SetCell(x+i, y, r, termbox.ColorDefault, termbox.ColorDefault)
		}

		x += len(word) + 1
	}

	tuning := &difficultyTuning[d.difficulty]
	x, y = 16, y+2
	for i, r := range difficultyStartingMoney {
		termbox.SetCell(x-len(difficultyStartingMoney)+i-1, y, r, termbox.ColorDefault|termbox.AttrBold, termbox.ColorDefault)
	}
	for i, r := range []rune(fmt.Sprintf("%d", tuning.startingmoney)) {
		termbox.SetCell(x+i, y, r, termbox.ColorDefault, termbox.ColorDefault)
	}

	x, y = 16, y+1
	for i, r := range difficultyGimmick {
		termbox.SetCell(x-len(difficultyGimmick)+i-1, y, r, termbox.ColorDefault|termbox.AttrBold, termbox.ColorDefault)
	}
	for _, word := range tuning.gimmick {
		if x+len(word)+1 >= w {
			x, y = 16, y+1
		}

		for i, r := range word {
			termbox.SetCell(x+i, y, r, termbox.ColorDefault, termbox.ColorDefault)
		}

		x += len(word) + 1
	}

	x, y = 16, y+2
	for _, word := range difficultyStartHint {
		if x+len(word)+1 >= w {
			x, y = 16, y+1
		}

		for i, r := range word {
			termbox.SetCell(x+i, y, r, termbox.ColorDefault, termbox.ColorDefault)
		}

		x += len(word) + 1
	}
}

func (d *difficultyConfirm) char(r rune) {
	switch r {
	case 'q', 'Q':
		paintCtx = d.parent
		repaint()
	case 'b', 'B':
		game := &gameUI{d.parent, &state.State{
			Cash: difficultyTuning[d.difficulty].startingmoney,
		}}
		paintCtx = game
		go game.process()
		repaint()
	}
}

func (d *difficultyConfirm) key(k termbox.Key) {
}
