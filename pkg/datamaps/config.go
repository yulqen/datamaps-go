package datamaps

import (
	"os"
	"path/filepath"
)

const (
	config_dir_name = "datamaps-go"
	db_name         = "datamaps.db"
)

func getUserConfigDir() (string, error) {
	dir, err := os.UserConfigDir()
	config_path := filepath.Join(dir, config_dir_name)
	if err != nil {
		return "", err
	}
	return config_path, nil
}

type Options struct {
	DBPath      string
	DMPath      string
	DMName      string
	DMOverwrite bool
	DMInitial   bool
	DMData      []DatamapLine
}

func defaultOptions() *Options {
	return &Options{
		DBPath:      "PATH TO DB",
		DMPath:      "PATH TO DATAMAP",
		DMName:      "Unnamed Datamap",
		DMOverwrite: false,
		DMInitial:   false,
		DMData:      make([]DatamapLine, 0),
	}
}

func ParseOptions() *Options {
	opts := defaultOptions()
	return opts

}
