debug-datamaps:
	dlv test ./datamaps/ --wd .`/datamaps/

dummy-import:
	./build/datamaps import --returnname "Hunkers" --datamapname "Tonk 1" --xlsxpath datamaps/testdata/

dummy-datamap-import:
	./build/datamaps datamap --datamapname "Tonk 1" --import  datamaps/testdata/datamap_matches_test_template.csv

build:
	go build -o build/datamaps ./cmd/datamaps/main.go

test-all:
	go test ./...

clean-config:
	rm -r ~/.config/datamaps/

godoc:
	godoc -http :6060 -goroot /usr/share/go-1.14/
