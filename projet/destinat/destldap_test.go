package destinat

import (
	"main/tools"
	"os"
	"testing"
)

// Commun package
func TestMain(m *testing.M) {
	tools.LoadConfig("../../data/test")
	os.Exit(m.Run())
}

/*
Prerequis Openldap    accessible
*/

func TestDestLdap(t *testing.T) {

	de, err := NewDestLdap()
	tools.TEcheck(t, err, "NewDestLdap")
	defer de.Close()

	testDestination(t, de, "groupeTestA", "groupeTestB")
}
