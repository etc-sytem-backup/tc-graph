// config.iniからパラメータを読み込む為の定義とメソッド
package iniread

import (
	"log"

	"gopkg.in/ini.v1"
)

// iniファイル取込用構造体
type ConfigList struct {

	// config.iniの設定内容にあわせる

	// 接続先に関する情報
	Mode                  string // 管理者 or 利用者
	User_Name             string // ユーザー名
	Ip_Address            string // 接続先のIPアドレス
	Sftp_Port             string // 情報取得用ポート番号
	Data_Dir              string // 画面に表示される情報元CSVのディレクトリパス
	Disp_Main_Csv         string // メインモニタ
	Disp_Avg_Csv          string // 統計(平均)
	Disp_Table_Week_Csv   string // 統計（一覧 リピート回数（当週））
	Disp_Table_Month_Csv  string // 統計（一覧 リピート回数（当月））
	Disp_Parking_Time_Csv string // 統計（一覧 長時間駐車）
	Disp_Alert_Csv        string // 逆走検知
	Disp_Passage_Csv      string // 断面交通量
	Disp_Settings_Csv     string // 設定

	// 情報取得間隔
	Disp_Main_Interval         string // メインモニタ
	Disp_Avg_Interval          string // 統計(平均)
	Disp_Table_Interval        string // 統計（一覧）
	Disp_Alert_Interval        string // 逆走検知
	Disp_Passage_Interval      string // 断面交通量
	Rsu_connect_reset_interval string // RSU回線切断警告の表示時間

	// コマンドファイル名
	Radio_Start   string // 無線開始
	Radio_Stop    string // 無線停止
	Alert_Reset   string // 逆走警報クリア
	Rsu_Connect   string // RSU回線切断通知
	Passage_Reset string // 断面交通量
	Large_Plus    string // 大型駐車数プラス1補正
	Large_Minus   string // 大型駐車数マイナス1補正
	Small_Plus    string // 小型駐車数プラス1補正
	Small_Minus   string // 小型駐車数マイナス1補正

	// ログ
	Bin_log_path string // 送受信データ（バイナリ）保存用ディレクトリパス
	Csv_log_path string // 送受信データ（CSV形式）保存用ディレクトリパス
	Run_log_path string // 動作ログ保存用ディレクトリパス
}

var Config ConfigList

// iniファイル読込
func LoadConfig() {

	// iniファイルを読み込む
	cfg, err := ini.Load("./config.ini")
	if err != nil {
		log.Fatalln(err)
	}

	// 構造体を初期化する
	Config = ConfigList{

		// iniファイルのデータを読み込む
		Mode:                       cfg.Section("target").Key("mode").String(),                         // 管理者 or 利用者
		User_Name:                  cfg.Section("target").Key("user_name").String(),                    // ユーザー名
		Ip_Address:                 cfg.Section("target").Key("ip_address").String(),                   // 接続先のIPアドレス
		Sftp_Port:                  cfg.Section("target").Key("sftp_port").String(),                    // 情報取得用ポート番号
		Data_Dir:                   cfg.Section("target").Key("data_dir").String(),                     // 画面に表示される情報元CSVのディレクトリパス
		Disp_Main_Csv:              cfg.Section("target").Key("disp_main_csv").String(),                // メインモニタ
		Disp_Avg_Csv:               cfg.Section("target").Key("disp_avg_csv").String(),                 // 統計(平均)
		Disp_Table_Week_Csv:        cfg.Section("target").Key("disp_table_week_csv").String(),          // 統計（一覧 リピート回数（当週））
		Disp_Table_Month_Csv:       cfg.Section("target").Key("disp_table_month_csv").String(),         // 統計（一覧 リピート回数（当月））
		Disp_Parking_Time_Csv:      cfg.Section("target").Key("disp_parking_time_csv").String(),        // 統計（一覧 長時間駐車）
		Disp_Alert_Csv:             cfg.Section("target").Key("disp_alert_csv").String(),               // 逆走検知
		Disp_Passage_Csv:           cfg.Section("target").Key("disp_passage_csv").String(),             // 断面交通量
		Disp_Settings_Csv:          cfg.Section("target").Key("disp_settings_csv").String(),            // 設定
		Disp_Main_Interval:         cfg.Section("interval").Key("disp_main_interval").String(),         // メインモニタ
		Disp_Avg_Interval:          cfg.Section("interval").Key("disp_avg_interval").String(),          // 統計(平均)
		Disp_Table_Interval:        cfg.Section("interval").Key("disp_table_interval").String(),        // 統計（一覧）
		Disp_Alert_Interval:        cfg.Section("interval").Key("disp_alert_interval").String(),        // 逆走検知
		Disp_Passage_Interval:      cfg.Section("interval").Key("disp_passage_interval").String(),      // 断面交通量
		Rsu_connect_reset_interval: cfg.Section("interval").Key("rsu_connect_reset_interval").String(), // RSU回線切断警告表示時間
		Radio_Start:                cfg.Section("command").Key("radio_start").String(),                 // 無線開始
		Radio_Stop:                 cfg.Section("command").Key("radio_stop").String(),                  // 無線停止
		Alert_Reset:                cfg.Section("command").Key("alert_reset").String(),                 // 逆走警報クリア
		Passage_Reset:              cfg.Section("command").Key("passage_reset").String(),               // 断面交通量リセット
		Rsu_Connect:                cfg.Section("command").Key("rsu_connect").String(),                 // RSU回線切断通知
		Large_Plus:                 cfg.Section("command").Key("large_plus").String(),                  // 大型駐車数プラス1補正
		Large_Minus:                cfg.Section("command").Key("large_minus").String(),                 // 大型駐車数マイナス1補正
		Small_Plus:                 cfg.Section("command").Key("small_plus").String(),                  // 小型駐車数プラス1補正
		Small_Minus:                cfg.Section("command").Key("small_minus").String(),                 // 小型駐車数マイナス1補正
		Bin_log_path:               cfg.Section("log").Key("bin_log").String(),                         // 送受信データ（バイナリ）保存用ディレクトリパス
		Csv_log_path:               cfg.Section("log").Key("csv_log").String(),                         // 送受信データ（CSV形式）保存用ディレクトリパス
		Run_log_path:               cfg.Section("log").Key("run_log").String(),                         // 動作ログ保存用ディレクトリパス
	}
}

/* iniファイル読込実行 */
func Run() {
	LoadConfig()
}
