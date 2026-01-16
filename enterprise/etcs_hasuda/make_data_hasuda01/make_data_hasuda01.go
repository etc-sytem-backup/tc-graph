// make_data_hasuda01.go
// 蓮田SA向けバージョン　メイン画面用データ作成
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
    //    "encoding/csv"

    //    "crypto/tls"
    //    "net/smtp"

	"localhost.com/iniread"
	"localhost.com/readcsv"
)


// プロジェクト定数
const (
    // Nothing.
)

// package変数
var log_run_path string = "0"      // 動作ログファイル格納用パス(初期値 : Off)
var duration_list_show []string    // 長時間車両一覧スライス（表示用）
var traffic_jam_flag string        // ランプ停滞中フラグ On:1 / Off:0
var disp_radio_status string = "0" // 電波発信中ステータス
var large_parking_offset = 0       // 大型車駐車車両数オフセット（メインモニタプラスマイナスボタン用）
var other_parking_offset = 0       // 大型車以外車両数オフセット（メインモニタプラスマイナスボタン用）
var large_setting_offset = 0       // 大型車駐車車両数オフセット（画面端末 設定画面）
var other_setting_offset = 0       // 大型車以外車両数オフセット（画面端末 設定画面）

/*
   各種ログファイル保存用ディレクトリ作成
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



// 満空管理ファイルを作成する。
// その他、display表示に連携する管理ファイルを作成する。
// この関数は、Goルーチンとして、一定間隔毎に繰り返し処理される（繰り返し間隔はConfig.iniに設定）
func make_parking_table() {

    // var wcn []string
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
                    log.Printf("make_driving_history.sh Error!! wcn: %s\n",wcn[0])
					//log.Fatal(err)
                    continue
				}

                // ★この処理は無理にbashにさせなくても。。。。
				// driving_history.csvの行数を取得（直近２レコードあれば２行のはず）
				lines, err := exec.Command("./script/get_line_count.sh").Output()
				if err != nil {
                    log.Printf("get_line_count.sh Error!!\n")
					//log.Fatal(err)
                    continue
				}

				// intに変換
				s := string(lines)
				s = strings.TrimRight(s, "\n")
				linecnt, err := strconv.Atoi(s)
				if err != nil {
                    log.Printf("strconv.Atoi Error!!\n")
					//log.Fatal(err)
                    continue
				}

                // -----------------------------------------------------------------------------
				// 1行しかない場合、通過アンテナがPARKだった場合は、駐車場入庫と見なして処理する。
                // -----------------------------------------------------------------------------
				if linecnt == 1 {

                    //debug
                    // fmt.Printf("debug linecnt = %d\n",linecnt)
                    // fmt.Printf("debug linecnt = %d\n",linecnt)
                    
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

                //debug
                // fmt.Printf("debug linecnt = %d\n",linecnt)
                // log.Printf("debug linecnt = %d\n",linecnt)


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
                    //debug
                    //fmt.Printf("debug check error so continue.... status1:%s status2:%s\n",status1,status2)

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
                log.Printf("Delay_time : %s - %s = %dsec -> speed : %v km\n",before_cardatas[0],now_cardatas[0],time_duration,speed)                

                //debug
                //fmt.Printf("debug check status  status1:%s -> status2:%s\n",status1,status2)

                
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

                    log.Printf("Directions. %s : %s --> %s\n",now_parking_cardata,status1,status2)

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
                    // 同様に、出ていった車両のWCN番号をWCN管理テーブルから削除。
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


                    // 駐車場から出た車両について、10分(初期値)以上駐車したかを調べ、10分以内の退場だった場合は駐車パス管理ファイルに車両情報を追加する。
                    if status1 == "PARK" && status2 == "OUT" {
                        var time_duration_min int
                        now_cardatas := strings.Split(now_parking_cardata,",")
                        before_cardatas := strings.Split(before_parking_cardata,",")
                        time_duration, _ := date_duration(before_cardatas[0], now_cardatas[0]) //引数フォーマット → 20221014142509891
                        time_duration_min = time_duration / 60                                 // 秒 → 分

                        // 駐車場利用時間が設定時間よりも短い場合は、駐車パステーブルに車両情報を追加する。
                        if time_duration < iniread.Config.Duration_time {
                            
                            results, err := exec.Command("./script/make_path_table.sh",now_parking_cardata).Output()
                            if err != nil {
                                log.Printf("make_path_table.sh Error!!\n")
                                log.Fatal(err)
                            }
                            result := string(results)
                            result = strings.TrimRight(result, "\n")
                            if result == "0" {  // 新規登録した
                                log.Printf("Parking Path. %s : %d min\n",now_parking_cardata, time_duration_min)                                
                            }
                        }
                    }
				}
            }
		}
	}
}

// display_server連携用ファイルを作成します。
func make_display_main_csv() {

    var large_in_parking_cnt string // 駐車台数（大型車）
    var other_in_parking_cnt string // 駐車台数（大型車以外）
    var large_parking_space string  // 駐車室数（大型車）
    var other_parking_space string  // 駐車室数（大型車以外）
    var large_drivepath_cnt string  // 駐車パス台数（大型車）
    var other_drivepath_cnt string  // 駐車パス台数（大型車以外）
    var radio_status string         // 電波発射中ステータス
    var stack int                   // 計算用一時変数

	t := time.NewTicker(time.Duration(iniread.Config.Request_interval) * time.Millisecond) // config.iniの設定毎
	defer t.Stop() // タイマを止める。 <- これがないとメモリリークを起こす
    
	for {
		select {
		case <-t.C:			// タイマー時間経過した。（Config.iniに設定の時間間隔）

            // 駐車台数（大型車両）
            result, err := exec.Command("./script/get_carcount_large.sh","./parking_list/parking_table.csv").Output()
            if err != nil {
                log.Printf("delete_parking_table.sh Error!! %v\n",err)
            }
            large_in_parking_cnt = string(result)
            large_in_parking_cnt = strings.TrimRight(large_in_parking_cnt, "\n")

            // 取得した「駐車台数（大型車両）」とオフセットを加算（設定画面のオフセット「disp_setting.csv」の内容も加える）
            stack,_ = strconv.Atoi(large_in_parking_cnt)
            stack = stack + large_parking_offset + large_setting_offset
            large_in_parking_cnt = strconv.Itoa(stack)
            
            // 駐車台数（大型車両以外）
            result, err = exec.Command("./script/get_carcount_other.sh","./parking_list/parking_table.csv").Output()
            if err != nil {
                log.Printf("delete_parking_table.sh Error!! %v\n",err)
            }
            other_in_parking_cnt = string(result)
            other_in_parking_cnt = strings.TrimRight(other_in_parking_cnt, "\n")

            // 取得した「駐車台数（大型車両以外）」とオフセットを加算（設定画面のオフセット「disp_setting.csv」の内容も加える）
            stack,_ = strconv.Atoi(other_in_parking_cnt)
            stack = stack + other_parking_offset + other_setting_offset
            other_in_parking_cnt = strconv.Itoa(stack)

            // 車室数（大型車）
            large_parking_space = strconv.Itoa(iniread.Config.Large_parking_space)


            // 車室数（大型車両以外）
            other_parking_space = strconv.Itoa(iniread.Config.Other_parking_space)


            // 駐車パス台数（大型車両）
            result, err = exec.Command("./script/get_carcount_large.sh","./parking_list/drive_path_table.csv").Output()
            if err != nil {
                log.Printf("delete_parking_table.sh Error!! %v\n",err)
            }
            large_drivepath_cnt = string(result)
            large_drivepath_cnt = strings.TrimRight(large_drivepath_cnt, "\n")

            // 駐車パス台数（大型車両以外）
            result, err = exec.Command("./script/get_carcount_other.sh","./parking_list/drive_path_table.csv").Output()
            if err != nil {
                log.Printf("delete_parking_table.sh Error!! %v\n",err)
            }
            other_drivepath_cnt = string(result)
            other_drivepath_cnt = strings.TrimRight(other_drivepath_cnt, "\n")

            // 電波発射中ステータス
            radio_status = disp_radio_status
            
            
            // display用のcsvファイル（disp_main.csv）を作成。
            // display_server参照ディレクトリ(~/opt/aps/disp_data)へcsvファイルをコピーする
            _, err = exec.Command("./script/make_main_data.sh",
                large_in_parking_cnt,
                other_in_parking_cnt,
                large_parking_space,
                other_parking_space,
                large_drivepath_cnt,
                other_drivepath_cnt,
                radio_status).Output()
            if err != nil {
                log.Printf("make_main_data.sh Error!! %v\n",err)
            }
		}
	}
}


// 指定時刻になったら何らかの処理をさせる。
// この関数はGoルーチンでコールされ、内部処理は250msec間隔で繰り返される。
func runOnceAfterTime() {

	t := time.NewTicker(time.Duration(250) * time.Millisecond) // config.iniの設定毎
	defer t.Stop() // タイマを止める。 <- これがないとメモリリークを起こす

    // 2023/11/26 iniread.Config.Path_reset_timeに、000000（00時00分00秒）をセットしている
    targetTime, _ := strconv.Atoi(iniread.Config.Path_reset_time)
	processed := false
    
	for {
		select {
		case <-t.C:			// 250msec経過した。
            
            // 現在時刻の取得
            now_date := get_datestr()                        // 現在日付を取得
            now_time_str := now_date[8:14]
            now_time,_ := strconv.Atoi(now_time_str)

            // 設定時刻を過ぎたら、一度だけ処理を実施
            if now_time >= targetTime && !processed {

                // 通過パス表示項目のリセット処理
                hold_day := iniread.Config.Goback_drive_path_day // 何日分残すか、設定値を取得
                _, err := exec.Command("./script/reset_drive_path.sh", hold_day, "./parking_list/drive_path_table.csv").Output()
                if err != nil {
                    log.Printf("reset_drive_path.shの処理に失敗しました  Error. : %v\n",err)
                }

                processed = true
                
            } else if now_time < targetTime && processed {

                // 翌日のためにフラグをリセット
                processed = false
            }
		}
	}
}

// radio_control関数は、display端末から作成される無線制御ファイルを監視し、開始/停止をSBOXへ指示出しする。
// 開始を検出した場合は、画面表示用のアンテナステータスフラグをtrueに、停止を検出した場合はfalseにする。
// アンテナステータスフラグは、make_display_main_csv()内で参照する。
func radio_control() {
	t := time.NewTicker(time.Duration(1000) * time.Millisecond) // 1秒間隔で監視
	defer t.Stop() // タイマを止める。 <- これがないとメモリリークを起こす
    
	for {
		select {
		case <-t.C:			// 1sec経過した。

            out, err := exec.Command("./script/check_radio_control.sh").Output()
            if err != nil {
                log.Printf("check_radio_control.shの処理に失敗しました  Error.")
            }

            // radio_start = "1" / radio_stop = "0"
            result := strings.TrimSpace(string(out))
            //fmt.Printf("radio_status : %s\n",result)
            switch(result) {
            case "1":
                disp_radio_status = "1"
            case "0":
                disp_radio_status = "0"
            default:
                // 何もしない
            }
        }
    }
}

// 表示用端末（display）が作成する、大型車両の駐車車両数調整要求ファイルを検出します。
// 検出した要求ファイルに応じて、オフセットカウンタを増減させます。
// オフセットカウンタは、make_display_main_csv()関数で参照されます。
func timer_large_offset_check() {
	t := time.NewTicker(time.Duration(300) * time.Millisecond) // 300msec間隔で監視
	defer t.Stop() // タイマを止める。 <- これがないとメモリリークを起こす
    
	for {
		select {
		case <-t.C:			// 300msec経過した。

            out, err := exec.Command("./script/large_offset_check.sh").Output()
            if err != nil {
                log.Printf("large_offset_check.shの処理に失敗しました  Error.")
            }

            // 加算 = "0" / 減算 = "1" / 何もしない = "2"
            result := strings.TrimSpace(string(out))
            switch(result) {
            case "0":
                large_parking_offset = large_parking_offset + 1
                //fmt.Printf("large_parking_offset : %d\n",large_parking_offset)
            case "1":
                large_parking_offset = large_parking_offset - 1
                //fmt.Printf("large_parking_offset : %d\n",large_parking_offset)
            default:
                // 何もしない
            }
        }
    }
}

// 表示用端末（display）が作成する、大型車両以外の駐車車両数調整要求ファイルを検出します。
// 検出した要求ファイルに応じて、オフセットカウンタを増減させます。
// オフセットカウンタは、make_display_main_csv()関数で参照されます。
func timer_other_offset_check() {
	t := time.NewTicker(time.Duration(300) * time.Millisecond) // 300msec間隔で監視
	defer t.Stop() // タイマを止める。 <- これがないとメモリリークを起こす
    
	for {
		select {
		case <-t.C:			// 300msec経過した。

            out, err := exec.Command("./script/other_offset_check.sh").Output()
            if err != nil {
                log.Printf("other_offset_check.shの処理に失敗しました  Error.")
            }

            // 加算 = "0" / 減算 = "1" / 何もしない = "2"
            result := strings.TrimSpace(string(out))
            switch(result) {
            case "0":
                other_parking_offset = other_parking_offset + 1
                //fmt.Printf("other_parking_offset : %d\n",other_parking_offset)
            case "1":
                other_parking_offset = other_parking_offset - 1
                //fmt.Printf("other_parking_offset : %d\n",other_parking_offset)
            default:
                // 何もしない
            }
        }
    }
}

// timer_setting_offset_checkは、画面端末アプリが作成する「disp_setting.csv」を監視し、内容をオフセット値として取得します。
func timer_setting_offset_check() {
	t := time.NewTicker(time.Duration(300) * time.Millisecond) // 300msec間隔で監視
	defer t.Stop() // タイマを止める。 <- これがないとメモリリークを起こす
    
	for {
		select {
		case <-t.C:			// 300msec経過した。

            out, err := exec.Command("./script/check_disp_setting_csv.sh").Output()
            if err != nil {
                log.Printf("check_disp_setting_csv.shの処理に失敗しました  Error.")
            }
            result := strings.TrimSpace(string(out))
            
            // csv文字列をバラし、オフセット値とする
            parts := strings.Split(result, ",")

            // 設定画面で作成したオフセット値を取得（大型車両駐車台数）
            large_setting_offset, err = strconv.Atoi(parts[0])
            if err != nil {
                log.Println("変換エラー:", err)
            }
            
            // 設定画面で作成したオフセット値を取得（大型車両以外駐車台数）
            other_setting_offset , err = strconv.Atoi(parts[1])
            if err != nil {
                log.Println("変換エラー:", err)
            }
        }
    }
}

/* Main */
func main() {

    // Goルーチン終了待ちオブジェクト作成
	var wg sync.WaitGroup
    wg.Add(1)                           // WaitGroupにGoルーチン登録

    // メイン処理
    go timer_large_offset_check()       // 画面端末アプリからの駐車場在車数調整要求ファイル検知(大型車駐車台数)
    go timer_other_offset_check()       // 画面端末アプリからの駐車場在車数調整要求ファイル検知(大型車以外駐車台数)
    go timer_setting_offset_check()     // 画面端末アプリからの駐車場在車数オフセットファイル検知（設定画面）
    
	go make_parking_table()             // 満空管理ファイル(./parking_list/parking_table.csv)を作成・更新

    go runOnceAfterTime()               // 指定時刻になったら何らかの処理をする

    go make_display_main_csv()          // メインモニタ表示用（displayプロセス用）csvファイル作成

    go radio_control()                  // 無線制御の開始と停止を制御

    // 全てのGoルーチンが終了するまで待つ（異常がない限り終わらない予定）
    wg.Wait()

}

