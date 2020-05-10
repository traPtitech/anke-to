<template>
  <div class="modal">
    <div class="modal-background" @click="disableModal"></div>
    <div class="modal-card">
      <header class="modal-card-head">
        <p class="modal-card-title">
          {{ activeModal.summary }} ({{ numberOfSelectedUsers }})
        </p>
        <div class="icon-button close round">
          <span class="ti-close" @click.prevent="disableModal"></span>
        </div>
      </header>
      <section class="modal-card-body">
        <!-- Content ... -->
        <div class="modal-body-top-wrapper">
          <label class="has-text-weight-bold">
            <checkbox :checked="isUserTrap" @input="isUserTrap = $event" />
            <span>traP</span>
          </label>
          <input
            v-model="searchQuery"
            placeholder="traQ ID でフィルター"
            class="input search"
          />
        </div>

        <!-- select type tab -->
        <div class="tabs is-centered">
          <ul>
            <li
              v-for="(tab, index) in getGroupTypes"
              :key="index"
              :class="{ 'is-active': selectedGroupType === tab }"
              class="tab"
              @click="selectedGroupType = tab"
            >
              <a>{{ tab !== '' ? tab : 'その他' }}</a>
            </li>
          </ul>
        </div>

        <!-- list -->
        <div class="user-list-wrapper">
          <span
            v-for="(group, index) in getGroupTypeMap[selectedGroupType]"
            :key="index"
          >
            <div class="has-text-weight-bold group-name">
              <checkbox
                v-show="searchQuery.length === 0"
                class="checkbox"
                :checked="isGroupSelectedMap[group.id]"
                @input="toggleIsGroupSelected(group.id)"
              />
              <span>
                {{ group.name }}
              </span>
            </div>

            <span
              v-for="(userId, index) in getFilteredActiveMembers(group.id)"
              :key="index"
            >
              <label class="checkbox-label">
                <input
                  v-model="userIsSelectedMap[getUsersMap[userId].name]"
                  type="checkbox"
                />
                <span>{{ getUsersMap[userId].name }}</span>
              </label>
            </span>
          </span>
        </div>
      </section>
      <footer class="modal-card-foot">
        <div
          v-for="(inputError, index) in errorsList"
          :key="index"
          class="error-message"
        >
          <span class="ti-alert"></span>
          <span>{{ inputError.message }}</span>
        </div>
        <button
          :class="{ disabled: !confirmOk }"
          class="button confirm"
          @click.prevent="confirmList"
        >
          <span class="ti-check" />
          決定
        </button>
      </footer>
    </div>
  </div>
</template>

<script>
import Checkbox from '@/components/Utils/Checkbox'
import common from '@/bin/common'
import { mapGetters } from 'vuex'

export default {
  name: 'UserListModal',
  components: {
    checkbox: Checkbox
  },
  props: {
    activeModal: {
      type: Object,
      required: true
    },
    userListProps: {
      type: Array,
      required: false,
      default: undefined
    },
    information: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      traq: null,
      selectedGroupType: '',
      userIsSelectedMap: {}, // userName をキー、そのユーザーが選択されているかどうかを値として持つ
      searchQuery: ''
    }
  },
  computed: {
    ...mapGetters(['getMyTraqId']),
    ...mapGetters('traq', [
      'getActiveUsers',
      'getSortedGroups',
      'getSortedGroupsMap',
      'getGroupTypes',
      'getGroupTypeMap',
      'getUsersMap'
    ]),
    isUserTrap: {
      get() {
        return !Object.values(this.userIsSelectedMap).includes(false)
      },
      set(newBool) {
        if (newBool) {
          this.selectAll()
        } else {
          this.removeAll()
        }
      }
    },
    isGroupSelectedMap() {
      // groupのidをキー、グループのすべてのメンバーが選択されているかどうかを値として持つ連想配列
      return Object.fromEntries(
        this.getSortedGroups.map(group => [
          group.id,
          this.isGroupSelected(group.id)
        ])
      )
    },
    numberOfSelectedUsers() {
      if (this.isUserTrap) {
        return Object.keys(this.getActiveUsers).length
      }
      return Object.values(this.userIsSelectedMap).filter(val => val).length
    },
    inputErrors() {
      return {
        noAdministrator: {
          isError:
            this.activeModal.name === 'administrators' &&
            this.numberOfSelectedUsers === 0,
          message: '管理者がいません'
        }
      }
    },
    errorsList() {
      return Object.values(this.inputErrors).filter(err => err.isError)
    },
    confirmOk() {
      return common.noErrors(this.inputErrors)
    },
    visibleUsersList() {
      let ret = Object.assign({}, this.allUsersList)
      delete ret.traP
      return ret
    }
  },
  watch: {},
  created() {
    // this.userIsSelectedMap を初期化
    if (this.information[this.activeModal.name] && this.getActiveUsers) {
      if (this.information[this.activeModal.name][0] === 'traP') {
        this.userIsSelectedMap = Object.fromEntries(
          this.getActiveUsers.map(user => [user.name, true])
        )
      } else {
        this.userIsSelectedMap = Object.fromEntries(
          this.getActiveUsers.map(user => [user.name, false])
        )
        this.information[this.activeModal.name].forEach(userName => {
          this.$set(this.userIsSelectedMap, userName, true)
        })
      }
    }

    // this.selectedGroupType を初期化
    this.selectedGroupType = this.getGroupTypes ? this.getGroupTypes[0] : ''
  },
  mounted() {},
  methods: {
    disableModal() {
      this.$emit('disable-modal')
    },
    confirmList() {
      if (this.confirmOk) {
        let selectedUsersList = []
        if (this.isUserTrap) {
          selectedUsersList = ['traP']
        } else {
          selectedUsersList = this.getActiveUsers
            .filter(user => this.userIsSelectedMap[user.name])
            .map(user => user.name)
        }
        this.$emit('set-user-list', this.activeModal.name, selectedUsersList)
        this.disableModal()
      }
    },
    selectAll() {
      Object.keys(this.userIsSelectedMap).forEach(userName => {
        this.userIsSelectedMap[userName] = true
      })
    },
    removeAll() {
      Object.keys(this.userIsSelectedMap).forEach(userName => {
        this.userIsSelectedMap[userName] = false
      })
    },
    isGroupSelected(groupId) {
      return !this.getSortedGroupsMap[groupId].activeMembers
        .map(userId => this.userIsSelectedMap[this.getUsersMap[userId].name])
        .includes(false)
    },
    toggleIsGroupSelected(groupId) {
      if (this.isGroupSelectedMap[groupId]) this.removeAllInGroup(groupId)
      else this.selectAllInGroup(groupId)
    },
    selectAllInGroup(groupId) {
      this.getActiveMembers(groupId).forEach(userId => {
        this.userIsSelectedMap[this.getUsersMap[userId].name] = true
      })
    },
    removeAllInGroup(groupId) {
      this.getActiveMembers(groupId).forEach(userId => {
        this.userIsSelectedMap[this.getUsersMap[userId].name] = false
      })
    },
    getActiveMembers(groupId) {
      return this.getSortedGroupsMap[groupId].activeMembers
    },
    getFilteredActiveMembers(groupId) {
      return this.getSortedGroupsMap[groupId].activeMembers.filter(userId =>
        this.getUsersMap[userId].name.includes(this.searchQuery)
      )
    }
  }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style lang="scss" scoped>
.icon-button {
  color: white;
  font-weight: bolder;
  width: 1.5rem;
  height: 1.5rem;
  padding: 0.25rem;
  margin-left: 1rem;
  display: flex;
  &.round {
    border-radius: 1rem;
  }
  &:hover {
    cursor: pointer;
  }
  &.close {
    background-color: $base-brown;
    &:hover {
      background-color: $var-red;
    }
  }
  span[class^='ti-'] {
    line-height: normal;
    font-size: small;
    margin: 10% auto 0 auto;
  }
}
.modal-body-top-wrapper {
  display: inline-flex;
  width: 100%;
  margin: 10px 0 20px 0;
  * {
    margin: auto 0;
  }
  label {
    display: flex;
  }
  .input.search {
    width: 13rem;
    margin-left: auto;
  }
}
.group-name {
  margin: 1.5rem 0 0.5rem 0;
  display: flex;
  * {
    margin: auto 0;
  }
  .checkbox {
    margin-right: 0.5rem;
  }
}
.modal-card-body {
  .checkbox-label {
    margin: 0.5rem;
    display: -webkit-inline-box;
    width: fit-content;
    span {
      margin-left: 0.2rem;
    }
  }
}
.error-message {
  color: $var-red;
  [class^='ti-'] {
    margin-right: 0.5rem;
  }
}
.button.confirm {
  margin-left: auto;
  background-color: $button-background-color-green;
  border: none;
  box-sizing: border-box;
  width: fit-content;
  &.disabled {
    border: $button-disabled-border-color 1px solid;
    color: $button-disabled-border-color;
    background-color: $button-disabled-background-color;
    pointer-events: none;
  }
  span {
    margin: 3px 4px 0 0;
  }
}
</style>
