<template>
  <div class="details is-fullheight">
    <div class="tabs is-centered">
      <ul>
        <li
          class="tab"
          :class="{ 'is-active': selectedTab===tab }"
          v-for="(tab, index) in detailTabs"
          :key="index"
          @click="selectedTab = tab"
        >
          <a>{{ tab }}</a>
        </li>
      </ul>
      <a
        @click="isEditing = !isEditing"
        id="edit-button"
        :class="{'is-editing': isEditing}"
        v-if="showEditButton"
      >
        <span class="ti-pencil"></span>
      </a>
    </div>
    <div
      :class="{'is-editing has-navbar-fixed-bottom' : isEditing}"
      class="details-child is-fullheight"
    >
      <information-summary v-if="currentTabComponent!=='information-edit'" :details="summaryProps"></information-summary>
      <component
        :is="currentTabComponent"
        :traqId="traqId"
        class="details-child is-fullheight"
        :name="currentTabComponent"
        :editMode="isEditing ? 'question' : undefined"
        :informationProps="informationProps"
        :questionsProps="questions"
        :title="title"
        :inputErrors="isEditing ? inputErrors: undefined"
        @set-data="setData"
        @set-question-content="setQuestionContent"
        @remove-question="removeQuestion"
      ></component>
      <edit-nav-bar
        v-if="isEditing"
        :editButtons="editButtons"
        @submit-questionnaire="submitQuestionnaire"
        @abort-editing="abortEditing"
      ></edit-nav-bar>
    </div>
  </div>
</template>

<script>

import moment from 'moment'
import router from '@/router'
import InformationSummary from '@/components/InformationSummary'
import Information from '@/components/QuestionnaireDetails/Information'
import InformationEdit from '@/components/QuestionnaireDetails/InformationEdit'
import Questions from '@/components/Questions'
import QuestionsEdit from '@/components/QuestionnaireDetails/QuestionsEdit'
import axios from '@/bin/axios'
import common from '@/util/common'
import EditNavBar from '@/components/Utils/EditNavBar.vue'

export default {
  name: 'QuestionnaireDetails',
  async created () {
    this.getInformation()
    this.getQuestions()
  },
  components: {
    'information-summary': InformationSummary,
    'information': Information,
    'information-edit': InformationEdit,
    'questions': Questions,
    'questions-edit': QuestionsEdit,
    'edit-nav-bar': EditNavBar
  },
  props: {
    traqId: {
      required: true
    }
  },
  data () {
    return {
      detailTabs: [ 'Information', 'Questions' ],
      selectedTab: 'Information',
      showEditButton: false,
      noTimeLimit: true,
      information: {},
      questions: [],
      newQuestionnaireId: undefined,
      removedQuestionIds: []
    }
  },
  methods: {
    alertNetworkError: common.alertNetworkError,
    getDateStr: common.customDateStr,
    getInformation () {
      // サーバーにアンケートの情報をリクエストする
      if (this.isNewQuestionnaire) {
        this.information = {
          title: '',
          description: '',
          res_shared_to: 'public',
          res_time_limit: this.newTimeLimit,
          respondents: [],
          administrators: [ this.traqId ],
          targets: [ 'traP' ]
        }
      } else {
        return axios
          .get('/questionnaires/' + this.questionnaireId)
          .then(res => {
            this.information = res.data
            if (this.administrates) {
              this.enableEditButton()
            } else {
              this.disableEditButton()
            }
            if (this.information.res_time_limit && this.information.res_time_limit !== 'NULL') {
              this.noTimeLimit = false
            }
          })
      }
    },
    getQuestions () {
      this.questions = []
      if (!this.isNewQuestionnaire) {
        axios
          .get('/questionnaires/' + this.questionnaireId + '/questions')
          .then(res => {
            this.questions = []
            res.data.forEach(data => {
              this.questions.push(common.convertDataToQuestion(data))
            })
          })
      }
    },
    submitQuestionnaire () {
      const informationData = {
        title: this.information.title,
        description: this.information.description,
        res_time_limit: this.noTimeLimit ? 'NULL' : new Date(this.information.res_time_limit).toLocaleString('ja-GB'),
        res_shared_to: this.information.res_shared_to,
        targets: this.information.targets,
        administrators: this.information.administrators
      }

      if (this.isNewQuestionnaire) {
        axios.post('/questionnaires', informationData)
          .then(resp => {
            // 返ってきたquestionnaireIDを保存
            this.newQuestionnaireId = resp.data.questionnaireID
          })
          .then(() => {
            // 質問をサーバーに送信
            return this.sendQuestions(0)
          })
          .then(() => {
            // 作成したアンケートの個別ページに遷移
            router.push('/questionnaires/' + this.newQuestionnaireId)
          })
          .catch(error => {
            // エラーが起きた場合は、送信済みのInformationを削除する
            console.log(error)
            axios.delete('/questionnaires/' + this.newQuestionnaireId)
            this.alertNetworkError()
          })
      } else {
        axios.patch('/questionnaires/' + this.questionnaireId, informationData)
          .then(() => {
            // 質問を送信
            return this.sendQuestions(0)
          })
          .then(() => {
            if (this.removedQuestionIds.length > 0) {
              // 削除された質問がある場合、それを送信
              return this.deleteRemovedQuestions(0)
            }
          })
          .then(this.getInformation) // 情報をアップデート
          .then(this.getQuestions) // 質問をアップデート
          .then(this.disableEditing) // 編集を終了
          .catch(error => {
            console.log(error)
            this.alertNetworkError()
          })
      }
    },
    sendQuestions (index) {
      // questions配列の、index番目以降の質問をサーバーに送信する
      const question = this.questions[ index ]
      const data = this.createQuestionData(index)

      if (this.isNewQuestion(question)) {
        return axios
          .post('/questions', data)
          .then(() => {
            if (index < this.questions.length - 1) {
              // 残りの質問を送信
              return this.sendQuestions(index + 1)
            }
          })
      } else {
        return axios
          .patch('/questions/' + question.questionId, data)
          .then(() => {
            if (index < this.questions.length - 1) {
              // 残りの質問を送信
              return this.sendQuestions(index + 1)
            }
          })
      }
    },
    deleteRemovedQuestions (index) {
      // removedQuestionIds配列の、index以降の質問について、DELETEリクエストを送る
      const id = this.removedQuestionIds[ index ]
      return axios
        .delete('/questions/' + id)
        .then(() => {
          if (index < this.removedQuestionIds.length - 1) {
            return this.deleteRemovedQuestions(index + 1)
          }
        })
    },
    deleteQuestionnaire () {
      if (this.isNewQuestionnaire) {
        router.push('/administrates')
      } else {
        axios
          .delete('/questionnaires/' + this.questionnaireId)
          .then(() => {
            router.push('/administrates')
            // アンケートを削除したら、Administratesページに戻る
          })
      }
    },
    createQuestionData (index) {
      // 与えられた質問1つ分のデータをサーバーに送るフォーマットのquestionDataにして返す
      const question = this.questions[ index ]
      let data = {
        questionnaireID: this.isNewQuestionnaire ? this.newQuestionnaireId : this.questionnaireId,
        question_type: question.type,
        question_num: index,
        page_num: question.pageNum,
        body: question.questionBody,
        is_required: question.isRequired,
        options: [],
        scale_label_left: '',
        scale_label_right: '',
        scale_min: 0,
        scale_max: 0
      }
      switch (question.type) {
        case 'Checkbox':
        case 'MultipleChoice':
          question.options.forEach(option => {
            data.options.push(option.label)
          })
          break
        case 'LinearScale':
          data.scale_label_left = question.scaleLabels.left
          data.scale_label_right = question.scaleLabels.right
          data.scale_min = question.scaleRange.left
          data.scale_max = question.scaleRange.right
          break
      }
      return data
    },
    isNewQuestion (question) {
      return question.questionId < 0
    },
    enableEditButton () {
      this.showEditButton = true
    },
    disableEditButton () {
      this.showEditButton = false
    },
    disableEditing () {
      this.isEditing = false
    },
    abortEditing () {
      if (this.isNewQuestionnaire) {
        router.push('/administrates')
      } else {
        this.getInformation()
          .then(this.getQuestions)
          .then(this.disableEditing)
      }
    },
    setData (name, data) {
      switch (name) {
        case 'questions':
          this.questions = data
          break
        case 'information':
          this.information = data
          break
        case 'noTimeLimit':
          this.noTimeLimit = data
          break
      }
    },
    setQuestionContent (index, label, value) {
      this.questions[ index ][ label ] = value
    },
    removeQuestion (index) {
      const id = this.questions[ index ].questionId
      if (id > 0) {
        // サーバーに存在する質問を削除した場合はリストに追加
        this.removedQuestionIds.push(id)
      }
      this.questions.splice(index, 1)
    }
  },
  computed: {
    administrates () {
      // 管理者かどうかを返す
      // getInformation() が完了する前は false を返す
      return this.information.administrators ? common.administrates(this.information.administrators, this.traqId) : false
    },
    questionnaireId () {
      if (this.isNewQuestionnaire) {
        return undefined
      } else {
        return Number(this.$route.params.id)
      }
    },
    isNewQuestionnaire () {
      return this.$route.params.id === 'new'
    },
    submitOk () {
      return common.noErrors(this.inputErrors)
    },
    isEditing: {
      get: function () {
        return this.isNewQuestionnaire || (this.$route.hash === '#edit' && this.administrates)
      },
      set: function (newBool) {
        // newBool : 閲覧 -> 編集
        // !newBool : 編集 -> 閲覧
        const newRoute = {
          name: 'QuestionnaireDetails',
          params: { id: this.questionnaireId },
          hash: newBool ? '#edit' : undefined
        }
        router.push(newRoute)
      }
    },
    currentTabComponent () {
      switch (this.selectedTab) {
        case 'Information': {
          if (this.isEditing) {
            return 'information-edit'
          } else {
            return 'information'
          }
        }
        case 'Questions': {
          if (this.isEditing) {
            return 'questions-edit'
          } else {
            return 'questions'
          }
        }
      }
    },
    title () {
      return this.information.title
    },
    editButtons () {
      return [
        {
          label: '送信',
          atClick: 'submit-questionnaire',
          disabled: !this.submitOk
        },
        {
          label: 'キャンセル',
          atClick: 'abort-editing',
          disabled: false
        }
      ]
    },
    informationProps () {
      return {
        details: this.information,
        administrates: this.administrates,
        deleteQuestionnaire: this.deleteQuestionnaire,
        questionnaireId: this.questionnaireId,
        noTimeLimit: this.noTimeLimit
      }
    },
    summaryProps () {
      let ret = {
        title: this.information.title
      }
      if (this.selectedTab === 'Information') {
        ret.description = this.information.description
        ret.timeLimit = this.getDateStr(this.information.res_time_limit)
      }
      return ret
    },
    newTimeLimit () {
      // 1週間後の23:59
      return moment().add(7, 'days').endOf('day').format().slice(0, -6)
    },
    inputErrors () {
      return {
        noTitle: {
          message: 'タイトルは入力必須です',
          isError: this.information.title === ''
        },
        noQuestions: {
          message: '質問がありません',
          isError: this.questions.length === 0
        }
      }
    }
  },
  watch: {
    $route: function (newRoute, oldRoute) {
      if (newRoute.params.id !== oldRoute.params.id) {
        this.showEditButton = false
        this.getInformation()
        this.getQuestions()
        this.newQuestionnaireId = undefined
      }
    },
    noTimeLimit: function (newBool, oldBool) {
      if (oldBool && !newBool && (this.information.res_time_limit === 'NULL' || this.information.res_time_limit === '')) {
        // 新しく回答期限を作ろうとしたとき
        this.information.res_time_limit = this.newTimeLimit
      }
    },
    information: function (newVal) {
      if (newVal.res_time_limit === '') {
        this.noTimeLimit = true
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
