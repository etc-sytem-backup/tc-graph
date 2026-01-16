#!/bin/bash

# tmuxセッションを作成
tmux new-session -d -s develop

##### 各セッションで特定のプログラムを実行
## develop
tmux send-keys -t develop '~/sk_prj/enterprise/etcs/etcs_hasuda/' C-m



