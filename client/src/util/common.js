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
    // console.log(data)
    // console.log(this.questionTypes)
    // console.log(data.question_type)
    let question = {
      questionId: data.questionID,
      type: data.question_type,
      component: this.questionTypes[data.question_type].component,
      questionBody: data.body,
      isRequired: data.is_required
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
    // console.log(responseData)
    let question = Object.assign({}, questionData)
    switch (question.type) {
      case 'Text':
        question.responseBody = responseData.response
        break
      case 'Number':
        question.responseBody = Number(responseData.response)
        break
      case 'Checkbox':
        responseData.option_response.forEach(selectedOption => {
          question.isSelected[selectedOption] = true
        })
        break
      case 'MultipleChoice':
        question.selected = responseData.option_response[0]
        break
      case 'LinearScale':
        console.log(responseData)
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
  }
}
