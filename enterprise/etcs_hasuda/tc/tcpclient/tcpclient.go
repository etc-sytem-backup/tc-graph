/*
   Traffic Counter TCPClient
*/

package tcpclient

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

/// ============== SBOX関連 Start ==============
// S-BOXへのヘッダー情報
type SBOXHeader_st struct {
	densou_time   [16]byte // 伝聞送信時刻
	command_kinds [2]byte  // コマンド種別
	sequence      [3]byte  // シーケンスNo
	total_size    [6]byte  // トータルデータサイズ
	data_size     [5]byte  // データ部サイズ
}
var SBOX_Header SBOXHeader_st

// S-BOX データ部 ETC情報取得要求
type SBOXData_ETCInfo_st struct {
	syasai_siji_flg  [2]byte // 車載器指示フラグ
	syasai_siji_info [2]byte // 車載器指示情報
}
var SBOXData_ETCInfo SBOXData_ETCInfo_st

// S-BOX データ部 車載器発話認証
type SBOXData_SPF_st struct {
	lid_info [8]byte // LID情報
}
var SBOXData_SPF SBOXData_SPF_st

// 無線制御
type SBOXData_RadioCtrl_st struct {
    Ctrl_info        [2]byte // 制御情報(ASCII)
	syasai_siji_flg  [2]byte // 車載器指示フラグ
	syasai_siji_info [2]byte // 車載器指示情報
}
var SBOXData_RadioCtrl SBOXData_RadioCtrl_st

/// ============== SBOX関連 End ==============

/// ============== RSU関連 Start ==============
// RSUへのヘッダー情報
type RSUHeader_st struct {

	/* 0x11,0x00,0x00,0x01：上位サーバ AP システム
	   0x12,0x00,0x00,0x01：上位サーバ ログシステム
	   0x21,0xXX,0xXX,0xXX：路側無線装置（RSU）*/
	send_rsu_num [4]byte // 送信先機器番号
	recv_rsu_num [4]byte // 送信元機器番号

	seq_num uint16 // シーケンス番号 0x0000 ～ 0xFFFF (65535) （要求－応答は同一番号）, 接続ポート毎に管理
	kinds   int8   // 電文種別 0x01:要求、0x02:応答、0x03:通知
	if_info byte   // I/F情報  0x0A:APシステム、0x0B:ログシステム
}

var RSUHeader RSUHeader_st

// RSU 死活監視/ASK切替/QPSK切替/リンク接続/電波停止  要求コマンドデータ部情報
type RSUDataStandard_st struct {
	info_kinds byte    // 電文種別情報
	yobi       [3]byte // 境界用パディング
	cmd_kinds  [2]byte // コマンド種別
	cmd_length int16   // コマンドデータ長
	cmd_data   [2]byte // 結果コード (0x00,0x00で固定)
}

var RSUDataStandard RSUDataStandard_st

// RSU 時刻構成 データ部情報
type RSUDataTime_st struct {
	info_kinds  byte    // 電文種別情報
	yobi        [3]byte // 境界用パディング
	cmd_kinds   [2]byte // コマンド種別
	cmd_length  int16   // コマンドデータ長
	result_code [2]byte // 結果コード
	year_bcd    [2]byte // 年
	month_bcd   [1]byte // 月
	day_bcd     [1]byte // 日
	hour_bcd    [1]byte // 時間
	min_bcd     [1]byte // 分
	sec_bcd     [1]byte // 秒
	yobi_2      [1]byte // 予備

}

var RSUDataTime RSUDataTime_st

// RSU 装置状態 データ部情報
type RSUDataRSUStatus_st struct {
}

var RSUDataRSUStatus RSUDataRSUStatus_st

// RSU リンク切断 データ部情報
type RSUDataLinkBreak_st struct {
}

var RSUDataLinkBreak RSUDataLinkBreak_st

// RSU VICS データ部情報
type RSUDataVICS_st struct {
	info_kinds byte    // 電文種別情報
	yobi       [3]byte // 境界用パディング
	cmd_kinds  [2]byte // コマンド種別
	cmd_length int16   // コマンドデータ長

	cmd_data       [2]byte // 結果コード (0x00,0x00で固定)
	info_flg       [1]byte // 情報登録削除フラグ
	yobi_2         [1]byte // 予備
	daikubun_count [2]byte // 大区分データカウント
	// 大区分データについては別構造体で定義
}

var RSUDataVICS RSUDataVICS_st

// RSU VICS データ部情報（大区分データ）
type RSUDataVICS_Data_st struct {
	data_size   [2]byte // データサイズ(下記、内容のサイズ)
	ippan_yusen [1]byte // 一般/優先
	yobi        [1]byte // 予備
	// 内容は別構造体で定義
}

var RSUDataVICS_Data RSUDataVICS_Data_st

// RSU VICS データ部情報（内容:画像）
type RSUDataVICS_DataGazo_st struct {
	// ヘッダー部(8 byte)
	id         [1]byte // 格納ID番号
	seigyo_flg [1]byte // 制御フラグ
	jyoho_menu [4]byte // 情報メニュー
	jitsu_data [2]byte // 実データ情報量(下記実データ部のバイト数が入る)

	// 実データ部(1567 byte)
	teikyo_date    [2]byte    // 提供時刻（時/分）
	teikyo_data    [1]byte    // 提供位置指定有無/情報提供方位
	douro_syubetsu [1]byte    // 道路種別
	sid            [2]byte    // SID関連
	service_speed  [1]byte    // 道路のサービス速度
	yuko_kyori     [2]byte    // 有効距離
	jyoho_bytes    [2]byte    // 情報バイト数
	mgo_flags      [1]byte    // 文字/画像/音声 有無フラグ
	moji_bytes     [1]byte    // 文字情報バイト数
	moji_datas     [1]byte    // 漢字文字データ(喋らせたい文字列？表示させたい文字列？)JIS、SJISどちらでも対応との事。
	gazo_bytes     [2]byte    // 画像情報バイト数
	gazo_sikibetsu [1]byte    // 画像形式識別フラグ
	gazo_data      [1545]byte // 画像データ(画像ファイル別に固定)
	onsei_syubetsu [1]byte    // 音声情報種別数
	onsei_bytes    [2]byte    // 音声情報バイト数
	go_sikibetsu   [1]byte    // 言語/音声 識別フラグ
	onsei_data     [1]byte    // 音声データ(音声ファイル別に固定)
}

var RSUDataVICS_DataGazo RSUDataVICS_DataGazo_st

// RSU VICS データ部情報（内容:文字）
type RSUDataVICS_DataMoji_st struct {

	// ヘッダー部(8 byte)
	id         [1]byte // 格納ID番号
	seigyo_flg [1]byte // 制御フラグ
	jyoho_menu [4]byte // 情報メニュー
	jitsu_data [2]byte // 実データ情報量(下記実データ部のバイト数が入る) <- 57 byte

	// 実データ部(60)
	teikyo_date    [2]byte  // 提供時刻（時/分）
	teikyo_data    [1]byte  // 提供位置指定有無/情報提供方位
	douro_syubetsu [1]byte  // 道路種別
	sid            [2]byte  // SID関連
	service_speed  [1]byte  // 道路のサービス速度
	yuko_kyori     [2]byte  // 有効距離
	jyoho_bytes    [2]byte  // 情報バイト数
	mgo_flags      [1]byte  // 文字/画像/音声 有無フラグ
	moji_bytes     [1]byte  // 文字情報バイト数 JISで38byte（改行含む）
	moji_datas     [32]byte // 漢字文字データ(逆走してます。停車してください。)
	//    moji_datas     [38]byte              // 漢字文字データ(予約を確認しました。入場してください。)
	//    moji_datas     [56]byte              // 漢字文字データ(地震が発生しました。あわてず左側路肩に停車してください。)
	gazo_bytes     [2]byte // 画像情報バイト数
	gazo_sikibetsu [1]byte // 画像形式識別フラグ
	gazo_data      [1]byte // 画像データ(画像ファイル別に固定)
	onsei_syubetsu [1]byte // 音声情報種別数
	onsei_bytes    [2]byte // 音声情報バイト数
	go_sikibetsu   [1]byte // 言語/音声 識別フラグ
	onsei_data     [1]byte // 音声データ(音声ファイル別に固定)
}

var RSUDataVICS_DataMoji RSUDataVICS_DataMoji_st

// RSU VICS データ部情報（内容:音声）
type RSUDataVICS_DataOnsei_st struct {
	// ヘッダー部(8 byte)
	id         [1]byte // 格納ID番号
	seigyo_flg [1]byte // 制御フラグ
	jyoho_menu [4]byte // 情報メニュー
	jitsu_data [2]byte // 実データ情報量(下記実データ部のバイト数が入る) <- 98 byte

	// 実データ部(102)
	teikyo_date    [2]byte  // 提供時刻（時/分）
	teikyo_data    [1]byte  // 提供位置指定有無/情報提供方位
	douro_syubetsu [1]byte  // 道路種別
	sid            [2]byte  // SID関連
	service_speed  [1]byte  // 道路のサービス速度
	yuko_kyori     [2]byte  // 有効距離
	jyoho_bytes    [2]byte  // 情報バイト数
	mgo_flags      [1]byte  // 文字/画像/音声 有無フラグ
	moji_bytes     [1]byte  // 文字情報バイト数
	moji_datas     [1]byte  // 漢字文字データ(喋らせたい文字列？表示させたい文字列？)JISコードじゃないとだめ。
	gazo_bytes     [2]byte  // 画像情報バイト数
	gazo_sikibetsu [1]byte  // 画像形式識別フラグ
	gazo_data      [1]byte  // 画像データ(画像ファイル別に固定)
	onsei_syubetsu [1]byte  // 音声情報種別数
	onsei_bytes    [2]byte  // 音声情報バイト数(例:F2S6_ｷﾞｬ&'ｸ%ｿｳ/ｼ%ﾃｲﾏｽ%%.)
	go_sikibetsu   [1]byte  // 言語/音声 識別フラグ
	onsei_data     [24]byte // 音声発話指示データ(バイト数 -> F2S6_ｷﾞｬ&'ｸ%ｿｳ/ｼ%ﾃｲﾏｽ%%.)
	//    onsei_data     [25]byte              // 音声発話指示データ(バイト数 -> F2S6_ﾖ'ﾔｸ%ｦ/ｶ*ｸﾆﾝｼ%ﾏｼ%ﾀ%.)
	//    onsei_data     [28]byte              // 音声発話指示データ(バイト数 -> F2S6_ｼﾞ'ｼ%ﾝｶﾞ&/ﾊ'ｯｾｲｼ%ﾏｼ%ﾀ%.)

}

var RSUDataVICS_DataOnsei RSUDataVICS_DataOnsei_st

/// ============== RSU関連 End ==============

// Traffic Counter クライアント構成
type Client struct {
	addr string // アドレス

	Timeout time.Duration // タイムアウト

	// TCP connection
	mu   sync.Mutex // 排他制御用
	conn net.Conn   // 通信コンテキスト
}

// アスキー変換
func make_ascii(bin uint8) (byte, byte) {

	val_high := bin >> 4
	val_low := bin & 0x0F

	if 0 <= val_high && val_high <= 9 {
		val_high = val_high + 0x30
	} else {

		switch val_high {
		case 0x0A:
			val_high = 0x41
		case 0x0B:
			val_high = 0x42
		case 0x0C:
			val_high = 0x43
		case 0x0D:
			val_high = 0x44
		case 0x0E:
			val_high = 0x45
		case 0x0F:
			val_high = 0x46
		}
	}

	if 0 <= val_low && val_low <= 9 {
		val_low = val_low + 0x30
	} else {

		switch val_low {
		case 0x0A:
			val_low = 0x41
		case 0x0B:
			val_low = 0x42
		case 0x0C:
			val_low = 0x43
		case 0x0D:
			val_low = 0x44
		case 0x0E:
			val_low = 0x45
		case 0x0F:
			val_low = 0x46
		}
	}

	return val_high, val_low
}

/*
   新しいクライアントを作成する
*/
func NewClient(addr string) *Client {
	return &Client{addr: addr}
}

/* SBOX向けヘッダー作成 */
func MakeSBOXHeader(cmdkind string, sqnum string, totalsize string, datasize string) SBOXHeader_st {
	const MilliFormat = "2006/01/02 15:04:05.000" // ← このフォーマット文字列は固定らしい。
	now := time.Now()
	nowUTC := now.UTC()

	year_val, month_val, day_val := now.Date() // 年月日を数字で取得してみる
	//    fmt.Printf("%04d年%02d月%02d日 %02d:%02d:%02d\n",year_val,int(month_val),day_val, now.Hour(), now.Minute(), now.Second())
	header_str := fmt.Sprintf("%v%02v%02v%02v%02v%02v%05v", year_val, int(month_val), int(day_val), now.Hour(), now.Minute(), now.Second(), nowUTC.Format(MilliFormat)[20:])

	copy(SBOX_Header.densou_time[:], []byte(header_str[1:])) // 電文送信時刻 “YYYMMDDHHMMSSTTT”※年月日時分秒ミリ秒。 年は西暦下3桁を使用。
	copy(SBOX_Header.command_kinds[:], []byte(cmdkind))      // コマンド種別
	copy(SBOX_Header.sequence[:], []byte(sqnum))             // シーケンス番号
	copy(SBOX_Header.total_size[:], []byte(totalsize))       // トータルデータサイズ
	copy(SBOX_Header.data_size[:], []byte(datasize))         // データ部サイズ

	return SBOX_Header
}

/* SBOX向け ETC情報取得データ部作成
   flag  : 指示フラグ      : 0x00:車載器指示しない  0x01:車載器指示する
   icode : 識別子          : bit7    0:課金なし  1:課金
           通行許可フラグ  : bit6    0:通行可    1:通行不可
           発話指示        : bit0~5  0x04～0x3C
           1, 課金あり
           2, 通行できません
           3, ETCを利用できません
           4, まもなく、ETC 料金所です。このままお進みください。
           5, ETC を利用できません。右側車線にお進みください。
           6, まもなく、ETC 料金所です。右側車線にお進みください。
           7, ETC を利用できません。左側車線にお進みください。
           8, まもなく、ETC 料金所です。左側車線にお進みください。
           9, ETC を利用できません。中央車線にお進みください。
           10, まもなく、ETC 料金所です。中央車線にお進みください。
           11, ETC を利用できません。徐行してください。
           12, まもなく、ETC 料金所です。徐行してください。

   ※OKI電気仕様書では、識別子がb0で、b2〜b7が情報コードになっているが、実際にセットするデータはその逆転になっている。
     説明はないが、OKI電気仕様書の例がそうなっている。
*/
func MakeSBOXData_ETCInfo(sflag byte, icode byte) SBOXData_ETCInfo_st {

	sflag = sflag + 0x30                                    // フラグをアスキーコード化
	SBOXData_ETCInfo.syasai_siji_flg = [2]byte{0x30, sflag} // 0x30:車載器指示しない、0x31:車載器指示する

	/*
	   bit 0  : 識別子　　　　 0:課金なし、1:課金あり
	   bit 1  : 通行可否フラグ 0:通行可、　1:通行不可
	   bit 2~7: INT(0..63)：値に応じてメッセージを表示
	*/
	//    high_byte := (icode >> 4) + 0x30
	//	low_byte := (icode & 0x0F) + 0x30
	val1, val2 := make_ascii(icode)
	SBOXData_ETCInfo.syasai_siji_info = [2]byte{val1, val2} // 車載器指示情報

	return SBOXData_ETCInfo
}

/* SBOX向け 電波制御
   ctrl_flg : 制御              → 停止 0x00 開始 0x01  引数はASCIIではなくバイナリなので注意！
   sflag    : 車載器指示フラグ  → 無し 0x00 有り 0x01  引数はASCIIではなくバイナリなので注意！
   icode    : 車載器指示情報
              b0 〜 b5 → 0〜63の整数(発話指示)
              b6 0:通行可、1:通行不可
              b7 0:課金なし、1:課金あり

              ★OKI仕様書記載の例１
              課金なし、通行不可、情報コード２
              0x42 → 0x34,0x32 ASCII化

              ★OKI仕様書記載の例２
              課金あり、通行可、情報コード12
              0x8c → 0x38,0x63 ASCII化

           ■発話指示
           1, 課金あり
           2, 通行できません
           3, ETCを利用できません
           4, まもなく、ETC 料金所です。このままお進みください。
           5, ETC を利用できません。右側車線にお進みください。
           6, まもなく、ETC 料金所です。右側車線にお進みください。
           7, ETC を利用できません。左側車線にお進みください。
           8, まもなく、ETC 料金所です。左側車線にお進みください。
           9, ETC を利用できません。中央車線にお進みください。
           10, まもなく、ETC 料金所です。中央車線にお進みください。
           11, ETC を利用できません。徐行してください。
           12, まもなく、ETC 料金所です。徐行してください。

  */
func MakeSBOXData_RadioCtrl(ctrl_flg byte, sflag byte, icode byte) SBOXData_RadioCtrl_st {

    //    copy(SBOXData_RadioCtrl.Ctrl_info[:],[]byte(ctrl_flg))  // 制御情報
    ctrl_flg = ctrl_flg + 0x30                              // 制御を２バイト幅でアスキーコード化
    SBOXData_RadioCtrl.Ctrl_info = [2]byte{0x30,ctrl_flg}   // 0x30,0x30:停止、0x30,0x31:開始

	sflag = sflag + 0x30                                    // 車載器指示フラグを２バイト幅でアスキーコード化
	SBOXData_RadioCtrl.syasai_siji_flg = [2]byte{0x30, sflag} // 0x30,0x30:車載器指示しない、0x30,0x31:車載器指示する

	/* 車載器指示情報のアスキーコード化
	   bit 0  : 識別子         0:課金なし、1:課金あり
	   bit 1  : 通行可否フラグ 0:通行可、  1:通行不可
	   bit 2~7: INT(0..63)：値に応じてメッセージを表示
	*/
	//  high_byte := (icode >> 4) + 0x30
	//	low_byte := (icode & 0x0F) + 0x30
	val1, val2 := make_ascii(icode)
	SBOXData_RadioCtrl.syasai_siji_info = [2]byte{val1, val2} // 車載器指示情報
    
    return SBOXData_RadioCtrl
}

/* SBOX向け 車載器発話認証
   spf_info  : LID情報(ASCII文字)
*/
func MakeSBOXData_SPFInfo(spf_info string) SBOXData_SPF_st {
	copy(SBOXData_SPF.lid_info[:], []byte(spf_info)) // LID情報（ASCII文字列）

	return SBOXData_SPF
}

/* RSU向けヘッダー作成 */
func MakeRSUHeader(seq_num uint16, machine_no int) RSUHeader_st {
	RSUHeader.send_rsu_num = [4]byte{0x21, 0x00, 0x00, byte(machine_no)} // 送信先機器番号 (初期値:RSUへ)
	RSUHeader.recv_rsu_num = [4]byte{0x11, 0x00, 0x00, byte(machine_no)} // 送信元機器番号 (初期値:APサーバーより)
	RSUHeader.seq_num = seq_num                                          // シーケンス番号
	RSUHeader.kinds = 0x01                                               // 電文種別 (初期値:要求)
	RSUHeader.if_info = 0x0A                                             // I/F情報（初期値:APサーバー）
	return RSUHeader
}

/* RSU向け死活監視データ部作成 */
func MakeRSUData_DeadOrAlive() RSUDataStandard_st {
	RSUDataStandard.info_kinds = 0x01                // 電文種別情報：共通電文
	RSUDataStandard.yobi = [3]byte{0x00, 0x00, 0x00} // 予備
	RSUDataStandard.cmd_kinds = [2]byte{0x49, 0x48}  // コマンド種別：IH
	RSUDataStandard.cmd_length = 2                   // コマンドデータ長：2バイト
	RSUDataStandard.cmd_data = [2]byte{0x00, 0x00}   // 結果コード：0x00,0x00固定
	return RSUDataStandard
}

/* RSU向け時刻校正データ部作成
   settime : 校正したい時間文字列 --> 2022
*/
func MakeRSUData_TimeCalibration(year uint64, month uint64, day uint64, hour uint64, min uint64, sec uint64) RSUDataTime_st {

	RSUDataTime.info_kinds = 0x01                // 電文種別情報：共通電文
	RSUDataTime.yobi = [3]byte{0x00, 0x00, 0x00} // 予備
	RSUDataTime.cmd_kinds = [2]byte{0x49, 0x54}  // コマンド種別：IT
	RSUDataTime.cmd_length = 10                  // コマンドデータ長：10バイト

	RSUDataTime.result_code = [2]byte{0x00, 0x00}   // 結果コード
	copy(RSUDataTime.year_bcd[:], ToBcd(year, 2))   // 年
	copy(RSUDataTime.month_bcd[:], ToBcd(month, 1)) // 月
	copy(RSUDataTime.day_bcd[:], ToBcd(day, 1))     // 日
	copy(RSUDataTime.hour_bcd[:], ToBcd(hour, 1))   // 時
	copy(RSUDataTime.min_bcd[:], ToBcd(min, 1))     // 分
	copy(RSUDataTime.sec_bcd[:], ToBcd(sec, 1))     // 秒
	RSUDataTime.yobi_2 = [1]byte{0x00}              // 予備

	//    fmt.Printf("BCD_Year:0x%X 0x%X\n",RSUDataTime.year_bcd[1],RSUDataTime.year_bcd[0] )

	return RSUDataTime
}

/* RSU向けASK切替データ部作成*/
func MakeRSUData_ASK() RSUDataStandard_st {
	RSUDataStandard.info_kinds = 0x12                // 電文種別情報：QPSK電文
	RSUDataStandard.yobi = [3]byte{0x00, 0x00, 0x00} // 予備
	RSUDataStandard.cmd_kinds = [2]byte{0x49, 0x41}  // コマンド種別：IA
	RSUDataStandard.cmd_length = 2                   // コマンドデータ長：2バイト
	RSUDataStandard.cmd_data = [2]byte{0x00, 0x00}   // 結果コード：0x00,0x00固定
	return RSUDataStandard
}

/* RSU向けQPSK切替データ部作成  */
func MakeRSUData_QPSK() RSUDataStandard_st {
	RSUDataStandard.info_kinds = 0x11                // 電文種別情報：ASK電文
	RSUDataStandard.yobi = [3]byte{0x00, 0x00, 0x00} // 予備
	RSUDataStandard.cmd_kinds = [2]byte{0x49, 0x51}  // コマンド種別：IQ
	RSUDataStandard.cmd_length = 2                   // コマンドデータ長：2バイト
	RSUDataStandard.cmd_data = [2]byte{0x00, 0x00}   // 結果コード：0x00,0x00固定
	return RSUDataStandard
}

/* RSU向けリンク接続データ部作成  */
func MakeRSUData_Link() RSUDataStandard_st {
	RSUDataStandard.info_kinds = 0x12                // 電文種別情報：QPSK電文
	RSUDataStandard.yobi = [3]byte{0x00, 0x00, 0x00} // 予備
	RSUDataStandard.cmd_kinds = [2]byte{0x49, 0x43}  // コマンド種別：IC
	RSUDataStandard.cmd_length = 2                   // コマンドデータ長：2バイト
	RSUDataStandard.cmd_data = [2]byte{0x00, 0x00}   // 結果コード：0x00,0x00固定
	return RSUDataStandard
}

/* RSU向けVICS要求データ部作成 （前半） */
func MakeRSUData_VICS(cmd_length int16) RSUDataVICS_st {

	RSUDataVICS.info_kinds = 0x12                // 電文種別情報：QPSK電文
	RSUDataVICS.yobi = [3]byte{0x00, 0x00, 0x00} // 予備
	RSUDataVICS.cmd_kinds = [2]byte{0x49, 0x49}  // コマンド種別：II
	RSUDataVICS.cmd_length = cmd_length          // コマンドデータ長：6 + ? = ??バイト  結果コード以下のデータ長を入れる。

	RSUDataVICS.cmd_data = [2]byte{0x00, 0x00}       // 結果コード：0x00,0x00固定
	RSUDataVICS.info_flg = [1]byte{0x00}             // 情報登録削除フラグ 0x00:登録  0x01:削除
	RSUDataVICS.yobi_2 = [1]byte{0x00}               // 予備
	RSUDataVICS.daikubun_count = [2]byte{0x00, 0x01} // 大区分データカウント
	return RSUDataVICS
}

/* RSU向けVICS要求データ部作成 （後半） */
func MakeRSUData_VICSData(data_length uint16) RSUDataVICS_Data_st {

	// データ内容のバイト数をバイナリ変換(BigEndian)
	len := uint16(data_length)
	len_byte := make([]byte, binary.MaxVarintLen16)
	binary.BigEndian.PutUint16(len_byte, len)
	fmt.Printf("VICS Data Length:%d -> 0x%X\n", len, len_byte)
	log.Printf("VICS Data Length:%d -> 0x%X\n", len, len_byte)

	copy(RSUDataVICS_Data.data_size[:], len_byte) // 「データの内容」のデータサイズ(0~65535 : 2byte)
	RSUDataVICS_Data.ippan_yusen = [1]byte{0x01}  // 一般 / 優先 0x00:一般　0x01:優先
	RSUDataVICS_Data.yobi = [1]byte{0x00}         // 予備
	//    RSUDataVICS_Data.data = [1]byte{0xFF}                  // データの内容(別構造体として定義)

	return RSUDataVICS_Data
}

/* RSU向けVICS要求データ部作成 （内容:画像） */
func MakeRSUData_VICSData_Gazo() RSUDataVICS_DataGazo_st {

	// Total Byte
	// 8 + 1567 = 1575
	// 8 + 1774 = 1782

	// ヘッダー部（8 byte）
	RSUDataVICS_DataGazo.id = [1]byte{0x34}                           // 格納ID番号(52)
	RSUDataVICS_DataGazo.seigyo_flg = [1]byte{0x0C}                   // 制御フラグ
	RSUDataVICS_DataGazo.jyoho_menu = [4]byte{0x00, 0x00, 0x00, 0x80} // 情報メニュー
	RSUDataVICS_DataGazo.jitsu_data = [2]byte{0x06, 0x1F}             // 実データ情報量(下記実データ部バイト数)  1567 byte
	//    RSUDataVICS_DataGazo.jitsu_data  = [2]byte{0x07,0x04}           // 実データ情報量(下記実データ部バイト数)  1796 byte

	//実データ部
	//  17 + 1545 + 5 = 1567 byte
	//  17 + 1774 + 5 = 1796 byte
	RSUDataVICS_DataGazo.teikyo_date = [2]byte{0x1F, 0x3F} // 提供時刻（時/分） 31,63
	RSUDataVICS_DataGazo.teikyo_data = [1]byte{0x00}       // 提供位置指定有無/情報提供方位
	RSUDataVICS_DataGazo.douro_syubetsu = [1]byte{0x00}    // 道路種別
	RSUDataVICS_DataGazo.sid = [2]byte{0x00}               // SID関連
	RSUDataVICS_DataGazo.service_speed = [1]byte{0x00}     // 道路のサービス速度
	RSUDataVICS_DataGazo.yuko_kyori = [2]byte{0x00, 0x00}  // 有効距離
	RSUDataVICS_DataGazo.jyoho_bytes = [2]byte{0x06, 0x14} // 情報バイト数(これ以降？のデータバイト数 : 1556 byte)
	RSUDataVICS_DataGazo.mgo_flags = [1]byte{0x02}         // 文字/画像/音声 有無フラグ   画像情報有無フラグ:1
	RSUDataVICS_DataGazo.moji_bytes = [1]byte{0x00}        // 文字情報バイト数
	RSUDataVICS_DataGazo.moji_datas = [1]byte{0x00}        // 漢字文字データ(喋らせたい文字列？表示させたい文字列？)JIS、SJISどちらでも対応との事。
	RSUDataVICS_DataGazo.gazo_bytes = [2]byte{0x06, 0x09}  // 画像情報バイト数
	RSUDataVICS_DataGazo.gazo_sikibetsu = [1]byte{0x03}    // 画像形式識別フラグ Ping = 3

	// 画像ファイル読み込み（読み込み専用）
	//  16238805-1.png : 1545 byte
	//  16240809-1.png : 1774 byte
	fp, err := os.OpenFile("./png/16238805-1.png", os.O_RDONLY, 755)
	//    fp, err := os.OpenFile("./png/16240809-1.png", os.O_RDONLY,755)
	if err != nil {
		log.Printf("16238805-1.png Open Error.\n")
		//		log.Printf("16240809-1.png Open Error.\n")
		RSUDataVICS_DataGazo.gazo_data = [1545]byte{0x00} // 画像データ(画像ファイル別に固定)(ファイルから読み込むためここではダミーデータセット)
		//        RSUDataVICS_DataGazo.gazo_data      = [1774]byte{0x00}          // 画像データ(画像ファイル別に固定)(ファイルから読み込むためここではダミーデータセット)
	}

	// 読み込んだ画像ファイルをバイトバッファにセット
	var png_img [1545]byte // 画像用バイトバッファ
	//    var png_img [1774]byte          // 画像用バイトバッファ
	errb := binary.Read(fp, binary.BigEndian, &png_img)
	if errb != nil {
		fmt.Println("error occured 'binary.Read()'")
		RSUDataVICS_DataGazo.gazo_data = [1545]byte{0x00} // 画像データ(画像ファイル別に固定)(ファイルから読み込むためここではダミーデータセット)
		//        RSUDataVICS_DataGazo.gazo_data  = [1774]byte{0x00}          // 画像データ(画像ファイル別に固定)(ファイルから読み込むためここではダミーデータセット)
	} else {
		RSUDataVICS_DataGazo.gazo_data = png_img // 画像データ(画像ファイル別に固定)(ファイルから読み込むためここではダミーデータセット)
	}

	// 開いた画像ファイルを閉じる
	if err := fp.Close(); err != nil {
		log.Fatal(err)
	}

	RSUDataVICS_DataGazo.onsei_syubetsu = [1]byte{0x00}    // 音声情報種別数
	RSUDataVICS_DataGazo.onsei_bytes = [2]byte{0x00, 0x00} // 音声情報バイト数
	RSUDataVICS_DataGazo.go_sikibetsu = [1]byte{0x00}      // 言語/音声 識別フラグ
	RSUDataVICS_DataGazo.onsei_data = [1]byte{0x00}        // 音声データ(音声ファイル別に固定)

	return RSUDataVICS_DataGazo
}

/* RSU向けVICS要求データ部作成 （内容:音声） */
func MakeRSUData_VICSData_Onsei() RSUDataVICS_DataOnsei_st {

	// Total Byte
	// F2S6_ｷﾞｬ&'ｸ%ｿｳ/ｼ%ﾃｲﾏｽ%%.     -> 8 + 46 = 54 byte
	// F2S6_ﾖ'ﾔｸ%ｦ/ｶ*ｸ%ﾆﾝｼ%ﾏｼ%ﾀ%.   -> 8 + 47 = 55 byte
	// F2S6_ｼﾞ'ｼ%ﾝｶﾞ&/ﾊ'ｯｾｲｼ%ﾏｼ%ﾀ%. -> 8 + 50 = 58 byte

	// ヘッダー部（8 byte）
	RSUDataVICS_DataOnsei.id = [1]byte{0x34}                           // 格納ID番号(52)
	RSUDataVICS_DataOnsei.seigyo_flg = [1]byte{0x0C}                   // 制御フラグ
	RSUDataVICS_DataOnsei.jyoho_menu = [4]byte{0x80, 0x00, 0x00, 0x00} // 情報メニュー
	RSUDataVICS_DataOnsei.jitsu_data = [2]byte{0x00, 0x2E}             // 実データ情報量(下記実データ部バイト数 22 + 24 = 46) F2S6_ｷﾞｬ&'ｸ%ｿｳ/ｼ%ﾃｲﾏｽ%%.
	//    RSUDataVICS_DataOnsei.jitsu_data  = [2]byte{0x00,0x2F}           // 実データ情報量(下記実データ部バイト数 22 + 25 = 47) F2S6_ﾖ'ﾔｸ%ｦ/ｶ*ｸ%ﾆﾝｼ%ﾏｼ%ﾀ%.
	//    RSUDataVICS_DataOnsei.jitsu_data  = [2]byte{0x00,0x32}           // 実データ情報量(下記実データ部バイト数 22 + 39 = 50) F2S6_ｼﾞ'ｼ%ﾝｶﾞ&/ﾊ'ｯｾｲｼ%ﾏｼ%ﾀ%.

	//実データ部
	//  22 + 24 -> 46 byte F2S6_ｷﾞｬ&'ｸ%ｿｳ/ｼ%ﾃｲﾏｽ%%.
	//  22 + 25 -> 47 byte F2S6_ﾖ'ﾔｸ%ｦ/ｶ*ｸ%ﾆﾝｼ%ﾏｼ%ﾀ%.
	//  22 + 28 -> 50 byte F2S6_ｼﾞ'ｼ%ﾝｶﾞ&/ﾊ'ｯｾｲｼ%ﾏｼ%ﾀ%.
	RSUDataVICS_DataOnsei.teikyo_date = [2]byte{0x1F, 0x3F} // 提供時刻（時/分） 31,63
	RSUDataVICS_DataOnsei.teikyo_data = [1]byte{0x00}       // 提供位置指定有無/情報提供方位
	RSUDataVICS_DataOnsei.douro_syubetsu = [1]byte{0x00}    // 道路種別
	RSUDataVICS_DataOnsei.sid = [2]byte{0x00}               // SID関連
	RSUDataVICS_DataOnsei.service_speed = [1]byte{0x00}     // 道路のサービス速度
	RSUDataVICS_DataOnsei.yuko_kyori = [2]byte{0x00, 0x00}  // 有効距離
	RSUDataVICS_DataOnsei.jyoho_bytes = [2]byte{0x00, 0x23} // 情報バイト数(これ以降？のデータバイト数 11 + 24 = 35byte) F2S6_ｷﾞｬ&'ｸ%ｿｳ/ｼ%ﾃｲﾏｽ%%.
	//    RSUDataVICS_DataOnsei.jyoho_bytes    = [2]byte{0x00,0x24}        // 情報バイト数(これ以降？のデータバイト数 11 + 25 = 36byte) F2S6_ﾖ'ﾔｸ%ｦ/ｶ*ｸ%ﾆﾝｼ%ﾏｼ%ﾀ%.
	//    RSUDataVICS_DataOnsei.jyoho_bytes    = [2]byte{0x00,0x1C}        // 情報バイト数(これ以降？のデータバイト数 11 + 28 = 39 byte) F2S6_ｼﾞ'ｼ%ﾝｶﾞ&/ﾊ'ｯｾｲｼ%ﾏｼ%ﾀ%.

	RSUDataVICS_DataOnsei.mgo_flags = [1]byte{0x01}         // 文字/画像/音声 有無フラグ   音声情報有無フラグ:1
	RSUDataVICS_DataOnsei.moji_bytes = [1]byte{0x00}        // 文字情報バイト数
	RSUDataVICS_DataOnsei.moji_datas = [1]byte{0x00}        // 漢字文字データ(喋らせたい文字列？表示させたい文字列？) JISコード。
	RSUDataVICS_DataOnsei.gazo_bytes = [2]byte{0x00, 0x00}  // 画像情報バイト数
	RSUDataVICS_DataOnsei.gazo_sikibetsu = [1]byte{0x00}    // 画像形式識別フラグ
	RSUDataVICS_DataOnsei.gazo_data = [1]byte{0x00}         // 画像データ
	RSUDataVICS_DataOnsei.onsei_syubetsu = [1]byte{0x01}    // 音声情報種別数
	RSUDataVICS_DataOnsei.onsei_bytes = [2]byte{0x00, 0x18} // 音声情報バイト数 24byte
	//    RSUDataVICS_DataOnsei.onsei_bytes    = [2]byte{0x00,0x19}        // 音声情報バイト数 25byte
	//    RSUDataVICS_DataOnsei.onsei_bytes    = [2]byte{0x00,0x1C}        // 音声情報バイト数 28byte
	RSUDataVICS_DataOnsei.go_sikibetsu = [1]byte{0x00} // 言語/音声 識別フラグ 日本語

	// 音声指示ファイル(内容はアスキーコード)読み込み（読み込み専用）
	//   F2S6_ｷﾞｬ&'ｸ%ｿｳ/ｼ%ﾃｲﾏｽ%%.     -> 24 byte
	//   F2S6_ﾖ'ﾔｸ%ｦ/ｶ*ｸ%ﾆﾝｼ%ﾏｼ%ﾀ%.   -> 25 byte
	//   F2S6_ｼﾞ'ｼ%ﾝｶﾞ&/ﾊ'ｯｾｲｼ%ﾏｼ%ﾀ%. -> 28 byte
	fp, err := os.OpenFile("./onsei/vics_onsei_gyaku.bin", os.O_RDONLY, 755)
	if err != nil {
		log.Printf("vics_onsei_gyaku.bin Open Error.\n")
		RSUDataVICS_DataOnsei.onsei_data = [24]byte{0x00} // 音声データ(音声ファイル別に固定) F2S6_ｷﾞｬ&'ｸ%ｿｳ/ｼ%ﾃｲﾏｽ%%.
		//        RSUDataVICS_DataOnsei.onsei_data     = [25]byte{0x00}            // 音声データ(音声ファイル別に固定) F2S6_ﾖ'ﾔｸ%ｦ/ｶ*ｸ%ﾆﾝｼ%ﾏｼ%ﾀ%.
		//        RSUDataVICS_DataOnsei.onsei_data     = [28]byte{0x00}            // 音声データ(音声ファイル別に固定) F2S6_ｼﾞ'ｼ%ﾝｶﾞ&/ﾊ'ｯｾｲｼ%ﾏｼ%ﾀ%.
	}

	// 読み込んだ音声指示ファイルをバイトバッファにセット
	var onsei [24]byte // 音声指示テキスト用バイトバッファ F2S6_ｷﾞｬ&'ｸ%ｿｳ/ｼ%ﾃｲﾏｽ%%.
	//    var onsei [25]byte          // 音声指示テキスト用バイトバッファ F2S6_ﾖ'ﾔｸ%ｦ/ｶ*ｸ%ﾆﾝｼ%ﾏｼ%ﾀ%.
	//    var onsei [28]byte          // 音声指示テキスト用バイトバッファ F2S6_ｼﾞ'ｼ%ﾝｶﾞ&/ﾊ'ｯｾｲｼ%ﾏｼ%ﾀ%.

	errb := binary.Read(fp, binary.BigEndian, &onsei)
	if errb != nil {
		log.Println("error occured 'binary.Read()'")
		log.Println(errb)
		RSUDataVICS_DataOnsei.onsei_data = [24]byte{0x00} // 音声データ(音声ファイル別に固定) F2S6_ｷﾞｬ&'ｸ%ｿｳ/ｼ%ﾃｲﾏｽ%%.
		//        RSUDataVICS_DataOnsei.onsei_data     = [25]byte{0x00}            // 音声データ(音声ファイル別に固定) F2S6_ﾖ'ﾔｸ%ｦ/ｶ*ｸ%ﾆﾝｼ%ﾏｼ%ﾀ%.
		//        RSUDataVICS_DataOnsei.onsei_data     = [28]byte{0x00}            // 音声データ(音声ファイル別に固定) F2S6_ｼﾞ'ｼ%ﾝｶﾞ&/ﾊ'ｯｾｲｼ%ﾏｼ%ﾀ%.

	} else {
		RSUDataVICS_DataOnsei.onsei_data = onsei // 音声指示テキストファイルデータ
	}

	// 開いた音声指示ファイルを閉じる
	if err := fp.Close(); err != nil {
		log.Fatal(err)
	}

	return RSUDataVICS_DataOnsei
}

/* RSU向けVICS要求データ部作成 （内容:文字） */
func MakeRSUData_VICSData_Moji() RSUDataVICS_DataMoji_st {

	// Total Byte
	// 逆走してます。停車してください。                           -> 8 + 54 = 62 byte
	// 予約を確認しました。入場してください。                    -> 8 + 60 = 68 byte
	// 地震が発生しました。あわてず左側路肩に停車してください。-> 8 + 78 = 86 byte

	// ヘッダー部（8 byte）
	RSUDataVICS_DataMoji.id = [1]byte{0x34}                           // 格納ID番号(52)
	RSUDataVICS_DataMoji.seigyo_flg = [1]byte{0x0C}                   // 制御フラグ
	RSUDataVICS_DataMoji.jyoho_menu = [4]byte{0x08, 0x00, 0x00, 0x00} // 情報メニュー
	RSUDataVICS_DataMoji.jitsu_data = [2]byte{0x00, 0x36}             // 実データ情報量(実データ部バイト数) 逆走してます。停車してください。                            -> 54 byte
	//    RSUDataVICS_DataMoji.jitsu_data  = [2]byte{0x00,0x3C}     // 実データ情報量(実データ部バイト数) 予約を確認しました。入場してください。                     -> 60 byte
	//    RSUDataVICS_DataMoji.jitsu_data  = [2]byte{0x00,0x4E}     // 実データ情報量(実データ部バイト数) 地震が発生しました。あわてず左側路肩に停車してください。 -> 78 byte

	//実データ部
	//  13 + 32 + 9 = 54 逆走してます。停車してください。
	//  13 + 38 + 9 = 60 予約を確認しました。入場してください。
	//  13 + 56 + 9 = 78 地震が発生しました。あわてず左側路肩に停車してください。
	RSUDataVICS_DataMoji.teikyo_date = [2]byte{0x1F, 0x3F} // 提供時刻（時/分） 31,63
	RSUDataVICS_DataMoji.teikyo_data = [1]byte{0x00}       // 提供位置指定有無/情報提供方位
	RSUDataVICS_DataMoji.douro_syubetsu = [1]byte{0x00}    // 道路種別
	RSUDataVICS_DataMoji.sid = [2]byte{0x00}               // SID関連
	RSUDataVICS_DataMoji.service_speed = [1]byte{0x00}     // 道路のサービス速度
	RSUDataVICS_DataMoji.yuko_kyori = [2]byte{0x00, 0x00}  // 有効距離
	RSUDataVICS_DataMoji.jyoho_bytes = [2]byte{0x31}       // 情報バイト数(これ以降？のデータバイト数 49 byte)
	RSUDataVICS_DataMoji.mgo_flags = [1]byte{0x04}         // 文字/画像/音声 有無フラグ   文字情報有無フラグ:1
	RSUDataVICS_DataMoji.moji_bytes = [1]byte{0x50}        // 文字情報バイト数(80 byte)

	// 表示文字テキストファイル読み込み（読み込み専用）
	// 　逆走してます。停車してください。                           -> 32 byte
	// 　予約を確認しました。入場してください。                    -> 38 byte
	// 　地震が発生しました。あわてず左側路肩に停車してください。-> 56 byte
	fp, err := os.OpenFile("./moji/jis_gyakusou.bin", os.O_RDONLY, 755)
	//    fp, err := os.OpenFile("./moji/jis_yoyaku.bin", os.O_RDONLY,755)
	//    fp, err := os.OpenFile("./moji/jis_jishin.bin", os.O_RDONLY,755)
	if err != nil {
		log.Printf("jis_gyakusou.bin Open Error.\n")
		//        log.Printf("jis_yoyaku.bin Open Error.\n")
		//        log.Printf("jis_jishin.bin Open Error.\n")
	}

	// 読み込んだ表示文字テキストファイルをバイトバッファにセット
	var moji [32]byte // 逆走してます。停車してください。
	//    var moji [38]byte          // 予約を確認しました。入場してください。
	//    var moji [56]byte          // 地震が発生しました。あわてず左側路肩に停車してください。
	errb := binary.Read(fp, binary.BigEndian, &moji)
	if errb != nil {
		log.Println("error occured 'binary.Read()'")
		log.Println(errb)
		RSUDataVICS_DataMoji.moji_datas = [32]byte{0x00} // 漢字文字データ(表示させたい文字列) JISコード。
		//        RSUDataVICS_DataMoji.moji_datas     = [38]byte{0x00}            // 漢字文字データ(表示させたい文字列) JISコード。
		//        RSUDataVICS_DataMoji.moji_datas     = [56]byte{0x00}            // 漢字文字データ(表示させたい文字列) JISコード。
	} else {
		RSUDataVICS_DataMoji.moji_datas = moji // 漢字文字データ(表示させたい文字列) JISコード。
	}

	// 開いた音声指示テキストファイルを閉じる
	if err := fp.Close(); err != nil {
		log.Fatal(err)
	}

	RSUDataVICS_DataMoji.gazo_bytes = [2]byte{0x00, 0x00}  // 画像情報バイト数
	RSUDataVICS_DataMoji.gazo_sikibetsu = [1]byte{0x00}    // 画像形式識別フラグ
	RSUDataVICS_DataMoji.gazo_data = [1]byte{0x00}         // 画像データ
	RSUDataVICS_DataMoji.onsei_syubetsu = [1]byte{0x00}    // 音声情報種別数
	RSUDataVICS_DataMoji.onsei_bytes = [2]byte{0x00, 0x00} // 音声情報バイト数
	RSUDataVICS_DataMoji.go_sikibetsu = [1]byte{0x00}      // 言語/音声 識別フラグ
	RSUDataVICS_DataMoji.onsei_data = [1]byte{0x00}        // 音声データ(音声ファイル別に固定)

	return RSUDataVICS_DataMoji
}

/* RSU向け電波停止データ部作成  */
func MakeRSUData_WSTOP() RSUDataStandard_st {
	RSUDataStandard.info_kinds = 0x12                // 電文種別情報：QPSK電文
	RSUDataStandard.yobi = [3]byte{0x00, 0x00, 0x00} // 予備
	RSUDataStandard.cmd_kinds = [2]byte{0x49, 0x44}  // コマンド種別：ID
	RSUDataStandard.cmd_length = 2                   // コマンドデータ長：2バイト
	RSUDataStandard.cmd_data = [2]byte{0x00, 0x00}   // 結果コード：0x00,0x00固定
	return RSUDataStandard
}

/*
   通知・コマンド応答の受信
*/
func (c *Client) Read_command() (string, error) {

	// コマンド応答受信
	buf := make([]byte, 1024)
	rlen, err := c.conn.Read(buf)
	if err != nil {
		_ = c.conn.Close()
		// fmt.Println("TCPDataReadError")
		return "", err
	}

	// 戻り：受信文字列, エラー有無
	return string(buf[:rlen]), nil
}

/*
   S-BOX向け任意コマンド送信
*/
func (c *Client) Send_command(cmd []byte) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.connect(); err != nil {
		fmt.Println("ConnectError")
		return "", err
	}

	if c.Timeout > 0 {
		if err := c.conn.SetDeadline(time.Now().Add(c.Timeout)); err != nil {
			_ = c.conn.Close()
			fmt.Println("TimeOut")
			return "", err
		}
	}

	// コマンドデータ送信
	_, err := c.conn.Write(cmd)
	if err != nil {
		_ = c.conn.Close()
		fmt.Println(err)
		fmt.Println("CommandSendError")
		return "", err
	}

	/*
	   // コマンド応答受信
	   buf := make([]byte, 1024)
	   rlen, err := c.conn.Read(buf)
	   if err != nil {
	       _ = c.conn.Close()
	       fmt.Println("CommandReadError")
	       return "", err
	   }

	   // 戻り：受信文字列, エラー有無
	   return string(buf[:rlen]), nil
	*/

	return "", nil
}

/*
   サーバーへの通信切断
*/
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.close()
}

/*
   サーバーへの接続
*/
func (c *Client) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.connect()
}

/*
   S-BOXへの接続(コネクション)がnilの場合は、接続を試みる。
   すでに接続されている場合は、そのコネクションを維持して使い回す。
*/
func (c *Client) connect() error {
	if c.conn == nil {
		conn, err := net.DialTimeout("tcp", c.addr, c.Timeout)
		if err != nil {
			return err
		}
		c.conn = conn
	}
	return nil
}

// closeはコネクションがnilでない場合にコネクションをCloseします。
/*
   S-BOXへの接続(コネクション)がnil以外ならば、接続をCloseする。
   S-BOXへの接続(コネクション)がnilの場合はCloseしない。
*/
func (c *Client) close() error {
	var err error
	if c.conn != nil {
		err = c.conn.Close()
		c.conn = nil
	}
	return err
}

/* uint64を指定バイトでBCD変換  */
func ToBcd(num uint64, byteCount int) []byte {
	var (
		bcd = make([]byte, byteCount)
	)

	for index := 1; index <= byteCount; index++ {
		mod := num % 100

		digit2 := mod % 10
		digit1 := (mod - digit2) / 10

		bcd[(byteCount - index)] = byte((digit1 * 16) + digit2)

		num = (num - mod) / 100
	}

	return bcd
}

// ToUInt64 -- bcd を uint64 に変換します.
func ToUInt64(bcd []byte) uint64 {
	var (
		result uint64 = 0
	)

	for _, b := range bcd {
		digit1 := b >> 4
		digit2 := b & 0x0f

		result = (result * 100) + uint64(digit1*10) + uint64(digit2)
	}

	return result
}
