# kubernetes-ek8 v2.0.0.0 （支持k8s v1.29.0 和v1.28.4）
ek8是一款k8s集群快速部署工具,目标是实现k8s的高可用集群 快捷安装、简易维护。

#### （一）ek8作为k8s集群安装管理工具，具有以下特点：
1. 具有生产和部署2种功能。

``` 生产功能通过命令创建支持不同版本的k8s部署工具，快速集成kubernetes新版本；生产和部署功能集成于“ek8”命令； ```

2. 安装便捷：
* 通过一台个人电脑作为安装机，配置config.cfg后，执行ek8命令，即可完成集群安装；
*  服务求无需做任何人工配置，安装启动时，根据提示提供ssh服务器连接账户的密码；
* 系统环境配置、内核检测与升级、服务器重启由安装程序自动完成，全程无人干预；
3. 集群服务器操作系统支持丰富，如centos7.x、centos8 stream、rockyLinux9.x、almaLinux9.x；
4. 支持多台（可配）服务器同时安装；
5. 支持已有集群的master服务器扩展、支持node服务器扩展，安全快捷；
6. 当前版本保留20% shell脚本，80%代码 ***GO*** 重写，实现模块化设计、模版化配置、规则引擎解析，能够快速集成k8s最新版本。
7.  ... ...

**ek8资源文件比较大,获取请先安装git-lfs：**

```
Mac: brew install git-lfs
ubuntu: apt-get install git-lfs
centos: yum install git-lfs

安装后，执行git lfs install命令，让当前环境支持全局的LFS配置
```


  
#### （二）使用ek8安装k8s集群前了解
1. ek8安装形式：

<img width="378" alt="image" src="https://github.com/catman002/kubernetes-ek8/assets/35952528/38770bef-3170-47f5-b639-401c56fd9a77">

```安装机是安装了ek8工具的电脑，是作为安装、管理集群环境的操作电脑。```

2.安装流程

<img width="682" alt="image" src="https://github.com/catman002/kubernetes-ek8/assets/35952528/b0f2ff42-87be-49ab-a13d-7580212481a2">



```先配置config.cfg文件，然后执行 ek8 create all 指令，最后行安装指令 ek8 install all(全新安装)```

3. 确认目标安装服务器是amd64架构，操作系统是centos7、centos8 stream、rockyLinux、almaLinux之一；
4. 安装机的系统要求：
* 支持amd64或arm64架构的mac系统；
* 支持amd64架构的linux,如：centos7、centos8 stream、rockyLinux、almaLinux、ubuntu等系统。
5. *****安装机的软件要求：ssh、sshpass、rsync 软件需要提前安装。*****
6. 安装前先配置config.cfg文件

  ***此版本仅作为学习～～交流使用～～，只支持最新2个大版本的最新构建与安装 ！***
  
#### （三）其他需要了解的：
1. 最新组件版本：
* docker:20.10.9
* containerd:1.7.10
* etcd:v3.5.9
* kubernetesv1.29.0
* pause:3.9
* coredns:1.11.1
* calico:3.26.4
* cni:
* haproxy:2.9.1
* keepalived:2.2.8
* crictl:1.29.0
* nerdctl:1.7.1
* openssl:1_1_0j
* buildkit:0.12.4

2. 了解构建内部流程和安装内部流程：

<img width="884" alt="image" src="https://github.com/catman002/kubernetes-ek8/assets/35952528/530c08c8-6624-4ade-9e00-de750e3db2d1">


<img width="932" alt="image" src="https://github.com/catman002/kubernetes-ek8/assets/35952528/4f9fa408-8a15-40e4-a487-b55cd5c32d47">


#### （四）ek8 命令格式: 
```
ek8 command [options] [mods]
```
###### （a）子命令说明 command:
```
create: 依据config.cfg配置，创建集群所需要配置。
install: 执行集群安装或指定模块；
installnew: 对已存在的集群安装扩展服务器；
appendmasters: 扩展一组master节点配置；
appendnodes: 扩展一组node节点配置；
delete: 快速删除集群或模块，但保留etdc、docker、containerd、registry的数据目录；
qdelete: 完全删除集群集群或模块，并且删除etdc、docker、containerd、registry的数据文件不保存；
stop: 停止集群或停止指定模块；
start: 启动集群或指定模块；
 help: ek8命令帮助；
info: 显示ek8版本信息。
```
###### （b）操作说明 options:
```
--exclude: 排除指定模块。
```
###### （c）模块说明 mods:
```
1）与证书的相关模块：caCONFIG、ca、registryCA、etcdCA、frontProxyCA、adminCERT、apiserverCERT、apiserverEtcdClientCERT
 apiserverKubeletClientCERT、bootstrapCERT、controllerManagerCERT、etcdClientCERT、etcdPeerCERT、etcdServerCERT,frontProxyClientCERT、
kubeProxyCERT、registryCERT、saCERT、schedulerCERT、tokenCSV
2）与服务的配置相关的模块：source、apiserverCONFIG、bootstrapCONFIG、calicoCONFIG、controllerManagerCONFIG、corednsCONFIG、dockerCONFIG、etcdCONFIG,haproxyCONFIG、
keepalivedCONFIG、kubeProxyCONFIG、registryCONFIG、schedulerCONFIG、initEnvCONFIG
3)install、installnew、delete、qdelete、stop、start子命令支持的模块：config、soft、deploy、initenv、hostname、cert、docker、containerd、etcd、
kubebin、apiserver、haproxy、keepalived、admin,controller、scheduler、bootstrap、kubeproxy、calico、coredns、dnsutils、registry
4）all 代表所有模块
```
#### （五）ek8命令安装k8s集群例子:
```
1) 全新安装，依据config.cfg配置，创建集群所需文件，包含证书、服务配置、脚本、集群安全等文件，并安装：
./ek8 create all && ./ek8 install all
2) 完全删除集群集群或模块，并且删除etdc、docker、containerd、registry的数据文件不保存：
./ek8 qdelete all
3）快速删除集群或模块，但保留etdc、docker、containerd、registry的数据目录：
./ek8 delete all：
4）新增master服务器，在现有集群环境增加新的master服务器(结合config.cfg的配置)：
./ek8 appendmasters all && ./ek8 installnew all
4）新增node服务器，在现有集群环境增加新的node服务器(结合config.cfg的配置)：
./ek8 appendnodes all && ./ek8 installnew all
2）在在现有集群环境，更新apiserver证书（其他证书类似）：
./ek8 create apiserverCERT && ./ek8 install cert,apiserver
3)集群环境安装，但不安装containerd、haproxy和keepalived(一般维护使用)。
./ek8 create all && ./ek8 install all --exclude=containerd,haproxy,keepalived

... ...

```
#### （六）config.cfg配置说明：
```
#ek8版本（不用改）
ek8_version=v2.0.0.0
#需要安装的k8s版本
kube_version=v1.29.0
k8s_common_ver=${kube_version%.*}.x
sub_conf=cfg/${k8s_common_ver}
#加载子配置：sub.config (version),system.config,cert.config
source $sub_conf
source cfg/base/system.config
source cfg/base/cert.config
#目标服务器系统target_system，目前支持: target_system(centos7/rocky/alma)
#目标服务器cpu类型，target_system：arm64/amd64
target_system=centos7
target_cpu=adm64
#网卡信息,通过ipconfig ,ip a等命令查看，不要配错了
#interface=ens34
interface=ens192
#远程服务器账户、端口(先确认可以通过ssh，使用root远程登录系统)
user=root
ssh_port=22
#kube servers，由服务器名称hostname和ip组成（集群服务器信息:etcd_servers 所有服务器列表;
#apiserver_servers 服务器列表;haproxy_servers 服务器列表;registry_* registry服务器配置信息
etcd_servers=(node1:192.168.0.131 node2:192.168.0.132)
apiserver_servers=(node1:192.168.0.131 node2:192.168.0.132)
node_servers=(node3:192.168.0.133)
haproxy_servers=(node1:192.168.0.131 node2:192.168.0.132)
#镜像库信息，根据情况配置
registry_server_domain=reg.domain.com
registry_server_ip=192.168.0.131
registry_server_port=5000
registry_server_repo_url=/home/docker/registry
registry_username=xxx
registry_passwd=xxx
#容器运行环境 runtime: docker/containerd
#docker_data_dir：docker环境数据目录;#containerd_data_dir：containerd环境数据目录
runtime=containerd
docker_data_dir=/home/var/lib/docker
containerd_data_dir=/home/var/lib/containerd

#需要新增的node(仅在需要新增node时修改;新增后，把新增的nodes合并到¥node_servers中)
#new_node_servers=(node5:192.168.0.135)
new_node_servers=()

#需要新增的master(仅在需要新增master时修改,新增master必须配置新的haproxy；
#新增后，把新增的apiserver合并$apiserver_servers中)
#new_api_servers=(node4:192.168.0.134)
new_api_servers=()

#keepalived 信息
# vip: keepalived虚拟ip.
#如果安装环境是在阿里云、腾讯云等公有云，这个地址必须配成 slb(负载均衡) 地址，同时阿安装增加参数 --exclude=keepalived
VIP=192.168.0.10
netmask_bit=24
mcast_group=224.0.0.22
vrid=192
#haproxy 信息
VIP_PORT=6444
#kube信息
#service_cluster_ip_range and  cluster_cid
cluster_cidr=10.244.0.0/16
service_cluster_ip_range=10.245.0.0/16
#coredns 信息
coredns_cluster_ip=10.245.0.10
#加载子配置： kube.config
source cfg/base/kube.config
```
### （七）V3.0 计划
* 支持ubuntu；
* 图形配置管理；
* GO的覆盖率达到90%以上；
