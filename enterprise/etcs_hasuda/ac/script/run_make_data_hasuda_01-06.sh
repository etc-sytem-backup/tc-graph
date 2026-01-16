#!/bin/bash

# tmuxセッションを作成
tmux new-session -d -s make_data_hasuda01
tmux new-session -d -s make_data_hasuda02
tmux new-session -d -s make_data_hasuda03
tmux new-session -d -s make_data_hasuda04
tmux new-session -d -s make_data_hasuda05
tmux new-session -d -s make_data_hasuda06


##### 各セッションで特定のプログラムを実行
## display用データ作成
tmux send-keys -t make_data_hasuda01 '~/opt/aps/make_data_hasuda01/script/run_make_data_hasuda01.sh' C-m
tmux send-keys -t make_data_hasuda02 '~/opt/aps/make_data_hasuda02/script/run_make_data_hasuda02.sh' C-m
tmux send-keys -t make_data_hasuda03 '~/opt/aps/make_data_hasuda03/script/run_make_data_hasuda03.sh' C-m
tmux send-keys -t make_data_hasuda04 '~/opt/aps/make_data_hasuda04/script/run_make_data_hasuda04.sh' C-m
tmux send-keys -t make_data_hasuda05 '~/opt/aps/make_data_hasuda05/script/run_make_data_hasuda05.sh' C-m
tmux send-keys -t make_data_hasuda06 '~/opt/aps/make_data_hasuda06/script/run_make_data_hasuda06.sh' C-m


