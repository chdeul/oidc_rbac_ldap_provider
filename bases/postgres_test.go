package bases

import (
	"database/sql"
	"fmt"
	obj "main/objets"
	"testing"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "antoine"
	password = "chdeul"
	dbname   = "test"
)

// Prerequis  une base postgres qui tourne
// avec une base "test", une table "users": uid, ident, password, name, email
// Une ligne d'ident=antoine

func TestPostgres(t *testing.T) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	echeck(t, err)
	defer db.Close()

	db.SetMaxOpenConns(10)

	rows, err := db.Query("SELECT uid, ident, password, name, email FROM users WHERE ident = $1", "antoine")
	echeck(t, err)
	defer rows.Close()
	check(t, rows)

	for rows.Next() {
		var user obj.User
		rows.Scan(&user.Uid, &user.Identifiant, &user.MotDePasse, &user.NomComplet, &user.AdresseMail)
		check(t, user.AdresseMail)
	}
}

func TestTrigger(t *testing.T) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	echeck(t, err)
	defer db.Close()

	go addTrigger(db, psqlInfo, "users")
	row := db.QueryRow(`insert into users(name, age) values('Scrooge McDuck', 93) returning id`)
	var userid int
	row.Scan(&userid)

}

// TODO Other CRUD tests
