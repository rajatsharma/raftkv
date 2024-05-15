package config

import (
	"log"
	"os"
)

type Config struct {
	Id       string
	HttpPort string
	RaftPort string
}

func GetConfig() Config {
	cfg := Config{}
	for i, arg := range os.Args[1:] {
		if arg == "--node-id" {
			cfg.Id = os.Args[i+2]
			i++
			continue
		}

		if arg == "--httpserver-port" {
			cfg.HttpPort = os.Args[i+2]
			i++
			continue
		}

		if arg == "--node-port" {
			cfg.RaftPort = os.Args[i+2]
			i++
			continue
		}
	}

	if cfg.Id == "" {
		log.Fatal("Missing required parameter: --node-id")
	}

	if cfg.RaftPort == "" {
		log.Fatal("Missing required parameter: --node-port")
	}

	if cfg.HttpPort == "" {
		log.Fatal("Missing required parameter: --httpserver-port")
	}

	return cfg
}
