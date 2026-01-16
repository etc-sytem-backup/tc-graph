#!/bin/bash
### ../disp_dataディレクトリ内に、passage_count_resetがあった場合、アンテナ毎の通過履歴ファイルを削除する。

## SBOXディレクトリの配列
declare -a sbox_dirs=("$HOME/opt/aps/sbox01" "$HOME/opt/aps/sbox02" "$HOME/opt/aps/sbox03")

## アンテナ毎の通過履歴ファイル保存ディレクトリ
wcn_files_dir="../ac/tc_csv_table/"


## チェックするディレクトリ
check_dir="../disp_data"

## ファイル(passage_count_reset)が存在するかチェック
if [ -f "${check_dir}/passage_count_reset" ]; then

    ## アンテナ毎の通過履歴ファイルを全て削除する
    # SBOX01〜03配下の通過履歴ファイル
    for sbox_dir in "${sbox_dirs[@]}"
    do
        rm -rf "${sbox_dir}/tc_csv_table/WCN_rireki.csv"
    done

    ## passage_count_resetを削除する
    rm "${check_dir}/passage_count_reset"
fi

exit 0
