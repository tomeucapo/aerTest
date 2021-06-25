package utils

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"

	"time"

	r "gopkg.in/gorethink/gorethink.v3"
)

type Config struct {
	LogServers []string
	Database   string
	AERConfig  AER `toml:"aer"`
}

type AER struct {
	Endpoint     string
	Timeout      time.Duration
	RequestTypes []string
	RsPrintAbove int64
}

func ReadConfig(configFile string) Config {
	_, err := os.Stat(configFile)
	if err != nil {
		log.Fatal("Config file is missing: ", configFile)
	}

	var config Config
	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		log.Fatal(err)
	}

	return config
}

func GetRDBConnectArguments(conf *Config) r.ConnectOpts {
	return r.ConnectOpts{
		Addresses:  conf.LogServers,
		InitialCap: 10,
		MaxOpen:    10,
		Database:   conf.Database}
}
