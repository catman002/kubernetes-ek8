#!/bin/bash
####################
#containerd 安装
###################



source config.cfg
cd source/containerd && \
mkdir -p /etc/cni/net.d /opt/cni/bin && \
tar xf cni-plugins-linux-amd64-v${cni_plugins_version}.tgz -C /opt/cni/bin/ && \
tar -C / -xzf cri-containerd-cni-${containerd_version}-linux-amd64.tar.gz && \
#tar xf crictl-v${crictl_version}-linux-amd64.tar.gz -C /usr/local/bin/ && \
tar xf nerdctl-${nerdctl_version}-linux-amd64.tar.gz -C /usr/local/bin/ && \
tar xf buildkit-v${buildkit_version}.linux-amd64.tar.gz -C /usr/local/ && \
cd .. && cd ..
cat > /etc/crictl.yaml <<EOF
runtime-endpoint: unix:///run/containerd/containerd.sock
image-endpoint: unix:///run/containerd/containerd.sock
timeout: 10
debug: false
EOF
cat > /etc/systemd/system/buildkit.service <<EOF
[Unit]
Description=BuildKit
Documentation=https://github.com/moby/buildkit
[Service]
ExecStart=/usr/local/bin/buildkitd --oci-worker=false --containerd-worker=true
[Install]
WantedBy=multi-user.target
EOF
echo "install.sh(containerd) ok"
