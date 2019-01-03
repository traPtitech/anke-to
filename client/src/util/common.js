import moment from 'moment'
/* eslint-disable */

export default {
  customDateStr: function(str) {
    return str === 'NULL' ? '-' : moment(str).format('YYYY/MM/DD HH:mm')
  },

  relativeDateStr: function(str) {
    return str === 'NULL'
      ? '-'
      : moment(str)
          .locale('ja')
          .fromNow()
  },

  swapOrder: function(arr, i0, i1) {
    if (i0 < 0 || i1 < 0 || i0 >= arr.length || i1 >= arr.length) return
    let tmp = arr[i0]
    this.$set(arr, i0, arr[i1])
    this.$set(arr, i1, tmp)
  }
}
