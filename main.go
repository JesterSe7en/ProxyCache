package main

import (
	"fmt"
	"os"
)

func main() {
	// parse cmd args
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Usage: cacheProxy <port> <redirectURL>")
		os.Exit(1)
	}

	redisConfig, err := loadConfig()
	if err != nil || redisConfig == nil {
		LogFatal("cannot load redis config", err)
	}

	// setup logging and handle basic configuration
	err = initLogger()
	if err != nil {
		LogError("cannot initialize logger", err)
	}

	// start the server and connect to other module (Redis)
}
