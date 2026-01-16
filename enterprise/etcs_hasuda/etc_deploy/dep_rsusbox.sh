#!/bin/bash

# 共通変数の読込
source ./param.sh 

# ~/opt/apsが存在していないならばapsディレクトリ作成
if [[ ! -d ${deptop} ]]; then

    # 存在しない場合は作成
    mkdir ${deptop}
fi

# ~/opt/aps/配下にディレクトリ作成
array=(rsu01 rsu02 rsu03 rsu04 sbox01 sbox02 sbox03 sbox04 log01 log02 log03 log04)
for i in ${array[@]}; do
    if [[ ! -d ${deptop}/${i} ]]; then

        # 存在しない場合は作成
        mkdir ${deptop}/${i}
    fi
done

if [[ ! -e ${devdir}/tc/traffic_counter ]]; then   # コンパイル済みか？（実行ファイルが存在する？）
    echo "${machine} ${devdir}/tc/ Go Build!"
    cd ${devdir}/tc/
    go build .
else
    echo ""
    echo "${machine} don't Compile."
fi

array=(rsu01 rsu02 rsu03 rsu04 sbox01 sbox02 sbox03 sbox04 log01 log02 log03 log04)
for machine in ${array[@]}; do
    echo "${devdir}/tc/* --> ${deptop}/${machine}/"
    cp -rfp ${devdir}/tc/* ${deptop}/${machine}/                                # files copy
    mv ${deptop}/${machine}/traffic_counter ${deptop}/${machine}/${machine}        # change bin name
    rm -rf ${deptop}/${machine}/*.go                                               # go file delete
    rm -rf ${deptop}/${machine}/*.mod                                              # go file delete
    rm -rf ${deptop}/${machine}/*.sum                                              # go file delete
    rm -rf ${deptop}/${machine}/*.md                                               # go file delete
    rm -rf ${deptop}/${machine}/conretry                                           # go dir delete
    rm -rf ${deptop}/${machine}/findcmd                                            # go dir delete
    rm -rf ${deptop}/${machine}/iniread                                            # go dir delete
    rm -rf ${deptop}/${machine}/makecmd                                            # go dir delete
    rm -rf ${deptop}/${machine}/tcpclient                                          # go dir delete

                                                                                   # 起動スクリプトはアプリケーションに対応したものを名前を変えてコピーする
    rm -rf ${deptop}/${machine}/script/*                                           # デプロイ先のスクリプトディレクトリを空にする
    cp -fp ${devdir}/tc/script/${machine}.sh ${deptop}/${machine}/script/       # 起動スクリプト
    cp -fp ${devdir}/tc/script/param.sh ${deptop}/${machine}/script/            # 共通パラメータ
    cp -fp ${devdir}/tc/script/make_csv_table.sh ${deptop}/${machine}/script/   # 通過履歴ファイルの一本化スクリプト
    cp -fp ${devdir}/tc/script/make_wcn_table.sh ${deptop}/${machine}/script/   # WCNファイルの一本化スクリプト
    cp -fp ${devdir}/tc/script/get_receive_time.sh ${deptop}/${machine}/script/ # WCN_rireki_a?.csvに登録されているWCNの、検出時刻を返すスクリプト
done



