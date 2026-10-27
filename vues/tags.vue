<script>
    app.component('ui-tags', {

        template : "#ui-tags",

        emit : ['update:items', 'listTags'],

        data () {
            return {
                pattern: "",
                selected : "",
                erreur : ""
            }
        },

        props : {
            items: {
                type: Array,
                default: []
            },
            slist: {
                type: Array,
                default: []
            },
            pattern: {
                type: String,
                default: ''
            },
            label: {
                type: String,
                default: 'Tags'
            }
        },
           
        watch : {
            slist : function (newVal) {
                if ( newVal.length == 1 ) {
                    this.selected = newVal[0]
                    this.pattern = ""
                    this.addTag()
                } else {
                    this.selected = ""
                }
             } 
        },

        methods : {

            listTags : function () {
                if ( this.pattern.trim() ) {
                  this.$emit( "listTags", this.pattern.trim() )
                } else {
                  this.pattern = ""
                  this.selected = ""
                }
            },

            addTag : function () {
                if (this.selected && !this.items.includes(this.selected)) {
                    this.$emit('update:items', [...this.items, this.selected])
                }
            },

            removeTag : function (index) {
                const newItems = [...this.items]
                newItems.splice(index, 1)
                this.$emit('update:items', newItems)
            }
        }
    })
//# sourceURL=http://localhost.vue/vues/tags.vue
</script>

<template id="ui-tags">

    <div class="input-group" >

        <!-- Pattern pour recherche -->
        <label class="input-group-text" for="id-input-tag">{{ label }}</label>
        <input type="text" class="form-control" id="id-input-tag" v-model="pattern" @keydown.enter.prevent="listTags" placeholder="pattern avec *"/>

        <!-- Recherche pour alimenter le sélecteur -->
        <button type="button" class="btn btn-light" @click="listTags"><span class="fa fa-search"></span></button>

        <!-- Selecteur de tags parmi une liste -->
        <select class="form-control" ref="Select" @change="addTag" v-model="selected" >
            <option value="">{{ slist.length > 0 ? 'Sélectionnez ↓' : '← Cherchez' }}</option>
            <option v-for="t in slist" :key="t" :value="t">{{t}}</option>
        </select>

    </div>        

    <!-- Liste des tags choisis -->
    <div class="d-flex flex-wrap gap-2 mb-2" style="padding-top:8px">
        <span v-for="(tag, index) in items" :key="index" class="badge bg-secondary d-flex align-items-center">
            {{ tag }} <button type="button" class="btn-close btn-close-white ms-2" @click="removeTag(index)"></button>
        </span>
    </div>

</template>
