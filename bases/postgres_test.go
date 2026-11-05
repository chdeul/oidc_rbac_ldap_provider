package bases

import (
	"database/sql"
	"fmt"
	obj "main/objets"
	"main/tools"
	"testing"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "mantis"
	password = "mantis"
	dbname   = "mantis"
)

// Prerequis  une base mantis postgres qui tourne

func TestPostgres(t *testing.T) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	tools.TEcheck(t, err, "Open")
	defer db.Close()

	db.SetMaxOpenConns(10)

	rows, err := db.Query("SELECT id, username, password, realname, email FROM mantis_user_table WHERE username = $1", "administrator")
	tools.TEcheck(t, err, "Query")
	defer rows.Close()
	tools.Tcheck(t, rows, "Query result")
	if rows != nil {
		for rows.Next() {
			var user obj.User
			rows.Scan(&user.Uid, &user.Identifiant, &user.MotDePasse, &user.NomComplet, &user.AdresseMail)
			tools.Tcheck(t, user.AdresseMail, "root@localhost")
		}
	}
}

func TestTrigger(t *testing.T) {

	t.Skip("TODO TestTrigger")

	//
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	tools.TEcheck(t, err, "Open")
	defer db.Close()

	addTrigger(db, psqlInfo, "mantis_user_table")

	waitTrigger(psqlInfo, "mantis_user_table")

	row := db.QueryRow("insert into mantis_user_table(username, password, realname, email) values('xxx', 'ppp','X Xxx', 'xxx@xxx.xx') returning id")
	var userid int
	row.Scan(&userid)

}

// TODO Other CRUD tests
