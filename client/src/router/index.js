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
import Blank from '@/pages/Blank'
import { getRequest2Callback, redirect2AuthEndpoint } from '@/util/api.js'

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
      path: '/questionnaires/new',
      name: 'QuestionnaireDetailsNew',
      component: QuestionnaireDetails,
      meta: {
        requiresTraqAuth: true
      }
    },
    {
      path: '/questionnaires/:id',
      name: 'QuestionnaireDetails',
      component: QuestionnaireDetails
    },
    {
      path: '/questionnaires/:id/edit',
      name: 'QuestionnaireDetailsEdit',
      component: QuestionnaireDetails,
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
    },
    {
      // traQトークン取得後のコールバックURL
      path: '/callback',
      name: 'Callback',
      component: Blank,
      beforeEnter: async (to, _, next) => {
        await getRequest2Callback(to)
        const destination = sessionStorage.getItem('destination')
        if (destination) {
          next(destination)
        }
        next()
      }
    }
  ],
  scrollBehavior(savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      // ページ遷移の時ページスクロールをトップに
      return { x: 0, y: 0 }
    }
  }
})

router.beforeEach(async (to, from, next) => {
  console.log(to.name)
  if (to.name === 'Callback') {
    next()
    return
  }
  // traQにログイン済みかどうか調べる
  if (!store.state.me) {
    await store.dispatch('whoAmI')
  }

  if (!store.state.me) {
    // 未ログインの場合、traQのログインページに飛ばす
    sessionStorage.setItem(`destination`, to.fullPath)
    await redirect2AuthEndpoint()
  }

  next()
})

export default router
