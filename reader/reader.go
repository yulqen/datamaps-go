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

const maxCols = 16384
const maxAlphabets = (maxCols / 26) - 1

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
	colL    string
	rowLidx int
	value   string
}

//alphSingle generates all the letters of the alphabet
func alphaSingle() []string {
	letters := make([]string, 26)
	for idx := range letters {
		letters[idx] = string('A' + byte(idx))
	}
	return letters
}

var alphaStream = alphas(maxAlphabets)

//alphas generates the alpha column compont of Excel cell references
//Adds n alphabets to the first (A..Z) alphabet.
func alphas(n int) []string {
	single := alphaSingle()
	slen := len(single)
	tmp := make([]string, len(single))
	copy(tmp, single)
	for cycle := 0; cycle < n; cycle++ {
		for y := 0; y < slen; y++ {
			single = append(single, single[(cycle+2)-2]+tmp[y])
		}
	}
	return single
}

//ReadXLSX reads an XLSX file
func ReadXLSX(excelFileName string) []ExtractedCell {
	var out []ExtractedCell
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		fmt.Printf("Cannot open %s", excelFileName)
	}
	for _, sheet := range xlFile.Sheets {
		for rowLidx, row := range sheet.Rows {
			for colLidx, cell := range row.Cells {
				ex := ExtractedCell{
					cell:    cell,
					colL:    alphaStream[colLidx],
					rowLidx: rowLidx + 1,
					value:   cell.Value}
				out = append(out, ex)
				text := cell.String()
				log.Printf("Sheet: %s Row: %d Col: %q Value: %s\n",
					sheet.Name, ex.rowLidx, alphaStream[colLidx], text)
			}
		}
	}
	return out
}
