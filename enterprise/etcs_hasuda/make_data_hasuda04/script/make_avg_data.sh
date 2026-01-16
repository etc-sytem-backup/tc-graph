#!/bin/bash
## input_fileを解析し下記のデータを導き出す。
## ・車種別にデータを分ける
## ・当週、当月での平均駐車時間を演算
## ランプ停滞中表示用のフラグ(1 or 0)は、引数として取得する。

# ファイルパス
input_file="./parking_list/parkout_table.csv"
output_file="./disp_data/disp_avg.csv"

# ランプ停滞中データ
traffic_jam=$1

# 入力ファイルが存在するかチェック
# まだファイルが作成されていない場合は処理せず抜ける。
if [[ ! -f "$input_file" ]]; then
    echo "Input file does not exist. Exiting."
    exit 0
fi

# 今日の日付と今週の最初の日を取得
today=$(date +%Y%m%d)
week_start=$(date --date="this monday" +%Y%m%d) # 今週の月曜日始まりで1週間
# week_start=$(date --date="this sunday" +%Y%m%d) # 今週の日曜日始まりで1週間
# week_start=$(date --date="last monday" +%Y%m%d) # 先週の月曜日始まりで1週間
# week_start=$(date --date="last sunday" +%Y%m%d) # 先週の日曜日始まりで1週間


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


