package rest

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	// tools.LoadConfig("data/test/config.yml")

	os.Exit(m.Run())
}

// Nothing testable here for REST Router
// See each route test
