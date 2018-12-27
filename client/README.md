# client

> A Vue.js project

## Build Setup

```bash
# install dependencies
npm install

# serve with hot reload at localhost:8080
npm run dev

# build for production with minification
npm run build

# build for production and view the bundle analyzer report
npm run build --report

# run unit tests
npm run unit

# run all tests
npm test
```

For a detailed explanation on how things work, check out the [guide](http://vuejs-templates.github.io/webpack/) and [docs for vue-loader](http://vuejs.github.io/vue-loader).

## URL

- `/` : 自分が対象になっているアンケート一覧
- `/administrates` : 自分が管理者になっているアンケート一覧
- `/responses` : 自分の回答一覧
- `/explorer` : すべてのアンケート一覧
- `/questionnaires/:id` : questionnaireID = id のアンケートの詳細
  - 編集のときは `#edit-form` を末尾につける
- `/results/:id` : questionnaireID = id のアンケートの結果
- `/responses/:id` : responseID = id の回答
  - 編集のときは `#edit-form` を末尾につける
