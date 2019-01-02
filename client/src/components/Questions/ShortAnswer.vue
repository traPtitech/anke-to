<template>
  <div>
    <div class="question-body">
      <p class="subtitle" v-if="editMode!=='question'">{{ content.questionBody }}</p>
      <input
        type="text"
        class="subtitle input has-underline"
        v-if="editMode==='question'"
        placeholder="質問文"
        v-model="content.questionBody"
      >
    </div>
    <div class="response-body">
      <p class="has-underline" v-if="editMode!=='response'">{{ content.responseBody }}</p>
      <input
        type="text"
        class="input has-underline"
        v-if="editMode==='response' && content.responseType==='text'"
        placeholder="回答"
        v-model="content.responseBody"
      >
      <input
        type="number"
        class="input has-underline"
        v-if="editMode==='response' && content.responseType==='number'"
        placeholder="0"
        v-model.number="content.responseBody"
      >
    </div>
  </div>
</template>

<script>

// import <componentname> from '<path to component file>'

export default {
  name: '',
  components: {
  },
  props: {
    content: {
      type: Object,
      required: true
    },
    editMode: {
      type: String, // 'question' or 'response'
      required: false // 渡されなかった場合はview
    }
  },
  data () {
    return {
    }
  },
  methods: {
  },
  computed: {
  },
  mounted () {
  }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style scoped>
</style>
