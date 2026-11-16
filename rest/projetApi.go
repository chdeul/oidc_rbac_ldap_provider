package rest

import (
	obj "main/objets"
	"main/projet"
	"main/tables"
	"main/tools"
	"net/http"
)

//  serve  GET /agrem/projet  :  retourne
//  serve  PUT POST  /agrem/projet  :  prise en compte

// REST call :
func apiProjet(w http.ResponseWriter, r *http.Request) {

	// Read users directory
	initAnnRef()

	switch r.Method {

	case http.MethodGet:

		// Recupere le nom du projet
		name := r.URL.Query().Get("name")
		if name != "" {
			// Return json table of project
			tab := projet.LoadCsvProject(name)
			// out csv
			if r.Header.Get("Accept") == "text/csv" {
				tables.WriteRights(tab, w)
				return
			}
			// default: out json
			tools.Encode(w, tab)
			return
		}
		// liste des noms possibles
		list := projet.ListOfCsvProject()
		tools.Encode(w, list)

	case http.MethodPut:

		// Recupere l'url du lien
		var url string
		tools.Decode(r, &url)

		// Recupere le csv par get sur cette url
		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		tools.CheckError(err, "http.NewRequest failed: %v", err)
		req.Header.Add("Authorization", "agpat-1234567890AbcDefGhij") // TODO ( config auth ? )
		req.Header.Add("Accept", "text/csv")
		resp, err := client.Do(req)
		if !tools.CheckError(err, "GET project table failed: %v", err) {
			w.WriteHeader(resp.StatusCode)
			return
		}
		defer resp.Body.Close()
		if !tools.Check(resp.StatusCode == 200, "GET project table failed: http status %d", resp.StatusCode) {
			w.WriteHeader(resp.StatusCode)
			return
		}
		// Analyze and exec the project table
		tab := tables.ReadRights(resp.Body)
		if tools.Check(tab, "GET project table failed: wrong data") {
			projet.LectureProjet(tab)
			w.WriteHeader(http.StatusOK)
		}
		// Return error
		w.WriteHeader(http.StatusNotAcceptable)

	case http.MethodPost:

		// Recupere les données directement json ou csv
		var tab []obj.Right
		if r.Header.Get("Content-Type") == "application/json" {
			// Already json table from body
			tools.Decode(r, &tab)
		} else {
			// Analyze the project csv table from body
			tab = tables.ReadRights(r.Body)
		}
		// Exec the project table
		if tools.Check(tab, "GET project table failed: wrong data") {
			projet.LectureProjet(tab)
			w.WriteHeader(http.StatusOK)
			return
		}
		// Return error
		w.WriteHeader(http.StatusNotAcceptable)

	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
