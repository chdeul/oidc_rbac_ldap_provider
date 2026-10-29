<script>
    app.component('ui-projet', {

        template : "#ui-projet",
        
        data () {
            return {

                erreur : ""
            }
        },

        methods : { 

            prendre: function() {

               // Appel serveur 
                myAxios.get("/projet").then( (rep) => {

                }).catch( (err) => {

                    this.erreur = "Projet non pris en compte !"
                })

            }
        }

     })
//# sourceURL=http://localhost:8080/vues/Admin/projet.vue
</script>

<template id="ui-projet">
    <div class="col-md-12 marge-haut" >
        <ui-tableau />
    </div>
    <div class="col-md-12 marge-haut" >
        <div class="col-md-3">
           &nbsp;  &nbsp; &nbsp; <button class="btn btn-warning" @click="prendre">Prendre en compte</button>
           <div v-show="erreur" class="mb-3 red">
                <span>{{erreur}}</span>
            </div>       
        </div>
     </div>
</template>

<style>

</style>