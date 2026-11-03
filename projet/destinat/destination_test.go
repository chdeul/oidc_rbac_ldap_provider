package destinat

import (
	obj "main/objets"
	"main/tools"
	"os"
	"testing"
)

// Commun package
func TestMain(m *testing.M) {
	tools.LoadConfig("../../data/test")
	os.Exit(m.Run())
}

// Common test Function  for all destination
func testDestination(t *testing.T, de Destination, role1, role2 string) {
	tools.Tcheck(t, de, "NewDestXxxxx")

	// users methods
	cuser := obj.User{Uid: "99", Identifiant: "xxxxx", NomComplet: "X Xxxx", AdresseMail: "xx@xxx.fr"}
	err := de.CreateUser(&cuser)
	tools.TEcheck(t, err, "create user")

	users := de.SearchUsers("xxxxx") // 3 letters min, no pattern
	tools.Tcheck(t, len(users) == 1, "search user")
	tools.Tcheck(t, users[0].AdresseMail == "xx@xxx.fr", "search user")

	cuser.AdresseMail = "zz@xyz.fr"
	cuser.NomComplet = "Z XYZ"
	err = de.UpdateUser(&cuser)
	tools.TEcheck(t, err, "update user")

	users = de.SearchUsers("xxxxx") // 3 letters min, no pattern
	tools.Tcheck(t, len(users) == 1, "search user")
	tools.Tcheck(t, users[0].AdresseMail == "zz@xyz.fr", "update user")
	tools.Tcheck(t, users[0].NomComplet == "Z XYZ", "update user")

	// projects methods
	prjs := de.SearchProjects("Mon_Projet")
	tools.Tcheck(t, prjs, "search projects")

	// roles methods
	err = de.CreateRole(prjs[0], users[0], role1)
	tools.TEcheck(t, err, "create role")

	err = de.UpdateRole(prjs[0], users[0], role2)
	tools.TEcheck(t, err, "update role")

	role := de.GetRole(prjs[0], users[0])
	tools.Tcheck(t, role, "get role")
	tools.Tcheck(t, role == role2, "get role")

	err = de.DeleteRole(prjs[0], users[0])
	tools.TEcheck(t, err, "delete role")

	// Clean
	err = de.DeleteUser(&cuser)
	tools.TEcheck(t, err, "delete user")
}
