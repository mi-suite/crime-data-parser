package cmd

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/lib/pq"
)

var db *sql.DB

// SetupDatabaseConnection with sql and postgres
func setupDatabaseConnection(databaseURL string) *sql.DB {
	dbURL, err := pq.ParseURL(databaseURL)

	if err != nil {
		log.Fatal("Error parsing DB URL:::", err)
	}

	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error opening connection to db:::", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func createTable(db *sql.DB, tableName string) {
	var query = fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS "%s" (
            "SN" INTEGER NOT NULL UNIQUE,
            "ID" INTEGER NOT NULL,
            "Case Number" CHARACTER VARYING,
            "Date" TIMESTAMP WITHOUT TIME ZONE,
            "Block" CHARACTER VARYING NOT NULL,
            "IUCR" CHARACTER VARYING NOT NULL,
            "Primary Type" CHARACTER VARYING NOT NULL,
            "Description" CHARACTER VARYING NOT NULL,
            "Location Description" CHARACTER VARYING,
            "Arrest" BOOLEAN NOT NULL,
            "Domestic" BOOLEAN NOT NULL,
            "Beat" INTEGER NOT NULL,
            "District" INTEGER,
            "Ward" INTEGER,
            "Community Area" INTEGER,
            "FBI Code" CHARACTER VARYING NOT NULL,
            "X Coordinate" INTEGER,
            "Y Coordinate" INTEGER,
            "Year" INTEGER NOT NULL,
            "Updated On" timestamp without time zone,
            "Latitude" FLOAT,
            "Longitude" FLOAT,
            "Location" CHARACTER VARYING
        )
    `, tableName)

	_, err := db.Query(query)
	if err != nil {
		panic(err)
	}

	return
}

func insertIntoTable(db *sql.DB, data *CSVData) {
	query := `
		INSERT into "crimes.all" (
            "SN", "ID", "Case Number", "Date", "Block", "IUCR", "Primary Type", "Description",
            "Location Description", "Arrest", "Domestic", "Beat", "District", "Ward", "Community Area",
            "FBI Code", "X Coordinate", "Y Coordinate", "Year", "Updated On", "Latitude",
            "Longitude", "Location"
        )
        VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
            $17, $18, $19, $20, $21, $22, $23
        )
        ON CONFLICT ("SN") DO NOTHING
        RETURNING *;
    `

	var ScanClient ScanClient

	err := db.QueryRow(
		query,
		data.SN, data.ID, data.CaseNumber, data.Date.Time, data.Block, data.IUCR, data.PrimaryType,
		data.Description, data.LocationDescription, data.Arrest, data.Domestic, data.Beat, data.District,
		data.Ward, data.CommunityArea, data.FBICode, data.XCoordinate, data.YCoordinate,
		data.Year, data.UpdatedOn.Time, data.Latitude, data.Longitude, data.Location,
	).Scan(
		&ScanClient.SN, &ScanClient.ID, &ScanClient.CaseNumber, &ScanClient.Date, &ScanClient.Block, &ScanClient.IUCR, &ScanClient.PrimaryType,
		&ScanClient.Description, &ScanClient.LocationDescription, &ScanClient.Arrest, &ScanClient.Domestic, &ScanClient.Beat, &ScanClient.District,
		&ScanClient.Ward, &ScanClient.CommunityArea, &ScanClient.FBICode, &ScanClient.XCoordinate, &ScanClient.YCoordinate,
		&ScanClient.Year, &ScanClient.UpdatedOn, &ScanClient.Latitude, &ScanClient.Longitude, &ScanClient.Location,
	)

	if err != nil && err.Error() != "sql: no rows in result set" {
		log.Fatal(err)
	}
}
