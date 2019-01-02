<template>
  <div class="questionnaire-details is-fullheight">
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
        v-show="showEditButton"
      >
        <span class="ti-pencil"></span>
      </a>
    </div>
    <component
      :is="currentTabComponent"
      :props="detailTabsProps"
      :class="{'is-editing' : isEditing, 'has-navbar-fixed-bottom': isEditing}"
      class="questionnaire-details-child is-fullheight"
      :name="currentTabComponent"
      @enable-edit-button="enableEditButton"
      @disable-editing="disableEditing"
    ></component>
  </div>
</template>

<script>

// import <componentname> from '<path to component file>'
import router from '@/router'
import Information from '@/components/QuestionnaireDetails/Information'
import InformationEdit from '@/components/QuestionnaireDetails/InformationEdit'
import Questions from '@/components/QuestionnaireDetails/Questions'
import QuestionsEdit from '@/components/QuestionnaireDetails/QuestionsEdit'

export default {
  name: 'QuestionnaireDetails',
  components: {
    'information': Information,
    'information-edit': InformationEdit,
    'questions': Questions,
    'questions-edit': QuestionsEdit
  },
  props: {
    traqId: {
      type: String,
      required: true
    }
  },
  data () {
    return {
      detailTabs: [ 'Information', 'Questions' ],
      selectedTab: 'Information',
      showEditButton: false
    }
  },
  methods: {
    enableEditButton () {
      this.showEditButton = true
    },
    disableEditing () {
      this.isEditing = false
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
    detailTabsProps () {
      return {
        traqId: this.traqId
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
.tabs {
  margin-bottom: 0;
  margin-right: 0.5rem;
  margin-left: 0.5rem;
}
.tabs:first-child {
  margin-top: 1rem;
}
.is-editing {
  background-color: #c2c2c2;
}
#edit-button {
  border: #dbdbdb solid 1px;
}
.is-fullheight {
  min-height: fit-content;
}
.has-navbar-fixed-bottom {
  padding-bottom: 100px;
}
.questionnaire-details /deep/ .questionnaire-details-child {
  // 子コンポーネント Information, InformationEdit, Questions, QuestionsEdit に適用される
  pre {
    white-space: pre-line;
    font-size: inherit;
    -webkit-font-smoothing: inherit;
    font-family: inherit;
    line-height: inherit;
    background-color: inherit;
    color: inherit;
    padding: 0.625em;
  }
  article.column {
    padding: 0;
  }
  .columns {
    margin-bottom: 0;
  }
  .columns:first-child {
    display: flex;
  }
  .card {
    max-width: 100%;
    padding: 0.7rem;
  }
  .card-content {
    .subtitle {
      margin: 0;
    }
    details {
      margin: 0.5rem;
      p {
        padding: 0 0.5rem;
      }
    }
  }
  .navbar.is-fixed-bottom {
    background-color: gray;
  }
  @media screen and (min-width: 769px) {
    // widthが大きいときは横並びのカードの間を狭くする
    .column:not(:last-child) > .card {
      margin-right: 0;
    }
  }
}
</style>
