package state

import (
	"math/rand"
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
	for i := 0; i < 20; i++ {
		s.Inventory[randomStartingItem()] = uint(rand.Intn(5) + 1)
	}

	s.Unlock()
}

func randomStartingItem() (i Item) {
	i.Category = ItemCategory(rand.Intn(int(Large-Small+9))) + Small

	switch rand.Intn(4) {
	case 0: // metal
		if rand.Intn(3) == 0 {
			i.Category |= Ingot
		} else {
			i.Category |= Ore
		}

		switch rand.Intn(5) {
		case 0:
			i.Category |= Copper
		case 1:
			i.Category |= Tin
		case 2:
			i.Category |= Bronze
		case 3:
			i.Category |= Lead
		case 4:
			i.Category |= Iron
		}

	case 1: // wood
		if rand.Intn(2) == 0 {
			i.Category |= Log
		} else {
			i.Category |= Plank
		}

		switch rand.Intn(4) {
		case 0:
			i.Category |= Birch
		case 1:
			i.Category |= Pine
		case 2:
			i.Category |= Maple
		case 3:
			i.Category |= Walnut
		}

	case 2: // generic craft materials
		switch rand.Intn(3) {
		case 0:
			i.Category |= Feather
		case 1:
			i.Category |= Leather
		case 2:
			i.Category |= Bone
		}

	case 3: // wool and wool products
		i.Category |= Wool

		switch rand.Intn(4) {
		case 0:
			i.Category |= Cloth
		case 1:
			i.Category |= Thread
		case 2, 3:
			// just plain wool
		}
	}

	i.Quality = uint16(rand.Intn(0x0300) + 0x0200)

	return
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

	current := s.Inventory[i]
	if delta < 0 && current < uint(-delta) {
		return false
	}

	if s.Inventory == nil {
		s.Inventory = make(map[Item]uint)
	}

	if delta < 0 && current == uint(-delta) {
		delete(s.Inventory, i)
	} else {
		s.Inventory[i] = uint(int(s.Inventory[i]) + delta)
	}

	return true
}

func (s *State) GetItemCount(i Item) (count uint) {
	s.Lock()
	count = s.Inventory[i]
	s.Unlock()
	return

}
