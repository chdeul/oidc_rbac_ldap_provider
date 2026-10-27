package referent

import (
	"fmt"
	obj "main/objets"
	"main/tables"
	"main/tools"
	"sort"
	"strings"
)

/*
 * Liste d'utilisateurs chargée en mémoire à partir d'un fichier csv
 *  nommé  users.csv    dans ./data
 */

type RefFile struct {

	// Referentiel d'utilisateur
	users   []obj.User
	maxuser int
}

func NewRefFile(filename string) (*RefFile, error) {

	// Lire le fichier csv
	tab := tables.ReadUsers(filename)
	max := len(tab)

	// Verifier syntaxe des usernames
	for _, u := range tab {
		if !tools.IsIdent(u.Identifiant) {
			return nil, fmt.Errorf("Mauvaise syntaxe d'identifiant %s", u.Identifiant)
		}
	}
	// Trier le tableau
	sort.Slice(tab, func(i, j int) bool {
		return tab[i].Identifiant < tab[j].Identifiant
	})
	return &RefFile{users: tab, maxuser: max}, nil
}

func (rf *RefFile) Close() {
	rf.users = nil
	rf.maxuser = 0
}

// Methode d'interface
func (rf *RefFile) SearchUser(ident string) []obj.User {

	// Cas username exact recherché
	if !strings.ContainsAny(ident, tools.MatchChars) {

		// Recherche rapide ( tableau trié )
		idx := sort.Search(rf.maxuser, func(i int) bool {
			return rf.users[i].Identifiant >= ident
		})
		// Trouvé
		if idx < rf.maxuser && rf.users[idx].Identifiant == ident {
			return []obj.User{rf.users[idx]}
		}
		return nil
	}
	// Regex de recherche % _ .*
	re := tools.InitMatch(ident)

	// Recherche par pattern matching
	tab := make([]obj.User, 0)
	for _, u := range rf.users {
		if tools.Match(re, u.Identifiant) {
			tab = append(tab, u)
		}
	}
	// Celui ou ceux trouvé
	return tab
}

// Methode d'interface
func (rf *RefFile) SearchGroup(ident string) []obj.User {

	// Regex de recherche % _ .*
	re := tools.InitMatch(ident)

	// Recherche par pattern matching
	tab := make([]obj.User, 0)
	for _, u := range rf.users {
		if tools.Match(re, u.Roles) {
			tab = append(tab, u)
		}
	}
	// Ceux trouvés
	return tab
}
