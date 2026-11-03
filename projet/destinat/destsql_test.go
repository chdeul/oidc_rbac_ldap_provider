package destinat

import (
	"main/tools"
	"testing"
)

/*
Prerequis Base SQL mantis  accessible avec un Projet "Mon_Projet"
*/

func TestDestSql(t *testing.T) {

	de, err := NewDestSql()
	tools.TEcheck(t, err, "NewDestSql")
	defer de.Close()

	testDestination(t, de, "developer", "manager")
}
