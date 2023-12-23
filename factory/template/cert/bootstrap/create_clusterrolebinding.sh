#!/bin/bash
#参考资料：
#https://segmentfault.com/a/1190000041869038?utm_source=sf-similar-article
#https://www.cnblogs.com/normanlin/p/10909703.html
#https://www.jianshu.com/p/bb973ab1029b

#允许 kubelet-bootstrap 用户创建首次启动的 CSR 请求 和RBAC授权规则
function a { kubectl delete clusterrolebinding my-clusterrolebinding_kubelet-bootstrap_system-node-bootstrapper;return 0;} && a
kubectl create clusterrolebinding my-clusterrolebinding_kubelet-bootstrap_system-node-bootstrapper --clusterrole=system:node-bootstrapper --user={{bootstrap_user}}
function a { kubectl delete clusterrolebinding my-clusterrolebinding_kubelet-bootstrap_cluster-admin;return 0;} && a
kubectl create clusterrolebinding my-clusterrolebinding_kubelet-bootstrap_cluster-admin --clusterrole=cluster-admin --user={{bootstrap_user}}
function a { kubectl delete  clusterrolebinding my-clusterrolebinding_system-node_system-node;return 0;} && a
kubectl create clusterrolebinding my-clusterrolebinding_system-node_system-node --clusterrole=system:node --group=system:nodes

# 自动批准 system:bootstrappers 组用户 TLS bootstrapping 首次申请证书的 CSR 请求
function a { kubectl delete clusterrolebinding my-clusterrolebinding_system-bootstrappers_nodeclient;return 0;} && a
kubectl create clusterrolebinding my-clusterrolebinding_system-bootstrappers_nodeclient --clusterrole=system:certificates.k8s.io:certificatesigningrequests:nodeclient --group={{bootstrap_group}}

# 自动批准 system:nodes 组用户更新 kubelet 自身与 apiserver 通讯证书的 CSR 请求
function a { kubectl delete my-clusterrolebinding_system-nodes_selfnodeclient;return 0;} && a
kubectl create clusterrolebinding my-clusterrolebinding_system-nodes_selfnodeclient --clusterrole=system:certificates.k8s.io:certificatesigningrequests:selfnodeclient --group=system:nodes

# 自动批准 system:nodes 组用户更新 kubelet 10250 api 端口证书的 CSR 请求
function a { kubectl delete my-clusterrolebinding_system-nodes_selfnodeserver;return 0;} && a
kubectl create clusterrolebinding my-clusterrolebinding_system-nodes_selfnodeserver --clusterrole=system:certificates.k8s.io:certificatesigningrequests:selfnodeserver --group=system:nodes
echo "create_clusterrolebinding.sh(bootstrap) ok"