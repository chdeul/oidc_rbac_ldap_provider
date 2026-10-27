package openidc

import (
	"context"
	"crypto/sha256"
	"fmt"
	"main/interf"
	storage "main/openidc/Storage"
	"main/tools"
	"net/http"
	"net/url"
	"strings"

	"github.com/zitadel/oidc/v3/pkg/op"
	"golang.org/x/text/language"
)

var store *storage.Storage
var provider *op.Provider

// Gère les Form submit POST  de oidc-ui/auth/login et renew
func formManagement(w http.ResponseWriter, r *http.Request) {

	// Check
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Data of form submitted
	r.ParseForm()
	identifiant := r.FormValue("Identifiant") // Ou Email
	password := r.FormValue("MotDePasse")
	email := ""
	if strings.Contains(identifiant, "@") {
		email = identifiant // oui Email
		identifiant = ""
	}
	requestId := r.FormValue("RequestId")
	uri := tools.Cfg.Server.UsedUrl + tools.Cfg.Openidc.Server.SsoUiPath

	// Requete de Login sur le SSO
	if strings.Contains(r.URL.Path, "/auth/login") {

		// get identifiant from email
		if identifiant == "" {
			u := interf.AnnReadUserFromEmail(email)
			if u != nil {
				identifiant = u.Identifiant
			}
		}
		var err error
		if identifiant == "" || password == "" || requestId == "" {
			err = fmt.Errorf("Login Form data incorrectes")
		} else {
			// Login verification  des données  et check requestID
			err = store.CheckUsernamePassword(identifiant, password, requestId)
		}
		if err == nil {
			// Ok , retourne sur l'application client callback
			callback := op.AuthCallbackURL(provider)
			path := tools.Cfg.Openidc.Server.ContextPath + callback(context.Background(), requestId)
			http.Redirect(w, r, path, http.StatusFound)
		} else {
			// Retourne l'erreur que l'interface web html
			url := uri + tools.Cfg.Openidc.Server.FormLogin
			query := "?requestId=" + requestId + "&error=Identifiant ou Mot de passe incorrect"
			http.Redirect(w, r, url+query, http.StatusFound)
		}
		return
	}

	// Renouvellement mot de passe SSO
	if strings.Contains(r.URL.Path, "/auth/renew") {

		// Url de retour dans le mail, qui sera appelée avec la clef
		url := uri + tools.Cfg.Openidc.Server.FormLogin + "?requestId="

		// Renouvellement mot de passe 2 étapes : send mail, puis setpassword
		msg, _ := interf.RenewPassword(identifiant, password, email, requestId, url)
		url = uri + tools.Cfg.Openidc.Server.LogoutPath + "?msg=" + msg
		http.Redirect(w, r, url, http.StatusFound)
		return
	}

	// Requete non prévue, iconnue
	fmt.Println("Requête inconnue", r.URL.Path)
	w.WriteHeader(http.StatusNoContent)
}

// We use the login UI and pass the (auth request) id
func loginURLMaker(id string) string {
	return tools.Cfg.Openidc.Server.SsoUiPath + tools.Cfg.Openidc.Server.FormLogin + "?requestId=" + id
}

// We use the login UI and logout page
func logoutURLMaker() string {

	// inject list of logout urls of all clients
	list := ""
	for _, cl := range tools.Cfg.Openidc.Server.Registered {
		if cl.FrontLogout != "" {
			list = list + "|" + cl.FrontLogout
		}
	}
	// Logout url of sso
	query := "?urls=" + url.QueryEscape(list)
	return tools.Cfg.Openidc.Server.SsoUiPath + tools.Cfg.Openidc.Server.LogoutPath + query
}

// Enregistre les clients Oidc, dont moi même
func registerClients() {

	// Clients from config files
	for _, cl := range tools.Cfg.Openidc.Server.Registered {
		client := storage.WebClient(cl.Ident, cl.Secret, loginURLMaker, cl.Callback, cl.PostLogout, cl.FrontLogout)
		storage.RegisterClients(client)
	}
}

// Register a new client in storage/clients.go to use it
func SetupServer(storage *storage.Storage) (*op.Provider, error) {

	// Glbal variable
	store = storage

	// Clients from config file
	registerClients()

	// the OpenID Provider requires a 32-byte key for (token) encryption
	key := sha256.Sum256([]byte("$agréments!"))

	config := &op.Config{

		CryptoKey: key,

		// will be used if the end_session endpoint is called without a post_logout_redirect_uri
		DefaultLogoutRedirectURI: logoutURLMaker(),

		//enables code_challenge_method S256 for PKCE (and therefore PKCE in general)
		CodeMethodS256: true,

		// enables additional client_id/client_secret authentication by form post (not only HTTP Basic Auth)
		AuthMethodPost: true,

		// enables additional authentication by using private_key_jwt
		AuthMethodPrivateKeyJWT: true,

		// enables refresh_token grant use
		GrantTypeRefreshToken: true,

		// enables use of the `request` Object parameter
		RequestObjectSupported: true,

		// we'll set the here accordingly
		SupportedUILocales: []language.Tag{language.French},

		SupportedClaims: []string{"openid", "profile", "email", "role"},
		SupportedScopes: []string{"sub", "email", "preferred_username", "name", "roles"},
	}
	// Cree le serveur OIDC
	var err error
	path := tools.Cfg.Openidc.Server.ContextPath
	providerPath := tools.Cfg.Server.BackUrl + path + "/"
	provider, err = op.NewProvider(config, store, op.StaticIssuer(providerPath), op.WithAllowInsecure())
	if err != nil {
		panic(err)
	}
	// Http Serve les requetes oidc standars
	//	/healthz	/ready	/.well-known/openid-configuration	/oauth/token	/oauth/introspect
	//	/callback	/authorize	/userinfo	/revoke		/end_session	/keys	/device_authorization
	http.Handle(path+"/", http.StripPrefix(path, provider))
	return provider, err
}

func SetupSsoUi() error {

	// Http serve static resources oidc-ui/* ( html, js, css, img, fonts )
	fs := http.FileServer(http.Dir("oidc-ui/"))
	path := tools.Cfg.Openidc.Server.SsoUiPath
	http.Handle(path+"/", http.StripPrefix(path+"/", fs))

	// Http Serve le submit form  /oidc-ui/auth/login
	http.HandleFunc(path+"/auth/", formManagement)
	return nil
}
