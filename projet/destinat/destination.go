package destinat

import (
	"log"
	obj "main/objets"
	"os"
)

// Interface des types de destinations gérés
type Destination interface {

	// TODO  *** Fonctions élémentaires
	Name() string

	// Search Project (or groups)
	SearchProjects(projet string) []*obj.Project

	// Crud Users
	SearchUsers(ident string) []*obj.User
	CreateUser(user *obj.User) error
	UpdateUser(user *obj.User) error
	DeleteUser(user *obj.User) error

	SearchUsersRole(prj *obj.Project) []*obj.User

	// Crud Roles (RBAC : Association Projet/Utilisateur -> Rôle)
	GetRole(prj *obj.Project, user *obj.User) string
	CreateRole(prj *obj.Project, user *obj.User, role string) error
	UpdateRole(prj *obj.Project, user *obj.User, role string) error
	DeleteRole(prj *obj.Project, user *obj.User) error

	Close()
}

// Construit des implementations de types de destinations  divers
//   - openldap
//   - rest api
//   - gitlab
//   - mantis
//   - sql base
//   - alfresco
//   - nexus
//   - harbor
//   - jenkins
//   - sonarqube
//   - rancher
//   - argocd
//   - redmine
//   - etc.
func DestFactory(destname string) Destination {

	var err error
	var dest Destination

	switch destname {
	case "ldap":
		dest, err = NewDestLdap() // done
	case "sqldb":
		dest, err = NewDestSql()
	case "gitlab":
		dest, err = NewDestGitlab() // done
	case "mantis":
		dest, err = NewDestMantis() // done
	default:
		log.Println("Destination-progiciel non pris en charge : " + destname)
	}
	if err != nil {
		os.Exit(-1)
	}
	return dest
}

// Crée les implementations vers les outils de destinations
func DestFactories(tools []string) []Destination {

	tab := make([]Destination, len(tools))
	for i, tool := range tools {
		tab[i] = DestFactory(tool)
	}
	return tab
}

func DestCloseAll(dest []Destination) {

	for _, de := range dest {
		de.Close()
	}
}
