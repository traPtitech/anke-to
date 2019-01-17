<template>
  <div>
    <nav class="navbar has-background-white-bis" role="navigation" aria-label="main navigation">
      <div class="navbar-brand">
        <a
          role="button"
          class="navbar-burger"
          @click="$emit('toggle-side-menu')"
          :class="{'is-active' : isSideMenuActive}"
        >
          <span></span>
          <span></span>
          <span></span>
        </a>
        <div @click="$emit('close-side-menu')">
          <router-link class="navbar-item" to="/targeted">
            <p class="title is-3">anke-to</p>
          </router-link>
        </div>
        <div class="subtitle pull-right">
          <img
            class="user-icon"
            v-if="traqId"
            v-bind:src="'https://q.trapti.tech/static/icon/' + traqId + '/64.png'"
          >
          <p class="user-name">{{traqId}}</p>
        </div>
      </div>
    </nav>
  </div>
</template>
<script>

export default {
  name: 'TopNavbar',
  props: {
    isSideMenuActive: {
      type: Boolean,
      required: true
    },
    traqId: {
      required: false
    }
  },
  methods: {
  }
}
</script>

<style scoped>
.title {
  margin-bottom: 0;
}
.subtitle {
  min-width: 100px;
  position: absolute;
  top: 0;
  bottom: 0;
  right: 0;
  margin: auto;
  padding: 0 1rem;
}
.navbar-burger {
  margin-left: 0;
}
.user-icon {
  width: 36px;
  height: 36px;
  border-radius: 18px;
  display: inline-block;
  vertical-align: middle;
}
.user-name {
  display: inline;
  vertical-align: middle;
}
</style>
