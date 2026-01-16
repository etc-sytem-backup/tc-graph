// make_data_hasuda05.go
// 蓮田SA向けバージョン　RSU回線断線検知処理
package main

import (
    "io/ioutil"
	"fmt"
	"log"
	"os"
    "sort"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
	"sync"

    "crypto/tls"
    "net/smtp"

	"localhost.com/iniread"
    //	"localhost.com/readcsv"
)


// プロジェクト定数
const (
    user     = "alert-mail@etc-system.jp"  // SMTPユーザー名
    password = "nozawana55"                // SMTPパスワード
    rcpt     = "hasuda@etc-system.jp"      // 送信先アドレス
    host     = "sv10460.xserver.jp:465"    // SMTPサーバー

    // user     = "kintaka@etc-system.jp"  // SMTPユーザー名
    // password = "kinChan55"              // SMTPパスワード
    // rcpt     = "kintaka@etc-system.jp"  // 送信先アドレス
    // host     = "sv10460.xserver.jp:465" // SMTPサーバー
)

// package変数
var log_run_path string = "0"                         // 動作ログファイル格納用パス(初期値 : Off)
var CheckStartDate = make(map[string]time.Time)       // RSU接続監視用：死活監視ファイルの最終更新時間
var ScCheckStartDate = make(map[string]time.Time)     // Sc通知監視用：Sc通知ファイルの最終更新時間
var rsuKeys = []string{"RSU01", "RSU02", "RSU03"}     // アンテナ識別用の名前
var sboxKeys = []string{"SBOX01", "SBOX02", "SBOX03"} // SBOX Port識別用の名前
var rsuDisconnected = make(map[string]bool)           // 接続状態を追跡するための変数（各RSU毎）
var sboxDisconnected = make(map[string]bool)          // SBOXのSc通知受信状態を追跡するための変数（各SBOX Port毎）
var lastModifiedTimes = make(map[string]time.Time)    // RSU接続監視用：死活監視ファイルの最終更新時間（前回時刻）
var sclastModifiedTimes = make(map[string]time.Time)  // Sc通知監視用：Sc通知ファイルの最終更新時間（前回時刻）
var rsuAlertSent = make(map[string]time.Time)         // 各RSUの最後のアラート送信時刻を記録
var sboxAlertSent = make(map[string]time.Time)        // SBOXのSc通知最後のアラート送信時刻を記録

// package構造体
type RsuStatus struct {
    FileDetected bool           // RSU接続状態ステータス true / false
    LastModifiedTime time.Time  // 死活監視応答ファイルの最終更新時間
    IsDisconnected bool         // RSU接続状態ステータス true / false  <-- 2023/11/23 修正前に利用していた変数。元に戻すことも考えて残しておく。
    FileName       string       // RSU死活監視要求応答ファイル名
}

type SboxStatus struct {
    FileDetected bool           // SboxのSc通知ステータス true / false
    LastModifiedTime time.Time  // Sc通知ファイルの最終更新時間
    IsDisconnected bool         // SboxのSc通知ステータス true / false  <-- 2023/11/23 修正前に利用していた変数。元に戻すことも考えて残しておく。
    FileName       string       // Sc通知ファイル名
}

// 指定のメールアドレスにメール送信する
// msg : メール本文
func send_mail(msg string) {

    server := host
    body := msg

    // TLS config
    tlsconfig := &tls.Config{
        InsecureSkipVerify: true,
        ServerName:         host,
    }

    // ログイン情報を用意する
    auth := smtp.PlainAuth("", user, password, host)

    // TLSで通信するためのコネクションを用意する
    con, err := tls.Dial("tcp", server, tlsconfig)
    if err != nil {
        fmt.Printf("tls.Dial Error : %v",err)
        log.Printf("Dial Error : %v",err)
    }
    // TLSのコネクションでSMTP接続する
    c, err := smtp.NewClient(con, host)
    if err != nil {
        fmt.Printf("smtp.NewClient Error : %v",err)
        log.Printf("smtp.NewClient Error : %v",err)
    }
    if err = c.Auth(auth); err != nil {
        fmt.Printf("c.Auth Error : %v",err)
        log.Printf("c.Auth Error : %v",err)
    }
    if err = c.Mail(user); err != nil {
        fmt.Printf("c.Mail Error : %v",err)
        log.Printf("c.Mail Error : %v",err)
    }
    if err = c.Rcpt(rcpt); err != nil {
        fmt.Printf("c.Rcpt Error : %v",err)
        log.Printf("c.Rcpt Error : %v",err)
    }
    w, err := c.Data()
    if err != nil {
        fmt.Printf("c.Data Error : %v",err)
        log.Printf("c.Data Error : %v",err)
    }

    message := "From: " + user + "\r\n"
    message += "To: " + rcpt + "\r\n"
    message += "Subject:" + "ETC-System RSU Connect Alert" + "\r\n"
    message += "\n" + body

    _, err = w.Write([]byte(message))
    if err != nil {
        fmt.Printf("w.Write Error : %v",err)
        log.Printf("w.Write Error : %v",err)
    }
    defer w.Close()

    //    log.Printf("Send Mail : Body -> %s\n",message)

    c.Quit()
}


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

// RSU回線状況ポーリング
func timer_rsu_connect_chk_old() {
    old_ans_filename := make(map[string]string)
	t := time.NewTicker(1 * time.Second) // 1秒おきに通知
	defer t.Stop() // タイマを止める。 <- これがないとメモリリークを起こす

	for {
		select {
		case <-t.C: // 1秒経過した。

            // 各RSUの回線状態と死活監視応答履歴ファイル名を取得
            rsuStatusMap := rsu_connect_check()

            // 各RSUの回線状態を確認
            for rsuKey, rsuStatus := range rsuStatusMap {
                fmt.Printf("%s connection check result: %v\n", rsuKey, rsuStatus.IsDisconnected)

                // rsu_connect_check()の戻り値がtrueであればMailを送信する。
                // 但し、検出ファイル名が前回と一緒の場合は、Mail送信しない。
                if rsuStatus.IsDisconnected && old_ans_filename[rsuKey] != rsuStatus.FileName {

                    // Mail送信
                    send_mail(rsuStatus.FileName)
                    fmt.Printf("SendMail : %v", rsuStatus.FileName)

                    // 死活監視応答ファイル名をバックアップ
                    old_ans_filename[rsuKey] = rsuStatus.FileName

                    // RSU切断通知ファイルを作成する
                    _, err := exec.Command("./script/make_rsu_connect_chk.sh").Output()
                    if err != nil {
                        log.Printf("make_rsu_connect_chk.sh Error!!: %v\n",err)
                    }
                }
            }
		}
	}
}

// RSU回線状況ポーリング
func timer_rsu_connect_chk() {
    t := time.NewTicker(1 * time.Second) // 1秒おきに無限ループ処理
    defer t.Stop()

    for {
        select {
        case <-t.C: // 1秒経過した
            rsuStatusMap := rsu_connect_check()

            for rsuKey, rsuStatus := range rsuStatusMap {

                fmt.Printf("rsuDisconnected[%s] = %v, rsuStatus.FileDetected = %v\n", rsuKey,rsuDisconnected[rsuKey], rsuStatus.FileDetected)
                
                // 接続切断と判断された後、死活監視応答ファイルのタイムスタンプが更新又は新規作成された場合、通信が回復したとみなし、回復メールを送信する。
                if rsuDisconnected[rsuKey] == true && rsuStatus.FileDetected == true{

                    // Mail送信
                    var send_message string
                    switch(rsuKey) {
                    case "RSU01":
                        send_message = fmt.Sprintf("アンテナ１号機との接続が回復しました。") 
                    case "RSU02":
                        send_message = fmt.Sprintf("アンテナ２号機との接続が回復しました。") 
                    case "RSU03":
                        send_message = fmt.Sprintf("アンテナ３号機との接続が回復しました。") 
                    }
                    fmt.Printf("Send Mail Message : %s\n", send_message)
                    log.Printf("Send Mail Message : %s\n", send_message)
                    send_mail(send_message)

                    rsuDisconnected[rsuKey] = false  // 状態をリセットします。
                    fmt.Printf("rsuDisconnected[%s] <-- %v\n",rsuKey,rsuDisconnected[rsuKey])
                    log.Printf("rsuDisconnected[%s] <-- %v\n",rsuKey,rsuDisconnected[rsuKey])
                    
                }

                // 死活監視応答ファイルが存在したら、最終更新時刻を取得（更新）します。
                if rsuStatus.FileDetected {
                    CheckStartDate[rsuKey] = rsuStatus.LastModifiedTime
                } 

                //log.Printf("CheckStartDate[%s] : %v\n", rsuKey, rsuStatus.LastModifiedTime)
            }
        }
    }
}

// Sc通知受信状況ポーリング
func timer_sc_receive_chk() {
    t := time.NewTicker(1 * time.Second) // 1秒おきに無限ループ処理
    defer t.Stop()

    for {
        select {
        case <-t.C: // 1秒経過した
            sboxStatusMap := sc_receive_check()

            for sboxKey, sboxStatus := range sboxStatusMap {

                fmt.Printf("sboxDisconnected[%s] = %v, sboxStatus.FileDetected = %v\n", sboxKey,sboxDisconnected[sboxKey], sboxStatus.FileDetected)
                //log.Printf("sboxDisconnected[%s] = %v, sboxStatus.FileDetected = %v\n", sboxKey,sboxDisconnected[sboxKey], sboxStatus.FileDetected)
                
                // Sc通知未確認と判断された後、Sc通知ファイルのタイムスタンプが更新又は新規作成された場合、車両の通過を確認したとみなし、回復メールを送信する。
                if sboxDisconnected[sboxKey] == true && sboxStatus.FileDetected == true{

                    // Mail送信
                    var send_message string
                    switch(sboxKey) {
                    case "SBOX01":
                        send_message = fmt.Sprintf("アンテナ１号機の車両通過を確認しました。") 
                    case "SBOX02":
                        send_message = fmt.Sprintf("アンテナ２号機の車両通過を確認しました。") 
                    case "SBOX03":
                        send_message = fmt.Sprintf("アンテナ３号機の車両通過を確認しました。") 
                    }
                    fmt.Printf("Send Mail Message : %s\n", send_message)
                    fmt.Printf("sboxDisconnected[%s] = %v, sboxStatus.FileDetected = %v\n", sboxKey,sboxDisconnected[sboxKey], sboxStatus.FileDetected)
                    log.Printf("Send Mail Message : %s\n", send_message)
                    log.Printf("sboxDisconnected[%s] = %v, sboxStatus.FileDetected = %v\n", sboxKey,sboxDisconnected[sboxKey], sboxStatus.FileDetected)
                    send_mail(send_message)

                    sboxDisconnected[sboxKey] = false  // 状態をリセットします。
                    fmt.Printf("sboxDisconnected[%s] <-- %v\n",sboxKey,sboxDisconnected[sboxKey])
                }

                // Sc通知監視応答ファイルが存在したら、最終更新時刻を取得（更新）します。
                if sboxStatus.FileDetected {
                    ScCheckStartDate[sboxKey] = sboxStatus.LastModifiedTime
                } 

                //log.Printf("ScCheckStartDate[%s] : %v\n", sboxKey, sboxStatus.LastModifiedTime)
                //fmt.Printf("ScCheckStartDate[%s] : %v\n", sboxKey, sboxStatus.LastModifiedTime)
            }
        }
    }
}

// 各RSUの死活監視ファイル最終更新時刻を監視し、切断されていると判断した場合はメールを送信する。
func timer_send_mail() {
    //    t := time.NewTicker(time.Duration(iniread.Config.Connect_chk_interval) * time.Minute)
    t := time.NewTicker(1 * time.Second) // 1秒おきに無限ループ処理
    defer t.Stop()

    for {
        select {
        case <-t.C:

            // 各RSUの死活監視ファイル最終更新時刻に対し、現在時刻が設定監視時間（分）を超過しているか調査
            // 超過している場合は、アラートメールを送信する
            for rsuKey, startDate := range CheckStartDate {
                elapsed_time := time.Since(startDate).Minutes()

                timeSinceLastAlert := time.Since(rsuAlertSent[rsuKey]).Minutes()
                
                //log.Printf("rsuKey = %s, elapsed_time : %vmin, Connect_chk_interval : %d\n", rsuKey, elapsed_time, iniread.Config.Connect_chk_interval)
                if elapsed_time >= float64(iniread.Config.Connect_chk_interval) && timeSinceLastAlert >= float64(iniread.Config.Connect_chk_interval) {

                    // Mail送信
                    var send_message string
                    switch(rsuKey) {
                    case "RSU01":
                        send_message = fmt.Sprintf("アンテナ１号機との接続が%d分間ありません。",iniread.Config.Connect_chk_interval) 
                    case "RSU02":
                        send_message = fmt.Sprintf("アンテナ２号機との接続が%d分間ありません。",iniread.Config.Connect_chk_interval) 
                    case "RSU03":
                        send_message = fmt.Sprintf("アンテナ３号機との接続が%d分間ありません。",iniread.Config.Connect_chk_interval) 
                    }
                    log.Printf("AH通知 Send Mail Message[%s] : %s\n", rsuKey, send_message)
                    fmt.Printf("AH通知 Send Mail Message[%s] : %s\n", rsuKey, send_message)
                    
                    send_mail(send_message)

                    // RSU切断通知ファイルを作成する
                    _, err := exec.Command("./script/make_rsu_connect_chk.sh").Output()
                    if err != nil {
                        log.Printf("make_rsu_connect_chk.sh Error!!: %v\n",err)
                    }

                    // 該当アンテナとの接続が切断されたとマークします。
                    rsuDisconnected[rsuKey] = true
                    fmt.Printf("rsuDisconnected[%s] <-- %v\n",rsuKey,rsuDisconnected[rsuKey])

                    rsuAlertSent[rsuKey] = time.Now()  // 最後のアラート送信時刻を更新
                }
            }
        }
    }
}


// SBOXからのSC通知ファイル最終更新時刻を監視し、切断されていると判断した場合はメールを送信する。
func timer_send_mail_sc() {
    //    t := time.NewTicker(time.Duration(iniread.Config.Connect_chk_interval) * time.Minute)
    t := time.NewTicker(1 * time.Second) // 1秒おきに無限ループ処理
    defer t.Stop()

    for {
        select {
        case <-t.C:

            // SBOXからのSc通知ファイル最終更新時刻に対し、現在時刻が設定監視時間（分）を超過しているか調査
            // 超過している場合は、アラートメールを送信する
            for sboxKey, startDate := range ScCheckStartDate {
                elapsed_time := time.Since(startDate).Minutes()

                timeSinceLastAlert := time.Since(sboxAlertSent[sboxKey]).Minutes()
                
                //log.Printf("sboxKey = %s, elapsed_time : %vmin, Sc_receive_interval : %d\n", sboxKey, elapsed_time, iniread.Config.Sc_receive_interval)
                if elapsed_time >= float64(iniread.Config.Sc_receive_interval) && timeSinceLastAlert >= float64(iniread.Config.Sc_receive_interval) {

                    // Mail送信
                    var send_message string
                    switch(sboxKey) {
                    case "SBOX01":
                        send_message = fmt.Sprintf("アンテナ１号機の車両通過検知が%d分間ありません。",iniread.Config.Sc_receive_interval) 
                    case "SBOX02":
                        send_message = fmt.Sprintf("アンテナ２号機の車両通過検知が%d分間ありません。",iniread.Config.Sc_receive_interval) 
                    case "SBOX03":
                        send_message = fmt.Sprintf("アンテナ３号機の車両通過検知が%d分間ありません。",iniread.Config.Sc_receive_interval) 
                    }
                    log.Printf("Sc通知 Send Mail Message[%s] : %s\n", sboxKey, send_message)
                    fmt.Printf("Sc通知 Send Mail Message[%s] : %s\n", sboxKey, send_message)
                    send_mail(send_message)

                    // 2023/12/08 RSU切断通知と同様に、Sc通知未確認ファイルを作成するか？
                    // 現在は、RSU切断通知のみ対応。（下記は、RSU切断通知ファイル作成処理のまま）
                    // 結果、メインモニタ上にアンテナとの接続がされていない通知メッセージが表示される。（接続切断通知）
                    _, err := exec.Command("./script/make_rsu_connect_chk.sh").Output()
                    if err != nil {
                        log.Printf("make_rsu_connect_chk.sh Error!!: %v\n",err)
                    }

                    // 該当アンテナとの接続が切断されたとマークします。
                    sboxDisconnected[sboxKey] = true
                    fmt.Printf("sboxDisconnected[%s] <-- %v\n",sboxKey,sboxDisconnected[sboxKey])

                    sboxAlertSent[sboxKey] = time.Now()  // 最後のアラート送信時刻を更新
                }
            }
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
func fixed_schedule() {

	t := time.NewTicker(time.Duration(250) * time.Millisecond) // config.iniの設定毎
	defer t.Stop() // タイマを止める。 <- これがないとメモリリークを起こす
    
	for {
		select {
		case <-t.C:			// 250msec経過した。

            // ★設定時間と現在時間が一致したときに何かをさせたい場合は、ここに処理を記述する★

            // // 現在時刻の取得
            // now_date := get_datestr()                        // 現在日付を取得
            // now_time_str := now_date[8:14]
            // now_time,_ := strconv.Atoi(now_time_str)

            // /* 駐車パス管理テーブルリセット(指定秒の間、複数回処理される可能性大)  */
            // judge_time_str := iniread.Config.Path_reset_time
            // judge_time,_ := strconv.Atoi(judge_time_str)         // 駐車パス台数のリセット時刻を取得
            // if now_time == judge_time {                          // 現在時刻が設定時刻と一致したら
            //     hold_day := iniread.Config.Goback_drive_path_day // 何日分残すか、設定値を取得
            //     _, err := exec.Command("./script/reset_drive_path.sh", hold_day, "./parking_list/drive_path_table.csv").Output()
            //     if err != nil {
            //         log.Printf("reset_drive_path.shの処理に失敗しました  Error. : %v\n",err)
            //     }
            // }
		}
	}
}

// ファイル情報操作用
type FileInfoList []os.FileInfo

// 
func (f FileInfoList) Len() int {
    return len(f)
}

// ファイルの更新時間取得
func (f FileInfoList) Less(i, j int) bool {
    return f[i].ModTime().Before(f[j].ModTime())
}

// ファイル情報入れ替え
func (f FileInfoList) Swap(i, j int) {
    f[i], f[j] = f[j], f[i]
}

// RSUからの死活監視応答を調査する
func getLastModifiedFile(directory, pattern string) (os.FileInfo, error) {
    files, err := ioutil.ReadDir(directory)
    if err != nil {
        return nil, err
    }

    var filteredFiles FileInfoList
    for _, file := range files {
        if strings.Contains(file.Name(), pattern) {
            filteredFiles = append(filteredFiles, file)
        }
    }

    sort.Sort(sort.Reverse(filteredFiles))

    if len(filteredFiles) == 0 {
        return nil, nil
    }

    return filteredFiles[0], nil
}

// 時間差10分判定
// 引数の時間と現在時間を比較し、10分以上の時間差であればtrue、そうで無ければfalseを返す
func isTimeDifferenceMoreThan10Mins(modTime time.Time) bool {
    fmt.Printf("LastTime : %.2fmin  Comparison time : %vmin\n",time.Since(modTime).Minutes(), float64(iniread.Config.Connect_chk_interval))
    return time.Since(modTime).Minutes() >= float64(iniread.Config.Connect_chk_interval)
}

// RSU回線接続チェックを行う。
// result_map : 各RSUの回線状態と死活監視応答履歴ファイル名を格納したマップ。キーは"RSU01"、"RSU02"、"RSU03"。値はRsuStatus構造体。
func rsu_connect_check_old() map[string]RsuStatus {
    directories := []string{iniread.Config.Rsu01_ah, iniread.Config.Rsu02_ah, iniread.Config.Rsu03_ah}
    patterns := []string{iniread.Config.Find_a1, iniread.Config.Find_a2, iniread.Config.Find_a3}
    rsuKeys := []string{"RSU01", "RSU02", "RSU03"}

    result_map := make(map[string]RsuStatus)

    // 検索対象ディレクトリの数だけループ
    for i := range directories {
        file, err := getLastModifiedFile(directories[i], patterns[i])
        if err != nil {
            fmt.Printf("Error reading directory: %v\n", err)
            continue
        }

        // RSUの死活監視ファイルが存在していたら、現在時間と比較。
        if file != nil {
            fmt.Printf("FileName : %v\n", file.Name())
            status := RsuStatus{
                IsDisconnected: isTimeDifferenceMoreThan10Mins(file.ModTime()),
                FileName:       file.Name(),
            }
            result_map[rsuKeys[i]] = status
        }
    }

    return result_map
}

// RSU回線接続チェックを行う。
// result_map : 各RSUの回線状態と死活監視応答履歴ファイル名を格納したマップ。キーは"RSU01"、"RSU02"、"RSU03"。値はRsuStatus構造体。
func rsu_connect_check_old2() map[string]RsuStatus {
    directories := []string{iniread.Config.Rsu01_ah, iniread.Config.Rsu02_ah, iniread.Config.Rsu03_ah}  // 検索対象ディレクトリ
    patterns := []string{iniread.Config.Find_a1, iniread.Config.Find_a2, iniread.Config.Find_a3}        // 検索対象ファイルのファイル名パターン
    rsuKeys := []string{"RSU01", "RSU02", "RSU03"}                                                      // アンテナ識別用の名前

    result_map := make(map[string]RsuStatus)

    for i := range directories {
        file, err := getLastModifiedFile(directories[i], patterns[i])
        if err != nil {
            fmt.Printf("Error reading directory: %v\n", err)
            continue
        }

        // 死活監視ファイルが見つかった
        if file != nil {
            result_map[rsuKeys[i]] = RsuStatus{
                FileDetected:     true,              // 死活監視ファイル検出した:true
                LastModifiedTime: file.ModTime(),    // 最終更新時間
                FileName:       file.Name(),         // 死活監視ファイル名
            }
        }
        fmt.Printf("getLastModifiedFile() => %v, FileName : %s\n",file, result_map[rsuKeys[i]].FileName)
        fmt.Printf("result_map[rsuKeys[%d]] : %v, %v, %v\n",i,result_map[rsuKeys[i]].FileDetected, result_map[rsuKeys[i]].LastModifiedTime, result_map[rsuKeys[i]].FileName)
    }

    return result_map
}

// RSU回線接続チェックを行う。
// result_map : 各RSUの回線状態と死活監視応答履歴ファイル名を格納したマップ。キーは"RSU01"、"RSU02"、"RSU03"。値はRsuStatus構造体。
func rsu_connect_check() map[string]RsuStatus {
    directories := []string{iniread.Config.Rsu01_ah, iniread.Config.Rsu02_ah, iniread.Config.Rsu03_ah}
    patterns := []string{iniread.Config.Find_a1, iniread.Config.Find_a2, iniread.Config.Find_a3}
    rsuKeys := []string{"RSU01", "RSU02", "RSU03"}

    result_map := make(map[string]RsuStatus)

    for i := range directories {
        file, err := getLastModifiedFile(directories[i], patterns[i])
        if err != nil {
            fmt.Printf("Error reading directory: %v\n", err)
            continue
        }

        // 死活監視ファイルが見つかった場合の処理
        if file != nil {
            lastModifiedTime := file.ModTime()

            // 前回の最終更新時間と比較します。
            fileDetected := lastModifiedTime.After(lastModifiedTimes[rsuKeys[i]])

            // 結果をマップに保存します。
            result_map[rsuKeys[i]] = RsuStatus{
                FileDetected:     fileDetected,      // 前回より新しい場合のみtrue
                LastModifiedTime: lastModifiedTime,  // 最終更新時間
                FileName:         file.Name(),       // 死活監視ファイル名
            }

            // 最終更新時間を更新します。
            if fileDetected {
                lastModifiedTimes[rsuKeys[i]] = lastModifiedTime
            }
        } else {
            // ファイルが見つからなかった場合、FileDetectedをfalseに設定
            result_map[rsuKeys[i]] = RsuStatus{
                FileDetected: false,
            }
        }

        fmt.Printf("getLastModifiedFile() => %v, FileName : %s\n", file, result_map[rsuKeys[i]].FileName)
        fmt.Printf("result_map[rsuKeys[%d]] : %v, %v, %v\n", i, result_map[rsuKeys[i]].FileDetected, result_map[rsuKeys[i]].LastModifiedTime, result_map[rsuKeys[i]].FileName)
    }

    return result_map
}

// Sc通知受信チェックを行う。
// result_map : SBOXのSc通知受信状態とSc通知受信ファイル名を格納したマップ。キーは"SBOX01"、"SBOX02"、"SBOX03"。値はSboxStatus構造体。
func sc_receive_check() map[string]SboxStatus {
    directories := []string{iniread.Config.Sbox01_sc, iniread.Config.Sbox02_sc, iniread.Config.Sbox03_sc}
    patterns := []string{iniread.Config.Find_car1, iniread.Config.Find_car2, iniread.Config.Find_car3}
    sboxKeys := []string{"SBOX01", "SBOX02", "SBOX03"}

    result_map := make(map[string]SboxStatus)

    for i := range directories {
        file, err := getLastModifiedFile(directories[i], patterns[i])
        if err != nil {
            fmt.Printf("Error reading directory: %v\n", err)
            continue
        }

        // Sc通知ファイルが見つかった場合の処理
        if file != nil {
            sclastModifiedTime := file.ModTime()

            // 前回の最終更新時間と比較します。
            fileDetected := sclastModifiedTime.After(sclastModifiedTimes[sboxKeys[i]])

            // 結果をマップに保存します。
            result_map[sboxKeys[i]] = SboxStatus{
                FileDetected:     fileDetected,       // 前回より新しい場合のみtrue
                LastModifiedTime: sclastModifiedTime, // 最終更新時間
                FileName:         file.Name(),        // Sc通知ファイル名
            }

            // 最終更新時間を更新します。
            if fileDetected {
                sclastModifiedTimes[sboxKeys[i]] = sclastModifiedTime
            }
        } else {
            // ファイルが見つからなかった場合、FileDetectedをfalseに設定
            result_map[sboxKeys[i]] = SboxStatus{
                FileDetected: false,
            }
        }

        fmt.Printf("getLastModifiedFile() => %v, FileName : %s\n", file, result_map[sboxKeys[i]].FileName)
        fmt.Printf("result_map[sboxKeys[%d]] : %v, %v, %v\n", i, result_map[sboxKeys[i]].FileDetected, result_map[sboxKeys[i]].LastModifiedTime, result_map[sboxKeys[i]].FileName)
    }

    return result_map
}

/* Main */
func main() {

    // Goルーチン終了待ちオブジェクト作成
	var wg sync.WaitGroup
    wg.Add(1)                           // WaitGroupにGoルーチン登録

    // 各RSUの初期CheckStartDateと、SBOX用のScCheckStartDateを設定
    // 初期値は現在時刻とする。
    for _, rsuKey := range rsuKeys {
        CheckStartDate[rsuKey] = time.Now()
    }
    for _, sboxkey := range sboxKeys {
        ScCheckStartDate[sboxkey] = time.Now()
    }

    
    log.Printf("CheckStartDate Initial : %v", CheckStartDate)
    log.Printf("ScCheckStartDate Initial : %v", ScCheckStartDate)

    go timer_rsu_connect_chk()             // RSU回線接続切断検知処理
    go timer_send_mail()                   // メール送信処理(RSU海鮮切断検知)

    go timer_sc_receive_chk()              // Sc通知検知処理（車両走行が一定時間検知されないとみなす）
    go timer_send_mail_sc()                // メール送信処理(Sc通知未検知：車両通貨未検知)
    
    // 全てのGoルーチンが終了するまで待つ（異常がない限り終わらない予定）
    wg.Wait()

}

