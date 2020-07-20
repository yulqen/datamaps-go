/*
datamaps-go is a simple tool to extract from and send data to spreadsheets.
*/
package main

import (
	"flag"
	"fmt"
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

// from datamaps-c

// { 0,0,0,0, "Datamap options: (when calling 'datamaps datamap')." },
// {"import", DM_IMPORT, "PATH", 0, "PATH to datamap file to import."},
// {"name", DM_NAME, "NAME", 0, "The name you want to give to the imported datamap."},
// {"overwrite", DM_OVERWRITE, 0, 0, "Start fresh with this datamap (erases existing datamap data)."},
// {"initial", DM_INITIAL, 0, 0, "This option must be used where no datamap table yet exists."},

// Checkout https://gobyexample.com/command-line-subcommands

func main() {
	// setup command
	setupCmd := flag.NewFlagSet("setup", flag.ExitOnError)
	setupCmd.Usage = func() { fmt.Println("No, you fucking idiot!") }

	// datamap command and its flags
	datamapCmd := flag.NewFlagSet("datamap", flag.ExitOnError)
	importFlg := datamapCmd.String("import", "/home/lemon/Documents/datamaps/input/datamap.csv", "Path to datamap")
	nameFlg := datamapCmd.String("name", "Unnamed datamap", "The name you want to give to the imported datamap.")
	overwriteFlg := datamapCmd.Bool("overwrite", false, "Start fresh with this datamap (erases existing datamap data).")
	initialFlg := datamapCmd.Bool("initial", false, "This option must be used where no datamap table yet exists.")

	if len(os.Args) < 2 {
		fmt.Println("expected 'datamap' or 'setup' subcommand")
		os.Exit(1)
	}

	switch os.Args[1] {

	case "datamap":
		datamapCmd.Parse(os.Args[2:])
		fmt.Println("subcommand 'datamap'")
		fmt.Println("  import:", *importFlg)
		fmt.Println("  name:", *nameFlg)
		fmt.Println("  overwrite:", *overwriteFlg)
		fmt.Println("  initial:", *initialFlg)
	case "setup":
		setupCmd.Parse(os.Args[2:])
		_, err := setUp()
		if err != nil {
			log.Fatal(err)
		}

	default:
		fmt.Println("Do not recognised that command. Expected 'datamap' subcommand.")
		os.Exit(1)
	}
}

// var setup bool

// // Example for go docs about how to set up short and long flags.
// const setupUsage = "Initialise configuration and database files"

// flag.BoolVar(&setup, "setup", false, setupUsage)
// flag.BoolVar(&setup, "s", false, setupUsage)
// flag.Parse()
// if setup == true {
// 	_, err := setUp()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// } else {
// 	fmt.Println("No command given.")
// 	flag.PrintDefaults()
// }

// func main() {
// 	app := &cli.App{
// 		Flags: []cli.Flag{
// 			&cli.StringFlag{
// 				Name:    "datamap, dm",
// 				Aliases: []string{"dm"},
// 				Value:   "/home/lemon/Documents/datamaps/input/datamap.csv",
// 				Usage:   "Path to a datamap file",
// 			},
// 			&cli.StringFlag{
// 				Name:    "master, m",
// 				Aliases: []string{"m"},
// 				Value:   "/home/lemon/Documents/datamaps/input/master.xlsx",
// 				Usage:   "Path to a master file",
// 			},
// 		},
// 		Commands: []*cli.Command{
// 			{
// 				Name:    "import",
// 				Aliases: []string{"i"},
// 				Usage:   "Import a bunch of populated templates",
// 				Action: func(c *cli.Context) error {
// 					return nil
// 				},
// 			},
// 			{
// 				Name:    "export",
// 				Aliases: []string{"e"},
// 				Usage:   "Export a master to populate blank templates",
// 				Action: func(c *cli.Context) error {
// 					return nil
// 				},
// 			},
// 		},
// 		Name:  "datamaps",
// 		Usage: "Import and export data to and from spreadsheets",
// 		Action: func(c *cli.Context) error {
// 			fmt.Println("DATAMAPS")
// 			return nil
// 		},
// 	}

// 	err := app.Run(os.Args)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	_, err = setUp()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
