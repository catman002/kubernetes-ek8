#!/bin/bash

source config.cfg

cd source/etcd && \
tar xvzf etcd-${etcd_version}-linux-amd64.tar.gz && \
systemctl daemon-reload && systemctl stop etcd && \
cd etcd-${etcd_version}-linux-amd64 && \
\cp -f etcd* ${bin_dir} && \
cd .. && \
rm -rf etcd-${etcd_version}-linux-amd64 && cd .. && cd ..
echo "install.sh(etcd) ok"