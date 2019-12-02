<template>
  <!-- eslint-disable-next-line vue/require-component-is -->
  <component
    :is="`icon-${name}`"
    :class="`icon-${name}`"
    :color="color"
    :size="size"
  />
</template>

<script>
import IconEye from './IconEye'
import IconEyeClosed from './IconEyeClosed'

export default {
  name: 'Icon',
  components: {
    IconEye,
    IconEyeClosed
  },
  props: {
    name: {
      type: String,
      required: true
    },
    size: {
      type: Number,
      default: 17
    },
    color: {
      type: String,
      default: 'black'
    }
  }
}
</script>

<style lang="scss" scoped></style>
