#!/bin/bash

# 引数が2つではない場合、エラーメッセージを表示し終了
if [ $# -ne 2 ]; then
    echo "Usage: $0 <SubtractionDays> <FileA path>"
    exit 1
fi

# 引数から減算日数とファイルパスを取得
subtractionDays=$1
fileA=$2

# 現在日付から減算日数を引き算し、削除対象日付を作成
targetDate=$(date -d "$subtractionDays days ago" "+%Y%m%d")

# FileAの1列目（日付）をチェックし、削除対象日付より小さい日付の行を削除
awk -F ',' -v target="$targetDate" 'substr($1, 1, 8) >= target' $fileA > temp.csv

# 変更をFileAに反映
mv temp.csv $fileA

