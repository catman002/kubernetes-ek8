package service

import (
	"ek8/common"
	"ek8/config"
)

func Create_scheduler_conf(apiServers map[string]string) {
	if common.ModFlags["schedulerCONFIG"] != 1 {
		return
	}
	common.TPrintln("创建kube-scheduler配置...")

	createSchedulerConf(apiServers)
}
func createSchedulerConf(servers map[string]string) {
	for _, ip := range servers {
		target_dir := config.Cfg["release_dir"] + "/" + ip + "/service/scheduler"

		kube_scheduler_conf_file_name := "kube-scheduler.conf"
		kube_scheduler_service_file_name := "kube-scheduler.service"
		kube_scheduler_conf_file := target_dir + "/" + kube_scheduler_conf_file_name
		kube_scheduler_service_file := target_dir + "/" + kube_scheduler_service_file_name

		template_kube_scheduler_conf_file := "factory/template/service/scheduler/" + config.K8S_VERSION + "/" + kube_scheduler_conf_file_name
		template_kube_scheduler_service_file := "factory/template/service/scheduler/" + kube_scheduler_service_file_name

		//创建 kube-scheduler.service
		new(common.FileAssist).
			Load(template_kube_scheduler_service_file).
			Println("创建 "+kube_scheduler_service_file_name+"...", ip).
			Set("scheduler_work_dir", config.Cfg["scheduler_work_dir"]).
			Set("kubernetes_etc_dir", config.Cfg["kubernetes_etc_dir"]).
			Set("scheduler_exec_start", config.Cfg["scheduler_exec_start"]).
			SaveAs(kube_scheduler_service_file).
			Close()
		//创建 kube-scheduler.conf
		new(common.FileAssist).
			Load(template_kube_scheduler_conf_file).
			Println("创建 "+kube_scheduler_conf_file_name+"...", ip).
			Set("bind-address", "0.0.0.0").
			Set("kubernetes_etc_dir", config.Cfg["kubernetes_etc_dir"]).
			Set("client-ca-file", config.Cfg["root_ca_file"]).
			Set("tls-cert-file", config.Cfg["scheduler_client_certificate"]).
			Set("tls-private-key-file", config.Cfg["scheduler_client_key"]).
			SaveAs(kube_scheduler_conf_file).
			Close()
	}
}
