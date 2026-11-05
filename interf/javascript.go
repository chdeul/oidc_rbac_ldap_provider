package interf

import (
	"main/tools"
	"os"

	"github.com/dop251/goja"
)

func executeJsScript(filename string, input []byte) []byte {

	// Load js code ( relative to Workspace Dir )
	javascript, err := os.ReadFile(tools.CfgFilename(filename))
	if err != nil {
		return nil
	}
	// js runtime embed in golang
	vm := goja.New()

	// Set a Go value into JS
	vm.Set("input", string(input))

	// Execute JS using that value
	result, err := vm.RunString(string(javascript))
	if err != nil {
		return nil
	}
	return []byte(result.String())
}
