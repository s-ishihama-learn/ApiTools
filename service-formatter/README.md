# XML / HTML Formatter (service-formatter)

XMLおよびHTMLテキストに対して、自動で改行やインデントを設定して書式を整えるツール・Web APIサービスです[cite: 1]。

---

## 開発環境・依存ライブラリ

モジュールの初期化と使用ライブラリの取得手順です[cite: 1]。

```bash
go mod init service-formatter

go get [github.com/go-xmlfmt/xmlfmt](https://github.com/go-xmlfmt/xmlfmt)
go get [github.com/labstack/echo/v4](https://github.com/labstack/echo/v4)

```

---

## 実行・利用方法

本プロジェクトには「コマンドライン型」と「Web API型」の2種類の実装が含まれています。

### 1. コマンドライン型

CLI環境で直接ファイルを整形・出力します。

* **構成ファイル**: `main.go`, `xmlFormatter.go`


#### デバッグ実行

```bash
# 標準出力へ確認
go run main.go xmlFormatter.go testdata/index-test.html

# ファイルへ出力
go run main.go xmlFormatter.go testdata/index-test.html testdata/index-test_format.html
go run main.go xmlFormatter.go testdata/sample.xml

```

#### コンパイル & 実行

```bash
go build -o XMLformatter.exe main.go xmlFormatter.go
./XMLformatter.exe testdata/index-test.html

```

---

### 2. Web API型

HTTP経由で整形処理を提供するサービスです。

* **構成ファイル**: `XmlFormatService.go`, `xmlFormatter.go`


#### Linux向けビルド

```cmd
set GOOS=linux
set GOARCH=amd64
go build -o xml-formatter.exe XmlFormatService.go xmlFormatter.go

```

#### デプロイ手順

```bash
# REST API (バイナリの配置とサービス再起動)
scp xml-formatter.exe beachstone@www.logware.jp:/home/beachstone/api/xml-formatter.exe
# サーバー側作業:
# chmod 755 xml-formatter.exe
# systemctl restart format.service

# HTML (フロントエンド画面の配置)
scp index.html beachstone@www.logware.jp:/tmp/.
# サーバー側作業:
# cp /tmp/index.html /var/www/html/formatter/.

```

#### Web API 動作確認 (curl)

```bash
# ファイル送信フォーム形式
curl -X POST http://localhost:6100/format/file -F "file=@index-test.html" -o index-testF.html

# JSON形式 (ローカル)
curl -X POST http://localhost:6100/format/json -H "Content-Type: application/json" -d '{"content": "<html><body><script>console.log(\"hello\");</script></body></html>"}'

# JSON形式 (本番環境)
curl -X POST [http://www.logware.jp/format/json](http://www.logware.jp/format/json) -H "Content-Type: application/json" -d '{"content": "<html><body><script>console.log(\"hello\");</script></body></html>"}'

```

#### Web画面動作確認

* **URL**: [http://www.logware.jp/formatter/](https://www.google.com/search?q=http://www.logware.jp/formatter/&utm_source=gemini)

* **テストデータ例**:
```html
<html><body><h1>Hello World!<br>(Type HTML)</h1><script>console.log("hello");</script></body></html>

```



---

## 保守・ログ確認

障害発生時や動作確認時のログ出力コマンドです。

```bash
tail -20 /var/log/messages
tail -20 /var/log/nginx/error.log
tail -20 /var/log/nginx/access.log

```

---

## 仕様・注意事項

* **対応フォーマット**: XML および HTML の両方に対応しています。


* **JavaScriptコメントの扱い**:
* HTML内のJavaScriptにおいて、コメント行の判定が難しい場合に後続のコードまでコメント化されてしまう制限があります。


* **対策**: コメント行の末尾に `;`（セミコロン）を記述すると、正しく改行され独立したコメント行として認識されます。




* **`<br>` および `<hr>` タグの内部処理**:
* 閉じタグのない `<br>` や `<hr>` はそのまま処理するとインデント崩れの原因となります。


* 処理内部で一時的に `<br/>` や `<hr/>` へ変換してインデント調整を行い、処理後に元の形式へ戻す仕様となっています。


* そのため、入力データ内で `<br>` と `<br/>`（または `<hr>` と `<hr/>`）が混在している場合は、全て `<br>`（`<hr>`）へ一括で統一出力されます。
