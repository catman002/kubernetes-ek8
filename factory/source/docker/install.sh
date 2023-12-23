#!/bin/bash
####################
#docker 安装
###################
#wget https://download.docker.com/linux/static/stable/x86_64/docker-18.06.3-ce.tgz
#tar zxf docker-18.06.3-ce.tgz
#cp docker/* /usr/bin/


source config.cfg
cd source/docker && \
tar xvzf docker-${docker_version}.tgz && \
\cp -f  docker/* $bin_dir && \
rm -rf docker && \
\cp -f docker-compose-Linux-x86_64 $bin_dir && \
cd .. && cd ..

echo "install.sh(docker) ok"
