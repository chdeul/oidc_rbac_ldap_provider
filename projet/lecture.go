package projet

import (
	obj "main/objets"
	"main/projet/destinat"
	"main/projet/referent"
	"main/tables"
	"main/tools"
)

// TODO  Réception de droits projet par API
// TODO  Sauvegarde / centralisation de droits projet par API
// TODO  GUI Web d'édtion de droits projet
// TODO  Page html, markdown, descrivant le projet et les liens

func LectureProjet(filename string) {

	// Droits lus dans un .csv
	users := tables.ReadUsersRights(filename)
	header := users[0]

	// Referentiel d'utilisateur indiqué
	ref := referent.RefFactory()
	defer referent.RefClose(ref)

	// Destinataires de droits sur projet
	dest := destinat.DestFactories(header.Tools)
	defer destinat.DestCloseAll(dest)

	// Projet dans chaque outil destination
	projet := searchProjects(header.Type, dest)

	// Pour chaque utilisateur ou groupe déclaré (sauf entete)
	for _, up := range users[1:] {

		// traitement
		switch up.Type {

		case "User": // Declare les users et leurs roles

			// Verifier syntaxe des identifiants
			if !tools.IsIdent(up.Name) {
				continue
			}

			// Recherche l'utilisateur dans le référentiel
			users := ref.SearchUser(up.Name)

			// non trouvé ou non unique
			if len(users) != 1 {
				uniq := ""
				if len(users) > 1 {
					uniq = "non-unique"
				}
				tools.Error("Utilisateur non trouvé dans l'annuaire : %s %s", up.Name, uniq)
				continue
			}
			// Trouvé
			user := &users[0]

			// add role if
			if referent.IntoEmails(ref, user, "external") {
				user.Roles += ",external"
			}

			switch up.Type[0] {
			case 'U':
				updateRights(projet, user, dest, up.Tools)
			}

		case "Group": // Déclare les groupes et leurs roles

			// Verifier syntaxe des noms de groupe
			if !tools.IsGroup(up.Name) {
				continue
			}

			// Groupe constitués sur filtrage email référentiel  !@internal.com, @societe.com
			users := referent.SearchEmails(ref, up.Name)

			// Sinon groupe de role: name
			if users == nil {
				// Recherche des utilisateurs de ce groupe dans le référentiel
				users = ref.SearchGroup(up.Name)
			}
			// Aucun trouvé
			if len(users) == 0 {
				tools.Warn("Groupe vide dans l'annuaire : %s", up.Name)
				continue
			}
			// Chaque user du groupe
			for _, user := range users {

				// TODO add role if
				//if referent.IntoEmails(ref,user,"external") {
				//	user.Roles += ",external"
				//}

				switch up.Type[0] {
				case 'G':
					updateRights(projet, &user, dest, up.Tools)
				}
			}

		default:
			tools.Error("Type dans projet inconnu %s", up.Type)
		}
	}
	// TODO Liste des utilisateurs déclarés ( avec ceux des groupes )
	ud := users[1:]

	// 2eme passe : Suppression ou Notification de droits non listés dans le projet
	// Pour chaque outil destination
	for i, d := range dest {

		// Projet dans l'outil
		if projet[i] == nil {
			continue
		}
		removeRights(d, projet[i], ud)
	}
	// Note: Suppression utilisateur dans l'outil pas ici : par batch
	// car peut avoir des droits sur un autre projet
}

// Recherche le projet de nom donné dans chaque outil
func searchProjects(prjname string, dest []destinat.Destination) []*obj.Project {

	// tableau des projets de chaque outil
	tab := make([]*obj.Project, len(dest))

	// Pour chaque outil destination
	for i, d := range dest {
		p := d.SearchProjects(prjname)
		if len(p) == 1 {
			tab[i] = p[0]
		} else if len(p) > 1 {
			tools.Error("multiple projets trouvés pour %s dans %s", prjname, d.Name)
		}
	}
	return tab
}

// Cree ou met à jour si besoin les roles, pour chaque destination, pour cet utilisateur
func updateRights(projet []*obj.Project, user *obj.User, dest []destinat.Destination, droit []string) {

	// Pour chaque outil destination
	for i, d := range dest {

		// Projet dans l'outil
		if projet[i] == nil {
			tools.Error("projet inexistant dans %s", d.Name())
			continue
		}

		// Droit dans outil non vide pour cet utilisateur
		if droit[i] != "" {

			// Cherche l'utilisateur dans l'outil
			ud := d.SearchUsers(user.Identifiant)

			// Si non trouvé : le Crée
			if len(ud) == 0 {
				err := d.CreateUser(user)
				if err != nil {
					continue
				}
			} else if len(ud) == 1 {
				// Trouvé et à mettre à jour
				if ud[0].AdresseMail != user.AdresseMail || ud[0].NomComplet != user.NomComplet {
					d.UpdateUser(user)
				}
			} else {
				tools.Error("utilisateur multiple %s dans %s", user.Identifiant, d.Name())
				continue
			}
			// Créer ou modifier le role pour cet utilisateur et ce projet
			role := d.GetRole(projet[i], user)
			if droit[i] != role {
				if role == "" {
					// Créer le role pour cet utilisateur et ce projet
					d.CreateRole(projet[i], user, droit[i])
				} else {
					// Modifier le role pour cet utilisateur et ce projet
					// TODO  Warning si ancien role < droit[i] nouveau + cas dans group et hors
					d.UpdateRole(projet[i], user, droit[i])
				}
			}
		}
	}

}

// Notifie, voire supprime des droits non déclarés
func removeRights(d destinat.Destination, prj *obj.Project, users []tables.UserRights) {

	// Recherche tous les utilisateurs avec droits sur ce projet dans l'outil d
	ud := d.SearchUsersRole(prj)
	for _, user := range ud {

		// Pour chaque utilisateur ou groupe déclaré
		found := false
		for i, up := range users[1:] {

			if up.Type == "User" {

				// Si dans la liste et avec droit : ok, on le gère
				found = up.Name == user.Identifiant
				if found && up.Tools[i] == user.Roles {
					break
				}
				// TODO droit updaté & déclassé : droit différent
				//if found && up.Tools[i] > user.Roles {
				//   warning
				//}
			}
			if up.Type == "Group" {
				// Notifie les groupes et leurs roles non listés
				// TODO droit group	non listé
			}
		}
		if !found {
			// TODO utilisateur/droit ajouté manuellement : suppression/ajout droit éventuel,à proposer
			// si pas dans la liste ou sans droit  : warning, voire suppression droit
			// d.DeleteRole(projet[i], user)
		}
	}
}
