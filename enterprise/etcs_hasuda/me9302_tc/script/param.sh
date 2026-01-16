#!/bin/bash

### 共通パラメータ設定
deptop=~/opt/aps                     # homeディレクトリ直下のopt/apsにプログラムを配置する

## me9302_tcプロセスが走るPC(APServer)のIPアドレスとポート番号
## me01_ip_add=192.168.1.48:59001
## me02_ip_add=192.168.1.48:59002
## me03_ip_add=192.168.1.48:59003
## me04_ip_add=192.168.1.48:59004
me01_ip_add=192.168.110.110:59001
me02_ip_add=192.168.110.110:59002
me03_ip_add=192.168.110.110:59003
me04_ip_add=192.168.110.110:59004

## ME9302からの受信データを一本化して保存する場所（各アンテナ毎にまとめる）
sbox01_rireki=../sbox01/tc_csv_table/WCN_rireki.csv
sbox02_rireki=../sbox02/tc_csv_table/WCN_rireki.csv
sbox03_rireki=../sbox03/tc_csv_table/WCN_rireki.csv
sbox04_rireki=../sbox04/tc_csv_table/WCN_rireki.csv

## ME9302から受信したWCN番号一覧を保存する場所（各アンテナ毎にまとめる）
sbox01_wcn=../sbox01/tc_wcn_table/WCN_Table.csv
sbox02_wcn=../sbox02/tc_wcn_table/WCN_Table.csv
sbox03_wcn=../sbox03/tc_wcn_table/WCN_Table.csv
sbox04_wcn=../sbox04/tc_wcn_table/WCN_Table.csv

## acが収集した各アンテナの「通過履歴ファイル」を一本化したもの。
ac_rireki=../ac/tc_csv_table/WCN_rireki.csv
ac_wcn=../ac/tc_wcn_table/WCN_table.csv
