kubectl config set-cluster {{cluster_name}} --certificate-authority={{kubeproxy_certificate_authority}} --embed-certs=true --server={{cluster_address}} --kubeconfig={{kubernetes_etc_dir}}/kube-proxy.kubeconfig
# 设置客户端认证参数
kubectl config set-credentials {{kubeproxy_CN}} --client-certificate={{kubeproxy_client_certificate}} --client-key={{kubeproxy_client_key}} --embed-certs=true --kubeconfig={{kubernetes_etc_dir}}/kube-proxy.kubeconfig
# 设置上下文参数
kubectl config set-context "{{kubeproxy_CN}}@{{cluster_name}}" --cluster={{cluster_name}} --user={{kubeproxy_CN}}  --kubeconfig={{kubernetes_etc_dir}}/kube-proxy.kubeconfig
# 选择默认上下文
kubectl config use-context "{{kubeproxy_CN}}@{{cluster_name}}" --kubeconfig={{kubernetes_etc_dir}}/kube-proxy.kubeconfig
echo "create_config.sh (kube-proxy) ok"