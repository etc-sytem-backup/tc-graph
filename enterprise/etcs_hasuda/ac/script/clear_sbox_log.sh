#!/bin/bash

# SBOX01~04の./log/csvを空にする
mkdir -p ../dust
mkdir -p ../dust/sbox01
mkdir -p ../dust/sbox02
mkdir -p ../dust/sbox03
mkdir -p ../dust/sbox04

# dustディレクトリへ移動（ゴミ箱の様な使い方）
mv -f ../sbox01/log/csv/* ../dust/sbox01
mv -f ../sbox02/log/csv/* ../dust/sbox02
mv -f ../sbox03/log/csv/* ../dust/sbox03
mv -f ../sbox04/log/csv/* ../dust/sbox04




