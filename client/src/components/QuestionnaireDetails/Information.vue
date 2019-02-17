<template>
  <div>
    <!-- <information-summary :details="summaryProps"></information-summary> -->
    <div class="columns details">
      <article class="column is-6">
        <div class="card">
          <!-- 情報 -->
          <div>
            <header class="card-header">
              <div class="card-header-title subtitle">情報</div>
            </header>
            <div class="card-content">
              <div class="has-text-weight-bold">
                <div>更新日時 : {{ getDateStr(details.modified_at) }}</div>
                <div>作成日時 : {{ getDateStr(details.created_at) }}</div>
              </div>

              <!-- user lists -->
              <details v-for="(userList, key) in userLists" :key="key">
                <summary>{{ userList.summary }}</summary>

                <p class="has-text-grey">{{ userList.liststr }}</p>
              </details>

              <div class="wrapper editable">
                <span class="label">結果の公開範囲:</span>
                <span>{{ resSharedToLabel }}</span>
              </div>
            </div>
          </div>
        </div>
      </article>

      <article class="column is-5">
        <div class="card">
          <!-- 操作 -->
          <div>
            <header class="card-header">
              <div class="card-header-title subtitle">操作</div>
            </header>
            <div class="card-content management-buttons">
              <router-link
                class="button"
                :class="{'is-disabled': timeLimitExceeded}"
                :disabled="timeLimitExceeded"
                :to="{name: 'NewResponseDetails', params: {questionnaireId: this.questionnaireId}}"
              >新しい回答を作成</router-link>
              <div class="new-response-link-panel">
                <input
                  id="new-response-link"
                  class="input"
                  type="text"
                  :value="newResponseLink"
                  ref="link"
                  @click="$refs.link.select()"
                  readonly
                >
                <span class="button" @click="copyNewResponseLink">
                  <span class="ti-clipboard"></span>
                </span>
              </div>
              <transition name="fade">
                <p class="copy-message" v-if="copyMessage.showMessage">{{ copyMessage.message }}</p>
              </transition>
              <router-link
                v-if="canViewResults"
                :to="{ name: 'Results', params: { id: questionnaireId }}"
                class="button"
              >結果を見る</router-link>
              <div v-if="!canViewResults" class="button is-disabled">結果を見る</div>
              <button
                class="button"
                @click.prevent="deleteQuestionnaire"
                :class="{'is-disabled' : !administrates}"
                :disabled="!administrates"
              >アンケートを削除</button>
            </div>
          </div>
        </div>
        <div class="card">
          <!-- 自分の回答一覧 -->
          <div>
            <header class="card-header">
              <div class="card-header-title subtitle">自分の回答</div>
            </header>
            <div class="card-content">
              <ul>
                <li v-for="(response, index) in responses" :key="index">
                  <span
                    :class="{'ti-save': response.submitted_at==='NULL', 'ti-check': response.submitted_at!=='NULL'}"
                  ></span>
                  <router-link
                    :to="'/responses/' + response.responseID"
                  >{{ getDateStr(response.modified_at) }}</router-link>
                  <a>
                    <span
                      class="ti-trash is-pulled-right"
                      @click="deleteResponse(response.responseID, index)"
                    ></span>
                  </a>
                </li>
              </ul>
            </div>
          </div>
        </div>
      </article>
    </div>
  </div>
</template>

<script>

import axios from '@/bin/axios'
import router from '@/router'
import common from '@/util/common'

export default {
  name: 'Information',
  components: {
  },
  async created () {
    axios
      .get('/users/me/responses/' + this.questionnaireId)
      .then(res => {
        this.responses = res.data
      })
  },
  props: {
    informationProps: {
      type: Object,
      required: true
    },
    traqId: {
      required: true
    }
  },
  data () {
    return {
      responses: [],
      activeModal: {},
      isModalActive: false,
      newQuestionnaire: false,
      copyMessage: {
        showMessage: false
      }
    }
  },
  methods: {
    getDateStr: common.customDateStr,
    createResponse () {
      router.push({
        name: 'NewResponseDetails',
        params: { questionnaireId: this.questionnaireId }
      })
    },
    deleteResponse (responseId, index) {
      if (window.confirm('この回答を削除しますか？')) {
        axios.delete('/responses/' + responseId, { method: 'delete', withCredentials: true })
        this.responses.splice(index, 1)
      }
    },
    copyNewResponseLink () {
      let link = document.querySelector('#new-response-link')
      link.select()
      if (document.execCommand('copy')) {
        this.copyMessage = {
          showMessage: true,
          message: '回答ページへのリンクをコピーしました'
        }
      } else {
        this.copyMessage = {
          showMessage: true,
          message: 'コピーに失敗しました'
        }
      }
    }
  },
  computed: {
    details () {
      return this.informationProps.details
    },
    administrates () {
      return this.informationProps.administrates
    },
    deleteQuestionnaire () {
      return this.informationProps.deleteQuestionnaire
    },
    questionnaireId () {
      return this.informationProps.questionnaireId
    },
    noTimeLimit () {
      return this.informationProps.noTimeLimit
    },
    canViewResults () {
      // 結果をみる権限があるかどうかを返す
      return common.canViewResults(this.details, this.administrates, this.responses.length > 0)
    },
    userLists () {
      return common.getUserLists(this.details)
    },
    resSharedToLabel () {
      const labels = {
        public: '全体',
        respondents: '回答済みの人',
        administrators: '管理者のみ'
      }
      return labels[ this.details.res_shared_to ]
    },
    timeLimitExceeded () {
      // 回答期限を過ぎていた場合はtrueを返す
      return new Date(this.details.res_time_limit).getTime() < new Date().getTime()
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
@import "@/css/variables.scss";
.card-header-title.is-editing {
  padding: 0;
}

.editable {
  span {
    width: fit-content;
    height: fit-content;
    top: 0;
    bottom: 0;
    margin: auto 0.2rem;
    white-space: nowrap;
  }
  &.wrapper {
    display: flex;
  }
}
.wrapper {
  .checkbox {
    width: 4rem;
    margin: 0.5rem;
  }
}
.management-buttons {
  .button:not(:first-child) {
    margin-top: 0.7rem;
  }
  .button {
    max-width: fit-content;
    display: block;
    &:hover {
      background-color: $base-pink;
    }
  }
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
}
.copy-message {
  font-size: smaller;
  margin: 0.3rem;
}
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s;
}
.fade-enter, .fade-leave-to /* .fade-leave-active below version 2.1.8 */ {
  opacity: 0;
}
.modal-card-head {
  .ti-check {
    background-color: darkgrey;
    color: white;
    font-weight: bolder;
    width: 1.5rem;
    height: 1.5rem;
    padding: 0.25rem;
    border-radius: 1rem;
  }
}
.editorbuttons {
  margin: auto;
  .button {
    margin: 0 1rem 2rem 1rem;
    width: 8rem;
    max-width: 100%;
  }
}
</style>
