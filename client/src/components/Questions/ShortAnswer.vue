<template>
  <div>
    <p class="has-underline" v-if="!editMode">{{ content.body }}</p>
    <p class="has-underline placeholder" v-if="editMode==='question'">回答</p>
    <input
      type="text"
      class="input has-underline"
      v-if="editMode==='response' && content.type==='Text'"
      placeholder="回答"
      v-model="content.body"
    >
    <input
      type="number"
      class="input has-underline"
      v-if="editMode==='response' && content.type==='Number'"
      placeholder="0"
      v-model.number="content.body"
    >
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
.placeholder {
  color: grey;
}
</style>
