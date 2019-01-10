<template>
  <div>
    <div class="columns">
      <article class="column is-11">
        <div class="card">
          <header class="card-header">
            <div id="title" class="card-header-title title">
              <div>タイトル</div>
            </div>
          </header>
          <transition-group name="list" tag="div" class="card-content questions">
            <div v-for="(question, index) in questions" :key="question.questionId" class="question">
              <div class="question-body columns">
                <div class="column left-bar">
                  <div class="sort-handle">
                    <span
                      class="ti-angle-up icon"
                      @click="swapOrder(questions, index, index-1)"
                      :class="{disabled: isFirstQuestion(index)}"
                    ></span>
                    <span
                      class="ti-angle-down icon"
                      @click="swapOrder(questions, index, index+1)"
                      :class="{disabled: isLastQuestion(index)}"
                    ></span>
                  </div>
                  <div class="delete-button">
                    <span class="ti-trash icon is-medium" @click="removeQuestion(index)"></span>
                  </div>
                </div>
                <div class="column question is-editable">
                  <div class="columns is-inline-block-mobile">
                    <div class="column">
                      <input
                        type="text"
                        class="subtitle input has-underline is-editable"
                        placeholder="質問文"
                        v-model="question.questionBody"
                      >
                    </div>
                    <div class="column is-2 required-question-checkbox is-pulled-right">
                      <label class="checkbox">
                        必須
                        <input type="checkbox" v-model="question.isRequired">
                      </label>
                    </div>
                  </div>
                  <component
                    :editMode="'question'"
                    :is="question.component"
                    :content="question"
                    class="response-body"
                  ></component>
                </div>
              </div>
              <hr>
            </div>
          </transition-group>
          <div class="add-question">
            <div class="add-question-button button" @click="toggleNewQuestionDropdown">
              <span>新しい質問を追加</span>
              <span
                class="icon is-small"
                :class=" newQuestionDropdownIsActive ? 'ti-angle-down' : 'ti-angle-right'"
              ></span>
            </div>
            <div v-show="newQuestionDropdownIsActive" class="question-type-buttons">
              <button
                v-for="questionType in questionTypes"
                :key="questionType.type"
                class="button"
                @click="insertQuestion(questionType)"
              >{{ questionType.label }}</button>
            </div>
          </div>
        </div>
      </article>
    </div>
  </div>
</template>

<script>

import MultipleChoice from '@/components/Questions/MultipleChoice'
import LinearScale from '@/components/Questions/LinearScale'
import ShortAnswer from '@/components/Questions/ShortAnswer'
import common from '@/util/common'

export default {
  name: 'Questions',
  components: {
    'multiple-choice': MultipleChoice,
    'linear-scale': LinearScale,
    'short-answer': ShortAnswer
  },
  props: {
    questions: {
      type: Array,
      required: false
    }
  },
  data () {
    return {
      newQuestionDropdownIsActive: false,
      questionTypes: [
        {
          type: 'Text',
          label: 'テキスト',
          component: 'short-answer'
        },
        {
          type: 'Number',
          label: '数値',
          component: 'short-answer'
        },
        {
          type: 'Checkbox',
          label: 'チェックボックス',
          component: 'multiple-choice'
        },
        {
          type: 'MultipleChoice',
          label: 'ラジオボタン',
          component: 'multiple-choice'
        },
        {
          type: 'LinearScale',
          label: '目盛り',
          component: 'linear-scale'
        }
      ],
      lastQuestionId: 0
    }
  },
  methods: {
    swapOrder: common.swapOrder,
    isFirstQuestion (index) {
      return index === 0
    },
    isLastQuestion (index) {
      return index === this.questions.length - 1
    },
    removeQuestion (index) {
      this.questions.splice(index, 1)
    },
    insertQuestion (questionType) {
      this.questions.push(this.getDefaultQuestion(questionType))
    },
    toggleNewQuestionDropdown () {
      this.newQuestionDropdownIsActive = !this.newQuestionDropdownIsActive
    },
    getDefaultQuestion (questionType) {
      let ret = {
        questionId: this.getNewQuestionId(),
        type: questionType.type,
        component: questionType.component,
        questionBody: '',
        isRequired: false
      }
      switch (questionType.type) {
        case 'Checkbox':
          ret.options = [ {label: '', id: 0} ]
          ret.isSelected = [ false ]
          break
        case 'MultipleChoice':
          ret.options = [ {label: '', id: 0} ]
          break
        case 'LinearScale':
          ret.scaleRange = {left: 1, right: 5}
          ret.scaleLabels = {left: '', right: ''}
          break
      }
      return ret
    },
    getNewQuestionId () {
      const ret = this.lastQuestionId - 1
      this.lastQuestionId -= 1
      return ret
    }
  },
  computed: {
  },
  mounted () {
  }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style lang="scss" scoped>
.checkbox {
  // 1行に1つの選択肢
  display: block;
}
.card-content {
  padding: 1rem 0 0 0;
}
.question-body > .column {
  padding: 0.5rem;
}
.question.is-editable {
  border-left: solid;
  padding-left: 1rem;
  padding: 0 0.5rem;
}
.wrapper.add-option {
  display: flex;
  height: 2.5rem;
}
.add-option-button {
  margin: auto;
  cursor: pointer;
}
.required-question-checkbox {
  margin: auto;
  label.checkbox {
    padding: 0.5rem;
    // margin: 0 0.2rem;
    width: fit-content;
    border-left: solid 1px gray;
  }
}
.sort-handle {
  height: 5rem;
}
.add-question {
  margin-bottom: 1rem;
}
.add-question-button {
  background: whitesmoke;
  .icon {
    margin-right: 0.3rem;
    margin-left: 0.3rem;
  }
  .is-small {
    font-size: 0.9rem;
  }
}
.question-type-buttons {
  background: whitesmoke;
  padding: 0.5rem;
  border-radius: 0.2rem;
  button {
    margin: 0.2rem;
  }
}
.question:first-child {
  padding-top: 1rem;
}
</style>
