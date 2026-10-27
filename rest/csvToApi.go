package rest

import (
	obj "main/objets"
	"main/tables"
	"main/tools"
	"net/http"
	"strings"
)

//  serve  GET  /agrem/csv/users ?search=xxx

// Http serve test api based on users.csv file
//http.Handle("/csv/", http.StripPrefix("/csv/", http.HandlerFunc(rest.ServeCsvApi)))

// REST call : search users, return json list
func ServeCsvApi(w http.ResponseWriter, r *http.Request) {

	// Check basic usage
	if r.Header.Get("PRIVATE_TOKEN") != "Good" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if r.Method != "GET" || !strings.HasPrefix(r.URL.Path, "users/") {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	ident := strings.Replace(r.URL.Path, "users/", "", 1)
	if !strings.ContainsAny(ident, "*._%") {

		// Read this user from file
		user := tables.ReadUser(ident, "./data/users.csv")
		if user == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		// Return as json array
		tools.Encode(w, []obj.User{*user})
		return
	}
	// Read all users from file
	users := tables.ReadUsers("./data/users.csv")
	if users == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// Filtre sur la regex ident
	result := []obj.User{}
	re := tools.InitMatch(ident)
	for _, u := range users {
		if tools.Match(re, u.Identifiant) {
			result = append(result, u)
		}
	}
	// Return as json array
	tools.Encode(w, result)
}

// Filtre un tableau
func Filter[T any](items []T, predicate func(T) bool) []T {
	var result []T
	for _, item := range items {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}
