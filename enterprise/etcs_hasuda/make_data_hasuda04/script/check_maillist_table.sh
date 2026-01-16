#!/bin/bash
## メール送信済みリストに対し、すでに同じデータが登録されていないかチェックする。
## 同じデータが登録されている場合は、戻り値として文字列"Hit"返す

# 引数の取得
check_string=$1

# 出力ファイル名の定義
target_file="./parking_list/maillist_table.csv"

# メール送信済みリストの存在チェックと作成
if [ ! -d "$(dirname "$target_file")" ]; then
  mkdir -p "$(dirname "$target_file")"
fi

# 変数check_stringと同じデータが、すでにtarget_fileに存在している場合は、戻り値"Hit"。
# 存在していない場合は、戻り値が空文字
if grep -Fxq "$check_string" "$target_file"
then
    echo "Hit"
    exit 0
fi

echo ""
exit 0

