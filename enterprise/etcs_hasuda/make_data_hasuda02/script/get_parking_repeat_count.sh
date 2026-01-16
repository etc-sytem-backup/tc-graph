#!/bin/bash

source ./param.sh # 共通パラメータ読み込み

# 入力ファイルパス
INPUT_FILE="./parking_list/parking_repeat_table.csv"

# 出力ファイルパス
MONTH_FILE="./disp_data/disp_table_month.csv"
WEEK_FILE="./disp_data/disp_table_week.csv"

# 入力ファイルが存在するか確認し、無ければ何もせずに抜ける。
if [ ! -f $INPUT_FILE ]; then
    echo "Input file does not exist. Exiting."
    exit 0
fi

# 現在の年月（yyyy-mm）、現在の年週（yyyy-ww）
CURRENT_MONTH=$(date +"%Y%m")
CURRENT_WEEK=$(date +"%Y%V")

# 入力ファイルからWCN番号を抽出し重複を削除
WCN_LIST=$(awk -F, '{print $4}' $INPUT_FILE | sort | uniq)

# 抽出したWCN番号についてループ
for WCN in $WCN_LIST; do

    # マッチする行のカウント、最新の行の取得
    COUNT=0
    LAST_LINE=""

    while IFS= read -r line; do
        if [[ "$line" == *"$WCN"* ]]; then
            COUNT=$((COUNT+1))
            LAST_LINE="$line"
        fi
    done < $INPUT_FILE

    # 最新の行からデータ抽出
    DATE=$(echo $LAST_LINE | awk -F, '{print substr($1,1,17)}')
    ETC=$(echo $LAST_LINE | awk -F, '{print $6}')
    BRANCH=$(echo $LAST_LINE | awk -F, '{print $7}')
    USE=$(echo $LAST_LINE | awk -F, '{print $8}')
    TYPE=$(echo $LAST_LINE | awk -F, '{print $9}')
    NUMBER=$(echo $LAST_LINE | awk -F, '{print $10}')

    # データ整形
    OUTPUT="${DATE:1},${WCN},${ETC},${BRANCH},${USE},${TYPE},${NUMBER},${COUNT}"

    # 当月のデータを月次ファイルに書き込む
    if [[ "${DATE:0:6}" == "$CURRENT_MONTH" ]]; then
        if grep -q "$WCN" $MONTH_FILE; then
            sed -i "/$WCN/d" $MONTH_FILE
        fi
        echo $OUTPUT >> $MONTH_FILE
    fi

    # 当週のデータを週次ファイルに書き込む
    DATE_WEEK=$(date -d"${DATE:0:4}-${DATE:4:2}-${DATE:6:2}" +"%Y%V")
    if [[ "$DATE_WEEK" == "$CURRENT_WEEK" ]]; then
        if grep -q "$WCN" $WEEK_FILE; then
            sed -i "/$WCN/d" $WEEK_FILE
        fi
        echo $OUTPUT >> $WEEK_FILE
    fi
done

# ファイルをカウントの降順でソート
sort -t, -k8 -nr -o $MONTH_FILE $MONTH_FILE
sort -t, -k8 -nr -o $WEEK_FILE $WEEK_FILE

# 当月以外のデータを削除
awk -v month="$CURRENT_MONTH" -F, '{if (substr($1,1,6) == month) print $0}' $INPUT_FILE > tmp.csv && mv tmp.csv $INPUT_FILE

# 当月データと当週データをdisplay連絡用ディレクトリへコピー
cp -p ${MONTH_FILE} ../disp_data/
cp -p ${WEEK_FILE} ../disp_data/

exit 0

