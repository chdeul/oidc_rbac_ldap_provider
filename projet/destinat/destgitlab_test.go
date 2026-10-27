package destinat

import (
	obj "main/objets"
	"main/tools"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	d := os.Getenv("WORKSPACE_DIR")
	os.Chdir(d)
	tools.LoggingLevel = tools.DebugLevel
	os.Exit(m.Run())
}

func TestGitlab(t *testing.T) {

	de, err := NewDestGitlab()
	if err != nil {
		t.Fatalf("NewDestGitlab failed: %v", err)
	}
	defer de.Close()

	// users methods
	users := de.SearchUsers("utilisateur") // 3 letters min, no pattern
	if len(users) == 0 {
		t.Fatal("SearchUsers is nil")
	}
	tools.Debug("SearchUsers nb %d", len(users))

	cuser := obj.User{Identifiant: "xxxxx", NomComplet: "X Xxxx", AdresseMail: "xx@xxx.fr"}
	err = de.CreateUser(&cuser)
	if err != nil {
		tools.Fatal("CreateUser failed")
	}
	tools.Debug("CreateUser id %s", cuser.Uid)

	// on ne peut pas changer l'email sous gitlab !
	cuser.NomComplet = "Z XYZ"
	err = de.UpdateUser(&cuser)
	if err != nil {
		tools.Debug("UpdateUser failed")
	}
	tools.Debug("UpdateUser email %s", cuser.AdresseMail)

	err = de.DeleteUser(&cuser)
	if err != nil {
		tools.Debug("DeleteUser failed")
	}
	tools.Debug("DeleteUser id %s", cuser.Uid)

	// projects methods
	prjs := de.SearchProjects("Proj")
	if len(prjs) == 0 {
		t.Fatal("SearchProjects is nil")
	}
	tools.Debug("SearchProjects nb %d", len(prjs))

	// roles methods
	err = de.CreateRole(prjs[0], users[0], "Guest")
	if err != nil {
		tools.Debug("CreateRole failed")
	}
	tools.Debug("CreateRole")

	err = de.UpdateRole(prjs[0], users[0], "Owner")
	if err != nil {
		tools.Debug("UpdateRole failed")
	}
	tools.Debug("UpdateRole")

	role := de.GetRole(prjs[0], users[0])
	if role == "" {
		t.Error("GetRole is empty")
	}
	tools.Debug("GetRole %s", role)

	err = de.DeleteRole(prjs[0], users[0])
	if err != nil {
		tools.Debug("DeleteRole failed")
	}
	tools.Debug("DeleteRole")
}
