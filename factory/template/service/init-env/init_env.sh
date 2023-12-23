#!/bin/bash

## graceful shutdown
trap "exit 0;" SIGTERM SIGINT

source config.cfg

function echo_s() {
 echo -e "\033[32m[info]\033[0m $1"
}

function echo_f() {
 echo -e "\033[31m[fail]\033[0m $1"
}

function checkR () {
 if [ $? -ne 0 ];then
    echo -e "\033[31m[fail] \033[0m $1 fail !"
    exit 1
  fi
}


##################
##default config
##################
core_ver_str=`uname -r`
core_ver_array=(${core_ver_str//-/ })
core_ver=${core_ver_array[0]%.*}

#####################
#init all server env
####################
#设置时间同步
if [[ $target_system == "centos7" || $target_system == "centos8" ]];then
ntpdate ntpdate ntp3.aliyun.com
echo "0 */6 * * * root ntpdate ntp.aliyun.com" >/etc/crontab
cp  /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
fi


#关闭防火墙
systemctl stop firewalld
systemctl disable --now firewalld

#关闭selinux
setenforce 0  >/dev/null 2>&1
sed -i "s/IPV6INIT=\"yes\"/IPV6INIT=\"no\"/" /etc/sysconfig/network-scripts/ifcfg-${interface}
sed -i "s/SELINUX=permissive/SELINUX=disabled/" /etc/sysconfig/selinux
sed -i "s/SELINUX=enforcing/SELINUX=disabled/" /etc/sysconfig/selinux

#关闭交换分区
sed -ri 's/.*swap.*/#&/' /etc/fstab
swapoff -a
sysctl -w vm.swappiness=0

#配置ulimit
ulimit_size=655350
ulimit -SHn $ulimit_size
cat > /etc/security/limits.conf <<EOF
* soft nofile $ulimit_size
* hard nofile $[$ulimit_size*2]
* soft nproc $ulimit_size
* hard nproc $ulimit_size
* soft memlock unlimited
* hard memlock unlimited
EOF
sysctl -p /etc/security/limits.conf

#ipvs.modules
#4.18以下使用nf_conntrack_ipv4即可
cat > /etc/sysconfig/modules/ipvs.modules <<EOF
#!/bin/bash
modprobe -- ip_vs
modprobe -- ip_vs_rr
modprobe -- ip_vs_wrr
modprobe -- ip_vs_sh
modprobe -- br_netfilter
modprobe -- $conntrack_str
modprobe -- ip_vs_lc
modprobe -- ip_vs_wlc
modprobe -- ip_vs_lblc
modprobe -- ip_vs_lblcr
modprobe -- ip_vs_dh
modprobe -- ip_vs_nq
modprobe -- ip_vs_sed
modprobe -- ip_vs_ftp
modprobe -- ip_vs_fo
EOF
chmod 755 /etc/sysconfig/modules/ipvs.modules
bash /etc/sysconfig/modules/ipvs.modules
modprobe br_netfilter

#安装ipvsadm
#4.18以下使用nf_conntrack_ipv4即可
cat > /etc/modules-load.d/ipvs.conf <<EOF
ip_vs
ip_vs_lc
ip_vs_wlc
ip_vs_rr
ip_vs_wrr
ip_vs_lblc
ip_vs_lblcr
ip_vs_dh
ip_vs_sh
ip_vs_fo
ip_vs_nq
ip_vs_sed
ip_vs_ftp
ip_vs_sh
$conntrack_str
ip_tables
ip_set
xt_set
ipt_set
ipt_rpfilter
ipt_REJECT
ipip
EOF
#containerd需要的模块
cat <<EOF | sudo tee /etc/modules-load.d/containerd.conf
overlay
br_netfilter
EOF
systemctl restart systemd-modules-load.service
systemctl enable systemd-modules-load.service
lsmod | grep -e ip_vs -e nf_conntrack

#修改内核参数
cat  > /etc/sysctl.d/k8s.conf << EOF
net.ipv4.ip_forward = 1
net.bridge.bridge-nf-call-iptables = 1
net.bridge.bridge-nf-call-ip6tables = 1
#fs.may_detach_mounts = 1
vm.overcommit_memory=1
vm.panic_on_oom=0
fs.inotify.max_user_watches=89100
fs.file-max=52706963
fs.nr_open=52706963
net.netfilter.nf_conntrack_max=2310720
net.ipv4.tcp_keepalive_time = 600
net.ipv4.tcp_keepalive_probes = 3
net.ipv4.tcp_keepalive_intvl =15
net.ipv4.tcp_max_tw_buckets = 36000
net.ipv4.tcp_tw_reuse = 1
net.ipv4.tcp_max_orphans = 327680
net.ipv4.tcp_orphan_retries = 3
net.ipv4.tcp_syncookies = 1
net.ipv4.tcp_max_syn_backlog = 16384
#net.ipv4.ip_conntrack_max = 65536
net.ipv4.tcp_max_syn_backlog = 16384
net.ipv4.tcp_timestamps = 0
#net.ipv4.tcp_tw_recycle=0
net.core.somaxconn = 16384
vm.swappiness=0
net.ipv6.conf.all.disable_ipv6=1
vm.max_map_count=262144
EOF
sysctl -p /etc/sysctl.d/k8s.conf

#修改 sysctl.conf
cat >/etc/sysctl.conf  << EOF
net.core.somaxconn= 1024
EOF
sysctl -p

#修改network
#cat >/etc/sysconfig/network << EOF
#NETWORKING_IPV6=no
#EOF
#sysctl -p /etc/sysconfig/network
sysctl --system

#修该sshd
sed -i "s/#   StrictHostKeyChecking ask/StrictHostKeyChecking no/g"  /etc/ssh/ssh_config
sed -i "s/#UseDNS yes/UseDNS no/g" /etc/ssh/sshd_config
systemctl restart sshd

#修改源
echo >/etc/yum.repos.d/kubernetes.repo << EOF
[kubernetes]
name=Kubernetes
baseurl=http://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64/
enabled=1
gpgcheck=0
EOF
yum makecache fast

#其他
echo 1 > /proc/sys/net/bridge/bridge-nf-call-iptables
function ai { systemctl disable ip6tables.service;return 0; } && ai

echo "init_env.sh ok"


