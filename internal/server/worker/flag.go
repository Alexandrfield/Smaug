package worker

import (
	"flag"
	"log"
)

func ParseFlags() (Config, error) {

	var config Config
	var key string
	flag.StringVar(&key, "k", "", "master key for encrypt data")
	var databaseURI string
	flag.StringVar(&databaseURI, "d", "", "uri for database [default:]")
	flag.StringVar(&config.Port, "s", "50051", "port to run server")
	flag.Parse()

	if len(key) == 0 {
		log.Fatalf("No key for start Server\n")
	}
	if len(key) == 0 {
		log.Fatalf("No database dsn for start Server\n")
	}
	config.DatabasDSN = databaseURI
	config.MasterKey = []byte(key)
	return config, nil
}
