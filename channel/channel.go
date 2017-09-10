package main

import (
	"fmt"
	"sync"
)

func isClosed(ch <-chan int) bool {
	select {
	case <-ch:
		return true
	default:
		return false
	}
}

// not suggest
func safeSend(ch chan<- int, value int) (closed bool) {
	defer func() {
		if recover() != nil {
			closed = true
		}
	}()

	ch <- value
	return false
}

// not suggest
func safeClose(ch chan int) (justClosed bool) {
	defer func() {
		if recover() != nil {
			justClosed = false
		}
	}()

	close(ch)
	return true
}

// use sync.once to force close channel once
type MyChannel struct {
	C    chan int
	once sync.Once
}

func NewMyChannel() *MyChannel {
	return &MyChannel{C: make(chan int)}
}

func (mc *MyChannel) SafeClose() {
	mc.once.Do(func() {
		close(mc.C)
	})
}

// use sync.mutext to force close channel once
type MyChannel2 struct {
	C      chan int
	closed bool
	mutex  sync.Mutex
}

func NewMyChannel2() *MyChannel2 {
	return &MyChannel2{C: make(chan int)}
}

func (mc *MyChannel2) SafeClose() {
	mc.mutex.Lock()
	if !mc.closed {
		close(mc.C)
		mc.closed = true
	}
	mc.mutex.Unlock()
}

func (mc *MyChannel2) IsClosed() bool {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	return mc.closed
}

func main() {
	c := make(chan int)
	fmt.Println(isClosed(c))
	for i := 0; i < 2; i++ {
		ret := safeClose(c)
		fmt.Println("close result:", ret)
	}
	fmt.Println(isClosed(c))

	mc := NewMyChannel()
	for i := 0; i < 2; i++ {
		mc.SafeClose()
	}

	mc2 := NewMyChannel2()
	fmt.Println("mc status:", mc2.IsClosed())
	for i := 0; i < 2; i++ {
		mc2.SafeClose()
	}
	fmt.Println("mc status:", mc2.IsClosed())
}
