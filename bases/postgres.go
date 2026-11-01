package bases

import (
	"database/sql"
	"fmt"
	"main/tools"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

// TODO  CRUD object : table row : with config file for request,  unshalling
// TODO  Flyway like (sql_migrate) Create roles, database, schema,  Drop
// TODO  config Pool of connections

// Add trigger on table change, to be notified as replacement for webhook
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

		PERFORM pg_notify('%s_pgsqlchannel', payload::text);
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

	CREATE OR REPLACE TRIGGER trigger_%s
	AFTER INSERT OR UPDATE ON %s
	FOR EACH ROW
	EXECUTE FUNCTION notify_%s();
	`, tableName, tableName, tableName, tableName, tableName)

	_, err := db.Query(install)
	tools.CheckError(err, "Cannot set postgresql trigger : %v", err)

	// Create a new listener
	listener := pq.NewListener(psqlInfo, 10*time.Second, time.Minute,
		func(ev pq.ListenerEventType, err error) {
			tools.CheckError(err, "Postgresql Listener error: %v", err)
		})

	// Subscribe to a channel
	err = listener.Listen(tableName + "_pgsqlchannel")
	tools.CheckError(err)

	tools.Info("Listening for Postgresql notifications...")
	for {
		select {
		case n := <-listener.Notify:
			if n != nil {

				// TODO   call a callback
				tools.Debug("Received Postgresql notification: channel=%s payload=%s", n.Channel, n.Extra)
			}
		case <-time.After(90 * time.Second):
			// Keep the connection alive
			go listener.Ping()
		}
	}
}
