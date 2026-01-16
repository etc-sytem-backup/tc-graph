#!/bin/bash
#
# $1: WCN番号
# $2: CSVファイル（フルパス）
#

# 引数が2つではない場合、エラーメッセージを表示し終了
if [ $# -ne 2 ]; then
    echo "Usage: $0 <WCN number> <CSV file path>"
    exit 1
fi

# 引数からWCN番号とCSVファイルのパスを取得
wcnNumber=$1
csvFile=$2

# CSVファイルの4列目（WCN番号）をチェックし、指定されたWCN番号と一致する行があるか確認
matchingLines=$(awk -F ',' -v wcn="$wcnNumber" '$4 == wcn' "$csvFile")

if [ -n "$matchingLines" ]
then
    echo 1
else
    echo 0
fi

exit 0

