#!/bin/bash
source config.cfg
rm -rf $kubelet_work_dir 2>.kubeleterr
if [[ $? -ne  0 ]];then
  cat .kubeleterr | while read line
  do
    line=${line#*\"}
    line=${line%\"*}
    umount $line
  done
  rm -rf .kubeleterr $kubelet_work_dir
fi