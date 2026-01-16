#!/bin/bash
#
# SBOX01〜04の通過履歴とWCN番号テーブルをacに取り込む
# ファイルをコピーしているのでコストはかかるが、取りこぼしは無い（SBOX側のファイルを消しているわけではない）
#

## ./取込先ディレクトリがない場合は作成する
if [[ ! -d ./tc_csv_table ]]; then
    mkdir ./tc_csv_table
fi

if [[ ! -d ./tc_csv_table/tmp ]]; then
    mkdir ./tc_csv_table/tmp
fi

if [[ ! -d ./tc_wcn_table ]]; then
    mkdir ./tc_wcn_table
fi

#### 通過履歴テーブル
## A1のWCN_rireki.csvをコピー(日付で降順にソートする)
\cp -f ../sbox01/tc_csv_table/WCN_rireki.csv ./tc_csv_table/tmp/WCN_rireki_a1.tmp
cat ./tc_csv_table/tmp/WCN_rireki_a1.tmp | sort -t, -k1,1nr > ./tc_csv_table/tmp/WCN_rireki_a1.csv
rm -rf ./tc_csv_table/tmp/WCN_rireki_a1.tmp

## A2のWCN_rireki.csvをコピー(日付で降順にソートする)
\cp -f ../sbox02/tc_csv_table/WCN_rireki.csv ./tc_csv_table/tmp/WCN_rireki_a2.tmp
cat ./tc_csv_table/tmp/WCN_rireki_a2.tmp | sort -t, -k1,1nr > ./tc_csv_table/tmp/WCN_rireki_a2.csv
rm -rf ./tc_csv_table/tmp/WCN_rireki_a2.tmp

## A3のWCN_rireki.csvをコピー(日付で降順にソートする)
\cp -f ../sbox03/tc_csv_table/WCN_rireki.csv ./tc_csv_table/tmp/WCN_rireki_a3.tmp
cat ./tc_csv_table/tmp/WCN_rireki_a3.tmp | sort -t, -k1,1nr > ./tc_csv_table/tmp/WCN_rireki_a3.csv
rm -rf ./tc_csv_table/tmp/WCN_rireki_a3.tmp

# ## A4のWCN_rireki.csvをコピー(日付で降順にソートする)
# \cp -f ../sbox04/tc_csv_table/WCN_rireki.csv ./tc_csv_table/WCN_rireki_a4.tmp
# cat ./tc_csv_table/WCN_rireki_a4.tmp | sort -t, -k1,1nr > ./tc_csv_table/WCN_rireki_a4.csv
# rm -rf ./tc_csv_table/WCN_rireki_a4.tmp

## コピーしたtc_WCN_rireki.csvを一本に纏める
# tmp.csv と WCN_rireki.csv の内容を結合し、それを一時的なファイル tmp_combined.csv に書き出す
cat ./tc_csv_table/tmp/* | tac > tmp.csv
cat ./tmp.csv ./tc_csv_table/WCN_rireki.csv > ./tmp_combined.csv

# tmp_combined.csv の内容をソートし（1列目を数値として）、重複行を削除する
# その結果を WCN_rireki.csv に書き出す
sort -t, -k1,1nr ./tmp_combined.csv | uniq > ./tc_csv_table/WCN_rireki.csv


# 不要になった一時的なファイルを削除する
rm ./tmp_combined.csv


## 不要ファイルは消す
rm -rf ./tmp.csv
rm -rf ./tmp2.csv

## 通過履歴表示に利用できるかも？ <- 2023/10/23 断面交通量表示の為、消さずに残す方針に変更
#rm -rf ./tc_csv_table/tmp/WCN_rireki_a1.csv
#rm -rf ./tc_csv_table/tmp/WCN_rireki_a2.csv
#rm -rf ./tc_csv_table/tmp/WCN_rireki_a3.csv
#rm -rf ./tc_csv_table/tmp/WCN_rireki_a4.csv
mv ./tc_csv_table/tmp/WCN_rireki_a1.csv ./tc_csv_table
mv ./tc_csv_table/tmp/WCN_rireki_a2.csv ./tc_csv_table
mv ./tc_csv_table/tmp/WCN_rireki_a3.csv ./tc_csv_table


#### WCN管理テーブル
## A1のtc_csv_table.csvをコピー
\cp -f ../sbox01/tc_wcn_table/WCN_table.csv ./tc_wcn_table/WCN_table_a1.csv

## A2のtc_csv_table.csvをコピー
\cp -f ../sbox02/tc_wcn_table/WCN_table.csv ./tc_wcn_table/WCN_table_a2.csv

## A3のtc_csv_table.csvをコピー
\cp -f ../sbox03/tc_wcn_table/WCN_table.csv ./tc_wcn_table/WCN_table_a3.csv

## A4のtc_csv_table.csvをコピー
\cp -f ../sbox04/tc_wcn_table/WCN_table.csv ./tc_wcn_table/WCN_table_a4.csv

## コピーしたtc_wcn_table.csvを一本に纏める(1カラム目のWCNでsortする)
cat ./tc_wcn_table/* | tac > ./tmp.csv                          # 最新のwcnファイル群で一時ファイル(収集した最新のファイル群で作成)を更新 
cat ./tmp.csv | awk -F, '!a[$1]++' > uniq.csv                   # WCN番号が重複している場合は削除する
cat ./uniq.csv | sort -t, -k1,1n > ./tc_wcn_table/WCN_table.csv # 区切り文字は「,」で、カラム１(WCN番号)を数値とみなして保存する。ソート結果はWCN番号の昇順で。

## まとめ終わったら、コピーしたファイルを消す(各アンテナを通過した車両数（通過数ではない）として利用できる？)
## <- 2023/10/23 断面交通量表示の為、消さずに残す方針に変更
rm -rf ./tmp.csv
rm -rf ./uniq.csv
#rm -rf ./tc_wcn_table/WCN_table_a1.csv
#rm -rf ./tc_wcn_table/WCN_table_a2.csv
#rm -rf ./tc_wcn_table/WCN_table_a3.csv
#rm -rf ./tc_wcn_table/WCN_table_a4.csv

