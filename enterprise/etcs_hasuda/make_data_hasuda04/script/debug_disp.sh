#!/bin/bash

# Debug用に各ファイルを端末に表示する。

echo "== parking_table.csv =="
echo " -> ./parking_list/parking_table.csv"
cat ./driving_history/driving_history.csv
echo " ↓"
cat ./parking_list/parking_table.csv
echo ""
echo "== disp_main.csv=="
echo " -> ./disp_data/disp_main.csv"
cat ./disp_data/disp_main.csv
echo ""
echo "== drive_path_table.csv =="
echo " -> ./parking_list/drive_path_table.csv"
cat ./parking_list/drive_path_table.csv
