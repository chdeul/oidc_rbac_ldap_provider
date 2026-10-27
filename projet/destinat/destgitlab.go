package destinat

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	obj "main/objets"
	"main/tools"
	"os"
	"strconv"
	"strings"

	"resty.dev/v3"
)

// Authorization: Bearer <token_oauth_admin>
// https://gitlab.example.com/api/v4/projects

type GitlabConf struct {
	BaseUrl  string
	HeadAuth string
	CertFile string
	Insecure bool
}

type DestGitlab struct {
	cfg  *GitlabConf
	glab *resty.Client
}

func (de *DestGitlab) Name() string {
	return "gitlab"
}

func NewDestGitlab() (*DestGitlab, error) {

	// Charge la configuration
	var cfg GitlabConf
	err := tools.LoadAnyConfig("data/gitlab.yml", &cfg)
	if err != nil {
		panic(err)
	}
	// cree le client d'API Rest gitlab
	client := resty.New().SetDebug(false)
	client.SetAuthToken(cfg.HeadAuth)

	// Trust gitlab cert
	if cfg.CertFile != "" {
		certPool := x509.NewCertPool()
		pem, _ := os.ReadFile(cfg.CertFile)
		certPool.AppendCertsFromPEM(pem)
		client.SetTLSClientConfig(&tls.Config{
			RootCAs:            certPool,
			InsecureSkipVerify: cfg.Insecure,
		})
	}
	// Memorize configuration
	dest := DestGitlab{cfg: &cfg, glab: client}
	return &dest, err
}

func (de *DestGitlab) Close() {
	if de.cfg != nil {
		de.glab.Close()
		de.cfg = nil
	}
}

// Roles / Levels gitlab-ce

func Role(level int) string {
	switch level {
	case 0:
		return "No access"
	case 5:
		return "Minimal access"
	case 10:
		return "Guest"
	case 15:
		return "Planner"
	case 20:
		return "Reporter"
	case 25:
		return "Security Manager"
	case 30:
		return "Developer"
	case 40:
		return "Maintainer"
	case 50:
		return "Owner"
	default:
		return strconv.Itoa(level)
	}
}

func Level(role string) int {
	switch role {
	case "No access":
		return 0
	case "Minimal access":
		return 5
	case "Guest":
		return 10
	case "Planner":
		return 15
	case "Reporter":
		return 20
	case "Security Manager":
		return 25
	case "Developer":
		return 30
	case "Maintainer":
		return 40
	case "Owner":
		return 50
	default:
		return 30
	}
}

func (de *DestGitlab) SearchProjects(projet string) []*obj.Project {

	// Chercher le projet ou group  à partir de son nom flou

	// Table de la réponse gitlab
	var result []struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		Path        string `json:"path_with_namespace"`
		GPath       string `json:"full_path"`
		Description string `json:"description"`
	}

	// Chercher les groups d'abord
	resp, err := de.glab.R().
		SetQueryParams(map[string]string{
			"search": projet,
		}).
		SetResult(&result).
		Get(de.cfg.BaseUrl + "/groups")

	if err != nil {
		tools.Error("Recherche groups gitlab %v", err)
		return nil
	}
	// Conversion en Table rendue
	var tab []*obj.Project
	if resp.StatusCode() == 200 { // Search ok
		for _, p := range result {
			tab = append(tab, obj.NewProject(strconv.Itoa(p.Id), p.GPath, p.Name, p.Description))
		}
	} else {

		tools.Error("Recherche groups gitlab %s", resp.String())
		return nil
	}

	// Chercher les projets ensuite
	resp, err = de.glab.R().
		SetQueryParams(map[string]string{
			"search": projet,
		}).
		SetResult(&result).
		Get(de.cfg.BaseUrl + "/projects")

	if err != nil {
		tools.Error("Recherche projets gitlab %v", err)
		return nil
	}
	if resp.StatusCode() == 200 { // Search ok

		// Conversion en Table rendue  ( .git for projects )
		for _, p := range result {
			tab = append(tab, obj.NewProject(strconv.Itoa(p.Id), p.Path+".git", p.Name, p.Description))
		}

	} else {

		tools.Error("Recherche projets gitlab %s", resp.String())
		return nil
	}
	return tab
}

// Methodes d'interface

func (de *DestGitlab) SearchUsersRole(prj *obj.Project) []*obj.User {
	return nil
}

func (de *DestGitlab) GetRole(prj *obj.Project, user *obj.User) string {

	// Table de la réponse gitlab
	var result struct {
		Id    int `json:"id"`
		Level int `json:"access_level"`
	}
	categ := "groups"
	if strings.HasSuffix(prj.Identifiant, ".git") {
		categ = "projects"
	}
	// Chercher les roles  GET /members
	resp, err := de.glab.R().
		SetPathParams(map[string]string{
			"pid": prj.Pid,
			"uid": user.Uid,
		}).
		SetResult(&result).
		Get(de.cfg.BaseUrl + "/" + categ + "/{pid}/members/all/{uid}")

	if err != nil {
		tools.Error("Recherche roles gitlab %v", err)
		return ""
	}
	if resp.StatusCode() == 200 { // Search ok

		return Role(result.Level)
	}
	tools.Error("Recherche roles gitlab %s", resp.String())
	return ""
}

func (de *DestGitlab) CreateRole(prj *obj.Project, user *obj.User, role string) error {

	categ := "groups"
	if strings.HasSuffix(prj.Identifiant, ".git") {
		categ = "projects"
	}

	// Creer un utilisateur  POST /users
	resp, err := de.glab.R().
		SetPathParams(map[string]string{
			"pid": prj.Pid,
		}).
		SetBody(map[string]string{
			"user_id":      user.Uid,
			"access_level": strconv.Itoa(Level(role)),
		}).
		Post(de.cfg.BaseUrl + "/" + categ + "/{pid}/members")

	if err != nil {
		tools.Error("Creer role gitlab %v", err)
		return err
	}
	if resp.StatusCode() == 201 { // Created)
		return nil
	}
	tools.Error("Creer role gitlab %s", resp.String())
	return errors.New(resp.String())
}

func (de *DestGitlab) UpdateRole(prj *obj.Project, user *obj.User, role string) error {

	if user.Uid == "1" {
		tools.Error("Mise à jour membre gitlab { Cannot update root member }")
		return errors.New("Cannot update root member")
	}
	categ := "groups"
	if strings.HasSuffix(prj.Identifiant, ".git") {
		categ = "projects"
	}

	// M%ise à jour un membre
	resp, err := de.glab.R().
		SetPathParams(map[string]string{
			"pid": prj.Pid,
			"uid": user.Uid,
		}).
		SetBody(map[string]string{
			"access_level": strconv.Itoa(Level(role)),
		}).
		Put(de.cfg.BaseUrl + "/" + categ + "/{pid}/members/{uid}")

	if err != nil {
		tools.Error("Mise à jour membre gitlab %v", err)
		return err
	}
	if resp.StatusCode() == 200 { // Ok
		return nil
	}
	tools.Error("Mise à jour membre gitlab %s", resp.String())
	return errors.New(resp.String())

}

func (de *DestGitlab) DeleteRole(prj *obj.Project, user *obj.User) error {

	if user.Uid == "1" {
		tools.Error("Supprimer membre gitlab { Cannot remove root member }")
		return errors.New("Cannot remove root member")
	}
	categ := "groups"
	if strings.HasSuffix(prj.Identifiant, ".git") {
		categ = "projects"
	}

	// Supprimer un membre
	resp, err := de.glab.R().
		SetPathParams(map[string]string{
			"pid": prj.Pid,
			"uid": user.Uid,
		}).
		Delete(de.cfg.BaseUrl + "/" + categ + "/{pid}/members/{uid}")

	if err != nil {
		tools.Error("Supprimer utilisateur gitlab %v", err)
		return err
	}
	if resp.StatusCode() == 204 { // No content , deleted
		return nil
	}
	tools.Error("Supprimer utilisateur gitlab %s", resp.String())
	return errors.New(resp.String())
}

// gitlab users methodes

func (de *DestGitlab) SearchUsers(ident string) []*obj.User {

	// Table de la réponse gitlab
	var result []struct {
		Id       int    `json:"id"`
		Username string `json:"username"`
		Name     string `json:"name"`
		Email    string `json:"email"`
	}
	// Chercher les utilisateurs  GET /users
	resp, err := de.glab.R().
		SetPathParams(map[string]string{
			"username": ident, // ou morceau
		}).
		SetResult(&result).
		Get(de.cfg.BaseUrl + "/users?search={username}")

	if err != nil {
		tools.Error("Recherche utilisateur gitlab %v", err)
		return nil
	}
	if resp.StatusCode() == 200 { // Search ok
		// Conversion en Table rendue
		tab := make([]*obj.User, len(result))
		for i, p := range result {
			tab[i] = obj.NewUser(strconv.Itoa(p.Id), p.Username, p.Name, p.Email, "")
		}
		return tab
	}
	tools.Error("Recherche utilisateur gitlab %s", resp.String())
	return nil
}

func (de *DestGitlab) CreateUser(user *obj.User) error {

	// Réponse gitlab
	var result struct {
		Id       int    `json:"id"`
		Username string `json:"username"`
	}

	// User has property external
	hasExternal := strings.Contains(user.Roles, "external")

	// Creer un utilisateur  POST /users
	resp, err := de.glab.R().
		SetBody(map[string]any{
			"username":              user.Identifiant,
			"name":                  user.NomComplet,
			"email":                 user.AdresseMail,
			"force_random_password": true,
			"external":              hasExternal,
		}).
		SetResult(&result).
		Post(de.cfg.BaseUrl + "/users")

	if err != nil {
		tools.Error("Creer utilisateur gitlab %v", err)
		return err
	}
	if resp.StatusCode() == 201 { // Created
		if user.Identifiant == result.Username {
			user.Uid = strconv.Itoa(result.Id)
			return nil
		}
	}
	tools.Error("Creer utilisateur gitlab %s", resp.String())
	return errors.New(resp.String())
}

func (de *DestGitlab) UpdateUser(user *obj.User) error {

	// Mettre à jour un utilisateur  PUT /users
	resp, err := de.glab.R().
		SetPathParams(map[string]string{
			"id": user.Uid,
		}).
		SetBody(map[string]string{
			"username": user.Identifiant,
			"name":     user.NomComplet,
			"email":    user.AdresseMail,
		}).
		Put(de.cfg.BaseUrl + "/users/{id}")

	if err != nil {
		tools.Error("Mise à jour utilisateur gitlab %v", err)
		return err
	}
	if resp.StatusCode() == 200 { // ok
		return nil
	}
	tools.Error("Mise à jour utilisateur gitlab %s", resp.String())
	return errors.New(resp.String())
}

func (de *DestGitlab) DeleteUser(user *obj.User) error {

	if user.Uid == "1" || user.Identifiant == "root" {
		tools.Error("Supprimer utilisateur gitlab { Cannot remove root user }")
		return errors.New("Cannot remove root user")
	}
	// Supprimer un utilisateur  DELETE /users  (but not his contributions !)
	resp, err := de.glab.R().
		SetPathParams(map[string]string{
			"id": user.Uid,
		}).
		Delete(de.cfg.BaseUrl + "/users/{id}")

	if err != nil {
		tools.Error("Supprimer utilisateur gitlab %v", err)
		return err
	}
	if resp.StatusCode() == 204 { // No content , deleted
		return nil
	}
	tools.Error("Supprimer utilisateur gitlab %s", resp.String())
	return errors.New(resp.String())
}
