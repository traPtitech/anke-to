<template>
  <div class="is-fullheight">
    <div class="tabs is-centered">
      <ul></ul>
      <a
        id="edit-button"
        :class="{'is-editing': isEditing}"
        @click.prevent="isEditing = !isEditing"
      >
        <span class="ti-pencil"></span>
      </a>
    </div>
    <div :class="{'is-editing' : isEditing}" class="is-fullheight">
      <p>responseId = {{ responseId }} の回答詳細画面</p>
    </div>
  </div>
</template>

<script>

// import <componentname> from '<path to component file>'
import router from '@/router'

export default {
  name: 'ResponseDetails',
  components: {
  },
  props: {
  },
  data () {
    return {
    }
  },
  methods: {
  },
  computed: {
    responseId () {
      return this.$route.params.id
    },
    isEditing: {
      get: function () {
        if (this.$route.hash === '#edit') {
          return true
        }
        return false
      },
      set: function (newBool) {
        if (newBool) {
          // 閲覧 -> 編集
          router.push('/responses/' + this.responseId + '#edit')
        } else {
          // 編集 -> 閲覧
          router.push('/responses/' + this.responseId)
        }
      }
    },
    editButtonLink () {
      if (!this.isEditing) {
        return this.responseId + '#edit'
      } else {
        return this.responseId
      }
    }
  },
  mounted () {
  }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style scoped>
.tabs {
  margin-bottom: 0;
  margin-right: 0.5rem;
  margin-left: 0.5rem;
}
.tabs:first-child {
  margin-top: 1rem;
}
.is-editing {
  background-color: #c2c2c2;
}
#edit-button {
  border: #dbdbdb solid 1px;
}
</style>
