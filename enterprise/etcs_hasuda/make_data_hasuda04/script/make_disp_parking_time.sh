#!/bin/bash 

### display連携用のcsvファイルを作成する。
## displayに連携する日付は、西暦の先頭文字を削除したものになっている・・・ので、先頭文字だけ削除してコピーする

# 入力ファイル
input_file="./parking_list/longtime_parking_table.csv"

# 出力ファイル
output_file="./disp_data/disp_parking_time.csv"
output_disp_file="../disp_data/disp_parking_time.csv"

# 出力ファイルが存在する場合、それを空にする
> "${output_file}"

# 入力ファイルを1行ずつ読み込む
while IFS= read -r line; do
    # 先頭の1文字を削除して出力ファイルに追加
    echo "${line:1}" >> "${output_file}"
done < "$input_file"

> "${output_disp_file}"
cp -rp ${output_file} ${output_disp_file}

exit 0
