package main

import (
	"fmt"
	"log"

	"os"
	"os/exec"
	"time"

	//	"localhost.com/application_counter/csvcontroller"
	"localhost.com/iniread"
)

/* Package Global var  */
var (
    log_bin_path string = "./"        // バイナリファイル格納用パス
    log_run_path string = "./"        // 動作ログファイル格納用パス
    log_csv_path string = "./"        // CSVファイル格納用パス
)


/*
   要求テキストファイル作成
   config.iniに設定しているRSU及びsboxのディレクトリに、引数で指示されている要求テキストファイルを作成する
*/
func make_reqfile(machine string, num int, cmd string) {

	switch machine {
	case "SBOX":
		switch num {
		case 1:
			var cmd string = iniread.Config.Sbox01_path + "/S" + cmd
			err := exec.Command("touch", cmd).Run()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Command : %s\n", cmd)
		case 2:
			var cmd string = iniread.Config.Sbox02_path + "/S" + cmd
			err := exec.Command("touch", cmd).Run()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Command : %s\n", cmd)
		case 3:
			var cmd string = iniread.Config.Sbox03_path + "/S" + cmd
			err := exec.Command("touch", cmd).Run()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Command : %s\n", cmd)
		case 4:
			var cmd string = iniread.Config.Sbox04_path + "/S" + cmd
			err := exec.Command("touch", cmd).Run()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Command : %s\n", cmd)
		}

	case "RSU":
		switch num {
		case 1:
			var cmd string = iniread.Config.Rsu01_path + "/R" + cmd
			err := exec.Command("touch", cmd).Run()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Command : %s\n", cmd)
		case 2:
			var cmd string = iniread.Config.Rsu02_path + "/R" + cmd
			err := exec.Command("touch", cmd).Run()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Command : %s\n", cmd)
		case 3:
			var cmd string = iniread.Config.Rsu03_path + "/R" + cmd
			err := exec.Command("touch", cmd).Run()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Command : %s\n", cmd)
		case 4:
			var cmd string = iniread.Config.Rsu04_path + "/R" + cmd
			err := exec.Command("touch", cmd).Run()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Command : %s\n", cmd)
		}
	}
}

/*
   300msec毎に、RSU01のCSVファイルをACにコピーする
*/
func cploop_rsu01_csv() {
	t := time.NewTicker(300 * time.Millisecond) // 300msecおきに通知
	for {
		select {
		case <-t.C:
			// 300msec経過した。

			/* tc_f2で、
			   　 YYYYMMDDHHMMSS_WCN_A1.csv
			      WCN_XXXXXXXXXXXX.csv
			     を作成するようにした。
			     さらに、oki_sbox01/scriptに、csvを収集するためのスクリプトも配置した。
			     なので、この関数で特に処理を行う必要がなくなった。
			     んが、今後の為に関数は残しておく。
			*/

		}
	}
	t.Stop() // タイマを止める。 <- これがないとメモリリークを起こす
}

/*
   300msec毎に、RSU02のCSVファイルをACにコピーする
*/
func cploop_rsu02_csv() {
	t := time.NewTicker(300 * time.Millisecond) // 300msecおきに通知
	for {
		select {
		case <-t.C:
			// 300msec経過した。

			// CSVファイルをACにコピー
			/* tc_f2で、
			   　 YYYYMMDDHHMMSS_WCN_A1.csv
			      WCN_XXXXXXXXXXXX.csv
			     を作成するようにした。
			     さらに、oki_sbox01/scriptに、csvを収集するためのスクリプトも配置した。
			     なので、この関数で特に処理を行う必要がなくなった。
			     んが、今後の為に関数は残しておく。
			*/

		}
	}
	t.Stop() // タイマを止める。 <- これがないとメモリリークを起こす
}

/*
   300msec毎に、RSU03のCSVファイルをACにコピーする
*/
func cploop_rsu03_csv() {
	t := time.NewTicker(300 * time.Millisecond) // 300msecおきに通知
	for {
		select {
		case <-t.C:
			// 300msec経過した。

			// CSVファイルをACにコピー
			/* tc_f2で、
			   　 YYYYMMDDHHMMSS_WCN_A1.csv
			      WCN_XXXXXXXXXXXX.csv
			     を作成するようにした。
			     さらに、oki_sbox01/scriptに、csvを収集するためのスクリプトも配置した。
			     なので、この関数で特に処理を行う必要がなくなった。
			     んが、今後の為に関数は残しておく。
			*/

		}
	}
	t.Stop() // タイマを止める。 <- これがないとメモリリークを起こす
}


// 300msec毎に、SBOX01~04の通過履歴とWCNテーブルを取り込む
func cploop_sbox_csv() {
	t := time.NewTicker(300 * time.Millisecond) // 300msecおきに通知
	defer t.Stop() // タイマを止める。 <- これがないとメモリリークを起こす

	for {
		select {
		case <-t.C:
			// 300msec経過した。

			// sbox01～sbox04から、通過履歴とWCN番号テーブルを取り込む
			err := exec.Command("./script/update_ac_data.sh").Run()
			if err != nil {
				fmt.Printf("SBOX01~04 WCN_rireki.csv, WCN_table.csv Update Error.\n")
			}
		}
	}
}

/*
   Config.request_interval毎に任意の処理を行う
*/
func timer_send_request() {
	t := time.NewTicker(time.Duration(iniread.Config.Request_interval) * time.Millisecond) // 300msecおきに通知
    defer t.Stop() // タイマを止める。 <- これがないとメモリリークを起こす

    var etc_req_name string
	for {
		select {
		case <-t.C:

            // sbox01
            // 要求中か？（要求中ファイルが存在するか？）
            // ファイルが既に存在している場合はスルー。
            // ファイルが存在していない(要求送信中ではない)場合は、要求を依頼し、要求中ファイルを新規作成
            etc_req_name = iniread.Config.Sbox01_path + "/IcReq"
            _ , err := os.Stat(etc_req_name)
			if err != nil {

                // ETC情報を取得するための要求を発行
                fmt.Printf("SBOX01_Reqest_ASK.\n")
                make_reqfile("RSU", 1, "IA")  // ASK切替要求
                make_reqfile("SBOX", 1, "Ic") // ETC情報取得

				// ファイルが無い場合は新規作成
				_ , err = os.OpenFile(etc_req_name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
                if err != nil {
                    fmt.Printf("etc_req_name:%s create error.\n",etc_req_name)
                }
			}

            // sbox02
            // 要求中か？（要求中ファイルが存在するか？）
            // ファイルが既に存在している場合はスルー。
            // ファイルが存在していない(要求送信中ではない)場合は、要求を依頼し、要求中ファイルを新規作成
            etc_req_name = iniread.Config.Sbox02_path + "/IcReq"
            _ , err = os.Stat(etc_req_name)
			if err != nil {

                // ETC情報を取得するための要求を発行
                fmt.Printf("SBOX02_Reqest_ASK.\n")
                make_reqfile("RSU", 2, "IA")  // ASK切替要求
                make_reqfile("SBOX", 2, "Ic") // ETC情報取得

				// ファイルが無い場合は新規作成
				_ , err = os.OpenFile(etc_req_name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
                if err != nil {
                    fmt.Printf("etc_req_name:%s create error.\n",etc_req_name)
                }
			}


            // sbox03
            // 要求中か？（要求中ファイルが存在するか？）
            // ファイルが既に存在している場合はスルー。
            // ファイルが存在していない(要求送信中ではない)場合は、要求を依頼し、要求中ファイルを新規作成
            etc_req_name = iniread.Config.Sbox03_path + "/IcReq"
            _ , err = os.Stat(etc_req_name)
			if err != nil {

                // ETC情報を取得するための要求を発行
                fmt.Printf("SBOX03_Reqest_ASK.\n")
                make_reqfile("RSU", 3, "IA")  // ASK切替要求
                make_reqfile("SBOX", 3, "Ic") // ETC情報取得

				// ファイルが無い場合は新規作成
				_ , err = os.OpenFile(etc_req_name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
                if err != nil {
                    fmt.Printf("etc_req_name:%s create error.\n",etc_req_name)
                }
			}
		}
	}
}

/*
   20秒ごとに任意の処理を行う
   → 20秒毎に死活監視を発信する
*/
func timer_20sec() {

	t := time.NewTicker(20 * time.Second) // 20秒おきに通知
	defer t.Stop() // タイマを止める。 <- これがないとメモリリークを起こす
	for {
		select {
		case <-t.C:
			// 20秒経過した。

			fmt.Printf("keep alive rsu.\n")
			make_reqfile("RSU", 1, "IH") // RSU死活監視
			make_reqfile("RSU", 2, "IH") // RSU死活監視
			make_reqfile("RSU", 3, "IH") // RSU死活監視
			make_reqfile("RSU", 4, "IH") // RSU死活監視

			fmt.Printf("keep alive sbox.\n")
			make_reqfile("SBOX", 1, "IH") // SBOX死活監視
			make_reqfile("SBOX", 2, "IH") // SBOX死活監視
			make_reqfile("SBOX", 3, "IH") // SBOX死活監視
			make_reqfile("SBOX", 4, "IH") // SBOX死活監視

            fmt.Printf("== ETC System Start!! ==\n")
            make_reqfile("RSU", 1, "IA") // RSUASK切替
            make_reqfile("RSU", 2, "IA") // RSUASK切替
            make_reqfile("RSU", 3, "IA") // RSUASK切替
            make_reqfile("RSU", 4, "IA") // RSUASK切替

            // SBOX無線制御開始は、画面操作からコントロールするためコメントアウト
            // make_reqfile("SBOX", 1, "IRONYES") // SBOX無線制御開始（車載器指示あり）
            // make_reqfile("SBOX", 2, "IRONYES") // SBOX無線制御開始（車載器指示あり）
            // make_reqfile("SBOX", 3, "IRONYES") // SBOX無線制御開始（車載器指示あり）
            // make_reqfile("SBOX", 4, "IRONYES") // SBOX無線制御開始（車載器指示あり）
		}
	}
}

/*
   初期化
*/
func init() {

	iniread.Run() // config.ini読込

    // ログファイル保存設定
    log_setup()
	go timer_10()           // 10秒タイマースタート(ログファイルのローテーション用)

    // データ受信などプログラムが利用するディレクトリを無ければ作成しておく
    work_dir_setup()

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
            fmt.Printf("os.Mkdir -> ./log : %v\n", err)
			return err
		}
	}

	_, err = os.Open(log_run_path)
	if os.IsNotExist(err) {

		// runフォルダ作成
		err = os.Mkdir(log_run_path, 0777)
		if err != nil {
            fmt.Printf("os.Mkdir -> %s : %v\n", log_run_path, err)
			return err
		}
	}

	_, err = os.Open(log_csv_path)
	if os.IsNotExist(err) {

		// csvフォルダ作成
		err = os.Mkdir(log_csv_path, 0777)
		if err != nil {
            fmt.Printf("os.Mkdir -> %s : %v\n", log_csv_path ,err)
			return err
		}
	}

	_, err = os.Open(log_bin_path)
	if os.IsNotExist(err) {

		// binフォルダ作成
		err = os.Mkdir(log_bin_path, 0777)
		if err != nil {
            fmt.Printf("os.Mkdir -> %s : %v\n",log_bin_path, err)
            return err
		}
	}

	return nil
}

/* ログファイル保存設定 */
func log_setup() {

	// Log/CSV用フォルダの指定
	log_bin_path = iniread.Config.Bin_log_path
	log_csv_path = iniread.Config.Csv_log_path
	log_run_path = iniread.Config.Run_log_path

    fmt.Printf("log_bin_path : %s\n",log_bin_path)
    fmt.Printf("log_csv_path : %s\n",log_csv_path)
    fmt.Printf("log_run_path : %s\n",log_run_path)

    // log保存ディレクトリの準備
	err := make_log_folder(log_bin_path, log_csv_path, log_run_path)
	if err != nil {
		panic(err)
	}

	// Log保存ファイル設定
	now := time.Now()
	year_val, month_val, day_val := now.Date() // 年月日を数字で取得してみる
	log_filename := fmt.Sprintf(log_run_path + "/" + "%04d%02d%02d.log", year_val, int(month_val), day_val)

    // fmt.Printf("Run_log_path : %s\n",iniread.Config.Run_log_path)
    // fmt.Printf("log_filename : %s\n",log_filename)

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

// workディレクトリ準備
func work_dir_setup() {

    // 作成するディレクトリのスライス
    paths := []string{
		iniread.Config.Rsu01_path,
		iniread.Config.Rsu02_path,
		iniread.Config.Rsu03_path,
		iniread.Config.Rsu04_path,
		iniread.Config.Sbox01_path,
		iniread.Config.Sbox02_path,
		iniread.Config.Sbox03_path,
		iniread.Config.Sbox04_path,
	}

    // 作成する全てのディレクトリについて、作成されていなければ作成する。
    for _, path := range paths {
		_, err := os.Open(path)
		if os.IsNotExist(err) {

			// 受信データ保存様フォルダ作成
			err = os.Mkdir(path, 0777)
			if err != nil {
				fmt.Printf("os.Mkdir() Error for path %s : %v\n", path, err)
			}
		}
	}
}

/*
   10秒毎に任意の処理を行う
*/
func timer_10() {
	t := time.NewTicker(10 * time.Second) // 10秒おきに通知
	defer t.Stop() // タイマを止める。 <- これがないとメモリリークを起こす

	for {
		select {
		case <-t.C:
			// 10秒経過した。
			now := time.Now()
			year_val, month_val, day_val := now.Date() // 年月日を数字で取得してみる
			log_filename := fmt.Sprintf(log_run_path+ "/" + "%04d%02d%02d.log", year_val, int(month_val), day_val)

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

func main() {

	go timer_20sec()        // AP → RSU/SBOX 死活監視要求送信(20秒間隔)

    //	go timer_send_request() // AP → RSU/SBOX 要求テキストファイルを作成   // 2022/08/25 こちらからの要求シーケンスが変わるため一旦コメントアウト
	//    go cploop_rsu01_csv()    // RSU01のCSVファイルをAC側にコピー
	//    go cploop_rsu02_csv()    // RSU02のCSVファイルをAC側にコピー
	//    go cploop_rsu03_csv()    // RSU03のCSVファイルをAC側にコピー
	//    go cploop_rsu04_csv()    // RSU01のCSVファイルをAC側にコピー    <-- 将来拡張するならこんな感じで

	go cploop_sbox_csv() //300msec毎に、SBOX01~04の通過履歴とWCNテーブルを取り込む


    // システム起動時にASK切替と無線制御開始要求を送信する
    // 以降、SBOXからSc通知が不定期で送信されてくる。
    fmt.Printf("== ETC System Start!! ==\n")
    make_reqfile("RSU", 1, "IA") // RSUASK切替
    make_reqfile("RSU", 2, "IA") // RSUASK切替
    make_reqfile("RSU", 3, "IA") // RSUASK切替
    make_reqfile("RSU", 4, "IA") // RSUASK切替

    //// 2023/11/20 システム起動時に電波発信コマンドを送らない。
    //// 電波発信は、画面端末アプリの発射ボタンで行う。
    //// どうしてもシステムで発信したい場合は、sbox01〜03のカレントディレクトリで、要求ファイルSIRONNOを作成する
    // make_reqfile("SBOX", 1, "IRONNO") // SBOX無線制御開始（車載器指示なし）
    // make_reqfile("SBOX", 2, "IRONNO") // SBOX無線制御開始（車載器指示なし）
    // make_reqfile("SBOX", 3, "IRONNO") // SBOX無線制御開始（車載器指示なし）
    // make_reqfile("SBOX", 4, "IRONNO") // SBOX無線制御開始（車載器指示なし）


	// メインスレッドを終わらせない
    select{}

}
