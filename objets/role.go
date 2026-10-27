package obj

import "strings"

var (
	// Niveaux de droits, profils d'utilisateur
	Publique       string = "Publique"
	Consultation   string = "Consultation"
	Edition        string = "Edition"
	Administration string = "Administration"
)

func GetProfil(ident, roles string) string {

	// admin ==> Administration
	if ident == "admin" {
		return Administration
	}
	// Must have roles
	if strings.Contains(roles, Administration) {
		return Administration
	}
	if strings.Contains(roles, Edition) {
		return Edition
	}
	// Authentifié ==> role par défaut  Consultation
	return Consultation
}

// Role nécessaire pour une requete api / Roles de l'utilisateur
func ProfilInclus(apiRole, userRoles string) bool {

	// Tout api administraiton
	if strings.Contains(userRoles, Administration) {
		return true
	}
	// api Edition ou consultation
	if strings.Contains(userRoles, Edition) {
		if apiRole == Edition || apiRole == Consultation {
			return true
		}
	}
	// api consultation
	if strings.Contains(userRoles, Consultation) {
		if apiRole == Consultation {
			return true
		}
	}
	// aucun role
	return false
}
