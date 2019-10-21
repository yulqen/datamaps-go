package reader

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
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

//ReadDML returns a pointer to a slice of DatamapLine structs
func ReadDML(path string) (*[]DatamapLine, error) {
	var s []DatamapLine
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return &s, &fileError{path, "Cannot open."}
	}
	r := csv.NewReader(strings.NewReader(string(data)))
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Cannot read line %s, ", err)
		}
		if record[0] == "cell_key" {
			// this must be the header
			continue
		}
		dml := DatamapLine{Key: record[0], Sheet: record[1], Cellref: record[2]}
		s = append(s, dml)
		// fmt.Printf("Key: %s; sheet: %s cellref: %s\n", dml.Key, dml.Sheet, dml.Cellref)
		// klen, slen := Keylens(dml)
		// fmt.Printf("Key length: %d\n", klen)
		// fmt.Printf("Sheet length: %d\n\n", slen)
	}
	return &s, nil
}

//ReadXLSX reads an XLSX file
func ReadXLSX(excelFileName string) {
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		fmt.Printf("Cannot open %s", excelFileName)
	}
	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {
			for _, cell := range row.Cells {
				text := cell.String()
				fmt.Printf("Sheet: %s\nValue: %s\n", sheet.Name, text)
			}
		}
	}
}
