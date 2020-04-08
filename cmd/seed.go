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
	},
}

func init() {
	rootCmd.AddCommand(seedCmd)
}
