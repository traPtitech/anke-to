<template>
  <div class="wrapper">
    <div class="card">
      <table class="table is-striped">
        <thead>
          <tr>
            <th
              v-for="(header, index) in headers.concat(questions)"
              :key="index"
            >{{ header }}</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="(result, index) in results"
            :key="index"
          >
            <td class="table-item-traqid">{{ result.traqID }}</td>
            <td class="table-item-time">{{ getDateStr(result.submitted_at) }}</td>
            <td v-for="response in result.response_body">
              {{getResponse(response)}}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script>
// import <componentname> from '<path to component file>'

import common from '@/util/common'

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
  data() {
    return {
      headers: ['traQID', '回答日時']
    }
  },
  methods: {
    getDateStr(str) {
      return common.customDateStr(str)
    },
    getResponse(body) {
      switch (body.question_type) {
        case 'MultipleChoice':
        case 'Checkbox':
        case 'Dropdown':
          let ret = ''
          for (const response of body.option_response) {
            if (ret != '') {
              ret += ', '
            }
            ret += response
          }
          return ret
        default:
          return body.response
      }
    }
  },
  computed: {},
  mounted() {}
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style lang="scss" scoped>
td {
  vertical-align: middle;
  font-size: 0.9em;
}
th,
td {
  min-width: 10em;
}
</style>
