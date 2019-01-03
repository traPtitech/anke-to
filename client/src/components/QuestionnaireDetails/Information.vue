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
                <span>{{ getDateStr(details.res_time_limit) }}</span>
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
                <span>{{ resSharedToStr }}</span>
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
              <button class="button" @click.prevent="createResponse">新しい回答を作成</button>
              <router-link
                :to="{ name: 'Results', params: { id: questionnaireId }}"
                class="button"
                :class="{'is-disabled' : !canViewResults}"
              >結果を見る</router-link>
              <button
                class="button"
                @click.prevent="deleteQuestionnaire"
                :class="{'is-disabled' : !administrates}"
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
                    title="po"
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
    this.getDetails()
    axios
      .get('/users/me/responses/' + this.questionnaireId)
      .then(res => {
        this.responses = res.data
      })
  },
  props: {
    traqId: {
      type: String,
      required: true
    }
  },
  data () {
    return {
      details: {},
      responses: [],
      activeModal: {},
      isModalActive: false,
      userTraqIdList: [ 'mds_boy', '60', 'xxkiritoxx', 'yamada' ], // テスト用
      noTimeLimit: true,
      newQuestionnaire: false
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
    deleteQuestionnaire () {
      axios
        .delete('/questionnaires/' + this.questionnaireId)
        .then(() => {
          router.push('/administrates') // アンケートを削除したら、Administratesページに戻る
        })
        .catch(error => {
          console.log(error)
        })
    },
    getDateStr (str) {
      return common.customDateStr(str)
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
          router.push('/responses/' + resp.data.responseID + '#edit')
        })
        .catch(error => {
          console.log(error)
        })
    },
    deleteResponse (responseId, index) {
      axios.delete('/responses/' + responseId, {method: 'delete', withCredentials: true})
      this.responses.splice(index, 1)
    }
  },
  computed: {
    questionnaireId () {
      return this.$route.params.id
    },
    administrates () {
      // 管理者かどうかを返す
      if (this.details.administrators) {
        for (let i = 0; i < this.details.administrators.length; i++) {
          if (this.traqId === this.details.administrators[ i ]) {
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
      if (typeof this.details.respondents === 'undefined') {
        return {}
      }
      return {
        targets: {
          name: 'targets',
          summary: '対象者',
          list: this.details.targets,
          liststr: this.toListString(this.details.targets),
          editable: true
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
          editable: true
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
    }
  },
  watch: {
    questionnaireId: function () {
      // 異なるquestionnaireIdのページに飛んだらdetailsをサーバーの状態に戻す
      this.getDetails()
    }
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
