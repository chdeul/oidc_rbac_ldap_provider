package bases

import (
	"fmt"
	obj "main/objets"
	"main/tools"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

// OpenLDAP , schémas avec :  inetOrgPerson   et   groupOfNames
// Simple organization : a users list, a static group list,  no hierachy

var usersDn string
var groupsDn string
var adminDn string
var adminPw string

func LdapOpen(url string, bind string, passwd string) (any, error) {

	// Connection open
	l, err := ldap.DialURL(url)
	if l == nil || err != nil {
		tools.Panic("LDAP Open: %v", err)
		return l, err
	}
	// Absolute dn : "bind" is root dn of openldap
	usersDn = "ou=users," + bind
	groupsDn = "ou=groups," + bind

	// Read/Write role : cn=admin on OpenLdap
	adminDn = "cn=admin," + bind
	adminPw = passwd
	if passwd != "" {
		// Bind with the admin dn
		err = l.Bind(adminDn, adminPw)
		if err == nil {
			// Check expected schema
			err = initSchema(l, bind)
		}
	} else { //  Read only
		err = l.UnauthenticatedBind("")
	}
	if err != nil {
		tools.Panic("LDAP Bind: %v", err)
	}
	return l, err
}

func LdapClose(conn any) {

	// Cast the any parameter
	l, _ := conn.(*ldap.Conn)
	if l != nil {
		l.Close()
	}
}

func initSchema(l *ldap.Conn, bind string) error {

	// Check if Org.Unit exist
	searchRequest := ldap.NewSearchRequest(bind,
		ldap.ScopeSingleLevel, ldap.NeverDerefAliases, 0, 0, false,
		"(objectClass=organizationalUnit)", []string{"ou"}, nil)
	ou, err := l.Search(searchRequest)

	if err == nil {
		// find it
		foundUsers := false
		foundGroups := false
		for _, e := range ou.Entries {
			if e.GetAttributeValue("ou") == "users" {
				foundUsers = true
			}
			if e.GetAttributeValue("ou") == "groups" {
				foundGroups = true
			}
		}
		// Else create it
		if !foundUsers {
			tools.Info("LDAP Init Schema : Add org.unit users")
			addRequest := ldap.NewAddRequest(usersDn, nil)
			addRequest.Attribute("objectClass", []string{"organizationalUnit"})
			addRequest.Attribute("ou", []string{"users"})
			err = l.Add(addRequest)
		}
		if err == nil && !foundGroups {
			tools.Info("LDAP Init Schema : Add org.unit groups")
			addRequest := ldap.NewAddRequest(groupsDn, nil)
			addRequest.Attribute("objectClass", []string{"organizationalUnit"})
			addRequest.Attribute("ou", []string{"groups"})
			err = l.Add(addRequest)
		}
	}
	return err
}

// ------------------ Users ------------------

func LdapGetUser(conn any, identifiant string) *obj.User {

	users := LdapSearch(conn, identifiant)

	// Must be uniq
	if len(users) != 1 {
		if len(users) > 1 {
			tools.Warn("LdapGetUser: Too many User entries returned")
		}
		return nil
	}
	return users[0]
}

func LdapGetUserFromEmail(conn any, email string) *obj.User {

	// Check pattern : email sans ( ) \
	if strings.ContainsAny(email, "()\\") {
		tools.Warn("LdapGetUserFromEmail: Malformed mail address")
		return nil
	}
	// Recherche sur mail
	rech := fmt.Sprintf("(&(objectClass=inetOrgPerson)(mail=%s))", email)
	users := ldapSearchExpr(conn, rech)

	// Must be uniq
	if len(users) != 1 {
		if len(users) > 1 {
			tools.Warn("LdapGetUserFromEmail: Too many User entries returned")
		}
		return nil
	}
	return users[0]
}

func LdapSearch(conn any, pattern string) []*obj.User {

	// Check pattern : identifiant + * , sans ( ) \
	if strings.ContainsAny(pattern, "()\\") {
		tools.Warn("LdapGetUserFromEmail: Malformed search pattern")
		return nil
	}
	rech := fmt.Sprintf("(&(objectClass=inetOrgPerson)(sn=%s))", pattern)
	return ldapSearchExpr(conn, rech)
}

func ldapSearchExpr(conn any, rech string) []*obj.User {
	l, _ := conn.(*ldap.Conn)

	// Filter on objets of users,and this identifiant pattern
	searchRequest := ldap.NewSearchRequest(usersDn,
		ldap.ScopeChildren, ldap.NeverDerefAliases, 0, 0, false,
		rech, []string{"uid", "sn", "cn", "mail", "title"}, nil)

	// Search for the matching usernames
	sr, err := l.Search(searchRequest)
	if err != nil {
		tools.Error("ldapSearchExpr: error %v", err)
		return nil
	}

	// Get tables
	ret := make([]*obj.User, len(sr.Entries))
	for i, entry := range sr.Entries {
		// Get essentials attributes
		user := obj.User{Uid: entry.GetAttributeValue("uid"),
			Identifiant: entry.GetAttributeValue("sn"),
			AdresseMail: entry.GetAttributeValue("mail"),
			Roles:       entry.GetAttributeValue("title"),
			NomComplet:  entry.GetAttributeValue("cn")}
		ret[i] = &user
	}
	// User build
	return ret
}

func LdapSetPassword(conn any, identifiant string, pass string) error {
	l, _ := conn.(*ldap.Conn)

	// Modification du mot de passe avec encryptage
	dn := "sn=" + identifiant + "," + usersDn
	passwordRequest := ldap.NewPasswordModifyRequest(dn, "", pass)

	// Set Password
	_, err := l.PasswordModify(passwordRequest)
	if err != nil {
		tools.Error("LdapSetPassword: error %v", err)
	}
	return err
}

func LdapCheckPassword(url, bind, identifiant, passwd string) bool {

	// Connection open
	l, err := ldap.DialURL(url)
	if l == nil || err != nil {
		tools.Debug("LdapCheckPassword: error  %v", err)
		return false
	}
	defer l.Close()

	// Bind as the user to verify their password
	usersDn = "ou=users," + bind
	userdn := "sn=" + identifiant + "," + usersDn
	err = l.Bind(userdn, passwd)
	return err == nil
}

func LdapAddUser(conn any, user *obj.User) error {
	l, _ := conn.(*ldap.Conn)

	// Ajout des attributs d'un user
	userdn := "sn=" + user.Identifiant + "," + usersDn
	addRequest := ldap.NewAddRequest(userdn, nil)
	addRequest.Attribute("objectClass", []string{"inetOrgPerson"})
	addRequest.Attribute("uid", []string{user.Uid})
	addRequest.Attribute("sn", []string{user.Identifiant})
	addRequest.Attribute("cn", []string{user.NomComplet})
	addRequest.Attribute("mail", []string{user.AdresseMail})
	if user.Roles != "" {
		addRequest.Attribute("title", []string{user.Roles})
	}

	// Add user
	err := l.Add(addRequest)
	if err != nil {
		tools.Error("LdapAdduser: error %v", err)
	} else if user.MotDePasse != "" {
		err = LdapSetPassword(conn, user.Identifiant, user.MotDePasse)
	}
	return err
}

func LdapUpdateUser(conn any, user *obj.User) error {
	l, _ := conn.(*ldap.Conn)

	// Création de la requête de modification
	dn := "sn=" + user.Identifiant + "," + usersDn
	modifyRequest := ldap.NewModifyRequest(dn, nil)
	modifyRequest.Replace("cn", []string{user.NomComplet})
	modifyRequest.Replace("mail", []string{user.AdresseMail})
	if user.Roles != "" {
		modifyRequest.Replace("title", []string{user.Roles})
	} else {
		// TODO reset title ?
	}
	// Modify attributes
	err := l.Modify(modifyRequest)
	if err != nil {
		tools.Error("LdapUpdateUser: error %v", err)
	}
	return err
}

func LdapSearchGroups(conn any, pattern string) []*obj.Project {
	l, _ := conn.(*ldap.Conn)

	// Check pattern : nom + * , sans ( ) \
	if strings.ContainsAny(pattern, "()\\") {
		return nil
	}
	// Filter on objets of groups,and this common name
	searchRequest := ldap.NewSearchRequest(groupsDn,
		ldap.ScopeSingleLevel, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=groupOfNames)(cn=%s))", pattern),
		[]string{"cn", "description"},
		nil,
	)
	// Search for the given group
	sr, err := l.Search(searchRequest)
	if err != nil {
		tools.Error("LdapSearchGroups: error %v", err)
		return nil
	}
	// Transform to list of names
	ret := make([]*obj.Project, len(sr.Entries))
	for i, gr := range sr.Entries {
		nom := gr.GetAttributeValue("cn")
		ret[i] = obj.NewProject("", nom, nom, gr.GetAttributeValue("description"))
	}
	// Group cn list
	return ret
}

func LdapGetGroup(conn any, groupe string) []string {
	l, _ := conn.(*ldap.Conn)

	// Filter on objets of groups,and this common name
	searchRequest := ldap.NewSearchRequest(groupsDn,
		ldap.ScopeSingleLevel, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=groupOfNames)(cn=%s))", groupe),
		[]string{"member"},
		nil,
	)
	// Search for the given group
	sr, err := l.Search(searchRequest)
	if err != nil {
		tools.Error("LdapGetGroup: error %v", err)
		return nil
	}

	// Must be uniq
	if len(sr.Entries) != 1 {
		if len(sr.Entries) > 1 {
			tools.Warn("LdapGetGroup: Too many User entries returned")
		}
		return nil
	}
	// Get essentials attributes
	attr := sr.Entries[0]
	members := attr.GetAttributeValues("member")

	// Transform to members sn
	for i, m := range members {
		uid := strings.Replace(m, "sn=", "", 1)
		uid = strings.Replace(uid, ","+usersDn, "", 1)
		uid = strings.Replace(uid, adminDn, "admin", 1)
		members[i] = uid
	}
	// Members sn list
	return members
}

func LdapAddGroup(conn any, groupe, description string) error {
	l, _ := conn.(*ldap.Conn)

	// Create LDAP Add Request
	dn := "cn=" + groupe + "," + groupsDn
	addRequest := ldap.NewAddRequest(dn, nil)
	addRequest.Attribute("objectClass", []string{"groupOfNames"})
	addRequest.Attribute("cn", []string{groupe})
	addRequest.Attribute("description", []string{description})

	// Ajoute le membre "admin" par défaut, même si n'existe pas
	defmember := "sn=admin," + usersDn
	addRequest.Attribute("member", []string{defmember})

	// Execute Add
	err := l.Add(addRequest)
	if err != nil {
		tools.Error("LdapAddGroup: error %v", err)
	}
	return err
}

func LdapUpdateGroup(conn any, groupe, description string) error {
	l, _ := conn.(*ldap.Conn)

	// Création de la requête de modification
	dn := "cn=" + groupe + "," + groupsDn
	modifyRequest := ldap.NewModifyRequest(dn, nil)
	modifyRequest.Replace("description", []string{description})

	// Modify attributes
	err := l.Modify(modifyRequest)
	if err != nil {
		tools.Error("LdapUpdateGroup: error %v", err)
	}
	return err
}

func LdapAddMember(conn any, groupe string, identifiant string) error {
	l, _ := conn.(*ldap.Conn)

	dn := "cn=" + groupe + "," + groupsDn
	mb := "sn=" + identifiant + "," + usersDn

	// Modify request
	modifyRequest := ldap.NewModifyRequest(dn, nil)
	modifyRequest.Add("member", []string{mb})

	// Execute le modify member
	err := l.Modify(modifyRequest)
	if err != nil {
		tools.Error("LdapAddMember: error %v", err)
	}
	return err
}
func LdapDeleteMember(conn any, groupe string, identifiant string) error {
	l, _ := conn.(*ldap.Conn)

	if identifiant == "admin" {
		return fmt.Errorf("L'utilisateur admin est toujours membre des groupes")
	}
	// Modify request
	dn := "cn=" + groupe + "," + groupsDn
	mb := "sn=" + identifiant + "," + usersDn
	modifyRequest := ldap.NewModifyRequest(dn, nil)
	modifyRequest.Delete("member", []string{mb})

	// Execute le modify member
	err := l.Modify(modifyRequest)
	if err != nil {
		tools.Error("LdapDeleteMember: error %v", err)
	}
	return nil
}

func LdapDeleteGroup(conn any, groupe string) error {
	l, _ := conn.(*ldap.Conn)

	// Delete request
	dn := "cn=" + groupe + "," + groupsDn
	deleteRequest := ldap.NewDelRequest(dn, nil)

	// Execute le delete
	err := l.Del(deleteRequest)
	if err != nil {
		tools.Error("LdapDeleteGroup: error %v", err)
	}
	return err
}

func LdapDeleteUser(conn any, identifiant string) error {
	l, _ := conn.(*ldap.Conn)

	// Delete request
	dn := "sn=" + identifiant + "," + usersDn
	deleteRequest := ldap.NewDelRequest(dn, nil)

	// Execute le delete
	err := l.Del(deleteRequest)
	if err != nil {
		tools.Error("LdapDeleteUser: error %v", err)
	}
	// Delete this member of groups
	groups := LdapGroupsWithMember(conn, identifiant)
	for _, grp := range groups {
		err2 := LdapDeleteMember(conn, grp, identifiant)
		if err2 != nil {
			tools.Error("LdapDeleteUser: error %v", err2)
		}
	}
	return err
}

func LdapGroupsWithMember(conn any, identifiant string) []string {
	l, _ := conn.(*ldap.Conn)

	// Search request for all groups containing the user
	searchRequest := ldap.NewSearchRequest(groupsDn,
		ldap.ScopeSingleLevel, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(member=%s)", identifiant),
		[]string{"cn"},
		nil,
	)
	// Search for the given groups
	sr, err := l.Search(searchRequest)
	if err != nil {
		tools.Error("LdapGroupsWithMember: error %v", err)
		return nil
	}
	// Transform to list of names
	ret := make([]string, len(sr.Entries))
	for i, gr := range sr.Entries {
		ret[i] = gr.GetAttributeValue("cn")
	}
	// Group cn list
	return ret
}
