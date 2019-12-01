<template type="functional">
  <svg
    xmlns="http://www.w3.org/2000/svg"
    xmlns:xlink="http://www.w3.org/1999/xlink"
    viewBox="0 0 17.3 17.3"
    :width="size"
    :height="size"
  >
    <path
      d="M16.5,12.7l0.7-0.7l-1.4-1.4C16.7,9.8,17,9.1,17.1,9l-0.9-0.4c0,0-1.8,3.6-7.5,3.6c0,0,0,0,0,0
	c-5.7,0-7.4-3.4-7.5-3.6L0.3,9c0,0.1,0.4,0.7,1.1,1.5l-1.4,1.4l0.7,0.7l1.4-1.5c0.6,0.4,1.2,0.8,2.1,1.2l-0.7,1.8l0.9,0.4l0.7-1.8
	c0.8,0.2,1.8,0.4,2.9,0.4v1.8h1v-1.8c1.1,0,2.1-0.2,2.9-0.4l0.7,1.8l0.9-0.4l-0.7-1.8c0.8-0.3,1.5-0.7,2-1.1L16.5,12.7z"
      :fill="color"
    />
  </svg>
</template>

<script>
export default {
  name: 'IconEyeClosed',
  props: {
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
