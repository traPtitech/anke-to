<template>
  <div>
    <!-- view only -->
    <div v-if="!editMode">
      <div v-for="(option, index) in content.options" :key="index" class="is-flex">
        <span :class="[readOnlyBoxClass, isSelected(option.label) ? 'checked' : 'empty']"></span>
        <span class="option-label">{{ option.label }}</span>
      </div>
    </div>

    <!-- edit question -->
    <div v-if="editMode==='question'">
      <transition-group name="list" tag="div">
        <div v-for="(option, index) in content.options" :key="option.id" class="is-flex">
          <span class="sort-handle">
            <span
              class="ti-angle-up icon"
              @click="swapOrder(content.options, index, index-1)"
              :class="{disabled: isFirstOption(index)}"
            ></span>
            <span
              class="ti-angle-down icon"
              @click="swapOrder(content.options, index, index+1)"
              :class="{disabled: isLastOption(index)}"
            ></span>
          </span>
          <span :class="readOnlyBoxClass"></span>
          <input
            type="text"
            class="input has-underline option-label is-editable"
            :value="content.options[index].label"
            @input="setOption(index, $event.target.value)"
          >
          <span class="delete-button">
            <span class="ti-trash icon is-medium" @click="removeOption(index)"></span>
          </span>
        </div>
      </transition-group>
      <div class="wrapper add-option">
        <div class="add-option-button" @click="addOption()">
          <span class="ti-plus circled icon"></span>
          <span>新しい選択肢を追加</span>
        </div>
      </div>
    </div>

    <!-- edit response -->
    <div v-if="editMode==='response'">
      <div class="is-flex" v-for="(option, index) in content.options" :key="index">
        <label class="option-label">
          <input
            v-if="content.type==='Checkbox'"
            type="checkbox"
            :value="option.label"
            v-model="contentProps.isSelected[option.label]"
          >
          <input
            v-if="content.type==='MultipleChoice'"
            type="radio"
            :value="option.label"
            v-model="contentProps.selected"
          >
          {{ option.label }}
        </label>
      </div>
    </div>
  </div>
</template>

<script>

// import <componentname> from '<path to component file>'
import common from '@/util/common'
import question from '@/util/question'

export default {
  name: 'MultipleChoice',
  components: {
  },
  props: {
    contentProps: {
      type: Object,
      required: true
    },
    questionIndex: {
      type: Number,
      required: false
    },
    editMode: {
      type: String,
      required: false
    }
  },
  data () {
    return {
      newId: -1
    }
  },
  methods: {
    swapOrder: common.swapOrder,
    setContent: question.setContent,
    setOption (index, value) {
      let newOptions = Object.assign({}, this.content.options)
      newOptions[ index ].label = value
      this.setContent('options', newOptions)
    },
    isFirstOption (index) {
      return index === 0
    },
    isLastOption (index) {
      return index === this.content.options.length - 1
    },
    isSelected (label) {
      if (this.content.type === 'Checkbox') {
        return this.content.isSelected[ label ]
      } else if (this.content.type === 'MultipleChoice') {
        return this.content.selected === label
      }
    },
    updateOptionId () {
      this.newId--
    },
    addOption () {
      let newOptions = Object.assign([], this.content.options)
      newOptions.push({
        id: this.newId,
        label: ''
      })
      this.updateOptionId()
      this.setContent('options', newOptions)
    },
    removeOption (index) {
      let newOptions = Object.assign([], this.content.options)
      newOptions.splice(index, 1)
      this.setContent('options', newOptions)
    }
  },
  computed: {
    content () {
      return this.contentProps
    },
    readOnlyBoxClass () {
      switch (this.content.type) {
        case 'Checkbox':
          return 'readonly-checkbox'
        case 'MultipleChoice':
          return 'readonly-radiobutton'
        default:
          return undefined
      }
    }
  },
  watch: {
  },
  mounted () {
  }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style lang="scss" scoped>
input[type="checkbox"] {
  margin: auto 0;
}
.option-label {
  width: 100%;
  padding: 0 0.5rem;
  margin: auto;
}
.sort-handle {
  margin-right: 0.7rem;
}
.delete-button {
  height: fit-content;
  margin: auto;
}
.wrapper.add-option {
  display: flex;
  height: 2.5rem;
}
.add-option-button {
  margin: auto;
  cursor: pointer;
}
</style>
