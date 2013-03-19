package state

import (
	"math/rand"
	"sort"
	"sync"
	"time"
	"unsafe"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func itemLeq(i, j *Item) bool {
	if i == j {
		return true
	}

	if i.Category != j.Category {
		if ci, cj := i.Category.Category(), j.Category.Category(); ci != cj {
			return ci < cj
		}

		return i.Category < j.Category
	}

	if i.Name != j.Name {
		return i.Name < j.Name
	}

	if i.Quality != j.Quality {
		return i.Quality < j.Quality
	}

	return uintptr(unsafe.Pointer(i.Components)) < uintptr(unsafe.Pointer(j.Components))
}

type State struct {
	DayLength   time.Duration
	Day         uint64
	DayProgress uint8
	Cash        uint64
	Score       uint64
	Workers     []*Worker
	Inventory   []*Item
	sync.Mutex
}

func (s *State) Init() {
	s.Lock()

	s.Workers = append(s.Workers, &Worker{})
	s.Workers = append(s.Workers, &Worker{})
	s.Workers = append(s.Workers, &Worker{})
	s.Workers = append(s.Workers, &Worker{})
	s.Workers = append(s.Workers, &Worker{})
	s.Workers = append(s.Workers, &Worker{})

	for i := 0; i < 100; i++ {
		s.addItem(randomStartingItem())
	}

	s.Unlock()
}

func randomStartingItem() (i *Item) {
	i = new(Item)

	i.Category = ItemCategory(rand.Intn(int(Large-Small+9))) + Small

	switch rand.Intn(4) {
	case 0: // metal
		if rand.Intn(3) == 0 {
			i.Category |= Ingot
		} else {
			i.Category |= Ore
		}

		switch rand.Intn(4) {
		case 0:
			i.Category |= Copper
		case 1:
			i.Category |= Tin
		case 2:
			i.Category |= Lead
		case 3:
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
	s.incrementScore(amount)
	s.Unlock()
}

func (s *State) incrementScore(amount uint64) {
	s.Score += amount
}

func (s *State) GetScore() (score uint64) {
	s.Lock()
	score = s.Score
	s.Unlock()
	return
}

func (s *State) AdvanceTime() {
	s.Lock()

	s.DayProgress++

	for _, w := range s.Workers {
		w.advanceDayPercent(s)
	}

	if s.DayProgress >= 100 {
		s.DayProgress = 0
		s.Day++

		for _, w := range s.Workers {
			w.advanceDay(s.Day, s)
		}
	}

	s.Unlock()
}

func (s *State) AddItem(i *Item) {
	s.Lock()
	s.addItem(i)
	s.Unlock()
}

func (s *State) indexItem(i *Item) int {
	return sort.Search(len(s.Inventory), func(j int) bool {
		return itemLeq(i, s.Inventory[j])
	})
}

func (s *State) addItem(i *Item) {
	index := s.indexItem(i)

	if index < len(s.Inventory) {
		s.Inventory = append(s.Inventory[:index+1], s.Inventory[index:]...)
		s.Inventory[index] = i
	} else {
		s.Inventory = append(s.Inventory[:index], i)
	}
}

func (s *State) RemoveItem(i *Item) (ok bool) {
	s.Lock()
	ok = s.removeItem(i)
	s.Unlock()
	return
}

func (s *State) removeItem(i *Item) bool {
	if i == nil {
		return false
	}

	index := s.indexItem(i)

	if index < len(s.Inventory) && s.Inventory[index] == i {
		s.Inventory = append(s.Inventory[:index], s.Inventory[index+1:]...)

		return true
	}

	return false
}
