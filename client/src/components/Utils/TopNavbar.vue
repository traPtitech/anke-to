<template>
  <div>
    <nav class="navbar" role="navigation" aria-label="main navigation">
      <div class="navbar-brand">
        <a
          role="button"
          class="navbar-burger"
          @click="$emit('toggle-side-menu')"
          :class="{ 'is-active': isSideMenuActive }"
        >
          <span></span>
          <span></span>
          <span></span>
        </a>
        <div @click="$emit('close-side-menu')" class="is-flex">
          <router-link class="navbar-item" to="/targeted">
            <img class="title" src="/static/img/logo.png" />
          </router-link>
        </div>
        <div class="subtitle pull-right user">
          <img
            class="user-icon"
            v-if="traqId"
            v-bind:src="'https://q.trap.jp/api/1.0/public/icon/' + traqId"
          />
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

<style lang="scss" scoped>
.navbar {
  background-color: $base-gray;
}
.title {
  margin-bottom: 0;
}
.subtitle {
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
.user {
  height: fit-content;
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
