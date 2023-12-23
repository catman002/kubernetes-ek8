#!/bin/bash
function a { kubectl delete clusterrolebinding my-clusterrolebinding_system-kube-apiserver_kubernetes-admin;return 0;} && a
kubectl create my-clusterrolebinding_system-kube-apiserver_kubernetes-admin  \
--clusterrole=system:kubelet-api-admin --user {{admin_CN}}  \
--kubeconfig={{kubernetes_etc_dir}}/admin.kubeconfig
echo "create_clusterrolebinding.sh (admin) ok"

#kube-apiserver:kubelet-apis