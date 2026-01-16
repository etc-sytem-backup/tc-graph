#!/bin/bash
#
# 引数 $1:WCN番号
#
# システムで保持しているWCN番号管理テーブルから、駐車場から出庫した車両のWCNを削除する

# 引数が2つではない場合、エラーメッセージを表示し終了
if [ $# -ne 2 ]; then
    echo "Usage: $0 <stringA> <FileA path>"
    exit 1
fi

# 引数から文字列とファイルパスを取得
stringA=$1
fileA=$2

# FileAから文字列Aを完全一致で検索し、一致した行を削除
awk -F ',' -v str="$stringA" '$1 != str' $fileA > temp.csv

# 左から１列目のデータで降順ソート（数値として）
sort -t',' -k1,1nr temp.csv > $fileA

# 一時ファイルを削除
rm -rf temp.csv

