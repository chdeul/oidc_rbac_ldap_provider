<script>
    app.component('ui-outils', {

        template : "#ui-outils",
        
        data () {
            return {

                // CRUD  Project, User, Role
                projname : "",
                username : "",
                rolename : "",

                // Table de strings
                items : [],

                // Error message
                erreur : ""
            }
        },

        methods : { 

            search: function() {

                // Que chercher 
                let query = { }
                query.project = this.projname 
                query.user = this.username 
                query.role = this.rolename
    
                // Appel serveur 
                myAxios.get("/outil/gitlab", { 
                    params: query 
                }).then( (rep) => {

                }).catch( (err) => {

                    this.erreur = "Projet non pris en compte !"
                })
            },
            select: function(item) {

            },
            edit: function(item) {

            }
        }

     })
//# sourceURL=http://localhost:8080/vues/Admin/outils.vue
</script>

<template id="ui-outils">
    <div class="col-md-12 marge-haut" >
 
        <div class="input-group">
            <label for="id-proj">Projet:</label>
            <input type="text" id="id-proj" v-model="projname">
        </div>
        <div class="input-group">
            <label for="id-user">Projet:</label>
            <input type="text" id="id-user" v-model="username">
        </div>
        <div class="input-group">
            <label for="id-role">Projet:</label>
            <input type="text" id="id-role" v-model="rolename">
        </div>
        <button class="btn btn-secondary" @click="search">Chercher</button>

        <table class="table table-striped">
            <thead>
                <tr>
                    <th class="col-icon"></th>
                    <th class="w-25">Name</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="item in items" :key="item" @click="select(item)">
                    <td><span class="fa fa-edit" @click="edit(item)"></span></td>
                    <td>{{item}}</td>
                </tr>
            </tbody>
        </table>

    </div>
</template>

<style>

</style>



