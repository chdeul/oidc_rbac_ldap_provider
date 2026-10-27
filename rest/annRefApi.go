package rest

import (
	"main/projet/referent"
	"main/tools"
	"net/http"
)

//  serve  GET  /agrem/annref  : browse l'annuaire référentiel

var ref referent.Referentiel

func initAnnRef() {
	if ref == nil {
		ref = referent.RefFactory()
	}
}

// REST call : search users ou groups, return json list
func apiAnnRefRead(w http.ResponseWriter, r *http.Request) {

	initAnnRef()

	// Read only annuary
	if r.Method != "GET" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// Search users
	id := r.URL.Query().Get("user")
	if id != "" {
		users := ref.SearchUser(id)
		tools.Encode(w, users)
		return
	}
	// Else search groups ( users with this role )
	group := r.URL.Query().Get("group")
	if group != "" {
		users := ref.SearchGroup(group)
		tools.Encode(w, users)
		return
	}
	// Else Return none
	w.WriteHeader(http.StatusNoContent)
}
