<template>
  <div class="wrapper">
    <div class="card">
      <table class="table is-striped">
        <thead>
          <tr>
            <th v-for="(header, index) in headers" :key="index">{{ header }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(row, index) in itemrows" :key="index">
            <td class="table-item-title" v-html="row.title"></td>
            <td class="table-item-date">{{ row.res_time_limit }}</td>
            <td class="table-item-date">{{ row.modified_at }}</td>
            <td v-html="row.resultsLinkHtml"></td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script>
import axios from '@/bin/axios'
import Table from '@/components/Utils/Table.vue'

export default {
  name: 'ExplorerLayout',
  components: {
    'customtable': Table
  },
  async created () {
    const resp = await axios.get('/questionnaires')
    this.questionnaires = resp.data
  },
  props: {
    traqId: {
      type: String,
      required: true
    }
  },
  data () {
    return {
      questionnaires: [],
      headers: ['Title', 'Time Limit', 'Modified At', 'Results']
    }
  },
  computed: {
    itemrows () {
      let row = {}
      let itemrows = []
      for (let i = 0; i < this.questionnaires.length; i++) {
        row.title = this.getTitleHtml(i)
        row.res_time_limit = this.getDateStr(this.questionnaires[i].res_time_limit)
        row.modified_at = this.getDateStr(this.questionnaires[i].modified_at)
        row.resultsLinkHtml = this.getResultsLinkHtml(this.questionnaires[i].questionnaireID) // 結果を見る権限があるかどうかでボタンの色を変えたりしたい
        itemrows.push(row)
      }
      return itemrows
    }
  },
  methods: {
    getTitleHtml (i) {
      return '<a href="/questionnaires/' + this.questionnaires[i].questionnaireID + '">' + this.questionnaires[i].title + '</a>'
    },
    getResultsLinkHtml (id) {
      return '<a href="/resuslts/' + id + '">' + '<span class="ti-new-window"></span><br>Open' + '</a>'
    },
    getDateStr (str) {
      return str === 'NULL' ? '-' : new Date(str).toLocaleDateString()
    }
  }
}
</script>

<style scoped>
td {
  vertical-align: middle;
  font-size: 0.9em;
}
.table-item-title {
  min-width: 10em;
  font-size: 1em;
}
.table-item-date {
  min-width: 8em;
  text-align: center;
}
</style>
