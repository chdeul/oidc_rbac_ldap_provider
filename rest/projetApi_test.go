package rest

import (
	"bytes"
	"encoding/json"
	"io"
	obj "main/objets"
	"main/tools"
	"net/http"
	"testing"
)

func httpTest(meth, url string, data any) (resp *http.Response, err error) {
	client := &http.Client{}
	var datab io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		tools.AssertError(err, "json.Marshal failed: %v", err)
		datab = bytes.NewBuffer(jsonData)
	}
	req, err := http.NewRequest(meth, url, datab)
	tools.AssertError(err, "http.NewRequest failed: %v", err)
	req.Header.Add("Authorization", "agpat-1234567890AbcDefGhij")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	return client.Do(req)
}

func TestProjetApi(t *testing.T) {

	// Get without auth
	resp, err := http.Get("http://localhost:8088/agrem/api/projet")
	tools.TEcheck(t, err, "GET failed: %v", err)
	defer resp.Body.Close()
	tools.Tcheck(t, resp.StatusCode == 401, "List projects")

	// Get projects list
	var pl []string
	resp, err = httpTest("GET", "http://localhost:8088/agrem/api/projet", nil)
	tools.TEcheck(t, err, "GET failed: %v", err)
	defer resp.Body.Close()
	tools.Tcheck(t, resp.StatusCode == 200, "List projects")
	json.NewDecoder(resp.Body).Decode(&pl)
	tools.Tcheck(t, len(pl) > 0, "Empty List of projects")

	// Get a project in json
	var prj []obj.Right
	resp, err = httpTest("GET", "http://localhost:8088/agrem/api/projet?name="+pl[0], nil)
	tools.TEcheck(t, err, "GET failed: %v", err)
	defer resp.Body.Close()
	tools.Tcheck(t, resp.StatusCode == 200, "Read a project users/rights")
	err = json.NewDecoder(resp.Body).Decode(&prj)
	tools.TEcheck(t, err, "Decode project")
	tools.Tcheck(t, len(prj) > 0, "Empty project")

	// Post a project ( table of user/rights )
	resp, err = httpTest("POST", "http://localhost:8088/agrem/api/projet", prj)
	tools.TEcheck(t, err, "POST failed: %v", err)
	tools.Tcheck(t, resp.StatusCode == 200, "POST project data")

	// Put a link to project ( an url where read table )
	data := "http://localhost:8088/agrem/api/projet?name=" + pl[0]
	resp, err = httpTest("PUT", "http://localhost:8088/agrem/api/projet", data)
	tools.TEcheck(t, err, "PUT failed: %v", err)
	tools.Tcheck(t, resp.StatusCode == 200, "PUT link project")
}
