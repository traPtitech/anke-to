<template>
  <div class="wrapper">
    <button class="button download-csv" v-on:click="downloadCSV">CSV形式でダウンロード</button>
    <div class="card">
      <table class="table is-striped">
        <thead>
          <tr>
            <th
              v-for="(header, index) in headers.concat(questions)"
              :key="index"
              @click="sort(index + 1)"
              :class="{ active: sorted == index + 1 || sorted == -1 - index }"
            >
              {{ header }}
              <span class="arrow" :class="sorted != index + 1 ? 'asc' : 'dsc'"></span>
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(result, index) in results" :key="index">
            <td class="table-item-traqid">{{ result.traqID }}</td>
            <td class="table-item-time">{{ getDateStr(result.submitted_at) }}</td>
            <td
              v-for="response in result.response_body"
              :key="response.responseID"
            >{{ getResponse(response) }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script>
/* eslint-disable */

import axios from '@/bin/axios'
import common from '@/bin/common'

export default {
  name: 'Spreadsheet',
  components: {},
  props: {
    results: {
      type: Array,
      required: true
    },
    questions: {
      type: Array,
      required: true
    }
  },
  data () {
    return {
      headers: [ 'traQID', '回答日時' ],
      sorted: ''
    }
  },
  methods: {
    getDateStr (str) {
      return common.getDateStr(str)
    },
    getResponse (body) {
      switch (body.question_type) {
        case 'MultipleChoice':
        case 'Checkbox':
        case 'Dropdown':
          let ret = ''
          body.option_response.forEach(response => {
            if (ret != '') {
              ret += ', '
            }
            ret += response
          })
          return ret
        default:
          return body.response
      }
    },
    downloadCSV () {
      let csv = '\ufeff'
      this.headers.concat(this.questions).forEach(header => {
        if (csv != '\ufeff') {
          csv += ','
        }
        csv += '"' + header + '"'
      })
      csv += '\n'
      this.results.forEach(result => {
        csv += result.traqID + ',' + this.getDateStr(result.submitted_at)
        result.response_body.forEach(response => {
          csv += ',' + '"' + this.getResponse(response) + '"'
        })
        csv += '\n'
      })
      const blob = new Blob([ csv ], { type: 'text/csv' })
      let link = document.createElement('a')
      link.href = window.URL.createObjectURL(blob)
      link.download = 'Result.csv'
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
    },
    sort (index) {
      let param = ''
      if (this.sorted != index) {
        param += '-'
        this.sorted = index
      } else {
        this.sorted = -index
      }
      switch (index) {
        case 1:
          param += 'traqid'
          break
        case 2:
          param += 'submitted_at'
          break
        default:
          param += index - 2
      }
      axios
        .get('/results/' + this.questionnaireId + '?sort=' + param)
        .then(res => {
          this.$emit('set-results', res.data)
        })
    }
  },
  computed: {
    questionnaireId () {
      return this.$route.params.id
    }
  },
  mounted () { }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style lang="scss" scoped>
th,
td {
  vertical-align: middle;
  font-size: 0.9em;
  min-width: 10em;
}

th.active {
  background-color: #fafafa;
}

.arrow {
  opacity: 1;
  color: black;
  display: inline-block;
  vertical-align: middle;
  margin: auto;
  opacity: 1;
}

.arrow.asc {
  border-left: 4px solid transparent;
  border-right: 4px solid transparent;
  border-bottom: 4px solid black;
}

.arrow.dsc {
  border-left: 4px solid transparent;
  border-right: 4px solid transparent;
  border-top: 4px solid black;
}

.download-csv {
  margin: 0 1.5rem;
}
</style>
