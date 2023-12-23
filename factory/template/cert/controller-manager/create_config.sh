#!/bin/bash
# 设置集群参数
kubectl config set-cluster {{cluster_name}} --certificate-authority={{controller_certificate_authority}}  --embed-certs=true --server={{cluster_address}} --kubeconfig={{kubernetes_etc_dir}}/controller-manager.kubeconfig
# 设置客户端认证参数
kubectl config set-credentials {{controller_CN}}  --client-certificate={{controller_client_certificate}}  --client-key={{controller_client_key}}  --embed-certs=true --kubeconfig={{kubernetes_etc_dir}}/controller-manager.kubeconfig
# 设置上下文参数
kubectl config set-context "{{controller_CN}}@{{cluster_name}}" --cluster={{cluster_name}} --user={{controller_CN}}  --kubeconfig={{kubernetes_etc_dir}}/controller-manager.kubeconfig
# 选择默认上下文
kubectl config use-context "{{controller_CN}}@{{cluster_name}}" --kubeconfig={{kubernetes_etc_dir}}/controller-manager.kubeconfig
echo "create_config.sh (controller-manager) ok"