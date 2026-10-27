package referent

import (
	"database/sql"
	"fmt"
	obj "main/objets"
	"main/tools"

	_ "github.com/lib/pq"
)

/*
 * Liste d'utilisateurs chargée à partir d'une base sql
 *  configuréeà partir de sqlInput.yml   dans ./data
 */

// Configuration postgres
type SqlConfig struct {
	Database struct {
		Driver   string
		Host     string
		Port     string
		User     string
		Password string
		Dbname   string
	}
	Annuary struct {
		Users string
		Group string
	}
}

type RefSql struct {
	conn *sql.DB
	cfg  *SqlConfig
}

func NewRefSql(filename string) (*RefSql, error) {

	// Config yml : doit mapper sur un tableau de User
	var cfgSql SqlConfig
	err := tools.LoadAnyConfig(tools.Cfg.Ref.Path, &cfgSql)
	if err != nil {
		return nil, err
	}
	ret := RefSql{cfg: &cfgSql}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		ret.cfg.Database.Host, ret.cfg.Database.Port, ret.cfg.Database.User, ret.cfg.Database.Password,
		ret.cfg.Database.Dbname)

	ret.conn, err = sql.Open(ret.cfg.Database.Driver, psqlInfo)
	if err != nil || ret.conn == nil {
		return nil, err
	}
	return &ret, nil
}

func (rf *RefSql) Close() {
	if rf.conn != nil {
		rf.conn.Close()
	}
	rf.cfg = nil
}

// Methode d'interface
func (rf *RefSql) SearchUser(id string) []obj.User {

	// Rien à chercher : retourne vide
	if id == "" {
		return nil
	}
	rows, err := rf.conn.Query(rf.cfg.Annuary.Users, id)
	if err != nil {
		return nil
	}
	defer rows.Close()

	result := []obj.User{}
	for rows.Next() {
		var user obj.User
		rows.Scan(&user.Uid, &user.Identifiant, &user.MotDePasse, &user.NomComplet, &user.AdresseMail)
		result = append(result, user)
	}
	return result
}

// Methode d'interface
func (rf *RefSql) SearchGroup(id string) []obj.User {

	// Rien à chercher : retourne vide
	if id == "" {
		return nil
	}
	rows, err := rf.conn.Query(rf.cfg.Annuary.Group, id)
	if err != nil {
		return nil
	}
	defer rows.Close()

	result := []obj.User{}
	for rows.Next() {
		var user obj.User
		rows.Scan(&user.Uid, &user.Identifiant, &user.MotDePasse, &user.NomComplet, &user.AdresseMail)
		result = append(result, user)
	}
	return result
}
