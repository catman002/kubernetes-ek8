#!/bin/bash
# 设置集群参数
kubectl config set-cluster {{cluster_name}} --certificate-authority={{kubectl_certificate_authority}} --embed-certs=true --server={{cluster_address}} --kubeconfig={{kubernetes_etc_dir}}/admin.kubeconfig
# 设置客户端认证参数
kubectl config set-credentials {{admin_CN}} --client-certificate={{kubectl_client_certificate}} --client-key={{kubectl_client_key}} --embed-certs=true --kubeconfig={{kubernetes_etc_dir}}/admin.kubeconfig
# 设置上下文参数
kubectl config set-context "{{admin_CN}}@{{cluster_name}}" --cluster={{cluster_name}}  --user={{admin_CN}}  --kubeconfig={{kubernetes_etc_dir}}/admin.kubeconfig
# 选择默认上下文
kubectl config use-context "{{admin_CN}}@{{cluster_name}}" --kubeconfig={{kubernetes_etc_dir}}/admin.kubeconfig


#生成config
mkdir -p  ~/.kube
\cp -f {{kubernetes_etc_dir}}/admin.kubeconfig ~/.kube/config
export KUBECONFIG=$HOME/.kube/config

#补全kubectl命令
source /usr/share/bash-completion/bash_completion
source <(kubectl completion bash)
kubectl completion bash > ~/.kube/completion.bash.inc
source '/root/.kube/completion.bash.inc'
source $HOME/.bash_profile

echo "create_config.sh (admin) ok"
