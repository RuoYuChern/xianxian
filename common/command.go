package common

import (
	"container/list"
	"sync"
)

type TAutoCloseable interface {
	Close()
}

var items *list.List = list.New()
var mu sync.Mutex = sync.Mutex{}

func TaddItem(c TAutoCloseable) {
	mu.Lock()
	defer mu.Unlock()
	items.PushBack(c)
}

func TcloseItems() {
	mu.Lock()
	defer mu.Unlock()
	for c := items.Front(); c != nil; c = c.Next() {
		if c != nil {
			ac := c.Value.(TAutoCloseable)
			ac.Close()
		}
	}
}
