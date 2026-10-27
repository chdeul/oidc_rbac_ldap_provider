package tables

import (
	"encoding/csv"
	"log"
	"os"
)

func readCsv(filename string) [][]string {

	file, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return nil
	}
	reader := csv.NewReader(file)
	reader.Comma = ';'
	reader.Comment = '#'
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true
	tab, err := reader.ReadAll()
	if err != nil {
		log.Println(err)
		return nil
	}
	file.Close()
	return tab
}

func writeCsv(filename string, tab [][]string) error {

	file, err := os.Create(filename)
	if err != nil {
		log.Println(err)
	}
	writer := csv.NewWriter(file)
	writer.Comma = ';'
	err = writer.WriteAll(tab)
	if err != nil {
		log.Println(err)
	}
	return nil
}
