<template>
  <div class="hello">
    {{msg}}
    <br>
    {{questionnaires}}
  </div>
</template>

<script>
import axios from 'axios'

if (process.env.NODE_ENV === 'development') {
  axios.defaults.baseURL = 'http://localhost:1323'
}

export default {
  name: 'HelloWorld',
  data () {
    return {
      msg: 'anke-to',
      questionnaires: []
    }
  },
  async created () {
    const resp = await axios.get('/questionnaire')
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
