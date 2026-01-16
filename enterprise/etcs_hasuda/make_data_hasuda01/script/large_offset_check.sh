#!/bin/bash
# 表示用端末（display）によるオフセット増減要求ファイルを検出し、戻り値を設定する。
# ポーリングするディレクトリは「../disp_data」とする。
# large_plus  : "0" を返す
# large_minus : "1" を返す
# 未検出      : "2" を返す

# チェックするディレクトリ
check_dir="../disp_data"

result="2"

# ファイル(large_plus)が存在するかチェック
if [ -f "${check_dir}/large_plus" ]; then
    rm "${check_dir}/large_plus"
    result="0"
fi

# ファイル(small_plus)が存在するかチェック
if [ -f "${check_dir}/large_minus" ]; then
    rm "${check_dir}/large_minus"
    result="1"
fi

echo ${result}

exit 0


