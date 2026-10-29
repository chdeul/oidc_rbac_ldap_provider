<script>

    // DO NOT Touch : Value injected at backend !
    const ContextPath = ''

    // Config axios, appels ajax (xhr)
    var myAxios = axios.create({ baseURL: ContextPath+'/api'})

    // Transmet le token d'authentication sur chaque requete axios
    myAxios.defaults.withCredentials = true;

    // Traitement standard d'erreur appel ajax
    myAxios.interceptors.response.use(
        function(response) {
            return response
        },
        function(error) {
            // Non autorisé
            if (error.response && error.response.status === 401) {
                if ( router.currentRoute.value.path !== "/vue/login" ) {
                  message( "Vous n'êtes pas autorisé à accéder à cette requête !" )
                  router.push("/vue/login" )
                  return Promise.reject(error);
                }
            }
            //console.log( "axios intercept: ", error)
            if (error.response && error.response.status >= 500) {
                router.push("/vue/error?code="+error.response.status)
            }
            return Promise.reject(error);
        }
    )

    // Bouton de scrolling arrière
    const btnbtt = document.getElementById("backToTop")
    if (btnbtt) {
        window.addEventListener("scroll", () => {
        if (window.scrollY > 150) {
            btnbtt.classList.add("show");
        } else {
            btnbtt.classList.remove("show");
        }
        });
        // remonte 
        btnbtt.addEventListener("click", () => {
            window.scrollTo({ top: 0,behavior: "smooth" });
        });
    }
    // Create vuejs
    var app = Vue.createApp({

        template : "#ui-app",

        // App global data
        // access by this.$root.xxx
        data : function() {
            return {
                logged : null,
                profil : null
            }
        },

        created : function () {

            // Vérif si pas déjà encore logué via token
            const authData = this.decodeJwt()
            if ( authData ) {
                if ( authData.identifiant )
                    this.$root.majLogged( authData.identifiant, authData.profil )
                 else
                    this.$root.majLogged( null, null )
            }
        },

        methods : {

            // Decode le cookie JWT si présent
            decodeJwt() {

                // Essaie TODO Recupere dans Bearer 

                // Recupere le cookie token
                const name = "JWTAgremToken"
                const value = `; ${document.cookie}`
                let  parts = value.split(`; ${name}=`)
                if (parts.length !== 2) return null
                const token = parts.pop().split(";").shift()

                // Decode le token JWT
                parts = token.split(".")
                if (parts.length !== 3) return null
                const payload = parts[1]
                const decoded = atob(payload.replace(/-/g, "+").replace(/_/g, "/"))
                return JSON.parse(decoded)
            },

            // voir $Emit event logged par login/logout
            majLogged : function (user, profil) {
                // Sauve auth et profil
                this.logged = user
                this.profil = profil;
            },
            isAdmin() {
                return this.logged && this.profil === "Administration"
            },
            isEdit() {
                return this.logged && (this.profil === "Edition" || this.profil === "Administration")
            },
            isConsult() {
                return this.logged && (this.profil === "Consultation" || this.profil === "Edition" || this.profil === "Administration")
            }
        }
    })

    // Mixin sur toutes Pages authentifiées
    var authent = {
        created : function () {
            // Vérif si bien authentifié sur cette page
            if ( !this.$root.logged ) {
                // se relogue pour cette page authentifiée
                 router.push("/vue/login")
            } else {
                // Controle des droits sur cette page
                if ( ! this.$root.profil ) router.push("/vue/error")
            }
        },
        methods : {
            formatDate : function(date) {
                if ( date === undefined ) return "?"
                return moment(date).format('DD/MM/YYYY')
            } ,
            formatFloat :  function(nombre) {
                if ( nombre === undefined ) return "?"
                return nombre.toLocaleString('fr-FR', {minimumFractionDigits: 2, maximumFractionDigits: 2})
            },
            formatFloatI :  function(nombre) {
                if ( nombre === undefined ) return "?"
                return nombre.toLocaleString('fr-FR', {minimumFractionDigits: 0, maximumFractionDigits: 0})
            },
            formatFloat2 :  function(nombre) {
                if ( nombre === undefined ) return "?"
                if ( Math.abs(nombre) < 0.005 ) return ""
                return nombre.toLocaleString('fr-FR', {minimumFractionDigits: 2, maximumFractionDigits: 2})
            }
        }
    }

//# sourceURL=http://localhost:8080/vues/app.vue
</script>

<!-- template de la single page -->
<template id="ui-app">
    <div>
        <!-- header + main content -->
        <ui-header></ui-header>
        <div class="container-fluid">
            <div class="row app_padding">
                <router-view @logged="majLogged"></router-view>
            </div>
        </div>
    </div>
</template>

<style>
.app_padding {
    padding-top: 10px; 
    padding-bottom: 10px;
}
</style>