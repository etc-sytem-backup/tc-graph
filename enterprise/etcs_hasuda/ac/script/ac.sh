#!/bin/bash

# スクリプトの場所を取得
script_dir=$(dirname "$(realpath "$0")")

# sourceコマンドでパラメータファイルを読み込む
source "${script_dir}/param.sh"

# SBOX01〜04ディレクトリ内部の不要な要求ファイルを削除する
cd ${deptop}/sbox01/
rm -rf ./IcReq
rm -rf ./SIc
rm -rf ./SIH

cd ${deptop}/sbox02/
rm -rf ./IcReq
rm -rf ./SIc
rm -rf ./SIH

cd ${deptop}/sbox03/
rm -rf ./IcReq
rm -rf ./SIc
rm -rf ./SIH

cd ${deptop}/sbox04/
rm -rf ./IcReq
rm -rf ./SIc
rm -rf ./SIH

## acが何らかの理由で強制終了したとしても、再び起動させる
while true
do
    cd ${deptop}/ac/
    command="ps aux | grep -ie \"./ac\s\" | wc -l"
    kanshi=$(eval ${command})
    if [[ ${kanshi} -eq 0 ]]; then # commandの戻り値は文字ではあるが改行が含まれている場合を考慮し、-eqとして数値判定を実施。（文字0の後ろに改行が含まれていたとしても、-eqが自動的に数値のゼロに変換してくれる。）
        command="./ac"
        eval ${command}
    fi
done

## 旧コード
# while true
# do
#     cd ~/opt/aps/ac/
#     command="ps -aux | grep -e ac | wc -l"
#     kanshi=eval  ${command}
#     if [ ${kanshi}="1" ]; then
#         command="./ac"
#         eval ${command}
#     fi
# done
