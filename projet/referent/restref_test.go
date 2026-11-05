package referent

import (
	"main/tools"
	"testing"
)

/*
Prerequis api rest mantis    accessible
*/

func TestRefRest(t *testing.T) {

	tools.Cfg.Ref.Type = "rest"
	tools.Cfg.Ref.Path = "data/test/restApi.yml"

	rf, err := NewRefRest()
	tools.TEcheck(t, err, "NewRefRest")
	defer rf.Close()

	testReferentiel(t, rf)
}
