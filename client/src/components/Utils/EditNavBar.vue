<template>
  <nav class="navbar is-fixed-bottom is-flex">
    <div class="editor-buttons">
      <button
        v-for="editButton in editButtons"
        :key="editButton.label"
        class="button is-medium"
        :class="editButton.class"
        :disabled="editButton.disabled"
        @click.prevent="atClick(editButton.atClick)"
      >{{ editButton.label }}</button>
    </div>
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
@import "@/css/variables.scss";
.editor-buttons {
  margin: auto;
  .button {
    margin: 1rem;
    width: 8rem;
    max-width: 100%;
    &:hover {
      border-color: $base-pink;
    }
  }
  @media screen and (max-width: 768px) {
    .button {
      width: fit-content;
    }
  }
}
.navbar.is-fixed-bottom {
  background-color: $var-indigo;
}
</style>
