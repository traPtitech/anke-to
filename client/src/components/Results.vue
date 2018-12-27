<template>
  <div>questionnaireId = {{ questionnaireId }} のアンケートの結果</div>
</template>

<script>

// import <componentname> from '<path to component file>'

export default {
  name: 'Results',
  components: {
  },
  props: {
    // name: {
    //   type: ,
    //   required:
    // }
  },
  data () {
    return {
    }
  },
  methods: {
  },
  computed: {
    questionnaireId () {
      return this.$route.params.id
    }
  },
  mounted () {
  }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style scoped>
</style>
