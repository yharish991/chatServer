package testhelpers

import (
	"log"
	"os"
	"strings"
)

// GetServerRootDir is a utility function that gets the project root directory
func GetServerRootDir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	srcIndex := strings.LastIndex(dir, "/src")
	dir = dir[:srcIndex]
	return dir
}
