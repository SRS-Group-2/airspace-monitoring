### Generating the sqlite3 database:
- Download the csv data from https://opensky-network.org/datasets/metadata/
- Strip not needed columns from the csv, (can be done very quickly in Excel)
- Remove duplicated icao24 from csv: <https://support.microsoft.com/en-us/office/find-and-remove-duplicates-00e35bea-b46a-4d5d-b28e-66a552dc138d>
- Install sqlite3
- In the folder with the csv run: 
```
sqlite3 aircraft_info.db
.mode csv
.import <aircraftDatabase.csv> aircraft_info
CREATE UNIQUE INDEX idx_icao24 ON aircraft_info (icao24);
.quit
```
- That should generate an `aircraft_info.db` file containing the `aircraft_info` table.
