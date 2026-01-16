#!/bin/bash

# tmuxセッションを作成
tmux new-session -d -s ac_shell
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
tmux new-session -d -s make_data_hasuda02
tmux new-session -d -s make_data_hasuda03
tmux new-session -d -s make_data_hasuda04
tmux new-session -d -s make_data_hasuda05
tmux new-session -d -s develop
tmux new-session -d -s deploy

##### 各セッションで特定のプログラムを実行
## main shell
tmux send-keys -t ac_shell 'cd ~/opt/aps/ac/' C-m

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

## ac
tmux send-keys -t ac '~/opt/aps/ac/script/ac.sh' C-m

## display用データ作成
tmux send-keys -t make_data_hasuda01 '~/opt/aps/make_data_hasuda01/script/run_make_data_hasuda01.sh' C-m
tmux send-keys -t make_data_hasuda02 '~/opt/aps/make_data_hasuda02/script/run_make_data_hasuda02.sh' C-m
tmux send-keys -t make_data_hasuda03 '~/opt/aps/make_data_hasuda03/script/run_make_data_hasuda03.sh' C-m
tmux send-keys -t make_data_hasuda04 '~/opt/aps/make_data_hasuda04/script/run_make_data_hasuda04.sh' C-m
tmux send-keys -t make_data_hasuda05 '~/opt/aps/make_data_hasuda05/script/run_make_data_hasuda05.sh' C-m

## etc_deploy, 開発環境ディレクトリ参照用
tmux send-keys -t develop 'cd ~/sk_prj/enterprise/etcs/etcs_hasuda/' C-m
tmux send-keys -t deploy 'cd ~/sk_prj/enterprise/etcs/etcs_hasuda/etc_deploy/' C-m


