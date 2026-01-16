#!/bin/bash

# tmuxセッションを作成
tmux new-session -d -s ac_shell

##### 各セッションで特定のプログラムを実行
## acディレクトリをカレントとして
tmux send-keys -t deploy '~/sk_prj/enterprise/etcs/etcs_hasuda/ac/' C-m



