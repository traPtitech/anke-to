<template>
  <div>
    <div class="columns">
      <article class="column is-11">
        <div class="card">
          <!-- タイトル、説明、回答期限 -->
          <div>
            <header class="card-header">
              <div class="card-header-title title" :class="{'is-editing' : isEditing}">
                <input v-show="isEditing" id="title" v-model="details.title" class="input">
                <div v-show="!isEditing">{{ details.title }}</div>
              </div>
            </header>
            <div class="card-content">
              <textarea
                v-show="isEditing"
                id="description"
                v-model="details.description"
                class="textarea"
                rows="5"
              ></textarea>
              <pre v-show="!isEditing">{{ details.description }}</pre>
            </div>
            <div class="has-text-right">回答期限 : {{ props.getDateStr(details.res_time_limit) }}</div>
          </div>
        </div>
      </article>
    </div>

    <div class="columns details">
      <article class="column is-4">
        <div class="card">
          <!-- 操作 -->
          <div>
            <header class="card-header">
              <div class="card-header-title subtitle">操作</div>
            </header>
            <div class="card-content management-buttons">
              <a class="button" :href="questionnaireId + '/new-response'">新しい回答を作成</a>
              <a
                class="button"
                :disabled="!canViewResults"
                :class="{'disabled' : !canViewResults}"
                :href="'/results/' + questionnaireId"
              >結果を見る</a>
            </div>
          </div>
        </div>
      </article>

      <article class="column is-7">
        <div class="card">
          <!-- 自分の回答一覧 -->
          <div>
            <header class="card-header">
              <div class="card-header-title subtitle">自分の回答</div>
            </header>
            <div class="card-content">このアンケートに対する自分の回答一覧</div>
          </div>
        </div>

        <div class="card">
          <!-- 情報 -->
          <div>
            <header class="card-header">
              <div class="card-header-title subtitle">情報</div>
            </header>
            <div class="card-content">
              <div class="has-text-weight-bold">
                <div>更新日時 : {{ props.getDateStr(details.modified_at) }}</div>
                <div>作成日時 : {{ props.getDateStr(details.created_at) }}</div>
              </div>
              <details v-for="(userList, index) in userLists" :key="index">
                <summary>{{ userList.summary }}</summary>
                <p class="has-text-grey">{{ toListString(userList.list) }}</p>
              </details>
              <div class="has-text-weight-bold">
                <div>結果の公開範囲: {{ resSharedToStr }}</div>
              </div>
            </div>
          </div>
        </div>
      </article>
    </div>
    {{ details }}
  </div>
</template>

<script>

// import <componentname> from '<path to component file>'
import axios from '@/bin/axios'

export default {
  name: 'Information',
  components: {
  },
  async created () {
    const resp = await axios.get('/questionnaires/' + this.questionnaireId)
    this.details = resp.data
    if (this.administrates) {
      this.$emit('enable-edit-button')
    }
  },
  props: {
    props: {
      type: Object,
      required: true
    }
  },
  data () {
    return {
      details: {}
    }
  },
  methods: {
    toListString (list) {
      if (typeof list === 'undefined' || list.length === 0) {
        return ''
      }
      let ret = ''
      for (let i = 0; i < list.length - 1; i++) {
        ret += list[ i ] + ', '
      }
      ret += list[ list.length - 1 ]
      return ret
    }
  },
  computed: {
    traqId () {
      return this.props.traqId
    },
    getDateStr () {
      return this.props.getDateStr
    },
    isEditing () {
      return this.props.isEditing
    },
    questionnaireId () {
      return this.$route.params.id
    },
    administrates () {
      // 管理者かどうかを返す
      if (typeof this.details.administrators !== 'undefined') {
        for (let i = 0; i < this.details.administrators.length; i++) {
          if (this.props.traqId === this.details.administrators[ i ]) {
            return true
          }
        }
      }
      return false
    },
    canViewResults () {
      // 結果をみる権限があるかどうかを返す
      return ((this.details.res_shared_to === 'public') ||
        (this.details.res_shared_to === 'administrators' && this.administrates) ||
        (this.details.res_shared_to === 'respondents' && this.responses.length > 0))
    },
    responses () {
      // このアンケートに対する自分の回答一覧を返す 未実装
      let ret = []
      return ret
    },
    userLists () {
      return [
        {
          summary: '対象者',
          list: this.details.targets
        },
        {
          summary: '回答済みの人',
          list: this.details.respondents
        },
        {
          summary: '管理者',
          list: this.details.administrators
        }
      ]
    },
    resSharedToStr () {
      switch (this.details.res_shared_to) {
        case 'public': return '全体'
        case 'administrators': return '管理者のみ'
        case 'respondents': return '回答済みの人'
      }
    }
  },
  mounted () {
  }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style lang="scss" scoped>
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
.card {
  max-width: 100%;
  padding: 0.7rem;
}
article.column {
  padding: 0;
}
.columns {
  margin-bottom: 0;
}
.columns:first-child {
  display: flex;
}

.card-header-title.is-editing {
  padding: 0;
}
.card-content {
  .subtitle {
    margin: 0;
  }
  > details {
    margin: 0.5rem;
    > p {
      padding: 0 0.5rem;
    }
  }
}
.management-buttons {
  > .button:not(:last-child) {
    margin-bottom: 0.7rem;
  }
  > .button {
    max-width: fit-content;
    display: block;
  }
}
#title.input {
  font-size: 2rem;
}
@media screen and (min-width: 769px) {
  // widthが大きいときは横並びのカードの間を狭くする
  .column:not(:last-child) > .card {
    margin-right: 0;
  }
}
</style>
