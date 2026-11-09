package rest

import (
	"main/projet/referent"
	"main/tables"
	"main/tools"
	"net/http"
	"strings"
)

//  serve  GET  /agrem/annref  : browse l'annuaire référentiel

var ref referent.Referentiel

func initAnnRef() {
	if ref == nil {
		ref = referent.RefFactory()
	}
}

// REST call : search users ou groups, return csv list (sep;)
func apiAnnRefRead(w http.ResponseWriter, r *http.Request) {

	initAnnRef()

	csv := false
	accept := r.Header.Get("Accept")
	if (strings.Contains(accept, "text/csv") || strings.Contains(accept, "*/*")) && !strings.Contains(accept, "application/json") {
		csv = true
	}

	// Read only directory
	if r.Method != "GET" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// Search groups ( users with this role )
	group := r.URL.Query().Get("group")
	if group != "" {
		users := ref.SearchGroup(group)
		if csv {
			tables.EncodeCsv(w, users)
		} else {
			tools.Encode(w, users)
		}
		return
	}
	// Else Search users (default)
	id := r.URL.Query().Get("user")
	if id == "" {
		id = "*"
	}
	users := ref.SearchUser(id)
	if csv {
		tables.EncodeCsv(w, users)
	} else {
		tools.Encode(w, users)
	}
}
