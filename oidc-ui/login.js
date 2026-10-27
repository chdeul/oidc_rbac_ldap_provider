// Au démarrage de la page : init
function init() {
    // Retransmet le Request ID dans les formdata
    const params = new URLSearchParams(window.location.search)
    const reqid  = params.get("requestId")
    if ( reqid ) document.getElementById("request-id").value = reqid
    // Clean error message
    let err = document.getElementById("msg-erreur")
    err.innerText = ""
    const error  = params.get("error")
    if ( error ) err.innerText = error
    let souv = document.getElementById("souvenir")
    let ident = document.getElementById("identifiant")
    let form = document.getElementById("formulaire")
    let oublie = document.getElementById("oublie")
    let but = document.getElementById("connect-but")
    let lib = document.getElementById("libelle-mdp")
    let mdp = document.getElementById("motdepasse")
    // Cas on a memorisé l'identifiant
    const saved = localStorage.getItem("savedIdent")
    if (saved) {
        // Préremplit identifiant et mot de passe (nav)
        form.autocomplete = "on"
        ident.autocomplete = "username"
        mdp.autocomplete = "current-password"
        ident.value = saved
        souv.checked = true
    } else {
        form.autocomplete = "off"
        ident.autocomplete = "off"
        mdp.autocomplete = "new-password"
        souv.checked = false
    }
    // cas réception lien mail renew
    const renew  = reqid.startsWith("agnew-")
    if ( renew ) {
        // Nouveau mot de passe à saisir
        oublie.checked = false
        oublie.disabled = true
        lib.innerText = "Nouveau mot de passe"
        but.innerText = "Enregistrer"
        mdp.value = ""
        mdp.autocomplete = "new-password"
        mdp.placeholder = "Veuillez saisir un nouveau mot de passe..."
        err.innerText = "Attention : vous allez changer de mot de passe !"
        form.action="./auth/renew"            
    }
}

// Toggle mot de passe visible (sauf si oublié coché)
function montrePswd() {
    const oublie = document.getElementById("oublie").checked
    let mdp = document.getElementById("motdepasse")
    if (mdp.type === "password" && !oublie) mdp.type = "text";
    else mdp.type = "password"
}

// Toggle Se souvenir des données : les sauvegardes ou purge
function souviensToi() {
    const ident = document.getElementById("identifiant")
    const souv = document.getElementById("souvenir").checked
    let form = document.getElementById("formulaire")
    let mdp = document.getElementById("motdepasse")
    if (souv) {
        form.autocomplete = "on"
        ident.autocomplete = "username"
        mdp.autocomplete = "current-password"
        localStorage.setItem("savedIdent", ident.value)
    } else {
        form.autocomplete = "off"
        ident.autocomplete = "off"
        mdp.autocomplete = "new-password"
        localStorage.removeItem("savedIdent")
    }
}

// Toggle Oublié : oui envoi mail, non login normal 
function renouvellePswd(obj) {
    let renew = obj.checked
    let mdp = document.getElementById("motdepasse")
    let but = document.getElementById("connect-but")
    let form = document.getElementById("formulaire")
    mdp.type = "password"
    if ( renew ) {
        mdp.readonly = true
        mdp.value = "!Demande!Renouvellement!"
        but.innerText = "Recevoir un email"
        form.action="./auth/renew"
    } else {
        mdp.readonly = false
        mdp.value = ""
        but.innerText = "Connexion"
        form.action="./auth/login"
    }
}

// Valide le formulaire
function valider() {

        // Reset pas d'erreur
        let err = document.getElementById("msg-erreur")
        err.innerText = ""

        // Contole des données sasie et cohérence
        const oublie = document.getElementById("oublie").checked;
        const ident = document.getElementById("identifiant").value.trim()
        const mdp = document.getElementById("motdepasse").value.trim()

        // Si Oublié : rappel email obligatoire (pas le simple identifiant)
        if (oublie && ident === "") {
            err.innerText = "Veuillez saisir l'e-mail associé au compte, ou l'identifiant !"
            return false
        }
        // Mot de passe ou token non vide
        if ( mdp === "" && !oublie.enabled ) {
            err.innerText = "Veuillez saisir un mot de passe !"
            return false
        }             
        // Mot de passe ou token non vide
        if ( mdp === "" ) {
            err.innerText = "Veuillez saisir un mot de passe ou token !"
            return false
        }
        // Identifiant ou mail non vide ( sauf si token dans mot de passe )
        if (ident === "" && !mdp.startsWith("agta-") ) {
            err.innerText = "Veuillez saisir un identifiant ou e-mail !"
            return false
        }
        // Mot de passe taille
        if ( mdp.length < 8) {
            err.innerText = "Veuillez saisir un mot de passe ou token valide (trop court) !"
            return false
        }
        // Identifiant ou mail  valide ( chars used )
        if (  ident.includes("@")  
        &&    ident.match("/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{3,132}$/")) {
            err.innerText = "Adresse mail invalide ou comportant des caractères interdits !"
            return false
        } else if ( ident.match("/^[a-zA-Z][a-zA-Z0-9._-]{2,64}$/") ) {
            err.innerText = "Identifiant invalide ou comportant des caractères interdits !"
            return false
        }

        // Resauve éventuellement les données si changé et demandé
        souviensToi()
        return true
}
