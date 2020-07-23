/*
datamaps-go is a simple tool to extract from and send data to spreadsheets.
*/
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/yulqen/datamaps-go/pkg/datamaps"
)

const (
	configDirName = "datamaps-go"
	dbName        = "datamaps.db"
)

func setUp() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	// check if config folder exists
	configPath := filepath.Join(dir, configDirName)
	dbPath := filepath.Join(configPath, dbName)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Println("Config directory does not exist.")
		log.Printf("Creating config directory %s\n", configPath)
		if err := os.Mkdir(filepath.Join(dir, "datamaps-go"), 0700); err != nil {
			return "", err
		}
	} else {
		log.Println("Config directory found.")
	}
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Println("Database does not exist.")
		_, err = os.Create(dbPath)
		if err != nil {
			return "", err
		}
		log.Printf("Creating database file at %s\n", dbPath)
		_, err := datamaps.SetupDB(dbPath)
		if err != nil {
			return "", err
		}
	} else {
		log.Println("Database file found.")
	}
	return dir, nil
}

func main() {

	opts := datamaps.ParseOptions()

	switch opts.Command {
	case "datamap":
		err := datamaps.DatamapToDB(opts)
		if err != nil {
			log.Fatal(err)
		}
	case "setup":
		_, err := setUp()
		if err != nil {
			log.Fatal(err)
		}
	case "server":
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				http.NotFound(w, r)
				return
			}
			fmt.Fprintf(w, "Welcome to datamaps!")
			// or you could write it thus
			// w.Write([]byte("Hello from datamaps"))
		})
		log.Println("Starting server on :8080")
		err := http.ListenAndServe(":8080", nil)
		log.Fatal(err)
	}
}
