#!/bin/bash
## input_fileを解析し下記のデータを導き出す。
## ・車種別にデータを分ける
## ・当週、当月での平均駐車時間を演算
## ランプ停滞中表示用のフラグ(1 or 0)は、引数として取得する。

# ファイルパス
input_file_A="./parking_list/parkout_table.csv"
input_file_B="./parking_list/parkingtime_table.csv"
input_file_tmp="./parking_list/input_file_tmp.csv"
input_file="./parking_list/all_parkingtime_table.csv"
output_file="./disp_data/disp_avg.csv"

# parkout_table.csv  : 駐車場を退場した車両の駐車時間一覧
# parktime_table.csv : 現在入庫している全車両の駐車時間一覧
# どちらも未作成ならば、何もせずに処理を抜ける。
if [[ ! -f "$input_file_A" ]] && [[ ! -f "$input_file_B" ]]; then
    echo "Neither of the input files exist. Exiting."
    exit 0
fi

# ファイルAとファイルBを結合して一時ファイルを作成
cat $input_file_A $input_file_B > $input_file_tmp

# 一時ファイルの第一列を数値とみなし、昇順に並び替えてinput_fileを作成
sort -n -t, -k1 $input_file_tmp > $input_file

# 一時ファイルを削除
rm $input_file_tmp

# 最終出力ファイル
output_file="./disp_data/disp_avg.csv"

# ランプ停滞中データ
traffic_jam=$1


## 2023/07/14 
# "last monday"とすれば、直近の月曜日を一週間のはじめとし、月曜日〜日曜日となる。んが、、、、
# プログラムを動作させる日が月曜日だった場合、月曜日が終わってない（火曜日になってない）ので、
# 先週の月曜日始まりになってしまう。。。
# よって、単純に"last monday"を指定するのでは無く、本スクリプトが月曜日に実行された場合は、
# week_startを当日の月曜日日付になるように処理する。
#
# # 今日の日付と今週の最初の日を取得
# today=$(date +%Y%m%d)
# # week_start=$(date --date="this monday" +%Y%m%d) # 今週の月曜日始まりで1週間
# # week_start=$(date --date="this sunday" +%Y%m%d) # 今週の日曜日始まりで1週間
# week_start=$(date --date="last monday" +%Y%m%d) # 先週の月曜日始まりで1週間
# # week_start=$(date --date="last sunday" +%Y%m%d) # 先週の日曜日始まりで1週間


# 今日の日付を取得
today=$(date +%Y%m%d)

# 今日の曜日を取得
day_of_week=$(date +%u)

# 今週の最初の日を取得
if [ "$day_of_week" -eq 1 ]; then
    # 今日が月曜日ならば、week_startは今日の日付
    week_start=$today
else
    # 今日が月曜日以外ならば、week_startは前の週の月曜日の日付
    week_start=$(date --date="last monday" +%Y%m%d)
fi


# デバッグログ：日付情報
echo "DEBUG: Today's date: $today"
echo "DEBUG: Week start date: $week_start"

# 変数を初期化
total_parking_time_today=0
total_parking_time_week=0
total_large_vehicle_parking_time_today=0
total_large_vehicle_parking_time_week=0
total_non_large_vehicle_parking_time_today=0
total_non_large_vehicle_parking_time_week=0

count_today=0
count_week=0
count_large_today=0
count_large_week=0
count_non_large_today=0
count_non_large_week=0

# ファイルを読み込む
while IFS=, read -r datetime antenna alias wcn status etc_card office_code usage vehicle_type serial_number parking_time; do
    declare -i parking_time

    # デバッグ：読み込んだデータを表示
    echo "DEBUG: Read data: $datetime, $antenna, $alias, $wcn, $status, $etc_card, $office_code, $usage, $vehicle_type, $serial_number, $parking_time"

    # 日付を抽出
    date=${datetime:0:8}

    echo "DEBUG: Read date: $date"
    echo "DEBUG: Comparing with today: $today and week start: $week_start"
    
    # 大型車かどうかを判定
    first_char_vehicle_type=${vehicle_type:0:1}
    is_large_vehicle=false
    if [ "$first_char_vehicle_type" == "1" ] || [ "$first_char_vehicle_type" == "9" ] || [ "$first_char_vehicle_type" == "2" ]; then
        is_large_vehicle=true
    fi

    echo "DEBUG: Is large vehicle: $is_large_vehicle"
    
    # 当日のデータ
    if [ "$date" == "$today" ]; then
        total_parking_time_today=$((total_parking_time_today + parking_time))
        count_today=$((count_today + 1))
        
        if $is_large_vehicle; then
            total_large_vehicle_parking_time_today=$((total_large_vehicle_parking_time_today + parking_time))
            count_large_today=$((count_large_today + 1))
        else
            total_non_large_vehicle_parking_time_today=$((total_non_large_vehicle_parking_time_today + parking_time))
            count_non_large_today=$((count_non_large_today + 1))
        fi
    fi
    
    # 当週のデータ
    if [ "$date" -ge "$week_start" ] && [ "$date" -le     "$today" ]; then
        total_parking_time_week=$((total_parking_time_week + parking_time))
        count_week=$((count_week + 1))
        
        if $is_large_vehicle; then
            total_large_vehicle_parking_time_week=$((total_large_vehicle_parking_time_week + parking_time))
            count_large_week=$((count_large_week + 1))
        else
            total_non_large_vehicle_parking_time_week=$((total_non_large_vehicle_parking_time_week + parking_time))
            count_non_large_week=$((count_non_large_week + 1))
        fi
    fi
done < "$input_file"

# デバッグログ：集計結果
echo "DEBUG: Total parking time today: $total_parking_time_today"
echo "DEBUG: Total parking time this week: $total_parking_time_week"
echo "DEBUG: Total large vehicle parking time today: $total_large_vehicle_parking_time_today"
echo "DEBUG: Total large vehicle parking time this week: $total_large_vehicle_parking_time_week"
echo "DEBUG: Total non large vehicle parking time today: $total_non_large_vehicle_parking_time_today"
echo "DEBUG: Total non large vehicle parking time this week: $total_non_large_vehicle_parking_time_week"

# 平均を計算
if [ $count_today -eq 0 ]; then
    avg_parking_time_today=0
else
    avg_parking_time_today=$((total_parking_time_today / count_today))
fi

if [ $count_week -eq 0 ]; then
    avg_parking_time_week=0
else
    avg_parking_time_week=$((total_parking_time_week / count_week))
fi

if [ $count_large_today -eq 0 ]; then
    avg_large_vehicle_parking_time_today=0
else
    avg_large_vehicle_parking_time_today=$((total_large_vehicle_parking_time_today / count_large_today))
fi

if [ $count_large_week -eq 0 ]; then
    avg_large_vehicle_parking_time_week=0
else
    avg_large_vehicle_parking_time_week=$((total_large_vehicle_parking_time_week / count_large_week))
fi

if [ $count_non_large_today -eq 0 ]; then
    avg_non_large_vehicle_parking_time_today=0
else
    avg_non_large_vehicle_parking_time_today=$((total_non_large_vehicle_parking_time_today / count_non_large_today))
fi

if [ $count_non_large_week -eq 0 ]; then
    avg_non_large_vehicle_parking_time_week=0
else
    avg_non_large_vehicle_parking_time_week=$((total_non_large_vehicle_parking_time_week / count_non_large_week))
fi


# 結果をCSVファイルに保存
echo "${traffic_jam},${avg_parking_time_today},${avg_parking_time_week},${avg_large_vehicle_parking_time_today},${avg_large_vehicle_parking_time_week},${avg_non_large_vehicle_parking_time_today},${avg_non_large_vehicle_parking_time_week}" > ${output_file}

# CSVファイルをdisplay用公開ディレクトリへコピー
cp ${output_file} ../disp_data/

exit 0


