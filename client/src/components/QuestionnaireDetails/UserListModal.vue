<template>
  <div class="modal">
    <div class="modal-background"></div>
    <div class="modal-card">
      <header class="modal-card-head">
        <p class="modal-card-title">{{ activeModal.summary }}</p>
        <span
          class="ti-check icon-button round confirm"
          @click.prevent="confirmList"
          :class="{'disabled': !confirmOk}"
        ></span>
        <span class="ti-close icon-button round close" @click.prevent="disableModal"></span>
      </header>
      <section class="modal-card-body">
        <!-- Content ... -->
        <input-error-message :inputError="inputErrors.noAdministrator"></input-error-message>
        <label class="checkbox user-trap has-text-weight-bold">
          <input type="checkbox" v-model="isUserTrap">
          traP
        </label>
        <div class="user-list-wrapper">
          <span v-for="(group, key) in groupTypes[selectedGroupType]" :key="key">
            <div class="has-text-weight-bold group-name">
              {{ key }}
              <span
                v-if="!isUserTrap"
                class="ti-check icon-button select-group"
                @click.prevent="selectGroup(selectedGroupType, key)"
              ></span>
            </div>

            <!-- not user: traP -->
            <span v-for="(traqId, index) in group.activeMembers" :key="index">
              <label v-if="!isUserTrap" class="checkbox">
                <input type="checkbox" v-model="selectedUsersList" :value="traqId">
                <span>{{ traqId }}</span>
              </label>

              <!-- user: traP -->
              <span v-if="isUserTrap" class="dummy-checkbox">
                <span class="readonly-checkbox checked"></span>
                <span>{{ traqId }}</span>
              </span>
            </span>
          </span>
        </div>
      </section>
    </div>
  </div>
</template>

<script>

import InputErrorMessage from '@/components/Utils/InputErrorMessage'
import common from '@/bin/common'

export default {
  name: 'UserListModal',
  created () {
    this.selectedUsersList = this.userListProps
  },
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
      required: false
    },
    traqId: {
      required: true
    },
    users: {
      type: Object,
      required: true
    },
    groupTypes: {
      type: Object,
      required: true
    }
  },
  data () {
    return {
      traq: null,
      selectedUsersList: [],
      selectedGroupType: 'grade'
    }
  },
  methods: {
    disableModal () {
      this.$emit('disable-modal')
    },
    confirmList () {
      if (this.confirmOk) {
        this.$emit('set-user-list', this.activeModal.name, this.selectedUsersList)
        this.disableModal()
      }
    },
    selectGroup (type, group) {
      this.selectedUsersList =
        this.selectedUsersList
          .concat(this.groupTypes[ type ][ group ].activeMembers) // 該当するグループのユーザーを追加
          .filter((user, index, array) => { return array.indexOf(user) === index }) // 重複除去
    }
  },
  computed: {
    isUserTrap: {
      get () {
        return this.selectedUsersList.length === 1 && this.selectedUsersList[ 0 ] === 'traP'
      },
      set (newBool) {
        if (newBool) {
          this.selectedUsersList = [ 'traP' ]
        } else {
          this.selectedUsersList = []
        }
      }
    },
    inputErrors () {
      return {
        noAdministrator: {
          isError: this.activeModal.name === 'administrators' && this.selectedUsersList.length === 0,
          message: '管理者がいません'
        }
      }
    },
    confirmOk () {
      return common.noErrors(this.inputErrors)
    },
    visibleUsersList () {
      let ret = Object.assign({}, this.allUsersList)
      delete ret.traP
      return ret
    }
  },
  mounted () {
  }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style lang="scss" scoped>
@import "@/css/variables.scss";

.icon-button {
  color: white;
  font-weight: bolder;
  width: 1.5rem;
  height: 1.5rem;
  padding: 0.25rem;
  margin-left: 1rem;
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
