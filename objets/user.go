package obj

type User struct {
	Uid         string
	Identifiant string
	MotDePasse  string
	AdresseMail string
	NomComplet  string
	Roles       string // Roles, Services, divers  ex:  Manager,Security
}

func NewUser(id, ident, nom, mail, roles string) *User {
	return &User{Uid: id, Identifiant: ident, NomComplet: nom, AdresseMail: mail, Roles: roles}
}
