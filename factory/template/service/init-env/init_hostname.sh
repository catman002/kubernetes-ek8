#!/bin/bash
source config.cfg
registry_server="${registry_server_ip} ${registry_server_domain}"
kubernetesVersion=$kube_version

host_default="127.0.0.1  localhost localhost.localdomain localhost4 localhost4.localdomain4  \n
::1 localhost localhost.localdomain localhost6 localhost6.localdomain6"

############################
##update server's  hostname
############################
##echo_s "Configuring server environment ..."
current_server_name={{hostname}}
#for hostport in ${all_servers[*]}
for hostport in {{allServers}}
do
  hostportArray=(${hostport//:/ })
  server_name=${hostportArray[0]}
  server_ip=${hostportArray[1]}

  #init host_content
  if [[ -z $hosts_content ]] ;then
    hosts_content="${server_ip}  ${server_name}"
  else
    hosts_content="${hosts_content}\n${server_ip}  ${server_name}"
  fi
done
hosts_content="${host_default}
$registry_server
$hosts_content"

#修改服务器名
hostnamectl set-hostname $current_server_name
echo -e "$hosts_content" > /etc/hosts
