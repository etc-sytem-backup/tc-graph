// make_data_hasuda02.go
// 蓮田SA向けバージョン　平均画面用データ作成
package main

import (
    //	"context"
	"fmt"
	"log"
    //	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
	"sync"
    "encoding/csv"
    "sort"
    //    "crypto/tls"
    //    "net/smtp"

	"localhost.com/iniread"
	"localhost.com/readcsv"
)


// プロジェクト定数
const (

    // Noting
)

// package変数
var log_run_path string = "0"   // 動作ログファイル格納用パス(初期値 : Off)
var duration_list_show []string // 長時間車両一覧スライス（表示用）
var traffic_jam_flag string     // ランプ停滞中フラグ On:1 / Off:0




/*
   受信バイナリ保存用。CSVファイル保存用。動作ログ保存用。
   log_bin_path : 送受信ログPATH（バイナリ）
   log_csv_path : 送受信ログPATH
   log_run_path : 動作ログPATH
*/
func make_log_folder(log_bin_path string, log_csv_path string, log_run_path string) error {

	//logフォルダ直下に、bin, csv, runのフォルダがあるか確認し、なければ作成する
	_, err := os.Open("./log")
	if os.IsNotExist(err) {

		// logフォルダ作成
		err = os.Mkdir("./log", 0777)
		if err != nil {
			return err
		}
	}

	_, err = os.Open(log_run_path)
	if os.IsNotExist(err) {

		// runフォルダ作成
		err = os.Mkdir(log_run_path, 0777)
		if err != nil {
			return err
		}
	}

	_, err = os.Open(log_csv_path)
	if os.IsNotExist(err) {

		// csvフォルダ作成
		err = os.Mkdir(log_csv_path, 0777)
		if err != nil {
			return err
		}
	}

	_, err = os.Open(log_bin_path)
	if os.IsNotExist(err) {

		// csvフォルダ作成
		err = os.Mkdir(log_bin_path, 0777)
		if err != nil {
			return err
		}
	}

	return nil
}

/* ログ設定 */
func log_setup() {

    // log保存ディレクトリの準備
	err := make_log_folder("./log/bin", "./log/csv", "./log/run")
	if err != nil {
		panic(err)
	}

	log_run_path = "./log/run/" // ログファイル保存用ディレクトリ

	// Log保存ファイル設定
	now := time.Now()
	year_val, month_val, day_val := now.Date() // 年月日を数字で取得してみる
	log_filename := fmt.Sprintf(iniread.Config.Run_log_path+"/"+"%04d%02d%02d.log", year_val, int(month_val), day_val)

    fmt.Printf("Run_log_path : %s\n",iniread.Config.Run_log_path)
    fmt.Printf("log_filename : %s\n",log_filename)

	// ファイルが既に存在している場合はスルー。
	// ファイルが存在していない場合は、新規作成してログデータの保存先にする
    // 結果、1日1ファイルのログが残る
	log_file, err := os.OpenFile(log_filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {

		// ファイルが無い場合は新規作成
		log_file, _ = os.OpenFile(log_filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	}

	// Logの保存先をファイルにする(デフォルトは標準出力)
	log.SetOutput(log_file)

}

/*
   初期化
*/
func init() {
	iniread.Run() // config.ini読込

    // ログファイル保存設定
    log_setup()
    go timer_10()   // logファイルのローテーション処理起動

    // 起動マシンの情報をログ出力
    log.Printf("NumCPU: %d\n", runtime.NumCPU())
	log.Printf("NumGoroutine: %d\n", runtime.NumGoroutine())
	log.Printf("Version: %s\n", runtime.Version())

    // // 予約テーブル初期化(ファイルを消す→予約情報(config.ini)で再構築)
    // file_remove(iniread.Config.Reserve_table_path,"*.csv")

    // ワークディレクトリ群の作成
    make_work_folder()
    
}

// ワークディレクトリ群の作成。
// ワークディレクトリ群が無い場合は新規に作成する。
func make_work_folder() error {

    // 満空管理テーブルと駐車パス管理テーブルの格納場所
	_, err := os.Open("./parking_list")
	if os.IsNotExist(err) {
		err := os.Mkdir("./parking_list", 0777)
		if err != nil {
			return err
		}
	}


    // 直近の通過履歴ファイル格納用
	_, err = os.Open("./driving_history")
	if os.IsNotExist(err) {
		err := os.Mkdir("./driving_history", 0777)
		if err != nil {
			return err
		}
	}


    // ac管理下のWCN_rireki.csvバックアップ格納場所
	_, err = os.Open("./old_wcn_rireki")
	if os.IsNotExist(err) {
		err := os.Mkdir("./old_wcn_rireki", 0777)
		if err != nil {
			return err
		}
	}

    // display_server用のファイル格納場所
	_, err = os.Open("./disp_data")
	if os.IsNotExist(err) {
		err := os.Mkdir("./disp_data", 0777)
		if err != nil {
			return err
		}
	}

	return nil
}

// パラメータで渡されたファイルを削除する
// Usage：file_remove("/home/k/opt/aps","*.csv")
func file_remove(target_file_path string, remove_file_name string) {

    log.Printf("cmd: rm -rf %s/%s\n",target_file_path,remove_file_name)

    _, err := exec.Command("bash", "-c", "rm -rf " + target_file_path + "/" + remove_file_name).Output()
    if err != nil {
        log.Fatal(err)
    }
}


/*
   10秒毎に任意の処理を行う
   ・ログ作成
*/
func timer_10() {
	t := time.NewTicker(10 * time.Second) // 10秒おきに通知
	defer t.Stop() // タイマを止める。 <- これがないとメモリリークを起こす

	for {
		select {
		case <-t.C: // 10秒経過した。

            /* 動作ログ作成 */
			now := time.Now()
			year_val, month_val, day_val := now.Date() // 年月日を数字で取得してみる
			log_filename := fmt.Sprintf(log_run_path + "/" + "%04d%02d%02d.log", year_val, int(month_val), day_val)

			// ファイルが既に存在している場合はスルー。
			// ファイルが存在していない(1日経過している)場合は、新規作成してログデータの保存先にする
			log_file, err := os.OpenFile(log_filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {

				// ファイルが無い場合は新規作成
				log_file, _ = os.OpenFile(log_filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			}

			log.SetOutput(log_file)
		}
	}
}


/* ミリ秒含みの日付文字列を作成する  */
func get_datestr() string{

    // 現在時間を取得
	now := time.Now()
	nowUTC := now.UTC()

    // ミリ秒の算出（文字列変換含む）
    t2 := nowUTC.UnixNano() / int64(time.Millisecond)    // 時間(ナノ秒)を時間(ミリ秒)に変換
    t2_str := strconv.Itoa(int(t2))                      // 時間(ミリ秒)を文字列に変換
    ms_str := t2_str[len(t2_str)-3:]                     // 時間(ミリ秒)文字列から、ミリ秒部分だけを切り出す


    // ミリ秒を含む日時データを作成
    year_val, month_val, day_val := now.Date() // 年月日を数字で取得してみる
    time_stamp := fmt.Sprintf("%04d%02d%02d%02d%02d%02d%s", year_val, int(month_val), day_val, now.Hour(), now.Minute(), now.Second(),ms_str)

    return time_stamp
}


// 与えられた年月日時分秒(ms含む)２つの差分（秒）を求める
//   before_time : 過去日付
//   after_time : 未来/現在日付
// Result
//   int : 秒数
func date_duration(before_time string, after_time string) (int, error) {
	t1_year_s := before_time[0:4]
	t1_month_s := before_time[4:6]
	t1_day_s := before_time[6:8]
	t1_hh_s := before_time[8:10]
	t1_mm_s := before_time[10:12]
	t1_ss_s := before_time[12:14]

	t2_year_s := after_time[0:4]
	t2_month_s := after_time[4:6]
	t2_day_s := after_time[6:8]
	t2_hh_s := after_time[8:10]
	t2_mm_s := after_time[10:12]
	t2_ss_s := after_time[12:14]

	t1_year, err := strconv.Atoi(t1_year_s)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	t1_month, err := strconv.Atoi(t1_month_s)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	t1_day, err := strconv.Atoi(t1_day_s)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	t1_hh, err := strconv.Atoi(t1_hh_s)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	t1_mm, err := strconv.Atoi(t1_mm_s)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	t1_ss, err := strconv.Atoi(t1_ss_s)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	t2_year, err := strconv.Atoi(t2_year_s)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	t2_month, err := strconv.Atoi(t2_month_s)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	t2_day, err := strconv.Atoi(t2_day_s)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	t2_hh, err := strconv.Atoi(t2_hh_s)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	t2_mm, err := strconv.Atoi(t2_mm_s)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	t2_ss, err := strconv.Atoi(t2_ss_s)
	if err != nil {
		log.Println(err)
		return 0, err
	}

    // 時間差を演算
	before := time.Date(t1_year, time.Month(t1_month), t1_day, t1_hh, t1_mm, t1_ss, 0, time.Local)
	after := time.Date(t2_year, time.Month(t2_month), t2_day, t2_hh, t2_mm, t2_ss, 0, time.Local)
	duration := after.Sub(before)

	// 必ず正の数とする。
	totalSeconds := int(duration.Seconds())
    if totalSeconds < 0 {
        totalSeconds = -totalSeconds  // 符号の反転 totalSeconds = totalSeconds * -1 と同意。
    }

	return totalSeconds, nil
}


/* 用途コードを平仮名に変換*/
func change_youto_code(code string) string {

    var result = ""

    switch code {

    // 自家用
    case "bb": // さ
        result = "さ"
    case "bd": // す
        result = "す"
    case "be": // せ
        result = "せ"
    case "bf": // そ
        result = "そ"
    case "c0": // た
        result = "た"
    case "c1": // ち
        result = "ち"
    case "c2": // つ
        result = "つ"
    case "c3": // て
        result = "て"
    case "c4": // と
        result = "と"
    case "c5": // な
        result = "な"
    case "c6": // に
        result = "に"
    case "c7": // ぬ
        result = "ぬ"
    case "c8": // ね
        result = "ね"
    case "c9": // の
        result = "の"
    case "ca": // は
        result = "は"
    case "cb": // ひ
        result = "ひ"
    case "cc": // ふ
        result = "ふ"
    case "ce": // ほ
        result = "ほ"
    case "cf": // ま
        result = "ま"
    case "d0": // み
        result = "み"
    case "d1": // む
        result = "む"
    case "d2": // め
        result = "め"
    case "d3": // も
        result = "も"
    case "d4": // や
        result = "や"
    case "d5": // ゆ
        result = "ゆ"
    case "d7": // ら
        result = "ら"
    case "d8": // り
        result = "り"
    case "d9": // る
        result = "る"
    case "db": // ろ
        result = "ろ"

    // 貸渡（レンタカー）
    case "da": // れ
        result = "れ"
    case "dc": // わ
        result = "わ"

    // 事業用
    case "b1": // あ
        result = "あ"
    case "b2": // い
        result = "い"
    case "b3": // う
        result = "う"
    case "b4": // え
        result = "え"
    case "b6": // か
        result = "か"
    case "b7": // き
        result = "き"
    case "b8": // く
        result = "く"
    case "b9": // け
        result = "け"
    case "ba": // こ
        result = "こ"
    case "a6": // を
        result = "を"

    // 駐留軍人軍属私有車両用等
    case "45": // E
        result = "Ｅ"
    case "48": // H
        result = "Ｈ"
    case "4b": // K
        result = "Ｋ"
    case "4d": // M
        result = "Ｍ"
    case "54": // T
        result = "Ｔ"
    case "59": // Y
        result = "Ｙ"
    case "d6": // よ
        result = "よ"
    default:
        result = code
    }

    return result
}

/*陸運局支局コードを地名に変換*/
func change_sikyoku_code(code string) string {

    var result = ""

    switch code {
    case "535053": // 札幌 SPS
        result = "札幌"
    case "535020": // 札   SP 
        result = "札"
    case "484448": // 函館 HDH
        result = "函館"
    case "484420": // 函   HD 
        result = "函"
    case "414b41": // 旭川 AKA
        result = "旭川"
    case "414b20": // 旭   AK 
        result = "旭"
    case "4d524d": // 室蘭 MRM
        result = "室蘭"
    case "4d5220": // 室   MR 
        result = "室"
    case "4b524b": // 釧路 KRK
        result = "釧路"
    case "4b5220": // 釧   KR 
        result = "釧"
    case "4f484f": // 帯広 OHO
        result = "帯広"
    case "4f4820": // 帯   OH 
        result = "帯"
    case "4b494b": // 北見 KIK
        result = "北見"
    case "4b4920": // 北   KI 
        result = "北"
    case "414d41": // 青森 AMA
        result = "青森"
    case "414d48": // 八戸 AMH
        result = "八戸"
    case "414d20": // 青   AM 
        result = "青"
    case "495449": // 岩手 ITI
        result = "岩手"
    case "495420": // 岩   IT 
        result = "岩"
    case "4D4753": // 仙台 MGS
        result = "仙台"
    case "4d474d": // 宮城 MGM
        result = "宮城"
    case "4d4720": // 宮   MG 
        result = "宮"
    case "415441": // 秋田 ATA
        result = "秋田"
    case "415420": // 秋   AT 
        result = "秋"
    case "594120": // 山形 YA 
        result = "山形"
    case "594153": // 庄内 YAS
        result = "庄内"
    case "465320": // 福島 FS 
        result = "福島"
    case "465341": // 会津 FSA
        result = "会津"
    case "465349": // いわきFSI
        result = "いわき"
    case "49474d": // 水戸 IGM
        result = "水戸"
    case "494754": // 土浦 IGT
        result = "土浦"
    case "49474b": // つくばIGK
        result = "つくば"
    case "494749": // 茨城 IGI
        result = "茨城"
    case "494720": // 茨   IG 
        result = "茨"
    case "544755": // 宇都宮TGU
        result = "宇都宮"
    case "54474e": // 那須 TGN
        result = "那須"
    case "544743": // とちぎTGC
        result = "とちぎ"
    case "544754": // 栃木 TGT
        result = "栃木"
    case "544720": // 栃   TG 
        result = "栃"
    case "474d47": // 群馬 GMG
        result = "群馬"
    case "474d54": // 高崎 GMT
        result = "高崎"
    case "474d20": // 群   GM 
        result = "群"
    case "53544f": // 大宮 STO
        result = "大宮"
    case "535447": // 川越 STG
        result = "川越"
    case "535454": // 所沢 STT
        result = "所沢"
    case "53544b": // 熊谷 STK
        result = "熊谷"
    case "535442": // 春日部STB
        result = "春日部"
    case "535453": // 埼玉 STS
        result = "埼玉"
    case "535420": // 埼   ST 
        result = "埼"
    case "434243": // 千葉 CBC
        result = "千葉"
    case "434254": // 成田 CBT
        result = "成田"
    case "43424e": // 習志野CBN
        result = "習志野"
    case "434253": // 袖ヶ浦CBS
        result = "袖ヶ浦"
    case "434244": // 野田 CBD
        result = "野田"
    case "43424b": // 柏   CBK
        result = "柏"
    case "434220": // 千   CB 
        result = "千"
    case "544b53": // 品川 TKS
        result = "品川"
    case "544f53": // 品   TOS
        result = "品"
    case "544b4e": // 練馬 TKN
        result = "練馬"
    case "544f4e": // 練   TON
        result = "練"
    case "544b41 ": // 足立 TKA
        result = "足立"
    case "544f41 ": // 足   TOA
        result = "足"
    case "544b48": // 八王子TKH
        result = "八王子"
    case "544b54": // 多摩 TKT
        result = "多摩"
    case "544f54": // 多   TOT
        result = "多"
    case "4b4e59": // 横浜 KNY
        result = "横浜"
    case "4b4e4b": // 川崎 KNK
        result = "川崎"
    case "4b4e4e": // 湘南 KNN
        result = "湘南"
    case "4b4e53": // 相模 KNS
        result = "相模"
    case "4b4e20": // 神   KN 
        result = "神"
    case "594e20": // 山梨 YN 
        result = "山梨"
    case "464a53": // 富士山FJS
        result = "富士山"
    case "4e474e": // 新潟 NGN
        result = "新潟"
    case "4e474f": // 長岡 NGO
        result = "長岡"
    case "4e4720": // 新   NG 
        result = "新"
    case "545954": // 富山 TYT
        result = "富山"
    case "545920": // 富   TY 
        result = "富"
    case "494b4b": // 金沢 IKK
        result = "金沢"
    case "494b49": // 石川 IKI
        result = "石川"
    case "494b20": // 石   IK 
        result = "石"
    case "4e4e4e": // 長野 NNN
        result = "長野"
    case "4e4e4d": // 松本 NNM
        result = "松本"
    case "4e4e53": // 諏訪 NNS
        result = "諏訪"
    case "4e4e20": // 長   NN 
        result = "長"
    case "464920": // 福井 FI 
        result = "福井"
    case "474647": // 岐阜 GFG
        result = "岐阜"
    case "474648": // 飛騨 GFH
        result = "飛騨"
    case "474620": // 岐   GF 
        result = "岐"
    case "535a53": // 静岡 SZS
        result = "静岡"
    case "535a48": // 浜松 SZH
        result = "浜松"
    case "535a4e": // 沼津 SZN
        result = "沼津"
    case "535a49": // 伊豆 SZI
        result = "伊豆"
    case "535a20": // 静   SZ 
        result = "静"
    case "41434e": // 名古屋ACN
        result = "名古屋"
    case "414354": // 豊橋 ACT
        result = "豊橋"
    case "41435a": // 岡崎 ACZ
        result = "岡崎"
    case "41434d": // 三河 ACM
        result = "三河"
    case "414359": // 豊田 ACY
        result = "豊田"
    case "414349": // 一宮 ACI
        result = "一宮"
    case "41434f": // 尾張小ACO牧
        result = "尾張小"
    case "414320": // 愛   AC 
        result = "愛"
    case "4d454d": // 三重 MEM
        result = "三重"
    case "4d4553": // 鈴鹿 MES
        result = "鈴鹿"
    case "4d4520": // 三   ME 
        result = "三"
    case "534953": // 滋賀 SIS
        result = "滋賀"
    case "534920": // 滋   SI 
        result = "滋"
    case "4b544b": // 京都 KTK
        result = "京都"
    case "4b5420": // 京   KT 
        result = "京"
    case "4f534f": // 大阪 OSO
        result = "大阪"
    case "4f534e": // なにわOSN
        result = "なにわ"
    case "4f5353": // 堺   OSS
        result = "堺"
    case "4f535a": // 和泉 OSZ
        result = "和泉"
    case "4f5320": // 大   OS 
        result = "大"
    case "4f5349": // 泉   OSI
        result = "泉"
    case "48474b": // 神戸 HGK
        result = "神戸"
    case "484748": // 姫路 HGH
        result = "姫路"
    case "484720": // 兵   HG 
        result = "兵"
    case "4e524e": // 奈良 NRN
        result = "奈良"
    case "4e5220": // 奈   NR 
        result = "奈"
    case "574b57": // 和歌山WKW
        result = "和歌山"
    case "574b20": // 和   WK 
        result = "和"
    case "545454": // 鳥取 TTT
        result = "鳥取"
    case "545420": // 鳥   TT 
        result = "鳥"
    case "534e20": // 島根 SN 
        result = "島根"
    case "534d20": // 島   SM 
        result = "島"
    case "4f594f": // 岡山 OYO
        result = "岡山"
    case "4f594b": // 倉敷 OYK
        result = "倉敷"
    case "4f5920": // 岡   OY 
        result = "岡"
    case "485348": // 広島 HSH
        result = "広島"
    case "485346": // 福山 HSF
        result = "福山"
    case "485320": // 広   HS 
        result = "広"
    case "595553": // 下関 YUS
        result = "下関"
    case "595559": // 山口 YUY
        result = "山口"
    case "595520": // 山   YU 
        result = "山"
    case "545354": // 徳島 TST
        result = "徳島"
    case "545320": // 徳   TS 
        result = "徳"
    case "4b414b": // 香川 KAK
        result = "香川"
    case "4b4120": // 香   KA 
        result = "香"
    case "454820": // 愛媛 EH 
        result = "愛媛"
    case "4b434b": // 高知 KCK
        result = "高知"
    case "4b4320": // 高   KC 
        result = "高"
    case "464f46": // 福岡 FOF
        result = "福岡"
    case "464f4b": // 北九州FOK
        result = "北九州"
    case "464f52": // 久留米FOR
        result = "久留米"
    case "464f43": // 筑豊 FOC
        result = "筑豊"
    case "464f20": // 福   FO 
        result = "福"
    case "534153": // 佐賀 SAS
        result = "佐賀"
    case "534120": // 佐   SA 
        result = "佐"
    case "4e5320": // 長崎 NS 
        result = "長崎"
    case "4e5353": // 佐世保NSS
        result = "佐世保"
    case "4b554b": // 熊本 KUK
        result = "熊本"
    case "4b5520": // 熊   KU 
        result = "熊"
    case "4f5420": // 大分 OT 
        result = "大分"
    case "4d5a20": // 宮崎 MZ 
        result = "宮崎"
    case "4b4f4b": // 鹿児島KOK
        result = "鹿児島"
    case "4b4f20": // 鹿   KO 
        result = "鹿"
    case "4f4e4f": // 沖縄 ONO
        result = "沖縄"
    case "4f4e20": // 沖   ON 
        result = "沖"
    default:
        result = code
    }

    return result
}


// 満空管理ファイルを作成する。
// その他、display表示に連携する管理ファイルを作成する。
// この関数は、Goルーチンとして、一定間隔毎に繰り返し処理される（繰り返し間隔はConfig.iniに設定）
func make_parking_table() {

    var in_park_list []string // IN→PARKの通過履歴と時速を保存

	t := time.NewTicker(time.Duration(iniread.Config.Request_interval) * time.Millisecond) // config.iniの設定毎
	defer t.Stop() // タイマを止める。 <- これがないとメモリリークを起こす
    
	for {
		select {
		case <-t.C:
			// タイマー時間経過した。（Config.iniに設定の時間間隔）

            // ../ac/tc_wcn_table/WCN_table.csv から、収集した全WCNをstring配列で取得する
            wcn := readcsv.Read("../ac/tc_wcn_table/WCN_table.csv")

            //            log.Printf("=== wcn check start ===\n")
            // 取得したWCNだけループ
			for _, w := range wcn {

                //log.Printf("w:%s",w)
                wcn := strings.Split(w,",") // WCN番号の取出し

				// 指定のWCNで、直近の２レコードを収集。
                // ../ac/ac_csv/WCN_rireki.csvを解析し、./driving_history/driving_history.csvを作成
				_, err := exec.Command("./script/make_driving_history.sh", wcn[0]).Output()
				if err != nil {
                    log.Printf("make_driving_history.sh Error!!\n")
					log.Fatal(err)
				}

                // ★この処理は無理にbashにさせなくても。。。。
				// driving_history.csvの行数を取得（直近２レコードあれば２行のはず）
				lines, err := exec.Command("./script/get_line_count.sh").Output()
				if err != nil {
                    log.Printf("get_line_count.sh Error!!\n")
					log.Fatal(err)
				}

				// intに変換
				s := string(lines)
				s = strings.TrimRight(s, "\n")
				linecnt, err := strconv.Atoi(s)
				if err != nil {
                    log.Printf("strconv.Atoi Error!!\n")
					log.Fatal(err)
				}

                // -----------------------------------------------------------------------------
				// 1行しかない場合、通過アンテナがPARKだった場合は、駐車場入庫と見なして処理する。
                // -----------------------------------------------------------------------------
				if linecnt == 1 {
                                        statuses_one, err := exec.Command("./script/get_driving_status.sh", "1P").Output() // １行目の５カラムのみを取得 (IN or PARK or OUT)
                    if err != nil {
                        log.Printf("get_driving_status.sh status_one Error!!\n")
                        log.Fatal(err)
                    }

                    // stringに変換して余分な改行を取る
                    status_one := string(statuses_one)
                    status_one = strings.TrimRight(status_one, "\n")

                    //debug
                    // fmt.Printf("debug status_one = %s\n",status_one)
                    // fmt.Printf("debug status_one = %s\n",status_one)

                    // status_oneが空(ファイル解析に失敗した)の場合は何もせずに次のループへ
                    if status_one == "" {
                        continue
                    }

                    // 通過履歴がPARK？（INを飛ばしてPARKだった？）
                    if status_one == "PARK" {
                        
                        // すでに満空管理テーブルに登録されている場合は次のwcn番号処理へ（continue）
                        results, err := exec.Command("./script/check_parking_table.sh", wcn[0],"./parking_list/parking_table.csv").Output()
                        if err != nil {
                            log.Printf("make_parking_table.sh Error!!\n")
                        }
                        result := string(results)
                        result = strings.TrimRight(result, "\n")

                        if result == "1" { // すでに登録されている
                            // log.Printf("debug: There is %s in parking_table.csv.",wcn[0])
                            // fmt.Printf("debug: There is %s in parking_table.csv.",wcn[0])
                            continue
                        }

                        // log.Printf("debug: There isn't %s in parking_table.csv.",wcn[0])
                        // fmt.Printf("debug: There isn't %s in parking_table.csv.",wcn[0])
                        
                        // 車両情報取得(./driving_history/driving_history.csv 1行目)を取得
                        parking_cardatas, err := exec.Command("bash", "-c", "awk -F',' '{print $0}' ./driving_history/driving_history.csv | sed -n 1P").Output()
                        if err != nil {
                            log.Printf(" 車両情報(./driving_history/driving_history.csv 1行目)を取得できませんでした。: IN無しでいきなりPARK。 Error!!\n")
                            log.Fatal(err)
                        }
                        // stringに変換
                        parking_cardata := string(parking_cardatas)
                        parking_cardata = strings.TrimRight(parking_cardata, "\n")


                        log.Printf("Directions. %s : %s --> %s\n",parking_cardata,"----",status_one)
                        fmt.Printf("Directions. %s : %s --> %s\n",parking_cardata,"----",status_one)

                        // parking_table.csvにレコードを追加する
                        // 同じ車両のデータがすでに登録されている場合、新しいデータで更新する。 ← この機能は必要ないが、今回はあっても害はないのでそのままにする。
                        _, err = exec.Command("./script/make_parking_table.sh", wcn[0],parking_cardata).Output()
                        if err != nil {
                            log.Printf("make_parking_table.sh Error!!\n")
                        }
                        
                    }
                    
                    // 次のWCN番号処理へ
                    continue
				}

				// driving_history.csvが2行以上なら満空管理テーブル作成処理を実施
				statuses1, err := exec.Command("./script/get_driving_status.sh", "1P").Output() // １行目の５カラムのみを取得 (IN or PARK or OUT)
				if err != nil {
                    log.Printf("get_driving_status.sh status1 Error!!\n")
					log.Fatal(err)
				}

				statuses2, err := exec.Command("./script/get_driving_status.sh", "2P").Output() // ２行目の５カラムのみを取得 (IN or PARK or OUT)
				if err != nil {
                    log.Printf("get_driving_status.sh status2 Error!!\n")
					log.Fatal(err)
				}

				// stringに変換して余分な改行を取る
				status1 := string(statuses1)
				status2 := string(statuses2)
				status1 = strings.TrimRight(status1, "\n")
				status2 = strings.TrimRight(status2, "\n")


				// status1及びstatus2が空(ファイル解析に失敗した)の場合は何もせずに次のループへ
				if status1 == "" || status2 == "" {
					continue
				}

				var now_parking_cardata string       // 車両情報（直近）
                var before_parking_cardata string    // 車両情報（一つ前）

                // 直近の車両情報取得(./driving_history/driving_history.csv 2行目)を取得
                now_parking_cardatas, err := exec.Command("bash", "-c", "awk -F',' '{print $0}' ./driving_history/driving_history.csv | sed -n 2P").Output()
                if err != nil {
                    log.Printf(" 直近の車両情報(./driving_history/driving_history.csv 2行目)を取得できませんでした Error!!\n")
                    log.Fatal(err)
                }
                // stringに変換
                now_parking_cardata = string(now_parking_cardatas)
                now_parking_cardata = strings.TrimRight(now_parking_cardata, "\n")

                // 前回の車両情報取得(./driving_history/driving_history.csv 1行目)を取得
                before_parking_cardatas, err := exec.Command("bash", "-c", "awk -F',' '{print $0}' ./driving_history/driving_history.csv | sed -n 1P").Output()
                if err != nil {
                    log.Printf(" 前回の車両情報(./driving_history/driving_history.csv 1行目)を取得できませんでした Error!!\n")
                    log.Fatal(err)
                }
                // stringに変換
                before_parking_cardata = string(before_parking_cardatas)
                before_parking_cardata = strings.TrimRight(before_parking_cardata, "\n")

                // now_parking_cardataとbefore_parking_cardataの先頭データが時刻になっている。（YYYYMMDDhhmmssttt）
                // 時間差を導き出し、アンテナ1→2の通過時速(speed)を計算する。
                now_cardatas := strings.Split(now_parking_cardata,",")
                before_cardatas := strings.Split(before_parking_cardata,",")
                time_duration, _ := date_duration(before_cardatas[0],now_cardatas[0])  //引数フォーマット → 20221014142509891

                var speed float64
                if time_duration > 0 {
                    speed = (float64(iniread.Config.Entrance_distance) /  float64(time_duration)) * 3.6   // 距離(m)÷時(s)＝速さ（秒速）　→　時速（約4倍）
                } else {
                    speed = 0
                }
                // log.Printf("Delay_time : %s - %s = %dsec -> speed : %v km\n",before_cardatas[0],now_cardatas[0],time_duration,speed)                


                /* 駐車場入口アンテナ（入口）の前でウロウロしている場合（渋滞などで連続検知するなど）
                   逆走プロセス側で逆走検知が行われているが、こちらでは駐車場に入庫したとみなす
                   ・検出している車両情報が駐車管理テーブルに存在していない場合は追加する。
                   　存在している場合は何もせずスルーする。
                */
                if status1 == "PARK" && status2 == "PARK" {

                    // すでに満空管理テーブルに登録されている場合は次のwcn番号処理へ（continue）
                    results, err := exec.Command("./script/check_parking_table.sh", wcn[0],"./parking_list/parking_table.csv").Output()
                    if err != nil {
                        log.Printf("make_parking_table.sh Error!!\n")
                    }
                    result := string(results)
                    result = strings.TrimRight(result, "\n")

                    if result == "1" { // すでに登録されている
                        continue
                    }

                    // 車両情報取得(./driving_history/driving_history.csv 1行目)を取得
                    parking_cardatas, err := exec.Command("bash", "-c", "awk -F',' '{print $0}' ./driving_history/driving_history.csv | sed -n 1P").Output()
                    if err != nil {
                        log.Printf(" 車両情報(./driving_history/driving_history.csv 1行目)を取得できませんでした。: IN無しでいきなりPARK。 Error!!\n")
                        log.Fatal(err)
                    }
                    // stringに変換
                    parking_cardata := string(parking_cardatas)
                    parking_cardata = strings.TrimRight(parking_cardata, "\n")


                    log.Printf("Directions. %s : %s --> %s\n",parking_cardata,"Continue",status1)
                    fmt.Printf("Directions. %s : %s --> %s\n",parking_cardata,"Continue",status1)

                    // parking_table.csvにレコードを追加する
                    // 同じ車両のデータがすでに登録されている場合、新しいデータで更新する。 ← この機能は必要ないが、今回はあっても害はないのでそのままにする。
                    _, err = exec.Command("./script/make_parking_table.sh", wcn[0],parking_cardata).Output()
                    if err != nil {
                        log.Printf("make_parking_table.sh Error!!\n")
                    }
                }
                
				/* 正常通行or異常通行判定
                  直近の走行履歴２件に絞って抽出しているので、IN → PARKの順番が守られているならば、それは正しく駐車場に入ってきていると判断できる。
                　IN → PARK  ： 通常通りの進行 入り口分岐から駐車場に入ってきた。　→　満空管理テーブル「parking_table.csv」に履歴レコード追加。
                　満空管理テーブル「parking_table.csv」のレコード数がそのまま駐車場に侵入した車両の数と一致する。
                */
				if status1 == "IN" && status2 == "PARK" {

                    // IN → PARK の通過履歴データと通過時速をスライスへ保存（後に、ランプ停滞中判定に利用する）
                    // in_park_listは、本関数のみで都度再作成される
                    var add_str string
                    add_str = fmt.Sprintf("%s,%v",now_parking_cardata,speed)
                    in_park_list = append(in_park_list,add_str)

                    // すでに満空管理テーブルに登録されている場合は次のwcn番号処理へ（continue）
                    results, err := exec.Command("./script/check_parking_table.sh", wcn[0],"./parking_list/parking_table.csv").Output()
                    if err != nil {
                        log.Printf("make_parking_table.sh Error!!\n")
                    }
                    result := string(results)
                    result = strings.TrimRight(result, "\n")

                    if result == "1" { // すでに登録されている
                        // log.Printf("登録済み : %s",wcn[0])
                        continue
                    }

                    log.Printf("Directions. %s : %s --> %s   flg : %s   speed : %vkm \n",now_parking_cardata,status1,status2,traffic_jam_flag,speed)

                    // parking_table.csvにレコードを追加する
                    // 同じ車両のデータがすでに登録されている場合、新しいデータで更新する。 ← この機能は必要ないが、今回はあっても害はないのでそのままにする。
                    _, err = exec.Command("./script/make_parking_table.sh", wcn[0],now_parking_cardata).Output()
                    if err != nil {
                        log.Printf("make_parking_table.sh Error!!\n")
                    }

				} else {
                    // PARK → IN  ： 駐車場から入り口分岐に逆走
                    // OUT → PARK ： 出口から駐車場に逆走
                    // PARK → OUT ： 駐車場から正しく出ていった
                    // OUT → IN   ： 一度出てから再び入ってくる正常な動線ではあるが、まだパーキングまで侵入していない

                    // 駐車場から出ていった車両は、満空管理テーブルから削除する。
                    // 逆走で駐車場から出てしまった場合も、満空管理テーブルから削除。
                    // 同様に、WCN管理テーブルから出ていった車両のWCN番号を削除。
                    if status1 == "PARK" && status2 == "OUT" ||
                       status1 == "PARK" && status2 == "IN" {

                        // すでに満空管理テーブルから削除されている場合は、次のwcn番号処理へ（continue）
                        results, err := exec.Command("./script/check_parking_table.sh", wcn[0],"./parking_list/parking_table.csv").Output()
                        if err != nil {
                            log.Printf("make_parking_table.sh Error!!\n")
                        }
                        result := string(results)
                        result = strings.TrimRight(result, "\n")
                        if result == "0" { // 存在していない
                            continue
                        }

                        log.Printf("Directions. %s : %s --> %s\n",now_parking_cardata,status1,status2)

                        // 満空管理テーブルから削除
                        _, err = exec.Command("./script/delete_parking_table.sh", wcn[0],now_parking_cardata).Output()
                        if err != nil {
                            log.Printf("delete_parking_table.sh Error!!\n")
                        }

                        // WCN管理テーブル(ac管理下)から削除
                        _, err = exec.Command("./script/delete_wcn_table.sh", wcn[0], "../ac/tc_wcn_table/WCN_table.csv").Output()
                        if err != nil {
                            log.Printf("delete_wcn_table.sh Error!!\n")
                        }

                    }


                    // 出庫時の時間から駐車時間を演算し、駐車時間管理テーブルに駐車時間と車両情報を追加する。
                    if status1 == "PARK" && status2 == "OUT" {
                        var time_duration_min int
                        now_cardatas := strings.Split(now_parking_cardata,",")
                        before_cardatas := strings.Split(before_parking_cardata,",")
                        time_duration, _ := date_duration(before_cardatas[0], now_cardatas[0]) //引数フォーマット → 20221014142509891
                        time_duration_min = time_duration / 60                                 // 秒 → 分

                        // 車両情報と駐車時間を「駐車時間管理テーブル」へ追記する。 → 統計（平均）モニタで表示するデータ
                        var str_time_durationtime_duration_min string
                        str_time_durationtime_duration_min = strconv.Itoa(time_duration_min)
                        _, err := exec.Command("./script/make_parkout_table.sh",now_parking_cardata, str_time_durationtime_duration_min).Output()
                        if err != nil {
                            log.Printf("make_parkout_table.sh Error!!\n")
                            log.Fatal(err)
                        }

                        log.Printf("Directions. %s : %s --> %s  ParkingTime : %v\n",now_parking_cardata,status1,status2,time_duration_min)

                    }
				}
            }

            // ★ランプ停滞中判定処理★
            // 前段の解析処理にて、IN→PARKの履歴データのみを、時速とともにスライスに保存している。
            // スライスの内容について、下記の条件を満たしていたら、「ランプ停滞中」とみなす。
            // ・一番直近の通過履歴が対象
            // ・通過時の時速が10km以下
            traffic_jam_flag_check(in_park_list)
            
		}
	}
}

// ランプ停滞中判定処理
// make_parking_table()処理内で作成したin_park_listを元に、ランプ停滞中判定を行う。
// 　→　in_park_list : 最新の通過履歴がIN→PARKになっている全車両リスト
// in_park_listに何も入っていなければ、何も処理せずに抜ける。
// 一番直近の通過車両履歴に記録されている時速が、指定の時速以下だった場合、ランプ停滞中とみなす。
func traffic_jam_flag_check(in_park_list []string) {

    // ソート処理用構造体
    type Record struct {
        Number int        // ソートする日付
        Data   string     // 日付以降のcsvデータ
    }
    var records []Record

    // IN→PARK履歴の車両がある場合のみ処理する
    if len(in_park_list) != 0 {

        // IN→PARK車両の検出データを、日付とそれ以外に分ける。（日付を数値データに変換）
        for _, car_data := range in_park_list {
            splitString := strings.SplitN(car_data, ",", 2)  // 文字列をカンマで分割
            num, err := strconv.Atoi(splitString[0])         // 分割した文字列の最初の要素を整数に変換
            if err != nil {
                fmt.Println(err)
                continue
            }
            records = append(records, Record{Number: num, Data: splitString[1]})                
        }

        // RecordスライスをNumberに基づいて降順にソート
        sort.Slice(records, func(i, j int) bool { return records[i].Number > records[j].Number })


        // 先頭要素(直近に入庫した車両)の時速が、指定のスピードより遅い場合、ランプ停滞中とみなす。
        // Number:[0]           Data:[0]    [1] [2]           [3]   [4]                   [5]     [6] [7]  [8]   [9]
        //        20230715000150000, RSU02, A2, 016105562150, PARK, 011934301405044328bb, 465349, c9, 500, 4034, 1380
        splitCarData := strings.Split(records[0].Data,",")
        get_speed, err := strconv.ParseFloat(splitCarData[9],64)
        if err != nil {
            log.Printf("splitCarData Atoi Error! %v",err)
        }

        if get_speed <= float64(iniread.Config.Traffic_jam_speed) {
            log.Printf("traffic_jam_flag Check speed = %vkm : %d,%s",get_speed, records[0].Number,records[0].Data)
            traffic_jam_flag = "1"  // ランプ停滞中：ON
        } else {
            traffic_jam_flag = "0"  // ランプ停滞中：OFF
        }

    }

}

// make_display_avg_csvは、平均モニタ用のcsvデータを作成する。
func make_display_avg_csv() {

	t := time.NewTicker(time.Duration(iniread.Config.Request_interval) * time.Millisecond) // config.iniの設定毎
	defer t.Stop() // タイマを止める。 <- これがないとメモリリークを起こす
    
	for {
		select {
		case <-t.C:			// タイマー時間経過した。（Config.iniに設定の時間間隔）

            // ./parking_list/parkout_table.csvを読み込み、統計（平均）モニタ用のデータを作成。
            // 作成したデータを、display公開ディレクトリ(~/opt/aps/disp_data)へcsvファイルをコピーする
            _, err := exec.Command("./script/make_avg_data.sh",traffic_jam_flag).Output()
            if err != nil {
                log.Printf("make_avg_data.sh Error!!\n")
                log.Fatal(err)
            }
        }
    }
}

// 長時間駐車管理テーブルを作成する。
// さらに、displayへの表示用連携ディレクトリにファイルを作成する
func make_parkingtime_table(input_file string, output_file string) {

    var duration_list_now []string  // 長時間駐車情報格納スライス

    // 進入車両一覧ファイルの読込
    // parking_table_info := readcsv.Read("./parking_list/parking_table.csv")
    parking_table_info := readcsv.Read(input_file)
    //log.Printf("line_cnt of parking_table.csv : %d\n",len(parking_table_info))

    /* 全ての進入車両に対して、滞留時間を調査。
       ファイル作成用のスライスにデータをセットする。
    */
    for _, csv_line := range parking_table_info {

        // ./parking_list/parking_table.csv
        // [0]                [1]    [2] [3]           [4]   [5]                   [6]     [7] [8]  [9]
        // 20230601001010000, RSU02, A2, 016072700261, PARK, 01198047052906908bbb, 465341, cb, 100, 4226
        fstatus := strings.Split(csv_line,",") // データを分解（データはcsv形式）

        // ミリ秒を含む日付文字列を作成(YYYYMMDDhhmmssxxx)
        time_stamp := get_datestr()

        // 現在時刻と通過記録時刻の差分を計算(秒)　※ミリ秒部分は無視される
        s_duration, err := date_duration(fstatus[0],time_stamp)   // 第1引数：過去　第2引数：現在
        if err != nil {
            log.Printf("date_duration Error!! : %s,%s\n",time_stamp, fstatus[0])
        }

        // 秒 → 分 へ時間変換。ただし１分未満の場合は、強制的に１分にする
        m_duration := s_duration / 60  // 秒 → 分 変換
        if m_duration <= 1 {
            m_duration = 1
        }
        

        // 2023/11/30 ここのlog.Printfとfmt.Printfは必要に応じてデバッグに利用する。logに出力するとすぐにファイルが肥大化するので注意。
        //log.Printf("NowDate : %s, ParkInDate : %s, Duration : %smin, WCN : %s",time_stamp, fstatus[0], strconv.Itoa(m_duration), fstatus[3])
        //fmt.Printf("NowDate : %s, ParkInDate : %s, Duration : %smin, WCN : %s",time_stamp, fstatus[0], strconv.Itoa(m_duration), fstatus[3])

        // 駐車時間情報格納スライスにデータを追加セット
        write_line := fstatus[0] + "," + fstatus[1] + "," + fstatus[2] + "," + fstatus[3] + "," + fstatus[4] + "," + fstatus[5] + "," + fstatus[6] + "," + fstatus[7] + "," + fstatus[8] + "," + fstatus[9] + "," + strconv.Itoa(m_duration)

        // 2023/11/30 ここのlog.Printfとfmt.Printfは必要に応じてデバッグに利用する。logに出力するとすぐにファイルが肥大化するので注意。
        //log.Printf("add duration_list : %s",write_line)
        //fmt.Printf("add duration_list : %s",write_line)

        duration_list_now = append(duration_list_now, write_line)

    }

    // duration_list_nowをcsvファイルとして保存
    if len(duration_list_now) != 0 {
        //        outputFile, err := os.Create("./parking_list/parkingtime_table.csv")
        outputFile, err := os.Create(output_file)
        if err != nil {
            fmt.Println("Error creating output file:", err)
            return
        }
        defer outputFile.Close()

        writer := csv.NewWriter(outputFile)

        for _, output_csv := range duration_list_now {

            record := strings.Split(output_csv, ",") // CSV文字列を解析して[]stringに変換
            err := writer.Write(record)              // []stringをCSVファイルに書き込む
            if err != nil {
               log.Println("Error writing to file:", err)
            }
        }
        writer.Flush()

        // 書き込み中にエラーが発生していないかチェック
        if err := writer.Error(); err != nil {
            log.Println("Error during write:", err)
        }
    }
}

// make_parkingtime_table_csvは、現在入庫している全車両の駐車時間を算出し、駐車時間テーブルを作成する。
func make_parkingtime_table_csv() {

	t := time.NewTicker(time.Duration(iniread.Config.Repeat_check_interval) * time.Second) // config.iniの設定毎
	defer t.Stop() // タイマを止める。 <- これがないとメモリリークを起こす
    
	for {
		select {
		case <-t.C:			// タイマー時間経過した。（Config.iniに設定の時間間隔）

            // 満空管理テーブルにある全車両について、何分駐車しているかを演算し別ファイルに保存する。
            make_parkingtime_table("./parking_list/parking_table.csv","./parking_list/parkingtime_table.csv")
        }
    }
}


// 指定時刻になったら何らかの処理をさせる。
// この関数はGoルーチンでコールされ、内部処理は250msec間隔で繰り返される。
func fixed_schedule() {

	t := time.NewTicker(time.Duration(250) * time.Millisecond) // config.iniの設定毎
	defer t.Stop() // タイマを止める。 <- これがないとメモリリークを起こす
    
	for {
		select {
		case <-t.C:			// 250msec経過した。

            // 現在時刻の取得
            now_date := get_datestr()                        // 現在日付を取得
            now_time_str := now_date[8:14]
            now_time,_ := strconv.Atoi(now_time_str)

            /* 指定の時刻(config.ini)になった時に何らかの処理を行う  */
            judge_time_str := iniread.Config.Path_reset_time     // 設定時刻取得
            judge_time,_ := strconv.Atoi(judge_time_str)         
            if now_time == judge_time {                          // 現在時刻が設定時刻と一致したら

                // ★設定時刻時の処理例★
                // hold_day := iniread.Config.Goback_drive_path_day // 何日分残すか、設定値を取得
                // _, err := exec.Command("./script/reset_drive_path.sh", hold_day, "./parking_list/drive_path_table.csv").Output()
                // if err != nil {
                //     log.Printf("reset_drive_path.shの処理に失敗しました  Error. : %v\n",err)
                // }
            }


		}
	}
}


/* Main */
func main() {

    // Goルーチン終了待ちオブジェクト作成
	var wg sync.WaitGroup
    wg.Add(1)                           // WaitGroupにGoルーチン登録

	go make_parking_table()             // 満空管理ファイル(./parking_list/parking_table.csv)を作成・更新

    go make_parkingtime_table_csv()     // 現在駐車中の全車両が、何分間駐車しているかを調査しファイル化する

    go make_display_avg_csv()           // 平均モニタ用csvファイル作成

    // 全てのGoルーチンが終了するまで待つ（異常がない限り終わらない予定）
    wg.Wait()

}

