<template>
  <div class="wrapper">
    <div class="card">
      <div class="tabs">
        <ul>
          <li
            v-for="tab in tableFormTabs"
            :key="tab"
            class="tab"
            :class="{ 'is-active': tableForm === tab }"
          >
            <a @click="tableForm = tab">{{ tab }}</a>
          </li>
        </ul>
        <button
          v-if="canDownload"
          class="button download"
          @click="downloadTable"
        >
          <span class="ti-download"></span>
        </button>
      </div>
      <div class="scroll-view">
        <!-- table view -->
        <table v-show="tableForm === 'view'" class="table is-striped">
          <thead>
            <tr>
              <th
                v-for="(header, k) in tableHeaders"
                :key="k"
                :class="{
                  'active-header': isColumnActive(k),
                  hidden: isColumnHidden(k)
                }"
              >
                <span class="header-wrapper">
                  <span class="header-icon-left" @click="toggleShowColumn(k)">
                    <Icon
                      :name="isColumnHidden(k) ? 'eye-closed' : 'eye'"
                      color="var(--base-darkbrown)"
                      class="clickable"
                    ></Icon>
                  </span>
                  <span class="header-label">
                    {{ header }}
                  </span>
                  <span
                    class="header-icon-right clickable"
                    :class="sorted !== k + 1 ? 'ti-angle-up' : 'ti-angle-down'"
                    @click="sort(k + 1)"
                  ></span>
                </span>
              </th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(row, j) in results" :key="j">
              <td
                v-for="(item, k) in getTableRow(j)"
                :key="k"
                :class="{ hidden: isColumnHidden(k) }"
              >
                {{ item }}
              </td>
            </tr>
          </tbody>
        </table>

        <!-- markdown view -->
        <textarea
          v-show="tableForm === 'markdown'"
          class="textarea"
          :value="markdownTable"
          :rows="results.length + 3"
          readonly
        ></textarea>

        <!-- csv view -->
        <textarea
          v-show="tableForm === 'csv'"
          class="textarea"
          :value="csvTable"
          :rows="results.length + 2"
          readonly
        ></textarea>
      </div>
    </div>
  </div>
</template>

<script>
import Icon from '@/components/Icons/Icon'

export default {
  name: 'Spreadsheet',
  components: {
    Icon
  },
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
  data() {
    return {
      defaultColumns: [
        { name: 'traqId', label: 'traQID' },
        { name: 'submittedAt', label: '回答日時' }
      ],
      sorted: '',
      downloadLabel: 'CSV形式でダウンロード',
      tableForm: 'view', // 'view', 'markdown', 'csv'
      tableFormTabs: ['view', 'markdown', 'csv'],
      showColumn: []
    }
  },
  computed: {
    tableWidth() {
      return this.defaultColumns.length + this.questions.length
    },
    questionnaireId() {
      return this.$route.params.id
    },
    tableHeaders() {
      const ret = this.defaultColumns
        .map(column => column.label)
        .concat(this.questions)
      return ret
    },
    markdownTable() {
      // results の表を markdown 形式にしたものを返す
      let ret = ''

      // 項目の行
      ret += this.arrayToMarkdown(this.tableHeaders)

      // 2行目
      ret += this.arrayToMarkdown(new Array(this.tableWidth).fill('-'))

      // 各回答の行
      for (let i = 0; i < this.results.length; i++) {
        let arr = []
        const result = this.results[i]
        this.defaultColumns.forEach(header => {
          arr.push(result[header.name])
        })
        result.responseBody.forEach(body => {
          arr.push(this.responseToString(body))
        })

        ret += this.arrayToMarkdown(arr)
      }

      ret = ret.slice(0, -1) // 末尾の改行を削除

      return ret
    },
    csvTable() {
      const arrayToCsv = function(arr) {
        let ret = ''
        arr.forEach(val => {
          ret += '"' + val + '",'
        })
        ret = ret.slice(0, -1) // 最後の ',' は取り除く
        ret += '\n'
        return ret
      }

      let csv = '\ufeff'

      csv += arrayToCsv(
        this.tableHeaders.filter((_, index) => !this.isColumnHidden(index))
      )

      this.results
        .filter((_, index) => !this.isColumnHidden(index))
        .forEach(result => {
          const defaultResults = [result.traqId, result.submittedAt]
          csv += arrayToCsv(
            defaultResults
              .concat(
                result.responseBody.map(response =>
                  this.responseToString(response)
                )
              )
              .filter((_, index) => !this.isColumnHidden(index))
          )
        })
      return csv
    },
    canDownload() {
      return this.tableForm === 'markdown' || this.tableForm === 'csv'
    }
  },
  beforeUpdate() {
    if (this.showColumn.length < this.tableWidth) {
      this.showColumn = new Array(this.tableWidth).fill(true)
    }
  },
  methods: {
    getTableRow(index) {
      // 表のindex行目に表示する文字列の配列を返す
      const ret = this.defaultColumns
        .map(column => this.results[index][column.name])
        .concat(
          this.results[index].responseBody.map(response =>
            this.responseToString(response)
          )
        )
      return ret
    },
    responseToString(body) {
      let ret = ''
      switch (body.question_type) {
        case 'MultipleChoice':
        case 'Checkbox':
        case 'Dropdown':
          body.option_response.forEach(response => {
            if (ret !== '') {
              ret += ', '
            }
            ret += response
          })
          return ret
        default:
          return body.response
      }
    },
    downloadTable() {
      if (!this.canDownload) return
      let form = {}
      switch (this.tableForm) {
        case 'markdown':
          form = {
            type: 'text/markdown',
            ext: '.md',
            data: this.markdownTable
          }
          break
        case 'csv':
          form = { type: 'text/csv', ext: '.csv', data: this.csvTable }
          break
      }
      const blob = new Blob([form.data], { type: form.type })
      let link = document.createElement('a')
      link.href = window.URL.createObjectURL(blob)
      link.download = 'Result' + form.ext
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
    },
    sort(index) {
      let query = ''
      if (this.sorted !== index) {
        query += '-'
        this.sorted = index
      } else {
        this.sorted = -index
      }
      switch (index) {
        case 1:
          query += 'traqid'
          break
        case 2:
          query += 'submitted_at'
          break
        default:
          query += index - 2
      }
      this.$emit('get-results', '?sort=' + query)
    },
    arrayToMarkdown(arr) {
      // 受け取った配列1行分のmarkdownを返す
      let ret = '|'
      arr
        .filter((val, index) => !this.isColumnHidden(index))
        .forEach(val => {
          ret += ' ' + val + ' |'
        })
      ret += '\n'
      return ret
    },
    toggleShowColumn(index) {
      this.$set(this.showColumn, index, !this.showColumn[index])
    },
    isColumnActive(index) {
      return this.sorted === Math.abs(index + 1)
    },
    isColumnHidden(index) {
      return (
        this.showColumn.length === this.tableWidth && !this.showColumn[index]
      )
    }
  }
}
</script>

<style lang="scss" scoped>
th,
td {
  vertical-align: middle;
  font-size: 0.9em;
  min-width: 10em;
  word-break: break-all;
  &.hidden {
    color: $base-gray;
  }
}

th.active-header {
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

.download {
  margin: auto 0.5rem;
}

.card {
  overflow-x: unset;
}

.scroll-view {
  overflow-x: auto;
}

.button {
  display: block;
}

.header-wrapper {
  display: flex;
  [class*='header-'] {
    margin-top: auto;
    margin-bottom: auto;
  }
  .header-icon-left {
    display: flex;
    margin: auto 5px auto 0;
  }
  .header-icon-right {
    margin-left: 15px;
  }
}
</style>
