<template>
  <div>
    <!-- view or edit question -->
    <p
      class="has-underline placeholder"
      v-if="editMode==='question' || (editMode!=='response' && typeof content.responseBody==='undefined')"
    >{{ responsePlaceholder }}</p>

    <!-- view response -->
    <p
      class="has-underline"
      v-if="!editMode && typeof content.responseBody !=='undefined'"
    >{{ content.responseBody }}</p>

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
    },
    responsePlaceholder () {
      if (this.content.type === 'Text') {
        return '回答 (テキスト)'
      } else if (this.content.type === 'Number') {
        return '回答 (数値)'
      } else {
        return ''
      }
    }
  },
  mounted () {
  }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style lang="scss" scoped>
@import "@/css/variables.scss";

.placeholder {
  color: $base-brown;
  &.has-underline {
    border-bottom: $base-brown dotted 0.5px;
  }
}
</style>
