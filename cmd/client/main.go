package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Alexandrfield/Smaug/internal/client/worker"
	"github.com/Alexandrfield/Smaug/internal/common"
)

func main() {
	fmt.Printf("Start Smaug aplication\n")
	logger := common.GetComponentLogger()
	done := make(chan struct{})

	go worker.ClientAppLoop(done, logger)

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	select {
	case <-done:
	case <-osSignals:
	}
	fmt.Printf("Stop Smaug aplication\n")
}
