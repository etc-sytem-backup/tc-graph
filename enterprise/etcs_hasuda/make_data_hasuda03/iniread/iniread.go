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

    Request_interval int         // 要求間隔(msec)
    Repeat_check_interval int    // 駐車場利用リピート回数。テーブルチェック間隔(sec)
    Parking_duration int         // 警告表示までの駐車場滞在時間

    Large_parking_space int      // 駐車室数（大型車）
    Other_parking_space int      // 駐車室数（大型車以外）

    Duration_time int            // 駐車パスとみなす駐車場利用時間(秒)
    Goback_drive_path_day string // 駐車パス管理テーブル内のデータを直近から何日分残すか指定
    Path_reset_time string       // 駐車パス管理テーブルのリセット時刻

    Entrance_distance int        // アンテナ1からアンテナ2までの距離（メートル）
    Traffic_jam_speed int        // ランプ1停滞中フラグ判定用
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
        Bin_log_path:        cfg.Section("log").Key("bin_log").String(),
        Csv_log_path:        cfg.Section("log").Key("csv_log").String(),
        Run_log_path:        cfg.Section("log").Key("run_log").String(),

        Request_interval:    cfg.Section("num").Key("request_interval").MustInt(),        // 要求発信間隔(msec)
        Repeat_check_interval: cfg.Section("num").Key("repeat_check_interval").MustInt(), // 駐車場利用リピート回数。テーブルチェック間隔(sec)
        Parking_duration:    cfg.Section("num").Key("parking_duration").MustInt(),        // 警告表示までの駐車場滞在時間

        Large_parking_space: cfg.Section("num").Key("large_parking_space").MustInt(),     // 駐車場室数（大型車）
        Other_parking_space: cfg.Section("num").Key("other_parking_space").MustInt(),     // 駐車場室数（大型車以外）

        Duration_time: cfg.Section("num").Key("duration_time").MustInt(),                 // 駐車パスとみなす駐車場利用時間
        Path_reset_time: cfg.Section("date").Key("path_reset_time").String(),             // 駐車パス管理テーブルのリセット時刻
        Goback_drive_path_day: cfg.Section("date").Key("goback_drive_path_day").String(), // 駐車パス管理テーブルのリセット時刻

        Entrance_distance: cfg.Section("num").Key("entrance_distance").MustInt(),         // アンテナ1からアンテナ2までの距離（メートル）
        Traffic_jam_speed: cfg.Section("num").Key("traffic_jam_speed").MustInt(),         // ランプ1停滞中フラグ判定用
    }
}

// inireadのエントリポイント
func Run() {
    LoadConfig()
}
