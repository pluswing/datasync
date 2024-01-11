# datasync

データ同期するツール
開発向け。

- mysql
  - ほかDB (postgres ...)
- data (file)

- s3 or google storage

datasync push # データ送る
datasync pull # データ受け取る
datasync apply # データ適用

設定ファイル
datasync.yaml
```yaml
target:
  kind: mysql
  ...
upload:
  kind: s3
  ...
```

各種フレームワークに対しては、糖衣構文的な形で、configレスみたいな形にしたい。
```yaml
target:
  kind: rails
```

.datasync
```jsonl
{"hash code": "", "timestamp": "", "comment": "...."}
{"hash code": "", "timestamp": "", "comment": "...."}
{"hash code": "", "timestamp": "", "comment": "...."}
```
hash code.zip
hash code.zip
hash code.zip


# これからやること

- バグ修正
  - mysqlファイルの対応 (済) 未テスト
- postgresへの対応
- mysql, postgresコマンドを使う設定を追加
  - mysqldump, pg_dumpコマンドが入っている前提でそれを使う設定
  -  use: "native"
- rmのremote対応
- aws, samba対応

- init サブコマンド
  - bubble tea 使う

- mysql importエラー
  - packet for query is too large. Try adjusting the `Config.MaxAllowedPacket`
  - set global max_allowed_packet = 2 * 1024 * 1024 * 1024
    - で回避可能かと思ったけどそうもいかない模様。
    - mysqldumpの挙動だと、↑の設定に応じて、insert文が分割される。。
