package tables

import (
	obj "main/objets"
	"main/tools"
	"os"
	"testing"
)

func TestCsvUsers(t *testing.T) {

	t.Run("step1", testCsvUsers)

	t.Run("clean", cleanCsvUsers)
}

func testCsvUsers(t *testing.T) {

	// ; separator
	fn := "test.csv"
	var tab []obj.User = make([]obj.User, 2)
	tab[0] = *obj.NewUser("1", "aa", "%|\"'", "bb_b:", "ccc xxx", "")
	tab[1] = *obj.NewUser("2", "aa%|\"'", "", "bb,b", "", "chef")
	err := WriteUsers(tab, fn)
	tools.TEcheck(t, err, "WriteUsers")

	res := ReadUsers(fn)
	tools.Tcheck(t, len(res) == 2, "ReadUsers")
	tools.Tcheck(t, res[0].AdresseMail == "ccc xxx" && res[1].AdresseMail == "", "ReadUsers")
}

func cleanCsvUsers(t *testing.T) {
	os.Remove("test.csv")
}
