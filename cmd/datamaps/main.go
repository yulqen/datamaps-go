/*
datamaps-go is a simple tool to extract from and send data to spreadsheets.
*/
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yulqen/datamaps-go/pkg/datamaps"
)

func main() {
	// env := datamaps.DetectConfig()
	// if !env {
	// 	datamaps.SetUp()
	// }
	opts := datamaps.ParseOptions()

	switch opts.Command {
	case "datamap":
		if err := datamaps.DatamapToDB(opts); err != nil {
			log.Fatal(err)
		}
	case "setup":
		_, err := datamaps.SetUp()
		if err != nil {
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
