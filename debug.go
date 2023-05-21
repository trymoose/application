package application

import (
	"context"
	"os/signal"
	"runtime/debug"
	"syscall"
)

func _PrintDebug() {
	debugCtx, debugCancel := signal.NotifyContext(context.Background(), syscall.SIGQUIT, syscall.SIGILL, syscall.SIGTRAP, syscall.SIGABRT, syscall.SIGSTKFLT /*syscall.SIGEMT,*/, syscall.SIGSYS)
	defer debugCancel()
	<-debugCtx.Done()
	debug.PrintStack()
	Exit(2)
}
