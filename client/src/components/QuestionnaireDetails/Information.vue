<template>
  <div>
    <form class="questionnaire-details" @submit="submitQuestionnaire">
      <div class="columns">
        <article class="column is-11">
          <div class="card">
            <!-- タイトル、説明、回答期限 -->
            <div>
              <header class="card-header">
                <div class="card-header-title title" :class="{'is-editing' : isEditing}">
                  <input v-show="isEditing" id="title" v-model="details.title" class="input">
                  <div v-show="!isEditing">{{ details.title }}</div>
                </div>
              </header>
              <div class="card-content">
                <textarea
                  v-show="isEditing"
                  id="description"
                  v-model="details.description"
                  class="textarea"
                  rows="5"
                ></textarea>
                <pre v-show="!isEditing">{{ details.description }}</pre>
              </div>
              <div class="is-pulled-right is-inline-block wrapper">
                <div class="wrapper editable">
                  <span class="label">回答期限 :</span>
                  <span v-show="!isEditing">{{ getDateStr(details.res_time_limit) }}</span>
                  <input
                    v-show="isEditing"
                    class="input"
                    type="datetime-local"
                    v-model="resTimeLimitEditStr"
                    :disabled="noTimeLimit"
                  >
                </div>
                <label class="checkbox is-pulled-right" v-show="isEditing">
                  <input type="checkbox" v-model="noTimeLimit">
                  なし
                </label>
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
                  <summary>
                    {{ userList.summary }}
                    <a>
                      <span
                        class="ti-pencil"
                        v-show="userList.editable"
                        @click="changeActiveModal(userList)"
                      ></span>
                    </a>
                  </summary>

                  <p class="has-text-grey">{{ userList.liststr }}</p>
                </details>

                <!-- modal -->
                <div class="modal" v-if="isModalActive" :class="{'is-active': isModalActive}">
                  <div class="modal-background"></div>
                  <div class="modal-card">
                    <header class="modal-card-head">
                      <p class="modal-card-title">{{ activeModal.summary }}</p>
                      <span class="ti-check" @click.prevent="disableModal"></span>
                    </header>
                    <section class="modal-card-body">
                      <!-- Content ... -->
                      <label class="checkbox" v-for="(user, index) in userTraqIdList" :key="index">
                        <input type="checkbox" v-model="details[activeModal.name]" :value="user">
                        {{ user }}
                      </label>
                    </section>
                  </div>
                </div>

                <div class="wrapper editable">
                  <span class="label">結果の公開範囲:</span>
                  <span v-show="!isEditing">{{ resSharedToStr }}</span>
                  <span v-show="isEditing" class="select">
                    <select v-model="resSharedToStr">
                      <option>全体</option>
                      <option>回答済みの人</option>
                      <option>管理者のみ</option>
                    </select>
                  </span>
                </div>
              </div>
            </div>
          </div>
        </article>
        <article class="column is-5">
          <div class="card" v-show="!isEditing">
            <!-- 操作 -->
            <div>
              <header class="card-header">
                <div class="card-header-title subtitle">操作</div>
              </header>
              <div class="card-content management-buttons">
                <!-- <button class="button" :href="questionnaireId + '/new-response'">新しい回答を作成</button> -->
                <button class="button" @click.prevent="createResponse">新しい回答を作成</button>
                <router-link
                  :to="{ name: 'Results', params: { id: questionnaireId }}"
                  class="button"
                  :class="{'is-disabled' : !canViewResults}"
                >結果を見る</router-link>
              </div>
            </div>
          </div>
          <div class="card" v-show="!isEditing">
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
                      title="po"
                    ></span>
                    <a
                      :href="'/responses/' + response.responseID"
                    >{{ getDateStr(response.modified_at) }}</a>
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
      <div class="columns" v-show="isEditing">
        <div class="column is-9"></div>
        <article class="column is-2 is-flex">
          <div class="editorbuttons">
            <input
              type="submit"
              value="送信"
              class="button is-medium is-pulled-right"
              id="submitbutton"
              @click.prevent="submitQuestionnaire"
            >
            <button
              class="button is-medium is-pulled-right"
              @click.prevent="$emit('disable-editing')"
            >キャンセル</button>
          </div>
        </article>
      </div>
    </form>
  </div>
</template>

<script>

// import <componentname> from '<path to component file>'
import axios from '@/bin/axios'
import router from '@/router'
import moment from 'moment'
import {customDateStr} from '@/util/common'

export default {
  name: 'Information',
  components: {
  },
  async created () {
    this.getDetails()
    axios
      .get('/users/me/responses/' + this.questionnaireId)
      .then(res => {
        this.responses = res.data
      })
  },
  props: {
    props: {
      type: Object,
      required: true
    }
  },
  data () {
    return {
      details: {},
      responses: [],
      questionnaireId: this.$route.params.id,
      activeModal: {},
      isModalActive: false,
      userTraqIdList: [ 'mds_boy', '60', 'xxkiritoxx', 'yamada' ], // テスト用
      noTimeLimit: true
    }
  },
  methods: {
    getDetails () {
      // サーバーにアンケートの情報をリクエストする
      axios
        .get('/questionnaires/' + this.questionnaireId)
        .then(res => {
          this.details = res.data
          if (this.administrates) {
            this.$emit('enable-edit-button')
          }
          if (this.details.res_time_limit && this.details.res_time_limit !== 'NULL') {
            this.noTimeLimit = false
          }
        })
    },
    submitQuestionnaire () {
      this.$emit('disable-editing') // 編集モード終了
      const data = {
        title: this.details.title,
        description: this.details.description,
        res_time_limit: this.noTimeLimit ? 'NULL' : new Date(this.details.res_time_limit).toLocaleString(),
        res_shared_to: this.details.res_shared_to,
        targets: this.details.targets,
        administrators: this.details.administrators
      }
      axios.patch('/questionnaires/' + this.questionnaireId, data)
        // PATCHリクエストを送る
        .then(this.getDetails)
        // detailsをアップデート
        .catch(function (error) {
          console.log(error)
        })
    },
    getDateStr (str) {
      return customDateStr(str)
    },
    toListString (list) {
      if (list && list.length === 0) {
        return ''
      }
      let ret = ''
      for (let i = 0; i < list.length - 1; i++) {
        ret += list[ i ] + ', '
      }
      ret += list[ list.length - 1 ]
      return ret
    },
    createResponse () {
      const data = {
        questionnaireID: parseInt(this.questionnaireId),
        submitted_at: 'NULL',
        body: []
      }
      axios
        .post('/responses', data)
        .then(resp => {
          // POSTリクエストで返ってきたresponseIDをもとに、responses/:responseID に編集モードで飛ぶ
          // router.push('/responses/' + resp.data.responseID + '#edit')
          router.push('/responses/' + 1 + '#edit') // テスト用
        })
        .catch(error => {
          console.log(error)
        })
    },
    deleteResponse (responseId, index) {
      axios.delete('/responses/' + responseId, {method: 'delete', withCredentials: true})
      this.responses.splice(index, 1)
    },
    changeActiveModal (obj) {
      this.activeModal = obj
      this.isModalActive = true
      console.log(this.activeModal)
    },
    disableModal () {
      this.isModalActive = false
    }
  },
  computed: {
    traqId () {
      return this.props.traqId
    },
    isEditing () {
      return this.props.isEditing
    },
    administrates () {
      // 管理者かどうかを返す
      if (this.details.administrators) {
        for (let i = 0; i < this.details.administrators.length; i++) {
          if (this.props.traqId === this.details.administrators[ i ]) {
            return true
          }
        }
      }
      return false
    },
    canViewResults () {
      // 結果をみる権限があるかどうかを返す
      return ((this.details.res_shared_to === 'public') ||
        (this.details.res_shared_to === 'administrators' && this.administrates) ||
        (this.details.res_shared_to === 'respondents' && this.responses.length > 0))
    },
    userLists () {
      if (!this.details.targets) {
        return {}
      }
      return {
        targets: {
          name: 'targets',
          summary: '対象者',
          list: this.details.targets,
          liststr: this.toListString(this.details.targets),
          editable: this.isEditing
        },
        respondents: {
          name: 'respondents',
          summary: '回答済みの人',
          list: this.details.respondents.filter((user, index, array) => {
            // 重複除去
            return array.indexOf(user) === index
          }),
          liststr: this.toListString(this.details.respondents.filter((user, index, array) => {
            // 重複除去
            return array.indexOf(user) === index
          })),
          editable: false
        },
        administrators: {
          name: 'administrators',
          summary: '管理者',
          list: this.details.administrators,
          liststr: this.toListString(this.details.administrators),
          editable: this.isEditing
        }
      }
    },
    resSharedToStr: {
      get: function () {
        switch (this.details.res_shared_to) {
          case 'public': return '全体'
          case 'administrators': return '管理者のみ'
          case 'respondents': return '回答済みの人'
        }
      },
      set: function (str) {
        switch (str) {
          case '全体': {
            this.details.res_shared_to = 'public'
            break
          }
          case '管理者のみ': {
            this.details.res_shared_to = 'administrators'
            break
          }
          case '回答済みの人': {
            this.details.res_shared_to = 'respondents'
            break
          }
        }
      }
    },
    resTimeLimitEditStr: {
      get: function () {
        if (!this.details.res_time_limit || this.details.res_time_limit === 'NULL') return ''
        return this.details.res_time_limit.slice(0, 16)
      },
      set: function (str) {
        this.details.res_time_limit = str
      }
    }
  },
  watch: {
    isEditing: function (newBool, oldBool) {
      if (oldBool && !newBool) {
        // 編集モードから閲覧モードに変わったときはdetailsをサーバーの状態に戻す
        this.$emit('disable-editing')
        this.getDetails()
      }
    },
    noTimeLimit: function (newBool, oldBool) {
      if (oldBool && !newBool && this.details.res_time_limit === 'NULL') {
        // 新しく回答期限を作ろうとすると、1週間後の日時が設定される
        this.details.res_time_limit = moment().add(7, 'days').format().slice(0, -6)
      }
    }
  },
  mounted () {
  }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style lang="scss" scoped>
pre {
  white-space: pre-line;
  font-size: inherit;
  -webkit-font-smoothing: inherit;
  font-family: inherit;
  line-height: inherit;
  background-color: inherit;
  color: inherit;
  padding: 0.625em;
}
.card {
  max-width: 100%;
  padding: 0.7rem;
}
article.column {
  padding: 0;
}
.columns {
  margin-bottom: 0;
}
.columns:first-child {
  display: flex;
}

.card-header-title.is-editing {
  padding: 0;
}
.card-content {
  .subtitle {
    margin: 0;
  }
  details {
    margin: 0.5rem;
    p {
      padding: 0 0.5rem;
    }
  }
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
#title.input {
  font-size: 2rem;
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
@media screen and (min-width: 769px) {
  // widthが大きいときは横並びのカードの間を狭くする
  .column:not(:last-child) > .card {
    margin-right: 0;
  }
}
</style>
