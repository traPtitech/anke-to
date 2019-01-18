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
    if (this.isNewResponse) {
      this.getQuestions()
    } else {
      this.getResponseData()
        .then(this.getQuestions)
        .then(this.setResponsesToQuestions)
    }
  },
  props: {
    traqId: {
      required: true
    },
    isNewResponse: {
      type: Boolean,
      required: false
    }
  },
  data () {
    return {
      questions: [],
      responseData: {}
    }
  },
  methods: {
    getResponseData () {
      return axios
        .get('/responses/' + this.responseId)
        .then(res => {
          this.responseData = res.data

          // questionIdをキーにしてresponseData.body の各要素をとれるようにする
          let newBody = {}
          this.responseData.body.forEach(data => {
            newBody[ data.questionID ] = data
          })
          this.responseData.body = newBody
        })
    },
    getQuestions () {
      return axios
        .get('/questionnaires/' + this.questionnaireId + '/questions')
        .then(res => {
          // convertDataToQuestion を通したものを this.questions に保存
          res.data.forEach(data => {
            this.questions.push(common.convertDataToQuestion(data))
          })
        })
    },
    setResponsesToQuestions () {
      // 各質問に対して、該当する回答の情報を this.questions に入れる
      this.questions.forEach((question, index) => {
        this.$set(this.questions, index, common.setResponseToQuestion(question, this.responseData.body[ question.questionId ]))
      })
    }
  },
  computed: {
    responseId () {
      return this.isNewResponse ? undefined : Number(this.$route.params.id)
    },
    questionnaireId () {
      if (this.isNewResponse) {
        return Number(this.$route.params.questionnaireId)
      } else {
        return this.responseData.questionnaireID
      }
    },
    isEditing: {
      get: function () {
        if (this.isNewResponse || this.$route.hash === '#edit') {
          return true
        }
        return false
      },
      set: function (newBool) {
        // newBool : 閲覧 -> 編集
        // !newBool : 編集 -> 閲覧
        const newRoute = {
          name: 'ResponseDetails',
          params: {id: this.responseId},
          hash: newBool ? '#edit' : undefined
        }
        router.push(newRoute)
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
