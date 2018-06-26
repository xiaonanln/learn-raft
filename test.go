package main

import (
	"runtime"
	"sync"
)

var num = 0

func main() {
	runtime.GOMAXPROCS(1)
	var wait sync.WaitGroup
	for i := 0; i < 100000; i++ {
		wait.Add(1)
		go jj(&wait)
	}

	wait.Wait()
	println(num)
}

func jj(wait *sync.WaitGroup) {
	num += 1
	wait.Done()
}
