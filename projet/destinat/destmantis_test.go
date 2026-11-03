package destinat

import (
	"main/tools"
	"testing"
)

/*
Prerequis Mantis    accessible avec un Projet "Mon_Projet"
*/

func TestDestMantis(t *testing.T) {

	de, err := NewDestMantis()
	tools.TEcheck(t, err, "NewDestMantis")
	defer de.Close()

	testDestination(t, de, "develper", "manager")
}
