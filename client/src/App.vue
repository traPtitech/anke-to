<template>
  <div id="app" class="is-fullheight">
    <top-navbar
      :is-side-menu-active="isSideMenuActive"
      @toggle-side-menu="toggleSideMenu"
      @close-side-menu="closeSideMenu"
    ></top-navbar>
    <div class="columns">
      <side-menu class="fixed-sidemenu desktop"></side-menu>
      <side-menu
        v-show="isSideMenuActive"
        class="sidemenu"
        @close-side-menu="closeSideMenu"
      ></side-menu>
      <div class="column app-main" @click="closeSideMenu">
        <router-view></router-view>
      </div>
    </div>
  </div>
</template>

<script>
import TopNavbar from './components/Utils/TopNavbar.vue'
import SideMenu from './components/Utils/SideMenu.vue'

export default {
  name: 'App',
  components: {
    'top-navbar': TopNavbar,
    'side-menu': SideMenu
  },
  data() {
    return {
      isSideMenuActive: false,
      user: {}
    }
  },
  computed: {},
  async created() {},
  methods: {
    toggleSideMenu() {
      this.isSideMenuActive = !this.isSideMenuActive
    },
    closeSideMenu() {
      this.isSideMenuActive = false
    }
  }
}
</script>

<style lang="scss">
@import '~bulma/bulma.sass';
@import '@/style/global.scss';
@import '@/style/colors.scss';
@import 'static/css/themify-icons.css';

.app-main {
  min-width: 0;
}
</style>
