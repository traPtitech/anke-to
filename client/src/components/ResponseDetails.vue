<template>
  <div class="is-fullheight details">
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
    <div :class="{'is-editing' : isEditing}" class="is-fullheight details-child">
      <questions :traqId="traqId" :editMode="isEditing? 'response' : undefined"></questions>
    </div>
  </div>
</template>

<script>

// import <componentname> from '<path to component file>'
import router from '@/router'
import Questions from '@/components/Questions'

export default {
  name: 'ResponseDetails',
  components: {
    'questions': Questions
  },
  props: {
    traqId: {
      type: String,
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
<style lang="scss" scoped>
</style>
