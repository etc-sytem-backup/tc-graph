// make_archive_hasuda.go
// 蓮田SA向けバージョン　指定するファイルのアーカイブを作成する
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
    //	"strings"
	"time"
	"sync"

	"localhost.com/iniread"
    //	"localhost.com/readcsv"
)


// プロジェクト定数
const (
    // Nothing
)

// package変数
var log_run_path string = "0"   // 動作ログファイル格納用パス(初期値 : Off)


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

    // ワークディレクトリ群の作成
    //make_work_folder()
    
}

// ワークディレクトリ群の作成。
// ワークディレクトリ群が無い場合は新規に作成する。
func make_work_folder() error {

    // display_server用のファイル格納場所
	_, err := os.Open("./disp_data")
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


// 指定時刻になったら何らかの処理をさせる。
// この関数はGoルーチンでコールされ、内部処理は250msec間隔で繰り返される。
func runOnceAfterTime() {

	t := time.NewTicker(time.Duration(250) * time.Millisecond) // config.iniの設定毎
	defer t.Stop() // タイマを止める。 <- これがないとメモリリークを起こす

    // iniread.Config.Script_start_timeに、000000（00時00分00秒）をセットしている
    targetTime, _ := strconv.Atoi(iniread.Config.Script_start_time)
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


                // make_data_hasuda01の動作ログをアーカイブしておく
                cmd := exec.Command("./script/make_zip.sh", "../make_data_hasuda01/log/run/", "log")
                msg, err := cmd.CombinedOutput() // 標準出力と標準エラー出力の両方を取得
                if err != nil {
                    log.Printf("make_zip.shの処理に失敗しました。msg: %s ... Error: %v\n", string(msg), err)
                }

                // make_data_hasuda02の動作ログをアーカイブしておく
                cmd = exec.Command("./script/make_zip.sh", "../make_data_hasuda02/log/run/", "log")
                msg, err = cmd.CombinedOutput() // 標準出力と標準エラー出力の両方を取得
                if err != nil {
                    log.Printf("make_zip.shの処理に失敗しました。msg: %s ... Error: %v\n", string(msg), err)
                }

                // make_data_hasuda03の動作ログをアーカイブしておく
                cmd = exec.Command("./script/make_zip.sh", "../make_data_hasuda03/log/run/", "log")
                msg, err = cmd.CombinedOutput() // 標準出力と標準エラー出力の両方を取得
                if err != nil {
                    log.Printf("make_zip.shの処理に失敗しました。msg: %s ... Error: %v\n", string(msg), err)
                }
                
                // make_data_hasuda05の動作ログをアーカイブしておく
                cmd = exec.Command("./script/make_zip.sh", "../make_data_hasuda05/log/run/", "log")
                msg, err = cmd.CombinedOutput() // 標準出力と標準エラー出力の両方を取得
                if err != nil {
                    log.Printf("make_zip.shの処理に失敗しました。msg: %s ... Error: %v\n", string(msg), err)
                }

                // make_data_hasuda06の動作ログをアーカイブしておく
                cmd = exec.Command("./script/make_zip.sh", "../make_data_hasuda06/log/run/", "log")
                msg, err = cmd.CombinedOutput() // 標準出力と標準エラー出力の両方を取得
                if err != nil {
                    log.Printf("make_zip.shの処理に失敗しました。msg: %s ... Error: %v\n", string(msg), err)
                }

                // log01のRSUログをアーカイブしておく
                cmd = exec.Command("./script/make_zip.sh", "../log01/log/csv/", "csv")
                msg, err = cmd.CombinedOutput() // 標準出力と標準エラー出力の両方を取得
                if err != nil {
                    log.Printf("make_zip.shの処理に失敗しました。msg: %s ... Error: %v\n", string(msg), err)
                }

                // log02のRSUログをアーカイブしておく
                cmd = exec.Command("./script/make_zip.sh", "../log02/log/csv/", "csv")
                msg, err = cmd.CombinedOutput() // 標準出力と標準エラー出力の両方を取得
                if err != nil {
                    log.Printf("make_zip.shの処理に失敗しました。msg: %s ... Error: %v\n", string(msg), err)
                }

                // log03のRSUログをアーカイブしておく
                cmd = exec.Command("./script/make_zip.sh", "../log03/log/csv/", "csv")
                msg, err = cmd.CombinedOutput() // 標準出力と標準エラー出力の両方を取得
                if err != nil {
                    log.Printf("make_zip.shの処理に失敗しました。msg: %s ... Error: %v\n", string(msg), err)
                }

                // rsu01のRSUログをアーカイブしておく
                cmd = exec.Command("./script/make_zip.sh", "../rsu01/log/run/", "log")
                msg, err = cmd.CombinedOutput() // 標準出力と標準エラー出力の両方を取得
                if err != nil {
                    log.Printf("make_zip.shの処理に失敗しました。msg: %s ... Error: %v\n", string(msg), err)
                }

                // rsu02のRSUログをアーカイブしておく
                cmd = exec.Command("./script/make_zip.sh", "../rsu02/log/run/", "log")
                msg, err = cmd.CombinedOutput() // 標準出力と標準エラー出力の両方を取得
                if err != nil {
                    log.Printf("make_zip.shの処理に失敗しました。msg: %s ... Error: %v\n", string(msg), err)
                }

                // rsu03のRSUログをアーカイブしておく
                cmd = exec.Command("./script/make_zip.sh", "../rsu03/log/run/", "log")
                msg, err = cmd.CombinedOutput() // 標準出力と標準エラー出力の両方を取得
                if err != nil {
                    log.Printf("make_zip.shの処理に失敗しました。msg: %s ... Error: %v\n", string(msg), err)
                }

                // sbox01のRSUログをアーカイブしておく
                cmd = exec.Command("./script/make_zip.sh", "../sbox01/log/run/", "log")
                msg, err = cmd.CombinedOutput() // 標準出力と標準エラー出力の両方を取得
                if err != nil {
                    log.Printf("make_zip.shの処理に失敗しました。msg: %s ... Error: %v\n", string(msg), err)
                }

                // sbox02のRSUログをアーカイブしておく
                cmd = exec.Command("./script/make_zip.sh", "../sbox02/log/run/", "log")
                msg, err = cmd.CombinedOutput() // 標準出力と標準エラー出力の両方を取得
                if err != nil {
                    log.Printf("make_zip.shの処理に失敗しました。msg: %s ... Error: %v\n", string(msg), err)
                }

                // sbox03のRSUログをアーカイブしておく
                cmd = exec.Command("./script/make_zip.sh", "../sbox03/log/run/", "log")
                msg, err = cmd.CombinedOutput() // 標準出力と標準エラー出力の両方を取得
                if err != nil {
                    log.Printf("make_zip.shの処理に失敗しました。msg: %s ... Error: %v\n", string(msg), err)
                }

                // make_archive_hasudaのRSUログをアーカイブしておく
                cmd = exec.Command("./script/make_zip.sh", "../make_archive_hasuda/log/run/", "log")
                msg, err = cmd.CombinedOutput() // 標準出力と標準エラー出力の両方を取得
                if err != nil {
                    log.Printf("make_zip.shの処理に失敗しました。msg: %s ... Error: %v\n", string(msg), err)
                }

                // 蓮田SAで収集したデータのバックアップを実施
                cmd = exec.Command("./script/hasuda_backup.sh")
                msg, err = cmd.CombinedOutput() // 標準出力と標準エラー出力の両方を取得
                if err != nil {
                    log.Printf("hasuda_backup.shの処理に失敗しました。msg: %s ... Error: %v\n", string(msg), err)
                }

                processed = true
                
            } else if now_time < targetTime && processed {

                // 翌日のためにフラグをリセット
                processed = false
            }
		}
	}
}

/* Main */
func main() {

    // Goルーチン終了待ちオブジェクト作成
	var wg sync.WaitGroup
    wg.Add(1)                           // WaitGroupにGoルーチン登録

    go runOnceAfterTime()               // 指定時刻になったら何らかの処理をする

    // 全てのGoルーチンが終了するまで待つ（異常がない限り終わらない予定）
    wg.Wait()

}

