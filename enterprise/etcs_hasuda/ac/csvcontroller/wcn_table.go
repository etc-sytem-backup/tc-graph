package csvcontroller

import (
	"encoding/csv"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

const (
	// wcnテーブル
	tcWcnTablePath = "./tc_wcn_table/"
)

var Ch_resultWcn chan []string

// 初期化処理
func init() {
	// チャネル作成
	Ch_resultWcn = make(chan []string, 10)
}

// Read wcnテーブルファイルを読み込む
func Wcn_table() {

	t := time.NewTicker(10 * time.Second)

	// 常にディレクトリを監視
	for {
		select {
		case <-t.C:

			// ディレクトリ配下のcsvファイルを取得する
			files, err := ioutil.ReadDir(tcWcnTablePath)
			if err != nil {
				log.Fatal("ディレクトリ取得エラー:", err)
			}
			// 存在するファイル文ループする
			for _, f := range files {

				// csvファイルを開く
				file, err := os.Open(tcWcnTablePath + f.Name())
				if err != nil {
					// ファイルを閉じる
					file.Close()
				}

				// 1行ずつ読み取る
				reader := csv.NewReader(file)
				var wcnTableRecord []string
				for {
					wcnTableRecord, err = reader.Read()
					if err != nil {
						break
					}
					// スライスをstringに変換
					wcn := strings.Join(wcnTableRecord, ",")

					// alert.csvを作成する
					// 交通をチェック
					resultWcn := CheckTraffic(wcn)

					// チャネルにセット
					Ch_resultWcn <- resultWcn

				}
				// ファイルを閉じる
				defer file.Close()
			}
		}
	}
	t.Stop()
}
