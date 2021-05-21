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
import { sendTokenRequest, sendCodeRequest } from '../bin/traqAuth'

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
      meta: {
        requiresTraqAuth: true
      }
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
        const clearSessionStorage = () => {
          sessionStorage.removeItem('nextRoute')
          sessionStorage.removeItem('previousRoute')
          sessionStorage.removeItem(`traq-auth-code-verifier-${state}`)
        }

        const code = to.query.code
        const state = to.query.state
        const codeVerifier = sessionStorage.getItem(
          `traq-auth-code-verifier-${state}`
        )
        if (!code || !codeVerifier) {
          let previousRoute = sessionStorage.getItem('previousRoute')
          if (!previousRoute) previousRoute = '/targeted'
          clearSessionStorage()
          next(previousRoute)
          return
        }

        const res = await sendTokenRequest(code, codeVerifier)
        store.commit('traq/setAccessToken', res.data.access_token)

        let nextRoute = sessionStorage.getItem('nextRoute')
        if (!nextRoute) nextRoute = '/targeted'
        clearSessionStorage()
        next(nextRoute)
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
  // traQにログイン済みかどうか調べる
  if (!store.state.me) {
    await store.dispatch('whoAmI')
  }

  if (!store.state.me) {
    // 未ログインの場合、traQのログインページに飛ばす
    const traQLoginURL = 'https://q.trap.jp/login?redirect=' + location.href
    location.href = traQLoginURL
  }

  if (to.meta.requiresTraqAuth) {
    await store.dispatch('traq/ensureToken')
    if (!store.state.traq.accessToken) {
      const message =
        'アンケートの編集・作成にはtraQアカウントへのアクセスが必要です。OKを押すとtraQに飛びます。'
      if (window.confirm(message)) {
        sessionStorage.setItem('nextRoute', to.path) // traQでのトークン取得後に飛ばすルート
        sessionStorage.setItem('previousRoute', from.path) // traQでのトークン取得失敗時に飛ばすルート
        await sendCodeRequest()

        // traQのconsentページに飛ぶ前にnextが表示されることを防ぐ
        next(false)
        return
      } else {
        // キャンセルを押された場合は元のルートに戻る
        if (from.path !== to.path) {
          next(from.path)
        } else {
          // url直打ちなどでアクセスされた場合
          next('/targeted')
        }
        return
      }
    }
  }

  next()
})

export default router
