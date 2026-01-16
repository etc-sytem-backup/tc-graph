#!/bin/bash

# 各アンテナ（RSU）で検出した入退履歴から、該当のWCN情報だけを抜き出し、直近の２つだけを導き出す。
cmd="echo $1 | sed s/,.*$//"
wcn=$(eval ${cmd})
cat ../ac/tc_csv_table/WCN_rireki.csv | grep -e ${wcn} | head -n 2 | tac > ./driving_history/driving_history.csv


