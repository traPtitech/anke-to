<template>
  <div class="wrapper">
    <p>{{ msg }}</p>
  </div>
</template>
<script>
// import axios from '@/bin/axios'

export default {
  name: 'CreatedLayout',
  data () {
    return {
      msg: '自分が作成したアンケートの一覧'
    }
  }
}
</script>
