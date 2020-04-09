/*
Copyright © 2020 csvPath HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/lib/pq"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

func validateCSVPath(csvPath string) {
	file, err := os.Open(csvPath)
	if err != nil {
		log.Panicf("failed reading file: %s", err)
	}
	defer file.Close()
}

// seedCmd represents the seed command
var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "A brief description of your command",
	Args: func(cmd *cobra.Command, args []string) error {
		databaseURL := viper.GetString("databaseURL")
		csvPath := viper.GetString("csvPath")
		tableName := viper.GetString("tableName")

		if csvPath == "" {
			return errors.New("csv file path must be specified in config.yml")
		}

		if databaseURL == "" {
			return errors.New("Database url must be specified in config.yml")
		}

		if tableName == "" {
			return errors.New("Table name must be specified in config.yml")
		}

		validateCSVPath(csvPath)

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		databaseURL := viper.GetString("databaseURL")
		csvPath := viper.GetString("csvPath")
		tableName := viper.GetString("tableName")

		fmt.Println(databaseURL + " " + csvPath + " " + tableName)

		setupDatabaseConnection(databaseURL)

		createTable(db, tableName)

		CSVReader(csvPath)
	},
}

func init() {
	rootCmd.AddCommand(seedCmd)
}

func CSVReader(csvPath string) {
	clientsFile, err := os.OpenFile(csvPath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer clientsFile.Close()

	csvData := []*CSVData{}

	if err := gocsv.UnmarshalFile(clientsFile, &csvData); err != nil {
		panic(err)
	}
	for _, client := range csvData {
		fmt.Println("____________", client.SN)
		insertIntoTable(db, client)
	}

	if _, err := clientsFile.Seek(0, 0); err != nil { // Go to the start of the file
		panic(err)
	}
}

func createTable(db *sql.DB, tableName string) {
	var query = fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS "%s" (
            "SN" INTEGER NOT NULL,
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

	if err != nil {
		log.Fatal(err)
	}
}
