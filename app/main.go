package main

import (
	"log"
	_config "user-service/config"
	_server "user-service/server"
)

func main() {
	config := _config.LoadConfig(".")

	db := _config.InitDB(config)

	cacheClient := _config.NewRedisClient(config)
	defer cacheClient.Close()
	server := _server.NewServer(&config, db, cacheClient)

	err := server.Run()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
}
