<template type="functional">
  <svg
    version="1.1"
    xmlns="http://www.w3.org/2000/svg"
    xmlns:xlink="http://www.w3.org/1999/xlink"
    viewBox="0 0 17 17"
    :width="size"
    :height="size"
  >
    <g :fill="color">
      <path
        d="M17,8.8c-1.3-3.3-4.7-5.5-8.5-5.5S1.3,5.5,0,8.8c0,0.1,0,0.2,0,0.4c1.3,3.3,4.7,5.5,8.5,5.5s7.2-2.2,8.5-5.5
		C17,9.1,17,8.9,17,8.8z M8.5,13.6C5.2,13.6,2.2,11.8,1,9c1.2-2.8,4.2-4.6,7.5-4.6S14.8,6.2,16,9C14.8,11.8,11.8,13.6,8.5,13.6z"
      />
      <path
        d="M8.5,5.4C6.5,5.4,4.9,7,4.9,9c0,2,1.6,3.6,3.6,3.6c2,0,3.6-1.6,3.6-3.6C12.1,7,10.5,5.4,8.5,5.4z"
      />
    </g>
  </svg>
</template>

<script>
export default {
  name: 'IconEye',
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
