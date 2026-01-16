#!/bin/bash

# 共通変数の読込
source ./param.sh 

# ~/opt/apsが存在していないならばapsディレクトリ作成
if [[ ! -d ${deptop} ]]; then

    # 存在しない場合は作成
    mkdir ${deptop}
fi

# ~/opt/aps/simulatorが存在していないならば作成する
if [[ ! -d ${deptop}/simulator ]]; then

    # 存在しない場合は作成
    mkdir ${deptop}/simulator
fi


# ~/opt/aps/simulator配下にディレクトリ作成
array=(ese-box)
for i in ${array[@]}; do
    if [[ ! -d ${deptop}/simulator/${i} ]]; then

        # 存在しない場合は作成
        mkdir ${deptop}/simulator/${i}
    fi
done

machine=ese-box
if [[ ! -e ${devdir}/ese-box/ese-box ]]; then   # コンパイル済みか？（実行ファイルが存在する？）
    echo "${machine} ${devdir}/ese-box/ Go Build!"
    cd ${devdir}/ese-box/
    go build .
else
    echo ""
    echo "${machine} don't Compile."
fi

echo "${devdir}/ese-box/* --> ${deptop}/simulator/${machine}/"
cp -rfp ${devdir}/ese-box/* ${deptop}/simulator/${machine}/                  # files copy
mv ${deptop}/simulator/${machine}/ese-box ${deptop}/simulator/${machine}/${machine}    # change bin name
rm -rf ${deptop}/simulator/${machine}/*.go                           # go file delete
rm -rf ${deptop}/simulator/${machine}/*.mod                          # go file delete
rm -rf ${deptop}/simulator/${machine}/*.sum                          # go file delete
rm -rf ${deptop}/simulator/${machine}/*.md                           # go file delete
rm -rf ${deptop}/simulator/${machine}/iniread                        # go dir delete
rm -rf ${deptop}/simulator/${machine}/readcsv                        # go dir delete
rm -rf ${deptop}/simulator/${machine}/disp_esebox                    # go dir delete

