#!/bin/bash 

### display連携用のcsvファイル(逆走検知モニタ用)を作成する。
## displayに連携する日付は、西暦の先頭文字を削除したものになっている・・・ので、先頭文字だけ削除してコピーする

# 入力ファイル
input_file="./parking_list/alert_table.csv"

# 出力ファイル
output_file="./disp_data/disp_alert.csv"
output_disp_file="../disp_data/disp_alert.csv"

# 出力ファイルが存在する場合、それを空にする
> "${output_file}"

# 入力ファイルを1行ずつ読み込む
while IFS= read -r line; do

    # CSVデータを配列に格納する
    IFS=',' read -ra line_val <<< "$line"

    # line_valの2，3，5列目を除外し、新しい行を作成する
    new_line="${line_val[0]},${line_val[3]},${line_val[5]},${line_val[6]},${line_val[7]},${line_val[8]},${line_val[9]},${line_val[10]}"

    # 先頭の1文字を削除
    new_line="${new_line:1}"

    # データを出力ファイルに追加
    echo "${new_line}" >> "${output_file}"

done < "$input_file"

> "${output_disp_file}"
cp -rp ${output_file} ${output_disp_file}

exit 0
