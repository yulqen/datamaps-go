package datamaps

import (
	"log"
	"path/filepath"

	"github.com/tealeg/xlsx/v3"
)

func ExportMaster(opts *Options) error {
	// A master represents a set of file data from the database. Actually, in terms of the database,
	// it should represent a "return".
	//
	// The meat of the master is of the format:
	//		Key	1 | Key_1_Value_for_FD_1 | Key_1_Value_for_FD_2 | Key_1_Value_for_FD_3 | ... etc
	//		Key 2 | Key_2_Value_for_FD_1 | Key_2_Value_for_FD_2 | Key_2_Value_for_FD_3 | ... etc
	//		Key 3 | Key_3_Value_for_FD_1 | Key_3_Value_for_FD_2 | Key_3_Value_for_FD_3 | ... etc
	//		...
	// Could be represented as a slice or a map[string][]string
	// SQL statement:
	//		SELECT datamap_line.key , return_data.value, return_data.filename
	//		FROM return_data INNER JOIN datamap_line on return_data.dml_id=datamap_line.id
	//		WHERE "key";

	// filename := filepath.Join(opts.MasterOutPutPath, "master.xlsx")
	wb := xlsx.NewFile()
	sh, err := wb.AddSheet("Master Data")

	testRow, err := sh.Row(1)
	if err != nil {
		log.Fatal(err)
	}
	if sl := testRow.WriteSlice([]string{"Hello", "Bollocks", "Knackers", "Bottyies"}, -1); sl == -1 {
		log.Printf("not a slice type")
	}
	log.Printf("writing slice to row\n")

	log.Printf("saving master at %s", opts.MasterOutPutPath)
	if err := wb.Save(filepath.Join(opts.MasterOutPutPath, "master.xlsx")); err != nil {
		log.Fatalf("cannot save file to %s", opts.MasterOutPutPath)
	}
	sh.Close()
	return nil
}
