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

//DatamapLine - a line from the datamap.
type DatamapLine struct {
	Key     string
	Sheet   string
	Cellref string
}

type fileError struct {
	file string
	msg  string
}

func (e *fileError) Error() string {
	return fmt.Sprintf("%s", e.msg)
}

//Keylens returns the length of a key
func Keylens(dml DatamapLine) (int, int) {
	return len(dml.Key), len(dml.Sheet)
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

//ExtractedCell is Data pulled from a cell
type ExtractedCell struct {
	cell    *xlsx.Cell
	colLidx string
	rowLidx int
	value   string
}

var lowerAlpha [2]string = [2]string{"A", "B"}

// two alphabet lengths
const ablen int = 52

func alphas() []string {
	acount := 0
	lets := make([]string, ablen)
	for i, _ := range lets {
		if i == 0 {
			lets[i] = "A"
			continue
		}
		if i%26 != 0 {
			lets[i] = string('A' + byte(i))
		} else {
			acount++
			lets[i] = string(64+acount) + string('A'+byte(i))
		}
	}
	return lets
}

//ReadXLSX reads an XLSX file
func ReadXLSX(excelFileName string) []ExtractedCell {
	var out []ExtractedCell
	alphs := alphas()
	for _, l := range alphs {
		fmt.Printf("Letter: %s", string(l))
	}
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		fmt.Printf("Cannot open %s", excelFileName)
	}
	for _, sheet := range xlFile.Sheets {
		for rowLidx, row := range sheet.Rows {
			for colLidx, cell := range row.Cells {
				ex := ExtractedCell{
					cell:    cell,
					colLidx: string(alphs[colLidx]),
					rowLidx: rowLidx + 1,
					value:   cell.Value}
				out = append(out, ex)
				text := cell.String()
				log.Printf("Sheet: %s Row: %d Col: %q Value: %s\n",
					sheet.Name, rowLidx, alphs[colLidx], text)
			}
		}
	}
	return out
}
