package datamaps

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	config_dir_name = "datamaps-go"
	db_name         = "datamaps.db"
)

func getUserConfigDir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	config_path := filepath.Join(dir, config_dir_name)
	return config_path, nil
}

func defaultDMPath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "Documents", "datamaps"), nil
}

// TODO - need a func to replace Options.Command with the one we pass
// Needs to use flag.NewFlagSet so we can Parse on it in main.

type Options struct {
	Command     string
	DBPath      string
	DMPath      string
	DMName      string
	DMOverwrite bool
	DMInitial   bool
	DMData      []DatamapLine
}

func defaultOptions() *Options {
	dbpath, err := getUserConfigDir()
	if err != nil {
		log.Fatalf("Unable to get user config directory %v", err)
	}
	dmpath, err := defaultDMPath()
	if err != nil {
		log.Fatalf("Unable to get default datamaps directory %v", err)
	}
	return &Options{
		Command:     "help",
		DBPath:      filepath.Join(dbpath, "datamaps.db"),
		DMPath:      dmpath,
		DMName:      "Unnamed Datamap",
		DMOverwrite: false,
		DMInitial:   false,
		DMData:      make([]DatamapLine, 0),
	}
}

func nextString(args []string, i *int, message string) string {
	if len(args) > *i+1 {
		*i++
	} else {
		log.Fatal(message)
	}
	return args[*i]
}

func processOptions(opts *Options, allArgs []string) {
	switch allArgs[0] {
	case "server":
		opts.Command = "server"
	case "datamap":
		opts.Command = "datamap"
	case "setup":
		opts.Command = "setup"
	default:
		log.Fatal("No relevant command provided.")
	}
	restArgs := allArgs[1:]
	for i := 0; i < len(allArgs[1:]); i++ {
		arg := restArgs[i]
		switch arg {
		case "-i", "--import":
			opts.DMPath = nextString(restArgs, &i, "import path required")
		case "-n", "--name":
			opts.DMName = nextString(restArgs, &i, "datamap name required")
		}

	}
}

func ParseOptions() *Options {
	opts := defaultOptions()
	processOptions(opts, os.Args[1:])

	// setup command
	setupCmd := flag.NewFlagSet("setup", flag.ExitOnError)
	setupCmd.Usage = func() { fmt.Println("No, you fucking idiot!") }

	// datamap command and its flags
	datamapCmd := flag.NewFlagSet("datamap", flag.ExitOnError)
	_ = datamapCmd.String("import", opts.DMPath, "Path to datamap")
	_ = datamapCmd.String("name", opts.DMName, "The name you want to give to the imported datamap.")
	_ = datamapCmd.Bool("overwrite", opts.DMOverwrite, "Start fresh with this datamap (erases existing datamap data).")
	_ = datamapCmd.Bool("initial", opts.DMInitial, "This option must be used where no datamap table yet exists.")

	// server command and its flags
	_ = flag.NewFlagSet("server", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("Expected 'datamap' or 'setup' subcommand")
		os.Exit(1)
	}

	return opts

}
