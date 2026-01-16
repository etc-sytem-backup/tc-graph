#!/bin/bash

### ./tc_wcn/に配置されているWCN番号ファイルを一覧ファイルに１本化する。
### すでに一覧ファイルに同じWCNの内容が含まれている場合は追加しない。（データを入れ替えて更新したほうが良い？）

## 通過履歴ファイル群の一覧を配列に取得する
files="./tc_wcn/*"                         # WCN番号ファイル格納ディレクトリ
out_file="./tc_wcn_table/WCN_table.csv"   # WCN番号一覧ファイル

# 一覧ファイルが存在しない場合は空のファイルを作る
if [[ ! -e ${out_file} ]]; then
    touch ${out_file}
fi

fileary=()  # ファイル名格納用
# dirary=() # ディレクトリ名格納用
for filepath in $files; do
    if [[ -f $filepath ]] ; then
        fileary+=("$filepath")
        #  elif [ -d $filepath ] ; then   # ファイルではなくディレクトリだった場合
        #    dirary+=("$filepath")
    fi
done

## 取得したファイルごとに処理
##   WCNが一覧ファイルに含まれているか？
##     含まれている → 何もせず、履歴ファイルを消す
##     含まれてない → 一覧ファイルに追記して、一覧ファイルをソートする
##                    ※一番左のカラムが日時を含む形式になっているので、ソートすると時系列に並ぶ
for i in ${fileary[@]}; do
    echo $i
    wcn=$(awk -F ',' '{print $1}' ./$i)
    result=$(cat ./${out_file} | grep -c ${wcn})  # WCN番号が、一覧ファイルに含まれている？
    echo "${wcn}$ => {result}"
    if [[ ${result} = "0"  ]]; then         # 一覧ファイルにまだ存在していないwcn番号だった

        # 一覧ファイルにWCN番号を追記してソートする
        cat ./$i >> ${out_file}
        cp ${out_file} ./tmp.csv
        cat ./tmp.csv | sort -t, -k1,1n > ${out_file}   # 区切り文字は「,」で、カラム１(日付自分秒マイクロ)を数値とみなして保存する
        rm -rf ./tmp.csv
    fi
    rm -rf ./$i  # チェックしたWCN番号ファイルを削除する（将来的にこれを削除ではなく移動にするとバックアップとなる）
done


