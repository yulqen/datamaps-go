/*datamaps is a simple tool to extract from and send data to spreadsheets.
 */
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/yulqen/datamaps-go/pkg/datamaps"
)

func main() {
	opts := datamaps.ParseOptions()
	if opts.Command == "help" {
		os.Stdout.WriteString(datamaps.Usage)
		os.Exit(0)
	}
	dbpc := datamaps.NewDBPathChecker(os.UserConfigDir)
	if !dbpc.Check() {
		datamaps.SetUp()
	}
	switch opts.Command {
	case "import":
		if err := datamaps.ImportToDB(opts); err != nil {
			log.Fatal(err)
		}
	case "datamap":
		if err := datamaps.DatamapToDB(opts); err != nil {
			log.Fatal(err)
		}
	case "setup":
		_, err := datamaps.SetUp()
		if err != nil {
			log.Fatal(err)
		}
	case "createmaster":
		if err := datamaps.CreateMaster(opts); err != nil {
			log.Fatal(err)
		}
	case "server":
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				http.NotFound(w, r)
				return
			}
			fmt.Fprintf(w, "Welcome to datamaps!")
			// or you could write it thus
			// w.Write([]byte("Hello from datamaps"))
		})
		log.Println("Starting server on :8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}
}
