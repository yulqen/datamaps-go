#!/usr/bin/env bash

echo "Deleting datamaps.db..."
if [ -f ~/.config/datamaps/datampaps.db ] ; then
  rm ~/.config/datamaps/datamaps.db
fi
#echo "Importing datamap.csv from testdata directory..."
#./bin/datamaps -- datamap --import ~/Documents/datamaps/input/datamap.csv --datamapname "Sept 2020"