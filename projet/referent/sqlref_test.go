package referent

import (
	"main/tools"
	"testing"
)

/*
Prerequis Base sql mantis    accessible
*/

func TestRefSql(t *testing.T) {

	rf, err := NewRefSql()
	tools.TEcheck(t, err, "NewRefSql")
	defer rf.Close()

	testReferentiel(t, rf)
}
