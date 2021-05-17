import axios from 'axios'

export const traQBaseURL = 'https://q.trap.jp/api/v3'
axios.defaults.baseURL =
  process.env.NODE_ENV === 'development'
    ? 'http://localhost:8080/api'
    : 'https://anke-to.trap.jp/api'

export async function redirect2AuthEndpoint() {
  const data = (await axios.get('/oauth2/generate/code')).data

  const authorizationEndpointUrl = new URL(data)

  window.location.assign(authorizationEndpointUrl.toString())
}

export async function getRequest2Callback(to) {
  return axios.get('/oauth2/callback', {
    params: {
      code: to.query.code,
      state: to.query.state
    }
  })
}
