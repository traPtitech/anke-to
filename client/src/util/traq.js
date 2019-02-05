/* eslint-disable */

import io from 'socket.io-client'
export default (endpoint, debug) => {
  const socket = io(endpoint)
  const events = ['connect', 'disconnect', 'user/list']
  const listeners = {}

  const log = (...logs) => {
    if (!debug) return
    logs.unshift(new Date().toLocaleString() + ' - ')
    console.info(...logs)
  }

  for (let event of events) {
    const evName = event
    listeners[evName] = []
    socket.on(evName, data => {
      log(`>>> ${evName} :`, data)
      for (let listener of listeners[evName]) listener(data)
    })
  }

  const emit = (job, data, callback) => {
    log(`<<< ${job} :`, data)
    socket.emit(job, data, data => {
      if (callback) callback(data)
      log(`>>> ${job} <finished> :`, data)
      for (let listener of listeners[job] || []) listener()
    })
  }

  const apis = {
    connect: _ => socket.connect(),
    disconnect: _ => socket.disconnect(),
    listen: (event, listener) => listeners[event].push(listener),
    user: {
      list: callback => emit('user/list', {}, callback)
    }
  }

  const makeAsync = obj => {
    const newObj = {}
    for (let key in obj) {
      if (typeof obj[key] === 'object') {
        newObj[key] = makeAsync(obj[key])
      } else if (typeof obj[key] === 'function') {
        newObj[key] = (...args) =>
          new Promise(resolve => obj[key](...args, resolve))
      }
    }
    return newObj
  }

  apis.async = makeAsync(apis)
  return apis
}
