package referent

import (
	"fmt"
	"log"
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
		Host string
		Port string
		Root string
	}
	Annuary struct {
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

func NewRefLdap(filename string) (*RefLdap, error) {

	// Config yml : doit mapper sur un tableau de User
	var cfgLdap LdapConfig
	err := tools.LoadAnyConfig(tools.Cfg.Ref.Path, &cfgLdap)
	if err != nil {
		return nil, err
	}
	// Url de connection  (TODO slapd)
	ldapInfo := fmt.Sprintf("ldap://%s:%s", cfgLdap.Database.Host, cfgLdap.Database.Port)

	// Connection open
	ret := RefLdap{cfg: &cfgLdap}
	ret.conn, err = ldap.DialURL(ldapInfo)
	if ret.conn == nil || err != nil {
		log.Println(err)
		return nil, err
	}

	// Anonymous read ldap (username agrements for trace)
	err = ret.conn.UnauthenticatedBind("agrements")
	return &ret, err

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
	usersDn := rf.cfg.Annuary.Users.Path + "," + rf.cfg.Database.Root

	att := strings.Split(rf.cfg.Annuary.Users.Fields, "|")
	if len(att) != 5 {
		return nil
	}
	// Filter on objets of users,and this identifiant pattern
	searchRequest := ldap.NewSearchRequest(usersDn,
		scope(rf.cfg.Annuary.Users.Scope), ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(rf.cfg.Annuary.Users.Search, id),
		att, nil)

	// Search for the matching usernames
	sr, err := rf.conn.Search(searchRequest)
	if err != nil {
		log.Println(err)
		return nil
	}
	// Get tables
	ret := make([]obj.User, len(sr.Entries))
	for i, entry := range sr.Entries {
		// Get essentials attributes
		user := obj.User{
			Uid:         entry.GetAttributeValue(att[0]),
			Identifiant: entry.GetAttributeValue(att[1]),
			NomComplet:  entry.GetAttributeValue(att[2]),
			AdresseMail: entry.GetAttributeValue(att[3]),
			Roles:       entry.GetAttributeValue(att[4])}
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
	att := strings.Split(rf.cfg.Annuary.Users.Fields, "|")
	if len(att) != 4 {
		return nil
	}

	groupDn := rf.cfg.Annuary.Group.Path + "," + rf.cfg.Database.Root

	// Filter on objets of users,and this identifiant pattern
	searchRequest := ldap.NewSearchRequest(groupDn,
		scope(rf.cfg.Annuary.Group.Scope), ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(rf.cfg.Annuary.Group.Search, id),
		[]string{rf.cfg.Annuary.Group.Member}, nil)

	// Search for the matching Group
	sr, err := rf.conn.Search(searchRequest)
	if err != nil {
		log.Println(err)
		return nil
	}
	// Make users tables
	ret := []obj.User{}
	for _, entry := range sr.Entries {
		// Members of a found group
		members := entry.GetAttributeValues(rf.cfg.Annuary.Group.Member)
		for _, member := range members {
			// Dn of member (if not already)
			mpath := fmt.Sprintf(rf.cfg.Annuary.Group.Mpath, member)
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
	if err != nil {
		log.Println(err)
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
		Roles:       sr.Entries[0].GetAttributeValue(att[4])}
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
