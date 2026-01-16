package csvcontroller

import (
	"encoding/csv"
	"log"
	"os"
)

func MakeReverseCsv(reverseRirekiWcn []string) {
	// スライスを配列に変換
	records := [][]string{reverseRirekiWcn}

	// csvを作成
	file, err := os.Create("./ac_reverse/reverse.csv")
	if err != nil {
		log.Fatal(err)
	}

	// 書き込み処理
	w := csv.NewWriter(file)
	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatal(err)
		}
	}

	// バッファに残っているデータをすべて書き込む
	w.Flush()

	// すべてのエラーハンドリング
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}
