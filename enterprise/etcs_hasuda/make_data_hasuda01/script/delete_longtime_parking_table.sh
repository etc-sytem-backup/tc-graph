#!/bin/bash

# 引数の確認
if [ "$#" -ne 1 ]; then
    echo "Usage: ./script.sh WCN_NUMBER"
    exit 1
fi

# 引数からWCN番号を取得
wcn_number=$1

# 解析ファイルのパス
input_file="./parking_list/longtime_parking_table.csv"

# 一時ファイルのパス
temp_file="./parking_list/temp.csv"

# WCN番号が一致する行を検索
matching_lines=$(awk -F, -v wcn="$wcn_number" '$4 == wcn' "$input_file")

# 一致する行がある場合
if [ -n "$matching_lines" ]; then
    # 一致する行を解析ファイルから削除
    awk -F, -v wcn="$wcn_number" '$4 != wcn' "$input_file" > "$temp_file" && mv "$temp_file" "$input_file"
fi

exit 0
