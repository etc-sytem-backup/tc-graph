#!/bin/bash

# display用のcsvファイル（disp_passage_count.csv）を作成。
# disp_passage_count.csvの作成場所は./disp_data/とする。
# display_server参照ディレクトリ(~/opt/aps/disp_data)へcsvファイルをコピーする

# 引数の数を確認
if [ "$#" -ne 5 ]; then
    echo "Usage: $0 cnt_antenna1 cnt_antenna2 cnt_antenna3 cnt_antenna4 cnt_antenna5"
    exit 1
fi

# パラメータを取得
cnt_antenna1=$1
cnt_antenna2=$2
cnt_antenna3=$3
cnt_antenna4=$4
cnt_antenna5=$5

# disp_dataディレクトリを作成（すでに存在していてもエラーにならない）
mkdir -p ./disp_data

# CSVファイルを作成
echo "${cnt_antenna1},${cnt_antenna2},${cnt_antenna3},${cnt_antenna4},${cnt_antenna5}" > ./disp_data/disp_passage_count.csv

# ../disp_dataディレクトリを作成（すでに存在していてもエラーにならない）
mkdir -p ../disp_data

# CSVファイルをコピー
cp ./disp_data/disp_passage_count.csv ../disp_data/





