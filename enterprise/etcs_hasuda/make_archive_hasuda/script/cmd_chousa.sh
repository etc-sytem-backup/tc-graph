#!/bin/bash
## make_zip.sh作成のためのテストプログラム。記録として残しておく。

# 引数の数をチェック
if [ "$#" -ne 2 ]; then
    echo "使用方法: $0 [ディレクトリ] [拡張子]"
    exit 1
fi

# 変数に引数を代入
DIRECTORY=$1
EXTENSION=$2

# ディレクトリの存在を確認
if [ ! -d "$DIRECTORY" ]; then
    echo "ディレクトリが存在しません: $DIRECTORY"
    exit 1
fi

# 前日の日付を取得
YESTERDAY=$(date -d "yesterday" +"%Y%m%d")

# ZIPファイル名を設定
ZIPFILE="$DIRECTORY/archive_$YESTERDAY.zip"

# 指定した拡張子を持ち、前日以前に作成または更新されたファイルを検出
FILES=$(find "$DIRECTORY" -type f -name "*.$EXTENSION" -newermt "$YESTERDAY" ! -newermt "$(date +"%Y-%m-%d")")

# ファイルが存在する場合、ZIPにまとめて元のファイルを削除
if [ -n "$FILES" ]; then
    zip -m "$ZIPFILE" $FILES
else
    echo "対象のファイルはありませんでした。"
fi


