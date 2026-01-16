#!/bin/bash

### 受信データのcsvファイル(１レコード１ファイル)を、１本のファイルに纏める。
###  tc_csvディレクトリ内のファイルは、tc側で非同期に作成される。
###  しかしながら、１本ファイル化を本スクリプトで行うことにより、ファイル作成中に取りまとめる事故を防ぐことができる。

## tc_csvディレクトリ内の全ファイルを、１本のファイルに纏める(追加する)
## １本のファイルにまとめたら、もとのレコードファイルを削除する（将来的に削除ではなく、移動するとバックアップになる）
## リダイレクトで作成しているので、csv_Table.csvは、新しいデータが末尾（日付：昇順並び）に追加される。

# スクリプトの場所を取得
script_dir=$(dirname "$(realpath "$0")")

# sourceコマンドでパラメータファイルを読み込む
source "${script_dir}/param.sh"

## アンテナ番号により取り出す対象の通過履歴ファイルを決定する
## sbox01_rireki〜sbox04_rirekiは、param.shに設定されている
# sbox01_rireki=../sbox01/tc_csv_table/WCN_rireki.csv
# sbox02_rireki=../sbox02/tc_csv_table/WCN_rireki.csv
# sbox03_rireki=../sbox03/tc_csv_table/WCN_rireki.csv
# sbox04_rireki=../sbox04/tc_csv_table/WCN_rireki.csv
if [[ $1 -eq 1 ]]; then 
    find_file_name=${sbox01_rireki}
elif [[ $1 -eq 2 ]]; then 
    find_file_name=${sbox02_rireki}
elif [[ $1 -eq 3 ]]; then
    find_file_name=${sbox03_rireki}
elif [[ $1 -eq 4 ]]; then
    find_file_name=${sbox04_rireki}
fi

directory=./tc_csv
if [[ -n "$(ls $directory)" ]]; then
    
    # １ファイル１レコードの通過履歴ファイルを一覧化
    cat ./tc_csv/* >> ${find_file_name} 

    # 重複データがある場合は削除（通常運用で、同じデータの複数受信はありえないが、シミュレータ利用の場合は任意のデータを作れるのでありえる。）
    sort ${find_file_name} | uniq > temp.csv && mv temp.csv ${find_file_name}

    rm -rf ./tc_csv/*                                                         # 一覧化した個別ファイルは削除する
    rm -rf ./temp.csv
fi



