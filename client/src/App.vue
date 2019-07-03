<template>
  <div id="app" class="is-fullheight">
    <top-navbar
      @toggle-side-menu="toggleSideMenu"
      @close-side-menu="closeSideMenu"
      :isSideMenuActive="isSideMenuActive"
    ></top-navbar>
    <div class="columns is-fullheight">
      <side-menu class="fixed-sidemenu desktop"></side-menu>
      <side-menu class="sidemenu" v-show="isSideMenuActive" @close-side-menu="closeSideMenu"></side-menu>
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
  async created () {
  },
  data () {
    return {
      isSideMenuActive: false,
      user: {}
    }
  },
  computed: {
  },
  methods: {
    toggleSideMenu () {
      this.isSideMenuActive = !this.isSideMenuActive
    },
    closeSideMenu () {
      this.isSideMenuActive = false
    }
  }
}
</script>

<style lang="scss">
@import "../node_modules/bulma/bulma.sass";
@import "@/style/global.scss";
@import "@/style/colors.scss";
@import "static/css/themify-icons.css";
</style>
