import axios from 'axios'

axios.defaults.baseURL = '/api'

// if (process.env.NODE_ENV === 'development') {
//   axios.defaults.baseURL =
//     'https://virtserver.swaggerhub.com/60-deg/anke-to/1.0.0/'
// }

if (process.env.NODE_ENV === 'development') {
  axios.defaults.baseURL = 'http://anke-to.sysad.trap.show/'
}

// if (process.env.NODE_ENV === 'development') {
//   axios.defaults.baseURL = 'http://localhost:1323/api'
// }

axios.defaults.withCredentials = true

export default axios
