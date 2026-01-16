#!/bin/bash
##■要件
## ・検出したファイル群をzipファイルにまとめる
## ・zipファイルにまとめた元ファイルは削除する
## ・zipファイル名のフォーマットに日付を含ませてarchive_yyyymmdd.zipとする。
## ・zipファイル名の日付はスクリプトを実行した前日とする。

# 使用方法を表示する関数
usage() {
    echo "使用方法: $0 <対象ディレクトリ> <ファイル拡張子（スペース区切り）>"
    echo "例: $0 /path/to/directory 'csv'"
}

# デバッグ情報を出力する関数
debug() {
    echo "[debug] $1"
}

# 必要なパラメータ数をチェック
if [[ $# -lt 2 ]]; then
    debug "引数が不足しています。"
    usage
    exit 1
fi


# 変数に引数を代入
DIRECTORY=$1
EXTENSION=$2

# ディレクトリの存在を確認
if [ ! -d "$DIRECTORY" ]; then
    debug "ディレクトリが存在しません: $DIRECTORY"
    usage
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
    debug "対象のファイルはありませんでした。"
fi

exit 0


