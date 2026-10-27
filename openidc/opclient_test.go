package openidc

import (
	"context"
	"fmt"
	storage "main/openidc/Storage"
	"main/tools"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestOidcClient(t *testing.T) {

	/*
		Configuration to OIDC Test Provider
			oidc test server   "https://oidctest.wsweet.org/"   id "private",  secret "tardis"
			scopes   "openid", "profile", "email", "roles"
			users:pwd    dwho:dwho  rtyler:rtyler   or  msmith:msmith
	*/
	// OPClient needs config
	tools.LoadConfig("..")

	// Serveur web à lancer
	srv := &http.Server{
		Addr:    ":8008",
		Handler: http.DefaultServeMux,
	}

	// Func to test  ( see config.yml )
	OpClient()

	// Déclenche le test dans 1 s
	time.AfterFunc(1*time.Second, func() {

		// Get Login path + redirect to sso  /authorize
		resp, err := http.Get("http://localhost:8008/login")
		if err != nil {
			t.Fatalf("GET failed: %v", err)
		}
		defer resp.Body.Close()

		// Reponse : 200 from sso  with cookie lemonldappdata
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
		// with cookie lemonldappdata
		c := resp.Header.Values("Set-Cookie")[0]
		if !strings.HasPrefix(c, "lemonldappdata=") {
			t.Errorf("expected Set-Cookie lemonldappdata=")
		}
		// extract token from cookie issuerRequestoauth2%22%3A%22 ( 1773411293_61837 )  %22
		i := strings.Index(c, "issuerRequestoauth2%22%3A%22")
		if i+28+16 >= len(c) {
			t.Errorf("expected issuerRequestoauth2 token")
		}
		// Get token : 1773411293_61837
		c = c[i+28 : i+28+16]

		// Simulate Received form data from sso login page  username dwho
		data := "inProgress=1&url=aHR0cHM6Ly9vaWRjdGVzdC53c3dlZXQub3JnL29hdXRoMg%3D%3D&timezone=1&skin=oidctest&token=1773411293_27801&user=dwho&password=dwho"
		data = strings.Replace(data, "1773411293_27801", c, 1)
		resp, err = http.Post(
			"https://oidctest.wsweet.org/oauth2/authorize?client_id=private&redirect_uri=http%3A%2F%2Flocalhost%3A8008%2Fopenidc%2Fcallback&response_mode=query&response_type=code&scope=openid+profile+email+roles&state=%2Fstate",
			"application/x-www-form-urlencoded",
			strings.NewReader(data),
		)
		// Reponse : 200 from sso
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
		// with callback called
		//if received == nil {
		//	t.Errorf("expected callback called")
		//}
		// Receive {"email":"dwho@badwolf.org","name":"Doctor Who","preferred_username":"dwho","roles":["users","timelords"],"sid":"***","sub":"dwho"}
		//if received.Name != "Doctor Who" {
		//	t.Errorf("expected name %s, received %s", "Doctor Who", received.Name)
		//}
		fmt.Println("test web end.")

		// Arret
		srv.Shutdown(context.Background())
	})

	// Serveur web
	fmt.Println("server web start...")
	srv.ListenAndServe()
	fmt.Println("server web end.")
}

func TestOidcServer(t *testing.T) {

	// OPServer needs config
	tools.LoadConfig(".")

	ustore := storage.NewUserStore("../data/users.json")
	if ustore == nil {
		t.Fatalf("Setup OIDC StoreFromFile failed")
	}
	_, err := SetupServer(storage.NewStorage(ustore))
	if err != nil {
		t.Fatalf("Setup OIDC Server failed: %v", err)
	}
	// Serveur web à lancer
	srv := &http.Server{
		Addr:    ":8008",
		Handler: nil,
	}
	// Déclenche le test dans 1 s
	time.AfterFunc(1*time.Second, func() {

		// Get Login path + redirect to sso  /authorize
		// /authorize?client_id=sid1&redirect_uri=http%3A%2F%2Flocalhost%3A8008%2Fopenidc%2Fcallback&response_mode=query&response_type=code&scope=openid+profile+email+roles&state=%2Fstate
		resp, err := http.Get("http://localhost:8008/authorize?client_id=sid1&redirect_uri=http%3A%2F%2Flocalhost%3A8008%2Fopenidc%2Fcallback&response_mode=query&response_type=code&scope=openid+profile+email+roles&state=%2Fstate")
		if err != nil {
			t.Fatalf("GET failed: %v", err)
		}
		defer resp.Body.Close()

		// Reponse : 200 from sso
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}

		fmt.Println("test web end.")

		// Arret
		srv.Shutdown(context.Background())
	})

	// Serveur web
	fmt.Println("server web start...")
	srv.ListenAndServe()
	fmt.Println("server web end.")
}
