#!/bin/bash

# tmuxセッションを作成
tmux new-session -d -s me01
tmux new-session -d -s me02
tmux new-session -d -s me03
tmux new-session -d -s ac
tmux new-session -d -s make_data_hasuda01
tmux new-session -d -s make_data_hasuda02
tmux new-session -d -s make_data_hasuda03
tmux new-session -d -s make_data_hasuda04
tmux new-session -d -s make_data_hasuda05
tmux new-session -d -s develop
tmux new-session -d -s deploy

##### 各セッションで特定のプログラムを実行
## etcs_tc
tmux send-keys -t me01 '~/opt/aps/me01/script/me01.sh' C-m
tmux send-keys -t me02 '~/opt/aps/me02/script/me02.sh' C-m
tmux send-keys -t me03 '~/opt/aps/me03/script/me03.sh' C-m

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


