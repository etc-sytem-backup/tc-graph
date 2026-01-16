#!/bin/bash

parking_file="./parking_list/parking_table.csv"
cnt_file="./parking_list/parking_table_cnt.csv"

# 引数が車室数
car_max=$1
car_cnt=$(wc -l ${parking_file} | awk '{print $1}')

echo ${car_cnt}","${car_max} > ${cnt_file}


