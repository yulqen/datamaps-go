/*
datamaps-go is a simple tool to extract from and send data to spreadsheets.
*/
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func createConfigDir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	// check if config folder exists
	config_path := filepath.Join(dir, "datamaps-go")
	if _, err := os.Stat(config_path); os.IsNotExist(err) {
		log.Println("Config directory does not exist.")
		log.Println("Creating config directory.")
		if err := os.Mkdir(filepath.Join(dir, "datamaps-go"), 0700); err != nil {
			return "", err
		}
	} else {
		log.Println("Config directory found.")
	}
	return dir, nil
}

// Entry point

func main() {
	dir, err := createConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
}
