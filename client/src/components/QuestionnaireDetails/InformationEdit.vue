<template>
  <div>
    <form>
      <div class="columns">
        <article class="column is-11">
          <div class="card">
            <!-- タイトル、説明、回答期限 -->
            <div>
              <header class="card-header">
                <div id="title" class="card-header-title title is-editing">
                  <div class="wrapper">
                    <textarea
                      :value="details.title"
                      @input="$set(details, 'title', $event.target.value)"
                      class="input"
                      placeholder="タイトル"
                    ></textarea>
                  </div>
                </div>
              </header>
              <input-error-message :inputError="inputErrors.noTitle"></input-error-message>
              <div class="card-content">
                <textarea
                  id="description"
                  v-model="details.description"
                  class="textarea"
                  rows="5"
                  placeholder="説明"
                ></textarea>
              </div>
              <div class="is-pulled-right is-inline-block wrapper">
                <div class="wrapper editable">
                  <span class="label">回答期限 :</span>
                  <input
                    class="input"
                    type="datetime-local"
                    v-model="resTimeLimitEditStr"
                    :disabled="noTimeLimit"
                  >
                </div>
                <label class="checkbox is-pulled-right">
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
                        v-if="userList.editable"
                        @click="changeActiveModal(userList)"
                      ></span>
                    </a>
                  </summary>

                  <p class="has-text-grey">{{ userList.liststr }}</p>
                </details>
                <input-error-message :inputError="inputErrors.noAdministrator"></input-error-message>

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
                        <span>{{ user }}</span>
                      </label>
                    </section>
                  </div>
                </div>

                <div class="wrapper editable">
                  <span class="label">結果の公開範囲:</span>
                  <span class="select">
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
          <div class="card">
            <!-- 操作 -->
            <div>
              <header class="card-header">
                <div class="card-header-title subtitle">操作</div>
              </header>
              <div class="card-content management-buttons">
                <button class="button" @click.prevent="deleteQuestionnaire">アンケートを削除</button>
              </div>
            </div>
          </div>
        </article>
      </div>
    </form>
  </div>
</template>

<script>

// import <componentname> from '<path to component file>'
import common from '@/util/common'
import InputErrorMessage from '@/components/Utils/InputErrorMessage'

export default {
  name: 'InformationEdit',
  components: {
    'input-error-message': InputErrorMessage
  },
  props: {
    informationProps: {
      type: Object,
      required: true
    },
    inputErrors: {
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
      userTraqIdList: [ 'mds_boy', '60', 'xxkiritoxx', 'yamada' ] // テスト用
      // noTimeLimit: true
    }
  },
  methods: {
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
    setInformation (information) {
      this.$emit('set-data', 'information', information)
    },
    changeActiveModal (obj) {
      this.activeModal = obj
      this.isModalActive = true
    },
    disableModal () {
      this.isModalActive = false
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
    noTimeLimit: {
      get () {
        return this.informationProps.noTimeLimit
      },
      set (newBool) {
        this.$emit('set-data', 'noTimeLimit', newBool)
      }
    },
    isNewQuestionnaire () {
      return this.$route.params.id === 'new'
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
          editable: true
        },
        respondents: {
          name: 'respondents',
          summary: '回答済みの人',
          list: this.details.respondents ? this.details.respondents.filter((user, index, array) => {
            // 重複除去
            return array.indexOf(user) === index
          }) : [],
          liststr: this.details.respondents ? this.toListString(this.details.respondents.filter((user, index, array) => {
            // 重複除去
            return array.indexOf(user) === index
          })) : '',
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
    traqId: function () {
      if (this.isNewQuestionnaire) {
        this.details.administrators = [ this.traqId ]
        this.details.targets = [ this.traqId ]
      }
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
.modal-card-body {
  .details.checkbox {
    margin: 0.5rem;
    display: -webkit-inline-box;
    width: fit-content;
  }
}
#title {
  .input {
    font-size: 2rem;
  }
  .wrapper {
    width: 100%;
  }
}
// .message {
//   margin-bottom: 0.5rem;
//   // .error-message {
//   //   font-size: 1rem;
//   //   margin: 0.5rem;
//   // }
// }
.editor-buttons {
  margin: auto;
  .button {
    margin: 1rem;
    // margin-bottom: 1rem;
    width: 8rem;
    max-width: 100%;
  }
}
</style>
