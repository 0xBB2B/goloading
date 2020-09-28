package progress

import (
	"fmt"
	"strings"
	"sync"
	"unicode"

	"github.com/nsf/termbox-go"
)

var (
	wg sync.WaitGroup
	mu sync.Mutex
)

// Progress is print message channel
type Progress chan<- string

// ManagerProgress is progress view master
type ManagerProgress interface {
	Add() Progress
	Wait()
}

// NewProgressManager is init new ProgressManager
func NewProgressManager() ManagerProgress {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()
	w, h := termbox.Size()
	lm := &loadManager{width: w, height: h}
	return lm
}

type loadManager struct {
	maxID  int
	width  int
	height int
}

func (lm *loadManager) Add() Progress {
	mu.Lock()
	defer mu.Unlock()

	wg.Add(1)
	lm.maxID++
	l := &load{no: lm.maxID, mother: lm}
	l.c = make(chan string, 1)
	go func(l *load) {
		defer wg.Done()
		for s := range l.c {
			l.print(s)
		}
	}(l)
	fmt.Printf("\n")
	fmt.Printf("\033[s")
	return l.c
}

func (lm loadManager) Wait() {
	wg.Wait()
}

type load struct {
	c      chan string
	no     int
	mother *loadManager
}

func (l *load) print(s string) {
	mu.Lock()
	defer mu.Unlock()

	if !isASCII(s) {
		panic(fmt.Errorf("Only ASCII is supported"))
	}
	str := strings.Replace(s, "\n", " ", -1)
	w, h := l.mother.width, l.mother.height
	if len(str) > w {
		str = fmt.Sprintf("%-*.*s", w, w, str)
	}
	m := l.mother.maxID - l.no + 1
	if m >= h {
		return
	}
	fmt.Printf("\033[%dA", m)
	var realArgs []interface{}
	realArgs = append(realArgs, "\r")
	realArgs = append(realArgs, fmt.Sprintf("%s", str))
	fmt.Print(realArgs...)
	fmt.Printf("\033[K")
	fmt.Printf("\033[u")
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}
