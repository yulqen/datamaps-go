package datamaps

import (
	"log"
	"os"
	"path/filepath"
)

const (
	configDirName = "datamaps-go"
	dbName        = "datamaps.db"
)

func getUserConfigDir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	configPath := filepath.Join(dir, configDirName)

	return configPath, nil
}

func defaultDMPath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "Documents", "datamaps"), nil
}

// Options for the whole CLI application.
type Options struct {
	Command     string
	DBPath      string
	DMPath      string
	DMName      string
	DMOverwrite bool
	DMInitial   bool
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
		DBPath:      filepath.Join(dbpath, dbName),
		DMPath:      dmpath,
		DMName:      "Unnamed Datamap",
		DMOverwrite: false,
		DMInitial:   false,
	}
}

// nextString get the next string in a slice.
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
		case "--overwrite":
			opts.DMOverwrite = true
		case "--initial":
			opts.DMInitial = true
		}
	}
}

//ParseOptions for CLI.
func ParseOptions() *Options {
	opts := defaultOptions()
	processOptions(opts, os.Args[1:])

	return opts
}
