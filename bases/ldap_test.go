package bases

import (
	obj "main/objets"
	"main/tools"
	"os"
	"testing"
)

// Testé avec OpenLDAP en container dinkel/openldap,  vide
// --env SLAPD_DOMAIN=ldap.example.org
// --env SLAPD_PASSWORD=admin123
// --publish 389:389

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestLdap(t *testing.T) {
	var l any
	var e error

	l, e = LdapOpen("ldap://localhost:389", "dc=ldap,dc=example,dc=org", "admin123")
	tools.Tcheck(t, l, "LdapOpen")
	tools.Tcheck(t, e, "LdapOpen")

	e = LdapAddUser(l, &obj.User{Uid: "111", Identifiant: "utest", NomComplet: "Un Test", AdresseMail: "un.test@localhost.org"})
	tools.Tcheck(t, e, "LdapAddUser")

	e = LdapSetPassword(l, "utest", "AbcdEf=$123")
	tools.Tcheck(t, e, "LdapSetPassword")

	LdapClose(l)
	l, e = LdapOpen("ldap://localhost:389", "dc=ldap,dc=example,dc=org", "")
	tools.Tcheck(t, l, "LdapOpen")
	tools.Tcheck(t, e, "LdapOpen")

	u := LdapGetUser(l, "utest")
	tools.Tcheck(t, u != nil && u.NomComplet == "Un Test" && u.AdresseMail == "un.test@localhost.org", "LdapGetUser")

	ok := LdapCheckPassword("ldap://localhost:389", "dc=ldap,dc=example,dc=org", "utest", "AbcdEf=$123")
	tools.Tcheck(t, ok, "LdapCheckPassword")

	LdapClose(l)
	l, e = LdapOpen("ldap://localhost:389", "dc=ldap,dc=example,dc=org", "admin123")
	tools.Tcheck(t, l, "LdapOpen")
	tools.Tcheck(t, e, "LdapOpen")

	e = LdapUpdateUser(l, &obj.User{Uid: "111", Identifiant: "utest", NomComplet: "Un-Test", AdresseMail: "un.test@localhost.org"})
	tools.Tcheck(t, e, "LdapUpdateUser")

	u = LdapGetUser(l, "utest")
	tools.Tcheck(t, u != nil && u.NomComplet == "Un-Test" && u.AdresseMail == "un.test@localhost.org", "LdapGetUser")

	e = LdapAddGroup(l, "UnGroupe", "description...")
	tools.Tcheck(t, e, "LdapAddGroup")

	e = LdapAddMember(l, "UnGroupe", "utest")
	tools.Tcheck(t, e, "LdapAddMember")

	m := LdapGetGroup(l, "UnGroupe")
	tools.Tcheck(t, len(m) == 2 && m[1] == "utest" && m[0] == "admin", "LdapGetGroup")

	e = LdapDeleteMember(l, "UnGroupe", "utest")
	tools.Tcheck(t, e, "LdapDeleteMember")

	m = LdapGetGroup(l, "UnGroupe")
	tools.Tcheck(t, len(m) == 1 && m[0] == "admin", "LdapGetGroup")

	e = LdapDeleteGroup(l, "UnGroupe")
	tools.Tcheck(t, e, "LdapDeleteGroup")

	m = LdapGetGroup(l, "UnGroupe")
	tools.Tcheck(t, len(m) == 0, "LdapGetGroup")

	e = LdapDeleteUser(l, "utest")
	tools.Tcheck(t, e, "LdapDeleteUser")

	u = LdapGetUser(l, "utest")
	tools.Tcheck(t, u == nil, "LdapGetUser")

	LdapClose(l)
}

func TestClean(t *testing.T) {
	l, _ := LdapOpen("ldap://localhost:389", "dc=ldap,dc=example,dc=org", "admin123")
	tools.Tcheck(t, l, "LdapOpen for clean")
	if l != nil {
		LdapDeleteGroup(l, "UnGroupe")
		LdapDeleteUser(l, "utest")
		LdapClose(l)
	}
}
