package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const INDENT_STR = "　"

var ICONS = [...]string{"■", "□", "▼", "▽", "●", "○"}

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("マークダウンファイルを指定してください。")
		return
	}
	mdfile := os.Args[1]

	// マークダウンファイルを読み込む
	fp, err := os.Open(mdfile)
	if err != nil {
		log.Fatal(err)
	}

	sls := scanLines(fp)                 // ファイルから行を読み取る
	menu := makeMenu(sls)                // マークダウンの表題から目次文字列を生成
	fls := formatLines(makeMdLines(sls)) // マークダウン文字列を整形
	ls := makeLines(sls, menu, fls)      // 最終的なテキストのスライスを返却

	for _, v := range ls {
		fmt.Println(v)
	}
}

// マークダウンとして処理するのに必要な業のみ取り出したスライスを返却
func makeMdLines(ls []string) (mdls []string) {
	readFlg := false
	for _, l := range ls {
		if l == "{code}" {
			readFlg = true
			continue
		}
		if l == "{/code}" {
			readFlg = false
			break
		}
		if readFlg {
			mdls = append(mdls, l)
		}
	}
	if len(mdls) == 0 {
		log.Fatal()
	}
	return
}

// 目次の文字列を生成する
func makeMenu(ls []string) (menu string) {
	mls := make([]string, 0)
	for _, l := range ls {
		cnt := strings.Count(l, "#")
		if 0 < cnt {
			if l == "#contents" {
				continue
			}
			indent := strings.Repeat(INDENT_STR, cnt-1)
			nl := indent + "- " + strings.Trim(l, "# ")
			mls = append(mls, nl)
		}
	}
	menu = strings.Join(mls, "\n")
	return
}

// ファイル内のすべての行テキストを読み込む
func scanLines(fp *os.File) (ls []string) {
	ls = make([]string, 0)
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		l := scanner.Text()
		ls = append(ls, l)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal()
	}
	return
}

// マークダウン行のスライスを整形する
func formatLines(ls []string) (fls []string) {
	fls = make([]string, 0)
	ic := 0
	for _, l := range ls {
		cnt := strings.Count(l, "#")
		var nl string
		switch cnt {
		case 0:
			indent := strings.Repeat(INDENT_STR, ic)
			nl = indent + l
		case 1:
			indent := strings.Repeat(INDENT_STR, cnt-1)
			nl = indent + ICONS[0] + strings.Trim(l, "# ")
			ic = cnt
		case 2:
			indent := strings.Repeat(INDENT_STR, cnt-1)
			nl = indent + ICONS[1] + strings.Trim(l, "# ")
			ic = cnt
		case 3:
			indent := strings.Repeat(INDENT_STR, cnt-1)
			nl = indent + ICONS[2] + strings.Trim(l, "# ")
			ic = cnt
		}
		fls = append(fls, nl)
	}
	return
}

// 最終的なテキストのリストを生成
func makeLines(ls []string, menu string, fls []string) (ts []string) {
	ts = make([]string, 0)
	pasteFlg := false
	for _, v := range ls {
		if v == "#contents" {
			ts = append(ts, ICONS[0]+"目次")
			ts = append(ts, menu)
			continue
		}
		if v == "{code}" {
			pasteFlg = true
			s := strings.Join(fls, "\n")
			ts = append(ts, s)
		}
		if v == "{/code}" {
			pasteFlg = false
			continue
		}
		if !pasteFlg {
			ts = append(ts, v)
		}
	}
	return
}
