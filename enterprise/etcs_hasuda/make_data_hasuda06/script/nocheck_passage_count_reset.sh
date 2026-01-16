#!/bin/bash
### アンテナ毎の通過履歴ファイルを、問答無用で削除する。

## SBOXディレクトリの配列
declare -a sbox_dirs=("$HOME/opt/aps/sbox01" "$HOME/opt/aps/sbox02" "$HOME/opt/aps/sbox03")

## アンテナ毎の通過履歴ファイルを全て削除する
# SBOX01〜03配下の通過履歴ファイル
for sbox_dir in "${sbox_dirs[@]}"
do
    rm -rf "${sbox_dir}/tc_csv_table/WCN_rireki.csv"
done


exit 0
