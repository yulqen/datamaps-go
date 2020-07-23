/*
datamaps-go is a simple tool to extract from and send data to spreadsheets.
*/
package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/yulqen/datamaps-go/pkg/datamaps"
)

const (
	config_dir_name = "datamaps-go"
	db_name         = "datamaps.db"
)

func setUp() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	// check if config folder exists
	config_path := filepath.Join(dir, config_dir_name)
	db_path := filepath.Join(config_path, db_name)
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
		_, err := datamaps.SetupDB(db_path)
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
		data, err := datamaps.ReadDML(opts.DMPath)
		if err != nil {
			log.Fatal(err)
		}
		opts.DMData = data

		err = datamaps.DatamapToDB(opts)
		if err != nil {
			log.Fatal(err)
		}
	case "setup":
		_, err := setUp()
		if err != nil {
			log.Fatal(err)
		}
	}

	// switch os.Args[1] {

	// case "server":
	// 	opts.Command.Parse(os.Args[2:])
	// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 		if r.URL.Path != "/" {
	// 			http.NotFound(w, r)
	// 			return
	// 		}
	// 		fmt.Fprintf(w, "Welcome to datamaps!")
	// 		// or you could write it thus
	// 		// w.Write([]byte("Hello from Snippetbox"))
	// 	})
	// 	log.Println("Starting server on :8080")
	// 	err := http.ListenAndServe(":8080", nil)
	// 	log.Fatal(err)

	// case "datamap":
	// 	datamapCmd.Parse(os.Args[2:])
	// 	fmt.Println("subcommand 'datamap'")
	// 	fmt.Println("  import:", *dmPath)
	// 	fmt.Println("  name:", *dmName)
	// 	fmt.Println("  overwrite:", *dmOverwrite)
	// 	fmt.Println("  initial:", *dmInitial)

	// 	dir, err := os.UserConfigDir()
	// 	if err != nil {
	// 		os.Exit(1)
	// 	}
	// 	// check if config folder exists
	// 	config_path := filepath.Join(dir, config_dir_name)
	// 	if _, err := os.Stat(config_path); os.IsNotExist(err) {
	// 		fmt.Println("Config directory and database does not exist. Run datamaps setup to fix.")
	// 		os.Exit(1)
	// 	}
	// 	// Here we actually read the data from the file
	// 	data, err := datamaps.ReadDML(opts.ImportPath)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	opts.DMData = data

	// 	opts.DBPath = filepath.Join(config_path, db_name)
	// 	err = datamaps.DatamapToDB(opts)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// case "setup":
	// 	setupCmd.Parse(os.Args[2:])
	// 	_, err := setUp()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// default:
	// 	fmt.Println("Do not recognised that command. Expected 'datamap' subcommand.")
	// 	os.Exit(1)
	// }
}
