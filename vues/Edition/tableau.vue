<script>
   app.component('ui-tableau', {

        template : "#ui-tableau",

        data() {
            return {
                head: [],
                data: [],
                editing: { row: null, col: null },
                editValue: ""
            }
        },

        created() {
            this.load()
        },

        methods: {
            load : function () {

                fetch("/agrem/resources/projet.csv")
                  .then(r => r.text())
                  .then(text => {
                    const rows = text.trim().split("\n").map(row => row.split(";"))
                    const newr = rows.map(row => row.map(s => ((s[0] === '"') ? s.slice(1, -1) : s) ))
                    this.head = newr[0]
                    this.data = newr.slice(1)
                  })
                    
            },

            isEditing: function (r, c) {
                return this.editing.row === r && this.editing.col === c
            },

            editCell: function (r, c) {
               this.editing = { row: r, col: c }
               this.editValue = this.data[r][c]

                this.$nextTick(() => {
                   const input = this.$refs.editor
                  if (input && input.focus) input.focus()
               })
            },

            saveEdit: function (r, c) {
               this.data[r][c] = this.editValue
               this.editing = { row: null, col: null }
            },

            addRow: function () {
               this.data.push( Array( this.head.length ).fill("") )
            }

        },

    })
//# sourceURL=http://localhost.vue/vues/Edition/tableau.vue
</script>

<template  id="ui-tableau">
  <div>

    <table class="spreadsheet">
      <thead>
        <tr>
          <th v-for="(col, cIndex) in head" :key="cIndex">
            {{ col }}
          </th>
        </tr>
      </thead>

      <tbody>
        <tr v-for="(row, rIndex) in data" :key="rIndex">
          <td v-for="(cell, cIndex) in row" :key="cIndex" @click="editCell(rIndex, cIndex)" >
            <input v-if="isEditing(rIndex, cIndex)" v-model="editValue" 
                @blur="saveEdit(rIndex, cIndex)" @keyup.enter="saveEdit(rIndex, cIndex)"
            ref="editor" />
            <span v-else>{{ cell }}</span>
          </td>
        </tr>
      </tbody>
    </table>

    <button @click="addRow">Ajouter une ligne</button>
  </div>
</template>


<style>
.spreadsheet {
  border-collapse: collapse;
  width: 100%;
}

.spreadsheet th,
.spreadsheet td {
  border: 1px solid #ccc;
  padding: 6px;
  text-align: center;
}

.spreadsheet td input {
  width: 100%;
  border: none;
  outline: none;
  padding: 4px;
}
</style>
