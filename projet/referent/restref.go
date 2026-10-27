package referent

import (
	"fmt"
	"main/interf"
	obj "main/objets"
	"main/tools"
)

/*
 * Liste d'utilisateurs chargée à partir d'une api rest
 *  configurée depuis  restApi.yml    dans ./data
 */

type RefRest struct {
	cfg *interf.RestInputConfig
}

func NewRefRest(filename string) (*RefRest, error) {

	// Config yml : doit mapper sur un tableau de User
	cfg, err := interf.LoadInputConfig(tools.Cfg.Ref.Path)
	return &RefRest{cfg: cfg}, err
}

func (rf *RefRest) Close() {
	rf.cfg = nil
}

// Methode d'interface
func (rf *RefRest) SearchUser(id string) (table []obj.User) {

	// Rien à chercher : retourne vide
	if id == "" {
		// Erreur interne
		return nil
	}
	// Catch exception de conversion
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("RestInput mal configuré %v", r)
			table = nil
		}
	}()
	// Appel rest avec remapping dynamique des données
	ret := interf.RestInputFct(rf.cfg, "GET", "users/"+id, nil)

	// Conversion array of map[string]string to array of User
	return convertToUsers(ret)
}

// Methode d'interface
func (rf *RefRest) SearchGroup(gr string) (table []obj.User) {

	// Rien à chercher : retourne vide
	if gr == "" {
		return nil
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("RestInput mal configuré %v", r)
			table = nil
		}
	}()

	// Appel rest avec remapping
	ret := interf.RestInputFct(rf.cfg, "GET", "group/"+gr, nil)

	// Conversion array of map[string]string to array of User
	return convertToUsers(ret)
}

// Conversion de array of any/map[string]any/string to array of User
func convertToUsers(ret any) []obj.User {

	// Conversion array
	tab := ret.([]any)

	// Allocation Tableau résultat
	result := make([]obj.User, len(tab))

	// Chaque user trouvé ( en général 1 élément )
	for i, t := range tab {
		// Map des champs mappés de User
		m := t.(map[string]any)
		result[i].Uid = m["Uid"].(string)
		result[i].Identifiant = m["Identifiant"].(string)
		result[i].NomComplet = m["NomComplet"].(string)
		result[i].AdresseMail = m["AdresseMail"].(string)
		result[i].Roles = m["Roles"].(string)
		// On ne remonte pas le MotDePasse éventuel
	}
	// Résultat de la fonction
	return result
}
