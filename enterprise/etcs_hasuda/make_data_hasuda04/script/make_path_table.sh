#!/bin/bash
#
# $1:追加データ
#

# パス通過管理ファイル
file_path=./parking_list/drive_path_table.csv

# 追加データ
data=$1

# データが既に存在するかチェック
if grep -Fxq "$data" "$file_path"; then
    echo 1
    exit 0
fi

# ファイルの内容を一時ファイルに退避
tmp_file=$(mktemp)
cp "$file_path" "$tmp_file"

# データを追加
echo "$data" > "$file_path"
cat "$tmp_file" >> "$file_path"

# 一時ファイルを削除
rm "$tmp_file"

echo 0
exit 0


