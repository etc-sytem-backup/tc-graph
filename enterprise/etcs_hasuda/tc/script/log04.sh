#!/bin/bash

# スクリプトの場所を取得
script_dir=$(dirname "$(realpath "$0")")

# sourceコマンドでパラメータファイルを読み込む
source "${script_dir}/param.sh"

while true
do
    # cd ~/opt/aps/log04/
    cd ${deptop}/log04/
    command="ps aux | grep -ie \"log04\s\" | wc -l"
    kanshi=$(eval ${command})
    if [[ ${kanshi} -eq 0 ]]; then # commandの戻り値は文字ではあるが改行が含まれている場合を考慮し、-eqとして数値判定を実施。（文字0の後ろに改行が含まれていたとしても、-eqが自動的に数値のゼロに変換してくれる。）
        command="./log04 -s=\"LOG\" -i=\"${rsu04_ip_add}:${log04_port}\" -c=\"IH\" -n=1"
        eval ${command}
    fi
done

