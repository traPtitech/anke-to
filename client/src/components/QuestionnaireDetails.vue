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
    <div :class="{'is-editing' : isEditing, 'has-navbar-fixed-bottom': isEditing}">
      <component
        :is="currentTabComponent"
        :traqId="traqId"
        class="details-child is-fullheight"
        :name="currentTabComponent"
        :editMode="isEditing ? 'question' : undefined"
        :informationProps="informationProps"
        :questionsProps="questions"
        :title="title"
        @set-data="setData"
        @set-question-content="setQuestionContent"
      ></component>
      <edit-nav-bar
        v-if="isEditing"
        :editButtons="editButtons"
        @submit-questionnaire="submitQuestionnaire"
        @disable-editing="disableEditing"
      ></edit-nav-bar>
    </div>
  </div>
</template>

<script>

import moment from 'moment'
import router from '@/router'
import Information from '@/components/QuestionnaireDetails/Information'
import InformationEdit from '@/components/QuestionnaireDetails/InformationEdit'
import Questions from '@/components/Questions'
import QuestionsEdit from '@/components/QuestionnaireDetails/QuestionsEdit'
import axios from '@/bin/axios'
import common from '@/util/common'

export default {
  name: 'QuestionnaireDetails',
  async created () {
    this.getQuestions()
  },
  components: {
    'information': Information,
    'information-edit': InformationEdit,
    'questions': Questions,
    'questions-edit': QuestionsEdit
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
      questions: []
      noTimeLimit: true,
    }
  },
  methods: {
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
          targets: [ this.traqId ]
        }
      } else {
        axios
          .get('/questionnaires/' + this.questionnaireId)
          .then(res => {
            this.information = res.data
            if (this.administrates) {
              this.enableEditButton()
            }
            if (this.information.res_time_limit && this.information.res_time_limit !== 'NULL') {
              this.noTimeLimit = false
            }
          })
      }
    },
    getQuestions () {
      axios
        .get('/questionnaires/' + this.questionnaireId + '/questions')
        .then(res => {
          this.questions = []
          res.data.forEach(data => {
            this.questions.push(common.convertDataToQuestion(data))
          })
        })
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
    },
    enableEditButton () {
      this.showEditButton = true
    },
    disableEditing () {
      this.isEditing = false
    },
    setQuestions (questions) {
      this.questions = questions
    },
    setQuestionContent (index, label, value) {
      this.questions[ index ][ label ] = value
    }
  },
  computed: {
    questionnaireId () {
      return this.isNewQuestionnaire ? '' : this.$route.params.id
    },
    isNewQuestionnaire () {
      return this.$route.params.id === 'new'
    },
    administrates () {
      // 管理者かどうかを返す
      if (this.information.administrators) {
        for (let i = 0; i < this.information.administrators.length; i++) {
          if (this.traqId === this.information.administrators[ i ]) {
            return true
          }
        }
      }
      return false
    },
    submitOk () {
      // 送信できるかどうかを返す
      return this.information.title !== '' && this.information.administrators && this.information.administrators.length > 0 &&
        this.questions.length > 0
    },
    isEditing: {
      get: function () {
        if (this.isNewQuestionnaire || this.$route.hash === '#edit') {
          return true
        }
        return false
      },
      set: function (newBool) {
        if (newBool) {
          // 閲覧 -> 編集
          router.push('/questionnaires/' + this.questionnaireId + '#edit')
        } else {
          // 編集 -> 閲覧
          router.push('/questionnaires/' + this.questionnaireId)
        }
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
          atClick: 'disable-editing',
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
    newTimeLimit () {
      // 1週間後の日時
      return moment().add(7, 'days').format().slice(0, -6)
    }
  },
  watch: {
    $route: function (newRoute, oldRoute) {
      if (newRoute.params.id !== oldRoute.params.id) {
        this.getInformation()
        this.getQuestions()
        this.newQuestionnaireId = undefined
      }
    },
    noTimeLimit: function (newBool, oldBool) {
      if (oldBool && !newBool && (this.information.res_time_limit === 'NULL' || this.information.res_time_limit === '')) {
        // 新しく回答期限を作ろうとすると、1週間後の日時が設定される
        this.information.res_time_limit = this.newTimeLimit
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
