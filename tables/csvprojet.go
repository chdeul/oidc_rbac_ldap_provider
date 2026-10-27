package tables

type UserRights struct {
	Type  string
	Name  string
	Tools []string
}

func userRightsToStrings(user UserRights) []string {
	t := make([]string, 2)
	t[0] = user.Type
	t[1] = user.Name
	t = append(t, user.Tools...)
	return t
}

func ReadUsersRights(filename string) []UserRights {

	// TODO  lecture par table markdown

	tab := readCsv(filename)
	users := make([]UserRights, len(tab))
	for i, t := range tab {
		users[i] = UserRights{t[0], t[1], t[2:]}
	}
	return users
}

func WriteUsersRights(users []UserRights, filename string) error {

	tab := make([][]string, len(users)+1)
	tab[0] = []string{"#type", "name", "tools..."}
	for i, user := range users {
		tab[i+1] = userRightsToStrings(user)
	}
	return writeCsv(filename, tab)
}
