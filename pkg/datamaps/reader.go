/*
reader reads datamap files and populated spreadsheets
*/
package datamaps

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"github.com/tealeg/xlsx"
	"github.com/yulqen/coords"
)

type (
	// SheetData is the data from the sheet
	SheetData map[string]ExtractedCell
	// FileData is the data from the file.
	FileData map[string]SheetData
	// ExtractedData is the Extraced data from the file, filtered by a Datamap.
	ExtractedData map[string]map[string]xlsx.Cell
)

//DatamapLine - a line from the datamap.
type DatamapLine struct {
	Key     string
	Sheet   string
	Cellref string
}

//ExtractedCell is data pulled from a cell.
type ExtractedCell struct {
	Cell  *xlsx.Cell
	Col   string
	Row   int
	Value string
}

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
// in a slice of DatamapLine structs.
func getSheetNames(dmls []DatamapLine) []string {
	var sheetNames []string
	for _, dml := range dmls {
		if sheetInSlice(sheetNames, dml.Sheet) == false {
			sheetNames = append(sheetNames, dml.Sheet)
		}
	}
	return sheetNames
}

//ReadDML returns a slice of DatamapLine structs.
func ReadDML(path string) ([]DatamapLine, error) {
	var s []DatamapLine
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
		dml := DatamapLine{
			Key:     strings.Trim(record[0], " "),
			Sheet:   strings.Trim(record[1], " "),
			Cellref: strings.Trim(record[2], " ")}
		s = append(s, dml)
	}
	return s, nil
}

//ReadXLSX returns the file's data as a map,
// keyed on sheet name. All values are returned as strings.
// Paths to a datamap and the spreadsheet file required.
func ReadXLSX(ssheet string) FileData {

	// open the files
	data, err := xlsx.OpenFile(ssheet)
	if err != nil {
		log.Fatal(err)
	}
	outer := make(FileData, 1)

	// get the data
	for _, sheet := range data.Sheets {
		inner := make(SheetData)
		for rowLidx, row := range sheet.Rows {
			for colLidx, cell := range row.Cells {
				colStr, err := coords.ColIndexToAlpha(colLidx)
				if err != nil {
					log.Fatal(err)
				}
				ex := ExtractedCell{
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

//Extract returns the file's data as a map,
// using the datamap as a filter, keyed on sheet name. All values
// are returned as strings.
// Paths to a datamap and the spreadsheet file required.
func Extract(dm string, ssheet string) ExtractedData {
	xdata := ReadXLSX(ssheet)
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

//GetTargetFiles finds all xlsx and xlsm files in directory.
func GetTargetFiles(path string) ([]string, error) {
	if lastchar := path[len(path)-1:]; lastchar != string(filepath.Separator) {
		return nil, fmt.Errorf("path must end in a %s character", string(filepath.Separator))
	}
	fullpath := strings.Join([]string{path, "*.xlsx"}, "")
	output, err := filepath.Glob(fullpath)
	if err != nil {
		return nil, err
	}
	if output == nil {
		return nil, fmt.Errorf("cannot find any xlsx files in %s", path)
	}
	return output, nil
}
