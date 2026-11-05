package referent

import (
	"main/tools"
	"testing"
)

/*
Prerequis Openldap    accessible
*/

func TestRefLdap(t *testing.T) {

	tools.Cfg.Ref.Type = "ldap"
	tools.Cfg.Ref.Path = "data/test/ldapData.yml"

	rf, err := NewRefLdap()
	tools.TEcheck(t, err, "NewRefLdap")
	defer rf.Close()

	testReferentiel(t, rf)
}
