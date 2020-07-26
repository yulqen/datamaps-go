# Datamaps in Go

Some useful queries:

### Joining and pattern matching
`select name, key, cellref from datamap_line join datamap on (datamap.id
= datamap_line.dm_id) where key like 'Total RDEL%';`
