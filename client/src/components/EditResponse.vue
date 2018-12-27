<template>
  <div>
    <p>questionnaireId = {{ questionnaireId }} のアンケートの回答編集画面</p>
  </div>
</template>

<script>

// import <componentname> from '<path to component file>'

export default {
  name: 'EditResponse',
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
