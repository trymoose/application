package application

import (
	"os"
	"sync"
	"time"
)

var exit struct {
	once   sync.Once
	cancel func()
}

func Exit(code int) {
	exit.once.Do(func() {
		exit.cancel()
		time.Sleep(time.Second * 5)
		os.Exit(code)
	})
}

func ExitSleepForever() {
	exit.once.Do(exit.cancel)
	select {}
}

func _SetExitCancel(cancel func()) {
	if exit.cancel == nil {
		exit.cancel = cancel
	}
	old := exit.cancel
	exit.cancel = func() {
		defer old()
		defer cancel()
	}
}
