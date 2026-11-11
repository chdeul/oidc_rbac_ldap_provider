package rest

import (
	obj "main/objets"
	"main/projet"
	"main/tables"
	"main/tools"
	"net/http"
)

//  serve  GET POST  /agrem/projet  :  prise en compte

// REST call : search users, return json list
func apiProjet(w http.ResponseWriter, r *http.Request) {

	// Read users directory
	initAnnRef()

	switch r.Method {

	case http.MethodGet:
		// Recupere le nom du projet
		name := r.URL.Query().Get("name")
		if !tools.Check(name, "No name in query of project") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// Return json table of project
		tab := projet.LoadCsvProject(name)
		tools.Encode(w, tab)

	case http.MethodPut:
		// recupere une url en query
		url := r.URL.Query().Get("data")
		if url == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// recupere le csv par get sur cette url
		resp, err := http.Get(url)
		if !tools.CheckError(err, "GET project table failed: %v", err) {
			w.WriteHeader(resp.StatusCode)
			return
		}
		defer resp.Body.Close()

		// Analyze and exec the project table
		tab := tables.ReadRights(resp.Body)
		projet.LectureProjet(tab)

		// Return none
		w.WriteHeader(http.StatusOK)

	case http.MethodPost:
		var tab []obj.Right
		if r.Header.Get("Content-Type") == "application/json" {
			// Already json table from body
			tools.Decode(r, &tab)
		} else {
			// Analyze the project csv table from body
			tab = tables.ReadRights(r.Body)
		}
		// Exec the project table
		projet.LectureProjet(tab)

	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
