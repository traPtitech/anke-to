<template>
  <div>
    <nav class="navbar has-background-white-bis" role="navigation" aria-label="main navigation">
      <!-- <div class="container"> -->
      <div class="navbar-brand">
        <a
          role="button"
          class="navbar-burger"
          @click="toggleMenu"
          :class="{'is-active' : isSideMenuActive}"
        >
          <span></span>
          <span></span>
          <span></span>
        </a>
        <a class="navbar-item" href="/">
          <p class="title is-3">anke-to</p>
        </a>
      </div>
      <!-- </div> -->
    </nav>
  </div>
</template>
<script>
// import axios from '@/bin/axios'

export default {
  name: 'TopNavbar',
  props: {
    isSideMenuActive: {
      type: Boolean,
      required: true
    }
  },
  methods: {
    toggleMenu () {
      this.$emit('toggle-side-menu')
      // this.isMenuActive = !this.isMenuActive
    },
    closeMenu () {
      this.$emit('close-side-menu')
      // this.isMenuActive = false
    }
  }
}
</script>

<style scoped>
.title {
  margin-bottom: 0;
}
.navbar-burger {
  margin-left: 0;
}
</style>
