package main

import (
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-xmlfmt/xmlfmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// リクエスト用構造体
type FormatRequest struct {
	Content string `json:"content"`
}

// レスポンス用構造体
type FormatResponse struct {
	FormattedContent string `json:"formatted_content"`
}

func formatJsonHandler(c echo.Context) error {
	req := new(FormatRequest)
	if err := c.Bind(req); err != nil || req.Content == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "リクエストボディに 'content' 文字列が含まれていません。",
		})
	}

	// リクエストパラメータの取得
	content := string(req.Content)

	// <br><hr>タグの書き換え
	useBr := strings.Contains(content, "<br>")
	useHr := strings.Contains(content, "<hr>")
	if useBr {
		content = strings.ReplaceAll(content, "<br>", "<br/>")
	}
	if useHr {
		content = strings.ReplaceAll(content, "<hr>", "<hr/>")
	}

	// XMLフォーマット処理
	x := xmlfmt.FormatXML(content, "", "\t")

	// <script>タグ内の整形処理
	scriptIndent := getScriptIndentTabs(x)
	if len(scriptIndent) > 0 {
		x = formatScriptTags(x, scriptIndent)
	}

	// <br><hr>タグの戻し処理
	if useBr {
		x = strings.ReplaceAll(x, "<br/>", "<br>")
	}
	if useHr {
		x = strings.ReplaceAll(x, "<hr/>", "<hr>")
	}

	// 先頭の改行を削除
	if strings.HasPrefix(x, "\n") {
		x = strings.TrimPrefix(x, "\n")
	}
	// 最終行改行
	x = x + "\n"

	// JSONで応答を返す
	return c.JSON(http.StatusOK, FormatResponse{
		FormattedContent: x,
	})
}

// フォーマット処理を実行するハンドラー関数
func formatFileHandler(c echo.Context) error {
	// フォームリクエストからファイルを取得（キー名: "file"）
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "ファイルが指定されていないか、リクエストが無効です。 'file' フィールドでアップロードしてください。",
		})
	}

	src, err := fileHeader.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "ファイルの読み込みに失敗しました。",
		})
	}
	defer src.Close()

	// ファイル内容の読み込み
	contentByte, err := io.ReadAll(src)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "ファイルデータの読み取りに失敗しました。",
		})
	}

	// 変換用文字列の取得
	content := string(contentByte)

	// <br><hr>タグの書き換え
	useBr := strings.Contains(content, "<br>")
	useHr := strings.Contains(content, "<hr>")
	if useBr {
		content = strings.ReplaceAll(content, "<br>", "<br/>")
	}
	if useHr {
		content = strings.ReplaceAll(content, "<hr>", "<hr/>")
	}

	// XMLフォーマット処理
	x := xmlfmt.FormatXML(content, "", "\t")

	// <script>タグ内の整形処理
	scriptIndent := getScriptIndentTabs(x)
	if len(scriptIndent) > 0 {
		x = formatScriptTags(x, scriptIndent)
	}

	// <br><hr>タグの戻し処理
	if useBr {
		x = strings.ReplaceAll(x, "<br/>", "<br>")
	}
	if useHr {
		x = strings.ReplaceAll(x, "<hr/>", "<hr>")
	}

	// 先頭の改行を削除
	if strings.HasPrefix(x, "\n") {
		x = strings.TrimPrefix(x, "\n")
	}
	// 最終行改行
	x = x + "\n"

	// 出力ファイル名の生成 (例: sample.html -> sampleF.html)
	origName := fileHeader.Filename
	ext := filepath.Ext(origName)
	base := origName[:len(origName)-len(ext)]
	downloadFileName := base + "F" + ext

	// ダウンロード応答ヘッダーを設定してレスポンスを返す
	c.Response().Header().Set(echo.HeaderContentDisposition, "attachment; filename=\""+downloadFileName+"\"")
	return c.Blob(http.StatusOK, "application/octet-stream", []byte(x))
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CORS 設定を追加
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{http.MethodPost, http.MethodOptions},
		AllowHeaders:  []string{echo.HeaderContentType},
		ExposeHeaders: []string{echo.HeaderContentDisposition},
	}))

	// POST
	e.POST("/format/json", formatJsonHandler)
	e.POST("/format/file", formatFileHandler)

	e.Logger.Fatal(e.Start(":6100"))
}
