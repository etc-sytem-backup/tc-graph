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
    if [ $? -ne 0 ]; then   # 指定のWCNが登録されていないなら何もせず。登録されていたらデータを削除。
        exit 0
    else
        num=$(grep -m1 -n -e $1 ./${file} | awk -F ':' '{print $1}') # ヒットした行番号を取得
        sed "${num}d" ${file} > parking_table.tmp  # 満空管理テーブルから、その行を削除する
        cat ./parking_table.tmp > ${file}
        rm -rf ./parking_table.tmp
        echo "$1を満空管理テーブルから削除しました。"
    fi
else                      # 駐車車両テーブルが存在していない場合は新規作成する
    echo "${file}を新規作成します。"
    touch ${file}
    echo ${res} >> ${file}
fi
exit 0
