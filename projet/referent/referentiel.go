package referent

import (
	obj "main/objets"
	"main/tools"
	"os"
	"strings"
)

// Interface des types de Référentiels gérés
type Referentiel interface {

	// Recherche un utilisateur via son identifiant, ou une partie
	// Retourne une table d'utilisateurs trouvés (un seul en général)
	SearchUser(id string) []obj.User

	// Retourne un groupe d'utilisateurs trouvés
	SearchGroup(gr string) []obj.User

	Close()
}

// Construit des implementations de référentiels d'utilisateurs divers
//   - csv file
//   - openldap
//   - rest api
//   - sql base
func RefFactory() Referentiel {

	var ref Referentiel

	// Annuaire d'utlisateurs dans un fichier  .csv  (fileref.go)
	switch tools.Cfg.Ref.Type {
	case "file":
		ref, _ = NewRefFile(tools.Cfg.Ref.Path) // Csv file

	case "rest":
		ref, _ = NewRefRest(tools.Cfg.Ref.Path) // Api rest

	case "sql":
		ref, _ = NewRefSql(tools.Cfg.Ref.Path) // Sql db

	case "ldap":
		ref, _ = NewRefLdap(tools.Cfg.Ref.Path) // Ldap

		// Rien d'autres implémenté pour le moment
	default:
		println("Référentiel d'utilisateur non pris en charge :" + tools.Cfg.Ref.Type)
		os.Exit(-1)
	}
	return ref
}

func RefClose(ref Referentiel) {
	ref.Close()
}

func SearchEmails(ref Referentiel, name string) []obj.User {

	// S'il y a un filter de nom name dans la conf referentiel
	for f, v := range tools.Cfg.Ref.Filters {
		if f == name {

			// Groupe à chercher
			var tab []obj.User

			// Matching suffixes
			fils := strings.Split(v, "|")

			// Cherche dans tous les utilisateurs
			users := ref.SearchUser("*")
			for _, u := range users {
				for _, fi := range fils {
					if strings.HasPrefix(fi, "!") {
						// Si matche le pattern négatif
						if !strings.HasSuffix(u.AdresseMail, fi[1:]) {

							// Ajoute au groupe
							tab = append(tab, u)
						}
					} else {
						// Si matche le pattern
						if strings.HasSuffix(u.AdresseMail, fi) {

							// Ajoute au groupe
							tab = append(tab, u)
						}
					}
				}
			}
			// Retourne le groupe
			return tab
		}
	}
	return nil
}

func IntoEmails(ref Referentiel, user *obj.User, name string) bool {

	// S'il y a un filter de nom name dans la conf referentiel
	for f, v := range tools.Cfg.Ref.Filters {
		if f == name {

			// Matching suffixes
			fils := strings.Split(v, "|")
			for _, fi := range fils {
				if strings.HasPrefix(strings.TrimSpace(fi), "!") {
					// Si matche le pattern négatif
					if !strings.HasSuffix(user.AdresseMail, strings.TrimSpace(fi[1:])) {
						return true
					}
				} else {
					// Si matche le pattern
					if strings.HasSuffix(user.AdresseMail, strings.TrimSpace(fi)) {
						return true
					}
				}
			}
		}
	}
	// Ne matche rien
	return false
}
