package obj

type User struct {
	Uid         string // Numéro (uid) propre à chaque outil ou idem identifiant
	Identifiant string // Nom court (sn) : champ "Clé" de la table
	MotDePasse  string // Crypté (userPassword) si possible  ou vide
	AdresseMail string // Adresse mail (mail) pas forcément unique
	NomComplet  string // Nom (cn) split(Prénom Nom)
	Roles       string // Roles (title), Services, divers  ex:  Manager,Security, propre à chaque outil
}

func NewUser(id, ident, nom, mail, roles string) *User {
	return &User{Uid: id, Identifiant: ident, NomComplet: nom, AdresseMail: mail, Roles: roles}
}
