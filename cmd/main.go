/*
datamaps-go is a simple tool to extract from and send data to spreadsheets.
*/
package main

import (
	"log"
	"os"
	"path/filepath"
)

func setUp() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	// check if config folder exists
	config_path := filepath.Join(dir, "datamaps-go")
	db_path := filepath.Join(config_path, "datamaps.db")
	if _, err := os.Stat(config_path); os.IsNotExist(err) {
		log.Println("Config directory does not exist.")
		log.Printf("Creating config directory %s\n", config_path)
		if err := os.Mkdir(filepath.Join(dir, "datamaps-go"), 0700); err != nil {
			return "", err
		}
	} else {
		log.Println("Config directory found.")
	}
	if _, err := os.Stat(db_path); os.IsNotExist(err) {
		log.Println("Database does not exist.")
		_, err = os.Create(db_path)
		if err != nil {
			return "", err
		}
		log.Printf("Creating database file at %s\n", db_path)
	} else {
		log.Println("Database file found.")
	}
	return dir, nil
}

// Entry point

func main() {
	dir, err := setUp()
	if err != nil {
		log.Fatal(err)
	}
}
