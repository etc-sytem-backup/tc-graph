#!/bin/bash

## ../disp_dataディレクトリ内に、alert_resetがあった場合、下記のファイルを削除。
## ./parking_list/alert_table.csv
## ./parking_list/disp_alert.csv
## このファイルを削除する事で、逆走検知モニタに表示されている履歴が削除される。
## 但し、WCN_rireki.csvに残っている最新の現状については、元データとなる車両を現在位置から移動させるなどの対応をしない限り消えない。
## 結論として、過去の履歴は消えるが現在進行形のデータは消えない。

# 削除対象ファイル名の定義
delete_file_01="./parking_list/alert_table.csv"
delete_file_02="./disp_data/disp_alert.csv"

# alert_resetが作成されるディレクトリ
check_dir="../disp_data"

# ファイル(alert_reset)が存在するかチェック
if [[ -f "${check_dir}/alert_reset" ]]; then

    # ファイルが存在するならば、delete_file_01及び02を削除し、alert_resetも削除する。
    rm -rf "${delete_file_01}"
    rm -rf "${delete_file_02}"
    rm -rf "${check_dir}/alert_reset"
fi

exit 0
