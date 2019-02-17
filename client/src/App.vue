<template>
  <div id="app" class="is-fullheight">
    <top-navbar
      @toggle-side-menu="toggleSideMenu"
      @close-side-menu="closeSideMenu"
      :isSideMenuActive="isSideMenuActive"
      :traqId="traqId"
    ></top-navbar>
    <div class="columns is-fullheight">
      <side-menu class="fixed-sidemenu desktop" :traqId="traqId"></side-menu>
      <side-menu
        class="sidemenu"
        v-show="isSideMenuActive"
        :traqId="traqId"
        @close-side-menu="closeSideMenu"
      ></side-menu>
      <div class="column app-main" @click="closeSideMenu">
        <router-view :traqId="traqId"></router-view>
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
      return this.user.traqID
    }
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
  max-width: 100%;
}

.is-fullheight {
  // height: 100%;
  min-height: -webkit-fill-available;
}

.button.is-disabled {
  background-color: white;
  border-color: #dbdbdb;
  -webkit-box-shadow: none;
  box-shadow: none;
  opacity: 0.5;
  pointer-events: none; // ポインター操作を無効化
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

summary {
  cursor: pointer;
}

a[disabled] {
  cursor: default;
  pointer-events: none;
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

.has-navbar-fixed-bottom {
  padding-bottom: 100px;
}

.is-editing {
  background-color: #c2c2c2;
}
.details {
  .tabs {
    margin-bottom: 0;
    margin-right: 0.5rem;
    margin-left: 0.5rem;
  }
  .tabs:first-child {
    margin-top: 1rem;
  }
  #edit-button {
    border: #dbdbdb solid 1px;
  }
  .checkbox {
    width: 4rem;
    margin: 0.5rem;
  }
  .is-fullheight {
    min-height: fit-content;
  }
  .details-child {
    article.column {
      padding: 0;
    }
    .columns {
      margin-bottom: 0;
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
    @media screen and (min-width: 769px) {
      // widthが大きいときは横並びのカードの間を狭くする
      .column:not(:last-child) > .card {
        margin-right: 0;
      }
    }
  }
}

.icon.circled {
  background-color: lightgray;
  border-radius: 1rem;
}

// readonly buttons
.readonly-checkbox,
.readonly-radiobutton {
  width: 0.8rem;
  height: 0.8rem;
  border: grey solid 1px;
  display: inline-block;
  margin: auto;
  &.checked {
    background-color: darkgray;
  }
}
.readonly-checkbox {
  border-radius: 0.1rem;
}
.readonly-radiobutton {
  border-radius: 0.5rem;
}

// list animation
.list-move {
  transition: transform 1s;
}
.list-leave-to,
.list-enter {
  transition: all 1s;
  opacity: 0;
  transform: translateX(30px);
}
.list-leave-active {
  position: absolute;
}

// sort handles and trash buttons
.sort-handle,
.delete-button {
  width: min-content;
  margin: 0 auto;
  .icon {
    cursor: pointer;
    margin: auto;
  }
  .icon.disabled {
    color: lightgray;
    pointer-events: none;
  }
  .icon.is-medium {
    font-size: 1.2rem;
  }
}

// questions
.questions {
  .question {
    .columns {
      width: 100%;
    }
    .column {
      padding: 0.2rem;
    }
    .column.question {
      padding: 0.5rem 1rem;
    }
    .column.left-bar {
      width: fit-content;
      max-width: fit-content;
      padding: 0;
    }
    input,
    textarea {
      border: none;
      outline: none;
      -webkit-box-shadow: none;
      box-shadow: none;
      border-radius: 0;
    }
    .question-body {
      margin-bottom: 0.5rem;
    }
    .response-body {
      margin: 0.5rem 0.5rem;
      p {
        margin-top: 1rem;
      }
      p.has-underline {
        padding-bottom: 0.25rem;
      }
    }
    .has-underline {
      padding: 0 0.5rem;
      border-bottom: grey dotted 0.5px;
    }
    .is-editable.has-underline {
      border-bottom: lightgrey solid 0.5px;
    }
    input:focus.has-underline {
      border-bottom: black solid 2px;
    }
  }
}
</style>
