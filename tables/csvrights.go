package tables

import (
	"io"
	obj "main/objets"
	"main/tools"
	"strings"
)

func RightToStrings(user obj.Right) []string {
	t := make([]string, 2)
	t[0] = user.Type
	t[1] = user.Name
	t = append(t, user.Tools...)
	return t
}

func ReadRights(data any) []obj.Right {

	var tab [][]string
	if filename, ok := data.(string); ok {
		if strings.HasSuffix(filename, ".md") {
			// TODO  lecture par table markdown
		} else {
			tab = readCsv(filename)
		}

	} else if tab, ok = data.([][]string); ok {
		// already

	} else if stream, ok := data.(io.Reader); ok {
		tab = readCsv(stream)

	} else {
		tools.Fatal("Bad ReadUsersRights parameter")
		return nil
	}
	nb := 0
	users := make([]obj.Right, len(tab))
	for _, t := range tab {
		if t[0] != "" && t[1] != "" {
			users[nb] = *obj.NewRight(t[0], t[1], t[2:])
			nb++
		}
	}
	if nb < len(tab) {
		// Reduit l'allocation si beaucoup de lignes vides
		newS := make([]obj.Right, nb)
		copy(newS, users)
		return newS
	}
	return users
}

func WriteRights(users []obj.Right, out any) error {

	tab := make([][]string, len(users))
	for i, user := range users {
		tab[i] = RightToStrings(user)
	}
	return writeCsv(out, tab)
}
