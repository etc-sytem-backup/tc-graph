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
    Connect_chk_interval int     // RSU回線切断判定用経過時間(分)
    Sc_receive_interval int      // 車両通過通知判定用経過時間(分)
    
    Rsu01_ah string             // 入口ランプアンテナ通信ログ保存ディレクトリ
    Rsu02_ah string             // 駐車場入口アンテナ通信ログ保存ディレクトリ
    Rsu03_ah string             // 駐車場出口アンテナ通信ログ保存ディレクトリ

    Sbox01_sc string            // 入口ランプアンテナSc通知ログ保存ディレクトリ
    Sbox02_sc string            // 駐車場入口アンテナSc通知ログ保存ディレクトリ
    Sbox03_sc string            // 駐車場出口アンテナSc通知ログ保存ディレクトリ
    
    Find_a1 string              // 入口ランプアンテナ通信ログファイル名検索文字
    Find_a2 string              // 駐車場入口アンテナ通信ログファイル名検索文字
    Find_a3 string              // 駐車場出口アンテナ通信ログファイル名検索文字

    Find_car1 string            // 入口ランプアンテナ最新通過通知ファイル検索文字
    Find_car2 string            // 駐車場入口アンテナ最新通過通知ファイル検索文字
    Find_car3 string            // 駐車場出口アンテナ最新通過通知ファイル検索文字

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
        Bin_log_path:         cfg.Section("log").Key("bin_log").String(),               // 送受信データ（バイナリ）保存用ディレクトリパス
        Csv_log_path:         cfg.Section("log").Key("csv_log").String(),               // 送受信データ（CSV形式）保存用ディレクトリパス 
        Run_log_path:         cfg.Section("log").Key("run_log").String(),               // 動作ログ保存用ディレクトリパス                

        Request_interval:     cfg.Section("num").Key("request_interval").MustInt(),     // 要求発信間隔(msec)
        Connect_chk_interval: cfg.Section("num").Key("connect_chk_interval").MustInt(), // RSU回線切断判定用経過時間
        Sc_receive_interval:  cfg.Section("num").Key("sc_receive_interval").MustInt(),  // 車両通過通知判定用経過時間(分)

        Rsu01_ah:             cfg.Section("path").Key("rsu01_ah").String(),             // 入口ランプアンテナ通信ログ保存ディレクトリ
        Rsu02_ah:             cfg.Section("path").Key("rsu02_ah").String(),             // 駐車場入口アンテナ通信ログ保存ディレクトリ
        Rsu03_ah:             cfg.Section("path").Key("rsu03_ah").String(),             // 駐車場出口アンテナ通信ログ保存ディレクトリ

        Sbox01_sc:            cfg.Section("path").Key("sbox01_sc").String(),             // 入口ランプアンテナSc通知ログ保存ディレクトリ
        Sbox02_sc:            cfg.Section("path").Key("sbox02_sc").String(),             // 駐車場入口アンテナSc通知ログ保存ディレクトリ
        Sbox03_sc:            cfg.Section("path").Key("sbox03_sc").String(),             // 駐車場出口アンテナSc通知ログ保存ディレクトリ
        
        Find_a1:              cfg.Section("find").Key("find_a1").String(),              // 入口ランプアンテナ通信ログファイル名検索文字
        Find_a2:              cfg.Section("find").Key("find_a2").String(),              // 駐車場入口アンテナ通信ログファイル名検索文字
        Find_a3:              cfg.Section("find").Key("find_a3").String(),              // 駐車場出口アンテナ通信ログファイル名検索文字

        Find_car1:            cfg.Section("find").Key("find_car1").String(),            // 入口ランプアンテナ通過通知ファイル名検索文字
        Find_car2:            cfg.Section("find").Key("find_car2").String(),            // 駐車場入口アンテナ通過通知ファイル名検索文字
        Find_car3:            cfg.Section("find").Key("find_car3").String(),            // 駐車場出口アンテナ通過通知ファイル名検索文字
        
    }
}

// inireadのエントリポイント
func Run() {
    LoadConfig()
}
