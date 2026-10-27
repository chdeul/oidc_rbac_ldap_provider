package tables

import (
	obj "main/objets"
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
	err := WriteUsers(tab, fn)
	if err != nil {
		t.Error("WriteUsers failed")
	}
	res := ReadUsers(fn)
	if len(res) != 2 {
		t.Fatal("ReadUsers failed")
	}
	if res[0].AdresseMail != "ccc xxx" || res[1].AdresseMail != "" {
		t.Error("ReadUser failed")
	}
	err = WriteUser(NewUser("3", "zz", "_", "xxxxxxx", "123.56", "chef"), fn)
	if err != nil {
		t.Error("WriteUser failed")
	}
	ru := ReadUser("zz", fn)
	if ru.Identifiant != "zz" || ru.MotDePasse != "_" || ru.AdresseMail != "123.56" {
		t.Error("ReadUser failed")
	}
}

func cleanCsvUsers(t *testing.T) {
	os.Remove("test.csv")
}
