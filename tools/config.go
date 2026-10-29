package tools

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"go.yaml.in/yaml/v4"
)

// Structure de données de config
type Config struct {
	Logging string
	LogTime bool

	Batch struct {
		Active bool
	}
	Server struct {
		Protocol    string
		ListenIp    string
		Port        string
		CertFile    string
		KeyFile     string
		UsedUrl     string
		BackUrl     string
		TokenSecret string
		Smtp        struct {
			Ident    string
			Password string
			Host     string
			Port     string
		}
	}
	Ref struct {
		Type    string
		Path    string
		Filters map[string]string
	}
	Dest struct {
		Ldap struct {
			Url    string
			Root   string
			Passwd string
		}
	}
	WebUi struct {
		Active      bool
		ContextPath string
		Cache       bool
		CsvApi      bool
		Login       struct {
			Active    bool
			AgeCookie string
		}
	}
	Openidc struct {
		Client struct {
			Active      bool
			Issuer      string
			Ident       string
			Secret      string
			Callback    string
			PostLogout  string
			FrontLogout string
			Scope       []string
			AuthMethod  string
			UidField    string
			LoginPath   string
			AgeCookie   string
			CreateUser  bool
			RootCA      string
		}
		Server struct {
			Active      bool
			ContextPath string
			SsoUiPath   string
			FormLogin   string
			LogoutPath  string

			Registered []struct {
				Ident       string
				Secret      string
				Callback    string
				PostLogout  string
				FrontLogout string
			}
		}
	}
}

// Config globale chargée
var Cfg Config

// Recuperer du json
func Decode(r *http.Request, obj interface{}) {
	json.NewDecoder(r.Body).Decode(&obj)
}

// Envoyer du json
func Encode(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func LoadConfig(dir string) {

	// Lecture fichier
	err := LoadAnyConfig(dir+"/config.yml", &Cfg)
	if err != nil {
		panic(err)
	}

	// Normalisation des context path configuré  "/"" => ""   "/toto" => "/toto"
	if Cfg.WebUi.ContextPath == "/" {
		Cfg.WebUi.ContextPath = ""
	}
	if Cfg.Openidc.Server.ContextPath == "/" {
		Cfg.Openidc.Server.ContextPath = ""
	}
	if Cfg.Openidc.Server.SsoUiPath == "/" {
		Cfg.Openidc.Server.SsoUiPath = ""
	}
	// Pas de date et heure sur les log et coloré (dev)
	if !Cfg.LogTime {
		LogWithColor(true)
		log.SetFlags(0)
	}
	// Positionne le niveau de log
	SetLogLevel(Cfg.Logging)
}

// Charge une configuration yaml
func LoadAnyConfig(filepath string, ptr any) error {

	data, err := os.ReadFile(filepath)
	if err != nil {
		Error("Fichier config %s non trouvé", filepath)
		return err
	}
	if err := yaml.Unmarshal(data, ptr); err != nil {
		Error("Fichier config erroné : %s", err)
		return err
	}
	return err
}
