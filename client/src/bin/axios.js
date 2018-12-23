import axios from 'axios'

if (process.env.NODE_ENV === 'development') {
  axios.defaults.baseURL =
    'https://virtserver.swaggerhub.com/60-deg/anke-to/1.0.0/'
} else {
  axios.defaults.baseURL = '/api'
}

// if (process.env.NODE_ENV === 'development') {
//   axios.defaults.baseURL = 'http://anke-to.sysad.trap.show/'
// } else {
//   axios.defaults.baseURL = '/api'
// }

axios.defaults.withCredentials = true

export default axios
