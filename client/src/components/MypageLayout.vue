<template>
  <div class="wrapper">
    <customtable :headers="headers" :itemrows="itemrows"></customtable>
    <!-- -->
  </div>
</template>

<script>
import axios from '@/bin/axios'
import Table from '@/components/Utils/Table.vue'

export default {
  name: 'MypageLayout',
  components: {
    'customtable': Table
  },
  async created () {
    const resp = await axios.get('/users/me/targeted')
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
      headers: ['Title', 'Time Limit', 'Response', 'Modified At', 'Results', 'Details']
    }
  },
  computed: {
    itemrows () {
      let row = []
      let rows = []
      for (var i = 0; i < this.questionnaires.length; i++) {
        row[0] = this.questionnaires[i].title
        row[1] = this.questionnaires[i].res_time_limit
        row[2] = this.hasResponded(i) ? 'sent' : '-' // saved も返せるようにしたさ
        row[3] = this.questionnaires[i].modified_at
        row[4] = this.getResultsLinkHtml(this.questionnaires[i].questionnaireID)
        row[5] = this.getDetailsLinkHtml(this.questionnaires[i].questionnaireID)
        rows.push(row)
      }
      return rows
    }
  },
  methods: {
    hasResponded (index) {
      // 回答送信済み : ✔︎, 未送 : !
      return this.questionnaires[index].responded_at != null
    },
    getResultsLinkHtml (id) {
      return '<a href="/results/' + id + '">link</a>'
    },
    getDetailsLinkHtml (id) {
      return '<a href="/questionnaires/' + id + '">link</a>'
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.box {
  width: fit-content;
  margin: 1rem 2rem;
}
</style>
