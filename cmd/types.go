package cmd

import "time"

type DateTime struct {
	time.Time
}

// UnmarshalCSV Convert the CSV string as internal date
func (date *DateTime) UnmarshalCSV(csv string) (err error) {
	date.Time, err = time.Parse("01/02/2006 15:04:05 AM", csv)
	if err != nil {
		date.Time, err = time.Parse("01/02/2006 15:04:05 PM", csv)
	}

	return err
}

type CSVData struct {
	SN                  int      `csv:"SN"`
	ID                  int      `csv:"ID"`
	CaseNumber          string   `csv:"Case Number"`
	Date                DateTime `csv:"Date"`
	Block               string   `csv:"Block"`
	IUCR                string   `csv:"IUCR"`
	PrimaryType         string   `csv:"Primary Type"`
	Description         string   `csv:"Description"`
	LocationDescription string   `csv:"Location Description"`
	Arrest              bool     `csv:"Arrest"`
	Domestic            bool     `csv:"Domestic"`
	Beat                int      `csv:"Beat"`
	District            int      `csv:"District"`
	Ward                int      `csv:"Ward"`
	CommunityArea       int      `csv:"Community Area"`
	FBICode             string   `csv:"FBI Code"`
	XCoordinate         int      `csv:"X Coordinate"`
	YCoordinate         int      `csv:"Y Coordinate"`
	Year                int      `csv:"Year"`
	UpdatedOn           DateTime `csv:"Updated On"`
	Latitude            int      `csv:"Latitude"`
	Longitude           int      `csv:"Longitude"`
	Location            string   `csv:"Location"`
}

type ScanClient struct {
	SN                  int       `csv:"SN"`
	ID                  int       `csv:"ID"`
	CaseNumber          string    `csv:"Case Number"`
	Date                time.Time `csv:"Date"`
	Block               string    `csv:"Block"`
	IUCR                string    `csv:"IUCR"`
	PrimaryType         string    `csv:"Primary Type"`
	Description         string    `csv:"Description"`
	LocationDescription string    `csv:"Location Description"`
	Arrest              bool      `csv:"Arrest"`
	Domestic            bool      `csv:"Domestic"`
	Beat                int       `csv:"Beat"`
	District            int       `csv:"District"`
	Ward                int       `csv:"Ward"`
	CommunityArea       int       `csv:"Community Area"`
	FBICode             string    `csv:"FBI Code"`
	XCoordinate         int       `csv:"X Coordinate"`
	YCoordinate         int       `csv:"Y Coordinate"`
	Year                int       `csv:"Year"`
	UpdatedOn           time.Time `csv:"Updated On"`
	Latitude            int       `csv:"Latitude"`
	Longitude           int       `csv:"Longitude"`
	Location            string    `csv:"Location"`
}
