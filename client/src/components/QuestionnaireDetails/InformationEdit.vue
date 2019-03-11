<template>
  <div>
    <form>
      <div class="columns is-flex">
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
                  期限なし
                </label>
              </div>
            </div>
          </div>
        </article>
      </div>

      <div class="columns details">
        <article class="column is-6">
          <div class="card information-card">
            <!-- 情報 -->
            <div>
              <header class="card-header">
                <div class="card-header-title subtitle">情報</div>
              </header>
              <div class="card-content">
                <div v-for="(userList, key) in userLists" :key="key" class="user-list-wrapper">
                  <div>
                    <span class="has-text-weight-bold">{{ userList.summary }}</span>
                    <a>
                      <span
                        class="ti-pencil"
                        v-if="userList.editable"
                        @click="changeActiveModal(userList)"
                      ></span>
                    </a>
                  </div>
                  <p class="has-text-grey user-list">{{ userList.liststr }}</p>
                </div>

                <!-- modal -->
                <user-list-modal
                  v-if="isModalActive"
                  :class="{'is-active': isModalActive}"
                  :activeModal="activeModal"
                  :userListProps="details[activeModal.name]"
                  :traqId="traqId"
                  :users="users"
                  :groupTypes="groupTypes"
                  :information="details"
                  @disable-modal="disableModal"
                  @set-user-list="setUserList"
                ></user-list-modal>

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

import common from '@/bin/common'
import axios from '@/bin/axios'
import InputErrorMessage from '@/components/Utils/InputErrorMessage'
import UserListModal from '@/components/QuestionnaireDetails/UserListModal'

export default {
  name: 'InformationEdit',
  created () {
    this.getUsers()
      .then(this.getGroupTypes)
  },
  components: {
    'input-error-message': InputErrorMessage,
    'user-list-modal': UserListModal
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
      users: {},
      groupTypes: {}
      // usersIsSelected: {}
    }
  },
  methods: {
    getDateStr (str) {
      return common.customDateStr(str)
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
    },
    setUserList (listName, newList) {
      this.details[ listName ] = newList
      this.setInformation(this.details)
    },
    getUsers () {
      return axios
        .get('https://q.trap.jp/api/1.0/users')
        .then(res => {
          res.data.forEach(user => {
            if (user.accountStatus === 1) {
              this.users[ user.userId ] = user
            }
          })
        })
        .catch(err => {
          console.log(err)
        })
    },
    getGroupTypes () {
      return axios
        .get('https://q.trap.jp/api/1.0/groups')
        .then(res => {
          let tmp = {}
          res.data.forEach(group => {
            if (typeof tmp[ group.type ] === 'undefined') {
              tmp[ group.type ] = []
            }
            // 除名されていないメンバーをtraQID順にソートしたtraQIDのリストactiveMembersを作成
            group.activeMembers =
              group.members.filter(userId => typeof this.users[userId] !== 'undefined' && this.users[ userId ].accountStatus === 1 && this.users[ userId ].name !== 'traP')
                .map(userId => this.users[ userId ].name)
                .sort((a, b) => { return a.toLowerCase().localeCompare(b.toLowerCase()) })
            tmp[ group.type ].push(group)
          })

          // typeごとに、group名をソートしたものをgroupTypesに入れる
          Object.keys(tmp).forEach(type => {
            this.$set(this.groupTypes, type, {})
            tmp[ type ]
              .sort((a, b) => { return a.name.toLowerCase().localeCompare(b.name.toLowerCase()) })
          })
          this.groupTypes = tmp
        })
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
      return common.getUserLists(this.details)
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
        if (str === '') {
          this.$emit('set-data', 'noTimeLimit', true)
        } else {
          this.details.res_time_limit = str
        }
      }
    }
  },
  watch: {
    traqId: function (newVal) {
      // traqIdがundefinedから変わった時に呼ばれる
      if (newVal && this.isNewQuestionnaire) {
        this.details.administrators = [ this.traqId ]
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
.details .checkbox {
  width: 6rem;
  margin: 0.5rem;
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
.user-list {
  margin: 0 0.5rem;
}
#title {
  .input {
    font-size: 2rem;
  }
  .wrapper {
    width: 100%;
  }
}
.editor-buttons {
  margin: auto;
  .button {
    margin: 1rem;
    // margin-bottom: 1rem;
    width: 8rem;
    max-width: 100%;
  }
}
.user-list-wrapper {
  margin-bottom: 0.5rem;
}
</style>
