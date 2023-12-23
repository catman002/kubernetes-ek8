#!/bin/bash
# 设置集群参数
kubectl config set-cluster {{cluster_name}} --certificate-authority={{scheduler_certificate_authority}} --embed-certs=true --server={{cluster_address}} --kubeconfig={{kubernetes_etc_dir}}/scheduler.kubeconfig
# 设置客户端认证参数
kubectl config set-credentials "{{scheduler_CN}}" --client-certificate={{scheduler_client_certificate}} --client-key={{scheduler_client_key}} --embed-certs=true --kubeconfig={{kubernetes_etc_dir}}/scheduler.kubeconfig
# 设置上下文参数
kubectl config set-context "{{scheduler_CN}}@{{cluster_name}} " --cluster={{cluster_name}}  --user="{{scheduler_CN}}"  --kubeconfig={{kubernetes_etc_dir}}/scheduler.kubeconfig
# 选择默认上下文
kubectl config use-context "{{scheduler_CN}}@{{cluster_name}} " --kubeconfig={{kubernetes_etc_dir}}/scheduler.kubeconfig
echo "create_config.sh (scheduler) ok"