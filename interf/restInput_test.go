package interf

import (
	"os"
	"testing"
)

func testRestInput(t *testing.T) {

	cfg, err := LoadInputConfig("./restInput.yml")
	if err != nil {
		os.Exit(-1)
	}
	ret := RestInputFct(cfg, "GET", "users/{id}", nil)
	if ret == nil {
		t.Error("RestInputFct failed")
	}
}
