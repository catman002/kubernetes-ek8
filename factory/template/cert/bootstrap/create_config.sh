#!/bin/bash
# 设置集群参数
kubectl config --kubeconfig={{kubernetes_etc_dir}}/kubelet-bootstrap.kubeconfig set-cluster {{cluster_name}} --server={{cluster_address}} --certificate-authority={{ca_file}} --embed-certs=true
kubectl config --kubeconfig={{kubernetes_etc_dir}}/kubelet-bootstrap.kubeconfig set-credentials {{bootstrap_user}} --token="{{bootstrap_token}}"
kubectl config --kubeconfig={{kubernetes_etc_dir}}/kubelet-bootstrap.kubeconfig set-context bootstrap --user={{bootstrap_user}} --cluster={{cluster_name}}
kubectl config --kubeconfig={{kubernetes_etc_dir}}/kubelet-bootstrap.kubeconfig use-context bootstrap
echo "create_config.sh (bootstrap) ok"