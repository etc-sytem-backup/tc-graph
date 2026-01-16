#!/bin/bash

## ../disp_dataディレクトリ内に、alert_resetがあった場合、下記のファイルを削除。
## ./parking_list/alert_table.csv
## このファイルを削除する事で、逆走検知モニタに表示されている履歴が削除される。

# 削除対象ファイル名の定義
delete_file="./parking_list/alert_table.csv"

# alert_resetが作成されるディレクトリ
check_dir="../disp_data"

# ファイル(alert_reset)が存在するかチェック
if [[ -f "${check_dir}/alert_reset" ]]; then

    # ファイルが存在するならば、delete_fileを削除し、alert_resetも削除する。
    rm -rf "${delete_file}"
    rm -rf "${check_dir}/alert_reset"
fi

exit 0
