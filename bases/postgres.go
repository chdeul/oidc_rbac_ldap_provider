package bases

import (
	"database/sql"
	"fmt"
	"log"
	obj "main/objets"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "antoine"
	password = "chdeul"
	dbname   = "test"
)

func LeTest() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT uid, ident, password, name, email FROM users WHERE ident = $1", "antoine")
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var user obj.User
		rows.Scan(&user.Uid, &user.Identifiant, &user.MotDePasse, &user.NomComplet, &user.AdresseMail)
		if user.AdresseMail == "" {
			return
		}
	}
	go addTrigger(db, psqlInfo, "users")
	/*
		    row := db.QueryRow(`insert into users(name, age) values('Scrooge McDuck', 93) returning id`)
			var userid int
			row.Scan(&userid)
	*/
}

func addTrigger(db *sql.DB, psqlInfo string, tableName string) {

	install := fmt.Sprintf(`
	CREATE OR REPLACE FUNCTION notify_%s()
	RETURNS trigger AS $$
	DECLARE
		payload json;
	BEGIN
		payload := json_build_object(
			'table', TG_TABLE_NAME,
			'oper' , TG_OP,
			'data' , row_to_json(NEW)
		);

		PERFORM pg_notify('%s_channel', payload::text);
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

	CREATE OR REPLACE TRIGGER trigger_%s
	AFTER INSERT OR UPDATE ON %s
	FOR EACH ROW
	EXECUTE FUNCTION notify_%s();
	`, tableName, tableName, tableName, tableName, tableName)

	_, err := db.Query(install)
	if err != nil {
		return
	}

	// Create a new listener
	listener := pq.NewListener(psqlInfo, 10*time.Second, time.Minute,
		func(ev pq.ListenerEventType, err error) {
			if err != nil {
				log.Println("Listener error:", err)
			}
		})

	// Subscribe to a channel
	err = listener.Listen(tableName + "_channel")
	if err != nil {
		return
	}

	log.Println("Listening for notifications...")
	for {
		select {
		case n := <-listener.Notify:
			if n != nil {
				fmt.Printf("Received notification: channel=%s payload=%s\n", n.Channel, n.Extra)
			}
		case <-time.After(90 * time.Second):
			// Keep the connection alive
			go listener.Ping()
		}
	}
	// log.Println("End Listening.")
}
