package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-xmlfmt/xmlfmt"
)

func main() {
	// コマンドライン引数
	args := os.Args[1:]

	var inputFile string
	var outputFile string
	switch len(args) {
	case 0:
		// 引数が指定されていない場合はエラーを出力して終了
		fmt.Fprintf(os.Stderr, "エラー: 入力ファイル名が指定されていません。\n使用方法: %s <入力ファイル> [出力ファイル]\n", os.Args[0])
		os.Exit(1)
	case 1:
		// 入力ファイルのみ指定された場合
		inputFile = args[0]

		// 拡張子の前に "F" を挿入
		ext := filepath.Ext(inputFile)              // 例: ".html"
		base := inputFile[:len(inputFile)-len(ext)] // 例: "index"
		outputFile = base + "F" + ext               // 例: "indexF.html"
	default:
		// 入力ファイルと出力ファイル（またはそれ以上）が指定された場合
		inputFile = args[0]
		outputFile = args[1]
	}

	// ファイル読込
	xml1 := readXML(inputFile)

	// <br><hr>タグの処理
	useBr := strings.Contains(xml1, "<br>")
	useHr := strings.Contains(xml1, "<hr>")
	if useBr {
		xml1 = strings.ReplaceAll(xml1, "<br>", "<br/>")
	}
	if useHr {
		xml1 = strings.ReplaceAll(xml1, "<hr>", "<hr/>")
	}

	// XMLフォーマット
	x := xmlfmt.FormatXML(xml1, "", "\t")

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

	// 最終行改行
	x = x + "\n"

	// ファイル出力
	err := os.WriteFile(outputFile, []byte(x), 0666)
	if err != nil {
		log.Fatal(err)
	}
}
