#!/bin/bash

# csvファイル名
filename=$1

# 列9の値の一番左端の文字が「1」「2」「9」以外になっている行をカウント
count=$(awk -F, '{if (substr($9, 1, 1) !~ /^[129]$/) print $0}' $filename | wc -l)

# 結果を表示
echo $count
