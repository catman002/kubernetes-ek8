#!/bin/bash
libseccomp_version=2.5.1-1
lua_version=5.4.3
haproxy_version=2.4.0
keepalived_version=2.2.2

etcd_version=v3.5.9
# in calicalico.yaml,PodDisruptionBudget:policy/v1beta1 -> policy/v1 (k8s 1.21+)
calico_version=3.26.4
pause_version=3.8
coredns_version=1.9.4

##https://github.com/containernetworking/plugins 下载cni-plugins
docker_version=20.10.9
containerd_version=1.7.10
cni_plugins_version=1.3.0
crictl_version=1.24.2
nerdctl_version=1.7.1
buildkit_version=0.12.4
openssl_version=1_1_0j
#runtime: docker/containerd
runtime=containerd
