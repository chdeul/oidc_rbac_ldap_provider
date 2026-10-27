<script>
    app.component('ui-logout', {

        template : "#ui-logout",

        emits : [ "logged" ],
     
        data : function() {
            return {
                oidcLogoutUrl : ""
            }
        },
        mounted () {

            // Géré côté backend Native 
            myAxios.delete( "/auth").then( (rep) => {

                // puis Côté vue.js si succès backend 
                this.$emit( "logged", null, null, null ) 

                this.oidcLogoutUrl = rep.data
            })
        },
        
        methods : {

            oidcLogout() {
                // Puis OIDC Logout éventuel
                if ( this.oidcLogoutUrl ) {
                   window.location.href = this.oidcLogoutUrl 
                }                 
            },

            accueil() {
                router.push( "/vue/home" )
            }
        }

    })
//# sourceURL=http://localhost.vue/vues/Auth/logout.vue
</script>

<template id="ui-logout">
    <div class="col-md-6 offset-md-3 marge-haut">

            <h4>Déconnecté !</h4>  
        <br/>
        <div>
          <button type="button" class="btn btn-primary" @click="accueil">Accueil</button> &nbsp; 
          <button v-if="oidcLogoutUrl" type="button" class="btn btn-warning" @click="oidcLogout">Logout SSO</button>
        </div>
     </div>
</template>
