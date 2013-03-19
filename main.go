package main

import (
	"flag"
	"github.com/nsf/termbox-go"
	"os"
	"runtime/pprof"
)

var (
	cpuprof = flag.String("cpuprof", "", "if non-empty, the filename for a pprof CPU profile")
	memprof = flag.String("memprof", "", "if non-empty, the filename for a pprof heap profile")
)

func main() {
	flag.Parse()

	if *cpuprof != "" {
		f, err := os.Create(*cpuprof)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		err = pprof.StartCPUProfile(f)
		if err != nil {
			panic(err)
		}
		defer pprof.StopCPUProfile()
	}

	if *memprof != "" {
		f, err := os.Create(*memprof)
		if err != nil {
			panic(err)
		}

		defer func() {
			pprof.Lookup("heap").WriteTo(f, 0)
			f.Close()
		}()
	}

	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	var shouldExit bool

	go painter()

	for {
		e := termbox.PollEvent()
		switch e.Type {
		case termbox.EventError:
			panic(e.Err) // TODO: better error handling?

		case termbox.EventResize:
			repaint()

		case termbox.EventKey:
			if e.Key != termbox.KeyEsc {
				shouldExit = false
			}

			if e.Key == 0 {
				paintLock.Lock()
				paintCtx.char(e.Ch)
				paintLock.Unlock()
			} else if e.Key == termbox.KeyEsc {
				if shouldExit {
					return
				} else {
					shouldExit = true
				}
			} else {
				paintLock.Lock()
				paintCtx.key(e.Key)
				paintLock.Unlock()
			}
		}
	}
}
