import Vue from 'vue'
import Router from 'vue-router'
// import HelloWorld from '@/components/HelloWorld'
import MypageLayout from '@/components/Mypage/MypageLayout'
import CreatedLayout from '@/components/Created/CreatedLayout'
import AnsweredLayout from '@/components/Answered/AnsweredLayout'
import ExplorerLayout from '@/components/Explorer/ExplorerLayout'
import NotFound from '@/components/Main/NotFound'

Vue.use(Router)

export default new Router({
  mode: 'history',
  routes: [
    {
      path: '/',
      name: 'MypageLayout',
      component: MypageLayout
    },
    {
      path: '/created',
      name: 'Created',
      component: CreatedLayout
    },
    {
      path: '/answered',
      name: 'Answered',
      component: AnsweredLayout
    },
    {
      path: '/explorer',
      name: 'Explorer',
      component: ExplorerLayout
    },
    {
      path: '*',
      name: 'NotFound',
      component: NotFound
    }
  ]
})
