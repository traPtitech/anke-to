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
        <p class="subtitle pull-right">{{ traqId }}</p>
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
      type: String,
      required: true
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
  height: min-content;
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
</style>
