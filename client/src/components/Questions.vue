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
          <div class="card-content questions">
            <div
              v-for="(question, index) in questions"
              :key="index"
              class="question"
              :class="{'is-selected' : isSelectedQuestion(index)}"
              @click="setSelectedQuestionIndex(index)"
            >
              <div class="question-body">
                <p class="subtitle" v-if="!isSelectedQuestion(index)">{{ question.questionBody }}</p>
                <input
                  type="text"
                  class="subtitle input has-underline is-editable"
                  v-if="isSelectedQuestion(index)"
                  placeholder="質問文"
                  v-model="question.questionBody"
                >
              </div>
              <component
                :editMode="getQuestionEditMode(index)"
                :is="question.component"
                :content="question"
                class="response-body"
              ></component>
              <hr>
            </div>
          </div>
        </div>
      </article>
    </div>
  </div>
</template>

<script>

// import <componentname> from '<path to component file>'
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
    // questions: {
    //   type: Array,
    //   required: false
    // },
    editMode: {
      type: String,
      required: false
    }
  },
  data () {
    return {
      selectedQuestionIndex: 0,
      questions: [ // テスト用
        {
          questionId: 1,
          type: 'Text',
          component: 'short-answer',
          questionBody: 'ペンは持っていますか？',
          body: 'はい'
        },
        {
          questionId: 2,
          type: 'Number',
          component: 'short-answer',
          questionBody: 'ペンは何本持っていますか？',
          body: '12'
        },
        {
          questionId: 3,
          type: 'Checkbox',
          component: 'multiple-choice',
          questionBody: '何色のペンを持っていますか？',
          options: [
            {
              label: '赤',
              id: 0
            },
            {
              label: '青',
              id: 1
            },
            {
              label: '黄色',
              id: 2
            }
          ],
          isSelected: [ false, false, false ]
        },
        {
          questionId: 4,
          type: 'MultipleChoice',
          component: 'multiple-choice',
          questionBody: '何色のペンが欲しいですか？',
          options: [
            {
              label: '赤',
              id: 0
            },
            {
              label: '青',
              id: 1
            },
            {
              label: '黄色',
              id: 2
            }
          ],
          selected: ''
        }
      ]
    }
  },
  methods: {
    swapOrder: common.swapOrder,
    componentEditMode (index) {
      return this.editMode !== 'question' ? this.editMode : (index === this.selectedQuestionIndex ? 'question' : undefined)
    },
    isSelectedQuestion (index) {
      return this.editMode === 'question' && index === this.selectedQuestionIndex
    },
    getQuestionEditMode (index) {
      if (this.editMode !== 'question') return this.editMode
      if (this.isSelectedQuestion(index)) return 'question'
      else return undefined
    },
    setSelectedQuestionIndex (index) {
      if (this.editMode === 'question') this.selectedQuestionIndex = index
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
.question.is-selected {
  border-left: solid;
  padding-left: 1rem;
}
.questions /deep/ .question {
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
    margin: 0 0.5rem;
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
  .list-move {
    transition: transform 1s;
  }
  .list-leave-to,
  .list-enter {
    transition: all 1s;
    opacity: 0;
    transform: translateX(30px);
  }
  .list-leave-active {
    position: absolute;
  }
}
</style>
