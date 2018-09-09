<template>
  <div class="wrapper">
    <table class="table is-striped">
      <thead>
        <tr>
          <th v-for="header in headers" :key="header.id">
            {{ header }}
          </th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="questionnaire in questionnaires" :key="questionnaire.questionnaireID">
          <td></td>
          <td>
            {{ questionnaire.title }}
          </td>
          <td>
            {{ questionnaire.res_time_limit }}
          </td>
          <td></td>
          <td>
            {{ questionnaire.modified_at }}
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
import axios from '@/bin/axios'

export default {
  name: 'MypageLayout',
  data () {
    return {
      msg: 'mypage',
      questionnaires: [],
      headers: ['', 'Title', 'Time Limit', 'Response', 'Modified At', 'Results', 'Details']
    }
  },
  async created () {
    const resp = await axios.get('/questionnaires')
    this.questionnaires = resp.data
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h1,
h2 {
  font-weight: normal;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>
