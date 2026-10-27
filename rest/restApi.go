package rest

import (
	"context"
	obj "main/objets"
	"main/tools"
	"net/http"
	"regexp"
	"strings"
)

type myReqCtx struct {
}

// Données contexte de requête
type agremReqCtx struct {
	path []string  // Url et param
	user *obj.User // Nil ou Logué
}

// Une route d'entrée d'API
type route struct {
	method  string           // GET POST PUT DELETE ...
	regex   *regexp.Regexp   // Url path with matching elements
	handler http.HandlerFunc // http func handler
	profils *string          // Auth profils required
}

// Tableau des routes / handler methods
var routes = []route{

	// REST Configuration
	Route("GET", "config", apiConfig, &obj.Publique),

	// REST Authentification
	Route("GET", "auth", apiLogged, &obj.Publique),
	Route("POST", "auth", apiLogin, &obj.Publique),
	Route("PUT", "auth", apiRenew, &obj.Publique),
	Route("DELETE", "auth", apiLogout, &obj.Publique),

	// REST Ldap agrements users and groups : browse, edit
	Route("GET", "ldap", apiLdapRead, &obj.Consultation),
	Route("POST", "ldap", apiLdapCreate, &obj.Administration),
	Route("PUT", "ldap", apiLdapUpdate, &obj.Administration),
	Route("DELETE", "ldap", apiLdapDelete, &obj.Administration),

	// REST Annuaire Ref users and group : browse
	Route("GET", "annref", apiAnnRefRead, &obj.Publique),

	// REST Outils destination agrements
	Route("GET", "outil/.*", apiOutilsRead, &obj.Consultation),
	Route("POST", "outil/.*", apiOutilsCreate, &obj.Administration),
	Route("PUT", "outil/.*", apiOutilsUpdate, &obj.Administration),
	Route("DELETE", "outil/.*", apiOutilsDelete, &obj.Administration),

	// REST Projet: prise en compte
	Route("GET", "projet", apiProjet, &obj.Edition),
}

// Declare une route
func Route(method, pattern string, handler http.HandlerFunc, profils *string) route {
	return route{method, regexp.MustCompile("^" + pattern + "$"), handler, profils}
}

// Extrait un element (matching) de la route
func getField(r *http.Request, index int) string {
	fields := r.Context().Value(myReqCtx{}).(agremReqCtx).path
	return fields[index]
}

// Aiguille sur la bonne route
func ServeRestApi(w http.ResponseWriter, r *http.Request) {

	// Catch rest errors
	defer func() {
		if r := recover(); r != nil {
			tools.Error("Exception dans  ServeRestApi %v", r)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()

	// Pour message erreur final
	var allow []string

	// Cherche parmi les routes
	for _, route := range routes {

		// Est-ce que cela match ?
		matches := route.regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {

			// Trouvé, mais pas bon verbe
			if r.Method != route.method {
				// Pour message erreur final
				allow = append(allow, route.method)
				continue
			}
			// Contexte de requete : avec url
			reqCtx := agremReqCtx{path: matches[1:]}

			// Jeton d'auth présent Dans Header ou Cookie
			jwtToken := tools.ExtractJwt(r)
			if jwtToken != "" {
				// Decode le jeton, et teste si match le profil requis pour cette url
				reqCtx.user = checkUtilisateur(jwtToken, route.profils)
			}
			// Trouvé avec bon verbe, verif si profil requis
			if *route.profils != obj.Publique {

				// Jeton d'auth non présent
				if jwtToken == "" {
					// Pas de jeton , Pas Authentifie
					http.Error(w, "Accès non authentifié", http.StatusUnauthorized)
					return
				}
				if reqCtx.user == nil {
					// Authentifie mais Pas les bons droits
					http.Error(w, "Vous n'avez pas les droits", http.StatusForbidden)
					return
				}
			}
			// aiguille vers la methode, avec le contexte de requete
			ctx := context.WithValue(r.Context(), myReqCtx{}, reqCtx)
			route.handler(w, r.WithContext(ctx))
			return
		}
	}
	if len(allow) > 0 {
		// Trouvé mauvais verbe
		w.Header().Set("Allow", strings.Join(allow, ", "))
		http.Error(w, "Verbe non autorisé", http.StatusMethodNotAllowed)
		return
	}
	// Url Non trouvée
	http.NotFound(w, r)
}
