#!/usr/bin/bash

# 共通変数の読込
source ./param.sh 

# ~/opt/apsが存在していないならばapsディレクトリ作成
if [ ! -d ${deptop} ]; then

    # 存在しない場合は作成
    mkdir ${deptop}
fi

# ~/opt/aps/配下にディレクトリ作成
#array=(rsu01 rsu02 rsu03 rsu04 sbox01 sbox02 sbox03 sbox04 ac alert area_carcount gate_carcount reserve html_alert html_reserve html_parking html_parking_cnt)
array=(alert)
for i in ${array[@]}; do
    if [ ! -d ${deptop}/${i} ]; then

        # 存在しない場合は作成
        mkdir ${deptop}/${i}
    fi
done

# Car Alert(Reverse check) Copy
machine=alert
if [ ! -e ${devdir}/ac_alert/ac_alert ]; then   # コンパイル済みか？（実行ファイルが存在する？）
    echo "${machine} ${devdir}/ac_alert/ Go Build!"
    cd ${devdir}/ac_alert/
    go build .
else
    echo ""
    echo "${machine} don't Compile."
fi

echo "${devdir}/ac_alert/* --> ${deptop}/${machine}/"
cp -rfp ${devdir}/ac_alert/* ${deptop}/${machine}/                  # files copy
mv ${deptop}/${machine}/ac_alert ${deptop}/${machine}/${machine}    # change bin name
rm -rf ${deptop}/${machine}/*.go                           # go file delete
rm -rf ${deptop}/${machine}/*.mod                          # go file delete
rm -rf ${deptop}/${machine}/*.sum                          # go file delete
rm -rf ${deptop}/${machine}/*.md                           # go file delete
rm -rf ${deptop}/${machine}/iniread                        # go dir delete
rm -rf ${deptop}/${machine}/readcsv                        # go dir delete
