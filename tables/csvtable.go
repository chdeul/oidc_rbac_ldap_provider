package tables

import (
	"encoding/csv"
	"io"
	"log"
	"main/tools"
	"os"
	"strings"
)

func readCsv(filename string) [][]string {
	file, err := os.Open(tools.CfgFilename(filename))
	if err != nil {
		log.Println(err)
		return nil
	}
	defer file.Close()

	return readCsvStream(file)
}

func writeCsv(filename string, tab [][]string) error {

	file, err := os.Create(tools.CfgFilename(filename))
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	return writeCsvStream(file, tab)
}

func readCsvStream(stream io.Reader) [][]string {

	reader := csv.NewReader(stream)
	reader.Comma = ','
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

func writeCsvStream(stream io.Writer, tab [][]string) error {

	writer := csv.NewWriter(stream)
	writer.Comma = ','
	err := writer.WriteAll(tab)
	if err != nil {
		log.Println(err)
	}
	return nil
}
