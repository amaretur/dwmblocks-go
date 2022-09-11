package main 

import (
	"os"
	"log"
)

func getHomeDir() string {

	homeDir, ok := os.LookupEnv("HOME")
	if ! ok {   
		log.Fatal("could not get the value of the environment variable $HOME")
	} 

	return homeDir
}

func contains(array []string, e string) bool {

	for _, element := range array {
		if element == e {
			return true
		}
	}

	return false
}

