package main

import (
	"github.com/nsf/termbox-go"
)

type mainMenu struct{}

func (m *mainMenu) parent() context {
	return m
}

func (m *mainMenu) paint(w, h int) {
}

func (m *mainMenu) char(ch rune) {
}

func (m *mainMenu) key(termbox.Key) {
}
