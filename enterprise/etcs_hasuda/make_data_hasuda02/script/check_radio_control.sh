#!/bin/bash

# ../disp_dataディレクトリ内に、radio_startがあった場合、SBOX01〜04に無線開始指示。
# ../disp_dataディレクトリ内に、radio_stopがあった場合、SBOX01〜04に無線停止指示。

#!/bin/bash

# SBOXディレクトリの配列
declare -a sbox_dirs=("$HOME/opt/aps/sbox01" "$HOME/opt/aps/sbox02" "$HOME/opt/aps/sbox03")

# チェックするディレクトリ
check_dir="../disp_data"

# ファイル(radio_start)が存在するかチェック
if [ -f "${check_dir}/radio_start" ]; then
    # ファイルが存在するならば、全SBOXディレクトリに、空ファイル（SIRONNO）を作成する
    for sbox_dir in "${sbox_dirs[@]}"
    do
        touch "${sbox_dir}/SIRONNO"
    done
    # 空ファイル（SIRONNO）作成後、radio_startを削除する
    rm "${check_dir}/radio_start"
fi

# ファイル(radio_stop)が存在するかチェック
if [ -f "${check_dir}/radio_stop" ]; then
    # ファイルが存在するならば、全SBOXディレクトリに、空ファイル（SIROFFNO）を作成する
    for sbox_dir in "${sbox_dirs[@]}"
    do
        touch "${sbox_dir}/SIROFFNO"
    done
    # 空ファイル（SIROFFNO）作成後、radio_stopを削除する
    rm "${check_dir}/radio_stop"
fi

exit 0
