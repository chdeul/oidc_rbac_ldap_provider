package tables

import obj "main/objets"

func NewUser(uid, ident, mp, nc, am, rt string) *obj.User {
	u := obj.User{Uid: uid, Identifiant: ident, MotDePasse: mp, NomComplet: nc, AdresseMail: am, Roles: rt}
	return &u
}

func userToStrings(user obj.User) []string {
	t := make([]string, 6)
	t[0] = user.Uid
	t[1] = user.Identifiant
	t[2] = user.MotDePasse
	t[3] = user.AdresseMail
	t[4] = user.NomComplet
	t[5] = user.Roles
	return t
}

func ReadUsers(filename string) []obj.User {

	tab := readCsv(filename)
	users := make([]obj.User, len(tab))
	for i, t := range tab {
		users[i] = obj.User{Uid: t[0], Identifiant: t[1], MotDePasse: t[2], AdresseMail: t[3], NomComplet: t[4], Roles: t[5]}
	}
	return users
}

func WriteUsers(users []obj.User, filename string) error {

	tab := make([][]string, len(users)+1)
	tab[0] = []string{"#id", "identifiant", "password", "mail", "nom complet", "roles titres"}
	for i, user := range users {
		tab[i+1] = userToStrings(user)
	}
	return writeCsv(filename, tab)
}

func ReadUser(identifiant string, filename string) *obj.User {

	// Cherche dans Table des utilisateurs
	users := ReadUsers(filename)
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
	users := ReadUsers(filename)
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
	users := ReadUsers(filename)
	for i := range users {

		// Recupere les données de l'utilisateur d'identifiant
		if users[i].Identifiant == data.Identifiant {
			users[i].Uid = data.Uid
			users[i].AdresseMail = data.AdresseMail
			users[i].MotDePasse = data.MotDePasse
			users[i].NomComplet = data.NomComplet
			users[i].Roles = data.Roles
			return WriteUsers(users, filename)
		}
	}
	// Non trouvé l'ajoute
	return WriteUsers(append(users, *data), filename)
}
