#!/bin/bash

source config.cfg

#keepalived
cd source/keepalived && \
tar xvzf keepalived-${keepalived_version}.tar.gz && \
cd keepalived-${keepalived_version} && \
 ./configure prefix=/usr/local/keepalived && \
make && make install  && \
\cp -f  /usr/local/keepalived/sbin/keepalived /${bin_dir} && \
cd .. && \
rm -rf keepalived-${keepalived_version} && \
cd .. &&  cd ..
echo "install.sh(keepalived) ok"
