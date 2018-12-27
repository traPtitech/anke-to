<template>
  <div id="app" class="is-fullheight">
    <top-navbar
      @toggle-side-menu="toggleSideMenu"
      @close-side-menu="closeSideMenu"
      :isSideMenuActive="isSideMenuActive"
      :traqId="traqId"
    ></top-navbar>
    <div class="columns is-fullheight">
      <side-menu class="fixed-sidemenu desktop"></side-menu>
      <side-menu class="sidemenu" v-show="isSideMenuActive"></side-menu>
      <div class="column app-main" @click="closeSideMenu">
        <router-view :traqId="traqId" :getDateStr="getDateStr"></router-view>
      </div>
    </div>
  </div>
</template>

<script>
import axios from '@/bin/axios'
import TopNavbar from './components/Utils/TopNavbar.vue'
import SideMenu from './components/Utils/SideMenu.vue'
export default {
  name: 'App',
  components: {
    'top-navbar': TopNavbar,
    'side-menu': SideMenu
  },
  async created () {
    const resp = await axios.get('/users/me')
    this.user = resp.data
  },
  data () {
    return {
      isSideMenuActive: false,
      user: {}
    }
  },
  computed: {
    traqId () {
      return String(this.user.traqID)
    }
  },
  methods: {
    toggleSideMenu () {
      this.isSideMenuActive = !this.isSideMenuActive
    },
    closeSideMenu () {
      this.isSideMenuActive = false
    },
    getDateStr (str) {
      return str === 'NULL' ? '-' : new Date(str).toLocaleString()
    }
  }
}
</script>

<style lang="scss">
@import "../node_modules/bulma/bulma.sass";
@import url("../static/css/themify-icons.css");
// @import "../node_modules/themify-icons/themify-icons/_themify-icons.scss";
#app {
  font-family: "Avenir", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  // text-align: center;
  color: #2c3e50;
}
.columns {
  padding-top: 0.75rem;
}

.column {
  padding: 1rem;
}

.is-fullheight {
  height: 100%;
}

@media screen and (max-width: 768px) {
  // mobile
  .sidemenu {
    height: fit-content;
  }
  .app-main {
    height: 100%;
  }
}

@media screen and (max-width: 1088px) {
  // 固定サイドメニューは非表示
  .fixed-sidemenu {
    display: none;
  }
}

@media screen and (min-width: 1088px) {
  // widthが大きいときは固定サイドメニューのみ表示
  .sidemenu {
    display: none;
  }
}

html {
  height: 100%;
}

body {
  height: 100%;
}

.app-main {
  padding: 0;
}

.columns {
  margin: 0;
  padding-top: 0;
}

.card {
  /* width: fit-content; */
  margin: 1rem 1.5rem;
  overflow-x: auto;
  width: auto;
  max-width: fit-content;
}
.card-header-title {
  color: #707880;
  font-weight: 400;
  padding: 1rem 1.5rem;
}
.card-content {
  padding: 1rem;
}
</style>
