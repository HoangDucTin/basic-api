package main

import (
	"basic-api/internal/logger"
	"basic-api/internal/mongo"
	"basic-api/internal/server"
	"time"
)

func connectMongo() {
	address := "127.0.0.1:27017"
	database := "test"
	username := ""
	password := ""
	timeout := 5 * time.Second

	logger.Info("Initialize MongoDB connection (%s).", address)

	if err := mongo.Setup(address, database, username, password, timeout);
		err != nil {
		logger.Exit(err.Error())
	}

	logger.Info("Successfully connected to MongoDB.")
}

func releaseMongo() {
	logger.Info("Release MongoDB connection.")
	mongo.Close()
	logger.Info("Successfully release MongoDB connections.")
}

func startService() {
	address := ":8080"
	timeout := 5 * time.Second

	logger.Info("Start serving at '%s'.", address)

	if err := server.ServeHttp(address, timeout);
		err != nil {
		logger.Fail(err.Error())
	}

	logger.Info("Stop serving at '%s'.", address)
}

func main() {
	connectMongo()
	startService()
	releaseMongo()
}

// End-of-file
