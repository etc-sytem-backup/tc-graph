#!/bin/bash
## 蓮田SAで収集したデータを、毎日バックアップ（zipファイル）
## このスクリプトを実行した日をファイル名に含める。
## 同日に複数回実行されたとしても、作成されるzipファイルが更新されるのみ。

#!/bin/bash

# Get yesterday's date in YYYYMMDD format
yesterday_date=$(date -d "yesterday" +%Y%m%d)

# Change current directory to /home/etcs_hasuda/opt/aps/ac
cd /home/etcs_hasuda/opt/aps/ac

# Rename directories with yesterday's date
mv ./tc_csv_table ./tc_csv_table.$yesterday_date
mv ./tc_wcn_table ./tc_wcn_table.$yesterday_date

# Compress the renamed directories
zip -r ./tc_csv_table.$yesterday_date.zip ./tc_csv_table.$yesterday_date
zip -r ./tc_wcn_table.$yesterday_date.zip ./tc_wcn_table.$yesterday_date

# Change current directory to /home/etcs_hasuda/opt/aps
cd /home/etcs_hasuda/opt/aps

# Copy and rename the directory
cp -rp ./disp_data ./disp_data.$yesterday_date

# Compress the copied directory
zip -r ./disp_data.$yesterday_date.zip ./disp_data.$yesterday_date

