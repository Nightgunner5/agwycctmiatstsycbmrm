package state

import (
	"math/rand"
)

type Task uint16

const (
	TaskNone Task = iota
	TaskCombineItems
	TaskMakePlanks
	TaskMakeBlocks
	TaskMakeSheets
	TaskSmeltOres
	TaskSmeltNuggets
	TaskSpinThread
	TaskWeaveCloth
	NumTasks
)

func (t Task) String() string {
	switch t {
	case TaskNone:
		return "none"

	case TaskCombineItems:
		return "combine items"

	case TaskMakePlanks:
		return "make planks"

	case TaskMakeBlocks:
		return "make blocks"

	case TaskMakeSheets:
		return "make sheets"

	case TaskSmeltOres:
		return "smelt ores"

	case TaskSmeltNuggets:
		return "smelt nuggets"

	case TaskSpinThread:
		return "spin thread"

	case TaskWeaveCloth:
		return "weave cloth"
	}

	return "ERROR"
}

func (t Task) complete(w *Worker, s *State) {
	// TODO: skill levels

	switch t {
	case TaskCombineItems:
		for i, ii := range w.Inventory {
			for j, ij := range w.Inventory {
				if j <= i || !ii.Category.CanCombine(ij.Category) {
					continue
				}

				ii = ii.Clone()
				ij = ii.Combine(ij.Clone())

				if ij == nil {
					w.Inventory = append(w.Inventory[:j], w.Inventory[j+1:]...)
				} else {
					w.Inventory[j] = ij
				}

				w.Inventory[i] = ii

				s.incrementScore(uint64(ii.Quality)>>4 + 1)

				// TODO: log task success

				for _, item := range w.Inventory[:i] {
					s.addItem(item)
				}
				w.Inventory = w.Inventory[i:]

				return
			}
		}

	case TaskMakePlanks:
		for i, ii := range w.Inventory {
			if ii.Category&0x000fff00 != Log {
				continue
			}

			ii = ii.Clone()
			ii.Category = ii.Category&0xfff000ff | Plank
			// TODO: adjust quality (and size?)
			w.Inventory[i] = ii

			s.incrementScore(uint64(ii.Quality)>>4 + 1)

			// TODO: log task success

			for _, item := range w.Inventory[:i] {
				s.addItem(item)
			}
			w.Inventory = w.Inventory[i:]
			return
		}

	case TaskMakeBlocks:
		for i, ii := range w.Inventory {
			if ii.Category&0x000fff00 != Log && ii.Category&0x000fff00 != Ingot {
				continue
			}

			ii = ii.Clone()
			ii.Category = ii.Category&0xfff000ff | Block
			// TODO: adjust quality (and size?)
			w.Inventory[i] = ii

			s.incrementScore(uint64(ii.Quality)>>4 + 1)

			// TODO: log task success

			for _, item := range w.Inventory[:i] {
				s.addItem(item)
			}
			w.Inventory = w.Inventory[i:]
			return
		}

	case TaskMakeSheets:
		for i, ii := range w.Inventory {
			if ii.Category&0x000fff00 != Log && ii.Category&0x000fff00 != Plank && ii.Category&0x000fff00 != Ingot {
				continue
			}

			ii = ii.Clone()
			ii.Category = ii.Category&0xfff000ff | Sheet
			// TODO: adjust quality (and size?)
			w.Inventory[i] = ii

			s.incrementScore(uint64(ii.Quality)>>4 + 1)

			// TODO: log task success

			for _, item := range w.Inventory[:i] {
				s.addItem(item)
			}
			w.Inventory = w.Inventory[i:]
			return
		}

	case TaskSmeltOres:
		for i, ii := range w.Inventory {
			if ii.Category&0x000fff00 != Ore {
				continue
			}

			ii = ii.Clone()
			ii.Category = ii.Category&0xfff000ff | Ingot
			// TODO: adjust quality (and size?)
			w.Inventory[i] = ii

			s.incrementScore(uint64(ii.Quality)>>4 + 1)

			// TODO: log task success

			for _, item := range w.Inventory[:i] {
				s.addItem(item)
			}
			w.Inventory = w.Inventory[i:]
			return
		}

	case TaskSmeltNuggets:
		for i, ii := range w.Inventory {
			if ii.Category&0x000fff00 != Nugget {
				continue
			}

			ii = ii.Clone()
			ii.Category = ii.Category&0xfff000ff | Ingot
			// TODO: adjust quality (and size?)
			w.Inventory[i] = ii

			s.incrementScore(uint64(ii.Quality)>>4 + 1)

			// TODO: log task success

			for _, item := range w.Inventory[:i] {
				s.addItem(item)
			}
			w.Inventory = w.Inventory[i:]
			return
		}

	case TaskSpinThread:
		for i, ii := range w.Inventory {
			if ii.Category&0xffffff00 != Wool {
				continue
			}

			ii = ii.Clone()
			ii.Category = ii.Category | Thread
			// TODO: adjust quality (and size?)
			w.Inventory[i] = ii

			s.incrementScore(uint64(ii.Quality)>>4 + 1)

			// TODO: log task success

			for _, item := range w.Inventory[:i] {
				s.addItem(item)
			}
			w.Inventory = w.Inventory[i:]
			return
		}

	case TaskWeaveCloth:
		for i, ii := range w.Inventory {
			if ii.Category&0x000fff00 != Thread {
				continue
			}

			ii = ii.Clone()
			ii.Category = ii.Category&0xfff000ff | Cloth
			// TODO: adjust quality (and size?)
			w.Inventory[i] = ii

			s.incrementScore(uint64(ii.Quality)>>4 + 1)

			// TODO: log task success

			for _, item := range w.Inventory[:i] {
				s.addItem(item)
			}
			w.Inventory = w.Inventory[i:]
			return
		}
	}
	// give all items back to the supply pile
	for _, item := range w.Inventory {
		s.addItem(item)
	}
	w.Inventory = nil

	// TODO: log task failure
	w.Task = TaskNone
}

type Worker struct {
	Task      Task
	Progress  uint16
	Inventory []*Item
}

func (w *Worker) advanceDay(current uint64, s *State) {
}

func (w *Worker) advanceTime(s *State) {
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
