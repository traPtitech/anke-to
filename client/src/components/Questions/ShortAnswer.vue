<template>
  <div>
    <!-- view or edit question -->
    <p
      v-if="
        editMode === 'question' ||
        (editMode !== 'response' && typeof content.responseBody === 'undefined')
      "
      class="has-underline placeholder"
    >
      {{ responsePlaceholder }}
    </p>

    <!-- view response -->
    <p
      v-if="!editMode && typeof content.responseBody !== 'undefined'"
      class="has-underline"
      :class="{ 'multi-line': content.type === 'TextArea' }"
    >
      {{ content.responseBody }}
    </p>

    <!-- edit response -->
    <input
      v-if="editMode === 'response' && content.type === 'Text'"
      v-model="content.responseBody"
      type="text"
      class="input has-underline"
      placeholder="回答"
    />
    <input
      v-if="editMode === 'response' && content.type === 'Number'"
      v-model.number="content.responseBody"
      type="number"
      class="input has-underline"
      placeholder="0"
    />
    <textarea
      v-if="editMode === 'response' && content.type === 'TextArea'"
      v-model="content.responseBody"
      class="input has-underline"
      placeholder="回答"
      rows="5"
    ></textarea>
  </div>
</template>

<script>
// import <componentname> from '<path to component file>'

export default {
  name: '',
  components: {},
  props: {
    contentProps: {
      type: Object,
      required: true
    },
    editMode: {
      type: String, // 'question' or 'response'
      required: false, // 渡されなかった場合はview
      default: undefined
    }
  },
  data() {
    return {}
  },
  computed: {
    content() {
      return this.contentProps
    },
    responsePlaceholder() {
      if (this.content.type === 'Text' || this.content.type === 'TextArea') {
        return '回答 (テキスト)'
      } else if (this.content.type === 'Number') {
        return '回答 (数値)'
      } else {
        return ''
      }
    }
  },
  mounted() {},
  methods: {}
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style lang="scss" scoped>
.placeholder {
  color: $base-brown;
  &.has-underline {
    border-bottom: $base-brown dotted 0.5px;
  }
}

.multi-line {
  white-space: pre-line;
  word-break: break-all;
}

textarea {
  height: 10em;
}
</style>
