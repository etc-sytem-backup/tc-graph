package iniread

import (
    "log"
	"gopkg.in/ini.v1"
)


// iniファイル取込用構造体
type ConfigList struct {

    // config.iniの設定内容にあわせる

    Bin_log_path string       // 受信データ保存用（バイナリ）
    Csv_log_path string       // 受信データ保存用（）
    Run_log_path string

    A01_tc_csv_path string    // A01のETC情報取得要求応答データ保存先
    A02_tc_csv_path string    // A02のETC情報取得要求応答データ保存先
    A03_tc_csv_path string    // A03のETC情報取得要求応答データ保存先
    A04_tc_csv_path string    // A04のETC情報取得要求応答データ保存先
    A01_tc_table_path string  // A01の通過履歴一覧ファイル保存先     
    A02_tc_table_path string  // A02の通過履歴一覧ファイル保存先     
    A03_tc_table_path string  // A03の通過履歴一覧ファイル保存先     
    A04_tc_table_path string  // A04の通過履歴一覧ファイル保存先     

    A01_tc_wcn_path string    // sbox01の通過WCN番号一覧保存先
    A02_tc_wcn_path string    // sbox02の通過WCN番号一覧保存先
    A03_tc_wcn_path string    // sbox03の通過WCN番号一覧保存先
    A04_tc_wcn_path string    // sbox04の通過WCN番号一覧保存先
    A01_wcn_table_path string // sbox01の通過WCN番号一覧保存先
    A02_wcn_table_path string // sbox02の通過WCN番号一覧保存先
    A03_wcn_table_path string // sbox03の通過WCN番号一覧保存先
    A04_wcn_table_path string // sbox04の通過WCN番号一覧保存先

    Ip_address string         // ME9302のIPアドレス
    Port_num_a1 string        // 受信ポート番号（A1用）
    Port_num_a2 string        // 受信ポート番号（A2用）
    Port_num_a3 string        // 受信ポート番号（A3用）
    Port_num_a4 string        // 受信ポート番号（A4用）
    Timer_interval int        // Goルーチン処理間隔(msec)
    Detection_interval int    // 渋滞による再検出か否かを判断する秒数
}
var Config ConfigList

/*
   iniファイル読込
*/
func LoadConfig() {

    // iniファイルを読み込む
    cfg, err := ini.Load("./config.ini")
    if err != nil {
        log.Fatalln(err)
    }

    // 構造体を初期化する
    Config = ConfigList{

        // iniファイルのデータを読み込む

        Bin_log_path: cfg.Section("log").Key("bin_log").String(),                     // 送受信データ（バイナリ）保存用ディレクトリパス
        Csv_log_path: cfg.Section("log").Key("csv_log").String(),                     // 送受信データ（CSV形式）保存用ディレクトリパス
        Run_log_path: cfg.Section("log").Key("run_log").String(),                     // 動作ログ保存用ディレクトリパス

        A01_tc_csv_path: cfg.Section("tc_csv").Key("a01_tc_csv_path").String(),       // A01のETC情報取得要求応答データ保存先
        A02_tc_csv_path: cfg.Section("tc_csv").Key("a02_tc_csv_path").String(),       // A02のETC情報取得要求応答データ保存先
        A03_tc_csv_path: cfg.Section("tc_csv").Key("a03_tc_csv_path").String(),       // A03のETC情報取得要求応答データ保存先
        A04_tc_csv_path: cfg.Section("tc_csv").Key("a04_tc_csv_path").String(),       // A04のETC情報取得要求応答データ保存先
        A01_tc_table_path: cfg.Section("tc_csv").Key("a01_tc_table_path").String(),   // A01の通過履歴一覧ファイル保存先     
        A02_tc_table_path: cfg.Section("tc_csv").Key("a02_tc_table_path").String(),   // A02の通過履歴一覧ファイル保存先     
        A03_tc_table_path: cfg.Section("tc_csv").Key("a03_tc_table_path").String(),   // A03の通過履歴一覧ファイル保存先     
        A04_tc_table_path: cfg.Section("tc_csv").Key("a04_tc_table_path").String(),   // A04の通過履歴一覧ファイル保存先     


        A01_tc_wcn_path: cfg.Section("tc_wcn").Key("a01_tc_wcn_path").String(),       // sbox01の通過WCN番号一覧保存先
        A02_tc_wcn_path: cfg.Section("tc_wcn").Key("a02_tc_wcn_path").String(),       // sbox02の通過WCN番号一覧保存先
        A03_tc_wcn_path: cfg.Section("tc_wcn").Key("a03_tc_wcn_path").String(),       // sbox03の通過WCN番号一覧保存先
        A04_tc_wcn_path: cfg.Section("tc_wcn").Key("a04_tc_wcn_path").String(),       // sbox04の通過WCN番号一覧保存先
        A01_wcn_table_path: cfg.Section("tc_wcn").Key("a01_wcn_table_path").String(), // sbox01の通過WCN番号一覧保存先
        A02_wcn_table_path: cfg.Section("tc_wcn").Key("a01_wcn_table_path").String(), // sbox02の通過WCN番号一覧保存先
        A03_wcn_table_path: cfg.Section("tc_wcn").Key("a01_wcn_table_path").String(), // sbox03の通過WCN番号一覧保存先
        A04_wcn_table_path: cfg.Section("tc_wcn").Key("a01_wcn_table_path").String(), // sbox04の通過WCN番号一覧保存先

        Ip_address:     cfg.Section("info").Key("ip_address").String(),               // ME9302からデータ受信するPC端末のIPアドレス
        Port_num_a1:    cfg.Section("info").Key("port_num_a1").String(),              // 受信ポート番号（A1用）
        Port_num_a2:    cfg.Section("info").Key("port_num_a2").String(),              // 受信ポート番号（A2用）
        Port_num_a3:    cfg.Section("info").Key("port_num_a3").String(),              // 受信ポート番号（A3用）
        Port_num_a4:    cfg.Section("info").Key("port_num_a4").String(),              // 受信ポート番号（A4用）
        Timer_interval: cfg.Section("num").Key("request_interval").MustInt(),         // 要求発信間隔(msec)
        Detection_interval: cfg.Section("num").Key("detection_interval").MustInt(),   // 渋滞による再検出か否かを判断する秒数
    }
}

/* iniファイル読込実行 */
func Run() {
    LoadConfig()
}
