#!/bin/bash

# 引数の取得
csv_string=$1
time_string=$2

# 出力ファイル名の定義
output_file="./parking_list/parkout_table.csv"

# 出力ディレクトリの存在チェックと作成
if [ ! -d "$(dirname "$output_file")" ]; then
  mkdir -p "$(dirname "$output_file")"
fi

# 書き込み文字列の作成
write_string="${csv_string},${time_string}"

# ファイルへの追記
echo "$write_string" >> "$output_file"

# 先週以前の日付データ（行）を全て削除する処理
# ここでは、日付は "YYYYMMDD" 形式と仮定します
# awk コマンドを使用して、1列目の日付が今週のものである行だけを抽出します
# そして、その結果を一時ファイルに書き出し、元のファイルに上書きします

# 今週の日曜日の日付を取得
this_sunday=$(date -d 'last sunday' '+%Y%m%d')

# awk コマンドで今週のデータだけを抽出
# 先頭から8文字を切り出して判断している（ミリ秒を無視している）
awk -v this_sunday="$this_sunday" -F, '{
  record_date = substr($1, 1, 8);
  if (record_date >= this_sunday) print $0
}' "$output_file" > temp.csv

# 元のファイルに上書き
mv temp.csv "$output_file"

exit 0

