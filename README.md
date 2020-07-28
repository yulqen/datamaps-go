# Datamaps in Go

Some useful queries:

### Joining and pattern matching
`select name, key, cellref from datamap_line join datamap on (datamap.id
= datamap_line.dm_id) where key like 'Total RDEL%';`

### Create return and return_data tables

`CREATE TABLE return(id INTEGER PRIMARY KEY, name TEXT, date_created TEXT);
CREATE TABLE return_data(id INTEGER PRIMARY KEY, dml_id INTEGER, value TEXT, FOREIGN KEY (dml_id) REFERENCES datamap_line(id) ON DELETE CASCADE);
`
### This crazy join works
`select dm.name, dml.key, dml.sheet, return.name, return_data.value from
datamap as dm inner join datamap_line as dml on dml.dm_id=dm.id inner join
return on return_data.ret_id=return.id inner join return_data on
return_data.ret_id=return.id;`

### More accurate SQL to test all data at this stage:
`select datamap.name, datamap_line.key, datamap_line.sheet, return.name,
return_data.filename, return_data.value from datamap, datamap_line, return,
return_data where datamap_line.dm_id=datamap.id;`
