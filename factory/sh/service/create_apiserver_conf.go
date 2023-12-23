package service

import (
	"ek8/common"
	"ek8/config"
	"strconv"
)

func Create_apiserver_conf(apiServers map[string]string, etcdServers map[string]string) {
	if common.ModFlags["apiserverCONFIG"] != 1 {
		return
	}
	common.TPrintln("创建kube-api-server配置...")

	//拼接 etcd_initial_cluster
	var etcd_servers_str string
	i := 0
	for _, ip := range etcdServers {
		if etcd_servers_str == "" {
			etcd_servers_str = "https://" + ip + ":" + config.Cfg["etcd_port"] //2379
		} else {
			etcd_servers_str = etcd_servers_str + ",https://" + ip + ":" + config.Cfg["etcd_port"] //2379
		}
		i++
	}
	//etcd_servers_str = strings.Replace(etcd_servers_str, ",", "\\,", -1)
	//etcd_servers_str = strings.Replace(etcd_servers_str, "/", "\\/", -1)

	apiserver_count := strconv.Itoa(len(apiServers))
	i = 0
	for _, ip := range apiServers {
		target_dir := config.Cfg["release_dir"] + "/" + ip + "/service/apiserver"

		kube_apiserver_conf_file_name := "kube-apiserver.conf"
		kube_apiserver__service_file_name := "kube-apiserver.service"
		kube_apiserver_conf_file := target_dir + "/" + kube_apiserver_conf_file_name
		kube_apiserver__service_file := target_dir + "/" + kube_apiserver__service_file_name

		template_kube_apiserver_conf_file := "factory/template/service/apiserver/" + config.K8S_VERSION + "/" + kube_apiserver_conf_file_name
		template_kube_apiserver__service_file := "factory/template/service/apiserver/" + kube_apiserver__service_file_name

		apiserver_advertise_address := ip

		//创建 kube-apiserver.service
		new(common.FileAssist).
			Load(template_kube_apiserver__service_file).
			Println("创建 "+kube_apiserver__service_file_name+"...", ip).
			Set("apiserver_work_dir", config.Cfg["apiserver_work_dir"]).
			Set("kubernetes_etc_dir", config.Cfg["kubernetes_etc_dir"]).
			Set("apiserver_exec_start", config.Cfg["apiserver_exec_start"]).
			SaveAs(kube_apiserver__service_file).Close()
		//创建kube-apiserver.conf
		new(common.FileAssist).
			Load(template_kube_apiserver_conf_file).
			Println("创建 "+kube_apiserver_conf_file_name+"...", ip).
			Set("bind_address", ip).
			Set("apiserver_secure_port", config.Cfg["apiserver_secure_port"]).
			Set("token_auth_file", config.Cfg["token_auth_file"]).
			Set("apiserver_advertise_address", apiserver_advertise_address).
			Set("service_cluster_ip_range", config.Cfg["service_cluster_ip_range"]).
			Set("service_node_port_range", config.Cfg["service_node_port_range"]).
			Set("tls_cert_file", config.Cfg["tls_cert_file"]).
			Set("tls_private_key_file", config.Cfg["tls_private_key_file"]).
			Set("client_ca_file", config.Cfg["client_ca_file"]).
			Set("apiserver_client_certificate", config.Cfg["apiserver_client_certificate"]).
			Set("apiserver_client_key", config.Cfg["apiserver_client_key"]).
			Set("service_account_key_file", config.Cfg["service_account_key_file"]).
			Set("service_account_signing_key_file", config.Cfg["service_account_signing_key_file"]).
			Set("etcd_cafile", config.Cfg["etcd_cafile"]).
			Set("etcd_certfile", config.Cfg["etcd_certfile"]).
			Set("etcd_keyfile", config.Cfg["etcd_keyfile"]).
			Set("etcd_servers", etcd_servers_str).
			Set("apiserver_count", apiserver_count).
			Set("requestheader_client_ca_file", config.Cfg["requestheader_client_ca_file"]).
			Set("proxy_client_cert_file", config.Cfg["proxy_client_cert_file"]).
			Set("proxy_client_key_file", config.Cfg["proxy_client_key_file"]).
			Set("service_account_issuer", config.Cfg["service_account_issuer"]).
			Set("api_audiences", config.Cfg["api_audiences"]).
			SaveAs(kube_apiserver_conf_file).Close()

		i++
	}
}
