import Vue from 'vue'
import store from '@/store'
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

const router = new Router({
  mode: 'history',
  routes: [
    {
      path: '/',
      redirect: '/targeted'
    },
    {
      path: '/targeted',
      name: 'Targeted',
      component: Targeted
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
      component: QuestionnaireDetails
    },
    {
      path: '/results/:id',
      name: 'Results',
      component: Results
    },
    {
      path: '/responses/:id',
      name: 'ResponseDetails',
      component: ResponseDetails
    },
    {
      path: '/responses/new/:questionnaireId',
      name: 'NewResponseDetails',
      component: ResponseDetails,
      props: { isNewResponse: true }
    },
    {
      path: '*',
      name: 'NotFound',
      component: NotFound
    }
  ],
  scrollBehavior (savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      // ページ遷移の時ページスクロールをトップに
      return { x: 0, y: 0 }
    }
  }
})

router.beforeEach(async (to, _, next) => {
  // traQにログイン済みかどうか調べる
  if (!store.state.me) {
    await store.dispatch('whoAmI')
  }

  if (!store.state.me) {
    // 未ログインの場合、traQのログインページに飛ばす
    const traQLoginURL = 'https://q.trap.jp/login?redirect=' + location.href
    location.href = traQLoginURL
  }

  next()
})

export default router
