#!/bin/bash

# display用のcsvファイル（disp_rmain.csv）を作成。
# disp_main.csvの作成場所は./disp_data/とする。
# display_server参照ディレクトリ(~/opt/aps/disp_data)へcsvファイルをコピーする

# 引数の数を確認
if [ "$#" -ne 7 ]; then
    echo "Usage: $0 large_in_parking_cnt other_in_parking_cnt large_parking_space other_parking_space large_drivepath_cnt other_drivepath_cnt radio_status"
    exit 1
fi

# パラメータを取得
large_in_parking_cnt=$1
other_in_parking_cnt=$2
large_parking_space=$3
other_parking_space=$4
large_drivepath_cnt=$5
other_drivepath_cnt=$6
radio_status=$7

# disp_dataディレクトリを作成（すでに存在していてもエラーにならない）
mkdir -p ./disp_data

# CSVファイルを作成
echo "${large_in_parking_cnt},${large_parking_space},${large_drivepath_cnt},${other_in_parking_cnt},${other_parking_space},${other_drivepath_cnt},${radio_status}" > ./disp_data/disp_main.csv

# ../disp_dataディレクトリを作成（すでに存在していてもエラーにならない）
mkdir -p ../disp_data

# CSVファイルをコピー
cp ./disp_data/disp_main.csv ../disp_data/





