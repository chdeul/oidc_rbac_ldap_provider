package bases

import (
	obj "main/objets"
	"testing"
)

// Testé avec OpenLDAP en container dinkel/openldap,  vide
// --env SLAPD_DOMAIN=ldap.example.org
// --env SLAPD_PASSWORD=admin123
// --publish 389:389

func TestLdap(t *testing.T) {

	t.Run("step1", testLdap)

	t.Run("clean", cleanLdap)
}

func testLdap(t *testing.T) {
	var l any
	var e error

	l, e = LdapOpen("ldap://localhost:389", "dc=ldap,dc=example,dc=org", "admin123")
	if e != nil {
		t.Fatal("-- LdapOpen failed")
	}
	e = LdapAddUser(l, &obj.User{Uid: "1", Identifiant: "utest", NomComplet: "Un Test", AdresseMail: "un.test@localhost.org"})
	if e != nil {
		t.Error("-- LdapAddUser failed")
	}
	e = LdapSetPassword(l, "utest", "AbcdEf=$123")
	if e != nil {
		t.Error("-- LdapSetPassword failed")
	}
	LdapClose(l)
	l, e = LdapOpen("ldap://localhost:389", "dc=ldap,dc=example,dc=org", "")
	if e != nil {
		t.Fatal("-- LdapOpen failed")
	}
	u := LdapGetUser(l, "utest")
	if u == nil || u.NomComplet != "Un Test" || u.AdresseMail != "un.test@localhost.org" {
		t.Error("-- LdapGetUser failed")
	}
	ok := LdapCheckPassword("ldap://localhost:389", "dc=ldap,dc=example,dc=org", "utest", "AbcdEf=$123")
	if !ok {
		t.Error("-- LdapCheckPassword failed")
	}
	LdapClose(l)
	l, e = LdapOpen("ldap://localhost:389", "dc=ldap,dc=example,dc=org", "admin123")
	if e != nil {
		t.Fatal("-- LdapOpen failed")
	}
	e = LdapUpdateUser(l, &obj.User{Uid: "1", Identifiant: "utest", NomComplet: "Un-Test", AdresseMail: "un.test@localhost.org"})
	if e != nil {
		t.Error("-- LdapUpdateUser failed")
	}
	u = LdapGetUser(l, "utest")
	if u == nil || u.NomComplet != "Un-Test" || u.AdresseMail != "un.test@localhost.org" {
		t.Error("-- LdapUpdateUser failed")
	}
	e = LdapAddGroup(l, "UnGroupe", "description...")
	if e != nil {
		t.Error("-- LdapAddGroup failed")
	}
	e = LdapAddMember(l, "UnGroupe", "utest")
	if e != nil {
		t.Error("-- LdapAddMember failed")
	}
	m := LdapGetGroup(l, "UnGroupe")
	if len(m) != 2 || m[1] != "utest" || m[0] != "admin" {
		t.Error("-- LdapGetGroup failed")
	}
	e = LdapDeleteMember(l, "UnGroupe", "utest")
	if e != nil {
		t.Error("-- LdapDeleteMember failed")
	}
	m = LdapGetGroup(l, "UnGroupe")
	if len(m) != 1 || m[0] != "admin" {
		t.Error("-- LdapDeleteMember failed")
	}
	e = LdapDeleteGroup(l, "UnGroupe")
	if e != nil {
		t.Error("-- LdapDeleteGroup failed")
	}
	m = LdapGetGroup(l, "UnGroupe")
	if len(m) != 0 {
		t.Error("-- LdapDeleteGroup failed")
	}
	e = LdapDeleteUser(l, "utest")
	if e != nil {
		t.Error("-- LdapLdapDeleteUserAddGroup failed")
	}
	u = LdapGetUser(l, "utest")
	if u != nil {
		t.Error("-- LdapDeleteUser failed")
	}
	LdapClose(l)
}

func cleanLdap(t *testing.T) {
	l, _ := LdapOpen("ldap://localhost:389", "dc=ldap,dc=example,dc=org", "admin123")
	if l != nil {
		LdapDeleteGroup(l, "UnGroupe")
		LdapDeleteUser(l, "utest")
		LdapClose(l)
	}
}
