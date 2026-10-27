package interf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"main/tools"
	"maps"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

// Configuration de la transformation ApiRest quelconque -> Api Annuaire Referentiel

type RestInputConfig struct {
	HeadAuth string
	BaseUrl  string
	CertFile string
	Insecure bool
	Routes   []RouteConfig
}

type RouteConfig struct {
	Source struct {
		Method string
		Path   string
		Header map[string]string
	}

	Target []struct {
		Method    string
		URL       string
		Header    map[string]string
		Body      map[string]any
		JsScript  string
		MapArray  string
		MapObject map[string]any
	}
}

// Charge la configuration
func LoadInputConfig(path string) (*RestInputConfig, error) {
	var cfg RestInputConfig
	err := tools.LoadAnyConfig(path, &cfg)
	return &cfg, err
}

func escape(str string) string {
	str = strings.ReplaceAll(str, "?", "\\?")
	str = strings.ReplaceAll(str, "+", "\\+")
	str = strings.ReplaceAll(str, ".", "\\.")
	str = strings.ReplaceAll(str, "*", "\\*")
	str = strings.ReplaceAll(str, "(", "\\(")
	str = strings.ReplaceAll(str, ")", "\\)")
	return str
}

// Recupère les paramètres appelés : {param}
func extractParams(pattern, path string, hpatt map[string]string, header any) map[string]any {

	// recherche les {xxx}
	re := regexp.MustCompile(`\{(\w+)\}`)
	keys := re.FindAllStringSubmatch(escape(pattern), -1)

	// regex de substitution : recherche les valeurs de {xxx}
	regexPattern := "^" + re.ReplaceAllString(escape(pattern), `([^/]+)`) + "$"
	re2 := regexp.MustCompile(regexPattern)
	values := re2.FindStringSubmatch(path)

	// renvoie une map [param] valeur
	params := map[string]any{}
	if len(values) > 1 {
		for i, key := range keys {
			params[key[1]] = values[i+1]
		}
	}
	// type de headers ?
	mh, ok := header.(map[string]any)
	if ok { // Type map
		// Fusionne avec params
		maps.Copy(params, mh)
	}
	// Pas d'autres headers
	if len(hpatt) == 0 {
		return params
	}

	// Idem pour chaque valeur de header trouvé
	for h, v := range hpatt {

		// Recupere la valeur du header passé
		var value string
		mh, ok := header.(map[string]string)
		if ok { // Type map
			value = mh[h]
		} else {
			hh, ok := header.(http.Header)
			if ok { // Type http
				value = hh.Get(h)
			}
		}
		// S'il existe une valeur
		if value != "" {
			// Récupère les valeurs de paramètres
			keys = re.FindAllStringSubmatch(v, -1)
			regexPattern = "^" + re.ReplaceAllString(v, `([^/]+)`) + "$"
			re2 = regexp.MustCompile(regexPattern)
			values = re2.FindStringSubmatch(value)
			if len(values) > 1 {
				for i, key := range keys {
					params[key[1]] = values[i+1]
				}
			}
		}
	}
	// Cookie : non , en REST Api stateless peu de cookie utilisés
	return params
}

// Convertit les paramm appelés , en param ou body appelant
func applyParamsToBody(body map[string]any, params map[string]any) map[string]any {
	result := map[string]any{}
	for k, v := range body {
		var data any
		if str, ok := v.(string); ok {
			for key, val := range params {
				if sval, ok2 := val.(string); ok2 {
					str = strings.ReplaceAll(str, "{"+key+"}", sval)
				} else if str == "{"+key+"}" {
					data = val
				}
			}
			if data == nil {
				result[k] = str
			} else {
				result[k] = data
			}
		} else {
			result[k] = v
		}
	}
	return result
}

// Convertit les paramm appelés , en param ou body appelant
func applyParamsToString(urlPath string, params map[string]any, enc bool) string {
	str := urlPath
	for key, val := range params {

		if sval, ok := val.(string); ok {
			if enc {
				sval = url.QueryEscape(sval)
			}
			str = strings.ReplaceAll(str, "{"+key+"}", sval)
		}
	}
	return str
}

func mapResponse(raw []byte, maparray string, mapobject map[string]any) any {

	res := gjson.ParseBytes(raw)
	if maparray == "" && !res.IsArray() {
		return mapObject(res, mapobject)
	}
	var tab []gjson.Result
	if maparray == "" {
		tab = res.Array()
	} else {
		tab = res.Get(maparray).Array()
	}
	ret := make([]map[string]any, len(tab))
	for i, obj := range tab {
		// Extraie et reconstruit
		ret[i] = mapObject(obj, mapobject)
	}
	return ret
}

// Transformation de la réponse en réponse json pour l'appelant
func mapObject(obj gjson.Result, mapping map[string]any) map[string]any {

	// Extraie et reconstruit
	ret := map[string]any{}
	for outKey, jsonPath := range mapping {
		var v any
		sjp, ok := jsonPath.(string)
		if !ok {
			v = jsonPath
		} else {
			o := obj.Get(sjp)
			t := o.Type
			switch t {
			case gjson.False:
				v = "false"
			case gjson.Number:
				v = strconv.FormatFloat(o.Num, 'f', -1, 64)
			case gjson.String:
				v = o.Str
			case gjson.True:
				v = "true"
			case gjson.JSON:
				v = o.Raw
			default:
				v = jsonPath
			}
		}
		ret[outKey] = v
	}
	return ret
}

// Transformation de la réponse en réponse json pour l'appelant
func mapArray(raw []byte, mapping map[string]any) []map[string]any {

	tab := gjson.ParseBytes(raw).Array()
	tabres := make([]map[string]any, len(tab))
	for i, obj := range tab {
		// Extraie et reconstruit
		result := map[string]any{}
		for outKey, jsonPath := range mapping {
			sjp, ok := jsonPath.(string)
			if !ok {
				result[outKey] = jsonPath
			} else {
				result[outKey] = obj.Get(sjp).Str
			}
		}
		tabres[i] = result
	}
	return tabres
}

// Fait l'appel(s) REST désigné(s) avec conversion du résultat
func RestInputFct(cfg *RestInputConfig, method, path string, values map[string]any) any {

	// Pour chaque appel rest Annuaire Referentiel ( get Users )
	for _, route := range cfg.Routes {

		// Prefixe du chemin demandé ex: users/toto/{id}?page={p}   jusque { ou ?
		prefix := route.Source.Path
		idx := strings.IndexAny(route.Source.Path, "{?}")
		if idx >= 0 {
			prefix = prefix[:idx]
		}
		// Requete ne correspondant pas à celle demandée
		if route.Source.Method != method || !strings.HasPrefix(path, prefix) {
			continue
		}
		// Extraie les valeurs des paramètres id et p
		params := extractParams(route.Source.Path, path, route.Source.Header, values)
		var mapped any

		// Enchaine les  requetes à appeler (si plusieurs)
		for _, target := range route.Target {

			// Extrait les param de mapped appel précédent
			mp, ok := mapped.(map[string]any)
			if ok {
				for m, v := range mp {
					params[m] = v
				}
			}
			// Si au moins un paramètres
			var targetUrl string = target.URL
			var jsonBody []byte = nil
			if len(params) > 0 {
				if target.Body != nil {
					// Construire le body cible en injectant des paramètres
					targetBody := applyParamsToBody(target.Body, params)
					jsonBody, _ = json.Marshal(targetBody)
					tools.Debug("body %s", jsonBody)
				}
				// Construire l'url cible en injectant des paramètres
				targetUrl = applyParamsToString(target.URL, params, true)
				tools.Debug("url %s", targetUrl)
			}
			// Appel API cible
			req, _ := http.NewRequest(target.Method, cfg.BaseUrl+targetUrl, bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			if cfg.HeadAuth != "" {
				req.Header.Set("Authorization", cfg.HeadAuth)
			}
			for h, v := range target.Header {
				val := applyParamsToString(v, params, false)
				req.Header.Set(h, val)
				tools.Debug("header %s: %s", h, val)
			}
			// Attends la réponse
			tools.Debug("http do %s", req.URL.RequestURI())
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return err
			}
			if resp.StatusCode > 299 {
				return fmt.Errorf("Rest input http error %d", resp.StatusCode)
			}
			// Lit les data json de la réponse
			rawResp, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			// Applique une conversion javascript eventuelle
			if target.JsScript != "" {
				rawResp = executeJsScript(target.JsScript, rawResp)
			}

			// Transformer la réponse
			mapped = mapResponse(rawResp, target.MapArray, target.MapObject)
		}
		// vers l'appelant
		return mapped
	}
	// Ne matche aucune url : réponse vide
	return nil
}
