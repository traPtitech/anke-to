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
        <div v-show="tableForm === 'view'">
          <div
            v-for="question in countedData"
            :key="question.title"
            class="card"
          >
            <header class="card-header">
              <p class="card-header-title">{{ question.title }}</p>
            </header>
            <div class="card-content">
              <div class="content">
                <div v-if="isNumberType(question.type)" class="content">
                  <ul>
                    <li>平均値: {{ question.total.average }}</li>
                    <li>標準偏差: {{ question.total.standardDeviation }}</li>
                    <li>中央値: {{ question.total.median }}</li>
                    <li>最頻値: {{ question.total.mode }}</li>
                  </ul>
                </div>
                <div class="table-container">
                  <table class="table is-striped">
                    <thead>
                      <td>回答</td>
                      <td>回答数</td>
                      <td v-if="isSelectType(question.type)">選択率</td>
                      <td>その回答をした人</td>
                    </thead>
                    <tbody>
                      <tr v-for="[choice, ids] of question.data" :key="choice">
                        <td>{{ choice }}</td>
                        <td>{{ ids.length }}</td>
                        <td v-if="isSelectType(question.type)">
                          {{
                            `${((ids.length / question.length) * 100).toFixed(
                              2
                            )}%`
                          }}
                        </td>
                        <td>{{ ids.join(', ') }}</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- markdown view -->
        <textarea
          v-show="tableForm === 'markdown'"
          class="textarea"
          :value="markdownTable"
          :rows="markdownTable.split('\n').length + 3"
          readonly
        ></textarea>
      </div>
    </div>
  </div>
</template>

<script>
const isSelectType = type =>
  ['MultipleChoice', 'Checkbox', 'Dropdown'].includes(type)
const isNumberType = type => ['LinearScale', 'Number'].includes(type)

const countData = (questions, results) => {
  const data = Array.from({ length: questions.length }, () => [])
  const questionTypes = new Array(questions.length)
  for (const result of results) {
    const answers = result.responseBody
    for (const [i, answer] of answers.entries()) {
      const type = answer.question_type
      if (!questionTypes[i]) {
        questionTypes[i] = type
      }

      const datum = {
        traqId: result.traqId,
        modifiedAt: result.modifiedAt
      }
      if (isSelectType(type)) {
        datum.answer = answer.option_response
      } else if (isNumberType(type)) {
        datum.answer = +answer.response
      } else {
        datum.answer = answer.response
      }

      data[i].push(datum)
    }
  }
  return questions.map((q, i) => ({
    title: q,
    type: questionTypes[i],
    data: generateData(questionTypes[i], data[i]),
    total: generateTotal(questionTypes[i], data[i]),
    length: data[i].length
  }))
}

const generateData = (type, data) => {
  const total = new Map()
  for (const datum of data) {
    if (isSelectType(type)) {
      for (const answer of datum.answer) {
        if (!total.has(answer)) total.set(answer, [])
        total.get(answer).push(datum.traqId)
      }
    } else {
      const { answer } = datum
      if (!total.has(answer)) total.set(answer, [])
      total.get(answer).push(datum.traqId)
    }
  }

  let arr = [...total]
  if (isNumberType(type)) {
    arr = arr.sort((a, b) => b[0] - a[0])
  }
  return arr
}

const generateTotal = (type, data) => {
  if (isNumberType(type)) {
    const average =
      data.reduce((acc, datum) => acc + datum.answer, 0) / data.length
    const variance =
      data
        .map(datum => (datum.answer - average) ** 2)
        .reduce((acc, v) => acc + v, 0) / data.length
    return {
      average: average.toFixed(2),
      standardDeviation: Math.sqrt(variance).toFixed(2),
      median: getMedian(data),
      mode: getMode(data)
    }
  }
  return null
}

const getMedian = data => {
  const middle = Math.floor(data.length / 2)
  const sorted = data.sort((a, b) => a.answer - b.answer)
  if (data.length % 2 !== 0) {
    return sorted[middle].answer
  }
  return (sorted[middle - 1].answer + sorted[middle].answer) / 2
}

const getMode = data => {
  const map = new Map()
  for (const { answer } of data) {
    if (!map.has(answer)) map.set(answer, 0)
    map.set(answer, map.get(answer) + 1)
  }
  const arr = [...map].sort((a, b) => b[1] - a[1])
  return arr
    .filter(v => arr[0][1] === v[1])
    .map(v => v[0])
    .join(', ')
}

export default {
  name: 'Statistics',
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
  data() {
    return {
      tableForm: 'view', // 'view', 'markdown'
      tableFormTabs: ['view', 'markdown']
    }
  },
  computed: {
    questionnaireId() {
      return this.$route.params.id
    },
    countedData() {
      if (this.questions.length <= 0 || this.results.length <= 0) {
        return null
      }
      return countData(this.questions, this.results)
    },
    markdownTable() {
      if (!this.countedData) return ''
      return this.countedData
        .flatMap(question => {
          const { total, data } = question
          let res = [`# ${question.title}`]
          if (isNumberType(question.type)) {
            res = res.concat([
              `**平均値**: ${total.average}`,
              `**標準偏差**: ${total.standardDeviation}`,
              `**中央値**: ${total.median}`,
              `**最頻値**: ${total.mode}`,
              ''
            ])
          }
          if (isSelectType(question.type)) {
            res = res.concat(
              [
                '| 回答 | 回答数 | 選択率 | その回答をした人 |',
                '| - | - | - | - |'
              ],
              data.map(
                ([choice, ids]) =>
                  `| ${choice ? choice : ''} | ${ids.length} | ${(
                    (ids.length / question.length) *
                    100
                  ).toFixed(2)}% | ${ids.join(', ')} |`
              )
            )
          } else {
            res = res.concat(
              ['| 回答 | 回答数 | その回答をした人 |', '| - | - | - |'],
              data.map(([choice, ids]) => {
                const c = choice ? choice : ''
                return `| ${
                  isNumberType(question.type) ? c : c.replace(/\n/g, '<br>')
                } | ${ids.length} | ${ids.join(', ')} |`
              })
            )
          }
          return res.concat([''])
        })
        .join('\n')
    },
    canDownload() {
      return this.tableForm === 'markdown'
    }
  },
  methods: {
    downloadTable() {
      if (!this.canDownload) return
      let form = { type: 'text/markdown', ext: '.md', data: this.markdownTable }
      const blob = new Blob([form.data], { type: form.type })
      let link = document.createElement('a')
      link.href = window.URL.createObjectURL(blob)
      link.download = 'Result' + form.ext
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
    },
    isSelectType(type) {
      return isSelectType(type)
    },
    isNumberType(type) {
      return isNumberType(type)
    }
  }
}
</script>

<style lang="scss" scoped></style>
