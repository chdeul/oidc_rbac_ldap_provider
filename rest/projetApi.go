package rest

import (
	"main/projet"
	"net/http"
)

//  serve  GET  /agrem/projet  :  prise en compte

// REST call : search users, return json list
func apiProjet(w http.ResponseWriter, r *http.Request) {

	initAnnRef()

	if r.Method != "GET" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	projet.LectureProjet("resources/projet.csv")

	// Return none
	w.WriteHeader(http.StatusOK)
}
