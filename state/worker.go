package state

import (
	"math/rand"
)

type Task uint16

const (
	TaskNone Task = iota
	TaskCombineItems
	NumTasks
)

func (t Task) String() string {
	switch t {
	case TaskNone:
		return "none"

	case TaskCombineItems:
		return "combine items"
	}

	return "ERROR"
}

func (t Task) complete(w *Worker, s *State) {
	// TODO: skill levels

	switch t {
	case TaskCombineItems:
		for i, ii := range w.Inventory {
			if !ii.Category.Category().CanCombine() {
				continue
			}

			if ii.Category&0x000000ff >= Gigantic {
				s.addItem(ii)
				w.Inventory = append(w.Inventory[:i], w.Inventory[i+1:]...)
				return
			}

			ci := ii.Category & 0xffffff00
			si := ci&0x000fff00 == Scrap
			if si {
				ci &= 0xfff00000
			}

			for j, ij := range w.Inventory[i+1:] {
				cj := ij.Category & 0xffffff00
				sj := cj&0x000fff00 == Scrap
				if sj {
					cj &= 0xfff00000
				}

				if si && sj {
					continue
				}

				if ci == cj || (si && ci == cj&0xfff00000) || (sj && ci&0xfff00000 == cj) {

					sizeHave := ii.Category & 0x000000ff
					sizeTaken := 0xff - sizeHave

					if sizeTaken > ij.Category&0x000000ff {
						sizeTaken = ij.Category & 0x000000ff
					}

					if ij.Category&0x000000ff == sizeTaken {
						w.Inventory = append(w.Inventory[:j], w.Inventory[j+1:]...)
					} else {
						ij = new(Item)
						*ij = *w.Inventory[j]
						w.Inventory[j] = ij

						ij.Category = ij.Category&0xffffff00 | (ij.Category&0x000000ff - sizeTaken)
					}

					ii = new(Item)
					*ii = *w.Inventory[i]
					w.Inventory[i] = ii

					if si {
						ii.Category = ij.Category&0xffffff00 | (sizeHave + sizeTaken)
					} else {
						ii.Category = ii.Category&0xffffff00 | (sizeHave + sizeTaken)
					}

					// TODO: lose some quality from stitching two things together
					ii.Quality = uint16(uint32(ii.Quality)*uint32(sizeHave)+uint32(ij.Quality)*uint32(sizeTaken)) >> 8

					s.incrementScore(uint64(ii.Quality))

					if sizeHave+sizeTaken >= Gigantic {
						s.addItem(ii)
						w.Inventory = append(w.Inventory[:i], w.Inventory[i+1:]...)
					}

					return
				}
			}
		}

		// TODO: log task failure
		w.Task = TaskNone
	}
}

type Worker struct {
	Task      Task
	Progress  uint16
	Inventory []*Item
}

func (w *Worker) advanceDay(current uint64, s *State) {
}

func (w *Worker) advanceDayPercent(s *State) {
	if w.Progress == ^uint16(0) {
		w.Progress = 0
		w.Task.complete(w, s)
	}

	if w.Task == TaskNone {
		w.Progress = 0

		// TODO: boredom
		return
	}

	delta := uint16(rand.Intn(0x0100) + 0x0010)

	if delta+w.Progress < w.Progress {
		w.Progress = ^uint16(0)
	} else {
		w.Progress += delta
	}
}
