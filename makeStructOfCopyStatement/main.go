package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

// Item 抽出データを保持する構造体
type Item struct {
	Level int
	Name  string
	Pic   string
}

func main() {
	// 出力用のファイルを1つ作成
	outFile, err := os.Create("analysis_result.txt")
	if err != nil {
		fmt.Printf("出力ファイル作成失敗: %v\n", err)
		return
	}
	defer outFile.Close()
	writer := bufio.NewWriter(outFile)

	// 正規表現の定義
	reLevel := regexp.MustCompile(`level\s*=\s*(\d+)`)
	reName := regexp.MustCompile(`name\s*=\s*"([^"]*)"`)
	rePic := regexp.MustCompile(`pic\s*=\s*"([^"]*)"`)
	reOccurs := regexp.MustCompile(`occurs\s*=\s*(\d+)`)

	// カレントディレクトリのJavaファイルを走査
	files, _ := os.ReadDir(".")

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".java" {
			// 1. ファイルごとにデータを抽出
			items := processJavaFile(file.Name(), reLevel, reName, rePic, reOccurs)

			// 2. 出力ファイルにセクション区切りを書き込む
			fmt.Fprintf(writer, "\n================================================================================\n")
			writer.WriteString(fmt.Sprintf(" SOURCE FILE: %s\n", file.Name()))
			fmt.Fprintf(writer, "================================================================================\n")
			writer.WriteString(fmt.Sprintf("%-3s | %-40s | %s\n", "LV", "NAME (EXTENDED)", "PIC"))
			fmt.Fprintf(writer, "--------------------------------------------------------------------------------\n")

			if len(items) == 0 {
				writer.WriteString(" (有効な項目は見つかりませんでした)\n")
			} else {
				for _, it := range items {
					writer.WriteString(fmt.Sprintf("%02d  | %-40s | %s\n", it.Level, it.Name, it.Pic))
				}
			}
			writer.Flush() // バッファを書き出し
		}
	}

	fmt.Println("抽出完了: analysis_result.txt を確認してください。")
}

// OCCURS展開対応
func processJavaFile(filename string, reL, reN, reP, reO *regexp.Regexp) []Item {
	var items []Item
	var occursBuffer []Item
	occursCount := 0
	occursParentLevel := 0
	inOccurs := false

	file, err := os.Open(filename)
	if err != nil {
		return items
	}
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

		// OCCURS終了判定
		if inOccurs && lvl <= occursParentLevel {
			expandOccurs(&items, occursBuffer, occursCount)
			inOccurs = false
			occursBuffer = nil
		}

		// 新規OCCURS開始判定
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
