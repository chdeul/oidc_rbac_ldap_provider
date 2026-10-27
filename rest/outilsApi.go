package rest

import (
	obj "main/objets"
	"main/projet/destinat"
	"main/tools"
	"net/http"
	"strings"
)

/*
// Crud Users
CreateUser(user *obj.User) error
UpdateUser(user *obj.User) error
DeleteUser(user *obj.User) error

// Crud Roles (RBAC : Association Projet/Utilisateur -> Rôle)
CreateRole(prj *obj.Project, user *obj.User, role string) error
UpdateRole(prj *obj.Project, user *obj.User, , role string) error
DeleteRole(prj *obj.Project, user *obj.User, , role string) error
*/

func openOutil(url, prj string) destinat.Destination {

	mots := strings.Split(url, "/")
	if len(mots) > 0 {
		outil := mots[len(mots)-1]
		app := destinat.DestFactory(outil)
		if app == nil {
			tools.Panic("Mauvais outil dans outilsApi : %s", outil)
		}
		return app
	} else {
		tools.Panic("Mauvaise url dans outilsApi : %s", url)
	}
	return nil
}

// GET /outils/gitlab?project=xxxx&user=yyyy&role=zzz
func apiOutilsRead(w http.ResponseWriter, r *http.Request) {

	var user *obj.User
	var project *obj.Project

	prj := r.URL.Query().Get("project")
	path := r.URL.Path
	app := openOutil(path, prj)
	defer app.Close()

	usr := r.URL.Query().Get("user")
	rol := r.URL.Query().Get("role")
	if prj == "" && usr == "" {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	if prj != "" {
		tabp := app.SearchProjects(prj)
		if usr == "" && rol == "" {
			tools.Encode(w, tabp)
			return
		} else if len(tabp) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		} else if len(tabp) > 1 {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
		project = tabp[0]
	}
	if usr != "" {
		tabu := app.SearchUsers(usr)
		if prj == "" {
			tools.Encode(w, tabu)
			return
		} else if len(tabu) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		} else if len(tabu) > 1 {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
		user = tabu[0]
	}
	if project == nil || user == nil {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	if rol == "" {
		rol = "*"
	}
	tabr := app.GetRole(project, user)
	tools.Encode(w, tabr)
}

func apiOutilsCreate(w http.ResponseWriter, r *http.Request) {

}

func apiOutilsUpdate(w http.ResponseWriter, r *http.Request) {

}

func apiOutilsDelete(w http.ResponseWriter, r *http.Request) {

}
