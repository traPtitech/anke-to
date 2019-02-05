<template>
  <div class="modal">
    <div class="modal-background"></div>
    <div class="modal-card">
      <header class="modal-card-head">
        <p class="modal-card-title">{{ activeModal.summary }}</p>
        <span
          class="ti-check modal-button"
          @click.prevent="confirmList"
          :class="{'disabled': !confirmOk}"
        ></span>
        <span class="ti-close modal-button" @click.prevent="disableModal"></span>
      </header>
      <section class="modal-card-body">
        <!-- Content ... -->
        <input-error-message :inputError="inputErrors.noAdministrator"></input-error-message>
        <label class="checkbox user-trap has-text-weight-bold">
          <input type="checkbox" v-model="isUserTrap">
          traP
        </label>
        <div class="user-list-wrapper">
          <span v-for="(user, index) in userTraqIdList" :key="index">
            <!-- user: traP -->
            <label v-if="!isUserTrap" class="checkbox">
              <input type="checkbox" v-model="selectedUserList" :value="user.Name">
              <span>{{ user.Name }}</span>
            </label>

            <!-- not user: traP -->
            <span v-if="isUserTrap" class="dummy-checkbox">
              <span class="readonly-checkbox checked"></span>
              {{ user.Name }}
            </span>
          </span>
        </div>
      </section>
    </div>
  </div>
</template>

<script>

import InputErrorMessage from '@/components/Utils/InputErrorMessage'
import common from '@/util/common'
import traQ from '@/util/traq'

export default {
  name: 'UserListModal',
  created () {
    this.selectedUserList = this.userListProps
    this.getUsersList()
  },
  beforeDestroy: function () {
    this.traq.disconnect()
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
    }
  },
  data () {
    return {
      traq: null,
      selectedUserList: [],
      userTraqIdList: []
    }
  },
  methods: {
    getUsersList () {
      this.traq = traQ('https://q.trapti.tech', true)
      this.traq.listen('connect', () => {
        this.traq.user.list(data => {
          this.userTraqIdList = data
          this.userTraqIdList.splice(0, 1)  // user: traP を取り除く
        })
      })
    },
    disableModal () {
      this.$emit('disable-modal')
    },
    confirmList () {
      if (this.confirmOk) {
        this.$emit('set-user-list', this.activeModal.name, this.selectedUserList)
        this.disableModal()
      }
    }
  },
  computed: {
    isUserTrap: {
      get () {
        return this.selectedUserList.length === 1 && this.selectedUserList[ 0 ] === 'traP'
      },
      set (newBool) {
        if (newBool) {
          this.selectedUserList = [ 'traP' ]
        } else {
          this.selectedUserList = [ this.traqId ]
        }
      }
    },
    inputErrors () {
      return {
        noAdministrator: {
          isError: this.activeModal.name === 'administrators' && this.selectedUserList.length === 0,
          message: '管理者がいません'
        }
      }
    },
    confirmOk () {
      return common.noErrors(this.inputErrors)
    }
  },
  mounted () {
  }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style lang="scss" scoped>
.modal-card-head {
  .modal-button {
    color: white;
    font-weight: bolder;
    width: 1.5rem;
    height: 1.5rem;
    padding: 0.25rem;
    margin-left: 1rem;
    border-radius: 1rem;
    &:hover {
      cursor: pointer;
    }
  }
  .ti-check {
    background-color: rgb(208, 255, 137);
    &:hover {
      background-color: greenyellow;
    }
    &.disabled {
      background-color: lightgray;
      pointer-events: none;
    }
  }
  .ti-close {
    background-color: rgb(255, 160, 160);
    &:hover {
      background-color: red;
    }
  }
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
