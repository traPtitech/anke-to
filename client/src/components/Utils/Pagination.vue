<template>
  <nav class="navbar is-fixed-bottom pagination is-rounded is-centered">
    <ul class="pagination-list">
      <li>
        <router-link
          class="pagination-previous"
          :disabled="disableFirstButton"
          :to="getPageLink('first')"
        >
          <span class="ti-angle-double-left"></span>
        </router-link>
      </li>
      <li v-for="pageNum in pages" :key="pageNum">
        <router-link
          class="pagination-link"
          v-if="pageNum"
          :class="{'is-current' : pageNum===currentPage}"
          :to="getPageLink(pageNum)"
        >{{ pageNum }}</router-link>
        <span v-if="!pageNum" class="pagination-link" disabled></span>
      </li>
      <li>
        <router-link
          class="pagination-next"
          :disabled="disableLastButton"
          :to="getPageLink('last')"
        >
          <span class="ti-angle-double-right"></span>
        </router-link>
      </li>
    </ul>
  </nav>
</template>

<script>

// import <componentname> from '<path to component file>'

export default {
  name: 'Pagination',
  components: {
  },
  props: {
    range: {
      type: Object,
      required: true
    },
    currentPage: {
      type: Number,
      required: true
    },
    getPageLink: {
      type: Function,
      required: true
    }
  },
  data () {
    return {
      paginationWidth: 1
    }
  },
  methods: {
  },
  computed: {
    pages () {
      let ret = []
      for (let i = this.currentPage - this.paginationWidth; i <= this.currentPage + this.paginationWidth; i++) {
        if (i >= this.range.first && i <= this.range.last) {
          ret.push(i)
        } else {
          ret.push(undefined)
        }
      }
      return ret
    },
    disableFirstButton () {
      return this.currentPage - this.paginationWidth <= this.range.first
    },
    disableLastButton () {
      return this.currentPage + this.paginationWidth >= this.range.last
    }
  },
  mounted () {
  }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style lang="scss" scoped>
.pagination {
  margin: 0.5rem 0;
}
.pagination-next,
.pagination-previous {
  margin: 0.25rem 0.75rem;
}
.pagination-link[disabled] {
  cursor: default;
}
</style>
