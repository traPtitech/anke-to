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
      <questions
        :traqId="traqId"
        :editMode="isEditing? 'response' : undefined"
        :questions="questions"
      ></questions>
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
      questions: [ // テスト用
        {
          questionId: 1,
          type: 'Text',
          component: 'short-answer',
          questionBody: 'ペンは持っていますか？',
          body: 'はい'
        },
        {
          questionId: 2,
          type: 'Number',
          component: 'short-answer',
          questionBody: 'ペンは何本持っていますか？',
          body: '12'
        },
        {
          questionId: 3,
          type: 'Checkbox',
          component: 'multiple-choice',
          questionBody: '何色のペンを持っていますか？',
          options: [
            {
              label: '赤',
              id: 0
            },
            {
              label: '青',
              id: 1
            },
            {
              label: '黄色',
              id: 2
            }
          ],
          isSelected: [ false, false, false ]
        },
        {
          questionId: 4,
          type: 'MultipleChoice',
          component: 'multiple-choice',
          questionBody: '何色のペンが欲しいですか？',
          options: [
            {
              label: '赤',
              id: 0
            },
            {
              label: '青',
              id: 1
            },
            {
              label: '黄色',
              id: 2
            }
          ],
          selected: ''
        },
        {
          questionId: 5,
          type: 'LinearScale',
          component: 'linear-scale',
          questionBody: '好きなペンの太さは？',
          scaleRange: {
            left: 0,
            right: 10
          },
          scaleLabels: {
            left: '細い',
            right: '太い'
          },
          selected: 3
        }
      ]
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
