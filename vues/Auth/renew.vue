<script>
    app.component('ui-renew', {

        template : "#ui-renew",

        data () {
            return {
                user : {
                    Identifiant : "",
                    AdresseMail  : "",
                    MotDePasse : "",
                    NomComplet : ""
                },
                MotDePasse2 : "",
                Clef : null,
                erreur : ""
            }
        },

        mounted () {
           // Clé de sécurité du renouvellement de mot de passe, éventuelle
           this.Clef =  this.$route.query.clef 
           this.NomComplet = this.Clef 
        },

        methods : { 
            renew : function() {

                // reset message d'erreur
                this.erreur = "";

                // Si on vient du lien dans email de renouvellement
                if ( this.Clef ) {

                    // Verif mots de passe égaux 
                    if ( this.user.MotDePasse !== this.MotDePasse2 ) {
                       this.erreur = "Mots de passe différents, ré-essayez !" 
                       return
                    }
                    // Transmet la Clef dans le champ email, et le mot de passe
                    this.user.NomComplet = this.Clef 
                } else {
                // Sinon, 
                    // demande initiale d'email de renouvellement
                    this.user.MotDePasse = "$Demande$Renouvellement$" 
                }
                // Appel serveur de renouvellement
                myAxios.put("/auth", this.user).then( (rep) => {

                    // Si Mail envoyé 
                    if ( !this.Clef ) {

                        // Dialog pour demander de vérifier les mails
                        message( rep.data )

                        // Puis page d'accueil
                        router.push( "/vue/home")

                    // Si on vient du lien dans email de renouvellement
                    } else {

                        // Dialog pour confirmer le bon renouvellement
                        message( rep.data )

                        // Puis login
                        router.push( "/vue/login")
                    }

                }).catch( (err) => {

                    this.erreur = "Erreur ou Incohérence des données"
                    console.log(err);
                })
            },
            cancel : function() {
                 router.push( "/vue/home")
            }
        }  
    })
//# sourceURL=http://localhost.vue/vues/Auth/renew.vue
</script>

<template id="ui-renew">
    <div class="col-md-6 offset-md-3 marge-haut">

        <div>
        <h4>Réinitialisation du mot de passe</h4>
        <p v-if="Clef">Vous avez oublié ou perdu votre mot de passe, veuillez saisir votre identifiant 
        de compte et un nouveau mot de passe pour ce compte (8 Caractères avec Maj, min, chiffre, symbole)... </p>    
        <p v-else>Vous avez oublié ou perdu votre mot de passe, veuillez saisir votre identifiant 
        de compte et l'adresse mail enregistrée pour ce compte...</p>   
        </div>

        <form @submit.prevent="renew">
            <div class="mb-3 mt-3">
                <label for="Identifiant" class="form-label">Identifiant:</label>
                <input type="text" class="form-control" v-model="user.Identifiant" id="Identifiant" required>
            </div>
            <div v-if="Clef" class="mb-3">
                <label for="MotDePasse" class="form-label">Mot de passe:</label>
                <input type="password" class="form-control" v-model="user.MotDePasse" id="MotDePasse" required>
            </div>
            <div v-if="Clef" class="mb-3">
                <label for="MotDePasse2" class="form-label">Mot de passe (confirmation):</label>
                <input type="password" class="form-control" v-model="MotDePasse2" id="MotDePasse2" required>
            </div>
            <div v-else class="mb-3">
                <label for="AdresseMail" class="form-label">Adresse mail:</label>
                <input type="email" class="form-control" v-model="user.AdresseMail" id="AdresseMail" required>
            </div>
            <div v-show="erreur" class="mb-3 red">
                <span>{{erreur}}</span>
            </div>
            <button type="submit" class="btn btn-primary">Envoyer</button> &nbsp;
            <button type="button" class="btn btn-secondary" @click="cancel">Annuler</button>
        </form>
    </div>
</template>
