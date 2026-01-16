#!/bin/bash

### 共通パラメータ設定
deptop=~/opt/aps                     # homeディレクトリ直下のopt/apsにプログラムを配置する

# Check if 'gdate' command exists
# APServerにインストールされているdateがBSD系かGNU系かによって、aliasを設定する。
# Linuxだと内部的にGNU系のdateが採用されているはず。（MacOSだとBSD系）
if command -v gdate >/dev/null 2>&1; then
    shopt -s expand_aliases
    alias date=gdate
    alias sed=gsed
    echo "Alias set: 'date' 'sed' now refers to 'gdate' 'gsed'"
else
    echo "'gdate' 'gsed' command not found. 'date' and 'sed' will use the system default."
fi



