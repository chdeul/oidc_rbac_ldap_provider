package main

import (
	"context"
	"crypto/tls"
	"embed"
	_ "embed"
	"io/fs"
	"main/interf"
	"main/openidc"
	storage "main/openidc/Storage"
	"main/rest"
	"main/tools"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// static web server content embed in binary.
//
//go:embed resources/*
var resources embed.FS

func main() {

	// ---- Initialisations générales ----
	tools.Info("Début Agréments")

	// TODO Internationnalisation English

	// Initialisation  Config
	//time.LoadLocation("Europe/Paris")
	tools.LoadConfig(".")

	// Config envoi Mail
	cs := tools.Cfg.Server.Smtp
	interf.SetupSmtp(cs.Host, cs.Port, cs.Ident, cs.Password)

	// ---- Batchs ----

	// Traitements batchs, TODO cron  batchs
	if tools.Cfg.Batch.Active {

	}

	// ---- Server OIDC  /oidc/*  ----
	if tools.Cfg.Openidc.Server.Active {

		// Users Annuary for OIDC  : on Agréments openldap
		users := storage.NewUserStore("")
		defer storage.CloseStore(users)

		// Store for OIDC features (ldap + config file + memory)
		store := storage.NewStorage(users)

		// Oidc protocol : /oidc/authorize	/oidc/userinfo etc
		openidc.SetupServer(store)

		// Oidc UI  : sso login page
		openidc.SetupSsoUi()
	}

	// ---- Agréments Web interface / Administration  / Api  ----
	if tools.Cfg.WebUi.Active {

		// Context path de la partie appli web  /agrem/*
		path := tools.Cfg.WebUi.ContextPath

		// Http serve static resources ( js, css, img, fonts )
		var fsr http.Handler
		if tools.Cfg.WebUi.Cache {
			sfs, e := fs.Sub(resources, "resources")
			if e != nil {
				tools.Fatal("FS Server fail :  %v", e)
			}
			fsr = http.FileServer(http.FS(sfs))
		} else {
			fsr = http.FileServer(http.Dir("resources"))
		}
		http.Handle(path+"/resources/", http.StripPrefix(path+"/resources/", fsr))

		// Http serve Rest API  /api/*     (include native /api/auth/login & logout)
		http.Handle(path+"/api/", http.StripPrefix(path+"/api/", http.HandlerFunc(rest.ServeRestApi)))

		// Http server refresh and favoris of browser on /vue/* routes
		http.Handle(path+"/vue/", http.StripPrefix(path+"/vue/", http.HandlerFunc(ServeVue)))

		// Point d'entree, Accueil ( single page, vues ) et autres  /  => /index.html
		http.Handle(path+"/", http.StripPrefix(path+"/", http.HandlerFunc(ServeIndex)))
		if path != "" {
			// Wrong context path
			http.Handle("/", http.StripPrefix("/", http.HandlerFunc(ServeNone)))
		}
	}

	// ---- Client OIDC  /client/login /client/callback (oidc login & logout ) ----
	if tools.Cfg.Openidc.Client.Active {
		time.AfterFunc(1*time.Second, func() {

			// -----  Enregistré 1 s apres le démarrage server oidc
			openidc.OpClient()
		})
	}

	// ---- Start composite http Server  : Oidc Server,  Oidc-UI,  Oidc Client,  Web-Ui ----
	if tools.Cfg.Openidc.Server.Active || tools.Cfg.Openidc.Client.Active || tools.Cfg.WebUi.Active {

		// Serveur HTTP
		svr := &http.Server{
			TLSConfig: &tls.Config{MinVersion: tls.VersionTLS12},
			Addr:      tools.Cfg.Server.ListenIp + ":" + tools.Cfg.Server.Port,
			Handler:   http.DefaultServeMux,
		}
		// Sur Ctrl C ou kill
		go detecterArret(svr)

		if tools.Cfg.Server.Protocol == "https" {
			// https : en production pour transfert mots de passe crypté
			err := svr.ListenAndServeTLS(tools.Cfg.Server.CertFile, tools.Cfg.Server.KeyFile)
			if err != nil {
				tools.Fatal("TLS Server fail to start :  %v", err)
			}
		} else {
			// Http pour test
			err := svr.ListenAndServe()
			if err != nil {
				tools.Fatal("Server fail to start :  %v", err)
			}
		}
	}
	tools.Info("Fin Agréments")
}

func detecterArret(svr *http.Server) {

	// Capture CTRL+C ou kill -SIGINT
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	tools.Info("Serveur web prêt : %s", tools.Cfg.Server.BackUrl)

	<-stop // On attend le signal

	tools.Info("Arrêt demandé...")

	// Timeout pour l'arrêt propre
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := svr.Shutdown(ctx); err != nil {
		tools.Fatal("Arrêt forcé: %v", err)
	}
}
