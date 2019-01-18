<template>
  <nav class="navbar is-fixed-bottom">
    <article class="is-2 is-flex">
      <div class="editor-buttons is-pulled-right">
        <button
          v-for="editButton in editButtons"
          :key="editButton.label"
          class="button is-medium"
          :disabled="editButton.disabled"
          @click.prevent="atClick(editButton.atClick)"
        >{{ editButton.label }}</button>
      </div>
    </article>
  </nav>
</template>

<script>

// import <componentname> from '<path to component file>'

export default {
  name: 'EditNavBar',
  components: {
  },
  props: {
    editButtons: {
      type: Array,
      required: true
    }
  },
  data () {
    return {
    }
  },
  methods: {
    atClick (func) {
      this.$emit(func)
    }
  },
  computed: {
  },
  mounted () {
  }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style lang="scss" scoped>
.editor-buttons {
  margin: auto;
  .button {
    margin: 1rem;
    // margin-bottom: 1rem;
    width: 8rem;
    max-width: 100%;
  }
}
.navbar.is-fixed-bottom {
  background-color: gray;
}
</style>
