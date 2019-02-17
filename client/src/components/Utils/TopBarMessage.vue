<template>
  <div
    v-if="message.showMessage"
    class="message"
    :class="{'is-danger': message.color==='red', 'is-success': message.color==='green'}"
  >
    <p class="message-body">{{ message.body }}</p>
  </div>
</template>

<script>

export default {
  name: 'top-bar-message',
  components: {
  },
  props: {
    message: {
      type: Object,
      required: true
    }
  },
  data () {
    return {
    }
  },
  methods: {
  },
  computed: {
  },
  mounted () {
  }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style lang="scss" scoped>
.message-body {
  font-size: 1rem;
  margin: 1rem;
}
</style>
