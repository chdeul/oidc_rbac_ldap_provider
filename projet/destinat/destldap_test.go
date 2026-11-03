package destinat

import (
	"main/tools"
	"testing"
)

/*
Prerequis Openldap    accessible
*/

func TestDestLdap(t *testing.T) {

	de, err := NewDestLdap()
	tools.TEcheck(t, err, "NewDestLdap")
	defer de.Close()

	testDestination(t, de, "groupeTestA", "groupeTestB")
}
