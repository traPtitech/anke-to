<template>
  <div class="wrapper checkbox icon-button" :class="{ checked: checked }">
    <span
      class="ti-check"
      :class="{ checked: checked }"
      @click="emitInput"
    ></span>
  </div>
</template>

<script>
export default {
  name: 'Checkbox',
  components: {},
  props: {
    checked: {
      type: Boolean,
      required: true
    }
  },
  methods: {
    emitInput() {
      this.$emit('input', !this.checked)
    }
  }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style lang="scss" scoped>
.checkbox.icon-button {
  width: 1.5rem;
  height: 1.5rem;
  padding: 0;
  display: flex;
  color: $base-gray;
  background-color: white;
  border: 2px solid $base-darkbrown;
  border-radius: 15%;
  &:hover {
    cursor: pointer;
    color: $base-gray;
  }
  &.checked {
    color: white;
    background-color: $base-darkbrown;
    border: none;
  }
}
span[class^='ti-'] {
  line-height: normal;
  font-size: small;
  font-weight: bolder;
  margin: 0 auto;
  margin-top: calc(0.4rem - 2px); // borderの分を引いた
  &.checked {
    margin-top: 0.4rem;
  }
}
</style>
