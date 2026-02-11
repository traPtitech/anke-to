# パッケージ構成

```
.
├── controller/    ビジネスロジック
├── handler/    echoのハンドラー
├── model/    dbからのデータの取り出し
│   └── mock_model/    modelのmockgenによるmock。直接編集してはいけない。
├── openapi/    oapi-codegenによって生成されたコード。直接編集してはいけない。
├── traq    traQとの通信関連
│   └── mock_traq/    traqのmockgenによるmock。直接編集してはいけない。
├── main.go
├── middleware.go    ミドルウェアの切替
├── wire.go    wireによるDI
└── wire_gen.go    wireによる生成コード。直接編集してはいけない。
```