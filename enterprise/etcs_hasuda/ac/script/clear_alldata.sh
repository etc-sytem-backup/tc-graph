#!/bin/bash

# 前日の日付を取得
yyyymmdd=$(date -d "yesterday" '+%Y%m%d')

# システムのデータをオールクリア

# ac
rm -rf ./tc_csv_table/*
rm -rf ./tc_wcn_table/* 

# SBOX01~04
rm -rf ../sbox01/tc_csv_table/*
rm -rf ../sbox02/tc_csv_table/*
rm -rf ../sbox03/tc_csv_table/*
rm -rf ../sbox04/tc_csv_table/*
rm -rf ../sbox01/tc_wcn_table/*
rm -rf ../sbox02/tc_wcn_table/*
rm -rf ../sbox03/tc_wcn_table/*
rm -rf ../sbox04/tc_wcn_table/*

# SBOX01~04の通信ログを移動
mkdir -p ../dust
mkdir -p ../dust/sbox01
mkdir -p ../dust/sbox02
mkdir -p ../dust/sbox03
mkdir -p ../dust/sbox04


# SBOX01
zip -j "../sbox01/log/csv/sbox01.${yyyymmdd}.zip" ../sbox01/log/csv/*.csv > ./sbox01.zip.txt
mv -f "../sbox01/log/csv/sbox01.${yyyymmdd}.zip" ../dust/sbox01 > ./sbox01.mv.txt

# SBOX02
zip -j "../sbox02/log/csv/sbox02.${yyyymmdd}.zip" ../sbox02/log/csv/*.csv > ./sbox02.zip.txt
mv -f "../sbox02/log/csv/sbox02.${yyyymmdd}.zip" ../dust/sbox02 > ./sbox02.mv.txt

# SBOX03
zip -j "../sbox03/log/csv/sbox03.${yyyymmdd}.zip" ../sbox03/log/csv/*.csv > ./sbox03.zip.txt
mv -f "../sbox03/log/csv/sbox03.${yyyymmdd}.zip" ../dust/sbox03 > ./sbox03.mv.txt

# SBOX04
zip -j "../sbox04/log/csv/sbox04.${yyyymmdd}.zip" ../sbox04/log/csv/*.csv > ./sbox04.zip.txt
mv -f "../sbox04/log/csv/sbox04.${yyyymmdd}.zip" ../dust/sbox04 > ./sbox04.mv.txt

# 圧縮→移動 済のcsvファイル達を削除
rm -rf ../sbox01/log/csv/*.csv
rm -rf ../sbox02/log/csv/*.csv
rm -rf ../sbox03/log/csv/*.csv
rm -rf ../sbox04/log/csv/*.csv


# make_data_hasuda01~04
rm -rf ../make_data_hasuda01/disp_data/*
rm -rf ../make_data_hasuda01/parking_list/*
rm -rf ../make_data_hasuda01/driving_history/*

rm -rf ../make_data_hasuda02/disp_data/*
rm -rf ../make_data_hasuda02/parking_list/*
rm -rf ../make_data_hasuda02/driving_history/*

rm -rf ../make_data_hasuda03/disp_data/*
rm -rf ../make_data_hasuda03/parking_list/*
rm -rf ../make_data_hasuda03/driving_history/*

rm -rf ../make_data_hasuda04/disp_data/*
rm -rf ../make_data_hasuda04/parking_list/*
rm -rf ../make_data_hasuda04/driving_history/*

rm -rf ../make_data_hasuda05/disp_data/*
rm -rf ../make_data_hasuda05/parking_list/*
rm -rf ../make_data_hasuda05/driving_history/*

rm -rf ../make_data_hasuda06/disp_data/*
rm -rf ../make_data_hasuda06/parking_list/*
rm -rf ../make_data_hasuda06/driving_history/*


# display連携用ディレクトリ
rm -rf ../disp_data/*

# display連携用初期ファイル作成とデータ挿入
# touch ../disp_data/disp_main_csv
# touch ../disp_data/disp_avg_csv
# touch ../disp_data/disp_table_week_csv
# touch ../disp_data/disp_table_month_csv
# touch ../disp_data/disp_parking_time_csv
# touch ../disp_data/disp_alert_csv
# echo "0,0,0,0,0,0" > ../disp_data/disp_main_csv
# echo "0,0,0,0,0,0,0" > ../disp_data/disp_avg_csv
# echo "0,0,0,0,0,0,0,0" > ../disp_data/disp_table_week_csv
# echo "0,0,0,0,0,0,0,0" > ../disp_data/disp_table_month_csv
# echo "0,0,0,0,0,0,0,0" > ../disp_data/disp_parking_time_csv
# echo "0,0,0,0,0,0,0,0" > ../disp_data/disp_alert_csv
