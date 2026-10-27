package storage

import (
	"crypto/rsa"
	"main/bases"
	"main/tools"
)

const (
	// ServiceUserID is the ID of the service user.
	ServiceUserID = "openldap"
	// ServiceUserKeyID is the key ID of the service user.
	ServiceUserKeyID = "dn=ldap,dn=example,dn=org"
)

type User struct { /// ex                // ldap               // oidc
	ID       string // antoine           // uid                // userInfo.Subject
	Username string // antoine           // sn                 // userInfo.PreferredUsername (login,identifiant)
	Name     string // Antoine Chdeul    // cn                 // userInfo.Name
	Email    string // email             // mail               // userInfo.Email
	Profiles string // manager,service   // title              // userInfo.Roles  (custom scope role)
}

type Service struct {
	keys map[string]*rsa.PublicKey
}

type UserStore interface {
	GetUserByID(string) *User
	CheckPassword(string, string) bool
	ServiceUserStoreID() string
}

// OIDC User Store use Agréments destination ldap

type userStore struct {
	conn any
}

func NewUserStore(issuer string) UserStore {

	conn, err := bases.LdapOpen(tools.Cfg.Dest.Ldap.Url, tools.Cfg.Dest.Ldap.Root, "")
	if err != nil {
		panic("Impossible de se connecter au LDAP d'agréments")
	}
	store := userStore{conn}
	return store
}

// ServiceUserStoreID is only used in the appli server
func (u userStore) ServiceUserStoreID() string {
	return ServiceUserID
}

func (u userStore) GetUserByID(ident string) *User {
	tu := bases.LdapGetUser(u.conn, ident)
	if tu != nil {
		user := User{ID: tu.Uid, Username: tu.Identifiant, Name: tu.NomComplet, Email: tu.AdresseMail, Profiles: tu.Roles}
		return &user
	}
	return nil
}

func (u userStore) CheckPassword(ident string, password string) bool {

	return bases.LdapCheckPassword(tools.Cfg.Dest.Ldap.Url, tools.Cfg.Dest.Ldap.Root, ident, password)
}

func CloseStore(us UserStore) {
	bases.LdapClose(us.(userStore).conn)
}
