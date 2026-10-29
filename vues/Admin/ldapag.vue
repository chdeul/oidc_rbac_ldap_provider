<script>
    app.component('ui-ldapag', {

        template : "#ui-ldapag",
        mixins : [ authent ],

        data () {
            return {

                // Recherche
                pattern : "*",
                isGroup : false,

                // Résultat recherche
                users  : [],    // Table des utilisateurs trouvés
                groups : [],    // Table des noms de groupes trouvés
                members: [],    // Identifiants des membres
                erreur : "",    // Message d'erreur de recherche

                // Modale d'édition utilisateur
                modal  : null,
                user   : {},     // Utilisateur édité
                optPwd : "none", // Option passwd none|check|change

                // Modale édition de Groupe
                modalGrp : null, 
                group    : "",   // Nom du groupe édité
                mlist    : [],    // Identifiants des membres
                slist    : []    // Liste de choix d'identifiants
            }
        },

        mounted : function() {
            this.modal = new bootstrap.Modal( this.$refs.editBox )
            this.modalGrp = new bootstrap.Modal( this.$refs.editGrpBox )
        },

        methods : { 

            // SWitch entre users et groups
            toggleGroup() {
                this.pattern = "*"
                this.erreur = ""
            },

            // Recherche users ou groups
            search : function() {

                // Reset msg d'erreur
                this.erreur = ""
                this.user = {}
                this.group = ""
                this.members = []

                // Que chercher
                let query = { }
                if ( this.isGroup ) {
                    query.group = this.pattern
                } else  {
                    query.user = this.pattern 
                }

                // Appel serveur search ldap users
                myAxios.get("/ldap", { 
                    params: query 
                }).then( (rep) => {
                    if ( this.isGroup ) 
                        this.groups = rep.data
                    else 
                        this.users = rep.data
                }).catch( (err) => {
                    this.erreur = "Erreur appel REST : " + err
                })
            },

            // Recherche users identifiant pour edition d'un groupe
            searchUsers : function( pattern ) {
                this.erreur = ""
                this.slist = ""
                // Appel serveur search ldap users
                myAxios.get("/ldap", { 
                    params: { user:pattern }
                }).then( (rep) => {
                    this.slist = rep.data.map( u => u.Identifiant )
                }).catch( (err) => {
                    this.slist = [ ]
                    this.erreur = "Erreur appel REST : " + err
                })
            },

            // Sélection d'un groupe : affiche ses membres
            groupe : function(group) {
                this.erreur = ""
                this.group = group
                this.members = []
                // Appel serveur search ldap users
                myAxios.get("/ldap", { 
                    params: { membersOf:group.Identifiant } 
                }).then( (rep) => {
                    this.members = rep.data
                }).catch( (err) => {
                    this.erreur = "Erreur appel REST : " + err
                })
            },

            // Affichage liste séparateur virgule 
            lister : function(tab) {
                let ret = ""
                tab.forEach(elt => { ret = ret + ", " + elt });
                return ret.replace(", ","")
            },

            // Initialise le formulaire d'édition utilisateur
            init : function( edit ) {
                this.isGroup = false
                this.erreur = ""
                this.$refs.Ident.readOnly = edit
                this.$refs.Uid.readOnly = edit
                this.$refs.Pass.readOnly = true
                this.optPwd = "none"               
                this.$refs.PassBc.style.backgroundColor = ""                
                this.$refs.VisiBc.style.backgroundColor = ""                
                this.$refs.Pass.placeholder = "Crypté - SHA256"
                this.user.MotDePasse == ""
                this.modal.show()            
            },

            // Initialise le formulaire d'édition d'un groupe
            initGroup : function( edit ) {
                this.isGroup = true
                this.erreur = ""
                this.$refs.Group.readOnly = edit
                this.slist = []

                this.modalGrp.show()            
            },

            // Nouvel utilisateur ( avec ident )
            addUser : function() {
                this.user = { Uid:"", Identifiant:"", NomComplet:"", MotDePasse:"", AdresseMail:"", Roles:"" }
                this.init( false )
            },

            // Edition utilisateur courant ( sauf ident )
            editUser : function(user) {
                this.user = JSON.parse(JSON.stringify(user))
                this.init( true )
            },
  
            addGroup : function() {
                this.group = ""
                this.mlist = [ "admin" ]
                this.initGroup( false )     
            },

            editGroup : function(group) {
                this.erreur = ""
                this.group = group
                this.members = []
                // Appel serveur search ldap users
                myAxios.get("/ldap", { 
                    params: { membersOf:group.Identifiant } 
                }).then( (rep) => {
                    this.members = rep.data
                    this.mlist = JSON.parse(JSON.stringify(this.members))
                    this.initGroup( true )  
                }).catch( (err) => {
                    this.erreur = "Erreur appel REST : " + err
                })
            },

            // Validation utilisateur edité ou ajouté
            valider : function() {

                this.erreur = ""
    
                // Vérifications
                if ( !this.user.Uid.match(/^[0-9]{1,12}$/) ) {
                    this.erreur = "Id incorrect : 0-9 et taille 1 à 12 caractères"
                    return
                }
                
                // Vérifications
                if ( !this.user.Identifiant.match(/^[a-z0-9_]{2,63}$/) ) {
                    this.erreur = "Identifiant incorrect : a-z 0-9 _ et taille 2 à 63 caractères"
                    return
                }

                // this.MotDePasse si saisie  
                if ( this.optPwd !== "none" ) {
                    if ( !this.user.MotDePasse.match(/^[a-zA-Z0-9_!@#$%^&*()=+.?-]{6,64}$/) ) {
                        this.erreur = "Mot de passe incorrect : 6 à 64 caractères, a-zA-Z0-9_!@#$%^&*()=+.?-"
                        return
                    }

                }  
                if ( this.user.NomComplet.length < 3 ) {
                    this.erreur = "Nom complet vide ou trop court ( 3 caractères ou plus )"
                    return
                }
                if ( this.user.AdresseMail.length < 3 || !this.user.AdresseMail.includes("@") ) {
                    this.erreur = "Adresse E-mail incorrecte"
                    return
                }                
                if ( this.user.Roles.length > 127 ) {
                    this.erreur = "Roles, titres trop long"
                    return
                }                
  
                this.modal.hide()   

                let query = { }
                if ( !this.$refs.Ident.readOnly ) {
                    // Nouvel utilisateur
                    query.user = true
                    myAxios.post("/ldap", this.user, { 
                        params: query 
                    }).then( (rep) => {
                        this.search()  // maj table 
                    }).catch( (err) => {
                        this.erreur = "Erreur appel REST : " + err
                    })                 
                    return 
                } 
                
                // Modification d'un utilisateur sauf Identifiant et mot de passe
                if ( this.optPwd === "change" ) query.change = true
                else if ( this.optPwd === "check" ) query.check = true
                else delete this.user.MotDePasse
                myAxios.put("/ldap", this.user, { 
                    params: query 
                }).then( (rep) => { 
                    this.search()  // maj table  
                }).catch( (err) => {
                    this.erreur = "Erreur appel REST : " + err
                })
            },

            // Suppression utilisateur edité
            supprimer : function() {

                this.modal.hide()   

                this.erreur = ""
                myAxios.delete("/ldap", { 
                    params: { user: this.user.Identifiant } 
                }).then( (rep) => {  
                    this.search()  // maj table 
                }).catch( (err) => {
                    this.erreur = "Erreur appel REST : " + err
                })
            },

            validerGroupe : function() {

                this.erreur = ""
    
                 // Vérifications
                if ( !this.group.Identifiant.match(/^[a-zA-Z0-9_]{2,63}$/) ) {
                    this.erreur = "Nom de groupe incorrect : a-z A-Z 0-9 _ et taille 2 à 63 caractères"
                    return
                }
                // Ferme d'abord la boite
                this.modalGrp.hide() 

                // Data for api call 
                const data = { Group: this.group.Identifiant, Members : this.mlist }
                const query = { group : true  }    

                if ( !this.$refs.Group.readOnly ) {
                    // Nouveau Groupe
                    myAxios.post("/ldap", data, { 
                        params: query 
                    }).then( (rep) => { 
                        this.search()  // maj table     
                    }).catch( (err) => {
                        this.erreur = "Erreur appel REST : " + err
                    })
                   
                } else {
                    // Modification d'un groupe sauf nom
                    myAxios.put("/ldap", data, { 
                        params: query 
                    }).then( (rep) => {  
                        this.search()  // maj table     
                    }).catch( (err) => {
                        this.erreur = "Erreur appel REST : " + err
                    })
                }
            },
            
            // Suppression groupe édité
            supprimerGroupe : function() {

                // Ferme d'abord la boite
                this.modalGrp.hide() 

                this.erreur = ""
                myAxios.delete("/ldap", { 
                   params: { group: this.group } 
                }).then( (rep) => {
                   this.search()  // maj table   
                }).catch( (err) => {
                    this.erreur = "Erreur appel REST : " + err
                })
            },

            // Change d'Option d'actions sur password
            setCheckPwd : function(event) {
                let elt = this.$refs.Pass
                let btn = event.currentTarget
                this.user.MotDePasse =""  
                // De none à check
                if ( this.optPwd === "none" && this.$refs.Ident.readOnly ) {
                    elt.readOnly = false
                    elt.placeholder = "Mot de passe actuel"
                    btn.style.backgroundColor = "lightgreen" 
                    this.optPwd = "check"
                // De Check à Change
                } else if ( this.optPwd === "check" || this.optPwd === "none") {
                    elt.readOnly = false
                    elt.placeholder = "Nouveau mot de passe"
                    btn.style.backgroundColor = "salmon"
                   this.optPwd = "change"
                // De Change à none
                } else {
                    elt.readOnly = true
                    elt.placeholder = "Crypté - SHA256"
                    btn.style.backgroundColor = ""                
                    this.optPwd = "none"
                 }
            },

            // Montre ou cache le password
            showPwd : function (event) {
                let elt = this.$refs.Pass
                let btn = event.currentTarget
                if ( elt.type === "text" ) {
                    elt.type = "password"
                    btn.style.backgroundColor = ""
                } else {
                    elt.type = "text"
                    btn.style.backgroundColor = "lightblue"
                }
             }
        } 
            
    })
//# sourceURL=http://localhost:8080/vues/Admin/ldapag.vue
</script>

<template id="ui-ldapag">
    <div class="col-md-10 offset-md-1 marge-haut" >

        <div class="col-md-8">
            <h3>Parcourir le LDAP</h3>
            <div>

                <div class="input-group">
                    <label class="input-group-text" for="grp-id">Groupe :</label>
                    <span class="input-group-text"><input type="checkbox" v-model="isGroup" id="grp-id" @change="toggleGroup"></span>
                    <label class="input-group-text" for="nom-id"> &nbsp; &nbsp; &nbsp; Nom :</label>
                    <input type="text" class="form-control" placeholder="pattern avec *" v-model="pattern" id="nom-id">
                    <button class="btn btn-primary" @click="search">Chercher</button>
                </div>

            </div>
        </div>
        <br/>
            <div class="m-3 text-danger">{{ erreur  }}</div>          
        <br/>
        <div v-if="isGroup"> 
             <h5>Groupes LDAP
                <button class="m-2 btn btn-sm btn-dafault" @click="addGroup" title="Nouveau groupe">
                    <span class="fa fa-plus"></span>
                </button> 
             </h5>
             <table class="table table-striped">
                <thead>
                    <tr>
                        <th class="col-icon"></th>
                        <th class="w-25">Common Name</th>
                        <th>Membres</th>
                   </tr>
                </thead>
                <tbody>
                    <tr v-for="g in groups" :key="g" @click="groupe(g)">
                        <td><span class="fa fa-edit" @click="editGroup(g)"></span></td>
                        <td>{{g.NomComplet}}</td>
                        <td>{{ group == g ? lister(members) : '...' }}</td>
                    </tr>
                </tbody>
            </table>
        </div>
        <div v-else>
             <h5>
                Utilisateurs LDAP
                <button class="m-2 btn btn-sm btn-dafault" @click="addUser" title="Nouvel utilisateur">
                    <span class="fa fa-plus"></span>
                </button> 
             </h5>
             <table class="table table-striped">
                <thead>
                    <tr>
                        <th class="col-icon"></th>
                        <th>Id</th>
                        <th>Short Name</th>
                        <th>Common Name</th>
                        <th>Mail</th>
                        <th>Profile</th>
                   </tr>
                </thead>
                <tbody>
                    <tr v-for="u in users" :key="u.sn">
                        <td><span class="fa fa-edit" @click="editUser(u)"></span></td>
                        <td>{{u.Uid}}</td>
                        <td>{{u.Identifiant}}</td>
                        <td>{{u.NomComplet}}</td>
                        <td>{{u.AdresseMail}}</td>
                        <td>{{u.Roles}}</td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>

    <div class="modal" ref="editBox" tabindex="-1">
        <div class="modal-dialog" style="max-width: 600px;">
            <div class="modal-content">
                <div class="modal-header">
                    <h5 class="modal-title">Edition d'un Utilisateur</h5>
                    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Fermer"></button>
                </div>
                <form onsubmit="return false" class="modal-body">
                    <div class="input-group">
                        <span class="input-group-text w-25">Id :</span>
                        <input type="text" class="form-control" v-model="user.Uid" ref="Uid" readonly 
                            autocomplete="off">
                    </div>
                    <div class="input-group">
                        <span class="input-group-text w-25">Identifiant :</span>
                        <input type="text" class="form-control" v-model="user.Identifiant" ref="Ident" readonly 
                            autocomplete="off">
                    </div>
                    <div class="input-group">
                        <span class="input-group-text w-25">Mot de passe :</span>
                        <input type="password" class="form-control" v-model="user.MotDePasse" ref="Pass" readonly 
                            placeholder="Crypté - SHA256" autocomplete="new-password">   
                        <button type="button" class="btn btn-light" @click="setCheckPwd($event)" ref="PassBc">
                            <span class="fa fa-edit"></span>
                        </button>
                        <button type="button" class="btn btn-light" @click="showPwd($event)" ref="VisiBc">
                            <span class="fa fa-eye"></span>
                        </button>
                    </div>
                    <div class="input-group">
                        <span class="input-group-text w-25">Nom complet :</span>
                        <input type="text" class="form-control" v-model="user.NomComplet">
                    </div>
                    <div class="input-group">
                        <span class="input-group-text w-25">Adresse mail :</span>
                        <input type="email" class="form-control" v-model="user.AdresseMail">
                    </div>
                    <div class="input-group">
                        <span class="input-group-text w-25">Roles, titres :</span>
                        <input type="text" class="form-control" v-model="user.Roles">
                    </div>              
                    <div class="m-3 text-danger">{{ erreur  }}</div>          
                </form>
                <div class="modal-footer">
                    <button type="button" class="btn btn-warning" @click="supprimer">Supprimer</button>
                    <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Annuler</button>
                    <button type="button" class="btn btn-primary" @click="valider">Valider</button>
                </div>
            </div>
        </div>
    </div>

    <div class="modal" ref="editGrpBox" tabindex="-1">
        <div class="modal-dialog" style="max-width: 600px;">
            <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title">Edition d'un groupe</h5>
                <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Fermer"></button>
            </div>
            <div class="modal-body">
                 <div class="input-group">
                    <span class="input-group-text w-25">Nom du groupe :</span>
                    <input type="text" class="form-control" v-model="group.Identifiant" ref="Group" readonly 
                       autocomplete="off">
                </div>

                <ui-tags v-model:items="mlist" :slist="slist" @list-tags="searchUsers" label="Membres"/>
           
                <div class="m-3 text-danger">{{ erreur  }}</div>          
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-warning" @click="supprimerGroupe">Supprimer</button>
                <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Annuler</button>
                <button type="button" class="btn btn-primary" @click="validerGroupe">Valider</button>
            </div>
            </div>
        </div>
    </div>

</template>

<style>

</style>