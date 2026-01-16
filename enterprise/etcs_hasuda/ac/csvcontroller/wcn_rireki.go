package csvcontroller

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	// wcn履歴
	acCsvPath = "./ac_csv/"
)

// Read wcn履歴ファイルを読み込む
func Wcn_rireki(resultWcn []string) {

	// ディレクトリ配下のcsvファイルを取得する
	files, err := ioutil.ReadDir(acCsvPath)
	if err != nil {
		log.Fatal("ディレクトリ取得エラー:", err)
	}
	// 存在するファイル文ループする
	for _, f := range files {

		// csvファイルを開く
		file, err := os.Open(acCsvPath + f.Name())
		if err != nil {
			file.Close()
		}

		// 1行ずつ読み取る
		reader := csv.NewReader(file)
		var rirekiWcn []string
		for {
			rirekiWcn, err = reader.Read()
			if err != nil {
				break
			}

			// 比較のためstringに変換
			rws := strings.Join(resultWcn, "")

			// 一致するwcnで逆走チェックを追加
			for _, wcn := range rirekiWcn {
				// wcnが一致。かつ、長時間駐車
				if strings.Contains(rws, wcn) && strings.Contains(rws, "LONG") {
					rirekiWcn = append(rirekiWcn, "", "", "○")
				}
				// wcnが一致。かつ、予約
				if strings.Contains(rws, wcn) && strings.Contains(rws, "RESERVE") {
					rirekiWcn = append(rirekiWcn, "", "○")
				}
				// wcnが一致。かつ、逆走
				if strings.Contains(rws, wcn) && strings.Contains(rws, "REVERSE") {
					rirekiWcn = append(rirekiWcn, "○")
				}
			}
			fmt.Println(rirekiWcn)

			// 逆走csv作成
			MakeReverseCsv(rirekiWcn)
		}
		// ファイルを閉じる
		defer file.Close()
	}
}
