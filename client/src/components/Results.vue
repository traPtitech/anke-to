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
    </div>
    <component
      :is="currentTabComponent"
      class="details-child is-fullheight"
      :name="currentTabComponent"
      :results="results"
      :questions="questions"
    ></component>
  </div>
</template>

<script>

// import <componentname> from '<path to component file>'
import axios from '@/bin/axios'
import Individual from '@/components/Results/Individual'
import Spreadsheet from '@/components/Results/Spreadsheet'

export default {
  name: 'Results',
  components: {
    individual: Individual,
    spreadsheet: Spreadsheet
  },
  async created () {
    axios
      .get('/results/' + this.questionnaireId)
      .then(res => {
        this.results = res.data
      })
    axios
      .get('/questionnaires/' + this.questionnaireId + '/questions')
      .then(res => {
        for(const question of res.data) {
          this.questions.push(question.body)
        }
      })
  },
  props: {
  },
  data () {
    return {
      results: [],
      questions: [],
      detailTabs: [ 'Spreadsheet', 'Individual' ],
      selectedTab: 'Spreadsheet'
    }
  },
  methods: {
  },
  computed: {
    questionnaireId () {
      return this.$route.params.id
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
<style scoped>
</style>
