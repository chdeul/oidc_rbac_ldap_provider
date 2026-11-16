package projet

import (
	obj "main/objets"
	"main/tables"
	"main/tools"
	"os"
	"strings"
)

func StoreCsvProject(users []obj.Right) {

	// TODO : nom abrégé     Nom projet  to filename
	projname := users[0].Type
	filename := tools.CfgFilename("data/projects/" + projname + ".csv")

	tables.WriteRights(users, filename)
}

func LoadCsvProject(projname string) []obj.Right {

	// TODO : nom abrégé     Nom projet  to filename
	filename := tools.CfgFilename("data/projects/" + projname + ".csv")

	return tables.ReadRights(filename)
}

// Find projects list
func ListOfCsvProject() []string {

	list := []string{}
	dir, err := os.ReadDir(tools.CfgFilename("data/projects"))
	if tools.CheckError(err, "Cannot read dir of projects") {
		for _, f := range dir {
			if !f.IsDir() && strings.HasSuffix(f.Name(), ".csv") {
				list = append(list, strings.Replace(f.Name(), ".csv", "", 1))
			}
		}
	}
	return list
}
