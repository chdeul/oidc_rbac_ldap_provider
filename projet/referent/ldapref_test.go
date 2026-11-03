package referent

import (
	"main/tools"
	"testing"
)

/*
Prerequis Openldap    accessible
*/

func TestRefLdap(t *testing.T) {

	rf, err := NewRefLdap()
	tools.TEcheck(t, err, "NewRefLdap")
	defer rf.Close()

	testReferentiel(t, rf)
}
