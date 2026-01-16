#!/bin/bash

# スクリプトの場所を取得
script_dir=$(dirname "$(realpath "$0")")

# sourceコマンドでパラメータファイルを読み込む
source "${script_dir}/param.sh"

## monitor_mainが何らかの理由で強制終了したとしても、再び起動させる
while true
do
    cd ${deptop}/make_data_hasuda01/
    command="ps aux | grep -ie \"./make_data_hasuda01\s\" | wc -l"
    kanshi=$(eval ${command})
    if [[ ${kanshi} -eq 0 ]]; then # commandの戻り値は文字ではあるが改行が含まれている場合を考慮し、-eqとして数値判定を実施。（文字0の後ろに改行が含まれていたとしても、-eqが自動的に数値のゼロに変換してくれる。）
        command="./make_data_hasuda01"
        eval ${command}
    fi
done

