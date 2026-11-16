package interf

import (
	"main/tools"
	"os"

	"github.com/dop251/goja"
)

func executeJsScript(filename string, input []byte) []byte {

	// Load js code ( relative to Workspace Dir )
	javascript, err := os.ReadFile(tools.CfgFilename(filename))
	if !tools.CheckError(err, "Fichier javascript non trouvé %v", err) {
		return nil
	}
	// js runtime embed in golang
	vm := goja.New()

	// Set a Go json value into JS variable "input"
	vm.Set("input", string(input))

	// Execute JS custom code to parse input and produce computed result
	result, err := vm.RunString(string(javascript))
	if !tools.CheckError(err, "Javascript error %v", err) {
		return nil
	}
	// Get json object computed
	return []byte(result.String())
}
