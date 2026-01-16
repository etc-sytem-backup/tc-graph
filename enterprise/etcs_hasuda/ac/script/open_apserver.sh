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
tmux new-session -d -s ac
tmux new-session -d -s make_data_hasuda01
#tmux new-session -d -s monitor_avg
#tmux new-session -d -s monitor_table
#tmux new-session -d -s monitor_alert
tmux new-session -d -s develop
tmux new-session -d -s deploy

## ac
tmux send-keys -t ac 'cd ~/opt/aps/ac' C-m

## etc_deploy, 開発環境ディレクトリ参照用
tmux send-keys -t develop 'cd ~/sk_prj/enterprise/etcs/etcs_hasuda/' C-m
tmux send-keys -t deploy 'cd ~/sk_prj/enterprise/etcs/etcs_hasuda/etc_deploy/' C-m


