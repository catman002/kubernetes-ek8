#!/bin/bash
nets=(`ip a |grep cali | awk '{print $2}' | awk -F@ '{print $1}'`)
for (( i=0;i<${#nets[@]};i++ )) do
  echo ${nets[i]}
  ip link delete ${nets[i]}
done