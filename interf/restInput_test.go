package interf

import (
	"main/tools"
	"testing"
)

// Prerequis lancer une API REST

func TestRestInput(t *testing.T) {

	cfg, err := LoadInputConfig("./restInput.yml")
	tools.TEcheck(t, err, "LoadInputConfig")

	ret := RestInputFct(cfg, "GET", "users/{id}", nil)
	tools.Tcheck(t, ret, "RestInputFct failed")
}
