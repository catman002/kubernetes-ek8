package service

import (
	"ek8/common"
	"ek8/config"
	"strconv"
)

func Create_controller_conf(apiServers map[string]string) {
	if common.ModFlags["controllerManagerCONFIG"] != 1 {
		return
	}
	common.TPrintln("创建kube-controller-manager配置...")

	createControllerConf(apiServers)
}

func createControllerConf(servers map[string]string) {
	for _, ip := range servers {
		target_dir := config.Cfg["release_dir"] + "/" + ip + "/service/controller-manager"

		kube_controller_manager_conf_file_name := "kube-controller-manager.conf"
		kube_controller_manager_service_file_name := "kube-controller-manager.service"
		kube_controller_manager_conf_file := target_dir + "/" + kube_controller_manager_conf_file_name
		kube_controller_manager_service_file := target_dir + "/" + kube_controller_manager_service_file_name

		template_kube_controller_manager_conf_file := "factory/template/service/controller-manager/" + config.K8S_VERSION + "/" + kube_controller_manager_conf_file_name
		template_kube_controller_manager_service_file := "factory/template/service/controller-manager/" + kube_controller_manager_service_file_name

		caExpireTime, _ := strconv.Atoi(config.CADAYS)
		caExpireTimeStr := strconv.Itoa(caExpireTime*24) + "h"

		//创建 kube-controller-manager.service
		new(common.FileAssist).
			Load(template_kube_controller_manager_service_file).
			Println("创建 "+kube_controller_manager_service_file_name+"...", ip).
			Set("controller_work_dir", config.Cfg["controller_work_dir"]).
			Set("kubernetes_etc_dir", config.Cfg["kubernetes_etc_dir"]).
			Set("controller_exec_start", config.Cfg["controller_exec_start"]).
			SaveAs(kube_controller_manager_service_file).Close()
		//创建kube-controller-manager.conf
		new(common.FileAssist).
			Load(template_kube_controller_manager_conf_file).
			Println("创建 "+kube_controller_manager_conf_file_name+"...", ip).
			Set("bind-address", "0.0.0.0").
			Set("caExpireTime", caExpireTimeStr).
			Set("service_cluster_ip_range", config.Cfg["service_cluster_ip_range"]).
			Set("requestheader_client_ca_file", config.Cfg["requestheader_client_ca_file"]).
			Set("cluster_siging_cert_file", config.Cfg["cluster_siging_cert_file"]).
			Set("cluster_siging_key_file", config.Cfg["cluster_siging_key_file"]).
			Set("cluster_cidr", config.Cfg["cluster_cidr"]).
			Set("root_ca_file", config.Cfg["root_ca_file"]).
			Set("client-ca-file", config.Cfg["root_ca_file"]).
			Set("service_account_private_key_file", config.Cfg["service_account_private_key_file"]).
			Set("controller_tls_cert_file", config.Cfg["controller_tls_cert_file"]).
			Set("controller_tls_private_key_file", config.Cfg["controller_tls_private_key_file"]).
			Set("controller_port", config.Cfg["controller_port"]).               // 1.24以下
			Set("controller_secure_port", config.Cfg["controller_secure_port"]). //1.24以下
			Set("kubernetes_etc_dir", config.Cfg["kubernetes_etc_dir"]).         //1.24以下
			Set("cluster_name", config.Cfg["cluster_name"]).                     //1.24以下
			SaveAs(kube_controller_manager_conf_file).Close()
	}
}
