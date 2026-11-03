package interf

import (
	"main/tools"
	"testing"
)

func TestJs(t *testing.T) {
	var input []byte = []byte("Oh")
	ret := executeJsScript("../data/test/testjs.js", input)
	tools.Tcheck(t, ret, "executeJsScript")
	val := string(ret)
	tools.Check(t, val == "Hello! Oh")
}
