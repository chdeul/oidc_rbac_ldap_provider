package destinat

import (
	"main/bases"
	obj "main/objets"
	"main/tools"
	"slices"
)

/*
 * Annuaire LDAP avec  users et groups statiques
 */

type DestLdap struct {
	ldap any // La connexion ldap
}

func (de *DestLdap) Name() string {
	return "ldap"
}

func NewDestLdap() (*DestLdap, error) {

	// Créer l'objet Destination
	dest := DestLdap{}

	// Ouvrir la connexion en lecture/ecriture
	var err error
	dest.ldap, err = bases.LdapOpen(tools.Cfg.Dest.Ldap.Url, tools.Cfg.Dest.Ldap.Root, tools.Cfg.Dest.Ldap.Passwd)
	return &dest, err
}

func (de *DestLdap) Close() {
	if de.ldap != nil {
		bases.LdapClose(de.ldap)
	}
}

// Methodes d'interface

func (de *DestLdap) SearchUsersRole(prj *obj.Project) []*obj.User {
	return nil
}

func (de *DestLdap) GetRole(prj *obj.Project, user *obj.User) string {
	// Dummy value to forge call of create role (get role is done into)
	return ""
}

func (de *DestLdap) CreateRole(prj *obj.Project, user *obj.User, role string) error {

	// role du projet : ajoute l'user au groupe du nom du projet
	if role == "!project" {
		projname := prj.Identifiant
		members := bases.LdapGetGroup(de.ldap, projname)
		if members == nil {
			// Crée le groupe pour l'ajouter membre
			bases.LdapAddGroup(de.ldap, projname, "")
			tools.Debug("LdapAddGroup %s", projname)
			bases.LdapAddMember(de.ldap, projname, user.Identifiant)
			tools.Debug("LdapAddMember %s %s", projname, user.Identifiant)
		} else if !slices.Contains(members, user.Identifiant) {
			// pas déjà membre : l'ajoute
			bases.LdapAddMember(de.ldap, projname, user.Identifiant)
			tools.Debug("LdapAddMember %s %s", projname, user.Identifiant)
		}
		// role "autre" que projet et users  : ajoute l'user à ce groupe "autre"
	} else if role != "!user" && role != "" {
		members := bases.LdapGetGroup(de.ldap, role)
		if members == nil {
			// Crée le groupe pour l'ajouter membre
			bases.LdapAddGroup(de.ldap, role, "")
			tools.Debug("LdapAddGroup %s", role)
			bases.LdapAddMember(de.ldap, role, user.Identifiant)
			tools.Debug("LdapAddMember %s %s", role, user.Identifiant)
		} else if !slices.Contains(members, user.Identifiant) {
			// pas déjà membre : l'ajoute
			bases.LdapAddMember(de.ldap, role, user.Identifiant)
			tools.Debug("LdapAddMember %s %s", role, user.Identifiant)
		}
	}
	// Rien d'autre à faire, pas errreur
	return nil
}

func (de *DestLdap) UpdateRole(prj *obj.Project, user *obj.User, role string) error {
	// done into Create role
	return nil
}

func (de *DestLdap) DeleteRole(prj *obj.Project, user *obj.User) error {
	projname := prj.Identifiant
	bases.LdapDeleteMember(de.ldap, projname, user.Identifiant)
	return nil
}

// ldap projects methodes
func (de *DestLdap) SearchProjects(projet string) []*obj.Project {
	return bases.LdapSearchGroups(de.ldap, projet)
}

// Ldap users methodes
func (de *DestLdap) SearchUsers(id string) []*obj.User {
	return bases.LdapSearch(de.ldap, id)
}

func (de *DestLdap) CreateUser(user *obj.User) error {
	return bases.LdapAddUser(de.ldap, user)
}

func (de *DestLdap) UpdateUser(user *obj.User) error {
	return bases.LdapUpdateUser(de.ldap, user)
}

func (de *DestLdap) DeleteUser(user *obj.User) error {
	return bases.LdapDeleteUser(de.ldap, user.Identifiant)
}
