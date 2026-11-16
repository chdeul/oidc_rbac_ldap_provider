package interf

import (
	"encoding/json"
	obj "main/objets"
	"main/tools"
	"testing"
)

// Prerequis lancer une API REST ( ex: Mantis )

func MapToStruct(m any) (obj.User, error) {
	var v obj.User
	b, err := json.Marshal(m)
	if tools.CheckError(err) {
		err = json.Unmarshal(b, &v)
	}
	return v, err
}

func TestRestInput(t *testing.T) {

	// Config appel api REST Mantis
	var cfg RestInputConfig
	err := tools.LoadAnyConfig("data/test/restInput.yml", &cfg)
	tools.TEcheck(t, err, "Load Rest Input Config")

	// Test 1 : recup Objet
	params := map[string]any{}
	params["token"] = "2YlTDV80nt8ZQhleYfi_CLfoyLZ1gR1T"
	ret := RestInputFct(&cfg, "GET", "test/1", params)
	tools.Tcheck(t, ret, "RestInputFct failed %v", ret)
	obj, err := MapToStruct(ret)
	tools.TEcheck(t, err, "RestInputFct type failed")
	tools.Tcheck(t, obj.Identifiant == "administrator")
	tools.Tcheck(t, obj.Uid == "1")
	tools.Tcheck(t, obj.Roles == "administrator")

	// Test 2 : recup tableau d'objets
	params = map[string]any{}
	params["prj"] = "2"
	ret = RestInputFct(&cfg, "GET", "test2", params)
	tools.Tcheck(t, ret, "RestInputFct failed %v", ret)
	obj2, ok2 := ret.([]map[string]any)
	tools.Tcheck(t, ok2, "RestInputFct type failed")
	tools.Tcheck(t, len(obj2) > 0, "RestInputFct type failed")
	obj, err = MapToStruct(obj2[0])
	tools.TEcheck(t, err, "RestInputFct type failed")
	tools.Tcheck(t, obj.Identifiant == "administrator")
	tools.Tcheck(t, obj.Uid == "1")
	tools.Tcheck(t, obj.Roles == "administrateur")

	// Test 3 : recup objet extrait en js
	params = map[string]any{}
	ret = RestInputFct(&cfg, "GET", "test3/test2?name=2", params)
	tools.Tcheck(t, ret, "RestInputFct failed %v", ret)
	obj, err = MapToStruct(ret)
	tools.TEcheck(t, err, "RestInputFct type failed")
	tools.Tcheck(t, obj.Identifiant == "administrator")
	tools.Tcheck(t, obj.Uid == "1")
	tools.Tcheck(t, obj.Roles == "90:administrateur")
}
