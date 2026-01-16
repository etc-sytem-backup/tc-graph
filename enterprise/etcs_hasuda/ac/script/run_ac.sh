#!/bin/bash

# tmuxセッションを作成
tmux new-session -d -s ac

##### 各セッションで特定のプログラムを実行
## ac
tmux send-keys -t ac '~/opt/aps/ac/script/ac.sh' C-m



