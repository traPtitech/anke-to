<template>
  <div>
    <pagination :range="range" :currentPage="currentPage" :getPageLink="getPageLink"></pagination>
  </div>
</template>

<script>

import Pagination from '@/components/Utils/Pagination'

export default {
  name: 'Individual',
  components: {
    'pagination': Pagination
  },
  props: {
    results: {
      type: Array,
      required: true
    }
  },
  data () {
    return {
    }
  },
  methods: {
    getPageLink (pageName) {
      let ret = {
        name: 'Results',
        params: {id: this.$route.params.id},
        hash: '#individual'
      }
      switch (pageName) {
        case 'first':
          ret.query = {page: this.range.first}
          break
        case 'last':
          ret.query = {page: this.range.last}
          break
        default:
          ret.query = {page: pageName}
          break
      }
      return ret
    }
  },
  computed: {
    currentPage () {
      return this.$route.query.page ? Number(this.$route.query.page) : this.range.first
    },
    range () {
      return {
        first: 1,
        last: this.results.length
      }
    }
  },
  mounted () {
  }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style lang="scss" scoped>
</style>
