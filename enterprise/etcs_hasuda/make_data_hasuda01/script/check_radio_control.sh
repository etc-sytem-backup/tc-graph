#!/bin/bash
# ./disp_radio_statusの内容を読み取り、戻り値を初期化。無ければ作成（内容は数値文字列の"0"）する。
# ../disp_dataディレクトリ内に、radio_startがあった場合、SBOX01〜04に無線開始指示し、戻り値を1とする
# ../disp_dataディレクトリ内に、radio_stopがあった場合、SBOX01〜04に無線停止指示し、戻り値を0とする。
# radio_startとradio_stopを検出した時、./disp_radio_statusの内容を変更する。（radio_start = 1, radio_stop = 0）

# SBOXディレクトリの配列
declare -a sbox_dirs=("$HOME/opt/aps/sbox01" "$HOME/opt/aps/sbox02" "$HOME/opt/aps/sbox03")

# チェックするディレクトリ
check_dir="../disp_data"
status_file="./disp_radio_status"

# ./disp_radio_statusの存在チェックと読み込み/作成
if [ -f "$status_file" ]; then
    result=$(cat "$status_file")
else
    echo "0" > "$status_file"
    result="0"
fi

# ファイル(radio_start)が存在するかチェック
if [ -f "${check_dir}/radio_start" ]; then
    for sbox_dir in "${sbox_dirs[@]}"
    do
        touch "${sbox_dir}/SIRONNO"
    done
    rm "${check_dir}/radio_start"
    echo "1" > "$status_file"
    result="1"
fi

# ファイル(radio_stop)が存在するかチェック
if [ -f "${check_dir}/radio_stop" ]; then
    for sbox_dir in "${sbox_dirs[@]}"
    do
        touch "${sbox_dir}/SIROFFNO"
    done
    rm "${check_dir}/radio_stop"
    echo "0" > "$status_file"
    result="0"
fi

# 終了前にresultを出力
echo "$result"

exit 0

