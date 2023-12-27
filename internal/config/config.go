package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

type cliOptions struct {
	Resolver string
	Config   string
}

type fileOptions struct {
	Records []Record `json:"records"`
}

type Record struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
	TTL   int    `json:"ttl"`
	Note  string `json:"note,omitempty"`
}

type Config struct {
	cliOptions
	fileOptions
}

var config Config

func Load() error {
	flag.StringVar(&config.Resolver, "resolver", "", "resolver address to forward queries to (ip:port)")
	flag.StringVar(&config.Config, "config", "", "config filepath")
	flag.Parse()

	// Config file is optional
	if config.Config == "" {
		return nil
	}

	file, err := os.Open(config.Config)
	if err != nil {
		return fmt.Errorf("failed to open config file: %s", err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Println("Failed to close config file:", err)
		}
	}()

	decoder := json.NewDecoder(file)

	var config Config
	err = decoder.Decode(&config)
	if err != nil {
		return fmt.Errorf("failed to decode config file: %s", err)
	}

	return nil
}

func Get() Config {
	return config
}
