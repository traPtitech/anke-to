<template>
  <div v-if="userList">
    <div>
      <span class="has-text-weight-bold"
        >{{ userList.summary }} ({{ userList.list.length }})</span
      >
      <a class="edit-list-icon">
        <span
          v-if="userList.editable"
          class="ti-pencil"
          @click="changeActiveModal(userList)"
        ></span>
      </a>
    </div>
    <p class="has-text-grey user-list">
      <span v-for="(user, index) in userList.list" :key="index">
        <span
          :class="{
            'highlight-name': user === 'traP' || user === $store.state.me.traqId
          }"
          >{{ user }}</span
        >
        <span>
          {{ index === userList.list.length - 1 ? '' : ', ' }}
        </span>
      </span>
    </p>
  </div>
</template>

<script>
export default {
  name: 'UserList',
  components: {},
  props: {
    userList: {
      type: Object,
      default: undefined
    }
  },
  methods: {
    changeActiveModal(list) {
      this.$emit('change-active-modal', list)
    }
  }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style lang="scss" scoped>
.user-list {
  margin: 0 0.5rem;
}
.edit-list-icon {
  margin-left: 0.2rem;
}
</style>
