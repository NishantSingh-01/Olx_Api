package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Port string
	Env  string
}

func MustLoad() Config {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT is Required") // for better logs and error checking we are using panic here.
	}
	env := os.Getenv("ENV")
	if env == "" {
		panic("ENV is Required")
	}
	return Config{
		Port: port,
		Env:  env,
	}
}
