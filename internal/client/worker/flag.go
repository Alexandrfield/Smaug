package worker

import (
	"flag"
)

func ParseFlags() (Config, error) {
	var config Config
	flag.StringVar(&config.ServerAddr, "s", "127.0.0.1:50051", "address and port to run server ")
	flag.Parse()
	return config, nil
}
