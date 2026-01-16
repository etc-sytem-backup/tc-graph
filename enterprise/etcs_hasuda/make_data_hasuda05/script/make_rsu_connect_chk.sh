#!/bin/bash

## disp_dataディレクトリに、display通知用の「RSU回線切断通知」ファイルを作成する。

# ファイル作成path
path="../disp_data"

# RSU回線切断通知ファイル名
filename="rsu_connect_false"

# ディレクトリが存在するか確認し、存在しない場合は作成
if [ ! -d "$path" ]; then
  mkdir -p "$path"
fi

# ファイルを作成
touch "${path}/${filename}"
