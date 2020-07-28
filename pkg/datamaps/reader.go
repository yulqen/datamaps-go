/* datmamaps packages handles datamap files and populated spreadsheets.
 */

package datamaps

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	// Required for the sqlite3 driver
	_ "github.com/mattn/go-sqlite3"

	"github.com/tealeg/xlsx"
	"github.com/yulqen/coords"
)

type (
	// sheetData is the data from the sheet.
	sheetData map[string]extractedCell

	// fileData is the data from the file.
	fileData map[string]sheetData

	// extractedData is the Extracted data from the file, filtered by a Datamap.
	extractedData map[string]map[string]xlsx.Cell
)

//datamapLine - a line from the datamap.
type datamapLine struct {
	Key     string
	Sheet   string
	Cellref string
}

//extractedCell is data pulled from a cell.
type extractedCell struct {
	Cell  *xlsx.Cell
	Col   string
	Row   int
	Value string
}

type extractedDatamapFile []datamapLine

//sheetInSlice is a helper which returns true
// if a string is in a slice of strings.
func sheetInSlice(list []string, key string) bool {
	for _, x := range list {
		if x == key {
			return true
		}
	}

	return false
}

//getSheetNames returns the number of Sheet field entries
// in a slice of datamapLine structs.
func getSheetNames(dmls extractedDatamapFile) []string {
	var sheetNames []string

	for _, dml := range dmls {
		if !sheetInSlice(sheetNames, dml.Sheet) {
			sheetNames = append(sheetNames, dml.Sheet)
		}
	}

	return sheetNames
}

// ReadDML returns a slice of datamapLine structs given a
// path to a datamap file.
func ReadDML(path string) (extractedDatamapFile, error) {
	var s extractedDatamapFile

	data, err := ioutil.ReadFile(path)

	if err != nil {
		return s, fmt.Errorf("Cannot find file: %s", path)
	}

	r := csv.NewReader(strings.NewReader(string(data)))

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return s, errors.New("Cannot read line %s")
		}

		if record[0] == "cell_key" {
			// this must be the header
			continue
		}

		dml := datamapLine{
			Key:     strings.Trim(record[0], " "),
			Sheet:   strings.Trim(record[1], " "),
			Cellref: strings.Trim(record[2], " ")}
		s = append(s, dml)
	}

	return s, nil
}

// ReadXLSX returns a file at path's data as a map,
// keyed on sheet name. All values are returned as strings.
// Paths to a datamap and the spreadsheet file required.
func ReadXLSX(path string) fileData {
	// open the files
	data, err := xlsx.OpenFile(path)
	if err != nil {
		log.Fatal(err)
	}

	outer := make(fileData, 1)

	// get the data
	for _, sheet := range data.Sheets {
		inner := make(sheetData)

		for rowLidx, row := range sheet.Rows {
			for colLidx, cell := range row.Cells {
				colStr, err := coords.ColIndexToAlpha(colLidx)
				if err != nil {
					log.Fatal(err)
				}

				ex := extractedCell{
					Cell:  cell,
					Col:   colStr,
					Row:   rowLidx + 1,
					Value: cell.Value}
				cellref := fmt.Sprintf("%s%d", ex.Col, ex.Row)
				inner[cellref] = ex
			}

			outer[sheet.Name] = inner
		}
	}

	return outer
}

// DatamapFromDB creates an extractedDatamapFile from the database given
// the name of a datamap. Of course, in this instance, the data is not
// coming from a datamap file (such as datamap.csv) but from datamap data
// previous stored in the database by DatamapToDB or similar.
func DatamapFromDB(name string, db *sql.DB) (extractedDatamapFile, error) {

	var out extractedDatamapFile

	query := `
	select
		key, sheet, cellref
	from datamap_line
		join datamap on datamap_line.dm_id = datamap.id where datamap.name = ?;
	`
	rows, err := db.Query(query, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			key     string
			sheet   string
			cellref string
		)
		if err := rows.Scan(&key, &sheet, &cellref); err != nil {
			return nil, err
		}

		out = append(out, datamapLine{Key: key, Sheet: sheet, Cellref: cellref})
	}

	return out, nil
}

// ExtractDBDatamap uses a datamap named from the database db to extract values
// from the populated spreadsheet file file.
func ExtractDBDatamap(name string, file string, db *sql.DB) (extractedData, error) {
	xdata := ReadXLSX(file)
	ddata, err := DatamapFromDB(name, db) // this will need to return an extractedDatamapFile

	if err != nil {
		return nil, err
	}

	names := getSheetNames(ddata)
	outer := make(extractedData, len(names))
	var inner map[string]xlsx.Cell

	seen := make(map[string]bool) // using this as a set to distinguish sheet names

	for _, i := range ddata {
		sheet := i.Sheet
		if !seen[sheet] {
			seen[sheet] = true
			inner = make(map[string]xlsx.Cell) // we create a new inner when we get a new sheet
		}
		cellref := i.Cellref

		if val, ok := xdata[sheet][cellref]; ok {
			inner[cellref] = *val.Cell
			outer[sheet] = inner
		}
	}

	return outer, nil
}

// extract returns the file at path's data as a map,
// using the datamap as a filter, keyed on sheet name. All values
// are returned as strings. (Currently deprecated in favour of
// ExtractDBDatamap.
func extract(dm string, path string) extractedData {
	xdata := ReadXLSX(path)
	ddata, err := ReadDML(dm)

	if err != nil {
		log.Fatal(err)
	}

	names := getSheetNames(ddata)
	outer := make(extractedData, len(names))
	inner := make(map[string]xlsx.Cell)

	for _, i := range ddata {
		sheet := i.Sheet
		cellref := i.Cellref

		if val, ok := xdata[sheet][cellref]; ok {
			inner[cellref] = *val.Cell
			outer[sheet] = inner
		}
	}

	return outer
}

//getTargetFiles finds all xlsx and xlsm files in directory.
func getTargetFiles(path string) ([]string, error) {
	if lastchar := path[len(path)-1:]; lastchar != string(filepath.Separator) {
		return nil, fmt.Errorf("path must end in a %s character", string(filepath.Separator))
	}

	fullpath := strings.Join([]string{path, "*.xls[xm]"}, "")
	output, err := filepath.Glob(fullpath)

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, fmt.Errorf("cannot find any xlsx files in %s", path)
	}

	return output, nil
}
