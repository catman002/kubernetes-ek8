#!/bin/bash

ver_str=`uname -r`
ver_array=(${ver_str//-/ })
ver=${ver_array[0]%.*}
echo ${ver}

if [[ $ver > 4.18 ]]; then
  echo $ver_str
else
  rpm --import https://www.elrepo.org/RPM-GPG-KEY-elrepo.org
  rpm -Uvh https://www.elrepo.org/elrepo-release-7.el7.elrepo.noarch.rpm
  yum --enablerepo=elrepo-kernel install -y kernel-ml
  awk -F \' '$1=="menuentry " {print i++ " : " $2}' /etc/grub2.cfg
  grub2-set-default 0
  grub2-editenv list
fi