# anke-to
[![codecov](https://codecov.io/gh/traPtitech/anke-to/branch/main/graph/badge.svg)](https://codecov.io/gh/traPtitech/anke-to)
[![](https://github.com/traPtitech/anke-to/workflows/Release/badge.svg?branch=release)](https://github.com/traPtitech/anke-to/actions)
[![swagger](https://img.shields.io/badge/swagger-docs-brightgreen)](https://apis.trap.jp/?urls.primaryName=anke-to)
[![go report](https://goreportcard.com/badge/traPtitech/anke-to)](https://goreportcard.com/report/traPtitech/anke-to)

部内アンケートシステム

## 開発
https://wiki.trapti.tech/SysAd/docs/anke-to/development
### サーバーサイド
Docker, Goが必要です
```
make dev
```
注意：本サービスはユーザー認証機能を持ちません。リバースプロキシなどを利用して、外部の認証サービスで取得したユーザーIDをHTTPヘッダーのX-Showcase-UserにユーザーIDを設定した上で、本サービスにリクエストを転送してください。

#### ベンチマーク
Docker,openapi-generator-cli,Goが必要です。
```
#ベンチマーク前のanke-to起動
# make tuning

#ベンチマーク
$ make bench

#750レコードinsert
$ make bench-init

#pprof
$ make pprof

#pt-query-digest
# make slow

#myprofiler
# make myprof ARGS="{引数}"
```

### クライアントサイド
Node.js が必要です
```
cd client
npm install
npm run serve
```

(詳しくは `client/README.md` を参照)