package main

import (
	"log"
	"os"
	"path"
	"strings"

	"chatServer/src/api"
	"chatServer/src/chatserver"
	"chatServer/src/config"
	"chatServer/src/connections"
)

// function that returns the path of chat server root
func getServerRootDir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	srcIndex := strings.LastIndex(dir, "/src")
	dir = dir[:srcIndex]
	return dir
}

func main() {
	log.Println("Starting the chat server!!!")

	// read the config
	reader := config.NewReaderImpl()
	cfg := reader.Read(path.Join(getServerRootDir(), "/resources/config/config.json"))

	// start the chat server
	chatService := chatserver.NewServiceImpl(path.Join(getServerRootDir(), cfg.LogFilePath))
	chatService.Run()

	// start the api server
	apiService := api.NewServiceImpl(chatService)
	apiController := api.NewControllerImpl(apiService)
	go apiController.Register()

	// handle incoming connections
	connectionsService := connections.NewServiceImpl(chatService, cfg)
	connectionsService.HandleConnections()
}
