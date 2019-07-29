<template>
  <div>
    <div>
      <span class="has-text-weight-bold">{{ userList.summary }} ({{ userList.list.length }})</span>
      <a>
        <span class="ti-pencil" v-if="userList.editable" @click="changeActiveModal(userList)"></span>
      </a>
    </div>
    <p class="has-text-grey user-list">
      <span v-for="(user, index) in userList.list" :key="index">
        <span :class="{'highlight-name': user === 'traP' || user === getMyTraqId}">{{ user }}</span>
        <span>
          {{
          index === userList.list.length - 1 ? "" : ", "
          }}
        </span>
      </span>
    </p>
  </div>
</template>

<script>

export default {
  name: 'UserList',
  components: {
  },
  props: {
    userList: {
      type: Array
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
  watch: {
  },
  mounted () {
  }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style lang="scss" scoped>
.user-list {
  margin: 0 0.5rem;
}
</style>
