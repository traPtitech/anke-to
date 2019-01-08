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
                <div class="column left-bar" v-if="isQuestionEditMode">
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
                <div class="column question" :class="{'is-editable' : isQuestionEditMode}">
                  <p class="subtitle" v-if="!isQuestionEditMode">
                    {{ question.questionBody }}
                    <span
                      class="ti-alert required-question-icon"
                      v-if="showRequiredIcon(index)"
                    >必須</span>
                  </p>
                  <div v-if="isQuestionEditMode" class="columns is-inline-block-mobile">
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
                    :editMode="editMode"
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
                @click="insertNewQuestion(questionType)"
              >{{ questionType.label }}</button>
            </div>
            <!-- <div class="add-question-button">
              <span class="select">
                <select v-model="newQuestionType">
                  <option
                    v-for="questionType in questionTypes"
                    :key="questionType.type"
                    :value="questionType.type"
                  >{{ questionType.label }}</option>
                </select>
              </span>
              <span class="ti-plus circled icon"></span>
              <span>新しい質問を追加</span>
            </div>-->
            <!-- <div class="dropdown" :class="{'is-active' : newQuestionDropdownIsActive}">
              <div class="dropdown-trigger">
                <button class="button" @click="toggleNewQuestionDropdown()">
                  <span>新しい質問を追加</span>
                  <span class="ti-angle-down icon is-small"></span>
                </button>
              </div>
              <div class="dropdown-menu" id="dropdown-menu" role="menu">
                <div class="dropdown-content is-flex">
                  <button
                    class="dropdown-item button"
                    v-for="questionType in questionTypes"
                    :key="questionType.type"
                  >{{ questionType.label }}</button>
                </div>
              </div>
            </div>-->
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
    },
    editMode: {
      type: String,
      required: false
    }
  },
  data () {
    return {
      // selectedQuestionIndex: 0
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
      // newQuestionType: 'Text',
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
    showRequiredIcon (index) {
      return this.editMode !== 'question' && this.questions[ index ].isRequired
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
      if (questionType.type === 'Checkbox' || questionType.type === 'MultipleChoice') {
        ret.options = [ {label: '', id: 0} ]
      } else if (questionType.type === 'LinearScale') {
        ret.scaleRange = {left: 1, right: 5}
        ret.scaleLabels = {left: '', right: ''}
      }
      console.log(ret)
      return ret
    },
    getNewQuestionId () {
      const ret = this.lastQuestionId - 1
      this.lastQuestionId -= 1
      return ret
    },
    insertNewQuestion (questionType) {
      this.questions.push(this.getDefaultQuestion(questionType))
    }

    // isSelectedQuestion (index) {
    //   return this.editMode === 'question' && index === this.selectedQuestionIndex
    // },
    // getQuestionEditMode (index) {
    //   if (this.editMode !== 'question') return this.editMode
    //   if (this.isSelectedQuestion(index)) return 'question'
    //   else return undefined
    // },
    // setSelectedQuestionIndex (index) {
    //   if (this.editMode === 'question') this.selectedQuestionIndex = index
    // }
  },
  computed: {
    isQuestionEditMode () {
      return this.editMode === 'question'
    }
  },
  mounted () {
  }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style lang="scss" scoped>
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
    background: lightgray;
    padding: 0.5rem;
    margin: 0;
    width: fit-content;
    border-radius: 0.3rem;
  }
}
.required-question-icon {
  font-size: 0.8rem;
  color: red;
  margin: auto 0.5rem;
  word-break: keep-all;
}
.required-question-icon::before {
  margin-right: 0.2rem;
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
.questions /deep/ .question {
  .columns {
    width: 100%;
  }
  .column {
    padding: 0.2rem;
  }
  .column.question {
    padding: 0.5rem 1rem;
  }
  .column.left-bar {
    width: fit-content;
    max-width: fit-content;
    padding: 0;
  }
  input,
  textarea {
    border: none;
    outline: none;
    -webkit-box-shadow: none;
    box-shadow: none;
    border-radius: 0;
  }
  .question-body {
    margin-bottom: 0.5rem;
  }
  .response-body {
    margin: 0.5rem 0.5rem;
    p {
      margin-top: 1rem;
    }
    p.has-underline {
      padding-bottom: 0.25rem;
    }
  }
  .has-underline {
    padding: 0 0.5rem;
    border-bottom: grey dotted 0.5px;
  }
  .is-editable.has-underline {
    border-bottom: lightgrey solid 0.5px;
  }
  input:focus.has-underline {
    border-bottom: black solid 2px;
  }
}
</style>
