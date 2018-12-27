<template>
  <div class="wrapper">
    <p>{{ msg }}</p>
  </div>
</template>
<script>
// import axios from '@/bin/axios'

export default {
  name: 'Administrates',
  data () {
    return {
      msg: '自分が管理者になっているアンケートの一覧'
    }
  }
}
</script>
