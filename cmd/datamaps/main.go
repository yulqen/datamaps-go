package main

import (
	"log"
	"net/http"
	"os"

	"git.yulqen.org/go/datamaps-go/internal/datamaps"
	"git.yulqen.org/go/datamaps-go/internal/models"
)

type application struct {
	datamaps *models.DatamapModel
}

func main() {
	opts := datamaps.ParseOptions()
	if opts.Command == "help" {
		os.Stdout.WriteString(datamaps.Usage)
		os.Exit(0)
	}
	// TODO - removed this to handle "setup" bug below.
	// Check that removing this has no consequences.
	// dbpc := datamaps.NewDBPathChecker(os.UserConfigDir)
	// if !dbpc.Check() {
	// 	datamaps.SetUp()
	// }
	switch opts.Command {
	case "checkdb":
		dbpc := datamaps.NewDBPathChecker(os.UserConfigDir)
		if !dbpc.Check() {
			log.Println("No database file exists. Please run datamaps setup")
		}
	case "import":
		if err := datamaps.ImportToDB(opts); err != nil {
			log.Fatal(err)
		}
	case "datamap":
		if err := datamaps.DatamapToDB(opts); err != nil {
			log.Fatal(err)
		}
	case "setup":
		// BUG This gets called twice if the !dbpc.Check()
		// call above reveals that the config dir is present
		_, err := datamaps.SetUp()
		if err != nil {
			log.Fatal(err)
		}
	case "createmaster":
		if err := datamaps.CreateMaster(opts); err != nil {
			log.Fatal(err)
		}
	case "server":
		// Database connection
		db, err := datamaps.OpenDB("postgres://postgres:example@localhost:5432/datamaps")
		if err != nil {
			log.Fatalf("cannot connect to database - %v", err)
		}

		defer db.Close()
		app := &application{
			datamaps: &models.DatamapModel{DB: db},
		}
		log.Println("Starting server on :8080")
		log.Fatal(http.ListenAndServe(":8080", app.routes()))
	}
}
