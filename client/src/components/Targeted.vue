<template>
  <div class="wrapper">
    <div class="list card content">
      <header class="card-header">
        <div class="card-header-title subtitle">回答対象になっているアンケート</div>
      </header>
      <div class="card-content">
        <article class="post" v-for="(row, index) in itemrows" :key="index">
          <div>
            <div class="questionnaire-title">
              <span :class="{'ti-check': row.status==='sent', 'ti-alert' : row.status==='unsent'}"></span>
              <span class="subtitle" v-html="row.title"></span>
            </div>
            <p>{{ row.description }}</p>
            <div class="media">
              <div class="media-content has-text-weight-bold columns">
                <div class="content column res-time-limit">回答期限: {{ row.res_time_limit }}</div>
                <div class="content column modified-at">更新日: {{ row.modified_at }}</div>
              </div>
            </div>
          </div>
          <!-- <div class="media-right column is-narrow" v-html="row.resultsLinkHtml"></div> -->
        </article>
      </div>
    </div>
    <!-- {{ questionnaires }} -->
  </div>
</template>

<script>
import axios from '@/bin/axios'
import {customDateStr, relativeDateStr} from '@/util/common'

export default {
  name: 'Mypage',
  components: {
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
      headers: [ 'Title', 'Time Limit', 'Response', 'Modified At', 'Results', 'Details' ]
    }
  },
  computed: {
    itemrows () {
      let rows = []
      for (var i = 0; i < this.questionnaires.length; i++) {
        let row = {}
        row.title = this.getTitleHtml(i)
        row.description = this.questionnaires[ i ].description
        row.res_time_limit = this.getDateStr(this.questionnaires[ i ].res_time_limit)
        row.status = this.getStatus(i)
        // row.status = this.hasResponded(i) ? '✔︎' : '-' // saved も返せるようにしたさ
        row.modified_at = this.getRelativeDateStr(this.questionnaires[ i ].modified_at)
        row.resultsLinkHtml = this.getResultsLinkHtml(this.questionnaires[ i ].questionnaireID) // 結果を見る権限があるかどうかでボタンの色を変えたりしたい

        rows.push(row)
      }
      return rows
    }
  },
  methods: {
    getDateStr (str) {
      return customDateStr(str)
    },
    getRelativeDateStr (str) {
      return relativeDateStr(str)
    },
    getStatus (i) {
      if (this.questionnaires[ i ].responded_at != null) {
        return 'sent'
      } else {
        return 'unsent'
      }
    },
    getTitleHtml (i) {
      return '<a href="/questionnaires/' + this.questionnaires[ i ].questionnaireID + '">' + this.questionnaires[ i ].title + '</a>'
    },
    getResultsLinkHtml (id) {
      return '<a href="/resuslts/' + id + '" class="button is-info">Results</a>'
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="scss" scoped>
.content {
  margin-left: 1.5rem;
  p {
    margin-bottom: 0.5em;
    word-break: break-all;
    line-height: 1.1em;
  }
  h4 {
    margin-bottom: 0.7em;
  }
}
.content.column {
  padding: 0;
  margin-bottom: 0;
  display: inline-block;
}
.content.column.res-time-limit {
  width: 15rem;
}
.content.column.modified-at {
  width: 10rem;
}
article.post {
  padding: 1rem;
  /* padding-bottom: 0; */
  border-bottom: 1px solid #e6eaee;
}
.columns {
  padding-top: 0;
  .media-content {
    padding-top: 1em;
  }
}
.column {
  padding-left: 0;
  .media-right {
    margin: auto;
  }
}
.questionnaire-title {
  padding-bottom: 1rem;
}
</style>
