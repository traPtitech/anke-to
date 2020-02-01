export default {
  namespaced: true,
  state: {
    accessToken: null
  },
  getters: {},
  mutations: {
    setAccessToken(state, token) {
      state.accessToken = token
    }
  }
}
