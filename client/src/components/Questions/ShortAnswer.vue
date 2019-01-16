<template>
  <div>
    <!-- view only -->
    <p class="has-underline" v-if="!editMode">{{ content.responseBody }}</p>

    <!-- edit question -->
    <p class="has-underline placeholder" v-if="editMode==='question'">回答</p>

    <!-- edit response -->
    <input
      type="text"
      class="input has-underline"
      v-if="editMode==='response' && content.type==='Text'"
      placeholder="回答"
      v-model="content.responseBody"
    >
    <input
      type="number"
      class="input has-underline"
      v-if="editMode==='response' && content.type==='Number'"
      placeholder="0"
      v-model.number="content.responseBody"
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
    contentProps: {
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
    content () {
      return this.contentProps
    }
  },
  mounted () {
  }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style lang="scss" scoped>
.placeholder {
  color: lightgrey;
  &.has-underline {
    border-bottom: lightgrey dotted 0.5px;
  }
}
</style>
