package iniread

import (
    "log"
	"gopkg.in/ini.v1"
)


// iniファイル取込用構造体
type ConfigList struct {

    Bin_log_path string       // 受信データ保存用（バイナリ）
    Csv_log_path string       // 受信データ保存用（）
    Run_log_path string

    Rsu01_path string         // RSU01
    Rsu02_path string         // RSU02
    Rsu03_path string         // RSU03
    Rsu04_path string         // RSU04
    Sbox01_path string        // SBOX01
    Sbox02_path string        // SBOX01
    Sbox03_path string        // SBOX01
    Sbox04_path string        // SBOX01
    
    Rsu01_csv_path string     // RSU01CSVフォルダパス
    Rsu02_csv_path string     // RSU02CSVフォルダパス
    Rsu03_csv_path string     // RSU03CSVフォルダパス
    Rsu04_csv_path string     // RSU04CSVフォルダパス

    Ac_rsu01_csv_path string  // acCSVフォルダパス
    Ac_rsu02_csv_path string  // acCSVフォルダパス
    Ac_rsu03_csv_path string  // acCSVフォルダパス
    Ac_rsu04_csv_path string  // acCSVフォルダパス

    Rsu01_name string          // RSU01の名前
    Rsu02_name string          // RSU02の名前
    Rsu03_name string          // RSU03の名前
    Rsu04_name string          // RSU04の名前

    Rsu01_alias string         // RSU01の名前変換
    Rsu02_alias string         // RSU02の名前変換
    Rsu03_alias string         // RSU03の名前変換
    Rsu04_alias string         // RSU04の名前変換

    Exitgate_name string       // 出口のRSU
    Exitgate_alias string      // 出口のRSU別名
    
    A1_A2_distance int      // A1-A2間距離(m)
    A2_A3_distance int      // A2-A3間距離(m)
    A3_A4_distance int      // A3-A4間距離(m)
    Request_interval int    // 要求間隔(msec)
    File_del_daycount int   // 断面交通量保存日数(日数)

    Car_reserve_0001 string     // 駐車予約1
    Car_reserve_0002 string     // 駐車予約2
    Car_reserve_0003 string     // 駐車予約3
    Car_reserve_0004 string     // 駐車予約4
    Car_reserve_0005 string     // 駐車予約5
    

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

        Rsu01_path:          cfg.Section("mpath").Key("rsu01_path").String(),     // RSU01
        Rsu02_path:          cfg.Section("mpath").Key("rsu02_path").String(),     // RSU02
        Rsu03_path:          cfg.Section("mpath").Key("rsu03_path").String(),     // RSU03
        Rsu04_path:          cfg.Section("mpath").Key("rsu04_path").String(),     // RSU04
        Sbox01_path:         cfg.Section("mpath").Key("sbox01_path").String(),    // SBOX 58001
        Sbox02_path:         cfg.Section("mpath").Key("sbox02_path").String(),    // SBOX 58002
        Sbox03_path:         cfg.Section("mpath").Key("sbox03_path").String(),    // SBOX 58003
        Sbox04_path:         cfg.Section("mpath").Key("sbox04_path").String(),    // SBOX 58004
        
        Rsu01_csv_path:      cfg.Section("csvpath").Key("rsu01_csv_path").String(),     // RSU01 CSVフォルダパス
        Rsu02_csv_path:      cfg.Section("csvpath").Key("rsu02_csv_path").String(),     // RSU02 CSVフォルダパス
        Rsu03_csv_path:      cfg.Section("csvpath").Key("rsu03_csv_path").String(),     // RSU03 CSVフォルダパス
        Rsu04_csv_path:      cfg.Section("csvpath").Key("rsu04_csv_path").String(),     // RSU04 CSVフォルダパス

        Ac_rsu01_csv_path:   cfg.Section("csvpath").Key("ac_rsu01_csv_path").String(),     // AC CSVフォルダパス RSU01
        Ac_rsu02_csv_path:   cfg.Section("csvpath").Key("ac_rsu02_csv_path").String(),     // AC CSVフォルダパス RSU02
        Ac_rsu03_csv_path:   cfg.Section("csvpath").Key("ac_rsu03_csv_path").String(),     // AC CSVフォルダパス RSU03
        Ac_rsu04_csv_path:   cfg.Section("csvpath").Key("ac_rsu04_csv_path").String(),     // AC CSVフォルダパス RSU04

        Rsu01_name:          cfg.Section("name").Key("rsu01_name").String(),            // RSU01の名前変換
        Rsu02_name:          cfg.Section("name").Key("rsu02_name").String(),            // RSU01の名前変換
        Rsu03_name:          cfg.Section("name").Key("rsu03_name").String(),            // RSU01の名前変換
        Rsu04_name:          cfg.Section("name").Key("rsu04_name").String(),            // RSU01の名前変換

        Rsu01_alias:         cfg.Section("name").Key("rsu01_alias").String(),            // RSU01の名前変換
        Rsu02_alias:         cfg.Section("name").Key("rsu02_alias").String(),            // RSU01の名前変換
        Rsu03_alias:         cfg.Section("name").Key("rsu03_alias").String(),            // RSU01の名前変換
        Rsu04_alias:         cfg.Section("name").Key("rsu04_alias").String(),            // RSU01の名前変換

        Exitgate_name:       cfg.Section("name").Key("exitgate_name").String(),          // 出口RSUの名前
        Exitgate_alias:      cfg.Section("name").Key("exitgate_alias").String(),         // 出口RSUの名前変換
        
        A1_A2_distance:      cfg.Section("num").Key("a1_a2_distance").MustInt(),        // A1-A2間距離(m)
        A2_A3_distance:      cfg.Section("num").Key("a2_a3_distance").MustInt(),        // A2-A3間距離(m)
        A3_A4_distance:      cfg.Section("num").Key("a3_a4_distance").MustInt(),        // A3-A4間距離(m)
        Request_interval:    cfg.Section("num").Key("request_interval").MustInt(),      // 要求発信間隔(msec)
        File_del_daycount:   cfg.Section("num").Key("file_del_daycount").MustInt(),     // ファイル削除日数

        Car_reserve_0001:    cfg.Section("reserve").Key("car_reserve_0001").String(),   // 駐車場予約1
        Car_reserve_0002:    cfg.Section("reserve").Key("car_reserve_0002").String(),   // 駐車場予約2
        Car_reserve_0003:    cfg.Section("reserve").Key("car_reserve_0003").String(),   // 駐車場予約3
        Car_reserve_0004:    cfg.Section("reserve").Key("car_reserve_0004").String(),   // 駐車場予約4
        Car_reserve_0005:    cfg.Section("reserve").Key("car_reserve_0005").String(),   // 駐車場予約5
    }
}

/**/
func Run() {
    LoadConfig()
}
