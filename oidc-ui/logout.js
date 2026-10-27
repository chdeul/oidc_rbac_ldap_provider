// Au démarrage de la page : init
function init() {

    // Query param
    const params = new URLSearchParams(window.location.search)
    const msg  = params.get("msg")
    if ( msg ) {
        // Message en rouge
        let err = document.getElementById("msg-erreur")
        err.innerHTML = msg
    }
    // All client logout 
    const urls  = params.get("urls")
    if ( urls ) {
        // "front Logout" des clients enregistrés du SSO 
        urlstab = urls.split("|")
        urlstab.forEach(url => {
            if ( url && url.startsWith("http") ) {
                fetch(url)
            }
        });
    }
}