#!/bin/bash
####################
#docker 安装
###################
#wget https://download.docker.com/linux/static/stable/x86_64/docker-18.06.3-ce.tgz
#tar zxf docker-18.06.3-ce.tgz
#cp docker/* /usr/bin/


source config.cfg
function a { yum erase -y libseccomp;return 0; } && a;
cd source/libseccomp && \
rpm -ivh libseccomp-${libseccomp_version}.el8.x86_64.rpm && \
cd ..
echo "install.sh(libseccomp) ok"
