package datamaps

import (
	"database/sql"
	"log"
	"path/filepath"

	"github.com/tealeg/xlsx/v3"

	// Needed for the sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
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
	//
	// SELECT datamap_line.key, return_data.value, return_data.filename
	//                                    FROM (((return_data
	//                                    INNER JOIN datamap_line ON return_data.dml_id=datamap_line.id)
	//                                    INNER JOIN datamap ON datamap_line.dm_id=datamap.id)
	//                                    INNER JOIN return on return_data.ret_id=return.id)
	//                                    WHERE datamap.name="Tonk 1" AND return.name="Hunkers";

	// filename := filepath.Join(opts.MasterOutPutPath, "master.xlsx")

	// a test key
	targetKey := "A Rabbit"

	wb := xlsx.NewFile()
	sh, err := wb.AddSheet("Master Data")

	testRow, err := sh.Row(1)
	if err != nil {
		log.Fatal(err)
	}

	// SQLITE CODE

	db, err := sql.Open("sqlite3", opts.DBPath)

	// Get number amount of datamap keys in target datamap
	keyCountRows := db.QueryRow("SELECT count(key) FROM datamap_line, datamap WHERE datamap.name=?;", opts.DMName)
	if err != nil {
		log.Fatal(err)
	}

	var datamapKeysNumber int64
	if err := keyCountRows.Scan(&datamapKeysNumber); err != nil {
		return err
	}

	datamapKeysRows, err := db.Query("SELECT key FROM datamap_line, datamap WHERE datamap.name=?;", opts.DMName)
	if err != nil {
		log.Fatal(err)
	}

	var datamapKeys []string

	var i int64
	for i = 0; i < datamapKeysNumber; i++ {
		if k := datamapKeysRows.Next(); k {
			var key string
			if err := datamapKeysRows.Scan(&key); err != nil {
				return err
			}
			datamapKeys = append(datamapKeys, key)
		}
	}

	sqlCount := `SELECT count(return_data.id)
                                          FROM (((return_data
                                          INNER JOIN datamap_line ON return_data.dml_id=datamap_line.id)
                                          INNER JOIN datamap ON datamap_line.dm_id=datamap.id)
                                          INNER JOIN return on return_data.ret_id=return.id)
                                          WHERE datamap.name=? AND return.name=? AND datamap_line.key=?
                                          GROUP BY datamap_line.key;`

	var rowCount int64
	rowCountRes := db.QueryRow(sqlCount, opts.DMName, opts.ReturnName, targetKey)
	if err != nil {
		return err
	}

	if err := rowCountRes.Scan(&rowCount); err != nil {
		return err
	}

	getDataSQL := `SELECT datamap_line.key, return_data.value, return_data.filename
                                          FROM (((return_data
                                          INNER JOIN datamap_line ON return_data.dml_id=datamap_line.id) 
                                          INNER JOIN datamap ON datamap_line.dm_id=datamap.id) 
                                          INNER JOIN return on return_data.ret_id=return.id) 
                                          WHERE datamap.name=? AND return.name=? AND datamap_line.key=?;`

	var values = make([]string, rowCount+1)
	for _, k := range datamapKeys {
		masterData, err := db.Query(getDataSQL, opts.DMName, opts.ReturnName, k)
		if err != nil {
			return err
		}
		var x int64
		for x = 0; x < rowCount; x++ {
			var key, filename, value string
			if b := masterData.Next(); b {
				if err := masterData.Scan(&key, &value, &filename); err != nil {
					return err
				}
				values = append(values, value)
			}
		}
	}

	if sl := testRow.WriteSlice(values, -1); sl == -1 {
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
