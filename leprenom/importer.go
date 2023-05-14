package leprenom

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

type CsvContent struct {
	Gender    int
	FirstName string
	Year      int
	Count     int
}

func fromRecord(record []string) (CsvContent, error) {
	if len(record) != 4 {
		return CsvContent{}, fmt.Errorf("Invalid size of record %v %d, expected 4", record, len(record))
	}
	gender, err := strconv.Atoi(record[0])
	if err != nil {
		return CsvContent{}, err
	}
	year := 0
	if record[2] != "XXXX" {
		year, err = strconv.Atoi(record[2])
		if err != nil {
			return CsvContent{}, err
		}
	}
	count, err := strconv.Atoi(record[3])
	if err != nil {
		return CsvContent{}, err
	}
	return CsvContent{Gender: gender, FirstName: record[1], Year: year, Count: count}, nil
}

func Import(r io.Reader) ([]CsvContent, error) {
	reader := csv.NewReader(r)
	reader.Comma = ';'
	records, err := reader.ReadAll()
	if err != nil {
		return []CsvContent{}, err
	}
	content := make([]CsvContent, 0, len(records))
	for _, record := range records {
		// Skip header
		if record[0] == "sexe" {
			continue
		}
		csvContent, err := fromRecord(record)
		if err != nil {
			return []CsvContent{}, err
		}
		if csvContent.FirstName == "_PRENOMS_RARES" {
			continue
		}
		content = append(content, csvContent)
	}
	return content, nil
}
