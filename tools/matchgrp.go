package tools

import (
	"regexp"
	"strings"
)

const MatchChars string = "%_.*"
const Separators string = "| \t\n"

func InitMatch(ident string) *regexp.Regexp {

	// Rien à chercher : retourne vide
	if ident == "" {
		return nil
	}
	ident = strings.ReplaceAll(ident, "%", "*")
	ident = strings.ReplaceAll(ident, "_", ".")
	ident = strings.ReplaceAll(ident, "*", ".*")
	ident = "^" + ident + "$"

	re := regexp.MustCompile(ident)
	return re
}

// Methode d'interface
func Match(re *regexp.Regexp, values string) bool {

	// Rien à chercher : retourne vide
	if re == nil {
		return false
	}
	mots := strings.Split(values, Separators)
	for _, m := range mots {
		// cherche une valeur
		if m != "" && re.MatchString(m) {
			return true
		}
	}
	// non trouvé
	return false
}

func IsIdent(ident string) bool {
	//  a-z 0-9 _ et taille 2 à 63 caractères
	if len(ident) < 2 || len(ident) > 63 {
		Warn("Taille identifiant incorrect %s", ident)
		return false
	}
	for _, r := range ident {
		if !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_') {
			Warn("Identifiant incorrect %s", ident)
			return false
		}
	}
	return true
}

func IsGroup(ident string) bool {
	return IsIdent(ident)
}

func IsPattern(ident string) bool {
	// TODO a-z A-Z 0-9 _ . * et taille 1 à 63 caractères
	// Warn( "Recherche incorrecte %s", ident )
	return true
}

// TODO  isProject  with /  or .
