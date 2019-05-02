import Vue from 'vue'
import Router from 'vue-router'
import Targeted from '@/pages/Targeted'
import Administrates from '@/pages/Administrates'
import Responses from '@/pages/Responses'
import Explorer from '@/pages/Explorer'
import QuestionnaireDetails from '@/pages/QuestionnaireDetails'
import Results from '@/pages/Results'
import ResponseDetails from '@/pages/ResponseDetails'
import NotFound from '@/pages/NotFound'

Vue.use(Router)

export default new Router({
  mode: 'history',
  props: {
    traqId: String
  },
  routes: [
    {
      path: '/',
      redirect: '/targeted'
    },
    {
      path: '/targeted',
      name: 'Targeted',
      component: Targeted,
      props: { traqId: String(this.traqId) }
    },
    {
      path: '/administrates',
      name: 'Administrates',
      component: Administrates
    },
    {
      path: '/responses',
      name: 'Responses',
      component: Responses
    },
    {
      path: '/explorer',
      name: 'Explorer',
      component: Explorer
    },
    {
      path: '/questionnaires/:id',
      name: 'QuestionnaireDetails',
      component: QuestionnaireDetails,
      props: true
    },
    {
      path: '/results/:id',
      name: 'Results',
      component: Results,
      props: true
    },
    {
      path: '/responses/:id',
      name: 'ResponseDetails',
      component: ResponseDetails,
      props: true
    },
    {
      path: '/responses/new/:questionnaireId',
      name: 'NewResponseDetails',
      component: ResponseDetails,
      props: { default: true, isNewResponse: true }
    },
    {
      path: '*',
      name: 'NotFound',
      component: NotFound
    }
  ],
  scrollBehavior (to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      // ページ遷移の時ページスクロールをトップに
      return { x: 0, y: 0 }
    }
  }
})
