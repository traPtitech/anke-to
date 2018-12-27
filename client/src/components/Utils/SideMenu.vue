<template>
  <div class="submenu column is-2 has-background-light is-fullheight">
    <aside class="nav menu">
      <ul class="menu-list">
        <li v-for="(menuItem, index) in menuItems" :key="index">
          <a :href="menuItem.url">{{ menuItem.name }}</a>
        </li>
      </ul>
    </aside>
  </div>
</template>
<script>
// import axios from '@/bin/axios'

export default {
  name: 'SideMenu',
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
  }
}
</script>

<style scoped>
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
  width: 5em;
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
