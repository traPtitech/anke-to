<template>
  <div class="columns">
    <article class="column is-11">
      <div class="card">
        <div class="card-content questions">
          <div v-for="(question, index) in questions" :key="index" class="question">
            <div class="question-body">
              <p class="subtitle">
                {{ question.questionBody }}
                <span
                  class="ti-alert required-question-icon"
                  v-if="showRequiredIcon(index)"
                >必須</span>
              </p>
            </div>
            <input-error-message v-if="inputErrors" :inputError="inputErrors[question.questionId]"></input-error-message>
            <component
              :editMode="editMode"
              :is="question.component"
              :contentProps="question"
              :questionIndex="index"
              class="response-body"
              @set-question-content="setQuestionContent"
            ></component>
            <hr>
          </div>
        </div>
      </div>
    </article>
  </div>
</template>

<script>

import MultipleChoice from '@/components/Questions/MultipleChoice'
import LinearScale from '@/components/Questions/LinearScale'
import ShortAnswer from '@/components/Questions/ShortAnswer'
import InputErrorMessage from '@/components/Utils/InputErrorMessage'

// import common from '@/util/common'

export default {
  name: 'Questions',
  components: {
    'multiple-choice': MultipleChoice,
    'linear-scale': LinearScale,
    'short-answer': ShortAnswer,
    'input-error-message': InputErrorMessage
  },
  props: {
    questionsProps: {
      type: Array,
      required: false
    },
    editMode: {
      type: String,
      required: false
    },
    inputErrors: {
      type: Object,
      required: false
    }
  },
  data () {
    return {
    }
  },
  methods: {
    showRequiredIcon (index) {
      return this.questions[ index ].isRequired
    },
    setQuestionContent (index, label, value) {
      this.$emit('set-question-content', index, label, value)
    }
  },
  computed: {
    questions () {
      return this.questionsProps
    }
  },
  mounted () {
  }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style lang="scss" scoped>
.required-question-icon {
  font-size: 0.8rem;
  color: red;
  margin: auto 0.5rem;
  word-break: keep-all;
}
.required-question-icon::before {
  margin-right: 0.2rem;
}
</style>
