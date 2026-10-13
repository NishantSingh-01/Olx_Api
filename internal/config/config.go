package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Env         string
	DatabaseURL string
}


func MustLoad() Config {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT is Required") // for better logs and error checking we are using panic here.
	}
	env := os.Getenv("ENV")
	if env == "" {
		panic("ENV is Required")
	}
	db_url := os.Getenv("DATABASE_URL")
	if db_url == "" {
		panic("DATABASE_URL is Required")
	}
	return Config{
		Port:        port,
		Env:         env,
		DatabaseURL: db_url,
	}
}
