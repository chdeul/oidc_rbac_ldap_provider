package projet

import (
	"main/tables"
	"main/tools"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	tools.LoadConfig("data/test/config.yml")
	os.Exit(m.Run())
}

func lectureFileProjet(filename string) {
	tab := tables.ReadRights(filename)
	LectureProjet(tab)
}

func TestLdap(t *testing.T) {

	// Nettoyer
	//ldapdelete -H 'ldap://localhost:389' -x -D 'cn=admin,dc=ldap,dc=example,dc=org' -w 'admin123' -r 'ou=users,dc=ldap,dc=example,dc=org'
	//ldapdelete -H 'ldap://localhost:389' -x -D 'cn=admin,dc=ldap,dc=example,dc=org' -w 'admin123' -r 'ou=groups,dc=ldap,dc=example,dc=org'

	// Creer
	lectureFileProjet("data/test/prjldap.csv")
}

func TestMantis(t *testing.T) {

	// Nettoyer Mantis, créer le projet Mon_Projet

	lectureFileProjet("data/test/prjmantis.csv")
}

func TestGitlab(t *testing.T) {

	// Nettoyer Gitlab, créer le projet Mon_Projet

	lectureFileProjet("data/test/prjgitlab.csv")
}
