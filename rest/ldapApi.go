package rest

import (
	"main/bases"
	obj "main/objets"
	"main/tools"
	"net/http"
)

var conn, wconn any

// Group with members
type groupData struct {
	Group   string
	Members []string
}

func open(write bool) {
	if !write && conn == nil {
		url := tools.Cfg.Dest.Ldap.Url
		bind := tools.Cfg.Dest.Ldap.Root
		conn, _ = bases.LdapOpen(url, bind, "")
	}
	if write && wconn == nil {
		url := tools.Cfg.Dest.Ldap.Url
		bind := tools.Cfg.Dest.Ldap.Root
		pass := tools.Cfg.Dest.Ldap.Passwd
		wconn, _ = bases.LdapOpen(url, bind, pass)
	}
}

// Query groups=  users=   user=  group=

// REST call : search users or group, return json list
func apiLdapRead(w http.ResponseWriter, r *http.Request) {
	open(false)
	group := r.URL.Query().Get("group")
	user := r.URL.Query().Get("user")
	membersOf := r.URL.Query().Get("membersOf")
	if membersOf != "" {
		users := bases.LdapGetGroup(conn, membersOf)
		tools.Encode(w, users)
	} else if user != "" {
		users := bases.LdapSearch(conn, user)
		tools.Encode(w, users)
	} else if group != "" {
		groups := bases.LdapSearchGroups(conn, group)
		tools.Encode(w, groups)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

// REST call : add user or group
func apiLdapCreate(w http.ResponseWriter, r *http.Request) {
	open(true)

	isGroup := r.URL.Query().Get("group") == "true"
	if isGroup {
		//  Group to create with members
		var data groupData
		tools.Decode(r, &data)
		err := bases.LdapAddGroup(wconn, data.Group, "...") // TODO Description à éditer et afficher
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
		// update members
		err = updateMembers(&data)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}

	} else {
		//  User to create
		var user obj.User
		tools.Decode(r, &user)
		err := bases.LdapAddUser(wconn, &user)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

// REST call : update user or group
func apiLdapUpdate(w http.ResponseWriter, r *http.Request) {
	open(false)
	open(true)
	isGroup := r.URL.Query().Get("group") == "true"
	isCheck := r.URL.Query().Get("check") == "true"
	isChange := r.URL.Query().Get("change") == "true"
	if isGroup {
		//  Group to update members
		var data groupData
		tools.Decode(r, &data)
		// update members
		err := updateMembers(&data)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	//  User modified
	var user obj.User
	tools.Decode(r, &user)

	if isCheck {
		// First check password
		url := tools.Cfg.Dest.Ldap.Url
		bind := tools.Cfg.Dest.Ldap.Root
		cheched := bases.LdapCheckPassword(url, bind, user.Identifiant, user.MotDePasse)
		if !cheched {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
	} else if isChange {
		// First change password
		err := bases.LdapSetPassword(wconn, user.Identifiant, user.MotDePasse)
		if err != nil {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}
	// Then update fields (except username and password)
	err := bases.LdapUpdateUser(wconn, &user)
	if err != nil {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// REST call : delete user or group
func apiLdapDelete(w http.ResponseWriter, r *http.Request) {
	open(true)
	var err error
	group := r.URL.Query().Get("group")
	user := r.URL.Query().Get("user")
	if group != "" {
		// Delete group
		err = bases.LdapDeleteGroup(wconn, group)
	} else if user != "" {
		// Delete user  and update group members
		err = bases.LdapDeleteUser(wconn, user)
	} else {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Update members of an existing group
func updateMembers(data *groupData) error {

	// Membres existants
	users := bases.LdapGetGroup(conn, data.Group)
	for _, u := range users {

		// Si Dans nouveaux membres
		var trouve bool = false
		for _, m := range data.Members {
			if u == m {
				//  Oui : déjà
				trouve = true
				break
			}
		}
		// Non : le retire  ( garde toujours admin )
		if !trouve && u != "admin" {
			err := bases.LdapDeleteMember(wconn, data.Group, u)
			if err != nil {
				return err
			}
		}
	}
	// Membres nouveaux
	for _, m := range data.Members {

		// Cherche dans memebres existants
		var trouve bool = false
		for _, u := range users {
			if u == m {
				trouve = true
				break
			}
		}
		// Non : l'ajoute
		if !trouve {
			err := bases.LdapAddMember(wconn, data.Group, m)
			if err != nil {
				return err
			}
		}
	}
	// Ok, tout s'est bien passé
	return nil
}
