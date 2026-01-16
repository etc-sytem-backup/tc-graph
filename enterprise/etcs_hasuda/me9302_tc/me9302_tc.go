/*
   ログに対応させる
   config.iniに対応させる
   OKIのtcと同様のカラム数でcsvファイルを吐き出せるようにする。
*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
    "os/exec"
	"strconv"
	"strings"
	"time"

	"etc-system.jp/iniread"
)

/* Package Construct */
const (
    UDP_READ_BUF = 64
)


/* Package Global var  */
var (
    log_bin_path string = "./"        // バイナリファイル格納用パス
    log_run_path string = "./"        // 動作ログファイル格納用パス
    log_csv_path string = "./"        // CSVファイル格納用パス
    ipAddress string    = ""             // IPアドレス:ポート番号
    MachineNo int       = 0                 // RSUの機器番号
)


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

// workディレクトリ準備
func work_dir_setup() {

    // 作成するディレクトリのスライス
    paths := []string{
		iniread.Config.A01_tc_csv_path,
		iniread.Config.A02_tc_csv_path,
		iniread.Config.A03_tc_csv_path,
		iniread.Config.A04_tc_csv_path,
		iniread.Config.A01_tc_wcn_path,
		iniread.Config.A02_tc_wcn_path,
		iniread.Config.A03_tc_wcn_path,
		iniread.Config.A04_tc_wcn_path,
        iniread.Config.A01_tc_table_path,
        iniread.Config.A02_tc_table_path,
        iniread.Config.A03_tc_table_path,
        iniread.Config.A04_tc_table_path,
        iniread.Config.A01_wcn_table_path,
        iniread.Config.A02_wcn_table_path,
        iniread.Config.A03_wcn_table_path,
        iniread.Config.A04_wcn_table_path,
	}

    // 作成する全てのディレクトリについて、作成されていなければ作成する。
    for _, path := range paths {
		_, err := os.Open(path)
		if os.IsNotExist(err) {

			// 受信データ保存フォルダ作成
			err = os.Mkdir(path, 0777)
			if err != nil {
				fmt.Printf("os.Mkdir() Error for path %s : %v\n", path, err)
			} else {
                fmt.Printf("Directory a [%s]\n",path)
            }
            
		} else {
            fmt.Printf("FileFind...There is a [%s]\n",path)
        }
	}
}

/* ログファイル保存設定 */
func log_setup() {

	// Log/CSV用フォルダの指定
	log_bin_path = iniread.Config.Bin_log_path
	log_csv_path = iniread.Config.Csv_log_path
	log_run_path = iniread.Config.Run_log_path

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

/* Package Initial  */
func init() {
   
    // Config.ini読み込み
    iniread.Run()

    // ログファイル保存設定
    log_setup()
	go timer_10()           // 10秒タイマースタート(ログファイルのローテーション用)

    // データ受信などプログラムが利用するディレクトリを無ければ作成しておく
    work_dir_setup()

	// コマンドオプション初期設定
	flag.StringVar(&ipAddress, "i", "192.168.1.202:11000", "IP Address")
	flag.IntVar(&MachineNo, "n", 1, "Machine Number")

	// コマンドオプション解析
	flag.Parse()

    // 与えられたオプションが2個以上または2個以下の場合は、使い方を表示してプログラム終了。
	var prm_cnt int = 3
	if (len(os.Args) < prm_cnt) || (len(os.Args) > prm_cnt) {
		fmt.Println("Usage: me0X -i=\"192.168.110.212:58001\" -n=1")
		fmt.Println("-i : IP Address:PortNo")
		fmt.Println("-n : Machine Number")
		fmt.Printf("ARGS: %d\n", len(os.Args))
		os.Exit(-1)
	}

	// 取得パラメータを標準出力に（確認）
	log.Printf("param -i : %s\n", ipAddress)
	log.Printf("param -n : %d\n", MachineNo)

    log.Printf("-- Package Inisial OK. --\n")

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

/* センサー通過情報保存（CSVファイル作成）
   SBOXからの応答受診時に作成されるデータを直接作成する。
   
   save_path     : データ保存先ファイルパス
   timestamp    : YYYYMMDDhhmmss
   rsu_name     : RSU01 ～ RSU04
   rsu_name2    : A1 ～ A2
   write_status : IN or PARK or OUT
   wcn_num      : XXXXXXXXXXXX <- 12 chars
   card_num     : XXXXXXXXXXXXXXXXXXX <- 19 chars

   <Exsample>
    日時(マイクロ秒含む),装置コード,アンテナID,WCN番号     ,ステータス,ETCカード番号       ,支局  ,用途,種別,一連番号」とする。
    20230710000001000   ,M0001     ,A1        ,016072700261,IN        ,01198047052906908bbb,465341,cb  ,100 ,4226

 */
func make_passage_data(save_path string,
                       timestamp string,
                       rsu_name string,
                       rsu_name2 string,
                       write_status string,
                       wcn_num string,
                       card_num string,
                       sikyoku string,
                       youto string,
                       syubetsu string,
                       ichiren string,
) {

    // YYYYMMDDHHMMSS_WCN_A1.csvファイル作成
    var filename string

    // CSVファイル名作成
    filename = save_path + "/" + timestamp + "_WCN_" + rsu_name2 + ".csv"

    // ファイルがなければ作成し、すでにあればスルーする
    // panic: open .: too many open files の対策として無名関数化する
    func() {
        fp, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
        if err != nil {
            //エラー処理
            log.Printf("OpenFile Error %s : %v",filename, err)
            log.Fatal(err)
        }
        defer fp.Close()

        // CSVファイルにデータ書き込み（逆走検知画面で活用予定。ETC情報要求応答受信時に保存される。）
        w := bufio.NewWriter(fp)
        writeline := timestamp + "," + rsu_name + "," + rsu_name2 + "," + wcn_num + "," + write_status + "," + card_num + "," + sikyoku + "," + youto + "," + syubetsu + "," + ichiren + "\n"
        fmt.Fprint(w, writeline)      // 作成データをファイルに書き込み
        w.Flush()                     // ファイル書き込みを確定させる
        fmt.Printf("File Save -> %s\n",filename)
    }()
}

/* 車両情報（WCN番号）保存（csvファイル作成）
   A1～A4を通過した車両のWCN番号を保存する。
   一度保存したWCN番号は保存しない

   wcn_file_path    : ファイル保存パス
   time_stamp       : 時刻
   wcn_number       : WCN番号
   etc_card_number  : カード番号
 */
func make_wcncar_data(wcn_file_path string, time_stamp string, wcn_number string, etc_card_number string) {

    // 受信したWCNをファイル名とする。
    filename := wcn_file_path + "/" + "WCN_" + wcn_number + ".csv"

    // ファイルがなければ作成し、すでにあればスルーする
    // panic: open .: too many open files の対策として無名関数化する
    func() {
        fp, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
        if err != nil {
            //エラー処理
            log.Fatal(err)
        }
        defer fp.Close()

        // CSVファイルにデータ書き込み(予約画面で利用予定。WCN番号とETCカードIDを紐付けている)
        w := bufio.NewWriter(fp)
        fmt.Fprint(w, wcn_number + "," + etc_card_number + "\n")      // データ WCN番号,ETCカードID
        w.Flush()
        fmt.Printf("File Save -> %s\n",filename)
    }()
}


/* 与えられた年月日時分秒(ms含む)２つの差分（秒）を求める */
func date_duration(time1 string, time2 string) (int, error) {

    //                                       01234567890123456
    // time1, time2の差分（秒）を取得する -> 20221014142509891
    t1_year_s  := time1[0:4]
    t1_month_s := time1[4:6]
    t1_day_s   := time1[6:8]
    t1_hh_s    := time1[8:10]
    t1_mm_s    := time1[10:12]
    t1_ss_s    := time1[12:14]

    t2_year_s  := time2[0:4]
    t2_month_s := time2[4:6]
    t2_day_s   := time2[6:8]
    t2_hh_s    := time2[8:10]
    t2_mm_s    := time2[10:12]
    t2_ss_s    := time2[12:14]

    t1_year,err := strconv.Atoi(t1_year_s)
    if err != nil {
        log.Printf("t1_year strconv.Atoi Error!!\n")
        log.Println(err)
    }

    t1_month,err := strconv.Atoi(t1_month_s)
    if err != nil {
        log.Printf("t1_month strconv.Atoi Error!!\n")
        log.Println(err)
    }

    t1_day,err := strconv.Atoi(t1_day_s)
    if err != nil {
        log.Printf("t1_day strconv.Atoi Error!!\n")
        log.Println(err)
    }

    t1_hh,err := strconv.Atoi(t1_hh_s)
    if err != nil {
        log.Printf("t1_hh strconv.Atoi Error!!\n")
        log.Println(err)
    }

    t1_mm,err := strconv.Atoi(t1_mm_s)
    if err != nil {
        log.Printf("t1_mm strconv.Atoi Error!!\n")
        log.Println(err)
    }

    t1_ss,err := strconv.Atoi(t1_ss_s)
    if err != nil {
        log.Printf("t1_ss strconv.Atoi Error!!\n")
        log.Println(err)
    }

    t2_year,err := strconv.Atoi(t2_year_s)
    if err != nil {
        log.Printf("t2_year strconv.Atoi Error!!\n")
        log.Println(err)
    }

    t2_month,err := strconv.Atoi(t2_month_s)
    if err != nil {
        log.Printf("t2_month strconv.Atoi Error!!\n")
        log.Println(err)
    }

    t2_day,err := strconv.Atoi(t2_day_s)
    if err != nil {
        log.Printf("t2_day strconv.Atoi Error!!\n")
        log.Println(err)
    }

    t2_hh,err := strconv.Atoi(t2_hh_s)
    if err != nil {
        log.Printf("t2_hh strconv.Atoi Error!!\n")
        log.Println(err)
    }

    t2_mm,err := strconv.Atoi(t2_mm_s)
    if err != nil {
        log.Printf("t2_mm strconv.Atoi Error!!\n")
        log.Println(err)
    }

    t2_ss,err := strconv.Atoi(t2_ss_s)
    if err != nil {
        log.Printf("t2_ss strconv.Atoi Error!!\n")
        log.Println(err)
    }

    before := time.Date(t1_year, time.Month(t1_month), t1_day, t1_hh, t1_mm, t1_ss, 0, time.Local)
    after := time.Date(t2_year, time.Month(t2_month), t2_day, t2_hh, t2_mm, t2_ss, 0, time.Local)
    duration := after.Sub(before)
    //log.Printf("Duration: %d\n",duration)
    
    //    hours0 := int(duration.Hours())
    //    days := hours0 / 24
    //    hours := hours0 % 24
    mins := int(duration.Minutes()) % 60
    secs := int(duration.Seconds()) % 60
    log.Printf("Duration_time : %d min, %d sec\n",mins,secs)

    // 1分以上時間が空いていた場合は、経過した分だけ戻り値の秒に加算する。
    if mins > 0 {
        add := mins * 60        // 分→秒 変換
        secs = secs + add       // 変換した秒数を加算
    }

    return secs,err

    /* Sample
    day1 := time.Date(2000, 12, 31, 0, 0, 0, 0, time.Local)
    day2 := time.Date(2001, 1, 2, 12, 30, 0, 0, time.Local)
    duration := day2.Sub(day1)
    fmt.Println(duration) // => "60h30m0s"

    hours0 := int(duration.Hours())
    days := hours0 / 24
    hours := hours0 % 24
    mins := int(duration.Minutes()) % 60
    secs := int(duration.Seconds()) % 60
    fmt.Printf("%d days + %d hours + %d minutes + %d seconds\n", days, hours, mins, secs)
    // => "2 days + 12 hours + 30 minutes + 0 seconds"
    */

}


// 受信したWCN番号を、指定保存先に保存する。
// ■保存フォーマット
//    日時(マイクロ秒含む),装置コード,アンテナID,WCN番号     ,ステータス,ETCカード番号       ,支局  ,用途,種別,一連番号」とする。
// → 20230710000001000   ,M0001     ,A1        ,016072700261,IN        ,01198047052906908bbb,465341,cb  ,100 ,4226
func wcn_save(wcn_num string,       // WCN番号
              a_num int) {          // アンテナ番号(1 〜 4)

    var (
        antenna_name string         // 装置コード
        antenna_alias string        // アンテナID
        antenna_field string        // ステータス
        tc_save_path string         // UDP受信データ保存path
        wcn_save_path string        // WCN番号ファイル保存path
    )

    // 車両情報にダミーデータをセット
    card_num := "********************" // ETCカードナンバー
    sikyoku := "******"                // 支局    
    youto := "**"                      // 用途    
    syubetsu := "***"                  // 種別    
    ichiren := "****"                  // 一連番号

    switch (a_num) {
    case 1: 
        tc_save_path = iniread.Config.A01_tc_csv_path
        wcn_save_path = iniread.Config.A01_tc_wcn_path
        antenna_name = "M0001"
        antenna_alias = "A1"
        antenna_field = "IN"
    case 2:
        tc_save_path = iniread.Config.A02_tc_csv_path
        wcn_save_path = iniread.Config.A02_tc_wcn_path
        antenna_name = "M0002"
        antenna_alias = "A2"
        antenna_field = "PARK"
    case 3:
        tc_save_path = iniread.Config.A03_tc_csv_path
        wcn_save_path = iniread.Config.A03_tc_wcn_path
        antenna_name = "M0003"
        antenna_alias = "A3"
        antenna_field = "OUT"
    case 4:
        tc_save_path = iniread.Config.A04_tc_csv_path
        wcn_save_path = iniread.Config.A04_tc_wcn_path
        antenna_name = "M0004"
        antenna_alias = "A4"
        antenna_field = "OTHER"
    default:
    }

    // ミリ秒を含む日付文字列を作成(YYYYMMDDhhmmssxxx)
    time_stamp := get_datestr()

    /* 受信したWCN番号がすでにシステムに検知されている場合、前回受信した時間から何秒経過しているかチェックする。
       指定秒数(INIファイル)以内に受信した場合、渋滞とみなしてcsvファイルは作成しない。
     */
    log.Printf("wcn_num:%s, a_num:%d\n",wcn_num, a_num)
    result, err := exec.Command("./script/get_receive_time.sh", wcn_num, strconv.Itoa(a_num)).Output()
    if err != nil {
        log.Printf("get_receive_time.sh Error!!: %v\n",err)
        log.Printf("WCN_Number Find Result : %s\n",string(result))
    }

    // WCN番号検出結果判定
    res := string(result)
    res = strings.TrimRight(res,"\n")
    log.Printf("./script/get_receive_time.sh -> %s\n",res)
    fmt.Printf("./script/get_receive_time.sh -> %s\n",res)
    if res == "NoHit" {    // WCN_rireki.csvファイルはあるけれど、受信したWCNは登録されていない
    
        //log.Printf("Receive WCN First Save : %s <- %s\n",res, wcn_num) // WCN番号の受信時間を保存する。
        log.Printf("Receive WCN First Save : %s <- %s\n", time_stamp, wcn_num) // WCN番号の受信時間を保存する。
        fmt.Printf("Receive WCN First Save : %s <- %s\n", time_stamp, wcn_num) // WCN番号の受信時間を保存する。

        // CSVデータ作成
        make_passage_data(tc_save_path,time_stamp,antenna_name,antenna_alias,antenna_field,wcn_num,card_num,sikyoku,youto,syubetsu,ichiren)

        // wcnデータ作成
        make_wcncar_data(wcn_save_path,time_stamp,wcn_num,card_num)
    }

    // WCN_rireki.csvファイルが存在しており、受信したWCN番号がすでに登録されている。アンテナ番号も間違っていない。
    if res != "NotFound" && res != "NoMachine" && res != "NoHit"{
        a_duration, err := date_duration(res,time_stamp) // 前回の受信時間と今回の受信時間を比較し、時間間隔(秒)を取得
        if err != nil {
            log.Printf("Error...result:%s, time_stamp:%s\n",res ,time_stamp)
            log.Println(err)
        }
        
        // 時間間隔が設定時間をオーバーしていたらCSVファイルを作成する（渋滞ではないと考える。acによってWCN_rireki_A?.csvに追加される）
        log.Printf("Receive Duration : %d Sec = %s - %s\n",a_duration, time_stamp, res)
        if a_duration >= iniread.Config.Detection_interval {

            log.Printf("Receive WCN Save Duration %d Sec : %s <- %s\n",a_duration, time_stamp, wcn_num)
            fmt.Printf("Receive WCN Save Duration %d Sec : %s <- %s\n",a_duration, time_stamp, wcn_num)

            // CSVデータ作成
            make_passage_data(tc_save_path,time_stamp,antenna_name,antenna_alias,antenna_field,wcn_num,card_num,sikyoku,youto,syubetsu,ichiren) // WCN番号の受信時間を保存する。

            // wcnデータ作成
            make_wcncar_data(wcn_save_path,time_stamp,wcn_num,card_num)
        }
    }
}


// WCNファイル及び通過履歴ファイルを一本化する。
// 直下の./tc_csvに受信ファイルが作成されている。
// 一本にまとめた物を配置するのは、../sbox01〜03/tc_csv_table/とする。
// ※OKIバージョンのacをそのまま利用するため。
func make_csv_table() {
    t := time.NewTicker(1 * time.Second) // 1秒おき

    defer t.Stop() // タイマを止める。 <- これがないとメモリリークを起こす

    for {
        select {
        case <-t.C:            // 1秒経過した。

            // WCNファイルの一本化
            _, err := exec.Command("./script/make_wcn_table.sh", strconv.Itoa(MachineNo)).Output()
            if err != nil {
                log.Printf("make_wcn_table.sh Error! : %v\n",err)
            }

            // 通過履歴ファイルの一本化
            _, err = exec.Command("./script/make_csv_table.sh", strconv.Itoa(MachineNo)).Output()
            if err != nil {
                log.Printf("make_csv_table.sh Error! : %v\n",err)
            }
        }
    }
}


/* ETC-System版 トラフィックカウンター  */
func main() {

    log.Printf("-- ETC-System Traffic Counter Start --")
    fmt.Printf("-- ETC-System Traffic Counter Start --")

    recAddr, err := net.ListenPacket("udp", ipAddress) // データ受信IPアドレスとポート(UDP)
    if err != nil {
        log.Println("ListenPacket Error!!")
        fmt.Println("ListenPacket Error!!")
        os.Exit(-1)
    }
    defer recAddr.Close()
    log.Println("-- UDP Listen OK --")
    fmt.Println("-- UDP Listen OK --")
    log.Printf("-- UDP Listen Ip/Port : %s --\n",ipAddress)
    fmt.Printf("-- UDP Listen Ip/Port : %s --\n",ipAddress)

    go make_csv_table()     // 通過履歴、WCN番号の一覧化（１本化）

    var buf [UDP_READ_BUF]byte            // 受信バッファ（固定）
    for {
        n, addr, err := recAddr.ReadFrom(buf[:])
        if err != nil {
            log.Println("ReadFrom Error!!")
            break
        }

        // 受信データから余分な改行コード「CRLF」「LF」を取り除く。正しいと思われるWCN番号以外は処理せず(continue)再び受信する。
        var wcn_num string
        wcn_num = string(buf[:n])
        wcn_num = strings.TrimRight(wcn_num,"\n")
        wcn_num = strings.TrimRight(wcn_num,"\r")

        // 受信データにマイナスが含まれている場合、それはME9302がアンテナとの接続に成功している通知らしい（死活監視）
        if strings.Index(wcn_num,"-") == -1 {
            log.Printf("Received from %v, Data : %s", addr, string(buf[:n]))            
            fmt.Printf("Received from %v, Data : %s", addr, string(buf[:n]))            
        } else {
            continue
        }

        // UDPで受信したWCN番号をファイルに保存する。
        // ■フォーマット
        //    日時(マイクロ秒含む),装置コード,アンテナID,WCN番号     ,ステータス,ETCカード番号       ,支局  ,用途,種別,一連番号」とする。
        // → 20230710000001000   ,M0001     ,A1        ,016072700261,IN        ,01198047052906908bbb,465341,cb  ,100 ,4226
        switch {
        case n >= 1: // WCN番号を検知した

            // 検知したwcn番号をファイルに保存する
            wcn_save(wcn_num,MachineNo)
        
        case n == -1: // WCN番号ではない何かを検知
            // 特に何も処理しない。

        }
    }
}

