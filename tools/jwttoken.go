package tools

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	// Config Private key, default
	secretkey = "ma_Clef!Secrète."
)

// Var globale pour client oidc local
var SsoLogoutUrl string

func ExtractToken(r *http.Request) string {

	// Jeton d'auth présent Dans Header ou Cookie
	token := r.Header.Get("Authorization")
	if token == "" {
		cook, err := r.Cookie("JWTAgremToken")
		if err == nil {
			token = cook.Value
		}
	}
	if token != "" {
		token = strings.Replace(token, "Bearer ", "", 1)
	}
	return token
}

// Agrement personal access token ( n months, no-info, for API, header  )
func CreateAgPat(ident, profil string, expdate uint64) (string, error) {

	// TODO  random and save -> Ldap -> with user and profil and exp date
	return "agpat-1234567890AbcDefGhij", nil
}

// Parse Agrement personal access token
func parseAgPat(tokenString string) (ident, profil, oidcToken string) {

	// TODO  random -> Ldap -> user, profil, out of date ?
	if strings.HasPrefix(tokenString, "agpat-") {

		// user, profil, oidc:nothing
		return "api_client", "Edition", ""
	}
	// no authorization
	return "", "", ""
}

// Agrement JWT token ( 12h max, visible info, signed, for web ui, cookie or header )
func CreateToken(ident, profil, oidcToken string) (string, error) {

	secret := Cfg.Server.TokenSecret
	if secret == "" {
		secret = secretkey
	}

	// Creation with expiration date, no refresh token
	maxAge := 60 * 60 * 10                                            // seconds => 10h => 1 jour ouvré
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{ // HMAC-SHA 256
		"identifiant": ident,                                                      // Username
		"profil":      profil,                                                     // Droits
		"ssotoken":    oidcToken,                                                  // empty  or  sso token
		"exp":         time.Now().Add(time.Duration(maxAge) * time.Second).Unix(), // Expiration
	})
	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString([]byte(secret))
}

func ParseToken(tokenString string) (ident, profil, oidcToken string) {

	// If Personal access token
	if strings.HasPrefix(tokenString, "agpat-") {
		return parseAgPat(tokenString)
	}

	secret := Cfg.Server.TokenSecret
	if secret == "" {
		secret = secretkey
	}

	// Decode
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	// Si mauvais
	if err != nil {
		// Utilisateur vide
		return "", "", ""
	}

	// Extrait les données
	claims, ok := token.Claims.(jwt.MapClaims)

	// Si Invalid ou Expiré
	if !ok || !token.Valid {
		// Utilisateur vide
		return "", "", ""
	}

	// Return data ( identifiant et profil de droits )
	return claims["identifiant"].(string), claims["profil"].(string), claims["ssotoken"].(string)
}

func JwtCreateAndTransmit(w http.ResponseWriter, identifiant string, idOidc string, profil string) (string, error) {

	// Génère un (nouveau) jeton
	token, err := CreateToken(identifiant, profil, idOidc)
	if err != nil {
		return "", err
	}
	// Url extérieure des outils
	parsed, err := url.Parse(Cfg.Server.UsedUrl)
	if err != nil {
		return "", err
	}
	isHttps := (parsed.Scheme == "https")
	domaine := parsed.Hostname()

	// Un cookie  valable pour tous les appli controlés : no context  /
	cook := http.Cookie{Name: "JWTAgremToken", Value: token, Domain: domaine, Path: "/", HttpOnly: false, Secure: isHttps}

	// Durée de session
	age := Cfg.WebUi.Login.AgeCookie
	if idOidc != "" {
		age = Cfg.Openidc.Client.AgeCookie
	}
	if age != "session" {
		cook.MaxAge, _ = strconv.Atoi(Cfg.WebUi.Login.AgeCookie)
	}

	// Retourne le jeton JWT D'authentification dans un cookie
	http.SetCookie(w, &cook)

	// Login ok,
	return token, err
}

func JwtRemoveAndTransmit(w http.ResponseWriter) {

	// Url extérieure des outils
	parsed, err := url.Parse(Cfg.Server.UsedUrl)
	if err != nil {
		return
	}
	isHttps := (parsed.Scheme == "https")
	domaine := parsed.Hostname()

	// Un cookie  valable pour tous les appli controlés : no context  /  ( pour le proxy)
	cook := http.Cookie{Name: "JWTAgremToken", Value: "", Domain: domaine, Path: "/", HttpOnly: false, Secure: isHttps, MaxAge: -1}

	// Effacer   Token: ""  MaxAge: -1
	http.SetCookie(w, &cook)
}

// For SSO client local login/logout ( opclient / rest )
func SetOidcLogoutUrl(url string) {
	SsoLogoutUrl = url
}

func GetOidcLogoutUrl() string {
	return SsoLogoutUrl
}
