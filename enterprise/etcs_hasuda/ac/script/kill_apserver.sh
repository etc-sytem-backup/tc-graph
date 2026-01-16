#!/bin/bash
# tmux上で 「$ tmux kill-server」 とすれば全てプロセスを閉じてくれるが、一応スクリプトでも閉じれるように作っておいた

# セッションを閉じる
tmux kill-session -t ac_shell
tmux kill-session -t sbox01
tmux kill-session -t sbox02
tmux kill-session -t sbox03
tmux kill-session -t rsu01
tmux kill-session -t rsu02
tmux kill-session -t rsu03
tmux kill-session -t log01
tmux kill-session -t log02
tmux kill-session -t log03
tmux kill-session -t ac
tmux kill-session -t make_data_hasuda01
tmux kill-session -t make_data_hasuda02
tmux kill-session -t make_data_hasuda03
tmux kill-session -t make_data_hasuda04
tmux kill-session -t make_data_hasuda05
tmux kill-session -t make_data_hasuda06
tmux kill-session -t make_archive_hasuda
tmux kill-session -t develop
tmux kill-session -t deploy



