#!/bin/bash

parking_duration_file="./parking_list/parking_duration_table.csv"

## すでに放置車両一覧ファイルが存在している場合は内容を削除

if [[ -e ${parking_duration_file} ]]; then # ファイルが存在していたら
    : > ${parking_duration_file}           # ファイルの中身をクリアする
else
    touch ${parking_duration_file}         # ファイルを新規作成する
fi
