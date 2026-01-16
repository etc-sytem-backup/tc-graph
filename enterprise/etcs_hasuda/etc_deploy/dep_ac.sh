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
array=(ac)
for i in ${array[@]}; do
    if [ ! -d ${deptop}/${i} ]; then

        # 存在しない場合は作成
        mkdir ${deptop}/${i}
    fi
done

# Application Counter Copy
machine=ac
if [ ! -e ${devdir}/ac/application_counter ]; then   # コンパイル済みか？（実行ファイルが存在する？）
    echo "${machine} ${devdir}/ac/ Go Build!"
    cd ${devdir}/ac/
    go build .
else
    echo ""
    echo "${machine} don't Compile."
fi

echo "${devdir}/ac/* --> ${deptop}/${machine}/"
rm -rf ${deptop}/${machine}/*                                               # files clear(remove)
cp -rfp ${devdir}/ac/* ${deptop}/${machine}/                                # files copy
mv ${deptop}/${machine}/application_counter ${deptop}/${machine}/${machine} # change bin name
rm -rf ${deptop}/${machine}/*.go                                            # go file delete
rm -rf ${deptop}/${machine}/*.mod                                           # go file delete
rm -rf ${deptop}/${machine}/*.sum                                           # go file delete
rm -rf ${deptop}/${machine}/*.md                                            # go file delete
rm -rf ${deptop}/${machine}/iniread                                         # go dir delete
rm -rf ${deptop}/${machine}/csvcontroller                                   # go dir delete
