#!/bin/bash

# スクリプトの場所を取得
script_dir=$(dirname "$(realpath "$0")")

# sourceコマンドでパラメータファイルを読み込む
source "${script_dir}/param.sh"


while true
do
    # cd ~/opt/aps/rsu02/
    cd ${deptop}/rsu02/
    command="ps aux | grep -ie \"rsu02\s\" | wc -l"
    kanshi=$(eval ${command})
    if [[ ${kanshi} -eq 0 ]]; then # commandの戻り値は文字ではあるが改行が含まれている場合を考慮し、-eqとして数値判定を実施。（文字0の後ろに改行が含まれていたとしても、-eqが自動的に数値のゼロに変換してくれる。）
        command="./rsu02 -s=\"RSU\" -i=\"${rsu02_ip_add}:${rsu02_port}\" -c=\"IH\" -n=1"
        eval ${command}
    fi
done

