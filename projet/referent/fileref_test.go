package referent

import (
	"main/tools"
	"testing"
)

func TestRefFile(t *testing.T) {

	rf, err := NewRefFile()
	tools.TEcheck(t, err, "NewRefFile")
	defer rf.Close()

	testReferentiel(t, rf)
}
