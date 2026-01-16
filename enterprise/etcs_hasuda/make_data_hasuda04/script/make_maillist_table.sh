#!/bin/bash
## メール送信済みリストを作成する。
## 同じデータは登録しない。
## メール送信済みリストが存在しない場合はtouchコマンドで作成する。

# 引数の取得
send_string=$1

# 出力ファイル名の定義
output_file="./parking_list/maillist_table.csv"

# 出力ディレクトリの存在チェックと作成
if [ ! -d "$(dirname "$output_file")" ]; then
  mkdir -p "$(dirname "$output_file")"
fi

# 書き込み文字列の作成
write_string="${send_string}"

# 変数write_stringと同じデータが、すでにoutput_fileに存在している場合は、データを追記せずにプログラムを終了させる
if grep -Fxq "$write_string" "$output_file"
then
    echo "There is ${write_string} in maillist_table.csv"
    exit 0
fi

# 一時ファイルへの書き込み
echo "$write_string" > "${output_file}.tmp"

# 元のファイルの内容を一時ファイルに追記
if [ -f "$output_file" ]; then
    cat "$output_file" >> "${output_file}.tmp"
fi

# 一時ファイルを元のファイルに移動
mv "${output_file}.tmp" "$output_file"


