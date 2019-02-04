<template>
  <div class="columns">
    <article class="column is-11">
      <div class="card">
        <!-- タイトル、説明、回答期限 -->
        <div>
          <header class="card-header">
            <div id="title" class="card-header-title title">
              <span
                v-if="details.responseIconClass"
                :class="details.responseIconClass"
                class="response-icon"
              ></span>
              <router-link v-if="details.titleLink" :to="details.titleLink">{{ details.title }}</router-link>
              <p v-if="!details.titleLink">{{ details.title }}</p>
            </div>
          </header>
          <div v-if="details.description" class="card-content">
            <pre>{{ details.description }}</pre>
          </div>
          <div v-if="details.timeLimit" class="is-pulled-right wrapper editable time-limit">
            <span class="label">回答期限 :</span>
            <span class="time-limit">{{ details.timeLimit }}</span>
          </div>
          <div v-if="details.responseDetails" class="response-details">
            <div class="is-pulled-left is-inline-block wrapper editable">
              <span class="respondent subtitle is-5">{{ details.responseDetails.respondent }}</span>
            </div>
            <div class="is-pulled-right is-flex editable">
              <span class="label">回答日時 :</span>
              <span class="time-limit">{{ details.responseDetails.time }}</span>
            </div>
          </div>
        </div>
      </div>
    </article>
  </div>
</template>

<script>

// import <componentname> from '<path to component file>'

export default {
  name: '',
  components: {
  },
  props: {
    details: {
      type: Object,
      required: true
    }
  },
  data () {
    return {
    }
  },
  methods: {
  },
  computed: {
  },
  mounted () {
  }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style lang="scss" scoped>
.card {
  max-width: 100%;
  padding: 0.7rem;
}
.columns {
  margin: 0;
}
.column {
  padding: 0;
}
.editable {
  span {
    width: fit-content;
    height: fit-content;
    top: 0;
    bottom: 0;
    margin: auto 0.2rem;
    white-space: nowrap;
  }
  &.wrapper {
    display: flex;
  }
}
.response-icon {
  font-size: large;
  margin-right: 0.5rem;
}
.response-details,
.time-limit {
  margin: 1rem;
}
pre {
  white-space: pre-line;
  font-size: inherit;
  -webkit-font-smoothing: inherit;
  font-family: inherit;
  line-height: inherit;
  background-color: inherit;
  color: inherit;
  padding: 0.625em;
}
</style>
