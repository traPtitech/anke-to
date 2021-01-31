# パッケージ構成

```
.
├── main.go
├── model/    dbからのデータの取り出し
│   └── mock_model/    modelのmockgenによるmock。直接編集してはいけない。
├── router/    echoのハンドラー・ビジネスロジック
├── router.go    echo routerの定義
├── traq    traQとの通信関連
│   └── mock_traq/    traqのmockgenによるmock。直接編集してはいけない。
├── wire.go    wireによるDI
└── wire_gen.go    wireによる生成コード。直接編集してはいけない。
```