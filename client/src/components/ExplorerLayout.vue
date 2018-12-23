<template>
  <div class="wrapper is-fullheight">
    <div class="dropdowns">
      <div class="dropdown" :class="{ 'is-active': DropdownIsActive.sortOrder }">
        <div class="dropdown-trigger">
          <button
            class="button"
            aria-haspopup="true"
            aria-controls="dropdown-menu"
            @click="DropdownIsActive.targetedOption = false; DropdownIsActive.sortOrder = !DropdownIsActive.sortOrder"
          >
            <p>並び替え</p>
            <span class="ti-angle-down"></span>
          </button>
        </div>
        <div class="dropdown-menu" id="dropdown-menu" role="menu">
          <div
            class="dropdown-content"
            v-for="(order, index) in sortOrders"
            :key="index"
            :class="{'is-selected' : order.opt===sortOrder}"
            @click="changeSortOrder(order.opt); DropdownIsActive.sortOrder = false"
          >
            <p class="dropdown-item">{{ order.str }}</p>
          </div>
        </div>
      </div>
      <div class="dropdown" :class="{ 'is-active': DropdownIsActive.targetedOption }">
        <div class="dropdown-trigger">
          <button
            class="button"
            aria-haspopup="true"
            aria-controls="dropdown-menu"
            @click="DropdownIsActive.sortOrder = false; DropdownIsActive.targetedOption = !DropdownIsActive.targetedOption"
          >
            <p>フィルター</p>
            <span class="ti-angle-down"></span>
          </button>
        </div>
        <div class="dropdown-menu" id="dropdown-menu" role="menu">
          <div
            class="dropdown-content"
            v-for="(option, index) in targetedOptions"
            :key="index"
            :class="{'is-selected' : option.opt===targetedOption}"
            @click="changetargetedOption(option.opt); DropdownIsActive.targetedOption = false"
          >
            <p class="dropdown-item">{{ option.str }}</p>
          </div>
        </div>
      </div>
    </div>
    <div
      class="card-wrapper is-fullheight"
      @click="DropdownIsActive.sortOrder = false; DropdownIsActive.targetedOption = false"
    >
      <div class="card">
        <!-- <button class="button" v-on:click="changeSortOrder('-title')">Button</button> -->
        <table class="table is-striped">
          <thead>
            <tr>
              <th v-for="(header, index) in headers" :key="index">{{ header }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(row, index) in itemrows" :key="index">
              <td class="table-item-title" v-html="row.title"></td>
              <td class="table-item-date">{{ row.res_time_limit }}</td>
              <td class="table-item-date">{{ row.modified_at }}</td>
              <td class="table-item-date">{{ row.created_at }}</td>
              <td v-html="row.resultsLinkHtml"></td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script>
import axios from '@/bin/axios'
import Table from '@/components/Utils/Table.vue'

export default {
  name: 'ExplorerLayout',
  components: {
    'customtable': Table
  },
  async created () {
    this.getQuestionnaires()
  },
  props: {
    traqId: {
      type: String,
      required: true
    }
  },
  data () {
    return {
      questionnaires: [],
      headers: ['', '回答期限', '更新日時', '作成日時', '結果'],
      sortOrder: '-modified_at',
      sortOrders: [
        {
          str: '最近更新された',
          opt: '-modified_at'
        },
        {
          str: '最近更新されていない',
          opt: 'modified_at'
        },
        {
          str: 'タイトル順',
          opt: 'title'
        },
        {
          str: 'タイトル逆順',
          opt: '-title'
        },
        {
          str: '最新',
          opt: '-created_at'
        },
        {
          str: '最も古い',
          opt: 'created_at'
        }
      ],
      targetedOption: false,
      targetedOptions: [
        {
          str: '全て',
          opt: false
        },
        {
          str: '対象外のもののみ',
          opt: true
        }
      ],
      DropdownIsActive: {
        sortOrder: false,
        targetedOption: false
      }
    }
  },
  computed: {
    itemrows () {
      let itemrows = []
      for (let i = 0; i < this.questionnaires.length; i++) {
        let row = {}
        row.title = this.getTitleHtml(i)
        row.res_time_limit = this.getDateStr(this.questionnaires[i].res_time_limit)
        row.modified_at = this.getDateStr(this.questionnaires[i].modified_at)
        row.created_at = this.getDateStr(this.questionnaires[i].created_at)
        row.resultsLinkHtml = this.getResultsLinkHtml(this.questionnaires[i].questionnaireID) // 結果を見る権限があるかどうかでボタンの色を変えたりしたい
        itemrows.push(row)
      }
      return itemrows
    }
  },
  methods: {
    getTitleHtml (i) {
      return '<a href="/questionnaires/' + this.questionnaires[i].questionnaireID + '">' + this.questionnaires[i].title + '</a>'
    },
    getResultsLinkHtml (id) {
      return '<a href="/resuslts/' + id + '">' + '<span class="ti-new-window"></span><br>Open' + '</a>'
    },
    getDateStr (str) {
      return str === 'NULL' ? '-' : new Date(str).toLocaleString()
    },
    getQuestionnaires () {
      axios
        .get('/questionnaires?sort=' + this.sortOrder + '&nontargeted=' + this.targetedOption)
        .then(response => (this.questionnaires = response.data))
        .catch(error => console.log(error))
    },
    changeSortOrder (sortOrder) {
      this.sortOrder = sortOrder
      this.getQuestionnaires()
    },
    changetargetedOption (targetedOption) {
      this.targetedOption = targetedOption
      this.getQuestionnaires()
    }

  }
}
</script>

<style scoped>
td {
  vertical-align: middle;
  font-size: 0.9em;
}
.table-item-title {
  min-width: 10em;
  font-size: 1em;
}
.table-item-date {
  min-width: 8em;
  text-align: center;
}
.dropdowns {
  padding-left: 1.5rem;
  padding-right: 1.5rem;
  padding-top: 1.5rem;
}
.dropdown-content:hover {
  background-color: lightgray;
}
.dropdown-content.is-selected {
  background-color: darkgray;
}
.button p {
  margin-right: 0.5em;
}
.card {
  margin: 0;
}
.card-wrapper {
  padding: 1rem 1.5rem;
}
</style>
