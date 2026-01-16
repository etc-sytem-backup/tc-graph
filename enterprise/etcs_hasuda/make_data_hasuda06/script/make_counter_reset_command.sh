#!/bin/bash
# 断面交通量カウンターのリセット要求ファイルを作成する。

# 要求ファイルを配置するディレクトリ
check_dir="../disp_data"

# カウンターリセット要求ファイル作成
touch "${check_dir}/passage_count_reset"

exit 0


