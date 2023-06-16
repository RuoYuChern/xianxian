package common

import (
	"container/list"
	"sync"
	"time"
)

type Action interface {
	Call() error
}

type Aactor struct {
	que   *list.List
	mu    sync.Mutex
	quit  chan string
	signa chan string
}

var actor *Aactor
var actOnce sync.Once

func (a *Aactor) add(act Action) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.que.PushBack(act)
	a.signa <- "One"
}

func (a *Aactor) pop() Action {
	get().mu.Lock()
	defer get().mu.Unlock()
	if a.que.Len() > 0 {
		first := a.que.Front()
		a.que.Remove(first)
		return first.Value.(Action)
	} else {
		return nil
	}
}

func (a *Aactor) Close() {
	a.signa <- "Quit"
	close(a.quit)
	close(a.signa)
}

func doWork(i int, msg string) {
	act := actor.pop()
	if act == nil {
		return
	}
	err := act.Call()
	if err != nil {
		Logger.Infof("work[%d], msg:[%s] do error:%s", i, msg, err.Error())
	}
}

func work(i int) {
	workNo := i
	stop := false
	Logger.Infof("Actor:%d start ......", workNo)
	for {
		select {
		case msgq := <-get().signa:
			doWork(workNo, msgq)
		case quit := <-get().quit:
			Logger.Infof("actor %d %s", workNo, quit)
			stop = true
		case <-time.After(600 * time.Second):
			Logger.Infof("actor %d, %s", workNo, "time after")
		}
		if stop {
			break
		}
	}
	Logger.Infof("Actor:%d quited", workNo)
}

func get() *Aactor {
	actOnce.Do(func() {
		actor = &Aactor{
			que:   list.New(),
			quit:  make(chan string),
			signa: make(chan string, 100),
		}
		for i := 0; i < 4; i++ {
			go work(i)
		}
		TaddItem(actor)
	})
	return actor
}

func AddAaction(act Action) {
	get().add(act)
}
