/* eslint-disable */
export default {
  setContent: function(label, value) {
    this.$emit('set-question-content', this.questionIndex, label, value)
  }
}
