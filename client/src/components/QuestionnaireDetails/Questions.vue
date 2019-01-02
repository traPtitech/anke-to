<template>
  <div>
    <div class="columns">
      <article class="column is-11">
        <div class="card">
          <div>
            <header class="card-header">
              <div id="title" class="card-header-title title">
                <div>タイトル</div>
              </div>
            </header>
            <div class="card-content questions">
              <div v-for="(question, index) in questions" :key="index">
                <component :is="question.component" :content="question.content" class="question"></component>
                <hr>
              </div>
              <!-- <short-answer class="question" :content="shortAnswerContent"></short-answer>
              <hr>
              <short-answer class="question" :content="shortAnswerContent" editMode="question"></short-answer>
              <hr>
              <short-answer class="question" :content="shortAnswerContent" editMode="response"></short-answer>-->
              <div v-for="(question, index) in questions" :key="index">
                <component
                  :is="question.component"
                  editMode="question"
                  :content="question.content"
                  class="question"
                ></component>
                <hr>
              </div>

              <div v-for="(question, index) in questions" :key="index">
                <component
                  :is="question.component"
                  editMode="response"
                  :content="question.content"
                  class="question"
                ></component>
                <hr>
              </div>
            </div>
          </div>
        </div>
      </article>
    </div>
  </div>
</template>

<script>

// import <componentname> from '<path to component file>'
import Checkbox from '@/components/Questions/Checkbox'
import Dropdown from '@/components/Questions/Dropdown'
import LinearScale from '@/components/Questions/LinearScale'
import Radiobutton from '@/components/Questions/Radiobutton'
import ShortAnswer from '@/components/Questions/ShortAnswer'

export default {
  name: 'Questions',
  components: {
    'checkbox': Checkbox,
    'dropdown': Dropdown,
    'linear-scale': LinearScale,
    'radiobutton': Radiobutton,
    'short-answer': ShortAnswer
  },
  props: {
    // name: {
    //   type: ,
    //   required:
    // }
  },
  data () {
    return {
      questions: [
        {
          questionId: 1,
          type: 'Number',
          component: 'short-answer',
          content: {
            questionBody: 'ペンは持っていますか？',
            responseBody: 'はい',
            responseType: 'text'
          }
        },
        {
          questionId: 2,
          type: 'Text',
          component: 'short-answer',
          content: {
            questionBody: 'ペンは何本持っていますか？',
            responseBody: '12',
            responseType: 'number'
          }
        }
      ],
      shortAnswerContent: {
        questionBody: 'ペンは持っていますか？',
        responseBody: '12',
        responseType: 'number'
      }
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
<style lang="scss" scoped>
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
    .has-underline {
      border-bottom: lightgrey solid 0.5px;
    }
    input:focus.has-underline {
      border-bottom: black solid 2px;
    }
  }
  .response-body {
    margin: 0 0.5rem;
    .has-underline {
      padding: 0 0.5rem;
      border-bottom: grey dotted 0.5px;
    }
    input:focus.has-underline {
      border-bottom: black solid 2px;
    }
    p {
      margin-top: 1rem;
    }
    p.has-underline {
      padding-bottom: 0.25rem;
    }
  }
}
</style>
