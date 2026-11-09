
### Fonctions IHM :

* Browser le référentiel utilisateur (__fait__)

* Browser et éditer le ldap  (__fait__)

* Edition dans Feuille Calc (__en cours__)

* Création, Modifs des rbac projet

* Auth et droits : admin, edition, lecture, public

* Présentation Notif, Confirmation choix sur écarts

* Admin les batchs ( supressions, notifs ),

* Config smtp, mapping roles, outils

* Webhooks sur modif rbac projet, création

* Liste des projets, organisation 3 niveaux

* Page de suvi d'un projet

* Page de liens vers outils projet

* Création des projets, Mise à jour
  - nom court, long, descript, id, 3 chefs de projets
  - risques, budget, statut, avancement, imputations
  - lots, versions, livrables, doc essentielles

* Page d'acceuil, résumé, mémorisé par utilisateur

* Liste des projets par utilisateur, et droits

* Suppression utilisateur, et Projet via ihm

### Fonctions Back, et Batch

* Lecture du référentiel 
  - fichier csv         (__fait__)
  - any ldap            (__fait__)
  - any base sql        (__fait__)
  - any rest api        (__fait__)

* Traitement rbac projet
  - création, update    (__fait__)
  - suppression, batchs
  - notifications écarts

* Connecteur des outils 
  - ldap groupes        (__fait__)
  - mantis api rest     (__fait__)
  - gitlab api          (__fait__)
  - alfresco
  - nextcloud
  - nexus oss
  - sonarqube
  - redmine
  - any rest api        (__en cours__)
  - any sql base        (__en cours__)
  - etc.

* Réception des rbac projets via REST

* Envoi : via Gitlab ci, alfresco notif, oOo, ...

* Stockage Bdd Sqlite,postgres
  - rbac projets
  - projets, pages
  - config ihm

* Création de projet 
  - dans les connecteurs des outils
  - via api rest 
  - désactivation, archivage, suppression

### Fonctions OpenIDConnect

* Serveur OIDC / LDAP (__fait__)

* Customs thèmes pages login/logout

* Client OIDC  (__fait__)

* Auth native / LDAP (__fait__)

* Logout centralisé (__fait__)

* Password renouvellement (__fait__)

### Reverse proxy OpenIDConnect

* Reverse proxy (__fait__)

* Filtrage ip, user (__fait__)

* Auth client oidc / header (__fait__)

* Automatismes javascript  (__fait__)

* Reesources statiques (__fait__)

* Dockerfile , Env var  

### Tests 

* Unitaires    (__en cours__)

* Intégration   (__en cours__)

* En Containeurisation 