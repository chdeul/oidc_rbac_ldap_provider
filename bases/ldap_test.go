package bases

import (
	obj "main/objets"
	"main/tools"
	"testing"
)

// Testé avec OpenLDAP en container dinkel/openldap,  vide
// --env SLAPD_DOMAIN=ldap.example.org
// --env SLAPD_PASSWORD=admin123
// --publish 389:389

func check(t *testing.T, val any, args ...any) {
	if !tools.Assert(val, args) {
		t.Fail()
	}
}
func echeck(t *testing.T, val error, args ...any) {
	if !tools.AssertError(val, args) {
		t.Fail()
	}
}
func TestLdap(t *testing.T) {
	var l any
	var e error

	l, e = LdapOpen("ldap://localhost:389", "dc=ldap,dc=example,dc=org", "admin123")
	check(t, l)
	echeck(t, e)

	e = LdapAddUser(l, &obj.User{Uid: "111", Identifiant: "utest", NomComplet: "Un Test", AdresseMail: "un.test@localhost.org"})
	echeck(t, e)

	e = LdapSetPassword(l, "utest", "AbcdEf=$123")
	echeck(t, e)

	LdapClose(l)
	l, e = LdapOpen("ldap://localhost:389", "dc=ldap,dc=example,dc=org", "")
	check(t, l)
	echeck(t, e)

	u := LdapGetUser(l, "utest")
	check(t, u != nil && u.NomComplet == "Un Test" && u.AdresseMail == "un.test@localhost.org")

	ok := LdapCheckPassword("ldap://localhost:389", "dc=ldap,dc=example,dc=org", "utest", "AbcdEf=$123")
	check(t, ok)

	LdapClose(l)
	l, e = LdapOpen("ldap://localhost:389", "dc=ldap,dc=example,dc=org", "admin123")
	check(t, l)
	echeck(t, e)

	e = LdapUpdateUser(l, &obj.User{Uid: "111", Identifiant: "utest", NomComplet: "Un-Test", AdresseMail: "un.test@localhost.org"})
	echeck(t, e)

	u = LdapGetUser(l, "utest")
	check(t, u != nil && u.NomComplet == "Un-Test" && u.AdresseMail == "un.test@localhost.org")

	e = LdapAddGroup(l, "UnGroupe", "description...")
	echeck(t, e)

	e = LdapAddMember(l, "UnGroupe", "utest")
	echeck(t, e)

	m := LdapGetGroup(l, "UnGroupe")
	check(t, len(m) == 2 && m[1] == "utest" && m[0] == "admin")

	e = LdapDeleteMember(l, "UnGroupe", "utest")
	echeck(t, e)

	m = LdapGetGroup(l, "UnGroupe")
	check(t, len(m) == 1 && m[0] == "admin")

	e = LdapDeleteGroup(l, "UnGroupe")
	echeck(t, e)

	m = LdapGetGroup(l, "UnGroupe")
	check(t, len(m) == 0)

	e = LdapDeleteUser(l, "utest")
	echeck(t, e)

	u = LdapGetUser(l, "utest")
	check(t, u == nil)

	LdapClose(l)
}

func TestClean(t *testing.T) {
	l, _ := LdapOpen("ldap://localhost:389", "dc=ldap,dc=example,dc=org", "admin123")
	check(t, l)
	if l != nil {
		LdapDeleteGroup(l, "UnGroupe")
		LdapDeleteUser(l, "utest")
		LdapClose(l)
	}
}
