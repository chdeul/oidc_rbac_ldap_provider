<script>
    app.component("ui-modale", {
        template: "#ui-modale",
        props : [ "titre" ],
        emits: [ "valider" ],
        methods: {
            valider() {
                this.$emit( "valider" )
                // hide() à appeler depuis l'appelant
            },
            show() {
                bootstrap.Modal.getOrCreateInstance(this.$refs.modale).show()
            },
            hide() {
                document.activeElement?.blur()
                bootstrap.Modal.getInstance(this.$refs.modale).hide()
            }
        }
   });
//# sourceURL=http://localhost:8088/vues/public/modale.vue
</script>

<template id="ui-modale">

    <div class="modal" data-bs-backdrop="static" ref="modale" tabindex="-1" >
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-header">
                    <h5 class="modal-title">{{titre}}</h5>
                    <button type="button" class="btn-close" @click="hide()" ></button>
                </div>
                <div class="modal-body">
                    <slot>Modal body</slot>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-secondary"  @click="hide()" >Annuler</button>
                    <button type="button" class="btn btn-primary" @click="valider()">Valider</button>
                </div>
            </div>
        </div>
    </div>

</template> 

<style>

</style>