<template>
  <div class="has-navbar-fixed-bottom">
    <questions :questionsProps="questionData"></questions>
    <pagination :range="range" :currentPage="currentPage" :defaultPageLink="defaultPageLink"></pagination>
  </div>
</template>

<script>

import Questions from '@/components/Questions'
import Pagination from '@/components/Utils/Pagination'

export default {
  name: 'Individual',
  async created () {
  },
  components: {
    'questions': Questions,
    'pagination': Pagination
  },
  props: {
    results: {
      type: Array,
      required: true
    },
    responseData: {
      type: Object,
      required: true
    },
    questionData: {
      type: Array,
      required: true
    },
    information: {
      type: Object,
      required: true
    }
  },
  data () {
    return {
    }
  },
  methods: {
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
    },
    defaultPageLink () {
      return {
        name: 'Results',
        params: { id: this.$route.params.id },
        query: { tab: 'individual' }
      }
    }
  },
  watch: {
  }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style lang="scss" scoped>
</style>
