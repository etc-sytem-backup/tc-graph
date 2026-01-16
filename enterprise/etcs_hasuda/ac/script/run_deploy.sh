#!/bin/bash

# tmuxセッションを作成
tmux new-session -d -s deploy

##### 各セッションで特定のプログラムを実行
## deploy
tmux send-keys -t deploy '~/sk_prj/enterprise/etcs/etcs_hasuda/etc_deploy/' C-m



