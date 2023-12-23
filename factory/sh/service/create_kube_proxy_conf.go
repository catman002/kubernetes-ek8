package service

import (
	"ek8/common"
	"ek8/config"
)

func Create_kubeproxy_conf(allservers map[string]string) {
	if common.ModFlags["kubeProxyCONFIG"] != 1 {
		return
	}
	common.TPrintln("创建kube-proxy配置...")

	CreatekubeproxyConf(allservers)
}

func CreatekubeproxyConf(servers map[string]string) {
	for _, ip := range servers {
		target_dir := config.Cfg["release_dir"] + "/" + ip + "/service/kube-proxy"

		kube_proxy_yaml_file_name := "kube-proxy.yaml"
		kube_proxy_service_file_name := "kube-proxy.service"
		kube_proxy_yaml_file := target_dir + "/" + kube_proxy_yaml_file_name
		kube_proxy_service_file := target_dir + "/" + kube_proxy_service_file_name

		template_kube_proxy_yaml_file := "factory/template/service/kube-proxy/" + config.K8S_VERSION + "/" + kube_proxy_yaml_file_name
		template_kube_proxy_service_file := "factory/template/service/kube-proxy/" + config.K8S_VERSION + "/" + kube_proxy_service_file_name

		//创建kube-proxy.service
		new(common.FileAssist).
			Load(template_kube_proxy_service_file).
			Println("创建 "+kube_proxy_service_file_name+"...", ip).
			Set("kubeproxy_work_dir", config.Cfg["kubeproxy_work_dir"]).
			Set("kubeproxy_exec_start", config.Cfg["kubeproxy_exec_start"]).
			Set("kubeproxy_dir", config.Cfg["kubeproxy_dir"]).
			SaveAs(kube_proxy_service_file).
			Close()
		//创建kube-proxy.yaml
		new(common.FileAssist).
			Load(template_kube_proxy_yaml_file).
			Println("创建 "+kube_proxy_yaml_file_name+"...", ip).
			Set("ip", ip).
			Set("kubernetes_etc_dir", config.Cfg["kubernetes_etc_dir"]).
			Set("cluster_cidr", config.Cfg["cluster_cidr"]).
			SaveAs(kube_proxy_yaml_file).
			Close()
	}
}
