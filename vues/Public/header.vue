<script>
    app.component("ui-header", {
        template: "#ui-header"
   });
//# sourceURL=http://localhost.vue/vues/Public/header.vue
</script>

<template id="ui-header">

  <!-- header sombre -->
  <header class="navbar navbar-expand-sm bg-dark navbar-dark d-print-none">
    <div class="container-fluid">

      <!-- logo , home page -->
      <a class="navbar-brand" href="#"><img src="resources/img/glue.png" height="40"/>&nbsp; Agréments </a>

      <!-- compact menu sur petit ecran -->
      <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#collapsibleNavbar">
        <span class="navbar-toggler-icon"></span>
      </button>

      <!-- menu -->
      <div class="collapse navbar-collapse" id="collapsibleNavbar">

        <!-- menu authentifié , partie à gauche --> 
        <ul class="navbar-nav">

           <li class="nav-item dropdown" v-show="$root.isConsult()">
              <span class="nav-link dropdown-toggle" role="button" data-bs-toggle="dropdown">Administration</span>
              <ul class="dropdown-menu  dropdown-menu-end">
                <!-- Les fonctions ldap -->
                <li><router-link to="/vue/admin/ldapag" class="dropdown-item">Ldap agréments</router-link></li>
                <li><router-link to="/vue/admin/refbrowse" class="dropdown-item">Annuaire</router-link></li>
                <li><router-link to="/vue/admin/outils" class="dropdown-item">Outils</router-link></li>
                 <li><router-link to="/vue/admin/projet" class="dropdown-item">Projet</router-link></li>
              </ul>
           </li>

           <li class="nav-item dropdown" v-show="$root.isEdit()" >
              <span class="nav-link dropdown-toggle" role="button" data-bs-toggle="dropdown">Edition</span>
              <ul class="dropdown-menu  dropdown-menu-end">
                <li><router-link to="/vue/home" class="dropdown-item">Accueil</router-link></li>
              </ul>
           </li>

           <li class="nav-item dropdown" v-show="$root.isConsult()" >
              <span class="nav-link dropdown-toggle" role="button" data-bs-toggle="dropdown">Consultation</span>
              <ul class="dropdown-menu  dropdown-menu-end">
                <li><router-link to="/vue/home" class="dropdown-item">Accueil</router-link></li>
              </ul>
           </li>

           <!-- menu public, à droite -->
           <li class="nav-item dropdown" >
              <span class="nav-link dropdown-toggle" role="button" data-bs-toggle="dropdown">Publique</span>
              <ul class="dropdown-menu  dropdown-menu-end">
                <li><router-link to="/vue/home" class="dropdown-item">Accueil</router-link></li>
              </ul>
           </li>

        </ul>

        <!-- menu public , partie à droite -->
        <ul class="navbar-nav ms-auto">

           <!-- Authentification -->
           <li v-if="$root.logged" class="nav-item dropdown">
              
              <span class="nav-link dropdown-toggle" role="button" data-bs-toggle="dropdown"><i class="fa fa-user"></i> {{$root.logged}}</span>
              <ul class="dropdown-menu  dropdown-menu-end">
                <li><router-link to="/vue/home" class="dropdown-item"><i class="fa fa-cog"></i> Paramètres</router-link></li>
                <li><router-link to="/vue/logout" class="dropdown-item"><i class="fa fa-power-off"></i> Déconnexion</router-link></li>
              </ul>
           </li>
           <li v-else class="nav-item">
              <router-link to="/vue/login" class="nav-link"><i class="fa fa-user"></i> Connexion</router-link>
           </li>

        </ul>

      </div>

    </div>
  </header>
</template> 

<style scoped>
.navbar-brand {
  padding-top : 0px;
  padding-bottom : 0px;
}
</style>