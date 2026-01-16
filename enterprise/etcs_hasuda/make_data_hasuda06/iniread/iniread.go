package iniread

import (
    "log"
	"gopkg.in/ini.v1"
)


// iniファイル取込用構造体
type ConfigList struct {

    // config.iniの設定内容にあわせる
    Bin_log_path string          // 送受信データ（バイナリ）保存用ディレクトリパス
    Csv_log_path string          // 送受信データ（CSV形式）保存用ディレクトリパス
    Run_log_path string          // 動作ログ保存用ディレクトリパス

    Ac_csv_path string           // ac直下にあるsboxデータ収集ディレクトリ（通信結果）
    Ac_wcn_path string           // ac直下にあるsboxデータ収集ディレクトリ（アンテナを通過した車両）
    Ac_wcn_table_path string     // ac直下にあるアンテナを通過した車両の一覧ファイル保存先

    Ac_wcn_rireki_a1_file string // アンテナA1の通過履歴ファイル
    Ac_wcn_rireki_a2_file string // アンテナA2の通過履歴ファイル
    Ac_wcn_rireki_a3_file string // アンテナA3の通過履歴ファイル
    Ac_wcn_rireki_a4_file string // アンテナA4の通過履歴ファイル
    Path_reset_time string       // 駐車パス管理テーブルのリセット時刻
    
    Request_interval int         // 要求間隔(msec)

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
        Bin_log_path:         cfg.Section("log").Key("bin_log").String(),                 // 送受信データ（バイナリ）保存用ディレクトリパス
        Csv_log_path:         cfg.Section("log").Key("csv_log").String(),                 // 送受信データ（CSV形式）保存用ディレクトリパス 
        Run_log_path:         cfg.Section("log").Key("run_log").String(),                 // 動作ログ保存用ディレクトリパス                

        Ac_csv_path:         cfg.Section("path").Key("ac_csv_path").String(),             //  ac直下にあるsboxデータ収集ディレクトリ（通信結果）
        Ac_wcn_path:         cfg.Section("path").Key("ac_csv_path").String(),             //  ac直下にあるsboxデータ収集ディレクトリ（アンテナを通過した車両）
        Ac_wcn_table_path:   cfg.Section("path").Key("ac_wcn_table_path").String(),       //  ac直下にあるアンテナを通過した車両の一覧ファイル保存先

        Ac_wcn_rireki_a1_file: cfg.Section("path").Key("ac_wcn_rireki_a1_file").String(), // アンテナA1の通過履歴ファイル
        Ac_wcn_rireki_a2_file: cfg.Section("path").Key("ac_wcn_rireki_a2_file").String(), // アンテナA2の通過履歴ファイル
        Ac_wcn_rireki_a3_file: cfg.Section("path").Key("ac_wcn_rireki_a3_file").String(), // アンテナA3の通過履歴ファイル
        Ac_wcn_rireki_a4_file: cfg.Section("path").Key("ac_wcn_rireki_a4_file").String(), // アンテナA4の通過履歴ファイル
        Path_reset_time: cfg.Section("path").Key("path_reset_time").String(),             // 駐車パス管理テーブルのリセット時刻
        
        Request_interval:     cfg.Section("num").Key("request_interval").MustInt(),       // 要求発信間隔(msec)

    }
}

// inireadのエントリポイント
func Run() {
    LoadConfig()
}
