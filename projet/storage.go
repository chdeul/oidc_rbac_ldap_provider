package projet

import (
	obj "main/objets"
	"main/tables"
	"main/tools"
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
