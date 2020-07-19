/*
datamaps-go is a simple tool to extract from and send data to spreadsheets.
*/
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
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
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "datamap, dm",
				Aliases: []string{"dm"},
				Value:   "/home/lemon/Documents/datamaps/input/datamap.csv",
				Usage:   "Path to a datamap file",
			},
			&cli.StringFlag{
				Name:    "master, m",
				Aliases: []string{"m"},
				Value:   "/home/lemon/Documents/datamaps/input/master.xlsx",
				Usage:   "Path to a master file",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "import",
				Aliases: []string{"i"},
				Usage:   "Import a bunch of populated templates",
				Action: func(c *cli.Context) error {
					return nil
				},
			},
			{
				Name:    "export",
				Aliases: []string{"e"},
				Usage:   "Export a master to populate blank templates",
				Action: func(c *cli.Context) error {
					return nil
				},
			},
		},
		Name:  "datamaps",
		Usage: "Import and export data to and from spreadsheets",
		Action: func(c *cli.Context) error {
			fmt.Println("DATAMAPS")
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	_, err = setUp()
	if err != nil {
		log.Fatal(err)
	}
}
