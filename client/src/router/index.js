import Vue from 'vue'
import Router from 'vue-router'
// import HelloWorld from '@/components/HelloWorld'
import MypageLayout from '@/components/Mypage/MypageLayout'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'MypageLayout',
      component: MypageLayout
    }
  ]
})
