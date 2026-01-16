/*
   コマンドバイト列作成パッケージ
*/

package makecmd

import (
	"bytes"
	"encoding/binary"
	"fmt"
    "os"
    "time"
    "log"

	"localhost.com/tcpclient"
)


/* SBOXからのETC情報通知後、作成されているLIDファイルからLID番号を取得する
   
  */
func read_lid_etcinfo(fname string) string {

    fp, err := os.Open(fname)
    if err != nil {
        fmt.Printf("Filename:%s",fname)
        fmt.Println("LID FileOpen Error!!")
        log.Printf("Filename:%s",fname)
        log.Println("LID FileOpen Error!!")
        
    }

    data := make([]byte, 8)
    _, err = fp.Read(data)
    if err != nil {
        panic(err)
    }
    //    fmt.Printf("Read %d bytes: %s\n", count, data)
    return string(data)
}

/* RSU01～RSU04のフォルダを参照し、LIDファイルからLID番号を取得する
   Input
     port_no : SBOXのポート番号

   Output
     result  : LID番号
*/
func read_lid(port_no string) string{
    var fname string
    
    switch port_no {
    case "58001":
        fname = "../oki_rsu01/LID.bin"
    case "58002":
        fname = "../oki_rsu02/LID.bin"
    case "58003":
        fname = "../oki_rsu03/LID.bin"
    case "58004":
        fname = "../oki_rsu04/LID.bin"
    default:
        fname = "../oki_rsu01/LID.bin"
    }

    fp, err := os.Open(fname)
    if err != nil {
        fmt.Printf("Filename:%s",fname)
        fmt.Println("LID FileOpen Error!!")
        log.Printf("Filename:%s",fname)
        log.Println("LID FileOpen Error!!")
        
    }

    data := make([]byte, 8)
    _, err = fp.Read(data)
    if err != nil {
        panic(err)
    }
    //    fmt.Printf("Read %d bytes: %s\n", count, data)
    return string(data)
}

/*
   与えられたパラメータ(機種種別とコマンド)から、送信用バイナリを作成する。

   ■Input
   client     *tcpclient.Client : TCPクライアント構造体のポインタ
   machine    string            : RSU or SBOX or ME93
   mcmd       string            : 各装置へのコマンド
   seq_no     int16             : 送受信用シーケンス番号
   machine_no int               : 機器番号
   ipaddr     string            : IPアドレスとポート番号
   send_byte  []byte            : 送信バイナリ

   ■Return
   send_byte  : 送信バイナリ
   flg_err    : コマンドエラー検知用フラグ(ありえないコマンド指定された時にfalseとなる)
*/
func Run(client *tcpclient.Client, machine string, mcmd string, seq_no uint16, machine_no int, ipaddr string, send_byte *bytes.Buffer) (*bytes.Buffer, bool) {

    var flg_err bool = false                            // 初期値：異常なし

    seq_str := fmt.Sprintf("%03d",seq_no)                // シーケンス番号作成
    
    /*
       コマンドライン引数:Machineによって、SBOXコマンドなのかRSUコマンドなのかを切り分ける。
       RSU、SBOXそれぞれのコマンドに応じてデータを作成する
    */
    switch machine {
    case "SBOX":
        switch mcmd {

        case "IH": // 死活監視

            // SBOX通信用ヘッダー作成
            tcpclient.SBOX_Header = tcpclient.MakeSBOXHeader("IH",seq_str,"000032","00000")

            // 死活監視データ情報作成
            // → 死活監視は、データ部無し

            // 要求ヘッダーとデータを結合し、一本のバイナリデータに
            binary.Write(send_byte, binary.BigEndian, &tcpclient.SBOX_Header)

        case "IT": // 状態通知

            // SBOX通信用ヘッダー作成
            tcpclient.SBOX_Header = tcpclient.MakeSBOXHeader("IT",seq_str,"000032","00000")

            // 状態通知データ情報作成
            // → 状態通知は、データ部無し
            
            // 要求ヘッダーとデータを結合し、一本のバイナリデータに
            binary.Write(send_byte, binary.BigEndian, &tcpclient.SBOX_Header)

        case "IQ": // 車載器発話認証（RSU：リンク接続応答のLIDを利用して認証）
            
            // SBOX通信用ヘッダー作成
            tcpclient.SBOX_Header = tcpclient.MakeSBOXHeader("IQ",seq_str,"000040","00008")

            //車載器発話認証データ情報作成(対象のRSUから取得)
            //fmt.Printf("Port:%s\n",ipaddr[len(ipaddr)-5:])
            lid := read_lid(ipaddr[len(ipaddr)-5:])
            tcpclient.SBOXData_SPF = tcpclient.MakeSBOXData_SPFInfo(lid)

            // 要求ヘッダーとデータを結合し、一本のバイナリデータに
            binary.Write(send_byte, binary.BigEndian, &tcpclient.SBOX_Header)
            binary.Write(send_byte, binary.BigEndian, &tcpclient.SBOXData_SPF)

        case "IQS": // 車載器発話認証 (Sc通知のLIDを利用して認証)
            
            // SBOX通信用ヘッダー作成
            tcpclient.SBOX_Header = tcpclient.MakeSBOXHeader("IQ",seq_str,"000040","00008")

            // 車載器発話認証データ情報作成(ETC情報通知結果の取得)
            fname := "./Sc_LID.bin"  // Sc通知によりLIDが取得されている
            lid := read_lid_etcinfo(fname)
            tcpclient.SBOXData_SPF = tcpclient.MakeSBOXData_SPFInfo(lid)

            // 要求ヘッダーとデータを結合し、一本のバイナリデータに
            binary.Write(send_byte, binary.BigEndian, &tcpclient.SBOX_Header)
            binary.Write(send_byte, binary.BigEndian, &tcpclient.SBOXData_SPF)
            
        case "Ic": // ETC情報取得

            // SBOX通信用ヘッダー作成
            tcpclient.SBOX_Header = tcpclient.MakeSBOXHeader("Ic",seq_str,"000036","00004")

            // ETC情報取得データ情報作成
            tcpclient.SBOXData_ETCInfo = tcpclient.MakeSBOXData_ETCInfo(0x01,0x43)

            // 要求ヘッダーとデータを結合し、一本のバイナリデータに
            binary.Write(send_byte, binary.BigEndian, &tcpclient.SBOX_Header)
            binary.Write(send_byte, binary.BigEndian, &tcpclient.SBOXData_ETCInfo)

        case "IRONYES":         // 無線制御：開始：車載器指示有り

            // SBOX通信用ヘッダー作成
            tcpclient.SBOX_Header = tcpclient.MakeSBOXHeader("IR",seq_str,"000038","00006")
            tcpclient.SBOXData_RadioCtrl = tcpclient.MakeSBOXData_RadioCtrl(0x01,0x01,0x42) // 制御開始,指示あり, 課金なし、通行可、情報コード1
            
            // 要求ヘッダーとデータを結合し、一本のバイナリデータに
            binary.Write(send_byte, binary.BigEndian, &tcpclient.SBOX_Header)
            binary.Write(send_byte, binary.BigEndian, &tcpclient.SBOXData_RadioCtrl)

        case "IRONNO":          // 無線制御：開始：車載器指示無し

            // SBOX通信用ヘッダー作成
            tcpclient.SBOX_Header = tcpclient.MakeSBOXHeader("IR",seq_str,"000038","00006")
            tcpclient.SBOXData_RadioCtrl = tcpclient.MakeSBOXData_RadioCtrl(0x01,0x00,0x42) // 制御開始,指示無し, 課金なし、通行可、情報コード1
            
            // 要求ヘッダーとデータを結合し、一本のバイナリデータに
            binary.Write(send_byte, binary.BigEndian, &tcpclient.SBOX_Header)
            binary.Write(send_byte, binary.BigEndian, &tcpclient.SBOXData_RadioCtrl)


        case "IROFFYES":        // 無線制御：停止：車載器指示有り
            // SBOX通信用ヘッダー作成
            tcpclient.SBOX_Header = tcpclient.MakeSBOXHeader("IR",seq_str,"000038","00006")
            tcpclient.SBOXData_RadioCtrl = tcpclient.MakeSBOXData_RadioCtrl(0x00,0x01,0x00) // 制御停止, 指示有り, 課金なし、通行可、情報コード0
            
            // 要求ヘッダーとデータを結合し、一本のバイナリデータに
            binary.Write(send_byte, binary.BigEndian, &tcpclient.SBOX_Header)
            binary.Write(send_byte, binary.BigEndian, &tcpclient.SBOXData_RadioCtrl)

        case "IROFFNO":        // 無線制御：停止：車載器指示無し
            // SBOX通信用ヘッダー作成
            tcpclient.SBOX_Header = tcpclient.MakeSBOXHeader("IR",seq_str,"000038","00006")
            tcpclient.SBOXData_RadioCtrl = tcpclient.MakeSBOXData_RadioCtrl(0x00,0x00,0x00) // 制御停止, 指示なし, 課金なし、通行可、情報コード0
            
            // 要求ヘッダーとデータを結合し、一本のバイナリデータに
            binary.Write(send_byte, binary.BigEndian, &tcpclient.SBOX_Header)
            binary.Write(send_byte, binary.BigEndian, &tcpclient.SBOXData_RadioCtrl)

        default:
            flg_err = true      // コマンド異常
            fmt.Printf("MCommandError:%s\n",mcmd)
        }

    case "RSU":
        // 要求ヘッダー情報作成
        tcpclient.RSUHeader = tcpclient.MakeRSUHeader(seq_no, machine_no)

        switch mcmd {

        case "IH":  // 死活監視

            // 死活監視データ情報作成
            tcpclient.RSUDataStandard = tcpclient.MakeRSUData_DeadOrAlive()

            // 要求ヘッダーとデータを結合し、一本のバイナリデータに
            binary.Write(send_byte, binary.BigEndian, &tcpclient.RSUHeader)
            binary.Write(send_byte, binary.BigEndian, &tcpclient.RSUDataStandard)


        case "IT": // 時刻校正

            now := time.Now()
            year_val, month_val, day_val := now.Date()   // 年月日を数字で取得してみる

            tcpclient.RSUDataTime = tcpclient.MakeRSUData_TimeCalibration(uint64(year_val),uint64(month_val),uint64(day_val),uint64(now.Hour()),uint64(now.Minute()),uint64(now.Second()))

            // 要求ヘッダーとデータを結合し、一本のバイナリデータに
            binary.Write(send_byte, binary.BigEndian, &tcpclient.RSUHeader)
            binary.Write(send_byte, binary.BigEndian, &tcpclient.RSUDataTime)
            log.Printf("IT: %v\n",tcpclient.RSUDataTime)

            
        case "IA": // ASK切替
            // ASK要求データ情報作成
            tcpclient.RSUDataStandard = tcpclient.MakeRSUData_ASK()

            // 要求ヘッダーとデータを結合し、一本のバイナリデータに
            binary.Write(send_byte, binary.BigEndian, &tcpclient.RSUHeader)
            binary.Write(send_byte, binary.BigEndian, &tcpclient.RSUDataStandard)

        case "IQ": // QPSK切替
            // QPSK要求データ情報作成
            tcpclient.RSUDataStandard = tcpclient.MakeRSUData_QPSK()

            // 要求ヘッダーとデータを結合し、一本のバイナリデータに
            binary.Write(send_byte, binary.BigEndian, &tcpclient.RSUHeader)
            binary.Write(send_byte, binary.BigEndian, &tcpclient.RSUDataStandard)

        case "IC": // リンク接続
            // リンク接続要求データ情報作成
            tcpclient.RSUDataStandard = tcpclient.MakeRSUData_Link()

            // 要求ヘッダーとデータを結合し、一本のバイナリデータに
            binary.Write(send_byte, binary.BigEndian, &tcpclient.RSUHeader)
            binary.Write(send_byte, binary.BigEndian, &tcpclient.RSUDataStandard)

        case "IIG": // VICS要求 画像
            // VICS要求データ情報作成
            tcpclient.RSUDataVICS = tcpclient.MakeRSUData_VICS(1585)                  // 結果コード以降の総データバイト数
            tcpclient.RSUDataVICS_Data = tcpclient.MakeRSUData_VICSData(1575)         // 大区分データ中の「内容」のデータバイト数
            tcpclient.RSUDataVICS_DataGazo = tcpclient.MakeRSUData_VICSData_Gazo()    // 大区分データ中の内容
            
            // 要求ヘッダーとデータを結合し、一本のバイナリデータに
            binary.Write(send_byte, binary.BigEndian, &tcpclient.RSUHeader)
            binary.Write(send_byte, binary.BigEndian, &tcpclient.RSUDataVICS)
            binary.Write(send_byte, binary.BigEndian, &tcpclient.RSUDataVICS_Data)
            binary.Write(send_byte, binary.BigEndian, &tcpclient.RSUDataVICS_DataGazo)

        case "IIM": // VICS要求 文字
            // VICS要求データ情報作成
            tcpclient.RSUDataVICS = tcpclient.MakeRSUData_VICS(72)                  // 結果コード以降の総データバイト数
            tcpclient.RSUDataVICS_Data = tcpclient.MakeRSUData_VICSData(62)         // 大区分データ中の「内容」のデータバイト数
            tcpclient.RSUDataVICS_DataMoji = tcpclient.MakeRSUData_VICSData_Moji()
            
            // 要求ヘッダーとデータを結合し、一本のバイナリデータに
            binary.Write(send_byte, binary.BigEndian, &tcpclient.RSUHeader)
            binary.Write(send_byte, binary.BigEndian, &tcpclient.RSUDataVICS)
            binary.Write(send_byte, binary.BigEndian, &tcpclient.RSUDataVICS_Data)
            binary.Write(send_byte, binary.BigEndian, &tcpclient.RSUDataVICS_DataMoji)

        case "IIO": // VICS要求 音声
            // VICS要求データ情報作成
            tcpclient.RSUDataVICS = tcpclient.MakeRSUData_VICS(64)                  // 結果コード以降の総データバイト数
            tcpclient.RSUDataVICS_Data = tcpclient.MakeRSUData_VICSData(54)         // 大区分データ中の「内容」のデータバイト数
            tcpclient.RSUDataVICS_DataOnsei = tcpclient.MakeRSUData_VICSData_Onsei()
            
            // 要求ヘッダーとデータを結合し、一本のバイナリデータに
            binary.Write(send_byte, binary.BigEndian, &tcpclient.RSUHeader)
            binary.Write(send_byte, binary.BigEndian, &tcpclient.RSUDataVICS)
            binary.Write(send_byte, binary.BigEndian, &tcpclient.RSUDataVICS_Data)
            binary.Write(send_byte, binary.BigEndian, &tcpclient.RSUDataVICS_DataOnsei)
            
        case "ID": // 電波停止
            // 電波停止データ情報作成
            tcpclient.RSUDataStandard = tcpclient.MakeRSUData_WSTOP()

            // 要求ヘッダーとデータを結合し、一本のバイナリデータに
            binary.Write(send_byte, binary.BigEndian, &tcpclient.RSUHeader)
            binary.Write(send_byte, binary.BigEndian, &tcpclient.RSUDataStandard)
            
        default:
            flg_err = true      // コマンド異常
            fmt.Printf("MCommandError:%s\n",mcmd)
        }
    case "LOG":
        fmt.Printf("LogSystem WakeUp\n")
        log.Printf("LogSystem WakeUp\n")        
    default:
        flg_err = true          // 機種選択異常
        fmt.Printf("MachineError:%s\n",machine)
        log.Printf("MachineError:%s\n",machine)        
    }

    return send_byte, flg_err
}

