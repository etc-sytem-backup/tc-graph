#!/bin/bash

# tmuxセッションを作成
tmux new-session -d -s make_data_hasuda03
#tmux new-session -d -s monitor_avg
#tmux new-session -d -s monitor_table
#tmux new-session -d -s monitor_alert

##### 各セッションで特定のプログラムを実行
## display用データ作成
tmux send-keys -t make_data_hasuda03 '~/opt/aps/make_data_hasuda03/script/run_make_data_hasuda03.sh' C-m
#tmux send-keys -t monitor_avg '~/opt/aps/monitor_avg/script/monitor_avg.sh' C-m
#tmux send-keys -t monitor_table '~/opt/aps/monitor_table/script/monitor_table.sh' C-m
#tmux send-keys -t monitor_alert '~/opt/aps/monitor_alert/script/monitor_alert.sh' C-m



