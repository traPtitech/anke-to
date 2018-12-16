import Vue from 'vue'
import Router from 'vue-router'
// import HelloWorld from '@/components/HelloWorld'
import MypageLayout from '@/components/MypageLayout'
import CreatedLayout from '@/components/CreatedLayout'
import AnsweredLayout from '@/components/AnsweredLayout'
import ExplorerLayout from '@/components/ExplorerLayout'
import NotFound from '@/components/Utils/NotFound'

Vue.use(Router)

export default new Router({
  mode: 'history',
  props: {
    traqId: String
  },
  routes: [
    {
      path: '/',
      name: 'MypageLayout',
      component: MypageLayout,
      props: { traqId: String(this.traqId) }
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
