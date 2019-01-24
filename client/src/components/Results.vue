<template>
  <div>
    <div v-if="canViewResults" class="details is-fullheight">
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
      </div>
      <component
        :is="currentTabComponent"
        class="details-child is-fullheight"
        :name="currentTabComponent"
        :results="results"
        :questions="questions"
      ></component>
    </div>

    <div v-if="this.information.administrators && !canViewResults" class="message is-danger">
      <p class="message-body error-message">結果を閲覧する権限がありません</p>
    </div>
  </div>
</template>

<script>

// import <componentname> from '<path to component file>'
import axios from '@/bin/axios'
import common from '@/util/common'
import Individual from '@/components/Results/Individual'
import Spreadsheet from '@/components/Results/Spreadsheet'

export default {
  name: 'Results',
  components: {
    individual: Individual,
    spreadsheet: Spreadsheet
  },
  async created () {
    this.getInformation()
      .then(this.getMyResponses)
      .then(() => {
        if (this.canViewResults) {
          this.getResults()
          this.getQuestions()
        }
      })
  },
  props: {
    traqId: {
      required: true
    }
  },
  data () {
    return {
      results: [],
      questions: [],
      information: {},
      hasResponded: false,
      detailTabs: [ 'Spreadsheet', 'Individual' ],
      selectedTab: 'Spreadsheet'
    }
  },
  methods: {
    getResults () {
      return axios
        .get('/results/' + this.questionnaireId)
        .then(res => {
          this.results = res.data
        })
    },
    getQuestions () {
      this.questions = []
      return axios
        .get('/questionnaires/' + this.questionnaireId + '/questions')
        .then(res => {
          for (const question of res.data) {
            this.questions.push(question.body)
          }
        })
    },
    getInformation () {
      return axios
        .get('/questionnaires/' + this.questionnaireId)
        .then(res => {
          this.information = res.data
        })
    },
    getMyResponses () {
      return axios
        .get('/users/me/responses/' + this.questionnaireId)
        .then(res => {
          if (res.data.length > 0) {
            this.hasResponded = true
          }
        })
    }
  },
  computed: {
    questionnaireId () {
      return this.$route.params.id
    },
    administrates () {
      if (!this.information.administrators) {
        return undefined
      }
      return common.administrates(this.information.administrators, this.traqId)
    },
    canViewResults () {
      return common.canViewResults(this.information, this.administrates, this.hasResponded)
    },
    currentTabComponent () {
      switch (this.selectedTab) {
        case 'Spreadsheet':
          return 'spreadsheet'
        case 'Individual':
          return 'individual'
      }
    }
  },
  mounted () {
  }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style lang="scss" scoped>
.error-message {
  font-size: 1rem;
  margin: 1rem;
}
</style>
