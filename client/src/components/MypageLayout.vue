<template>
  <div class="wrapper">
    <div class="list card content">
      <header class="card-header">
        <div class="card-header-title subtitle">回答対象になっているアンケート</div>
      </header>
      <div class="card-content">
        <article class="post columns" v-for="(row, index) in itemrows" :key="index">
          <div class="column">
            <h4 v-html="row.title"></h4>
            <p>{{ row.description }}</p>
            <div class="media">
              <div class="media-left">{{ row.status }}</div>
              <div class="media-content has-text-weight-bold columns">
                <div class="content column">回答期限: {{ row.res_time_limit }}</div>
                <div class="content column">更新日: {{ row.modified_at }}</div>
              </div>
            </div>
          </div>
          <!-- <div class="media-right column is-narrow" v-html="row.resultsLinkHtml"></div> -->
        </article>
      </div>
    </div>
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
      let row = {}
      let rows = []
      for (var i = 0; i < this.questionnaires.length; i++) {
        row.title = this.getTitleHtml(i)
        row.description = this.questionnaires[i].description
        row.res_time_limit = this.questionnaires[i].res_time_limit
        row.status = this.hasResponded(i) ? '✔︎' : '-' // saved も返せるようにしたさ
        row.modified_at = this.questionnaires[i].modified_at
        row.resultsLinkHtml = this.getResultsLinkHtml(this.questionnaires[i].questionnaireID) // 結果を見る権限があるかどうかでボタンの色を変えたりしたい

        rows.push(row)
      }
      return rows
    }
  },
  methods: {
    hasResponded (index) {
      return this.questionnaires[index].responded_at != null
    },
    getTitleHtml (i) {
      return '<a href="/questionnaires/' + this.questionnaires[i].questionnaireID + '">' + this.questionnaires[i].title + '</a>'
    },
    getResultsLinkHtml (id) {
      return '<a href="/resuslts/' + id + '" class="button is-info">Results</a>'
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.card {
  /* width: fit-content; */
  margin: 1rem 1.5rem;
  overflow-x: auto;
  width: auto;
  max-width: 800px;
}
.card-header-title {
  color: #707880;
  font-weight: 400;
  padding: 1rem 1.5rem;
}
.card-content {
  padding: 1rem;
}
.content {
  margin-left: 1.5rem;
}
.content p {
  margin-bottom: 0.5em;
  word-break: break-all;
  line-height: 1.1em;
}
.content h4 {
  margin-bottom: 0.7em;
}
.content.column {
  padding: 0;
  margin-bottom: 0;
  max-width: fit-content;
  display: inline-block;
}
.columns.media-content {
  padding-top: 1em;
}
article.post {
  margin: 1rem;
  padding-bottom: 0;
  border-bottom: 1px solid #e6eaee;
}
.columns {
  padding-top: 0;
}
.column {
  padding-left: 0;
}
.media-right.column {
  margin: auto;
}
</style>
