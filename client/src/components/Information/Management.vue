<template>
  <div>
    <header class="card-header">
      <div class="card-header-title subtitle">操作</div>
    </header>
    <div class="card-content management-buttons">
      <!-- 回答する -->
      <management-button
        type="newResponse"
        class="button-wrapper"
        :disabled="timeLimitExceeded"
        :questionnaireId="questionnaireId"
      >
      </management-button>
      <div class="new-response-link-panel">
        <input
          id="new-response-link"
          class="input"
          type="text"
          :value="newResponseLink"
          ref="link"
          @click="$refs.link.select()"
          readonly
        />
        <span class="button" @click="copyNewResponseLink">
          <span class="ti-clipboard"></span>
        </span>
      </div>
      <transition name="fade">
        <p class="copy-message" v-if="copyMessage.showMessage">
          {{ copyMessage.message }}
        </p>
      </transition>

      <!-- 結果を見る -->
      <management-button
        class="button-wrapper"
        type="viewResults"
        :disabled="!canViewResults"
        :questionnaireId="questionnaireId"
      >
      </management-button>

      <!-- アンケートを削除 -->
      <management-button
        class="button-wrapper"
        type="deleteQuestionnaire"
        :disabled="!administrates"
        :questionnaireId="questionnaireId"
      >
      </management-button>
    </div>
  </div>
</template>

<script>
import ManagementButton from '@/components/Information/ManagementButton'

export default {
  name: 'Management',
  components: {
    'management-button': ManagementButton
  },
  props: {
    res_time_limit: {
      type: String
    },
    questionnaireId: {
      type: Number
    },
    canViewResults: {
      type: Boolean
    },
    administrates: {
      type: Boolean
    }
  },
  data () {
    return {
      copyMessage: {
        showMessage: false
      }
    }
  },
  methods: {
    copyNewResponseLink () {
      let link = document.querySelector('#new-response-link')
      // link.select()
      let range = document.createRange()
      range.selectNode(link)
      window.getSelection().addRange(range)
      if (document.execCommand('copy')) {
        this.showCopyMessage('リンクをコピーしました！')
      } else {
        this.showCopyMessage('コピーに失敗しました')
      }
    },
    async showCopyMessage (message) {
      this.copyMessage = {
        showMessage: true,
        message: message
      }
      await new Promise(resolve => setTimeout(resolve, 3000))
      this.resetCopyMessage()
    },
    resetCopyMessage () {
      this.copyMessage = {
        showMessage: false
      }
    }
  },
  computed: {
    timeLimitExceeded () {
      // 回答期限を過ぎていた場合はtrueを返す
      return new Date(this.res_time_limit).getTime() < new Date().getTime()
    },
    newResponseLink () {
      return location.protocol + '//' + location.host + '/responses/new/' + this.questionnaireId
    }
  },
  watch: {
  },
  mounted () {
  }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style lang="scss" scoped>
.management-buttons {
  .new-response-link-panel {
    display: flex;
    input {
      width: -webkit-fill-available;
      max-width: 20rem;
      margin: 0;
      display: inherit;
      border-radius: 4px 0 0 4px;
      border-color: $base-gray;
    }
    .button {
      min-width: fit-content;
      margin: 0;
      display: inline-block;
      border-radius: 0 4px 4px 0;
      border-color: $base-gray;
      &:hover {
        background-color: $base-pink;
      }
    }
  }
  .copy-message {
    font-size: smaller;
    margin: 0.3rem;
  }
  .button-wrapper:not(:first-child) {
    margin-top: 0.7rem;
  }
}
</style>
