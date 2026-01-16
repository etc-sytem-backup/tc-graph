#!/bin/bash

## アンテナ毎に通過履歴ファイルのレコード数をカウントする。
## 通過履歴一覧ファイルは、acディレクトリ配下に保存されている。
## 行数(通過車両数)取得コマンド → cat ../ac/tc_csv_table/WCN_rireki.csv | wc -l
## ファイル名を引数とし、通過車両数を戻り値とする。

table_file=$1     # 通過履歴一覧ファイル名(引数)

result=$(cat ${table_file} | wc -l)    # 車両数検出
echo ${result}                         # 車両集を返す

exit 0



