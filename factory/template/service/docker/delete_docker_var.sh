#!/bin/bash
source config.cfg
rm -rf $docker_data_dir 2>.dockererr
if [[ $? -ne  0 ]];then
  cat .dockererr | while read line
  do
    line=${line#*\"}
    line=${line%\"*}
    umount $line
  done
  rm -rf .dockererr $docker_data_dir $docker_work_dir
fi
