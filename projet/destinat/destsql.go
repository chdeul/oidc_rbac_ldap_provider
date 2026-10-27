package destinat

import obj "main/objets"

/*
 * Base de données SQL avec tables : projects, users, rbac
 */

type DestSql struct {
}

func (de *DestSql) Name() string {
	return "sqldb"
}

func NewDestSql() (*DestSql, error) {
	return nil, nil
}

func (de *DestSql) Close() {

}

// TODO  sql projects methodes

func (de *DestSql) SearchUsersRole(prj *obj.Project) []*obj.User {
	return nil
}

func (de *DestSql) SearchProjects(projet string) []*obj.Project {
	return nil
}

// TODO  sql users methodes

func (de *DestSql) SearchUsers(ident string) []*obj.User {
	return nil
}

func (de *DestSql) CreateUser(user *obj.User) error {
	return nil
}

func (de *DestSql) UpdateUser(user *obj.User) error {
	return nil
}

func (de *DestSql) DeleteUser(ident *obj.User) error {
	return nil
}

// TODO Sql role methodes

func (de *DestSql) GetRole(prj *obj.Project, user *obj.User) string {
	return ""
}

func (de *DestSql) CreateRole(prj *obj.Project, user *obj.User, role string) error {
	return nil
}

func (de *DestSql) UpdateRole(prj *obj.Project, user *obj.User, role string) error {
	return nil
}

func (de *DestSql) DeleteRole(prj *obj.Project, user *obj.User) error {
	return nil
}
