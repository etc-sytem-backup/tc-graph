/*
   トラフィックカウンター

   ■2024/03/31 車種判定について
   現車種判定の処理は、ナンバープレート情報から導き出している。
   そうではなく、車載器固有情報に含まれる車種情報から、車種を導き出したい場合についてメモを残す。

   SBOXからの受信データに対し、車載器固有情報から車種を導き出すコードのサンプルが入っている。
   　→　./syasyuhantei/convert_number_plate_data.zip

   1005行目付近にて、SBOXからのバイナリデータを切り出している。
   そこに、切り出す処理を追加する。
   具体的には、保存される原始データcsvの15列目に対し、先頭から15文字目から3文字がそれに当たる。
   詳細は、GitHubリポジトリ「2401_Documents」の下記ファイルを参照。
   2401_Documents/10_作業用資料/10_開発設計メモ_アプリ毎_drawio/data_analysis_genshi.drawio
*/

package main

import (
    "bufio"
    "bytes"
    "context"
    "encoding/binary"
    "flag"
    "fmt"
    "log"
    "os"
    "runtime"
    "strings"
    "time"
    "strconv"
    "os/exec"

    "localhost.com/conretry"
    "localhost.com/findcmd"
    "localhost.com/makecmd"
    "localhost.com/tcpclient"
	"etc-system.jp/iniread"
)

/*
   Traffic Counter Initialize
*/
var Machine string   // SBOX or ME93 or RSU
var MachineNo int    // RSUの機器番号
var IpAddress string // IPアドレス:ポート番号
var MCommand string  // コマンド
var SeqNumber uint16 // 送受信シーケンス番号

var RecLID [8]byte // LID格納用(ASCII)（RSUリンク接続応答）

var log_bin_path string // バイナリファイル格納用パス
var log_run_path string // 動作ログファイル格納用パス
var log_csv_path string // CSVファイル格納用パス
var tc_wcn_path string  // WCNファイル格納用パス(ac_f2連携用)
var tc_csv_path string  // ac_f2で必要な「YYYYMMDDHHMMSS_WCN_A1.csv」作成用パス
var tc_wcn_table_path string  // WCNファイル（一覧）格納用パス(ac_f2連携用)
var tc_csv_table_path string  // ac_f2で必要な「YYYYMMDDHHMMSS_WCN_A1.csv」の一覧ファイル作成用パス

/*
   traffic_counterイニシャル
*/
func init() {

    // Config.ini読み込み
    iniread.Run()

    // コマンドオプション初期設定
    flag.StringVar(&Machine, "s", "SBOX", "SBOX or ME93 or RSU or LOG")
    flag.IntVar(&MachineNo, "n", 1, "Machine Number")
    flag.StringVar(&IpAddress, "i", "192.168.1.201:50000", "IP Address")
    flag.StringVar(&MCommand, "c", "IH", "Machine Command")

    // コマンドオプション解析
    flag.Parse()

    // パラメータ少ない場合はメッセージ表示して終了
    // LOGシステムの場合は-cオプション無くてもオッケー
    var prm_cnt int
    if Machine == "LOG" {
        prm_cnt = 3
    } else {
        prm_cnt = 4
    }

    if len(os.Args) < prm_cnt {
        fmt.Println("Usage: tc_f2 -s=\"RSU\" -i=\"192.168.110.11:50001\" -c=\"IH\" -n=1")
        fmt.Println("-s : SBOX or RSU or LOG")
        fmt.Println("-i : IP Address:PortNo")
        fmt.Println("-c : Command")
        fmt.Println("-n : Machine Number")
        fmt.Printf("ARGS: %d\n", len(os.Args))
        os.Exit(-1)
    }

    // 取得パラメータを標準出力に（確認）
    fmt.Printf("param -s : %s\n", Machine)
    fmt.Printf("param -i : %s\n", IpAddress)
    fmt.Printf("param -c : %s\n", MCommand)
    fmt.Printf("param -n : %d\n", MachineNo)

    // シーケンスナンバー初期化
    SeqNumber = 0

    // Log/CSV用フォルダの作成
    log_bin_path = "./log/bin/"
    log_csv_path = "./log/csv/"
    log_run_path = "./log/run/"
    tc_wcn_path = "./tc_wcn/"
    tc_csv_path = "./tc_csv/"
    tc_wcn_table_path = "./tc_wcn_table/"
    tc_csv_table_path = "./tc_csv_table/"
    err := make_log_folder()
    if err != nil {
        panic(err)
    }

    // Log保存ファイル設定
    now := time.Now()
    year_val, month_val, day_val := now.Date() // 年月日を数字で取得してみる
    log_filename := fmt.Sprintf(log_run_path+"%04d%02d%02d.log", year_val, int(month_val), day_val)

    // ファイルが既に存在している場合はスルー。
    // ファイルが存在していない(1分過ぎている)場合は、新規作成してログデータの保存先にする
    log_file, err := os.OpenFile(log_filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {

        // ファイルが無い場合は新規作成
        log_file, _ = os.OpenFile(log_filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    }

    // Logの保存先をファイルにする(デフォルトは標準出力)
    log.SetOutput(log_file)

    // 実行環境情報をログに残す
    log.Printf("NumCPU: %d\n", runtime.NumCPU())
    log.Printf("NumGoroutine: %d\n", runtime.NumGoroutine())
    log.Printf("Version: %s\n", runtime.Version())

}

/*
   受信バイナリ保存用。CSVファイル保存用。動作ログ保存用。
*/
func make_log_folder() error {

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

        // binフォルダ作成
        err = os.Mkdir(log_bin_path, 0777)
        if err != nil {
            return err
        }
    }

    /* tc_wcn_path, tc_wcn_path、tc_wcn_table_path、tc_csv_table_path はac_f2連携ファイル格納用  */
    _, err = os.Open(tc_wcn_path)
    if os.IsNotExist(err) {

        // csvフォルダ作成
        err = os.Mkdir(tc_wcn_path, 0777)
        if err != nil {
            return err
        }
    }

    _, err = os.Open(tc_csv_path)
    if os.IsNotExist(err) {

        // csvフォルダ作成
        err = os.Mkdir(tc_csv_path, 0777)
        if err != nil {
            return err
        }
    }

    _, err = os.Open(tc_csv_table_path)
    if os.IsNotExist(err) {

        // csvフォルダ作成
        err = os.Mkdir(tc_csv_table_path, 0777)
        if err != nil {
            return err
        }
    }

    _, err = os.Open(tc_wcn_table_path)
    if os.IsNotExist(err) {

        // csvフォルダ作成
        err = os.Mkdir(tc_wcn_table_path, 0777)
        if err != nil {
            return err
        }
    }


    return nil
}

/*
   1バイトバイナリを上位4ビットと下位4ビットに分けて、それぞれをアスキーコード化
   0x12　→　0x31,0x32
   0xCD　→　0x43,0x44
*/
func make_ascii(bin uint8) (byte, byte) {

    val_high := bin >> 4
    val_low := bin & 0x0F

    if 0 <= val_high && val_high <= 9 {
        val_high = val_high + 0x30
    } else {

        switch val_high {
        case 0x0A:
            val_high = 0x41
        case 0x0B:
            val_high = 0x42
        case 0x0C:
            val_high = 0x43
        case 0x0D:
            val_high = 0x44
        case 0x0E:
            val_high = 0x45
        case 0x0F:
            val_high = 0x46
        }
    }

    if 0 <= val_low && val_low <= 9 {
        val_low = val_low + 0x30
    } else {

        switch val_low {
        case 0x0A:
            val_low = 0x41
        case 0x0B:
            val_low = 0x42
        case 0x0C:
            val_low = 0x43
        case 0x0D:
            val_low = 0x44
        case 0x0E:
            val_low = 0x45
        case 0x0F:
            val_low = 0x46
        }
    }

    return val_high, val_low
}

/*
   コネクションリードし続ける。
   すでにコネクション確立されている前提
*/
func receive_loop(client *tcpclient.Client) (string, error) {

    var req_str string

    // ひたすら受信処理
    for {
        resp, err := client.Read_command()

        // 応答/通知 種別判定
        //        fmt.Printf("GetReceiveLoop  : %X\n",resp)
        //        fmt.Printf("GetReceiveLength: %d\n",len([]byte(resp)))
        if len([]byte(resp)) < 18 {
            continue
        } // 受信したデータが18バイト未満はありえないので、18バイト未満受信した場合は無視して次のループへ

        if err != nil {
            fmt.Println("ReadCommandError")
        } else {

            tmp := []byte(resp)
            switch {
            case tmp[16] == 0x4e && tmp[17] == 0x4a: // NJ
                req_str = "NJ"

            case tmp[16] == 0x4e && tmp[17] == 0x46: // NF
                req_str = "NF"

            case tmp[16] == 0x4e && tmp[17] == 0x49: // NI
                req_str = "NI"

            case tmp[16] == 0x53 && tmp[17] == 0x54: // ST
                req_str = "ST"

            case tmp[16] == 0x56 && tmp[17] == 0x54: // VT
                req_str = "VT"

            case tmp[16] == 0x41 && tmp[17] == 0x48: // AH
                req_str = "AH"

            case tmp[16] == 0x41 && tmp[17] == 0x54: // AT
                req_str = "AT"

            case tmp[16] == 0x41 && tmp[17] == 0x41: // AA
                req_str = "AA"

            case tmp[16] == 0x41 && tmp[17] == 0x51: // AQ
                req_str = "AQ"

            case tmp[16] == 0x41 && tmp[17] == 0x43: // AC
                req_str = "AC"

                // LIDのアスキー化
                val1, val2 := make_ascii(tmp[22])
                RecLID[0] = val1
                RecLID[1] = val2
                val1, val2 = make_ascii(tmp[23])
                RecLID[2] = val1
                RecLID[3] = val2
                val1, val2 = make_ascii(tmp[24])
                RecLID[4] = val1
                RecLID[5] = val2
                val1, val2 = make_ascii(tmp[25])
                RecLID[6] = val1
                RecLID[7] = val2

                log.Printf("LID = %X\n", RecLID[:])
                fpw, err := os.Create("LID.bin") // ファイル作成
                if err != nil {
                    fmt.Println("LID BinFileCreate Error!!")
                    log.Println("LID BinFileCreate Error!!")
                }
                // 固定長部分を書き込み
                err = binary.Write(fpw, binary.BigEndian, &RecLID)
                if err != nil {
                    fmt.Println("LID Write Error!!")
                    log.Println("LID Write Error!!")
                }
                fpw.Close()

            case tmp[16] == 0x53 && tmp[17] == 0x63: // Sc
                req_str = "Sc"

                // // LIDのアスキー化
                // val1, val2 := make_ascii(tmp[22])
                // RecLID[0] = val1
                // RecLID[1] = val2
                // val1, val2 = make_ascii(tmp[23])
                // RecLID[2] = val1
                // RecLID[3] = val2
                // val1, val2 = make_ascii(tmp[24])
                // RecLID[4] = val1
                // RecLID[5] = val2
                // val1, val2 = make_ascii(tmp[25])
                // RecLID[6] = val1
                // RecLID[7] = val2

                RecLID[0] = tmp[34]
                RecLID[1] = tmp[35]
                RecLID[2] = tmp[36]
                RecLID[3] = tmp[37]
                RecLID[4] = tmp[38]
                RecLID[5] = tmp[39]
                RecLID[6] = tmp[40]
                RecLID[7] = tmp[41]

                log.Printf("Sc_LID = %X\n", RecLID[:])
                fpw, err := os.Create("Sc_LID.bin") // ファイル作成
                if err != nil {
                    fmt.Println("Sc_LID BinFileCreate Error!!")
                    log.Println("Sc_LID BinFileCreate Error!!")
                }
                // 固定長部分を書き込み
                err = binary.Write(fpw, binary.BigEndian, &RecLID)
                if err != nil {
                    fmt.Println("Sc_LID Write Error!!")
                    log.Println("Sc_LID Write Error!!")
                }
                fpw.Close()


            case tmp[16] == 0x41 && tmp[17] == 0x49: // AI
                req_str = "AI"

            case tmp[16] == 0x41 && tmp[17] == 0x44: // AD
                req_str = "AD"

            case tmp[16] == 0x41 && tmp[17] == 0x63: // Ac
                req_str = "Ac"

            case tmp[16] == 0x41 && tmp[17] == 0x52: // AR
                req_str = "AR"

            default:
                req_str = "Err"
            }

            // 応答情報を画面表示
            fmt.Printf("===> %s -> AP : %s %X\n", Machine, req_str, []byte(resp))

            //                log.Printf("===> %s -> AP : %s %X\n", Machine, req_str ,[]byte(resp))    // 容量削減のためコメントアウト

            /* 2022/05/22 容量削減のため。binは残さないこととする。特に必要とされてもいないので。
            now := time.Now()
            year_val, month_val, day_val := now.Date() // 年月日を数字で取得してみる
            //    fmt.Printf("%04d年%02d月%02d日 %02d:%02d:%02d\n",year_val,int(month_val),day_val, now.Hour(), now.Minute(), now.Second())
            fp_header := fmt.Sprintf("%04d%02d%02d_%02d%02d%02d", year_val, int(month_val), day_val, now.Hour(), now.Minute(), now.Second())
            //    fp_header :=fmt.Sprintf("%04d%02d%02d",year_val,int(month_val),day_val)
            filename := log_bin_path + fp_header + "_" + Machine + "->AP_" + req_str + "_receive.bin"

            // ファイルがなければ作成し、あれば追記する。
            // panic: open .: too many open files の対策として無名関数化
            func() {
                fp, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
                if err != nil {
                    //エラー処理
                    log.Fatal(err)
                }
                defer fp.Close()
                binary.Write(fp, binary.BigEndian, []byte(resp))
            }()
            */

            // CSVファイル作成
            // コマンドによって、バイナリファイルを切り分け、CSVフォーマットでデータを保存する
            err = make_csv(Machine, req_str, []byte(resp), log_csv_path)

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

/*
   30秒ごとに任意の処理を行う
   → 30秒毎に死活監視を発信する
   
   2022/05/22 現在この関数は呼ばれていない。
   ac_f2側にて死活監視を制御するように変更したため。
   とはいえ、今後のことを考えて消さずに残してはおくようにした。
*/
func timer_30(client *tcpclient.Client) {

    send_byte := new(bytes.Buffer) // 送信バイナリ格納用

    // TCP通信キャンセル指示用のコンテキストを作成する
    ctx := context.Background()

    t := time.NewTicker(30 * time.Second) // 30秒おきに通知
    for {
        select {
        case <-t.C:
            // 30秒経過した。

            send_byte = make_request(client, Machine, "IH", IpAddress, send_byte)

            fmt.Printf("===> AP -> %s : IH %X\n", Machine, send_byte)
            //            log.Printf("===> AP -> %s : IH %X\n", Machine, send_byte)            // 容量削減の為コメントアウト

            // コマンド送信処理（リトライ3回）
            err := conretry.Retry(ctx, 3, 0, func() error {
                var ierr error
                _, ierr = client.Send_command(send_byte.Bytes()) // AP → RSU/SBOX 死活監視要求送信
                return ierr
            })
            if err != nil {
                panic(err)
            }

            /* 2022/05/22 容量削減のため、バイナリファイルは作成しない。現在、特に必要ともされていない。
            // バイナリファイル作成
            now := time.Now()
            year_val, month_val, day_val := now.Date() // 年月日を数字で取得してみる
            //    fmt.Printf("%04d年%02d月%02d日 %02d:%02d:%02d\n",year_val,int(month_val),day_val, now.Hour(), now.Minute(), now.Second())
            fp_header := fmt.Sprintf("%04d%02d%02d_%02d%02d%02d", year_val, int(month_val), day_val, now.Hour(), now.Minute(), now.Second())
            //    fp_header :=fmt.Sprintf("%04d%02d%02d",year_val,int(month_val),day_val)
            filename := log_bin_path + fp_header + "_AP->" + Machine + "_" + "IH" + "_send.bin"

            // ファイルがなければ作成し、あれば追記する。
            // panic: open .: too many open files の対策として無名関数化
            func() {
                fp, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
                if err != nil {
                    //エラー処理
                    log.Fatal(err)
                }
                defer fp.Close()
                binary.Write(fp, binary.BigEndian, []byte(send_byte.Bytes()))
            }()

            */

            err = make_csv(Machine, "IH", send_byte.Bytes(), log_csv_path)

        }
    }
    t.Stop() // タイマを止める。 <- これがないとメモリリークを起こす
}

/*
   要求の作成
   Input
   client        TCP/IPクライアント
   machine       送信機器
   request_str   要求コマンド種別
   ipaddr        IPアドレスとポート番号
   send_byte     要求コマンドデータ列を詰めるbytes.Buffer
*/
func make_request(client *tcpclient.Client, machine string, request_str string, ipaddr string, send_byte *bytes.Buffer) *bytes.Buffer {

    // 機種がSBOXの場合、シーケンスナンバーをインクリメント(但し、999を限界値とする)
    // 機種がRSUの場合は、SeqNumberがuint16(0-65535)なので、何も考えずにインクリメントする。
    if machine == "SBOX" {
        if SeqNumber < 999 {
            SeqNumber++
        } else {
            SeqNumber = 1
        }
    } else {
        SeqNumber++ // RSUは何も考えずにインクリメント(uint16なので0x0000～0xFFFF)
    }

    //    fmt.Printf("SeqNumber:%d\n",SeqNumber)

    send_byte.Reset()
    send_byte, _ = makecmd.Run(client, machine, request_str, SeqNumber, MachineNo, ipaddr, send_byte)

    return send_byte
}

/*
   要求送信ループ
*/
func send_loop(client *tcpclient.Client) {

    var getcmd findcmd.Req_cmd_st
    send_byte := new(bytes.Buffer) // 送信バイナリ格納用

    // TCP通信キャンセル指示用のコンテキストを作成する
    ctx := context.Background()

    for {
        // コマンド検知
        if len(findcmd.Ch_req_cmd) != 0 {
            getcmd = <-findcmd.Ch_req_cmd // 検知したコマンド取り出し

            // 取り出したコマンドで送信データ作成
            send_byte = make_request(client, getcmd.Machine, getcmd.Command, IpAddress, send_byte)

            //            fmt.Printf("Machine:%s, Command:%s\n",getcmd.Machine, getcmd.Command)

            // コマンド送信処理（リトライ3回）
            err := conretry.Retry(ctx, 3, 0, func() error {
                var ierr error
                _, ierr = client.Send_command(send_byte.Bytes()) // AP → RSU/SBOX 要求送信
                return ierr
            })
            if err != nil {
                panic(err)
            }

            // 要求送信したら、要求依頼ファイルを消す
            if getcmd.Machine == "RSU" {
                switch getcmd.Command {
                case "IH":
                    _ = os.Remove("RIH") // 死活監視要求

                case "IT":
                    _ = os.Remove("RIT") // 時刻校正要求

                case "IA":
                    _ = os.Remove("RIA") // ASK切り替え要求

                case "IQ":
                    _ = os.Remove("RIQ") // QPSK切り替え要求

                case "IC":
                    _ = os.Remove("RIC") // リンク接続要求

                case "IIG":
                    _ = os.Remove("RIIG") // VICS要求 画像

                case "IIM":
                    _ = os.Remove("RIIM") // VICS要求 文字

                case "IIO":
                    _ = os.Remove("RIIO") // VICS要求 音声

                case "ID":
                    _ = os.Remove("RID") // 電波停止要求
                }

            }

            if getcmd.Machine == "SBOX" {
                switch getcmd.Command {
                case "IH":
                    _ = os.Remove("SIH")       // 死活監視要求

                case "IT":
                    _ = os.Remove("SIT")       // 状態通知要求

                case "Ic":
                    _ = os.Remove("SIc")       // ETC情報取得要求

                case "IQ":
                    _ = os.Remove("SIQ")       // 車載機発話認証（SPF認証）要求 : RSUリンク接続応答のLIDをセット

                case "IQS":
                    _ = os.Remove("SIQS")       // 車載機発話認証（SPF認証）要求 : Sc通知のLIDをセット

                case "IRONYES":
                    _ = os.Remove("SIRONYES")  // 無線制御：開始：車載器指示有り

                case "IRONNO":
                    _ = os.Remove("SIRONNO")   // 無線制御：開始：車載器指示無し

                case "IROFFYES":
                    _ = os.Remove("SIROFFYES") // 無線制御：停止：車載器指示有り

                case "IROFFNO":
                    _ = os.Remove("SIROFFNO")  // 無線制御：停止：車載器指示無し
                }

            }

            fmt.Printf("===> AP -> %s : %s %X\n", getcmd.Machine, getcmd.Command, send_byte)
            //log.Printf("===> AP -> %s : %s %X\n", getcmd.Machine, getcmd.Command, send_byte)            // 容量削減のためコメントアウト

            /* 20220829 バイナリファイルの保存は不要と判断。コメント化
            // バイナリファイル作成
            now := time.Now()
            year_val, month_val, day_val := now.Date() // 年月日を数字で取得してみる
            //    fmt.Printf("%04d年%02d月%02d日 %02d:%02d:%02d\n",year_val,int(month_val),day_val, now.Hour(), now.Minute(), now.Second())
            fp_header := fmt.Sprintf("%04d%02d%02d_%02d%02d%02d", year_val, int(month_val), day_val, now.Hour(), now.Minute(), now.Second())
            //    fp_header :=fmt.Sprintf("%04d%02d%02d",year_val,int(month_val),day_val)
            filename := log_bin_path + fp_header + "_AP->" + getcmd.Machine + "_" + getcmd.Command + "_send.bin"

            // ファイルがなければ作成し、あれば追記する。
            // panic: open .: too many open files の対策として無名関数化する
            func() {
                fp, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
                if err != nil {
                    //エラー処理
                    log.Fatal(err)
                }
                defer fp.Close()
                binary.Write(fp, binary.BigEndian, []byte(send_byte.Bytes()))
            }()
            */
            
            // 送信データ履歴(CSVファイル)作成
            // コマンドによって、バイナリファイルを切り分け、CSVフォーマットでデータを保存する
            // fmt.Printf("getcmd.Command:%v\n",getcmd.Command)
            err = make_csv(getcmd.Machine, getcmd.Command, send_byte.Bytes(), log_csv_path)
        }
    }
}


/*
   WCNファイル及び通過履歴ファイルを一本化する。
*/
func make_csv_table() {
    t := time.NewTicker(200 * time.Millisecond) // 200msecタイマー

    defer t.Stop() // タイマを止める。 <- これがないとメモリリークを起こす

    for {
        select {
        case <-t.C:            // タイマー経過した。

            // WCNファイルの一本化
            _, err := exec.Command("./script/make_wcn_table.sh").Output()
            if err != nil {
                log.Fatal(err)
            }

            // 通過履歴ファイルの一本化
            _, err = exec.Command("./script/make_csv_table.sh").Output()
            if err != nil {
                log.Fatal(err)
            }
        }
    }
}


/*
   CSVファイルの作成と保存
   コマンドによってバイナリファイルの分割を制御し、CSVファイルを作成する
   mch         : SBOX / RSU / LOG
   cmd         : 要求 / 応答
   write_data  : 送受信バイト列
   fpath       : 保存ファイルパス

   Retugn      : error or nil
*/
func make_csv(mch string, cmd string, write_data []byte, fpath string) error {

    var req_flg bool      // 要求:true 応答:false
    var sc_csv_save bool  // Sc通知受診時、正しいデータを受信できたか否か判定用。  保存する:true　保存しない:false

    // SBOX電文用ヘッダー
    var head_dtime []byte
    var head_cmd []byte
    var head_seqno []byte
    var head_total_size []byte
    var head_data_size []byte

    // CSV保存用
    var write_title_csv string
    var write_data_csv string
    var write_wcn_csv string
    //    var write_ic_card_no_csv string
    var write_ic_card_id string           // ETCカードID
    var sikyoku_code string               // ナンバープレート：陸運支局コード
    var youto_code string                 // ナンバープレート：用途コード
    var bunrui_number string              // ナンバープレート：車種分類番号
    var ichiren_number string             // ナンバープレート：一連番号

    // RSU電文用ヘッダー
    var head_send_mno []byte
    var head_receive_mno []byte
    //    var head_seqno []byte     // <-- SBOX側にも同じ名前の変数があるのでそっちを利用する
    var head_syubetsu []byte
    var head_if_info []byte
    var head_syubetu_info []byte
    var head_yobi []byte
    var head_cmd_syubetsu []byte
    var head_cmd_data_size []byte

    req_flg = true
    sc_csv_save = false


    // アンテナ関連
    var a_num = 0               // アンテナ番号
    var rsu_name string         // RSU1〜RSU4
    var rsu_name2 string        // A1〜A4
    var write_status string     // IN、PARK、OUT

    // SBOXの場合
    if strings.Contains(IpAddress, "58001") == true {
        a_num = 1
        rsu_name = "RSU01"
        rsu_name2 = "A1"
        write_status = "IN"
    }

    if strings.Contains(IpAddress, "58002") == true {
        a_num = 2
        rsu_name = "RSU02"
        rsu_name2 = "A2"
        write_status = "PARK"
    }

    if strings.Contains(IpAddress, "58003") == true {
        a_num = 3
        rsu_name = "RSU03"
        rsu_name2 = "A3"
        write_status = "OUT"
    }

    if strings.Contains(IpAddress, "58004") == true {
        a_num = 4
        rsu_name = "RSU04"
        rsu_name2 = "A4"
        write_status = "OUT"
    }

    // RSUの場合
    if strings.Contains(IpAddress, "192.168.110.11:") == true {
        a_num = 1
        rsu_name = "RSU01"
        rsu_name2 = "A1"
        write_status = "IN"
    }

    if strings.Contains(IpAddress, "192.168.110.12:") == true {
        a_num = 2
        rsu_name = "RSU02"
        rsu_name2 = "A2"
        write_status = "PARK"
    }

    if strings.Contains(IpAddress, "192.168.110.13") == true {
        a_num = 3
        rsu_name = "RSU03"
        rsu_name2 = "A3"
        write_status = "OUT"
    }

    if strings.Contains(IpAddress, "192.168.110.14") == true {
        a_num = 4
        rsu_name = "RSU04"
        rsu_name2 = "A4"
        write_status = "OUT"
    }

    
    // バイナリ→文字列変換
    switch mch {
    case "SBOX":

        // SBOXヘッダー部のCSV化
        head_dtime = write_data[0:16]       //  0~15: 16byte 電文送信時刻
        head_cmd = write_data[16:18]        // 16~17:  2byte コマンド種別
        head_seqno = write_data[18:21]      // 18~20:  3byte シーケンスNo
        head_total_size = write_data[21:27] // 21~26:  6byte トータルサイズ
        head_data_size = write_data[27:32]  // 27~31:  5byte データサイズ

        switch cmd {
        case "IH": // 死活監視

            write_title_csv = fmt.Sprintf("電文送信時刻,コマンド種別,シーケンスNo,トータルサイズ,データサイズ\n")
            write_data_csv = fmt.Sprintf("%s,%s,%s,%s,%s\n", string(head_dtime), string(head_cmd), string(head_seqno), string(head_total_size), string(head_data_size))

        case "IT": // 状態通知

            write_title_csv = fmt.Sprintf("電文送信時刻,コマンド種別,シーケンスNo,トータルサイズ,データサイズ\n")
            write_data_csv = fmt.Sprintf("%s,%s,%s,%s,%s\n", string(head_dtime), string(head_cmd), string(head_seqno), string(head_total_size), string(head_data_size))

        case "IQ","IQS": // 車載器発話認証
            data_lid := write_data[32:40] // LID 8byte

            write_title_csv = fmt.Sprintf("電文送信時刻,コマンド種別,シーケンスNo,トータルサイズ,データサイズ,LID\n")
            write_data_csv = fmt.Sprintf("%s,%s,%s,%s,%s,%s\n", string(head_dtime), string(head_cmd), string(head_seqno), string(head_total_size), string(head_data_size), string(data_lid))

        case "Ic": // ETC情報取得
            data_shiji := write_data[32:34] // 車載器指示フラグ
            data_info := write_data[34:36]  // 車載器指示情報

            write_title_csv = fmt.Sprintf("電文送信時刻,コマンド種別,シーケンスNo,トータルサイズ,データサイズ,車載器指示フラグ,車載器指示情報\n")
            write_data_csv = fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s\n", string(head_dtime), string(head_cmd), string(head_seqno), string(head_total_size), string(head_data_size), string(data_shiji), string(data_info))

        case "IRONYES","IRONNO","IROFFYES","IROFFNO": // 無線制御
            data_ctrl := write_data[32:34]  // 制御
            data_shiji := write_data[34:36] // 車載器指示フラグ
            data_info := write_data[36:38]  // 車載器指示情報

            write_title_csv = fmt.Sprintf("電文送信時刻,コマンド種別,シーケンスNo,トータルサイズ,データサイズ,制御,車載器指示フラグ,車載器指示情報\n")
            write_data_csv = fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s\n", string(head_dtime), string(head_cmd), string(head_seqno), string(head_total_size), string(head_data_size),string(data_ctrl) ,string(data_shiji), string(data_info))

        case "AH": // 死活監視（応答）
            data_kekka := write_data[32:34] // 結果コード

            write_title_csv = fmt.Sprintf("電文送信時刻,コマンド種別,シーケンスNo,トータルサイズ,データサイズ,結果コード\n")
            write_data_csv = fmt.Sprintf("%s,%s,%s,%s,%s,%s\n", string(head_dtime), string(head_cmd), string(head_seqno), string(head_total_size), string(head_data_size), string(data_kekka))

            req_flg = false

        case "ST": // 状態通知（通知）
            data_kekka := write_data[32:34] // 結果コード

            //結果コードによってデータ部のあり方が異なる
            if string(data_kekka) == "00" { // 正常

                data_tuchi := write_data[34:38] // 通知コード
                data_sam := write_data[38:40]   // SAM鍵データ状態
                data_spf := write_data[40:42]   // SPF鍵データ状態
                data_musen := write_data[42:44] // 無線部状態
                data_yobi := write_data[44:48]  // 予備

                write_title_csv = fmt.Sprintf("電文送信時刻,コマンド種別,シーケンスNo,トータルサイズ,データサイズ,結果コード,通知コード,SAM鍵データ状態,SPF鍵データ状態,無線部状態,予備\n")
                write_data_csv = fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s\n", string(head_dtime), string(head_cmd), string(head_seqno), string(head_total_size), string(head_data_size), string(data_kekka), string(data_tuchi), string(data_sam), string(data_spf), string(data_musen), string(data_yobi))
            } else {
                write_title_csv = fmt.Sprintf("電文送信時刻,コマンド種別,シーケンスNo,トータルサイズ,データサイズ,結果コード\n")
                write_data_csv = fmt.Sprintf("%s,%s,%s,%s,%s,%s\n", string(head_dtime), string(head_cmd), string(head_seqno), string(head_total_size), string(head_data_size), string(data_kekka))
            }

            req_flg = false

        case "AT": // 状態通知（応答）
            data_kekka := write_data[32:34] // 結果コード

            //結果コードによってデータ部のあり方が異なる
            if string(data_kekka) == "00" { // 正常

                data_tuchi := write_data[34:38] // 通知コード
                data_sam := write_data[38:40]   // SAM鍵データ状態
                data_spf := write_data[40:42]   // SPF鍵データ状態
                data_musen := write_data[42:44] // 無線部状態
                data_yobi := write_data[44:48]  // 予備

                write_title_csv = fmt.Sprintf("電文送信時刻,コマンド種別,シーケンスNo,トータルサイズ,データサイズ,結果コード,通知コード,SAM鍵データ状態,SPF鍵データ状態,無線部状態,予備\n")
                write_data_csv = fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s\n", string(head_dtime), string(head_cmd), string(head_seqno), string(head_total_size), string(head_data_size), string(data_kekka), string(data_tuchi), string(data_sam), string(data_spf), string(data_musen), string(data_yobi))
            } else {
                write_title_csv = fmt.Sprintf("電文送信時刻,コマンド種別,シーケンスNo,トータルサイズ,データサイズ,結果コード\n")
                write_data_csv = fmt.Sprintf("%s,%s,%s,%s,%s,%s\n", string(head_dtime), string(head_cmd), string(head_seqno), string(head_total_size), string(head_data_size), string(data_kekka))
            }

            req_flg = false

        case "AQ": // 車載器発話認証（応答）
            data_kekka := write_data[32:34] // 結果コード

            write_title_csv = fmt.Sprintf("電文送信時刻,コマンド種別,シーケンスNo,トータルサイズ,データサイズ,結果コード\n")
            write_data_csv = fmt.Sprintf("%s,%s,%s,%s,%s,%s\n", string(head_dtime), string(head_cmd), string(head_seqno), string(head_total_size), string(head_data_size), string(data_kekka))

            req_flg = false

        case "Ac": // ETC情報取得（応答） 2022/08/30 Icコマンドが廃止になった為、この応答が帰ってくることはないが、消さずにコードは残しておく。
            data_kekka := write_data[32:34] // 結果コード        2byte

            //結果コードによってデータ部のあり方が異なる
            if string(data_kekka) == "00" || string(data_kekka) == "22" { // 正常
                data_lid := write_data[34:42]                    // LID                8byte
                data_wcn := write_data[42:54]                    // WCN               12byte
                write_wcn_csv = string(write_data[42:54])        // WCN(別ファイル用) 12byte
                data_idnr := write_data[54:70]                   // IDNA              16byte
                data_ic_card_no := write_data[70:86]             // IC Card No        16byte
                //                write_ic_card_no_csv = string(write_data[70:86]) // IC Card No        16byte (ac_f2連携用)
                data_car_kanri_no := write_data[86:102]          // 車載器管理番号(別ファイル用)   16byte
                data_ic_card_id := write_data[102:122]           // IC Card ID        20byte
                write_ic_card_id = string(write_data[102:122])   // IC Card ID        20byte (ac_f2連携用)
                data_keiyaku := write_data[122:304]              // 契約情報        182byte
                data_car_koyu := write_data[304:402]             // 車載器固有情報  98byte
                data_car_status := write_data[402:404]           // 車載器状態有無   2byte
                data_car_struct := write_data[404:416]           // 車載器構成情報  12byte
                data_car_yobi := write_data[416:418]             // 予備              2byte

                write_title_csv = fmt.Sprintf("電文送信時刻,コマンド種別,シーケンスNo,トータルサイズ,データサイズ,結果コード,LID,WCN,IDNA,ICカード管理番号,車載器管理番号,ICカードID,契約情報,車載器固有情報,車載器状態有無,車載器構成情報,予備\n")
                write_data_csv = fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s\n",
                    string(head_dtime),
                    string(head_cmd),
                    string(head_seqno),
                    string(head_total_size),
                    string(head_data_size),
                    string(data_kekka),
                    string(data_lid),
                    string(data_wcn),
                    string(data_idnr),
                    string(data_ic_card_no),
                    string(data_car_kanri_no),
                    string(data_ic_card_id),
                    string(data_keiyaku),
                    string(data_car_koyu),
                    string(data_car_status),
                    string(data_car_struct),
                    string(data_car_yobi),
                )

            } else {
                write_title_csv = fmt.Sprintf("電文送信時刻,コマンド種別,シーケンスNo,トータルサイズ,データサイズ,結果コード\n")
                write_data_csv = fmt.Sprintf("%s,%s,%s,%s,%s,%s\n", string(head_dtime), string(head_cmd), string(head_seqno), string(head_total_size), string(head_data_size), string(data_kekka))
            }

            req_flg = false

        case "Sc": // ETC情報取得（通知）
            data_kekka := write_data[32:34] // 結果コード        2byte

            //結果コードによってデータ部のあり方が異なる
            if string(data_kekka) == "00" || string(data_kekka) == "22" { // 正常

                data_lid := write_data[34:42]                    // LID                8byte
                data_wcn_info := write_data[42:44]               // WCN格納情報        2byte
                data_wcn := write_data[44:56]                    // WCN               12byte
                write_wcn_csv = string(write_data[44:56])        // WCN(別ファイル用) 12byte
                data_idnr := write_data[56:72]                   // IDNA              16byte
                data_ic_card_no := write_data[72:88]             // IC Card No        16byte
                //                write_ic_card_no_csv = string(write_data[70:86]) // IC Card No        16byte (ac_f2連携用)
                data_car_kanri_no := write_data[88:104]          // 車載器管理番号(別ファイル用)   16byte
                data_ic_card_id := write_data[104:124]           // IC Card ID        20byte
                write_ic_card_id = string(write_data[104:124])   // IC Card ID        20byte (ac_f2連携用)

                sikyoku_code = string(write_data[325:331])       // 陸運局支局コード
                youto_code = string(write_data[331:333])         // 用途コード
                bunrui_number = string(write_data[333:336])      // 車種分類番号
                ichiren_number = string(write_data[336:340])     // 一連番号          

                data_keiyaku := write_data[124:306]              // 契約情報        182byte
                data_car_koyu := write_data[306:404]             // 車載器固有情報  98byte(ここにナンバープレート情報が入っている)
                data_car_status := write_data[404:406]           // 車載器状態有無   2byte
                data_car_struct := write_data[406:418]           // 車載器構成情報  12byte

                write_title_csv = fmt.Sprintf("電文送信時刻,コマンド種別,シーケンスNo,トータルサイズ,データサイズ,結果コード,LID,WCN格納情報,WCN,IDNA,ICカード管理番号,車載器管理番号,ICカードID,契約情報,車載器固有情報,車載器状態有無,車載器構成情報\n")
                write_data_csv = fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s\n",
                    string(head_dtime),
                    string(head_cmd),
                    string(head_seqno),
                    string(head_total_size),
                    string(head_data_size),
                    string(data_kekka),
                    string(data_lid),
                    string(data_wcn_info),
                    string(data_wcn),
                    string(data_idnr),
                    string(data_ic_card_no),
                    string(data_car_kanri_no),
                    string(data_ic_card_id),
                    string(data_keiyaku),
                    string(data_car_koyu),
                    string(data_car_status),
                    string(data_car_struct),
                )
                
                sc_csv_save = true

            } else {
                write_title_csv = fmt.Sprintf("電文送信時刻,コマンド種別,シーケンスNo,トータルサイズ,データサイズ,結果コード\n")
                write_data_csv = fmt.Sprintf("%s,%s,%s,%s,%s,%s\n", string(head_dtime), string(head_cmd), string(head_seqno), string(head_total_size), string(head_data_size), string(data_kekka))
            }

            req_flg = false

        case "AR":
            data_kekka := write_data[32:34] // 結果コード

            write_title_csv = fmt.Sprintf("電文送信時刻,コマンド種別,シーケンスNo,トータルサイズ,データサイズ,結果コード\n")
            write_data_csv = fmt.Sprintf("%s,%s,%s,%s,%s,%s\n", string(head_dtime), string(head_cmd), string(head_seqno), string(head_total_size), string(head_data_size), string(data_kekka))

            req_flg = false

        }

    case "RSU":

        // RSUヘッダー部のCSV化
        head_send_mno = write_data[0:4]        // 送信先機器番号    4byte
        head_receive_mno = write_data[4:8]     // 送信元機器番号    4byte
        head_seqno = write_data[8:10]          // シーケンス番号    2byte
        head_syubetsu = write_data[10:11]      // 電文種別           1byte
        head_if_info = write_data[11:12]       // I/F 情報           1byte
        head_syubetu_info = write_data[12:13]  // 電文種別情報      1byte
        head_yobi = write_data[13:16]          // 予備               3byte
        head_cmd_syubetsu = write_data[16:18]  // コマンド種別      2byte
        head_cmd_data_size = write_data[18:20] // コマンドデータ長  2byte

        switch cmd {
        case "IH": // 死活監視
            data_kekka := write_data[20:22] // 結果コード        2byte
            write_title_csv = fmt.Sprintf("送信先機器番号,送信元機器番号,シーケンス番号,電文種別,I/F情報,電文種別情報,予備,コマンド種別,コマンドデータ長,結果コード\n")
            write_data_csv = fmt.Sprintf("%X,%X,%X,%X,%X,%X,%X,%X,%X,%X\n",
                head_send_mno,
                head_receive_mno,
                head_seqno,
                head_syubetsu,
                head_if_info,
                head_syubetu_info,
                head_yobi,
                head_cmd_syubetsu,
                head_cmd_data_size,
                data_kekka)

        case "IT": // 時刻校正
            data_kekka := write_data[20:22] // 結果コード        2byte
            data_year := write_data[22:24]  // 年                 2byte
            data_month := write_data[24:25] // 月                 1byte
            data_day := write_data[25:26]   // 日                 1byte
            data_hour := write_data[26:27]  // 時                 1byte
            data_min := write_data[27:28]   // 分                 1byte
            data_sec := write_data[28:29]   // 秒                 1byte
            data_yobi := write_data[29:30]  // 予備               1byte

            write_title_csv = fmt.Sprintf("送信先機器番号,送信元機器番号,シーケンス番号,電文種別,I/F情報,電文種別情報,予備,コマンド種別,コマンドデータ長,結果コード,年,月,日,時,分,秒,予備\n")
            write_data_csv = fmt.Sprintf("%X,%X,%X,%X,%X,%X,%X,%X,%X,%X,%X,%X,%X,%X,%X,%X,%X\n",
                head_send_mno,
                head_receive_mno,
                head_seqno,
                head_syubetsu,
                head_if_info,
                head_syubetu_info,
                head_yobi,
                head_cmd_syubetsu,
                head_cmd_data_size,
                data_kekka,
                data_year,
                data_month,
                data_day,
                data_hour,
                data_min,
                data_sec,
                data_yobi)

        case "IA": // ASK切替
            data_kekka := write_data[20:22] // 結果コード        2byte
            write_title_csv = fmt.Sprintf("送信先機器番号,送信元機器番号,シーケンス番号,電文種別,I/F情報,電文種別情報,予備,コマンド種別,コマンドデータ長,結果コード\n")
            write_data_csv = fmt.Sprintf("%X,%X,%X,%X,%X,%X,%X,%X,%X,%X\n",
                head_send_mno,
                head_receive_mno,
                head_seqno,
                head_syubetsu,
                head_if_info,
                head_syubetu_info,
                head_yobi,
                head_cmd_syubetsu,
                head_cmd_data_size,
                data_kekka)

        case "IQ": // QPSK切替
            data_kekka := write_data[20:22] // 結果コード        2byte
            write_title_csv = fmt.Sprintf("送信先機器番号,送信元機器番号,シーケンス番号,電文種別,I/F情報,電文種別情報,予備,コマンド種別,コマンドデータ長,結果コード\n")
            write_data_csv = fmt.Sprintf("%X,%X,%X,%X,%X,%X,%X,%X,%X,%X\n",
                head_send_mno,
                head_receive_mno,
                head_seqno,
                head_syubetsu,
                head_if_info,
                head_syubetu_info,
                head_yobi,
                head_cmd_syubetsu,
                head_cmd_data_size,
                data_kekka)

        case "IC": // リンク接続
            data_kekka := write_data[20:22] // 結果コード        2byte
            write_title_csv = fmt.Sprintf("送信先機器番号,送信元機器番号,シーケンス番号,電文種別,I/F情報,電文種別情報,予備,コマンド種別,コマンドデータ長,結果コード\n")
            write_data_csv = fmt.Sprintf("%X,%X,%X,%X,%X,%X,%X,%X,%X,%X\n",
                head_send_mno,
                head_receive_mno,
                head_seqno,
                head_syubetsu,
                head_if_info,
                head_syubetu_info,
                head_yobi,
                head_cmd_syubetsu,
                head_cmd_data_size,
                data_kekka)

        case "IIG", "IIM", "IIO":
            data_kekka := write_data[20:22]    // 結果コード            2 byte
            data_del_flg := write_data[22:23]  // 情報登録削除フラグ   1 byte
            data_yobi := write_data[23:24]     // 予備                   1 byte
            data_daikubun := write_data[24:26] // 大区分データ数        2 byte

            data_size := write_data[26:28]  // データサイズ          2 byte
            data_ippan := write_data[28:29] // 一般/優先              1 byte
            data_yobi2 := write_data[29:30] // 予備                   1 byte

            // データサイズがバイト配列なので、数値に変換。
            // ※固定バイト配列 → スライス → 数値
            var size_slice []byte = data_size[:]
            size := binary.BigEndian.Uint16(size_slice)

            /*
               log.Printf("結果コード : 0x%X\n",data_kekka)
               log.Printf("情報登録削除フラグ : 0x0%X\n",data_del_flg)
               log.Printf("予備 : 0x%X\n",data_yobi)
               log.Printf("大区分データ数 : 0x0%X\n",data_daikubun)
               log.Printf("データサイズ : %d\n",size)
               log.Printf("一般/優先 : 0x%X\n",data_ippan)
               log.Printf("予備 : 0x%X\n",data_yobi2)
            */

            data_val := write_data[30 : 30+size] // データの内容（上記データサイズがバイト範囲）
            //            log.Printf("内容 : 0x%X\n",data_val)

            // まずは大区分が1つだけと仮定してCSV文字列作成(改行なし)
            write_title_csv = fmt.Sprintf("送信先機器番号,送信元機器番号,シーケンス番号,電文種別,I/F情報,電文種別情報,予備,コマンド種別,コマンドデータ長,結果コード,情報登録削除フラグ,予備,大区分データ数,データサイズ,一般/優先,予備,データの内容\n")
            write_data_csv = fmt.Sprintf("%X,%X,%X,%X,%X,%X,%X,%X,%X,%X,%X,%X,%X,%X,%X,%X,%X",
                head_send_mno,
                head_receive_mno,
                head_seqno,
                head_syubetsu,
                head_if_info,
                head_syubetu_info,
                head_yobi,
                head_cmd_syubetsu,
                head_cmd_data_size,
                data_kekka,
                data_del_flg,
                data_yobi,
                data_daikubun,
                data_size,
                data_ippan,
                data_yobi2,
                data_val)

            write_data_csv = write_data_csv + "\n" // 改行を追加

        case "ID": // 電波停止
            data_kekka := write_data[20:22] // 結果コード        2byte
            write_title_csv = fmt.Sprintf("送信先機器番号,送信元機器番号,シーケンス番号,電文種別,I/F情報,電文種別情報,予備,コマンド種別,コマンドデータ長,結果コード\n")
            write_data_csv = fmt.Sprintf("%X,%X,%X,%X,%X,%X,%X,%X,%X,%X\n",
                head_send_mno,
                head_receive_mno,
                head_seqno,
                head_syubetsu,
                head_if_info,
                head_syubetu_info,
                head_yobi,
                head_cmd_syubetsu,
                head_cmd_data_size,
                data_kekka)

        case "AH": // 死活監視(応答)
            data_kekka := write_data[20:22] // 結果コード        2byte
            write_title_csv = fmt.Sprintf("送信先機器番号,送信元機器番号,シーケンス番号,電文種別,I/F情報,電文種別情報,予備,コマンド種別,コマンドデータ長,結果コード\n")
            write_data_csv = fmt.Sprintf("%X,%X,%X,%X,%X,%X,%X,%X,%X,%X\n",
                head_send_mno,
                head_receive_mno,
                head_seqno,
                head_syubetsu,
                head_if_info,
                head_syubetu_info,
                head_yobi,
                head_cmd_syubetsu,
                head_cmd_data_size,
                data_kekka)

            req_flg = false

        case "AT": // 時刻校正(応答)
            data_kekka := write_data[20:22] // 結果コード        2byte
            write_title_csv = fmt.Sprintf("送信先機器番号,送信元機器番号,シーケンス番号,電文種別,I/F情報,電文種別情報,予備,コマンド種別,コマンドデータ長,結果コード\n")
            write_data_csv = fmt.Sprintf("%X,%X,%X,%X,%X,%X,%X,%X,%X,%X\n",
                head_send_mno,
                head_receive_mno,
                head_seqno,
                head_syubetsu,
                head_if_info,
                head_syubetu_info,
                head_yobi,
                head_cmd_syubetsu,
                head_cmd_data_size,
                data_kekka)

            req_flg = false

        case "AA": // ASK切替(応答)
            data_kekka := write_data[20:22] // 結果コード        2byte
            write_title_csv = fmt.Sprintf("送信先機器番号,送信元機器番号,シーケンス番号,電文種別,I/F情報,電文種別情報,予備,コマンド種別,コマンドデータ長,結果コード\n")
            write_data_csv = fmt.Sprintf("%X,%X,%X,%X,%X,%X,%X,%X,%X,%X\n",
                head_send_mno,
                head_receive_mno,
                head_seqno,
                head_syubetsu,
                head_if_info,
                head_syubetu_info,
                head_yobi,
                head_cmd_syubetsu,
                head_cmd_data_size,
                data_kekka)

            req_flg = false

        case "AQ": // QPSK切替(応答)
            data_kekka := write_data[20:22] // 結果コード        2byte
            write_title_csv = fmt.Sprintf("送信先機器番号,送信元機器番号,シーケンス番号,電文種別,I/F情報,電文種別情報,予備,コマンド種別,コマンドデータ長,結果コード\n")
            write_data_csv = fmt.Sprintf("%X,%X,%X,%X,%X,%X,%X,%X,%X,%X\n",
                head_send_mno,
                head_receive_mno,
                head_seqno,
                head_syubetsu,
                head_if_info,
                head_syubetu_info,
                head_yobi,
                head_cmd_syubetsu,
                head_cmd_data_size,
                data_kekka)

            req_flg = false

        case "AC": // リンク接続(応答)
            data_kekka := write_data[20:22] // 結果コード        2byte
            data_lid := write_data[22:26]   // LID                4byte
            data_aslid := write_data[26:32] // ASL-ID             6byte

            write_title_csv = fmt.Sprintf("送信先機器番号,送信元機器番号,シーケンス番号,電文種別,I/F情報,電文種別情報,予備,コマンド種別,コマンドデータ長,結果コード,LID,ASL-ID\n")
            write_data_csv = fmt.Sprintf("%X,%X,%X,%X,%X,%X,%X,%X,%X,%X,%X,%X\n",
                head_send_mno,
                head_receive_mno,
                head_seqno,
                head_syubetsu,
                head_if_info,
                head_syubetu_info,
                head_yobi,
                head_cmd_syubetsu,
                head_cmd_data_size,
                data_kekka,
                data_lid,
                data_aslid)

            req_flg = false

        case "AI": // VICS要求(応答)

            data_kekka := write_data[20:22] // 結果コード        2byte
            write_title_csv = fmt.Sprintf("送信先機器番号,送信元機器番号,シーケンス番号,電文種別,I/F情報,電文種別情報,予備,コマンド種別,コマンドデータ長,結果コード\n")
            write_data_csv = fmt.Sprintf("%X,%X,%X,%X,%X,%X,%X,%X,%X,%X\n",
                head_send_mno,
                head_receive_mno,
                head_seqno,
                head_syubetsu,
                head_if_info,
                head_syubetu_info,
                head_yobi,
                head_cmd_syubetsu,
                head_cmd_data_size,
                data_kekka)

            req_flg = false

        case "AD": // 電波停止(応答)
            data_kekka := write_data[20:22] // 結果コード        2byte
            write_title_csv = fmt.Sprintf("送信先機器番号,送信元機器番号,シーケンス番号,電文種別,I/F情報,電文種別情報,予備,コマンド種別,コマンドデータ長,結果コード\n")
            write_data_csv = fmt.Sprintf("%X,%X,%X,%X,%X,%X,%X,%X,%X,%X\n",
                head_send_mno,
                head_receive_mno,
                head_seqno,
                head_syubetsu,
                head_if_info,
                head_syubetu_info,
                head_yobi,
                head_cmd_syubetsu,
                head_cmd_data_size,
                data_kekka)

            req_flg = false

        case "NJ": // 装置状態(通知)
            data_kekka := write_data[20:22]  // 結果コード        2byte
            data_seigyo := write_data[22:23] // 制御部異常        1byte
            data_musen := write_data[23:24]  // 無線部異常        1byte
            data_mode := write_data[24:25]   // モード            1byte

            write_title_csv = fmt.Sprintf("送信先機器番号,送信元機器番号,シーケンス番号,電文種別,I/F情報,電文種別情報,予備,コマンド種別,コマンドデータ長,結果コード,制御部異常,無線部異常,モード\n")
            write_data_csv = fmt.Sprintf("%X,%X,%X,%X,%X,%X,%X,%X,%X,%X,%X,%X,%X\n",
                head_send_mno,
                head_receive_mno,
                head_seqno,
                head_syubetsu,
                head_if_info,
                head_syubetu_info,
                head_yobi,
                head_cmd_syubetsu,
                head_cmd_data_size,
                data_kekka,
                data_seigyo,
                data_musen,
                data_mode)

            req_flg = false
            log.Printf("NJ:title: %s\n",write_title_csv)
            log.Printf("NJ:data : %s\n",write_data_csv)

        case "NF": // リンク切断(通知)
            data_kekka := write_data[20:22] // 結果コード        2byte
            data_lid := write_data[22:26]   // LID                4byte

            write_title_csv = fmt.Sprintf("送信先機器番号,送信元機器番号,シーケンス番号,電文種別,I/F情報,電文種別情報,予備,コマンド種別,コマンドデータ長,結果コード,LID\n")
            write_data_csv = fmt.Sprintf("%X,%X,%X,%X,%X,%X,%X,%X,%X,%X,%X\n",
                head_send_mno,
                head_receive_mno,
                head_seqno,
                head_syubetsu,
                head_if_info,
                head_syubetu_info,
                head_yobi,
                head_cmd_syubetsu,
                head_cmd_data_size,
                data_kekka,
                data_lid)

            req_flg = false

        }
    case "LOG":

        // RSUヘッダー部のCSV化
        head_send_mno = write_data[0:4]        // 送信先機器番号    4byte
        head_receive_mno = write_data[4:8]     // 送信元機器番号    4byte
        head_seqno = write_data[8:10]          // シーケンス番号    2byte
        head_syubetsu = write_data[10:11]      // 電文種別           1byte
        head_if_info = write_data[11:12]       // I/F 情報           1byte
        head_syubetu_info = write_data[12:13]  // 電文種別情報      1byte
        head_yobi = write_data[13:16]          // 予備               3byte
        head_cmd_syubetsu = write_data[16:18]  // コマンド種別      2byte
        head_cmd_data_size = write_data[18:20] // コマンドデータ長  2byte

        switch cmd {
        case "NI": // ログ通知
            data_kekka := write_data[20:22] // 結果コード        2byte
            log_data := write_data[22:]     // 日付以降はひとまとめに保存。データは固定長ではない。 ?byte
            write_title_csv = fmt.Sprintf("送信先機器番号,送信元機器番号,シーケンス番号,電文種別,I/F情報,電文種別情報,予備,コマンド種別,コマンドデータ長,結果コード,ログデータ\n")
            write_data_csv = fmt.Sprintf("%X,%X,%X,%X,%X,%X,%X,%X,%X,%X,%X\n",
                head_send_mno,
                head_receive_mno,
                head_seqno,
                head_syubetsu,
                head_if_info,
                head_syubetu_info,
                head_yobi,
                head_cmd_syubetsu,
                head_cmd_data_size,
                data_kekka,
                log_data)
            
            req_flg = false
        }

    default:
        // Nothing
    }

    // 送信/受信に応じて、ファイル名の文字列を変える
    var send_or_receive string
    if req_flg == true {
        send_or_receive = "send"
    } else {
        send_or_receive = "receive"
    }

    /* --- CSVファイル名作成 ---
    　送信・受信の両方でファイル名を作成している。
    　send_or_receiveの内容を上記のIF文で切り替えている。
    　保存先は「log_csv_path → log/csv」に保存する。
    */ 

    var fp_header string = ""
    switch mch {
    case "SBOX":
        fp_header = fmt.Sprintf("2%s", string(head_dtime))  // 受信したデータの電文送信時刻
    default:
        // 現在時間を取得
        now := time.Now()
        nowUTC := now.UTC()

        // ミリ秒の算出（文字列変換含む）
        t2 := nowUTC.UnixNano() / int64(time.Millisecond)    // 時間(ナノ秒)を時間(ミリ秒)に変換
        t2_str := strconv.Itoa(int(t2))                      // 時間(ミリ秒)を文字列に変換
        ms_str := t2_str[len(t2_str)-3:]                     // 時間(ミリ秒)文字列から、ミリ秒部分だけを切り出す

        // ミリ秒を含む日時データを作成
        year_val, month_val, day_val := now.Date() // 年月日を数字で取得してみる
        fp_header = fmt.Sprintf("%04d%02d%02d%02d%02d%02d%s", year_val, int(month_val), day_val, now.Hour(), now.Minute(), now.Second(),ms_str)
    }

    // 通信内容ファイル名作成(csvファイル)
    var filename = ""
    if send_or_receive == "_send" {
        filename = log_csv_path + fp_header + "_AP_to_" + mch + "_" + cmd + "_" + rsu_name2 + "_" + send_or_receive + ".csv"
    } else {
        filename = log_csv_path + fp_header + "_" + mch + "_to_AP_" + cmd + "_" + rsu_name2 + "_" + send_or_receive + ".csv"
    }
        
    // すでに同時刻のファイルがなければ作成。
    // msecオーダーでの同時刻ファイルはできないと考えているが、もし同時刻のファイルが存在した場合は何もファイル操作せずに終了する。
    func() {                         // panic: open .: too many open files の対策として無名関数化する

        // CSVファイルにデータ書き込み
        /* --- CSVファイル保存 ---
              保存先はfilenameに文字列として作成されている。
        */ 
        fp, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
        if err != nil {
            //エラー処理
            log.Fatal(err)
        }
        defer fp.Close()

        w := bufio.NewWriter(fp)
        fmt.Fprint(w, write_title_csv) // タイトル
        fmt.Fprint(w, write_data_csv)  // データ
        w.Flush()
    }()


    // SBOXのAc応答、またはETC情報通知だった場合は、受信したWCN情報でファイルを作成する(ac連携用)
    // このif文内部で作成されるcsvファイルは、acによって、全アンテナの通過履歴「WCN_rireki.csv」及び全アンテナのWCN番号履歴「WCN_Table.csv」の作成材料として活用される。
    if mch == "SBOX" && cmd == "Ac" || 
       mch == "SBOX" && cmd == "Sc" && sc_csv_save == true {

        /*
           
           連続でSc通知を受信した場合、下記の条件にてcsvファイルとして保存するか否かを判断する。

           ■条件
           Sc通知されて来たWCN番号が「通過履歴」にまだ登録されていなければ、ファイル作成。
           登録されている場合、登録時の時刻(res)とSc通知された時刻(fp_header)を比較し、config.iniに設定された秒数を経過していたらファイル作成。
           ※渋滞等で車両がアンテナの前でじっとしている場合、2秒間隔でSc通知が発行されるため、車両を検知してから一定時間は通知を受信してもファイルを保存しないようにする。
        */

       /*
           受信したWCN番号がすでにシステムに検知されている場合、前回受信した時間から何秒経過しているかチェックする。
           指定秒数(INIファイル)以内に受信した場合、渋滞とみなしてcsvファイルは作成しない。
       */

        csv_make_flg := false
        log.Printf("wcn_num: %s, a_num: %d\n",write_wcn_csv, a_num )
        result, err := exec.Command("./script/get_receive_time.sh", write_wcn_csv, strconv.Itoa(a_num)).Output()
        if err != nil {
            log.Printf("get_receive_time.sh Error!!: %v\n",err)
            log.Printf("WCN_Number Find Result : %s",string(result))
            log.Fatal(err)
        }

        
        // WCN番号検出結果判定
        res := string(result)
        res = strings.TrimRight(res,"\n")
        log.Printf("./script/get_receive_time.sh -> %s\n",res)
        if res == "NoHit" {    // WCN_rireki.csvファイルはあるけれど、受信したWCNは登録されていない
            
            log.Printf("Receive WCN First Save : %s <- %s\n", fp_header, write_wcn_csv) // WCN番号の受信時間を保存する。

            // CSVデータ作成
            csv_make_flg = true
        }

        // WCN_rireki_a?.csvファイルが存在しており、受信したWCN番号がすでに登録されている。アンテナ番号も間違っていない。
        if res != "NotFound" && res != "NoMachine" && res != "NoHit"{
            a_duration, err := date_duration(res, fp_header) // 前回の受信時間と今回の受信時間を比較し、時間間隔(秒)を取得
            if err != nil {
                log.Printf("Error...result:%s, time_stamp:%s\n",res ,fp_header)
                log.Println(err)
            }
            
            // 時間間隔が設定時間をオーバーしていたらCSVファイルを作成する（渋滞ではないと考える。acによってWCN_rireki.csvに追加される）
            log.Printf("Receive Duration : %d Sec = %s - %s\n",a_duration, fp_header, res)
            if a_duration >= iniread.Config.Detection_interval {

                log.Printf("Receive WCN Save Duration %d Sec : %s <- %s\n",a_duration, fp_header, write_wcn_csv)
                
                // CSVデータ作成
                csv_make_flg = true
            }
        }
          

        // 渋滞と判断する時間を経過しているか？
        if csv_make_flg == true {

            // 受信したWCNをファイル名とする。　→　 tc_wcn_path = "./tc_wcn/"
            filename := tc_wcn_path + "WCN_" + write_wcn_csv + ".csv"

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
                fmt.Fprint(w, write_wcn_csv+","+write_ic_card_id+"\n") // データ WCN番号,ETCカードID
                w.Flush()
            }()

            // YYYYMMDDHHMMSS_WCN_A1.csvファイル作成
            filename = tc_csv_path + fp_header + "_WCN_" + rsu_name2 + ".csv"

            /* 同時刻(ミリ秒まで)のファイルがすでに存在している場合は何もしない。
           ファイルがなければ作成する。
         マイクロ秒単位で同時に複数のTCPパケットを受信することは無いと想定している。
         結果、ミリ秒管理で、ETC情報を受信する度にファイルが作成される。
            */

            func() {        // panic: open .: too many open files の対策として無名関数化する（この無名関数がFatalで死んでも、呼び元のtcにはなんら影響はない）
                fp, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
                if err != nil {
                    //エラー処理
                    fmt.Printf("os.OpenFile Error : %v\n",err)
                    log.Fatal(err)
                }
                defer fp.Close()

                // CSVファイルにデータ書き込み（逆走検知画面で活用予定。ETC情報要求応答受信時に保存される。）
                w := bufio.NewWriter(fp)
                writeline := fp_header + "," + rsu_name + "," + rsu_name2 + "," + write_wcn_csv + "," + write_status + "," + write_ic_card_id + "," + sikyoku_code + "," + youto_code + "," + bunrui_number + "," + ichiren_number + "\n"
                fmt.Fprint(w, writeline) // データ
                if _, err := fmt.Fprint(w, writeline); err != nil {
                    fmt.Printf("fmt.Fprint(w, writeline) Error : %v\n", err)
                }
                
                //w.Flush()
                if err := w.Flush(); err != nil {
                    fmt.Printf("w.Flush() Error : %v\n", err)
                }                
                
                fmt.Printf("Save csv -> %s\n",filename)
            }()
            
        } else {
            fmt.Printf("Don't Save CSV of Sc. Traffic jam reaction. : %s <- %s\n", fp_header, write_wcn_csv)
        }

    }

    return nil
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


/*
   Traffic Counter メイン処理
*/
func main() {

    var flg_err bool = false       // 初期値：異常なし
    send_byte := new(bytes.Buffer) // 送信バイナリ格納用

    // コマンドファイル検出起動（コルーチン）
    findcmd.Run()

    // TCP通信キャンセル指示用のコンテキストを作成する
    ctx := context.Background()

    // S-BOX / RSUへの接続クライアント作成
    //    client := tcpclient.NewClient("192.168.1.150:8080")
    client := tcpclient.NewClient(IpAddress)
    defer client.Close()

    // パラメータに不備があった場合（上記フィルタ処理に引っかからずdefaultで抜けてしまった）
    send_byte, flg_err = makecmd.Run(client, Machine, MCommand, SeqNumber, MachineNo, IpAddress, send_byte)
    if flg_err == true {
        fmt.Println("Usage: rsu01 -s=\"SBOX\" -i=\"192.168.1.201:58001\" -c=\"IH\"")
        fmt.Println("-s : SBOX or RSU")
        fmt.Println("-i : IP Address:PortNo")
        fmt.Println("-c : Command")
        os.Exit(-1)
    }

    go timer_10()           // 10秒タイマースタート(ログファイルのローテーション用)
    go make_csv_table()     // 通過履歴、WCN番号の一覧化（１本化）

    // ポートコネクトしたまま受信を待ち続け、コマンド送信ファイルを検知した時のみ、要求コマンドを送信する。
    port_open_err := conretry.Retry(ctx, 3, 0, func() error {
        ierr := client.Connect()
        return ierr
    })
    if port_open_err != nil {
        panic(port_open_err)
    }

    go receive_loop(client) // RSU/SBOX → AP 応答/通知送信

    // ログシステムとして起動する場合は、要求は発信しない
    if Machine != "LOG" {
        //        fmt.Printf("LOG?? Machine:%s\n", Machine)
        go send_loop(client) // AP → RSU/SBOX 要求送信
        //        go timer_30(client)         // AP → RSU/SBOX 死活監視要求送信(30秒間隔)   // 20200324 30秒毎の死活監視はac側で行う事とする
    }

    // RSU/SBOXとの接続Portが閉じてしまった場合はつなぎ直す
    // for {

    //     port_open_err := conretry.Retry(ctx, 3, 0, func() error {
    //         ierr := client.Connect()
    //         return ierr
    //     })
    //     if port_open_err != nil {
    //         panic(port_open_err)
    //     }

    // }

    select {}
    
}
