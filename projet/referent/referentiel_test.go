package referent

import (
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
func testReferentiel(t *testing.T, rf Referentiel) {
	tools.Tcheck(t, rf, "NewRefXxxxx")

	// Recherche un utilisateur via son identifiant, ou une partie
	// Retourne une table d'utilisateurs trouvés (un seul en général)
	users := rf.SearchUser("admin")
	tools.Tcheck(t, len(users) == 1, "SearchUser")
	tools.Tcheck(t, users[0].Identifiant == "admin", "SearchUser")
	tools.Tcheck(t, users[0].NomComplet == "Super Administrateur", "SearchUser")
	tools.Tcheck(t, users[0].AdresseMail == "root@local.loc", "SearchUser")
	tools.Tcheck(t, users[0].Roles == "Administrateur", "SearchUser")
	tools.Tcheck(t, users[0].Uid == "1", "SearchUser")
	tools.Tcheck(t, users[0].MotDePasse == "chdeul", "SearchUser")

	// Retourne un groupe d'utilisateurs trouvés
	users = rf.SearchGroup("Edition")
	tools.Tcheck(t, len(users) == 2, "SearchUser")
}
