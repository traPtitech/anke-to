<template>
  <div class="card-wrapper is-fullheight">
    <div class="card">
      <header class="card-header">
        <div class="card-header-title subtitle">自分の回答</div>
      </header>
      <div class="card-content">
        <table class="table is-striped">
          <thead>
            <tr>
              <th v-for="(header, index) in headers" :key="index">{{ header }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(response, index) in responses" :key="index">
              <td class="table-item-title">
                <router-link
                  :to="'/questionnaires/' + response.questionnaireID"
                >{{ response.questionnaire_title }}</router-link>
              </td>
              <td class="table-item-date">{{ getDateStr(response.res_time_limit) }}</td>
              <td
                class="table-item-date"
              >{{ response.submitted_at == 'NULL' ? '未提出' : getRelativeDateStr(response.submitted_at) }}</td>
              <td class="table-item-date">{{ getRelativeDateStr(response.modified_at) }}</td>
              <td>
                <router-link :to="'/responses/' + response.responseID" target="_blank">
                  <span class="ti-new-window"></span>
                  <br>Open
                </router-link>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script>
import axios from '@/bin/axios'
import Table from '@/components/Utils/Table.vue'
import common from '@/util/common'

export default {
  name: 'Responses',
  components: {
    customtable: Table
  },
  async created () {
    axios.get('/users/me/responses').then(resp => {
      this.responses = resp.data
    })
  },
  props: {
    traqId: {
      required: true
    }
  },
  data () {
    return {
      responses: [],
      headers: [ '', '回答期限', '回答日時', '更新日時', '回答' ]
    }
  },
  computed: {},
  methods: {
    getDateStr (str) {
      return common.customDateStr(str)
    },
    getRelativeDateStr (str) {
      return common.relativeDateStr(str)
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
.dropdowns {
  padding-left: 1.5rem;
  padding-right: 1.5rem;
  padding-top: 1.5rem;
}
.dropdown-content:hover {
  background-color: lightgray;
}
.dropdown-content.is-selected {
  background-color: darkgray;
}
.button p {
  margin-right: 0.5em;
}
.card {
  margin: 0;
}
.card-wrapper {
  padding: 1rem 1.5rem;
}
</style>
