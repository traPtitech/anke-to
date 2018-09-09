import axios from 'axios'

if (process.env.NODE_ENV === 'development') {
  axios.defaults.baseURL =
    'https://virtserver.swaggerhub.com/60-deg/anke-to/1.0.0/'
}

axios.defaults.withCredentials = true

export default axios
