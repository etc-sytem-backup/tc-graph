#!/bin/bash
#
# 引数 $1:WCN番号、$2:追加するレコード情報
#
# 駐車場に侵入している車両の数を管理するテーブルを作成。
# 重複チェックを行っているので、parking_table.csvのレコード数が、駐車している車両の数と一致する。

# 2023/05/25 前バージョンでは、ac管理下の通過履歴を直接見ていたが、呼び元のプログラム中にすでに保持しているのでそれを利用。
# res=$(grep -m1 -e $1 ../ac/ac_csv/WCN_rireki.csv)
# if [ $? -ne 0 ]; then
#     echo ${res}
#     echo "$1が ../ac/ac_csv/WCN_rireki.csv 内に見つかりません。"
#     exit 0   # 見つからなければ抜ける
# fi

res=$2
file="./parking_list/parking_table.csv"
if [ -e ${file} ]; then     # 駐車車両テーブルがすでに存在していたら
    grep -e $1 ./${file}    # 駐車車両テーブルにWCNが登録されているかチェック
    if [ $? -ne 0 ]; then   # 登録されていないなら登録。登録されていたら古いデータを削除し、新しいデータを追加する

        # 履歴レコードをparking_table.csvの１行目に追加する
        echo "${res}を登録します。"
        echo ${res} > ./parking_table.tmp; cat ${file} >> parking_table.tmp; cat ./parking_table.tmp > ${file}
        rm -rf ./parking_table.tmp
        exit 0
    else
        num=$(grep -m1 -n -e $1 ./${file} | awk -F ':' '{print $1}') # ヒットした行番号を取得
        sed "${num}d" ${file} > WCN_rireki.tmp  # その行を削除した一時ファイルを作成
        echo ${res} > ./line.tmp; cat ./WCN_rireki.tmp >> ./line.tmp; cat ./line.tmp > ${file}  # parking_table.csvの１行目に最新の情報を追加する。
        rm -rf ./WCN_rireki.tmp
        rm -rf ./line.tmp
        exit 0
    fi
else                      # 駐車車両テーブルが存在していない場合は新規作成する

    echo "${file}を新規作成します。"
    touch ${file}
    echo ${res} >> ${file}
fi
exit 0
