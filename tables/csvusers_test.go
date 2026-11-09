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
	var fn = "test.csv"
	var tab []obj.User = make([]obj.User, 2)
	tab[0] = *NewUser("1", "aa", "%|\"'", "bb_b:", "ccc xxx", "")
	tab[1] = *NewUser("2", "aa%|\"'", "", "bb,b", "", "chef")
	err := WriteUsersFile(tab, fn)
	tools.TEcheck(t, err, "WriteUsers")

	res := ReadUsersFile(fn)
	tools.Tcheck(t, len(res) == 2, "ReadUsers")
	tools.Tcheck(t, res[0].AdresseMail == "ccc xxx" && res[1].AdresseMail == "", "ReadUsers")

	err = WriteUser(NewUser("3", "zz", "_", "xxxxxxx", "123.56", "chef"), fn)
	tools.TEcheck(t, err, "WriteUser")

	ru := ReadUser("zz", fn)
	tools.Tcheck(t, ru.Identifiant == "zz" && ru.MotDePasse == "_" && ru.AdresseMail == "123.56", "ReadUser")
}

func cleanCsvUsers(t *testing.T) {
	os.Remove("test.csv")
}
