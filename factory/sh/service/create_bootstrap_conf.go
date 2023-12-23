package service

import (
	"ek8/common"
	"ek8/config"
	"strings"
)

func Create_bootstrap_conf(allservers map[string]string) {
	if common.ModFlags["bootstrapCONFIG"] != 1 {
		return
	}
	common.TPrintln("创建bootstrap配置...")

	createKubeletConf(allservers)

}
func createKubeletConf(allservers map[string]string) {
	cluster_address := "https://" + config.Cfg["VIP"] + ":" + config.Cfg["VIP_PORT"]
	cluster_address = strings.Replace(cluster_address, "/", "\\/", -1)

	for _, ip := range allservers {
		target_dir := config.Cfg["release_dir"] + "/" + ip + "/service/bootstrap"

		kubelet_conf_file_name := "kubelet.yaml"
		kubelet_service_file_name := "kubelet.service"
		delete_kubelet_var_sh_file_name := "delete_kubelet_var.sh"
		kubelet_conf_file := target_dir + "/" + kubelet_conf_file_name
		kubelet_service_file := target_dir + "/" + kubelet_service_file_name
		delete_kubelet_var_sh_file := target_dir + "/" + delete_kubelet_var_sh_file_name

		template_kubelet_conf_file := "factory/template/service/bootstrap/" + config.K8S_VERSION + "/" + kubelet_conf_file_name
		template_kubelet_service_file := "factory/template/service/bootstrap/" + config.K8S_VERSION + "/" + kubelet_service_file_name
		template_delete_kubelet_var_sh_file := "factory/template/service/bootstrap/" + config.K8S_VERSION + "/" + delete_kubelet_var_sh_file_name

		//创建 kubelet.yaml
		new(common.FileAssist).
			Load(template_kubelet_conf_file).
			Println("创建 "+kubelet_conf_file_name+"...", ip).
			Set("kubelet_client_ca_file", config.Cfg["kubelet_client_ca_file"]).
			Set("ip", ip).
			Set("coredns_cluster_ip", config.Cfg["coredns_cluster_ip"]).
			Set("kubelet_static_pod_path", config.Cfg["kubelet_static_pod_path"]).
			SaveAs(kubelet_conf_file).Close()
		//创建 kubelet.service
		new(common.FileAssist).
			Load(template_kubelet_service_file).
			Println("创建 "+kubelet_service_file_name+"...", ip).
			Set("kubelet_work_dir", config.Cfg["kubelet_work_dir"]).
			Set("kubelet_exec_start", config.Cfg["kubelet_exec_start"]).
			Set("kubernetes_etc_dir", config.Cfg["kubernetes_etc_dir"]).
			Set("kubernetes_etc_pki_dir", config.Cfg["kubernetes_etc_pki_dir"]).
			Set("kubelet_dir", config.Cfg["kubelet_dir"]).
			Set("pod_infra_container_image", config.Cfg["pod_infra_container_image"]).
			SaveAs(kubelet_service_file).Close()
		//创建 delete_kubelet_var.sh
		new(common.FileAssist).
			CopyFrom(template_delete_kubelet_var_sh_file).
			Println("创建 "+delete_kubelet_var_sh_file_name+"...", ip).
			CopyTo(delete_kubelet_var_sh_file).Close()
	}

}
