<script>
    app.component('ui-projet', {

        template : "#ui-projet",

        data () {
            return {
                // Tableau csv
                headers: [],
                items: [],

                // Red Message 
                errorMsg : ""
            }
        },

        mounted() {
            this.load()
        },

        methods : { 

            load() {
                myAxios.get("projet?name=Mon_Projet").then( (rep) => {

                //this.headers = { text: "LAST", value: "lastAttended", sortable: true,width: 200},
                //this.items =

                }).catch( (err) => {
                    this.errorMsg = "Projet non chargé !"
                })
            }, 
            
            exec() {
                myAxios.post("projet", "").then( (rep) => {
                }).catch( (err) => {
                    this.errorMsg = "Projet non pris en compte !"
                })
            }
        }

     })
//# sourceURL=http://localhost:8088/vues/edition/projet.vue
</script>

<template id="ui-projet">

    <div class="col-md-12 marge-haut" >
        <h4>Table des droits du projet</h4>
        <easy-data-table id="table" :headers="headers" :items="items"></easy-data-table>
    </div>

    <div v-show="errorMsg" class="mb-3 red">
        <span>{{errorMsg}}</span>
    </div>

    <div class="col-md-12 marge-haut" >
        <div class="col-md-3">
           <button class="btn btn-warning" @click="exec">Prendre en compte</button>  
        </div>
     </div>  

</template>

<style>
</style>