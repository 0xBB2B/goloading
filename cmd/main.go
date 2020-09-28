package main

import (
	"fmt"
	"goloading/pkg/progress"
	"math/rand"
	"time"
)

func main() {
	pm := progress.NewProgressManager()
	defer fmt.Println("Done!")
	defer pm.Wait()

	go func() {
		p := pm.Add()
		defer close(p)
		for i := 0; i <= 100; i++ {
			p <- fmt.Sprintf("test1: %d/%d", i, 100)
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		}
		p <- "test1: done test1test1test1test1test1test1test1test1test1test1test1test1test1test1"
	}()
	go func() {
		p := pm.Add()
		defer close(p)
		for i := 0; i <= 200; i++ {
			p <- fmt.Sprintf("test2: %d/%d test2test2test2test2test2test2test2test2test2test2test2test2test2", i, 200)
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		}
		p <- "test2: done"
	}()
	time.Sleep(time.Second)
	go func() {
		p := pm.Add()
		defer close(p)
		for i := 0; i <= 150; i++ {
			p <- fmt.Sprintf("test3: %d/%d", i, 150)
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		}
		p <- "test3: done"
	}()
}
