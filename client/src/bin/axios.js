import axios from 'axios'

axios.defaults.baseURL = '/api'

if (process.env.NODE_ENV === 'development') {
  axios.defaults.baseURL = 'http://localhost:1323/api'
  // axios.defaults.baseURL = 'http://client.anke-to.sysad.trap.show/api'
  //    'https://virtserver.swaggerhub.com/60-deg/anke-to/1.0.0/'
}

axios.defaults.withCredentials = true

export default axios
