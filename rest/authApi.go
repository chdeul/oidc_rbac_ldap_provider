package rest

import (
	"main/interf"
	obj "main/objets"
	"main/tools"
	"net/http"
)

const usersfile string = "./data/users.csv"

type formData struct {
	Identifiant string
	MotDePasse  string
	RequestId   string
}

type profilData struct {
	Identifiant string
	Profil      string
}

// Retourne l'utilisateur correspondant au token, si le profil est autorisé
func checkUtilisateur(token string, profils *string) *obj.User {

	// Extraire Utilisateur et Profil du jeton
	identifiant, profil, _ := tools.ParseToken(token)
	if identifiant == "" {
		// Jeton expiré ou Plus logué
		return nil
	}
	// Bon profil ?
	if obj.ProfilInclus(*profils, profil) {

		// Retourne l'utilisateur
		return &obj.User{Identifiant: identifiant, MotDePasse: token}
	}
	// Pas du bon profil
	return nil
}

// REST call : login this user , retourne le token
func apiLogin(w http.ResponseWriter, r *http.Request) {

	// Utilisateur entrant
	var data formData
	tools.Decode(r, &data)

	// Cherche dans Table des utilisateurs
	user := interf.AnnReadUser(data.Identifiant)

	// Trouvé et password ok
	if user != nil && interf.AnnCheckPassword(data.Identifiant, data.MotDePasse) {

		// Génère un (nouveau) jeton
		profil := obj.GetProfil(user.Identifiant, user.Roles)
		token, err := tools.JwtCreateAndTransmit(w, user.Identifiant, "", profil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Retourne Aussi le jeton dans le mot de passe de l'User
		user.MotDePasse = token
		w.WriteHeader(http.StatusOK)
		tools.Encode(w, profilData{Identifiant: user.Identifiant, Profil: profil})
		return
	}
	// Non trouvé : accès interdit
	w.WriteHeader(http.StatusUnauthorized)
}

// REST Call : logout user from token
func apiLogout(w http.ResponseWriter, r *http.Request) {

	// Jeton d'auth présent Dans Header ou Cookie
	jwtToken := tools.ExtractJwt(r)
	_, _, ssotoken := tools.ParseToken(jwtToken)

	// Oidc logout
	oidcurl := ""
	if ssotoken != "" {
		oidcurl = tools.GetOidcLogoutUrl()
		oidcurl = oidcurl + "?id_token_hint=" + ssotoken

		// On ne reviens pas si logout sso
		// oidcurl = oidcurl + "&state=789456123" // Unused None
		// oidcurl = oidcurl + "&post_logout_redirect_uri=" + tools.Cfg.Openidc.Client.PostLogout  // None
	}

	// le jeton sera supprimé côté client (stateless app)
	tools.JwtRemoveAndTransmit(w)
	w.WriteHeader(http.StatusOK)

	// la redirection vers sso logout sera côté client
	tools.Encode(w, oidcurl)
}

// REST Call : retourne l'utilisateur ou vide
func apiLogged(w http.ResponseWriter, r *http.Request) {

	// Recupere l'utilisateur de la requete (du token)
	user := r.Context().Value(myReqCtx{}).(agremReqCtx).user
	if user != nil {
		// renvoie l'Utilisateur sans mot de passe, sans email
		profil := obj.GetProfil(user.Identifiant, user.Roles)
		w.WriteHeader(http.StatusOK)
		tools.Encode(w, profilData{Identifiant: user.Identifiant, Profil: profil})
	} else {
		// Pas d'utilisateur logué
		w.WriteHeader(http.StatusNoContent)
	}
}

// REST Call : Envoie ou traite un mail de reinit du mot de passe
func apiRenew(w http.ResponseWriter, r *http.Request) {

	// Utilisateur entrant
	var data obj.User
	tools.Decode(r, &data)

	// Url indiqué dans le mail, appelée en retour avec la clef
	url := tools.Cfg.Server.UsedUrl + tools.Cfg.WebUi.ContextPath
	url = url + "/vue/renew?clef=" // from webui

	// Do the job : send mail or change password  (Clef dans NomComplet)
	msg, _ := interf.RenewPassword(data.Identifiant, data.MotDePasse, data.AdresseMail, data.NomComplet, url)
	w.WriteHeader(http.StatusOK)
	tools.Encode(w, msg)
}
