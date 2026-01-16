#!/bin/bash

wc -l ./driving_history/driving_history.csv | awk '{print $1}'
