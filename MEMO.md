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

