package main

import (
	"datamaps-go/reader"
)

func main() {
	//reader.ReadDML("/home/lemon/Documents/datamaps/input/datamap.csv")
	data, err := reader.ReadXLSX("/home/lemon/Documents/datamaps/input/A417%20Air%20Balloon_Q1%20Apr%20-%20June%202019_Return.xlsm")
}
