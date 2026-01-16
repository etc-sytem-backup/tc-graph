#!/bin/bash

## 引数：　$1:WCN番号　$2アンテナ番号：1〜4
## 指定のWCN番号が、過去の受信データに含まれている場合、受信した日時を返す。
## 見つからない場合は、異常「-1」を返す。
## 引数$2に対し、1〜4以外が渡されたら異常「-1」を返す。

# sourceコマンドでパラメータファイルを読み込む
source "./script/param.sh"

## アンテナ番号により取り出す対象の通過履歴ファイルを決定する
## 1号〜4号の該当する通過履歴ファイル名は、param.shに設定されている
## 各アンテナ毎に、アンテナ毎の通過履歴ファイルを指定できるようにif文で分けている。
## ※現状、アンテナ毎ではなく、全アンテナの通過履歴ファイルを検索対象としている。
##   アンテナがよほど近くに配置されていなければ、渋滞タイマーに引っかかることは無いと判断。
if [[ $2 -eq 1 ]]; then 
    find_file_name=${ac_rireki_a1}
elif [[ $2 -eq 2 ]]; then 
    find_file_name=${ac_rireki_a2}
elif [[ $2 -eq 3 ]]; then
    find_file_name=${ac_rireki_a3}
elif [[ $2 -eq 4 ]]; then
    find_file_name=${ac_rireki}
else
    echo "NoMachine"
fi

## 調査対象の通過履歴ファイルが存在していれば、検索実施。
if [[ -e ${find_file_name} ]]; then

    # 引数(WCN番号)を指定して、受信日付を検索
    cmd=$(cat ${find_file_name} | grep -e $1 | head -n 1 | awk -F ',' '{print $1}')
    if [[ -n ${cmd} ]]; then
        echo ${cmd} # 受信日時を呼び元に返す
    else
        echo "NoHit"
    fi

else
    echo "NotFound"
fi

exit 0
