# DataSync
[English](https://github.com/pluswing/datasync/blob/develop/README.md) [日本語](https://github.com/pluswing/datasync/blob/develop/README_ja.md)


開発者のためのデータベース共有ツール

## これはなに？
開発作業中にデータベースの内容をチームメンバーと共有すること、ありますよね？
DataSyncは、開発者が直面する一般的な問題を解決します。
データベースのバックアップ作成、共有、適用が面倒で時間がかかり、しばしば複雑です。
DataSyncを使用すると、これらのプロセスがシンプルかつ効率的になります。
例えば、新しい機能のテストやバグの再現に必要なデータを、チームメンバーに簡単に渡すことができます。
これにより、開発サイクルを高速化し、より生産的な活動に注力できます。

## 概要
DataSyncは、mysqlおよびファイルのバックアップ、履歴管理、そして任意のバックアップの迅速な適用を可能にするツールです。

主要な機能には以下のものがあります：

- データベースとファイルのバックアップ
- バックアップの履歴管理
- 任意のバックアップの適用
- クラウドストレージとの連携によるバックアップの簡単な共有と取得

## 使い方の例

### バックアップ実行
```
$ datasync dump -m "feature_test"
✔️ mysql dump completed (database: sample)
✔︎ compress data completed.
Dump succeeded. Version ID = 35ca8d497d334891b2ff627174a2b88a
```

### バックアップの一覧を表示
```
$ datasync ls -a
-- Remote versions --
224120fe68d14f6eaf2b4ea0533c497f 2024-01-30 13:53:21 test001
-- local versions --
35ca8d497d334891b2ff627174a2b88a 2024-02-10 09:33:27 test002
```

### バックアップの適用
```
$ datasync apply 35ca8d4
✔︎ decompress data completed.
✔︎ mysql import completed (database: sample)
Apply succeeded. Version ID = 35ca8d497d334891b2ff627174a2b88a
```

### クラウドへバックアップを送信
```
$ datasync push
```

### クラウドからバックアップを取得
```
$ datasync pull
```

## インストール
DataSyncは依存のない1ファイルバイナリです。以下の手順に従ってください：

1. [Releasesページ](https://github.com/pluswing/datasync/releases)から最新のDataSyncをダウンロードします。
2. ダウンロードしたファイルを任意の場所に保存します。
3. コマンドラインから、保存した場所に移動し、以下のコマンドを実行してDataSyncを初期化します：
```
datasync init
```
これで、準備完了です！

## License
MIT License
