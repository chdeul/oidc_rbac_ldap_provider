package referent

import (
	"fmt"
	obj "main/objets"
	"main/tools"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

/*
 * Liste d'utilisateurs chargée à partir d'un annuaire ldap
 *  configuré à partir de  ldapInput.yml    dans ./data
 */

type LdapConfig struct {
	Database struct {
		Host     string
		Port     string
		Root     string
		Bind     string
		Password string
	}
	Directory struct {
		Users struct {
			Scope  string
			Path   string
			Search string
			Fields string
		}
		Group struct {
			Scope  string
			Path   string
			Search string
			Member string
			Mpath  string
		}
	}
}

type RefLdap struct {
	conn *ldap.Conn
	cfg  *LdapConfig
}

func NewRefLdap() (*RefLdap, error) {

	// Config yml : doit mapper sur un tableau de User
	var cfgLdap LdapConfig
	err := tools.LoadAnyConfig(tools.Cfg.Ref.Path, &cfgLdap)
	tools.AssertError(err, "Referentiel config")

	// Url de connection  (TODO slapd)
	ldapInfo := fmt.Sprintf("ldap://%s:%s", cfgLdap.Database.Host, cfgLdap.Database.Port)

	// Connection open
	ret := RefLdap{cfg: &cfgLdap}
	ret.conn, err = ldap.DialURL(ldapInfo)
	if !tools.Check(ret.conn != nil && err == nil, "Referentiel connection") {
		return nil, err
	}

	if cfgLdap.Database.Password == "" {
		// Anonymous read ldap : userPassword will be empty
		err = ret.conn.UnauthenticatedBind("")
	} else {
		// Auth : with reading of user Password
		err = ret.conn.Bind(cfgLdap.Database.Bind, cfgLdap.Database.Password)
	}
	return &ret, err

}

func (rf *RefLdap) Type() string {
	return "ldap"
}

func (rf *RefLdap) Close() {
	if rf.conn != nil {
		rf.conn.Close()
		rf.conn = nil
	}
}

// Methode d'interface
func (rf *RefLdap) SearchUser(id string) []obj.User {

	// Rien à chercher : retourne vide
	if id == "" {
		return nil
	}
	usersDn := rf.cfg.Directory.Users.Path + "," + rf.cfg.Database.Root

	att := strings.Split(rf.cfg.Directory.Users.Fields, "|")
	if !tools.Assert(len(att) == 6, "Ref Ldap fields") {
		return nil
	}
	// Filter on objets of users,and this identifiant pattern
	searchRequest := ldap.NewSearchRequest(usersDn,
		scope(rf.cfg.Directory.Users.Scope), ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(rf.cfg.Directory.Users.Search, id),
		att, nil)

	// Search for the matching usernames
	sr, err := rf.conn.Search(searchRequest)
	if !tools.CheckError(err, "Ldap search user") {
		return nil
	}

	// Get tables
	ret := make([]obj.User, len(sr.Entries))
	for i, entry := range sr.Entries {
		tools.Check(len(entry.Attributes) == 6, "Ldap search att")
		// Get essentials attributes
		user := obj.User{
			Uid:         entry.GetAttributeValue(att[0]),
			Identifiant: entry.GetAttributeValue(att[1]),
			NomComplet:  entry.GetAttributeValue(att[2]),
			AdresseMail: entry.GetAttributeValue(att[3]),
			MotDePasse:  entry.GetAttributeValue(att[4]),
			Roles:       entry.GetAttributeValue(att[5])}
		ret[i] = user
	}
	return ret
}

// Methode d'interface
func (rf *RefLdap) SearchGroup(id string) []obj.User {

	// Rien à chercher : retourne vide
	if id == "" {
		return nil
	}
	att := strings.Split(rf.cfg.Directory.Users.Fields, "|")
	if !tools.Assert(len(att) == 6, "Ref Ldap fields") {
		return nil
	}

	groupDn := rf.cfg.Directory.Group.Path + "," + rf.cfg.Database.Root

	// Filter on objets of users,and this identifiant pattern
	searchRequest := ldap.NewSearchRequest(groupDn,
		scope(rf.cfg.Directory.Group.Scope), ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(rf.cfg.Directory.Group.Search, id),
		[]string{rf.cfg.Directory.Group.Member}, nil)

	// Search for the matching Group
	sr, err := rf.conn.Search(searchRequest)
	if !tools.CheckError(err, "Ldap search group") {
		return nil
	}
	// Make users tables
	ret := []obj.User{}
	for _, entry := range sr.Entries {
		// Members of a found group
		members := entry.GetAttributeValues(rf.cfg.Directory.Group.Member)
		for _, member := range members {
			// Dn of member (if not already)
			mpath := fmt.Sprintf(rf.cfg.Directory.Group.Mpath, member)
			// Get user from dn member
			user := getUser(rf.conn, mpath, att)
			// If exist and not already in table
			if user != nil && !contains(&ret, user.Identifiant) {
				ret = append(ret, *user)
			}
		}
	}
	// Return tables
	return ret
}

// Get simple user from dn
func getUser(conn *ldap.Conn, dn string, att []string) *obj.User {

	// Filter on objets of users,and this identifiant pattern
	searchRequest := ldap.NewSearchRequest(dn,
		ldap.ScopeBaseObject, ldap.NeverDerefAliases, 0, 0, false,
		"(objectClass=*)", att, nil)

	// Search for the username
	sr, err := conn.Search(searchRequest)
	if !tools.CheckError(err, "Ldap search user") {
		return nil
	}
	// Must be unique
	if len(sr.Entries) != 1 {
		return nil
	}
	// Get essentials attributes
	user := obj.User{
		Uid:         sr.Entries[0].GetAttributeValue(att[0]),
		Identifiant: sr.Entries[0].GetAttributeValue(att[1]),
		NomComplet:  sr.Entries[0].GetAttributeValue(att[2]),
		AdresseMail: sr.Entries[0].GetAttributeValue(att[3]),
		MotDePasse:  sr.Entries[0].GetAttributeValue(att[4]),
		Roles:       sr.Entries[0].GetAttributeValue(att[5])}
	return &user
}

// Table of users contains this identifiant
func contains(slice *[]obj.User, val string) bool {
	for _, item := range *slice {
		if item.Identifiant == val {
			return true
		}
	}
	return false
}

// Deep of searching from config
func scope(name string) int {

	switch name {
	case "BaseObject":
		return ldap.ScopeBaseObject
	case "SingleLevel":
		return ldap.ScopeSingleLevel
	case "WholeSubtree":
		return ldap.ScopeWholeSubtree
	case "Children":
		return ldap.ScopeChildren
	default:
		return ldap.ScopeChildren
	}
}
