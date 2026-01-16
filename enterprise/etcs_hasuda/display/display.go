package main

import (
	"bytes"
	"context"
	"embed"
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"etc-system.jp/iniread"

	"github.com/pkg/sftp"
	"github.com/zserge/lorca"
	"golang.org/x/crypto/ssh"
)

// グローバル定数
const (
	ADMIN_MODE = "admin"
	USER_MODE  = "user"
)

// グローバル変数
var (
	//go:embed www
	fs                         embed.FS
	mode                       string
	ip_address                 string
	sftp_port                  string
	sftp_username              string
	sftp_password              string
	remote_dir_path            string
	main_csv_filename          string
	avg_csv_filename           string
	table_week_csv_filename    string
	table_month_csv_filename   string
	parking_time_csv_file_name string
	alert_csv_file_name        string
	passage_csv_file_name      string
	settings_csv_file_name     string
	main_interval_ms           int
	avg_interval_ms            int
	table_interval_ms          int
	alert_interval_ms          int
	passage_interval_ms        int
	rsu_connect_reset_min      int

	// 画面はグローバル
	ui lorca.UI
	ln net.Listener
)

// 初期化
func init() {

	// Config.ini読み込み
	iniread.Run()

	// 接続先IP取得
	ip_address = iniread.Config.Ip_Address

	// モードの決定
	mode = iniread.Config.Mode

	// SFTP用接続先情報取得
	sftp_port = iniread.Config.Sftp_Port
	sftp_username = iniread.Config.User_Name
	sftp_password = "qpsk5.8G"
	remote_dir_path = iniread.Config.Data_Dir
	main_csv_filename = iniread.Config.Disp_Main_Csv
	avg_csv_filename = iniread.Config.Disp_Avg_Csv
	passage_csv_file_name = iniread.Config.Disp_Passage_Csv
	table_week_csv_filename = iniread.Config.Disp_Table_Week_Csv
	table_month_csv_filename = iniread.Config.Disp_Table_Month_Csv
	parking_time_csv_file_name = iniread.Config.Disp_Parking_Time_Csv
	alert_csv_file_name = iniread.Config.Disp_Alert_Csv
	settings_csv_file_name = iniread.Config.Disp_Settings_Csv

	// 情報取得間隔
	var err error
	main_interval_ms, err = atoi(iniread.Config.Disp_Main_Interval)
	if err != nil {
		fmt.Println("Invalid interval time in config.ini")
		log.Printf("atoi error:%v", err)
	}

	avg_interval_ms, err = atoi(iniread.Config.Disp_Avg_Interval)
	if err != nil {
		fmt.Println("Invalid interval time in config.ini")
		log.Printf("atoi error:%v", err)
	}

	table_interval_ms, err = atoi(iniread.Config.Disp_Table_Interval)
	if err != nil {
		fmt.Println("Invalid interval time in config.ini")
		log.Printf("atoi error:%v", err)
	}

	alert_interval_ms, err = atoi(iniread.Config.Disp_Alert_Interval)
	if err != nil {
		fmt.Println("Invalid interval time in config.ini")
		log.Printf("atoi error:%v", err)
	}

	passage_interval_ms, err = atoi(iniread.Config.Disp_Passage_Interval)
	if err != nil {
		fmt.Println("Invalid interval time in config.ini")
		log.Printf("atoi error:%v", err)
	}

	rsu_connect_reset_min, err = atoi(iniread.Config.Rsu_connect_reset_interval)
	if err != nil {
		fmt.Println("Invalid interval time in config.ini")
		log.Printf("atoi error:%v", err)
	}

	// ログファイル保存設定
	log_setup()
	go timer_10() // 10秒タイマースタート(ログファイルのローテーション用)

	// 起動マシンの情報をログ出力
	log.Printf("NumCPU: %d\n", runtime.NumCPU())
	log.Printf("NumGoroutine: %d\n", runtime.NumGoroutine())
	log.Printf("Version: %s\n", runtime.Version())
}

// log_setup()は、各種ログファイルの保存先ディレクトリを作成します。
func log_setup() {

	// log保存ディレクトリの準備
	err := make_log_folder(iniread.Config.Bin_log_path, iniread.Config.Csv_log_path, iniread.Config.Run_log_path)
	if err != nil {
		panic(err)
	}

	// Log保存ファイル設定
	now := time.Now()
	year_val, month_val, day_val := now.Date() // 年月日を数字で取得してみる
	log_filename := fmt.Sprintf(iniread.Config.Run_log_path+"/"+"%04d%02d%02d.log", year_val, int(month_val), day_val)

	//    fmt.Printf("Run_log_path : %s\n",iniread.Config.Run_log_path)
	//    fmt.Printf("log_filename : %s\n",log_filename)

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

// 受信バイナリ保存用。CSVファイル保存用。動作ログ保存用。
// Parameters
//   - log_bin_path : 送受信ログPATH（バイナリ）
//   - log_csv_path : 送受信ログPATH
//   - log_run_path : 動作ログPATH
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

// timer_10関数は、10秒おきに評価されるタイマー関数です。
// ログファイルを1日でローテーションします。
// Parameters
//   - Nothing
//
// Returns
//   - Nothing
func timer_10() {
	t := time.NewTicker(10 * time.Second) // 10秒おきに通知
	defer t.Stop()                        // タイマを止める。 <- これがないとメモリリークを起こす

	for {
		select {
		case <-t.C:
			// 10秒経過した。
			now := time.Now()
			year_val, month_val, day_val := now.Date() // 年月日を数字で取得してみる
			//			log_filename := fmt.Sprintf(log_run_path+"%04d%02d%02d.log", year_val, int(month_val), day_val)
			log_filename := fmt.Sprintf(iniread.Config.Run_log_path+"/"+"%04d%02d%02d.log", year_val, int(month_val), day_val)

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

// RSU回線切断警告文を、指定時間経過後にクリアする
func timer_rsuConnectMessage(ui lorca.UI, interval int) {
	t := time.NewTicker(time.Duration(interval) * time.Minute) // 指定時間経過タイマー
	defer t.Stop()                                             // タイマを止める。 <- これがないとメモリリークを起こす

	for {
		select {
		case <-t.C:
			// 指定時間経過した。
			// 警告文をクリア
			log.Printf("debug timer_rsuConnectMessage : innerText is clear.\n")
			ui.Eval(fmt.Sprintf("document.querySelector('.rsu-connect').innerText = ''"))
			return // 自身を終了させる（Goルーチン終了）
		}
	}
}

// atoi関数は、文字列を数値に変換します。
// エラーハンドリングの共通化のために実装。
// Parameters
//   - s
//
// Returns
//   - int
func atoi(s string) (int, error) {
	num, err := strconv.Atoi(s)
	if err != nil {
		fmt.Printf("%s", err)
		log.Printf("%s", err)
	}
	return num, err
}

// SFTPで読み込んだファイルをスライスに展開します
func sftpCSVToSlice(remoteCSV *sftp.File) ([][]string, error) {
	reader := csv.NewReader(remoteCSV)
	records, err := reader.ReadAll()
	if err != nil {
		return [][]string{}, err
	}
	return records, nil
}

// スライスからi番目の要素を削除する
func remove(slice [][]string, i int) [][]string {
	// 削除する要素の前と後でスライスを分割し、それらを結合する
	return append(slice[:i], slice[i+1:]...)
}

// abs関数は、intの絶対値を取得します。
// math.Abs(float64)の変換を共通化するために実装。
// Parameters
//   - n
//
// Returns
//   - int
func abs(n int) int {
	return int(math.Abs(float64(n)))
}

// initLorca関数は、Lorcaの初期設定および、初期画面のロードまでを行います。。
// Parameters
//   - Nothing
//
// Returns
//   - lorca.UI, net.Listener
func initLorca() (lorca.UI, net.Listener) {

	// 画面ライブラリ(Google Chrome)用引数定義
	args := []string{}
	if runtime.GOOS == "linux" {
		args = append(args, "--class=Lorca")
	}

	// Ubuntuのバージョンアップに伴うLorcaの既知のエラーの対応
	args = append(args, "--remote-allow-origins=*")
	args = append(args, "--start-fullscreen")

	// 画面オブジェクト作成
	ui, err := lorca.New("", "", 1920, 1080, args...)
	if err != nil {
		log.Fatal(err)
	}

	// Load HTML.
	// You may also use `data:text/html,<base64>` approach to load initial HTML,
	// e.g: ui.Load("data:text/html," + url.PathEscape(html))
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}

	go http.Serve(ln, http.FileServer(http.FS(fs)))
	ui.Load(fmt.Sprintf("http://%s/www", ln.Addr()))

	return ui, ln
}

/*
SFTP用のコネクション
*/
func connectSftpTarget(target string) (*sftp.Client, error) {
	key, key_err := ssh.ParsePrivateKey([]byte(PRIVATE_KEY))
	if key_err != nil {
		return nil, key_err
	}

	config := &ssh.ClientConfig{
		User: sftp_username,
		Auth: []ssh.AuthMethod{
			// ssh.Password(sftp_password),
			ssh.PublicKeys(key),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 中間者攻撃対策を無効化
	}

	conn, dial_err := ssh.Dial("tcp", target, config)
	if dial_err != nil {
		return nil, dial_err
	}

	client, create_client_err := sftp.NewClient(conn)
	if create_client_err != nil {
		return nil, create_client_err
	}
	fmt.Printf("SFTP client has connected to %s\n", target)

	return client, nil
}

// 指定ファイルを削除する
func removeSFTPFileWithTimeOut(client *sftp.Client, remoteFilePath string, timeout int) (bool, error) {

	var err error

	// タイムアウト用のcontextを定義
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	// 結果を受け取るためのチャネル
	result := make(chan string)

	go func() {
		err = client.Remove(remoteFilePath)
		result <- "File Remove is Done"
	}()

	select {
	case <-result:
		log.Printf("debug removeSFTPFileWithTimeOut File Remove is Done!\n")
		return true, err
	case <-ctx.Done():
		fmt.Printf("Timeout with remove remote file, over %d seconds\n", timeout)
		//		log.Printf("Timeout with remove remote file, over %d seconds\n", timeout)
		ui.Close()
		os.Exit(1)
	}
	return false, err
}

// 指定ファイルを開く
func openSFTPFileWithTimeOut(client *sftp.Client, remoteFilePath string, timeout int) (*sftp.File, error) {
	var remoteFile *sftp.File
	var err error

	// タイムアウト用のcontextを定義
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	// 結果を受け取るためのチャネル
	result := make(chan string)

	go func() {
		remoteFile, err = client.Open(remoteFilePath)
		result <- "File Open is Done"
	}()

	select {
	case <-result:
		return remoteFile, err
	case <-ctx.Done():
		fmt.Printf("Timeout with opening remote file, over %d seconds\n", timeout)
		//		log.Printf("Timeout with opening remote file, over %d seconds\n", timeout)

		// 2023/09/26 SFTP接続に失敗しても、プログラムを落とさないように対応(下記の2行をコメントアウト)
		//ui.Close()
		//os.Exit(1)
	}
	return nil, fmt.Errorf("Some error occured...\n")
}

// SFTPに接続し、RSU回線断線通知ファイル（初期 : rsu_connect_false）の有無を検知
// Result
//
//	true  : RSU回線切断（通知ファイルが存在していた）
//	false : RSU回線正常（通知ファイルが無かった）
func readRSUConnect(client *sftp.Client) (bool, error) {

	rsuConnectFilePath := remote_dir_path + "/" + iniread.Config.Rsu_Connect

	// SFTP接続先のRSU回線切断通知ファイルについて存在を確認する。
	// ★注意★ rsuConnectを取得できなかった時に、defer rsuConnect.Close()を評価すると例外が発生して死ぬ。
	// err
	//   != nil : ファイルが見つからない
	//   == nil : ファイルが見つかった
	rsuConnect, err := openSFTPFileWithTimeOut(client, rsuConnectFilePath, 60)
	if err != nil {
		fmt.Printf("Failed to open remote %s : %v\n", rsuConnectFilePath, err)
		return false, err

	} else {

		log.Printf("debug : openSFTPFileWithTimeOut : File open success. : %s\n", rsuConnectFilePath)

		defer rsuConnect.Close() // ファイルポインタを取得できたら、きちんとCloseするように。

		// RSU回線切断通知ファイルを検知後、SFTP経由で通知ファイルを削除する。
		result, err := removeSFTPFileWithTimeOut(client, rsuConnectFilePath, 3)
		if err != nil {
			fmt.Printf("%s remove error : %v\n", rsuConnectFilePath, err)
			log.Printf("%s remove error : %v\n", rsuConnectFilePath, err)
		}

		if result == true {
			fmt.Printf("File remove success. : %s\n", rsuConnectFilePath)
			log.Printf("debug : removeSFTPFileWithTimeOut : File remove success. : %s\n", rsuConnectFilePath)
		}
	}

	return true, nil
}

/*
SFTPに接続し、メイン画面用データを構築する。
*/
func readMainCSV(ma_data map[string]string, client *sftp.Client) error {

	var file_get_flg bool = true // SFTP経由で、ファイルポインタを取得できたか否かの判定用

	result := make(map[string]string)
	remoteCSVPath := remote_dir_path + "/" + main_csv_filename

	// SFTP接続先のCSVファイルについて存在を確認する。
	// ★注意★ mainCSVを取得できなかった時に、defer mainCSV.Close()を評価すると例外が発生して死ぬ。
	mainCSV, err := openSFTPFileWithTimeOut(client, remoteCSVPath, 60)
	if err != nil {
		fmt.Println("Failed to open remote Main csv: ", err)
		log.Println("Failed to open remote Main csv: ", err)
		//	return err
		file_get_flg = false // ファイルポインタ取得失敗
	} else {
		defer mainCSV.Close() // ファイルポインタ「mainCSV」を取得できたら、きちんとCloseするように。
	}

	// ファイルポインタを取得できていたら、内容を収集する。
	// ファイルポインタを取得できなかった場合は、空データをセットする。
	if file_get_flg == true {

		// ファイルの情報を取得する(最終更新時刻)
		fileInfo, err := client.Stat(remoteCSVPath)
		if err != nil {
			fmt.Println("Failed to stat remote Main csv: ", err)
			log.Println("Failed to stat remote Main csv: ", err)
			return err
		}
		modTime := formatTimeJST(fileInfo.ModTime())

		// 読み込んだCSVをスライスに入れる
		records, err := sftpCSVToSlice(mainCSV)
		if err != nil {
			fmt.Println("Error reading Main csv file: ", err)
			log.Println("Error reading Main csv file: ", err)
			return err
		}

		// データチェック
		if (len(records) == 0) || (len(records[0]) != 7) {
			return fmt.Errorf("Invalid Main csv length, expected 7\n")
		}

		// Print all records(テスト用)
		for _, record := range records {
			fmt.Print("Main\t")
			for _, item := range record {
				fmt.Printf("%s\t", item)
			}
			fmt.Println()
		}

		// 最終受信時刻
		result["last-receive-time"] = modTime

		// 受信データを加工したデータを作るための変換
		l_car_cnt, err1 := atoi(records[0][0])
		if err1 != nil {
			//return err1
			l_car_cnt = 0 // エラーの場合はゼロセット
		}

		l_park_cnt, err2 := atoi(records[0][1])
		if err2 != nil {
			// return err2
			l_park_cnt = 0 // エラーの場合はゼロセット
		}

		l_pass_cnt, err3 := atoi(records[0][2])
		if err3 != nil {
			// return err3
			l_pass_cnt = 0 // エラーの場合はゼロセット
		}

		s_car_cnt, err4 := atoi(records[0][3])
		if err4 != nil {
			// return err4
			s_car_cnt = 0 // エラーの場合はゼロセット
		}

		s_park_cnt, err5 := atoi(records[0][4])
		if err5 != nil {
			// return err5
			s_park_cnt = 0 // エラーの場合はゼロセット
		}

		s_pass_cnt, err6 := atoi(records[0][5])
		if err6 != nil {
			// return err6
			s_pass_cnt = 0 // エラーの場合はゼロセット
		}

		// 駐車台数、駐車室数、パス数、電波状態
		result["l-car-cnt"] = fmt.Sprint(l_car_cnt)
		result["l-park-cnt"] = fmt.Sprint(l_park_cnt)
		result["l-pass-cnt"] = fmt.Sprint(l_pass_cnt)
		result["s-car-cnt"] = fmt.Sprint(s_car_cnt)
		result["s-park-cnt"] = fmt.Sprint(s_park_cnt)
		result["s-pass-cnt"] = fmt.Sprint(s_pass_cnt)
		result["radio-status"] = records[0][6]

		// 満空率(小数点以下切り捨て)(ゼロ除算の場合は強制的にゼロセット)
		if l_car_cnt != 0 && l_park_cnt != 0 {
			result["l-manku-ratio"] = fmt.Sprint(int(math.Round(float64(l_car_cnt) / float64(l_park_cnt) * 100)))
		} else {
			result["l-manku-ratio"] = fmt.Sprint(0)
		}

		if s_car_cnt != 0 && s_park_cnt != 0 {
			result["s-manku-ratio"] = fmt.Sprint(int(math.Round(float64(s_car_cnt) / float64(s_park_cnt) * 100)))
		} else {
			result["s-manku-ratio"] = fmt.Sprint(0)
		}

		// 満車フラグ、超過台数
		if l_park_cnt-l_car_cnt <= 0 {
			result["l-mansya-flag"] = fmt.Sprint(1)
			result["l-tyouka-cnt"] = fmt.Sprint(l_car_cnt - l_park_cnt)
		} else {
			result["l-mansya-flag"] = fmt.Sprint(0)
			result["l-tyouka-cnt"] = fmt.Sprint(0)
		}
		if s_park_cnt-s_car_cnt <= 0 {
			result["s-mansya-flag"] = fmt.Sprint(1)
			result["s-tyouka-cnt"] = fmt.Sprint(s_car_cnt - s_park_cnt)
		} else {
			result["s-mansya-flag"] = fmt.Sprint(0)
			result["s-tyouka-cnt"] = fmt.Sprint(0)
		}

	} else {

		fmt.Printf("Can't open disp_main.csv\n")

		// 取得できなかった場合は、SFTP接続日付をマスクしておく
		result["last-receive-time"] = fmt.Sprint("---")

		// 駐車台数、駐車室数、パス数
		result["l-car-cnt"] = fmt.Sprint(0)
		result["l-park-cnt"] = fmt.Sprint(0)
		result["l-pass-cnt"] = fmt.Sprint(0)
		result["s-car-cnt"] = fmt.Sprint(0)
		result["s-park-cnt"] = fmt.Sprint(0)
		result["s-pass-cnt"] = fmt.Sprint(0)

		// 満空率(小数点以下切り捨て)
		result["l-manku-ratio"] = fmt.Sprint(0)
		result["s-manku-ratio"] = fmt.Sprint(0)

		// 満車フラグ、超過台数
		result["l-mansya-flag"] = fmt.Sprint(0)
		result["l-tyouka-cnt"] = fmt.Sprint(0)

		result["s-mansya-flag"] = fmt.Sprint(0)
		result["s-tyouka-cnt"] = fmt.Sprint(0)

		// 電波状態
		result["radio-status"] = "-1"

	}

	// ma_data(の参照)に入れる
	for k, v := range result {
		ma_data[k] = v
	}

	return nil
}

/*
メイン画面を一定間隔でリフレッシュする。
*/
func refreshMaDisp(ui lorca.UI, interval int, client *sftp.Client) {
	t := time.NewTicker(time.Duration(interval) * time.Millisecond)
	defer t.Stop()

	// HTML側ボタン(大型駐車台数プラス1補正)クリック処理紐付け
	ui.Bind("plusLCarCnt", func() {
		err := sendCmd(iniread.Config.Large_Plus, client)
		if err != nil {
			log.Printf("%s", err)
			fmt.Printf("%s", err)
		}
	})
	// HTML側ボタン(大型駐車台数マイナス1補正)クリック処理紐付け
	ui.Bind("minusLCarCnt", func() {
		err := sendCmd(iniread.Config.Large_Minus, client)
		if err != nil {
			log.Printf("%s", err)
			fmt.Printf("%s", err)
		}
	})
	// HTML側ボタン(小型駐車台数プラス1補正)クリック処理紐付け
	ui.Bind("plusSCarCnt", func() {
		err := sendCmd(iniread.Config.Small_Plus, client)
		if err != nil {
			log.Printf("%s", err)
			fmt.Printf("%s", err)
		}
	})
	// HTML側ボタン(小型駐車台数マイナス1補正)クリック処理紐付け
	ui.Bind("minusSCarCnt", func() {
		err := sendCmd(iniread.Config.Small_Minus, client)
		if err != nil {
			log.Printf("%s", err)
			fmt.Printf("%s", err)
		}
	})

	// HTML側ボタン(電波発射)クリック処理紐付け
	ui.Bind("sendRadioStartCmd", func() {
		err := sendCmd(iniread.Config.Radio_Start, client)
		if err != nil {
			log.Printf("%s", err)
			fmt.Printf("%s", err)
		}
	})

	// HTML側ボタン（電波停止）クリック処理紐付け
	ui.Bind("sendRadioStopCmd", func() {
		err := sendCmd(iniread.Config.Radio_Stop, client)
		if err != nil {
			log.Printf("%s", err)
			fmt.Printf("%s", err)
		}
	})

	// 画面更新データ定義
	ma_data := make(map[string]string)

	// 画面更新ループ
	for {
		select {
		case <-t.C:

			// CSVを読みに行く
			err := readMainCSV(ma_data, client)
			if err != nil {
				fmt.Printf("%s", err)
				log.Printf("%s", err)
				break
			}

			ui.Eval(fmt.Sprintf("updateReceiveTime(\"%s\", \"main\")", ma_data["last-receive-time"]))
			ui.Eval(fmt.Sprintf("updateLCarCnt(\"%s\")", ma_data["l-car-cnt"]))
			ui.Eval(fmt.Sprintf("updateSCarCnt(\"%s\")", ma_data["s-car-cnt"]))

			ui.Eval(fmt.Sprintf("updateLParkCnt(\"%s\")", ma_data["l-park-cnt"]))
			ui.Eval(fmt.Sprintf("updateSParkCnt(\"%s\")", ma_data["s-park-cnt"]))

			ui.Eval(fmt.Sprintf("updateLMankuRatio(\"%s ％\", \"%s\")", ma_data["l-manku-ratio"], ma_data["l-mansya-flag"]))
			ui.Eval(fmt.Sprintf("updateSMankuRatio(\"%s ％\", \"%s\")", ma_data["s-manku-ratio"], ma_data["s-mansya-flag"]))

			ui.Eval(fmt.Sprintf("updateLPassCnt(\"%s\")", ma_data["l-pass-cnt"]))
			ui.Eval(fmt.Sprintf("updateSPassCnt(\"%s\")", ma_data["s-pass-cnt"]))

			ui.Eval(fmt.Sprintf("updateLMansyaflg(\"%s\")", ma_data["l-mansya-flag"]))
			ui.Eval(fmt.Sprintf("updateSMansyaflg(\"%s\")", ma_data["s-mansya-flag"]))

			ui.Eval(fmt.Sprintf("updateLTyoukaCnt(\"%s\")", ma_data["l-tyouka-cnt"]))
			ui.Eval(fmt.Sprintf("updateSTyoukaCnt(\"%s\")", ma_data["s-tyouka-cnt"]))

			ui.Eval(fmt.Sprintf("updateRadioStatus(\"%s\")", ma_data["radio-status"]))

			// RSU回線状態ファイル（rsu_connect_false）を読みに行く
			//   true  : RSU回線切断（通知ファイルが存在していた）
			//   false : RSU回線正常（通知ファイルが無かった）
			status, _ := readRSUConnect(client)
			if status == true {
				fmt.Printf("RSU Connect Error!\n")
				log.Printf("debug readRSUConnect : there is rsu_connect! %v\n", status)
				now_time_org := time.Now()
				now_time := now_time_org.Format("2006/01/02 15:04:05")
				log.Printf("debug now_time is %s\n", now_time)
				ui.Eval(fmt.Sprintf("document.querySelector('.rsu-connect').innerText = 'システム停止。もしくは車両の通行が%d分以上ありません。判定時刻:%s'", rsu_connect_reset_min, now_time))

				//ui.Eval(fmt.Sprintf("document.getElementById(\"rsu-connect\").innerText = システム停止。もしくは車両の通行が%d分以上ありません。判定時刻:%s",rsu_connect_reset_min , now_time.Format("2006/01/02 15:04:05")))
				go timer_rsuConnectMessage(ui, rsu_connect_reset_min)
			}
		}
	}
}

/*
SFTPに接続し、統計（平均）画面用データを構築する。
*/
func readAvgCSV(av_data map[string]string, client *sftp.Client) error {

	var file_get_flg bool = true // SFTP経由で、ファイルポインタを取得できたか否かの判定用

	result := make(map[string]string)
	remoteCSVPath := remote_dir_path + "/" + avg_csv_filename

	// SFTP接続先のCSVファイルについて存在を確認する。
	// ★注意★ mainCSVを取得できなかった時に、defer mainCSV.Close()を評価すると例外が発生して死ぬ。
	avgCSV, err := openSFTPFileWithTimeOut(client, remoteCSVPath, 60)
	if err != nil {
		fmt.Println("Failed to open remote Avg csv: ", err)
		log.Println("Failed to open remote Avg csv: ", err)
		// return err
		file_get_flg = false // ファイルポインタ取得失敗
	} else {
		defer avgCSV.Close() // ファイルポインタ「mainCSV」を取得できたら、きちんとCloseするように。
	}

	// ファイルポインタを取得できていたら、内容を収集する。
	// ファイルポインタを取得できなかった場合は、空データをセットする。
	if file_get_flg == true {

		// ファイルの情報を取得する(最終更新時刻)
		fileInfo, err := client.Stat(remoteCSVPath)
		if err != nil {
			fmt.Println("Failed to stat remote Avg csv: ", err)
			log.Println("Failed to stat remote Avg csv: ", err)
			return err
		}
		modTime := formatTimeJST(fileInfo.ModTime())

		// 読み込んだCSVをスライスに入れる
		records, err := sftpCSVToSlice(avgCSV)
		if err != nil {
			fmt.Println("Error reading Avg csv file: ", err)
			log.Println("Error reading Avg csv file: ", err)
			return err
		}

		// データチェック
		if (len(records) == 0) || (len(records[0]) != 7) {
			return fmt.Errorf("Invalid Avg csv length, expected 7\n")
		}

		// Print all records(テスト用)
		for _, record := range records {
			fmt.Print("Avg\t")
			for _, item := range record {
				fmt.Printf("%s\t", item)
			}
			fmt.Println()
		}

		// 最終受信時刻
		result["last-receive-time"] = modTime

		// 受信データを表示するために変換(010などを10にするため)
		jam_flg := records[0][0]
		a_parking_min_avg_day, err1 := atoi(records[0][1])
		if err1 != nil {
			//return err1
			a_parking_min_avg_day = 0
		}

		a_parking_min_avg_week, err2 := atoi(records[0][2])
		if err2 != nil {
			//return err2
			a_parking_min_avg_week = 0
		}

		l_parking_min_avg_day, err3 := atoi(records[0][3])
		if err3 != nil {
			//return err3
			l_parking_min_avg_day = 0
		}

		l_parking_min_avg_week, err4 := atoi(records[0][4])
		if err4 != nil {
			//return err4
			l_parking_min_avg_week = 0
		}

		s_parking_min_avg_day, err5 := atoi(records[0][5])
		if err5 != nil {
			//return err5
			s_parking_min_avg_day = 0
		}

		s_parking_min_avg_week, err6 := atoi(records[0][6])
		if err6 != nil {
			//return err6
			s_parking_min_avg_week = 0
		}

		// 駐車台数、駐車室数、パス数
		result["jam-flg"] = jam_flg
		result["a-parking-min-avg-day"] = fmt.Sprint(a_parking_min_avg_day)
		result["a-parking-min-avg-week"] = fmt.Sprint(a_parking_min_avg_week)
		result["l-parking-min-avg-day"] = fmt.Sprint(l_parking_min_avg_day)
		result["l-parking-min-avg-week"] = fmt.Sprint(l_parking_min_avg_week)
		result["s-parking-min-avg-day"] = fmt.Sprint(s_parking_min_avg_day)
		result["s-parking-min-avg-week"] = fmt.Sprint(s_parking_min_avg_week)

	} else {

		fmt.Printf("Can't open disp_avg.csv\n")

		// 最終受信時刻
		result["last-receive-time"] = "---"

		// 駐車台数、駐車室数、パス数
		result["jam-flg"] = "0"
		result["a-parking-min-avg-day"] = fmt.Sprint(0)
		result["a-parking-min-avg-week"] = fmt.Sprint(0)
		result["l-parking-min-avg-day"] = fmt.Sprint(0)
		result["l-parking-min-avg-week"] = fmt.Sprint(0)
		result["s-parking-min-avg-day"] = fmt.Sprint(0)
		result["s-parking-min-avg-week"] = fmt.Sprint(0)
	}

	// av_data(の参照)に入れる
	for k, v := range result {
		av_data[k] = v
	}

	return nil
}

/*
統計(平均)を一定間隔でリフレッシュする。
*/
func refreshAvDisp(ui lorca.UI, interval int, client *sftp.Client) {
	t := time.NewTicker(time.Duration(interval) * time.Millisecond)
	defer t.Stop()

	// 画面更新データ定義
	av_data := make(map[string]string)

	for {
		select {
		case <-t.C:
			// CSVを読みに行く
			err := readAvgCSV(av_data, client)
			if err != nil {
				fmt.Printf("%s", err)
				log.Printf("%s", err)
				break
			}

			ui.Eval(fmt.Sprintf("updateReceiveTime(\"%s\", \"avg\")", av_data["last-receive-time"]))
			ui.Eval(fmt.Sprintf("updateJamFlg(\"%s\")", av_data["jam-flg"]))
			ui.Eval(fmt.Sprintf("updateAParkingMinAvgDay(\"%s\")", av_data["a-parking-min-avg-day"]))
			ui.Eval(fmt.Sprintf("updateAParkingMinAvgWeek(\"%s\")", av_data["a-parking-min-avg-week"]))
			ui.Eval(fmt.Sprintf("updateLParkingMinAvgDay(\"%s\")", av_data["l-parking-min-avg-day"]))
			ui.Eval(fmt.Sprintf("updateLParkingMinAvgWeek(\"%s\")", av_data["l-parking-min-avg-week"]))
			ui.Eval(fmt.Sprintf("updateSParkingMinAvgDay(\"%s\")", av_data["s-parking-min-avg-day"]))
			ui.Eval(fmt.Sprintf("updateSParkingMinAvgWeek(\"%s\")", av_data["s-parking-min-avg-week"]))
		}
	}
}

/*
SFTPに接続し、統計（一覧）/ 統計（逆走）用データを取得する。
これらはデータの形式が全て同じため、第二引数でCSVファイルのタイプを指定
*/
func readTableOrReverseCSV(table_or_reverse_data *[][]string, csvFileName string, client *sftp.Client) error {

	var file_get_flg bool = true // SFTP経由で、ファイルポインタを取得できたか否かの判定用

	// 二次元スライスを戻す
	result := make([][]string, 0)
	remoteCSVPath := remote_dir_path + "/" + csvFileName

	// SFTP接続先のCSVファイルについて存在を確認する。
	// ★注意★ tableOrReverseCSVを取得できなかった時に、defer tableOrReverseCSV.Close()を評価すると例外が発生して死ぬ。
	tableOrReverseCSV, err := openSFTPFileWithTimeOut(client, remoteCSVPath, 60)
	if err != nil {
		fmt.Printf("Failed to open remote %s: %s\n", csvFileName, err)
		log.Printf("Failed to open remote %s: %s\n", csvFileName, err)
		//		return err
		file_get_flg = false // ファイルポインタ取得失敗
	} else {
		defer tableOrReverseCSV.Close() // ファイルポインタ「mainCSV」を取得できたら、きちんとCloseするように。
	}

	// ファイルポインタを取得できていたら、内容を収集する。
	// ファイルポインタを取得できなかった場合は、空データをセットする。
	if file_get_flg == true {

		// ファイルの情報を取得する(最終更新時刻)
		fileInfo, err := client.Stat(remoteCSVPath)
		if err != nil {
			fmt.Printf("Failed to open remote %s: %s\n", csvFileName, err)
			log.Printf("Failed to open remote %s: %s\n", csvFileName, err)
			return err
		}
		modTime := formatTimeJST(fileInfo.ModTime())

		// 読み込んだCSVをスライスに入れる
		records, err := sftpCSVToSlice(tableOrReverseCSV)
		if err != nil {
			fmt.Printf("Failed to open remote %s: %s\n", csvFileName, err)
			log.Printf("Failed to open remote %s: %s\n", csvFileName, err)
			return err
		}

		// Print all records(テスト用)
		fmt.Printf("%s\n", csvFileName)
		for _, record := range records {
			for _, item := range record {
				fmt.Printf("%s\t", item)
			}
			fmt.Println()
		}

		// データチェック
		if len(records) == 0 {
			return fmt.Errorf("%s is empty, expected 8 value per row.\n", csvFileName)
		}
		for i, record := range records {
			if len(record) != 8 {
				fmt.Printf("Invalid data size with row %d in %s, expected 8 value per row. skipped!\n", i+1, csvFileName)
				records = remove(records, i)
			}
			if len(records) == 0 {
				return fmt.Errorf("Invalid all data size in %s, expected 8 value per row.\n", csvFileName)
			}
		}

		// スライスにマップを追加していく
		for _, record := range records {
			resultRowSlice := make([]string, 0)

			// ファイルは存在したものの、ファイルの中身が想定のcsvデータではなかった（record[0]が空）なら次のデータへ。
			// 管理者によるデータクリアのタイミングによっては、空のcsvデータになっている場合がある。
			if record[0] == "" {
				continue
			}

			// 最終受信時刻
			resultRowSlice = append(resultRowSlice, modTime)

			// 受信データを変換
			densou_time, err1 := parseDensouTime(record[0])
			if err1 != nil {
				return err1
			}

			wcn := record[1]
			card_id := record[2]
			sikyoku, err2 := searchSikyokuFromCode(record[3])
			if err2 != nil {
				//                return err2
				sikyoku = "******"
			}

			youto, err3 := searchYoutoFromCode(record[4])
			if err3 != nil {
				//                return err3
				youto = "*"
			}

			bunrui_number := record[5]
			ichiren_number := record[6]
			var value string
			switch csvFileName {
			case table_week_csv_filename, table_month_csv_filename, parking_time_csv_file_name:
				v, err4 := atoi(record[7])
				if err4 != nil {
					return err4
				}
				value = fmt.Sprint(v)
			case alert_csv_file_name:
				value = record[7]
			}

			// データをマップに代入。ユーザーモードなら「****」を代入。
			resultRowSlice = append(resultRowSlice, densou_time)
			if mode == ADMIN_MODE {
				resultRowSlice = append(resultRowSlice, wcn)
				resultRowSlice = append(resultRowSlice, card_id)
				resultRowSlice = append(resultRowSlice, sikyoku)
				resultRowSlice = append(resultRowSlice, youto)
				resultRowSlice = append(resultRowSlice, bunrui_number)
				resultRowSlice = append(resultRowSlice, ichiren_number)
			} else {
				resultRowSlice = append(resultRowSlice, "****")
				resultRowSlice = append(resultRowSlice, "****")
				resultRowSlice = append(resultRowSlice, "****")
				resultRowSlice = append(resultRowSlice, "****")
				resultRowSlice = append(resultRowSlice, "****")
				resultRowSlice = append(resultRowSlice, "****")
			}
			resultRowSlice = append(resultRowSlice, value)

			// resultへ代入
			result = append(result, resultRowSlice)
		}

	} else {

		fmt.Printf("Can't open %s\n", csvFileName)

		resultRowSlice := make([]string, 0)

		// 最終受信時刻
		resultRowSlice = append(resultRowSlice, "---")

		// 受信データを変換
		//densou_time, err1 := parseDensouTime(record[0])
		densou_time := ""
		wcn := ""
		card_id := ""
		sikyoku := ""
		youto := ""
		bunrui_number := ""
		ichiren_number := ""
		value := ""

		// 空データをマップに代入。
		resultRowSlice = append(resultRowSlice, densou_time)

		resultRowSlice = append(resultRowSlice, wcn)
		resultRowSlice = append(resultRowSlice, card_id)
		resultRowSlice = append(resultRowSlice, sikyoku)
		resultRowSlice = append(resultRowSlice, youto)
		resultRowSlice = append(resultRowSlice, bunrui_number)
		resultRowSlice = append(resultRowSlice, ichiren_number)

		resultRowSlice = append(resultRowSlice, value)

		// resultへ代入
		result = append(result, resultRowSlice)

	}

	// エラーが無ければtable_or_reverse_data(の参照)に入れる
	*table_or_reverse_data = [][]string{}
	*table_or_reverse_data = append(*table_or_reverse_data, result...)

	return nil
}

/*
yyyMMddhhmmsstttをyyyy/MM/dd hh:mm:ssに変換
*/
func parseDensouTime(t string) (string, error) {
	// 2を加えて、ミリ秒部分を無視
	t = "2" + t[:len(t)-3]

	// パースする時間のレイアウト
	layoutIn := "20060102150405"

	// 出力する時間のレイアウト
	layoutOut := "2006/01/02 15:04:05"

	// パース
	tt, err := time.Parse(layoutIn, t)
	if err != nil {
		return "", err
	}

	// フォーマットして返す
	return tt.Format(layoutOut), nil
}

/*
統計（一覧）を一定間隔でリフレッシュする
*/
func refreshTableDisp(ui lorca.UI, interval int, client *sftp.Client) {
	t := time.NewTicker(time.Duration(interval) * time.Millisecond)
	defer t.Stop()

	// 画面更新データ定義
	table_week_data := make([][]string, 0)
	table_month_data := make([][]string, 0)
	pariking_time_data := make([][]string, 0)

	for {
		select {
		case <-t.C:
			// CSVを3つ読みに行く
			err := readTableOrReverseCSV(&table_week_data, table_week_csv_filename, client)
			if err != nil {
				fmt.Printf("%s", err)
				log.Printf("%s", err)
			}
			err = readTableOrReverseCSV(&table_month_data, table_month_csv_filename, client)
			if err != nil {
				fmt.Printf("%s", err)
				log.Printf("%s", err)
			}
			err = readTableOrReverseCSV(&pariking_time_data, parking_time_csv_file_name, client)
			if err != nil {
				fmt.Printf("%s", err)
				log.Printf("%s", err)
			}

			// 画面に反映する
			ui.Eval(fmt.Sprintf("updateTableWeek(%s)", formatTwoDimSliceToString(table_week_data)))
			ui.Eval(fmt.Sprintf("updateTableMonth(%s)", formatTwoDimSliceToString(table_month_data)))
			ui.Eval(fmt.Sprintf("updateParkingTime(%s)", formatTwoDimSliceToString(pariking_time_data)))
		}
	}
}

/*
Goの二次元文字列スライスを、Javascriptの引数のフォーマットに変換する
*/
func formatTwoDimSliceToString(twoDimSlice [][]string) string {
	var lines []string
	for _, row := range twoDimSlice {
		var line []string
		for _, element := range row {
			line = append(line, "\""+element+"\"")
		}
		lines = append(lines, "["+strings.Join(line, ",")+"]")
	}

	// 文字列のスライスをカンマで結合し、配列フォーマットにする
	result := "[" + strings.Join(lines, ",") + "]"
	return result
}

// 逆走検知テーブル（HTML）の更新
// 20230721 JavaScript側の処理をGo側でやろうとしたが旨く動作せず。
// 当面、JavaScript側の処理で行い、後日こっちに変更を試みるつもりでいる。
func reverse_update(ui lorca.UI, reverse_data [][]string) {
	script := `
		var tbody = document.querySelector("#disp-reverse");
		tbody.innerHTML = "";
	`

	for _, row := range reverse_data {
		script += `var rowElement = document.createElement("tr");`

		for _, cell := range row {
			script += fmt.Sprintf(`var cellElement = document.createElement("td");
									cellElement.textContent = "%s"; 
									rowElement.appendChild(cellElement);`, cell)
		}

		script += "tbody.appendChild(rowElement);"
	}

	err := ui.Eval(script)
	if err != nil {
		log.Fatal(err)
	}
}

/*
統計（逆走検知）画面を一定間隔でリフレッシュする
*/
func refreshReverseDisp(ui lorca.UI, interval int, client *sftp.Client) {
	t := time.NewTicker(time.Duration(interval) * time.Millisecond)
	defer t.Stop()

	// データ送信用の関数バインド
	ui.Bind("sendAlertResetCmd", func() {
		err := sendCmd(iniread.Config.Alert_Reset, client)
		if err != nil {
			log.Printf("%s", err)
			fmt.Printf("%s", err)
		}
	})

	// 画面更新データ定義
	reverse_data := make([][]string, 0)

	// 画面更新ループ
	for {
		select {
		case <-t.C:
			// CSVを読みに行く
			err := readTableOrReverseCSV(&reverse_data, alert_csv_file_name, client)
			if err != nil {
				fmt.Printf("%s", err)
				log.Printf("%s", err)
			}

			ui.Eval(fmt.Sprintf("updateReverse(%s)", formatTwoDimSliceToString(reverse_data)))
			//reverse_update(ui, reverse_data)
		}
	}
}

/*
SFTPに接続し、断面交通量画面用データを構築する。
*/
func readPassageCSV(passage_data map[string]string, client *sftp.Client) error {

	var file_get_flg bool = true // SFTP経由で、ファイルポインタを取得できたか否かの判定用

	result := make(map[string]string)
	remoteCSVPath := remote_dir_path + "/" + passage_csv_file_name

	// SFTP接続先のCSVファイルについて存在を確認する。
	passageCSV, err := openSFTPFileWithTimeOut(client, remoteCSVPath, 60)
	if err != nil {
		fmt.Println("Failed to open remote Passage Count csv: ", err)
		log.Println("Failed to open remote Passage Count csv: ", err)
		// return err
		file_get_flg = false // ファイルポインタ取得失敗
	} else {
		defer passageCSV.Close()
	}

	// ファイルポインタを取得できていたら、内容を収集する。
	// ファイルポインタを取得できなかった場合は、空データをセットする。
	if file_get_flg == true {

		// ファイルの情報を取得する(最終更新時刻)
		fileInfo, err := client.Stat(remoteCSVPath)
		if err != nil {
			fmt.Println("Failed to stat Passage Count csv: ", err)
			log.Println("Failed to stat Passage Count csv: ", err)
			return err
		}
		modTime := formatTimeJST(fileInfo.ModTime())

		// 読み込んだCSVをスライスに入れる
		records, err := sftpCSVToSlice(passageCSV)
		if err != nil {
			fmt.Println("Error reading Passage Count file: ", err)
			log.Println("Error reading Passage Count file: ", err)
			return err
		}

		// データチェック
		if (len(records) == 0) || (len(records[0]) != 5) {
			return fmt.Errorf("Invalid Passage Count csv length, expected 5\n")
		}

		// Print all records(テスト用)
		for _, record := range records {
			fmt.Print("Passage Count\t")
			for _, item := range record {
				fmt.Printf("%s\t", item)
			}
			fmt.Println()
		}

		// 最終受信時刻
		result["last-receive-time"] = modTime

		// 受信データを表示するために変換(010などが入っていた時に10として扱うため)
		antenna1, err1 := atoi(records[0][0])
		if err1 != nil {
			antenna1 = 0
		}
		antenna2, err2 := atoi(records[0][1])
		if err2 != nil {
			antenna2 = 0
		}
		antenna3, err3 := atoi(records[0][2])
		if err3 != nil {
			antenna3 = 0
		}

		// アンテナごとの断面交通量
		result["antenna1"] = fmt.Sprint(antenna1)
		result["antenna2"] = fmt.Sprint(antenna2)
		result["antenna3"] = fmt.Sprint(antenna3)
	} else {

		fmt.Printf("Can't open disp_passage_count.csv\n")

		// 最終受信時刻
		result["last-receive-time"] = "---"

		// アンテナごとの断面交通量
		result["antenna1"] = fmt.Sprint(0)
		result["antenna2"] = fmt.Sprint(0)
		result["antenna3"] = fmt.Sprint(0)
	}

	// av_data(の参照)に入れる
	for k, v := range result {
		passage_data[k] = v
	}

	return nil
}

/*
断面交通量を一定間隔でリフレッシュする。
*/
func refreshPassageDisp(ui lorca.UI, interval int, client *sftp.Client) {
	t := time.NewTicker(time.Duration(interval) * time.Millisecond)
	defer t.Stop()

	// データ送信用の関数バインド
	ui.Bind("sendPassageResetCmd", func() {
		err := sendCmd(iniread.Config.Passage_Reset, client)
		if err != nil {
			log.Printf("%s", err)
			fmt.Printf("%s", err)
		}
	})

	// 画面更新データ定義
	passage_data := make(map[string]string)

	for {
		select {
		case <-t.C:
			// CSVを読みに行く
			err := readPassageCSV(passage_data, client)
			if err != nil {
				fmt.Printf("%s", err)
				log.Printf("%s", err)
				break
			}

			ui.Eval(fmt.Sprintf("updateReceiveTime(\"%s\", \"passage\")", passage_data["last-receive-time"]))
			ui.Eval(fmt.Sprintf("updatePassageCnt1(\"%s\")", passage_data["antenna1"]))
			ui.Eval(fmt.Sprintf("updatePassageCnt2(\"%s\")", passage_data["antenna2"]))
			ui.Eval(fmt.Sprintf("updatePassageCnt3(\"%s\")", passage_data["antenna3"]))
		}
	}
}

/*
SFTPに接続し、設定画面用データを構築する。
*/
func readSettingsCSV(passage_data map[string]string, client *sftp.Client) error {

	var file_get_flg bool = true // SFTP経由で、ファイルポインタを取得できたか否かの判定用

	result := make(map[string]string)
	remoteCSVPath := remote_dir_path + "/" + settings_csv_file_name

	// SFTP接続先のCSVファイルについて存在を確認する。
	settingsCSV, err := openSFTPFileWithTimeOut(client, remoteCSVPath, 60)
	if err != nil {
		fmt.Println("Failed to open remote Settings csv: ", err)
		log.Println("Failed to open remote Settings csv: ", err)
		// return err
		file_get_flg = false // ファイルポインタ取得失敗
	} else {
		defer settingsCSV.Close()
	}

	// ファイルポインタを取得できていたら、内容を収集する。
	// ファイルポインタを取得できなかった場合は新規ファイルを作成する。
	if file_get_flg == true {

		// ファイルの情報を取得する(最終更新時刻)
		fileInfo, err := client.Stat(remoteCSVPath)
		if err != nil {
			fmt.Println("Failed to stat Settings csv: ", err)
			log.Println("Failed to stat Settings csv: ", err)
			return err
		}
		modTime := formatTimeJST(fileInfo.ModTime())

		// 読み込んだCSVをスライスに入れる
		records, err := sftpCSVToSlice(settingsCSV)
		if err != nil {
			fmt.Println("Error reading Settings file: ", err)
			log.Println("Error reading Settings file: ", err)
			return err
		}

		// データチェック
		if (len(records) == 0) || (len(records[0]) != 5) {
			return fmt.Errorf("Invalid Settings csv length, expected 5\n")
		}

		// Print all records(テスト用)
		for _, record := range records {
			fmt.Print("Settings\t")
			for _, item := range record {
				fmt.Printf("%s\t", item)
			}
			fmt.Println()
		}

		// 最終受信時刻
		result["last-receive-time"] = modTime

		// 受信データを表示するために変換(010などが入っていた時に10として扱うため)
		l_park_offset, err1 := atoi(records[0][0])
		if err1 != nil {
			l_park_offset = 0
		}
		s_park_offset, err2 := atoi(records[0][1])
		if err2 != nil {
			s_park_offset = 0
		}

		// 設定
		result["l-park-offset"] = fmt.Sprint(l_park_offset)
		result["s-park-offset"] = fmt.Sprint(s_park_offset)
	} else {

		fmt.Printf("Can't open disp_setting.csv\n")

		// 最終受信時刻
		result["last-receive-time"] = "---"

		// アンテナごとの断面交通量
		result["l-park-offset"] = fmt.Sprint(0)
		result["s-park-offset"] = fmt.Sprint(0)

		// 新規ファイル作成
		file, err := client.Create(remoteCSVPath)
		if err != nil {
			return fmt.Errorf("Failed to create Settings file: %v", err)
		}
		defer file.Close()

		// 初期値が入ったファイルを作成
		_, err = file.Write(bytes.NewBufferString("0,0,0,0,0").Bytes())
		if err != nil {
			return fmt.Errorf("Failed to write to file: %v", err)
		}
	}

	// av_data(の参照)に入れる
	for k, v := range result {
		passage_data[k] = v
	}

	return nil
}

/*
設定画面の値を基にファイルを作成する
*/
func createSettingsCSV(l_park_offset string, s_park_offset string, client *sftp.Client) {
	remoteCSVPath := remote_dir_path + "/" + settings_csv_file_name

	file, err := client.Create(remoteCSVPath)
	if err != nil {
		fmt.Println("Failed to create Settings file: ", err)
		log.Println("Failed to create Settings file: ", err)
		return
	}
	defer file.Close()

	settings := l_park_offset + "," + s_park_offset + ",0,0,0"
	_, err = file.Write(bytes.NewBufferString(settings).Bytes())
	if err != nil {
		fmt.Println("Failed to update Settings file: ", err)
		log.Println("Failed to update Settings file: ", err)
		return
	}
}

/*
設定画面を更新する。
*/
func updateSettingsDisp(ui lorca.UI, client *sftp.Client) {
	// 画面更新データ定義
	settings_data := make(map[string]string)
	err := readSettingsCSV(settings_data, client)
	if err != nil {
		fmt.Printf("%s", err)
		log.Printf("%s", err)
		return
	}
	ui.Eval(fmt.Sprintf("updateReceiveTime(\"%s\", \"settings\")", settings_data["last-receive-time"]))
	ui.Eval(fmt.Sprintf("updateLParkingOffset(\"%s\")", settings_data["l-park-offset"]))
	ui.Eval(fmt.Sprintf("updateSParkingOffset(\"%s\")", settings_data["s-park-offset"]))
}

/*
設定画面を表示する。
*/
func showSettingsDisp(ui lorca.UI, client *sftp.Client) {
	// 初期表示用の関数バインド
	ui.Bind("initParkingOffsetSettings", func() {
		updateSettingsDisp(ui, client)
	})

	// データ送信用の関数バインド
	ui.Bind("sendParkingOffsetSettings", func(l_park_offset string, s_park_offset string) {
		// 設定ファイルを作成する
		createSettingsCSV(l_park_offset, s_park_offset, client)
		updateSettingsDisp(ui, client)
	})
}

/*
time.TimeをJST時刻にしてフォーマットする
*/
func formatTimeJST(t time.Time) string {
	const format = "2006/01/02 15:04:05"
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		// WindowsではLoadLocationが動かないため追加
		jst = time.FixedZone("JST", 9*60*60)
	}
	return t.In(jst).Format(format)
}

/*
APServerの指定ディレクトリに、電波発射、電波停止、逆走警報クリアを連絡するテキストファイルを作成する。
iniread.Config.Large_Plus		：大型駐車数プラス1補正
iniread.Config.Large_Minus	：大型駐車数マイナス1補正
iniread.Config.Small_Plus		：小型駐車数プラス1補正
iniread.Config.Small_Minus	：小型駐車数マイナス1補正
iniread.Config.Radio_Start	：電波発射
iniread.Config.Radio_Stop		：電波停止
iniread.Config.Alert_Reset	：逆走警報クリア
iniread.Config.Passage_Reset：断面交通量リセット
*/
func sendCmd(cmd string, client *sftp.Client) error {
	switch cmd {
	case iniread.Config.Large_Plus: // 大型駐車数プラス1補正
		cmdFile, err := client.Create(remote_dir_path + "/" + iniread.Config.Large_Plus)
		if err != nil {
			fmt.Println("Failed to create remote Large Plus file: ", err)
			log.Println("Failed to create remote Large Plus file: ", err)
			return err
		}
		fmt.Println("Large Plus file created!")
		defer cmdFile.Close()
	case iniread.Config.Large_Minus: // 大型駐車数マイナス1補正
		cmdFile, err := client.Create(remote_dir_path + "/" + iniread.Config.Large_Minus)
		if err != nil {
			fmt.Println("Failed to create remote Large Minus file: ", err)
			log.Println("Failed to create remote Large Minus file: ", err)
			return err
		}
		fmt.Println("Large Minus file created!")
		defer cmdFile.Close()
	case iniread.Config.Small_Plus: // 小型駐車数プラス1補正
		cmdFile, err := client.Create(remote_dir_path + "/" + iniread.Config.Small_Plus)
		if err != nil {
			fmt.Println("Failed to create remote Small Plus file: ", err)
			log.Println("Failed to create remote Small Plus file: ", err)
			return err
		}
		fmt.Println("Small Plus file created!")
		defer cmdFile.Close()
	case iniread.Config.Small_Minus: // 小型駐車数マイナス1補正
		cmdFile, err := client.Create(remote_dir_path + "/" + iniread.Config.Small_Minus)
		if err != nil {
			fmt.Println("Failed to create remote Small Minus file: ", err)
			log.Println("Failed to create remote Small Minus file: ", err)
			return err
		}
		fmt.Println("Small Minus file created!")
		defer cmdFile.Close()

	case iniread.Config.Radio_Start: // 電波発射
		cmdFile, err := client.Create(remote_dir_path + "/" + iniread.Config.Radio_Start)
		if err != nil {
			fmt.Println("Failed to create remote Radio Start file: ", err)
			log.Println("Failed to create remote Radio Start file: ", err)
			return err
		}
		fmt.Println("Radio Start file created!")
		defer cmdFile.Close()
	case iniread.Config.Radio_Stop: // 電波停止
		cmdFile, err := client.Create(remote_dir_path + "/" + iniread.Config.Radio_Stop)
		if err != nil {
			fmt.Println("Failed to create remote Radio Stop file: ", err)
			log.Println("Failed to create remote Radio Stop file: ", err)
			return err
		}
		fmt.Println("Radio Stop file created!")
		defer cmdFile.Close()
	case iniread.Config.Alert_Reset: // 逆走警報クリア
		cmdFile, err := client.Create(remote_dir_path + "/" + iniread.Config.Alert_Reset)
		if err != nil {
			fmt.Println("Failed to create remote Alert Reset file: ", err)
			log.Println("Failed to create remote Alert Reset file: ", err)
			return err
		}
		fmt.Println("Alert Reset file created!")
		defer cmdFile.Close()
	case iniread.Config.Passage_Reset: // 断面交通量リセット
		cmdFile, err := client.Create(remote_dir_path + "/" + iniread.Config.Passage_Reset)
		if err != nil {
			fmt.Println("Failed to create remote Passage Count Reset file: ", err)
			log.Println("Failed to create remote Passage Count Reset file: ", err)
			return err
		}
		fmt.Println("Passage Count Reset file created!")
		defer cmdFile.Close()
	}
	return nil
}

func main() {

	// SFTP用接続クライアント作成
	sftp_target := ip_address + ":" + sftp_port
	fmt.Printf("SFTP Connect Start... : %s\n", sftp_target)
	sftp_client, sftp_err := connectSftpTarget(sftp_target)
	if sftp_err != nil {
		fmt.Printf("%s\n", sftp_err)
		log.Printf("%s\n", sftp_err)
		time.Sleep(1 * time.Second)
		return
	}
	defer sftp_client.Close()
	fmt.Printf("SFTP Connect OK : %s\n", sftp_target)

	// モードをコンソールに表示
	fmt.Printf("mode: %s\n", mode)

	// 画面の初期化
	ui, ln = initLorca()
	defer ui.Close()
	defer ln.Close()

	// 画面を指定間隔で更新
	go refreshMaDisp(ui, main_interval_ms, sftp_client)         // メインモニタ
	go refreshAvDisp(ui, avg_interval_ms, sftp_client)          // 統計（平均）
	go refreshTableDisp(ui, table_interval_ms, sftp_client)     // 統計（一覧）
	go refreshReverseDisp(ui, alert_interval_ms, sftp_client)   // 逆走検知
	go refreshPassageDisp(ui, passage_interval_ms, sftp_client) // 断面交通量
	go showSettingsDisp(ui, sftp_client)                        // 設定

	// 画面が閉じられると終了
	func() {
		<-ui.Done()
		os.Exit(0)
	}()
}
