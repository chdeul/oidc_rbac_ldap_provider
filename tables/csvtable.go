package tables

import (
	"encoding/csv"
	"io"
	"log"
	"main/tools"
	"os"
	"strings"
)

func readCsv(data any) [][]string {

	var stream io.Reader
	if filename, ok := data.(string); ok {

		file, err := os.Open(tools.CfgFilename(filename))
		if err != nil {
			log.Println(err)
			return nil
		}
		defer file.Close()
		stream = file

	} else if stream, ok = data.(io.Reader); ok {
		// already

	} else {
		tools.Fatal("Bad readCsv parameter")
		return nil
	}
	return readCsvStream(stream)
}

func writeCsv(out any, tab [][]string) error {

	var stream io.Writer
	if filename, ok := out.(string); ok {

		file, err := os.Create(tools.CfgFilename(filename))
		tools.CheckError(err, "Création projet.csv")
		defer file.Close()
		stream = file

	} else if stream, ok = out.(io.Writer); ok {
		// already

	} else {
		tools.Fatal("Bad writeCsv parameter")
		return nil
	}

	return writeCsvStream(stream, tab)
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
