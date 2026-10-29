<script>
    app.component('ui-refbrowse', {

        template : "#ui-refbrowse",
        mixins : [ authent ],

        data () {
            return {

                // Recherche
                pattern : "*",
                isGroup : false,
                erreur : "",

                // Résultat recherche
                users  : [],
            }
        },

        methods : { 

            search : function() {

                // Reset msg d'erreur
                this.erreur = ""
                this.user = {}

                // Que chercher
                let query = { }
                if ( this.isGroup ) {
                    query.group = this.pattern
                } else  {
                    query.user = this.pattern 
                }

                // Appel serveur search ref users
                myAxios.get("/annref", { 
                    params: query 
                }).then( (rep) => {
                    this.users = rep.data
                }).catch( (err) => {
                    this.erreur = "Erreur appel REST : " + err
                })
            },
            differe : function(u,user) {
                let diff = "Utilisateur " + user.NomComplet + " existant dans Agréments\n"
                if ( u.NomComplet !== user.NomComplet )
                    diff = diff + "\n - Le nom différe : " + u.NomComplet 
                if ( user.AdresseMail !== user.AdresseMail )
                        diff = diff + "\n - Le mail différe : " + u.AdresseMail 
                if ( user.Uid !== user.Uid )
                        diff = diff + "\n - L'ID différe : " + u.Uid
                if ( user.Roles !== user.Roles )
                    diff = diff + "\n - Les roles différent : " + u.Roles                    
                return diff
            },
            creerUser : function( user ) {
                confirme( "Voulez-vous créer l'utilisateur " + user.NomComplet + " dans Agréments ?", () => {
                    // Nouvel utilisateur
                    query = { }
                    query.user = true
                    myAxios.post("/ldap", user, { 
                        params: query 
                    }).then( (rep) => {
                        message( "Utilisateur " + user.NomComplet + " créé dans Agréments" )
                    }).catch( (err) => {
                        this.errval = "Erreur appel REST : " + err
                    })
                })
            },
            toLdap : function(user) {

                // Que chercher
                let query = { }
                query.user = user.Identifiant 
              
                // Appel serveur search ldap users
                myAxios.get("/ldap", { 
                    params: query 
                }).then( (rep) => {
                    let users = rep.data
                    if ( users.length == 1 ) {
                        message( this.differe( users[0], user ) )
                    } else if ( users.length > 1 ) { 
                        this.erreur = "Plusieurs utilisateurs dupliqués dans le ldap !" 
                    } else {
                        this.creerUser( user )
                    }
                }).catch( (err) => {
                    this.erreur = "Erreur appel REST : " + err
                })  
            }
        }      
    })
//# sourceURL=http://localhost:8080/vues/Admin/refbrowse.vue
</script>

<template id="ui-refbrowse">
    <div class="col-md-10 offset-md-1 marge-haut" >

        <div class="col-md-8">
            <h3>Parcourir l'annuaire</h3>
            <div>

                <div class="input-group mb-3">
                    <label class="input-group-text" for="grp-id">Role :</label>
                    <span class="input-group-text"><input class="form-check-input mt-0" type="checkbox" v-model="isGroup" id="grp-id"></span>
                    <label class="input-group-text" for="nom-id"> &nbsp; &nbsp; &nbsp; Identifiant :</label>
                    <input type="text" class="form-control" placeholder="pattern avec . _ * ou %" v-model="pattern" id="nom-id">
                    <span class="input-group-text btn btn-primary" @click="search">Chercher</span>
                </div>

            </div>
            <div class="m-3 text-danger">{{ erreur  }}</div>          
        </div>
        <br/>
        <br/>
        <div>
             <h5>Utilisateurs</h5>
             <table class="table table-striped">
                <thead>
                    <tr>
                        <th scope="col">Id</th>
                        <th scope="col">Identifiant</th>
                        <th scope="col">Nom complet</th>
                        <th scope="col">Adresse mail</th>
                        <th scope="col">Roles</th>
                   </tr>
                </thead>
                <tbody>
                    <tr v-for="u in users" :key="u.sn" @dblclick="toLdap(u)">
                        <td scope="row">{{u.Uid}}</td>
                        <td>{{u.Identifiant}}</td>
                        <td>{{u.NomComplet}}</td>
                        <td>{{u.AdresseMail}}</td>
                        <td>{{u.Roles}}</td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>

</template>

<style>

</style>