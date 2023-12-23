#!/bin/bash

ver_str=`uname -r`
ver_array=(${ver_str//-/ })
ver=${ver_array[0]%.*}
echo ${ver}

if [[ $ver > 4.18 ]]; then
  echo $ver_str
else
  yum install https://www.elrepo.org/elrepo-release-8.el8.elrepo.noarch.rpm -y ; \
  sed -i "s@mirrorlist@#mirrorlist@g" /etc/yum.repos.d/elrepo.repo ; \
  sed -i "s@elrepo.org/linux@mirrors.tuna.tsinghua.edu.cn/elrepo@g" /etc/yum.repos.d/elrepo.repo ; \
  yum  --disablerepo="*"  --enablerepo="elrepo-kernel"  list  available -y ; \
  yum  --enablerepo=elrepo-kernel  install  kernel-ml -y ; \
  grubby --default-kernel
fi