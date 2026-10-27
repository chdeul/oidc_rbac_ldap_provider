<script>
    app.component('ui-login', {

        template : "#ui-login",

        emits : [ "logged" ],

        data () {
            return {
                cfg : { 
                    Form: false, 	
                    Oidc: false,
                    OidcPath: ""
                },
                formData : {
                    Identifiant : "",
                    MotDePasse  : "",
                    RequestId : ""
                },
                isSouv : false,
                erreur : ""
            }
        },
        created () {
            // Appel serveur recup config
            myAxios.get("/config?login").then( (rep) => {
                this.cfg = rep.data
                if ( this.cfg.Oidc && !this.cfg.Form ) {
                    window.location.replace( this.cfg.OidcPath )
                }
            }).catch( (err) => {
                this.erreur = "Le serveur ne répond pas !"
            })
        },

        mounted () {
            // On stocke le requestId Oidc dans le nom complet
            this.formData.RequestId = this.$route.query.requestID

            let saved = localStorage.getItem("savedIdent")
            this.isSouv = (saved != null)
            if ( saved != null ) {
                this.formData.Identifiant = saved
             } else {
                this.formData.Identifiant = ""
             }
        },

        methods : { 
            login : function() {

                // Reset msg d'erreur
                this.erreur = ""

                // Appel serveur login
                myAxios.post("/auth", this.formData ).then( (rep) => {

                    let user = rep.data

                    // OK, stoke le token (mot de passe) et maj user logué
                    this.$emit( "logged", user.Identifiant, user.Profil )

                    // Se souvenir update identifiant if neccessary
                    this.souvenir( this.isSouv )

                    // Page d'accueil
                    router.push( "/vue/home")

                }).catch( (err) => {

                    this.erreur = "Identifiant ou mot de passe incorrect !"
                })
            },
            cancel : function() {
                // Retourne à la page publique
                router.push( "/vue/home")
            },

            sso : function() {
                			
                // Redirigine ici après authent SSO
			    const path = this.$root.contextPath + encodeURIComponent("/index.html?route=/vue/home")

                // Se logue sur le sso
                window.location = this.cfg.OidcPath + "?redirect=" + path
            },

            autoComplete( event )
            {
                // Choix souviens-toi
                const checked = event.target.checked
                this.souvenir( checked )
            },

            souvenir( checked ) {

                if (checked) {
                    // Oui:  autocomplete pour password, et stockage en local identifiant
                    localStorage.setItem("savedIdent", this.formData.Identifiant)
                } else {
                    // Non Efface 
                    localStorage.removeItem("savedIdent")
                }
            },

            montrePswd()
            {
                this.$refs.MotDePasse.type = (this.$refs.MotDePasse.type === "password") ? "text" : "password"
            }
        } 
    })
//# sourceURL=http://localhost.vue/vues/Auth/login.vue
</script>

<template id="ui-login">
    <div class="col-md-4 offset-md-4 marge-haut" >
        <div v-if="cfg.Form">
            <h4>Authentification</h4>
        </div>
        <form v-if="cfg.Form" @submit.prevent="login" :autocomplete="isSouv?'on':'off'">
            <div class="mb-3 mt-3">
                <label for="Identifiant" class="form-label">Identifiant:</label>
                <input type="text" class="form-control" v-model="formData.Identifiant" 
                       id="Identifiant" required :autocomplete="isSouv?'username':'off'">
                <input type="hidden" v-model="formData.RequestId">
            </div>             
            <div class="mb-3" style="position:relative">
                <label for="MotDePasse" class="form-label">Mot de passe:</label>            
                <input type="password" class="form-control" v-model="formData.MotDePasse" 
                       id="MotDePasse" ref="MotDePasse" required :autocomplete="isSouv?'current-password':'new-password'">
                <span class="fa fa-lg fa-eye password-eye" @click="montrePswd"></span>
            </div>
            <div class="mb-3">
                <router-link to="/vue/renew">Mot de passe perdu...</router-link> 
                <input type="checkbox" checked @change="autoComplete" v-model="isSouv"
                       id="SeSouvenir" style="float: right;margin-top:5px">
                <label for="SeSouvenir" class="form-label" 
                       style="float: right">Se souvenir &nbsp;</label>
            </div>
            <div v-show="erreur" class="mb-3 red">
                <span>{{erreur}}</span>
            </div>        
            <button type="submit" class="loginmarge btn btn-primary">Connexion</button>  &nbsp;
            <button type="button" class="loginmarge btn btn-secondary" @click="cancel">Annuler</button> &nbsp;
            <span v-if="cfg.Oidc && !formData.RequestId" class="oidc-panel">
                <button type="button" class="loginmarge btn btn-warning" @click="sso">Login SSO</button>
            </span>
        </form>
    </div>
</template>

<style>
.password-eye {
 color:lightgrey;
 padding-left:1em;
 position:absolute;
 top:2em;
 right:0.5em;
}
.oidc-panel {
    float: right;
}
.loginmarge {
    margin-bottom: 8px;
}
</style>