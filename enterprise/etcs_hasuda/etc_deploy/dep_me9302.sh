#!/bin/bash

# 共通変数の読込
source ./param.sh 

# ~/opt/apsが存在していないならばapsディレクトリ作成
if [ ! -d ${deptop} ]; then

    # 存在しない場合は作成
    mkdir ${deptop}
fi

# ~/opt/aps/配下にディレクトリ作成
array=(me01 me02 me03 me04)
for i in ${array[@]}; do
    if [ ! -d ${deptop}/${i} ]; then

        # 存在しない場合は作成
        mkdir ${deptop}/${i}
    fi
done

machine=me9302_tc
if [ ! -e ${devdir}/${machine}/${machine} ]; then   # コンパイル済みか？（実行ファイルが存在する？）
    echo "${machine} ${devdir}/${machine}/ Go Build!"
    cd ${devdir}/${machine}/
    go build .
else
    echo ""
    echo "${machine} don't Compile."
fi

array=(me01 me02 me03 me04)
target=me9302_tc
for machine in ${array[@]}; do
    echo "${devdir}/${target}/* --> ${deptop}/${machine}/"
    cp -rfp ${devdir}/${target}/* ${deptop}/${machine}/                  # files copy
    mv ${deptop}/${machine}/${target} ${deptop}/${machine}/${machine}    # change bin name
    rm -rf ${deptop}/${machine}/*.go                           # go file delete
    rm -rf ${deptop}/${machine}/*.mod                          # go file delete
    rm -rf ${deptop}/${machine}/*.sum                          # go file delete
    rm -rf ${deptop}/${machine}/*.md                           # go file delete
    rm -rf ${deptop}/${machine}/conretry                       # go dir delete
    rm -rf ${deptop}/${machine}/findcmd                        # go dir delete
    rm -rf ${deptop}/${machine}/makecmd                        # go dir delete
    rm -rf ${deptop}/${machine}/tcpclient                      # go dir delete
    rm -rf ${deptop}/${machine}/iniread                        # go dir delete

    # 起動スクリプトはアプリケーションに対応したものを名前を変えてコピーする
    rm -rf ${deptop}/${machine}/script/*                                                  # デプロイ先のスクリプトディレクトリを空にする
    cp -fp ${devdir}/${target}/script/${machine}.sh ${deptop}/${machine}/script/          # 起動スクリプト
    cp -fp ${devdir}/${target}/script/param.sh ${deptop}/${machine}/script/               # スクリプト共通パラメータ
    cp -fp ${devdir}/${target}/script/get_receive_time.sh ${deptop}/${machine}/script/    # 通過時刻取得スクリプト
    cp -fp ${devdir}/${target}/script/make_csv_table.sh ${deptop}/${machine}/script/      # 通過履歴ファイルの一本化スクリプト
    cp -fp ${devdir}/${target}/script/make_wcn_table.sh ${deptop}/${machine}/script/      # WCN番号ファイルの一本化スクリプト
done
