#!/bin/bash
# 表示用端末（display）が作成する「在車数オフセットファイル」を検出し、戻り値を設定する。
# ファイル名: disp_setting.csv
# ファイルは1行だけのcsvファイル
# 　→例: 36,64,0,0,0
# 　　　第一項目: 大型車両駐車台数初期値
# 　　　第二項目: 大型車両以外の駐車台数初期値
#
# ポーリングするディレクトリは「../disp_data」とする。
# 1行目のcsv文字列を返す


# チェックするディレクトリ
check_dir="../disp_data"
target_file="disp_setting.csv"

# ファイルが存在するかチェック
# 無ければファイルを作成する(初期値 0,0,0,0,0)
if [ ! -f "${check_dir}/${target_file}" ]; then
    touch "${check_dir}/${target_file}
    echo "0,0,0,0,0" > "${check_dir}/${target_file}
fi

cat "${check_dir}/${target_file}"

exit 0


