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

	"github.com/tealeg/xlsx/v3"
)

type (
	// sheetData is the data from the sheet.
	sheetData map[string]extractedCell

	// FileData is the data from the file.
	FileData map[string]sheetData

	// ExtractedData is the Extracted data from the file, filtered by a Datamap.
	ExtractedData map[string]map[string]xlsx.Cell
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

var (
	inner = make(sheetData)
	exc   extractedCell
)

// ExtractedDatamapFile is a slice of datamapLine structs, each of which encodes a single line
// in the datamap file/database table.
type ExtractedDatamapFile []datamapLine

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
func getSheetNames(dmls ExtractedDatamapFile) []string {
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
func ReadDML(path string) (ExtractedDatamapFile, error) {
	var s ExtractedDatamapFile

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

// cellVisitor is used by datamaps.rowVisitor() and is called
// on every cell in the target xlsx file in order to extract
// the data.
func cellVisitor(c *xlsx.Cell) error {
	x, y := c.GetCoordinates()
	cellref := xlsx.GetCellIDStringFromCoords(x, y)

	ex := extractedCell{
		Cell:  c,
		Value: c.Value,
	}

	inner[cellref] = ex

	return nil
}

// rowVisitor is used as a callback by xlsx.sheet.ForEachRow(). It wraps
// a call to xlsx.Row.ForEachCell() which actually extracts the data.
func rowVisitor(r *xlsx.Row) error {
	if err := r.ForEachCell(cellVisitor, xlsx.SkipEmptyCells); err != nil {
		return err
	}
	return nil
}

// ReadXLSX returns a file at path's data as a map,
// keyed on sheet name. All values are returned as strings.
// Paths to a datamap and the spreadsheet file required.
func ReadXLSX(path string) FileData {
	wb, err := xlsx.OpenFile(path)
	if err != nil {
		log.Fatal(err)
	}

	outer := make(FileData, 1)

	// get the data
	for _, sheet := range wb.Sheets {

		if err := sheet.ForEachRow(rowVisitor); err != nil {
			log.Fatal(err)
		}
		outer[sheet.Name] = inner
		inner = make(sheetData)
	}

	return outer
}

// DatamapFromDB creates an ExtractedDatamapFile from the database given
// the name of a datamap. Of course, in this instance, the data is not
// coming from a datamap file (such as datamap.csv) but from datamap data
// previous stored in the database by DatamapToDB or similar.
func DatamapFromDB(name string, db *sql.DB) (ExtractedDatamapFile, error) {

	var out ExtractedDatamapFile

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
func ExtractDBDatamap(name string, file string, db *sql.DB) (ExtractedData, error) {
	ddata, err := DatamapFromDB(name, db) // this will need to return an ExtractedDatamapFile
	if err != nil {
		return nil, err
	}
	if len(ddata) == 0 {
		return nil, fmt.Errorf("there is no datamap in the database matching name '%s'. Try running 'datamaps datamap --import...'", name)
	}
	xdata := ReadXLSX(file)

	names := getSheetNames(ddata)
	outer := make(ExtractedData, len(names))
	// var inner map[string]xlsx.Cell

	for _, s := range names {
		outer[s] = make(map[string]xlsx.Cell)
	}

	for _, i := range ddata {
		sheet := i.Sheet
		cellref := i.Cellref

		if val, ok := xdata[sheet][cellref]; ok {
			outer[sheet][cellref] = *val.Cell
		}
	}

	return outer, nil
}

// extract returns the file at path's data as a map,
// using the datamap as a filter, keyed on sheet name. All values
// are returned as strings. (Currently deprecated in favour of
// ExtractDBDatamap.
func extract(dm string, path string) ExtractedData {
	xdata := ReadXLSX(path)
	ddata, err := ReadDML(dm)

	if err != nil {
		log.Fatal(err)
	}

	names := getSheetNames(ddata)
	outer := make(ExtractedData, len(names))
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
