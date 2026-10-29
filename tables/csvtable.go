package tables

import (
	"encoding/csv"
	"log"
	"os"
	"strings"
)

func readCsv(filename string) [][]string {

	file, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer file.Close()

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
	// TrimEndSpace
	for i, l := range tab {
		for j, v := range l {
			s := v
			var b bool
			for {
				s, b = strings.CutSuffix(s, " ")
				if !b {
					tab[i][j] = s
					break
				}
			}
		}
	}
	return tab
}

func writeCsv(filename string, tab [][]string) error {

	file, err := os.Create(filename)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	err = writer.WriteAll(tab)
	if err != nil {
		log.Println(err)
	}
	return nil
}
