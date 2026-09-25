package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

// ファイル読込
func readXML(fileName string) string {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	return string(b)
}

// インデントの数を数える
func countIndentTabs(line string) int {
	// 行頭にある '\t' 以外の文字（半角スペースなど）を取り除かないよう、'\t' のみを削除
	trimmed := strings.TrimLeft(line, "\t")

	// 削る前の長さ - 削った後の長さ = 先頭の '\t' の個数
	return len(line) - len(trimmed)
}

// インデントの数を数える
func getScriptIndentTabs(htmlData string) []int {
	var tabCounts []int

	// 1つ以上のタブ (\t+) で始まり、その直後に <script が続く行にマッチ
	scriptRegex := regexp.MustCompile(`^\t+<script>`)

	scanner := bufio.NewScanner(strings.NewReader(htmlData))
	for scanner.Scan() {
		line := scanner.Text()

		if scriptRegex.MatchString(line) {
			tabs := countIndentTabs(line)
			tabCounts = append(tabCounts, tabs)
		}
	}

	return tabCounts
}

// <script>?</script> 内のJavaScriptに改行とシンプルなインデントを適用する関数
func formatScriptTags(htmlContent string, indent []int) string {
	// <script>タグの開き・閉じを含む範囲を置換対象にする正規表現
	re := regexp.MustCompile(`(?s)<script>(.*?)</script>`)
	idx := 0

	return re.ReplaceAllStringFunc(htmlContent, func(match string) string {
		submatches := re.FindStringSubmatch(match)
		if len(submatches) < 2 {
			return match
		}
		jsCode := submatches[1]

		// 既に複数行に分かれていない1行状のコードをセミコロンや括弧で改行分割
		jsCode = strings.TrimSpace(jsCode)
		if jsCode == "" {
			return match
		}

		// セミコロン、波括弧の前後に改行を差し込む
		jsCode = strings.ReplaceAll(jsCode, ";", ";\n")
		jsCode = strings.ReplaceAll(jsCode, "{", "{\n")

		// 各行をスキャンしてインデントを整える
		scanner := bufio.NewScanner(strings.NewReader(jsCode))
		var formattedLines []string

		// <script>タグのインデント
		indentLevel := indent[idx] + 1
		endIndent := strings.Repeat("\t", indentLevel-1)
		idx++

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}

			// 閉じ括弧の場合はインデントを1段下げる
			if strings.HasPrefix(line, "}") {
				if indentLevel > 1 {
					indentLevel--
				}
			}

			// インデント（タブ）を付与
			indent := strings.Repeat("\t", indentLevel)
			formattedLines = append(formattedLines, indent+line)

			// 開き括弧で終わる場合は次の行のインデントを増やす
			if strings.HasSuffix(line, "{") {
				indentLevel++
			}
		}

		// 組み立て
		result := "<script>\n" + strings.Join(formattedLines, "\n") + "\n" + endIndent + "</script>"
		return result
	})
}

// XML最下位階層の改行を削除
func formatXML(outputFile string) (ss []byte) {
	fp, err := os.Open(outputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	s0 := ""
	for scanner.Scan() {
		s := scanner.Text()
		if len(strings.TrimSpace(s)) == 0 {
			continue
		}
		if strings.HasSuffix(s0, ">") {
			s0 = s0 + "\n"
			ss = append(ss, s0...)
			s0 = s
		} else {
			s0 = s0 + strings.TrimSpace(s)
		}
	}
	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}
	s0 = s0 + "\n"
	ss = append(ss, s0...)
	return ss
}
