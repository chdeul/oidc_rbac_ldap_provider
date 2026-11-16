<script>
    app.component('ui-projet', {

        template : "#ui-projet",

        data () {
            return {

                // Projects names list
                prjNames : [], 

                // Projet, raw data
                name: "",
                projet : [],
                tools: [],

                // Projet, easydatatable
                headers: [],
                items: [],
                itemEdit: {},
            }
        },
        created () {
            // Load list of projects, that current user is manager
            this.prjNames = [ "Mon_Projet", "Autre_Projet" ]
        
            // Get row data from backend
            myAxios.get("projet").then( (rep) => {

                // Check and store raw data : not empty
                this.prjNames = rep.data
                if ( this.prjNames.length == 0 ) {
                    message( "Aucun projet que vous pouvez gérer !" )
                    return
                }

            }).catch( (err) => {
                message( "Erreur, liste des projets non chargée !")
            })
        },

        methods : { 

            init() {
                // Pas de projet chargé
                this.projet = []
                this.items = []
                this.headers = []
                this.tools = []
                this.itemEdit = {}
            },

            loadProject(event) {
                // Clear
                this.init()
                // Get slected projet
                this.name = event.target.value
                if ( !this.name ) return
                // Get row data from backend
                myAxios.get("projet?name="+this.name).then( (rep) => {

                    // Check and store raw data : not empty
                    this.projet = rep.data
                    if ( this.projet.length == 0 ) {
                          message( "Projet inexistant !" )
                          return
                    }
                    // Generate data for easydatatable
                    this.generateTable(this.projet)

                }).catch( (err) => {
                    message( "Erreur, projet non chargé !")
                })
            }, 

            generateTable( data ) {
                   // Header (1st element) of table
                    const hd = data[0]
                    const rows = data.slice(1)
                    this.headers = [
                         { text: "", value: "edit", width: 20 }, //icon
                         { text: "", value: "del", width: 20 },  //icon
                         { text: hd.Type, value: "Type", width: 120 },
                         { text: hd.Name, value: "Name", sortable: true, width: 100 },
                    ]
                    this.tools = hd.Tools
                    hd.Tools.forEach( (t,j) => {
                        this.headers.push( { text: t, value: t, sortable: true, width: 100 } )
                    });
                    // Space 
                    this.headers.push( { text: "", value: "Dummy" } )
                    
                    // Rows User/Rights of table
                    this.items = [ ]
                    rows.forEach((ur) => {
                        item = { Type: ur.Type, Name: ur.Name }
                        ur.Tools.forEach( (t,i) => { item[hd.Tools[i]] = t })
                        this.items.push( item )
                    })
            },

            editUserRights(item) {
                if ( item )  { // Modification
                    if ( !item.Name ) return
                    this.itemEdit = JSON.parse(JSON.stringify(item)) 
                } else {  // Nouveau
                    this.itemEdit = { Type: "User", Name: "" } 
                }
                this.$refs.urModale.show()
            },

            checkUserRights(item) {

                // Type : 
                if ( item.Type !== "User" && item.Type !== "Group" )
                    return message("Type invalide")

                // Username :
                if ( !/^\w{3,32}$/.test(item.Name) )
                    return message("Name invalide")

                let ret = true
                this.tools.forEach( (t) => {
                    right = item[t]
                    // Role in tools
                    if ( !/^\w{0,32}$/.test(right) ) {
                        message("Droit invalide : " + t+"/"+right)  
                        ret = false 
                    }                 
                })
                return ret
            },

            applyUserRights() {

                // Controls input contents
                if ( !this.checkUserRights(this.itemEdit) ) return

                // Old or new row 
                let item = this.projet.find( (ur) => { return ur.Name === this.itemEdit.Name } )
                if ( !item ) item = { Tools : [] }

                // New fields values
                item.Type = this.itemEdit.Type
                item.Name = this.itemEdit.Name
                this.tools.forEach ( (t,i) => {
                   item.Tools[i] = this.itemEdit[t]
                })
                // Hide modale and refresh table
                message( "Droits utilisateurs modifiés" )
                this.$refs.urModale.hide()
                this.itemEdit = { }
                this.generateTable( this.projet )
            },

            delUserRights(item) {

                if ( !item  || !item.Name ) return

                confirme( "Voulez-vous supprimer les droits de [ " + item.Name + " ]", () => {
                   this.projet = this.projet.filter( ur => ur.Name !== item.Name )
                   this.items = this.items.filter( ur => ur.Name !== item.Name )
                   message( "Droits utilisateur supprimés" )
                })
            },

            execUserRights() {

                if ( this.projet.length == 0 ) return
                
                myAxios.post("projet", this.projet).then( (rep) => {
                   message( "Droits projet appliqués et enregistrés" )
                }).catch( (err) => {
                   message( "Erreur, droits projet non appliqués !" )
                })
            }
        }

     })
//# sourceURL=http://localhost:8088/vues/edition/projet.vue
</script>

<template id="ui-projet">

    <div class="col-md-12 marge-haut" >
        <div>
            <select class="" style="float:right" v-model="name" @change="loadProject">
                <option value="" disabled>Select...</option>
                <option v-for="prj in prjNames" :key="prj" :value="prj">{{prj}}</option>
            </select>
            <h4>Table des droits du projet</h4>
        </div>
        <easy-data-table id="table" :headers="headers" :items="items">

            <template #header-edit="header">
                <div class="customize-header">
                    <span class="fa fa-plus" @click="editUserRights(null)"></span>
                </div>
            </template>

            <template #item-edit="item" >
                 <span class="fa fa-edit" @click="editUserRights(item)"></span>
            </template>

            <template #item-del="item" >
                 <span class="fa fa-trash" @click="delUserRights(item)"></span>
            </template>

        </easy-data-table>

        <div class="col-md-12 marge-haut" >
            
            <button type="button" class="marge btn btn-primary" @click="execUserRights()">Prendre en compte</button>  
            <button type="button" class="marge btn btn-secondary" @click="load()">Annuler</button>  
            
        </div>

    </div>
    
    <!-- modale d'edition de ligne -->
    <ui-modale id="ur-modale" ref="urModale" titre="Utilisateur / Droits" @valider="applyUserRights()">

        <div class="mb-3 row">
            <label for="type" class="col-sm-2 col-form-label">Type</label>
            <div class="col-sm-10">
                <input type="text" class="form-control-plaintext" id="type"  v-model="itemEdit.Type">
            </div>
        </div>
        <div class="mb-3 row">
            <label for="name" class="col-sm-2 col-form-label">Name</label>
            <div class="col-sm-10">
                <input type="text" class="form-control-plaintext" id="name"  v-model="itemEdit.Name">
            </div>
        </div>
        <div class="mb-3 row" v-for="tool,i in tools" key="tool">
            <label :for="tool" class="col-sm-2 col-form-label">{{ tool }}</label>
            <div class="col-sm-10">
                <input type="text" class="form-control-plaintext" :id="tool"  v-model="itemEdit[tool]" >
            </div>
        </div>   
    </ui-modale>

</template>

<style>
</style>