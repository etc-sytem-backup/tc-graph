#!/bin/bash

# tmuxセッションを作成
tmux new-session -d -s me01
tmux new-session -d -s me02
tmux new-session -d -s me03

##### 各セッションで特定のプログラムを実行
## etcs_tc
tmux send-keys -t me01 '~/opt/aps/me01/script/me01.sh' C-m
tmux send-keys -t me02 '~/opt/aps/me02/script/me02.sh' C-m
tmux send-keys -t me03 '~/opt/aps/me03/script/me03.sh' C-m



