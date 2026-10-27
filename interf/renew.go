package interf

import (
	"fmt"
	"main/bases"
	obj "main/objets"
	"main/tools"
	"strings"

	"github.com/google/uuid"
)

// Envoie ou traite un mail de reinit du mot de passe
func RenewPassword(identifiant, motDePasse, adresseMail, requestId, renewUrl string) (string, error) {

	var err error
	var user *obj.User

	if identifiant == "" {
		if adresseMail == "" {
			err = fmt.Errorf("Identifiant ou Email vide")
			return "Identifiant ou E-mail vide", err
		}
		// Verifie email  en base
		user = AnnReadUserFromEmail(adresseMail)
		if user == nil {
			err = fmt.Errorf("Email incorrect")
			return "E-mail inconnu, pas de compte associé", err
		}
	} else {
		// Verifie identifiant en base
		user = AnnReadUser(identifiant)
		if user == nil {
			err = fmt.Errorf("Identifiant incorrect")
			return "Identifiant inconnu, pas de compte associé", err
		}

		// Verifie mail comme en base
		if adresseMail != "" && adresseMail != user.AdresseMail {
			err = fmt.Errorf("Identifiant et Email incohérent")
			return "E-mail non associé à cet identifiant", err
		}
	}

	// Si mot de passe : à renouveller (webui)  ou à renouveler (oidc)
	if motDePasse == "$Demande$Renouvellement$" || motDePasse == "!Demande!Renouvellement!" {

		// Genere un jeton provisoire valable 7 jours
		url := generateLink(user, renewUrl)
		if url == "" {
			err = fmt.Errorf("generation lien renouvellement incorrect")
			return "Erreur interne de demande de renouvellement", err
		}
		err = sendRenewMail(user.AdresseMail, user.Identifiant, url)
		if err != nil {
			return "Echec d'envoi du e-mail de renouvellement", err
		}

		tools.Debug("Message de renouvellement envoyé vers %s\n", user.AdresseMail)

		return "E-mail de renouvellement envoyé", err

	} else { // sinon Motdepasse + cle de securité

		// verifie la clé et enregistre le nouveau mot de passe
		err = saveNewPassword(requestId, user, motDePasse)
		if err != nil {
			return "Echec du contrôle du lien de renouvellement", err
		}
	}
	return "Nouveau mot de passe enregistré", err
}

func generateLink(user *obj.User, renewUrl string) string {

	// Genere un jeton provisoire valable 10 h
	clef, err := tools.CreateToken(user.Identifiant, "Renouvellement", "")
	if err != nil {
		return ""
	}
	// Prefix clef de renouvellement (voir login.html)
	clef = "agnew-" + clef

	// Reset du mot de passe
	user.MotDePasse = "agnew-" + uuid.NewString()
	err = AnnWritePassword(user)
	if err != nil {
		return ""
	}

	//  envoi un mail avec une url de reinit du mot de passe
	url := renewUrl + clef
	return url
}

func saveNewPassword(clef string, user *obj.User, password string) error {

	// Verifie la nature de la clef
	if clef == "" || !strings.Contains(clef, "agnew-") {
		err := fmt.Errorf("Clef de renouvellement invalide")
		return err
	}
	clef = strings.Replace(clef, "agnew-", "", 1)

	// Decode le token
	ident, profil, _ := tools.ParseToken(clef)
	if ident == "" || ident != user.Identifiant || profil != "Renouvellement" {
		err := fmt.Errorf("Clef de renouvellement ou identifiant incorrect")
		return err
	}
	// Positionne le nouveau mot de passe
	user.MotDePasse = password
	err := AnnWritePassword(user)
	if err != nil {
		return err
	}
	return err
}

// Main function
func sendRenewMail(email string, identifiant string, url string) error {

	from := "noreply@localhost.loc"

	subject := EncodeSubject("ATTENTION: Vous avez demandé un Renouvellement de mot de passe")
	body := "<p>Bonjour,</p>\r\n"
	body += "<p>Vous avez demandé un Renouvellement de mot de passe<br/>"
	body += "pour le compte d'identifiant <b>" + identifiant + "</b><br/></p>\r\n"
	body += "<p>Aller sur " + url + "  dans les 10h, pour renseigner un nouveau mot de passe.</p>\r\n"
	body += "<p>A bientôt !</p>\r\n"

	msg := fmt.Sprintf("From: %s\r\n", from)
	msg += fmt.Sprintf("To: %s\r\n", email)
	msg += fmt.Sprintf("Subject: %s\r\n", subject)
	msg += "MIME-Version: 1.0;\r\n"
	msg += "Content-Transfer-Encoding: quoted-printable\r\n"
	msg += "Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n"

	msg += fmt.Sprint(body)

	tools.Debug("Renouvellement %s : %s", identifiant, url)

	var err error = nil

	// TODO Envoi du mail de renouvellement
	//err := MySendMail(from, email, []byte(msg))
	//if err != nil {
	//	Error("Echec d'envoi de mail  :%v", err)
	//}
	return err
}

// -------------- annuaire externe (sso) --> notre annuaire (openldap agrément)

func CheckOrCreateUserFromOidcLogin(user *obj.User) error {

	// Cherche cet utilisateur OIDC dans notre annuaire Agréments
	u := AnnReadUser(user.Identifiant)
	if u == nil {
		// Si Pas Auto Create user at oidc login
		if !tools.Cfg.Openidc.Client.CreateUser {
			return fmt.Errorf("Utilisateur non référencé dans Agréments")
		}
		// Créer cet utilisateur
		return AnnAddUser(user)
	}
	// Trouvé mais d'ID différent : pas le même utilisateur ( username supprimé puis réutilisé )
	if u.Uid != user.Uid {
		// Il faudra à l'admin supprimer l'ancien username (et ses droits), pour créer le nouveau (et droits)
		return fmt.Errorf("Un utilisateur de même identifiant, mais d'id différent a été détecté")
	}
	//  Changement de Roles : Ré affectation => pertes ou ajout de droits
	if u.Roles != user.Roles {
		// TODO Chgt Roles : Signaler à la gestion des droits projets
		// Pour mise à jour
	}

	// Adresse mail, Nom, Roles à mettre à jour ?
	if u.NomComplet != user.NomComplet || u.AdresseMail != user.AdresseMail || u.Roles != user.Roles {
		return AnnUpdateUser(user)
	}
	// OK l'utilisateur est dans notre annuaire Agréments à jour
	return nil
}

// ---------------- Renew use Destination Agréments openldap Annuary

func AnnReadUser(identifiant string) *obj.User {
	c := annOpen("")
	defer annClose(c)
	return bases.LdapGetUser(c, identifiant)
}

func AnnReadUserFromEmail(adresseMail string) *obj.User {
	c := annOpen("")
	defer annClose(c)
	return bases.LdapGetUserFromEmail(c, adresseMail)
}

func AnnWritePassword(user *obj.User) error {
	c := annOpen(tools.Cfg.Dest.Ldap.Passwd)
	defer annClose(c)
	return bases.LdapSetPassword(c, user.Identifiant, user.MotDePasse)
}

func AnnAddUser(user *obj.User) error {
	c := annOpen(tools.Cfg.Dest.Ldap.Passwd)
	defer annClose(c)
	return bases.LdapAddUser(c, user)
}

func AnnUpdateUser(user *obj.User) error {
	c := annOpen(tools.Cfg.Dest.Ldap.Passwd)
	defer annClose(c)
	return bases.LdapUpdateUser(c, user)
}

func AnnCheckPassword(identifiant, password string) bool {
	url := tools.Cfg.Dest.Ldap.Url
	bind := tools.Cfg.Dest.Ldap.Root
	return bases.LdapCheckPassword(url, bind, identifiant, password)
}

func annOpen(pass string) any {
	url := tools.Cfg.Dest.Ldap.Url
	bind := tools.Cfg.Dest.Ldap.Root
	conn, _ := bases.LdapOpen(url, bind, pass)
	return conn
}

func annClose(conn any) {
	bases.LdapClose(conn)
}

// -------------------
