package interf

import (
	"main/tools"
	"testing"
)

// Prerequis lancer une API REST

func TestRestInput(t *testing.T) {

	t.Skip("TODO TestRestInput")

	var cfg RestInputConfig
	err := tools.LoadAnyConfig("data/test/restInput.yml", &cfg)
	tools.TEcheck(t, err, "Load Rest Input Config")

	ret := RestInputFct(&cfg, "GET", "users/{id}", nil)
	tools.Tcheck(t, ret, "RestInputFct failed")
}
