#!/bin/bash
cd `dirname $0`

app_dir=".."
cd $app_dir

while true
do
    command="ps aux | grep -E \"display\s\" | wc -l"
    kanshi=$(eval ${command})
    if [ ${kanshi} = "0" ]; then
        command="./display"
        eval ${command}
    fi
done