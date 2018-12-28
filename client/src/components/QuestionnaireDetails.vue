<template>
  <div>
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
      :class="{'is-editing' : isEditing}"
      :name="currentTabComponent"
      @enable-edit-button="enableEditButton"
      @disable-editing="disableEditing"
    ></component>
  </div>
</template>

<script>

// import <componentname> from '<path to component file>'
import Information from '@/components/QuestionnaireDetails/Information'
import Questions from '@/components/QuestionnaireDetails/Questions'

export default {
  name: 'QuestionnaireDetails',
  components: {
    'information': Information,
    'questions': Questions
  },
  props: {
    traqId: {
      type: String,
      required: true
    },
    getDateStr: {
      type: Function,
      required: true
    }
  },
  data () {
    return {
      detailTabs: [ 'Information', 'Questions' ],
      selectedTab: 'Information',
      isEditing: false,
      showEditButton: false
    }
  },
  methods: {
    enableEditButton () {
      this.showEditButton = true
    },
    disableEditing () {
      console.log('QuestionnaireDetails.vue disableEditing()')
      this.isEditing = false
    }
  },
  computed: {
    detailTabsProps () {
      return {
        traqId: this.traqId,
        getDateStr: this.getDateStr,
        isEditing: this.isEditing
      }
    },
    currentTabComponent () {
      // return 'product-reviews'
      return this.selectedTab.replace(/([a-z])([A-Z])/g, '$1-$2').toLowerCase()
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
</style>
