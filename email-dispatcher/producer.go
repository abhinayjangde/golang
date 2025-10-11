package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

// private function
func loadRecipients(filePath string) error {
	f, err := os.Open(filePath)

	if err != nil {
		return err
	}

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()

	if err != nil {
		return err
	}

	for _, record := range records[1:] { // Skip header row
		fmt.Println(record)
		// recipient := Recipient{
		// 	Name:  record[0],
		// 	Email: record[1],
		// }
		// // Here you would typically send the recipient to a channel or process it
		// _ = recipient // Placeholder to avoid unused variable error
	}

	return nil
}
