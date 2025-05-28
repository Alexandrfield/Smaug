package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Alexandrfield/Smaug/internal/common"
	"github.com/Alexandrfield/Smaug/internal/server/worker"
)

func main() {
	fmt.Printf("Start Smaug server\n")
	logger := common.GetComponentLogger("debug")
	done := make(chan struct{})
	osSignals := make(chan os.Signal, 1)
	ctx, cancle := context.WithCancel(context.Background())
	config, err := worker.ParseFlags()
	if err != nil {
		logger.Fatalf("Can't start application. err:%s", err)
	}
	go worker.ServerLoop(ctx, done, config, logger)

	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	<-osSignals
	cancle()
	<-done
	fmt.Printf("Stop Smaug server\n")
}
