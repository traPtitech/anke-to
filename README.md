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
#### ツールのインストール
```
make init
```
#### 開発
```
make dev
```
#### テスト
```
make test
```
注意：本サービスはユーザー認証機能を持ちません。リバースプロキシなどを利用して、外部の認証サービスで取得したユーザーIDをHTTPヘッダーのX-Forwarded-UserにユーザーIDを設定した上で、本サービスにリクエストを転送してください