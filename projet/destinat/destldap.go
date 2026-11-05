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

	ldapCfg := tools.Cfg.Dest["ldap"]
	if tools.Check(ldapCfg) {

		// Ouvrir la connexion en lecture/ecriture
		conn, err := bases.LdapOpen(ldapCfg.BaseUrl, ldapCfg.IdentAccess, ldapCfg.TokenPasswd)

		// Créer l'objet Destination
		dest := DestLdap{ldap: conn}
		return &dest, err
	}
	return nil, nil
}

func (de *DestLdap) Close() {
	if de.ldap != nil {
		bases.LdapClose(de.ldap)
	}
}

// Methodes d'interface

func (de *DestLdap) SearchUsersRole(prj *obj.Project) []*obj.User {
	// membres du groupe projet  (pas des autres groupes !)
	projname := prj.Identifiant
	members := bases.LdapGetGroup(de.ldap, projname)
	if len(members) > 0 {
		tab := make([]*obj.User, len(members))
		for i, m := range members {
			ul := bases.LdapSearch(de.ldap, m)
			if len(ul) == 1 {
				tab[i] = ul[0]
			}
		}
		return tab
	}
	return nil
}

func (de *DestLdap) GetRole(prj *obj.Project, user *obj.User) string {
	// Dummy value to force call of create role (get role is done into)
	// Ldap has no role : role <==> (groupe named as role) has member
	return ""
}

func (de *DestLdap) CreateRole(prj *obj.Project, user *obj.User, role string) error {

	// role du projet : ajoute l'user au groupe du nom du projet
	if role == "groups" {
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
	} else if role != "users" && role != "" {
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
	tools.Debug("LdapDeleteMember %s %s", projname, user.Identifiant)
	return nil
}

// ldap projects methodes
func (de *DestLdap) SearchProjects(projet string) []*obj.Project {
	// projet = group dans le ldap
	prj := bases.LdapSearchGroups(de.ldap, projet)
	if len(prj) == 0 {
		// Pour le ldap, si le projet n'existe pas, on le créé automatiquement
		err := bases.LdapAddGroup(de.ldap, projet, "Project group")
		if err == nil {
			prj = append(prj, obj.NewProject("", projet, projet, "Project group"))
		}
	}
	return prj
}

// Ldap users methodes
func (de *DestLdap) SearchUsers(id string) []*obj.User {
	return bases.LdapSearch(de.ldap, id)
}

func (de *DestLdap) CreateUser(user *obj.User) error {
	u := bases.LdapAddUser(de.ldap, user)
	tools.Debug("LdapAddUser %s", user.Identifiant)
	return u
}

func (de *DestLdap) UpdateUser(user *obj.User) error {
	u := bases.LdapUpdateUser(de.ldap, user)
	tools.Debug("LdapUpdateUser %s", user.Identifiant)
	return u
}

func (de *DestLdap) DeleteUser(user *obj.User) error {
	u := bases.LdapDeleteUser(de.ldap, user.Identifiant)
	tools.Debug("LdapDeleteUser %s", user.Identifiant)
	return u
}
