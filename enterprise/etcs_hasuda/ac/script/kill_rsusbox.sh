#!/bin/bash

# sbox,rsuのセッションを閉じる
tmux kill-session -t sbox01
tmux kill-session -t sbox02
tmux kill-session -t sbox03
tmux kill-session -t rsu01
tmux kill-session -t rsu02
tmux kill-session -t rsu03
tmux kill-session -t log01
tmux kill-session -t log02
tmux kill-session -t log03


