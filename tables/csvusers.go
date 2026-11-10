package tables

import (
	"io"
	obj "main/objets"
	"main/tools"
)

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

func ReadUsers(data any) []obj.User {

	var tab [][]string
	if filename, ok := data.(string); ok {
		tab = readCsv(filename)

	} else if tab, ok = data.([][]string); ok {
		// already

	} else if stream, ok := data.(io.Reader); ok {
		tab = readCsv(stream)

	} else {
		tools.Fatal("Bad ReadUsersRights parameter")
		return nil
	}
	users := make([]obj.User, len(tab))
	for i, t := range tab {
		// See order userToStrings
		users[i] = obj.User{Uid: t[4], Identifiant: t[2], MotDePasse: t[3], AdresseMail: t[1], NomComplet: t[0], Roles: t[5]}
	}
	return users
}

func WriteUsers(users []obj.User, out any) error {
	tab := make([][]string, len(users)+1)
	tab[0] = []string{"#nom complet", "mail", "identifiant", "password", "id", "roles titres"}
	for i, user := range users {
		tab[i+1] = userToStrings(user)
	}
	return writeCsv(out, tab)
}
