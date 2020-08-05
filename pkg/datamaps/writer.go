package datamaps

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	"github.com/tealeg/xlsx/v3"

	// Needed for the sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

// CreateMaster creates a master spreadsheet for a specific return given,
// based on a datamap name - both of which already need to be in the database,
// along with the data associated with the return. The datamap and return data
// must already have been imported.
func CreateMaster(opts *Options) error {
	wb := xlsx.NewFile()
	sh, err := wb.AddSheet("Master Data")
	if err != nil {
		return fmt.Errorf("cannot add 'Master Data' sheet to new XLSX file: %v", err)
	}

	db, err := sql.Open("sqlite3", opts.DBPath)
	if err != nil {
		return fmt.Errorf("cannot open database %v", err)
	}

	// Get number amount of datamap keys in target datamap
	keyCountRows := db.QueryRow("SELECT count(key) FROM datamap_line, datamap WHERE datamap.name=?;", opts.DMName)

	var datamapKeysNumber int64
	if err := keyCountRows.Scan(&datamapKeysNumber); err != nil {
		return err
	}

	datamapKeysRows, err := db.Query("SELECT key FROM datamap_line, datamap WHERE datamap.name=?;", opts.DMName)
	if err != nil {
		return fmt.Errorf("cannot query for keys in database - %v", err)
	}

	var datamapKeys []string

	var i int64
	for i = 0; i < datamapKeysNumber; i++ {
		if k := datamapKeysRows.Next(); k {
			var key string
			if err := datamapKeysRows.Scan(&key); err != nil {
				return fmt.Errorf("cannot Scan for key %s - %v", key, err)
			}
			datamapKeys = append(datamapKeys, key)
		}
	}

	sqlCount := `SELECT count(return_data.id)
                                          FROM (((return_data
                                          INNER JOIN datamap_line ON return_data.dml_id=datamap_line.id)
                                          INNER JOIN datamap ON datamap_line.dm_id=datamap.id)
                                          INNER JOIN return on return_data.ret_id=return.id)
                                          WHERE datamap.name=? AND return.name=?;`

	var rowCount int64
	rowCountRes := db.QueryRow(sqlCount, opts.DMName, opts.ReturnName)
	if err != nil {
		return fmt.Errorf("cannot query for row count of return data - %v", err)
	}

	if err := rowCountRes.Scan(&rowCount); err != nil {
		return err
	}

	getDataSQL := `SELECT datamap_line.key, return_data.value, return_data.filename
                                          FROM (((return_data
                                          INNER JOIN datamap_line ON return_data.dml_id=datamap_line.id) 
                                          INNER JOIN datamap ON datamap_line.dm_id=datamap.id) 
                                          INNER JOIN return on return_data.ret_id=return.id) 
                                          WHERE datamap.name=? AND return.name=? AND datamap_line.key=?
										  ORDER BY return_data.filename;`

	seen := make(map[string]struct{}) // homemade Set https://emersion.fr/blog/2017/sets-in-go/

	var values = make(map[string][]string)
	headerSlice := make([]string, 0)
	for _, k := range datamapKeys {
		masterData, err := db.Query(getDataSQL, opts.DMName, opts.ReturnName, k)
		if err != nil {
			return err
		}
		for masterData.Next() {
			var key, filename, value string
			if err := masterData.Scan(&key, &value, &filename); err != nil {
				return err
			}
			values, err = appendValueMap(key, value, values)
			if _, ok := seen[filename]; !ok {
				headerSlice = append(headerSlice, filename)
				seen[filename] = struct{}{}
			}
			if err != nil {
				return err
			}
		}
	}

	for masterRow := 0; masterRow <= len(datamapKeys); masterRow++ {
		r, err := sh.AddRowAtIndex(masterRow)
		if err != nil {
			return fmt.Errorf("cannot create row %d in output spreadsheet: %v", masterRow, err)
		}
		if masterRow == 0 {
			if hdr := r.WriteSlice(append([]string{""}, headerSlice...), -1); hdr == -1 {
				return fmt.Errorf("cannot write header values into header row: %v", err)
			}
			continue
		}
		dmlKey := datamapKeys[masterRow-1]
		if sl := r.WriteSlice(append([]string{dmlKey}, values[dmlKey]...), -1); sl == -1 {
			log.Printf("not a slice type")
		}
	}

	log.Printf("saving master at %s", opts.MasterOutPutPath)
	if err := wb.Save(filepath.Join(opts.MasterOutPutPath, "master.xlsx")); err != nil {
		log.Fatalf("cannot save file to %s", opts.MasterOutPutPath)
	}
	sh.Close()
	return nil
}

func appendValueMap(k, v string, values map[string][]string) (map[string][]string, error) {

	var keyIn bool
	for kv := range values {
		if kv == k {
			keyIn = true
			break
		}
	}

	if keyIn {
		slice, ok := values[k]
		if !ok {
			return nil, fmt.Errorf("%s is not a key in this map", k)
		}
		slice = append(slice, v)
		values[k] = slice
		return values, nil
	} else {
		values[k] = []string{v}
		return values, nil
	}
}
