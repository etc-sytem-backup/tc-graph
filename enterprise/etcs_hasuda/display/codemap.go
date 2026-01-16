package main

import "fmt"

/*
支局コードから支局を探し、返す
*/
func searchSikyokuFromCode(code string) (string, error) {
	switch code {
	case "535053": // 札幌 SPS
		return "札幌", nil
	case "535020": // 札   SP
		return "札", nil
	case "484448": // 函館 HDH
		return "函館", nil
	case "484420": // 函   HD
		return "函", nil
	case "414b41": // 旭川 AKA
		return "旭川", nil
	case "414b20": // 旭   AK
		return "旭", nil
	case "4d524d": // 室蘭 MRM
		return "室蘭", nil
	case "4d5220": // 室   MR
		return "室", nil
	case "4b524b": // 釧路 KRK
		return "釧路", nil
	case "4b5220": // 釧   KR
		return "釧", nil
	case "4f484f": // 帯広 OHO
		return "帯広", nil
	case "4f4820": // 帯   OH
		return "帯", nil
	case "4b494b": // 北見 KIK
		return "北見", nil
	case "4b4920": // 北   KI
		return "北", nil
	case "414d41": // 青森 AMA
		return "青森", nil
	case "414d48": // 八戸 AMH
		return "八戸", nil
	case "414d20": // 青   AM
		return "青", nil
	case "495449": // 岩手 ITI
		return "岩手", nil
	case "495420": // 岩   IT
		return "岩", nil
	case "4D4753": // 仙台 MGS
		return "仙台", nil
	case "4d474d": // 宮城 MGM
		return "宮城", nil
	case "4d4720": // 宮   MG
		return "宮", nil
	case "415441": // 秋田 ATA
		return "秋田", nil
	case "415420": // 秋   AT
		return "秋", nil
	case "594120": // 山形 YA
		return "山形", nil
	case "594153": // 庄内 YAS
		return "庄内", nil
	case "465320": // 福島 FS
		return "福島", nil
	case "465341": // 会津 FSA
		return "会津", nil
	case "465349": // いわきFSI
		return "いわき", nil
	case "49474d": // 水戸 IGM
		return "水戸", nil
	case "494754": // 土浦 IGT
		return "土浦", nil
	case "49474b": // つくばIGK
		return "つくば", nil
	case "494749": // 茨城 IGI
		return "茨城", nil
	case "494720": // 茨   IG
		return "茨", nil
	case "544755": // 宇都宮TGU
		return "宇都宮", nil
	case "54474e": // 那須 TGN
		return "那須", nil
	case "544743": // とちぎTGC
		return "とちぎ", nil
	case "544754": // 栃木 TGT
		return "栃木", nil
	case "544720": // 栃   TG
		return "栃", nil
	case "474d47": // 群馬 GMG
		return "群馬", nil
	case "474d54": // 高崎 GMT
		return "高崎", nil
	case "474d20": // 群   GM
		return "群", nil
	case "53544f": // 大宮 STO
		return "大宮", nil
	case "535447": // 川越 STG
		return "川越", nil
	case "535454": // 所沢 STT
		return "所沢", nil
	case "53544b": // 熊谷 STK
		return "熊谷", nil
	case "535442": // 春日部STB
		return "春日部", nil
	case "535453": // 埼玉 STS
		return "埼玉", nil
	case "535420": // 埼   ST
		return "埼", nil
	case "434243": // 千葉 CBC
		return "千葉", nil
	case "434254": // 成田 CBT
		return "成田", nil
	case "43424e": // 習志野CBN
		return "習志野", nil
	case "434253": // 袖ヶ浦CBS
		return "袖ヶ浦", nil
	case "434244": // 野田 CBD
		return "野田", nil
	case "43424b": // 柏   CBK
		return "柏", nil
	case "434220": // 千   CB
		return "千", nil
	case "544b53": // 品川 TKS
		return "品川", nil
	case "544f53": // 品   TOS
		return "品", nil
	case "544b4e": // 練馬 TKN
		return "練馬", nil
	case "544f4e": // 練   TON
		return "練", nil
	case "544b41 ": // 足立 TKA
		return "足立", nil
	case "544f41 ": // 足   TOA
		return "足", nil
	case "544b48": // 八王子TKH
		return "八王子", nil
	case "544b54": // 多摩 TKT
		return "多摩", nil
	case "544f54": // 多   TOT
		return "多", nil
	case "4b4e59": // 横浜 KNY
		return "横浜", nil
	case "4b4e4b": // 川崎 KNK
		return "川崎", nil
	case "4b4e4e": // 湘南 KNN
		return "湘南", nil
	case "4b4e53": // 相模 KNS
		return "相模", nil
	case "4b4e20": // 神   KN
		return "神", nil
	case "594e20": // 山梨 YN
		return "山梨", nil
	case "464a53": // 富士山FJS
		return "富士山", nil
	case "4e474e": // 新潟 NGN
		return "新潟", nil
	case "4e474f": // 長岡 NGO
		return "長岡", nil
	case "4e4720": // 新   NG
		return "新", nil
	case "545954": // 富山 TYT
		return "富山", nil
	case "545920": // 富   TY
		return "富", nil
	case "494b4b": // 金沢 IKK
		return "金沢", nil
	case "494b49": // 石川 IKI
		return "石川", nil
	case "494b20": // 石   IK
		return "石", nil
	case "4e4e4e": // 長野 NNN
		return "長野", nil
	case "4e4e4d": // 松本 NNM
		return "松本", nil
	case "4e4e53": // 諏訪 NNS
		return "諏訪", nil
	case "4e4e20": // 長   NN
		return "長", nil
	case "464920": // 福井 FI
		return "福井", nil
	case "474647": // 岐阜 GFG
		return "岐阜", nil
	case "474648": // 飛騨 GFH
		return "飛騨", nil
	case "474620": // 岐   GF
		return "岐", nil
	case "535a53": // 静岡 SZS
		return "静岡", nil
	case "535a48": // 浜松 SZH
		return "浜松", nil
	case "535a4e": // 沼津 SZN
		return "沼津", nil
	case "535a49": // 伊豆 SZI
		return "伊豆", nil
	case "535a20": // 静   SZ
		return "静", nil
	case "41434e": // 名古屋ACN
		return "名古屋", nil
	case "414354": // 豊橋 ACT
		return "豊橋", nil
	case "41435a": // 岡崎 ACZ
		return "岡崎", nil
	case "41434d": // 三河 ACM
		return "三河", nil
	case "414359": // 豊田 ACY
		return "豊田", nil
	case "414349": // 一宮 ACI
		return "一宮", nil
	case "41434f": // 尾張小ACO牧
		return "尾張小", nil
	case "414320": // 愛   AC
		return "愛", nil
	case "4d454d": // 三重 MEM
		return "三重", nil
	case "4d4553": // 鈴鹿 MES
		return "鈴鹿", nil
	case "4d4520": // 三   ME
		return "三", nil
	case "534953": // 滋賀 SIS
		return "滋賀", nil
	case "534920": // 滋   SI
		return "滋", nil
	case "4b544b": // 京都 KTK
		return "京都", nil
	case "4b5420": // 京   KT
		return "京", nil
	case "4f534f": // 大阪 OSO
		return "大阪", nil
	case "4f534e": // なにわOSN
		return "なにわ", nil
	case "4f5353": // 堺   OSS
		return "堺", nil
	case "4f535a": // 和泉 OSZ
		return "和泉", nil
	case "4f5320": // 大   OS
		return "大", nil
	case "4f5349": // 泉   OSI
		return "泉", nil
	case "48474b": // 神戸 HGK
		return "神戸", nil
	case "484748": // 姫路 HGH
		return "姫路", nil
	case "484720": // 兵   HG
		return "兵", nil
	case "4e524e": // 奈良 NRN
		return "奈良", nil
	case "4e5220": // 奈   NR
		return "奈", nil
	case "574b57": // 和歌山WKW
		return "和歌山", nil
	case "574b20": // 和   WK
		return "和", nil
	case "545454": // 鳥取 TTT
		return "鳥取", nil
	case "545420": // 鳥   TT
		return "鳥", nil
	case "534e20": // 島根 SN
		return "島根", nil
	case "534d20": // 島   SM
		return "島", nil
	case "4f594f": // 岡山 OYO
		return "岡山", nil
	case "4f594b": // 倉敷 OYK
		return "倉敷", nil
	case "4f5920": // 岡   OY
		return "岡", nil
	case "485348": // 広島 HSH
		return "広島", nil
	case "485346": // 福山 HSF
		return "福山", nil
	case "485320": // 広   HS
		return "広", nil
	case "595553": // 下関 YUS
		return "下関", nil
	case "595559": // 山口 YUY
		return "山口", nil
	case "595520": // 山   YU
		return "山", nil
	case "545354": // 徳島 TST
		return "徳島", nil
	case "545320": // 徳   TS
		return "徳", nil
	case "4b414b": // 香川 KAK
		return "香川", nil
	case "4b4120": // 香   KA
		return "香", nil
	case "454820": // 愛媛 EH
		return "愛媛", nil
	case "4b434b": // 高知 KCK
		return "高知", nil
	case "4b4320": // 高   KC
		return "高", nil
	case "464f46": // 福岡 FOF
		return "福岡", nil
	case "464f4b": // 北九州FOK
		return "北九州", nil
	case "464f52": // 久留米FOR
		return "久留米", nil
	case "464f43": // 筑豊 FOC
		return "筑豊", nil
	case "464f20": // 福   FO
		return "福", nil
	case "534153": // 佐賀 SAS
		return "佐賀", nil
	case "534120": // 佐   SA
		return "佐", nil
	case "4e5320": // 長崎 NS
		return "長崎", nil
	case "4e5353": // 佐世保NSS
		return "佐世保", nil
	case "4b554b": // 熊本 KUK
		return "熊本", nil
	case "4b5520": // 熊   KU
		return "熊", nil
	case "4f5420": // 大分 OT
		return "大分", nil
	case "4d5a20": // 宮崎 MZ
		return "宮崎", nil
	case "4b4f4b": // 鹿児島KOK
		return "鹿児島", nil
	case "4b4f20": // 鹿   KO
		return "鹿", nil
	case "4f4e4f": // 沖縄 ONO
		return "沖縄", nil
	case "4f4e20": // 沖   ON
		return "沖", nil
	default:
		return "", fmt.Errorf("Not found Sikyoku from Code\n")
	}
}

/*
用途コードから用途を探し、返す
*/
func searchYoutoFromCode(code string) (string, error) {
	switch code {

	// 自家用
	case "bb": // さ
		return "さ", nil
	case "bd": // す
		return "す", nil
	case "be": // せ
		return "せ", nil
	case "bf": // そ
		return "そ", nil
	case "c0": // た
		return "た", nil
	case "c1": // ち
		return "ち", nil
	case "c2": // つ
		return "つ", nil
	case "c3": // て
		return "て", nil
	case "c4": // と
		return "と", nil
	case "c5": // な
		return "な", nil
	case "c6": // に
		return "に", nil
	case "c7": // ぬ
		return "ぬ", nil
	case "c8": // ね
		return "ね", nil
	case "c9": // の
		return "の", nil
	case "ca": // は
		return "は", nil
	case "cb": // ひ
		return "ひ", nil
	case "cc": // ふ
		return "ふ", nil
	case "ce": // ほ
		return "ほ", nil
	case "cf": // ま
		return "ま", nil
	case "d0": // み
		return "み", nil
	case "d1": // む
		return "む", nil
	case "d2": // め
		return "め", nil
	case "d3": // も
		return "も", nil
	case "d4": // や
		return "や", nil
	case "d5": // ゆ
		return "ゆ", nil
	case "d7": // ら
		return "ら", nil
	case "d8": // り
		return "り", nil
	case "d9": // る
		return "る", nil
	case "db": // ろ
		return "ろ", nil

	// 貸渡（レンタカー）
	case "da": // れ
		return "れ", nil
	case "dc": // わ
		return "わ", nil

	// 事業用
	case "b1": // あ
		return "あ", nil
	case "b2": // い
		return "い", nil
	case "b3": // う
		return "う", nil
	case "b4": // え
		return "え", nil
	case "b6": // か
		return "か", nil
	case "b7": // き
		return "き", nil
	case "b8": // く
		return "く", nil
	case "b9": // け
		return "け", nil
	case "ba": // こ
		return "こ", nil
	case "a6": // を
		return "を", nil

	// 駐留軍人軍属私有車両用等
	case "45": // E
		return "Ｅ", nil
	case "48": // H
		return "Ｈ", nil
	case "4b": // K
		return "Ｋ", nil
	case "4d": // M
		return "Ｍ", nil
	case "54": // T
		return "Ｔ", nil
	case "59": // Y
		return "Ｙ", nil
	case "d6": // よ
		return "よ", nil
	default:
		return "", fmt.Errorf("Not found Youto from Code.\n")
	}
}
