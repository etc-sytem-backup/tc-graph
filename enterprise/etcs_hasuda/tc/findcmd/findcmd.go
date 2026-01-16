/**/
package findcmd

import (

    "time"
    "os"
)

/* コマンド確認チャネル用構造体  */
type Req_cmd_st struct {
    Machine string
    Command string
}
var Req_cmd Req_cmd_st
var Ch_req_cmd chan Req_cmd_st

func init() {
    Ch_req_cmd = make(chan Req_cmd_st, 100) // コマンド100個格納できるチャネルとして作成
}

/* コマンド発信指示ファイルの有無をチェック
   Input
   cmd : 発信コマンド名(ファイル名)

   Output
   error  : nil 存在する。  error 存在しない。
   string : 発信コマンド名(ファイル名)
*/
func chk_cmd_filename(cmd string) (string, error) {
    // コマンド発信指示ファイルが存在するか常にチェック。
    // ファイルが存在した場合、そのファイルに対応するコマンドをチャネルにセット
    fp, err := os.OpenFile(cmd, os.O_RDONLY, 0)
    if err != nil {
        if os.IsNotExist(err) {
            return cmd, err // ファイルが存在しない
        }
        return cmd, err // それ以外のエラー(例えばパーミッションがない)
    }
    defer fp.Close()
    
    // ファイルが正しく読み込める(つまり存在する)
    return cmd, nil
}

func find_cmd() {
    interval := time.NewTicker(300 * time.Millisecond) // 300msec間隔
    for {
        select {
        case <-interval.C:  // 300msec毎にコマンドチェック(コマンド名ファイルが存在するか否か)

            // =========
            // RSU 要求
            // =========
            _, err := chk_cmd_filename("RIH") // 死活監視
            if err == nil {
                Req_cmd.Machine = "RSU"
                Req_cmd.Command = "IH"
                Ch_req_cmd <- Req_cmd
            }

            _, err = chk_cmd_filename("RIT") // 状態通知
            if err == nil {
                Req_cmd.Machine = "RSU"
                Req_cmd.Command = "IT"
                Ch_req_cmd <- Req_cmd
            }

            _, err = chk_cmd_filename("RIA") // ASK切替
            if err == nil {
                Req_cmd.Machine = "RSU"
                Req_cmd.Command = "IA"
                Ch_req_cmd <- Req_cmd
            }

            _, err = chk_cmd_filename("RIQ") // QPSK切替
            if err == nil {
                Req_cmd.Machine = "RSU"
                Req_cmd.Command = "IQ"
                Ch_req_cmd <- Req_cmd
            }

            _, err = chk_cmd_filename("RIC") // リンク接続
            if err == nil {
                Req_cmd.Machine = "RSU"
                Req_cmd.Command = "IC"
                Ch_req_cmd <- Req_cmd
            }

            _, err = chk_cmd_filename("RIIG") // VICS情報（画像） ← 2022/05  廃止
            if err == nil {
                Req_cmd.Machine = "RSU"
                Req_cmd.Command = "IIG"
                Ch_req_cmd <- Req_cmd
            }

            _, err = chk_cmd_filename("RIIM") // VICS情報（文字） ← 2022/05  廃止
            if err == nil {
                Req_cmd.Machine = "RSU"
                Req_cmd.Command = "IIM"
                Ch_req_cmd <- Req_cmd
            }

            _, err = chk_cmd_filename("RIIO") // VICS情報（音声） ← 2022/05  廃止
            if err == nil {
                Req_cmd.Machine = "RSU"
                Req_cmd.Command = "IIO"
                Ch_req_cmd <- Req_cmd
            }
            
            _, err = chk_cmd_filename("RID") // 電波停止
            if err == nil {
                Req_cmd.Machine = "RSU"
                Req_cmd.Command = "ID"
                Ch_req_cmd <- Req_cmd
                //                _ = os.Remove("RID")
            }

            // ==========
            // SBOX 要求
            // ==========
            _, err = chk_cmd_filename("SIH")          // 死活監視
            if err == nil {
                Req_cmd.Machine = "SBOX"
                Req_cmd.Command = "IH"
                Ch_req_cmd <- Req_cmd
            }

            _, err = chk_cmd_filename("SIT")          // 状態通知
            if err == nil {
                Req_cmd.Machine = "SBOX"
                Req_cmd.Command = "IT"
                Ch_req_cmd <- Req_cmd
            }

            _, err = chk_cmd_filename("SIc")          // ETC情報取得 ← 2022/8/15 廃止
            if err == nil {
                Req_cmd.Machine = "SBOX"
                Req_cmd.Command = "Ic"
                Ch_req_cmd <- Req_cmd
            }

            _, err = chk_cmd_filename("SIQ")          // 車載器発話認証（SPF認証） RSUリンク接続LID利用
            if err == nil {
                Req_cmd.Machine = "SBOX"
                Req_cmd.Command = "IQ"
                Ch_req_cmd <- Req_cmd
            }

            _, err = chk_cmd_filename("SIQS")          // 車載器発話認証（SPF認証） Sc通知LID利用
            if err == nil {
                Req_cmd.Machine = "SBOX"
                Req_cmd.Command = "IQS"
                Ch_req_cmd <- Req_cmd
            }

            _, err = chk_cmd_filename("SIRONYES")     // 無線制御：開始：車載器指示有り
            if err == nil {
                Req_cmd.Machine = "SBOX"
                Req_cmd.Command = "IRONYES"
                Ch_req_cmd <- Req_cmd
            }

            _, err = chk_cmd_filename("SIRONNO")      // 無線制御：開始：車載器指示無し
            if err == nil {
                Req_cmd.Machine = "SBOX"
                Req_cmd.Command = "IRONNO"
                Ch_req_cmd <- Req_cmd
            }

            _, err = chk_cmd_filename("SIROFFYES")    // 無線制御：停止：車載器指示有り
            if err == nil {
                Req_cmd.Machine = "SBOX"
                Req_cmd.Command = "IROFFYES"
                Ch_req_cmd <- Req_cmd
            }

            _, err = chk_cmd_filename("SIROFFNO")     // 無線制御：停止：車載器指示無し
            if err == nil {
                Req_cmd.Machine = "SBOX"
                Req_cmd.Command = "IROFFNO"
                Ch_req_cmd <- Req_cmd
            }

        }
    }
}

func Run() {

    //    var rcmd Req_cmd_st
    
    // コマンドファイル検出コルーチン
    go find_cmd()

    // 延々とカウントアップを実施(無限に処理させてみる)
    //    go func () {
    //        for {
    //            rcmd := <-Ch_req_cmd
    //            fmt.Println(rcmd.Machine)
    //            fmt.Println(rcmd.Command)
    //        }
    //    }()
}

