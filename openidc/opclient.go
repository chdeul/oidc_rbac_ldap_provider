package openidc

import (
	"context"
	"main/interf"
	obj "main/objets"
	"main/tools"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/zitadel/oidc/v3/pkg/client/rp"
	httphelper "github.com/zitadel/oidc/v3/pkg/http"
	"github.com/zitadel/oidc/v3/pkg/oidc"
)

var key = []byte("addef187ef1783bc46abc906h2d4z6f9")

var client rp.RelyingParty

// Options handler oidc client
var urlOptions = []rp.URLParamOpt{
	//rp.WithPromptURLParam("Welcome back!"),
	rp.WithResponseModeURLParam(oidc.ResponseMode(oidc.ResponseModeQuery)),
}

// Start login procedure
func clientLogin(w http.ResponseWriter, r *http.Request) {

	// Avec query path referer
	p := r.URL.Query().Get("redirect")
	theState := func() string {
		// Un simple multiple de 7 pour anti csrf
		alea := rand.Intn(65535*65535/7) * 7
		// plus path initial de redirection final
		return strconv.Itoa(alea) + "~" + p
	}
	// passe Client oidc qui va rediriger vers Server oidc
	ah := rp.AuthURLHandler(theState, client, urlOptions...)
	ah.ServeHTTP(w, r)
}

// login from oidc client
func clientCallback(w http.ResponseWriter, r *http.Request,
	tokens *oidc.Tokens[*oidc.IDTokenClaims],
	state string, rp rp.RelyingParty, info *oidc.UserInfo) {

	// state : jeton antcsrf , Url initiale demandée
	path := tools.Cfg.WebUi.ContextPath + "/index.html" // default redirection home page
	values := strings.SplitN(state, "~", 2)
	if len(values) != 2 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Verifier anti csrf : multiple de 7
	alea, err := strconv.Atoi(values[0])
	if err != nil || (alea%7) != 0 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Mémorise pour logout
	tools.SetOidcLogoutUrl(rp.GetEndSessionEndpoint())
	idtoken := tokens.IDToken

	// Search or Create user into agrem ldap from OIDC Infos
	user := obj.User{Uid: info.Subject, Identifiant: info.PreferredUsername, NomComplet: info.Name, AdresseMail: info.Email, Roles: info.Claims["Roles"].(string)}
	err = interf.CheckOrCreateUserFromOidcLogin(&user)
	if err != nil {
		http.Redirect(w, r, path, http.StatusForbidden)
		return
	}

	// Génère un (nouveau) jeton d'uthentification JWT
	_, err = tools.JwtCreateAndTransmit(w, user.Identifiant, idtoken, obj.GetProfil(user.Identifiant, user.Roles))
	if err != nil {
		http.Redirect(w, r, path, http.StatusInternalServerError)
		return
	}
	// Redirige vers la page d'accueil de l'appli web Agrements : Authentifié JWT
	http.Redirect(w, r, path, http.StatusFound)
}

// logout from oidc server
func clientLogout(w http.ResponseWriter, r *http.Request) {

	// le jeton sera supprimé côté client (stateless app)
	tools.JwtRemoveAndTransmit(w)
	w.WriteHeader(http.StatusOK)
}

func OpClient() error {

	// WithJWTProfile option
	//keyID := "changeit"
	//signingKey := []byte{cert.pem}

	// Option oidc client
	cookieHandler := httphelper.NewCookieHandler(key, key, httphelper.WithUnsecure())

	options := []rp.Option{
		// rp.With ....
		rp.WithCookieHandler(cookieHandler),
		rp.WithVerifierOpts(rp.WithIssuedAtOffset(5 * time.Second)),
		rp.WithSigningAlgsFromDiscovery(),
		//rp.WithJWTProfile(rp.SignerFromKeyAndKeyID(signingKey, keyID)),
	}

	// Trusted https certificates
	ctx := tools.LoadTrustedCertFiles(context.Background())

	// Register as oidc client ( call server oidc .well-known/openid-configuration )
	var err error
	client, err = rp.NewRelyingPartyOIDC(ctx, tools.Cfg.Openidc.Client.Issuer+"/",
		tools.Cfg.Openidc.Client.Ident, tools.Cfg.Openidc.Client.Secret,
		tools.Cfg.Openidc.Client.Callback, tools.Cfg.Openidc.Client.Scope, options...)
	if err != nil {
		panic(err)
	}

	// Register /client/login  (declenche le login sso)
	http.Handle(tools.Cfg.Openidc.Client.LoginPath, http.HandlerFunc(clientLogin))

	// Register the frontchannel_logout_uri   /client/logout  (sur logout sso)
	if tools.Cfg.Openidc.Client.FrontLogout != "" {
		u, _ := url.Parse(tools.Cfg.Openidc.Client.FrontLogout)
		http.Handle(u.Path+"/", http.HandlerFunc(clientLogout))
	}
	// Register Callback endpoint   /client/callback  (sur bon login sso)
	u, _ := url.Parse(tools.Cfg.Openidc.Client.Callback)
	http.HandleFunc(u.Path+"/", rp.CodeExchangeHandler(rp.UserinfoCallback(clientCallback), client))
	return nil
}
