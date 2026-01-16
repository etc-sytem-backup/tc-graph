#!/bin/bash
## 逆走車両情報と逆走コメントをつなげて、逆走テーブルの先頭行に追加する。
## このスクリプトで作成されるのは、display用データではない。
## display表示用データは、別スクリプトを作成し、このスクリプトで作成される出力ファイルを加工して作成すること。

# 引数の取得
csv_string=$1
alert_string=$2

# 出力ファイル名の定義
output_file="./parking_list/alert_table.csv"

# 出力ディレクトリの存在チェックと作成
if [ ! -d "$(dirname "$output_file")" ]; then
  mkdir -p "$(dirname "$output_file")"
fi

# 書き込み文字列の作成
write_string="${csv_string},${alert_string}"

# 一時ファイルへの書き込み
echo "$write_string" > "${output_file}.tmp"

# 元のファイルの内容を一時ファイルに追記
if [ -f "$output_file" ]; then
    cat "$output_file" >> "${output_file}.tmp"
fi

# 一時ファイルを元のファイルに移動
mv "${output_file}.tmp" "$output_file"

