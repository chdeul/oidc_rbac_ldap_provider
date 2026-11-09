package tables

import (
	"io"
	obj "main/objets"
	"net/http"
)

func NewUser(uid, ident, mp, nc, am, rt string) *obj.User {
	u := obj.User{Uid: uid, Identifiant: ident, MotDePasse: mp, NomComplet: nc, AdresseMail: am, Roles: rt}
	return &u
}

func userToStrings(user obj.User) []string {
	t := make([]string, 6)

	// Order into .csv file
	t[0] = user.NomComplet
	t[1] = user.AdresseMail
	t[2] = user.Identifiant
	t[3] = user.MotDePasse
	t[4] = user.Uid
	t[5] = user.Roles
	return t
}

func ReadUsers(into io.Reader) []obj.User {
	tab := readCsvStream(into)
	users := make([]obj.User, len(tab))
	for i, t := range tab {
		// See order userToStrings
		users[i] = obj.User{Uid: t[4], Identifiant: t[2], MotDePasse: t[3], AdresseMail: t[1], NomComplet: t[0], Roles: t[5]}
	}
	return users
}

func WriteUsers(out io.Writer, users []obj.User) error {
	tab := make([][]string, len(users)+1)
	tab[0] = []string{"#nom complet", "mail", "identifiant", "password", "id", "roles titres"}
	for i, user := range users {
		tab[i+1] = userToStrings(user)
	}
	return writeCsvStream(out, tab)
}

func ReadUsersFile(filename string) []obj.User {

	tab := readCsv(filename)
	users := make([]obj.User, len(tab))
	for i, t := range tab {
		// See order userToStrings
		users[i] = obj.User{Uid: t[4], Identifiant: t[2], MotDePasse: t[3], AdresseMail: t[1], NomComplet: t[0], Roles: t[5]}
	}
	return users
}

func WriteUsersFile(users []obj.User, filename string) error {

	tab := make([][]string, len(users)+1)
	tab[0] = []string{"#nom complet", "mail", "identifiant", "password", "id", "roles titres"}
	for i, user := range users {
		tab[i+1] = userToStrings(user)
	}
	return writeCsv(filename, tab)
}

func ReadUser(identifiant string, filename string) *obj.User {

	// Cherche dans Table des utilisateurs
	users := ReadUsersFile(filename)
	for _, user := range users {

		// Recupere les données de l'utilisateur d'identifiant
		if user.Identifiant == identifiant {
			return &user
		}
	}
	return nil
}

func ReadUserFromEmail(email string, filename string) *obj.User {

	// Cherche dans Table des utilisateurs
	users := ReadUsersFile(filename)
	for _, user := range users {

		// Recupere les données de l'utilisateur d'identifiant
		if user.AdresseMail == email {
			return &user
		}
	}
	return nil
}

func WriteUser(data *obj.User, filename string) error {

	// Cherche dans Table des utilisateurs
	users := ReadUsersFile(filename)
	for i := range users {

		// Recupere les données de l'utilisateur d'identifiant
		if users[i].Identifiant == data.Identifiant {
			users[i].Uid = data.Uid
			users[i].AdresseMail = data.AdresseMail
			users[i].MotDePasse = data.MotDePasse
			users[i].NomComplet = data.NomComplet
			users[i].Roles = data.Roles
			return WriteUsersFile(users, filename)
		}
	}
	// Non trouvé l'ajoute
	return WriteUsersFile(append(users, *data), filename)
}

// Recuperer du csv
func DecodeCsv(r *http.Request, obj interface{}) []obj.User {
	return ReadUsers(r.Body)
}

// Envoyer du csv
func EncodeCsv(w http.ResponseWriter, tab []obj.User) {
	w.Header().Set("Content-Type", "text/csv")
	WriteUsers(w, tab)
}
