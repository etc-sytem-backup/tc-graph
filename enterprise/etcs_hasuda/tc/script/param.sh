#!/bin/bash

### 共通パラメータ設定
deptop=~/opt/aps                     # homeディレクトリ直下のopt/apsにプログラムを配置する

###########################################################
#各アンテナ（RSU）のIPアドレス(OKI : 192.168.110.11 ~ 14)
rsu01_ip_add=192.168.110.11             # RSU01のIPアドレス
rsu02_ip_add=192.168.110.12             # RSU02のIPアドレス
rsu03_ip_add=192.168.110.13             # RSU03のIPアドレス
rsu04_ip_add=192.168.110.14             # RSU04のIPアドレス
rsu01_port=50001                        # 各RSUの指定Port番号
rsu02_port=50001                        # 各RSUの指定Port番号
rsu03_port=50001                        # 各RSUの指定Port番号
rsu04_port=50001                        # 各RSUの指定Port番号
log01_port=50003                        # 各RSUの指定LOGPort番号
log02_port=50003                        # 各RSUの指定LOGPort番号
log03_port=50003                        # 各RSUの指定LOGPort番号
log04_port=50003                        # 各RSUの指定LOGPort番号

## 設定：シミュレーター用（シミュレーターで動かす時は、本番設定をコメントアウトしてこっちを有効化）
# rsu01_ip_add=192.168.1.45             # RSU01のIPアドレス
# rsu02_ip_add=192.168.1.45             # RSU02のIPアドレス
# rsu03_ip_add=192.168.1.45             # RSU03のIPアドレス
# rsu04_ip_add=192.168.1.45             # RSU04のIPアドレス
# rsu01_ip_add=192.168.170.123             # RSU01のIPアドレス
# rsu02_ip_add=192.168.170.123             # RSU02のIPアドレス
# rsu03_ip_add=192.168.170.123             # RSU03のIPアドレス
# rsu04_ip_add=192.168.170.123             # RSU04のIPアドレス
# rsu01_port=51001                         # 各RSUの指定Port番号
# rsu02_port=51002                         # 各RSUの指定Port番号
# rsu03_port=51003                         # 各RSUの指定Port番号
# rsu04_port=51004                         # 各RSUの指定Port番号
# log01_port=51011                         # 各RSUの指定LOGPort番号
# log02_port=51012                         # 各RSUの指定LOGPort番号
# log03_port=51013                         # 各RSUの指定LOGPort番号
# log04_port=51014                         # 各RSUの指定LOGPort番号


##########################################
#SBOXのIPアドレス(OKI : 192.168.110.100)
sbox_ip_add=192.168.110.100                # SBOX01のIPアドレス
sbox01_port=58001                          # SBOXのPort番号(RSU01)
sbox02_port=58002                          # SBOXのPort番号(RSU02)
sbox03_port=58003                          # SBOXのPort番号(RSU03)
sbox04_port=58004                          # SBOXのPort番号(RSU04)

## 設定：シミュレーター用（シミュレーターで動かす時は、本番設定をコメントアウトしてこっちを有効化）
#sbox_ip_add=192.168.1.45              # ese-boxのIPアドレス
#sbox_ip_add=192.168.170.123            # ese-boxのIPアドレス
# sbox01_port=58001                     # ese-boxのPort番号(RSU01)
# sbox02_port=58002                     # ese-boxのPort番号(RSU02)
# sbox03_port=58003                     # ese-boxのPort番号(RSU03)
# sbox04_port=58004                     # ese-boxのPort番号(RSU04)

## acが収集した各アンテナの「通過履歴ファイル」を一本化したもの。
ac_rireki=../ac/tc_csv_table/WCN_rireki.csv
ac_rireki_a1=../ac/tc_csv_table/WCN_rireki_a1.csv
ac_rireki_a2=../ac/tc_csv_table/WCN_rireki_a2.csv
ac_rireki_a3=../ac/tc_csv_table/WCN_rireki_a3.csv
ac_wcn=../ac/tc_wcn_table/WCN_table.csv
