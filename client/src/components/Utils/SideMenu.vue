<template>
  <div class="submenu column is-2 has-background-light is-fullheight">
    <aside class="nav menu">
      <ul class="menu-list">
        <li v-for="(menuItem, index) in menuItems" :key="index" @click="$emit('close-side-menu')">
          <router-link :to="menuItem.url">{{ menuItem.name }}</router-link>
        </li>
        <hr>
        <li @click.prevent="createQuestionnaire">
          <a>New Questionnaire</a>
        </li>
      </ul>
    </aside>
  </div>
</template>
<script>
import router from '@/router'

export default {
  name: 'SideMenu',
  props: {
    traqId: {
      type: String,
      required: true
    }
  },
  data () {
    return {
      menuItems: [
        {
          name: 'Targeted',
          url: '/targeted'
        },
        {
          name: 'Administrates',
          url: '/administrates'
        },
        {
          name: 'Responses',
          url: '/responses'
        },
        {
          name: 'Explorer',
          url: '/explorer'
        }
      ]
    }
  },
  methods: {
    createQuestionnaire () {
      this.$emit('close-side-menu')
      router.push('/questionnaires/new')
    }
  }
}
</script>

<style lang="scss" scoped>
.submenu {
  min-width: fit-content;
}
.menu-list {
  hr {
    border-style: outset;
    margin: 0.3rem 0;
  }
  .button {
    margin: 0.5rem;
  }
}
.menu-list a {
  position: relative;
  display: inline-block;
  text-decoration: none;
}
.menu-list a::after {
  position: absolute;
  bottom: 3.5px;
  left: 0.4em;
  content: "";
  width: 100%;
  height: 2px;
  background: #333;
  transform: scale(0, 1);
  transform-origin: right top;
  transition: transform 0.3s;
}
.menu-list :hover::after {
  transform-origin: left top;
  transform: scale(1, 1);
}
</style>
