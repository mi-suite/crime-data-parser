package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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
