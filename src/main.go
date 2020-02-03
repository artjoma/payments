package main

import (
	"fmt"
	"os"
	"os/signal"
	"payments/src/api"
	"payments/src/repo"
	"payments/src/service"
	"runtime"
	"syscall"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func init() {
	if os.Getenv("ENV_NAME") == "" {
		//for local development
		os.Setenv("DB_HOST", "127.0.0.1")
		//for development
		fmt.Println("not found ENV_NAME key try expose env var from .env file")
		if err := godotenv.Load(); err != nil {
			if err := godotenv.Load("../.env"); err != nil {
				panic(".env file not found")
			}
		}
	}
	//init logging
	envName := os.Getenv("ENV_NAME")
	if envName == "local" {
		log.SetFormatter(&log.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05.00000",
			FullTimestamp:   true,
		})
	} else {
		log.SetFormatter(&log.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05.00000",
		})
	}
	log.Infof("ENV: %s", envName)
}

func main() {
	sigs := make(chan os.Signal, 1)
	if runtime.GOOS != "windows" {
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	}

	createAppContext()
	sig := <-sigs
	fmt.Println("shutdown. Signal: ", sig)
	close(sigs)

	destroy()
}

// layers initialization ordering important !
func createAppContext() {
	//init data access layer
	repo.InitRepo(
		os.Getenv("APP_CONNECTION_POOL_SIZE"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))
	//init services layer
	service.InitService()
	//init HTTP API layer
	api.InitHttp()
}

//ordering important
func destroy() {
	api.ShutdownHttp()
	service.ShutdownService()
	repo.ShutdownRepo()
}
