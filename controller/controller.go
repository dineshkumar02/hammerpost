package controller

import (
	"os"
	"os/signal"
	"syscall"
)

type fn func()

func HandleCntrlC(f fn) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGABRT, syscall.SIGSEGV)
	go func() {
		<-c
		f()
		os.Exit(1)
	}()
}
