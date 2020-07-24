package datamaps

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	configDirName = "datamaps-go"
	dbName        = "datamaps.db"
)

// mocking funcs in go https://stackoverflow.com/questions/19167970/mock-functions-in-go
// we only need the func signature to create the type. This is pretty weird, we're want to mock
// os.UserHomeDir so that we can set it to something like /tmp in our tests. Here we are creating
// two types: GetUserConfigDir to represent the os.UserConfigDir function, DBPathChecker as a wrapper
// which which we can assign methods to that holds the value of the func os.UserConfigDir and the
// method, check(), which does the work, using the passed in func to determine the user $HOME/.config
// directory.
// Which is a lot of work for what it is, but it does make this testable and serves as an example
// of how things could be done in Go.

type GetUserConfigDir func() (string, error)

type DBPathChecker struct {
	getUserConfigDir GetUserConfigDir
}

func NewDBPathChecker(h GetUserConfigDir) *DBPathChecker {
	return &DBPathChecker{getUserConfigDir: h}
}

func (db *DBPathChecker) check() bool {
	userConfig, err := db.getUserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	configPath := filepath.Join(userConfig, "datamaps")
	dbPath := filepath.Join(configPath, "datamaps.db")
	fmt.Fprintf(os.Stderr, "dbPath is definitely %s\n", dbPath)
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "db does not exist\n")
		return false
	}
	return true
}

// DetectConfig looks for the configuration directory and
// files, and the database file needed to run the application.
func DetectDBFile() bool {
	dir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	// check if config folder exists
	configPath := filepath.Join(dir, configDirName)
	dbPath := filepath.Join(configPath, dbName)
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "db does not exist\n")
		return false
	}
	return true
}

// SetUp creates the config directory and requisite files
func SetUp() (string, error) {
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
		_, err := SetupDB(dbPath)
		if err != nil {
			return "", err
		}
	} else {
		log.Println("Database file found.")
	}
	return dir, nil
}

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
