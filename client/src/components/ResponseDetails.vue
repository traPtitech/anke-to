<template>
  <div class="is-fullheight details">
    <div class="tabs is-centered">
      <ul></ul>
      <a
        id="edit-button"
        :class="{'is-editing': isEditing}"
        @click.prevent="isEditing = !isEditing"
        v-if="!isNewResponse"
      >
        <span class="ti-pencil"></span>
      </a>
    </div>
    <div :class="{'is-editing' : isEditing}" class="is-fullheight details-child">
      <questions
        :traqId="traqId"
        :editMode="isEditing? 'response' : undefined"
        :questionsProps="questions"
        @set-questions="setQuestions"
        @set-question-content="setQuestionContent"
      ></questions>
    </div>
  </div>
</template>

<script>

// import <componentname> from '<path to component file>'
import axios from 'axios'
import router from '@/router'
import common from '@/util/common'
import Questions from '@/components/Questions'

export default {
  name: 'ResponseDetails',
  components: {
    'questions': Questions
  },
  async created () {
    this.getQuestions()
  },
  props: {
    traqId: {
      type: String,
      required: true
    }
  },
  data () {
    return {
      questions: []
    }
  },
  methods: {
    getQuestions () {
      let responseData = {}

      // 該当する回答のデータを取得して responseData に保存
      axios
        .get('/responses/' + this.responseId)
        .then(res => {
          responseData = res.data
          // console.log(responseData)
        })
        .then(() => {
          // 該当するアンケートの質問一覧を取得して、convertDataToQuestion を通したものを this.questions に保存
          axios
            .get('/questionnaires/' + responseData.questionnaireID + '/questions')
            .then(res => {
              res.data.forEach(data => {
                this.questions.push(common.convertDataToQuestion(data))
              })
            })
            .then(() => {
              // 各質問に対して、該当する回答の情報を this.questions に入れる
              // questions[i] の questionId と responseData.body[i] の questionID は一致するはず (怪しい)
              this.questions.forEach((question, index) => {
                if (question.questionId === responseData.body[ index ].questionID) {
                  this.$set(this.questions, index, common.setResponseToQuestion(question, responseData.body[ index ]))
                } else {
                  // questionとresponseのquestionIDが一致しなかった場合の処理 (未実装)
                }
              })
            })
        })
    },
    setQuestions (questions) {
      this.questions = questions
    },
    setQuestionContent (index, label, value) {
      console.log(index)
      console.log(value)
      // this.questions[ index ][ label ] = value
      let newQuestion = Object.assign({}, this.questions[ index ])
      newQuestion[ label ] = value
      this.$set(this.questions, index, newQuestion)
    }
  },
  computed: {
    responseId () {
      return this.isNewResponse ? '' : this.$route.params.id
    },
    isNewResponse () {
      return this.$route.params.id === 'new'
    },
    isEditing: {
      get: function () {
        if (this.isNewResponse || this.$route.hash === '#edit') {
          return true
        }
        return false
      },
      set: function (newBool) {
        if (newBool) {
          // 閲覧 -> 編集
          router.push('/responses/' + this.responseId + '#edit')
        } else {
          // 編集 -> 閲覧
          router.push('/responses/' + this.responseId)
        }
      }
    },
    editButtonLink () {
      if (!this.isEditing) {
        return this.responseId + '#edit'
      } else {
        return this.responseId
      }
    }
  },
  mounted () {
  }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style lang="scss" scoped>
</style>
