### Generating the sqlite3 database:
- Download the csv data from https://opensky-network.org/datasets/metadata/
- Strip not needed columns from the csv, (can be done very quickly in Excel)
- Install sqlite3
- In the folder with the csv run: 
```
sqlite3 aircraft_info.db
.mode csv
.import <aircraftDatabase.csv> aircraft_info
.quit
```
- That should generate an `aircraft_info.db` file containing the `aircraft_info` table.