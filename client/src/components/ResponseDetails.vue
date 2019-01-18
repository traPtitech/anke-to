<template>
  <div class="is-fullheight details" :class="{'has-navbar-fixed-bottom': isEditing}">
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
    <edit-nav-bar
      v-if="isEditing"
      :editButtons="editButtons"
      @submit-response="submitResponse"
      @save-response="saveResponse"
      @disable-editing="disableEditing"
    ></edit-nav-bar>
  </div>
</template>

<script>

// import <componentname> from '<path to component file>'
import axios from 'axios'
import router from '@/router'
import common from '@/util/common'
import Questions from '@/components/Questions'
import EditNavBar from '@/components/Utils/EditNavBar.vue'

export default {
  name: 'ResponseDetails',
  components: {
    'questions': Questions,
    'edit-nav-bar': EditNavBar
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
      responseData: {}.isEditing
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
    },
    sendResponse (data) {
      if (this.isNewResponse) {
        axios
          .post('/responses', data)
          .then(resp => {
            const responseId = resp.data.responseID
            router.push({
              name: 'ResponseDetails',
              params: {id: responseId}
            })
          })
      } else {
        axios
          .patch('/responses/' + this.responseId, data)
          .then(() => {
            this.isEditing = false
          })
      }
    },
    submitResponse () {
      // 回答の送信
      let data = this.createResponseData()
      data.submitted_at = new Date().toLocaleString('ja-GB')
      this.sendResponse(data)
    },
    saveResponse () {
      // 回答の保存
      let data = this.createResponseData()
      data.submitted_at = 'NULL'
      this.sendResponse(data)
    },
    disableEditing () {
      if (this.isNewResponse) {
        // 新しい回答の場合は、アンケートの詳細画面に戻る
        router.push({
          name: 'QuestionnaireDetails',
          params: {id: this.questionnaireId}
        })
      } else {
        this.isEditing = false
      }
    },
    createResponseData () {
      // サーバーに送るフォーマットのresponseDataを作成する
      let data = {
        questionnaireID: this.questionnaireId,
        body: []
      }
      this.questions.forEach(question => {
        let body = {
          questionID: question.questionId,
          question_type: question.type,
          response: '',
          option_response: []
        }
        switch (question.type) {
          case 'MultipleChoice':
            body.option_response = [ question.selected ]
            break
          case 'Checkbox':
            Object.keys(question.isSelected).forEach(key => {
              if (question.isSelected[ key ]) {
                body.option_response.push(key)
              }
            })
            break
          case 'Text':
          case 'Number':
            body.option_response = [ '' ]
            body.response = String(question.responseBody)
            break
          case 'LinearScale':
            body.option_response = [ '' ]
            body.response = String(question.selected)
            break
        }
        data.body.push(body)
      })
      return data
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
    },
    submitOk () {
      // 未実装
      return true
    },
    editButtons () {
      return [
        {
          label: '送信',
          atClick: 'submit-response',
          disabled: !this.submitOk
        },
        {
          label: '保存',
          atClick: 'save-response',
          disabled: false
        },
        {
          label: 'キャンセル',
          atClick: 'disable-editing',
          disabled: false
        }
      ]
    }
  },
  mounted () {
  }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style lang="scss" scoped>
</style>
