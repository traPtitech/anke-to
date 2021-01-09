# anke-to
[![codecov](https://codecov.io/gh/traPtitech/anke-to/branch/master/graph/badge.svg)](https://codecov.io/gh/traPtitech/anke-to)
[![](https://github.com/traPtitech/anke-to/workflows/Release/badge.svg?branch=release)](https://github.com/traPtitech/anke-to/actions)
[![swagger](https://img.shields.io/badge/swagger-docs-brightgreen)](https://traptitech.github.io/anke-to/swagger/)
[![go report](https://goreportcard.com/badge/traPtitech/anke-to)](https://goreportcard.com/report/traPtitech/anke-to)

部内アンケートシステム

## 開発
https://wiki.trapti.tech/SysAd/docs/anke-to/development
### サーバーサイド
Docker, Goが必要です
```
make dev
```

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