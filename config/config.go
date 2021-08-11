package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Debug bool
}

var config Config

func InitConfig() {
	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))

	if err != nil {
		log.Fatal(err)
	}

	config = Config{
		Debug: debug,
	}
}

func GetConfig() Config {
	return config
}
