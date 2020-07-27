# Datamaps in Go

Some useful queries:

### Joining and pattern matching
`select name, key, cellref from datamap_line join datamap on (datamap.id
= datamap_line.dm_id) where key like 'Total RDEL%';`

### Create return and return_data tables

`CREATE TABLE return(id INTEGER PRIMARY KEY, name TEXT, date_created TEXT);
CREATE TABLE return_data(id INTEGER PRIMARY KEY, dml_id INTEGER, value TEXT, FOREIGN KEY (dml_id) REFERENCES datamap_line(id) ON DELETE CASCADE);
`
