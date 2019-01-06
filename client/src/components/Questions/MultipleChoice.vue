<template>
  <div>
    <!-- view only -->
    <div v-if="!editMode">
      <div v-for="(option, index) in content.options" :key="index" class="is-flex">
        <span :class="[readOnlyBoxClass, isSelected(index) ? 'checked' : 'empty']"></span>
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
            v-model="content.options[index].label"
          >
          <span class="ti-trash icon is-medium" @click="removeOption(index)"></span>
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
        <input
          v-if="content.type==='Checkbox'"
          type="checkbox"
          :value="option"
          v-model="content.isSelected[index]"
        >
        <input
          v-if="content.type==='MultipleChoice'"
          type="radio"
          :value="option"
          v-model="content.selected"
        >
        <label class="option-label">{{ option.label }}</label>
      </div>
    </div>
  </div>
</template>

<script>

// import <componentname> from '<path to component file>'
import common from '@/util/common'

export default {
  name: 'MultipleChoice',
  components: {
  },
  props: {
    content: {
      type: Object,
      required: true
    },
    editMode: {
      type: String,
      required: false
    }
  },
  data () {
    return {
      options: []
    }
  },
  methods: {
    swapOrder: common.swapOrder,
    isFirstOption (index) {
      return index === 0
    },
    isLastOption (index) {
      return index === this.content.options.length - 1
    },
    isSelected (index) {
      if (this.content.type === 'Checkbox') {
        return this.content.isSelected[ index ]
      } else if (this.content.type === 'MultipleChoice') {
        return this.content.selected === this.content.options[ index ]
      }
    },
    addOption () {
      let newId = this.content.options.length
      this.content.options.push({
        id: newId,
        label: ''
      })
    },
    removeOption (index) {
      this.content.options.splice(index, 1)
    }
  },
  computed: {
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
  width: min-content;
}
.icon {
  cursor: pointer;
  margin: auto;
}
.icon.disabled {
  color: lightgray;
  pointer-events: none;
}
.icon.is-medium {
  font-size: 1.2rem;
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
