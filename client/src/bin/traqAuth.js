import axios from 'axios'
import sha256 from 'js-sha256'

const baseUrl = 'https://q.trap.jp/api/v3/oauth2'

const getTraqClientId = () => {
  const clientId = process.env.VUE_APP_TRAQ_CLIENT_ID
  if (!clientId) {
    console.error('client ID not set')
    return
  }
  return clientId
}

const randomString = len => {
  const characters =
    'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
  const charactarsLength = characters.length

  let array = window.crypto.getRandomValues(new Uint32Array(len))
  array = array.map(val => characters.charCodeAt(val % charactarsLength))
  return String.fromCharCode(...array)
}

const getCodeChallenge = codeVerifier => {
  const sha256Result = sha256(codeVerifier)
  const bytes = new Uint8Array(sha256Result.length / 2)
  for (let i = 0; i < sha256Result.length; i += 2) {
    bytes[i / 2] = parseInt(sha256Result.substring(i, i + 2), 16)
  }
  const base64 = btoa(String.fromCharCode(...bytes))
  const base64url = base64
    .replace(/\+/g, '-')
    .replace(/\//g, '_')
    .replace(/=/g, '')
  return base64url
}

export function sendCodeRequest() {
  const url = baseUrl + '/authorize'

  const state = randomString(10)
  const codeVerifier = randomString(43)
  const codeChallenge = getCodeChallenge(codeVerifier)

  sessionStorage.setItem(`traq-auth-code-verifier-${state}`, codeVerifier)

  const params = new URLSearchParams({
    response_type: 'code',
    client_id: getTraqClientId(),
    state: state,
    code_challenge: codeChallenge,
    code_challenge_method: 'S256'
  })
  window.location.assign(new URL(url + '?' + params.toString()))
}

export async function sendTokenRequest(code, codeVerifier) {
  const url = baseUrl + '/token'

  const params = new URLSearchParams({
    grant_type: 'authorization_code',
    client_id: getTraqClientId(),
    code: code,
    code_verifier: codeVerifier
  })
  return axios.post(url, params, { withCredentials: false })
}
