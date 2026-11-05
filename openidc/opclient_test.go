package openidc

import (
	"context"
	"crypto/tls"
	storage "main/openidc/Storage"
	"main/tools"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

/*
	Configuration to OIDC Test Provider
		oidc test server   "https://oidctest.wsweet.org/"   id "private",  secret "tardis"
		scopes   "openid", "profile", "email", "roles"
		users:pwd    dwho:dwho  rtyler:rtyler   or  msmith:msmith
*/

// Commun package
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestOidcClient(t *testing.T) {
	//tools.Recover()

	// OPClient needs config
	err := tools.LoadAnyConfig("data/test/wsweetoidc.yml", &tools.Cfg)
	tools.TEcheck(t, err, "Load Any Config")

	// Serveur web à lancer
	srv := &http.Server{
		Addr:    tools.Cfg.Server.ListenIp + ":" + tools.Cfg.Server.Port,
		Handler: http.DefaultServeMux,
	}

	// Func to test ( see config.yml )
	OpClient()

	// Déclenche le test dans 1 s
	time.AfterFunc(1*time.Second, func() {

		// Get Login path + redirect to sso  /authorize
		resp, err := http.Get("http://localhost:8080/client/login")
		tools.TEcheck(t, err, "GET failed: %v", err)
		defer resp.Body.Close()

		// Reponse : 200 from sso  with cookie lemonldappdata
		tools.Tcheck(t, resp.StatusCode == http.StatusOK, "expected 200, got %d", resp.StatusCode)

		// with cookie lemonldappdata
		c := resp.Header.Values("Set-Cookie")[0]
		tools.Tcheck(t, strings.HasPrefix(c, "lemonldappdata="), "expected Set-Cookie lemonldappdata=")

		// extract token from cookie issuerRequestoauth2%22%3A%22 ( 1773411293_61837 )  %22
		i := strings.Index(c, "issuerRequestoauth2%22%3A%22")
		tools.Tcheck(t, i+28+16 < len(c), "expected issuerRequestoauth2 token")

		// Get token : 1773411293_61837
		c = c[i+28 : i+28+16]

		// Simulate Received form data from sso login page  username dwho
		data := "inProgress=1&url=aHR0cHM6Ly9vaWRjdGVzdC53c3dlZXQub3JnL29hdXRoMg%3D%3D&timezone=1&skin=oidctest&token=1773411293_27801&user=dwho&password=dwho"
		data = strings.Replace(data, "1773411293_27801", c, 1)
		resp, err = http.Post(
			"https://oidctest.wsweet.org/oauth2/authorize?client_id=private&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fopenidc%2Fcallback&response_mode=query&response_type=code&scope=openid+profile+email+roles&state=%2Fstate",
			"application/x-www-form-urlencoded",
			strings.NewReader(data),
		)
		// Reponse : 200 from sso
		tools.Tcheck(t, resp.StatusCode == http.StatusOK, "expected 200, got %d", resp.StatusCode)

		tools.Info("test web end.")

		// Arret
		srv.Shutdown(context.Background())
	})

	// Serveur web
	tools.Info("server web start...")
	srv.ListenAndServe()
	tools.Info("server web end.")
}

// Require : openldap running

func TestOidcServer(t *testing.T) {

	// OPClient needs config
	err := tools.LoadAnyConfig("data/test/locsvroidc.yml", &tools.Cfg)
	tools.TEcheck(t, err, "Load Any Config")

	ustore := storage.NewUserStore("https://localhost:8080/oidc/")
	tools.Tcheck(t, ustore != nil, "Setup OIDC StoreFromFile failed")

	_, err = SetupServer(storage.NewStorage(ustore))
	tools.TEcheck(t, err, "Setup OIDC Server failed: %v", err)

	// Serveur web à lancer
	srv := &http.Server{
		TLSConfig: &tls.Config{MinVersion: tls.VersionTLS12},
		Addr:      tools.Cfg.Server.ListenIp + ":" + tools.Cfg.Server.Port,
		Handler:   http.DefaultServeMux,
	}
	// Déclenche le test dans 1 s
	time.AfterFunc(1*time.Second, func() {

		// Get Login path + redirect to sso  /authorize
		resp, err := http.Get("https://localhost:8080/oidc/authorize?client_id=test&redirect_uri=https%3A%2F%2Flocalhost%3A8080%2Foidc%2Fcallback&response_mode=query&response_type=code&scope=openid+profile+email+roles&state=%2Fstate")
		tools.TEcheck(t, err, "GET failed: %v", err)
		defer resp.Body.Close()

		// Reponse : 404 from sso  (car pas de page de login enregistré pour ce test )
		tools.Tcheck(t, resp.StatusCode == http.StatusNotFound, "expected 404, got %d", resp.StatusCode)
		tools.Tcheck(t, resp.Request.URL.Path == "/oidc-ui/login.html", "login page not built")

		tools.Info("test web end.")

		// Arret
		srv.Shutdown(context.Background())
	})

	// Serveur web
	tools.Info("server web start...")
	srv.ListenAndServeTLS(tools.Cfg.Server.CertFile, tools.Cfg.Server.KeyFile)
	tools.Info("server web end.")
}
