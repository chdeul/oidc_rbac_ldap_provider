package destinat

import (
	"main/tools"
	"testing"
)

/*
Prerequis Gitlab    accessible avec un Projet "Mon_Projet"
*/

func TestGitlab(t *testing.T) {

	de, err := NewDestGitlab()
	tools.TEcheck(t, err, "NewDestGitlab")
	defer de.Close()

	testDestination(t, de, "developer", "owner")
}
