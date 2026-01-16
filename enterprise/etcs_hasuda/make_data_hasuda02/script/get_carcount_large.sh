#!/bin/bash

# csvファイル名
#filename=./parking_list/drive_path_table.csv
filename=$1

# 列9の値の一番左端の文字が「1」「2」「9」になっている行をカウント
count=$(awk -F, '{if (substr($9, 1, 1) ~ /^[129]$/) print $0}' $filename | wc -l)

# 結果を表示
echo $count
