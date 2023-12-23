
source ../config.cfg

calico_flag=0
coredns_flag=0
docker_flag=0
etcd_flag=0
haproxy_flag=0
lua_flag=0
keepalived_flag=0
pause_flag=0

#calico
if [[ $calico_flag -eq 1 ]]; then
  docker pull calico/cni:v${calico_version}
  docker pull calico/pod2daemon-flexvol:v${calico_version}
  docker pull calico/node:v${calico_version}
  docker pull calico/kube-controllers:v${calico_version}

  mkdir -p calico/$calico_version
  docker save -o calico/$calico_version/cni.rar calico/cni:v${calico_version}
  docker save -o calico/$calico_version/pod2daemon-flexvol.rar calico/pod2daemon-flexvol:v${calico_version}
  docker save -o calico/$calico_version/node.rar calico/node:v${calico_version}
  docker save -o calico/$calico_version/kube-controllers.rar calico/kube-controllers:v${calico_version}

  docker rmi calico/cni:v${calico_version}
  docker rmi calico/pod2daemon-flexvol:v${calico_version}
  docker rmi calico/node:v${calico_version}
  docker rmi calico/kube-controllers:v${calico_version}
fi

#coredns
if [[ $coredns_flag -eq 1 ]]; then
  docker pull coredns/coredns:$coredns_version
  docker save -o coredns/coredns_$coredns_version.rar coredns/coredns:$coredns_version
  docker rmi coredns/coredns:$coredns_version
fi

#docker
if [[ $docker_flag -eq 1 ]]; then
  wget https://download.docker.com/linux/static/stable/x86_64/docker-${docker_version}.tgz -O docker/docker-${docker_version}.tgz
fi

#etcd
if [[ $etcd_flag -eq 1 ]]; then
  wget https://github.com/etcd-io/etcd/releases/download/${etcd_version}/etcd-${etcd_version}-linux-amd64.tar.gz -O etcd/etcd-${etcd_version}-linux-amd64.tar.gz 
fi

#haproxy
if [[ $haproxy_flag -eq 1 ]]; then
  wget https://github.com/haproxy/haproxy/archive/refs/tags/v${haproxy_version}.tar.gz haproxy/haproxy-${haproxy_version}.tar.gz
fi

#lua
if [[ $lua_flag -eq 1 ]]; then
  wget https://github.com/lua/lua/archive/refs/tags/v${lua_version}.tar.gz -O haproxy/lua-${lua_version}.tar.gz
fi

#keepalived
if [[ $keepalived_flag -eq 1 ]]; then
  wget https://github.com/acassen/keepalived/archive/refs/tags/v${keepalived_version}.tar.gz -O keepalived/keepalived-${keepalived_version}.tar.gz
fi

#pause
if [[ $pause_flag -eq 1 ]]; then
  docker pull registry.aliyuncs.com/google_containers/pause:$pause_version
  docker save -o pause/pause_$pause_version.rar registry.aliyuncs.com/google_containers/pause:$pause_version
  docker rmi registry.aliyuncs.com/google_containers/pause:$pause_version
fi
