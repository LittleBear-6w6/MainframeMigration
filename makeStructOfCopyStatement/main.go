package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Item struct {
	Level int
	Name  string
	Pic   string
}

func main() {
	//---設定項目---
	inputDir := "./test"
	outputDir := "./output"
	outputFile := "struct_base.def"
	//---設定項目ここまで---

	// 出力ディレクトリが存在しない場合は作成
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("出力ディレクトリ作成失敗: %v\n", err)
		return
	}

	outfile, err := os.Create(filepath.Join(outputDir, outputFile))
	if err != nil {
		fmt.Printf("出力ファイル作成失敗: %v\n", err)
		return
	}
	defer outfile.Close()
	writer := bufio.NewWriter(outfile)

	reLevel := regexp.MustCompile(`level\s*=\s*(\d+)`)
	reName := regexp.MustCompile(`name\s*=\s*"([^"]*)"`)
	rePic := regexp.MustCompile(`pic\s*=\s*"([^"]*)"`)
	reOccurs := regexp.MustCompile(`occurs\s*=\s*(\d+)`)

	files, err := os.ReadDir(inputDir)
	if err != nil {
		fmt.Printf("入力ディレクトリの読み込み失敗: %v\n", err)
		return
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".java" {
			targetPath := filepath.Join(inputDir, file.Name())
			items := processJavaFile(targetPath, reLevel, reName, rePic, reOccurs)

			// 拡張子を除いたファイル名を取得
			baseName := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))

			// 指定の構造体レイアウトで書き出し
			writer.WriteString(fmt.Sprintf("struct %s {\n", baseName))

			for _, it := range items {
				byteSize := calculateByteSize(it.Pic)
				if byteSize > 0 {
					// 縦を揃えず、BYTEの後と変数名の後に1つずつスペースを入れる形式
					writer.WriteString(fmt.Sprintf("\tBYTE\t%s[%d];\n", it.Name, byteSize))
				}
			}
			writer.WriteString("};\n\n")
			writer.Flush()
		}
	}
	fmt.Printf("処理完了: %s 内のファイルを解析し、%s に出力しました\n", inputDir, filepath.Join(outputDir, outputFile))
}

// PIC句からバイト数を算出するロジック
func calculateByteSize(pic string) int {
	if pic == "" {
		return 0
	}

	reNum := regexp.MustCompile(`\((\d+)\)`)
	match := reNum.FindStringSubmatch(pic)
	if len(match) < 2 {
		return 0
	}
	n, _ := strconv.Atoi(match[1])

	isPacked := strings.Contains(pic, "PACKED-DECIMAL")

	switch {
	case strings.HasPrefix(pic, "X"):
		return n
	case strings.HasPrefix(pic, "N"):
		return n * 2
	case strings.HasPrefix(pic, "9") || strings.HasPrefix(pic, "S9"):
		if isPacked {
			// (n + 1) / 2 の整数除算（切り上げ相当）
			return (n + 1) / 2
		}
		return n
	default:
		return 0
	}
}

// ファイル読み込みとOCCURS展開ロジック
func processJavaFile(filename string, reL, reN, reP, reO *regexp.Regexp) []Item {
	var items []Item
	var occursBuffer []Item
	occursCount, occursParentLevel, inOccurs := 0, 0, false

	file, _ := os.Open(filename)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		mLevel := reL.FindStringSubmatch(line)
		mName := reN.FindStringSubmatch(line)
		if len(mLevel) < 2 || len(mName) < 2 {
			continue
		}

		lvl, _ := strconv.Atoi(mLevel[1])
		name := mName[1]
		mPic := reP.FindStringSubmatch(line)
		mOccurs := reO.FindStringSubmatch(line)

		if inOccurs && lvl <= occursParentLevel {
			expandOccurs(&items, occursBuffer, occursCount)
			inOccurs = false
			occursBuffer = nil
		}

		if len(mOccurs) == 2 {
			inOccurs = true
			occursCount, _ = strconv.Atoi(mOccurs[1])
			occursParentLevel = lvl
			continue
		}

		pic := ""
		if len(mPic) == 2 {
			pic = mPic[1]
		}
		newItem := Item{Level: lvl, Name: name, Pic: pic}

		if inOccurs {
			occursBuffer = append(occursBuffer, newItem)
		} else if pic != "" {
			items = append(items, newItem)
		}
	}
	if inOccurs {
		expandOccurs(&items, occursBuffer, occursCount)
	}
	return items
}

func expandOccurs(target *[]Item, buffer []Item, count int) {
	for i := 1; i <= count; i++ {
		for _, b := range buffer {
			*target = append(*target, Item{
				Level: b.Level,
				Name:  fmt.Sprintf("%s_%d", b.Name, i),
				Pic:   b.Pic,
			})
		}
	}
}
