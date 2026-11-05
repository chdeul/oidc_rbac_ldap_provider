package referent

import (
	"main/tools"
	"os"
	"testing"
)

// Commun package
func TestMain(m *testing.M) {
	tools.LoadConfig("data/test/config.yml")
	os.Exit(m.Run())
}

// prerequis : user admin existe avec les bons attributs (dans mantis,ldap)

// Common test Function  for all destination
func testReferentiel(t *testing.T, rf Referentiel) {

	tools.Recover()

	tools.Tcheck(t, rf, "NewRefXxxxx")

	// Recherche un utilisateur via son identifiant, ou une partie
	// Retourne une table d'utilisateurs trouvés (un seul en général)
	users := rf.SearchUser("admin")
	tools.Tcheck(t, len(users) == 1, "SearchUser nb")
	tools.Tcheck(t, users[0].Identifiant == "admin", "SearchUser ident")
	tools.Tcheck(t, users[0].NomComplet == "Super Administrateur", "SearchUser name")
	tools.Tcheck(t, users[0].AdresseMail == "root@localhost.loc", "SearchUser mail")
	tools.Tcheck(t, users[0].Roles == "Administrateur", "SearchUser roles")
	tools.Tcheck(t, users[0].MotDePasse != "", "SearchUser password")

	tools.Tcheck(t, users[0].Uid != "", "SearchUser uid")

	// Retourne un groupe d'utilisateurs trouvés
	users = rf.SearchGroup("Edition")
	tools.Tcheck(t, len(users) == 2, "SearchGroup")
}
