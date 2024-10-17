package main

import (
	"flag"
	"os"
)

type ProxyConfig struct {
	port        int
	redirectURL string
}

func main() {
	// parse cmd args
	var port = flag.Int("port", -1, "Required. The port on which the server will listen for incoming requests")
	var redirectURL = flag.String("redirectURL", "", "Required. The URL of the external service to be proxied")

	flag.Parse()

	if *port == -1 || *redirectURL == "" {
		flag.Usage()
		os.Exit(1)
	}

	var config ProxyConfig
	config.port = *port
	config.redirectURL = *redirectURL

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
