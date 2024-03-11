package models

import (
	"database/sql"
	"time"
)

// datamapLine - a line from the datamap.
type datamapLine struct {
	Key     string
	Sheet   string
	Cellref string
}

type Datamap struct {
	ID          int
	Name        string
	Description string
	Created     time.Time
}

type DatamapModel struct {
	DB *sql.DB
}
