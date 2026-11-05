package referent

import (
	"main/tools"
	"testing"
)

/*
Prerequis Base sql mantis    accessible
*/

func TestRefSql(t *testing.T) {

	tools.Cfg.Ref.Type = "sql"
	tools.Cfg.Ref.Path = "data/test/sqlData.yml"

	rf, err := NewRefSql()
	tools.TEcheck(t, err, "NewRefSql")
	defer rf.Close()

	testReferentiel(t, rf)
}
