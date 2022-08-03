package main

import (
	"Rms/database"
	"Rms/server"
	"github.com/sirupsen/logrus"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	errV := godotenv.Load(".env")
	if errV != nil {
		logrus.Error("error in loading env file")
		return
	}

	err := database.ConnectAndMigrate(os.Getenv("host"), os.Getenv("port"), os.Getenv("databaseName"), os.Getenv("user"), os.Getenv("password"), database.SSLModeDisable)
	if err != nil {
		log.Fatal(err)
		return
	}

	logrus.Println("connected")
	srv := server.SetupRoutes()
	err = srv.Run(":8081")
	if err != nil {
		log.Fatal(err)
		return
	}
}
