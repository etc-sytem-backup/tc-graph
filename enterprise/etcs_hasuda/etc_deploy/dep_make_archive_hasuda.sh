#!/bin/bash

# 共通変数の読込
source ./param.sh 

# ~/opt/apsが存在していないならばapsディレクトリ作成
if [ ! -d ${deptop} ]; then

    # 存在しない場合は作成
    mkdir ${deptop}
fi

# ~/opt/aps/配下にディレクトリ作成
#array=(rsu01 rsu02 rsu03 rsu04 sbox01 sbox02 sbox03 sbox04 ac alert area_carcount gate_carcount reserve html_alert html_reserve html_parking html_parking_cnt)
array=(make_archive_hasuda)
for i in ${array[@]}; do
    if [ ! -d ${deptop}/${i} ]; then

        # 存在しない場合は作成
        mkdir ${deptop}/${i}
    fi
done

# make_archive_hasuda Copy
machine=make_archive_hasuda
if [ ! -e ${devdir}/make_archive_hasuda/${machine} ]; then   # コンパイル済みか？（実行ファイルが存在する？）
echo "${machine} ${devdir}/make_archive_hasuda/ Go Build!"
    cd ${devdir}/make_archive_hasuda/
    go build .
else
    echo ""
    echo "${machine} don't Compile."
fi

echo "${devdir}/make_archive_hasuda/* --> ${deptop}/${machine}/"
cp -rfp ${devdir}/make_archive_hasuda/* ${deptop}/${machine}/                  # files copy
rm -rf ${deptop}/${machine}/*.go                           # go file delete
rm -rf ${deptop}/${machine}/*.mod                          # go file delete
rm -rf ${deptop}/${machine}/*.sum                          # go file delete
rm -rf ${deptop}/${machine}/*.md                           # go file delete
rm -rf ${deptop}/${machine}/iniread                        # go dir delete
rm -rf ${deptop}/${machine}/readcsv                        # go dir delete



