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
        <!-- error message -->
        <input-error-message
          :input-error="inputErrors.noAdministrator"
        ></input-error-message>

        <!-- user traP -->
        <label class="checkbox user-trap has-text-weight-bold">
          <input v-model="isUserTrap" type="checkbox" />
          traP
        </label>

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
              {{ group.name }}
              <span
                v-if="!isUserTrap && group.activeMembers.length > 0"
                class="ti-check icon-button select-group"
                @click.prevent="selectAllInGroup(selectedGroupType, index)"
              ></span>
              <span
                v-if="!isUserTrap && group.activeMembers.length > 0"
                class="ti-close icon-button select-group"
                @click.prevent="removeAllInGroup(selectedGroupType, index)"
              ></span>
            </div>

            <!-- not user: traP -->
            <span v-for="(userId, index) in group.activeMembers" :key="index">
              <label
                v-if="!isUserTrap && userId !== getMyTraqId"
                class="checkbox"
              >
                <input
                  v-model="usersIsSelected[getUsersMap[userId].name]"
                  type="checkbox"
                />
                <span>{{ getUsersMap[userId].name }}</span>
              </label>
            </span>
          </span>
        </div>
      </section>
      <footer class="modal-card-foot">
        <button
          :class="{ disabled: !confirmOk }"
          class="button"
          @click.prevent="confirmList"
        >
          決定
        </button>
      </footer>
    </div>
  </div>
</template>

<script>
import InputErrorMessage from '@/components/Utils/InputErrorMessage'
import common from '@/bin/common'
import { mapGetters } from 'vuex'

export default {
  name: 'UserListModal',
  components: {
    'input-error-message': InputErrorMessage
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
      selectedGroupType: this.getGroupTypes ? this.getGroupTypes[0] : '',
      usersIsSelected: {}
    }
  },
  computed: {
    ...mapGetters(['getMyTraqId']),
    ...mapGetters('traq', [
      'getActiveUsers',
      'getGroupTypes',
      'getGroupTypeMap',
      'getUsersMap'
    ]),
    isUserTrap: {
      get() {
        return this.usersIsSelected.traP === true
      },
      set(newBool) {
        if (newBool) {
          Object.keys(this.usersIsSelected).forEach(userName => {
            this.usersIsSelected[userName] = false
          })
        } else {
          this.usersIsSelected[this.getMyTraqId] = true
        }
        this.usersIsSelected.traP = newBool
      }
    },
    numberOfSelectedUsers() {
      if (this.isUserTrap) {
        return Object.keys(this.getActiveUsers).length
      }
      let count = 0
      Object.keys(this.usersIsSelected).forEach(userName => {
        if (this.usersIsSelected[userName]) {
          count++
        }
      })
      return count
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
    this.setUsersIsSelected(this.getActiveUsers)
  },
  mounted() {},
  methods: {
    disableModal() {
      this.$emit('disable-modal')
    },
    confirmList() {
      if (this.confirmOk) {
        let selectedUsersList = []
        Object.keys(this.usersIsSelected).forEach(userName => {
          if (this.usersIsSelected[userName]) {
            selectedUsersList.push(userName)
          }
        })
        this.$emit('set-user-list', this.activeModal.name, selectedUsersList)
        this.disableModal()
      }
    },
    selectAllInGroup(type, index) {
      this.getGroupTypeMap[type][index].activeMembers.forEach(userId => {
        this.usersIsSelected[this.getUsersMap[userId].name] = true
      })
    },
    removeAllInGroup(type, index) {
      this.getGroupTypeMap[type][index].activeMembers.forEach(userId => {
        this.usersIsSelected[this.getUsersMap[userId].name] = false
      })
    },
    setUsersIsSelected(users) {
      let tmp = {}
      if (
        Object.keys(users).length > 0 &&
        this.information.administrators &&
        this.information.targets
      ) {
        Object.keys(users).forEach(userId => {
          tmp[users[userId].name] = false
        })
        this.information[this.activeModal.name].forEach(userName => {
          tmp[userName] = true
        })
      }
      this.usersIsSelected = tmp
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
  &.confirm {
    background-color: $base-bluegray;
    &:hover {
      background-color: $var-indigo;
    }
    &.disabled {
      background-color: lightgray;
      pointer-events: none;
    }
  }
  &.close {
    background-color: $base-brown;
    &:hover {
      background-color: $var-red;
    }
  }
  &.select-group {
    background-color: $base-gray;
    &:hover {
      background-color: $base-darkbrown;
    }
  }
  span[class^='ti-'] {
    line-height: normal;
    font-size: small;
    margin: 10% auto 0 auto;
  }
}
.group-name {
  margin: 1.5rem 0 0.5rem 0;
}
.modal-card-body {
  .details.checkbox {
    margin: 0.5rem;
    display: -webkit-inline-box;
    width: fit-content;
  }
}
.checkbox,
.dummy-checkbox {
  width: 180px;
  display: inline-block;
  line-height: 1.25;
  margin: 0.5rem;
  &.user-trap {
    display: block;
    margin-bottom: 1rem;
  }
}
</style>
