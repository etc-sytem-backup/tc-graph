#!/bin/bash

# old_WCN_rirekiディレクトリにWCN_rireki.csvが存在していない場合はacからコピーしておく
now_file=../ac/tc_csv_rireki/WCN_rireki.csv
old_file=./old_wcn_rireki/WCN_rireki.csv
if [[ ! -e ${old_file} ]]; then     # 前回の逆走テーブルがコピーされていない場合は、コピーを行う
    cp -rf ${now_file} ${old_file}
fi

# acのWCN_rireki.csvとalert側のWCN_rireki.csvが異なる場合で戻り値を変える
if [[ $(diff -q ${now_file} ${old_file}) ]]; then
    # 不一致
    echo "1"
else
    # 一致
    echo "0"
fi
    
