#!/bin/bash
#
# 引数 $1:追加するレコード情報
#
# 駐車場に侵入した車両情報を./parking_list/parking_repeat_table.csvへ追記する。


# CSV文字列
csv_string=$1

# 対象ファイル
target_file="./parking_list/parking_repeat_table.csv"

# ファイルが存在しない場合は新規作成
if [ ! -f $target_file ]; then
  touch $target_file
fi

# 現在の年と月を取得
current_year=$(date +%Y)
current_month=$(date +%m)

# 当月のデータのみを一時ファイルに書き出す
awk -F, -v year=$current_year -v month=$current_month \
  'substr($1,1,4) == year && substr($1,5,2) == month' $target_file > temp.csv

# 一時ファイルを元のファイルに移動
mv temp.csv $target_file

# CSV文字列をファイルの末尾に追加
echo $csv_string >> $target_file

exit 0
