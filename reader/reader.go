package reader

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"github.com/tealeg/xlsx"
)

const (
	maxCols      = 16384
	maxAlphabets = (maxCols / 26) - 1
)

var colstream = cols(maxAlphabets)

//DatamapLine - a line from the datamap.
type DatamapLine struct {
	Key     string
	Sheet   string
	Cellref string
}

//ExtractedCell is Data pulled from a cell
type ExtractedCell struct {
	cell    *xlsx.Cell
	colL    string
	rowLidx int
	value   string
}

//Keylens returns the length of a key
func Keylens(dml DatamapLine) (int, int) {
	return len(dml.Key), len(dml.Sheet)
}

// sheetInSlice is a helper which returns true
// if a string is in a slice of strings
func sheetInSlice(list []string, key string) bool {
	for _, x := range list {
		if x == key {
			return true
		}
	}
	return false
}

// getSheetNames returns the number of Sheet field entries
// in a slice of DatamapLine structs
func getSheetNames(dmls []DatamapLine) []string {
	var sheetNames []string
	for _, dml := range dmls {
		if sheetInSlice(sheetNames, dml.Sheet) == false {
			sheetNames = append(sheetNames, dml.Sheet)
		}
	}
	fmt.Printf("sheetNames is %v", sheetNames)
	return sheetNames
}

//ReadDML returns a slice of DatamapLine structs
func ReadDML(path string) ([]DatamapLine, error) {
	var s []DatamapLine
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return s, errors.New("Cannot find file")
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

//alphabet generates all the letters of the alphabet
func alphabet() []string {
	letters := make([]string, 26)
	for idx := range letters {
		letters[idx] = string('A' + byte(idx))
	}
	return letters
}

//cols generates the alpha column compont of Excel cell references
//Adds n alphabets to the first (A..Z) alphabet.
func cols(n int) []string {
	out := alphabet()
	alen := len(out)
	tmp := make([]string, alen)
	copy(tmp, out)
	for cycle := 0; cycle < n; cycle++ {
		for y := 0; y < alen; y++ {
			out = append(out, out[(cycle+2)-2]+tmp[y])
		}
	}
	return out
}

type (
	SheetData map[string]ExtractedCell
	FileData  map[string]SheetData
)

func ReadXLSToMap(dm string, tm string) FileData {

	// open the files
	excelData, err := xlsx.OpenFile(tm)
	if err != nil {
		log.Fatalf("Cannot open %v", excelData)
	}

	dmlData, err := ReadDML(dm)
	if err != nil {
		log.Fatalf("Cannot read datamap file %v", dm)
	}

	sheetNames := getSheetNames(dmlData)

	output := make(FileData, len(sheetNames))

	// get the data
	for _, sheet := range excelData.Sheets {
		data := make(SheetData)
		for rowLidx, row := range sheet.Rows {
			for colLidx, cell := range row.Cells {
				ex := ExtractedCell{
					cell:    cell,
					colL:    colstream[colLidx],
					rowLidx: rowLidx + 1,
					value:   cell.Value}
				cellref := fmt.Sprintf("%s%d", ex.colL, ex.rowLidx)
				data[cellref] = ex
			}
			output[sheet.Name] = data
		}
	}
	return output
}

//ReadXLSX reads an XLSX file
func ReadXLSX(fn string) []ExtractedCell {
	var out []ExtractedCell
	f, err := xlsx.OpenFile(fn)
	if err != nil {
		fmt.Printf("Cannot open %s", fn)
	}
	for _, sheet := range f.Sheets {
		for rowLidx, row := range sheet.Rows {
			for colLidx, cell := range row.Cells {
				ex := ExtractedCell{
					cell:    cell,
					colL:    colstream[colLidx],
					rowLidx: rowLidx + 1,
					value:   cell.Value}
				out = append(out, ex)
				text := cell.String()
				log.Printf("Sheet: %s Row: %d Col: %q Value: %s\n",
					sheet.Name, ex.rowLidx, colstream[colLidx], text)
			}
		}
	}
	return out
}
