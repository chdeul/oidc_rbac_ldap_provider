package projet

import (
	"main/tools"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	tools.LoadConfig("../data/test")
	os.Exit(m.Run())
}

func TestLdap(t *testing.T) {

	// Nettoyer
	//ldapdelete -H 'ldap://localhost:389' -x -D 'cn=admin,dc=ldap,dc=example,dc=org' -w 'admin123' -r 'ou=users,dc=ldap,dc=example,dc=org'
	//ldapdelete -H 'ldap://localhost:389' -x -D 'cn=admin,dc=ldap,dc=example,dc=org' -w 'admin123' -r 'ou=groups,dc=ldap,dc=example,dc=org'

	// Creer
	LectureProjet("../data/test/prjldap.csv")

}
func TestMantis(t *testing.T) {
	LectureProjet("../data/test/prjmantis.csv")
}
func TestGitlab(t *testing.T) {
	LectureProjet("../data/test/prjgitlab.csv")
}
