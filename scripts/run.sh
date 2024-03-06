#!/bin/bash

# This is required to handle return pNew; bug in sqlite
export CGO_CFLAGS="-g -O2 -Wno-return-local-addr"

# Ensure correct params passed into the script
if [ ! -n "$1" ] || [ $1 == "help" ]
then
    echo "You must pass a parameter to this script."
    echo "Use of 'test', 'test-all', 'build', 'clean-config', 'dummy-import', 'dummy-datamap-import' and 'debug-datamaps' is currently permitted."
fi

# Here is where we launch all the commands
[[ $1 == "test" ]] && go test ./datamaps
[[ $1 == "test-all" ]] && go test ./...
[[ $1 == "build" ]] &&  go build -o build/datamaps ./cmd/datamaps/main.go
[[ $1 == "clean-config" ]] && rm -r ~/.config/datamaps/
[[ $1 == "dummy-import" ]] && go run ./cmd/datamaps/main.go import --returnname "Hunkers" --datamapname "Tonk 1" --xlsxpath datamaps/testdata/
[[ $1 == "dummy-datamap-import" ]] && go run ./cmd/datamaps/main.go datamap --datamapname "Tonk 1" --import  datamaps/testdata/datamap_matches_test_template.csv
[[ $1 == "debug-datamaps" ]] &&  dlv test ./datamaps/ --wd ./datamaps/
