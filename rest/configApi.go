package rest

import (
	"main/tools"
	"net/http"
)

type UiLoginCfg struct {
	Form     bool
	Oidc     bool
	OidcPath string
}

// REST call : retourne des élément de configuration
func apiConfig(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	if query.Has("login") {
		var cfg UiLoginCfg
		cfg.Form = tools.Cfg.WebUi.Login.Active
		cfg.Oidc = tools.Cfg.Openidc.Client.Active
		cfg.OidcPath = tools.Cfg.Openidc.Client.LoginPath
		w.WriteHeader(http.StatusOK)
		tools.Encode(w, cfg)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}
