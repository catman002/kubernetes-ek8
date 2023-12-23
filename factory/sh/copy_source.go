package sh

import (
	"ek8/common"
	"ek8/config"
	"ek8/shell"
)

func Copy_source() {
	if common.ModFlags["source"] != 1 {
		return
	}
	common.TPrintln("拷贝源文件...")
	target_dir := config.Cfg["release_dir"] + "/source"
	cmd := new(common.Command).
		Add("source config.cfg").
		//Add("rm -rf " + target_dir).
		Add("mkdir -p " + target_dir).
		Add("mkdir -p " + target_dir + "/libseccomp").
		Add("\\cp -Rf factory/source/libseccomp/libseccomp-${libseccomp_version}.el8.x86_64.rpm  " + target_dir + "/libseccomp/").
		Add("\\cp -Rf factory/source/libseccomp/install.sh " + target_dir + "/libseccomp/").
		Add("mkdir -p " + target_dir + "/registry").
		Add("\\cp -Rf factory/source/registry/registry.rar  " + target_dir + "/registry/").
		Add("mkdir -p " + target_dir + "/calico").
		Add("\\cp -Rf factory/source/calico/$calico_version  " + target_dir + "/calico/").
		Add("mkdir -p " + target_dir + "/coredns").
		Add("\\cp -f factory/source/coredns/coredns_$coredns_version.rar " + target_dir + "/coredns/").
		Add("\\cp -f factory/source/coredns/dnsutils.rar " + target_dir + "/coredns/").
		Add("mkdir -p " + target_dir + "/docker").
		Add("\\cp -f factory/source/docker/docker-$docker_version.tgz " + target_dir + "/docker/").
		Add("\\cp -f factory/source/docker/install.sh " + target_dir + "/docker/").
		Add("\\cp -f factory/source/docker/docker-compose-Linux-x86_64 " + target_dir + "/docker/").
		Add("mkdir -p " + target_dir + "/containerd").
		Add("\\cp -f factory/source/containerd/crictl-v${crictl_version}-linux-amd64.tar.gz " + target_dir + "/containerd/").
		Add("\\cp -f factory/source/containerd/nerdctl-${nerdctl_version}-linux-amd64.tar.gz " + target_dir + "/containerd/").
		Add("\\cp -f factory/source/containerd/buildkit-v${buildkit_version}.linux-amd64.tar.gz " + target_dir + "/containerd/").
		Add("\\cp -f factory/source/containerd/cni-plugins-linux-amd64-v${cni_plugins_version}.tgz " + target_dir + "/containerd/").
		Add("\\cp -f factory/source/containerd/cri-containerd-cni-${containerd_version}-linux-amd64.tar.gz " + target_dir + "/containerd/").
		Add("\\cp -f factory/source/containerd/install.sh " + target_dir + "/containerd/").
		Add("mkdir -p " + target_dir + "/etcd").
		Add("\\cp -f factory/source/etcd/etcd-$etcd_version-linux-amd64.tar.gz " + target_dir + "/etcd/").
		Add("\\cp -f factory/source/etcd/install.sh " + target_dir + "/etcd/").
		Add("mkdir -p " + target_dir + "/openssl").
		Add("\\cp -f factory/source/openssl/OpenSSL_${openssl_version}.tar.gz " + target_dir + "/openssl/").
		Add("\\cp -f factory/source/openssl/install.sh " + target_dir + "/openssl/").
		Add("mkdir -p " + target_dir + "/haproxy").
		Add("\\cp -f factory/source/haproxy/haproxy-$haproxy_version.tar.gz " + target_dir + "/haproxy/").
		Add("\\cp -f factory/source/haproxy/lua-$lua_version.tar.gz " + target_dir + "/haproxy/").
		Add("\\cp -f factory/source/haproxy/install.sh " + target_dir + "/haproxy/").
		Add("mkdir -p " + target_dir + "/keepalived").
		Add("\\cp -f factory/source/keepalived/keepalived-$keepalived_version.tar.gz " + target_dir + "/keepalived/").
		Add("\\cp -f factory/source/keepalived/install.sh " + target_dir + "/keepalived/").
		Add("mkdir -p " + target_dir + "/pause").
		Add("\\cp -f factory/source/pause/pause_$pause_version.rar " + target_dir + "/pause/").
		Add("mkdir -p " + target_dir + "/kubernetes/$target_cpu").
		Add("\\cp -Rf factory/source/kubernetes/$target_cpu/$kube_version " + target_dir + "/kubernetes/$target_cpu").
		Add("\\cp -Rf factory/source/kubernetes/install.sh " + target_dir + "/kubernetes/").
		Add("echo source 拷贝完成").
		ToString()
	state, _ := shell.Exec(cmd, "local")
	if !state {
		common.EPrintln(cmd)
		panic("拷贝源文件出错！")

	}
}
