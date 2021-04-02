import axios from 'axios'
import { Apis } from '@traptitech/traq'

let apis = null

const setAuthToken = token => {
  if (token && !apis) {
    apis = new Apis(
      { accessToken: token }, // configuration
      undefined, // base path (デフォルト値を入れるように)
      axios.create({ withCredentials: false }) // axios instance (withCredentials: false が必要)
    )
  }
}

export default {
  namespaced: true,
  state: {
    accessToken: null,
    accessTokenEnsured: false,
    users: null,
    groups: null
  },
  getters: {
    getUsersMap(state) {
      // userのidをキーとするusersの連想配列
      if (!state.users) return null
      return Object.fromEntries(state.users.map(user => [user.id, user]))
    },
    getActiveUsers(state) {
      if (!state.users) return null
      return state.users
        .filter(user => !user.suspended && user.name !== 'traP' && !user.bot)
        .sort((a, b) => {
          return a.name.toLowerCase().localeCompare(b.name.toLowerCase())
        })
    },
    getActiveUsersMap(_, getters) {
      // userのidをキーとするactiveなusersの連想配列
      if (!getters.getActiveUsers) return null
      return Object.fromEntries(
        getters.getActiveUsers.map(user => [user.id, user])
      )
    },
    getGroupsMap(state) {
      // groupのidをキーとするgroupsの連想配列
      if (!state.groups) return null
      return Object.fromEntries(state.groups.map(group => [group.id, group]))
    },
    getSortedGroups(state, getters) {
      if (!state.groups || !getters.getActiveUsersMap) return null

      // グループ名でソート
      let groups = [...state.groups]
      groups.sort((a, b) =>
        a.name.toLowerCase().localeCompare(b.name.toLowerCase())
      )

      // グループメンバーをソート
      const sortGroupMembers = data => {
        let group = JSON.parse(JSON.stringify(data)) // deep copy
        group.members = group.members
          .filter(member => getters.getUsersMap[member.id])
          .sort((a, b) => {
            const nameA = getters.getUsersMap[a.id].name
            const nameB = getters.getUsersMap[b.id].name
            return nameA.toLowerCase().localeCompare(nameB.toLowerCase())
          })
        // activeなメンバーのみのリストを作る
        group.activeMembers = group.members
          .filter(
            member =>
              getters.getActiveUsersMap[member.id] &&
              getters.getActiveUsersMap[member.id].name !== 'traP'
          )
          .map(member => member.id)
        return group
      }

      groups = groups
        .map(group => sortGroupMembers(group))
        .filter(group => group.activeMembers.length > 0) // activeMembersが一人もいないグループを削除

      return groups
    },
    getSortedGroupsMap(_, getters) {
      // groupのidをキーとするgroupsの連想配列
      if (!getters.getSortedGroups) return null
      return Object.fromEntries(
        getters.getSortedGroups.map(group => [group.id, group])
      )
    },
    getGroupTypes(state) {
      // groupのtypeのリストを返す
      if (!state.groups) return null
      return state.groups
        .map(group => group.type)
        .filter((type, i, self) => self.indexOf(type) === i) // 重複除去
    },
    getGroupTypeMap(state, getters) {
      // typeをキー、そのtypeを持つgroupの配列を値として持つ連想配列
      if (!state.users || !state.groups) return null

      let ret = Object.fromEntries(
        getters.getGroupTypes.map(type => [type, []])
      )

      getters.getSortedGroups.forEach(group => {
        ret[group.type].push(group)
      })

      return ret
    }
  },
  mutations: {
    setAccessToken(state, token) {
      state.accessToken = token
    },
    setAccessTokenEnsured(state, ensured) {
      state.accessTokenEnsured = ensured
    },
    setUsers(state, users) {
      state.users = users
    },
    setGroups(state, groups) {
      state.groups = groups
    }
  },
  actions: {
    async ensureToken({ state, commit }) {
      if (!state.accessToken) {
        return
      }
      if (state.accessTokenEnsured) {
        return
      }
      setAuthToken(state.accessToken)

      try {
        await apis.getMe()
        commit('setAccessTokenEnsured', true)
      } catch {
        commit('setAccessToken', null)
      }
    },
    async updateUsers({ state, commit }) {
      if (!state.accessToken) {
        console.error('no access token')
        return
      }
      setAuthToken(state.accessToken)

      await apis
        .getUsers()
        .then(res => {
          commit('setUsers', res.data)
        })
        .catch(err => {
          console.log(err)
        })
    },
    async updateGroups({ state, commit }) {
      if (!state.accessToken) {
        console.error('no access token')
        return
      }
      setAuthToken(state.accessToken)

      await apis
        .getUserGroups()
        .then(res => {
          commit('setGroups', res.data)
        })
        .catch(err => {
          console.log(err)
        })
    }
  }
}
