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

        Timer_interval: cfg.Section("num").Key("request_interval").MustInt(),         // 要求発信間隔(msec)
        Detection_interval: cfg.Section("num").Key("detection_interval").MustInt(),   // 渋滞による再検出か否かを判断する秒数
    }
}

/* iniファイル読込実行 */
func Run() {
    LoadConfig()
}
