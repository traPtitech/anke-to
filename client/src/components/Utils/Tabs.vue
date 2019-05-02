<template>
  <div>
    <div class="tabs is-centered">
      <ul>
        <li
          class="tab"
          :class="{ 'is-active': selectedTab === tab }"
          v-for="(tab, index) in tabs"
          :key="index"
          @click="selectedTab = tab"
        >
          <a>{{ tab }}</a>
        </li>
      </ul>
    </div>
    <component :is="currentTabComponent" :props="props"></component>
  </div>
</template>

<script>

// import <componentname> from '<path to component file>'
import Information from '@/components/Information/Information'
import Questions from '@/components/Questions/Questions'

export default {
  name: 'Tabs',
  components: {
    'information': Information,
    'questions': Questions
  },
  props: {
    tabs: {
      type: Array,
      required: true
    },
    props: {
      type: Object,
      required: true
    }
  },
  data () {
    return {
      selectedTab: this.tabs[ 0 ]
    }
  },
  methods: {
  },
  computed: {
    currentTabComponent: function () {
      // return 'product-reviews'
      return this.selectedTab.replace(/([a-z])([A-Z])/g, '$1-$2').toLowerCase()
    }
  },
  mounted () {
  }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style scoped>
.tabs:not(:last-child) {
  margin-bottom: 1rem;
}
.tabs:first-child {
  margin-top: 1rem;
}
</style>
