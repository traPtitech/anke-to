import moment from 'moment'
/* eslint-disable */

function customDateStr(str) {
  return str === 'NULL' ? '-' : moment(str).format('YYYY/MM/DD HH:mm')
}

function relativeDateStr(str) {
  return str === 'NULL'
    ? '-'
    : moment(str)
        .locale('ja')
        .fromNow()
}

export { customDateStr, relativeDateStr }
