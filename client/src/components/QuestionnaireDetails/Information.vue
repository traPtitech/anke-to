<template>
  <div>
    <div class="columns">
      <article class="column is-11">
        <div class="card">
          <!-- タイトル、説明、回答期限 -->
          <div>
            <header class="card-header">
              <div id="title" class="card-header-title title">
                <div>{{ details.title }}</div>
              </div>
            </header>
            <div class="card-content">
              <pre>{{ details.description }}</pre>
            </div>
            <div class="is-pulled-right is-inline-block wrapper">
              <div class="wrapper editable">
                <span class="label">回答期限 :</span>
                <span class="time-limit-str">{{ getTimeLimitStr(details.res_time_limit) }}</span>
              </div>
            </div>
          </div>
        </div>
      </article>
    </div>

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
              <button
                class="button"
                @click.prevent="createResponse"
                :class="{'is-disabled': timeLimitExceeded}"
                :disabled="timeLimitExceeded"
              >新しい回答を作成</button>
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

// import <componentname> from '<path to component file>'
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
      newQuestionnaire: false
    }
  },
  methods: {
    getDateStr (str) {
      return common.customDateStr(str)
    },
    getTimeLimitStr (str) {
      return this.noTimeLimit ? 'なし' : common.customDateStr(str)
    },
    createResponse () {
      router.push({
        name: 'NewResponseDetails',
        params: {questionnaireId: this.questionnaireId}
      })
    },
    deleteResponse (responseId, index) {
      axios.delete('/responses/' + responseId, {method: 'delete', withCredentials: true})
      this.responses.splice(index, 1)
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
}
.editable.wrapper {
  display: flex;
}
.wrapper {
  .checkbox {
    width: 4rem;
    margin: 0.5rem;
  }
}
.management-buttons {
  .button:not(:last-child) {
    margin-bottom: 0.7rem;
  }
  .button {
    max-width: fit-content;
    display: block;
  }
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
#title {
  .input {
    font-size: 2rem;
  }
  .wrapper {
    width: 100%;
  }
  .error-message {
    font-size: 1rem;
    margin: 0.5rem;
  }
}
.editorbuttons {
  margin: auto;
  .button {
    margin: 0 1rem 2rem 1rem;
    // margin-bottom: 1rem;
    width: 8rem;
    max-width: 100%;
  }
}
</style>
