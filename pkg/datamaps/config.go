package datamaps

import (
	"log"
	"os"
	"path/filepath"
)

const (
	configDirName = "datamaps"
	dbName        = "datamaps.db"
)

// mocking funcs in go https://stackoverflow.com/questions/19167970/mock-functions-in-go
// we only need the func signature to create the type. This is pretty weird, we're want to mock
// os.UserHomeDir so that we can set it to something like /tmp in our tests. Here we are creating
// two types: getUserConfigDir to represent the os.UserConfigDir function, dbPathChecker as a wrapper
// which which we can assign methods to that holds the value of the func os.UserConfigDir and the
// method, check(), which does the work, using the passed in func to determine the user $HOME/.config
// directory.
// Which is a lot of work for what it is, but it does make this testable and serves as an example
// of how things could be done in Go.

// getUserConfigDir allows replaces os.UserConfigDir
// for testing purposes.
type getUserConfigDir func() (string, error)

// DBPathChecker contains the func used to create the user config dir.
type DBPathChecker struct {
	userConfig getUserConfigDir
}

// NewDBPathChecker creates a DBPathChecker using whatever
// func you want as the argument, as long as it matches the
// type os.UserConfigDir. This makes it convenient for testing
// and was done as an experiment here to practice mocking in Go.
func NewDBPathChecker(h getUserConfigDir) *DBPathChecker {
	return &DBPathChecker{userConfig: h}
}

// Check returns true if the necessary config files (including
// the database) are in place - false if not
func (db *DBPathChecker) Check() bool {
	userConfig, err := db.userConfig()
	if err != nil {
		log.Fatal(err)
	}
	dbPath := filepath.Join(userConfig, "datamaps.db")
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
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
		log.Printf("Creating config directory %s.\n", configPath)
		if err := os.Mkdir(filepath.Join(dir, "datamaps"), 0700); err != nil {
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
		log.Printf("Creating database file at %s.\n", dbPath)
		_, err := setupDB(dbPath)
		if err != nil {
			return "", err
		}
	} else {
		log.Println("Database file found.")
	}
	return dir, nil
}

func userConfigDir() (string, error) {
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

func defaultXLSXPath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "Documents", "datamaps", "import"), nil
}

// Options for the whole CLI application.
type Options struct {
	// Command is the main CLI sub-command (e.g. "datamap" handles all datamap
	// operations and the flags that follow pertain only to that operation.
	Command string

	// DBPath is the path to the database file.
	DBPath string

	// DMPath is the path to a datamap file.
	DMPath string

	// DMNane is the name of a datamap, whether setting or querying.
	DMName string

	// XLSXPath is the path to a directory containing ".xlsx" files for
	// importing.
	XLSXPath string

	// ReturnName is the name of a Return, whether setting or querying.
	ReturnName string

	// DMOverwrite is currently not used.
	DMOverwrite bool

	// DMInitial is currently not used.
	DMInitial bool

	// MasterOutPutPath is where the master.xlsx file is to be saved
	MasterOutPutPath string
}

func defaultOptions() *Options {
	dbpath, err := userConfigDir()
	if err != nil {
		log.Fatalf("Unable to get user config directory %v", err)
	}

	dmPath, err := defaultDMPath()
	xlsxPath, err := defaultXLSXPath()

	if err != nil {
		log.Fatalf("Unable to get default datamaps directory %v", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("unable to get user home directory")
	}

	return &Options{
		Command:          "help",
		DBPath:           filepath.Join(dbpath, dbName),
		DMPath:           dmPath,
		DMName:           "Unnamed Datamap",
		XLSXPath:         xlsxPath,
		ReturnName:       "Unnamed Return",
		DMOverwrite:      false,
		DMInitial:        false,
		MasterOutPutPath: filepath.Join(homeDir, "Desktop"),
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
	case "import":
		opts.Command = "import"
	case "datamap":
		opts.Command = "datamap"
	case "setup":
		opts.Command = "setup"
	case "server":
		opts.Command = "server"
	default:
		log.Fatal("No relevant command provided.")
	}

	restArgs := allArgs[1:]

	for i := 0; i < len(allArgs[1:]); i++ {
		arg := restArgs[i]
		switch arg {
		case "--xlsxpath":
			opts.XLSXPath = nextString(restArgs, &i, "xlsx directory path required")
		case "--returnname":
			opts.ReturnName = nextString(restArgs, &i, "return name required")
		case "--import":
			opts.DMPath = nextString(restArgs, &i, "import path required")
		case "--datamapname":
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
