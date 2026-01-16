#!/bin/bash
# tmux上で 「$ tmux kill-server」 とすれば全てプロセスを閉じてくれるが、一応スクリプトでも閉じれるように作っておいた

# セッションを閉じる
tmux kill-session -t make_data_hasuda01
tmux kill-session -t make_data_hasuda02
tmux kill-session -t make_data_hasuda03
tmux kill-session -t make_data_hasuda04
tmux kill-session -t make_data_hasuda05
tmux kill-session -t make_data_hasuda06




