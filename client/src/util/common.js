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
  },

  convertDataToQuestion(data) {
    // サーバーから送られてきた質問1つ分のデータを、Questionsで使えるフォーマットに変換して返す
    let question = {
      questionId: data.questionID,
      type: data.question_type,
      component: this.questionTypes[data.question_type].component,
      questionBody: data.body,
      isRequired: data.is_required,
      pageNum: data.page_num
    }
    switch (data.question_type) {
      case 'Checkbox':
        question.options = []
        question.isSelected = {}
        data.options.forEach((option, index) => {
          question.options.push({
            id: index,
            label: option
          })
          question.isSelected[option] = false
        })
        break
      case 'MultipleChoice':
        question.options = []
        data.options.forEach((option, index) => {
          question.options.push({
            id: index,
            label: option
          })
        })
        break
      case 'LinearScale':
        question.scaleLabels = {
          left: data.scale_label_left,
          right: data.scale_label_right
        }
        question.scaleRange = {
          left: data.scale_min,
          right: data.scale_max
        }
      default:
        break
    }
    return question
  },

  setResponseToQuestion(questionData, responseData) {
    // サーバーから送られてきた回答1つ分のデータを、指定されたquestionに入れる
    let question = Object.assign({}, questionData)
    switch (question.type) {
      case 'Text':
        question.responseBody = responseData.response
        break
      case 'Number':
        question.responseBody = Number(responseData.response)
        break
      case 'Checkbox':
        question.isSelected = {}
        responseData.option_response.forEach(selectedOption => {
          question.isSelected[selectedOption] = true
        })
        break
      case 'MultipleChoice':
        question.selected = responseData.option_response[0]
        break
      case 'LinearScale':
        question.selected = Number(responseData.response)
        break
      default:
        break
    }
    return question
  },

  questionTypes: {
    Text: {
      type: 'Text',
      label: 'テキスト',
      component: 'short-answer'
    },
    Number: {
      type: 'Number',
      label: '数値',
      component: 'short-answer'
    },
    Checkbox: {
      type: 'Checkbox',
      label: 'チェックボックス',
      component: 'multiple-choice'
    },
    MultipleChoice: {
      type: 'MultipleChoice',
      label: 'ラジオボタン',
      component: 'multiple-choice'
    },
    LinearScale: {
      type: 'LinearScale',
      label: '目盛り',
      component: 'linear-scale'
    }
  },

  alertNetworkError() {
    alert('Network Error')
  },

  administrates(administrators, traqId) {
    if (administrators[0] === 'traP') {
      return true
    }
    for (let i = 0; i < administrators.length; i++) {
      if (traqId === administrators[i]) {
        return true
      }
    }
    return false
  },

  canViewResults(information, administrates, hasResponded) {
    return (
      information.res_shared_to === 'public' ||
      (information.res_shared_to === 'administrators' && administrates) ||
      (information.res_shared_to === 'respondents' && hasResponded)
    )
  },
  toListString(list) {
    let ret = ''
    if (list.length > 0) {
      for (let i = 0; i < list.length - 1; i++) {
        ret += list[i] + ', '
      }
      ret += list[list.length - 1]
    }
    return ret
  },
  getUserLists(details) {
    if (details.targets && details.respondents && details.administrators) {
      return {
        targets: {
          name: 'targets',
          summary: '対象者',
          list: details.targets,
          liststr: this.toListString(details.targets),
          editable: true
        },
        administrators: {
          name: 'administrators',
          summary: '管理者',
          list: details.administrators,
          liststr: this.toListString(details.administrators),
          editable: true
        },
        respondents: {
          name: 'respondents',
          summary: '回答済みの人',
          list: details.respondents.filter((user, index, array) => {
            // 重複除去
            return array.indexOf(user) === index
          }),
          liststr: this.toListString(
            details.respondents.filter((user, index, array) => {
              // 重複除去
              return array.indexOf(user) === index
            })
          ),
          editable: false
        }
      }
    }
    return {}
  },

  noErrors(errors) {
    // 送信できるかどうかを返す
    const keys = Object.keys(errors)
    for (let i = 0; i < keys.length; i++) {
      if (errors[keys[i]].isError) {
        return false
      }
    }
    return true
  }
}
