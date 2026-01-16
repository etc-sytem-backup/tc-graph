package csvcontroller

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"localhost.com/iniread"
)

var (
	// 2点間交通量
	Count int

	// 精査済みwcn
	resultWcn []string
)

/*
   初期化
*/
func init() {
	iniread.Run() // config.ini読込
}

func CheckTraffic(wcn string) []string {
	// alert.csvを作成
	_, err := exec.Command("sh", "./script/make_alert.sh", wcn).Output()
	if err != nil {
		log.Fatal(err)
	}

	// alert.csvを開く
	file, err := os.Open("./ac_alert/alert.csv")
	if err != nil {
		// ファイルを閉じる
		file.Close()
	}
	// ファイルを閉じる
	defer file.Close()

	// 1行ずつ読み取り出力する
	scanner := bufio.NewScanner(file)
	var alertRecord []string
	for scanner.Scan() {
		// スライスに行をstringで詰める
		alertRecord = append(alertRecord, scanner.Text())
	}
	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// alert.csvが2行無い場合はnilを返す
	if len(alertRecord) != 2 {
		return nil
	}

	// 各行をスライスに変換
	line1 := strings.Split(alertRecord[0], ",")
	line2 := strings.Split(alertRecord[1], ",")

	// ステータスを取得
	status1 := line1[len(line1)-1]
	status2 := line2[len(line2)-1]
	time1 := line1[2]
	time2 := line2[2]
	// 時刻の時以下を取得しintに変換
	s1 := time1[8:14]
	s2 := time2[8:14]
	minutes1, _ := strconv.Atoi(s1)
	minutes2, _ := strconv.Atoi(s2)

	// ステータスを精査
	resultWcn = validate(status1, status2, minutes1, minutes2, wcn)

	// 交通量
	fmt.Println("交通量：", Count)

	// 予約チェック
	resultWcn = reserve(wcn, resultWcn)

	return resultWcn

}

// 交通精査
func validate(status1 string, status2 string, minutes1 int, minutes2 int, wcn string) []string {
	var resultWcn []string
	if status1 == "IN" && status2 == "PARK" {
		// 正常
		fmt.Println("正常通行")
		// 交通量カウント
		Count++

		// 長時間駐車チェック
		if minutes2-minutes1 > 5000 {
			// 5分超過の場合
			resultWcn = append(resultWcn, wcn, "LONG")
			fmt.Println("長時間駐車")
		}

		// 速度計算
		// day1 := time.Date(2000, 12, 31, 0, 0, 0, 0, time.Local)
		// day2 := time.Date(2001, 1, 2, 12, 30, 0, 0, time.Local)
		// duration := day2.Sub(day1)
	}
	if status1 == "IN" && status2 == "IN" {
		// 滞在時間を計算する
		result := minutes1 - minutes2
		// 1分超過なら異常
		if result > 100 {
			resultWcn = append(resultWcn, wcn, "REVERSE")
			fmt.Println("逆走")
		} else {
			// 1分以内なら渋滞
			fmt.Println("渋滞")
		}
	}
	if status1 == "IN" && status2 == "OUT" {
		// 異常
		fmt.Println("異常通行")
		// 逆走wcnに追加
		resultWcn = append(resultWcn, wcn, "REVERSE")
	}
	if status1 == "OUT" && status2 == "IN" {
		// 滞在時間を計算する
		result := minutes1 - minutes2
		// 30分超過なら正常
		if result > 3000 {
			fmt.Println("出戻り(正常)")
			// 交通量カウント
			Count++
		} else {
			// 30分以内なら異常
			resultWcn = append(resultWcn, wcn, "REVERSE")
			fmt.Println("異常通行")
		}
	}
	if status1 == "OUT" && status2 == "PARK" {
		// 異常
		fmt.Println("逆走")
		// 逆走wcnに追加
		resultWcn = append(resultWcn, wcn, "REVERSE")
	}
	if status1 == "OUT" && status2 == "OUT" {
		// 滞在時間を計算する
		result := minutes1 - minutes2
		// 1分超過なら異常
		if result > 100 {
			resultWcn = append(resultWcn, wcn, "REVERSE")
			fmt.Println("逆走")
		} else {
			// それ以外なら渋滞
			fmt.Println("渋滞")
		}
	}
	if status1 == "PARK" && status2 == "IN" {
		// 異常
		fmt.Println("逆走")
		// 逆走wcnに追加
		resultWcn = append(resultWcn, wcn, "REVERSE")
	}
	if status1 == "PARK" && status2 == "PARK" {
		// 滞在時間を計算する
		result := minutes1 - minutes2
		// 1分超過なら異常
		if result > 100 {
			resultWcn = append(resultWcn, wcn, "REVERSE")
			fmt.Println("逆走")
		} else {
			// それ以外なら渋滞
			fmt.Println("渋滞")
		}
	}
	if status1 == "PARK" && status2 == "OUT" {
		// 正常
		fmt.Println("正常通行")
	}
	return resultWcn
}

// 予約確認
func reserve(wcn string, resultWcn []string) []string {
	var reserveList []string
	// config.iniから予約wcnを取得
	r1 := iniread.Config.Car_reserve_0001
	r2 := iniread.Config.Car_reserve_0002
	r3 := iniread.Config.Car_reserve_0003
	r4 := iniread.Config.Car_reserve_0004
	r5 := iniread.Config.Car_reserve_0005
	reserveList = append(reserveList, r1, r2, r3, r4, r5)
	rs := strings.Join(reserveList, ",")

	// wcnが予約wcnと一致の場合、予約車両。
	if strings.Contains(rs, wcn) {
		resultWcn = append(resultWcn, wcn, "RESERVE")
		fmt.Println("予約車両です")
	}
	return resultWcn
}
