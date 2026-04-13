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

## 必要な環境変数
```
ENV：
PORT: :
MARIADB_USERNAME: root
MARIADB_PASSWORD: password
MARIADB_HOSTNAME: 127.0.0.1
MARIADB_DATABASE: anke-to
MARIADB_PORT: 3306
TRAQ_BOT_TOKEN: ""
TRAQ_WEBHOOK_ID: ""
TRAQ_WEBHOOK_SECRET: ""
```

### 環境変数
- `ENV`：実行環境。`ENV == production` のときはログレベルが異なります。`ENV == neoshowcase` のときは NeoShowcase でデプロイするため、DB 関連の変数名が変わります
- `PORT`：サービスのポート（例：`:1323`）
- `MARIADB_USERNAME`：データベースのユーザー名。`ENV == neoshowcase` のときは `NS_MARIADB_USER`
- `MARIADB_PASSWORD`：データベースのパスワード。`ENV == neoshowcase` のときは `NS_MARIADB_PASSWORD`
- `MARIADB_HOSTNAME`：データベースのホスト名または IP。`ENV == neoshowcase` のときは `NS_MARIADB_HOSTNAME`
- `MARIADB_PORT`：データベースのポート。`ENV == neoshowcase` のときは `NS_MARIADB_PORT`
- `MARIADB_DATABASE`：データベース名。`ENV == neoshowcase` のときは `NS_MARIADB_DATABASE`
- `TRAQ_BOT_TOKEN`：traQ API の認証トークン（未使用時は空で可）
- `TRAQ_WEBHOOK_ID`：traQ Webhook の Client ID（未使用時は空で可）
- `TRAQ_WEBHOOK_SECRET`：traQ Webhook の Client Secret（未使用時は空で可）