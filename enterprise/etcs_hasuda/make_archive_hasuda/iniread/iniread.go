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
    Script_start_time string     // スクリプトの実行時間 000000〜235959（00:00:00~23:59:59）
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

        Script_start_time:    cfg.Section("path").Key("script_start_time").String(),
        
        Request_interval:     cfg.Section("num").Key("request_interval").MustInt(),       // 要求発信間隔(msec)

    }
}

// inireadのエントリポイント
func Run() {
    LoadConfig()
}
