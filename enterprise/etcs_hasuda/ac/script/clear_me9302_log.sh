#!/bin/bash

# 設計としてまだ未完成

## 2023/08/01
# 現状のme01〜me03は、受信データを./log/csvに残していない。
# ./tc_csvと./tc_wcnに受信データを記録するが、Bashスクリプトによってsbox01〜03の./tc_csv_tableとtc_wcn_tableに移動(1本化)する

# ME01~04の./log/csvを空にする
mkdir -p ../dust
mkdir -p ../dust/me01
mkdir -p ../dust/me02
mkdir -p ../dust/me03
mkdir -p ../dust/me04

# dustディレクトリへ移動（ゴミ箱の様な使い方）
mv -f ../me01/log/csv/* ../dust/me01
mv -f ../me02/log/csv/* ../dust/me02
mv -f ../me03/log/csv/* ../dust/me03
mv -f ../me04/log/csv/* ../dust/me04




