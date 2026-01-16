#!/bin/bash

point=$1

awk -F',' '{print $5}' ./driving_history/driving_history.csv | sed -n ${point}


