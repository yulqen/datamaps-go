debug-datamaps:
	dlv test ./pkg/datamaps/ --wd ./pkg/datamaps/

dummy-import:
	./datamaps import --returnname "Hunkers" --datamapname "Tonk 1" --xlsxpath pkg/datamaps/testdata/

dummy-datamap-import:
	./datamaps datamap --datamapname "Tonk 1" --import  pkg/datamaps/testdata/short/datamap_matches_test_template.csv

build:
	go build -o datamaps ./cmd/datamaps/main.go

test-all:
	go test ./...

clean-config:
	rm -r ~/.config/datamaps/
