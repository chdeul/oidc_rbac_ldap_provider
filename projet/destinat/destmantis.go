package destinat

import (
	"crypto/tls"
	"crypto/x509"
	obj "main/objets"
	"main/tools"
	"os"
	"strconv"
	"strings"

	"resty.dev/v3"
)

/*
 * Mantis bugs tracker : projects, users, rbac
 */

type MantisConf struct {
	BaseUrl  string
	HeadAuth string
	CertFile string
	Insecure bool
}

type DestMantis struct {
	cfg *MantisConf
	mbt *resty.Client
}

func (de *DestMantis) Name() string {
	return "mantis"
}

func NewDestMantis() (*DestMantis, error) {

	// Charge la configuration
	var cfg MantisConf
	err := tools.LoadAnyConfig("data/mantis.yml", &cfg)
	if err != nil {
		panic(err)
	}
	// cree le client d'API Rest Mantisbt
	client := resty.New().SetDebug(false)
	// client.SetAuthToken(cfg.HeadAuth)
	client.SetHeader("Authorization", cfg.HeadAuth)

	// Trust mantis cert
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
	dest := DestMantis{cfg: &cfg, mbt: client}
	return &dest, err
}

func (de *DestMantis) Close() {
	if de.cfg != nil {
		de.mbt.Close()
		de.cfg = nil
	}
}

// mantis projects methodes

func (de *DestMantis) SearchUsersRole(prj *obj.Project) []*obj.User {
	return nil
}

func (de *DestMantis) SearchProjects(projet string) []*obj.Project {

	// Chercher le projet à partir de son nom flou

	// Table de la réponse gitlab
	var result struct {
		Projects []struct {
			Id          int    `json:"id"`
			Name        string `json:"name"`
			Enabled     bool   `json:"enabled"`
			Description string `json:"description"`
		} `json:"projects"`
	}

	// Chercher les groups d'abord
	resp, err := de.mbt.R().
		SetResult(&result).
		Get(de.cfg.BaseUrl + "/projects")

	if err != nil {
		tools.Error("Recherche projet mantis %v", err)
		return nil
	}
	// Match en Table rendue
	var tab []*obj.Project
	if resp.StatusCode() == 200 { // Search ok
		for _, p := range result.Projects {
			if p.Enabled && strings.Contains(p.Name, projet) {
				tab = append(tab, obj.NewProject(strconv.Itoa(p.Id), p.Name, p.Name, p.Description))
			}
		}
	} else {
		tools.Error("Recherche projet mantis %s", resp.String())
		return nil
	}
	return tab
}

// mantis users methodes

func (de *DestMantis) SearchUsers(ident string) []*obj.User {

	// Table de la réponse
	var result struct {
		Users []struct {
			Id       int    `json:"id"`
			Name     string `json:"name"`
			RealName string `json:"real_name"`
			Email    string `json:"email"`
			Level    struct {
				Name string `json:"name"`
			} `json:"access_level"`
		} `json:"users"`
	}

	// Chercher cet utilisateur  ( exact )
	resp, err := de.mbt.R().
		SetPathParam("name", ident).
		SetResult(&result).
		Get(de.cfg.BaseUrl + "/users/username/{name}")

	if err != nil {
		tools.Error("Recherche utilisateur mantis %v", err)
		return nil
	}
	// Conversion en Table de 1 rendue
	var tab []*obj.User
	if resp.StatusCode() == 200 { // Search ok
		for _, p := range result.Users {
			tab = append(tab, obj.NewUser(strconv.Itoa(p.Id), p.Name, p.RealName, p.Email, p.Level.Name))
		}
	} else {
		tools.Error("Recherche utilisateur mantis %s", resp.String())
		return nil
	}
	return tab
}

func (de *DestMantis) CreateUser(user *obj.User) error {

	var result struct {
		User struct {
			Id int `json:"id"`
		} `json:"user"`
	}
	resp, err := de.mbt.R().
		SetBody(map[string]any{
			"name":         user.Identifiant,
			"real_name":    user.NomComplet,
			"email":        user.AdresseMail,
			"enabled":      true,
			"protected":    false,
			"access_level": map[string]string{"name": "viewer"},
			"password":     "",
		}).
		SetResult(&result).
		Post(de.cfg.BaseUrl + "/users")

	if err != nil {
		tools.Error("Création utilisateur mantis %v", err)
		return err
	}
	if resp.StatusCode() != 201 {
		tools.Error("Création utilisateur mantis %s", resp.String())
	}
	user.Uid = strconv.Itoa(result.User.Id)
	return nil
}

func (de *DestMantis) UpdateUser(user *obj.User) error {

	resp, err := de.mbt.R().
		SetPathParam("uid", user.Uid).
		SetBody(map[string]any{
			"user": map[string]any{
				"name":         user.Identifiant,
				"real_name":    user.NomComplet,
				"email":        user.AdresseMail,
				"enabled":      true,
				"protected":    false,
				"access_level": map[string]string{"name": "viewer"},
			},
		}).
		Patch(de.cfg.BaseUrl + "/users/{uid}")

	if err != nil {
		tools.Error("Mise à jour utilisateur mantis %v", err)
		return err
	}
	if resp.StatusCode() != 200 {
		tools.Error("Mise à jour utilisateur mantis %s", resp.String())
	}
	return nil

}

func (de *DestMantis) DeleteUser(user *obj.User) error {

	resp, err := de.mbt.R().
		SetPathParam("uid", user.Uid).
		Delete(de.cfg.BaseUrl + "/users/{uid}")

	if err != nil {
		tools.Error("Suppression utilisateur mantis %v", err)
		return err
	}
	if resp.StatusCode() != 204 {
		tools.Error("Suppression utilisateur mantis %s", resp.String())
	}
	return nil
}

// mantis role methodes

func (de *DestMantis) GetRole(prj *obj.Project, user *obj.User) string {

	var result struct {
		Users []struct {
			Name  string `json:"name"`
			Level struct {
				Name string `json:"name"`
			} `json:"access_level"`
		} `json:"users"`
	}
	resp, err := de.mbt.R().
		SetPathParam("pid", prj.Pid).
		SetResult(&result).
		Get(de.cfg.BaseUrl + "/projects/{pid}/users")

	if err != nil {
		tools.Error("Role utilisateur mantis %v", err)
		return ""
	}
	if resp.StatusCode() != 200 {
		tools.Error("Role utilisateur mantis %s", resp.String())
		return ""
	}
	// Recherche cet utilisateur
	for _, r := range result.Users {
		if r.Name == user.Identifiant {
			return r.Level.Name
		}
	}
	return ""
}

func (de *DestMantis) CreateRole(prj *obj.Project, user *obj.User, role string) error {

	//id, _ := strconv.Atoi(user.Uid)
	resp, err := de.mbt.R().
		SetPathParam("pid", prj.Pid).
		SetBody(map[string]any{
			"user":         map[string]string{"name": user.Identifiant},
			"access_level": map[string]string{"name": role},
		}).
		Put(de.cfg.BaseUrl + "/projects/{pid}/users")

	if err != nil {
		tools.Error("Création/Mise à jour Role utilisateur mantis %v", err)
		return err
	}
	if resp.StatusCode() != 204 {
		tools.Error("Création/Mise à jour Role utilisateur mantis %s", resp.String())
	}
	return nil
}

func (de *DestMantis) UpdateRole(prj *obj.Project, user *obj.User, role string) error {
	return de.CreateRole(prj, user, role)
}

func (de *DestMantis) DeleteRole(prj *obj.Project, user *obj.User) error {
	resp, err := de.mbt.R().
		SetPathParams(map[string]string{
			"pid": prj.Pid,
			"uid": user.Uid,
		}).
		Delete(de.cfg.BaseUrl + "/projects/{pid}/users/{uid}")

	if err != nil {
		tools.Error("Suppression Role utilisateur mantis %v", err)
		return err
	}
	if resp.StatusCode() != 204 {
		tools.Error("Suppression Role utilisateur mantis %s", resp.String())
	}
	return nil
}
