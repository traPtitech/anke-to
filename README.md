# anke-to

部内アンケートシステム

## 開発
https://wiki.trapti.tech/SysAd/docs/anke-to/development
### サーバーサイド
Dockerが必要です
```
docker-compose -f development/docker-compose.yaml up --build
```

### クライアントサイド
Node.js が必要です
```
cd client
npm run serve
```

(詳しくは `client/README.md` を参照)