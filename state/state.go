package state

import (
	"sync"
	"time"
)

type State struct {
	DayLength   time.Duration
	Day         uint64
	Cash        uint64
	Score       uint64
	scoreInterp uint64
	Workers     []*Worker
	Inventory   map[Item]uint
	sync.Mutex
}

func (s *State) Init() {
	s.Lock()

	s.Workers = append(s.Workers, &Worker{s})
	s.Workers = append(s.Workers, &Worker{s})
	s.Workers = append(s.Workers, &Worker{s})
	s.Workers = append(s.Workers, &Worker{s})
	s.Workers = append(s.Workers, &Worker{s})
	s.Workers = append(s.Workers, &Worker{s})

	s.Inventory = make(map[Item]uint)

	s.Unlock()
}

func (s *State) GetDay() (date uint64) {
	s.Lock()
	date = s.Day
	s.Unlock()
	return
}

func (s *State) ChangeCash(amount int64) bool {
	s.Lock()
	defer s.Unlock()

	if amount < 0 && s.Cash < uint64(-amount) {
		return false
	}

	s.Cash = uint64(int64(s.Cash) + amount)

	return true
}

func (s *State) GetCash() (cash uint64) {
	s.Lock()
	cash = s.Cash
	s.Unlock()
	return
}

func (s *State) IncrementScore(amount uint64) {
	s.Lock()
	s.Score += amount
	s.Unlock()
}

func (s *State) GetRealScore() (score uint64) {
	s.Lock()
	score = s.Score
	s.Unlock()
	return
}

func (s *State) GetInterpolatedScore() (uint64, bool) {
	s.Lock()
	defer s.Unlock()

	switch {
	case s.scoreInterp > s.Score:
		s.scoreInterp = s.Score

		fallthrough
	case s.scoreInterp == s.Score:
		return s.Score, false

	case s.scoreInterp >= s.Score-10:
		s.scoreInterp += 1
		return s.scoreInterp, true

	case s.scoreInterp >= s.Score-100:
		s.scoreInterp += 10
		return s.scoreInterp, true

	case s.scoreInterp >= s.Score-1000:
		s.scoreInterp += 100
		return s.scoreInterp, true

	case s.scoreInterp >= s.Score-10000:
		s.scoreInterp += 1000
		return s.scoreInterp, true
	}
	s.scoreInterp += 10000
	return s.scoreInterp, true
}

func (s *State) AdvanceDay() {
	s.Lock()

	s.Day++

	for _, w := range s.Workers {
		w.advanceDay(s.Day)
	}

	s.Unlock()
}

func (s *State) AddItem(i Item) {
	s.ChangeItemCount(i, 1)
}

func (s *State) RemoveItem(i Item) bool {
	return s.ChangeItemCount(i, -1)
}

func (s *State) ChangeItemCount(i Item, delta int) bool {
	s.Lock()
	defer s.Unlock()

	if delta < 0 && s.Inventory[i] < uint(-delta) {
		return false
	}

	if s.Inventory == nil {
		s.Inventory = make(map[Item]uint)
	}

	s.Inventory[i] = uint(int(s.Inventory[i]) + delta)

	return true
}

func (s *State) GetItemCount(i Item) (count uint) {
	s.Lock()
	count = s.Inventory[i]
	s.Unlock()
	return

}
