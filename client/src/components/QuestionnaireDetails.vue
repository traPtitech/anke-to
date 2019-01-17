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
    <component
      :is="currentTabComponent"
      :traqId="traqId"
      :class="{'is-editing' : isEditing, 'has-navbar-fixed-bottom': isEditing}"
      class="details-child is-fullheight"
      :name="currentTabComponent"
      :editMode="isEditing ? 'question' : undefined"
      :questionsProps="questions"
      @enable-edit-button="enableEditButton"
      @disable-editing="disableEditing"
      @set-questions="setQuestions"
      @set-question-content="setQuestionContent"
    ></component>
  </div>
</template>

<script>

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
    }
  },
  methods: {
    getQuestions () {
      axios
        .get('/questionnaires/' + this.questionnaireId + '/questions')
        .then(res => {
          this.questions = []
          res.data.forEach(data => {
            this.questions.push(common.convertDataToQuestion(data))
          })
        })
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
    }
  },
  mounted () {
  }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style lang="scss" scoped>
</style>
