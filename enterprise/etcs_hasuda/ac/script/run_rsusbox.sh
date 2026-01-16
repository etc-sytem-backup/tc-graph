#!/bin/bash

# tmuxセッションを作成
tmux new-session -d -s sbox01
tmux new-session -d -s sbox02
tmux new-session -d -s sbox03
tmux new-session -d -s rsu01
tmux new-session -d -s rsu02
tmux new-session -d -s rsu03
tmux new-session -d -s log01
tmux new-session -d -s log02
tmux new-session -d -s log03

##### 各セッションで特定のプログラムを実行
## tc
tmux send-keys -t sbox01 '~/opt/aps/sbox01/script/sbox01.sh' C-m
tmux send-keys -t sbox02 '~/opt/aps/sbox02/script/sbox02.sh' C-m
tmux send-keys -t sbox03 '~/opt/aps/sbox03/script/sbox03.sh' C-m
tmux send-keys -t rsu01 '~/opt/aps/rsu01/script/rsu01.sh' C-m
tmux send-keys -t rsu02 '~/opt/aps/rsu02/script/rsu02.sh' C-m
tmux send-keys -t rsu03 '~/opt/aps/rsu03/script/rsu03.sh' C-m
tmux send-keys -t log01 '~/opt/aps/log01/script/log01.sh' C-m
tmux send-keys -t log02 '~/opt/aps/log02/script/log02.sh' C-m
tmux send-keys -t log03 '~/opt/aps/log03/script/log03.sh' C-m



